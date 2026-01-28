package seeds

import (
	"context"
	"fmt"
	"log"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/models/shared"
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

	// VIP用户（5个）
	{
		Username:    "vip_user01",
		Email:       "vip01@qingyu.com",
		Password:    "Vip@123456",
		Role:        "reader",
		IsPremium:   true,
		Description: "VIP测试用户1",
	},
	{
		Username:    "vip_user02",
		Email:       "vip02@qingyu.com",
		Password:    "Vip@123456",
		Role:        "reader",
		IsPremium:   true,
		Description: "VIP测试用户2",
	},
	{
		Username:    "vip_user03",
		Email:       "vip03@qingyu.com",
		Password:    "Vip@123456",
		Role:        "reader",
		IsPremium:   true,
		Description: "VIP测试用户3",
	},
	{
		Username:    "vip_user04",
		Email:       "vip04@qingyu.com",
		Password:    "Vip@123456",
		Role:        "reader",
		IsPremium:   true,
		Description: "VIP测试用户4",
	},
	{
		Username:    "vip_user05",
		Email:       "vip05@qingyu.com",
		Password:    "Vip@123456",
		Role:        "reader",
		IsPremium:   true,
		Description: "VIP测试用户5",
	},

	// 普通用户（8个）
	{
		Username:    "user01",
		Email:       "user01@qingyu.com",
		Password:    "User@123456",
		Role:        "reader",
		IsPremium:   false,
		Description: "普通用户1",
	},
	{
		Username:    "user02",
		Email:       "user02@qingyu.com",
		Password:    "User@123456",
		Role:        "reader",
		IsPremium:   false,
		Description: "普通用户2",
	},
	{
		Username:    "user03",
		Email:       "user03@qingyu.com",
		Password:    "User@123456",
		Role:        "reader",
		IsPremium:   false,
		Description: "普通用户3",
	},
	{
		Username:    "user04",
		Email:       "user04@qingyu.com",
		Password:    "User@123456",
		Role:        "reader",
		IsPremium:   false,
		Description: "普通用户4",
	},
	{
		Username:    "user05",
		Email:       "user05@qingyu.com",
		Password:    "User@123456",
		Role:        "reader",
		IsPremium:   false,
		Description: "普通用户5",
	},
	{
		Username:    "user06",
		Email:       "user06@qingyu.com",
		Password:    "User@123456",
		Role:        "reader",
		IsPremium:   false,
		Description: "普通用户6",
	},
	{
		Username:    "user07",
		Email:       "user07@qingyu.com",
		Password:    "User@123456",
		Role:        "reader",
		IsPremium:   false,
		Description: "普通用户7",
	},
	{
		Username:    "user08",
		Email:       "user08@qingyu.com",
		Password:    "User@123456",
		Role:        "reader",
		IsPremium:   false,
		Description: "普通用户8",
	},
}

// CreateTestUsers 创建内测账号
func CreateTestUsers() error {
	ctx := context.Background()

	// 初始化MongoDB连接
	cfg, err := config.LoadConfig("")
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 获取MongoDB配置
	mongoConfig, err := cfg.Database.GetMongoConfig()
	if err != nil {
		return fmt.Errorf("获取MongoDB配置失败: %w", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConfig.URI))
	if err != nil {
		return fmt.Errorf("连接MongoDB失败: %w", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(mongoConfig.Database)
	collection := db.Collection("users")

	// 清空现有测试用户
	_, err = collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空用户失败: %w", err)
	}

	log.Println("开始创建内测账号...")

	// 创建测试用户
	for _, testUser := range testUsers {
		// 加密密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testUser.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("  ✗ 用户 %s 密码加密失败: %v", testUser.Username, err)
			continue
		}

		// 创建用户对象
		now := time.Now()
		user := users.User{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
			Username:         testUser.Username,
			Email:            testUser.Email,
			Phone:            "",
			Password:         string(hashedPassword),
			Roles:            []string{testUser.Role},
			Status:           users.UserStatusActive,
		}

		// 插入数据库
		_, err = collection.InsertOne(ctx, user)
		if err != nil {
			log.Printf("  ✗ 创建用户 %s 失败: %v", testUser.Username, err)
			continue
		}

		log.Printf("  ✓ 创建用户: %s (%s)", testUser.Username, testUser.Description)
	}

	log.Printf("成功创建 %d 个内测账号", len(testUsers))
	return nil
}
