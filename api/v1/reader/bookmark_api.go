package reader

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	readerModels "Qingyu_backend/models/reader"
	"Qingyu_backend/service/interfaces"
	readerservice "Qingyu_backend/service/reader"
	"Qingyu_backend/pkg/response"
)

// BookmarkAPI 书签API
type BookmarkAPI struct {
	bookmarkService interfaces.BookmarkService
}

// NewBookmarkAPI 创建书签API实例
func NewBookmarkAPI(bookmarkService interfaces.BookmarkService) *BookmarkAPI {
	return &BookmarkAPI{
		bookmarkService: bookmarkService,
	}
}

// CreateBookmark 创建书签
//
//	@Summary		创建书签
//	@Description	在指定章节位置创建书签
//	@Tags			Reader-Bookmark
//	@Accept			json
//	@Produce		json
//	@Param		bookId	path	string	true	"书籍ID"
//	@Param			request	body	CreateBookmarkRequest	true	"书签信息"
//	@Success		201		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Failure		409		{object}	response.APIResponse	"书签已存在"
//	@Router			/api/v1/reader/books/{bookId}/bookmarks [post]
func (api *BookmarkAPI) CreateBookmark(c *gin.Context) {
	var req CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 解析ID
	bookmark, err := req.ToBookmark(userID.(string))
	if err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 创建书签
	if err := api.bookmarkService.CreateBookmark(c.Request.Context(), bookmark); err != nil {
		if err == readerservice.ErrBookmarkAlreadyExists {
			response.Conflict(c, "书签已存在", err.Error())
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Created(c, bookmark)
}

// GetBookmarks 获取书签列表
//
//	@Summary		获取书签列表
//	@Description	获取用户的书签列表，支持筛选和分页
//	@Tags			Reader-Bookmark
//	@Accept			json
//	@Produce		json
//	@Param		bookId	path	string	true	"书籍ID"
//	@Param			color		query	string	false	"颜色筛选"
//	@Param			tag			query	string	false	"标签筛选"
//	@Param			isPublic	query	bool	false	"是否公开"
//	@Param			page		query	int	false	"页码"	default(1)
//	@Param			size		query	int	false	"每页数量"	default(20)
//	@Success		200		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Router			/api/v1/reader/bookmarks [get]
//	@Router			/api/v1/reader/books/{bookId}/bookmarks [get]
func (api *BookmarkAPI) GetBookmarks(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// 构建筛选条件
	filter := &readerModels.BookmarkFilter{
		Color:    c.Query("color"),
		Tag:      c.Query("tag"),
		IsPublic: nil,
	}

	if isPublic := c.Query("isPublic"); isPublic != "" {
		if isPublic == "true" {
			trueVal := true
			filter.IsPublic = &trueVal
		} else {
			falseVal := false
			filter.IsPublic = &falseVal
		}
	}

	// 检查是否是获取某本书的书签
	bookID := c.Param("bookId")
	var result interface{}
	var err error

	if bookID != "" {
		result, err = api.bookmarkService.GetBookBookmarks(c.Request.Context(), userID.(string), bookID, page, size)
	} else {
		result, err = api.bookmarkService.GetUserBookmarks(c.Request.Context(), userID.(string), filter, page, size)
	}

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, result)
}

// GetBookmark 获取书签详情
//
//	@Summary		获取书签详情
//	@Description	获取单个书签的详细信息
//	@Tags			Reader-Bookmark
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"书签ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		401	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Router			/api/v1/reader/bookmarks/{id} [get]
func (api *BookmarkAPI) GetBookmark(c *gin.Context) {
	bookmarkID := c.Param("id")

	bookmark, err := api.bookmarkService.GetBookmark(c.Request.Context(), bookmarkID)
	if err != nil {
		if err == readerservice.ErrBookmarkNotFound {
			response.NotFound(c, "书签不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, bookmark)
}

// UpdateBookmark 更新书签
//
//	@Summary		更新书签
//	@Description	更新书签信息
//	@Tags			Reader-Bookmark
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string				true	"书签ID"
//	@Param			request	body	UpdateBookmarkRequest	true	"更新信息"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Failure		404		{object}	response.APIResponse
//	@Router			/api/v1/reader/bookmarks/{id} [put]
func (api *BookmarkAPI) UpdateBookmark(c *gin.Context) {
	bookmarkID := c.Param("id")

	var req UpdateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	bookmark := &readerModels.Bookmark{
		Note:     req.Note,
		Color:    req.Color,
		Quote:    req.Quote,
		IsPublic: req.IsPublic,
		Tags:     req.Tags,
	}

	if err := api.bookmarkService.UpdateBookmark(c.Request.Context(), bookmarkID, bookmark); err != nil {
		if err == readerservice.ErrBookmarkNotFound {
			response.NotFound(c, "书签不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// DeleteBookmark 删除书签
//
//	@Summary		删除书签
//	@Description	删除指定书签
//	@Tags			Reader-Bookmark
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"书签ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		401	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Router			/api/v1/reader/bookmarks/{id} [delete]
func (api *BookmarkAPI) DeleteBookmark(c *gin.Context) {
	bookmarkID := c.Param("id")

	if err := api.bookmarkService.DeleteBookmark(c.Request.Context(), bookmarkID); err != nil {
		if err == readerservice.ErrBookmarkNotFound {
			response.NotFound(c, "书签不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// ExportBookmarks 导出书签
//
//	@Summary		导出书签
//	@Description	导出用户的所有书签
//	@Tags			Reader-Bookmark
//	@Accept			json
//	@Produce		json
//	@Param			format	query	string	true	"导出格式"	Enums(json,csv)
//	@Success		200		{file}		file
//	@Failure		401		{object}	response.APIResponse
//	@Router			/api/v1/reader/bookmarks/export [get]
func (api *BookmarkAPI) ExportBookmarks(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	format := c.DefaultQuery("format", "json")
	data, contentType, err := api.bookmarkService.ExportBookmarks(c.Request.Context(), userID.(string), format)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 设置响应头
	filename := "bookmarks_" + format
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}

// GetBookmarkStats 获取书签统计
//
//	@Summary		书签统计
//	@Description	获取用户的书签统计信息
//	@Tags			Reader-Bookmark
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.APIResponse
//	@Failure		401	{object}	response.APIResponse
//	@Router			/api/v1/reader/bookmarks/stats [get]
func (api *BookmarkAPI) GetBookmarkStats(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	stats, err := api.bookmarkService.GetBookmarkStats(c.Request.Context(), userID.(string))
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, stats)
}

// SearchBookmarks 搜索书签
//
//	@Summary		搜索书签
//	@Description	在笔记、引用、标签中搜索关键词
//	@Tags			Reader-Bookmark
//	@Accept			json
//	@Produce		json
//	@Param			keyword	query	string	true	"搜索关键词"
//	@Param			page		query	int	false	"页码"	default(1)
//	@Param			size		query	int	false	"每页数量"	default(20)
//	@Success		200		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Router			/api/v1/reader/bookmarks/search [get]
func (api *BookmarkAPI) SearchBookmarks(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "参数错误", "关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	result, err := api.bookmarkService.SearchBookmarks(c.Request.Context(), userID.(string), keyword, page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, result)
}

// CreateBookmarkRequest 创建书签请求
type CreateBookmarkRequest struct {
	BookID    string   `json:"bookId" binding:"required"`
	ChapterID string   `json:"chapterId" binding:"required"`
	Position  int      `json:"position" binding:"required,min=0"`
	Note      string   `json:"note"`
	Color     string   `json:"color"`
	Quote     string   `json:"quote"`
	IsPublic  bool     `json:"isPublic"`
	Tags      []string `json:"tags"`
}

// ToBookmark 转换为书签模型
func (r *CreateBookmarkRequest) ToBookmark(userID string) (*readerModels.Bookmark, error) {
	bookID, err := primitive.ObjectIDFromHex(r.BookID)
	if err != nil {
		return nil, err
	}

	chapterID, err := primitive.ObjectIDFromHex(r.ChapterID)
	if err != nil {
		return nil, err
	}

	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	return &readerModels.Bookmark{
		UserID:    userOID,
		BookID:    bookID,
		ChapterID: chapterID,
		Position:  r.Position,
		Note:      r.Note,
		Color:     r.Color,
		Quote:     r.Quote,
		IsPublic:  r.IsPublic,
		Tags:      r.Tags,
	}, nil
}

// UpdateBookmarkRequest 更新书签请求
type UpdateBookmarkRequest struct {
	Note     string   `json:"note"`
	Color    string   `json:"color"`
	Quote    string   `json:"quote"`
	IsPublic bool     `json:"isPublic"`
	Tags     []string `json:"tags"`
}
