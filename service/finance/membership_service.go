package finance

import (
	financeModel "Qingyu_backend/models/finance"
	"Qingyu_backend/repository/interfaces/finance"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MembershipService 会员服务接口
type MembershipService interface {
	// 套餐管理
	GetPlans(ctx context.Context) ([]*financeModel.MembershipPlan, error)
	GetPlan(ctx context.Context, planID string) (*financeModel.MembershipPlan, error)

	// 订阅管理
	Subscribe(ctx context.Context, userID string, planID string, paymentMethod string) (*financeModel.UserMembership, error)
	GetMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error)
	CancelMembership(ctx context.Context, userID string) error
	RenewMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error)

	// 会员权益
	GetBenefits(ctx context.Context, level string) ([]*financeModel.MembershipBenefit, error)
	GetUsage(ctx context.Context, userID string) ([]*financeModel.MembershipUsage, error)

	// 会员卡管理
	ActivateCard(ctx context.Context, userID string, cardCode string) (*financeModel.UserMembership, error)
	ListCards(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.MembershipCard, int64, error)

	// 会员检查
	CheckMembership(ctx context.Context, userID string, level string) (bool, error)
	IsVIP(ctx context.Context, userID string) (bool, error)
}

// MembershipServiceImpl 会员服务实现
type MembershipServiceImpl struct {
	membershipRepo finance.MembershipRepository
}

// NewMembershipService 创建会员服务
func NewMembershipService(membershipRepo finance.MembershipRepository) MembershipService {
	return &MembershipServiceImpl{
		membershipRepo: membershipRepo,
	}
}

// ============ 套餐管理 ============

// GetPlans 获取套餐列表
func (s *MembershipServiceImpl) GetPlans(ctx context.Context) ([]*financeModel.MembershipPlan, error) {
	plans, err := s.membershipRepo.ListPlans(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("获取套餐列表失败: %w", err)
	}

	return plans, nil
}

// GetPlan 获取套餐详情
func (s *MembershipServiceImpl) GetPlan(ctx context.Context, planID string) (*financeModel.MembershipPlan, error) {
	oid, err := primitive.ObjectIDFromHex(planID)
	if err != nil {
		return nil, fmt.Errorf("无效的套餐ID: %w", err)
	}

	plan, err := s.membershipRepo.GetPlan(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("获取套餐失败: %w", err)
	}

	return plan, nil
}

// ============ 订阅管理 ============

// Subscribe 订阅会员
func (s *MembershipServiceImpl) Subscribe(ctx context.Context, userID string, planID string, paymentMethod string) (*financeModel.UserMembership, error) {
	// 1. 验证套餐ID
	planOID, err := primitive.ObjectIDFromHex(planID)
	if err != nil {
		return nil, fmt.Errorf("无效的套餐ID: %w", err)
	}

	// 2. 获取套餐信息
	plan, err := s.membershipRepo.GetPlan(ctx, planOID)
	if err != nil {
		return nil, fmt.Errorf("获取套餐失败: %w", err)
	}

	if !plan.IsEnabled {
		return nil, fmt.Errorf("套餐已下架")
	}

	// 3. 检查是否已有会员
	existingMembership, err := s.membershipRepo.GetMembership(ctx, userID)
	if err == nil && existingMembership != nil {
		// 如果已有会员，检查是否可以续费
		if existingMembership.Status == financeModel.MembershipStatusActive {
			return nil, fmt.Errorf("您已是会员，无需重复订阅")
		}
	}

	// 4. 计算会员时间
	now := time.Now()
	var startTime time.Time
	var endTime time.Time

	if existingMembership != nil && existingMembership.EndTime.After(now) {
		// 续费：从当前到期时间开始
		startTime = existingMembership.EndTime
	} else {
		// 新订阅：从现在开始
		startTime = now
	}

	endTime = startTime.AddDate(0, 0, plan.Duration)

	// 5. 创建会员记录
	membership := &financeModel.UserMembership{
		UserID:      userID,
		PlanID:      plan.ID,
		PlanName:    plan.Name,
		PlanType:    plan.Type,
		Level:       s.getLevelFromType(plan.Type),
		StartTime:   startTime,
		EndTime:     endTime,
		AutoRenew:   false,
		Status:      financeModel.MembershipStatusActive,
		PaymentID:   primitive.NewObjectID(), // TODO: 创建支付记录
		ActivatedAt: now,
	}

	err = s.membershipRepo.CreateMembership(ctx, membership)
	if err != nil {
		return nil, fmt.Errorf("创建会员失败: %w", err)
	}

	// 6. TODO: 创建支付记录并扣款
	// 这里需要集成钱包服务或支付网关

	return membership, nil
}

// GetMembership 获取会员状态
func (s *MembershipServiceImpl) GetMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error) {
	membership, err := s.membershipRepo.GetMembership(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取会员状态失败: %w", err)
	}

	// 检查会员是否过期
	if membership.Status == financeModel.MembershipStatusActive && time.Now().After(membership.EndTime) {
		// 自动更新为过期状态
		updates := map[string]interface{}{
			"status": financeModel.MembershipStatusExpired,
		}
		_ = s.membershipRepo.UpdateMembership(ctx, membership.ID, updates)
		membership.Status = financeModel.MembershipStatusExpired
	}

	return membership, nil
}

// CancelMembership 取消自动续费
func (s *MembershipServiceImpl) CancelMembership(ctx context.Context, userID string) error {
	membership, err := s.membershipRepo.GetMembership(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取会员信息失败: %w", err)
	}

	if membership.Status != financeModel.MembershipStatusActive {
		return fmt.Errorf("会员未激活或已过期")
	}

	// 取消自动续费
	updates := map[string]interface{}{
		"auto_renew": false,
	}

	err = s.membershipRepo.UpdateMembership(ctx, membership.ID, updates)
	if err != nil {
		return fmt.Errorf("取消自动续费失败: %w", err)
	}

	return nil
}

// RenewMembership 手动续费
func (s *MembershipServiceImpl) RenewMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error) {
	membership, err := s.membershipRepo.GetMembership(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取会员信息失败: %w", err)
	}

	// 获取原套餐
	plan, err := s.membershipRepo.GetPlan(ctx, membership.PlanID)
	if err != nil {
		return nil, fmt.Errorf("获取套餐信息失败: %w", err)
	}

	// 计算新的到期时间
	var newEndTime time.Time
	if membership.EndTime.After(time.Now()) {
		newEndTime = membership.EndTime.AddDate(0, 0, plan.Duration)
	} else {
		newEndTime = time.Now().AddDate(0, 0, plan.Duration)
	}

	// 更新会员信息
	updates := map[string]interface{}{
		"end_time": newEndTime,
		"status":   financeModel.MembershipStatusActive,
	}

	err = s.membershipRepo.UpdateMembership(ctx, membership.ID, updates)
	if err != nil {
		return nil, fmt.Errorf("续费失败: %w", err)
	}

	// TODO: 创建支付记录并扣款

	// 重新获取更新后的会员信息
	return s.membershipRepo.GetMembership(ctx, userID)
}

// ============ 会员权益 ============

// GetBenefits 获取权益列表
func (s *MembershipServiceImpl) GetBenefits(ctx context.Context, level string) ([]*financeModel.MembershipBenefit, error) {
	benefits, err := s.membershipRepo.ListBenefits(ctx, level, true)
	if err != nil {
		return nil, fmt.Errorf("获取权益列表失败: %w", err)
	}

	return benefits, nil
}

// GetUsage 获取权益使用情况
func (s *MembershipServiceImpl) GetUsage(ctx context.Context, userID string) ([]*financeModel.MembershipUsage, error) {
	usages, err := s.membershipRepo.ListUsages(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取权益使用情况失败: %w", err)
	}

	return usages, nil
}

// ============ 会员卡管理 ============

// ActivateCard 激活会员卡
func (s *MembershipServiceImpl) ActivateCard(ctx context.Context, userID string, cardCode string) (*financeModel.UserMembership, error) {
	// 1. 查找会员卡
	card, err := s.membershipRepo.GetMembershipCardByCode(ctx, cardCode)
	if err != nil {
		return nil, fmt.Errorf("会员卡不存在或已失效: %w", err)
	}

	// 2. 检查会员卡状态
	if card.Status != financeModel.CardStatusUnused {
		return nil, fmt.Errorf("会员卡已被使用")
	}

	// 检查是否过期
	if card.ExpireAt != nil && time.Now().After(*card.ExpireAt) {
		return nil, fmt.Errorf("会员卡已过期")
	}

	// 3. 检查是否已有会员
	existingMembership, err := s.membershipRepo.GetMembership(ctx, userID)

	// 4. 计算会员时间
	now := time.Now()
	var startTime time.Time
	var endTime time.Time

	if err == nil && existingMembership != nil && existingMembership.EndTime.After(now) {
		// 已有会员且未过期：从到期时间开始
		startTime = existingMembership.EndTime
	} else {
		// 新会员或已过期：从现在开始
		startTime = now
	}

	endTime = startTime.AddDate(0, 0, card.Duration)

	// 5. 创建或更新会员
	var membership *financeModel.UserMembership
	if err == nil && existingMembership != nil {
		// 更新现有会员
		updates := map[string]interface{}{
			"plan_id":    card.PlanID,
			"plan_name":  card.PlanType,
			"plan_type":  card.PlanType,
			"level":      s.getLevelFromType(card.PlanType),
			"start_time": startTime,
			"end_time":   endTime,
			"status":     financeModel.MembershipStatusActive,
		}

		err = s.membershipRepo.UpdateMembership(ctx, existingMembership.ID, updates)
		if err != nil {
			return nil, fmt.Errorf("激活会员卡失败: %w", err)
		}

		membership, _ = s.membershipRepo.GetMembershipByID(ctx, existingMembership.ID)
	} else {
		// 创建新会员
		membership = &financeModel.UserMembership{
			UserID:      userID,
			PlanID:      card.PlanID,
			PlanName:    card.PlanType,
			PlanType:    card.PlanType,
			Level:       s.getLevelFromType(card.PlanType),
			StartTime:   startTime,
			EndTime:     endTime,
			AutoRenew:   false,
			Status:      financeModel.MembershipStatusActive,
			ActivatedAt: now,
		}

		err = s.membershipRepo.CreateMembership(ctx, membership)
		if err != nil {
			return nil, fmt.Errorf("创建会员失败: %w", err)
		}
	}

	// 6. 更新会员卡状态
	nowTime := now
	cardUpdates := map[string]interface{}{
		"status":       financeModel.CardStatusUsed,
		"activated_by": userID,
		"activated_at": nowTime,
	}

	err = s.membershipRepo.UpdateMembershipCard(ctx, card.ID, cardUpdates)
	if err != nil {
		return nil, fmt.Errorf("更新会员卡状态失败: %w", err)
	}

	return membership, nil
}

// ListCards 列出会员卡
func (s *MembershipServiceImpl) ListCards(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.MembershipCard, int64, error) {
	cards, err := s.membershipRepo.ListMembershipCards(ctx, filter, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("查询会员卡列表失败: %w", err)
	}

	total, err := s.membershipRepo.CountMembershipCards(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计会员卡失败: %w", err)
	}

	return cards, total, nil
}

// ============ 会员检查 ============

// CheckMembership 检查会员等级
func (s *MembershipServiceImpl) CheckMembership(ctx context.Context, userID string, level string) (bool, error) {
	membership, err := s.membershipRepo.GetMembership(ctx, userID)
	if err != nil {
		return false, nil // 没有会员不算错误
	}

	// 检查是否激活且未过期
	if membership.Status != financeModel.MembershipStatusActive {
		return false, nil
	}

	if time.Now().After(membership.EndTime) {
		return false, nil
	}

	// 检查等级
	if level != "" && membership.Level != level {
		return false, nil
	}

	return true, nil
}

// IsVIP 检查是否是VIP
func (s *MembershipServiceImpl) IsVIP(ctx context.Context, userID string) (bool, error) {
	membership, err := s.membershipRepo.GetMembership(ctx, userID)
	if err != nil {
		return false, nil
	}

	if membership.Status != financeModel.MembershipStatusActive {
		return false, nil
	}

	if time.Now().After(membership.EndTime) {
		return false, nil
	}

	// 检查是否是VIP级别
	return membership.Level == financeModel.MembershipLevelVIPMonthly ||
		membership.Level == financeModel.MembershipLevelVIPYearly ||
		membership.Level == financeModel.MembershipLevelSuperVIP, nil
}

// ============ 辅助函数 ============

// getLevelFromType 从套餐类型获取会员等级
func (s *MembershipServiceImpl) getLevelFromType(planType string) string {
	switch planType {
	case financeModel.MembershipTypeMonthly:
		return financeModel.MembershipLevelVIPMonthly
	case financeModel.MembershipTypeYearly:
		return financeModel.MembershipLevelVIPYearly
	case financeModel.MembershipTypeSuper:
		return financeModel.MembershipLevelSuperVIP
	default:
		return financeModel.MembershipLevelNormal
	}
}
