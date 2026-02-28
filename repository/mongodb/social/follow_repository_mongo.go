package social

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/social"
	socialRepo "Qingyu_backend/repository/interfaces/social"
	"Qingyu_backend/repository/mongodb/base"
)

// MongoFollowRepository MongoDB关注仓储实现
type MongoFollowRepository struct {
	*base.BaseMongoRepository                   // 嵌入基类，管理主collection (follows)
	authorFollowsCollection   *mongo.Collection // 作者关注collection独立管理
}

func sanitizeFollowQueryToken(field, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("%s不能为空", field)
	}
	objectID, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return value, nil
	}
	return objectID.Hex(), nil
}

func sanitizeFollowType(followType string) (string, error) {
	switch followType {
	case "user", "author":
		return followType, nil
	default:
		return "", fmt.Errorf("不支持的关注类型")
	}
}

// NewMongoFollowRepository 创建MongoDB关注仓储实例
func NewMongoFollowRepository(db *mongo.Database) socialRepo.FollowRepository {
	authorFollowsCollection := db.Collection("author_follows")

	// 创建索引
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// follows集合索引
	followIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "follower_id", Value: 1},
				{Key: "following_id", Value: 1},
				{Key: "follow_type", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "following_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	// author_follows集合索引
	authorFollowIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "author_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "author_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	// 使用基类的GetCollection()创建follows索引
	_, err := base.NewBaseMongoRepository(db, "follows").GetCollection().Indexes().CreateMany(ctx, followIndexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create follows indexes: %v\n", err)
	}

	_, err = authorFollowsCollection.Indexes().CreateMany(ctx, authorFollowIndexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create author_follows indexes: %v\n", err)
	}

	return &MongoFollowRepository{
		BaseMongoRepository:     base.NewBaseMongoRepository(db, "follows"),
		authorFollowsCollection: authorFollowsCollection,
	}
}

// ========== 用户关注 ==========

// CreateFollow 创建关注关系
func (r *MongoFollowRepository) CreateFollow(ctx context.Context, follow *social.Follow) error {
	if follow.ID.IsZero() {
		follow.ID = primitive.NewObjectID()
	}
	now := time.Now()
	follow.CreatedAt = now
	follow.UpdatedAt = now

	_, err := r.GetCollection().InsertOne(ctx, follow)
	if err != nil {
		return fmt.Errorf("failed to create follow: %w", err)
	}
	return nil
}

// DeleteFollow 删除关注关系
func (r *MongoFollowRepository) DeleteFollow(ctx context.Context, followerID, followingID, followType string) error {
	safeFollowerID, err := sanitizeFollowQueryToken("follower_id", followerID)
	if err != nil {
		return err
	}
	safeFollowingID, err := sanitizeFollowQueryToken("following_id", followingID)
	if err != nil {
		return err
	}
	safeFollowType, err := sanitizeFollowType(followType)
	if err != nil {
		return err
	}
	_, err = r.GetCollection().DeleteOne(ctx, bson.D{
		{Key: "follower_id", Value: safeFollowerID},
		{Key: "following_id", Value: safeFollowingID},
		{Key: "follow_type", Value: safeFollowType},
	})
	if err != nil {
		return fmt.Errorf("failed to delete follow: %w", err)
	}
	return nil
}

// GetFollow 获取关注关系
func (r *MongoFollowRepository) GetFollow(ctx context.Context, followerID, followingID, followType string) (*social.Follow, error) {
	safeFollowerID, err := sanitizeFollowQueryToken("follower_id", followerID)
	if err != nil {
		return nil, err
	}
	safeFollowingID, err := sanitizeFollowQueryToken("following_id", followingID)
	if err != nil {
		return nil, err
	}
	safeFollowType, err := sanitizeFollowType(followType)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"follower_id":  safeFollowerID,
		"following_id": safeFollowingID,
		"follow_type":  safeFollowType,
	}

	var follow social.Follow
	err = r.GetCollection().FindOne(ctx, filter).Decode(&follow)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get follow: %w", err)
	}
	return &follow, nil
}

// IsFollowing 检查是否已关注
func (r *MongoFollowRepository) IsFollowing(ctx context.Context, followerID, followingID, followType string) (bool, error) {
	safeFollowerID, err := sanitizeFollowQueryToken("follower_id", followerID)
	if err != nil {
		return false, err
	}
	safeFollowingID, err := sanitizeFollowQueryToken("following_id", followingID)
	if err != nil {
		return false, err
	}
	safeFollowType, err := sanitizeFollowType(followType)
	if err != nil {
		return false, err
	}
	filter := bson.M{
		"follower_id":  safeFollowerID,
		"following_id": safeFollowingID,
		"follow_type":  safeFollowType,
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check follow status: %w", err)
	}
	return count > 0, nil
}

// GetFollowers 获取粉丝列表
func (r *MongoFollowRepository) GetFollowers(ctx context.Context, userID string, followType string, page, size int) ([]*social.FollowInfo, int64, error) {
	safeUserID, err := sanitizeFollowQueryToken("following_id", userID)
	if err != nil {
		return nil, 0, err
	}
	safeFollowType, err := sanitizeFollowType(followType)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{
		"following_id": safeUserID,
		"follow_type":  safeFollowType,
	}

	// 统计总数
	total, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count followers: %w", err)
	}

	// 分页查询
	skip := (page - 1) * size
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(size)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find followers: %w", err)
	}
	defer cursor.Close(ctx)

	var follows []*social.FollowInfo
	if err = cursor.All(ctx, &follows); err != nil {
		return nil, 0, fmt.Errorf("failed to decode followers: %w", err)
	}

	// 确保返回空切片而不是nil
	if follows == nil {
		follows = []*social.FollowInfo{}
	}

	return follows, total, nil
}

// GetFollowing 获取关注列表
func (r *MongoFollowRepository) GetFollowing(ctx context.Context, userID string, followType string, page, size int) ([]*social.FollowingInfo, int64, error) {
	safeUserID, err := sanitizeFollowQueryToken("follower_id", userID)
	if err != nil {
		return nil, 0, err
	}
	safeFollowType, err := sanitizeFollowType(followType)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{
		"follower_id": safeUserID,
		"follow_type": safeFollowType,
	}

	// 统计总数
	total, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count following: %w", err)
	}

	// 分页查询
	skip := (page - 1) * size
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(size)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find following: %w", err)
	}
	defer cursor.Close(ctx)

	var follows []*social.FollowingInfo
	if err = cursor.All(ctx, &follows); err != nil {
		return nil, 0, fmt.Errorf("failed to decode following: %w", err)
	}

	// 确保返回空切片而不是nil
	if follows == nil {
		follows = []*social.FollowingInfo{}
	}

	return follows, total, nil
}

// UpdateMutualStatus 更新互相关注状态
func (r *MongoFollowRepository) UpdateMutualStatus(ctx context.Context, followerID, followingID, followType string, isMutual bool) error {
	safeFollowerID, err := sanitizeFollowQueryToken("follower_id", followerID)
	if err != nil {
		return err
	}
	safeFollowingID, err := sanitizeFollowQueryToken("following_id", followingID)
	if err != nil {
		return err
	}
	safeFollowType, err := sanitizeFollowType(followType)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"is_mutual":  isMutual,
			"updated_at": time.Now(),
		},
	}

	_, err = r.GetCollection().UpdateOne(ctx, bson.D{
		{Key: "follower_id", Value: safeFollowerID},
		{Key: "following_id", Value: safeFollowingID},
		{Key: "follow_type", Value: safeFollowType},
	}, update)
	if err != nil {
		return fmt.Errorf("failed to update mutual status: %w", err)
	}
	return nil
}

// ========== 作者关注 ==========

// CreateAuthorFollow 创建作者关注
func (r *MongoFollowRepository) CreateAuthorFollow(ctx context.Context, authorFollow *social.AuthorFollow) error {
	if authorFollow.ID.IsZero() {
		authorFollow.ID = primitive.NewObjectID()
	}
	now := time.Now()
	authorFollow.CreatedAt = now
	authorFollow.UpdatedAt = now

	_, err := r.authorFollowsCollection.InsertOne(ctx, authorFollow)
	if err != nil {
		return fmt.Errorf("failed to create author follow: %w", err)
	}
	return nil
}

// DeleteAuthorFollow 删除作者关注
func (r *MongoFollowRepository) DeleteAuthorFollow(ctx context.Context, userID, authorID string) error {
	safeUserID, err := sanitizeFollowQueryToken("user_id", userID)
	if err != nil {
		return err
	}
	safeAuthorID, err := sanitizeFollowQueryToken("author_id", authorID)
	if err != nil {
		return err
	}
	_, err = r.authorFollowsCollection.DeleteOne(ctx, bson.D{
		{Key: "user_id", Value: safeUserID},
		{Key: "author_id", Value: safeAuthorID},
	})
	if err != nil {
		return fmt.Errorf("failed to delete author follow: %w", err)
	}
	return nil
}

// GetAuthorFollow 获取作者关注
func (r *MongoFollowRepository) GetAuthorFollow(ctx context.Context, userID, authorID string) (*social.AuthorFollow, error) {
	safeUserID, err := sanitizeFollowQueryToken("user_id", userID)
	if err != nil {
		return nil, err
	}
	safeAuthorID, err := sanitizeFollowQueryToken("author_id", authorID)
	if err != nil {
		return nil, err
	}
	var authorFollow social.AuthorFollow
	err = r.authorFollowsCollection.FindOne(ctx, bson.D{
		{Key: "user_id", Value: safeUserID},
		{Key: "author_id", Value: safeAuthorID},
	}).Decode(&authorFollow)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get author follow: %w", err)
	}
	return &authorFollow, nil
}

// GetAuthorFollowers 获取作者的粉丝列表
func (r *MongoFollowRepository) GetAuthorFollowers(ctx context.Context, authorID string, page, size int) ([]*social.FollowInfo, int64, error) {
	safeAuthorID, err := sanitizeFollowQueryToken("author_id", authorID)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{
		"author_id": safeAuthorID,
	}

	// 统计总数
	total, err := r.authorFollowsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count author followers: %w", err)
	}

	// 分页查询
	skip := (page - 1) * size
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(size)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.authorFollowsCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find author followers: %w", err)
	}
	defer cursor.Close(ctx)

	var follows []*social.FollowInfo
	if err = cursor.All(ctx, &follows); err != nil {
		return nil, 0, fmt.Errorf("failed to decode author followers: %w", err)
	}

	// 确保返回空切片而不是nil
	if follows == nil {
		follows = []*social.FollowInfo{}
	}

	return follows, total, nil
}

// GetUserFollowingAuthors 获取用户关注的作者列表
func (r *MongoFollowRepository) GetUserFollowingAuthors(ctx context.Context, userID string, page, size int) ([]*social.AuthorFollow, int64, error) {
	safeUserID, err := sanitizeFollowQueryToken("user_id", userID)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{
		"user_id": safeUserID,
	}

	// 统计总数
	total, err := r.authorFollowsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count following authors: %w", err)
	}

	// 分页查询
	skip := (page - 1) * size
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(size)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.authorFollowsCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find following authors: %w", err)
	}
	defer cursor.Close(ctx)

	var authorFollows []*social.AuthorFollow
	if err = cursor.All(ctx, &authorFollows); err != nil {
		return nil, 0, fmt.Errorf("failed to decode following authors: %w", err)
	}

	// 确保返回空切片而不是nil
	if authorFollows == nil {
		authorFollows = []*social.AuthorFollow{}
	}

	return authorFollows, total, nil
}

// ========== 统计 ==========

// GetFollowStats 获取关注统计
func (r *MongoFollowRepository) GetFollowStats(ctx context.Context, userID string) (*social.FollowStats, error) {
	// 统计粉丝数
	followerCount, err := r.GetCollection().CountDocuments(ctx, bson.M{
		"following_id": userID,
		"follow_type":  "user",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to count followers: %w", err)
	}

	// 统计关注数
	followingCount, err := r.GetCollection().CountDocuments(ctx, bson.M{
		"follower_id": userID,
		"follow_type": "user",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to count following: %w", err)
	}

	return &social.FollowStats{
		UserID:         userID,
		FollowerCount:  int(followerCount),
		FollowingCount: int(followingCount),
		UpdatedAt:      time.Now(),
	}, nil
}

// UpdateFollowStats 更新关注统计（简化实现，实际应该在user表中维护）
func (r *MongoFollowRepository) UpdateFollowStats(ctx context.Context, userID string, followerDelta, followingDelta int) error {
	// 简化实现：不在这里更新统计，统计信息实时计算
	// 实际项目中可以在user表中添加follower_count和following_count字段
	return nil
}

// CountFollowers 统计粉丝数
func (r *MongoFollowRepository) CountFollowers(ctx context.Context, userID, followType string) (int64, error) {
	count, err := r.GetCollection().CountDocuments(ctx, bson.M{
		"following_id": userID,
		"follow_type":  followType,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count followers: %w", err)
	}
	return count, nil
}

// CountFollowing 统计关注数
func (r *MongoFollowRepository) CountFollowing(ctx context.Context, userID, followType string) (int64, error) {
	count, err := r.GetCollection().CountDocuments(ctx, bson.M{
		"follower_id": userID,
		"follow_type": followType,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count following: %w", err)
	}
	return count, nil
}

// Health 健康检查
func (r *MongoFollowRepository) Health(ctx context.Context) error {
	// 简单的健康检查：执行一次count操作
	_, err := r.GetCollection().EstimatedDocumentCount(ctx)
	if err != nil {
		return fmt.Errorf("follows collection health check failed: %w", err)
	}

	_, err = r.authorFollowsCollection.EstimatedDocumentCount(ctx)
	if err != nil {
		return fmt.Errorf("author_follows collection health check failed: %w", err)
	}

	return nil
}
