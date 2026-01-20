//go:build e2e
// +build e2e

package layer1_basic

import (
	"testing"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/users"
	e2e "Qingyu_backend/test/e2e/framework"
)

// TestReadingFlow 测试阅读流程
// 流程: 浏览书城 -> 查看书籍详情 -> 获取章节列表 -> 阅读章节内容 -> 保存阅读进度
// TestReadingFlow 测试阅读流程
// 流程: 浏览书城 -> 查看书籍详情 -> 获取章节列表 -> 阅读章节内容 -> 保存阅读进度
func TestReadingFlow(t *testing.T) {
	RunReadingFlow(t)
}



