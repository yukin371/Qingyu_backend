package mongodb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/repository/mongodb/base"

	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

type chapterParagraph struct {
	Order   int
	Content string
	Format  string
}

// MongoChapterContentRepository MongoDB 章节内容仓储实现
type MongoChapterContentRepository struct {
	*base.BaseMongoRepository
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoChapterContentRepository 创建 MongoDB 章节内容仓储实例
func NewMongoChapterContentRepository(database *mongo.Database) BookstoreRepo.ChapterContentRepository {
	return &MongoChapterContentRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(database, "chapter_contents"),
		client:              database.Client(),
		database:            database,
	}
}

// Create 创建章节内容
func (r *MongoChapterContentRepository) Create(ctx context.Context, content *bookstore.ChapterContent) error {
	if content == nil {
		return errors.New("content cannot be nil")
	}

	content.BeforeCreate()

	result, err := r.GetCollection().InsertOne(ctx, content)
	if err != nil {
		return fmt.Errorf("failed to create chapter content: %w", err)
	}

	content.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 根据 ID 获取章节内容
func (r *MongoChapterContentRepository) GetByID(ctx context.Context, id string) (*bookstore.ChapterContent, error) {
	objectID, err := r.ParseID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid chapter content ID: %w", err)
	}

	var content bookstore.ChapterContent
	err = r.GetCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get chapter content: %w", err)
	}
	return &content, nil
}

// GetByChapterID 根据章节 ID 获取内容
func (r *MongoChapterContentRepository) GetByChapterID(ctx context.Context, chapterID string) (*bookstore.ChapterContent, error) {
	contents, err := r.ListByChapterID(ctx, chapterID)
	if err != nil {
		return nil, err
	}
	return aggregateChapterContents(contents), nil
}

// ListByChapterID 根据章节 ID 获取分段内容
func (r *MongoChapterContentRepository) ListByChapterID(ctx context.Context, chapterID string) ([]*bookstore.ChapterContent, error) {
	objectID, err := r.ParseID(chapterID)
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetSort(bson.D{
		{Key: "paragraph_order", Value: 1},
		{Key: "created_at", Value: 1},
		{Key: "_id", Value: 1},
	})

	cursor, err := r.GetCollection().Find(ctx, bson.M{"chapter_id": objectID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list chapter content by chapter_id: %w", err)
	}
	defer cursor.Close(ctx)

	var contents []*bookstore.ChapterContent
	if err := cursor.All(ctx, &contents); err != nil {
		return nil, fmt.Errorf("failed to decode chapter contents: %w", err)
	}

	return contents, nil
}

// Update 更新章节内容
func (r *MongoChapterContentRepository) Update(ctx context.Context, content *bookstore.ChapterContent) error {
	if content == nil || content.ID.IsZero() {
		return errors.New("invalid content")
	}

	content.BeforeUpdate()

	filter := bson.M{"_id": content.ID}
	update := bson.M{"$set": content}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update chapter content: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("chapter content not found")
	}

	return nil
}

// Delete 删除章节内容
func (r *MongoChapterContentRepository) Delete(ctx context.Context, id string) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return err
	}

	result, err := r.GetCollection().DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete chapter content: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("chapter content not found")
	}

	return nil
}

// GetContentByChapterIDs 批量获取章节内容
func (r *MongoChapterContentRepository) GetContentByChapterIDs(ctx context.Context, chapterIDs []string) ([]*bookstore.ChapterContent, error) {
	if len(chapterIDs) == 0 {
		return []*bookstore.ChapterContent{}, nil
	}

	// 使用基类方法批量转换 string IDs to ObjectIDs
	objectIDs, err := r.ParseIDs(chapterIDs)
	if err != nil {
		return nil, fmt.Errorf("invalid chapter IDs: %w", err)
	}

	filter := bson.M{"chapter_id": bson.M{"$in": objectIDs}}

	cursor, err := r.GetCollection().Find(ctx, filter)
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
			if content.ID.IsZero() {
				content.ID = primitive.NewObjectID()
			}
			content.BeforeCreate()
			docs[i] = content
		}
	}

	_, err := r.GetCollection().InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to batch create chapter contents: %w", err)
	}

	return nil
}

// BatchDelete 批量删除章节内容
func (r *MongoChapterContentRepository) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	// 使用基类方法批量转换 string ID 为 ObjectID
	objectIDs, err := r.ParseIDs(ids)
	if err != nil {
		return fmt.Errorf("failed to parse IDs: %w", err)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}

	_, err = r.GetCollection().DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to batch delete chapter contents: %w", err)
	}

	return nil
}

// GetLatestVersion 获取最新版本号
func (r *MongoChapterContentRepository) GetLatestVersion(ctx context.Context, chapterID string) (int, error) {
	objectID, err := r.ParseID(chapterID)
	if err != nil {
		return 0, fmt.Errorf("invalid chapter ID: %w", err)
	}
	options := options.Find().SetSort(bson.D{{Key: "version", Value: -1}}).SetLimit(1)
	cursor, err := r.GetCollection().Find(ctx, bson.M{"chapter_id": objectID}, options)
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
func (r *MongoChapterContentRepository) UpdateContent(ctx context.Context, chapterID string, content string) (*bookstore.ChapterContent, error) {
	objectID, err := r.ParseID(chapterID)
	if err != nil {
		return nil, fmt.Errorf("invalid chapter ID: %w", err)
	}
	existingContents, err := r.ListByChapterID(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("failed to find existing content: %w", err)
	}

	paragraphs := splitLegacyContent(content)
	now := time.Now()
	nextVersion := 1
	if len(existingContents) > 0 {
		latestVersion, latestErr := r.GetLatestVersion(ctx, chapterID)
		if latestErr != nil {
			return nil, fmt.Errorf("failed to get latest version: %w", latestErr)
		}
		nextVersion = latestVersion + 1
		if _, err := r.GetCollection().DeleteMany(ctx, bson.M{"chapter_id": objectID}); err != nil {
			return nil, fmt.Errorf("failed to delete previous chapter contents: %w", err)
		}
	}

	docs := make([]interface{}, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		wordCount := len([]rune(paragraph.Content))
		doc := &bookstore.ChapterContent{
			ID:             primitive.NewObjectID(),
			ChapterID:      objectID,
			Content:        paragraph.Content,
			Format:         paragraph.Format,
			Version:        nextVersion,
			ParagraphOrder: paragraph.Order,
			WordCount:      wordCount,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		doc.Hash = doc.CalculateHash()
		docs = append(docs, doc)
	}

	if len(docs) > 0 {
		if _, err := r.GetCollection().InsertMany(ctx, docs); err != nil {
			return nil, fmt.Errorf("failed to replace chapter contents: %w", err)
		}
	}

	aggregate := aggregateParagraphs(objectID, paragraphs, nextVersion, now)
	if aggregate != nil {
		aggregate.Hash = aggregate.CalculateHash()
	}
	return aggregate, nil
}

// Health 健康检查
func (r *MongoChapterContentRepository) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return r.client.Ping(ctx, nil)
}

func aggregateChapterContents(contents []*bookstore.ChapterContent) *bookstore.ChapterContent {
	if len(contents) == 0 {
		return nil
	}

	paragraphs := make([]chapterParagraph, 0, len(contents))
	chapterID := contents[0].ChapterID
	version := contents[0].Version
	createdAt := contents[0].CreatedAt
	updatedAt := contents[0].UpdatedAt
	for idx, item := range contents {
		order := item.ParagraphOrder
		if order <= 0 {
			order = idx + 1
		}
		paragraphs = append(paragraphs, chapterParagraph{
			Order:   order,
			Content: item.Content,
			Format:  item.Format,
		})
		if item.Version > version {
			version = item.Version
		}
		if item.CreatedAt.Before(createdAt) {
			createdAt = item.CreatedAt
		}
		if item.UpdatedAt.After(updatedAt) {
			updatedAt = item.UpdatedAt
		}
	}

	aggregate := aggregateParagraphs(chapterID, paragraphs, version, updatedAt)
	if aggregate != nil {
		aggregate.CreatedAt = createdAt
	}
	return aggregate
}

func aggregateParagraphs(chapterID primitive.ObjectID, paragraphs []chapterParagraph, version int, ts time.Time) *bookstore.ChapterContent {
	if len(paragraphs) == 0 {
		return nil
	}

	contentParts := make([]string, 0, len(paragraphs))
	wordCount := 0
	format := paragraphs[0].Format
	for _, paragraph := range paragraphs {
		contentParts = append(contentParts, paragraph.Content)
		wordCount += len([]rune(paragraph.Content))
		if format == "" && paragraph.Format != "" {
			format = paragraph.Format
		}
	}

	aggregate := &bookstore.ChapterContent{
		ChapterID: chapterID,
		Content:   strings.Join(contentParts, "\n\n"),
		Format:    format,
		Version:   version,
		WordCount: wordCount,
		CreatedAt: ts,
		UpdatedAt: ts,
	}
	aggregate.Hash = aggregate.CalculateHash()
	return aggregate
}

func splitLegacyContent(content string) []chapterParagraph {
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	normalized = strings.TrimSpace(normalized)
	if normalized == "" {
		return []chapterParagraph{{
			Order:   1,
			Content: "",
			Format:  bookstore.ContentFormatMarkdown,
		}}
	}

	blocks := strings.Split(normalized, "\n\n")
	paragraphs := make([]chapterParagraph, 0, len(blocks))
	order := 1
	for _, block := range blocks {
		paragraph := strings.TrimSpace(block)
		if paragraph == "" {
			continue
		}
		paragraphs = append(paragraphs, chapterParagraph{
			Order:   order,
			Content: paragraph,
			Format:  bookstore.ContentFormatMarkdown,
		})
		order++
	}

	if len(paragraphs) == 0 {
		paragraphs = append(paragraphs, chapterParagraph{
			Order:   1,
			Content: normalized,
			Format:  bookstore.ContentFormatMarkdown,
		})
	}

	return paragraphs
}
