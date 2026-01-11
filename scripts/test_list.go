//go:build ignore

package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	coll := db.Collection("books")

	// 查询总数
	total, _ := coll.CountDocuments(ctx, bson.M{})
	fmt.Printf("Total books in DB: %d\n", total)

	// 查询前10本
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}, {Key: "title", Value: 1}})
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		fmt.Printf("Find error: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	// 使用bson.M解码避免类型问题
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		fmt.Printf("All error: %v\n", err)
		return
	}

	fmt.Printf("Found %d books:\n", len(results))
	for i, b := range results {
		if i >= 10 {
			break
		}
		title := b["title"]
		author := b["author"]
		fmt.Printf("%d. %v - %v\n", i+1, title, author)
	}
}
