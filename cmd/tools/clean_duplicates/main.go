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

	// 查询所有 arc 类型的大纲
	cursor, err := db.Collection("outlines").Find(ctx, bson.M{
		"project_id": projectOID,
		"type":       "arc",
	}, options.Find().SetSort(bson.M{"order": 1}))
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	defer cursor.Close(ctx)

	var outlines []bson.M
	if err = cursor.All(ctx, &outlines); err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	// 按 document_id 分组
	docToOutlines := make(map[string][]bson.M)
	for _, outline := range outlines {
		docID := ""
		if docIDField, ok := outline["document_id"]; ok && docIDField != nil {
			docID = docIDField.(string)
		}
		if docID != "" {
			docToOutlines[docID] = append(docToOutlines[docID], outline)
		}
	}

	// 找出并清理重复的
	fmt.Printf("=== 清理重复大纲节点 ===\n\n")
	totalDeleted := 0

	for docID, outlineList := range docToOutlines {
		if len(outlineList) > 1 {
			fmt.Printf("文档 %s 有 %d 个关联大纲:\n", docID, len(outlineList))
			// 保留最早创建的，删除其他的
			var earliest bson.M
			var earliestTime time.Time
			first := true
			for _, outline := range outlineList {
				id := outline["_id"].(primitive.ObjectID)
				createdAt := time.Time{}
				if ct, ok := outline["created_at"]; ok && ct != nil {
					switch v := ct.(type) {
					case time.Time:
						createdAt = v
					case primitive.DateTime:
						createdAt = v.Time()
					}
				}
				fmt.Printf("  - %s (created_at: %v)\n", id.Hex(), createdAt)

				if first || createdAt.Before(earliestTime) {
					earliest = outline
					earliestTime = createdAt
					first = false
				}
			}

			// 删除除了最早创建之外的所有大纲
			fmt.Printf("  保留: %s\n", earliest["_id"].(primitive.ObjectID).Hex())
			for _, outline := range outlineList {
				id := outline["_id"].(primitive.ObjectID)
				if id.Hex() != earliest["_id"].(primitive.ObjectID).Hex() {
					fmt.Printf("  删除: %s\n", id.Hex())
					_, err := db.Collection("outlines").DeleteOne(ctx, bson.M{"_id": id})
					if err != nil {
						log.Printf("  删除失败: %v", err)
					} else {
						totalDeleted++
					}
				}
			}
			fmt.Println()
		}
	}

	fmt.Printf("清理完成: 共删除 %d 个重复大纲节点\n", totalDeleted)

	// 验证
	count, _ := db.Collection("outlines").CountDocuments(ctx, bson.M{
		"project_id": projectOID,
		"type":       "arc",
	})
	fmt.Printf("剩余 arc 大纲数量: %d\n", count)
}
