package writer

// ============================================================================
// Writer 模块 DTO 类型定义
// ============================================================================

// ============================================================================
// 项目管理 DTO
// ============================================================================

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	CoverURL  string   `json:"cover_url"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	ProjectID string   `json:"project_id"` // 用于模板创建
}

// CreateProjectResponse 创建项目响应
type CreateProjectResponse struct {
	ProjectID string `json:"project_id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// ListProjectsRequest 列出项目请求
type ListProjectsRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Status   string `json:"status"`
}

// ListProjectsResponse 列出项目响应
type ListProjectsResponse struct {
	Projects interface{} `json:"projects"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	Title    string   `json:"title"`
	Summary  string   `json:"summary"`
	CoverURL string   `json:"cover_url"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Status   string   `json:"status"`
}

// ============================================================================
// 文档管理 DTO
// ============================================================================

// CreateDocumentRequest 创建文档请求
type CreateDocumentRequest struct {
	ProjectID    string   `json:"project_id"`
	ParentID     string   `json:"parent_id"`
	Title        string   `json:"title"`
	Type         string   `json:"type"`
	Level        int      `json:"level"`
	Order        int      `json:"order"`
	CharacterIDs []string `json:"character_ids"`
	LocationIDs  []string `json:"location_ids"`
	TimelineIDs  []string `json:"timeline_ids"`
	Tags         []string `json:"tags"`
	Notes        string   `json:"notes"`
}

// CreateDocumentResponse 创建文档响应
type CreateDocumentResponse struct {
	DocumentID string `json:"document_id"`
	Title      string `json:"title"`
	Type       string `json:"type"`
	CreatedAt  string `json:"created_at"`
}

// DocumentTreeResponse 文档树响应
type DocumentTreeResponse struct {
	ProjectID string       `json:"project_id"`
	Documents []*TreeNode `json:"documents"`
}

// TreeNode 树节点
type TreeNode struct {
	Document interface{}   `json:"document"`
	Children []*TreeNode   `json:"children"`
}

// ListDocumentsRequest 列出文档请求
type ListDocumentsRequest struct {
	ProjectID string `json:"project_id"`
	ParentID  string `json:"parent_id"`
	Page      string `json:"page"`
	PageSize  string `json:"page_size"`
	Status    string `json:"status"`
}

// ListDocumentsResponse 列出文档响应
type ListDocumentsResponse struct {
	Documents interface{} `json:"documents"`
	Total     int         `json:"total"`
	Page      int         `json:"page"`
	PageSize  int         `json:"page_size"`
}

// UpdateDocumentRequest 更新文档请求
type UpdateDocumentRequest struct {
	Title         string   `json:"title"`
	Notes         string   `json:"notes"`
	Status        string   `json:"status"`
	CharacterIDs  []string `json:"character_ids"`
	LocationIDs   []string `json:"location_ids"`
	TimelineIDs   []string `json:"timeline_ids"`
	Tags          []string `json:"tags"`
}

// MoveDocumentRequest 移动文档请求
type MoveDocumentRequest struct {
	DocumentID string `json:"document_id"`
	NewParentID string `json:"new_parent_id"`
	Order      int    `json:"order"`
}

// ReorderDocumentsRequest 重新排序文档请求
type ReorderDocumentsRequest struct {
	ProjectID string         `json:"project_id"`
	ParentID  string         `json:"parent_id"`
	Orders    map[string]int `json:"orders"`
}

// AutoSaveRequest 自动保存请求
type AutoSaveRequest struct {
	DocumentID     string `json:"document_id"`
	Content        string `json:"content"`
	CurrentVersion int    `json:"current_version"`
	SaveType       string `json:"save_type"`
}

// AutoSaveResponse 自动保存响应
type AutoSaveResponse struct {
	Saved       bool   `json:"saved"`
	NewVersion  int    `json:"new_version"`
	WordCount   int    `json:"word_count"`
	SavedAt     string `json:"saved_at"`
	HasConflict bool   `json:"has_conflict"`
}

// SaveStatusResponse 保存状态响应
type SaveStatusResponse struct {
	DocumentID     string `json:"document_id"`
	LastSavedAt    string `json:"last_saved_at"`
	CurrentVersion int    `json:"current_version"`
	IsSaving       bool   `json:"is_saving"`
	WordCount      int    `json:"word_count"`
}

// DocumentContentResponse 文档内容响应
type DocumentContentResponse struct {
	DocumentID string `json:"document_id"`
	Content    string `json:"content"`
	Version    int    `json:"version"`
	WordCount  int    `json:"word_count"`
	UpdatedAt  string `json:"updated_at"`
}

// UpdateContentRequest 更新内容请求
type UpdateContentRequest struct {
	DocumentID string `json:"document_id"`
	Content    string `json:"content"`
	Version    int    `json:"version"`
}

// DuplicateRequest 复制请求
type DuplicateRequest struct {
	NewTitle      string   `json:"new_title"`
	TargetProject string   `json:"target_project"`
	CopyContent   bool     `json:"copy_content"`
}

// DuplicateResponse 复制响应
type DuplicateResponse struct {
	NewDocumentID string `json:"new_document_id"`
	Title         string `json:"title"`
}

// 版本控制 DTO

// VersionHistoryResponse 版本历史响应
type VersionHistoryResponse struct {
	Versions []*VersionInfo `json:"versions"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// VersionInfo 版本信息
type VersionInfo struct {
	VersionID  string `json:"version_id"`
	Version    int    `json:"version"`
	Message    string `json:"message"`
	CreatedAt  string `json:"created_at"`
	CreatedBy  string `json:"created_by"`
	WordCount  int    `json:"word_count"`
}

// VersionDetail 版本详情
type VersionDetail struct {
	VersionID  string `json:"version_id"`
	DocumentID string `json:"document_id"`
	Version    int    `json:"version"`
	Content    string `json:"content"`
	Message    string `json:"message"`
	CreatedAt  string `json:"created_at"`
	CreatedBy  string `json:"created_by"`
	WordCount  int    `json:"word_count"`
}

// VersionDiff 版本差异
type VersionDiff struct {
	FromVersion   string      `json:"from_version"`
	ToVersion     string      `json:"to_version"`
	Changes       []ChangeItem `json:"changes"`
	AddedLines    int         `json:"added_lines"`
	DeletedLines  int         `json:"deleted_lines"`
}

// ChangeItem 变更项
type ChangeItem struct {
	Type    string `json:"type"`    // added, deleted, modified
	Line    int    `json:"line"`
	Content string `json:"content"`
}

// ============================================================================
// 内容管理 DTO
// ============================================================================

// CreateCharacterRequest 创建角色请求
type CreateCharacterRequest struct {
	Name             string   `json:"name"`
	Alias            string   `json:"alias"`
	Summary          string   `json:"summary"`
	Traits           []string `json:"traits"`
	Background       string   `json:"background"`
	AvatarURL        string   `json:"avatar_url"`
	PersonalityPrompt string  `json:"personality_prompt"`
	SpeechPattern    string   `json:"speech_pattern"`
	CurrentState     string   `json:"current_state"`
}

// UpdateCharacterRequest 更新角色请求
type UpdateCharacterRequest struct {
	Name             *string  `json:"name"`
	Alias            *string  `json:"alias"`
	Summary          *string  `json:"summary"`
	Traits           []string `json:"traits"`
	Background       *string  `json:"background"`
	AvatarURL        *string  `json:"avatar_url"`
	PersonalityPrompt *string  `json:"personality_prompt"`
	SpeechPattern    *string  `json:"speech_pattern"`
	CurrentState     *string  `json:"current_state"`
}

// CreateRelationRequest 创建关系请求
type CreateRelationRequest struct {
	FromID   string  `json:"from_id"`
	ToID     string  `json:"to_id"`
	Type     string  `json:"type"`
	Strength int    `json:"strength"`
	Notes    string  `json:"notes"`
}

// CharacterGraph 角色关系图
type CharacterGraph struct {
	Nodes interface{} `json:"nodes"`
	Edges interface{} `json:"edges"`
}

// CreateLocationRequest 创建地点请求
type CreateLocationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Climate     string `json:"climate"`
	Culture     string `json:"culture"`
	Geography   string `json:"geography"`
	Atmosphere  string `json:"atmosphere"`
	ParentID    string `json:"parent_id"`
	ImageURL    string `json:"image_url"`
}

// UpdateLocationRequest 更新地点请求
type UpdateLocationRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Climate     *string `json:"climate"`
	Culture     *string `json:"culture"`
	Geography   *string `json:"geography"`
	Atmosphere  *string `json:"atmosphere"`
	ParentID    *string `json:"parent_id"`
	ImageURL    *string `json:"image_url"`
}

// LocationNode 地点节点
type LocationNode struct {
	Location interface{}     `json:"location"`
	Children []*LocationNode `json:"children"`
}

// CreateLocationRelationRequest 创建地点关系请求
type CreateLocationRelationRequest struct {
	FromID   string  `json:"from_id"`
	ToID     string  `json:"to_id"`
	Type     string  `json:"type"`
	Distance *int    `json:"distance"`
	Notes    string  `json:"notes"`
}

// CreateTimelineRequest 创建时间线请求
type CreateTimelineRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

// CreateTimelineEventRequest 创建时间线事件请求
type CreateTimelineEventRequest struct {
	TimelineID   string   `json:"timeline_id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	EventType    string   `json:"event_type"`
	Importance   int      `json:"importance"`
	Participants []string `json:"participants"`
	LocationIDs  []string `json:"location_ids"`
	ChapterIDs   []string `json:"chapter_ids"`
	StoryTime    string   `json:"story_time"`
	Duration     *int     `json:"duration"`
	Impact       string   `json:"impact"`
}

// UpdateTimelineEventRequest 更新时间线事件请求
type UpdateTimelineEventRequest struct {
	Title        *string   `json:"title"`
	Description  *string   `json:"description"`
	EventType    *string   `json:"event_type"`
	Importance   *int      `json:"importance"`
	Participants []string  `json:"participants"`
	LocationIDs  []string  `json:"location_ids"`
	ChapterIDs   []string  `json:"chapter_ids"`
	StoryTime    *string   `json:"story_time"`
	Duration     *int      `json:"duration"`
	Impact       *string   `json:"impact"`
}

// TimelineVisualization 时间线可视化
type TimelineVisualization struct {
	Events      []*TimelineEventNode  `json:"events"`
	Connections []*EventConnection    `json:"connections"`
}

// TimelineEventNode 时间线事件节点
type TimelineEventNode struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	StoryTime  string   `json:"story_time"`
	EventType  string   `json:"event_type"`
	Importance int      `json:"importance"`
	Characters []string `json:"characters"`
}

// EventConnection 事件连接
type EventConnection struct {
	FromEventID string `json:"from_event_id"`
	ToEventID   string `json:"to_event_id"`
	Type        string `json:"type"`
}

// ============================================================================
// 发布导出 DTO
// ============================================================================

// PublishProjectRequest 发布项目请求
type PublishProjectRequest struct {
	BookstoreID    string   `json:"bookstore_id"`
	CategoryID     string   `json:"category_id"`
	Tags           []string `json:"tags"`
	Description    string   `json:"description"`
	CoverImage     string   `json:"cover_image"`
	PublishType    string   `json:"publish_type"`
	Price          *float64 `json:"price"`
	FreeChapters   int      `json:"free_chapters"`
	AuthorNote     string   `json:"author_note"`
	EnableComment  bool     `json:"enable_comment"`
	EnableShare    bool     `json:"enable_share"`
	PublishTime    *string  `json:"publish_time"`
}

// PublishDocumentRequest 发布文档请求
type PublishDocumentRequest struct {
	ChapterTitle  string  `json:"chapter_title"`
	ChapterNumber int     `json:"chapter_number"`
	IsFree        bool    `json:"is_free"`
	AuthorNote    string  `json:"author_note"`
	PublishTime   *string `json:"publish_time"`
}

// UpdateDocumentPublishStatusRequest 更新文档发布状态请求
type UpdateDocumentPublishStatusRequest struct {
	IsPublished     bool    `json:"is_published"`
	IsFree          bool    `json:"is_free"`
	ChapterNumber   *int    `json:"chapter_number"`
	UnpublishReason string  `json:"unpublish_reason"`
}

// BatchPublishDocumentsRequest 批量发布文档请求
type BatchPublishDocumentsRequest struct {
	DocumentIDs    []string `json:"document_ids"`
	IsFree         bool     `json:"is_free"`
	AutoNumbering  bool     `json:"auto_numbering"`
	StartNumber    int      `json:"start_number"`
	PublishTime    *string  `json:"publish_time"`
}

// BatchPublishResult 批量发布结果
type BatchPublishResult struct {
	SuccessCount int                  `json:"success_count"`
	FailCount    int                  `json:"fail_count"`
	Results      []BatchPublishItem   `json:"results"`
}

// BatchPublishItem 批量发布项
type BatchPublishItem struct {
	DocumentID string `json:"document_id"`
	Success    bool   `json:"success"`
	RecordID   string `json:"record_id"`
	Error      string `json:"error"`
}

// PublicationRecord 发布记录
type PublicationRecord struct {
	ID            string                `json:"id"`
	Type          string                `json:"type"`
	ResourceID    string                `json:"resource_id"`
	ResourceTitle string                `json:"resource_title"`
	BookstoreID   string                `json:"bookstore_id"`
	BookstoreName string                `json:"bookstore_name"`
	Status        string                `json:"status"`
	PublishTime   *string               `json:"publish_time"`
	CreatedBy     string                `json:"created_by"`
	CreatedAt     string                `json:"created_at"`
	UpdatedAt     string                `json:"updated_at"`
}

// PublicationStatus 发布状态
type PublicationStatus struct {
	ProjectID         string                `json:"project_id"`
	ProjectTitle      string                `json:"project_title"`
	IsPublished       bool                  `json:"is_published"`
	BookstoreID       string                `json:"bookstore_id"`
	BookstoreName     string                `json:"bookstore_name"`
	PublishedAt       *string               `json:"published_at"`
	Statistics        interface{}           `json:"statistics"`
	TotalChapters     int                   `json:"total_chapters"`
	PublishedChapters int                   `json:"published_chapters"`
}

// ExportDocumentRequest 导出文档请求
type ExportDocumentRequest struct {
	Format  string            `json:"format"`
	Options *ExportOptions     `json:"options"`
}

// ExportOptions 导出选项
type ExportOptions struct {
	TOC        bool `json:"toc"`        // 目录
	IncludeMeta bool `json:"include_meta"` // 包含元数据
}

// ExportProjectRequest 导出项目请求
type ExportProjectRequest struct {
	Format          string  `json:"format"`
	IncludeDocuments bool    `json:"include_documents"`
	IncludeMeta     bool    `json:"include_meta"`
}

// ExportTask 导出任务
type ExportTask struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	ResourceID    string    `json:"resource_id"`
	ResourceTitle string    `json:"resource_title"`
	Format        string    `json:"format"`
	Status        string    `json:"status"`
	Progress      int       `json:"progress"`
	FileURL       string    `json:"file_url"`
	FileSize      int64     `json:"file_size"`
	ErrorMsg      string    `json:"error_msg"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	CompletedAt   *string  `json:"completed_at"`
	ExpiresAt     string    `json:"expires_at"`
}

// ExportFile 导出文件
type ExportFile struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
	FileSize int64  `json:"file_size"`
}
