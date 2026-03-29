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

	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?directConnection=true"))
	db := client.Database("qingyu")

	projectID := "69c790d7be2b6d81b914bfed"
	projectOID, _ := primitive.ObjectIDFromHex(projectID)

	// 查询没有关联文档的大纲节点
	cursor, err := db.Collection("outlines").Find(ctx, bson.M{
		"project_id": projectOID,
		"$or": []bson.M{
			{"document_id": bson.M{"$exists": false}},
			{"document_id": ""},
		},
	})
	if err != nil {
		log.Printf("查询失败: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var outlines []bson.M
	if err = cursor.All(ctx, &outlines); err != nil {
		log.Printf("解析大纲失败: %v", err)
		return
	}

	fmt.Printf("找到 %d 个没有关联文档的大纲节点\n", len(outlines))

	// 查询没有关联大纲的文档
	docCursor, err := db.Collection("documents").Find(ctx, bson.M{
		"project_id": projectOID,
		"$or": []bson.M{
			{"outline_node_id": bson.M{"$exists": false}},
			{"outline_node_id": ""},
		},
	})
	if err != nil {
		log.Printf("查询文档失败: %v", err)
		return
	}
	defer docCursor.Close(ctx)

	var documents []bson.M
	if err = docCursor.All(ctx, &documents); err != nil {
		log.Printf("解析文档失败: %v", err)
		return
	}

	fmt.Printf("找到 %d 个没有关联大纲的文档\n", len(documents))

	// 统计匹配情况
	type MatchStatus struct {
		OutlineTitle      string
		DocumentTitle     string
		OutlineID         string
		DocumentID        string
		OutlineHasDoc     bool
		DocumentHasOutline bool
	}

	// 检查双向绑定情况
	fmt.Println("\n=== 双向绑定检查 ===")

	// 检查有 outline_node_id 的文档
	matchedDocs, _ := db.Collection("documents").CountDocuments(ctx, bson.M{
		"project_id":      projectOID,
		"outline_node_id": bson.M{"$ne": ""},
	})
	fmt.Printf("文档有 outline_node_id: %d\n", matchedDocs)

	// 检查有 document_id 的大纲
	matchedOutlines, _ := db.Collection("outlines").CountDocuments(ctx, bson.M{
		"project_id": projectOID,
		"document_id": bson.M{"$ne": ""},
	})
	fmt.Printf("大纲有 document_id: %d\n", matchedOutlines)

	// 检查 volumes 和 chapters
	volumes, _ := db.Collection("documents").CountDocuments(ctx, bson.M{
		"project_id": projectOID,
		"type":       "volume",
	})
	fmt.Printf("卷文档数量: %d\n", volumes)

	chapters, _ := db.Collection("documents").CountDocuments(ctx, bson.M{
		"project_id": projectOID,
		"type":       "chapter",
	})
	fmt.Printf("章节文档数量: %d\n", chapters)

	arcOutlines, _ := db.Collection("outlines").CountDocuments(ctx, bson.M{
		"project_id": projectOID,
		"type":       "arc",
	})
	fmt.Printf("arc 大纲数量: %d\n", arcOutlines)

	sceneOutlines, _ := db.Collection("outlines").CountDocuments(ctx, bson.M{
		"project_id": projectOID,
		"type":       "scene",
	})
	fmt.Printf("scene 大纲数量: %d\n", sceneOutlines)

	fmt.Println("\n建议:")
	if len(outlines) > 0 {
		fmt.Printf("- 需要为 %d 个大纲创建对应文档\n", len(outlines))
	}
	if len(documents) > 0 {
		fmt.Printf("- 需要为 %d 个文档创建对应大纲\n", len(documents))
	}
}
