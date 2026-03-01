package shared

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"Qingyu_backend/models/users"
	"Qingyu_backend/repository"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	userrepo "Qingyu_backend/repository/interfaces/user"
)

const UserCollection = "users"

// MongoUserRepository 统一的MongoDB用户仓储实现
// 实现 BaseUserRepository 接口，供 admin 和 user 模块共享使用
type MongoUserRepository struct {
	db *mongo.Database
}

// NewMongoUserRepository 创建统一的用户仓储实例
func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{db: db}
}

// === 实现 base.CRUDRepository 接口 ===

func (r *MongoUserRepository) Create(ctx context.Context, user *users.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.db.Collection(UserCollection).InsertOne(ctx, user)
	if err != nil {
		return transformMongoError(err)
	}
	return nil
}

func (r *MongoUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	objID, err := repository.StringToObjectId(id)
	if err != nil {
		return nil, userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeValidation, "invalid ID format", err)
	}

	var user users.User
	err = r.db.Collection(UserCollection).FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, userrepo.NewUserRepositoryError(
				userrepo.ErrorTypeNotFound, "user not found", nil)
		}
		return nil, userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeInternal, "database error", err)
	}
	return &user, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objID, err := repository.StringToObjectId(id)
	if err != nil {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeValidation, "invalid ID format", err)
	}

	updates["updated_at"] = time.Now()
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": updates}

	result, err := r.db.Collection(UserCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return transformMongoError(err)
	}
	if result.MatchedCount == 0 {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeNotFound, "user not found", nil)
	}
	return nil
}

func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
	objID, err := repository.StringToObjectId(id)
	if err != nil {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeValidation, "invalid ID format", err)
	}

	result, err := r.db.Collection(UserCollection).DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return transformMongoError(err)
	}
	if result.DeletedCount == 0 {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeNotFound, "user not found", nil)
	}
	return nil
}

func (r *MongoUserRepository) List(ctx context.Context, filter base.Filter) ([]*users.User, error) {
	// 将 base.Filter 转换为 bson.M 用于查询
	bsonFilter := bson.M{}
	if filter != nil {
		conditions := filter.GetConditions()
		for k, v := range conditions {
			bsonFilter[k] = v
		}
	}
	cursor, err := r.db.Collection(UserCollection).Find(ctx, bsonFilter)
	if err != nil {
		return nil, transformMongoError(err)
	}
	defer cursor.Close(ctx)

	var results []*users.User
	if err = cursor.All(ctx, &results); err != nil {
		return nil, transformMongoError(err)
	}
	return results, nil
}

func (r *MongoUserRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	// 将 base.Filter 转换为 bson.M 用于查询
	bsonFilter := bson.M{}
	if filter != nil {
		conditions := filter.GetConditions()
		for k, v := range conditions {
			bsonFilter[k] = v
		}
	}
	count, err := r.db.Collection(UserCollection).CountDocuments(ctx, bsonFilter)
	if err != nil {
		return 0, transformMongoError(err)
	}
	return count, nil
}

func (r *MongoUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	objID, err := repository.StringToObjectId(id)
	if err != nil {
		return false, userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeValidation, "invalid ID format", err)
	}

	count, err := r.db.Collection(UserCollection).CountDocuments(ctx, bson.M{"_id": objID})
	if err != nil {
		return false, transformMongoError(err)
	}
	return count > 0, nil
}

// === 实现 BaseUserRepository 共享方法 ===

func (r *MongoUserRepository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	var user users.User
	err := r.db.Collection(UserCollection).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, userrepo.NewUserRepositoryError(
				userrepo.ErrorTypeNotFound, "user not found", nil)
		}
		return nil, transformMongoError(err)
	}
	return &user, nil
}

func (r *MongoUserRepository) GetByUsername(ctx context.Context, username string) (*users.User, error) {
	var user users.User
	err := r.db.Collection(UserCollection).FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, userrepo.NewUserRepositoryError(
				userrepo.ErrorTypeNotFound, "user not found", nil)
		}
		return nil, transformMongoError(err)
	}
	return &user, nil
}

func (r *MongoUserRepository) UpdateStatus(ctx context.Context, id string, status users.UserStatus) error {
	objID, err := repository.StringToObjectId(id)
	if err != nil {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeValidation, "invalid ID format", err)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"status":     status,
		"updated_at": time.Now(),
	}}

	result, err := r.db.Collection(UserCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return transformMongoError(err)
	}
	if result.MatchedCount == 0 {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeNotFound, "user not found", nil)
	}
	return nil
}

func (r *MongoUserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	objID, err := repository.StringToObjectId(id)
	if err != nil {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeValidation, "invalid ID format", err)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"password":   hashedPassword,
		"updated_at": time.Now(),
	}}

	result, err := r.db.Collection(UserCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return transformMongoError(err)
	}
	if result.MatchedCount == 0 {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeNotFound, "user not found", nil)
	}
	return nil
}

func (r *MongoUserRepository) SetEmailVerified(ctx context.Context, id string, verified bool) error {
	objID, err := repository.StringToObjectId(id)
	if err != nil {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeValidation, "invalid ID format", err)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"email_verified": verified,
		"updated_at":     time.Now(),
	}}

	result, err := r.db.Collection(UserCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return transformMongoError(err)
	}
	if result.MatchedCount == 0 {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeNotFound, "user not found", nil)
	}
	return nil
}

func (r *MongoUserRepository) BatchUpdateStatus(ctx context.Context, ids []string, status users.UserStatus) error {
	objIDs, err := repository.StringSliceToObjectIDSlice(ids)
	if err != nil {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeValidation, "invalid ID in batch", err)
	}

	if len(objIDs) == 0 {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeValidation, "no valid IDs to update", nil)
	}

	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	update := bson.M{"$set": bson.M{
		"status":     status,
		"updated_at": time.Now(),
	}}

	_, err = r.db.Collection(UserCollection).UpdateMany(ctx, filter, update)
	return err
}

func (r *MongoUserRepository) BatchDelete(ctx context.Context, ids []string) error {
	objIDs, err := repository.StringSliceToObjectIDSlice(ids)
	if err != nil {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeValidation, "invalid ID in batch", err)
	}

	if len(objIDs) == 0 {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeValidation, "no valid IDs to delete", nil)
	}

	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	_, err = r.db.Collection(UserCollection).DeleteMany(ctx, filter)
	return err
}

func (r *MongoUserRepository) CountByStatus(ctx context.Context, status users.UserStatus) (int64, error) {
	count, err := r.db.Collection(UserCollection).CountDocuments(ctx, bson.M{"status": status})
	if err != nil {
		return 0, transformMongoError(err)
	}
	return count, nil
}

func (r *MongoUserRepository) CountByRole(ctx context.Context, role string) (int64, error) {
	count, err := r.db.Collection(UserCollection).CountDocuments(ctx, bson.M{"roles": role})
	if err != nil {
		return 0, transformMongoError(err)
	}
	return count, nil
}

// === 辅助方法 ===

func transformMongoError(err error) error {
	if err == nil {
		return nil
	}
	if err == mongo.ErrNoDocuments {
		return userrepo.NewUserRepositoryError(
			userrepo.ErrorTypeNotFound, "document not found", err)
	}
	return userrepo.NewUserRepositoryError(
		userrepo.ErrorTypeInternal, "database operation failed", err)
}
