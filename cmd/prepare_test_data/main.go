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

	"Qingyu_backend/config"
)

// 此文件用于准备测试数据，包括测试用户和AI配额

func main() {
	fmt.Println("=========================================")
	fmt.Println("  测试数据准备工具")
	fmt.Println("=========================================")
	fmt.Println()

	// 加载配置
	fmt.Println("1. 加载配置...")
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}
	fmt.Println("✓ 配置加载成功")
	fmt.Println()

	// 连接MongoDB
	fmt.Println("2. 连接MongoDB...")
	ctx := context.Background()

	// 获取MongoDB配置
	mongoConfig := config.GlobalConfig.Database.Primary.MongoDB
	if mongoConfig == nil {
		log.Fatal("❌ MongoDB配置未找到")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConfig.URI))
	if err != nil {
		log.Fatalf("❌ 连接MongoDB失败: %v", err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ MongoDB ping失败: %v", err)
	}
	fmt.Println("✓ MongoDB连接成功")
	fmt.Println()

	db := client.Database(mongoConfig.Database)

	// 检查并创建测试用户
	fmt.Println("3. 检查测试用户...")
	err = ensureTestUsers(ctx, db)
	if err != nil {
		log.Fatalf("❌ 创建测试用户失败: %v", err)
	}
	fmt.Println("✓ 测试用户准备完成")
	fmt.Println()

	// 激活AI配额
	fmt.Println("4. 清理并激活AI配额...")

	// 先删除所有旧配额记录
	quotas := db.Collection("ai_quotas")
	result, _ := quotas.DeleteMany(ctx, bson.M{})
	if result.DeletedCount > 0 {
		fmt.Printf("   已删除 %d 条旧配额记录\n", result.DeletedCount)
	}

	err = activateAIQuotas(ctx, db)
	if err != nil {
		log.Fatalf("❌ 激活AI配额失败: %v", err)
	}
	fmt.Println("✓ AI配额激活成功")
	fmt.Println()

	// 检查分类数据
	fmt.Println("5. 检查分类数据...")
	err = ensureCategories(ctx, db)
	if err != nil {
		log.Fatalf("❌ 创建分类失败: %v", err)
	}
	fmt.Println("✓ 分类数据准备完成")
	fmt.Println()

	// 统计数据
	fmt.Println("=========================================")
	fmt.Println("  ✓ 测试数据准备完成")
	fmt.Println("=========================================")
	fmt.Println()

	printStats(ctx, db)

	fmt.Println()
	fmt.Println("可以开始运行测试了！")
	fmt.Println("  go test ./test/integration -v -count=1")
	fmt.Println()
}

func ensureTestUsers(ctx context.Context, db *mongo.Database) error {
	users := db.Collection("users")

	// 检查test_user01
	var user bson.M
	err := users.FindOne(ctx, bson.M{"username": "test_user01"}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		fmt.Println("   ⚠ 测试用户不存在，请先运行: go run cmd/create_beta_users/main.go")
		return fmt.Errorf("测试用户不存在")
	} else if err != nil {
		return err
	}

	fmt.Println("   ✓ 测试用户已存在")
	return nil
}

func activateAIQuotas(ctx context.Context, db *mongo.Database) error {
	users := db.Collection("users")
	quotas := db.Collection("ai_quotas")

	// 获取测试用户ID
	testUsers := []string{"test_user01", "test_user02", "test_user03", "vip_user01", "vip_user02"}

	for _, username := range testUsers {
		var user bson.M
		err := users.FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil {
			fmt.Printf("   ! 跳过用户 %s (不存在)\n", username)
			continue
		}

		// 处理_id字段（统一转换为string，因为数据库中存储的是string）
		var userIDStr string
		switch v := user["_id"].(type) {
		case primitive.ObjectID:
			userIDStr = v.Hex()
		case string:
			userIDStr = v
		default:
			fmt.Printf("   ! 跳过用户 %s (ID类型不支持: %T)\n", username, v)
			continue
		}

		// 检查配额是否存在
		var quota bson.M
		err = quotas.FindOne(ctx, bson.M{"user_id": userIDStr}).Decode(&quota)

		now := time.Now()

		// VIP用户和管理员获得超大配额（模拟无限制）
		var totalQuota, remainingQuota int
		if username == "vip_user01" || username == "vip_user02" {
			totalQuota = 999999 // VIP用户：百万级配额
			remainingQuota = 999999
		} else {
			totalQuota = 100000 // 普通测试用户：十万级配额
			remainingQuota = 100000
		}

		// 设置重置时间：明天（daily配额每天重置）
		resetAt := now.AddDate(0, 0, 1) // 明天同一时间

		quotaData := bson.M{
			"user_id":         userIDStr,
			"quota_type":      "daily",        // daily quota
			"total_quota":     totalQuota,     // 大配额
			"used_quota":      0,              // haven't used any
			"remaining_quota": remainingQuota, // all available
			"status":          "active",       // active status
			"reset_at":        resetAt,        // 重置时间（关键！）
			"last_reset_date": now,
			"updated_at":      now,
			// Also set monthly fields for reference
			"monthly_limit":   totalQuota * 30, // 月配额更大
			"used_this_month": 0,
		}

		if err == mongo.ErrNoDocuments {
			// 创建新配额
			quotaData["created_at"] = now
			_, err = quotas.InsertOne(ctx, quotaData)
			if err != nil {
				return fmt.Errorf("创建配额失败 (%s): %v", username, err)
			}
			fmt.Printf("   ✓ %s: AI配额已创建并激活\n", username)
		} else {
			// 更新现有配额
			_, err = quotas.UpdateOne(
				ctx,
				bson.M{"user_id": userIDStr},
				bson.M{"$set": quotaData},
			)
			if err != nil {
				return fmt.Errorf("更新配额失败 (%s): %v", username, err)
			}
			fmt.Printf("   ✓ %s: AI配额已更新并激活\n", username)
		}
	}

	return nil
}

func ensureCategories(ctx context.Context, db *mongo.Database) error {
	categories := db.Collection("categories")

	// 检查分类数量
	count, err := categories.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count >= 5 {
		fmt.Printf("   ✓ 分类数据充足 (%d个)\n", count)
		return nil
	}

	fmt.Println("   正在创建基础分类...")

	defaultCategories := []bson.M{
		{
			"name":        "玄幻",
			"slug":        "xuanhuan",
			"description": "玄幻小说",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  1,
			"is_active":   true,
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		},
		{
			"name":        "都市",
			"slug":        "dushi",
			"description": "都市小说",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  2,
			"is_active":   true,
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		},
		{
			"name":        "仙侠",
			"slug":        "xianxia",
			"description": "仙侠小说",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  3,
			"is_active":   true,
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		},
		{
			"name":        "科幻",
			"slug":        "kehuan",
			"description": "科幻小说",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  4,
			"is_active":   true,
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		},
		{
			"name":        "历史",
			"slug":        "lishi",
			"description": "历史小说",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  5,
			"is_active":   true,
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		},
	}

	docs := make([]interface{}, len(defaultCategories))
	for i, cat := range defaultCategories {
		docs[i] = cat
	}

	_, err = categories.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	fmt.Printf("   ✓ 创建了 %d 个分类\n", len(defaultCategories))
	return nil
}

func printStats(ctx context.Context, db *mongo.Database) {
	fmt.Println("数据统计:")

	// 书籍
	count, _ := db.Collection("books").CountDocuments(ctx, bson.M{})
	fmt.Printf("  - 书籍: %d 本\n", count)

	// 章节
	count, _ = db.Collection("chapters").CountDocuments(ctx, bson.M{})
	fmt.Printf("  - 章节: %d 个\n", count)

	// 用户
	count, _ = db.Collection("users").CountDocuments(ctx, bson.M{
		"username": bson.M{"$regex": "^(test_user|vip_user)"},
	})
	fmt.Printf("  - 测试用户: %d 个\n", count)

	// 分类
	count, _ = db.Collection("categories").CountDocuments(ctx, bson.M{})
	fmt.Printf("  - 分类: %d 个\n", count)

	// AI配额
	count, _ = db.Collection("ai_quotas").CountDocuments(ctx, bson.M{"status": "active"})
	fmt.Printf("  - 激活的AI配额: %d 个\n", count)
}
