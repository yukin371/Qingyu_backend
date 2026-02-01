// Package validator provides data validation functionality for the seeder tool.
// It validates data relationships, detects orphaned records, and ensures ID format consistency.
package validator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DataValidator 数据验证器
type DataValidator struct {
	db *mongo.Database
}

// ValidationReport 验证报告
type ValidationReport struct {
	IsValid               bool      `json:"is_valid"`                 // 总体是否通过验证
	TotalOrphanedRecords  int       `json:"total_orphaned_records"`   // 孤儿记录总数
	InconsistentFields    []string  `json:"inconsistent_fields"`      // 格式不一致的字段
	InconsistentRecords   []string  `json:"inconsistent_records"`     // 格式不一致的记录
	OrphanDetails         []string  `json:"orphan_details"`           // 孤儿记录详情
	Errors                []error   `json:"errors"`                   // 错误列表
	CollectionStats       map[string]int64 `json:"collection_stats"` // 集合统计信息
	Summary               string    `json:"summary"`                  // 验证摘要
	ValidatedAt           time.Time `json:"validated_at"`             // 验证时间
}

// OrphanRecordDetail 孤儿记录详情
type OrphanRecordDetail struct {
	Collection string    `json:"collection"`  // 集合名称
	RecordID   string    `json:"record_id"`   // 记录ID
	Field      string    `json:"field"`       // 外键字段名
	ForeignKey string    `json:"foreign_key"` // 外键值
	TargetColl string    `json:"target_coll"` // 目标集合
}

// IDFormatIssue ID格式问题
type IDFormatIssue struct {
	Collection string `json:"collection"` // 集合名称
	Field      string `json:"field"`      // 字段名
	Expected   string `json:"expected"`   // 期望格式
	Actual     string `json:"actual"`     // 实际格式
}

// NewDataValidator 创建数据验证器
func NewDataValidator(db *mongo.Database) *DataValidator {
	return &DataValidator{
		db: db,
	}
}

// ValidateNoOrphanedRecords 验证无孤儿记录
// 检查所有外键引用是否都指向存在的记录
func (v *DataValidator) ValidateNoOrphanedRecords(ctx context.Context) (*ValidationReport, error) {
	report := &ValidationReport{
		ValidatedAt: time.Now(),
		OrphanDetails: make([]string, 0),
		Errors:       make([]error, 0),
	}

	// 验证各个集合的孤儿记录
	orphanCounts := make(map[string]int64)

	// 1. 检查 reading_progress 孤儿记录
	readingProgressOrphans, err := v.ValidateReadingProgressOrphans(ctx)
	if err != nil {
		report.Errors = append(report.Errors, fmt.Errorf("检查 reading_progress 失败: %w", err))
	} else if readingProgressOrphans > 0 {
		orphanCounts["reading_progress"] = readingProgressOrphans
		report.OrphanDetails = append(report.OrphanDetails,
			fmt.Sprintf("reading_progress: %d 条孤儿记录", readingProgressOrphans))
	}

	// 2. 检查 reading_histories 孤儿记录
	readingHistoryOrphans, err := v.ValidateReadingHistoryOrphans(ctx)
	if err != nil {
		report.Errors = append(report.Errors, fmt.Errorf("检查 reading_histories 失败: %w", err))
	} else if readingHistoryOrphans > 0 {
		orphanCounts["reading_histories"] = readingHistoryOrphans
		report.OrphanDetails = append(report.OrphanDetails,
			fmt.Sprintf("reading_histories: %d 条孤儿记录", readingHistoryOrphans))
	}

	// 3. 检查 bookmarks 孤儿记录
	bookmarkOrphans, err := v.ValidateBookmarkOrphans(ctx)
	if err != nil {
		report.Errors = append(report.Errors, fmt.Errorf("检查 bookmarks 失败: %w", err))
	} else if bookmarkOrphans > 0 {
		orphanCounts["bookmarks"] = bookmarkOrphans
		report.OrphanDetails = append(report.OrphanDetails,
			fmt.Sprintf("bookmarks: %d 条孤儿记录", bookmarkOrphans))
	}

	// 4. 检查 likes 孤儿记录
	likeOrphans, err := v.ValidateLikeOrphans(ctx)
	if err != nil {
		report.Errors = append(report.Errors, fmt.Errorf("检查 likes 失败: %w", err))
	} else if likeOrphans > 0 {
		orphanCounts["likes"] = likeOrphans
		report.OrphanDetails = append(report.OrphanDetails,
			fmt.Sprintf("likes: %d 条孤儿记录", likeOrphans))
	}

	// 5. 检查 notifications 孤儿记录
	notificationOrphans, err := v.ValidateNotificationOrphans(ctx)
	if err != nil {
		report.Errors = append(report.Errors, fmt.Errorf("检查 notifications 失败: %w", err))
	} else if notificationOrphans > 0 {
		orphanCounts["notifications"] = notificationOrphans
		report.OrphanDetails = append(report.OrphanDetails,
			fmt.Sprintf("notifications: %d 条孤儿记录", notificationOrphans))
	}

	// 计算总孤儿记录数
	for _, count := range orphanCounts {
		report.TotalOrphanedRecords += int(count)
	}

	// 生成摘要
	if report.TotalOrphanedRecords == 0 {
		report.IsValid = true
		report.Summary = "未发现孤儿记录，数据完整性良好"
	} else {
		report.IsValid = false
		report.Summary = fmt.Sprintf("发现 %d 条孤儿记录，需要清理", report.TotalOrphanedRecords)
	}

	return report, nil
}

// ValidateIDFormatConsistency 验证ID格式一致性
// 检查外键字段是否使用正确的类型（string vs ObjectID）
func (v *DataValidator) ValidateIDFormatConsistency(ctx context.Context) (*ValidationReport, error) {
	report := &ValidationReport{
		ValidatedAt:       time.Now(),
		InconsistentFields: make([]string, 0),
		InconsistentRecords: make([]string, 0),
		Errors:            make([]error, 0),
	}

	// 检查 likes 集合（应该使用 string 类型）
	if err := v.validateLikesIDFormat(ctx, report); err != nil {
		report.Errors = append(report.Errors, fmt.Errorf("检查 likes ID格式失败: %w", err))
	}

	// 检查 notifications 集合（应该使用 string 类型）
	if err := v.validateNotificationsIDFormat(ctx, report); err != nil {
		report.Errors = append(report.Errors, fmt.Errorf("检查 notifications ID格式失败: %w", err))
	}

	// 检查 reading_progress 集合（使用 ObjectID - 旧格式）
	if err := v.validateReadingProgressIDFormat(ctx, report); err != nil {
		report.Errors = append(report.Errors, fmt.Errorf("检查 reading_progress ID格式失败: %w", err))
	}

	// 生成摘要
	if len(report.InconsistentRecords) == 0 {
		report.IsValid = true
		report.Summary = "ID格式一致性检查通过"
	} else {
		report.IsValid = false
		report.Summary = fmt.Sprintf("发现 %d 条记录ID格式不一致", len(report.InconsistentRecords))
	}

	return report, nil
}

// ValidateRelationships 综合验证所有关系
func (v *DataValidator) ValidateRelationships(ctx context.Context) (*ValidationReport, error) {
	// 1. 验证孤儿记录
	orphanReport, err := v.ValidateNoOrphanedRecords(ctx)
	if err != nil {
		return nil, err
	}

	// 2. 验证ID格式一致性
	formatReport, err := v.ValidateIDFormatConsistency(ctx)
	if err != nil {
		return nil, err
	}

	// 3. 合并报告
	report := &ValidationReport{
		ValidatedAt:     time.Now(),
		TotalOrphanedRecords: orphanReport.TotalOrphanedRecords,
		OrphanDetails:   orphanReport.OrphanDetails,
		InconsistentFields: formatReport.InconsistentFields,
		InconsistentRecords: formatReport.InconsistentRecords,
		Errors:          append(orphanReport.Errors, formatReport.Errors...),
	}

	// 获取集合统计信息
	stats, err := v.GetCollectionStats(ctx)
	if err == nil {
		report.CollectionStats = stats
	}

	// 生成总体摘要
	if report.TotalOrphanedRecords == 0 && len(report.InconsistentRecords) == 0 && len(report.Errors) == 0 {
		report.IsValid = true
		report.Summary = "数据验证完成，所有检查通过"
	} else {
		report.IsValid = false
		summaryParts := []string{}
		if report.TotalOrphanedRecords > 0 {
			summaryParts = append(summaryParts, fmt.Sprintf("%d条孤儿记录", report.TotalOrphanedRecords))
		}
		if len(report.InconsistentRecords) > 0 {
			summaryParts = append(summaryParts, fmt.Sprintf("%d条格式不一致", len(report.InconsistentRecords)))
		}
		if len(report.Errors) > 0 {
			summaryParts = append(summaryParts, fmt.Sprintf("%d个错误", len(report.Errors)))
		}
		report.Summary = "数据验证完成，发现问题: " + strings.Join(summaryParts, ", ")
	}

	return report, nil
}

// ValidateReadingProgressOrphans 验证阅读进度孤儿记录
func (v *DataValidator) ValidateReadingProgressOrphans(ctx context.Context) (int64, error) {
	return v.validateOrphansByCollection(
		ctx,
		"reading_progress",
		"user_id", "users",
		"book_id", "books",
	)
}

// ValidateReadingHistoryOrphans 验证阅读历史孤儿记录
func (v *DataValidator) ValidateReadingHistoryOrphans(ctx context.Context) (int64, error) {
	return v.validateOrphansByCollection(
		ctx,
		"reading_histories",
		"user_id", "users",
		"book_id", "books",
		"chapter_id", "chapters",
	)
}

// ValidateBookmarkOrphans 验证书签孤儿记录
func (v *DataValidator) ValidateBookmarkOrphans(ctx context.Context) (int64, error) {
	return v.validateOrphansByCollection(
		ctx,
		"bookmarks",
		"user_id", "users",
		"book_id", "books",
		"chapter_id", "chapters",
	)
}

// ValidateLikeOrphans 验证点赞孤儿记录
func (v *DataValidator) ValidateLikeOrphans(ctx context.Context) (int64, error) {
	collection := v.db.Collection("likes")

	// 聚合查询：查找孤儿记录
	// Like 使用 string 类型的 user_id 和 target_id
	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$lookup", bson.D{
			{"from", "books"},
			{"localField", "target_id"},
			{"foreignField", "_id"},
			{"as", "book"},
		}}},
		bson.D{{"$match", bson.D{
			{"user", bson.D{{"$size", 0}}},
		}}},
		bson.D{{"$count", "orphan_count"}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return 0, err
	}

	if len(results) == 0 {
		return 0, nil
	}

	orphanCount := int64(0)
	if count, ok := results[0]["orphan_count"].(int32); ok {
		orphanCount = int64(count)
	}

	return orphanCount, nil
}

// ValidateNotificationOrphans 验证通知孤儿记录
func (v *DataValidator) ValidateNotificationOrphans(ctx context.Context) (int64, error) {
	collection := v.db.Collection("notifications")

	// 聚合查询：查找孤儿记录
	// Notification 使用 string 类型的 user_id
	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$match", bson.D{
			{"user", bson.D{{"$size", 0}}},
		}}},
		bson.D{{"$count", "orphan_count"}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return 0, err
	}

	if len(results) == 0 {
		return 0, nil
	}

	orphanCount := int64(0)
	if count, ok := results[0]["orphan_count"].(int32); ok {
		orphanCount = int64(count)
	}

	return orphanCount, nil
}

// GetCollectionStats 获取集合统计信息
func (v *DataValidator) GetCollectionStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	collections := []string{
		"users", "books", "chapters",
		"reading_progress", "reading_histories", "bookmarks",
		"likes", "notifications",
	}

	for _, collName := range collections {
		count, err := v.db.Collection(collName).CountDocuments(ctx, bson.M{})
		if err == nil {
			stats[collName] = count
		}
	}

	return stats, nil
}

// ========== Private Helper Methods ==========

// validateOrphansByCollection 通过集合验证孤儿记录
// 支持多个外键字段
func (v *DataValidator) validateOrphansByCollection(
	ctx context.Context,
	collectionName string,
	fieldPairs ...string,
) (int64, error) {
	if len(fieldPairs)%2 != 0 {
		return 0, fmt.Errorf("fieldPairs 必须是偶数个参数")
	}

	collection := v.db.Collection(collectionName)

	// 构建聚合管道来查找孤儿记录
	pipeline := mongo.Pipeline{}

	// 为每个外键字段添加 $lookup
	for i := 0; i < len(fieldPairs); i += 2 {
		localField := fieldPairs[i]
		targetColl := fieldPairs[i+1]

		lookupStage := bson.D{{"$lookup", bson.D{
			{"from", targetColl},
			{"localField", localField},
			{"foreignField", "_id"},
			{"as", fmt.Sprintf("ref_%s", targetColl)},
		}}}
		pipeline = append(pipeline, lookupStage)
	}

	// 添加匹配条件：任意引用为空
	matchConditions := bson.A{}
	for i := 0; i < len(fieldPairs); i += 2 {
		targetColl := fieldPairs[i+1]
		matchConditions = append(matchConditions, bson.D{
			{fmt.Sprintf("ref_%s", targetColl), bson.D{{"$size", 0}}},
		})
	}

	pipeline = append(pipeline, bson.D{{"$match", bson.D{
		{"$or", matchConditions},
	}}})

	pipeline = append(pipeline, bson.D{{"$count", "orphan_count"}})

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return 0, err
	}

	if len(results) == 0 {
		return 0, nil
	}

	orphanCount := int64(0)
	if count, ok := results[0]["orphan_count"].(int32); ok {
		orphanCount = int64(count)
	}

	return orphanCount, nil
}

// validateLikesIDFormat 验证 likes 集合的ID格式
// Like 模型应该使用 string 类型存储外键
func (v *DataValidator) validateLikesIDFormat(ctx context.Context, report *ValidationReport) error {
	collection := v.db.Collection("likes")

	// 查找前100条记录检查格式
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetLimit(100))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return err
	}

	for _, doc := range results {
		// 检查 user_id 是否为 string 类型
		if userID, ok := doc["user_id"].(string); ok {
			// 验证是否为有效的 ObjectID.Hex()
			if _, err := primitive.ObjectIDFromHex(userID); err != nil {
				report.InconsistentRecords = append(report.InconsistentRecords,
					fmt.Sprintf("likes[%s]: user_id '%s' 不是有效的 ObjectID.Hex()",
						doc["_id"], userID))
			}
		} else {
			report.InconsistentRecords = append(report.InconsistentRecords,
				fmt.Sprintf("likes[%s]: user_id 应该是 string 类型", doc["_id"]))
		}

		// 检查 target_id 是否为 string 类型
		if targetID, ok := doc["target_id"].(string); ok {
			if _, err := primitive.ObjectIDFromHex(targetID); err != nil {
				report.InconsistentRecords = append(report.InconsistentRecords,
					fmt.Sprintf("likes[%s]: target_id '%s' 不是有效的 ObjectID.Hex()",
						doc["_id"], targetID))
			}
		} else {
			report.InconsistentRecords = append(report.InconsistentRecords,
				fmt.Sprintf("likes[%s]: target_id 应该是 string 类型", doc["_id"]))
		}
	}

	return nil
}

// validateNotificationsIDFormat 验证 notifications 集合的ID格式
// Notification 模型应该使用 string 类型存储外键
func (v *DataValidator) validateNotificationsIDFormat(ctx context.Context, report *ValidationReport) error {
	collection := v.db.Collection("notifications")

	// 查找前100条记录检查格式
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetLimit(100))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return err
	}

	for _, doc := range results {
		// 检查 user_id 是否为 string 类型
		if userID, ok := doc["user_id"].(string); ok {
			// 验证是否为有效的 ObjectID.Hex()
			if _, err := primitive.ObjectIDFromHex(userID); err != nil {
				report.InconsistentRecords = append(report.InconsistentRecords,
					fmt.Sprintf("notifications[%s]: user_id '%s' 不是有效的 ObjectID.Hex()",
						doc["_id"], userID))
			}
		} else {
			report.InconsistentRecords = append(report.InconsistentRecords,
				fmt.Sprintf("notifications[%s]: user_id 应该是 string 类型", doc["_id"]))
		}
	}

	return nil
}

// validateReadingProgressIDFormat 验证 reading_progress 集合的ID格式
// ReadingProgress 使用 ObjectID 类型（旧格式）
func (v *DataValidator) validateReadingProgressIDFormat(ctx context.Context, report *ValidationReport) error {
	collection := v.db.Collection("reading_progress")

	// 查找前100条记录检查格式
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetLimit(100))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return err
	}

	for _, doc := range results {
		// 检查 user_id, book_id, chapter_id 是否为 ObjectID 类型
		for _, field := range []string{"user_id", "book_id", "chapter_id"} {
			if _, ok := doc[field].(primitive.ObjectID); !ok {
				report.InconsistentRecords = append(report.InconsistentRecords,
					fmt.Sprintf("reading_progress[%s]: %s 应该是 ObjectID 类型", doc["_id"], field))
			}
		}
	}

	return nil
}
