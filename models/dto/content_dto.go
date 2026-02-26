package dto

import "time"

// ============================================================================
// 内容管理模块统一 DTO 定义
// ============================================================================

// 本文件包含内容管理模块（Writer模块）的所有DTO定义
// 用于替代 service/interfaces/writer/dto.go 和 service/writer/*/project_dto.go 中的重复定义
//
// 迁移指南：
// 1. 使用 dto.CreateProjectRequest 替代原 service/interfaces/writer.CreateProjectRequest
// 2. 使用 dto.CreateDocumentRequest 替代原 service/interfaces/writer.CreateDocumentRequest
// 3. 其他DTO类似迁移
//
// 注意：
// - 字段名统一使用 camelCase（与Go标准一致）
// - 时间字段统一使用 time.Time 类型
// - 添加完整的 validation 标签

// ============================================================================
// 项目管理 DTO
// ============================================================================

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Title    string   `json:"title" validate:"required,min=1,max=100"`
	Summary  string   `json:"summary,omitempty" validate:"max=500"`
	CoverURL string   `json:"coverUrl,omitempty" validate:"omitempty,url"`
	Category string   `json:"category,omitempty" validate:"max=50"`
	Tags     []string `json:"tags,omitempty" validate:"max=10"`
}

// CreateProjectResponse 创建项目响应
type CreateProjectResponse struct {
	ProjectID string    `json:"projectId"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// UpdateProjectRequest 更新项目请求
// 使用指针类型表示可选字段，允许部分更新
type UpdateProjectRequest struct {
	Title    *string   `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
	Summary  *string   `json:"summary,omitempty" validate:"omitempty,max=500"`
	CoverURL *string   `json:"coverUrl,omitempty" validate:"omitempty,url"`
	Category *string   `json:"category,omitempty" validate:"omitempty,max=50"`
	Tags     []string  `json:"tags,omitempty" validate:"max=10"`
	Status   *string   `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
}

// ListProjectsRequest 项目列表请求
type ListProjectsRequest struct {
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"pageSize" validate:"min=1,max=100"`
	Status   string `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	Category string `json:"category,omitempty"`
}

// ListProjectsResponse 项目列表响应
type ListProjectsResponse struct {
	Projects interface{} `json:"projects"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// GetProjectResponse 获取项目详情响应
// 使用匿名嵌入方式包含完整的项目信息
type GetProjectResponse struct {
	// 这里应该嵌入实际的项目模型
	// Project *writer.Project
	Title  string `json:"title"`
	Status string `json:"status"`
}

// UpdateStatisticsRequest 更新统计请求
type UpdateStatisticsRequest struct {
	ProjectID string `json:"projectId" validate:"required"`
}

// AddCollaboratorRequest 添加协作者请求
type AddCollaboratorRequest struct {
	ProjectID string `json:"projectId" validate:"required"`
	UserID    string `json:"userId" validate:"required"`
	Role      string `json:"role" validate:"required,oneof=owner editor viewer"`
}

// RemoveCollaboratorRequest 移除协作者请求
type RemoveCollaboratorRequest struct {
	ProjectID string `json:"projectId" validate:"required"`
	UserID    string `json:"userId" validate:"required"`
}

// ============================================================================
// 文档管理 DTO
// ============================================================================

// CreateDocumentRequest 创建文档请求
type CreateDocumentRequest struct {
	ProjectID    string   `json:"projectId" validate:"required"`
	ParentID     string   `json:"parentId,omitempty"`
	Title        string   `json:"title" validate:"required,min=1,max=200"`
	Type         string   `json:"type" validate:"required,oneof=chapter scene note"`
	Level        int      `json:"level,omitempty" validate:"min=0,max=10"`
	Order        int      `json:"order,omitempty" validate:"min=0"`
	CharacterIDs []string `json:"characterIds,omitempty" validate:"max=50"`
	LocationIDs  []string `json:"locationIds,omitempty" validate:"max=50"`
	TimelineIDs  []string `json:"timelineIds,omitempty" validate:"max=50"`
	Tags         []string `json:"tags,omitempty" validate:"max=20"`
	Notes        string   `json:"notes,omitempty" validate:"max=1000"`
}

// CreateDocumentResponse 创建文档响应
type CreateDocumentResponse struct {
	DocumentID string    `json:"documentId"`
	Title      string    `json:"title"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"createdAt"`
}

// UpdateDocumentRequest 更新文档请求
// 所有字段都是可选的，允许部分更新
type UpdateDocumentRequest struct {
	Title        *string   `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Notes        *string   `json:"notes,omitempty" validate:"omitempty,max=1000"`
	Status       *string   `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	CharacterIDs []string  `json:"characterIds,omitempty" validate:"max=50"`
	LocationIDs  []string  `json:"locationIds,omitempty" validate:"max=50"`
	TimelineIDs  []string  `json:"timelineIds,omitempty" validate:"max=50"`
	Tags         []string  `json:"tags,omitempty" validate:"max=20"`
}

// ListDocumentsRequest 文档列表请求
type ListDocumentsRequest struct {
	ProjectID string `json:"projectId" validate:"required"`
	ParentID  string `json:"parentId,omitempty"`
	Page      int    `json:"page" validate:"min=1"`
	PageSize  int    `json:"pageSize" validate:"min=1,max=100"`
	Status    string `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
}

// ListDocumentsResponse 文档列表响应
type ListDocumentsResponse struct {
	Documents interface{} `json:"documents"`
	Total     int         `json:"total"`
	Page      int         `json:"page"`
	PageSize  int         `json:"pageSize"`
}

// MoveDocumentRequest 移动文档请求
type MoveDocumentRequest struct {
	DocumentID  string `json:"documentId" validate:"required"`
	NewParentID string `json:"newParentId,omitempty"`
	Order       int    `json:"order" validate:"min=0"`
}

// ReorderDocumentsRequest 重新排序文档请求
type ReorderDocumentsRequest struct {
	ProjectID string         `json:"projectId" validate:"required"`
	ParentID  string         `json:"parentId,omitempty"`
	Orders    map[string]int `json:"orders" validate:"required,min=1"` // documentID -> order
}

// ============================================================================
// 文档内容 DTO
// ============================================================================

// AutoSaveRequest 自动保存请求
type AutoSaveRequest struct {
	DocumentID     string `json:"documentId" validate:"required"`
	Content        string `json:"content" validate:"required"`
	CurrentVersion int    `json:"currentVersion" validate:"min=0"` // 客户端当前版本号
	SaveType       string `json:"saveType" validate:"omitempty,oneof=auto manual"` // auto|manual
}

// AutoSaveResponse 自动保存响应
type AutoSaveResponse struct {
	Saved       bool      `json:"saved"`
	NewVersion  int       `json:"newVersion"`
	WordCount   int       `json:"wordCount"`
	SavedAt     time.Time `json:"savedAt"`
	HasConflict bool      `json:"hasConflict"`
}

// SaveStatusResponse 保存状态响应
type SaveStatusResponse struct {
	DocumentID     string    `json:"documentId"`
	LastSavedAt    time.Time `json:"lastSavedAt"`
	CurrentVersion int       `json:"currentVersion"`
	IsSaving       bool      `json:"isSaving"`
	WordCount      int       `json:"wordCount"`
}

// DocumentContentResponse 文档内容响应
type DocumentContentResponse struct {
	DocumentID string    `json:"documentId"`
	Content    string    `json:"content"`
	Version    int       `json:"version"`
	WordCount  int       `json:"wordCount"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// UpdateContentRequest 更新内容请求
type UpdateContentRequest struct {
	DocumentID string `json:"documentId" validate:"required"`
	Content    string `json:"content" validate:"required"`
	Version    int    `json:"version" validate:"min=0"` // 用于版本冲突检测
}

// DuplicateRequest 复制文档/项目请求
type DuplicateRequest struct {
	NewTitle      string `json:"newTitle" validate:"required,min=1,max=200"`
	TargetProject string `json:"targetProject,omitempty"` // 目标项目ID，空则表示当前项目
	CopyContent   bool   `json:"copyContent"`             // 是否复制内容
}

// DuplicateResponse 复制响应
type DuplicateResponse struct {
	NewDocumentID string `json:"newDocumentId"`
	Title         string `json:"title"`
}

// ============================================================================
// 版本控制 DTO
// ============================================================================

// VersionHistoryResponse 版本历史响应
type VersionHistoryResponse struct {
	Versions []*VersionInfo `json:"versions"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

// VersionInfo 版本信息
type VersionInfo struct {
	VersionID  string    `json:"versionId"`
	Version    int       `json:"version"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
	CreatedBy  string    `json:"createdBy"`
	WordCount  int       `json:"wordCount"`
}

// VersionDetail 版本详情
type VersionDetail struct {
	VersionID  string    `json:"versionId"`
	DocumentID string    `json:"documentId"`
	Version    int       `json:"version"`
	Content    string    `json:"content"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
	CreatedBy  string    `json:"createdBy"`
	WordCount  int       `json:"wordCount"`
}

// VersionDiff 版本差异
type VersionDiff struct {
	FromVersion  string      `json:"fromVersion"`
	ToVersion    string      `json:"toVersion"`
	Changes      []ChangeItem `json:"changes"`
	AddedLines   int         `json:"addedLines"`
	DeletedLines int         `json:"deletedLines"`
}

// ChangeItem 变更项
type ChangeItem struct {
	Type    string `json:"type"` // added, deleted, modified
	Line    int    `json:"line"`
	Content string `json:"content"`
}

// ============================================================================
// 大纲管理 DTO
// ============================================================================

// CreateOutlineRequest 创建大纲请求
type CreateOutlineRequest struct {
	Title      string   `json:"title" validate:"required,min=1,max=200"`
	ParentID   string   `json:"parentId,omitempty"`
	Summary    string   `json:"summary,omitempty" validate:"max=1000"`
	Type       string   `json:"type,omitempty" validate:"omitempty,oneof=arc chapter scene"`
	Tension    int      `json:"tension,omitempty" validate:"min=0,max=10"`
	ChapterID  string   `json:"chapterId,omitempty"`
	Characters []string `json:"characters,omitempty" validate:"max=50"`
	Items      []string `json:"items,omitempty" validate:"max=100"`
	Order      int      `json:"order,omitempty" validate:"min=0"`
}

// UpdateOutlineRequest 更新大纲请求
type UpdateOutlineRequest struct {
	Title      *string   `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	ParentID   *string   `json:"parentId,omitempty"`
	Summary    *string   `json:"summary,omitempty" validate:"omitempty,max=1000"`
	Type       *string   `json:"type,omitempty" validate:"omitempty,oneof=arc chapter scene"`
	Tension    *int      `json:"tension,omitempty" validate:"omitempty,min=0,max=10"`
	ChapterID  *string   `json:"chapterId,omitempty"`
	Characters *[]string `json:"characters,omitempty" validate:"omitempty,max=50"`
	Items      *[]string `json:"items,omitempty" validate:"omitempty,max=100"`
	Order      *int      `json:"order,omitempty" validate:"omitempty,min=0"`
}

// OutlineTreeNode 大纲树节点
type OutlineTreeNode struct {
	// 这里应该嵌入实际的大纲节点模型
	// *writer.OutlineNode
	Title     string             `json:"title"`
	Type      string             `json:"type"`
	Children  []*OutlineTreeNode `json:"children,omitempty"`
}

// ============================================================================
// 角色管理 DTO
// ============================================================================

// CreateCharacterRequest 创建角色请求
type CreateCharacterRequest struct {
	Name              string   `json:"name" validate:"required,min=1,max=100"`
	Alias             []string `json:"alias,omitempty" validate:"max=10"`
	Summary           string   `json:"summary,omitempty" validate:"max=500"`
	Traits            []string `json:"traits,omitempty" validate:"max=20"`
	Background        string   `json:"background,omitempty" validate:"max=2000"`
	AvatarURL         string   `json:"avatarUrl,omitempty" validate:"omitempty,url"`
	PersonalityPrompt string   `json:"personalityPrompt,omitempty" validate:"max=1000"`
	SpeechPattern     string   `json:"speechPattern,omitempty" validate:"max=500"`
	CurrentState      string   `json:"currentState,omitempty" validate:"max=500"`
}

// UpdateCharacterRequest 更新角色请求
type UpdateCharacterRequest struct {
	Name              *string   `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Alias             *[]string `json:"alias,omitempty" validate:"omitempty,max=10"`
	Summary           *string   `json:"summary,omitempty" validate:"omitempty,max=500"`
	Traits            *[]string `json:"traits,omitempty" validate:"omitempty,max=20"`
	Background        *string   `json:"background,omitempty" validate:"omitempty,max=2000"`
	AvatarURL         *string   `json:"avatarUrl,omitempty" validate:"omitempty,url"`
	PersonalityPrompt *string   `json:"personalityPrompt,omitempty" validate:"omitempty,max=1000"`
	SpeechPattern     *string   `json:"speechPattern,omitempty" validate:"omitempty,max=500"`
	CurrentState      *string   `json:"currentState,omitempty" validate:"omitempty,max=500"`
}

// CreateRelationRequest 创建关系请求
type CreateRelationRequest struct {
	FromID   string `json:"fromId" validate:"required"`
	ToID     string `json:"toId" validate:"required"`
	Type     string `json:"type" validate:"required,max=50"`
	Strength int    `json:"strength" validate:"min=0,max=10"`
	Notes    string `json:"notes,omitempty" validate:"max=500"`
}

// CharacterGraph 角色关系图
type CharacterGraph struct {
	Nodes interface{} `json:"nodes"`
	Edges interface{} `json:"edges"`
}

// ============================================================================
// 地点管理 DTO
// ============================================================================

// CreateLocationRequest 创建地点请求
type CreateLocationRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	Climate     string `json:"climate,omitempty" validate:"max=100"`
	Culture     string `json:"culture,omitempty" validate:"max=500"`
	Geography   string `json:"geography,omitempty" validate:"max=500"`
	Atmosphere  string `json:"atmosphere,omitempty" validate:"max=500"`
	ParentID    string `json:"parentId,omitempty"`
	ImageURL    string `json:"imageUrl,omitempty" validate:"omitempty,url"`
}

// UpdateLocationRequest 更新地点请求
type UpdateLocationRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Climate     *string `json:"climate,omitempty" validate:"omitempty,max=100"`
	Culture     *string `json:"culture,omitempty" validate:"omitempty,max=500"`
	Geography   *string `json:"geography,omitempty" validate:"omitempty,max=500"`
	Atmosphere  *string `json:"atmosphere,omitempty" validate:"omitempty,max=500"`
	ParentID    *string `json:"parentId,omitempty"`
	ImageURL    *string `json:"imageUrl,omitempty" validate:"omitempty,url"`
}

// LocationNode 地点节点
type LocationNode struct {
	Location interface{}     `json:"location"`
	Children []*LocationNode `json:"children"`
}

// CreateLocationRelationRequest 创建地点关系请求
type CreateLocationRelationRequest struct {
	FromID   string  `json:"fromId" validate:"required"`
	ToID     string  `json:"toId" validate:"required"`
	Type     string  `json:"type" validate:"required,max=50"`
	Distance *int    `json:"distance,omitempty" validate:"omitempty,min=0"`
	Notes    string  `json:"notes,omitempty" validate:"max=500"`
}

// ============================================================================
// 时间线管理 DTO
// ============================================================================

// CreateTimelineRequest 创建时间线请求
type CreateTimelineRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty" validate:"max=500"`
	StartTime   string `json:"startTime,omitempty"`
	EndTime     string `json:"endTime,omitempty"`
}

// CreateTimelineEventRequest 创建时间线事件请求
type CreateTimelineEventRequest struct {
	TimelineID   string   `json:"timelineId" validate:"required"`
	Title        string   `json:"title" validate:"required,min=1,max=200"`
	Description  string   `json:"description,omitempty" validate:"max=1000"`
	EventType    string   `json:"eventType,omitempty" validate:"max=50"`
	Importance   int      `json:"importance,omitempty" validate:"min=0,max=10"`
	Participants []string `json:"participants,omitempty" validate:"max=50"`
	LocationIDs  []string `json:"locationIds,omitempty" validate:"max=50"`
	ChapterIDs   []string `json:"chapterIds,omitempty" validate:"max=50"`
	StoryTime    string   `json:"storyTime,omitempty"`
	Duration     *int     `json:"duration,omitempty" validate:"omitempty,min=0"`
	Impact       string   `json:"impact,omitempty" validate:"max=500"`
}

// UpdateTimelineEventRequest 更新时间线事件请求
type UpdateTimelineEventRequest struct {
	Title        *string   `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Description  *string   `json:"description,omitempty" validate:"omitempty,max=1000"`
	EventType    *string   `json:"eventType,omitempty" validate:"omitempty,max=50"`
	Importance   *int      `json:"importance,omitempty" validate:"omitempty,min=0,max=10"`
	Participants *[]string `json:"participants,omitempty" validate:"omitempty,max=50"`
	LocationIDs  *[]string `json:"locationIds,omitempty" validate:"omitempty,max=50"`
	ChapterIDs   *[]string `json:"chapterIds,omitempty" validate:"omitempty,max=50"`
	StoryTime    *string   `json:"storyTime,omitempty"`
	Duration     *int      `json:"duration,omitempty" validate:"omitempty,min=0"`
	Impact       *string   `json:"impact,omitempty" validate:"omitempty,max=500"`
}

// TimelineVisualization 时间线可视化
type TimelineVisualization struct {
	Events      []*TimelineEventNode `json:"events"`
	Connections []*EventConnection   `json:"connections"`
}

// TimelineEventNode 时间线事件节点
type TimelineEventNode struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	StoryTime  string   `json:"storyTime"`
	EventType  string   `json:"eventType"`
	Importance int      `json:"importance"`
	Characters []string `json:"characters"`
}

// EventConnection 事件连接
type EventConnection struct {
	FromEventID string `json:"fromEventId"`
	ToEventID   string `json:"toEventId"`
	Type        string `json:"type"`
}

// ============================================================================
// 发布导出 DTO
// ============================================================================

// PublishProjectRequest 发布项目请求
type PublishProjectRequest struct {
	BookstoreID    string     `json:"bookstoreId" validate:"required"`
	CategoryID     string     `json:"categoryId" validate:"required"`
	Tags           []string   `json:"tags" validate:"max=10"`
	Description    string     `json:"description" validate:"max=1000"`
	CoverImage     string     `json:"coverImage" validate:"omitempty,url"`
	PublishType    string     `json:"publishType" validate:"required,oneof=serial complete"`
	PublishTime    *time.Time `json:"publishTime,omitempty"`
	Price          *float64   `json:"price,omitempty" validate:"omitempty,min=0"`
	FreeChapters   int        `json:"freeChapters" validate:"min=0"`
	AuthorNote     string     `json:"authorNote" validate:"max=500"`
	EnableComment  bool       `json:"enableComment"`
	EnableShare    bool       `json:"enableShare"`
}

// PublishDocumentRequest 发布文档（章节）请求
type PublishDocumentRequest struct {
	ChapterTitle  string     `json:"chapterTitle" validate:"required,min=1,max=200"`
	ChapterNumber int        `json:"chapterNumber" validate:"required,min=1"`
	IsFree        bool       `json:"isFree"`
	AuthorNote    string     `json:"authorNote" validate:"max=500"`
	PublishTime   *time.Time `json:"publishTime,omitempty"`
}

// UpdateDocumentPublishStatusRequest 更新文档发布状态请求
type UpdateDocumentPublishStatusRequest struct {
	IsPublished     bool       `json:"isPublished"`
	IsFree          bool       `json:"isFree"`
	ChapterNumber   *int       `json:"chapterNumber,omitempty" validate:"omitempty,min=1"`
	PublishTime     *time.Time `json:"publishTime,omitempty"`
	UnpublishReason string     `json:"unpublishReason,omitempty" validate:"max=500"`
}

// BatchPublishDocumentsRequest 批量发布文档请求
type BatchPublishDocumentsRequest struct {
	DocumentIDs   []string   `json:"documentIds" validate:"required,min=1,max=50"`
	AutoNumbering bool       `json:"autoNumbering"`
	StartNumber   int        `json:"startNumber" validate:"min=1"`
	IsFree        bool       `json:"isFree"`
	PublishTime   *time.Time `json:"publishTime,omitempty"`
}

// BatchPublishResult 批量发布结果
type BatchPublishResult struct {
	SuccessCount int                `json:"successCount"`
	FailCount    int                `json:"failCount"`
	Results      []BatchPublishItem `json:"results"`
}

// BatchPublishItem 批量发布项
type BatchPublishItem struct {
	DocumentID string `json:"documentId"`
	Success    bool   `json:"success"`
	RecordID   string `json:"recordId,omitempty"`
	Error      string `json:"error,omitempty"`
}

// PublicationRecord 发布记录
type PublicationRecord struct {
	ID              string              `json:"id"`
	Type            string              `json:"type"` // project, document
	ResourceID      string              `json:"resourceId"`
	ResourceTitle   string              `json:"resourceTitle"`
	BookstoreID     string              `json:"bookstoreId"`
	BookstoreName   string              `json:"bookstoreName"`
	Status          string              `json:"status"`
	PublishTime     *time.Time          `json:"publishTime"`
	ScheduledTime   *time.Time          `json:"scheduledTime"`
	UnpublishTime   *time.Time          `json:"unpublishTime,omitempty"`
	UnpublishReason string              `json:"unpublishReason,omitempty"`
	Metadata        PublicationMetadata  `json:"metadata"`
	CreatedBy       string              `json:"createdBy"`
	CreatedAt       time.Time           `json:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt"`
}

// PublicationMetadata 发布元数据
type PublicationMetadata struct {
	// 项目发布元数据
	CategoryID    string   `json:"categoryId,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Description   string   `json:"description,omitempty"`
	CoverImage    string   `json:"coverImage,omitempty"`
	PublishType   string   `json:"publishType,omitempty"`
	Price         float64  `json:"price,omitempty"`
	FreeChapters  int      `json:"freeChapters,omitempty"`
	AuthorNote    string   `json:"authorNote,omitempty"`
	EnableComment bool     `json:"enableComment,omitempty"`
	EnableShare   bool     `json:"enableShare,omitempty"`

	// 文档发布元数据
	ChapterTitle  string `json:"chapterTitle,omitempty"`
	ChapterNumber int    `json:"chapterNumber,omitempty"`
	IsFree        bool   `json:"isFree,omitempty"`
}

// PublicationStatus 发布状态
type PublicationStatus struct {
	ProjectID           string                `json:"projectId"`
	ProjectTitle        string                `json:"projectTitle"`
	IsPublished         bool                  `json:"isPublished"`
	BookstoreID         string                `json:"bookstoreId,omitempty"`
	BookstoreName       string                `json:"bookstoreName,omitempty"`
	PublishedAt         *time.Time            `json:"publishedAt,omitempty"`
	TotalChapters       int                   `json:"totalChapters"`
	PublishedChapters   int                   `json:"publishedChapters"`
	UnpublishedChapters int                   `json:"unpublishedChapters"`
	PendingChapters     int                   `json:"pendingChapters"`
	Statistics          PublicationStatistics `json:"statistics"`
}

// PublicationStatistics 发布统计
type PublicationStatistics struct {
	TotalViews     int64     `json:"totalViews"`
	TotalLikes     int64     `json:"totalLikes"`
	TotalComments  int64     `json:"totalComments"`
	TotalShares    int64     `json:"totalShares"`
	TotalFavorites int64     `json:"totalFavorites"`
	LastSyncTime   time.Time `json:"lastSyncTime"`
}

// ============================================================================
// 发布状态常量
// ============================================================================

const (
	// PublicationStatusPending 等待发布
	PublicationStatusPending = "pending"
	// PublicationStatusPublished 已发布
	PublicationStatusPublished = "published"
	// PublicationStatusUnpublished 已取消发布
	PublicationStatusUnpublished = "unpublished"
	// PublicationStatusFailed 发布失败
	PublicationStatusFailed = "failed"
)

// ============================================================================
// 发布类型常量
// ============================================================================

const (
	// PublishTypeSerial 连载
	PublishTypeSerial = "serial"
	// PublishTypeComplete 完本
	PublishTypeComplete = "complete"
)

// ============================================================================
// 导出 DTO
// ============================================================================

// ExportDocumentRequest 导出文档请求
type ExportDocumentRequest struct {
	Format  string         `json:"format" validate:"required,oneof=txt md pdf docx epub"`
	Options *ExportOptions `json:"options,omitempty"`
}

// ExportOptions 导出选项
type ExportOptions struct {
	TOC         bool `json:"toc"`         // 目录
	IncludeMeta bool `json:"includeMeta"` // 包含元数据
}

// ExportProjectRequest 导出项目请求
type ExportProjectRequest struct {
	Format          string `json:"format" validate:"required,oneof=txt md pdf docx epub"`
	IncludeDocuments bool   `json:"includeDocuments"`
	IncludeMeta     bool   `json:"includeMeta"`
}

// ExportTask 导出任务
type ExportTask struct {
	ID            string     `json:"id"`
	Type          string     `json:"type"` // project, document
	ResourceID    string     `json:"resourceId"`
	ResourceTitle string     `json:"resourceTitle"`
	Format        string     `json:"format"`
	Status        string     `json:"status"`
	Progress      int        `json:"progress"`
	FileURL       string     `json:"fileUrl"`
	FileSize      int64      `json:"fileSize"`
	ErrorMsg      string     `json:"errorMsg"`
	CreatedBy     string     `json:"createdBy"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
	ExpiresAt     time.Time  `json:"expiresAt"`
}

// ExportFile 导出文件
type ExportFile struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
	MimeType string `json:"mimeType"`
	FileSize int64  `json:"fileSize"`
}

// ============================================================================
// 文档树 DTO
// ============================================================================

// DocumentTreeNode 文档树节点
type DocumentTreeNode struct {
	Document interface{}        `json:"document"`
	Children []*DocumentTreeNode `json:"children"`
}

// DocumentTreeResponse 文档树响应
type DocumentTreeResponse struct {
	ProjectID string              `json:"projectId"`
	Documents []*DocumentTreeNode `json:"documents"`
}
