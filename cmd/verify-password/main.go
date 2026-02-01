package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
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

	// 查询reader1
	var user bson.M
	err = collection.FindOne(context.Background(), bson.M{"username": "reader1"}).Decode(&user)
	if err != nil {
		log.Fatal("查询失败:", err)
	}

	hash := user["password"].(string)
	fmt.Printf("数据库中的哈希: %s\n", hash)

	// 验证password123
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte("password123"))
	if err == nil {
		fmt.Println("✓ 密码 'password123' 验证成功")
	} else {
		fmt.Printf("✗ 密码 'password123' 验证失败: %v\n", err)
	}

	// 尝试生成新哈希并测试
	newHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("生成哈希失败:", err)
	}

	fmt.Printf("\n新生成的哈希: %s\n", string(newHash))
	err = bcrypt.CompareHashAndPassword(newHash, []byte("password123"))
	if err == nil {
		fmt.Println("✓ 新哈希验证成功")
	}
}
