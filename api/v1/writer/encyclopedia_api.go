package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	writerModels "Qingyu_backend/models/writer"
	writerInterface "Qingyu_backend/repository/interfaces/writer"
)

// EncyclopediaApi 设定百科API
type EncyclopediaApi struct {
	conceptRepo writerInterface.ConceptRepository
}

// NewEncyclopediaApi 创建设定百科API
func NewEncyclopediaApi(conceptRepo writerInterface.ConceptRepository) *EncyclopediaApi {
	return &EncyclopediaApi{conceptRepo: conceptRepo}
}

// ListConcepts 获取项目下的概念列表
// @Summary 获取概念列表
// @Description 获取项目下的所有概念
// @Tags 设定百科
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} response.APIResponse{data=[]*writerModels.Concept}
// @Router /api/v1/writer/projects/{id}/concepts [get]
func (api *EncyclopediaApi) ListConcepts(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", nil)
		return
	}

	if api.conceptRepo == nil {
		response.Success(c, []*writerModels.Concept{})
		return
	}

	concepts, err := api.conceptRepo.ListByProject(c.Request.Context(), projectID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	if concepts == nil {
		concepts = []*writerModels.Concept{}
	}
	response.Success(c, concepts)
}

// SearchConcepts 搜索概念
// @Summary 搜索概念
// @Description 在项目下搜索概念
// @Tags 设定百科
// @Produce json
// @Param id path string true "项目ID"
// @Param q query string false "搜索关键词"
// @Success 200 {object} response.APIResponse{data=[]*writerModels.Concept}
// @Router /api/v1/writer/projects/{id}/concepts/search [get]
func (api *EncyclopediaApi) SearchConcepts(c *gin.Context) {
	projectID := c.Param("id")
	keyword := c.Query("q")

	if api.conceptRepo == nil {
		response.Success(c, []*writerModels.Concept{})
		return
	}

	var concepts []*writerModels.Concept
	var err error

	// Search(ctx, projectID, category, keyword)
	raw, repoErr := api.conceptRepo.Search(c.Request.Context(), projectID, "", keyword)
	if repoErr != nil {
		err = repoErr
	} else {
		concepts = raw
	}

	if err != nil {
		response.InternalError(c, err)
		return
	}

	if concepts == nil {
		concepts = []*writerModels.Concept{}
	}
	response.Success(c, concepts)
}

// GetConcept 获取单个概念详情
// @Summary 获取概念详情
// @Description 获取指定概念的详细信息
// @Tags 设定百科
// @Produce json
// @Param id path string true "项目ID"
// @Param conceptId path string true "概念ID"
// @Success 200 {object} response.APIResponse{data=*writerModels.Concept}
// @Router /api/v1/writer/projects/{id}/concepts/{conceptId} [get]
func (api *EncyclopediaApi) GetConcept(c *gin.Context) {
	conceptID := c.Param("conceptId")

	if api.conceptRepo == nil {
		response.NotFound(c, "概念未找到")
		return
	}

	concept, err := api.conceptRepo.GetByID(c.Request.Context(), conceptID)
	if err != nil || concept == nil {
		response.NotFound(c, "概念未找到")
		return
	}

	response.Success(c, concept)
}
