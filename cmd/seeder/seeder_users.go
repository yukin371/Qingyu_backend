// Package main 提供用户数据填充功能
package main

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/generators"
	"Qingyu_backend/cmd/seeder/models"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserSeeder 用户数据填充器
type UserSeeder struct {
	db       *utils.Database
	config   *config.Config
	gen      *generators.UserGenerator
	inserter *utils.BulkInserter
}

// NewUserSeeder 创建用户填充器
func NewUserSeeder(db *utils.Database, cfg *config.Config) *UserSeeder {
	collection := db.Collection("users")
	return &UserSeeder{
		db:     db,
		config: cfg,
		gen:    generators.NewUserGenerator(),
		inserter: utils.NewBulkInserter(collection, cfg.BatchSize),
	}
}

// SeedRealUsers 填充固定测试账号（用于E2E测试）
func (s *UserSeeder) SeedRealUsers() error {
	now := time.Now()
	defaultPassword := utils.DefaultPasswordHash()

	// 固定测试账号列表 - 使用 ObjectID 而不是字符串 ID
	testUsers := []models.User{
		{
			ID:        primitive.NewObjectID(),
			Username:  "testuser001",
			Password:  defaultPassword, // bcrypt加密的 "password"
			Roles:     []string{generators.RoleReader},
			Status:    models.UserStatusActive,
			Email:     "testuser001@example.com",
			Nickname:  "测试读者001",
			Bio:       "E2E测试用读者账号，密码：password",
			Avatar:    "/images/avatars/default.png",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        primitive.NewObjectID(),
			Username:  "testauthor001",
			Password:  defaultPassword, // bcrypt加密的 "password"
			Roles:     []string{generators.RoleAuthor},
			Status:    models.UserStatusActive,
			Email:     "testauthor001@example.com",
			Nickname:  "测试作者001",
			Bio:       "E2E测试用作者账号，密码：password",
			Avatar:    "/images/avatars/default.png",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        primitive.NewObjectID(),
			Username:  "testadmin001",
			Password:  defaultPassword, // bcrypt加密的 "password"
			Roles:     []string{generators.RoleAdmin},
			Status:    models.UserStatusActive,
			Email:     "testadmin001@example.com",
			Nickname:  "测试管理员001",
			Bio:       "E2E测试用管理员账号，密码：password",
			Avatar:    "/images/avatars/default.png",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// 尝试读取并合并原有的 users.json 数据（向后兼容）
	if data, err := os.ReadFile("data/users.json"); err == nil {
		var testData struct {
			TestUsers []map[string]interface{} `json:"test_users"`
		}
		if json.Unmarshal(data, &testData) == nil {
			for _, u := range testData.TestUsers {
				testUsers = append(testUsers, models.User{
					ID:        primitive.NewObjectID(),
					Username:  u["username"].(string),
					Email:     u["email"].(string),
					Password:  u["password"].(string),
					Roles:     []string{u["role"].(string)},
					Status:    models.UserStatusActive,
					Nickname:  u["nickname"].(string),
					Bio:       u["bio"].(string),
					Avatar:    "/images/avatars/default.png",
					CreatedAt: now,
					UpdatedAt: now,
				})
			}
		}
	}

	return s.inserter.InsertMany(context.Background(), testUsers)
}

// SeedGeneratedUsers 填充生成的用户
func (s *UserSeeder) SeedGeneratedUsers() error {
	scale := config.GetScaleConfig(s.config.Scale)
	totalUsers := scale.Users

	// 按比例计算各角色数量
	readerCount := int(math.Round(float64(totalUsers) * scale.ReaderPercent))
	authorCount := int(math.Round(float64(totalUsers) * scale.AuthorPercent))
	adminCount := totalUsers - readerCount - authorCount // 确保总数正确

	var allUsers []models.User

	// 生成 reader 角色用户
	if readerCount > 0 {
		readerUsers := s.gen.GenerateUsers(readerCount, generators.RoleReader)
		allUsers = append(allUsers, readerUsers...)
	}

	// 生成 author 角色用户
	if authorCount > 0 {
		authorUsers := s.gen.GenerateUsers(authorCount, generators.RoleAuthor)
		allUsers = append(allUsers, authorUsers...)
	}

	// 生成 admin 角色用户
	if adminCount > 0 {
		adminUsers := s.gen.GenerateUsers(adminCount, generators.RoleAdmin)
		allUsers = append(allUsers, adminUsers...)
	}

	return s.inserter.InsertMany(context.Background(), allUsers)
}

// Clean 清空用户数据
func (s *UserSeeder) Clean() error {
	ctx := context.Background()
	_, err := s.db.Collection("users").DeleteMany(ctx, bson.M{})
	return err
}

// Count 统计用户数量
func (s *UserSeeder) Count() (int64, error) {
	return s.inserter.CountDocuments(context.Background(), bson.M{})
}

// GetGenerator 获取用户生成器
func (s *UserSeeder) GetGenerator() *generators.UserGenerator {
	return s.gen
}

// GetInserter 获取批量插入器
func (s *UserSeeder) GetInserter() *utils.BulkInserter {
	return s.inserter
}
