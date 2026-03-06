package wallet

import (
	"context"

	pkgtransaction "Qingyu_backend/pkg/transaction"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// TransactionRunner 定义钱包域的事务入口，屏蔽底层仓储事务实现细节。
type TransactionRunner interface {
	Run(ctx context.Context, fn func(context.Context) error) error
}

type repositoryTransactionRunner struct {
	walletRepo sharedRepo.WalletRepository
}

type genericTransactionRunner struct {
	runner pkgtransaction.Runner
}

// NewRepositoryTransactionRunner 使用钱包仓储提供事务执行能力。
func NewRepositoryTransactionRunner(walletRepo sharedRepo.WalletRepository) TransactionRunner {
	return &repositoryTransactionRunner{walletRepo: walletRepo}
}

// NewGenericTransactionRunner 将通用事务执行器适配为钱包域事务入口。
func NewGenericTransactionRunner(runner pkgtransaction.Runner) TransactionRunner {
	return &genericTransactionRunner{runner: runner}
}

func (r *repositoryTransactionRunner) Run(ctx context.Context, fn func(context.Context) error) error {
	return r.walletRepo.RunInTransaction(ctx, fn)
}

func (r *genericTransactionRunner) Run(ctx context.Context, fn func(context.Context) error) error {
	return r.runner.Run(ctx, fn)
}

func runWalletTransaction(ctx context.Context, runner TransactionRunner, fn func(context.Context) error) error {
	return runner.Run(ctx, fn)
}
