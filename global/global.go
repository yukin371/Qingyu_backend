package global

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// 全局变量
//
// Deprecated: 这些全局变量已废弃，将在v2.0版本中删除
//
// ⚠️ 重要说明（ARCH-002重构后）：
//
// 1. 生产代码已完全迁移到依赖注入模式，使用 ServiceContainer
// 2. 测试代码仍在使用这些变量，需要逐步迁移
// 3. 新代码禁止使用这些全局变量
//
// 迁移指南：
//
//	旧方式（已废弃）:
//	global.DB.Collection("users")
//	global.MongoClient
//
//	新方式（推荐）:
//	service.ServiceManager.GetMongoDB().Collection("users")
//	service.ServiceManager.GetMongoClient()
//
//	或者在服务中通过依赖注入:
//	func NewYourService(db *mongo.Database) *YourService {
//	    return &YourService{db: db}
//	}
//
// 需要重构的测试文件（按优先级）：
// - P0: test/integration/helpers.go ✅ 已完成
// - P0: test/e2e/data/helper.go ✅ 已完成
// - P1: test/e2e/data/*.go (consistency_validator.go, factory.go)
// - P1: test/e2e/framework/*.go
// - P2: test/integration/*_test.go (各个集成测试文件)
var (
	// MongoDB 客户端
	// Deprecated: 使用 ServiceContainer.GetMongoClient() 代替
	MongoClient *mongo.Client

	// MongoDB 数据库
	// Deprecated: 使用 ServiceContainer.GetMongoDB() 代替
	DB *mongo.Database
)
