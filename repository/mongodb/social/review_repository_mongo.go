package social

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/social"
)

// MongoReviewRepository MongoDB书评仓储实现
type MongoReviewRepository struct {
	reviewCollection *mongo.Collection
	likeCollection   *mongo.Collection
}

var reviewSafeQueryTokenPattern = regexp.MustCompile(`^[A-Za-z0-9:_-]{1,128}$`)

func sanitizeReviewQueryToken(field, value string) (string, error) {
	if !reviewSafeQueryTokenPattern.MatchString(value) {
		return "", fmt.Errorf("%s格式不合法", field)
	}
	return value, nil
}

func sanitizeReviewFilter(filter bson.M) (bson.M, error) {
	safeFilter := make(bson.M, len(filter))
	for key, value := range filter {
		switch key {
		case "book_id", "user_id":
			valueStr, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid %s filter type", key)
			}
			objectID, err := primitive.ObjectIDFromHex(valueStr)
			if err != nil {
				return nil, fmt.Errorf("invalid id: %w", err)
			}
			safeFilter[key] = objectID.Hex()
		case "is_public":
			boolValue, ok := value.(bool)
			if !ok {
				return nil, fmt.Errorf("invalid is_public filter type")
			}
			safeFilter[key] = boolValue
		case "rating":
			safeFilter[key] = value
		default:
			return nil, fmt.Errorf("unsupported filter key: %s", key)
		}
	}
	return safeFilter, nil
}

// NewMongoReviewRepository 创建MongoDB书评仓储实例
func NewMongoReviewRepository(db *mongo.Database) *MongoReviewRepository {
	reviewCollection := db.Collection("reviews")
	likeCollection := db.Collection("review_likes")

	// 创建索引
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Reviews索引
	reviewIndexes := []mongo.IndexModel{
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
			Keys: bson.D{
				{Key: "is_public", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "book_id", Value: 1},
				{Key: "rating", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "book_id", Value: 1},
				{Key: "like_count", Value: -1},
			},
		},
	}

	_, err := reviewCollection.Indexes().CreateMany(ctx, reviewIndexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create review indexes: %v\n", err)
	}

	// ReviewLikes索引
	likeIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "review_id", Value: 1},
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "review_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
	}

	_, err = likeCollection.Indexes().CreateMany(ctx, likeIndexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create review_like indexes: %v\n", err)
	}

	return &MongoReviewRepository{
		reviewCollection: reviewCollection,
		likeCollection:   likeCollection,
	}
}

// ========== 书评管理 ==========

// CreateReview 创建书评
func (r *MongoReviewRepository) CreateReview(ctx context.Context, review *social.Review) error {
	if review.ID.IsZero() {
		review.ID = primitive.NewObjectID()
	}

	if review.CreatedAt.IsZero() {
		review.CreatedAt = time.Now()
	}
	review.UpdatedAt = time.Now()

	// 初始化统计字段
	if review.LikeCount == 0 {
		review.LikeCount = 0
	}
	if review.CommentCount == 0 {
		review.CommentCount = 0
	}

	_, err := r.reviewCollection.InsertOne(ctx, review)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

// GetReviewByID 根据ID获取书评
func (r *MongoReviewRepository) GetReviewByID(ctx context.Context, reviewID string) (*social.Review, error) {
	objectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return nil, fmt.Errorf("review not found")
	}

	var review social.Review
	err = r.reviewCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&review)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	return &review, nil
}

// GetReviewsByBook 获取书籍的书评列表
func (r *MongoReviewRepository) GetReviewsByBook(ctx context.Context, bookID string, page, size int) ([]*social.Review, int64, error) {
	safeBookID, err := sanitizeReviewQueryToken("book_id", bookID)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{"book_id": safeBookID}
	sort := bson.D{{Key: "created_at", Value: -1}}

	return r.findReviews(ctx, filter, page, size, sort)
}

// GetReviewsByUser 获取用户的书评列表
func (r *MongoReviewRepository) GetReviewsByUser(ctx context.Context, userID string, page, size int) ([]*social.Review, int64, error) {
	safeUserID, err := sanitizeReviewQueryToken("user_id", userID)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{"user_id": safeUserID}
	sort := bson.D{{Key: "created_at", Value: -1}}

	return r.findReviews(ctx, filter, page, size, sort)
}

// GetPublicReviews 获取公开书评列表
func (r *MongoReviewRepository) GetPublicReviews(ctx context.Context, page, size int) ([]*social.Review, int64, error) {
	filter := bson.M{"is_public": true}
	sort := bson.D{{Key: "created_at", Value: -1}}

	return r.findReviews(ctx, filter, page, size, sort)
}

// GetReviewsByRating 根据评分获取书评
func (r *MongoReviewRepository) GetReviewsByRating(ctx context.Context, bookID string, rating int, page, size int) ([]*social.Review, int64, error) {
	safeBookID, err := sanitizeReviewQueryToken("book_id", bookID)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{
		"book_id": safeBookID,
		"rating":  rating,
	}
	sort := bson.D{{Key: "created_at", Value: -1}}

	return r.findReviews(ctx, filter, page, size, sort)
}

// UpdateReview 更新书评
func (r *MongoReviewRepository) UpdateReview(ctx context.Context, reviewID string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	updates["updated_at"] = time.Now()

	result, err := r.reviewCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("review not found")
	}

	return nil
}

// DeleteReview 删除书评
func (r *MongoReviewRepository) DeleteReview(ctx context.Context, reviewID string) error {
	objectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	result, err := r.reviewCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("review not found")
	}

	return nil
}

// ========== 书评点赞 ==========

// CreateReviewLike 创建书评点赞
func (r *MongoReviewRepository) CreateReviewLike(ctx context.Context, reviewLike *social.ReviewLike) error {
	if reviewLike.ID.IsZero() {
		reviewLike.ID = primitive.NewObjectID()
	}

	if reviewLike.CreatedAt.IsZero() {
		reviewLike.CreatedAt = time.Now()
	}

	_, err := r.likeCollection.InsertOne(ctx, reviewLike)
	if err != nil {
		return fmt.Errorf("failed to create review like: %w", err)
	}

	return nil
}

// DeleteReviewLike 删除书评点赞
func (r *MongoReviewRepository) DeleteReviewLike(ctx context.Context, reviewID, userID string) error {
	safeReviewID, err := sanitizeReviewQueryToken("review_id", reviewID)
	if err != nil {
		return err
	}
	safeUserID, err := sanitizeReviewQueryToken("user_id", userID)
	if err != nil {
		return err
	}
	_, err = r.likeCollection.DeleteOne(ctx, bson.M{
		"review_id": safeReviewID,
		"user_id":   safeUserID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete review like: %w", err)
	}

	return nil
}

// GetReviewLike 获取书评点赞记录
func (r *MongoReviewRepository) GetReviewLike(ctx context.Context, reviewID, userID string) (*social.ReviewLike, error) {
	safeReviewID, err := sanitizeReviewQueryToken("review_id", reviewID)
	if err != nil {
		return nil, err
	}
	safeUserID, err := sanitizeReviewQueryToken("user_id", userID)
	if err != nil {
		return nil, err
	}
	var reviewLike social.ReviewLike
	err = r.likeCollection.FindOne(ctx, bson.M{
		"review_id": safeReviewID,
		"user_id":   safeUserID,
	}).Decode(&reviewLike)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("review like not found")
		}
		return nil, fmt.Errorf("failed to get review like: %w", err)
	}

	return &reviewLike, nil
}

// IsReviewLiked 检查是否已点赞
func (r *MongoReviewRepository) IsReviewLiked(ctx context.Context, reviewID, userID string) (bool, error) {
	safeReviewID, err := sanitizeReviewQueryToken("review_id", reviewID)
	if err != nil {
		return false, err
	}
	safeUserID, err := sanitizeReviewQueryToken("user_id", userID)
	if err != nil {
		return false, err
	}
	count, err := r.likeCollection.CountDocuments(ctx, bson.M{
		"review_id": safeReviewID,
		"user_id":   safeUserID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check review like: %w", err)
	}

	return count > 0, nil
}

// GetReviewLikes 获取书评点赞列表
func (r *MongoReviewRepository) GetReviewLikes(ctx context.Context, reviewID string, page, size int) ([]*social.ReviewLike, int64, error) {
	safeReviewID, err := sanitizeReviewQueryToken("review_id", reviewID)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{"review_id": safeReviewID}

	// 计算总数
	total, err := r.likeCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count review likes: %w", err)
	}

	// 计算跳过数
	skip := int64((page - 1) * size)

	// 查询
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(size))

	cursor, err := r.likeCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find review likes: %w", err)
	}
	defer cursor.Close(ctx)

	var likes []*social.ReviewLike
	if err := cursor.All(ctx, &likes); err != nil {
		return nil, 0, fmt.Errorf("failed to decode review likes: %w", err)
	}

	return likes, total, nil
}

// IncrementReviewLikeCount 增加书评点赞数
func (r *MongoReviewRepository) IncrementReviewLikeCount(ctx context.Context, reviewID string) error {
	objectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	_, err = r.reviewCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$inc": bson.M{"like_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to increment review like count: %w", err)
	}

	return nil
}

// DecrementReviewLikeCount 减少书评点赞数
func (r *MongoReviewRepository) DecrementReviewLikeCount(ctx context.Context, reviewID string) error {
	objectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	_, err = r.reviewCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$inc": bson.M{"like_count": -1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to decrement review like count: %w", err)
	}

	return nil
}

// ========== 统计 ==========

// GetAverageRating 获取书籍平均评分
func (r *MongoReviewRepository) GetAverageRating(ctx context.Context, bookID string) (float64, error) {
	safeBookID, err := sanitizeReviewQueryToken("book_id", bookID)
	if err != nil {
		return 0, err
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"book_id": safeBookID,
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"avg":   bson.M{"$avg": "$rating"},
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.reviewCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("failed to get average rating: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0, fmt.Errorf("failed to decode average rating: %w", err)
	}

	if len(results) == 0 {
		return 0, nil
	}

	avg, ok := results[0]["avg"].(float64)
	if !ok {
		return 0, nil
	}

	return avg, nil
}

// GetRatingDistribution 获取评分分布
func (r *MongoReviewRepository) GetRatingDistribution(ctx context.Context, bookID string) (map[int]int64, error) {
	safeBookID, err := sanitizeReviewQueryToken("book_id", bookID)
	if err != nil {
		return nil, err
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"book_id": safeBookID,
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$rating",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.reviewCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating distribution: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode rating distribution: %w", err)
	}

	distribution := make(map[int]int64)
	for _, result := range results {
		rating, ok := result["_id"].(int32)
		if !ok {
			continue
		}
		count, ok := result["count"].(int32)
		if !ok {
			continue
		}
		distribution[int(rating)] = int64(count)
	}

	// 确保所有评分都有值
	for i := 1; i <= 5; i++ {
		if _, exists := distribution[i]; !exists {
			distribution[i] = 0
		}
	}

	return distribution, nil
}

// CountReviews 统计书评数
func (r *MongoReviewRepository) CountReviews(ctx context.Context, bookID string) (int64, error) {
	safeBookID, err := sanitizeReviewQueryToken("book_id", bookID)
	if err != nil {
		return 0, err
	}
	count, err := r.reviewCollection.CountDocuments(ctx, bson.M{"book_id": safeBookID})
	if err != nil {
		return 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	return count, nil
}

// CountUserReviews 统计用户书评数
func (r *MongoReviewRepository) CountUserReviews(ctx context.Context, userID string) (int64, error) {
	safeUserID, err := sanitizeReviewQueryToken("user_id", userID)
	if err != nil {
		return 0, err
	}
	count, err := r.reviewCollection.CountDocuments(ctx, bson.M{"user_id": safeUserID})
	if err != nil {
		return 0, fmt.Errorf("failed to count user reviews: %w", err)
	}

	return count, nil
}

// Health 健康检查
func (r *MongoReviewRepository) Health(ctx context.Context) error {
	return r.reviewCollection.Database().Client().Ping(ctx, nil)
}

// findReviews 通用的查询书评方法
func (r *MongoReviewRepository) findReviews(ctx context.Context, filter bson.M, page, size int, sort bson.D) ([]*social.Review, int64, error) {
	safeFilter, err := sanitizeReviewFilter(filter)
	if err != nil {
		return nil, 0, err
	}

	// 计算总数
	total, err := r.reviewCollection.CountDocuments(ctx, safeFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	// 计算跳过数
	skip := int64((page - 1) * size)

	// 查询
	opts := options.Find().
		SetSort(sort).
		SetSkip(skip).
		SetLimit(int64(size))

	cursor, err := r.reviewCollection.Find(ctx, safeFilter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find reviews: %w", err)
	}
	defer cursor.Close(ctx)

	var reviews []*social.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, 0, fmt.Errorf("failed to decode reviews: %w", err)
	}

	return reviews, total, nil
}
