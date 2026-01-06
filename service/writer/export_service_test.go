package writer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/service/interfaces"
	"Qingyu_backend/service/writer/mocks"
)

// ============ 测试辅助函数 ============

// createTestExportTask 创建测试导出任务
func createTestExportTask(id, resourceID, resourceTitle, format, createdBy string) *interfaces.ExportTask {
	return &interfaces.ExportTask{
		ID:            id,
		Type:          interfaces.ExportTypeDocument,
		ResourceID:    resourceID,
		ResourceTitle: resourceTitle,
		Format:        format,
		Status:        interfaces.ExportStatusPending,
		Progress:      0,
		CreatedBy:     createdBy,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(24 * time.Hour),
	}
}

// ============ ExportService 测试 ============

// TestNewExportService 测试服务创建
func TestNewExportService(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	assert.NotNil(t, service, "服务不应为空")
}

// ============ ExportDocument 测试 ============

// TestExportDocument_Success 测试导出文档成功
func TestExportDocument_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	documentID := "doc123"
	projectID := "proj123"
	userID := "user123"
	document := createTestDocument(documentID, projectID, "Test Chapter")

	// 设置Mock期望
	mockDocRepo.On("FindByID", mock.Anything, documentID).Return(document, nil)
	mockExportRepo.On("Create", mock.Anything, mock.MatchedBy(func(task *interfaces.ExportTask) bool {
		return task.ResourceID == documentID && task.Format == interfaces.ExportFormatTXT
	})).Return(nil)

	// 执行测试
	req := &interfaces.ExportDocumentRequest{
		Format:     interfaces.ExportFormatTXT,
		IncludeMeta: false,
		Options:    nil,
	}

	task, err := service.(*ExportService).ExportDocument(context.Background(), documentID, projectID, userID, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, documentID, task.ResourceID)
	assert.Equal(t, "Test Chapter", task.ResourceTitle)
	assert.Equal(t, interfaces.ExportStatusPending, task.Status)
	assert.Equal(t, userID, task.CreatedBy)

	mockDocRepo.AssertExpectations(t)
	mockExportRepo.AssertExpectations(t)
}

// TestExportDocument_DocumentNotFound 测试文档不存在
func TestExportDocument_DocumentNotFound(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	documentID := "nonexistent"
	projectID := "proj123"
	userID := "user123"

	// 设置Mock期望 - 文档不存在
	mockDocRepo.On("FindByID", mock.Anything, documentID).Return(nil, errors.New("文档不存在"))

	// 执行测试
	req := &interfaces.ExportDocumentRequest{
		Format: interfaces.ExportFormatTXT,
	}

	task, err := service.(*ExportService).ExportDocument(context.Background(), documentID, projectID, userID, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "文档不存在")

	mockDocRepo.AssertExpectations(t)
}

// TestExportDocument_Forbidden 测试权限验证失败
func TestExportDocument_Forbidden(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	documentID := "doc123"
	projectID := "proj123"
	wrongProjectID := "wrongproj"
	userID := "user123"
	document := createTestDocument(documentID, wrongProjectID, "Test Chapter")

	// 设置Mock期望 - 项目ID不匹配
	mockDocRepo.On("FindByID", mock.Anything, documentID).Return(document, nil)

	// 执行测试
	req := &interfaces.ExportDocumentRequest{
		Format: interfaces.ExportFormatTXT,
	}

	task, err := service.(*ExportService).ExportDocument(context.Background(), documentID, projectID, userID, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "无权访问")

	mockDocRepo.AssertExpectations(t)
}

// TestExportDocument_CreateTaskFailed 测试创建导出任务失败
func TestExportDocument_CreateTaskFailed(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	documentID := "doc123"
	projectID := "proj123"
	userID := "user123"
	document := createTestDocument(documentID, projectID, "Test Chapter")

	// 设置Mock期望
	mockDocRepo.On("FindByID", mock.Anything, documentID).Return(document, nil)
	mockExportRepo.On("Create", mock.Anything, mock.AnythingOfType("*interfaces.ExportTask")).Return(errors.New("数据库错误"))

	// 执行测试
	req := &interfaces.ExportDocumentRequest{
		Format: interfaces.ExportFormatTXT,
	}

	task, err := service.(*ExportService).ExportDocument(context.Background(), documentID, projectID, userID, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "创建导出任务失败")

	mockDocRepo.AssertExpectations(t)
	mockExportRepo.AssertExpectations(t)
}

// ============ GetExportTask 测试 ============

// TestGetExportTask_Success 测试获取导出任务成功
func TestGetExportTask_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	taskID := "task123"
	task := createTestExportTask(taskID, "doc123", "Test Chapter", interfaces.ExportFormatTXT, "user123")
	task.Status = interfaces.ExportStatusCompleted

	// 设置Mock期望
	mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)

	// 执行测试
	result, err := service.(*ExportService).GetExportTask(context.Background(), taskID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, taskID, result.ID)
	assert.Equal(t, interfaces.ExportStatusCompleted, result.Status)

	mockExportRepo.AssertExpectations(t)
}

// TestGetExportTask_NotFound 测试获取不存在的导出任务
func TestGetExportTask_NotFound(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	taskID := "nonexistent"

	// 设置Mock期望 - 任务不存在
	mockExportRepo.On("FindByID", mock.Anything, taskID).Return(nil, errors.New("任务不存在"))

	// 执行测试
	result, err := service.(*ExportService).GetExportTask(context.Background(), taskID)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "导出任务不存在")

	mockExportRepo.AssertExpectations(t)
}

// ============ DownloadExportFile 测试 ============

// TestDownloadExportFile_Success 测试下载导出文件成功
func TestDownloadExportFile_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	taskID := "task123"
	task := createTestExportTask(taskID, "doc123", "Test Chapter", interfaces.ExportFormatTXT, "user123")
	task.Status = interfaces.ExportStatusCompleted
	task.FileURL = "/exports/test.txt"
	task.FileSize = 1024

	signedURL := "https://storage.example.com/exports/test.txt?signature=xxx"

	// 设置Mock期望
	mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)
	mockFileStorage.On("GetSignedURL", mock.Anything, task.FileURL, 1*time.Hour).Return(signedURL, nil)

	// 执行测试
	result, err := service.(*ExportService).DownloadExportFile(context.Background(), taskID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Chapter.txt", result.Filename)
	assert.Equal(t, signedURL, result.URL)
	assert.Equal(t, "text/plain", result.MimeType)
	assert.Equal(t, int64(1024), result.FileSize)

	mockExportRepo.AssertExpectations(t)
	mockFileStorage.AssertExpectations(t)
}

// TestDownloadExportFile_TaskNotCompleted 测试任务未完成时下载
func TestDownloadExportFile_TaskNotCompleted(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	taskID := "task123"
	task := createTestExportTask(taskID, "doc123", "Test Chapter", interfaces.ExportFormatTXT, "user123")
	task.Status = interfaces.ExportStatusProcessing

	// 设置Mock期望
	mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)

	// 执行测试
	result, err := service.(*ExportService).DownloadExportFile(context.Background(), taskID)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "导出任务未完成")

	mockExportRepo.AssertExpectations(t)
}

// TestDownloadExportFile_FileExpired 测试文件已过期
func TestDownloadExportFile_FileExpired(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	taskID := "task123"
	task := createTestExportTask(taskID, "doc123", "Test Chapter", interfaces.ExportFormatTXT, "user123")
	task.Status = interfaces.ExportStatusCompleted
	task.ExpiresAt = time.Now().Add(-1 * time.Hour) // 已过期

	// 设置Mock期望
	mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)

	// 执行测试
	result, err := service.(*ExportService).DownloadExportFile(context.Background(), taskID)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "导出文件已过期")

	mockExportRepo.AssertExpectations(t)
}

// ============ ListExportTasks 测试 ============

// TestListExportTasks_Success 测试列出导出任务成功
func TestListExportTasks_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	projectID := "proj123"
	page := 1
	pageSize := 20

	tasks := []*interfaces.ExportTask{
		createTestExportTask("task1", "doc1", "Chapter 1", interfaces.ExportFormatTXT, "user1"),
		createTestExportTask("task2", "doc2", "Chapter 2", interfaces.ExportFormatMD, "user1"),
	}

	// 设置Mock期望
	mockExportRepo.On("FindByProjectID", mock.Anything, projectID, page, pageSize).Return(tasks, int64(2), nil)

	// 执行测试
	result, total, err := service.(*ExportService).ListExportTasks(context.Background(), projectID, page, pageSize)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)

	mockExportRepo.AssertExpectations(t)
}

// TestListExportTasks_InvalidPageParams 测试无效的分页参数
func TestListExportTasks_InvalidPageParams(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	projectID := "proj123"

	tests := []struct {
		name      string
		page      int
		pageSize  int
		expectPage int
		expectPageSize int
	}{
		{"页码为0", 0, 20, 1, 20},
		{"页码为负", -1, 20, 1, 20},
		{"页大小为0", 1, 0, 1, 20},
		{"页大小过大", 1, 200, 1, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置Mock期望 - 使用期望的参数
			mockExportRepo.On("FindByProjectID", mock.Anything, projectID, tt.expectPage, tt.expectPageSize).
				Return([]*interfaces.ExportTask{}, int64(0), nil).Once()

			// 执行测试
			result, total, err := service.(*ExportService).ListExportTasks(context.Background(), projectID, tt.page, tt.pageSize)

			// 验证结果
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, int64(0), total)

			mockExportRepo.AssertExpectations(t)

			// 重置Mock
			mockExportRepo.ExpectedCalls = nil
		})
	}
}

// ============ DeleteExportTask 测试 ============

// TestDeleteExportTask_Success 测试删除导出任务成功
func TestDeleteExportTask_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	taskID := "task123"
	userID := "user123"
	task := createTestExportTask(taskID, "doc123", "Test Chapter", interfaces.ExportFormatTXT, userID)
	task.FileURL = "/exports/test.txt"

	// 设置Mock期望
	mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)
	mockFileStorage.On("Delete", mock.Anything, task.FileURL).Return(nil)
	mockExportRepo.On("Delete", mock.Anything, taskID).Return(nil)

	// 执行测试
	err := service.(*ExportService).DeleteExportTask(context.Background(), taskID, userID)

	// 验证结果
	assert.NoError(t, err)

	mockExportRepo.AssertExpectations(t)
	mockFileStorage.AssertExpectations(t)
}

// TestDeleteExportTask_Forbidden 测试无权限删除导出任务
func TestDeleteExportTask_Forbidden(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	taskID := "task123"
	userID := "user123"
	otherUserID := "otheruser"
	task := createTestExportTask(taskID, "doc123", "Test Chapter", interfaces.ExportFormatTXT, otherUserID)

	// 设置Mock期望 - 创建者不是当前用户
	mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)

	// 执行测试
	err := service.(*ExportService).DeleteExportTask(context.Background(), taskID, userID)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无权删除")

	mockExportRepo.AssertExpectations(t)
}

// TestDeleteExportTask_NotFound 测试删除不存在的导出任务
func TestDeleteExportTask_NotFound(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	taskID := "nonexistent"
	userID := "user123"

	// 设置Mock期望 - 任务不存在
	mockExportRepo.On("FindByID", mock.Anything, taskID).Return(nil, errors.New("任务不存在"))

	// 执行测试
	err := service.(*ExportService).DeleteExportTask(context.Background(), taskID, userID)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "导出任务不存在")

	mockExportRepo.AssertExpectations(t)
}

// ============ ExportProject 测试 ============

// TestExportProject_Success 测试导出项目成功
func TestExportProject_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	projectID := "proj123"
	userID := "user123"
	project := createTestProject(projectID, userID, "Test Project")

	// 设置Mock期望
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	mockExportRepo.On("Create", mock.Anything, mock.MatchedBy(func(task *interfaces.ExportTask) bool {
		return task.ResourceID == projectID && task.Type == interfaces.ExportTypeProject
	})).Return(nil)

	// 执行测试
	req := &interfaces.ExportProjectRequest{
		IncludeDocuments: true,
		DocumentFormats:  interfaces.ExportFormatTXT,
	}

	task, err := service.(*ExportService).ExportProject(context.Background(), projectID, userID, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, projectID, task.ResourceID)
	assert.Equal(t, "Test Project", task.ResourceTitle)
	assert.Equal(t, interfaces.ExportTypeProject, task.Type)
	assert.Equal(t, interfaces.ExportFormatZIP, task.Format)

	mockProjectRepo.AssertExpectations(t)
	mockExportRepo.AssertExpectations(t)
}

// TestExportProject_ProjectNotFound 测试项目不存在
func TestExportProject_ProjectNotFound(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	projectID := "nonexistent"
	userID := "user123"

	// 设置Mock期望 - 项目不存在
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(nil, errors.New("项目不存在"))

	// 执行测试
	req := &interfaces.ExportProjectRequest{
		IncludeDocuments: true,
	}

	task, err := service.(*ExportService).ExportProject(context.Background(), projectID, userID, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "项目不存在")

	mockProjectRepo.AssertExpectations(t)
}

// ============ CancelExportTask 测试 ============

// TestCancelExportTask_Success 测试取消导出任务成功
func TestCancelExportTask_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	taskID := "task123"
	userID := "user123"
	task := createTestExportTask(taskID, "doc123", "Test Chapter", interfaces.ExportFormatTXT, userID)
	task.Status = interfaces.ExportStatusPending

	// 设置Mock期望
	mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)
	mockExportRepo.On("Update", mock.Anything, mock.MatchedBy(func(t *interfaces.ExportTask) bool {
		return t.Status == interfaces.ExportStatusCancelled
	})).Return(nil)

	// 执行测试
	err := service.(*ExportService).CancelExportTask(context.Background(), taskID, userID)

	// 验证结果
	assert.NoError(t, err)

	mockExportRepo.AssertExpectations(t)
}

// TestCancelExportTask_Forbidden 测试无权限取消导出任务
func TestCancelExportTask_Forbidden(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	taskID := "task123"
	userID := "user123"
	otherUserID := "otheruser"
	task := createTestExportTask(taskID, "doc123", "Test Chapter", interfaces.ExportFormatTXT, otherUserID)
	task.Status = interfaces.ExportStatusPending

	// 设置Mock期望 - 创建者不是当前用户
	mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)

	// 执行测试
	err := service.(*ExportService).CancelExportTask(context.Background(), taskID, userID)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无权取消")

	mockExportRepo.AssertExpectations(t)
}

// TestCancelExportTask_InvalidStatus 测试取消不允许取消状态的任务
func TestCancelExportTask_InvalidStatus(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockContentRepo := new(mocks.MockDocumentContentRepository)
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockExportRepo := new(mocks.MockExportTaskRepository)
	mockFileStorage := new(mocks.MockFileStorage)

	service := NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	taskID := "task123"
	userID := "user123"
	task := createTestExportTask(taskID, "doc123", "Test Chapter", interfaces.ExportFormatTXT, userID)
	task.Status = interfaces.ExportStatusCompleted // 已完成，不能取消

	// 设置Mock期望
	mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)

	// 执行测试
	err := service.(*ExportService).CancelExportTask(context.Background(), taskID, userID)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "任务状态不允许取消")

	mockExportRepo.AssertExpectations(t)
}
