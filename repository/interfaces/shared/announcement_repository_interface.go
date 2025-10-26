package shared

import (
	"Qingyu_backend/models/shared"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AnnouncementRepository 公告仓储接口
type AnnouncementRepository interface {
	// 继承基础CRUD接口
	base.CRUDRepository[*shared.Announcement, primitive.ObjectID]

	// 公告特定查询方法
	GetActive(ctx context.Context, targetUsers string, limit, offset int) ([]*shared.Announcement, error)
	GetByType(ctx context.Context, announcementType string, limit, offset int) ([]*shared.Announcement, error)
	GetByTimeRange(ctx context.Context, startTime, endTime *time.Time, limit, offset int) ([]*shared.Announcement, error)
	GetEffective(ctx context.Context, targetUsers string, limit int) ([]*shared.Announcement, error)

	// 统计方法
	IncrementViewCount(ctx context.Context, announcementID primitive.ObjectID) error
	GetViewStats(ctx context.Context, announcementID primitive.ObjectID) (int64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, announcementIDs []primitive.ObjectID, isActive bool) error
	BatchDelete(ctx context.Context, announcementIDs []primitive.ObjectID) error

	// 事务支持
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
