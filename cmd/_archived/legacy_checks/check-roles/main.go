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

	// 连接到MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	usersCol := db.Collection("users")

	// 查询所有测试用户及其角色
	users := []string{"author1", "reader1", "admin1"}

	for _, username := range users {
		filter := bson.M{"username": username}
		var user bson.M
		err := usersCol.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			fmt.Printf("❌ 用户 %s 查询失败: %v\n", username, err)
			continue
		}

		fmt.Printf("✅ 用户: %s\n", username)
		fmt.Printf("   Roles: %v\n", user["roles"])
		fmt.Printf("   Status: %v\n\n", user["status"])
	}
}
