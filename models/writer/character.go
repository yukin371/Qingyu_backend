package writer

import (
	"Qingyu_backend/models/writer/base"
)

// Character 角色卡片
type Character struct {
	base.IdentifiedEntity    `bson:",inline"` // ID
	base.Timestamps          `bson:",inline"` // 时间戳
	base.ProjectScopedEntity `bson:",inline"` // ProjectID
	base.NamedEntity         `bson:",inline"` // Name

	// 基本信息（保持BSON字段名不变，确保数据库兼容）
	Alias         []string `bson:"alias,omitempty" json:"alias,omitempty" validate:"max=10"`
	Summary       string   `bson:"summary,omitempty" json:"summary,omitempty" validate:"max=500"`
	Traits        []string `bson:"traits,omitempty" json:"traits,omitempty" validate:"max=20"`
	Background    string   `bson:"background,omitempty" json:"background,omitempty" validate:"max=2000"`
	AvatarURL     string   `bson:"avatar_url,omitempty" json:"avatarUrl,omitempty" validate:"omitempty,url"`

	// AI相关字段
	PersonalityPrompt string `bson:"personality_prompt,omitempty" json:"personalityPrompt,omitempty" validate:"max=1000"` // 角色性格提示
	SpeechPattern     string `bson:"speech_pattern,omitempty" json:"speechPattern,omitempty" validate:"max=500"`         // 角色语音模式
	CurrentState      string `bson:"current_state,omitempty" json:"currentState,omitempty" validate:"max=200"`           // 角色当前状态
	ShortDescription  string `bson:"short_description,omitempty" json:"shortDescription,omitempty" validate:"max=200"`    // 角色摘要
}

// TouchForCreate 创建时设置默认值
func (c *Character) TouchForCreate() {
	c.IdentifiedEntity.GenerateID()
	c.Timestamps.TouchForCreate()
}

// Validate 验证角色数据
func (c *Character) Validate() error {
	if err := base.ValidateName(c.Name, 100); err != nil {
		return err
	}
	if c.ProjectID == "" {
		return base.ErrProjectIDRequired
	}
	if err := base.ValidateURL(c.AvatarURL); err != nil {
		return err
	}
	return nil
}

// RelationType 关系类型
type RelationType string

const (
	RelationFriend  RelationType = "朋友"
	RelationFamily  RelationType = "家庭"
	RelationEnemy   RelationType = "敌人"
	RelationRomance RelationType = "恋人"
	RelationAlly    RelationType = "盟友"
	RelationOther   RelationType = "其他"
)

// CharacterRelation 角色关系（边表）
// 允许有方向，也可约定对称关系由 Service 保持双向两条边
type CharacterRelation struct {
	base.IdentifiedEntity    `bson:",inline"`
	base.Timestamps          `bson:",inline"`
	base.ProjectScopedEntity `bson:",inline"`

	FromID   string       `bson:"from_id" json:"fromId" validate:"required"`
	ToID     string       `bson:"to_id" json:"toId" validate:"required"`
	Type     RelationType `bson:"type" json:"type" validate:"required"`
	Strength int          `bson:"strength" json:"strength" validate:"min=0,max=100"`
	Notes    string       `bson:"notes,omitempty" json:"notes,omitempty" validate:"max=500"`
}

// TouchForCreate 创建时设置默认值
func (cr *CharacterRelation) TouchForCreate() {
	cr.IdentifiedEntity.GenerateID()
	cr.Timestamps.TouchForCreate()
}

// Validate 验证角色关系
func (cr *CharacterRelation) Validate() error {
	if cr.ProjectID == "" {
		return base.ErrProjectIDRequired
	}
	if !cr.Type.IsValid() {
		return base.ErrInvalidEnum
	}
	return nil
}

// IsValid 验证关系类型是否合法
func (rt RelationType) IsValid() bool {
	switch rt {
	case RelationFriend, RelationFamily, RelationEnemy, RelationRomance, RelationAlly, RelationOther:
		return true
	default:
		return false
	}
}

// IsValidRelationType 验证关系类型是否合法（保持向后兼容）
func IsValidRelationType(t string) bool {
	return RelationType(t).IsValid()
}
