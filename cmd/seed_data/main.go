//go:build !auto

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"
	"Qingyu_backend/migration/seeds"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║    青羽写作平台 - 测试数据更新工具     ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// 加载配置
	fmt.Println("📁 加载配置文件...")
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v\n", err)
	}
	fmt.Printf("✓ 配置加载成功\n")

	// 使用GetMongoConfig获取MongoDB配置
	mongoConfig, err := cfg.Database.GetMongoConfig()
	if err != nil {
		log.Fatalf("❌ 获取MongoDB配置失败: %v\n", err)
	}
	fmt.Printf("  数据库: %s\n", mongoConfig.URI)
	fmt.Printf("  数据库名: %s\n", mongoConfig.Database)
	fmt.Println()

	// 初始化数据库
	fmt.Println("🔗 连接数据库...")
	if err := core.InitDB(); err != nil {
		log.Fatalf("❌ 初始化数据库失败: %v\n", err)
	}
	fmt.Println("✓ 数据库连接成功")
	fmt.Println()

	ctx := context.Background()

	// 检查数据库连接
	if global.DB == nil {
		log.Fatal("❌ 数据库未初始化")
	}

	// 显示菜单
	for {
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("请选择操作：")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("1. 全部更新（清理旧数据 + 创建新数据）")
		fmt.Println("2. 仅创建新数据（跳过已存在的）")
		fmt.Println("3. 清理所有测试数据")
		fmt.Println("4. 查看数据统计")
		fmt.Println("5. 退出")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		var choice int
		fmt.Print("请输入选项 (1-5): ")
		fmt.Scanf("%d", &choice)

		fmt.Println()

		switch choice {
		case 1:
			cleanAllData(ctx)
			createAllData(ctx, global.DB)
		case 2:
			createAllData(ctx, global.DB)
		case 3:
			cleanAllData(ctx)
		case 4:
			showStatistics(ctx)
		case 5:
			fmt.Println("👋 再见！")
			os.Exit(0)
		default:
			fmt.Println("❌ 无效选项，请重新选择")
		}

		fmt.Println()
		fmt.Println("按 Enter 继续...")
		fmt.Scanln()
	}
}

func cleanAllData(ctx context.Context) {
	fmt.Println("🗑️  开始清理测试数据...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	collections := []string{
		"users",
		"books",
		"chapters",
		"chapter_contents",
		"wallets",
		"transactions",
		"comments",
		"likes",
		"collections",
		"follows",
	}

	totalDeleted := 0

	for _, collName := range collections {
		collection := global.DB.Collection(collName)
		result, err := collection.DeleteMany(ctx, bson.M{})
		if err != nil {
			fmt.Printf("❌ 清理 %s 失败: %v\n", collName, err)
			continue
		}
		totalDeleted += int(result.DeletedCount)
		if result.DeletedCount > 0 {
			fmt.Printf("✓ 清理 %s: %d 条\n", collName, result.DeletedCount)
		}
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("✓ 总计清理 %d 条数据\n", totalDeleted)
	fmt.Println()
}

func createAllData(ctx context.Context, db *mongo.Database) {
	fmt.Println("✨ 开始创建测试数据...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 1. 创建用户
	fmt.Println("【1/5】创建用户数据")
	if err := seeds.SeedEnhancedUsers(ctx, db); err != nil {
		fmt.Printf("❌ 创建用户失败: %v\n", err)
	} else {
		fmt.Println("✓ 用户数据创建完成")
	}
	fmt.Println()

	// 2. 创建书籍
	fmt.Println("【2/5】创建书籍数据")
	if err := seeds.SeedBooks(ctx, db); err != nil {
		fmt.Printf("❌ 创建书籍失败: %v\n", err)
	} else {
		fmt.Println("✓ 书籍数据创建完成")
	}
	fmt.Println()

	// 3. 创建章节
	fmt.Println("【3/5】创建章节数据")
	if err := seeds.SeedChapters(ctx, db); err != nil {
		fmt.Printf("❌ 创建章节失败: %v\n", err)
	} else {
		fmt.Println("✓ 章节数据创建完成")
	}
	fmt.Println()

	// 4. 创建钱包
	fmt.Println("【4/5】创建钱包数据")
	if err := seeds.SeedWallets(ctx, db); err != nil {
		fmt.Printf("❌ 创建钱包失败: %v\n", err)
	} else {
		fmt.Println("✓ 钱包数据创建完成")
	}
	fmt.Println()

	// 5. 创建社交数据
	fmt.Println("【5/5】创建社交数据")
	if err := seeds.SeedSocialData(ctx, db); err != nil {
		fmt.Printf("❌ 创建社交数据失败: %v\n", err)
	} else {
		fmt.Println("✓ 社交数据创建完成")
	}
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ 所有测试数据创建完成！")
	fmt.Println()

	// 打印测试账号
	printTestAccounts()
}

func showStatistics(ctx context.Context) {
	fmt.Println("📊 数据统计")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	collections := []struct {
		Name  string
		Alias string
	}{
		{"users", "用户"},
		{"books", "书籍"},
		{"chapters", "章节"},
		{"chapter_contents", "章节内容"},
		{"wallets", "钱包"},
		{"transactions", "交易记录"},
		{"comments", "评论"},
		{"likes", "点赞"},
		{"collections", "收藏"},
		{"follows", "关注"},
	}

	total := 0
	for _, coll := range collections {
		count, err := global.DB.Collection(coll.Name).CountDocuments(ctx, bson.M{})
		if err != nil {
			fmt.Printf("❌ 统计 %s 失败: %v\n", coll.Alias, err)
			continue
		}
		total += int(count)
		fmt.Printf("  %s: %d 条\n", coll.Alias, count)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  总计: %d 条数据\n", total)
	fmt.Println()
}

func printTestAccounts() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📝 测试账号列表")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	users := seeds.GetTestUserCredentials()

	// 管理员账号
	fmt.Println("【管理员账号】")
	adminCount := 0
	for _, user := range users {
		if user.Role == "admin" {
			fmt.Printf("  用户名: %s\n", user.Username)
			fmt.Printf("  邮箱: %s\n", user.Email)
			fmt.Printf("  密码: %s\n", user.Password)
			fmt.Printf("  说明: %s\n", user.Description)
			fmt.Println()
			adminCount++
			if adminCount >= 2 {
				break
			}
		}
	}

	// 作者账号
	fmt.Println("【作者账号】")
	authorCount := 0
	for _, user := range users {
		if user.Role == "author" {
			fmt.Printf("  %s | %s\n", user.Username, user.Password)
			authorCount++
			if authorCount >= 3 {
				break
			}
		}
	}
	fmt.Println()

	// VIP读者
	fmt.Println("【VIP读者】")
	vipCount := 0
	for _, user := range users {
		if user.Role == "reader" && user.Email[:3] == "vip" {
			fmt.Printf("  %s | %s\n", user.Username, user.Password)
			vipCount++
			if vipCount >= 3 {
				break
			}
		}
	}
	fmt.Println()

	// 普通读者
	fmt.Println("【普通读者】")
	readerCount := 0
	for _, user := range users {
		if user.Role == "reader" && user.Email[:6] == "reader" {
			fmt.Printf("  %s | %s\n", user.Username, user.Password)
			readerCount++
			if readerCount >= 3 {
				break
			}
		}
	}
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 提示：")
	fmt.Println("1. 所有账号可直接登录")
	fmt.Println("2. 建议定期更换测试密码")
	fmt.Println("3. 生产环境请使用强密码")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
