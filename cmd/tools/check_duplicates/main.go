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

	fmt.Printf("=== Arc 大纲节点列表 ===\n")
	fmt.Printf("总数: %d\n\n", len(outlines))

	// 检查标题重复
	titleCount := make(map[string]int)
	titleToOutlines := make(map[string][]bson.M)
	for _, outline := range outlines {
		title := outline["title"].(string)
		titleCount[title]++
		titleToOutlines[title] = append(titleToOutlines[title], outline)
	}

	fmt.Println("标题重复统计:")
	hasDuplicates := false
	for title, count := range titleCount {
		if count > 1 {
			hasDuplicates = true
			fmt.Printf("  \"%s\": %d 个\n", title, count)
		}
	}
	if !hasDuplicates {
		fmt.Println("  无重复")
	}

	fmt.Println("\n详细列表:")
	for i, outline := range outlines {
		id := outline["_id"].(primitive.ObjectID).Hex()
		title := outline["title"].(string)
		docID := ""
		if docIDField, ok := outline["document_id"]; ok && docIDField != nil {
			docID = docIDField.(string)
		}
		parentID := ""
		if parentIDField, ok := outline["parent_id"]; ok && parentIDField != nil {
			parentID = parentIDField.(string)
		}
		order := outline["order"].(int32)

		// 检测是否为重复项
		isDup := titleCount[title] > 1
		dupMark := ""
		if isDup {
			dupMark = " [重复!]"
		}

		fmt.Printf("%d. [%s] \"%s\"%s doc=%s parent=%s order=%d\n",
			i+1, id[:8], title, dupMark, docID, parentID, order)
	}

	// 显示重复项的详细信息
	fmt.Println("\n=== 重复项详细对比 ===")
	for title, outlinesWithTitle := range titleToOutlines {
		if len(outlinesWithTitle) > 1 {
			fmt.Printf("\n标题: \"%s\"\n", title)
			for i, outline := range outlinesWithTitle {
				id := outline["_id"].(primitive.ObjectID).Hex()
				docID := ""
				if docIDField, ok := outline["document_id"]; ok && docIDField != nil {
					docID = docIDField.(string)
				}
				parentID := ""
				if parentIDField, ok := outline["parent_id"]; ok && parentIDField != nil {
					parentID = parentIDField.(string)
				}
				order := outline["order"].(int32)
				createdAt, updatedAt := time.Time{}, time.Time{}
				if ct, ok := outline["created_at"]; ok && ct != nil {
					// Handle both time.Time and primitive.DateTime
					switch v := ct.(type) {
					case time.Time:
						createdAt = v
					case primitive.DateTime:
						createdAt = v.Time()
					}
				}
				if ut, ok := outline["updated_at"]; ok && ut != nil {
					switch v := ut.(type) {
					case time.Time:
						updatedAt = v
					case primitive.DateTime:
						updatedAt = v.Time()
					}
				}

				fmt.Printf("  [%d] ID=%s\n", i+1, id)
				fmt.Printf("      document_id=%s\n", docID)
				fmt.Printf("      parent_id=%s\n", parentID)
				fmt.Printf("      order=%d\n", order)
				fmt.Printf("      created_at=%v\n", createdAt)
				fmt.Printf("      updated_at=%v\n", updatedAt)
			}
		}
	}
}
