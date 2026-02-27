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
