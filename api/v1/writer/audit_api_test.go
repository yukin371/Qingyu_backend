package writer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	auditModel "Qingyu_backend/models/audit"
	auditInterface "Qingyu_backend/service/interfaces/audit"
)

type mockWriterContentAuditService struct {
	mock.Mock
}

func (m *mockWriterContentAuditService) CheckContent(ctx context.Context, content string) (*auditInterface.AuditCheckResult, error) {
	args := m.Called(ctx, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditInterface.AuditCheckResult), args.Error(1)
}

func (m *mockWriterContentAuditService) AuditDocument(ctx context.Context, documentID string, content string, authorID string) (*auditModel.AuditRecord, error) {
	args := m.Called(ctx, documentID, content, authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditModel.AuditRecord), args.Error(1)
}

func (m *mockWriterContentAuditService) GetAuditResult(ctx context.Context, targetType, targetID string) (*auditModel.AuditRecord, error) {
	args := m.Called(ctx, targetType, targetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditModel.AuditRecord), args.Error(1)
}

func (m *mockWriterContentAuditService) BatchAuditDocuments(ctx context.Context, documentIDs []string) ([]*auditModel.AuditRecord, error) {
	args := m.Called(ctx, documentIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModel.AuditRecord), args.Error(1)
}

func (m *mockWriterContentAuditService) ReviewAudit(ctx context.Context, auditID string, reviewerID string, approved bool, note string) error {
	args := m.Called(ctx, auditID, reviewerID, approved, note)
	return args.Error(0)
}

func (m *mockWriterContentAuditService) SubmitAppeal(ctx context.Context, auditID string, authorID string, reason string) error {
	args := m.Called(ctx, auditID, authorID, reason)
	return args.Error(0)
}

func (m *mockWriterContentAuditService) ReviewAppeal(ctx context.Context, auditID string, reviewerID string, approved bool, note string) error {
	args := m.Called(ctx, auditID, reviewerID, approved, note)
	return args.Error(0)
}

func (m *mockWriterContentAuditService) GetUserViolations(ctx context.Context, userID string) ([]*auditModel.ViolationRecord, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModel.ViolationRecord), args.Error(1)
}

func (m *mockWriterContentAuditService) GetUserViolationSummary(ctx context.Context, userID string) (*auditModel.UserViolationSummary, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditModel.UserViolationSummary), args.Error(1)
}

func (m *mockWriterContentAuditService) GetPendingReviews(ctx context.Context, limit int) ([]*auditModel.AuditRecord, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModel.AuditRecord), args.Error(1)
}

func (m *mockWriterContentAuditService) GetHighRiskAudits(ctx context.Context, minRiskLevel int, limit int) ([]*auditModel.AuditRecord, error) {
	args := m.Called(ctx, minRiskLevel, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auditModel.AuditRecord), args.Error(1)
}

func (m *mockWriterContentAuditService) GetAuditStatistics(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0), args.Error(1)
}

func setupWriterAuditTestRouter(auditService *mockWriterContentAuditService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := NewAuditApi(auditService)

	router.GET("/pending", api.GetPendingReviews)
	router.GET("/high-risk", api.GetHighRiskAudits)

	return router
}

func TestAuditApiGetPendingReviewsParsesLimitQuery(t *testing.T) {
	service := new(mockWriterContentAuditService)
	service.On("GetPendingReviews", mock.Anything, 20).Return([]*auditModel.AuditRecord{}, nil)
	router := setupWriterAuditTestRouter(service)

	req, _ := http.NewRequest(http.MethodGet, "/pending?limit=20", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	service.AssertExpectations(t)
}

func TestAuditApiGetHighRiskAuditsParsesNumericQueries(t *testing.T) {
	service := new(mockWriterContentAuditService)
	service.On("GetHighRiskAudits", mock.Anything, 4, 20).Return([]*auditModel.AuditRecord{}, nil)
	router := setupWriterAuditTestRouter(service)

	req, _ := http.NewRequest(http.MethodGet, "/high-risk?minRiskLevel=4&limit=20", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	service.AssertExpectations(t)
}
