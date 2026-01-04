package reader

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Bookmark 书签
type Bookmark struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"userId"`
	BookID      primitive.ObjectID `bson:"book_id" json:"bookId"`
	ChapterID   primitive.ObjectID `bson:"chapter_id" json:"chapterId"`
	Position    int               `bson:"position" json:"position"`      // 字符位置
	Note        string            `bson:"note" json:"note"`               // 书签笔记
	Color       string            `bson:"color" json:"color"`             // 书签颜色（默认yellow）
	Quote       string            `bson:"quote,omitempty" json:"quote,omitempty"` // 引用的文本
	IsPublic    bool              `bson:"is_public" json:"isPublic"`      // 是否公开
	Tags        []string          `bson:"tags,omitempty" json:"tags,omitempty"` // 标签
	CreatedAt   time.Time         `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updatedAt"`
}

// BookmarkExport 书签导出格式
type BookmarkExport struct {
	BookTitle    string    `json:"bookTitle"`
	ChapterTitle string    `json:"chapterTitle"`
	Position     int       `json:"position"`
	Note         string    `json:"note"`
	Color        string    `json:"color"`
	Quote        string    `json:"quote"`
	CreatedAt    time.Time `json:"createdAt"`
}

// BookmarkFilter 书签筛选条件
type BookmarkFilter struct {
	BookID    *primitive.ObjectID `bson:"book_id,omitempty"`
	ChapterID *primitive.ObjectID `bson:"chapter_id,omitempty"`
	Color     string              `bson:"color,omitempty"`
	Tag       string              `bson:"tags,omitempty"`
	IsPublic  *bool               `bson:"is_public,omitempty"`
	DateFrom  *time.Time          `bson:"-"`
	DateTo    *time.Time          `bson:"-"`
}

// BookmarkStats 书签统计
type BookmarkStats struct {
	TotalCount      int64            `json:"totalCount"`
	ByColor         map[string]int64 `json:"byColor"`
	ByBook          map[string]int64 `json:"byBook"`
	PublicCount     int64            `json:"publicCount"`
	PrivateCount    int64            `json:"privateCount"`
	ThisMonthCount  int64            `json:"thisMonthCount"`
	ThisWeekCount   int64            `json:"thisWeekCount"`
	RecentBookmarks []Bookmark       `json:"recentBookmarks"`
}
