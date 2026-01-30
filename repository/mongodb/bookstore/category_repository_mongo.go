package mongodb

import (
	"Qingyu_backend/models/bookstore"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/repository/mongodb/base"

	BookstoreInterface "Qingyu_backend/repository/interfaces/bookstore"
	infra "Qingyu_backend/repository/interfaces/infrastructure"
)

// MongoCategoryRepository MongoDB分类仓储实现
type MongoCategoryRepository struct {
	*base.BaseMongoRepository
	client *mongo.Client
}

// NewMongoCategoryRepository 创建MongoDB分类仓储实例
func NewMongoCategoryRepository(client *mongo.Client, database string) BookstoreInterface.CategoryRepository {
	db := client.Database(database)
	return &MongoCategoryRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "categories"),
		client:              client,
	}
}

// Create 创建分类
func (r *MongoCategoryRepository) Create(ctx context.Context, category *bookstore.Category) error {
	if category == nil {
		return errors.New("category cannot be nil")
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	_, err := r.GetCollection().InsertOne(ctx, category)
	if err != nil {
		return err
	}

	// ID 已在调用方设置（或自动生成），无需从结果获取
	return nil
}

// GetByID 根据ID获取分类
func (r *MongoCategoryRepository) GetByID(ctx context.Context, id string) (*bookstore.Category, error) {
	objectID, err := r.ParseID(id)
	if err != nil {
		return nil, err
	}

	var category bookstore.Category
	err = r.GetCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// Update 更新分类
func (r *MongoCategoryRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return err
	}

	updates["updated_at"] = time.Now()

	result, err := r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("category not found")
	}

	return nil
}

// Delete 删除分类
func (r *MongoCategoryRepository) Delete(ctx context.Context, id string) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return err
	}

	result, err := r.GetCollection().DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("category not found")
	}

	return nil
}

// Health 健康检查
func (r *MongoCategoryRepository) Health(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}

// Count 统计分类数量（满足Filter接口）
func (r *MongoCategoryRepository) Count(ctx context.Context, filter infra.Filter) (int64, error) {
	var query bson.M
	if filter != nil {
		query = bson.M(filter.GetConditions())
	} else {
		query = bson.M{}
	}
	return r.GetCollection().CountDocuments(ctx, query)
}

// List 根据过滤条件列出分类
func (r *MongoCategoryRepository) List(ctx context.Context, filter infra.Filter) ([]*bookstore.Category, error) {
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
	}
	cursor, err := r.GetCollection().Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []*bookstore.Category
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Exists 判断分类是否存在
func (r *MongoCategoryRepository) Exists(ctx context.Context, id string) (bool, error) {
	objectID, err := r.ParseID(id)
	if err != nil {
		return false, err
	}

	count, err := r.GetCollection().CountDocuments(ctx, bson.M{"_id": objectID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetByName 根据名称获取分类
func (r *MongoCategoryRepository) GetByName(ctx context.Context, name string) (*bookstore.Category, error) {
	var category bookstore.Category
	err := r.GetCollection().FindOne(ctx, bson.M{"name": name}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// GetByParent 根据父分类获取子分类列表
func (r *MongoCategoryRepository) GetByParent(ctx context.Context, parentID string, limit, offset int) ([]*bookstore.Category, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "sort_order", Value: 1}, {Key: "created_at", Value: 1}})

	cursor, err := r.GetCollection().Find(ctx, bson.M{"parent_id": parentID}, opts)
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

// GetByLevel 根据层级获取分类列表
func (r *MongoCategoryRepository) GetByLevel(ctx context.Context, level int, limit, offset int) ([]*bookstore.Category, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "sort_order", Value: 1}, {Key: "created_at", Value: 1}})

	cursor, err := r.GetCollection().Find(ctx, bson.M{"level": level}, opts)
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

// GetRootCategories 获取根分类列表
func (r *MongoCategoryRepository) GetRootCategories(ctx context.Context) ([]*bookstore.Category, error) {
	opts := options.Find().SetSort(bson.D{{Key: "sort_order", Value: 1}, {Key: "created_at", Value: 1}})

	cursor, err := r.GetCollection().Find(ctx, bson.M{"level": 0, "is_active": true}, opts)
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

// GetCategoryTree 获取分类树
func (r *MongoCategoryRepository) GetCategoryTree(ctx context.Context) ([]*bookstore.CategoryTree, error) {
	// 获取所有激活的分类
	opts := options.Find().SetSort(bson.D{{Key: "level", Value: 1}, {Key: "sort_order", Value: 1}})
	cursor, err := r.GetCollection().Find(ctx, bson.M{"is_active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*bookstore.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	// 构建分类树
	categoryMap := make(map[string]*bookstore.CategoryTree)
	var roots []*bookstore.CategoryTree

	// 创建所有节点
	for _, cat := range categories {
		tree := &bookstore.CategoryTree{
			Category: *cat,
			Children: make([]*bookstore.CategoryTree, 0),
		}
		categoryMap[cat.ID] = tree

		if cat.Level == 0 {
			roots = append(roots, tree)
		}
	}

	// 建立父子关系
	for _, cat := range categories {
		if cat.ParentID != nil {
			if parent, exists := categoryMap[*cat.ParentID]; exists {
				parent.Children = append(parent.Children, categoryMap[cat.ID])
			}
		}
	}

	return roots, nil
}

// CountByParent 统计父分类下的子分类数量
func (r *MongoCategoryRepository) CountByParent(ctx context.Context, parentID string) (int64, error) {
	return r.GetCollection().CountDocuments(ctx, bson.M{"parent_id": parentID})
}

// UpdateBookCount 更新分类下的书籍数量
func (r *MongoCategoryRepository) UpdateBookCount(ctx context.Context, categoryID string, count int64) error {
	objectID, err := r.ParseID(categoryID)
	if err != nil {
		return err
	}

	_, err = r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"book_count": count,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// GetChildren 获取子分类
func (r *MongoCategoryRepository) GetChildren(ctx context.Context, parentID string) ([]*bookstore.Category, error) {
	opts := options.Find().SetSort(bson.D{{Key: "sort_order", Value: 1}, {Key: "created_at", Value: 1}})
	cursor, err := r.GetCollection().Find(ctx, bson.M{"parent_id": parentID, "is_active": true}, opts)
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

// GetAncestors 获取祖先分类
func (r *MongoCategoryRepository) GetAncestors(ctx context.Context, categoryID string) ([]*bookstore.Category, error) {
	var ancestors []*bookstore.Category
	currentID := categoryID

	for {
		objectID, err := r.ParseID(currentID)
		if err != nil {
			return nil, err
		}

		var category bookstore.Category
		err = r.GetCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&category)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				break
			}
			return nil, err
		}

		if category.ParentID == nil {
			ancestors = append([]*bookstore.Category{&category}, ancestors...)
			break
		}

		ancestors = append([]*bookstore.Category{&category}, ancestors...)
		currentID = *category.ParentID
	}

	return ancestors, nil
}

// GetDescendants 获取后代分类
func (r *MongoCategoryRepository) GetDescendants(ctx context.Context, categoryID string) ([]*bookstore.Category, error) {
	var descendants []*bookstore.Category

	// 递归获取所有后代
	var getChildren func(parentID string) error
	getChildren = func(parentID string) error {
		children, err := r.GetChildren(ctx, parentID)
		if err != nil {
			return err
		}

		for _, child := range children {
			descendants = append(descendants, child)
			if err := getChildren(child.ID); err != nil {
				return err
			}
		}

		return nil
	}

	if err := getChildren(categoryID); err != nil {
		return nil, err
	}

	return descendants, nil
}

// BatchUpdateStatus 批量更新状态
func (r *MongoCategoryRepository) BatchUpdateStatus(ctx context.Context, categoryIDs []string, isActive bool) error {
	objectIDs := make([]primitive.ObjectID, 0, len(categoryIDs))
	for _, id := range categoryIDs {
		objectID, err := r.ParseID(id)
		if err != nil {
			return err
		}
		objectIDs = append(objectIDs, objectID)
	}

	_, err := r.GetCollection().UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": objectIDs}},
		bson.M{
			"$set": bson.M{
				"is_active":  isActive,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// Transaction 执行事务
func (r *MongoCategoryRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
