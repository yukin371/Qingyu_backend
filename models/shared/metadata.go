package shared

import "time"

// Pinned 置顶状态混入
type Pinned struct {
	IsPinned bool       `bson:"is_pinned" json:"isPinned"`
	PinnedAt *time.Time `bson:"pinned_at,omitempty" json:"pinnedAt,omitempty"`
	PinnedBy *string    `bson:"pinned_by,omitempty" json:"pinnedBy,omitempty"`
}

// Pin 置顶
func (p *Pinned) Pin(operatorID string) {
	p.IsPinned = true
	now := time.Now()
	p.PinnedAt = &now
	p.PinnedBy = &operatorID
}

// Unpin 取消置顶
func (p *Pinned) Unpin() {
	p.IsPinned = false
	p.PinnedAt = nil
	p.PinnedBy = nil
}

// Expirable 有效期混入
type Expirable struct {
	ExpiresAt *time.Time `bson:"expires_at,omitempty" json:"expiresAt,omitempty"`
}

// IsExpired 判断是否已过期
func (e *Expirable) IsExpired() bool {
	if e.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*e.ExpiresAt)
}

// SetExpiration 设置过期时间
func (e *Expirable) SetExpiration(duration time.Duration) {
	expiresAt := time.Now().Add(duration)
	e.ExpiresAt = &expiresAt
}

// TargetEntity 关联实体混入
type TargetEntity struct {
	TargetType string `bson:"target_type" json:"targetType"`
	TargetID   string `bson:"target_id" json:"targetId"`
}
