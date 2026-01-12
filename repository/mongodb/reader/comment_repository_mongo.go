package reader

import (
	"Qingyu_backend/models/social"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCommentRepository MongoDB评论仓储实现
type MongoCommentRepository struct {
	collection *mongo.Collection
}

// NewMongoCommentRepository 创建MongoDB评论仓储实例
func NewMongoCommentRepository(db *mongo.Database) *MongoCommentRepository {
	collection := db.Collection("comments")

	// 创建索引
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "book_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "parent_id", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "book_id", Value: 1},
				{Key: "like_count", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "chapter_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		// 索引创建失败不影响启动，只记录日志
		fmt.Printf("Warning: Failed to create comment indexes: %v\n", err)
	}

	return &MongoCommentRepository{
		collection: collection,
	}
}

// Create 创建评论
func (r *MongoCommentRepository) Create(ctx context.Context, comment *social.Comment) error {
	if comment.ID == "" {
		comment.ID = primitive.NewObjectID().Hex()
	}

	if comment.CreatedAt.IsZero() {
		comment.CreatedAt = time.Now()
	}
	comment.UpdatedAt = time.Now()

	// 初始化统计字段（使用mixin字段）
	if comment.ReplyCount == 0 {
		comment.ReplyCount = 0
	}

	// 默认状态为正常
	if comment.State == "" {
		comment.State = social.CommentStateNormal
	}

	_, err := r.collection.InsertOne(ctx, comment)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

// GetByID 根据ID获取评论
func (r *MongoCommentRepository) GetByID(ctx context.Context, id string) (*social.Comment, error) {
	var comment social.Comment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&comment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return &comment, nil
}

// Update 更新评论
func (r *MongoCommentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

// Delete 删除评论（软删除）
func (r *MongoCommentRepository) Delete(ctx context.Context, id string) error {
	// 软删除：标记为已删除状态
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"state":      social.CommentStateDeleted,
				"updated_at": time.Now(),
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

// GetCommentsByBookID 获取书籍的评论列表
func (r *MongoCommentRepository) GetCommentsByBookID(ctx context.Context, bookID string, page, size int) ([]*social.Comment, int64, error) {
	filter := bson.M{
		"target_id":   bookID,
		"target_type": social.CommentTargetTypeBook,
		"state":       social.CommentStateNormal,
		"parent_id":   nil, // 只获取顶级评论
	}

	return r.findComments(ctx, filter, page, size, bson.D{{Key: "created_at", Value: -1}})
}

// GetCommentsByUserID 获取用户的评论列表
func (r *MongoCommentRepository) GetCommentsByUserID(ctx context.Context, userID string, page, size int) ([]*social.Comment, int64, error) {
	filter := bson.M{"author_id": userID}
	return r.findComments(ctx, filter, page, size, bson.D{{Key: "created_at", Value: -1}})
}

// GetRepliesByCommentID 获取评论的回复列表
func (r *MongoCommentRepository) GetRepliesByCommentID(ctx context.Context, commentID string) ([]*social.Comment, error) {
	filter := bson.M{
		"parent_id": commentID,
		"state":     social.CommentStateNormal,
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}
	defer cursor.Close(ctx)

	var replies []*social.Comment
	if err := cursor.All(ctx, &replies); err != nil {
		return nil, fmt.Errorf("failed to decode replies: %w", err)
	}

	return replies, nil
}

// GetCommentsByChapterID 获取章节的评论列表
func (r *MongoCommentRepository) GetCommentsByChapterID(ctx context.Context, chapterID string, page, size int) ([]*social.Comment, int64, error) {
	filter := bson.M{
		"target_id":   chapterID,
		"target_type": social.CommentTargetTypeChapter,
		"state":       social.CommentStateNormal,
		"parent_id":   nil,
	}

	return r.findComments(ctx, filter, page, size, bson.D{{Key: "created_at", Value: -1}})
}

// GetCommentsByBookIDSorted 获取书籍的排序评论列表
func (r *MongoCommentRepository) GetCommentsByBookIDSorted(ctx context.Context, bookID string, sortBy string, page, size int) ([]*social.Comment, int64, error) {
	filter := bson.M{
		"target_id":   bookID,
		"target_type": social.CommentTargetTypeBook,
		"state":       social.CommentStateNormal,
		"parent_id":   nil,
	}

	var sort bson.D
	switch sortBy {
	case social.CommentSortByHot:
		sort = bson.D{{Key: "like_count", Value: -1}, {Key: "created_at", Value: -1}}
	case social.CommentSortByLatest:
		fallthrough
	default:
		sort = bson.D{{Key: "created_at", Value: -1}}
	}

	return r.findComments(ctx, filter, page, size, sort)
}

// UpdateCommentStatus 更新评论审核状态
func (r *MongoCommentRepository) UpdateCommentStatus(ctx context.Context, id, status, reason string) error {
	updates := bson.M{
		"state":      social.CommentState(status),
		"updated_at": time.Now(),
	}

	if reason != "" {
		updates["reject_reason"] = reason
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to update comment status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

// GetPendingComments 获取待审核评论列表
func (r *MongoCommentRepository) GetPendingComments(ctx context.Context, page, size int) ([]*social.Comment, int64, error) {
	filter := bson.M{"state": social.CommentStateNormal} // 或根据业务需求使用pending状态
	return r.findComments(ctx, filter, page, size, bson.D{{Key: "created_at", Value: 1}})
}

// IncrementLikeCount 增加点赞数
func (r *MongoCommentRepository) IncrementLikeCount(ctx context.Context, id string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{"like_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to increment like count: %w", err)
	}

	return nil
}

// DecrementLikeCount 减少点赞数
func (r *MongoCommentRepository) DecrementLikeCount(ctx context.Context, id string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{"like_count": -1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to decrement like count: %w", err)
	}

	return nil
}

// IncrementReplyCount 增加回复数
func (r *MongoCommentRepository) IncrementReplyCount(ctx context.Context, id string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{"thread_size": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to increment reply count: %w", err)
	}

	return nil
}

// DecrementReplyCount 减少回复数
func (r *MongoCommentRepository) DecrementReplyCount(ctx context.Context, id string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{"thread_size": -1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to decrement reply count: %w", err)
	}

	return nil
}

// GetBookRatingStats 获取书籍评分统计
func (r *MongoCommentRepository) GetBookRatingStats(ctx context.Context, bookID string) (map[string]interface{}, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"target_id":   bookID,
			"target_type": social.CommentTargetTypeBook,
			"state":       social.CommentStateNormal,
			"rating":      bson.M{"$gt": 0},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":         nil,
			"total_count": bson.M{"$sum": 1},
			"average":     bson.M{"$avg": "$rating"},
			"five_star":   bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$rating", 5}}, 1, 0}}},
			"four_star":   bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$rating", 4}}, 1, 0}}},
			"three_star":  bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$rating", 3}}, 1, 0}}},
			"two_star":    bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$rating", 2}}, 1, 0}}},
			"one_star":    bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$rating", 1}}, 1, 0}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating stats: %w", err)
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode rating stats: %w", err)
	}

	if len(results) == 0 {
		return map[string]interface{}{
			"total_count": 0,
			"average":     0.0,
			"five_star":   0,
			"four_star":   0,
			"three_star":  0,
			"two_star":    0,
			"one_star":    0,
		}, nil
	}

	return results[0], nil
}

// GetCommentCount 获取书籍评论总数
func (r *MongoCommentRepository) GetCommentCount(ctx context.Context, bookID string) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"target_id":   bookID,
		"target_type": social.CommentTargetTypeBook,
		"state":       social.CommentStateNormal,
		"parent_id":   nil,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count comments: %w", err)
	}

	return count, nil
}

// GetCommentsByIDs 批量获取评论
func (r *MongoCommentRepository) GetCommentsByIDs(ctx context.Context, ids []string) ([]*social.Comment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by IDs: %w", err)
	}
	defer cursor.Close(ctx)

	var comments []*social.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, fmt.Errorf("failed to decode comments: %w", err)
	}

	return comments, nil
}

// DeleteCommentsByBookID 删除书籍的所有评论
func (r *MongoCommentRepository) DeleteCommentsByBookID(ctx context.Context, bookID string) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{
			"target_id":   bookID,
			"target_type": social.CommentTargetTypeBook,
		},
		bson.M{
			"$set": bson.M{
				"state":      social.CommentStateDeleted,
				"updated_at": time.Now(),
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to delete comments by book ID: %w", err)
	}

	return nil
}

// Health 健康检查
func (r *MongoCommentRepository) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}

// findComments 通用的查询评论方法
func (r *MongoCommentRepository) findComments(ctx context.Context, filter bson.M, page, size int, sort bson.D) ([]*social.Comment, int64, error) {
	// 计算总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	// 计算跳过数
	skip := int64((page - 1) * size)

	// 查询
	opts := options.Find().
		SetSort(sort).
		SetSkip(skip).
		SetLimit(int64(size))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find comments: %w", err)
	}
	defer cursor.Close(ctx)

	var comments []*social.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, 0, fmt.Errorf("failed to decode comments: %w", err)
	}

	return comments, total, nil
}
