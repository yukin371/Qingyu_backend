package admin

import (
	adminModel "Qingyu_backend/models/users"
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"time"
)

// TODO: 完善审核流程
//   - 多级审核支持
//   - 审核规则引擎
//   - 自动审核（敏感词过滤、AI内容审核）
// TODO: 完善日志记录功能
//   - 操作日志详细记录
//   - 日志搜索和分析
//   - 日志归档和备份
// TODO: 添加权限管理功能
// TODO: 实现数据统计和报表导出

// AdminServiceImpl 管理服务实现
type AdminServiceImpl struct {
	auditRepo   AuditRepository
	logRepo     LogRepository
	userRepo    UserRepository
	initialized bool // 初始化标志
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
	_ = s.LogOperation(ctx, &LogOperationRequest{
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
		_ = s.LogOperation(ctx, &LogOperationRequest{
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
	_ = s.LogOperation(ctx, &LogOperationRequest{
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
	_ = s.LogOperation(ctx, &LogOperationRequest{
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

	_ = s.LogOperation(ctx, &LogOperationRequest{
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
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("error writing CSV headers: %w", err)
	}

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
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("error writing CSV record: %w", err)
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("生成CSV失败: %w", err)
	}

	return []byte(builder.String()), nil
}

// ============ 健康检查 ============

// GetSystemStats 获取系统统计
func (s *AdminServiceImpl) GetSystemStats(ctx context.Context) (interface{}, error) {
	// 返回系统统计数据（临时实现，后续可连接Repository获取真实数据）
	stats := map[string]interface{}{
		"totalUsers":    0,
		"activeUsers":   0,
		"totalBooks":    0,
		"totalRevenue":  0.0,
		"pendingAudits": 0,
	}
	return stats, nil
}

// GetSystemConfig 获取系统配置
func (s *AdminServiceImpl) GetSystemConfig(ctx context.Context) (interface{}, error) {
	// 返回系统配置（临时实现，后续可连接ConfigService获取真实数据）
	config := map[string]interface{}{
		"allowRegistration":        true,
		"requireEmailVerification": true,
		"maxUploadSize":            10485760,
		"enableAudit":              true,
	}
	return config, nil
}

// UpdateSystemConfig 更新系统配置
func (s *AdminServiceImpl) UpdateSystemConfig(ctx context.Context, req interface{}) error {
	// 后续可连接ConfigService实现配置持久化
	return nil
}

// CreateAnnouncement 创建公告
func (s *AdminServiceImpl) CreateAnnouncement(ctx context.Context, adminID, title, content, announceType, priority string) (interface{}, error) {
	// 后续可连接AnnouncementService实现公告创建
	announcement := map[string]interface{}{
		"id":        "new-announcement-id",
		"title":     title,
		"content":   content,
		"type":      announceType,
		"priority":  priority,
		"createdBy": adminID,
		"createdAt": time.Now(),
	}
	return announcement, nil
}

// GetAnnouncements 获取公告列表
func (s *AdminServiceImpl) GetAnnouncements(ctx context.Context, page, pageSize int) (interface{}, int64, error) {
	// 设置默认分页
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	// TODO: 实现分页查询时使用这些参数
	_ = page     // 显式标记为有意未使用（待实现）
	_ = pageSize // 显式标记为有意未使用（待实现）

	// 后续可连接AnnouncementService实现公告列表查询
	announcements := []interface{}{}
	return announcements, 0, nil
}

// GetAuditStatistics 获取审核统计
func (s *AdminServiceImpl) GetAuditStatistics(ctx context.Context) (interface{}, error) {
	// 返回审核统计数据（临时实现，后续可连接Repository获取真实数据）
	stats := map[string]interface{}{
		"pending":     0,
		"approved":    0,
		"rejected":    0,
		"highRisk":    0,
		"approveRate": 0.0,
	}
	return stats, nil
}

// ============ BaseService 接口实现 ============

// Initialize 初始化服务
func (s *AdminServiceImpl) Initialize(ctx context.Context) error {
	if s.initialized {
		return nil
	}

	// 验证依赖项
	if s.auditRepo == nil {
		return fmt.Errorf("auditRepo is nil")
	}
	if s.logRepo == nil {
		return fmt.Errorf("logRepo is nil")
	}
	if s.userRepo == nil {
		return fmt.Errorf("userRepo is nil")
	}

	// 初始化完成
	s.initialized = true
	return nil
}

// Health 健康检查
func (s *AdminServiceImpl) Health(ctx context.Context) error {
	if !s.initialized {
		return fmt.Errorf("service not initialized")
	}

	// 可以扩展检查各个依赖的健康状态
	return nil
}

// Close 关闭服务，清理资源
func (s *AdminServiceImpl) Close(ctx context.Context) error {
	// 管理服务暂无需要清理的资源
	s.initialized = false
	return nil
}

// GetServiceName 获取服务名称
func (s *AdminServiceImpl) GetServiceName() string {
	return "AdminService"
}

// GetVersion 获取服务版本
func (s *AdminServiceImpl) GetVersion() string {
	return "v1.0.0"
}
