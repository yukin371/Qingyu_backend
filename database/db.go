package database

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "Qingyu_backend/config"
)

var Client *mongo.Client

// ConnectDB 连接到 MongoDB
func ConnectDB(cfg *config.Config) {
	var err error
	// 创建上下文
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

	// 创建客户端选项
    clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	// 连接到MongoDB
    Client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal("连接MongoDB失败: %v", err)
    }

    // 检查连接
    err = Client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("MongoDB连接测试失败: %v", err)
    }

    log.Println("MongoDB连接成功")
}

// GetCollection 获取一个集合
func GetCollection(collectionName string) *mongo.Collection {
    return Client.Database(config.LoadConfig().DBName).Collection(collectionName)
}