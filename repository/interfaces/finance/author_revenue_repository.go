package finance

import (
	financeModel "Qingyu_backend/models/finance"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthorRevenueRepository 作者收入仓储接口
type AuthorRevenueRepository interface {
	// 收入记录管理
	CreateEarning(ctx context.Context, earning *financeModel.AuthorEarning) error
	GetEarning(ctx context.Context, earningID primitive.ObjectID) (*financeModel.AuthorEarning, error)
	ListEarnings(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error)
	UpdateEarning(ctx context.Context, earningID primitive.ObjectID, updates map[string]interface{}) error
	BatchUpdateEarnings(ctx context.Context, earningIDs []primitive.ObjectID, updates map[string]interface{}) error
	GetEarningsByAuthor(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error)
	GetEarningsByBook(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error)

	// 提现申请管理
	CreateWithdrawalRequest(ctx context.Context, request *financeModel.WithdrawalRequest) error
	GetWithdrawalRequest(ctx context.Context, requestID primitive.ObjectID) (*financeModel.WithdrawalRequest, error)
	ListWithdrawalRequests(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.WithdrawalRequest, int64, error)
	UpdateWithdrawalRequest(ctx context.Context, requestID primitive.ObjectID, updates map[string]interface{}) error
	GetUserWithdrawalRequests(ctx context.Context, userID string, page, pageSize int) ([]*financeModel.WithdrawalRequest, int64, error)

	// 结算管理
	CreateSettlement(ctx context.Context, settlement *financeModel.Settlement) error
	GetSettlement(ctx context.Context, settlementID primitive.ObjectID) (*financeModel.Settlement, error)
	ListSettlements(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.Settlement, int64, error)
	UpdateSettlement(ctx context.Context, settlementID primitive.ObjectID, updates map[string]interface{}) error
	GetAuthorSettlements(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.Settlement, int64, error)
	GetPendingSettlements(ctx context.Context) ([]*financeModel.Settlement, error)

	// 收入统计
	GetRevenueStatistics(ctx context.Context, authorID string, period string, limit int) ([]*financeModel.RevenueStatistics, error)
	CreateRevenueStatistics(ctx context.Context, statistics *financeModel.RevenueStatistics) error
	UpdateRevenueStatistics(ctx context.Context, statisticsID primitive.ObjectID, updates map[string]interface{}) error

	// 收入明细
	GetRevenueDetails(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.RevenueDetail, int64, error)
	GetRevenueDetailByBook(ctx context.Context, authorID string, bookID primitive.ObjectID) (*financeModel.RevenueDetail, error)
	CreateRevenueDetail(ctx context.Context, detail *financeModel.RevenueDetail) error
	UpdateRevenueDetail(ctx context.Context, detailID primitive.ObjectID, updates map[string]interface{}) error

	// 税务信息
	CreateTaxInfo(ctx context.Context, taxInfo *financeModel.TaxInfo) error
	GetTaxInfo(ctx context.Context, userID string) (*financeModel.TaxInfo, error)
	UpdateTaxInfo(ctx context.Context, userID string, updates map[string]interface{}) error
}
