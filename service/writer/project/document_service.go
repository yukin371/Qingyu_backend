package project

import (
	model "Qingyu_backend/models/writer"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DocumentService 处理文档相关业务逻辑
type DocumentService struct {
	db *mongo.Database
}

// NewDocumentService 创建设置服务
func NewDocumentService(db *mongo.Database) *DocumentService {
	return &DocumentService{db: db}
}

// getCollection 获取数据库集合
func (s *DocumentService) getCollection() *mongo.Collection {
	return s.db.Collection("documents")
}

// Create 创建文档
func (s *DocumentService) Create(doc *model.Document) (*model.Document, error) {
	if doc == nil {
		return nil, errors.New("文档不能为空")
	}
	doc.ID = primitive.NewObjectID()
	doc.TouchForCreate()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.getCollection().InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// List 返回文档列表（可按用户过滤）
func (s *DocumentService) List(userID string, limit, offset int64) ([]*model.Document, error) {
	filter := bson.M{}
	if userID != "" {
		filter["user_id"] = userID
	}

	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(limit)
	}
	if offset > 0 {
		findOptions.SetSkip(offset)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.getCollection().Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []*model.Document
	for cursor.Next(ctx) {
		var d model.Document
		if err := cursor.Decode(&d); err != nil {
			return nil, err
		}
		docs = append(docs, &d)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return docs, nil
}

// GetByID 根据ID获取文档
func (s *DocumentService) GetByID(id string) (*model.Document, error) {
	if id == "" {
		return nil, errors.New("id为空")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var doc model.Document
	err := s.getCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// Update 更新文档元数据（不更新 UserID 和 ID）
// 注意：此方法只更新Document元数据，不更新内容
// 如需更新文档内容，请使用DocumentContentRepository
func (s *DocumentService) Update(id string, update *model.Document) (*model.Document, error) {
	if id == "" || update == nil {
		return nil, errors.New("id 或 update 为空")
	}
	update.TouchForUpdate()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 只更新元数据字段
	updateDoc := bson.M{
		"$set": bson.M{
			"title":      update.Title,
			"status":     update.Status,
			"word_count": update.WordCount,
			// 可以添加其他元数据字段
			"character_ids": update.CharacterIDs,
			"location_ids":  update.LocationIDs,
			"timeline_ids":  update.TimelineIDs,
			"plot_threads":  update.PlotThreads,
			"key_points":    update.KeyPoints,
			"writing_hints": update.WritingHints,
			"tags":          update.Tags,
			"notes":         update.Notes,
			"updated_at":    update.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated model.Document
	err := s.getCollection().FindOneAndUpdate(ctx, bson.M{"_id": id}, updateDoc, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

// Delete 删除文档
func (s *DocumentService) Delete(id string) (bool, error) {
	if id == "" {
		return false, errors.New("id is empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.getCollection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
