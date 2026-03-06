package wallet

import (
	"context"

	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

func runWalletTransaction(ctx context.Context, walletRepo sharedRepo.WalletRepository, fn func(context.Context) error) error {
	return walletRepo.RunInTransaction(ctx, fn)
}
