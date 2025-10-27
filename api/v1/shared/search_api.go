package shared

import (
	"Qingyu_backend/service/shared/search"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SearchAPI 搜索API处理器
type SearchAPI struct {
	searchService search.SearchService
}

// NewSearchAPI 创建搜索API实例
func NewSearchAPI(searchService search.SearchService) *SearchAPI {
	return &SearchAPI{
		searchService: searchService,
	}
}

// ============ 搜索API ============

// SearchBooks 搜索书籍
//
//	@Summary		搜索书籍
//	@Description	根据关键词搜索书籍，支持标题、作者、描述、标签等字段
//	@Tags			搜索
//	@Accept			json
//	@Produce		json
//	@Param			q			query		string	true	"搜索关键词"
//	@Param			category	query		string	false	"分类过滤"
//	@Param			author		query		string	false	"作者过滤"
//	@Param			sort_by		query		string	false	"排序方式: relevance(相关度), time(时间), popularity(热度)"	default(relevance)
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Success		200			{object}	APIResponse{data=search.SearchResult}
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/search/books [get]
func (api *SearchAPI) SearchBooks(c *gin.Context) {
	// 1. 获取搜索参数
	keyword := c.Query("q")
	if keyword == "" {
		BadRequest(c, "参数错误", "搜索关键词不能为空")
		return
	}

	category := c.Query("category")
	author := c.Query("author")
	sortBy := c.DefaultQuery("sort_by", "relevance")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 2. 构建搜索请求
	req := &search.SearchRequest{
		Keyword:  keyword,
		Category: category,
		Author:   author,
		SortBy:   sortBy,
		Page:     page,
		PageSize: pageSize,
	}

	// 3. 执行搜索
	result, err := api.searchService.SearchBooks(c.Request.Context(), req)
	if err != nil {
		Error(c, 500, "搜索失败", err.Error())
		return
	}

	// 4. 返回结果
	Success(c, 200, "搜索成功", result)
}

// SearchDocuments 搜索文档
//
//	@Summary		搜索文档
//	@Description	根据关键词搜索用户项目中的文档
//	@Tags			搜索
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			q			query		string	true	"搜索关键词"
//	@Param			project_id	query		string	false	"项目ID过滤"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Success		200			{object}	APIResponse{data=search.SearchResult}
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/search/documents [get]
func (api *SearchAPI) SearchDocuments(c *gin.Context) {
	// 1. 验证用户登录
	_, exists := c.Get("userId")
	if !exists {
		Unauthorized(c, "未授权")
		return
	}

	// 2. 获取搜索参数
	keyword := c.Query("q")
	if keyword == "" {
		BadRequest(c, "参数错误", "搜索关键词不能为空")
		return
	}

	projectID := c.Query("project_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 3. 构建搜索请求
	req := &search.SearchRequest{
		Keyword:   keyword,
		ProjectID: projectID,
		Page:      page,
		PageSize:  pageSize,
	}

	// 4. 执行搜索
	result, err := api.searchService.SearchDocuments(c.Request.Context(), req)
	if err != nil {
		Error(c, 500, "搜索失败", err.Error())
		return
	}

	// 5. 返回结果
	Success(c, 200, "搜索成功", result)
}

// GetSearchSuggestions 获取搜索建议
//
//	@Summary		获取搜索建议
//	@Description	根据关键词前缀获取搜索建议（自动补全）
//	@Tags			搜索
//	@Accept			json
//	@Produce		json
//	@Param			q		query		string	true	"搜索关键词前缀"
//	@Param			limit	query		int		false	"返回数量"	default(10)
//	@Success		200		{object}	APIResponse{data=[]string}
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/v1/search/suggest [get]
func (api *SearchAPI) GetSearchSuggestions(c *gin.Context) {
	// 1. 获取参数
	keyword := c.Query("q")
	if keyword == "" {
		BadRequest(c, "参数错误", "搜索关键词不能为空")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 2. 获取建议
	suggestions, err := api.searchService.GetSuggestions(c.Request.Context(), keyword, limit)
	if err != nil {
		Error(c, 500, "获取搜索建议失败", err.Error())
		return
	}

	// 3. 返回结果
	Success(c, 200, "获取成功", suggestions)
}

// TODO(Phase3): 搜索历史API
// GetSearchHistory 获取用户搜索历史
// func (api *SearchAPI) GetSearchHistory(c *gin.Context) {
//     userID, _ := c.Get("userId")
//     history, err := api.searchService.GetSearchHistory(c.Request.Context(), userID.(string), 20)
//     if err != nil {
//         Error(c, 500, "获取搜索历史失败", err.Error())
//         return
//     }
//     Success(c, 200, "获取成功", history)
// }

// TODO(Phase3): 热门搜索API
// GetHotSearches 获取热门搜索词
// func (api *SearchAPI) GetHotSearches(c *gin.Context) {
//     hotSearches, err := api.searchService.GetHotSearches(c.Request.Context(), 10)
//     if err != nil {
//         Error(c, 500, "获取热门搜索失败", err.Error())
//         return
//     }
//     Success(c, 200, "获取成功", hotSearches)
// }
