package admin

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/pkg/response"
	adminRepo "Qingyu_backend/repository/interfaces/admin"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	adminService "Qingyu_backend/service/admin"
)

// ContentExportAPI 内容导出API处理器
type ContentExportAPI struct {
	bookRepo          bookstoreRepo.BookRepository
	chapterRepo       bookstoreRepo.ChapterRepository
	exportHistoryRepo adminRepo.ExportHistoryRepository
	exportService     adminService.ExportService
}

// NewContentExportAPI 创建内容导出API实例
func NewContentExportAPI(
	bookRepo bookstoreRepo.BookRepository,
	chapterRepo bookstoreRepo.ChapterRepository,
	exportHistoryRepo adminRepo.ExportHistoryRepository,
) *ContentExportAPI {
	return &ContentExportAPI{
		bookRepo:          bookRepo,
		chapterRepo:       chapterRepo,
		exportHistoryRepo: exportHistoryRepo,
		exportService:     adminService.NewExportService(),
	}
}

// ==================== 书籍导出 ====================

// ExportBooks 导出书籍数据
//
//	@Summary		导出书籍数据
//	@Description	管理员导出书籍数据，支持CSV和Excel格式
//	@Tags			Admin-Content
//	@Accept			json
//	@Produce		json
//	@Param			format		query		string	false	"导出格式"	Enums(csv,excel)
//	@Param			status		query		string	false	"状态筛选"
//	@Param			author		query		string	false	"作者筛选"
//	@Param			start_date	query		string	false	"开始日期"	format(date)
//	@Param			end_date		query		string	false	"结束日期"	format(date)
//	@Success		200			{file}		file
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/admin/content/books/export [get]
func (api *ContentExportAPI) ExportBooks(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")

	// 验证格式
	if format != "csv" && format != "excel" {
		response.BadRequest(c, "参数错误", "不支持的导出格式，仅支持csv和excel")
		return
	}

	// 构建过滤器
	filter := &base.BaseFilter{
		Conditions: make(map[string]interface{}),
	}

	// 添加状态筛选
	if status := c.Query("status"); status != "" {
		filter.Conditions["status"] = status
	}

	// 添加作者筛选
	if author := c.Query("author"); author != "" {
		filter.Conditions["author"] = author
	}

	// 添加日期范围筛选
	if startDate := c.Query("start_date"); startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			if filter.CreatedAt == nil {
				filter.CreatedAt = &base.TimeRange{}
			}
			filter.CreatedAt.Start = &parsedDate
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			// 设置为当天的23:59:59
			endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 0, parsedDate.Location())
			if filter.CreatedAt == nil {
				filter.CreatedAt = &base.TimeRange{}
			}
			filter.CreatedAt.End = &endOfDay
		}
	}

	// 获取书籍列表
	books, err := api.bookRepo.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 转换为可导出格式
	exportableBooks := make([]adminService.Exportable, len(books))
	for i, book := range books {
		exportableBooks[i] = &BookExportAdapter{book: book}
	}

	// 定义导出列
	columns := []adminService.ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "Title", Title: "书名"},
		{Key: "Author", Title: "作者"},
		{Key: "Introduction", Title: "简介"},
		{Key: "Status", Title: "状态"},
		{Key: "WordCount", Title: "字数"},
		{Key: "ChapterCount", Title: "章节数"},
		{Key: "Price", Title: "价格(分)"},
		{Key: "IsFree", Title: "是否免费"},
		{Key: "CreatedAt", Title: "创建时间"},
	}

	// 根据格式导出
	switch format {
	case "csv":
		api.exportBooksToCSV(c, exportableBooks, columns)
	case "excel":
		api.exportBooksToExcel(c, exportableBooks, columns)
	}
}

// ==================== 章节导出 ====================

// ExportChapters 导出章节数据
//
//	@Summary		导出章节数据
//	@Description	管理员导出章节数据，支持CSV和Excel格式
//	@Tags			Admin-Content
//	@Accept			json
//	@Produce		json
//	@Param			format		query		string	false	"导出格式"	Enums(csv,excel)
//	@Param			book_id		query		string	false	"书籍ID"
//	@Param			start_date	query		string	false	"开始日期"	format(date)
//	@Param			end_date		query		string	false	"结束日期"	format(date)
//	@Success		200			{file}		file
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/admin/content/chapters/export [get]
func (api *ContentExportAPI) ExportChapters(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")

	// 验证格式
	if format != "csv" && format != "excel" {
		response.BadRequest(c, "参数错误", "不支持的导出格式，仅支持csv和excel")
		return
	}

	// 构建过滤器
	filter := &base.BaseFilter{
		Conditions: make(map[string]interface{}),
	}

	// 添加书籍ID筛选
	if bookID := c.Query("book_id"); bookID != "" {
		filter.Conditions["book_id"] = bookID
	}

	// 添加日期范围筛选（按发布时间）
	if startDate := c.Query("start_date"); startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.Conditions["publishTime.$gte"] = parsedDate
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 0, parsedDate.Location())
			filter.Conditions["publishTime.$lte"] = endOfDay
		}
	}

	// 获取章节列表
	chapters, err := api.chapterRepo.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 转换为可导出格式
	exportableChapters := make([]adminService.Exportable, len(chapters))
	for i, chapter := range chapters {
		exportableChapters[i] = &ChapterExportAdapter{chapter: chapter}
	}

	// 定义导出列
	columns := []adminService.ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "BookID", Title: "书籍ID"},
		{Key: "Title", Title: "章节标题"},
		{Key: "ChapterNum", Title: "章节号"},
		{Key: "WordCount", Title: "字数"},
		{Key: "Price", Title: "价格(分)"},
		{Key: "IsFree", Title: "是否免费"},
		{Key: "PublishTime", Title: "发布时间"},
		{Key: "CreatedAt", Title: "创建时间"},
	}

	// 根据格式导出
	switch format {
	case "csv":
		api.exportChaptersToCSV(c, exportableChapters, columns)
	case "excel":
		api.exportChaptersToExcel(c, exportableChapters, columns)
	}
}

// ==================== 导出模板 ====================

// GetBookExportTemplate 获取书籍导出模板
//
//	@Summary		获取书籍导出模板
//	@Description	获取书籍导出的字段模板
//	@Tags			Admin-Content
//	@Accept			json
//	@Produce		json
//	@Success		200	{file}		file
//	@Router			/api/v1/admin/content/books/export/template [get]
func (api *ContentExportAPI) GetBookExportTemplate(c *gin.Context) {
	columns := []adminService.ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "Title", Title: "书名"},
		{Key: "Author", Title: "作者"},
		{Key: "Introduction", Title: "简介"},
		{Key: "Status", Title: "状态"},
		{Key: "WordCount", Title: "字数"},
		{Key: "ChapterCount", Title: "章节数"},
		{Key: "Price", Title: "价格(分)"},
		{Key: "IsFree", Title: "是否免费"},
		{Key: "CreatedAt", Title: "创建时间"},
	}

	template, err := api.exportService.GetExportTemplate(columns)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=books_template.csv")
	c.Data(200, "text/csv", template)
}

// GetChapterExportTemplate 获取章节导出模板
//
//	@Summary		获取章节导出模板
//	@Description	获取章节导出的字段模板
//	@Tags			Admin-Content
//	@Accept			json
//	@Produce		json
//	@Success		200	{file}		file
//	@Router			/api/v1/admin/content/chapters/export/template [get]
func (api *ContentExportAPI) GetChapterExportTemplate(c *gin.Context) {
	columns := []adminService.ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "BookID", Title: "书籍ID"},
		{Key: "Title", Title: "章节标题"},
		{Key: "ChapterNum", Title: "章节号"},
		{Key: "WordCount", Title: "字数"},
		{Key: "Price", Title: "价格(分)"},
		{Key: "IsFree", Title: "是否免费"},
		{Key: "PublishTime", Title: "发布时间"},
		{Key: "CreatedAt", Title: "创建时间"},
	}

	template, err := api.exportService.GetExportTemplate(columns)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=chapters_template.csv")
	c.Data(200, "text/csv", template)
}

// ==================== 私有导出方法 ====================

// exportBooksToCSV 导出书籍为CSV
func (api *ContentExportAPI) exportBooksToCSV(c *gin.Context, data []adminService.Exportable, columns []adminService.ExportColumn) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=books.csv")

	result, err := api.exportService.ExportToCSV(c.Request.Context(), data, columns)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	c.Data(200, "text/csv", result)
}

// exportBooksToExcel 导出书籍为Excel
func (api *ContentExportAPI) exportBooksToExcel(c *gin.Context, data []adminService.Exportable, columns []adminService.ExportColumn) {
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=books.xlsx")

	result, err := api.exportService.ExportToExcel(c.Request.Context(), data, columns)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", result)
}

// exportChaptersToCSV 导出章节为CSV
func (api *ContentExportAPI) exportChaptersToCSV(c *gin.Context, data []adminService.Exportable, columns []adminService.ExportColumn) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=chapters.csv")

	result, err := api.exportService.ExportToCSV(c.Request.Context(), data, columns)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	c.Data(200, "text/csv", result)
}

// exportChaptersToExcel 导出章节为Excel
func (api *ContentExportAPI) exportChaptersToExcel(c *gin.Context, data []adminService.Exportable, columns []adminService.ExportColumn) {
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=chapters.xlsx")

	result, err := api.exportService.ExportToExcel(c.Request.Context(), data, columns)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", result)
}

// ==================== 导出数据适配器 ====================

// BookExportAdapter 书籍导出适配器
type BookExportAdapter struct {
	book *bookstore.Book
}

func (a *BookExportAdapter) ToExportRow() []string {
	return []string{
		a.book.ID.Hex(),
		a.book.Title,
		a.book.Author,
		a.book.Introduction,
		string(a.book.Status),
		strconv.FormatInt(a.book.WordCount, 10),
		strconv.Itoa(a.book.ChapterCount),
		strconv.FormatFloat(a.book.Price, 'f', 2, 64),
		strconv.FormatBool(a.book.IsFree),
		a.book.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ChapterExportAdapter 章节导出适配器
type ChapterExportAdapter struct {
	chapter *bookstore.Chapter
}

func (a *ChapterExportAdapter) ToExportRow() []string {
	return []string{
		a.chapter.ID,
		a.chapter.BookID,
		a.chapter.Title,
		strconv.Itoa(a.chapter.ChapterNum),
		strconv.Itoa(a.chapter.WordCount),
		strconv.FormatFloat(a.chapter.Price, 'f', 2, 64),
		strconv.FormatBool(a.chapter.IsFree),
		a.chapter.PublishTime.Format("2006-01-02 15:04:05"),
		a.chapter.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
