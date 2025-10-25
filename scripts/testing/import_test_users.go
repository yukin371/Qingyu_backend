package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/models/users"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// TestUser 测试账号结构
type TestUser struct {
	Username    string
	Email       string
	Password    string
	Role        string
	Description string
}

// 测试账号列表
var testUsers = []TestUser{
	// 管理员账号
	{
		Username:    "admin",
		Email:       "admin@qingyu.com",
		Password:    "Admin@123456",
		Role:        "admin",
		Description: "系统管理员",
	},

	// VIP用户
	{
		Username:    "vip_user01",
		Email:       "vip01@qingyu.com",
		Password:    "Vip@123456",
		Role:        "vip",
		Description: "VIP测试用户1",
	},
	{
		Username:    "vip_user02",
		Email:       "vip02@qingyu.com",
		Password:    "Vip@123456",
		Role:        "vip",
		Description: "VIP测试用户2",
	},

	// 普通用户
	{
		Username:    "test_user01",
		Email:       "test01@qingyu.com",
		Password:    "Test@123456",
		Role:        "user",
		Description: "普通测试用户1",
	},
	{
		Username:    "test_user02",
		Email:       "test02@qingyu.com",
		Password:    "Test@123456",
		Role:        "user",
		Description: "普通测试用户2",
	},
	{
		Username:    "test_user03",
		Email:       "test03@qingyu.com",
		Password:    "Test@123456",
		Role:        "user",
		Description: "普通测试用户3",
	},
	{
		Username:    "test_user04",
		Email:       "test04@qingyu.com",
		Password:    "Test@123456",
		Role:        "user",
		Description: "普通测试用户4",
	},
	{
		Username:    "test_user05",
		Email:       "test05@qingyu.com",
		Password:    "Test@123456",
		Role:        "user",
		Description: "普通测试用户5",
	},
}

func main() {
	fmt.Println("========================================")
	fmt.Println("青羽后端 - 测试用户导入工具")
	fmt.Println("========================================")
	fmt.Println()

	// 加载配置（使用测试配置）
	_, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 直接连接数据库（不依赖 ServiceContainer）
	db, err := connectDB()
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	ctx := context.Background()

	// 获取用户集合
	userCollection := db.Collection("users")

	fmt.Println("开始创建测试账号...")
	fmt.Println()

	successCount := 0
	skipCount := 0
	errorCount := 0

	for i, testUser := range testUsers {
		fmt.Printf("[%d/%d] 创建账号: %s (%s)...", i+1, len(testUsers), testUser.Username, testUser.Description)

		// 检查用户是否已存在
		var existingUser users.User
		err := userCollection.FindOne(ctx, bson.M{
			"$or": []bson.M{
				{"username": testUser.Username},
				{"email": testUser.Email},
			},
		}).Decode(&existingUser)

		if err == nil {
			fmt.Println(" [已存在，跳过]")
			skipCount++
			continue
		} else if err != mongo.ErrNoDocuments {
			fmt.Printf(" [查询错误: %v]\n", err)
			errorCount++
			continue
		}

		// 密码加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testUser.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf(" [密码加密失败: %v]\n", err)
			errorCount++
			continue
		}

		// 创建用户
		now := time.Now()
		user := users.User{
			ID:        primitive.NewObjectID().Hex(),
			Username:  testUser.Username,
			Email:     testUser.Email,
			Password:  string(hashedPassword),
			Role:      testUser.Role,
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err = userCollection.InsertOne(ctx, user)
		if err != nil {
			fmt.Printf(" [创建失败: %v]\n", err)
			errorCount++
			continue
		}

		fmt.Println(" [成功]")
		successCount++
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("账号创建完成")
	fmt.Println("========================================")
	fmt.Printf("✓ 成功创建: %d 个\n", successCount)
	fmt.Printf("○ 已存在跳过: %d 个\n", skipCount)
	fmt.Printf("✗ 创建失败: %d 个\n", errorCount)
	fmt.Printf("总计: %d 个\n", len(testUsers))
	fmt.Println()

	// 打印账号信息
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("测试账号清单")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	fmt.Println("【管理员账号】")
	for _, user := range testUsers {
		if user.Role == "admin" {
			fmt.Printf("  用户名: %s\n", user.Username)
			fmt.Printf("  邮箱: %s\n", user.Email)
			fmt.Printf("  密码: %s\n", user.Password)
			fmt.Printf("  说明: %s\n", user.Description)
			fmt.Println()
		}
	}

	fmt.Println("【VIP用户】")
	for _, user := range testUsers {
		if user.Role == "vip" {
			fmt.Printf("  用户名: %s | 邮箱: %s | 密码: %s\n", user.Username, user.Email, user.Password)
		}
	}
	fmt.Println()

	fmt.Println("【普通用户】")
	for _, user := range testUsers {
		if user.Role == "user" {
			fmt.Printf("  用户名: %s | 邮箱: %s | 密码: %s\n", user.Username, user.Email, user.Password)
		}
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("提示：")
	fmt.Println("1. 所有账号均已激活可直接登录")
	fmt.Println("2. 管理员和VIP用户享有高级权限")
	fmt.Println("3. 账号信息已保存在数据库中")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// connectDB 连接数据库
func connectDB() (*mongo.Database, error) {
	cfg := config.GlobalConfig.Database
	if cfg == nil || cfg.Primary.MongoDB == nil {
		return nil, fmt.Errorf("数据库配置缺失")
	}

	mongoCfg := cfg.Primary.MongoDB

	clientOptions := options.Client().
		ApplyURI(mongoCfg.URI).
		SetConnectTimeout(mongoCfg.ConnectTimeout).
		SetMaxPoolSize(mongoCfg.MaxPoolSize)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("连接 MongoDB 失败: %w", err)
	}

	// 验证连接
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("Ping MongoDB 失败: %w", err)
	}

	fmt.Printf("✓ 已连接到数据库: %s\n", mongoCfg.Database)
	fmt.Println()

	return client.Database(mongoCfg.Database), nil
}
