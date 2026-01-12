package messaging

import (
	"Qingyu_backend/models/messaging"
	infra "Qingyu_backend/repository/interfaces/infrastructure"
	interfaces "Qingyu_backend/repository/interfaces/messaging"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoAnnouncementRepository MongoDB公告仓储实现
type MongoAnnouncementRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

// NewMongoAnnouncementRepository 创建MongoDB公告仓储实例
func NewMongoAnnouncementRepository(client *mongo.Client, database string) interfaces.AnnouncementRepository {
	return &MongoAnnouncementRepository{
		collection: client.Database(database).Collection("announcements"),
		client:     client,
	}
}

// Create 创建公告
func (r *MongoAnnouncementRepository) Create(ctx context.Context, announcement *messaging.Announcement) error {
	if announcement == nil {
		return errors.New("announcement cannot be nil")
	}

	announcement.CreatedAt = time.Now()
	announcement.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, announcement)
	if err != nil {
		return err
	}

	// ID现在是string类型，从MongoDB的ObjectID转换
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		announcement.ID = oid.Hex()
	}
	return nil
}

// GetByID 根据ID获取公告
func (r *MongoAnnouncementRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*messaging.Announcement, error) {
	var announcement messaging.Announcement
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&announcement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &announcement, nil
}

// Update 更新公告
func (r *MongoAnnouncementRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
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
		return errors.New("announcement not found")
	}

	return nil
}

// Delete 删除公告
func (r *MongoAnnouncementRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("announcement not found")
	}

	return nil
}

// Health 健康检查
func (r *MongoAnnouncementRepository) Health(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}

// Count 统计公告数量
func (r *MongoAnnouncementRepository) Count(ctx context.Context, filter infra.Filter) (int64, error) {
	var query bson.M
	if filter != nil {
		query = bson.M(filter.GetConditions())
	} else {
		query = bson.M{}
	}
	return r.collection.CountDocuments(ctx, query)
}

// List 根据过滤条件列出公告
func (r *MongoAnnouncementRepository) List(ctx context.Context, filter infra.Filter) ([]*messaging.Announcement, error) {
	var query bson.M
	if filter != nil {
		query = bson.M(filter.GetConditions())
	} else {
		query = bson.M{}
	}

	opts := options.Find()
	if filter != nil {
		sort := filter.GetSort()
		if len(sort) > 0 {
			var sortDoc bson.D
			for k, v := range sort {
				sortDoc = append(sortDoc, bson.E{Key: k, Value: v})
			}
			opts.SetSort(sortDoc)
		}

		// Type assert to access pagination methods
		if announcementFilter, ok := filter.(*messaging.AnnouncementFilter); ok {
			if limit := announcementFilter.GetLimit(); limit > 0 {
				opts.SetLimit(int64(limit))
			}

			if offset := announcementFilter.GetOffset(); offset > 0 {
				opts.SetSkip(int64(offset))
			}
		}
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var announcements []*messaging.Announcement
	if err := cursor.All(ctx, &announcements); err != nil {
		return nil, err
	}
	return announcements, nil
}

// Exists 判断公告是否存在
func (r *MongoAnnouncementRepository) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetActive 获取激活的公告列表
func (r *MongoAnnouncementRepository) GetActive(ctx context.Context, targetRole string, limit, offset int) ([]*messaging.Announcement, error) {
	now := time.Now()
	query := bson.M{
		"is_active": true,
		"$or": []bson.M{
			{"target_role": targetRole},
			{"target_role": "all"},
		},
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"start_time": nil},
					{"start_time": bson.M{"$lte": now}},
				},
			},
			{
				"$or": []bson.M{
					{"end_time": nil},
					{"end_time": bson.M{"$gte": now}},
				},
			},
		},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "priority", Value: -1}, {Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var announcements []*messaging.Announcement
	if err = cursor.All(ctx, &announcements); err != nil {
		return nil, err
	}

	return announcements, nil
}

// GetByType 根据类型获取公告列表
func (r *MongoAnnouncementRepository) GetByType(ctx context.Context, announcementType string, limit, offset int) ([]*messaging.Announcement, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "priority", Value: -1}, {Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"type": announcementType, "is_active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var announcements []*messaging.Announcement
	if err = cursor.All(ctx, &announcements); err != nil {
		return nil, err
	}

	return announcements, nil
}

// GetByTimeRange 根据时间范围获取公告列表
func (r *MongoAnnouncementRepository) GetByTimeRange(ctx context.Context, startTime, endTime *time.Time, limit, offset int) ([]*messaging.Announcement, error) {
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
		SetSort(bson.D{{Key: "priority", Value: -1}, {Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var announcements []*messaging.Announcement
	if err = cursor.All(ctx, &announcements); err != nil {
		return nil, err
	}

	return announcements, nil
}

// GetEffective 获取当前有效的公告
func (r *MongoAnnouncementRepository) GetEffective(ctx context.Context, targetRole string, limit int) ([]*messaging.Announcement, error) {
	now := time.Now()
	query := bson.M{
		"is_active": true,
		"$or": []bson.M{
			{"target_role": targetRole},
			{"target_role": "all"},
		},
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"start_time": nil},
					{"start_time": bson.M{"$lte": now}},
				},
			},
			{
				"$or": []bson.M{
					{"end_time": nil},
					{"end_time": bson.M{"$gte": now}},
				},
			},
		},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "priority", Value: -1}, {Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var announcements []*messaging.Announcement
	if err = cursor.All(ctx, &announcements); err != nil {
		return nil, err
	}

	return announcements, nil
}

// IncrementViewCount 增加查看次数
func (r *MongoAnnouncementRepository) IncrementViewCount(ctx context.Context, announcementID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": announcementID},
		bson.M{
			"$inc": bson.M{"view_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// GetViewStats 获取查看统计
func (r *MongoAnnouncementRepository) GetViewStats(ctx context.Context, announcementID primitive.ObjectID) (int64, error) {
	var announcement messaging.Announcement
	err := r.collection.FindOne(ctx, bson.M{"_id": announcementID}, options.FindOne().SetProjection(bson.M{"view_count": 1})).Decode(&announcement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}
	return announcement.ViewCount, nil
}

// BatchUpdateStatus 批量更新状态
func (r *MongoAnnouncementRepository) BatchUpdateStatus(ctx context.Context, announcementIDs []primitive.ObjectID, isActive bool) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": announcementIDs}},
		bson.M{
			"$set": bson.M{
				"is_active":  isActive,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// BatchDelete 批量删除
func (r *MongoAnnouncementRepository) BatchDelete(ctx context.Context, announcementIDs []primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(
		ctx,
		bson.M{"_id": bson.M{"$in": announcementIDs}},
	)
	return err
}

// Transaction 执行事务
func (r *MongoAnnouncementRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		return nil, fn(sc)
	})
	return err
}
