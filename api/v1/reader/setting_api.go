package reader

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/reading/reader"
	"Qingyu_backend/service/reading"
)

// SettingAPI 设置API
type SettingAPI struct {
	readerService *reading.ReaderService
}

// NewSettingAPI 创建设置API实例
func NewSettingAPI(readerService *reading.ReaderService) *SettingAPI {
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
//	@Success	200	{object}	response.Response
//	@Router		/api/v1/reader/settings [get]
func (api *SettingAPI) GetReadingSettings(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	settings, err := api.readerService.GetReadingSettings(c.Request.Context(), userID.(string))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取阅读设置失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", settings)
}

// SaveReadingSettings 保存阅读设置
//
//	@Summary	保存阅读设置
//	@Tags		阅读器
//	@Param		request	body		reader.ReadingSettings	true	"阅读设置"
//	@Success	200		{object}	response.Response
//	@Router		/api/v1/reader/settings [post]
func (api *SettingAPI) SaveReadingSettings(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	var settings reader.ReadingSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		shared.ValidationError(c, err)
		return
	}

	settings.UserID = userID.(string)

	err := api.readerService.SaveReadingSettings(c.Request.Context(), &settings)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "保存阅读设置失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "保存成功", nil)
}

// UpdateReadingSettings 更新阅读设置
//
//	@Summary	更新阅读设置
//	@Tags		阅读器
//	@Param		request	body		UpdateSettingsRequest	true	"更新设置请求"
//	@Success	200		{object}	response.Response
//	@Router		/api/v1/reader/settings [put]
func (api *SettingAPI) UpdateReadingSettings(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
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

	err := api.readerService.UpdateReadingSettings(c.Request.Context(), userID.(string), updates)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新阅读设置失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}
