package admin

import (
	"time"
)

// ExportHistory 导出历史记录
type ExportHistory struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	AdminID     string    `bson:"admin_id" json:"admin_id"`         // 管理员ID
	ExportType  string    `bson:"export_type" json:"export_type"`   // 导出类型: books, chapters, users
	Format      string    `bson:"format" json:"format"`             // 导出格式: csv, excel
	RecordCount int       `bson:"record_count" json:"record_count"` // 记录数量
	FilePath    string    `bson:"file_path" json:"file_path"`       // 文件路径
	FileSize    int64     `bson:"file_size" json:"file_size"`       // 文件大小(字节)
	Status      string    `bson:"status" json:"status"`             // 状态: pending, completed, failed
	ErrorMsg    string    `bson:"error_msg,omitempty" json:"error_msg,omitempty"` // 错误信息
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`     // 创建时间
	CompletedAt *time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"` // 完成时间
}

// ExportHistoryStatus 导出状态
const (
	ExportStatusPending   = "pending"   // 待处理
	ExportStatusCompleted = "completed" // 已完成
	ExportStatusFailed    = "failed"    // 失败
)

// ExportType 导出类型
const (
	ExportTypeBooks    = "books"    // 书籍导出
	ExportTypeChapters = "chapters" // 章节导出
	ExportTypeUsers    = "users"    // 用户导出
)

// ExportFormat 导出格式
const (
	ExportFormatCSV   = "csv"   // CSV格式
	ExportFormatExcel = "excel" // Excel格式
)

// IsPending 检查是否为待处理状态
func (e *ExportHistory) IsPending() bool {
	return e.Status == ExportStatusPending
}

// IsCompleted 检查是否为已完成状态
func (e *ExportHistory) IsCompleted() bool {
	return e.Status == ExportStatusCompleted
}

// IsFailed 检查是否为失败状态
func (e *ExportHistory) IsFailed() bool {
	return e.Status == ExportStatusFailed
}

// MarkCompleted 标记为已完成
func (e *ExportHistory) MarkCompleted() {
	e.Status = ExportStatusCompleted
	now := time.Now()
	e.CompletedAt = &now
}

// MarkFailed 标记为失败
func (e *ExportHistory) MarkFailed() {
	e.Status = ExportStatusFailed
}

// Validate 验证导出历史记录
func (e *ExportHistory) Validate() error {
	if e.AdminID == "" {
		return ErrExportAdminIDRequired
	}
	if e.ExportType == "" {
		return ErrExportTypeRequired
	}
	if e.Format == "" {
		return ErrExportFormatRequired
	}
	return nil
}

// 导出历史错误
var (
	ErrExportAdminIDRequired = NewExportError("admin ID is required")
	ErrExportTypeRequired    = NewExportError("export type is required")
	ErrExportFormatRequired  = NewExportError("export format is required")
)

// ExportError 导出错误类型
type ExportError string

func NewExportError(msg string) ExportError {
	return ExportError(msg)
}

func (e ExportError) Error() string {
	return string(e)
}
