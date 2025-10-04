package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// PermissionServiceImpl 权限服务实现
type PermissionServiceImpl struct {
	authRepo    sharedRepo.AuthRepository
	cacheClient CacheClient // 缓存客户端（Redis）
	cacheTTL    time.Duration
}

// CacheClient 缓存客户端接口
type CacheClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// NewPermissionService 创建权限服务
func NewPermissionService(authRepo sharedRepo.AuthRepository, cacheClient CacheClient) PermissionService {
	return &PermissionServiceImpl{
		authRepo:    authRepo,
		cacheClient: cacheClient,
		cacheTTL:    5 * time.Minute, // 权限缓存5分钟
	}
}

// ============ 权限检查 ============

// CheckPermission 检查用户是否有指定权限
func (s *PermissionServiceImpl) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	// 1. 获取用户权限
	permissions, err := s.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("获取用户权限失败: %w", err)
	}

	// 2. 检查是否有通配符权限 "*"
	for _, perm := range permissions {
		if perm == "*" {
			return true, nil
		}
	}

	// 3. 检查是否有精确匹配
	for _, perm := range permissions {
		if perm == permission {
			return true, nil
		}
	}

	// 4. 检查通配符匹配（例如：book.* 匹配 book.read）
	for _, perm := range permissions {
		if strings.HasSuffix(perm, ".*") {
			prefix := strings.TrimSuffix(perm, ".*")
			if strings.HasPrefix(permission, prefix+".") {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetUserPermissions 获取用户权限
func (s *PermissionServiceImpl) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	// 1. 尝试从缓存获取
	if s.cacheClient != nil {
		cacheKey := s.getPermissionCacheKey(userID)
		cached, err := s.cacheClient.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			// 解析缓存的权限列表（用逗号分隔）
			if cached == "[]" {
				return []string{}, nil
			}
			return strings.Split(cached, ","), nil
		}
	}

	// 2. 从数据库获取
	permissions, err := s.authRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户权限失败: %w", err)
	}

	// 3. 存入缓存
	if s.cacheClient != nil {
		cacheKey := s.getPermissionCacheKey(userID)
		cacheValue := "[]"
		if len(permissions) > 0 {
			cacheValue = strings.Join(permissions, ",")
		}
		_ = s.cacheClient.Set(ctx, cacheKey, cacheValue, s.cacheTTL)
	}

	return permissions, nil
}

// GetRolePermissions 获取角色权限
func (s *PermissionServiceImpl) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	permissions, err := s.authRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("获取角色权限失败: %w", err)
	}

	return permissions, nil
}

// HasRole 检查用户是否有指定角色
func (s *PermissionServiceImpl) HasRole(ctx context.Context, userID, role string) (bool, error) {
	// 1. 获取用户角色
	roles, err := s.authRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("获取用户角色失败: %w", err)
	}

	// 2. 检查角色
	for _, r := range roles {
		if r.Name == role {
			return true, nil
		}
	}

	return false, nil
}

// ============ 缓存管理 ============

// InvalidateUserPermissionsCache 清除用户权限缓存
func (s *PermissionServiceImpl) InvalidateUserPermissionsCache(ctx context.Context, userID string) error {
	if s.cacheClient == nil {
		return nil
	}

	cacheKey := s.getPermissionCacheKey(userID)
	return s.cacheClient.Delete(ctx, cacheKey)
}

// getPermissionCacheKey 获取权限缓存Key
func (s *PermissionServiceImpl) getPermissionCacheKey(userID string) string {
	return fmt.Sprintf("user:permissions:%s", userID)
}
