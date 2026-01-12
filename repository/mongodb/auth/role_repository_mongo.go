package auth

import (
	authModel "Qingyu_backend/models/auth"
	authInterface "Qingyu_backend/repository/interfaces/auth"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RoleRepositoryImpl 角色Repository实现
type RoleRepositoryImpl struct {
	db             *mongo.Database
	roleCollection *mongo.Collection
}

// NewRoleRepository 创建角色Repository
func NewRoleRepository(db *mongo.Database) authInterface.RoleRepository {
	return &RoleRepositoryImpl{
		db:             db,
		roleCollection: db.Collection("roles"),
	}
}

// ============ 角色管理 ============

// CreateRole 创建角色
func (r *RoleRepositoryImpl) CreateRole(ctx context.Context, role *authModel.Role) error {
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
func (r *RoleRepositoryImpl) GetRole(ctx context.Context, roleID string) (*authModel.Role, error) {
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
func (r *RoleRepositoryImpl) GetRoleByName(ctx context.Context, name string) (*authModel.Role, error) {
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
func (r *RoleRepositoryImpl) UpdateRole(ctx context.Context, roleID string, updates map[string]interface{}) error {
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
func (r *RoleRepositoryImpl) DeleteRole(ctx context.Context, roleID string) error {
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
func (r *RoleRepositoryImpl) ListRoles(ctx context.Context) ([]*authModel.Role, error) {
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
func (r *RoleRepositoryImpl) AssignUserRole(ctx context.Context, userID, roleID string) error {
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
func (r *RoleRepositoryImpl) RemoveUserRole(ctx context.Context, userID, roleID string) error {
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
func (r *RoleRepositoryImpl) GetUserRoles(ctx context.Context, userID string) ([]*authModel.Role, error) {
	userCollection := r.db.Collection("users")

	// 查询用户 - 同时支持 role（字符串）和 roles（数组）字段
	// 注意：users集合的_id字段是string类型，不是ObjectID，所以直接使用userID
	var user struct {
		Role  string   `bson:"role"`  // 单个角色（旧格式）
		Roles []string `bson:"roles"` // 多个角色（新格式）
	}
	err := userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("用户不存在: %s", userID)
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 收集所有角色名称
	roleNames := []string{}

	// 1. 处理旧的 role 字段（字符串）
	if user.Role != "" {
		roleNames = append(roleNames, user.Role)
	}

	// 2. 处理新的 roles 字段（数组）
	for _, role := range user.Roles {
		if role != "" && role != user.Role {
			roleNames = append(roleNames, role)
		}
	}

	// 如果没有任何角色，返回空列表
	if len(roleNames) == 0 {
		return []*authModel.Role{}, nil
	}

	// 尝试从 roles 集合查询角色对象（如果存在）
	roleObjectIDs := make([]primitive.ObjectID, 0, len(roleNames))
	for _, roleID := range roleNames {
		if oid, err := primitive.ObjectIDFromHex(roleID); err == nil {
			roleObjectIDs = append(roleObjectIDs, oid)
		}
	}

	var roles []*authModel.Role

	// 如果有有效的ObjectID，尝试从roles集合查询
	if len(roleObjectIDs) > 0 {
		cursor, err := r.roleCollection.Find(ctx, bson.M{"_id": bson.M{"$in": roleObjectIDs}})
		if err == nil {
			defer cursor.Close(ctx)
			if err = cursor.All(ctx, &roles); err == nil && len(roles) > 0 {
				return roles, nil
			}
		}
	}

	// 如果roles集合不存在或没有数据，从角色名称创建虚拟角色对象（向后兼容）
	roles = make([]*authModel.Role, 0, len(roleNames))
	for _, roleName := range roleNames {
		// 跳过ObjectID格式的字符串（已经尝试查询过了）
		if _, err := primitive.ObjectIDFromHex(roleName); err == nil {
			continue
		}

		// 为角色名称创建虚拟角色对象
		role := &authModel.Role{
			ID:          roleName, // 使用名称作为ID
			Name:        roleName,
			Description: roleName + " role",
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// HasUserRole 检查用户是否有指定角色
func (r *RoleRepositoryImpl) HasUserRole(ctx context.Context, userID, roleID string) (bool, error) {
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
func (r *RoleRepositoryImpl) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	role, err := r.GetRole(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return role.Permissions, nil
}

// GetUserPermissions 获取用户权限（所有角色的权限并集）
func (r *RoleRepositoryImpl) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
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
func (r *RoleRepositoryImpl) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}
