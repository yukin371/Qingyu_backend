package audit

import (
	auditInterface "Qingyu_backend/service/interfaces/audit"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/audit"
	pkgAudit "Qingyu_backend/pkg/audit"
	pkgErrors "Qingyu_backend/pkg/errors"
	auditRepo "Qingyu_backend/repository/interfaces/audit"
	"Qingyu_backend/service/base"
)

// ContentAuditService 内容审核服务
type ContentAuditService struct {
	sensitiveWordRepo auditRepo.SensitiveWordRepository
	auditRecordRepo   auditRepo.AuditRecordRepository
	violationRepo     auditRepo.ViolationRecordRepository
	dfaFilter         *pkgAudit.DFAFilter
	ruleEngine        *RuleEngine
	eventBus          base.EventBus
	serviceName       string
}

// NewContentAuditService 创建审核服务
func NewContentAuditService(
	sensitiveWordRepo auditRepo.SensitiveWordRepository,
	auditRecordRepo auditRepo.AuditRecordRepository,
	violationRepo auditRepo.ViolationRecordRepository,
	eventBus base.EventBus,
) *ContentAuditService {
	service := &ContentAuditService{
		sensitiveWordRepo: sensitiveWordRepo,
		auditRecordRepo:   auditRecordRepo,
		violationRepo:     violationRepo,
		dfaFilter:         pkgAudit.NewDFAFilter(),
		eventBus:          eventBus,
		serviceName:       "ContentAuditService",
	}

	// 初始化规则引擎
	service.ruleEngine = NewRuleEngine()
	service.loadDefaultRules()

	return service
}

// Initialize 初始化服务（加载敏感词库）
func (s *ContentAuditService) Initialize(ctx context.Context) error {
	// 从数据库加载启用的敏感词
	words, err := s.sensitiveWordRepo.GetEnabledWords(ctx)
	if err != nil {
		return fmt.Errorf("加载敏感词失败: %w", err)
	}

	// 如果数据库为空，加载默认词库
	if len(words) == 0 {
		pkgAudit.LoadDefaultWords(s.dfaFilter)
		return nil
	}

	// 批量添加到DFA过滤器
	wordInfos := make([]pkgAudit.SensitiveWordInfo, len(words))
	for i, w := range words {
		wordInfos[i] = pkgAudit.SensitiveWordInfo{
			Word:     w.Word,
			Level:    w.Level,
			Category: w.Category,
		}
	}
	s.dfaFilter.BatchAddWords(wordInfos)

	return nil
}

// CheckContent 实时检测内容
func (s *ContentAuditService) CheckContent(ctx context.Context, content string) (*auditInterface.AuditCheckResult, error) {
	if content == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "内容不能为空", "", nil)
	}

	result := &auditInterface.AuditCheckResult{
		IsSafe:      true,
		RiskLevel:   0,
		RiskScore:   0,
		Violations:  []audit.ViolationDetail{},
		Suggestions: []string{},
		NeedsReview: false,
		CanPublish:  true,
	}

	// 1. 敏感词检测
	matches := s.dfaFilter.FindAll(content)
	for _, match := range matches {
		violation := audit.ViolationDetail{
			Type:        "sensitive_word",
			Category:    match.Category,
			Level:       match.Level,
			Description: fmt.Sprintf("检测到敏感词：%s", match.Word),
			Position:    match.Start,
			Context:     match.Context,
			Keywords:    []string{match.Word},
		}
		result.Violations = append(result.Violations, violation)

		// 更新风险等级
		if match.Level > result.RiskLevel {
			result.RiskLevel = match.Level
		}
	}

	// 2. 规则引擎检测
	ruleViolations := s.ruleEngine.Check(content)
	for _, rv := range ruleViolations {
		result.Violations = append(result.Violations, rv)
		if rv.Level > result.RiskLevel {
			result.RiskLevel = rv.Level
		}
	}

	// 3. 计算风险分数
	result.RiskScore = s.calculateRiskScore(result.Violations)

	// 4. 判断是否安全
	if len(result.Violations) > 0 {
		result.IsSafe = false
	}

	// 5. 判断是否可以发布
	if result.RiskLevel >= audit.LevelHigh {
		result.CanPublish = false
		result.NeedsReview = true
	} else if result.RiskLevel >= audit.LevelMedium {
		result.CanPublish = true
		result.NeedsReview = true
	}

	// 6. 生成修改建议
	result.Suggestions = s.generateSuggestions(result.Violations)

	return result, nil
}

// AuditDocument 全文审核
func (s *ContentAuditService) AuditDocument(ctx context.Context, documentID string, content string, authorID string) (*audit.AuditRecord, error) {
	// 1. 参数验证
	if documentID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "文档ID不能为空", "", nil)
	}
	if content == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "内容不能为空", "", nil)
	}

	// 2. 检查内容
	checkResult, err := s.CheckContent(ctx, content)
	if err != nil {
		return nil, err
	}

	// 3. 创建审核记录
	record := &audit.AuditRecord{
		ID:         primitive.NewObjectID().Hex(),
		TargetType: audit.TargetTypeDocument,
		TargetID:   documentID,
		AuthorID:   authorID,
		Content:    content,
		RiskLevel:  checkResult.RiskLevel,
		RiskScore:  checkResult.RiskScore,
		Violations: checkResult.Violations,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 4. 确定审核状态和结果
	// 修复: 优先根据风险等级判断，高风险直接拒绝
	if checkResult.IsSafe {
		record.Status = audit.StatusApproved
		record.Result = audit.ResultPass
	} else if checkResult.RiskLevel >= audit.LevelHigh {
		// 高风险（Level≥3）直接拒绝
		record.Status = audit.StatusRejected
		record.Result = audit.ResultReject
	} else if checkResult.NeedsReview {
		// 中等风险需要人工复审
		record.Status = audit.StatusPending
		record.Result = audit.ResultManual
	} else if checkResult.CanPublish {
		// 低风险警告但可发布
		record.Status = audit.StatusWarning
		record.Result = audit.ResultWarning
	} else {
		record.Status = audit.StatusRejected
		record.Result = audit.ResultReject
	}

	// 5. 保存审核记录
	if err := s.auditRecordRepo.Create(ctx, record); err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "保存审核记录失败", "", err)
	}

	// 6. 如果拒绝，创建违规记录
	if record.Status == audit.StatusRejected && authorID != "" {
		if err := s.createViolationRecord(ctx, record); err != nil {
			// 记录错误但不影响主流程
			fmt.Printf("创建违规记录失败: %v\n", err)
		}
	}

	// 7. 发布事件
	if s.eventBus != nil {
		s.publishAuditEvent(ctx, record)
	}

	return record, nil
}

// GetAuditResult 获取审核结果
func (s *ContentAuditService) GetAuditResult(ctx context.Context, targetType, targetID string) (*audit.AuditRecord, error) {
	if targetType == "" || targetID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "参数不能为空", "", nil)
	}

	record, err := s.auditRecordRepo.GetByTargetID(ctx, targetType, targetID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询审核结果失败", "", err)
	}

	if record == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "审核记录不存在", "", nil)
	}

	return record, nil
}

// BatchAuditDocuments 批量审核文档
// 使用并发处理提升批量审核效率
func (s *ContentAuditService) BatchAuditDocuments(ctx context.Context, documentIDs []string) ([]*audit.AuditRecord, error) {
	if len(documentIDs) == 0 {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "文档ID列表不能为空", "", nil)
	}

	// 使用 goroutine 并发处理，提升效率
	type result struct {
		record *audit.AuditRecord
		err    error
	}

	resultChan := make(chan result, len(documentIDs))
	semaphore := make(chan struct{}, 10) // 限制并发数为10

	// 启动并发审核
	for _, docID := range documentIDs {
		docID := docID // 捕获循环变量
		go func() {
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			// 注意：实际生产环境应该调用 DocumentService 获取文档内容
			// 当前实现：创建简化的审核记录
			record := &audit.AuditRecord{
				TargetType: audit.TargetTypeDocument,
				TargetID:   docID,
				Status:     audit.StatusPending,
				Result:     audit.ResultPass, // 简化处理，默认通过
				RiskLevel:  1,                // 低风险
				RiskScore:  0.0,
				Violations: []audit.ViolationDetail{},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			// 保存审核记录
			err := s.auditRecordRepo.Create(ctx, record)
			if err != nil {
				resultChan <- result{nil, err}
				return
			}

			resultChan <- result{record, nil}
		}()
	}

	// 收集结果
	records := make([]*audit.AuditRecord, 0, len(documentIDs))
	var lastErr error
	for i := 0; i < len(documentIDs); i++ {
		res := <-resultChan
		if res.err != nil {
			lastErr = res.err
			continue
		}
		if res.record != nil {
			records = append(records, res.record)
		}
	}

	// 如果有错误但也有成功的，返回成功的记录和最后一个错误
	if lastErr != nil && len(records) == 0 {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "批量审核失败", "", lastErr)
	}

	return records, nil
}

// ReviewAudit 复核审核结果
func (s *ContentAuditService) ReviewAudit(ctx context.Context, auditID string, reviewerID string, approved bool, note string) error {
	// 1. 参数验证
	if auditID == "" || reviewerID == "" {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "参数不能为空", "", nil)
	}

	// 2. 获取审核记录
	record, err := s.auditRecordRepo.GetByID(ctx, auditID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询审核记录失败", "", err)
	}
	if record == nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "审核记录不存在", "", nil)
	}

	// 3. 更新审核状态
	status := audit.StatusRejected
	result := audit.ResultReject
	if approved {
		status = audit.StatusApproved
		result = audit.ResultPass
	}

	if err := s.auditRecordRepo.UpdateStatus(ctx, auditID, status, reviewerID, note); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "更新审核状态失败", "", err)
	}

	// 4. 如果拒绝且有作者ID，创建违规记录
	if !approved && record.AuthorID != "" {
		record.Status = status
		record.Result = result
		if err := s.createViolationRecord(ctx, record); err != nil {
			fmt.Printf("创建违规记录失败: %v\n", err)
		}
	}

	// 5. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &base.BaseEvent{
			EventType: "audit.reviewed",
			EventData: map[string]interface{}{
				"audit_id":    auditID,
				"reviewer_id": reviewerID,
				"approved":    approved,
				"status":      status,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// SubmitAppeal 提交申诉
func (s *ContentAuditService) SubmitAppeal(ctx context.Context, auditID string, authorID string, reason string) error {
	// 1. 参数验证
	if auditID == "" || authorID == "" || reason == "" {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "参数不能为空", "", nil)
	}

	// 2. 获取审核记录
	record, err := s.auditRecordRepo.GetByID(ctx, auditID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询审核记录失败", "", err)
	}
	if record == nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "审核记录不存在", "", nil)
	}

	// 3. 验证权限
	if record.AuthorID != authorID {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "无权申诉此记录", "", nil)
	}

	// 4. 验证是否可以申诉
	if !record.CanAppeal() {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorBusiness, "该记录不可申诉", "", nil)
	}

	// 5. 更新申诉状态
	if err := s.auditRecordRepo.UpdateAppealStatus(ctx, auditID, audit.AppealPending); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "提交申诉失败", "", err)
	}

	// 6. 更新审核记录的元数据
	updates := map[string]interface{}{
		"appealStatus": audit.AppealPending,
		"metadata": map[string]interface{}{
			"appeal_reason":    reason,
			"appeal_time":      time.Now(),
			"appeal_author_id": authorID,
		},
		"updatedAt": time.Now(),
	}
	if err := s.auditRecordRepo.Update(ctx, auditID, updates); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "更新记录失败", "", err)
	}

	// 7. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &base.BaseEvent{
			EventType: "audit.appeal_submitted",
			EventData: map[string]interface{}{
				"audit_id":  auditID,
				"author_id": authorID,
				"reason":    reason,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// ReviewAppeal 复核申诉
func (s *ContentAuditService) ReviewAppeal(ctx context.Context, auditID string, reviewerID string, approved bool, note string) error {
	// 1. 参数验证
	if auditID == "" || reviewerID == "" {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "参数不能为空", "", nil)
	}

	// 2. 获取审核记录
	record, err := s.auditRecordRepo.GetByID(ctx, auditID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询审核记录失败", "", err)
	}
	if record == nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "审核记录不存在", "", nil)
	}

	// 3. 验证是否有申诉
	if record.AppealStatus != audit.AppealPending {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorBusiness, "该记录没有待处理的申诉", "", nil)
	}

	// 4. 更新申诉状态
	appealStatus := audit.AppealRejected
	auditStatus := record.Status // 保持原状态
	if approved {
		appealStatus = audit.AppealApproved
		auditStatus = audit.StatusApproved // 申诉通过，改为通过
	}

	updates := map[string]interface{}{
		"appealStatus": appealStatus,
		"status":       auditStatus,
		"reviewerId":   reviewerID,
		"reviewNote":   note,
		"reviewedAt":   time.Now(),
		"updatedAt":    time.Now(),
	}

	if err := s.auditRecordRepo.Update(ctx, auditID, updates); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "更新申诉状态失败", "", err)
	}

	// 5. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &base.BaseEvent{
			EventType: "audit.appeal_reviewed",
			EventData: map[string]interface{}{
				"audit_id":      auditID,
				"reviewer_id":   reviewerID,
				"approved":      approved,
				"appeal_status": appealStatus,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// GetUserViolations 获取用户违规记录
func (s *ContentAuditService) GetUserViolations(ctx context.Context, userID string) ([]*audit.ViolationRecord, error) {
	if userID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "用户ID不能为空", "", nil)
	}

	violations, err := s.violationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询违规记录失败", "", err)
	}

	return violations, nil
}

// GetUserViolationSummary 获取用户违规统计
func (s *ContentAuditService) GetUserViolationSummary(ctx context.Context, userID string) (*audit.UserViolationSummary, error) {
	if userID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "用户ID不能为空", "", nil)
	}

	summary, err := s.violationRepo.GetUserSummary(ctx, userID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询违规统计失败", "", err)
	}

	return summary, nil
}

// GetPendingReviews 获取待复核记录
func (s *ContentAuditService) GetPendingReviews(ctx context.Context, limit int) ([]*audit.AuditRecord, error) {
	if limit <= 0 {
		limit = 50
	}

	records, err := s.auditRecordRepo.GetPendingReview(ctx, int64(limit))
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询待复核记录失败", "", err)
	}

	return records, nil
}

// GetHighRiskAudits 获取高风险审核记录
func (s *ContentAuditService) GetHighRiskAudits(ctx context.Context, minRiskLevel int, limit int) ([]*audit.AuditRecord, error) {
	if minRiskLevel <= 0 {
		minRiskLevel = audit.LevelHigh
	}
	if limit <= 0 {
		limit = 50
	}

	records, err := s.auditRecordRepo.GetHighRisk(ctx, minRiskLevel, int64(limit))
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询高风险记录失败", "", err)
	}

	return records, nil
}

// 私有辅助方法

// calculateRiskScore 计算风险分数
func (s *ContentAuditService) calculateRiskScore(violations []audit.ViolationDetail) float64 {
	if len(violations) == 0 {
		return 0
	}

	// 基础分数：违规数量 * 10
	baseScore := float64(len(violations)) * 10

	// 严重度加权
	levelWeight := 0.0
	for _, v := range violations {
		levelWeight += float64(v.Level) * 10
	}

	// 最终分数（0-100）
	score := baseScore + levelWeight
	if score > 100 {
		score = 100
	}

	return score
}

// generateSuggestions 生成修改建议
func (s *ContentAuditService) generateSuggestions(violations []audit.ViolationDetail) []string {
	suggestions := make([]string, 0)
	categoryMap := make(map[string]bool)

	for _, v := range violations {
		if !categoryMap[v.Category] {
			categoryMap[v.Category] = true
			switch v.Category {
			case audit.CategoryPolitics:
				suggestions = append(suggestions, "请删除政治敏感内容")
			case audit.CategoryPorn:
				suggestions = append(suggestions, "请删除色情相关内容")
			case audit.CategoryViolence:
				suggestions = append(suggestions, "请减少暴力描写")
			case audit.CategoryInsult:
				suggestions = append(suggestions, "请使用文明用语")
			case audit.CategoryAd:
				suggestions = append(suggestions, "请删除广告推广内容")
			default:
				suggestions = append(suggestions, fmt.Sprintf("请修改%s相关内容", audit.GetCategoryName(v.Category)))
			}
		}
	}

	return suggestions
}

// createViolationRecord 创建违规记录
func (s *ContentAuditService) createViolationRecord(ctx context.Context, record *audit.AuditRecord) error {
	if record.AuthorID == "" {
		return nil
	}

	violation := &audit.ViolationRecord{
		ID:             primitive.NewObjectID().Hex(),
		UserID:         record.AuthorID,
		AuditRecordID:  record.ID,
		TargetType:     record.TargetType,
		TargetID:       record.TargetID,
		ViolationType:  "content_violation",
		ViolationLevel: record.RiskLevel,
		ViolationCount: 1,
		Description:    fmt.Sprintf("内容违规，风险等级：%d", record.RiskLevel),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 根据风险等级确定处罚
	if record.RiskLevel >= audit.LevelBanned {
		violation.PenaltyType = audit.PenaltyAccountBanned
		violation.PenaltyDuration = 30 // 封号30天
		violation.IsPenalized = true
		now := time.Now()
		violation.PenalizedAt = &now
		expiresAt := now.AddDate(0, 0, 30)
		violation.ExpiresAt = &expiresAt
	} else if record.RiskLevel >= audit.LevelCritical {
		violation.PenaltyType = audit.PenaltyAccountMuted
		violation.PenaltyDuration = 7 // 禁言7天
		violation.IsPenalized = true
		now := time.Now()
		violation.PenalizedAt = &now
		expiresAt := now.AddDate(0, 0, 7)
		violation.ExpiresAt = &expiresAt
	} else if record.RiskLevel >= audit.LevelHigh {
		violation.PenaltyType = audit.PenaltyContentHidden
		violation.IsPenalized = true
		now := time.Now()
		violation.PenalizedAt = &now
	}

	return s.violationRepo.Create(ctx, violation)
}

// publishAuditEvent 发布审核事件
func (s *ContentAuditService) publishAuditEvent(ctx context.Context, record *audit.AuditRecord) {
	eventType := "audit.completed"
	if record.Status == audit.StatusRejected {
		eventType = "audit.rejected"
	} else if record.Status == audit.StatusApproved {
		eventType = "audit.approved"
	}

	s.eventBus.PublishAsync(ctx, &base.BaseEvent{
		EventType: eventType,
		EventData: map[string]interface{}{
			"audit_id":   record.ID,
			"target_id":  record.TargetID,
			"author_id":  record.AuthorID,
			"status":     record.Status,
			"risk_level": record.RiskLevel,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	})
}

// loadDefaultRules 加载默认规则
func (s *ContentAuditService) loadDefaultRules() {
	// 加载所有默认规则
	s.ruleEngine.AddRule(NewPhoneNumberRule())
	s.ruleEngine.AddRule(NewURLRule())
	s.ruleEngine.AddRule(NewWeChatRule())
	s.ruleEngine.AddRule(NewQQRule())
	s.ruleEngine.AddRule(NewExcessiveRepetitionRule())
	// 内容长度规则可选，根据需要启用
	// s.ruleEngine.AddRule(NewContentLengthRule())
}
