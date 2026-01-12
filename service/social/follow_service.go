package social

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/social"
	socialRepo "Qingyu_backend/repository/interfaces/social"
	"Qingyu_backend/service/base"
)

// FollowService 关注服务
type FollowService struct {
	followRepo  socialRepo.FollowRepository
	eventBus    base.EventBus
	serviceName string
	version     string
}

// NewFollowService 创建关注服务实例
func NewFollowService(
	followRepo socialRepo.FollowRepository,
	eventBus base.EventBus,
) *FollowService {
	return &FollowService{
		followRepo:  followRepo,
		eventBus:    eventBus,
		serviceName: "FollowService",
		version:     "1.0.0",
	}
}

// =========================
// BaseService 接口实现
// =========================

// Initialize 初始化服务
func (s *FollowService) Initialize(ctx context.Context) error {
	return nil
}

// Health 健康检查
func (s *FollowService) Health(ctx context.Context) error {
	if err := s.followRepo.Health(ctx); err != nil {
		return fmt.Errorf("关注Repository健康检查失败: %w", err)
	}
	return nil
}

// Close 关闭服务
func (s *FollowService) Close(ctx context.Context) error {
	return nil
}

// GetServiceName 获取服务名称
func (s *FollowService) GetServiceName() string {
	return s.serviceName
}

// GetVersion 获取服务版本
func (s *FollowService) GetVersion() string {
	return s.version
}

// =========================
// 用户关注
// =========================

// FollowUser 关注用户
func (s *FollowService) FollowUser(ctx context.Context, followerID, followingID string) error {
	// 参数验证
	if followerID == "" || followingID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if followerID == followingID {
		return fmt.Errorf("不能关注自己")
	}

	// 检查是否已关注
	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID, "user")
	if err != nil {
		return fmt.Errorf("检查关注状态失败: %w", err)
	}
	if isFollowing {
		return fmt.Errorf("已经关注过该用户")
	}

	// 检查是否被对方关注（判断是否互相关注）
	isFollowed, err := s.followRepo.IsFollowing(ctx, followingID, followerID, "user")
	if err != nil {
		return fmt.Errorf("检查互相关注状态失败: %w", err)
	}

	// 创建关注关系
	follow := &social.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
		FollowType:  "user",
		IsMutual:    isFollowed,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.followRepo.CreateFollow(ctx, follow); err != nil {
		return fmt.Errorf("关注失败: %w", err)
	}

	// 如果是互相关注，更新对方的关注状态
	if isFollowed {
		if err := s.followRepo.UpdateMutualStatus(ctx, followingID, followerID, "user", true); err != nil {
			fmt.Printf("Warning: Failed to update mutual status: %v\n", err)
		}
	}

	// 更新统计
	if err := s.followRepo.UpdateFollowStats(ctx, followerID, 0, 1); err != nil {
		fmt.Printf("Warning: Failed to update follower stats: %v\n", err)
	}
	if err := s.followRepo.UpdateFollowStats(ctx, followingID, 1, 0); err != nil {
		fmt.Printf("Warning: Failed to update following stats: %v\n", err)
	}

	// 发布事件
	s.publishFollowEvent(ctx, "follow.created", followerID, followingID, "user")

	return nil
}

// UnfollowUser 取消关注用户
func (s *FollowService) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	// 参数验证
	if followerID == "" || followingID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	// 检查是否已关注
	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID, "user")
	if err != nil {
		return fmt.Errorf("检查关注状态失败: %w", err)
	}
	if !isFollowing {
		return fmt.Errorf("未关注该用户")
	}

	// 检查是否互相关注
	isMutual, err := s.followRepo.IsFollowing(ctx, followingID, followerID, "user")
	if err != nil {
		return fmt.Errorf("检查互相关注状态失败: %w", err)
	}

	// 删除关注关系
	if err := s.followRepo.DeleteFollow(ctx, followerID, followingID, "user"); err != nil {
		return fmt.Errorf("取消关注失败: %w", err)
	}

	// 如果之前是互相关注，更新对方的关注状态
	if isMutual {
		if err := s.followRepo.UpdateMutualStatus(ctx, followingID, followerID, "user", false); err != nil {
			fmt.Printf("Warning: Failed to update mutual status: %v\n", err)
		}
	}

	// 更新统计
	if err := s.followRepo.UpdateFollowStats(ctx, followerID, 0, -1); err != nil {
		fmt.Printf("Warning: Failed to update follower stats: %v\n", err)
	}
	if err := s.followRepo.UpdateFollowStats(ctx, followingID, -1, 0); err != nil {
		fmt.Printf("Warning: Failed to update following stats: %v\n", err)
	}

	// 发布事件
	s.publishFollowEvent(ctx, "follow.deleted", followerID, followingID, "user")

	return nil
}

// CheckFollowStatus 检查关注状态
func (s *FollowService) CheckFollowStatus(ctx context.Context, followerID, followingID string) (bool, error) {
	if followerID == "" || followingID == "" {
		return false, fmt.Errorf("用户ID不能为空")
	}

	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID, "user")
	if err != nil {
		return false, fmt.Errorf("检查关注状态失败: %w", err)
	}

	return isFollowing, nil
}

// GetFollowers 获取粉丝列表
func (s *FollowService) GetFollowers(ctx context.Context, userID string, page, size int) ([]*social.FollowInfo, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}

	// 参数验证和默认值
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	followers, total, err := s.followRepo.GetFollowers(ctx, userID, "user", page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取粉丝列表失败: %w", err)
	}

	return followers, total, nil
}

// GetFollowing 获取关注列表
func (s *FollowService) GetFollowing(ctx context.Context, userID string, page, size int) ([]*social.FollowingInfo, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}

	// 参数验证和默认值
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	following, total, err := s.followRepo.GetFollowing(ctx, userID, "user", page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取关注列表失败: %w", err)
	}

	return following, total, nil
}

// GetFollowStats 获取关注统计
func (s *FollowService) GetFollowStats(ctx context.Context, userID string) (*social.FollowStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	stats, err := s.followRepo.GetFollowStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取关注统计失败: %w", err)
	}

	return stats, nil
}

// =========================
// 作者关注
// =========================

// FollowAuthor 关注作者
func (s *FollowService) FollowAuthor(ctx context.Context, userID, authorID, authorName, authorAvatar string, notifyNewBook bool) error {
	// 参数验证
	if userID == "" || authorID == "" {
		return fmt.Errorf("用户ID和作者ID不能为空")
	}

	// 检查是否已关注
	existing, err := s.followRepo.GetAuthorFollow(ctx, userID, authorID)
	if err != nil {
		return fmt.Errorf("检查关注状态失败: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("已经关注过该作者")
	}

	// 创建作者关注
	authorFollow := &social.AuthorFollow{
		UserID:        userID,
		AuthorID:      authorID,
		AuthorName:    authorName,
		AuthorAvatar:  authorAvatar,
		NotifyNewBook: notifyNewBook,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.followRepo.CreateAuthorFollow(ctx, authorFollow); err != nil {
		return fmt.Errorf("关注作者失败: %w", err)
	}

	// 更新统计
	if err := s.followRepo.UpdateFollowStats(ctx, userID, 0, 1); err != nil {
		fmt.Printf("Warning: Failed to update follow stats: %v\n", err)
	}

	// 发布事件
	s.publishFollowEvent(ctx, "author_follow.created", userID, authorID, "author")

	return nil
}

// UnfollowAuthor 取消关注作者
func (s *FollowService) UnfollowAuthor(ctx context.Context, userID, authorID string) error {
	// 参数验证
	if userID == "" || authorID == "" {
		return fmt.Errorf("用户ID和作者ID不能为空")
	}

	// 检查是否已关注
	existing, err := s.followRepo.GetAuthorFollow(ctx, userID, authorID)
	if err != nil {
		return fmt.Errorf("检查关注状态失败: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("未关注该作者")
	}

	// 删除关注
	if err := s.followRepo.DeleteAuthorFollow(ctx, userID, authorID); err != nil {
		return fmt.Errorf("取消关注失败: %w", err)
	}

	// 更新统计
	if err := s.followRepo.UpdateFollowStats(ctx, userID, 0, -1); err != nil {
		fmt.Printf("Warning: Failed to update follow stats: %v\n", err)
	}

	// 发布事件
	s.publishFollowEvent(ctx, "author_follow.deleted", userID, authorID, "author")

	return nil
}

// GetFollowingAuthors 获取关注的作者列表
func (s *FollowService) GetFollowingAuthors(ctx context.Context, userID string, page, size int) ([]*social.AuthorFollow, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}

	// 参数验证和默认值
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	authors, total, err := s.followRepo.GetUserFollowingAuthors(ctx, userID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取关注作者列表失败: %w", err)
	}

	return authors, total, nil
}

// =========================
// 私有辅助方法
// =========================

// publishFollowEvent 发布关注事件
func (s *FollowService) publishFollowEvent(ctx context.Context, eventType string, followerID, targetID, followType string) {
	if s.eventBus == nil {
		return
	}

	event := &base.BaseEvent{
		EventType: eventType,
		EventData: map[string]interface{}{
			"follower_id": followerID,
			"target_id":   targetID,
			"follow_type": followType,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}
