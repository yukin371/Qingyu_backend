package bookstore

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookStatus 书籍状态
type BookStatus string

const (
	BookStatusOngoing   BookStatus = "ongoing"   // 连载中
	BookStatusCompleted BookStatus = "completed" // 已完结
	BookStatusPaused    BookStatus = "paused"    // 暂停更新
)

// BookDetail 书籍详情模型
type BookDetail struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title        string             `bson:"title" json:"title"`
	Subtitle     string             `bson:"subtitle" json:"subtitle"`
	Author       string             `bson:"author" json:"author"`
	AuthorID     primitive.ObjectID `bson:"author_id" json:"author_id"`
	Description  string             `bson:"description" json:"description"`
	CoverURL     string             `bson:"cover_url" json:"cover_url"`
	Publisher    string             `bson:"publisher" json:"publisher"`
	PublishDate  time.Time          `bson:"publish_date" json:"publish_date"`
	ISBN         string             `bson:"isbn" json:"isbn"`
	Categories   []string           `bson:"categories" json:"categories"`
	Tags         []string           `bson:"tags" json:"tags"`
	Status       BookStatus         `bson:"status" json:"status"`
	WordCount    int64              `bson:"word_count" json:"word_count"`
	ChapterCount int                `bson:"chapter_count" json:"chapter_count"`
	Price        float64            `bson:"price" json:"price"`
	IsFree       bool               `bson:"is_free" json:"is_free"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// BeforeCreate 在创建前设置时间戳
func (bd *BookDetail) BeforeCreate() {
	now := time.Now()
	bd.CreatedAt = now
	bd.UpdatedAt = now
}

// BeforeUpdate 在更新前刷新更新时间戳
func (bd *BookDetail) BeforeUpdate() {
	bd.UpdatedAt = time.Now()
}

// IsCompleted 检查书籍是否已完结
func (bd *BookDetail) IsCompleted() bool {
	return bd.Status == BookStatusCompleted
}

// IsOngoing 检查书籍是否连载中
func (bd *BookDetail) IsOngoing() bool {
	return bd.Status == BookStatusOngoing
}

// IsPaused 检查书籍是否暂停更新
func (bd *BookDetail) IsPaused() bool {
	return bd.Status == BookStatusPaused
}