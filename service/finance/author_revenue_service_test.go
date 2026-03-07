package finance

import (
	"context"
	"errors"
	"testing"

	financeModel "Qingyu_backend/models/finance"
	"Qingyu_backend/models/shared/types"
	financeRepo "Qingyu_backend/repository/interfaces/finance"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockAuthorRevenueRepository struct {
	withdrawals map[string]*financeModel.WithdrawalRequest
	failCreate  error
	counter     int
}

func newMockAuthorRevenueRepository() *mockAuthorRevenueRepository {
	return &mockAuthorRevenueRepository{
		withdrawals: make(map[string]*financeModel.WithdrawalRequest),
	}
}

func (m *mockAuthorRevenueRepository) CreateEarning(ctx context.Context, earning *financeModel.AuthorEarning) error {
	return nil
}

func (m *mockAuthorRevenueRepository) GetEarning(ctx context.Context, earningID primitive.ObjectID) (*financeModel.AuthorEarning, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAuthorRevenueRepository) ListEarnings(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error) {
	return nil, 0, nil
}

func (m *mockAuthorRevenueRepository) UpdateEarning(ctx context.Context, earningID primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}

func (m *mockAuthorRevenueRepository) BatchUpdateEarnings(ctx context.Context, earningIDs []primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}

func (m *mockAuthorRevenueRepository) GetEarningsByAuthor(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error) {
	return nil, 0, nil
}

func (m *mockAuthorRevenueRepository) GetEarningsByBook(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error) {
	return nil, 0, nil
}

func (m *mockAuthorRevenueRepository) CreateWithdrawalRequest(ctx context.Context, request *financeModel.WithdrawalRequest) error {
	if m.failCreate != nil {
		return m.failCreate
	}
	m.counter++
	request.ID = primitive.NewObjectID()
	m.withdrawals[request.ID.Hex()] = cloneAuthorRevenueWithdrawal(request)
	return nil
}

func (m *mockAuthorRevenueRepository) GetWithdrawalRequest(ctx context.Context, requestID primitive.ObjectID) (*financeModel.WithdrawalRequest, error) {
	request, ok := m.withdrawals[requestID.Hex()]
	if !ok {
		return nil, errors.New("not found")
	}
	return cloneAuthorRevenueWithdrawal(request), nil
}

func (m *mockAuthorRevenueRepository) ListWithdrawalRequests(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.WithdrawalRequest, int64, error) {
	return nil, 0, nil
}

func (m *mockAuthorRevenueRepository) UpdateWithdrawalRequest(ctx context.Context, requestID primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}

func (m *mockAuthorRevenueRepository) GetUserWithdrawalRequests(ctx context.Context, userID string, page, pageSize int) ([]*financeModel.WithdrawalRequest, int64, error) {
	return nil, 0, nil
}

func (m *mockAuthorRevenueRepository) CreateSettlement(ctx context.Context, settlement *financeModel.Settlement) error {
	return nil
}

func (m *mockAuthorRevenueRepository) GetSettlement(ctx context.Context, settlementID primitive.ObjectID) (*financeModel.Settlement, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAuthorRevenueRepository) ListSettlements(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.Settlement, int64, error) {
	return nil, 0, nil
}

func (m *mockAuthorRevenueRepository) UpdateSettlement(ctx context.Context, settlementID primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}

func (m *mockAuthorRevenueRepository) GetAuthorSettlements(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.Settlement, int64, error) {
	return nil, 0, nil
}

func (m *mockAuthorRevenueRepository) GetPendingSettlements(ctx context.Context) ([]*financeModel.Settlement, error) {
	return nil, nil
}

func (m *mockAuthorRevenueRepository) GetRevenueStatistics(ctx context.Context, authorID string, period string, limit int) ([]*financeModel.RevenueStatistics, error) {
	return nil, nil
}

func (m *mockAuthorRevenueRepository) CreateRevenueStatistics(ctx context.Context, statistics *financeModel.RevenueStatistics) error {
	return nil
}

func (m *mockAuthorRevenueRepository) UpdateRevenueStatistics(ctx context.Context, statisticsID primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}

func (m *mockAuthorRevenueRepository) GetRevenueDetails(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.RevenueDetail, int64, error) {
	return nil, 0, nil
}

func (m *mockAuthorRevenueRepository) GetRevenueDetailByBook(ctx context.Context, authorID string, bookID primitive.ObjectID) (*financeModel.RevenueDetail, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAuthorRevenueRepository) CreateRevenueDetail(ctx context.Context, detail *financeModel.RevenueDetail) error {
	return nil
}

func (m *mockAuthorRevenueRepository) UpdateRevenueDetail(ctx context.Context, detailID primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}

func (m *mockAuthorRevenueRepository) CreateTaxInfo(ctx context.Context, taxInfo *financeModel.TaxInfo) error {
	return nil
}

func (m *mockAuthorRevenueRepository) GetTaxInfo(ctx context.Context, userID string) (*financeModel.TaxInfo, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAuthorRevenueRepository) UpdateTaxInfo(ctx context.Context, userID string, updates map[string]interface{}) error {
	return nil
}

type mockWalletRepository struct {
	wallets             map[string]*financeModel.Wallet
	withdraws           map[string]*financeModel.WithdrawRequest
	failCreateWithdraw  error
	failUpdateBalance   error
	withdrawRequestSeed int
}

func newMockWalletRepository() *mockWalletRepository {
	return &mockWalletRepository{
		wallets:   make(map[string]*financeModel.Wallet),
		withdraws: make(map[string]*financeModel.WithdrawRequest),
	}
}

func (m *mockWalletRepository) CreateWallet(ctx context.Context, wallet *financeModel.Wallet) error {
	m.wallets[wallet.UserID] = cloneWallet(wallet)
	return nil
}

func (m *mockWalletRepository) GetWallet(ctx context.Context, userID string) (*financeModel.Wallet, error) {
	wallet, ok := m.wallets[userID]
	if !ok {
		return nil, errors.New("wallet not found")
	}
	return cloneWallet(wallet), nil
}

func (m *mockWalletRepository) UpdateWallet(ctx context.Context, userID string, updates map[string]interface{}) error {
	return nil
}

func (m *mockWalletRepository) UpdateBalance(ctx context.Context, userID string, amount int64) error {
	if m.failUpdateBalance != nil {
		return m.failUpdateBalance
	}
	wallet, ok := m.wallets[userID]
	if !ok {
		return errors.New("wallet not found")
	}
	wallet.Balance += types.Money(amount)
	return nil
}

func (m *mockWalletRepository) UpdateBalanceWithCheck(ctx context.Context, userID string, amount int64) error {
	if amount < 0 {
		wallet, ok := m.wallets[userID]
		if !ok {
			return errors.New("wallet not found")
		}
		if wallet.Balance < types.Money(-amount) {
			return errors.New("余额不足")
		}
	}
	return m.UpdateBalance(ctx, userID, amount)
}

func (m *mockWalletRepository) CreateTransaction(ctx context.Context, transaction *financeModel.Transaction) error {
	return nil
}

func (m *mockWalletRepository) GetTransaction(ctx context.Context, transactionID string) (*financeModel.Transaction, error) {
	return nil, errors.New("not implemented")
}

func (m *mockWalletRepository) ListTransactions(ctx context.Context, filter *financeRepo.TransactionFilter) ([]*financeModel.Transaction, error) {
	return nil, nil
}

func (m *mockWalletRepository) CountTransactions(ctx context.Context, filter *financeRepo.TransactionFilter) (int64, error) {
	return 0, nil
}

func (m *mockWalletRepository) CreateWithdrawRequest(ctx context.Context, request *financeModel.WithdrawRequest) error {
	if m.failCreateWithdraw != nil {
		return m.failCreateWithdraw
	}
	m.withdrawRequestSeed++
	request.ID = primitive.NewObjectID()
	m.withdraws[request.ID.Hex()] = cloneWalletWithdraw(request)
	return nil
}

func (m *mockWalletRepository) GetWithdrawRequest(ctx context.Context, withdrawID string) (*financeModel.WithdrawRequest, error) {
	request, ok := m.withdraws[withdrawID]
	if !ok {
		return nil, errors.New("withdraw not found")
	}
	return cloneWalletWithdraw(request), nil
}

func (m *mockWalletRepository) UpdateWithdrawRequest(ctx context.Context, withdrawID string, updates map[string]interface{}) error {
	return nil
}

func (m *mockWalletRepository) ListWithdrawRequests(ctx context.Context, filter *financeRepo.WithdrawFilter) ([]*financeModel.WithdrawRequest, error) {
	return nil, nil
}

func (m *mockWalletRepository) CountWithdrawRequests(ctx context.Context, filter *financeRepo.WithdrawFilter) (int64, error) {
	return 0, nil
}

func (m *mockWalletRepository) Health(ctx context.Context) error {
	return nil
}

func (m *mockWalletRepository) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type rollbackRunner struct {
	revenueRepo *mockAuthorRevenueRepository
	walletRepo  *mockWalletRepository
}

func (r rollbackRunner) Run(ctx context.Context, fn func(context.Context) error) error {
	revenueSnapshot := cloneAuthorRevenueWithdrawMap(r.revenueRepo.withdrawals)
	walletSnapshot := cloneWalletMap(r.walletRepo.wallets)
	walletWithdrawSnapshot := cloneWalletWithdrawMap(r.walletRepo.withdraws)

	if err := fn(ctx); err != nil {
		r.revenueRepo.withdrawals = revenueSnapshot
		r.walletRepo.wallets = walletSnapshot
		r.walletRepo.withdraws = walletWithdrawSnapshot
		return err
	}

	return nil
}

func TestAuthorRevenueCreateWithdrawalRequestSuccess(t *testing.T) {
	revenueRepo := newMockAuthorRevenueRepository()
	walletRepo := newMockWalletRepository()
	walletRepo.wallets["author-1"] = &financeModel.Wallet{ID: primitive.NewObjectID(), UserID: "author-1", Balance: types.NewMoneyFromYuan(100)}
	service := NewAuthorRevenueServiceWithDependencies(revenueRepo, walletRepo, rollbackRunner{
		revenueRepo: revenueRepo,
		walletRepo:  walletRepo,
	}).(*AuthorRevenueServiceImpl)

	request, err := service.CreateWithdrawalRequest(context.Background(), "author-1", 50, financeModel.WithdrawMethodAlipay, financeModel.WithdrawAccount{
		AccountType: "alipay",
		AccountName: "Author",
		AccountNo:   "author@example.com",
	})

	assert.NoError(t, err)
	assert.NotNil(t, request)
	assert.Len(t, revenueRepo.withdrawals, 1)
	assert.Len(t, walletRepo.withdraws, 1)
	assert.NotEmpty(t, request.TransactionID)
	assert.Equal(t, int64(types.NewMoneyFromYuan(50)), int64(walletRepo.wallets["author-1"].Balance))
}

func TestAuthorRevenueCreateWithdrawalRequestRollbackOnWalletFailure(t *testing.T) {
	revenueRepo := newMockAuthorRevenueRepository()
	walletRepo := newMockWalletRepository()
	walletRepo.wallets["author-1"] = &financeModel.Wallet{ID: primitive.NewObjectID(), UserID: "author-1", Balance: types.NewMoneyFromYuan(100)}
	walletRepo.failCreateWithdraw = errors.New("mock wallet withdraw failure")
	service := NewAuthorRevenueServiceWithDependencies(revenueRepo, walletRepo, rollbackRunner{
		revenueRepo: revenueRepo,
		walletRepo:  walletRepo,
	}).(*AuthorRevenueServiceImpl)

	request, err := service.CreateWithdrawalRequest(context.Background(), "author-1", 50, financeModel.WithdrawMethodAlipay, financeModel.WithdrawAccount{
		AccountType: "alipay",
		AccountName: "Author",
		AccountNo:   "author@example.com",
	})

	assert.Error(t, err)
	assert.Nil(t, request)
	assert.Empty(t, revenueRepo.withdrawals)
	assert.Empty(t, walletRepo.withdraws)
	assert.Equal(t, int64(types.NewMoneyFromYuan(100)), int64(walletRepo.wallets["author-1"].Balance))
}

func TestAuthorRevenueCreateWithdrawalRequestRollbackOnBalanceFailure(t *testing.T) {
	revenueRepo := newMockAuthorRevenueRepository()
	walletRepo := newMockWalletRepository()
	walletRepo.wallets["author-1"] = &financeModel.Wallet{ID: primitive.NewObjectID(), UserID: "author-1", Balance: types.NewMoneyFromYuan(100)}
	walletRepo.failUpdateBalance = errors.New("mock balance failure")
	service := NewAuthorRevenueServiceWithDependencies(revenueRepo, walletRepo, rollbackRunner{
		revenueRepo: revenueRepo,
		walletRepo:  walletRepo,
	}).(*AuthorRevenueServiceImpl)

	request, err := service.CreateWithdrawalRequest(context.Background(), "author-1", 50, financeModel.WithdrawMethodAlipay, financeModel.WithdrawAccount{
		AccountType: "alipay",
		AccountName: "Author",
		AccountNo:   "author@example.com",
	})

	assert.Error(t, err)
	assert.Nil(t, request)
	assert.Empty(t, revenueRepo.withdrawals)
	assert.Empty(t, walletRepo.withdraws)
	assert.Equal(t, int64(types.NewMoneyFromYuan(100)), int64(walletRepo.wallets["author-1"].Balance))
}

func cloneAuthorRevenueWithdrawal(request *financeModel.WithdrawalRequest) *financeModel.WithdrawalRequest {
	if request == nil {
		return nil
	}
	cloned := *request
	return &cloned
}

func cloneAuthorRevenueWithdrawMap(source map[string]*financeModel.WithdrawalRequest) map[string]*financeModel.WithdrawalRequest {
	cloned := make(map[string]*financeModel.WithdrawalRequest, len(source))
	for key, value := range source {
		cloned[key] = cloneAuthorRevenueWithdrawal(value)
	}
	return cloned
}

func cloneWallet(wallet *financeModel.Wallet) *financeModel.Wallet {
	if wallet == nil {
		return nil
	}
	cloned := *wallet
	return &cloned
}

func cloneWalletMap(source map[string]*financeModel.Wallet) map[string]*financeModel.Wallet {
	cloned := make(map[string]*financeModel.Wallet, len(source))
	for key, value := range source {
		cloned[key] = cloneWallet(value)
	}
	return cloned
}

func cloneWalletWithdraw(request *financeModel.WithdrawRequest) *financeModel.WithdrawRequest {
	if request == nil {
		return nil
	}
	cloned := *request
	return &cloned
}

func cloneWalletWithdrawMap(source map[string]*financeModel.WithdrawRequest) map[string]*financeModel.WithdrawRequest {
	cloned := make(map[string]*financeModel.WithdrawRequest, len(source))
	for key, value := range source {
		cloned[key] = cloneWalletWithdraw(value)
	}
	return cloned
}
