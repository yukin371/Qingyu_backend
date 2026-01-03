package bookstore

import (
	"Qingyu_backend/service/bookstore"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChapterCatalogAPI 章节目录和购买API处理器
type ChapterCatalogAPI struct {
	chapterService  bookstore.ChapterService
	purchaseService bookstore.ChapterPurchaseService
}

// NewChapterCatalogAPI 创建章节目录API实例
func NewChapterCatalogAPI(
	chapterService bookstore.ChapterService,
	purchaseService bookstore.ChapterPurchaseService,
) *ChapterCatalogAPI {
	return &ChapterCatalogAPI{
		chapterService:  chapterService,
		purchaseService: purchaseService,
	}
}

// GetChapterCatalog 获取书籍章节目录
//
//	@Summary		获取书籍章节目录
//	@Description	获取指定书籍的章节目录，返回树形结构的章节列表
//	@Tags			章节目录
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		404	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/bookstore/books/{id}/chapters [get]
func (api *ChapterCatalogAPI) GetChapterCatalog(c *gin.Context) {
	bookIDStr := c.Param("id")
	if bookIDStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "书籍ID不能为空",
		})
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的书籍ID格式",
		})
		return
	}

	// 获取用户ID（从中间件设置的上下文中）
	var userID primitive.ObjectID
	if userIDValue, exists := c.Get("userId"); exists {
		if uid, ok := userIDValue.(string); ok {
			userID, _ = primitive.ObjectIDFromHex(uid)
		}
	}

	catalog, err := api.purchaseService.GetChapterCatalog(c.Request.Context(), userID, bookID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, APIResponse{
				Code:    404,
				Message: "书籍不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取章节目录失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    catalog,
	})
}

// GetChapterInfo 获取单个章节信息
//
//	@Summary		获取单个章节信息
//	@Description	根据章节ID获取章节详细信息
//	@Tags			章节目录
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"书籍ID"
//	@Param			chapterId	path		string	true	"章节ID"
//	@Success		200			{object}	APIResponse
//	@Failure		400			{object}	APIResponse
//	@Failure		404			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/bookstore/books/{id}/chapters/{chapterId} [get]
func (api *ChapterCatalogAPI) GetChapterInfo(c *gin.Context) {
	chapterIdStr := c.Param("chapterId")
	if chapterIdStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "章节ID不能为空",
		})
		return
	}

	chapterID, err := primitive.ObjectIDFromHex(chapterIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的章节ID格式",
		})
		return
	}

	chapter, err := api.chapterService.GetChapterByID(c.Request.Context(), chapterID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, APIResponse{
				Code:    404,
				Message: "章节不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取章节信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    chapter,
	})
}

// GetTrialChapters 获取试读章节列表
//
//	@Summary		获取试读章节列表
//	@Description	获取书籍的试读章节列表（通常前N章免费）
//	@Tags			章节目录
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"书籍ID"
//	@Param			count	query		int		false	"试读章节数量"	default(10)
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/bookstore/books/{id}/trial-chapters [get]
func (api *ChapterCatalogAPI) GetTrialChapters(c *gin.Context) {
	bookIDStr := c.Param("id")
	if bookIDStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "书籍ID不能为空",
		})
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的书籍ID格式",
		})
		return
	}

	trialCount, _ := strconv.Atoi(c.DefaultQuery("count", "10"))
	if trialCount <= 0 {
		trialCount = 10
	}

	chapters, err := api.purchaseService.GetTrialChapters(c.Request.Context(), bookID, trialCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取试读章节失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"book_id":   bookID.Hex(),
			"count":     len(chapters),
			"chapters":  chapters,
		},
	})
}

// GetVIPChapters 获取VIP章节列表
//
//	@Summary		获取VIP章节列表
//	@Description	获取书籍的VIP章节列表（需要VIP权限）
//	@Tags			章节目录
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/bookstore/books/{id}/vip-chapters [get]
func (api *ChapterCatalogAPI) GetVIPChapters(c *gin.Context) {
	bookIDStr := c.Param("id")
	if bookIDStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "书籍ID不能为空",
		})
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的书籍ID格式",
		})
		return
	}

	chapters, err := api.purchaseService.GetVIPChapters(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取VIP章节失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"book_id":   bookID.Hex(),
			"count":     len(chapters),
			"chapters":  chapters,
		},
	})
}

// GetChapterPrice 获取章节价格
//
//	@Summary		获取章节价格
//	@Description	获取指定章节的价格信息
//	@Tags			章节购买
//	@Accept			json
//	@Produce		json
//	@Param			chapterId	path		string	true	"章节ID"
//	@Success		200			{object}	APIResponse
//	@Failure		400			{object}	APIResponse
//	@Failure		404			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/bookstore/chapters/{chapterId}/price [get]
func (api *ChapterCatalogAPI) GetChapterPrice(c *gin.Context) {
	chapterIdStr := c.Param("chapterId")
	if chapterIdStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "章节ID不能为空",
		})
		return
	}

	chapterID, err := primitive.ObjectIDFromHex(chapterIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的章节ID格式",
		})
		return
	}

	price, err := api.purchaseService.GetChapterPrice(c.Request.Context(), chapterID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, APIResponse{
				Code:    404,
				Message: "章节不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取章节价格失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"chapter_id": chapterID.Hex(),
			"price":      price,
		},
	})
}

// PurchaseChapter 购买章节
//
//	@Summary		购买章节
//	@Description	购买单个付费章节
//	@Tags			章节购买
//	@Accept			json
//	@Produce		json
//	@Param			chapterId	path		string	true	"章节ID"
//	@Success		200			{object}	APIResponse
//	@Failure		400			{object}	APIResponse
//	@Failure		403			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/reader/chapters/{chapterId}/purchase [post]
func (api *ChapterCatalogAPI) PurchaseChapter(c *gin.Context) {
	chapterIdStr := c.Param("chapterId")
	if chapterIdStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "章节ID不能为空",
		})
		return
	}

	chapterID, err := primitive.ObjectIDFromHex(chapterIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的章节ID格式",
		})
		return
	}

	// 获取用户ID
	userIDValue, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未授权访问",
		})
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "无效的用户信息",
		})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的用户ID格式",
		})
		return
	}

	purchase, err := api.purchaseService.PurchaseChapter(c.Request.Context(), userID, chapterID)
	if err != nil {
		if strings.Contains(err.Error(), "already purchased") {
			c.JSON(http.StatusConflict, APIResponse{
				Code:    409,
				Message: "章节已购买",
			})
			return
		}
		if strings.Contains(err.Error(), "insufficient balance") {
			c.JSON(http.StatusForbidden, APIResponse{
				Code:    403,
				Message: "余额不足",
			})
			return
		}
		if strings.Contains(err.Error(), "free chapter") {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:    400,
				Message: "免费章节无需购买",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "购买章节失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "购买成功",
		Data:    purchase,
	})
}

// PurchaseBook 购买全书
//
//	@Summary		购买全书
//	@Description	批量购买书籍的所有付费章节（享受折扣）
//	@Tags			章节购买
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		403	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/reader/books/{id}/buy-all [post]
func (api *ChapterCatalogAPI) PurchaseBook(c *gin.Context) {
	bookIDStr := c.Param("id")
	if bookIDStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "书籍ID不能为空",
		})
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的书籍ID格式",
		})
		return
	}

	// 获取用户ID
	userIDValue, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未授权访问",
		})
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "无效的用户信息",
		})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的用户ID格式",
		})
		return
	}

	purchase, err := api.purchaseService.PurchaseBook(c.Request.Context(), userID, bookID)
	if err != nil {
		if strings.Contains(err.Error(), "already purchased") {
			c.JSON(http.StatusConflict, APIResponse{
				Code:    409,
				Message: "全书已购买",
			})
			return
		}
		if strings.Contains(err.Error(), "insufficient balance") {
			c.JSON(http.StatusForbidden, APIResponse{
				Code:    403,
				Message: "余额不足",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "购买全书失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "购买成功",
		Data:    purchase,
	})
}

// GetPurchases 获取购买记录
//
//	@Summary		获取购买记录
//	@Description	获取用户的章节购买记录列表
//	@Tags			章节购买
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"页码"	default(1)
//	@Param			size	query		int	false	"每页数量"	default(20)
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/reader/purchases [get]
func (api *ChapterCatalogAPI) GetPurchases(c *gin.Context) {
	// 获取用户ID
	userIDValue, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未授权访问",
		})
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "无效的用户信息",
		})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的用户ID格式",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	purchases, total, err := api.purchaseService.GetChapterPurchases(c.Request.Context(), userID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取购买记录失败",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "获取成功",
		Data:    purchases,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// GetBookPurchases 获取某本书的购买记录
//
//	@Summary		获取某本书的购买记录
//	@Description	获取用户对指定书籍的购买记录
//	@Tags			章节购买
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"书籍ID"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/reader/purchases/{bookId} [get]
func (api *ChapterCatalogAPI) GetBookPurchases(c *gin.Context) {
	bookIDStr := c.Param("id")
	if bookIDStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "书籍ID不能为空",
		})
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的书籍ID格式",
		})
		return
	}

	// 获取用户ID
	userIDValue, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未授权访问",
		})
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "无效的用户信息",
		})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的用户ID格式",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	purchases, total, err := api.purchaseService.GetBookPurchases(c.Request.Context(), userID, bookID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取购买记录失败",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "获取成功",
		Data:    purchases,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// CheckChapterAccess 检查章节访问权限
//
//	@Summary		检查章节访问权限
//	@Description	检查用户对指定章节的访问权限
//	@Tags			章节购买
//	@Accept			json
//	@Produce		json
//	@Param			chapterId	path		string	true	"章节ID"
//	@Success		200			{object}	APIResponse
//	@Failure		400			{object}	APIResponse
//	@Failure		404			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/bookstore/chapters/{chapterId}/access [get]
func (api *ChapterCatalogAPI) CheckChapterAccess(c *gin.Context) {
	chapterIdStr := c.Param("chapterId")
	if chapterIdStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "章节ID不能为空",
		})
		return
	}

	chapterID, err := primitive.ObjectIDFromHex(chapterIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的章节ID格式",
		})
		return
	}

	// 获取用户ID（可选）
	var userID primitive.ObjectID
	if userIDValue, exists := c.Get("userId"); exists {
		if uid, ok := userIDValue.(string); ok {
			userID, _ = primitive.ObjectIDFromHex(uid)
		}
	}

	accessInfo, err := api.purchaseService.CheckChapterAccess(c.Request.Context(), userID, chapterID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, APIResponse{
				Code:    404,
				Message: "章节不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "检查访问权限失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "检查成功",
		Data:    accessInfo,
	})
}
