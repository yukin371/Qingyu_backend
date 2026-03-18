package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/migration/seeds"
	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/users"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func main() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("         青羽写作平台 - 内测数据初始化工具")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 1. 加载配置
	fmt.Println("【步骤 1/6】加载配置...")
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}
	fmt.Println("✓ 配置加载成功")
	fmt.Println()

	// 2. 连接数据库
	fmt.Println("【步骤 2/6】连接数据库...")
	db, err := connectDB()
	if err != nil {
		log.Fatalf("❌ 连接数据库失败: %v", err)
	}
	fmt.Println("✓ 数据库连接成功")
	fmt.Println()

	ctx := context.Background()

	// 3. 创建内测用户
	fmt.Println("【步骤 3/6】创建内测用户...")
	err = createBetaUsers(ctx, db)
	if err != nil {
		log.Fatalf("❌ 创建内测用户失败: %v", err)
	}
	fmt.Println("✓ 内测用户创建完成")
	fmt.Println()

	// 4. 创建分类数据
	fmt.Println("【步骤 4/6】创建分类数据...")
	err = createCategories(ctx, db)
	if err != nil {
		log.Fatalf("❌ 创建分类失败: %v", err)
	}
	fmt.Println("✓ 分类数据创建完成")
	fmt.Println()

	// 5. 导入书籍数据
	fmt.Println("【步骤 5/6】导入书籍数据...")
	fmt.Println("提示：如果已有 data/novels_100.json 文件，将自动导入")
	fmt.Println("     如果没有，请先运行: python scripts/data/import_novels.py --max-novels 100 --output data/novels_100.json")
	err = importBooks(ctx, db)
	if err != nil {
		log.Printf("⚠ 书籍导入失败或跳过: %v", err)
		fmt.Println("  提示: 你可以稍后手动运行书籍导入命令")
	} else {
		fmt.Println("✓ 书籍数据导入完成")
	}
	fmt.Println()

	// 6. 创建轮播图
	fmt.Println("【步骤 6/6】创建轮播图...")
	err = createBanners(ctx, db)
	if err != nil {
		log.Printf("⚠ 轮播图创建失败: %v", err)
		fmt.Println("  提示: 需要先有书籍数据才能创建轮播图")
	} else {
		fmt.Println("✓ 轮播图创建完成")
	}
	fmt.Println()

	// 7. 激活AI配额
	fmt.Println("【步骤 7/7】激活AI配额...")
	err = activateAIQuotas(ctx, db)
	if err != nil {
		log.Fatalf("❌ 激活AI配额失败: %v", err)
	}
	fmt.Println("✓ AI配额激活完成")
	fmt.Println()

	// 显示统计信息
	printFinalStats(ctx, db)

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("         ✓ 内测数据初始化完成！")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	printQuickStartGuide()
}

// connectDB 连接数据库
func connectDB() (*mongo.Database, error) {
	dbConfig := config.GlobalConfig.Database
	if dbConfig == nil {
		return nil, fmt.Errorf("数据库配置未找到")
	}

	mongoConfig, err := dbConfig.GetMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("获取MongoDB配置失败: %w", err)
	}

	clientOptions := options.Client().ApplyURI(mongoConfig.URI)
	clientOptions.SetMaxPoolSize(mongoConfig.MaxPoolSize)
	clientOptions.SetMinPoolSize(mongoConfig.MinPoolSize)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("连接失败: %w", err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("ping失败: %w", err)
	}

	return client.Database(mongoConfig.Database), nil
}

// createBetaUsers 创建内测用户
func createBetaUsers(ctx context.Context, db *mongo.Database) error {
	userCollection := db.Collection("users")

	betaUsers := []BetaUser{
		// 管理员账号（3个）
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
		// VIP用户（5个）
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
		// 普通写作用户（5个）
		{
			Username:    "writer_xuanhuan",
			Email:       "writer_xuanhuan@qingyu.com",
			Password:    "Writer@123456",
			Role:        "reader",
			Nickname:    "玄幻小说家",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=writer_xuanhuan",
			Description: "测试玄幻小说创作功能",
		},
		{
			Username:    "writer_yanqing",
			Email:       "writer_yanqing@qingyu.com",
			Password:    "Writer@123456",
			Role:        "reader",
			Nickname:    "言情作家",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=writer_yanqing",
			Description: "测试言情小说创作功能",
		},
		{
			Username:    "writer_dushi",
			Email:       "writer_dushi@qingyu.com",
			Password:    "Writer@123456",
			Role:        "reader",
			Nickname:    "都市创作者",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=writer_dushi",
			Description: "测试都市小说创作功能",
		},
		{
			Username:    "writer_newbie",
			Email:       "writer_newbie@qingyu.com",
			Password:    "Writer@123456",
			Role:        "reader",
			Nickname:    "新手作者",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=writer_newbie",
			Description: "测试新手写作流程",
		},
		{
			Username:    "writer_pro",
			Email:       "writer_pro@qingyu.com",
			Password:    "Writer@123456",
			Role:        "reader",
			Nickname:    "专业作家",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=writer_pro",
			Description: "测试高级写作功能",
		},
		// 普通阅读用户（5个）
		{
			Username:    "reader01",
			Email:       "reader01@qingyu.com",
			Password:    "Reader@123456",
			Role:        "reader",
			Nickname:    "书虫小白",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=reader01",
			Description: "测试基础阅读功能",
		},
		{
			Username:    "reader02",
			Email:       "reader02@qingyu.com",
			Password:    "Reader@123456",
			Role:        "reader",
			Nickname:    "阅读达人",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=reader02",
			Description: "测试书架管理功能",
		},
		{
			Username:    "reader03",
			Email:       "reader03@qingyu.com",
			Password:    "Reader@123456",
			Role:        "reader",
			Nickname:    "小说爱好者",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=reader03",
			Description: "测试搜索和推荐功能",
		},
		{
			Username:    "reader04",
			Email:       "reader04@qingyu.com",
			Password:    "Reader@123456",
			Role:        "reader",
			Nickname:    "评论家",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=reader04",
			Description: "测试评论和互动功能",
		},
		{
			Username:    "reader05",
			Email:       "reader05@qingyu.com",
			Password:    "Reader@123456",
			Role:        "reader",
			Nickname:    "随性读者",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=reader05",
			Description: "测试用户行为追踪",
		},
		// 综合测试用户（2个）
		{
			Username:    "tester_all",
			Email:       "tester_all@qingyu.com",
			Password:    "Test@123456",
			Role:        "reader",
			Nickname:    "全功能测试员",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=tester_all",
			Description: "测试所有用户功能",
		},
		{
			Username:    "tester_api",
			Email:       "tester_api@qingyu.com",
			Password:    "Test@123456",
			Role:        "reader",
			Nickname:    "API测试专员",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=tester_api",
			Description: "API接口测试专用账号",
		},
	}

	successCount := 0
	skipCount := 0

	for _, betaUser := range betaUsers {
		// 检查用户是否已存在
		var existingUser users.User
		err := userCollection.FindOne(ctx, bson.M{
			"$or": []bson.M{
				{"username": betaUser.Username},
				{"email": betaUser.Email},
			},
		}).Decode(&existingUser)

		if err == nil {
			skipCount++
			continue
		} else if err != mongo.ErrNoDocuments {
			return fmt.Errorf("查询用户失败: %w", err)
		}

		// 密码加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(betaUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("密码加密失败: %w", err)
		}

		// 创建用户
		now := time.Now()
		user := users.User{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
			Username:         betaUser.Username,
			Email:            betaUser.Email,
			Password:         string(hashedPassword),
			Nickname:         betaUser.Nickname,
			Avatar:           betaUser.Avatar,
			Roles:            []string{betaUser.Role},
			Status:           "active",
		}

		_, err = userCollection.InsertOne(ctx, user)
		if err != nil {
			return fmt.Errorf("创建用户失败 (%s): %w", betaUser.Username, err)
		}

		successCount++
	}

	fmt.Printf("  创建成功: %d 个用户\n", successCount)
	fmt.Printf("  已存在跳过: %d 个用户\n", skipCount)

	return nil
}

// createCategories 创建分类数据
func createCategories(ctx context.Context, db *mongo.Database) error {
	categories := db.Collection("categories")

	// 检查分类数量
	count, err := categories.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count >= 5 {
		fmt.Printf("  分类数据已存在 (%d个)，跳过创建\n", count)
		return nil
	}

	now := time.Now()
	defaultCategories := []interface{}{
		bson.M{
			"name":        "玄幻",
			"slug":        "xuanhuan",
			"description": "玄幻小说",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  1,
			"is_active":   true,
			"created_at":  now,
			"updated_at":  now,
		},
		bson.M{
			"name":        "都市",
			"slug":        "dushi",
			"description": "都市小说",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  2,
			"is_active":   true,
			"created_at":  now,
			"updated_at":  now,
		},
		bson.M{
			"name":        "仙侠",
			"slug":        "xianxia",
			"description": "仙侠小说",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  3,
			"is_active":   true,
			"created_at":  now,
			"updated_at":  now,
		},
		bson.M{
			"name":        "科幻",
			"slug":        "kehuan",
			"description": "科幻小说",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  4,
			"is_active":   true,
			"created_at":  now,
			"updated_at":  now,
		},
		bson.M{
			"name":        "历史",
			"slug":        "lishi",
			"description": "历史小说",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  5,
			"is_active":   true,
			"created_at":  now,
			"updated_at":  now,
		},
		bson.M{
			"name":        "言情",
			"slug":        "yanqing",
			"description": "言情小说",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  6,
			"is_active":   true,
			"created_at":  now,
			"updated_at":  now,
		},
		bson.M{
			"name":        "奇幻",
			"slug":        "qihuan",
			"description": "奇幻小说",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  7,
			"is_active":   true,
			"created_at":  now,
			"updated_at":  now,
		},
		bson.M{
			"name":        "其他",
			"slug":        "other",
			"description": "其他类型",
			"parent_id":   nil,
			"level":       1,
			"sort_order":  99,
			"is_active":   true,
			"created_at":  now,
			"updated_at":  now,
		},
	}

	_, err = categories.InsertMany(ctx, defaultCategories)
	if err != nil {
		return fmt.Errorf("插入分类失败: %w", err)
	}

	fmt.Printf("  创建了 %d 个分类\n", len(defaultCategories))
	return nil
}

// importBooks 导入书籍数据
func importBooks(ctx context.Context, db *mongo.Database) error {
	// 检查是否已有书籍
	bookCollection := db.Collection("books")
	count, err := bookCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("查询书籍失败: %w", err)
	}

	if count > 0 {
		fmt.Printf("  数据库已有 %d 本书籍，跳过导入\n", count)
		return nil
	}

	// 尝试从 JSON 文件导入
	filepath := "data/novels_100.json"
	importer := seeds.NewNovelImporter(db, false)

	fmt.Printf("  正在从 %s 导入书籍...\n", filepath)
	err = importer.ImportFromJSON(ctx, filepath)
	if err != nil {
		return fmt.Errorf("导入书籍失败: %w", err)
	}

	// 创建索引
	err = importer.CreateIndexes(ctx)
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	return nil
}

// createBanners 创建轮播图
func createBanners(ctx context.Context, db *mongo.Database) error {
	bannerCollection := db.Collection("banners")
	bookCollection := db.Collection("books")

	// 检查是否已有轮播图
	count, err := bannerCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("查询轮播图失败: %w", err)
	}

	if count > 0 {
		fmt.Printf("  轮播图已存在 (%d个)，跳过创建\n", count)
		return nil
	}

	// 获取高评分书籍
	cursor, err := bookCollection.Find(ctx, bson.M{},
		options.Find().SetLimit(5).SetSort(bson.D{{Key: "word_count", Value: -1}}))
	if err != nil {
		return fmt.Errorf("查询书籍失败: %w", err)
	}
	defer cursor.Close(ctx)

	var books []struct {
		ID     primitive.ObjectID `bson:"_id"`
		Title  string             `bson:"title"`
		Cover  string             `bson:"cover"`
		Rating float64            `bson:"rating"`
	}
	if err = cursor.All(ctx, &books); err != nil {
		return fmt.Errorf("读取书籍失败: %w", err)
	}

	if len(books) == 0 {
		return fmt.Errorf("没有找到书籍，无法创建轮播图")
	}

	now := time.Now()
	banners := []interface{}{}

	bannerTitles := []string{"精品推荐", "热门连载", "编辑精选", "新书上架", "人气爆款"}
	bannerDescs := []string{
		"高分力作，不容错过！",
		"最热门的小说，千万读者的选择",
		"编辑部精心挑选的优质作品",
		"新鲜出炉的精彩小说，抢先阅读",
		"超人气作品，口碑爆棚",
	}

	for i, book := range books {
		startTime := now.Add(time.Duration(-7+i) * 24 * time.Hour)
		endTime := now.Add(30 * 24 * time.Hour)

		banner := &bookstore.Banner{
			ID:          primitive.NewObjectID(),
			Title:       bannerTitles[i] + "：" + book.Title,
			Description: bannerDescs[i],
			Image:       book.Cover,
			Target:      book.ID.Hex(),
			TargetType:  "book",
			SortOrder:   i + 1,
			IsActive:    true,
			StartTime:   &startTime,
			EndTime:     &endTime,
			ClickCount:  150 + int64(i*50),
			CreatedAt:   startTime,
			UpdatedAt:   now,
		}
		banners = append(banners, banner)
	}

	_, err = bannerCollection.InsertMany(ctx, banners)
	if err != nil {
		return fmt.Errorf("插入轮播图失败: %w", err)
	}

	fmt.Printf("  创建了 %d 个轮播图\n", len(banners))
	return nil
}

// activateAIQuotas 激活AI配额
func activateAIQuotas(ctx context.Context, db *mongo.Database) error {
	users := db.Collection("users")
	quotas := db.Collection("ai_quotas")

	// 删除所有旧配额
	_, _ = quotas.DeleteMany(ctx, bson.M{})

	testUsers := []string{
		"admin", "admin_test", "admin_audit",
		"vip_writer01", "vip_writer02", "vip_reader01", "vip_reader02", "vip_tester",
		"writer_xuanhuan", "writer_yanqing", "writer_dushi", "writer_newbie", "writer_pro",
		"reader01", "reader02", "reader03", "reader04", "reader05",
		"tester_all", "tester_api",
	}

	activatedCount := 0
	for _, username := range testUsers {
		var user bson.M
		err := users.FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil {
			continue
		}

		var userIDStr string
		switch v := user["_id"].(type) {
		case primitive.ObjectID:
			userIDStr = v.Hex()
		case string:
			userIDStr = v
		default:
			continue
		}

		now := time.Now()
		resetAt := now.AddDate(0, 0, 1)

		// 根据角色设置配额
		var totalQuota int
		role, _ := user["role"].(string)
		if role == "admin" || role == "vip" {
			totalQuota = 999999 // 管理员和VIP：超大配额
		} else {
			totalQuota = 100000 // 普通用户：十万配额
		}

		quotaData := bson.M{
			"user_id":         userIDStr,
			"quota_type":      "daily",
			"total_quota":     totalQuota,
			"used_quota":      0,
			"remaining_quota": totalQuota,
			"status":          "active",
			"reset_at":        resetAt,
			"last_reset_date": now,
			"created_at":      now,
			"updated_at":      now,
			"monthly_limit":   totalQuota * 30,
			"used_this_month": 0,
		}

		_, err = quotas.InsertOne(ctx, quotaData)
		if err != nil {
			continue
		}

		activatedCount++
	}

	fmt.Printf("  激活了 %d 个用户的AI配额\n", activatedCount)
	return nil
}

// printFinalStats 打印最终统计
func printFinalStats(ctx context.Context, db *mongo.Database) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("数据统计")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 用户统计
	userCount, _ := db.Collection("users").CountDocuments(ctx, bson.M{})
	adminCount, _ := db.Collection("users").CountDocuments(ctx, bson.M{"role": "admin"})
	vipCount, _ := db.Collection("users").CountDocuments(ctx, bson.M{"role": "vip"})
	normalCount, _ := db.Collection("users").CountDocuments(ctx, bson.M{"roles": "reader"})

	fmt.Printf("【用户数据】\n")
	fmt.Printf("  总用户数: %d\n", userCount)
	fmt.Printf("  - 管理员: %d\n", adminCount)
	fmt.Printf("  - VIP用户: %d\n", vipCount)
	fmt.Printf("  - 普通用户: %d\n", normalCount)
	fmt.Println()

	// 分类统计
	categoryCount, _ := db.Collection("categories").CountDocuments(ctx, bson.M{})
	fmt.Printf("【分类数据】\n")
	fmt.Printf("  分类总数: %d\n", categoryCount)
	fmt.Println()

	// 书籍统计
	bookCount, _ := db.Collection("books").CountDocuments(ctx, bson.M{})
	chapterCount, _ := db.Collection("chapters").CountDocuments(ctx, bson.M{})
	fmt.Printf("【书籍数据】\n")
	fmt.Printf("  书籍总数: %d\n", bookCount)
	fmt.Printf("  章节总数: %d\n", chapterCount)
	fmt.Println()

	// 轮播图统计
	bannerCount, _ := db.Collection("banners").CountDocuments(ctx, bson.M{})
	fmt.Printf("【轮播图】\n")
	fmt.Printf("  轮播图数: %d\n", bannerCount)
	fmt.Println()

	// AI配额统计
	quotaCount, _ := db.Collection("ai_quotas").CountDocuments(ctx, bson.M{"status": "active"})
	fmt.Printf("【AI配额】\n")
	fmt.Printf("  激活配额: %d\n", quotaCount)
}

// printQuickStartGuide 打印快速开始指南
func printQuickStartGuide() {
	fmt.Println("🚀 快速开始")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("【测试账号】（用户名 / 密码）")
	fmt.Println()
	fmt.Println("管理员账号：")
	fmt.Println("  admin / Admin@123456")
	fmt.Println()
	fmt.Println("VIP作家账号：")
	fmt.Println("  vip_writer01 / Vip@123456")
	fmt.Println("  vip_writer02 / Vip@123456")
	fmt.Println()
	fmt.Println("VIP读者账号：")
	fmt.Println("  vip_reader01 / Vip@123456")
	fmt.Println("  vip_reader02 / Vip@123456")
	fmt.Println()
	fmt.Println("普通作家账号：")
	fmt.Println("  writer_xuanhuan / Writer@123456")
	fmt.Println("  writer_yanqing / Writer@123456")
	fmt.Println()
	fmt.Println("普通读者账号：")
	fmt.Println("  reader01 / Reader@123456")
	fmt.Println("  reader02 / Reader@123456")
	fmt.Println()
	fmt.Println("测试专员账号：")
	fmt.Println("  tester_all / Test@123456")
	fmt.Println("  tester_api / Test@123456")
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("【下一步】")
	fmt.Println("1. 启动服务器: go run cmd/server/main.go")
	fmt.Println("2. 访问 API: http://localhost:9090")
	fmt.Println("3. 使用上述账号登录测试")
	fmt.Println()
	fmt.Println("【如需导入更多书籍】")
	fmt.Println("运行: python scripts/data/import_novels.py --max-novels 100 --output data/novels_100.json")
	fmt.Println("然后: go run cmd/migrate/main.go -command=import-novels -file=data/novels_100.json")
	fmt.Println()
}
