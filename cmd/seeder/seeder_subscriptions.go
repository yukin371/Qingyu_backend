// Package main 提供书籍订阅关系填充功能
package main

import (
	"context"
	"fmt"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/relationships"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// SubscriptionSeeder 书籍订阅关系填充器
type SubscriptionSeeder struct {
	db      *utils.Database
	config  *config.Config
	builder *relationships.RelationshipBuilder
}

// NewSubscriptionSeeder 创建订阅关系填充器
func NewSubscriptionSeeder(db *utils.Database, cfg *config.Config) *SubscriptionSeeder {
	return &SubscriptionSeeder{
		db:      db,
		config:  cfg,
		builder: relationships.NewRelationshipBuilder(db),
	}
}

// SeedSubscriptions 填充用户订阅书籍关系。
func (s *SubscriptionSeeder) SeedSubscriptions() error {
	return s.builder.BuildSubscriptions()
}

// Clean 清空订阅关系。
func (s *SubscriptionSeeder) Clean() error {
	if _, err := s.db.Collection("subscriptions").DeleteMany(context.Background(), bson.M{}); err != nil {
		return fmt.Errorf("清空 subscriptions 集合失败: %w", err)
	}

	fmt.Println("已清空 subscriptions 集合")
	return nil
}
