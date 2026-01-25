package user

import (
	"context"
	"testing"
	"time"
)

// TestEmailVerificationTokenManager_GenerateCode 测试验证码生成
func TestEmailVerificationTokenManager_GenerateCode(t *testing.T) {
	manager := NewEmailVerificationTokenManager()
	ctx := context.Background()

	// 测试生成验证码
	code, err := manager.GenerateCode(ctx, "user123", "test@example.com")
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}

	// 验证码应该是6位数字
	if len(code) != 6 {
		t.Errorf("验证码长度应该是6位，实际为: %d", len(code))
	}

	// 验证码应该全是数字
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("验证码应该只包含数字，发现: %c", c)
		}
	}
}

// TestEmailVerificationTokenManager_ValidateCode 测试验证码验证
func TestEmailVerificationTokenManager_ValidateCode(t *testing.T) {
	manager := NewEmailVerificationTokenManager()
	ctx := context.Background()

	// 生成验证码
	userID := "user123"
	email := "test@example.com"
	code, err := manager.GenerateCode(ctx, userID, email)
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}

	// 验证正确的验证码
	err = manager.ValidateCode(ctx, userID, email, code)
	if err != nil {
		t.Errorf("验证正确验证码失败: %v", err)
	}

	// 验证错误的验证码
	err = manager.ValidateCode(ctx, userID, email, "000000")
	if err == nil {
		t.Error("应该返回验证码错误")
	}

	// 验证错误用户ID的验证码
	err = manager.ValidateCode(ctx, "wrong_user", email, code)
	if err == nil {
		t.Error("应该返回用户不匹配错误")
	}

	// 验证错误邮箱的验证码
	err = manager.ValidateCode(ctx, userID, "wrong@example.com", code)
	if err == nil {
		t.Error("应该返回验证码不存在错误")
	}
}

// TestEmailVerificationTokenManager_MarkCodeAsUsed 测试标记验证码已使用
func TestEmailVerificationTokenManager_MarkCodeAsUsed(t *testing.T) {
	manager := NewEmailVerificationTokenManager()
	ctx := context.Background()

	// 生成验证码
	userID := "user123"
	email := "test@example.com"
	code, err := manager.GenerateCode(ctx, userID, email)
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}

	// 验证验证码应该成功
	err = manager.ValidateCode(ctx, userID, email, code)
	if err != nil {
		t.Errorf("验证正确验证码失败: %v", err)
	}

	// 标记验证码已使用
	err = manager.MarkCodeAsUsed(ctx, email)
	if err != nil {
		t.Fatalf("标记验证码已使用失败: %v", err)
	}

	// 验证已使用的验证码应该失败
	err = manager.ValidateCode(ctx, userID, email, code)
	if err == nil {
		t.Error("已使用的验证码应该验证失败")
	}
}

// TestEmailVerificationTokenManager_CleanExpiredCodes 测试清理过期验证码
func TestEmailVerificationTokenManager_CleanExpiredCodes(t *testing.T) {
	manager := NewEmailVerificationTokenManager()
	ctx := context.Background()

	// 生成验证码
	userID := "user123"
	_, err := manager.GenerateCode(ctx, userID, "test1@example.com")
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}

	// 手动添加一个过期的验证码
	manager.tokens["expired@example.com"] = &VerificationTokenInfo{
		Code:      "123456",
		Email:     "expired@example.com",
		UserID:    userID,
		ExpiresAt: time.Now().Add(-time.Hour), // 1小时前过期
		Used:      false,
	}

	// 记录清理前的数量
	beforeCount := len(manager.tokens)

	// 清理过期验证码
	manager.CleanExpiredCodes(ctx)

	// 记录清理后的数量
	afterCount := len(manager.tokens)

	// 应该清理了至少1个过期验证码
	if afterCount >= beforeCount {
		t.Error("应该清理至少1个过期验证码")
	}

	// 验证过期的验证码已被清理
	_, exists := manager.tokens["expired@example.com"]
	if exists {
		t.Error("过期验证码应该被清理")
	}
}

// TestEmailVerificationTokenManager_ConcurrentAccess 测试并发访问
func TestEmailVerificationTokenManager_ConcurrentAccess(t *testing.T) {
	manager := NewEmailVerificationTokenManager()
	ctx := context.Background()

	// 并发生成验证码
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(index int) {
			email := "user" + string(rune('0'+index)) + "@example.com"
			userID := "user" + string(rune('0'+index))
			_, err := manager.GenerateCode(ctx, userID, email)
			if err != nil {
				t.Errorf("并发生成验证码失败: %v", err)
			}
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}
