package interfaces

import (
	"context"
	financeModel "Qingyu_backend/models/finance"

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

// AuthorRevenueService 作者收入服务接口
type AuthorRevenueService interface {
	// 收入查询
	GetEarnings(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error)
	GetBookEarnings(ctx context.Context, authorID string, bookID string) ([]*financeModel.AuthorEarning, int64, error)
	GetRevenueDetails(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.RevenueDetail, int64, error)
	GetRevenueStatistics(ctx context.Context, authorID string, period string) ([]*financeModel.RevenueStatistics, error)

	// 收入记录
	CreateEarning(ctx context.Context, earning *financeModel.AuthorEarning) error
	CalculateEarning(ctx context.Context, earningType string, amount float64, authorID string, bookID primitive.ObjectID) (float64, float64, error)

	// 提现管理
	CreateWithdrawalRequest(ctx context.Context, userID string, amount float64, method string, account financeModel.WithdrawAccount) (*financeModel.WithdrawalRequest, error)
	GetWithdrawals(ctx context.Context, userID string, page, pageSize int) ([]*financeModel.WithdrawalRequest, int64, error)

	// 结算管理
	GetSettlements(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.Settlement, int64, error)
	GetSettlement(ctx context.Context, settlementID string) (*financeModel.Settlement, error)

	// 税务信息
	GetTaxInfo(ctx context.Context, userID string) (*financeModel.TaxInfo, error)
	UpdateTaxInfo(ctx context.Context, userID string, taxInfo *financeModel.TaxInfo) error
}
