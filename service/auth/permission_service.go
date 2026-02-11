package auth

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	authModel "Qingyu_backend/models/auth"
	middlewareAuth "Qingyu_backend/internal/middleware/auth"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// PermissionServiceImpl 权限服务实现
type PermissionServiceImpl struct {
	authRepo    sharedRepo.AuthRepository
	cacheClient CacheClient // 缓存客户端（Redis）
	cacheTTL    time.Duration
	logger      *zap.Logger

	// RBAC集成
	checker *middlewareAuth.RBACChecker

	// 内存缓存
	roleCache       map[string]*authModel.Role
	cacheMutex      sync.RWMutex
	cacheLastUpdate time.Time

	// singleflight防止缓存击穿
	cacheGroup singleflight.Group
}

// CacheClient 缓存客户端接口
type CacheClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// NewPermissionService 创建权限服务
func NewPermissionService(authRepo sharedRepo.AuthRepository, cacheClient CacheClient, logger *zap.Logger) PermissionService {
	return &PermissionServiceImpl{
		authRepo:    authRepo,
		cacheClient: cacheClient,
		cacheTTL:    5 * time.Minute, // 权限缓存5分钟
		logger:      logger,
		roleCache:   make(map[string]*authModel.Role),
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

	// 2. 检查是否有完全通配符权限 "*:*" 或 "*"
	for _, perm := range permissions {
		if perm == "*:*" || perm == "*" {
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
		// 支持 ".*" 或 ":*" 格式
		if strings.HasSuffix(perm, ".*") {
			prefix := strings.TrimSuffix(perm, ".*")
			if strings.HasPrefix(permission, prefix+".") {
				return true, nil
			}
		} else if strings.HasSuffix(perm, ":*") {
			prefix := strings.TrimSuffix(perm, ":*")
			if strings.HasPrefix(permission, prefix+":") {
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

// ============ RBAC集成 ============

// SetChecker 设置RBAC检查器，用于动态权限更新
func (s *PermissionServiceImpl) SetChecker(checker interface{}) {
	if rbacChecker, ok := checker.(*middlewareAuth.RBACChecker); ok {
		s.checker = rbacChecker
		if s.logger != nil {
			s.logger.Info("RBACChecker已设置到PermissionService")
		}
	}
}

// LoadPermissionsToChecker 从数据库加载权限到RBACChecker
func (s *PermissionServiceImpl) LoadPermissionsToChecker(ctx context.Context) error {
	if s.checker == nil {
		return fmt.Errorf("RBACChecker未设置，请先调用SetChecker")
	}

	if s.logger != nil {
		s.logger.Info("从数据库加载权限到RBACChecker")
	}

	// 1. 加载所有角色
	roles, err := s.authRepo.ListRoles(ctx)
	if err != nil {
		return fmt.Errorf("加载角色列表失败: %w", err)
	}

	// 2. 为每个角色设置权限
	for _, role := range roles {
		if s.logger != nil {
			s.logger.Debug("加载角色权限",
				zap.String("role", role.Name),
				zap.Int("permissions", len(role.Permissions)),
			)
		}

		// 批量授予权限
		if len(role.Permissions) > 0 {
			// 转换权限格式（从 "user.read" 到 "user:read"）
			convertedPerms := s.convertPermissions(role.Permissions)
			s.checker.BatchGrantPermissions(role.Name, convertedPerms)
		}
	}

	if s.logger != nil {
		s.logger.Info("权限加载完成",
			zap.Int("roles", len(roles)),
		)
	}

	return nil
}

// LoadUserRolesToChecker 从数据库加载用户角色到RBACChecker
func (s *PermissionServiceImpl) LoadUserRolesToChecker(ctx context.Context, userID string) error {
	if s.checker == nil {
		return fmt.Errorf("RBACChecker未设置，请先调用SetChecker")
	}

	// 获取用户角色列表
	roles, err := s.authRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户角色失败: %w", err)
	}

	// 分配角色到checker
	for _, role := range roles {
		s.checker.AssignRole(userID, role.Name)
	}

	if s.logger != nil {
		s.logger.Debug("用户角色已加载到checker",
			zap.String("user", userID),
			zap.Int("roles", len(roles)),
		)
	}

	return nil
}

// ReloadAllFromDatabase 从数据库重新加载所有权限和角色
func (s *PermissionServiceImpl) ReloadAllFromDatabase(ctx context.Context) error {
	if s.logger != nil {
		s.logger.Info("从数据库重新加载所有权限和角色")
	}

	// 加载到checker
	if s.checker != nil {
		if err := s.LoadPermissionsToChecker(ctx); err != nil {
			return fmt.Errorf("加载权限到checker失败: %w", err)
		}
	}

	if s.logger != nil {
		s.logger.Info("权限重新加载完成")
	}

	return nil
}

// convertPermissions 转换权限格式（从 "user.read" 到 "user:read"）
func (s *PermissionServiceImpl) convertPermissions(permissions []string) []string {
	converted := make([]string, len(permissions))
	for i, perm := range permissions {
		// 如果使用 "." 分隔符，转换为 ":"
		converted[i] = strings.ReplaceAll(perm, ".", ":")
	}
	return converted
}

// ============ 内存缓存管理 ============

// initializeMemoryCache 初始化内存缓存
func (s *PermissionServiceImpl) initializeMemoryCache(ctx context.Context) error {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// 使用singleflight防止并发加载
	_, err, _ := s.cacheGroup.Do("init_roles", func() (interface{}, error) {
		// 加载所有角色
		roles, err := s.authRepo.ListRoles(ctx)
		if err != nil {
			return nil, fmt.Errorf("加载角色列表失败: %w", err)
		}

		// 构建缓存
		s.roleCache = make(map[string]*authModel.Role)
		for _, role := range roles {
			s.roleCache[role.Name] = role
		}

		s.cacheLastUpdate = time.Now()

		if s.logger != nil {
			s.logger.Info("角色缓存已初始化",
				zap.Int("roles", len(s.roleCache)),
			)
		}

		return nil, nil
	})

	return err
}

// getRoleFromCache 从缓存获取角色
func (s *PermissionServiceImpl) getRoleFromCache(roleName string) (*authModel.Role, bool) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	role, ok := s.roleCache[roleName]
	return role, ok
}

// RefreshRoleCache 刷新角色缓存
func (s *PermissionServiceImpl) RefreshRoleCache(ctx context.Context) error {
	return s.initializeMemoryCache(ctx)
}

