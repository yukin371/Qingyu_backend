package reading

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reading/reader"
	"Qingyu_backend/service/base"
	"Qingyu_backend/service/reading"
)

// MinimalMockAnnotationRepo 最小Mock注记Repository
type MinimalMockAnnotationRepo struct {
	mock.Mock
}

func (m *MinimalMockAnnotationRepo) Create(ctx context.Context, annotation *reader.Annotation) error {
	args := m.Called(ctx, annotation)
	return args.Error(0)
}

func (m *MinimalMockAnnotationRepo) GetByID(ctx context.Context, id string) (*reader.Annotation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Annotation), args.Error(1)
}

func (m *MinimalMockAnnotationRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MinimalMockAnnotationRepo) GetByBook(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MinimalMockAnnotationRepo) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return nil
}

func (m *MinimalMockAnnotationRepo) GetByChapter(ctx context.Context, userID, chapterID string) ([]*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) GetBookmarks(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) GetNotes(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) SearchNotes(ctx context.Context, userID, keyword string) ([]*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) Health(ctx context.Context) error {
	return nil
}

func (m *MinimalMockAnnotationRepo) BatchCreate(ctx context.Context, annotations []*reader.Annotation) error {
	args := m.Called(ctx, annotations)
	return args.Error(0)
}

func (m *MinimalMockAnnotationRepo) BatchDelete(ctx context.Context, annotationIDs []string) error {
	args := m.Called(ctx, annotationIDs)
	return args.Error(0)
}

func (m *MinimalMockAnnotationRepo) CountByBook(ctx context.Context, userID, bookID string) (int64, error) {
	return 0, nil
}

func (m *MinimalMockAnnotationRepo) CountByType(ctx context.Context, userID string, annotationType int) (int64, error) {
	return 0, nil
}

func (m *MinimalMockAnnotationRepo) CountByUser(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}

func (m *MinimalMockAnnotationRepo) DeleteByBook(ctx context.Context, userID, bookID string) error {
	return nil
}

func (m *MinimalMockAnnotationRepo) DeleteByChapter(ctx context.Context, userID, bookID, chapterID string) error {
	return nil
}

func (m *MinimalMockAnnotationRepo) GetBookmarkByPosition(ctx context.Context, userID, bookID, chapterID string, position int) (*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) GetByType(ctx context.Context, userID, bookID string, annotationType int) ([]*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) GetByUserAndBook(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MinimalMockAnnotationRepo) GetByUserAndChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) GetNotesByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) GetLatestBookmark(ctx context.Context, userID, bookID string) (*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) GetHighlights(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) GetHighlightsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) SyncAnnotations(ctx context.Context, userID string, annotations []*reader.Annotation) error {
	return nil
}

func (m *MinimalMockAnnotationRepo) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*reader.Annotation, error) {
	return nil, nil
}

func (m *MinimalMockAnnotationRepo) GetSharedAnnotations(ctx context.Context, userID string) ([]*reader.Annotation, error) {
	return nil, nil
}

// 创建测试ReaderService
func createTestReaderService(annotationRepo *MinimalMockAnnotationRepo) *reading.ReaderService {
	eventBus := base.NewSimpleEventBus()

	// 使用真实的Repository接口
	return reading.NewReaderService(
		nil,            // chapter repo
		nil,            // progress repo
		annotationRepo, // annotation repo
		nil,            // settings repo
		eventBus,
		nil, // cache service
		nil, // vip service
	)
}

// ========== 批量创建注记测试 ==========

// TestBatchCreateAnnotations_Success 测试批量创建注记成功
func TestBatchCreateAnnotations_Success(t *testing.T) {
	mockRepo := new(MinimalMockAnnotationRepo)
	service := createTestReaderService(mockRepo)
	ctx := context.Background()

	// 准备测试数据
	annotations := []*reader.Annotation{
		{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    "user-123",
			BookID:    "book-456",
			ChapterID: "chapter-789",
			Type:      "bookmark",
			Text:      "测试书签1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    "user-123",
			BookID:    "book-456",
			ChapterID: "chapter-790",
			Type:      "note",
			Text:      "测试笔记2",
			Note:      "这是我的笔记",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// 设置Mock期望
	mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Annotation")).Return(nil).Times(2)

	// 执行批量创建
	err := service.BatchCreateAnnotations(ctx, annotations)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBatchCreateAnnotations_Empty 测试批量创建空列表
func TestBatchCreateAnnotations_Empty(t *testing.T) {
	mockRepo := new(MinimalMockAnnotationRepo)
	service := createTestReaderService(mockRepo)
	ctx := context.Background()

	// 执行批量创建（空列表）
	err := service.BatchCreateAnnotations(ctx, []*reader.Annotation{})

	// 验证结果（应该直接返回nil，不调用repo）
	assert.NoError(t, err)
	mockRepo.AssertNotCalled(t, "Create")
}

// TestBatchCreateAnnotations_TooMany 测试批量创建过多注记
func TestBatchCreateAnnotations_TooMany(t *testing.T) {
	mockRepo := new(MinimalMockAnnotationRepo)
	service := createTestReaderService(mockRepo)
	ctx := context.Background()

	// 创建51个注记（超过限制）
	annotations := make([]*reader.Annotation, 51)
	for i := 0; i < 51; i++ {
		annotations[i] = &reader.Annotation{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    "user-123",
			BookID:    "book-456",
			ChapterID: "chapter-789",
			Type:      "bookmark",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	// 执行批量创建
	err := service.BatchCreateAnnotations(ctx, annotations)

	// 验证结果（应该返回错误）
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不能超过50个")
	mockRepo.AssertNotCalled(t, "Create")
}

// ========== 批量删除注记测试 ==========

// TestBatchDeleteAnnotations_Success 测试批量删除注记成功
func TestBatchDeleteAnnotations_Success(t *testing.T) {
	mockRepo := new(MinimalMockAnnotationRepo)
	service := createTestReaderService(mockRepo)
	ctx := context.Background()

	annotationIDs := []string{"ann-1", "ann-2", "ann-3"}

	// 设置Mock期望 - BatchDeleteAnnotations只调用Delete，不调用GetByID
	for _, id := range annotationIDs {
		mockRepo.On("Delete", mock.Anything, id).Return(nil).Once()
	}

	// 执行批量删除
	err := service.BatchDeleteAnnotations(ctx, annotationIDs)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBatchDeleteAnnotations_Empty 测试批量删除空列表
func TestBatchDeleteAnnotations_Empty(t *testing.T) {
	mockRepo := new(MinimalMockAnnotationRepo)
	service := createTestReaderService(mockRepo)
	ctx := context.Background()

	// 执行批量删除（空列表）
	err := service.BatchDeleteAnnotations(ctx, []string{})

	// 验证结果（应该直接返回nil）
	assert.NoError(t, err)
	mockRepo.AssertNotCalled(t, "Delete")
}

// TestBatchDeleteAnnotations_TooMany 测试批量删除过多注记
func TestBatchDeleteAnnotations_TooMany(t *testing.T) {
	mockRepo := new(MinimalMockAnnotationRepo)
	service := createTestReaderService(mockRepo)
	ctx := context.Background()

	// 创建101个ID（超过限制）
	annotationIDs := make([]string, 101)
	for i := 0; i < 101; i++ {
		annotationIDs[i] = primitive.NewObjectID().Hex()
	}

	// 执行批量删除
	err := service.BatchDeleteAnnotations(ctx, annotationIDs)

	// 验证结果（应该返回错误）
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不能超过100个")
	mockRepo.AssertNotCalled(t, "Delete")
}

// ========== 获取注记统计测试 ==========

// TestGetAnnotationStats_Success 测试获取注记统计成功
func TestGetAnnotationStats_Success(t *testing.T) {
	mockRepo := new(MinimalMockAnnotationRepo)
	service := createTestReaderService(mockRepo)
	ctx := context.Background()

	// 准备Mock数据
	annotations := []*reader.Annotation{
		{ID: "ann-1", Type: "bookmark", UserID: "user-123", BookID: "book-456"},
		{ID: "ann-2", Type: "bookmark", UserID: "user-123", BookID: "book-456"},
		{ID: "ann-3", Type: "highlight", UserID: "user-123", BookID: "book-456"},
		{ID: "ann-4", Type: "highlight", UserID: "user-123", BookID: "book-456"},
		{ID: "ann-5", Type: "highlight", UserID: "user-123", BookID: "book-456"},
		{ID: "ann-6", Type: "note", UserID: "user-123", BookID: "book-456"},
	}

	// 修复：应该Mock GetByUserAndBook而不是GetByBook
	mockRepo.On("GetByUserAndBook", mock.Anything, "user-123", "book-456").Return(annotations, nil)

	// 执行获取统计
	stats, err := service.GetAnnotationStats(ctx, "user-123", "book-456")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 6, stats["totalCount"])
	assert.Equal(t, 2, stats["bookmarkCount"])
	assert.Equal(t, 3, stats["highlightCount"])
	assert.Equal(t, 1, stats["noteCount"])

	mockRepo.AssertExpectations(t)
}

// TestGetAnnotationStats_Empty 测试获取空注记统计
func TestGetAnnotationStats_Empty(t *testing.T) {
	mockRepo := new(MinimalMockAnnotationRepo)
	service := createTestReaderService(mockRepo)
	ctx := context.Background()

	// 返回空列表 - 修复：使用GetByUserAndBook
	mockRepo.On("GetByUserAndBook", mock.Anything, "user-123", "book-456").Return([]*reader.Annotation{}, nil)

	// 执行获取统计
	stats, err := service.GetAnnotationStats(ctx, "user-123", "book-456")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 0, stats["totalCount"])
	assert.Equal(t, 0, stats["bookmarkCount"])
	assert.Equal(t, 0, stats["highlightCount"])
	assert.Equal(t, 0, stats["noteCount"])

	mockRepo.AssertExpectations(t)
}

// ========== 同步注记测试 ==========

// TestSyncAnnotations_Success 测试同步注记成功
func TestSyncAnnotations_Success(t *testing.T) {
	mockRepo := new(MinimalMockAnnotationRepo)
	service := createTestReaderService(mockRepo)
	ctx := context.Background()

	// 准备服务器端注记（模拟已有3条注记）
	serverAnnotations := []*reader.Annotation{
		{
			ID:        "ann-1",
			UserID:    "user-123",
			BookID:    "book-456",
			Type:      "bookmark",
			CreatedAt: time.Unix(1700000100, 0), // 比lastSyncTime新
		},
		{
			ID:        "ann-2",
			UserID:    "user-123",
			BookID:    "book-456",
			Type:      "note",
			CreatedAt: time.Unix(1699999900, 0), // 比lastSyncTime旧
		},
	}

	// 修复：使用GetByUserAndBook和mock.Anything
	mockRepo.On("GetByUserAndBook", mock.Anything, "user-123", "book-456").Return(serverAnnotations, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*reader.Annotation")).Return(nil).Once()

	// 准备同步请求
	syncReq := &reading.SyncAnnotationsRequest{
		BookID:       "book-456",
		LastSyncTime: 1700000000, // Unix时间戳
		LocalAnnotations: []*reader.Annotation{
			{
				Type:      "highlight",
				Text:      "本地新增高亮",
				BookID:    "book-456",
				ChapterID: "chapter-789",
			},
		},
	}

	// 执行同步
	result, err := service.SyncAnnotations(ctx, "user-123", syncReq)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result["uploadedCount"])   // 上传1条
	assert.Equal(t, 1, result["downloadedCount"]) // 下载1条（只有ann-1比lastSyncTime新）
	assert.Contains(t, result, "syncTime")

	mockRepo.AssertExpectations(t)
}

// TestSyncAnnotations_NoLocalChanges 测试无本地变更的同步
func TestSyncAnnotations_NoLocalChanges(t *testing.T) {
	mockRepo := new(MinimalMockAnnotationRepo)
	service := createTestReaderService(mockRepo)
	ctx := context.Background()

	// 准备服务器端注记
	serverAnnotations := []*reader.Annotation{
		{
			ID:        "ann-1",
			UserID:    "user-123",
			BookID:    "book-456",
			Type:      "bookmark",
			CreatedAt: time.Unix(1700000100, 0),
		},
	}

	// 修复：使用GetByUserAndBook和mock.Anything
	mockRepo.On("GetByUserAndBook", mock.Anything, "user-123", "book-456").Return(serverAnnotations, nil)

	// 准备同步请求（无本地注记）
	syncReq := &reading.SyncAnnotationsRequest{
		BookID:           "book-456",
		LastSyncTime:     1700000000,
		LocalAnnotations: []*reader.Annotation{}, // 空列表
	}

	// 执行同步
	result, err := service.SyncAnnotations(ctx, "user-123", syncReq)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result["uploadedCount"])
	assert.Equal(t, 1, result["downloadedCount"])

	mockRepo.AssertExpectations(t)
}

// ========== 性能基准测试 ==========

// BenchmarkBatchCreateAnnotations 批量创建性能测试
func BenchmarkBatchCreateAnnotations(b *testing.B) {
	mockRepo := new(MinimalMockAnnotationRepo)
	service := createTestReaderService(mockRepo)
	ctx := context.Background()

	// 准备测试数据
	annotations := make([]*reader.Annotation, 10)
	for i := 0; i < 10; i++ {
		annotations[i] = &reader.Annotation{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    "user-123",
			BookID:    "book-456",
			ChapterID: "chapter-789",
			Type:      "bookmark",
			CreatedAt: time.Now(),
		}
	}

	mockRepo.On("Create", ctx, mock.Anything).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.BatchCreateAnnotations(ctx, annotations)
	}
}

// BenchmarkGetAnnotationStats 获取统计性能测试
func BenchmarkGetAnnotationStats(b *testing.B) {
	mockRepo := new(MinimalMockAnnotationRepo)
	service := createTestReaderService(mockRepo)
	ctx := context.Background()

	// 准备测试数据（100条注记）
	annotations := make([]*reader.Annotation, 100)
	for i := 0; i < 100; i++ {
		annotations[i] = &reader.Annotation{
			ID:     primitive.NewObjectID().Hex(),
			Type:   []string{"bookmark", "highlight", "note"}[i%3],
			UserID: "user-123",
			BookID: "book-456",
		}
	}

	mockRepo.On("GetByBook", ctx, "user-123", "book-456").Return(annotations, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetAnnotationStats(ctx, "user-123", "book-456")
	}
}
