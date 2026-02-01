// Package models 提供统计相关数据模型
package models

import "time"

// BookStats 书籍统计
type BookStats struct {
	ID             string    `bson:"_id"`
	BookID         string    `bson:"book_id"`
	ViewCount      int       `bson:"view_count"`
	ReadCount      int       `bson:"read_count"`
	FavoriteCount  int       `bson:"favorite_count"`
	ShareCount     int       `bson:"share_count"`
	AvgReadTime    int       `bson:"avg_read_time"`
	CompletionRate float64   `bson:"completion_rate"`
	Date           string    `bson:"date"`
	CreatedAt      time.Time `bson:"created_at"`
}

// ChapterStats 章节统计
type ChapterStats struct {
	ID        string    `bson:"_id"`
	ChapterID string    `bson:"chapter_id"`
	BookID    string    `bson:"book_id"`
	ViewCount int       `bson:"view_count"`
	ReadCount int       `bson:"read_count"`
	StayTime  int       `bson:"stay_time"`
	Date      string    `bson:"date"`
	CreatedAt time.Time `bson:"created_at"`
}
