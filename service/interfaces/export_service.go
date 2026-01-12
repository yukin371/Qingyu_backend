package interfaces

import (
	"context"
	"time"
)

// ExportService 导出服务接口
type ExportService interface {
	// 文档导出
	ExportDocument(ctx context.Context, documentID, projectID, userID string, req *ExportDocumentRequest) (*ExportTask, error)
	GetExportTask(ctx context.Context, taskID string) (*ExportTask, error)
	DownloadExportFile(ctx context.Context, taskID string) (*ExportFile, error)
	ListExportTasks(ctx context.Context, projectID string, page, pageSize int) ([]*ExportTask, int64, error)
	DeleteExportTask(ctx context.Context, taskID, userID string) error

	// 项目导出
	ExportProject(ctx context.Context, projectID, userID string, req *ExportProjectRequest) (*ExportTask, error)

	// 任务管理
	CancelExportTask(ctx context.Context, taskID, userID string) error
}

// ExportDocumentRequest 导出文档请求
type ExportDocumentRequest struct {
	Format      string         `json:"format" validate:"required,oneof=txt md docx"` // 导出格式
	IncludeMeta bool           `json:"includeMeta"`                                  // 是否包含元数据
	Options     *ExportOptions `json:"options,omitempty"`                            // 导出选项
}

// ExportProjectRequest 导出项目请求
type ExportProjectRequest struct {
	IncludeDocuments  bool           `json:"includeDocuments"`                             // 是否包含文档
	IncludeCharacters bool           `json:"includeCharacters"`                            // 是否包含角色
	IncludeLocations  bool           `json:"includeLocations"`                             // 是否包含地点
	IncludeTimeline   bool           `json:"includeTimeline"`                              // 是否包含时间线
	DocumentFormats   string         `json:"documentFormats" validate:"oneof=txt md docx"` // 文档导出格式
	Options           *ExportOptions `json:"options,omitempty"`
}

// ExportOptions 导出选项
type ExportOptions struct {
	TOC             bool     `json:"toc"`                       // 是否生成目录
	PageNumbers     bool     `json:"pageNumbers"`               // 是否添加页码
	IncludeNotes    bool     `json:"includeNotes"`              // 是否包含注释
	IncludeTags     bool     `json:"includeTags"`               // 是否包含标签
	Header          string   `json:"header,omitempty"`          // 页眉
	Footer          string   `json:"footer,omitempty"`          // 页脚
	FontSize        int      `json:"fontSize,omitempty"`        // 字号
	LineSpacing     float64  `json:"lineSpacing,omitempty"`     // 行间距
	ExcludeChapters []string `json:"excludeChapters,omitempty"` // 排除的章节ID
}

// ExportTask 导出任务
type ExportTask struct {
	ID            string     `json:"id"`
	Type          string     `json:"type"`       // document, project
	ResourceID    string     `json:"resourceId"` // 文档ID或项目ID
	ResourceTitle string     `json:"resourceTitle"`
	Format        string     `json:"format"`
	Status        string     `json:"status"`            // pending, processing, completed, failed, cancelled
	Progress      int        `json:"progress"`          // 0-100
	FileSize      int64      `json:"fileSize"`          // 文件大小（字节）
	FileURL       string     `json:"fileUrl,omitempty"` // 下载链接
	ExpiresAt     time.Time  `json:"expiresAt"`         // 文件过期时间
	ErrorMsg      string     `json:"errorMsg,omitempty"`
	CreatedBy     string     `json:"createdBy"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
}

// ExportFile 导出文件
type ExportFile struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content,omitempty"`
	URL      string `json:"url,omitempty"`
	MimeType string `json:"mimeType"`
	FileSize int64  `json:"fileSize"`
}

// ExportTaskStatus 导出任务状态常量
const (
	ExportStatusPending    = "pending"    // 等待处理
	ExportStatusProcessing = "processing" // 处理中
	ExportStatusCompleted  = "completed"  // 已完成
	ExportStatusFailed     = "failed"     // 失败
	ExportStatusCancelled  = "cancelled"  // 已取消
)

// ExportFormat 导出格式常量
const (
	ExportFormatTXT  = "txt"
	ExportFormatMD   = "md"
	ExportFormatDOCX = "docx"
	ExportFormatZIP  = "zip"
)

// ExportType 导出类型常量
const (
	ExportTypeDocument = "document"
	ExportTypeProject  = "project"
)
