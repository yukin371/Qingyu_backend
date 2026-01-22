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

	// 查询所有书籍的标题和封面
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"title": 1, "cover": 1}))
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal("解析结果失败:", err)
	}

	fmt.Printf("找到 %d 本书:\n\n", len(results))

	for _, book := range results {
		title := book["title"]
		cover := book["cover"]
		fmt.Printf("标题: %v\n", title)
		fmt.Printf("封面: %v\n\n", cover)
	}
}
