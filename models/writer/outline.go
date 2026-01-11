package writer

import (
	"Qingyu_backend/models/writer/base"
)

// OutlineNode 大纲节点（写作规划专用）
// 这是专门用于写作规划的大纲节点，包含情节、紧张度等写作相关属性
// 与通用的Node不同，OutlineNode专注于写作内容的规划
type OutlineNode struct {
	base.IdentifiedEntity    `bson:",inline"` // ID
	base.Timestamps          `bson:",inline"` // CreatedAt, UpdatedAt
	base.ProjectScopedEntity `bson:",inline"` // ProjectID
	base.TitledEntity        `bson:",inline"` // Title

	// 层级结构（保持BSON字段名不变，确保数据库兼容）
	ParentID string `bson:"parent_id,omitempty" json:"parentId,omitempty"` // 父级大纲节点ID（如：卷/幕）
	Order    int    `bson:"order" json:"order" validate:"min=0"`           // 同级排序

	// 写作规划属性
	Summary   string `bson:"summary,omitempty" json:"summary,omitempty" validate:"max=1000"` // 本节摘要
	Type      string `bson:"type,omitempty" json:"type,omitempty" validate:"max=50"`         // 结构类型（如：英雄之旅阶段、起承转合）
	Tension   int    `bson:"tension" json:"tension" validate:"min=0,max=10"`                 // 紧张度/情绪值，用于生成情绪曲线

	// 关联实体（保持BSON字段名不变，确保数据库兼容）
	ChapterID  string   `bson:"chapter_id,omitempty" json:"chapterId,omitempty"`  // 对应实际写的章节ID
	Characters []string `bson:"characters,omitempty" json:"characters,omitempty" validate:"max=20"` // 本节登场人物ID列表
	Items      []string `bson:"items,omitempty" json:"items,omitempty" validate:"max=20"`           // 涉及道具ID列表
}

// TouchForCreate 创建时设置默认值
func (o *OutlineNode) TouchForCreate() {
	o.IdentifiedEntity.GenerateID()
	o.Timestamps.TouchForCreate()
}

// TouchForUpdate 更新时设置默认值
func (o *OutlineNode) TouchForUpdate() {
	o.Timestamps.Touch()
}

// IsRoot 判断是否为根节点
func (o *OutlineNode) IsRoot() bool {
	return o.ParentID == ""
}

// Validate 验证大纲节点数据
func (o *OutlineNode) Validate() error {
	// 验证标题
	if err := base.ValidateTitle(o.Title, 200); err != nil {
		return err
	}

	// 验证项目ID
	if o.ProjectID == "" {
		return base.ErrProjectIDRequired
	}

	return nil
}
