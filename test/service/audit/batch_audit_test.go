package audit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/models/audit"
	auditService "Qingyu_backend/service/audit"
	"Qingyu_backend/service/base"
)

// MockAuditRecordRepository Mock审核记录Repository
type MockAuditRecordRepository struct {
	mock.Mock
}

func (m *MockAuditRecordRepository) Create(ctx context.Context, record *audit.AuditRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockAuditRecordRepository) GetByID(ctx context.Context, id string) (*audit.AuditRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*audit.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockAuditRecordRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAuditRecordRepository) List(ctx context.Context, filter interface{}) ([]*audit.AuditRecord, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) Count(ctx context.Context, filter interface{}) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuditRecordRepository) FindWithPagination(ctx context.Context, filter, pagination interface{}) (interface{}, error) {
	args := m.Called(ctx, filter, pagination)
	return args.Get(0), args.Error(1)
}

func (m *MockAuditRecordRepository) GetByTargetID(ctx context.Context, targetType, targetID string) (*audit.AuditRecord, error) {
	args := m.Called(ctx, targetType, targetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*audit.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) GetByAuthor(ctx context.Context, authorID string, limit int64) ([]*audit.AuditRecord, error) {
	args := m.Called(ctx, authorID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) GetByStatus(ctx context.Context, status string, limit int64) ([]*audit.AuditRecord, error) {
	args := m.Called(ctx, status, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) GetByResult(ctx context.Context, result string, limit int64) ([]*audit.AuditRecord, error) {
	args := m.Called(ctx, result, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) GetPendingReview(ctx context.Context, limit int64) ([]*audit.AuditRecord, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) GetHighRisk(ctx context.Context, minRiskLevel int, limit int64) ([]*audit.AuditRecord, error) {
	args := m.Called(ctx, minRiskLevel, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) UpdateStatus(ctx context.Context, id string, status string, reviewerID string, note string) error {
	args := m.Called(ctx, id, status, reviewerID, note)
	return args.Error(0)
}

func (m *MockAuditRecordRepository) UpdateAppealStatus(ctx context.Context, id string, appealStatus string) error {
	args := m.Called(ctx, id, appealStatus)
	return args.Error(0)
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

func (m *MockAuditRecordRepository) GetRecentRecords(ctx context.Context, targetType string, limit int64) ([]*audit.AuditRecord, error) {
	args := m.Called(ctx, targetType, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.AuditRecord), args.Error(1)
}

func (m *MockAuditRecordRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockSensitiveWordRepository Mock敏感词Repository
type MockSensitiveWordRepository struct {
	mock.Mock
}

func (m *MockSensitiveWordRepository) GetEnabledWords(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockSensitiveWordRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockViolationRecordRepository Mock违规记录Repository
type MockViolationRecordRepository struct {
	mock.Mock
}

func (m *MockViolationRecordRepository) Create(ctx context.Context, record interface{}) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockViolationRecordRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
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

func (m *MockEventBus) Publish(event base.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(event base.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

// TestBatchAuditDocuments 测试批量审核文档
func TestBatchAuditDocuments(t *testing.T) {
	ctx := context.Background()

	t.Run("成功批量审核", func(t *testing.T) {
		mockAuditRepo := new(MockAuditRecordRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockViolationRepo := new(MockViolationRecordRepository)
		mockEventBus := new(MockEventBus)

		// 创建服务
		service := auditService.NewContentAuditService(
			mockSensitiveRepo,
			mockAuditRepo,
			mockViolationRepo,
			mockEventBus,
		)

		// Mock Create 方法
		mockAuditRepo.On("Create", ctx, mock.AnythingOfType("*audit.AuditRecord")).
			Return(nil).Times(3)

		// 批量审核3个文档
		documentIDs := []string{"doc1", "doc2", "doc3"}
		records, err := service.BatchAuditDocuments(ctx, documentIDs)

		assert.NoError(t, err)
		assert.Len(t, records, 3)

		// 验证每个记录
		for i, record := range records {
			assert.Equal(t, documentIDs[i], record.TargetID)
			assert.Equal(t, audit.TargetTypeDocument, record.TargetType)
			assert.Equal(t, audit.StatusPending, record.Status)
		}

		mockAuditRepo.AssertExpectations(t)
	})

	t.Run("空文档ID列表", func(t *testing.T) {
		mockAuditRepo := new(MockAuditRecordRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockViolationRepo := new(MockViolationRecordRepository)
		mockEventBus := new(MockEventBus)

		service := auditService.NewContentAuditService(
			mockSensitiveRepo,
			mockAuditRepo,
			mockViolationRepo,
			mockEventBus,
		)

		records, err := service.BatchAuditDocuments(ctx, []string{})

		assert.Error(t, err)
		assert.Nil(t, records)
		assert.Contains(t, err.Error(), "不能为空")
	})

	t.Run("部分审核失败", func(t *testing.T) {
		mockAuditRepo := new(MockAuditRecordRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockViolationRepo := new(MockViolationRecordRepository)
		mockEventBus := new(MockEventBus)

		service := auditService.NewContentAuditService(
			mockSensitiveRepo,
			mockAuditRepo,
			mockViolationRepo,
			mockEventBus,
		)

		// Mock: 第2个Create失败
		mockAuditRepo.On("Create", ctx, mock.MatchedBy(func(r *audit.AuditRecord) bool {
			return r.TargetID == "doc1" || r.TargetID == "doc3"
		})).Return(nil)

		mockAuditRepo.On("Create", ctx, mock.MatchedBy(func(r *audit.AuditRecord) bool {
			return r.TargetID == "doc2"
		})).Return(assert.AnError)

		documentIDs := []string{"doc1", "doc2", "doc3"}
		records, err := service.BatchAuditDocuments(ctx, documentIDs)

		// 应该返回成功的记录
		assert.NoError(t, err)    // 部分成功不应返回错误
		assert.Len(t, records, 2) // 只有2个成功

		mockAuditRepo.AssertExpectations(t)
	})

	t.Run("大批量审核", func(t *testing.T) {
		mockAuditRepo := new(MockAuditRecordRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockViolationRepo := new(MockViolationRecordRepository)
		mockEventBus := new(MockEventBus)

		service := auditService.NewContentAuditService(
			mockSensitiveRepo,
			mockAuditRepo,
			mockViolationRepo,
			mockEventBus,
		)

		// 生成100个文档ID
		documentIDs := make([]string, 100)
		for i := 0; i < 100; i++ {
			documentIDs[i] = "doc" + string(rune(i))
		}

		mockAuditRepo.On("Create", ctx, mock.AnythingOfType("*audit.AuditRecord")).
			Return(nil).Times(100)

		start := time.Now()
		records, err := service.BatchAuditDocuments(ctx, documentIDs)
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.Len(t, records, 100)

		// 并发处理应该快于串行（粗略估计）
		t.Logf("批量审核100个文档耗时: %v", duration)

		mockAuditRepo.AssertExpectations(t)
	})
}

// BenchmarkBatchAuditDocuments 基准测试批量审核性能
func BenchmarkBatchAuditDocuments(b *testing.B) {
	ctx := context.Background()
	mockAuditRepo := new(MockAuditRecordRepository)
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockViolationRepo := new(MockViolationRecordRepository)
	mockEventBus := new(MockEventBus)

	service := auditService.NewContentAuditService(
		mockSensitiveRepo,
		mockAuditRepo,
		mockViolationRepo,
		mockEventBus,
	)

	mockAuditRepo.On("Create", ctx, mock.AnythingOfType("*audit.AuditRecord")).
		Return(nil)

	documentIDs := make([]string, 10)
	for i := 0; i < 10; i++ {
		documentIDs[i] = "benchmark_doc" + string(rune(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.BatchAuditDocuments(ctx, documentIDs)
	}
}
