package finance

import (
	financeModel "Qingyu_backend/models/finance"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MembershipRepository 会员仓储接口
type MembershipRepository interface {
	// 套餐管理
	CreatePlan(ctx context.Context, plan *financeModel.MembershipPlan) error
	GetPlan(ctx context.Context, planID primitive.ObjectID) (*financeModel.MembershipPlan, error)
	GetPlanByType(ctx context.Context, planType string) (*financeModel.MembershipPlan, error)
	ListPlans(ctx context.Context, enabledOnly bool) ([]*financeModel.MembershipPlan, error)
	UpdatePlan(ctx context.Context, planID primitive.ObjectID, updates map[string]interface{}) error
	DeletePlan(ctx context.Context, planID primitive.ObjectID) error

	// 用户会员管理
	CreateMembership(ctx context.Context, membership *financeModel.UserMembership) error
	GetMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error)
	GetMembershipByID(ctx context.Context, membershipID primitive.ObjectID) (*financeModel.UserMembership, error)
	UpdateMembership(ctx context.Context, membershipID primitive.ObjectID, updates map[string]interface{}) error
	DeleteMembership(ctx context.Context, membershipID primitive.ObjectID) error
	ListMemberships(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.UserMembership, error)

	// 会员卡管理
	CreateMembershipCard(ctx context.Context, card *financeModel.MembershipCard) error
	GetMembershipCard(ctx context.Context, cardID primitive.ObjectID) (*financeModel.MembershipCard, error)
	GetMembershipCardByCode(ctx context.Context, code string) (*financeModel.MembershipCard, error)
	UpdateMembershipCard(ctx context.Context, cardID primitive.ObjectID, updates map[string]interface{}) error
	ListMembershipCards(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.MembershipCard, error)
	CountMembershipCards(ctx context.Context, filter map[string]interface{}) (int64, error)
	BatchCreateMembershipCards(ctx context.Context, cards []*financeModel.MembershipCard) error

	// 会员权益管理
	CreateBenefit(ctx context.Context, benefit *financeModel.MembershipBenefit) error
	GetBenefit(ctx context.Context, benefitID primitive.ObjectID) (*financeModel.MembershipBenefit, error)
	GetBenefitByCode(ctx context.Context, code string) (*financeModel.MembershipBenefit, error)
	ListBenefits(ctx context.Context, level string, enabledOnly bool) ([]*financeModel.MembershipBenefit, error)
	UpdateBenefit(ctx context.Context, benefitID primitive.ObjectID, updates map[string]interface{}) error
	DeleteBenefit(ctx context.Context, benefitID primitive.ObjectID) error

	// 会员权益使用情况
	CreateUsage(ctx context.Context, usage *financeModel.MembershipUsage) error
	GetUsage(ctx context.Context, userID string, benefitCode string) (*financeModel.MembershipUsage, error)
	UpdateUsage(ctx context.Context, usageID primitive.ObjectID, updates map[string]interface{}) error
	ListUsages(ctx context.Context, userID string) ([]*financeModel.MembershipUsage, error)
}
