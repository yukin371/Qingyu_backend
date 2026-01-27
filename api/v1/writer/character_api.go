package writer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	writerModels "Qingyu_backend/models/writer" // Import for Swagger annotations
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
// @Success 201 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 401 {object} shared.APIResponse
// @Router /api/v1/projects/{projectId}/characters [post]
func (api *CharacterApi) CreateCharacter(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "项目ID不能为空", "")
		return
	}

	var req interfaces.CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
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
		shared.Error(c, http.StatusInternalServerError, "创建角色失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "创建成功", character)
}

// GetCharacter 获取角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详细信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param characterId path string true "角色ID"
// @Param projectId query string true "项目ID"
// @Success 200 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/characters/{characterId} [get]
func (api *CharacterApi) GetCharacter(c *gin.Context) {
	characterID := c.Param("characterId")
	projectID := c.Query("projectId")

	if characterID == "" || projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "characterId和projectId不能为空")
		return
	}

	character, err := api.characterService.GetByID(c.Request.Context(), characterID, projectID)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "角色不存在", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", character)
}

// ListCharacters 获取项目角色列表
// @Summary 获取项目角色列表
// @Description 获取指定项目的所有角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/projects/{projectId}/characters [get]
func (api *CharacterApi) ListCharacters(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "项目ID不能为空", "")
		return
	}

	characters, err := api.characterService.List(c.Request.Context(), projectID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取角色列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", characters)
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
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/characters/{characterId} [put]
func (api *CharacterApi) UpdateCharacter(c *gin.Context) {
	characterID := c.Param("characterId")
	projectID := c.Query("projectId")

	if characterID == "" || projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "characterId和projectId不能为空")
		return
	}

	var req interfaces.UpdateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	character, err := api.characterService.Update(c.Request.Context(), characterID, projectID, &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新角色失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", character)
}

// DeleteCharacter 删除角色
// @Summary 删除角色
// @Description 删除指定角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param characterId path string true "角色ID"
// @Param projectId query string true "项目ID"
// @Success 200 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/characters/{characterId} [delete]
func (api *CharacterApi) DeleteCharacter(c *gin.Context) {
	characterID := c.Param("characterId")
	projectID := c.Query("projectId")

	if characterID == "" || projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "characterId和projectId不能为空")
		return
	}

	err := api.characterService.Delete(c.Request.Context(), characterID, projectID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "删除角色失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}

// CreateCharacterRelation 创建角色关系
// @Summary 创建角色关系
// @Description 创建两个角色之间的关系
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param projectId query string true "项目ID"
// @Param request body object true "创建关系请求"
// @Success 201 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Router /api/v1/characters/relations [post]
func (api *CharacterApi) CreateCharacterRelation(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "项目ID不能为空", "")
		return
	}

	var req interfaces.CreateRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	relation, err := api.characterService.CreateRelation(c.Request.Context(), projectID, &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "创建关系失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "创建成功", relation)
}

// ListCharacterRelations 获取角色关系列表
// @Summary 获取角色关系列表
// @Description 获取项目或指定角色的关系列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Param characterId query string false "角色ID（可选，不传则返回项目所有关系）"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/projects/{projectId}/characters/relations [get]
func (api *CharacterApi) ListCharacterRelations(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "项目ID不能为空", "")
		return
	}

	characterID := c.Query("characterId")
	var charIDPtr *string
	if characterID != "" {
		charIDPtr = &characterID
	}

	relations, err := api.characterService.ListRelations(c.Request.Context(), projectID, charIDPtr)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取关系列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", relations)
}

// DeleteCharacterRelation 删除角色关系
// @Summary 删除角色关系
// @Description 删除指定的角色关系
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param relationId path string true "关系ID"
// @Param projectId query string true "项目ID"
// @Success 200 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/characters/relations/{relationId} [delete]
func (api *CharacterApi) DeleteCharacterRelation(c *gin.Context) {
	relationID := c.Param("relationId")
	projectID := c.Query("projectId")

	if relationID == "" || projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "relationId和projectId不能为空")
		return
	}

	err := api.characterService.DeleteRelation(c.Request.Context(), relationID, projectID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "删除关系失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}

// GetCharacterGraph 获取角色关系图
// @Summary 获取角色关系图
// @Description 获取项目的角色关系图（包含所有角色和关系）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/projects/{projectId}/characters/graph [get]
func (api *CharacterApi) GetCharacterGraph(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "项目ID不能为空", "")
		return
	}

	graph, err := api.characterService.GetCharacterGraph(c.Request.Context(), projectID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取关系图失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", graph)
}

var _ = writerModels.Character{}
