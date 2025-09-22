package ai

import (
	"net/http"
	"strconv"

	"Qingyu_backend/models/ai"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
)

// NovelContextApi 小说上下文API控制器
type NovelContextApi struct {
	contextService *aiService.NovelContextService
}

// NewNovelContextApi 创建小说上下文API控制器
func NewNovelContextApi(contextService *aiService.NovelContextService) *NovelContextApi {
	return &NovelContextApi{
		contextService: contextService,
	}
}

// GetNovelContext 获取小说写作上下文
// GET /api/v1/novel/{projectId}/context?position={position}&maxTokens={maxTokens}
func (a *NovelContextApi) GetNovelContext(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "项目ID不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	position := c.Query("position")
	maxTokensStr := c.DefaultQuery("maxTokens", "4000")
	focusType := c.Query("focusType")
	focusID := c.Query("focusId")

	maxTokens, err := strconv.Atoi(maxTokensStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "maxTokens参数格式错误",
			"timestamp": getTimestamp(),
		})
		return
	}

	req := &ai.ContextBuildRequest{
		ProjectID:       projectID,
		CurrentPosition: position,
		MaxTokens:       maxTokens,
		FocusType:       focusType,
		FocusID:         focusID,
	}

	response, err := a.contextService.BuildContext(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "构建上下文失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      response,
		"timestamp": getTimestamp(),
	})
}

// CreateNovelMemory 创建小说记忆
// POST /api/v1/novel/{projectId}/memory
func (a *NovelContextApi) CreateNovelMemory(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "项目ID不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	var req struct {
		Type       string                 `json:"type" binding:"required"` // character, plot, setting, chapter
		Title      string                 `json:"title" binding:"required"`
		Content    string                 `json:"content" binding:"required"`
		Summary    string                 `json:"summary"`
		Metadata   map[string]interface{} `json:"metadata"`
		Importance int                    `json:"importance"` // 1-10
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "请求参数错误: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 验证类型
	validTypes := map[string]bool{
		"character": true,
		"plot":      true,
		"setting":   true,
		"chapter":   true,
	}
	if !validTypes[req.Type] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "无效的记忆类型",
			"timestamp": getTimestamp(),
		})
		return
	}

	// 设置默认重要性
	if req.Importance == 0 {
		req.Importance = 5
	}
	if req.Importance < 1 || req.Importance > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "重要性必须在1-10之间",
			"timestamp": getTimestamp(),
		})
		return
	}

	novelContext := &ai.NovelContext{
		ProjectID:  projectID,
		Type:       req.Type,
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Metadata:   req.Metadata,
		Importance: req.Importance,
	}

	err := a.contextService.StoreContext(c.Request.Context(), novelContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "创建记忆失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      novelContext,
		"timestamp": getTimestamp(),
	})
}

// SearchNovelContext 搜索小说相关内容
// GET /api/v1/novel/{projectId}/search?q={query}&type={type}&limit={limit}
func (a *NovelContextApi) SearchNovelContext(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "项目ID不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "搜索查询不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	contextType := c.Query("type")
	limitStr := c.DefaultQuery("limit", "20")
	vectorWeightStr := c.DefaultQuery("vectorWeight", "0.6")
	keywordWeightStr := c.DefaultQuery("keywordWeight", "0.3")
	metadataWeightStr := c.DefaultQuery("metadataWeight", "0.1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "limit参数格式错误",
			"timestamp": getTimestamp(),
		})
		return
	}

	vectorWeight, _ := strconv.ParseFloat(vectorWeightStr, 32)
	keywordWeight, _ := strconv.ParseFloat(keywordWeightStr, 32)
	metadataWeight, _ := strconv.ParseFloat(metadataWeightStr, 32)

	options := &aiService.SearchOptions{
		VectorWeight:   float32(vectorWeight),
		KeywordWeight:  float32(keywordWeight),
		MetadataWeight: float32(metadataWeight),
		MaxResults:     limit,
		MinScore:       0.1,
	}

	results, err := a.contextService.SearchContext(c.Request.Context(), projectID, query, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "搜索失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 如果指定了类型，进行过滤
	if contextType != "" {
		filteredResults := make([]*ai.RetrievalResult, 0)
		for _, result := range results {
			if result.Context.Type == contextType {
				filteredResults = append(filteredResults, result)
			}
		}
		results = filteredResults
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"results": results,
			"total":   len(results),
			"query":   query,
		},
		"timestamp": getTimestamp(),
	})
}

// UpdateNovelMemory 更新小说记忆
// PUT /api/v1/novel/{projectId}/memory/{memoryId}
func (a *NovelContextApi) UpdateNovelMemory(c *gin.Context) {
	projectID := c.Param("projectId")
	memoryID := c.Param("memoryId")

	if projectID == "" || memoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "项目ID和记忆ID不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	var req struct {
		Title      string                 `json:"title"`
		Content    string                 `json:"content"`
		Summary    string                 `json:"summary"`
		Metadata   map[string]interface{} `json:"metadata"`
		Importance int                    `json:"importance"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "请求参数错误: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 这里需要先从数据库获取现有记忆，然后更新
	// 简化实现，实际需要添加数据库查询逻辑
	novelContext := &ai.NovelContext{
		ID:         memoryID,
		ProjectID:  projectID,
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Metadata:   req.Metadata,
		Importance: req.Importance,
	}

	err := a.contextService.UpdateContext(c.Request.Context(), novelContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "更新记忆失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      novelContext,
		"timestamp": getTimestamp(),
	})
}

// DeleteNovelMemory 删除小说记忆
// DELETE /api/v1/novel/{projectId}/memory/{memoryId}
func (a *NovelContextApi) DeleteNovelMemory(c *gin.Context) {
	projectID := c.Param("projectId")
	memoryID := c.Param("memoryId")

	if projectID == "" || memoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "项目ID和记忆ID不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	err := a.contextService.DeleteContext(c.Request.Context(), memoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "删除记忆失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"timestamp": getTimestamp(),
	})
}
