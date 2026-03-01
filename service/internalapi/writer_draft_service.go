package internalapi

import (
	"context"
	"errors"

	"Qingyu_backend/models/writer"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
)

// WriterDraftService AI写作草稿服务
// 处理WriterDraft的CRUD操作，支持AI写作助手功能中的文档草稿管理
type WriterDraftService struct {
	repo writerRepo.WriterDraftRepository
}

// NewWriterDraftService 创建WriterDraftService实例
func NewWriterDraftService(repo writerRepo.WriterDraftRepository) *WriterDraftService {
	return &WriterDraftService{
		repo: repo,
	}
}

// CreateOrUpdateRequest 创建或更新请求
type CreateOrUpdateRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	ProjectID string `json:"project_id" binding:"required"`
	Action    string `json:"action" binding:"required"` // create, update, create_or_update, append
	Document  WriterDraftData
}

// WriterDraftData 文档数据
type WriterDraftData struct {
	ChapterNum int    `json:"chapter_num"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content"`
	Format     string `json:"format"`
}

// CreateOrUpdate 创建或更新文档
// 根据Action参数执行不同的操作：
//   - create: 创建新文档（如果已存在则返回错误）
//   - update: 更新现有文档（如果不存在则返回错误）
//   - create_or_update: 创建或更新文档（智能判断）
//   - append: 追加内容到现有文档
func (s *WriterDraftService) CreateOrUpdate(ctx context.Context, req *CreateOrUpdateRequest) (*writer.WriterDraft, error) {
	switch req.Action {
	case "create":
		return s.create(ctx, req)
	case "update":
		return s.update(ctx, req)
	case "create_or_update":
		return s.createOrUpdate(ctx, req)
	case "append":
		return s.append(ctx, req)
	default:
		return nil, errors.New("invalid action")
	}
}

// create 创建新文档
func (s *WriterDraftService) create(ctx context.Context, req *CreateOrUpdateRequest) (*writer.WriterDraft, error) {
	// 检查是否已存在
	existingDoc, err := s.repo.GetByProjectAndChapter(ctx, req.ProjectID, req.Document.ChapterNum)
	if err == nil && existingDoc != nil {
		return nil, errors.New("document already exists")
	}

	doc := &writer.WriterDraft{}
	doc.ProjectID = req.ProjectID
	doc.ChapterNum = req.Document.ChapterNum
	doc.Title = req.Document.Title
	doc.Content = req.Document.Content
	doc.Format = req.Document.Format

	// 设置默认值和计算字段
	doc.BeforeCreate()
	doc.UpdateContent(req.Document.Content)

	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// update 更新现有文档
func (s *WriterDraftService) update(ctx context.Context, req *CreateOrUpdateRequest) (*writer.WriterDraft, error) {
	doc, err := s.repo.GetByProjectAndChapter(ctx, req.ProjectID, req.Document.ChapterNum)
	if err != nil {
		return nil, errors.New("document not found")
	}

	if req.Document.Title != "" {
		doc.Title = req.Document.Title
	}
	if req.Document.Content != "" {
		doc.UpdateContent(req.Document.Content)
	}
	if req.Document.Format != "" {
		doc.Format = req.Document.Format
	}

	if err := s.repo.Update(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// createOrUpdate 创建或更新文档（智能判断）
func (s *WriterDraftService) createOrUpdate(ctx context.Context, req *CreateOrUpdateRequest) (*writer.WriterDraft, error) {
	doc, err := s.repo.GetByProjectAndChapter(ctx, req.ProjectID, req.Document.ChapterNum)
	if err != nil || doc == nil {
		// 不存在则创建
		return s.create(ctx, req)
	}

	// 存在则更新
	if req.Document.Title != "" {
		doc.Title = req.Document.Title
	}
	if req.Document.Content != "" {
		doc.UpdateContent(req.Document.Content)
	}
	if req.Document.Format != "" {
		doc.Format = req.Document.Format
	}

	if err := s.repo.Update(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// append 追加内容到现有文档
func (s *WriterDraftService) append(ctx context.Context, req *CreateOrUpdateRequest) (*writer.WriterDraft, error) {
	doc, err := s.repo.GetByProjectAndChapter(ctx, req.ProjectID, req.Document.ChapterNum)
	if err != nil {
		return nil, errors.New("document not found")
	}

	// 追加内容
	newContent := doc.Content + "\n\n" + req.Document.Content
	doc.UpdateContent(newContent)

	if err := s.repo.Update(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// GetDocument 获取文档
func (s *WriterDraftService) GetDocument(ctx context.Context, userID, projectID, documentID string) (*writer.WriterDraft, error) {
	return s.repo.GetByID(ctx, documentID)
}

// ListDocuments 列出文档
func (s *WriterDraftService) ListDocuments(ctx context.Context, userID, projectID string, limit int) ([]*writer.WriterDraft, int, error) {
	docs, err := s.repo.ListByProject(ctx, projectID, limit)
	return docs, len(docs), err
}

// DeleteDocument 删除文档
func (s *WriterDraftService) DeleteDocument(ctx context.Context, userID, projectID, documentID string) error {
	return s.repo.Delete(ctx, documentID)
}

// BatchGetDocuments 批量获取文档
func (s *WriterDraftService) BatchGetDocuments(ctx context.Context, userID, projectID string, documentIDs []string) ([]*writer.WriterDraft, error) {
	return s.repo.BatchGetByIDs(ctx, documentIDs)
}

// GetDocumentByProjectAndChapter 根据项目和章节获取文档
func (s *WriterDraftService) GetDocumentByProjectAndChapter(ctx context.Context, projectID string, chapterNum int) (*writer.WriterDraft, error) {
	return s.repo.GetByProjectAndChapter(ctx, projectID, chapterNum)
}
