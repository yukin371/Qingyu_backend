package document

// ============================================================================
// 文档模块 DTO 类型定义（已废弃）
// ============================================================================

// 本文件中的DTO定义已迁移至 models/dto/content_dto.go
// 为了保持向后兼容，这里重新导出新定义
//
// 迁移指南：
// - 将 import "Qingyu_backend/service/writer/document" 改为 "Qingyu_backend/models/dto"
// - 或者直接使用 document.CreateDocumentRequest（它会自动指向 dto.CreateDocumentRequest）
//
// 废弃时间：2026-02-26
// 计划移除时间：2026-06-01

import (
	"Qingyu_backend/models/dto"
	"Qingyu_backend/models/writer"
)

// CreateDocumentRequest 创建文档请求
// Deprecated: 使用 dto.CreateDocumentRequest 替代
type CreateDocumentRequest = dto.CreateDocumentRequest

// CreateDocumentResponse 创建文档响应
// Deprecated: 使用 dto.CreateDocumentResponse 替代
type CreateDocumentResponse = dto.CreateDocumentResponse

// DocumentTreeNode 文档树节点
// Deprecated: 保留文档模块特定类型
type DocumentTreeNode struct {
	*writer.Document
	Children []*DocumentTreeNode `json:"children"`
}

// DocumentTreeResponse 文档树响应
// Deprecated: 保留文档模块特定响应类型
type DocumentTreeResponse struct {
	ProjectID string              `json:"projectId"`
	Documents []*DocumentTreeNode `json:"documents"`
}

// MoveDocumentRequest 移动文档请求
// Deprecated: 使用 dto.MoveDocumentRequest 替代
type MoveDocumentRequest = dto.MoveDocumentRequest

// ReorderDocumentsRequest 重新排序请求
// Deprecated: 使用 dto.ReorderDocumentsRequest 替代
type ReorderDocumentsRequest = dto.ReorderDocumentsRequest

// UpdateDocumentRequest 更新文档请求
// Deprecated: 使用 dto.UpdateDocumentRequest 替代
type UpdateDocumentRequest = dto.UpdateDocumentRequest

// ListDocumentsRequest 文档列表请求
// Deprecated: 使用 dto.ListDocumentsRequest 替代
type ListDocumentsRequest = dto.ListDocumentsRequest

// ListDocumentsResponse 文档列表响应
// Deprecated: 保留文档模块特定响应类型
type ListDocumentsResponse struct {
	Documents []*writer.Document `json:"documents"`
	Total     int                `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"pageSize"`
}

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
