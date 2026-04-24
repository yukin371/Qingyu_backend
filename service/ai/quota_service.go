package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/models/ai"
	"Qingyu_backend/pkg/cache"
	"Qingyu_backend/pkg/quota"
	aiRepo "Qingyu_backend/repository/interfaces/ai"
	"Qingyu_backend/service/base"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuotaService 配额服务（增强版，支持Redis缓存和预警）
type QuotaService struct {
	quotaRepo       aiRepo.QuotaRepository
	redisClient     cache.RedisClient   // Redis客户端用于缓存
	eventBus        base.EventBus       // 事件总线用于预警通知
	cacheTTL        time.Duration       // 缓存过期时间（默认5分钟）
	policyService   *QuotaPolicyService // 策略服务（可选）
	profileResolver *QuotaProfileResolver

	// 预警阈值配置
	warningThreshold  float64 // 预警阈值（默认20%）
	criticalThreshold float64 // 严重阈值（默认10%）
}

// NewQuotaService 创建配额服务（基础版，无缓存）
func NewQuotaService(quotaRepo aiRepo.QuotaRepository) *QuotaService {
	return &QuotaService{
		quotaRepo:         quotaRepo,
		cacheTTL:          5 * time.Minute,
		warningThreshold:  0.2, // 20%
		criticalThreshold: 0.1, // 10%
	}
}

// NewQuotaServiceWithCache 创建配额服务（增强版，支持Redis缓存）
func NewQuotaServiceWithCache(
	quotaRepo aiRepo.QuotaRepository,
	redisClient cache.RedisClient,
	eventBus base.EventBus,
) *QuotaService {
	return &QuotaService{
		quotaRepo:         quotaRepo,
		redisClient:       redisClient,
		eventBus:          eventBus,
		cacheTTL:          5 * time.Minute,
		warningThreshold:  0.2, // 20%
		criticalThreshold: 0.1, // 10%
	}
}

// InitializeUserQuota 初始化用户配额
func (s *QuotaService) InitializeUserQuota(ctx context.Context, userID, userRole, membershipLevel string) error {
	profile := s.resolveQuotaProfile(ctx, userID, userRole, membershipLevel)

	// 检查是否已存在日配额
	existing, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, ai.QuotaTypeDaily)
	if err == nil && existing != nil {
		return nil // 已存在，不需要初始化
	}

	// 如果策略服务已加载，从策略获取配额
	defaultQuota := s.resolveEffectiveDailyQuota(ctx, profile.UserRole, profile.MembershipLevel)

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
			UserRole:          profile.UserRole,
			MembershipLevel:   profile.MembershipLevel,
			TotalConsumptions: 0,
			AveragePerDay:     0,
			CustomFields: map[string]interface{}{
				quotaCustomFieldProfileManaged: true,
			},
		},
	}

	return s.quotaRepo.CreateQuota(ctx, quota)
}

// Check 检查配额是否可用（实现 quota.Checker 接口）
func (s *QuotaService) Check(ctx context.Context, userID string, amount int) error {
	err := s.CheckQuota(ctx, userID, amount)
	if err == nil {
		return nil
	}
	// 统一对外错误域，避免中间件与业务层错误类型耦合。
	switch {
	case errors.Is(err, ai.ErrQuotaExhausted):
		return quota.ErrQuotaExhausted
	case errors.Is(err, ai.ErrQuotaSuspended):
		return quota.ErrQuotaSuspended
	case errors.Is(err, ai.ErrInsufficientQuota):
		return quota.ErrInsufficientQuota
	default:
		return err
	}
}

// CheckQuota 检查配额是否可用
func (s *QuotaService) CheckQuota(ctx context.Context, userID string, amount int) error {
	quota, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, ai.QuotaTypeDaily)
	if err != nil {
		if err == ai.ErrQuotaNotFound {
			// 配额不存在，尝试初始化
			if initErr := s.InitializeUserQuota(ctx, userID, "", ""); initErr != nil {
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

	// 自动升级配额：如果配置文件中的配额更大，自动升级
	if config.GlobalConfig != nil && config.GlobalConfig.AIQuota != nil {
		// 从配额元数据获取用户角色，如果没有则默认为reader/normal
		userRole := "reader"
		membershipLevel := "normal"
		if quota.Metadata != nil {
			if quota.Metadata.UserRole != "" {
				userRole = quota.Metadata.UserRole
			}
			if quota.Metadata.MembershipLevel != "" {
				membershipLevel = quota.Metadata.MembershipLevel
			}
		}

		configQuota := config.GlobalConfig.AIQuota.GetDefaultQuota(userRole, membershipLevel)
		if configQuota > quota.TotalQuota {
			// 配置中的配额更大，自动升级
			oldTotal := quota.TotalQuota
			quota.TotalQuota = configQuota
			// 增加剩余配额（按照增加的比例）
			increase := configQuota - oldTotal
			quota.RemainingQuota = quota.RemainingQuota + increase
			if quota.RemainingQuota < 0 {
				quota.RemainingQuota = 0
			}
			if quota.RemainingQuota > configQuota {
				quota.RemainingQuota = configQuota
			}
			// 更新到数据库
			if updateErr := s.quotaRepo.UpdateQuota(ctx, quota); updateErr != nil {
				// 升级失败不影响检查，继续使用旧值
				fmt.Printf("警告: 配额升级失败: %v\n", updateErr)
			}
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

// ConsumeQuota 消费配额（增强版：支持缓存失效和预警）
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

	// 清除缓存
	s.invalidateQuotaCache(ctx, userID, ai.QuotaTypeDaily)

	// 检查并触发预警
	s.checkAndPublishWarning(ctx, quota)

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

	// 清除缓存
	s.invalidateQuotaCache(ctx, userID, ai.QuotaTypeDaily)

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

// GetQuotaInfo 获取配额信息（支持Redis缓存）
func (s *QuotaService) GetQuotaInfo(ctx context.Context, userID string) (*ai.UserQuota, error) {
	// 尝试从缓存获取
	if s.redisClient != nil {
		cacheKey := fmt.Sprintf("quota:user:%s:daily", userID)
		cached, err := s.redisClient.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var quota ai.UserQuota
			if jsonErr := json.Unmarshal([]byte(cached), &quota); jsonErr == nil {
				return &quota, nil
			}
		}
	}

	// 缓存未命中，从数据库获取
	quota, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, ai.QuotaTypeDaily)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if s.redisClient != nil {
		s.cacheQuota(ctx, quota)
	}

	return quota, nil
}

// GetAllQuotas 获取用户所有配额
func (s *QuotaService) GetAllQuotas(ctx context.Context, userID string) ([]*ai.UserQuota, error) {
	return s.quotaRepo.GetAllQuotasByUserID(ctx, userID)
}

// GetQuotaStatistics 获取配额统计
func (s *QuotaService) GetQuotaStatistics(ctx context.Context, userID string) (*aiRepo.QuotaStatistics, error) {
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
			setQuotaManualOverride(newQuota, true)

			return s.quotaRepo.CreateQuota(ctx, newQuota)
		}
		return err
	}

	// 更新配额
	quota.TotalQuota = totalQuota
	quota.RemainingQuota = totalQuota - quota.UsedQuota
	setQuotaManualOverride(quota, true)

	if err := s.quotaRepo.UpdateQuota(ctx, quota); err != nil {
		return err
	}

	// 清除缓存
	s.invalidateQuotaCache(ctx, userID, quotaType)

	return nil
}

// SuspendUserQuota 暂停用户配额
func (s *QuotaService) SuspendUserQuota(ctx context.Context, userID string) error {
	quotas, err := s.quotaRepo.GetAllQuotasByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, quota := range quotas {
		setQuotaManualOverride(quota, true)
		quota.Status = ai.QuotaStatusSuspended
		if err := s.quotaRepo.UpdateQuota(ctx, quota); err != nil {
			return err
		}
		// 清除缓存
		s.invalidateQuotaCache(ctx, userID, quota.QuotaType)
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
		setQuotaManualOverride(quota, true)
		quota.Status = ai.QuotaStatusActive
		if err := s.quotaRepo.UpdateQuota(ctx, quota); err != nil {
			return err
		}
		// 清除缓存
		s.invalidateQuotaCache(ctx, userID, quota.QuotaType)
	}

	return nil
}

// RechargeQuota 配额充值（管理员或自助充值）
func (s *QuotaService) RechargeQuota(ctx context.Context, userID string, amount int, reason, operatorID string) error {
	if amount <= 0 {
		return fmt.Errorf("充值金额必须大于0")
	}

	quota, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, ai.QuotaTypeDaily)
	if err != nil {
		return fmt.Errorf("获取配额失败: %w", err)
	}

	beforeBalance := quota.RemainingQuota

	// 增加配额
	quota.TotalQuota += amount
	quota.RemainingQuota += amount
	setQuotaManualOverride(quota, true)

	// 更新数据库
	if err := s.quotaRepo.UpdateQuota(ctx, quota); err != nil {
		return fmt.Errorf("更新配额失败: %w", err)
	}

	// 清除缓存
	s.invalidateQuotaCache(ctx, userID, ai.QuotaTypeDaily)

	// 创建充值记录
	transaction := &ai.QuotaTransaction{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		QuotaType:     ai.QuotaTypeDaily,
		Amount:        amount,
		Type:          "recharge",
		Service:       "system",
		Description:   reason,
		BeforeBalance: beforeBalance,
		AfterBalance:  quota.RemainingQuota,
		Timestamp:     time.Now(),
	}

	return s.quotaRepo.CreateTransaction(ctx, transaction)
}

// ============ 私有辅助方法 ============

// cacheQuota 缓存配额信息
func (s *QuotaService) cacheQuota(ctx context.Context, quota *ai.UserQuota) {
	if s.redisClient == nil || quota == nil {
		return
	}

	cacheKey := fmt.Sprintf("quota:user:%s:%s", quota.UserID, quota.QuotaType)
	data, err := json.Marshal(quota)
	if err != nil {
		// 缓存失败不影响业务，仅记录日志
		fmt.Printf("缓存配额失败: %v\n", err)
		return
	}

	if err := s.redisClient.Set(ctx, cacheKey, string(data), s.cacheTTL); err != nil {
		fmt.Printf("写入Redis失败: %v\n", err)
	}
}

// invalidateQuotaCache 清除配额缓存
func (s *QuotaService) invalidateQuotaCache(ctx context.Context, userID string, quotaType ai.QuotaType) {
	if s.redisClient == nil {
		return
	}

	cacheKey := fmt.Sprintf("quota:user:%s:%s", userID, quotaType)
	if err := s.redisClient.Delete(ctx, cacheKey); err != nil {
		fmt.Printf("清除缓存失败: %v\n", err)
	}
}

// checkAndPublishWarning 检查并发布配额预警
func (s *QuotaService) checkAndPublishWarning(ctx context.Context, quota *ai.UserQuota) {
	if s.eventBus == nil || quota.TotalQuota == 0 {
		return
	}

	// 计算剩余百分比
	remainingPercent := float64(quota.RemainingQuota) / float64(quota.TotalQuota)

	var level string
	var shouldAlert bool

	if remainingPercent <= s.criticalThreshold {
		level = "critical"
		shouldAlert = true
	} else if remainingPercent <= s.warningThreshold {
		level = "warning"
		shouldAlert = true
	}

	if shouldAlert {
		// 发布预警事件
		event := &QuotaWarningEvent{
			UserID:           quota.UserID,
			QuotaType:        string(quota.QuotaType),
			TotalQuota:       quota.TotalQuota,
			RemainingQuota:   quota.RemainingQuota,
			UsedQuota:        quota.UsedQuota,
			RemainingPercent: remainingPercent * 100,
			Level:            level,
			Timestamp:        time.Now(),
		}

		// 异步发布事件
		if err := s.eventBus.PublishAsync(ctx, event); err != nil {
			fmt.Printf("发布配额预警事件失败: %v\n", err)
		}
	}
}

// QuotaWarningEvent 配额预警事件
type QuotaWarningEvent struct {
	UserID           string
	QuotaType        string
	TotalQuota       int
	RemainingQuota   int
	UsedQuota        int
	RemainingPercent float64
	Level            string // warning, critical
	Timestamp        time.Time
}

// GetEventType 实现Event接口
func (e *QuotaWarningEvent) GetEventType() string {
	return "quota.warning"
}

// GetEventData 实现Event接口
func (e *QuotaWarningEvent) GetEventData() interface{} {
	return e
}

// GetTimestamp 实现Event接口
func (e *QuotaWarningEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetSource 实现Event接口
func (e *QuotaWarningEvent) GetSource() string {
	return "QuotaService"
}

// SetPolicyService 设置策略服务
func (s *QuotaService) SetPolicyService(ps *QuotaPolicyService) {
	s.policyService = ps
}

// SetProfileResolver 设置用户配额画像解析器。
func (s *QuotaService) SetProfileResolver(resolver *QuotaProfileResolver) {
	s.profileResolver = resolver
}

// RefreshUserQuotaProfile 刷新用户配额画像并同步默认配额。
func (s *QuotaService) RefreshUserQuotaProfile(ctx context.Context, userID string) error {
	profile := s.resolveQuotaProfile(ctx, userID, "", "")
	quotas, err := s.quotaRepo.GetAllQuotasByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if len(quotas) == 0 {
		return s.InitializeUserQuota(ctx, userID, profile.UserRole, profile.MembershipLevel)
	}

	resolvedDailyQuota := s.resolveEffectiveDailyQuota(ctx, profile.UserRole, profile.MembershipLevel)

	for _, quota := range quotas {
		metadata := ensureQuotaMetadata(quota)
		metadata.UserRole = profile.UserRole
		metadata.MembershipLevel = profile.MembershipLevel
		metadata.CustomFields[quotaCustomFieldProfileManaged] = true

		if quota.QuotaType == ai.QuotaTypeDaily && !isQuotaManualOverride(quota) {
			quota.TotalQuota = resolvedDailyQuota
			quota.RemainingQuota = resolvedDailyQuota - quota.UsedQuota
			if quota.RemainingQuota < 0 {
				quota.RemainingQuota = 0
			}
		}

		previousStatus := quota.Status
		quota.BeforeUpdate()
		if previousStatus == ai.QuotaStatusSuspended {
			quota.Status = previousStatus
		}

		if err := s.quotaRepo.UpdateQuota(ctx, quota); err != nil {
			return err
		}
		s.invalidateQuotaCache(ctx, userID, quota.QuotaType)
	}

	return nil
}

func (s *QuotaService) resolveQuotaProfile(ctx context.Context, userID, fallbackRole, fallbackMembershipLevel string) *QuotaProfile {
	profile := &QuotaProfile{
		UserID:          userID,
		UserRole:        normalizeQuotaUserRole([]string{fallbackRole}),
		MembershipLevel: normalizeQuotaMembershipLevel(fallbackMembershipLevel),
	}

	if fallbackRole == "" {
		profile.UserRole = string(ai.UserRoleReader)
	}
	if fallbackMembershipLevel == "" {
		profile.MembershipLevel = "normal"
	}

	if s.profileResolver == nil {
		return profile
	}

	resolvedProfile, err := s.profileResolver.Resolve(ctx, userID)
	if err != nil || resolvedProfile == nil {
		return profile
	}

	if resolvedProfile.UserRole != "" {
		profile.UserRole = resolvedProfile.UserRole
	}
	if resolvedProfile.MembershipLevel != "" {
		profile.MembershipLevel = resolvedProfile.MembershipLevel
	}

	return profile
}

func (s *QuotaService) resolveEffectiveDailyQuota(ctx context.Context, userRole, membershipLevel string) int {
	if s.policyService != nil {
		policyQuota, err := s.policyService.GetEffectiveDailyQuota(ctx, userRole, membershipLevel)
		if err == nil && policyQuota != 0 {
			return policyQuota
		}
	}

	if config.GlobalConfig != nil && config.GlobalConfig.AIQuota != nil {
		return config.GlobalConfig.AIQuota.GetDefaultQuota(userRole, membershipLevel)
	}

	return ai.GetDefaultQuota(userRole, membershipLevel)
}
