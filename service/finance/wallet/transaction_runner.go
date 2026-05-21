package wallet

import (
	"context"
	"strings"

	pkgtransaction "Qingyu_backend/pkg/transaction"
	financeRepo "Qingyu_backend/repository/interfaces/finance"
)

// TransactionRunner 定义钱包域的事务入口，屏蔽底层仓储事务实现细节。
type TransactionRunner interface {
	Run(ctx context.Context, fn func(context.Context) error) error
}

type repositoryTransactionRunner struct {
	walletRepo financeRepo.WalletRepository
}

type genericTransactionRunner struct {
	runner pkgtransaction.Runner
}

// NewRepositoryTransactionRunner 使用钱包仓储提供事务执行能力。
func NewRepositoryTransactionRunner(walletRepo financeRepo.WalletRepository) TransactionRunner {
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
	if runner == nil {
		return fn(ctx)
	}

	err := runner.Run(ctx, fn)
	if err == nil {
		return nil
	}

	// 本地单机 Mongo 默认不支持事务，允许降级到顺序执行以完成开发联调。
	if isTransactionUnsupported(err) {
		return fn(ctx)
	}

	return err
}

func isTransactionUnsupported(err error) bool {
	if err == nil {
		return false
	}

	message := err.Error()
	return strings.Contains(message, "Transaction numbers are only allowed on a replica set member or mongos") ||
		strings.Contains(message, "transactions are not supported")
}
