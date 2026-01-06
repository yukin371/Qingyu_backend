package bookstore

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChapterContent 章节内容模型（与元数据分离）
type ChapterContent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ChapterID primitive.ObjectID `bson:"chapter_id" json:"chapterId" binding:"required"`
	Content   string             `bson:"content" json:"content" binding:"required"`
	Format    string             `bson:"format" json:"format"`                     // markdown, html, txt
	Version   int                `bson:"version" json:"version"`                   // 版本号

	// 元数据
	WordCount int       `bson:"word_count" json:"wordCount"` // 字数统计
	Hash      string    `bson:"hash,omitempty" json:"hash,omitempty"` // 内容哈希（用于校验和去重）

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

// ContentFormat 内容格式常量
const (
	ContentFormatMarkdown string = "markdown" // Markdown 格式
	ContentFormatHTML     string = "html"     // HTML 格式
	ContentFormatText     string = "text"     // 纯文本
)

// BeforeCreate 在创建前设置默认值
func (cc *ChapterContent) BeforeCreate() {
	now := time.Now()
	cc.CreatedAt = now
	cc.UpdatedAt = now
	if cc.Format == "" {
		cc.Format = ContentFormatMarkdown
	}
	if cc.Version == 0 {
		cc.Version = 1
	}
}

// BeforeUpdate 在更新前刷新时间戳和版本号
func (cc *ChapterContent) BeforeUpdate() {
	cc.UpdatedAt = time.Now()
	cc.Version++
}

// UpdateContent 更新内容
func (cc *ChapterContent) UpdateContent(content string) {
	cc.Content = content
	cc.WordCount = len([]rune(content))
	cc.BeforeUpdate()
}

// IsMarkdown 判断是否为 Markdown 格式
func (cc *ChapterContent) IsMarkdown() bool {
	return cc.Format == ContentFormatMarkdown
}

// CalculateHash 计算内容哈希
func (cc *ChapterContent) CalculateHash() string {
	// 简单哈希实现，生产环境建议使用 crypto/sha256
	return fmt.Sprintf("%s:%d", cc.ChapterID.Hex(), cc.Version)
}
