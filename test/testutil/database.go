package testutil

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"

	"Qingyu_backend/config"
	"Qingyu_backend/service"
	"Qingyu_backend/service/container"
)

// SetupTestDB 设置测试数据库
// 返回数据库实例和清理函数
func SetupTestDB(t *testing.T) (*mongo.Database, func()) {
	t.Helper()

	// 加载配置 - 使用相对于项目根目录的路径
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		// 如果失败，尝试使用相对路径
		cfg, err = config.LoadConfig("../../config/config.yaml")
		if err != nil {
			cfg, err = config.LoadConfig("../../../config/config.yaml")
			if err != nil {
				t.Fatalf("加载配置失败: %v", err)
			}
		}
	}

	// 初始化全局配置
	config.GlobalConfig = cfg

	// 创建服务容器并初始化
	c := container.NewServiceContainer()
	err = c.Initialize(context.Background())
	if err != nil {
		t.Fatalf("初始化ServiceContainer失败: %v", err)
	}

	// 获取MongoDB数据库
	db := c.GetMongoDB()
	if db == nil {
		t.Fatal("获取MongoDB数据库失败")
	}

	// 清理函数
	cleanup := func() {
		// 清理测试集合
		ctx := context.Background()

		// 推荐系统测试集合
		_ = db.Collection("user_behaviors").Drop(ctx)
		_ = db.Collection("user_profiles").Drop(ctx)
		_ = db.Collection("item_features").Drop(ctx)
		_ = db.Collection("book_statistics").Drop(ctx)
		_ = db.Collection("books").Drop(ctx)

		// Writing相关测试集合
		_ = db.Collection("projects").Drop(ctx)
		_ = db.Collection("documents").Drop(ctx)
		_ = db.Collection("document_contents").Drop(ctx)

		// Shared相关测试集合
		_ = db.Collection("wallets").Drop(ctx)
		_ = db.Collection("transactions").Drop(ctx)
		_ = db.Collection("withdraw_requests").Drop(ctx)

		// Reading相关测试集合
		_ = db.Collection("reading_settings").Drop(ctx)
		_ = db.Collection("reading_progress").Drop(ctx)
		_ = db.Collection("annotations").Drop(ctx)
		_ = db.Collection("chapters").Drop(ctx)

		// Messaging相关测试集合
		_ = db.Collection("announcements").Drop(ctx)

		// 关闭服务容器
		_ = c.Close(ctx)
	}

	return db, cleanup
}

// SetupTestContainer 设置完整的测试服务容器
// 返回服务容器和清理函数
func SetupTestContainer(t *testing.T) (*container.ServiceContainer, func()) {
	t.Helper()

	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		cfg, err = config.LoadConfig("../../config/config.yaml")
		if err != nil {
			cfg, err = config.LoadConfig("../../../config/config.yaml")
			if err != nil {
				t.Fatalf("加载配置失败: %v", err)
			}
		}
	}

	// 初始化全局配置
	config.GlobalConfig = cfg

	// 初始化服务容器
	err = service.InitializeServices()
	if err != nil {
		t.Fatalf("初始化服务失败: %v", err)
	}

	c := service.GetServiceContainer()
	if c == nil {
		t.Fatal("获取ServiceContainer失败")
	}

	cleanup := func() {
		ctx := context.Background()
		db := c.GetMongoDB()

		// 清理测试集合
		_ = db.Collection("user_behaviors").Drop(ctx)
		_ = db.Collection("user_profiles").Drop(ctx)
		_ = db.Collection("item_features").Drop(ctx)
		_ = db.Collection("book_statistics").Drop(ctx)
		_ = db.Collection("books").Drop(ctx)
		_ = db.Collection("projects").Drop(ctx)
		_ = db.Collection("documents").Drop(ctx)
		_ = db.Collection("document_contents").Drop(ctx)
		_ = db.Collection("wallets").Drop(ctx)
		_ = db.Collection("transactions").Drop(ctx)
		_ = db.Collection("withdraw_requests").Drop(ctx)
		_ = db.Collection("reading_settings").Drop(ctx)
		_ = db.Collection("reading_progress").Drop(ctx)
		_ = db.Collection("annotations").Drop(ctx)
		_ = db.Collection("chapters").Drop(ctx)
		_ = db.Collection("announcements").Drop(ctx)

		// 关闭服务容器
		_ = service.CloseServices(ctx)
	}

	return c, cleanup
}
