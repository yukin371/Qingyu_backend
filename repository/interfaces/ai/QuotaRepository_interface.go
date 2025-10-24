package ai

import (
	"context"
	"time"

	aiModels "Qingyu_backend/models/ai"
)

// QuotaRepository 配额Repository接口
type QuotaRepository interface {
	// 配额管理
	CreateQuota(ctx context.Context, quota *aiModels.UserQuota) error
	GetQuotaByUserID(ctx context.Context, userID string, quotaType aiModels.QuotaType) (*aiModels.UserQuota, error)
	UpdateQuota(ctx context.Context, quota *aiModels.UserQuota) error
	DeleteQuota(ctx context.Context, userID string, quotaType aiModels.QuotaType) error

	// 批量操作
	GetAllQuotasByUserID(ctx context.Context, userID string) ([]*aiModels.UserQuota, error)
	BatchResetQuotas(ctx context.Context, quotaType aiModels.QuotaType) error

	// 配额事务记录
	CreateTransaction(ctx context.Context, transaction *aiModels.QuotaTransaction) error
	GetTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*aiModels.QuotaTransaction, error)
	GetTransactionsByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*aiModels.QuotaTransaction, error)

	// 统计查询
	GetQuotaStatistics(ctx context.Context, userID string) (*QuotaStatistics, error)
	GetTotalConsumption(ctx context.Context, userID string, quotaType aiModels.QuotaType, startTime, endTime time.Time) (int, error)

	// Health 健康检查
	Health(ctx context.Context) error
}

// QuotaStatistics 配额统计信息
type QuotaStatistics struct {
	UserID            string         `json:"userId"`
	TotalQuota        int            `json:"totalQuota"`        // 总配额
	UsedQuota         int            `json:"usedQuota"`         // 已用配额
	RemainingQuota    int            `json:"remainingQuota"`    // 剩余配额
	UsagePercentage   float64        `json:"usagePercentage"`   // 使用百分比
	TotalTransactions int            `json:"totalTransactions"` // 总交易次数
	DailyAverage      float64        `json:"dailyAverage"`      // 日均消费
	QuotaByType       map[string]int `json:"quotaByType"`       // 按类型统计
	QuotaByService    map[string]int `json:"quotaByService"`    // 按服务统计
}
