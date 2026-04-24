package ai

import "time"

// QuotaReconciliationRecord 表示 AI 服务侧返回的消费记录摘要。
type QuotaReconciliationRecord struct {
	ID           string `json:"id"`
	WorkflowType string `json:"workflowType"`
	TokensUsed   int    `json:"tokensUsed"`
	ConsumedAt   string `json:"consumedAt"`
}

// UserQuotaReconciliation 表示单用户的 backend 与 AI 服务配额对账结果。
type UserQuotaReconciliation struct {
	UserID               string                      `json:"userId"`
	TimeRange            string                      `json:"timeRange"`
	WorkflowType         string                      `json:"workflowType,omitempty"`
	BackendQuotaType     string                      `json:"backendQuotaType"`
	BackendTotalTokens   int                         `json:"backendTotalTokens"`
	BackendRecordCount   int                         `json:"backendRecordCount"`
	AIServiceTotalTokens int                         `json:"aiServiceTotalTokens"`
	AIServiceRecordCount int                         `json:"aiServiceRecordCount"`
	DifferenceTokens     int                         `json:"differenceTokens"`
	DifferenceRatio      float64                     `json:"differenceRatio"`
	AlertLevel           string                      `json:"alertLevel"`
	ShouldAlert          bool                        `json:"shouldAlert"`
	WindowStartAt        time.Time                   `json:"windowStartAt"`
	WindowEndAt          time.Time                   `json:"windowEndAt"`
	CheckedAt            time.Time                   `json:"checkedAt"`
	Records              []QuotaReconciliationRecord `json:"records,omitempty"`
}
