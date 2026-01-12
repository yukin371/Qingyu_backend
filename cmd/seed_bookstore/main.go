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
)

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║      青羽书店 - 测试数据填充工具        ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// 连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	fmt.Println("🔗 数据库连接成功")
	fmt.Println()

	// 1. 创建分类
	fmt.Println("【1/4】创建分类数据")
	if err := seedCategories(ctx, db); err != nil {
		log.Printf("创建分类失败: %v", err)
	} else {
		fmt.Println("✓ 分类数据创建完成")
	}
	fmt.Println()

	// 2. 创建 Banner
	fmt.Println("【2/4】创建 Banner 数据")
	if err := seedBanners(ctx, db); err != nil {
		log.Printf("创建 Banner 失败: %v", err)
	} else {
		fmt.Println("✓ Banner 数据创建完成")
	}
	fmt.Println()

	// 3. 创建榜单
	fmt.Println("【3/4】创建榜单数据")
	if err := seedRankings(ctx, db); err != nil {
		log.Printf("创建榜单失败: %v", err)
	} else {
		fmt.Println("✓ 榜单数据创建完成")
	}
	fmt.Println()

	// 4. 显示统计
	fmt.Println("【4/4】数据统计")
	showStatistics(ctx, db)
	fmt.Println()

	fmt.Println("✨ 所有测试数据填充完成！")
}

// Category 分类结构
type Category struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	Name        string              `bson:"name"`
	Description string              `bson:"description"`
	Icon        string              `bson:"icon"`
	ParentID    *primitive.ObjectID `bson:"parent_id,omitempty"`
	Level       int                 `bson:"level"`
	SortOrder   int                 `bson:"sort_order"`
	BookCount   int64               `bson:"book_count"`
	IsActive    bool                `bson:"is_active"`
	CreatedAt   time.Time           `bson:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at"`
}

func seedCategories(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("categories")

	// 检查是否已有数据
	count, _ := collection.CountDocuments(ctx, bson.M{})
	if count > 0 {
		fmt.Printf("  已存在 %d 条分类数据，跳过\n", count)
		return nil
	}

	now := time.Now()

	// 一级分类
	categories := []Category{
		{
			Name:        "玄幻",
			Description: "东方玄幻、异世大陆、高武世界",
			Icon:        "/icons/xuanhuan.png",
			Level:       0,
			SortOrder:   1,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "都市",
			Description: "都市生活、都市异能、恋爱日常",
			Icon:        "/icons/dushi.png",
			Level:       0,
			SortOrder:   2,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "仙侠",
			Description: "古典仙侠、现代修真、洪荒封神",
			Icon:        "/icons/xianxia.png",
			Level:       0,
			SortOrder:   3,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "科幻",
			Description: "未来世界、星际战争、时空穿梭",
			Icon:        "/icons/kehuan.png",
			Level:       0,
			SortOrder:   4,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "武侠",
			Description: "传统武侠、武侠幻想、国术无双",
			Icon:        "/icons/wuxia.png",
			Level:       0,
			SortOrder:   5,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "历史",
			Description: "架空历史、历史传记、两晋隋唐",
			Icon:        "/icons/lishi.png",
			Level:       0,
			SortOrder:   6,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "游戏",
			Description: "虚拟网游、电子竞技、游戏异界",
			Icon:        "/icons/youxi.png",
			Level:       0,
			SortOrder:   7,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "奇幻",
			Description: "西方奇幻、剑与魔法、黑暗幻想",
			Icon:        "/icons/qihuan.png",
			Level:       0,
			SortOrder:   8,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	docs := make([]interface{}, len(categories))
	for i, cat := range categories {
		docs[i] = cat
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("插入分类失败: %w", err)
	}

	fmt.Printf("  创建 %d 条分类数据\n", len(result.InsertedIDs))
	return nil
}

// Banner Banner 结构
type Banner struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Image       string             `bson:"image"`
	Target      string             `bson:"target"`
	TargetType  string             `bson:"target_type"`
	SortOrder   int                `bson:"sort_order"`
	IsActive    bool               `bson:"is_active"`
	StartTime   *time.Time         `bson:"start_time,omitempty"`
	EndTime     *time.Time         `bson:"end_time,omitempty"`
	ClickCount  int64              `bson:"click_count"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func seedBanners(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("banners")

	// 检查是否已有数据
	count, _ := collection.CountDocuments(ctx, bson.M{})
	if count > 0 {
		fmt.Printf("  已存在 %d 条 Banner 数据，跳过\n", count)
		return nil
	}

	now := time.Now()

	banners := []Banner{
		{
			Title:       "修真世界 - 热门推荐",
			Description: "凡人流经典之作，不可错过",
			Image:       "https://images.unsplash.com/photo-1518709268805-4e9042af9f23?w=1200&h=400&fit=crop",
			Target:      "6956392cfe350a59abae6607",
			TargetType:  "book",
			SortOrder:   1,
			IsActive:    true,
			ClickCount:  0,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Title:       "新人福利大礼包",
			Description: "注册即送 100 青羽币",
			Image:       "https://images.unsplash.com/photo-1557683316-973673baf926?w=1200&h=400&fit=crop",
			Target:      "/promo/newbie",
			TargetType:  "url",
			SortOrder:   2,
			IsActive:    true,
			ClickCount:  0,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Title:       "本周排行榜",
			Description: "查看最热门的小说",
			Image:       "https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=1200&h=400&fit=crop",
			Target:      "/bookstore/rankings",
			TargetType:  "url",
			SortOrder:   3,
			IsActive:    true,
			ClickCount:  0,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Title:       "VIP 会员限时优惠",
			Description: "月卡仅需 9.9 元",
			Image:       "https://images.unsplash.com/photo-1559526324-4b87b5e36e44?w=1200&h=400&fit=crop",
			Target:      "/vip",
			TargetType:  "url",
			SortOrder:   4,
			IsActive:    true,
			ClickCount:  0,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	docs := make([]interface{}, len(banners))
	for i, banner := range banners {
		docs[i] = banner
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("插入 Banner 失败: %w", err)
	}

	fmt.Printf("  创建 %d 条 Banner 数据\n", len(result.InsertedIDs))
	return nil
}

// RankingItem 榜单项目结构
type RankingItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	BookID    primitive.ObjectID `bson:"book_id"`
	Type      string             `bson:"type"`
	Rank      int                `bson:"rank"`
	Score     float64            `bson:"score"`
	ViewCount int64              `bson:"view_count"`
	LikeCount int64              `bson:"like_count"`
	Period    string             `bson:"period"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func seedRankings(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("rankings")

	// 检查是否已有数据
	count, _ := collection.CountDocuments(ctx, bson.M{})
	if count > 0 {
		fmt.Printf("  已存在 %d 条榜单数据，跳过\n", count)
		return nil
	}

	now := time.Now()

	// 获取一本书的 ID 作为示例
	bookID, _ := primitive.ObjectIDFromHex("6956392cfe350a59abae6607")

	// 生成当前周期
	today := now.Format("2006-01-02")
	year, week := now.ISOWeek()
	weeklyPeriod := fmt.Sprintf("%d-W%02d", year, week)
	monthlyPeriod := now.Format("2006-01")

	rankings := []RankingItem{
		// 实时榜
		{
			BookID:    bookID,
			Type:      "realtime",
			Rank:      1,
			Score:     9.5,
			ViewCount: 100000,
			LikeCount: 5000,
			Period:    today,
			CreatedAt: now,
			UpdatedAt: now,
		},
		// 周榜
		{
			BookID:    bookID,
			Type:      "weekly",
			Rank:      1,
			Score:     9.6,
			ViewCount: 500000,
			LikeCount: 25000,
			Period:    weeklyPeriod,
			CreatedAt: now,
			UpdatedAt: now,
		},
		// 月榜
		{
			BookID:    bookID,
			Type:      "monthly",
			Rank:      1,
			Score:     9.7,
			ViewCount: 2000000,
			LikeCount: 100000,
			Period:    monthlyPeriod,
			CreatedAt: now,
			UpdatedAt: now,
		},
		// 新人榜
		{
			BookID:    bookID,
			Type:      "newbie",
			Rank:      1,
			Score:     9.4,
			ViewCount: 50000,
			LikeCount: 2500,
			Period:    monthlyPeriod,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	docs := make([]interface{}, len(rankings))
	for i, ranking := range rankings {
		docs[i] = ranking
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("插入榜单失败: %w", err)
	}

	fmt.Printf("  创建 %d 条榜单数据\n", len(result.InsertedIDs))
	return nil
}

func showStatistics(ctx context.Context, db *mongo.Database) {
	collections := []struct {
		Name  string
		Alias string
	}{
		{"categories", "分类"},
		{"banners", "Banner"},
		{"rankings", "榜单"},
		{"books", "书籍"},
	}

	total := 0
	for _, coll := range collections {
		count, err := db.Collection(coll.Name).CountDocuments(ctx, bson.M{})
		if err != nil {
			fmt.Printf("  ❌ 统计 %s 失败: %v\n", coll.Alias, err)
			continue
		}
		total += int(count)
		fmt.Printf("  %s: %d 条\n", coll.Alias, count)
	}
	fmt.Printf("  总计: %d 条数据\n", total)
}
