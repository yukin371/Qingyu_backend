package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Printf("连接MongoDB失败: %v\n", err)
		return
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	fmt.Println("开始填充E2E测试数据...")

	// 1. 创建测试用户
	usersCollection := db.Collection("users")
	userCount, _ := usersCollection.CountDocuments(ctx, bson.M{"username": "testuser"})
	if userCount == 0 {
		now := time.Now()
		// 使用bcrypt哈希的密码（密码：123456）
		bcryptHash := "$2a$10$ilq82USQvkysphMM3WZ0UeA5YBTAXUu/SZPBczEdFOlRLzAVs6hw."
		user := bson.M{
			"_id":             "test-user-001", // 使用string类型的_id
			"username":        "testuser",
			"email":           "testuser@qingyu.com",
			"password":        bcryptHash,
			"roles":           bson.A{"reader"}, // 多角色支持
			"vip_level":       0,
			"status":          "active",
			"nickname":        "测试用户",
			"avatar":          "/images/avatars/default.png",
			"bio":             "这是一个测试账号",
			"email_verified":  false,
			"phone_verified":  false,
			"created_at":      now,
			"updated_at":      now,
		}
		_, err := usersCollection.InsertOne(ctx, user)
		if err != nil {
			fmt.Printf("创建测试用户失败: %v\n", err)
		} else {
			fmt.Println("✓ 已创建测试用户: testuser/123456")
		}
	} else {
		fmt.Println("✓ 测试用户已存在")
	}

	// 2. 创建测试分类
	categoriesCollection := db.Collection("categories")
	categoryCount, _ := categoriesCollection.CountDocuments(ctx, bson.M{"name": bson.M{"$in": bson.A{"玄幻", "修仙"}}})
	if categoryCount == 0 {
		now := time.Now()
		categories := []interface{}{
			bson.M{
				"name":        "玄幻",
				"slug":        "xuanhuan",
				"description": "奇幻玄幻，想象力无限",
				"icon":        "/images/icons/xuanhuan.png",
				"sort_order":  1,
				"is_active":   true,
				"created_at":  now,
				"updated_at":  now,
			},
			bson.M{
				"name":        "修仙",
				"slug":        "xiuxian",
				"description": "修仙问道，长生不老",
				"icon":        "/images/icons/xiuxian.png",
				"sort_order":  2,
				"is_active":   true,
				"created_at":  now,
				"updated_at":  now,
			},
		}
		_, err := categoriesCollection.InsertMany(ctx, categories)
		if err != nil {
			fmt.Printf("创建测试分类失败: %v\n", err)
		} else {
			fmt.Println("✓ 已创建测试分类: 玄幻、修仙")
		}
	} else {
		fmt.Println("✓ 测试分类已存在")
	}

	// 3. 创建测试书籍
	booksCollection := db.Collection("books")
	bookCount, _ := booksCollection.CountDocuments(ctx, bson.M{"title": bson.M{"$regex": "修仙"}})
	if bookCount == 0 {
		now := time.Now()
		publishedAt := now.Add(-180 * 24 * time.Hour)

		books := []interface{}{
			bson.M{
				"title":          "修仙世界",
				"author":         "飞升作者",
				"introduction":   "一个普通少年，意外获得神秘传承，踏上修仙之路。",
				"cover":          "/images/covers/xiuxian_shijie.jpg",
				"categories":     bson.A{"玄幻", "修仙"},
				"tags":           bson.A{"修仙", "玄幻", "升级", "热血"},
				"status":         "ongoing",
				"rating":         8.5,
				"rating_count":   1250,
				"view_count":     45000,
				"word_count":     1500000,
				"chapter_count":  450,
				"price":          0.0,
				"is_free":        true,
				"is_recommended": true,
				"is_featured":    true,
				"is_hot":         true,
				"published_at":   publishedAt,
				"last_update_at": now.Add(-24 * time.Hour),
				"created_at":     now,
				"updated_at":     now,
			},
			bson.M{
				"title":          "修仙归来",
				"author":         "逍遥子",
				"introduction":   "一代仙尊渡劫失败，重生回到地球，再登巅峰！",
				"cover":          "/images/covers/xiuxian_guilai.jpg",
				"categories":     bson.A{"玄幻", "修仙"},
				"tags":           bson.A{"修仙", "玄幻", "重生", "爽文"},
				"status":         "ongoing",
				"rating":         9.2,
				"rating_count":   8900,
				"view_count":     120000,
				"word_count":     2800000,
				"chapter_count":  820,
				"price":          9.9,
				"is_free":        false,
				"is_recommended": true,
				"is_featured":    true,
				"is_hot":         true,
				"published_at":   now.Add(-30 * 24 * time.Hour),
				"last_update_at": now.Add(-12 * time.Hour),
				"created_at":     now,
				"updated_at":     now,
			},
		}

		_, err := booksCollection.InsertMany(ctx, books)
		if err != nil {
			fmt.Printf("创建测试书籍失败: %v\n", err)
		} else {
			fmt.Printf("✓ 已创建 %d 本修仙小说测试书籍\n", len(books))
		}
	} else {
		fmt.Println("✓ 测试书籍已存在")
	}

	// 4. 验证数据
	userCount, _ = usersCollection.CountDocuments(ctx, bson.M{"username": "testuser"})
	bookCount, _ = booksCollection.CountDocuments(ctx, bson.M{"title": bson.M{"$regex": "修仙"}})
	categoryCount, _ = categoriesCollection.CountDocuments(ctx, bson.M{"name": bson.M{"$in": bson.A{"玄幻", "修仙"}}})

	fmt.Println("\n数据验证:")
	fmt.Printf("- 测试用户: %d\n", userCount)
	fmt.Printf("- 修仙书籍: %d\n", bookCount)
	fmt.Printf("- 测试分类: %d\n", categoryCount)

	if userCount > 0 && bookCount > 0 && categoryCount > 0 {
		fmt.Println("\n✅ 所有E2E测试数据添加成功!")
		fmt.Println("\n测试账号: testuser / 123456")
		fmt.Println("测试书籍: 修仙世界、修仙归来")
		fmt.Println("测试分类: 玄幻、修仙")
	} else {
		fmt.Println("\n❌ 数据添加不完整，请检查")
	}
}
