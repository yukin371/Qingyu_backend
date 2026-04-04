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
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	dbName := "qingyu"
	// 列出所有collections
	collections, err := client.Database(dbName).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Collections:", collections)

	// 检查outlines集合
	collection := client.Database(dbName).Collection("outlines")

	// 统计文档数量
	total, _ := collection.CountDocuments(ctx, bson.M{})
	fmt.Printf("\noutline_nodes总文档数: %d\n", total)

	// 查询一个示例文档
	var result bson.M
	collection.FindOne(ctx, bson.M{}).Decode(&result)
	fmt.Printf("\n示例文档字段: %v\n", result)

	// 统计有document_id的节点
	withDocId, _ := collection.CountDocuments(ctx, bson.M{
		"document_id": bson.M{"$exists": true, "$ne": ""},
	})
	fmt.Printf("\n有document_id的节点: %d\n", withDocId)

	// 统计有tags的节点
	withTags, _ := collection.CountDocuments(ctx, bson.M{
		"tags": bson.M{"$exists": true, "$ne": []interface{}{}},
	})
	fmt.Printf("有tags的节点: %d\n", withTags)

	// 查询一个有document_id的文档
	var sample bson.M
	collection.FindOne(ctx, bson.M{
		"document_id": bson.M{"$exists": true, "$ne": ""},
	}, options.FindOne().SetProjection(bson.M{
		"title":       1,
		"document_id": 1,
		"tags":        1,
	})).Decode(&sample)

	fmt.Printf("\n示例有document_id的节点:\n")
	fmt.Printf("  Title: %v\n", sample["title"])
	fmt.Printf("  DocumentID: %v\n", sample["document_id"])
	fmt.Printf("  Tags: %v\n", sample["tags"])
}
