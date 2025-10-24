package global

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// 全局变量
// Deprecated: 这些全局变量已废弃，将在v2.0版本中删除
// 请使用 ServiceContainer.GetMongoDB() 和 ServiceContainer.GetMongoClient() 代替
//
// 迁移指南：
//   旧方式: global.DB.Collection("users")
//   新方式: serviceContainer.GetMongoDB().Collection("users")
//
// 或者在服务中通过依赖注入：
//   func NewYourService(db *mongo.Database) *YourService {
//       return &YourService{db: db}
//   }
var (
	// MongoDB 客户端
	// Deprecated: 使用 ServiceContainer.GetMongoClient() 代替
	MongoClient *mongo.Client

	// MongoDB 数据库
	// Deprecated: 使用 ServiceContainer.GetMongoDB() 代替
	DB *mongo.Database
)
