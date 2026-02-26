package content

import (
	"context"

	"Qingyu_backend/models/dto"
)

// ============================================================================
// 项目服务接口详细定义
// ============================================================================

// ProjectService 项目管理服务接口
// 提供项目的完整生命周期管理功能
type ProjectService interface {
	// === 基础 CRUD 操作 ===

	// CreateProject 创建新项目
	// 创建新的写作项目，初始化项目结构和统计信息
	CreateProject(ctx context.Context, req *dto.CreateProjectRequest) (*dto.ProjectResponse, error)

	// UpdateProject 更新项目信息
	// 更新项目的标题、简介、封面、分类、标签等信息
	UpdateProject(ctx context.Context, id string, req *dto.UpdateProjectRequest) (*dto.ProjectResponse, error)

	// GetProject 获取项目详情
	// 返回项目的完整信息，包括统计数据
	GetProject(ctx context.Context, id string) (*dto.ProjectResponse, error)

	// DeleteProject 删除项目
	// 执行软删除，将项目标记为已删除状态
	DeleteProject(ctx context.Context, id string) error

	// ListProjects 获取项目列表
	// 支持按状态、分类分页查询项目列表
	ListProjects(ctx context.Context, req *dto.ListProjectsRequest) (*dto.ListProjectsResponse, error)

	// === 项目统计 ===

	// GetProjectStatistics 获取项目统计信息
	// 返回项目的详细统计数据，包括字数、章节数、角色数等
	GetProjectStatistics(ctx context.Context, projectID string) (*dto.ProjectStatistics, error)

	// UpdateProjectStatistics 更新项目统计
	// 更新项目的统计数据，通常在文档变更时调用
	UpdateProjectStatistics(ctx context.Context, projectID string, stats *dto.ProjectStatistics) error

	// RecalculateProjectStatistics 重新计算项目统计
	// 重新计算并更新项目的所有统计数据
	RecalculateProjectStatistics(ctx context.Context, projectID string) (*dto.ProjectStatistics, error)

	// === 项目操作 ===

	// DuplicateProject 复制项目
	// 创建项目的副本，包括文档结构和内容
	DuplicateProject(ctx context.Context, projectID string, req *dto.DuplicateProjectRequest) (*dto.ProjectResponse, error)

	// ArchiveProject 归档项目
	// 将项目标记为归档状态
	ArchiveProject(ctx context.Context, projectID string) error

	// UnarchiveProject 取消归档项目
	// 将项目从归档状态恢复
	UnarchiveProject(ctx context.Context, projectID string) error

	// PublishProject 发布项目
	// 将项目发布到书城
	PublishProject(ctx context.Context, projectID string, req *dto.PublishProjectRequest) (*dto.PublicationRecord, error)

	// UnpublishProject 取消发布项目
	// 将项目从书城下架
	UnpublishProject(ctx context.Context, projectID string) error

	// === 协作管理 ===

	// AddCollaborator 添加协作者
	// 添加用户到项目协作团队
	AddCollaborator(ctx context.Context, projectID, userID, role string) error

	// RemoveCollaborator 移除协作者
	// 从项目协作团队移除用户
	RemoveCollaborator(ctx context.Context, projectID, userID string) error

	// UpdateCollaboratorRole 更新协作者角色
	// 更新协作者在项目中的角色权限
	UpdateCollaboratorRole(ctx context.Context, projectID, userID, role string) error

	// ListCollaborators 获取协作者列表
	// 返回项目的所有协作者及其角色
	ListCollaborators(ctx context.Context, projectID string) ([]*dto.CollaboratorInfo, error)

	// === 项目导出 ===

	// ExportProject 导出项目
	// 将项目导出为指定格式的文件
	ExportProject(ctx context.Context, projectID string, req *dto.ExportProjectRequest) (*dto.ExportTask, error)

	// GetExportTasks 获取导出任务列表
	// 返回项目的所有导出任务
	GetExportTasks(ctx context.Context, projectID string) ([]*dto.ExportTask, error)

	// === 批量操作 ===

	// BatchUpdateProjects 批量更新项目
	// 批量更新多个项目的属性
	BatchUpdateProjects(ctx context.Context, projectIDs []string, updates map[string]interface{}) error

	// BatchDeleteProjects 批量删除项目
	// 批量软删除多个项目
	BatchDeleteProjects(ctx context.Context, projectIDs []string) error

	// BatchArchiveProjects 批量归档项目
	// 批量将项目标记为归档状态
	BatchArchiveProjects(ctx context.Context, projectIDs []string) error

	// === 项目搜索 ===

	// SearchProjects 搜索项目
	// 根据关键词搜索项目标题和简介
	SearchProjects(ctx context.Context, userID, keyword string, page, pageSize int) (*dto.ListProjectsResponse, error)

	// GetProjectsByTags 根据标签获取项目
	// 按标签筛选项目列表
	GetProjectsByTags(ctx context.Context, userID string, tags []string, page, pageSize int) (*dto.ListProjectsResponse, error)
}
