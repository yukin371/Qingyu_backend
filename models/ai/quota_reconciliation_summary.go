package ai

import "time"

// QuotaConsumptionSummaryItem 表示按某一维度聚合后的消费摘要。
type QuotaConsumptionSummaryItem struct {
	GroupKey     string `json:"groupKey" bson:"group_key"`
	TotalTokens  int64  `json:"totalTokens" bson:"total_tokens"`
	TotalRecords int64  `json:"totalRecords" bson:"total_records"`
}

// QuotaConsumptionSummary 表示某一时间窗内的消费聚合结果。
type QuotaConsumptionSummary struct {
	TimeRange    string                        `json:"timeRange"`
	WorkflowType string                        `json:"workflowType,omitempty"`
	GroupBy      string                        `json:"groupBy"`
	Page         int                           `json:"page"`
	PageSize     int                           `json:"pageSize"`
	TotalGroups  int64                         `json:"totalGroups"`
	TotalTokens  int64                         `json:"totalTokens"`
	TotalRecords int64                         `json:"totalRecords"`
	Items        []QuotaConsumptionSummaryItem `json:"items"`
}

// QuotaConsumptionReconciliationItem 表示某个聚合分组下的跨服务对账结果。
type QuotaConsumptionReconciliationItem struct {
	GroupKey         string  `json:"groupKey"`
	BackendTokens    int64   `json:"backendTokens"`
	BackendRecords   int64   `json:"backendRecords"`
	AIServiceTokens  int64   `json:"aiServiceTokens"`
	AIServiceRecords int64   `json:"aiServiceRecords"`
	DifferenceTokens int64   `json:"differenceTokens"`
	DifferenceRatio  float64 `json:"differenceRatio"`
	AlertLevel       string  `json:"alertLevel"`
	ShouldAlert      bool    `json:"shouldAlert"`
}

// QuotaConsumptionReconciliationSummary 表示全局消费聚合的跨服务对账结果。
type QuotaConsumptionReconciliationSummary struct {
	TimeRange             string                               `json:"timeRange"`
	WorkflowType          string                               `json:"workflowType,omitempty"`
	GroupBy               string                               `json:"groupBy"`
	Page                  int                                  `json:"page"`
	PageSize              int                                  `json:"pageSize"`
	TotalGroups           int64                                `json:"totalGroups"`
	BackendTotalTokens    int64                                `json:"backendTotalTokens"`
	BackendTotalRecords   int64                                `json:"backendTotalRecords"`
	AIServiceTotalTokens  int64                                `json:"aiServiceTotalTokens"`
	AIServiceTotalRecords int64                                `json:"aiServiceTotalRecords"`
	DifferenceTokens      int64                                `json:"differenceTokens"`
	DifferenceRatio       float64                              `json:"differenceRatio"`
	AlertLevel            string                               `json:"alertLevel"`
	ShouldAlert           bool                                 `json:"shouldAlert"`
	CheckedAt             time.Time                            `json:"checkedAt"`
	Items                 []QuotaConsumptionReconciliationItem `json:"items"`
}
