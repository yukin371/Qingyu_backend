//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CleanupResult 清理结果统计
type CleanupResult struct {
	CollectionName string
	BeforeCount    int64
	AfterCount     int64
	DeletedCount   int64
	Error          error
}

// CleanupReport 清理报告
type CleanupReport struct {
	StartTime    time.Time
	EndTime      time.Time
	Results      []CleanupResult
	TotalDeleted int64
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	report := &CleanupReport{
		StartTime: time.Now(),
		Results:   make([]CleanupResult, 0),
	}

	fmt.Println("========================================")
	fmt.Println("孤儿记录清理工具")
	fmt.Println("========================================")
	fmt.Printf("开始时间: %s\n", report.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("数据库: qingyu\n")
	fmt.Println()

	// 1. 清理 book_lists 孤儿记录
	fmt.Println("【1/4】清理 book_lists 孤儿记录...")
	result := cleanupBookLists(ctx, db)
	report.Results = append(report.Results, result)
	fmt.Printf("  清理前: %d 条\n", result.BeforeCount)
	fmt.Printf("  清理后: %d 条\n", result.AfterCount)
	fmt.Printf("  删除: %d 条\n", result.DeletedCount)
	if result.Error != nil {
		fmt.Printf("  错误: %v\n", result.Error)
	}
	fmt.Println()

	// 2. 清理 author_revenue 孤儿记录
	fmt.Println("【2/4】清理 author_revenue 孤儿记录...")
	result = cleanupAuthorRevenue(ctx, db)
	report.Results = append(report.Results, result)
	fmt.Printf("  清理前: %d 条\n", result.BeforeCount)
	fmt.Printf("  清理后: %d 条\n", result.AfterCount)
	fmt.Printf("  删除: %d 条\n", result.DeletedCount)
	if result.Error != nil {
		fmt.Printf("  错误: %v\n", result.Error)
	}
	fmt.Println()

	// 3. 清理 likes 孤儿记录
	fmt.Println("【3/4】清理 likes 孤儿记录...")
	result = cleanupLikes(ctx, db)
	report.Results = append(report.Results, result)
	fmt.Printf("  清理前: %d 条\n", result.BeforeCount)
	fmt.Printf("  清理后: %d 条\n", result.AfterCount)
	fmt.Printf("  删除: %d 条\n", result.DeletedCount)
	if result.Error != nil {
		fmt.Printf("  错误: %v\n", result.Error)
	}
	fmt.Println()

	// 4. 清理 notifications 孤儿记录
	fmt.Println("【4/4】清理 notifications 孤儿记录...")
	result = cleanupNotifications(ctx, db)
	report.Results = append(report.Results, result)
	fmt.Printf("  清理前: %d 条\n", result.BeforeCount)
	fmt.Printf("  清理后: %d 条\n", result.AfterCount)
	fmt.Printf("  删除: %d 条\n", result.DeletedCount)
	if result.Error != nil {
		fmt.Printf("  错误: %v\n", result.Error)
	}
	fmt.Println()

	// 计算总删除数
	for _, r := range report.Results {
		report.TotalDeleted += r.DeletedCount
	}

	report.EndTime = time.Now()

	// 打印汇总
	fmt.Println("========================================")
	fmt.Println("清理汇总")
	fmt.Println("========================================")
	fmt.Printf("结束时间: %s\n", report.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("总耗时: %s\n", report.EndTime.Sub(report.StartTime).Round(time.Second))
	fmt.Printf("总删除: %d 条\n", report.TotalDeleted)
	fmt.Println()

	// 生成Markdown报告
	generateReport(report)

	fmt.Println("报告已生成: docs/reports/2026-02-01-orphan-cleanup-report.md")
}

// cleanupBookLists 清理 book_lists 孤儿记录
func cleanupBookLists(ctx context.Context, db *mongo.Database) CleanupResult {
	result := CleanupResult{CollectionName: "book_lists"}

	collection := db.Collection("book_lists")
	usersCollection := db.Collection("users")

	// 统计清理前数量
	beforeCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("统计清理前数量失败: %w", err)
		return result
	}
	result.BeforeCount = beforeCount

	// 获取所有有效的 user_id
	validUserIDs, err := usersCollection.Distinct(ctx, "_id", bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("获取有效用户ID失败: %w", err)
		return result
	}

	// 删除孤儿记录
	deleteResult, err := collection.DeleteMany(ctx, bson.M{
		"user_id": bson.M{"$nin": validUserIDs},
	})
	if err != nil {
		result.Error = fmt.Errorf("删除孤儿记录失败: %w", err)
		return result
	}

	// 统计清理后数量
	afterCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("统计清理后数量失败: %w", err)
		return result
	}

	result.AfterCount = afterCount
	result.DeletedCount = deleteResult.DeletedCount
	return result
}

// cleanupAuthorRevenue 清理 author_revenue 孤儿记录
func cleanupAuthorRevenue(ctx context.Context, db *mongo.Database) CleanupResult {
	result := CleanupResult{CollectionName: "author_revenue"}

	collection := db.Collection("author_revenue")
	usersCollection := db.Collection("users")

	// 统计清理前数量
	beforeCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("统计清理前数量失败: %w", err)
		return result
	}
	result.BeforeCount = beforeCount

	// 获取所有有效的 user_id
	validUserIDs, err := usersCollection.Distinct(ctx, "_id", bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("获取有效用户ID失败: %w", err)
		return result
	}

	// 删除孤儿记录
	deleteResult, err := collection.DeleteMany(ctx, bson.M{
		"user_id": bson.M{"$nin": validUserIDs},
	})
	if err != nil {
		result.Error = fmt.Errorf("删除孤儿记录失败: %w", err)
		return result
	}

	// 统计清理后数量
	afterCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("统计清理后数量失败: %w", err)
		return result
	}

	result.AfterCount = afterCount
	result.DeletedCount = deleteResult.DeletedCount
	return result
}

// cleanupLikes 清理 likes 孤儿记录
func cleanupLikes(ctx context.Context, db *mongo.Database) CleanupResult {
	result := CleanupResult{CollectionName: "likes"}

	collection := db.Collection("likes")
	usersCollection := db.Collection("users")
	booksCollection := db.Collection("books")
	commentsCollection := db.Collection("comments")

	// 统计清理前数量
	beforeCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("统计清理前数量失败: %w", err)
		return result
	}
	result.BeforeCount = beforeCount

	// 获取所有有效的 user_id
	validUserIDs, err := usersCollection.Distinct(ctx, "_id", bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("获取有效用户ID失败: %w", err)
		return result
	}

	// 获取所有有效的 book_id
	validBookIDs, err := booksCollection.Distinct(ctx, "_id", bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("获取有效书籍ID失败: %w", err)
		return result
	}

	// 获取所有有效的 comment_id
	validCommentIDs, err := commentsCollection.Distinct(ctx, "_id", bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("获取有效评论ID失败: %w", err)
		return result
	}

	// 删除孤儿记录（user_id无效 或 (target_type=book且target_id无效) 或 (target_type=comment且target_id无效)）
	deleteResult, err := collection.DeleteMany(ctx, bson.M{
		"$or": []bson.M{
			{"user_id": bson.M{"$nin": validUserIDs}},
			{
				"$and": []bson.M{
					{"target_type": "book"},
					{"target_id": bson.M{"$nin": validBookIDs}},
				},
			},
			{
				"$and": []bson.M{
					{"target_type": "comment"},
					{"target_id": bson.M{"$nin": validCommentIDs}},
				},
			},
		},
	})
	if err != nil {
		result.Error = fmt.Errorf("删除孤儿记录失败: %w", err)
		return result
	}

	// 统计清理后数量
	afterCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("统计清理后数量失败: %w", err)
		return result
	}

	result.AfterCount = afterCount
	result.DeletedCount = deleteResult.DeletedCount
	return result
}

// cleanupNotifications 清理 notifications 孤儿记录
func cleanupNotifications(ctx context.Context, db *mongo.Database) CleanupResult {
	result := CleanupResult{CollectionName: "notifications"}

	collection := db.Collection("notifications")
	usersCollection := db.Collection("users")

	// 统计清理前数量
	beforeCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("统计清理前数量失败: %w", err)
		return result
	}
	result.BeforeCount = beforeCount

	// 获取所有有效的 user_id
	validUserIDs, err := usersCollection.Distinct(ctx, "_id", bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("获取有效用户ID失败: %w", err)
		return result
	}

	// 删除孤儿记录
	deleteResult, err := collection.DeleteMany(ctx, bson.M{
		"user_id": bson.M{"$nin": validUserIDs},
	})
	if err != nil {
		result.Error = fmt.Errorf("删除孤儿记录失败: %w", err)
		return result
	}

	// 统计清理后数量
	afterCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		result.Error = fmt.Errorf("统计清理后数量失败: %w", err)
		return result
	}

	result.AfterCount = afterCount
	result.DeletedCount = deleteResult.DeletedCount
	return result
}

// generateReport 生成Markdown报告
func generateReport(report *CleanupReport) {
	// 确保报告目录存在
	os.MkdirAll("docs/reports", 0755)

	filename := "docs/reports/2026-02-01-orphan-cleanup-report.md"
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("创建报告文件失败: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "# 孤儿记录清理报告\n\n")
	fmt.Fprintf(file, "> **生成时间:** %s\n", report.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "> **执行时长:** %s\n", report.EndTime.Sub(report.StartTime).Round(time.Second))
	fmt.Fprintf(file, "> **数据库:** qingyu\n\n")

	fmt.Fprintf(file, "## 清理前统计\n\n")
	for _, result := range report.Results {
		fmt.Fprintf(file, "- **%s**: %,d 条\n", result.CollectionName, result.BeforeCount)
	}

	fmt.Fprintf(file, "\n## 清理后统计\n\n")
	for _, result := range report.Results {
		fmt.Fprintf(file, "- **%s**: %,d 条\n", result.CollectionName, result.AfterCount)
	}

	fmt.Fprintf(file, "\n## 清理数量\n\n")
	for _, result := range report.Results {
		fmt.Fprintf(file, "- **%s**: %,d 条", result.CollectionName, result.DeletedCount)
		if result.Error != nil {
			fmt.Fprintf(file, " ❌ (错误: %v)", result.Error)
		} else {
			fmt.Fprintf(file, " ✅")
		}
		fmt.Fprintf(file, "\n")
	}

	fmt.Fprintf(file, "\n- **总计**: %,d 条\n\n", report.TotalDeleted)

	// 验证结果
	fmt.Fprintf(file, "## 验证结果\n\n")
	allSuccess := true
	for _, result := range report.Results {
		if result.Error != nil {
			allSuccess = false
			fmt.Fprintf(file, "- **%s**: 清理失败 ❌\n", result.CollectionName)
		} else {
			fmt.Fprintf(file, "- **%s**: 清理成功 ✅\n", result.CollectionName)
		}
	}

	if allSuccess {
		fmt.Fprintf(file, "\n### 总体状态: 成功 ✅\n\n")
		fmt.Fprintf(file, "所有孤儿记录已成功清理，数据库关联完整性已恢复。\n")
	} else {
		fmt.Fprintf(file, "\n### 总体状态: 部分失败 ❌\n\n")
		fmt.Fprintf(file, "部分集合清理失败，请检查错误信息并重新执行清理。\n")
	}

	fmt.Fprintf(file, "\n## 清理详情\n\n")
	for _, result := range report.Results {
		fmt.Fprintf(file, "### %s\n\n", result.CollectionName)
		fmt.Fprintf(file, "| 指标 | 数值 |\n")
		fmt.Fprintf(file, "|------|------|\n")
		fmt.Fprintf(file, "| 清理前 | %,d |\n", result.BeforeCount)
		fmt.Fprintf(file, "| 清理后 | %,d |\n", result.AfterCount)
		fmt.Fprintf(file, "| 删除 | %,d |\n", result.DeletedCount)
		if result.Error != nil {
			fmt.Fprintf(file, "| 状态 | 失败 ❌ |\n")
			fmt.Fprintf(file, "\n**错误信息:**\n```\n%v\n```\n", result.Error)
		} else {
			fmt.Fprintf(file, "| 状态 | 成功 ✅ |\n")
		}
		fmt.Fprintf(file, "\n")
	}
}
