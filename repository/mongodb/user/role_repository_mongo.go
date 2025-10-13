package user

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	usersModel "Qingyu_backend/models/users"
	infrastructure "Qingyu_backend/repository/interfaces/infrastructure"
	UserInterface "Qingyu_backend/repository/interfaces/user"
)

// MongoRoleRepository MongoDB角色仓储实现
type MongoRoleRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoRoleRepository 创建新的MongoDB角色仓储实例
func NewMongoRoleRepository(db *mongo.Database) UserInterface.RoleRepository {
	return &MongoRoleRepository{
		db:         db,
		collection: db.Collection("roles"),
	}
}

// Health 健康检查
func (r *MongoRoleRepository) Health(ctx context.Context) error {
	if err := r.db.Client().Ping(ctx, nil); err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeConnection,
			"数据库连接失败",
			err,
		)
	}
	return nil
}

// Create 创建角色
func (r *MongoRoleRepository) Create(ctx context.Context, role *usersModel.Role) error {
	if role == nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"角色对象不能为空",
			nil,
		)
	}

	// 设置创建时间
	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now

	// 生成ObjectID
	if role.ID == "" {
		role.ID = primitive.NewObjectID().Hex()
	}

	// 插入文档
	_, err := r.collection.InsertOne(ctx, role)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeDuplicate,
				"角色名已存在",
				err,
			)
		}
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"创建角色失败",
			err,
		)
	}

	return nil
}

// GetByID 根据ID获取角色
func (r *MongoRoleRepository) GetByID(ctx context.Context, id string) (*usersModel.Role, error) {
	var role usersModel.Role

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeNotFound,
				"角色不存在",
				err,
			)
		}
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"查询角色失败",
			err,
		)
	}

	return &role, nil
}

// Update 更新角色信息
func (r *MongoRoleRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"更新角色失败",
			err,
		)
	}

	if result.MatchedCount == 0 {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeNotFound,
			"角色不存在",
			nil,
		)
	}

	return nil
}

// Delete 删除角色
func (r *MongoRoleRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"删除角色失败",
			err,
		)
	}

	if result.DeletedCount == 0 {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeNotFound,
			"角色不存在",
			nil,
		)
	}

	return nil
}

// List 获取角色列表
func (r *MongoRoleRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*usersModel.Role, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"查询角色列表失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	var roles []*usersModel.Role
	for cursor.Next(ctx) {
		var role usersModel.Role
		if err := cursor.Decode(&role); err != nil {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeInternal,
				"解析角色数据失败",
				err,
			)
		}
		roles = append(roles, &role)
	}

	if err := cursor.Err(); err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"遍历角色数据失败",
			err,
		)
	}

	return roles, nil
}

// Count 统计角色数量
func (r *MongoRoleRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"统计角色数量失败",
			err,
		)
	}
	return count, nil
}

// Exists 检查角色是否存在
func (r *MongoRoleRepository) Exists(ctx context.Context, id string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"检查角色存在性失败",
			err,
		)
	}
	return count > 0, nil
}

// GetByName 根据名称获取角色
func (r *MongoRoleRepository) GetByName(ctx context.Context, name string) (*usersModel.Role, error) {
	var role usersModel.Role

	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeNotFound,
				fmt.Sprintf("角色 %s 不存在", name),
				err,
			)
		}
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"查询角色失败",
			err,
		)
	}

	return &role, nil
}

// GetDefaultRole 获取默认角色
func (r *MongoRoleRepository) GetDefaultRole(ctx context.Context) (*usersModel.Role, error) {
	var role usersModel.Role

	err := r.collection.FindOne(ctx, bson.M{"is_default": true}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeNotFound,
				"默认角色不存在",
				err,
			)
		}
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"查询默认角色失败",
			err,
		)
	}

	return &role, nil
}

// ExistsByName 检查角色名是否存在
func (r *MongoRoleRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		return false, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"检查角色名存在性失败",
			err,
		)
	}
	return count > 0, nil
}

// ListAllRoles 列出所有角色
func (r *MongoRoleRepository) ListAllRoles(ctx context.Context) ([]*usersModel.Role, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"查询所有角色失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	var roles []*usersModel.Role
	for cursor.Next(ctx) {
		var role usersModel.Role
		if err := cursor.Decode(&role); err != nil {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeInternal,
				"解析角色数据失败",
				err,
			)
		}
		roles = append(roles, &role)
	}

	if err := cursor.Err(); err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"遍历角色数据失败",
			err,
		)
	}

	return roles, nil
}

// ListDefaultRoles 列出默认角色
func (r *MongoRoleRepository) ListDefaultRoles(ctx context.Context) ([]*usersModel.Role, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{"is_default": true}, opts)
	if err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"查询默认角色失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	var roles []*usersModel.Role
	for cursor.Next(ctx) {
		var role usersModel.Role
		if err := cursor.Decode(&role); err != nil {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeInternal,
				"解析角色数据失败",
				err,
			)
		}
		roles = append(roles, &role)
	}

	if err := cursor.Err(); err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"遍历角色数据失败",
			err,
		)
	}

	return roles, nil
}

// GetRolePermissions 获取角色的权限列表
func (r *MongoRoleRepository) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	role, err := r.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return role.Permissions, nil
}

// UpdateRolePermissions 更新角色的全部权限
func (r *MongoRoleRepository) UpdateRolePermissions(ctx context.Context, roleID string, permissions []string) error {
	now := time.Now()
	filter := bson.M{"_id": roleID}
	update := bson.M{"$set": bson.M{
		"permissions": permissions,
		"updated_at":  now,
	}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"更新角色权限失败",
			err,
		)
	}

	if result.MatchedCount == 0 {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeNotFound,
			fmt.Sprintf("角色ID %s 不存在", roleID),
			nil,
		)
	}

	return nil
}

// AddPermission 为角色添加权限
func (r *MongoRoleRepository) AddPermission(ctx context.Context, roleID string, permission string) error {
	now := time.Now()
	filter := bson.M{"_id": roleID}
	update := bson.M{
		"$addToSet": bson.M{"permissions": permission},
		"$set":      bson.M{"updated_at": now},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"添加角色权限失败",
			err,
		)
	}

	if result.MatchedCount == 0 {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeNotFound,
			fmt.Sprintf("角色ID %s 不存在", roleID),
			nil,
		)
	}

	return nil
}

// RemovePermission 从角色移除权限
func (r *MongoRoleRepository) RemovePermission(ctx context.Context, roleID string, permission string) error {
	now := time.Now()
	filter := bson.M{"_id": roleID}
	update := bson.M{
		"$pull": bson.M{"permissions": permission},
		"$set":  bson.M{"updated_at": now},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"移除角色权限失败",
			err,
		)
	}

	if result.MatchedCount == 0 {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeNotFound,
			fmt.Sprintf("角色ID %s 不存在", roleID),
			nil,
		)
	}

	return nil
}

// CountByName 按名称统计角色数量
func (r *MongoRoleRepository) CountByName(ctx context.Context, name string) (int64, error) {
	filter := bson.M{"name": name}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"统计角色数量失败",
			err,
		)
	}

	return count, nil
}
