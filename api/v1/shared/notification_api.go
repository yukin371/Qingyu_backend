package shared

import (
	messagingModel "Qingyu_backend/models/messaging"
	"Qingyu_backend/service/shared/messaging"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NotificationAPI 通知API处理器
type NotificationAPI struct {
	notificationService messaging.NotificationService
}

// NewNotificationAPI 创建通知API实例
func NewNotificationAPI(notificationService messaging.NotificationService) *NotificationAPI {
	return &NotificationAPI{
		notificationService: notificationService,
	}
}

// ============ 通知API ============

// GetNotifications 获取用户通知列表
//
//	@Summary		获取通知列表
//	@Description	获取当前用户的通知列表（分页）
//	@Tags			通知
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"页码（默认1）"
//	@Param			page_size	query		int		false	"每页大小（默认20）"
//	@Param			is_read		query		bool	false	"是否已读过滤"
//	@Success		200			{object}	PaginatedResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/notifications [get]
func (api *NotificationAPI) GetNotifications(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("userId")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}
	userID := userIDInterface.(string)

	// 2. 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 3. 获取通知列表
	notifications, total, err := api.notificationService.ListNotifications(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		Error(c, 500, "获取通知列表失败", err.Error())
		return
	}

	// 4. 返回结果
	Paginated(c, notifications, total, page, pageSize, "获取成功")
}

// GetUnreadCount 获取未读通知数量
//
//	@Summary		获取未读数量
//	@Description	获取当前用户的未读通知数量
//	@Tags			通知
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	APIResponse{data=map[string]int64}
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/notifications/unread-count [get]
func (api *NotificationAPI) GetUnreadCount(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("userId")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}
	userID := userIDInterface.(string)

	// 2. 获取未读数量
	count, err := api.notificationService.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		Error(c, 500, "获取未读数量失败", err.Error())
		return
	}

	// 3. 返回结果
	Success(c, 200, "获取成功", map[string]int64{
		"unread_count": count,
	})
}

// MarkAsRead 标记通知为已读
//
//	@Summary		标记为已读
//	@Description	标记指定通知为已读
//	@Tags			通知
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"通知ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/notifications/{id}/read [put]
func (api *NotificationAPI) MarkAsRead(c *gin.Context) {
	// 1. 获取通知ID
	notificationID := c.Param("id")
	if notificationID == "" {
		BadRequest(c, "参数错误", "通知ID不能为空")
		return
	}

	// 2. 标记为已读
	err := api.notificationService.MarkAsRead(c.Request.Context(), notificationID)
	if err != nil {
		Error(c, 500, "标记失败", err.Error())
		return
	}

	// 3. 返回成功
	Success(c, 200, "标记成功", nil)
}

// MarkAllAsRead 标记所有通知为已读
//
//	@Summary		标记所有为已读
//	@Description	标记当前用户的所有通知为已读
//	@Tags			通知
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	APIResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/notifications/read-all [put]
func (api *NotificationAPI) MarkAllAsRead(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("userId")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}
	userID := userIDInterface.(string)

	// 2. 标记所有为已读
	err := api.notificationService.MarkAllAsRead(c.Request.Context(), userID)
	if err != nil {
		Error(c, 500, "标记失败", err.Error())
		return
	}

	// 3. 返回成功
	Success(c, 200, "标记成功", nil)
}

// DeleteNotification 删除通知
//
//	@Summary		删除通知
//	@Description	删除指定通知
//	@Tags			通知
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"通知ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/notifications/{id} [delete]
func (api *NotificationAPI) DeleteNotification(c *gin.Context) {
	// 1. 获取通知ID
	notificationID := c.Param("id")
	if notificationID == "" {
		BadRequest(c, "参数错误", "通知ID不能为空")
		return
	}

	// 2. 删除通知
	err := api.notificationService.DeleteNotification(c.Request.Context(), notificationID)
	if err != nil {
		Error(c, 500, "删除失败", err.Error())
		return
	}

	// 3. 返回成功
	Success(c, 200, "删除成功", nil)
}

// CreateNotification 创建系统通知（管理员API）
//
//	@Summary		创建系统通知
//	@Description	创建系统通知（仅管理员）
//	@Tags			通知
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateNotificationRequest	true	"创建请求"
//	@Success		200		{object}	APIResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/v1/notifications [post]
func (api *NotificationAPI) CreateNotification(c *gin.Context) {
	// 1. 解析请求
	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	// 2. 创建通知
	notification := &messagingModel.NotificationDelivery{
		UserID:  req.UserID,
		Type:    req.Type,
		Title:   req.Title,
		Content: req.Content,
	}

	err := api.notificationService.CreateNotification(c.Request.Context(), notification)
	if err != nil {
		Error(c, 500, "创建通知失败", err.Error())
		return
	}

	// 3. 返回成功
	Success(c, 200, "创建成功", notification)
}

// ============ 请求/响应结构 ============

// CreateNotificationRequest 创建通知请求
type CreateNotificationRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content" binding:"required"`
}

// TODO(Phase3): 推送通知API
// SendPushNotification 发送推送通知
// func (api *NotificationAPI) SendPushNotification(c *gin.Context) {
//     // 推送通知功能
// }

// TODO(Phase3): 短信通知API
// SendSMSNotification 发送短信通知
// func (api *NotificationAPI) SendSMSNotification(c *gin.Context) {
//     // 短信通知功能
// }
