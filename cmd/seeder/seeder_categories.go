// Package main 提供分类数据填充功能
package main

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/utils"
	bookstoreModel "Qingyu_backend/models/bookstore"

	"go.mongodb.org/mongo-driver/bson"
)

type categorySeed struct {
	Name        string
	Description string
	SortOrder   int
}

var standardCategorySeeds = []categorySeed{
	{Name: "玄幻", Description: "东方玄幻、异世大陆、高武世界", SortOrder: 1},
	{Name: "都市", Description: "都市生活、都市异能、恋爱日常", SortOrder: 2},
	{Name: "仙侠", Description: "古典仙侠、现代修真、洪荒封神", SortOrder: 3},
	{Name: "科幻", Description: "未来世界、星际战争、时空穿梭", SortOrder: 4},
	{Name: "武侠", Description: "传统武侠、武侠幻想、国术无双", SortOrder: 5},
	{Name: "历史", Description: "架空历史、历史传记、两晋隋唐", SortOrder: 6},
	{Name: "游戏", Description: "虚拟网游、电子竞技、游戏异界", SortOrder: 7},
	{Name: "奇幻", Description: "西方奇幻、剑与魔法、黑暗幻想", SortOrder: 8},
}

// CategorySeeder 分类数据填充器
type CategorySeeder struct {
	db     *utils.Database
	config *config.Config
}

// NewCategorySeeder 创建分类填充器
func NewCategorySeeder(db *utils.Database, cfg *config.Config) *CategorySeeder {
	return &CategorySeeder{
		db:     db,
		config: cfg,
	}
}

// SeedCategories 填充标准分类数据。
func (s *CategorySeeder) SeedCategories() error {
	ctx := context.Background()

	count, err := s.db.Collection("categories").CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("统计分类数量失败: %w", err)
	}
	if count > 0 {
		fmt.Printf("分类集合已有 %d 条数据，跳过填充\n", count)
		return nil
	}

	now := time.Now()
	categories := make([]interface{}, 0, len(standardCategorySeeds))
	for _, item := range standardCategorySeeds {
		categories = append(categories, bookstoreModel.Category{
			Name:        item.Name,
			Description: item.Description,
			Icon:        "",
			Level:       0,
			SortOrder:   item.SortOrder,
			BookCount:   0,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	if _, err := s.db.Collection("categories").InsertMany(ctx, categories); err != nil {
		return fmt.Errorf("插入分类失败: %w", err)
	}

	fmt.Printf("成功插入 %d 个标准分类\n", len(categories))
	return nil
}

// GetActiveCategories 获取当前启用分类列表。
func (s *CategorySeeder) GetActiveCategories() ([]*bookstoreModel.Category, error) {
	ctx := context.Background()

	cursor, err := s.db.Collection("categories").Find(
		ctx,
		bson.M{"is_active": true},
	)
	if err != nil {
		return nil, fmt.Errorf("查询分类失败: %w", err)
	}
	defer cursor.Close(ctx)

	var categories []*bookstoreModel.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, fmt.Errorf("解析分类失败: %w", err)
	}

	return categories, nil
}

// Clean 清空分类数据。
func (s *CategorySeeder) Clean() error {
	ctx := context.Background()
	if _, err := s.db.Collection("categories").DeleteMany(ctx, bson.M{}); err != nil {
		return fmt.Errorf("清空 categories 集合失败: %w", err)
	}

	fmt.Println("已清空 categories 集合")
	return nil
}
