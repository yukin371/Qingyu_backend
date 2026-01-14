// Package generators 提供数据生成器
package generators

import (
	"time"

	"github.com/google/uuid"
	"Qingyu_backend/cmd/seeder/models"
)

// UserGenerator 用户数据生成器
type UserGenerator struct {
	*BaseGenerator
}

// NewUserGenerator 创建用户生成器
func NewUserGenerator() *UserGenerator {
	return &UserGenerator{
		BaseGenerator: NewBaseGenerator(),
	}
}

// GenerateUser 生成单个用户
func (g *UserGenerator) GenerateUser(role string) models.User {
	username := g.Username(role)
	now := time.Now()

	return models.User{
		ID:       uuid.New().String(),
		Username: username,
		Email:    g.Email(),
		// bcrypt hash of "password"
		Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
		Role:     role,
		Nickname: username,
		Avatar:   "/images/avatars/default.png",
		Bio:      g.faker.Paragraph(1, 3, 20, " "),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// GenerateUsers 批量生成用户
func (g *UserGenerator) GenerateUsers(count int, role string) []models.User {
	users := make([]models.User, count)
	for i := 0; i < count; i++ {
		users[i] = g.GenerateUser(role)
	}
	return users
}
