package writer

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// ExportService 导出服务实现
type ExportService struct {
	documentRepo        DocumentRepository
	documentContentRepo DocumentContentRepository
	projectRepo         ProjectRepository
	exportTaskRepo      ExportTaskRepository // 需要创建
	fileStorage         FileStorage          // 文件存储接口
}

// DocumentRepository 文档仓储接口
type DocumentRepository interface {
	FindByID(ctx context.Context, id string) (*writer.Document, error)
	FindByProjectID(ctx context.Context, projectID string) ([]*writer.Document, error)
}

// DocumentContentRepository 文档内容仓储接口
type DocumentContentRepository interface {
	FindByID(ctx context.Context, id string) (*writer.DocumentContent, error)
}

// ProjectRepository 项目仓储接口
type ProjectRepository interface {
	FindByID(ctx context.Context, id string) (*writer.Project, error)
}

// ExportTaskRepository 导出任务仓储接口
type ExportTaskRepository interface {
	Create(ctx context.Context, task *serviceInterfaces.ExportTask) error
	FindByID(ctx context.Context, id string) (*serviceInterfaces.ExportTask, error)
	FindByProjectID(ctx context.Context, projectID string, page, pageSize int) ([]*serviceInterfaces.ExportTask, int64, error)
	Update(ctx context.Context, task *serviceInterfaces.ExportTask) error
	Delete(ctx context.Context, id string) error
	FindByUser(ctx context.Context, userID string, page, pageSize int) ([]*serviceInterfaces.ExportTask, int64, error)
}

// FileStorage 文件存储接口
type FileStorage interface {
	Upload(ctx context.Context, filename string, content io.Reader, mimeType string) (string, error)
	Download(ctx context.Context, url string) (io.ReadCloser, error)
	Delete(ctx context.Context, url string) error
	GetSignedURL(ctx context.Context, url string, expiration time.Duration) (string, error)
}

// NewExportService 创建ExportService实例
func NewExportService(
	documentRepo DocumentRepository,
	documentContentRepo DocumentContentRepository,
	projectRepo ProjectRepository,
	exportTaskRepo ExportTaskRepository,
	fileStorage FileStorage,
) serviceInterfaces.ExportService {
	return &ExportService{
		documentRepo:        documentRepo,
		documentContentRepo: documentContentRepo,
		projectRepo:         projectRepo,
		exportTaskRepo:      exportTaskRepo,
		fileStorage:         fileStorage,
	}
}

// ExportDocument 导出文档
func (s *ExportService) ExportDocument(
	ctx context.Context,
	documentID, projectID, userID string,
	req *serviceInterfaces.ExportDocumentRequest,
) (*serviceInterfaces.ExportTask, error) {
	// 验证文档是否存在
	document, err := s.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorNotFound, "文档不存在", "", err)
	}

	// 验证项目权限
	if document.ProjectID.Hex() != projectID {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorForbidden, "无权访问此文档", "", nil)
	}

	// 创建导出任务
	task := &serviceInterfaces.ExportTask{
		ID:            primitive.NewObjectID().Hex(),
		Type:          serviceInterfaces.ExportTypeDocument,
		ResourceID:    documentID,
		ResourceTitle: document.Title,
		Format:        req.Format,
		Status:        serviceInterfaces.ExportStatusPending,
		Progress:      0,
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(24 * time.Hour), // 24小时后过期
	}

	// 保存任务到数据库
	if err := s.exportTaskRepo.Create(ctx, task); err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorInternal, "创建导出任务失败", "", err)
	}

	// 异步处理导出任务
	go s.processDocumentExport(context.Background(), task, document, req)

	return task, nil
}

// processDocumentExport 处理文档导出
func (s *ExportService) processDocumentExport(
	ctx context.Context,
	task *serviceInterfaces.ExportTask,
	document *writer.Document,
	req *serviceInterfaces.ExportDocumentRequest,
) {
	// 更新任务状态为处理中
	task.Status = serviceInterfaces.ExportStatusProcessing
	task.Progress = 10
	task.UpdatedAt = time.Now()
	s.exportTaskRepo.Update(ctx, task)

	// 获取文档内容
	content, err := s.documentContentRepo.FindByID(ctx, document.ID.Hex())
	if err != nil {
		s.failTask(ctx, task, "获取文档内容失败: "+err.Error())
		return
	}

	task.Progress = 30
	task.UpdatedAt = time.Now()
	s.exportTaskRepo.Update(ctx, task)

	// 根据格式生成文件
	var fileContent []byte
	var filename string

	switch req.Format {
	case serviceInterfaces.ExportFormatTXT:
		fileContent, _, filename = s.generateTXT(content.Content, document, req)
	case serviceInterfaces.ExportFormatMD:
		fileContent, _, filename = s.generateMarkdown(content.Content, document, req)
	case serviceInterfaces.ExportFormatDOCX:
		fileContent, _, filename = s.generateDOCX(content.Content, document, req)
	default:
		s.failTask(ctx, task, "不支持的导出格式")
		return
	}

	task.Progress = 80
	task.UpdatedAt = time.Now()
	s.exportTaskRepo.Update(ctx, task)

	// 上传文件
	// 这里需要实现文件上传逻辑
	// fileURL, err := s.fileStorage.Upload(ctx, filename, bytes.NewReader(fileContent), mimeType)
	// if err != nil {
	//     s.failTask(ctx, task, "文件上传失败: "+err.Error())
	//     return
	// }

	// 暂时使用占位符
	fileURL := fmt.Sprintf("/exports/%s/%s", task.ID, filename)

	// 更新任务为完成状态
	task.Status = serviceInterfaces.ExportStatusCompleted
	task.Progress = 100
	task.FileURL = fileURL
	task.FileSize = int64(len(fileContent))
	completedAt := time.Now()
	task.CompletedAt = &completedAt
	task.UpdatedAt = time.Now()

	s.exportTaskRepo.Update(ctx, task)
}

// generateTXT 生成TXT格式
func (s *ExportService) generateTXT(
	content string,
	document *writer.Document,
	req *serviceInterfaces.ExportDocumentRequest,
) ([]byte, string, string) {
	result := content

	if req.Options != nil && req.Options.TOC {
		result = "# " + document.Title + "\n\n" + result
	}

	filename := fmt.Sprintf("%s.txt", document.Title)
	return []byte(result), "text/plain", filename
}

// generateMarkdown 生成Markdown格式
func (s *ExportService) generateMarkdown(
	content string,
	document *writer.Document,
	req *serviceInterfaces.ExportDocumentRequest,
) ([]byte, string, string) {
	result := fmt.Sprintf("# %s\n\n", document.Title)

	if req.Options != nil && req.IncludeMeta {
		result += fmt.Sprintf("**字数**: %d\n\n", document.WordCount)
		result += fmt.Sprintf("**创建时间**: %s\n\n", document.CreatedAt.Format("2006-01-02 15:04:05"))
		if len(document.Tags) > 0 {
			result += fmt.Sprintf("**标签**: %v\n\n", document.Tags)
		}
	}

	result += content

	filename := fmt.Sprintf("%s.md", document.Title)
	return []byte(result), "text/markdown", filename
}

// generateDOCX 生成DOCX格式
func (s *ExportService) generateDOCX(
	content string,
	document *writer.Document,
	req *serviceInterfaces.ExportDocumentRequest,
) ([]byte, string, string) {
	// 这里应该使用库来生成真正的DOCX文件
	// 如 github.com/unidoc/unioffice 或其他库
	// 暂时返回简单文本
	result := fmt.Sprintf("%s\n\n%s", document.Title, content)
	filename := fmt.Sprintf("%s.docx", document.Title)
	return []byte(result), "application/vnd.openxmlformats-officedocument.wordprocessingml.document", filename
}

// failTask 标记任务失败
func (s *ExportService) failTask(ctx context.Context, task *serviceInterfaces.ExportTask, errorMsg string) {
	task.Status = serviceInterfaces.ExportStatusFailed
	task.ErrorMsg = errorMsg
	task.UpdatedAt = time.Now()
	s.exportTaskRepo.Update(ctx, task)
}

// GetExportTask 获取导出任务
func (s *ExportService) GetExportTask(ctx context.Context, taskID string) (*serviceInterfaces.ExportTask, error) {
	task, err := s.exportTaskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorNotFound, "导出任务不存在", "", err)
	}
	return task, nil
}

// DownloadExportFile 下载导出文件
func (s *ExportService) DownloadExportFile(ctx context.Context, taskID string) (*serviceInterfaces.ExportFile, error) {
	task, err := s.GetExportTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.Status != serviceInterfaces.ExportStatusCompleted {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorValidation, "导出任务未完成", "", nil)
	}

	if time.Now().After(task.ExpiresAt) {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorValidation, "导出文件已过期", "", nil)
	}

	// 生成签名URL
	signedURL, err := s.fileStorage.GetSignedURL(ctx, task.FileURL, 1*time.Hour)
	if err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorInternal, "生成下载链接失败", "", err)
	}

	return &serviceInterfaces.ExportFile{
		Filename: fmt.Sprintf("%s.%s", task.ResourceTitle, task.Format),
		URL:      signedURL,
		MimeType: s.getMimeType(task.Format),
		FileSize: task.FileSize,
	}, nil
}

// getMimeType 根据格式获取MIME类型
func (s *ExportService) getMimeType(format string) string {
	switch format {
	case serviceInterfaces.ExportFormatTXT:
		return "text/plain"
	case serviceInterfaces.ExportFormatMD:
		return "text/markdown"
	case serviceInterfaces.ExportFormatDOCX:
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case serviceInterfaces.ExportFormatZIP:
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}

// ListExportTasks 列出导出任务
func (s *ExportService) ListExportTasks(
	ctx context.Context,
	projectID string,
	page, pageSize int,
) ([]*serviceInterfaces.ExportTask, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	return s.exportTaskRepo.FindByProjectID(ctx, projectID, page, pageSize)
}

// DeleteExportTask 删除导出任务
func (s *ExportService) DeleteExportTask(ctx context.Context, taskID, userID string) error {
	task, err := s.exportTaskRepo.FindByID(ctx, taskID)
	if err != nil {
		return errors.NewServiceError("ExportService", errors.ServiceErrorNotFound, "导出任务不存在", "", err)
	}

	// 验证权限
	if task.CreatedBy != userID {
		return errors.NewServiceError("ExportService", errors.ServiceErrorForbidden, "无权删除此导出任务", "", nil)
	}

	// 删除文件
	if task.FileURL != "" {
		_ = s.fileStorage.Delete(ctx, task.FileURL)
	}

	// 删除任务记录
	return s.exportTaskRepo.Delete(ctx, taskID)
}

// ExportProject 导出项目
func (s *ExportService) ExportProject(
	ctx context.Context,
	projectID, userID string,
	req *serviceInterfaces.ExportProjectRequest,
) (*serviceInterfaces.ExportTask, error) {
	// 验证项目是否存在
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorNotFound, "项目不存在", "", err)
	}

	// 创建导出任务
	task := &serviceInterfaces.ExportTask{
		ID:            primitive.NewObjectID().Hex(),
		Type:          serviceInterfaces.ExportTypeProject,
		ResourceID:    projectID,
		ResourceTitle: project.Title,
		Format:        serviceInterfaces.ExportFormatZIP,
		Status:        serviceInterfaces.ExportStatusPending,
		Progress:      0,
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(24 * time.Hour),
	}

	// 保存任务
	if err := s.exportTaskRepo.Create(ctx, task); err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorInternal, "创建导出任务失败", "", err)
	}

	// 异步处理
	go s.processProjectExport(context.Background(), task, project, req)

	return task, nil
}

// processProjectExport 处理项目导出
func (s *ExportService) processProjectExport(
	ctx context.Context,
	task *serviceInterfaces.ExportTask,
	project *writer.Project,
	req *serviceInterfaces.ExportProjectRequest,
) {
	task.Status = serviceInterfaces.ExportStatusProcessing
	task.Progress = 10
	task.UpdatedAt = time.Now()
	s.exportTaskRepo.Update(ctx, task)

	// 获取所有文档
	if req.IncludeDocuments {
		documents, err := s.documentRepo.FindByProjectID(ctx, project.ID.Hex())
		if err != nil {
			s.failTask(ctx, task, "获取项目文档失败: "+err.Error())
			return
		}

		task.Progress = 50
		task.UpdatedAt = time.Now()
		s.exportTaskRepo.Update(ctx, task)

		// 导出所有文档（这里应该打包成ZIP）
		// 实际实现需要使用archive/zip包
		_ = documents
	}

	// 标记完成
	task.Status = serviceInterfaces.ExportStatusCompleted
	task.Progress = 100
	task.FileURL = fmt.Sprintf("/exports/%s/%s.zip", task.ID, project.Title)
	task.FileSize = 0 // 实际应该计算文件大小
	completedAt := time.Now()
	task.CompletedAt = &completedAt
	task.UpdatedAt = time.Now()

	s.exportTaskRepo.Update(ctx, task)
}

// CancelExportTask 取消导出任务
func (s *ExportService) CancelExportTask(ctx context.Context, taskID, userID string) error {
	task, err := s.exportTaskRepo.FindByID(ctx, taskID)
	if err != nil {
		return errors.NewServiceError("ExportService", errors.ServiceErrorNotFound, "导出任务不存在", "", err)
	}

	// 验证权限
	if task.CreatedBy != userID {
		return errors.NewServiceError("ExportService", errors.ServiceErrorForbidden, "无权取消此导出任务", "", nil)
	}

	// 只能取消处理中或等待中的任务
	if task.Status != serviceInterfaces.ExportStatusPending && task.Status != serviceInterfaces.ExportStatusProcessing {
		return errors.NewServiceError("ExportService", errors.ServiceErrorValidation, "任务状态不允许取消", "", nil)
	}

	task.Status = serviceInterfaces.ExportStatusCancelled
	task.UpdatedAt = time.Now()
	return s.exportTaskRepo.Update(ctx, task)
}
