package wallet

import (
	"context"
	"fmt"
	"time"

	walletModel "Qingyu_backend/models/shared/wallet"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// WithdrawServiceImpl 提现服务实现
type WithdrawServiceImpl struct {
	walletRepo sharedRepo.WalletRepository
}

// NewWithdrawService 创建提现服务
func NewWithdrawService(walletRepo sharedRepo.WalletRepository) WithdrawService {
	return &WithdrawServiceImpl{
		walletRepo: walletRepo,
	}
}

// ============ 提现操作 ============

// CreateWithdrawRequest 创建提现请求
func (s *WithdrawServiceImpl) CreateWithdrawRequest(ctx context.Context, userID, walletID string, amount float64, method, account string) (*WithdrawRequest, error) {
	// 1. 验证金额
	if amount <= 0 {
		return nil, fmt.Errorf("提现金额必须大于0")
	}

	// 最小提现金额
	if amount < 10 {
		return nil, fmt.Errorf("提现金额不能小于10元")
	}

	// 2. 获取钱包
	wallet, err := s.walletRepo.GetWallet(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("钱包不存在: %w", err)
	}

	// 3. 检查钱包状态
	if wallet.Status != "active" {
		return nil, fmt.Errorf("钱包已冻结，无法提现")
	}

	// 4. 检查余额
	if wallet.Balance < amount {
		return nil, fmt.Errorf("余额不足")
	}

	// 5. 创建提现请求
	request := &walletModel.WithdrawRequest{
		UserID:   userID,
		WalletID: walletID,
		Amount:   amount,
		Method:   method,
		Account:  account,
		Status:   "pending",
	}

	if err := s.walletRepo.CreateWithdrawRequest(ctx, request); err != nil {
		return nil, fmt.Errorf("创建提现请求失败: %w", err)
	}

	// 6. 冻结提现金额（从余额中扣除，待审核通过后实际提现）
	if err := s.walletRepo.UpdateBalance(ctx, walletID, -amount); err != nil {
		return nil, fmt.Errorf("冻结提现金额失败: %w", err)
	}

	return convertToWithdrawResponse(request), nil
}

// ApproveWithdraw 审核通过提现
func (s *WithdrawServiceImpl) ApproveWithdraw(ctx context.Context, requestID, reviewerID, remark string) error {
	// 1. 获取提现请求
	request, err := s.walletRepo.GetWithdrawRequest(ctx, requestID)
	if err != nil {
		return fmt.Errorf("提现请求不存在: %w", err)
	}

	// 2. 检查状态
	if request.Status != "pending" {
		return fmt.Errorf("提现请求状态异常: %s", request.Status)
	}

	// 3. 更新状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":       "approved",
		"reviewer_id":  reviewerID,
		"reviewed_at":  now,
		"remark":       remark,
		"processed_at": now,
	}

	if err := s.walletRepo.UpdateWithdrawRequest(ctx, requestID, updates); err != nil {
		return fmt.Errorf("更新提现请求失败: %w", err)
	}

	// 4. 创建提现交易记录
	transaction := &walletModel.Transaction{
		WalletID:    request.WalletID,
		UserID:      request.UserID,
		Type:        "withdraw",
		Amount:      -request.Amount,
		Method:      request.Method,
		Status:      "success",
		Description: "提现",
	}

	if err := s.walletRepo.CreateTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("创建交易记录失败: %w", err)
	}

	return nil
}

// RejectWithdraw 拒绝提现
func (s *WithdrawServiceImpl) RejectWithdraw(ctx context.Context, requestID, reviewerID, reason string) error {
	// 1. 获取提现请求
	request, err := s.walletRepo.GetWithdrawRequest(ctx, requestID)
	if err != nil {
		return fmt.Errorf("提现请求不存在: %w", err)
	}

	// 2. 检查状态
	if request.Status != "pending" {
		return fmt.Errorf("提现请求状态异常: %s", request.Status)
	}

	// 3. 退还金额
	if err := s.walletRepo.UpdateBalance(ctx, request.WalletID, request.Amount); err != nil {
		return fmt.Errorf("退还金额失败: %w", err)
	}

	// 4. 更新状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":      "rejected",
		"reviewer_id": reviewerID,
		"reviewed_at": now,
		"remark":      reason,
	}

	if err := s.walletRepo.UpdateWithdrawRequest(ctx, requestID, updates); err != nil {
		return fmt.Errorf("更新提现请求失败: %w", err)
	}

	return nil
}

// GetWithdrawRequest 获取提现请求
func (s *WithdrawServiceImpl) GetWithdrawRequest(ctx context.Context, requestID string) (*WithdrawRequest, error) {
	request, err := s.walletRepo.GetWithdrawRequest(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("获取提现请求失败: %w", err)
	}

	return convertToWithdrawResponse(request), nil
}

// ListWithdrawRequests 列出提现请求
func (s *WithdrawServiceImpl) ListWithdrawRequests(ctx context.Context, userID, status string, limit, offset int) ([]*WithdrawRequest, error) {
	filter := &sharedRepo.WithdrawFilter{
		UserID: userID,
		Status: status,
		Limit:  limit,
		Offset: offset,
	}

	requests, err := s.walletRepo.ListWithdrawRequests(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("获取提现列表失败: %w", err)
	}

	result := make([]*WithdrawRequest, len(requests))
	for i, r := range requests {
		result[i] = convertToWithdrawResponse(r)
	}

	return result, nil
}

// ============ 辅助函数 ============

// convertToWithdrawResponse 转换为响应格式
func convertToWithdrawResponse(request *walletModel.WithdrawRequest) *WithdrawRequest {
	return &WithdrawRequest{
		ID:          request.ID,
		UserID:      request.UserID,
		WalletID:    request.WalletID,
		Amount:      request.Amount,
		Method:      request.Method,
		Account:     request.Account,
		Status:      request.Status,
		ReviewerID:  request.ReviewerID,
		ReviewedAt:  request.ReviewedAt,
		ProcessedAt: request.ProcessedAt,
		Remark:      request.Remark,
		CreatedAt:   request.CreatedAt,
		UpdatedAt:   request.UpdatedAt,
	}
}
