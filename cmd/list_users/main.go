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
	coll := db.Collection("users")

	// 查找所有用户
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	fmt.Printf("数据库中的用户列表:\n")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for i, user := range results {
		username := user["username"]
		email := user["email"]
		userID := user["_id"]
		fmt.Printf("%d. 用户名: %v\n", i+1, username)
		fmt.Printf("   邮箱: %v\n", email)
		fmt.Printf("   ID: %v\n", userID)
		fmt.Println()
	}
}
