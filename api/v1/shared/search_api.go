package shared

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/service/shared/search"
)

// SearchAPI 搜索API处理器
// 注意：此模块只提供通用的搜索建议服务
// - 书籍搜索已迁移到 bookstore 模块
// - 文档搜索已迁移到 writer 模块
type SearchAPI struct {
	searchService search.SearchService
}

// NewSearchAPI 创建搜索API实例
func NewSearchAPI(searchService search.SearchService) *SearchAPI {
	return &SearchAPI{
		searchService: searchService,
	}
}

// GetSearchSuggestions 获取搜索建议
//
//	@Summary		获取搜索建议
//	@Description	根据关键词前缀获取搜索建议（自动补全），支持书籍和文档的搜索建议
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
		Error(c, http.StatusInternalServerError, "获取搜索建议失败", err.Error())
		return
	}

	// 3. 返回结果
	Success(c, http.StatusOK, "获取成功", suggestions)
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
