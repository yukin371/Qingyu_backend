package bookstore

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/bookstore"
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
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	// 获取用户ID（从中间件设置的上下文中）
	var userID primitive.ObjectID
	if userIDValue, exists := c.Get("userId"); exists {
		if uid, ok := userIDValue.(string); ok {
			userID, _ = primitive.ObjectIDFromHex(uid)
		}
	}

	catalog, err := api.purchaseService.GetChapterCatalog(c.Request.Context(), userID.Hex(), bookID.Hex())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			shared.NotFound(c, "书籍不存在")
			return
		}
		shared.InternalError(c, "获取章节目录失败", err)
		return
	}

	shared.SuccessData(c, catalog)
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
		shared.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	chapterID, err := primitive.ObjectIDFromHex(chapterIdStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的章节ID格式")
		return
	}

	chapter, err := api.chapterService.GetChapterByID(c.Request.Context(), chapterID.Hex())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			shared.NotFound(c, "章节不存在")
			return
		}
		shared.InternalError(c, "获取章节信息失败", err)
		return
	}

	shared.SuccessData(c, chapter)
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
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	trialCount, _ := strconv.Atoi(c.DefaultQuery("count", "10"))
	if trialCount <= 0 {
		trialCount = 10
	}

	chapters, err := api.purchaseService.GetTrialChapters(c.Request.Context(), bookID.Hex(), trialCount)
	if err != nil {
		shared.InternalError(c, "获取试读章节失败", err)
		return
	}

	shared.Success(c, 200, "获取成功", map[string]interface{}{
		"book_id":  bookID.Hex(),
		"count":    len(chapters),
		"chapters": chapters,
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
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	chapters, err := api.purchaseService.GetVIPChapters(c.Request.Context(), bookID.Hex())
	if err != nil {
		shared.InternalError(c, "获取VIP章节失败", err)
		return
	}

	shared.Success(c, 200, "获取成功", map[string]interface{}{
		"book_id":  bookID.Hex(),
		"count":    len(chapters),
		"chapters": chapters,
	})
}

// GetChapterPrice 获取章节价格
//
//	@Summary		获取章节价格
//	@Description	获取指定章节的价格信息
//	@Tags			章节购买
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"章节ID"
//	@Success		200			{object}	APIResponse
//	@Failure		400			{object}	APIResponse
//	@Failure		404			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/bookstore/chapters/{id}/price [get]
func (api *ChapterCatalogAPI) GetChapterPrice(c *gin.Context) {
	chapterIdStr := c.Param("chapterId")
	if chapterIdStr == "" {
		shared.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	chapterID, err := primitive.ObjectIDFromHex(chapterIdStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的章节ID格式")
		return
	}

	price, err := api.purchaseService.GetChapterPrice(c.Request.Context(), chapterID.Hex())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			shared.NotFound(c, "章节不存在")
			return
		}
		shared.InternalError(c, "获取章节价格失败", err)
		return
	}

	shared.Success(c, 200, "获取成功", map[string]interface{}{
		"chapter_id": chapterID.Hex(),
		"price":      price,
	})
}

// PurchaseChapter 购买章节
//
//	@Summary		购买章节
//	@Description	购买单个付费章节
//	@Tags			章节购买
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"章节ID"
//	@Success		200			{object}	APIResponse
//	@Failure		400			{object}	APIResponse
//	@Failure		403			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/reader/chapters/{id}/purchase [post]
func (api *ChapterCatalogAPI) PurchaseChapter(c *gin.Context) {
	chapterIdStr := c.Param("chapterId")
	if chapterIdStr == "" {
		shared.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	chapterID, err := primitive.ObjectIDFromHex(chapterIdStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的章节ID格式")
		return
	}

	// 获取用户ID
	userIDValue, exists := c.Get("userId")
	if !exists {
		shared.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		shared.Unauthorized(c, "无效的用户信息")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的用户ID格式")
		return
	}

	purchase, err := api.purchaseService.PurchaseChapter(c.Request.Context(), userID.Hex(), chapterID.Hex())
	if err != nil {
		if strings.Contains(err.Error(), "already purchased") {
			shared.Error(c, http.StatusConflict, "章节已购买", "")
			return
		}
		if strings.Contains(err.Error(), "insufficient balance") {
			shared.Forbidden(c, "余额不足")
			return
		}
		if strings.Contains(err.Error(), "free chapter") {
			shared.BadRequest(c, "参数错误", "免费章节无需购买")
			return
		}
		shared.InternalError(c, "购买章节失败", err)
		return
	}

	shared.Success(c, 200, "购买成功", purchase)
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
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	// 获取用户ID
	userIDValue, exists := c.Get("userId")
	if !exists {
		shared.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		shared.Unauthorized(c, "无效的用户信息")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的用户ID格式")
		return
	}

	purchase, err := api.purchaseService.PurchaseBook(c.Request.Context(), userID.Hex(), bookID.Hex())
	if err != nil {
		if strings.Contains(err.Error(), "already purchased") {
			shared.Error(c, http.StatusConflict, "全书已购买", "")
			return
		}
		if strings.Contains(err.Error(), "insufficient balance") {
			shared.Forbidden(c, "余额不足")
			return
		}
		shared.InternalError(c, "购买全书失败", err)
		return
	}

	shared.Success(c, 200, "购买成功", purchase)
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
		shared.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		shared.Unauthorized(c, "无效的用户信息")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的用户ID格式")
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

	purchases, total, err := api.purchaseService.GetChapterPurchases(c.Request.Context(), userID.Hex(), page, size)
	if err != nil {
		shared.InternalError(c, "获取购买记录失败", err)
		return
	}

	shared.Paginated(c, purchases, total, page, size, "获取成功")
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
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	// 获取用户ID
	userIDValue, exists := c.Get("userId")
	if !exists {
		shared.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		shared.Unauthorized(c, "无效的用户信息")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的用户ID格式")
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

	purchases, total, err := api.purchaseService.GetBookPurchases(c.Request.Context(), userID.Hex(), bookID.Hex(), page, size)
	if err != nil {
		shared.InternalError(c, "获取购买记录失败", err)
		return
	}

	shared.Paginated(c, purchases, total, page, size, "获取成功")
}

// CheckChapterAccess 检查章节访问权限
//
//	@Summary		检查章节访问权限
//	@Description	检查用户对指定章节的访问权限
//	@Tags			章节购买
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"章节ID"
//	@Success		200			{object}	APIResponse
//	@Failure		400			{object}	APIResponse
//	@Failure		404			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/bookstore/chapters/{id}/access [get]
func (api *ChapterCatalogAPI) CheckChapterAccess(c *gin.Context) {
	chapterIdStr := c.Param("chapterId")
	if chapterIdStr == "" {
		shared.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	chapterID, err := primitive.ObjectIDFromHex(chapterIdStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的章节ID格式")
		return
	}

	// 获取用户ID（可选）
	var userID primitive.ObjectID
	if userIDValue, exists := c.Get("userId"); exists {
		if uid, ok := userIDValue.(string); ok {
			userID, _ = primitive.ObjectIDFromHex(uid)
		}
	}

	accessInfo, err := api.purchaseService.CheckChapterAccess(c.Request.Context(), userID.Hex(), chapterID.Hex())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			shared.NotFound(c, "章节不存在")
			return
		}
		shared.InternalError(c, "检查访问权限失败", err)
		return
	}

	shared.SuccessData(c, accessInfo)
}
