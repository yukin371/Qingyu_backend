package recommendation

import (
	"context"
	"fmt"
	"time"

	reco "Qingyu_backend/models/recommendation"
	recoRepo "Qingyu_backend/repository/interfaces/recommendation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const recommendationTableCollection = "recommendation_tables"

type MongoTableRepository struct {
	col *mongo.Collection
}

func NewMongoTableRepository(db *mongo.Database) recoRepo.TableRepository {
	return &MongoTableRepository{col: db.Collection(recommendationTableCollection)}
}

func (r *MongoTableRepository) Create(ctx context.Context, table *reco.RecommendationTable) error {
	if table == nil {
		return fmt.Errorf("table cannot be nil")
	}
	now := time.Now()
	if table.CreatedAt.IsZero() {
		table.CreatedAt = now
	}
	table.UpdatedAt = now
	_, err := r.col.InsertOne(ctx, table)
	return err
}

func (r *MongoTableRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	if updates == nil {
		updates = map[string]interface{}{}
	}
	updates["updated_at"] = time.Now()
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": updates})
	return err
}

func (r *MongoTableRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *MongoTableRepository) GetByID(ctx context.Context, id string) (*reco.RecommendationTable, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var table reco.RecommendationTable
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&table); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &table, nil
}

func (r *MongoTableRepository) GetByTypePeriod(ctx context.Context, tableType reco.TableType, period string, source reco.TableSource) (*reco.RecommendationTable, error) {
	var table reco.RecommendationTable
	filter := bson.M{
		"type":   tableType,
		"period": period,
		"source": source,
	}
	if err := r.col.FindOne(ctx, filter).Decode(&table); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &table, nil
}

func (r *MongoTableRepository) List(ctx context.Context, tableType *reco.TableType, source *reco.TableSource, page, pageSize int) ([]*reco.RecommendationTable, int64, error) {
	filter := bson.M{}
	if tableType != nil {
		filter["type"] = *tableType
	}
	if source != nil {
		filter["source"] = *source
	}

	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tables []*reco.RecommendationTable
	if err := cursor.All(ctx, &tables); err != nil {
		return nil, 0, err
	}
	return tables, total, nil
}

func (r *MongoTableRepository) UpsertByTypePeriod(ctx context.Context, table *reco.RecommendationTable) error {
	if table == nil {
		return fmt.Errorf("table cannot be nil")
	}
	now := time.Now()
	if table.CreatedAt.IsZero() {
		table.CreatedAt = now
	}
	table.UpdatedAt = now

	filter := bson.M{
		"type":   table.Type,
		"period": table.Period,
		"source": table.Source,
	}
	update := bson.M{
		"$set": bson.M{
			"name":       table.Name,
			"status":     table.Status,
			"items":      table.Items,
			"metadata":   table.Metadata,
			"updated_by": table.UpdatedBy,
			"updated_at": table.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"created_at": table.CreatedAt,
		},
	}

	_, err := r.col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true)) // codeql[go/sql-injection]: MongoDB query, not SQL - IDs are validated ObjectIDs
	return err
}

func (r *MongoTableRepository) Health(ctx context.Context) error {
	return r.col.Database().Client().Ping(ctx, nil)
}
