package mongodb

import (
	"context"
	"errors"
	"log"

	"Qingyu_backend/models/bookstore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BookStreamRepository 书籍流式仓储
type BookStreamRepository struct {
	baseRepo *MongoBookRepository
	cursorMgr *CursorManager
}

// NewBookStreamRepository 创建书籍流式仓储
func NewBookStreamRepository(baseRepo *MongoBookRepository) *BookStreamRepository {
	return &BookStreamRepository{
		baseRepo:  baseRepo,
		cursorMgr: NewCursorManager(),
	}
}

// StreamSearch 流式搜索书籍
func (r *BookStreamRepository) StreamSearch(ctx context.Context, filter *bookstore.BookFilter) (*mongo.Cursor, error) {
	if filter == nil {
		filter = &bookstore.BookFilter{}
	}

	// 构建查询条件
	query := r.buildQuery(filter)

	// 构建排序选项
	opts := r.buildOptions(filter)

	// 执行查询
	collection := r.baseRepo.GetCollection()
	cursor, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

// StreamByCursor 根据游标继续流式读取
func (r *BookStreamRepository) StreamByCursor(ctx context.Context, filter *bookstore.BookFilter) (*mongo.Cursor, error) {
	if filter == nil {
		filter = &bookstore.BookFilter{}
	}

	// 如果没有游标，使用普通的StreamSearch
	if filter.Cursor == nil || *filter.Cursor == "" {
		return r.StreamSearch(ctx, filter)
	}

	// 解析游标
	cursorFilter, err := r.buildCursorFilter(filter)
	if err != nil {
		return nil, err
	}

	// 构建基础查询条件
	query := r.buildQuery(filter)

	// 合并游标过滤条件
	if len(cursorFilter) > 0 {
		for k, v := range cursorFilter {
			query[k] = v
		}
	}

	// 构建排序选项
	opts := r.buildOptions(filter)

	// 执行查询
	collection := r.baseRepo.GetCollection()
	cursor, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

// buildQuery 构建MongoDB查询条件
func (r *BookStreamRepository) buildQuery(filter *bookstore.BookFilter) bson.M {
	query := bson.M{}

	// 状态过滤
	if filter.Status != nil {
		query["status"] = *filter.Status
	} else {
		// 默认只返回已发布的书籍
		query["status"] = bson.M{"$in": []string{"published", "ongoing", "completed"}}
	}

	// 分类过滤
	if filter.CategoryID != nil {
		objectID, err := primitive.ObjectIDFromHex(*filter.CategoryID)
		if err == nil {
			query["category_ids"] = objectID
		}
	}

	// 作者过滤
	if filter.Author != nil {
		query["author"] = *filter.Author
	}

	// 标签过滤 (ANY语义：只要匹配任一标签即可)
	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}

	// 推荐过滤
	if filter.IsRecommended != nil {
		query["is_recommended"] = *filter.IsRecommended
	}

	// 精选过滤
	if filter.IsFeatured != nil {
		query["is_featured"] = *filter.IsFeatured
	}

	// 热门过滤
	if filter.IsHot != nil {
		query["is_hot"] = *filter.IsHot
	}

	// 免费过滤
	if filter.IsFree != nil {
		query["is_free"] = *filter.IsFree
	}

	// 价格范围过滤
	if filter.MinPrice != nil || filter.MaxPrice != nil {
		priceQuery := bson.M{}
		if filter.MinPrice != nil {
			priceQuery["$gte"] = *filter.MinPrice
		}
		if filter.MaxPrice != nil {
			priceQuery["$lte"] = *filter.MaxPrice
		}
		query["price"] = priceQuery
	}

	// 关键词搜索
	if filter.Keyword != nil && *filter.Keyword != "" {
		keyword := *filter.Keyword
		orConditions := []bson.M{
			{"title": bson.M{"$regex": keyword, "$options": "i"}},
			{"author": bson.M{"$regex": keyword, "$options": "i"}},
			{"introduction": bson.M{"$regex": keyword, "$options": "i"}},
		}
		query["$or"] = orConditions
	}

	log.Printf("[DEBUG] StreamSearch query: %+v", query)
	return query
}

// buildOptions 构建查询选项
func (r *BookStreamRepository) buildOptions(filter *bookstore.BookFilter) *options.FindOptions {
	opts := options.Find()

	// 设置Limit
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	} else {
		opts.SetLimit(20) // 默认每批20条
	}

	// 设置排序
	sortField := filter.SortBy
	if sortField == "" {
		sortField = "created_at" // 默认按创建时间排序
	}

	sortOrder := -1 // 默认降序
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}

	opts.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	// 设置Offset (仅用于offset游标类型)
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}

	return opts
}

// buildCursorFilter 构建游标过滤条件
func (r *BookStreamRepository) buildCursorFilter(filter *bookstore.BookFilter) (bson.M, error) {
	if filter.Cursor == nil || *filter.Cursor == "" {
		return bson.M{}, nil
	}

	// 确定排序字段和方向
	sortField := filter.SortBy
	if sortField == "" {
		sortField = "created_at"
	}

	sortOrder := -1
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}

	// 使用游标管理器构建过滤条件
	cursorFilter, err := r.cursorMgr.BuildCursorFilter(*filter.Cursor, sortField, sortOrder)
	if err != nil {
		return nil, err
	}

	return cursorFilter, nil
}

// GetCursorManager 获取游标管理器
func (r *BookStreamRepository) GetCursorManager() *CursorManager {
	return r.cursorMgr
}

// CountWithFilter 使用过滤条件统计书籍数量
func (r *BookStreamRepository) CountWithFilter(ctx context.Context, filter *bookstore.BookFilter) (int64, error) {
	query := r.buildQuery(filter)
	collection := r.baseRepo.GetCollection()
	return collection.CountDocuments(ctx, query)
}

// StreamBatch 批量流式读取并处理
func (r *BookStreamRepository) StreamBatch(ctx context.Context, filter *bookstore.BookFilter, batchSize int, callback func(*bookstore.Book) error) error {
	if filter == nil {
		return errors.New("filter cannot be nil")
	}

	if batchSize <= 0 {
		batchSize = 20
	}

	// 设置batch size
	filter.Limit = batchSize

	cursor, err := r.StreamSearch(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var book bookstore.Book
		if err := cursor.Decode(&book); err != nil {
			return err
		}

		if err := callback(&book); err != nil {
			return err
		}
	}

	return cursor.Err()
}
