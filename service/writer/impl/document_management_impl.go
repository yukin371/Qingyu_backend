package impl

import (
	"context"

	"Qingyu_backend/models/dto"
	"Qingyu_backend/models/writer"
	serviceWriter "Qingyu_backend/service/interfaces/writer"
	writerdocument "Qingyu_backend/service/writer/document"
	writerproject "Qingyu_backend/service/writer/project"
)

// DocumentManagementImpl 文档管理端口实现
type DocumentManagementImpl struct {
	documentService  *writerdocument.DocumentService
	projectService   *writerproject.ProjectService
	nodeService      *writerproject.NodeService
	versionService   *writerproject.VersionService
	autosaveService  *writerproject.AutoSaveService
	duplicateService *writerdocument.DuplicateService
	serviceName      string
	version          string
}

// NewDocumentManagementImpl 创建文档管理端口实现
func NewDocumentManagementImpl(
	documentService *writerdocument.DocumentService,
	projectService *writerproject.ProjectService,
	nodeService *writerproject.NodeService,
	versionService *writerproject.VersionService,
	autosaveService *writerproject.AutoSaveService,
	duplicateService *writerdocument.DuplicateService,
) serviceWriter.DocumentManagementPort {
	return &DocumentManagementImpl{
		documentService:  documentService,
		projectService:   projectService,
		nodeService:      nodeService,
		versionService:   versionService,
		autosaveService:  autosaveService,
		duplicateService: duplicateService,
		serviceName:      "DocumentManagementPort",
		version:          "1.0.0",
	}
}

// ============================================================================
// BaseService 生命周期方法实现
// ============================================================================

func (d *DocumentManagementImpl) Initialize(ctx context.Context) error {
	return d.documentService.Initialize(ctx)
}

func (d *DocumentManagementImpl) Health(ctx context.Context) error {
	return d.documentService.Health(ctx)
}

func (d *DocumentManagementImpl) Close(ctx context.Context) error {
	return d.documentService.Close(ctx)
}

func (d *DocumentManagementImpl) GetServiceName() string {
	return d.serviceName
}

func (d *DocumentManagementImpl) GetVersion() string {
	return d.version
}

// ============================================================================
// DocumentManagementPort 方法实现
// ============================================================================

// CreateDocument 创建文档
func (d *DocumentManagementImpl) CreateDocument(ctx context.Context, req *serviceWriter.CreateDocumentRequest) (*serviceWriter.CreateDocumentResponse, error) {
	// 转换请求类型
	documentReq := &writerdocument.CreateDocumentRequest{
		ProjectID:    req.ProjectID,
		ParentID:     req.ParentID,
		Title:        req.Title,
		Type:         req.Type,
		Order:        req.Order,
		CharacterIDs: req.CharacterIDs,
		LocationIDs:  req.LocationIDs,
		TimelineIDs:  req.TimelineIDs,
		Tags:         req.Tags,
		Notes:        req.Notes,
	}
	documentResp, err := d.documentService.CreateDocument(ctx, documentReq)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	return &serviceWriter.CreateDocumentResponse{
		DocumentID: documentResp.DocumentID,
		Title:      documentResp.Title,
		Type:       documentResp.Type,
		CreatedAt:  documentResp.CreatedAt,
	}, nil
}

// GetDocument 获取文档详情
func (d *DocumentManagementImpl) GetDocument(ctx context.Context, documentID string) (*writer.Document, error) {
	return d.documentService.GetDocument(ctx, documentID)
}

// GetDocumentTree 获取文档树
func (d *DocumentManagementImpl) GetDocumentTree(ctx context.Context, projectID string) (*serviceWriter.DocumentTreeResponse, error) {
	documentResp, err := d.documentService.GetDocumentTree(ctx, projectID)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	return &serviceWriter.DocumentTreeResponse{
		ProjectID: documentResp.ProjectID,
		Documents: d.convertDocumentTreeNodes(documentResp.Documents),
	}, nil
}

// convertDocumentTreeNodes 递归转换文档树节点
func (d *DocumentManagementImpl) convertDocumentTreeNodes(nodes []*writerdocument.DocumentTreeNode) []*dto.DocumentTreeNode {
	if nodes == nil {
		return nil
	}
	result := make([]*dto.DocumentTreeNode, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, &dto.DocumentTreeNode{
			Document: node.Document,
			Children: d.convertDocumentTreeNodes(node.Children),
		})
	}
	return result
}

// ListDocuments 获取文档列表
func (d *DocumentManagementImpl) ListDocuments(ctx context.Context, req *serviceWriter.ListDocumentsRequest) (*serviceWriter.ListDocumentsResponse, error) {
	// 转换请求类型
	documentReq := &writerdocument.ListDocumentsRequest{
		ProjectID: req.ProjectID,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}
	documentResp, err := d.documentService.ListDocuments(ctx, documentReq)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	return &serviceWriter.ListDocumentsResponse{
		Documents: documentResp.Documents,
		Total:     int(documentResp.Total),
		Page:      documentResp.Page,
		PageSize:  documentResp.PageSize,
	}, nil
}

// UpdateDocument 更新文档
func (d *DocumentManagementImpl) UpdateDocument(ctx context.Context, documentID string, req *serviceWriter.UpdateDocumentRequest) error {
	// 转换请求类型
	documentReq := &writerdocument.UpdateDocumentRequest{
		Title:        req.Title,
		Notes:        req.Notes,
		Status:       req.Status,
		CharacterIDs: req.CharacterIDs,
		LocationIDs:  req.LocationIDs,
		TimelineIDs:  req.TimelineIDs,
		Tags:         req.Tags,
	}
	return d.documentService.UpdateDocument(ctx, documentID, documentReq)
}

// DeleteDocument 删除文档
func (d *DocumentManagementImpl) DeleteDocument(ctx context.Context, documentID string) error {
	return d.documentService.DeleteDocument(ctx, documentID)
}

// MoveDocument 移动文档
func (d *DocumentManagementImpl) MoveDocument(ctx context.Context, req *serviceWriter.MoveDocumentRequest) error {
	// 转换请求类型
	documentReq := &writerdocument.MoveDocumentRequest{
		DocumentID:  req.DocumentID,
		NewParentID: req.NewParentID,
		Order:       req.Order,
	}
	return d.documentService.MoveDocument(ctx, documentReq)
}

// ReorderDocuments 重新排序文档
func (d *DocumentManagementImpl) ReorderDocuments(ctx context.Context, req *serviceWriter.ReorderDocumentsRequest) error {
	// 转换请求类型
	documentReq := &writerdocument.ReorderDocumentsRequest{
		ProjectID: req.ProjectID,
		ParentID:  req.ParentID,
		Orders:    req.Orders,
	}
	return d.documentService.ReorderDocuments(ctx, documentReq)
}

// AutoSaveDocument 自动保存文档
func (d *DocumentManagementImpl) AutoSaveDocument(ctx context.Context, req *serviceWriter.AutoSaveRequest) (*serviceWriter.AutoSaveResponse, error) {
	// 转换请求类型
	documentReq := &writerdocument.AutoSaveRequest{
		DocumentID:     req.DocumentID,
		Content:        req.Content,
		CurrentVersion: req.CurrentVersion,
		SaveType:       req.SaveType,
	}
	documentResp, err := d.documentService.AutoSaveDocument(ctx, documentReq)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	return &serviceWriter.AutoSaveResponse{
		Saved:       documentResp.Saved,
		NewVersion:  documentResp.NewVersion,
		WordCount:   documentResp.WordCount,
		SavedAt:     documentResp.SavedAt,
		HasConflict: documentResp.HasConflict,
	}, nil
}

// GetSaveStatus 获取保存状态
func (d *DocumentManagementImpl) GetSaveStatus(ctx context.Context, documentID string) (*serviceWriter.SaveStatusResponse, error) {
	documentResp, err := d.documentService.GetSaveStatus(ctx, documentID)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	return &serviceWriter.SaveStatusResponse{
		DocumentID:     documentResp.DocumentID,
		LastSavedAt:    documentResp.LastSavedAt,
		CurrentVersion: documentResp.CurrentVersion,
		IsSaving:       documentResp.IsSaving,
		WordCount:      documentResp.WordCount,
	}, nil
}

// GetDocumentContent 获取文档内容
func (d *DocumentManagementImpl) GetDocumentContent(ctx context.Context, documentID string) (*serviceWriter.DocumentContentResponse, error) {
	documentResp, err := d.documentService.GetDocumentContent(ctx, documentID)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	return &serviceWriter.DocumentContentResponse{
		DocumentID: documentResp.DocumentID,
		Content:    documentResp.Content,
		Version:    documentResp.Version,
		WordCount:  documentResp.WordCount,
		UpdatedAt:  documentResp.UpdatedAt,
	}, nil
}

// UpdateDocumentContent 更新文档内容
func (d *DocumentManagementImpl) UpdateDocumentContent(ctx context.Context, req *serviceWriter.UpdateContentRequest) error {
	// 转换请求类型
	documentReq := &writerdocument.UpdateContentRequest{
		DocumentID: req.DocumentID,
		Content:    req.Content,
		Version:    req.Version,
	}
	return d.documentService.UpdateDocumentContent(ctx, documentReq)
}

// DuplicateDocument 复制文档
func (d *DocumentManagementImpl) DuplicateDocument(ctx context.Context, documentID string, req *serviceWriter.DuplicateRequest) (*serviceWriter.DuplicateResponse, error) {
	// 转换请求类型 - Service API 和 Port API 的语义不同，需要做适配
	var targetParentID *string
	if req.TargetProject != "" {
		targetParentID = &req.TargetProject
	}
	duplicateReq := &writerdocument.DuplicateRequest{
		TargetParentID: targetParentID,
		Position:       "inner", // 默认放在内部
		CopyContent:    req.CopyContent,
	}
	documentResp, err := d.duplicateService.Duplicate(ctx, documentID, duplicateReq)
	if err != nil {
		return nil, err
	}
	// 从 Document 中提取信息构建响应
	newTitle := req.NewTitle
	if newTitle == "" {
		newTitle = documentResp.Title
	}
	return &serviceWriter.DuplicateResponse{
		NewDocumentID: documentResp.ID.Hex(),
		Title:         newTitle,
	}, nil
}

// ============================================================================
// 版本控制方法
// ============================================================================

// BumpVersionAndCreateRevision 创建新版本并记录修订
func (d *DocumentManagementImpl) BumpVersionAndCreateRevision(projectID, nodeID, authorID, message string) (*writer.FileRevision, error) {
	return d.versionService.BumpVersionAndCreateRevision(projectID, nodeID, authorID, message)
}

// UpdateContentWithVersion 使用乐观并发控制更新内容
func (d *DocumentManagementImpl) UpdateContentWithVersion(projectID, nodeID, authorID, message, newContent string, expectedVersion int) (*writer.FileRevision, error) {
	return d.versionService.UpdateContentWithVersion(projectID, nodeID, authorID, message, newContent, expectedVersion)
}

// RollbackToVersion 回滚到指定的历史版本
func (d *DocumentManagementImpl) RollbackToVersion(projectID, nodeID string, targetVersion int, authorID, message string) (*writer.FileRevision, error) {
	return d.versionService.RollbackToVersion(projectID, nodeID, targetVersion, authorID, message)
}

// ListRevisions 列表修订
func (d *DocumentManagementImpl) ListRevisions(ctx context.Context, projectID, nodeID string, limit, offset int64) ([]*writer.FileRevision, error) {
	return d.versionService.ListRevisions(ctx, projectID, nodeID, limit, offset)
}

// GetVersionHistory 获取版本历史
func (d *DocumentManagementImpl) GetVersionHistory(ctx context.Context, documentID string, page, pageSize int) (*serviceWriter.VersionHistoryResponse, error) {
	versionResp, err := d.versionService.GetVersionHistory(ctx, documentID, page, pageSize)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	versions := make([]*serviceWriter.VersionInfo, 0, len(versionResp.Versions))
	for _, v := range versionResp.Versions {
		versions = append(versions, &serviceWriter.VersionInfo{
			VersionID: v.VersionID,
			Version:   v.Version,
			Message:   v.Message,
			CreatedAt: v.CreatedAt,
			CreatedBy: v.CreatedBy,
			WordCount: v.WordCount,
		})
	}
	return &serviceWriter.VersionHistoryResponse{
		Versions: versions,
		Total:    versionResp.Total,
		Page:     versionResp.Page,
		PageSize: versionResp.PageSize,
	}, nil
}

// GetVersionDetail 获取特定版本
func (d *DocumentManagementImpl) GetVersionDetail(ctx context.Context, documentID, versionID string) (*serviceWriter.VersionDetail, error) {
	versionResp, err := d.versionService.GetVersion(ctx, documentID, versionID)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	return &serviceWriter.VersionDetail{
		VersionID:  versionResp.VersionID,
		DocumentID: versionResp.DocumentID,
		Version:    versionResp.Version,
		Content:    versionResp.Content,
		Message:    versionResp.Message,
		CreatedAt:  versionResp.CreatedAt,
		CreatedBy:  versionResp.CreatedBy,
		WordCount:  versionResp.WordCount,
	}, nil
}

// CompareVersions 比较两个版本
func (d *DocumentManagementImpl) CompareVersions(ctx context.Context, documentID, fromVersionID, toVersionID string) (*serviceWriter.VersionDiff, error) {
	versionResp, err := d.versionService.CompareVersions(ctx, documentID, fromVersionID, toVersionID)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	changes := make([]serviceWriter.ChangeItem, 0, len(versionResp.Changes))
	for _, c := range versionResp.Changes {
		changes = append(changes, serviceWriter.ChangeItem{
			Type:    c.Type,
			Line:    c.Line,
			Content: c.Content,
		})
	}
	return &serviceWriter.VersionDiff{
		FromVersion:  versionResp.FromVersion,
		ToVersion:    versionResp.ToVersion,
		Changes:      changes,
		AddedLines:   versionResp.AddedLines,
		DeletedLines: versionResp.DeletedLines,
	}, nil
}

// RestoreVersion 恢复到特定版本
func (d *DocumentManagementImpl) RestoreVersion(ctx context.Context, documentID, versionID string) error {
	return d.versionService.RestoreVersion(ctx, documentID, versionID)
}

// CreatePatch 提交一个候选补丁
func (d *DocumentManagementImpl) CreatePatch(projectID, nodeID string, baseVersion int, diffFormat, diffPayload, createdBy, message string) (*writer.FilePatch, error) {
	return d.versionService.CreatePatch(projectID, nodeID, baseVersion, diffFormat, diffPayload, createdBy, message)
}

// ApplyPatch 审核并应用补丁
func (d *DocumentManagementImpl) ApplyPatch(projectID, patchID, applierID string) (*writer.FileRevision, error) {
	return d.versionService.ApplyPatch(projectID, patchID, applierID)
}

// DetectConflicts 检测文件的版本冲突
func (d *DocumentManagementImpl) DetectConflicts(ctx context.Context, projectID, nodeID string, expectedVersion int) (*writer.ConflictInfo, error) {
	return d.versionService.DetectConflicts(ctx, projectID, nodeID, expectedVersion)
}

// CreateCommit 创建批量提交
func (d *DocumentManagementImpl) CreateCommit(ctx context.Context, projectID, authorID, message string, files []writer.CommitFile) (*writer.Commit, error) {
	return d.versionService.CreateCommit(ctx, projectID, authorID, message, files)
}

// ListCommits 查询提交历史
func (d *DocumentManagementImpl) ListCommits(ctx context.Context, projectID string, authorID string, limit, offset int64) ([]*writer.Commit, error) {
	return d.versionService.ListCommits(ctx, projectID, authorID, limit, offset)
}

// GetCommitDetails 获取提交详情
func (d *DocumentManagementImpl) GetCommitDetails(ctx context.Context, projectID, commitID string) (*writer.Commit, []*writer.FileRevision, error) {
	return d.versionService.GetCommitDetails(ctx, projectID, commitID)
}

// ============================================================================
// 自动保存方法
// ============================================================================

// StartAutoSave 启动文档自动保存
func (d *DocumentManagementImpl) StartAutoSave(documentID, projectID, nodeID, userID string) error {
	return d.autosaveService.StartAutoSave(documentID, projectID, nodeID, userID)
}

// StopAutoSave 停止文档自动保存
func (d *DocumentManagementImpl) StopAutoSave(documentID string) error {
	return d.autosaveService.StopAutoSave(documentID)
}

// GetAutoSaveStatus 获取自动保存状态
func (d *DocumentManagementImpl) GetAutoSaveStatus(documentID string) (isRunning bool, lastSaved *interface{}) {
	running, saved := d.autosaveService.GetStatus(documentID)
	if saved != nil {
		// 转换 time.Time 为字符串，使用 interface{} 类型
		formatted := saved.Format("2006-01-02T15:04:05Z07:00")
		result := interface{}(formatted)
		return running, &result
	}
	return running, nil
}

// SaveImmediately 立即执行一次保存（手动触发）
func (d *DocumentManagementImpl) SaveImmediately(documentID string) error {
	return d.autosaveService.SaveImmediately(documentID)
}
