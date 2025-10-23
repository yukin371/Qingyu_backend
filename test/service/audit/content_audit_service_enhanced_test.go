package audit_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	auditModels "Qingyu_backend/models/audit"
	pkgAudit "Qingyu_backend/pkg/audit"
	"Qingyu_backend/repository/interfaces/infrastructure"
	auditService "Qingyu_backend/service/audit"
	"Qingyu_backend/service/base"
)

// ============================================
// Phase 4.1: 内容审核Service测试
// ============================================
//
// 对应SRS需求：REQ-AI-AGENT-003（内容安全，P0优先级）
//
// 测试分类：
// - Phase 1: 敏感词检测（4个测试）
// - Phase 2: 合规性检查（4个测试）
// - Phase 3: 性能与优化（4个测试）
// ============================================

// ==============================================
// Mock实现
// ==============================================

// MockSensitiveWordRepository Mock敏感词Repository
type MockSensitiveWordRepository struct {
	mock.Mock
}

func (m *MockSensitiveWordRepository) GetEnabledWords(ctx context.Context) ([]*auditModels.SensitiveWord, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.SensitiveWord), args.Error(1)
}

func (m *MockSensitiveWordRepository) Create(ctx context.Context, word *auditModels.SensitiveWord) error {
	args := m.Called(ctx, word)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) GetByID(ctx context.Context, id string) (*auditModels.SensitiveWord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditModels.SensitiveWord), args.Error(1)
}

func (m *MockSensitiveWordRepository) GetByWord(ctx context.Context, word string) (*auditModels.SensitiveWord, error) {
	args := m.Called(ctx, word)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditModels.SensitiveWord), args.Error(1)
}

func (m *MockSensitiveWordRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*auditModels.SensitiveWord, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.SensitiveWord), args.Error(1)
}

func (m *MockSensitiveWordRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSensitiveWordRepository) FindWithPagination(ctx context.Context, filter infrastructure.Filter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[auditModels.SensitiveWord], error) {
	args := m.Called(ctx, filter, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*infrastructure.PagedResult[auditModels.SensitiveWord]), args.Error(1)
}

func (m *MockSensitiveWordRepository) GetByCategory(ctx context.Context, category string) ([]*auditModels.SensitiveWord, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.SensitiveWord), args.Error(1)
}

func (m *MockSensitiveWordRepository) GetByLevel(ctx context.Context, minLevel int) ([]*auditModels.SensitiveWord, error) {
	args := m.Called(ctx, minLevel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.SensitiveWord), args.Error(1)
}

func (m *MockSensitiveWordRepository) BatchCreate(ctx context.Context, words []*auditModels.SensitiveWord) error {
	args := m.Called(ctx, words)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) BatchUpdate(ctx context.Context, ids []string, updates map[string]interface{}) error {
	args := m.Called(ctx, ids, updates)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) BatchDelete(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) CountByCategory(ctx context.Context) (map[string]int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockSensitiveWordRepository) CountByLevel(ctx context.Context) (map[int]int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[int]int64), args.Error(1)
}

// MockAuditRecordRepository Mock审核记录Repository
type MockAuditRecordRepository struct {
	mock.Mock
}

func (m *MockAuditRecordRepository) Create(ctx context.Context, record *auditModels.AuditRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockAuditRecordRepository) GetByID(ctx context.Context, id string) (*auditModels.AuditRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditModels.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) GetByTargetID(ctx context.Context, targetType, targetID string) (*auditModels.AuditRecord, error) {
	args := m.Called(ctx, targetType, targetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditModels.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) UpdateStatus(ctx context.Context, id, status, reviewerID, note string) error {
	args := m.Called(ctx, id, status, reviewerID, note)
	return args.Error(0)
}

func (m *MockAuditRecordRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockAuditRecordRepository) GetPendingReview(ctx context.Context, limit int64) ([]*auditModels.AuditRecord, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) GetHighRisk(ctx context.Context, minRiskLevel int, limit int64) ([]*auditModels.AuditRecord, error) {
	args := m.Called(ctx, minRiskLevel, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) UpdateAppealStatus(ctx context.Context, id, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockAuditRecordRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAuditRecordRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*auditModels.AuditRecord, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuditRecordRepository) FindWithPagination(ctx context.Context, filter infrastructure.Filter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[auditModels.AuditRecord], error) {
	args := m.Called(ctx, filter, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*infrastructure.PagedResult[auditModels.AuditRecord]), args.Error(1)
}

func (m *MockAuditRecordRepository) GetByAuthor(ctx context.Context, authorID string, limit int64) ([]*auditModels.AuditRecord, error) {
	args := m.Called(ctx, authorID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) GetByStatus(ctx context.Context, status string, limit int64) ([]*auditModels.AuditRecord, error) {
	args := m.Called(ctx, status, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) GetByResult(ctx context.Context, result string, limit int64) ([]*auditModels.AuditRecord, error) {
	args := m.Called(ctx, result, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) BatchUpdateStatus(ctx context.Context, ids []string, status string) error {
	args := m.Called(ctx, ids, status)
	return args.Error(0)
}

func (m *MockAuditRecordRepository) CountByStatus(ctx context.Context, startTime, endTime time.Time) (map[string]int64, error) {
	args := m.Called(ctx, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockAuditRecordRepository) CountByAuthor(ctx context.Context, authorID string) (int64, error) {
	args := m.Called(ctx, authorID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuditRecordRepository) CountHighRiskByAuthor(ctx context.Context, authorID string, minLevel int) (int64, error) {
	args := m.Called(ctx, authorID, minLevel)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuditRecordRepository) GetRecentRecords(ctx context.Context, targetType string, limit int64) ([]*auditModels.AuditRecord, error) {
	args := m.Called(ctx, targetType, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockViolationRecordRepository Mock违规记录Repository
type MockViolationRecordRepository struct {
	mock.Mock
}

func (m *MockViolationRecordRepository) Create(ctx context.Context, record *auditModels.ViolationRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockViolationRecordRepository) GetByUserID(ctx context.Context, userID string) ([]*auditModels.ViolationRecord, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.ViolationRecord), args.Error(1)
}

func (m *MockViolationRecordRepository) GetUserSummary(ctx context.Context, userID string) (*auditModels.UserViolationSummary, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditModels.UserViolationSummary), args.Error(1)
}

func (m *MockViolationRecordRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockViolationRecordRepository) GetByID(ctx context.Context, id string) (*auditModels.ViolationRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditModels.ViolationRecord), args.Error(1)
}

func (m *MockViolationRecordRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockViolationRecordRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockViolationRecordRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*auditModels.ViolationRecord, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.ViolationRecord), args.Error(1)
}

func (m *MockViolationRecordRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockViolationRecordRepository) FindWithPagination(ctx context.Context, filter infrastructure.Filter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[auditModels.ViolationRecord], error) {
	args := m.Called(ctx, filter, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*infrastructure.PagedResult[auditModels.ViolationRecord]), args.Error(1)
}

func (m *MockViolationRecordRepository) GetActiveViolations(ctx context.Context, userID string) ([]*auditModels.ViolationRecord, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.ViolationRecord), args.Error(1)
}

func (m *MockViolationRecordRepository) GetRecentViolations(ctx context.Context, userID string, since time.Time) ([]*auditModels.ViolationRecord, error) {
	args := m.Called(ctx, userID, since)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.ViolationRecord), args.Error(1)
}

func (m *MockViolationRecordRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockViolationRecordRepository) CountHighRiskByUser(ctx context.Context, userID string, minLevel int) (int64, error) {
	args := m.Called(ctx, userID, minLevel)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockViolationRecordRepository) GetTopViolators(ctx context.Context, limit int) ([]*auditModels.UserViolationSummary, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.UserViolationSummary), args.Error(1)
}

func (m *MockViolationRecordRepository) ApplyPenalty(ctx context.Context, userID string, penaltyType string, duration int, description string) error {
	args := m.Called(ctx, userID, penaltyType, duration, description)
	return args.Error(0)
}

func (m *MockViolationRecordRepository) RemovePenalty(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockViolationRecordRepository) GetActivePenalties(ctx context.Context) ([]*auditModels.ViolationRecord, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModels.ViolationRecord), args.Error(1)
}

func (m *MockViolationRecordRepository) CleanExpiredPenalties(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockEventBus Mock事件总线
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	args := m.Called(eventType, handlerName)
	return args.Error(0)
}

func (m *MockEventBus) Publish(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// ==============================================
// Phase 1: 敏感词检测（4个测试用例）
// ==============================================

// TestContentAudit_SensitiveWordMatch_Politics 测试政治敏感词匹配
func TestContentAudit_SensitiveWordMatch_Politics(t *testing.T) {
	// 创建Mock
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockAuditRepo := new(MockAuditRecordRepository)
	mockViolationRepo := new(MockViolationRecordRepository)
	mockEventBus := new(MockEventBus)

	// 创建Service
	auditSvc := auditService.NewContentAuditService(
		mockSensitiveRepo,
		mockAuditRepo,
		mockViolationRepo,
		mockEventBus,
	)

	ctx := context.Background()

	// Mock加载敏感词库
	sensitiveWords := []*auditModels.SensitiveWord{
		{
			ID:       "word_politics_1",
			Word:     "反动",
			Category: auditModels.CategoryPolitics,
			Level:    auditModels.LevelCritical,
		},
		{
			ID:       "word_politics_2",
			Word:     "暴乱",
			Category: auditModels.CategoryPolitics,
			Level:    auditModels.LevelHigh,
		},
	}
	mockSensitiveRepo.On("GetEnabledWords", ctx).Return(sensitiveWords, nil)

	// 初始化服务（加载敏感词）
	err := auditSvc.Initialize(ctx)
	assert.NoError(t, err)

	// 测试内容检测
	content := "这是一段包含反动言论和暴乱煽动的内容"
	result, err := auditSvc.CheckContent(ctx, content)
	assert.NoError(t, err)

	// 验证检测结果
	assert.False(t, result.IsSafe, "包含敏感词的内容不应该被判定为安全")
	assert.GreaterOrEqual(t, len(result.Violations), 2, "应该检测到至少2个敏感词")
	assert.Equal(t, auditModels.LevelCritical, result.RiskLevel, "最高风险等级应该是Critical")
	assert.False(t, result.CanPublish, "严重违规内容不应该允许发布")
	assert.True(t, result.NeedsReview, "高风险内容需要人工复核")

	t.Logf("检测到 %d 个违规项，风险等级：%d，风险分数：%.2f",
		len(result.Violations), result.RiskLevel, result.RiskScore)
}

// TestContentAudit_SensitiveWordMatch_PornAndViolence 测试色情和暴力敏感词
func TestContentAudit_SensitiveWordMatch_PornAndViolence(t *testing.T) {
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockAuditRepo := new(MockAuditRecordRepository)
	mockViolationRepo := new(MockViolationRecordRepository)
	mockEventBus := new(MockEventBus)

	auditSvc := auditService.NewContentAuditService(
		mockSensitiveRepo,
		mockAuditRepo,
		mockViolationRepo,
		mockEventBus,
	)

	ctx := context.Background()

	// Mock不同类别的敏感词
	sensitiveWords := []*auditModels.SensitiveWord{
		{
			Word:     "黄色内容",
			Category: auditModels.CategoryPorn,
			Level:    auditModels.LevelHigh,
		},
		{
			Word:     "血腥暴力",
			Category: auditModels.CategoryViolence,
			Level:    auditModels.LevelMedium,
		},
	}
	mockSensitiveRepo.On("GetEnabledWords", ctx).Return(sensitiveWords, nil)

	err := auditSvc.Initialize(ctx)
	assert.NoError(t, err)

	// 测试色情内容
	content1 := "这里有黄色内容链接"
	result1, err := auditSvc.CheckContent(ctx, content1)
	assert.NoError(t, err)
	assert.False(t, result1.IsSafe)
	assert.Equal(t, auditModels.LevelHigh, result1.RiskLevel)

	// 测试暴力内容
	content2 := "故事中描述了血腥暴力场面"
	result2, err := auditSvc.CheckContent(ctx, content2)
	assert.NoError(t, err)
	assert.False(t, result2.IsSafe)
	assert.Equal(t, auditModels.LevelMedium, result2.RiskLevel)

	t.Logf("色情内容风险等级：%d，暴力内容风险等级：%d", result1.RiskLevel, result2.RiskLevel)
}

// TestContentAudit_ReplacementSuggestions 测试敏感词替换建议
func TestContentAudit_ReplacementSuggestions(t *testing.T) {
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockAuditRepo := new(MockAuditRecordRepository)
	mockViolationRepo := new(MockViolationRecordRepository)
	mockEventBus := new(MockEventBus)

	auditSvc := auditService.NewContentAuditService(
		mockSensitiveRepo,
		mockAuditRepo,
		mockViolationRepo,
		mockEventBus,
	)

	ctx := context.Background()

	// Mock敏感词
	sensitiveWords := []*auditModels.SensitiveWord{
		{
			Word:        "脏话",
			Category:    auditModels.CategoryInsult,
			Level:       auditModels.LevelLow,
			Replacement: "***",
		},
	}
	mockSensitiveRepo.On("GetEnabledWords", ctx).Return(sensitiveWords, nil)

	err := auditSvc.Initialize(ctx)
	assert.NoError(t, err)

	// 测试检测和建议
	content := "你这个脏话满口的人"
	result, err := auditSvc.CheckContent(ctx, content)
	assert.NoError(t, err)

	// 验证建议
	assert.NotEmpty(t, result.Suggestions, "应该提供修改建议")
	assert.Contains(t, result.Suggestions[0], "文明用语", "建议应该包含文明用语提示")

	t.Logf("修改建议：%v", result.Suggestions)
}

// TestContentAudit_SensitiveWordLibraryUpdate 测试敏感词库更新
func TestContentAudit_SensitiveWordLibraryUpdate(t *testing.T) {
	// 测试DFAFilter的动态更新能力
	filter := pkgAudit.NewDFAFilter()

	// 初始添加敏感词
	filter.AddWord("测试词1", auditModels.LevelHigh, auditModels.CategoryPolitics)
	assert.True(t, filter.Check("包含测试词1的内容"), "应该检测到敏感词")

	// 批量添加
	words := []pkgAudit.SensitiveWordInfo{
		{Word: "测试词2", Level: auditModels.LevelMedium, Category: auditModels.CategoryPorn},
		{Word: "测试词3", Level: auditModels.LevelLow, Category: auditModels.CategoryViolence},
	}
	filter.BatchAddWords(words)

	content := "测试词1 测试词2 测试词3"
	matches := filter.FindAll(content)
	assert.Equal(t, 3, len(matches), "应该找到3个敏感词")

	// 移除敏感词
	filter.RemoveWord("测试词2")
	matches = filter.FindAll(content)
	assert.Equal(t, 2, len(matches), "移除后应该只找到2个敏感词")

	// 获取统计信息
	stats := filter.GetStatistics()
	assert.Equal(t, 2, stats.TotalWords, "总词数应该是2")
	assert.Equal(t, 1, stats.ByCategory[auditModels.CategoryPolitics])
	assert.Equal(t, 1, stats.ByCategory[auditModels.CategoryViolence])

	t.Logf("敏感词库统计：总词数=%d，分类统计=%v", stats.TotalWords, stats.ByCategory)
}

// ==============================================
// Phase 2: 合规性检查（4个测试用例）
// ==============================================

// TestContentAudit_PhoneNumberAndURLDetection 测试手机号和URL检测
func TestContentAudit_PhoneNumberAndURLDetection(t *testing.T) {
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockAuditRepo := new(MockAuditRecordRepository)
	mockViolationRepo := new(MockViolationRecordRepository)
	mockEventBus := new(MockEventBus)

	auditSvc := auditService.NewContentAuditService(
		mockSensitiveRepo,
		mockAuditRepo,
		mockViolationRepo,
		mockEventBus,
	)

	ctx := context.Background()

	// Mock空敏感词库
	mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*auditModels.SensitiveWord{}, nil)
	err := auditSvc.Initialize(ctx)
	assert.NoError(t, err)

	// 测试手机号检测
	content1 := "联系我：13812345678，欢迎咨询"
	result1, err := auditSvc.CheckContent(ctx, content1)
	assert.NoError(t, err)
	assert.False(t, result1.IsSafe, "包含手机号应该被检测出")

	foundPhone := false
	for _, v := range result1.Violations {
		if v.Type == "phone_number" {
			foundPhone = true
			assert.Equal(t, auditModels.CategoryAd, v.Category)
			assert.Contains(t, v.Keywords[0], "138")
			break
		}
	}
	assert.True(t, foundPhone, "应该检测到手机号违规")

	// 测试URL检测
	content2 := "详情请访问 https://example.com/promo"
	result2, err := auditSvc.CheckContent(ctx, content2)
	assert.NoError(t, err)
	assert.False(t, result2.IsSafe, "包含URL应该被检测出")

	foundURL := false
	for _, v := range result2.Violations {
		if v.Type == "url" {
			foundURL = true
			assert.Equal(t, auditModels.CategoryAd, v.Category)
			break
		}
	}
	assert.True(t, foundURL, "应该检测到URL违规")

	t.Logf("手机号检测：%d个违规，URL检测：%d个违规",
		len(result1.Violations), len(result2.Violations))
}

// TestContentAudit_WeChatAndQQDetection 测试微信和QQ号检测
func TestContentAudit_WeChatAndQQDetection(t *testing.T) {
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockAuditRepo := new(MockAuditRecordRepository)
	mockViolationRepo := new(MockViolationRecordRepository)
	mockEventBus := new(MockEventBus)

	auditSvc := auditService.NewContentAuditService(
		mockSensitiveRepo,
		mockAuditRepo,
		mockViolationRepo,
		mockEventBus,
	)

	ctx := context.Background()

	mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*auditModels.SensitiveWord{}, nil)
	err := auditSvc.Initialize(ctx)
	assert.NoError(t, err)

	// 测试微信检测
	content1 := "加我微信：wx123456"
	result1, err := auditSvc.CheckContent(ctx, content1)
	assert.NoError(t, err)
	assert.False(t, result1.IsSafe)

	foundWeChat := false
	for _, v := range result1.Violations {
		if v.Type == "wechat" {
			foundWeChat = true
			break
		}
	}
	assert.True(t, foundWeChat, "应该检测到微信号违规")

	// 测试QQ检测
	content2 := "加我QQ群：12345678"
	result2, err := auditSvc.CheckContent(ctx, content2)
	assert.NoError(t, err)
	assert.False(t, result2.IsSafe)

	foundQQ := false
	for _, v := range result2.Violations {
		if v.Type == "qq" {
			foundQQ = true
			break
		}
	}
	assert.True(t, foundQQ, "应该检测到QQ号违规")

	t.Logf("微信检测成功：%v，QQ检测成功：%v", foundWeChat, foundQQ)
}

// TestContentAudit_ViolationRecordCreation 测试违规记录创建
func TestContentAudit_ViolationRecordCreation(t *testing.T) {
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockAuditRepo := new(MockAuditRecordRepository)
	mockViolationRepo := new(MockViolationRecordRepository)
	mockEventBus := new(MockEventBus)

	auditSvc := auditService.NewContentAuditService(
		mockSensitiveRepo,
		mockAuditRepo,
		mockViolationRepo,
		mockEventBus,
	)

	ctx := context.Background()

	// Mock敏感词库（严重违规）
	sensitiveWords := []*auditModels.SensitiveWord{
		{
			Word:     "严重违规词",
			Category: auditModels.CategoryPolitics,
			Level:    auditModels.LevelBanned, // 最高等级，会触发封号
		},
	}
	mockSensitiveRepo.On("GetEnabledWords", ctx).Return(sensitiveWords, nil)
	err := auditSvc.Initialize(ctx)
	assert.NoError(t, err)

	// Mock审核记录保存
	mockAuditRepo.On("Create", ctx, mock.AnythingOfType("*audit.AuditRecord")).Return(nil)

	// Mock违规记录保存
	mockViolationRepo.On("Create", ctx, mock.AnythingOfType("*audit.ViolationRecord")).
		Run(func(args mock.Arguments) {
			violation := args.Get(1).(*auditModels.ViolationRecord)
			// 验证违规记录内容
			assert.Equal(t, "author_test", violation.UserID)
			assert.Equal(t, auditModels.LevelBanned, violation.ViolationLevel)
			assert.Equal(t, auditModels.PenaltyAccountBanned, violation.PenaltyType)
			assert.True(t, violation.IsPenalized)
			assert.Equal(t, 30, violation.PenaltyDuration) // 封号30天
		}).
		Return(nil)

	// Mock事件发布
	mockEventBus.On("PublishAsync", ctx, mock.AnythingOfType("*base.BaseEvent")).Return(nil)

	// 执行审核
	content := "这里包含严重违规词的内容"
	authorID := "author_test"
	documentID := "doc_test"

	record, err := auditSvc.AuditDocument(ctx, documentID, content, authorID)
	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, auditModels.StatusRejected, record.Status)
	assert.Equal(t, auditModels.ResultReject, record.Result)
	assert.Equal(t, auditModels.LevelBanned, record.RiskLevel)

	// 验证Mock调用
	mockAuditRepo.AssertExpectations(t)
	mockViolationRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)

	t.Logf("违规记录创建成功，状态：%s，结果：%s", record.Status, record.Result)
}

// TestContentAudit_ManualReviewTrigger 测试人工复审触发
func TestContentAudit_ManualReviewTrigger(t *testing.T) {
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockAuditRepo := new(MockAuditRecordRepository)
	mockViolationRepo := new(MockViolationRecordRepository)
	mockEventBus := new(MockEventBus)

	auditSvc := auditService.NewContentAuditService(
		mockSensitiveRepo,
		mockAuditRepo,
		mockViolationRepo,
		mockEventBus,
	)

	ctx := context.Background()

	// Mock中等风险敏感词（触发人工复核）
	sensitiveWords := []*auditModels.SensitiveWord{
		{
			Word:     "可疑内容",
			Category: auditModels.CategoryOther,
			Level:    auditModels.LevelMedium, // 中等风险
		},
	}
	mockSensitiveRepo.On("GetEnabledWords", ctx).Return(sensitiveWords, nil)
	err := auditSvc.Initialize(ctx)
	assert.NoError(t, err)

	// 检测内容
	content := "这是一段包含可疑内容的文本"
	result, err := auditSvc.CheckContent(ctx, content)
	assert.NoError(t, err)

	// 验证触发人工复核
	assert.Equal(t, auditModels.LevelMedium, result.RiskLevel)
	assert.True(t, result.NeedsReview, "中等风险内容应该触发人工复核")
	assert.True(t, result.CanPublish, "中等风险内容可以发布但需要复核")

	t.Logf("风险等级：%d，需要复核：%v，可以发布：%v",
		result.RiskLevel, result.NeedsReview, result.CanPublish)
}

// ==============================================
// Phase 3: 性能与优化（4个测试用例）
// ==============================================

// TestContentAudit_LargeDocumentPerformance 测试大文档审核性能（<3秒/万字）
func TestContentAudit_LargeDocumentPerformance(t *testing.T) {
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockAuditRepo := new(MockAuditRecordRepository)
	mockViolationRepo := new(MockViolationRecordRepository)
	mockEventBus := new(MockEventBus)

	auditSvc := auditService.NewContentAuditService(
		mockSensitiveRepo,
		mockAuditRepo,
		mockViolationRepo,
		mockEventBus,
	)

	ctx := context.Background()

	// Mock敏感词库（实际场景）
	sensitiveWords := make([]*auditModels.SensitiveWord, 0)
	for i := 0; i < 100; i++ { // 100个敏感词
		sensitiveWords = append(sensitiveWords, &auditModels.SensitiveWord{
			Word:     fmt.Sprintf("敏感词%d", i),
			Category: auditModels.CategoryOther,
			Level:    auditModels.LevelLow,
		})
	}
	mockSensitiveRepo.On("GetEnabledWords", ctx).Return(sensitiveWords, nil)
	err := auditSvc.Initialize(ctx)
	assert.NoError(t, err)

	// 生成1万字的测试文档
	var content strings.Builder
	baseText := "这是一段正常的文本内容，用于性能测试。"
	for i := 0; i < 500; i++ { // 约1万字
		content.WriteString(baseText)
	}
	testContent := content.String()
	wordCount := len([]rune(testContent))

	// 性能测试
	startTime := time.Now()
	result, err := auditSvc.CheckContent(ctx, testContent)
	duration := time.Since(startTime)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// 验证性能（<3秒/万字）
	expectedMaxDuration := time.Duration(wordCount/10000*3) * time.Second
	if expectedMaxDuration < 3*time.Second {
		expectedMaxDuration = 3 * time.Second
	}

	assert.Less(t, duration, expectedMaxDuration,
		fmt.Sprintf("审核%d字应在%v内完成", wordCount, expectedMaxDuration))

	t.Logf("性能测试：%d字文档，耗时%v，平均速度：%.0f字/秒",
		wordCount, duration, float64(wordCount)/duration.Seconds())
}

// TestContentAudit_BatchAudit 测试批量审核
func TestContentAudit_BatchAudit(t *testing.T) {
	t.Skip("批量审核功能待完善，当前实现为TODO")

	// TODO: 当BatchAuditDocuments完整实现后，补充此测试
	// 预期测试内容：
	// 1. 批量提交多个文档ID
	// 2. 并发审核处理
	// 3. 返回所有审核结果
	// 4. 验证审核记录正确保存
}

// TestContentAudit_ResultCache 测试审核结果缓存
func TestContentAudit_ResultCache(t *testing.T) {
	t.Skip("TDD: 审核结果缓存功能未实现，待开发")

	// TODO: 实现审核结果缓存机制
	// 实现要点：
	// 1. 相同内容的审核结果缓存（使用内容hash作为key）
	// 2. 缓存有效期（如1小时）
	// 3. 敏感词库更新时清除缓存
	// 4. 缓存命中率统计

	// 测试流程：
	// 1. 首次审核内容A → 无缓存，执行完整审核
	// 2. 再次审核相同内容A → 命中缓存，直接返回结果
	// 3. 审核不同内容B → 无缓存
	// 4. 更新敏感词库 → 清除缓存
	// 5. 再次审核内容A → 无缓存（已清除）
}

// TestContentAudit_AsyncAuditQueue 测试异步审核队列
func TestContentAudit_AsyncAuditQueue(t *testing.T) {
	t.Skip("TDD: 异步审核队列功能未实现，待开发")

	// TODO: 实现异步审核队列
	// 实现要点：
	// 1. 审核任务加入队列（Redis/RabbitMQ）
	// 2. 后台Worker异步处理
	// 3. 审核完成后通知回调
	// 4. 任务状态查询（pending/processing/completed/failed）
	// 5. 失败重试机制（最多3次）

	// 使用场景：
	// - 大文档审核（>10万字）
	// - 批量审核
	// - 定时全量审核

	// 测试流程：
	// 1. 提交异步审核任务 → 返回任务ID
	// 2. 查询任务状态 → "pending"
	// 3. Worker处理 → 状态变为"processing"
	// 4. 审核完成 → 状态变为"completed"
	// 5. 获取审核结果
}

// ==============================================
// 辅助测试：规则引擎和DFA过滤器
// ==============================================

// TestDFAFilter_BasicFunctionality 测试DFA过滤器基础功能
func TestDFAFilter_BasicFunctionality(t *testing.T) {
	filter := pkgAudit.NewDFAFilter()

	// 添加敏感词
	filter.AddWord("测试敏感词", auditModels.LevelHigh, auditModels.CategoryPolitics)

	// 检测功能
	assert.True(t, filter.Check("这里有测试敏感词"), "应该检测到敏感词")
	assert.False(t, filter.Check("这里没有问题"), "不应该检测到敏感词")

	// 查找所有匹配
	content := "测试敏感词在前，中间还有测试敏感词，后面也有测试敏感词"
	matches := filter.FindAll(content)
	assert.Equal(t, 3, len(matches), "应该找到3个匹配")

	// 验证匹配结果
	for i, match := range matches {
		assert.Equal(t, "测试敏感词", match.Word)
		assert.Equal(t, auditModels.LevelHigh, match.Level)
		assert.Equal(t, auditModels.CategoryPolitics, match.Category)
		t.Logf("匹配%d: 位置[%d:%d], 上下文: %s", i+1, match.Start, match.End, match.Context)
	}

	// 替换功能
	replaced := filter.Replace(content, "***")
	assert.NotContains(t, replaced, "测试敏感词")
	assert.Contains(t, replaced, "***")

	// 掩码替换
	masked := filter.ReplaceWithMask(content, '*')
	assert.Equal(t, 3, strings.Count(masked, "*****")) // "测试敏感词"被替换为5个*

	t.Logf("原文：%s", content)
	t.Logf("替换：%s", replaced)
	t.Logf("掩码：%s", masked)
}

// TestDFAFilter_CaseInsensitive 测试大小写不敏感
func TestDFAFilter_CaseInsensitive(t *testing.T) {
	filter := pkgAudit.NewDFAFilter()
	filter.AddWord("Test", auditModels.LevelMedium, auditModels.CategoryOther)

	// DFA过滤器应该对大小写不敏感
	assert.True(t, filter.Check("test"))
	assert.True(t, filter.Check("TEST"))
	assert.True(t, filter.Check("Test"))
	assert.True(t, filter.Check("TeSt"))

	t.Log("DFA过滤器大小写不敏感测试通过")
}

// TestRuleEngine_ComplexRules 测试规则引擎复杂规则
func TestRuleEngine_ComplexRules(t *testing.T) {
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockAuditRepo := new(MockAuditRecordRepository)
	mockViolationRepo := new(MockViolationRecordRepository)
	mockEventBus := new(MockEventBus)

	auditSvc := auditService.NewContentAuditService(
		mockSensitiveRepo,
		mockAuditRepo,
		mockViolationRepo,
		mockEventBus,
	)

	ctx := context.Background()

	mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*auditModels.SensitiveWord{}, nil)
	err := auditSvc.Initialize(ctx)
	assert.NoError(t, err)

	// 测试包含多种违规的复杂内容
	content := "联系电话13812345678，详情访问https://spam.com，加我微信wx666，QQ群123456，aaaaaaaaaaaa（过度重复）"
	result, err := auditSvc.CheckContent(ctx, content)
	assert.NoError(t, err)

	// 应该检测到多种违规类型
	violationTypes := make(map[string]bool)
	for _, v := range result.Violations {
		violationTypes[v.Type] = true
	}

	assert.True(t, violationTypes["phone_number"], "应该检测到手机号")
	assert.True(t, violationTypes["url"], "应该检测到URL")
	assert.True(t, violationTypes["wechat"], "应该检测到微信")
	assert.True(t, violationTypes["qq"], "应该检测到QQ")
	assert.True(t, violationTypes["excessive_repetition"], "应该检测到过度重复")

	t.Logf("复杂内容检测到%d种违规类型：%v", len(violationTypes), violationTypes)
}

// ==============================================
// 边界测试
// ==============================================

// TestContentAudit_EmptyContent 测试空内容处理
func TestContentAudit_EmptyContent(t *testing.T) {
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockAuditRepo := new(MockAuditRecordRepository)
	mockViolationRepo := new(MockViolationRecordRepository)
	mockEventBus := new(MockEventBus)

	auditSvc := auditService.NewContentAuditService(
		mockSensitiveRepo,
		mockAuditRepo,
		mockViolationRepo,
		mockEventBus,
	)

	ctx := context.Background()

	// 测试空内容
	_, err := auditSvc.CheckContent(ctx, "")
	assert.Error(t, err, "空内容应该返回错误")
	assert.Contains(t, err.Error(), "不能为空")

	t.Log("空内容验证测试通过")
}

// TestContentAudit_SafeContent 测试安全内容
func TestContentAudit_SafeContent(t *testing.T) {
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockAuditRepo := new(MockAuditRecordRepository)
	mockViolationRepo := new(MockViolationRecordRepository)
	mockEventBus := new(MockEventBus)

	auditSvc := auditService.NewContentAuditService(
		mockSensitiveRepo,
		mockAuditRepo,
		mockViolationRepo,
		mockEventBus,
	)

	ctx := context.Background()

	mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*auditModels.SensitiveWord{}, nil)
	err := auditSvc.Initialize(ctx)
	assert.NoError(t, err)

	// 测试完全安全的内容
	content := "这是一篇正常的小说内容，讲述了一个温馨的故事。主人公经历了成长的挑战，最终收获了友谊和成功。"
	result, err := auditSvc.CheckContent(ctx, content)
	assert.NoError(t, err)

	// 验证结果
	assert.True(t, result.IsSafe, "安全内容应该被判定为安全")
	assert.Equal(t, 0, result.RiskLevel, "安全内容风险等级应该是0")
	assert.Equal(t, float64(0), result.RiskScore, "安全内容风险分数应该是0")
	assert.Empty(t, result.Violations, "安全内容不应有违规项")
	assert.True(t, result.CanPublish, "安全内容可以发布")
	assert.False(t, result.NeedsReview, "安全内容不需要复核")

	t.Log("安全内容测试通过")
}

// ==============================================
// 总结
// ==============================================
//
// 已实现测试（10个可运行）：
// - ✅ 政治敏感词匹配测试
// - ✅ 色情和暴力敏感词测试
// - ✅ 敏感词替换建议测试
// - ✅ 敏感词库更新测试
// - ✅ 手机号和URL检测测试
// - ✅ 微信和QQ号检测测试
// - ✅ 违规记录创建测试
// - ✅ 人工复审触发测试
// - ✅ 大文档性能测试（<3秒/万字）
// - ✅ DFA过滤器基础功能测试
// - ✅ 规则引擎复杂规则测试
// - ✅ 空内容和安全内容测试
//
// TDD待开发功能（2个）：
// - ⏸️ 审核结果缓存
// - ⏸️ 异步审核队列
//
// 待完善功能（1个）：
// - ⏸️ 批量审核（实现为TODO）
//
// 对应SRS需求：
// - REQ-AI-AGENT-003（内容安全检测）
// - REQ-CONTENT-AUDIT-001（敏感词检测）
// - REQ-CONTENT-AUDIT-002（合规性检查）
// - REQ-CONTENT-AUDIT-003（性能要求）
//
// ==============================================
