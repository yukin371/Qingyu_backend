package writer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/search"
	searchservice "Qingyu_backend/service/search"
	"Qingyu_backend/pkg/response"
	"errors"
)

// SearchAPI 搜索API处理器（写作端）
type SearchAPI struct {
	searchService *searchservice.SearchService
}

// NewSearchAPI 创建搜索API实例
func NewSearchAPI(searchService *searchservice.SearchService) *SearchAPI {
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
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		401			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/writer/search/documents [get]
func (api *SearchAPI) SearchDocuments(c *gin.Context) {
	// 1. 验证用户登录
	userID, exists := c.Get("userId")
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
	req := &search.SearchRequest{
		Type:     search.SearchTypeDocuments,
		Query:    keyword,
		Page:     page,
		PageSize: pageSize,
		Filter:   make(map[string]interface{}),
	}

	// 添加用户ID过滤（必需）
	req.Filter["user_id"] = userID

	// 如果指定了项目ID，添加项目过滤
	if projectID != "" {
		req.Filter["project_id"] = projectID
	}

	// 4. 执行搜索
	resp, err := api.searchService.Search(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 5. 检查响应
	if !resp.Success {
		response.InternalError(c, errors.New(resp.Error.Message))
		return
	}

	// 6. 返回结果
	shared.Success(c, http.StatusOK, "搜索成功", resp.Data)
}
