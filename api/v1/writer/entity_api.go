package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/interfaces"
)

// EntityApi 统一实体API处理器
type EntityApi struct {
	entityService interfaces.EntityService
}

func getEntityProjectID(c *gin.Context) string {
	if projectID := c.Param("projectId"); projectID != "" {
		return projectID
	}
	return c.Param("id")
}

// NewEntityApi 创建EntityApi实例
func NewEntityApi(entityService interfaces.EntityService) *EntityApi {
	return &EntityApi{
		entityService: entityService,
	}
}

// ListEntities 获取项目下所有实体
// @Summary 获取项目实体列表
// @Description 查询项目下所有实体，支持按类型筛选
// @Tags 实体管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Param type query string false "实体类型筛选（character/item/location）"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/writer/projects/{projectId}/entities [get]
func (api *EntityApi) ListEntities(c *gin.Context) {
	projectID := getEntityProjectID(c)
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", "")
		return
	}

	// 可选的实体类型筛选
	var entityType *string
	if t := c.Query("type"); t != "" {
		et := writer.EntityType(t)
		if !et.IsValid() {
			response.BadRequest(c, "无效的实体类型", "可选值: character, item, location, organization, foreshadowing")
			return
		}
		entityType = &t
	}

	summaries, err := api.entityService.ListEntities(c.Request.Context(), projectID, entityType)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, summaries)
}

// GetEntityGraph 获取项目实体图谱
// @Summary 获取项目实体图谱
// @Description 获取项目下所有实体及其关系，以图谱形式返回
// @Tags 实体管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/writer/projects/{projectId}/entities/graph [get]
func (api *EntityApi) GetEntityGraph(c *gin.Context) {
	projectID := getEntityProjectID(c)
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", "")
		return
	}

	graph, err := api.entityService.GetEntityGraph(c.Request.Context(), projectID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, graph)
}

// UpdateEntityStateFields 更新实体状态字段
// @Summary 更新实体状态字段
// @Description 更新指定实体的状态字段（如生命值、心情等）
// @Tags 实体管理
// @Accept json
// @Produce json
// @Param entityId path string true "实体ID"
// @Param request body map[string]writer.StateValue true "状态字段"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/writer/entities/{entityId}/state-fields [put]
func (api *EntityApi) UpdateEntityStateFields(c *gin.Context) {
	entityID := c.Param("entityId")
	if entityID == "" {
		response.BadRequest(c, "实体ID不能为空", "")
		return
	}

	var stateFields map[string]writer.StateValue
	if !shared.BindJSON(c, &stateFields) {
		return
	}

	if len(stateFields) == 0 {
		response.BadRequest(c, "状态字段不能为空", "")
		return
	}

	err := api.entityService.UpdateEntityStateFields(c.Request.Context(), entityID, stateFields)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "状态字段更新成功"})
}
