package core

import (
	"context"
	"log"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/global"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongoDB 初始化MongoDB连接
func InitMongoDB() error {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置MongoDB连接选项
	clientOptions := options.Client().ApplyURI(cfg.Database.MongoURI)

	// 应用高级配置选项
	clientOptions.SetConnectTimeout(cfg.Database.ConnectTimeout)
	clientOptions.SetMaxPoolSize(cfg.Database.MaxPoolSize)
	clientOptions.SetMinPoolSize(cfg.Database.MinPoolSize)
	clientOptions.SetMaxConnIdleTime(cfg.Database.MaxConnIdleTime)
	clientOptions.SetRetryWrites(cfg.Database.RetryWrites)
	clientOptions.SetRetryReads(cfg.Database.RetryReads)

	// 设置连接超时
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.ConnectTimeout)
	defer cancel()

	// 连接到MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		return err
	}

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		return err
	}

	// 设置全局变量
	global.MongoClient = client
	global.DB = client.Database(cfg.Database.DBName)

	log.Printf("Connected to MongoDB: %s, Database: %s", cfg.Database.MongoURI, cfg.Database.DBName)
	log.Printf("MongoDB connection pool configured with MaxSize: %d, MinSize: %d", cfg.Database.MaxPoolSize, cfg.Database.MinPoolSize)
	return nil
}

// CloseMongoDB 关闭MongoDB连接
func CloseMongoDB() {
	if global.MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := global.MongoClient.Disconnect(ctx)
		if err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		} else {
			log.Println("Disconnected from MongoDB")
		}
	}
}
