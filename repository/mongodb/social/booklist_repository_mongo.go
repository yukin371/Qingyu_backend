package social

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/social"
)

// MongoBookListRepository MongoDB书单仓储实现
type MongoBookListRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoBookListRepository 创建MongoDB书单仓储
func NewMongoBookListRepository(db *mongo.Database) *MongoBookListRepository {
	return &MongoBookListRepository{
		db:         db,
		collection: db.Collection("book_lists"),
	}
}

// ========== 书单管理 ==========

// CreateBookList 创建书单
func (r *MongoBookListRepository) CreateBookList(ctx context.Context, bookList *social.BookList) error {
	bookList.CreatedAt = time.Now()
	bookList.UpdatedAt = time.Now()

	// 初始化Books字段为空切片，避免null值
	if bookList.Books == nil {
		bookList.Books = []social.BookListItem{}
	}

	result, err := r.collection.InsertOne(ctx, bookList)
	if err != nil {
		return fmt.Errorf("创建书单失败: %w", err)
	}

	// 获取插入后的ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		bookList.ID = oid
	}

	return nil
}

// GetBookListByID 根据ID获取书单
func (r *MongoBookListRepository) GetBookListByID(ctx context.Context, bookListID string) (*social.BookList, error) {
	objectID, err := primitive.ObjectIDFromHex(bookListID)
	if err != nil {
		return nil, fmt.Errorf("无效的书单ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	var bookList social.BookList

	err = r.collection.FindOne(ctx, filter).Decode(&bookList)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("获取书单失败: %w", err)
	}

	return &bookList, nil
}

// GetBookListsByUser 获取用户的书单列表
func (r *MongoBookListRepository) GetBookListsByUser(ctx context.Context, userID string, page, size int) ([]*social.BookList, int64, error) {
	filter := bson.M{"user_id": userID}

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计书单数量失败: %w", err)
	}

	// 分页查询
	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(size)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询书单列表失败: %w", err)
	}

	var bookLists []*social.BookList
	if err = cursor.All(ctx, &bookLists); err != nil {
		return nil, 0, fmt.Errorf("解析书单列表失败: %w", err)
	}

	return bookLists, total, nil
}

// GetPublicBookLists 获取公开书单列表
func (r *MongoBookListRepository) GetPublicBookLists(ctx context.Context, page, size int) ([]*social.BookList, int64, error) {
	filter := bson.M{"is_public": true}

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计公开书单数量失败: %w", err)
	}

	// 分页查询
	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(size)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询公开书单列表失败: %w", err)
	}

	var bookLists []*social.BookList
	if err = cursor.All(ctx, &bookLists); err != nil {
		return nil, 0, fmt.Errorf("解析公开书单列表失败: %w", err)
	}

	return bookLists, total, nil
}

// GetBookListsByCategory 根据分类获取书单
func (r *MongoBookListRepository) GetBookListsByCategory(ctx context.Context, category string, page, size int) ([]*social.BookList, int64, error) {
	filter := bson.M{
		"is_public": true,
		"category":  category,
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计分类书单数量失败: %w", err)
	}

	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(size)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询分类书单列表失败: %w", err)
	}

	var bookLists []*social.BookList
	if err = cursor.All(ctx, &bookLists); err != nil {
		return nil, 0, fmt.Errorf("解析分类书单列表失败: %w", err)
	}

	return bookLists, total, nil
}

// GetBookListsByTag 根据标签获取书单
func (r *MongoBookListRepository) GetBookListsByTag(ctx context.Context, tag string, page, size int) ([]*social.BookList, int64, error) {
	filter := bson.M{
		"is_public": true,
		"tags":      tag,
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计标签书单数量失败: %w", err)
	}

	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(size)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询标签书单列表失败: %w", err)
	}

	var bookLists []*social.BookList
	if err = cursor.All(ctx, &bookLists); err != nil {
		return nil, 0, fmt.Errorf("解析标签书单列表失败: %w", err)
	}

	return bookLists, total, nil
}

// SearchBookLists 搜索书单
func (r *MongoBookListRepository) SearchBookLists(ctx context.Context, keyword string, page, size int) ([]*social.BookList, int64, error) {
	filter := bson.M{
		"is_public": true,
		"$or": []bson.M{
			{"title": bson.M{"$regex": keyword, "$options": "i"}},
			{"description": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计搜索结果数量失败: %w", err)
	}

	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(size)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("搜索书单失败: %w", err)
	}

	var bookLists []*social.BookList
	if err = cursor.All(ctx, &bookLists); err != nil {
		return nil, 0, fmt.Errorf("解析搜索结果失败: %w", err)
	}

	return bookLists, total, nil
}

// UpdateBookList 更新书单
func (r *MongoBookListRepository) UpdateBookList(ctx context.Context, bookListID string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(bookListID)
	if err != nil {
		return fmt.Errorf("无效的书单ID: %w", err)
	}

	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新书单失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("书单不存在")
	}

	return nil
}

// DeleteBookList 删除书单
func (r *MongoBookListRepository) DeleteBookList(ctx context.Context, bookListID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookListID)
	if err != nil {
		return fmt.Errorf("无效的书单ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除书单失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("书单不存在")
	}

	return nil
}

// ========== 书单书籍管理 ==========

// AddBookToList 添加书籍到书单
func (r *MongoBookListRepository) AddBookToList(ctx context.Context, bookListID string, bookItem *social.BookListItem) error {
	objectID, err := primitive.ObjectIDFromHex(bookListID)
	if err != nil {
		return fmt.Errorf("无效的书单ID: %w", err)
	}

	bookItem.AddTime = time.Now()

	filter := bson.M{"_id": objectID}
	update := bson.M{"$push": bson.M{"books": bookItem}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("添加书籍到书单失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("书单不存在")
	}

	return nil
}

// RemoveBookFromList 从书单中移除书籍
func (r *MongoBookListRepository) RemoveBookFromList(ctx context.Context, bookListID, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookListID)
	if err != nil {
		return fmt.Errorf("无效的书单ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$pull": bson.M{"books": bson.M{"book_id": bookID}}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("从书单移除书籍失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("书单不存在")
	}

	return nil
}

// UpdateBookInList 更新书单中的书籍
func (r *MongoBookListRepository) UpdateBookInList(ctx context.Context, bookListID, bookID string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(bookListID)
	if err != nil {
		return fmt.Errorf("无效的书单ID: %w", err)
	}

	// 更新嵌套数组中的元素
	filter := bson.M{
		"_id":           objectID,
		"books.book_id": bookID,
	}

	update := bson.M{"$set": bson.M{}}
	for key, value := range updates {
		update["$set"].(bson.M)["books.$."+key] = value
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新书单书籍失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("书单或书籍不存在")
	}

	return nil
}

// ReorderBooks 重新排序书籍
func (r *MongoBookListRepository) ReorderBooks(ctx context.Context, bookListID string, bookOrders map[string]int) error {
	objectID, err := primitive.ObjectIDFromHex(bookListID)
	if err != nil {
		return fmt.Errorf("无效的书单ID: %w", err)
	}

	// 获取当前书单
	bookList, err := r.GetBookListByID(ctx, bookListID)
	if err != nil {
		return err
	}
	if bookList == nil {
		return errors.New("书单不存在")
	}

	// 更新排序
	for i := range bookList.Books {
		if newOrder, exists := bookOrders[bookList.Books[i].BookID]; exists {
			bookList.Books[i].Order = newOrder
		}
	}

	// 保存更新
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"books": bookList.Books, "updated_at": time.Now()}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("重新排序书籍失败: %w", err)
	}

	return nil
}

// GetBooksInList 获取书单中的书籍
func (r *MongoBookListRepository) GetBooksInList(ctx context.Context, bookListID string) ([]*social.BookListItem, error) {
	bookList, err := r.GetBookListByID(ctx, bookListID)
	if err != nil {
		return nil, err
	}
	if bookList == nil {
		return nil, errors.New("书单不存在")
	}

	// 转换为指针切片
	result := make([]*social.BookListItem, len(bookList.Books))
	for i := range bookList.Books {
		result[i] = &bookList.Books[i]
	}

	return result, nil
}

// ========== 书单点赞 ==========

// CreateBookListLike 创建书单点赞
func (r *MongoBookListRepository) CreateBookListLike(ctx context.Context, bookListLike *social.BookListLike) error {
	likesCollection := r.db.Collection("book_list_likes")
	bookListLike.CreatedAt = time.Now()

	result, err := likesCollection.InsertOne(ctx, bookListLike)
	if err != nil {
		return fmt.Errorf("创建点赞记录失败: %w", err)
	}

	// 获取插入后的ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		bookListLike.ID = oid
	}

	// 增加书单点赞数
	return r.IncrementBookListLikeCount(ctx, bookListLike.BookListID)
}

// DeleteBookListLike 删除书单点赞
func (r *MongoBookListRepository) DeleteBookListLike(ctx context.Context, bookListID, userID string) error {
	likesCollection := r.db.Collection("book_list_likes")

	filter := bson.M{
		"booklist_id": bookListID,
		"user_id":      userID,
	}

	result, err := likesCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除点赞记录失败: %w", err)
	}

	if result.DeletedCount > 0 {
		// 减少书单点赞数
		return r.DecrementBookListLikeCount(ctx, bookListID)
	}

	return nil
}

// GetBookListLike 获取书单点赞记录
func (r *MongoBookListRepository) GetBookListLike(ctx context.Context, bookListID, userID string) (*social.BookListLike, error) {
	likesCollection := r.db.Collection("book_list_likes")

	filter := bson.M{
		"booklist_id": bookListID,
		"user_id":      userID,
	}

	var bookListLike social.BookListLike
	err := likesCollection.FindOne(ctx, filter).Decode(&bookListLike)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("获取点赞记录失败: %w", err)
	}

	return &bookListLike, nil
}

// IsBookListLiked 检查是否已点赞
func (r *MongoBookListRepository) IsBookListLiked(ctx context.Context, bookListID, userID string) (bool, error) {
	likesCollection := r.db.Collection("book_list_likes")

	filter := bson.M{
		"booklist_id": bookListID,
		"user_id":      userID,
	}

	count, err := likesCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("检查点赞状态失败: %w", err)
	}

	return count > 0, nil
}

// GetBookListLikes 获取书单点赞列表
func (r *MongoBookListRepository) GetBookListLikes(ctx context.Context, bookListID string, page, size int) ([]*social.BookListLike, int64, error) {
	likesCollection := r.db.Collection("book_list_likes")

	filter := bson.M{"booklist_id": bookListID}

	total, err := likesCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计点赞数量失败: %w", err)
	}

	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(size)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := likesCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询点赞列表失败: %w", err)
	}

	var likes []*social.BookListLike
	if err = cursor.All(ctx, &likes); err != nil {
		return nil, 0, fmt.Errorf("解析点赞列表失败: %w", err)
	}

	return likes, total, nil
}

// IncrementBookListLikeCount 增加书单点赞数
func (r *MongoBookListRepository) IncrementBookListLikeCount(ctx context.Context, bookListID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookListID)
	if err != nil {
		return fmt.Errorf("无效的书单ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"like_count": 1}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("增加点赞数失败: %w", err)
	}

	return nil
}

// DecrementBookListLikeCount 减少书单点赞数
func (r *MongoBookListRepository) DecrementBookListLikeCount(ctx context.Context, bookListID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookListID)
	if err != nil {
		return fmt.Errorf("无效的书单ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"like_count": -1}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("减少点赞数失败: %w", err)
	}

	return nil
}

// ========== 书单复制 ==========

// ForkBookList 复制书单
func (r *MongoBookListRepository) ForkBookList(ctx context.Context, originalID, userID string) (*social.BookList, error) {
	// 获取原始书单
	original, err := r.GetBookListByID(ctx, originalID)
	if err != nil {
		return nil, err
	}
	if original == nil {
		return nil, errors.New("原始书单不存在")
	}

	// 创建新书单
	forked := &social.BookList{
		UserID:      userID,
		Title:       original.Title,
		Description: original.Description,
		Cover:       original.Cover,
		Category:    original.Category,
		Tags:        original.Tags,
		IsPublic:    original.IsPublic,
		Books:       original.Books,
		LikeCount:   0,
		ForkCount:   0,
		ViewCount:   0,
		OriginalID:  &original.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = r.CreateBookList(ctx, forked)
	if err != nil {
		return nil, err
	}

	// 增加原始书单的被复制次数
	err = r.IncrementForkCount(ctx, originalID)
	if err != nil {
		return nil, err
	}

	return forked, nil
}

// IncrementForkCount 增加被复制次数
func (r *MongoBookListRepository) IncrementForkCount(ctx context.Context, bookListID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookListID)
	if err != nil {
		return fmt.Errorf("无效的书单ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"fork_count": 1}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("增加被复制次数失败: %w", err)
	}

	return nil
}

// GetForkedBookLists 获取复制的书单列表
func (r *MongoBookListRepository) GetForkedBookLists(ctx context.Context, originalID string, page, size int) ([]*social.BookList, int64, error) {
	objectID, err := primitive.ObjectIDFromHex(originalID)
	if err != nil {
		return nil, 0, fmt.Errorf("无效的原始书单ID: %w", err)
	}

	filter := bson.M{
		"is_public":   true,
		"original_id": objectID,
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计复制书单数量失败: %w", err)
	}

	skip := int64((page - 1) * size)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(size)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询复制书单列表失败: %w", err)
	}

	var bookLists []*social.BookList
	if err = cursor.All(ctx, &bookLists); err != nil {
		return nil, 0, fmt.Errorf("解析复制书单列表失败: %w", err)
	}

	return bookLists, total, nil
}

// ========== 统计 ==========

// IncrementViewCount 增加浏览次数
func (r *MongoBookListRepository) IncrementViewCount(ctx context.Context, bookListID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookListID)
	if err != nil {
		return fmt.Errorf("无效的书单ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"view_count": 1}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("增加浏览次数失败: %w", err)
	}

	return nil
}

// CountUserBookLists 统计用户书单数
func (r *MongoBookListRepository) CountUserBookLists(ctx context.Context, userID string) (int64, error) {
	filter := bson.M{"user_id": userID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计用户书单数失败: %w", err)
	}

	return count, nil
}

// Health 健康检查
func (r *MongoBookListRepository) Health(ctx context.Context) error {
	// 执行一个简单的查询来检查连接
	_, err := r.collection.FindOne(ctx, bson.M{}).Raw()
	return err
}
