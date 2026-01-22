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
	booksCollection := db.Collection("books")

	// 查询修仙书籍的详细信息
	fmt.Println("=== 查询修仙书籍的详细信息 ===")
	var books []bson.M
	cursor, err := booksCollection.Find(ctx, bson.M{"title": bson.M{"$regex": "修仙"}})
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &books); err != nil {
		fmt.Printf("解析失败: %v\n", err)
		return
	}

	for i, book := range books {
		fmt.Printf("\n书籍 %d:\n", i+1)
		fmt.Printf("  _id: %v\n", book["_id"])
		fmt.Printf("  title: %s\n", book["title"])
		fmt.Printf("  author: %s\n", book["author"])
		fmt.Printf("  status: %s\n", book["status"])
		fmt.Printf("  categories: %v\n", book["categories"])
		fmt.Printf("  tags: %v\n", book["tags"])
		fmt.Printf("  is_free: %v\n", book["is_free"])
		fmt.Printf("  is_recommended: %v\n", book["is_recommended"])
		fmt.Printf("  is_featured: %v\n", book["is_featured"])
		fmt.Printf("  is_hot: %v\n", book["is_hot"])
	}

	// 测试不同的status查询
	fmt.Println("\n\n=== 测试不同的status查询 ===")

	// 1. 查询status=ongoing的书籍
	count1, _ := booksCollection.CountDocuments(ctx, bson.M{"status": "ongoing"})
	fmt.Printf("status=ongoing的书籍数量: %d\n", count1)

	// 2. 查询status=published的书籍
	count2, _ := booksCollection.CountDocuments(ctx, bson.M{"status": "published"})
	fmt.Printf("status=published的书籍数量: %d\n", count2)

	// 3. 查询status in [published, ongoing, completed]的书籍
	count3, _ := booksCollection.CountDocuments(ctx, bson.M{
		"status": bson.M{"$in": []string{"published", "ongoing", "completed"}},
	})
	fmt.Printf("status in [published, ongoing, completed]的书籍数量: %d\n", count3)

	// 4. 无status限制的书籍数量
	count4, _ := booksCollection.CountDocuments(ctx, bson.M{})
	fmt.Printf("所有书籍数量（无status限制）: %d\n", count4)
}
