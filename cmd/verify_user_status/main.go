// Package main 提供验证用户状态字段的工具
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
	// 连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	usersCol := db.Collection("users")

	fmt.Println("========================================")
	fmt.Println("    用户状态字段验证报告")
	fmt.Println("========================================")
	fmt.Println()

	// 统计总用户数
	totalUsers, err := usersCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatalf("统计总用户数失败: %v", err)
	}
	fmt.Printf("总用户数: %d\n", totalUsers)
	fmt.Println()

	// 检查状态为空的用户
	emptyFilter := bson.M{
		"$or": []bson.M{
			{"status": bson.M{"$exists": false}},
			{"status": nil},
			{"status": ""},
		},
	}
	emptyStatusCount, err := usersCol.CountDocuments(ctx, emptyFilter)
	if err != nil {
		log.Fatalf("统计状态为空的用户数失败: %v", err)
	}
	fmt.Printf("状态为空的用户数: %d\n", emptyStatusCount)
	fmt.Println()

	if emptyStatusCount > 0 {
		fmt.Println("⚠ 警告: 仍有用户状态为空！")
		fmt.Println("状态为空的用户详情:")
		cursor, err := usersCol.Find(ctx, emptyFilter, options.Find().SetProjection(bson.M{
			"_id":      1,
			"username": 1,
			"status":   1,
		}))
		if err == nil {
			var results []bson.M
			if err = cursor.All(ctx, &results); err == nil {
				for _, r := range results {
					fmt.Printf("  ID: %v, Username: %v, Status: %v\n",
						r["_id"], r["username"], r["status"])
				}
			}
			cursor.Close(ctx)
		}
	} else {
		fmt.Println("✓ 所有用户都已有状态字段")
	}
	fmt.Println()

	// 状态分布
	fmt.Println("用户状态分布:")
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$status", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}},
	}
	cursor, err := usersCol.Aggregate(ctx, pipeline)
	if err == nil {
		var results []bson.M
		if err = cursor.All(ctx, &results); err == nil {
			for _, r := range results {
				status := r["_id"]
				var count int64
				switch v := r["count"].(type) {
				case int32:
					count = int64(v)
				case int64:
					count = v
				case int:
					count = int64(v)
				}
				fmt.Printf("  %v: %d\n", status, count)
			}
		}
		cursor.Close(ctx)
	}
	fmt.Println()

	// 示例用户
	fmt.Println("示例用户（前5条）:")
	cursor, err = usersCol.Find(ctx, bson.M{}, options.Find().SetLimit(5).SetProjection(bson.M{
		"_id":      1,
		"username": 1,
		"status":   1,
	}))
	if err == nil {
		var results []bson.M
		if err = cursor.All(ctx, &results); err == nil {
			for _, r := range results {
				fmt.Printf("  ID: %v, Username: %v, Status: %v\n",
					r["_id"], r["username"], r["status"])
			}
		}
		cursor.Close(ctx)
	}
	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("    验证结果")
	fmt.Println("========================================")
	if emptyStatusCount == 0 {
		fmt.Println("✓ 验证通过：所有用户都已有状态字段")
	} else {
		fmt.Printf("⚠ 验证失败：仍有 %d 个用户状态为空\n", emptyStatusCount)
	}
}
