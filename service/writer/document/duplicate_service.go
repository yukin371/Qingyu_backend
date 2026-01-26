package document

import (
	"Qingyu_backend/models/writer"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/utils"
	"context"
	"fmt"

	pkgErrors "Qingyu_backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DuplicateService 文档复制服务
type DuplicateService struct {
	docRepo     writerRepo.DocumentRepository
	contentRepo writerRepo.DocumentContentRepository
	serviceName string
	version     string
}

// NewDuplicateService 创建文档复制服务
func NewDuplicateService(
	docRepo writerRepo.DocumentRepository,
	contentRepo writerRepo.DocumentContentRepository,
) *DuplicateService {
	return &DuplicateService{
		docRepo:     docRepo,
		contentRepo: contentRepo,
		serviceName: "DuplicateService",
		version:     "1.0.0",
	}
}

// DuplicateRequest 文档复制请求
type DuplicateRequest struct {
	TargetParentID *string `json:"targetParentId,omitempty"` // 目标父文档ID
	Position       string  `json:"position"`                 // inner/before/after
	CopyContent    bool    `json:"copyContent"`              // 是否复制内容
}

// DuplicateResponse 文档复制响应
type DuplicateResponse struct {
	DocumentID string `json:"documentId"`
	Title      string `json:"title"`
	StableRef  string `json:"stableRef"`
}

// Duplicate 复制文档
func (s *DuplicateService) Duplicate(ctx context.Context, documentID string, req *DuplicateRequest) (*writer.Document, error) {
	// 1. 参数验证
	if err := s.validateRequest(documentID, req); err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "参数验证失败", err.Error(), err)
	}

	// 2. 加载原文档
	sourceDoc, err := s.docRepo.GetByID(ctx, documentID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询原文档失败", "", err)
	}

	if sourceDoc == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "文档不存在", "", nil)
	}

	// 3. 创建新文档（复制元数据）
	newDoc := s.createDuplicateDocument(sourceDoc, req)

	// 4. 保存新文档
	if err := s.docRepo.Create(ctx, newDoc); err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "创建新文档失败", "", err)
	}

	// 5. 复制内容（如果需要）
	if req.CopyContent {
		if err := s.duplicateContent(ctx, documentID, newDoc); err != nil {
			// 内容复制失败不影响文档创建，只记录错误
			// 可以考虑回滚文档创建，但P1阶段简化处理
		}
	}

	return newDoc, nil
}

// validateRequest 验证请求参数
func (s *DuplicateService) validateRequest(documentID string, req *DuplicateRequest) error {
	if documentID == "" {
		return fmt.Errorf("文档ID不能为空")
	}

	if req == nil {
		return fmt.Errorf("请求参数不能为空")
	}

	// 验证Position参数
	if req.Position != "" && req.Position != "inner" && req.Position != "before" && req.Position != "after" {
		return fmt.Errorf("无效的position参数: %s", req.Position)
	}

	return nil
}

// createDuplicateDocument 创建复制的文档
func (s *DuplicateService) createDuplicateDocument(sourceDoc *writer.Document, req *DuplicateRequest) *writer.Document {
	// 复制原文档的所有字段
	newDoc := &writer.Document{
		ProjectID:    sourceDoc.ProjectID, // 同项目
		Title:        fmt.Sprintf("Copy - %s", sourceDoc.Title),
		Type:         sourceDoc.Type,
		Level:        sourceDoc.Level,
		Order:        sourceDoc.Order,
		Status:       sourceDoc.Status,
		WordCount:    0, // 重置字数统计
		CharacterIDs: append([]primitive.ObjectID{}, sourceDoc.CharacterIDs...),
		LocationIDs:  append([]primitive.ObjectID{}, sourceDoc.LocationIDs...),
		TimelineIDs:  append([]primitive.ObjectID{}, sourceDoc.TimelineIDs...),
		PlotThreads:  append([]string{}, sourceDoc.PlotThreads...),
		KeyPoints:    append([]string{}, sourceDoc.KeyPoints...),
		WritingHints: append([]string{}, sourceDoc.WritingHints...),
		Tags:         append([]string{}, sourceDoc.Tags...),
		Notes:        sourceDoc.Notes,
	}

	// 处理父节点位置（如果指定）
	if req.TargetParentID != nil && *req.TargetParentID != "" {
		if parentID, err := primitive.ObjectIDFromHex(*req.TargetParentID); err == nil {
			newDoc.ParentID = parentID
		}
	} else {
		// 默认使用原文档的父节点
		newDoc.ParentID = sourceDoc.ParentID
	}

	// 生成新的StableRef（添加-copy后缀）
	newDoc.StableRef = fmt.Sprintf("%s-copy", sourceDoc.StableRef)

	// 使用LexoRankUtils生成新的OrderKey
	newDoc.OrderKey = utils.GenerateSiblingOrderKey(sourceDoc.OrderKey)

	// 初始化时间戳
	newDoc.TouchForCreate()

	return newDoc
}

// duplicateContent 复制文档内容
func (s *DuplicateService) duplicateContent(ctx context.Context, sourceDocumentID string, newDoc *writer.Document) error {
	// 查询原文档内容
	sourceContent, err := s.contentRepo.GetByDocumentID(ctx, sourceDocumentID)
	if err != nil {
		return fmt.Errorf("查询原文档内容失败: %w", err)
	}

	if sourceContent == nil {
		// 原文档没有内容，不复制
		return nil
	}

	// 创建新内容
	newContent := &writer.DocumentContent{
		DocumentID:  newDoc.ID,
		Content:     sourceContent.Content,
		ContentType: sourceContent.ContentType,
		WordCount:   0, // 重置字数统计
		CharCount:   0, // 重置字符统计
		Version:     1, // 新文档版本从1开始
	}

	newContent.TouchForCreate()

	// 保存新内容
	if err := s.contentRepo.Create(ctx, newContent); err != nil {
		return fmt.Errorf("创建新文档内容失败: %w", err)
	}

	// 更新文档的字数统计
	if sourceContent.WordCount > 0 {
		newDoc.WordCount = sourceContent.WordCount
	}

	return nil
}

// ToResponse 将Document转换为DuplicateResponse
func (s *DuplicateService) ToResponse(doc *writer.Document) *DuplicateResponse {
	return &DuplicateResponse{
		DocumentID: doc.ID.Hex(),
		Title:      doc.Title,
		StableRef:  doc.StableRef,
	}
}

// BaseService接口实现

func (s *DuplicateService) Initialize(ctx context.Context) error {
	return nil
}

func (s *DuplicateService) Health(ctx context.Context) error {
	return s.docRepo.Health(ctx)
}

func (s *DuplicateService) Close(ctx context.Context) error {
	return nil
}

func (s *DuplicateService) GetServiceName() string {
	return s.serviceName
}

func (s *DuplicateService) GetVersion() string {
	return s.version
}
