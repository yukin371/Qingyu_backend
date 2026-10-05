package reader

import (
	"net/http"
	"time"

	"Qingyu_backend/api/v1/shared"
	readerModels "Qingyu_backend/models/reader"

	"Qingyu_backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func respondFontNotImplemented(c *gin.Context, message string) {
	response.JSON(c, http.StatusNotImplemented, response.APIResponse{
		Code:      response.CodeInternalError,
		Message:   message,
		Timestamp: time.Now().UnixMilli(),
		RequestID: c.GetString("requestId"),
	})
}

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
//	@Success	200		{object}	response.APIResponse
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

	response.Success(c, gin.H{
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
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/fonts/{name} [get]
func (api *FontAPI) GetFontByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "参数错误", "字体名称不能为空")
		return
	}

	// 从内置字体中查找
	for _, font := range readerModels.BuiltInFonts {
		if font.Name == name {
			response.Success(c, font)
			return
		}
	}

	response.NotFound(c, "字体不存在")
}

// CreateCustomFont 创建自定义字体
//
//	@Summary	创建自定义字体
//	@Tags		阅读器-字体
//	@Param		request	body object	true	"创建字体请求"
//	@Success	501		{object}	response.APIResponse
//	@Router		/api/v1/reader/fonts [post]
func (api *FontAPI) CreateCustomFont(c *gin.Context) {
	// 获取用户ID
	_, ok := shared.GetUserIDWithMessage(c, "请先登录")
	if !ok {
		return
	}

	var req readerModels.CreateCustomFontRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	respondFontNotImplemented(c, "自定义字体持久化功能开发中")
}

// TODO(T3-018-security): 在确认字体归属规则并补 user-scope service 契约后，执行当前用户 owner 校验。
// UpdateFont 更新自定义字体
//
//	@Summary	更新自定义字体
//	@Tags		阅读器-字体
//	@Param		id		path	string						true	"字体ID"
//	@Param		request	body object	true	"更新字体请求"
//	@Success	501		{object}	response.APIResponse
//	@Router		/api/v1/reader/fonts/{id} [put]
func (api *FontAPI) UpdateFont(c *gin.Context) {
	_, ok := shared.GetUserIDWithMessage(c, "请先登录")
	if !ok {
		return
	}

	fontID := c.Param("id")
	if fontID == "" {
		response.BadRequest(c, "参数错误", "字体ID不能为空")
		return
	}

	var req readerModels.UpdateFontRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	respondFontNotImplemented(c, "自定义字体持久化功能开发中")
}

// TODO(T3-018-security): 在确认字体归属规则并补 user-scope service 契约后，执行当前用户 owner 校验。
// DeleteFont 删除自定义字体
//
//	@Summary	删除自定义字体
//	@Tags		阅读器-字体
//	@Param		id	path	string	true	"字体ID"
//	@Success	501	{object}	response.APIResponse
//	@Router		/api/v1/reader/fonts/{id} [delete]
func (api *FontAPI) DeleteFont(c *gin.Context) {
	_, ok := shared.GetUserIDWithMessage(c, "请先登录")
	if !ok {
		return
	}

	fontID := c.Param("id")
	if fontID == "" {
		response.BadRequest(c, "参数错误", "字体ID不能为空")
		return
	}

	respondFontNotImplemented(c, "自定义字体持久化功能开发中")
}

// SetFontPreference 设置字体偏好
//
//	@Summary	设置字体偏好
//	@Tags		阅读器-字体
//	@Param		request	body object	true	"字体偏好"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/settings/font [post]
func (api *FontAPI) SetFontPreference(c *gin.Context) {
	userID, ok := shared.GetUserIDWithMessage(c, "请先登录")
	if !ok {
		return
	}

	var preference readerModels.FontPreference
	if err := c.ShouldBindJSON(&preference); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
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
		response.NotFound(c, "字体不存在")
		return
	}

	// 设置用户ID
	preference.UserID = userID

	// 实际应用中应该：
	// 1. 更新用户的阅读设置
	// 2. 保存到数据库
	// 3. 清除缓存

	response.Success(c, gin.H{
		"message":    "字体偏好设置成功",
		"preference": preference,
	})
}
