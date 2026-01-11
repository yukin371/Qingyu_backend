package reader

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
	"go.mongodb.org/mongo-driver/bson/primitive"

	progressSync "Qingyu_backend/pkg/sync"
	"Qingyu_backend/service/interfaces"
	ws "Qingyu_backend/pkg/websocket"
)

// MockProgressSyncService 模拟ProgressSyncService
type MockProgressSyncService struct {
	mock.Mock
}

func (m *MockProgressSyncService) GetHub() *ws.ProgressHub {
	// Return a mock hub - for testing we don't actually need to connect to it
	return ws.NewProgressHub()
}

func (m *MockProgressSyncService) SyncProgress(ctx context.Context, userID, bookID, chapterID, deviceID string, progress float64) error {
	args := m.Called(ctx, userID, bookID, chapterID, deviceID, progress)
	return args.Error(0)
}

func (m *MockProgressSyncService) MergeOfflineProgresses(ctx context.Context, userID string, progresses []progressSync.OfflineProgress) error {
	args := m.Called(ctx, userID, progresses)
	return args.Error(0)
}

func (m *MockProgressSyncService) GetSyncStatus(userID string) *progressSync.SyncStatus {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return &progressSync.SyncStatus{
			UserID:          userID,
			ConnectedDevices: []string{},
			DeviceCount:      0,
			IsSyncing:        false,
		}
	}
	return args.Get(0).(*progressSync.SyncStatus)
}

func setupSyncTestRouter(syncService interfaces.ProgressSyncService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userId", userID)
		}
		c.Next()
	})

	api := NewSyncAPI(syncService)

	v1 := r.Group("/api/v1/reader/progress")
	{
		v1.POST("/sync", api.SyncProgress)
		v1.POST("/merge", api.MergeOfflineProgresses)
		v1.GET("/sync-status", api.GetSyncStatus)
	}

	return r
}

func TestSyncAPI_SyncProgress_Success(t *testing.T) {
	// Given
	mockService := new(MockProgressSyncService)
	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)

	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	deviceID := "test-device"

	reqBody := SyncProgressRequest{
		BookID:    bookID,
		ChapterID: chapterID,
		Progress:  0.5,
		DeviceID:  deviceID,
	}

	mockService.On("SyncProgress", mock.Anything, mock.AnythingOfType("string"), bookID, chapterID, deviceID, 0.5).Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress/sync", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSyncAPI_SyncProgress_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockProgressSyncService)
	router := setupSyncTestRouter(mockService, "") // No userID

	reqBody := SyncProgressRequest{
		BookID:    primitive.NewObjectID().Hex(),
		ChapterID: primitive.NewObjectID().Hex(),
		Progress:  0.5,
		DeviceID:  "test-device",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress/sync", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSyncAPI_SyncProgress_MissingRequiredFields(t *testing.T) {
	// Given
	mockService := new(MockProgressSyncService)
	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"progress": 0.5,
		// Missing bookId, chapterId, deviceId
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress/sync", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSyncAPI_SyncProgress_InvalidProgress(t *testing.T) {
	// Given
	mockService := new(MockProgressSyncService)
	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)

	reqBody := SyncProgressRequest{
		BookID:    primitive.NewObjectID().Hex(),
		ChapterID: primitive.NewObjectID().Hex(),
		Progress:  1.5, // Invalid - should be between 0 and 1
		DeviceID:  "test-device",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress/sync", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSyncAPI_SyncProgress_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockProgressSyncService)
	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)

	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	deviceID := "test-device"

	reqBody := SyncProgressRequest{
		BookID:    bookID,
		ChapterID: chapterID,
		Progress:  0.5,
		DeviceID:  deviceID,
	}

	mockService.On("SyncProgress", mock.Anything, mock.AnythingOfType("string"), bookID, chapterID, deviceID, 0.5).Return(assert.AnError)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress/sync", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestSyncAPI_MergeOfflineProgresses_Success(t *testing.T) {
	// Given
	mockService := new(MockProgressSyncService)
	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)

	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()

	reqBody := MergeProgressRequest{
		Progresses: []OfflineProgressItem{
			{
				BookID:    bookID,
				ChapterID: chapterID,
				Progress:  0.3,
				Timestamp: time.Now().Format(time.RFC3339),
				DeviceID:  "device-1",
			},
			{
				BookID:    bookID,
				ChapterID: chapterID,
				Progress:  0.7,
				Timestamp: time.Now().Add(1 * time.Minute).Format(time.RFC3339),
				DeviceID:  "device-2",
			},
		},
	}

	mockService.On("MergeOfflineProgresses", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]sync.OfflineProgress")).Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress/merge", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSyncAPI_MergeOfflineProgresses_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockProgressSyncService)
	router := setupSyncTestRouter(mockService, "") // No userID

	reqBody := MergeProgressRequest{
		Progresses: []OfflineProgressItem{},
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress/merge", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSyncAPI_MergeOfflineProgresses_InvalidTimestamp(t *testing.T) {
	// Given
	mockService := new(MockProgressSyncService)
	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)

	reqBody := MergeProgressRequest{
		Progresses: []OfflineProgressItem{
			{
				BookID:    primitive.NewObjectID().Hex(),
				ChapterID: primitive.NewObjectID().Hex(),
				Progress:  0.5,
				Timestamp: "invalid-timestamp",
				DeviceID:  "device-1",
			},
		},
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress/merge", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSyncAPI_GetSyncStatus_Success(t *testing.T) {
	// Given
	mockService := new(MockProgressSyncService)
	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)

	expectedStatus := &progressSync.SyncStatus{
		UserID:          userID,
		ConnectedDevices: []string{"device-1", "device-2"},
		DeviceCount:      2,
		IsSyncing:        true,
	}

	mockService.On("GetSyncStatus", userID).Return(expectedStatus)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/progress/sync-status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSyncAPI_GetSyncStatus_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockProgressSyncService)
	router := setupSyncTestRouter(mockService, "") // No userID

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/progress/sync-status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSyncAPI_MergeOfflineProgresses_EmptyProgresses(t *testing.T) {
	// Given
	mockService := new(MockProgressSyncService)
	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)

	reqBody := MergeProgressRequest{
		Progresses: []OfflineProgressItem{},
	}

	mockService.On("MergeOfflineProgresses", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]sync.OfflineProgress")).Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress/merge", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
