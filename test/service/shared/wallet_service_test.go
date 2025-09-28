package shared

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockWallet 模拟钱包模型
type MockWallet struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`
	Balance   int64              `bson:"balance" json:"balance"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockTransaction 模拟交易模型
type MockTransaction struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"userId"`
	WalletID      primitive.ObjectID `bson:"wallet_id" json:"walletId"`
	Type          string             `bson:"type" json:"type"`
	Amount        int64              `bson:"amount" json:"amount"`
	Status        string             `bson:"status" json:"status"`
	Description   string             `bson:"description" json:"description"`
	OrderID       string             `bson:"order_id" json:"orderId"`
	PaymentMethod string             `bson:"payment_method" json:"paymentMethod"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockWithdrawRecord 模拟提现记录模型
type MockWithdrawRecord struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"userId"`
	WalletID      primitive.ObjectID `bson:"wallet_id" json:"walletId"`
	Amount        int64              `bson:"amount" json:"amount"`
	Fee           int64              `bson:"fee" json:"fee"`
	ActualAmount  int64              `bson:"actual_amount" json:"actualAmount"`
	Status        string             `bson:"status" json:"status"`
	PaymentMethod string             `bson:"payment_method" json:"paymentMethod"`
	AccountInfo   map[string]string  `bson:"account_info" json:"accountInfo"`
	ProcessedAt   *time.Time         `bson:"processed_at" json:"processedAt"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockWalletRepository 模拟钱包仓储
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) Create(ctx context.Context, wallet *MockWallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) (*MockWallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockWallet), args.Error(1)
}

func (m *MockWalletRepository) UpdateBalance(ctx context.Context, walletID primitive.ObjectID, amount int64) error {
	args := m.Called(ctx, walletID, amount)
	return args.Error(0)
}

func (m *MockWalletRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*MockWallet, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockWallet), args.Error(1)
}

// MockTransactionRepository 模拟交易仓储
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *MockTransaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*MockTransaction, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockTransaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByOrderID(ctx context.Context, orderID string) (*MockTransaction, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockTransaction), args.Error(1)
}

func (m *MockTransactionRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// MockWithdrawRepository 模拟提现仓储
type MockWithdrawRepository struct {
	mock.Mock
}

func (m *MockWithdrawRepository) Create(ctx context.Context, record *MockWithdrawRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockWithdrawRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*MockWithdrawRecord, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockWithdrawRecord), args.Error(1)
}

func (m *MockWithdrawRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// MockWalletService 模拟钱包服务
type MockWalletService struct {
	walletRepo     *MockWalletRepository
	transactionRepo *MockTransactionRepository
	withdrawRepo   *MockWithdrawRepository
}

func NewMockWalletService(walletRepo *MockWalletRepository, transactionRepo *MockTransactionRepository, withdrawRepo *MockWithdrawRepository) *MockWalletService {
	return &MockWalletService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		withdrawRepo:    withdrawRepo,
	}
}

// CreateWallet 创建钱包
func (s *MockWalletService) CreateWallet(ctx context.Context, userID primitive.ObjectID) (*MockWallet, error) {
	wallet := &MockWallet{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Balance:   0,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.walletRepo.Create(ctx, wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetBalance 获取余额
func (s *MockWalletService) GetBalance(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return wallet.Balance, nil
}

// Recharge 充值
func (s *MockWalletService) Recharge(ctx context.Context, userID primitive.ObjectID, amount int64, orderID, paymentMethod string) (*MockTransaction, error) {
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}

	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 创建交易记录
	transaction := &MockTransaction{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		WalletID:      wallet.ID,
		Type:          "recharge",
		Amount:        amount,
		Status:        "pending",
		Description:   "钱包充值",
		OrderID:       orderID,
		PaymentMethod: paymentMethod,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// 更新钱包余额
	err = s.walletRepo.UpdateBalance(ctx, wallet.ID, wallet.Balance+amount)
	if err != nil {
		return nil, err
	}

	// 更新交易状态
	transaction.Status = "completed"
	err = s.transactionRepo.UpdateStatus(ctx, transaction.ID, "completed")
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// Consume 消费
func (s *MockWalletService) Consume(ctx context.Context, userID primitive.ObjectID, amount int64, description string) (*MockTransaction, error) {
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}

	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if wallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// 创建交易记录
	transaction := &MockTransaction{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		WalletID:    wallet.ID,
		Type:        "consume",
		Amount:      -amount,
		Status:      "completed",
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// 更新钱包余额
	err = s.walletRepo.UpdateBalance(ctx, wallet.ID, wallet.Balance-amount)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// Withdraw 提现
func (s *MockWalletService) Withdraw(ctx context.Context, userID primitive.ObjectID, amount int64, paymentMethod string, accountInfo map[string]string) (*MockWithdrawRecord, error) {
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}

	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 计算手续费（假设1%）
	fee := amount / 100
	actualAmount := amount - fee

	if wallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// 创建提现记录
	withdrawRecord := &MockWithdrawRecord{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		WalletID:      wallet.ID,
		Amount:        amount,
		Fee:           fee,
		ActualAmount:  actualAmount,
		Status:        "pending",
		PaymentMethod: paymentMethod,
		AccountInfo:   accountInfo,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.withdrawRepo.Create(ctx, withdrawRecord)
	if err != nil {
		return nil, err
	}

	// 冻结余额
	err = s.walletRepo.UpdateBalance(ctx, wallet.ID, wallet.Balance-amount)
	if err != nil {
		return nil, err
	}

	return withdrawRecord, nil
}

// GetTransactionHistory 获取交易历史
func (s *MockWalletService) GetTransactionHistory(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*MockTransaction, error) {
	return s.transactionRepo.GetByUserID(ctx, userID, limit, offset)
}

// 测试用例

func TestWalletService_CreateWallet_Success(t *testing.T) {
	walletRepo := new(MockWalletRepository)
	transactionRepo := new(MockTransactionRepository)
	withdrawRepo := new(MockWithdrawRepository)
	walletService := NewMockWalletService(walletRepo, transactionRepo, withdrawRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()

	// Mock 设置
	walletRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockWallet")).Return(nil)

	// 执行测试
	wallet, err := walletService.CreateWallet(ctx, userID)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, userID, wallet.UserID)
	assert.Equal(t, int64(0), wallet.Balance)
	assert.Equal(t, "active", wallet.Status)

	walletRepo.AssertExpectations(t)
}

func TestWalletService_GetBalance_Success(t *testing.T) {
	walletRepo := new(MockWalletRepository)
	transactionRepo := new(MockTransactionRepository)
	withdrawRepo := new(MockWithdrawRepository)
	walletService := NewMockWalletService(walletRepo, transactionRepo, withdrawRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	expectedBalance := int64(1000)

	wallet := &MockWallet{
		ID:      primitive.NewObjectID(),
		UserID:  userID,
		Balance: expectedBalance,
		Status:  "active",
	}

	// Mock 设置
	walletRepo.On("GetByUserID", ctx, userID).Return(wallet, nil)

	// 执行测试
	balance, err := walletService.GetBalance(ctx, userID)

	// 断言
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)

	walletRepo.AssertExpectations(t)
}

func TestWalletService_Recharge_Success(t *testing.T) {
	walletRepo := new(MockWalletRepository)
	transactionRepo := new(MockTransactionRepository)
	withdrawRepo := new(MockWithdrawRepository)
	walletService := NewMockWalletService(walletRepo, transactionRepo, withdrawRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	amount := int64(500)
	orderID := "order_123"
	paymentMethod := "alipay"

	wallet := &MockWallet{
		ID:      primitive.NewObjectID(),
		UserID:  userID,
		Balance: 1000,
		Status:  "active",
	}

	// Mock 设置
	walletRepo.On("GetByUserID", ctx, userID).Return(wallet, nil)
	transactionRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockTransaction")).Return(nil)
	walletRepo.On("UpdateBalance", ctx, wallet.ID, int64(1500)).Return(nil)
	transactionRepo.On("UpdateStatus", ctx, mock.AnythingOfType("primitive.ObjectID"), "completed").Return(nil)

	// 执行测试
	transaction, err := walletService.Recharge(ctx, userID, amount, orderID, paymentMethod)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, userID, transaction.UserID)
	assert.Equal(t, amount, transaction.Amount)
	assert.Equal(t, "recharge", transaction.Type)
	assert.Equal(t, orderID, transaction.OrderID)
	assert.Equal(t, paymentMethod, transaction.PaymentMethod)

	walletRepo.AssertExpectations(t)
	transactionRepo.AssertExpectations(t)
}

func TestWalletService_Recharge_InvalidAmount(t *testing.T) {
	walletRepo := new(MockWalletRepository)
	transactionRepo := new(MockTransactionRepository)
	withdrawRepo := new(MockWithdrawRepository)
	walletService := NewMockWalletService(walletRepo, transactionRepo, withdrawRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	amount := int64(-100)
	orderID := "order_123"
	paymentMethod := "alipay"

	// 执行测试
	transaction, err := walletService.Recharge(ctx, userID, amount, orderID, paymentMethod)

	// 断言
	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.Equal(t, "invalid amount", err.Error())
}

func TestWalletService_Consume_Success(t *testing.T) {
	walletRepo := new(MockWalletRepository)
	transactionRepo := new(MockTransactionRepository)
	withdrawRepo := new(MockWithdrawRepository)
	walletService := NewMockWalletService(walletRepo, transactionRepo, withdrawRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	amount := int64(200)
	description := "购买书籍"

	wallet := &MockWallet{
		ID:      primitive.NewObjectID(),
		UserID:  userID,
		Balance: 1000,
		Status:  "active",
	}

	// Mock 设置
	walletRepo.On("GetByUserID", ctx, userID).Return(wallet, nil)
	transactionRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockTransaction")).Return(nil)
	walletRepo.On("UpdateBalance", ctx, wallet.ID, int64(800)).Return(nil)

	// 执行测试
	transaction, err := walletService.Consume(ctx, userID, amount, description)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, userID, transaction.UserID)
	assert.Equal(t, -amount, transaction.Amount)
	assert.Equal(t, "consume", transaction.Type)
	assert.Equal(t, description, transaction.Description)
	assert.Equal(t, "completed", transaction.Status)

	walletRepo.AssertExpectations(t)
	transactionRepo.AssertExpectations(t)
}

func TestWalletService_Consume_InsufficientBalance(t *testing.T) {
	walletRepo := new(MockWalletRepository)
	transactionRepo := new(MockTransactionRepository)
	withdrawRepo := new(MockWithdrawRepository)
	walletService := NewMockWalletService(walletRepo, transactionRepo, withdrawRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	amount := int64(1500)
	description := "购买书籍"

	wallet := &MockWallet{
		ID:      primitive.NewObjectID(),
		UserID:  userID,
		Balance: 1000,
		Status:  "active",
	}

	// Mock 设置
	walletRepo.On("GetByUserID", ctx, userID).Return(wallet, nil)

	// 执行测试
	transaction, err := walletService.Consume(ctx, userID, amount, description)

	// 断言
	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.Equal(t, "insufficient balance", err.Error())

	walletRepo.AssertExpectations(t)
}

func TestWalletService_Withdraw_Success(t *testing.T) {
	walletRepo := new(MockWalletRepository)
	transactionRepo := new(MockTransactionRepository)
	withdrawRepo := new(MockWithdrawRepository)
	walletService := NewMockWalletService(walletRepo, transactionRepo, withdrawRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	amount := int64(1000)
	paymentMethod := "alipay"
	accountInfo := map[string]string{
		"account": "test@example.com",
		"name":    "测试用户",
	}

	wallet := &MockWallet{
		ID:      primitive.NewObjectID(),
		UserID:  userID,
		Balance: 2000,
		Status:  "active",
	}

	// Mock 设置
	walletRepo.On("GetByUserID", ctx, userID).Return(wallet, nil)
	withdrawRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockWithdrawRecord")).Return(nil)
	walletRepo.On("UpdateBalance", ctx, wallet.ID, int64(1000)).Return(nil)

	// 执行测试
	withdrawRecord, err := walletService.Withdraw(ctx, userID, amount, paymentMethod, accountInfo)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, withdrawRecord)
	assert.Equal(t, userID, withdrawRecord.UserID)
	assert.Equal(t, amount, withdrawRecord.Amount)
	assert.Equal(t, int64(10), withdrawRecord.Fee) // 1% 手续费
	assert.Equal(t, int64(990), withdrawRecord.ActualAmount)
	assert.Equal(t, "pending", withdrawRecord.Status)
	assert.Equal(t, paymentMethod, withdrawRecord.PaymentMethod)
	assert.Equal(t, accountInfo, withdrawRecord.AccountInfo)

	walletRepo.AssertExpectations(t)
	withdrawRepo.AssertExpectations(t)
}

func TestWalletService_GetTransactionHistory_Success(t *testing.T) {
	walletRepo := new(MockWalletRepository)
	transactionRepo := new(MockTransactionRepository)
	withdrawRepo := new(MockWithdrawRepository)
	walletService := NewMockWalletService(walletRepo, transactionRepo, withdrawRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	limit := 10
	offset := 0

	expectedTransactions := []*MockTransaction{
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Type:        "recharge",
			Amount:      500,
			Status:      "completed",
			Description: "钱包充值",
		},
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Type:        "consume",
			Amount:      -200,
			Status:      "completed",
			Description: "购买书籍",
		},
	}

	// Mock 设置
	transactionRepo.On("GetByUserID", ctx, userID, limit, offset).Return(expectedTransactions, nil)

	// 执行测试
	transactions, err := walletService.GetTransactionHistory(ctx, userID, limit, offset)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Len(t, transactions, 2)
	assert.Equal(t, expectedTransactions, transactions)

	transactionRepo.AssertExpectations(t)
}