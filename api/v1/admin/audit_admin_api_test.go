package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	apperrors "Qingyu_backend/pkg/errors"
	"Qingyu_backend/internal/middleware/builtin"
	"Qingyu_backend/pkg/logger"
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

	// 添加错误处理中间件
	errorHandler := builtin.NewErrorHandlerMiddleware(logger.Get().Logger)
	r.Use(errorHandler.Handler())

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

	// 添加错误处理中间件
	errorHandler := builtin.NewErrorHandlerMiddleware(logger.Get().Logger)
	r.Use(errorHandler.Handler())

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
		v1.POST("/batch-review", api.BatchReviewAudit)
	}

	return r
}

// ==================== GetPendingAudits Tests ====================

func TestAuditAdminAPI_GetPendingAudits_Success(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	expectedRecords := []*audit.AuditRecord{
		func() *audit.AuditRecord {
			r := &audit.AuditRecord{TargetType: "document", Status: "pending"}
			r.ID = primitive.NewObjectID()
			return r
		}(),
		func() *audit.AuditRecord {
			r := &audit.AuditRecord{TargetType: "chapter", Status: "pending"}
			r.ID = primitive.NewObjectID()
			return r
		}(),
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
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "操作成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_GetPendingAudits_WithTargetType(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	expectedRecords := []*audit.AuditRecord{
		func() *audit.AuditRecord {
			r := &audit.AuditRecord{TargetType: "document", Status: "pending"}
			r.ID = primitive.NewObjectID()
			return r
		}(),
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

	mockService.On("GetPendingReviews", mock.Anything, 50).Return(
		nil,
		apperrors.BookstoreServiceFactory.InternalError("GET_PENDING_AUDITS_FAILED", "获取待审核列表失败", errors.New("database error")),
	)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/pending", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 中间件会自动映射InternalError为500状态码
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
	assert.Equal(t, float64(0), response["code"])
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
	assert.Equal(t, float64(0), response["code"])
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

	mockService.On("ReviewAudit", mock.Anything, auditID, userID, true, "内容符合规范").Return(
		apperrors.BookstoreServiceFactory.InternalError("REVIEW_AUDIT_FAILED", "审核失败", errors.New("service error")),
	)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit/"+auditID+"/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 中间件会自动映射InternalError为500状态码
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
	assert.Equal(t, float64(0), response["code"])
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
	assert.Equal(t, float64(0), response["code"])
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

	mockService.On("ReviewAppeal", mock.Anything, auditID, userID, true, "申诉理由充分").Return(
		apperrors.BookstoreServiceFactory.InternalError("REVIEW_APPEAL_FAILED", "审核申诉失败", errors.New("service error")),
	)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/audit/"+auditID+"/appeal/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 中间件会自动映射InternalError为500状态码
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetHighRiskAudits Tests ====================

func TestAuditAdminAPI_GetHighRiskAudits_Success(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	expectedRecords := []*audit.AuditRecord{
		func() *audit.AuditRecord {
			r := &audit.AuditRecord{TargetType: "document", RiskLevel: 5}
			r.ID = primitive.NewObjectID()
			return r
		}(),
		func() *audit.AuditRecord {
			r := &audit.AuditRecord{TargetType: "chapter", RiskLevel: 4}
			r.ID = primitive.NewObjectID()
			return r
		}(),
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
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "操作成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_GetHighRiskAudits_WithFilters(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	expectedRecords := []*audit.AuditRecord{
		func() *audit.AuditRecord {
			r := &audit.AuditRecord{TargetType: "document", RiskLevel: 5}
			r.ID = primitive.NewObjectID()
			return r
		}(),
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

	mockService.On("GetHighRiskAudits", mock.Anything, 3, 50).Return(
		nil,
		apperrors.BookstoreServiceFactory.InternalError("GET_HIGH_RISK_AUDITS_FAILED", "获取高风险审核记录失败", errors.New("database error")),
	)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/high-risk", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 中间件会自动映射InternalError为500状态码
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetAuditStatistics Tests ====================

func TestAuditAdminAPI_GetAuditStatistics_Success(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	expectedStats := map[string]interface{}{
		"totalAudits":        int64(1000),
		"pendingAudits":      int64(50),
		"approvedAudits":     int64(800),
		"rejectedAudits":     int64(150),
		"highRiskAudits":     int64(20),
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
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "操作成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_GetAuditStatistics_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouter(mockService)

	mockService.On("GetAuditStatistics", mock.Anything).Return(
		nil,
		apperrors.BookstoreServiceFactory.InternalError("GET_AUDIT_STATISTICS_FAILED", "获取审核统计失败", errors.New("database error")),
	)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/statistics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 中间件会自动映射InternalError为500状态码
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== Batch Review Audit Tests ====================

func TestAuditAdminAPI_BatchReviewAudit_Success(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouterWithAuth(mockService, "admin123")

	auditID1 := primitive.NewObjectID().Hex()
	auditID2 := primitive.NewObjectID().Hex()

	reqBody := BatchReviewAuditRequest{
		AuditIDs:   []string{auditID1, auditID2},
		Action:     "approve",
		ReviewNote: "批量审核通过",
	}

	mockService.On("ReviewAudit", mock.Anything, auditID1, "admin123", true, "批量审核通过").Return(nil)
	mockService.On("ReviewAudit", mock.Anything, auditID2, "admin123", true, "批量审核通过").Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	reqHTTP, _ := http.NewRequest("POST", "/api/v1/admin/audit/batch-review", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.Equal(t, float64(2), data["success"])
	assert.Equal(t, float64(0), data["failed"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_BatchReviewAudit_ApproveAll(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouterWithAuth(mockService, "admin123")

	auditIDs := []string{
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
	}

	reqBody := BatchReviewAuditRequest{
		AuditIDs:   auditIDs,
		Action:     "approve",
		ReviewNote: "全部通过",
	}

	for _, id := range auditIDs {
		mockService.On("ReviewAudit", mock.Anything, id, "admin123", true, "全部通过").Return(nil)
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	reqHTTP, _ := http.NewRequest("POST", "/api/v1/admin/audit/batch-review", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(3), data["total"])
	assert.Equal(t, float64(3), data["success"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_BatchReviewAudit_RejectAll(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouterWithAuth(mockService, "admin123")

	auditIDs := []string{
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
	}

	reqBody := BatchReviewAuditRequest{
		AuditIDs:   auditIDs,
		Action:     "reject",
		ReviewNote: "全部拒绝",
	}

	for _, id := range auditIDs {
		mockService.On("ReviewAudit", mock.Anything, id, "admin123", false, "全部拒绝").Return(nil)
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	reqHTTP, _ := http.NewRequest("POST", "/api/v1/admin/audit/batch-review", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.Equal(t, float64(2), data["success"])

	mockService.AssertExpectations(t)
}

func TestAuditAdminAPI_BatchReviewAudit_EmptyList_Error(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouterWithAuth(mockService, "admin123")

	reqBody := BatchReviewAuditRequest{
		AuditIDs:   []string{},
		Action:     "approve",
		ReviewNote: "测试",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	reqHTTP, _ := http.NewRequest("POST", "/api/v1/admin/audit/batch-review", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuditAdminAPI_BatchReviewAudit_InvalidAction_Error(t *testing.T) {
	// Given
	mockService := new(MockContentAuditService)
	router := setupAuditAdminTestRouterWithAuth(mockService, "admin123")

	reqBody := BatchReviewAuditRequest{
		AuditIDs:   []string{primitive.NewObjectID().Hex()},
		Action:     "invalid_action",
		ReviewNote: "测试",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	reqHTTP, _ := http.NewRequest("POST", "/api/v1/admin/audit/batch-review", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
