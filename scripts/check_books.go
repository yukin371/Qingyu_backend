package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	Title     string `bson:"title"`
	Status    string `bson:"status"`
	CreatedAt string `bson:"created_at"`
}

func main() {
	ctx := context.Background()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	books := db.Collection("books")

	// 统计各状态数量
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$status", "count": bson.M{"$sum": 1}}},
	}
	cursor, _ := books.Aggregate(ctx, pipeline)
	var results []bson.M
	cursor.All(ctx, &results)

	fmt.Println("📊 书籍状态分布:")
	for _, r := range results {
		fmt.Printf("   %s: %d 本\n", r["_id"], r["count"])
	}

	// 检查前5本书
	cursor2, _ := books.Find(ctx, bson.M{}, options.Find().SetLimit(5).SetSort(bson.D{{Key: "created_at", Value: -1}}))
	var booksData []Book
	cursor2.All(ctx, &booksData)

	fmt.Println("\n📖 最新5本书:")
	for _, b := range booksData {
		fmt.Printf("   - %s (状态: %s)\n", b.Title, b.Status)
	}
}
