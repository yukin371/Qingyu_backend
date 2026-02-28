//go:build e2e
// +build e2e

package data

import (
	"testing"

	"github.com/stretchr/testify/require"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"
	"Qingyu_backend/service"
	"Qingyu_backend/test/testutil"
)

// SetupTestEnvironment 设置测试环境
// 这个函数初始化数据库连接和其他必要的服务
// 必须在运行任何测试之前调用
func SetupTestEnvironment(t *testing.T) {
	t.Helper()

	// 加载配置
	cfg, err := config.LoadConfig("../../../config")
	require.NoError(t, err, "加载配置失败")

	// 设置全局配置
	config.GlobalConfig = cfg
	testutil.EnableStrictLogAssertionsIgnoreWarn(t)

	// 初始化服务（包括数据库连接）
	err = core.InitServices()
	require.NoError(t, err, "初始化服务失败")

	// 设置全局DB变量（向后兼容）
	if service.ServiceManager != nil {
		global.DB = service.ServiceManager.GetMongoDB()
		global.MongoClient = service.ServiceManager.GetMongoClient()
	}
}

