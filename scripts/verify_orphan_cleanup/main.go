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

type VerificationResult struct {
	CollectionName string
	TotalCount     int64
	OrphanCount    int64
	Status         string
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	fmt.Println("========================================")
	fmt.Println("孤儿记录清理验证")
	fmt.Println("========================================")
	fmt.Printf("验证时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("数据库: qingyu\n")
	fmt.Println()

	results := make([]VerificationResult, 0)

	// 1. 验证 book_lists
	fmt.Println("【1/4】验证 book_lists...")
	result := verifyBookLists(ctx, db)
	results = append(results, result)
	fmt.Printf("  总记录: %d, 孤儿记录: %d, 状态: %s\n", result.TotalCount, result.OrphanCount, result.Status)
	fmt.Println()

	// 2. 验证 author_revenue
	fmt.Println("【2/4】验证 author_revenue...")
	result = verifyAuthorRevenue(ctx, db)
	results = append(results, result)
	fmt.Printf("  总记录: %d, 孤儿记录: %d, 状态: %s\n", result.TotalCount, result.OrphanCount, result.Status)
	fmt.Println()

	// 3. 验证 likes
	fmt.Println("【3/4】验证 likes...")
	result = verifyLikes(ctx, db)
	results = append(results, result)
	fmt.Printf("  总记录: %d, 孤儿记录: %d, 状态: %s\n", result.TotalCount, result.OrphanCount, result.Status)
	fmt.Println()

	// 4. 验证 notifications
	fmt.Println("【4/4】验证 notifications...")
	result = verifyNotifications(ctx, db)
	results = append(results, result)
	fmt.Printf("  总记录: %d, 孤儿记录: %d, 状态: %s\n", result.TotalCount, result.OrphanCount, result.Status)
	fmt.Println()

	// 汇总
	fmt.Println("========================================")
	fmt.Println("验证汇总")
	fmt.Println("========================================")

	totalOrphans := int64(0)
	allClean := true
	for _, r := range results {
		totalOrphans += r.OrphanCount
		if r.OrphanCount > 0 {
			allClean = false
		}
	}

	fmt.Printf("总孤儿记录数: %d\n", totalOrphans)
	if allClean {
		fmt.Println("验证结果: ✅ 所有集合无孤儿记录")
	} else {
		fmt.Println("验证结果: ❌ 仍有孤儿记录需要处理")
	}
	fmt.Println()

	// 详细结果
	fmt.Println("详细结果:")
	fmt.Println()
	for _, r := range results {
		fmt.Printf("- %s: %d/%d 孤儿记录 %s\n", r.CollectionName, r.OrphanCount, r.TotalCount, r.Status)
	}
}

func verifyBookLists(ctx context.Context, db *mongo.Database) VerificationResult {
	result := VerificationResult{CollectionName: "book_lists"}

	collection := db.Collection("book_lists")

	// 统计总数
	totalCount, _ := collection.CountDocuments(ctx, bson.D{})
	result.TotalCount = totalCount

	// 统计孤儿记录
	orphanCount, err := collection.CountDocuments(ctx, bson.M{
		"user_id": bson.M{"$nin": bson.A{"$users.distinct", "_id"}},
	})
	if err != nil {
		// 使用聚合管道方式
		pipeline := mongo.Pipeline{
			bson.D{{"$lookup", bson.D{
				{"from", "users"},
				{"localField", "user_id"},
				{"foreignField", "_id"},
				{"as", "user"},
			}}},
			bson.D{{"$match", bson.D{{"user", bson.D{{"$size", 0}}}}}},
			bson.D{{"$count", "orphaned"}},
		}
		cursor, _ := collection.Aggregate(ctx, pipeline)
		var counts []bson.M
		cursor.All(ctx, &counts)
		if len(counts) > 0 {
			result.OrphanCount = counts[0]["orphaned"].(int64)
		}
	} else {
		result.OrphanCount = orphanCount
	}

	if result.OrphanCount == 0 {
		result.Status = "✅"
	} else {
		result.Status = "❌"
	}

	return result
}

func verifyAuthorRevenue(ctx context.Context, db *mongo.Database) VerificationResult {
	result := VerificationResult{CollectionName: "author_revenue"}

	collection := db.Collection("author_revenue")

	// 统计总数
	totalCount, _ := collection.CountDocuments(ctx, bson.D{})
	result.TotalCount = totalCount

	// 使用聚合管道统计孤儿记录
	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$match", bson.D{{"user", bson.D{{"$size", 0}}}}}},
		bson.D{{"$count", "orphaned"}},
	}
	cursor, _ := collection.Aggregate(ctx, pipeline)
	var counts []bson.M
	cursor.All(ctx, &counts)
	if len(counts) > 0 {
		result.OrphanCount = counts[0]["orphaned"].(int64)
	}

	if result.OrphanCount == 0 {
		result.Status = "✅"
	} else {
		result.Status = "❌"
	}

	return result
}

func verifyLikes(ctx context.Context, db *mongo.Database) VerificationResult {
	result := VerificationResult{CollectionName: "likes"}

	collection := db.Collection("likes")

	// 统计总数
	totalCount, _ := collection.CountDocuments(ctx, bson.D{})
	result.TotalCount = totalCount

	// 使用聚合管道统计孤儿记录（user_id无效）
	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$match", bson.D{{"user", bson.D{{"$size", 0}}}}}},
		bson.D{{"$count", "orphaned"}},
	}
	cursor, _ := collection.Aggregate(ctx, pipeline)
	var counts []bson.M
	cursor.All(ctx, &counts)
	if len(counts) > 0 {
		result.OrphanCount = counts[0]["orphaned"].(int64)
	}

	if result.OrphanCount == 0 {
		result.Status = "✅"
	} else {
		result.Status = "❌"
	}

	return result
}

func verifyNotifications(ctx context.Context, db *mongo.Database) VerificationResult {
	result := VerificationResult{CollectionName: "notifications"}

	collection := db.Collection("notifications")

	// 统计总数
	totalCount, _ := collection.CountDocuments(ctx, bson.D{})
	result.TotalCount = totalCount

	// 使用聚合管道统计孤儿记录
	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$match", bson.D{{"user", bson.D{{"$size", 0}}}}}},
		bson.D{{"$count", "orphaned"}},
	}
	cursor, _ := collection.Aggregate(ctx, pipeline)
	var counts []bson.M
	cursor.All(ctx, &counts)
	if len(counts) > 0 {
		result.OrphanCount = counts[0]["orphaned"].(int64)
	}

	if result.OrphanCount == 0 {
		result.Status = "✅"
	} else {
		result.Status = "❌"
	}

	return result
}
