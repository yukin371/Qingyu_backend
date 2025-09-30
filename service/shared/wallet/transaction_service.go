package wallet

import (
	"context"
	"fmt"

	walletModel "Qingyu_backend/models/shared/wallet"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// TransactionServiceImpl 交易服务实现
type TransactionServiceImpl struct {
	walletRepo sharedRepo.WalletRepository
}

// NewTransactionService 创建交易服务
func NewTransactionService(walletRepo sharedRepo.WalletRepository) TransactionService {
	return &TransactionServiceImpl{
		walletRepo: walletRepo,
	}
}

// ============ 交易操作 ============

// Recharge 充值
func (s *TransactionServiceImpl) Recharge(ctx context.Context, walletID string, amount float64, method, orderNo string) (*Transaction, error) {
	// 1. 验证金额
	if amount <= 0 {
		return nil, fmt.Errorf("充值金额必须大于0")
	}

	// 2. 获取钱包
	wallet, err := s.walletRepo.GetWallet(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("钱包不存在: %w", err)
	}

	// 3. 检查钱包状态
	if wallet.Status != "active" {
		return nil, fmt.Errorf("钱包已冻结，无法充值")
	}

	// 4. 创建交易记录
	transaction := &walletModel.Transaction{
		WalletID:    walletID,
		UserID:      wallet.UserID,
		Type:        "recharge",
		Amount:      amount,
		Method:      method,
		OrderNo:     orderNo,
		Status:      "success",
		Description: "充值",
	}

	if err := s.walletRepo.CreateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("创建交易记录失败: %w", err)
	}

	// 5. 更新余额
	if err := s.walletRepo.UpdateBalance(ctx, walletID, amount); err != nil {
		return nil, fmt.Errorf("更新余额失败: %w", err)
	}

	return convertToTransactionResponse(transaction), nil
}

// Consume 消费
func (s *TransactionServiceImpl) Consume(ctx context.Context, walletID string, amount float64, reason string) (*Transaction, error) {
	// 1. 验证金额
	if amount <= 0 {
		return nil, fmt.Errorf("消费金额必须大于0")
	}

	// 2. 获取钱包
	wallet, err := s.walletRepo.GetWallet(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("钱包不存在: %w", err)
	}

	// 3. 检查钱包状态
	if wallet.Status != "active" {
		return nil, fmt.Errorf("钱包已冻结，无法消费")
	}

	// 4. 检查余额
	if wallet.Balance < amount {
		return nil, fmt.Errorf("余额不足")
	}

	// 5. 创建交易记录
	transaction := &walletModel.Transaction{
		WalletID:    walletID,
		UserID:      wallet.UserID,
		Type:        "consume",
		Amount:      -amount, // 负数表示消费
		Status:      "success",
		Description: reason,
	}

	if err := s.walletRepo.CreateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("创建交易记录失败: %w", err)
	}

	// 6. 更新余额
	if err := s.walletRepo.UpdateBalance(ctx, walletID, -amount); err != nil {
		return nil, fmt.Errorf("更新余额失败: %w", err)
	}

	return convertToTransactionResponse(transaction), nil
}

// Transfer 转账
func (s *TransactionServiceImpl) Transfer(ctx context.Context, fromWalletID, toWalletID string, amount float64, reason string) error {
	// 1. 验证金额
	if amount <= 0 {
		return fmt.Errorf("转账金额必须大于0")
	}

	// 2. 获取源钱包
	fromWallet, err := s.walletRepo.GetWallet(ctx, fromWalletID)
	if err != nil {
		return fmt.Errorf("源钱包不存在: %w", err)
	}

	// 3. 获取目标钱包
	toWallet, err := s.walletRepo.GetWallet(ctx, toWalletID)
	if err != nil {
		return fmt.Errorf("目标钱包不存在: %w", err)
	}

	// 4. 检查钱包状态
	if fromWallet.Status != "active" || toWallet.Status != "active" {
		return fmt.Errorf("钱包已冻结，无法转账")
	}

	// 5. 检查余额
	if fromWallet.Balance < amount {
		return fmt.Errorf("余额不足")
	}

	// 6. 创建转出交易记录
	outTransaction := &walletModel.Transaction{
		WalletID:    fromWalletID,
		UserID:      fromWallet.UserID,
		Type:        "transfer_out",
		Amount:      -amount,
		Status:      "success",
		Description: "转账给 " + toWallet.UserID + ": " + reason,
	}

	if err := s.walletRepo.CreateTransaction(ctx, outTransaction); err != nil {
		return fmt.Errorf("创建转出记录失败: %w", err)
	}

	// 7. 创建转入交易记录
	inTransaction := &walletModel.Transaction{
		WalletID:    toWalletID,
		UserID:      toWallet.UserID,
		Type:        "transfer_in",
		Amount:      amount,
		Status:      "success",
		Description: "来自 " + fromWallet.UserID + " 的转账: " + reason,
	}

	if err := s.walletRepo.CreateTransaction(ctx, inTransaction); err != nil {
		return fmt.Errorf("创建转入记录失败: %w", err)
	}

	// 8. 更新余额
	if err := s.walletRepo.UpdateBalance(ctx, fromWalletID, -amount); err != nil {
		return fmt.Errorf("更新源钱包余额失败: %w", err)
	}

	if err := s.walletRepo.UpdateBalance(ctx, toWalletID, amount); err != nil {
		// TODO: 需要回滚
		return fmt.Errorf("更新目标钱包余额失败: %w", err)
	}

	return nil
}

// GetTransaction 获取交易记录
func (s *TransactionServiceImpl) GetTransaction(ctx context.Context, transactionID string) (*Transaction, error) {
	transaction, err := s.walletRepo.GetTransaction(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("获取交易记录失败: %w", err)
	}

	return convertToTransactionResponse(transaction), nil
}

// ListTransactions 列出交易记录
func (s *TransactionServiceImpl) ListTransactions(ctx context.Context, walletID string, limit, offset int) ([]*Transaction, error) {
	filter := &sharedRepo.TransactionFilter{
		WalletID: walletID,
		Limit:    limit,
		Offset:   offset,
	}

	transactions, err := s.walletRepo.ListTransactions(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("获取交易列表失败: %w", err)
	}

	result := make([]*Transaction, len(transactions))
	for i, t := range transactions {
		result[i] = convertToTransactionResponse(t)
	}

	return result, nil
}

// ============ 辅助函数 ============

// convertToTransactionResponse 转换为响应格式
func convertToTransactionResponse(transaction *walletModel.Transaction) *Transaction {
	return &Transaction{
		ID:          transaction.ID,
		WalletID:    transaction.WalletID,
		UserID:      transaction.UserID,
		Type:        transaction.Type,
		Amount:      transaction.Amount,
		Method:      transaction.Method,
		OrderNo:     transaction.OrderNo,
		Status:      transaction.Status,
		Description: transaction.Description,
		CreatedAt:   transaction.CreatedAt,
	}
}
