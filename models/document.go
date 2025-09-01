package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Document 表示系统中的一个文档实体
// 包含文档的基本信息、内容以及版本历史
type Document struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // MongoDB 唯一标识符
	UserID    string             `bson:"userId" json:"userId"`              // 文档所属用户ID
	Title     string             `bson:"title" json:"title"`                // 文档标题
	Content   string             `bson:"content" json:"content"`            // 文档当前内容
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`        // 文档创建时间
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`        // 文档最后更新时间
	Versions  []DocumentVersion  `bson:"versions" json:"versions"`          // 文档的历史版本集合
}

// DocumentVersion 表示文档的一个历史版本
// 记录了特定时间点的文档内容
type DocumentVersion struct {
	Title     string             `bson:"title" json:"title"`         // 历史标题
	VersionID primitive.ObjectID `bson:"versionId" json:"versionId"` // 版本唯一标识符
	Content   string             `bson:"content" json:"content"`     // 历史内容
	Time      time.Time          `bson:"time" json:"time"`           // 版本创建时间
}
