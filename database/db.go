package database

import (
	"context"
	"log"
	"time"

	"Qingyu_backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client 全局 MongoDB 客户端
var Client *mongo.Client

// 保存当前配置
var currentConfig *config.Config

// ConnectDB 连接到 MongoDB
func ConnectDB(cfg *config.Config) error {
	var err error
	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建客户端选项
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	// 连接到MongoDB
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("连接MongoDB失败: ", err)
		return err
	}

	// 检查连接
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Println("MongoDB连接测试失败: ", err)
		return err
	}

	log.Println("MongoDB连接成功")

	// 保存配置以供后续使用
	currentConfig = cfg

	return nil
}

// GetCollection 获取一个集合
func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database(currentConfig.DBName).Collection(collectionName)
}

// DisconnectDB 断开与MongoDB的连接
func DisconnectDB() error {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := Client.Disconnect(ctx)
		if err != nil {
			log.Println("断开MongoDB连接失败:", err)
			return err
		}
		log.Println("MongoDB连接已关闭")
	}
	return nil
}
