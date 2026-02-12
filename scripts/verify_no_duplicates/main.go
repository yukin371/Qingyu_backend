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

const (
	collectionName = "reading_progress"
	databaseName   = "qingyu"
	mongoURI       = "mongodb://localhost:27017"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer client.Disconnect(ctx)

	coll := client.Database(databaseName).Collection(collectionName)

	// 验证无重复记录
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", bson.D{
				{"user_id", "$user_id"},
				{"book_id", "$book_id"},
			}},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$match", bson.D{
			{"count", bson.D{{"$gt", 1}}},
		}}},
		{{"$count", "duplicates"}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatalf("解析结果失败: %v", err)
	}

	if len(results) == 0 {
		fmt.Println("✅ 验证成功: 无重复记录")
	} else {
		fmt.Printf("❌ 验证失败: 发现 %d 组重复记录\n", results[0]["duplicates"])
	}

	// 获取当前集合文档总数
	totalDocs, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("获取文档总数失败: %v\n", err)
	} else {
		fmt.Printf("当前集合文档总数: %d\n", totalDocs)
	}
}
