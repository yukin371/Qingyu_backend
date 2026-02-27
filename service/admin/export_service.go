package admin

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// ==================== 导出格式定义 ====================

// ExportFormat 导出格式
type ExportFormat string

const (
	ExportFormatCSV   ExportFormat = "csv"
	ExportFormatExcel ExportFormat = "excel"
)

// ==================== 导出列定义 ====================

// ExportColumn 导出列定义
type ExportColumn struct {
	Key   string // 数据键
	Title string // 显示标题
}

// ==================== 可导出数据接口 ====================

// Exportable 可导出数据接口
type Exportable interface {
	ToExportRow() []string
}

// ==================== 导出服务接口 ====================

// ExportService 导出服务接口
type ExportService interface {
	// ExportToCSV 导出为CSV格式
	ExportToCSV(ctx context.Context, data []Exportable, columns []ExportColumn) ([]byte, error)
	// ExportToExcel 导出为Excel格式
	ExportToExcel(ctx context.Context, data []Exportable, columns []ExportColumn) ([]byte, error)
	// GetExportTemplate 获取导出模板（仅包含表头）
	GetExportTemplate(columns []ExportColumn) ([]byte, error)
}

// ==================== 导出服务实现 ====================

// exportServiceImpl 导出服务实现
type exportServiceImpl struct {
	// 默认Excel工作表名称
	defaultSheetName string
}

// NewExportService 创建导出服务
func NewExportService() ExportService {
	return &exportServiceImpl{
		defaultSheetName: "Sheet1",
	}
}

// ==================== CSV 导出实现 ====================

// ExportToCSV 导出为CSV格式
func (s *exportServiceImpl) ExportToCSV(ctx context.Context, data []Exportable, columns []ExportColumn) ([]byte, error) {
	// 验证参数
	if err := s.validateExportParams(data, columns); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 写入表头
	headers := make([]string, len(columns))
	for i, col := range columns {
		headers[i] = col.Title
	}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("写入表头失败: %w", err)
	}

	// 写入数据行
	for _, item := range data {
		row := item.ToExportRow()
		// 确保只导出指定列的数据
		if len(row) >= len(columns) {
			exportRow := row[:len(columns)]
			if err := writer.Write(exportRow); err != nil {
				return nil, fmt.Errorf("写入数据行失败: %w", err)
			}
		} else {
			// 如果数据列不足，填充空字符串
			exportRow := make([]string, len(columns))
			copy(exportRow, row)
			if err := writer.Write(exportRow); err != nil {
				return nil, fmt.Errorf("写入数据行失败: %w", err)
			}
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV生成失败: %w", err)
	}

	return buf.Bytes(), nil
}

// ==================== Excel 导出实现 ====================

// ExportToExcel 导出为Excel格式
func (s *exportServiceImpl) ExportToExcel(ctx context.Context, data []Exportable, columns []ExportColumn) ([]byte, error) {
	// 验证参数
	if err := s.validateExportParams(data, columns); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}

	// 创建Excel文件
	f := excelize.NewFile()
	sheetName := s.defaultSheetName

	// 写入表头
	for i, col := range columns {
		cell := s.columnToLetter(i) + "1"
		if err := f.SetCellValue(sheetName, cell, col.Title); err != nil {
			return nil, fmt.Errorf("写入表头失败: %w", err)
		}
	}

	// 写入数据行
	for rowIndex, item := range data {
		row := item.ToExportRow()
		excelRow := rowIndex + 2 // Excel行号从1开始，表头占1行

		for colIndex := range columns {
			cell := s.columnToLetter(colIndex) + strconv.Itoa(excelRow)
			value := ""
			if colIndex < len(row) {
				value = row[colIndex]
			}
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				return nil, fmt.Errorf("写入数据单元格失败: %w", err)
			}
		}
	}

	// 生成到缓冲区
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("Excel生成失败: %w", err)
	}

	return buf.Bytes(), nil
}

// ==================== 导出模板实现 ====================

// GetExportTemplate 获取导出模板（仅包含表头）
func (s *exportServiceImpl) GetExportTemplate(columns []ExportColumn) ([]byte, error) {
	if len(columns) == 0 {
		return nil, fmt.Errorf("列定义不能为空")
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 写入表头
	headers := make([]string, len(columns))
	for i, col := range columns {
		headers[i] = col.Title
	}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("写入表头失败: %w", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("生成模板失败: %w", err)
	}

	return buf.Bytes(), nil
}

// ==================== 辅助方法 ====================

// validateExportParams 验证导出参数
func (s *exportServiceImpl) validateExportParams(data []Exportable, columns []ExportColumn) error {
	if columns == nil || len(columns) == 0 {
		return fmt.Errorf("列定义不能为空")
	}
	return nil
}

// columnToLetter 将列索引转换为Excel列字母（0->A, 1->B, ..., 26->AA, ...）
func (s *exportServiceImpl) columnToLetter(index int) string {
	letter := ""
	for index >= 0 {
		remainder := index % 26
		letter = string(rune('A'+remainder)) + letter
		index = index/26 - 1
		if index < 0 {
			break
		}
	}
	return letter
}

// ==================== 导出数据适配器 ====================

// ExportDataAdapter 导出数据适配器，用于将任意结构转换为Exportable
type ExportDataAdapter struct {
	data []map[string]string
}

// NewExportDataAdapter 创建导出数据适配器
func NewExportDataAdapter(data []map[string]string) *ExportDataAdapter {
	return &ExportDataAdapter{data: data}
}

// ToExportRow 实现Exportable接口
func (a *ExportDataAdapter) ToExportRow(index int) []string {
	if index >= len(a.data) {
		return []string{}
	}
	row := make([]string, len(a.data[index]))
	i := 0
	for _, v := range a.data[index] {
		row[i] = v
		i++
	}
	return row
}

// ==================== 导出配置 ====================

// ExportConfig 导出配置
type ExportConfig struct {
	// 导出格式
	Format ExportFormat
	// 是否包含表头
	IncludeHeader bool
	// 文件名（不含扩展名）
	FileName string
	// 工作表名称（仅Excel）
	SheetName string
	// 最大导出行数限制
	MaxRows int
}

// DefaultExportConfig 默认导出配置
func DefaultExportConfig() *ExportConfig {
	return &ExportConfig{
		Format:         ExportFormatCSV,
		IncludeHeader:  true,
		FileName:       "export",
		SheetName:      "Sheet1",
		MaxRows:        100000,
	}
}
