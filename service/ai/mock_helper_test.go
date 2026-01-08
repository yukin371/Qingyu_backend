package ai

import (
	"time"

	"Qingyu_backend/models/ai"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ Mock辅助函数 ============

// NewTestContextService 创建用于测试的ContextService
func NewTestContextService() *ContextService {
	return &ContextService{
		documentService:     nil,
		projectService:      nil,
		nodeService:         nil,
		versionService:      nil,
		documentContentRepo: nil,
	}
}

// CreateTestUserQuota 创建测试用的用户配额
func CreateTestUserQuota(userID string) *ai.UserQuota {
	return &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      0,
		RemainingQuota: 1000,
		Status:         ai.QuotaStatusActive,
		ResetAt:        time.Now().AddDate(0, 0, 1),
		Metadata: &ai.QuotaMetadata{
			UserRole:        "reader",
			MembershipLevel: "normal",
		},
	}
}

// CreateTestQuotaTransaction 创建测试用的配额事务
func CreateTestQuotaTransaction(userID string) *ai.QuotaTransaction {
	return &ai.QuotaTransaction{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		QuotaType:     ai.QuotaTypeDaily,
		Amount:        100,
		Type:          "consume",
		Service:       "chat",
		Model:         "claude-3",
		RequestID:     "req-123",
		BeforeBalance: 1000,
		AfterBalance:  900,
		Timestamp:     time.Now(),
	}
}
