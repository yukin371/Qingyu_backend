package content

import (
	"context"
	"fmt"

	"Qingyu_backend/models/dto"
	"Qingyu_backend/models/writer"
	writerService "Qingyu_backend/service/writer/document"
)

// DocumentAdapter 文档适配器
// 将现有DocumentService的功能适配到新的DocumentService接口
type DocumentAdapter struct {
	documentService *writerService.DocumentService
}

// NewDocumentAdapter 创建文档适配器
func NewDocumentAdapter(documentService *writerService.DocumentService) *DocumentAdapter {
	return &DocumentAdapter{
		documentService: documentService,
	}
}

// =========================
// 基础 CRUD 操作
// =========================

// CreateDocument 创建新文档
func (a *DocumentAdapter) CreateDocument(ctx context.Context, req *dto.CreateDocumentRequest) (*dto.DocumentResponse, error) {
	// 转换为旧版请求类型（已经统一到dto包）
	oldReq := &writerService.CreateDocumentRequest{
		ProjectID:    req.ProjectID,
		ParentID:     req.ParentID,
		Title:        req.Title,
		Type:         req.Type,
		Level:        req.Level,
		Order:        req.Order,
		CharacterIDs: req.CharacterIDs,
		LocationIDs:  req.LocationIDs,
		TimelineIDs:  req.TimelineIDs,
		Tags:         req.Tags,
		Notes:        req.Notes,
	}

	// 调用现有服务
	resp, err := a.documentService.CreateDocument(ctx, oldReq)
	if err != nil {
		return nil, fmt.Errorf("创建文档失败: %w", err)
	}

	// 转换为新的响应类型
	return &dto.DocumentResponse{
		DocumentID: resp.DocumentID,
		Title:      resp.Title,
		Type:       resp.Type,
		CreatedAt:  resp.CreatedAt,
	}, nil
}

// UpdateDocument 更新文档元数据
func (a *DocumentAdapter) UpdateDocument(ctx context.Context, id string, req *dto.UpdateDocumentRequest) (*dto.DocumentResponse, error) {
	// 转换为旧版请求类型
	oldReq := &writerService.UpdateDocumentRequest{
		Title:        req.Title,
		Notes:        req.Notes,
		Status:       req.Status,
		CharacterIDs: req.CharacterIDs,
		LocationIDs:  req.LocationIDs,
		TimelineIDs:  req.TimelineIDs,
		Tags:         req.Tags,
	}

	// 调用现有服务
	err := a.documentService.UpdateDocument(ctx, id, oldReq)
	if err != nil {
		return nil, fmt.Errorf("更新文档失败: %w", err)
	}

	// 获取更新后的文档信息
	doc, err := a.documentService.GetDocument(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取更新后的文档失败: %w", err)
	}

	// 转换为新的响应类型
	return a.convertToDocumentResponse(doc), nil
}

// GetDocument 获取文档详情
func (a *DocumentAdapter) GetDocument(ctx context.Context, id string) (*dto.DocumentResponse, error) {
	// 调用现有服务
	doc, err := a.documentService.GetDocument(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取文档失败: %w", err)
	}

	// 转换为新的响应类型
	return a.convertToDocumentResponse(doc), nil
}

// DeleteDocument 删除文档
func (a *DocumentAdapter) DeleteDocument(ctx context.Context, id string) error {
	// 调用现有服务
	err := a.documentService.DeleteDocument(ctx, id)
	if err != nil {
		return fmt.Errorf("删除文档失败: %w", err)
	}

	return nil
}

// ListDocuments 获取文档列表
func (a *DocumentAdapter) ListDocuments(ctx context.Context, req *dto.ListDocumentsRequest) (*dto.ListDocumentsResponse, error) {
	// 转换为旧版请求类型
	oldReq := &writerService.ListDocumentsRequest{
		ProjectID: req.ProjectID,
		ParentID:  req.ParentID,
		Page:      req.Page,
		PageSize:  req.PageSize,
		Status:    req.Status,
	}

	// 调用现有服务
	resp, err := a.documentService.ListDocuments(ctx, oldReq)
	if err != nil {
		return nil, fmt.Errorf("获取文档列表失败: %w", err)
	}

	// 转换文档列表为DTO格式
	documents := make([]*dto.DocumentResponse, 0, len(resp.Documents))
	for _, doc := range resp.Documents {
		documents = append(documents, a.convertToDocumentResponse(doc))
	}

	return &dto.ListDocumentsResponse{
		Documents: documents,
		Total:     resp.Total,
		Page:      resp.Page,
		PageSize:  resp.PageSize,
	}, nil
}

// =========================
// 文档操作
// =========================

// DuplicateDocument 复制文档
func (a *DocumentAdapter) DuplicateDocument(ctx context.Context, id string) (*dto.DocumentResponse, error) {
	// 创建默认的复制请求（复制内容）
	dupReq := &writerService.DuplicateRequest{
		TargetParentID: nil, // 使用默认父节点
		Position:       "inner",
		CopyContent:    true,
	}

	// 调用现有服务的复制功能
	resp, err := a.documentService.DuplicateDocument(ctx, id, dupReq)
	if err != nil {
		return nil, fmt.Errorf("复制文档失败: %w", err)
	}

	// 获取新创建的文档
	newDoc, err := a.documentService.GetDocument(ctx, resp.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("获取复制的文档失败: %w", err)
	}

	// 转换为新的响应类型
	return a.convertToDocumentResponse(newDoc), nil
}

// MoveDocument 移动文档
func (a *DocumentAdapter) MoveDocument(ctx context.Context, id string, newParentID string, order int) error {
	// 转换为旧版请求类型
	oldReq := &writerService.MoveDocumentRequest{
		DocumentID:  id,
		NewParentID: newParentID,
		Order:       order,
	}

	// 调用现有服务
	err := a.documentService.MoveDocument(ctx, oldReq)
	if err != nil {
		return fmt.Errorf("移动文档失败: %w", err)
	}

	return nil
}

// GetDocumentTree 获取文档树
func (a *DocumentAdapter) GetDocumentTree(ctx context.Context, projectID string) (*dto.DocumentTreeResponse, error) {
	// 调用现有服务
	resp, err := a.documentService.GetDocumentTree(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("获取文档树失败: %w", err)
	}

	// 转换文档树为DTO格式
	documents := make([]*dto.DocumentTreeNode, 0, len(resp.Documents))
	for _, docNode := range resp.Documents {
		documents = append(documents, a.convertDocumentTreeNode(docNode))
	}

	return &dto.DocumentTreeResponse{
		ProjectID: projectID,
		Documents: documents,
	}, nil
}

// =========================
// 内容管理
// =========================

// GetDocumentContent 获取文档内容
func (a *DocumentAdapter) GetDocumentContent(ctx context.Context, documentID string) (*dto.DocumentContentResponse, error) {
	// 调用现有服务
	resp, err := a.documentService.GetDocumentContent(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("获取文档内容失败: %w", err)
	}

	// 已经是DTO类型，直接返回
	return resp, nil
}

// UpdateDocumentContent 更新文档内容
func (a *DocumentAdapter) UpdateDocumentContent(ctx context.Context, req *dto.UpdateContentRequest) error {
	// 转换为旧版请求类型（字段名相同）
	oldReq := &writerService.UpdateContentRequest{
		DocumentID: req.DocumentID,
		Content:    req.Content,
		Version:    req.Version,
	}

	// 调用现有服务
	err := a.documentService.UpdateDocumentContent(ctx, oldReq)
	if err != nil {
		return fmt.Errorf("更新文档内容失败: %w", err)
	}

	return nil
}

// AutoSaveDocument 自动保存
func (a *DocumentAdapter) AutoSaveDocument(ctx context.Context, req *dto.AutoSaveRequest) (*dto.AutoSaveResponse, error) {
	// 转换为旧版请求类型
	oldReq := &writerService.AutoSaveRequest{
		DocumentID:     req.DocumentID,
		Content:        req.Content,
		CurrentVersion: req.CurrentVersion,
		SaveType:       req.SaveType,
	}

	// 调用现有服务
	resp, err := a.documentService.AutoSaveDocument(ctx, oldReq)
	if err != nil {
		return nil, fmt.Errorf("自动保存失败: %w", err)
	}

	// 已经是DTO类型，直接返回
	return resp, nil
}

// =========================
// 版本控制
// =========================

// GetVersionHistory 获取版本历史
func (a *DocumentAdapter) GetVersionHistory(ctx context.Context, documentID string, page, pageSize int) (*dto.VersionHistoryResponse, error) {
	// 现有服务暂未实现完整的版本控制功能
	// 返回空列表或待实现错误
	return nil, fmt.Errorf("版本历史功能暂未实现，需要扩展DocumentService")
}

// GetVersionDetail 获取版本详情
func (a *DocumentAdapter) GetVersionDetail(ctx context.Context, documentID, versionID string) (*dto.VersionDetail, error) {
	// 现有服务暂未实现完整的版本控制功能
	return nil, fmt.Errorf("版本详情功能暂未实现，需要扩展DocumentService")
}

// CompareVersions 比较版本
func (a *DocumentAdapter) CompareVersions(ctx context.Context, documentID, fromVersionID, toVersionID string) (*dto.VersionDiff, error) {
	// 现有服务暂未实现完整的版本控制功能
	return nil, fmt.Errorf("版本比较功能暂未实现，需要扩展DocumentService")
}

// RestoreVersion 恢复版本
func (a *DocumentAdapter) RestoreVersion(ctx context.Context, documentID, versionID string) error {
	// 现有服务暂未实现完整的版本控制功能
	return fmt.Errorf("版本恢复功能暂未实现，需要扩展DocumentService")
}

// =========================
// 批量操作
// =========================

// BatchUpdateDocuments 批量更新文档
func (a *DocumentAdapter) BatchUpdateDocuments(ctx context.Context, documentIDs []string, updates map[string]interface{}) error {
	// 现有服务暂未实现批量更新功能
	return fmt.Errorf("批量更新功能暂未实现，需要扩展DocumentService")
}

// BatchDeleteDocuments 批量删除文档
func (a *DocumentAdapter) BatchDeleteDocuments(ctx context.Context, documentIDs []string) error {
	// 现有服务暂未实现批量删除功能
	// 可以逐个调用删除方法
	for _, id := range documentIDs {
		if err := a.documentService.DeleteDocument(ctx, id); err != nil {
			return fmt.Errorf("批量删除文档失败（ID: %s）: %w", id, err)
		}
	}
	return nil
}

// ReorderDocuments 重新排序文档
func (a *DocumentAdapter) ReorderDocuments(ctx context.Context, req *dto.ReorderDocumentsRequest) error {
	// 转换为旧版请求类型
	oldReq := &writerService.ReorderDocumentsRequest{
		ProjectID: req.ProjectID,
		ParentID:  req.ParentID,
		Orders:    req.Orders,
	}

	// 调用现有服务
	err := a.documentService.ReorderDocuments(ctx, oldReq)
	if err != nil {
		return fmt.Errorf("重新排序文档失败: %w", err)
	}

	return nil
}

// =========================
// 私有辅助方法
// =========================

// convertToDocumentResponse 将Document模型转换为DocumentResponse DTO
func (a *DocumentAdapter) convertToDocumentResponse(doc *writer.Document) *dto.DocumentResponse {
	if doc == nil {
		return nil
	}

	// 转换ObjectID数组为字符串数组
	characterIDs := make([]string, 0, len(doc.CharacterIDs))
	for _, id := range doc.CharacterIDs {
		characterIDs = append(characterIDs, id.Hex())
	}

	locationIDs := make([]string, 0, len(doc.LocationIDs))
	for _, id := range doc.LocationIDs {
		locationIDs = append(locationIDs, id.Hex())
	}

	timelineIDs := make([]string, 0, len(doc.TimelineIDs))
	for _, id := range doc.TimelineIDs {
		timelineIDs = append(timelineIDs, id.Hex())
	}

	// 转换ParentID（如果是零值则不设置）
	parentID := ""
	if !doc.ParentID.IsZero() {
		parentID = doc.ParentID.Hex()
	}

	return &dto.DocumentResponse{
		DocumentID:    doc.ID.Hex(),
		ProjectID:     doc.ProjectID.Hex(),
		ParentID:      parentID,
		Title:         doc.Title,
		Type:          doc.Type,
		Level:         doc.Level,
		Order:         doc.Order,
		Status:        doc.Status,
		CharacterIDs:  characterIDs,
		LocationIDs:   locationIDs,
		TimelineIDs:   timelineIDs,
		Tags:          doc.Tags,
		Notes:         doc.Notes,
		CreatedBy:     "", // 需要从doc获取，如果Document模型有CreatedBy字段
		CreatedAt:     doc.CreatedAt,
		UpdatedAt:     doc.UpdatedAt,
		WordCount:     doc.WordCount,
		Version:       1, // 默认版本，实际应从DocumentContent获取
	}
}

// convertDocumentTreeNode 递归转换文档树节点
func (a *DocumentAdapter) convertDocumentTreeNode(node *writerService.DocumentTreeNode) *dto.DocumentTreeNode {
	if node == nil {
		return nil
	}

	// 转换当前节点的Document为DocumentResponse
	docResponse := a.convertToDocumentResponse(node.Document)

	// 转换当前节点
	convertedNode := &dto.DocumentTreeNode{
		Document: docResponse,
		Children: make([]*dto.DocumentTreeNode, 0, len(node.Children)),
	}

	// 递归转换子节点
	for _, child := range node.Children {
		convertedNode.Children = append(convertedNode.Children, a.convertDocumentTreeNode(child))
	}

	return convertedNode
}
