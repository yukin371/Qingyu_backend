package finance

import (
	"context"
	"errors"
	"testing"
	"time"

	financeModel "Qingyu_backend/models/finance"
	"Qingyu_backend/models/shared/types"
	financeRepo "Qingyu_backend/repository/interfaces/finance"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type membershipStateRepository struct {
	plan                 *financeModel.MembershipPlan
	membershipByUser     map[string]*financeModel.UserMembership
	failCreateMembership error
	failUpdateMembership error
}

func newMembershipStateRepository(plan *financeModel.MembershipPlan) *membershipStateRepository {
	return &membershipStateRepository{
		plan:             plan,
		membershipByUser: make(map[string]*financeModel.UserMembership),
	}
}

func (m *membershipStateRepository) CreatePlan(ctx context.Context, plan *financeModel.MembershipPlan) error {
	return nil
}
func (m *membershipStateRepository) GetPlan(ctx context.Context, planID primitive.ObjectID) (*financeModel.MembershipPlan, error) {
	if m.plan != nil && m.plan.ID == planID {
		return cloneUserPlan(m.plan), nil
	}
	return nil, errors.New("plan not found")
}
func (m *membershipStateRepository) GetPlanByType(ctx context.Context, planType string) (*financeModel.MembershipPlan, error) {
	return nil, errors.New("not implemented")
}
func (m *membershipStateRepository) ListPlans(ctx context.Context, enabledOnly bool) ([]*financeModel.MembershipPlan, error) {
	return nil, nil
}
func (m *membershipStateRepository) UpdatePlan(ctx context.Context, planID primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}
func (m *membershipStateRepository) DeletePlan(ctx context.Context, planID primitive.ObjectID) error {
	return nil
}
func (m *membershipStateRepository) CreateMembership(ctx context.Context, membership *financeModel.UserMembership) error {
	if m.failCreateMembership != nil {
		return m.failCreateMembership
	}
	if membership.ID.IsZero() {
		membership.ID = primitive.NewObjectID()
	}
	m.membershipByUser[membership.UserID] = cloneUserMembership(membership)
	return nil
}
func (m *membershipStateRepository) GetMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error) {
	membership, ok := m.membershipByUser[userID]
	if !ok {
		return nil, errors.New("membership not found")
	}
	return cloneUserMembership(membership), nil
}
func (m *membershipStateRepository) GetMembershipByID(ctx context.Context, membershipID primitive.ObjectID) (*financeModel.UserMembership, error) {
	for _, membership := range m.membershipByUser {
		if membership.ID == membershipID {
			return cloneUserMembership(membership), nil
		}
	}
	return nil, errors.New("membership not found")
}
func (m *membershipStateRepository) UpdateMembership(ctx context.Context, membershipID primitive.ObjectID, updates map[string]interface{}) error {
	if m.failUpdateMembership != nil {
		return m.failUpdateMembership
	}
	for userID, membership := range m.membershipByUser {
		if membership.ID != membershipID {
			continue
		}
		if endTime, ok := updates["end_time"].(time.Time); ok {
			membership.EndTime = endTime
		}
		if status, ok := updates["status"].(string); ok {
			membership.Status = status
		}
		if paymentID, ok := updates["payment_id"].(primitive.ObjectID); ok {
			membership.PaymentID = paymentID
		}
		m.membershipByUser[userID] = cloneUserMembership(membership)
		return nil
	}
	return errors.New("membership not found")
}
func (m *membershipStateRepository) DeleteMembership(ctx context.Context, membershipID primitive.ObjectID) error {
	return nil
}
func (m *membershipStateRepository) ListMemberships(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.UserMembership, error) {
	return nil, nil
}
func (m *membershipStateRepository) CreateMembershipCard(ctx context.Context, card *financeModel.MembershipCard) error {
	return nil
}
func (m *membershipStateRepository) GetMembershipCard(ctx context.Context, cardID primitive.ObjectID) (*financeModel.MembershipCard, error) {
	return nil, errors.New("not implemented")
}
func (m *membershipStateRepository) GetMembershipCardByCode(ctx context.Context, code string) (*financeModel.MembershipCard, error) {
	return nil, errors.New("not implemented")
}
func (m *membershipStateRepository) UpdateMembershipCard(ctx context.Context, cardID primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}
func (m *membershipStateRepository) ListMembershipCards(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.MembershipCard, error) {
	return nil, nil
}
func (m *membershipStateRepository) CountMembershipCards(ctx context.Context, filter map[string]interface{}) (int64, error) {
	return 0, nil
}
func (m *membershipStateRepository) BatchCreateMembershipCards(ctx context.Context, cards []*financeModel.MembershipCard) error {
	return nil
}
func (m *membershipStateRepository) CreateBenefit(ctx context.Context, benefit *financeModel.MembershipBenefit) error {
	return nil
}
func (m *membershipStateRepository) GetBenefit(ctx context.Context, benefitID primitive.ObjectID) (*financeModel.MembershipBenefit, error) {
	return nil, errors.New("not implemented")
}
func (m *membershipStateRepository) GetBenefitByCode(ctx context.Context, code string) (*financeModel.MembershipBenefit, error) {
	return nil, errors.New("not implemented")
}
func (m *membershipStateRepository) ListBenefits(ctx context.Context, level string, enabledOnly bool) ([]*financeModel.MembershipBenefit, error) {
	return nil, nil
}
func (m *membershipStateRepository) UpdateBenefit(ctx context.Context, benefitID primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}
func (m *membershipStateRepository) DeleteBenefit(ctx context.Context, benefitID primitive.ObjectID) error {
	return nil
}
func (m *membershipStateRepository) CreateUsage(ctx context.Context, usage *financeModel.MembershipUsage) error {
	return nil
}
func (m *membershipStateRepository) GetUsage(ctx context.Context, userID string, benefitCode string) (*financeModel.MembershipUsage, error) {
	return nil, errors.New("not implemented")
}
func (m *membershipStateRepository) UpdateUsage(ctx context.Context, usageID primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}
func (m *membershipStateRepository) ListUsages(ctx context.Context, userID string) ([]*financeModel.MembershipUsage, error) {
	return nil, nil
}

type membershipWalletRepository struct {
	wallets           map[string]*financeModel.Wallet
	transactions      map[string]*financeModel.Transaction
	failUpdateBalance error
	failCreateTx      error
}

func newMembershipWalletRepository() *membershipWalletRepository {
	return &membershipWalletRepository{
		wallets:      make(map[string]*financeModel.Wallet),
		transactions: make(map[string]*financeModel.Transaction),
	}
}

func (m *membershipWalletRepository) CreateWallet(ctx context.Context, wallet *financeModel.Wallet) error {
	m.wallets[wallet.UserID] = cloneWallet(wallet)
	return nil
}
func (m *membershipWalletRepository) GetWallet(ctx context.Context, userID string) (*financeModel.Wallet, error) {
	wallet, ok := m.wallets[userID]
	if !ok {
		return nil, errors.New("wallet not found")
	}
	return cloneWallet(wallet), nil
}
func (m *membershipWalletRepository) UpdateWallet(ctx context.Context, userID string, updates map[string]interface{}) error {
	return nil
}
func (m *membershipWalletRepository) UpdateBalance(ctx context.Context, userID string, amount int64) error {
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
func (m *membershipWalletRepository) CreateTransaction(ctx context.Context, transaction *financeModel.Transaction) error {
	if m.failCreateTx != nil {
		return m.failCreateTx
	}
	if transaction.ID.IsZero() {
		transaction.ID = primitive.NewObjectID()
	}
	m.transactions[transaction.ID.Hex()] = cloneFinanceTransaction(transaction)
	return nil
}
func (m *membershipWalletRepository) GetTransaction(ctx context.Context, transactionID string) (*financeModel.Transaction, error) {
	return nil, errors.New("not implemented")
}
func (m *membershipWalletRepository) ListTransactions(ctx context.Context, filter *financeRepo.TransactionFilter) ([]*financeModel.Transaction, error) {
	return nil, nil
}
func (m *membershipWalletRepository) CountTransactions(ctx context.Context, filter *financeRepo.TransactionFilter) (int64, error) {
	return 0, nil
}
func (m *membershipWalletRepository) CreateWithdrawRequest(ctx context.Context, request *financeModel.WithdrawRequest) error {
	return nil
}
func (m *membershipWalletRepository) GetWithdrawRequest(ctx context.Context, withdrawID string) (*financeModel.WithdrawRequest, error) {
	return nil, errors.New("not implemented")
}
func (m *membershipWalletRepository) UpdateWithdrawRequest(ctx context.Context, withdrawID string, updates map[string]interface{}) error {
	return nil
}
func (m *membershipWalletRepository) ListWithdrawRequests(ctx context.Context, filter *financeRepo.WithdrawFilter) ([]*financeModel.WithdrawRequest, error) {
	return nil, nil
}
func (m *membershipWalletRepository) CountWithdrawRequests(ctx context.Context, filter *financeRepo.WithdrawFilter) (int64, error) {
	return 0, nil
}
func (m *membershipWalletRepository) Health(ctx context.Context) error { return nil }
func (m *membershipWalletRepository) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type membershipRollbackRunner struct {
	membershipRepo *membershipStateRepository
	walletRepo     *membershipWalletRepository
}

func (r membershipRollbackRunner) Run(ctx context.Context, fn func(context.Context) error) error {
	membershipSnapshot := cloneMembershipMap(r.membershipRepo.membershipByUser)
	walletSnapshot := cloneWalletMap(r.walletRepo.wallets)
	transactionSnapshot := cloneFinanceTransactionMap(r.walletRepo.transactions)
	if err := fn(ctx); err != nil {
		r.membershipRepo.membershipByUser = membershipSnapshot
		r.walletRepo.wallets = walletSnapshot
		r.walletRepo.transactions = transactionSnapshot
		return err
	}
	return nil
}

func TestMembershipSubscribeTransactionSuccess(t *testing.T) {
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	repo := newMembershipStateRepository(plan)
	walletRepo := newMembershipWalletRepository()
	walletRepo.wallets["user123"] = &financeModel.Wallet{ID: primitive.NewObjectID(), UserID: "user123", Balance: types.NewMoneyFromYuan(100)}

	service := NewMembershipServiceWithDependencies(repo, walletRepo, membershipRollbackRunner{
		membershipRepo: repo,
		walletRepo:     walletRepo,
	}).(*MembershipServiceImpl)

	membership, err := service.Subscribe(context.Background(), "user123", plan.ID.Hex(), "wallet")
	assert.NoError(t, err)
	assert.NotNil(t, membership)
	assert.Equal(t, financeModel.MembershipStatusActive, membership.Status)
	assert.Len(t, repo.membershipByUser, 1)
	assert.Len(t, walletRepo.transactions, 1)
	assert.Equal(t, int64(types.NewMoneyFromYuan(80.1)), int64(walletRepo.wallets["user123"].Balance))
}

func TestMembershipSubscribeRollbackOnWalletTransactionFailure(t *testing.T) {
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	repo := newMembershipStateRepository(plan)
	walletRepo := newMembershipWalletRepository()
	walletRepo.wallets["user123"] = &financeModel.Wallet{ID: primitive.NewObjectID(), UserID: "user123", Balance: types.NewMoneyFromYuan(100)}
	walletRepo.failCreateTx = errors.New("mock tx failure")

	service := NewMembershipServiceWithDependencies(repo, walletRepo, membershipRollbackRunner{
		membershipRepo: repo,
		walletRepo:     walletRepo,
	}).(*MembershipServiceImpl)

	membership, err := service.Subscribe(context.Background(), "user123", plan.ID.Hex(), "wallet")
	assert.Error(t, err)
	assert.Nil(t, membership)
	assert.Empty(t, repo.membershipByUser)
	assert.Empty(t, walletRepo.transactions)
	assert.Equal(t, int64(types.NewMoneyFromYuan(100)), int64(walletRepo.wallets["user123"].Balance))
}

func TestMembershipRenewRollbackOnBalanceFailure(t *testing.T) {
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	repo := newMembershipStateRepository(plan)
	now := time.Now()
	repo.membershipByUser["user123"] = &financeModel.UserMembership{
		ID:        primitive.NewObjectID(),
		UserID:    "user123",
		PlanID:    plan.ID,
		PlanName:  plan.Name,
		PlanType:  plan.Type,
		Level:     financeModel.MembershipLevelVIPMonthly,
		StartTime: now.Add(-10 * 24 * time.Hour),
		EndTime:   now.Add(5 * 24 * time.Hour),
		Status:    financeModel.MembershipStatusActive,
	}

	walletRepo := newMembershipWalletRepository()
	walletRepo.wallets["user123"] = &financeModel.Wallet{ID: primitive.NewObjectID(), UserID: "user123", Balance: types.NewMoneyFromYuan(100)}
	walletRepo.failUpdateBalance = errors.New("mock balance failure")

	service := NewMembershipServiceWithDependencies(repo, walletRepo, membershipRollbackRunner{
		membershipRepo: repo,
		walletRepo:     walletRepo,
	}).(*MembershipServiceImpl)

	membership, err := service.RenewMembership(context.Background(), "user123")
	assert.Error(t, err)
	assert.Nil(t, membership)
	assert.Equal(t, int64(types.NewMoneyFromYuan(100)), int64(walletRepo.wallets["user123"].Balance))
	assert.True(t, repo.membershipByUser["user123"].PaymentID.IsZero())
}

func cloneUserPlan(plan *financeModel.MembershipPlan) *financeModel.MembershipPlan {
	if plan == nil {
		return nil
	}
	cloned := *plan
	return &cloned
}

func cloneUserMembership(membership *financeModel.UserMembership) *financeModel.UserMembership {
	if membership == nil {
		return nil
	}
	cloned := *membership
	return &cloned
}

func cloneMembershipMap(source map[string]*financeModel.UserMembership) map[string]*financeModel.UserMembership {
	cloned := make(map[string]*financeModel.UserMembership, len(source))
	for key, value := range source {
		cloned[key] = cloneUserMembership(value)
	}
	return cloned
}

func cloneFinanceTransaction(transaction *financeModel.Transaction) *financeModel.Transaction {
	if transaction == nil {
		return nil
	}
	cloned := *transaction
	return &cloned
}

func cloneFinanceTransactionMap(source map[string]*financeModel.Transaction) map[string]*financeModel.Transaction {
	cloned := make(map[string]*financeModel.Transaction, len(source))
	for key, value := range source {
		cloned[key] = cloneFinanceTransaction(value)
	}
	return cloned
}
