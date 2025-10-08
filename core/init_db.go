package core

import (
	"context"
	"fmt"

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

	// 检查主数据库配置
	if cfg.Primary.Type != config.DatabaseTypeMongoDB || cfg.Primary.MongoDB == nil {
		return fmt.Errorf("MongoDB configuration is missing or invalid")
	}

	mongoCfg := cfg.Primary.MongoDB

	// 创建MongoDB客户端配置
	clientOptions := options.Client().
		ApplyURI(mongoCfg.URI).
		SetConnectTimeout(mongoCfg.ConnectTimeout).
		SetMaxPoolSize(mongoCfg.MaxPoolSize).
		SetMinPoolSize(mongoCfg.MinPoolSize).
		SetServerSelectionTimeout(mongoCfg.ServerTimeout)

	// 连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), mongoCfg.ConnectTimeout)
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
	global.DB = client.Database(mongoCfg.Database)

	fmt.Printf("Successfully connected to MongoDB: %s/%s\n", mongoCfg.URI, mongoCfg.Database)
	return nil
}
