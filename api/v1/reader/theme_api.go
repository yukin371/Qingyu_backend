package reader

import (
	"net/http"

	"Qingyu_backend/api/v1/shared"
	readerModels "Qingyu_backend/models/reader"

	"Qingyu_backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// ThemeAPI 主题API
type ThemeAPI struct {
	// 主题写操作在 service/repository 接入前保持显式未实现，避免返回未持久化的假成功。
}

func respondThemeNotImplemented(c *gin.Context, message string) {
	response.JSON(c, http.StatusNotImplemented, response.APIResponse{
		Code:    http.StatusNotImplemented,
		Message: message,
	})
}

// NewThemeAPI 创建主题API实例
func NewThemeAPI() *ThemeAPI {
	return &ThemeAPI{}
}

// GetThemes 获取可用主题列表
//
//	@Summary	获取可用主题列表
//	@Tags		阅读器-主题
//	@Param		builtin	query	bool	false	"仅显示内置主题"
//	@Param		public	query	bool	false	"仅显示公开主题"
//	@Success	200		{object} response.APIResponse
//	@Router		/api/v1/reader/themes [get]
func (api *ThemeAPI) GetThemes(c *gin.Context) {
	// 获取查询参数
	builtinOnly := c.Query("builtin") == "true"
	publicOnly := c.Query("public") == "true"

	// 获取内置主题（实际应用中应该从数据库查询）
	themes := readerModels.BuiltInThemes

	// 过滤主题
	filteredThemes := make([]*readerModels.ReaderTheme, 0)
	for _, theme := range themes {
		if builtinOnly && !theme.IsBuiltIn {
			continue
		}
		if publicOnly && !theme.IsPublic {
			continue
		}
		filteredThemes = append(filteredThemes, theme)
	}

	response.Success(c, gin.H{
		"themes": filteredThemes,
		"total":  len(filteredThemes),
	})
}

// GetThemeByName 根据名称获取主题
//
//	@Summary	根据名称获取主题
//	@Tags		阅读器-主题
//	@Param		name	path	string	true	"主题名称"
//	@Success	200		{object} response.APIResponse
//	@Router		/api/v1/reader/themes/{name} [get]
func (api *ThemeAPI) GetThemeByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "参数错误", "主题名称不能为空")
		return
	}

	// 从内置主题中查找
	for _, theme := range readerModels.BuiltInThemes {
		if theme.Name == name {
			response.Success(c, theme)
			return
		}
	}

	response.NotFound(c, "主题不存在")
}

// CreateCustomTheme 创建自定义主题
//
//	@Summary	创建自定义主题
//	@Tags		阅读器-主题
//	@Param		request	body object	true	"创建主题请求"
//	@Failure	501		{object} response.APIResponse
//	@Router		/api/v1/reader/themes [post]
func (api *ThemeAPI) CreateCustomTheme(c *gin.Context) {
	// 获取用户ID
	if _, ok := shared.GetUserIDWithMessage(c, "未授权"); !ok {
		return
	}

	var req readerModels.CreateCustomThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 主题 service/repository 尚未接入，不能把内存对象伪装成持久化结果。
	respondThemeNotImplemented(c, "自定义主题持久化功能开发中")
}

// UpdateTheme 更新自定义主题
//
//	@Summary	更新自定义主题
//	@Tags		阅读器-主题
//	@Param		id		path	string						true	"主题ID"
//	@Param		request	body object	true	"更新主题请求"
//	@Failure	501		{object} response.APIResponse
//	@Router		/api/v1/reader/themes/{id} [put]
func (api *ThemeAPI) UpdateTheme(c *gin.Context) {
	if _, ok := shared.GetUserIDWithMessage(c, "未授权"); !ok {
		return
	}

	themeID := c.Param("id")
	if themeID == "" {
		response.BadRequest(c, "参数错误", "主题ID不能为空")
		return
	}

	var req readerModels.UpdateThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 请求和认证校验保留，实际 owner 校验与持久化等待主题 service 接入。
	respondThemeNotImplemented(c, "自定义主题持久化功能开发中")
}

// DeleteTheme 删除自定义主题
//
//	@Summary	删除自定义主题
//	@Tags		阅读器-主题
//	@Param		id	path	string	true	"主题ID"
//	@Failure	501	{object} response.APIResponse
//	@Router		/api/v1/reader/themes/{id} [delete]
func (api *ThemeAPI) DeleteTheme(c *gin.Context) {
	if _, ok := shared.GetUserIDWithMessage(c, "未授权"); !ok {
		return
	}

	themeID := c.Param("id")
	if themeID == "" {
		response.BadRequest(c, "参数错误", "主题ID不能为空")
		return
	}

	// 实际 owner 校验与删除等待主题 service 接入。
	respondThemeNotImplemented(c, "自定义主题持久化功能开发中")
}

// ActivateTheme 激活主题
//
//	@Summary	激活主题（应用到阅读设置）
//	@Tags		阅读器-主题
//	@Param		name	path	string	true	"主题名称"
//	@Failure	501	{object} response.APIResponse
//	@Router		/api/v1/reader/themes/{name}/activate [post]
func (api *ThemeAPI) ActivateTheme(c *gin.Context) {
	if _, ok := shared.GetUserIDWithMessage(c, "未授权"); !ok {
		return
	}

	themeName := c.Param("name")
	if themeName == "" {
		response.BadRequest(c, "参数错误", "主题名称不能为空")
		return
	}

	// 验证主题是否存在
	themeExists := false
	for _, theme := range readerModels.BuiltInThemes {
		if theme.Name == themeName {
			themeExists = true
			break
		}
	}

	if !themeExists {
		response.NotFound(c, "主题不存在")
		return
	}

	// 激活状态和用户阅读设置尚未接入持久化 service。
	respondThemeNotImplemented(c, "主题激活持久化功能开发中")
}
