package writing

import (
	"Qingyu_backend/models/writer"
	"context"
)

// LocationRepository 地点Repository接口
type LocationRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, location *writer.Location) error
	FindByID(ctx context.Context, locationID string) (*writer.Location, error)
	FindByProjectID(ctx context.Context, projectID string) ([]*writer.Location, error)
	FindByParentID(ctx context.Context, parentID string) ([]*writer.Location, error)
	Update(ctx context.Context, location *writer.Location) error
	Delete(ctx context.Context, locationID string) error

	// 关系管理
	CreateRelation(ctx context.Context, relation *writer.LocationRelation) error
	FindRelations(ctx context.Context, projectID string, locationID *string) ([]*writer.LocationRelation, error)
	DeleteRelation(ctx context.Context, relationID string) error

	// 辅助方法
	ExistsByID(ctx context.Context, locationID string) (bool, error)
	CountByProjectID(ctx context.Context, projectID string) (int64, error)
}
