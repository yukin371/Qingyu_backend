package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) *mongo.Database {
	// 从环境变量获取MongoDB连接字符串
	mongoURI := os.Getenv("MONGODB_TEST_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// 创建MongoDB客户端
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("连接MongoDB失败: %v", err)
	}

	// Ping测试
	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("Ping MongoDB失败: %v", err)
	}

	// 使用测试数据库
	testDBName := fmt.Sprintf("qingyu_test_%d", time.Now().Unix())
	db := client.Database(testDBName)

	t.Logf("✓ 测试数据库已创建: %s", testDBName)

	return db
}

// cleanupTestDB 清理测试数据库
func cleanupTestDB(t *testing.T, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 删除测试数据库
	if err := db.Drop(ctx); err != nil {
		t.Logf("⚠ 清理测试数据库失败: %v", err)
	} else {
		t.Logf("✓ 测试数据库已清理: %s", db.Name())
	}

	// 断开连接
	if err := db.Client().Disconnect(ctx); err != nil {
		t.Logf("⚠ 断开MongoDB连接失败: %v", err)
	}
}
