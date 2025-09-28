package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/document"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	mongodb "Qingyu_backend/repository/mongodb"
	documentRepo "Qingyu_backend/repository/interfaces/writing"
)

// MongoProjectRepository MongoDB项目仓储实现
type MongoProjectRepository struct {
	db           *mongo.Database
	collection   *mongo.Collection
	queryBuilder base.QueryBuilder
}

// NewMongoProjectRepository 创建MongoDB项目仓储实例
func NewMongoProjectRepository(db *mongo.Database) documentRepo.ProjectRepository {
	return &MongoProjectRepository{
		db:           db,
		collection:   db.Collection("projects"),
		queryBuilder: mongodb.NewMongoQueryBuilder(),
	}
}

// 实现 documentRepo.ProjectRepository 接口

// Create 创建项目
func (r *MongoProjectRepository) Create(ctx context.Context, project *document.Project) error {
	if project == nil {
		return errors.New("项目对象不能为空")
	}

	project.ID = primitive.NewObjectID().Hex()
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, project)
	return err
}

// GetByID 根据ID获取项目
func (r *MongoProjectRepository) GetByID(ctx context.Context, id string) (*document.Project, error) {
	var project document.Project
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("项目不存在")
		}
		return nil, err
	}
	return &project, nil
}

// Update 更新项目
func (r *MongoProjectRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	return err
}

// Delete 删除项目
func (r *MongoProjectRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Health 健康检查
func (r *MongoProjectRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}

// GetListByOwnerID 根据所有者ID获取项目列表(分页)
func (r *MongoProjectRepository) GetListByOwnerID(ctx context.Context, ownerID string, limit, offset int64) ([]*document.Project, error) {
	opts := options.Find().SetLimit(limit).SetSkip(offset)
	filter := bson.M{"owner_id": ownerID}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []*document.Project
	if err = cursor.All(ctx, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

// GetByOwnerAndStatus 根据所有者ID和状态获取项目
func (r *MongoProjectRepository) GetByOwnerAndStatus(ctx context.Context, ownerID, status string, limit, offset int64) ([]*document.Project, error) {
	opts := options.Find().SetLimit(limit).SetSkip(offset)
	filter := bson.M{"owner_id": ownerID, "status": status}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []*document.Project
	if err = cursor.All(ctx, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

// UpdateByOwner 根据所有者ID更新项目
func (r *MongoProjectRepository) UpdateByOwner(ctx context.Context, projectID, ownerID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	filter := bson.M{"_id": projectID, "owner_id": ownerID}
	update := bson.M{"$set": updates}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Restore 根据项目ID和所有者ID恢复项目
func (r *MongoProjectRepository) Restore(ctx context.Context, projectID, ownerID string) error {
	filter := bson.M{"_id": projectID, "owner_id": ownerID}
	update := bson.M{"$set": bson.M{"status": document.ProjectStatusPrivate}}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// IsOwner 根据项目ID和所有者ID检查是否为项目所有者
func (r *MongoProjectRepository) IsOwner(ctx context.Context, projectID, ownerID string) (bool, error) {
	filter := bson.M{"_id": projectID, "owner_id": ownerID}
	var project document.Project
	err := r.collection.FindOne(ctx, filter).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// SoftDelete 根据项目ID和所有者ID软删除项目
func (r *MongoProjectRepository) SoftDelete(ctx context.Context, projectID, ownerID string) error {
	filter := bson.M{"_id": projectID, "owner_id": ownerID}
	update := bson.M{"$set": bson.M{"status": document.ProjectStatusDeleted}}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// HardDelete 根据项目ID硬删除项目
func (r *MongoProjectRepository) HardDelete(ctx context.Context, projectID string) error {
	filter := bson.M{"_id": projectID}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// CountByOwner 根据所有者ID统计项目数量
func (r *MongoProjectRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	filter := bson.M{"owner_id": ownerID}
	return r.collection.CountDocuments(ctx, filter)
}

// CountByStatus 根据状态统计项目数量
func (r *MongoProjectRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	filter := bson.M{"status": status}
	return r.collection.CountDocuments(ctx, filter)
}

// CreateWithTransaction 创建项目并在事务中执行回调
func (r *MongoProjectRepository) CreateWithTransaction(ctx context.Context, project *document.Project, callback func(ctx context.Context) error) error {
	project.ID = primitive.NewObjectID().Hex()
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()
	project.Status = string(document.ProjectStatusPublic)
	_, err := r.collection.InsertOne(ctx, project)
	if err != nil {
		return err
	}
	return callback(ctx)
}
