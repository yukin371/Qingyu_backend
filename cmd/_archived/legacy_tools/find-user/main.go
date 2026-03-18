package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 连接到MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer client.Disconnect(context.Background())

	// 获取所有数据库
	databases, err := client.ListDatabaseNames(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("获取数据库列表失败:", err)
	}

	username := "testuser002"
	fmt.Printf("搜索用户: %s\n", username)

	for _, dbName := range databases {
		db := client.Database(dbName)
		collection := db.Collection("users")

		var user bson.M
		err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
		if err == nil {
			fmt.Printf("\n✓ 在数据库 '%s' 中找到用户!\n", dbName)
			fmt.Printf("  用户名: %v\n", user["username"])
			fmt.Printf("  邮箱: %v\n", user["email"])
			fmt.Printf("  角色: %v\n", user["roles"])
			fmt.Printf("  状态: %v\n", user["status"])
		}
	}
}
