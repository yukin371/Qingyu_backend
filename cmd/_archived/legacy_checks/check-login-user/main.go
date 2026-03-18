package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 连接MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	// 检查qingyu数据库
	db := client.Database("qingyu")
	usersColl := db.Collection("users")

	// 查找testuser002
	fmt.Println("=== 检查 testuser002 用户 ===")
	var user bson.M
	err = usersColl.FindOne(ctx, bson.M{"username": "testuser002"}).Decode(&user)
	if err != nil {
		fmt.Printf("用户 testuser002 未找到: %v\n", err)
	} else {
		fmt.Printf("找到用户: %+v\n", user)
	}

	// 查找author1
	fmt.Println("\n=== 检查 author1 用户 ===")
	err = usersColl.FindOne(ctx, bson.M{"username": "author1"}).Decode(&user)
	if err != nil {
		fmt.Printf("用户 author1 未找到: %v\n", err)
	} else {
		fmt.Printf("找到用户: %+v\n", user)
	}

	// 列出所有用户
	fmt.Println("\n=== 所有用户列表 ===")
	cursor, err := usersColl.Find(ctx, bson.M{})
	if err != nil {
		fmt.Printf("查询用户列表失败: %v\n", err)
		os.Exit(1)
	}
	defer cursor.Close(ctx)

	var users []bson.M
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatal("获取用户列表失败:", err)
	}

	for _, u := range users {
		username := u["username"]
		role := u["role"]
		fmt.Printf("- %s (role: %v)\n", username, role)
	}
}
