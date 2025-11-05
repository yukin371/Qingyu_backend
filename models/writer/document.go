package writer

import (
	"fmt"
	"time"
)

// Document 文档模型（用于文档内容）
type Document struct {
	ID        string       `bson:"_id,omitempty" json:"id"`
	ProjectID string       `bson:"project_id" json:"projectId" validate:"required"`
	ParentID  string       `bson:"parent_id,omitempty" json:"parentId,omitempty"` // 父文档ID，空表示根文档
	Title     string       `bson:"title" json:"title" validate:"required,min=1,max=200"`
	Type      DocumentType `bson:"type" json:"type" validate:"required"`
	Level     int          `bson:"level" json:"level"`   // 层级深度（0-2）
	Order     int          `bson:"order" json:"order"`   // 同级排序
	Status    string       `bson:"status" json:"status"` // planned | writing | completed

	// 统计信息（从DocumentContent同步）
	WordCount int `bson:"word_count" json:"wordCount"` // 字数统计

	// 关联信息
	CharacterIDs []string `bson:"character_ids,omitempty" json:"characterIds,omitempty"` // 关联角色
	LocationIDs  []string `bson:"location_ids,omitempty" json:"locationIds,omitempty"`   // 关联地点
	TimelineIDs  []string `bson:"timeline_ids,omitempty" json:"timelineIds,omitempty"`   // 关联时间线

	// 写作辅助（AI上下文需要）
	PlotThreads  []string `bson:"plot_threads,omitempty" json:"plotThreads,omitempty"`   // 情节线索
	KeyPoints    []string `bson:"key_points,omitempty" json:"keyPoints,omitempty"`       // 关键点
	WritingHints []string `bson:"writing_hints,omitempty" json:"writingHints,omitempty"` // 写作提示

	// 标签和备注
	Tags  []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Notes string   `bson:"notes,omitempty" json:"notes,omitempty"`

	// 时间戳
	CreatedAt time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}

// DocumentType 文档类型
type DocumentType string

const (
	TypeVolume  DocumentType = "volume"  // 卷
	TypeChapter DocumentType = "chapter" // 章
	TypeSection DocumentType = "section" // 节
	TypeScene   DocumentType = "scene"   // 场景
)

// IsRoot 判断是否为根文档
func (d *Document) IsRoot() bool {
	return d.ParentID == ""
}

// CanHaveChildren 判断是否可以有子文档
func (d *Document) CanHaveChildren() bool {
	return d.Level < 2 // 最多3层，0-2
}

// GetNextLevel 获取下一层级
func (d *Document) GetNextLevel() int {
	return d.Level + 1
}

// UpdateWordCount 更新字数统计
func (d *Document) UpdateWordCount(count int) {
	d.WordCount = count
	d.UpdatedAt = time.Now()
}

// Validate 验证文档数据
func (d *Document) Validate() error {
	if d.ProjectID == "" {
		return fmt.Errorf("项目ID不能为空")
	}
	if d.Title == "" {
		return fmt.Errorf("文档标题不能为空")
	}
	if len(d.Title) > 200 {
		return fmt.Errorf("文档标题不能超过200字符")
	}
	if !d.Type.IsValid() {
		return fmt.Errorf("无效的文档类型: %s", d.Type)
	}
	if d.Level < 0 || d.Level > 2 {
		return fmt.Errorf("文档层级必须在0-2之间")
	}
	return nil
}

// IsValid 验证文档类型是否有效
func (t DocumentType) IsValid() bool {
	switch t {
	case TypeVolume, TypeChapter, TypeSection, TypeScene:
		return true
	}
	return false
}

// String 返回文档类型的字符串表示
func (t DocumentType) String() string {
	return string(t)
}

// TouchForCreate 设置创建时的默认值
func (d *Document) TouchForCreate() {
	now := time.Now()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	if d.UpdatedAt.IsZero() {
		d.UpdatedAt = now
	}
	if d.Status == "" {
		d.Status = "planned"
	}
}

// TouchForUpdate 设置更新时的默认值
func (d *Document) TouchForUpdate() {
	d.UpdatedAt = time.Now()
}
