package document

import "time"

// Character 角色卡片
type Character struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	ProjectID  string    `bson:"project_id" json:"projectId"`
	Name       string    `bson:"name" json:"name"`
	Alias      []string  `bson:"alias,omitempty" json:"alias,omitempty"`
	Summary    string    `bson:"summary,omitempty" json:"summary,omitempty"`
	Traits     []string  `bson:"traits,omitempty" json:"traits,omitempty"` // 性格标签
	Background string    `bson:"background,omitempty" json:"background,omitempty"`
	AvatarURL  string    `bson:"avatar_url,omitempty" json:"avatarUrl,omitempty"`
	CreatedAt  time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updatedAt"`
}

// RelationType 关系类型
type RelationType string

const (
	RelationFriend  RelationType = "friend"
	RelationFamily  RelationType = "family"
	RelationEnemy   RelationType = "enemy"
	RelationRomance RelationType = "romance"
	RelationAlly    RelationType = "ally"
	RelationOther   RelationType = "other"
)

// CharacterRelation 角色关系（边表）
// 允许有方向，也可约定对称关系由 Service 保持双向两条边
type CharacterRelation struct {
	ID        string       `bson:"_id,omitempty" json:"id"`
	ProjectID string       `bson:"project_id" json:"projectId"`
	FromID    string       `bson:"from_id" json:"fromId"`
	ToID      string       `bson:"to_id" json:"toId"`
	Type      RelationType `bson:"type" json:"type"`
	Strength  int          `bson:"strength" json:"strength"` // 0-100 强度
	Notes     string       `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt time.Time    `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time    `bson:"updated_at" json:"updatedAt"`
}
