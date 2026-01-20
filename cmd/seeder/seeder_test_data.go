// Package main 提供测试数据填充功能
package main

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/cmd/seeder/models"
	"Qingyu_backend/cmd/seeder/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// TestDataSeeder 测试数据填充器
type TestDataSeeder struct {
	db *utils.Database
}

// NewTestDataSeeder 创建测试数据填充器
func NewTestDataSeeder(db *utils.Database) *TestDataSeeder {
	return &TestDataSeeder{db: db}
}

// SeedTestData 填充测试所需的数据
func (s *TestDataSeeder) SeedTestData() error {
	fmt.Println("开始填充测试数据...")

	if err := s.seedTestUser(); err != nil {
		return fmt.Errorf("填充测试用户失败: %w", err)
	}

	if err := s.seedTestBooks(); err != nil {
		return fmt.Errorf("填充测试书籍失败: %w", err)
	}

	if err := s.seedTestCategories(); err != nil {
		return fmt.Errorf("填充测试分类失败: %w", err)
	}

	fmt.Println("✅ 测试数据填充完成!")
	return nil
}

// seedTestUser 创建测试用户
func (s *TestDataSeeder) seedTestUser() error {
	ctx := context.Background()
	collection := s.db.Collection("users")

	now := time.Now()

	// 检查用户是否已存在
	count, _ := collection.CountDocuments(ctx, bson.M{"username": "testuser"})
	if count > 0 {
		fmt.Println("✓ 测试用户 testuser 已存在")
		return nil
	}

	// 创建测试用户（密码: 123456）
	// 注意：这是一个示例，实际应用中应该使用bcrypt等哈希算法
	// 为了简化，这里直接存储明文密码（仅用于测试环境）
	user := models.User{
		ID:        uuid.New().String(),
		Username:  "testuser",
		Email:     "testuser@qingyu.com",
		Password:  "123456", // ⚠️ 测试环境使用明文密码
		Role:      "user",
		Nickname:  "测试用户",
		Avatar:    "/images/avatars/default.png",
		Bio:       "这是一个测试账号",
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	fmt.Println("✓ 已创建测试用户: testuser/123456")
	return nil
}

// seedTestBooks 创建测试书籍
func (s *TestDataSeeder) seedTestBooks() error {
	ctx := context.Background()
	collection := s.db.Collection("books")

	// 检查是否已有测试书籍
	count, _ := collection.CountDocuments(ctx, bson.M{"title": bson.M{"$regex": "修仙"}})
	if count > 0 {
		fmt.Println("✓ 测试书籍已存在")
		return nil
	}

	now := time.Now()
	publishedAt := now.Add(-180 * 24 * time.Hour)

	books := []models.Book{
		{
			ID:           uuid.New().String(),
			Title:        "修仙世界",
			Author:       "飞升作者",
			Introduction: "一个普通少年，意外获得神秘传承，踏上修仙之路。历经千辛万苦，最终飞升成仙，成为一代传奇。",
			Cover:        "/images/covers/xiuxian_shijie.jpg",
			Categories:   []string{"玄幻", "修仙"},
			Tags:         []string{"修仙", "玄幻", "升级", "热血"},
			Status:       "ongoing",
			Rating:       8.5,
			RatingCount:  1250,
			ViewCount:    45000,
			WordCount:    1500000,
			ChapterCount: 450,
			Price:        0,
			IsFree:       true,
			IsRecommended: true,
			IsFeatured:   true,
			IsHot:        true,
			PublishedAt:  publishedAt,
			LastUpdateAt: now.Add(-24 * time.Hour),
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           uuid.New().String(),
			Title:        "修仙归来",
			Author:       "逍遥子",
			Introduction: "一代仙尊渡劫失败，重生回到地球。这一世，他要弥补所有遗憾，守护所爱之人，再登巅峰！",
			Cover:        "/images/covers/xiuxian Guilai.jpg",
			Categories:   []string{"玄幻", "修仙"},
			Tags:         []string{"修仙", "玄幻", "重生", "爽文"},
			Status:       "ongoing",
			Rating:       9.2,
			RatingCount:  8900,
			ViewCount:    120000,
			WordCount:    2800000,
			ChapterCount: 820,
			Price:        9.9,
			IsFree:       false,
			IsRecommended: true,
			IsFeatured:   true,
			IsHot:        true,
			PublishedAt:  publishedAt.Add(-30 * 24 * time.Hour),
			LastUpdateAt: now.Add(-12 * time.Hour),
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           uuid.New().String(),
			Title:        "修仙传说之无敌天下",
			Author:       "剑气纵横",
			Introduction: "天地不仁，以万物为刍狗。既然天道不公，那我便逆天而行，成就无上霸业！",
			Cover:        "/images/covers/xiuxian_chuanshuo.jpg",
			Categories:   []string{"玄幻", "修仙"},
			Tags:         []string{"修仙", "玄幻", "热血", "冒险"},
			Status:       "completed",
			Rating:       7.8,
			RatingCount:  560,
			ViewCount:    28000,
			WordCount:    980000,
			ChapterCount: 320,
			Price:        0,
			IsFree:       true,
			IsRecommended: false,
			IsFeatured:   false,
			IsHot:        false,
			PublishedAt:  publishedAt.Add(-90 * 24 * time.Hour),
			LastUpdateAt: now.Add(-60 * 24 * time.Hour),
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           uuid.New().String(),
			Title:        "万古修仙",
			Author:       "虚无居士",
			Introduction: "上古修仙界，强者如林。少年叶凡，偶得神秘小鼎，开启了一段波澜壮阔的修仙之旅。",
			Cover:        "/images/covers/wangu_xiuxian.jpg",
			Categories:   []string{"玄幻", "修仙"},
			Tags:         []string{"修仙", "玄幻", "冒险", "升级"},
			Status:       "ongoing",
			Rating:       8.8,
			RatingCount:  3200,
			ViewCount:    78000,
			WordCount:    1850000,
			ChapterCount: 580,
			Price:        19.9,
			IsFree:       false,
			IsRecommended: true,
			IsFeatured:   true,
			IsHot:        true,
			PublishedAt:  publishedAt.Add(-60 * 24 * time.Hour),
			LastUpdateAt: now.Add(-6 * time.Hour),
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           uuid.New().String(),
			Title:        "修仙从娃娃抓起",
			Author:       "童心未泯",
			Introduction: "穿越到修仙世界，发现自己竟然变成了婴儿。不过没关系，修仙就要从娃娃抓起！",
			Cover:        "/images/covers/xiuxian_wawa.jpg",
			Categories:   []string{"玄幻", "修仙"},
			Tags:         []string{"修仙", "玄幻", "搞笑", "轻松"},
			Status:       "ongoing",
			Rating:       7.2,
			RatingCount:  890,
			ViewCount:    15000,
			WordCount:    650000,
			ChapterCount: 210,
			Price:        0,
			IsFree:       true,
			IsRecommended: false,
			IsFeatured:   false,
			IsHot:        false,
			PublishedAt:  publishedAt.Add(-45 * 24 * time.Hour),
			LastUpdateAt: now.Add(-3 * 24 * time.Hour),
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	// 转换为 []interface{}
	booksInterface := make([]interface{}, len(books))
	for i := range books {
		booksInterface[i] = books[i]
	}

	_, err := collection.InsertMany(ctx, booksInterface)
	if err != nil {
		return err
	}

	fmt.Printf("✓ 已创建 %d 本修仙小说测试书籍\n", len(books))
	return nil
}

// seedTestCategories 创建测试分类
func (s *TestDataSeeder) seedTestCategories() error {
	ctx := context.Background()
	collection := s.db.Collection("categories")

	// 检查是否已有"玄幻"分类
	count, _ := collection.CountDocuments(ctx, bson.M{"name": "玄幻"})
	if count > 0 {
		fmt.Println("✓ 测试分类已存在")
		return nil
	}

	now := time.Now()

	categories := []interface{}{
		bson.M{
			"_id":        uuid.New().String(),
			"name":       "玄幻",
			"slug":       "xuanhuan",
			"description": "奇幻玄幻，想象力无限",
			"icon":       "/images/icons/xuanhuan.png",
			"parent_id":  nil,
			"sort_order": 1,
			"is_active":  true,
			"created_at": now,
			"updated_at": now,
		},
		bson.M{
			"_id":        uuid.New().String(),
			"name":       "修仙",
			"slug":       "xiuxian",
			"description": "修仙问道，长生不老",
			"icon":       "/images/icons/xiuxian.png",
			"parent_id":  nil,
			"sort_order": 2,
			"is_active":  true,
			"created_at": now,
			"updated_at": now,
		},
	}

	_, err := collection.InsertMany(ctx, categories)
	if err != nil {
		return err
	}

	fmt.Println("✓ 已创建测试分类: 玄幻、修仙")
	return nil
}
