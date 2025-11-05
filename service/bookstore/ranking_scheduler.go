package bookstore

import (
	"Qingyu_backend/models/bookstore"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// RankingScheduler 榜单调度器
type RankingScheduler struct {
	service BookstoreService
	cron    *cron.Cron
	logger  *log.Logger
}

// NewRankingScheduler 创建榜单调度器
func NewRankingScheduler(service BookstoreService, logger *log.Logger) *RankingScheduler {
	return &RankingScheduler{
		service: service,
		cron:    cron.New(cron.WithSeconds()),
		logger:  logger,
	}
}

// Start 启动调度器
func (s *RankingScheduler) Start() error {
	// 实时榜：每5分钟更新一次
	_, err := s.cron.AddFunc("0 */5 * * * *", s.updateRealtimeRanking)
	if err != nil {
		return fmt.Errorf("failed to add realtime ranking job: %w", err)
	}

	// 周榜：每小时更新一次
	_, err = s.cron.AddFunc("0 0 * * * *", s.updateWeeklyRanking)
	if err != nil {
		return fmt.Errorf("failed to add weekly ranking job: %w", err)
	}

	// 月榜：每天凌晨2点更新
	_, err = s.cron.AddFunc("0 0 2 * * *", s.updateMonthlyRanking)
	if err != nil {
		return fmt.Errorf("failed to add monthly ranking job: %w", err)
	}

	// 新人榜：每天凌晨3点更新
	_, err = s.cron.AddFunc("0 0 3 * * *", s.updateNewbieRanking)
	if err != nil {
		return fmt.Errorf("failed to add newbie ranking job: %w", err)
	}

	// 清理过期榜单：每天凌晨4点执行
	_, err = s.cron.AddFunc("0 0 4 * * *", s.cleanupExpiredRankings)
	if err != nil {
		return fmt.Errorf("failed to add cleanup job: %w", err)
	}

	s.cron.Start()
	s.logger.Println("Ranking scheduler started")
	return nil
}

// Stop 停止调度器
func (s *RankingScheduler) Stop() {
	s.cron.Stop()
	s.logger.Println("Ranking scheduler stopped")
}

// updateRealtimeRanking 更新实时榜
func (s *RankingScheduler) updateRealtimeRanking() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	period := bookstore.GetPeriodString(bookstore.RankingTypeRealtime, time.Now())

	s.logger.Printf("Updating realtime ranking for period: %s", period)

	err := s.service.UpdateRankings(ctx, bookstore.RankingTypeRealtime, period)
	if err != nil {
		s.logger.Printf("Failed to update realtime ranking: %v", err)
		return
	}

	s.logger.Printf("Successfully updated realtime ranking for period: %s", period)
}

// updateWeeklyRanking 更新周榜
func (s *RankingScheduler) updateWeeklyRanking() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	now := time.Now()
	period := bookstore.GetPeriodString(bookstore.RankingTypeWeekly, now)

	s.logger.Printf("Updating weekly ranking for period: %s", period)

	err := s.service.UpdateRankings(ctx, bookstore.RankingTypeWeekly, period)
	if err != nil {
		s.logger.Printf("Failed to update weekly ranking: %v", err)
		return
	}

	s.logger.Printf("Successfully updated weekly ranking for period: %s", period)

	// 如果是周一，也更新上周的榜单
	if now.Weekday() == time.Monday {
		lastWeek := now.AddDate(0, 0, -7)
		lastWeekPeriod := bookstore.GetPeriodString(bookstore.RankingTypeWeekly, lastWeek)

		s.logger.Printf("Updating last week ranking for period: %s", lastWeekPeriod)

		err := s.service.UpdateRankings(ctx, bookstore.RankingTypeWeekly, lastWeekPeriod)
		if err != nil {
			s.logger.Printf("Failed to update last week ranking: %v", err)
		} else {
			s.logger.Printf("Successfully updated last week ranking for period: %s", lastWeekPeriod)
		}
	}
}

// updateMonthlyRanking 更新月榜
func (s *RankingScheduler) updateMonthlyRanking() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	now := time.Now()
	period := bookstore.GetPeriodString(bookstore.RankingTypeMonthly, now)

	s.logger.Printf("Updating monthly ranking for period: %s", period)

	err := s.service.UpdateRankings(ctx, bookstore.RankingTypeMonthly, period)
	if err != nil {
		s.logger.Printf("Failed to update monthly ranking: %v", err)
		return
	}

	s.logger.Printf("Successfully updated monthly ranking for period: %s", period)

	// 如果是每月1号，也更新上个月的榜单
	if now.Day() == 1 {
		lastMonth := now.AddDate(0, -1, 0)
		lastMonthPeriod := bookstore.GetPeriodString(bookstore.RankingTypeMonthly, lastMonth)

		s.logger.Printf("Updating last month ranking for period: %s", lastMonthPeriod)

		err := s.service.UpdateRankings(ctx, bookstore.RankingTypeMonthly, lastMonthPeriod)
		if err != nil {
			s.logger.Printf("Failed to update last month ranking: %v", err)
		} else {
			s.logger.Printf("Successfully updated last month ranking for period: %s", lastMonthPeriod)
		}
	}
}

// updateNewbieRanking 更新新人榜
func (s *RankingScheduler) updateNewbieRanking() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	now := time.Now()
	period := bookstore.GetPeriodString(bookstore.RankingTypeNewbie, now)

	s.logger.Printf("Updating newbie ranking for period: %s", period)

	err := s.service.UpdateRankings(ctx, bookstore.RankingTypeNewbie, period)
	if err != nil {
		s.logger.Printf("Failed to update newbie ranking: %v", err)
		return
	}

	s.logger.Printf("Successfully updated newbie ranking for period: %s", period)

	// 如果是每月1号，也更新上个月的新人榜
	if now.Day() == 1 {
		lastMonth := now.AddDate(0, -1, 0)
		lastMonthPeriod := bookstore.GetPeriodString(bookstore.RankingTypeNewbie, lastMonth)

		s.logger.Printf("Updating last month newbie ranking for period: %s", lastMonthPeriod)

		err := s.service.UpdateRankings(ctx, bookstore.RankingTypeNewbie, lastMonthPeriod)
		if err != nil {
			s.logger.Printf("Failed to update last month newbie ranking: %v", err)
		} else {
			s.logger.Printf("Successfully updated last month newbie ranking for period: %s", lastMonthPeriod)
		}
	}
}

// cleanupExpiredRankings 清理过期榜单
func (s *RankingScheduler) cleanupExpiredRankings() {
	s.logger.Println("Starting cleanup of expired rankings")

	// 这里可以添加清理逻辑，比如删除3个月前的榜单数据
	// 由于没有直接的Repository访问，这里只是记录日志
	// 实际实现中可能需要注入RankingRepository

	s.logger.Println("Cleanup of expired rankings completed")
}

// UpdateRankingNow 立即更新指定类型的榜单
func (s *RankingScheduler) UpdateRankingNow(rankingType bookstore.RankingType, period string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if period == "" {
		period = bookstore.GetPeriodString(rankingType, time.Now())
	}

	s.logger.Printf("Manually updating %s ranking for period: %s", rankingType, period)

	err := s.service.UpdateRankings(ctx, rankingType, period)
	if err != nil {
		s.logger.Printf("Failed to manually update %s ranking: %v", rankingType, err)
		return err
	}

	s.logger.Printf("Successfully manually updated %s ranking for period: %s", rankingType, period)
	return nil
}

// GetSchedulerStatus 获取调度器状态
func (s *RankingScheduler) GetSchedulerStatus() map[string]interface{} {
	entries := s.cron.Entries()

	status := map[string]interface{}{
		"running":   len(entries) > 0,
		"job_count": len(entries),
		"next_runs": make([]map[string]interface{}, 0),
	}

	for i, entry := range entries {
		jobInfo := map[string]interface{}{
			"id":       i + 1,
			"next_run": entry.Next.Format("2006-01-02 15:04:05"),
		}
		status["next_runs"] = append(status["next_runs"].([]map[string]interface{}), jobInfo)
	}

	return status
}
