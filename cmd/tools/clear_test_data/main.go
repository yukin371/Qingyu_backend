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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?directConnection=true"))
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	// 清空相关集合
	collections := []string{"outlines", "documents", "projects", "chapters", "characters", "locations"}
	totalDeleted := int64(0)

	for _, collName := range collections {
		result, err := db.Collection(collName).DeleteMany(ctx, bson.M{})
		if err != nil {
			log.Printf("清空 %s 失败: %v", collName, err)
			continue
		}
		totalDeleted += result.DeletedCount
		fmt.Printf("✓ 清空 %s: 删除 %d 条记录\n", collName, result.DeletedCount)
	}

	fmt.Printf("\n总计删除 %d 条记录，数据库已清空\n", totalDeleted)
}
