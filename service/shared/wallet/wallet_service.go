package wallet

import (
	"context"
	"fmt"

	walletModel "Qingyu_backend/models/shared/wallet"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// WalletServiceImpl 钱包服务实现
type WalletServiceImpl struct {
	walletRepo sharedRepo.WalletRepository
}

// NewWalletService 创建钱包服务
func NewWalletService(walletRepo sharedRepo.WalletRepository) WalletService {
	return &WalletServiceImpl{
		walletRepo: walletRepo,
	}
}

// ============ 钱包管理 ============

// CreateWallet 创建钱包
func (s *WalletServiceImpl) CreateWallet(ctx context.Context, userID string) (*Wallet, error) {
	// 1. 检查是否已存在钱包
	_, err := s.walletRepo.GetWalletByUserID(ctx, userID)
	if err == nil {
		return nil, fmt.Errorf("用户已有钱包")
	}

	// 2. 创建钱包
	wallet := &walletModel.Wallet{
		UserID:  userID,
		Balance: 0,
		Status:  "active",
	}

	if err := s.walletRepo.CreateWallet(ctx, wallet); err != nil {
		return nil, fmt.Errorf("创建钱包失败: %w", err)
	}

	return convertToWalletResponse(wallet), nil
}

// GetWallet 获取钱包
func (s *WalletServiceImpl) GetWallet(ctx context.Context, walletID string) (*Wallet, error) {
	wallet, err := s.walletRepo.GetWallet(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("获取钱包失败: %w", err)
	}

	return convertToWalletResponse(wallet), nil
}

// GetWalletByUserID 根据用户ID获取钱包
func (s *WalletServiceImpl) GetWalletByUserID(ctx context.Context, userID string) (*Wallet, error) {
	wallet, err := s.walletRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取钱包失败: %w", err)
	}

	return convertToWalletResponse(wallet), nil
}

// GetBalance 获取余额
func (s *WalletServiceImpl) GetBalance(ctx context.Context, walletID string) (float64, error) {
	wallet, err := s.walletRepo.GetWallet(ctx, walletID)
	if err != nil {
		return 0, fmt.Errorf("获取余额失败: %w", err)
	}

	return wallet.Balance, nil
}

// FreezeWallet 冻结钱包
func (s *WalletServiceImpl) FreezeWallet(ctx context.Context, walletID string, reason string) error {
	updates := map[string]interface{}{
		"status": "frozen",
	}

	if err := s.walletRepo.UpdateWallet(ctx, walletID, updates); err != nil {
		return fmt.Errorf("冻结钱包失败: %w", err)
	}

	return nil
}

// UnfreezeWallet 解冻钱包
func (s *WalletServiceImpl) UnfreezeWallet(ctx context.Context, walletID string) error {
	updates := map[string]interface{}{
		"status": "active",
	}

	if err := s.walletRepo.UpdateWallet(ctx, walletID, updates); err != nil {
		return fmt.Errorf("解冻钱包失败: %w", err)
	}

	return nil
}

// ============ 辅助函数 ============

// convertToWalletResponse 转换为响应格式
func convertToWalletResponse(wallet *walletModel.Wallet) *Wallet {
	return &Wallet{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		Status:    wallet.Status,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}
