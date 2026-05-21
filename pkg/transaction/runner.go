package transaction

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

// Runner 提供领域无关的事务执行入口。
type Runner interface {
	Run(ctx context.Context, fn func(context.Context) error) error
}

type mongoRunner struct {
	client *mongo.Client
}

// NewMongoRunner 使用 MongoDB session 执行事务。
func NewMongoRunner(client *mongo.Client) Runner {
	return &mongoRunner{client: client}
}

func (r *mongoRunner) Run(ctx context.Context, fn func(context.Context) error) error {
	err := RunMongoTransaction(ctx, r.client, fn)
	if err != nil {
		return fmt.Errorf("事务执行失败: %w", err)
	}
	return nil
}
