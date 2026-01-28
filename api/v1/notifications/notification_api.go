package notifications

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/notification"
	notifService "Qingyu_backend/service/notification"
)

// NotificationAPI 通知API
type NotificationAPI struct {
	notificationService notifService.NotificationService
}

// NewNotificationAPI 创建通知API实例
func NewNotificationAPI(notificationService notifService.NotificationService) *NotificationAPI {
	return &NotificationAPI{
		notificationService: notificationService,
	}
}

// GetNotifications 获取通知列表
// @Summary 获取通知列表
// @Description 获取当前用户的通知列表，支持分页和筛选
// @Tags notifications
// @Accept json
// @Produce json
// @Param type query string false "通知类型 (system, social, content, reward, message, update, membership)"
// @Param read query bool false "是否已读"
// @Param priority query string false "优先级 (low, normal, high, urgent)"
// @Param keyword query string false "关键词搜索"
// @Param limit query int false "每页数量 (default: 20, max: 100)"
// @Param offset query int false "偏移量 (default: 0)"
// @Param sortBy query string false "排序字段 (created_at, priority, read_at)"
// @Param sortDesc query bool false "是否降序 (default: true)"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications [get]
func (api *NotificationAPI) GetNotifications(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 解析查询参数
	var notifType *notification.NotificationType
	if typeStr := c.Query("type"); typeStr != "" {
		t := notification.NotificationType(typeStr)
		if !t.IsValid() {
			shared.Error(c, http.StatusBadRequest, "INVALID_TYPE", "无效的通知类型")
			return
		}
		notifType = &t
	}

	var read *bool
	if readStr := c.Query("read"); readStr != "" {
		if readBool, err := strconv.ParseBool(readStr); err == nil {
			read = &readBool
		}
	}

	var priority *notification.NotificationPriority
	if priorityStr := c.Query("priority"); priorityStr != "" {
		p := notification.NotificationPriority(priorityStr)
		priority = &p
	}

	var keyword *string
	if keywordStr := c.Query("keyword"); keywordStr != "" {
		keyword = &keywordStr
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortDesc, _ := strconv.ParseBool(c.DefaultQuery("sortDesc", "true"))

	// 构建请求
	req := &notifService.GetNotificationsRequest{
		UserID:   userIDStr,
		Type:     notifType,
		Read:     read,
		Priority: priority,
		Keyword:  keyword,
		Limit:    limit,
		Offset:   offset,
		SortBy:   sortBy,
		SortDesc: sortDesc,
	}

	// 验证请求
	if err := shared.GetValidator().Struct(req); err != nil {
		shared.HandleValidationError(c, err)
		return
	}

	// 调用服务
	result, err := api.notificationService.GetNotifications(c.Request.Context(), req)
	if err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, result)
}

// GetNotification 获取通知详情
// @Summary 获取通知详情
// @Description 根据ID获取通知详情
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "通知ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/{id} [get]
func (api *NotificationAPI) GetNotification(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 获取通知ID
	id := c.Param("id")
	if id == "" {
		shared.Error(c, http.StatusBadRequest, "INVALID_ID", "通知ID不能为空")
		return
	}

	// 调用服务
	notif, err := api.notificationService.GetNotification(c.Request.Context(), id)
	if err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	// 验证权限
	if notif.UserID != userIDStr {
		shared.Error(c, http.StatusForbidden, "FORBIDDEN", "无权访问此通知")
		return
	}

	shared.SuccessData(c, notif)
}

// MarkAsRead 标记通知为已读
// @Summary 标记通知为已读
// @Description 将指定通知标记为已读
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "通知ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/{id}/read [put]
func (api *NotificationAPI) MarkAsRead(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 获取通知ID
	id := c.Param("id")
	if id == "" {
		shared.Error(c, http.StatusBadRequest, "INVALID_ID", "通知ID不能为空")
		return
	}

	// 调用服务
	if err := api.notificationService.MarkAsRead(c.Request.Context(), id, userIDStr); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"message": "标记成功"})
}

// MarkMultipleAsRead 批量标记通知为已读
// @Summary 批量标记通知为已读
// @Description 将多个通知标记为已读
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body object{ids=[]string} true "通知ID列表"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/mark-read [put]
func (api *NotificationAPI) MarkMultipleAsRead(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 解析请求
	var req struct {
		IDs []string `json:"ids" validate:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 验证请求
	if err := shared.GetValidator().Struct(req); err != nil {
		shared.HandleValidationError(c, err)
		return
	}

	// 调用服务
	if err := api.notificationService.MarkMultipleAsRead(c.Request.Context(), req.IDs, userIDStr); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"message": "批量标记成功"})
}

// MarkAllAsRead 标记所有通知为已读
// @Summary 标记所有通知为已读
// @Description 将当前用户的所有未读通知标记为已读
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/read-all [put]
func (api *NotificationAPI) MarkAllAsRead(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 调用服务
	if err := api.notificationService.MarkAllAsRead(c.Request.Context(), userIDStr); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"message": "全部标记成功"})
}

// DeleteNotification 删除通知
// @Summary 删除通知
// @Description 删除指定通知
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "通知ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/{id} [delete]
func (api *NotificationAPI) DeleteNotification(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 获取通知ID
	id := c.Param("id")
	if id == "" {
		shared.Error(c, http.StatusBadRequest, "INVALID_ID", "通知ID不能为空")
		return
	}

	// 调用服务
	if err := api.notificationService.DeleteNotification(c.Request.Context(), id, userIDStr); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"message": "删除成功"})
}

// BatchDeleteNotifications 批量删除通知
// @Summary 批量删除通知
// @Description 批量删除多个通知
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body object{ids=[]string} true "通知ID列表"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/batch-delete [delete]
func (api *NotificationAPI) BatchDeleteNotifications(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 解析请求
	var req struct {
		IDs []string `json:"ids" validate:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 验证请求
	if err := shared.GetValidator().Struct(req); err != nil {
		shared.HandleValidationError(c, err)
		return
	}

	// 调用服务
	if err := api.notificationService.BatchDeleteNotifications(c.Request.Context(), req.IDs, userIDStr); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"message": "批量删除成功"})
}

// DeleteAllNotifications 删除所有通知
// @Summary 删除所有通知
// @Description 删除当前用户的所有通知
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/delete-all [delete]
func (api *NotificationAPI) DeleteAllNotifications(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 调用服务
	if err := api.notificationService.DeleteAllNotifications(c.Request.Context(), userIDStr); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"message": "全部删除成功"})
}

// GetUnreadCount 获取未读通知数量
// @Summary 获取未读通知数量
// @Description 获取当前用户的未读通知数量
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/unread-count [get]
func (api *NotificationAPI) GetUnreadCount(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 调用服务
	count, err := api.notificationService.GetUnreadCount(c.Request.Context(), userIDStr)
	if err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"count": count})
}

// GetNotificationStats 获取通知统计
// @Summary 获取通知统计
// @Description 获取当前用户的通知统计信息
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/stats [get]
func (api *NotificationAPI) GetNotificationStats(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 调用服务
	stats, err := api.notificationService.GetNotificationStats(c.Request.Context(), userIDStr)
	if err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, stats)
}

// GetNotificationPreference 获取通知偏好设置
// @Summary 获取通知偏好设置
// @Description 获取当前用户的通知偏好设置
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/preferences [get]
func (api *NotificationAPI) GetNotificationPreference(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 调用服务
	preference, err := api.notificationService.GetNotificationPreference(c.Request.Context(), userIDStr)
	if err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, preference)
}

// UpdateNotificationPreference 更新通知偏好设置
// @Summary 更新通知偏好设置
// @Description 更新当前用户的通知偏好设置
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body object true "偏好设置"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/preferences [put]
func (api *NotificationAPI) UpdateNotificationPreference(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 解析请求
	var req notifService.UpdateNotificationPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 验证请求
	if err := shared.GetValidator().Struct(req); err != nil {
		shared.HandleValidationError(c, err)
		return
	}

	// 调用服务
	if err := api.notificationService.UpdateNotificationPreference(c.Request.Context(), userIDStr, &req); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"message": "更新成功"})
}

// ResetNotificationPreference 重置通知偏好设置
// @Summary 重置通知偏好设置
// @Description 重置当前用户的通知偏好设置为默认值
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/preferences/reset [post]
func (api *NotificationAPI) ResetNotificationPreference(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 调用服务
	if err := api.notificationService.ResetNotificationPreference(c.Request.Context(), userIDStr); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"message": "重置成功"})
}

// GetEmailNotificationSettings 获取邮件通知设置
// @Summary 获取邮件通知设置
// @Description 获取当前用户的邮件通知设置
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/user-management/email-notifications [get]
func (api *NotificationAPI) GetEmailNotificationSettings(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 调用服务
	settings, err := api.notificationService.GetEmailNotificationSettings(c.Request.Context(), userIDStr)
	if err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, settings)
}

// UpdateEmailNotificationSettings 更新邮件通知设置
// @Summary 更新邮件通知设置
// @Description 更新当前用户的邮件通知设置
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body object true "邮件通知设置"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/user-management/email-notifications [put]
func (api *NotificationAPI) UpdateEmailNotificationSettings(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 解析请求
	var settings notification.EmailNotificationSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		shared.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 调用服务
	if err := api.notificationService.UpdateEmailNotificationSettings(c.Request.Context(), userIDStr, &settings); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"message": "更新成功"})
}

// GetSMSNotificationSettings 获取短信通知设置
// @Summary 获取短信通知设置
// @Description 获取当前用户的短信通知设置
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/user-management/sms-notifications [get]
func (api *NotificationAPI) GetSMSNotificationSettings(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 调用服务
	settings, err := api.notificationService.GetSMSNotificationSettings(c.Request.Context(), userIDStr)
	if err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, settings)
}

// UpdateSMSNotificationSettings 更新短信通知设置
// @Summary 更新短信通知设置
// @Description 更新当前用户的短信通知设置
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body object true "短信通知设置"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/user-management/sms-notifications [put]
func (api *NotificationAPI) UpdateSMSNotificationSettings(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 解析请求
	var settings notification.SMSNotificationSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		shared.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 调用服务
	if err := api.notificationService.UpdateSMSNotificationSettings(c.Request.Context(), userIDStr, &settings); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"message": "更新成功"})
}

// RegisterPushDevice 注册推送设备
// @Summary 注册推送设备
// @Description 注册新的推送设备
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body object true "设备信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/push/register [post]
func (api *NotificationAPI) RegisterPushDevice(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 解析请求
	var req notifService.RegisterPushDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 设置用户ID
	req.UserID = userIDStr

	// 验证请求
	if err := shared.GetValidator().Struct(req); err != nil {
		shared.HandleValidationError(c, err)
		return
	}

	// 调用服务
	device, err := api.notificationService.RegisterPushDevice(c.Request.Context(), &req)
	if err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, device)
}

// UnregisterPushDevice 取消注册推送设备
// @Summary 取消注册推送设备
// @Description 取消注册推送设备
// @Tags notifications
// @Accept json
// @Produce json
// @Param deviceId path string true "设备ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/push/unregister/{deviceId} [delete]
func (api *NotificationAPI) UnregisterPushDevice(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 获取设备ID
	deviceID := c.Param("deviceId")
	if deviceID == "" {
		shared.Error(c, http.StatusBadRequest, "INVALID_DEVICE_ID", "设备ID不能为空")
		return
	}

	// 调用服务
	if err := api.notificationService.UnregisterPushDevice(c.Request.Context(), deviceID, userIDStr); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"message": "取消注册成功"})
}

// GetPushDevices 获取推送设备列表
// @Summary 获取推送设备列表
// @Description 获取当前用户的推送设备列表
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/push/devices [get]
func (api *NotificationAPI) GetPushDevices(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 调用服务
	devices, err := api.notificationService.GetUserPushDevices(c.Request.Context(), userIDStr)
	if err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"devices": devices})
}

// ClearReadNotifications 清除已读通知
// @Summary 清除已读通知
// @Description 清除当前用户的所有已读通知
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/clear-read [post]
func (api *NotificationAPI) ClearReadNotifications(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 调用服务
	count, err := api.notificationService.ClearReadNotifications(c.Request.Context(), userIDStr)
	if err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{
		"message": "清除成功",
		"count":   count,
	})
}

// ResendNotification 重发通知
// @Summary 重发通知
// @Description 重新发送指定的通知给用户
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "通知ID"
// @Param request body object{method=string} true "重发方式 (email, push, sms)"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/{id}/resend [post]
func (api *NotificationAPI) ResendNotification(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "INVALID_USER_ID", "无效的用户ID")
		return
	}

	// 获取通知ID
	id := c.Param("id")
	if id == "" {
		shared.Error(c, http.StatusBadRequest, "INVALID_ID", "通知ID不能为空")
		return
	}

	// 解析请求
	var req struct {
		Method string `json:"method" validate:"required,oneof=email push sms"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 验证请求
	if err := shared.GetValidator().Struct(&req); err != nil {
		shared.HandleValidationError(c, err)
		return
	}

	// 调用服务
	if err := api.notificationService.ResendNotification(c.Request.Context(), id, userIDStr, req.Method); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.SuccessData(c, gin.H{"message": "重发成功"})
}

// GetWSEndpoint 获取WebSocket端点
// @Summary 获取WebSocket端点
// @Description 获取通知WebSocket连接的端点信息
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/ws-endpoint [get]
func (api *NotificationAPI) GetWSEndpoint(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		shared.Error(c, http.StatusUnauthorized, "MISSING_TOKEN", "缺少认证令牌")
		return
	}

	// 提取Bearer token
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		shared.Error(c, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "无效的令牌格式")
		return
	}

	// 确定协议（ws或wss）
	scheme := "ws"
	if c.Request.TLS != nil {
		scheme = "wss"
	}

	// 获取主机
	host := c.Request.Host
	if host == "" {
		host = "localhost:8080"
	}

	// 构建WebSocket URL
	wsURL := fmt.Sprintf("%s://%s/api/v1/notifications/ws?token=%s", scheme, host, token)

	shared.SuccessData(c, gin.H{
		"url": wsURL,
		"message": "WebSocket端点获取成功",
	})
}
