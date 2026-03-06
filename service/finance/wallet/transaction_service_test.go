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

func TestTransactionServiceTransferSuccess(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user_from"] = &financeModel.Wallet{ID: "wallet_from", UserID: "user_from", Balance: types.Money(1000)}
	repo.wallets["user_to"] = &financeModel.Wallet{ID: "wallet_to", UserID: "user_to", Balance: types.Money(200)}

	service := NewTransactionService(repo)

	err := service.Transfer(context.Background(), "user_from", "user_to", 300, "gift")
	require.NoError(t, err)

	assert.Equal(t, types.Money(700), repo.wallets["user_from"].Balance)
	assert.Equal(t, types.Money(500), repo.wallets["user_to"].Balance)
	assert.Len(t, repo.transactions, 2)
}

func TestTransactionServiceTransferRollbackOnTargetBalanceFailure(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user_from"] = &financeModel.Wallet{ID: "wallet_from", UserID: "user_from", Balance: types.Money(1000)}
	repo.wallets["user_to"] = &financeModel.Wallet{ID: "wallet_to", UserID: "user_to", Balance: types.Money(200)}
	repo.SetUpdateBalanceError("user_to", errors.New("mock target balance failure"))

	service := NewTransactionService(repo)

	err := service.Transfer(context.Background(), "user_from", "user_to", 300, "gift")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "更新目标钱包余额失败")

	assert.Equal(t, types.Money(1000), repo.wallets["user_from"].Balance)
	assert.Equal(t, types.Money(200), repo.wallets["user_to"].Balance)
	assert.Empty(t, repo.transactions)
}

func TestTransactionServiceRechargeRollbackOnBalanceFailure(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user_a"] = &financeModel.Wallet{ID: "wallet_a", UserID: "user_a", Balance: types.Money(100)}
	repo.SetUpdateBalanceError("user_a", errors.New("mock balance failure"))

	service := NewTransactionService(repo)

	tx, err := service.Recharge(context.Background(), "user_a", 50, "alipay", "ORDER001")
	require.Error(t, err)
	assert.Nil(t, tx)

	assert.Equal(t, types.Money(100), repo.wallets["user_a"].Balance)
	assert.Empty(t, repo.transactions)
}
