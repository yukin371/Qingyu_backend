package reader

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnnotationType string

const (
	AnnotationTypeNote      AnnotationType = "note"      // 笔记
	AnnotationTypeBookmark  AnnotationType = "bookmark"  // 书签
	AnnotationTypeHighlight AnnotationType = "highlight" // 划线/高亮
)

type Annotation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`       // 用户ID
	BookID    primitive.ObjectID `bson:"book_id" json:"bookId"`       // 书籍ID
	ChapterID primitive.ObjectID `bson:"chapter_id" json:"chapterId"` // 章节ID
	Range     string             `bson:"range" json:"range"`          // 标注范围：start-end
	Text      string             `bson:"text" json:"text"`            // 标注文本
	Note      string             `bson:"note" json:"note"`            // 注释
	Type      string             `bson:"type" json:"type"`            // 标注类型 书签 | 划线
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}
