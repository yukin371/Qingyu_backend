package booklist

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/social"
)

// BookListAPI 书单API处理器
type BookListAPI struct {
	bookListService *social.BookListService
}

// NewBookListAPI 创建书单API实例
func NewBookListAPI(bookListService *social.BookListService) *BookListAPI {
	return &BookListAPI{
		bookListService: bookListService,
	}
}

// CreateBookListRequest 创建书单请求
type CreateBookListRequest struct {
	Title       string   `json:"title" binding:"required,max=100"`
	Description string   `json:"description" binding:"max=500"`
	Cover       string   `json:"cover"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	IsPublic    bool     `json:"is_public"`
}

// CreateBookList 创建书单
// @Summary 创建书单
// @Tags 书单管理
// @Accept json
// @Produce json
// @Param request body CreateBookListRequest true "书单信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/booklists [post]
// @Security Bearer
func (api *BookListAPI) CreateBookList(c *gin.Context) {
	var req CreateBookListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 获取用户信息
	userName := ""
	userAvatar := ""
	if name, ok := c.Get("username"); ok {
		userName = name.(string)
	}

	bookList, err := api.bookListService.CreateBookList(
		c.Request.Context(),
		userID.(string),
		userName,
		userAvatar,
		req.Title,
		req.Description,
		req.Cover,
		req.Category,
		req.Tags,
		req.IsPublic,
	)

	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "创建书单失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "创建书单成功", bookList)
}

// GetBookLists 获取书单列表
// @Summary 获取书单列表
// @Tags 书单管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/booklists [get]
func (api *BookListAPI) GetBookLists(c *gin.Context) {
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

	bookLists, total, err := api.bookListService.GetBookLists(
		c.Request.Context(),
		params.Page,
		params.Size,
	)

	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取书单列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取书单列表成功", gin.H{
		"list":  bookLists,
		"total": total,
		"page":  params.Page,
		"size":  params.Size,
	})
}

// GetBookListDetail 获取书单详情
// @Summary 获取书单详情
// @Tags 书单管理
// @Accept json
// @Produce json
// @Param id path string true "书单ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/booklists/{id} [get]
func (api *BookListAPI) GetBookListDetail(c *gin.Context) {
	bookListID := c.Param("id")
	if bookListID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书单ID不能为空")
		return
	}

	bookList, err := api.bookListService.GetBookListByID(c.Request.Context(), bookListID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取书单详情失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取书单详情成功", bookList)
}

// UpdateBookListRequest 更新书单请求
type UpdateBookListRequest struct {
	Title       *string   `json:"title" binding:"omitempty,max=100"`
	Description *string   `json:"description" binding:"omitempty,max=500"`
	Cover       *string   `json:"cover"`
	Category    *string   `json:"category"`
	Tags        *[]string `json:"tags"`
	IsPublic    *bool     `json:"is_public"`
}

// UpdateBookList 更新书单
// @Summary 更新书单
// @Tags 书单管理
// @Accept json
// @Produce json
// @Param id path string true "书单ID"
// @Param request body UpdateBookListRequest true "更新信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/booklists/{id} [put]
// @Security Bearer
func (api *BookListAPI) UpdateBookList(c *gin.Context) {
	bookListID := c.Param("id")
	if bookListID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书单ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	var req UpdateBookListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Cover != nil {
		updates["cover"] = *req.Cover
	}
	if req.Category != nil {
		updates["category"] = *req.Category
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

	err := api.bookListService.UpdateBookList(c.Request.Context(), userID.(string), bookListID, updates)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新书单失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新书单成功", nil)
}

// DeleteBookList 删除书单
// @Summary 删除书单
// @Tags 书单管理
// @Accept json
// @Produce json
// @Param id path string true "书单ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/booklists/{id} [delete]
// @Security Bearer
func (api *BookListAPI) DeleteBookList(c *gin.Context) {
	bookListID := c.Param("id")
	if bookListID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书单ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	err := api.bookListService.DeleteBookList(c.Request.Context(), userID.(string), bookListID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "删除书单失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除书单成功", nil)
}

// LikeBookList 点赞书单
// @Summary 点赞书单
// @Tags 书单管理
// @Accept json
// @Produce json
// @Param id path string true "书单ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/booklists/{id}/like [post]
// @Security Bearer
func (api *BookListAPI) LikeBookList(c *gin.Context) {
	bookListID := c.Param("id")
	if bookListID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书单ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	err := api.bookListService.LikeBookList(c.Request.Context(), userID.(string), bookListID)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "已经点赞过该书单" {
			shared.Error(c, http.StatusBadRequest, "操作失败", errMsg)
		} else {
			shared.Error(c, http.StatusInternalServerError, "点赞失败", errMsg)
		}
		return
	}

	shared.Success(c, http.StatusOK, "点赞成功", nil)
}

// ForkBookList 复制书单
// @Summary 复制书单
// @Tags 书单管理
// @Accept json
// @Produce json
// @Param id path string true "书单ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/booklists/{id}/fork [post]
// @Security Bearer
func (api *BookListAPI) ForkBookList(c *gin.Context) {
	bookListID := c.Param("id")
	if bookListID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书单ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	bookList, err := api.bookListService.ForkBookList(c.Request.Context(), userID.(string), bookListID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "复制书单失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "复制书单成功", bookList)
}

// GetBooksInList 获取书单中的书籍
// @Summary 获取书单中的书籍
// @Tags 书单管理
// @Accept json
// @Produce json
// @Param id path string true "书单ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/booklists/{id}/books [get]
func (api *BookListAPI) GetBooksInList(c *gin.Context) {
	bookListID := c.Param("id")
	if bookListID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书单ID不能为空")
		return
	}

	books, err := api.bookListService.GetBooksInList(c.Request.Context(), bookListID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取书单书籍失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取书单书籍成功", gin.H{
		"list": books,
	})
}
