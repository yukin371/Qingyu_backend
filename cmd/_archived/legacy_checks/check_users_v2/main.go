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
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	usersCollection := db.Collection("users")

	// 查询 testuser001 和 testauthor001
	usernames := []string{"testuser001", "testauthor001", "testadmin001"}
	
	for _, username := range usernames {
		var user bson.M
		err := usersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil {
			fmt.Printf("❌ 用户 %s 不存在或查询失败: %v\n", username, err)
		} else {
			fmt.Printf("✓ 用户 %s 存在: id=%v, role=%v, nickname=%v\n", 
				username, user["_id"], user["role"], user["nickname"])
		}
	}

	// 统计用户总数
	total, _ := usersCollection.CountDocuments(ctx, bson.M{})
	fmt.Printf("\n数据库用户总数: %d\n", total)
}
