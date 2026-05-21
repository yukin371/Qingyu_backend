package ai

import (
	"context"
	"io"
	"log"
	"testing"

	"Qingyu_backend/config"
	aiModels "Qingyu_backend/models/ai"
	pb "Qingyu_backend/pkg/grpc/pb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuotaSchedulerRunConsistencyCheckCreatesGlobalAlertFromAggregatedSummary(t *testing.T) {
	previousConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{
		AIQuota: &config.AIQuotaConfig{
			ConsistencyThresholds: &config.QuotaConsistencyThresholdsConfig{
				Global: &config.QuotaConsistencyThresholdConfig{
					WarningTokens:  50,
					CriticalTokens: 150,
					WarningRatio:   0.2,
					CriticalRatio:  0.6,
				},
				Workflow: &config.QuotaConsistencyThresholdConfig{
					WarningTokens:  1000,
					CriticalTokens: 2000,
					WarningRatio:   1.0,
					CriticalRatio:  1.0,
				},
			},
		},
	}
	t.Cleanup(func() {
		config.GlobalConfig = previousConfig
	})

	backendItems := make([]aiModels.QuotaConsumptionSummaryItem, 0, 20)
	aiItems := make([]*pb.QuotaConsumptionSummaryItem, 0, 20)
	for i := 1; i <= 20; i++ {
		groupKey := "workflow-" + string(rune('a'+i-1))
		backendItems = append(backendItems, aiModels.QuotaConsumptionSummaryItem{
			GroupKey:     groupKey,
			TotalTokens:  10,
			TotalRecords: 1,
		})
		aiItems = append(aiItems, &pb.QuotaConsumptionSummaryItem{
			GroupKey:     groupKey,
			TotalTokens:  5,
			TotalRecords: 1,
		})
	}

	repo := &quotaDashboardRepoStub{
		consumption: &aiModels.QuotaConsumptionSummary{
			GroupBy:      "workflow",
			TotalGroups:  20,
			TotalTokens:  200,
			TotalRecords: 20,
			Items:        backendItems,
		},
	}
	reader := &quotaSummaryReaderStub{
		response: &pb.QuotaConsumptionSummaryResponse{
			Success:      true,
			TimeRange:    "day",
			GroupBy:      "workflow",
			Page:         1,
			PageSize:     20,
			TotalGroups:  20,
			TotalTokens:  100,
			TotalRecords: 20,
			Items:        aiItems,
		},
	}
	alertRepo := newQuotaAlertRepoStub()
	dashboardService := NewQuotaDashboardService(repo, alertRepo, nil)
	dashboardService.SetConsumptionSummaryReader(reader)

	scheduler := &QuotaScheduler{
		dashboardService: dashboardService,
		alertService:     NewQuotaAlertService(alertRepo),
		logger:           log.New(io.Discard, "", 0),
	}

	err := scheduler.RunConsistencyCheck(context.Background())
	require.NoError(t, err)
	require.Len(t, alertRepo.alerts, 1)

	for _, alert := range alertRepo.alerts {
		assert.Equal(t, aiModels.QuotaAlertTypeConsistency, alert.Type)
		assert.Equal(t, aiModels.QuotaAlertLevelWarning, alert.Level)
		assert.Equal(t, aiModels.QuotaAlertStatusPending, alert.Status)
		assert.Equal(t, "global", alert.Data["scope"])
		assert.Equal(t, "workflow", alert.Data["groupBy"])
	}
}
