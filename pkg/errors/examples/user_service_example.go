package user

import (
	"context"
	"fmt"

	"Qingyu_backend/pkg/errors"
)

// ============================================================================
// 用户服务 - 使用统一错误系统
// 
// 本文件展示如何迁移到统一错误系统
// ============================================================================

// UserService 用户服务
type UserService struct {
	repo    UserRepository
	factory *errors.ErrorFactory
}

// NewUserService 创建用户服务
func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo:    repo,
		factory: errors.UserServiceFactory,
	}
}

// ============================================================================
// 用户查询相关方法
// ============================================================================

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*User, error) {
	// 参数验证
	if userID == "" {
		return nil, s.factory.ValidationError(
			"1008",
			"用户ID不能为空",
			"field: user_id",
		)
	}

	// 查询用户
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, s.factory.InternalError(
			"5001",
			"查询用户失败",
			err,
		)
	}

	// 检查用户是否存在
	if user == nil {
		return nil, s.factory.NotFoundError("用户", userID)
	}

	return user, nil
}

// GetByEmail 根据邮箱获取用户
func (s *UserService) GetByEmail(ctx context.Context, email string) (*User, error) {
	// 参数验证
	if email == "" {
		return nil, s.factory.ValidationError(
			"1008",
			"邮箱不能为空",
			"field: email",
		)
	}

	// 查询用户
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, s.factory.InternalError(
			"5001",
			"查询用户失败",
			err,
		)
	}

	// 检查用户是否存在
	if user == nil {
		return nil, s.factory.NotFoundError("用户", "邮箱:"+email)
	}

	return user, nil
}

// ============================================================================
// 用户创建相关方法
// ============================================================================

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
	// 验证请求参数
	if err := s.validateCreateUserRequest(ctx, req); err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, s.factory.InternalError(
			"5001",
			"检查用户名失败",
			err,
		)
	}
	if exists {
		return nil, s.factory.BusinessError(
			"2003",
			"用户名已被使用",
			"用户名 '"+req.Username+"' 已被使用，请选择其他用户名",
		)
	}

	// 检查邮箱是否已存在
	exists, err = s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, s.factory.InternalError(
			"5001",
			"检查邮箱失败",
			err,
		)
	}
	if exists {
		return nil, s.factory.BusinessError(
			"2004",
			"邮箱已被使用",
			"邮箱 '"+req.Email+"' 已被使用，请使用其他邮箱",
		)
	}

	// 创建用户对象
	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 应该是加密后的密码
	}

	// 保存到数据库
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, s.factory.InternalError(
			"5001",
			"创建用户失败",
			err,
		)
	}

	return user, nil
}

// validateCreateUserRequest 验证创建用户请求
func (s *UserService) validateCreateUserRequest(ctx context.Context, req *CreateUserRequest) error {
	// 检查请求是否为空
	if req == nil {
		return s.factory.ValidationError(
			"1001",
			"请求参数无效",
			"请求体不能为空",
		)
	}

	// 验证用户名
	if req.Username == "" {
		return s.factory.ValidationError(
			"1008",
			"用户名不能为空",
			"field: username",
		)
	}

	if len(req.Username) < 3 {
		return s.factory.ValidationError(
			"1001",
			"用户名长度无效",
			"用户名长度必须至少为3个字符",
		)
	}

	if len(req.Username) > 20 {
		return s.factory.ValidationError(
			"1001",
			"用户名长度无效",
			"用户名长度不能超过20个字符",
		)
	}

	// 验证邮箱
	if req.Email == "" {
		return s.factory.ValidationError(
			"1008",
			"邮箱不能为空",
			"field: email",
		)
	}

	// 这里应该有邮箱格式验证
	// if !emailRegex.MatchString(req.Email) {
	//     return s.factory.ValidationError(
	//         "1009",
	//         "邮箱格式无效",
	//         "field: email",
	//     )
	// }

	// 验证密码
	if req.Password == "" {
		return s.factory.ValidationError(
			"1008",
			"密码不能为空",
			"field: password",
		)
	}

	if len(req.Password) < 6 {
		return s.factory.ValidationError(
			"2014",
			"密码强度不足",
			"密码长度必须至少为6个字符",
		)
	}

	return nil
}

// ============================================================================
// 用户认证相关方法
// ============================================================================

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*User, string, error) {
	// 验证请求参数
	if req == nil || req.Username == "" || req.Password == "" {
		return nil, "", s.factory.ValidationError(
			"1008",
			"用户名或密码不能为空",
			"fields: username, password",
		)
	}

	// 查询用户
	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, "", s.factory.InternalError(
			"5001",
			"查询用户失败",
			err,
		)
	}

	// 检查用户是否存在
	if user == nil {
		return nil, "", s.factory.AuthError(
			"2002",
			"用户名或密码错误",
		)
	}

	// 验证密码
	if !s.verifyPassword(user.Password, req.Password) {
		return nil, "", s.factory.AuthError(
			"2002",
			"用户名或密码错误",
		)
	}

	// 检查账户状态
	if user.Status == "locked" {
		return nil, "", s.factory.ForbiddenError(
			"2015",
			"账户已被锁定，请联系管理员",
		)
	}

	if user.Status == "disabled" {
		return nil, "", s.factory.ForbiddenError(
			"2016",
			"账户已被禁用",
		)
	}

	// 生成Token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", s.factory.InternalError(
			"5000",
			"生成Token失败",
			err,
		)
	}

	return user, token, nil
}

// ============================================================================
// 用户更新相关方法
// ============================================================================

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, userID string, req *UpdateUserRequest) (*User, error) {
	// 获取用户
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}

	if req.Bio != nil {
		user.Bio = *req.Bio
	}

	// 如果要更新用户名，检查是否已被使用
	if req.Username != nil && *req.Username != user.Username {
		exists, err := s.repo.ExistsByUsername(ctx, *req.Username)
		if err != nil {
			return nil, s.factory.InternalError(
				"5001",
				"检查用户名失败",
				err,
			)
		}
		if exists {
			return nil, s.factory.BusinessError(
				"2003",
				"用户名已被使用",
			)
		}
		user.Username = *req.Username
	}

	// 如果要更新邮箱，检查是否已被使用
	if req.Email != nil && *req.Email != user.Email {
		exists, err := s.repo.ExistsByEmail(ctx, *req.Email)
		if err != nil {
			return nil, s.factory.InternalError(
				"5001",
				"检查邮箱失败",
				err,
			)
		}
		if exists {
			return nil, s.factory.BusinessError(
				"2004",
				"邮箱已被使用",
			)
		}
		user.Email = *req.Email
	}

	// 保存更新
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, s.factory.InternalError(
			"5001",
			"更新用户失败",
			err,
		)
	}

	return user, nil
}

// ============================================================================
// 用户删除相关方法
// ============================================================================

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	// 检查用户是否存在
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// 删除用户
	if err := s.repo.Delete(ctx, userID); err != nil {
		return s.factory.InternalError(
			"5001",
			"删除用户失败",
			err,
		)
	}

	_ = user // 避免未使用变量警告
	return nil
}

// ============================================================================
// 辅助方法
// ============================================================================

// verifyPassword 验证密码
func (s *UserService) verifyPassword(hashedPassword, password string) bool {
	// 实际实现应该使用bcrypt等加密算法
	// 这里简化处理
	return hashedPassword == password
}

// generateToken 生成Token
func (s *UserService) generateToken(user *User) (string, error) {
	// 实际实现应该使用JWT等
	// 这里简化处理
	return fmt.Sprintf("token_%s", user.ID), nil
}

// ============================================================================
// 请求和响应结构
// ============================================================================

// User 用户模型
type User struct {
	ID       string
	Username string
	Email    string
	Password string
	Nickname string
	Bio      string
	Status   string
}

// UserRepository 用户仓库接口
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Nickname *string `json:"nickname,omitempty"`
	Bio      *string `json:"bio,omitempty"`
}
