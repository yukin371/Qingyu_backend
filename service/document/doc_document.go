package document

import (
	"context"
	"errors"
	"time"

	"Qingyu_backend/global"
	model "Qingyu_backend/models/document"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DocumentService 处理文档相关业务逻辑
type DocumentService struct{}

func getCollection() *mongo.Collection {
	return global.DB.Collection("documents")
}

// Create 创建文档
func (s *DocumentService) Create(doc *model.Document) (*model.Document, error) {
	if doc == nil {
		return nil, errors.New("document is nil")
	}
	doc.ID = primitive.NewObjectID().Hex()
	doc.TouchForCreate()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := getCollection().InsertOne(ctx, doc)
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

	cursor, err := getCollection().Find(ctx, filter, findOptions)
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
		return nil, errors.New("id is empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var doc model.Document
	err := getCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// Update 更新文档（不更新 UserID 和 ID）
func (s *DocumentService) Update(id string, update *model.Document) (*model.Document, error) {
	if id == "" || update == nil {
		return nil, errors.New("invalid arguments")
	}
	update.TouchForUpdate()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateDoc := bson.M{
		"$set": bson.M{
			"title":      update.Title,
			"content":    update.Content,
			"tags":       update.Tags,
			"updated_at": update.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated model.Document
	err := getCollection().FindOneAndUpdate(ctx, bson.M{"_id": id}, updateDoc, opts).Decode(&updated)
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

	res, err := getCollection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
