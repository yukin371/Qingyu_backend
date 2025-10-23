package writer

import (
	"fmt"
	"time"
)

// DocumentContent 文档内容
// 用于存储文档的实际内容，与Document分离以提升性能
type DocumentContent struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	DocumentID   string    `bson:"document_id" json:"documentId" validate:"required"`
	Content      string    `bson:"content" json:"content"`                        // 文档内容
	ContentType  string    `bson:"content_type" json:"contentType"`               // markdown | richtext
	WordCount    int       `bson:"word_count" json:"wordCount"`                   // 字数统计
	CharCount    int       `bson:"char_count" json:"charCount"`                   // 字符统计
	GridFSID     string    `bson:"gridfs_id,omitempty" json:"gridfsId,omitempty"` // 大文件GridFS ID
	Version      int       `bson:"version" json:"version"`                        // 版本号（乐观锁）
	LastSavedAt  time.Time `bson:"last_saved_at" json:"lastSavedAt"`              // 最后保存时间
	LastEditedBy string    `bson:"last_edited_by" json:"lastEditedBy"`            // 最后编辑人
	UpdatedAt    time.Time `bson:"updated_at" json:"updatedAt"`
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`
}

// IsLargeDocument 判断是否为大文档（>1MB）
func (d *DocumentContent) IsLargeDocument() bool {
	return len(d.Content) > 1024*1024
}

// GetDisplayWordCount 获取显示用的字数
func (d *DocumentContent) GetDisplayWordCount() int {
	if d.WordCount > 0 {
		return d.WordCount
	}
	return len([]rune(d.Content))
}

// Validate 验证文档内容数据
func (d *DocumentContent) Validate() error {
	if d.DocumentID == "" {
		return fmt.Errorf("文档ID不能为空")
	}
	if d.ContentType == "" {
		return fmt.Errorf("内容类型不能为空")
	}
	if d.ContentType != "markdown" && d.ContentType != "richtext" {
		return fmt.Errorf("无效的内容类型: %s", d.ContentType)
	}
	return nil
}

// TouchForCreate 设置创建时的默认值
func (d *DocumentContent) TouchForCreate() {
	now := time.Now()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	if d.UpdatedAt.IsZero() {
		d.UpdatedAt = now
	}
	if d.LastSavedAt.IsZero() {
		d.LastSavedAt = now
	}
	if d.Version == 0 {
		d.Version = 1
	}
}

// TouchForUpdate 设置更新时的默认值
func (d *DocumentContent) TouchForUpdate() {
	d.UpdatedAt = time.Now()
	d.LastSavedAt = time.Now()
}
