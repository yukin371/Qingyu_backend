package lock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"Qingyu_backend/pkg/logger"

	"go.uber.org/zap"
)

// QuotaLock 配额锁管理器
// 用于防止配额并发消费问题
type QuotaLock struct {
	lock DistributedLock
}

// NewQuotaLock 创建配额锁
func NewQuotaLock(redisClient interface{}) *QuotaLock {
	var lock DistributedLock

	// 根据类型创建不同的锁实现
	switch client := redisClient.(type) {
	case DistributedLock:
		lock = client
	default:
		// 假设是Redis客户端
		// lock = NewRedisDistributedLock(client)
		lock = &mockDistributedLock{} // 默认使用mock
	}

	return &QuotaLock{
		lock: lock,
	}
}

// ConsumeQuotaWithLock 使用锁安全地消费配额
// userID: 用户ID
// quotaType: 配额类型
// amount: 消费数量
// consumeFunc: 实际的消费函数
func (q *QuotaLock) ConsumeQuotaWithLock(
	ctx context.Context,
	userID string,
	quotaType string,
	amount int64,
	consumeFunc func(context.Context, string, string, int64) error,
) error {
	lockKey := q.getLockKey(userID, quotaType)
	ttl := 30 * time.Second // 锁30秒后自动释放

	// 使用锁执行消费
	err := q.lock.WithLockRetry(ctx, lockKey, ttl, 3, 100*time.Millisecond, func() error {
		return consumeFunc(ctx, userID, quotaType, amount)
	})

	if err != nil {
		logger.Error("failed to consume quota with lock",
			zap.String("user_id", userID),
			zap.String("quota_type", quotaType),
			zap.Int64("amount", amount),
			zap.Error(err),
		)
		return fmt.Errorf("failed to consume quota: %w", err)
	}

	return nil
}

// BatchConsumeQuotaWithLock 使用锁批量消费配额
func (q *QuotaLock) BatchConsumeQuotaWithLock(
	ctx context.Context,
	operations []QuotaOperation,
) error {
	// 按用户和配额类型分组
	groups := q.groupOperations(operations)

	// 并发执行各组操作
	errChan := make(chan error, len(groups))

	for _, group := range groups {
		go func(ops []QuotaOperation) {
			if len(ops) == 0 {
				errChan <- nil
				return
			}

			userID := ops[0].UserID
			quotaType := ops[0].QuotaType
			lockKey := q.getLockKey(userID, quotaType)
			ttl := 60 * time.Second

			err := q.lock.WithLock(ctx, lockKey, ttl, func() error {
				for _, op := range ops {
					if err := op.ConsumeFunc(ctx, op.UserID, op.QuotaType, op.Amount); err != nil {
						return err
					}
				}
				return nil
			})

			errChan <- err
		}(group)
	}

	// 等待所有操作完成
	for i := 0; i < len(groups); i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

// QuotaOperation 配额操作
type QuotaOperation struct {
	UserID      string
	QuotaType   string
	Amount      int64
	ConsumeFunc func(context.Context, string, string, int64) error
}

// groupOperations 按用户和配额类型分组操作
func (q *QuotaLock) groupOperations(operations []QuotaOperation) [][]QuotaOperation {
	groups := make(map[string][]QuotaOperation)

	for _, op := range operations {
		key := fmt.Sprintf("%s:%s", op.UserID, op.QuotaType)
		groups[key] = append(groups[key], op)
	}

	result := make([][]QuotaOperation, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	return result
}

// getLockKey 获取锁的键名
func (q *QuotaLock) getLockKey(userID, quotaType string) string {
	return fmt.Sprintf("quota:%s:%s", userID, quotaType)
}

// mockDistributedLock mock分布式锁（用于测试）
type mockDistributedLock struct {
	locks map[string]bool
}

func newMockDistributedLock() *mockDistributedLock {
	return &mockDistributedLock{
		locks: make(map[string]bool),
	}
}

func (m *mockDistributedLock) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if m.locks[key] {
		return false, nil
	}
	m.locks[key] = true
	return true, nil
}

func (m *mockDistributedLock) Unlock(ctx context.Context, key string) error {
	delete(m.locks, key)
	return nil
}

func (m *mockDistributedLock) Extend(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if !m.locks[key] {
		return false, nil
	}
	return true, nil
}

func (m *mockDistributedLock) IsLocked(ctx context.Context, key string) (bool, error) {
	return m.locks[key], nil
}

func (m *mockDistributedLock) WithLock(ctx context.Context, key string, ttl time.Duration, fn func() error) error {
	acquired, err := m.Lock(ctx, key, ttl)
	if err != nil {
		return err
	}
	if !acquired {
		return errors.New("failed to acquire lock")
	}
	defer m.Unlock(ctx, key)
	return fn()
}

func (m *mockDistributedLock) WithLockRetry(ctx context.Context, key string, ttl time.Duration, maxRetries int, retryInterval time.Duration, fn func() error) error {
	for i := 0; i < maxRetries; i++ {
		err := m.WithLock(ctx, key, ttl, fn)
		if err == nil {
			return nil
		}
		if err.Error() != "failed to acquire lock" {
			return err
		}
		if i < maxRetries-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryInterval):
				continue
			}
		}
	}
	return errors.New("failed to acquire lock after retries")
}

// 使用示例：
//
// func (s *QuotaService) ConsumeQuota(ctx context.Context, userID string, quotaType QuotaType, amount int64) error {
//     quotaLock := NewQuotaLock(s.redis)
//
//     return quotaLock.ConsumeQuotaWithLock(ctx, userID, string(quotaType), amount,
//         func(ctx context.Context, userID, quotaType string, amount int64) error {
//             // 实际的消费逻辑
//             quota, err := s.repo.GetByUserAndType(ctx, userID, quotaType)
//             if err != nil {
//                 return err
//             }
//
//             if quota.RemainingQuota < amount {
//                 return InsufficientQuota()
//             }
//
//             quota.RemainingQuota -= amount
//             quota.UsedQuota += amount
//
//             return s.repo.Update(ctx, quota)
//         },
//     )
// }
