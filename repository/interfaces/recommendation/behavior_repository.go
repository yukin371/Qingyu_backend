package recommendation

import (
	reco "Qingyu_backend/models/recommendation/reco"
	"context"
)

// BehaviorRepository 用户行为仓储接口
type BehaviorRepository interface {
	Create(ctx context.Context, b *reco.Behavior) error
	BatchCreate(ctx context.Context, bs []*reco.Behavior) error
	GetByUser(ctx context.Context, userID string, limit int) ([]*reco.Behavior, error)
}
