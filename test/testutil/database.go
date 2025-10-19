package testutil

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"
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

	// 初始化数据库
	err = core.InitDB()
	if err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}

	// 清理函数
	cleanup := func() {
		// 清理测试集合
		ctx := context.Background()

		// 推荐系统测试集合
		_ = global.DB.Collection("user_behaviors").Drop(ctx)
		_ = global.DB.Collection("user_profiles").Drop(ctx)
		_ = global.DB.Collection("item_features").Drop(ctx)
		_ = global.DB.Collection("book_statistics").Drop(ctx)
		_ = global.DB.Collection("books").Drop(ctx)

		// Writing相关测试集合
		_ = global.DB.Collection("projects").Drop(ctx)
		_ = global.DB.Collection("documents").Drop(ctx)
		_ = global.DB.Collection("document_contents").Drop(ctx)
	}

	return global.DB, cleanup
}
