package core

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/global"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitDB 初始化数据库连接
func InitDB() error {
	cfg := config.GlobalConfig.Database
	if cfg == nil {
		return fmt.Errorf("database configuration is missing")
	}

	// 创建MongoDB客户端配置
	clientOptions := options.Client().
		ApplyURI(cfg.MongoURI).
		SetConnectTimeout(cfg.ConnectTimeout).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxConnIdleTime(cfg.MaxConnIdleTime).
		SetRetryWrites(cfg.RetryWrites).
		SetRetryReads(cfg.RetryReads)

	// 连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 验证连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// 设置全局数据库实例
	global.DB = client.Database(cfg.DBName)

	fmt.Printf("Successfully connected to MongoDB: %s/%s\n", cfg.MongoURI, cfg.DBName)
	return nil
}
