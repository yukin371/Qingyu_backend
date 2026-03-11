// Package main 提供书城数据填充功能
package main

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/generators"
	"Qingyu_backend/cmd/seeder/models"
	"Qingyu_backend/cmd/seeder/utils"
	bookstoreModel "Qingyu_backend/models/bookstore"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type categoryDistribution struct {
	Name  string
	Ratio float64
}

var defaultCategoryDistributions = []categoryDistribution{
	{Name: "仙侠", Ratio: 0.25},
	{Name: "都市", Ratio: 0.20},
	{Name: "科幻", Ratio: 0.15},
	{Name: "历史", Ratio: 0.10},
	{Name: "玄幻", Ratio: 0.10},
	{Name: "武侠", Ratio: 0.08},
	{Name: "游戏", Ratio: 0.07},
	{Name: "奇幻", Ratio: 0.05},
}

// BookstoreSeeder 书城数据填充器
type BookstoreSeeder struct {
	db       *utils.Database
	config   *config.Config
	gen      *generators.BookGenerator
	inserter *utils.BulkInserter
}

// NewBookstoreSeeder 创建书城填充器
func NewBookstoreSeeder(db *utils.Database, cfg *config.Config) *BookstoreSeeder {
	collection := db.Collection("books")
	return &BookstoreSeeder{
		db:       db,
		config:   cfg,
		gen:      generators.NewBookGenerator(),
		inserter: utils.NewBulkInserter(collection, cfg.BatchSize),
	}
}

// SeedGeneratedBooks 填充生成的书籍数据
func (s *BookstoreSeeder) SeedGeneratedBooks() error {
	// 获取配置的规模
	scale := config.GetScaleConfig(s.config.Scale)
	totalBooks := scale.Books

	// 1. 获取真实的author用户
	fmt.Println("正在获取author角色用户...")
	authorIDs, err := s.getAuthorUsers()
	if err != nil {
		return fmt.Errorf("获取author用户失败: %w", err)
	}

	if len(authorIDs) == 0 {
		return fmt.Errorf("没有找到author角色的用户，请先运行用户填充命令")
	}

	fmt.Printf("找到 %d 个author用户\n", len(authorIDs))

	// 2. 读取真实分类，建立标准分类映射
	categories, err := s.getStandardCategories()
	if err != nil {
		return fmt.Errorf("获取分类失败: %w", err)
	}
	if len(categories) == 0 {
		return fmt.Errorf("没有找到分类数据，请先运行分类填充命令")
	}

	// 3. 存储所有生成的书籍
	var allBooks []models.Book

	remaining := totalBooks
	for i, distribution := range defaultCategoryDistributions {
		category, ok := categories[distribution.Name]
		if !ok {
			return fmt.Errorf("缺少标准分类: %s", distribution.Name)
		}

		count := int(float64(totalBooks) * distribution.Ratio)
		if i == len(defaultCategoryDistributions)-1 {
			count = remaining
		}
		remaining -= count
		if count <= 0 {
			continue
		}

		// 生成该分类的书籍，使用真实分类和 author ID
		books := s.gen.GenerateBooksFromCategory(count, category, authorIDs)
		allBooks = append(allBooks, books...)

		fmt.Printf("已生成 %d 本%s类书籍\n", count, category.Name)
	}

	// 5. 批量插入所有书籍
	ctx := context.Background()
	if err := s.inserter.InsertMany(ctx, allBooks); err != nil {
		return fmt.Errorf("插入书籍失败: %w", err)
	}

	fmt.Printf("成功插入 %d 本书籍\n", len(allBooks))

	// 6. 输出author分配统计
	s.printAuthorDistributionStats(allBooks, authorIDs)

	return nil
}

func (s *BookstoreSeeder) getStandardCategories() (map[string]*bookstoreModel.Category, error) {
	ctx := context.Background()
	cursor, err := s.db.Collection("categories").Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, fmt.Errorf("查询分类失败: %w", err)
	}
	defer cursor.Close(ctx)

	var categories []*bookstoreModel.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, fmt.Errorf("解析分类失败: %w", err)
	}

	result := make(map[string]*bookstoreModel.Category, len(categories))
	for _, category := range categories {
		result[category.Name] = category
	}
	return result, nil
}

// getAuthorUsers 获取author角色的用户ID列表
func (s *BookstoreSeeder) getAuthorUsers() ([]primitive.ObjectID, error) {
	ctx := context.Background()

	// 查询roles包含author的用户
	cursor, err := s.db.Collection("users").Find(ctx, bson.M{"roles": bson.M{"$in": []string{"author"}}})
	if err != nil {
		return nil, fmt.Errorf("查询author用户失败: %w", err)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var users []struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("解析author用户失败: %w", err)
	}

	// 提取ID列表
	authorIDs := make([]primitive.ObjectID, len(users))
	for i, u := range users {
		authorIDs[i] = u.ID
	}

	return authorIDs, nil
}

// printAuthorDistributionStats 打印author分配统计信息
func (s *BookstoreSeeder) printAuthorDistributionStats(books []models.Book, authorIDs []primitive.ObjectID) {
	// 统计每个author的书籍数量
	authorBookCount := make(map[primitive.ObjectID]int)
	for _, book := range books {
		if !book.AuthorID.IsZero() {
			authorBookCount[book.AuthorID]++
		}
	}

	// 计算统计信息
	totalBooks := len(books)
	totalAuthors := len(authorIDs)
	avgBooksPerAuthor := float64(totalBooks) / float64(totalAuthors)
	minBooks := totalBooks
	maxBooks := 0

	for _, count := range authorBookCount {
		if count < minBooks {
			minBooks = count
		}
		if count > maxBooks {
			maxBooks = count
		}
	}

	fmt.Println("\n📊 作者书籍分配统计:")
	fmt.Printf("  总书籍数: %d\n", totalBooks)
	fmt.Printf("  总作者数: %d\n", totalAuthors)
	fmt.Printf("  平均每作者: %.1f 本\n", avgBooksPerAuthor)
	fmt.Printf("  最少: %d 本, 最多: %d 本\n", minBooks, maxBooks)
}

// SeedBanners 填充 banner 数据
func (s *BookstoreSeeder) SeedBanners() error {
	ctx := context.Background()
	collection := s.db.Collection("banners")

	now := time.Now()

	// 定义 banners - 字段名与Banner模型匹配
	banners := []interface{}{
		map[string]interface{}{
			"_id":         primitive.NewObjectID(),
			"title":       "新书推荐",
			"description": "最新上架的精品好书",
			"image":       "/images/banners/new_books.jpg",
			"target":      "/books/new",
			"target_type": "url",
			"sort_order":  1,
			"is_active":   true,
			"start_time":  now,
			"end_time":    now.Add(30 * 24 * time.Hour),
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"_id":         primitive.NewObjectID(),
			"title":       "限时免费",
			"description": "限时免费阅读热门作品",
			"image":       "/images/banners/free_books.jpg",
			"target":      "/books/free",
			"target_type": "url",
			"sort_order":  2,
			"is_active":   true,
			"start_time":  now,
			"end_time":    now.Add(7 * 24 * time.Hour),
			"created_at":  now,
			"updated_at":  now,
		},
	}

	// 批量插入 banners
	_, err := collection.InsertMany(ctx, banners)
	if err != nil {
		return fmt.Errorf("插入 banners 失败: %w", err)
	}

	fmt.Printf("成功插入 %d 个 banner\n", len(banners))
	return nil
}

// Clean 清空书城数据（books 和 banners 集合）
func (s *BookstoreSeeder) Clean() error {
	ctx := context.Background()

	// 清空 books 集合
	_, err := s.db.Collection("books").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 books 集合失败: %w", err)
	}

	// 清空 banners 集合
	_, err = s.db.Collection("banners").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 banners 集合失败: %w", err)
	}

	fmt.Println("已清空 books 和 banners 集合")
	return nil
}

// Count 统计书籍数量
func (s *BookstoreSeeder) Count() (int64, error) {
	ctx := context.Background()
	count, err := s.db.Collection("books").CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("统计书籍数量失败: %w", err)
	}
	return count, nil
}
