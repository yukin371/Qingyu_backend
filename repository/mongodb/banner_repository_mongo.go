package mongodb

import (
	"context"
	"time"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"Qingyu_backend/models/reading/bookstore"
	"Qingyu_backend/repository/interfaces"
)

// MongoBannerRepository MongoDB Banner仓储实现
type MongoBannerRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

// NewMongoBannerRepository 创建MongoDB Banner仓储实例
func NewMongoBannerRepository(client *mongo.Client, database string) interfaces.BannerRepository {
	return &MongoBannerRepository{
		collection: client.Database(database).Collection("banners"),
		client:     client,
	}
}

// Create 创建Banner
func (r *MongoBannerRepository) Create(ctx context.Context, banner *bookstore.Banner) error {
	if banner == nil {
		return errors.New("banner cannot be nil")
	}
	
	banner.CreatedAt = time.Now()
	banner.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, banner)
	if err != nil {
		return err
	}
	
	banner.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 根据ID获取Banner
func (r *MongoBannerRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Banner, error) {
	var banner bookstore.Banner
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&banner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &banner, nil
}

// Update 更新Banner
func (r *MongoBannerRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	if err != nil {
		return err
	}
	
	if result.MatchedCount == 0 {
		return errors.New("banner not found")
	}
	
	return nil
}

// Delete 删除Banner
func (r *MongoBannerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	
	if result.DeletedCount == 0 {
		return errors.New("banner not found")
	}
	
	return nil
}

// Health 健康检查
func (r *MongoBannerRepository) Health(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}

// GetActive 获取激活的Banner列表
func (r *MongoBannerRepository) GetActive(ctx context.Context, limit, offset int) ([]*bookstore.Banner, error) {
	now := time.Now()
	query := bson.M{
		"is_active": true,
		"$or": []bson.M{
			{"start_time": nil, "end_time": nil},
			{"start_time": bson.M{"$lte": now}, "end_time": nil},
			{"start_time": nil, "end_time": bson.M{"$gte": now}},
			{"start_time": bson.M{"$lte": now}, "end_time": bson.M{"$gte": now}},
		},
	}
	
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{"sort_order", 1}, {"created_at", -1}})
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var banners []*bookstore.Banner
	if err = cursor.All(ctx, &banners); err != nil {
		return nil, err
	}
	
	return banners, nil
}

// GetByTargetType 根据目标类型获取Banner列表
func (r *MongoBannerRepository) GetByTargetType(ctx context.Context, targetType string, limit, offset int) ([]*bookstore.Banner, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{"sort_order", 1}, {"created_at", -1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{"target_type": targetType, "is_active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var banners []*bookstore.Banner
	if err = cursor.All(ctx, &banners); err != nil {
		return nil, err
	}
	
	return banners, nil
}

// GetByTimeRange 根据时间范围获取Banner列表
func (r *MongoBannerRepository) GetByTimeRange(ctx context.Context, startTime, endTime *time.Time, limit, offset int) ([]*bookstore.Banner, error) {
	query := bson.M{"is_active": true}
	
	if startTime != nil && endTime != nil {
		query["$or"] = []bson.M{
			{"start_time": nil, "end_time": nil},
			{"start_time": bson.M{"$lte": *endTime}, "end_time": bson.M{"$gte": *startTime}},
		}
	} else if startTime != nil {
		query["$or"] = []bson.M{
			{"start_time": nil},
			{"start_time": bson.M{"$gte": *startTime}},
		}
	} else if endTime != nil {
		query["$or"] = []bson.M{
			{"end_time": nil},
			{"end_time": bson.M{"$lte": *endTime}},
		}
	}
	
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{"sort_order", 1}, {"created_at", -1}})
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var banners []*bookstore.Banner
	if err = cursor.All(ctx, &banners); err != nil {
		return nil, err
	}
	
	return banners, nil
}

// IncrementClickCount 增加点击次数
func (r *MongoBannerRepository) IncrementClickCount(ctx context.Context, bannerID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": bannerID},
		bson.M{
			"$inc": bson.M{"click_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// GetClickStats 获取点击统计
func (r *MongoBannerRepository) GetClickStats(ctx context.Context, bannerID primitive.ObjectID) (int64, error) {
	var banner bookstore.Banner
	err := r.collection.FindOne(ctx, bson.M{"_id": bannerID}, options.FindOne().SetProjection(bson.M{"click_count": 1})).Decode(&banner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}
	return banner.ClickCount, nil
}

// BatchUpdateStatus 批量更新状态
func (r *MongoBannerRepository) BatchUpdateStatus(ctx context.Context, bannerIDs []primitive.ObjectID, isActive bool) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": bannerIDs}},
		bson.M{
			"$set": bson.M{
				"is_active": isActive,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// Transaction 执行事务
func (r *MongoBannerRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	
	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		return session.WithTransaction(sc, func(sc mongo.SessionContext) (interface{}, error) {
			return nil, fn(sc)
		})
	})
}