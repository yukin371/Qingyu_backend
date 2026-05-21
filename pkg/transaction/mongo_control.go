package transaction

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

const disableMongoTransactionsEnv = "QINGYU_MONGODB_DISABLE_TRANSACTIONS"

// MongoTransactionsDisabled 返回当前进程是否显式关闭 MongoDB 事务。
func MongoTransactionsDisabled() bool {
	if value, ok := lookupBoolValue(os.Getenv(disableMongoTransactionsEnv)); ok {
		return value
	}

	if value, ok := lookupBoolInDotEnv(filepath.Join(".", ".env"), disableMongoTransactionsEnv); ok {
		return value
	}

	return false
}

// RunMongoTransaction 在启用事务时使用 WithTransaction，关闭事务时保留 session/context 但顺序执行。
func RunMongoTransaction(ctx context.Context, client *mongo.Client, fn func(context.Context) error) error {
	return RunMongoTransactionWithSession(ctx, client, func(sessCtx mongo.SessionContext) error {
		return fn(sessCtx)
	})
}

// RunMongoTransactionWithSession 为需要 mongo.SessionContext 的场景提供统一执行入口。
func RunMongoTransactionWithSession(ctx context.Context, client *mongo.Client, fn func(mongo.SessionContext) error) error {
	return RunMongoSession(ctx, client, func(_ mongo.Session, sessCtx mongo.SessionContext) error {
		return fn(sessCtx)
	})
}

// RunMongoSession 统一管理 session 生命周期，并根据开关决定是否开启事务。
func RunMongoSession(ctx context.Context, client *mongo.Client, fn func(mongo.Session, mongo.SessionContext) error) error {
	if client == nil {
		return fmt.Errorf("mongo client is nil")
	}

	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	if MongoTransactionsDisabled() {
		return fn(session, mongo.NewSessionContext(ctx, session))
	}

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		if err := fn(session, sessCtx); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

// RunMongoSessionWithResult 为需要返回结果对象的场景提供统一执行入口。
func RunMongoSessionWithResult(ctx context.Context, client *mongo.Client, fn func(mongo.Session, mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	if client == nil {
		return nil, fmt.Errorf("mongo client is nil")
	}

	session, err := client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	if MongoTransactionsDisabled() {
		return fn(session, mongo.NewSessionContext(ctx, session))
	}

	return session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return fn(session, sessCtx)
	})
}

func lookupBoolInDotEnv(path, key string) (bool, bool) {
	file, err := os.Open(path)
	if err != nil {
		return false, false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		if strings.TrimSpace(parts[0]) != key {
			continue
		}

		return lookupBoolValue(parts[1])
	}

	return false, false
}

func lookupBoolValue(raw string) (bool, bool) {
	value := strings.TrimSpace(strings.ToLower(raw))
	switch value {
	case "1", "true", "yes", "on":
		return true, true
	case "0", "false", "no", "off":
		return false, true
	default:
		return false, false
	}
}
