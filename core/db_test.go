package core

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/service"
)

// TestMongoDBConnection 测试MongoDB连接
func TestMongoDBConnection(t *testing.T) {
	// 加载配置
	_, err := config.LoadConfig(".")
	if err != nil {
		t.Skipf("Skipping test: cannot load config: %v", err)
	}

	// 初始化服务容器，这会自动创建MongoDB连接
	err = service.InitializeServices()
	if err != nil {
		t.Skipf("Skipping test: cannot initialize services: %v", err)
	}

	// 获取服务容器实例
	container := service.GetServiceContainer()
	if container == nil {
		t.Fatal("Service container is nil")
	}

	// 测试数据库连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用容器中的数据库
	result := container.GetMongoClient().Database("admin").RunCommand(ctx, map[string]interface{}{"ping": 1})
	if result.Err() != nil {
		t.Skipf("Skipping test: MongoDB not available - %v", result.Err())
	}

	t.Log("MongoDB connection test passed")
}
