package context

import (
	"net/http"

	"github.com/gin-gonic/gin"

	internalService "Qingyu_backend/service/internalapi"
)

// ContextAPI 上下文数据API处理器
// 供AI服务内部调用，获取项目上下文数据（角色、大纲、文档内容、角色关系等）
type ContextAPI struct {
	aggregator *internalService.ContextAggregator
}

// NewContextAPI 创建ContextAPI实例
func NewContextAPI(aggregator *internalService.ContextAggregator) *ContextAPI {
	return &ContextAPI{aggregator: aggregator}
}

// GetProjectContext 获取项目上下文汇总
// @Summary 获取项目上下文汇总
// @Description 获取项目基本信息、统计、设置等上下文数据
// @Tags Internal-Context
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} internalService.ProjectContext
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/internal/projects/{id}/context [get]
func (api *ContextAPI) GetProjectContext(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目ID不能为空"})
		return
	}

	result, err := api.aggregator.GetProjectContext(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCharacters 获取项目角色列表
// @Summary 获取项目角色列表
// @Description 获取指定项目下所有角色卡片数据
// @Tags Internal-Context
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} map[string][]internalService.CharacterInfo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/internal/projects/{id}/characters [get]
func (api *ContextAPI) GetCharacters(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目ID不能为空"})
		return
	}

	characters, err := api.aggregator.GetCharacters(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"characters": characters,
	})
}

// GetOutline 获取项目大纲树
// @Summary 获取项目大纲树
// @Description 获取指定项目的完整大纲树结构
// @Tags Internal-Context
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} map[string][]internalService.OutlineNodeInfo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/internal/projects/{id}/outline [get]
func (api *ContextAPI) GetOutline(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目ID不能为空"})
		return
	}

	nodes, err := api.aggregator.GetOutline(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"outline": nodes,
	})
}

// GetDocumentContent 获取文档内容
// @Summary 获取文档内容
// @Description 根据文档ID获取文档正文内容和字数统计
// @Tags Internal-Context
// @Produce json
// @Param id path string true "文档ID"
// @Success 200 {object} internalService.DocumentContentInfo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/internal/documents/{id}/content [get]
func (api *ContextAPI) GetDocumentContent(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文档ID不能为空"})
		return
	}

	content, err := api.aggregator.GetDocumentContent(c.Request.Context(), documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, content)
}

// GetCharacterRelations 获取角色关系列表
// @Summary 获取角色关系列表
// @Description 获取指定项目下所有角色关系数据
// @Tags Internal-Context
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} map[string][]internalService.RelationInfo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/internal/projects/{id}/relations [get]
func (api *ContextAPI) GetCharacterRelations(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目ID不能为空"})
		return
	}

	relations, err := api.aggregator.GetCharacterRelations(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"relations": relations,
	})
}
