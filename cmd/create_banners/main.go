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
	fmt.Println("=== 创建书城轮播图测试数据 ===\n")

	// 加载配置
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	db, err := connectDB()
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	ctx := context.Background()

	// 获取一些高评分书籍的ID
	bookCollection := db.Collection("books")
	cursor, err := bookCollection.Find(ctx, bson.M{},
		options.Find().SetLimit(10).SetSort(bson.D{{Key: "rating", Value: -1}}))
	if err != nil {
		log.Fatalf("查询书籍失败: %v", err)
	}
	defer cursor.Close(ctx)

	var books []struct {
		ID     primitive.ObjectID `bson:"_id"`
		Title  string             `bson:"title"`
		Cover  string             `bson:"cover"`
		Rating float64            `bson:"rating"`
	}
	if err = cursor.All(ctx, &books); err != nil {
		log.Fatalf("读取书籍失败: %v", err)
	}

	if len(books) == 0 {
		log.Fatal("没有找到推荐书籍，请先导入书籍数据")
	}

	fmt.Printf("找到 %d 本推荐书籍\n\n", len(books))

	// 创建轮播图数据
	bannerCollection := db.Collection("banners")

	// 清空现有轮播图
	_, err = bannerCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("清空现有轮播图失败: %v", err)
	}

	now := time.Now()
	startTime1 := now.Add(-7 * 24 * time.Hour)
	endTime1 := now.Add(30 * 24 * time.Hour)

	banners := []interface{}{
		&bookstore.Banner{
			ID:          primitive.NewObjectID(),
			Title:       "精品推荐：" + books[0].Title,
			Description: "高分力作，不容错过！评分：" + fmt.Sprintf("%.1f", books[0].Rating),
			Image:       books[0].Cover,
			Target:      books[0].ID.Hex(),
			TargetType:  "book",
			SortOrder:   1,
			IsActive:    true,
			StartTime:   &startTime1,
			EndTime:     &endTime1,
			ClickCount:  150,
			CreatedAt:   startTime1,
			UpdatedAt:   now,
		},
	}

	if len(books) > 1 {
		startTime2 := now.Add(-5 * 24 * time.Hour)
		endTime2 := now.Add(30 * 24 * time.Hour)
		banners = append(banners, &bookstore.Banner{
			ID:          primitive.NewObjectID(),
			Title:       "热门连载：" + books[1].Title,
			Description: "最热门的小说，千万读者的选择",
			Image:       books[1].Cover,
			Target:      books[1].ID.Hex(),
			TargetType:  "book",
			SortOrder:   2,
			IsActive:    true,
			StartTime:   &startTime2,
			EndTime:     &endTime2,
			ClickCount:  230,
			CreatedAt:   startTime2,
			UpdatedAt:   now,
		})
	}

	if len(books) > 2 {
		startTime3 := now.Add(-3 * 24 * time.Hour)
		endTime3 := now.Add(30 * 24 * time.Hour)
		banners = append(banners, &bookstore.Banner{
			ID:          primitive.NewObjectID(),
			Title:       "编辑精选：" + books[2].Title,
			Description: "编辑部精心挑选的优质作品",
			Image:       books[2].Cover,
			Target:      books[2].ID.Hex(),
			TargetType:  "book",
			SortOrder:   3,
			IsActive:    true,
			StartTime:   &startTime3,
			EndTime:     &endTime3,
			ClickCount:  180,
			CreatedAt:   startTime3,
			UpdatedAt:   now,
		})
	}

	if len(books) > 3 {
		startTime4 := now.Add(-1 * 24 * time.Hour)
		endTime4 := now.Add(30 * 24 * time.Hour)
		banners = append(banners, &bookstore.Banner{
			ID:          primitive.NewObjectID(),
			Title:       "新书上架：" + books[3].Title,
			Description: "新鲜出炉的精彩小说，抢先阅读",
			Image:       books[3].Cover,
			Target:      books[3].ID.Hex(),
			TargetType:  "book",
			SortOrder:   4,
			IsActive:    true,
			StartTime:   &startTime4,
			EndTime:     &endTime4,
			ClickCount:  95,
			CreatedAt:   startTime4,
			UpdatedAt:   now,
		})
	}

	if len(books) > 4 {
		startTime5 := now.Add(-2 * 24 * time.Hour)
		endTime5 := now.Add(30 * 24 * time.Hour)
		banners = append(banners, &bookstore.Banner{
			ID:          primitive.NewObjectID(),
			Title:       "人气爆款：" + books[4].Title,
			Description: "超人气作品，口碑爆棚",
			Image:       books[4].Cover,
			Target:      books[4].ID.Hex(),
			TargetType:  "book",
			SortOrder:   5,
			IsActive:    true,
			StartTime:   &startTime5,
			EndTime:     &endTime5,
			ClickCount:  310,
			CreatedAt:   startTime5,
			UpdatedAt:   now,
		})
	}

	// 插入轮播图
	result, err := bannerCollection.InsertMany(ctx, banners)
	if err != nil {
		log.Fatalf("插入轮播图失败: %v", err)
	}

	fmt.Printf("✓ 成功创建 %d 个轮播图\n", len(result.InsertedIDs))

	// 创建索引
	_, err = bannerCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "sort_order", Value: 1}}},
		{Keys: bson.D{{Key: "is_active", Value: 1}}},
		{Keys: bson.D{{Key: "start_time", Value: 1}, {Key: "end_time", Value: 1}}},
	})
	if err != nil {
		log.Printf("创建索引失败: %v", err)
	} else {
		fmt.Println("✓ 索引创建成功")
	}

	fmt.Println("\n=== 轮播图数据创建完成 ===")
}

func connectDB() (*mongo.Database, error) {
	dbConfig := config.GlobalConfig.Database
	if dbConfig == nil {
		return nil, fmt.Errorf("database configuration is missing")
	}

	if dbConfig.Primary.MongoDB == nil {
		return nil, fmt.Errorf("MongoDB configuration is missing")
	}

	mongoConfig := dbConfig.Primary.MongoDB
	clientOptions := options.Client().ApplyURI(mongoConfig.URI)
	clientOptions.SetMaxPoolSize(mongoConfig.MaxPoolSize)
	clientOptions.SetMinPoolSize(mongoConfig.MinPoolSize)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping database
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client.Database(mongoConfig.Database), nil
}
