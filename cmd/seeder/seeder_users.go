// Package main 提供用户数据填充功能
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/generators"
	"Qingyu_backend/cmd/seeder/models"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
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

// SeedRealUsers 填充真实测试用户
func (s *UserSeeder) SeedRealUsers() error {
	// 读取真实用户数据
	data, err := os.ReadFile("data/users.json")
	if err != nil {
		return fmt.Errorf("读取 users.json 失败: %w", err)
	}

	var testData struct {
		TestUsers []map[string]interface{} `json:"test_users"`
	}

	if err := json.Unmarshal(data, &testData); err != nil {
		return fmt.Errorf("解析 users.json 失败: %w", err)
	}

	// 转换为 User 对象
	realUsers := make([]models.User, len(testData.TestUsers))
	for i, u := range testData.TestUsers {
		realUsers[i] = models.User{
			ID:        u["username"].(string), // 使用用户名作为 ID
			Username:  u["username"].(string),
			Email:     u["email"].(string),
			Password:  u["password"].(string),
			Role:      u["role"].(string),
			Nickname:  u["nickname"].(string),
			Bio:       u["bio"].(string),
			Avatar:    "/images/avatars/default.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	return s.inserter.InsertMany(context.Background(), realUsers)
}

// SeedGeneratedUsers 填充生成的用户
func (s *UserSeeder) SeedGeneratedUsers() error {
	scale := config.GetScaleConfig(s.config.Scale)

	var allUsers []models.User

	// 计算各角色用户数量，确保不出现负数
	normalCount := scale.Users - scale.Authors - 100
	if normalCount < 0 {
		normalCount = 0
	}
	
	// 生成普通用户
	if normalCount > 0 {
		normalUsers := s.gen.GenerateUsers(normalCount, "reader")
		allUsers = append(allUsers, normalUsers...)
	}

	// 生成作者用户
	authors := s.gen.GenerateUsers(scale.Authors, "author")
	allUsers = append(allUsers, authors...)

	// 生成 VIP 用户
	vipCount := 100
	if scale.Users < scale.Authors+vipCount {
		vipCount = scale.Users - scale.Authors
		if vipCount < 0 {
			vipCount = 0
		}
	}
	if vipCount > 0 {
		vipUsers := s.gen.GenerateUsers(vipCount, "vip")
		allUsers = append(allUsers, vipUsers...)
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
