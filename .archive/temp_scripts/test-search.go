package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Printf("连接MongoDB失败: %v\n", err)
		return
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	booksCollection := db.Collection("books")

	fmt.Println("=== 模拟SearchWithFilter方法 ===")

	// 1. 构建基础查询（不包括关键词）
	query := bson.M{}
	fmt.Printf("初始查询条件: %+v\n", query)

	// 2. 如果有关键词，使用Go代码进行过滤
	keyword := "修仙"
	fmt.Printf("\n搜索关键词: %s\n", keyword)

	// 先获取符合其他条件的所有书籍（不设置limit以获取所有数据）
	opts := options.Find()
	sortBy := "created_at"
	sortOrder := -1
	opts.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})
	fmt.Printf("排序: %s %d\n", sortBy, sortOrder)

	cursor, err := booksCollection.Find(ctx, query, opts)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	// 定义Book结构体（简化版）
	type Book struct {
		ID           string   `bson:"_id"`
		Title        string   `bson:"title"`
		Author       string   `bson:"author"`
		Introduction string   `bson:"introduction"`
		Status       string   `bson:"status"`
		Categories   []string `bson:"categories"`
		Tags         []string `bson:"tags"`
	}

	var allBooks []*Book
	if err = cursor.All(ctx, &allBooks); err != nil {
		fmt.Printf("解析失败: %v\n", err)
		return
	}

	fmt.Printf("\n从数据库获取到 %d 本书\n", len(allBooks))

	// 在Go代码中进行关键词过滤
	var filteredBooks []*Book
	for _, book := range allBooks {
		// 直接使用 Contains 进行匹配
		titleMatch := strings.Contains(book.Title, keyword)
		authorMatch := strings.Contains(book.Author, keyword)
		introMatch := strings.Contains(book.Introduction, keyword)

		if titleMatch || authorMatch || introMatch {
			filteredBooks = append(filteredBooks, book)
			fmt.Printf("✓ 匹配: %s (title:%v, author:%v, intro:%v)\n",
				book.Title, titleMatch, authorMatch, introMatch)
		}
	}

	fmt.Printf("\n过滤后匹配 %d 本书\n", len(filteredBooks))

	// 应用分页限制
	limit := 20
	if len(filteredBooks) > limit {
		filteredBooks = filteredBooks[:limit]
	}
	fmt.Printf("应用分页限制(%d)后返回 %d 本书\n", limit, len(filteredBooks))
}
