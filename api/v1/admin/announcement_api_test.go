package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	messagingModel "Qingyu_backend/models/messaging"
	messagingBase "Qingyu_backend/models/messaging/base"
	messagingService "Qingyu_backend/service/messaging"
	apperrors "Qingyu_backend/pkg/errors"
)

// MockAnnouncementService 模拟AnnouncementService
type MockAnnouncementService struct {
	mock.Mock
}

func (m *MockAnnouncementService) GetAnnouncementByID(ctx context.Context, id string) (*messagingModel.Announcement, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*messagingModel.Announcement), args.Error(1)
}

func (m *MockAnnouncementService) GetAnnouncements(ctx context.Context, req *messagingService.GetAnnouncementsRequest) (*messagingService.GetAnnouncementsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*messagingService.GetAnnouncementsResponse), args.Error(1)
}

func (m *MockAnnouncementService) GetEffectiveAnnouncements(ctx context.Context, targetUsers string, limit int) ([]*messagingModel.Announcement, error) {
	args := m.Called(ctx, targetUsers, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*messagingModel.Announcement), args.Error(1)
}

func (m *MockAnnouncementService) CreateAnnouncement(ctx context.Context, req *messagingService.CreateAnnouncementRequest) (*messagingModel.Announcement, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*messagingModel.Announcement), args.Error(1)
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

// setupAnnouncementAPITestRouter 设置公告管理测试路由
func setupAnnouncementAPITestRouter(service *MockAnnouncementService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := NewAnnouncementAPI(service)

	v1 := r.Group("/api/v1/admin")
	{
		v1.GET("/announcements", api.GetAnnouncements)
		v1.GET("/announcements/:id", api.GetAnnouncementByID)
		v1.POST("/announcements", api.CreateAnnouncement)
		v1.PUT("/announcements/:id", api.UpdateAnnouncement)
		v1.DELETE("/announcements/:id", api.DeleteAnnouncement)
		v1.PUT("/announcements/batch-status", api.BatchUpdateStatus)
		v1.DELETE("/announcements/batch-delete", api.BatchDelete)
	}

	return r
}

// ==================== GetAnnouncementByID Tests ====================

func TestAnnouncementAPI_GetAnnouncementByID_Success(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	announcementID := primitive.NewObjectID().Hex()
	objID := primitive.NewObjectID()
	title := "测试公告"
	announcementType := messagingModel.AnnouncementTypeInfo
	expectedAnnouncement := &messagingModel.Announcement{
		IdentifiedEntity: messagingBase.IdentifiedEntity{ID: objID},
		TitledEntity:     messagingBase.TitledEntity{Title: title},
		Content:          "这是一个测试公告",
		Type:             announcementType,
	}

	mockService.On("GetAnnouncementByID", mock.Anything, announcementID).Return(expectedAnnouncement, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/announcements/"+announcementID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // pkg/response uses 0 for success
	assert.Equal(t, "操作成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestAnnouncementAPI_GetAnnouncementByID_InvalidID(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	invalidID := "invalid-object-id"
	mockService.On("GetAnnouncementByID", mock.Anything, invalidID).Return(
		nil,
		apperrors.BookstoreServiceFactory.ValidationError("INVALID_ID", "无效的公告ID", invalidID),
	)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/announcements/"+invalidID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(5000), response["code"])

	mockService.AssertExpectations(t)
}

func TestAnnouncementAPI_GetAnnouncementByID_NotFound(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	announcementID := primitive.NewObjectID().Hex()
	mockService.On("GetAnnouncementByID", mock.Anything, announcementID).Return(
		nil,
		apperrors.BookstoreServiceFactory.NotFoundError("Announcement", announcementID),
	)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/announcements/"+announcementID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(5000), response["code"])

	mockService.AssertExpectations(t)
}

func TestAnnouncementAPI_GetAnnouncementByID_InternalError(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	announcementID := primitive.NewObjectID().Hex()
	mockService.On("GetAnnouncementByID", mock.Anything, announcementID).Return(
		nil,
		apperrors.BookstoreServiceFactory.InternalError("ANNOUNCEMENT_GET_FAILED", "获取公告失败", errors.New("database error")),
	)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/announcements/"+announcementID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(5000), response["code"])

	mockService.AssertExpectations(t)
}

func TestAnnouncementAPI_GetAnnouncementByID_EmptyID(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/announcements/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	// Gin redirects trailing slash, so we expect 301
	assert.Equal(t, http.StatusMovedPermanently, w.Code)
}

// ==================== GetAnnouncements Tests ====================

func TestAnnouncementAPI_GetAnnouncements_Success(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	objID1 := primitive.NewObjectID()
	objID2 := primitive.NewObjectID()
	title1 := "公告1"
	title2 := "公告2"
	expectedResponse := &messagingService.GetAnnouncementsResponse{
		Announcements: []*messagingModel.Announcement{
			{IdentifiedEntity: messagingBase.IdentifiedEntity{ID: objID1}, TitledEntity: messagingBase.TitledEntity{Title: title1}, Type: messagingModel.AnnouncementTypeInfo},
			{IdentifiedEntity: messagingBase.IdentifiedEntity{ID: objID2}, TitledEntity: messagingBase.TitledEntity{Title: title2}, Type: messagingModel.AnnouncementTypeWarning},
		},
		Total: 2,
	}

	mockService.On("GetAnnouncements", mock.Anything, mock.AnythingOfType("*messaging.GetAnnouncementsRequest")).Return(expectedResponse, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/announcements", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	mockService.AssertExpectations(t)
}

func TestAnnouncementAPI_GetAnnouncements_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	mockService.On("GetAnnouncements", mock.Anything, mock.AnythingOfType("*messaging.GetAnnouncementsRequest")).Return(
		nil,
		apperrors.BookstoreServiceFactory.InternalError("GET_ANNOUNCEMENTS_FAILED", "获取公告列表失败", errors.New("database error")),
	)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/announcements", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== CreateAnnouncement Tests ====================

func TestAnnouncementAPI_CreateAnnouncement_Success(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	endTime := time.Now().Add(24 * time.Hour)
	title := "新公告"
	req := messagingService.CreateAnnouncementRequest{
		Title:       title,
		Content:     "公告内容",
		Type:        string(messagingModel.AnnouncementTypeInfo),
		TargetRole:  "all",
		Priority:    1,
		IsActive:    true,
		EndTime:     &endTime,
	}

	objID := primitive.NewObjectID()
	expectedAnnouncement := &messagingModel.Announcement{
		IdentifiedEntity: messagingBase.IdentifiedEntity{ID: objID},
		TitledEntity:     messagingBase.TitledEntity{Title: title},
		Content:          req.Content,
		Type:             messagingModel.AnnouncementTypeInfo,
	}

	mockService.On("CreateAnnouncement", mock.Anything, mock.MatchedBy(func(r *messagingService.CreateAnnouncementRequest) bool {
		return r.Title == req.Title && r.Content == req.Content
	})).Return(expectedAnnouncement, nil)

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("POST", "/api/v1/admin/announcements", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	mockService.AssertExpectations(t)
}

func TestAnnouncementAPI_CreateAnnouncement_InvalidJSON(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/admin/announcements", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnnouncementAPI_CreateAnnouncement_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	title := "新公告"
	req := messagingService.CreateAnnouncementRequest{
		Title:      title,
		Content:    "公告内容",
		Type:       string(messagingModel.AnnouncementTypeInfo),
		TargetRole: "all",
		Priority:   1,
		IsActive:   true,
	}

	mockService.On("CreateAnnouncement", mock.Anything, mock.AnythingOfType("*messaging.CreateAnnouncementRequest")).Return(
		nil,
		apperrors.BookstoreServiceFactory.ValidationError("INVALID_TITLE", "标题格式错误", "标题太短"),
	)

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("POST", "/api/v1/admin/announcements", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== UpdateAnnouncement Tests ====================

func TestAnnouncementAPI_UpdateAnnouncement_Success(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	announcementID := primitive.NewObjectID().Hex()
	title := "更新后的标题"
	content := "更新后的内容"
	req := messagingService.UpdateAnnouncementRequest{
		Title:   &title,
		Content: &content,
	}

	mockService.On("UpdateAnnouncement", mock.Anything, announcementID, mock.MatchedBy(func(r *messagingService.UpdateAnnouncementRequest) bool {
		return *r.Title == *req.Title && *r.Content == *req.Content
	})).Return(nil)

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("PUT", "/api/v1/admin/announcements/"+announcementID, bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestAnnouncementAPI_UpdateAnnouncement_NotFound(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	announcementID := primitive.NewObjectID().Hex()
	title := "更新后的标题"
	req := messagingService.UpdateAnnouncementRequest{
		Title: &title,
	}

	mockService.On("UpdateAnnouncement", mock.Anything, announcementID, mock.AnythingOfType("*messaging.UpdateAnnouncementRequest")).Return(
		apperrors.BookstoreServiceFactory.NotFoundError("Announcement", announcementID),
	)

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("PUT", "/api/v1/admin/announcements/"+announcementID, bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== DeleteAnnouncement Tests ====================

func TestAnnouncementAPI_DeleteAnnouncement_Success(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	announcementID := primitive.NewObjectID().Hex()
	mockService.On("DeleteAnnouncement", mock.Anything, announcementID).Return(nil)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/announcements/"+announcementID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestAnnouncementAPI_DeleteAnnouncement_NotFound(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	announcementID := primitive.NewObjectID().Hex()
	mockService.On("DeleteAnnouncement", mock.Anything, announcementID).Return(
		apperrors.BookstoreServiceFactory.NotFoundError("Announcement", announcementID),
	)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/announcements/"+announcementID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== BatchUpdateStatus Tests ====================

func TestAnnouncementAPI_BatchUpdateStatus_Success(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	req := messagingService.BatchUpdateAnnouncementStatusRequest{
		AnnouncementIDs: []string{primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()},
		IsActive: false,
	}

	mockService.On("BatchUpdateStatus", mock.Anything, mock.MatchedBy(func(r *messagingService.BatchUpdateAnnouncementStatusRequest) bool {
		return len(r.AnnouncementIDs) == len(req.AnnouncementIDs) && r.IsActive == req.IsActive
	})).Return(nil)

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("PUT", "/api/v1/admin/announcements/batch-status", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestAnnouncementAPI_BatchUpdateStatus_ValidationError(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	req := messagingService.BatchUpdateAnnouncementStatusRequest{
		AnnouncementIDs: []string{"invalid-id"},
		IsActive: true,
	}

	mockService.On("BatchUpdateStatus", mock.Anything, mock.AnythingOfType("*messaging.BatchUpdateAnnouncementStatusRequest")).Return(
		apperrors.BookstoreServiceFactory.ValidationError("INVALID_IDS", "无效的ID列表", "包含无效的ObjectID"),
	)

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("PUT", "/api/v1/admin/announcements/batch-status", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== BatchDelete Tests ====================

func TestAnnouncementAPI_BatchDelete_Success(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	req := BatchDeleteRequest{
		IDs: []string{primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()},
	}

	mockService.On("BatchDelete", mock.Anything, req.IDs).Return(nil)

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("DELETE", "/api/v1/admin/announcements/batch-delete", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestAnnouncementAPI_BatchDelete_EmptyIDs(t *testing.T) {
	// Given
	mockService := new(MockAnnouncementService)
	router := setupAnnouncementAPITestRouter(mockService)

	req := BatchDeleteRequest{
		IDs: []string{},
	}

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("DELETE", "/api/v1/admin/announcements/batch-delete", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
