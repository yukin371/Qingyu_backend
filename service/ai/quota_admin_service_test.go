package ai

import (
	"context"
	"errors"
	"testing"
	"time"

	aiModels "Qingyu_backend/models/ai"
	pb "Qingyu_backend/pkg/grpc/pb"
	aiRepo "Qingyu_backend/repository/interfaces/ai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type quotaAdminRepoStub struct {
	transactions           []*aiModels.QuotaTransaction
	quotasByUserID         map[string]map[aiModels.QuotaType]*aiModels.UserQuota
	allQuotasByUserID      map[string][]*aiModels.UserQuota
	getQuotaErrByUserID    map[string]error
	updateQuotaErrByUserID map[string]error
	createQuotaErrByUserID map[string]error
	createTransactionErr   error
}

func newQuotaAdminRepoStub() *quotaAdminRepoStub {
	return &quotaAdminRepoStub{
		quotasByUserID:    make(map[string]map[aiModels.QuotaType]*aiModels.UserQuota),
		allQuotasByUserID: make(map[string][]*aiModels.UserQuota),
	}
}

func quotaAdminStubQuota(userID string, quotaType aiModels.QuotaType) string {
	return userID + "|" + string(quotaType)
}

func (s *quotaAdminRepoStub) CreateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	if quota == nil {
		return nil
	}
	if s.createQuotaErrByUserID != nil {
		if err, ok := s.createQuotaErrByUserID[quota.UserID]; ok {
			return err
		}
	}
	if s.quotasByUserID == nil {
		s.quotasByUserID = make(map[string]map[aiModels.QuotaType]*aiModels.UserQuota)
	}
	if s.quotasByUserID[quota.UserID] == nil {
		s.quotasByUserID[quota.UserID] = make(map[aiModels.QuotaType]*aiModels.UserQuota)
	}
	s.quotasByUserID[quota.UserID][quota.QuotaType] = quota
	if s.allQuotasByUserID != nil {
		s.allQuotasByUserID[quota.UserID] = append(s.allQuotasByUserID[quota.UserID], quota)
	}
	return nil
}

func (s *quotaAdminRepoStub) GetQuotaByUserID(ctx context.Context, userID string, quotaType aiModels.QuotaType) (*aiModels.UserQuota, error) {
	if s.getQuotaErrByUserID != nil {
		if err, ok := s.getQuotaErrByUserID[userID]; ok {
			return nil, err
		}
	}
	if s.quotasByUserID != nil {
		if quotaByType, ok := s.quotasByUserID[userID]; ok {
			if quota, ok := quotaByType[quotaType]; ok {
				return quota, nil
			}
		}
	}
	return nil, aiModels.ErrQuotaNotFound
}

func (s *quotaAdminRepoStub) UpdateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	if quota == nil {
		return nil
	}
	if s.updateQuotaErrByUserID != nil {
		if err, ok := s.updateQuotaErrByUserID[quota.UserID]; ok {
			return err
		}
	}
	if s.quotasByUserID == nil {
		s.quotasByUserID = make(map[string]map[aiModels.QuotaType]*aiModels.UserQuota)
	}
	if s.quotasByUserID[quota.UserID] == nil {
		s.quotasByUserID[quota.UserID] = make(map[aiModels.QuotaType]*aiModels.UserQuota)
	}
	s.quotasByUserID[quota.UserID][quota.QuotaType] = quota
	if s.allQuotasByUserID != nil {
		quotas := s.allQuotasByUserID[quota.UserID]
		replaced := false
		for i, existing := range quotas {
			if existing != nil && existing.QuotaType == quota.QuotaType {
				quotas[i] = quota
				replaced = true
				break
			}
		}
		if !replaced {
			quotas = append(quotas, quota)
		}
		s.allQuotasByUserID[quota.UserID] = quotas
	}
	return nil
}

func (s *quotaAdminRepoStub) DeleteQuota(ctx context.Context, userID string, quotaType aiModels.QuotaType) error {
	if s.quotasByUserID != nil {
		if quotaByType, ok := s.quotasByUserID[userID]; ok {
			delete(quotaByType, quotaType)
		}
	}
	return nil
}

func (s *quotaAdminRepoStub) GetAllQuotasByUserID(ctx context.Context, userID string) ([]*aiModels.UserQuota, error) {
	if s.allQuotasByUserID != nil {
		if quotas, ok := s.allQuotasByUserID[userID]; ok {
			return quotas, nil
		}
	}
	if s.quotasByUserID == nil {
		return nil, nil
	}
	quotaByType, ok := s.quotasByUserID[userID]
	if !ok {
		return nil, nil
	}
	quotas := make([]*aiModels.UserQuota, 0, len(quotaByType))
	for _, quota := range quotaByType {
		quotas = append(quotas, quota)
	}
	return quotas, nil
}

func (s *quotaAdminRepoStub) BatchResetQuotas(ctx context.Context, quotaType aiModels.QuotaType) error {
	return nil
}

func (s *quotaAdminRepoStub) CreateTransaction(ctx context.Context, transaction *aiModels.QuotaTransaction) error {
	if s.createTransactionErr != nil {
		return s.createTransactionErr
	}
	s.transactions = append(s.transactions, transaction)
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
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	elapsed := now.Sub(dayStart)
	firstConsumeAt := dayStart.Add(elapsed / 4)
	secondConsumeAt := dayStart.Add(elapsed / 2)
	otherWorkflowAt := dayStart.Add((elapsed * 3) / 4)
	restoreAt := now
	repo := &quotaAdminRepoStub{
		transactions: []*aiModels.QuotaTransaction{
			{
				UserID:    "user-1",
				Type:      "consume",
				Service:   "story_write",
				Amount:    120,
				Timestamp: firstConsumeAt,
			},
			{
				UserID:    "user-1",
				Type:      "consume",
				Service:   "story_write",
				Amount:    80,
				Timestamp: secondConsumeAt,
			},
			{
				UserID:    "user-1",
				Type:      "consume",
				Service:   "outline",
				Amount:    50,
				Timestamp: otherWorkflowAt,
			},
			{
				UserID:    "user-1",
				Type:      "restore",
				Service:   "story_write",
				Amount:    -20,
				Timestamp: restoreAt,
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
				{Id: "r1", WorkflowType: "story_write", TokensUsed: 70, ConsumedAt: firstConsumeAt.Format(time.RFC3339)},
				{Id: "r2", WorkflowType: "story_write", TokensUsed: 80, ConsumedAt: secondConsumeAt.Format(time.RFC3339)},
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

func TestQuotaAdminServiceGetUserQuotaReconciliationRequiresReader(t *testing.T) {
	service := NewQuotaAdminService(&quotaAdminRepoStub{}, nil)

	result, err := service.GetUserQuotaReconciliation(context.Background(), "user-1", "day", "")
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "AI配额对账客户端未配置")
}

func TestQuotaAdminServiceGetUserQuotaReconciliationHandlesFailedAIResponse(t *testing.T) {
	service := NewQuotaAdminService(&quotaAdminRepoStub{}, nil)
	service.SetConsumptionReader(&quotaConsumptionReaderStub{
		response: &pb.QuotaConsumptionResponse{
			Success:      false,
			ErrorMessage: "grpc unavailable",
		},
	})

	result, err := service.GetUserQuotaReconciliation(context.Background(), "user-1", "day", "")
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "grpc unavailable")
}

func TestQuotaAdminServiceGetUserQuotaReconciliationReturnsGenericErrorWhenAIResponseHasNoMessage(t *testing.T) {
	service := NewQuotaAdminService(&quotaAdminRepoStub{}, nil)
	service.SetConsumptionReader(&quotaConsumptionReaderStub{
		response: &pb.QuotaConsumptionResponse{
			Success: false,
		},
	})

	result, err := service.GetUserQuotaReconciliation(context.Background(), "user-1", "day", "")
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "AI服务返回未成功状态")
}

func TestQuotaAdminServiceBatchOperationsValidateAndAggregateResults(t *testing.T) {
	t.Run("batch recharge validates input and aggregates partial failures", func(t *testing.T) {
		service := NewQuotaAdminService(&quotaAdminRepoStub{}, newQuotaAlertRepoStub())

		result, err := service.BatchRecharge(context.Background(), nil, 10, "daily", "manual", "operator-1")
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "用户列表不能为空")

		result, err = service.BatchRecharge(context.Background(), []string{"user-1"}, 0, "daily", "manual", "operator-1")
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "充值金额必须大于0")

		repo := &quotaAdminRepoStub{
			quotasByUserID: map[string]map[aiModels.QuotaType]*aiModels.UserQuota{
				"user-1": {
					aiModels.QuotaTypeDaily: &aiModels.UserQuota{
						UserID:         "user-1",
						QuotaType:      aiModels.QuotaTypeDaily,
						TotalQuota:     100,
						RemainingQuota: 80,
						Status:         aiModels.QuotaStatusActive,
					},
				},
			},
			getQuotaErrByUserID: map[string]error{
				"user-2": errors.New("db unavailable"),
			},
		}
		service = NewQuotaAdminService(repo, newQuotaAlertRepoStub())

		result, err = service.BatchRecharge(context.Background(), []string{"user-1", "user-2"}, 20, "daily", "manual", "operator-1")
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, result.Total)
		assert.Equal(t, 1, result.Success)
		assert.Equal(t, 1, result.Failed)
		require.Len(t, result.Errors, 1)
		assert.Contains(t, result.Errors[0], "user-2")
	})

	t.Run("batch update validates quota type and aggregates partial failures", func(t *testing.T) {
		service := NewQuotaAdminService(&quotaAdminRepoStub{}, newQuotaAlertRepoStub())

		result, err := service.BatchUpdateQuota(context.Background(), nil, 100, "daily")
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "用户列表不能为空")

		result, err = service.BatchUpdateQuota(context.Background(), []string{"user-1"}, 100, "quarterly")
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "不支持的配额类型")

		repo := &quotaAdminRepoStub{
			quotasByUserID: map[string]map[aiModels.QuotaType]*aiModels.UserQuota{
				"user-1": {
					aiModels.QuotaTypeDaily: &aiModels.UserQuota{
						UserID:         "user-1",
						QuotaType:      aiModels.QuotaTypeDaily,
						TotalQuota:     50,
						UsedQuota:      10,
						RemainingQuota: 40,
						Status:         aiModels.QuotaStatusActive,
					},
				},
				"user-2": {
					aiModels.QuotaTypeDaily: &aiModels.UserQuota{
						UserID:         "user-2",
						QuotaType:      aiModels.QuotaTypeDaily,
						TotalQuota:     60,
						UsedQuota:      20,
						RemainingQuota: 40,
						Status:         aiModels.QuotaStatusActive,
					},
				},
			},
			updateQuotaErrByUserID: map[string]error{
				"user-2": errors.New("update failed"),
			},
		}
		service = NewQuotaAdminService(repo, newQuotaAlertRepoStub())

		result, err = service.BatchUpdateQuota(context.Background(), []string{"user-1", "user-2"}, 120, "daily")
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, result.Total)
		assert.Equal(t, 1, result.Success)
		assert.Equal(t, 1, result.Failed)
		require.Len(t, result.Errors, 1)
		assert.Contains(t, result.Errors[0], "user-2")
		assert.Contains(t, result.Errors[0], "更新配额失败")
	})

	t.Run("batch suspend and activate aggregate partial failures", func(t *testing.T) {
		suspendQuota := &aiModels.UserQuota{
			UserID:         "user-1",
			QuotaType:      aiModels.QuotaTypeDaily,
			TotalQuota:     80,
			RemainingQuota: 80,
			Status:         aiModels.QuotaStatusActive,
		}
		activateQuota := &aiModels.UserQuota{
			UserID:         "user-3",
			QuotaType:      aiModels.QuotaTypeDaily,
			TotalQuota:     40,
			RemainingQuota: 0,
			Status:         aiModels.QuotaStatusSuspended,
		}

		repo := &quotaAdminRepoStub{
			quotasByUserID: map[string]map[aiModels.QuotaType]*aiModels.UserQuota{
				"user-1": {
					aiModels.QuotaTypeDaily: suspendQuota,
				},
				"user-3": {
					aiModels.QuotaTypeDaily: activateQuota,
				},
			},
			allQuotasByUserID: map[string][]*aiModels.UserQuota{
				"user-1": []*aiModels.UserQuota{suspendQuota},
				"user-3": []*aiModels.UserQuota{activateQuota},
			},
			getQuotaErrByUserID: map[string]error{
				"user-2": errors.New("lookup failed"),
			},
		}
		service := NewQuotaAdminService(repo, newQuotaAlertRepoStub())

		suspendResult, err := service.BatchSuspend(context.Background(), []string{"user-1", "user-2"})
		require.NoError(t, err)
		require.NotNil(t, suspendResult)
		assert.Equal(t, 2, suspendResult.Total)
		assert.Equal(t, 1, suspendResult.Success)
		assert.Equal(t, 1, suspendResult.Failed)
		require.Len(t, suspendResult.Errors, 1)
		assert.Contains(t, suspendResult.Errors[0], "user-2")
		assert.Equal(t, aiModels.QuotaStatusSuspended, suspendQuota.Status)

		activateResult, err := service.BatchActivate(context.Background(), []string{"user-3", "user-2"})
		require.NoError(t, err)
		require.NotNil(t, activateResult)
		assert.Equal(t, 2, activateResult.Total)
		assert.Equal(t, 1, activateResult.Success)
		assert.Equal(t, 1, activateResult.Failed)
		require.Len(t, activateResult.Errors, 1)
		assert.Contains(t, activateResult.Errors[0], "user-2")
		assert.Equal(t, aiModels.QuotaStatusActive, activateQuota.Status)
	})
}
