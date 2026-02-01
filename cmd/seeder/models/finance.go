// Package models 提供财务相关数据模型
package models

import "time"

// AuthorRevenue 作者收入
type AuthorRevenue struct {
	ID        string    `bson:"_id" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	BookID    string    `bson:"book_id" json:"book_id"`
	Amount    float64   `bson:"amount" json:"amount"`
	Type      string    `bson:"type" json:"type"`
	Status    string    `bson:"status" json:"status"`
	Period    string    `bson:"period" json:"period"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	SettledAt time.Time `bson:"settled_at" json:"settled_at"`
}

// Membership 会员
type Membership struct {
	ID        string    `bson:"_id" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	Type      string    `bson:"type" json:"type"`
	Status    string    `bson:"status" json:"status"`
	StartAt   time.Time `bson:"start_at" json:"start_at"`
	ExpireAt  time.Time `bson:"expire_at" json:"expire_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
