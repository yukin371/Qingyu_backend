package seeds

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Category 分类模型
type Category struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	ParentID    *primitive.ObjectID `bson:"parent_id,omitempty"`
	Sort        int                `bson:"sort"`
	IsActive    bool               `bson:"is_active"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// SeedCategories 分类种子数据
func SeedCategories(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("categories")

	// 检查是否已有数据
	count, err := collection.CountDocuments(ctx, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("failed to count categories: %w", err)
	}

	if count > 0 {
		fmt.Printf("Categories collection already has %d documents, skipping seed\n", count)
		return nil
	}

	now := time.Now()

	// 准备分类数据
	categories := []Category{
		{
			Name:        "玄幻",
			Description: "东方玄幻、异世大陆、高武世界",
			Sort:        1,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "都市",
			Description: "都市生活、都市异能、恋爱日常",
			Sort:        2,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "仙侠",
			Description: "古典仙侠、现代修真、洪荒封神",
			Sort:        3,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "科幻",
			Description: "未来世界、星际战争、时空穿梭",
			Sort:        4,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "武侠",
			Description: "传统武侠、武侠幻想、国术无双",
			Sort:        5,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "历史",
			Description: "架空历史、历史传记、两晋隋唐",
			Sort:        6,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "游戏",
			Description: "虚拟网游、电子竞技、游戏异界",
			Sort:        7,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "奇幻",
			Description: "西方奇幻、剑与魔法、黑暗幻想",
			Sort:        8,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	// 插入分类
	docs := make([]interface{}, len(categories))
	for i, category := range categories {
		docs[i] = category
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to insert categories: %w", err)
	}

	fmt.Printf("✓ Seeded %d categories\n", len(result.InsertedIDs))
	for i, id := range result.InsertedIDs {
		fmt.Printf("  - %s (ID: %v)\n", categories[i].Name, id)
	}

	return nil
}







