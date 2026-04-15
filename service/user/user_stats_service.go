package user

import (
	"Qingyu_backend/models/users"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// UserStats 用户统计数据
type UserStats struct {
	UserID        string    `json:"user_id"`
	TotalProjects int64     `json:"total_projects"`
	TotalBooks    int64     `json:"total_books"`
	TotalWords    int64     `json:"total_words"`
	TotalReading  int64     `json:"total_reading"`
	TotalLikes    int64     `json:"total_likes"`
	TotalComments int64     `json:"total_comments"`
	TotalRevenue  float64   `json:"total_revenue"`
	MemberLevel   string    `json:"member_level"`
	RegisteredAt  time.Time `json:"registered_at"`
	LastActiveAt  time.Time `json:"last_active_at"`
	ActiveDays    int       `json:"active_days"`
}

// UserStatsUserRepository 用户统计所需的最小用户仓储能力。
type UserStatsUserRepository interface {
	GetByID(ctx context.Context, id string) (*users.User, error)
}

// UserStatsProjectRepository 用户统计所需的最小项目仓储能力。
type UserStatsProjectRepository interface {
	CountByOwner(ctx context.Context, ownerID string) (int64, error)
}

// UserStatsService 用户域统计服务端口（切片 B：仅承接 GetUserStats）。
type UserStatsService interface {
	GetUserStats(ctx context.Context, userID string) (*UserStats, error)
}

// UserStatsServiceImpl 用户域统计服务实现。
type UserStatsServiceImpl struct {
	userRepo    UserStatsUserRepository
	projectRepo UserStatsProjectRepository
}

var _ UserStatsService = (*UserStatsServiceImpl)(nil)

// NewUserStatsService 创建用户域统计服务。
func NewUserStatsService(userRepo UserStatsUserRepository, projectRepo UserStatsProjectRepository) UserStatsService {
	return &UserStatsServiceImpl{
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}

// GetUserStats 获取用户统计。
func (s *UserStatsServiceImpl) GetUserStats(ctx context.Context, userID string) (*UserStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if s.userRepo == nil {
		return nil, fmt.Errorf("用户仓储未配置")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		zap.L().Error("获取用户信息失败",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	projectCount := int64(0)
	if s.projectRepo != nil {
		projectCount, err = s.projectRepo.CountByOwner(ctx, userID)
		if err != nil {
			zap.L().Warn("统计项目数失败，使用默认值0",
				zap.String("user_id", userID),
				zap.Error(err),
			)
			projectCount = 0
		}
	}

	memberLevel := "普通读者"
	if user.IsAdmin() {
		memberLevel = "管理员"
	} else if user.IsAuthor() {
		memberLevel = "作者"
	}
	if user.IsVIP() {
		memberLevel += fmt.Sprintf(" (VIP Level %d)", user.GetVIPLevel())
	}

	return &UserStats{
		UserID:        userID,
		TotalProjects: projectCount,
		TotalBooks:    0, // TODO(Task3): 需要支持string author_id查询
		TotalWords:    0, // TODO(Task3): 需要支持string author_id查询
		TotalReading:  0, // TODO(Task3): 需要阅读行为统计
		TotalLikes:    0, // TODO(Task3): 需要点赞统计
		TotalComments: 0, // TODO(Task3): 需要评论统计
		TotalRevenue:  0, // TODO(Task3): 需要钱包统计
		MemberLevel:   memberLevel,
		RegisteredAt:  user.CreatedAt,
		LastActiveAt:  user.UpdatedAt,
		ActiveDays:    0, // TODO(Task3): 需要活跃度统计
	}, nil
}
