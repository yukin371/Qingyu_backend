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

// BookFilter 模拟BookFilter结构
type BookFilter struct {
	Keyword *string
	Status  *string
	Author  *string
	Limit   int
	Offset  int
	SortBy  string
}

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

	fmt.Println("=== 端到端搜索测试 ===\n")

	// 1. 模拟API层构建filter
	keyword := "修仙"
	filter := &BookFilter{
		Keyword: &keyword,
		Limit:   20,
		Offset:  0,
		SortBy:  "created_at",
	}
	fmt.Printf("1. API层构建filter:\n")
	fmt.Printf("   Keyword: %s\n", *filter.Keyword)
	fmt.Printf("   Status: %v\n", filter.Status)
	fmt.Printf("   Limit: %d\n", filter.Limit)
	fmt.Printf("   Offset: %d\n", filter.Offset)

	// 2. Repository层构建query
	query := bson.M{}
	if filter.Status != nil {
		query["status"] = *filter.Status
		fmt.Printf("\n2. Repository层添加status过滤: %s\n", *filter.Status)
	} else {
		fmt.Printf("\n2. Repository层: 无status过滤\n")
	}
	fmt.Printf("   MongoDB query: %+v\n", query)

	// 3. 执行查询（获取所有书籍）
	opts := options.Find()
	sortOrder := -1
	opts.SetSort(bson.D{{Key: filter.SortBy, Value: sortOrder}})
	fmt.Printf("\n3. 执行MongoDB查询，排序: %s %d\n", filter.SortBy, sortOrder)

	cursor, err := booksCollection.Find(ctx, query, opts)
	if err != nil {
		fmt.Printf("   查询失败: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	// 定义Book结构
	type Book struct {
		ID           string   `bson:"_id"`
		Title        string   `bson:"title"`
		Author       string   `bson:"author"`
		Introduction string   `bson:"introduction"`
		Status       string   `bson:"status"`
	}

	var allBooks []*Book
	if err = cursor.All(ctx, &allBooks); err != nil {
		fmt.Printf("   解析失败: %v\n", err)
		return
	}
	fmt.Printf("   获取到 %d 本书\n", len(allBooks))

	// 4. 在Go代码中进行关键词过滤
	fmt.Printf("\n4. Go代码过滤关键词: '%s'\n", *filter.Keyword)
	var filteredBooks []*Book
	for i, book := range allBooks {
		titleMatch := strings.Contains(book.Title, *filter.Keyword)
		authorMatch := strings.Contains(book.Author, *filter.Keyword)
		introMatch := strings.Contains(book.Introduction, *filter.Keyword)

		if titleMatch || authorMatch || introMatch {
			filteredBooks = append(filteredBooks, book)
			fmt.Printf("   ✓ [%d] %s (匹配字段: title=%v, author=%v, intro=%v)\n",
				i+1, book.Title, titleMatch, authorMatch, introMatch)
		}
	}

	fmt.Printf("\n   过滤后: %d 本书\n", len(filteredBooks))

	// 5. 应用分页限制
	if filter.Limit > 0 && len(filteredBooks) > filter.Limit {
		filteredBooks = filteredBooks[:filter.Limit]
	}
	fmt.Printf("\n5. 应用分页限制(%d)后: %d 本书\n", filter.Limit, len(filteredBooks))

	// 6. 最终结果
	fmt.Printf("\n=== 最终结果 ===\n")
	fmt.Printf("返回 %d 本书:\n", len(filteredBooks))
	for i, book := range filteredBooks {
		fmt.Printf("%d. %s - %s (status: %s)\n", i+1, book.Title, book.Author, book.Status)
	}
}
