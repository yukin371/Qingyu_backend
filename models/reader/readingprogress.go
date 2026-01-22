package reader

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/shared/types"
)

// ReadingProgress 阅读进度

type ReadingProgress struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserID      primitive.ObjectID `bson:"user_id" json:"userId"`           // 用户ID
	BookID      primitive.ObjectID `bson:"book_id" json:"bookId"`           // 书籍ID
	ChapterID   primitive.ObjectID `bson:"chapter_id" json:"chapterId"`     // 章节ID
	Progress    types.Progress     `bson:"progress" json:"progress"`        // 进度：0-1之间的小数
	ReadingTime int64              `bson:"reading_time" json:"readingTime"` // 阅读时间（秒）
	LastReadAt  time.Time          `bson:"last_read_at" json:"lastReadAt"`  // 最后阅读时间
	Status      string             `bson:"status" json:"status"`            // 书籍状态: reading(在读), want_read(想读), finished(读完)
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}
