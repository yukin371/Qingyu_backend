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
	if r.client == nil {
		return fmt.Errorf("mongo client is nil")
	}

	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("启动事务失败: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		if err := fn(sessCtx); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return fmt.Errorf("事务执行失败: %w", err)
	}
	return nil
}
