package social

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/interfaces"
)

// BookListAPI 书单API处理器
type BookListAPI struct {
	bookListService interfaces.BookListService
}

// NewBookListAPI 创建书单API实例
func NewBookListAPI(bookListService interfaces.BookListService) *BookListAPI {
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
// @Tags 社交-书单
// @Accept json
// @Produce json
// @Param request body CreateBookListRequest true "书单信息"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/booklists [post]
// @Security Bearer
func (api *BookListAPI) CreateBookList(c *gin.Context) {
	var req CreateBookListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	// 获取用户信息
	userName := ""
	userAvatar := ""
	if name, exists := c.Get("username"); exists {
		userName = name.(string)
	}

	bookList, err := api.bookListService.CreateBookList(
		c.Request.Context(),
		userID,
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
		response.InternalError(c, err)
		return
	}

	response.Created(c, bookList)
}

// GetBookLists 获取书单列表
// @Summary 获取书单列表
// @Tags 社交-书单
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/booklists [get]
func (api *BookListAPI) GetBookLists(c *gin.Context) {
	pagination := shared.GetPaginationParamsStandard(c)

	bookLists, total, err := api.bookListService.GetBookLists(
		c.Request.Context(),
		pagination.Page,
		pagination.PageSize,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.RespondWithPaginated(c, bookLists, int(total), pagination.Page, pagination.PageSize, "")
}

// GetBookListDetail 获取书单详情
// @Summary 获取书单详情
// @Tags 社交-书单
// @Accept json
// @Produce json
// @Param id path string true "书单ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/booklists/{id} [get]
func (api *BookListAPI) GetBookListDetail(c *gin.Context) {
	bookListID, ok := shared.GetRequiredParam(c, "id", "书单ID")
	if !ok {
		return
	}

	bookList, err := api.bookListService.GetBookListByID(c.Request.Context(), bookListID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, bookList)
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
// @Tags 社交-书单
// @Accept json
// @Produce json
// @Param id path string true "书单ID"
// @Param request body UpdateBookListRequest true "更新信息"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/booklists/{id} [put]
// @Security Bearer
func (api *BookListAPI) UpdateBookList(c *gin.Context) {
	bookListID, ok := shared.GetRequiredParam(c, "id", "书单ID")
	if !ok {
		return
	}

	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	var req UpdateBookListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
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
		response.BadRequest(c, "参数错误", "没有要更新的字段")
		return
	}

	err := api.bookListService.UpdateBookList(c.Request.Context(), userID, bookListID, updates)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// DeleteBookList 删除书单
// @Summary 删除书单
// @Tags 社交-书单
// @Accept json
// @Produce json
// @Param id path string true "书单ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/booklists/{id} [delete]
// @Security Bearer
func (api *BookListAPI) DeleteBookList(c *gin.Context) {
	bookListID, ok := shared.GetRequiredParam(c, "id", "书单ID")
	if !ok {
		return
	}

	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	err := api.bookListService.DeleteBookList(c.Request.Context(), userID, bookListID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// LikeBookList 点赞书单
// @Summary 点赞书单
// @Tags 社交-书单
// @Accept json
// @Produce json
// @Param id path string true "书单ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/booklists/{id}/like [post]
// @Security Bearer
func (api *BookListAPI) LikeBookList(c *gin.Context) {
	bookListID, ok := shared.GetRequiredParam(c, "id", "书单ID")
	if !ok {
		return
	}

	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	err := api.bookListService.LikeBookList(c.Request.Context(), userID, bookListID)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "已经点赞过该书单" {
			response.BadRequest(c, "操作失败", errMsg)
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, nil)
}

// ForkBookList 复制书单
// @Summary 复制书单
// @Tags 社交-书单
// @Accept json
// @Produce json
// @Param id path string true "书单ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/booklists/{id}/fork [post]
// @Security Bearer
func (api *BookListAPI) ForkBookList(c *gin.Context) {
	bookListID, ok := shared.GetRequiredParam(c, "id", "书单ID")
	if !ok {
		return
	}

	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	bookList, err := api.bookListService.ForkBookList(c.Request.Context(), userID, bookListID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, bookList)
}

// GetBooksInList 获取书单中的书籍
// @Summary 获取书单中的书籍
// @Tags 社交-书单
// @Accept json
// @Produce json
// @Param id path string true "书单ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/booklists/{id}/books [get]
func (api *BookListAPI) GetBooksInList(c *gin.Context) {
	bookListID, ok := shared.GetRequiredParam(c, "id", "书单ID")
	if !ok {
		return
	}

	books, err := api.bookListService.GetBooksInList(c.Request.Context(), bookListID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"list": books,
	})
}
