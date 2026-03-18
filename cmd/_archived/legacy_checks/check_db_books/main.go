package main

import (
	"context"
	"fmt"
	"log"
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
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Disconnect(ctx)

	// 检查各个数据库的书籍数量
	databases := []string{"qingyu", "Qingyu_writer", "Qingyu_backend"}

	for _, dbName := range databases {
		db := client.Database(dbName)
		count, err := db.Collection("books").CountDocuments(ctx, bson.M{})
		if err != nil {
			fmt.Printf("%s.books: 错误 - %v\n", dbName, err)
		} else {
			fmt.Printf("%s.books: %d 本\n", dbName, count)

			// 如果有数据，查看前2本的标题
			if count > 0 {
				opts := options.Find().SetLimit(2).SetProjection(bson.M{"title": 1, "author": 1})
				cursor, err := db.Collection("books").Find(ctx, bson.M{}, opts)
				if err == nil {
					var results []bson.M
					cursor.All(ctx, &results)
					for _, book := range results {
						fmt.Printf("  - %s by %s\n", book["title"], book["author"])
					}
				}
			}
		}
	}
}
