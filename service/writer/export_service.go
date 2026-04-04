package writer

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	serviceInterfaces "Qingyu_backend/service/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ExportService 导出服务实现
type ExportService struct {
	documentRepo        DocumentRepository
	documentContentRepo DocumentContentRepository
	projectRepo         ProjectRepository
	exportTaskRepo      ExportTaskRepository
	fileStorage         FileStorage
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

// NewExportService 创建 ExportService 实例
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
	document, err := s.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorNotFound, "文档不存在", "", err)
	}

	if document.ProjectID.Hex() != projectID {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorForbidden, "无权访问此文档", "", nil)
	}

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
		ExpiresAt:     time.Now().Add(24 * time.Hour),
	}

	if err := s.exportTaskRepo.Create(ctx, task); err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorInternal, "创建导出任务失败", "", err)
	}

	go s.processDocumentExport(context.Background(), task, document, req)
	return task, nil
}

func (s *ExportService) processDocumentExport(
	ctx context.Context,
	task *serviceInterfaces.ExportTask,
	document *writer.Document,
	req *serviceInterfaces.ExportDocumentRequest,
) {
	s.updateTaskProgress(ctx, task, serviceInterfaces.ExportStatusProcessing, 10)

	content, err := s.documentContentRepo.FindByID(ctx, document.ID.Hex())
	if err != nil {
		s.failTask(ctx, task, "获取文档内容失败: "+err.Error())
		return
	}

	s.updateTaskProgress(ctx, task, serviceInterfaces.ExportStatusProcessing, 30)

	fileContent, mimeType, filename, err := s.renderDocumentExport(content.Content, document, req)
	if err != nil {
		s.failTask(ctx, task, err.Error())
		return
	}

	s.updateTaskProgress(ctx, task, serviceInterfaces.ExportStatusProcessing, 80)

	if err := s.completeTaskWithFile(ctx, task, filename, mimeType, fileContent); err != nil {
		s.failTask(ctx, task, err.Error())
	}
}

func (s *ExportService) renderDocumentExport(
	content string,
	document *writer.Document,
	req *serviceInterfaces.ExportDocumentRequest,
) ([]byte, string, string, error) {
	switch req.Format {
	case serviceInterfaces.ExportFormatTXT:
		fileContent, mimeType, filename := s.generateTXT(content, document, req)
		return fileContent, mimeType, filename, nil
	case serviceInterfaces.ExportFormatMD:
		fileContent, mimeType, filename := s.generateMarkdown(content, document, req)
		return fileContent, mimeType, filename, nil
	case serviceInterfaces.ExportFormatDOCX:
		fileContent, mimeType, filename := s.generateDOCX(content, document, req)
		return fileContent, mimeType, filename, nil
	default:
		return nil, "", "", fmt.Errorf("不支持的导出格式: %s", req.Format)
	}
}

// generateTXT 生成 TXT 格式
func (s *ExportService) generateTXT(
	content string,
	document *writer.Document,
	req *serviceInterfaces.ExportDocumentRequest,
) ([]byte, string, string) {
	result := tiptapToPlainText(content)
	if req.Options != nil && req.Options.TOC {
		result = document.Title + "\n\n" + result
	}

	filename := fmt.Sprintf("%s.txt", sanitizeFileName(document.Title))
	return []byte(result), "text/plain; charset=utf-8", filename
}

// generateMarkdown 生成 Markdown 格式
func (s *ExportService) generateMarkdown(
	content string,
	document *writer.Document,
	req *serviceInterfaces.ExportDocumentRequest,
) ([]byte, string, string) {
	var builder strings.Builder
	builder.WriteString("# ")
	builder.WriteString(document.Title)
	builder.WriteString("\n\n")

	if req.IncludeMeta || hasIncludeMeta(req.Options) {
		builder.WriteString(fmt.Sprintf("**字数**: %d\n\n", document.WordCount))
		builder.WriteString(fmt.Sprintf("**创建时间**: %s\n\n", document.CreatedAt.Format("2006-01-02 15:04:05")))
		if len(document.Tags) > 0 {
			builder.WriteString(fmt.Sprintf("**标签**: %s\n\n", strings.Join(document.Tags, "、")))
		}
		builder.WriteString("---\n\n")
	}

	builder.WriteString(tiptapToMarkdown(content))

	filename := fmt.Sprintf("%s.md", sanitizeFileName(document.Title))
	return []byte(builder.String()), "text/markdown; charset=utf-8", filename
}

// generateDOCX 生成 DOCX 格式
func (s *ExportService) generateDOCX(
	content string,
	document *writer.Document,
	_ *serviceInterfaces.ExportDocumentRequest,
) ([]byte, string, string) {
	plainText := tiptapToPlainText(content)
	docxData, err := generateMinimalDOCX(document.Title, plainText)
	if err != nil {
		// 降级到可打开的纯文本内容，避免任务直接失败
		docxData = []byte(document.Title + "\n\n" + plainText)
	}

	filename := fmt.Sprintf("%s.docx", sanitizeFileName(document.Title))
	return docxData, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", filename
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
	if task.FileURL == "" {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorNotFound, "导出文件不存在", "", nil)
	}

	exportFile := &serviceInterfaces.ExportFile{
		Filename: fmt.Sprintf("%s.%s", sanitizeFileName(task.ResourceTitle), task.Format),
		MimeType: s.getMimeType(task.Format),
		FileSize: task.FileSize,
	}

	if s.fileStorage != nil {
		reader, err := s.fileStorage.Download(ctx, task.FileURL)
		if err == nil {
			defer reader.Close()
			content, readErr := io.ReadAll(reader)
			if readErr == nil {
				exportFile.Content = content
				return exportFile, nil
			}
		}

		signedURL, err := s.fileStorage.GetSignedURL(ctx, task.FileURL, time.Hour)
		if err == nil {
			exportFile.URL = signedURL
		}
	} else {
		exportFile.URL = task.FileURL
	}

	if len(exportFile.Content) == 0 && exportFile.URL == "" {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorInternal, "读取导出文件失败", "", nil)
	}

	return exportFile, nil
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

	if task.CreatedBy != userID {
		return errors.NewServiceError("ExportService", errors.ServiceErrorForbidden, "无权删除此导出任务", "", nil)
	}

	if task.FileURL != "" && s.fileStorage != nil {
		_ = s.fileStorage.Delete(ctx, task.FileURL)
	}

	return s.exportTaskRepo.Delete(ctx, taskID)
}

// ExportProject 导出项目
func (s *ExportService) ExportProject(
	ctx context.Context,
	projectID, userID string,
	req *serviceInterfaces.ExportProjectRequest,
) (*serviceInterfaces.ExportTask, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorNotFound, "项目不存在", "", err)
	}

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

	if err := s.exportTaskRepo.Create(ctx, task); err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorInternal, "创建导出任务失败", "", err)
	}

	go s.processProjectExport(context.Background(), task, project, req)
	return task, nil
}

func (s *ExportService) processProjectExport(
	ctx context.Context,
	task *serviceInterfaces.ExportTask,
	project *writer.Project,
	req *serviceInterfaces.ExportProjectRequest,
) {
	s.updateTaskProgress(ctx, task, serviceInterfaces.ExportStatusProcessing, 10)

	archiveData, err := s.buildProjectArchive(ctx, project, req)
	if err != nil {
		s.failTask(ctx, task, err.Error())
		return
	}

	s.updateTaskProgress(ctx, task, serviceInterfaces.ExportStatusProcessing, 80)

	filename := fmt.Sprintf("%s.zip", sanitizeFileName(project.Title))
	if err := s.completeTaskWithFile(ctx, task, filename, "application/zip", archiveData); err != nil {
		s.failTask(ctx, task, err.Error())
	}
}

// CancelExportTask 取消导出任务
func (s *ExportService) CancelExportTask(ctx context.Context, taskID, userID string) error {
	task, err := s.exportTaskRepo.FindByID(ctx, taskID)
	if err != nil {
		return errors.NewServiceError("ExportService", errors.ServiceErrorNotFound, "导出任务不存在", "", err)
	}

	if task.CreatedBy != userID {
		return errors.NewServiceError("ExportService", errors.ServiceErrorForbidden, "无权取消此导出任务", "", nil)
	}

	if task.Status != serviceInterfaces.ExportStatusPending && task.Status != serviceInterfaces.ExportStatusProcessing {
		return errors.NewServiceError("ExportService", errors.ServiceErrorValidation, "任务状态不允许取消", "", nil)
	}

	task.Status = serviceInterfaces.ExportStatusCancelled
	task.UpdatedAt = time.Now()
	return s.exportTaskRepo.Update(ctx, task)
}

// ExportProjectAsZip 将项目导出为 ZIP 字节数据（直接返回，不创建异步任务）
func (s *ExportService) ExportProjectAsZip(ctx context.Context, projectID, _ string) ([]byte, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorNotFound, "项目不存在", "", err)
	}

	defaultReq := &serviceInterfaces.ExportProjectRequest{
		IncludeDocuments: true,
		DocumentFormats:  serviceInterfaces.ExportFormatTXT,
	}
	return s.buildProjectArchive(ctx, project, defaultReq)
}

func (s *ExportService) buildProjectArchive(
	ctx context.Context,
	project *writer.Project,
	req *serviceInterfaces.ExportProjectRequest,
) ([]byte, error) {
	documents, err := s.documentRepo.FindByProjectID(ctx, project.ID.Hex())
	if err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorInternal, "获取项目文档失败", "", err)
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	rootFolder := sanitizeZipFileName(project.Title)

	if req == nil || req.IncludeDocuments {
		docTree := buildDocumentTree(documents)
		format := normalizeProjectDocumentFormat(req)
		if err := s.addDocumentsToZip(ctx, zipWriter, docTree, rootFolder, "", format, req); err != nil {
			_ = zipWriter.Close()
			return nil, errors.NewServiceError("ExportService", errors.ServiceErrorInternal, "打包文档失败", "", err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorInternal, "生成 ZIP 失败", "", err)
	}

	return buf.Bytes(), nil
}

// ImportProject 从 ZIP 数据导入项目
func (s *ExportService) ImportProject(ctx context.Context, userID string, zipData []byte) (*serviceInterfaces.ImportResult, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorValidation, "无效的 ZIP 文件", "", err)
	}

	rootFolder := findZipRootFolder(reader)
	if rootFolder == "" {
		return nil, errors.NewServiceError("ExportService", errors.ServiceErrorValidation, "无法找到项目根目录", "", nil)
	}

	projectTitle := strings.TrimSuffix(rootFolder, "/")
	project := &writer.Project{
		WritingType: "novel",
		Status:      writer.StatusDraft,
		Visibility:  writer.VisibilityPrivate,
		Summary:     fmt.Sprintf("从文件导入于 %s", time.Now().Format("2006-01-02 15:04:05")),
	}
	project.Title = projectTitle
	project.AuthorID, _ = primitive.ObjectIDFromHex(userID)
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()
	project.ID = primitive.NewObjectID()

	documentCount := 0
	folderDocMap := make(map[string]string)

	for _, file := range reader.File {
		if file.FileInfo().IsDir() || !strings.HasPrefix(file.Name, rootFolder) {
			continue
		}

		relativePath := strings.TrimPrefix(strings.TrimPrefix(file.Name, rootFolder), "/")
		if relativePath == "" {
			continue
		}

		pathParts := strings.Split(relativePath, "/")
		parentID := ""
		for i := 0; i < len(pathParts)-1; i++ {
			folderPath := strings.Join(pathParts[:i+1], "/")
			if folderDocID, exists := folderDocMap[folderPath]; exists {
				parentID = folderDocID
			}
		}

		if strings.HasSuffix(file.Name, ".txt") || strings.HasSuffix(file.Name, ".md") || strings.HasSuffix(file.Name, ".docx") {
			rc, openErr := file.Open()
			if openErr != nil {
				continue
			}

			var fileBuf bytes.Buffer
			_, readErr := fileBuf.ReadFrom(rc)
			_ = rc.Close()
			if readErr != nil {
				continue
			}

			title := strings.TrimSuffix(pathParts[len(pathParts)-1], filepathExt(pathParts[len(pathParts)-1]))
			docID := primitive.NewObjectID().Hex()
			documentCount++
			if len(pathParts) > 1 {
				folderPath := strings.Join(pathParts[:len(pathParts)-1], "/")
				folderDocMap[folderPath] = docID
			}

			_ = parentID
			_ = title
			_ = fileBuf.String()
			continue
		}

		folderPath := strings.TrimSuffix(relativePath, "/")
		folderDocMap[folderPath] = primitive.NewObjectID().Hex()
	}

	return &serviceInterfaces.ImportResult{
		ProjectID:     project.ID.Hex(),
		Title:         projectTitle,
		DocumentCount: documentCount,
	}, nil
}

func (s *ExportService) addDocumentsToZip(
	ctx context.Context,
	zipWriter *zip.Writer,
	docTree map[string][]*writer.Document,
	currentPath string,
	parentID string,
	format string,
	req *serviceInterfaces.ExportProjectRequest,
) error {
	children := docTree[parentID]
	for _, doc := range children {
		docID := doc.ID.Hex()
		fileName := sanitizeZipFileName(doc.Title)
		nextPath := currentPath
		if nextPath != "" {
			nextPath += "/"
		}
		nextPath += fileName

		hasChildren := len(docTree[docID]) > 0

		content, err := s.documentContentRepo.FindByID(ctx, docID)
		if err == nil && content != nil && content.Content != "" {
			fileReq := projectRequestToDocumentRequest(format, req)
			fileContent, _, exportedFileName, renderErr := s.renderDocumentExport(content.Content, doc, fileReq)
			if renderErr != nil {
				return renderErr
			}

			filePath := nextPath + "/" + exportedFileName
			if !hasChildren {
				filePath = currentPath + "/" + exportedFileName
			}

			writerHandle, createErr := zipWriter.Create(strings.TrimPrefix(filePath, "/"))
			if createErr != nil {
				return createErr
			}

			if _, writeErr := writerHandle.Write(fileContent); writeErr != nil {
				return writeErr
			}
		}

		if hasChildren {
			if err := s.addDocumentsToZip(ctx, zipWriter, docTree, nextPath, docID, format, req); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ExportService) completeTaskWithFile(
	ctx context.Context,
	task *serviceInterfaces.ExportTask,
	filename string,
	mimeType string,
	fileContent []byte,
) error {
	if s.fileStorage == nil {
		return fmt.Errorf("文件存储未配置")
	}

	fileURL, err := s.fileStorage.Upload(ctx, filename, bytes.NewReader(fileContent), mimeType)
	if err != nil {
		return fmt.Errorf("文件上传失败: %w", err)
	}

	task.Status = serviceInterfaces.ExportStatusCompleted
	task.Progress = 100
	task.FileURL = fileURL
	task.FileSize = int64(len(fileContent))
	completedAt := time.Now()
	task.CompletedAt = &completedAt
	task.UpdatedAt = time.Now()

	return s.exportTaskRepo.Update(ctx, task)
}

func (s *ExportService) updateTaskProgress(
	ctx context.Context,
	task *serviceInterfaces.ExportTask,
	status string,
	progress int,
) {
	task.Status = status
	task.Progress = progress
	task.UpdatedAt = time.Now()
	_ = s.exportTaskRepo.Update(ctx, task)
}

func (s *ExportService) failTask(ctx context.Context, task *serviceInterfaces.ExportTask, errorMsg string) {
	task.Status = serviceInterfaces.ExportStatusFailed
	task.ErrorMsg = errorMsg
	task.UpdatedAt = time.Now()
	_ = s.exportTaskRepo.Update(ctx, task)
}

func (s *ExportService) getMimeType(format string) string {
	switch format {
	case serviceInterfaces.ExportFormatTXT:
		return "text/plain; charset=utf-8"
	case serviceInterfaces.ExportFormatMD:
		return "text/markdown; charset=utf-8"
	case serviceInterfaces.ExportFormatDOCX:
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case serviceInterfaces.ExportFormatZIP:
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}

func normalizeProjectDocumentFormat(req *serviceInterfaces.ExportProjectRequest) string {
	if req == nil {
		return serviceInterfaces.ExportFormatTXT
	}

	switch req.DocumentFormats {
	case serviceInterfaces.ExportFormatTXT, serviceInterfaces.ExportFormatMD, serviceInterfaces.ExportFormatDOCX:
		return req.DocumentFormats
	default:
		return serviceInterfaces.ExportFormatTXT
	}
}

func projectRequestToDocumentRequest(
	format string,
	req *serviceInterfaces.ExportProjectRequest,
) *serviceInterfaces.ExportDocumentRequest {
	return &serviceInterfaces.ExportDocumentRequest{
		Format: format,
		Options: func() *serviceInterfaces.ExportOptions {
			if req == nil {
				return nil
			}
			return req.Options
		}(),
	}
}

func buildDocumentTree(documents []*writer.Document) map[string][]*writer.Document {
	tree := make(map[string][]*writer.Document)
	for _, doc := range documents {
		parentID := ""
		if !doc.ParentID.IsZero() {
			parentID = doc.ParentID.Hex()
		}
		tree[parentID] = append(tree[parentID], doc)
	}

	for parentID := range tree {
		sort.Slice(tree[parentID], func(i, j int) bool {
			if tree[parentID][i].Order == tree[parentID][j].Order {
				return tree[parentID][i].Title < tree[parentID][j].Title
			}
			return tree[parentID][i].Order < tree[parentID][j].Order
		})
	}

	return tree
}

func findZipRootFolder(reader *zip.Reader) string {
	for _, file := range reader.File {
		if idx := strings.Index(file.Name, "/"); idx > 0 {
			return file.Name[:idx+1]
		}
	}
	return ""
}

func sanitizeFileName(name string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	result := strings.TrimSpace(replacer.Replace(name))
	if len(result) > 100 {
		result = result[:100]
	}
	if result == "" {
		return "untitled"
	}
	return result
}

func sanitizeZipFileName(name string) string {
	return sanitizeFileName(name)
}

func filepathExt(name string) string {
	idx := strings.LastIndex(name, ".")
	if idx < 0 {
		return ""
	}
	return name[idx:]
}

func hasIncludeMeta(options *serviceInterfaces.ExportOptions) bool {
	if options == nil {
		return false
	}
	return options.IncludeTags
}

type tiptapNode struct {
	Type    string       `json:"type"`
	Text    string       `json:"text"`
	Content []tiptapNode `json:"content"`
}

func tiptapToPlainText(content string) string {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return ""
	}

	if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
		return trimmed
	}

	var root tiptapNode
	if err := json.Unmarshal([]byte(trimmed), &root); err == nil {
		text := strings.TrimSpace(extractTipTapText(root))
		if text != "" {
			return text
		}
	}

	var nodes []tiptapNode
	if err := json.Unmarshal([]byte(trimmed), &nodes); err == nil {
		parts := make([]string, 0, len(nodes))
		for _, node := range nodes {
			text := strings.TrimSpace(extractTipTapText(node))
			if text != "" {
				parts = append(parts, text)
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, "\n\n")
		}
	}

	return trimmed
}

func extractTipTapText(node tiptapNode) string {
	if node.Type == "hardBreak" {
		return "\n"
	}
	if node.Text != "" {
		return node.Text
	}

	parts := make([]string, 0, len(node.Content))
	for _, child := range node.Content {
		text := extractTipTapText(child)
		if text != "" {
			parts = append(parts, text)
		}
	}

	switch node.Type {
	case "paragraph", "heading", "blockquote", "codeBlock":
		return strings.Join(parts, "") + "\n\n"
	case "bulletList", "orderedList", "doc":
		return strings.Join(parts, "")
	case "listItem":
		return "- " + strings.TrimSpace(strings.Join(parts, "")) + "\n"
	default:
		return strings.Join(parts, "")
	}
}

func tiptapToMarkdown(content string) string {
	return tiptapToPlainText(content)
}

func generateMinimalDOCX(title string, plainText string) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`,
		"word/document.xml": buildDOCXDocumentXML(title, plainText),
	}

	for path, content := range files {
		fileWriter, err := zipWriter.Create(path)
		if err != nil {
			_ = zipWriter.Close()
			return nil, err
		}
		if _, err := fileWriter.Write([]byte(content)); err != nil {
			_ = zipWriter.Close()
			return nil, err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buildDOCXDocumentXML(title string, plainText string) string {
	var paragraphs strings.Builder
	paragraphs.WriteString(docxParagraph(title, true))

	for _, block := range splitParagraphs(plainText) {
		paragraphs.WriteString(docxParagraph(block, false))
	}

	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<w:document xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas"` +
		` xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"` +
		` xmlns:o="urn:schemas-microsoft-com:office:office"` +
		` xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"` +
		` xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math"` +
		` xmlns:v="urn:schemas-microsoft-com:vml"` +
		` xmlns:wp14="http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing"` +
		` xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"` +
		` xmlns:w10="urn:schemas-microsoft-com:office:word"` +
		` xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"` +
		` xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml"` +
		` xmlns:wpg="http://schemas.microsoft.com/office/word/2010/wordprocessingGroup"` +
		` xmlns:wpi="http://schemas.microsoft.com/office/word/2010/wordprocessingInk"` +
		` xmlns:wne="http://schemas.microsoft.com/office/word/2006/wordml"` +
		` xmlns:wps="http://schemas.microsoft.com/office/word/2010/wordprocessingShape" mc:Ignorable="w14 wp14">` +
		`<w:body>` + paragraphs.String() + `<w:sectPr/></w:body></w:document>`
}

func docxParagraph(text string, heading bool) string {
	escaped := xmlEscape(text)
	if heading {
		return `<w:p><w:pPr><w:pStyle w:val="Heading1"/></w:pPr><w:r><w:t xml:space="preserve">` + escaped + `</w:t></w:r></w:p>`
	}
	return `<w:p><w:r><w:t xml:space="preserve">` + escaped + `</w:t></w:r></w:p>`
}

func splitParagraphs(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	raw := strings.Split(text, "\n")
	parts := make([]string, 0, len(raw))
	for _, line := range raw {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	if len(parts) == 0 {
		return []string{""}
	}
	return parts
}

func xmlEscape(text string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(text))
	return buf.String()
}
