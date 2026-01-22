package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 连接MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	collection := db.Collection("books")

	// 查询所有使用placeholder.com的书籍
	filter := bson.M{"cover": bson.M{"$regex": "placeholder.com"}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	defer cursor.Close(ctx)

	var books []bson.M
	if err = cursor.All(ctx, &books); err != nil {
		log.Fatal("解析结果失败:", err)
	}

	fmt.Printf("找到 %d 本使用placeholder封面的书籍\n\n", len(books))

	// 更新每本书的封面路径
	updatedCount := 0
	for _, book := range books {
		id := book["_id"]
		title := book["title"]

		// 使用默认封面
		newCover := "/images/covers/default-book-cover.jpg"

		// 更新文档
		update := bson.M{"$set": bson.M{"cover": newCover}}
		_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
		if err != nil {
			log.Printf("更新书籍 %v 失败: %v\n", title, err)
			continue
		}

		updatedCount++
		fmt.Printf("✓ 已更新: %v\n", title)
	}

	fmt.Printf("\n成功更新 %d/%d 本书的封面路径\n", updatedCount, len(books))
	fmt.Println("所有书籍封面已更新为: /images/covers/default-book-cover.jpg")
}
