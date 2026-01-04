package reader

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/reader"
	readerRepo "Qingyu_backend/repository/interfaces/reader"
)

// BookmarkMongoRepository 书签MongoDB仓储实现
type BookmarkMongoRepository struct {
	collection *mongo.Collection
}

// NewBookmarkMongoRepository 创建书签仓储
func NewBookmarkMongoRepository(db *mongo.Database) readerRepo.BookmarkRepository {
	return &BookmarkMongoRepository{
		collection: db.Collection("bookmarks"),
	}
}

// Create 创建书签
func (r *BookmarkMongoRepository) Create(ctx context.Context, bookmark *reader.Bookmark) error {
	bookmark.ID = primitive.NewObjectID()
	bookmark.CreatedAt = time.Now()
	bookmark.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, bookmark)
	return err
}

// GetByID 根据ID获取书签
func (r *BookmarkMongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*reader.Bookmark, error) {
	var bookmark reader.Bookmark
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bookmark)
	if err != nil {
		return nil, err
	}
	return &bookmark, nil
}

// GetByUser 获取用户的所有书签
func (r *BookmarkMongoRepository) GetByUser(ctx context.Context, userID primitive.ObjectID, filter *reader.BookmarkFilter, page, size int) ([]*reader.Bookmark, int64, error) {
	// 构建查询条件
	query := bson.M{"user_id": userID}

	if filter != nil {
		if filter.BookID != nil {
			query["book_id"] = *filter.BookID
		}
		if filter.ChapterID != nil {
			query["chapter_id"] = *filter.ChapterID
		}
		if filter.Color != "" {
			query["color"] = filter.Color
		}
		if filter.Tag != "" {
			query["tags"] = bson.M{"$in": []string{filter.Tag}}
		}
		if filter.IsPublic != nil {
			query["is_public"] = *filter.IsPublic
		}
		if filter.DateFrom != nil {
			query["created_at"] = bson.M{"$gte": *filter.DateFrom}
		}
		if filter.DateTo != nil {
			if query["created_at"] == nil {
				query["created_at"] = bson.M{}
			}
			query["created_at"].(bson.M)["$lte"] = *filter.DateTo
		}
	}

	// 获取总数
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(size)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var bookmarks []*reader.Bookmark
	if err = cursor.All(ctx, &bookmarks); err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}

// GetByBook 获取某本书的书签
func (r *BookmarkMongoRepository) GetByBook(ctx context.Context, userID, bookID primitive.ObjectID, page, size int) ([]*reader.Bookmark, int64, error) {
	filter := &reader.BookmarkFilter{BookID: &bookID}
	return r.GetByUser(ctx, userID, filter, page, size)
}

// GetByChapter 获取某章节的书签
func (r *BookmarkMongoRepository) GetByChapter(ctx context.Context, userID, chapterID primitive.ObjectID, page, size int) ([]*reader.Bookmark, int64, error) {
	filter := &reader.BookmarkFilter{ChapterID: &chapterID}
	return r.GetByUser(ctx, userID, filter, page, size)
}

// Update 更新书签
func (r *BookmarkMongoRepository) Update(ctx context.Context, bookmark *reader.Bookmark) error {
	bookmark.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"note":       bookmark.Note,
			"color":      bookmark.Color,
			"quote":      bookmark.Quote,
			"is_public":  bookmark.IsPublic,
			"tags":       bookmark.Tags,
			"position":   bookmark.Position,
			"updated_at": bookmark.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": bookmark.ID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("bookmark not found")
	}
	return nil
}

// Delete 删除书签
func (r *BookmarkMongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("bookmark not found")
	}
	return nil
}

// DeleteByBook 删除某本书的所有书签
func (r *BookmarkMongoRepository) DeleteByBook(ctx context.Context, userID, bookID primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{
		"user_id": userID,
		"book_id": bookID,
	})
	return err
}

// Exists 检查书签是否存在
func (r *BookmarkMongoRepository) Exists(ctx context.Context, userID, chapterID primitive.ObjectID, position int) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"user_id":    userID,
		"chapter_id": chapterID,
		"position":   position,
	})
	return count > 0, err
}

// GetStats 获取书签统计
func (r *BookmarkMongoRepository) GetStats(ctx context.Context, userID primitive.ObjectID) (*reader.BookmarkStats, error) {
	stats := &reader.BookmarkStats{
		ByColor: make(map[string]int64),
		ByBook:  make(map[string]int64),
	}

	// 总数
	total, err := r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	stats.TotalCount = total

	// 按颜色分组
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{"_id": "$color", "count": bson.M{"$sum": 1}}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err == nil {
		var results []struct {
			Color string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		cursor.All(ctx, &results)
		for _, r := range results {
			stats.ByColor[r.Color] = r.Count
		}
		cursor.Close(ctx)
	}

	// 按书籍分组
	pipeline = []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{"_id": "$book_id", "count": bson.M{"$sum": 1}}},
	}
	cursor, err = r.collection.Aggregate(ctx, pipeline)
	if err == nil {
		var results []struct {
			BookID string `bson:"_id"`
			Count  int64  `bson:"count"`
		}
		cursor.All(ctx, &results)
		for _, r := range results {
			stats.ByBook[r.BookID] = r.Count
		}
		cursor.Close(ctx)
	}

	// 公开/私有统计
	publicCount, _ := r.collection.CountDocuments(ctx, bson.M{"user_id": userID, "is_public": true})
	stats.PublicCount = publicCount
	stats.PrivateCount = total - publicCount

	// 本月/本周统计
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))

	stats.ThisMonthCount, _ = r.collection.CountDocuments(ctx, bson.M{
		"user_id":    userID,
		"created_at": bson.M{"$gte": monthStart},
	})
	stats.ThisWeekCount, _ = r.collection.CountDocuments(ctx, bson.M{
		"user_id":    userID,
		"created_at": bson.M{"$gte": weekStart},
	})

	// 最近书签
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(10)
	cursor, err = r.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err == nil {
		var bookmarks []reader.Bookmark
		cursor.All(ctx, &bookmarks)
		stats.RecentBookmarks = bookmarks
		cursor.Close(ctx)
	}

	return stats, nil
}

// GetPublicBookmarks 获取公开书签
func (r *BookmarkMongoRepository) GetPublicBookmarks(ctx context.Context, page, size int) ([]*reader.Bookmark, int64, error) {
	query := bson.M{"is_public": true}

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(size)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var bookmarks []*reader.Bookmark
	if err = cursor.All(ctx, &bookmarks); err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}

// Search 搜索书签
func (r *BookmarkMongoRepository) Search(ctx context.Context, userID primitive.ObjectID, keyword string, page, size int) ([]*reader.Bookmark, int64, error) {
	query := bson.M{
		"user_id": userID,
		"$or": []bson.M{
			{"note": bson.M{"$regex": keyword, "$options": "i"}},
			{"quote": bson.M{"$regex": keyword, "$options": "i"}},
			{"tags": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(size)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var bookmarks []*reader.Bookmark
	if err = cursor.All(ctx, &bookmarks); err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}
