package channels

import (
	"context"
	"fmt"
	"time"

	messagingModel "Qingyu_backend/models/messaging"
	messagingRepo "Qingyu_backend/repository/mongodb/messaging"
)

// InboxNotificationService 站内通知服务接口
type InboxNotificationService interface {
	// SendNotification 发送站内通知
	SendNotification(ctx context.Context, notification *messagingModel.InboxNotification) error

	// SendBatchNotifications 批量发送站内通知
	SendBatchNotifications(ctx context.Context, notifications []*messagingModel.InboxNotification) error

	// GetUserNotifications 获取用户通知列表
	GetUserNotifications(ctx context.Context, userID string, filter *InboxNotificationFilter) ([]*messagingModel.InboxNotification, int64, error)

	// GetNotificationByID 根据ID获取通知
	GetNotificationByID(ctx context.Context, id string) (*messagingModel.InboxNotification, error)

	// GetUnreadCount 获取未读通知数量
	GetUnreadCount(ctx context.Context, userID string) (int64, error)

	// MarkAsRead 标记为已读
	MarkAsRead(ctx context.Context, id, userID string) error

	// MarkAllAsRead 标记所有为已读
	MarkAllAsRead(ctx context.Context, userID string) (int64, error)

	// DeleteNotification 删除通知
	DeleteNotification(ctx context.Context, id, userID string) error

	// PinNotification 置顶通知
	PinNotification(ctx context.Context, id, userID string) error

	// UnpinNotification 取消置顶
	UnpinNotification(ctx context.Context, id, userID string) error

	// CleanupExpired 清理过期通知
	CleanupExpired(ctx context.Context) (int64, error)
}

// InboxNotificationServiceImpl 站内通知服务实现
type InboxNotificationServiceImpl struct {
	repo *messagingRepo.InboxNotificationRepository
}

// NewInboxNotificationService 创建站内通知服务
func NewInboxNotificationService(repo *messagingRepo.InboxNotificationRepository) InboxNotificationService {
	return &InboxNotificationServiceImpl{
		repo: repo,
	}
}

// InboxNotificationFilter 服务层过滤器
type InboxNotificationFilter struct {
	Type      messagingModel.InboxNotificationType
	IsRead    *bool
	IsPinned  *bool
	Priority  messagingModel.InboxNotificationPriority
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	PageSize  int
}

// SendNotification 发送站内通知
func (s *InboxNotificationServiceImpl) SendNotification(ctx context.Context, notification *messagingModel.InboxNotification) error {
	// 参数验证
	if err := s.validateNotification(notification); err != nil {
		return fmt.Errorf("通知验证失败: %w", err)
	}

	// 设置默认值
	s.setDefaults(notification)

	// 保存到数据库
	_, err := s.repo.Create(ctx, notification)
	if err != nil {
		return fmt.Errorf("发送站内通知失败: %w", err)
	}

	return nil
}

// SendBatchNotifications 批量发送站内通知
func (s *InboxNotificationServiceImpl) SendBatchNotifications(ctx context.Context, notifications []*messagingModel.InboxNotification) error {
	if len(notifications) == 0 {
		return fmt.Errorf("通知列表不能为空")
	}

	// 验证并发送每个通知
	for _, notification := range notifications {
		if err := s.SendNotification(ctx, notification); err != nil {
			// 记录错误但继续处理其他通知
			// TODO: 添加日志记录
			continue
		}
	}

	return nil
}

// GetUserNotifications 获取用户通知列表
func (s *InboxNotificationServiceImpl) GetUserNotifications(ctx context.Context, userID string, filter *InboxNotificationFilter) ([]*messagingModel.InboxNotification, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}

	// 构建Repository过滤器
	repoFilter := &messagingRepo.InboxNotificationFilter{}

	if filter != nil {
		repoFilter.Type = filter.Type
		repoFilter.IsRead = filter.IsRead
		repoFilter.IsPinned = filter.IsPinned
		repoFilter.Priority = filter.Priority
		repoFilter.StartDate = filter.StartDate
		repoFilter.EndDate = filter.EndDate
		repoFilter.Page = filter.Page
		repoFilter.PageSize = filter.PageSize
	}

	// 查询通知列表
	notifications, total, err := s.repo.GetByReceiverID(ctx, userID, repoFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户通知列表失败: %w", err)
	}

	return notifications, total, nil
}

// GetNotificationByID 根据ID获取通知
func (s *InboxNotificationServiceImpl) GetNotificationByID(ctx context.Context, id string) (*messagingModel.InboxNotification, error) {
	if id == "" {
		return nil, fmt.Errorf("通知ID不能为空")
	}

	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取通知失败: %w", err)
	}

	if notification == nil {
		return nil, fmt.Errorf("通知不存在")
	}

	return notification, nil
}

// GetUnreadCount 获取未读通知数量
func (s *InboxNotificationServiceImpl) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, fmt.Errorf("用户ID不能为空")
	}

	count, err := s.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("获取未读数量失败: %w", err)
	}

	return count, nil
}

// MarkAsRead 标记为已读
func (s *InboxNotificationServiceImpl) MarkAsRead(ctx context.Context, id, userID string) error {
	if id == "" {
		return fmt.Errorf("通知ID不能为空")
	}
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	err := s.repo.MarkAsRead(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("标记已读失败: %w", err)
	}

	return nil
}

// MarkAllAsRead 标记所有为已读
func (s *InboxNotificationServiceImpl) MarkAllAsRead(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, fmt.Errorf("用户ID不能为空")
	}

	count, err := s.repo.MarkAllAsRead(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("标记全部已读失败: %w", err)
	}

	return count, nil
}

// DeleteNotification 删除通知
func (s *InboxNotificationServiceImpl) DeleteNotification(ctx context.Context, id, userID string) error {
	if id == "" {
		return fmt.Errorf("通知ID不能为空")
	}
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	err := s.repo.SoftDelete(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("删除通知失败: %w", err)
	}

	return nil
}

// PinNotification 置顶通知
func (s *InboxNotificationServiceImpl) PinNotification(ctx context.Context, id, userID string) error {
	if id == "" {
		return fmt.Errorf("通知ID不能为空")
	}
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	err := s.repo.PinNotification(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("置顶通知失败: %w", err)
	}

	return nil
}

// UnpinNotification 取消置顶
func (s *InboxNotificationServiceImpl) UnpinNotification(ctx context.Context, id, userID string) error {
	if id == "" {
		return fmt.Errorf("通知ID不能为空")
	}
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	err := s.repo.UnpinNotification(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("取消置顶失败: %w", err)
	}

	return nil
}

// CleanupExpired 清理过期通知
func (s *InboxNotificationServiceImpl) CleanupExpired(ctx context.Context) (int64, error) {
	count, err := s.repo.DeleteExpired(ctx)
	if err != nil {
		return 0, fmt.Errorf("清理过期通知失败: %w", err)
	}

	return count, nil
}

// validateNotification 验证通知数据
func (s *InboxNotificationServiceImpl) validateNotification(notification *messagingModel.InboxNotification) error {
	if notification == nil {
		return fmt.Errorf("通知不能为空")
	}

	if notification.ReceiverID == "" {
		return fmt.Errorf("接收者ID不能为空")
	}

	if notification.Type == "" {
		return fmt.Errorf("通知类型不能为空")
	}

	if notification.Title == "" {
		return fmt.Errorf("通知标题不能为空")
	}

	if notification.Content == "" {
		return fmt.Errorf("通知内容不能为空")
	}

	return nil
}

// setDefaults 设置默认值
func (s *InboxNotificationServiceImpl) setDefaults(notification *messagingModel.InboxNotification) {
	// 默认优先级
	if notification.Priority == "" {
		notification.Priority = messagingModel.InboxNotificationPriorityNormal
	}

	// 默认未读
	if notification.IsRead == false && notification.ReadAt == nil {
		notification.IsRead = false
	}
}
