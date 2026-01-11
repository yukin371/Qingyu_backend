//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/models/bookstore"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.Primary.MongoDB.URI))
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(cfg.Database.Primary.MongoDB.Database)
	booksCollection := db.Collection("books")
	bannersCollection := db.Collection("banners")

	fmt.Println("=== 创建测试Banner数据 ===")

	// 1. 查找现有书籍
	cursor, err := booksCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("查询书籍失败: %v", err)
	}
	defer cursor.Close(ctx)

	var books []bookstore.Book
	if err = cursor.All(ctx, &books); err != nil {
		log.Fatalf("解析书籍数据失败: %v", err)
	}

	if len(books) == 0 {
		fmt.Println("❌ 数据库中没有书籍，请先运行 create_test_books.go 创建测试书籍")
		return
	}

	fmt.Printf("✓ 找到 %d 本书籍\n", len(books))

	// 2. 清除旧的Banner
	result, err := bannersCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatalf("删除旧Banner失败: %v", err)
	}
	fmt.Printf("✓ 删除了 %d 个旧Banner\n", result.DeletedCount)

	// 3. 创建新的Banner（最多取前5本书）
	bannerCount := 5
	if len(books) < bannerCount {
		bannerCount = len(books)
	}

	banners := []interface{}{}
	now := time.Now()
	startTime := now.Add(-24 * time.Hour)
	endTime := now.Add(30 * 24 * time.Hour)

	for i := 0; i < bannerCount; i++ {
		book := books[i]
		banner := bookstore.Banner{
			ID:          primitive.NewObjectID(),
			Title:       fmt.Sprintf("推荐 - %s", book.Title),
			Description: book.Introduction,
			Image:       book.Cover,
			Target:      book.ID.Hex(),
			TargetType:  "book",
			SortOrder:   i + 1,
			IsActive:    true,
			StartTime:   &startTime,
			EndTime:     &endTime,
			ClickCount:  0,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		banners = append(banners, banner)
	}

	// 4. 插入Banner
	if len(banners) > 0 {
		insertResult, err := bannersCollection.InsertMany(ctx, banners)
		if err != nil {
			log.Fatalf("插入Banner失败: %v", err)
		}
		fmt.Printf("✓ 成功创建 %d 个Banner\n", len(insertResult.InsertedIDs))
	}

	// 5. 显示创建的Banner
	fmt.Println("\n=== 创建的Banner列表 ===")
	for i, b := range banners {
		banner := b.(bookstore.Banner)
		fmt.Printf("%d. %s\n", i+1, banner.Title)
		fmt.Printf("   书籍ID: %s\n", banner.Target)
		fmt.Printf("   目标类型: %s\n", banner.TargetType)
		fmt.Printf("   图片: %s\n", banner.Image)
		fmt.Println()
	}

	fmt.Println("✅ Banner数据创建完成！")
	fmt.Println("\n现在可以在前端首页看到这些Banner了")
	fmt.Println("访问: http://localhost:5173/")
}
