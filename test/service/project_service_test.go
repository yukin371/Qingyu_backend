package service_test

import (
	"Qingyu_backend/models/writer"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	pkgErrors "Qingyu_backend/pkg/errors"
	"Qingyu_backend/repository/interfaces/infrastructure"
	"Qingyu_backend/service/interfaces/base"
	"Qingyu_backend/service/project"
)

// MockProjectRepository 模拟ProjectRepository
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(ctx context.Context, project *writer.Project) error {
	args := m.Called(ctx, project)
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

func (m *MockProjectRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.Project, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) GetListByOwnerID(ctx context.Context, ownerID string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, ownerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByOwnerAndStatus(ctx context.Context, ownerID, status string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, ownerID, status, limit, offset)
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
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) CreateWithTransaction(ctx context.Context, project *writer.Project, callback func(ctx context.Context) error) error {
	args := m.Called(ctx, project, callback)
	return args.Error(0)
}

func (m *MockProjectRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockProjectRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// MockEventBus 模拟EventBus
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	args := m.Called(eventType, handlerName)
	return args.Error(0)
}

// TestProjectService_CreateProject 测试创建项目
func TestProjectService_CreateProject(t *testing.T) {
	// 1. 创建Mock Repository
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	// 2. 创建Service
	service := project.NewProjectService(mockProjectRepo, mockEventBus)

	ctx := context.WithValue(context.Background(), "userID", "user123")

	// 3. 设置Mock期望
	mockProjectRepo.On("Create", mock.Anything, mock.AnythingOfType("*writer.Project")).Return(nil).Run(func(args mock.Arguments) {
		proj := args.Get(1).(*writer.Project)
		proj.ID = "project123"
		proj.CreatedAt = time.Now()
	})
	mockEventBus.On("PublishAsync", mock.Anything, mock.Anything).Return(nil)

	// 4. 执行测试
	t.Run("正常创建项目", func(t *testing.T) {
		req := &project.CreateProjectRequest{
			Title:    "测试项目",
			Summary:  "测试简介",
			Category: "玄幻",
			Tags:     []string{"热血"},
		}

		resp, err := service.CreateProject(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "测试项目", resp.Title)
		assert.Equal(t, "draft", resp.Status)
		assert.NotEmpty(t, resp.ProjectID)
	})

	// 5. 测试参数验证
	t.Run("标题为空应该失败", func(t *testing.T) {
		req := &project.CreateProjectRequest{
			Title: "",
		}

		resp, err := service.CreateProject(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)

		// 验证错误类型
		serviceErr, ok := err.(*pkgErrors.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkgErrors.ServiceErrorValidation, serviceErr.Type)
	})

	// 6. 测试未登录用户
	t.Run("未登录用户应该失败", func(t *testing.T) {
		ctxNoUser := context.Background()

		req := &project.CreateProjectRequest{
			Title: "测试项目",
		}

		resp, err := service.CreateProject(ctxNoUser, req)

		assert.Error(t, err)
		assert.Nil(t, resp)

		serviceErr, ok := err.(*pkgErrors.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkgErrors.ServiceErrorUnauthorized, serviceErr.Type)
	})

	// 7. 验证Mock调用
	mockProjectRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestProjectService_GetProject 测试获取项目
func TestProjectService_GetProject(t *testing.T) {
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := project.NewProjectService(mockProjectRepo, mockEventBus)

	ctx := context.WithValue(context.Background(), "userID", "user123")

	// 测试获取存在的项目
	t.Run("获取存在的项目", func(t *testing.T) {
		testProject := &writer.Project{
			ID:         "project123",
			AuthorID:   "user123",
			Title:      "测试项目",
			Visibility: writer.VisibilityPrivate,
		}

		mockProjectRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil).Once()

		result, err := service.GetProject(ctx, "project123")

		require.NoError(t, err)
		assert.Equal(t, "测试项目", result.Title)
	})

	// 测试项目不存在
	t.Run("项目不存在", func(t *testing.T) {
		mockProjectRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil).Once()

		result, err := service.GetProject(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Nil(t, result)

		serviceErr, ok := err.(*pkgErrors.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkgErrors.ServiceErrorNotFound, serviceErr.Type)
	})

	// 测试无权限访问私密项目
	t.Run("无权限访问私密项目", func(t *testing.T) {
		testProject := &writer.Project{
			ID:         "project123",
			AuthorID:   "other_user",
			Title:      "他人项目",
			Visibility: writer.VisibilityPrivate,
		}

		mockProjectRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil).Once()

		result, err := service.GetProject(ctx, "project123")

		assert.Error(t, err)
		assert.Nil(t, result)

		serviceErr, ok := err.(*pkgErrors.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkgErrors.ServiceErrorForbidden, serviceErr.Type)
	})

	// 测试访问公开项目
	t.Run("访问公开项目", func(t *testing.T) {
		testProject := &writer.Project{
			ID:         "project123",
			AuthorID:   "other_user",
			Title:      "公开项目",
			Visibility: writer.VisibilityPublic,
		}

		mockProjectRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil).Once()

		result, err := service.GetProject(ctx, "project123")

		require.NoError(t, err)
		assert.Equal(t, "公开项目", result.Title)
	})

	mockProjectRepo.AssertExpectations(t)
}

// TestProjectService_ListMyProjects 测试获取项目列表
func TestProjectService_ListMyProjects(t *testing.T) {
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := project.NewProjectService(mockProjectRepo, mockEventBus)

	ctx := context.WithValue(context.Background(), "userID", "user123")

	t.Run("获取所有项目", func(t *testing.T) {
		testProjects := []*writer.Project{
			{ID: "proj1", Title: "项目1"},
			{ID: "proj2", Title: "项目2"},
		}

		mockProjectRepo.On("GetListByOwnerID", mock.Anything, "user123", int64(10), int64(0)).Return(testProjects, nil).Once()
		mockProjectRepo.On("CountByOwner", mock.Anything, "user123").Return(int64(2), nil).Once()

		req := &project.ListProjectsRequest{
			Page:     1,
			PageSize: 10,
		}

		resp, err := service.ListMyProjects(ctx, req)

		require.NoError(t, err)
		assert.Equal(t, 2, len(resp.Projects))
		assert.Equal(t, int64(2), resp.Total)
		assert.Equal(t, 1, resp.Page)
	})

	t.Run("按状态筛选", func(t *testing.T) {
		testProjects := []*writer.Project{
			{ID: "proj1", Title: "项目1", Status: writer.StatusDraft},
		}

		mockProjectRepo.On("GetByOwnerAndStatus", mock.Anything, "user123", "draft", int64(10), int64(0)).Return(testProjects, nil).Once()
		mockProjectRepo.On("CountByOwner", mock.Anything, "user123").Return(int64(1), nil).Once()

		req := &project.ListProjectsRequest{
			Page:     1,
			PageSize: 10,
			Status:   "draft",
		}

		resp, err := service.ListMyProjects(ctx, req)

		require.NoError(t, err)
		assert.Equal(t, 1, len(resp.Projects))
	})

	mockProjectRepo.AssertExpectations(t)
}

// TestProjectService_UpdateProject 测试更新项目
func TestProjectService_UpdateProject(t *testing.T) {
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := project.NewProjectService(mockProjectRepo, mockEventBus)

	ctx := context.WithValue(context.Background(), "userID", "user123")

	t.Run("所有者更新成功", func(t *testing.T) {
		testProject := &writer.Project{
			ID:       "project123",
			AuthorID: "user123",
			Title:    "原标题",
		}

		mockProjectRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil).Once()
		mockProjectRepo.On("Update", mock.Anything, "project123", mock.Anything).Return(nil).Once()
		mockEventBus.On("PublishAsync", mock.Anything, mock.Anything).Return(nil)

		req := &project.UpdateProjectRequest{
			Title: "新标题",
		}

		err := service.UpdateProject(ctx, "project123", req)

		assert.NoError(t, err)
	})

	t.Run("非所有者无法更新", func(t *testing.T) {
		testProject := &writer.Project{
			ID:       "project123",
			AuthorID: "other_user",
			Title:    "原标题",
		}

		mockProjectRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil).Once()

		req := &project.UpdateProjectRequest{
			Title: "新标题",
		}

		err := service.UpdateProject(ctx, "project123", req)

		assert.Error(t, err)

		serviceErr, ok := err.(*pkgErrors.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkgErrors.ServiceErrorForbidden, serviceErr.Type)
	})

	t.Run("编辑者可以更新", func(t *testing.T) {
		now := time.Now()
		testProject := &writer.Project{
			ID:       "project123",
			AuthorID: "other_user",
			Title:    "原标题",
			Collaborators: []writer.Collaborator{
				{
					UserID:     "user123",
					Role:       writer.RoleEditor,
					InvitedAt:  now,
					AcceptedAt: &now,
				},
			},
		}

		mockProjectRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil).Once()
		mockProjectRepo.On("Update", mock.Anything, "project123", mock.Anything).Return(nil).Once()
		mockEventBus.On("PublishAsync", mock.Anything, mock.Anything).Return(nil)

		req := &project.UpdateProjectRequest{
			Title: "新标题",
		}

		err := service.UpdateProject(ctx, "project123", req)

		assert.NoError(t, err)
	})

	mockProjectRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestProjectService_DeleteProject 测试删除项目
func TestProjectService_DeleteProject(t *testing.T) {
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := project.NewProjectService(mockProjectRepo, mockEventBus)

	ctx := context.WithValue(context.Background(), "userID", "user123")

	t.Run("所有者删除成功", func(t *testing.T) {
		testProject := &writer.Project{
			ID:       "project123",
			AuthorID: "user123",
			Title:    "测试项目",
		}

		mockProjectRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil).Once()
		mockProjectRepo.On("SoftDelete", mock.Anything, "project123", "user123").Return(nil).Once()
		mockEventBus.On("PublishAsync", mock.Anything, mock.Anything).Return(nil)

		err := service.DeleteProject(ctx, "project123")

		assert.NoError(t, err)
	})

	t.Run("非所有者无法删除", func(t *testing.T) {
		testProject := &writer.Project{
			ID:       "project123",
			AuthorID: "other_user",
			Title:    "测试项目",
		}

		mockProjectRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil).Once()

		err := service.DeleteProject(ctx, "project123")

		assert.Error(t, err)

		serviceErr, ok := err.(*pkgErrors.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkgErrors.ServiceErrorForbidden, serviceErr.Type)
		assert.Contains(t, serviceErr.Message, "只有项目所有者")
	})

	t.Run("项目不存在", func(t *testing.T) {
		mockProjectRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil).Once()

		err := service.DeleteProject(ctx, "nonexistent")

		assert.Error(t, err)

		serviceErr, ok := err.(*pkgErrors.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkgErrors.ServiceErrorNotFound, serviceErr.Type)
	})

	mockProjectRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestProjectService_UpdateProjectStatistics 测试更新统计
func TestProjectService_UpdateProjectStatistics(t *testing.T) {
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := project.NewProjectService(mockProjectRepo, mockEventBus)

	ctx := context.Background()

	t.Run("更新统计成功", func(t *testing.T) {
		testProject := &writer.Project{
			ID:       "project123",
			AuthorID: "user123",
			Title:    "测试项目",
		}

		stats := &writer.ProjectStats{
			TotalWords:    10000,
			ChapterCount:  10,
			DocumentCount: 20,
		}

		mockProjectRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil).Once()
		mockProjectRepo.On("Update", mock.Anything, "project123", mock.Anything).Return(nil).Once()

		err := service.UpdateProjectStatistics(ctx, "project123", stats)

		assert.NoError(t, err)
	})

	mockProjectRepo.AssertExpectations(t)
}

// TestProjectService_ServiceInterface 测试BaseService接口实现
func TestProjectService_ServiceInterface(t *testing.T) {
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := project.NewProjectService(mockProjectRepo, mockEventBus)

	ctx := context.Background()

	t.Run("Initialize方法", func(t *testing.T) {
		err := service.Initialize(ctx)
		assert.NoError(t, err)
	})

	t.Run("Health方法", func(t *testing.T) {
		mockProjectRepo.On("Health", mock.Anything).Return(nil).Once()

		err := service.Health(ctx)
		assert.NoError(t, err)
	})

	t.Run("Close方法", func(t *testing.T) {
		err := service.Close(ctx)
		assert.NoError(t, err)
	})

	t.Run("GetServiceName方法", func(t *testing.T) {
		name := service.GetServiceName()
		assert.Equal(t, "ProjectService", name)
	})

	t.Run("GetVersion方法", func(t *testing.T) {
		version := service.GetVersion()
		assert.Equal(t, "1.0.0", version)
	})

	mockProjectRepo.AssertExpectations(t)
}
