package writer

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/dto"
	serviceInterfaces "Qingyu_backend/service/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type publicationMetadataDocument struct {
	CategoryID    string   `bson:"category_id,omitempty"`
	Tags          []string `bson:"tags,omitempty"`
	Description   string   `bson:"description,omitempty"`
	CoverImage    string   `bson:"cover_image,omitempty"`
	PublishType   string   `bson:"publish_type,omitempty"`
	Price         float64  `bson:"price,omitempty"`
	FreeChapters  int      `bson:"free_chapters,omitempty"`
	AuthorNote    string   `bson:"author_note,omitempty"`
	EnableComment bool     `bson:"enable_comment,omitempty"`
	EnableShare   bool     `bson:"enable_share,omitempty"`
	ChapterTitle  string   `bson:"chapter_title,omitempty"`
	ChapterNumber int      `bson:"chapter_number,omitempty"`
	IsFree        bool     `bson:"is_free,omitempty"`
}

type publicationRecordDocument struct {
	ID              primitive.ObjectID          `bson:"_id,omitempty"`
	ProjectID       string                      `bson:"project_id"`
	Type            string                      `bson:"type"`
	ResourceID      string                      `bson:"resource_id"`
	ResourceTitle   string                      `bson:"resource_title"`
	BookstoreID     string                      `bson:"bookstore_id"`
	BookstoreName   string                      `bson:"bookstore_name"`
	ExternalID      string                      `bson:"external_id,omitempty"`
	Status          string                      `bson:"status"`
	PublishTime     *time.Time                  `bson:"publish_time,omitempty"`
	ScheduledTime   *time.Time                  `bson:"scheduled_time,omitempty"`
	UnpublishTime   *time.Time                  `bson:"unpublish_time,omitempty"`
	UnpublishReason string                      `bson:"unpublish_reason,omitempty"`
	Metadata        publicationMetadataDocument `bson:"metadata,omitempty"`
	CreatedBy       string                      `bson:"created_by"`
	ReviewedBy      string                      `bson:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time                  `bson:"reviewed_at,omitempty"`
	ReviewNote      string                      `bson:"review_note,omitempty"`
	CreatedAt       time.Time                   `bson:"created_at"`
	UpdatedAt       time.Time                   `bson:"updated_at"`
}

// MongoPublicationRepository 发布记录 Mongo 仓储。
type MongoPublicationRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewMongoPublicationRepository(db *mongo.Database) *MongoPublicationRepository {
	return &MongoPublicationRepository{
		db:         db,
		collection: db.Collection("publication_records"),
	}
}

func (r *MongoPublicationRepository) Create(ctx context.Context, record *serviceInterfaces.PublicationRecord) error {
	doc, err := r.toDocument(ctx, record)
	if err != nil {
		return err
	}
	if doc.ID.IsZero() {
		doc.ID = primitive.NewObjectID()
	}

	if record.ID == "" {
		record.ID = doc.ID.Hex()
	}

	_, err = r.collection.InsertOne(ctx, doc)
	return err
}

func (r *MongoPublicationRepository) FindByID(ctx context.Context, id string) (*serviceInterfaces.PublicationRecord, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid publication record id: %w", err)
	}

	var doc publicationRecordDocument
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return r.toRecord(&doc), nil
}

func (r *MongoPublicationRepository) FindByProjectID(ctx context.Context, projectID string, page, pageSize int) ([]*serviceInterfaces.PublicationRecord, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	filter := bson.M{"project_id": projectID}
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var docs []publicationRecordDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}

	records := make([]*serviceInterfaces.PublicationRecord, 0, len(docs))
	for i := range docs {
		records = append(records, r.toRecord(&docs[i]))
	}
	return records, total, nil
}

func (r *MongoPublicationRepository) FindPending(ctx context.Context, page, pageSize int) ([]*serviceInterfaces.PublicationRecord, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	filter := bson.M{"status": serviceInterfaces.PublicationStatusPending}
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var docs []publicationRecordDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}

	records := make([]*serviceInterfaces.PublicationRecord, 0, len(docs))
	for i := range docs {
		records = append(records, r.toRecord(&docs[i]))
	}
	return records, total, nil
}

func (r *MongoPublicationRepository) FindByResourceID(ctx context.Context, resourceID string) (*serviceInterfaces.PublicationRecord, error) {
	// resource_id 存储为 string 类型，直接使用 string 查询
	var doc publicationRecordDocument
	if err := r.collection.FindOne(
		ctx,
		bson.M{"resource_id": resourceID},
		options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return r.toRecord(&doc), nil
}

func (r *MongoPublicationRepository) Update(ctx context.Context, record *serviceInterfaces.PublicationRecord) error {
	doc, err := r.toDocument(ctx, record)
	if err != nil {
		return err
	}
	if doc.ID.IsZero() {
		return fmt.Errorf("publication record id is required")
	}

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": doc.ID}, doc)
	return err
}

func (r *MongoPublicationRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid publication record id: %w", err)
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoPublicationRepository) FindPublishedByProjectID(ctx context.Context, projectID string) (*serviceInterfaces.PublicationRecord, error) {
	filter := bson.M{
		"project_id": projectID,
		"type":       "project",
		"status":     serviceInterfaces.PublicationStatusPublished,
	}

	var doc publicationRecordDocument
	if err := r.collection.FindOne(ctx, filter, options.FindOne().SetSort(bson.D{{Key: "updated_at", Value: -1}})).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return r.toRecord(&doc), nil
}

func (r *MongoPublicationRepository) toDocument(ctx context.Context, record *serviceInterfaces.PublicationRecord) (*publicationRecordDocument, error) {
	projectID, err := r.resolveProjectID(ctx, record)
	if err != nil {
		return nil, err
	}

	doc := &publicationRecordDocument{
		ProjectID:       projectID,
		Type:            record.Type,
		ResourceID:      record.ResourceID,
		ResourceTitle:   record.ResourceTitle,
		BookstoreID:     record.BookstoreID,
		BookstoreName:   record.BookstoreName,
		ExternalID:      record.ExternalID,
		Status:          record.Status,
		PublishTime:     record.PublishTime,
		ScheduledTime:   record.ScheduledTime,
		UnpublishTime:   record.UnpublishTime,
		UnpublishReason: record.UnpublishReason,
		Metadata: publicationMetadataDocument{
			CategoryID:    record.Metadata.CategoryID,
			Tags:          record.Metadata.Tags,
			Description:   record.Metadata.Description,
			CoverImage:    record.Metadata.CoverImage,
			PublishType:   record.Metadata.PublishType,
			Price:         record.Metadata.Price,
			FreeChapters:  record.Metadata.FreeChapters,
			AuthorNote:    record.Metadata.AuthorNote,
			EnableComment: record.Metadata.EnableComment,
			EnableShare:   record.Metadata.EnableShare,
			ChapterTitle:  record.Metadata.ChapterTitle,
			ChapterNumber: record.Metadata.ChapterNumber,
			IsFree:        record.Metadata.IsFree,
		},
		CreatedBy:  record.CreatedBy,
		ReviewedBy: record.ReviewedBy,
		ReviewedAt: record.ReviewedAt,
		ReviewNote: record.ReviewNote,
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.UpdatedAt,
	}

	if record.ID != "" {
		objectID, err := primitive.ObjectIDFromHex(record.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid publication record id: %w", err)
		}
		doc.ID = objectID
	}
	return doc, nil
}

func (r *MongoPublicationRepository) toRecord(doc *publicationRecordDocument) *serviceInterfaces.PublicationRecord {
	return &serviceInterfaces.PublicationRecord{
		ID:              doc.ID.Hex(),
		Type:            doc.Type,
		ResourceID:      doc.ResourceID,
		ResourceTitle:   doc.ResourceTitle,
		BookstoreID:     doc.BookstoreID,
		BookstoreName:   doc.BookstoreName,
		ExternalID:      doc.ExternalID,
		Status:          doc.Status,
		PublishTime:     doc.PublishTime,
		ScheduledTime:   doc.ScheduledTime,
		UnpublishTime:   doc.UnpublishTime,
		UnpublishReason: doc.UnpublishReason,
		Metadata: dto.PublicationMetadata{
			CategoryID:    doc.Metadata.CategoryID,
			Tags:          doc.Metadata.Tags,
			Description:   doc.Metadata.Description,
			CoverImage:    doc.Metadata.CoverImage,
			PublishType:   doc.Metadata.PublishType,
			Price:         doc.Metadata.Price,
			FreeChapters:  doc.Metadata.FreeChapters,
			AuthorNote:    doc.Metadata.AuthorNote,
			EnableComment: doc.Metadata.EnableComment,
			EnableShare:   doc.Metadata.EnableShare,
			ChapterTitle:  doc.Metadata.ChapterTitle,
			ChapterNumber: doc.Metadata.ChapterNumber,
			IsFree:        doc.Metadata.IsFree,
		},
		CreatedBy:  doc.CreatedBy,
		ReviewedBy: doc.ReviewedBy,
		ReviewedAt: doc.ReviewedAt,
		ReviewNote: doc.ReviewNote,
		CreatedAt:  doc.CreatedAt,
		UpdatedAt:  doc.UpdatedAt,
	}
}

func (r *MongoPublicationRepository) resolveProjectID(ctx context.Context, record *serviceInterfaces.PublicationRecord) (string, error) {
	if record.Type == "project" {
		return record.ResourceID, nil
	}

	docID, err := primitive.ObjectIDFromHex(record.ResourceID)
	if err != nil {
		return "", fmt.Errorf("invalid document id: %w", err)
	}

	var result struct {
		ProjectID primitive.ObjectID `bson:"project_id"`
	}
	if err := r.db.Collection("documents").FindOne(ctx, bson.M{"_id": docID}).Decode(&result); err != nil {
		return "", err
	}
	return result.ProjectID.Hex(), nil
}
