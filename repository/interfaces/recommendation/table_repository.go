package recommendation

import (
	reco "Qingyu_backend/models/recommendation"
	"context"
)

// TableRepository 推荐榜表仓储接口
type TableRepository interface {
	Create(ctx context.Context, table *reco.RecommendationTable) error
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*reco.RecommendationTable, error)
	GetByTypePeriod(ctx context.Context, tableType reco.TableType, period string, source reco.TableSource) (*reco.RecommendationTable, error)
	List(ctx context.Context, tableType *reco.TableType, source *reco.TableSource, page, pageSize int) ([]*reco.RecommendationTable, int64, error)
	UpsertByTypePeriod(ctx context.Context, table *reco.RecommendationTable) error
	Health(ctx context.Context) error
}
