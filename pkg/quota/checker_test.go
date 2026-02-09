package quota

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCheckerInterface 测试配额检查器接口契约
// 这个测试定义了配额检查器应该有的行为，而不依赖具体实现
func TestCheckerInterface(t *testing.T) {
	t.Run("应该能够检查配额并返回nil表示通过", func(t *testing.T) {
		checker := &mockChecker{shouldPass: true}
		err := checker.Check(context.Background(), "user123", 1000)
		assert.NoError(t, err)
	})

	t.Run("配额不足时应该返回错误", func(t *testing.T) {
		checker := &mockChecker{shouldPass: false}
		err := checker.Check(context.Background(), "user123", 1000)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额")
	})

	t.Run("应该传递用户ID和预估数量", func(t *testing.T) {
		checker := &mockChecker{}
		_ = checker.Check(context.Background(), "user456", 2000)

		assert.Equal(t, "user456", checker.lastUserID)
		assert.Equal(t, 2000, checker.lastAmount)
	})
}

// mockChecker 是一个用于测试的模拟配额检查器
type mockChecker struct {
	shouldPass bool
	lastUserID string
	lastAmount int
}

func (m *mockChecker) Check(ctx context.Context, userID string, amount int) error {
	m.lastUserID = userID
	m.lastAmount = amount
	if !m.shouldPass {
		return errors.New("配额不足")
	}
	return nil
}

// TestMockCheckerImplementsInterface 确保mock实现了接口
func TestMockCheckerImplementsInterface(t *testing.T) {
	var checker Checker = &mockChecker{}
	assert.NotNil(t, checker)
}
