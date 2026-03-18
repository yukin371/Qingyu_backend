package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	// 使用qingyu_test数据库（E2E测试使用的数据库）
	db := client.Database("qingyu_test")
	collection := db.Collection("users")

	// 生成 password123 的哈希
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("生成哈希失败:", err)
	}

	// 定义测试用户
	testUsers := []struct {
		username string
		email    string
		roles    []string
		nickname string
	}{
		{"reader1", "reader1@qingyu.com", []string{"reader"}, "书虫小王"},
		{"author1", "author1@qingyu.com", []string{"reader", "author"}, "墨小白"},
	}

	for _, user := range testUsers {
		// 检查用户是否存在
		var existing bson.M
		err := collection.FindOne(context.Background(), bson.M{"username": user.username}).Decode(&existing)

		if err == mongo.ErrNoDocuments {
			// 用户不存在，创建新用户
			newUser := bson.M{
				"username":  user.username,
				"email":     user.email,
				"password":  string(passwordHash),
				"roles":     user.roles,
				"nickname":  user.nickname,
				"avatar":    "/images/avatars/default.png",
				"status":    "active",
				"createdAt": primitive.NewDateTimeFromTime(time.Now()),
				"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
			}
			_, err := collection.InsertOne(context.Background(), newUser)
			if err != nil {
				log.Printf("创建用户 %s 失败: %v\n", user.username, err)
			} else {
				fmt.Printf("✓ 创建用户 %s 成功 (数据库: qingyu_test)\n", user.username)
			}
		} else if err == nil {
			// 用户存在，更新密码和角色
			result, err := collection.UpdateOne(
				context.Background(),
				bson.M{"username": user.username},
				bson.M{"$set": bson.M{"password": string(passwordHash), "roles": user.roles, "status": "active"}},
			)
			if err != nil {
				log.Printf("更新用户 %s 失败: %v\n", user.username, err)
			} else {
				fmt.Printf("✓ 更新用户 %s 密码和角色 (ModifiedCount: %d, 数据库: qingyu_test)\n", user.username, result.ModifiedCount)
			}
		} else {
			log.Printf("查询用户 %s 失败: %v\n", user.username, err)
		}
	}

	fmt.Println("\n测试用户设置完成!")
	fmt.Println("数据库: qingyu_test")
	fmt.Println("reader1 密码: password123")
	fmt.Println("author1 密码: password123")
	fmt.Println("testuser002 密码: password123 (已存在)")
}
