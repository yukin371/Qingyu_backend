package admin

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/bookstore"
)

var categoryNamePattern = regexp.MustCompile(`^[\p{L}\p{N}\p{Han}\s._\-()]{1,64}$`)

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
	if err := validateCategoryObjectID(categoryID); err != nil {
		return nil, err
	}

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
	if err := validateCategoryObjectID(categoryID); err != nil {
		return false, err
	}

	allCategories, err := r.List(ctx, bson.M{})
	if err != nil {
		return false, err
	}
	for _, cat := range allCategories {
		if cat.ParentID != nil && *cat.ParentID == categoryID {
			return true, nil
		}
	}
	return false, nil
}

// NameExistsAtLevel 检查同级分类名称是否存在
func (r *CategoryAdminMongoRepository) NameExistsAtLevel(ctx context.Context, parentID *string, name string, excludeID string) (bool, error) {
	safeName, err := sanitizeCategoryName(name)
	if err != nil {
		return false, err
	}

	filter := bson.M{"name": safeName}

	if parentID == nil {
		// 顶级分类：parent_id 为 null 或不存在
		filter["$or"] = []bson.M{
			{"parent_id": bson.M{"$exists": false}},
			{"parent_id": nil},
		}
	} else {
		if err := validateCategoryObjectID(*parentID); err != nil {
			return false, err
		}
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

	categories, err := r.List(ctx, bson.M{})
	if err != nil {
		return false, err
	}
	for _, cat := range categories {
		if cat == nil {
			continue
		}
		if cat.Name != safeName {
			continue
		}
		if parentID == nil {
			if cat.ParentID != nil {
				continue
			}
		} else {
			if cat.ParentID == nil || *cat.ParentID != *parentID {
				continue
			}
		}
		if excludeID != "" && cat.ID == excludeID {
			continue
		}
		return true, nil
	}
	return false, nil
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
	safeFilter, err := sanitizeCategoryListFilter(filter)
	if err != nil {
		return nil, err
	}

	var findOpts *options.FindOptions
	if len(opts) > 0 {
		if fo, ok := opts[0].(*options.FindOptions); ok {
			findOpts = fo
		}
	}

	cursor, err := r.collection.Find(ctx, bson.M{}, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var allCategories []*bookstore.Category
	if err = cursor.All(ctx, &allCategories); err != nil {
		return nil, err
	}

	return filterCategoriesInMemory(allCategories, safeFilter), nil
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

func validateCategoryObjectID(id string) error {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return fmt.Errorf("invalid category id: %w", err)
	}
	return nil
}

func sanitizeCategoryName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if !categoryNamePattern.MatchString(trimmed) {
		return "", errors.New("invalid category name format")
	}
	return trimmed, nil
}

func sanitizeCategoryListFilter(filter interface{}) (bson.M, error) {
	if filter == nil {
		return bson.M{}, nil
	}

	allowedKeys := map[string]struct{}{
		"_id":       {},
		"name":      {},
		"parent_id": {},
		"is_active": {},
		"level":     {},
	}

	out := bson.M{}
	var input map[string]interface{}

	switch f := filter.(type) {
	case bson.M:
		input = map[string]interface{}(f)
	case map[string]interface{}:
		input = f
	default:
		return nil, errors.New("invalid filter type")
	}

	for key, value := range input {
		if strings.HasPrefix(key, "$") {
			return nil, errors.New("operator filters are not allowed")
		}
		if _, ok := allowedKeys[key]; !ok {
			return nil, fmt.Errorf("unsupported filter key: %s", key)
		}

		switch key {
		case "_id", "parent_id":
			id, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("%s must be string", key)
			}
			if err := validateCategoryObjectID(id); err != nil {
				return nil, err
			}
			out[key] = id
		case "name":
			name, ok := value.(string)
			if !ok {
				return nil, errors.New("name must be string")
			}
			safeName, err := sanitizeCategoryName(name)
			if err != nil {
				return nil, err
			}
			out[key] = safeName
		case "is_active":
			active, ok := value.(bool)
			if !ok {
				return nil, errors.New("is_active must be bool")
			}
			out[key] = active
		case "level":
			switch lv := value.(type) {
			case int:
				out[key] = lv
			case int32:
				out[key] = int(lv)
			case int64:
				out[key] = int(lv)
			default:
				return nil, errors.New("level must be integer")
			}
		}
	}

	return out, nil
}

func filterCategoriesInMemory(categories []*bookstore.Category, filter bson.M) []*bookstore.Category {
	if len(filter) == 0 {
		return categories
	}

	result := make([]*bookstore.Category, 0, len(categories))
	for _, cat := range categories {
		if cat == nil {
			continue
		}
		if !matchesCategoryFilter(cat, filter) {
			continue
		}
		result = append(result, cat)
	}
	return result
}

func matchesCategoryFilter(cat *bookstore.Category, filter bson.M) bool {
	for key, value := range filter {
		switch key {
		case "_id":
			id, _ := value.(string)
			if cat.ID != id {
				return false
			}
		case "name":
			name, _ := value.(string)
			if cat.Name != name {
				return false
			}
		case "parent_id":
			parentID, _ := value.(string)
			if cat.ParentID == nil || *cat.ParentID != parentID {
				return false
			}
		case "is_active":
			active, _ := value.(bool)
			if cat.IsActive != active {
				return false
			}
		case "level":
			level, _ := value.(int)
			if cat.Level != level {
				return false
			}
		}
	}
	return true
}
