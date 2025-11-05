package shared

import (
	"context"
	"fmt"
	"time"

	authModel "Qingyu_backend/models/shared/auth"
	sharedInterfaces "Qingyu_backend/repository/interfaces/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuthRepositoryImpl 认证Repository实现
type AuthRepositoryImpl struct {
	db             *mongo.Database
	roleCollection *mongo.Collection
}

// NewAuthRepository 创建认证Repository
func NewAuthRepository(db *mongo.Database) sharedInterfaces.AuthRepository {
	return &AuthRepositoryImpl{
		db:             db,
		roleCollection: db.Collection("roles"),
	}
}

// ============ 角色管理 ============

// CreateRole 创建角色
func (r *AuthRepositoryImpl) CreateRole(ctx context.Context, role *authModel.Role) error {
	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now

	result, err := r.roleCollection.InsertOne(ctx, role)
	if err != nil {
		return fmt.Errorf("插入角色失败: %w", err)
	}

	// 设置ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		role.ID = oid.Hex()
	}

	return nil
}

// GetRole 获取角色
func (r *AuthRepositoryImpl) GetRole(ctx context.Context, roleID string) (*authModel.Role, error) {
	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return nil, fmt.Errorf("无效的角色ID: %w", err)
	}

	var role authModel.Role
	err = r.roleCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("角色不存在: %s", roleID)
		}
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	return &role, nil
}

// GetRoleByName 根据名称获取角色
func (r *AuthRepositoryImpl) GetRoleByName(ctx context.Context, name string) (*authModel.Role, error) {
	var role authModel.Role
	err := r.roleCollection.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("角色不存在: %s", name)
		}
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	return &role, nil
}

// UpdateRole 更新角色
func (r *AuthRepositoryImpl) UpdateRole(ctx context.Context, roleID string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return fmt.Errorf("无效的角色ID: %w", err)
	}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	result, err := r.roleCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return fmt.Errorf("更新角色失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("角色不存在: %s", roleID)
	}

	return nil
}

// DeleteRole 删除角色
func (r *AuthRepositoryImpl) DeleteRole(ctx context.Context, roleID string) error {
	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return fmt.Errorf("无效的角色ID: %w", err)
	}

	// 检查是否是系统角色
	var role authModel.Role
	err = r.roleCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("角色不存在: %s", roleID)
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}

	if role.IsSystem {
		return fmt.Errorf("不能删除系统角色: %s", role.Name)
	}

	result, err := r.roleCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("删除角色失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("角色不存在: %s", roleID)
	}

	return nil
}

// ListRoles 列出所有角色
func (r *AuthRepositoryImpl) ListRoles(ctx context.Context) ([]*authModel.Role, error) {
	cursor, err := r.roleCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{bson.E{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("查询角色列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var roles []*authModel.Role
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, fmt.Errorf("解析角色列表失败: %w", err)
	}

	return roles, nil
}

// ============ 用户角色关联 ============

// AssignUserRole 分配用户角色
func (r *AuthRepositoryImpl) AssignUserRole(ctx context.Context, userID, roleID string) error {
	// 验证角色是否存在
	_, err := r.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	// 获取users集合
	userCollection := r.db.Collection("users")

	// 更新用户的roles数组
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("无效的用户ID: %w", err)
	}

	result, err := userCollection.UpdateOne(
		ctx,
		bson.M{"_id": userObjectID},
		bson.M{"$addToSet": bson.M{"roles": roleID}},
	)
	if err != nil {
		return fmt.Errorf("分配角色失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("用户不存在: %s", userID)
	}

	return nil
}

// RemoveUserRole 移除用户角色
func (r *AuthRepositoryImpl) RemoveUserRole(ctx context.Context, userID, roleID string) error {
	userCollection := r.db.Collection("users")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("无效的用户ID: %w", err)
	}

	result, err := userCollection.UpdateOne(
		ctx,
		bson.M{"_id": userObjectID},
		bson.M{"$pull": bson.M{"roles": roleID}},
	)
	if err != nil {
		return fmt.Errorf("移除角色失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("用户不存在: %s", userID)
	}

	return nil
}

// GetUserRoles 获取用户角色
func (r *AuthRepositoryImpl) GetUserRoles(ctx context.Context, userID string) ([]*authModel.Role, error) {
	userCollection := r.db.Collection("users")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID: %w", err)
	}

	// 查询用户的roles数组
	var user struct {
		Roles []string `bson:"roles"`
	}
	err = userCollection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("用户不存在: %s", userID)
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 如果没有角色，返回空列表
	if len(user.Roles) == 0 {
		return []*authModel.Role{}, nil
	}

	// 查询角色详情
	roleObjectIDs := make([]primitive.ObjectID, 0, len(user.Roles))
	for _, roleID := range user.Roles {
		if oid, err := primitive.ObjectIDFromHex(roleID); err == nil {
			roleObjectIDs = append(roleObjectIDs, oid)
		}
	}

	cursor, err := r.roleCollection.Find(ctx, bson.M{"_id": bson.M{"$in": roleObjectIDs}})
	if err != nil {
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}
	defer cursor.Close(ctx)

	var roles []*authModel.Role
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, fmt.Errorf("解析角色失败: %w", err)
	}

	return roles, nil
}

// HasUserRole 检查用户是否有指定角色
func (r *AuthRepositoryImpl) HasUserRole(ctx context.Context, userID, roleID string) (bool, error) {
	userCollection := r.db.Collection("users")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, fmt.Errorf("无效的用户ID: %w", err)
	}

	count, err := userCollection.CountDocuments(
		ctx,
		bson.M{
			"_id":   userObjectID,
			"roles": roleID,
		},
	)
	if err != nil {
		return false, fmt.Errorf("查询用户角色失败: %w", err)
	}

	return count > 0, nil
}

// ============ 权限查询 ============

// GetRolePermissions 获取角色权限
func (r *AuthRepositoryImpl) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	role, err := r.GetRole(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return role.Permissions, nil
}

// GetUserPermissions 获取用户权限（所有角色的权限并集）
func (r *AuthRepositoryImpl) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	roles, err := r.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 使用map去重
	permissionMap := make(map[string]bool)
	for _, role := range roles {
		for _, perm := range role.Permissions {
			permissionMap[perm] = true
		}
	}

	// 转换为数组
	permissions := make([]string, 0, len(permissionMap))
	for perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// ============ 健康检查 ============

// Health 健康检查
func (r *AuthRepositoryImpl) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}
