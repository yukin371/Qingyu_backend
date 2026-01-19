package seeds

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Book 书籍模型（简化版）
type Book struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Title         string             `bson:"title"`
	Author        string             `bson:"author"`
	Introduction  string             `bson:"introduction"`
	Cover         string             `bson:"cover"`
	Categories    []string           `bson:"categories"`
	Tags          []string           `bson:"tags"`
	Status        string             `bson:"status"`
	WordCount     int64              `bson:"word_count"`
	ChapterCount  int                `bson:"chapter_count"`
	Price         float64            `bson:"price"`
	IsFree        bool               `bson:"is_free"`
	IsRecommended bool               `bson:"is_recommended"`
	IsFeatured    bool               `bson:"is_featured"`
	IsHot         bool               `bson:"is_hot"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

// SeedBooks 书籍种子数据
func SeedBooks(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("books")

	// 检查是否已有数据
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to count books: %w", err)
	}

	if count > 0 {
		fmt.Printf("Books collection already has %d documents, skipping seed\n", count)
		return nil
	}

	now := time.Now()

	// 准备书籍数据
	books := []Book{
		{
			Title:         "修真世界",
			Author:        "方想",
			Introduction:  "一个卖力气的普通小子，在修真界摸爬滚打，不断成长的故事。",
			Cover:         "/covers/xiuzhen.jpg",
			Categories:    []string{"仙侠", "修真"},
			Tags:          []string{"凡人流", "逆袭", "修真"},
			Status:        "completed",
			WordCount:     6000000,
			ChapterCount:  1320,
			Price:         0.05,
			IsFree:        false,
			IsRecommended: true,
			IsFeatured:    true,
			IsHot:         true,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			Title:         "诡秘之主",
			Author:        "爱潜水的乌贼",
			Introduction:  "蒸汽与机械的浪潮中，谁能触及非凡？历史和黑暗的迷雾里，又是谁在耳语？",
			Cover:         "/covers/guimi.jpg",
			Categories:    []string{"奇幻", "克苏鲁"},
			Tags:          []string{"蒸汽朋克", "克苏鲁", "神秘"},
			Status:        "completed",
			WordCount:     5000000,
			ChapterCount:  1394,
			Price:         0.05,
			IsFree:        false,
			IsRecommended: true,
			IsFeatured:    true,
			IsHot:         true,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			Title:         "全职高手",
			Author:        "蝴蝶蓝",
			Introduction:  "荣耀中被誉为教科书级别的顶尖高手，因为种种原因遭到俱乐部的驱逐。",
			Cover:         "/covers/quanzhi.jpg",
			Categories:    []string{"游戏", "电竞"},
			Tags:          []string{"电竞", "热血", "团队"},
			Status:        "completed",
			WordCount:     5300000,
			ChapterCount:  1728,
			Price:         0.05,
			IsFree:        false,
			IsRecommended: true,
			IsFeatured:    false,
			IsHot:         true,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			Title:         "斗破苍穹",
			Author:        "天蚕土豆",
			Introduction:  "这里是属于斗气的世界，没有花俏艳丽的魔法，有的，仅仅是繁衍到巅峰的斗气！",
			Cover:         "/covers/doupo.jpg",
			Categories:    []string{"玄幻", "异世大陆"},
			Tags:          []string{"废材流", "逆袭", "热血"},
			Status:        "completed",
			WordCount:     5300000,
			ChapterCount:  1648,
			Price:         0.05,
			IsFree:        true,
			IsRecommended: false,
			IsFeatured:    false,
			IsHot:         true,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			Title:         "大奉打更人",
			Author:        "卖报小郎君",
			Introduction:  "这个世界，有儒；有道；有佛；有妖；有术士。",
			Cover:         "/covers/dafeng.jpg",
			Categories:    []string{"玄幻", "东方玄幻"},
			Tags:          []string{"探案", "推理", "玄幻"},
			Status:        "ongoing",
			WordCount:     4200000,
			ChapterCount:  1100,
			Price:         0.05,
			IsFree:        false,
			IsRecommended: true,
			IsFeatured:    true,
			IsHot:         true,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	// 插入书籍
	docs := make([]interface{}, len(books))
	for i, book := range books {
		docs[i] = book
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to insert books: %w", err)
	}

	fmt.Printf("✓ Seeded %d books\n", len(result.InsertedIDs))
	for i, id := range result.InsertedIDs {
		fmt.Printf("  - %s by %s (ID: %v)\n", books[i].Title, books[i].Author, id)
	}

	return nil
}
