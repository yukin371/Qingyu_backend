package reader

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	progressSync "Qingyu_backend/pkg/sync"
	ws "Qingyu_backend/pkg/websocket"
	"Qingyu_backend/service/interfaces"
)

// MockProgressSyncService 模拟ProgressSyncService
type MockProgressSyncService struct {
	mock.Mock
	hub *ws.ProgressHub
}

func (m *MockProgressSyncService) GetHub() *ws.ProgressHub {
	if m.hub == nil {
		m.hub = ws.NewProgressHub()
	}
	return m.hub
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
	for _, call := range m.ExpectedCalls {
		if call.Method == "GetSyncStatus" {
			args := m.Called(userID)
			if args.Get(0) == nil {
				return &progressSync.SyncStatus{
					UserID:           userID,
					ConnectedDevices: []string{},
					DeviceCount:      0,
					IsSyncing:        false,
				}
			}
			return args.Get(0).(*progressSync.SyncStatus)
		}
	}

	devices := m.GetHub().GetConnectedDevices(userID)
	return &progressSync.SyncStatus{
		UserID:           userID,
		ConnectedDevices: devices,
		DeviceCount:      len(devices),
		IsSyncing:        len(devices) > 1,
	}
}

type syncStatusResponse struct {
	Data progressSync.SyncStatus `json:"data"`
}

func getSyncStatusResponse(t *testing.T, router *gin.Engine) progressSync.SyncStatus {
	t.Helper()

	req, err := http.NewRequest("GET", "/api/v1/reader/progress/sync-status", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var responseBody syncStatusResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &responseBody))

	return responseBody.Data
}

func assertSyncStatusMatches(
	t *testing.T,
	status progressSync.SyncStatus,
	userID string,
	expectedDevices []string,
	expectedCount int,
	expectedSyncing bool,
) {
	t.Helper()

	assert.Equal(t, userID, status.UserID)
	assert.Len(t, status.ConnectedDevices, expectedCount)
	assert.Equal(t, expectedCount, status.DeviceCount)
	assert.Equal(t, expectedSyncing, status.IsSyncing)
	for _, deviceID := range expectedDevices {
		assert.Contains(t, status.ConnectedDevices, deviceID)
	}
}

func waitForSyncStatus(
	t *testing.T,
	router *gin.Engine,
	userID string,
	expectedDevices []string,
	expectedCount int,
	expectedSyncing bool,
) {
	t.Helper()

	assert.Eventually(t, func() bool {
		status := getSyncStatusResponse(t, router)
		if status.UserID != userID || status.DeviceCount != expectedCount || status.IsSyncing != expectedSyncing {
			return false
		}
		if len(status.ConnectedDevices) != expectedCount {
			return false
		}
		for _, deviceID := range expectedDevices {
			if !assert.Contains(t, status.ConnectedDevices, deviceID) {
				return false
			}
		}
		return true
	}, time.Second, 10*time.Millisecond)
}

func setupSyncTestRouter(syncService interfaces.ProgressSyncService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 错误处理中间件
	r.Use(func(c *gin.Context) {
		c.Next()
		// 检查是否有错误
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": c.Errors.String(),
			})
		}
	})

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := NewSyncAPI(syncService)

	v1 := r.Group("/api/v1/reader/progress")
	{
		v1.GET("/ws", api.SyncWebSocket)
		v1.POST("/sync", api.SyncProgress)
		v1.POST("/merge", api.MergeOfflineProgresses)
		v1.GET("/sync-status", api.GetSyncStatus)
	}

	return r
}

func waitForConnectedDevices(t *testing.T, hub *ws.ProgressHub, userID string, expected []string) {
	t.Helper()
	assert.Eventually(t, func() bool {
		devices := hub.GetConnectedDevices(userID)
		if len(devices) != len(expected) {
			return false
		}
		for _, device := range expected {
			if !assert.Contains(t, devices, device) {
				return false
			}
		}
		return true
	}, time.Second, 10*time.Millisecond)
}

func closeSyncWebSocket(t *testing.T, conn *websocket.Conn, hub *ws.ProgressHub, userID string) {
	t.Helper()
	if conn == nil {
		return
	}

	_ = conn.Close()
	waitForConnectedDevices(t, hub, userID, []string{})
}

func dialSyncWebSocket(t *testing.T, serverURL string, headers map[string]string) (*websocket.Conn, *http.Response, error) {
	t.Helper()

	dialer := websocket.DefaultDialer
	wsURL := strings.Replace(serverURL, "http", "ws", 1) + "/api/v1/reader/progress/ws"
	requestHeader := http.Header{}
	for key, value := range headers {
		requestHeader.Set(key, value)
	}

	return dialer.Dial(wsURL, requestHeader)
}

func performMergeOfflineProgressesRequest(t *testing.T, router *gin.Engine, body MergeProgressRequest) *httptest.ResponseRecorder {
	t.Helper()

	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/reader/progress/merge", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
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

func TestSyncAPI_SyncWebSocket_Unauthorized(t *testing.T) {
	mockService := new(MockProgressSyncService)
	hub := mockService.GetHub()
	go hub.Run()

	router := setupSyncTestRouter(mockService, "")
	server := httptest.NewServer(router)
	defer server.Close()

	conn, resp, err := dialSyncWebSocket(t, server.URL, nil)
	if conn != nil {
		defer conn.Close()
	}

	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestSyncAPI_SyncWebSocket_DefaultsMissingDeviceIDToUnknown(t *testing.T) {
	mockService := new(MockProgressSyncService)
	hub := mockService.GetHub()
	go hub.Run()

	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)
	server := httptest.NewServer(router)
	defer server.Close()

	conn, resp, err := dialSyncWebSocket(t, server.URL, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
	t.Cleanup(func() {
		closeSyncWebSocket(t, conn, hub, userID)
	})

	waitForConnectedDevices(t, hub, userID, []string{"unknown"})
}

func TestSyncAPI_SyncWebSocket_ReconnectSameDeviceDoesNotDuplicateStatus(t *testing.T) {
	mockService := new(MockProgressSyncService)
	hub := mockService.GetHub()
	go hub.Run()

	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)
	server := httptest.NewServer(router)
	defer server.Close()

	headers := map[string]string{"X-Device-ID": "device-1"}
	conn1, _, err := dialSyncWebSocket(t, server.URL, headers)
	require.NoError(t, err)
	waitForConnectedDevices(t, hub, userID, []string{"device-1"})

	closeSyncWebSocket(t, conn1, hub, userID)

	conn2, _, err := dialSyncWebSocket(t, server.URL, headers)
	require.NoError(t, err)
	t.Cleanup(func() {
		closeSyncWebSocket(t, conn2, hub, userID)
	})

	waitForConnectedDevices(t, hub, userID, []string{"device-1"})
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

func TestSyncAPI_MergeOfflineProgresses_ConvertsTimestampsAndUserID(t *testing.T) {
	mockService := new(MockProgressSyncService)
	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)

	firstTimestamp := time.Now().Add(-2 * time.Minute).UTC().Truncate(time.Second)
	secondTimestamp := time.Now().UTC().Truncate(time.Second)
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()

	reqBody := MergeProgressRequest{
		Progresses: []OfflineProgressItem{
			{
				BookID:    bookID,
				ChapterID: chapterID,
				Progress:  0.25,
				Timestamp: firstTimestamp.Format(time.RFC3339),
				DeviceID:  "device-1",
			},
			{
				BookID:    bookID,
				ChapterID: chapterID,
				Progress:  0.8,
				Timestamp: secondTimestamp.Format(time.RFC3339),
				DeviceID:  "device-2",
			},
		},
	}

	mockService.On(
		"MergeOfflineProgresses",
		mock.Anything,
		userID,
		mock.MatchedBy(func(items []progressSync.OfflineProgress) bool {
			require.Len(t, items, 2)
			return items[0].UserID == userID &&
				items[0].Timestamp.Equal(firstTimestamp) &&
				items[1].UserID == userID &&
				items[1].Timestamp.Equal(secondTimestamp)
		}),
	).Return(nil)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress/merge", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

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

func TestSyncAPI_MergeOfflineProgresses_ServiceError(t *testing.T) {
	mockService := new(MockProgressSyncService)
	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)

	reqBody := MergeProgressRequest{
		Progresses: []OfflineProgressItem{
			{
				BookID:    primitive.NewObjectID().Hex(),
				ChapterID: primitive.NewObjectID().Hex(),
				Progress:  0.5,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				DeviceID:  "device-1",
			},
		},
	}

	mockService.On("MergeOfflineProgresses", mock.Anything, userID, mock.Anything).Return(assert.AnError)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress/merge", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestSyncAPI_GetSyncStatus_Success(t *testing.T) {
	// Given
	mockService := new(MockProgressSyncService)
	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)

	expectedStatus := &progressSync.SyncStatus{
		UserID:           userID,
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

func TestSyncAPI_GetSyncStatus_NoConnectionsReturnsIdle(t *testing.T) {
	mockService := new(MockProgressSyncService)
	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)

	status := getSyncStatusResponse(t, router)

	assertSyncStatusMatches(t, status, userID, []string{}, 0, false)
}

func TestSyncAPI_GetSyncStatus_TracksWebSocketLifecycle(t *testing.T) {
	mockService := new(MockProgressSyncService)
	hub := mockService.GetHub()
	go hub.Run()

	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)
	server := httptest.NewServer(router)
	defer server.Close()

	waitForSyncStatus(t, router, userID, []string{}, 0, false)

	conn1, _, err := dialSyncWebSocket(t, server.URL, map[string]string{"X-Device-ID": "device-1"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn1.Close()
	})
	waitForConnectedDevices(t, hub, userID, []string{"device-1"})
	waitForSyncStatus(t, router, userID, []string{"device-1"}, 1, false)

	conn2, _, err := dialSyncWebSocket(t, server.URL, map[string]string{"X-Device-ID": "device-2"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn2.Close()
	})
	waitForConnectedDevices(t, hub, userID, []string{"device-1", "device-2"})
	waitForSyncStatus(t, router, userID, []string{"device-1", "device-2"}, 2, true)

	_ = conn2.Close()
	waitForConnectedDevices(t, hub, userID, []string{"device-1"})
	waitForSyncStatus(t, router, userID, []string{"device-1"}, 1, false)

	_ = conn1.Close()
	waitForConnectedDevices(t, hub, userID, []string{})
	waitForSyncStatus(t, router, userID, []string{}, 0, false)
}

func TestSyncAPI_GetSyncStatus_AfterMergeOfflineProgresses_PreservesConnectionState(t *testing.T) {
	mockService := new(MockProgressSyncService)
	hub := mockService.GetHub()
	go hub.Run()

	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)
	server := httptest.NewServer(router)
	defer server.Close()

	conn, _, err := dialSyncWebSocket(t, server.URL, map[string]string{"X-Device-ID": "device-1"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	waitForConnectedDevices(t, hub, userID, []string{"device-1"})
	waitForSyncStatus(t, router, userID, []string{"device-1"}, 1, false)

	reqBody := MergeProgressRequest{
		Progresses: []OfflineProgressItem{
			{
				BookID:    primitive.NewObjectID().Hex(),
				ChapterID: primitive.NewObjectID().Hex(),
				Progress:  0.65,
				Timestamp: time.Now().UTC().Truncate(time.Second).Format(time.RFC3339),
				DeviceID:  "device-1",
			},
		},
	}

	mockService.On("MergeOfflineProgresses", mock.Anything, userID, mock.AnythingOfType("[]sync.OfflineProgress")).Return(nil).Once()

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress/merge", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
	waitForSyncStatus(t, router, userID, []string{"device-1"}, 1, false)
}

func TestSyncAPI_GetSyncStatus_IsSyncingStaysTrueAcrossTwoDeviceMerge(t *testing.T) {
	mockService := new(MockProgressSyncService)
	hub := mockService.GetHub()
	go hub.Run()

	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)
	server := httptest.NewServer(router)
	defer server.Close()

	conn1, _, err := dialSyncWebSocket(t, server.URL, map[string]string{"X-Device-ID": "device-1"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn1.Close()
	})

	conn2, _, err := dialSyncWebSocket(t, server.URL, map[string]string{"X-Device-ID": "device-2"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn2.Close()
	})

	waitForConnectedDevices(t, hub, userID, []string{"device-1", "device-2"})
	waitForSyncStatus(t, router, userID, []string{"device-1", "device-2"}, 2, true)

	firstTimestamp := time.Now().Add(-time.Minute).UTC().Truncate(time.Second)
	secondTimestamp := time.Now().UTC().Truncate(time.Second)
	reqBody := MergeProgressRequest{
		Progresses: []OfflineProgressItem{
			{
				BookID:    primitive.NewObjectID().Hex(),
				ChapterID: primitive.NewObjectID().Hex(),
				Progress:  0.35,
				Timestamp: firstTimestamp.Format(time.RFC3339),
				DeviceID:  "device-1",
			},
			{
				BookID:    primitive.NewObjectID().Hex(),
				ChapterID: primitive.NewObjectID().Hex(),
				Progress:  0.82,
				Timestamp: secondTimestamp.Format(time.RFC3339),
				DeviceID:  "device-2",
			},
		},
	}

	mockService.On(
		"MergeOfflineProgresses",
		mock.Anything,
		userID,
		mock.MatchedBy(func(items []progressSync.OfflineProgress) bool {
			require.Len(t, items, 2)
			return items[0].DeviceID == "device-1" &&
				items[1].DeviceID == "device-2" &&
				items[0].Timestamp.Equal(firstTimestamp) &&
				items[1].Timestamp.Equal(secondTimestamp)
		}),
	).Return(nil).Once()

	w := performMergeOfflineProgressesRequest(t, router, reqBody)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
	waitForSyncStatus(t, router, userID, []string{"device-1", "device-2"}, 2, true)
}

func TestSyncAPI_GetSyncStatus_OfflineMergeDoesNotAddPhantomConnectedDevice(t *testing.T) {
	mockService := new(MockProgressSyncService)
	hub := mockService.GetHub()
	go hub.Run()

	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)
	server := httptest.NewServer(router)
	defer server.Close()

	conn1, _, err := dialSyncWebSocket(t, server.URL, map[string]string{"X-Device-ID": "device-1"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn1.Close()
	})

	conn2, _, err := dialSyncWebSocket(t, server.URL, map[string]string{"X-Device-ID": "device-2"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn2.Close()
	})

	waitForConnectedDevices(t, hub, userID, []string{"device-1", "device-2"})
	waitForSyncStatus(t, router, userID, []string{"device-1", "device-2"}, 2, true)

	reqBody := MergeProgressRequest{
		Progresses: []OfflineProgressItem{
			{
				BookID:    primitive.NewObjectID().Hex(),
				ChapterID: primitive.NewObjectID().Hex(),
				Progress:  0.91,
				Timestamp: time.Now().UTC().Truncate(time.Second).Format(time.RFC3339),
				DeviceID:  "device-3",
			},
		},
	}

	mockService.On(
		"MergeOfflineProgresses",
		mock.Anything,
		userID,
		mock.MatchedBy(func(items []progressSync.OfflineProgress) bool {
			require.Len(t, items, 1)
			return items[0].DeviceID == "device-3"
		}),
	).Return(nil).Once()

	w := performMergeOfflineProgressesRequest(t, router, reqBody)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	status := getSyncStatusResponse(t, router)
	assertSyncStatusMatches(t, status, userID, []string{"device-1", "device-2"}, 2, true)
	assert.NotContains(t, status.ConnectedDevices, "device-3")
}

func TestSyncAPI_GetSyncStatus_DisconnectedDeviceMergeKeepsRemainingDeviceIdle(t *testing.T) {
	mockService := new(MockProgressSyncService)
	hub := mockService.GetHub()
	go hub.Run()

	userID := primitive.NewObjectID().Hex()
	router := setupSyncTestRouter(mockService, userID)
	server := httptest.NewServer(router)
	defer server.Close()

	conn1, _, err := dialSyncWebSocket(t, server.URL, map[string]string{"X-Device-ID": "device-1"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn1.Close()
	})

	conn2, _, err := dialSyncWebSocket(t, server.URL, map[string]string{"X-Device-ID": "device-2"})
	require.NoError(t, err)

	waitForConnectedDevices(t, hub, userID, []string{"device-1", "device-2"})
	waitForSyncStatus(t, router, userID, []string{"device-1", "device-2"}, 2, true)

	_ = conn2.Close()
	waitForConnectedDevices(t, hub, userID, []string{"device-1"})
	waitForSyncStatus(t, router, userID, []string{"device-1"}, 1, false)

	reqBody := MergeProgressRequest{
		Progresses: []OfflineProgressItem{
			{
				BookID:    primitive.NewObjectID().Hex(),
				ChapterID: primitive.NewObjectID().Hex(),
				Progress:  0.67,
				Timestamp: time.Now().UTC().Truncate(time.Second).Format(time.RFC3339),
				DeviceID:  "device-2",
			},
		},
	}

	mockService.On(
		"MergeOfflineProgresses",
		mock.Anything,
		userID,
		mock.MatchedBy(func(items []progressSync.OfflineProgress) bool {
			require.Len(t, items, 1)
			return items[0].DeviceID == "device-2"
		}),
	).Return(nil).Once()

	w := performMergeOfflineProgressesRequest(t, router, reqBody)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
	waitForSyncStatus(t, router, userID, []string{"device-1"}, 1, false)
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
