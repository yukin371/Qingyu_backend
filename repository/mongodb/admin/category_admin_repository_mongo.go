package admin

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/bookstore"
)

// CategoryAdminMongoRepository 分类管理 MongoDB 实现
type CategoryAdminMongoRepository struct {
	collection *mongo.Collection
}

// NewCategoryAdminMongoRepository 创建分类管理仓储实例
func NewCategoryAdminMongoRepository(db *mongo.Database) *CategoryAdminMongoRepository {
	return &CategoryAdminMongoRepository{
		collection: db.Collection("categories"),
	}
}

// UpdateBookCount 更新分类下的作品数量
func (r *CategoryAdminMongoRepository) UpdateBookCount(ctx context.Context, categoryID string, count int64) error {
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"book_count": count, "updated_at": time.Now()}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// BatchUpdateStatus 批量更新分类启用状态
func (r *CategoryAdminMongoRepository) BatchUpdateStatus(ctx context.Context, categoryIDs []string, isActive bool) error {
	objectIDs := make([]primitive.ObjectID, 0, len(categoryIDs))
	for _, id := range categoryIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		objectIDs = append(objectIDs, objectID)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	update := bson.M{"$set": bson.M{"is_active": isActive, "updated_at": time.Now()}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// GetDescendantIDs 获取所有子孙分类ID
func (r *CategoryAdminMongoRepository) GetDescendantIDs(ctx context.Context, categoryID string) ([]string, error) {
	var descendantIDs []string

	// 获取直接子分类
	filter := bson.M{"parent_id": categoryID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var children []bookstore.Category
	if err = cursor.All(ctx, &children); err != nil {
		return nil, err
	}

	for _, child := range children {
		descendantIDs = append(descendantIDs, child.ID)
		// 递归获取子分类的子孙
		subDescendants, err := r.GetDescendantIDs(ctx, child.ID)
		if err != nil {
			return nil, err
		}
		descendantIDs = append(descendantIDs, subDescendants...)
	}

	return descendantIDs, nil
}

// HasChildren 检查是否有子分类
func (r *CategoryAdminMongoRepository) HasChildren(ctx context.Context, categoryID string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"parent_id": categoryID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// NameExistsAtLevel 检查同级分类名称是否存在
func (r *CategoryAdminMongoRepository) NameExistsAtLevel(ctx context.Context, parentID *string, name string, excludeID string) (bool, error) {
	filter := bson.M{"name": name}

	if parentID == nil {
		// 顶级分类：parent_id 为 null 或不存在
		filter["$or"] = []bson.M{
			{"parent_id": bson.M{"$exists": false}},
			{"parent_id": nil},
		}
	} else {
		filter["parent_id"] = *parentID
	}

	// 排除当前分类（用于更新时校验）
	if excludeID != "" {
		excludeObjectID, err := primitive.ObjectIDFromHex(excludeID)
		if err != nil {
			return false, err
		}
		filter["_id"] = bson.M{"$ne": excludeObjectID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ============ 基础CRUD方法实现 ============

// Create 创建分类
func (r *CategoryAdminMongoRepository) Create(ctx context.Context, category *bookstore.Category) error {
	if category.ID == "" {
		category.ID = primitive.NewObjectID().Hex()
	}
	if category.CreatedAt.IsZero() {
		category.CreatedAt = time.Now()
	}
	if category.UpdatedAt.IsZero() {
		category.UpdatedAt = time.Now()
	}

	objectID, err := primitive.ObjectIDFromHex(category.ID)
	if err != nil {
		return err
	}

	category.ID = objectID.Hex()

	_, err = r.collection.InsertOne(ctx, category)
	return err
}

// GetByID 根据ID获取分类
func (r *CategoryAdminMongoRepository) GetByID(ctx context.Context, id string) (*bookstore.Category, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var category bookstore.Category
	filter := bson.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&category)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("category not found")
	}
	return &category, err
}

// Update 更新分类
func (r *CategoryAdminMongoRepository) Update(ctx context.Context, category *bookstore.Category) error {
	objectID, err := primitive.ObjectIDFromHex(category.ID)
	if err != nil {
		return err
	}

	category.UpdatedAt = time.Now()

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": category}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("category not found")
	}
	return nil
}

// Delete 删除分类
func (r *CategoryAdminMongoRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("category not found")
	}
	return nil
}

// List 获取分类列表
func (r *CategoryAdminMongoRepository) List(ctx context.Context, filter interface{}, opts ...interface{}) ([]*bookstore.Category, error) {
	var findOpts *options.FindOptions
	if len(opts) > 0 {
		if fo, ok := opts[0].(*options.FindOptions); ok {
			findOpts = fo
		}
	}

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*bookstore.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetTree 获取分类树
func (r *CategoryAdminMongoRepository) GetTree(ctx context.Context) ([]*bookstore.Category, error) {
	// 获取所有分类，按 level 和 sort_order 排序
	opts := options.Find().SetSort(bson.M{"level": 1, "sort_order": 1})
	allCategories, err := r.List(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	// 构建树形结构（返回顶级分类）
	var rootCategories []*bookstore.Category
	for _, cat := range allCategories {
		if cat.ParentID == nil {
			rootCategories = append(rootCategories, cat)
		}
	}

	return rootCategories, nil
}
