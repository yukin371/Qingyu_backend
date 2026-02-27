package document

import (
	"Qingyu_backend/models/writer"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	pkgErrors "Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	serviceBase "Qingyu_backend/service/base"
)

// DocumentService 文档服务
type DocumentService struct {
	documentRepo        writerRepo.DocumentRepository
	documentContentRepo writerRepo.DocumentContentRepository
	projectRepo         writerRepo.ProjectRepository
	eventBus            serviceBase.EventBus
	duplicateService    *DuplicateService
	authHelper          *AuthHelper
	serviceName         string
	version             string
}

// NewDocumentService 创建文档服务
func NewDocumentService(
	documentRepo writerRepo.DocumentRepository,
	documentContentRepo writerRepo.DocumentContentRepository,
	projectRepo writerRepo.ProjectRepository,
	eventBus serviceBase.EventBus,
) *DocumentService {
	duplicateSvc := NewDuplicateService(documentRepo, documentContentRepo)
	serviceName := "DocumentService"

	return &DocumentService{
		documentRepo:        documentRepo,
		documentContentRepo: documentContentRepo,
		projectRepo:         projectRepo,
		eventBus:            eventBus,
		duplicateService:    duplicateSvc,
		authHelper:          NewAuthHelper(projectRepo, documentRepo, serviceName),
		serviceName:         serviceName,
		version:             "1.0.0",
	}
}

// CreateDocument 创建文档
func (s *DocumentService) CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*CreateDocumentResponse, error) {
	// 1. 参数验证
	if err := s.validateCreateDocumentRequest(req); err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "参数验证失败", err.Error(), err)
	}

	// 2. 验证项目编辑权限
	userID, project, err := s.authHelper.VerifyProjectEdit(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}
	_ = userID // userID 已通过权限验证，后续可能使用

	// 3. 验证父文档和层级
	var level int
	if req.ParentID != "" {
		parent, err := s.documentRepo.GetByID(ctx, req.ParentID)
		if err != nil {
			return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询父文档失败", "", err)
		}

		if parent == nil {
			return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "父文档不存在", "", nil)
		}

		if !parent.CanHaveChildren(project.WritingType) {
			return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorBusiness, "该文档不能添加子文档（达到最大层级）", "", nil)
		}

		level = parent.GetNextLevel()
	}

	// 4. 创建文档对象（使用base mixins）
	doc := &writer.Document{}
	
	// 转换 ProjectID string -> ObjectID
	projectID, err := primitive.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "无效的项目ID", "", err)
	}
	doc.ProjectID = projectID
	doc.Title = req.Title
	
	// 转换 ParentID string -> ObjectID（如果有）
	var parentID primitive.ObjectID
	if req.ParentID != "" {
		parentID, err = primitive.ObjectIDFromHex(req.ParentID)
		if err != nil {
			return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "无效的父文档ID", "", err)
		}
	}
	doc.ParentID = parentID
	
	doc.Type = req.Type // DocumentType现在是string类型
	doc.Level = level
	doc.Order = req.Order
	doc.Status = "planned"
	doc.WordCount = 0
	
	// 转换 CharacterIDs []string -> []ObjectID
	doc.CharacterIDs = make([]primitive.ObjectID, 0, len(req.CharacterIDs))
	for _, id := range req.CharacterIDs {
		charID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "无效的角色ID: "+id, "", err)
		}
		doc.CharacterIDs = append(doc.CharacterIDs, charID)
	}
	
	// 转换 LocationIDs []string -> []ObjectID
	doc.LocationIDs = make([]primitive.ObjectID, 0, len(req.LocationIDs))
	for _, id := range req.LocationIDs {
		locID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "无效的地点ID: "+id, "", err)
		}
		doc.LocationIDs = append(doc.LocationIDs, locID)
	}
	
	// 转换 TimelineIDs []string -> []ObjectID
	doc.TimelineIDs = make([]primitive.ObjectID, 0, len(req.TimelineIDs))
	for _, id := range req.TimelineIDs {
		timeID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "无效的时间线ID: "+id, "", err)
		}
		doc.TimelineIDs = append(doc.TimelineIDs, timeID)
	}
	
	doc.Tags = req.Tags
	doc.Notes = req.Notes

	// 5. 保存文档
	if err := s.documentRepo.Create(ctx, doc); err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "创建文档失败", "", err)
	}

	// 6. 更新项目统计（异步）
	go s.updateProjectStatistics(context.Background(), req.ProjectID)

	// 7. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &serviceBase.BaseEvent{
			EventType: "document.created",
			EventData: map[string]interface{}{
				"document_id": doc.ID.Hex(),
				"project_id":  doc.ProjectID.Hex(),
				"title":       doc.Title,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return &CreateDocumentResponse{
		DocumentID: doc.ID.Hex(),
		Title:      doc.Title,
		Type:       string(doc.Type),
		CreatedAt:  doc.CreatedAt,
	}, nil
}

// GetDocument 获取文档详情
func (s *DocumentService) GetDocument(ctx context.Context, documentID string) (*writer.Document, error) {
	// 1. 验证文档查看权限
	_, doc, _, err := s.authHelper.VerifyDocumentView(ctx, documentID)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// GetDocumentTree 获取文档树
func (s *DocumentService) GetDocumentTree(ctx context.Context, projectID string) (*DocumentTreeResponse, error) {
	// 1. 验证项目查看权限
	_, _, err := s.authHelper.VerifyProjectView(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 2. 获取所有文档
	documents, err := s.documentRepo.GetByProjectID(ctx, projectID, 1000, 0)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询文档树失败", "", err)
	}

	// 3. 构建树形结构
	tree := s.buildDocumentTree(documents)

	return &DocumentTreeResponse{
		ProjectID: projectID,
		Documents: tree,
	}, nil
}

// UpdateDocument 更新文档
func (s *DocumentService) UpdateDocument(ctx context.Context, documentID string, req *UpdateDocumentRequest) error {
	// 1. 验证文档编辑权限
	_, doc, _, err := s.authHelper.VerifyDocumentEdit(ctx, documentID)
	if err != nil {
		return err
	}

	// 2. 构建更新数据
	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.CharacterIDs != nil {
		updates["character_ids"] = req.CharacterIDs
	}
	if req.LocationIDs != nil {
		updates["location_ids"] = req.LocationIDs
	}
	if req.TimelineIDs != nil {
		updates["timeline_ids"] = req.TimelineIDs
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}

	// 3. 更新文档
	if err := s.documentRepo.UpdateByProject(ctx, documentID, doc.ProjectID.Hex(), updates); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "更新文档失败", "", err)
	}

	// 4. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &serviceBase.BaseEvent{
			EventType: "document.updated",
			EventData: map[string]interface{}{
				"document_id": documentID,
				"project_id":  doc.ProjectID.Hex(),
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// DeleteDocument 删除文档
func (s *DocumentService) DeleteDocument(ctx context.Context, documentID string) error {
	// 1. 验证文档编辑权限
	_, doc, _, err := s.authHelper.VerifyDocumentEdit(ctx, documentID)
	if err != nil {
		return err
	}

	// 2. 软删除文档
	if err := s.documentRepo.SoftDelete(ctx, documentID, doc.ProjectID.Hex()); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "删除文档失败", "", err)
	}

	// 3. 更新项目统计（异步）
	go s.updateProjectStatistics(context.Background(), doc.ProjectID.Hex())

	// 4. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &serviceBase.BaseEvent{
			EventType: "document.deleted",
			EventData: map[string]interface{}{
				"document_id": documentID,
				"project_id":  doc.ProjectID.Hex(),
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// 私有方法
func (s *DocumentService) buildDocumentTree(documents []*writer.Document) []*DocumentTreeNode {
	// 构建树形结构
	nodeMap := make(map[string]*DocumentTreeNode)
	var rootNodes []*DocumentTreeNode

	// 第一遍遍历：创建所有节点
	for _, doc := range documents {
		node := &DocumentTreeNode{
			Document: doc,
			Children: []*DocumentTreeNode{},
		}
		nodeMap[doc.ID.Hex()] = node
	}

	// 第二遍遍历：建立父子关系
	for _, doc := range documents {
		node := nodeMap[doc.ID.Hex()]
		if doc.ParentID.IsZero() {
			rootNodes = append(rootNodes, node)
		} else {
			if parent, exists := nodeMap[doc.ParentID.Hex()]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return rootNodes
}

func (s *DocumentService) validateCreateDocumentRequest(req *CreateDocumentRequest) error {
	if req.ProjectID == "" {
		return fmt.Errorf("项目ID不能为空")
	}
	if req.Title == "" {
		return fmt.Errorf("文档标题不能为空")
	}
	if len(req.Title) > 200 {
		return fmt.Errorf("文档标题不能超过200字符")
	}
	if req.Type == "" {
		return fmt.Errorf("文档类型不能为空")
	}

	// 注意：文档类型的具体验证需要在服务层进行，因为需要项目的writing_type
	// 这里只做基础验证

	return nil
}

func (s *DocumentService) updateProjectStatistics(ctx context.Context, projectID string) {
	// 查询项目所有文档
	documents, err := s.documentRepo.GetByProjectID(ctx, projectID, 10000, 0)
	if err != nil {
		return
	}

	// 统计总字数和文档数
	totalWords := 0
	documentCount := len(documents)
	chapterCount := 0

	for _, doc := range documents {
		totalWords += doc.WordCount
		if doc.Type == writer.TypeChapter {
			chapterCount++
		}
	}

	// 更新项目统计
	stats := &writer.ProjectStats{
		TotalWords:    totalWords,
		ChapterCount:  chapterCount,
		DocumentCount: documentCount,
		LastUpdateAt:  time.Now(),
	}

	updates := map[string]interface{}{
		"statistics": stats,
	}

	s.projectRepo.Update(ctx, projectID, updates)
}

// BaseService接口实现
func (s *DocumentService) Initialize(ctx context.Context) error {
	return nil
}

func (s *DocumentService) Health(ctx context.Context) error {
	return s.documentRepo.Health(ctx)
}

func (s *DocumentService) Close(ctx context.Context) error {
	return nil
}

func (s *DocumentService) GetServiceName() string {
	return s.serviceName
}

func (s *DocumentService) GetVersion() string {
	return s.version
}

// ListDocuments 获取文档列表
func (s *DocumentService) ListDocuments(ctx context.Context, req *ListDocumentsRequest) (*ListDocumentsResponse, error) {
	// 1. 参数验证
	if req.ProjectID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "项目ID不能为空", "", nil)
	}

	// 2. 解析分页参数
	page := 1
	pageSize := 20
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	// 3. 查询文档列表
	offset := (page - 1) * pageSize
	documents, err := s.documentRepo.GetByProjectID(ctx, req.ProjectID, int64(pageSize), int64(offset))
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询文档列表失败", "", err)
	}

	// 4. 统计总数
	allDocs, err := s.documentRepo.GetByProjectID(ctx, req.ProjectID, 10000, 0)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "统计文档总数失败", "", err)
	}

	return &ListDocumentsResponse{
		Documents: documents,
		Total:     len(allDocs),
		Page:      page,
		PageSize:  pageSize,
	}, nil
}

// MoveDocument 移动文档
func (s *DocumentService) MoveDocument(ctx context.Context, req *MoveDocumentRequest) error {
	// 1. 参数验证
	if req.DocumentID == "" {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "文档ID不能为空", "", nil)
	}

	// 2. 验证文档编辑权限
	_, doc, project, err := s.authHelper.VerifyDocumentEdit(ctx, req.DocumentID)
	if err != nil {
		return err
	}

	// 3. 验证新父节点
	var level int
	var newParentID primitive.ObjectID
	if req.NewParentID != "" {
		parent, err := s.documentRepo.GetByID(ctx, req.NewParentID)
		if err != nil {
			return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询父文档失败", "", err)
		}

		if parent == nil {
			return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "目标父文档不存在", "", nil)
		}

		if !parent.CanHaveChildren(project.WritingType) {
			return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorBusiness, "目标节点不能添加子文档（达到最大层级）", "", nil)
		}

		level = parent.GetNextLevel()
		newParentID = parent.ID
	}

	// 4. 更新文档
	updates := map[string]interface{}{
		"parent_id": newParentID,
		"level":     level,
		"order":     req.Order,
	}

	if err := s.documentRepo.Update(ctx, req.DocumentID, updates); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "移动文档失败", "", err)
	}

	// 5. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &serviceBase.BaseEvent{
			EventType: "document.moved",
			EventData: map[string]interface{}{
				"document_id":   req.DocumentID,
				"new_parent_id": newParentID.Hex(),
				"project_id":    doc.ProjectID.Hex(),
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// ReorderDocuments 重新排序文档
func (s *DocumentService) ReorderDocuments(ctx context.Context, req *ReorderDocumentsRequest) error {
	// 1. 参数验证
	if req.ProjectID == "" {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "项目ID不能为空", "", nil)
	}

	if len(req.Orders) == 0 {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "排序列表不能为空", "", nil)
	}

	// 2. 验证项目编辑权限
	_, _, err := s.authHelper.VerifyProjectEdit(ctx, req.ProjectID)
	if err != nil {
		return err
	}

	// 3. 批量更新排序
	for docID, order := range req.Orders {
		updates := map[string]interface{}{
			"order": order,
		}
		if err := s.documentRepo.Update(ctx, docID, updates); err != nil {
			return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, fmt.Sprintf("更新文档 %s 排序失败", docID), "", err)
		}
	}

	// 4. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &serviceBase.BaseEvent{
			EventType: "documents.reordered",
			EventData: map[string]interface{}{
				"project_id": req.ProjectID,
				"parent_id":  req.ParentID,
				"count":      len(req.Orders),
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// AutoSaveDocument 自动保存文档
func (s *DocumentService) AutoSaveDocument(ctx context.Context, req *AutoSaveRequest) (*AutoSaveResponse, error) {
	// 1. 参数验证
	if req.DocumentID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "文档ID不能为空", "", nil)
	}
	if req.Content == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "内容不能为空", "", nil)
	}

	// 2. 验证文档编辑权限
	_, _, _, err := s.authHelper.VerifyDocumentEdit(ctx, req.DocumentID)
	if err != nil {
		return nil, err
	}

	// 3. 获取或创建DocumentContent
	content, err := s.documentContentRepo.GetByDocumentID(ctx, req.DocumentID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询文档内容失败", "", err)
	}

	var newVersion int
	hasConflict := false

	if content == nil {
		// 首次保存，创建新DocumentContent
		documentID, err := primitive.ObjectIDFromHex(req.DocumentID)
		if err != nil {
			return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "无效的文档ID", "", err)
		}
		newContent := &writer.DocumentContent{
			DocumentID:  documentID,
			Content:     req.Content,
			ContentType: "markdown", // 默认使用markdown格式
			Version:     1,
			WordCount:   len([]rune(req.Content)),
			CharCount:   len(req.Content),
		}
		if err := s.documentContentRepo.Create(ctx, newContent); err != nil {
			return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "创建文档内容失败", "", err)
		}
		newVersion = 1
	} else {
		// 更新现有内容（带版本检测）
		expectedVersion := req.CurrentVersion
		if expectedVersion == 0 {
			// 如果客户端未提供版本号，使用当前版本
			expectedVersion = content.Version
		}

		err := s.documentContentRepo.UpdateWithVersion(
			ctx,
			req.DocumentID,
			req.Content,
			expectedVersion,
		)

		if err != nil {
			if err.Error() == "版本冲突，请重新获取最新内容" {
				// 版本冲突
				return &AutoSaveResponse{
					Saved:       false,
					NewVersion:  content.Version,
					WordCount:   len([]rune(req.Content)),
					SavedAt:     time.Now(),
					HasConflict: true,
				}, nil
			}
			return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "保存文档内容失败", "", err)
		}

		// 更新成功，版本号+1
		newVersion = expectedVersion + 1
	}

	// 4. 计算字数并更新Document元数据
	wordCount := len([]rune(req.Content))
	updates := map[string]interface{}{
		"word_count": wordCount,
		"updated_at": time.Now(),
	}

	if err := s.documentRepo.Update(ctx, req.DocumentID, updates); err != nil {
		// 内容已保存，但元数据更新失败，记录错误但不返回失败
		// 下次获取文档时会自动同步
		fmt.Printf("警告：更新文档元数据失败: %v\n", err)
	}

	// 5. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &serviceBase.BaseEvent{
			EventType: "document.autosaved",
			EventData: map[string]interface{}{
				"document_id": req.DocumentID,
				"word_count":  wordCount,
				"save_type":   req.SaveType,
				"version":     newVersion,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	// 6. 返回响应
	return &AutoSaveResponse{
		Saved:       true,
		NewVersion:  newVersion,
		WordCount:   wordCount,
		SavedAt:     time.Now(),
		HasConflict: hasConflict,
	}, nil
}

// GetSaveStatus 获取保存状态
func (s *DocumentService) GetSaveStatus(ctx context.Context, documentID string) (*SaveStatusResponse, error) {
	// 1. 验证参数
	if documentID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "文档ID不能为空", "", nil)
	}

	// 2. 获取文档
	doc, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询文档失败", "", err)
	}

	if doc == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "文档不存在", "", nil)
	}

	// 3. 获取文档版本号
	// TODO:注意：完整的版本控制系统应该使用 VersionService
	// 当前实现：基于文档创建和更新时间差来估算版本号
	currentVersion := s.calculateDocumentVersion(doc)

	// 4. 返回保存状态
	return &SaveStatusResponse{
		DocumentID:     documentID,
		LastSavedAt:    doc.UpdatedAt,
		CurrentVersion: currentVersion,
		IsSaving:       false,
		WordCount:      doc.WordCount,
	}, nil
}

// calculateDocumentVersion 计算文档版本号
// 注意：这是简化实现。生产环境应该使用 VersionService 管理真实的版本控制
func (s *DocumentService) calculateDocumentVersion(doc *writer.Document) int {
	// 方案1：基于时间差估算版本号（每次更新视为一个新版本）
	// 假设平均每5分钟保存一次
	timeDiff := doc.UpdatedAt.Sub(doc.CreatedAt)
	estimatedVersion := int(timeDiff.Minutes()/5) + 1

	// 限制版本号范围（1-999）
	if estimatedVersion < 1 {
		estimatedVersion = 1
	}
	if estimatedVersion > 999 {
		estimatedVersion = 999
	}

	return estimatedVersion

	// TODO(Production): 使用 VersionService 获取真实版本号
	// if s.versionService != nil {
	// 	version, err := s.versionService.GetCurrentVersion(ctx, doc.ID)
	// 	if err == nil && version > 0 {
	// 		return version
	// 	}
	// }
}

// GetDocumentContent 获取文档内容
func (s *DocumentService) GetDocumentContent(ctx context.Context, documentID string) (*DocumentContentResponse, error) {
	// 1. 验证参数
	if documentID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "文档ID不能为空", "", nil)
	}

	// 2. 验证文档查看权限
	_, doc, _, err := s.authHelper.VerifyDocumentView(ctx, documentID)
	if err != nil {
		return nil, err
	}

	// 3. 获取内容
	content, err := s.documentContentRepo.GetByDocumentID(ctx, documentID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询文档内容失败", "", err)
	}

	// 4. 构建响应
	actualContent := ""
	version := 1
	if content != nil {
		actualContent = content.Content
		version = content.Version
	}

	// 5. 返回内容响应
	return &DocumentContentResponse{
		DocumentID: documentID,
		Content:    actualContent,
		Version:    version,
		WordCount:  doc.WordCount,
		UpdatedAt:  doc.UpdatedAt,
	}, nil
}

// UpdateDocumentContent 更新文档内容
func (s *DocumentService) UpdateDocumentContent(ctx context.Context, req *UpdateContentRequest) error {
	// 1. 参数验证
	if req.DocumentID == "" {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "文档ID不能为空", "", nil)
	}

	// 2. 验证文档编辑权限
	_, _, _, err := s.authHelper.VerifyDocumentEdit(ctx, req.DocumentID)
	if err != nil {
		return err
	}

	// 3. 获取现有DocumentContent检查版本
	existingContent, err := s.documentContentRepo.GetByDocumentID(ctx, req.DocumentID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询文档内容失败", "", err)
	}

	// 4. 保存内容
	if existingContent == nil {
		// 首次保存，创建新内容
		documentID, err := primitive.ObjectIDFromHex(req.DocumentID)
		if err != nil {
			return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "无效的文档ID", "", err)
		}
		newContent := &writer.DocumentContent{
			DocumentID:  documentID,
			Content:     req.Content,
			ContentType: "markdown", // 默认使用markdown格式
			Version:     1,
			WordCount:   len([]rune(req.Content)),
			CharCount:   len(req.Content),
		}
		if err := s.documentContentRepo.Create(ctx, newContent); err != nil {
			return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "创建文档内容失败", "", err)
		}
	} else {
		// 更新现有内容（带版本检测）
		expectedVersion := req.Version
		if expectedVersion == 0 {
			// 如果未提供版本号，使用当前版本
			expectedVersion = existingContent.Version
		}

		err := s.documentContentRepo.UpdateWithVersion(
			ctx,
			req.DocumentID,
			req.Content,
			expectedVersion,
		)

		if err != nil {
			if err.Error() == "版本冲突，请重新获取最新内容" {
				return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorBusiness, "版本冲突，请刷新后重试", "", err)
			}
			return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "保存内容失败", "", err)
		}
	}

	// 5. 计算字数并更新Document元数据
	wordCount := len([]rune(req.Content))
	updates := map[string]interface{}{
		"word_count": wordCount,
		"updated_at": time.Now(),
	}

	if err := s.documentRepo.Update(ctx, req.DocumentID, updates); err != nil {
		// 内容已保存，但元数据更新失败，记录错误但不返回失败
		fmt.Printf("警告：更新文档元数据失败: %v\n", err)
	}

	// 6. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &serviceBase.BaseEvent{
			EventType: "document.content_updated",
			EventData: map[string]interface{}{
				"document_id": req.DocumentID,
				"word_count":  wordCount,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// DuplicateDocument 复制文档
func (s *DocumentService) DuplicateDocument(ctx context.Context, documentID string, req *DuplicateRequest) (*DuplicateResponse, error) {
	// 1. 参数验证
	if documentID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "文档ID不能为空", "", nil)
	}

	if req == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "请求参数不能为空", "", nil)
	}

	// 2. 验证文档编辑权限
	_, sourceDoc, _, err := s.authHelper.VerifyDocumentEdit(ctx, documentID)
	if err != nil {
		return nil, err
	}

	// 3. 调用DuplicateService执行复制
	newDoc, err := s.duplicateService.Duplicate(ctx, documentID, req)
	if err != nil {
		return nil, err
	}

	// 4. 更新项目统计（异步）
	go s.updateProjectStatistics(context.Background(), sourceDoc.ProjectID.Hex())

	// 5. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &serviceBase.BaseEvent{
			EventType: "document.duplicated",
			EventData: map[string]interface{}{
				"source_document_id": documentID,
				"new_document_id":    newDoc.ID.Hex(),
				"project_id":         newDoc.ProjectID.Hex(),
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	// 6. 转换为响应
	return s.duplicateService.ToResponse(newDoc), nil
}
