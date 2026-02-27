package admin

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/admin"
)

// ==================== Mock ExportHistoryRepository ====================

// MockExportHistoryRepository 导出历史仓储Mock
type MockExportHistoryRepository struct {
	mock.Mock
}

func (m *MockExportHistoryRepository) Create(ctx context.Context, history *admin.ExportHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockExportHistoryRepository) GetByID(ctx context.Context, id string) (*admin.ExportHistory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*admin.ExportHistory), args.Error(1)
}

func (m *MockExportHistoryRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockExportHistoryRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockExportHistoryRepository) List(ctx context.Context, filter *ExportHistoryFilter) ([]*admin.ExportHistory, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*admin.ExportHistory), args.Get(1).(int64), args.Error(2)
}

func (m *MockExportHistoryRepository) Count(ctx context.Context, filter *ExportHistoryFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockExportHistoryRepository) ListByUser(ctx context.Context, adminID string, page, pageSize int) ([]*admin.ExportHistory, int64, error) {
	args := m.Called(ctx, adminID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*admin.ExportHistory), args.Get(1).(int64), args.Error(2)
}

func (m *MockExportHistoryRepository) ListByDateRange(ctx context.Context, startDate, endDate time.Time, page, pageSize int) ([]*admin.ExportHistory, int64, error) {
	args := m.Called(ctx, startDate, endDate, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*admin.ExportHistory), args.Get(1).(int64), args.Error(2)
}

func (m *MockExportHistoryRepository) CleanOldRecords(ctx context.Context, beforeDate time.Time) (int64, error) {
	args := m.Called(ctx, beforeDate)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockExportHistoryRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ==================== 测试辅助函数 ====================

// 创建测试导出历史记录
func createTestExportHistory(id, adminID, exportType, format string) *admin.ExportHistory {
	now := time.Now()
	return &admin.ExportHistory{
		ID:          id,
		AdminID:     adminID,
		ExportType:  exportType,
		Format:      format,
		RecordCount: 100,
		FilePath:    "/exports/test.csv",
		FileSize:    1024,
		Status:      admin.ExportStatusCompleted,
		CreatedAt:   now,
		CompletedAt: &now,
	}
}

// ==================== TestExportHistoryRepository_Create_Success ====================

func TestExportHistoryRepository_Create_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	mockRepo := new(MockExportHistoryRepository)

	history := createTestExportHistory(primitive.NewObjectID().Hex(), "admin-1", admin.ExportTypeBooks, admin.ExportFormatCSV)
	mockRepo.On("Create", ctx, mock.MatchedBy(func(h *admin.ExportHistory) bool {
		return h.AdminID == "admin-1" && h.ExportType == admin.ExportTypeBooks
	})).Return(nil)

	// When
	err := mockRepo.Create(ctx, history)

	// Then
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ==================== TestExportHistoryRepository_ListByUser_Success ====================

func TestExportHistoryRepository_ListByUser_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	mockRepo := new(MockExportHistoryRepository)

	adminID := "admin-1"
	expectedHistories := []*admin.ExportHistory{
		createTestExportHistory("1", adminID, admin.ExportTypeBooks, admin.ExportFormatCSV),
		createTestExportHistory("2", adminID, admin.ExportTypeChapters, admin.ExportFormatExcel),
	}

	mockRepo.On("ListByUser", ctx, adminID, 1, 10).Return(expectedHistories, int64(2), nil)

	// When
	histories, total, err := mockRepo.ListByUser(ctx, adminID, 1, 10)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, len(histories))
	assert.Equal(t, adminID, histories[0].AdminID)

	mockRepo.AssertExpectations(t)
}

// ==================== TestExportHistoryRepository_ListByDateRange_Success ====================

func TestExportHistoryRepository_ListByDateRange_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	mockRepo := new(MockExportHistoryRepository)

	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	expectedHistories := []*admin.ExportHistory{
		createTestExportHistory("1", "admin-1", admin.ExportTypeBooks, admin.ExportFormatCSV),
	}

	mockRepo.On("ListByDateRange", ctx, mock.MatchedBy(func(start time.Time) bool {
		return start.Before(endDate)
	}), mock.MatchedBy(func(end time.Time) bool {
		return end.After(startDate)
	}), 1, 10).Return(expectedHistories, int64(1), nil)

	// When
	histories, total, err := mockRepo.ListByDateRange(ctx, startDate, endDate, 1, 10)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, len(histories))

	mockRepo.AssertExpectations(t)
}

// ==================== TestExportHistoryRepository_CleanOldRecords_Success ====================

func TestExportHistoryRepository_CleanOldRecords_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	mockRepo := new(MockExportHistoryRepository)

	beforeDate := time.Now().AddDate(0, -1, 0)
	mockRepo.On("CleanOldRecords", ctx, mock.MatchedBy(func(date time.Time) bool {
		return date.Before(time.Now())
	})).Return(int64(5), nil)

	// When
	deletedCount, err := mockRepo.CleanOldRecords(ctx, beforeDate)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, int64(5), deletedCount)

	mockRepo.AssertExpectations(t)
}

// ==================== TestExportHistoryRepository_GetByID_Success ====================

func TestExportHistoryRepository_GetByID_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	mockRepo := new(MockExportHistoryRepository)

	id := primitive.NewObjectID().Hex()
	expectedHistory := createTestExportHistory(id, "admin-1", admin.ExportTypeBooks, admin.ExportFormatCSV)

	mockRepo.On("GetByID", ctx, id).Return(expectedHistory, nil)

	// When
	history, err := mockRepo.GetByID(ctx, id)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, history)
	assert.Equal(t, id, history.ID)
	assert.Equal(t, "admin-1", history.AdminID)

	mockRepo.AssertExpectations(t)
}

// ==================== TestExportHistoryRepository_Update_Success ====================

func TestExportHistoryRepository_Update_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	mockRepo := new(MockExportHistoryRepository)

	id := primitive.NewObjectID().Hex()
	updates := map[string]interface{}{
		"status":       admin.ExportStatusCompleted,
		"record_count": 200,
	}

	mockRepo.On("Update", ctx, id, updates).Return(nil)

	// When
	err := mockRepo.Update(ctx, id, updates)

	// Then
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// ==================== TestExportHistoryRepository_Delete_Success ====================

func TestExportHistoryRepository_Delete_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	mockRepo := new(MockExportHistoryRepository)

	id := primitive.NewObjectID().Hex()
	mockRepo.On("Delete", ctx, id).Return(nil)

	// When
	err := mockRepo.Delete(ctx, id)

	// Then
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// ==================== TestExportHistoryRepository_List_Success ====================

func TestExportHistoryRepository_List_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	mockRepo := new(MockExportHistoryRepository)

	filter := &ExportHistoryFilter{
		AdminID:    "admin-1",
		ExportType: admin.ExportTypeBooks,
	}

	expectedHistories := []*admin.ExportHistory{
		createTestExportHistory("1", "admin-1", admin.ExportTypeBooks, admin.ExportFormatCSV),
	}

	mockRepo.On("List", ctx, mock.MatchedBy(func(f *ExportHistoryFilter) bool {
		return f.AdminID == "admin-1" && f.ExportType == admin.ExportTypeBooks
	})).Return(expectedHistories, int64(1), nil)

	// When
	histories, total, err := mockRepo.List(ctx, filter)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, len(histories))

	mockRepo.AssertExpectations(t)
}

// ==================== TestExportHistoryRepository_Count_Success ====================

func TestExportHistoryRepository_Count_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	mockRepo := new(MockExportHistoryRepository)

	filter := &ExportHistoryFilter{
		AdminID: "admin-1",
	}

	mockRepo.On("Count", ctx, filter).Return(int64(10), nil)

	// When
	count, err := mockRepo.Count(ctx, filter)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)

	mockRepo.AssertExpectations(t)
}
