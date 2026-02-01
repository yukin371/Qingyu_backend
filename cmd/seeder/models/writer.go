// Package models 提供写作相关数据模型
package models

import "time"

// WriterProject 写作项目
type WriterProject struct {
	ID          string    `bson:"_id"`
	UserID      string    `bson:"user_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Genre       string    `bson:"genre"`
	Status      string    `bson:"status"`
	WordCount   int       `bson:"word_count"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

// WriterDocument 写作文档
type WriterDocument struct {
	ID        string    `bson:"_id"`
	ProjectID string    `bson:"project_id"`
	UserID    string    `bson:"user_id"`
	Title     string    `bson:"title"`
	Content   string    `bson:"content"`
	WordCount int       `bson:"word_count"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
