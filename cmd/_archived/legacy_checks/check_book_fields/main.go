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
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	coll := db.Collection("books")

	// 查看包含"爱"的书籍
	filter := bson.M{"title": bson.M{"$regex": "爱"}}
	cursor, err := coll.Find(ctx, filter, options.Find().SetLimit(3))
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	var results []bson.M
	cursor.All(ctx, &results)

	fmt.Printf("找到 %d 本包含'爱'的书:\n", len(results))
	for _, book := range results {
		fmt.Printf("  - _id: %v\n", book["_id"])
		fmt.Printf("    title: %v\n", book["title"])
		fmt.Printf("    title(类型): %T\n", book["title"])
		fmt.Printf("    author: %v\n", book["author"])
		fmt.Printf("    字段: %v\n", getKeys(book))
	}
}

func getKeys(m bson.M) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
