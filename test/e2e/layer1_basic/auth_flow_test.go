//go:build e2e
// +build e2e

package layer1_basic

import (
	"testing"
)

// TestAuthFlow 测试认证流程
// 流程: 注册 -> 登录 -> 获取用户信息 -> 登出
// TestAuthFlow 测试认证流程
// 流程: 注册 -> 登录 -> 获取用户信息 -> 登出
func TestAuthFlow(t *testing.T) {
	RunAuthFlow(t)
}



