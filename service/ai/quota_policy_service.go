package ai

import (
	"context"
	"errors"
	"fmt"

	aiModels "Qingyu_backend/models/ai"
	aiInterfaces "Qingyu_backend/repository/interfaces/ai"

	"Qingyu_backend/config"
)

// QuotaPolicyService 配额策略服务
type QuotaPolicyService struct {
	policyRepo aiInterfaces.QuotaPolicyRepository
}

// NewQuotaPolicyService 创建配额策略服务
func NewQuotaPolicyService(policyRepo aiInterfaces.QuotaPolicyRepository) *QuotaPolicyService {
	return &QuotaPolicyService{
		policyRepo: policyRepo,
	}
}

// CreatePolicy 创建配额策略
func (s *QuotaPolicyService) CreatePolicy(ctx context.Context, policy *aiModels.QuotaPolicy) error {
	if policy == nil {
		return errors.New("策略不能为空")
	}

	// 检查是否已存在相同角色和等级的策略
	existing, err := s.policyRepo.GetByRoleAndLevel(ctx, policy.UserRole, policy.MembershipLevel)
	if err == nil && existing != nil {
		return fmt.Errorf("已存在角色[%s]等级[%s]的策略", policy.UserRole, policy.MembershipLevel)
	}

	// 设置默认值
	if policy.DailyQuota <= 0 {
		policy.DailyQuota = 1000 // 默认日配额
	}
	if policy.MonthlyQuota <= 0 {
		policy.MonthlyQuota = 30000 // 默认月配额
	}
	if policy.TotalQuota < 0 {
		policy.TotalQuota = -1 // 无限配额
	}

	return s.policyRepo.Create(ctx, policy)
}

// GetPolicy 根据ID获取策略
func (s *QuotaPolicyService) GetPolicy(ctx context.Context, id string) (*aiModels.QuotaPolicy, error) {
	return s.policyRepo.GetByID(ctx, id)
}

// GetPolicyByRole 根据角色和等级获取策略
func (s *QuotaPolicyService) GetPolicyByRole(ctx context.Context, role aiModels.UserRole, level aiModels.MembershipLevel) (*aiModels.QuotaPolicy, error) {
	return s.policyRepo.GetByRoleAndLevel(ctx, role, level)
}

// ListPolicies 获取策略列表
func (s *QuotaPolicyService) ListPolicies(ctx context.Context, role, status string, page, limit int) ([]*aiModels.QuotaPolicy, int64, error) {
	return s.policyRepo.List(ctx, role, status, page, limit)
}

// UpdatePolicy 更新策略
func (s *QuotaPolicyService) UpdatePolicy(ctx context.Context, policy *aiModels.QuotaPolicy) error {
	if policy == nil {
		return errors.New("策略不能为空")
	}

	// 检查策略是否存在
	existing, err := s.policyRepo.GetByID(ctx, policy.ID.Hex())
	if err != nil {
		return fmt.Errorf("策略不存在: %w", err)
	}

	// 如果角色或等级发生了变化，检查是否与其他策略冲突
	if existing.UserRole != policy.UserRole || existing.MembershipLevel != policy.MembershipLevel {
		_, err = s.policyRepo.GetByRoleAndLevel(ctx, policy.UserRole, policy.MembershipLevel)
		if err == nil {
			return fmt.Errorf("已存在角色[%s]等级[%s]的策略", policy.UserRole, policy.MembershipLevel)
		}
	}

	// 验证配额值
	if policy.DailyQuota <= 0 {
		policy.DailyQuota = 1000 // 默认日配额
	}
	if policy.MonthlyQuota <= 0 {
		policy.MonthlyQuota = 30000 // 默认月配额
	}
	if policy.TotalQuota < 0 && policy.TotalQuota != -1 {
		return fmt.Errorf("总配额不能为负数（-1表示无限）")
	}

	return s.policyRepo.Update(ctx, policy)
}

// DeletePolicy 删除策略
func (s *QuotaPolicyService) DeletePolicy(ctx context.Context, id string) error {
	// 检查策略是否存在
	policy, err := s.policyRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("策略不存在: %w", err)
	}

	// 不能删除默认策略
	if policy.IsDefault {
		return fmt.Errorf("不能删除默认策略")
	}

	return s.policyRepo.Delete(ctx, id)
}

// InitializeDefaultPolicies 初始化默认策略
func (s *QuotaPolicyService) InitializeDefaultPolicies(ctx context.Context) error {
	if config.GlobalConfig == nil || config.GlobalConfig.AIQuota == nil {
		return errors.New("配置未加载，无法获取默认配额配置")
	}

	// 获取所有角色和等级组合
	roles := []aiModels.UserRole{
		aiModels.UserRoleReader,
		aiModels.UserRoleWriter,
		aiModels.UserRoleAdmin,
	}

	levels := []aiModels.MembershipLevel{
		aiModels.MembershipLevelNormal,
		aiModels.MembershipLevelVipMonthly,
		aiModels.MembershipLevelVipYearly,
		aiModels.MembershipLevelSuperVip,
	}

	// 检查是否已存在策略
	existingPolicies, _, err := s.policyRepo.List(ctx, "", "", 1, 1000)
	if err != nil {
		return fmt.Errorf("获取现有策略失败: %w", err)
	}

	// 如果已有策略，直接返回
	if len(existingPolicies) > 0 {
		return nil
	}

	// 创建默认策略
	for _, role := range roles {
		for _, level := range levels {
			dailyQuota := config.GlobalConfig.AIQuota.GetDefaultQuota(string(role), string(level))

			policy := &aiModels.QuotaPolicy{
				UserRole:        role,
				MembershipLevel: level,
				DailyQuota:      dailyQuota,
				MonthlyQuota:    dailyQuota * 30, // 月配额为日配额的30倍
				TotalQuota:      -1,              // 无限总配额
				IsDefault:       true,
				Status:          aiModels.QuotaPolicyStatusActive,
				Description:     fmt.Sprintf("默认策略 - 角色%s-%s", role, level),
			}

			if err := s.policyRepo.Create(ctx, policy); err != nil {
				return fmt.Errorf("创建策略失败[%s-%s]: %w", role, level, err)
			}
		}
	}

	return nil
}

// GetEffectiveDailyQuota 获取角色的有效日配额
func (s *QuotaPolicyService) GetEffectiveDailyQuota(ctx context.Context, userRole, membershipLevel string) (int, error) {
	// 转换类型
	role, err := s.parseUserRole(userRole)
	if err != nil {
		return 0, fmt.Errorf("无效的用户角色: %w", err)
	}

	level, err := s.parseMembershipLevel(membershipLevel)
	if err != nil {
		return 0, fmt.Errorf("无效的会员等级: %w", err)
	}

	// 尝试从策略获取
	policy, err := s.policyRepo.GetByRoleAndLevel(ctx, role, level)
	if err == nil && policy != nil && policy.Status == aiModels.QuotaPolicyStatusActive {
		return policy.DailyQuota, nil
	}

	// 策略不存在或未激活，使用默认配置
	if config.GlobalConfig != nil && config.GlobalConfig.AIQuota != nil {
		return config.GlobalConfig.AIQuota.GetDefaultQuota(userRole, membershipLevel), nil
	}

	// 回退到模型默认值
	return aiModels.GetDefaultQuota(userRole, membershipLevel), nil
}

// 私有辅助方法

// parseUserRole 解析用户角色
func (s *QuotaPolicyService) parseUserRole(role string) (aiModels.UserRole, error) {
	switch role {
	case "reader":
		return aiModels.UserRoleReader, nil
	case "writer":
		return aiModels.UserRoleWriter, nil
	case "admin":
		return aiModels.UserRoleAdmin, nil
	default:
		return "", fmt.Errorf("不支持的角色: %s", role)
	}
}

// parseMembershipLevel 解析会员等级
func (s *QuotaPolicyService) parseMembershipLevel(level string) (aiModels.MembershipLevel, error) {
	switch level {
	case "normal":
		return aiModels.MembershipLevelNormal, nil
	case "vip_monthly":
		return aiModels.MembershipLevelVipMonthly, nil
	case "vip_yearly":
		return aiModels.MembershipLevelVipYearly, nil
	case "super_vip":
		return aiModels.MembershipLevelSuperVip, nil
	default:
		return "", fmt.Errorf("不支持的会员等级: %s", level)
	}
}