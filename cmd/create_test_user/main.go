package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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
	usersColl := db.Collection("users")

	// 密码
	password := "Reader123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	// 创建测试用户
	testUser := bson.M{
		"username": "reader01",
		"email":    "reader01@qingyu.com",
		"password": string(hashedPassword),
		"role":     "reader",
		"status":   "active",
		"nickname": "测试读者01",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	// 检查用户是否已存在
	var existing bson.M
	err = usersColl.FindOne(ctx, bson.M{"username": "reader01"}).Decode(&existing)
	if err == nil {
		fmt.Println("用户 reader01 已存在，更新密码...")
		_, err = usersColl.UpdateOne(
			ctx,
			bson.M{"username": "reader01"},
			bson.M{"$set": bson.M{"password": string(hashedPassword)}},
		)
		if err != nil {
			log.Fatalf("更新用户失败: %v", err)
		}
		fmt.Println("✓ 密码已更新")
	} else {
		// 创建新用户
		_, err = usersColl.InsertOne(ctx, testUser)
		if err != nil {
			log.Fatalf("创建用户失败: %v", err)
		}
		fmt.Println("✓ 测试用户已创建")
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("测试用户信息:")
	fmt.Println("用户名: reader01")
	fmt.Println("邮箱: reader01@qingyu.com")
	fmt.Println("密码: Reader123456")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
