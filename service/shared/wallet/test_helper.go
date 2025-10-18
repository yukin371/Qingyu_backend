package wallet

import (
	"context"
	"fmt"
	"time"

	walletModel "Qingyu_backend/models/shared/wallet"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// MockWalletRepositoryV2 模拟WalletRepository（用于统一测试）
type MockWalletRepositoryV2 struct {
	wallets          map[string]*walletModel.Wallet
	transactions     map[string]*walletModel.Transaction
	withdrawRequests map[string]*walletModel.WithdrawRequest
	txCounter        int // 交易ID计数器，确保唯一性
}

// NewMockWalletRepositoryV2 创建Mock Repository
func NewMockWalletRepositoryV2() *MockWalletRepositoryV2 {
	return &MockWalletRepositoryV2{
		wallets:          make(map[string]*walletModel.Wallet),
		transactions:     make(map[string]*walletModel.Transaction),
		withdrawRequests: make(map[string]*walletModel.WithdrawRequest),
	}
}

// CreateWallet 创建钱包
func (m *MockWalletRepositoryV2) CreateWallet(ctx context.Context, wallet *walletModel.Wallet) error {
	wallet.ID = "wallet_" + wallet.UserID
	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()
	m.wallets[wallet.UserID] = wallet
	return nil
}

// GetWallet 获取钱包（根据用户ID）
func (m *MockWalletRepositoryV2) GetWallet(ctx context.Context, userID string) (*walletModel.Wallet, error) {
	if wallet, ok := m.wallets[userID]; ok {
		return wallet, nil
	}
	return nil, nil
}

// UpdateWallet 更新钱包
func (m *MockWalletRepositoryV2) UpdateWallet(ctx context.Context, walletID string, updates map[string]interface{}) error {
	for _, wallet := range m.wallets {
		if wallet.ID == walletID {
			if frozen, ok := updates["frozen"].(bool); ok {
				wallet.Frozen = frozen
			}
			wallet.UpdatedAt = time.Now()
			return nil
		}
	}
	return nil
}

// UpdateBalance 更新余额
// 注意：在当前实现中，walletID参数实际上可以是userID或wallet.ID
func (m *MockWalletRepositoryV2) UpdateBalance(ctx context.Context, walletID string, amount float64) error {
	for _, wallet := range m.wallets {
		// 支持通过wallet.ID或userID查找
		if wallet.ID == walletID || wallet.UserID == walletID {
			wallet.Balance += amount
			wallet.UpdatedAt = time.Now()
			return nil
		}
	}
	return nil
}

// CreateTransaction 创建交易
func (m *MockWalletRepositoryV2) CreateTransaction(ctx context.Context, transaction *walletModel.Transaction) error {
	// 使用计数器和时间戳确保ID唯一性
	m.txCounter++
	transaction.ID = fmt.Sprintf("tx_%s_%d", time.Now().Format("20060102150405"), m.txCounter)
	transaction.CreatedAt = time.Now()
	transaction.TransactionTime = time.Now()
	m.transactions[transaction.ID] = transaction
	return nil
}

// GetTransaction 获取交易
func (m *MockWalletRepositoryV2) GetTransaction(ctx context.Context, transactionID string) (*walletModel.Transaction, error) {
	if tx, ok := m.transactions[transactionID]; ok {
		return tx, nil
	}
	return nil, nil
}

// ListTransactions 列出交易
func (m *MockWalletRepositoryV2) ListTransactions(ctx context.Context, filter *sharedRepo.TransactionFilter) ([]*walletModel.Transaction, error) {
	var result []*walletModel.Transaction
	for _, tx := range m.transactions {
		if filter.UserID != "" && tx.UserID != filter.UserID {
			continue
		}
		result = append(result, tx)
	}

	// 按时间倒序排序（最新的在前）
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].CreatedAt.Before(result[j].CreatedAt) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	// 应用分页
	if filter.Limit > 0 {
		start := int(filter.Offset)
		end := start + int(filter.Limit)
		if start > len(result) {
			return []*walletModel.Transaction{}, nil
		}
		if end > len(result) {
			end = len(result)
		}
		result = result[start:end]
	}

	return result, nil
}

// CountTransactions 统计交易
func (m *MockWalletRepositoryV2) CountTransactions(ctx context.Context, filter *sharedRepo.TransactionFilter) (int64, error) {
	count := int64(0)
	for _, tx := range m.transactions {
		if filter.UserID != "" && tx.UserID != filter.UserID {
			continue
		}
		count++
	}
	return count, nil
}

// CreateWithdrawRequest 创建提现请求
func (m *MockWalletRepositoryV2) CreateWithdrawRequest(ctx context.Context, request *walletModel.WithdrawRequest) error {
	request.ID = "wd_" + time.Now().Format("20060102150405")
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	m.withdrawRequests[request.ID] = request
	return nil
}

// GetWithdrawRequest 获取提现请求
func (m *MockWalletRepositoryV2) GetWithdrawRequest(ctx context.Context, withdrawID string) (*walletModel.WithdrawRequest, error) {
	if req, ok := m.withdrawRequests[withdrawID]; ok {
		return req, nil
	}
	return nil, nil
}

// UpdateWithdrawRequest 更新提现请求
func (m *MockWalletRepositoryV2) UpdateWithdrawRequest(ctx context.Context, withdrawID string, updates map[string]interface{}) error {
	if req, ok := m.withdrawRequests[withdrawID]; ok {
		if status, ok := updates["status"].(string); ok {
			req.Status = status
		}
		req.UpdatedAt = time.Now()
		return nil
	}
	return nil
}

// ListWithdrawRequests 列出提现请求
func (m *MockWalletRepositoryV2) ListWithdrawRequests(ctx context.Context, filter *sharedRepo.WithdrawFilter) ([]*walletModel.WithdrawRequest, error) {
	var result []*walletModel.WithdrawRequest
	for _, req := range m.withdrawRequests {
		if filter.UserID != "" && req.UserID != filter.UserID {
			continue
		}
		if filter.Status != "" && req.Status != filter.Status {
			continue
		}
		result = append(result, req)
	}
	return result, nil
}

// CountWithdrawRequests 统计提现请求
func (m *MockWalletRepositoryV2) CountWithdrawRequests(ctx context.Context, filter *sharedRepo.WithdrawFilter) (int64, error) {
	count := int64(0)
	for _, req := range m.withdrawRequests {
		if filter.UserID != "" && req.UserID != filter.UserID {
			continue
		}
		if filter.Status != "" && req.Status != filter.Status {
			continue
		}
		count++
	}
	return count, nil
}

// Health 健康检查
func (m *MockWalletRepositoryV2) Health(ctx context.Context) error {
	return nil
}
