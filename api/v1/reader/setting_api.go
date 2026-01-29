package reader

import (
	readerModels "Qingyu_backend/models/reader"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/service/interfaces"
	"Qingyu_backend/pkg/response"
	"errors"
)

// SettingAPI 设置API
type SettingAPI struct {
	readerService interfaces.ReaderService
}

// NewSettingAPI 创建设置API实例
func NewSettingAPI(readerService interfaces.ReaderService) *SettingAPI {
	return &SettingAPI{
		readerService: readerService,
	}
}

// UpdateSettingsRequest 更新设置请求
type UpdateSettingsRequest struct {
	FontSize        *int     `json:"fontSize"`
	FontFamily      *string  `json:"fontFamily"`
	LineHeight      *float64 `json:"lineHeight"`
	BackgroundColor *string  `json:"backgroundColor"`
	TextColor       *string  `json:"textColor"`
	PageMode        *string  `json:"pageMode"` // scroll, paginate
	AutoSave        *bool    `json:"autoSave"`
	ShowProgress    *bool    `json:"showProgress"`
}

// GetReadingSettings 获取阅读设置
//
//	@Summary	获取阅读设置
//	@Tags		阅读器
//	@Success	200	{object}	response.APIResponse
//	@Router		/api/v1/reader/settings [get]
func (api *SettingAPI) GetReadingSettings(c *gin.Context) {
	// 检查服务是否初始化
	if api.readerService == nil {
		response.InternalError(c, errors.New("服务未初始化: 阅读器服务未正确初始化"))
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 类型断言安全检查
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		response.BadRequest(c, "参数错误", "无效的用户ID")
		return
	}

	settings, err := api.readerService.GetReadingSettings(c.Request.Context(), userIDStr)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, settings)
}

// SaveReadingSettings 保存阅读设置
//
//	@Summary	保存阅读设置
//	@Tags		阅读器
//	@Param		request	body object	true	"阅读设置"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/settings [post]
func (api *SettingAPI) SaveReadingSettings(c *gin.Context) {
	// 检查服务是否初始化
	if api.readerService == nil {
		response.InternalError(c, errors.New("服务未初始化: 阅读器服务未正确初始化"))
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 类型断言安全检查
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		response.BadRequest(c, "参数错误", "无效的用户ID")
		return
	}

	var settings readerModels.ReadingSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	settings.UserID = userIDStr

	err := api.readerService.SaveReadingSettings(c.Request.Context(), &settings)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// UpdateReadingSettings 更新阅读设置
//
//	@Summary	更新阅读设置
//	@Tags		阅读器
//	@Param		request	body		UpdateSettingsRequest	true	"更新设置请求"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/settings [put]
func (api *SettingAPI) UpdateReadingSettings(c *gin.Context) {
	// 检查服务是否初始化
	if api.readerService == nil {
		response.InternalError(c, errors.New("服务未初始化: 阅读器服务未正确初始化"))
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 类型断言安全检查
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		response.BadRequest(c, "参数错误", "无效的用户ID")
		return
	}

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.FontSize != nil {
		updates["font_size"] = *req.FontSize
	}
	if req.FontFamily != nil {
		updates["font_family"] = *req.FontFamily
	}
	if req.LineHeight != nil {
		updates["line_height"] = *req.LineHeight
	}
	if req.BackgroundColor != nil {
		updates["background_color"] = *req.BackgroundColor
	}
	if req.TextColor != nil {
		updates["text_color"] = *req.TextColor
	}
	if req.PageMode != nil {
		updates["page_mode"] = *req.PageMode
	}
	if req.AutoSave != nil {
		updates["auto_save"] = *req.AutoSave
	}
	if req.ShowProgress != nil {
		updates["show_progress"] = *req.ShowProgress
	}

	err := api.readerService.UpdateReadingSettings(c.Request.Context(), userIDStr, updates)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}
