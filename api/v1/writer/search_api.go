package writer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/shared/search_legacy"
)

// SearchAPI 搜索API处理器（写作端）
type SearchAPI struct {
	searchService search_legacy.SearchService
}

// NewSearchAPI 创建搜索API实例
func NewSearchAPI(searchService search_legacy.SearchService) *SearchAPI {
	return &SearchAPI{
		searchService: searchService,
	}
}

// SearchDocuments 搜索文档
//
//	@Summary		搜索文档
//	@Description	根据关键词搜索用户项目中的文档
//	@Tags			写作端-搜索
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			q			query		string	true	"搜索关键词"
//	@Param			project_id	query		string	false	"项目ID过滤"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Success		200			{object}	shared.APIResponse{data=search.SearchResult}
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		401			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/writer/search/documents [get]
func (api *SearchAPI) SearchDocuments(c *gin.Context) {
	// 1. 验证用户登录
	_, exists := c.Get("userId")
	if !exists {
		shared.Unauthorized(c, "未授权")
		return
	}

	// 2. 获取搜索参数
	keyword := c.Query("q")
	if keyword == "" {
		shared.BadRequest(c, "参数错误", "搜索关键词不能为空")
		return
	}

	projectID := c.Query("project_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 3. 构建搜索请求
	req := &search_legacy.SearchRequest{
		Keyword:   keyword,
		ProjectID: projectID,
		Page:      page,
		PageSize:  pageSize,
	}

	// 4. 执行搜索
	result, err := api.searchService.SearchDocuments(c.Request.Context(), req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "搜索失败", err.Error())
		return
	}

	// 5. 返回结果
	shared.Success(c, http.StatusOK, "搜索成功", result)
}
