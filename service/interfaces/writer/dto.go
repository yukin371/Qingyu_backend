package writer

// ============================================================================
// Writer 模块 DTO 类型定义（已废弃）
// ============================================================================

// 本文件中的DTO定义已迁移至 models/dto/content_dto.go
// 为了保持向后兼容，这里重新导出新定义
//
// 迁移指南：
// - 将 import "Qingyu_backend/service/interfaces/writer" 改为 "Qingyu_backend/models/dto"
// - 或者直接使用 writer.CreateProjectRequest（它会自动指向 dto.CreateProjectRequest）
//
// 废弃时间：2026-02-26
// 计划移除时间：2026-06-01

import "Qingyu_backend/models/dto"

// ============================================================================
// 项目管理 DTO（重新导出）
// ============================================================================

// CreateProjectRequest 创建项目请求
// Deprecated: 使用 dto.CreateProjectRequest 替代
type CreateProjectRequest = dto.CreateProjectRequest

// CreateProjectResponse 创建项目响应
// Deprecated: 使用 dto.CreateProjectResponse 替代
type CreateProjectResponse = dto.CreateProjectResponse

// ListProjectsRequest 列出项目请求
// Deprecated: 使用 dto.ListProjectsRequest 替代
type ListProjectsRequest = dto.ListProjectsRequest

// ListProjectsResponse 列出项目响应
// Deprecated: 使用 dto.ListProjectsResponse 替代
type ListProjectsResponse = dto.ListProjectsResponse

// UpdateProjectRequest 更新项目请求
// Deprecated: 使用 dto.UpdateProjectRequest 替代
type UpdateProjectRequest = dto.UpdateProjectRequest

// GetProjectResponse 获取项目详情响应
// Deprecated: 使用 dto.GetProjectResponse 替代
type GetProjectResponse = dto.GetProjectResponse

// UpdateStatisticsRequest 更新统计请求
// Deprecated: 使用 dto.UpdateStatisticsRequest 替代
type UpdateStatisticsRequest = dto.UpdateStatisticsRequest

// AddCollaboratorRequest 添加协作者请求
// Deprecated: 使用 dto.AddCollaboratorRequest 替代
type AddCollaboratorRequest = dto.AddCollaboratorRequest

// RemoveCollaboratorRequest 移除协作者请求
// Deprecated: 使用 dto.RemoveCollaboratorRequest 替代
type RemoveCollaboratorRequest = dto.RemoveCollaboratorRequest

// ============================================================================
// 文档管理 DTO（重新导出）
// ============================================================================

// CreateDocumentRequest 创建文档请求
// Deprecated: 使用 dto.CreateDocumentRequest 替代
type CreateDocumentRequest = dto.CreateDocumentRequest

// CreateDocumentResponse 创建文档响应
// Deprecated: 使用 dto.CreateDocumentResponse 替代
type CreateDocumentResponse = dto.CreateDocumentResponse

// DocumentTreeResponse 文档树响应
// Deprecated: 使用 dto.DocumentTreeResponse 替代
type DocumentTreeResponse = dto.DocumentTreeResponse

// TreeNode 树节点（保留旧定义以兼容）
// Deprecated: 直接使用 models 中的类型
type TreeNode struct {
	Document interface{}   `json:"document"`
	Children []*TreeNode   `json:"children"`
}

// ListDocumentsRequest 列出文档请求
// Deprecated: 使用 dto.ListDocumentsRequest 替代
type ListDocumentsRequest = dto.ListDocumentsRequest

// ListDocumentsResponse 列出文档响应
// Deprecated: 使用 dto.ListDocumentsResponse 替代
type ListDocumentsResponse = dto.ListDocumentsResponse

// UpdateDocumentRequest 更新文档请求
// Deprecated: 使用 dto.UpdateDocumentRequest 替代
type UpdateDocumentRequest = dto.UpdateDocumentRequest

// MoveDocumentRequest 移动文档请求
// Deprecated: 使用 dto.MoveDocumentRequest 替代
type MoveDocumentRequest = dto.MoveDocumentRequest

// ReorderDocumentsRequest 重新排序文档请求
// Deprecated: 使用 dto.ReorderDocumentsRequest 替代
type ReorderDocumentsRequest = dto.ReorderDocumentsRequest

// AutoSaveRequest 自动保存请求
// Deprecated: 使用 dto.AutoSaveRequest 替代
type AutoSaveRequest = dto.AutoSaveRequest

// AutoSaveResponse 自动保存响应
// Deprecated: 使用 dto.AutoSaveResponse 替代
type AutoSaveResponse = dto.AutoSaveResponse

// SaveStatusResponse 保存状态响应
// Deprecated: 使用 dto.SaveStatusResponse 替代
type SaveStatusResponse = dto.SaveStatusResponse

// DocumentContentResponse 文档内容响应
// Deprecated: 使用 dto.DocumentContentResponse 替代
type DocumentContentResponse = dto.DocumentContentResponse

// UpdateContentRequest 更新内容请求
// Deprecated: 使用 dto.UpdateContentRequest 替代
type UpdateContentRequest = dto.UpdateContentRequest

// DuplicateRequest 复制请求
// Deprecated: 使用 dto.DuplicateRequest 替代
type DuplicateRequest = dto.DuplicateRequest

// DuplicateResponse 复制响应
// Deprecated: 使用 dto.DuplicateResponse 替代
type DuplicateResponse = dto.DuplicateResponse

// 版本控制 DTO

// VersionHistoryResponse 版本历史响应
// Deprecated: 使用 dto.VersionHistoryResponse 替代
type VersionHistoryResponse = dto.VersionHistoryResponse

// VersionInfo 版本信息
// Deprecated: 使用 dto.VersionInfo 替代
type VersionInfo = dto.VersionInfo

// VersionDetail 版本详情
// Deprecated: 使用 dto.VersionDetail 替代
type VersionDetail = dto.VersionDetail

// VersionDiff 版本差异
// Deprecated: 使用 dto.VersionDiff 替代
type VersionDiff = dto.VersionDiff

// ChangeItem 变更项
// Deprecated: 使用 dto.ChangeItem 替代
type ChangeItem = dto.ChangeItem

// ============================================================================
// 内容管理 DTO（重新导出）
// ============================================================================

// CreateCharacterRequest 创建角色请求
// Deprecated: 使用 dto.CreateCharacterRequest 替代
type CreateCharacterRequest = dto.CreateCharacterRequest

// UpdateCharacterRequest 更新角色请求
// Deprecated: 使用 dto.UpdateCharacterRequest 替代
type UpdateCharacterRequest = dto.UpdateCharacterRequest

// CreateRelationRequest 创建关系请求
// Deprecated: 使用 dto.CreateRelationRequest 替代
type CreateRelationRequest = dto.CreateRelationRequest

// CharacterGraph 角色关系图
// Deprecated: 使用 dto.CharacterGraph 替代
type CharacterGraph = dto.CharacterGraph

// CreateLocationRequest 创建地点请求
// Deprecated: 使用 dto.CreateLocationRequest 替代
type CreateLocationRequest = dto.CreateLocationRequest

// UpdateLocationRequest 更新地点请求
// Deprecated: 使用 dto.UpdateLocationRequest 替代
type UpdateLocationRequest = dto.UpdateLocationRequest

// LocationNode 地点节点
// Deprecated: 使用 dto.LocationNode 替代
type LocationNode = dto.LocationNode

// CreateLocationRelationRequest 创建地点关系请求
// Deprecated: 使用 dto.CreateLocationRelationRequest 替代
type CreateLocationRelationRequest = dto.CreateLocationRelationRequest

// CreateTimelineRequest 创建时间线请求
// Deprecated: 使用 dto.CreateTimelineRequest 替代
type CreateTimelineRequest = dto.CreateTimelineRequest

// CreateTimelineEventRequest 创建时间线事件请求
// Deprecated: 使用 dto.CreateTimelineEventRequest 替代
type CreateTimelineEventRequest = dto.CreateTimelineEventRequest

// UpdateTimelineEventRequest 更新时间线事件请求
// Deprecated: 使用 dto.UpdateTimelineEventRequest 替代
type UpdateTimelineEventRequest = dto.UpdateTimelineEventRequest

// TimelineVisualization 时间线可视化
// Deprecated: 使用 dto.TimelineVisualization 替代
type TimelineVisualization = dto.TimelineVisualization

// TimelineEventNode 时间线事件节点
// Deprecated: 使用 dto.TimelineEventNode 替代
type TimelineEventNode = dto.TimelineEventNode

// EventConnection 事件连接
// Deprecated: 使用 dto.EventConnection 替代
type EventConnection = dto.EventConnection

// ============================================================================
// 发布导出 DTO（重新导出）
// ============================================================================

// PublishProjectRequest 发布项目请求
// Deprecated: 使用 dto.PublishProjectRequest 替代
type PublishProjectRequest = dto.PublishProjectRequest

// PublishDocumentRequest 发布文档请求
// Deprecated: 使用 dto.PublishDocumentRequest 替代
type PublishDocumentRequest = dto.PublishDocumentRequest

// UpdateDocumentPublishStatusRequest 更新文档发布状态请求
// Deprecated: 使用 dto.UpdateDocumentPublishStatusRequest 替代
type UpdateDocumentPublishStatusRequest = dto.UpdateDocumentPublishStatusRequest

// BatchPublishDocumentsRequest 批量发布文档请求
// Deprecated: 使用 dto.BatchPublishDocumentsRequest 替代
type BatchPublishDocumentsRequest = dto.BatchPublishDocumentsRequest

// BatchPublishResult 批量发布结果
// Deprecated: 使用 dto.BatchPublishResult 替代
type BatchPublishResult = dto.BatchPublishResult

// BatchPublishItem 批量发布项
// Deprecated: 使用 dto.BatchPublishItem 替代
type BatchPublishItem = dto.BatchPublishItem

// PublicationRecord 发布记录
// Deprecated: 使用 dto.PublicationRecord 替代
type PublicationRecord = dto.PublicationRecord

// PublicationStatus 发布状态
// Deprecated: 使用 dto.PublicationStatus 替代
type PublicationStatus = dto.PublicationStatus

// ExportDocumentRequest 导出文档请求
// Deprecated: 使用 dto.ExportDocumentRequest 替代
type ExportDocumentRequest = dto.ExportDocumentRequest

// ExportOptions 导出选项
// Deprecated: 使用 dto.ExportOptions 替代
type ExportOptions = dto.ExportOptions

// ExportProjectRequest 导出项目请求
// Deprecated: 使用 dto.ExportProjectRequest 替代
type ExportProjectRequest = dto.ExportProjectRequest

// ExportTask 导出任务
// Deprecated: 使用 dto.ExportTask 替代
type ExportTask = dto.ExportTask

// ExportFile 导出文件
// Deprecated: 使用 dto.ExportFile 替代
type ExportFile = dto.ExportFile
