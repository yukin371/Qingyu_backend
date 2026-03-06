package wallet

import (
	"context"

	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// TransactionRunner 定义钱包域的事务入口，屏蔽底层仓储事务实现细节。
type TransactionRunner interface {
	Run(ctx context.Context, fn func(context.Context) error) error
}

type repositoryTransactionRunner struct {
	walletRepo sharedRepo.WalletRepository
}

// NewRepositoryTransactionRunner 使用钱包仓储提供事务执行能力。
func NewRepositoryTransactionRunner(walletRepo sharedRepo.WalletRepository) TransactionRunner {
	return &repositoryTransactionRunner{walletRepo: walletRepo}
}

func (r *repositoryTransactionRunner) Run(ctx context.Context, fn func(context.Context) error) error {
	return r.walletRepo.RunInTransaction(ctx, fn)
}

func runWalletTransaction(ctx context.Context, runner TransactionRunner, fn func(context.Context) error) error {
	return runner.Run(ctx, fn)
}
