package admin

import (
	financeModel "Qingyu_backend/models/finance"
	adminModel "Qingyu_backend/models/users"
	financeRepo "Qingyu_backend/repository/interfaces/finance"
	financeWalletService "Qingyu_backend/service/finance/wallet"
	"context"
	"encoding/csv"
	"fmt"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	auditRepo         AuditRepository
	logRepo           LogRepository
	userRepo          UserRepository
	walletService     financeWalletService.WalletService
	walletRepo        financeRepo.WalletRepository
	authorRevenueRepo financeRepo.AuthorRevenueRepository
	initialized       bool // 初始化标志
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
	GetBasicProfile(ctx context.Context, userID string) (*UserBasicProfile, error)
	// Dashboard 统计方法
	GetTotalCount(ctx context.Context) (int64, error)
	CountActiveUsers(ctx context.Context, days int) (int64, error)
	CountNewUsersToday(ctx context.Context) (int64, error)
	CountAuthors(ctx context.Context) (int64, error)
}

// UserBasicProfile 用户基础信息
type UserBasicProfile struct {
	UserID      string
	Username    string
	Email       string
	DisplayName string
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
func NewAdminService(
	auditRepo AuditRepository,
	logRepo LogRepository,
	userRepo UserRepository,
	walletService financeWalletService.WalletService,
	walletRepo financeRepo.WalletRepository,
	authorRevenueRepo financeRepo.AuthorRevenueRepository,
) AdminService {
	return &AdminServiceImpl{
		auditRepo:         auditRepo,
		logRepo:           logRepo,
		userRepo:          userRepo,
		walletService:     walletService,
		walletRepo:        walletRepo,
		authorRevenueRepo: authorRevenueRepo,
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
	if withdrawID == "" {
		return fmt.Errorf("提现单ID不能为空")
	}

	operation := adminModel.OperationApproveWithdraw
	if !approved {
		operation = adminModel.OperationRejectWithdraw
	}

	walletRequest, walletErr := s.walletService.GetWithdrawRequest(ctx, withdrawID)
	if walletErr == nil && walletRequest != nil {
		if approved {
			if err := s.walletService.ApproveWithdraw(ctx, withdrawID, adminID); err != nil {
				return fmt.Errorf("审核钱包提现失败: %w", err)
			}
		} else {
			if strings.TrimSpace(reason) == "" {
				return fmt.Errorf("拒绝提现时必须填写原因")
			}
			if err := s.walletService.RejectWithdraw(ctx, withdrawID, adminID, reason); err != nil {
				return fmt.Errorf("拒绝钱包提现失败: %w", err)
			}
		}

		_ = s.LogOperation(ctx, &LogOperationRequest{
			AdminID:   adminID,
			Operation: operation,
			Target:    withdrawID,
			Details: map[string]interface{}{
				"approved": approved,
				"reason":   reason,
				"source":   "wallet",
				"user_id":  walletRequest.UserID,
			},
		})
		return nil
	}

	objectID, err := primitive.ObjectIDFromHex(withdrawID)
	if err != nil {
		return fmt.Errorf("提现单ID格式无效")
	}

	authorRequest, err := s.authorRevenueRepo.GetWithdrawalRequest(ctx, objectID)
	if err != nil {
		if walletErr != nil {
			return fmt.Errorf("提现申请不存在")
		}
		return fmt.Errorf("查询作者提现失败: %w", err)
	}

	if authorRequest.Status != financeModel.WithdrawStatusPending {
		return fmt.Errorf("该提现申请已处理，当前状态: %s", authorRequest.Status)
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":      financeModel.WithdrawStatusApproved,
		"approved_by": adminID,
		"approved_at": now,
		"updated_at":  now,
	}
	if !approved {
		if strings.TrimSpace(reason) == "" {
			return fmt.Errorf("拒绝提现时必须填写原因")
		}
		updates["status"] = financeModel.WithdrawStatusRejected
		updates["reject_reason"] = reason
		updates["note"] = reason
	}

	if err := s.authorRevenueRepo.UpdateWithdrawalRequest(ctx, objectID, updates); err != nil {
		return fmt.Errorf("审核作者提现失败: %w", err)
	}

	_ = s.LogOperation(ctx, &LogOperationRequest{
		AdminID:   adminID,
		Operation: operation,
		Target:    withdrawID,
		Details: map[string]interface{}{
			"approved": approved,
			"reason":   reason,
			"source":   "author",
			"user_id":  authorRequest.UserID,
		},
	})
	return nil
}

// ListWithdrawals 获取提现列表
func (s *AdminServiceImpl) ListWithdrawals(ctx context.Context, req *ListWithdrawalsRequest) ([]*WithdrawalRecord, int64, error) {
	page, pageSize := normalizeWithdrawalPagination(req)

	records, err := s.collectWithdrawalRecords(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	total := int64(len(records))
	start := (page - 1) * pageSize
	if start >= len(records) {
		return []*WithdrawalRecord{}, total, nil
	}
	end := start + pageSize
	if end > len(records) {
		end = len(records)
	}

	return records[start:end], total, nil
}

// GetWithdrawalStatistics 获取提现统计
func (s *AdminServiceImpl) GetWithdrawalStatistics(ctx context.Context, req *WithdrawalStatisticsRequest) (*WithdrawalStatistics, error) {
	listReq := &ListWithdrawalsRequest{
		Page:      1,
		PageSize:  1000,
		Status:    req.Status,
		Source:    req.Source,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	records, err := s.collectWithdrawalRecords(ctx, listReq)
	if err != nil {
		return nil, err
	}

	stats := &WithdrawalStatistics{}
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for _, record := range records {
		stats.TotalCount++
		switch record.Status {
		case financeModel.WithdrawStatusPending:
			stats.PendingCount++
			stats.PendingAmount += record.Amount
		case financeModel.WithdrawStatusApproved, "processed":
			stats.ApprovedCount++
			stats.ApprovedAmount += record.Amount
			if record.ReviewedAt != nil && !record.ReviewedAt.Before(todayStart) {
				stats.ApprovedTodayCount++
			}
		case financeModel.WithdrawStatusRejected, financeModel.WithdrawStatusFailed:
			stats.RejectedCount++
		}
	}

	return stats, nil
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
	// 从 repository 获取统计数据
	totalUsers, err := s.userRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取总用户数失败: %w", err)
	}

	// 统计最近7天活跃用户
	activeUsers, err := s.userRepo.CountActiveUsers(ctx, 7)
	if err != nil {
		return nil, fmt.Errorf("获取活跃用户数失败: %w", err)
	}

	// 统计今日新增用户
	newUsersToday, err := s.userRepo.CountNewUsersToday(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取今日新增用户失败: %w", err)
	}

	// 统计作者数量
	authorsCount, err := s.userRepo.CountAuthors(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取作者数量失败: %w", err)
	}

	stats := map[string]interface{}{
		"totalUsers":    totalUsers,
		"activeUsers":   activeUsers,
		"totalBooks":    0,   // TODO: 需要从 BookService 获取
		"totalRevenue":  0.0, // TODO: 需要从 FinanceService 获取
		"pendingAudits": 0,   // TODO: 需要从 auditRepo 获取
		"newUsersToday": newUsersToday,
		"authorsCount":  authorsCount,
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
	if s.walletService == nil {
		return fmt.Errorf("walletService is nil")
	}
	if s.walletRepo == nil {
		return fmt.Errorf("walletRepo is nil")
	}
	if s.authorRevenueRepo == nil {
		return fmt.Errorf("authorRevenueRepo is nil")
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

func (s *AdminServiceImpl) collectWithdrawalRecords(ctx context.Context, req *ListWithdrawalsRequest) ([]*WithdrawalRecord, error) {
	result := make([]*WithdrawalRecord, 0)

	if req == nil || req.Source == "" || req.Source == "wallet" {
		walletRecords, err := s.listWalletWithdrawals(ctx, req)
		if err != nil {
			return nil, err
		}
		result = append(result, walletRecords...)
	}

	if req == nil || req.Source == "" || req.Source == "author" {
		authorRecords, err := s.listAuthorWithdrawals(ctx, req)
		if err != nil {
			return nil, err
		}
		result = append(result, authorRecords...)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

func (s *AdminServiceImpl) listWalletWithdrawals(ctx context.Context, req *ListWithdrawalsRequest) ([]*WithdrawalRecord, error) {
	filter := &financeRepo.WithdrawFilter{
		Limit: 1000,
	}
	if req != nil {
		filter.Status = req.Status
		filter.StartDate = req.StartDate
		filter.EndDate = req.EndDate
	}

	requests, err := s.walletRepo.ListWithdrawRequests(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询钱包提现列表失败: %w", err)
	}

	records := make([]*WithdrawalRecord, 0, len(requests))
	for _, item := range requests {
		profile, _ := s.userRepo.GetBasicProfile(ctx, item.UserID)
		reviewedAt := timePtrFromValue(item.ReviewedAt)
		processedAt := timePtrFromValue(item.ProcessedAt)
		records = append(records, &WithdrawalRecord{
			ID:            item.ID.Hex(),
			UserID:        item.UserID,
			Username:      safeProfileField(profile, func(p *UserBasicProfile) string { return p.Username }),
			Email:         safeProfileField(profile, func(p *UserBasicProfile) string { return p.Email }),
			DisplayName:   safeProfileField(profile, func(p *UserBasicProfile) string { return p.DisplayName }),
			Source:        "wallet",
			Status:        item.Status,
			Amount:        item.Amount.ToYuan(),
			Fee:           item.Fee.ToYuan(),
			ActualAmount:  item.ActualAmount.ToYuan(),
			Method:        item.AccountType,
			Account:       item.Account,
			AccountType:   item.AccountType,
			AccountName:   item.AccountName,
			OrderNo:       item.OrderNo,
			ReviewedBy:    item.ReviewedBy,
			ReviewedAt:    reviewedAt,
			RejectReason:  item.RejectReason,
			ProcessedAt:   processedAt,
			TransactionID: item.TransactionID,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
		})
	}

	return records, nil
}

func (s *AdminServiceImpl) listAuthorWithdrawals(ctx context.Context, req *ListWithdrawalsRequest) ([]*WithdrawalRecord, error) {
	filter := map[string]interface{}{}
	if req != nil {
		if req.Status != "" {
			filter["status"] = req.Status
		}
		if !req.StartDate.IsZero() || !req.EndDate.IsZero() {
			createdAt := map[string]interface{}{}
			if !req.StartDate.IsZero() {
				createdAt["$gte"] = req.StartDate
			}
			if !req.EndDate.IsZero() {
				createdAt["$lte"] = req.EndDate
			}
			filter["created_at"] = createdAt
		}
	}

	requests, _, err := s.authorRevenueRepo.ListWithdrawalRequests(ctx, filter, 1, 1000)
	if err != nil {
		return nil, fmt.Errorf("查询作者提现列表失败: %w", err)
	}

	records := make([]*WithdrawalRecord, 0, len(requests))
	for _, item := range requests {
		profile, _ := s.userRepo.GetBasicProfile(ctx, item.UserID)
		records = append(records, &WithdrawalRecord{
			ID:            item.ID.Hex(),
			UserID:        item.UserID,
			Username:      safeProfileField(profile, func(p *UserBasicProfile) string { return p.Username }),
			Email:         safeProfileField(profile, func(p *UserBasicProfile) string { return p.Email }),
			DisplayName:   safeProfileField(profile, func(p *UserBasicProfile) string { return p.DisplayName }),
			Source:        "author",
			Status:        item.Status,
			Amount:        item.Amount.ToYuan(),
			Fee:           item.Fee.ToYuan(),
			ActualAmount:  item.ActualAmount.ToYuan(),
			Method:        item.Method,
			Account:       item.AccountInfo.AccountNo,
			AccountType:   item.AccountInfo.AccountType,
			AccountName:   item.AccountInfo.AccountName,
			BankName:      item.AccountInfo.BankName,
			ReviewedBy:    item.ApprovedBy,
			ReviewedAt:    item.ApprovedAt,
			RejectReason:  item.RejectReason,
			ProcessedAt:   item.CompletedAt,
			TransactionID: item.TransactionID,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
		})
	}

	return records, nil
}

func normalizeWithdrawalPagination(req *ListWithdrawalsRequest) (int, int) {
	page := 1
	pageSize := 20
	if req == nil {
		return page, pageSize
	}
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func timePtrFromValue(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	copied := value
	return &copied
}

func safeProfileField(profile *UserBasicProfile, getter func(*UserBasicProfile) string) string {
	if profile == nil {
		return ""
	}
	return getter(profile)
}
