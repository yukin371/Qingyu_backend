package main

import (
	"context"
	"fmt"
	"log"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/service"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 加载配置
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v\n", err)
	}

	// 初始化服务
	if err := core.InitServices(); err != nil {
		log.Fatalf("❌ 初始化服务失败: %v\n", err)
	}

	ctx := context.Background()
	db := service.GetServiceContainer().GetMongoDB()
	collection := db.Collection("books")

	// 统计总书籍数
	total, _ := collection.CountDocuments(ctx, bson.M{})
	fmt.Printf("书籍总数: %d\n\n", total)

	// 统计推荐书籍
	recommended, _ := collection.CountDocuments(ctx, bson.M{"is_recommended": true})
	fmt.Printf("推荐书籍数 (is_recommended=true): %d\n", recommended)

	// 统计精选书籍
	featured, _ := collection.CountDocuments(ctx, bson.M{"is_featured": true})
	fmt.Printf("精选书籍数 (is_featured=true): %d\n", featured)

	// 统计热门书籍
	hot, _ := collection.CountDocuments(ctx, bson.M{"is_hot": true})
	fmt.Printf("热门书籍数 (is_hot=true): %d\n\n", hot)

	// 查看前5本书的数据
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetLimit(5))
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	fmt.Println("前5本书的数据:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for cursor.Next(ctx) {
		var book bson.M
		cursor.Decode(&book)
		fmt.Printf("标题: %v\n", book["title"])
		fmt.Printf("作者: %v\n", book["author"])
		fmt.Printf("推荐: %v, 精选: %v, 热门: %v\n\n",
			book["is_recommended"], book["is_featured"], book["is_hot"])
	}
}
