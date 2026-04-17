package ai

import (
	"context"
	"fmt"
	"log"
	"time"

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

	if err := s.checkCrossServiceConsistency(ctx); err != nil {
		s.logger.Printf("quota consistency check failed: %v", err)
		return
	}

	s.logger.Println("quota consistency check completed")
}

func (s *QuotaScheduler) checkCrossServiceConsistency(ctx context.Context) error {
	if s.phase3Client == nil {
		s.logger.Println("skip quota consistency check: Phase3Client not configured")
		return nil
	}

	// TODO: 当前 Phase3Client 仅暴露创作工作流和健康检查接口，尚未提供跨服务 quota 统计查询。
	// 确认路径：为 AI gRPC 服务补充 quota/report 或 transactions 汇总接口，然后在此比对 MongoDB 聚合结果。
	_ = ctx
	_ = s.alertService
	s.logger.Println("skip quota consistency check: Phase3Client has no quota statistics API yet")
	return nil
}
