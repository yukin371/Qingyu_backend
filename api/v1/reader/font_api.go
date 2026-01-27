package reader

import (
	"net/http"

	"Qingyu_backend/api/v1/shared"
	readerModels "Qingyu_backend/models/reader"

	"github.com/gin-gonic/gin"
)

// FontAPI 字体API
type FontAPI struct {
	// 可以注入FontService，暂时先使用内置数据
}

// NewFontAPI 创建字体API实例
func NewFontAPI() *FontAPI {
	return &FontAPI{}
}

// GetFonts 获取可用字体列表
//
//	@Summary	获取可用字体列表
//	@Tags		阅读器-字体
//	@Param		category	query	string	false	"字体分类：serif/sans-serif/monospace"
//	@Param		builtin		query	bool	false	"仅显示内置字体"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/fonts [get]
func (api *FontAPI) GetFonts(c *gin.Context) {
	// 获取查询参数
	category := c.Query("category")
	builtinOnly := c.Query("builtin") == "true"

	// 获取内置字体
	fonts := readerModels.BuiltInFonts

	// 过滤字体
	filteredFonts := make([]*readerModels.ReaderFont, 0)
	for _, font := range fonts {
		if builtinOnly && !font.IsBuiltIn {
			continue
		}
		if category != "" && font.Category != category {
			continue
		}
		if !font.IsActive {
			continue
		}
		filteredFonts = append(filteredFonts, font)
	}

	// 按分类组织
	fontsByCategory := make(map[string][]*readerModels.ReaderFont)
	for _, font := range filteredFonts {
		fontsByCategory[font.Category] = append(fontsByCategory[font.Category], font)
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"fonts":           filteredFonts,
		"total":           len(filteredFonts),
		"fontsByCategory": fontsByCategory,
	})
}

// GetFontByName 根据名称获取字体
//
//	@Summary	根据名称获取字体
//	@Tags		阅读器-字体
//	@Param		name	path	string	true	"字体名称"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/fonts/{name} [get]
func (api *FontAPI) GetFontByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "字体名称不能为空")
		return
	}

	// 从内置字体中查找
	for _, font := range readerModels.BuiltInFonts {
		if font.Name == name {
			shared.Success(c, http.StatusOK, "获取成功", font)
			return
		}
	}

	shared.Error(c, http.StatusNotFound, "字体不存在", "未找到指定字体")
}

// CreateCustomFont 创建自定义字体
//
//	@Summary	创建自定义字体
//	@Tags		阅读器-字体
//	@Param		request	body object	true	"创建字体请求"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/fonts [post]
func (api *FontAPI) CreateCustomFont(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	var req readerModels.CreateCustomFontRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 创建自定义字体
	font := &readerModels.ReaderFont{
		ID:          "", // 应该使用ObjectID生成
		Name:        req.Name,
		DisplayName: req.DisplayName,
		FontFamily:  req.FontFamily,
		Description: req.Description,
		Category:    req.Category,
		FontURL:     req.FontURL,
		PreviewText: req.PreviewText,
		IsBuiltIn:   false,
		IsActive:    true,
		SupportSize: []int{12, 14, 16, 18, 20, 22, 24},
		UseCount:    0,
	}

	// 实际应用中应该保存到数据库
	shared.Success(c, http.StatusCreated, "创建成功", gin.H{
		"font":    font,
		"message": "自定义字体创建成功",
		"userId":  userID,
	})
}

// UpdateFont 更新自定义字体
//
//	@Summary	更新自定义字体
//	@Tags		阅读器-字体
//	@Param		id		path	string						true	"字体ID"
//	@Param		request	body object	true	"更新字体请求"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/fonts/{id} [put]
func (api *FontAPI) UpdateFont(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	fontID := c.Param("id")
	if fontID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "字体ID不能为空")
		return
	}

	var req readerModels.UpdateFontRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 实际应用中应该：
	// 1. 从数据库获取字体
	// 2. 验证用户权限
	// 3. 更新字体

	shared.Success(c, http.StatusOK, "更新成功", gin.H{
		"message": "字体更新成功",
		"fontId":  fontID,
		"userId":  userID,
	})
}

// DeleteFont 删除自定义字体
//
//	@Summary	删除自定义字体
//	@Tags		阅读器-字体
//	@Param		id	path	string	true	"字体ID"
//	@Success	200	{object}	shared.APIResponse
//	@Router		/api/v1/reader/fonts/{id} [delete]
func (api *FontAPI) DeleteFont(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	fontID := c.Param("id")
	if fontID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "字体ID不能为空")
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", gin.H{
		"message": "字体删除成功",
		"fontId":  fontID,
		"userId":  userID,
	})
}

// SetFontPreference 设置字体偏好
//
//	@Summary	设置字体偏好
//	@Tags		阅读器-字体
//	@Param		request	body object	true	"字体偏好"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/settings/font [post]
func (api *FontAPI) SetFontPreference(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	var preference readerModels.FontPreference
	if err := c.ShouldBindJSON(&preference); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 验证字体是否存在
	fontExists := false
	for _, font := range readerModels.BuiltInFonts {
		if font.Name == preference.FontName {
			fontExists = true
			break
		}
	}

	if !fontExists {
		shared.Error(c, http.StatusNotFound, "字体不存在", "未找到指定字体")
		return
	}

	// 设置用户ID
	preference.UserID = userID.(string)

	// 实际应用中应该：
	// 1. 更新用户的阅读设置
	// 2. 保存到数据库
	// 3. 清除缓存

	shared.Success(c, http.StatusOK, "设置成功", gin.H{
		"message":    "字体偏好设置成功",
		"preference": preference,
	})
}
