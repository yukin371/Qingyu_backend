package writer

import "time"

// Concept 设定概念
type Concept struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	ProjectID string `bson:"project_id" json:"projectId"`
	Name      string `bson:"name" json:"name"`
	Category  string `bson:"category" json:"category"` // 分类：魔法系统、地理环境、生物图鉴、历史背景
	Content   string `bson:"content" json:"content"`   // 详细设定内容 (Markdown/RichText)

	// 关联
	RelatedDocs []string `bson:"related_docs,omitempty" json:"relatedDocs,omitempty"` // 关联的文档ID
	Tags        []string `bson:"tags,omitempty" json:"tags,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
