// Package utils 提供数据验证功能
package utils

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// VerificationReport 验证报告
type VerificationReport struct {
	Category string   // 验证类别
	Passed   bool     // 是否通过验证
	Issues   []string // 发现的问题列表
}

// VerifyData 执行所有数据验证
func VerifyData(db *Database) []VerificationReport {
	var reports []VerificationReport

	// 验证用户数据
	reports = append(reports, verifyUsers(db))

	// 验证书籍数据
	reports = append(reports, verifyBooks(db))

	// 验证订阅关系
	reports = append(reports, verifySubscriptions(db))

	return reports
}

// verifyUsers 验证用户数据
func verifyUsers(db *Database) VerificationReport {
	report := VerificationReport{
		Category: "用户数据",
		Passed:   true,
		Issues:   []string{},
	}

	ctx := context.Background()
	collection := db.Collection("users")

	// 检查用户名唯一性
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":       "$username",
				"count":     bson.M{"$sum": 1},
				"userIds":   bson.M{"$push": "$_id"},
			},
		},
		{
			"$match": bson.M{
				"count": bson.M{"$gt": 1},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		report.Passed = false
		report.Issues = append(report.Issues, fmt.Sprintf("查询失败: %v", err))
		return report
	}
	defer cursor.Close(ctx)

	var duplicates []bson.M
	if err = cursor.All(ctx, &duplicates); err != nil {
		report.Passed = false
		report.Issues = append(report.Issues, fmt.Sprintf("解析结果失败: %v", err))
		return report
	}

	if len(duplicates) > 0 {
		report.Passed = false
		for _, dup := range duplicates {
			username := dup["_id"]
			count := dup["count"]
			report.Issues = append(report.Issues,
				fmt.Sprintf("重复用户名 '%s' 出现 %d 次", username, count))
		}
	} else {
		report.Issues = append(report.Issues, "所有用户名唯一")
	}

	return report
}

// verifyBooks 验证书籍数据
func verifyBooks(db *Database) VerificationReport {
	report := VerificationReport{
		Category: "书籍数据",
		Passed:   true,
		Issues:   []string{},
	}

	ctx := context.Background()
	collection := db.Collection("books")

	// 检查评分范围 (0-10)
	filter := bson.M{
		"$or": []bson.M{
			{"rating": bson.M{"$lt": 0}},
			{"rating": bson.M{"$gt": 10}},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		report.Passed = false
		report.Issues = append(report.Issues, fmt.Sprintf("查询失败: %v", err))
		return report
	}
	defer cursor.Close(ctx)

	var invalidBooks []bson.M
	if err = cursor.All(ctx, &invalidBooks); err != nil {
		report.Passed = false
		report.Issues = append(report.Issues, fmt.Sprintf("解析结果失败: %v", err))
		return report
	}

	if len(invalidBooks) > 0 {
		report.Passed = false
		report.Issues = append(report.Issues,
			fmt.Sprintf("发现 %d 本评分超出范围 (0-10) 的书籍", len(invalidBooks)))
		for _, book := range invalidBooks {
			title := book["title"]
			rating := book["rating"]
			report.Issues = append(report.Issues,
				fmt.Sprintf("  - '%s' 评分: %.1f", title, rating))
		}
	} else {
		report.Issues = append(report.Issues, "所有书籍评分在有效范围内 (0-10)")
	}

	return report
}

// verifySubscriptions 验证订阅关系
func verifySubscriptions(db *Database) VerificationReport {
	report := VerificationReport{
		Category: "订阅关系",
		Passed:   true,
		Issues:   []string{},
	}

	ctx := context.Background()
	collection := db.Collection("subscriptions")

	// 使用 $lookup 检查孤儿订阅（用户不存在）
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		{
			"$match": bson.M{
				"user": bson.M{"$size": 0},
			},
		},
		{
			"$count": "orphanCount",
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		report.Passed = false
		report.Issues = append(report.Issues, fmt.Sprintf("查询失败: %v", err))
		return report
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		report.Passed = false
		report.Issues = append(report.Issues, fmt.Sprintf("解析结果失败: %v", err))
		return report
	}

	if len(results) > 0 && results[0]["orphanCount"].(int32) > 0 {
		report.Passed = false
		count := results[0]["orphanCount"]
		report.Issues = append(report.Issues,
			fmt.Sprintf("发现 %d 个孤儿订阅（用户不存在）", count))
	} else {
		report.Issues = append(report.Issues, "所有订阅关系有效")
	}

	return report
}
