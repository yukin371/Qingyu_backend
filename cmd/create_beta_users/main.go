package main

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

// BetaUser 内测用户结构
type BetaUser struct {
	Username    string
	Email       string
	Password    string
	Role        string
	Nickname    string
	Avatar      string
	Description string
}

// 内测用户列表
var betaUsers = []BetaUser{
	// ============ 管理员账号（3个）============
	{
		Username:    "admin",
		Email:       "admin@qingyu.com",
		Password:    "Admin@123456",
		Role:        "admin",
		Nickname:    "系统管理员",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=admin",
		Description: "超级管理员，拥有所有权限",
	},
	{
		Username:    "admin_test",
		Email:       "admin_test@qingyu.com",
		Password:    "Admin@123456",
		Role:        "admin",
		Nickname:    "测试管理员",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=admin_test",
		Description: "测试用管理员账号",
	},
	{
		Username:    "admin_audit",
		Email:       "admin_audit@qingyu.com",
		Password:    "Admin@123456",
		Role:        "admin",
		Nickname:    "审核管理员",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=admin_audit",
		Description: "负责内容审核的管理员",
	},

	// ============ VIP用户（5个）============
	{
		Username:    "vip_writer01",
		Email:       "vip_writer01@qingyu.com",
		Password:    "Vip@123456",
		Role:        "vip",
		Nickname:    "VIP作家1号",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=vip_writer01",
		Description: "VIP写作用户，测试高级写作功能",
	},
	{
		Username:    "vip_writer02",
		Email:       "vip_writer02@qingyu.com",
		Password:    "Vip@123456",
		Role:        "vip",
		Nickname:    "VIP作家2号",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=vip_writer02",
		Description: "VIP写作用户，测试AI辅助功能",
	},
	{
		Username:    "vip_reader01",
		Email:       "vip_reader01@qingyu.com",
		Password:    "Vip@123456",
		Role:        "vip",
		Nickname:    "VIP读者1号",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=vip_reader01",
		Description: "VIP阅读用户，测试付费阅读功能",
	},
	{
		Username:    "vip_reader02",
		Email:       "vip_reader02@qingyu.com",
		Password:    "Vip@123456",
		Role:        "vip",
		Nickname:    "VIP读者2号",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=vip_reader02",
		Description: "VIP阅读用户，测试书架和推荐功能",
	},
	{
		Username:    "vip_tester",
		Email:       "vip_tester@qingyu.com",
		Password:    "Vip@123456",
		Role:        "vip",
		Nickname:    "VIP全能测试员",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=vip_tester",
		Description: "VIP综合测试账号",
	},

	// ============ 普通写作用户（5个）============
	{
		Username:    "writer_xuanhuan",
		Email:       "writer_xuanhuan@qingyu.com",
		Password:    "Writer@123456",
		Role:        "user",
		Nickname:    "玄幻小说家",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=writer_xuanhuan",
		Description: "测试玄幻小说创作功能",
	},
	{
		Username:    "writer_yanqing",
		Email:       "writer_yanqing@qingyu.com",
		Password:    "Writer@123456",
		Role:        "user",
		Nickname:    "言情作家",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=writer_yanqing",
		Description: "测试言情小说创作功能",
	},
	{
		Username:    "writer_dushi",
		Email:       "writer_dushi@qingyu.com",
		Password:    "Writer@123456",
		Role:        "user",
		Nickname:    "都市创作者",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=writer_dushi",
		Description: "测试都市小说创作功能",
	},
	{
		Username:    "writer_newbie",
		Email:       "writer_newbie@qingyu.com",
		Password:    "Writer@123456",
		Role:        "user",
		Nickname:    "新手作者",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=writer_newbie",
		Description: "测试新手写作流程",
	},
	{
		Username:    "writer_pro",
		Email:       "writer_pro@qingyu.com",
		Password:    "Writer@123456",
		Role:        "user",
		Nickname:    "专业作家",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=writer_pro",
		Description: "测试高级写作功能",
	},

	// ============ 普通阅读用户（5个）============
	{
		Username:    "reader01",
		Email:       "reader01@qingyu.com",
		Password:    "Reader@123456",
		Role:        "user",
		Nickname:    "书虫小白",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=reader01",
		Description: "测试基础阅读功能",
	},
	{
		Username:    "reader02",
		Email:       "reader02@qingyu.com",
		Password:    "Reader@123456",
		Role:        "user",
		Nickname:    "阅读达人",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=reader02",
		Description: "测试书架管理功能",
	},
	{
		Username:    "reader03",
		Email:       "reader03@qingyu.com",
		Password:    "Reader@123456",
		Role:        "user",
		Nickname:    "小说爱好者",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=reader03",
		Description: "测试搜索和推荐功能",
	},
	{
		Username:    "reader04",
		Email:       "reader04@qingyu.com",
		Password:    "Reader@123456",
		Role:        "user",
		Nickname:    "评论家",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=reader04",
		Description: "测试评论和互动功能",
	},
	{
		Username:    "reader05",
		Email:       "reader05@qingyu.com",
		Password:    "Reader@123456",
		Role:        "user",
		Nickname:    "随性读者",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=reader05",
		Description: "测试用户行为追踪",
	},

	// ============ 综合测试用户（2个）============
	{
		Username:    "tester_all",
		Email:       "tester_all@qingyu.com",
		Password:    "Test@123456",
		Role:        "user",
		Nickname:    "全功能测试员",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=tester_all",
		Description: "测试所有用户功能",
	},
	{
		Username:    "tester_api",
		Email:       "tester_api@qingyu.com",
		Password:    "Test@123456",
		Role:        "user",
		Nickname:    "API测试专员",
		Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=tester_api",
		Description: "API接口测试专用账号",
	},
}

func main() {
	fmt.Println("====================================")
	fmt.Println("创建青羽内测用户")
	fmt.Println("====================================")
	fmt.Println()

	// 加载配置
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	// 初始化数据库
	if err := core.InitDB(); err != nil {
		log.Fatalf("❌ 初始化数据库失败: %v", err)
	}

	ctx := context.Background()

	// 检查数据库连接
	if global.DB == nil {
		log.Fatal("❌ 数据库未初始化")
	}

	// 获取用户集合
	userCollection := global.DB.Collection("users")

	fmt.Println("开始创建内测账号...")
	fmt.Println()

	successCount := 0
	skipCount := 0
	errorCount := 0

	for i, betaUser := range betaUsers {
		fmt.Printf("[%d/%d] 创建账号: %-20s (%s)...",
			i+1, len(betaUsers), betaUser.Username, betaUser.Nickname)

		// 检查用户是否已存在
		var existingUser users.User
		err := userCollection.FindOne(ctx, bson.M{
			"$or": []bson.M{
				{"username": betaUser.Username},
				{"email": betaUser.Email},
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
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(betaUser.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf(" [密码加密失败: %v]\n", err)
			errorCount++
			continue
		}

		// 创建用户
		now := time.Now()
		user := users.User{
			ID:        primitive.NewObjectID().Hex(),
			Username:  betaUser.Username,
			Email:     betaUser.Email,
			Password:  string(hashedPassword),
			Nickname:  betaUser.Nickname,
			Avatar:    betaUser.Avatar,
			Role:      betaUser.Role,
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

		fmt.Println(" ✓")
		successCount++
	}

	// 统计结果
	fmt.Println()
	fmt.Println("====================================")
	fmt.Println("账号创建完成")
	fmt.Println("====================================")
	fmt.Printf("✓ 成功创建: %d 个\n", successCount)
	fmt.Printf("⊘ 已存在跳过: %d 个\n", skipCount)
	if errorCount > 0 {
		fmt.Printf("✗ 失败: %d 个\n", errorCount)
	}
	fmt.Printf("总计: %d 个\n", len(betaUsers))
	fmt.Println()

	// 打印账号信息
	printAccountList()
}

// printAccountList 打印账号列表
func printAccountList() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("青羽内测账号列表")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 管理员账号
	fmt.Println("【管理员账号】- 拥有所有权限")
	fmt.Println("──────────────────────────────────────────────")
	for _, user := range betaUsers {
		if user.Role == "admin" {
			fmt.Printf("  用户名: %-20s 昵称: %s\n", user.Username, user.Nickname)
			fmt.Printf("  邮箱: %-30s 密码: %s\n", user.Email, user.Password)
			fmt.Printf("  说明: %s\n", user.Description)
			fmt.Println()
		}
	}

	// VIP用户
	fmt.Println("【VIP用户】- 享有高级功能权限")
	fmt.Println("──────────────────────────────────────────────")
	for _, user := range betaUsers {
		if user.Role == "vip" {
			fmt.Printf("  用户名: %-20s 昵称: %-15s 密码: %s\n",
				user.Username, user.Nickname, user.Password)
			fmt.Printf("  说明: %s\n", user.Description)
			fmt.Println()
		}
	}

	// 写作用户
	fmt.Println("【写作用户】- 测试写作相关功能")
	fmt.Println("──────────────────────────────────────────────")
	for _, user := range betaUsers {
		if user.Role == "user" && len(user.Username) >= 6 && (user.Username[:6] == "writer" || user.Username[:6] == "tester") {
			fmt.Printf("  %-20s | %-15s | %s\n",
				user.Username, user.Nickname, user.Password)
		}
	}
	fmt.Println()

	// 阅读用户
	fmt.Println("【阅读用户】- 测试阅读相关功能")
	fmt.Println("──────────────────────────────────────────────")
	for _, user := range betaUsers {
		if user.Role == "user" && len(user.Username) >= 6 && user.Username[:6] == "reader" {
			fmt.Printf("  %-20s | %-15s | %s\n",
				user.Username, user.Nickname, user.Password)
		}
	}
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 使用提示")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("1. 所有账号均已激活，可直接登录")
	fmt.Println("2. 管理员账号拥有所有权限，可进行系统管理")
	fmt.Println("3. VIP账号享有高级功能（AI辅助、高级搜索等）")
	fmt.Println("4. 普通用户账号用于测试标准功能")
	fmt.Println("5. 建议定期更换测试环境密码")
	fmt.Println("6. 生产环境请使用更强的密码策略")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 快速登录信息
	fmt.Println("🚀 快速登录")
	fmt.Println("──────────────────────────────────────────────")
	fmt.Println("管理员: admin / Admin@123456")
	fmt.Println("VIP作家: vip_writer01 / Vip@123456")
	fmt.Println("VIP读者: vip_reader01 / Vip@123456")
	fmt.Println("普通作家: writer_xuanhuan / Writer@123456")
	fmt.Println("普通读者: reader01 / Reader@123456")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
