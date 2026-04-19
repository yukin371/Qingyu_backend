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

// Post 动态结构（与 models/social/post.go 一致）
type Post struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       string           `bson:"user_id" json:"userId"`
	UserName     string           `bson:"user_name" json:"userName"`
	UserAvatar   string           `bson:"user_avatar,omitempty" json:"userAvatar,omitempty"`
	UserLevel    int              `bson:"user_level" json:"userLevel"`
	Type         string           `bson:"type" json:"type"`
	Content      string           `bson:"content" json:"content"`
	Images       []string         `bson:"images,omitempty" json:"images,omitempty"`
	BookID       string           `bson:"book_id,omitempty" json:"bookId,omitempty"`
	BookTitle    string           `bson:"book_title,omitempty" json:"bookTitle,omitempty"`
	BookCover    string           `bson:"book_cover,omitempty" json:"bookCover,omitempty"`
	BookAuthor   string           `bson:"book_author,omitempty" json:"bookAuthor,omitempty"`
	ChapterID    string           `bson:"chapter_id,omitempty" json:"chapterId,omitempty"`
	ChapterTitle string           `bson:"chapter_title,omitempty" json:"chapterTitle,omitempty"`
	Progress     int              `bson:"progress,omitempty" json:"progress,omitempty"`
	Topics       []string         `bson:"topics,omitempty" json:"topics,omitempty"`
	LikeCount    int              `bson:"like_count" json:"likeCount"`
	CommentCount int              `bson:"comment_count" json:"commentCount"`
	ShareCount   int              `bson:"share_count" json:"shareCount"`
	CreatedAt    time.Time       `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time       `bson:"updated_at" json:"updatedAt"`
}

func main() {
	fmt.Println("开始插入社区动态数据到 MongoDB...")

	// 读取测试数据
	data, err := os.ReadFile("./tmp/community_posts_seed.json")
	if err != nil {
		log.Fatalf("读取测试数据文件失败: %v", err)
	}

	var rawPosts []map[string]interface{}
	if err := json.Unmarshal(data, &rawPosts); err != nil {
		log.Fatalf("解析测试数据失败: %v", err)
	}

	// 连接 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("连接 MongoDB 失败: %v", err)
	}
	defer client.Disconnect(ctx)

	// Ping 检查
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB Ping 失败: %v", err)
	}
	fmt.Println("✅ MongoDB 连接成功")

	db := client.Database("qingyu_test")
	collection := db.Collection("posts")

	// 转换数据
	var posts []interface{}
	for _, raw := range rawPosts {
		post := Post{
			ID:            primitive.NewObjectID(),
			UserID:       getString(raw, "userId"),
			UserName:     getString(raw, "userName"),
			UserAvatar:   getString(raw, "userAvatar"),
			UserLevel:    getInt(raw, "userLevel"),
			Type:         getString(raw, "type"),
			Content:      getString(raw, "content"),
			Images:       getStringArray(raw, "images"),
			BookID:       getString(raw, "bookId"),
			BookTitle:    getString(raw, "bookTitle"),
			BookCover:    getString(raw, "bookCover"),
			BookAuthor:   getString(raw, "bookAuthor"),
			ChapterID:    getString(raw, "chapterId"),
			ChapterTitle: getString(raw, "chapterTitle"),
			Progress:     getInt(raw, "progress"),
			Topics:       getStringArray(raw, "topics"),
			LikeCount:    getInt(raw, "likeCount"),
			CommentCount: getInt(raw, "commentCount"),
			ShareCount:   getInt(raw, "shareCount"),
			CreatedAt:    mustParseTime(getString(raw, "createdAt")),
			UpdatedAt:    mustParseTime(getString(raw, "updatedAt")),
		}
		posts = append(posts, post)
	}

	// 插入数据
	if len(posts) > 0 {
		result, err := collection.InsertMany(ctx, posts)
		if err != nil {
			log.Fatalf("插入数据失败: %v", err)
		}
		fmt.Printf("✅ 成功插入 %d 条动态数据！\n", len(result.InsertedIDs))
	} else {
		fmt.Println("⚠️ 没有数据需要插入")
	}

	// 创建索引
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "topics", Value: 1}}},
		{Keys: bson.D{{Key: "like_count", Value: -1}}},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Printf("⚠️ 创建索引失败: %v", err)
	} else {
		fmt.Println("✅ 索引创建成功")
	}

	fmt.Println("\n🎉 社区动态数据填充完成！")
}

// Helper functions
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return 0
}

func getStringArray(m map[string]interface{}, key string) []string {
	if v, ok := m[key].([]interface{}); ok {
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Now()
	}
	return t
}
