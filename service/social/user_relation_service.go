package social

import (
	"context"
	"errors"
	"fmt"

	socialModel "Qingyu_backend/models/social"
	socialRepo "Qingyu_backend/repository/interfaces/social"
)

// UserRelationService 用户关系服务接口
type UserRelationService interface {
	// 关注操作
	FollowUser(ctx context.Context, followerID, followeeID string) error
	UnfollowUser(ctx context.Context, followerID, followeeID string) error

	// 关系查询
	IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
	GetRelation(ctx context.Context, followerID, followeeID string) (*socialModel.UserRelation, error)

	// 列表查询
	GetFollowers(ctx context.Context, userID string, page, pageSize int) ([]*socialModel.UserRelation, int64, error)
	GetFollowing(ctx context.Context, userID string, page, pageSize int) ([]*socialModel.UserRelation, int64, error)

	// 统计
	GetFollowerCount(ctx context.Context, userID string) (int64, error)
	GetFollowingCount(ctx context.Context, userID string) (int64, error)
}

// UserRelationServiceImpl 用户关系服务实现
type UserRelationServiceImpl struct {
	repo socialRepo.UserRelationRepository
}

// NewUserRelationService 创建用户关系服务
func NewUserRelationService(repo socialRepo.UserRelationRepository) UserRelationService {
	return &UserRelationServiceImpl{
		repo: repo,
	}
}

// FollowUser 关注用户
func (s *UserRelationServiceImpl) FollowUser(ctx context.Context, followerID, followeeID string) error {
	// 参数验证
	if followerID == "" || followeeID == "" {
		return errors.New("用户ID不能为空")
	}
	if followerID == followeeID {
		return errors.New("不能关注自己")
	}

	// 检查是否已关注
	isFollowing, err := s.repo.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("检查关注关系失败: %w", err)
	}
	if isFollowing {
		return errors.New("已经关注了该用户")
	}

	// 检查是否有历史关系（取消关注后又重新关注）
	existingRelation, err := s.repo.GetRelation(ctx, followerID, followeeID)
	if err == nil && existingRelation != nil {
		// 重新激活关系
		updates := map[string]interface{}{
			"status": "active",
		}
		if err := s.repo.Update(ctx, existingRelation.ID, updates); err != nil {
			return fmt.Errorf("更新关注关系失败: %w", err)
		}
		return nil
	}

	// 创建新的关注关系
	relation := &socialModel.UserRelation{
		FollowerID: followerID,
		FolloweeID: followeeID,
		Status:     socialModel.RelationStatusActive,
	}

	if err := s.repo.Create(ctx, relation); err != nil {
		return fmt.Errorf("创建关注关系失败: %w", err)
	}

	return nil
}

// UnfollowUser 取消关注用户
func (s *UserRelationServiceImpl) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	// 参数验证
	if followerID == "" || followeeID == "" {
		return errors.New("用户ID不能为空")
	}

	// 获取关系
	relation, err := s.repo.GetRelation(ctx, followerID, followeeID)
	if err != nil {
		return errors.New("关注关系不存在")
	}

	// 检查是否已取消关注
	if relation.Status == socialModel.RelationStatusInactive {
		return errors.New("已经取消关注了")
	}

	// 更新状态为取消关注
	updates := map[string]interface{}{
		"status": "inactive",
	}
	if err := s.repo.Update(ctx, relation.ID, updates); err != nil {
		return fmt.Errorf("取消关注失败: %w", err)
	}

	return nil
}

// IsFollowing 检查是否关注
func (s *UserRelationServiceImpl) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	return s.repo.IsFollowing(ctx, followerID, followeeID)
}

// GetRelation 获取关系
func (s *UserRelationServiceImpl) GetRelation(ctx context.Context, followerID, followeeID string) (*socialModel.UserRelation, error) {
	return s.repo.GetRelation(ctx, followerID, followeeID)
}

// GetFollowers 获取粉丝列表
func (s *UserRelationServiceImpl) GetFollowers(ctx context.Context, userID string, page, pageSize int) ([]*socialModel.UserRelation, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.repo.GetFollowers(ctx, userID, pageSize, offset)
}

// GetFollowing 获取关注列表
func (s *UserRelationServiceImpl) GetFollowing(ctx context.Context, userID string, page, pageSize int) ([]*socialModel.UserRelation, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.repo.GetFollowing(ctx, userID, pageSize, offset)
}

// GetFollowerCount 获取粉丝数
func (s *UserRelationServiceImpl) GetFollowerCount(ctx context.Context, userID string) (int64, error) {
	return s.repo.CountFollowers(ctx, userID)
}

// GetFollowingCount 获取关注数
func (s *UserRelationServiceImpl) GetFollowingCount(ctx context.Context, userID string) (int64, error) {
	return s.repo.CountFollowing(ctx, userID)
}
