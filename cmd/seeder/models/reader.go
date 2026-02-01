// Package models 提供阅读相关数据模型
package models

import "time"

// ReadingHistory 阅读历史记录
type ReadingHistory struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	BookID    string    `bson:"book_id"`
	ChapterID string    `bson:"chapter_id"`
	ReadTime  time.Time `bson:"read_time"`
	Duration  int       `bson:"duration"`
	Device    string    `bson:"device"`
	CreatedAt time.Time `bson:"created_at"`
}

// ReadingProgress 阅读进度
type ReadingProgress struct {
	ID            string    `bson:"_id"`
	UserID        string    `bson:"user_id"`
	BookID        string    `bson:"book_id"`
	ChapterNum    int       `bson:"chapter_num"`
	Progress      float64   `bson:"progress"`
	LastReadAt    time.Time `bson:"last_read_at"`
	TotalReadTime int       `bson:"total_read_time"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

// Bookmark 书签
type Bookmark struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	BookID    string    `bson:"book_id"`
	ChapterID string    `bson:"chapter_id"`
	Position  int       `bson:"position"`
	Note      string    `bson:"note"`
	CreatedAt time.Time `bson:"created_at"`
}

// Annotation 批注
type Annotation struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	BookID    string    `bson:"book_id"`
	ChapterID string    `bson:"chapter_id"`
	Position  int       `bson:"position"`
	Content   string    `bson:"content"`
	Color     string    `bson:"color"`
	CreatedAt time.Time `bson:"created_at"`
}
