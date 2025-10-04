package container

import (
	"context"
	"io"
	"time"

	"Qingyu_backend/service/shared/admin"
	"Qingyu_backend/service/shared/auth"
	"Qingyu_backend/service/shared/messaging"
	"Qingyu_backend/service/shared/recommendation"
	"Qingyu_backend/service/shared/storage"
	"Qingyu_backend/service/shared/wallet"

	"github.com/stretchr/testify/mock"
)

// ============ Mock Auth Service ============

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.RegisterResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.LoginResponse), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, token string) error {
	return m.Called(ctx, token).Error(0)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenClaims), args.Error(1)
}

func (m *MockAuthService) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	args := m.Called(ctx, userID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthService) HasRole(ctx context.Context, userID, role string) (bool, error) {
	args := m.Called(ctx, userID, role)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthService) CreateRole(ctx context.Context, req *auth.CreateRoleRequest) (*auth.Role, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Role), args.Error(1)
}

func (m *MockAuthService) UpdateRole(ctx context.Context, roleID string, req *auth.UpdateRoleRequest) error {
	return m.Called(ctx, roleID, req).Error(0)
}

func (m *MockAuthService) DeleteRole(ctx context.Context, roleID string) error {
	return m.Called(ctx, roleID).Error(0)
}

func (m *MockAuthService) AssignRole(ctx context.Context, userID, roleID string) error {
	return m.Called(ctx, userID, roleID).Error(0)
}

func (m *MockAuthService) RemoveRole(ctx context.Context, userID, roleID string) error {
	return m.Called(ctx, userID, roleID).Error(0)
}

func (m *MockAuthService) CreateSession(ctx context.Context, userID string) (*auth.Session, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Session), args.Error(1)
}

func (m *MockAuthService) GetSession(ctx context.Context, sessionID string) (*auth.Session, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Session), args.Error(1)
}

func (m *MockAuthService) DestroySession(ctx context.Context, sessionID string) error {
	return m.Called(ctx, sessionID).Error(0)
}

func (m *MockAuthService) RefreshSession(ctx context.Context, sessionID string) error {
	return m.Called(ctx, sessionID).Error(0)
}

func (m *MockAuthService) Health(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// ============ Mock Wallet Service ============

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) CreateWallet(ctx context.Context, userID string) (*wallet.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wallet.Wallet), args.Error(1)
}

func (m *MockWalletService) GetWallet(ctx context.Context, userID string) (*wallet.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wallet.Wallet), args.Error(1)
}

func (m *MockWalletService) GetBalance(ctx context.Context, userID string) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockWalletService) FreezeWallet(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *MockWalletService) UnfreezeWallet(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *MockWalletService) Recharge(ctx context.Context, userID string, amount float64, method string) (*wallet.Transaction, error) {
	args := m.Called(ctx, userID, amount, method)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wallet.Transaction), args.Error(1)
}

func (m *MockWalletService) Consume(ctx context.Context, userID string, amount float64, reason string) (*wallet.Transaction, error) {
	args := m.Called(ctx, userID, amount, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wallet.Transaction), args.Error(1)
}

func (m *MockWalletService) Transfer(ctx context.Context, fromUserID, toUserID string, amount float64, reason string) (*wallet.Transaction, error) {
	args := m.Called(ctx, fromUserID, toUserID, amount, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wallet.Transaction), args.Error(1)
}

func (m *MockWalletService) GetTransaction(ctx context.Context, transactionID string) (*wallet.Transaction, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wallet.Transaction), args.Error(1)
}

func (m *MockWalletService) ListTransactions(ctx context.Context, userID string, req *wallet.ListTransactionsRequest) ([]*wallet.Transaction, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*wallet.Transaction), args.Error(1)
}

func (m *MockWalletService) RequestWithdraw(ctx context.Context, userID string, amount float64, account string) (*wallet.WithdrawRequest, error) {
	args := m.Called(ctx, userID, amount, account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wallet.WithdrawRequest), args.Error(1)
}

func (m *MockWalletService) GetWithdrawRequest(ctx context.Context, withdrawID string) (*wallet.WithdrawRequest, error) {
	args := m.Called(ctx, withdrawID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wallet.WithdrawRequest), args.Error(1)
}

func (m *MockWalletService) ListWithdrawRequests(ctx context.Context, req *wallet.ListWithdrawRequestsRequest) ([]*wallet.WithdrawRequest, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*wallet.WithdrawRequest), args.Error(1)
}

func (m *MockWalletService) ApproveWithdraw(ctx context.Context, withdrawID, adminID string) error {
	return m.Called(ctx, withdrawID, adminID).Error(0)
}

func (m *MockWalletService) RejectWithdraw(ctx context.Context, withdrawID, adminID, reason string) error {
	return m.Called(ctx, withdrawID, adminID, reason).Error(0)
}

func (m *MockWalletService) ProcessWithdraw(ctx context.Context, withdrawID string) error {
	return m.Called(ctx, withdrawID).Error(0)
}

func (m *MockWalletService) Health(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// ============ Mock Recommendation Service ============

type MockRecommendationService struct {
	mock.Mock
}

func (m *MockRecommendationService) GetPersonalizedRecommendations(ctx context.Context, userID string, limit int) ([]*recommendation.RecommendedItem, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*recommendation.RecommendedItem), args.Error(1)
}

func (m *MockRecommendationService) GetSimilarItems(ctx context.Context, itemID string, limit int) ([]*recommendation.RecommendedItem, error) {
	args := m.Called(ctx, itemID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*recommendation.RecommendedItem), args.Error(1)
}

func (m *MockRecommendationService) GetHotItems(ctx context.Context, itemType string, limit int) ([]*recommendation.RecommendedItem, error) {
	args := m.Called(ctx, itemType, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*recommendation.RecommendedItem), args.Error(1)
}

func (m *MockRecommendationService) RecordUserBehavior(ctx context.Context, req *recommendation.RecordBehaviorRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *MockRecommendationService) GetUserBehaviors(ctx context.Context, userID string, limit int) ([]*recommendation.UserBehavior, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*recommendation.UserBehavior), args.Error(1)
}

func (m *MockRecommendationService) RefreshRecommendations(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *MockRecommendationService) RefreshHotItems(ctx context.Context, itemType string) error {
	return m.Called(ctx, itemType).Error(0)
}

func (m *MockRecommendationService) Health(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// ============ Mock Messaging Service ============

type MockMessagingService struct {
	mock.Mock
}

func (m *MockMessagingService) Publish(ctx context.Context, topic string, message interface{}) error {
	return m.Called(ctx, topic, message).Error(0)
}

func (m *MockMessagingService) PublishDelayed(ctx context.Context, topic string, message interface{}, delay time.Duration) error {
	return m.Called(ctx, topic, message, delay).Error(0)
}

func (m *MockMessagingService) Subscribe(ctx context.Context, topic string, handler messaging.MessageHandler) error {
	return m.Called(ctx, topic, handler).Error(0)
}

func (m *MockMessagingService) Unsubscribe(ctx context.Context, topic string) error {
	return m.Called(ctx, topic).Error(0)
}

func (m *MockMessagingService) CreateTopic(ctx context.Context, topic string) error {
	return m.Called(ctx, topic).Error(0)
}

func (m *MockMessagingService) DeleteTopic(ctx context.Context, topic string) error {
	return m.Called(ctx, topic).Error(0)
}

func (m *MockMessagingService) ListTopics(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockMessagingService) Health(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// ============ Mock Storage Service ============

type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) Upload(ctx context.Context, req *storage.UploadRequest) (*storage.FileInfo, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*storage.FileInfo), args.Error(1)
}

func (m *MockStorageService) Download(ctx context.Context, fileID string) (io.ReadCloser, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorageService) Delete(ctx context.Context, fileID string) error {
	return m.Called(ctx, fileID).Error(0)
}

func (m *MockStorageService) GetFileInfo(ctx context.Context, fileID string) (*storage.FileInfo, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*storage.FileInfo), args.Error(1)
}

func (m *MockStorageService) GrantAccess(ctx context.Context, fileID, userID string) error {
	return m.Called(ctx, fileID, userID).Error(0)
}

func (m *MockStorageService) RevokeAccess(ctx context.Context, fileID, userID string) error {
	return m.Called(ctx, fileID, userID).Error(0)
}

func (m *MockStorageService) CheckAccess(ctx context.Context, fileID, userID string) (bool, error) {
	args := m.Called(ctx, fileID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorageService) ListFiles(ctx context.Context, req *storage.ListFilesRequest) ([]*storage.FileInfo, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*storage.FileInfo), args.Error(1)
}

func (m *MockStorageService) GetDownloadURL(ctx context.Context, fileID string, expiresIn time.Duration) (string, error) {
	args := m.Called(ctx, fileID, expiresIn)
	return args.String(0), args.Error(1)
}

func (m *MockStorageService) Health(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// ============ Mock Admin Service ============

type MockAdminService struct {
	mock.Mock
}

func (m *MockAdminService) ReviewContent(ctx context.Context, req *admin.ReviewContentRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *MockAdminService) GetPendingReviews(ctx context.Context, contentType string) ([]*admin.AuditRecord, error) {
	args := m.Called(ctx, contentType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*admin.AuditRecord), args.Error(1)
}

func (m *MockAdminService) ManageUser(ctx context.Context, req *admin.ManageUserRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *MockAdminService) BanUser(ctx context.Context, userID, reason string, duration time.Duration) error {
	return m.Called(ctx, userID, reason, duration).Error(0)
}

func (m *MockAdminService) UnbanUser(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *MockAdminService) GetUserStatistics(ctx context.Context, userID string) (*admin.UserStatistics, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*admin.UserStatistics), args.Error(1)
}

func (m *MockAdminService) ReviewWithdraw(ctx context.Context, withdrawID, adminID string, approved bool, reason string) error {
	return m.Called(ctx, withdrawID, adminID, approved, reason).Error(0)
}

func (m *MockAdminService) LogOperation(ctx context.Context, req *admin.LogOperationRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *MockAdminService) GetOperationLogs(ctx context.Context, req *admin.GetLogsRequest) ([]*admin.AdminLog, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*admin.AdminLog), args.Error(1)
}

func (m *MockAdminService) ExportLogs(ctx context.Context, startDate, endDate time.Time) ([]byte, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockAdminService) Health(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
