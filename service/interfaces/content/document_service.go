package content

import (
	"context"

	"Qingyu_backend/models/dto"
)

// ============================================================================
// 文档服务接口详细定义
// ============================================================================

// DocumentService 文档管理服务接口
// 提供文档的完整生命周期管理功能
type DocumentService interface {
	// === 基础 CRUD 操作 ===

	// CreateDocument 创建新文档
	// 创建章节、场景或笔记类型的文档节点
	CreateDocument(ctx context.Context, req *dto.CreateDocumentRequest) (*dto.DocumentResponse, error)

	// UpdateDocument 更新文档元数据
	// 更新文档标题、状态、标签、关联的角色/地点/时间线等
	UpdateDocument(ctx context.Context, id string, req *dto.UpdateDocumentRequest) (*dto.DocumentResponse, error)

	// GetDocument 获取文档详情
	// 返回文档的完整信息，包括元数据但不含内容
	GetDocument(ctx context.Context, id string) (*dto.DocumentResponse, error)

	// DeleteDocument 删除文档
	// 执行软删除，将文档标记为已删除状态
	DeleteDocument(ctx context.Context, id string) error

	// ListDocuments 获取文档列表
	// 支持按父节点、状态分页查询
	ListDocuments(ctx context.Context, req *dto.ListDocumentsRequest) (*dto.ListDocumentsResponse, error)

	// === 文档操作 ===

	// DuplicateDocument 复制文档
	// 创建文档的副本，可选择是否复制内容
	DuplicateDocument(ctx context.Context, id string) (*dto.DocumentResponse, error)

	// MoveDocument 移动文档
	// 将文档移动到新的父节点下，并更新排序
	MoveDocument(ctx context.Context, id string, newParentID string, order int) error

	// GetDocumentTree 获取文档树
	// 返回项目的完整文档树形结构
	GetDocumentTree(ctx context.Context, projectID string) (*dto.DocumentTreeResponse, error)

	// === 内容管理 ===

	// GetDocumentContent 获取文档内容
	// 返回文档的实际文本内容和版本信息
	GetDocumentContent(ctx context.Context, documentID string) (*dto.DocumentContentResponse, error)

	// UpdateDocumentContent 更新文档内容
	// 更新文档的文本内容，支持版本冲突检测
	UpdateDocumentContent(ctx context.Context, req *dto.UpdateContentRequest) error

	// AutoSaveDocument 自动保存
	// 定时自动保存文档内容，返回保存结果和版本信息
	AutoSaveDocument(ctx context.Context, req *dto.AutoSaveRequest) (*dto.AutoSaveResponse, error)

	// === 版本控制 ===

	// GetVersionHistory 获取版本历史
	// 返回文档的版本历史记录，支持分页
	GetVersionHistory(ctx context.Context, documentID string, page, pageSize int) (*dto.VersionHistoryResponse, error)

	// GetVersionDetail 获取版本详情
	// 返回指定版本的完整内容和元数据
	GetVersionDetail(ctx context.Context, documentID, versionID string) (*dto.VersionDetail, error)

	// CompareVersions 比较版本
	// 返回两个版本之间的差异对比
	CompareVersions(ctx context.Context, documentID, fromVersionID, toVersionID string) (*dto.VersionDiff, error)

	// RestoreVersion 恢复版本
	// 将文档内容恢复到指定版本
	RestoreVersion(ctx context.Context, documentID, versionID string) error

	// === 批量操作 ===

	// BatchUpdateDocuments 批量更新文档
	// 批量更新多个文档的属性
	BatchUpdateDocuments(ctx context.Context, documentIDs []string, updates map[string]interface{}) error

	// BatchDeleteDocuments 批量删除文档
	// 批量软删除多个文档
	BatchDeleteDocuments(ctx context.Context, documentIDs []string) error

	// ReorderDocuments 重新排序文档
	// 批量更新同父节点下文档的排序
	ReorderDocuments(ctx context.Context, req *dto.ReorderDocumentsRequest) error
}
