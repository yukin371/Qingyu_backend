package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/config"
	"Qingyu_backend/pkg/monitor"
)

func main() {
	log.Println("=== Qingyu 数据质量监控 ===")
	log.Println("正在启动数据质量监控服务...")

	// 加载数据库配置
	cfg, err := config.LoadDatabaseConfig("")
	if err != nil {
		log.Fatalf("加载数据库配置失败: %v", err)
	}

	mongoConfig, err := cfg.GetMongoConfig()
	if err != nil {
		log.Fatalf("获取MongoDB配置失败: %v", err)
	}

	// 连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConfig.URI))
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("断开MongoDB连接失败: %v", err)
		}
	}()

	// 验证连接
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Ping MongoDB失败: %v", err)
	}

	db := client.Database(mongoConfig.Database)
	log.Printf("✓ MongoDB连接成功: %s (数据库: %s)", mongoConfig.URI, mongoConfig.Database)

	// 创建监控器
	mongodb := monitor.NewMongoDB(db)
	alerter := monitor.GetAlerterFromEnv()
	qualityMonitor := monitor.NewDataQualityMonitor(mongodb, alerter)

	log.Println("✓ 数据质量监控器初始化完成")

	// 设置检查间隔
	checkInterval := 24 * time.Hour
	if intervalStr := os.Getenv("DATA_QUALITY_CHECK_INTERVAL"); intervalStr != "" {
		if duration, err := time.ParseDuration(intervalStr); err == nil {
			checkInterval = duration
			log.Printf("✓ 自定义检查间隔: %v", checkInterval)
		}
	}

	// 创建定时器
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	// 设置优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Println("\n=== 开始执行数据质量检查 ===")
	log.Printf("检查间隔: %v\n", checkInterval)

	// 立即执行一次检查
	if err := runCheck(context.Background(), qualityMonitor); err != nil {
		log.Printf("⚠ 首次检查失败: %v\n", err)
	}

	// 主循环
	for {
		select {
		case <-ticker.C:
			log.Println("\n=== 执行定时检查 ===")
			if err := runCheck(context.Background(), qualityMonitor); err != nil {
				log.Printf("⚠ 检查失败: %v\n", err)
			}
		case sig := <-sigChan:
			log.Printf("\n收到信号: %v，正在退出...\n", sig)
			return
		}
	}
}

// runCheck 执行单次检查
func runCheck(ctx context.Context, qualityMonitor *monitor.DataQualityMonitor) error {
	startTime := time.Now()
	log.Printf("检查开始时间: %s", startTime.Format(time.RFC3339))

	report, err := qualityMonitor.RunDailyCheck(ctx)
	if err != nil {
		return fmt.Errorf("数据质量检查失败: %w", err)
	}

	duration := time.Since(startTime)
	log.Printf("检查完成，耗时: %v", duration)

	// 输出检查结果
	printReport(report)

	return nil
}

// printReport 打印报告
func printReport(report *monitor.DataQualityReport) {
	log.Println("\n=== 数据质量报告 ===")
	log.Printf("检查时间: %s", report.CheckTime.Format(time.RFC3339))
	log.Printf("孤儿记录数: %d", report.TotalOrphanedRecords)
	log.Printf("统计不准确数: %d", report.InaccurateStatistics)
	log.Printf("是否有问题: %t", report.HasIssues)

	if len(report.Details) > 0 {
		log.Println("\n问题详情:")
		for i, detail := range report.Details {
			log.Printf("  %d. %s.%s: %d条", i+1, detail.Collection, detail.Field, detail.Count)
		}
	}

	if report.HasIssues {
		log.Println("\n⚠ 发现数据质量问题，已发送告警")
	} else {
		log.Println("\n✓ 数据质量检查通过")
	}
	log.Println("==================")
}
