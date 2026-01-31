package user

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	user2 "Qingyu_backend/service/interfaces/user"
	"context"
	"fmt"
	"time"

	"Qingyu_backend/internal/middleware/auth"
	usersModel "Qingyu_backend/models/users"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"

	"go.uber.org/zap"
)

// isDuplicateKeyError 检查是否是MongoDB唯一索引冲突错误
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// MongoDB Duplicate Key错误通常包含"11000"或"duplicate key"
	return contains(errStr, "11000") || contains(errStr, "duplicate key") ||
	       contains(errStr, "E11000")
}

// contains 检查字符串是否包含子串（忽略大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (
		s[:len(substr)] == substr ||
		s[len(s)-len(substr):] == substr ||
		indexOf(s, substr) >= 0))
}

// indexOf 查找子串位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepo repoInterfaces.UserRepository
	authRepo sharedRepo.AuthRepository
	name     string
	version  string
}

// NewUserService 创建用户服务
func NewUserService(userRepo repoInterfaces.UserRepository, authRepo sharedRepo.AuthRepository) user2.UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
		authRepo: authRepo,
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
		Status:   usersModel.UserStatusActive, // 默认设置为活跃状态
		Roles:    []string{"reader"}, // 默认角色为reader
	}

	// 4. 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "设置密码失败", err)
	}

	// 5. 保存到数据库（带重试机制处理并发冲突）
	maxRetries := 3
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := s.userRepo.Create(ctx, user)
		if err == nil {
			// 创建成功
			break
		}

		lastErr = err

		// 检查是否是唯一索引冲突（MongoDB错误代码11000）
		if isDuplicateKeyError(err) {
			// 重新检查用户是否已存在（可能是其他并发请求创建了）
			exists, checkErr := s.userRepo.ExistsByUsername(ctx, req.Username)
			if checkErr == nil && exists {
				// 用户确实存在了，返回友好错误
				return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeBusiness, "用户名已存在", nil)
			}

			// 如果是最后一次尝试，返回错误
			if attempt == maxRetries-1 {
				return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "创建用户失败（并发冲突）", err)
			}

			// 等待一小段时间后重试
			time.Sleep(time.Duration(attempt+1) * 50 * time.Millisecond)
			continue
		}

		// 其他类型的错误，直接返回
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "创建用户失败", err)
	}

	// 如果所有重试都失败，返回最后一次错误
	if lastErr != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "创建用户失败", lastErr)
	}

	return &user2.CreateUserResponse{
		User: ToUserDTO(user),
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
		User: ToUserDTO(user),
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
		User: ToUserDTO(updatedUser),
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

	// 5. 转换用户列表为 DTO
	userDTOs := ToUserDTOs(users)

	// 6. 计算总页数
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &user2.ListUsersResponse{
		Users:      userDTOs,
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
		Roles:    []string{"reader"},         // 默认角色
		Status:   usersModel.UserStatusActive, // 默认状态
	}

	// 4. 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "设置密码失败", err)
	}

	// 5. 保存到数据库（带重试机制处理并发冲突）
	maxRetries := 3
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := s.userRepo.Create(ctx, user)
		if err == nil {
			// 创建成功
			break
		}

		lastErr = err

		// 检查是否是唯一索引冲突（MongoDB错误代码11000）
		if isDuplicateKeyError(err) {
			// 重新检查用户是否已存在（可能是其他并发请求创建了）
			exists, checkErr := s.userRepo.ExistsByUsername(ctx, req.Username)
			if checkErr == nil && exists {
				// 用户确实存在了，返回友好错误
				return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeBusiness, "用户名已存在", nil)
			}

			// 如果是最后一次尝试，返回错误
			if attempt == maxRetries-1 {
				return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "创建用户失败（并发冲突）", err)
			}

			// 等待一小段时间后重试
			time.Sleep(time.Duration(attempt+1) * 50 * time.Millisecond)
			continue
		}

		// 其他类型的错误，直接返回
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "创建用户失败", err)
	}

	// 如果所有重试都失败，返回最后一次错误
	if lastErr != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "创建用户失败", lastErr)
	}

	// 6. 生成JWT令牌 - 使用用户的实际角色列表
	token, err := s.generateToken(user.ID.Hex(), user.Roles)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "生成Token失败", err)
	}

	return &user2.RegisterUserResponse{
		User:  ToUserDTO(user),
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
	zap.L().Debug("登录尝试", zap.String("username", req.Username))

	// 2. 获取用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		zap.L().Debug("获取用户失败", zap.Error(err))
		if repoInterfaces.IsNotFoundError(err) {
			return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
		}
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户失败", err)
	}

	zap.L().Debug("用户找到",
		zap.String("user_id", user.ID.Hex()),
		zap.String("username", user.Username),
		zap.String("status", string(user.Status)))
	// 安全地截取密码哈希前缀
	hashPrefix := "..."
	if len(user.Password) > 20 {
		hashPrefix = user.Password[:20] + "..."
	} else {
		hashPrefix = user.Password
	}
	zap.L().Debug("密码哈希", zap.String("hash_prefix", hashPrefix))
	zap.L().Debug("输入密码", zap.String("password", req.Password))

	// 3. 验证密码
	if !user.ValidatePassword(req.Password) {
		zap.L().Debug("密码验证失败")
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeUnauthorized, "密码错误", nil)
	}

	zap.L().Debug("密码验证成功")

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
	ip := req.ClientIP
	if ip == "" {
		ip = "unknown"
	}
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID.Hex(), ip); err != nil {
		// 记录错误但不影响登录流程
		zap.L().Warn("更新最后登录时间失败",
			zap.String("user_id", user.ID.Hex()),
			zap.String("ip", ip),
			zap.Error(err),
		)
	}

	// 6. 生成JWT令牌 - 使用用户的实际角色列表
	token, err := s.generateToken(user.ID.Hex(), user.Roles)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "生成Token失败", err)
	}

	return &user2.LoginUserResponse{
		User:  ToUserDTO(user),
		Token: token,
	}, nil
}

// LogoutUser 登出用户
func (s *UserServiceImpl) LogoutUser(ctx context.Context, req *user2.LogoutUserRequest) (*user2.LogoutUserResponse, error) {
	// 注意：完整的实现应该：
	// 1. 将 JWT Token 加入黑名单（Redis）
	// 2. 设置过期时间等于 Token 的剩余有效期
	// 当前实现：简化处理，仅返回成功
	// TODO(Production): 集成 TokenBlacklistRepository
	// if s.tokenBlacklistRepo != nil {
	// 	err := s.tokenBlacklistRepo.AddToBlacklist(ctx, req.Token, tokenExpiry)
	// 	if err != nil {
	// 		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "加入黑名单失败", err)
	// 	}
	// }
	return &user2.LogoutUserResponse{
		Success: true,
	}, nil
}

// ValidateToken 验证令牌
func (s *UserServiceImpl) ValidateToken(ctx context.Context, req *user2.ValidateTokenRequest) (*user2.ValidateTokenResponse, error) {
	// 注意：完整的实现应该：
	// 1. 验证 JWT Token 的签名和过期时间
	// 2. 检查 Token 是否在黑名单中
	// 3. 返回 Token 中的用户信息
	// 当前实现：简化处理，JWT 验证在中间件中完成
	// TODO(Production): 集成 JWT 验证库和黑名单检查
	// if s.tokenBlacklistRepo != nil {
	// 	isBlacklisted, _ := s.tokenBlacklistRepo.IsBlacklisted(ctx, req.Token)
	// 	if isBlacklisted {
	// 		return &user2.ValidateTokenResponse{Valid: false}, nil
	// 	}
	// }
	return &user2.ValidateTokenResponse{
		Valid: false, // 暂时返回false，实际验证在 JWT 中间件中完成
	}, nil
}

// UpdateLastLogin 更新最后登录时间
func (s *UserServiceImpl) UpdateLastLogin(ctx context.Context, req *user2.UpdateLastLoginRequest) (*user2.UpdateLastLoginResponse, error) {
	// 1. 验证请求数据
	if req.ID == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户ID不能为空", nil)
	}

	// 2. 更新最后登录时间
	// IP 地址应该从 API 层通过参数传递
	// API 层使用 utils.GetClientIP(c) 获取真实客户端 IP
	ip := "unknown" // 默认值，实际应该从请求参数中获取
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
// TODO:注意：这是简化实现。完整的密码重置流程应该包含：
// 1. 用户请求重置 -> 生成Token并发送到邮箱
// 2. 用户点击邮件链接 -> 验证Token
// 3. 用户输入新密码 -> 更新密码
// 当前实现：生成Token并模拟发送邮件（实际未发送）
func (s *UserServiceImpl) ResetPassword(ctx context.Context, req *user2.ResetPasswordRequest) (*user2.ResetPasswordResponse, error) {
	// 1. 验证请求数据
	if req.Email == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "邮箱不能为空", nil)
	}

	// 2. 检查用户是否存在
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if repoInterfaces.IsNotFoundError(err) {
			// 为了安全，即使用户不存在也返回成功（防止邮箱枚举攻击）
			return &user2.ResetPasswordResponse{
				Success: true,
			}, nil
		}
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "检查用户失败", err)
	}

	// 3. 生成密码重置Token
	tokenManager := NewPasswordResetTokenManager()
	resetToken, err := tokenManager.GenerateToken(ctx, req.Email)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "生成重置Token失败", err)
	}

	// 4. 构建重置邮件内容
	resetLink := fmt.Sprintf("https://qingyu.example.com/reset-password?token=%s&email=%s", resetToken, req.Email)
	emailBody := fmt.Sprintf(`
		<h2>密码重置请求</h2>
		<p>您好，%s，</p>
		<p>我们收到了您的密码重置请求。请点击下面的链接重置您的密码：</p>
		<p><a href="%s">重置密码</a></p>
		<p>该链接将在1小时后过期。</p>
		<p>如果您没有请求重置密码，请忽略此邮件。</p>
		<p>青羽写作团队</p>
	`, user.Username, resetLink)

	// 5. 发送重置邮件（当前为模拟发送）
	// 注意：EmailService 需要在 ServiceContainer 中注入
	// TODO(Production): 集成真实的邮件发送服务
	// if s.emailService != nil {
	// 	err := s.emailService.SendEmail(ctx, &messaging.EmailRequest{
	// 		To:      []string{req.Email},
	// 		Subject: "青羽写作 - 密码重置",
	// 		Body:    emailBody,
	// 		IsHTML:  true,
	// 	})
	// 	if err != nil {
	// 		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "发送重置邮件失败", err)
	// 	}
	// }

	// 模拟：打印日志代替发送邮件
	fmt.Printf("[Password Reset] Token generated for %s: %s\n", req.Email, resetToken)
	fmt.Printf("[Password Reset] Email content:\n%s\n", emailBody)

	return &user2.ResetPasswordResponse{
		Success: true,
	}, nil
}

// AssignRole 分配角色
func (s *UserServiceImpl) AssignRole(ctx context.Context, req *user2.AssignRoleRequest) (*user2.AssignRoleResponse, error) {
	// 1. 验证请求数据
	if req.UserID == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户ID不能为空", nil)
	}
	if req.RoleID == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "角色ID不能为空", nil)
	}

	// 2. 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
	}

	// 3. 检查角色是否存在
	_, err = s.authRepo.GetRole(ctx, req.RoleID)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "角色不存在", err)
	}

	// 4. 分配角色
	if err := s.authRepo.AssignUserRole(ctx, req.UserID, req.RoleID); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "分配角色失败", err)
	}

	return &user2.AssignRoleResponse{
		Assigned: true,
	}, nil
}

// RemoveRole 移除角色
func (s *UserServiceImpl) RemoveRole(ctx context.Context, req *user2.RemoveRoleRequest) (*user2.RemoveRoleResponse, error) {
	// 1. 验证请求数据
	if req.UserID == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户ID不能为空", nil)
	}
	if req.RoleID == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "角色ID不能为空", nil)
	}

	// 2. 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
	}

	// 3. 移除角色
	if err := s.authRepo.RemoveUserRole(ctx, req.UserID, req.RoleID); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "移除角色失败", err)
	}

	return &user2.RemoveRoleResponse{
		Removed: true,
	}, nil
}

// GetUserRoles 获取用户角色
func (s *UserServiceImpl) GetUserRoles(ctx context.Context, req *user2.GetUserRolesRequest) (*user2.GetUserRolesResponse, error) {
	// 1. 验证请求数据
	if req.UserID == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户ID不能为空", nil)
	}

	// 2. 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
	}

	// 3. 获取用户角色
	roles, err := s.authRepo.GetUserRoles(ctx, req.UserID)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户角色失败", err)
	}

	// 4. 转换为角色名称列表
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return &user2.GetUserRolesResponse{
		Roles: roleNames,
	}, nil
}

// GetUserPermissions 获取用户权限
func (s *UserServiceImpl) GetUserPermissions(ctx context.Context, req *user2.GetUserPermissionsRequest) (*user2.GetUserPermissionsResponse, error) {
	// 1. 验证请求数据
	if req.UserID == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户ID不能为空", nil)
	}

	// 2. 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
	}

	// 3. 获取用户权限（通过角色获取）
	permissions, err := s.authRepo.GetUserPermissions(ctx, req.UserID)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户权限失败", err)
	}

	return &user2.GetUserPermissionsResponse{
		Permissions: permissions,
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
func (s *UserServiceImpl) generateToken(userID string, roles []string) (string, error) {
	// 使用middleware包中的GenerateToken函数
	// 导入: "Qingyu_backend/middleware"
	return auth.GenerateToken(userID, "", roles)
}

// ==================== 邮箱验证相关方法 ====================

// SendEmailVerification 发送邮箱验证码
func (s *UserServiceImpl) SendEmailVerification(ctx context.Context, req *user2.SendEmailVerificationRequest) (*user2.SendEmailVerificationResponse, error) {
	// 1. 验证请求数据
	if req.UserID == "" || req.Email == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户ID和邮箱不能为空", nil)
	}

	// 2. 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		if repoInterfaces.IsNotFoundError(err) {
			return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
		}
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户失败", err)
	}

	// 3. 验证邮箱是否匹配
	if user.Email != req.Email {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "邮箱不匹配", nil)
	}

	// 4. 如果已经验证过，直接返回
	if user.EmailVerified {
		return &user2.SendEmailVerificationResponse{
			Success:   true,
			Message:   "邮箱已验证",
			ExpiresIn: 0,
		}, nil
	}

	// 5. 生成验证码
	tokenManager := NewEmailVerificationTokenManager()
	code, err := tokenManager.GenerateCode(ctx, req.UserID, req.Email)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "生成验证码失败", err)
	}

	// 6. 构建验证邮件内容
	emailBody := fmt.Sprintf(`
		<h2>邮箱验证</h2>
		<p>您好，%s，</p>
		<p>感谢您注册青羽写作平台！</p>
		<p>您的邮箱验证码是：<strong style="font-size: 24px; color: #007bff;">%s</strong></p>
		<p>验证码有效期为30分钟，请尽快完成验证。</p>
		<p>如果您没有注册青羽写作平台，请忽略此邮件。</p>
		<p>青羽写作团队</p>
	`, user.Username, code)

	// 7. 发送验证邮件（当前为模拟发送）
	// TODO(Production): 集成真实的邮件发送服务
	// if s.emailService != nil {
	// 	err := s.emailService.SendEmail(ctx, &messaging.EmailRequest{
	// 		To:      []string{req.Email},
	// 		Subject: "青羽写作 - 邮箱验证",
	// 		Body:    emailBody,
	// 		IsHTML:  true,
	// 	})
	// 	if err != nil {
	// 		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "发送验证邮件失败", err)
	// 	}
	// }

	// 模拟：打印日志代替发送邮件
	fmt.Printf("[Email Verification] Code generated for %s (UserID: %s): %s\n", req.Email, req.UserID, code)
	fmt.Printf("[Email Verification] Email content:\n%s\n", emailBody)

	return &user2.SendEmailVerificationResponse{
		Success:   true,
		Message:   "验证码已发送到您的邮箱",
		ExpiresIn: 1800, // 30分钟 = 1800秒
	}, nil
}

// VerifyEmail 验证邮箱
func (s *UserServiceImpl) VerifyEmail(ctx context.Context, req *user2.VerifyEmailRequest) (*user2.VerifyEmailResponse, error) {
	// 1. 验证请求数据
	if req.UserID == "" || req.Code == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户ID和验证码不能为空", nil)
	}

	// 2. 获取用户信息
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		if repoInterfaces.IsNotFoundError(err) {
			return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
		}
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户失败", err)
	}

	// 3. 验证验证码
	tokenManager := NewEmailVerificationTokenManager()
	if err := tokenManager.ValidateCode(ctx, req.UserID, user.Email, req.Code); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "验证码无效: "+err.Error(), nil)
	}

	// 4. 标记验证码为已使用
	if err := tokenManager.MarkCodeAsUsed(ctx, user.Email); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "标记验证码失败", err)
	}

	// 5. 更新用户邮箱验证状态
	updates := map[string]interface{}{
		"email_verified": true,
		"status":         usersModel.UserStatusActive, // 激活账户
	}
	if err := s.userRepo.Update(ctx, req.UserID, updates); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "更新用户状态失败", err)
	}

	return &user2.VerifyEmailResponse{
		Success: true,
		Message: "邮箱验证成功",
	}, nil
}

// ==================== 完整密码重置流程相关方法 ====================

// RequestPasswordReset 请求密码重置（发送重置邮件）
func (s *UserServiceImpl) RequestPasswordReset(ctx context.Context, req *user2.RequestPasswordResetRequest) (*user2.RequestPasswordResetResponse, error) {
	// 1. 验证请求数据
	if req.Email == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "邮箱不能为空", nil)
	}

	// 2. 检查用户是否存在
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if repoInterfaces.IsNotFoundError(err) {
			// 为了安全，即使用户不存在也返回成功（防止邮箱枚举攻击）
			return &user2.RequestPasswordResetResponse{
				Success:   true,
				Message:   "如果该邮箱已注册，您将收到密码重置邮件",
				ExpiresIn: 3600,
			}, nil
		}
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "检查用户失败", err)
	}

	// 3. 生成密码重置Token
	tokenManager := NewPasswordResetTokenManager()
	resetToken, err := tokenManager.GenerateToken(ctx, req.Email)
	if err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "生成重置Token失败", err)
	}

	// 4. 构建重置邮件内容
	resetLink := fmt.Sprintf("https://qingyu.example.com/reset-password?token=%s&email=%s", resetToken, req.Email)
	emailBody := fmt.Sprintf(`
		<h2>密码重置请求</h2>
		<p>您好，%s，</p>
		<p>我们收到了您的密码重置请求。请点击下面的链接重置您的密码：</p>
		<p><a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #007bff; color: white; text-decoration: none; border-radius: 4px;">重置密码</a></p>
		<p>或者复制以下链接到浏览器：</p>
		<p><code>%s</code></p>
		<p>该链接将在1小时后过期。</p>
		<p>如果您没有请求重置密码，请忽略此邮件，您的密码不会被更改。</p>
		<p>青羽写作团队</p>
	`, user.Username, resetLink, resetLink)

	// 5. 发送重置邮件（当前为模拟发送）
	// TODO(Production): 集成真实的邮件发送服务
	// if s.emailService != nil {
	// 	err := s.emailService.SendEmail(ctx, &messaging.EmailRequest{
	// 		To:      []string{req.Email},
	// 		Subject: "青羽写作 - 密码重置",
	// 		Body:    emailBody,
	// 		IsHTML:  true,
	// 	})
	// 	if err != nil {
	// 		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "发送重置邮件失败", err)
	// 	}
	// }

	// 模拟：打印日志代替发送邮件
	fmt.Printf("[Password Reset] Token generated for %s: %s\n", req.Email, resetToken)
	fmt.Printf("[Password Reset] Reset link: %s\n", resetLink)
	fmt.Printf("[Password Reset] Email content:\n%s\n", emailBody)

	return &user2.RequestPasswordResetResponse{
		Success:   true,
		Message:   "密码重置邮件已发送",
		ExpiresIn: 3600, // 1小时 = 3600秒
	}, nil
}

// ConfirmPasswordReset 确认密码重置（使用Token设置新密码）
func (s *UserServiceImpl) ConfirmPasswordReset(ctx context.Context, req *user2.ConfirmPasswordResetRequest) (*user2.ConfirmPasswordResetResponse, error) {
	// 1. 验证请求数据
	if req.Email == "" || req.Token == "" || req.Password == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "邮箱、Token和新密码不能为空", nil)
	}

	// 2. 检查用户是否存在
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if repoInterfaces.IsNotFoundError(err) {
			return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
		}
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户失败", err)
	}

	// 3. 验证Token
	tokenManager := GetGlobalPasswordResetTokenManager()
	if err := tokenManager.ValidateToken(ctx, req.Email, req.Token); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "Token验证失败: "+err.Error(), nil)
	}

	// 4. 设置新密码
	if err := user.SetPassword(req.Password); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "密码哈希失败", err)
	}
	hashedPassword := user.Password

	// 5. 更新密码
	if err := s.userRepo.UpdatePassword(ctx, user.ID.Hex(), hashedPassword); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "更新密码失败", err)
	}

	// 6. 标记Token为已使用
	if err := tokenManager.MarkTokenAsUsed(ctx, req.Email); err != nil {
		// 记录警告但不影响流程
		zap.L().Warn("标记重置Token为已使用失败",
			zap.String("email", req.Email),
			zap.Error(err),
		)
	}

	return &user2.ConfirmPasswordResetResponse{
		Success: true,
		Message: "密码重置成功，请使用新密码登录",
	}, nil
}

// ==================== 新增方法：邮箱/手机/设备管理 ====================

// EmailExists 检查邮箱是否存在
func (s *UserServiceImpl) EmailExists(ctx context.Context, email string) (bool, error) {
	return s.userRepo.ExistsByEmail(ctx, email)
}

// UnbindEmail 解绑邮箱
func (s *UserServiceImpl) UnbindEmail(ctx context.Context, userID string) error {
	// 1. 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if repoInterfaces.IsNotFoundError(err) {
			return serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
		}
		return serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户失败", err)
	}

	// 2. 检查用户是否有邮箱
	if user.Email == "" {
		return serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户未绑定邮箱", nil)
	}

	// 3. 清除邮箱
	updates := map[string]interface{}{
		"email":      "",
		"verified":   false,
		"updated_at": time.Now(),
	}

	if err := s.userRepo.Update(ctx, userID, updates); err != nil {
		return serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "解绑邮箱失败", err)
	}

	return nil
}

// UnbindPhone 解绑手机
func (s *UserServiceImpl) UnbindPhone(ctx context.Context, userID string) error {
	// 1. 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if repoInterfaces.IsNotFoundError(err) {
			return serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
		}
		return serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户失败", err)
	}

	// 2. 检查用户是否有手机号
	if user.Phone == "" {
		return serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户未绑定手机", nil)
	}

	// 3. 清除手机号
	updates := map[string]interface{}{
		"phone":         "",
		"phone_verified": false,
		"updated_at":    time.Now(),
	}

	if err := s.userRepo.Update(ctx, userID, updates); err != nil {
		return serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "解绑手机失败", err)
	}

	return nil
}

// DeleteDevice 删除设备
func (s *UserServiceImpl) DeleteDevice(ctx context.Context, userID string, deviceID string) error {
	// 当前 User 模型不支持设备管理，这是一个占位实现
	// 实际实现需要：
	// 1. 检查用户是否存在
	// 2. 检查设备是否存在
	// 3. 删除设备
	return serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "设备管理功能尚未实现", nil)
}

// VerifyPassword 验证密码
func (s *UserServiceImpl) VerifyPassword(ctx context.Context, userID string, password string) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
	}

	// 验证密码
	if !user.ValidatePassword(password) {
		return serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "密码错误", nil)
	}

	return nil
}

// DowngradeRole 角色降级
// 将用户角色降级到指定角色（只能降级到reader或author）
func (s *UserServiceImpl) DowngradeRole(ctx context.Context, req *user2.DowngradeRoleRequest) (*user2.DowngradeRoleResponse, error) {
	// 1. 验证请求数据
	if req.UserID == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "用户ID不能为空", nil)
	}
	if req.TargetRole == "" {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "目标角色不能为空", nil)
	}
	if !req.Confirm {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "请确认降级操作", nil)
	}

	// 2. 验证目标角色是否合法
	validRoles := map[string]bool{
		"reader": true,
		"author": true,
	}
	if !validRoles[req.TargetRole] {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeValidation, "目标角色无效，只能降级为reader或author", nil)
	}

	// 3. 获取用户信息
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		if repoInterfaces.IsNotFoundError(err) {
			return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeNotFound, "用户不存在", err)
		}
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "获取用户失败", err)
	}

	// 4. 检查用户当前角色（只能从author或admin降级）
	currentRoles := user.Roles
	hasAuthorOrAdmin := false
	for _, role := range currentRoles {
		if role == "author" || role == "admin" {
			hasAuthorOrAdmin = true
			break
		}
	}

	if !hasAuthorOrAdmin {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeBusiness, "已经是读者，无法降级", nil)
	}

	// 5. 检查是否已经是目标角色
	for _, role := range currentRoles {
		if role == req.TargetRole {
			return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeBusiness, "当前已是目标角色", nil)
		}
	}

	// 6. 执行降级操作
	// 如果目标角色是reader，则清空所有角色，只保留reader
	// 如果目标角色是author，则移除admin角色（如果有）
	var newRoles []string
	if req.TargetRole == "reader" {
		newRoles = []string{"reader"}
	} else if req.TargetRole == "author" {
		// 保留author角色，移除admin
		for _, role := range currentRoles {
			if role == "author" || role == "reader" {
				newRoles = append(newRoles, role)
			}
		}
		// 确保至少有author角色
		if len(newRoles) == 0 {
			newRoles = []string{"author"}
		}
	}

	// 7. 更新用户角色
	updates := map[string]interface{}{
		"roles": newRoles,
	}
	if err := s.userRepo.Update(ctx, req.UserID, updates); err != nil {
		return nil, serviceInterfaces.NewServiceError(s.name, serviceInterfaces.ErrorTypeInternal, "更新用户角色失败", err)
	}

	return &user2.DowngradeRoleResponse{
		CurrentRoles: newRoles,
	}, nil
}
