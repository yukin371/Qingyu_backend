package user

import (
	"context"
	"testing"
)

// MockPasswordUserRepository 模拟用户仓库（用于密码服务测试）
type MockPasswordUserRepository struct {
	users map[string]*MockPasswordUser
}

type MockPasswordUser struct {
	ID             string
	Email          string
	HashedPassword string
}

func (m *MockPasswordUserRepository) GetByID(ctx context.Context, id string) (*MockPasswordUser, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, context.Canceled
	}
	return user, nil
}

func (m *MockPasswordUserRepository) GetByEmail(ctx context.Context, email string) (*MockPasswordUser, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, context.Canceled
}

func (m *MockPasswordUserRepository) UpdatePasswordByEmail(ctx context.Context, email string, hashedPassword string) error {
	for _, user := range m.users {
		if user.Email == email {
			user.HashedPassword = hashedPassword
			return nil
		}
	}
	return context.Canceled
}

func (m *MockPasswordUserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	user, exists := m.users[id]
	if !exists {
		return context.Canceled
	}
	user.HashedPassword = hashedPassword
	return nil
}

// TestNewPasswordService 测试创建密码服务
func TestNewPasswordService(t *testing.T) {
	verificationService := NewVerificationService(nil, nil)
	mockRepo := &MockPasswordUserRepository{users: make(map[string]*MockPasswordUser)}

	service := NewPasswordService(verificationService, mockRepo)

	if service == nil {
		t.Fatal("密码服务创建失败")
	}

	if service.verificationService == nil {
		t.Error("验证服务未正确设置")
	}

	if service.userRepo == nil {
		t.Error("用户仓库未正确设置")
	}
}

// TestPasswordService_SendResetCode 测试发送重置验证码
func TestPasswordService_SendResetCode(t *testing.T) {
	verificationService := NewVerificationService(nil, nil)
	mockRepo := &MockPasswordUserRepository{
		users: map[string]*MockPasswordUser{
			"user1": {
				ID:             "user1",
				Email:          "test@example.com",
				HashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			},
		},
	}

	service := NewPasswordService(verificationService, mockRepo)

	ctx := context.Background()
	err := service.SendResetCode(ctx, "test@example.com")
	if err != nil {
		t.Errorf("发送重置验证码失败: %v", err)
	}

	// 测试不存在的邮箱
	err = service.SendResetCode(ctx, "nonexistent@example.com")
	if err == nil {
		t.Error("应该返回邮箱不存在错误")
	}
}

// TestPasswordService_ResetPassword 测试重置密码
func TestPasswordService_ResetPassword(t *testing.T) {
	verificationService := NewVerificationService(nil, nil)
	mockRepo := &MockPasswordUserRepository{
		users: map[string]*MockPasswordUser{
			"user1": {
				ID:             "user1",
				Email:          "test@example.com",
				HashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			},
		},
	}

	service := NewPasswordService(verificationService, mockRepo)
	ctx := context.Background()

	// 首先生成一个验证码
	tokenManager := verificationService.GetVerificationTokenManager()
	code, err := tokenManager.GenerateCode(ctx, "user1", "test@example.com")
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}

	// 使用正确的验证码重置密码
	err = service.ResetPassword(ctx, "test@example.com", code, "NewPassword123")
	if err != nil {
		t.Errorf("重置密码失败: %v", err)
	}

	// 使用错误的验证码重置密码
	err = service.ResetPassword(ctx, "test@example.com", "000000", "NewPassword123")
	if err != ErrInvalidCode {
		t.Error("应该返回验证码无效错误")
	}
}

// TestPasswordService_UpdatePassword 测试修改密码
func TestPasswordService_UpdatePassword(t *testing.T) {
	verificationService := NewVerificationService(nil, nil)
	mockRepo := &MockPasswordUserRepository{
		users: map[string]*MockPasswordUser{
			"user1": {
				ID:             "user1",
				Email:          "test@example.com",
				HashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // "password"
			},
		},
	}

	service := NewPasswordService(verificationService, mockRepo)
	ctx := context.Background()

	// 使用正确的旧密码修改密码
	err := service.UpdatePassword(ctx, "user1", "password", "NewPassword123")
	if err != nil {
		t.Errorf("修改密码失败: %v", err)
	}

	// 使用错误的旧密码修改密码
	err = service.UpdatePassword(ctx, "user1", "WrongPassword", "NewPassword123")
	if err != ErrOldPasswordMismatch {
		t.Error("应该返回旧密码错误")
	}

	// 修改不存在的用户密码
	err = service.UpdatePassword(ctx, "nonexistent", "password", "NewPassword123")
	if err == nil {
		t.Error("应该返回用户不存在错误")
	}
}

// TestPasswordService_GetUserByEmail 测试根据邮箱获取用户
func TestPasswordService_GetUserByEmail(t *testing.T) {
	verificationService := NewVerificationService(nil, nil)
	mockRepo := &MockPasswordUserRepository{
		users: map[string]*MockPasswordUser{
			"user1": {
				ID:             "user1",
				Email:          "test@example.com",
				HashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			},
		},
	}

	service := NewPasswordService(verificationService, mockRepo)
	ctx := context.Background()

	// 获取存在的用户
	user, err := service.GetUserByEmail(ctx, "test@example.com")
	if err != nil {
		t.Errorf("获取用户失败: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Error("获取的用户邮箱不匹配")
	}

	// 获取不存在的用户
	_, err = service.GetUserByEmail(ctx, "nonexistent@example.com")
	if err == nil {
		t.Error("应该返回用户不存在错误")
	}
}
