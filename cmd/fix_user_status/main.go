// Package main 提供修复用户状态字段的工具
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

// FixReport 修复报告
type FixReport struct {
	TotalUsers       int64              `json:"totalUsers"`       // 总用户数
	EmptyStatusCount int64              `json:"emptyStatusCount"` // 状态为空的用户数
	UpdatedCount     int64              `json:"updatedCount"`     // 更新的用户数
	StatusBefore     []bson.M           `json:"statusBefore"`     // 更新前的状态分布示例
	StatusAfter      map[string]int64   `json:"statusAfter"`      // 更新后的状态分布
	VerifiedCount    int64              `json:"verifiedCount"`    // 验证后仍为空的用户数
}

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

	report := FixReport{
		StatusAfter: make(map[string]int64),
	}

	fmt.Println("========================================")
	fmt.Println("    用户状态字段修复工具")
	fmt.Println("========================================")
	fmt.Println()

	// 步骤1: 检查用户状态现状
	fmt.Println("步骤1: 检查用户状态现状...")
	fmt.Println("----------------------------------------")

	// 统计总用户数
	totalUsers, err := usersCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatalf("统计总用户数失败: %v", err)
	}
	report.TotalUsers = totalUsers
	fmt.Printf("总用户数: %d\n", totalUsers)
	fmt.Println()

	// 统计状态为空的用户数
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
	report.EmptyStatusCount = emptyStatusCount
	fmt.Printf("状态为空的用户数: %d\n", emptyStatusCount)
	fmt.Println()

	if emptyStatusCount == 0 {
		fmt.Println("✓ 所有用户都已有状态字段，无需修复")
		return
	}

	// 查看状态为空的用户详情（前10条）
	fmt.Println("状态为空的用户详情（前10条）:")
	cursor, err := usersCol.Find(ctx, emptyFilter, options.Find().SetLimit(10).SetProjection(bson.M{
		"_id":      1,
		"username": 1,
		"roles":    1,
		"status":   1,
	}))
	if err == nil {
		var results []bson.M
		if err = cursor.All(ctx, &results); err == nil {
			report.StatusBefore = results
			for _, r := range results {
				fmt.Printf("  ID: %v, Username: %v, Roles: %v, Status: %v\n",
					r["_id"], r["username"], r["roles"], r["status"])
			}
		}
		cursor.Close(ctx)
	}
	fmt.Println()

	// 步骤2: 设置默认状态
	fmt.Println("步骤2: 为状态为空的用户设置默认状态 'active'...")
	fmt.Println("----------------------------------------")

	updateResult, err := usersCol.UpdateMany(ctx, emptyFilter, bson.M{
		"$set": bson.M{"status": "active"},
	})
	if err != nil {
		log.Fatalf("更新用户状态失败: %v", err)
	}
	report.UpdatedCount = updateResult.ModifiedCount
	fmt.Printf("更新的记录数: %d\n", updateResult.ModifiedCount)
	fmt.Println()

	// 步骤3: 验证结果
	fmt.Println("步骤3: 验证修复结果...")
	fmt.Println("----------------------------------------")

	// 确认无状态为空的用户
	verifiedEmptyCount, err := usersCol.CountDocuments(ctx, emptyFilter)
	if err != nil {
		log.Fatalf("验证失败: %v", err)
	}
	report.VerifiedCount = verifiedEmptyCount
	fmt.Printf("验证后状态为空的用户数: %d\n", verifiedEmptyCount)

	if verifiedEmptyCount > 0 {
		fmt.Println("⚠ 警告: 仍有用户状态为空！")
	} else {
		fmt.Println("✓ 所有用户都已有状态字段")
	}
	fmt.Println()

	// 查看更新后的状态分布
	fmt.Println("更新后的状态分布:")
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$status", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}},
	}
	cursor, err = usersCol.Aggregate(ctx, pipeline)
	if err == nil {
		var results []bson.M
		if err = cursor.All(ctx, &results); err == nil {
			for _, r := range results {
				status := r["_id"]
				// 处理不同的整数类型
				var count int64
				switch v := r["count"].(type) {
				case int32:
					count = int64(v)
				case int64:
					count = v
				case int:
					count = int64(v)
				default:
					count = 0
				}
				report.StatusAfter[fmt.Sprintf("%v", status)] = count
				fmt.Printf("  %v: %d\n", status, count)
			}
		}
		cursor.Close(ctx)
	}
	fmt.Println()

	// 打印修复摘要
	fmt.Println("========================================")
	fmt.Println("    修复摘要")
	fmt.Println("========================================")
	fmt.Printf("总用户数: %d\n", report.TotalUsers)
	fmt.Printf("修复前状态为空的用户数: %d\n", report.EmptyStatusCount)
	fmt.Printf("更新的记录数: %d\n", report.UpdatedCount)
	fmt.Printf("验证后状态为空的用户数: %d\n", report.VerifiedCount)
	fmt.Println()

	if report.VerifiedCount == 0 && report.UpdatedCount == report.EmptyStatusCount {
		fmt.Println("✓ 修复成功！所有用户都已设置默认状态 'active'")
	} else {
		fmt.Println("⚠ 修复未完全成功，请检查日志")
	}
}
