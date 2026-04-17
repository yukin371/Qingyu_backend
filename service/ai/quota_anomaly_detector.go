package ai

import (
	"context"
	"fmt"
	"math"
	"time"

	aiModels "Qingyu_backend/models/ai"
	aiInterfaces "Qingyu_backend/repository/interfaces/ai"

	"github.com/redis/go-redis/v9"
)

// QuotaAnomalyDetector 配额异常检测器
type QuotaAnomalyDetector struct {
	quotaRepo   aiInterfaces.QuotaRepository
	alertService *QuotaAlertService
	redisClient interface{} // 可以是 *redis.Client
}

// NewQuotaAnomalyDetector 创建配额异常检测器
func NewQuotaAnomalyDetector(quotaRepo aiInterfaces.QuotaRepository, alertService *QuotaAlertService, redisClient interface{}) *QuotaAnomalyDetector {
	return &QuotaAnomalyDetector{
		quotaRepo:   quotaRepo,
		alertService: alertService,
		redisClient: redisClient,
	}
}

// DetectAnomalies 运行所有异常检测算法
func (d *QuotaAnomalyDetector) DetectAnomalies(ctx context.Context) error {
	// 检测突发性异常
	if err := d.detectSpikeUsage(ctx); err != nil {
		fmt.Printf("突发性异常检测失败: %v\n", err)
	}

	// 检测滥用异常
	if err := d.detectRateAbuse(ctx); err != nil {
		fmt.Printf("滥用异常检测失败: %v\n", err)
	}

	// 检测配额耗尽
	if err := d.detectQuotaExhaustion(ctx); err != nil {
		fmt.Printf("配额耗尽检测失败: %v\n", err)
	}

	return nil
}

// detectSpikeUsage 检测突发性异常（Z-Score检测）
func (d *QuotaAnomalyDetector) detectSpikeUsage(ctx context.Context) error {
	// 获取过去7天的交易记录
	startTime := time.Now().AddDate(0, 0, -7)
	endTime := time.Now()

	// 查询用户日消费统计
	// 这里需要先获取所有活跃用户
	users, _, err := d.quotaRepo.ListUserQuotas(ctx, "", "", "", 1, 1000)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 按用户分组统计日均消费
	userConsumption := make(map[string][]int)
	for _, user := range users {
		// 模拟获取用户过去7天的消费数据
		dailyConsumption, err := d.getUserDailyConsumption(ctx, user.UserID, startTime, endTime)
		if err != nil {
			continue
		}
		if err != nil {
			continue
		}
		userConsumption[user.UserID] = dailyConsumption
	}

	// 计算每个用户的Z-Score
	for userID, consumptions := range userConsumption {
		if len(consumptions) < 7 {
			continue // 数据不足，跳过
		}

		// 计算均值和标准差
		mean := calculateMean(consumptions)
		std := calculateStd(consumptions, mean)

		if std == 0 {
			continue // 标准差为0，无波动
		}

		// 获取今天的消费数据
		todayConsumption, err := d.getUserDailyConsumption(ctx, userID, time.Now().AddDate(0, 0, -1), time.Now())
		if err != nil || len(todayConsumption) == 0 {
			continue
		}

		// 计算Z-Score
		zScore := float64(todayConsumption[0]-int(mean)) / std

		// Z-Score > 2 视为异常
		if zScore > 2 {
			// 创建异常告警
			alert := &aiModels.QuotaAlert{
				Type:    aiModels.QuotaAlertTypeAnomaly,
				UserID:  userID,
				Level:   aiModels.QuotaAlertLevelCritical,
				Title:   "突发性消费异常",
				Message: fmt.Sprintf("今日消费(%d)显著高于历史均值(%.2f)，Z-Score: %.2f", todayConsumption[0], mean, zScore),
				Data: map[string]interface{}{
					"zScore":        zScore,
					"todayValue":    todayConsumption[0],
					"mean":         mean,
					"std":          std,
					"consumptions":  consumptions,
				},
			}

			if err := d.alertService.CreateAlert(ctx, alert); err != nil {
				fmt.Printf("创建异常告警失败[%s]: %v\n", userID, err)
			}
		}
	}

	return nil
}

// detectRateAbuse 检测滥用异常
func (d *QuotaAnomalyDetector) detectRateAbuse(ctx context.Context) error {
	// 如果Redis不可用，跳过检测
	if d.redisClient == nil {
		return nil
	}

	// 获取所有活跃用户
	users, _, err := d.quotaRepo.ListUserQuotas(ctx, "", "", "", 1, 1000)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 检查每个用户的速率限制
	for _, user := range users {
		// 检查5分钟速率限制
		fiveMinKey := fmt.Sprintf("quota:rate:%s:5m", user.UserID)
		redisClient := d.redisClient.(*redis.Client) // 类型断言

		count5m, err := redisClient.Get(ctx, fiveMinKey).Int()
		if err != nil {
			continue // 读取失败，跳过
		}

		// 检查1小时速率限制
		oneHourKey := fmt.Sprintf("quota:rate:%s:1h", user.UserID)
		count1h, err := redisClient.Get(ctx, oneHourKey).Int()
		if err != nil {
			continue // 读取失败，跳过
		}

		// 阈值：5分钟超过50次，1小时超过200次
		if count5m > 50 {
			// 创建滥用告警
			alert := &aiModels.QuotaAlert{
				Type:    aiModels.QuotaAlertTypeAbuse,
				UserID:  user.UserID,
				Level:   aiModels.QuotaAlertLevelWarning,
				Title:   "API调用频率过高",
				Message: fmt.Sprintf("5分钟内调用%d次，超过阈值50", count5m),
				Data: map[string]interface{}{
					"count5m":     count5m,
					"count1h":     count1h,
					"threshold5m": 50,
				},
			}

			if err := d.alertService.CreateAlert(ctx, alert); err != nil {
				fmt.Printf("创建滥用告警失败[%s]: %v\n", user.UserID, err)
			}
		}

		if count1h > 200 {
			// 创建滥用告警（升级到严重级别）
			alert := &aiModels.QuotaAlert{
				Type:    aiModels.QuotaAlertTypeAbuse,
				UserID:  user.UserID,
				Level:   aiModels.QuotaAlertLevelCritical,
				Title:   "API调用频率严重超标",
				Message: fmt.Sprintf("1小时内调用%d次，超过阈值200", count1h),
				Data: map[string]interface{}{
					"count5m":     count5m,
					"count1h":     count1h,
					"threshold1h": 200,
				},
			}

			if err := d.alertService.CreateAlert(ctx, alert); err != nil {
				fmt.Printf("创建滥用告警失败[%s]: %v\n", user.UserID, err)
			}
		}
	}

	return nil
}

// detectQuotaExhaustion 检测配额耗尽
func (d *QuotaAnomalyDetector) detectQuotaExhaustion(ctx context.Context) error {
	// 获取所有用户的日配额
	userQuotas, _, err := d.quotaRepo.ListUserQuotas(ctx, "", "", "", 1, 1000)
	if err != nil {
		return fmt.Errorf("获取用户配额列表失败: %w", err)
	}

	for _, user := range userQuotas {
		if user.DailyQuota == 0 {
			continue // 无配额限制的用户
		}

		usedPercent := float64(user.DailyUsed) / float64(user.DailyQuota)
		remainingPercent := float64(user.DailyQuota-user.DailyUsed) / float64(user.DailyQuota)

		// 创建告警
		if remainingPercent == 0 {
			// 配额耗尽
			alert := &aiModels.QuotaAlert{
				Type:    aiModels.QuotaAlertTypeThreshold,
				UserID:  user.UserID,
				Level:   aiModels.QuotaAlertLevelCritical,
				Title:   "配额已耗尽",
				Message: fmt.Sprintf("用户配额已全部使用（已用%d/总量%d）", user.DailyUsed, user.DailyQuota),
				Data: map[string]interface{}{
					"usedQuota":    user.DailyUsed,
					"totalQuota":   user.DailyQuota,
					"usedPercent":  usedPercent * 100,
					"remaining":    user.DailyQuota - user.DailyUsed,
				},
			}

			if err := d.alertService.CreateAlert(ctx, alert); err != nil {
				fmt.Printf("创建配额耗尽告警失败[%s]: %v\n", user.UserID, err)
			}
		} else if remainingPercent < 0.1 {
			// 严重不足（剩余<10%）
			alert := &aiModels.QuotaAlert{
				Type:    aiModels.QuotaAlertTypeThreshold,
				UserID:  user.UserID,
				Level:   aiModels.QuotaAlertLevelCritical,
				Title:   "配额严重不足",
				Message: fmt.Sprintf("用户配额剩余%d（总量%d，剩余%.1f%%），即将耗尽",
					user.DailyQuota-user.DailyUsed, user.DailyQuota, remainingPercent*100),
				Data: map[string]interface{}{
					"usedQuota":    user.DailyUsed,
					"totalQuota":   user.DailyQuota,
					"usedPercent":  usedPercent * 100,
					"remaining":    user.DailyQuota - user.DailyUsed,
					"remainingPercent": remainingPercent * 100,
				},
			}

			if err := d.alertService.CreateAlert(ctx, alert); err != nil {
				fmt.Printf("创建严重不足告警失败[%s]: %v\n", user.UserID, err)
			}
		} else if remainingPercent < 0.2 {
			// 警告级别（剩余<20%）
			alert := &aiModels.QuotaAlert{
				Type:    aiModels.QuotaAlertTypeThreshold,
				UserID:  user.UserID,
				Level:   aiModels.QuotaAlertLevelWarning,
				Title:   "配额接近耗尽",
				Message: fmt.Sprintf("用户配额剩余%d（总量%d，剩余%.1f%%）",
					user.DailyQuota-user.DailyUsed, user.DailyQuota, remainingPercent*100),
				Data: map[string]interface{}{
					"usedQuota":    user.DailyUsed,
					"totalQuota":   user.DailyQuota,
					"usedPercent":  usedPercent * 100,
					"remaining":    user.DailyQuota - user.DailyUsed,
					"remainingPercent": remainingPercent * 100,
				},
			}

			if err := d.alertService.CreateAlert(ctx, alert); err != nil {
				fmt.Printf("创建接近耗尽告警失败[%s]: %v\n", user.UserID, err)
			}
		}
	}

	return nil
}

// 私有辅助方法

// getUserDailyConsumption 获取用户每日消费数据（模拟实现）
func (d *QuotaAnomalyDetector) getUserDailyConsumption(ctx context.Context, userID string, startTime, endTime time.Time) ([]int, error) {
	// 这里应该从仓库获取真实数据，为了演示，我们模拟一些数据
	days := int(endTime.Sub(startTime).Hours() / 24)
	consumption := make([]int, 0, days+1)

	for i := 0; i <= days; i++ {
		day := startTime.AddDate(0, 0, i)
		// 模拟消费数据：50-200之间的随机数
		value := 50 + (day.Unix() % 150)
		consumption = append(consumption, int(value))
	}

	return consumption, nil
}

// 计算平均值
func calculateMean(values []int) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range values {
		sum += float64(v)
	}

	return sum / float64(len(values))
}

// 计算标准差
func calculateStd(values []int, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sumSquares := 0.0
	for _, v := range values {
		diff := float64(v) - mean
		sumSquares += diff * diff
	}

	return math.Sqrt(sumSquares / float64(len(values)))
}