package seeds

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NovelCleaner 小说数据清理器
type NovelCleaner struct {
	db                *mongo.Database
	bookCollection    *mongo.Collection
	chapterCollection *mongo.Collection
}

// NewNovelCleaner 创建清理器
func NewNovelCleaner(db *mongo.Database) *NovelCleaner {
	return &NovelCleaner{
		db:                db,
		bookCollection:    db.Collection("books"),
		chapterCollection: db.Collection("chapters"),
	}
}

// Clean 清理所有导入的小说数据
func (nc *NovelCleaner) Clean(ctx context.Context) error {
	log.Println("开始清理小说数据...")

	// 获取清理前的统计
	bookCount, _ := nc.bookCollection.CountDocuments(ctx, primitive.M{})
	chapterCount, _ := nc.chapterCollection.CountDocuments(ctx, primitive.M{})

	log.Printf("清理前统计:")
	log.Printf("  书籍数量: %d", bookCount)
	log.Printf("  章节数量: %d", chapterCount)

	// 清理章节
	chapterResult, err := nc.chapterCollection.DeleteMany(ctx, primitive.M{})
	if err != nil {
		return fmt.Errorf("清理章节失败: %w", err)
	}
	log.Printf("✓ 已删除 %d 个章节", chapterResult.DeletedCount)

	// 清理书籍
	bookResult, err := nc.bookCollection.DeleteMany(ctx, primitive.M{})
	if err != nil {
		return fmt.Errorf("清理书籍失败: %w", err)
	}
	log.Printf("✓ 已删除 %d 本书籍", bookResult.DeletedCount)

	log.Println("\n✓ 清理完成")
	return nil
}

// CleanByCategory 按分类清理
func (nc *NovelCleaner) CleanByCategory(ctx context.Context, category string) error {
	log.Printf("开始清理分类 [%s] 的小说数据...", category)

	// 查找该分类的书籍ID
	filter := primitive.M{"categories": category}
	cursor, err := nc.bookCollection.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("查询书籍失败: %w", err)
	}
	defer cursor.Close(ctx)

	var bookIDs []interface{}
	for cursor.Next(ctx) {
		var book struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&book); err != nil {
			continue
		}
		bookIDs = append(bookIDs, book.ID)
	}

	if len(bookIDs) == 0 {
		log.Printf("未找到分类 [%s] 的书籍", category)
		return nil
	}

	// 删除章节
	chapterFilter := primitive.M{"book_id": primitive.M{"$in": bookIDs}}
	chapterResult, err := nc.chapterCollection.DeleteMany(ctx, chapterFilter)
	if err != nil {
		return fmt.Errorf("删除章节失败: %w", err)
	}

	// 删除书籍
	bookResult, err := nc.bookCollection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除书籍失败: %w", err)
	}

	log.Printf("✓ 已删除 %d 本书籍, %d 个章节", bookResult.DeletedCount, chapterResult.DeletedCount)
	return nil
}

// GetStats 获取当前统计
func (nc *NovelCleaner) GetStats(ctx context.Context) error {
	bookCount, err := nc.bookCollection.CountDocuments(ctx, primitive.M{})
	if err != nil {
		return err
	}

	chapterCount, err := nc.chapterCollection.CountDocuments(ctx, primitive.M{})
	if err != nil {
		return err
	}

	log.Printf("当前数据库统计:")
	log.Printf("  书籍总数: %d", bookCount)
	log.Printf("  章节总数: %d", chapterCount)

	return nil
}
