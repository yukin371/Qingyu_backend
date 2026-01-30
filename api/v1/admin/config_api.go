package admin

import (
	"net/http"

	"Qingyu_backend/api/v1/shared"
	sharedService "Qingyu_backend/service/shared"
	"Qingyu_backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// ConfigAPI 配置管理API
type ConfigAPI struct {
	configService *sharedService.ConfigService
}

// NewConfigAPI 创建配置管理API实例
func NewConfigAPI(configService *sharedService.ConfigService) *ConfigAPI {
	return &ConfigAPI{
		configService: configService,
	}
}

// GetAllConfigsResponse 获取所有配置响应
type GetAllConfigsResponse struct {
	Groups []*sharedService.ConfigGroup `json:"groups"`
}

// GetAllConfigs 获取所有配置
// @Summary 获取所有配置
// @Description 获取系统所有配置项（分组显示）
// @Tags 管理员-配置管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} shared.APIResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/config [get]
func (api *ConfigAPI) GetAllConfigs(c *gin.Context) {
	groups, err := api.configService.GetAllConfigs(c.Request.Context())
	if err != nil {
		shared.InternalError(c, "获取配置失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "获取配置成功", GetAllConfigsResponse{
		Groups: groups,
	})
}

// GetConfigByKeyRequest 根据Key获取配置请求
type GetConfigByKeyRequest struct {
	Key string `uri:"key" binding:"required"`
}

// GetConfigByKey 根据Key获取单个配置
// @Summary 根据Key获取单个配置
// @Description 获取指定的配置项
// @Tags 管理员-配置管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param key path string true "配置键（如 server.port）"
// @Success 200 {object} shared.APIResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/config/{key} [get]
func (api *ConfigAPI) GetConfigByKey(c *gin.Context) {
	var req GetConfigByKeyRequest
	if err := c.ShouldBindUri(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	item, err := api.configService.GetConfigByKey(c.Request.Context(), req.Key)
	if err != nil {
		shared.NotFound(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "获取配置成功", item)
}

// UpdateConfigRequest 更新配置请求
type UpdateConfigRequest struct {
	Key   string      `json:"key" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}

// UpdateConfig 更新配置
// @Summary 更新配置
// @Description 更新指定的配置项（需要管理员权限）
// @Tags 管理员-配置管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body UpdateConfigRequest true "更新配置请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 403 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/config [put]
func (api *ConfigAPI) UpdateConfig(c *gin.Context) {
	var req UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 转换为Service层请求
	serviceReq := &sharedService.UpdateConfigRequest{
		Key:   req.Key,
		Value: req.Value,
	}

	if err := api.configService.UpdateConfig(c.Request.Context(), serviceReq); err != nil {
		shared.InternalError(c, "更新配置失败", err)
		return
	}

	response.SuccessWithMessage(c, "配置更新成功", nil)
}

// BatchUpdateConfigRequest 批量更新配置请求
type BatchUpdateConfigRequest struct {
	Updates []UpdateConfigRequest `json:"updates" binding:"required,min=1"`
}

// BatchUpdateConfig 批量更新配置
// @Summary 批量更新配置
// @Description 批量更新多个配置项（需要管理员权限）
// @Tags 管理员-配置管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body BatchUpdateConfigRequest true "批量更新配置请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 403 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/config/batch [put]
func (api *ConfigAPI) BatchUpdateConfig(c *gin.Context) {
	var req BatchUpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 转换为Service层请求
	var serviceReqs []*sharedService.UpdateConfigRequest
	for _, update := range req.Updates {
		serviceReqs = append(serviceReqs, &sharedService.UpdateConfigRequest{
			Key:   update.Key,
			Value: update.Value,
		})
	}

	if err := api.configService.BatchUpdateConfig(c.Request.Context(), serviceReqs); err != nil {
		shared.InternalError(c, "批量更新配置失败", err)
		return
	}

	response.SuccessWithMessage(c, "配置批量更新成功", nil)
}

// ValidateConfigRequest 验证配置请求
type ValidateConfigRequest struct {
	YAMLContent string `json:"yaml_content" binding:"required"`
}

// ValidateConfig 验证配置
// @Summary 验证配置
// @Description 验证YAML配置格式是否正确
// @Tags 管理员-配置管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body ValidateConfigRequest true "验证配置请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Router /api/v1/admin/config/validate [post]
func (api *ConfigAPI) ValidateConfig(c *gin.Context) {
	var req ValidateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.configService.ValidateConfig(c.Request.Context(), req.YAMLContent); err != nil {
		shared.BadRequest(c, "配置验证失败", err.Error())
		return
	}

	response.SuccessWithMessage(c, "配置验证通过", nil)
}

// GetConfigBackupsResponse 获取配置备份响应
type GetConfigBackupsResponse struct {
	Backups []string `json:"backups"`
}

// GetConfigBackups 获取配置备份列表
// @Summary 获取配置备份列表
// @Description 获取所有可用的配置备份
// @Tags 管理员-配置管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} shared.APIResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/config/backups [get]
func (api *ConfigAPI) GetConfigBackups(c *gin.Context) {
	backups, err := api.configService.GetConfigBackups(c.Request.Context())
	if err != nil {
		shared.InternalError(c, "获取备份列表失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "获取备份列表成功", GetConfigBackupsResponse{
		Backups: backups,
	})
}

// RestoreConfigBackup 恢复配置备份
// @Summary 恢复配置备份
// @Description 将配置恢复到最近的备份（需要管理员权限）
// @Tags 管理员-配置管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} shared.APIResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 403 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/config/restore [post]
func (api *ConfigAPI) RestoreConfigBackup(c *gin.Context) {
	if err := api.configService.RestoreConfigBackup(c.Request.Context()); err != nil {
		shared.InternalError(c, "恢复配置失败", err)
		return
	}

	response.SuccessWithMessage(c, "配置恢复成功", nil)
}
