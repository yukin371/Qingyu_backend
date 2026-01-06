package storage

import (
	storageModel "Qingyu_backend/models/storage"
	"Qingyu_backend/repository/interfaces/shared"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoStorageRepository MongoDB文件存储Repository实现
type MongoStorageRepository struct {
	db                *mongo.Database
	filesCollection   *mongo.Collection
	accessCollection  *mongo.Collection
	uploadsCollection *mongo.Collection
}

// NewMongoStorageRepository 创建MongoDB Storage Repository
func NewMongoStorageRepository(db *mongo.Database) shared.StorageRepository {
	return &MongoStorageRepository{
		db:                db,
		filesCollection:   db.Collection("files"),
		accessCollection:  db.Collection("file_access"),
		uploadsCollection: db.Collection("multipart_uploads"),
	}
}

// ============ 文件元数据管理 ============

// CreateFile 创建文件记录
func (r *MongoStorageRepository) CreateFile(ctx context.Context, file *storageModel.FileInfo) error {
	if file == nil {
		return fmt.Errorf("file info cannot be nil")
	}

	now := time.Now()
	file.CreatedAt = now
	file.UpdatedAt = now

	if file.ID == "" {
		file.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.filesCollection.InsertOne(ctx, file)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}

// GetFile 根据ID获取文件信息
func (r *MongoStorageRepository) GetFile(ctx context.Context, fileID string) (*storageModel.FileInfo, error) {
	var file storageModel.FileInfo
	err := r.filesCollection.FindOne(ctx, bson.M{"_id": fileID}).Decode(&file)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("file not found: %s", fileID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return &file, nil
}

// GetFileByMD5 根据MD5哈希获取文件（用于去重）
func (r *MongoStorageRepository) GetFileByMD5(ctx context.Context, md5Hash string) (*storageModel.FileInfo, error) {
	var file storageModel.FileInfo
	err := r.filesCollection.FindOne(ctx, bson.M{"md5": md5Hash}).Decode(&file)

	if err == mongo.ErrNoDocuments {
		return nil, nil // 不存在不算错误
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get file by MD5: %w", err)
	}

	return &file, nil
}

// UpdateFile 更新文件信息
func (r *MongoStorageRepository) UpdateFile(ctx context.Context, fileID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.filesCollection.UpdateOne(
		ctx,
		bson.M{"_id": fileID},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("file not found: %s", fileID)
	}

	return nil
}

// DeleteFile 删除文件记录
func (r *MongoStorageRepository) DeleteFile(ctx context.Context, fileID string) error {
	result, err := r.filesCollection.DeleteOne(ctx, bson.M{"_id": fileID})

	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("file not found: %s", fileID)
	}

	return nil
}

// ListFiles 列表查询文件
func (r *MongoStorageRepository) ListFiles(ctx context.Context, filter *shared.FileFilter) ([]*storageModel.FileInfo, int64, error) {
	query := r.buildFileQuery(filter)

	// 计算总数
	total, err := r.filesCollection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count files: %w", err)
	}

	// 查询列表
	opts := options.Find()
	if filter != nil && filter.Page > 0 && filter.PageSize > 0 {
		opts.SetSkip(int64((filter.Page - 1) * filter.PageSize))
		opts.SetLimit(int64(filter.PageSize))
	} else if filter != nil && filter.Limit > 0 {
		opts.SetSkip(filter.Offset)
		opts.SetLimit(filter.Limit)
	}

	// 排序
	sortField := "created_at"
	if filter != nil && filter.SortBy != "" {
		sortField = filter.SortBy
	}
	sortOrder := -1 // 默认降序
	if filter != nil && filter.SortOrder == "asc" {
		sortOrder = 1
	}
	opts.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	cursor, err := r.filesCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list files: %w", err)
	}
	defer cursor.Close(ctx)

	var files []*storageModel.FileInfo
	if err = cursor.All(ctx, &files); err != nil {
		return nil, 0, fmt.Errorf("failed to decode files: %w", err)
	}

	return files, total, nil
}

// CountFiles 统计文件数量
func (r *MongoStorageRepository) CountFiles(ctx context.Context, filter *shared.FileFilter) (int64, error) {
	query := r.buildFileQuery(filter)
	return r.filesCollection.CountDocuments(ctx, query)
}

// buildFileQuery 构建文件查询条件
func (r *MongoStorageRepository) buildFileQuery(filter *shared.FileFilter) bson.M {
	query := bson.M{}

	if filter == nil {
		return query
	}

	if filter.UserID != "" {
		query["user_id"] = filter.UserID
	}
	if filter.Category != "" {
		query["category"] = filter.Category
	}
	if filter.FileType != "" {
		query["file_type"] = filter.FileType
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if filter.IsPublic != nil {
		query["is_public"] = *filter.IsPublic
	}
	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}
	if filter.Keyword != "" {
		query["$or"] = []bson.M{
			{"filename": bson.M{"$regex": filter.Keyword, "$options": "i"}},
			{"original_name": bson.M{"$regex": filter.Keyword, "$options": "i"}},
		}
	}
	if filter.StartDate != nil || filter.EndDate != nil {
		dateQuery := bson.M{}
		if filter.StartDate != nil {
			dateQuery["$gte"] = *filter.StartDate
		}
		if filter.EndDate != nil {
			dateQuery["$lte"] = *filter.EndDate
		}
		query["created_at"] = dateQuery
	}
	if filter.MinSize != nil || filter.MaxSize != nil {
		sizeQuery := bson.M{}
		if filter.MinSize != nil {
			sizeQuery["$gte"] = *filter.MinSize
		}
		if filter.MaxSize != nil {
			sizeQuery["$lte"] = *filter.MaxSize
		}
		query["size"] = sizeQuery
	}

	return query
}

// ============ 权限管理 ============

// GrantAccess 授予用户访问权限
func (r *MongoStorageRepository) GrantAccess(ctx context.Context, fileID, userID string) error {
	access := &storageModel.FileAccess{
		FileID:     fileID,
		UserID:     userID,
		Permission: storageModel.PermissionRead,
		GrantedAt:  time.Now(),
	}

	_, err := r.accessCollection.InsertOne(ctx, access)
	if err != nil {
		return fmt.Errorf("failed to grant access: %w", err)
	}

	return nil
}

// RevokeAccess 撤销用户访问权限
func (r *MongoStorageRepository) RevokeAccess(ctx context.Context, fileID, userID string) error {
	_, err := r.accessCollection.DeleteOne(ctx, bson.M{
		"file_id": fileID,
		"user_id": userID,
	})

	if err != nil {
		return fmt.Errorf("failed to revoke access: %w", err)
	}

	return nil
}

// CheckAccess 检查用户是否有访问权限
func (r *MongoStorageRepository) CheckAccess(ctx context.Context, fileID, userID string) (bool, error) {
	count, err := r.accessCollection.CountDocuments(ctx, bson.M{
		"file_id": fileID,
		"user_id": userID,
		"$or": []bson.M{
			{"expires_at": bson.M{"$exists": false}},
			{"expires_at": bson.M{"$gt": time.Now()}},
		},
	})

	if err != nil {
		return false, fmt.Errorf("failed to check access: %w", err)
	}

	return count > 0, nil
}

// ============ 分片上传管理 ============

// CreateMultipartUpload 创建分片上传任务
func (r *MongoStorageRepository) CreateMultipartUpload(ctx context.Context, upload *storageModel.MultipartUpload) error {
	if upload.ID == "" {
		upload.ID = primitive.NewObjectID().Hex()
	}
	if upload.UploadID == "" {
		upload.UploadID = primitive.NewObjectID().Hex()
	}
	upload.CreatedAt = time.Now()
	upload.UpdatedAt = time.Now()
	upload.Status = "pending"
	upload.UploadedChunks = []int{}

	_, err := r.uploadsCollection.InsertOne(ctx, upload)
	if err != nil {
		return fmt.Errorf("failed to create multipart upload: %w", err)
	}

	return nil
}

// GetMultipartUpload 获取分片上传任务
func (r *MongoStorageRepository) GetMultipartUpload(ctx context.Context, uploadID string) (*storageModel.MultipartUpload, error) {
	var upload storageModel.MultipartUpload
	err := r.uploadsCollection.FindOne(ctx, bson.M{"upload_id": uploadID}).Decode(&upload)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("multipart upload not found: %s", uploadID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get multipart upload: %w", err)
	}

	return &upload, nil
}

// UpdateMultipartUpload 更新分片上传任务
func (r *MongoStorageRepository) UpdateMultipartUpload(ctx context.Context, uploadID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.uploadsCollection.UpdateOne(
		ctx,
		bson.M{"upload_id": uploadID},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to update multipart upload: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("multipart upload not found: %s", uploadID)
	}

	return nil
}

// CompleteMultipartUpload 完成分片上传任务
func (r *MongoStorageRepository) CompleteMultipartUpload(ctx context.Context, uploadID string) error {
	now := time.Now()
	return r.UpdateMultipartUpload(ctx, uploadID, map[string]interface{}{
		"status":       "completed",
		"completed_at": now,
	})
}

// AbortMultipartUpload 中止分片上传任务
func (r *MongoStorageRepository) AbortMultipartUpload(ctx context.Context, uploadID string) error {
	return r.UpdateMultipartUpload(ctx, uploadID, map[string]interface{}{
		"status": "aborted",
	})
}

// ListMultipartUploads 列出用户的分片上传任务
func (r *MongoStorageRepository) ListMultipartUploads(ctx context.Context, userID string, status string) ([]*storageModel.MultipartUpload, error) {
	filter := bson.M{"uploaded_by": userID}
	if status != "" {
		filter["status"] = status
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.uploadsCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list multipart uploads: %w", err)
	}
	defer cursor.Close(ctx)

	var uploads []*storageModel.MultipartUpload
	if err = cursor.All(ctx, &uploads); err != nil {
		return nil, fmt.Errorf("failed to decode multipart uploads: %w", err)
	}

	return uploads, nil
}

// ============ 统计功能 ============

// IncrementDownloadCount 增加下载次数
func (r *MongoStorageRepository) IncrementDownloadCount(ctx context.Context, fileID string) error {
	_, err := r.filesCollection.UpdateOne(
		ctx,
		bson.M{"_id": fileID},
		bson.M{
			"$inc": bson.M{"downloads": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to increment download count: %w", err)
	}

	return nil
}

// ============ 健康检查 ============

// Health 健康检查
func (r *MongoStorageRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}
