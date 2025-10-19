package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	writerAPI "Qingyu_backend/api/v1/writer"
	"Qingyu_backend/models/document"
	"Qingyu_backend/repository/interfaces/infrastructure"
	"Qingyu_backend/service/base"
	"Qingyu_backend/service/project"
)

// === Mock Repositories ===

type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(ctx context.Context, project *document.Project) error {
	args := m.Called(ctx, project)
	// 模拟设置ID和时间戳（只要不返回错误就设置）
	if args.Error(0) == nil {
		if project.ID == "" {
			project.ID = "mock_project_id"
		}
		project.CreatedAt = time.Now()
		project.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id string) (*document.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepository) GetListByOwnerID(ctx context.Context, ownerID string, limit, offset int64) ([]*document.Project, error) {
	args := m.Called(ctx, ownerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*document.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByOwnerAndStatus(ctx context.Context, ownerID string, status string, limit, offset int64) ([]*document.Project, error) {
	args := m.Called(ctx, ownerID, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*document.Project), args.Error(1)
}

func (m *MockProjectRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) UpdateStatistics(ctx context.Context, id string, wordCount, chapterCount int64) error {
	args := m.Called(ctx, id, wordCount, chapterCount)
	return args.Error(0)
}

func (m *MockProjectRepository) GetByCategory(ctx context.Context, category string, limit, offset int64) ([]*document.Project, error) {
	args := m.Called(ctx, category, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*document.Project), args.Error(1)
}

func (m *MockProjectRepository) GetPublicProjects(ctx context.Context, limit, offset int64) ([]*document.Project, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*document.Project), args.Error(1)
}

func (m *MockProjectRepository) SearchProjects(ctx context.Context, keyword string, limit, offset int64) ([]*document.Project, error) {
	args := m.Called(ctx, keyword, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*document.Project), args.Error(1)
}

func (m *MockProjectRepository) UpdateStatus(ctx context.Context, id string, status document.ProjectStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockProjectRepository) UpdateVisibility(ctx context.Context, id string, visibility document.Visibility) error {
	args := m.Called(ctx, id, visibility)
	return args.Error(0)
}

func (m *MockProjectRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
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

func (m *MockProjectRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) CreateWithTransaction(ctx context.Context, project *document.Project, callback func(ctx context.Context) error) error {
	args := m.Called(ctx, project, callback)
	return args.Error(0)
}

// 实现其他接口方法
func (m *MockProjectRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*document.Project, error) {
	return nil, nil
}

func (m *MockProjectRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	return 0, nil
}

func (m *MockProjectRepository) Exists(ctx context.Context, id string) (bool, error) {
	return false, nil
}

// === Mock EventBus ===

type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	return nil
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	return nil
}

func (m *MockEventBus) Publish(ctx context.Context, event base.Event) error {
	return nil
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	return nil
}

// === 测试辅助函数 ===

func setupProjectTestRouter(projectService *project.ProjectService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := writerAPI.NewProjectApi(projectService)

	// 注册路由
	v1 := r.Group("/api/v1")
	{
		projects := v1.Group("/projects")
		{
			projects.POST("", api.CreateProject)
			projects.GET("", api.ListProjects)
			projects.GET("/:id", api.GetProject)
			projects.PUT("/:id", api.UpdateProject)
			projects.DELETE("/:id", api.DeleteProject)
			projects.PUT("/:id/statistics", api.UpdateProjectStatistics)
		}
	}

	return r
}

func createTestProject(authorID string) *document.Project {
	now := time.Now()
	return &document.Project{
		ID:         "test_project_id",
		AuthorID:   authorID,
		Title:      "测试项目",
		Summary:    "这是一个测试项目",
		CoverURL:   "https://example.com/cover.jpg",
		Category:   "科幻",
		Tags:       []string{"科幻", "冒险"},
		Status:     document.StatusDraft,
		Visibility: document.VisibilityPrivate,
		Statistics: document.ProjectStats{
			TotalWords:    0,
			ChapterCount:  0,
			DocumentCount: 0,
			LastUpdateAt:  now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// === 测试用例 ===

func TestProjectApi_CreateProject(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    project.CreateProjectRequest
		setupMock      func(*MockProjectRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "成功创建项目",
			userID: "user123",
			requestBody: project.CreateProjectRequest{
				Title:    "新项目",
				Summary:  "项目描述",
				CoverURL: "https://example.com/cover.jpg",
				Category: "科幻",
				Tags:     []string{"科幻", "冒险"},
			},
			setupMock: func(repo *MockProjectRepository) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*document.Project")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(201), resp["code"])
				assert.Equal(t, "创建成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.NotEmpty(t, data["projectId"])
				assert.Equal(t, "新项目", data["title"])
			},
		},
		{
			name:   "缺少必填字段",
			userID: "user123",
			requestBody: project.CreateProjectRequest{
				Summary: "项目描述",
			},
			setupMock:      func(repo *MockProjectRepository) {},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
				assert.Equal(t, "创建失败", resp["message"])
				// Service层验证失败，返回的是Internal Server Error
			},
		},
		{
			name:   "未登录用户",
			userID: "",
			requestBody: project.CreateProjectRequest{
				Title:   "新项目",
				Summary: "项目描述",
			},
			setupMock:      func(repo *MockProjectRepository) {},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock
			mockRepo := new(MockProjectRepository)
			mockEventBus := new(MockEventBus)
			tt.setupMock(mockRepo)

			// 创建Service和Router
			projectService := project.NewProjectService(mockRepo, mockEventBus)
			router := setupProjectTestRouter(projectService)

			// 构造请求
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// 添加上下文
			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "userID", tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectApi_GetProject(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		userID         string
		setupMock      func(*MockProjectRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "成功获取项目",
			projectID: "project123",
			userID:    "user123",
			setupMock: func(repo *MockProjectRepository) {
				testProject := createTestProject("user123")
				testProject.ID = "project123"
				repo.On("GetByID", mock.Anything, "project123").Return(testProject, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "获取成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "project123", data["id"])
				assert.Equal(t, "测试项目", data["title"])
			},
		},
		{
			name:      "项目不存在",
			projectID: "nonexistent",
			userID:    "user123",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
			},
		},
		{
			name:      "无权限访问",
			projectID: "project123",
			userID:    "other_user",
			setupMock: func(repo *MockProjectRepository) {
				testProject := createTestProject("user123")
				testProject.ID = "project123"
				repo.On("GetByID", mock.Anything, "project123").Return(testProject, nil)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock
			mockRepo := new(MockProjectRepository)
			mockEventBus := new(MockEventBus)
			tt.setupMock(mockRepo)

			// 创建Service和Router
			projectService := project.NewProjectService(mockRepo, mockEventBus)
			router := setupProjectTestRouter(projectService)

			// 构造请求
			req := httptest.NewRequest("GET", "/api/v1/projects/"+tt.projectID, nil)

			// 添加上下文
			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "userID", tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectApi_ListProjects(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		queryParams    string
		setupMock      func(*MockProjectRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:        "成功获取项目列表",
			userID:      "user123",
			queryParams: "?page=1&pageSize=10",
			setupMock: func(repo *MockProjectRepository) {
				projects := []*document.Project{
					createTestProject("user123"),
					createTestProject("user123"),
				}
				repo.On("GetListByOwnerID", mock.Anything, "user123", int64(10), int64(0)).Return(projects, nil)
				repo.On("CountByOwner", mock.Anything, "user123").Return(int64(2), nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "获取成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				projects := data["projects"].([]interface{})
				assert.Len(t, projects, 2)
				assert.Equal(t, float64(2), data["total"])
			},
		},
		{
			name:        "按状态筛选",
			userID:      "user123",
			queryParams: "?page=1&pageSize=10&status=draft",
			setupMock: func(repo *MockProjectRepository) {
				projects := []*document.Project{
					createTestProject("user123"),
				}
				repo.On("GetByOwnerAndStatus", mock.Anything, "user123", "draft", int64(10), int64(0)).Return(projects, nil)
				repo.On("CountByOwner", mock.Anything, "user123").Return(int64(1), nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				data := resp["data"].(map[string]interface{})
				projects := data["projects"].([]interface{})
				assert.Len(t, projects, 1)
			},
		},
		{
			name:        "空列表",
			userID:      "user123",
			queryParams: "?page=1&pageSize=10",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetListByOwnerID", mock.Anything, "user123", int64(10), int64(0)).Return([]*document.Project{}, nil)
				repo.On("CountByOwner", mock.Anything, "user123").Return(int64(0), nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				data := resp["data"].(map[string]interface{})
				projects := data["projects"].([]interface{})
				assert.Len(t, projects, 0)
				assert.Equal(t, float64(0), data["total"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock
			mockRepo := new(MockProjectRepository)
			mockEventBus := new(MockEventBus)
			tt.setupMock(mockRepo)

			// 创建Service和Router
			projectService := project.NewProjectService(mockRepo, mockEventBus)
			router := setupProjectTestRouter(projectService)

			// 构造请求
			req := httptest.NewRequest("GET", "/api/v1/projects"+tt.queryParams, nil)

			// 添加上下文
			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "userID", tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectApi_UpdateProject(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		userID         string
		requestBody    project.UpdateProjectRequest
		setupMock      func(*MockProjectRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "成功更新项目",
			projectID: "project123",
			userID:    "user123",
			requestBody: project.UpdateProjectRequest{
				Title:   "更新后的标题",
				Summary: "更新后的描述",
			},
			setupMock: func(repo *MockProjectRepository) {
				testProject := createTestProject("user123")
				testProject.ID = "project123"
				repo.On("GetByID", mock.Anything, "project123").Return(testProject, nil)
				repo.On("Update", mock.Anything, "project123", mock.AnythingOfType("map[string]interface {}")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "更新成功", resp["message"])
			},
		},
		{
			name:      "无权限更新",
			projectID: "project123",
			userID:    "other_user",
			requestBody: project.UpdateProjectRequest{
				Title: "更新后的标题",
			},
			setupMock: func(repo *MockProjectRepository) {
				testProject := createTestProject("user123")
				testProject.ID = "project123"
				repo.On("GetByID", mock.Anything, "project123").Return(testProject, nil)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
			},
		},
		{
			name:      "项目不存在",
			projectID: "nonexistent",
			userID:    "user123",
			requestBody: project.UpdateProjectRequest{
				Title: "更新后的标题",
			},
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock
			mockRepo := new(MockProjectRepository)
			mockEventBus := new(MockEventBus)
			tt.setupMock(mockRepo)

			// 创建Service和Router
			projectService := project.NewProjectService(mockRepo, mockEventBus)
			router := setupProjectTestRouter(projectService)

			// 构造请求
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/api/v1/projects/"+tt.projectID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// 添加上下文
			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "userID", tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectApi_DeleteProject(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		userID         string
		setupMock      func(*MockProjectRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "成功删除项目",
			projectID: "project123",
			userID:    "user123",
			setupMock: func(repo *MockProjectRepository) {
				testProject := createTestProject("user123")
				testProject.ID = "project123"
				repo.On("GetByID", mock.Anything, "project123").Return(testProject, nil)
				repo.On("SoftDelete", mock.Anything, "project123", "user123").Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "删除成功", resp["message"])
			},
		},
		{
			name:      "无权限删除",
			projectID: "project123",
			userID:    "other_user",
			setupMock: func(repo *MockProjectRepository) {
				testProject := createTestProject("user123")
				testProject.ID = "project123"
				repo.On("GetByID", mock.Anything, "project123").Return(testProject, nil)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
			},
		},
		{
			name:      "项目不存在",
			projectID: "nonexistent",
			userID:    "user123",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock
			mockRepo := new(MockProjectRepository)
			mockEventBus := new(MockEventBus)
			tt.setupMock(mockRepo)

			// 创建Service和Router
			projectService := project.NewProjectService(mockRepo, mockEventBus)
			router := setupProjectTestRouter(projectService)

			// 构造请求
			req := httptest.NewRequest("DELETE", "/api/v1/projects/"+tt.projectID, nil)

			// 添加上下文
			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "userID", tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectApi_UpdateProjectStatistics(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		userID         string
		setupMock      func(*MockProjectRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "成功更新统计信息",
			projectID: "project123",
			userID:    "user123",
			setupMock: func(repo *MockProjectRepository) {
				testProject := createTestProject("user123")
				testProject.ID = "project123"
				repo.On("GetByID", mock.Anything, "project123").Return(testProject, nil)
				// UpdateProjectStatistics调用Update方法而不是UpdateStatistics
				repo.On("Update", mock.Anything, "project123", mock.AnythingOfType("map[string]interface {}")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "更新成功", resp["message"])
			},
		},
		{
			name:      "项目不存在",
			projectID: "nonexistent",
			userID:    "user123",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock
			mockRepo := new(MockProjectRepository)
			mockEventBus := new(MockEventBus)
			tt.setupMock(mockRepo)

			// 创建Service和Router
			projectService := project.NewProjectService(mockRepo, mockEventBus)
			router := setupProjectTestRouter(projectService)

			// 构造请求
			req := httptest.NewRequest("PUT", "/api/v1/projects/"+tt.projectID+"/statistics", nil)

			// 添加上下文
			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "userID", tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)

			mockRepo.AssertExpectations(t)
		})
	}
}
