package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	usersModel "Qingyu_backend/models/users"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"
)

var (
	ErrInvalidCode         = fmt.Errorf("验证码无效或已过期")
	ErrOldPasswordMismatch = fmt.Errorf("旧密码错误")
)

// PasswordService 密码服务
type PasswordService struct {
	verificationService *VerificationService
	userRepo            repoInterfaces.UserRepository
}

// NewPasswordService 创建密码服务
func NewPasswordService(
	verificationService *VerificationService,
	userRepo repoInterfaces.UserRepository,
) *PasswordService {
	return &PasswordService{
		verificationService: verificationService,
		userRepo:            userRepo,
	}
}

// SendResetCode 发送密码重置验证码
func (s *PasswordService) SendResetCode(ctx context.Context, email string) error {
	// 检查邮箱是否存在
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return fmt.Errorf("邮箱不存在")
	}

	// 发送重置验证码
	return s.verificationService.SendEmailCode(ctx, email, "reset_password")
}

// ResetPassword 重置密码
func (s *PasswordService) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	// 验证验证码
	if err := s.verificationService.VerifyCode(ctx, email, code, "reset_password"); err != nil {
		return ErrInvalidCode
	}

	// ✅ 添加：标记验证码为已使用（防止重复使用）
	if err := s.verificationService.MarkCodeAsUsed(ctx, email); err != nil {
		return fmt.Errorf("标记验证码失败: %w", err)
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	return s.userRepo.UpdatePasswordByEmail(ctx, email, string(hashedPassword))
}

// UpdatePassword 修改密码（需要旧密码）
func (s *PasswordService) UpdatePassword(ctx context.Context, userID string, oldPassword, newPassword string) error {
	// 验证userID格式
	if !primitive.IsValidObjectID(userID) {
		return fmt.Errorf("无效的用户ID")
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}

	// 验证旧密码
	if !user.ValidatePassword(oldPassword) {
		return ErrOldPasswordMismatch
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	return s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}

// checkPassword 验证密码
func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GetUserByEmail 根据邮箱获取用户（辅助方法）
func (s *PasswordService) GetUserByEmail(ctx context.Context, email string) (*usersModel.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// GetUserByID 根据ID获取用户（辅助方法）
func (s *PasswordService) GetUserByID(ctx context.Context, userID string) (*usersModel.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
