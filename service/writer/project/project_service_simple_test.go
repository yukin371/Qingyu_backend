package project

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/models/writer"
	writerBase "Qingyu_backend/models/writer/base"
	serviceBase "Qingyu_backend/service/base"
	base "Qingyu_backend/repository/interfaces/infrastructure"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockProjectRepository Mock项目仓储
type MockProjectRepository struct {
	mock.Mock
}

// CRUDRepository接口实现
func (m *MockProjectRepository) Create(ctx context.Context, project *writer.Project) error {
	args := m.Called(ctx, project)
	// 如果调用成功，为项目生成ID
	if args.Error(0) == nil && project.ID.IsZero() {
		project.ID = primitive.NewObjectID()
	}
	return args.Error(0)
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id string) (*writer.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepository) List(ctx context.Context, filter base.Filter) ([]*writer.Project, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// HealthRepository接口实现
func (m *MockProjectRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ProjectRepository特定方法实现
func (m *MockProjectRepository) GetListByOwnerID(ctx context.Context, ownerID string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, ownerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByOwnerAndStatus(ctx context.Context, ownerID, status string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, ownerID, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) UpdateByOwner(ctx context.Context, projectID, ownerID string, updates map[string]interface{}) error {
	args := m.Called(ctx, projectID, ownerID, updates)
	return args.Error(0)
}

func (m *MockProjectRepository) IsOwner(ctx context.Context, projectID, ownerID string) (bool, error) {
	args := m.Called(ctx, projectID, ownerID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) SoftDelete(ctx context.Context, projectID, ownerID string) error {
	args := m.Called(ctx, projectID, ownerID)
	return args.Error(0)
}

func (m *MockProjectRepository) HardDelete(ctx context.Context, projectID string) error {
	args := m.Called(ctx, projectID)
	return args.Error(0)
}

func (m *MockProjectRepository) Restore(ctx context.Context, projectID, ownerID string) error {
	args := m.Called(ctx, projectID, ownerID)
	return args.Error(0)
}

func (m *MockProjectRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	args := m.Called(ctx, ownerID)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) CreateWithTransaction(ctx context.Context, project *writer.Project, callback func(ctx context.Context) error) error {
	args := m.Called(ctx, project, callback)
	return args.Error(0)
}

// MockEventBus Mock事件总线
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, event serviceBase.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event serviceBase.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventType string, handler serviceBase.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	args := m.Called(eventType, handlerName)
	return args.Error(0)
}

// createTestProject 辅助函数：创建测试项目
func createTestProject(id, authorID, title string, status writer.ProjectStatus) *writer.Project {
	return &writer.Project{
		IdentifiedEntity: writerBase.IdentifiedEntity{ID: id},
		OwnedEntity:      writerBase.OwnedEntity{AuthorID: authorID},
		TitledEntity:     writerBase.TitledEntity{Title: title},
		Status:           status,
		WritingType:      "novel",
		Visibility:       writer.VisibilityPrivate,
	}
}

// TestProjectService_CreateProject_Success 测试创建项目成功
func TestProjectService_CreateProject_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	userID := primitive.NewObjectID().Hex()
	ctx = context.WithValue(ctx, "userId", userID)

	req := &CreateProjectRequest{
		Title:    "测试项目",
		Summary:  "这是一个测试项目",
		CoverURL: "https://example.com/cover.jpg",
		Category: "玄幻",
		Tags:     []string{"修仙", "热血"},
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*writer.Project")).Return(nil).Once()
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Once()

	// Act
	resp, err := projectService.CreateProject(ctx, req)

	// Assert
	assert.NoError(t, err, "创建项目应该成功")
	assert.NotNil(t, resp, "响应不应为空")
	assert.NotEmpty(t, resp.ProjectID, "项目ID不应为空")
	assert.Equal(t, req.Title, resp.Title, "项目标题应该匹配")
	assert.Equal(t, "draft", resp.Status, "状态应该是草稿")

	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)

	t.Log("✅ 创建项目成功测试通过")
}

// TestProjectService_CreateProject_EmptyTitle_ReturnError 测试空标题创建项目失败
func TestProjectService_CreateProject_EmptyTitle_ReturnError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	userID := primitive.NewObjectID().Hex()
	ctx = context.WithValue(ctx, "userId", userID)

	// Act
	tests := []struct {
		name    string
		title   string
		wantErr string
	}{
		{"empty title", "", "项目标题不能为空"},
		{"title too long", string(make([]byte, 101)), "项目标题不能超过100字符"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &CreateProjectRequest{
				Title: tt.title,
			}

			// Act
			resp, err := projectService.CreateProject(ctx, req)

			// Assert
			assert.Error(t, err, "应该返回错误")
			assert.Nil(t, resp, "响应应该为空")
			assert.Contains(t, err.Error(), tt.wantErr, "错误信息应该匹配")
		})
	}

	t.Log("✅ 空标题验证测试通过")
}

// TestProjectService_CreateProject_NoUserID_ReturnError 测试未登录创建项目失败
func TestProjectService_CreateProject_NoUserID_ReturnError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	req := &CreateProjectRequest{
		Title: "测试项目",
	}

	// Act
	resp, err := projectService.CreateProject(ctx, req)

	// Assert
	assert.Error(t, err, "应该返回错误")
	assert.Nil(t, resp, "响应应该为空")
	assert.Contains(t, err.Error(), "用户未登录", "错误信息应包含用户未登录")

	t.Log("✅ 未登录创建项目测试通过")
}

// TestProjectService_GetProject_Success 测试获取项目成功
func TestProjectService_GetProject_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	userID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()

	expectedProject := createTestProject(projectID, userID, "测试项目", writer.StatusDraft)

	ctx = context.WithValue(ctx, "userId", userID)
	mockRepo.On("GetByID", ctx, projectID).Return(expectedProject, nil).Once()

	// Act
	project, err := projectService.GetProject(ctx, projectID)

	// Assert
	assert.NoError(t, err, "获取项目应该成功")
	assert.NotNil(t, project, "项目不应为空")
	assert.Equal(t, expectedProject.ID, project.ID, "项目ID应该匹配")
	assert.Equal(t, expectedProject.Title, project.Title, "项目标题应该匹配")

	mockRepo.AssertExpectations(t)

	t.Log("✅ 获取项目成功测试通过")
}

// TestProjectService_GetProject_NotFound_ReturnError 测试获取不存在的项目
func TestProjectService_GetProject_NotFound_ReturnError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	userID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()

	ctx = context.WithValue(ctx, "userId", userID)
	mockRepo.On("GetByID", ctx, projectID).Return(nil, nil).Once()

	// Act
	project, err := projectService.GetProject(ctx, projectID)

	// Assert
	assert.Error(t, err, "应该返回错误")
	assert.Nil(t, project, "项目应该为空")
	assert.Contains(t, err.Error(), "项目不存在", "错误信息应包含项目不存在")

	mockRepo.AssertExpectations(t)

	t.Log("✅ 获取不存在项目测试通过")
}

// TestProjectService_UpdateProject_Success 测试更新项目成功
func TestProjectService_UpdateProject_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	userID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()

	existingProject := createTestProject(projectID, userID, "原标题", writer.StatusDraft)

	ctx = context.WithValue(ctx, "userId", userID)
	mockRepo.On("GetByID", ctx, projectID).Return(existingProject, nil).Once()
	mockRepo.On("Update", ctx, projectID, mock.Anything).Return(nil).Once()
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Once()

	// Act
	req := &UpdateProjectRequest{
		Title:   "新标题",
		Summary: "新简介",
	}
	err := projectService.UpdateProject(ctx, projectID, req)

	// Assert
	assert.NoError(t, err, "更新项目应该成功")

	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)

	t.Log("✅ 更新项目成功测试通过")
}

// TestProjectService_DeleteProject_Success 测试删除项目成功
func TestProjectService_DeleteProject_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	userID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()

	existingProject := createTestProject(projectID, userID, "测试项目", writer.StatusDraft)

	ctx = context.WithValue(ctx, "userId", userID)
	mockRepo.On("GetByID", ctx, projectID).Return(existingProject, nil).Once()
	mockRepo.On("SoftDelete", ctx, projectID, userID).Return(nil).Once()
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Once()

	// Act
	err := projectService.DeleteProject(ctx, projectID)

	// Assert
	assert.NoError(t, err, "删除项目应该成功")

	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)

	t.Log("✅ 删除项目成功测试通过")
}

// TestProjectService_ListMyProjects_Success 测试获取我的项目列表成功
func TestProjectService_ListMyProjects_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	userID := primitive.NewObjectID().Hex()

	expectedProjects := []*writer.Project{
		createTestProject(primitive.NewObjectID().Hex(), userID, "项目1", writer.StatusDraft),
		createTestProject(primitive.NewObjectID().Hex(), userID, "项目2", writer.StatusSerializing),
	}

	ctx = context.WithValue(ctx, "userId", userID)
	mockRepo.On("GetListByOwnerID", ctx, userID, int64(10), int64(0)).Return(expectedProjects, nil).Once()
	mockRepo.On("CountByOwner", ctx, userID).Return(int64(2), nil).Once()

	// Act
	req := &ListProjectsRequest{
		Page:     1,
		PageSize: 10,
	}
	resp, err := projectService.ListMyProjects(ctx, req)

	// Assert
	assert.NoError(t, err, "获取项目列表应该成功")
	assert.NotNil(t, resp, "响应不应为空")
	assert.Len(t, resp.Projects, 2, "应该有2个项目")
	assert.Equal(t, int64(2), resp.Total, "总数应该是2")
	assert.Equal(t, 1, resp.Page, "页码应该是1")
	assert.Equal(t, 10, resp.PageSize, "每页数量应该是10")

	mockRepo.AssertExpectations(t)

	t.Log("✅ 获取我的项目列表成功测试通过")
}

// TestProjectService_ListMyProjects_WithStatus_Success 测试按状态获取项目列表成功
func TestProjectService_ListMyProjects_WithStatus_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	userID := primitive.NewObjectID().Hex()

	expectedProjects := []*writer.Project{
		createTestProject(primitive.NewObjectID().Hex(), userID, "连载中项目", writer.StatusSerializing),
	}

	ctx = context.WithValue(ctx, "userId", userID)
	mockRepo.On("GetByOwnerAndStatus", ctx, userID, "serializing", int64(10), int64(0)).Return(expectedProjects, nil).Once()
	mockRepo.On("CountByOwner", ctx, userID).Return(int64(1), nil).Once()

	// Act
	req := &ListProjectsRequest{
		Page:     1,
		PageSize: 10,
		Status:   "serializing",
	}
	resp, err := projectService.ListMyProjects(ctx, req)

	// Assert
	assert.NoError(t, err, "获取项目列表应该成功")
	assert.NotNil(t, resp, "响应不应为空")
	assert.Len(t, resp.Projects, 1, "应该有1个项目")
	assert.Equal(t, writer.StatusSerializing, resp.Projects[0].Status, "状态应该是连载中")

	mockRepo.AssertExpectations(t)

	t.Log("✅ 按状态获取项目列表成功测试通过")
}

// TestProjectService_UpdateProjectStatistics_Success 测试更新项目统计成功
func TestProjectService_UpdateProjectStatistics_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()

	existingProject := createTestProject(projectID, primitive.NewObjectID().Hex(), "测试项目", writer.StatusDraft)

	stats := &writer.ProjectStats{
		TotalWords:    10000,
		ChapterCount:  20,
		DocumentCount: 15,
		LastUpdateAt:  time.Now(),
	}

	mockRepo.On("GetByID", ctx, projectID).Return(existingProject, nil).Once()
	mockRepo.On("Update", ctx, projectID, mock.Anything).Return(nil).Once()

	// Act
	err := projectService.UpdateProjectStatistics(ctx, projectID, stats)

	// Assert
	assert.NoError(t, err, "更新统计应该成功")

	mockRepo.AssertExpectations(t)

	t.Log("✅ 更新项目统计成功测试通过")
}

// TestProjectService_RestoreProjectByID_Success 测试恢复项目成功
func TestProjectService_RestoreProjectByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	userID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()

	existingProject := createTestProject(projectID, userID, "测试项目", writer.StatusDraft)

	mockRepo.On("GetByID", ctx, projectID).Return(existingProject, nil).Once()
	mockRepo.On("Restore", ctx, projectID, userID).Return(nil).Once()
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Once()

	// Act
	err := projectService.RestoreProjectByID(ctx, projectID, userID)

	// Assert
	assert.NoError(t, err, "恢复项目应该成功")

	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)

	t.Log("✅ 恢复项目成功测试通过")
}

// TestProjectService_DeleteHard_Success 测试物理删除项目成功
func TestProjectService_DeleteHard_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()

	mockRepo.On("HardDelete", ctx, projectID).Return(nil).Once()
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Once()

	// Act
	err := projectService.DeleteHard(ctx, projectID)

	// Assert
	assert.NoError(t, err, "物理删除项目应该成功")

	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)

	t.Log("✅ 物理删除项目成功测试通过")
}

// TestProjectService_GetByIDWithoutAuth_Success 测试无权限验证获取项目成功
func TestProjectService_GetByIDWithoutAuth_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	projectService := NewProjectService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	authorID := primitive.NewObjectID().Hex()

	expectedProject := createTestProject(projectID, authorID, "测试项目", writer.StatusDraft)

	mockRepo.On("GetByID", ctx, projectID).Return(expectedProject, nil).Once()

	// Act
	project, err := projectService.GetByIDWithoutAuth(ctx, projectID)

	// Assert
	assert.NoError(t, err, "获取项目应该成功")
	assert.NotNil(t, project, "项目不应为空")
	assert.Equal(t, expectedProject.ID, project.ID, "项目ID应该匹配")

	mockRepo.AssertExpectations(t)

	t.Log("✅ 无权限验证获取项目成功测试通过")
}
