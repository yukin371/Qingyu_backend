package ai

import (
	"context"
	"errors"
	"fmt"
	"time"

	aiModels "Qingyu_backend/models/ai"
	pb "Qingyu_backend/pkg/grpc/pb"
	aiInterfaces "Qingyu_backend/repository/interfaces/ai"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	Total   int      `json:"total"`
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}

// QuotaAdminService 配额管理服务
type QuotaAdminService struct {
	quotaRepo         aiInterfaces.QuotaRepository
	alertRepo         aiInterfaces.QuotaAlertRepository
	consumptionReader quotaConsumptionReader
}

type quotaConsumptionReader interface {
	GetQuotaConsumption(ctx context.Context, userID string, timeRange string, workflowType string) (*pb.QuotaConsumptionResponse, error)
}

// NewQuotaAdminService 创建配额管理服务
func NewQuotaAdminService(quotaRepo aiInterfaces.QuotaRepository, alertRepo aiInterfaces.QuotaAlertRepository) *QuotaAdminService {
	return &QuotaAdminService{
		quotaRepo: quotaRepo,
		alertRepo: alertRepo,
	}
}

// SetConsumptionReader 设置 AI 服务对账读取客户端。
func (s *QuotaAdminService) SetConsumptionReader(reader quotaConsumptionReader) {
	s.consumptionReader = reader
}

// RechargeUserQuota 为用户充值配额
func (s *QuotaAdminService) RechargeUserQuota(ctx context.Context, userID string, amount int, quotaType string, reason, operatorID string) error {
	if amount <= 0 {
		return fmt.Errorf("充值金额必须大于0")
	}

	// 验证配额类型
	var quotaTypeEnum aiModels.QuotaType
	switch quotaType {
	case "daily":
		quotaTypeEnum = aiModels.QuotaTypeDaily
	case "monthly":
		quotaTypeEnum = aiModels.QuotaTypeMonthly
	default:
		return fmt.Errorf("不支持的配额类型: %s", quotaType)
	}

	// 获取用户配额
	quota, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, quotaTypeEnum)
	if err != nil && err != aiModels.ErrQuotaNotFound {
		return fmt.Errorf("获取配额失败: %w", err)
	}

	// 如果配额不存在，创建新配额
	if err == aiModels.ErrQuotaNotFound {
		now := time.Now()
		switch quotaTypeEnum {
		case aiModels.QuotaTypeDaily:
			quota = &aiModels.UserQuota{
				UserID:         userID,
				QuotaType:      quotaTypeEnum,
				TotalQuota:     amount,
				UsedQuota:      0,
				RemainingQuota: amount,
				Status:         aiModels.QuotaStatusActive,
				ResetAt:        now.AddDate(0, 0, 1),
			}
		case aiModels.QuotaTypeMonthly:
			quota = &aiModels.UserQuota{
				UserID:         userID,
				QuotaType:      quotaTypeEnum,
				TotalQuota:     amount,
				UsedQuota:      0,
				RemainingQuota: amount,
				Status:         aiModels.QuotaStatusActive,
				ResetAt:        now.AddDate(0, 1, 0),
			}
		}
		setQuotaManualOverride(quota, true)

		if err := s.quotaRepo.CreateQuota(ctx, quota); err != nil {
			return fmt.Errorf("创建配额失败: %w", err)
		}
	} else {
		// 更新配额
		beforeBalance := quota.RemainingQuota
		quota.TotalQuota += amount
		quota.RemainingQuota += amount
		setQuotaManualOverride(quota, true)

		if err := s.quotaRepo.UpdateQuota(ctx, quota); err != nil {
			return fmt.Errorf("更新配额失败: %w", err)
		}

		// 创建充值记录
		transaction := &aiModels.QuotaTransaction{
			ID:            primitive.NewObjectID(),
			UserID:        userID,
			QuotaType:     quotaTypeEnum,
			Amount:        amount,
			Type:          "recharge",
			Service:       "admin",
			Description:   reason,
			BeforeBalance: beforeBalance,
			AfterBalance:  quota.RemainingQuota,
			Timestamp:     time.Now(),
		}

		if err := s.quotaRepo.CreateTransaction(ctx, transaction); err != nil {
			zap.L().Warn("创建充值记录失败",
				zap.String("user_id", userID),
				zap.String("quota_type", string(quotaTypeEnum)),
				zap.Error(err),
			)
		}
	}

	// 创建告警记录（可选）
	alert := &aiModels.QuotaAlert{
		Type:    aiModels.QuotaAlertTypeConsistency,
		UserID:  userID,
		Level:   aiModels.QuotaAlertLevelInfo,
		Title:   "配额充值",
		Message: fmt.Sprintf("管理员为用户充值%s配额%d个，原因：%s", quotaType, amount, reason),
		Data: map[string]interface{}{
			"operatorID": operatorID,
			"amount":     amount,
			"quotaType":  quotaType,
		},
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		zap.L().Warn("创建配额充值告警记录失败",
			zap.String("user_id", userID),
			zap.Error(err),
		)
	}

	return nil
}

// BatchRecharge 批量充值
func (s *QuotaAdminService) BatchRecharge(ctx context.Context, userIDs []string, amount int, quotaType, reason, operatorID string) (*BatchOperationResult, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("用户列表不能为空")
	}

	if amount <= 0 {
		return nil, errors.New("充值金额必须大于0")
	}

	result := &BatchOperationResult{
		Total: len(userIDs),
	}

	// 并发执行充值操作
	for _, userID := range userIDs {
		err := s.RechargeUserQuota(ctx, userID, amount, quotaType, reason, operatorID)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("用户%s充值失败: %v", userID, err))
		} else {
			result.Success++
		}
	}

	return result, nil
}

// BatchUpdateQuota 批量更新配额
func (s *QuotaAdminService) BatchUpdateQuota(ctx context.Context, userIDs []string, totalQuota int, quotaType string) (*BatchOperationResult, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("用户列表不能为空")
	}

	if totalQuota < 0 {
		return nil, errors.New("总配额不能为负数")
	}

	result := &BatchOperationResult{
		Total: len(userIDs),
	}

	var quotaTypeEnum aiModels.QuotaType
	switch quotaType {
	case "daily":
		quotaTypeEnum = aiModels.QuotaTypeDaily
	case "monthly":
		quotaTypeEnum = aiModels.QuotaTypeMonthly
	default:
		return nil, fmt.Errorf("不支持的配额类型: %s", quotaType)
	}

	for _, userID := range userIDs {
		err := s.updateUserQuota(ctx, userID, totalQuota, quotaTypeEnum)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("用户%s更新配额失败: %v", userID, err))
		} else {
			result.Success++
		}
	}

	return result, nil
}

// BatchSuspend 批量暂停用户配额
func (s *QuotaAdminService) BatchSuspend(ctx context.Context, userIDs []string) (*BatchOperationResult, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("用户列表不能为空")
	}

	result := &BatchOperationResult{
		Total: len(userIDs),
	}

	for _, userID := range userIDs {
		_, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, aiModels.QuotaTypeDaily)
		if err != nil && err != aiModels.ErrQuotaNotFound {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("用户%s获取配额失败: %v", userID, err))
			continue
		}

		err = s.SuspendUserQuota(ctx, userID)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("用户%s暂停失败: %v", userID, err))
		} else {
			result.Success++
		}
	}

	return result, nil
}

// BatchActivate 批量激活用户配额
func (s *QuotaAdminService) BatchActivate(ctx context.Context, userIDs []string) (*BatchOperationResult, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("用户列表不能为空")
	}

	result := &BatchOperationResult{
		Total: len(userIDs),
	}

	for _, userID := range userIDs {
		_, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, aiModels.QuotaTypeDaily)
		if err != nil && err != aiModels.ErrQuotaNotFound {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("用户%s获取配额失败: %v", userID, err))
			continue
		}

		err = s.ActivateUserQuota(ctx, userID)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("用户%s激活失败: %v", userID, err))
		} else {
			result.Success++
		}
	}

	return result, nil
}

// ListUserQuotas 列出用户配额列表
func (s *QuotaAdminService) ListUserQuotas(ctx context.Context, role, status, search string, page, limit int) ([]*aiModels.UserQuotaListItem, int64, error) {
	return s.quotaRepo.ListUserQuotas(ctx, role, status, search, page, limit)
}

// GetUserQuotaDetail 获取用户所有配额详情
func (s *QuotaAdminService) GetUserQuotaDetail(ctx context.Context, userID string) (map[string]*aiModels.UserQuota, error) {
	quotas, err := s.quotaRepo.GetAllQuotasByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*aiModels.UserQuota)
	for _, quota := range quotas {
		key := string(quota.QuotaType)
		result[key] = quota
	}

	return result, nil
}

// GetUserQuotaReconciliation 获取单用户配额对账结果。
func (s *QuotaAdminService) GetUserQuotaReconciliation(
	ctx context.Context,
	userID string,
	timeRange string,
	workflowType string,
) (*aiModels.UserQuotaReconciliation, error) {
	if userID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if s.consumptionReader == nil {
		return nil, errors.New("AI配额对账客户端未配置")
	}

	normalizedRange, windowStart, windowEnd, err := resolveQuotaReconciliationWindow(timeRange, time.Now())
	if err != nil {
		return nil, err
	}

	transactions, err := s.quotaRepo.GetTransactionsByTimeRange(ctx, userID, windowStart, windowEnd)
	if err != nil {
		return nil, fmt.Errorf("获取后端配额事务失败: %w", err)
	}

	backendTotal := 0
	backendRecordCount := 0
	for _, tx := range transactions {
		if tx == nil || tx.Type != "consume" {
			continue
		}
		if workflowType != "" && tx.Service != workflowType {
			continue
		}
		backendTotal += tx.Amount
		backendRecordCount++
	}

	resp, err := s.consumptionReader.GetQuotaConsumption(ctx, userID, normalizedRange, workflowType)
	if err != nil {
		return nil, fmt.Errorf("获取AI服务配额消费失败: %w", err)
	}
	if !resp.GetSuccess() {
		message := resp.GetErrorMessage()
		if message == "" {
			message = resp.GetMessage()
		}
		if message == "" {
			message = "AI服务返回未成功状态"
		}
		return nil, errors.New(message)
	}

	level, shouldAlert := determineConsistencyAlertLevel(backendTotal, int(resp.GetTotalTokens()))
	diff := absInt(backendTotal - int(resp.GetTotalTokens()))
	diffRatio := calculateDifferenceRatioInt64(int64(backendTotal), int64(resp.GetTotalTokens()))

	records := make([]aiModels.QuotaReconciliationRecord, 0, len(resp.GetRecords()))
	for _, record := range resp.GetRecords() {
		records = append(records, aiModels.QuotaReconciliationRecord{
			ID:           record.GetId(),
			WorkflowType: record.GetWorkflowType(),
			TokensUsed:   int(record.GetTokensUsed()),
			ConsumedAt:   record.GetConsumedAt(),
		})
	}

	return &aiModels.UserQuotaReconciliation{
		UserID:               userID,
		TimeRange:            normalizedRange,
		WorkflowType:         workflowType,
		BackendQuotaType:     string(aiModels.QuotaTypeDaily),
		BackendTotalTokens:   backendTotal,
		BackendRecordCount:   backendRecordCount,
		AIServiceTotalTokens: int(resp.GetTotalTokens()),
		AIServiceRecordCount: int(resp.GetTotalRecords()),
		DifferenceTokens:     diff,
		DifferenceRatio:      diffRatio,
		AlertLevel:           string(level),
		ShouldAlert:          shouldAlert,
		WindowStartAt:        windowStart,
		WindowEndAt:          windowEnd,
		CheckedAt:            time.Now(),
		Records:              records,
	}, nil
}

// 私有辅助方法
func (s *QuotaAdminService) updateUserQuota(ctx context.Context, userID string, totalQuota int, quotaType aiModels.QuotaType) error {
	quota, err := s.quotaRepo.GetQuotaByUserID(ctx, userID, quotaType)
	if err != nil {
		if err == aiModels.ErrQuotaNotFound {
			// 创建新配额
			now := time.Now()
			switch quotaType {
			case aiModels.QuotaTypeDaily:
				quota = &aiModels.UserQuota{
					UserID:         userID,
					QuotaType:      quotaType,
					TotalQuota:     totalQuota,
					UsedQuota:      0,
					RemainingQuota: totalQuota,
					Status:         aiModels.QuotaStatusActive,
					ResetAt:        now.AddDate(0, 0, 1),
				}
			case aiModels.QuotaTypeMonthly:
				quota = &aiModels.UserQuota{
					UserID:         userID,
					QuotaType:      quotaType,
					TotalQuota:     totalQuota,
					UsedQuota:      0,
					RemainingQuota: totalQuota,
					Status:         aiModels.QuotaStatusActive,
					ResetAt:        now.AddDate(0, 1, 0),
				}
			}
			setQuotaManualOverride(quota, true)

			return s.quotaRepo.CreateQuota(ctx, quota)
		}
		return err
	}

	// 更新配额
	quota.TotalQuota = totalQuota
	quota.RemainingQuota = totalQuota - quota.UsedQuota
	setQuotaManualOverride(quota, true)

	return s.quotaRepo.UpdateQuota(ctx, quota)
}

// SuspendUserQuota 暂停用户配额
func (s *QuotaAdminService) SuspendUserQuota(ctx context.Context, userID string) error {
	quotas, err := s.quotaRepo.GetAllQuotasByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, quota := range quotas {
		setQuotaManualOverride(quota, true)
		quota.Status = aiModels.QuotaStatusSuspended
		if err := s.quotaRepo.UpdateQuota(ctx, quota); err != nil {
			return err
		}
	}

	// 创建告警
	alert := &aiModels.QuotaAlert{
		Type:    aiModels.QuotaAlertTypeConsistency,
		UserID:  userID,
		Level:   aiModels.QuotaAlertLevelWarning,
		Title:   "配额暂停",
		Message: fmt.Sprintf("管理员暂停了用户的配额"),
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		zap.L().Warn("创建配额暂停告警记录失败",
			zap.String("user_id", userID),
			zap.Error(err),
		)
	}

	return nil
}

// ActivateUserQuota 激活用户配额
func (s *QuotaAdminService) ActivateUserQuota(ctx context.Context, userID string) error {
	quotas, err := s.quotaRepo.GetAllQuotasByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, quota := range quotas {
		setQuotaManualOverride(quota, true)
		quota.Status = aiModels.QuotaStatusActive
		if err := s.quotaRepo.UpdateQuota(ctx, quota); err != nil {
			return err
		}
	}

	// 创建告警
	alert := &aiModels.QuotaAlert{
		Type:    aiModels.QuotaAlertTypeConsistency,
		UserID:  userID,
		Level:   aiModels.QuotaAlertLevelInfo,
		Title:   "配额激活",
		Message: fmt.Sprintf("管理员激活了用户的配额"),
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		zap.L().Warn("创建配额激活告警记录失败",
			zap.String("user_id", userID),
			zap.Error(err),
		)
	}

	return nil
}
