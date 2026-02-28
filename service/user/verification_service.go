package user

import (
	"context"
	"fmt"

	"Qingyu_backend/config"
	usersModel "Qingyu_backend/models/users"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"
	"Qingyu_backend/service/channels"
	"go.uber.org/zap"
)

// VerificationService 验证服务
type VerificationService struct {
	emailService channels.EmailService // 邮件服务（可选，用于发送验证邮件）
	userRepo     repoInterfaces.UserRepository
	authRepo     sharedRepo.AuthRepository
	tokenManager *EmailVerificationTokenManager
}

// NewVerificationService 创建验证服务
func NewVerificationService(
	userRepo repoInterfaces.UserRepository,
	authRepo sharedRepo.AuthRepository,
	emailService channels.EmailService,
) *VerificationService {
	return &VerificationService{
		userRepo:     userRepo,
		authRepo:     authRepo,
		emailService: emailService,
		tokenManager: NewEmailVerificationTokenManager(),
	}
}

// SendEmailCode 发送邮箱验证码
func (s *VerificationService) SendEmailCode(ctx context.Context, email, purpose string) error {
	// 检查邮箱是否存在
	if purpose == "verify_email" {
		// 验证邮箱时，不需要检查邮箱是否存在
	} else if purpose == "reset_password" {
		// 重置密码时，需要检查邮箱是否存在
		user, err := s.userRepo.GetByEmail(ctx, email)
		if err != nil || user == nil {
			return fmt.Errorf("邮箱不存在")
		}
	}

	// 生成验证码
	// TODO: 从上下文中获取userID（如果已登录）
	var userID string
	// 如果是已登录用户发送验证码，从上下文获取userID

	code, err := s.tokenManager.GenerateCode(ctx, userID, email)
	if err != nil {
		return fmt.Errorf("生成验证码失败: %w", err)
	}

	// 发送邮件
	// TODO: 待EmailService实现后调用
	// if s.emailService != nil {
	//     err := s.emailService.SendEmail(ctx, &messaging.EmailRequest{
	//         To:      []string{email},
	//         Subject: "青羽写作 - 邮箱验证码",
	//         Body:    fmt.Sprintf("您的验证码是: %s", code),
	//         IsHTML:  false,
	//     })
	//     if err != nil {
	//         return fmt.Errorf("发送邮件失败: %w", err)
	//     }
	// }

	// 临时方案：仅在开发环境打印到控制台
	if config.GetEnvBool("APP_DEBUG", false) {
		zap.L().Debug("[VerificationService] 发送邮箱验证码",
			zap.Bool("has_email", email != ""),
			zap.Bool("has_code", code != ""),
			zap.Bool("has_purpose", purpose != ""),
		)
	}

	return nil
}

// SendPhoneCode 发送手机验证码（模拟实现）
func (s *VerificationService) SendPhoneCode(ctx context.Context, phone, purpose string) error {
	// 生成验证码
	var userID string
	code, err := s.tokenManager.GenerateCode(ctx, userID, phone) // 复用邮箱验证码管理器
	if err != nil {
		return fmt.Errorf("生成验证码失败: %w", err)
	}

	// 模拟实现：仅在开发环境打印到控制台
	if config.GetEnvBool("APP_DEBUG", false) {
		zap.L().Debug("[VerificationService] 发送手机验证码（模拟）",
			zap.Bool("has_phone", phone != ""),
			zap.Bool("has_code", code != ""),
			zap.Bool("has_purpose", purpose != ""),
		)
	}

	return nil
}

// VerifyCode 验证验证码
func (s *VerificationService) VerifyCode(ctx context.Context, target, code, purpose string) error {
	// 验证验证码（purpose参数暂时未使用，保留以备将来扩展）
	err := s.tokenManager.ValidateCode(ctx, target, code, purpose)
	if err != nil {
		return fmt.Errorf("验证码无效或已过期: %w", err)
	}

	return nil
}

// MarkCodeAsUsed 标记验证码为已使用（防止重复使用）
func (s *VerificationService) MarkCodeAsUsed(ctx context.Context, target string) error {
	return s.tokenManager.MarkCodeAsUsed(ctx, target)
}

// SetEmailVerified 设置邮箱已验证
func (s *VerificationService) SetEmailVerified(ctx context.Context, userID string, email string) error {
	// 更新用户邮箱验证状态
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}

	// 检查邮箱是否匹配
	if user.Email != email {
		return fmt.Errorf("邮箱不匹配")
	}

	// 更新验证状态
	return s.userRepo.SetEmailVerified(ctx, userID, true)
}

// SetPhoneVerified 设置手机已验证
func (s *VerificationService) SetPhoneVerified(ctx context.Context, userID string, phone string) error {
	// 更新用户手机验证状态
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}

	// 检查手机是否匹配
	if user.Phone != phone {
		return fmt.Errorf("手机号不匹配")
	}

	// 更新验证状态
	return s.userRepo.SetPhoneVerified(ctx, userID, true)
}

// CheckPassword 检查密码是否正确
func (s *VerificationService) CheckPassword(ctx context.Context, userID string, password string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}

	// 验证密码
	if !user.ValidatePassword(password) {
		return fmt.Errorf("密码错误")
	}

	return nil
}

// GetVerificationTokenManager 获取验证码管理器（用于测试）
func (s *VerificationService) GetVerificationTokenManager() *EmailVerificationTokenManager {
	return s.tokenManager
}

// EmailExists 检查邮箱是否存在
func (s *VerificationService) EmailExists(ctx context.Context, email string) (bool, error) {
	return s.userRepo.ExistsByEmail(ctx, email)
}

// PhoneExists 检查手机是否存在
func (s *VerificationService) PhoneExists(ctx context.Context, phone string) (bool, error) {
	return s.userRepo.ExistsByPhone(ctx, phone)
}

// GetUserByEmail 根据邮箱获取用户
func (s *VerificationService) GetUserByEmail(ctx context.Context, email string) (*usersModel.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// GetUserIDFromContext 从上下文中获取用户ID
func getuserIDFromContext(ctx context.Context) string {
	// TODO: 从JWT或session中获取用户ID
	return ""
}
