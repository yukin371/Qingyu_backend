package stats

import (
	"context"
	"time"
)

// StatsFilter 统计查询过滤条件
type StatsFilter struct {
	UserID    string     `json:"user_id,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Limit     int        `json:"limit,omitempty"`
}

// StatsPort 统计服务端口定义
// 用于上层依赖反转，屏蔽具体实现细节
type StatsPort interface {
	GetUserStats(ctx context.Context, userID string) (*UserStats, error)
	GetPlatformUserStats(ctx context.Context, startDate, endDate time.Time) (*PlatformUserStats, error)
	GetContentStats(ctx context.Context, userID string) (*ContentStats, error)
	GetPlatformContentStats(ctx context.Context, startDate, endDate time.Time) (*PlatformContentStats, error)
	GetUserActivityStats(ctx context.Context, userID string, days int) (*ActivityStats, error)
	GetRevenueStats(ctx context.Context, userID string, startDate, endDate time.Time) (*RevenueStats, error)
	Health(ctx context.Context) error
}

// UserAggregatorPort 用户维度聚合器端口
type UserAggregatorPort interface {
	AggregateUserStats(ctx context.Context, filter *StatsFilter) (*PlatformUserStats, error)
}

// ContentAggregatorPort 内容维度聚合器端口
type ContentAggregatorPort interface {
	AggregateContentStats(ctx context.Context, filter *StatsFilter) (*PlatformContentStats, error)
}

// AggregatorPort 兼容组合接口（建议新代码优先依赖细粒度端口）
type AggregatorPort interface {
	UserAggregatorPort
	ContentAggregatorPort
}
