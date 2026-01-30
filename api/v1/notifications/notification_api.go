package notifications

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/notifications/dto"
	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/notification"
	notifService "Qingyu_backend/service/notification"
	"Qingyu_backend/pkg/response"
	"errors"
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 解析查询参数
	var notifType *notification.NotificationType
	if typeStr := c.Query("type"); typeStr != "" {
		t := notification.NotificationType(typeStr)
		if !t.IsValid() {
			response.BadRequest(c,  "INVALID_TYPE", "无效的通知类型")
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

	// 添加输入验证，防止DoS攻击
	if limit > 100 {
		limit = 100
	}
	if limit < 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

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
		response.InternalError(c, err)
		return
	}

	response.Success(c, result)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 获取通知ID
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c,  "INVALID_ID", "通知ID不能为空")
		return
	}

	// 调用服务
	notif, err := api.notificationService.GetNotification(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 验证权限
	if notif.UserID != userIDStr {
		response.Forbidden(c, "无权访问此通知")
		return
	}

	response.Success(c, notif)
}

// MarkAsRead 标记通知为已读
// @Summary 标记通知为已读
// @Description 将指定通知标记为已读
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "通知ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/{id}/read [post]
func (api *NotificationAPI) MarkAsRead(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 获取通知ID
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c,  "INVALID_ID", "通知ID不能为空")
		return
	}

	// 调用服务
	if err := api.notificationService.MarkAsRead(c.Request.Context(), id, userIDStr); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "标记成功", nil)
}

// MarkMultipleAsRead 批量标记通知为已读
// @Summary 批量标记通知为已读
// @Description 将多个通知标记为已读
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body dto.BatchMarkReadRequest true "批量标记已读请求"
// @Success 200 {object} shared.APIResponse{data=dto.BatchOperationResponse}
// @Router /api/v1/notifications/batch-read [post]
func (api *NotificationAPI) MarkMultipleAsRead(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 解析请求
	var req dto.BatchMarkReadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 验证请求
	if err := shared.GetValidator().Struct(req); err != nil {
		shared.HandleValidationError(c, err)
		return
	}

	// 调用服务
	succeeded, failed, err := api.notificationService.MarkMultipleAsReadWithResult(c.Request.Context(), req.NotificationIDs, userIDStr)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 构建响应
	resp := dto.BatchOperationResponse{
		Success:   failed == 0,
		Total:     len(req.NotificationIDs),
		Succeeded: succeeded,
		Failed:    failed,
	}

	if failed > 0 {
		resp.Errors = []string{"部分通知标记失败"}
	}

	response.Success(c, resp)
}

// MarkAllAsRead 标记所有通知为已读
// @Summary 标记所有通知为已读
// @Description 将当前用户的所有未读通知标记为已读
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/notifications/read-all [post]
func (api *NotificationAPI) MarkAllAsRead(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 调用服务
	if err := api.notificationService.MarkAllAsRead(c.Request.Context(), userIDStr); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "全部标记成功", nil)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 获取通知ID
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c,  "INVALID_ID", "通知ID不能为空")
		return
	}

	// 调用服务
	if err := api.notificationService.DeleteNotification(c.Request.Context(), id, userIDStr); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// BatchDeleteNotifications 批量删除通知
// @Summary 批量删除通知
// @Description 批量删除多个通知
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body dto.BatchDeleteRequest true "批量删除请求"
// @Success 200 {object} shared.APIResponse{data=dto.BatchOperationResponse}
// @Router /api/v1/notifications/batch-delete [post]
func (api *NotificationAPI) BatchDeleteNotifications(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 解析请求
	var req dto.BatchDeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 验证请求
	if err := shared.GetValidator().Struct(req); err != nil {
		shared.HandleValidationError(c, err)
		return
	}

	// 调用服务
	succeeded, failed, err := api.notificationService.BatchDeleteNotificationsWithResult(c.Request.Context(), req.NotificationIDs, userIDStr)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 构建响应
	resp := dto.BatchOperationResponse{
		Success:   failed == 0,
		Total:     len(req.NotificationIDs),
		Succeeded: succeeded,
		Failed:    failed,
	}

	if failed > 0 {
		resp.Errors = []string{"部分通知删除失败"}
	}

	response.Success(c, resp)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 调用服务
	if err := api.notificationService.DeleteAllNotifications(c.Request.Context(), userIDStr); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "全部删除成功", nil)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 调用服务
	count, err := api.notificationService.GetUnreadCount(c.Request.Context(), userIDStr)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{"count": count})
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 调用服务
	stats, err := api.notificationService.GetNotificationStats(c.Request.Context(), userIDStr)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, stats)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 调用服务
	preference, err := api.notificationService.GetNotificationPreference(c.Request.Context(), userIDStr)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, preference)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 解析请求
	var req notifService.UpdateNotificationPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 验证请求
	if err := shared.GetValidator().Struct(req); err != nil {
		shared.HandleValidationError(c, err)
		return
	}

	// 调用服务
	if err := api.notificationService.UpdateNotificationPreference(c.Request.Context(), userIDStr, &req); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 调用服务
	if err := api.notificationService.ResetNotificationPreference(c.Request.Context(), userIDStr); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "重置成功", nil)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 调用服务
	settings, err := api.notificationService.GetEmailNotificationSettings(c.Request.Context(), userIDStr)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, settings)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 解析请求
	var settings notification.EmailNotificationSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		response.BadRequest(c,  "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 调用服务
	if err := api.notificationService.UpdateEmailNotificationSettings(c.Request.Context(), userIDStr, &settings); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 调用服务
	settings, err := api.notificationService.GetSMSNotificationSettings(c.Request.Context(), userIDStr)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, settings)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 解析请求
	var settings notification.SMSNotificationSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		response.BadRequest(c,  "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 调用服务
	if err := api.notificationService.UpdateSMSNotificationSettings(c.Request.Context(), userIDStr, &settings); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 解析请求
	var req notifService.RegisterPushDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "INVALID_REQUEST", "请求参数错误")
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
		response.InternalError(c, err)
		return
	}

	response.Success(c, device)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 获取设备ID
	deviceID := c.Param("deviceId")
	if deviceID == "" {
		response.BadRequest(c,  "INVALID_DEVICE_ID", "设备ID不能为空")
		return
	}

	// 调用服务
	if err := api.notificationService.UnregisterPushDevice(c.Request.Context(), deviceID, userIDStr); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "取消注册成功", nil)
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
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 调用服务
	devices, err := api.notificationService.GetUserPushDevices(c.Request.Context(), userIDStr)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{"devices": devices})
}

// ClearReadNotifications 清除已读通知
// @Summary 清除已读通知
// @Description 删除用户所有已读的通知
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} shared.APIResponse{data=dto.BatchOperationResponse}
// @Router /api/v1/notifications/clear-read [post]
func (api *NotificationAPI) ClearReadNotifications(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 调用服务
	affected, err := api.notificationService.ClearReadNotifications(c.Request.Context(), userIDStr)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, dto.BatchOperationResponse{
		Success:   true,
		Total:     int(affected),
		Succeeded: int(affected),
		Failed:    0,
	})
}

// ResendNotification 重新发送通知
// @Summary 重新发送通知
// @Description 重新发送通知（邮件/推送）
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "通知ID"
// @Param request body object{method=string} true "重新发送请求"
// @Success 200 {object} shared.APIResponse{data=dto.MarkAsReadResponse}
// @Failure 400 {object} shared.APIResponse "参数错误"
// @Failure 404 {object} shared.APIResponse "通知不存在"
// @Router /api/v1/notifications/{id}/resend [post]
func (api *NotificationAPI) ResendNotification(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 获取通知ID
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c,  "INVALID_ID", "通知ID不能为空")
		return
	}

	// 解析请求
	var req struct {
		Method string `json:"method" validate:"required,oneof=email push sms"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "INVALID_REQUEST", "请求参数错误")
		return
	}

	// 验证请求
	if err := shared.GetValidator().Struct(&req); err != nil {
		shared.HandleValidationError(c, err)
		return
	}

	// 调用服务
	err := api.notificationService.ResendNotification(c.Request.Context(), id, userIDStr, req.Method)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, dto.MarkAsReadResponse{
		Success: true,
		Message: "重新发送成功",
	})
}

// GetWSEndpoint WebSocket端点
// @Summary WebSocket端点
// @Description 获取WebSocket连接地址，用于实时推送通知
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} shared.APIResponse{data=dto.WSEndpointResponse}
// @Router /api/v1/notifications/ws-endpoint [get]
func (api *NotificationAPI) GetWSEndpoint(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	_, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("INVALID_USER_ID: 无效的用户ID"))
		return
	}

	// 验证 Authorization Header 存在
	if c.GetHeader("Authorization") == "" {
		response.Unauthorized(c, "缺少认证令牌")
		return
	}

	// 生成WebSocket连接URL（token 通过 Authorization Header 传递）
	scheme := "ws"
	if c.Request.TLS != nil {
		scheme = "wss"
	}

	wsURL := fmt.Sprintf("%s://%s/ws/notifications", scheme, c.Request.Host)

	response.Success(c, dto.WSEndpointResponse{
		URL:     wsURL,
		Message: "请使用Authorization Header或子协议传递token",
	})
}
