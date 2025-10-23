package project_test

import (
	"Qingyu_backend/models/writer"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	base_infra "Qingyu_backend/repository/interfaces/infrastructure"
	writingRepo "Qingyu_backend/repository/interfaces/writing"
	"Qingyu_backend/service/base"
	projectService "Qingyu_backend/service/project"
)

// ============ Mock ============

// MockProjectRepository 简化的Mock
type MockProjectRepo struct {
	mock.Mock
	writingRepo.ProjectRepository
}

func (m *MockProjectRepo) Create(ctx context.Context, project *writer.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepo) GetByID(ctx context.Context, id string) (*writer.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Project), args.Error(1)
}

func (m *MockProjectRepo) GetListByOwnerID(ctx context.Context, ownerID string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, ownerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepo) GetByOwnerAndStatus(ctx context.Context, ownerID, status string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, ownerID, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepo) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepo) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockProjectRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepo) SoftDelete(ctx context.Context, projectID, ownerID string) error {
	args := m.Called(ctx, projectID, ownerID)
	return args.Error(0)
}

func (m *MockProjectRepo) Count(ctx context.Context, filter base_infra.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// MockEventBus 简化的Mock
type MockEventBus struct {
	mock.Mock
	base.EventBus
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// ============ Test Helpers ============

func createContext(userID string) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, "userID", userID)
}

func createProject(id, title, authorID string) *writer.Project {
	return &writer.Project{
		ID:         id,
		Title:      title,
		AuthorID:   authorID,
		Status:     writer.StatusDraft,
		Visibility: writer.VisibilityPrivate,
	}
}

// ============ CreateProject Tests ============

func TestCreateProject_Success(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	mockEventBus := new(MockEventBus)
	service := projectService.NewProjectService(mockRepo, mockEventBus)

	ctx := createContext("user123")
	req := &projectService.CreateProjectRequest{
		Title:   "测试项目",
		Summary: "简介",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*document.Project")).Return(nil)
	mockEventBus.On("PublishAsync", ctx, mock.AnythingOfType("*base.BaseEvent")).Return(nil)

	// Act
	result, err := service.CreateProject(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "测试项目", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestCreateProject_EmptyTitle(t *testing.T) {
	service := projectService.NewProjectService(nil, nil)
	ctx := createContext("user123")
	req := &projectService.CreateProjectRequest{Title: ""}

	result, err := service.CreateProject(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateProject_NoUser(t *testing.T) {
	service := projectService.NewProjectService(nil, nil)
	ctx := context.Background()
	req := &projectService.CreateProjectRequest{Title: "测试"}

	result, err := service.CreateProject(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============ GetProject Tests ============

func TestGetProject_Success(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	service := projectService.NewProjectService(mockRepo, nil)

	userID := "user123"
	ctx := createContext(userID)
	project := createProject("proj1", "测试项目", userID)

	mockRepo.On("GetByID", ctx, "proj1").Return(project, nil)

	result, err := service.GetProject(ctx, "proj1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "测试项目", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestGetProject_NotFound(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	service := projectService.NewProjectService(mockRepo, nil)

	ctx := createContext("user123")
	mockRepo.On("GetByID", ctx, "notfound").Return(nil, nil)

	result, err := service.GetProject(ctx, "notfound")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetProject_NoPermission(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	service := projectService.NewProjectService(mockRepo, nil)

	ctx := createContext("user123")
	project := createProject("proj1", "别人的项目", "user456")

	mockRepo.On("GetByID", ctx, "proj1").Return(project, nil)

	result, err := service.GetProject(ctx, "proj1")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// ============ ListMyProjects Tests ============

func TestListMyProjects_Success(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	service := projectService.NewProjectService(mockRepo, nil)

	userID := "user123"
	ctx := createContext(userID)

	projects := []*writer.Project{
		createProject("proj1", "项目1", userID),
		createProject("proj2", "项目2", userID),
	}

	req := &projectService.ListProjectsRequest{Page: 1, PageSize: 10}
	mockRepo.On("GetListByOwnerID", ctx, userID, int64(10), int64(0)).Return(projects, nil)
	mockRepo.On("CountByOwner", ctx, userID).Return(int64(2), nil)

	result, err := service.ListMyProjects(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, result.Projects, 2)
	assert.Equal(t, int64(2), result.Total)
	mockRepo.AssertExpectations(t)
}

func TestListMyProjects_WithStatus(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	service := projectService.NewProjectService(mockRepo, nil)

	userID := "user123"
	ctx := createContext(userID)

	projects := []*writer.Project{
		createProject("proj1", "草稿项目", userID),
	}

	req := &projectService.ListProjectsRequest{Page: 1, PageSize: 10, Status: "draft"}
	mockRepo.On("GetByOwnerAndStatus", ctx, userID, "draft", int64(10), int64(0)).Return(projects, nil)
	mockRepo.On("CountByOwner", ctx, userID).Return(int64(1), nil)

	result, err := service.ListMyProjects(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, result.Projects, 1)
	mockRepo.AssertExpectations(t)
}

// ============ UpdateProject Tests ============

func TestUpdateProject_Success(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	mockEventBus := new(MockEventBus)
	service := projectService.NewProjectService(mockRepo, mockEventBus)

	userID := "user123"
	ctx := createContext(userID)
	project := createProject("proj1", "旧标题", userID)

	mockRepo.On("GetByID", ctx, "proj1").Return(project, nil)
	mockRepo.On("Update", ctx, "proj1", mock.AnythingOfType("map[string]interface {}")).Return(nil)
	mockEventBus.On("PublishAsync", ctx, mock.AnythingOfType("*base.BaseEvent")).Return(nil)

	req := &projectService.UpdateProjectRequest{Title: "新标题"}
	err := service.UpdateProject(ctx, "proj1", req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProject_NoPermission(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	service := projectService.NewProjectService(mockRepo, nil)

	ctx := createContext("user123")
	project := createProject("proj1", "别人的项目", "user456")

	mockRepo.On("GetByID", ctx, "proj1").Return(project, nil)

	req := &projectService.UpdateProjectRequest{Title: "新标题"}
	err := service.UpdateProject(ctx, "proj1", req)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// ============ DeleteProject Tests ============

func TestDeleteProject_Success(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	mockEventBus := new(MockEventBus)
	service := projectService.NewProjectService(mockRepo, mockEventBus)

	userID := "user123"
	ctx := createContext(userID)
	project := createProject("proj1", "要删除的项目", userID)

	mockRepo.On("GetByID", ctx, "proj1").Return(project, nil)
	mockRepo.On("SoftDelete", ctx, "proj1", userID).Return(nil)
	mockEventBus.On("PublishAsync", ctx, mock.AnythingOfType("*base.BaseEvent")).Return(nil)

	err := service.DeleteProject(ctx, "proj1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProject_NoPermission(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	service := projectService.NewProjectService(mockRepo, nil)

	ctx := createContext("user123")
	project := createProject("proj1", "别人的项目", "user456")

	mockRepo.On("GetByID", ctx, "proj1").Return(project, nil)

	err := service.DeleteProject(ctx, "proj1")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
