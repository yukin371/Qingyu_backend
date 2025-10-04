package admin

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	adminModel "Qingyu_backend/models/shared/admin"
)

// AdminServiceImpl 管理服务实现
type AdminServiceImpl struct {
	auditRepo AuditRepository
	logRepo   LogRepository
	userRepo  UserRepository
}

// AuditRepository 审核记录Repository
type AuditRepository interface {
	Create(ctx context.Context, record *AuditRecord) error
	Get(ctx context.Context, recordID string) (*AuditRecord, error)
	Update(ctx context.Context, recordID string, updates map[string]interface{}) error
	ListByStatus(ctx context.Context, contentType, status string) ([]*AuditRecord, error)
	ListByContent(ctx context.Context, contentID string) ([]*AuditRecord, error)
}

// LogRepository 日志Repository
type LogRepository interface {
	Create(ctx context.Context, log *AdminLog) error
	List(ctx context.Context, filter *LogFilter) ([]*AdminLog, error)
}

// UserRepository 用户Repository（简化版）
type UserRepository interface {
	GetStatistics(ctx context.Context, userID string) (*UserStatistics, error)
	BanUser(ctx context.Context, userID, reason string, until time.Time) error
	UnbanUser(ctx context.Context, userID string) error
}

// LogFilter 日志过滤条件
type LogFilter struct {
	AdminID   string
	Operation string
	StartDate time.Time
	EndDate   time.Time
	Page      int
	PageSize  int
}

// NewAdminService 创建管理服务
func NewAdminService(auditRepo AuditRepository, logRepo LogRepository, userRepo UserRepository) AdminService {
	return &AdminServiceImpl{
		auditRepo: auditRepo,
		logRepo:   logRepo,
		userRepo:  userRepo,
	}
}

// ============ 内容审核 ============

// ReviewContent 审核内容
func (s *AdminServiceImpl) ReviewContent(ctx context.Context, req *ReviewContentRequest) error {
	// 1. 验证操作
	if req.Action != "approve" && req.Action != "reject" {
		return fmt.Errorf("无效的审核操作: %s", req.Action)
	}

	// 2. 获取或创建审核记录
	records, err := s.auditRepo.ListByContent(ctx, req.ContentID)
	if err != nil {
		return fmt.Errorf("查询审核记录失败: %w", err)
	}

	var record *AuditRecord
	if len(records) > 0 {
		// 使用最新的记录
		record = records[0]
	} else {
		// 创建新记录
		record = &AuditRecord{
			ContentID:   req.ContentID,
			ContentType: req.ContentType,
			Status:      adminModel.AuditStatusPending,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := s.auditRepo.Create(ctx, record); err != nil {
			return fmt.Errorf("创建审核记录失败: %w", err)
		}
	}

	// 3. 更新审核状态
	status := adminModel.AuditStatusApproved
	if req.Action == "reject" {
		status = adminModel.AuditStatusRejected
	}

	updates := map[string]interface{}{
		"status":      status,
		"reviewer_id": req.ReviewerID,
		"reason":      req.Reason,
		"reviewed_at": time.Now(),
		"updated_at":  time.Now(),
	}

	if err := s.auditRepo.Update(ctx, record.ID, updates); err != nil {
		return fmt.Errorf("更新审核记录失败: %w", err)
	}

	// 4. 记录操作日志
	s.LogOperation(ctx, &LogOperationRequest{
		AdminID:   req.ReviewerID,
		Operation: adminModel.OperationReviewContent,
		Target:    req.ContentID,
		Details: map[string]interface{}{
			"content_type": req.ContentType,
			"action":       req.Action,
			"reason":       req.Reason,
		},
	})

	return nil
}

// GetPendingReviews 获取待审核内容
func (s *AdminServiceImpl) GetPendingReviews(ctx context.Context, contentType string) ([]*AuditRecord, error) {
	return s.auditRepo.ListByStatus(ctx, contentType, adminModel.AuditStatusPending)
}

// ============ 用户管理 ============

// ManageUser 管理用户
func (s *AdminServiceImpl) ManageUser(ctx context.Context, req *ManageUserRequest) error {
	switch req.Action {
	case "ban":
		duration := time.Duration(req.Duration) * time.Second
		return s.BanUser(ctx, req.UserID, req.Reason, duration)
	case "unban":
		return s.UnbanUser(ctx, req.UserID)
	case "delete":
		// 删除用户操作（这里简化处理）
		s.LogOperation(ctx, &LogOperationRequest{
			AdminID:   req.AdminID,
			Operation: adminModel.OperationDeleteUser,
			Target:    req.UserID,
			Details: map[string]interface{}{
				"reason": req.Reason,
			},
		})
		return nil
	default:
		return fmt.Errorf("无效的用户管理操作: %s", req.Action)
	}
}

// BanUser 封禁用户
func (s *AdminServiceImpl) BanUser(ctx context.Context, userID, reason string, duration time.Duration) error {
	// 1. 计算封禁到期时间
	until := time.Now().Add(duration)
	if duration == 0 {
		// 永久封禁
		until = time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)
	}

	// 2. 执行封禁
	if err := s.userRepo.BanUser(ctx, userID, reason, until); err != nil {
		return fmt.Errorf("封禁用户失败: %w", err)
	}

	// 3. 记录日志
	s.LogOperation(ctx, &LogOperationRequest{
		AdminID:   "system", // 可以从context中获取
		Operation: adminModel.OperationBanUser,
		Target:    userID,
		Details: map[string]interface{}{
			"reason":   reason,
			"duration": duration.String(),
			"until":    until.Format(time.RFC3339),
		},
	})

	return nil
}

// UnbanUser 解封用户
func (s *AdminServiceImpl) UnbanUser(ctx context.Context, userID string) error {
	// 1. 执行解封
	if err := s.userRepo.UnbanUser(ctx, userID); err != nil {
		return fmt.Errorf("解封用户失败: %w", err)
	}

	// 2. 记录日志
	s.LogOperation(ctx, &LogOperationRequest{
		AdminID:   "system",
		Operation: adminModel.OperationUnbanUser,
		Target:    userID,
	})

	return nil
}

// GetUserStatistics 获取用户统计信息
func (s *AdminServiceImpl) GetUserStatistics(ctx context.Context, userID string) (*UserStatistics, error) {
	stats, err := s.userRepo.GetStatistics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户统计失败: %w", err)
	}
	return stats, nil
}

// ============ 提现审核 ============

// ReviewWithdraw 审核提现
func (s *AdminServiceImpl) ReviewWithdraw(ctx context.Context, withdrawID, adminID string, approved bool, reason string) error {
	// 记录操作日志
	operation := adminModel.OperationApproveWithdraw
	if !approved {
		operation = adminModel.OperationRejectWithdraw
	}

	s.LogOperation(ctx, &LogOperationRequest{
		AdminID:   adminID,
		Operation: operation,
		Target:    withdrawID,
		Details: map[string]interface{}{
			"approved": approved,
			"reason":   reason,
		},
	})

	// 实际的提现审核逻辑应该调用Wallet服务
	// 这里简化处理
	return nil
}

// ============ 操作日志 ============

// LogOperation 记录操作日志
func (s *AdminServiceImpl) LogOperation(ctx context.Context, req *LogOperationRequest) error {
	log := &AdminLog{
		AdminID:   req.AdminID,
		Operation: req.Operation,
		Target:    req.Target,
		Details:   req.Details,
		IP:        req.IP,
		CreatedAt: time.Now(),
	}

	if err := s.logRepo.Create(ctx, log); err != nil {
		return fmt.Errorf("记录操作日志失败: %w", err)
	}

	return nil
}

// GetOperationLogs 获取操作日志
func (s *AdminServiceImpl) GetOperationLogs(ctx context.Context, req *GetLogsRequest) ([]*AdminLog, error) {
	// 设置默认分页
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}

	filter := &LogFilter{
		AdminID:   req.AdminID,
		Operation: req.Operation,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Page:      page,
		PageSize:  pageSize,
	}

	logs, err := s.logRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询操作日志失败: %w", err)
	}

	return logs, nil
}

// ExportLogs 导出日志（CSV格式）
func (s *AdminServiceImpl) ExportLogs(ctx context.Context, startDate, endDate time.Time) ([]byte, error) {
	// 1. 查询日志
	filter := &LogFilter{
		StartDate: startDate,
		EndDate:   endDate,
		Page:      1,
		PageSize:  10000, // 导出限制
	}

	logs, err := s.logRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询日志失败: %w", err)
	}

	// 2. 生成CSV
	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	// 写入表头
	headers := []string{"时间", "管理员ID", "操作", "目标", "IP", "详情"}
	writer.Write(headers)

	// 写入数据
	for _, log := range logs {
		details := fmt.Sprintf("%v", log.Details)
		record := []string{
			log.CreatedAt.Format("2006-01-02 15:04:05"),
			log.AdminID,
			log.Operation,
			log.Target,
			log.IP,
			details,
		}
		writer.Write(record)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("生成CSV失败: %w", err)
	}

	return []byte(builder.String()), nil
}

// ============ 健康检查 ============

// Health 健康检查
func (s *AdminServiceImpl) Health(ctx context.Context) error {
	// 简单的健康检查，可以扩展检查各个依赖
	return nil
}
