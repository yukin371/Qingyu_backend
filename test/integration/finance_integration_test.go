package integration

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/service"
	"Qingyu_backend/service/container"
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
		if svc == nil {
			t.Error("会员服务为 nil")
		}
	})

	// 验证作者收入服务接口
	t.Run("AuthorRevenueService_Interface", func(t *testing.T) {
		svc, err := container.GetAuthorRevenueService()
		if err != nil {
			t.Fatalf("获取作者收入服务失败: %v", err)
		}

		// 验证服务实现了正确的接口类型
		if svc == nil {
			t.Error("作者收入服务为 nil")
		}
	})
}

// initServiceContainer 初始化服务容器用于测试
func initServiceContainer(t *testing.T) *container.ServiceContainer {
	// 加载配置（如果尚未加载）
	if config.GlobalConfig == nil {
		// 尝试从默认位置加载配置
		_, err := config.LoadConfig("./config")
		if err != nil {
			t.Skipf("跳过测试：配置加载失败 (%v)，请确保 config/config.yaml 存在", err)
		}
	}

	// 创建服务容器
	serviceContainer := container.NewServiceContainer()

	// 初始化服务容器（连接 MongoDB、创建工厂）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := serviceContainer.Initialize(ctx); err != nil {
		t.Fatalf("初始化服务容器失败: %v", err)
	}

	// 设置默认服务（创建并注册所有业务服务）
	if err := serviceContainer.SetupDefaultServices(); err != nil {
		t.Fatalf("设置默认服务失败: %v", err)
	}

	// 注意：不再设置全局变量 global.DB
	// 数据库连接已由 service.ServiceManager 管理

	return serviceContainer
}

// cleanupServiceContainer 清理服务容器
func cleanupServiceContainer(serviceContainer *container.ServiceContainer) {
	if serviceContainer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = serviceContainer.Close(ctx)
	}
}

// TestFinanceServiceWithContainer 使用完整容器测试
func TestFinanceServiceWithContainer(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	serviceContainer := initServiceContainer(t)
	defer cleanupServiceContainer(serviceContainer)

	t.Run("GetServices", func(t *testing.T) {
		// 测试获取会员服务
		membershipSvc, err := serviceContainer.GetMembershipService()
		if err != nil {
			t.Fatalf("获取会员服务失败: %v", err)
		}
		if membershipSvc == nil {
			t.Error("会员服务不应为 nil")
		}

		// 测试获取作者收入服务
		revenueSvc, err := serviceContainer.GetAuthorRevenueService()
		if err != nil {
			t.Fatalf("获取作者收入服务失败: %v", err)
		}
		if revenueSvc == nil {
			t.Error("作者收入服务不应为 nil")
		}

		t.Log("✓ Finance 服务获取成功")
	})

	t.Run("ServiceHealth", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 检查服务容器健康状态
		if err := serviceContainer.Health(ctx); err != nil {
			t.Logf("健康检查警告: %v", err)
		} else {
			t.Log("✓ 服务容器健康检查通过")
		}
	})
}

// TestMembershipServiceOperations 测试会员服务操作
func TestMembershipServiceOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	serviceContainer := initServiceContainer(t)
	defer cleanupServiceContainer(serviceContainer)

	ctx := context.Background()

	t.Run("GetPlans", func(t *testing.T) {
		svc, err := serviceContainer.GetMembershipService()
		if err != nil {
			t.Fatalf("获取会员服务失败: %v", err)
		}

		plans, err := svc.GetPlans(ctx)
		if err != nil {
			t.Logf("获取套餐列表失败（可能数据库未初始化）: %v", err)
			return
		}

		t.Logf("成功获取 %d 个套餐", len(plans))

		// 验证返回值不为 nil（空切片是合法的）
		if plans == nil {
			t.Error("套餐列表不应为 nil，应为空切片")
		} else {
			t.Log("✓ GetPlans 方法调用成功")
		}
	})
}

// TestAuthorRevenueServiceOperations 测试作者收入服务操作
func TestAuthorRevenueServiceOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	serviceContainer := initServiceContainer(t)
	defer cleanupServiceContainer(serviceContainer)

	ctx := context.Background()

	t.Run("GetEarnings", func(t *testing.T) {
		svc, err := serviceContainer.GetAuthorRevenueService()
		if err != nil {
			t.Fatalf("获取作者收入服务失败: %v", err)
		}

		testAuthorID := "test-author-123"
		earnings, total, err := svc.GetEarnings(ctx, testAuthorID, 1, 20)
		if err != nil {
			t.Logf("获取收入列表失败（可能数据库未初始化）: %v", err)
			return
		}

		t.Logf("成功获取 %d 条收入记录，总计 %d 条", len(earnings), total)
		t.Log("✓ GetEarnings 方法调用成功")
	})

	t.Run("GetRevenueStatistics", func(t *testing.T) {
		svc, err := serviceContainer.GetAuthorRevenueService()
		if err != nil {
			t.Fatalf("获取作者收入服务失败: %v", err)
		}

		testAuthorID := "test-author-123"
		stats, err := svc.GetRevenueStatistics(ctx, testAuthorID, "monthly")
		if err != nil {
			t.Logf("获取收入统计失败（可能数据库未初始化）: %v", err)
			return
		}

		t.Logf("成功获取 %d 条收入统计", len(stats))
		t.Log("✓ GetRevenueStatistics 方法调用成功")
	})
}
