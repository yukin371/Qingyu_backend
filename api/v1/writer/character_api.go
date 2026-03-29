package writer

import (
	"log"

	"github.com/gin-gonic/gin"

	writerModels "Qingyu_backend/models/writer" // Import for Swagger annotations
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/interfaces"
)

// CharacterApi 角色API处理器
type CharacterApi struct {
	characterService interfaces.CharacterService
}

// NewCharacterApi 创建CharacterApi实例
func NewCharacterApi(characterService interfaces.CharacterService) *CharacterApi {
	return &CharacterApi{
		characterService: characterService,
	}
}

// CreateCharacter 创建角色
// @Summary 创建角色
// @Description 在项目中创建一个新角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Param request body object true "创建角色请求"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/v1/projects/{projectId}/characters [post]
func (api *CharacterApi) CreateCharacter(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", "")
		return
	}

	var req interfaces.CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	character, err := api.characterService.Create(c.Request.Context(), projectID, userID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, character)
}

// GetCharacter 获取角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详细信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param characterId path string true "角色ID"
// @Param projectId query string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/characters/{characterId} [get]
func (api *CharacterApi) GetCharacter(c *gin.Context) {
	characterID := c.Param("characterId")
	projectID := c.Query("projectId")

	if characterID == "" || projectID == "" {
		response.BadRequest(c, "参数错误", "characterId和projectId不能为空")
		return
	}

	character, err := api.characterService.GetByID(c.Request.Context(), characterID, projectID)
	if err != nil {
		response.NotFound(c, "角色不存在")
		return
	}

	response.Success(c, character)
}

// ListCharacters 获取项目角色列表
// @Summary 获取项目角色列表
// @Description 获取指定项目的所有角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/projects/{projectId}/characters [get]
func (api *CharacterApi) ListCharacters(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", "")
		return
	}

	log.Printf("[ListCharacters] 获取项目角色列表, projectID=%s", projectID)

	characters, err := api.characterService.List(c.Request.Context(), projectID)
	if err != nil {
		log.Printf("[ListCharacters] 获取角色列表失败: %v", err)
		c.Error(err)
		return
	}

	log.Printf("[ListCharacters] 成功获取角色, 数量=%d", len(characters))
	if len(characters) > 0 {
		log.Printf("[ListCharacters] 第一个角色: ID=%s, Name=%s", characters[0].ID, characters[0].Name)
	}

	response.Success(c, characters)
}

// UpdateCharacter 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param characterId path string true "角色ID"
// @Param projectId query string true "项目ID"
// @Param request body object true "更新角色请求"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/characters/{characterId} [put]
func (api *CharacterApi) UpdateCharacter(c *gin.Context) {
	characterID := c.Param("characterId")
	projectID := c.Query("projectId")

	if characterID == "" || projectID == "" {
		response.BadRequest(c, "参数错误", "characterId和projectId不能为空")
		return
	}

	var req interfaces.UpdateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	character, err := api.characterService.Update(c.Request.Context(), characterID, projectID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, character)
}

// DeleteCharacter 删除角色
// @Summary 删除角色
// @Description 删除指定角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param characterId path string true "角色ID"
// @Param projectId query string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/characters/{characterId} [delete]
func (api *CharacterApi) DeleteCharacter(c *gin.Context) {
	characterID := c.Param("characterId")
	projectID := c.Query("projectId")

	if characterID == "" || projectID == "" {
		response.BadRequest(c, "参数错误", "characterId和projectId不能为空")
		return
	}

	err := api.characterService.Delete(c.Request.Context(), characterID, projectID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// CreateCharacterRelation 创建角色关系
// @Summary 创建角色关系
// @Description 创建两个角色之间的关系
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param projectId query string true "项目ID"
// @Param request body object true "创建关系请求"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/characters/relations [post]
func (api *CharacterApi) CreateCharacterRelation(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", "")
		return
	}

	var req interfaces.CreateRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	relation, err := api.characterService.CreateRelation(c.Request.Context(), projectID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, relation)
}

// ListCharacterRelations 获取角色关系列表
// @Summary 获取角色关系列表
// @Description 获取项目或指定角色的关系列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Param characterId query string false "角色ID（可选，不传则返回项目所有关系）"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/projects/{projectId}/characters/relations [get]
func (api *CharacterApi) ListCharacterRelations(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", "")
		return
	}

	characterID := c.Query("characterId")
	var charIDPtr *string
	if characterID != "" {
		charIDPtr = &characterID
	}

	log.Printf("[ListCharacterRelations] 获取项目关系列表, projectID=%s, characterID=%s", projectID, characterID)

	relations, err := api.characterService.ListRelations(c.Request.Context(), projectID, charIDPtr)
	if err != nil {
		log.Printf("[ListCharacterRelations] 获取关系列表失败: %v", err)
		c.Error(err)
		return
	}

	log.Printf("[ListCharacterRelations] 成功获取关系, 数量=%d", len(relations))
	if len(relations) > 0 {
		log.Printf("[ListCharacterRelations] 第一个关系: ID=%s, FromID=%s, ToID=%s", relations[0].ID, relations[0].FromID, relations[0].ToID)
	}

	response.Success(c, relations)
}

// DeleteCharacterRelation 删除角色关系
// @Summary 删除角色关系
// @Description 删除指定的角色关系
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param relationId path string true "关系ID"
// @Param projectId query string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/characters/relations/{relationId} [delete]
func (api *CharacterApi) DeleteCharacterRelation(c *gin.Context) {
	relationID := c.Param("relationId")
	projectID := c.Query("projectId")

	if relationID == "" || projectID == "" {
		response.BadRequest(c, "参数错误", "relationId和projectId不能为空")
		return
	}

	err := api.characterService.DeleteRelation(c.Request.Context(), relationID, projectID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// GetCharacterGraph 获取角色关系图
// @Summary 获取角色关系图
// @Description 获取项目的角色关系图（包含所有角色和关系）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/projects/{projectId}/characters/graph [get]
func (api *CharacterApi) GetCharacterGraph(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", "")
		return
	}

	graph, err := api.characterService.GetCharacterGraph(c.Request.Context(), projectID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, graph)
}

// CreateRelationTimelineEvent 创建关系时序变化事件
// @Summary 创建关系时序变化事件
// @Description 在指定章节创建关系时序变化事件
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param relationId path string true "关系ID"
// @Param projectId query string true "项目ID"
// @Param request body object true "时序事件请求"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/characters/relations/{relationId}/timeline [post]
func (api *CharacterApi) CreateRelationTimelineEvent(c *gin.Context) {
	relationID := c.Param("relationId")
	projectID := c.Query("projectId")

	if relationID == "" || projectID == "" {
		response.BadRequest(c, "参数错误", "relationId和projectId不能为空")
		return
	}

	var req interfaces.CreateRelationTimelineEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 确保relationId一致
	req.RelationID = relationID

	event, err := api.characterService.CreateRelationTimelineEvent(c.Request.Context(), projectID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, event)
}

// GetRelationTimeline 获取关系时序历史
// @Summary 获取关系时序历史
// @Description 获取指定关系的所有时序变化事件
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param relationId path string true "关系ID"
// @Param projectId query string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/characters/relations/{relationId}/timeline [get]
func (api *CharacterApi) GetRelationTimeline(c *gin.Context) {
	relationID := c.Param("relationId")
	projectID := c.Query("projectId")

	if relationID == "" || projectID == "" {
		response.BadRequest(c, "参数错误", "relationId和projectId不能为空")
		return
	}

	events, err := api.characterService.GetRelationTimeline(c.Request.Context(), relationID, projectID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, events)
}

// UpdateRelationTimelineEvent 更新关系时序事件
// @Summary 更新关系时序事件
// @Description 更新指定关系时序事件
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param eventId path string true "事件ID"
// @Param projectId query string true "项目ID"
// @Param request body object true "更新请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/characters/relations/timeline-events/{eventId} [put]
func (api *CharacterApi) UpdateRelationTimelineEvent(c *gin.Context) {
	eventID := c.Param("eventId")
	projectID := c.Query("projectId")

	if eventID == "" || projectID == "" {
		response.BadRequest(c, "参数错误", "eventId和projectId不能为空")
		return
	}

	var req interfaces.UpdateRelationTimelineEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	event, err := api.characterService.UpdateRelationTimelineEvent(c.Request.Context(), eventID, projectID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, event)
}

// DeleteRelationTimelineEvent 删除关系时序事件
// @Summary 删除关系时序事件
// @Description 删除指定关系时序事件
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param eventId path string true "事件ID"
// @Param projectId query string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/characters/relations/timeline-events/{eventId} [delete]
func (api *CharacterApi) DeleteRelationTimelineEvent(c *gin.Context) {
	eventID := c.Param("eventId")
	projectID := c.Query("projectId")

	if eventID == "" || projectID == "" {
		response.BadRequest(c, "参数错误", "eventId和projectId不能为空")
		return
	}

	err := api.characterService.DeleteRelationTimelineEvent(c.Request.Context(), eventID, projectID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

var _ = writerModels.Character{}
