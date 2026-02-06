// Package models 定义数据模型
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
	UserStatusDeleted  UserStatus = "deleted"
)

// User 用户模型
type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"` // 不序列化到 JSON
	Roles     []string           `bson:"roles" json:"roles"`
	Status    UserStatus         `bson:"status" json:"status"`
	Nickname  string             `bson:"nickname" json:"nickname"`
	Bio       string             `bson:"bio" json:"bio"`
	Avatar    string             `bson:"avatar" json:"avatar"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
