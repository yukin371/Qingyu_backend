package reader

import (
	"net/http"
	"strconv"

	"Qingyu_backend/api/v1/shared"
	readerModels "Qingyu_backend/models/reader"

	"github.com/gin-gonic/gin"
	"Qingyu_backend/pkg/response"
)

// ThemeAPI 主题API
type ThemeAPI struct {
	// 可以注入ThemeService，暂时先使用内置数据
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
//	@Success	200		{object}	shared.APIResponse
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

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"themes": filteredThemes,
		"total":  len(filteredThemes),
	})
}

// GetThemeByName 根据名称获取主题
//
//	@Summary	根据名称获取主题
//	@Tags		阅读器-主题
//	@Param		name	path	string	true	"主题名称"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/themes/{name} [get]
func (api *ThemeAPI) GetThemeByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c,  "参数错误", "主题名称不能为空")
		return
	}

	// 从内置主题中查找
	for _, theme := range readerModels.BuiltInThemes {
		if theme.Name == name {
			shared.Success(c, http.StatusOK, "获取成功", theme)
			return
		}
	}

	shared.Error(c, http.StatusNotFound, "主题不存在", "未找到指定主题")
}

// CreateCustomTheme 创建自定义主题
//
//	@Summary	创建自定义主题
//	@Tags		阅读器-主题
//	@Param		request	body object	true	"创建主题请求"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/themes [post]
func (api *ThemeAPI) CreateCustomTheme(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	var req readerModels.CreateCustomThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 创建自定义主题
	theme := &readerModels.ReaderTheme{
		ID:          generateID(), // 应该使用ObjectID生成
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IsBuiltIn:   false,
		IsPublic:    req.IsPublic,
		CreatorID:   userID.(string),
		Colors:      req.Colors,
		IsActive:    true,
		UseCount:    0,
	}

	// 实际应用中应该保存到数据库
	// 这里暂时返回成功
	shared.Success(c, http.StatusCreated, "创建成功", gin.H{
		"theme":   theme,
		"message": "自定义主题创建成功",
	})
}

// UpdateTheme 更新自定义主题
//
//	@Summary	更新自定义主题
//	@Tags		阅读器-主题
//	@Param		id		path	string						true	"主题ID"
//	@Param		request	body object	true	"更新主题请求"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/themes/{id} [put]
func (api *ThemeAPI) UpdateTheme(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	themeID := c.Param("id")
	if themeID == "" {
		response.BadRequest(c,  "参数错误", "主题ID不能为空")
		return
	}

	var req readerModels.UpdateThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 实际应用中应该：
	// 1. 从数据库获取主题
	// 2. 验证用户是否为创建者
	// 3. 更新主题
	// 4. 保存到数据库

	shared.Success(c, http.StatusOK, "更新成功", gin.H{
		"message": "主题更新成功",
		"themeId": themeID,
		"userId":  userID,
	})
}

// DeleteTheme 删除自定义主题
//
//	@Summary	删除自定义主题
//	@Tags		阅读器-主题
//	@Param		id	path	string	true	"主题ID"
//	@Success	200	{object}	shared.APIResponse
//	@Router		/api/v1/reader/themes/{id} [delete]
func (api *ThemeAPI) DeleteTheme(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	themeID := c.Param("id")
	if themeID == "" {
		response.BadRequest(c,  "参数错误", "主题ID不能为空")
		return
	}

	// 实际应用中应该：
	// 1. 从数据库获取主题
	// 2. 验证用户是否为创建者
	// 3. 删除主题

	shared.Success(c, http.StatusOK, "删除成功", gin.H{
		"message": "主题删除成功",
		"themeId": themeID,
		"userId":  userID,
	})
}

// ActivateTheme 激活主题
//
//	@Summary	激活主题（应用到阅读设置）
//	@Tags		阅读器-主题
//	@Param		name	path	string	true	"主题名称"
//	@Success	200	{object}	shared.APIResponse
//	@Router		/api/v1/reader/themes/{name}/activate [post]
func (api *ThemeAPI) ActivateTheme(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	themeName := c.Param("name")
	if themeName == "" {
		response.BadRequest(c,  "参数错误", "主题名称不能为空")
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
		shared.Error(c, http.StatusNotFound, "主题不存在", "未找到指定主题")
		return
	}

	// 实际应用中应该：
	// 1. 更新用户的阅读设置，将theme字段设置为themeName
	// 2. 清除设置缓存

	shared.Success(c, http.StatusOK, "激活成功", gin.H{
		"message":   "主题已激活",
		"themeName": themeName,
		"userId":    userID,
	})
}

// generateID 生成ID（临时方法，实际应使用primitive.NewObjectID）
func generateID() string {
	return strconv.FormatInt(int64(len("temp")), 10)
}
