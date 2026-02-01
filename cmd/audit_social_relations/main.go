// Package main 社交用户视角数据关联审查工具
// 审查 likes, comments, follows, conversations, messages, notifications 集合的数据关联完整性
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuditResult 审查结果结构
type AuditResult struct {
	TotalRecords    int64   `json:"totalRecords"`
	OrphanedRecords int64   `json:"orphanedRecords"`
	Issues          []Issue `json:"issues,omitempty"`
	Score           float64 `json:"score"` // 0-100
	Status          string  `json:"status"` // excellent, good, fair, poor
}

// Issue 数据问题
type Issue struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	Field    string `json:"field"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Message  string `json:"message"`
}

// SocialAuditReport 社交用户视角审查报告
type SocialAuditReport struct {
	AuditTime time.Time `json:"auditTime"`
	Auditor   string    `json:"auditor"`

	// 外键关系完整性
	LikesOrphanedUsers       AuditResult `json:"likesOrphanedUsers"`
	LikesOrphanedTargets     AuditResult `json:"likesOrphanedTargets"`
	CommentsOrphanedAuthors  AuditResult `json:"commentsOrphanedAuthors"`
	CommentsOrphanedTargets  AuditResult `json:"commentsOrphanedTargets"`
	FollowsOrphanedFollowers AuditResult `json:"followsOrphanedFollowers"`
	FollowsOrphanedFollowing AuditResult `json:"followsOrphanedFollowing"`
	ConversationsOrphaned    AuditResult `json:"conversationsOrphanedParticipants"`
	MessagesOrphanedConv      AuditResult `json:"messagesOrphanedConversations"`
	MessagesOrphanedSender    AuditResult `json:"messagesOrphanedSenders"`
	NotificationsOrphaned     AuditResult `json:"notificationsOrphanedUsers"`

	// 业务规则一致性
	DuplicateLikes AuditResult `json:"duplicateLikes"`
	EmptyComments  AuditResult `json:"emptyComments"`
	SelfFollows    AuditResult `json:"selfFollows"`

	// 统计数据准确性
	BooksLikeCountMismatch      AuditResult `json:"booksLikeCountMismatch"`
	UsersFollowersCountMismatch AuditResult `json:"usersFollowersCountMismatch"`
	UsersFollowingCountMismatch AuditResult `json:"usersFollowingCountMismatch"`

	// 整体评分
	OverallScore  float64 `json:"overallScore"`
	OverallStatus string  `json:"overallStatus"`
	TotalIssues   int64   `json:"totalIssues"`
}

var (
	db *mongo.Database
)

func main() {
	fmt.Println("========================================")
	fmt.Println("   社交用户视角数据关联审查工具")
	fmt.Println("========================================")
	fmt.Println()

	// 连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uri := "mongodb://localhost:27017"
	databaseName := "qingyu"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Printf("连接数据库失败: %v\n", err)
		return
	}
	defer client.Disconnect(ctx)

	db = client.Database(databaseName)
	fmt.Printf("已连接数据库: %s\n\n", databaseName)

	// 开始审查
	report := &SocialAuditReport{
		AuditTime: time.Now(),
		Auditor:   "数据审查专家女仆",
	}

	// 外键关系完整性审查
	fmt.Println("【1/9】审查 likes.user_id 外键...")
	report.LikesOrphanedUsers = auditLikesUserID(ctx)

	fmt.Println("【2/9】审查 likes.target_id 外键...")
	report.LikesOrphanedTargets = auditLikesTargetID(ctx)

	fmt.Println("【3/9】审查 comments.author_id 外键...")
	report.CommentsOrphanedAuthors = auditCommentsAuthorID(ctx)

	fmt.Println("【4/9】审查 comments.target_id 外键...")
	report.CommentsOrphanedTargets = auditCommentsTargetID(ctx)

	fmt.Println("【5/9】审查 follows 外键...")
	auditFollowsForeignKeys(ctx, report)

	fmt.Println("【6/9】审查 conversations 参与者...")
	report.ConversationsOrphaned = auditConversationsParticipants(ctx)

	fmt.Println("【7/9】审查 messages 外键...")
	auditMessagesForeignKeys(ctx, report)

	fmt.Println("【8/9】审查 notifications.user_id 外键...")
	report.NotificationsOrphaned = auditNotificationsUserID(ctx)

	// 业务规则一致性审查
	fmt.Println("【9/9】审查业务规则...")
	auditBusinessRules(ctx, report)

	// 统计数据准确性审查
	fmt.Println("【10/10】审查统计数据准确性...")
	auditStatisticsAccuracy(ctx, report)

	// 计算整体评分
	calculateOverallScore(report)

	// 打印报告
	printReport(report)

	// 保存报告到文件
	saveReportToFile(report)
}

// auditLikesUserID 审查 likes.user_id 外键
func auditLikesUserID(ctx context.Context) AuditResult {
	collection := db.Collection("likes")

	// 统计总记录数
	total, _ := collection.CountDocuments(ctx, bson.M{})

	// 查找孤儿记录 (user_id 不在 users._id 中)
	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		{{"$match", bson.D{
			{"user", bson.D{{"$size", 0}}},
		}}},
		{{"$count", "orphaned"}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return AuditResult{
			TotalRecords: total,
			Status:       "error",
			Score:        0,
		}
	}
	defer cursor.Close(ctx)

	var result []bson.M
	cursor.All(ctx, &result)

	orphaned := int64(0)
	if len(result) > 0 {
		// MongoDB 可能返回 int32 或 int64
		if val, ok := result[0]["orphaned"].(int64); ok {
			orphaned = val
		} else if val32, ok := result[0]["orphaned"].(int32); ok {
			orphaned = int64(val32)
		}
	}

	// 获取示例问题
	issues := getOrphanedIssues(ctx, collection, "user_id", "users", 5)

	return calculateResult(total, orphaned, issues)
}

// auditLikesTargetID 审查 likes.target_id 外键
func auditLikesTargetID(ctx context.Context) AuditResult {
	collection := db.Collection("likes")

	total, _ := collection.CountDocuments(ctx, bson.M{})

	// 分别审查 book 和 comment 类型
	orphaned := int64(0)

	// 审查 book 类型
	pipelineBook := mongo.Pipeline{
		{{"$match", bson.D{{"target_type", "book"}}}},
		{{"$lookup", bson.D{
			{"from", "books"},
			{"localField", "target_id"},
			{"foreignField", "_id"},
			{"as", "book"},
		}}},
		{{"$match", bson.D{
			{"book", bson.D{{"$size", 0}}},
		}}},
		{{"$count", "orphaned"}},
	}

	cursor, _ := collection.Aggregate(ctx, pipelineBook)
	if cursor.Next(ctx) {
		var result []bson.M
		cursor.All(ctx, &result)
		if len(result) > 0 {
			orphaned += getInt64FromBSON(result[0], "orphaned")
		}
		cursor.Close(ctx)
	}

	// 审查 comment 类型
	pipelineComment := mongo.Pipeline{
		{{"$match", bson.D{{"target_type", "comment"}}}},
		{{"$lookup", bson.D{
			{"from", "comments"},
			{"localField", "target_id"},
			{"foreignField", "_id"},
			{"as", "comment"},
		}}},
		{{"$match", bson.D{
			{"comment", bson.D{{"$size", 0}}},
		}}},
		{{"$count", "orphaned"}},
	}

	cursor2, _ := collection.Aggregate(ctx, pipelineComment)
	if cursor2.Next(ctx) {
		var result []bson.M
		cursor2.All(ctx, &result)
		if len(result) > 0 {
			orphaned += getInt64FromBSON(result[0], "orphaned")
		}
		cursor2.Close(ctx)
	}

	return calculateResult(total, orphaned, []Issue{})
}

// auditCommentsAuthorID 审查 comments.author_id 外键
func auditCommentsAuthorID(ctx context.Context) AuditResult {
	collection := db.Collection("comments")

	total, _ := collection.CountDocuments(ctx, bson.M{})

	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "author_id"},
			{"foreignField", "_id"},
			{"as", "author"},
		}}},
		{{"$match", bson.D{
			{"author", bson.D{{"$size", 0}}},
		}}},
		{{"$count", "orphaned"}},
	}

	cursor, _ := collection.Aggregate(ctx, pipeline)
	var result []bson.M
	if cursor.Next(ctx) {
		cursor.All(ctx, &result)
	}
	cursor.Close(ctx)

	orphaned := int64(0)
	if len(result) > 0 {
		// MongoDB 可能返回 int32 或 int64
		if val, ok := result[0]["orphaned"].(int64); ok {
			orphaned = val
		} else if val32, ok := result[0]["orphaned"].(int32); ok {
			orphaned = int64(val32)
		}
	}

	issues := getOrphanedIssues(ctx, collection, "author_id", "users", 5)

	return calculateResult(total, orphaned, issues)
}

// auditCommentsTargetID 审查 comments.target_id 外键
func auditCommentsTargetID(ctx context.Context) AuditResult {
	collection := db.Collection("comments")

	total, _ := collection.CountDocuments(ctx, bson.M{})

	orphaned := int64(0)

	// 审查 book 类型
	pipelineBook := mongo.Pipeline{
		{{"$match", bson.D{{"target_type", "book"}}}},
		{{"$lookup", bson.D{
			{"from", "books"},
			{"localField", "target_id"},
			{"foreignField", "_id"},
			{"as", "target"},
		}}},
		{{"$match", bson.D{
			{"target", bson.D{{"$size", 0}}},
		}}},
		{{"$count", "orphaned"}},
	}

	cursor, _ := collection.Aggregate(ctx, pipelineBook)
	if cursor.Next(ctx) {
		var result []bson.M
		cursor.All(ctx, &result)
		if len(result) > 0 {
			orphaned += getInt64FromBSON(result[0], "orphaned")
		}
		cursor.Close(ctx)
	}

	return calculateResult(total, orphaned, []Issue{})
}

// auditFollowsForeignKeys 审查 follows 外键
func auditFollowsForeignKeys(ctx context.Context, report *SocialAuditReport) {
	collection := db.Collection("user_relations")

	total, _ := collection.CountDocuments(ctx, bson.M{})

	// 审查 follower_id
	pipelineFollower := mongo.Pipeline{
		{{"$match", bson.D{{"status", "active"}}}},
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "follower_id"},
			{"foreignField", "_id"},
			{"as", "follower"},
		}}},
		{{"$match", bson.D{
			{"follower", bson.D{{"$size", 0}}},
		}}},
		{{"$count", "orphaned"}},
	}

	cursor, _ := collection.Aggregate(ctx, pipelineFollower)
	var result []bson.M
	if cursor.Next(ctx) {
		cursor.All(ctx, &result)
	}
	cursor.Close(ctx)

	orphanedFollowers := int64(0)
	if len(result) > 0 {
		orphanedFollowers = getInt64FromBSON(result[0], "orphaned")
	}
	report.FollowsOrphanedFollowers = calculateResult(total, orphanedFollowers, []Issue{})

	// 审查 following_id (followee_id)
	pipelineFollowing := mongo.Pipeline{
		{{"$match", bson.D{{"status", "active"}}}},
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "followee_id"},
			{"foreignField", "_id"},
			{"as", "following"},
		}}},
		{{"$match", bson.D{
			{"following", bson.D{{"$size", 0}}},
		}}},
		{{"$count", "orphaned"}},
	}

	cursor2, _ := collection.Aggregate(ctx, pipelineFollowing)
	var resultFollowing []bson.M
	if cursor2.Next(ctx) {
		cursor2.All(ctx, &resultFollowing)
	}
	cursor2.Close(ctx)

	orphanedFollowing := int64(0)
	if len(resultFollowing) > 0 {
		orphanedFollowing = resultFollowing[0]["orphaned"].(int64)
	}
	report.FollowsOrphanedFollowing = calculateResult(total, orphanedFollowing, []Issue{})
}

// auditConversationsParticipants 审查 conversations 参与者
func auditConversationsParticipants(ctx context.Context) AuditResult {
	collection := db.Collection("conversations")

	total, _ := collection.CountDocuments(ctx, bson.M{})

	// 检查参与者中是否有不存在的用户
	pipeline := mongo.Pipeline{
		{{"$unwind", "$participant_ids"}},
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "participant_ids"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		{{"$match", bson.D{
			{"user", bson.D{{"$size", 0}}},
		}}},
		{{"$count", "invalid_participants"}},
	}

	cursor, _ := collection.Aggregate(ctx, pipeline)
	var result []bson.M
	if cursor.Next(ctx) {
		cursor.All(ctx, &result)
	}
	cursor.Close(ctx)

	orphaned := int64(0)
	if len(result) > 0 {
		orphaned = getInt64FromBSON(result[0], "invalid_participants")
	}

	return calculateResult(total, orphaned, []Issue{})
}

// auditMessagesForeignKeys 审查 messages 外键
func auditMessagesForeignKeys(ctx context.Context, report *SocialAuditReport) {
	collection := db.Collection("direct_messages")

	total, _ := collection.CountDocuments(ctx, bson.M{})

	// 审查 conversation_id
	pipelineConv := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "conversations"},
			{"localField", "conversation_id"},
			{"foreignField", "_id"},
			{"as", "conversation"},
		}}},
		{{"$match", bson.D{
			{"conversation", bson.D{{"$size", 0}}},
		}}},
		{{"$count", "orphaned"}},
	}

	cursor, _ := collection.Aggregate(ctx, pipelineConv)
	var result []bson.M
	if cursor.Next(ctx) {
		cursor.All(ctx, &result)
	}
	cursor.Close(ctx)

	orphanedConv := int64(0)
	if len(result) > 0 {
		orphanedConv = getInt64FromBSON(result[0], "orphaned")
	}
	report.MessagesOrphanedConv = calculateResult(total, orphanedConv, []Issue{})

	// 审查 sender_id
	pipelineSender := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "sender_id"},
			{"foreignField", "_id"},
			{"as", "sender"},
		}}},
		{{"$match", bson.D{
			{"sender", bson.D{{"$size", 0}}},
		}}},
		{{"$count", "orphaned"}},
	}

	cursor2, _ := collection.Aggregate(ctx, pipelineSender)
	var result2 []bson.M
	if cursor2.Next(ctx) {
		cursor2.All(ctx, &result2)
	}
	cursor2.Close(ctx)

	orphanedSender := int64(0)
	if len(result2) > 0 {
		orphanedSender = getInt64FromBSON(result2[0], "orphaned")
	}
	report.MessagesOrphanedSender = calculateResult(total, orphanedSender, []Issue{})
}

// auditNotificationsUserID 审查 notifications.user_id 外键
func auditNotificationsUserID(ctx context.Context) AuditResult {
	collection := db.Collection("notifications")

	total, _ := collection.CountDocuments(ctx, bson.M{})

	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		{{"$match", bson.D{
			{"user", bson.D{{"$size", 0}}},
		}}},
		{{"$count", "orphaned"}},
	}

	cursor, _ := collection.Aggregate(ctx, pipeline)
	var result []bson.M
	if cursor.Next(ctx) {
		cursor.All(ctx, &result)
	}
	cursor.Close(ctx)

	orphaned := int64(0)
	if len(result) > 0 {
		// MongoDB 可能返回 int32 或 int64
		if val, ok := result[0]["orphaned"].(int64); ok {
			orphaned = val
		} else if val32, ok := result[0]["orphaned"].(int32); ok {
			orphaned = int64(val32)
		}
	}

	return calculateResult(total, orphaned, []Issue{})
}

// auditBusinessRules 审查业务规则
func auditBusinessRules(ctx context.Context, report *SocialAuditReport) {
	// 检查重复点赞
	collection := db.Collection("likes")
	totalLikes, _ := collection.CountDocuments(ctx, bson.M{})

	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", bson.D{
				{"user_id", "$user_id"},
				{"target_id", "$target_id"},
				{"target_type", "$target_type"},
			}},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$match", bson.D{
			{"count", bson.D{{"$gt", 1}}},
		}}},
		{{"$count", "duplicates"}},
	}

	cursor, _ := collection.Aggregate(ctx, pipeline)
	var result []bson.M
	if cursor.Next(ctx) {
		cursor.All(ctx, &result)
	}
	cursor.Close(ctx)

	duplicates := int64(0)
	if len(result) > 0 {
		duplicates = getInt64FromBSON(result[0], "duplicates")
	}
	report.DuplicateLikes = calculateResult(totalLikes, duplicates, []Issue{})

	// 检查空评论
	collection = db.Collection("comments")
	totalComments, _ := collection.CountDocuments(ctx, bson.M{})

	emptyComments, _ := collection.CountDocuments(ctx, bson.D{
		{"$or", bson.A{
			bson.D{{"content", ""}},
			bson.D{{"content", nil}},
		}},
	})
	report.EmptyComments = calculateResult(totalComments, emptyComments, []Issue{})

	// 检查自己关注自己
	collection = db.Collection("user_relations")
	totalFollows, _ := collection.CountDocuments(ctx, bson.D{{"status", "active"}})

	selfFollows, _ := collection.CountDocuments(ctx, bson.D{
		{"$expr", bson.D{{"$eq", bson.A{"$follower_id", "$followee_id"}}}},
		{"status", "active"},
	})
	report.SelfFollows = calculateResult(totalFollows, selfFollows, []Issue{})
}

// auditStatisticsAccuracy 审查统计数据准确性
func auditStatisticsAccuracy(ctx context.Context, report *SocialAuditReport) {
	// 检查书籍点赞数
	collection := db.Collection("books")
	totalBooks, _ := collection.CountDocuments(ctx, bson.M{})

	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "likes"},
			{"localField", "_id"},
			{"foreignField", "target_id"},
			{"as", "likes"},
		}}},
		{{"$project", bson.D{
			{"_id", 1},
			{"title", 1},
			{"stored_count", "$like_count"},
			{"actual_count", bson.D{{"$size", bson.A{"$likes"}}}},
			{"diff", bson.D{{"$subtract", bson.A{"$like_count", bson.D{{"$size", bson.A{"$likes"}}}}}}},
		}}},
		{{"$match", bson.D{
			{"diff", bson.D{{"$ne", 0}}},
		}}},
		{{"$count", "mismatch"}},
	}

	cursor, _ := collection.Aggregate(ctx, pipeline)
	var result []bson.M
	if cursor.Next(ctx) {
		cursor.All(ctx, &result)
	}
	cursor.Close(ctx)

	mismatch := int64(0)
	if len(result) > 0 {
		mismatch = getInt64FromBSON(result[0], "mismatch")
	}
	report.BooksLikeCountMismatch = calculateResult(totalBooks, mismatch, []Issue{})

	// 检查用户粉丝数
	collection = db.Collection("users")
	totalUsers, _ := collection.CountDocuments(ctx, bson.M{})

	pipelineFollowers := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "user_relations"},
			{"localField", "_id"},
			{"foreignField", "followee_id"},
			{"as", "followers"},
		}}},
		{{"$project", bson.D{
			{"_id", 1},
			{"username", 1},
			{"stored_count", "$followers_count"},
			{"actual_count", bson.D{{"$size", bson.A{"$followers"}}}},
			{"diff", bson.D{{"$subtract", bson.A{"$followers_count", bson.D{{"$size", bson.A{"$followers"}}}}}}},
		}}},
		{{"$match", bson.D{
			{"diff", bson.D{{"$ne", 0}}},
		}}},
		{{"$count", "mismatch"}},
	}

	cursor2, _ := collection.Aggregate(ctx, pipelineFollowers)
	var result2 []bson.M
	if cursor2.Next(ctx) {
		cursor2.All(ctx, &result2)
	}
	cursor2.Close(ctx)

	mismatchFollowers := int64(0)
	if len(result2) > 0 {
		mismatchFollowers = getInt64FromBSON(result2[0], "mismatch")
	}
	report.UsersFollowersCountMismatch = calculateResult(totalUsers, mismatchFollowers, []Issue{})

	// 检查用户关注数
	pipelineFollowing := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "user_relations"},
			{"localField", "_id"},
			{"foreignField", "follower_id"},
			{"as", "following"},
		}}},
		{{"$project", bson.D{
			{"_id", 1},
			{"username", 1},
			{"stored_count", "$following_count"},
			{"actual_count", bson.D{{"$size", bson.A{"$following"}}}},
			{"diff", bson.D{{"$subtract", bson.A{"$following_count", bson.D{{"$size", bson.A{"$following"}}}}}}},
		}}},
		{{"$match", bson.D{
			{"diff", bson.D{{"$ne", 0}}},
		}}},
		{{"$count", "mismatch"}},
	}

	cursor3, _ := collection.Aggregate(ctx, pipelineFollowing)
	var result3 []bson.M
	if cursor3.Next(ctx) {
		cursor3.All(ctx, &result3)
	}
	cursor3.Close(ctx)

	mismatchFollowing := int64(0)
	if len(result3) > 0 {
		mismatchFollowing = getInt64FromBSON(result3[0], "mismatch")
	}
	report.UsersFollowingCountMismatch = calculateResult(totalUsers, mismatchFollowing, []Issue{})
}

// getInt64FromBSON 从 BSON 中安全地获取 int64 值
func getInt64FromBSON(m bson.M, key string) int64 {
	if val, ok := m[key].(int64); ok {
		return val
	}
	if val32, ok := m[key].(int32); ok {
		return int64(val32)
	}
	if valFloat, ok := m[key].(float64); ok {
		return int64(valFloat)
	}
	return 0
}

// getOrphanedIssues 获取孤儿记录示例
func getOrphanedIssues(ctx context.Context, collection *mongo.Collection, field, targetCollection string, limit int64) []Issue {
	issues := []Issue{}

	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", targetCollection},
			{"localField", field},
			{"foreignField", "_id"},
			{"as", "ref"},
		}}},
		{{"$match", bson.D{
			{"ref", bson.D{{"$size", 0}}},
		}}},
		{{"$limit", limit}},
		{{"$project", bson.D{
			{"_id", 1},
			{field, 1},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return issues
	}
	defer cursor.Close(ctx)

	var results []bson.M
	cursor.All(ctx, &results)

	for _, r := range results {
		id := ""
		if oid, ok := r["_id"].(primitive.ObjectID); ok {
			id = oid.Hex()
		} else if str, ok := r["_id"].(string); ok {
			id = str
		}

		fieldValue := ""
		if val, ok := r[field].(string); ok {
			fieldValue = val
		} else if oid, ok := r[field].(primitive.ObjectID); ok {
			fieldValue = oid.Hex()
		}

		issues = append(issues, Issue{
			Type:    "orphaned_record",
			ID:      id,
			Field:   field,
			Actual:  fieldValue,
			Message: fmt.Sprintf("%s 引用的记录不存在于 %s 集合", field, targetCollection),
		})
	}

	return issues
}

// calculateResult 计算审查结果
func calculateResult(total, orphaned int64, issues []Issue) AuditResult {
	score := 100.0
	status := "excellent"

	if total > 0 {
		orphanRate := float64(orphaned) / float64(total) * 100
		score = 100 - orphanRate

		if orphanRate > 0 {
			status = "poor"
		} else if orphanRate == 0 {
			status = "excellent"
		}
	} else {
		status = "empty"
	}

	if score < 0 {
		score = 0
	}

	return AuditResult{
		TotalRecords:    total,
		OrphanedRecords: orphaned,
		Issues:          issues,
		Score:           score,
		Status:          status,
	}
}

// calculateOverallScore 计算整体评分
func calculateOverallScore(report *SocialAuditReport) {
	results := []float64{
		report.LikesOrphanedUsers.Score,
		report.LikesOrphanedTargets.Score,
		report.CommentsOrphanedAuthors.Score,
		report.CommentsOrphanedTargets.Score,
		report.FollowsOrphanedFollowers.Score,
		report.FollowsOrphanedFollowing.Score,
		report.ConversationsOrphaned.Score,
		report.MessagesOrphanedConv.Score,
		report.MessagesOrphanedSender.Score,
		report.NotificationsOrphaned.Score,
		report.DuplicateLikes.Score,
		report.EmptyComments.Score,
		report.SelfFollows.Score,
		report.BooksLikeCountMismatch.Score,
		report.UsersFollowersCountMismatch.Score,
		report.UsersFollowingCountMismatch.Score,
	}

	sum := 0.0
	for _, s := range results {
		sum += s
	}
	report.OverallScore = sum / float64(len(results))

	if report.OverallScore >= 95 {
		report.OverallStatus = "excellent"
	} else if report.OverallScore >= 80 {
		report.OverallStatus = "good"
	} else if report.OverallScore >= 60 {
		report.OverallStatus = "fair"
	} else {
		report.OverallStatus = "poor"
	}

	// 统计总问题数
	report.TotalIssues =
		report.LikesOrphanedUsers.OrphanedRecords +
			report.LikesOrphanedTargets.OrphanedRecords +
			report.CommentsOrphanedAuthors.OrphanedRecords +
			report.CommentsOrphanedTargets.OrphanedRecords +
			report.FollowsOrphanedFollowers.OrphanedRecords +
			report.FollowsOrphanedFollowing.OrphanedRecords +
			report.ConversationsOrphaned.OrphanedRecords +
			report.MessagesOrphanedConv.OrphanedRecords +
			report.MessagesOrphanedSender.OrphanedRecords +
			report.NotificationsOrphaned.OrphanedRecords +
			report.DuplicateLikes.OrphanedRecords +
			report.EmptyComments.OrphanedRecords +
			report.SelfFollows.OrphanedRecords +
			report.BooksLikeCountMismatch.OrphanedRecords +
			report.UsersFollowersCountMismatch.OrphanedRecords +
			report.UsersFollowingCountMismatch.OrphanedRecords
}

// printReport 打印报告
func printReport(report *SocialAuditReport) {
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("        社交用户视角数据关联审查报告")
	fmt.Println("========================================")
	fmt.Printf("审查时间: %s\n", report.AuditTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("审查人: %s\n", report.Auditor)
	fmt.Println()

	fmt.Println("----------------------------------------")
	fmt.Println("1. 外键关系完整性")
	fmt.Println("----------------------------------------")
	fmt.Printf("likes.user_id 孤儿记录: %d / %d (%.1f%%)\n",
		report.LikesOrphanedUsers.OrphanedRecords,
		report.LikesOrphanedUsers.TotalRecords,
		100-report.LikesOrphanedUsers.Score)
	fmt.Printf("likes.target_id 孤儿记录: %d / %d (%.1f%%)\n",
		report.LikesOrphanedTargets.OrphanedRecords,
		report.LikesOrphanedTargets.TotalRecords,
		100-report.LikesOrphanedTargets.Score)
	fmt.Printf("comments.author_id 孤儿记录: %d / %d (%.1f%%)\n",
		report.CommentsOrphanedAuthors.OrphanedRecords,
		report.CommentsOrphanedAuthors.TotalRecords,
		100-report.CommentsOrphanedAuthors.Score)
	fmt.Printf("comments.target_id 孤儿记录: %d / %d (%.1f%%)\n",
		report.CommentsOrphanedTargets.OrphanedRecords,
		report.CommentsOrphanedTargets.TotalRecords,
		100-report.CommentsOrphanedTargets.Score)
	fmt.Printf("follows.follower_id 孤儿记录: %d / %d (%.1f%%)\n",
		report.FollowsOrphanedFollowers.OrphanedRecords,
		report.FollowsOrphanedFollowers.TotalRecords,
		100-report.FollowsOrphanedFollowers.Score)
	fmt.Printf("follows.following_id 孤儿记录: %d / %d (%.1f%%)\n",
		report.FollowsOrphanedFollowing.OrphanedRecords,
		report.FollowsOrphanedFollowing.TotalRecords,
		100-report.FollowsOrphanedFollowing.Score)
	fmt.Printf("conversations 参与者问题: %d / %d (%.1f%%)\n",
		report.ConversationsOrphaned.OrphanedRecords,
		report.ConversationsOrphaned.TotalRecords,
		100-report.ConversationsOrphaned.Score)
	fmt.Printf("messages.conversation_id 孤儿记录: %d / %d (%.1f%%)\n",
		report.MessagesOrphanedConv.OrphanedRecords,
		report.MessagesOrphanedConv.TotalRecords,
		100-report.MessagesOrphanedConv.Score)
	fmt.Printf("messages.sender_id 孤儿记录: %d / %d (%.1f%%)\n",
		report.MessagesOrphanedSender.OrphanedRecords,
		report.MessagesOrphanedSender.TotalRecords,
		100-report.MessagesOrphanedSender.Score)
	fmt.Printf("notifications.user_id 孤儿记录: %d / %d (%.1f%%)\n",
		report.NotificationsOrphaned.OrphanedRecords,
		report.NotificationsOrphaned.TotalRecords,
		100-report.NotificationsOrphaned.Score)

	fmt.Println()
	fmt.Println("----------------------------------------")
	fmt.Println("2. 业务规则一致性")
	fmt.Println("----------------------------------------")
	fmt.Printf("重复点赞: %d / %d (%.1f%%)\n",
		report.DuplicateLikes.OrphanedRecords,
		report.DuplicateLikes.TotalRecords,
		100-report.DuplicateLikes.Score)
	fmt.Printf("空评论: %d / %d (%.1f%%)\n",
		report.EmptyComments.OrphanedRecords,
		report.EmptyComments.TotalRecords,
		100-report.EmptyComments.Score)
	fmt.Printf("自己关注自己: %d / %d (%.1f%%)\n",
		report.SelfFollows.OrphanedRecords,
		report.SelfFollows.TotalRecords,
		100-report.SelfFollows.Score)

	fmt.Println()
	fmt.Println("----------------------------------------")
	fmt.Println("3. 统计数据准确性")
	fmt.Println("----------------------------------------")
	fmt.Printf("书籍点赞数不一致: %d / %d (%.1f%%)\n",
		report.BooksLikeCountMismatch.OrphanedRecords,
		report.BooksLikeCountMismatch.TotalRecords,
		100-report.BooksLikeCountMismatch.Score)
	fmt.Printf("用户粉丝数不一致: %d / %d (%.1f%%)\n",
		report.UsersFollowersCountMismatch.OrphanedRecords,
		report.UsersFollowersCountMismatch.TotalRecords,
		100-report.UsersFollowersCountMismatch.Score)
	fmt.Printf("用户关注数不一致: %d / %d (%.1f%%)\n",
		report.UsersFollowingCountMismatch.OrphanedRecords,
		report.UsersFollowingCountMismatch.TotalRecords,
		100-report.UsersFollowingCountMismatch.Score)

	fmt.Println()
	fmt.Println("----------------------------------------")
	fmt.Println("4. 数据质量评估")
	fmt.Println("----------------------------------------")
	fmt.Printf("整体评分: %.1f/100\n", report.OverallScore)
	fmt.Printf("整体状态: %s\n", report.OverallStatus)
	fmt.Printf("总问题数: %d\n", report.TotalIssues)

	fmt.Println("========================================")
}

// saveReportToFile 保存报告到文件
func saveReportToFile(report *SocialAuditReport) {
	// 生成 Markdown 报告
	md := fmt.Sprintf(`# 社交用户视角数据关联审查报告

## 审查时间
%s

## 审查人
%s

## 外键关系完整性

### 点赞 (likes)
- **user_id 孤儿记录**: %d / %d (%.1f%%)
- **target_id 孤儿记录**: %d / %d (%.1f%%)

### 评论 (comments)
- **author_id 孤儿记录**: %d / %d (%.1f%%)
- **target_id 孤儿记录**: %d / %d (%.1f%%)

### 关注 (follows/user_relations)
- **follower_id 孤儿记录**: %d / %d (%.1f%%)
- **following_id 孤儿记录**: %d / %d (%.1f%%)

### 对话 (conversations)
- **参与者问题**: %d / %d (%.1f%%)

### 消息 (direct_messages)
- **conversation_id 孤儿记录**: %d / %d (%.1f%%)
- **sender_id 孤儿记录**: %d / %d (%.1f%%)

### 通知 (notifications)
- **user_id 孤儿记录**: %d / %d (%.1f%%)

## 业务规则一致性

- **重复点赞**: %d / %d (%.1f%%)
- **空评论**: %d / %d (%.1f%%)
- **自己关注自己**: %d / %d (%.1f%%)

## 统计数据准确性

- **书籍点赞数不一致**: %d / %d (%.1f%%)
- **用户粉丝数不一致**: %d / %d (%.1f%%)
- **用户关注数不一致**: %d / %d (%.1f%%)

## 数据质量评估

### 整体评分
**%.1f/100**

### 整体状态
**%s**

### 总问题数
**%d**

### 数据质量评估标准
- **excellent (90-100分)**: 数据质量优秀，问题极少
- **good (80-89分)**: 数据质量良好，有个别问题
- **fair (60-79分)**: 数据质量一般，存在较多问题
- **poor (0-59分)**: 数据质量差，存在严重问题

### 主要问题汇总
%s

### 修复建议
%s

---
*本报告由数据审查专家女仆自动生成*
`,
		report.AuditTime.Format("2006-01-02 15:04:05"),
		report.Auditor,
		report.LikesOrphanedUsers.OrphanedRecords, report.LikesOrphanedUsers.TotalRecords, 100-report.LikesOrphanedUsers.Score,
		report.LikesOrphanedTargets.OrphanedRecords, report.LikesOrphanedTargets.TotalRecords, 100-report.LikesOrphanedTargets.Score,
		report.CommentsOrphanedAuthors.OrphanedRecords, report.CommentsOrphanedAuthors.TotalRecords, 100-report.CommentsOrphanedAuthors.Score,
		report.CommentsOrphanedTargets.OrphanedRecords, report.CommentsOrphanedTargets.TotalRecords, 100-report.CommentsOrphanedTargets.Score,
		report.FollowsOrphanedFollowers.OrphanedRecords, report.FollowsOrphanedFollowers.TotalRecords, 100-report.FollowsOrphanedFollowers.Score,
		report.FollowsOrphanedFollowing.OrphanedRecords, report.FollowsOrphanedFollowing.TotalRecords, 100-report.FollowsOrphanedFollowing.Score,
		report.ConversationsOrphaned.OrphanedRecords, report.ConversationsOrphaned.TotalRecords, 100-report.ConversationsOrphaned.Score,
		report.MessagesOrphanedConv.OrphanedRecords, report.MessagesOrphanedConv.TotalRecords, 100-report.MessagesOrphanedConv.Score,
		report.MessagesOrphanedSender.OrphanedRecords, report.MessagesOrphanedSender.TotalRecords, 100-report.MessagesOrphanedSender.Score,
		report.NotificationsOrphaned.OrphanedRecords, report.NotificationsOrphaned.TotalRecords, 100-report.NotificationsOrphaned.Score,
		report.DuplicateLikes.OrphanedRecords, report.DuplicateLikes.TotalRecords, 100-report.DuplicateLikes.Score,
		report.EmptyComments.OrphanedRecords, report.EmptyComments.TotalRecords, 100-report.EmptyComments.Score,
		report.SelfFollows.OrphanedRecords, report.SelfFollows.TotalRecords, 100-report.SelfFollows.Score,
		report.BooksLikeCountMismatch.OrphanedRecords, report.BooksLikeCountMismatch.TotalRecords, 100-report.BooksLikeCountMismatch.Score,
		report.UsersFollowersCountMismatch.OrphanedRecords, report.UsersFollowersCountMismatch.TotalRecords, 100-report.UsersFollowersCountMismatch.Score,
		report.UsersFollowingCountMismatch.OrphanedRecords, report.UsersFollowingCountMismatch.TotalRecords, 100-report.UsersFollowingCountMismatch.Score,
		report.OverallScore,
		report.OverallStatus,
		report.TotalIssues,
		generateIssueSummary(report),
		generateFixSuggestions(report),
	)

	// 保存到文件
	reportPath := "E:\\Github\\Qingyu\\docs\\reports\\2026-02-01-data-relation-audit-social.md"
	err := os.WriteFile(reportPath, []byte(md), 0644)
	if err != nil {
		fmt.Printf("\n保存报告失败: %v\n", err)
		return
	}

	fmt.Printf("\n报告已保存到: %s\n", reportPath)
	fmt.Printf("报告长度: %d 字符\n", len(md))
}

// generateIssueSummary 生成问题摘要
func generateIssueSummary(report *SocialAuditReport) string {
	if report.TotalIssues == 0 {
		return "未发现数据关联问题，数据质量优秀。"
	}

	var summary string
	summary += fmt.Sprintf("- 发现 %d 个数据关联问题\n", report.TotalIssues)

	// 统计主要问题类型
	if report.LikesOrphanedUsers.OrphanedRecords > 0 {
		summary += fmt.Sprintf("- likes.user_id 有 %d 条孤儿记录\n", report.LikesOrphanedUsers.OrphanedRecords)
	}
	if report.LikesOrphanedTargets.OrphanedRecords > 0 {
		summary += fmt.Sprintf("- likes.target_id 有 %d 条孤儿记录\n", report.LikesOrphanedTargets.OrphanedRecords)
	}
	if report.CommentsOrphanedAuthors.OrphanedRecords > 0 {
		summary += fmt.Sprintf("- comments.author_id 有 %d 条孤儿记录\n", report.CommentsOrphanedAuthors.OrphanedRecords)
	}
	if report.DuplicateLikes.OrphanedRecords > 0 {
		summary += fmt.Sprintf("- 发现 %d 条重复点赞记录\n", report.DuplicateLikes.OrphanedRecords)
	}
	if report.EmptyComments.OrphanedRecords > 0 {
		summary += fmt.Sprintf("- 发现 %d 条空评论记录\n", report.EmptyComments.OrphanedRecords)
	}

	return summary
}

// generateFixSuggestions 生成修复建议
func generateFixSuggestions(report *SocialAuditReport) string {
	if report.TotalIssues == 0 {
		return "数据质量优秀，建议继续保持当前的数据管理规范。"
	}

	var suggestions string

	suggestions += "### 立即处理\n"
	if report.LikesOrphanedUsers.OrphanedRecords > 0 || report.LikesOrphanedTargets.OrphanedRecords > 0 {
		suggestions += "- 删除或修复 likes 表中的孤儿记录\n"
	}
	if report.CommentsOrphanedAuthors.OrphanedRecords > 0 {
		suggestions += "- 删除或修复 comments 表中的孤儿记录\n"
	}
	if report.DuplicateLikes.OrphanedRecords > 0 {
		suggestions += "- 清理重复的点赞记录，保留最早的一条\n"
	}
	if report.EmptyComments.OrphanedRecords > 0 {
		suggestions += "- 删除空评论记录\n"
	}
	if report.SelfFollows.OrphanedRecords > 0 {
		suggestions += "- 删除自己关注自己的记录\n"
	}

	suggestions += "\n### 预防措施\n"
	suggestions += "- 在 Seeder 工具中添加 ID 类型统一验证\n"
	suggestions += "- 在 API 层添加外键关系验证\n"
	suggestions += "- 定期执行数据关联完整性检查\n"
	suggestions += "- 添加数据库约束和索引\n"

	return suggestions
}
