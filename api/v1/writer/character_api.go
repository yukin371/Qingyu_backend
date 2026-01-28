package writer

import (
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
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/projects/{projectId}/characters [post]
func (api *CharacterApi) CreateCharacter(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.BadRequest(c,  "项目ID不能为空", "")
		return
	}

	var req interfaces.CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	character, err := api.characterService.Create(c.Request.Context(), projectID, userID, &req)
	if err != nil {
		response.InternalError(c, err)
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
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/characters/{characterId} [get]
func (api *CharacterApi) GetCharacter(c *gin.Context) {
	characterID := c.Param("characterId")
	projectID := c.Query("projectId")

	if characterID == "" || projectID == "" {
		response.BadRequest(c,  "参数错误", "characterId和projectId不能为空")
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
// @Success 200 {object} response.Response
// @Router /api/v1/projects/{projectId}/characters [get]
func (api *CharacterApi) ListCharacters(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.BadRequest(c,  "项目ID不能为空", "")
		return
	}

	characters, err := api.characterService.List(c.Request.Context(), projectID)
	if err != nil {
		response.InternalError(c, err)
		return
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
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/characters/{characterId} [put]
func (api *CharacterApi) UpdateCharacter(c *gin.Context) {
	characterID := c.Param("characterId")
	projectID := c.Query("projectId")

	if characterID == "" || projectID == "" {
		response.BadRequest(c,  "参数错误", "characterId和projectId不能为空")
		return
	}

	var req interfaces.UpdateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	character, err := api.characterService.Update(c.Request.Context(), characterID, projectID, &req)
	if err != nil {
		response.InternalError(c, err)
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
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/characters/{characterId} [delete]
func (api *CharacterApi) DeleteCharacter(c *gin.Context) {
	characterID := c.Param("characterId")
	projectID := c.Query("projectId")

	if characterID == "" || projectID == "" {
		response.BadRequest(c,  "参数错误", "characterId和projectId不能为空")
		return
	}

	err := api.characterService.Delete(c.Request.Context(), characterID, projectID)
	if err != nil {
		response.InternalError(c, err)
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
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/characters/relations [post]
func (api *CharacterApi) CreateCharacterRelation(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		response.BadRequest(c,  "项目ID不能为空", "")
		return
	}

	var req interfaces.CreateRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	relation, err := api.characterService.CreateRelation(c.Request.Context(), projectID, &req)
	if err != nil {
		response.InternalError(c, err)
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
// @Success 200 {object} response.Response
// @Router /api/v1/projects/{projectId}/characters/relations [get]
func (api *CharacterApi) ListCharacterRelations(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.BadRequest(c,  "项目ID不能为空", "")
		return
	}

	characterID := c.Query("characterId")
	var charIDPtr *string
	if characterID != "" {
		charIDPtr = &characterID
	}

	relations, err := api.characterService.ListRelations(c.Request.Context(), projectID, charIDPtr)
	if err != nil {
		response.InternalError(c, err)
		return
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
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/characters/relations/{relationId} [delete]
func (api *CharacterApi) DeleteCharacterRelation(c *gin.Context) {
	relationID := c.Param("relationId")
	projectID := c.Query("projectId")

	if relationID == "" || projectID == "" {
		response.BadRequest(c,  "参数错误", "relationId和projectId不能为空")
		return
	}

	err := api.characterService.DeleteRelation(c.Request.Context(), relationID, projectID)
	if err != nil {
		response.InternalError(c, err)
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
// @Success 200 {object} response.Response
// @Router /api/v1/projects/{projectId}/characters/graph [get]
func (api *CharacterApi) GetCharacterGraph(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.BadRequest(c,  "项目ID不能为空", "")
		return
	}

	graph, err := api.characterService.GetCharacterGraph(c.Request.Context(), projectID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, graph)
}

var _ = writerModels.Character{}
