package seeds

import (
	"context"
	"fmt"
	"log"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"
	"Qingyu_backend/models/users"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// TestUser 测试账号结构
type TestUser struct {
	Username    string
	Email       string
	Password    string
	Role        string
	IsPremium   bool
	Description string
}

// 内测账号列表
var testUsers = []TestUser{
	// 管理员账号（2个）
	{
		Username:    "admin",
		Email:       "admin@qingyu.com",
		Password:    "Admin@123456",
		Role:        "admin",
		IsPremium:   true,
		Description: "系统管理员",
	},
	{
		Username:    "admin_test",
		Email:       "admin_test@qingyu.com",
		Password:    "Admin@123456",
		Role:        "admin",
		IsPremium:   true,
		Description: "测试管理员",
	},

	// VIP用户（3个）
	{
		Username:    "vip_user01",
		Email:       "vip01@qingyu.com",
		Password:    "Vip@123456",
		Role:        "vip",
		IsPremium:   true,
		Description: "VIP测试用户1",
	},
	{
		Username:    "vip_user02",
		Email:       "vip02@qingyu.com",
		Password:    "Vip@123456",
		Role:        "vip",
		IsPremium:   true,
		Description: "VIP测试用户2",
	},
	{
		Username:    "vip_user03",
		Email:       "vip03@qingyu.com",
		Password:    "Vip@123456",
		Role:        "vip",
		IsPremium:   true,
		Description: "VIP测试用户3",
	},

	// 普通用户（5个）
	{
		Username:    "test_user01",
		Email:       "test01@qingyu.com",
		Password:    "Test@123456",
		Role:        "user",
		IsPremium:   false,
		Description: "普通测试用户1",
	},
	{
		Username:    "test_user02",
		Email:       "test02@qingyu.com",
		Password:    "Test@123456",
		Role:        "user",
		IsPremium:   false,
		Description: "普通测试用户2",
	},
	{
		Username:    "test_user03",
		Email:       "test03@qingyu.com",
		Password:    "Test@123456",
		Role:        "user",
		IsPremium:   false,
		Description: "普通测试用户3",
	},
	{
		Username:    "test_user04",
		Email:       "test04@qingyu.com",
		Password:    "Test@123456",
		Role:        "user",
		IsPremium:   false,
		Description: "普通测试用户4",
	},
	{
		Username:    "test_user05",
		Email:       "test05@qingyu.com",
		Password:    "Test@123456",
		Role:        "user",
		IsPremium:   false,
		Description: "普通测试用户5",
	},
}

func main() {
	fmt.Println("====================================")
	fmt.Println("创建MVP内测账号")
	fmt.Println("====================================")
	fmt.Println()

	// 加载配置
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	if err := core.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	ctx := context.Background()

	// 检查数据库连接
	if global.DB == nil {
		log.Fatal("数据库未初始化")
	}

	// 获取用户集合
	userCollection := global.DB.Collection("users")

	fmt.Println("开始创建测试账号...")
	fmt.Println()

	successCount := 0
	skipCount := 0

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
			fmt.Printf(" [错误: %v]\n", err)
			continue
		}

		// 密码加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testUser.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf(" [密码加密失败: %v]\n", err)
			continue
		}

		// 创建用户（User模型没有IsPremium字段，使用Role区分）
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
			continue
		}

		fmt.Println(" [成功]")
		successCount++
	}

	fmt.Println()
	fmt.Println("====================================")
	fmt.Println("账号创建完成")
	fmt.Println("====================================")
	fmt.Printf("成功创建: %d 个\n", successCount)
	fmt.Printf("已存在跳过: %d 个\n", skipCount)
	fmt.Printf("总计: %d 个\n", len(testUsers))
	fmt.Println()

	// 打印账号信息
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("内测账号列表")
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
			fmt.Printf("  用户名: %s\n", user.Username)
			fmt.Printf("  邮箱: %s\n", user.Email)
			fmt.Printf("  密码: %s\n", user.Password)
			fmt.Printf("  说明: %s\n", user.Description)
			fmt.Println()
		}
	}

	fmt.Println("【普通用户】")
	for _, user := range testUsers {
		if user.Role == "user" {
			fmt.Printf("  用户名: %s | 密码: %s\n", user.Username, user.Password)
		}
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("提示：")
	fmt.Println("1. 所有账号均已激活可直接登录")
	fmt.Println("2. 管理员和VIP用户享有premium权限")
	fmt.Println("3. 建议定期更换测试密码")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
