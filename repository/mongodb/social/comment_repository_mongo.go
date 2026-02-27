package social

import (
	"Qingyu_backend/models/social"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/repository/mongodb/base"
)

// MongoCommentRepository MongoDB评论仓储实现
type MongoCommentRepository struct {
	*base.BaseMongoRepository
}

func sanitizeSocialCommentQueryToken(field, value string) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return "", fmt.Errorf("%s格式不合法", field)
	}
	return objectID.Hex(), nil
}

func sanitizeSocialCommentFilter(filter bson.M) (bson.M, error) {
	safeFilter := make(bson.M, len(filter))
	for key, value := range filter {
		switch key {
		case "_id":
			switch idValue := value.(type) {
			case primitive.ObjectID:
				safeFilter[key] = idValue
			case string:
				objectID, err := primitive.ObjectIDFromHex(idValue)
				if err != nil {
					return nil, fmt.Errorf("invalid id: %w", err)
				}
				safeFilter[key] = objectID
			case bson.M:
				inRaw, ok := idValue["$in"]
				if !ok {
					return nil, fmt.Errorf("unsupported _id operator")
				}
				switch inList := inRaw.(type) {
				case []primitive.ObjectID:
					safeFilter[key] = bson.M{"$in": inList}
				case []string:
					objectIDs := make([]primitive.ObjectID, 0, len(inList))
					for _, id := range inList {
						objectID, err := primitive.ObjectIDFromHex(id)
						if err != nil {
							return nil, fmt.Errorf("invalid id: %w", err)
						}
						objectIDs = append(objectIDs, objectID)
					}
					safeFilter[key] = bson.M{"$in": objectIDs}
				default:
					return nil, fmt.Errorf("unsupported _id $in value")
				}
			default:
				return nil, fmt.Errorf("unsupported _id filter value")
			}
		case "target_id", "author_id":
			valueStr, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid %s filter type", key)
			}
			objectID, err := primitive.ObjectIDFromHex(valueStr)
			if err != nil {
				return nil, fmt.Errorf("invalid id: %w", err)
			}
			safeFilter[key] = objectID.Hex()
		case "parent_id":
			if value == nil {
				safeFilter[key] = nil
				continue
			}
			valueStr, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid parent_id filter type")
			}
			objectID, err := primitive.ObjectIDFromHex(valueStr)
			if err != nil {
				return nil, fmt.Errorf("invalid id: %w", err)
			}
			safeFilter[key] = objectID.Hex()
		case "target_type":
			valueStr, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid target_type filter type")
			}
			switch social.CommentTargetType(valueStr) {
			case social.CommentTargetTypeBook, social.CommentTargetTypeChapter, social.CommentTargetTypeArticle, social.CommentTargetTypeAnnouncement, social.CommentTargetTypeProject:
				safeFilter[key] = valueStr
			default:
				return nil, fmt.Errorf("invalid target_type")
			}
		case "state":
			valueStr, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid state filter type")
			}
			switch social.CommentState(valueStr) {
			case social.CommentStateNormal, social.CommentStateHidden, social.CommentStateDeleted, social.CommentStateRejected:
				safeFilter[key] = valueStr
			default:
				return nil, fmt.Errorf("invalid state")
			}
		case "rating":
			safeFilter[key] = value
		default:
			return nil, fmt.Errorf("unsupported filter key: %s", key)
		}
	}
	return safeFilter, nil
}

// NewMongoCommentRepository 创建MongoDB评论仓储实例
func NewMongoCommentRepository(db *mongo.Database) *MongoCommentRepository {
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

	collection := db.Collection("comments")
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		// 索引创建失败不影响启动，只记录日志
		fmt.Printf("Warning: Failed to create comment indexes: %v\n", err)
	}

	return &MongoCommentRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "comments"),
	}
}

// Create 创建评论
func (r *MongoCommentRepository) Create(ctx context.Context, comment *social.Comment) error {
	if comment.ID.IsZero() {
		comment.ID = primitive.NewObjectID()
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

	_, err := r.GetCollection().InsertOne(ctx, comment)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

// GetByID 根据ID获取评论
func (r *MongoCommentRepository) GetByID(ctx context.Context, id string) (*social.Comment, error) {
	objectID, err := r.ParseID(id)
	if err != nil {
		return nil, fmt.Errorf("comment not found")
	}

	var comment social.Comment
	err = r.GetCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&comment)
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
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	updates["updated_at"] = time.Now()

	result, err := r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
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
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	// 软删除：标记为已删除状态
	result, err := r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
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

// Exists 检查评论是否存在（用于外键验证）
func (r *MongoCommentRepository) Exists(ctx context.Context, id string) (bool, error) {
	objectID, err := r.ParseID(id)
	if err != nil {
		return false, err
	}

	count, err := r.GetCollection().CountDocuments(ctx, bson.M{
		"_id":   objectID,
		"state": social.CommentStateNormal, // 只统计正常状态的评论
	})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetCommentsByBookID 获取书籍的评论列表
func (r *MongoCommentRepository) GetCommentsByBookID(ctx context.Context, bookID string, page, size int) ([]*social.Comment, int64, error) {
	safeBookID, err := sanitizeSocialCommentQueryToken("target_id", bookID)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{
		"target_id":   safeBookID,
		"target_type": social.CommentTargetTypeBook,
		"state":       social.CommentStateNormal,
		"parent_id":   nil, // 只获取顶级评论
	}
	total, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}
	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(size))
	cursor, err := r.GetCollection().Find(ctx, filter, opts)
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

// GetCommentsByUserID 获取用户的评论列表
func (r *MongoCommentRepository) GetCommentsByUserID(ctx context.Context, userID string, page, size int) ([]*social.Comment, int64, error) {
	safeUserID, err := sanitizeSocialCommentQueryToken("author_id", userID)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{"author_id": safeUserID}
	total, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}
	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(size))
	cursor, err := r.GetCollection().Find(ctx, filter, opts)
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

// GetRepliesByCommentID 获取评论的回复列表
func (r *MongoCommentRepository) GetRepliesByCommentID(ctx context.Context, commentID string) ([]*social.Comment, error) {
	parentID, err := r.ParseID(commentID)
	if err != nil {
		return nil, fmt.Errorf("invalid comment id: %w", err)
	}

	filter := bson.M{
		"parent_id": parentID.Hex(),
		"state":     social.CommentStateNormal,
	}

	cursor, err := r.GetCollection().Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
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
	safeChapterID, err := sanitizeSocialCommentQueryToken("target_id", chapterID)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{
		"target_id":   safeChapterID,
		"target_type": social.CommentTargetTypeChapter,
		"state":       social.CommentStateNormal,
		"parent_id":   nil,
	}
	total, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}
	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(size))
	cursor, err := r.GetCollection().Find(ctx, filter, opts)
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

// GetCommentsByBookIDSorted 获取书籍的排序评论列表
func (r *MongoCommentRepository) GetCommentsByBookIDSorted(ctx context.Context, bookID string, sortBy string, page, size int) ([]*social.Comment, int64, error) {
	safeBookID, err := sanitizeSocialCommentQueryToken("target_id", bookID)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{
		"target_id":   safeBookID,
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
	total, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}
	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSort(sort).
		SetSkip(skip).
		SetLimit(int64(size))
	cursor, err := r.GetCollection().Find(ctx, filter, opts)
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

// UpdateCommentStatus 更新评论审核状态
func (r *MongoCommentRepository) UpdateCommentStatus(ctx context.Context, id, status, reason string) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	updates := bson.M{
		"state":      social.CommentState(status),
		"updated_at": time.Now(),
	}

	if reason != "" {
		updates["reject_reason"] = reason
	}

	result, err := r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
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
	total, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}
	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: 1}}).
		SetSkip(skip).
		SetLimit(int64(size))
	cursor, err := r.GetCollection().Find(ctx, filter, opts)
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

// IncrementLikeCount 增加点赞数
func (r *MongoCommentRepository) IncrementLikeCount(ctx context.Context, id string) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	_, err = r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
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
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	_, err = r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
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
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	_, err = r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
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
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	_, err = r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
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
	safeBookID, err := sanitizeSocialCommentQueryToken("target_id", bookID)
	if err != nil {
		return nil, err
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"target_id":   safeBookID,
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

	cursor, err := r.GetCollection().Aggregate(ctx, pipeline)
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
	safeBookID, err := sanitizeSocialCommentQueryToken("target_id", bookID)
	if err != nil {
		return 0, err
	}
	count, err := r.GetCollection().CountDocuments(ctx, bson.M{
		"target_id":   safeBookID,
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
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objectID, err := r.ParseID(id)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %w", err)
		}
		objectIDs = append(objectIDs, objectID)
	}

	cursor, err := r.GetCollection().Find(ctx, bson.M{"_id": bson.M{"$in": objectIDs}})
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
	safeBookID, err := sanitizeSocialCommentQueryToken("target_id", bookID)
	if err != nil {
		return err
	}
	_, err = r.GetCollection().UpdateMany(
		ctx,
		bson.M{
			"target_id":   safeBookID,
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
	return r.GetDB().Client().Ping(ctx, nil)
}
