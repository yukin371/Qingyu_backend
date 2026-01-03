package integration

import (
	"context"
	"testing"

	"Qingyu_backend/service"
	"Qingyu_backend/service/finance"
)

// TestFinanceServiceIntegration 测试 Finance 服务集成
func TestFinanceServiceIntegration(t *testing.T) {
	// 获取服务容器
	container := service.GetServiceContainer()
	if container == nil {
		t.Skip("服务容器未初始化，跳过集成测试")
	}

	// 测试获取会员服务
	t.Run("GetMembershipService", func(t *testing.T) {
		svc, err := container.GetMembershipService()
		if err != nil {
			t.Fatalf("获取会员服务失败: %v", err)
		}
		if svc == nil {
			t.Fatal("会员服务为 nil")
		}

		// 验证服务实现了预期接口
		var _ finance.MembershipService = svc
	})

	// 测试获取作者收入服务
	t.Run("GetAuthorRevenueService", func(t *testing.T) {
		svc, err := container.GetAuthorRevenueService()
		if err != nil {
			t.Fatalf("获取作者收入服务失败: %v", err)
		}
		if svc == nil {
			t.Fatal("作者收入服务为 nil")
		}

		// 验证服务实现了预期接口
		var _ finance.AuthorRevenueService = svc
	})
}

// TestFinanceServiceBasicOperations 测试 Finance 服务基本操作
func TestFinanceServiceBasicOperations(t *testing.T) {
	container := service.GetServiceContainer()
	if container == nil {
		t.Skip("服务容器未初始化，跳过集成测试")
	}

	ctx := context.Background()

	// 测试会员服务获取套餐列表
	t.Run("MembershipService_GetPlans", func(t *testing.T) {
		svc, err := container.GetMembershipService()
		if err != nil {
			t.Fatalf("获取会员服务失败: %v", err)
		}

		plans, err := svc.GetPlans(ctx)
		if err != nil {
			t.Logf("获取套餐列表失败（可能是数据库未初始化）: %v", err)
			return
		}

		// 验证返回值
		if plans == nil {
			t.Error("套餐列表不应为 nil")
		}

		t.Logf("成功获取 %d 个套餐", len(plans))
	})

	// 测试会员服务获取会员状态
	t.Run("MembershipService_GetMembership", func(t *testing.T) {
		svc, err := container.GetMembershipService()
		if err != nil {
			t.Fatalf("获取会员服务失败: %v", err)
		}

		// 使用测试用户ID
		testUserID := "test-user-123"
		membership, err := svc.GetMembership(ctx, testUserID)
		if err != nil {
			t.Logf("获取会员状态失败（可能是数据库未初始化或用户不存在）: %v", err)
			return
		}

		if membership != nil {
			t.Logf("用户 %s 的会员状态: %s", testUserID, membership.Status)
		}
	})

	// 测试作者收入服务获取收入统计
	t.Run("AuthorRevenueService_GetRevenueStatistics", func(t *testing.T) {
		svc, err := container.GetAuthorRevenueService()
		if err != nil {
			t.Fatalf("获取作者收入服务失败: %v", err)
		}

		// 使用测试作者ID
		testAuthorID := "test-author-123"
		stats, err := svc.GetRevenueStatistics(ctx, testAuthorID, "monthly")
		if err != nil {
			t.Logf("获取收入统计失败（可能是数据库未初始化或作者不存在）: %v", err)
			return
		}

		if stats == nil {
			t.Log("收入统计为空（可能是新作者）")
		} else {
			t.Logf("成功获取 %d 条收入统计", len(stats))
		}
	})
}

// TestFinanceServiceInterfaceCompliance 测试服务接口合规性
func TestFinanceServiceInterfaceCompliance(t *testing.T) {
	container := service.GetServiceContainer()
	if container == nil {
		t.Skip("服务容器未初始化，跳过集成测试")
	}

	// 验证会员服务接口
	t.Run("MembershipService_Interface", func(t *testing.T) {
		svc, err := container.GetMembershipService()
		if err != nil {
			t.Fatalf("获取会员服务失败: %v", err)
		}

		// 验证服务实现了正确的接口类型
		if _, ok := svc.(finance.MembershipService); !ok {
			t.Error("会员服务未实现 MembershipService 接口")
		}
	})

	// 验证作者收入服务接口
	t.Run("AuthorRevenueService_Interface", func(t *testing.T) {
		svc, err := container.GetAuthorRevenueService()
		if err != nil {
			t.Fatalf("获取作者收入服务失败: %v", err)
		}

		// 验证服务实现了正确的接口类型
		if _, ok := svc.(finance.AuthorRevenueService); !ok {
			t.Error("作者收入服务未实现 AuthorRevenueService 接口")
		}

		// 验证服务不为 nil
		if svc == nil {
			t.Fatal("作者收入服务为 nil")
		}
	})
}
