package writer

import (
	"context"
	"time"

	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/service/interfaces"
)

// DashboardStats 仪表板统计数据
type DashboardStats struct {
	TotalWords int64 `json:"totalWords"`
	BookCount  int64 `json:"bookCount"`
	TodayWords int64 `json:"todayWords"`
	Pending    int64 `json:"pending"`
	Streak     int   `json:"streak"`
}

// DashboardService 仪表板统计服务
type DashboardService struct {
	projectRepo   writerRepo.ProjectRepository
	publishService interfaces.PublishService
}

// NewDashboardService 创建仪表板统计服务
func NewDashboardService(projectRepo writerRepo.ProjectRepository, publishService interfaces.PublishService) *DashboardService {
	return &DashboardService{
		projectRepo:   projectRepo,
		publishService: publishService,
	}
}

// GetStats 获取作者仪表板统计数据
func (s *DashboardService) GetStats(ctx context.Context, userID string) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 项目数量
	bookCount, err := s.projectRepo.CountByOwner(ctx, userID)
	if err != nil {
		return nil, err
	}
	stats.BookCount = bookCount

	// 获取所有项目，汇总字数
	projects, err := s.projectRepo.GetListByOwnerID(ctx, userID, 1000, 0)
	if err != nil {
		return nil, err
	}

	var totalWords int64
	today := time.Now().Truncate(24 * time.Hour)

	var todayWords int64

	for _, p := range projects {
		totalWords += int64(p.Statistics.TotalWords)
		if p.UpdatedAt.After(today) {
			todayWords += int64(p.Statistics.TotalWords)
		}
	}
	stats.TotalWords = totalWords
	stats.TodayWords = todayWords

	// 待审核数量
	if s.publishService != nil {
		_, pendingCount, err := s.publishService.GetPendingPublicationRecords(ctx, 1, 1)
		if err == nil {
			stats.Pending = pendingCount
		}
	}

	// 连续写作天数 - 暂时返回0，后续通过 Redis 实现
	stats.Streak = 0

	return stats, nil
}
