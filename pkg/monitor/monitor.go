package monitor

import (
	"context"
	"fmt"
	"time"
)

// DataQualityMonitor 数据质量监控器
type DataQualityMonitor struct {
	db     Database
	alert  Alerter
	logger Logger
}

// Database 数据库接口
type Database interface {
	// OrphanedRecords 检查孤儿记录数量
	OrphanedRecords(ctx context.Context, collection, foreignKey, targetCollection string) (int, error)
	// StatisticsAccuracy 检查统计数据不准确的数量
	StatisticsAccuracy(ctx context.Context, collection, countField string) (int, error)
}

// Alerter 告警接口
type Alerter interface {
	// SendAlert 发送告警
	SendAlert(ctx context.Context, message string) error
}

// Logger 日志接口
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// CheckResult 检查结果
type CheckResult struct {
	Collection string `json:"collection"`
	Field      string `json:"field"`
	Count      int    `json:"count"`
}

// DataQualityReport 数据质量报告
type DataQualityReport struct {
	CheckTime            time.Time     `json:"check_time"`
	TotalOrphanedRecords int           `json:"total_orphaned_records"`
	InaccurateStatistics int           `json:"inaccurate_statistics"`
	Details              []CheckResult `json:"details"`
	HasIssues            bool          `json:"has_issues"`
}

// ToJSON 将报告转换为JSON字符串
func (r *DataQualityReport) ToJSON() string {
	return fmt.Sprintf(
		`{"check_time":"%s","total_orphaned_records":%d,"inaccurate_statistics":%d,"has_issues":%t,"details":%d}`,
		r.CheckTime.Format(time.RFC3339),
		r.TotalOrphanedRecords,
		r.InaccurateStatistics,
		r.HasIssues,
		len(r.Details),
	)
}

// NewDataQualityMonitor 创建数据质量监控器
func NewDataQualityMonitor(db Database, alert Alerter) *DataQualityMonitor {
	return &DataQualityMonitor{
		db:    db,
		alert: alert,
	}
}

// CheckOrphanedRecords 检查孤儿记录
func (m *DataQualityMonitor) CheckOrphanedRecords(ctx context.Context) (*DataQualityReport, error) {
	report := &DataQualityReport{
		CheckTime: time.Now(),
		Details:   make([]CheckResult, 0),
	}

	// 定义需要检查的外键关系
	checks := []struct {
		collection        string
		foreignKey        string
		targetCollection  string
	}{
		// 阅读进度
		{"reading_progress", "user_id", "users"},
		{"reading_progress", "book_id", "books"},
		// 阅读历史
		{"reading_history", "user_id", "users"},
		{"reading_history", "book_id", "books"},
		{"reading_history", "chapter_id", "chapters"},
		// 书签
		{"bookmarks", "user_id", "users"},
		{"bookmarks", "book_id", "books"},
		{"bookmarks", "chapter_id", "chapters"},
		// 点赞
		{"likes", "user_id", "users"},
		{"likes", "target_id", "users"}, // 简化处理，实际需要根据target_type判断
		// 通知
		{"notifications", "user_id", "users"},
		// 作者收益
		{"author_revenue", "user_id", "users"},
	}

	for _, check := range checks {
		count, err := m.db.OrphanedRecords(ctx, check.collection, check.foreignKey, check.targetCollection)
		if err != nil {
			return nil, fmt.Errorf("检查 %s.%s 失败: %w", check.collection, check.foreignKey, err)
		}

		if count > 0 {
			report.TotalOrphanedRecords += count
			report.Details = append(report.Details, CheckResult{
				Collection: check.collection,
				Field:      check.foreignKey,
				Count:      count,
			})
		}
	}

	report.HasIssues = report.TotalOrphanedRecords > 0

	return report, nil
}

// CheckStatisticsAccuracy 检查统计数据准确性
func (m *DataQualityMonitor) CheckStatisticsAccuracy(ctx context.Context) (*DataQualityReport, error) {
	report := &DataQualityReport{
		CheckTime: time.Now(),
		Details:   make([]CheckResult, 0),
	}

	// 检查 books.likes_count
	count, err := m.db.StatisticsAccuracy(ctx, "books", "likes_count")
	if err != nil {
		return nil, err
	}
	if count > 0 {
		report.InaccurateStatistics += count
		report.Details = append(report.Details, CheckResult{
			Collection: "books",
			Field:      "likes_count",
			Count:      count,
		})
	}

	// 检查 books.comments_count
	count, err = m.db.StatisticsAccuracy(ctx, "books", "comments_count")
	if err != nil {
		return nil, err
	}
	if count > 0 {
		report.InaccurateStatistics += count
		report.Details = append(report.Details, CheckResult{
			Collection: "books",
			Field:      "comments_count",
			Count:      count,
		})
	}

	// 检查 users.followers_count
	count, err = m.db.StatisticsAccuracy(ctx, "users", "followers_count")
	if err != nil {
		return nil, err
	}
	if count > 0 {
		report.InaccurateStatistics += count
		report.Details = append(report.Details, CheckResult{
			Collection: "users",
			Field:      "followers_count",
			Count:      count,
		})
	}

	report.HasIssues = report.InaccurateStatistics > 0

	return report, nil
}

// RunDailyCheck 执行每日检查
func (m *DataQualityMonitor) RunDailyCheck(ctx context.Context) (*DataQualityReport, error) {
	// 1. 检查孤儿记录
	orphanReport, err := m.CheckOrphanedRecords(ctx)
	if err != nil {
		return nil, err
	}

	// 2. 检查统计数据
	statsReport, err := m.CheckStatisticsAccuracy(ctx)
	if err != nil {
		return nil, err
	}

	// 3. 合并报告
	finalReport := &DataQualityReport{
		CheckTime:            time.Now(),
		TotalOrphanedRecords: orphanReport.TotalOrphanedRecords,
		InaccurateStatistics: statsReport.InaccurateStatistics,
		Details:              append(orphanReport.Details, statsReport.Details...),
		HasIssues:            orphanReport.HasIssues || statsReport.HasIssues,
	}

	// 4. 生成报告
	m.logReport(finalReport)

	// 5. 发送告警
	if m.ShouldAlert(finalReport) {
		if m.alert != nil {
			if err := m.alert.SendAlert(ctx, m.formatAlertMessage(finalReport)); err != nil {
				m.logError("发送告警失败", err)
			}
		}
	}

	return finalReport, nil
}

// ShouldAlert 判断是否需要告警
func (m *DataQualityMonitor) ShouldAlert(report *DataQualityReport) bool {
	return report.HasIssues
}

// formatAlertMessage 格式化告警消息
func (m *DataQualityMonitor) formatAlertMessage(report *DataQualityReport) string {
	return fmt.Sprintf("数据质量问题告警: 孤儿记录=%d, 统计不准确=%d",
		report.TotalOrphanedRecords, report.InaccurateStatistics)
}

// logReport 记录报告
func (m *DataQualityMonitor) logReport(report *DataQualityReport) {
	if m.logger != nil {
		m.logger.Info("数据质量检查完成",
			"孤儿记录数", report.TotalOrphanedRecords,
			"统计不准确数", report.InaccurateStatistics,
			"是否有问题", report.HasIssues,
		)
	}
}

// logError 记录错误
func (m *DataQualityMonitor) logError(msg string, err error) {
	if m.logger != nil {
		m.logger.Error(msg, "error", err)
	}
}
