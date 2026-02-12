package admin

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

	"Qingyu_backend/service/base"
	"Qingyu_backend/service/events"
)

// MockPersistedEventBus 模拟PersistedEventBus
type MockPersistedEventBus struct {
	mock.Mock
}

func (m *MockPersistedEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockPersistedEventBus) Unsubscribe(eventType string, handlerName string) error {
	args := m.Called(eventType, handlerName)
	return args.Error(0)
}

func (m *MockPersistedEventBus) Publish(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockPersistedEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Replay 模拟Replay方法
func (m *MockPersistedEventBus) Replay(ctx context.Context, handler base.EventHandler, filter events.EventFilter) (*events.ReplayResult, error) {
	args := m.Called(ctx, handler, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*events.ReplayResult), args.Error(1)
}

// setupEventsAdminTestRouter 设置事件管理测试路由
func setupEventsAdminTestRouter(eventBus *MockPersistedEventBus) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := NewEventsAdminAPI(eventBus)

	v1 := r.Group("/api/v1/admin/events")
	{
		v1.POST("/replay", api.ReplayEvents)
	}

	return r
}

// ==================== ReplayEvents Tests ====================

func TestEventsAPI_Replay_Success(t *testing.T) {
	// Given
	mockEventBus := new(MockPersistedEventBus)
	router := setupEventsAdminTestRouter(mockEventBus)

	reqBody := ReplayEventsRequest{
		EventType: "test_event",
		Limit:     100,
		DryRun:    false,
	}

	expectedResult := &events.ReplayResult{
		ReplayedCount: 50,
		FailedCount:   0,
		SkippedCount:  0,
		Duration:      100 * time.Millisecond,
	}

	mockEventBus.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(filter events.EventFilter) bool {
		return filter.EventType == "test_event" && filter.Limit == 100 && !filter.DryRun
	})).Return(expectedResult, nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/events/replay", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "操作成功", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(50), data["replayed_count"])
	assert.Equal(t, float64(0), data["failed_count"])
	assert.Equal(t, float64(0), data["skipped_count"])

	mockEventBus.AssertExpectations(t)
}

func TestEventsAPI_Replay_DryRun(t *testing.T) {
	// Given
	mockEventBus := new(MockPersistedEventBus)
	router := setupEventsAdminTestRouter(mockEventBus)

	reqBody := ReplayEventsRequest{
		Limit:  1000,
		DryRun: true,
	}

	expectedResult := &events.ReplayResult{
		ReplayedCount: 0,
		FailedCount:   0,
		SkippedCount:  1000,
		Duration:      50 * time.Millisecond,
	}

	mockEventBus.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(filter events.EventFilter) bool {
		return filter.DryRun && filter.Limit == 1000
	})).Return(expectedResult, nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/events/replay", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1000), data["skipped_count"])

	mockEventBus.AssertExpectations(t)
}

func TestEventsAPI_Replay_WithEventTypeFilter(t *testing.T) {
	// Given
	mockEventBus := new(MockPersistedEventBus)
	router := setupEventsAdminTestRouter(mockEventBus)

	reqBody := ReplayEventsRequest{
		EventType: "user_created",
		Limit:     500,
	}

	expectedResult := &events.ReplayResult{
		ReplayedCount: 30,
		FailedCount:   0,
		SkippedCount:  0,
		Duration:      80 * time.Millisecond,
	}

	mockEventBus.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(filter events.EventFilter) bool {
		return filter.EventType == "user_created" && filter.Limit == 500
	})).Return(expectedResult, nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/events/replay", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockEventBus.AssertExpectations(t)
}

func TestEventsAPI_Replay_WithTimeRange(t *testing.T) {
	// Given
	mockEventBus := new(MockPersistedEventBus)
	router := setupEventsAdminTestRouter(mockEventBus)

	fromTime := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	toTime := time.Now().Format(time.RFC3339)

	reqBody := ReplayEventsRequest{
		From:  fromTime,
		To:    toTime,
		Limit: 1000,
	}

	expectedResult := &events.ReplayResult{
		ReplayedCount: 200,
		FailedCount:   5,
		SkippedCount:  0,
		Duration:      200 * time.Millisecond,
	}

	mockEventBus.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(filter events.EventFilter) bool {
		if filter.StartTime == nil || filter.EndTime == nil {
			return false
		}
		return filter.Limit == 1000
	})).Return(expectedResult, nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/events/replay", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(200), data["replayed_count"])
	assert.Equal(t, float64(5), data["failed_count"])

	mockEventBus.AssertExpectations(t)
}

func TestEventsAPI_Replay_InvalidLimit_TooSmall(t *testing.T) {
	// Given
	mockEventBus := new(MockPersistedEventBus)
	router := setupEventsAdminTestRouter(mockEventBus)

	reqBody := ReplayEventsRequest{
		Limit: -1, // 使用负数来测试小于1的情况
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/events/replay", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEqual(t, float64(0), response["code"])
}

func TestEventsAPI_Replay_InvalidLimit_TooLarge(t *testing.T) {
	// Given
	mockEventBus := new(MockPersistedEventBus)
	router := setupEventsAdminTestRouter(mockEventBus)

	reqBody := ReplayEventsRequest{
		Limit: 10001,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/events/replay", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEqual(t, float64(0), response["code"])
}

func TestEventsAPI_Replay_InvalidTimeFormat(t *testing.T) {
	// Given
	mockEventBus := new(MockPersistedEventBus)
	router := setupEventsAdminTestRouter(mockEventBus)

	reqBody := ReplayEventsRequest{
		From: "invalid-time-format",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/events/replay", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEqual(t, float64(0), response["code"])
}

func TestEventsAPI_Replay_InvalidJSON(t *testing.T) {
	// Given
	mockEventBus := new(MockPersistedEventBus)
	router := setupEventsAdminTestRouter(mockEventBus)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/admin/events/replay", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventsAPI_Replay_ServiceError(t *testing.T) {
	// Given
	mockEventBus := new(MockPersistedEventBus)
	router := setupEventsAdminTestRouter(mockEventBus)

	reqBody := ReplayEventsRequest{
		Limit: 100,
	}

	mockEventBus.On("Replay", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/events/replay", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockEventBus.AssertExpectations(t)
}

func TestEventsAPI_Replay_DefaultLimit(t *testing.T) {
	// Given
	mockEventBus := new(MockPersistedEventBus)
	router := setupEventsAdminTestRouter(mockEventBus)

	reqBody := ReplayEventsRequest{} // 不设置Limit，使用默认值

	expectedResult := &events.ReplayResult{
		ReplayedCount: 1000,
		FailedCount:   0,
		SkippedCount:  0,
		Duration:      150 * time.Millisecond,
	}

	mockEventBus.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(filter events.EventFilter) bool {
		return filter.Limit == 1000 // 默认值
	})).Return(expectedResult, nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/events/replay", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockEventBus.AssertExpectations(t)
}

func TestEventsAPI_Replay_WithAllFilters(t *testing.T) {
	// Given
	mockEventBus := new(MockPersistedEventBus)
	router := setupEventsAdminTestRouter(mockEventBus)

	fromTime := time.Now().Add(-7 * 24 * time.Hour).Format(time.RFC3339)
	toTime := time.Now().Format(time.RFC3339)

	reqBody := ReplayEventsRequest{
		EventType: "chapter_published",
		From:      fromTime,
		To:        toTime,
		Limit:     5000,
		DryRun:    true,
	}

	expectedResult := &events.ReplayResult{
		ReplayedCount: 0,
		FailedCount:   0,
		SkippedCount:  5000,
		Duration:      300 * time.Millisecond,
	}

	mockEventBus.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(filter events.EventFilter) bool {
		return filter.EventType == "chapter_published" &&
			filter.StartTime != nil &&
			filter.EndTime != nil &&
			filter.Limit == 5000 &&
			filter.DryRun
	})).Return(expectedResult, nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/events/replay", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(300), data["duration_ms"]) // 300ms

	mockEventBus.AssertExpectations(t)
}
