package seeds

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/users"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// EnhancedUser 增强用户结构
type EnhancedUser struct {
	Username    string
	Email       string
	Password    string
	Phone       string
	Role        string
	Status      string
	Avatar      string
	Bio         string
	Description string
}

// 增强版测试用户列表
var enhancedTestUsers = []EnhancedUser{
	// ========== 管理员（2个）==========
	{
		Username:    "admin",
		Email:       "admin@qingyu.com",
		Password:    "Admin@123456",
		Phone:       "13800000001",
		Role:        "admin",
		Status:      "active",
		Avatar:      "/avatars/admin.jpg",
		Bio:         "青羽写作平台系统管理员",
		Description: "超级管理员",
	},
	{
		Username:    "admin02",
		Email:       "admin02@qingyu.com",
		Password:    "Admin@123456",
		Phone:       "13800000002",
		Role:        "admin",
		Status:      "active",
		Avatar:      "/avatars/admin02.jpg",
		Bio:         "内容审核管理员",
		Description: "内容管理员",
	},

	// ========== 作者（5个）==========
	{
		Username:    "author_famous",
		Email:       "author_famous@qingyu.com",
		Password:    "Author@123456",
		Phone:       "13800001001",
		Role:        "author",
		Status:      "active",
		Avatar:      "/avatars/author_famous.jpg",
		Bio:         "知名网络小说作家，代表作《修真世界》",
		Description: "知名作家",
	},
	{
		Username:    "author_new",
		Email:       "author_new@qingyu.com",
		Password:    "Author@123456",
		Phone:       "13800001002",
		Role:        "author",
		Status:      "active",
		Avatar:      "/avatars/author_new.jpg",
		Bio:         "新人作家，正在努力创作第一部作品",
		Description: "新人作家",
	},
	{
		Username:    "author_veteran",
		Email:       "author_veteran@qingyu.com",
		Password:    "Author@123456",
		Phone:       "13800001003",
		Role:        "author",
		Status:      "active",
		Avatar:      "/avatars/author_veteran.jpg",
		Bio:         "资深作家，已创作10部作品",
		Description: "资深作家",
	},
	{
		Username:    "author_pro",
		Email:       "author_pro@qingyu.com",
		Password:    "Author@123456",
		Phone:       "13800001004",
		Role:        "author",
		Status:      "active",
		Avatar:      "/avatars/author_pro.jpg",
		Bio:         "专业全职作家，日更一万",
		Description: "全职作家",
	},
	{
		Username:    "author_parttime",
		Email:       "author_parttime@qingyu.com",
		Password:    "Author@123456",
		Phone:       "13800001005",
		Role:        "author",
		Status:      "active",
		Avatar:      "/avatars/author_parttime.jpg",
		Bio:         "业余爱好者，周末写作",
		Description: "业余作家",
	},

	// ========== VIP读者（3个）==========
	{
		Username:    "reader_vip01",
		Email:       "vip01@qingyu.com",
		Password:    "Vip@123456",
		Phone:       "13800002001",
		Role:        "reader",
		Status:      "active",
		Avatar:      "/avatars/vip01.jpg",
		Bio:         "资深书虫，年会员用户",
		Description: "VIP用户1",
	},
	{
		Username:    "reader_vip02",
		Email:       "vip02@qingyu.com",
		Password:    "Vip@123456",
		Phone:       "13800002002",
		Role:        "reader",
		Status:      "active",
		Avatar:      "/avatars/vip02.jpg",
		Bio:         "月会员用户，喜欢玄幻小说",
		Description: "VIP用户2",
	},
	{
		Username:    "reader_vip03",
		Email:       "vip03@qingyu.com",
		Password:    "Vip@123456",
		Phone:       "13800002003",
		Role:        "reader",
		Status:      "active",
		Avatar:      "/avatars/vip03.jpg",
		Bio:         "书城资深用户，收藏千本书",
		Description: "VIP用户3",
	},

	// ========== 普通读者（5个）==========
	{
		Username:    "reader_normal01",
		Email:       "reader01@qingyu.com",
		Password:    "Reader@123456",
		Phone:       "13800003001",
		Role:        "reader",
		Status:      "active",
		Avatar:      "/avatars/reader01.jpg",
		Bio:         "普通读者",
		Description: "普通用户1",
	},
	{
		Username:    "reader_normal02",
		Email:       "reader02@qingyu.com",
		Password:    "Reader@123456",
		Phone:       "13800003002",
		Role:        "reader",
		Status:      "active",
		Avatar:      "/avatars/reader02.jpg",
		Bio:         "喜欢看修真小说",
		Description: "普通用户2",
	},
	{
		Username:    "reader_normal03",
		Email:       "reader03@qingyu.com",
		Password:    "Reader@123456",
		Phone:       "13800003003",
		Role:        "reader",
		Status:      "active",
		Avatar:      "/avatars/reader03.jpg",
		Bio:         "都市小说爱好者",
		Description: "普通用户3",
	},
	{
		Username:    "reader_normal04",
		Email:       "reader04@qingyu.com",
		Password:    "Reader@123456",
		Phone:       "13800003004",
		Role:        "reader",
		Status:      "active",
		Avatar:      "/avatars/reader04.jpg",
		Bio:         "科幻迷",
		Description: "普通用户4",
	},
	{
		Username:    "reader_normal05",
		Email:       "reader05@qingyu.com",
		Password:    "Reader@123456",
		Phone:       "13800003005",
		Role:        "reader",
		Status:      "active",
		Avatar:      "/avatars/reader05.jpg",
		Bio:         "历史小说爱好者",
		Description: "普通用户5",
	},

	// ========== 特殊状态用户（2个）==========
	{
		Username:    "user_banned",
		Email:       "banned@qingyu.com",
		Password:    "Banned@123456",
		Phone:       "13800009001",
		Role:        "reader",
		Status:      "banned",
		Avatar:      "/avatars/banned.jpg",
		Bio:         "违规用户",
		Description: "被封禁用户",
	},
	{
		Username:    "user_inactive",
		Email:       "inactive@qingyu.com",
		Password:    "Inactive@123456",
		Phone:       "13800009002",
		Role:        "reader",
		Status:      "inactive",
		Avatar:      "/avatars/inactive.jpg",
		Bio:         "长期未登录",
		Description: "未激活用户",
	},
}

// SeedEnhancedUsers 增强版用户种子数据
func SeedEnhancedUsers(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	fmt.Println("========================================")
	fmt.Println("开始创建增强版测试用户数据...")
	fmt.Println("========================================")

	successCount := 0
	skipCount := 0
	errorCount := 0

	for i, testUser := range enhancedTestUsers {
		fmt.Printf("[%d/%d] 创建用户: %s (%s)...", i+1, len(enhancedTestUsers), testUser.Username, testUser.Description)

		// 检查用户是否已存在
		var existingUser users.User
		err := collection.FindOne(ctx, bson.M{
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
			Phone:     testUser.Phone,
			Password:  string(hashedPassword),
			Role:      testUser.Role,
			Status:    users.UserStatus(testUser.Status),
			Avatar:    testUser.Avatar,
			Bio:       testUser.Bio,
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err = collection.InsertOne(ctx, user)
		if err != nil {
			fmt.Printf(" [创建失败: %v]\n", err)
			errorCount++
			continue
		}

		fmt.Println(" [成功]")
		successCount++
	}

	fmt.Println("========================================")
	fmt.Println("用户创建完成")
	fmt.Println("========================================")
	fmt.Printf("成功创建: %d 个\n", successCount)
	fmt.Printf("已存在跳过: %d 个\n", skipCount)
	fmt.Printf("创建失败: %d 个\n", errorCount)
	fmt.Printf("总计: %d 个\n", len(enhancedTestUsers))
	fmt.Println()

	return nil
}

// GetTestUserCredentials 获取测试用户凭证列表
func GetTestUserCredentials() []EnhancedUser {
	return enhancedTestUsers
}
