package core

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/global"
)

// TestMongoDBConnection 测试MongoDB连接
func TestMongoDBConnection(t *testing.T) {
	// 加载配置
	_, err := config.LoadConfig(".")
	if err != nil {
		t.Skipf("Skipping test: cannot load config: %v", err)
	}
	
	// 初始化MongoDB连接
	err = InitDB()
	if err != nil {
		t.Skipf("Skipping test: cannot initialize MongoDB: %v", err)
	}

	// 检查全局变量是否已设置
	if global.MongoClient == nil {
		t.Fatal("MongoDB client is nil")
	}

	if global.DB == nil {
		t.Fatal("MongoDB database is nil")
	}

	// 测试数据库连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	result := global.MongoClient.Database("admin").RunCommand(ctx, map[string]interface{}{"ping": 1})
	if result.Err() != nil {
		t.Fatalf("Failed to ping MongoDB: %v", result.Err())
	}

	t.Log("MongoDB connection test passed")
}