package writer

import (
	"Qingyu_backend/models/writer"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectDetailResponse 项目详情响应
type ProjectDetailResponse struct {
	ID             primitive.ObjectID       `json:"id"`
	Title          string                   `json:"title"`
	Description    string                   `json:"description"`
	CoverImage     string                   `json:"coverImage"`
	Genre          string                   `json:"genre"`
	Tags           []string                 `json:"tags"`
	Status         string                   `json:"status"`
	Visibility     string                   `json:"visibility"`
	TotalWords     int64                    `json:"totalWords"`
	ChapterCount   int                      `json:"chapterCount"`
	LastUpdateTime time.Time                `json:"lastUpdateTime"`
	CreatedAt      time.Time                `json:"createdAt"`
	UpdatedAt      time.Time                `json:"updatedAt"`
	Documents      []DocumentSummary      `json:"documents"`
	Characters     []writer.Character     `json:"characters"`
	Locations      []writer.Location      `json:"locations"`
	Timeline       []writer.TimelineEvent `json:"timeline"`
}

// DocumentSummary 文档摘要
type DocumentSummary struct {
	ID         primitive.ObjectID `json:"id"`
	Title      string             `json:"title"`
	Type       string             `json:"type"`
	WordCount  int                `json:"wordCount"`
	LastEditAt time.Time          `json:"lastEditAt"`
	Status     string             `json:"status"`
	SortOrder  int                `json:"sortOrder"`
}

// ProjectListResponse 项目列表响应
type ProjectListResponse struct {
	Projects []ProjectSummary `json:"projects"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Size     int              `json:"size"`
}

// ProjectSummary 项目摘要
type ProjectSummary struct {
	ID             primitive.ObjectID `json:"id"`
	Title          string             `json:"title"`
	CoverImage     string             `json:"coverImage"`
	Genre          string             `json:"genre"`
	Status         string             `json:"status"`
	TotalWords     int64              `json:"totalWords"`
	ChapterCount   int                `json:"chapterCount"`
	LastUpdateTime time.Time          `json:"lastUpdateTime"`
}

// DocumentDetailResponse 文档详情响应
type DocumentDetailResponse struct {
	ID         primitive.ObjectID `json:"id"`
	ProjectID  primitive.ObjectID `json:"projectId"`
	Title      string             `json:"title"`
	Content    string             `json:"content"`
	Type       string             `json:"type"`
	Status     string             `json:"status"`
	WordCount  int                `json:"wordCount"`
	Version    int                `json:"version"`
	LastEditAt time.Time          `json:"lastEditAt"`
	CreatedAt  time.Time          `json:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt"`
}

// VersionHistoryResponse 版本历史响应
type VersionHistoryResponse struct {
	Versions []VersionInfo `json:"versions"`
	Total    int64         `json:"total"`
}

// VersionInfo 版本信息
type VersionInfo struct {
	Version    int       `json:"version"`
	Comment    string    `json:"comment"`
	WordCount  int       `json:"wordCount"`
	CreatedAt  time.Time `json:"createdAt"`
	CreatedBy  string    `json:"createdBy"`
	IsSnapshot bool      `json:"isSnapshot"`
}

// VersionDiffResponse 版本对比响应
type VersionDiffResponse struct {
	OldVersion int      `json:"oldVersion"`
	NewVersion int      `json:"newVersion"`
	Additions  []string `json:"additions"`
	Deletions  []string `json:"deletions"`
	Changes    []string `json:"changes"`
}

// StatsResponse 统计数据响应
type StatsResponse struct {
	TotalWords     int64              `json:"totalWords"`
	ChapterCount   int                `json:"chapterCount"`
	AvgChapterLen  int                `json:"avgChapterLen"`
	ReadCount      int64              `json:"readCount"`
	LikeCount      int64              `json:"likeCount"`
	CommentCount   int64              `json:"commentCount"`
	DailyStats     []DailyStat        `json:"dailyStats"`
	ChapterStats   []ChapterStat      `json:"chapterStats"`
	ReaderBehavior ReaderBehaviorStat `json:"readerBehavior"`
}

// DailyStat 每日统计
type DailyStat struct {
	Date      string `json:"date"`
	WordCount int    `json:"wordCount"`
	ReadCount int64  `json:"readCount"`
}

// ChapterStat 章节统计
type ChapterStat struct {
	ChapterID   primitive.ObjectID `json:"chapterId"`
	ChapterName string             `json:"chapterName"`
	ReadCount   int64              `json:"readCount"`
	AvgReadTime int                `json:"avgReadTime"` // 秒
	FinishRate  float64            `json:"finishRate"`  // 完读率
}

// ReaderBehaviorStat 读者行为统计
type ReaderBehaviorStat struct {
	AvgReadTime      int      `json:"avgReadTime"`      // 平均阅读时长（秒）
	RetentionRate    float64  `json:"retentionRate"`    // 留存率
	SkipRate         float64  `json:"skipRate"`         // 跳章率
	FavoriteChapters []string `json:"favoriteChapters"` // 最受欢迎的章节
}

// AuditResultResponse 审核结果响应
type AuditResultResponse struct {
	DocumentID primitive.ObjectID `json:"documentId"`
	Status     string             `json:"status"`    // pending, approved, rejected
	RiskLevel  string             `json:"riskLevel"` // low, medium, high
	Issues     []AuditIssue       `json:"issues"`
	ReviewedAt time.Time          `json:"reviewedAt"`
	ReviewedBy string             `json:"reviewedBy"`
	Comment    string             `json:"comment"`
}

// AuditIssue 审核问题
type AuditIssue struct {
	Type        string `json:"type"` // sensitive_word, violence, etc.
	Description string `json:"description"`
	Position    int    `json:"position"` // 在文本中的位置
	Severity    string `json:"severity"` // low, medium, high
	Suggestion  string `json:"suggestion"`
}

// CheckContentRequest 检测内容请求
type CheckContentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=100000"`
}

// AuditDocumentRequest 审核文档请求
type AuditDocumentRequest struct {
	DocumentID string `json:"documentId" validate:"required"`
	Content    string `json:"content" validate:"required"`
}

// UpdateAuditResultRequest 更新审核结果请求
type UpdateAuditResultRequest struct {
	Status  string `json:"status" validate:"required,oneof=approved rejected"`
	Comment string `json:"comment,omitempty"`
}

// AddSensitiveWordRequest 添加敏感词请求
type AddSensitiveWordRequest struct {
	Word     string   `json:"word" validate:"required"`
	Category string   `json:"category" validate:"required"`
	Level    string   `json:"level" validate:"required,oneof=low medium high"`
	Tags     []string `json:"tags,omitempty"`
}

// SubmitAppealRequest 申诉请求
type SubmitAppealRequest struct {
	Reason string `json:"reason" validate:"required,min=10,max=500"`
}

// ReviewAuditRequest 复核请求
type ReviewAuditRequest struct {
	Approved bool   `json:"approved"`
	Note     string `json:"note" validate:"max=500"`
}

// ReviewAppealRequest 复核申诉请求
type ReviewAppealRequest struct {
	Approved bool   `json:"approved"`
	Note     string `json:"note" validate:"max=500"`
}
