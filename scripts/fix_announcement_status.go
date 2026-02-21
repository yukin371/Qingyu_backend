//go:build ignore

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

// Announcement 公告结构
type Announcement struct {
	ID      primitive.ObjectID `bson:"_id"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
	Status  string             `bson:"status"`
}

type StatusStats struct {
	Status string
	Count  int64
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 连接数据库
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/qingyu"))
	if err != nil {
		log.Fatal("连接 MongoDB 失败:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	collection := db.Collection("announcements")

	fmt.Println("========================================")
	fmt.Println("公告状态修复工具")
	fmt.Println("========================================")
	fmt.Printf("执行时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// 1. 检查公告总数
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("统计公告总数失败:", err)
	}
	fmt.Printf("公告总数: %d\n\n", totalCount)

	// 2. 检查状态为空的公告数量
	emptyStatusCount, err := collection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"status": bson.M{"$exists": false}},
			{"status": nil},
			{"status": ""},
		},
	})
	if err != nil {
		log.Fatal("统计空状态公告失败:", err)
	}
	fmt.Printf("状态为空的公告数量: %d\n\n", emptyStatusCount)

	// 3. 查看状态为空的公告详情
	fmt.Println("状态为空的公告详情:")
	if emptyStatusCount > 0 {
		cursor, err := collection.Find(ctx, bson.M{
			"$or": []bson.M{
				{"status": bson.M{"$exists": false}},
				{"status": nil},
				{"status": ""},
			},
		}, options.Find().SetLimit(20))
		if err != nil {
			log.Fatal("查询空状态公告失败:", err)
		}
		defer cursor.Close(ctx)

		var announcements []Announcement
		if err = cursor.All(ctx, &announcements); err != nil {
			log.Fatal("解析空状态公告失败:", err)
		}

		for _, ann := range announcements {
			fmt.Printf("  ID: %s, Title: %s, Status: '%s'\n", ann.ID.Hex(), ann.Title, ann.Status)
		}
		if emptyStatusCount > 20 {
			fmt.Printf("  ... 还有 %d 条\n", emptyStatusCount-20)
		}
	} else {
		fmt.Println("  无状态为空的公告")
	}
	fmt.Println()

	// 4. 查看更新前的状态分布
	fmt.Println("更新前的状态分布:")
	preStats := getStatusStats(ctx, collection)
	for _, stat := range preStats {
		fmt.Printf("  %s: %d\n", stat.Status, stat.Count)
	}
	fmt.Println()

	// 5. 更新状态为空的公告
	if emptyStatusCount > 0 {
		fmt.Println("正在更新状态为空的公告...")
		result, err := collection.UpdateMany(ctx,
			bson.M{
				"$or": []bson.M{
					{"status": bson.M{"$exists": false}},
					{"status": nil},
					{"status": ""},
				},
			},
			bson.M{"$set": bson.M{"status": "published"}},
		)
		if err != nil {
			log.Fatal("更新公告状态失败:", err)
		}
		fmt.Printf("成功更新 %d 条公告\n\n", result.ModifiedCount)
	} else {
		fmt.Println("无需更新\n")
	}

	// 6. 验证结果
	fmt.Println("========================================")
	fmt.Println("验证结果")
	fmt.Println("========================================\n")

	// 确认无状态为空的公告
	postEmptyCount, err := collection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"status": bson.M{"$exists": false}},
			{"status": nil},
			{"status": ""},
		},
	})
	if err != nil {
		log.Fatal("验证空状态失败:", err)
	}
	fmt.Printf("状态为空的公告数量: %d ", postEmptyCount)
	if postEmptyCount == 0 {
		fmt.Println("✓")
	} else {
		fmt.Println("✗ (仍有空状态公告)")
	}
	fmt.Println()

	// 查看更新后的状态分布
	fmt.Println("更新后的状态分布:")
	postStats := getStatusStats(ctx, collection)
	for _, stat := range postStats {
		fmt.Printf("  %s: %d\n", stat.Status, stat.Count)
	}
	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("修复完成")
	fmt.Println("========================================")
}

// getStatusStats 获取状态分布统计
func getStatusStats(ctx context.Context, collection *mongo.Collection) []StatusStats {
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$status", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("获取状态分布失败: %v", err)
		return nil
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		log.Printf("解析状态分布失败: %v", err)
		return nil
	}

	stats := make([]StatusStats, 0, len(results))
	for _, r := range results {
		status := r.ID
		if status == "" {
			status = "(empty)"
		}
		stats = append(stats, StatusStats{
			Status: status,
			Count:  r.Count,
		})
	}

	return stats
}
