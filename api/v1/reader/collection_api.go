package reader

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/reading"
)

// CollectionAPI 收藏API处理器
type CollectionAPI struct {
	collectionService *reading.CollectionService
}

// NewCollectionAPI 创建收藏API实例
func NewCollectionAPI(collectionService *reading.CollectionService) *CollectionAPI {
	return &CollectionAPI{
		collectionService: collectionService,
	}
}

// =========================
// 收藏管理
// =========================

// AddCollectionRequest 添加收藏请求
type AddCollectionRequest struct {
	BookID   string   `json:"book_id" binding:"required"`
	FolderID string   `json:"folder_id"`
	Note     string   `json:"note" binding:"max=500"`
	Tags     []string `json:"tags"`
	IsPublic bool     `json:"is_public"`
}

// AddCollection 添加收藏
// @Summary 添加收藏
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Param request body AddCollectionRequest true "收藏信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections [post]
// @Security Bearer
func (api *CollectionAPI) AddCollection(c *gin.Context) {
	var req AddCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	collection, err := api.collectionService.AddToCollection(
		c.Request.Context(),
		userID.(string),
		req.BookID,
		req.FolderID,
		req.Note,
		req.Tags,
		req.IsPublic,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "添加收藏失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "添加收藏成功", collection)
}

// GetCollections 获取收藏列表
// @Summary 获取收藏列表
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Param folder_id query string false "收藏夹ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections [get]
// @Security Bearer
func (api *CollectionAPI) GetCollections(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	folderID := c.Query("folder_id")

	var params struct {
		Page int `form:"page" binding:"min=1"`
		Size int `form:"size" binding:"min=1,max=100"`
	}
	params.Page = 1
	params.Size = 20

	if err := c.ShouldBindQuery(&params); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	collections, total, err := api.collectionService.GetUserCollections(
		c.Request.Context(),
		userID.(string),
		folderID,
		params.Page,
		params.Size,
	)

	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取收藏列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取收藏列表成功", gin.H{
		"list":  collections,
		"total": total,
		"page":  params.Page,
		"size":  params.Size,
	})
}

// UpdateCollectionRequest 更新收藏请求
type UpdateCollectionRequest struct {
	FolderID *string   `json:"folder_id"`
	Note     *string   `json:"note" binding:"omitempty,max=500"`
	Tags     *[]string `json:"tags"`
	IsPublic *bool     `json:"is_public"`
}

// UpdateCollection 更新收藏
// @Summary 更新收藏
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Param id path string true "收藏ID"
// @Param request body UpdateCollectionRequest true "更新信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections/{id} [put]
// @Security Bearer
func (api *CollectionAPI) UpdateCollection(c *gin.Context) {
	collectionID := c.Param("id")
	if collectionID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "收藏ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	var req UpdateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.FolderID != nil {
		updates["folder_id"] = *req.FolderID
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}
	if req.Tags != nil {
		updates["tags"] = *req.Tags
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	if len(updates) == 0 {
		shared.Error(c, http.StatusBadRequest, "参数错误", "没有要更新的字段")
		return
	}

	err := api.collectionService.UpdateCollection(
		c.Request.Context(),
		userID.(string),
		collectionID,
		updates,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "更新收藏失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新收藏成功", nil)
}

// DeleteCollection 删除收藏
// @Summary 删除收藏
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Param id path string true "收藏ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections/{id} [delete]
// @Security Bearer
func (api *CollectionAPI) DeleteCollection(c *gin.Context) {
	collectionID := c.Param("id")
	if collectionID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "收藏ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	err := api.collectionService.RemoveFromCollection(
		c.Request.Context(),
		userID.(string),
		collectionID,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "删除收藏失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除收藏成功", nil)
}

// CheckCollected 检查是否已收藏
// @Summary 检查是否已收藏
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Param book_id path string true "书籍ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections/check/{book_id} [get]
// @Security Bearer
func (api *CollectionAPI) CheckCollected(c *gin.Context) {
	bookID := c.Param("book_id")
	if bookID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	isCollected, err := api.collectionService.IsCollected(
		c.Request.Context(),
		userID.(string),
		bookID,
	)

	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "检查收藏状态失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "检查收藏状态成功", gin.H{
		"is_collected": isCollected,
	})
}

// GetCollectionsByTag 根据标签获取收藏
// @Summary 根据标签获取收藏
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Param tag path string true "标签"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections/tags/{tag} [get]
// @Security Bearer
func (api *CollectionAPI) GetCollectionsByTag(c *gin.Context) {
	tag := c.Param("tag")
	if tag == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "标签不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	var params struct {
		Page int `form:"page" binding:"min=1"`
		Size int `form:"size" binding:"min=1,max=100"`
	}
	params.Page = 1
	params.Size = 20

	if err := c.ShouldBindQuery(&params); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	collections, total, err := api.collectionService.GetCollectionsByTag(
		c.Request.Context(),
		userID.(string),
		tag,
		params.Page,
		params.Size,
	)

	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取收藏列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取收藏列表成功", gin.H{
		"list":  collections,
		"total": total,
		"page":  params.Page,
		"size":  params.Size,
	})
}

// =========================
// 收藏夹管理
// =========================

// CreateFolderRequest 创建收藏夹请求
type CreateFolderRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"max=200"`
	IsPublic    bool   `json:"is_public"`
}

// CreateFolder 创建收藏夹
// @Summary 创建收藏夹
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Param request body CreateFolderRequest true "收藏夹信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections/folders [post]
// @Security Bearer
func (api *CollectionAPI) CreateFolder(c *gin.Context) {
	var req CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	folder, err := api.collectionService.CreateFolder(
		c.Request.Context(),
		userID.(string),
		req.Name,
		req.Description,
		req.IsPublic,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "创建收藏夹失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "创建收藏夹成功", folder)
}

// GetFolders 获取收藏夹列表
// @Summary 获取收藏夹列表
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections/folders [get]
// @Security Bearer
func (api *CollectionAPI) GetFolders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	folders, err := api.collectionService.GetUserFolders(
		c.Request.Context(),
		userID.(string),
	)

	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取收藏夹列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取收藏夹列表成功", gin.H{
		"list": folders,
	})
}

// UpdateFolderRequest 更新收藏夹请求
type UpdateFolderRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=50"`
	Description *string `json:"description" binding:"omitempty,max=200"`
	IsPublic    *bool   `json:"is_public"`
}

// UpdateFolder 更新收藏夹
// @Summary 更新收藏夹
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Param id path string true "收藏夹ID"
// @Param request body UpdateFolderRequest true "更新信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections/folders/{id} [put]
// @Security Bearer
func (api *CollectionAPI) UpdateFolder(c *gin.Context) {
	folderID := c.Param("id")
	if folderID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "收藏夹ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	var req UpdateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	if len(updates) == 0 {
		shared.Error(c, http.StatusBadRequest, "参数错误", "没有要更新的字段")
		return
	}

	err := api.collectionService.UpdateFolder(
		c.Request.Context(),
		userID.(string),
		folderID,
		updates,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "更新收藏夹失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新收藏夹成功", nil)
}

// DeleteFolder 删除收藏夹
// @Summary 删除收藏夹
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Param id path string true "收藏夹ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections/folders/{id} [delete]
// @Security Bearer
func (api *CollectionAPI) DeleteFolder(c *gin.Context) {
	folderID := c.Param("id")
	if folderID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "收藏夹ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	err := api.collectionService.DeleteFolder(
		c.Request.Context(),
		userID.(string),
		folderID,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "删除收藏夹失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除收藏夹成功", nil)
}

// =========================
// 收藏分享
// =========================

// ShareCollection 分享收藏
// @Summary 分享收藏
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Param id path string true "收藏ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections/{id}/share [post]
// @Security Bearer
func (api *CollectionAPI) ShareCollection(c *gin.Context) {
	collectionID := c.Param("id")
	if collectionID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "收藏ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	err := api.collectionService.ShareCollection(
		c.Request.Context(),
		userID.(string),
		collectionID,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "分享收藏失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "分享收藏成功", nil)
}

// UnshareCollection 取消分享收藏
// @Summary 取消分享收藏
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Param id path string true "收藏ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections/{id}/share [delete]
// @Security Bearer
func (api *CollectionAPI) UnshareCollection(c *gin.Context) {
	collectionID := c.Param("id")
	if collectionID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "收藏ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	err := api.collectionService.UnshareCollection(
		c.Request.Context(),
		userID.(string),
		collectionID,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "取消分享失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "取消分享成功", nil)
}

// GetPublicCollections 获取公开收藏列表
// @Summary 获取公开收藏列表
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections/public [get]
func (api *CollectionAPI) GetPublicCollections(c *gin.Context) {
	var params struct {
		Page int `form:"page" binding:"min=1"`
		Size int `form:"size" binding:"min=1,max=100"`
	}
	params.Page = 1
	params.Size = 20

	if err := c.ShouldBindQuery(&params); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	collections, total, err := api.collectionService.GetPublicCollections(
		c.Request.Context(),
		params.Page,
		params.Size,
	)

	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取公开收藏列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取公开收藏列表成功", gin.H{
		"list":  collections,
		"total": total,
		"page":  params.Page,
		"size":  params.Size,
	})
}

// =========================
// 统计
// =========================

// GetCollectionStats 获取收藏统计
// @Summary 获取收藏统计
// @Tags 阅读端-收藏
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reader/collections/stats [get]
// @Security Bearer
func (api *CollectionAPI) GetCollectionStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	stats, err := api.collectionService.GetUserCollectionStats(
		c.Request.Context(),
		userID.(string),
	)

	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取收藏统计失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取收藏统计成功", stats)
}
