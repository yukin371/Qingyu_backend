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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("qingyu").Collection("outlines")

	// 统计各层级节点数
	roots, _ := collection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"parent_id": ""},
			{"parent_id": bson.M{"$exists": false}},
		},
	})

	l2Nodes, _ := collection.CountDocuments(ctx, bson.M{
		"parent_id": bson.M{"$ne": "", "$exists": true},
	})

	fmt.Printf("根节点(parent_id为空): %d\n", roots)
	fmt.Printf("子节点(parent_id不为空): %d\n", l2Nodes)

	// 查询一个示例根节点
	var root bson.M
	collection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"parent_id": ""},
			{"parent_id": bson.M{"$exists": false}},
		},
	}, options.FindOne().SetProjection(bson.M{
		"title":     1,
		"parent_id": 1,
	})).Decode(&root)

	fmt.Printf("\n示例根节点:\n")
	fmt.Printf("  Title: %v\n", root["title"])
	fmt.Printf("  ParentID: %v\n", root["parent_id"])

	// 查询这个根节点的子节点数量
	if rootID, ok := root["_id"].(primitive.ObjectID); ok {
		childCount, _ := collection.CountDocuments(ctx, bson.M{
			"parent_id": rootID.Hex(),
		})
		fmt.Printf("  子节点数量: %d\n", childCount)
	}

	// 列出所有parent_id值
	pipeline := mongo.Pipeline{
		bson.D{bson.E{Key: "$group", Value: bson.D{
			bson.E{Key: "_id", Value: "$parent_id"},
			bson.E{Key: "count", Value: bson.D{bson.E{Key: "$sum", Value: 1}}},
		}}},
	}
	cursor, _ := collection.Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		log.Printf("聚合查询失败: %v", err)
		return
	}

	fmt.Printf("\nParent ID分布:\n")
	for _, r := range results {
		parentID := r["_id"]
		if parentID == nil || parentID == "" {
			fmt.Printf("  (根节点): %v\n", r["count"])
		} else {
			fmt.Printf("  %v: %v\n", parentID, r["count"])
		}
	}
}
