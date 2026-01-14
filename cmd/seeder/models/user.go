// Package models 定义数据模型
package models

import "time"

// User 用户模型
type User struct {
	ID        string    `bson:"_id" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"-"` // 不序列化到 JSON
	Role      string    `bson:"role" json:"role"`
	Nickname  string    `bson:"nickname" json:"nickname"`
	Bio       string    `bson:"bio" json:"bio"`
	Avatar    string    `bson:"avatar" json:"avatar"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
