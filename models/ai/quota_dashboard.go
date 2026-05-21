package ai

// QuotaDashboard 配额仪表盘DTO
type QuotaDashboard struct {
	Summary        DashboardSummary     `json:"summary"`         // 汇总信息
	Distribution  QuotaDistribution    `json:"distribution"`    // 分布统计
	TopConsumers  []UserQuotaRanking   `json:"topConsumers"`    // 消费排行
	RecentAlerts  []AlertSummary       `json:"recentAlerts"`    // 近期告警
	TrendData     []TrendPoint         `json:"trendData"`       // 趋势数据
}

// DashboardSummary 仪表盘汇总数据
type DashboardSummary struct {
	TotalUsers      int64   `json:"totalUsers"`       // 总用户数
	ActiveUsers     int64   `json:"activeUsers"`      // 活跃用户数
	ExhaustedUsers  int64   `json:"exhaustedUsers"`   // 配额耗尽用户数
	NearExhaustUsers int64   `json:"nearExhaustUsers"` // 接近耗尽用户数
	SuspendedUsers  int64   `json:"suspendedUsers"`   // 暂停用户数
	TotalConsumption int64   `json:"totalConsumption"` // 总消费量
	AvgConsumption  float64 `json:"avgConsumption"`   // 平均消费量
}

// QuotaDistribution 配额分布统计
type QuotaDistribution struct {
	ByRole   map[string]int `json:"byRole"`   // 按角色分布
	ByLevel  map[string]int `json:"byLevel"`  // 按等级分布
	ByService map[string]int `json:"byService"` // 按服务分布
	ByStatus map[string]int `json:"byStatus"` // 按状态分布
}

// UserQuotaRanking 用户配额排行
type UserQuotaRanking struct {
	UserID       string  `json:"userId"`       // 用户ID
	Username     string  `json:"username"`     // 用户名
	Role         string  `json:"role"`         // 角色
	UsedQuota    int     `json:"usedQuota"`    // 已用配额
	TotalQuota   int     `json:"totalQuota"`   // 总配额
	UsagePercent float64 `json:"usagePercent"` // 使用百分比
}

// TrendPoint 趋势数据点
type TrendPoint struct {
	Date         string `json:"date"`          // 日期
	Consumption  int    `json:"consumption"`   // 消费量
	Users        int    `json:"users"`         // 用户数
}

// UserQuotaListItem 用户配额列表项DTO
type UserQuotaListItem struct {
	UserID       string  `json:"userId"`
	Username     string  `json:"username"`
	Role         string  `json:"role"`
	MemberLevel  string  `json:"memberLevel"`
	DailyQuota   int     `json:"dailyQuota"`
	DailyUsed    int     `json:"dailyUsed"`
	UsagePercent float64 `json:"usagePercent"`
	Status       string  `json:"status"`
}

// AlertSummary 告警摘要
type AlertSummary struct {
	ID        string           `json:"id"`         // 告警ID
	Type      string           `json:"type"`       // 告警类型
	Level     string           `json:"level"`      // 告警级别
	Title     string           `json:"title"`      // 告警标题
	UserID    string           `json:"userId"`     // 用户ID
	Status    string           `json:"status"`     // 状态
	CreatedAt string           `json:"createdAt"`  // 创建时间
}
