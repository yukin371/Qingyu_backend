//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// 查询testuser
	var user bson.M
	err = usersCollection.FindOne(ctx, bson.M{"username": "testuser"}).Decode(&user)
	if err != nil {
		fmt.Printf("查询用户失败: %v\n", err)
		return
	}

	fmt.Println("=== 测试用户信息 ===")
	fmt.Printf("用户名: %s\n", user["username"])
	fmt.Printf("邮箱: %s\n", user["email"])
	fmt.Printf("角色: %v\n", user["roles"])
	fmt.Printf("状态: %s\n", user["status"])
	fmt.Printf("昵称: %s\n", user["nickname"])

	password, ok := user["password"].(string)
	if ok {
		fmt.Printf("密码哈希: %s\n", password)
		fmt.Printf("哈希长度: %d\n", len(password))

		// 检查是否是bcrypt格式（以$2a$或$2b$开头）
		if len(password) == 60 && (password[:4] == "$2a$" || password[:4] == "$2b$") {
			fmt.Printf("✓ 密码格式正确（bcrypt哈希）\n")
		} else {
			fmt.Printf("❌ 密码格式不正确\n")
		}
	}
}
