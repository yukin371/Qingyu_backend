package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	aiModels "Qingyu_backend/models/ai"
	"Qingyu_backend/pkg/cache"
	pb "Qingyu_backend/pkg/grpc/pb"
	aiInterfaces "Qingyu_backend/repository/interfaces/ai"
)

// QuotaDashboardService 配额仪表盘服务
type QuotaDashboardService struct {
	quotaRepo         aiInterfaces.QuotaRepository
	alertRepo         aiInterfaces.QuotaAlertRepository
	redisClient       cache.RedisClient
	cacheTTL          time.Duration
	consumptionReader quotaConsumptionSummaryReader
	consistencyRunner quotaConsistencyRunner
}

type quotaConsumptionSummaryReader interface {
	GetQuotaConsumptionSummary(
		ctx context.Context,
		timeRange string,
		workflowType string,
		groupBy string,
		page int32,
		pageSize int32,
	) (*pb.QuotaConsumptionSummaryResponse, error)
}

type quotaConsistencyRunner interface {
	RunConsistencyCheck(ctx context.Context) error
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

// SetConsumptionSummaryReader 设置 AI 服务聚合对账读取客户端。
func (s *QuotaDashboardService) SetConsumptionSummaryReader(reader quotaConsumptionSummaryReader) {
	s.consumptionReader = reader
}

// SetConsistencyRunner 设置手动触发一致性检查的执行器。
func (s *QuotaDashboardService) SetConsistencyRunner(runner quotaConsistencyRunner) {
	s.consistencyRunner = runner
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
		Summary:      *summary,
		Distribution: *distribution,
		TopConsumers: topConsumers,
		RecentAlerts: s.formatAlertsToSummary(recentAlerts),
		TrendData:    trendData,
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

// GetReconciliationSummary 获取全局聚合对账摘要。
func (s *QuotaDashboardService) GetReconciliationSummary(
	ctx context.Context,
	timeRange string,
	workflowType string,
	groupBy string,
	page int,
	pageSize int,
) (*aiModels.QuotaConsumptionReconciliationSummary, error) {
	if s.consumptionReader == nil {
		return nil, fmt.Errorf("AI配额聚合对账客户端未配置")
	}

	normalizedRange, windowStart, windowEnd, err := resolveQuotaReconciliationWindow(timeRange, time.Now())
	if err != nil {
		return nil, err
	}
	normalizedGroupBy := normalizeQuotaSummaryGroupBy(groupBy)
	page, pageSize = normalizeQuotaSummaryPagination(page, pageSize)

	backendSummary, err := s.quotaRepo.GetConsumptionSummary(ctx, windowStart, windowEnd, workflowType, normalizedGroupBy, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取后端消费聚合摘要失败: %w", err)
	}
	backendSummary.TimeRange = normalizedRange

	aiResp, err := s.consumptionReader.GetQuotaConsumptionSummary(
		ctx,
		normalizedRange,
		workflowType,
		normalizedGroupBy,
		int32(page),
		int32(pageSize),
	)
	if err != nil {
		return nil, fmt.Errorf("获取AI服务消费聚合摘要失败: %w", err)
	}
	if !aiResp.GetSuccess() {
		if aiResp.GetErrorMessage() != "" {
			return nil, fmt.Errorf("AI服务消费聚合摘要失败: %s", aiResp.GetErrorMessage())
		}
		return nil, fmt.Errorf("AI服务消费聚合摘要返回未成功状态")
	}

	aiSummary := convertProtoQuotaConsumptionSummary(aiResp)
	items := buildQuotaConsumptionReconciliationItems(normalizedGroupBy, backendSummary.Items, aiSummary.Items)
	level, shouldAlert := determineConsistencyAlertLevelForScopeInt64(
		quotaConsistencyScopeGlobal,
		backendSummary.TotalTokens,
		aiSummary.TotalTokens,
	)

	return &aiModels.QuotaConsumptionReconciliationSummary{
		TimeRange:             normalizedRange,
		WorkflowType:          workflowType,
		GroupBy:               normalizedGroupBy,
		Page:                  page,
		PageSize:              pageSize,
		TotalGroups:           maxInt64(backendSummary.TotalGroups, aiSummary.TotalGroups),
		BackendTotalTokens:    backendSummary.TotalTokens,
		BackendTotalRecords:   backendSummary.TotalRecords,
		AIServiceTotalTokens:  aiSummary.TotalTokens,
		AIServiceTotalRecords: aiSummary.TotalRecords,
		DifferenceTokens:      absInt64(backendSummary.TotalTokens - aiSummary.TotalTokens),
		DifferenceRatio:       calculateDifferenceRatioInt64(backendSummary.TotalTokens, aiSummary.TotalTokens),
		AlertLevel:            string(level),
		ShouldAlert:           shouldAlert,
		CheckedAt:             time.Now(),
		Items:                 items,
	}, nil
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

// RunConsistencyCheck 手动执行一次跨服务一致性对账检查。
func (s *QuotaDashboardService) RunConsistencyCheck(ctx context.Context) error {
	if s.consistencyRunner == nil {
		return fmt.Errorf("AI配额对账执行器未配置")
	}
	return s.consistencyRunner.RunConsistencyCheck(ctx)
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

func convertProtoQuotaConsumptionSummary(resp *pb.QuotaConsumptionSummaryResponse) *aiModels.QuotaConsumptionSummary {
	items := make([]aiModels.QuotaConsumptionSummaryItem, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, aiModels.QuotaConsumptionSummaryItem{
			GroupKey:     item.GetGroupKey(),
			TotalTokens:  int64(item.GetTotalTokens()),
			TotalRecords: int64(item.GetTotalRecords()),
		})
	}

	return &aiModels.QuotaConsumptionSummary{
		TimeRange:    resp.GetTimeRange(),
		WorkflowType: resp.GetWorkflowType(),
		GroupBy:      resp.GetGroupBy(),
		Page:         int(resp.GetPage()),
		PageSize:     int(resp.GetPageSize()),
		TotalGroups:  int64(resp.GetTotalGroups()),
		TotalTokens:  int64(resp.GetTotalTokens()),
		TotalRecords: int64(resp.GetTotalRecords()),
		Items:        items,
	}
}

func buildQuotaConsumptionReconciliationItems(
	groupBy string,
	backendItems []aiModels.QuotaConsumptionSummaryItem,
	aiItems []aiModels.QuotaConsumptionSummaryItem,
) []aiModels.QuotaConsumptionReconciliationItem {
	backendMap := make(map[string]aiModels.QuotaConsumptionSummaryItem, len(backendItems))
	aiMap := make(map[string]aiModels.QuotaConsumptionSummaryItem, len(aiItems))
	order := make([]string, 0, len(backendItems)+len(aiItems))

	for _, item := range backendItems {
		backendMap[item.GroupKey] = item
		order = append(order, item.GroupKey)
	}
	for _, item := range aiItems {
		aiMap[item.GroupKey] = item
		if _, exists := backendMap[item.GroupKey]; !exists {
			order = append(order, item.GroupKey)
		}
	}

	result := make([]aiModels.QuotaConsumptionReconciliationItem, 0, len(order))
	scope := quotaConsistencyScopeUser
	if normalizeQuotaSummaryGroupBy(groupBy) == quotaConsistencyScopeWorkflow {
		scope = quotaConsistencyScopeWorkflow
	}
	for _, groupKey := range order {
		backendItem := backendMap[groupKey]
		aiItem := aiMap[groupKey]
		level, shouldAlert := determineConsistencyAlertLevelForScopeInt64(scope, backendItem.TotalTokens, aiItem.TotalTokens)
		result = append(result, aiModels.QuotaConsumptionReconciliationItem{
			GroupKey:         groupKey,
			BackendTokens:    backendItem.TotalTokens,
			BackendRecords:   backendItem.TotalRecords,
			AIServiceTokens:  aiItem.TotalTokens,
			AIServiceRecords: aiItem.TotalRecords,
			DifferenceTokens: absInt64(backendItem.TotalTokens - aiItem.TotalTokens),
			DifferenceRatio:  calculateDifferenceRatioInt64(backendItem.TotalTokens, aiItem.TotalTokens),
			AlertLevel:       string(level),
			ShouldAlert:      shouldAlert,
		})
	}

	return result
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
