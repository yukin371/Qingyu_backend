package mongodb

import (
	"context"
	"time"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"Qingyu_backend/models/reading/reader"
	"Qingyu_backend/repository/base"
	"Qingyu_backend/repository/reading"
)

// MongoReadingSettingsRepository MongoDB阅读设置仓储实现
type MongoReadingSettingsRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
	queryBuilder base.QueryBuilder
}

// NewMongoReadingSettingsRepository 创建MongoDB阅读设置仓储实例
func NewMongoReadingSettingsRepository(db *mongo.Database) reading.ReadingSettingsRepository {
	return &MongoReadingSettingsRepository{
		db:         db,
		collection: db.Collection("reading_settings"),
		queryBuilder: base.NewMongoQueryBuilder(),
	}
}

// Create 创建阅读设置
func (r *MongoReadingSettingsRepository) Create(ctx context.Context, settings **reader.ReadingSettings) error {
	if settings == nil || *settings == nil {
		return nil
	}
	
	actualSettings := *settings
	actualSettings.CreatedAt = time.Now()
	actualSettings.UpdatedAt = time.Now()
	
	_, err := r.collection.InsertOne(ctx, actualSettings)
	return err
}

// GetByID 根据ID获取阅读设置
func (r *MongoReadingSettingsRepository) GetByID(ctx context.Context, id string) (**reader.ReadingSettings, error) {
	var settings reader.ReadingSettings
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	result := &settings
	return &result, nil
}

// Update 更新阅读设置
func (r *MongoReadingSettingsRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	return err
}

// Delete 删除阅读设置
func (r *MongoReadingSettingsRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// List 获取阅读设置列表
func (r *MongoReadingSettingsRepository) List(ctx context.Context, filter base.Filter) ([]**reader.ReadingSettings, error) {
	mongoFilter := filter.GetConditions()
	
	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var settings []*reader.ReadingSettings
	if err = cursor.All(ctx, &settings); err != nil {
		return nil, err
	}
	
	// 转换为双指针类型
	result := make([]**reader.ReadingSettings, len(settings))
	for i, setting := range settings {
		result[i] = &setting
	}
	
	return result, nil
}

// Count 统计阅读设置数量
func (r *MongoReadingSettingsRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	mongoFilter := filter.GetConditions()
	return r.collection.CountDocuments(ctx, mongoFilter)
}

// Health 健康检查
func (r *MongoReadingSettingsRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}

// GetByUserID 根据用户ID获取阅读设置
func (r *MongoReadingSettingsRepository) GetByUserID(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	var settings reader.ReadingSettings
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &settings, nil
}

// UpdateByUserID 根据用户ID更新阅读设置
func (r *MongoReadingSettingsRepository) UpdateByUserID(ctx context.Context, userID string, settings *reader.ReadingSettings) error {
	settings.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": settings},
		options.Update().SetUpsert(true),
	)
	return err
}

// CreateDefaultSettings 为用户创建默认阅读设置
func (r *MongoReadingSettingsRepository) CreateDefaultSettings(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	settings := &reader.ReadingSettings{
		UserID:      userID,
		FontFamily:  "Microsoft YaHei",
		FontSize:    16,
		LineHeight:  1.5,
		Theme:       "light",
		Background:  "#ffffff",
		PageMode:    1, // 滑动模式
		AutoScroll:  false,
		ScrollSpeed: 3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	err := r.Create(ctx, &settings)
	if err != nil {
		return nil, err
	}
	
	return settings, nil
}

// ExistsByUserID 检查用户是否已有阅读设置
func (r *MongoReadingSettingsRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// BatchCreate 批量创建阅读设置
func (r *MongoReadingSettingsRepository) BatchCreate(ctx context.Context, settings []**reader.ReadingSettings) error {
	if len(settings) == 0 {
		return nil
	}
	
	// 验证设置对象
	for _, setting := range settings {
		if setting == nil || *setting == nil {
			continue
		}
	}
	
	// 准备批量插入的文档
	var documents []interface{}
	now := time.Now()
	for _, setting := range settings {
		if setting == nil || *setting == nil {
			continue
		}
		(*setting).CreatedAt = now
		(*setting).UpdatedAt = now
		documents = append(documents, *setting)
	}
	
	if len(documents) == 0 {
		return nil
	}
	
	_, err := r.collection.InsertMany(ctx, documents)
	return err
}

// BatchUpdate 批量更新阅读设置
func (r *MongoReadingSettingsRepository) BatchUpdate(ctx context.Context, ids []string, updates map[string]interface{}) error {
	if len(ids) == 0 {
		return nil
	}
	
	updates["updated_at"] = time.Now()
	filter := bson.M{"_id": bson.M{"$in": ids}}
	update := bson.M{"$set": updates}
	
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// BatchDelete 批量删除阅读设置
func (r *MongoReadingSettingsRepository) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	
	filter := bson.M{"_id": bson.M{"$in": ids}}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// Exists 检查阅读设置是否存在
func (r *MongoReadingSettingsRepository) Exists(ctx context.Context, id string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindWithPagination 分页查询阅读设置
func (r *MongoReadingSettingsRepository) FindWithPagination(ctx context.Context, filter base.Filter, pagination base.Pagination) (*base.PagedResult[*reader.ReadingSettings], error) {
	mongoFilter := filter.GetConditions()
	
	// 计算总数
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	
	// 构建查询选项
	opts := options.Find()
	if pagination.Skip > 0 {
		opts.SetSkip(int64(pagination.Skip))
	} else if pagination.Page > 0 && pagination.PageSize > 0 {
		skip := (pagination.Page - 1) * pagination.PageSize
		opts.SetSkip(int64(skip))
	}
	
	if pagination.PageSize > 0 {
		opts.SetLimit(int64(pagination.PageSize))
	}
	
	// 执行查询
	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	// 解析结果为双指针类型
	var settingsData []**reader.ReadingSettings
	for cursor.Next(ctx) {
		var settings reader.ReadingSettings
		if err := cursor.Decode(&settings); err != nil {
			return nil, err
		}
		settingsPtr := &settings
		settingsData = append(settingsData, &settingsPtr)
	}
	
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	
	// 使用NewPagedResult创建分页结果
	return base.NewPagedResult[*reader.ReadingSettings](settingsData, total, pagination), nil
}