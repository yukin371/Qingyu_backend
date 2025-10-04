package wallet

import (
	"context"
	"fmt"
	"time"

	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// UnifiedWalletService 统一的钱包服务实现
// 实现 WalletService 接口，整合钱包、交易、提现三大功能
type UnifiedWalletService struct {
	walletRepo sharedRepo.WalletRepository

	// 内部组件服务
	walletMgr      *WalletServiceImpl
	transactionMgr *TransactionServiceImpl
	withdrawMgr    *WithdrawServiceImpl
}

// NewUnifiedWalletService 创建统一钱包服务
func NewUnifiedWalletService(walletRepo sharedRepo.WalletRepository) WalletService {
	return &UnifiedWalletService{
		walletRepo:     walletRepo,
		walletMgr:      &WalletServiceImpl{walletRepo: walletRepo},
		transactionMgr: &TransactionServiceImpl{walletRepo: walletRepo},
		withdrawMgr:    &WithdrawServiceImpl{walletRepo: walletRepo},
	}
}

// ============ 钱包管理 ============

// CreateWallet 创建钱包
func (s *UnifiedWalletService) CreateWallet(ctx context.Context, userID string) (*Wallet, error) {
	return s.walletMgr.CreateWallet(ctx, userID)
}

// GetWallet 获取钱包（根据用户ID）
func (s *UnifiedWalletService) GetWallet(ctx context.Context, userID string) (*Wallet, error) {
	return s.walletMgr.GetWallet(ctx, userID)
}

// GetBalance 获取余额（根据用户ID）
func (s *UnifiedWalletService) GetBalance(ctx context.Context, userID string) (float64, error) {
	return s.walletMgr.GetBalance(ctx, userID)
}

// FreezeWallet 冻结钱包（根据用户ID）
func (s *UnifiedWalletService) FreezeWallet(ctx context.Context, userID string) error {
	return s.walletMgr.FreezeWallet(ctx, userID, "管理员冻结")
}

// UnfreezeWallet 解冻钱包（根据用户ID）
func (s *UnifiedWalletService) UnfreezeWallet(ctx context.Context, userID string) error {
	return s.walletMgr.UnfreezeWallet(ctx, userID)
}

// ============ 交易操作 ============

// Recharge 充值（根据用户ID）
func (s *UnifiedWalletService) Recharge(ctx context.Context, userID string, amount float64, method string) (*Transaction, error) {
	// 1. 获取钱包
	wallet, err := s.walletMgr.GetWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取钱包失败: %w", err)
	}

	// 2. 执行充值
	orderNo := generateOrderNo() // 生成订单号
	return s.transactionMgr.Recharge(ctx, wallet.ID, amount, method, orderNo)
}

// Consume 消费（根据用户ID）
func (s *UnifiedWalletService) Consume(ctx context.Context, userID string, amount float64, reason string) (*Transaction, error) {
	// 1. 获取钱包
	wallet, err := s.walletMgr.GetWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取钱包失败: %w", err)
	}

	// 2. 执行消费
	return s.transactionMgr.Consume(ctx, wallet.ID, amount, reason)
}

// Transfer 转账（根据用户ID）
func (s *UnifiedWalletService) Transfer(ctx context.Context, fromUserID, toUserID string, amount float64, reason string) (*Transaction, error) {
	// 1. 获取源钱包
	fromWallet, err := s.walletMgr.GetWallet(ctx, fromUserID)
	if err != nil {
		return nil, fmt.Errorf("获取源钱包失败: %w", err)
	}

	// 2. 获取目标钱包
	toWallet, err := s.walletMgr.GetWallet(ctx, toUserID)
	if err != nil {
		return nil, fmt.Errorf("获取目标钱包失败: %w", err)
	}

	// 3. 执行转账
	if err := s.transactionMgr.Transfer(ctx, fromWallet.ID, toWallet.ID, amount, reason); err != nil {
		return nil, err
	}

	// 4. 返回转出交易记录（最新的一条transfer_out记录）
	// 注意：Transfer方法内部创建了两条记录，这里返回转出记录
	transactions, err := s.transactionMgr.ListTransactions(ctx, fromWallet.ID, 1, 0)
	if err != nil || len(transactions) == 0 {
		return nil, fmt.Errorf("获取交易记录失败")
	}

	return transactions[0], nil
}

// ============ 交易查询 ============

// GetTransaction 获取交易记录
func (s *UnifiedWalletService) GetTransaction(ctx context.Context, transactionID string) (*Transaction, error) {
	return s.transactionMgr.GetTransaction(ctx, transactionID)
}

// ListTransactions 列出交易记录（根据用户ID）
func (s *UnifiedWalletService) ListTransactions(ctx context.Context, userID string, req *ListTransactionsRequest) ([]*Transaction, error) {
	// 1. 获取钱包
	wallet, err := s.walletMgr.GetWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取钱包失败: %w", err)
	}

	// 2. 计算分页参数
	limit := req.PageSize
	if limit <= 0 {
		limit = 20 // 默认每页20条
	}
	offset := (req.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// 3. 查询交易列表
	return s.transactionMgr.ListTransactions(ctx, wallet.ID, limit, offset)
}

// ============ 提现管理 ============

// RequestWithdraw 申请提现（根据用户ID）
func (s *UnifiedWalletService) RequestWithdraw(ctx context.Context, userID string, amount float64, account string) (*WithdrawRequest, error) {
	// 1. 获取钱包
	wallet, err := s.walletMgr.GetWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取钱包失败: %w", err)
	}

	// 2. 创建提现申请（默认支付宝）
	method := "alipay"
	return s.withdrawMgr.CreateWithdrawRequest(ctx, userID, wallet.ID, amount, method, account)
}

// GetWithdrawRequest 获取提现申请
func (s *UnifiedWalletService) GetWithdrawRequest(ctx context.Context, withdrawID string) (*WithdrawRequest, error) {
	return s.withdrawMgr.GetWithdrawRequest(ctx, withdrawID)
}

// ListWithdrawRequests 列出提现申请
func (s *UnifiedWalletService) ListWithdrawRequests(ctx context.Context, req *ListWithdrawRequestsRequest) ([]*WithdrawRequest, error) {
	// 计算分页参数
	limit := req.PageSize
	if limit <= 0 {
		limit = 20 // 默认每页20条
	}
	offset := (req.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	return s.withdrawMgr.ListWithdrawRequests(ctx, req.UserID, req.Status, limit, offset)
}

// ApproveWithdraw 批准提现
func (s *UnifiedWalletService) ApproveWithdraw(ctx context.Context, withdrawID, adminID string) error {
	return s.withdrawMgr.ApproveWithdraw(ctx, withdrawID, adminID, "")
}

// RejectWithdraw 拒绝提现
func (s *UnifiedWalletService) RejectWithdraw(ctx context.Context, withdrawID, adminID, reason string) error {
	return s.withdrawMgr.RejectWithdraw(ctx, withdrawID, adminID, reason)
}

// ProcessWithdraw 处理提现（标记为已打款）
func (s *UnifiedWalletService) ProcessWithdraw(ctx context.Context, withdrawID string) error {
	// 获取提现申请
	request, err := s.withdrawMgr.GetWithdrawRequest(ctx, withdrawID)
	if err != nil {
		return err
	}

	// 检查状态
	if request.Status != "approved" {
		return fmt.Errorf("只能处理已批准的提现申请，当前状态: %s", request.Status)
	}

	// 更新状态为processed
	updates := map[string]interface{}{
		"status": "processed",
	}

	return s.walletRepo.UpdateWithdrawRequest(ctx, withdrawID, updates)
}

// ============ 健康检查 ============

// Health 健康检查
func (s *UnifiedWalletService) Health(ctx context.Context) error {
	return s.walletRepo.Health(ctx)
}

// ============ 辅助函数 ============

// generateOrderNo 生成订单号（简化版）
// TODO: 实际应使用雪花算法或其他分布式ID生成方案
func generateOrderNo() string {
	return fmt.Sprintf("ORDER%d", getCurrentTimestamp())
}

// getCurrentTimestamp 获取当前时间戳（毫秒）
func getCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}
