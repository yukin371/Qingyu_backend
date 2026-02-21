package writer

import (
	"archive/zip"
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/writer"
	baseModel "Qingyu_backend/models/writer/base"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// MockDocumentRepositoryForExport Mock文档仓储
type MockDocumentRepositoryForExport struct {
	mock.Mock
}

func (m *MockDocumentRepositoryForExport) FindByID(ctx context.Context, id string) (*writer.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Document), args.Error(1)
}

func (m *MockDocumentRepositoryForExport) FindByProjectID(ctx context.Context, projectID string) ([]*writer.Document, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Document), args.Error(1)
}

// MockDocumentContentRepositoryForExport Mock文档内容仓储
type MockDocumentContentRepositoryForExport struct {
	mock.Mock
}

func (m *MockDocumentContentRepositoryForExport) FindByID(ctx context.Context, id string) (*writer.DocumentContent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.DocumentContent), args.Error(1)
}

// MockProjectRepositoryForExport Mock项目仓储
type MockProjectRepositoryForExport struct {
	mock.Mock
}

func (m *MockProjectRepositoryForExport) FindByID(ctx context.Context, id string) (*writer.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Project), args.Error(1)
}

// MockExportTaskRepositoryForExport Mock导出任务仓储
type MockExportTaskRepositoryForExport struct {
	mock.Mock
}

func (m *MockExportTaskRepositoryForExport) Create(ctx context.Context, task *serviceInterfaces.ExportTask) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockExportTaskRepositoryForExport) FindByID(ctx context.Context, id string) (*serviceInterfaces.ExportTask, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*serviceInterfaces.ExportTask), args.Error(1)
}

func (m *MockExportTaskRepositoryForExport) FindByProjectID(ctx context.Context, projectID string, page, pageSize int) ([]*serviceInterfaces.ExportTask, int64, error) {
	args := m.Called(ctx, projectID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*serviceInterfaces.ExportTask), args.Get(1).(int64), args.Error(2)
}

func (m *MockExportTaskRepositoryForExport) Update(ctx context.Context, task *serviceInterfaces.ExportTask) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockExportTaskRepositoryForExport) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockExportTaskRepositoryForExport) FindByUser(ctx context.Context, userID string, page, pageSize int) ([]*serviceInterfaces.ExportTask, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*serviceInterfaces.ExportTask), args.Get(1).(int64), args.Error(2)
}

// TestExportService_ExportProjectAsZip_Success 测试成功导出项目为ZIP
func TestExportService_ExportProjectAsZip_Success(t *testing.T) {
	// Given
	mockDocRepo := new(MockDocumentRepositoryForExport)
	mockContentRepo := new(MockDocumentContentRepositoryForExport)
	mockProjectRepo := new(MockProjectRepositoryForExport)
	mockTaskRepo := new(MockExportTaskRepositoryForExport)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID().Hex()

	// 创建项目
	project := &writer.Project{}
	project.ID = projectID
	project.Title = "测试项目"
	project.Summary = "这是一个测试项目"
	project.Status = writer.StatusDraft
	project.Visibility = writer.VisibilityPrivate
	project.WritingType = "novel"
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	// 创建测试文档
	doc1ID := primitive.NewObjectID()
	doc2ID := primitive.NewObjectID()

	doc1 := &writer.Document{
		ProjectID: projectID,
		Title:     "第一章",
		Type:      "chapter",
		Level:     0,
		Order:     1,
		StableRef: "ref-1",
		OrderKey:  "a1",
	}
	doc1.ID = doc1ID

	doc2 := &writer.Document{
		ProjectID: projectID,
		Title:     "第二章",
		Type:      "chapter",
		Level:     0,
		Order:     2,
		StableRef: "ref-2",
		OrderKey:  "a2",
	}
	doc2.ID = doc2ID

	documents := []*writer.Document{doc1, doc2}

	mockProjectRepo.On("FindByID", mock.Anything, projectID.Hex()).Return(project, nil)
	mockDocRepo.On("FindByProjectID", mock.Anything, projectID.Hex()).Return(documents, nil)

	// Mock文档内容
	content1 := &writer.DocumentContent{
		Content: "这是第一章的内容",
	}
	content1.ID = doc1ID

	content2 := &writer.DocumentContent{
		Content: "这是第二章的内容",
	}
	content2.ID = doc2ID

	mockContentRepo.On("FindByID", mock.Anything, doc1ID.Hex()).Return(content1, nil)
	mockContentRepo.On("FindByID", mock.Anything, doc2ID.Hex()).Return(content2, nil)

	service := NewExportService(mockDocRepo, mockContentRepo, mockProjectRepo, mockTaskRepo, nil)

	// When
	zipData, err := service.ExportProjectAsZip(context.Background(), projectID.Hex(), userID)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, zipData)
	assert.True(t, len(zipData) > 0)

	// 验证是有效的ZIP文件
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	assert.NoError(t, err)
	assert.True(t, len(reader.File) > 0)

	mockProjectRepo.AssertExpectations(t)
	mockDocRepo.AssertExpectations(t)
	mockContentRepo.AssertExpectations(t)
}

// TestExportService_ExportProjectAsZip_ProjectNotFound 测试项目不存在
func TestExportService_ExportProjectAsZip_ProjectNotFound(t *testing.T) {
	// Given
	mockDocRepo := new(MockDocumentRepositoryForExport)
	mockContentRepo := new(MockDocumentContentRepositoryForExport)
	mockProjectRepo := new(MockProjectRepositoryForExport)
	mockTaskRepo := new(MockExportTaskRepositoryForExport)

	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()

	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(nil, assert.AnError)

	service := NewExportService(mockDocRepo, mockContentRepo, mockProjectRepo, mockTaskRepo, nil)

	// When
	zipData, err := service.ExportProjectAsZip(context.Background(), projectID, userID)

	// Then
	assert.Error(t, err)
	assert.Nil(t, zipData)

	mockProjectRepo.AssertExpectations(t)
}

// TestExportService_ImportProject_Success 测试成功导入项目
func TestExportService_ImportProject_Success(t *testing.T) {
	// Given
	mockDocRepo := new(MockDocumentRepositoryForExport)
	mockContentRepo := new(MockDocumentContentRepositoryForExport)
	mockProjectRepo := new(MockProjectRepositoryForExport)
	mockTaskRepo := new(MockExportTaskRepositoryForExport)

	userID := primitive.NewObjectID().Hex()

	// 创建测试用的 ZIP 数据
	zipData := createTestZipData(t, "导入测试项目", map[string]string{
		"第一章.txt": "这是第一章的内容",
		"第二章.txt": "这是第二章的内容",
	})

	service := NewExportService(mockDocRepo, mockContentRepo, mockProjectRepo, mockTaskRepo, nil)

	// When
	result, err := service.ImportProject(context.Background(), userID, zipData)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ProjectID)
	assert.Equal(t, "导入测试项目", result.Title)
	assert.GreaterOrEqual(t, result.DocumentCount, 2)
}

// createTestZipData 创建测试用的ZIP数据
func createTestZipData(t *testing.T, projectName string, files map[string]string) []byte {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for filename, content := range files {
		w, err := zipWriter.Create(projectName + "/" + filename)
		assert.NoError(t, err)
		_, err = w.Write([]byte(content))
		assert.NoError(t, err)
	}

	err := zipWriter.Close()
	assert.NoError(t, err)

	return buf.Bytes()
}

// 确保Mock实现接口
var _ DocumentRepository = (*MockDocumentRepositoryForExport)(nil)
var _ DocumentContentRepository = (*MockDocumentContentRepositoryForExport)(nil)
var _ ProjectRepository = (*MockProjectRepositoryForExport)(nil)
var _ ExportTaskRepository = (*MockExportTaskRepositoryForExport)(nil)

// 确保ID可以被设置
var _ = baseModel.IdentifiedEntity{}
