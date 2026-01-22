package wallet

import (
	financeModel "Qingyu_backend/models/finance"
	"context"
	"fmt"

	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// WalletServiceImpl 钱包服务实现
type WalletServiceImpl struct {
	walletRepo sharedRepo.WalletRepository
}

// NewWalletService 创建钱包服务（内部使用）
// 注意：对外应该使用NewUnifiedWalletService
func NewWalletService(walletRepo sharedRepo.WalletRepository) *WalletServiceImpl {
	return &WalletServiceImpl{
		walletRepo: walletRepo,
	}
}

// ============ 钱包管理 ============

// CreateWallet 创建钱包
func (s *WalletServiceImpl) CreateWallet(ctx context.Context, userID string) (*Wallet, error) {
	// 1. 检查是否已存在钱包
	existingWallet, err := s.walletRepo.GetWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("检查钱包失败: %w", err)
	}
	if existingWallet != nil {
		return nil, fmt.Errorf("用户已有钱包")
	}

	// 2. 创建钱包
	wallet := &financeModel.Wallet{
		UserID:  userID,
		Balance: 0,
		Frozen:  false,
	}

	if err := s.walletRepo.CreateWallet(ctx, wallet); err != nil {
		return nil, fmt.Errorf("创建钱包失败: %w", err)
	}

	return convertToWalletResponse(wallet), nil
}

// GetWallet 根据用户ID获取钱包
func (s *WalletServiceImpl) GetWallet(ctx context.Context, userID string) (*Wallet, error) {
	wallet, err := s.walletRepo.GetWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取钱包失败: %w", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("钱包不存在")
	}

	return convertToWalletResponse(wallet), nil
}

// GetWalletByID 根据钱包ID获取钱包（内部使用）
func (s *WalletServiceImpl) GetWalletByID(ctx context.Context, walletID string) (*Wallet, error) {
	// 遍历查找钱包ID对应的钱包
	// 注意：这是一个内部辅助方法，实际应该通过Repository的GetWalletByID
	// 这里简化处理
	return nil, fmt.Errorf("未实现")
}

// GetBalance 获取余额（根据用户ID）
func (s *WalletServiceImpl) GetBalance(ctx context.Context, userID string) (int64, error) {
	wallet, err := s.walletRepo.GetWallet(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("获取余额失败: %w", err)
	}
	if wallet == nil {
		return 0, fmt.Errorf("钱包不存在")
	}

	return wallet.Balance, nil
}

// FreezeWallet 冻结钱包（根据用户ID）
func (s *WalletServiceImpl) FreezeWallet(ctx context.Context, userID string, reason string) error {
	// 1. 获取钱包
	wallet, err := s.walletRepo.GetWallet(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取钱包失败: %w", err)
	}
	if wallet == nil {
		return fmt.Errorf("钱包不存在")
	}

	// 2. 冻结钱包
	updates := map[string]interface{}{
		"frozen": true,
	}

	if err := s.walletRepo.UpdateWallet(ctx, wallet.ID, updates); err != nil {
		return fmt.Errorf("冻结钱包失败: %w", err)
	}

	return nil
}

// UnfreezeWallet 解冻钱包（根据用户ID）
func (s *WalletServiceImpl) UnfreezeWallet(ctx context.Context, userID string) error {
	// 1. 获取钱包
	wallet, err := s.walletRepo.GetWallet(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取钱包失败: %w", err)
	}
	if wallet == nil {
		return fmt.Errorf("钱包不存在")
	}

	// 2. 解冻钱包
	updates := map[string]interface{}{
		"frozen": false,
	}

	if err := s.walletRepo.UpdateWallet(ctx, wallet.ID, updates); err != nil {
		return fmt.Errorf("解冻钱包失败: %w", err)
	}

	return nil
}

// ============ 辅助函数 ============

// convertToWalletResponse 转换为响应格式
func convertToWalletResponse(wallet *financeModel.Wallet) *Wallet {
	return &Wallet{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		Frozen:    wallet.Frozen,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}
