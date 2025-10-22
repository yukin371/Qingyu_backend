package recommendation

import (
	reco "Qingyu_backend/models/recommendation/reco"
	"context"
)

// ProfileRepository 用户画像仓储接口
type ProfileRepository interface {
	Upsert(ctx context.Context, p *reco.UserProfile) error
	GetByUserID(ctx context.Context, userID string) (*reco.UserProfile, error)
}
