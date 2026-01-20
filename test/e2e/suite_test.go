//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"testing"

	layer1 "Qingyu_backend/test/e2e/layer1_basic"
	layer2 "Qingyu_backend/test/e2e/layer2_consistency"
	layer3 "Qingyu_backend/test/e2e/layer3_boundary"
)

// TestE2ESuite E2E测试套件入口
// 运行所有三层E2E测试
func TestE2ESuite(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试套件")
	}

	t.Log("========================================")
	t.Log("E2E 测试套件启动")
	t.Log("========================================")

	// Layer 1: 基础流程测试 (2-3分钟)
	t.Run("Layer1_基础流程测试", func(t *testing.T) {
		t.Log("运行 Layer 1 基础流程测试...")
		testLayer1Basic(t)
	})

	// Layer 2: 数据一致性测试 (3-5分钟)
	t.Run("Layer2_数据一致性测试", func(t *testing.T) {
		t.Log("运行 Layer 2 数据一致性测试...")
		testLayer2Consistency(t)
	})

	// Layer 3: 边界场景测试 (5-8分钟)
	t.Run("Layer3_边界场景测试", func(t *testing.T) {
		t.Log("运行 Layer 3 边界场景测试...")
		testLayer3Boundary(t)
	})

	t.Log("========================================")
	t.Log("E2E 测试套件完成")
	t.Log("========================================")
}

// testLayer1Basic 运行 Layer 1 基础流程测试
func testLayer1Basic(t *testing.T) {
	t.Log("----------------------------------------")
	t.Log("Layer 1: 基础流程测试")
	t.Log("测试目标: 验证核心业务流程的正确性")
	t.Log("预计耗时: 2-3分钟")
	t.Log("----------------------------------------")

	testCases := []struct {
		name string
		fn   func(*testing.T)
	}{
		{"认证流程", layer1.RunAuthFlow},
		{"阅读流程", layer1.RunReadingFlow},
		{"社交流程", layer1.RunSocialFlow},
		{"写作流程", layer1.RunWritingFlow},
	}

	passed := 0
	failed := 0

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fn(t)
			if !t.Failed() {
				passed++
			} else {
				failed++
			}
		})
	}

	t.Logf("----------------------------------------")
	t.Logf("Layer 1 完成: 通过 %d, 失败 %d", passed, failed)
}

// testLayer2Consistency 运行 Layer 2 数据一致性测试
func testLayer2Consistency(t *testing.T) {
	t.Log("----------------------------------------")
	t.Log("Layer 2: 数据一致性测试")
	t.Log("测试目标: 验证跨模块数据一致性")
	t.Log("预计耗时: 3-5分钟")
	t.Log("----------------------------------------")

	testCases := []struct {
		name string
		fn   func(*testing.T)
	}{
		{"用户阅读一致性", layer2.RunUserReadingConsistency},
		{"书籍章节一致性", layer2.RunBookChapterConsistency},
		{"社交互动一致性", layer2.RunSocialInteractionConsistency},
	}

	passed := 0
	failed := 0

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fn(t)
			if !t.Failed() {
				passed++
			} else {
				failed++
			}
		})
	}

	t.Logf("----------------------------------------")
	t.Logf("Layer 2 完成: 通过 %d, 失败 %d", passed, failed)
}

// testLayer3Boundary 运行 Layer 3 边界场景测试
func testLayer3Boundary(t *testing.T) {
	t.Log("----------------------------------------")
	t.Log("Layer 3: 边界场景测试")
	t.Log("测试目标: 验证边界条件和并发场景")
	t.Log("预计耗时: 5-8分钟")
	t.Log("----------------------------------------")

	testCases := []struct {
		name string
		fn   func(*testing.T)
	}{
		{"并发阅读", layer3.RunConcurrentReading},
		{"并发社交互动", layer3.RunConcurrentSocialInteraction},
		{"边界数据量", layer3.RunBoundaryDataSizes},
	}

	passed := 0
	failed := 0

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fn(t)
			if !t.Failed() {
				passed++
			} else {
				failed++
			}
		})
	}

	t.Logf("----------------------------------------")
	t.Logf("Layer 3 完成: 通过 %d, 失败 %d", passed, failed)
}

// TestE2EQuick 快速E2E测试（仅Layer 1）
// 用于CI/CD流水线中的快速验证
func TestE2EQuick(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	t.Log("运行快速E2E测试（仅Layer 1）...")
	testLayer1Basic(t)
}

// TestE2EStandard 标准E2E测试（Layer 1 + Layer 2）
// 用于CI/CD流水线中的标准验证
func TestE2EStandard(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	t.Log("运行标准E2E测试（Layer 1 + Layer 2）...")
	testLayer1Basic(t)
	testLayer2Consistency(t)
}

// PrintE2ESummary 打印E2E测试摘要
func PrintE2ESummary() {
	fmt.Println(`========================================
E2E 测试套件
========================================

测试层级:
  Layer 1: 基础流程测试 (2-3分钟)
    - 认证流程: 注册 -> 登录 -> 验证
    - 阅读流程: 浏览 -> 阅读 -> 保存进度
    - 社交流程: 评论 -> 收藏 -> 点赞
    - 写作流程: 创建项目 -> 验证

  Layer 2: 数据一致性测试 (3-5分钟)
    - 用户阅读一致性
    - 书籍章节一致性
    - 社交互动一致性

  Layer 3: 边界场景测试 (5-8分钟)
    - 并发阅读
    - 并发社交互动
    - 边界数据量

运行命令:
  make test-e2e              # 运行全部测试
  make test-e2e-quick        # 仅Layer 1
  make test-e2e-standard     # Layer 1 + Layer 2
  make test-e2e-layer1       # 仅Layer 1
  make test-e2e-layer2       # 仅Layer 2
  make test-e2e-layer3       # 仅Layer 3

注意事项:
  - 需要MongoDB服务运行
  - Redis降级模式下可运行
  - 测试数据会自动清理

========================================`)
}
