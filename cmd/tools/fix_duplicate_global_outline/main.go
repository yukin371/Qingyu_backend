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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?directConnection=true"))
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	coll := db.Collection("outlines")

	// 1. 查找所有有重复 global outline 的项目
	fmt.Println("=== 查找重复的 global outline ===")

	// 按 project_id 分组，查找 type="global" 的记录
	pipeline := []bson.M{
		{"$match": bson.M{"type": "global"}},
		{"$group": bson.M{
			"_id": bson.M{
				"project_id": "$project_id",
			},
			"count": bson.M{"$sum": 1},
			"ids": bson.M{"$push": "$_id"},
			"first_id": bson.M{"$first": "$_id"},
		}},
		{"$match": bson.M{"count": bson.M{"$gt": 1}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatalf("聚合查询失败: %v", err)
	}
	defer cursor.Close(ctx)

	var duplicates []struct {
		ID       struct {
			ProjectID primitive.ObjectID `bson:"project_id"`
		} `bson:"_id"`
		Count    int                  `bson:"count"`
		IDs      []primitive.ObjectID `bson:"ids"`
		FirstID  primitive.ObjectID   `bson:"first_id"`
	}
	if err = cursor.All(ctx, &duplicates); err != nil {
		log.Fatalf("解码失败: %v", err)
	}

	fmt.Printf("找到 %d 个项目有重复的 global outline\n", len(duplicates))

	// 2. 为每个有重复的项目，保留第一个，删除其他的
	for _, dup := range duplicates {
		projectID := dup.ID.ProjectID.Hex()
		fmt.Printf("\n处理项目 %s:\n", projectID)
		fmt.Printf("  重复数量: %d\n", dup.Count)
		fmt.Printf("  保留 ID: %s\n", dup.FirstID.Hex())

		// 删除除第一个之外的所有 global outline
		idsToDelete := dup.IDs[1:] // 跳过第一个
		for i, id := range idsToDelete {
			fmt.Printf("  删除 ID[%d]: %s\n", i+1, id.Hex())
			_, err := coll.DeleteOne(ctx, bson.M{"_id": id})
			if err != nil {
				log.Printf("  删除失败: %v", err)
			}
		}
	}

	// 3. 验证结果
	fmt.Println("\n=== 验证结果 ===")
	cursor2, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatalf("验证查询失败: %v", err)
	}
	defer cursor2.Close(ctx)

	var remainingDuplicates []struct {
		ID       struct {
			ProjectID primitive.ObjectID `bson:"project_id"`
		} `bson:"_id"`
		Count int `bson:"count"`
	}
	if err = cursor2.All(ctx, &remainingDuplicates); err != nil {
		log.Fatalf("验证解码失败: %v", err)
	}

	if len(remainingDuplicates) == 0 {
		fmt.Println("✓ 所有重复的 global outline 已清理完成！")
	} else {
		fmt.Printf("⚠ 仍有 %d 个项目存在重复\n", len(remainingDuplicates))
	}

	// 4. 显示最终状态
	fmt.Println("\n=== 最终状态 ===")
	cursor3, err := coll.Find(ctx, bson.M{"type": "global"})
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	defer cursor3.Close(ctx)

	var globalOutlines []bson.M
	if err = cursor3.All(ctx, &globalOutlines); err != nil {
		log.Fatalf("解码失败: %v", err)
	}

	fmt.Printf("总共 %d 个 global outline:\n", len(globalOutlines))
	for _, outline := range globalOutlines {
		projectID := "unknown"
		if pid, ok := outline["project_id"].(primitive.ObjectID); ok {
			projectID = pid.Hex()
		}
		fmt.Printf("  - ID: %s, Project: %s\n", outline["_id"], projectID[:8]+"...")
	}
}