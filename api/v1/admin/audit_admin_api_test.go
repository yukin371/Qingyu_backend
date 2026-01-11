package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/audit"
	auditInterface "Qingyu_backend/service/interfaces/audit"
)

// MockContentAuditService 模拟ContentAuditService
type MockContentAuditService struct {
	mock.Mock
}

func (m *MockContentAuditService) CheckContent(ctx context.Context, content string) (*auditInterface.AuditCheckResult, error) {
	args := m.Called(ctx, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditInterface.AuditCheckResult), args.Error(1)
}

func (m *MockContentAuditService) AuditDocument(ctx context.Context, documentID string, content string, authorID string) (*audit.AuditRecord, error) {
	args := m.Called(ctx, documentID, content, authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*audit.AuditRecord), args.Error(1)
}

func (m *MockContentAuditService) GetAuditResult(ctx context.Context, targetType, targetID string) (*audit.AuditRecord, error) {
	args := m.Called(ctx, targetType, targetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*audit.AuditRecord), args.Error(1)
}

func (m *MockContentAuditService) BatchAuditDocuments(ctx context.Context, documentIDs []string) ([]*audit.AuditRecord, error) {
	args := m.Called(ctx, documentIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.AuditRecord), args.Error(1)
}

func (m *MockContentAuditService) ReviewAudit(ctx context.Context, auditID string, reviewerID string, approved bool, note string) error {
	args := m.Called(ctx, auditID, reviewerID, approved, note)
	return args.Error(0)
}

func (m *MockContentAuditService) SubmitAppeal(ctx context.Context, auditID string, authorID string, reason string) error {
	args := m.Called(ctx, auditID, authorID, reason)
	return args.Error(0)
}

func (m *MockContentAuditService) ReviewAppeal(ctx context.Context, auditID string, reviewerID string, approved bool, note string) error {
	args := m.Called(ctx, auditID, reviewerID, approved, note)
	return args.Error(0)
}

func (m *MockContentAuditService) GetUserViolations(ctx context.Context, userID string) ([]*audit.ViolationRecord, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.ViolationRecord), args.Error(1)
}

func (m *MockContentAuditService) GetUserViolationSummary(ctx context.Context, userID string) (*audit.UserViolationSummary, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*audit.UserViolationSummary), args.Error(1)
}

func (m *MockContentAuditService) GetPendingReviews(ctx context.Context, limit int) ([]*audit.AuditRecord, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.AuditRecord), args.Error(1)
}

func (m *MockContentAuditService) GetHighRiskAudits(ctx context.Context, minRiskLevel int, limit int) ([]*audit.AuditRecord, error) {
	args := m.Called(ctx, minRiskLevel, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.AuditRecord), args.Error(1)
}

func (m *MockContentAuditService) GetAuditStatistics(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

// setupAuditAdminTestRouter 设置审核管理测试路由
func setupAuditAdminTestRouter(auditService *MockContentAuditService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := NewAuditAdminAPI(auditService)

	v1 := r.Group("/api/v1/admin/audit")
	{
		v1.GET("/pending", api.GetPendingAudits)
		v1.POST("/:id/review", api.ReviewAudit)
		v1.POST("/:id/appeal/review", api.ReviewAppeal)
		v1.GET("/high-risk", api.GetHighRiskAudits)
		v1.GET("/statistics", api.GetAuditStatistics)
	}

	return r
}

// setupAuditAdminTestRouterWithAuth 设置带认证中间件的测试路由
func setupAuditAdminTestRouterWithAuth(auditService *MockContentAuditService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 模拟认证中间件
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := NewAuditAdminAPI(auditService)

	v1 := r.Group("/api/v1/admin/audit")
	{
		v1.POST("/:id/review", api.ReviewAudit)
		v1.POST("/:id/appeal/review", api.ReviewAppeal)
	}

	return r
}

// ==================== GetPendingAudits Tests ====================

func TestAuditAdminAPI_GetPendingAudits_Success(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	expectedRecords := []*audit.AuditRecord{
		{ID: primitive.NewObjectID().Hex(), TargetType: "document", Status: "pending"},
		{ID: primitive.NewObjectID().Hex(), TargetType: "chapter", Status: "pending"},
	}

	mockService.On("GetPendingReviews", mock.Anything, 50).Return(expectedRecords, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/pending", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "获取成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_GetPendingAudits_WithTargetType(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	expectedRecords := []*audit.AuditRecord{
		{ID: primitive.NewObjectID().Hex(), TargetType: "document", Status: "pending"},
	}

	mockService.On("GetPendingReviews", mock.Anything, 20).Return(expectedRecords, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/pending?targetType=document&page=1&limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_GetPendingAudits_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	mockService.On("GetPendingReviews", mock.Anything, 50).Return(nil, assert.AnError)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/pending", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== ReviewAudit Tests ====================

func TestAuditAdminAPI_ReviewAudit_Success_Approve(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	userID := "admin123"
	router := setupAuditAdminTestRouterWithAuth(mockService, userID)

	auditID := primitive.NewObjectID().Hex()
	reqBody := ReviewAuditRequest{
		Action:     "approve",
		ReviewNote: "内容符合规范",
	}

	mockService.On("ReviewAudit", mock.Anything, auditID, userID, true, "内容符合规范").Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit/"+auditID+"/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "审核已通过", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_ReviewAudit_Success_Reject(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	userID := "admin123"
	router := setupAuditAdminTestRouterWithAuth(mockService, userID)

	auditID := primitive.NewObjectID().Hex()
	reqBody := ReviewAuditRequest{
		Action:     "reject",
		ReviewNote: "内容包含违规信息",
	}

	mockService.On("ReviewAudit", mock.Anything, auditID, userID, false, "内容包含违规信息").Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit/"+auditID+"/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "审核已拒绝", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_ReviewAudit_EmptyAuditID(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	userID := "admin123"
	router := setupAuditAdminTestRouterWithAuth(mockService, userID)

	reqBody := ReviewAuditRequest{
		Action:     "approve",
		ReviewNote: "内容符合规范",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit//review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "参数错误", response["message"])
}

func TestAuditAdminAPI_ReviewAudit_InvalidJSON(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	userID := "admin123"
	router := setupAuditAdminTestRouterWithAuth(mockService, userID)

	auditID := primitive.NewObjectID().Hex()

	// When
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit/"+auditID+"/review", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuditAdminAPI_ReviewAudit_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouterWithAuth(mockService, "") // No user ID

	auditID := primitive.NewObjectID().Hex()
	reqBody := ReviewAuditRequest{
		Action:     "approve",
		ReviewNote: "内容符合规范",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit/"+auditID+"/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "未授权", response["message"])
}

func TestAuditAdminAPI_ReviewAudit_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	userID := "admin123"
	router := setupAuditAdminTestRouterWithAuth(mockService, userID)

	auditID := primitive.NewObjectID().Hex()
	reqBody := ReviewAuditRequest{
		Action:     "approve",
		ReviewNote: "内容符合规范",
	}

	mockService.On("ReviewAudit", mock.Anything, auditID, userID, true, "内容符合规范").Return(assert.AnError)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit/"+auditID+"/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== ReviewAppeal Tests ====================

func TestAuditAdminAPI_ReviewAppeal_Success_Approve(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	userID := "admin123"
	router := setupAuditAdminTestRouterWithAuth(mockService, userID)

	auditID := primitive.NewObjectID().Hex()
	reqBody := ReviewAppealRequest{
		Action:     "approve",
		ReviewNote: "申诉理由充分",
	}

	mockService.On("ReviewAppeal", mock.Anything, auditID, userID, true, "申诉理由充分").Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit/"+auditID+"/appeal/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "申诉已通过", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_ReviewAppeal_Success_Reject(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	userID := "admin123"
	router := setupAuditAdminTestRouterWithAuth(mockService, userID)

	auditID := primitive.NewObjectID().Hex()
	reqBody := ReviewAppealRequest{
		Action:     "reject",
		ReviewNote: "申诉理由不充分",
	}

	mockService.On("ReviewAppeal", mock.Anything, auditID, userID, false, "申诉理由不充分").Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit/"+auditID+"/appeal/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "申诉已驳回", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_ReviewAppeal_EmptyAuditID(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	userID := "admin123"
	router := setupAuditAdminTestRouterWithAuth(mockService, userID)

	reqBody := ReviewAppealRequest{
		Action:     "approve",
		ReviewNote: "申诉理由充分",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit//appeal/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "参数错误", response["message"])
}

func TestAuditAdminAPI_ReviewAppeal_InvalidJSON(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	userID := "admin123"
	router := setupAuditAdminTestRouterWithAuth(mockService, userID)

	auditID := primitive.NewObjectID().Hex()

	// When
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit/"+auditID+"/appeal/review", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuditAdminAPI_ReviewAppeal_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouterWithAuth(mockService, "") // No user ID

	auditID := primitive.NewObjectID().Hex()
	reqBody := ReviewAppealRequest{
		Action:     "approve",
		ReviewNote: "申诉理由充分",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit/"+auditID+"/appeal/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "未授权", response["message"])
}

func TestAuditAdminAPI_ReviewAppeal_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	userID := "admin123"
	router := setupAuditAdminTestRouterWithAuth(mockService, userID)

	auditID := primitive.NewObjectID().Hex()
	reqBody := ReviewAppealRequest{
		Action:     "approve",
		ReviewNote: "申诉理由充分",
	}

	mockService.On("ReviewAppeal", mock.Anything, auditID, userID, true, "申诉理由充分").Return(assert.AnError)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit/"+auditID+"/appeal/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetHighRiskAudits Tests ====================

func TestAuditAdminAPI_GetHighRiskAudits_Success(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	expectedRecords := []*audit.AuditRecord{
		{ID: primitive.NewObjectID().Hex(), TargetType: "document", RiskLevel: 5},
		{ID: primitive.NewObjectID().Hex(), TargetType: "chapter", RiskLevel: 4},
	}

	mockService.On("GetHighRiskAudits", mock.Anything, 3, 50).Return(expectedRecords, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/high-risk", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "获取成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_GetHighRiskAudits_WithFilters(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	expectedRecords := []*audit.AuditRecord{
		{ID: primitive.NewObjectID().Hex(), TargetType: "document", RiskLevel: 5},
	}

	mockService.On("GetHighRiskAudits", mock.Anything, 4, 20).Return(expectedRecords, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/high-risk?minRiskLevel=4&limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_GetHighRiskAudits_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	mockService.On("GetHighRiskAudits", mock.Anything, 3, 50).Return(nil, assert.AnError)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/high-risk", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetAuditStatistics Tests ====================

func TestAuditAdminAPI_GetAuditStatistics_Success(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	expectedStats := map[string]interface{}{
		"totalAudits":      int64(1000),
		"pendingAudits":    int64(50),
		"approvedAudits":   int64(800),
		"rejectedAudits":   int64(150),
		"highRiskAudits":   int64(20),
		"averageProcessTime": float64(120.5),
	}

	mockService.On("GetAuditStatistics", mock.Anything).Return(expectedStats, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/statistics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "获取成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_GetAuditStatistics_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	mockService.On("GetAuditStatistics", mock.Anything).Return(nil, assert.AnError)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/statistics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}
