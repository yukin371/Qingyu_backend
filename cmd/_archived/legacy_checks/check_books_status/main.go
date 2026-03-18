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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	coll := db.Collection("books")

	// 检查前10本书的状态
	cursor, err := coll.Find(ctx, bson.M{}, options.Find().SetLimit(10))
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	fmt.Printf("数据库中的书籍状态（前10本）:\n")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for i, book := range results {
		title := book["title"]
		status := book["status"]
		fmt.Printf("%d. 标题: %v\n", i+1, title)
		fmt.Printf("   状态: %v\n\n", status)
	}
}
