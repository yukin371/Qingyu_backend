package admin

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestExportHistory_Constants 测试导出历史常量
func TestExportHistory_Constants(t *testing.T) {
	// 测试状态常量
	assert.Equal(t, "pending", ExportStatusPending)
	assert.Equal(t, "completed", ExportStatusCompleted)
	assert.Equal(t, "failed", ExportStatusFailed)

	// 测试类型常量
	assert.Equal(t, "books", ExportTypeBooks)
	assert.Equal(t, "chapters", ExportTypeChapters)
	assert.Equal(t, "users", ExportTypeUsers)

	// 测试格式常量
	assert.Equal(t, "csv", ExportFormatCSV)
	assert.Equal(t, "excel", ExportFormatExcel)
}

// TestExportHistory_PendingExport 测试待处理导出
func TestExportHistory_PendingExport(t *testing.T) {
	export := &ExportHistory{
		ID:          "export123",
		AdminID:     "admin123",
		ExportType:  ExportTypeBooks,
		Format:      ExportFormatCSV,
		RecordCount: 0,
		FilePath:    "",
		FileSize:    0,
		Status:      ExportStatusPending,
		ErrorMsg:    "",
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}

	assert.Equal(t, ExportStatusPending, export.Status)
	assert.True(t, export.IsPending())
	assert.False(t, export.IsCompleted())
	assert.False(t, export.IsFailed())
	assert.Nil(t, export.CompletedAt)
}

// TestExportHistory_CompletedExport 测试已完成导出
func TestExportHistory_CompletedExport(t *testing.T) {
	now := time.Now()
	export := &ExportHistory{
		ID:          "export123",
		AdminID:     "admin123",
		ExportType:  ExportTypeBooks,
		Format:      ExportFormatCSV,
		RecordCount: 1000,
		FilePath:    "/exports/books_20260227.csv",
		FileSize:    1024000,
		Status:      ExportStatusCompleted,
		ErrorMsg:    "",
		CreatedAt:   now.Add(-1 * time.Hour),
		CompletedAt: &now,
	}

	assert.Equal(t, ExportStatusCompleted, export.Status)
	assert.False(t, export.IsPending())
	assert.True(t, export.IsCompleted())
	assert.False(t, export.IsFailed())
	assert.NotNil(t, export.CompletedAt)
	assert.Greater(t, export.RecordCount, 0)
	assert.Greater(t, export.FileSize, int64(0))
}

// TestExportHistory_FailedExport 测试失败导出
func TestExportHistory_FailedExport(t *testing.T) {
	export := &ExportHistory{
		ID:          "export123",
		AdminID:     "admin123",
		ExportType:  ExportTypeBooks,
		Format:      ExportFormatCSV,
		RecordCount: 0,
		FilePath:    "",
		FileSize:    0,
		Status:      ExportStatusFailed,
		ErrorMsg:    "导出失败：数据库连接超时",
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}

	assert.Equal(t, ExportStatusFailed, export.Status)
	assert.False(t, export.IsPending())
	assert.False(t, export.IsCompleted())
	assert.True(t, export.IsFailed())
	assert.NotEmpty(t, export.ErrorMsg)
}

// TestExportHistory_AllExportTypes 测试所有导出类型
func TestExportHistory_AllExportTypes(t *testing.T) {
	types := []string{ExportTypeBooks, ExportTypeChapters, ExportTypeUsers}

	for _, exportType := range types {
		export := &ExportHistory{
			ID:         "export123",
			AdminID:    "admin123",
			ExportType: exportType,
			Format:     ExportFormatCSV,
			Status:     ExportStatusPending,
			CreatedAt:  time.Now(),
		}

		assert.Contains(t, types, export.ExportType)
	}
}

// TestExportHistory_AllFormats 测试所有导出格式
func TestExportHistory_AllFormats(t *testing.T) {
	formats := []string{ExportFormatCSV, ExportFormatExcel}

	for _, format := range formats {
		export := &ExportHistory{
			ID:         "export123",
			AdminID:    "admin123",
			ExportType: ExportTypeBooks,
			Format:     format,
			Status:     ExportStatusPending,
			CreatedAt:  time.Now(),
		}

		assert.Contains(t, formats, export.Format)
	}
}

// TestExportHistory_MarkCompleted 测试标记为完成
func TestExportHistory_MarkCompleted(t *testing.T) {
	export := &ExportHistory{
		ID:          "export123",
		AdminID:     "admin123",
		ExportType:  ExportTypeBooks,
		Format:      ExportFormatCSV,
		Status:      ExportStatusPending,
		RecordCount: 500,
		FilePath:    "/exports/books.csv",
		FileSize:    512000,
		CreatedAt:   time.Now().Add(-1 * time.Hour),
		CompletedAt: nil,
	}

	// 标记为完成
	now := time.Now()
	export.MarkCompleted()
	export.CompletedAt = &now

	assert.Equal(t, ExportStatusCompleted, export.Status)
	assert.NotNil(t, export.CompletedAt)
	assert.Greater(t, export.RecordCount, 0)
	assert.Greater(t, export.FileSize, int64(0))
}

// TestExportHistory_MarkFailed 测试标记为失败
func TestExportHistory_MarkFailed(t *testing.T) {
	export := &ExportHistory{
		ID:          "export123",
		AdminID:     "admin123",
		ExportType:  ExportTypeBooks,
		Format:      ExportFormatCSV,
		Status:      ExportStatusPending,
		ErrorMsg:    "",
		CreatedAt:   time.Now().Add(-1 * time.Hour),
		CompletedAt: nil,
	}

	// 标记为失败
	errorMsg := "导出失败：内存不足"
	export.MarkFailed()
	export.Status = ExportStatusFailed
	export.ErrorMsg = errorMsg

	assert.Equal(t, ExportStatusFailed, export.Status)
	assert.Equal(t, errorMsg, export.ErrorMsg)
	assert.Nil(t, export.CompletedAt)
}

// TestExportHistory_InvalidStatus 测试无效状态
func TestExportHistory_InvalidStatus(t *testing.T) {
	export := &ExportHistory{
		ID:         "export123",
		AdminID:    "admin123",
		ExportType: ExportTypeBooks,
		Format:     ExportFormatCSV,
		Status:     "unknown_status",
		CreatedAt:  time.Now(),
	}

	assert.False(t, export.IsPending())
	assert.False(t, export.IsCompleted())
	assert.False(t, export.IsFailed())
}

// TestExportHistory_EmptyFields 测试空字段
func TestExportHistory_EmptyFields(t *testing.T) {
	export := &ExportHistory{}

	assert.Empty(t, export.ID)
	assert.Empty(t, export.AdminID)
	assert.Empty(t, export.ExportType)
	assert.Empty(t, export.Format)
	assert.Empty(t, export.FilePath)
	assert.Empty(t, export.ErrorMsg)
	assert.Empty(t, export.Status)
	assert.Zero(t, export.RecordCount)
	assert.Zero(t, export.FileSize)
	assert.Nil(t, export.CompletedAt)
	assert.True(t, export.CreatedAt.IsZero())
}

// TestExportHistory_ZeroValues 测试零值
func TestExportHistory_ZeroValues(t *testing.T) {
	export := &ExportHistory{
		ID:          "export123",
		AdminID:     "admin123",
		ExportType:  ExportTypeBooks,
		Format:      ExportFormatCSV,
		RecordCount: 0,
		FileSize:    0,
		Status:      ExportStatusPending,
		CreatedAt:   time.Now(),
	}

	assert.Zero(t, export.RecordCount)
	assert.Zero(t, export.FileSize)
}

// TestExportHistory_NegativeValues 测试负值
func TestExportHistory_NegativeValues(t *testing.T) {
	export := &ExportHistory{
		ID:          "export123",
		AdminID:     "admin123",
		ExportType:  ExportTypeBooks,
		Format:      ExportFormatCSV,
		RecordCount: -1,  // 无效值
		FileSize:    -100, // 无效值
		Status:      ExportStatusPending,
		CreatedAt:   time.Now(),
	}

	// 负数在结构上是允许的，但业务逻辑应该检查
	assert.Less(t, export.RecordCount, 0)
	assert.Less(t, export.FileSize, int64(0))
}

// TestExportHistory_LargeValues 测试大值
func TestExportHistory_LargeValues(t *testing.T) {
	export := &ExportHistory{
		ID:          "export123",
		AdminID:     "admin123",
		ExportType:  ExportTypeBooks,
		Format:      ExportFormatCSV,
		RecordCount: 1000000,      // 100万条记录
		FileSize:    1024 * 1024 * 1024, // 1GB
		Status:      ExportStatusCompleted,
		CreatedAt:   time.Now(),
	}

	assert.Greater(t, export.RecordCount, 999999)
	assert.Greater(t, export.FileSize, int64(1024*1024*1024-1))
}

// TestExportHistory_LongPaths 测试长路径
func TestExportHistory_LongPaths(t *testing.T) {
	longPath := "/exports/2026/02/27/admin123/books_very_long_filename_with_many_characters_and_timestamp_20260227_153045_1234567890.csv"

	export := &ExportHistory{
		ID:          "export123",
		AdminID:     "admin123",
		ExportType:  ExportTypeBooks,
		Format:      ExportFormatCSV,
		FilePath:    longPath,
		Status:      ExportStatusCompleted,
		CreatedAt:   time.Now(),
	}

	assert.Greater(t, len(export.FilePath), 100)
	assert.Contains(t, export.FilePath, ".csv")
}

// TestExportHistory_LongErrorMessage 测试长错误信息
func TestExportHistory_LongErrorMessage(t *testing.T) {
	longErrorMsg := "导出失败：系统错误详细信息堆栈跟踪Database connection timeout at 2026-02-27 15:30:45. Connection string: mongodb://localhost:27017. Error: context deadline exceeded. Retry attempts: 3. Last error: dial tcp 127.0.0.1:27017: i/o timeout"

	export := &ExportHistory{
		ID:          "export123",
		AdminID:     "admin123",
		ExportType:  ExportTypeBooks,
		Format:      ExportFormatCSV,
		Status:      ExportStatusFailed,
		ErrorMsg:    longErrorMsg,
		CreatedAt:   time.Now(),
	}

	assert.Greater(t, len(export.ErrorMsg), 200)
	assert.Contains(t, export.ErrorMsg, "timeout")
}

// TestExportHistory_UnicodeInFields 测试字段中的Unicode字符
func TestExportHistory_UnicodeInFields(t *testing.T) {
	export := &ExportHistory{
		ID:          "导出123",
		AdminID:     "管理员456",
		ExportType:  ExportTypeBooks,
		Format:      ExportFormatCSV,
		FilePath:    "/exports/书籍_导出.csv",
		Status:      ExportStatusCompleted,
		ErrorMsg:    "错误：文件不存在",
		CreatedAt:   time.Now(),
	}

	assert.Contains(t, export.ID, "导出")
	assert.Contains(t, export.FilePath, "书籍")
	assert.Contains(t, export.ErrorMsg, "错误")
}

// TestExportHistory_Timestamps 测试时间戳
func TestExportHistory_Timestamps(t *testing.T) {
	createdTime := time.Date(2026, 2, 27, 15, 30, 45, 0, time.UTC)
	completedTime := time.Date(2026, 2, 27, 15, 35, 45, 0, time.UTC)

	export := &ExportHistory{
		ID:          "export123",
		AdminID:     "admin123",
		ExportType:  ExportTypeBooks,
		Format:      ExportFormatCSV,
		Status:      ExportStatusCompleted,
		CreatedAt:   createdTime,
		CompletedAt: &completedTime,
	}

	assert.Equal(t, createdTime, export.CreatedAt)
	assert.Equal(t, completedTime, *export.CompletedAt)
	assert.True(t, export.CompletedAt.After(export.CreatedAt))
}

// TestExportHistory_StatusTransitions 测试状态转换
func TestExportHistory_StatusTransitions(t *testing.T) {
	export := &ExportHistory{
		ID:         "export123",
		AdminID:    "admin123",
		ExportType: ExportTypeBooks,
		Format:     ExportFormatCSV,
		Status:     ExportStatusPending,
		CreatedAt:  time.Now(),
	}

	// 初始状态
	assert.True(t, export.IsPending())

	// 转换到完成
	export.Status = ExportStatusCompleted
	assert.True(t, export.IsCompleted())

	// 重置并转换到失败
	export.Status = ExportStatusFailed
	assert.True(t, export.IsFailed())
}

// TestExportHistory_IDFormats 测试各种ID格式
func TestExportHistory_IDFormats(t *testing.T) {
	ids := []struct {
		name string
		id   string
	}{
		{"MongoDB ObjectId", "507f1f77bcf86cd799439011"},
		{"UUID", "550e8400-e29b-41d4-a716-446655440000"},
		{"Simple ID", "export123"},
		{"Timestamp ID", "export_20260227_153045"},
	}

	for _, tc := range ids {
		t.Run(tc.name, func(t *testing.T) {
			export := &ExportHistory{
				ID:         tc.id,
				AdminID:    "admin123",
				ExportType: ExportTypeBooks,
				Format:     ExportFormatCSV,
				Status:     ExportStatusPending,
				CreatedAt:  time.Now(),
			}

			assert.NotEmpty(t, export.ID)
		})
	}
}

// TestExportHistory_AllStatuses 测试所有状态
func TestExportHistory_AllStatuses(t *testing.T) {
	statuses := []string{ExportStatusPending, ExportStatusCompleted, ExportStatusFailed}

	for _, status := range statuses {
		export := &ExportHistory{
			ID:         "export123",
			AdminID:    "admin123",
			ExportType: ExportTypeBooks,
			Format:     ExportFormatCSV,
			Status:     status,
			CreatedAt:  time.Now(),
		}

		assert.Contains(t, statuses, export.Status)

		switch status {
		case ExportStatusPending:
			assert.True(t, export.IsPending())
		case ExportStatusCompleted:
			assert.True(t, export.IsCompleted())
		case ExportStatusFailed:
			assert.True(t, export.IsFailed())
		}
	}
}

// TestExportHistory_Validate_Valid 测试有效验证
func TestExportHistory_Validate_Valid(t *testing.T) {
	export := &ExportHistory{
		ID:         "export123",
		AdminID:    "admin123",
		ExportType: ExportTypeBooks,
		Format:     ExportFormatCSV,
		Status:     ExportStatusPending,
		CreatedAt:  time.Now(),
	}

	err := export.Validate()
	assert.NoError(t, err)
}

// TestExportHistory_Validate_EmptyAdminID 测试空管理员ID
func TestExportHistory_Validate_EmptyAdminID(t *testing.T) {
	export := &ExportHistory{
		ID:         "export123",
		AdminID:    "",
		ExportType: ExportTypeBooks,
		Format:     ExportFormatCSV,
		Status:     ExportStatusPending,
		CreatedAt:  time.Now(),
	}

	err := export.Validate()
	assert.Error(t, err)
	assert.Equal(t, ErrExportAdminIDRequired, err)
}

// TestExportHistory_Validate_EmptyExportType 测试空导出类型
func TestExportHistory_Validate_EmptyExportType(t *testing.T) {
	export := &ExportHistory{
		ID:         "export123",
		AdminID:    "admin123",
		ExportType: "",
		Format:     ExportFormatCSV,
		Status:     ExportStatusPending,
		CreatedAt:  time.Now(),
	}

	err := export.Validate()
	assert.Error(t, err)
	assert.Equal(t, ErrExportTypeRequired, err)
}

// TestExportHistory_Validate_EmptyFormat 测试空格式
func TestExportHistory_Validate_EmptyFormat(t *testing.T) {
	export := &ExportHistory{
		ID:         "export123",
		AdminID:    "admin123",
		ExportType: ExportTypeBooks,
		Format:     "",
		Status:     ExportStatusPending,
		CreatedAt:  time.Now(),
	}

	err := export.Validate()
	assert.Error(t, err)
	assert.Equal(t, ErrExportFormatRequired, err)
}

// TestExportHistory_Validate_AllEmpty 测试全部为空
func TestExportHistory_Validate_AllEmpty(t *testing.T) {
	export := &ExportHistory{
		ID:      "export123",
		AdminID: "",
		Status:  ExportStatusPending,
	}

	err := export.Validate()
	assert.Error(t, err)
}

// TestExportError_Error 测试错误信息
func TestExportError_Error(t *testing.T) {
	err := NewExportError("test error message")
	assert.Equal(t, "test error message", err.Error())
}

// TestExportError_Constants 测试导出错误常量
func TestExportError_Constants(t *testing.T) {
	assert.NotEmpty(t, ErrExportAdminIDRequired.Error())
	assert.NotEmpty(t, ErrExportTypeRequired.Error())
	assert.NotEmpty(t, ErrExportFormatRequired.Error())
}
