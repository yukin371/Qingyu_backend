package monitor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDatabase 模拟数据库接口
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) OrphanedRecords(ctx context.Context, collection, foreignKey, targetCollection string) (int, error) {
	args := m.Called(ctx, collection, foreignKey, targetCollection)
	return args.Int(0), args.Error(1)
}

func (m *MockDatabase) StatisticsAccuracy(ctx context.Context, collection, countField string) (int, error) {
	args := m.Called(ctx, collection, countField)
	return args.Int(0), args.Error(1)
}

// MockAlerter 模拟告警接口
type MockAlerter struct {
	mock.Mock
	alertSent bool
	message   string
}

func (m *MockAlerter) SendAlert(ctx context.Context, message string) error {
	args := m.Called(ctx, message)
	m.alertSent = true
	m.message = message
	return args.Error(0)
}

// TestDataQualityMonitor_CheckOrphanedRecords 测试孤儿记录检查
func TestDataQualityMonitor_CheckOrphanedRecords(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDB := new(MockDatabase)
	monitor := NewDataQualityMonitor(mockDB, nil)

	// Mock不同集合的孤儿记录
	mockDB.On("OrphanedRecords", ctx, "reading_progress", "user_id", "users").Return(10, nil)
	mockDB.On("OrphanedRecords", ctx, "reading_progress", "book_id", "books").Return(5, nil)
	mockDB.On("OrphanedRecords", ctx, "reading_history", "user_id", "users").Return(3, nil)
	mockDB.On("OrphanedRecords", ctx, "reading_history", "book_id", "books").Return(2, nil)
	mockDB.On("OrphanedRecords", ctx, "reading_history", "chapter_id", "chapters").Return(1, nil)
	mockDB.On("OrphanedRecords", ctx, "bookmarks", "user_id", "users").Return(0, nil)
	mockDB.On("OrphanedRecords", ctx, "bookmarks", "book_id", "books").Return(0, nil)
	mockDB.On("OrphanedRecords", ctx, "bookmarks", "chapter_id", "chapters").Return(0, nil)
	mockDB.On("OrphanedRecords", ctx, "likes", "user_id", "users").Return(0, nil)
	mockDB.On("OrphanedRecords", ctx, "likes", "target_id", mock.Anything, mock.Anything).Return(0, nil)
	mockDB.On("OrphanedRecords", ctx, "notifications", "user_id", "users").Return(0, nil)
	mockDB.On("OrphanedRecords", ctx, "author_revenue", "user_id", "users").Return(0, nil)

	// Act
	result, err := monitor.CheckOrphanedRecords(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 21, result.TotalOrphanedRecords)
	assert.True(t, result.HasIssues)
	assert.Equal(t, 5, len(result.Details), "只有count>0的记录才会被添加到Details")
	mockDB.AssertExpectations(t)
}

// TestDataQualityMonitor_CheckOrphanedRecords_NoOrphans 测试无孤儿记录的情况
func TestDataQualityMonitor_CheckOrphanedRecords_NoOrphans(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDB := new(MockDatabase)
	monitor := NewDataQualityMonitor(mockDB, nil)

	// Mock无孤儿记录
	mockDB.On("OrphanedRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, nil)

	// Act
	result, err := monitor.CheckOrphanedRecords(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, result.TotalOrphanedRecords)
	assert.False(t, result.HasIssues)
}

// TestDataQualityMonitor_CheckOrphanedRecords_Error 测试检查出错的情况
func TestDataQualityMonitor_CheckOrphanedRecords_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDB := new(MockDatabase)
	monitor := NewDataQualityMonitor(mockDB, nil)

	mockDB.On("OrphanedRecords", ctx, "reading_progress", "user_id", "users").Return(0, errors.New("database error"))

	// Act
	result, err := monitor.CheckOrphanedRecords(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "检查")
}

// TestDataQualityMonitor_CheckStatisticsAccuracy 测试统计数据准确性检查
func TestDataQualityMonitor_CheckStatisticsAccuracy(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDB := new(MockDatabase)
	monitor := NewDataQualityMonitor(mockDB, nil)

	// Mock统计数据不准确的情况
	mockDB.On("StatisticsAccuracy", ctx, "books", "likes_count").Return(15, nil)
	mockDB.On("StatisticsAccuracy", ctx, "books", "comments_count").Return(8, nil)
	mockDB.On("StatisticsAccuracy", ctx, "users", "followers_count").Return(5, nil)

	// Act
	result, err := monitor.CheckStatisticsAccuracy(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 28, result.InaccurateStatistics)
	assert.True(t, result.HasIssues)
	assert.Equal(t, 3, len(result.Details))
	mockDB.AssertExpectations(t)
}

// TestDataQualityMonitor_CheckStatisticsAccuracy_NoIssues 测试统计数据准确的情况
func TestDataQualityMonitor_CheckStatisticsAccuracy_NoIssues(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDB := new(MockDatabase)
	monitor := NewDataQualityMonitor(mockDB, nil)

	mockDB.On("StatisticsAccuracy", ctx, "books", "likes_count").Return(0, nil)
	mockDB.On("StatisticsAccuracy", ctx, "books", "comments_count").Return(0, nil)
	mockDB.On("StatisticsAccuracy", ctx, "users", "followers_count").Return(0, nil)

	// Act
	result, err := monitor.CheckStatisticsAccuracy(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, result.InaccurateStatistics)
	assert.False(t, result.HasIssues)
}

// TestDataQualityMonitor_GenerateReport 测试报告生成
func TestDataQualityMonitor_GenerateReport(t *testing.T) {
	// Arrange
	report := &DataQualityReport{
		CheckTime:            time.Now(),
		TotalOrphanedRecords: 100,
		InaccurateStatistics: 50,
		Details: []CheckResult{
			{Collection: "reading_progress", Field: "user_id", Count: 50},
			{Collection: "reading_progress", Field: "book_id", Count: 30},
			{Collection: "books", Field: "likes_count", Count: 20},
		},
		HasIssues: true,
	}

	// Act
	json := report.ToJSON()

	// Assert
	assert.Contains(t, json, "total_orphaned_records")
	assert.Contains(t, json, "100")
	assert.Contains(t, json, "inaccurate_statistics")
	assert.Contains(t, json, "50")
	assert.Contains(t, json, "has_issues")
	assert.Contains(t, json, "true")
}

// TestDataQualityMonitor_ShouldAlert 测试告警判断
func TestDataQualityMonitor_ShouldAlert(t *testing.T) {
	monitor := NewDataQualityMonitor(nil, nil)

	tests := []struct {
		name     string
		report   *DataQualityReport
		expected bool
	}{
		{
			name: "有孤儿记录应该告警",
			report: &DataQualityReport{
				TotalOrphanedRecords: 100,
				HasIssues:            true,
			},
			expected: true,
		},
		{
			name: "有统计数据问题应该告警",
			report: &DataQualityReport{
				TotalOrphanedRecords: 0,
				InaccurateStatistics: 50,
				HasIssues:           true,
			},
			expected: true,
		},
		{
			name: "无问题不应告警",
			report: &DataQualityReport{
				TotalOrphanedRecords: 0,
				InaccurateStatistics: 0,
				HasIssues:           false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := monitor.ShouldAlert(tt.report)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDataQualityMonitor_RunDailyCheck 测试每日检查
func TestDataQualityMonitor_RunDailyCheck(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDB := new(MockDatabase)
	mockAlert := new(MockAlerter)
	monitor := NewDataQualityMonitor(mockDB, mockAlert)

	// Mock孤儿记录检查
	mockDB.On("OrphanedRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, nil)

	// Mock统计数据检查
	mockDB.On("StatisticsAccuracy", ctx, "books", "likes_count").Return(10, nil)
	mockDB.On("StatisticsAccuracy", ctx, "books", "comments_count").Return(0, nil)
	mockDB.On("StatisticsAccuracy", ctx, "users", "followers_count").Return(0, nil)

	// Mock告警
	mockAlert.On("SendAlert", ctx, mock.AnythingOfType("string")).Return(nil)

	// Act
	result, err := monitor.RunDailyCheck(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.InaccurateStatistics)
	assert.True(t, result.HasIssues)
	assert.True(t, mockAlert.alertSent)
	assert.Contains(t, mockAlert.message, "统计不准确")
	mockDB.AssertExpectations(t)
}

// TestDataQualityMonitor_RunDailyCheck_NoAlert 测试每日检查无告警
func TestDataQualityMonitor_RunDailyCheck_NoAlert(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDB := new(MockDatabase)
	mockAlert := new(MockAlerter)
	monitor := NewDataQualityMonitor(mockDB, mockAlert)

	// Mock无问题
	mockDB.On("OrphanedRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, nil)
	mockDB.On("StatisticsAccuracy", mock.Anything, mock.Anything, mock.Anything).Return(0, nil)

	// Act
	result, err := monitor.RunDailyCheck(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.TotalOrphanedRecords)
	assert.Equal(t, 0, result.InaccurateStatistics)
	assert.False(t, result.HasIssues)
	assert.False(t, mockAlert.alertSent)
	mockDB.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything)
}

// TestDataQualityMonitor_FormatAlertMessage 测试告警消息格式化
func TestDataQualityMonitor_FormatAlertMessage(t *testing.T) {
	// Arrange
	monitor := NewDataQualityMonitor(nil, nil)
	report := &DataQualityReport{
		CheckTime:           time.Now(),
		TotalOrphanedRecords: 100,
		InaccurateStatistics: 50,
		HasIssues:           true,
	}

	// Act
	message := monitor.formatAlertMessage(report)

	// Assert
	assert.Contains(t, message, "数据质量")
	assert.Contains(t, message, "孤儿记录")
	assert.Contains(t, message, "100")
	assert.Contains(t, message, "统计不准确")
	assert.Contains(t, message, "50")
}
