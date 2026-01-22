package writer

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentContent 文档内容
// 用于存储文档的实际内容，与Document分离以提升性能
type DocumentContent struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
	DeletedAt    time.Time          `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
	LastSavedAt  time.Time          `bson:"last_saved_at" json:"lastSavedAt"`
	LastEditedBy string             `bson:"last_edited_by" json:"lastEditedBy"`

	DocumentID  primitive.ObjectID `bson:"document_id" json:"documentId" validate:"required"`
	Content     string             `bson:"content" json:"content"`                        // 文档内容
	ContentType string             `bson:"content_type" json:"contentType"`               // markdown | richtext
	WordCount   int                `bson:"word_count" json:"wordCount"`                   // 字数统计
	CharCount   int                `bson:"char_count" json:"charCount"`                   // 字符统计
	GridFSID    primitive.ObjectID `bson:"gridfs_id,omitempty" json:"gridfsId,omitempty"` // 大文件GridFS ID
	Version     int                `bson:"version" json:"version"`                        // 版本号（乐观锁）
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
	if d.DocumentID.IsZero() {
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
		d.LastSavedAt = d.CreatedAt
	}
	if d.Version == 0 {
		d.Version = 1
	}
	if d.ID.IsZero() {
		d.ID = primitive.NewObjectID()
	}
}

// TouchForUpdate 设置更新时的默认值
func (d *DocumentContent) TouchForUpdate() {
	d.UpdatedAt = time.Now()
	d.LastSavedAt = time.Now()
}

// GetID 获取ID
func (d *DocumentContent) GetID() primitive.ObjectID {
	return d.ID
}

// SetID 设置ID
func (d *DocumentContent) SetID(id primitive.ObjectID) {
	d.ID = id
}
