package admin

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/users"
	adminrepo "Qingyu_backend/repository/interfaces/admin"
)

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = fmt.Errorf("user not found")
	// ErrInvalidUserID 无效的用户ID
	ErrInvalidUserID = fmt.Errorf("invalid user ID")
	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	// ErrCannotModifySuperAdmin 不能修改超级管理员
	ErrCannotModifySuperAdmin = fmt.Errorf("cannot modify super admin")
	// ErrInvalidRole 无效的角色
	ErrInvalidRole = fmt.Errorf("invalid role")
	// ErrInvalidBatchCount 无效的批量创建数量
	ErrInvalidBatchCount = fmt.Errorf("invalid batch count")
)

const (
	maxBatchCreateUsers = 100
)

// UserAdminService 用户管理服务接口
type UserAdminService interface {
	// GetUserList 获取用户列表
	GetUserList(ctx context.Context, filter *adminrepo.UserFilter, page, pageSize int) ([]*users.User, int64, error)

	// CreateUser 创建用户
	CreateUser(ctx context.Context, req *CreateUserRequest) (*users.User, error)

	// BatchCreateUsers 批量创建用户
	BatchCreateUsers(ctx context.Context, req *BatchCreateUserRequest) ([]*users.User, error)

	// GetUserDetail 获取用户详情
	GetUserDetail(ctx context.Context, userID string) (*users.User, error)

	// UpdateUserStatus 更新用户状态
	UpdateUserStatus(ctx context.Context, userID string, status users.UserStatus) error

	// UpdateUserRole 更新用户角色
	UpdateUserRole(ctx context.Context, userID, role string) error

	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, userID string) error

	// BatchUpdateStatus 批量更新用户状态
	BatchUpdateStatus(ctx context.Context, userIDs []string, status users.UserStatus) error

	// BatchDeleteUsers 批量删除用户
	BatchDeleteUsers(ctx context.Context, userIDs []string) error

	// GetUserActivities 获取用户活动记录
	GetUserActivities(ctx context.Context, userID string, page, pageSize int) ([]*users.UserActivity, int64, error)

	// GetUserStatistics 获取用户统计信息
	GetUserStatistics(ctx context.Context, userID string) (*users.UserStatistics, error)

	// ResetUserPassword 重置用户密码
	ResetUserPassword(ctx context.Context, userID string) (string, error)

	// SearchUsers 搜索用户
	SearchUsers(ctx context.Context, keyword string, page, pageSize int) ([]*users.User, int64, error)

	// GetUsersByRole 根据角色获取用户
	GetUsersByRole(ctx context.Context, role string, page, pageSize int) ([]*users.User, int64, error)

	// CountByStatus 按状态统计用户数量
	CountByStatus(ctx context.Context) (map[string]int64, error)

	// GetRecentUsers 获取最近注册的用户
	GetRecentUsers(ctx context.Context, limit int) ([]*users.User, error)

	// GetActiveUsers 获取活跃用户
	GetActiveUsers(ctx context.Context, days int, limit int) ([]*users.User, error)
}

// UserAdminServiceImpl 用户管理服务实现
type UserAdminServiceImpl struct {
	userRepo adminrepo.UserAdminRepository
}

// NewUserAdminService 创建用户管理服务
func NewUserAdminService(userRepo adminrepo.UserAdminRepository) UserAdminService {
	return &UserAdminServiceImpl{
		userRepo: userRepo,
	}
}

// GetUserList 获取用户列表
func (s *UserAdminServiceImpl) GetUserList(ctx context.Context, filter *adminrepo.UserFilter, page, pageSize int) ([]*users.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.userRepo.List(ctx, filter, page, pageSize)
}

// GetUserDetail 获取用户详情
func (s *UserAdminServiceImpl) GetUserDetail(ctx context.Context, userID string) (*users.User, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// UpdateUserStatus 更新用户状态
func (s *UserAdminServiceImpl) UpdateUserStatus(ctx context.Context, userID string, status users.UserStatus) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidUserID
	}

	// 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	// 不能修改管理员的状态
	if user.HasRole("admin") {
		return ErrCannotModifySuperAdmin
	}

	return s.userRepo.UpdateStatus(ctx, id, status)
}

// UpdateUserRole 更新用户角色
func (s *UserAdminServiceImpl) UpdateUserRole(ctx context.Context, userID, role string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidUserID
	}

	// 验证角色
	if !isValidRole(role) {
		return ErrInvalidRole
	}

	// 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	// 不能修改管理员的角色
	if user.HasRole("admin") {
		return ErrCannotModifySuperAdmin
	}

	return s.userRepo.UpdateRoles(ctx, id, role)
}

// DeleteUser 删除用户
func (s *UserAdminServiceImpl) DeleteUser(ctx context.Context, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidUserID
	}

	// 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	// 不能删除管理员
	if user.HasRole("admin") {
		return ErrCannotModifySuperAdmin
	}

	return s.userRepo.Delete(ctx, id)
}

// BatchUpdateStatus 批量更新用户状态
func (s *UserAdminServiceImpl) BatchUpdateStatus(ctx context.Context, userIDs []string, status users.UserStatus) error {
	ids := make([]primitive.ObjectID, 0, len(userIDs))
	for _, idStr := range userIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			continue // 跳过无效ID
		}

		// 检查是否是管理员
		user, err := s.userRepo.GetByID(ctx, id)
		if err == nil && user.HasRole("admin") {
			continue // 跳过管理员
		}

		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return fmt.Errorf("no valid user IDs")
	}

	return s.userRepo.BatchUpdateStatus(ctx, ids, status)
}

// BatchDeleteUsers 批量删除用户
func (s *UserAdminServiceImpl) BatchDeleteUsers(ctx context.Context, userIDs []string) error {
	return s.BatchUpdateStatus(ctx, userIDs, users.UserStatusDeleted)
}

// GetUserActivities 获取用户活动记录
func (s *UserAdminServiceImpl) GetUserActivities(ctx context.Context, userID string, page, pageSize int) ([]*users.UserActivity, int64, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, ErrInvalidUserID
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.userRepo.GetActivities(ctx, id, page, pageSize)
}

// GetUserStatistics 获取用户统计信息
func (s *UserAdminServiceImpl) GetUserStatistics(ctx context.Context, userID string) (*users.UserStatistics, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	return s.userRepo.GetStatistics(ctx, id)
}

// ResetUserPassword 重置用户密码
func (s *UserAdminServiceImpl) ResetUserPassword(ctx context.Context, userID string) (string, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", ErrInvalidUserID
	}

	// 检查用户是否存在
	_, err = s.userRepo.GetByID(ctx, id)
	if err != nil {
		return "", ErrUserNotFound
	}

	// 生成随机密码
	newPassword, err := generateRandomPassword(12)
	if err != nil {
		return "", fmt.Errorf("failed to generate password: %w", err)
	}

	if err := s.userRepo.ResetPassword(ctx, id, newPassword); err != nil {
		return "", err
	}

	return newPassword, nil
}

// SearchUsers 搜索用户
func (s *UserAdminServiceImpl) SearchUsers(ctx context.Context, keyword string, page, pageSize int) ([]*users.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.userRepo.SearchUsers(ctx, keyword, page, pageSize)
}

// GetUsersByRole 根据角色获取用户
func (s *UserAdminServiceImpl) GetUsersByRole(ctx context.Context, role string, page, pageSize int) ([]*users.User, int64, error) {
	if !isValidRole(role) {
		return nil, 0, ErrInvalidRole
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.userRepo.GetUsersByRole(ctx, role, page, pageSize)
}

// CountByStatus 按状态统计用户数量
func (s *UserAdminServiceImpl) CountByStatus(ctx context.Context) (map[string]int64, error) {
	return s.userRepo.CountByStatus(ctx)
}

// GetRecentUsers 获取最近注册的用户
func (s *UserAdminServiceImpl) GetRecentUsers(ctx context.Context, limit int) ([]*users.User, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.userRepo.GetRecentUsers(ctx, limit)
}

// GetActiveUsers 获取活跃用户
func (s *UserAdminServiceImpl) GetActiveUsers(ctx context.Context, days int, limit int) ([]*users.User, error) {
	if days < 1 {
		days = 7 // 默认7天
	}
	if days > 365 {
		days = 365
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.userRepo.GetActiveUsers(ctx, days, limit)
}

// isValidRole 验证角色是否有效
func isValidRole(role string) bool {
	validRoles := map[string]bool{
		"user":        true,
		"author":      true,
		"admin":       true,
		"super_admin": true,
		"vip":         true,
	}
	return validRoles[role]
}

// generateRandomPassword 生成随机密码
func generateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b), nil
}

// ============ 创建用户相关 ============

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string           `json:"username" binding:"required,min=3,max=50"`
	Email    string           `json:"email" binding:"required,email"`
	Password string           `json:"password" binding:"min=6,max=100"`
	Nickname string           `json:"nickname"`
	Role     string           `json:"role" binding:"required,oneof=reader author admin"`
	Status   users.UserStatus `json:"status"`
	Bio      string           `json:"bio"`
}

// BatchCreateUserRequest 批量创建用户请求
type BatchCreateUserRequest struct {
	Count  int              `json:"count" binding:"required,min=1,max=100"`
	Prefix string           `json:"prefix"`
	Role   string           `json:"role" binding:"required,oneof=reader author admin"`
	Status users.UserStatus `json:"status"`
}

// CreateUser 创建用户
func (s *UserAdminServiceImpl) CreateUser(ctx context.Context, req *CreateUserRequest) (*users.User, error) {
	// 验证角色
	if !isValidRole(req.Role) {
		return nil, ErrInvalidRole
	}

	// 检查邮箱是否已存在
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 设置默认密码
	password := req.Password
	if password == "" {
		var err error
		password, err = generateRandomPassword(12)
		if err != nil {
			return nil, fmt.Errorf("failed to generate password: %w", err)
		}
	}

	// 创建用户
	user := &users.User{
		Username:      req.Username,
		Email:         req.Email,
		Nickname:      req.Nickname,
		Roles:         []string{req.Role},
		Status:        req.Status,
		Bio:           req.Bio,
		EmailVerified: false,
	}

	if user.Status == "" {
		user.Status = users.UserStatusActive
	}
	if user.Nickname == "" {
		user.Nickname = req.Username
	}

	// 设置密码
	if err := user.SetPassword(password); err != nil {
		return nil, fmt.Errorf("failed to set password: %w", err)
	}

	// 保存用户
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// BatchCreateUsers 批量创建用户
func (s *UserAdminServiceImpl) BatchCreateUsers(ctx context.Context, req *BatchCreateUserRequest) ([]*users.User, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	if req.Count < 1 || req.Count > maxBatchCreateUsers {
		return nil, ErrInvalidBatchCount
	}

	// 验证角色
	if !isValidRole(req.Role) {
		return nil, ErrInvalidRole
	}

	// 设置默认值
	prefix := req.Prefix
	if prefix == "" {
		prefix = "batch_user"
	}
	status := req.Status
	if status == "" {
		status = users.UserStatusActive
	}

	// 获取当前最大用户数来确定起始ID
	// 这里简化处理，使用时间戳作为唯一标识
	baseID := time.Now().Unix()

	count := req.Count
	usersList := make([]*users.User, count)
	for i := 0; i < count; i++ {
		userID := baseID + int64(i)

		// 生成随机密码
		password, err := generateRandomPassword(12)
		if err != nil {
			return nil, fmt.Errorf("failed to generate password: %w", err)
		}

		user := &users.User{
			Username:      fmt.Sprintf("%s_%d", prefix, userID),
			Email:         fmt.Sprintf("%s_%d@example.com", prefix, userID),
			Nickname:      fmt.Sprintf("批量用户%d", i+1),
			Roles:         []string{req.Role},
			Status:        status,
			Bio:           "批量创建的用户",
			EmailVerified: false,
		}

		if err := user.SetPassword(password); err != nil {
			return nil, fmt.Errorf("failed to set password: %w", err)
		}

		usersList[i] = user
	}

	// 批量保存
	if err := s.userRepo.BatchCreate(ctx, usersList); err != nil {
		return nil, fmt.Errorf("failed to batch create users: %w", err)
	}

	return usersList, nil
}
