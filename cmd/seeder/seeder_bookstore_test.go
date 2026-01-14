// Package main 提供书城数据填充功能测试
package main

import (
	"testing"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/models"
)

// TestBookstoreSeeder_SeedsGeneratedBooks 测试书籍生成功能
func TestBookstoreSeeder_SeedsGeneratedBooks(t *testing.T) {
	// 创建测试配置
	cfg := &config.Config{
		Scale:     "small", // 使用小规模测试
		BatchSize: 10,
	}

	// 创建 BookstoreSeeder
	seeder := NewBookstoreSeeder(nil, cfg)

	// 测试分类比例计算
	scale := config.GetScaleConfig(cfg.Scale)
	totalBooks := scale.Books

	categoryRatios := map[string]float64{
		"仙侠": 0.30,
		"都市": 0.25,
		"科幻": 0.20,
		"历史": 0.15,
		"其他": 0.10,
	}

	var totalCount int
	for _, ratio := range categoryRatios {
		count := int(float64(totalBooks) * ratio)
		totalCount += count
		t.Logf("分类比例: %.2f%%, 书籍数量: %d", ratio*100, count)
	}

	t.Logf("总书籍数: %d (期望: %d)", totalCount, totalBooks)

	// 验证生成器能正确生成书籍
	gen := seeder.gen
	testBooks := gen.GenerateBooks(5, "仙侠")

	if len(testBooks) != 5 {
		t.Errorf("期望生成 5 本书，实际生成 %d 本", len(testBooks))
	}

	// 验证书籍属性
	for i, book := range testBooks {
		if book.ID == "" {
			t.Errorf("第 %d 本书缺少 ID", i)
		}
		if book.Title == "" {
			t.Errorf("第 %d 本书缺少标题", i)
		}
		if book.Author == "" {
			t.Errorf("第 %d 本书缺少作者", i)
		}
		if len(book.Categories) == 0 {
			t.Errorf("第 %d 本书缺少分类", i)
		}
		if book.Categories[0] != "仙侠" {
			t.Errorf("第 %d 本书分类错误，期望 '仙侠'，实际 '%s'", i, book.Categories[0])
		}
		t.Logf("书籍 %d: %s by %s (评分: %.1f)", i+1, book.Title, book.Author, book.Rating)
	}
}

// TestBookstoreSeeder_VerifyCategoryDistribution 测试分类分布
func TestBookstoreSeeder_VerifyCategoryDistribution(t *testing.T) {
	cfg := &config.Config{
		Scale:     "medium",
		BatchSize: 100,
	}

	scale := config.GetScaleConfig(cfg.Scale)
	totalBooks := scale.Books

	categoryRatios := map[string]float64{
		"仙侠": 0.30,
		"都市": 0.25,
		"科幻": 0.20,
		"历史": 0.15,
		"其他": 0.10,
	}

	t.Logf("测试规模: %s (总书籍数: %d)", cfg.Scale, totalBooks)

	for category, ratio := range categoryRatios {
		count := int(float64(totalBooks) * ratio)
		expectedMin := int(float64(totalBooks)*ratio - 1)
		expectedMax := int(float64(totalBooks)*ratio + 1)

		if count < expectedMin || count > expectedMax {
			t.Errorf("%s 类书籍数量 %d 超出预期范围 [%d, %d]", category, count, expectedMin, expectedMax)
		}

		t.Logf("%s: %d 本 (%.1f%%)", category, count, ratio*100)
	}
}

// TestBookstoreSeeder_BookStructure 测试书籍结构
func TestBookstoreSeeder_BookStructure(t *testing.T) {
	gen := NewBookstoreSeeder(nil, &config.Config{}).gen

	// 测试每个分类的书籍生成
	categories := []string{"仙侠", "都市", "科幻", "历史", "其他"}

	for _, category := range categories {
		book := gen.GenerateBook(category)

		// 验证必填字段
		if book.ID == "" {
			t.Errorf("[%s] 缺少 ID", category)
		}
		if book.Title == "" {
			t.Errorf("[%s] 缺少标题", category)
		}
		if book.Author == "" {
			t.Errorf("[%s] 缺少作者", category)
		}
		if book.Introduction == "" {
			t.Errorf("[%s] 缺少简介", category)
		}
		if book.Cover == "" {
			t.Errorf("[%s] 缺少封面", category)
		}
		if len(book.Categories) == 0 {
			t.Errorf("[%s] 缺少分类", category)
		}
		if len(book.Tags) == 0 {
			t.Errorf("[%s] 缺少标签", category)
		}
		if book.Status == "" {
			t.Errorf("[%s] 缺少状态", category)
		}
		if book.Rating < 0 || book.Rating > 10 {
			t.Errorf("[%s] 评分超出范围 [0, 10]: %.1f", category, book.Rating)
		}
		if book.WordCount <= 0 {
			t.Errorf("[%s] 字数必须大于 0: %d", category, book.WordCount)
		}
		if book.ChapterCount <= 0 {
			t.Errorf("[%s] 章节数必须大于 0: %d", category, book.ChapterCount)
		}
		if book.Price < 0 {
			t.Errorf("[%s] 价格不能为负: %.2f", category, book.Price)
		}
		if book.PublishedAt.IsZero() {
			t.Errorf("[%s] 缺少发布时间", category)
		}
		if book.LastUpdateAt.IsZero() {
			t.Errorf("[%s] 缺少更新时间", category)
		}
		if book.CreatedAt.IsZero() {
			t.Errorf("[%s] 缺少创建时间", category)
		}
		if book.UpdatedAt.IsZero() {
			t.Errorf("[%s] 缺少更新时间", category)
		}

		// 验证时间逻辑
		if book.LastUpdateAt.Before(book.PublishedAt) {
			t.Errorf("[%s] 更新时间早于发布时间", category)
		}

		t.Logf("[%s] 书籍结构验证通过: %s", category, book.Title)
	}
}

// TestBookstoreSeeder_BannerStructure 测试 Banner 结构
func TestBookstoreSeeder_BannerStructure(t *testing.T) {
	now := time.Now()

	// 模拟生成 banners
	banners := []map[string]interface{}{
		{
			"_id":         "banner_new_books",
			"title":       "新书推荐",
			"description": "最新上架的精品好书",
			"image":       "/images/banners/new_books.jpg",
			"link":        "/books/new",
			"type":        "new_books",
			"priority":    1,
			"is_active":   true,
			"start_at":    now,
			"end_at":      now.Add(30 * 24 * time.Hour),
			"created_at":  now,
			"updated_at":  now,
		},
		{
			"_id":         "banner_free_books",
			"title":       "限时免费",
			"description": "限时免费阅读热门作品",
			"image":       "/images/banners/free_books.jpg",
			"link":        "/books/free",
			"type":        "free_books",
			"priority":    2,
			"is_active":   true,
			"start_at":    now,
			"end_at":      now.Add(7 * 24 * time.Hour),
			"created_at":  now,
			"updated_at":  now,
		},
	}

	// 验证 banner 结构
	if len(banners) != 2 {
		t.Errorf("期望 2 个 banner，实际 %d 个", len(banners))
	}

	for i, banner := range banners {
		// 验证必填字段
		requiredFields := []string{"_id", "title", "description", "image", "link", "type", "priority", "is_active"}
		for _, field := range requiredFields {
			if _, ok := banner[field]; !ok {
				t.Errorf("Banner %d 缺少必填字段: %s", i, field)
			}
		}

		// 验证类型
		if priority, ok := banner["priority"].(int); !ok || priority <= 0 {
			t.Errorf("Banner %d priority 必须是正整数", i)
		}

		if isActive, ok := banner["is_active"].(bool); !ok {
			t.Errorf("Banner %d is_active 必须是布尔值", i)
		}

		t.Logf("Banner %d: %s (优先级: %d)", i+1, banner["title"], banner["priority"])
	}
}
