package ai

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/config"
	aiModels "Qingyu_backend/models/ai"
	pb "Qingyu_backend/pkg/grpc/pb"
	aiRepo "Qingyu_backend/repository/interfaces/ai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type quotaDashboardRepoStub struct {
	summary         *aiModels.DashboardSummary
	distribution    *aiModels.QuotaDistribution
	topConsumers    []aiModels.UserQuotaRanking
	topConsumersErr error
	trend           []aiModels.TrendPoint
	consumption     *aiModels.QuotaConsumptionSummary
}

func (s *quotaDashboardRepoStub) CreateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	return nil
}

func (s *quotaDashboardRepoStub) GetQuotaByUserID(ctx context.Context, userID string, quotaType aiModels.QuotaType) (*aiModels.UserQuota, error) {
	return nil, aiModels.ErrQuotaNotFound
}

func (s *quotaDashboardRepoStub) UpdateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	return nil
}

func (s *quotaDashboardRepoStub) DeleteQuota(ctx context.Context, userID string, quotaType aiModels.QuotaType) error {
	return nil
}

func (s *quotaDashboardRepoStub) GetAllQuotasByUserID(ctx context.Context, userID string) ([]*aiModels.UserQuota, error) {
	return nil, nil
}

func (s *quotaDashboardRepoStub) BatchResetQuotas(ctx context.Context, quotaType aiModels.QuotaType) error {
	return nil
}

func (s *quotaDashboardRepoStub) CreateTransaction(ctx context.Context, transaction *aiModels.QuotaTransaction) error {
	return nil
}

func (s *quotaDashboardRepoStub) GetTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*aiModels.QuotaTransaction, error) {
	return nil, nil
}

func (s *quotaDashboardRepoStub) GetTransactionsByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*aiModels.QuotaTransaction, error) {
	return nil, nil
}

func (s *quotaDashboardRepoStub) GetQuotaStatistics(ctx context.Context, userID string) (*aiRepo.QuotaStatistics, error) {
	return nil, nil
}

func (s *quotaDashboardRepoStub) GetTotalConsumption(ctx context.Context, userID string, quotaType aiModels.QuotaType, startTime, endTime time.Time) (int, error) {
	return 0, nil
}

func (s *quotaDashboardRepoStub) GetDashboardSummary(ctx context.Context) (*aiModels.DashboardSummary, error) {
	if s.summary != nil {
		return s.summary, nil
	}
	return &aiModels.DashboardSummary{}, nil
}

func (s *quotaDashboardRepoStub) GetQuotaDistribution(ctx context.Context) (*aiModels.QuotaDistribution, error) {
	if s.distribution != nil {
		return s.distribution, nil
	}
	return &aiModels.QuotaDistribution{}, nil
}

func (s *quotaDashboardRepoStub) GetTopConsumers(ctx context.Context, limit int) ([]aiModels.UserQuotaRanking, error) {
	if s.topConsumersErr != nil {
		return nil, s.topConsumersErr
	}
	if len(s.topConsumers) > 0 {
		return s.topConsumers, nil
	}
	return nil, nil
}

func (s *quotaDashboardRepoStub) GetConsumptionTrend(ctx context.Context, days int) ([]aiModels.TrendPoint, error) {
	if len(s.trend) > 0 {
		return s.trend, nil
	}
	return nil, nil
}

func (s *quotaDashboardRepoStub) GetConsumptionSummary(ctx context.Context, startTime, endTime time.Time, workflowType, groupBy string, page, pageSize int) (*aiModels.QuotaConsumptionSummary, error) {
	if s.consumption != nil {
		return s.consumption, nil
	}
	return &aiModels.QuotaConsumptionSummary{}, nil
}

func (s *quotaDashboardRepoStub) ListUserQuotas(ctx context.Context, role, status, search string, page, limit int) ([]*aiModels.UserQuotaListItem, int64, error) {
	return nil, 0, nil
}

func (s *quotaDashboardRepoStub) Health(ctx context.Context) error {
	return nil
}

type quotaSummaryReaderStub struct {
	response         *pb.QuotaConsumptionSummaryResponse
	err              error
	lastTimeRange    string
	lastWorkflowType string
	lastGroupBy      string
	lastPage         int32
	lastPageSize     int32
}

func (s *quotaSummaryReaderStub) GetQuotaConsumptionSummary(
	ctx context.Context,
	timeRange string,
	workflowType string,
	groupBy string,
	page int32,
	pageSize int32,
) (*pb.QuotaConsumptionSummaryResponse, error) {
	s.lastTimeRange = timeRange
	s.lastWorkflowType = workflowType
	s.lastGroupBy = groupBy
	s.lastPage = page
	s.lastPageSize = pageSize
	if s.err != nil {
		return nil, s.err
	}
	return s.response, nil
}

type quotaConsistencyRunnerStub struct {
	called int
	err    error
}

func (s *quotaConsistencyRunnerStub) RunConsistencyCheck(ctx context.Context) error {
	s.called++
	return s.err
}

func TestQuotaDashboardServiceGetReconciliationSummaryUsesConfiguredGlobalThresholds(t *testing.T) {
	previousConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{
		AIQuota: &config.AIQuotaConfig{
			ConsistencyThresholds: &config.QuotaConsistencyThresholdsConfig{
				Global: &config.QuotaConsistencyThresholdConfig{
					WarningTokens:  500,
					CriticalTokens: 2000,
					WarningRatio:   0.4,
					CriticalRatio:  0.8,
				},
				Workflow: &config.QuotaConsistencyThresholdConfig{
					WarningTokens:  500,
					CriticalTokens: 2000,
					WarningRatio:   0.4,
					CriticalRatio:  0.8,
				},
			},
		},
	}
	t.Cleanup(func() {
		config.GlobalConfig = previousConfig
	})

	repo := &quotaDashboardRepoStub{
		consumption: &aiModels.QuotaConsumptionSummary{
			TimeRange:    "day",
			GroupBy:      "workflow",
			Page:         1,
			PageSize:     20,
			TotalGroups:  1,
			TotalTokens:  1000,
			TotalRecords: 2,
			Items: []aiModels.QuotaConsumptionSummaryItem{
				{GroupKey: "rewrite", TotalTokens: 1000, TotalRecords: 2},
			},
		},
	}
	service := NewQuotaDashboardService(repo, nil, nil)
	service.SetConsumptionSummaryReader(&quotaSummaryReaderStub{
		response: &pb.QuotaConsumptionSummaryResponse{
			Success:      true,
			TimeRange:    "day",
			GroupBy:      "workflow",
			Page:         1,
			PageSize:     20,
			TotalGroups:  1,
			TotalTokens:  700,
			TotalRecords: 2,
			Items: []*pb.QuotaConsumptionSummaryItem{
				{GroupKey: "rewrite", TotalTokens: 700, TotalRecords: 2},
			},
		},
	})

	result, err := service.GetReconciliationSummary(context.Background(), "day", "", "workflow", 1, 20)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "workflow", result.GroupBy)
	assert.Equal(t, string(aiModels.QuotaAlertLevelInfo), result.AlertLevel)
	assert.False(t, result.ShouldAlert)
	require.Len(t, result.Items, 1)
	assert.Equal(t, string(aiModels.QuotaAlertLevelInfo), result.Items[0].AlertLevel)
}

func TestQuotaDashboardServiceGetReconciliationSummaryRequiresReader(t *testing.T) {
	service := NewQuotaDashboardService(&quotaDashboardRepoStub{}, nil, nil)

	result, err := service.GetReconciliationSummary(context.Background(), "day", "", "user", 1, 20)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "AI配额聚合对账客户端未配置")
}

func TestQuotaDashboardServiceGetReconciliationSummaryNormalizesInputAndHandlesFailedResponse(t *testing.T) {
	repo := &quotaDashboardRepoStub{
		consumption: &aiModels.QuotaConsumptionSummary{
			Items: []aiModels.QuotaConsumptionSummaryItem{{GroupKey: "user-1", TotalTokens: 10, TotalRecords: 1}},
		},
	}
	reader := &quotaSummaryReaderStub{
		response: &pb.QuotaConsumptionSummaryResponse{
			Success:      false,
			ErrorMessage: "ai unavailable",
		},
	}
	service := NewQuotaDashboardService(repo, nil, nil)
	service.SetConsumptionSummaryReader(reader)

	result, err := service.GetReconciliationSummary(context.Background(), "", "story_write", "unknown", 0, 999)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "day", reader.lastTimeRange)
	assert.Equal(t, "story_write", reader.lastWorkflowType)
	assert.Equal(t, "user", reader.lastGroupBy)
	assert.EqualValues(t, 1, reader.lastPage)
	assert.EqualValues(t, 100, reader.lastPageSize)
	assert.Contains(t, err.Error(), "ai unavailable")
}

func TestQuotaDashboardServiceGetReconciliationSummaryReturnsGenericFailureWhenAIResponseHasNoMessage(t *testing.T) {
	repo := &quotaDashboardRepoStub{
		consumption: &aiModels.QuotaConsumptionSummary{
			Items: []aiModels.QuotaConsumptionSummaryItem{{GroupKey: "user-1", TotalTokens: 10, TotalRecords: 1}},
		},
	}
	reader := &quotaSummaryReaderStub{
		response: &pb.QuotaConsumptionSummaryResponse{
			Success: false,
		},
	}
	service := NewQuotaDashboardService(repo, nil, nil)
	service.SetConsumptionSummaryReader(reader)

	result, err := service.GetReconciliationSummary(context.Background(), "day", "", "user", 1, 20)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "AI服务消费聚合摘要返回未成功状态")
}

func TestQuotaDashboardServiceGetDashboardReturnsEmptyCollectionsWhenRepositoriesAreEmpty(t *testing.T) {
	repo := &quotaDashboardRepoStub{
		topConsumers: []aiModels.UserQuotaRanking{},
		trend:        []aiModels.TrendPoint{},
	}
	alertRepo := newQuotaAlertRepoStub()
	service := NewQuotaDashboardService(repo, alertRepo, nil)

	result, err := service.GetDashboard(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(0), result.Summary.TotalUsers)
	assert.Equal(t, int64(0), result.Summary.TotalConsumption)
	require.Len(t, result.TopConsumers, 0)
	require.Len(t, result.RecentAlerts, 0)
	require.Len(t, result.TrendData, 0)
}

func TestQuotaDashboardServiceRunConsistencyCheckRequiresRunnerAndDelegates(t *testing.T) {
	service := NewQuotaDashboardService(&quotaDashboardRepoStub{}, nil, nil)

	err := service.RunConsistencyCheck(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "AI配额对账执行器未配置")

	runner := &quotaConsistencyRunnerStub{}
	service.SetConsistencyRunner(runner)

	err = service.RunConsistencyCheck(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, runner.called)
}
