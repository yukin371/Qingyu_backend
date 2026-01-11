package reader

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	readerModel "Qingyu_backend/models/reader"
	readerRepo "Qingyu_backend/repository/interfaces/reader"
)

// ReaderThemeRepositoryMongo MongoDB实现的读者主题仓储
type ReaderThemeRepositoryMongo struct {
	collection *mongo.Collection
}

// NewReaderThemeRepositoryMongo 创建MongoDB读者主题仓储
func NewReaderThemeRepositoryMongo(db *mongo.Database) readerRepo.ReaderThemeRepository {
	return &ReaderThemeRepositoryMongo{
		collection: db.Collection("reader_themes"),
	}
}

// GetBuiltInThemes 获取所有内置主题
func (r *ReaderThemeRepositoryMongo) GetBuiltInThemes(ctx context.Context) ([]*readerModel.ReaderTheme, error) {
	filter := bson.M{"is_built_in": true}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("获取内置主题失败: %w", err)
	}
	defer cursor.Close(ctx)

	var themes []*readerModel.ReaderTheme
	if err = cursor.All(ctx, &themes); err != nil {
		return nil, fmt.Errorf("解析内置主题失败: %w", err)
	}

	// 如果数据库中没有内置主题，返回默认的内置主题
	if len(themes) == 0 {
		return readerModel.BuiltInThemes, nil
	}

	return themes, nil
}

// GetThemeByName 根据名称获取主题
func (r *ReaderThemeRepositoryMongo) GetThemeByName(ctx context.Context, name string) (*readerModel.ReaderTheme, error) {
	filter := bson.M{"name": name}
	var theme readerModel.ReaderTheme
	err := r.collection.FindOne(ctx, filter).Decode(&theme)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("主题不存在: %s", name)
		}
		return nil, fmt.Errorf("查询主题失败: %w", err)
	}
	return &theme, nil
}

// CreateTheme 创建主题
func (r *ReaderThemeRepositoryMongo) CreateTheme(ctx context.Context, theme *readerModel.ReaderTheme) error {
	now := time.Now()
	theme.CreatedAt = now
	theme.UpdatedAt = now
	theme.UseCount = 0

	result, err := r.collection.InsertOne(ctx, theme)
	if err != nil {
		return fmt.Errorf("创建主题失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		theme.ID = oid.Hex()
	}

	return nil
}

// GetTheme 获取主题
func (r *ReaderThemeRepositoryMongo) GetTheme(ctx context.Context, themeID string) (*readerModel.ReaderTheme, error) {
	oid, err := primitive.ObjectIDFromHex(themeID)
	if err != nil {
		return nil, fmt.Errorf("无效的主题ID: %w", err)
	}

	filter := bson.M{"_id": oid}
	var theme readerModel.ReaderTheme
	err = r.collection.FindOne(ctx, filter).Decode(&theme)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("主题不存在: %s", themeID)
		}
		return nil, fmt.Errorf("查询主题失败: %w", err)
	}
	return &theme, nil
}

// UpdateTheme 更新主题
func (r *ReaderThemeRepositoryMongo) UpdateTheme(ctx context.Context, themeID string, updates map[string]interface{}) error {
	oid, err := primitive.ObjectIDFromHex(themeID)
	if err != nil {
		return fmt.Errorf("无效的主题ID: %w", err)
	}

	updates["updated_at"] = time.Now()
	filter := bson.M{"_id": oid}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新主题失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("主题不存在: %s", themeID)
	}

	return nil
}

// DeleteTheme 删除主题
func (r *ReaderThemeRepositoryMongo) DeleteTheme(ctx context.Context, themeID string) error {
	oid, err := primitive.ObjectIDFromHex(themeID)
	if err != nil {
		return fmt.Errorf("无效的主题ID: %w", err)
	}

	filter := bson.M{"_id": oid, "is_built_in": false} // 不允许删除内置主题
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除主题失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("主题不存在或为内置主题: %s", themeID)
	}

	return nil
}

// GetUserThemes 获取用户主题
func (r *ReaderThemeRepositoryMongo) GetUserThemes(ctx context.Context, userID string) ([]*readerModel.ReaderTheme, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"creator_id": userID},
			{"is_built_in": true},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("获取用户主题失败: %w", err)
	}
	defer cursor.Close(ctx)

	var themes []*readerModel.ReaderTheme
	if err = cursor.All(ctx, &themes); err != nil {
		return nil, fmt.Errorf("解析用户主题失败: %w", err)
	}

	return themes, nil
}

// GetActiveTheme 获取用户激活的主题
func (r *ReaderThemeRepositoryMongo) GetActiveTheme(ctx context.Context, userID string) (*readerModel.ReaderTheme, error) {
	// 查找用户自定义的激活主题
	filter := bson.M{
		"creator_id": userID,
		"is_active":  true,
	}

	var theme readerModel.ReaderTheme
	err := r.collection.FindOne(ctx, filter).Decode(&theme)
	if err == nil {
		return &theme, nil
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("查询激活主题失败: %w", err)
	}

	// 如果没有自定义激活主题，返回默认的内置主题
	return r.GetThemeByName(ctx, "light")
}

// SetActiveTheme 设置激活主题
func (r *ReaderThemeRepositoryMongo) SetActiveTheme(ctx context.Context, userID, themeID string) error {
	// 首先验证主题存在
	theme, err := r.GetTheme(ctx, themeID)
	if err != nil {
		return err
	}

	// 取消用户所有主题的激活状态
	filter := bson.M{"creator_id": userID}
	update := bson.M{"$set": bson.M{"is_active": false}}
	_, err = r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("取消激活状态失败: %w", err)
	}

	// 如果主题是内置主题，需要为用户创建一个副本
	if theme.IsBuiltIn {
		newTheme := *theme
		newTheme.ID = "" // 让MongoDB生成新ID
		newTheme.CreatorID = userID
		newTheme.IsBuiltIn = false
		newTheme.IsActive = true
		newTheme.UseCount = 0
		return r.CreateTheme(ctx, &newTheme)
	}

	// 激活自定义主题
	filter = bson.M{"_id": theme.ID}
	update = bson.M{"$set": bson.M{"is_active": true, "updated_at": time.Now()}}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("设置激活主题失败: %w", err)
	}

	return nil
}

// GetPublicThemes 获取公开主题
func (r *ReaderThemeRepositoryMongo) GetPublicThemes(ctx context.Context, page, pageSize int) ([]*readerModel.ReaderTheme, int64, error) {
	filter := bson.M{"is_public": true}

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计公开主题失败: %w", err)
	}

	// 分页查询
	skip := int64((page - 1) * pageSize)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "use_count", Value: -1}, {Key: "created_at", Value: -1}}) // 按使用次数和创建时间排序

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询公开主题失败: %w", err)
	}
	defer cursor.Close(ctx)

	var themes []*readerModel.ReaderTheme
	if err = cursor.All(ctx, &themes); err != nil {
		return nil, 0, fmt.Errorf("解析公开主题失败: %w", err)
	}

	return themes, total, nil
}

// IncrementUseCount 增加使用次数
func (r *ReaderThemeRepositoryMongo) IncrementUseCount(ctx context.Context, themeID string) error {
	oid, err := primitive.ObjectIDFromHex(themeID)
	if err != nil {
		return fmt.Errorf("无效的主题ID: %w", err)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$inc": bson.M{"use_count": 1}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("增加使用次数失败: %w", err)
	}

	return nil
}

// BatchGetThemes 批量获取主题
func (r *ReaderThemeRepositoryMongo) BatchGetThemes(ctx context.Context, themeIDs []string) (map[string]*readerModel.ReaderTheme, error) {
	oids := make([]primitive.ObjectID, 0, len(themeIDs))
	for _, id := range themeIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue // 跳过无效ID
		}
		oids = append(oids, oid)
	}

	filter := bson.M{"_id": bson.M{"$in": oids}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("批量查询主题失败: %w", err)
	}
	defer cursor.Close(ctx)

	themes := make(map[string]*readerModel.ReaderTheme)
	for cursor.Next(ctx) {
		var theme readerModel.ReaderTheme
		if err := cursor.Decode(&theme); err != nil {
			continue
		}
		themes[theme.ID] = &theme
	}

	return themes, nil
}

// Health 健康检查
func (r *ReaderThemeRepositoryMongo) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}
