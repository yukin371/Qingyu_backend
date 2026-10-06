package reader

import (
	readerModels "Qingyu_backend/models/reader"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/interfaces"
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
	Theme           *string  `json:"theme"`
	BackgroundColor *string  `json:"backgroundColor"`
	TextColor       *string  `json:"textColor"`
	PageMode        *string  `json:"pageMode"` // scroll, paginate
	AutoScroll      *bool    `json:"autoScroll"`
	ScrollSpeed     *int     `json:"scrollSpeed"`
	AutoSave        *bool    `json:"autoSave"`
	ShowProgress    *bool    `json:"showProgress"`
}

// buildReadingSettingsUpdates converts the HTTP request contract into the
// persistence-owned field names used by ReaderService.
func buildReadingSettingsUpdates(req UpdateSettingsRequest) (map[string]interface{}, error) {
	if req.TextColor != nil || req.AutoSave != nil || req.ShowProgress != nil {
		return nil, fmt.Errorf("textColor、autoSave、showProgress 尚未由阅读设置后端持久化，请仅在本地保存")
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
	if req.Theme != nil {
		updates["theme"] = *req.Theme
	}
	if req.BackgroundColor != nil {
		updates["background"] = *req.BackgroundColor
	}
	if req.AutoScroll != nil {
		updates["auto_scroll"] = *req.AutoScroll
	}
	if req.ScrollSpeed != nil {
		updates["scroll_speed"] = *req.ScrollSpeed
	}
	if req.PageMode != nil {
		switch strings.ToLower(strings.TrimSpace(*req.PageMode)) {
		case "scroll", "1":
			updates["page_mode"] = 1
		case "paginate", "2":
			updates["page_mode"] = 2
		default:
			return nil, fmt.Errorf("pageMode 必须是 scroll 或 paginate")
		}
	}

	return updates, nil
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
		c.Error(errors.New("服务未初始化: 阅读器服务未正确初始化"))
		return
	}

	// 获取并校验用户ID
	userIDStr, ok := shared.GetUserIDWithMessage(c, "请先登录")
	if !ok || userIDStr == "" {
		return
	}

	settings, err := api.readerService.GetReadingSettings(c.Request.Context(), userIDStr)
	if err != nil {
		c.Error(err)
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
		c.Error(errors.New("服务未初始化: 阅读器服务未正确初始化"))
		return
	}

	// 获取并校验用户ID
	userIDStr, ok := shared.GetUserIDWithMessage(c, "请先登录")
	if !ok || userIDStr == "" {
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
		c.Error(err)
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
//	@Router		/api/v1/reader/settings [patch]
func (api *SettingAPI) UpdateReadingSettings(c *gin.Context) {
	// 检查服务是否初始化
	if api.readerService == nil {
		c.Error(errors.New("服务未初始化: 阅读器服务未正确初始化"))
		return
	}

	// 获取并校验用户ID
	userIDStr, ok := shared.GetUserIDWithMessage(c, "请先登录")
	if !ok || userIDStr == "" {
		return
	}

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	updates, err := buildReadingSettingsUpdates(req)
	if err != nil {
		response.BadRequest(c, "阅读设置字段不受支持", err.Error())
		return
	}

	err = api.readerService.UpdateReadingSettings(c.Request.Context(), userIDStr, updates)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}
