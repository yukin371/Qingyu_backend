package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	coll := db.Collection("books")

	// 查询所有已发布的书籍
	cursor, err := coll.Find(ctx, bson.M{
		"status": bson.M{"$in": []string{"published", "ongoing", "completed"}},
	}, nil)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	defer cursor.Close(ctx)

	var allBooks []bson.M
	if err = cursor.All(ctx, &allBooks); err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	// 在 Go 代码中进行过滤
	keyword := "爱"
	keywordLower := strings.ToLower(keyword)
	var filteredBooks []bson.M
	for _, book := range allBooks {
		title := fmt.Sprintf("%v", book["title"])
		author := fmt.Sprintf("%v", book["author"])
		intro := fmt.Sprintf("%v", book["introduction"])
		
		if strings.Contains(strings.ToLower(title), keywordLower) ||
			strings.Contains(strings.ToLower(author), keywordLower) ||
			strings.Contains(strings.ToLower(intro), keywordLower) {
			filteredBooks = append(filteredBooks, book)
		}
	}

	fmt.Printf("搜索关键词: %s\n", keyword)
	fmt.Printf("匹配书籍数: %d\n\n", len(filteredBooks))
	
	if len(filteredBooks) > 0 {
		fmt.Println("匹配的书籍:")
		for i, book := range filteredBooks {
			fmt.Printf("%d. %v\n", i+1, book["title"])
		}
	}
}
