package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
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

	// 查找所有节点
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	type OutlineNode struct {
		ID         primitive.ObjectID `bson:"_id"`
		DocumentID string             `bson:"document_id"`
		Status     string             `bson:"status"`
		Title      string             `bson:"title"`
		Tags       []string           `bson:"tags"`
	}

	var nodes []OutlineNode
	if err = cursor.All(ctx, &nodes); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("找到 %d 个节点需要检查\n", len(nodes))

	// 编译正则表达式匹配 chapter-binding:{id}
	chapterBindingRegex := regexp.MustCompile(`^chapter-binding:([a-f0-9]{24})$`)

	// 更新每个节点
	successCount := 0
	migratedDocIDCount := 0
	removedTagCount := 0
	setStatusCount := 0

	for _, node := range nodes {
		update := bson.M{}
		needsUpdate := false

		// 处理 chapter-binding tags
		newTags := make([]string, 0, len(node.Tags))
		for _, tag := range node.Tags {
			matches := chapterBindingRegex.FindStringSubmatch(tag)
			if matches != nil {
				// 找到 chapter-binding tag
				bindingID := matches[1]

				// 如果 DocumentID 为空，从 tag 中提取 ID
				if node.DocumentID == "" {
					update["document_id"] = bindingID
					node.DocumentID = bindingID
					migratedDocIDCount++
					needsUpdate = true
				}

				// 移除 chapter-binding tag
				removedTagCount++
				continue // 不添加到 newTags 中，相当于移除
			}

			// 保留非 chapter-binding 的 tag
			newTags = append(newTags, tag)
		}

		// 如果 tags 有变化，更新 tags 字段
		if len(newTags) != len(node.Tags) {
			if len(newTags) == 0 {
				// 如果没有其他 tags，设置为空数组而不是 nil
				update["tags"] = []string{}
			} else {
				update["tags"] = newTags
			}
			needsUpdate = true
		}

		// 处理 Status 字段
		if node.Status == "" {
			var newStatus string
			if node.DocumentID != "" {
				newStatus = "outlined"
			} else {
				newStatus = "draft"
			}
			update["status"] = newStatus
			setStatusCount++
			needsUpdate = true
		}

		// 执行更新
		if needsUpdate {
			filter := bson.M{"_id": node.ID}
			updateDoc := bson.M{"$set": update}

			result, err := collection.UpdateOne(ctx, filter, updateDoc)
			if err != nil {
				log.Printf("更新节点 %s (%s) 失败: %v\n", node.ID, node.Title, err)
				continue
			}

			if result.MatchedCount > 0 {
				successCount++
				fmt.Printf("✓ 更新节点: %s\n", node.Title)
				if docID, ok := update["document_id"]; ok {
					fmt.Printf("  - 迁移 DocumentID: %s\n", docID)
				}
				if status, ok := update["status"]; ok {
					fmt.Printf("  - 设置 Status: %s\n", status)
				}
				if _, ok := update["tags"]; ok {
					fmt.Printf("  - 移除 chapter-binding tags (移除 %d 个)\n", removedTagCount)
				}
			}
		}
	}

	fmt.Printf("\n迁移完成！成功更新 %d/%d 个节点\n", successCount, len(nodes))
	fmt.Printf("  - 从 tags 迁移 DocumentID: %d 个\n", migratedDocIDCount)
	fmt.Printf("  - 移除 chapter-binding tags: %d 个\n", removedTagCount)
	fmt.Printf("  - 设置 Status 字段: %d 个\n", setStatusCount)

	// 验证结果
	withDocID, _ := collection.CountDocuments(ctx, bson.M{
		"document_id": bson.M{"$exists": true, "$ne": ""},
	})
	fmt.Printf("\n现在有 DocumentID 的节点数量: %d\n", withDocID)

	withStatus, _ := collection.CountDocuments(ctx, bson.M{
		"status": bson.M{"$exists": true, "$ne": ""},
	})
	fmt.Printf("现在有 Status 的节点数量: %d\n", withStatus)

	withChapterBinding, _ := collection.CountDocuments(ctx, bson.M{
		"tags": "chapter-binding",
	})
	if withChapterBinding > 0 {
		fmt.Printf("警告：仍有 %d 个节点包含 chapter-binding tags\n", withChapterBinding)
	} else {
		fmt.Printf("✓ 所有节点都已移除 chapter-binding tags\n")
	}
}
