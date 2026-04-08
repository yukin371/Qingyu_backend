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
	Alias      []string `bson:"alias,omitempty" json:"alias,omitempty" validate:"max=10"`
	Summary    string   `bson:"summary,omitempty" json:"summary,omitempty" validate:"max=500"`
	Traits     []string `bson:"traits,omitempty" json:"traits,omitempty" validate:"max=20"`
	Background string   `bson:"background,omitempty" json:"background,omitempty" validate:"max=2000"`
	AvatarURL  string   `bson:"avatar_url,omitempty" json:"avatarUrl,omitempty" validate:"omitempty,url"`

	// AI相关字段
	PersonalityPrompt string `bson:"personality_prompt,omitempty" json:"personalityPrompt,omitempty" validate:"max=1000"` // 角色性格提示
	SpeechPattern     string `bson:"speech_pattern,omitempty" json:"speechPattern,omitempty" validate:"max=500"`          // 角色语音模式
	CurrentState      string `bson:"current_state,omitempty" json:"currentState,omitempty" validate:"max=200"`            // 角色当前状态
	ShortDescription  string `bson:"short_description,omitempty" json:"shortDescription,omitempty" validate:"max=200"`    // 角色摘要

	// 可选：汇总的登场信息（用于快速查询，非核心字段）
	AppearanceCount    int     `bson:"appearance_count,omitempty" json:"appearanceCount,omitempty"`         // 登场次数
	FirstAppearanceID  *string `bson:"first_appearance_id,omitempty" json:"firstAppearanceId,omitempty"`   // 首次登场的大纲节点ID

	// 实体类型扩展
	EntityType    EntityType            `bson:"entity_type" json:"entityType"`
	StateFields   map[string]StateValue `bson:"state_fields,omitempty" json:"stateFields,omitempty"`
	EquippedItems []string              `bson:"equipped_items,omitempty" json:"equippedItems,omitempty"`
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
	if c.ProjectID.IsZero() {
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
	FromType EntityType   `bson:"from_type" json:"fromType"`
	ToType   EntityType   `bson:"to_type" json:"toType"`
	Type     RelationType `bson:"type" json:"type" validate:"required"`
	Strength int          `bson:"strength" json:"strength" validate:"min=0,max=100"`
	Notes    string       `bson:"notes,omitempty" json:"notes,omitempty" validate:"max=500"`

	// 时序控制字段（支持关系随剧情发展变化）
	ValidFromChapterID  *string                 `bson:"valid_from_chapter_id,omitempty" json:"validFromChapterId,omitempty"`   // 关系生效的起始章节ID
	ValidUntilChapterID *string                 `bson:"valid_until_chapter_id,omitempty" json:"validUntilChapterId,omitempty"` // 关系失效的章节ID
	TimelineEvents      []RelationTimelineEvent `bson:"timeline_events,omitempty" json:"timelineEvents,omitempty"`               // 变化历史记录
}

// TouchForCreate 创建时设置默认值
func (cr *CharacterRelation) TouchForCreate() {
	cr.IdentifiedEntity.GenerateID()
	cr.Timestamps.TouchForCreate()
}

// Validate 验证角色关系
func (cr *CharacterRelation) Validate() error {
	if cr.ProjectID.IsZero() {
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

// IsValidAtChapter 判断关系在某章节是否有效
// chapterOrder: 当前章节的序号
// chapterOrderMap: 章节ID到序号的映射表
func (cr *CharacterRelation) IsValidAtChapter(chapterOrder int, chapterOrderMap map[string]int) bool {
	// 如果没有设置生效章节，默认全局有效
	if cr.ValidFromChapterID == nil && cr.ValidUntilChapterID == nil {
		return true
	}

	// 检查是否已经生效
	if cr.ValidFromChapterID != nil {
		fromOrder, exists := chapterOrderMap[*cr.ValidFromChapterID]
		if !exists {
			// 如果起始章节不存在，保守处理：认为关系有效
			return true
		}
		if fromOrder > chapterOrder {
			// 还未到生效章节
			return false
		}
	}

	// 检查是否已经失效
	if cr.ValidUntilChapterID != nil {
		untilOrder, exists := chapterOrderMap[*cr.ValidUntilChapterID]
		if !exists {
			// 如果失效章节不存在，保守处理：认为关系有效
			return true
		}
		if untilOrder <= chapterOrder {
			// 已经过了失效章节
			return false
		}
	}

	return true
}
