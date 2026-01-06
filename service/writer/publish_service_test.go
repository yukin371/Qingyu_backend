package writer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/service/interfaces"
	"Qingyu_backend/service/writer/mocks"
)

// ============ 测试辅助函数 ============

// createTestPublicationRecord 创建测试发布记录
func createTestPublicationRecord(id, resourceID, resourceTitle, bookstoreID, createdBy string) *interfaces.PublicationRecord {
	return &interfaces.PublicationRecord{
		ID:            id,
		Type:          "project",
		ResourceID:    resourceID,
		ResourceTitle: resourceTitle,
		BookstoreID:   bookstoreID,
		Status:        interfaces.PublicationStatusPending,
		CreatedBy:     createdBy,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// createPublishProjectRequest 创建发布项目请求
func createPublishProjectRequest(bookstoreID, categoryID string) *interfaces.PublishProjectRequest {
	price := 9.99
	return &interfaces.PublishProjectRequest{
		BookstoreID:    bookstoreID,
		CategoryID:     categoryID,
		Tags:           []string{"小说", "玄幻"},
		Description:    "测试项目描述",
		PublishType:    interfaces.PublishTypeSerial,
		Price:          &price,
		FreeChapters:   10,
		AuthorNote:     "作者的话",
		EnableComment:  true,
		EnableShare:    true,
	}
}

// createPublishDocumentRequest 创建发布文档请求
func createPublishDocumentRequest(title string, number int) *interfaces.PublishDocumentRequest {
	return &interfaces.PublishDocumentRequest{
		ChapterTitle:  title,
		ChapterNumber: number,
		IsFree:        false,
		AuthorNote:    "章节作者的话",
	}
}

// ============ PublishProject 测试 ============

// TestPublishProject_Success 测试发布项目成功
func TestPublishProject_Success(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	userID := "user123"
	project := createTestProject(projectID, userID, "Test Project")

	req := createPublishProjectRequest("bookstore123", "category123")

	// 设置Mock期望
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(nil, nil)
	mockPublicationRepo.On("Create", mock.Anything, mock.MatchedBy(func(record *interfaces.PublicationRecord) bool {
		return record.ResourceID == projectID && record.BookstoreID == req.BookstoreID
	})).Return(nil)

	// 执行测试
	record, err := service.(*PublishService).PublishProject(context.Background(), projectID, userID, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, projectID, record.ResourceID)
	assert.Equal(t, "Test Project", record.ResourceTitle)
	assert.Equal(t, req.BookstoreID, record.BookstoreID)
	assert.Equal(t, interfaces.PublicationStatusPending, record.Status)
	assert.Equal(t, userID, record.CreatedBy)

	mockProjectRepo.AssertExpectations(t)
	mockPublicationRepo.AssertExpectations(t)
}

// TestPublishProject_ProjectNotFound 测试项目不存在
func TestPublishProject_ProjectNotFound(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "nonexistent"
	userID := "user123"
	req := createPublishProjectRequest("bookstore123", "category123")

	// 设置Mock期望 - 项目不存在
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(nil, errors.New("项目不存在"))

	// 执行测试
	record, err := service.(*PublishService).PublishProject(context.Background(), projectID, userID, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Contains(t, err.Error(), "项目不存在")

	mockProjectRepo.AssertExpectations(t)
}

// TestPublishProject_Forbidden 测试无权限发布项目
func TestPublishProject_Forbidden(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	userID := "user123"
	otherUserID := "otheruser"
	project := createTestProject(projectID, otherUserID, "Test Project")

	req := createPublishProjectRequest("bookstore123", "category123")

	// 设置Mock期望 - 用户不是所有者
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)

	// 执行测试
	record, err := service.(*PublishService).PublishProject(context.Background(), projectID, userID, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Contains(t, err.Error(), "无权发布")

	mockProjectRepo.AssertExpectations(t)
}

// TestPublishProject_AlreadyPublished 测试项目已发布
func TestPublishProject_AlreadyPublished(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	userID := "user123"
	project := createTestProject(projectID, userID, "Test Project")

	existingRecord := createTestPublicationRecord("record123", projectID, "Test Project", "bookstore123", userID)
	existingRecord.Status = interfaces.PublicationStatusPublished

	req := createPublishProjectRequest("bookstore123", "category123")

	// 设置Mock期望
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(existingRecord, nil)

	// 执行测试
	record, err := service.(*PublishService).PublishProject(context.Background(), projectID, userID, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Contains(t, err.Error(), "项目已发布")

	mockProjectRepo.AssertExpectations(t)
	mockPublicationRepo.AssertExpectations(t)
}

// TestPublishProject_CreateRecordFailed 测试创建发布记录失败
func TestPublishProject_CreateRecordFailed(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	userID := "user123"
	project := createTestProject(projectID, userID, "Test Project")

	req := createPublishProjectRequest("bookstore123", "category123")

	// 设置Mock期望
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(nil, nil)
	mockPublicationRepo.On("Create", mock.Anything, mock.AnythingOfType("*interfaces.PublicationRecord")).
		Return(errors.New("数据库错误"))

	// 执行测试
	record, err := service.(*PublishService).PublishProject(context.Background(), projectID, userID, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Contains(t, err.Error(), "创建发布记录失败")

	mockProjectRepo.AssertExpectations(t)
	mockPublicationRepo.AssertExpectations(t)
}

// ============ UnpublishProject 测试 ============

// TestUnpublishProject_Success 测试取消发布项目成功
func TestUnpublishProject_Success(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	userID := "user123"
	bookstoreID := "bookstore123"

	record := createTestPublicationRecord("record123", projectID, "Test Project", bookstoreID, userID)
	record.Status = interfaces.PublicationStatusPublished

	// 设置Mock期望
	mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(record, nil)
	mockBookstoreClient.On("UnpublishProject", mock.Anything, projectID, bookstoreID).Return(nil)
	mockPublicationRepo.On("Update", mock.Anything, mock.MatchedBy(func(r *interfaces.PublicationRecord) bool {
		return r.Status == interfaces.PublicationStatusUnpublished
	})).Return(nil)
	mockEventBus.On("PublishAsync", mock.Anything, mock.Anything).Return(nil)

	// 执行测试
	err := service.(*PublishService).UnpublishProject(context.Background(), projectID, userID)

	// 验证结果
	assert.NoError(t, err)

	mockPublicationRepo.AssertExpectations(t)
	mockBookstoreClient.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestUnpublishProject_NotPublished 测试项目未发布
func TestUnpublishProject_NotPublished(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	userID := "user123"

	// 设置Mock期望 - 没有发布记录
	mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(nil, errors.New("未发布"))

	// 执行测试
	err := service.(*PublishService).UnpublishProject(context.Background(), projectID, userID)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "项目未发布")

	mockPublicationRepo.AssertExpectations(t)
}

// TestUnpublishProject_Forbidden 测试无权限取消发布
func TestUnpublishProject_Forbidden(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	userID := "user123"
	otherUserID := "otheruser"
	bookstoreID := "bookstore123"

	record := createTestPublicationRecord("record123", projectID, "Test Project", bookstoreID, otherUserID)

	// 设置Mock期望 - 创建者不是当前用户
	mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(record, nil)

	// 执行测试
	err := service.(*PublishService).UnpublishProject(context.Background(), projectID, userID)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无权取消")

	mockPublicationRepo.AssertExpectations(t)
}

// ============ GetProjectPublicationStatus 测试 ============

// TestGetProjectPublicationStatus_Success 测试获取项目发布状态成功
func TestGetProjectPublicationStatus_Success(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	project := createTestProject(projectID, "user123", "Test Project")

	bookstoreID := "bookstore123"
	publishTime := time.Now()
	record := createTestPublicationRecord("record123", projectID, "Test Project", bookstoreID, "user123")
	record.Status = interfaces.PublicationStatusPublished
	record.PublishTime = &publishTime

	documents := []*writer.Document{
		createTestDocument("doc1", projectID, "Chapter 1"),
		createTestDocument("doc2", projectID, "Chapter 2"),
	}

	statistics := &interfaces.PublicationStatistics{
		TotalViews:    1000,
		TotalLikes:    100,
		TotalComments: 50,
	}

	// 设置Mock期望
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(record, nil)
	mockBookstoreClient.On("GetStatistics", mock.Anything, projectID, bookstoreID).Return(statistics, nil)
	mockDocumentRepo.On("FindByProjectID", mock.Anything, projectID).Return(documents, nil)

	// 执行测试
	status, err := service.(*PublishService).GetProjectPublicationStatus(context.Background(), projectID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, projectID, status.ProjectID)
	assert.Equal(t, "Test Project", status.ProjectTitle)
	assert.True(t, status.IsPublished)
	assert.Equal(t, bookstoreID, status.BookstoreID)
	assert.Equal(t, 2, status.TotalChapters)
	assert.Equal(t, 2, status.PublishedChapters)
	assert.Equal(t, int64(1000), status.Statistics.TotalViews)

	mockProjectRepo.AssertExpectations(t)
	mockPublicationRepo.AssertExpectations(t)
	mockBookstoreClient.AssertExpectations(t)
	mockDocumentRepo.AssertExpectations(t)
}

// TestGetProjectPublicationStatus_NotPublished 测试项目未发布时的状态
func TestGetProjectPublicationStatus_NotPublished(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	project := createTestProject(projectID, "user123", "Test Project")

	documents := []*writer.Document{
		createTestDocument("doc1", projectID, "Chapter 1"),
	}

	// 设置Mock期望
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(nil, nil)
	mockDocumentRepo.On("FindByProjectID", mock.Anything, projectID).Return(documents, nil)

	// 执行测试
	status, err := service.(*PublishService).GetProjectPublicationStatus(context.Background(), projectID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.False(t, status.IsPublished)
	assert.Equal(t, 1, status.TotalChapters)

	mockProjectRepo.AssertExpectations(t)
	mockPublicationRepo.AssertExpectations(t)
	mockDocumentRepo.AssertExpectations(t)
}

// ============ PublishDocument 测试 ============

// TestPublishDocument_Success 测试发布文档成功
func TestPublishDocument_Success(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	documentID := "doc123"
	projectID := "proj123"
	userID := "user123"

	document := createTestDocument(documentID, projectID, "Test Chapter")
	project := createTestProject(projectID, userID, "Test Project")

	req := createPublishDocumentRequest("第一章", 1)

	// 设置Mock期望
	mockDocumentRepo.On("FindByID", mock.Anything, documentID).Return(document, nil)
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	mockPublicationRepo.On("FindByResourceID", mock.Anything, documentID).Return(nil, nil)
	mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(nil, nil)
	mockPublicationRepo.On("Create", mock.Anything, mock.MatchedBy(func(record *interfaces.PublicationRecord) bool {
		return record.ResourceID == documentID && record.Type == "document"
	})).Return(nil)

	// 执行测试
	record, err := service.(*PublishService).PublishDocument(context.Background(), documentID, projectID, userID, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, documentID, record.ResourceID)
	assert.Equal(t, req.ChapterTitle, record.ResourceTitle)
	assert.Equal(t, "document", record.Type)
	assert.Equal(t, interfaces.PublicationStatusPending, record.Status)

	mockDocumentRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
	mockPublicationRepo.AssertExpectations(t)
}

// TestPublishDocument_DocumentNotFound 测试文档不存在
func TestPublishDocument_DocumentNotFound(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	documentID := "nonexistent"
	projectID := "proj123"
	userID := "user123"

	req := createPublishDocumentRequest("第一章", 1)

	// 设置Mock期望 - 文档不存在
	mockDocumentRepo.On("FindByID", mock.Anything, documentID).Return(nil, errors.New("文档不存在"))

	// 执行测试
	record, err := service.(*PublishService).PublishDocument(context.Background(), documentID, projectID, userID, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Contains(t, err.Error(), "文档不存在")

	mockDocumentRepo.AssertExpectations(t)
}

// TestPublishDocument_Forbidden 测试无权限发布文档
func TestPublishDocument_Forbidden(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	documentID := "doc123"
	projectID := "proj123"
	wrongProjectID := "wrongproj"
	userID := "user123"

	document := createTestDocument(documentID, wrongProjectID, "Test Chapter")

	req := createPublishDocumentRequest("第一章", 1)

	// 设置Mock期望 - 项目ID不匹配
	mockDocumentRepo.On("FindByID", mock.Anything, documentID).Return(document, nil)

	// 执行测试
	record, err := service.(*PublishService).PublishDocument(context.Background(), documentID, projectID, userID, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Contains(t, err.Error(), "无权访问")

	mockDocumentRepo.AssertExpectations(t)
}

// TestPublishDocument_AlreadyPublished 测试文档已发布
func TestPublishDocument_AlreadyPublished(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	documentID := "doc123"
	projectID := "proj123"
	userID := "user123"

	document := createTestDocument(documentID, projectID, "Test Chapter")
	project := createTestProject(projectID, userID, "Test Project")

	existingRecord := createTestPublicationRecord("record123", documentID, "Test Chapter", "bookstore123", userID)
	existingRecord.Type = "document"
	existingRecord.Status = interfaces.PublicationStatusPublished

	req := createPublishDocumentRequest("第一章", 1)

	// 设置Mock期望
	mockDocumentRepo.On("FindByID", mock.Anything, documentID).Return(document, nil)
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	mockPublicationRepo.On("FindByResourceID", mock.Anything, documentID).Return(existingRecord, nil)

	// 执行测试
	record, err := service.(*PublishService).PublishDocument(context.Background(), documentID, projectID, userID, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Contains(t, err.Error(), "文档已发布")

	mockDocumentRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
	mockPublicationRepo.AssertExpectations(t)
}

// ============ BatchPublishDocuments 测试 ============

// TestBatchPublishDocuments_Success 测试批量发布文档成功
func TestBatchPublishDocuments_Success(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	userID := "user123"

	project := createTestProject(projectID, userID, "Test Project")

	documentIDs := []string{"doc1", "doc2", "doc3"}

	doc1 := createTestDocument("doc1", projectID, "Chapter 1")
	doc2 := createTestDocument("doc2", projectID, "Chapter 2")
	doc3 := createTestDocument("doc3", projectID, "Chapter 3")

	req := &interfaces.BatchPublishDocumentsRequest{
		DocumentIDs:   documentIDs,
		AutoNumbering: true,
		StartNumber:   1,
		IsFree:        false,
	}

	// 设置Mock期望 - 每个文档的完整流程
	for i, docID := range documentIDs {
		var doc *writer.Document
		switch i {
		case 0:
			doc = doc1
		case 1:
			doc = doc2
		case 2:
			doc = doc3
		}

		mockDocumentRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
		mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil).Once()
		mockPublicationRepo.On("FindByResourceID", mock.Anything, docID).Return(nil, nil).Once()
		mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(nil, nil).Once()
		mockPublicationRepo.On("Create", mock.Anything, mock.AnythingOfType("*interfaces.PublicationRecord")).Return(nil).Once()
	}

	// 执行测试
	result, err := service.(*PublishService).BatchPublishDocuments(context.Background(), projectID, userID, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailCount)
	assert.Equal(t, 3, len(result.Results))

	for i, item := range result.Results {
		assert.True(t, item.Success)
		assert.Equal(t, documentIDs[i], item.DocumentID)
		assert.NotEmpty(t, item.RecordID)
	}

	mockDocumentRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
	mockPublicationRepo.AssertExpectations(t)
}

// TestBatchPublishDocuments_PartialFailure 测试批量发布部分失败
func TestBatchPublishDocuments_PartialFailure(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	userID := "user123"

	project := createTestProject(projectID, userID, "Test Project")

	documentIDs := []string{"doc1", "doc2"}

	doc1 := createTestDocument("doc1", projectID, "Chapter 1")

	req := &interfaces.BatchPublishDocumentsRequest{
		DocumentIDs:   documentIDs,
		AutoNumbering: true,
		StartNumber:   1,
		IsFree:        false,
	}

	// 设置Mock期望 - 第一个成功，第二个失败
	mockDocumentRepo.On("FindByID", mock.Anything, "doc1").Return(doc1, nil).Once()
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil).Once()
	mockPublicationRepo.On("FindByResourceID", mock.Anything, "doc1").Return(nil, nil).Once()
	mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(nil, nil).Once()
	mockPublicationRepo.On("Create", mock.Anything, mock.AnythingOfType("*interfaces.PublicationRecord")).Return(nil).Once()

	// 第二个文档不存在
	mockDocumentRepo.On("FindByID", mock.Anything, "doc2").Return(nil, errors.New("文档不存在")).Once()

	// 执行测试
	result, err := service.(*PublishService).BatchPublishDocuments(context.Background(), projectID, userID, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailCount)
	assert.Equal(t, 2, len(result.Results))

	assert.True(t, result.Results[0].Success)
	assert.False(t, result.Results[1].Success)
	assert.Contains(t, result.Results[1].Error, "文档不存在")

	mockDocumentRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
	mockPublicationRepo.AssertExpectations(t)
}

// TestBatchPublishDocuments_Forbidden 测试无权限批量发布
func TestBatchPublishDocuments_Forbidden(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	userID := "user123"
	otherUserID := "otheruser"

	project := createTestProject(projectID, otherUserID, "Test Project")

	documentIDs := []string{"doc1"}

	req := &interfaces.BatchPublishDocumentsRequest{
		DocumentIDs:   documentIDs,
		AutoNumbering: true,
		StartNumber:   1,
		IsFree:        false,
	}

	// 设置Mock期望 - 用户不是所有者
	mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)

	// 执行测试
	result, err := service.(*PublishService).BatchPublishDocuments(context.Background(), projectID, userID, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "无权批量发布")

	mockProjectRepo.AssertExpectations(t)
}

// ============ GetPublicationRecords 测试 ============

// TestGetPublicationRecords_Success 测试获取发布记录列表成功
func TestGetPublicationRecords_Success(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	projectID := "proj123"
	page := 1
	pageSize := 20

	records := []*interfaces.PublicationRecord{
		createTestPublicationRecord("record1", "doc1", "Chapter 1", "bookstore123", "user123"),
		createTestPublicationRecord("record2", "doc2", "Chapter 2", "bookstore123", "user123"),
	}

	// 设置Mock期望
	mockPublicationRepo.On("FindByProjectID", mock.Anything, projectID, page, pageSize).Return(records, int64(2), nil)

	// 执行测试
	result, total, err := service.(*PublishService).GetPublicationRecords(context.Background(), projectID, page, pageSize)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)

	mockPublicationRepo.AssertExpectations(t)
}

// TestGetPublicationRecord_Success 测试获取发布记录详情成功
func TestGetPublicationRecord_Success(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	recordID := "record123"
	record := createTestPublicationRecord(recordID, "doc1", "Chapter 1", "bookstore123", "user123")

	// 设置Mock期望
	mockPublicationRepo.On("FindByID", mock.Anything, recordID).Return(record, nil)

	// 执行测试
	result, err := service.(*PublishService).GetPublicationRecord(context.Background(), recordID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, recordID, result.ID)

	mockPublicationRepo.AssertExpectations(t)
}

// TestGetPublicationRecord_NotFound 测试获取不存在的发布记录
func TestGetPublicationRecord_NotFound(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepository)
	mockDocumentRepo := new(mocks.MockDocumentRepository)
	mockPublicationRepo := new(mocks.MockPublicationRepository)
	mockBookstoreClient := new(MockBookstoreClient)
	mockEventBus := new(mocks.MockEventBus)

	service := NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	recordID := "nonexistent"

	// 设置Mock期望 - 记录不存在
	mockPublicationRepo.On("FindByID", mock.Anything, recordID).Return(nil, errors.New("记录不存在"))

	// 执行测试
	result, err := service.(*PublishService).GetPublicationRecord(context.Background(), recordID)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "发布记录不存在")

	mockPublicationRepo.AssertExpectations(t)
}
