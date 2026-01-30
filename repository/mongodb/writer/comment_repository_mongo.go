package writer

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/repository/mongodb/base"
	writerrepo "Qingyu_backend/repository/interfaces/writer"
)

const (
	CommentCollection = "document_comments"
)

type MongoCommentRepository struct {
	*base.BaseMongoRepository
}

// NewMongoCommentRepository 创建批注仓储
func NewMongoCommentRepository(db *mongo.Database) writerrepo.CommentRepository {
	return &MongoCommentRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(db, CommentCollection),
	}
}

// Create 创建批注
func (r *MongoCommentRepository) Create(ctx context.Context, comment *writer.DocumentComment) error {
	comment.ID = primitive.NewObjectID()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	// 如果是顶级评论，设置 threadID 为自身ID
	if comment.ParentID == nil && comment.ThreadID == nil {
		threadID := comment.ID
		comment.ThreadID = &threadID
	}

	_, err := r.GetCollection().InsertOne(ctx, comment)
	return err
}

// GetByID 根据ID获取批注
func (r *MongoCommentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*writer.DocumentComment, error) {
	var comment writer.DocumentComment
	filter := bson.M{"_id": id, "deleted_at": nil}

	err := r.GetCollection().FindOne(ctx, filter).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// Update 更新批注
func (r *MongoCommentRepository) Update(ctx context.Context, id primitive.ObjectID, comment *writer.DocumentComment) error {
	comment.UpdatedAt = time.Now()

	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{"$set": comment}

	_, err := r.GetCollection().UpdateOne(ctx, filter, update)
	return err
}

// Delete 删除批注（软删除）
func (r *MongoCommentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	_, err := r.GetCollection().UpdateOne(ctx, filter, update)
	return err
}

// HardDelete 硬删除批注
func (r *MongoCommentRepository) HardDelete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.GetCollection().DeleteOne(ctx, filter)
	return err
}

// List 查询批注列表（支持分页和筛选）
func (r *MongoCommentRepository) List(ctx context.Context, filter *writer.CommentFilter, page, pageSize int) ([]*writer.DocumentComment, int64, error) {
	mongoFilter := r.buildFilter(filter)

	// 获取总数
	total, err := r.GetCollection().CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.GetCollection().Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var comments []*writer.DocumentComment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// GetByDocument 获取文档的所有批注
func (r *MongoCommentRepository) GetByDocument(ctx context.Context, documentID primitive.ObjectID, includeResolved bool) ([]*writer.DocumentComment, error) {
	filter := bson.M{"document_id": documentID, "deleted_at": nil}
	if !includeResolved {
		filter["resolved"] = false
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*writer.DocumentComment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

// GetByChapter 获取章节的所有批注
func (r *MongoCommentRepository) GetByChapter(ctx context.Context, chapterID primitive.ObjectID, includeResolved bool) ([]*writer.DocumentComment, error) {
	filter := bson.M{"chapter_id": chapterID, "deleted_at": nil}
	if !includeResolved {
		filter["resolved"] = false
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*writer.DocumentComment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

// GetThread 获取批注线程
func (r *MongoCommentRepository) GetThread(ctx context.Context, threadID primitive.ObjectID) (*writer.CommentThread, error) {
	// 获取线程中的所有批注
	filter := bson.M{"thread_id": threadID, "deleted_at": nil}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var allComments []*writer.DocumentComment
	if err = cursor.All(ctx, &allComments); err != nil {
		return nil, err
	}

	if len(allComments) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	// 找到根评论（parent_id为nil的）
	var rootComment *writer.DocumentComment
	var replies []writer.DocumentComment

	for _, comment := range allComments {
		if comment.ParentID == nil {
			rootComment = comment
		} else {
			replies = append(replies, *comment)
		}
	}

	// 统计未解决的回复数
	unresolvedCount := 0
	for _, reply := range replies {
		if !reply.Resolved {
			unresolvedCount++
		}
	}

	thread := &writer.CommentThread{
		ThreadID:    threadID,
		RootComment: rootComment,
		Replies:     replies,
		ReplyCount:  len(replies),
		Unresolved:  unresolvedCount,
	}

	return thread, nil
}

// GetReplies 获取批注的回复
func (r *MongoCommentRepository) GetReplies(ctx context.Context, parentID primitive.ObjectID) ([]*writer.DocumentComment, error) {
	filter := bson.M{"parent_id": parentID, "deleted_at": nil}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var replies []*writer.DocumentComment
	if err = cursor.All(ctx, &replies); err != nil {
		return nil, err
	}

	return replies, nil
}

// MarkAsResolved 标记为已解决
func (r *MongoCommentRepository) MarkAsResolved(ctx context.Context, id primitive.ObjectID, resolvedBy primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"resolved":    true,
		"resolved_by": resolvedBy,
		"resolved_at": now,
		"updated_at":  now,
	}}

	_, err := r.GetCollection().UpdateOne(ctx, filter, update)
	return err
}

// MarkAsUnresolved 标记为未解决
func (r *MongoCommentRepository) MarkAsUnresolved(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"resolved":    false,
		"resolved_by": nil,
		"resolved_at": nil,
		"updated_at":  now,
	}}

	_, err := r.GetCollection().UpdateOne(ctx, filter, update)
	return err
}

// GetStats 获取批注统计
func (r *MongoCommentRepository) GetStats(ctx context.Context, documentID primitive.ObjectID) (*writer.CommentStats, error) {
	filter := bson.M{"document_id": documentID, "deleted_at": nil}

	// 总数
	totalCount, _ := r.GetCollection().CountDocuments(ctx, filter)

	// 已解决数
	resolvedFilter := bson.M{"document_id": documentID, "resolved": true, "deleted_at": nil}
	resolvedCount, _ := r.GetCollection().CountDocuments(ctx, resolvedFilter)

	// 未解决数
	unresolvedCount := totalCount - resolvedCount

	// 按类型统计
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id":   "$type",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, _ := r.GetCollection().Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)

	byType := make(map[string]int)
	var results []bson.M
	cursor.All(ctx, &results)
	for _, result := range results {
		typeStr := result["_id"].(string)
		count := result["count"].(int32)
		byType[typeStr] = int(count)
	}

	// 按用户统计
	pipeline = []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id":   "$user_id",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, _ = r.GetCollection().Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)

	byUser := make(map[string]int)
	results = nil
	cursor.All(ctx, &results)
	for _, result := range results {
		userID := r.IDToHex(result["_id"].(primitive.ObjectID))
		count := result["count"].(int32)
		byUser[userID] = int(count)
	}

	stats := &writer.CommentStats{
		TotalCount:      int(totalCount),
		ResolvedCount:   int(resolvedCount),
		UnresolvedCount: int(unresolvedCount),
		ByType:          byType,
		ByUser:          byUser,
	}

	return stats, nil
}

// GetUserComments 获取用户的批注
func (r *MongoCommentRepository) GetUserComments(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*writer.DocumentComment, int64, error) {
	filter := bson.M{"user_id": userID, "deleted_at": nil}

	total, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var comments []*writer.DocumentComment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// BatchDelete 批量删除批注
func (r *MongoCommentRepository) BatchDelete(ctx context.Context, ids []primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{"_id": bson.M{"$in": ids}, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	_, err := r.GetCollection().UpdateMany(ctx, filter, update)
	return err
}

// Search 搜索批注
func (r *MongoCommentRepository) Search(ctx context.Context, keyword string, documentID primitive.ObjectID, page, pageSize int) ([]*writer.DocumentComment, int64, error) {
	filter := bson.M{
		"document_id": documentID,
		"deleted_at":  nil,
		"content":     bson.M{"$regex": keyword, "$options": "i"},
	}

	total, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var comments []*writer.DocumentComment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// buildFilter 构建查询过滤器
func (r *MongoCommentRepository) buildFilter(filter *writer.CommentFilter) bson.M {
	mongoFilter := bson.M{"deleted_at": nil}

	if filter.DocumentID != nil {
		mongoFilter["document_id"] = *filter.DocumentID
	}
	if filter.ChapterID != nil {
		mongoFilter["chapter_id"] = *filter.ChapterID
	}
	if filter.UserID != nil {
		mongoFilter["user_id"] = *filter.UserID
	}
	if filter.Type != "" {
		mongoFilter["type"] = filter.Type
	}
	if filter.Resolved != nil {
		mongoFilter["resolved"] = *filter.Resolved
	}
	if filter.ParentID != nil {
		mongoFilter["parent_id"] = *filter.ParentID
	} else {
		// 如果没有指定ParentID，默认查询顶级评论
		mongoFilter["parent_id"] = nil
	}
	if filter.ThreadID != nil {
		mongoFilter["thread_id"] = *filter.ThreadID
	}
	if filter.StartDate != nil {
		mongoFilter["created_at"] = bson.M{"$gte": *filter.StartDate}
	}
	if filter.EndDate != nil {
		if _, exists := mongoFilter["created_at"]; exists {
			mongoFilter["created_at"].(bson.M)["$lte"] = *filter.EndDate
		} else {
			mongoFilter["created_at"] = bson.M{"$lte": *filter.EndDate}
		}
	}
	if filter.Keyword != "" {
		mongoFilter["content"] = bson.M{"$regex": filter.Keyword, "$options": "i"}
	}
	if len(filter.Labels) > 0 {
		mongoFilter["metadata.labels"] = bson.M{"$in": filter.Labels}
	}

	return mongoFilter
}
