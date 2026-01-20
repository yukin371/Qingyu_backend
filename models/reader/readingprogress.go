package reader

import (
	"time"
)

// ReadingProgress 阅读进度

type ReadingProgress struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	UserID      string    `bson:"user_id" json:"userId"`           // 用户ID
	BookID      string    `bson:"book_id" json:"bookId"`           // 书籍ID
	ChapterID   string    `bson:"chapter_id" json:"chapterId"`     // 章节ID
	Progress    float64   `bson:"progress" json:"progress"`        // 进度：0-1之间的小数
	ReadingTime int64     `bson:"reading_time" json:"readingTime"` // 阅读时间（秒）
	LastReadAt  time.Time `bson:"last_read_at" json:"lastReadAt"`  // 最后阅读时间
	Status      string    `bson:"status" json:"status"`            // 书籍状态: reading(在读), want_read(想读), finished(读完)
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}
