package project

// ============================================================================
// 项目模块 DTO 类型定义（已废弃）
// ============================================================================

// 本文件中的DTO定义已迁移至 models/dto/content_dto.go
// 为了保持向后兼容，这里重新导出新定义
//
// 迁移指南：
// - 将请求 DTO 的 import 从 "Qingyu_backend/service/writer/project" 改为 "Qingyu_backend/models/dto"
// - 本文件仅保留项目模块特定响应类型与少量兼容别名
//
// 废弃时间：2026-02-26
// 计划移除时间：2026-06-01

import (
	"Qingyu_backend/models/dto"
	"Qingyu_backend/models/writer"
)

// CreateProjectResponse 创建项目响应
// Deprecated: 使用 dto.CreateProjectResponse 替代
type CreateProjectResponse = dto.CreateProjectResponse

// ListProjectsResponse 项目列表响应
// Deprecated: 保留项目模块特定响应类型
type ListProjectsResponse struct {
	Projects []*writer.Project `json:"projects"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

// GetProjectResponse 获取项目详情响应
// Deprecated: 保留项目模块特定响应类型
type GetProjectResponse struct {
	*writer.Project
}

// UpdateStatisticsRequest 更新统计请求
// Deprecated: 使用 dto.UpdateStatisticsRequest 替代
type UpdateStatisticsRequest = dto.UpdateStatisticsRequest

// AddCollaboratorRequest 添加协作者请求
// Deprecated: 使用 dto.AddCollaboratorRequest 替代
type AddCollaboratorRequest = dto.AddCollaboratorRequest

// RemoveCollaboratorRequest 移除协作者请求
// Deprecated: 使用 dto.RemoveCollaboratorRequest 替代
type RemoveCollaboratorRequest = dto.RemoveCollaboratorRequest
