package document

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/document"
	pkgErrors "Qingyu_backend/pkg/errors"
	writingRepo "Qingyu_backend/repository/interfaces/writing"
	"Qingyu_backend/service/base"
)

// DocumentService 文档服务
type DocumentService struct {
	documentRepo writingRepo.DocumentRepository
	projectRepo  writingRepo.ProjectRepository
	eventBus     base.EventBus
	serviceName  string
	version      string
}

// NewDocumentService 创建文档服务
func NewDocumentService(
	documentRepo writingRepo.DocumentRepository,
	projectRepo writingRepo.ProjectRepository,
	eventBus base.EventBus,
) *DocumentService {
	return &DocumentService{
		documentRepo: documentRepo,
		projectRepo:  projectRepo,
		eventBus:     eventBus,
		serviceName:  "DocumentService",
		version:      "1.0.0",
	}
}

// CreateDocument 创建文档
func (s *DocumentService) CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*CreateDocumentResponse, error) {
	// 1. 参数验证
	if err := s.validateCreateDocumentRequest(req); err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "参数验证失败", err.Error(), err)
	}

	// 2. 验证项目权限
	project, err := s.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	if project == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "项目不存在", "", nil)
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	if !project.CanEdit(userID) {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "无权限编辑该项目", "", nil)
	}

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

		if !parent.CanHaveChildren() {
			return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorBusiness, "该文档不能添加子文档（最多3层）", "", nil)
		}

		level = parent.GetNextLevel()
	}

	// 4. 创建文档对象
	doc := &document.Document{
		ProjectID:    req.ProjectID,
		ParentID:     req.ParentID,
		Title:        req.Title,
		Type:         document.DocumentType(req.Type),
		Level:        level,
		Order:        req.Order,
		Status:       "planned",
		WordCount:    0,
		CharacterIDs: req.CharacterIDs,
		LocationIDs:  req.LocationIDs,
		TimelineIDs:  req.TimelineIDs,
		Tags:         req.Tags,
		Notes:        req.Notes,
	}

	// 5. 保存文档
	if err := s.documentRepo.Create(ctx, doc); err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "创建文档失败", "", err)
	}

	// 6. 更新项目统计（异步）
	go s.updateProjectStatistics(context.Background(), req.ProjectID)

	// 7. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &base.BaseEvent{
			EventType: "document.created",
			EventData: map[string]interface{}{
				"document_id": doc.ID,
				"project_id":  doc.ProjectID,
				"title":       doc.Title,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return &CreateDocumentResponse{
		DocumentID: doc.ID,
		Title:      doc.Title,
		Type:       string(doc.Type),
		CreatedAt:  doc.CreatedAt,
	}, nil
}

// GetDocument 获取文档详情
func (s *DocumentService) GetDocument(ctx context.Context, documentID string) (*document.Document, error) {
	// 1. 查询文档
	doc, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询文档失败", "", err)
	}

	if doc == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "文档不存在", "", nil)
	}

	// 2. 验证项目权限
	project, err := s.projectRepo.GetByID(ctx, doc.ProjectID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	if project == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "项目不存在", "", nil)
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	if !project.CanView(userID) {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "无权限查看该文档", "", nil)
	}

	return doc, nil
}

// GetDocumentTree 获取文档树
func (s *DocumentService) GetDocumentTree(ctx context.Context, projectID string) (*DocumentTreeResponse, error) {
	// 1. 验证项目权限
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	if project == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "项目不存在", "", nil)
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	if !project.CanView(userID) {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "无权限查看该项目", "", nil)
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
	// 1. 查询文档
	doc, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询文档失败", "", err)
	}

	if doc == nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "文档不存在", "", nil)
	}

	// 2. 验证项目权限
	project, err := s.projectRepo.GetByID(ctx, doc.ProjectID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	if !project.CanEdit(userID) {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "无权限编辑该文档", "", nil)
	}

	// 3. 构建更新数据
	updates := make(map[string]interface{})

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}
	if req.Status != "" {
		updates["status"] = req.Status
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

	// 4. 更新文档
	if err := s.documentRepo.UpdateByProject(ctx, documentID, doc.ProjectID, updates); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "更新文档失败", "", err)
	}

	// 5. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &base.BaseEvent{
			EventType: "document.updated",
			EventData: map[string]interface{}{
				"document_id": documentID,
				"project_id":  doc.ProjectID,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// DeleteDocument 删除文档
func (s *DocumentService) DeleteDocument(ctx context.Context, documentID string) error {
	// 1. 获取文档
	doc, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询文档失败", "", err)
	}

	if doc == nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "文档不存在", "", nil)
	}

	// 2. 验证项目权限
	project, err := s.projectRepo.GetByID(ctx, doc.ProjectID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	if !project.CanEdit(userID) {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "无权限编辑该项目", "", nil)
	}

	// 3. 软删除文档
	if err := s.documentRepo.SoftDelete(ctx, documentID, doc.ProjectID); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "删除文档失败", "", err)
	}

	// 4. 更新项目统计（异步）
	go s.updateProjectStatistics(context.Background(), doc.ProjectID)

	// 5. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &base.BaseEvent{
			EventType: "document.deleted",
			EventData: map[string]interface{}{
				"document_id": documentID,
				"project_id":  doc.ProjectID,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// 私有方法
func (s *DocumentService) buildDocumentTree(documents []*document.Document) []*DocumentTreeNode {
	// 构建树形结构
	nodeMap := make(map[string]*DocumentTreeNode)
	var rootNodes []*DocumentTreeNode

	// 第一遍遍历：创建所有节点
	for _, doc := range documents {
		node := &DocumentTreeNode{
			Document: doc,
			Children: []*DocumentTreeNode{},
		}
		nodeMap[doc.ID] = node
	}

	// 第二遍遍历：建立父子关系
	for _, doc := range documents {
		node := nodeMap[doc.ID]
		if doc.ParentID == "" {
			rootNodes = append(rootNodes, node)
		} else {
			if parent, exists := nodeMap[doc.ParentID]; exists {
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

	// 验证文档类型
	docType := document.DocumentType(req.Type)
	if !docType.IsValid() {
		return fmt.Errorf("无效的文档类型: %s", req.Type)
	}

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
		if doc.Type == document.TypeChapter {
			chapterCount++
		}
	}

	// 更新项目统计
	stats := &document.ProjectStats{
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
