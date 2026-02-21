package aggregators

import (
	"context"
	"time"

	"Qingyu_backend/service/shared/stats"
)

// UserAggregator 用户维度聚合器（P1骨架实现）
type UserAggregator struct{}

var _ stats.UserAggregatorPort = (*UserAggregator)(nil)

// NewUserAggregator 创建用户聚合器
func NewUserAggregator() *UserAggregator {
	return &UserAggregator{}
}

// AggregateUserStats 聚合平台用户统计（骨架实现）
func (a *UserAggregator) AggregateUserStats(ctx context.Context, filter *stats.StatsFilter) (*stats.PlatformUserStats, error) {
	if filter == nil {
		filter = &stats.StatsFilter{}
	}
	_ = ctx

	return &stats.PlatformUserStats{
		TotalUsers:       0,
		NewUsers:         0,
		ActiveUsers:      0,
		VIPUsers:         0,
		RetentionRate:    0,
		AverageActiveDay: 0,
	}, nil
}

// BuildRange 归一化统计时间范围
func (a *UserAggregator) BuildRange(filter *stats.StatsFilter) (time.Time, time.Time) {
	now := time.Now()
	if filter == nil || filter.StartDate == nil || filter.EndDate == nil {
		return now.AddDate(0, 0, -7), now
	}
	return *filter.StartDate, *filter.EndDate
}
