package ai

import (
	"context"
	"fmt"
	"log"
	"time"

	aiModels "Qingyu_backend/models/ai"
	pb "Qingyu_backend/pkg/grpc/pb"

	"github.com/robfig/cron/v3"
)

// QuotaScheduler 配额相关定时任务调度器。
type QuotaScheduler struct {
	anomalyDetector  *QuotaAnomalyDetector
	dashboardService *QuotaDashboardService
	alertService     *QuotaAlertService
	phase3Client     *Phase3Client
	cron             *cron.Cron
	logger           *log.Logger
}

// NewQuotaScheduler 创建配额调度器。
func NewQuotaScheduler(
	anomalyDetector *QuotaAnomalyDetector,
	dashboardService *QuotaDashboardService,
	alertService *QuotaAlertService,
	phase3Client *Phase3Client,
	logger *log.Logger,
) *QuotaScheduler {
	if logger == nil {
		logger = log.Default()
	}

	return &QuotaScheduler{
		anomalyDetector:  anomalyDetector,
		dashboardService: dashboardService,
		alertService:     alertService,
		phase3Client:     phase3Client,
		cron:             cron.New(cron.WithSeconds()),
		logger:           logger,
	}
}

// Start 启动调度器。
func (s *QuotaScheduler) Start() error {
	if s.anomalyDetector != nil {
		if _, err := s.cron.AddFunc("0 */5 * * * *", s.runAnomalyDetection); err != nil {
			return fmt.Errorf("failed to add anomaly detection job: %w", err)
		}
	}

	if s.dashboardService != nil {
		if _, err := s.cron.AddFunc("0 */5 * * * *", s.refreshDashboardCache); err != nil {
			return fmt.Errorf("failed to add dashboard refresh job: %w", err)
		}
	}

	if _, err := s.cron.AddFunc("0 0 2 * * *", s.runConsistencyCheck); err != nil {
		return fmt.Errorf("failed to add consistency check job: %w", err)
	}

	s.cron.Start()
	s.logger.Println("quota scheduler started")
	return nil
}

// Stop 停止调度器。
func (s *QuotaScheduler) Stop() {
	s.cron.Stop()
	s.logger.Println("quota scheduler stopped")
}

func (s *QuotaScheduler) runAnomalyDetection() {
	if s.anomalyDetector == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := s.anomalyDetector.DetectAnomalies(ctx); err != nil {
		s.logger.Printf("quota anomaly detection failed: %v", err)
		return
	}

	s.logger.Println("quota anomaly detection completed")
}

func (s *QuotaScheduler) refreshDashboardCache() {
	if s.dashboardService == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := s.dashboardService.RefreshDashboardCache(ctx); err != nil {
		s.logger.Printf("quota dashboard refresh failed: %v", err)
		return
	}

	s.logger.Println("quota dashboard cache refreshed")
}

func (s *QuotaScheduler) runConsistencyCheck() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.RunConsistencyCheck(ctx); err != nil {
		s.logger.Printf("quota consistency check failed: %v", err)
		return
	}

	s.logger.Println("quota consistency check completed")
}

// RunConsistencyCheck 手动执行一次跨服务对账检查。
func (s *QuotaScheduler) RunConsistencyCheck(ctx context.Context) error {
	return s.checkCrossServiceConsistency(ctx)
}

func (s *QuotaScheduler) checkCrossServiceConsistency(ctx context.Context) error {
	if s.dashboardService == nil {
		s.logger.Println("skip quota consistency check: dashboard service not configured")
		return nil
	}

	if err := s.checkUserLevelConsistency(ctx); err != nil {
		return err
	}

	if err := s.checkAggregatedConsistency(ctx); err != nil {
		return err
	}

	return nil
}

func (s *QuotaScheduler) checkUserLevelConsistency(ctx context.Context) error {
	if s.phase3Client == nil {
		s.logger.Println("skip user-level quota consistency check: Phase3Client not configured")
		return nil
	}

	topConsumers, err := s.dashboardService.GetTopConsumers(ctx, 10)
	if err != nil {
		return fmt.Errorf("load top consumers failed: %w", err)
	}
	if len(topConsumers) == 0 {
		s.logger.Println("skip user-level quota consistency check: no quota consumers found")
		return nil
	}

	userIDs := make([]string, 0, len(topConsumers))
	for _, consumer := range topConsumers {
		if consumer.UserID != "" && consumer.UsedQuota > 0 {
			userIDs = append(userIDs, consumer.UserID)
		}
	}
	if len(userIDs) == 0 {
		s.logger.Println("skip user-level quota consistency check: no valid quota consumers found")
		return nil
	}

	batchResp, err := s.phase3Client.GetQuotaConsumptionBatch(ctx, userIDs, "day", "")
	if err != nil {
		return fmt.Errorf("load AI service quota summaries failed: %w", err)
	}
	if !batchResp.GetSuccess() {
		if batchResp.GetErrorMessage() != "" {
			return fmt.Errorf("AI service batch quota query failed: %s", batchResp.GetErrorMessage())
		}
		return fmt.Errorf("AI service batch quota query returned unsuccessful response")
	}

	summaryByUserID := buildQuotaConsumptionSummaryMap(batchResp.GetSummaries())
	alerts := buildUserConsistencyAlertRequests(topConsumers, summaryByUserID, "day")
	checkedKeys := buildCheckedUserConsistencyAlertKeys(topConsumers, "day")
	return s.emitConsistencyAlerts(ctx, alerts, checkedKeys)
}

func (s *QuotaScheduler) checkAggregatedConsistency(ctx context.Context) error {
	summary, err := s.dashboardService.GetReconciliationSummary(ctx, "day", "", "workflow", 1, 20)
	if err != nil {
		s.logger.Printf("skip aggregated quota consistency check: %v", err)
		return nil
	}

	alerts := buildReconciliationAlertRequests(summary)
	checkedKeys := buildCheckedReconciliationAlertKeys(summary)
	return s.emitConsistencyAlerts(ctx, alerts, checkedKeys)
}

type quotaConsistencyAlertRequest struct {
	userID  string
	level   string
	title   string
	message string
	data    map[string]interface{}
}

func buildUserConsistencyAlertRequests(
	topConsumers []aiModels.UserQuotaRanking,
	summaryByUserID map[string]*pb.QuotaUserConsumptionSummary,
	timeRange string,
) []quotaConsistencyAlertRequest {
	alerts := make([]quotaConsistencyAlertRequest, 0)

	for _, consumer := range topConsumers {
		if consumer.UserID == "" || consumer.UsedQuota <= 0 {
			continue
		}

		summary := summaryByUserID[consumer.UserID]
		aiTokens := 0
		if summary != nil {
			aiTokens = int(summary.GetTotalTokens())
		}

		level, shouldAlert := determineConsistencyAlertLevel(consumer.UsedQuota, aiTokens)
		if !shouldAlert {
			continue
		}

		diff := absInt(consumer.UsedQuota - aiTokens)
		alerts = append(alerts, quotaConsistencyAlertRequest{
			userID: consumer.UserID,
			level:  string(level),
			title:  "AI 配额对账偏差",
			message: fmt.Sprintf(
				"用户 %s 的 AI 配额对账出现偏差：后端=%d，AI服务=%d，差值=%d",
				consumer.UserID,
				consumer.UsedQuota,
				aiTokens,
				diff,
			),
			data: map[string]interface{}{
				"scope":            "user",
				"backendUsedQuota": consumer.UsedQuota,
				"aiServiceTokens":  aiTokens,
				"diff":             diff,
				"timeRange":        timeRange,
			},
		})
	}

	return alerts
}

func buildReconciliationAlertRequests(
	summary *aiModels.QuotaConsumptionReconciliationSummary,
) []quotaConsistencyAlertRequest {
	if summary == nil {
		return nil
	}

	alerts := make([]quotaConsistencyAlertRequest, 0)
	if summary.ShouldAlert {
		alerts = append(alerts, quotaConsistencyAlertRequest{
			level: string(summary.AlertLevel),
			title: "AI 配额全局对账偏差",
			message: fmt.Sprintf(
				"AI 配额全局对账出现偏差：后端=%d，AI服务=%d，差值=%d",
				summary.BackendTotalTokens,
				summary.AIServiceTotalTokens,
				summary.DifferenceTokens,
			),
			data: map[string]interface{}{
				"scope":                 "global",
				"timeRange":             summary.TimeRange,
				"groupBy":               summary.GroupBy,
				"backendTotalTokens":    summary.BackendTotalTokens,
				"backendTotalRecords":   summary.BackendTotalRecords,
				"aiServiceTotalTokens":  summary.AIServiceTotalTokens,
				"aiServiceTotalRecords": summary.AIServiceTotalRecords,
				"diff":                  summary.DifferenceTokens,
				"diffRatio":             summary.DifferenceRatio,
			},
		})
	}

	for _, item := range summary.Items {
		if !item.ShouldAlert {
			continue
		}
		alerts = append(alerts, quotaConsistencyAlertRequest{
			level: string(item.AlertLevel),
			title: "AI 配额工作流对账偏差",
			message: fmt.Sprintf(
				"工作流 %s 的 AI 配额对账出现偏差：后端=%d，AI服务=%d，差值=%d",
				item.GroupKey,
				item.BackendTokens,
				item.AIServiceTokens,
				item.DifferenceTokens,
			),
			data: map[string]interface{}{
				"scope":            "workflow",
				"timeRange":        summary.TimeRange,
				"groupBy":          summary.GroupBy,
				"groupKey":         item.GroupKey,
				"backendTokens":    item.BackendTokens,
				"backendRecords":   item.BackendRecords,
				"aiServiceTokens":  item.AIServiceTokens,
				"aiServiceRecords": item.AIServiceRecords,
				"diff":             item.DifferenceTokens,
				"diffRatio":        item.DifferenceRatio,
			},
		})
	}

	return alerts
}

func (s *QuotaScheduler) emitConsistencyAlerts(
	ctx context.Context,
	alerts []quotaConsistencyAlertRequest,
	checkedKeys map[string]struct{},
) error {
	activeKeys := buildActiveConsistencyAlertKeys(alerts)

	if s.alertService != nil {
		if err := s.alertService.ResolveRecoveredConsistencyAlerts(
			ctx,
			checkedKeys,
			activeKeys,
			"system-consistency-check",
		); err != nil {
			return err
		}
	}

	if len(alerts) == 0 {
		return nil
	}

	for _, alert := range alerts {
		if s.alertService == nil {
			s.logger.Printf(
				"quota consistency mismatch detected without alert service: scope=%v user=%s message=%s",
				alert.data["scope"],
				alert.userID,
				alert.message,
			)
			continue
		}

		if err := s.alertService.CreateConsistencyAlert(
			ctx,
			alert.userID,
			alert.level,
			alert.title,
			alert.message,
			alert.data,
		); err != nil {
			s.logger.Printf(
				"create quota consistency alert failed: scope=%v user=%s err=%v",
				alert.data["scope"],
				alert.userID,
				err,
			)
		}
	}

	return nil
}

func buildCheckedUserConsistencyAlertKeys(
	topConsumers []aiModels.UserQuotaRanking,
	timeRange string,
) map[string]struct{} {
	result := make(map[string]struct{})
	for _, consumer := range topConsumers {
		if consumer.UserID == "" || consumer.UsedQuota <= 0 {
			continue
		}
		key := buildConsistencyAlertKey(consumer.UserID, map[string]interface{}{
			"scope":     quotaConsistencyScopeUser,
			"timeRange": timeRange,
		})
		result[key] = struct{}{}
	}
	return result
}

func buildCheckedReconciliationAlertKeys(
	summary *aiModels.QuotaConsumptionReconciliationSummary,
) map[string]struct{} {
	result := make(map[string]struct{})
	if summary == nil {
		return result
	}

	globalKey := buildConsistencyAlertKey("", map[string]interface{}{
		"scope":     quotaConsistencyScopeGlobal,
		"timeRange": summary.TimeRange,
		"groupBy":   summary.GroupBy,
	})
	result[globalKey] = struct{}{}

	for _, item := range summary.Items {
		workflowKey := buildConsistencyAlertKey("", map[string]interface{}{
			"scope":     quotaConsistencyScopeWorkflow,
			"timeRange": summary.TimeRange,
			"groupBy":   summary.GroupBy,
			"groupKey":  item.GroupKey,
		})
		result[workflowKey] = struct{}{}
	}

	return result
}

func buildActiveConsistencyAlertKeys(alerts []quotaConsistencyAlertRequest) map[string]struct{} {
	result := make(map[string]struct{}, len(alerts))
	for _, alert := range alerts {
		key := buildConsistencyAlertKey(alert.userID, alert.data)
		result[key] = struct{}{}
	}
	return result
}

func buildQuotaConsumptionSummaryMap(summaries []*pb.QuotaUserConsumptionSummary) map[string]*pb.QuotaUserConsumptionSummary {
	result := make(map[string]*pb.QuotaUserConsumptionSummary, len(summaries))
	for _, summary := range summaries {
		if summary == nil || summary.GetUserId() == "" {
			continue
		}
		result[summary.GetUserId()] = summary
	}
	return result
}
