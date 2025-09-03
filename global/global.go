package global

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// 全局变量
var (
	// MongoDB 客户端
	MongoClient *mongo.Client
	
	// MongoDB 数据库
	DB *mongo.Database
)