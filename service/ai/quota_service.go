package ai

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/ai"
	"Qingyu_backend/repository/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuotaService 配额服务
type QuotaService struct {
	quotaRepo interfaces.QuotaRepository
}

// NewQuotaService 创建配额服务
func NewQuotaService(quotaRepo interfaces.QuotaRepository) *QuotaService {
	return &QuotaService{
		quotaRepo: quotaRepo,
	}
}

// InitializeUserQuota 初始化用户配额
func (s *QuotaService) InitializeUserQuota(ctx context.Context, userID, userRole, membershipLevel string) error {
	// 检查是否已存在日配额
	existing, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, ai.QuotaTypeDaily)
	if err == nil && existing != nil {
		return nil // 已存在，不需要初始化
	}

	// 获取默认配额
	defaultQuota := ai.GetDefaultQuota(userRole, membershipLevel)

	// 创建日配额
	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     defaultQuota,
		UsedQuota:      0,
		RemainingQuota: defaultQuota,
		Status:         ai.QuotaStatusActive,
		ResetAt:        time.Now().AddDate(0, 0, 1), // 明天重置
		Metadata: &ai.QuotaMetadata{
			UserRole:          userRole,
			MembershipLevel:   membershipLevel,
			TotalConsumptions: 0,
			AveragePerDay:     0,
		},
	}

	return s.quotaRepo.CreateQuota(ctx, quota)
}

// CheckQuota 检查配额是否可用
func (s *QuotaService) CheckQuota(ctx context.Context, userID string, amount int) error {
	quota, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, ai.QuotaTypeDaily)
	if err != nil {
		if err == ai.ErrQuotaNotFound {
			// 配额不存在，尝试初始化
			if initErr := s.InitializeUserQuota(ctx, userID, "reader", "normal"); initErr != nil {
				return fmt.Errorf("初始化配额失败: %w", initErr)
			}
			// 重新获取
			quota, err = s.quotaRepo.GetQuotaByUserID(ctx, userID, ai.QuotaTypeDaily)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// 检查是否可以消费
	if !quota.CanConsume(amount) {
		if quota.Status == ai.QuotaStatusExhausted {
			return ai.ErrQuotaExhausted
		}
		if quota.Status == ai.QuotaStatusSuspended {
			return ai.ErrQuotaSuspended
		}
		return ai.ErrInsufficientQuota
	}

	return nil
}

// ConsumeQuota 消费配额
func (s *QuotaService) ConsumeQuota(ctx context.Context, userID string, amount int, service, model, requestID string) error {
	quota, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, ai.QuotaTypeDaily)
	if err != nil {
		return err
	}

	// 记录消费前余额
	beforeBalance := quota.RemainingQuota

	// 消费配额
	if err := quota.Consume(amount); err != nil {
		return err
	}

	// 更新配额
	if err := s.quotaRepo.UpdateQuota(ctx, quota); err != nil {
		return fmt.Errorf("更新配额失败: %w", err)
	}

	// 创建事务记录
	transaction := &ai.QuotaTransaction{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		QuotaType:     ai.QuotaTypeDaily,
		Amount:        amount,
		Type:          "consume",
		Service:       service,
		Model:         model,
		RequestID:     requestID,
		Description:   fmt.Sprintf("消费%d配额用于%s服务", amount, service),
		BeforeBalance: beforeBalance,
		AfterBalance:  quota.RemainingQuota,
		Timestamp:     time.Now(),
	}

	return s.quotaRepo.CreateTransaction(ctx, transaction)
}

// RestoreQuota 恢复配额（用于错误回滚）
func (s *QuotaService) RestoreQuota(ctx context.Context, userID string, amount int, reason string) error {
	quota, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, ai.QuotaTypeDaily)
	if err != nil {
		return err
	}

	beforeBalance := quota.RemainingQuota
	quota.Restore(amount)

	if err := s.quotaRepo.UpdateQuota(ctx, quota); err != nil {
		return fmt.Errorf("更新配额失败: %w", err)
	}

	// 创建事务记录
	transaction := &ai.QuotaTransaction{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		QuotaType:     ai.QuotaTypeDaily,
		Amount:        -amount, // 负数表示恢复
		Type:          "restore",
		Service:       "system",
		Description:   reason,
		BeforeBalance: beforeBalance,
		AfterBalance:  quota.RemainingQuota,
		Timestamp:     time.Now(),
	}

	return s.quotaRepo.CreateTransaction(ctx, transaction)
}

// GetQuotaInfo 获取配额信息
func (s *QuotaService) GetQuotaInfo(ctx context.Context, userID string) (*ai.UserQuota, error) {
	return s.quotaRepo.GetQuotaByUserID(ctx, userID, ai.QuotaTypeDaily)
}

// GetAllQuotas 获取用户所有配额
func (s *QuotaService) GetAllQuotas(ctx context.Context, userID string) ([]*ai.UserQuota, error) {
	return s.quotaRepo.GetAllQuotasByUserID(ctx, userID)
}

// GetQuotaStatistics 获取配额统计
func (s *QuotaService) GetQuotaStatistics(ctx context.Context, userID string) (*interfaces.QuotaStatistics, error) {
	return s.quotaRepo.GetQuotaStatistics(ctx, userID)
}

// GetTransactionHistory 获取配额事务历史
func (s *QuotaService) GetTransactionHistory(ctx context.Context, userID string, limit, offset int) ([]*ai.QuotaTransaction, error) {
	return s.quotaRepo.GetTransactionsByUserID(ctx, userID, limit, offset)
}

// ResetDailyQuotas 重置日配额（定时任务调用）
func (s *QuotaService) ResetDailyQuotas(ctx context.Context) error {
	return s.quotaRepo.BatchResetQuotas(ctx, ai.QuotaTypeDaily)
}

// ResetMonthlyQuotas 重置月配额（定时任务调用）
func (s *QuotaService) ResetMonthlyQuotas(ctx context.Context) error {
	return s.quotaRepo.BatchResetQuotas(ctx, ai.QuotaTypeMonthly)
}

// UpdateUserQuota 更新用户配额（管理员操作）
func (s *QuotaService) UpdateUserQuota(ctx context.Context, userID string, quotaType ai.QuotaType, totalQuota int) error {
	quota, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, quotaType)
	if err != nil {
		if err == ai.ErrQuotaNotFound {
			// 创建新配额
			newQuota := &ai.UserQuota{
				UserID:     userID,
				QuotaType:  quotaType,
				TotalQuota: totalQuota,
				UsedQuota:  0,
				Status:     ai.QuotaStatusActive,
			}

			// 设置重置时间
			now := time.Now()
			switch quotaType {
			case ai.QuotaTypeDaily:
				newQuota.ResetAt = now.AddDate(0, 0, 1)
			case ai.QuotaTypeMonthly:
				newQuota.ResetAt = now.AddDate(0, 1, 0)
			}

			return s.quotaRepo.CreateQuota(ctx, newQuota)
		}
		return err
	}

	// 更新配额
	quota.TotalQuota = totalQuota
	quota.RemainingQuota = totalQuota - quota.UsedQuota
	return s.quotaRepo.UpdateQuota(ctx, quota)
}

// SuspendUserQuota 暂停用户配额
func (s *QuotaService) SuspendUserQuota(ctx context.Context, userID string) error {
	quotas, err := s.quotaRepo.GetAllQuotasByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, quota := range quotas {
		quota.Status = ai.QuotaStatusSuspended
		if err := s.quotaRepo.UpdateQuota(ctx, quota); err != nil {
			return err
		}
	}

	return nil
}

// ActivateUserQuota 激活用户配额
func (s *QuotaService) ActivateUserQuota(ctx context.Context, userID string) error {
	quotas, err := s.quotaRepo.GetAllQuotasByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, quota := range quotas {
		quota.Status = ai.QuotaStatusActive
		if err := s.quotaRepo.UpdateQuota(ctx, quota); err != nil {
			return err
		}
	}

	return nil
}
