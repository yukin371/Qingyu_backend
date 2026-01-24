package writer

import (
	"Qingyu_backend/models/writer/base"
	"Qingyu_backend/models/writer/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Document 文档模型（用于文档内容）
// 文档类型现在是动态的，由项目的writing_type决定
type Document struct {
	base.IdentifiedEntity `bson:",inline"` // ID
	base.Timestamps       `bson:",inline"` // CreatedAt, UpdatedAt, DeletedAt

	// 项目基础信息
	ProjectID primitive.ObjectID `bson:"project_id" json:"projectId" validate:"required"`
	Title     string             `bson:"title" json:"title" validate:"required,max=200"`

	// ===== v1.1 新增：稳定引用标识 =====
	StableRef string `bson:"stable_ref" json:"stableRef" validate:"required"`
	OrderKey  string `bson:"order_key" json:"orderKey" validate:"required"` // LexoRank排序键
	// ===== v1.1 新增结束 =====

	// 层级结构（保持BSON字段名不变，确保数据库兼容）
	ParentID primitive.ObjectID `bson:"parent_id,omitempty" json:"parentId,omitempty"` // 父文档ID，空值表示根文档
	Type     string             `bson:"type" json:"type" validate:"required"`          // 动态类型（如：volume, chapter, section, scene, article, act, beat等）
	Level    int                `bson:"level" json:"level" validate:"min=0"`           // 层级深度
	Order    int                `bson:"order" json:"order" validate:"min=0"`           // 同级排序
	Status   string             `bson:"status" json:"status"`                          // planned | writing | completed

	// 统计信息（从DocumentContent同步）
	WordCount int `bson:"word_count" json:"wordCount"` // 字数统计

	// 关联信息（保持BSON字段名不变，确保数据库兼容）
	CharacterIDs []primitive.ObjectID `bson:"character_ids,omitempty" json:"characterIds,omitempty" validate:"max=50"` // 关联角色
	LocationIDs  []primitive.ObjectID `bson:"location_ids,omitempty" json:"locationIds,omitempty" validate:"max=50"`   // 关联地点
	TimelineIDs  []primitive.ObjectID `bson:"timeline_ids,omitempty" json:"timelineIds,omitempty" validate:"max=20"`   // 关联时间线

	// 写作辅助（AI上下文需要）
	PlotThreads  []string `bson:"plot_threads,omitempty" json:"plotThreads,omitempty" validate:"max=20"`   // 情节线索
	KeyPoints    []string `bson:"key_points,omitempty" json:"keyPoints,omitempty" validate:"max=20"`       // 关键点
	WritingHints []string `bson:"writing_hints,omitempty" json:"writingHints,omitempty" validate:"max=20"` // 写作提示

	// 标签和备注
	Tags  []string `bson:"tags,omitempty" json:"tags,omitempty" validate:"max=10"`
	Notes string   `bson:"notes,omitempty" json:"notes,omitempty" validate:"max=500"`
}

// DocumentStatus 文档状态常量（保持向后兼容）
const (
	DocumentStatusPlanned   = "planned"   // 计划中
	DocumentStatusWriting   = "writing"   // 写作中
	DocumentStatusCompleted = "completed" // 已完成

	// 默认OrderKey（用于新文档）
	DefaultOrderKey = "a0"
)

// 历史文档类型常量（保持向后兼容，但已废弃，建议使用WritingType系统）
// 这些常量仅用于小说类型的项目
const (
	TypeVolume  = "volume"  // 卷（小说专用）
	TypeChapter = "chapter" // 章（小说专用）
	TypeSection = "section" // 节（小说/文章通用）
	TypeScene   = "scene"   // 场景（小说/剧本通用）
)

// IsRoot 判断是否为根文档
func (d *Document) IsRoot() bool {
	return d.ParentID.IsZero()
}

// CanHaveChildren 判断是否可以有子文档
// 需要传入项目的writing_type来动态判断
func (d *Document) CanHaveChildren(writingType string) bool {
	wt, ok := types.GlobalRegistry.Get(writingType)
	if !ok {
		return false // 未知的writing_type，不允许子节点
	}
	return wt.CanHaveChildren(d.Type)
}

// GetParentType 获取父级文档类型
func (d *Document) GetParentType(writingType string) (string, bool) {
	wt, ok := types.GlobalRegistry.Get(writingType)
	if !ok {
		return "", false
	}
	return wt.GetParentType(d.Type)
}

// GetChildTypes 获取允许的子文档类型列表
func (d *Document) GetChildTypes(writingType string) []string {
	wt, ok := types.GlobalRegistry.Get(writingType)
	if !ok {
		return nil
	}
	return wt.GetChildTypes(d.Type)
}

// GetNextLevel 获取下一层级
func (d *Document) GetNextLevel() int {
	return d.Level + 1
}

// UpdateWordCount 更新字数统计
func (d *Document) UpdateWordCount(count int) {
	d.WordCount = count
	d.Timestamps.Touch()
}

// TouchForCreate 创建时设置默认值
func (d *Document) TouchForCreate() {
	d.IdentifiedEntity.GenerateID()
	d.Timestamps.TouchForCreate()
	if d.Status == "" {
		d.Status = DocumentStatusPlanned
	}
}

// TouchForUpdate 更新时设置默认值
func (d *Document) TouchForUpdate() {
	d.Timestamps.Touch()
}

// Validate 验证文档数据
// 需要传入项目的writing_type来验证document type是否有效
func (d *Document) Validate(writingType string) error {
	// 验证标题
	if err := base.ValidateTitle(d.Title, 200); err != nil {
		return err
	}

	// 验证项目ID
	if d.ProjectID.IsZero() {
		return base.ErrProjectIDRequired
	}

	// 验证文档类型（基于writing_type）
	wt, ok := types.GlobalRegistry.Get(writingType)
	if !ok {
		return base.ErrInvalidWritingType
	}
	if !wt.ValidateDocumentType(d.Type) {
		return base.ErrInvalidDocumentType
	}

	// 验证层级（基于writing_type的最大层级）
	maxDepth := wt.GetMaxDepth()
	if d.Level < 0 || d.Level >= maxDepth {
		return base.ErrInvalidLevel
	}

	return nil
}

// ValidateWithoutType 验证文档数据（不验证类型，用于向后兼容）
func (d *Document) ValidateWithoutType() error {
	// 验证标题
	if err := base.ValidateTitle(d.Title, 200); err != nil {
		return err
	}

	// 验证项目ID
	if d.ProjectID.IsZero() {
		return base.ErrProjectIDRequired
	}

	// 验证层级（通用验证）
	if d.Level < 0 || d.Level > 10 {
		return base.ErrInvalidLevel
	}

	return nil
}
