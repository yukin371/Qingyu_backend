package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	dbName := "qingyu"
	collection := client.Database(dbName).Collection("outlines")

	// 查找所有有document_id但没有tags的节点
	cursor, err := collection.Find(ctx, bson.M{
		"document_id": bson.M{"$exists": true, "$ne": ""},
		"tags":        bson.M{"$exists": false},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	type OutlineNode struct {
		ID         primitive.ObjectID `bson:"_id"`
		DocumentID string             `bson:"document_id"`
		Title      string             `bson:"title"`
	}

	var nodes []OutlineNode
	if err = cursor.All(ctx, &nodes); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("找到 %d 个需要迁移的节点\n", len(nodes))

	// 更新每个节点
	successCount := 0
	for _, node := range nodes {
		tag := fmt.Sprintf("chapter-binding:%s", node.DocumentID)
		filter := bson.M{"_id": node.ID}
		update := bson.M{"$set": bson.M{"tags": []string{tag}}}

		result, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Printf("更新节点 %s (%s) 失败: %v\n", node.ID, node.Title, err)
			continue
		}

		if result.MatchedCount > 0 {
			successCount++
			fmt.Printf("✓ 更新节点: %s -> tags: [%s]\n", node.Title, tag)
		}
	}

	fmt.Printf("\n迁移完成！成功更新 %d/%d 个节点\n", successCount, len(nodes))

	// 验证结果
	withTags, _ := collection.CountDocuments(ctx, bson.M{
		"tags": bson.M{"$exists": true, "$ne": []interface{}{}},
	})
	fmt.Printf("现在有tags的节点数量: %d\n", withTags)
}
