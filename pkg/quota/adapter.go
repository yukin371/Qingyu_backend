package quota

import (
	"context"

	aiService "Qingyu_backend/service/ai"
)

// QuotaServiceAdapter 将原有的QuotaService适配为Checker接口
// 这是Adapter模式的应用，使旧代码可以与新接口兼容
type QuotaServiceAdapter struct {
	service *aiService.QuotaService
}

// NewQuotaServiceAdapter 创建适配器
func NewQuotaServiceAdapter(service *aiService.QuotaService) Checker {
	return &QuotaServiceAdapter{service: service}
}

// Check 实现Checker接口
func (a *QuotaServiceAdapter) Check(ctx context.Context, userID string, amount int) error {
	return a.service.CheckQuota(ctx, userID, amount)
}
