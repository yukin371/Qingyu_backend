package social

import (
	"context"
	"fmt"
	"time"

	pkgtransaction "Qingyu_backend/pkg/transaction"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/social"
)

// MongoPostRepository MongoDB动态仓储实现
type MongoPostRepository struct {
	postCollection *mongo.Collection
	likeCollection *mongo.Collection
}

// NewMongoPostRepository 创建MongoDB动态仓储实例
func NewMongoPostRepository(db *mongo.Database) *MongoPostRepository {
	postCollection := db.Collection("posts")
	likeCollection := db.Collection("post_likes")

	// 创建索引
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Posts 索引
	postIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "topics", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "like_count", Value: -1},
				{Key: "created_at", Value: -1},
			},
		},
	}

	_, err := postCollection.Indexes().CreateMany(ctx, postIndexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create post indexes: %v\n", err)
	}

	// PostLikes 索引
	likeIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "post_id", Value: 1},
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "post_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
	}

	_, err = likeCollection.Indexes().CreateMany(ctx, likeIndexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create post_like indexes: %v\n", err)
	}

	return &MongoPostRepository{
		postCollection: postCollection,
		likeCollection: likeCollection,
	}
}

// ========== 动态管理 ==========

// Create 创建动态
func (r *MongoPostRepository) Create(ctx context.Context, post *social.Post) error {
	if post.ID.IsZero() {
		post.ID = primitive.NewObjectID()
	}

	if post.CreatedAt.IsZero() {
		post.CreatedAt = time.Now()
	}
	post.UpdatedAt = time.Now()

	// 初始化统计字段
	if post.LikeCount == 0 {
		post.LikeCount = 0
	}
	if post.CommentCount == 0 {
		post.CommentCount = 0
	}
	if post.ShareCount == 0 {
		post.ShareCount = 0
	}

	_, err := r.postCollection.InsertOne(ctx, post)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

// GetByID 根据ID获取动态
func (r *MongoPostRepository) GetByID(ctx context.Context, id string) (*social.Post, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid post id")
	}

	var post social.Post
	err = r.postCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

// List 获取动态列表
func (r *MongoPostRepository) List(ctx context.Context, page, size int, topic, sort string) ([]*social.Post, int64, error) {
	filter := bson.M{}

	// 话题筛选
	if topic != "" {
		filter["topics"] = topic
	}

	// 计算总数
	total, err := r.postCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// 排序
	var sortOpts bson.D
	switch sort {
	case "hottest":
		// 按点赞数倒序，点赞数相同按时间倒序
		sortOpts = bson.D{
			{Key: "like_count", Value: -1},
			{Key: "created_at", Value: -1},
		}
	default:
		// 默认按时间倒序
		sortOpts = bson.D{{Key: "created_at", Value: -1}}
	}

	// 分页
	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSort(sortOpts).
		SetSkip(skip).
		SetLimit(int64(size))

	cursor, err := r.postCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find posts: %w", err)
	}
	defer cursor.Close(ctx)

	var posts []*social.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, 0, fmt.Errorf("failed to decode posts: %w", err)
	}

	return posts, total, nil
}

// Update 更新动态
func (r *MongoPostRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid post id: %w", err)
	}

	updates["updated_at"] = time.Now()

	result, err := r.postCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}

// Delete 删除动态
func (r *MongoPostRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid post id: %w", err)
	}

	result, err := r.postCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}

// GetByUser 获取用户发布的动态
func (r *MongoPostRepository) GetByUser(ctx context.Context, userID string, page, size int) ([]*social.Post, int64, error) {
	filter := bson.M{"user_id": userID}

	// 计算总数
	total, err := r.postCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// 分页
	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(size))

	cursor, err := r.postCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find posts: %w", err)
	}
	defer cursor.Close(ctx)

	var posts []*social.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, 0, fmt.Errorf("failed to decode posts: %w", err)
	}

	return posts, total, nil
}

// ========== 动态点赞 ==========

// Like 点赞动态
func (r *MongoPostRepository) Like(ctx context.Context, postID, userID string) error {
	postLike := &social.PostLike{
		ID:        primitive.NewObjectID(),
		PostID:    postID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	_, err := r.likeCollection.InsertOne(ctx, postLike)
	if err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	return nil
}

// Unlike 取消点赞动态
func (r *MongoPostRepository) Unlike(ctx context.Context, postID, userID string) error {
	_, err := r.likeCollection.DeleteOne(ctx, bson.M{
		"post_id": postID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}

	return nil
}

// IsLiked 检查是否已点赞
func (r *MongoPostRepository) IsLiked(ctx context.Context, postID, userID string) (bool, error) {
	count, err := r.likeCollection.CountDocuments(ctx, bson.M{
		"post_id": postID,
		"user_id": userID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check like status: %w", err)
	}

	return count > 0, nil
}

// GetLikedPostIDs 批量获取动态的点赞状态
func (r *MongoPostRepository) GetLikedPostIDs(ctx context.Context, userID string, postIDs []string) (map[string]bool, error) {
	if len(postIDs) == 0 {
		return make(map[string]bool), nil
	}

	filter := bson.M{
		"post_id": bson.M{"$in": postIDs},
		"user_id": userID,
	}

	cursor, err := r.likeCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find likes: %w", err)
	}
	defer cursor.Close(ctx)

	result := make(map[string]bool)
	for cursor.Next(ctx) {
		var like social.PostLike
		if err := cursor.Decode(&like); err != nil {
			continue
		}
		result[like.PostID] = true
	}

	return result, nil
}

// IncrementLikeCount 增加点赞数
func (r *MongoPostRepository) IncrementLikeCount(ctx context.Context, postID string) error {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return fmt.Errorf("invalid post id: %w", err)
	}

	_, err = r.postCollection.UpdateOne(
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
func (r *MongoPostRepository) DecrementLikeCount(ctx context.Context, postID string) error {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return fmt.Errorf("invalid post id: %w", err)
	}

	_, err = r.postCollection.UpdateOne(
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

// IncrementCommentCount 增加评论数
func (r *MongoPostRepository) IncrementCommentCount(ctx context.Context, postID string) error {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return fmt.Errorf("invalid post id: %w", err)
	}

	_, err = r.postCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$inc": bson.M{"comment_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to increment comment count: %w", err)
	}

	return nil
}

// DecrementCommentCount 减少评论数
func (r *MongoPostRepository) DecrementCommentCount(ctx context.Context, postID string) error {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return fmt.Errorf("invalid post id: %w", err)
	}

	_, err = r.postCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$inc": bson.M{"comment_count": -1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to decrement comment count: %w", err)
	}

	return nil
}

// RunInTransaction 在事务中执行操作
func (r *MongoPostRepository) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	err := pkgtransaction.RunMongoTransaction(ctx, r.postCollection.Database().Client(), fn)
	if err != nil {
		return fmt.Errorf("post transaction failed: %w", err)
	}

	return nil
}

// Health 健康检查
func (r *MongoPostRepository) Health(ctx context.Context) error {
	return r.postCollection.Database().Client().Ping(ctx, nil)
}
