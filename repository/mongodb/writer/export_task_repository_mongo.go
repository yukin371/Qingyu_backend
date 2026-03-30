package writer

import (
	"context"
	"fmt"
	"time"

	serviceInterfaces "Qingyu_backend/service/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// exportTaskDocument 导出任务MongoDB文档
type exportTaskDocument struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Type          string             `bson:"type"`
	ResourceID    string             `bson:"resource_id"`
	ResourceTitle string             `bson:"resource_title"`
	Format        string             `bson:"format"`
	Status        string             `bson:"status"`
	Progress      int                `bson:"progress"`
	FileSize      int64              `bson:"file_size,omitempty"`
	FileURL       string             `bson:"file_url,omitempty"`
	ExpiresAt     time.Time          `bson:"expires_at,omitempty"`
	ErrorMsg      string             `bson:"error_msg,omitempty"`
	CreatedBy     string             `bson:"created_by"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
	CompletedAt   *time.Time         `bson:"completed_at,omitempty"`
}

// MongoExportTaskRepository 导出任务MongoDB仓储
type MongoExportTaskRepository struct {
	collection *mongo.Collection
}

// NewMongoExportTaskRepository 创建导出任务仓储
func NewMongoExportTaskRepository(db *mongo.Database) *MongoExportTaskRepository {
	return &MongoExportTaskRepository{
		collection: db.Collection("export_tasks"),
	}
}

func (r *MongoExportTaskRepository) toModel(doc *exportTaskDocument) *serviceInterfaces.ExportTask {
	return &serviceInterfaces.ExportTask{
		ID:            doc.ID.Hex(),
		Type:          doc.Type,
		ResourceID:    doc.ResourceID,
		ResourceTitle: doc.ResourceTitle,
		Format:        doc.Format,
		Status:        doc.Status,
		Progress:      doc.Progress,
		FileSize:      doc.FileSize,
		FileURL:       doc.FileURL,
		ExpiresAt:     doc.ExpiresAt,
		ErrorMsg:      doc.ErrorMsg,
		CreatedBy:     doc.CreatedBy,
		CreatedAt:     doc.CreatedAt,
		UpdatedAt:     doc.UpdatedAt,
		CompletedAt:   doc.CompletedAt,
	}
}

func (r *MongoExportTaskRepository) toDoc(task *serviceInterfaces.ExportTask) (*exportTaskDocument, error) {
	var oid primitive.ObjectID
	if task.ID != "" {
		var err error
		oid, err = primitive.ObjectIDFromHex(task.ID)
		if err != nil {
			oid = primitive.NewObjectID()
		}
	} else {
		oid = primitive.NewObjectID()
	}

	return &exportTaskDocument{
		ID:            oid,
		Type:          task.Type,
		ResourceID:    task.ResourceID,
		ResourceTitle: task.ResourceTitle,
		Format:        task.Format,
		Status:        task.Status,
		Progress:      task.Progress,
		FileSize:      task.FileSize,
		FileURL:       task.FileURL,
		ExpiresAt:     task.ExpiresAt,
		ErrorMsg:      task.ErrorMsg,
		CreatedBy:     task.CreatedBy,
		CreatedAt:     task.CreatedAt,
		UpdatedAt:     task.UpdatedAt,
		CompletedAt:   task.CompletedAt,
	}, nil
}

func (r *MongoExportTaskRepository) Create(ctx context.Context, task *serviceInterfaces.ExportTask) error {
	doc, err := r.toDoc(task)
	if err != nil {
		return fmt.Errorf("convert task failed: %w", err)
	}
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = time.Now()
		doc.UpdatedAt = time.Now()
	}
	task.ID = doc.ID.Hex()

	_, err = r.collection.InsertOne(ctx, doc)
	return err
}

func (r *MongoExportTaskRepository) FindByID(ctx context.Context, id string) (*serviceInterfaces.ExportTask, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	var doc exportTaskDocument
	if err := r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return r.toModel(&doc), nil
}

func (r *MongoExportTaskRepository) FindByProjectID(ctx context.Context, projectID string, page, pageSize int) ([]*serviceInterfaces.ExportTask, int64, error) {
	filter := bson.M{"resource_id": projectID}
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var docs []exportTaskDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}

	tasks := make([]*serviceInterfaces.ExportTask, len(docs))
	for i, doc := range docs {
		tasks[i] = r.toModel(&doc)
	}
	return tasks, total, nil
}

func (r *MongoExportTaskRepository) Update(ctx context.Context, task *serviceInterfaces.ExportTask) error {
	doc, err := r.toDoc(task)
	if err != nil {
		return fmt.Errorf("convert task failed: %w", err)
	}
	doc.UpdatedAt = time.Now()

	update := bson.M{"$set": doc}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": doc.ID}, update)
	return err
}

func (r *MongoExportTaskRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *MongoExportTaskRepository) FindByUser(ctx context.Context, userID string, page, pageSize int) ([]*serviceInterfaces.ExportTask, int64, error) {
	filter := bson.M{"created_by": userID}
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var docs []exportTaskDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}

	tasks := make([]*serviceInterfaces.ExportTask, len(docs))
	for i, doc := range docs {
		tasks[i] = r.toModel(&doc)
	}
	return tasks, total, nil
}
