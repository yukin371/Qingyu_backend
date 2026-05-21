package ai

import (
	"context"
	"errors"
	"log"
	"testing"

	"Qingyu_backend/config"
	aiModels "Qingyu_backend/models/ai"
	pb "Qingyu_backend/pkg/grpc/pb"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

func TestDetermineConsistencyAlertLevel(t *testing.T) {
	tests := []struct {
		name        string
		backend     int
		ai          int
		expected    aiModels.QuotaAlertLevel
		shouldAlert bool
	}{
		{
			name:        "small diff no alert",
			backend:     100,
			ai:          95,
			expected:    aiModels.QuotaAlertLevelInfo,
			shouldAlert: false,
		},
		{
			name:        "warning by ratio",
			backend:     500,
			ai:          430,
			expected:    aiModels.QuotaAlertLevelWarning,
			shouldAlert: true,
		},
		{
			name:        "critical by amount",
			backend:     2500,
			ai:          1000,
			expected:    aiModels.QuotaAlertLevelCritical,
			shouldAlert: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, shouldAlert := determineConsistencyAlertLevel(tt.backend, tt.ai)
			assert.Equal(t, tt.expected, level)
			assert.Equal(t, tt.shouldAlert, shouldAlert)
		})
	}
}

func TestDetermineConsistencyAlertLevelForScopeUsesConfigThresholds(t *testing.T) {
	previousConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{
		AIQuota: &config.AIQuotaConfig{
			ConsistencyThresholds: &config.QuotaConsistencyThresholdsConfig{
				Global: &config.QuotaConsistencyThresholdConfig{
					WarningTokens:  500,
					CriticalTokens: 2500,
					WarningRatio:   0.25,
					CriticalRatio:  0.5,
				},
			},
		},
	}
	t.Cleanup(func() {
		config.GlobalConfig = previousConfig
	})

	level, shouldAlert := determineConsistencyAlertLevelForScope(quotaConsistencyScopeGlobal, 1000, 700)
	assert.Equal(t, aiModels.QuotaAlertLevelWarning, level)
	assert.True(t, shouldAlert)

	level, shouldAlert = determineConsistencyAlertLevelForScope(quotaConsistencyScopeGlobal, 1000, 850)
	assert.Equal(t, aiModels.QuotaAlertLevelInfo, level)
	assert.False(t, shouldAlert)
}

func TestBuildQuotaConsumptionSummaryMap(t *testing.T) {
	summaryMap := buildQuotaConsumptionSummaryMap([]*pb.QuotaUserConsumptionSummary{
		{UserId: "user-1", TotalTokens: 120},
		nil,
		{UserId: "", TotalTokens: 20},
		{UserId: "user-2", TotalTokens: 45},
	})

	assert.Len(t, summaryMap, 2)
	assert.Equal(t, int32(120), summaryMap["user-1"].GetTotalTokens())
	assert.Equal(t, int32(45), summaryMap["user-2"].GetTotalTokens())
}

func TestBuildUserConsistencyAlertRequests(t *testing.T) {
	topConsumers := []aiModels.UserQuotaRanking{
		{UserID: "user-1", UsedQuota: 800},
		{UserID: "user-2", UsedQuota: 2500},
		{UserID: "user-3", UsedQuota: 90},
	}
	summaryMap := map[string]*pb.QuotaUserConsumptionSummary{
		"user-1": {UserId: "user-1", TotalTokens: 700},
		"user-2": {UserId: "user-2", TotalTokens: 1000},
		"user-3": {UserId: "user-3", TotalTokens: 86},
	}

	alerts := buildUserConsistencyAlertRequests(topConsumers, summaryMap, "day")

	assert.Len(t, alerts, 2)
	assert.Equal(t, "user-1", alerts[0].userID)
	assert.Equal(t, "warning", alerts[0].level)
	assert.Equal(t, "user", alerts[0].data["scope"])
	assert.Equal(t, "day", alerts[0].data["timeRange"])
	assert.Equal(t, "user-2", alerts[1].userID)
	assert.Equal(t, "critical", alerts[1].level)
}

func TestBuildReconciliationAlertRequests(t *testing.T) {
	summary := &aiModels.QuotaConsumptionReconciliationSummary{
		TimeRange:             "day",
		GroupBy:               "workflow",
		BackendTotalTokens:    4000,
		BackendTotalRecords:   40,
		AIServiceTotalTokens:  2500,
		AIServiceTotalRecords: 35,
		DifferenceTokens:      1500,
		DifferenceRatio:       0.375,
		AlertLevel:            "critical",
		ShouldAlert:           true,
		Items: []aiModels.QuotaConsumptionReconciliationItem{
			{
				GroupKey:         "rewrite",
				BackendTokens:    2600,
				BackendRecords:   21,
				AIServiceTokens:  1500,
				AIServiceRecords: 18,
				DifferenceTokens: 1100,
				DifferenceRatio:  0.42,
				AlertLevel:       "critical",
				ShouldAlert:      true,
			},
			{
				GroupKey:         "chat",
				BackendTokens:    1400,
				BackendRecords:   19,
				AIServiceTokens:  1385,
				AIServiceRecords: 17,
				DifferenceTokens: 15,
				DifferenceRatio:  0.01,
				AlertLevel:       "info",
				ShouldAlert:      false,
			},
		},
	}

	alerts := buildReconciliationAlertRequests(summary)

	assert.Len(t, alerts, 2)
	assert.Equal(t, "", alerts[0].userID)
	assert.Equal(t, "critical", alerts[0].level)
	assert.Equal(t, "global", alerts[0].data["scope"])
	assert.Equal(t, "AI 配额全局对账偏差", alerts[0].title)
	assert.Equal(t, "workflow", alerts[1].data["scope"])
	assert.Equal(t, "rewrite", alerts[1].data["groupKey"])
	assert.Equal(t, "AI 配额工作流对账偏差", alerts[1].title)
}

func TestEmitConsistencyAlertsResolvesRecoveredCheckedKeysWithoutNewAlerts(t *testing.T) {
	recoveredAlert := &aiModels.QuotaAlert{
		ID:      primitive.NewObjectID(),
		Type:    aiModels.QuotaAlertTypeConsistency,
		Level:   aiModels.QuotaAlertLevelCritical,
		Title:   "workflow mismatch",
		Status:  aiModels.QuotaAlertStatusPending,
		Message: "workflow mismatch",
		Data: map[string]interface{}{
			"scope":     quotaConsistencyScopeWorkflow,
			"timeRange": "day",
			"groupBy":   "workflow",
			"groupKey":  "chapter_expand",
		},
	}

	repo := newQuotaAlertRepoStub(recoveredAlert)
	scheduler := &QuotaScheduler{
		alertService: NewQuotaAlertService(repo),
		logger:       log.Default(),
	}

	checkedKeys := map[string]struct{}{
		buildConsistencyAlertKey("", recoveredAlert.Data): {},
	}

	err := scheduler.emitConsistencyAlerts(context.Background(), nil, checkedKeys)
	assert.NoError(t, err)
	assert.Equal(t, aiModels.QuotaAlertStatusResolved, recoveredAlert.Status)
	assert.Equal(t, "system-consistency-check", recoveredAlert.ResolvedBy)
	assert.NotNil(t, recoveredAlert.ResolvedAt)
}

func TestQuotaSchedulerConsistencyChecksHandleMissingDependenciesAndDownstreamErrors(t *testing.T) {
	t.Run("skips when dashboard service is missing", func(t *testing.T) {
		scheduler := &QuotaScheduler{logger: log.Default()}

		err := scheduler.checkCrossServiceConsistency(context.Background())
		assert.NoError(t, err)
	})

	t.Run("skips user-level checks when phase3 client is missing", func(t *testing.T) {
		scheduler := &QuotaScheduler{
			dashboardService: &QuotaDashboardService{
				quotaRepo: &quotaDashboardRepoStub{
					topConsumers: []aiModels.UserQuotaRanking{
						{UserID: "user-1", UsedQuota: 100},
					},
				},
			},
			logger: log.Default(),
		}

		err := scheduler.checkUserLevelConsistency(context.Background())
		assert.NoError(t, err)
	})

	t.Run("returns top consumer load error", func(t *testing.T) {
		scheduler := &QuotaScheduler{
			dashboardService: &QuotaDashboardService{
				quotaRepo: &quotaDashboardRepoStub{
					topConsumersErr: errors.New("dashboard unavailable"),
				},
			},
			phase3Client: &Phase3Client{},
			logger:       log.Default(),
		}

		err := scheduler.checkUserLevelConsistency(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "load top consumers failed")
		assert.Contains(t, err.Error(), "dashboard unavailable")
	})

	t.Run("skips aggregated checks when summary reader is missing", func(t *testing.T) {
		scheduler := &QuotaScheduler{
			dashboardService: &QuotaDashboardService{
				quotaRepo: &quotaDashboardRepoStub{},
			},
			logger: log.Default(),
		}

		err := scheduler.checkAggregatedConsistency(context.Background())
		assert.NoError(t, err)
	})
}

type quotaSchedulerAIClientConnStub struct {
	invokeCount int
	batchResp   *pb.QuotaConsumptionBatchResponse
	invokeErr   error
}

func (s *quotaSchedulerAIClientConnStub) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	s.invokeCount++
	if s.invokeErr != nil {
		return s.invokeErr
	}
	if resp, ok := reply.(*pb.QuotaConsumptionBatchResponse); ok && s.batchResp != nil {
		*resp = *s.batchResp
	}
	return nil
}

func (s *quotaSchedulerAIClientConnStub) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, errors.New("not implemented")
}

func TestQuotaSchedulerCheckUserLevelConsistencySkipsInvalidConsumersAndSurfacesBatchFailures(t *testing.T) {
	t.Run("skips when top consumers are invalid", func(t *testing.T) {
		conn := &quotaSchedulerAIClientConnStub{}
		scheduler := &QuotaScheduler{
			dashboardService: &QuotaDashboardService{
				quotaRepo: &quotaDashboardRepoStub{
					topConsumers: []aiModels.UserQuotaRanking{
						{UserID: "", UsedQuota: 100},
						{UserID: "user-zero", UsedQuota: 0},
					},
				},
			},
			phase3Client: &Phase3Client{client: pb.NewAIServiceClient(conn)},
			logger:       log.Default(),
		}

		err := scheduler.checkUserLevelConsistency(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 0, conn.invokeCount)
	})

	t.Run("surfaces unsuccessful ai batch response", func(t *testing.T) {
		conn := &quotaSchedulerAIClientConnStub{
			batchResp: &pb.QuotaConsumptionBatchResponse{
				Success:      false,
				ErrorMessage: "batch unavailable",
			},
		}
		scheduler := &QuotaScheduler{
			dashboardService: &QuotaDashboardService{
				quotaRepo: &quotaDashboardRepoStub{
					topConsumers: []aiModels.UserQuotaRanking{
						{UserID: "user-1", UsedQuota: 100},
					},
				},
			},
			phase3Client: &Phase3Client{client: pb.NewAIServiceClient(conn)},
			logger:       log.Default(),
		}

		err := scheduler.checkUserLevelConsistency(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "AI service batch quota query failed")
		assert.Contains(t, err.Error(), "batch unavailable")
		assert.Equal(t, 1, conn.invokeCount)
	})

	t.Run("surfaces ai batch rpc errors", func(t *testing.T) {
		conn := &quotaSchedulerAIClientConnStub{
			invokeErr: errors.New("rpc down"),
		}
		scheduler := &QuotaScheduler{
			dashboardService: &QuotaDashboardService{
				quotaRepo: &quotaDashboardRepoStub{
					topConsumers: []aiModels.UserQuotaRanking{
						{UserID: "user-1", UsedQuota: 100},
					},
				},
			},
			phase3Client: &Phase3Client{client: pb.NewAIServiceClient(conn)},
			logger:       log.Default(),
		}

		err := scheduler.checkUserLevelConsistency(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "load AI service quota summaries failed")
		assert.Contains(t, err.Error(), "rpc down")
		assert.Equal(t, 1, conn.invokeCount)
	})
}
