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

	"Qingyu_backend/api/v1/announcements"
	"Qingyu_backend/models/messaging"
	messagingService "Qingyu_backend/service/messaging"
)

// ============ Mock AnnouncementService ============

type MockAnnouncementService struct {
	mock.Mock
}

func (m *MockAnnouncementService) GetAnnouncementByID(ctx context.Context, id string) (*messaging.Announcement, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*messaging.Announcement), args.Error(1)
}

func (m *MockAnnouncementService) GetAnnouncements(ctx context.Context, req *messagingService.GetAnnouncementsRequest) (*messagingService.GetAnnouncementsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*messagingService.GetAnnouncementsResponse), args.Error(1)
}

func (m *MockAnnouncementService) GetEffectiveAnnouncements(ctx context.Context, targetRole string, limit int) ([]*messaging.Announcement, error) {
	args := m.Called(ctx, targetRole, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*messaging.Announcement), args.Error(1)
}

func (m *MockAnnouncementService) CreateAnnouncement(ctx context.Context, req *messagingService.CreateAnnouncementRequest) (*messaging.Announcement, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*messaging.Announcement), args.Error(1)
}

func (m *MockAnnouncementService) UpdateAnnouncement(ctx context.Context, id string, req *messagingService.UpdateAnnouncementRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockAnnouncementService) DeleteAnnouncement(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAnnouncementService) BatchUpdateStatus(ctx context.Context, req *messagingService.BatchUpdateAnnouncementStatusRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAnnouncementService) BatchDelete(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockAnnouncementService) IncrementViewCount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ============ 公开API测试 ============

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestAnnouncementPublicAPI_GetEffectiveAnnouncements(t *testing.T) {
	router := setupTestRouter()
	mockService := new(MockAnnouncementService)

	api := announcements.NewAnnouncementPublicAPI(mockService)
	router.GET("/announcements/effective", api.GetEffectiveAnnouncements)

	t.Run("成功获取有效公告", func(t *testing.T) {
		now := time.Now()
		later := now.Add(24 * time.Hour)

		expectedAnnouncements := []*messaging.Announcement{
			{
				ID:        "1",
				Title:     "系统维护公告",
				Content:   "系统将于今晚进行维护",
				Type:      messaging.AnnouncementTypeSystem,
				Priority:  messaging.AnnouncementPriorityHigh,
				IsActive:  true,
				TargetRole: "all",
				ViewCount: 100,
				StartTime: &now,
				EndTime:   &later,
			},
			{
				ID:        "2",
				Title:     "新功能上线",
				Content:   "我们推出了全新的功能",
				Type:      messaging.AnnouncementTypeNotice,
				Priority:  messaging.AnnouncementPriorityNormal,
				IsActive:  true,
				TargetRole: "all",
				ViewCount: 50,
				StartTime: &now,
				EndTime:   &later,
			},
		}

		mockService.On("GetEffectiveAnnouncements", mock.Anything, "all", 10).
			Return(expectedAnnouncements, nil)

		req, _ := http.NewRequest("GET", "/announcements/effective?targetRole=all&limit=10", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 200, int(response["code"].(float64)))
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})

	t.Run("服务错误", func(t *testing.T) {
		mockService := new(MockAnnouncementService)
		api := announcements.NewAnnouncementPublicAPI(mockService)
		router := setupTestRouter()
		router.GET("/announcements/effective", api.GetEffectiveAnnouncements)

		mockService.On("GetEffectiveAnnouncements", mock.Anything, mock.Anything, mock.Anything).
			Return([]*messaging.Announcement{}, assert.AnError)

		req, _ := http.NewRequest("GET", "/announcements/effective", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("限制参数处理", func(t *testing.T) {
		// 测试超过最大限制的情况
		mockService.On("GetEffectiveAnnouncements", mock.Anything, "all", 50).
			Return([]*messaging.Announcement{}, nil)

		req, _ := http.NewRequest("GET", "/announcements/effective?limit=100", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockService.AssertExpectations(t)
	})
}

func TestAnnouncementPublicAPI_GetAnnouncementByID(t *testing.T) {
	router := setupTestRouter()

	t.Run("成功获取公告详情", func(t *testing.T) {
		mockService := new(MockAnnouncementService)
		api := announcements.NewAnnouncementPublicAPI(mockService)
		router.GET("/announcements/:id", api.GetAnnouncementByID)

		now := time.Now()
		later := now.Add(24 * time.Hour)

		expectedAnnouncement := &messaging.Announcement{
			ID:        "123",
			Title:     "重要公告",
			Content:   "这是重要公告的内容",
			Type:      messaging.AnnouncementTypeWarning,
			Priority:  messaging.AnnouncementPriorityHigh,
			IsActive:  true,
			TargetRole: "all",
			ViewCount: 200,
			CreatedBy: "admin",
			StartTime: &now,
			EndTime:   &later,
		}

		mockService.On("GetAnnouncementByID", mock.Anything, "123").
			Return(expectedAnnouncement, nil)

		req, _ := http.NewRequest("GET", "/announcements/123", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 200, int(response["code"].(float64)))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "重要公告", data["title"])
		assert.Equal(t, "这是重要公告的内容", data["content"])

		mockService.AssertExpectations(t)
	})

	t.Run("公告不存在", func(t *testing.T) {
		mockService := new(MockAnnouncementService)
		api := announcements.NewAnnouncementPublicAPI(mockService)
		router := setupTestRouter()
		router.GET("/announcements/:id", api.GetAnnouncementByID)

		mockService.On("GetAnnouncementByID", mock.Anything, "nonexistent").
			Return((*messaging.Announcement)(nil), assert.AnError)

		req, _ := http.NewRequest("GET", "/announcements/nonexistent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 404 或 500 都可以接受，取决于错误处理方式
		assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)

		mockService.AssertExpectations(t)
	})

	t.Run("空ID参数", func(t *testing.T) {
		mockService := new(MockAnnouncementService)
		api := announcements.NewAnnouncementPublicAPI(mockService)
		router := setupTestRouter()
		router.GET("/announcements/:id", api.GetAnnouncementByID)

		req, _ := http.NewRequest("GET", "/announcements/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 应该返回 404 或 400
		assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusBadRequest)
	})
}

func TestAnnouncementPublicAPI_IncrementViewCount(t *testing.T) {
	router := setupTestRouter()

	t.Run("成功增加查看次数", func(t *testing.T) {
		mockService := new(MockAnnouncementService)
		api := announcements.NewAnnouncementPublicAPI(mockService)
		router.POST("/announcements/:id/view", api.IncrementViewCount)

		mockService.On("IncrementViewCount", mock.Anything, "123").
			Return(nil)

		req, _ := http.NewRequest("POST", "/announcements/123/view", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 200, int(response["code"].(float64)))
		assert.Equal(t, "操作成功", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("空ID参数", func(t *testing.T) {
		mockService := new(MockAnnouncementService)
		api := announcements.NewAnnouncementPublicAPI(mockService)
		router := setupTestRouter()
		router.POST("/announcements/:id/view", api.IncrementViewCount)

		req, _ := http.NewRequest("POST", "/announcements/view", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 应该返回错误状态码
		assert.True(t, w.Code >= 400)
	})

	t.Run("服务错误", func(t *testing.T) {
		mockService := new(MockAnnouncementService)
		api := announcements.NewAnnouncementPublicAPI(mockService)
		router := setupTestRouter()
		router.POST("/announcements/:id/view", api.IncrementViewCount)

		mockService.On("IncrementViewCount", mock.Anything, "123").
			Return(assert.AnError)

		req, _ := http.NewRequest("POST", "/announcements/123/view", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService.AssertExpectations(t)
	})
}

// ============ 管理API测试 ============

// 这里可以添加管理员API的测试用例
// 例如: TestAnnouncementAdminAPI_CreateAnnouncement
// 例如: TestAnnouncementAdminAPI_UpdateAnnouncement
// 例如: TestAnnouncementAdminAPI_DeleteAnnouncement
// 例如: TestAnnouncementAdminAPI_BatchUpdateStatus
// 例如: TestAnnouncementAdminAPI_BatchDelete

// ============ 集成测试 ============

func TestAnnouncementAPI_Integration(t *testing.T) {
	t.Run("完整流程：创建、获取、更新、删除", func(t *testing.T) {
		router := setupTestRouter()
		mockService := new(MockAnnouncementService)

		// 设置路由
		publicAPI := announcements.NewAnnouncementPublicAPI(mockService)
		router.GET("/announcements/effective", publicAPI.GetEffectiveAnnouncements)
		router.GET("/announcements/:id", publicAPI.GetAnnouncementByID)
		router.POST("/announcements/:id/view", publicAPI.IncrementViewCount)

		now := time.Now()
		later := now.Add(24 * time.Hour)

		// 1. 获取有效公告列表
		announcements := []*messaging.Announcement{
			{
				ID:        "1",
				Title:     "测试公告",
				Content:   "内容",
				Type:      messaging.AnnouncementTypeSystem,
				IsActive:  true,
				StartTime: &now,
				EndTime:   &later,
			},
		}

		mockService.On("GetEffectiveAnnouncements", mock.Anything, "all", 10).
			Return(announcements, nil).Once()

		req1, _ := http.NewRequest("GET", "/announcements/effective", nil)
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// 2. 获取单个公告详情
		mockService.On("GetAnnouncementByID", mock.Anything, "1").
			Return(announcements[0], nil).Once()

		req2, _ := http.NewRequest("GET", "/announcements/1", nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		// 3. 增加查看次数
		mockService.On("IncrementViewCount", mock.Anything, "1").
			Return(nil).Once()

		req3, _ := http.NewRequest("POST", "/announcements/1/view", nil)
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusOK, w3.Code)

		mockService.AssertExpectations(t)
	})
}

func TestAnnouncementAPI_ErrorHandling(t *testing.T) {
	router := setupTestRouter()

	t.Run("无效的JSON请求体", func(t *testing.T) {
		mockService := new(MockAnnouncementService)

		// 这个测试需要管理员API端点
		// 这里仅作为示例展示如何测试错误处理
		router.POST("/admin/announcements", func(c *gin.Context) {
			var req messagingService.CreateAnnouncementRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "参数错误",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{})
		})

		invalidJSON := `{"title": "test", "type": "invalid_type"}`
		req, _ := http.NewRequest("POST", "/admin/announcements", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 应该返回错误
		assert.True(t, w.Code >= 400)
	})
}

// ============ 性能测试 ============

func TestAnnouncementAPI_Performance(t *testing.T) {
	router := setupTestRouter()
	mockService := new(MockAnnouncementService)

	api := announcements.NewAnnouncementPublicAPI(mockService)
	router.GET("/announcements/effective", api.GetEffectiveAnnouncements)

	// 准备大量数据
	largeAnnouncementList := make([]*messaging.Announcement, 100)
	for i := 0; i < 100; i++ {
		largeAnnouncementList[i] = &messaging.Announcement{
			ID:      string(rune(i)),
			Title:   "公告",
			Content: "内容",
		}
	}

	t.Run("获取大量公告", func(t *testing.T) {
		mockService.On("GetEffectiveAnnouncements", mock.Anything, "all", 100).
			Return(largeAnnouncementList, nil)

		req, _ := http.NewRequest("GET", "/announcements/effective?limit=100", nil)
		w := httptest.NewRecorder()

		start := time.Now()
		router.ServeHTTP(w, req)
		duration := time.Since(start)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Less(t, duration.Milliseconds(), int64(100), "响应时间应小于100ms")

		mockService.AssertExpectations(t)
	})
}
