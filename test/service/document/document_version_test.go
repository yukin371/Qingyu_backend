package document

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/service/base"
	documentService "Qingyu_backend/service/document"
)

// MockDocumentRepository Mock文档Repository
type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) GetByID(ctx context.Context, id string) (*writer.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Document), args.Error(1)
}

func (m *MockDocumentRepository) Create(ctx context.Context, doc *writer.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockDocumentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDocumentRepository) List(ctx context.Context, filter interface{}) ([]*writer.Document, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Document), args.Error(1)
}

func (m *MockDocumentRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockDocumentContentRepository Mock文档内容Repository
type MockDocumentContentRepository struct {
	mock.Mock
}

func (m *MockDocumentContentRepository) GetByDocumentID(ctx context.Context, documentID string) (*writer.DocumentContent, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.DocumentContent), args.Error(1)
}

func (m *MockDocumentContentRepository) Create(ctx context.Context, content *writer.DocumentContent) error {
	args := m.Called(ctx, content)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) Update(ctx context.Context, documentID string, updates map[string]interface{}) error {
	args := m.Called(ctx, documentID, updates)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockProjectRepository Mock项目Repository
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Health(ctx context.Context) error {
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

// TestDocumentVersionCalculation 测试文档版本号计算
func TestDocumentVersionCalculation(t *testing.T) {
	ctx := context.Background()

	t.Run("新创建的文档版本号为1", func(t *testing.T) {
		mockDocRepo := new(MockDocumentRepository)
		mockContentRepo := new(MockDocumentContentRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockEventBus := new(MockEventBus)

		service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

		// 创建刚好创建的文档
		doc := &writer.Document{
			ID:        "doc1",
			Title:     "测试文档",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockDocRepo.On("GetByID", ctx, "doc1").Return(doc, nil)

		status, err := service.GetSaveStatus(ctx, "doc1")

		assert.NoError(t, err)
		assert.Equal(t, 1, status.CurrentVersion)
	})

	t.Run("更新5分钟后版本号为2", func(t *testing.T) {
		mockDocRepo := new(MockDocumentRepository)
		mockContentRepo := new(MockDocumentContentRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockEventBus := new(MockEventBus)

		service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

		now := time.Now()
		doc := &writer.Document{
			ID:        "doc2",
			Title:     "测试文档2",
			CreatedAt: now.Add(-5 * time.Minute),
			UpdatedAt: now,
		}

		mockDocRepo.On("GetByID", ctx, "doc2").Return(doc, nil)

		status, err := service.GetSaveStatus(ctx, "doc2")

		assert.NoError(t, err)
		assert.Equal(t, 2, status.CurrentVersion)
	})

	t.Run("更新25分钟后版本号为6", func(t *testing.T) {
		mockDocRepo := new(MockDocumentRepository)
		mockContentRepo := new(MockDocumentContentRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockEventBus := new(MockEventBus)

		service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

		now := time.Now()
		doc := &writer.Document{
			ID:        "doc3",
			Title:     "测试文档3",
			CreatedAt: now.Add(-25 * time.Minute),
			UpdatedAt: now,
		}

		mockDocRepo.On("GetByID", ctx, "doc3").Return(doc, nil)

		status, err := service.GetSaveStatus(ctx, "doc3")

		assert.NoError(t, err)
		assert.Equal(t, 6, status.CurrentVersion) // 25/5 + 1 = 6
	})

	t.Run("版本号最大为999", func(t *testing.T) {
		mockDocRepo := new(MockDocumentRepository)
		mockContentRepo := new(MockDocumentContentRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockEventBus := new(MockEventBus)

		service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

		now := time.Now()
		doc := &writer.Document{
			ID:        "doc4",
			Title:     "测试文档4",
			CreatedAt: now.Add(-10000 * time.Minute), // 非常久以前
			UpdatedAt: now,
		}

		mockDocRepo.On("GetByID", ctx, "doc4").Return(doc, nil)

		status, err := service.GetSaveStatus(ctx, "doc4")

		assert.NoError(t, err)
		assert.LessOrEqual(t, status.CurrentVersion, 999) // 不应超过999
	})

	t.Run("版本号最小为1", func(t *testing.T) {
		mockDocRepo := new(MockDocumentRepository)
		mockContentRepo := new(MockDocumentContentRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockEventBus := new(MockEventBus)

		service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

		now := time.Now()
		doc := &writer.Document{
			ID:        "doc5",
			Title:     "测试文档5",
			CreatedAt: now.Add(1 * time.Minute), // 未来时间（错误数据）
			UpdatedAt: now,
		}

		mockDocRepo.On("GetByID", ctx, "doc5").Return(doc, nil)

		status, err := service.GetSaveStatus(ctx, "doc5")

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, status.CurrentVersion, 1) // 至少为1
	})

	t.Run("文档不存在", func(t *testing.T) {
		mockDocRepo := new(MockDocumentRepository)
		mockContentRepo := new(MockDocumentContentRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockEventBus := new(MockEventBus)

		service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

		mockDocRepo.On("GetByID", ctx, "notexist").Return(nil, assert.AnError)

		status, err := service.GetSaveStatus(ctx, "notexist")

		assert.Error(t, err)
		assert.Nil(t, status)
	})

	t.Run("批量文档版本计算", func(t *testing.T) {
		mockDocRepo := new(MockDocumentRepository)
		mockContentRepo := new(MockDocumentContentRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockEventBus := new(MockEventBus)

		service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

		now := time.Now()
		testCases := []struct {
			docID           string
			minutesAgo      int
			expectedVersion int
		}{
			{"doc_a", 0, 1},
			{"doc_b", 5, 2},
			{"doc_c", 10, 3},
			{"doc_d", 15, 4},
			{"doc_e", 20, 5},
		}

		for _, tc := range testCases {
			doc := &writer.Document{
				ID:        tc.docID,
				Title:     "测试文档",
				CreatedAt: now.Add(-time.Duration(tc.minutesAgo) * time.Minute),
				UpdatedAt: now,
			}

			mockDocRepo.On("GetByID", ctx, tc.docID).Return(doc, nil).Once()

			status, err := service.GetSaveStatus(ctx, tc.docID)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedVersion, status.CurrentVersion,
				"文档 %s 期望版本 %d，实际 %d", tc.docID, tc.expectedVersion, status.CurrentVersion)
		}
	})
}

// TestGetSaveStatusIntegration 集成测试
func TestGetSaveStatusIntegration(t *testing.T) {
	ctx := context.Background()

	t.Run("完整的保存状态信息", func(t *testing.T) {
		mockDocRepo := new(MockDocumentRepository)
		mockContentRepo := new(MockDocumentContentRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockEventBus := new(MockEventBus)

		service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

		now := time.Now()
		doc := &writer.Document{
			ID:        "full_doc",
			Title:     "完整测试文档",
			CreatedAt: now.Add(-10 * time.Minute),
			UpdatedAt: now,
			WordCount: 1500,
		}

		mockDocRepo.On("GetByID", ctx, "full_doc").Return(doc, nil)

		status, err := service.GetSaveStatus(ctx, "full_doc")

		assert.NoError(t, err)
		assert.Equal(t, "full_doc", status.DocumentID)
		assert.Equal(t, now, status.LastSavedAt)
		assert.Equal(t, 3, status.CurrentVersion) // 10分钟 / 5 + 1
		assert.Equal(t, false, status.IsSaving)
		assert.Equal(t, 1500, status.WordCount)
	})

	t.Run("空文档ID", func(t *testing.T) {
		mockDocRepo := new(MockDocumentRepository)
		mockContentRepo := new(MockDocumentContentRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockEventBus := new(MockEventBus)

		service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

		status, err := service.GetSaveStatus(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, status)
		assert.Contains(t, err.Error(), "不能为空")
	})
}

// BenchmarkCalculateDocumentVersion 基准测试版本计算性能
func BenchmarkCalculateDocumentVersion(b *testing.B) {
	ctx := context.Background()
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	now := time.Now()
	doc := &writer.Document{
		ID:        "bench_doc",
		Title:     "基准测试文档",
		CreatedAt: now.Add(-30 * time.Minute),
		UpdatedAt: now,
		WordCount: 2000,
	}

	mockDocRepo.On("GetByID", ctx, "bench_doc").Return(doc, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetSaveStatus(ctx, "bench_doc")
	}
}
