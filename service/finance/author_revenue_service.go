package finance

import (
	financeModel "Qingyu_backend/models/finance"
	"Qingyu_backend/models/shared/types"
	"Qingyu_backend/repository/interfaces/finance"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

// AuthorRevenueServiceImpl 作者收入服务实现
type AuthorRevenueServiceImpl struct {
	revenueRepo finance.AuthorRevenueRepository
}

// NewAuthorRevenueService 创建作者收入服务
func NewAuthorRevenueService(revenueRepo finance.AuthorRevenueRepository) AuthorRevenueService {
	return &AuthorRevenueServiceImpl{
		revenueRepo: revenueRepo,
	}
}

// ============ 收入查询 ============

// GetEarnings 获取作者收入列表
func (s *AuthorRevenueServiceImpl) GetEarnings(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error) {
	earnings, total, err := s.revenueRepo.GetEarningsByAuthor(ctx, authorID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("获取收入列表失败: %w", err)
	}

	return earnings, total, nil
}

// GetBookEarnings 获取某本书的收入
func (s *AuthorRevenueServiceImpl) GetBookEarnings(ctx context.Context, authorID string, bookID string) ([]*financeModel.AuthorEarning, int64, error) {
	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, 0, fmt.Errorf("无效的书籍ID: %w", err)
	}

	earnings, _, err := s.revenueRepo.GetEarningsByBook(ctx, bookOID, 1, 100) // 默认返回前100条
	if err != nil {
		return nil, 0, fmt.Errorf("获取书籍收入失败: %w", err)
	}

	// 过滤出该作者的收入（因为GetEarningsByBook没有过滤作者）
	var authorEarnings []*financeModel.AuthorEarning
	for _, earning := range earnings {
		if earning.AuthorID == authorID {
			authorEarnings = append(authorEarnings, earning)
		}
	}

	return authorEarnings, int64(len(authorEarnings)), nil
}

// GetRevenueDetails 获取收入明细
func (s *AuthorRevenueServiceImpl) GetRevenueDetails(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.RevenueDetail, int64, error) {
	details, total, err := s.revenueRepo.GetRevenueDetails(ctx, authorID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("获取收入明细失败: %w", err)
	}

	return details, total, nil
}

// GetRevenueStatistics 获取收入统计
func (s *AuthorRevenueServiceImpl) GetRevenueStatistics(ctx context.Context, authorID string, period string) ([]*financeModel.RevenueStatistics, error) {
	statistics, err := s.revenueRepo.GetRevenueStatistics(ctx, authorID, period, 0)
	if err != nil {
		return nil, fmt.Errorf("获取收入统计失败: %w", err)
	}

	return statistics, nil
}

// ============ 收入记录 ============

// CreateEarning 创建收入记录
func (s *AuthorRevenueServiceImpl) CreateEarning(ctx context.Context, earning *financeModel.AuthorEarning) error {
	err := s.revenueRepo.CreateEarning(ctx, earning)
	if err != nil {
		return fmt.Errorf("创建收入记录失败: %w", err)
	}

	// TODO: 更新收入明细和统计数据

	return nil
}

// CalculateEarning 计算收入分成
// 返回: (作者收入, 平台收入, error)
func (s *AuthorRevenueServiceImpl) CalculateEarning(ctx context.Context, earningType string, amount float64, authorID string, bookID primitive.ObjectID) (float64, float64, error) {
	var authorRate float64
	var platformRate float64

	switch earningType {
	case financeModel.EarningTypeChapterPurchase:
		authorRate = financeModel.ChapterPurchaseAuthorRate
		platformRate = financeModel.ChapterPurchasePlatformRate
	case financeModel.EarningTypeReward:
		authorRate = financeModel.RewardAuthorRate
		platformRate = financeModel.RewardPlatformRate
	case financeModel.EarningTypeVIPReading:
		authorRate = financeModel.VIPReadingAuthorRate
		platformRate = financeModel.VIPReadingPlatformRate
	default:
		return 0, 0, fmt.Errorf("未知的收入类型: %s", earningType)
	}

	authorIncome := amount * authorRate
	platformIncome := amount * platformRate

	return authorIncome, platformIncome, nil
}

// ============ 提现管理 ============

// CreateWithdrawalRequest 创建提现申请
func (s *AuthorRevenueServiceImpl) CreateWithdrawalRequest(ctx context.Context, userID string, amount float64, method string, account financeModel.WithdrawAccount) (*financeModel.WithdrawalRequest, error) {
	// 1. 验证提现金额
	if amount <= 0 {
		return nil, fmt.Errorf("提现金额必须大于0")
	}

	// 2. 计算手续费（示例：1%，最低1元）
	fee := amount * 0.01
	if fee < 1 {
		fee = 1
	}
	actualAmount := amount - fee

	// 3. TODO: 检查用户钱包余额是否足够
	// 这里需要集成钱包服务

	// 4. 创建提现申请
	request := &financeModel.WithdrawalRequest{
		UserID:       userID,
		Amount:       types.NewMoneyFromYuan(amount),
		Fee:          types.NewMoneyFromYuan(fee),
		ActualAmount: types.NewMoneyFromYuan(actualAmount),
		Method:       method,
		AccountInfo:  account,
		Status:       financeModel.WithdrawStatusPending,
	}

	err := s.revenueRepo.CreateWithdrawalRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("创建提现申请失败: %w", err)
	}

	// 5. TODO: 冻结钱包余额
	// 这里需要集成钱包服务

	return request, nil
}

// GetWithdrawals 获取提现记录
func (s *AuthorRevenueServiceImpl) GetWithdrawals(ctx context.Context, userID string, page, pageSize int) ([]*financeModel.WithdrawalRequest, int64, error) {
	requests, total, err := s.revenueRepo.GetUserWithdrawalRequests(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("获取提现记录失败: %w", err)
	}

	return requests, total, nil
}

// ============ 结算管理 ============

// GetSettlements 获取结算记录
func (s *AuthorRevenueServiceImpl) GetSettlements(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.Settlement, int64, error) {
	settlements, total, err := s.revenueRepo.GetAuthorSettlements(ctx, authorID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("获取结算记录失败: %w", err)
	}

	return settlements, total, nil
}

// GetSettlement 获取结算详情
func (s *AuthorRevenueServiceImpl) GetSettlement(ctx context.Context, settlementID string) (*financeModel.Settlement, error) {
	oid, err := primitive.ObjectIDFromHex(settlementID)
	if err != nil {
		return nil, fmt.Errorf("无效的结算ID: %w", err)
	}

	settlement, err := s.revenueRepo.GetSettlement(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("获取结算详情失败: %w", err)
	}

	return settlement, nil
}

// ============ 税务信息 ============

// GetTaxInfo 获取税务信息
func (s *AuthorRevenueServiceImpl) GetTaxInfo(ctx context.Context, userID string) (*financeModel.TaxInfo, error) {
	taxInfo, err := s.revenueRepo.GetTaxInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取税务信息失败: %w", err)
	}

	return taxInfo, nil
}

// UpdateTaxInfo 更新税务信息
func (s *AuthorRevenueServiceImpl) UpdateTaxInfo(ctx context.Context, userID string, taxInfo *financeModel.TaxInfo) error {
	// 检查是否已存在
	_, err := s.revenueRepo.GetTaxInfo(ctx, userID)

	if err != nil {
		// 不存在，创建新的
		taxInfo.UserID = userID
		err := s.revenueRepo.CreateTaxInfo(ctx, taxInfo)
		if err != nil {
			return fmt.Errorf("创建税务信息失败: %w", err)
		}
		return nil
	}

	// 存在，更新
	updates := map[string]interface{}{
		"id_type":   taxInfo.IDType,
		"id_number": taxInfo.IDNumber,
		"name":      taxInfo.Name,
		"tax_type":  taxInfo.TaxType,
		"tax_rate":  taxInfo.TaxRate,
	}

	err = s.revenueRepo.UpdateTaxInfo(ctx, userID, updates)
	if err != nil {
		return fmt.Errorf("更新税务信息失败: %w", err)
	}

	return nil
}
