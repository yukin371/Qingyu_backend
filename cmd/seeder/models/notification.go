// Package models 提供通知相关数据模型
package models

import "time"

// Notification 用户通知
type Notification struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	Type      string    `bson:"type"`
	Title     string    `bson:"title"`
	Content   string    `bson:"content"`
	Link      string    `bson:"link"`
	IsRead    bool      `bson:"is_read"`
	ReadAt    time.Time `bson:"read_at"`
	CreatedAt time.Time `bson:"created_at"`
}
