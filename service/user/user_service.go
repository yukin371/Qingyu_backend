package user

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	user2 "Qingyu_backend/service/interfaces/user"
	"context"
	"fmt"
	"time"

	"Qingyu_backend/middleware"
	usersModel "Qingyu_backend/models/users"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"
)

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepo repoInterfaces.UserRepository
	name     string
	version  string
}

// NewUserService 创建用户服务
func NewUserService(userRepo repoInterfaces.UserRepository) user2.UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
		name:     "UserService",
		version:  "1.0.0",
	}
}

// Initialize 初始化服务
func (s *UserServiceImpl) Initialize(ctx context.Context) error {
	return s.userRepo.Health(ctx)
}

// Health 健康检查
func (s *UserServiceImpl) Health(ctx context.Context) error {
	return s.userRepo.Health(ctx)
}

// Close 关闭服务
func (s *UserServiceImpl) Close(ctx context.Context) error {
	return nil
}

// GetServiceName 获取服务名称
func (s *UserServiceImpl) GetServiceName() string {
	return s.name
}

// GetVersion 获取服务版本
func (s *UserServiceImpl) GetVersion() string {
	return s.version
}

// CreateUser 创建用户
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user2.CreateUserRequest) (*user2.CreateUserResponse, error) {
	// 1. 验证请求数据
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "请求数据验证失败", err)
	}

	// 2. 检查用户是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "检查用户名失败", err)
	}
	if exists {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeBusiness, "用户名已存在", nil)
	}

	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "检查邮箱失败", err)
	}
	if exists {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeBusiness, "邮箱已存在", nil)
	}

	// 3. 创建用户对象
	user := &usersModel.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// 4. 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "设置密码失败", err)
	}

	// 5. 保存到数据库
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "创建用户失败", err)
	}

	return &user2.CreateUserResponse{
		User: user,
	}, nil
}

// GetUser 获取用户
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user2.GetUserRequest) (*user2.GetUserResponse, error) {
	// 1. 验证请求数据
	if req.ID == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户ID不能为空", nil)
	}

	// 2. 从数据库获取用户
	user, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		if repoInterfaces.IsNotFoundError(err) {
			return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
		}
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户失败", err)
	}

	return &user2.GetUserResponse{
		User: user,
	}, nil
}

// UpdateUser 更新用户
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *user2.UpdateUserRequest) (*user2.UpdateUserResponse, error) {
	// 1. 验证请求数据
	if req.ID == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户ID不能为空", nil)
	}
	if len(req.Updates) == 0 {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "更新数据不能为空", nil)
	}

	// 2. 检查用户是否存在
	exists, err := s.userRepo.Exists(ctx, req.ID)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "检查用户存在性失败", err)
	}
	if !exists {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", nil)
	}

	// 3. 更新用户信息
	if err := s.userRepo.Update(ctx, req.ID, req.Updates); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "更新用户失败", err)
	}

	// 4. 获取更新后的用户信息
	updatedUser, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取更新后的用户信息失败", err)
	}

	return &user2.UpdateUserResponse{
		User: *updatedUser,
	}, nil
}

// DeleteUser 删除用户
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *user2.DeleteUserRequest) (*user2.DeleteUserResponse, error) {
	// 1. 验证请求数据
	if req.ID == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户ID不能为空", nil)
	}

	// 2. 检查用户是否存在
	exists, err := s.userRepo.Exists(ctx, req.ID)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "检查用户存在性失败", err)
	}
	if !exists {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", nil)
	}

	// 3. 删除用户
	if err := s.userRepo.Delete(ctx, req.ID); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "删除用户失败", err)
	}

	return &user2.DeleteUserResponse{
		Deleted:   true,
		DeletedAt: time.Now(),
	}, nil
}

// ListUsers 列出用户
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *user2.ListUsersRequest) (*user2.ListUsersResponse, error) {
	// 1. 构建过滤器
	filter := repoInterfaces.UserFilter{
		Username: req.Username,
		Email:    req.Email,
		Status:   req.Status,
		FromDate: req.FromDate,
		ToDate:   req.ToDate,
	}

	// 2. 设置分页参数
	if req.PageSize > 0 {
		filter.Limit = int64(req.PageSize)
	}
	if req.Page > 0 {
		filter.Offset = int64((req.Page - 1) * req.PageSize)
	}

	// 3. 获取用户列表
	users, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户列表失败", err)
	}

	// 4. 获取总数
	total, err := s.userRepo.Count(ctx, filter)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户总数失败", err)
	}

	// 5. 转换用户列表类型
	var userList []*usersModel.User
	userList = append(userList, users...)

	// 6. 计算总页数
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &user2.ListUsersResponse{
		Users:      userList,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// RegisterUser 注册用户
func (s *UserServiceImpl) RegisterUser(ctx context.Context, req *user2.RegisterUserRequest) (*user2.RegisterUserResponse, error) {
	// 1. 验证请求数据
	if err := s.validateRegisterUserRequest(req); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "请求数据验证失败", err)
	}

	// 2. 检查用户是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "检查用户名失败", err)
	}
	if exists {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeBusiness, "用户名已存在", nil)
	}

	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "检查邮箱失败", err)
	}
	if exists {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeBusiness, "邮箱已存在", nil)
	}

	// 3. 创建用户对象
	user := &usersModel.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",                      // 默认角色
		Status:   usersModel.UserStatusActive, // 默认状态
	}

	// 4. 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "设置密码失败", err)
	}

	// 5. 保存到数据库
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "创建用户失败", err)
	}

	// 6. 生成JWT令牌
	token, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "生成Token失败", err)
	}

	return &user2.RegisterUserResponse{
		User:  user,
		Token: token,
	}, nil
}

// LoginUser 登录用户
func (s *UserServiceImpl) LoginUser(ctx context.Context, req *user2.LoginUserRequest) (*user2.LoginUserResponse, error) {
	// 1. 验证请求数据
	if req.Username == "" || req.Password == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户名和密码不能为空", nil)
	}

	// DEBUG: 记录登录尝试
	fmt.Printf("[DEBUG] 登录尝试 - 用户名: %s\n", req.Username)

	// 2. 获取用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		fmt.Printf("[DEBUG] 获取用户失败 - 错误: %v\n", err)
		if repoInterfaces.IsNotFoundError(err) {
			return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
		}
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户失败", err)
	}

	fmt.Printf("[DEBUG] 用户找到 - ID: %s, 用户名: %s, 状态: %s\n", user.ID, user.Username, user.Status)
	fmt.Printf("[DEBUG] 密码哈希: %s\n", user.Password[:20]+"...")
	fmt.Printf("[DEBUG] 输入密码: %s\n", req.Password)

	// 3. 验证密码
	if !user.ValidatePassword(req.Password) {
		fmt.Printf("[DEBUG] 密码验证失败\n")
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeUnauthorized, "密码错误", nil)
	}

	fmt.Printf("[DEBUG] 密码验证成功\n")

	// 4. 检查用户状态
	switch user.Status {
	case usersModel.UserStatusInactive:
		return nil, serviceInterfaces.NewServiceError(
			s.name,
			serviceInterfaces.ErrorTypeUnauthorized,
			"账号未激活，请先验证邮箱",
			nil,
		)
	case usersModel.UserStatusBanned:
		return nil, serviceInterfaces.NewServiceError(
			s.name,
			serviceInterfaces.ErrorTypeUnauthorized,
			"账号已被封禁，请联系管理员",
			nil,
		)
	case usersModel.UserStatusDeleted:
		return nil, serviceInterfaces.NewServiceError(
			s.name,
			serviceInterfaces.ErrorTypeUnauthorized,
			"账号已删除",
			nil,
		)
	case usersModel.UserStatusActive:
		// 允许登录，继续执行
	default:
		return nil, serviceInterfaces.NewServiceError(
			s.name,
			serviceInterfaces.ErrorTypeInternal,
			"未知的用户状态",
			nil,
		)
	}

	// 5. 更新最后登录时间
	// IP 地址应该从 context 中获取，这里暂时使用默认值
	ip := "unknown" // TODO: 从 context 中获取客户端 IP
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID, ip); err != nil {
		// 记录错误但不影响登录流程
		// 注意：不要使用fmt.Printf，会污染HTTP响应
		// TODO: 使用logger记录错误
		_ = err
	}

	// 6. 生成JWT令牌
	token, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "生成Token失败", err)
	}

	return &user2.LoginUserResponse{
		User:  user,
		Token: token,
	}, nil
}

// LogoutUser 登出用户
func (s *UserServiceImpl) LogoutUser(ctx context.Context, req *user2.LogoutUserRequest) (*user2.LogoutUserResponse, error) {
	// 这里简化处理，实际应该将令牌加入黑名单
	// TODO: 实现JWT令牌黑名单机制
	return &user2.LogoutUserResponse{
		Success: true,
	}, nil
}

// ValidateToken 验证令牌
func (s *UserServiceImpl) ValidateToken(ctx context.Context, req *user2.ValidateTokenRequest) (*user2.ValidateTokenResponse, error) {
	// 这里简化处理，实际应该验证JWT令牌
	// TODO: 实现JWT令牌验证
	return &user2.ValidateTokenResponse{
		Valid: false, // 暂时返回false
	}, nil
}

// UpdateLastLogin 更新最后登录时间
func (s *UserServiceImpl) UpdateLastLogin(ctx context.Context, req *user2.UpdateLastLoginRequest) (*user2.UpdateLastLoginResponse, error) {
	// 1. 验证请求数据
	if req.ID == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户ID不能为空", nil)
	}

	// 2. 更新最后登录时间
	// IP 地址应该从 context 中获取，这里暂时使用默认值
	ip := "unknown" // TODO: 从 context 中获取客户端 IP
	if err := s.userRepo.UpdateLastLogin(ctx, req.ID, ip); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "更新最后登录时间失败", err)
	}

	return &user2.UpdateLastLoginResponse{
		Updated: true,
	}, nil
}

// UpdatePassword 更新密码
func (s *UserServiceImpl) UpdatePassword(ctx context.Context, req *user2.UpdatePasswordRequest) (*user2.UpdatePasswordResponse, error) {
	// 1. 验证请求数据
	if err := s.validateUpdatePasswordRequest(req); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "请求数据验证失败", err)
	}

	// 2. 获取用户
	user, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		if repoInterfaces.IsNotFoundError(err) {
			return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
		}
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户失败", err)
	}

	// 3. 验证旧密码
	if !(*user).ValidatePassword(req.OldPassword) {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeUnauthorized, "旧密码错误", nil)
	}

	// 4. 更新密码
	(*user).Password = req.NewPassword
	if err := (*user).SetPassword(req.NewPassword); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "密码哈希失败", err)
	}
	hashedPassword := (*user).Password

	if err := s.userRepo.UpdatePassword(ctx, req.ID, hashedPassword); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "更新密码失败", err)
	}

	return &user2.UpdatePasswordResponse{
		Updated: true,
	}, nil
}

// ResetPassword 重置密码
func (s *UserServiceImpl) ResetPassword(ctx context.Context, req *user2.ResetPasswordRequest) (*user2.ResetPasswordResponse, error) {
	// 1. 验证请求数据
	if req.Email == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "邮箱不能为空", nil)
	}

	// 2. 检查用户是否存在
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if repoInterfaces.IsNotFoundError(err) {
			// 为了安全，即使用户不存在也返回成功
			return &user2.ResetPasswordResponse{
				Success: true,
			}, nil
		}
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "检查用户失败", err)
	}

	// 3. 生成新密码（这里简化处理，实际应该发送邮件）
	newPassword := "new_password_placeholder" // TODO: 实现密码重置邮件发送

	// 4. 更新密码
	(*user).Password = newPassword
	if err := (*user).SetPassword(newPassword); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "密码哈希失败", err)
	}
	hashedPassword := (*user).Password

	if err := s.userRepo.UpdatePassword(ctx, (*user).ID, hashedPassword); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "更新密码失败", err)
	}

	return &user2.ResetPasswordResponse{
		Success: true,
	}, nil
}

// AssignRole 分配角色
func (s *UserServiceImpl) AssignRole(ctx context.Context, req *user2.AssignRoleRequest) (*user2.AssignRoleResponse, error) {
	// TODO: 实现角色分配逻辑
	return &user2.AssignRoleResponse{
		Assigned: false, // 暂时返回false
	}, nil
}

// RemoveRole 移除角色
func (s *UserServiceImpl) RemoveRole(ctx context.Context, req *user2.RemoveRoleRequest) (*user2.RemoveRoleResponse, error) {
	// TODO: 实现角色移除逻辑
	return &user2.RemoveRoleResponse{
		Removed: false, // 暂时返回false
	}, nil
}

// GetUserRoles 获取用户角色
func (s *UserServiceImpl) GetUserRoles(ctx context.Context, req *user2.GetUserRolesRequest) (*user2.GetUserRolesResponse, error) {
	// TODO: 实现获取用户角色逻辑
	return &user2.GetUserRolesResponse{
		Roles: []string{}, // 暂时返回空列表
	}, nil
}

// GetUserPermissions 获取用户权限
func (s *UserServiceImpl) GetUserPermissions(ctx context.Context, req *user2.GetUserPermissionsRequest) (*user2.GetUserPermissionsResponse, error) {
	// TODO: 实现获取用户权限逻辑
	return &user2.GetUserPermissionsResponse{
		Permissions: []string{}, // 暂时返回空列表
	}, nil
}

// 私有方法

// validateCreateUserRequest 验证创建用户请求
func (s *UserServiceImpl) validateCreateUserRequest(req *user2.CreateUserRequest) error {
	if req.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if req.Email == "" {
		return fmt.Errorf("邮箱不能为空")
	}
	if req.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	return nil
}

// validateRegisterUserRequest 验证注册用户请求
func (s *UserServiceImpl) validateRegisterUserRequest(req *user2.RegisterUserRequest) error {
	if req.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if req.Email == "" {
		return fmt.Errorf("邮箱不能为空")
	}
	if req.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	return nil
}

// validateUpdatePasswordRequest 验证更新密码请求
func (s *UserServiceImpl) validateUpdatePasswordRequest(req *user2.UpdatePasswordRequest) error {
	if req.ID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if req.OldPassword == "" {
		return fmt.Errorf("旧密码不能为空")
	}
	if req.NewPassword == "" {
		return fmt.Errorf("新密码不能为空")
	}
	return nil
}

// generateToken 生成JWT令牌（辅助方法）
func (s *UserServiceImpl) generateToken(userID, role string) (string, error) {
	// 使用middleware包中的GenerateToken函数
	// 导入: "Qingyu_backend/middleware"
	return middleware.GenerateToken(userID, "", []string{role})
}
