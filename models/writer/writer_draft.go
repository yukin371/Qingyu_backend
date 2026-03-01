package writer

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WriterDraft AI写作草稿模型
// 用于AI写作助手功能中的文档草稿管理
type WriterDraft struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID string             `bson:"project_id" json:"project_id"`
	ChapterNum int               `bson:"chapter_num" json:"chapter_num"`
	Title      string            `bson:"title" json:"title"`
	Content    string            `bson:"content" json:"content"`
	Format     string            `bson:"format" json:"format"` // markdown, html, text
	WordCount  int               `bson:"word_count" json:"word_count"`
	Version    int               `bson:"version" json:"version"`

	// 元数据
	Status string   `bson:"status" json:"status"` // draft, reviewing, completed
	Tags   []string `bson:"tags,omitempty" json:"tags,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// DraftStatus 草稿状态常量
const (
	DraftStatusDraft      = "draft"      // 草稿中
	DraftStatusReviewing  = "reviewing"  // 审核中
	DraftStatusCompleted  = "completed"  // 已完成
)

// ContentFormat 内容格式常量
const (
	FormatMarkdown = "markdown" // Markdown格式
	FormatHTML     = "html"     // HTML格式
	FormatText     = "text"     // 纯文本
)

// BeforeCreate 创建前设置默认值
func (d *WriterDraft) BeforeCreate() {
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	if d.Format == "" {
		d.Format = FormatMarkdown
	}
	if d.Status == "" {
		d.Status = DraftStatusDraft
	}
	if d.Version == 0 {
		d.Version = 1
	}
}

// BeforeUpdate 更新前刷新时间戳
func (d *WriterDraft) BeforeUpdate() {
	d.UpdatedAt = time.Now()
	d.Version++
}

// UpdateContent 更新内容并计算字数
func (d *WriterDraft) UpdateContent(content string) {
	d.Content = content
	d.WordCount = len([]rune(content))
	d.BeforeUpdate()
}

// IsMarkdown 判断是否为Markdown格式
func (d *WriterDraft) IsMarkdown() bool {
	return d.Format == FormatMarkdown
}

// IsValidStatus 验证状态是否有效
func (d *WriterDraft) IsValidStatus() bool {
	switch d.Status {
	case DraftStatusDraft, DraftStatusReviewing, DraftStatusCompleted:
		return true
	}
	return false
}

// IsValidFormat 验证格式是否有效
func (d *WriterDraft) IsValidFormat() bool {
	switch d.Format {
	case FormatMarkdown, FormatHTML, FormatText:
		return true
	}
	return false
}
