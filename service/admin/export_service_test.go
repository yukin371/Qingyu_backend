package admin

import (
	"bytes"
	"context"
	"encoding/csv"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

// ==================== 测试数据结构 ====================

// MockExportable 模拟可导出数据
type MockExportable struct {
	ID    string
	Name  string
	Email string
	Age   int
}

func (m *MockExportable) ToExportRow() []string {
	return []string{m.ID, m.Name, m.Email, strconv.Itoa(m.Age)}
}

// ==================== TestExportService_ExportUsersToCSV_Success ====================

func TestExportService_ExportUsersToCSV_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	service := NewExportService()

	data := []Exportable{
		&MockExportable{ID: "1", Name: "Alice", Email: "alice@example.com", Age: 25},
		&MockExportable{ID: "2", Name: "Bob", Email: "bob@example.com", Age: 30},
	}

	columns := []ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "Name", Title: "姓名"},
		{Key: "Email", Title: "邮箱"},
		{Key: "Age", Title: "年龄"},
	}

	// When
	result, err := service.ExportToCSV(ctx, data, columns)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// 验证CSV内容
	reader := csv.NewReader(bytes.NewReader(result))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(records)) // Header + 2 rows

	// 验证表头
	assert.Equal(t, "ID", records[0][0])
	assert.Equal(t, "姓名", records[0][1])
	assert.Equal(t, "邮箱", records[0][2])
	assert.Equal(t, "年龄", records[0][3])

	// 验证数据
	assert.Equal(t, "1", records[1][0])
	assert.Equal(t, "Alice", records[1][1])
	assert.Equal(t, "alice@example.com", records[1][2])
}

// ==================== TestExportService_ExportUsersToExcel_Success ====================

func TestExportService_ExportUsersToExcel_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	service := NewExportService()

	data := []Exportable{
		&MockExportable{ID: "1", Name: "Alice", Email: "alice@example.com", Age: 25},
		&MockExportable{ID: "2", Name: "Bob", Email: "bob@example.com", Age: 30},
	}

	columns := []ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "Name", Title: "姓名"},
		{Key: "Email", Title: "邮箱"},
		{Key: "Age", Title: "年龄"},
	}

	// When
	result, err := service.ExportToExcel(ctx, data, columns)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Greater(t, len(result), 0)

	// 验证Excel文件可以打开
	_, err = excelize.OpenReader(bytes.NewReader(result))
	assert.NoError(t, err)
}

// ==================== TestExportService_ExportWithCustomColumns_Success ====================

func TestExportService_ExportWithCustomColumns_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	service := NewExportService()

	data := []Exportable{
		&MockExportable{ID: "1", Name: "Alice", Email: "alice@example.com", Age: 25},
	}

	// 自定义列 - 只导出部分字段
	columns := []ExportColumn{
		{Key: "Name", Title: "姓名"},
		{Key: "Email", Title: "邮箱"},
	}

	// When
	csvResult, err := service.ExportToCSV(ctx, data, columns)

	// Then
	assert.NoError(t, err)

	reader := csv.NewReader(bytes.NewReader(csvResult))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(records)) // Header + 1 row
	assert.Equal(t, 2, len(records[0])) // 只有两列
	assert.Equal(t, "姓名", records[0][0])
	assert.Equal(t, "邮箱", records[0][1])
}

// ==================== TestExportService_HandleLargeDataset_Success ====================

func TestExportService_HandleLargeDataset_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	service := NewExportService()

	// 创建1000条数据
	data := make([]Exportable, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = &MockExportable{
			ID:    string(rune(i)),
			Name:  "User",
			Email: "user@example.com",
			Age:   25,
		}
	}

	columns := []ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "Name", Title: "姓名"},
	}

	// When
	result, err := service.ExportToCSV(ctx, data, columns)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, result)

	reader := csv.NewReader(bytes.NewReader(result))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 1001, len(records)) // Header + 1000 rows
}

// ==================== TestExportService_ExportEmptyData_Success ====================

func TestExportService_ExportEmptyData_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	service := NewExportService()

	data := []Exportable{}
	columns := []ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "Name", Title: "姓名"},
	}

	// When
	result, err := service.ExportToCSV(ctx, data, columns)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, result)

	reader := csv.NewReader(bytes.NewReader(result))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(records)) // Only header
	assert.Equal(t, "ID", records[0][0])
	assert.Equal(t, "姓名", records[0][1])
}

// ==================== TestExportService_ExportWithSpecialCharacters_Success ====================

func TestExportService_ExportWithSpecialCharacters_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	service := NewExportService()

	data := []Exportable{
		&MockExportable{ID: "1", Name: "Alice, Bob", Email: "test@example.com", Age: 25},
		&MockExportable{ID: "2", Name: "Charlie \"Chuck\"", Email: "charlie@example.com", Age: 30},
	}

	columns := []ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "Name", Title: "姓名"},
		{Key: "Email", Title: "邮箱"},
	}

	// When
	result, err := service.ExportToCSV(ctx, data, columns)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// 验证特殊字符被正确处理
	reader := csv.NewReader(bytes.NewReader(result))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(records))
	// CSV会自动处理特殊字符，带逗号的字段会被引号包围
	assert.Equal(t, "Alice, Bob", records[1][1])
	assert.Equal(t, "Charlie \"Chuck\"", records[2][1])
}

// ==================== TestExportService_GetExportTemplate_Success ====================

func TestExportService_GetExportTemplate_Success(t *testing.T) {
	// Given
	service := NewExportService()

	columns := []ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "Name", Title: "姓名"},
		{Key: "Email", Title: "邮箱"},
	}

	// When
	result, err := service.GetExportTemplate(columns)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// 验证CSV包含表头
	reader := csv.NewReader(bytes.NewReader(result))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(records)) // Only header
	assert.Equal(t, "ID", records[0][0])
	assert.Equal(t, "姓名", records[0][1])
	assert.Equal(t, "邮箱", records[0][2])
}

// ==================== TestExportService_ExportToExcelWithMultipleSheets ====================

func TestExportService_ExportToExcelWithMultipleSheets(t *testing.T) {
	// Given
	ctx := context.Background()
	service := NewExportService()

	data := []Exportable{
		&MockExportable{ID: "1", Name: "Alice", Email: "alice@example.com", Age: 25},
	}

	columns := []ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "Name", Title: "姓名"},
	}

	// When
	result, err := service.ExportToExcel(ctx, data, columns)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// 验证Excel文件结构和内容
	f, err := excelize.OpenReader(bytes.NewReader(result))
	assert.NoError(t, err)

	// 获取所有工作表
	sheets := f.GetSheetList()
	assert.Greater(t, len(sheets), 0)

	// 读取第一个工作表的数据
	rows, err := f.GetRows("Sheet1")
	assert.NoError(t, err)
	assert.Greater(t, len(rows), 0)
}

// ==================== TestExportService_ExportNilColumns_Error ====================

func TestExportService_ExportNilColumns_Error(t *testing.T) {
	// Given
	ctx := context.Background()
	service := NewExportService()

	data := []Exportable{
		&MockExportable{ID: "1", Name: "Alice", Email: "alice@example.com", Age: 25},
	}

	// When - 空列定义
	_, err := service.ExportToCSV(ctx, data, nil)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "列定义")
}

// ==================== TestExportService_ExportWithChineseCharacters_Success ====================

func TestExportService_ExportWithChineseCharacters_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	service := NewExportService()

	data := []Exportable{
		&MockExportable{ID: "1", Name: "张三", Email: "zhangsan@example.com", Age: 25},
		&MockExportable{ID: "2", Name: "李四", Email: "lisi@example.com", Age: 30},
	}

	columns := []ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "Name", Title: "姓名"},
		{Key: "Email", Title: "邮箱"},
	}

	// When
	result, err := service.ExportToCSV(ctx, data, columns)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// 验证中文字符被正确处理
	reader := csv.NewReader(bytes.NewReader(result))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, "1", records[1][0])
	assert.Equal(t, "张三", records[1][1])
	assert.Equal(t, "zhangsan@example.com", records[1][2])
	assert.Equal(t, "2", records[2][0])
	assert.Equal(t, "李四", records[2][1])
	assert.Equal(t, "lisi@example.com", records[2][2])
}

// ==================== TestExportService_ExportDataWithInconsistentColumns ====================

func TestExportService_ExportDataWithInconsistentColumns_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	service := NewExportService()

	// 数据行数与列数不匹配
	data := []Exportable{
		&MockExportable{ID: "1", Name: "Alice", Email: "alice@example.com", Age: 25},
		&MockExportable{ID: "2", Name: "Bob", Email: "bob@example.com", Age: 30},
	}

	columns := []ExportColumn{
		{Key: "ID", Title: "ID"},
		{Key: "Name", Title: "姓名"},
	}

	// When
	result, err := service.ExportToCSV(ctx, data, columns)

	// Then - 应该成功，只导出指定列
	assert.NoError(t, err)
	assert.NotNil(t, result)

	reader := csv.NewReader(bytes.NewReader(result))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(records)) // Header + 2 rows
	assert.Equal(t, 2, len(records[0])) // Only 2 columns
}
