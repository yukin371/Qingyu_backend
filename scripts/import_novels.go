package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NovelData JSON数据结构
type NovelData struct {
	Metadata Metadata         `json:"metadata"`
	Novels   []Novel          `json:"novels"`
}

type Metadata struct {
	Source          string    `json:"source"`
	TotalNovels     int       `json:"total_novels"`
	TotalChapters   int       `json:"total_chapters"`
	GeneratedAt     time.Time `json:"generated_at"`
	ChapterSize     int       `json:"chapter_size"`
}

type Novel struct {
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Introduction  string    `json:"introduction"`
	Category      string    `json:"category"`
	WordCount     int       `json:"word_count"`
	ChapterCount  int       `json:"chapter_count"`
	Rating        float64   `json:"rating"`
	Status        string    `json:"status"`
	IsFree        bool      `json:"is_free"`
	Chapters      []Chapter `json:"chapters"`
}

type Chapter struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	WordCount int    `json:"word_count"`
}

// Book 书籍模型
type Book struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	Title         string               `bson:"title"`
	Author        string               `bson:"author"`
	AuthorID      primitive.ObjectID   `bson:"author_id,omitempty"`
	Introduction  string               `bson:"introduction"`
	Cover         string               `bson:"cover"`
	CategoryIDs   []primitive.ObjectID `bson:"category_ids"`
	Categories    []string             `bson:"categories"`
	Tags          []string             `bson:"tags"`
	Status        string               `bson:"status"`
	Rating        float64              `bson:"rating"`
	RatingCount   int64                `bson:"rating_count"`
	ViewCount     int64                `bson:"view_count"`
	WordCount     int64                `bson:"word_count"`
	ChapterCount  int                  `bson:"chapter_count"`
	Price         float64              `bson:"price"`
	IsFree        bool                 `bson:"is_free"`
	IsRecommended bool                 `bson:"is_recommended"`
	IsFeatured    bool                 `bson:"is_featured"`
	IsHot         bool                 `bson:"is_hot"`
	PublishedAt   *time.Time           `bson:"published_at,omitempty"`
	LastUpdateAt  *time.Time           `bson:"last_update_at,omitempty"`
	CreatedAt     time.Time            `bson:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at"`
}

// ChapterContent 章节内容模型
type ChapterContent struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	BookID         primitive.ObjectID `bson:"book_id"`
	ChapterNum     int                `bson:"chapter_num"`
	Title          string             `bson:"title"`
	Content        string             `bson:"content"`
	WordCount      int                `bson:"word_count"`
	IsFree         bool               `bson:"is_free"`
	Price          float64            `bson:"price"`
	PublishedAt    time.Time          `bson:"published_at"`
	CreatedAt      time.Time          `bson:"created_at"`
}

func main() {
	// 读取JSON文件
	filePath := "data/novels_100.json"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	var novelData NovelData
	if err := json.Unmarshal(data, &novelData); err != nil {
		log.Fatalf("解析JSON失败: %v", err)
	}

	fmt.Printf("📚 从文件加载了 %d 本小说\n", len(novelData.Novels))

	// 连接MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	booksCollection := db.Collection("books")
	chaptersCollection := db.Collection("chapter_contents")

	// 默认作者ID（用于所有导入的小说）
	defaultAuthorID := primitive.NewObjectID()

	fmt.Println("🔄 开始导入小说数据...")

	importedCount := 0
	totalChapters := 0

	for i, novel := range novelData.Novels {
		// 创建书籍ID
		bookID := primitive.NewObjectID()

		// 映射状态
		status := novel.Status
		if status == "" {
			status = "published"
		}

		// 创建Book对象
		book := Book{
			ID:            bookID,
			Title:         novel.Title,
			Author:        novel.Author,
			AuthorID:      defaultAuthorID,
			Introduction:  novel.Introduction,
			Cover:         "/covers/default.jpg",
			Categories:    []string{novel.Category},
			Tags:          []string{},
			Status:        status,
			Rating:        novel.Rating,
			RatingCount:   0,
			ViewCount:     0,
			WordCount:     int64(novel.WordCount),
			ChapterCount:  len(novel.Chapters),
			Price:         0.0,
			IsFree:        novel.IsFree,
			IsRecommended: false,
			IsFeatured:    false,
			IsHot:         false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// 插入书籍
		_, err := booksCollection.InsertOne(ctx, book)
		if err != nil {
			log.Printf("❌ 插入书籍失败 [%d] %s: %v", i+1, novel.Title, err)
			continue
		}

		// 批量插入章节
		if len(novel.Chapters) > 0 {
			var chapters []interface{}
			for j, ch := range novel.Chapters {
				chapter := ChapterContent{
					BookID:      bookID,
					ChapterNum:  j + 1,
					Title:       ch.Title,
					Content:     ch.Content,
					WordCount:   ch.WordCount,
					IsFree:      novel.IsFree,
					Price:       0.0,
					PublishedAt: time.Now(),
					CreatedAt:   time.Now(),
				}
				chapters = append(chapters, chapter)
			}

			_, err := chaptersCollection.InsertMany(ctx, chapters)
			if err != nil {
				log.Printf("⚠️  插入章节失败 [%d] %s: %v", i+1, novel.Title, err)
			} else {
				totalChapters += len(chapters)
			}
		}

		importedCount++
		if importedCount%10 == 0 {
			fmt.Printf("  已导入 %d/%d 本小说...\n", importedCount, len(novelData.Novels))
		}
	}

	fmt.Printf("\n✅ 导入完成！\n")
	fmt.Printf("   📖 成功导入: %d 本小说\n", importedCount)
	fmt.Printf("   📝 总章节数: %d 章\n", totalChapters)

	// 验证导入
	count, _ := booksCollection.CountDocuments(ctx, bson.M{})
	chapterCount, _ := chaptersCollection.CountDocuments(ctx, bson.M{})
	fmt.Printf("   📊 数据库统计:\n")
	fmt.Printf("      - 书籍总数: %d\n", count)
	fmt.Printf("      - 章节总数: %d\n", chapterCount)
}
