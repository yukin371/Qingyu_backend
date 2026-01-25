package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID     string   `bson:"_id,omitempty"`
	Username string  `bson:"username"`
	Email    string  `bson:"email"`
	Password string  `bson:"password"`
	Role     string  `bson:"role"`
	Roles    []string `bson:"roles"`
	Nickname string  `bson:"nickname"`
	Status   string  `bson:"status"`
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

	// 检查测试管理员是否存在
	var existingUser User
	err = usersCol.FindOne(ctx, bson.M{"username": "test_admin"}).Decode(&existingUser)
	if err == nil {
		// 删除旧的测试管理员
		usersCol.DeleteOne(ctx, bson.M{"username": "test_admin"})
		fmt.Println("删除旧的测试管理员")
	}

	// 生成密码哈希
	password := "TestAdmin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	// 创建测试管理员
	objectID := primitive.NewObjectID()
	admin := User{
		ID:       objectID.Hex(),
		Username: "test_admin",
		Email:    "test_admin@qingyu.com",
		Password: string(hashedPassword),
		Role:     "admin",
		Roles:    []string{"admin"},
		Nickname: "测试管理员",
		Status:   "active",
	}

	_, err = usersCol.InsertOne(ctx, admin)
	if err != nil {
		log.Fatalf("插入用户失败: %v", err)
	}

	// 输出到文件，方便后续使用
	output := fmt.Sprintf("TEST_ADMIN_USERNAME=test_admin\nTEST_ADMIN_PASSWORD=%s", password)
	os.WriteFile("test_admin_credentials.txt", []byte(output), 0644)

	fmt.Println("✓ 测试管理员创建成功")
	fmt.Printf("  用户名: %s\n", admin.Username)
	fmt.Printf("  密码: %s\n", password)
	fmt.Printf("  角色: %s\n", admin.Role)
	fmt.Printf("  ID: %s\n", admin.ID)
}
