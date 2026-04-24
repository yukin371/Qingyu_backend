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
	summary      *aiModels.DashboardSummary
	distribution *aiModels.QuotaDistribution
	topConsumers []aiModels.UserQuotaRanking
	trend        []aiModels.TrendPoint
	consumption  *aiModels.QuotaConsumptionSummary
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
	response *pb.QuotaConsumptionSummaryResponse
}

func (s *quotaSummaryReaderStub) GetQuotaConsumptionSummary(
	ctx context.Context,
	timeRange string,
	workflowType string,
	groupBy string,
	page int32,
	pageSize int32,
) (*pb.QuotaConsumptionSummaryResponse, error) {
	return s.response, nil
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
