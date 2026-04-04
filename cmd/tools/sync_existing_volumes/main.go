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

	// 获取项目ID（从命令行参数或默认值）
	projectID := "69c790d7be2b6d81b914bfed"
	projectOID, _ := primitive.ObjectIDFromHex(projectID)

	// 查询所有 volume 类型且 outline_node_id 为空的文档
	cursor, err := db.Collection("documents").Find(ctx, bson.M{
		"project_id": projectOID,
		"type":       "volume",
		"$or": []bson.M{
			{"outline_node_id": bson.M{"$exists": false}},
			{"outline_node_id": ""},
		},
	})
	if err != nil {
		log.Printf("查询失败: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var documents []bson.M
	if err = cursor.All(ctx, &documents); err != nil {
		log.Printf("解析文档失败: %v", err)
		return
	}

	fmt.Printf("找到 %d 个需要同步的卷文档\n", len(documents))

	// 为每个文档创建对应的大纲节点
	syncedCount := 0
	for _, doc := range documents {
		docID := doc["_id"].(primitive.ObjectID).Hex()
		title := doc["title"].(string)
		order := int(doc["order"].(int32))

		// 创建大纲节点
		outlineID := primitive.NewObjectID()
		outline := bson.M{
			"_id":        outlineID,
			"project_id": projectOID,
			"title":      title,
			"parent_id":  "",
			"order":      order,
			"type":       "arc",
			"summary":    "",
			"tension":    5,
			"document_id": docID,
			"characters": bson.A{},
			"items":      bson.A{},
			"tags":       bson.A{},
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}

		_, err := db.Collection("outlines").InsertOne(ctx, outline)
		if err != nil {
			log.Printf("创建大纲失败 [%s]: %v", title, err)
			continue
		}

		// 更新文档的 outline_node_id
		_, err = db.Collection("documents").UpdateOne(ctx,
			bson.M{"_id": doc["_id"]},
			bson.M{"$set": bson.M{"outline_node_id": outlineID.Hex()}})
		if err != nil {
			log.Printf("更新文档 outline_node_id 失败 [%s]: %v", title, err)
		}

		fmt.Printf("✓ 同步完成: %s → outline=%s\n", title, outlineID.Hex())
		syncedCount++
	}

	fmt.Printf("\n同步完成: %d/%d\n", syncedCount, len(documents))
}
