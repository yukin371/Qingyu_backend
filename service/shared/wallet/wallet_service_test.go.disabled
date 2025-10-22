package wallet

import (
	"context"
	"fmt"
	"testing"
	"time"

	walletModel "Qingyu_backend/models/shared/wallet"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// MockWalletRepository Mock钱包Repository
type MockWalletRepository struct {
	wallets           map[string]*walletModel.Wallet
	userWallets       map[string]*walletModel.Wallet
	transactions      map[string]*walletModel.Transaction
	withdrawRequests  map[string]*walletModel.WithdrawRequest
	shouldReturnError bool
}

func NewMockWalletRepository() *MockWalletRepository {
	return &MockWalletRepository{
		wallets:          make(map[string]*walletModel.Wallet),
		userWallets:      make(map[string]*walletModel.Wallet),
		transactions:     make(map[string]*walletModel.Transaction),
		withdrawRequests: make(map[string]*walletModel.WithdrawRequest),
	}
}

func (m *MockWalletRepository) CreateWallet(ctx context.Context, wallet *walletModel.Wallet) error {
	if m.shouldReturnError {
		return fmt.Errorf("mock error")
	}
	wallet.ID = "wallet_" + wallet.UserID
	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()
	m.wallets[wallet.ID] = wallet
	m.userWallets[wallet.UserID] = wallet
	return nil
}

func (m *MockWalletRepository) GetWallet(ctx context.Context, walletID string) (*walletModel.Wallet, error) {
	if m.shouldReturnError {
		return nil, fmt.Errorf("mock error")
	}
	wallet, exists := m.wallets[walletID]
	if !exists {
		return nil, fmt.Errorf("钱包不存在")
	}
	return wallet, nil
}

func (m *MockWalletRepository) GetWalletByUserID(ctx context.Context, userID string) (*walletModel.Wallet, error) {
	if m.shouldReturnError {
		return nil, fmt.Errorf("mock error")
	}
	wallet, exists := m.userWallets[userID]
	if !exists {
		return nil, fmt.Errorf("钱包不存在")
	}
	return wallet, nil
}

func (m *MockWalletRepository) UpdateWallet(ctx context.Context, walletID string, updates map[string]interface{}) error {
	if m.shouldReturnError {
		return fmt.Errorf("mock error")
	}
	wallet, exists := m.wallets[walletID]
	if !exists {
		return fmt.Errorf("钱包不存在")
	}

	if frozen, ok := updates["frozen"].(bool); ok {
		wallet.Frozen = frozen
	}
	wallet.UpdatedAt = time.Now()
	return nil
}

func (m *MockWalletRepository) UpdateBalance(ctx context.Context, walletID string, amount float64) error {
	if m.shouldReturnError {
		return fmt.Errorf("mock error")
	}
	wallet, exists := m.wallets[walletID]
	if !exists {
		return fmt.Errorf("钱包不存在")
	}
	wallet.Balance += amount
	wallet.UpdatedAt = time.Now()
	return nil
}

func (m *MockWalletRepository) CreateTransaction(ctx context.Context, transaction *walletModel.Transaction) error {
	if m.shouldReturnError {
		return fmt.Errorf("mock error")
	}
	transaction.ID = fmt.Sprintf("tx_%d", len(m.transactions)+1)
	transaction.CreatedAt = time.Now()
	m.transactions[transaction.ID] = transaction
	return nil
}

func (m *MockWalletRepository) GetTransaction(ctx context.Context, transactionID string) (*walletModel.Transaction, error) {
	if m.shouldReturnError {
		return nil, fmt.Errorf("mock error")
	}
	tx, exists := m.transactions[transactionID]
	if !exists {
		return nil, fmt.Errorf("交易不存在")
	}
	return tx, nil
}

func (m *MockWalletRepository) ListTransactions(ctx context.Context, filter *sharedRepo.TransactionFilter) ([]*walletModel.Transaction, error) {
	if m.shouldReturnError {
		return nil, fmt.Errorf("mock error")
	}
	var result []*walletModel.Transaction
	for _, tx := range m.transactions {
		if filter.UserID != "" && tx.UserID != filter.UserID {
			continue
		}
		result = append(result, tx)
	}
	return result, nil
}

func (m *MockWalletRepository) CountTransactions(ctx context.Context, filter *sharedRepo.TransactionFilter) (int64, error) {
	return int64(len(m.transactions)), nil
}

func (m *MockWalletRepository) CreateWithdrawRequest(ctx context.Context, request *walletModel.WithdrawRequest) error {
	if m.shouldReturnError {
		return fmt.Errorf("mock error")
	}
	request.ID = fmt.Sprintf("wd_%d", len(m.withdrawRequests)+1)
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	m.withdrawRequests[request.ID] = request
	return nil
}

func (m *MockWalletRepository) GetWithdrawRequest(ctx context.Context, requestID string) (*walletModel.WithdrawRequest, error) {
	if m.shouldReturnError {
		return nil, fmt.Errorf("mock error")
	}
	req, exists := m.withdrawRequests[requestID]
	if !exists {
		return nil, fmt.Errorf("提现请求不存在")
	}
	return req, nil
}

func (m *MockWalletRepository) UpdateWithdrawRequest(ctx context.Context, requestID string, updates map[string]interface{}) error {
	if m.shouldReturnError {
		return fmt.Errorf("mock error")
	}
	req, exists := m.withdrawRequests[requestID]
	if !exists {
		return fmt.Errorf("提现请求不存在")
	}

	if status, ok := updates["status"].(string); ok {
		req.Status = status
	}
	if reviewedBy, ok := updates["reviewed_by"].(string); ok {
		req.ReviewedBy = reviewedBy
	}
	req.UpdatedAt = time.Now()
	return nil
}

func (m *MockWalletRepository) ListWithdrawRequests(ctx context.Context, filter *sharedRepo.WithdrawFilter) ([]*walletModel.WithdrawRequest, error) {
	if m.shouldReturnError {
		return nil, fmt.Errorf("mock error")
	}
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

func (m *MockWalletRepository) Health(ctx context.Context) error {
	return nil
}

// ============ 测试用例 ============

// TestCreateWallet 测试创建钱包
func TestCreateWallet(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	service := NewWalletService(mockRepo)

	ctx := context.Background()
	wallet, err := service.CreateWallet(ctx, "user123")

	if err != nil {
		t.Fatalf("创建钱包失败: %v", err)
	}

	if wallet.UserID != "user123" {
		t.Errorf("用户ID错误: %s", wallet.UserID)
	}

	if wallet.Balance != 0 {
		t.Errorf("初始余额应为0，实际: %f", wallet.Balance)
	}

	t.Logf("创建钱包成功: %+v", wallet)
}

// TestCreateWallet_Duplicate 测试重复创建钱包
func TestCreateWallet_Duplicate(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	service := NewWalletService(mockRepo)

	ctx := context.Background()

	// 第一次创建
	_, err := service.CreateWallet(ctx, "user123")
	if err != nil {
		t.Fatalf("第一次创建钱包失败: %v", err)
	}

	// 第二次创建（应该失败）
	_, err = service.CreateWallet(ctx, "user123")
	if err == nil {
		t.Error("重复创建钱包应该失败")
	}

	t.Logf("重复创建钱包测试通过: %v", err)
}

// TestGetWallet 测试获取钱包
func TestGetWallet(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	service := NewWalletService(mockRepo)

	ctx := context.Background()

	// 先创建钱包
	created, _ := service.CreateWallet(ctx, "user123")

	// 获取钱包
	wallet, err := service.GetWallet(ctx, created.ID)
	if err != nil {
		t.Fatalf("获取钱包失败: %v", err)
	}

	if wallet.UserID != "user123" {
		t.Errorf("用户ID错误: %s", wallet.UserID)
	}

	t.Logf("获取钱包成功: %+v", wallet)
}

// TestGetWalletByUserID 测试根据用户ID获取钱包
func TestGetWalletByUserID(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	service := NewWalletService(mockRepo)

	ctx := context.Background()

	// 先创建钱包
	_, _ = service.CreateWallet(ctx, "user123")

	// 根据用户ID获取钱包
	wallet, err := service.GetWalletByUserID(ctx, "user123")
	if err != nil {
		t.Fatalf("获取钱包失败: %v", err)
	}

	if wallet.UserID != "user123" {
		t.Errorf("用户ID错误: %s", wallet.UserID)
	}

	t.Logf("根据用户ID获取钱包成功: %+v", wallet)
}

// TestGetBalance 测试获取余额
func TestGetBalance(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	service := NewWalletService(mockRepo)

	ctx := context.Background()

	// 创建钱包
	wallet, _ := service.CreateWallet(ctx, "user123")

	// 获取余额
	balance, err := service.GetBalance(ctx, wallet.ID)
	if err != nil {
		t.Fatalf("获取余额失败: %v", err)
	}

	if balance != 0 {
		t.Errorf("初始余额应为0，实际: %f", balance)
	}

	// 增加余额
	mockRepo.UpdateBalance(ctx, wallet.ID, 100.0)

	// 再次获取余额
	balance, err = service.GetBalance(ctx, wallet.ID)
	if err != nil {
		t.Fatalf("获取余额失败: %v", err)
	}

	if balance != 100.0 {
		t.Errorf("余额应为100，实际: %f", balance)
	}

	t.Logf("获取余额测试通过: balance=%.2f", balance)
}

// TestFreezeWallet 测试冻结钱包
func TestFreezeWallet(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	service := NewWalletService(mockRepo)

	ctx := context.Background()

	// 创建钱包
	wallet, _ := service.CreateWallet(ctx, "user123")

	// 冻结钱包
	err := service.FreezeWallet(ctx, wallet.ID, "违规操作")
	if err != nil {
		t.Fatalf("冻结钱包失败: %v", err)
	}

	// 验证冻结状态
	frozen, _ := mockRepo.GetWallet(ctx, wallet.ID)
	if !frozen.Frozen {
		t.Error("钱包应该被冻结")
	}

	t.Logf("冻结钱包测试通过")
}

// TestUnfreezeWallet 测试解冻钱包
func TestUnfreezeWallet(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	service := NewWalletService(mockRepo)

	ctx := context.Background()

	// 创建并冻结钱包
	wallet, _ := service.CreateWallet(ctx, "user123")
	service.FreezeWallet(ctx, wallet.ID, "测试")

	// 解冻钱包
	err := service.UnfreezeWallet(ctx, wallet.ID)
	if err != nil {
		t.Fatalf("解冻钱包失败: %v", err)
	}

	// 验证解冻状态
	unfrozen, _ := mockRepo.GetWallet(ctx, wallet.ID)
	if unfrozen.Frozen {
		t.Error("钱包应该被解冻")
	}

	t.Logf("解冻钱包测试通过")
}

// TestGetWallet_NotFound 测试获取不存在的钱包
func TestGetWallet_NotFound(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	service := NewWalletService(mockRepo)

	ctx := context.Background()

	_, err := service.GetWallet(ctx, "nonexistent")
	if err == nil {
		t.Error("获取不存在的钱包应该失败")
	}

	t.Logf("获取不存在的钱包测试通过: %v", err)
}
