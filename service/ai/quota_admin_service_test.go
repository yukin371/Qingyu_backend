package ai

import (
	"context"
	"testing"
	"time"

	aiModels "Qingyu_backend/models/ai"
	pb "Qingyu_backend/pkg/grpc/pb"
	aiRepo "Qingyu_backend/repository/interfaces/ai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type quotaAdminRepoStub struct {
	transactions []*aiModels.QuotaTransaction
}

func (s *quotaAdminRepoStub) CreateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	return nil
}

func (s *quotaAdminRepoStub) GetQuotaByUserID(ctx context.Context, userID string, quotaType aiModels.QuotaType) (*aiModels.UserQuota, error) {
	return nil, aiModels.ErrQuotaNotFound
}

func (s *quotaAdminRepoStub) UpdateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	return nil
}

func (s *quotaAdminRepoStub) DeleteQuota(ctx context.Context, userID string, quotaType aiModels.QuotaType) error {
	return nil
}

func (s *quotaAdminRepoStub) GetAllQuotasByUserID(ctx context.Context, userID string) ([]*aiModels.UserQuota, error) {
	return nil, nil
}

func (s *quotaAdminRepoStub) BatchResetQuotas(ctx context.Context, quotaType aiModels.QuotaType) error {
	return nil
}

func (s *quotaAdminRepoStub) CreateTransaction(ctx context.Context, transaction *aiModels.QuotaTransaction) error {
	return nil
}

func (s *quotaAdminRepoStub) GetTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*aiModels.QuotaTransaction, error) {
	return nil, nil
}

func (s *quotaAdminRepoStub) GetTransactionsByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*aiModels.QuotaTransaction, error) {
	result := make([]*aiModels.QuotaTransaction, 0, len(s.transactions))
	for _, tx := range s.transactions {
		if tx == nil || tx.UserID != userID {
			continue
		}
		if tx.Timestamp.Before(startTime) || tx.Timestamp.After(endTime) {
			continue
		}
		result = append(result, tx)
	}
	return result, nil
}

func (s *quotaAdminRepoStub) GetQuotaStatistics(ctx context.Context, userID string) (*aiRepo.QuotaStatistics, error) {
	return nil, nil
}

func (s *quotaAdminRepoStub) GetTotalConsumption(ctx context.Context, userID string, quotaType aiModels.QuotaType, startTime, endTime time.Time) (int, error) {
	return 0, nil
}

func (s *quotaAdminRepoStub) GetDashboardSummary(ctx context.Context) (*aiModels.DashboardSummary, error) {
	return nil, nil
}

func (s *quotaAdminRepoStub) GetQuotaDistribution(ctx context.Context) (*aiModels.QuotaDistribution, error) {
	return nil, nil
}

func (s *quotaAdminRepoStub) GetTopConsumers(ctx context.Context, limit int) ([]aiModels.UserQuotaRanking, error) {
	return nil, nil
}

func (s *quotaAdminRepoStub) GetConsumptionTrend(ctx context.Context, days int) ([]aiModels.TrendPoint, error) {
	return nil, nil
}

func (s *quotaAdminRepoStub) GetConsumptionSummary(ctx context.Context, startTime, endTime time.Time, workflowType, groupBy string, page, pageSize int) (*aiModels.QuotaConsumptionSummary, error) {
	return nil, nil
}

func (s *quotaAdminRepoStub) ListUserQuotas(ctx context.Context, role, status, search string, page, limit int) ([]*aiModels.UserQuotaListItem, int64, error) {
	return nil, 0, nil
}

func (s *quotaAdminRepoStub) Health(ctx context.Context) error {
	return nil
}

type quotaConsumptionReaderStub struct {
	response *pb.QuotaConsumptionResponse
	err      error
}

func (s *quotaConsumptionReaderStub) GetQuotaConsumption(ctx context.Context, userID string, timeRange string, workflowType string) (*pb.QuotaConsumptionResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.response, nil
}

func TestQuotaAdminServiceGetUserQuotaReconciliation(t *testing.T) {
	now := time.Now()
	repo := &quotaAdminRepoStub{
		transactions: []*aiModels.QuotaTransaction{
			{
				UserID:    "user-1",
				Type:      "consume",
				Service:   "story_write",
				Amount:    120,
				Timestamp: now.Add(-time.Hour),
			},
			{
				UserID:    "user-1",
				Type:      "consume",
				Service:   "story_write",
				Amount:    80,
				Timestamp: now.Add(-2 * time.Hour),
			},
			{
				UserID:    "user-1",
				Type:      "consume",
				Service:   "outline",
				Amount:    50,
				Timestamp: now.Add(-3 * time.Hour),
			},
			{
				UserID:    "user-1",
				Type:      "restore",
				Service:   "story_write",
				Amount:    -20,
				Timestamp: now.Add(-30 * time.Minute),
			},
		},
	}
	service := NewQuotaAdminService(repo, nil)
	service.SetConsumptionReader(&quotaConsumptionReaderStub{
		response: &pb.QuotaConsumptionResponse{
			Success:      true,
			TotalTokens:  150,
			TotalRecords: 2,
			Records: []*pb.QuotaRecord{
				{Id: "r1", WorkflowType: "story_write", TokensUsed: 70, ConsumedAt: now.Add(-time.Hour).Format(time.RFC3339)},
				{Id: "r2", WorkflowType: "story_write", TokensUsed: 80, ConsumedAt: now.Add(-2 * time.Hour).Format(time.RFC3339)},
			},
		},
	})

	result, err := service.GetUserQuotaReconciliation(context.Background(), "user-1", "day", "story_write")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "user-1", result.UserID)
	assert.Equal(t, "day", result.TimeRange)
	assert.Equal(t, "story_write", result.WorkflowType)
	assert.Equal(t, string(aiModels.QuotaTypeDaily), result.BackendQuotaType)
	assert.Equal(t, 200, result.BackendTotalTokens)
	assert.Equal(t, 2, result.BackendRecordCount)
	assert.Equal(t, 150, result.AIServiceTotalTokens)
	assert.Equal(t, 2, result.AIServiceRecordCount)
	assert.Equal(t, 50, result.DifferenceTokens)
	assert.InDelta(t, 0.25, result.DifferenceRatio, 0.0001)
	assert.Equal(t, string(aiModels.QuotaAlertLevelCritical), result.AlertLevel)
	assert.True(t, result.ShouldAlert)
	require.Len(t, result.Records, 2)
	assert.Equal(t, "r1", result.Records[0].ID)
}

func TestResolveQuotaReconciliationWindowRejectsInvalidRange(t *testing.T) {
	_, _, _, err := resolveQuotaReconciliationWindow("quarter", time.Now())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "不支持的时间范围")
}
