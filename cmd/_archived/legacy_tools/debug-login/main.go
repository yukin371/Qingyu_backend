package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 连接到MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("qingyu")
	collection := db.Collection("users")

	// 查询reader1和testuser002的密码
	users := []string{"reader1", "testuser002"}

	for _, username := range users {
		var user bson.M
		err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
		if err != nil {
			fmt.Printf("用户 %s 不存在\n", username)
			continue
		}

		hash, _ := user["password"].(string)
		roles, _ := user["roles"].(primitive.A)
		status, _ := user["status"].(string)

		fmt.Printf("\n=== 用户: %s ===\n", username)
		fmt.Printf("密码哈希: %s\n", hash)
		fmt.Printf("角色: %v\n", roles)
		fmt.Printf("状态: %s\n", status)

		// 测试密码验证
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte("password123"))
		if err == nil {
			fmt.Println("✓ 密码 'password123' 验证成功")
		} else {
			fmt.Printf("✗ 密码 'password123' 验证失败: %v\n", err)
		}
	}
}
