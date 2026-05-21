package discovery

import (
	discoveryModel "Qingyu_backend/models/discovery"
	"context"
)

// PreferenceRepository 持有 discovery 页自有轻量偏好配置。
type PreferenceRepository interface {
	GetByUserID(ctx context.Context, userID string) (*discoveryModel.PreferenceProfile, error)
	UpsertByUserID(ctx context.Context, userID string, profile *discoveryModel.PreferenceProfile) error
}
