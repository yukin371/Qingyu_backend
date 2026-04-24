package ai

import (
	"context"
	"fmt"
	"time"

	aiModels "Qingyu_backend/models/ai"
	aiInterfaces "Qingyu_backend/repository/interfaces/ai"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuotaAlertService 配额告警服务
type QuotaAlertService struct {
	alertRepo aiInterfaces.QuotaAlertRepository
}

// NewQuotaAlertService 创建配额告警服务
func NewQuotaAlertService(alertRepo aiInterfaces.QuotaAlertRepository) *QuotaAlertService {
	return &QuotaAlertService{
		alertRepo: alertRepo,
	}
}

// CreateAlert 创建告警
func (s *QuotaAlertService) CreateAlert(ctx context.Context, alert *aiModels.QuotaAlert) error {
	if alert == nil {
		return fmt.Errorf("告警不能为空")
	}

	// 设置默认值
	if alert.Level == "" {
		alert.Level = aiModels.QuotaAlertLevelInfo
	}
	if alert.Status == "" {
		alert.Status = aiModels.QuotaAlertStatusPending
	}

	return s.alertRepo.Create(ctx, alert)
}

// GetAlert 根据ID获取告警
func (s *QuotaAlertService) GetAlert(ctx context.Context, id string) (*aiModels.QuotaAlert, error) {
	return s.alertRepo.GetByID(ctx, id)
}

// ListAlerts 获取告警列表
func (s *QuotaAlertService) ListAlerts(ctx context.Context, alertType, level, status string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
	return s.alertRepo.List(ctx, alertType, level, status, page, limit)
}

// AcknowledgeAlert 确认告警
func (s *QuotaAlertService) AcknowledgeAlert(ctx context.Context, id, operatorID string) error {
	alert, err := s.alertRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取告警失败: %w", err)
	}

	// 更新告警状态
	alert.Status = aiModels.QuotaAlertStatusAcknowledged
	alert.ResolvedBy = operatorID

	return s.alertRepo.Update(ctx, alert)
}

// ResolveAlert 解决告警
func (s *QuotaAlertService) ResolveAlert(ctx context.Context, id, operatorID string) error {
	alert, err := s.alertRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取告警失败: %w", err)
	}

	// 更新告警状态
	alert.Resolve(operatorID)

	return s.alertRepo.Update(ctx, alert)
}

// IgnoreAlert 忽略告警
func (s *QuotaAlertService) IgnoreAlert(ctx context.Context, id, operatorID string) error {
	alert, err := s.alertRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取告警失败: %w", err)
	}

	// 更新告警状态
	alert.Status = aiModels.QuotaAlertStatusIgnored
	alert.ResolvedBy = operatorID

	return s.alertRepo.Update(ctx, alert)
}

// GetUserAlerts 获取用户告警列表
func (s *QuotaAlertService) GetUserAlerts(ctx context.Context, userID string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
	return s.alertRepo.GetByUserID(ctx, userID, page, limit)
}

// GetRecentAlerts 获取最近的告警
func (s *QuotaAlertService) GetRecentAlerts(ctx context.Context, limit int) ([]*aiModels.QuotaAlert, error) {
	return s.alertRepo.GetRecentGlobal(ctx, limit)
}

// CountPendingAlerts 统计待处理告警数量
func (s *QuotaAlertService) CountPendingAlerts(ctx context.Context) (map[aiModels.QuotaAlertStatus]int64, error) {
	return s.alertRepo.CountByStatus(ctx)
}

// ResolveRecoveredConsistencyAlerts 将本轮已复核且已恢复一致的 consistency 告警自动标记为 resolved。
func (s *QuotaAlertService) ResolveRecoveredConsistencyAlerts(
	ctx context.Context,
	checkedKeys map[string]struct{},
	activeKeys map[string]struct{},
	operatorID string,
) error {
	if len(checkedKeys) == 0 {
		return nil
	}
	if operatorID == "" {
		operatorID = "system"
	}

	alerts, _, err := s.alertRepo.List(ctx, string(aiModels.QuotaAlertTypeConsistency), "", "", 1, 1000)
	if err != nil {
		return fmt.Errorf("查询 consistency 告警失败: %w", err)
	}

	for _, alert := range alerts {
		if alert == nil || !isOpenConsistencyAlertStatus(alert.Status) {
			continue
		}

		alertKey := buildConsistencyAlertKey(alert.UserID, alert.Data)
		if _, shouldSync := checkedKeys[alertKey]; !shouldSync {
			continue
		}
		if _, stillActive := activeKeys[alertKey]; stillActive {
			continue
		}

		alert.Resolve(operatorID)
		if err := s.alertRepo.Update(ctx, alert); err != nil {
			return fmt.Errorf("更新 recovered consistency 告警失败: %w", err)
		}
	}

	return nil
}

// AlertManager 告警管理器
type AlertManager struct {
	alertService *QuotaAlertService
}

// NewAlertManager 创建告警管理器
func NewAlertManager(alertService *QuotaAlertService) *AlertManager {
	return &AlertManager{
		alertService: alertService,
	}
}

// ProcessAlerts 处理待处理告警
func (am *AlertManager) ProcessAlerts(ctx context.Context, threshold time.Duration) error {
	// 获取待处理告警
	alerts, _, err := am.alertService.ListAlerts(ctx, "", "", "pending", 1, 1000)
	if err != nil {
		return fmt.Errorf("获取待处理告警失败: %w", err)
	}

	now := time.Now()

	// 处理超时的告警
	for _, alert := range alerts {
		if now.Sub(alert.CreatedAt) > threshold {
			// 自动降级为已确认
			err := am.alertService.AcknowledgeAlert(ctx, alert.ID.Hex(), "system")
			if err != nil {
				fmt.Printf("自动确认告警失败[%s]: %v\n", alert.ID.Hex(), err)
				continue
			}
			fmt.Printf("告警超时自动确认: %s\n", alert.ID.Hex())
		}
	}

	return nil
}

// AutoResolveResolvedAlerts 自动解决已解决的告警
func (am *AlertManager) AutoResolveResolvedAlerts(ctx context.Context, maxAge time.Duration) error {
	// 获取已解决的告警
	alerts, _, err := am.alertService.ListAlerts(ctx, "", "", "resolved", 1, 1000)
	if err != nil {
		return fmt.Errorf("获取已解决告警失败: %w", err)
	}

	now := time.Now()

	for _, alert := range alerts {
		if alert.ResolvedAt != nil && now.Sub(*alert.ResolvedAt) > maxAge {
			// 删除旧告警（软删除）
			err := am.alertService.IgnoreAlert(ctx, alert.ID.Hex(), "system")
			if err != nil {
				fmt.Printf("删除旧告警失败[%s]: %v\n", alert.ID.Hex(), err)
			}
		}
	}

	return nil
}

// 私有辅助方法

// CreateThresholdAlert 创建阈值告警
func (s *QuotaAlertService) CreateThresholdAlert(ctx context.Context, userID, level, title, message string, data map[string]interface{}) error {
	alert := &aiModels.QuotaAlert{
		Type:    aiModels.QuotaAlertTypeThreshold,
		UserID:  userID,
		Level:   aiModels.QuotaAlertLevel(level),
		Title:   title,
		Message: message,
		Data:    data,
	}

	return s.CreateAlert(ctx, alert)
}

// CreateAnomalyAlert 创建异常告警
func (s *QuotaAlertService) CreateAnomalyAlert(ctx context.Context, userID, level, title, message string, data map[string]interface{}) error {
	alert := &aiModels.QuotaAlert{
		Type:    aiModels.QuotaAlertTypeAnomaly,
		UserID:  userID,
		Level:   aiModels.QuotaAlertLevel(level),
		Title:   title,
		Message: message,
		Data:    data,
	}

	return s.CreateAlert(ctx, alert)
}

// CreateAbuseAlert 创建滥用告警
func (s *QuotaAlertService) CreateAbuseAlert(ctx context.Context, userID, level, title, message string, data map[string]interface{}) error {
	alert := &aiModels.QuotaAlert{
		Type:    aiModels.QuotaAlertTypeAbuse,
		UserID:  userID,
		Level:   aiModels.QuotaAlertLevel(level),
		Title:   title,
		Message: message,
		Data:    data,
	}

	return s.CreateAlert(ctx, alert)
}

// CreateConsistencyAlert 创建一致性告警
func (s *QuotaAlertService) CreateConsistencyAlert(ctx context.Context, userID, level, title, message string, data map[string]interface{}) error {
	alert := &aiModels.QuotaAlert{
		Type:    aiModels.QuotaAlertTypeConsistency,
		UserID:  userID,
		Level:   aiModels.QuotaAlertLevel(level),
		Title:   title,
		Message: message,
		Data:    data,
	}

	return s.CreateAlert(ctx, alert)
}

// GenerateTransactionID 生成事务ID
func GenerateTransactionID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func isOpenConsistencyAlertStatus(status aiModels.QuotaAlertStatus) bool {
	return status == aiModels.QuotaAlertStatusPending || status == aiModels.QuotaAlertStatusAcknowledged
}
