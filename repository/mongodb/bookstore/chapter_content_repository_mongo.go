package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/bookstore"
	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

// MongoChapterContentRepository MongoDB 章节内容仓储实现
type MongoChapterContentRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
	database   *mongo.Database
}

// NewMongoChapterContentRepository 创建 MongoDB 章节内容仓储实例
func NewMongoChapterContentRepository(database *mongo.Database) BookstoreRepo.ChapterContentRepository {
	return &MongoChapterContentRepository{
		collection: database.Collection("chapter_contents"),
		client:     database.Client(),
		database:   database,
	}
}

// Create 创建章节内容
func (r *MongoChapterContentRepository) Create(ctx context.Context, content *bookstore.ChapterContent) error {
	if content == nil {
		return errors.New("content cannot be nil")
	}

	content.BeforeCreate()

	result, err := r.collection.InsertOne(ctx, content)
	if err != nil {
		return fmt.Errorf("failed to create chapter content: %w", err)
	}

	content.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 根据 ID 获取章节内容
func (r *MongoChapterContentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.ChapterContent, error) {
	var content bookstore.ChapterContent
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get chapter content: %w", err)
	}
	return &content, nil
}

// GetByChapterID 根据章节 ID 获取内容
func (r *MongoChapterContentRepository) GetByChapterID(ctx context.Context, chapterID primitive.ObjectID) (*bookstore.ChapterContent, error) {
	var content bookstore.ChapterContent
	err := r.collection.FindOne(ctx, bson.M{"chapter_id": chapterID}).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get chapter content by chapter_id: %w", err)
	}
	return &content, nil
}

// Update 更新章节内容
func (r *MongoChapterContentRepository) Update(ctx context.Context, content *bookstore.ChapterContent) error {
	if content == nil || content.ID.IsZero() {
		return errors.New("invalid content")
	}

	content.BeforeUpdate()

	filter := bson.M{"_id": content.ID}
	update := bson.M{"$set": content}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update chapter content: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("chapter content not found")
	}

	return nil
}

// Delete 删除章节内容
func (r *MongoChapterContentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete chapter content: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("chapter content not found")
	}

	return nil
}

// GetContentByChapterIDs 批量获取章节内容
func (r *MongoChapterContentRepository) GetContentByChapterIDs(ctx context.Context, chapterIDs []primitive.ObjectID) ([]*bookstore.ChapterContent, error) {
	if len(chapterIDs) == 0 {
		return []*bookstore.ChapterContent{}, nil
	}

	filter := bson.M{"chapter_id": bson.M{"$in": chapterIDs}}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapter contents: %w", err)
	}
	defer cursor.Close(ctx)

	var contents []*bookstore.ChapterContent
	if err = cursor.All(ctx, &contents); err != nil {
		return nil, fmt.Errorf("failed to decode chapter contents: %w", err)
	}

	return contents, nil
}

// BatchCreate 批量创建章节内容
func (r *MongoChapterContentRepository) BatchCreate(ctx context.Context, contents []*bookstore.ChapterContent) error {
	if len(contents) == 0 {
		return nil
	}

	docs := make([]interface{}, len(contents))
	for i, content := range contents {
		if content != nil {
			content.BeforeCreate()
			docs[i] = content
		}
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to batch create chapter contents: %w", err)
	}

	return nil
}

// BatchDelete 批量删除章节内容
func (r *MongoChapterContentRepository) BatchDelete(ctx context.Context, ids []primitive.ObjectID) error {
	if len(ids) == 0 {
		return nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}

	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to batch delete chapter contents: %w", err)
	}

	return nil
}

// GetLatestVersion 获取最新版本号
func (r *MongoChapterContentRepository) GetLatestVersion(ctx context.Context, chapterID primitive.ObjectID) (int, error) {
	options := options.Find().SetSort(bson.D{{"version", -1}}).SetLimit(1)
	cursor, err := r.collection.Find(ctx, bson.M{"chapter_id": chapterID}, options)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest version: %w", err)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var content bookstore.ChapterContent
		if err := cursor.Decode(&content); err == nil {
			return content.Version, nil
		}
	}

	return 0, nil // 没有内容，版本为 0
}

// UpdateContent 更新章节内容
func (r *MongoChapterContentRepository) UpdateContent(ctx context.Context, chapterID primitive.ObjectID, content string) (*bookstore.ChapterContent, error) {
	// 先查找现有内容
	existing, err := r.GetByChapterID(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("failed to find existing content: %w", err)
	}

	if existing == nil {
		// 创建新内容
		newContent := &bookstore.ChapterContent{
			ChapterID: chapterID,
			Content:   content,
			Format:    bookstore.ContentFormatMarkdown,
		}
		newContent.BeforeCreate()

		if err := r.Create(ctx, newContent); err != nil {
			return nil, fmt.Errorf("failed to create content: %w", err)
		}

		return newContent, nil
	}

	// 更新现有内容
	existing.UpdateContent(content)

	if err := r.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update content: %w", err)
	}

	return existing, nil
}

// Health 健康检查
func (r *MongoChapterContentRepository) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return r.client.Ping(ctx, nil)
}
