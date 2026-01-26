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

type User struct {
	ID       string   `bson:"_id,omitempty"`
	Username string   `bson:"username"`
	Email    string   `bson:"email"`
	Roles    []string `bson:"roles"`
	Status   string   `bson:"status"`
}

func main() {
	// 连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	usersCol := db.Collection("users")

	// 查询 test_admin 用户
	var user User
	err = usersCol.FindOne(ctx, bson.M{"username": "test_admin"}).Decode(&user)
	if err != nil {
		log.Fatalf("查询用户失败: %v", err)
	}

	fmt.Println("========================================")
	fmt.Println("test_admin 用户信息")
	fmt.Println("========================================")
	fmt.Printf("用户ID: %s\n", user.ID)
	fmt.Printf("用户名: %s\n", user.Username)
	fmt.Printf("邮箱: %s\n", user.Email)
	fmt.Printf("状态: %s\n", user.Status)
	fmt.Printf("角色: %v\n", user.Roles)
	fmt.Println("========================================")

	// 验证管理员角色
	hasAdminRole := false
	for _, role := range user.Roles {
		if role == "admin" || role == "super_admin" {
			hasAdminRole = true
			break
		}
	}

	if hasAdminRole {
		fmt.Println("✓ test_admin 具有管理员角色")
	} else {
		fmt.Println("✗ test_admin 不具有管理员角色")
	}

	// 验证其他测试用户
	fmt.Println("\n========================================")
	fmt.Println("所有管理员用户列表")
	fmt.Println("========================================")

	cursor, err := usersCol.Find(ctx, bson.M{"roles": "admin"})
	if err != nil {
		log.Printf("查询管理员用户失败: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var adminUsers []User
	if err = cursor.All(ctx, &adminUsers); err != nil {
		log.Printf("解析管理员用户失败: %v", err)
		return
	}

	for i, admin := range adminUsers {
		fmt.Printf("%d. %s (%s) - 角色: %v\n", i+1, admin.Username, admin.Email, admin.Roles)
	}

	fmt.Printf("总计: %d 个管理员用户\n", len(adminUsers))
	fmt.Println("========================================")
}
