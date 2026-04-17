package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	aiModels "Qingyu_backend/models/ai"
	aiInterfaces "Qingyu_backend/repository/interfaces/ai"
	"Qingyu_backend/pkg/cache"

	)

// QuotaDashboardService 配额仪表盘服务
type QuotaDashboardService struct {
	quotaRepo   aiInterfaces.QuotaRepository
	alertRepo   aiInterfaces.QuotaAlertRepository
	redisClient cache.RedisClient
	cacheTTL    time.Duration
}

// NewQuotaDashboardService 创建配额仪表盘服务
func NewQuotaDashboardService(quotaRepo aiInterfaces.QuotaRepository, alertRepo aiInterfaces.QuotaAlertRepository, redisClient cache.RedisClient) *QuotaDashboardService {
	return &QuotaDashboardService{
		quotaRepo:   quotaRepo,
		alertRepo:   alertRepo,
		redisClient: redisClient,
		cacheTTL:    5 * time.Minute, // 默认5分钟缓存
	}
}

// GetDashboard 获取完整的仪表盘数据
func (s *QuotaDashboardService) GetDashboard(ctx context.Context) (*aiModels.QuotaDashboard, error) {
	// 尝试从缓存获取
	var dashboard *aiModels.QuotaDashboard
	cacheKey := "quota:dashboard"

	if s.redisClient != nil {
		cached, err := s.redisClient.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			if jsonErr := json.Unmarshal([]byte(cached), &dashboard); jsonErr == nil {
				return dashboard, nil
			}
		}
	}

	// 缓存未命中，从数据库获取
	summary, err := s.quotaRepo.GetDashboardSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取仪表盘汇总失败: %w", err)
	}

	distribution, err := s.quotaRepo.GetQuotaDistribution(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取配额分布失败: %w", err)
	}

	topConsumers, err := s.quotaRepo.GetTopConsumers(ctx, 10)
	if err != nil {
		return nil, fmt.Errorf("获取消费排行失败: %w", err)
	}

	trendData, err := s.quotaRepo.GetConsumptionTrend(ctx, 7)
	if err != nil {
		return nil, fmt.Errorf("获取趋势数据失败: %w", err)
	}

	recentAlerts, err := s.alertRepo.GetRecentGlobal(ctx, 5)
	if err != nil {
		return nil, fmt.Errorf("获取近期告警失败: %w", err)
	}

	dashboard = &aiModels.QuotaDashboard{
		Summary:       *summary,
		Distribution:   *distribution,
		TopConsumers:   topConsumers,
		RecentAlerts:   s.formatAlertsToSummary(recentAlerts),
		TrendData:      trendData,
	}

	// 写入缓存
	if s.redisClient != nil {
		s.cacheDashboard(ctx, dashboard)
	}

	return dashboard, nil
}

// GetDashboardSummary 获取仪表盘汇总数据
func (s *QuotaDashboardService) GetDashboardSummary(ctx context.Context) (*aiModels.DashboardSummary, error) {
	return s.quotaRepo.GetDashboardSummary(ctx)
}

// GetQuotaDistribution 获取配额分布统计
func (s *QuotaDashboardService) GetQuotaDistribution(ctx context.Context) (*aiModels.QuotaDistribution, error) {
	return s.quotaRepo.GetQuotaDistribution(ctx)
}

// GetTopConsumers 获取消费排行
func (s *QuotaDashboardService) GetTopConsumers(ctx context.Context, limit int) ([]aiModels.UserQuotaRanking, error) {
	return s.quotaRepo.GetTopConsumers(ctx, limit)
}

// GetConsumptionTrend 获取消费趋势
func (s *QuotaDashboardService) GetConsumptionTrend(ctx context.Context, days int) ([]aiModels.TrendPoint, error) {
	return s.quotaRepo.GetConsumptionTrend(ctx, days)
}

// RefreshDashboardCache 刷新仪表盘缓存
func (s *QuotaDashboardService) RefreshDashboardCache(ctx context.Context) error {
	// 强制清除缓存
	if s.redisClient != nil {
		cacheKey := "quota:dashboard"
		if err := s.redisClient.Delete(ctx, cacheKey); err != nil {
			fmt.Printf("清除仪表盘缓存失败: %v\n", err)
			return err
		}
	}

	// 重新构建缓存
	_, err := s.GetDashboard(ctx)
	return err
}

// 私有辅助方法

// cacheDashboard 缓存仪表盘数据
func (s *QuotaDashboardService) cacheDashboard(ctx context.Context, dashboard *aiModels.QuotaDashboard) {
	if s.redisClient == nil || dashboard == nil {
		return
	}

	cacheKey := "quota:dashboard"
	data, err := json.Marshal(dashboard)
	if err != nil {
		fmt.Printf("序列化仪表盘数据失败: %v\n", err)
		return
	}

	if err := s.redisClient.Set(ctx, cacheKey, string(data), s.cacheTTL); err != nil {
		fmt.Printf("写入仪表盘缓存失败: %v\n", err)
	}
}

// formatAlertsToSummary 将告警列表转换为摘要
func (s *QuotaDashboardService) formatAlertsToSummary(alerts []*aiModels.QuotaAlert) []aiModels.AlertSummary {
	result := make([]aiModels.AlertSummary, 0, len(alerts))
	for _, alert := range alerts {
		summary := aiModels.AlertSummary{
			ID:        alert.ID.Hex(),
			Type:      string(alert.Type),
			Level:     string(alert.Level),
			Title:     alert.Title,
			UserID:    alert.UserID,
			Status:    string(alert.Status),
			CreatedAt: alert.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		result = append(result, summary)
	}
	return result
}