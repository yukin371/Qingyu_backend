package migration

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// FixCommentFields 修复评论集合中的字段名: status -> state, user_id -> author_id
type FixCommentFields struct{}

// Up 执行迁移:将 comments 集合中旧字段修复为新字段
func (m *FixCommentFields) Up(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("comments")
	log.Println("  [comments] 开始修复评论字段...")

	// Step 1: 修复 status -> state
	statusFilter := bson.M{"status": bson.M{"$exists": true}}
	statusCount, err := collection.CountDocuments(ctx, statusFilter)
	if err != nil {
		return fmt.Errorf("统计需要修复的评论(status字段): %w", err)
	}
	if statusCount > 0 {
		log.Printf("  [comments] 修复 %d 条评论(status -> state)...", statusCount)
		_, err := collection.UpdateMany(ctx, statusFilter, bson.M{
			"$set":   bson.M{"state": "normal"},
			"$unset": bson.M{"status": ""},
		})
		if err != nil {
			return fmt.Errorf("修复评论 state 字段失败: %w", err)
		}
		log.Printf("  [comments] 成功修复 %d 条评论(status -> state)", statusCount)
	} else {
		log.Println("  [comments] 无需修复 status 字段")
	}

	// Step 2: 修复 user_id -> author_id
	userIDFilter := bson.M{
		"user_id":   bson.M{"$exists": true},
		"author_id": bson.M{"$exists": false},
	}
	userIDCount, err := collection.CountDocuments(ctx, userIDFilter)
	if err != nil {
		return fmt.Errorf("统计需要修复的评论(user_id字段): %w", err)
	}
	if userIDCount > 0 {
		log.Printf("  [comments] 修复 %d 条评论(user_id -> author_id)...", userIDCount)
		// 用 aggregation pipeline 复制 user_id 到 author_id
		pipeline := []bson.M{
			{"$match": userIDFilter},
			{"$set": bson.M{"author_id": "$user_id"}},
			{"$unset": "user_id"},
			{"$merge": bson.M{"into": "comments", "whenMatched": "replace", "whenNotMatched": "discard"}},
		}
		_, err := collection.Aggregate(ctx, pipeline)
		if err != nil {
			return fmt.Errorf("修复评论 author_id 字段失败: %w", err)
		}
		log.Printf("  [comments] 成功修复 %d 条评论(user_id -> author_id)", userIDCount)
	} else {
		log.Println("  [comments] 无需修复 user_id 字段")
	}

	return nil
}

// Down 回滚迁移
func (m *FixCommentFields) Down(ctx context.Context, db *mongo.Database) error {
	log.Println("  [comments] 注意: 字段重命名迁移无法完全回滚，请手动检查数据")
	return nil
}
