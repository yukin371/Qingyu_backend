package monitor

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB MongoDB数据库实现
type MongoDB struct {
	db *mongo.Database
}

// NewMongoDB 创建MongoDB数据库实例
func NewMongoDB(db *mongo.Database) *MongoDB {
	return &MongoDB{db: db}
}

// OrphanedRecords 检查孤儿记录数量
// 使用聚合管道查找外键指向不存在记录的文档
func (m *MongoDB) OrphanedRecords(ctx context.Context, collection, foreignKey, targetCollection string) (int, error) {
	// 构建聚合管道
	// 1. $lookup: 关联目标集合
	// 2. $match: 筛选出关联为空的记录（孤儿记录）
	// 3. $count: 统计数量

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{"from", targetCollection},
			{"localField", foreignKey},
			{"foreignField", "_id"},
			{"as", "ref"},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{"ref", bson.D{{"$size", 0}}},
		}}},
		bson.D{{Key: "$count", Value: "total"}},
	}

	cursor, err := m.db.Collection(collection).Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("执行聚合查询失败: %w", err)
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return 0, fmt.Errorf("解析聚合结果失败: %w", err)
	}

	if len(result) > 0 {
		if total, ok := result[0]["total"].(int32); ok {
			return int(total), nil
		}
	}

	return 0, nil
}

// StatisticsAccuracy 检查统计数据不准确的数量
// 比较count字段与实际关联记录的数量不一致的文档
func (m *MongoDB) StatisticsAccuracy(ctx context.Context, collection, countField string) (int, error) {
	// 针对不同的count字段使用不同的验证逻辑
	switch collection {
	case "books":
		return m.checkBooksStatistics(ctx, countField)
	case "users":
		return m.checkUsersStatistics(ctx, countField)
	default:
		return 0, fmt.Errorf("不支持的集合: %s", collection)
	}
}

// checkBooksStatistics 检查书籍统计数据准确性
func (m *MongoDB) checkBooksStatistics(ctx context.Context, countField string) (int, error) {
	var targetCollection string

	switch countField {
	case "likes_count":
		targetCollection = "likes"
	case "comments_count":
		targetCollection = "comments"
	default:
		return 0, nil
	}

	// 使用聚合管道比较count字段与实际数量
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{"from", targetCollection},
		{"let", bson.D{{"bookId", "$_id"}}},
		{"pipeline", bson.A{
			bson.D{{Key: "$match", Value: bson.D{
				{"$expr", bson.D{{"$eq", bson.A{"$$bookId", "$target_id"}}}},
			}}},
		}},
		{"as", "refs"},
	}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{"_id", 1},
		{"stored_count", bson.D{{"$ifNull", bson.A{fmt.Sprintf("$%s", countField), 0}}}},
		{"actual_count", bson.D{{"$size", "$refs"}}},
		{"diff", bson.D{{"$subtract", bson.A{
			bson.D{{"$ifNull", bson.A{fmt.Sprintf("$%s", countField), 0}}},
			bson.D{{"$size", "$refs"}},
		}}}},
	}}}

	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{"diff", bson.D{{"$ne", 0}}},
	}}}

	countStage := bson.D{{Key: "$count", Value: "total"}}

	pipeline := mongo.Pipeline{lookupStage, projectStage, matchStage, countStage}

	cursor, err := m.db.Collection("books").Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("执行统计准确性检查失败: %w", err)
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return 0, fmt.Errorf("解析统计结果失败: %w", err)
	}

	if len(result) > 0 {
		if total, ok := result[0]["total"].(int32); ok {
			return int(total), nil
		}
	}

	return 0, nil
}

// checkUsersStatistics 检查用户统计数据准确性
func (m *MongoDB) checkUsersStatistics(ctx context.Context, countField string) (int, error) {
	if countField != "followers_count" {
		return 0, nil
	}

	// 使用聚合管道比较followers_count与实际的followers数量
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{"from", "user_relations"},
		{"let", bson.D{{"userId", "$_id"}}},
		{"pipeline", bson.A{
			bson.D{{Key: "$match", Value: bson.D{
				{"$expr", bson.D{{"$eq", bson.A{"$$userId", "$followed_id"}}}},
			}}},
		}},
		{"as", "followers"},
	}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{"_id", 1},
		{"stored_count", bson.D{{"$ifNull", bson.A{"$followers_count", 0}}}},
		{"actual_count", bson.D{{"$size", "$followers"}}},
		{"diff", bson.D{{"$subtract", bson.A{
			bson.D{{"$ifNull", bson.A{"$followers_count", 0}}},
			bson.D{{"$size", "$followers"}},
		}}}},
	}}}

	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{"diff", bson.D{{"$ne", 0}}},
	}}}

	countStage := bson.D{{Key: "$count", Value: "total"}}

	pipeline := mongo.Pipeline{lookupStage, projectStage, matchStage, countStage}

	cursor, err := m.db.Collection("users").Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("执行用户统计准确性检查失败: %w", err)
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return 0, fmt.Errorf("解析用户统计结果失败: %w", err)
	}

	if len(result) > 0 {
		if total, ok := result[0]["total"].(int32); ok {
			return int(total), nil
		}
	}

	return 0, nil
}
