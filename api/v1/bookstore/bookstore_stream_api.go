package bookstore

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Qingyu_backend/models/bookstore"
	streamService "Qingyu_backend/service/bookstore"
	"Qingyu_backend/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StreamSearchBooks 流式搜索书籍
//
//	@Summary		流式搜索书籍
//	@Description	使用游标机制流式返回搜索结果，NDJSON格式
//	@Tags			书籍
//	@Accept			json
//	@Produce		application/x-ndjson
//	@Param			keyword		query		string		false	"搜索关键词"
//	@Param			cursor		query		string		false	"游标"
//	@Param			limit		query		int			false	"每批数量"	default(20)
//	@Param			categoryId	query		string		false	"分类ID"
//	@Param			author		query		string		false	"作者"
//	@Param			status		query		string		false	"状态"
//	@Param			tags		query		[]string	false	"标签"
//	@Param			sortBy		query		string		false	"排序字段"	default(created_at)
//	@Param			sortOrder	query		string		false	"排序方向"	default(desc)
//	@Success		200			{object}	object	"流式响应"
//	@Failure		400			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/bookstore/books/stream [get]
func (api *BookstoreAPI) StreamSearchBooks(c *gin.Context) {
	// 解析参数
	keyword := c.Query("keyword")
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	categoryID := c.Query("categoryId")
	author := c.Query("author")
	status := c.Query("status")
	tags := c.QueryArray("tags")
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	// 验证limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// 构建过滤器
	filter := &bookstore.BookFilter{
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Limit:     limit,
	}

	if keyword != "" {
		filter.Keyword = &keyword
	}
	if cursor != "" {
		filter.Cursor = &cursor
	}
	if categoryID != "" {
		filter.CategoryID = &categoryID
	}
	if author != "" {
		filter.Author = &author
	}
	if status != "" {
		bookStatus := bookstore.BookStatus(status)
		filter.Status = &bookStatus
	}
	if len(tags) > 0 {
		filter.Tags = tags
	}

	// 设置流式响应头
	c.Writer.Header().Set("Content-Type", "application/x-ndjson")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// 创建流式写入器
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		response.InternalError(c, nil)
		return
	}

	// 创建流式服务
	svc := streamService.NewBookstoreStreamService(nil) // TODO: 注入真实的repository

	// 发送元数据
	api.sendStreamMeta(c.Writer, cursor, nil, true)
	flusher.Flush()

	// 执行搜索
	result, err := svc.StreamSearch(c.Request.Context(), filter)
	if err != nil {
		api.sendStreamError(c.Writer, err)
		flusher.Flush()
		return
	}

	// 发送数据
	for _, book := range result.Books {
		api.sendStreamData(c.Writer, []*bookstore.Book{book})
		flusher.Flush()
	}

	// 发送完成信号
	api.sendStreamDone(c.Writer, result.NextCursor, len(result.Books))
	flusher.Flush()
}

// sendStreamMeta 发送流式元数据
func (api *BookstoreAPI) sendStreamMeta(w http.ResponseWriter, cursor string, total *int64, hasMore bool) {
	msg := map[string]interface{}{
		"type":    "meta",
		"cursor":  cursor,
		"total":   total,
		"hasMore": hasMore,
	}
	json.NewEncoder(w).Encode(msg)
	w.Write([]byte("\n"))
}

// sendStreamData 发送流式数据
func (api *BookstoreAPI) sendStreamData(w http.ResponseWriter, books []*bookstore.Book) {
	msg := map[string]interface{}{
		"type":  "data",
		"books": books,
	}
	json.NewEncoder(w).Encode(msg)
	w.Write([]byte("\n"))
}

// sendStreamProgress 发送流式进度
func (api *BookstoreAPI) sendStreamProgress(w http.ResponseWriter, loaded int, total *int64) {
	msg := map[string]interface{}{
		"type":   "progress",
		"loaded": loaded,
		"total":  total,
	}
	json.NewEncoder(w).Encode(msg)
	w.Write([]byte("\n"))
}

// sendStreamDone 发送流式完成信号
func (api *BookstoreAPI) sendStreamDone(w http.ResponseWriter, cursor string, total int) {
	msg := map[string]interface{}{
		"type":   "done",
		"cursor": cursor,
		"total":  total,
	}
	json.NewEncoder(w).Encode(msg)
	w.Write([]byte("\n"))
}

// sendStreamError 发送流式错误
func (api *BookstoreAPI) sendStreamError(w http.ResponseWriter, err error) {
	msg := map[string]interface{}{
		"type":  "error",
		"error": err.Error(),
	}
	json.NewEncoder(w).Encode(msg)
	w.Write([]byte("\n"))
}

// StreamSearchBooksBatch 批量流式搜索书籍（用于客户端不支持流式的情况）
//
//	@Summary		批量流式搜索书籍
//	@Description	使用游标机制返回搜索结果，非流式响应
//	@Tags			书籍
//	@Accept			json
//	@Produce		json
//	@Param			keyword		query		string		false	"搜索关键词"
//	@Param			cursor		query		string		false	"游标"
//	@Param			limit		query		int			false	"每批数量"	default(20)
//	@Success		200			{object}	StreamSearchResult
//	@Failure		400			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/bookstore/books/stream-batch [get]
func (api *BookstoreAPI) StreamSearchBooksBatch(c *gin.Context) {
	// 解析参数
	keyword := c.Query("keyword")
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	categoryID := c.Query("categoryId")
	author := c.Query("author")
	status := c.Query("status")
	tags := c.QueryArray("tags")
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	// 验证limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// 构建过滤器
	filter := &bookstore.BookFilter{
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Limit:     limit,
	}

	if keyword != "" {
		filter.Keyword = &keyword
	}
	if cursor != "" {
		filter.Cursor = &cursor
	}
	if categoryID != "" {
		filter.CategoryID = &categoryID
	}
	if author != "" {
		filter.Author = &author
	}
	if status != "" {
		bookStatus := bookstore.BookStatus(status)
		filter.Status = &bookStatus
	}
	if len(tags) > 0 {
		filter.Tags = tags
	}

	// 创建流式服务
	svc := streamService.NewBookstoreStreamService(nil) // TODO: 注入真实的repository

	// 执行搜索
	result, err := svc.StreamSearchWithCursor(c.Request.Context(), filter)
	if err != nil {
		api.logger.Error("Stream search failed",
			zap.Error(err),
			zap.String("keyword", keyword),
			zap.String("cursor", cursor))
		response.InternalError(c, err)
		return
	}

	// 转换为DTO
	bookDTOs := ToBookDTOsFromPtrSlice(result.Books)

	responseData := map[string]interface{}{
		"books":      bookDTOs,
		"nextCursor": result.NextCursor,
		"hasMore":    result.HasMore,
		"total":      result.Total,
	}

	response.SuccessWithMessage(c, "流式搜索成功", responseData)
}
