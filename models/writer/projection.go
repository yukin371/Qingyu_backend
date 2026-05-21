package writer

import (
	"time"

	"Qingyu_backend/models/writer/base"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CharacterSnapshot 章节级角色投影快照。
type CharacterSnapshot struct {
	CharacterID   string                `bson:"character_id" json:"characterId"`
	CharacterName string                `bson:"character_name" json:"characterName"`
	EntityType    EntityType            `bson:"entity_type" json:"entityType"`
	Summary       string                `bson:"summary,omitempty" json:"summary,omitempty"`
	CurrentState  string                `bson:"current_state,omitempty" json:"currentState,omitempty"`
	StateFields   map[string]StateValue `bson:"state_fields,omitempty" json:"stateFields,omitempty"`
}

// RelationSnapshot 章节级关系投影快照。
type RelationSnapshot struct {
	FromID   string `bson:"from_id" json:"fromId"`
	ToID     string `bson:"to_id" json:"toId"`
	FromName string `bson:"from_name,omitempty" json:"fromName,omitempty"`
	ToName   string `bson:"to_name,omitempty" json:"toName,omitempty"`
	Relation string `bson:"relation,omitempty" json:"relation,omitempty"`
	Strength int    `bson:"strength,omitempty" json:"strength,omitempty"`
}

// ProjectionCheckpoint 记录最近一次投影刷新来源。
type ProjectionCheckpoint struct {
	LastRequestID string                `bson:"last_request_id,omitempty" json:"lastRequestId,omitempty"`
	LastCategory  ChangeRequestCategory `bson:"last_category,omitempty" json:"lastCategory,omitempty"`
	RefreshedAt   time.Time             `bson:"refreshed_at,omitempty" json:"refreshedAt,omitempty"`
}

// ChapterProjection 章节上下文投影。
type ChapterProjection struct {
	base.IdentifiedEntity    `bson:",inline"`
	base.Timestamps          `bson:",inline"`
	base.ProjectScopedEntity `bson:",inline"`

	ChapterID  primitive.ObjectID   `bson:"chapter_id" json:"chapterId"`
	Characters []CharacterSnapshot  `bson:"characters,omitempty" json:"characters,omitempty"`
	Relations  []RelationSnapshot   `bson:"relations,omitempty" json:"relations,omitempty"`
	Checkpoint ProjectionCheckpoint `bson:"checkpoint,omitempty" json:"checkpoint,omitempty"`
}

// TouchForCreate 创建时设置默认值。
func (p *ChapterProjection) TouchForCreate() {
	p.IdentifiedEntity.GenerateID()
	p.Timestamps.TouchForCreate()
}
