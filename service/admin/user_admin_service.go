package admin

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"Qingyu_backend/models/admin"
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
	// ErrBanReasonRequired 封禁时必须提供原因
	ErrBanReasonRequired = fmt.Errorf("ban reason is required when banning user")
)

const (
	maxBatchCreateUsers = 100
)

// MaxBatchCreateUsersCount 批量创建用户的最大数量
const MaxBatchCreateUsersCount = 1000

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

	// UpdateUserStatusWithReason 更新用户状态（带封禁原因记录）
	UpdateUserStatusWithReason(ctx context.Context, userID string, status users.UserStatus, operatorID string, banReason *string) error

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
	userRepo     adminrepo.UserAdminRepository
	banRecordRepo adminrepo.BanRecordRepository
}

// NewUserAdminService 创建用户管理服务
func NewUserAdminService(userRepo adminrepo.UserAdminRepository) UserAdminService {
	return &UserAdminServiceImpl{
		userRepo:     userRepo,
		banRecordRepo: nil, // 需要通过依赖注入设置
	}
}

// NewUserAdminServiceWithBanRepo 创建带封禁记录仓储的用户管理服务
func NewUserAdminServiceWithBanRepo(userRepo adminrepo.UserAdminRepository, banRecordRepo adminrepo.BanRecordRepository) UserAdminService {
	return &UserAdminServiceImpl{
		userRepo:     userRepo,
		banRecordRepo: banRecordRepo,
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

	// 使用新接口 ListWithPagination
	return s.userRepo.ListWithPagination(ctx, filter, page, pageSize)
}

// GetUserDetail 获取用户详情
func (s *UserAdminServiceImpl) GetUserDetail(ctx context.Context, userID string) (*users.User, error) {
	// 验证 ID 格式
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// UpdateUserStatus 更新用户状态
func (s *UserAdminServiceImpl) UpdateUserStatus(ctx context.Context, userID string, status users.UserStatus) error {
	// 验证 ID 格式
	if userID == "" {
		return ErrInvalidUserID
	}

	// 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 不能修改管理员的状态
	if user.HasRole("admin") {
		return ErrCannotModifySuperAdmin
	}

	return s.userRepo.UpdateStatus(ctx, userID, status)
}

// UpdateUserStatusWithReason 更新用户状态（带封禁原因记录）
func (s *UserAdminServiceImpl) UpdateUserStatusWithReason(ctx context.Context, userID string, status users.UserStatus, operatorID string, banReason *string) error {
	// 验证 ID 格式
	if userID == "" {
		return ErrInvalidUserID
	}

	// 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 不能修改管理员的状态
	if user.HasRole("admin") {
		return ErrCannotModifySuperAdmin
	}

	// 封禁时必须提供原因
	if status == users.UserStatusBanned && (banReason == nil || *banReason == "") {
		return ErrBanReasonRequired
	}

	oldStatus := user.Status
	now := time.Now()

	// 准备更新字段
	updates := make(map[string]interface{})

	if status == users.UserStatusBanned {
		// 封禁逻辑：设置封禁字段
		updates["banned_at"] = &now
		updates["banned_by"] = operatorID
		if banReason != nil {
			updates["ban_reason"] = *banReason
		}
	} else if oldStatus == users.UserStatusBanned {
		// 解封：清除封禁字段
		updates["banned_at"] = nil
		updates["banned_by"] = ""
		updates["ban_reason"] = ""
	}

	// 更新用户状态
	if err := s.userRepo.UpdateStatus(ctx, userID, status); err != nil {
		return err
	}

	// 更新其他字段（如果有）
	if len(updates) > 0 {
		if err := s.userRepo.Update(ctx, userID, updates); err != nil {
			return err
		}
	}

	// 记录封禁历史
	if s.banRecordRepo != nil {
		s.recordBanHistory(ctx, userID, status, operatorID, banReason)
	}

	return nil
}

// recordBanHistory 记录封禁历史
func (s *UserAdminServiceImpl) recordBanHistory(ctx context.Context, userID string, status users.UserStatus, operatorID string, banReason *string) {
	action := "unban"
	reason := "解除封禁"
	if status == users.UserStatusBanned {
		action = "ban"
		if banReason != nil {
			reason = *banReason
		}
	}

	record := &admin.BanRecord{
		UserID:     userID,
		Action:     action,
		Reason:     reason,
		OperatorID: operatorID,
	}

	// 忽略错误，记录失败不应影响主流程
	_ = s.banRecordRepo.Create(ctx, record)
}

// UpdateUserRole 更新用户角色
func (s *UserAdminServiceImpl) UpdateUserRole(ctx context.Context, userID, role string) error {
	// 验证 ID 格式
	if userID == "" {
		return ErrInvalidUserID
	}

	// 验证角色
	if !isValidRole(role) {
		return ErrInvalidRole
	}

	// 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 不能修改管理员的角色
	if user.HasRole("admin") {
		return ErrCannotModifySuperAdmin
	}

	// 注意：UpdateRoles 仍使用 ObjectID，需要转换
	// 这里是一个临时方案，实际上应该让 Repository 也接受 string ID
	// 但为了保持 admin 特有方法的一致性，这里保留 ObjectID
	// TODO: 将 UpdateRoles 也改为接受 string ID
	return s.userRepo.UpdateRoles(ctx, user.ID, role)
}

// DeleteUser 删除用户
func (s *UserAdminServiceImpl) DeleteUser(ctx context.Context, userID string) error {
	// 验证 ID 格式
	if userID == "" {
		return ErrInvalidUserID
	}

	// 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 不能删除管理员
	if user.HasRole("admin") {
		return ErrCannotModifySuperAdmin
	}

	return s.userRepo.Delete(ctx, userID)
}

// BatchUpdateStatus 批量更新用户状态
func (s *UserAdminServiceImpl) BatchUpdateStatus(ctx context.Context, userIDs []string, status users.UserStatus) error {
	// 过滤有效的 ID
	validIDs := make([]string, 0, len(userIDs))
	for _, idStr := range userIDs {
		if idStr == "" {
			continue // 跳过空ID
		}

		// 检查是否是管理员
		user, err := s.userRepo.GetByID(ctx, idStr)
		if err == nil && user.HasRole("admin") {
			continue // 跳过管理员
		}

		validIDs = append(validIDs, idStr)
	}

	if len(validIDs) == 0 {
		return fmt.Errorf("no valid user IDs")
	}

	// 新接口直接接受 []string
	return s.userRepo.BatchUpdateStatus(ctx, validIDs, status)
}

// BatchDeleteUsers 批量删除用户
func (s *UserAdminServiceImpl) BatchDeleteUsers(ctx context.Context, userIDs []string) error {
	return s.BatchUpdateStatus(ctx, userIDs, users.UserStatusDeleted)
}

// GetUserActivities 获取用户活动记录
func (s *UserAdminServiceImpl) GetUserActivities(ctx context.Context, userID string, page, pageSize int) ([]*users.UserActivity, int64, error) {
	// 验证 ID 格式
	if userID == "" {
		return nil, 0, ErrInvalidUserID
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// GetActivities 仍使用 ObjectID，需要转换
	// TODO: 将 GetActivities 也改为接受 string ID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, 0, ErrUserNotFound
	}

	return s.userRepo.GetActivities(ctx, user.ID, page, pageSize)
}

// GetUserStatistics 获取用户统计信息
func (s *UserAdminServiceImpl) GetUserStatistics(ctx context.Context, userID string) (*users.UserStatistics, error) {
	// 验证 ID 格式
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	// GetStatistics 仍使用 ObjectID，需要转换
	// TODO: 将 GetStatistics 也改为接受 string ID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.userRepo.GetStatistics(ctx, user.ID)
}

// ResetUserPassword 重置用户密码
func (s *UserAdminServiceImpl) ResetUserPassword(ctx context.Context, userID string) (string, error) {
	// 验证 ID 格式
	if userID == "" {
		return "", ErrInvalidUserID
	}

	// 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", ErrUserNotFound
	}

	// 生成随机密码
	newPassword, err := generateRandomPassword(12)
	if err != nil {
		return "", fmt.Errorf("failed to generate password: %w", err)
	}

	// ResetPassword 仍使用 ObjectID，需要转换
	// TODO: 将 ResetPassword 也改为接受 string ID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", ErrUserNotFound
	}

	if err := s.userRepo.ResetPassword(ctx, user.ID, newPassword); err != nil {
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
	// 使用新接口 CountByStatusMap
	return s.userRepo.CountByStatusMap(ctx)
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

	// 校验批量创建数量，防止过大的内存分配
	if req.Count <= 0 || req.Count > MaxBatchCreateUsersCount {
		return nil, ErrInvalidBatchCount
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
