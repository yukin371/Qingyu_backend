package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Printf("连接MongoDB失败: %v\n", err)
		return
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	usersCollection := db.Collection("users")

	// 1. 查找testuser
	fmt.Println("=== 1. 查找testuser ===")
	var user bson.M
	err = usersCollection.FindOne(ctx, bson.M{"username": "testuser"}).Decode(&user)
	if err != nil {
		fmt.Printf("❌ 查找用户失败: %v\n", err)
		return
	}

	fmt.Printf("✓ 找到用户:\n")
	fmt.Printf("  ID: %v\n", user["_id"])
	fmt.Printf("  用户名: %s\n", user["username"])
	fmt.Printf("  密码哈希: %s\n", user["password"])
	fmt.Printf("  状态: %s\n", user["status"])
	fmt.Printf("  角色: %v\n", user["roles"])

	// 2. 验证密码
	fmt.Println("\n=== 2. 验证密码 ===")
	password := "123456"
	passwordHash, ok := user["password"].(string)
	if !ok {
		fmt.Printf("❌ 密码字段类型错误\n")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err == nil {
		fmt.Printf("✓ 密码验证成功\n")
	} else {
		fmt.Printf("❌ 密码验证失败: %v\n", err)
	}

	// 3. 测试错误的密码
	fmt.Println("\n=== 3. 测试错误的密码 ===")
	wrongPassword := "wrongpassword"
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(wrongPassword))
	if err == nil {
		fmt.Printf("❌ 错误密码竟然验证成功了（不应该）\n")
	} else {
		fmt.Printf("✓ 错误密码正确地被拒绝: %v\n", err)
	}
}
