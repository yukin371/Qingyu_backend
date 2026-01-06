package reader

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// VIPPermissionService VIP权限验证服务接口
type VIPPermissionService interface {
	// CheckVIPAccess 检查用户是否有VIP权限访问章节
	CheckVIPAccess(ctx context.Context, userID, chapterID string, isVIPChapter bool) (bool, error)

	// CheckUserVIPStatus 检查用户是否为VIP用户
	CheckUserVIPStatus(ctx context.Context, userID string) (bool, error)

	// CheckChapterPurchased 检查用户是否已购买该章节
	CheckChapterPurchased(ctx context.Context, userID, chapterID string) (bool, error)

	// GrantVIPAccess 授予用户VIP权限（用于测试或管理）
	GrantVIPAccess(ctx context.Context, userID string, duration time.Duration) error

	// GrantChapterAccess 授予用户章节访问权限
	GrantChapterAccess(ctx context.Context, userID, chapterID string) error
}

// VIPPermissionServiceImpl VIP权限服务实现
type VIPPermissionServiceImpl struct {
	redisClient *redis.Client
	prefix      string
}

// NewVIPPermissionService 创建VIP权限服务
func NewVIPPermissionService(redisClient *redis.Client, prefix string) VIPPermissionService {
	if prefix == "" {
		prefix = "qingyu"
	}
	return &VIPPermissionServiceImpl{
		redisClient: redisClient,
		prefix:      prefix,
	}
}

// 缓存键生成
func (s *VIPPermissionServiceImpl) getKey(key string) string {
	return fmt.Sprintf("%s:vip:%s", s.prefix, key)
}

// CheckVIPAccess 检查用户是否有VIP权限访问章节
func (s *VIPPermissionServiceImpl) CheckVIPAccess(ctx context.Context, userID, chapterID string, isVIPChapter bool) (bool, error) {
	// 1. 如果不是VIP章节，任何人都可以访问
	if !isVIPChapter {
		return true, nil
	}

	// 2. 检查用户是否为VIP用户
	isVIP, err := s.CheckUserVIPStatus(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("检查用户VIP状态失败: %w", err)
	}
	if isVIP {
		return true, nil
	}

	// 3. 检查用户是否已购买该章节
	purchased, err := s.CheckChapterPurchased(ctx, userID, chapterID)
	if err != nil {
		return false, fmt.Errorf("检查章节购买状态失败: %w", err)
	}
	if purchased {
		return true, nil
	}

	// 4. 无权限访问
	return false, nil
}

// CheckUserVIPStatus 检查用户是否为VIP用户
func (s *VIPPermissionServiceImpl) CheckUserVIPStatus(ctx context.Context, userID string) (bool, error) {
	key := s.getKey(fmt.Sprintf("user:%s:status", userID))

	// 从Redis获取VIP状态
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("检查VIP状态失败: %w", err)
	}

	return exists > 0, nil
}

// CheckChapterPurchased 检查用户是否已购买该章节
func (s *VIPPermissionServiceImpl) CheckChapterPurchased(ctx context.Context, userID, chapterID string) (bool, error) {
	key := s.getKey(fmt.Sprintf("purchase:%s:chapters", userID))

	// 使用Redis Set检查章节是否在用户购买列表中
	isMember, err := s.redisClient.SIsMember(ctx, key, chapterID).Result()
	if err != nil {
		return false, fmt.Errorf("检查章节购买状态失败: %w", err)
	}

	return isMember, nil
}

// GrantVIPAccess 授予用户VIP权限
func (s *VIPPermissionServiceImpl) GrantVIPAccess(ctx context.Context, userID string, duration time.Duration) error {
	key := s.getKey(fmt.Sprintf("user:%s:status", userID))

	// 设置VIP状态，带过期时间
	err := s.redisClient.Set(ctx, key, "vip", duration).Err()
	if err != nil {
		return fmt.Errorf("授予VIP权限失败: %w", err)
	}

	return nil
}

// GrantChapterAccess 授予用户章节访问权限
func (s *VIPPermissionServiceImpl) GrantChapterAccess(ctx context.Context, userID, chapterID string) error {
	key := s.getKey(fmt.Sprintf("purchase:%s:chapters", userID))

	// 将章节ID添加到用户购买集合中
	err := s.redisClient.SAdd(ctx, key, chapterID).Err()
	if err != nil {
		return fmt.Errorf("授予章节访问权限失败: %w", err)
	}

	// 设置过期时间（永久购买可以设置很长的时间，如1年）
	_ = s.redisClient.Expire(ctx, key, 365*24*time.Hour)

	return nil
}

// =========================
// 辅助方法
// =========================

// RevokeVIPAccess 撤销用户VIP权限
func (s *VIPPermissionServiceImpl) RevokeVIPAccess(ctx context.Context, userID string) error {
	key := s.getKey(fmt.Sprintf("user:%s:status", userID))
	return s.redisClient.Del(ctx, key).Err()
}

// RevokeChapterAccess 撤销用户章节访问权限
func (s *VIPPermissionServiceImpl) RevokeChapterAccess(ctx context.Context, userID, chapterID string) error {
	key := s.getKey(fmt.Sprintf("purchase:%s:chapters", userID))
	return s.redisClient.SRem(ctx, key, chapterID).Err()
}

// GetUserPurchasedChapters 获取用户购买的所有章节
func (s *VIPPermissionServiceImpl) GetUserPurchasedChapters(ctx context.Context, userID string) ([]string, error) {
	key := s.getKey(fmt.Sprintf("purchase:%s:chapters", userID))

	chapters, err := s.redisClient.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("获取用户购买章节失败: %w", err)
	}

	return chapters, nil
}

// GetVIPExpireTime 获取用户VIP过期时间
func (s *VIPPermissionServiceImpl) GetVIPExpireTime(ctx context.Context, userID string) (time.Duration, error) {
	key := s.getKey(fmt.Sprintf("user:%s:status", userID))

	ttl, err := s.redisClient.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("获取VIP过期时间失败: %w", err)
	}

	return ttl, nil
}
