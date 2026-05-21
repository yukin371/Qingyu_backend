package writer

import (
	"Qingyu_backend/models/users"
	"context"
	"fmt"

	"go.uber.org/zap"
)

// ContentStats 内容统计
type ContentStats struct {
	UserID             string  `json:"user_id"`
	TotalProjects      int64   `json:"total_projects"`
	PublishedBooks     int64   `json:"published_books"`
	DraftBooks         int64   `json:"draft_books"`
	TotalChapters      int64   `json:"total_chapters"`
	TotalWords         int64   `json:"total_words"`
	AverageWordsPerDay float64 `json:"average_words_per_day"`
	TotalViews         int64   `json:"total_views"`
	TotalCollections   int64   `json:"total_collections"`
	AverageRating      float64 `json:"average_rating"`
}

// ContentStatsUserRepository 内容统计所需的最小用户仓储能力。
type ContentStatsUserRepository interface {
	GetByID(ctx context.Context, id string) (*users.User, error)
}

// ContentStatsProjectRepository 内容统计所需的最小项目仓储能力。
type ContentStatsProjectRepository interface {
	CountByOwner(ctx context.Context, ownerID string) (int64, error)
}

// ContentStatsService writer 域内容统计服务端口（切片 C：仅承接 GetContentStats）。
type ContentStatsService interface {
	GetContentStats(ctx context.Context, userID string) (*ContentStats, error)
}

// ContentStatsServiceImpl writer 域内容统计服务实现。
type ContentStatsServiceImpl struct {
	userRepo    ContentStatsUserRepository
	projectRepo ContentStatsProjectRepository
}

var _ ContentStatsService = (*ContentStatsServiceImpl)(nil)

// NewContentStatsService 创建 writer 域内容统计服务。
func NewContentStatsService(
	userRepository ContentStatsUserRepository,
	projectRepository ContentStatsProjectRepository,
) ContentStatsService {
	return &ContentStatsServiceImpl{
		userRepo:    userRepository,
		projectRepo: projectRepository,
	}
}

// GetContentStats 获取内容统计。
func (s *ContentStatsServiceImpl) GetContentStats(ctx context.Context, userID string) (*ContentStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if s.userRepo == nil {
		return nil, fmt.Errorf("用户仓储未配置")
	}

	_, err := s.userRepo.GetByID(ctx, userID)
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

	return &ContentStats{
		UserID:             userID,
		TotalProjects:      projectCount,
		PublishedBooks:     0, // TODO(Task3): 需要支持 string author_id 查询
		DraftBooks:         0, // TODO(Task3): 需要支持 string author_id 查询
		TotalChapters:      0, // TODO(Task3): 需要支持 string author_id 查询
		TotalWords:         0, // TODO(Task3): 需要支持 string author_id 查询
		AverageWordsPerDay: 0,
		TotalViews:         0, // TODO(Task3): 需要阅读统计
		TotalCollections:   0, // TODO(Task3): 需要收藏统计
		AverageRating:      0, // TODO(Task3): 需要评分统计
	}, nil
}
