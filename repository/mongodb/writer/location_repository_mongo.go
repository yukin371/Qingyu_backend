package writer

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// LocationRepositoryMongo Location Repository的MongoDB实现
type LocationRepositoryMongo struct {
	db                 *mongo.Database
	locationCollection *mongo.Collection
	relationCollection *mongo.Collection
}

// NewLocationRepository 创建LocationRepository实例
func NewLocationRepository(db *mongo.Database) writerRepo.LocationRepository {
	return &LocationRepositoryMongo{
		db:                 db,
		locationCollection: db.Collection("locations"),
		relationCollection: db.Collection("location_relations"),
	}
}

// Create 创建地点
func (r *LocationRepositoryMongo) Create(ctx context.Context, location *writer.Location) error {
	if location.ID.IsZero() {
		location.ID = primitive.NewObjectID()
	}

	now := time.Now()
	location.CreatedAt = now
	location.UpdatedAt = now

	_, err := r.locationCollection.InsertOne(ctx, location)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create location failed", err)
	}

	return nil
}

// FindByID 根据ID查询地点
func (r *LocationRepositoryMongo) FindByID(ctx context.Context, locationID string) (*writer.Location, error) {
	var location writer.Location
	filter := bson.M{"_id": locationID}

	err := r.locationCollection.FindOne(ctx, filter).Decode(&location)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "location not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find location failed", err)
	}

	return &location, nil
}

// FindByProjectID 查询项目下的所有地点
func (r *LocationRepositoryMongo) FindByProjectID(ctx context.Context, projectID string) ([]*writer.Location, error) {
	filter := bson.M{"project_id": projectID}

	cursor, err := r.locationCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find locations failed", err)
	}
	defer cursor.Close(ctx)

	var locations []*writer.Location
	if err = cursor.All(ctx, &locations); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode locations failed", err)
	}

	return locations, nil
}

// FindByParentID 查询父地点下的子地点
func (r *LocationRepositoryMongo) FindByParentID(ctx context.Context, parentID string) ([]*writer.Location, error) {
	filter := bson.M{"parent_id": parentID}

	cursor, err := r.locationCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find child locations failed", err)
	}
	defer cursor.Close(ctx)

	var locations []*writer.Location
	if err = cursor.All(ctx, &locations); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode child locations failed", err)
	}

	return locations, nil
}

// Update 更新地点
func (r *LocationRepositoryMongo) Update(ctx context.Context, location *writer.Location) error {
	location.UpdatedAt = time.Now()

	filter := bson.M{"_id": location.ID}
	update := bson.M{"$set": location}

	result, err := r.locationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update location failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "location not found", nil)
	}

	return nil
}

// Delete 删除地点
func (r *LocationRepositoryMongo) Delete(ctx context.Context, locationID string) error {
	filter := bson.M{"_id": locationID}

	result, err := r.locationCollection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete location failed", err)
	}

	if result.DeletedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "location not found", nil)
	}

	return nil
}

// CreateRelation 创建地点关系
func (r *LocationRepositoryMongo) CreateRelation(ctx context.Context, relation *writer.LocationRelation) error {
	if relation.ID.IsZero() {
		relation.ID = primitive.NewObjectID()
	}

	now := time.Now()
	relation.CreatedAt = now
	relation.UpdatedAt = now

	_, err := r.relationCollection.InsertOne(ctx, relation)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create location relation failed", err)
	}

	return nil
}

// FindRelations 查询地点关系
func (r *LocationRepositoryMongo) FindRelations(ctx context.Context, projectID string, locationID *string) ([]*writer.LocationRelation, error) {
	filter := bson.M{"project_id": projectID}

	if locationID != nil {
		filter["$or"] = []bson.M{
			{"from_id": *locationID},
			{"to_id": *locationID},
		}
	}

	cursor, err := r.relationCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find location relations failed", err)
	}
	defer cursor.Close(ctx)

	var relations []*writer.LocationRelation
	if err = cursor.All(ctx, &relations); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode location relations failed", err)
	}

	return relations, nil
}

// DeleteRelation 删除地点关系
func (r *LocationRepositoryMongo) DeleteRelation(ctx context.Context, relationID string) error {
	filter := bson.M{"_id": relationID}

	result, err := r.relationCollection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete location relation failed", err)
	}

	if result.DeletedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "location relation not found", nil)
	}

	return nil
}

// ExistsByID 检查地点是否存在
func (r *LocationRepositoryMongo) ExistsByID(ctx context.Context, locationID string) (bool, error) {
	filter := bson.M{"_id": locationID}
	count, err := r.locationCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, errors.NewRepositoryError(errors.RepositoryErrorInternal, fmt.Sprintf("check location exists failed: %s", locationID), err)
	}
	return count > 0, nil
}

// CountByProjectID 统计项目下的地点数量
func (r *LocationRepositoryMongo) CountByProjectID(ctx context.Context, projectID string) (int64, error) {
	filter := bson.M{"project_id": projectID}
	count, err := r.locationCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, errors.NewRepositoryError(errors.RepositoryErrorInternal, "count locations failed", err)
	}
	return count, nil
}
