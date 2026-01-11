package writer

import (
	"time"

	"Qingyu_backend/models/writer/base"
	"Qingyu_backend/models/writer/types"
)

// Project 项目模型
type Project struct {
	base.IdentifiedEntity `bson:",inline"` // ID
	base.OwnedEntity      `bson:",inline"` // AuthorID
	base.TitledEntity     `bson:",inline"` // Title
	base.Timestamps       `bson:",inline"` // CreatedAt, UpdatedAt, DeletedAt

	// 写作类型（保持BSON字段名不变，确保数据库兼容）
	WritingType string `bson:"writing_type" json:"writingType" validate:"required"` // novel, article, script

	// 详细信息（保持BSON字段名不变，确保数据库兼容）
	Summary    string        `bson:"summary,omitempty" json:"summary,omitempty" validate:"max=500"`
	CoverURL   string        `bson:"cover_url,omitempty" json:"coverUrl,omitempty" validate:"omitempty,url"`
	Status     ProjectStatus `bson:"status" json:"status" validate:"required"`
	Category   string        `bson:"category,omitempty" json:"category,omitempty" validate:"max=50"`
	Tags       []string      `bson:"tags,omitempty" json:"tags,omitempty" validate:"max=10"`
	Visibility Visibility    `bson:"visibility" json:"visibility" validate:"required"`

	// 统计信息
	Statistics ProjectStats `bson:"statistics" json:"statistics"`

	// 设置信息
	Settings ProjectSettings `bson:"settings" json:"settings"`

	// 协作信息
	Collaborators []Collaborator `bson:"collaborators,omitempty" json:"collaborators,omitempty"`

	// 发布时间（项目特有字段）
	PublishedAt *time.Time `bson:"published_at,omitempty" json:"publishedAt,omitempty"`
}

// ProjectStatus 项目状态
type ProjectStatus string

const (
	StatusDraft       ProjectStatus = "draft"       // 草稿
	StatusSerializing ProjectStatus = "serializing" // 连载中
	StatusCompleted   ProjectStatus = "completed"   // 已完结
	StatusSuspended   ProjectStatus = "suspended"   // 暂停
	StatusArchived    ProjectStatus = "archived"    // 已归档
)

// Visibility 可见性
type Visibility string

const (
	VisibilityPrivate Visibility = "private" // 私密
	VisibilityPublic  Visibility = "public"  // 公开
)

// ProjectStats 项目统计
type ProjectStats struct {
	TotalWords    int       `bson:"total_words" json:"totalWords"`       // 总字数
	ChapterCount  int       `bson:"chapter_count" json:"chapterCount"`   // 章节数
	DocumentCount int       `bson:"document_count" json:"documentCount"` // 文档数
	LastUpdateAt  time.Time `bson:"last_update_at" json:"lastUpdateAt"`  // 最后更新时间
}

// ProjectSettings 项目设置
type ProjectSettings struct {
	AutoBackup     bool `bson:"auto_backup" json:"autoBackup"`                            // 自动备份
	BackupInterval int  `bson:"backup_interval" json:"backupInterval"`                    // 备份间隔（小时）
	WordCountGoal  int  `bson:"word_count_goal,omitempty" json:"wordCountGoal,omitempty"` // 字数目标
}

// Collaborator 协作者
type Collaborator struct {
	UserID     string           `bson:"user_id" json:"userId"`
	Role       CollaboratorRole `bson:"role" json:"role"`
	InvitedAt  time.Time        `bson:"invited_at" json:"invitedAt"`
	AcceptedAt *time.Time       `bson:"accepted_at,omitempty" json:"acceptedAt,omitempty"`
}

// CollaboratorRole 协作者角色
type CollaboratorRole string

const (
	RoleOwner  CollaboratorRole = "owner"  // 所有者
	RoleEditor CollaboratorRole = "editor" // 编辑者
	RoleViewer CollaboratorRole = "viewer" // 查看者
)

// IsOwner 判断用户是否为项目所有者
func (p *Project) IsOwner(userID string) bool {
	return p.OwnedEntity.AuthorID == userID
}

// CanEdit 判断用户是否可以编辑项目
func (p *Project) CanEdit(userID string) bool {
	if p.IsOwner(userID) {
		return true
	}

	for _, collab := range p.Collaborators {
		if collab.UserID == userID && collab.Role == RoleEditor && collab.AcceptedAt != nil {
			return true
		}
	}

	return false
}

// CanView 判断用户是否可以查看项目
func (p *Project) CanView(userID string) bool {
	if p.CanEdit(userID) {
		return true
	}

	for _, collab := range p.Collaborators {
		if collab.UserID == userID && collab.AcceptedAt != nil {
			return true
		}
	}

	return p.Visibility == VisibilityPublic
}

// UpdateStatistics 更新项目统计信息
func (p *Project) UpdateStatistics(stats ProjectStats) {
	p.Statistics = stats
	p.UpdatedAt = time.Now()
}

// IsValid 验证项目状态值是否有效
func (s ProjectStatus) IsValid() bool {
	switch s {
	case StatusDraft, StatusSerializing, StatusCompleted, StatusSuspended, StatusArchived:
		return true
	}
	return false
}

// String 返回项目状态的字符串表示
func (s ProjectStatus) String() string {
	return string(s)
}

// IsValid 验证可见性值是否有效
func (v Visibility) IsValid() bool {
	switch v {
	case VisibilityPrivate, VisibilityPublic:
		return true
	}
	return false
}

// String 返回可见性的字符串表示
func (v Visibility) String() string {
	return string(v)
}

// IsValid 验证协作者角色是否有效
func (r CollaboratorRole) IsValid() bool {
	switch r {
	case RoleOwner, RoleEditor, RoleViewer:
		return true
	}
	return false
}

// String 返回角色的字符串表示
func (r CollaboratorRole) String() string {
	return string(r)
}

// Validate 验证项目数据的有效性
func (p *Project) Validate() error {
	// 验证标题
	if err := base.ValidateTitle(p.Title, 100); err != nil {
		return err
	}

	// 验证作者ID
	if p.AuthorID == "" {
		return base.ErrAuthorIDRequired
	}

	// 验证写作类型
	if p.WritingType == "" {
		p.WritingType = "novel" // 默认为小说类型
	}
	if !types.GlobalRegistry.ValidateWritingType(p.WritingType) {
		return base.ErrInvalidWritingType
	}

	// 验证封面URL
	if err := base.ValidateURL(p.CoverURL); err != nil {
		return err
	}

	// 验证状态和可见性
	if !p.Status.IsValid() {
		return base.ErrInvalidStatus
	}
	if !p.Visibility.IsValid() {
		return base.ErrInvalidVisibility
	}

	return nil
}
