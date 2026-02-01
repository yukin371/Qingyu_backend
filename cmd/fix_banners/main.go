package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Banner struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Image       string             `bson:"image"`
	Target      string             `bson:"target"`
	TargetType  string             `bson:"target_type"`
	SortOrder   int                `bson:"sort_order"`
	IsActive    bool               `bson:"is_active"`
	StartTime   *time.Time         `bson:"start_time,omitempty"`
	EndTime     *time.Time         `bson:"end_time,omitempty"`
	ClickCount  int64              `bson:"click_count"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║      修复Banner数据 - ObjectId修复       ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// 连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	collection := db.Collection("banners")

	fmt.Println("🔗 数据库连接成功")
	fmt.Println()

	// 1. 清空现有banners
	fmt.Println("【1/2】清空现有Banner数据")
	deleteResult, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatalf("清空banners失败: %v", err)
	}
	fmt.Printf("  已删除 %d 条旧数据\n", deleteResult.DeletedCount)
	fmt.Println()

	// 2. 创建新的banners
	fmt.Println("【2/2】创建新的Banner数据")
	now := time.Now()

	banners := []interface{}{
		Banner{
			Title:       "修真世界 - 热门推荐",
			Description: "凡人流经典之作，不可错过",
			Image:       "https://images.unsplash.com/photo-1518709268805-4e9042af9f23?w=1200&h=400&fit=crop",
			Target:      "6956392cfe350a59abae6607",
			TargetType:  "book",
			SortOrder:   1,
			IsActive:    true,
			ClickCount:  0,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Banner{
			Title:       "新人福利大礼包",
			Description: "注册即送 100 青羽币",
			Image:       "https://images.unsplash.com/photo-1557683316-973673baf926?w=1200&h=400&fit=crop",
			Target:      "/promo/newbie",
			TargetType:  "url",
			SortOrder:   2,
			IsActive:    true,
			ClickCount:  0,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Banner{
			Title:       "限时免费阅读",
			Description: "精选作品限时免费，不容错过",
			Image:       "https://images.unsplash.com/photo-1512820790803-83ca734da794?w=1200&h=400&fit=crop",
			Target:      "/books/free",
			TargetType:  "url",
			SortOrder:   3,
			IsActive:    true,
			ClickCount:  0,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	// 插入数据
	insertResult, err := collection.InsertMany(ctx, banners)
	if err != nil {
		log.Fatalf("插入banners失败: %v", err)
	}

	fmt.Printf("  成功插入 %d 条新数据\n", len(insertResult.InsertedIDs))
	fmt.Println()

	// 3. 验证数据
	fmt.Println("【验证】检查Banner数据")
	count, err := collection.CountDocuments(ctx, bson.M{"is_active": true})
	if err != nil {
		log.Printf("统计数据失败: %v", err)
	} else {
		fmt.Printf("  当前有 %d 条活跃的Banner\n", count)
	}
	fmt.Println()

	// 4. 查看一条数据样本
	fmt.Println("【样本】查看第一条Banner数据")
	var sample Banner
	err = collection.FindOne(ctx, bson.M{}).Decode(&sample)
	if err != nil {
		log.Printf("查询样本数据失败: %v", err)
	} else {
		fmt.Printf("  ID: %s (类型: %T)\n", sample.ID.Hex(), sample.ID)
		fmt.Printf("  标题: %s\n", sample.Title)
		fmt.Printf("  目标: %s (%s)\n", sample.Target, sample.TargetType)
	}
	fmt.Println()

	fmt.Println("✨ Banner数据修复完成！")
	fmt.Println()
	fmt.Println("现在可以测试首页API: GET /api/v1/bookstore/homepage")
}
