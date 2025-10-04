package reader

import (
	"time"
)

type AnnotationType string

const (
	AnnotationTypeBookmark  AnnotationType = "bookmark"  // 书签
	AnnotationTypeHighlight AnnotationType = "highlight" // 划线
)

type Annotation struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	UserID    string    `bson:"user_id" json:"userId"`       // 用户ID
	BookID    string    `bson:"book_id" json:"bookId"`       // 书籍ID
	ChapterID string    `bson:"chapter_id" json:"chapterId"` // 章节ID
	Range     string    `bson:"range" json:"range"`          // 标注范围：start-end
	Text      string    `bson:"text" json:"text"`            // 标注文本
	Note      string    `bson:"note" json:"note"`            // 注释
	Type      string    `bson:"type" json:"type"`            // 标注类型 书签 | 划线
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
