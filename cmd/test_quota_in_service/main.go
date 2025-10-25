package main

import (
	"context"
	"fmt"
	"log"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/service/container"
)

func main() {
	fmt.Println("========== 测试服务中的配额检查 ==========\n")

	// 1. 加载配置
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}
	fmt.Println("✓ 配置加载成功")

	// 2. 初始化数据库
	err = core.InitDB()
	if err != nil {
		log.Fatal("初始化数据库失败:", err)
	}
	fmt.Println("✓ 数据库连接成功")

	// 3. 初始化服务容器
	ctx := context.Background()
	serviceContainer := container.NewServiceContainer()
	err = serviceContainer.Initialize(ctx)
	if err != nil {
		log.Fatal("初始化服务容器失败:", err)
	}
	fmt.Println("✓ 服务容器初始化成功")

	// 4. 设置默认服务（包括QuotaService）
	err = serviceContainer.SetupDefaultServices()
	if err != nil {
		log.Fatal("设置默认服务失败:", err)
	}
	fmt.Println("✓ 默认服务设置成功")

	// 5. 获取QuotaService
	quotaService, err := serviceContainer.GetQuotaService()
	if err != nil {
		log.Fatal("获取QuotaService失败:", err)
	}
	if quotaService == nil {
		log.Fatal("QuotaService未初始化")
	}
	fmt.Println("✓ QuotaService获取成功")

	fmt.Println()

	// 6. 测试配额检查
	// 使用已创建的ctx

	// 使用vip_user01的ID (从MongoDB获取)
	userID := "68fc402acd736a40d4220ba4" // vip_user01

	fmt.Printf("测试用户: vip_user01\n")
	fmt.Printf("用户ID: %s\n\n", userID)

	// 7. 检查配额
	fmt.Println("检查配额 (请求100 tokens)...")
	err = quotaService.CheckQuota(ctx, userID, 100)

	if err != nil {
		fmt.Printf("❌ 配额检查失败: %v\n", err)
		fmt.Println("\n这就是为什么测试返回429错误！")

		// 尝试获取配额详情
		fmt.Println("\n获取配额详情...")
		quota, getErr := quotaService.GetQuotaInfo(ctx, userID)
		if getErr != nil {
			fmt.Printf("获取配额失败: %v\n", getErr)
		} else {
			fmt.Printf("  - QuotaType: %s\n", quota.QuotaType)
			fmt.Printf("  - TotalQuota: %d\n", quota.TotalQuota)
			fmt.Printf("  - UsedQuota: %d\n", quota.UsedQuota)
			fmt.Printf("  - RemainingQuota: %d\n", quota.RemainingQuota)
			fmt.Printf("  - Status: %s\n", quota.Status)
			fmt.Printf("  - ResetAt: %v\n", quota.ResetAt)
			fmt.Printf("  - IsAvailable(): %v\n", quota.IsAvailable())
			fmt.Printf("  - CanConsume(100): %v\n", quota.CanConsume(100))
		}
	} else {
		fmt.Println("✅ 配额检查通过！")

		// 获取配额详情
		quota, _ := quotaService.GetQuotaInfo(ctx, userID)
		fmt.Printf("\n配额详情:\n")
		fmt.Printf("  - QuotaType: %s\n", quota.QuotaType)
		fmt.Printf("  - TotalQuota: %d\n", quota.TotalQuota)
		fmt.Printf("  - RemainingQuota: %d\n", quota.RemainingQuota)
		fmt.Printf("  - Status: %s\n", quota.Status)
	}

	fmt.Println("\n==========================================")
}
