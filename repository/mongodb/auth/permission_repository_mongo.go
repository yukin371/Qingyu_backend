package auth

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/auth"
	authrepo "Qingyu_backend/repository/interfaces/auth"
)

const (
	PermissionCollection = "permissions"
	RoleCollection       = "roles"
	UserRoleCollection   = "user_roles"
)

// MongoPermissionRepository MongoDB权限仓储实现
type MongoPermissionRepository struct {
	db *mongo.Database
}

// NewMongoPermissionRepository 创建MongoDB权限仓储
func NewMongoPermissionRepository(db *mongo.Database) authrepo.PermissionRepository {
	return &MongoPermissionRepository{db: db}
}

// ==================== 权限管理 ====================

func (r *MongoPermissionRepository) GetAllPermissions(ctx context.Context) ([]*auth.Permission, error) {
	cursor, err := r.db.Collection(PermissionCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*auth.Permission
	if err = cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *MongoPermissionRepository) GetPermissionByCode(ctx context.Context, code string) (*auth.Permission, error) {
	var permission auth.Permission
	filter := bson.M{"code": code}

	err := r.db.Collection(PermissionCollection).FindOne(ctx, filter).Decode(&permission)
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

func (r *MongoPermissionRepository) CreatePermission(ctx context.Context, permission *auth.Permission) error {
	permission.CreatedAt = time.Now()

	_, err := r.db.Collection(PermissionCollection).InsertOne(ctx, permission)
	return err
}

func (r *MongoPermissionRepository) UpdatePermission(ctx context.Context, permission *auth.Permission) error {
	filter := bson.M{"code": permission.Code}
	update := bson.M{"$set": permission}

	_, err := r.db.Collection(PermissionCollection).UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoPermissionRepository) DeletePermission(ctx context.Context, code string) error {
	filter := bson.M{"code": code}
	_, err := r.db.Collection(PermissionCollection).DeleteOne(ctx, filter)
	return err
}

// ==================== 角色管理 ====================

func (r *MongoPermissionRepository) GetAllRoles(ctx context.Context) ([]*auth.Role, error) {
	cursor, err := r.db.Collection(RoleCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []*auth.Role
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *MongoPermissionRepository) GetRoleByID(ctx context.Context, roleID string) (*auth.Role, error) {
	var role auth.Role
	filter := bson.M{"_id": roleID}

	err := r.db.Collection(RoleCollection).FindOne(ctx, filter).Decode(&role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *MongoPermissionRepository) GetRoleByName(ctx context.Context, name string) (*auth.Role, error) {
	var role auth.Role
	filter := bson.M{"name": name}

	err := r.db.Collection(RoleCollection).FindOne(ctx, filter).Decode(&role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *MongoPermissionRepository) CreateRole(ctx context.Context, role *auth.Role) error {
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	_, err := r.db.Collection(RoleCollection).InsertOne(ctx, role)
	return err
}

func (r *MongoPermissionRepository) UpdateRole(ctx context.Context, role *auth.Role) error {
	role.UpdatedAt = time.Now()

	filter := bson.M{"_id": role.ID}
	update := bson.M{"$set": role}

	_, err := r.db.Collection(RoleCollection).UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoPermissionRepository) DeleteRole(ctx context.Context, roleID string) error {
	filter := bson.M{"_id": roleID, "is_system": false} // 不能删除系统角色
	result, err := r.db.Collection(RoleCollection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("role not found or is system role")
	}
	return nil
}

func (r *MongoPermissionRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionCode string) error {
	filter := bson.M{"_id": roleID}
	update := bson.M{
		"$addToSet": bson.M{"permissions": permissionCode},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	_, err := r.db.Collection(RoleCollection).UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoPermissionRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionCode string) error {
	filter := bson.M{"_id": roleID}
	update := bson.M{
		"$pull": bson.M{"permissions": permissionCode},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err := r.db.Collection(RoleCollection).UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoPermissionRepository) GetRolePermissions(ctx context.Context, roleID string) ([]*auth.Permission, error) {
	// 先获取角色
	role, err := r.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	if len(role.Permissions) == 0 {
		return []*auth.Permission{}, nil
	}

	// 获取权限详情
	filter := bson.M{"code": bson.M{"$in": role.Permissions}}
	cursor, err := r.db.Collection(PermissionCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*auth.Permission
	if err = cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

// ==================== 用户角色管理 ====================

func (r *MongoPermissionRepository) GetUserRoles(ctx context.Context, userID primitive.ObjectID) ([]string, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.db.Collection(UserRoleCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userRoles []*auth.UserRole
	if err = cursor.All(ctx, &userRoles); err != nil {
		return nil, err
	}

	roles := make([]string, 0, len(userRoles))
	for _, ur := range userRoles {
		// UserRole.RoleID 存储的是角色名称
		roles = append(roles, ur.RoleID)
	}

	return roles, nil
}

func (r *MongoPermissionRepository) AssignRoleToUser(ctx context.Context, userID primitive.ObjectID, roleName string) error {
	userRole := &auth.UserRole{
		UserID:     userID.Hex(),
		RoleID:     roleName, // 存储角色名称
		AssignedAt: time.Now(),
	}

	filter := bson.M{"user_id": userID, "role_id": roleName}
	update := bson.M{"$setOnInsert": userRole}

	_, err := r.db.Collection(UserRoleCollection).UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *MongoPermissionRepository) RemoveRoleFromUser(ctx context.Context, userID primitive.ObjectID, roleName string) error {
	filter := bson.M{"user_id": userID, "role_id": roleName}
	_, err := r.db.Collection(UserRoleCollection).DeleteOne(ctx, filter)
	return err
}

func (r *MongoPermissionRepository) ClearUserRoles(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID}
	_, err := r.db.Collection(UserRoleCollection).DeleteMany(ctx, filter)
	return err
}

// ==================== 权限检查 ====================

func (r *MongoPermissionRepository) UserHasPermission(ctx context.Context, userID primitive.ObjectID, permissionCode string) (bool, error) {
	// 获取用户角色
	roles, err := r.GetUserRoles(ctx, userID)
	if err != nil {
		return false, err
	}

	if len(roles) == 0 {
		return false, nil
	}

	// 获取角色详情
	filter := bson.M{"name": bson.M{"$in": roles}}
	cursor, err := r.db.Collection(RoleCollection).Find(ctx, filter)
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	var roleList []*auth.Role
	if err = cursor.All(ctx, &roleList); err != nil {
		return false, err
	}

	// 检查权限
	for _, role := range roleList {
		for _, permCode := range role.Permissions {
			if permCode == permissionCode {
				return true, nil
			}
		}
	}

	return false, nil
}

func (r *MongoPermissionRepository) UserHasAnyPermission(ctx context.Context, userID primitive.ObjectID, permissionCodes []string) (bool, error) {
	for _, perm := range permissionCodes {
		has, err := r.UserHasPermission(ctx, userID, perm)
		if err != nil {
			return false, err
		}
		if has {
			return true, nil
		}
	}
	return false, nil
}

func (r *MongoPermissionRepository) UserHasAllPermissions(ctx context.Context, userID primitive.ObjectID, permissionCodes []string) (bool, error) {
	for _, perm := range permissionCodes {
		has, err := r.UserHasPermission(ctx, userID, perm)
		if err != nil {
			return false, err
		}
		if !has {
			return false, nil
		}
	}
	return true, nil
}

func (r *MongoPermissionRepository) GetUserPermissions(ctx context.Context, userID primitive.ObjectID) ([]*auth.Permission, error) {
	// 获取用户角色
	roles, err := r.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return []*auth.Permission{}, nil
	}

	// 获取角色详情
	filter := bson.M{"name": bson.M{"$in": roles}}
	cursor, err := r.db.Collection(RoleCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roleList []*auth.Role
	if err = cursor.All(ctx, &roleList); err != nil {
		return nil, err
	}

	// 收集所有权限代码
	permissionCodes := make([]string, 0)
	for _, role := range roleList {
		permissionCodes = append(permissionCodes, role.Permissions...)
	}

	if len(permissionCodes) == 0 {
		return []*auth.Permission{}, nil
	}

	// 去重
	uniqueCodes := make(map[string]bool)
	for _, code := range permissionCodes {
		uniqueCodes[code] = true
	}
	permissionCodes = make([]string, 0, len(uniqueCodes))
	for code := range uniqueCodes {
		permissionCodes = append(permissionCodes, code)
	}

	// 获取权限详情
	filter = bson.M{"code": bson.M{"$in": permissionCodes}}
	cursor, err = r.db.Collection(PermissionCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*auth.Permission
	if err = cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}
