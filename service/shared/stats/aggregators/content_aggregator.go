package aggregators

import (
	"context"

	"Qingyu_backend/service/shared/stats"
)

// ContentAggregator 内容维度聚合器（P1骨架实现）
type ContentAggregator struct{}

var _ stats.ContentAggregatorPort = (*ContentAggregator)(nil)

// NewContentAggregator 创建内容聚合器
func NewContentAggregator() *ContentAggregator {
	return &ContentAggregator{}
}

// AggregateContentStats 聚合平台内容统计（骨架实现）
func (a *ContentAggregator) AggregateContentStats(ctx context.Context, filter *stats.StatsFilter) (*stats.PlatformContentStats, error) {
	_ = ctx
	_ = filter

	return &stats.PlatformContentStats{
		TotalBooks:        0,
		NewBooks:          0,
		TotalChapters:     0,
		TotalWords:        0,
		TotalViews:        0,
		AverageRating:     0,
		PopularCategories: []string{},
	}, nil
}
