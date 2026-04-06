package document

// ============================================================================
// 文档模块 DTO 类型定义（已废弃）
// ============================================================================

// 本文件中的DTO定义已迁移至 models/dto/content_dto.go
// 为了保持向后兼容，这里重新导出新定义
//
// 迁移指南：
// - 将请求 DTO 的 import 从 "Qingyu_backend/service/writer/document" 改为 "Qingyu_backend/models/dto"
// - 本文件仅保留文档模块特定树结构响应与少量兼容别名
//
// 废弃时间：2026-02-26
// 计划移除时间：2026-06-01

import (
	"Qingyu_backend/models/dto"
	"Qingyu_backend/models/writer"
)

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

// ListDocumentsResponse 文档列表响应
// Deprecated: 保留文档模块特定响应类型
type ListDocumentsResponse struct {
	Documents []*writer.Document `json:"documents"`
	Total     int                `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"pageSize"`
}

// AutoSaveResponse 自动保存响应
// Deprecated: 使用 dto.AutoSaveResponse 替代
type AutoSaveResponse = dto.AutoSaveResponse

// SaveStatusResponse 保存状态响应
// Deprecated: 使用 dto.SaveStatusResponse 替代
type SaveStatusResponse = dto.SaveStatusResponse

// DocumentContentResponse 文档内容响应
// Deprecated: 使用 dto.DocumentContentResponse 替代
type DocumentContentResponse = dto.DocumentContentResponse

// ParagraphContent 段落内容
// Deprecated: 使用 dto.ParagraphContent 替代
type ParagraphContent = dto.ParagraphContent

// DocumentContentsResponse 文档段落内容响应
// Deprecated: 使用 dto.DocumentContentsResponse 替代
type DocumentContentsResponse = dto.DocumentContentsResponse

// ReplaceDocumentContentsResponse 替换文档段落响应
// Deprecated: 使用 dto.ReplaceDocumentContentsResponse 替代
type ReplaceDocumentContentsResponse = dto.ReplaceDocumentContentsResponse

// ReindexDocumentContentsResponse 重新编号响应
// Deprecated: 使用 dto.ReindexDocumentContentsResponse 替代
type ReindexDocumentContentsResponse = dto.ReindexDocumentContentsResponse
