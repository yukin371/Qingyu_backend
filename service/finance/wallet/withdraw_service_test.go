package wallet

import (
	"context"
	"errors"
	"testing"

	financeModel "Qingyu_backend/models/finance"
	"Qingyu_backend/models/shared/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithdrawServiceCreateRequestRollbackOnBalanceFailure(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user_a"] = &financeModel.Wallet{ID: "wallet_a", UserID: "user_a", Balance: types.Money(5000)}
	repo.SetUpdateBalanceError("user_a", errors.New("mock freeze failure"))

	service := NewWithdrawService(repo)

	request, err := service.CreateWithdrawRequest(context.Background(), "user_a", "user_a", 2000, "alipay", "acc-1")
	require.Error(t, err)
	assert.Nil(t, request)
	assert.Equal(t, types.Money(5000), repo.wallets["user_a"].Balance)
	assert.Empty(t, repo.withdrawRequests)
}

func TestWithdrawServiceApproveRollbackOnTransactionFailure(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user_a"] = &financeModel.Wallet{ID: "wallet_a", UserID: "user_a", Balance: types.Money(5000)}
	repo.withdrawRequests["wd_1"] = &financeModel.WithdrawRequest{
		ID:     "wd_1",
		UserID: "user_a",
		Amount: types.Money(2000),
		Status: "pending",
	}
	repo.SetCreateTransactionError("withdraw", errors.New("mock transaction failure"))

	service := NewWithdrawService(repo)

	err := service.ApproveWithdraw(context.Background(), "wd_1", "admin", "ok")
	require.Error(t, err)
	assert.Equal(t, "pending", repo.withdrawRequests["wd_1"].Status)
	assert.Empty(t, repo.transactions)
}

func TestWithdrawServiceRejectRollbackOnStatusUpdateFailure(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user_a"] = &financeModel.Wallet{ID: "wallet_a", UserID: "user_a", Balance: types.Money(3000)}
	repo.withdrawRequests["wd_1"] = &financeModel.WithdrawRequest{
		ID:     "wd_1",
		UserID: "user_a",
		Amount: types.Money(1000),
		Status: "pending",
	}
	repo.SetUpdateWithdrawError("wd_1", errors.New("mock status update failure"))

	service := NewWithdrawService(repo)

	err := service.RejectWithdraw(context.Background(), "wd_1", "admin", "bad account")
	require.Error(t, err)
	assert.Equal(t, types.Money(3000), repo.wallets["user_a"].Balance)
	assert.Equal(t, "pending", repo.withdrawRequests["wd_1"].Status)
}
