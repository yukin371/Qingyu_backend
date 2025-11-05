package api

import (
	"Qingyu_backend/models/writer"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	writerAPI "Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/project"
)

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
				repo.On("Create", mock.Anything, mock.AnythingOfType("*writer.Project")).Return(nil)
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
				projects := []*writer.Project{
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
				projects := []*writer.Project{
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
				repo.On("GetListByOwnerID", mock.Anything, "user123", int64(10), int64(0)).Return([]*writer.Project{}, nil)
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
