package interfaces

import (
	"Qingyu_backend/models/reading/bookstore"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RankingRepository 榜单仓储接口
type RankingRepository interface {
	// 继承基础CRUD接口
	base.CRUDRepository[*bookstore.RankingItem, primitive.ObjectID]

	// 榜单特定查询方法
	GetByType(ctx context.Context, rankingType bookstore.RankingType, period string, limit, offset int) ([]*bookstore.RankingItem, error)
	GetByTypeWithBooks(ctx context.Context, rankingType bookstore.RankingType, period string, limit, offset int) ([]*bookstore.RankingItem, error)
	GetByBookID(ctx context.Context, bookID primitive.ObjectID, rankingType bookstore.RankingType, period string) (*bookstore.RankingItem, error)
	GetByPeriod(ctx context.Context, period string, limit, offset int) ([]*bookstore.RankingItem, error)

	// 榜单统计方法
	GetRankingStats(ctx context.Context, rankingType bookstore.RankingType, period string) (*bookstore.RankingStats, error)
	CountByType(ctx context.Context, rankingType bookstore.RankingType, period string) (int64, error)
	GetTopBooks(ctx context.Context, rankingType bookstore.RankingType, period string, limit int) ([]*bookstore.RankingItem, error)

	// 榜单更新方法
	UpsertRankingItem(ctx context.Context, item *bookstore.RankingItem) error
	BatchUpsertRankingItems(ctx context.Context, items []*bookstore.RankingItem) error
	UpdateRankings(ctx context.Context, rankingType bookstore.RankingType, period string, items []*bookstore.RankingItem) error

	// 榜单维护方法
	DeleteByPeriod(ctx context.Context, period string) error
	DeleteByType(ctx context.Context, rankingType bookstore.RankingType) error
	DeleteExpiredRankings(ctx context.Context, beforeDate time.Time) error

	// 实时榜单计算
	CalculateRealtimeRanking(ctx context.Context, period string) ([]*bookstore.RankingItem, error)
	CalculateWeeklyRanking(ctx context.Context, period string) ([]*bookstore.RankingItem, error)
	CalculateMonthlyRanking(ctx context.Context, period string) ([]*bookstore.RankingItem, error)
	CalculateNewbieRanking(ctx context.Context, period string) ([]*bookstore.RankingItem, error)

	// 事务支持
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
