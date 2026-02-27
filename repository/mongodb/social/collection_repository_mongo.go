package social

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/social"
	"Qingyu_backend/repository/mongodb/base"
)

var safeCollectionTagPattern = regexp.MustCompile(`^[\p{L}\p{N}_\-\s]{1,32}$`)

func normalizeObjectIDHex(field, value string) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return "", fmt.Errorf("%s格式非法", field)
	}
	return objectID.Hex(), nil
}

func validateCollectionTag(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("标签不能为空")
	}
	if !safeCollectionTagPattern.MatchString(value) {
		return fmt.Errorf("标签格式非法")
	}
	return nil
}

// MongoCollectionRepository MongoDB收藏仓储实现
type MongoCollectionRepository struct {
	*base.BaseMongoRepository // 嵌入基类，管理主collection (collections)
	folderColl *mongo.Collection // 独立管理folder collection
}

// NewMongoCollectionRepository 创建MongoDB收藏仓储实例
func NewMongoCollectionRepository(db *mongo.Database) *MongoCollectionRepository {
	// 创建索引
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Collection索引
	collectionIndexes := []mongo.IndexModel{
		{
			// 防重索引：用户不能重复收藏同一本书
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "book_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			// 用户收藏列表查询
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			// 收藏夹下的收藏查询
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "folder_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			// 公开收藏查询
			Keys: bson.D{
				{Key: "is_public", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			// 标签查询
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "tags", Value: 1},
			},
		},
	}

	collectionColl := db.Collection("collections")
	_, err := collectionColl.Indexes().CreateMany(ctx, collectionIndexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create collection indexes: %v\n", err)
	}

	// Folder索引
	folderIndexes := []mongo.IndexModel{
		{
			// 用户收藏夹列表
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			// 公开收藏夹查询
			Keys: bson.D{
				{Key: "is_public", Value: 1},
				{Key: "book_count", Value: -1},
			},
		},
	}

	folderColl := db.Collection("collection_folders")
	_, err = folderColl.Indexes().CreateMany(ctx, folderIndexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create folder indexes: %v\n", err)
	}

	return &MongoCollectionRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "collections"),
		folderColl:          folderColl,
	}
}

// =========================
// 收藏管理
// =========================

// Create 创建收藏
func (r *MongoCollectionRepository) Create(ctx context.Context, collection *social.Collection) error {
	if collection.ID.IsZero() {
		collection.ID = primitive.NewObjectID()
	}

	if collection.CreatedAt.IsZero() {
		collection.CreatedAt = time.Now()
	}
	collection.UpdatedAt = time.Now()

	// 参数验证
	if collection.UserID == "" || collection.BookID == "" {
		return fmt.Errorf("用户ID和书籍ID不能为空")
	}

	_, err := r.GetCollection().InsertOne(ctx, collection)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("已经收藏过该书籍")
		}
		return fmt.Errorf("failed to create collection: %w", err)
	}

	return nil
}

// GetByID 根据ID获取收藏
func (r *MongoCollectionRepository) GetByID(ctx context.Context, id string) (*social.Collection, error) {
	objectID, err := r.ParseID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid collection ID: %w", err)
	}

	var collection social.Collection
	err = r.GetCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&collection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("collection not found")
		}
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	return &collection, nil
}

// GetByUserAndBook 根据用户ID和书籍ID获取收藏
func (r *MongoCollectionRepository) GetByUserAndBook(ctx context.Context, userID, bookID string) (*social.Collection, error) {
	if userID == "" || bookID == "" {
		return nil, fmt.Errorf("用户ID和书籍ID不能为空")
	}
	userIDHex, err := normalizeObjectIDHex("user_id", userID)
	if err != nil {
		return nil, err
	}
	bookIDHex, err := normalizeObjectIDHex("book_id", bookID)
	if err != nil {
		return nil, err
	}

	var collection social.Collection
	err = r.GetCollection().FindOne(ctx, bson.D{
		{Key: "user_id", Value: userIDHex},
		{Key: "book_id", Value: bookIDHex},
	}).Decode(&collection)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未收藏不报错，返回nil
		}
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	return &collection, nil
}

// GetCollectionsByUser 获取用户的收藏列表
func (r *MongoCollectionRepository) GetCollectionsByUser(ctx context.Context, userID string, folderID string, page, size int) ([]*social.Collection, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}
	userIDHex, err := normalizeObjectIDHex("user_id", userID)
	if err != nil {
		return nil, 0, err
	}

	// 构建过滤条件
	filter := bson.M{"user_id": userIDHex}
	if folderID != "" {
		folderIDHex, normalizeErr := normalizeObjectIDHex("folder_id", folderID)
		if normalizeErr != nil {
			return nil, 0, normalizeErr
		}
		filter["folder_id"] = folderIDHex
	}

	// 计算总数
	total, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count collections: %w", err)
	}

	// 计算跳过数
	skip := int64((page - 1) * size)

	// 查询
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(size))

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find collections: %w", err)
	}
	defer cursor.Close(ctx)

	var collections []*social.Collection
	if err := cursor.All(ctx, &collections); err != nil {
		return nil, 0, fmt.Errorf("failed to decode collections: %w", err)
	}

	return collections, total, nil
}

// GetCollectionsByTag 根据标签获取收藏
func (r *MongoCollectionRepository) GetCollectionsByTag(ctx context.Context, userID string, tag string, page, size int) ([]*social.Collection, int64, error) {
	if userID == "" || tag == "" {
		return nil, 0, fmt.Errorf("用户ID和标签不能为空")
	}
	userIDHex, err := normalizeObjectIDHex("user_id", userID)
	if err != nil {
		return nil, 0, err
	}
	if err := validateCollectionTag(tag); err != nil {
		return nil, 0, err
	}

	// 为避免将用户输入直接拼到查询条件中，先按用户拉取，再在内存中过滤标签。
	userFilter := bson.D{
		{Key: "user_id", Value: userIDHex},
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.GetCollection().Find(ctx, userFilter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find collections by user: %w", err)
	}
	defer cursor.Close(ctx)

	var allCollections []*social.Collection
	if err := cursor.All(ctx, &allCollections); err != nil {
		return nil, 0, fmt.Errorf("failed to decode collections: %w", err)
	}

	filtered := make([]*social.Collection, 0, len(allCollections))
	for _, collection := range allCollections {
		for _, existingTag := range collection.Tags {
			if existingTag == tag {
				filtered = append(filtered, collection)
				break
			}
		}
	}

	total := int64(len(filtered))
	start := (page - 1) * size
	if start < 0 || start >= len(filtered) {
		return []*social.Collection{}, total, nil
	}
	end := start + size
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

// Update 更新收藏
func (r *MongoCollectionRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("invalid collection ID: %w", err)
	}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	result, err := r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to update collection: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("collection not found")
	}

	return nil
}

// Delete 删除收藏
func (r *MongoCollectionRepository) Delete(ctx context.Context, id string) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("invalid collection ID: %w", err)
	}

	result, err := r.GetCollection().DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("collection not found")
	}

	return nil
}

// =========================
// 收藏夹管理
// =========================

// CreateFolder 创建收藏夹
func (r *MongoCollectionRepository) CreateFolder(ctx context.Context, folder *social.CollectionFolder) error {
	if folder.ID.IsZero() {
		folder.ID = primitive.NewObjectID()
	}

	if folder.CreatedAt.IsZero() {
		folder.CreatedAt = time.Now()
	}
	folder.UpdatedAt = time.Now()

	// 参数验证
	if folder.UserID == "" || folder.Name == "" {
		return fmt.Errorf("用户ID和收藏夹名称不能为空")
	}

	// 初始化计数器
	folder.BookCount = 0

	_, err := r.folderColl.InsertOne(ctx, folder)
	if err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	return nil
}

// GetFolderByID 根据ID获取收藏夹
func (r *MongoCollectionRepository) GetFolderByID(ctx context.Context, id string) (*social.CollectionFolder, error) {
	objectID, err := r.ParseID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid folder ID: %w", err)
	}

	var folder social.CollectionFolder
	err = r.folderColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&folder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("folder not found")
		}
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	return &folder, nil
}

// GetFoldersByUser 获取用户的收藏夹列表
func (r *MongoCollectionRepository) GetFoldersByUser(ctx context.Context, userID string) ([]*social.CollectionFolder, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.folderColl.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find folders: %w", err)
	}
	defer cursor.Close(ctx)

	var folders []*social.CollectionFolder
	if err := cursor.All(ctx, &folders); err != nil {
		return nil, fmt.Errorf("failed to decode folders: %w", err)
	}

	return folders, nil
}

// UpdateFolder 更新收藏夹
func (r *MongoCollectionRepository) UpdateFolder(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("invalid folder ID: %w", err)
	}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	result, err := r.folderColl.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to update folder: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("folder not found")
	}

	return nil
}

// DeleteFolder 删除收藏夹
func (r *MongoCollectionRepository) DeleteFolder(ctx context.Context, id string) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("invalid folder ID: %w", err)
	}

	result, err := r.folderColl.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("folder not found")
	}

	return nil
}

// IncrementFolderBookCount 增加收藏夹书籍数量
func (r *MongoCollectionRepository) IncrementFolderBookCount(ctx context.Context, folderID string) error {
	if folderID == "" {
		return nil // 不在收藏夹中，不需要更新
	}

	objectID, err := r.ParseID(folderID)
	if err != nil {
		return fmt.Errorf("invalid folder ID: %w", err)
	}

	_, err = r.folderColl.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$inc": bson.M{"book_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to increment folder book count: %w", err)
	}

	return nil
}

// DecrementFolderBookCount 减少收藏夹书籍数量
func (r *MongoCollectionRepository) DecrementFolderBookCount(ctx context.Context, folderID string) error {
	if folderID == "" {
		return nil // 不在收藏夹中，不需要更新
	}

	objectID, err := r.ParseID(folderID)
	if err != nil {
		return fmt.Errorf("invalid folder ID: %w", err)
	}

	_, err = r.folderColl.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$inc": bson.M{"book_count": -1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to decrement folder book count: %w", err)
	}

	return nil
}

// =========================
// 公开收藏
// =========================

// GetPublicCollections 获取公开收藏列表
func (r *MongoCollectionRepository) GetPublicCollections(ctx context.Context, page, size int) ([]*social.Collection, int64, error) {
	filter := bson.M{"is_public": true}

	// 计算总数
	total, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count public collections: %w", err)
	}

	// 计算跳过数
	skip := int64((page - 1) * size)

	// 查询
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(size))

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find public collections: %w", err)
	}
	defer cursor.Close(ctx)

	var collections []*social.Collection
	if err := cursor.All(ctx, &collections); err != nil {
		return nil, 0, fmt.Errorf("failed to decode public collections: %w", err)
	}

	return collections, total, nil
}

// GetPublicFolders 获取公开收藏夹列表
func (r *MongoCollectionRepository) GetPublicFolders(ctx context.Context, page, size int) ([]*social.CollectionFolder, int64, error) {
	filter := bson.M{"is_public": true}

	// 计算总数
	total, err := r.folderColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count public folders: %w", err)
	}

	// 计算跳过数
	skip := int64((page - 1) * size)

	// 查询（按书籍数量排序）
	opts := options.Find().
		SetSort(bson.D{{Key: "book_count", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(size))

	cursor, err := r.folderColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find public folders: %w", err)
	}
	defer cursor.Close(ctx)

	var folders []*social.CollectionFolder
	if err := cursor.All(ctx, &folders); err != nil {
		return nil, 0, fmt.Errorf("failed to decode public folders: %w", err)
	}

	return folders, total, nil
}

// =========================
// 统计
// =========================

// CountUserCollections 统计用户收藏数
func (r *MongoCollectionRepository) CountUserCollections(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, fmt.Errorf("用户ID不能为空")
	}

	count, err := r.GetCollection().CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return 0, fmt.Errorf("failed to count user collections: %w", err)
	}

	return count, nil
}

// Health 健康检查
func (r *MongoCollectionRepository) Health(ctx context.Context) error {
	return r.GetDB().Client().Ping(ctx, nil)
}
