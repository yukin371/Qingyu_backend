package dto

// ===========================
// 阅读 DTO（符合分层架构规范）
// ===========================

// ReadingProgressDTO 阅读进度数据传输对象
// 用于：Service 层和 API 层数据传输，ID 和时间字段使用字符串类型
type ReadingProgressDTO struct {
	ID        string `json:"id" validate:"required"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	// 关联信息
	UserID   string `json:"userId" validate:"required"`
	BookID   string `json:"bookId" validate:"required"`
	ChapterID string `json:"chapterId" validate:"required"`

	// 阅读进度
	Progress    int    `json:"progress" validate:"min=0,max=100"`        // 阅读进度 0-100
	ReadingTime int64  `json:"readingTime" validate:"min=0"`              // 阅读时长（秒）
	LastReadAt  string `json:"lastReadAt,omitempty"`                      // 最后阅读时间（ISO8601）

	// 状态
	Status string `json:"status" validate:"required,oneof=reading paused completed abandoned"`
}
