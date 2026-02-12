package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	"Qingyu_backend/service/admin"
	"Qingyu_backend/service/finance/wallet"
	"Qingyu_backend/service/auth"
	"Qingyu_backend/service/shared/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============ Mock服务（使用与container测试相同的Mock） ============

// MockAuthService Auth服务Mock
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

// MockWalletService Wallet服务Mock
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

// MockStorageService Storage服务Mock（简化版）
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

func (m *MockStorageService) CheckAccess(ctx context.Context, fileID, userID string) (bool, error) {
	args := m.Called(ctx, fileID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorageService) GrantAccess(ctx context.Context, fileID, userID string) error {
	return m.Called(ctx, fileID, userID).Error(0)
}

func (m *MockStorageService) GetFileInfo(ctx context.Context, fileID string) (*storage.FileInfo, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*storage.FileInfo), args.Error(1)
}

// MockAdminService Admin服务Mock（简化版）
type MockAdminService struct {
	mock.Mock
}

func (m *MockAdminService) ReviewWithdraw(ctx context.Context, withdrawID, adminID string, approved bool, reason string) error {
	return m.Called(ctx, withdrawID, adminID, approved, reason).Error(0)
}

func (m *MockAdminService) GetUserStatistics(ctx context.Context, userID string) (*admin.UserStatistics, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*admin.UserStatistics), args.Error(1)
}

// ============ 测试场景1: Auth + Wallet 集成 ============

// TestAuthWallet_UserRegisterAndCreateWallet 测试用户注册后创建钱包
func TestAuthWallet_UserRegisterAndCreateWallet(t *testing.T) {
	ctx := context.Background()

	// 创建Mock服务
	authService := new(MockAuthService)
	walletService := new(MockWalletService)

	// 模拟用户注册
	registerReq := &auth.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "reader",
	}

	registerResp := &auth.RegisterResponse{
		User: &auth.UserInfo{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Roles:    []string{"reader"},
		},
		Token: "jwt_token_12345",
	}

	authService.On("Register", ctx, registerReq).Return(registerResp, nil)

	// 模拟为新用户创建钱包
	expectedWallet := &wallet.Wallet{
		ID:        "wallet123",
		UserID:    "user123",
		Balance:   0.0,
		Frozen:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	walletService.On("CreateWallet", ctx, "user123").Return(expectedWallet, nil)

	// 执行业务流程：注册 -> 创建钱包
	// 1. 用户注册
	user, err := authService.Register(ctx, registerReq)
	assert.NoError(t, err, "注册应该成功")
	assert.NotNil(t, user, "返回的用户信息不应为空")
	assert.Equal(t, "user123", user.User.ID)

	// 2. 为新用户创建钱包
	userWallet, err := walletService.CreateWallet(ctx, user.User.ID)
	assert.NoError(t, err, "创建钱包应该成功")
	assert.NotNil(t, userWallet, "钱包不应为空")
	assert.Equal(t, "user123", userWallet.UserID)
	assert.Equal(t, 0.0, userWallet.Balance, "新钱包余额应为0")
	assert.False(t, userWallet.Frozen, "新钱包不应被冻结")

	// 验证Mock调用
	authService.AssertExpectations(t)
	walletService.AssertExpectations(t)
}

// TestAuthWallet_UserLoginAndRecharge 测试用户登录后充值
func TestAuthWallet_UserLoginAndRecharge(t *testing.T) {
	ctx := context.Background()

	// 创建Mock服务
	authService := new(MockAuthService)
	walletService := new(MockWalletService)

	// 模拟用户登录
	loginReq := &auth.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	loginResp := &auth.LoginResponse{
		User: &auth.UserInfo{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Roles:    []string{"reader"},
		},
		Token: "jwt_token_67890",
	}

	authService.On("Login", ctx, loginReq).Return(loginResp, nil)

	// 模拟充值
	expectedTransaction := &wallet.Transaction{
		ID:              "txn123",
		UserID:          "user123",
		Type:            "recharge",
		Amount:          100.0,
		Balance:         100.0,
		Method:          "alipay",
		Status:          "success",
		TransactionTime: time.Now(),
		CreatedAt:       time.Now(),
	}

	walletService.On("Recharge", ctx, "user123", 100.0, "alipay").Return(expectedTransaction, nil)

	// 执行业务流程：登录 -> 充值
	// 1. 用户登录
	user, err := authService.Login(ctx, loginReq)
	assert.NoError(t, err, "登录应该成功")
	assert.NotNil(t, user, "返回的用户信息不应为空")
	assert.Equal(t, "user123", user.User.ID)

	// 2. 用户充值
	transaction, err := walletService.Recharge(ctx, user.User.ID, 100.0, "alipay")
	assert.NoError(t, err, "充值应该成功")
	assert.NotNil(t, transaction, "交易记录不应为空")
	assert.Equal(t, "recharge", transaction.Type)
	assert.Equal(t, 100.0, transaction.Amount)
	assert.Equal(t, "success", transaction.Status)

	// 验证Mock调用
	authService.AssertExpectations(t)
	walletService.AssertExpectations(t)
}

// TestAuthWallet_UserConsumeAfterAuth 测试用户认证后消费
func TestAuthWallet_UserConsumeAfterAuth(t *testing.T) {
	ctx := context.Background()

	// 创建Mock服务
	authService := new(MockAuthService)
	walletService := new(MockWalletService)

	// 模拟Token验证
	token := "jwt_token_abc123"
	tokenClaims := &auth.TokenClaims{
		UserID: "user123",
		Roles:  []string{"reader"},
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	authService.On("ValidateToken", ctx, token).Return(tokenClaims, nil)

	// 模拟消费
	expectedTransaction := &wallet.Transaction{
		ID:              "txn456",
		UserID:          "user123",
		Type:            "consume",
		Amount:          -50.0,
		Balance:         50.0,
		Reason:          "购买书籍",
		Status:          "success",
		TransactionTime: time.Now(),
		CreatedAt:       time.Now(),
	}

	walletService.On("Consume", ctx, "user123", 50.0, "购买书籍").Return(expectedTransaction, nil)

	// 执行业务流程：验证Token -> 消费
	// 1. 验证用户Token
	claims, err := authService.ValidateToken(ctx, token)
	assert.NoError(t, err, "Token验证应该成功")
	assert.NotNil(t, claims, "Token声明不应为空")
	assert.Equal(t, "user123", claims.UserID)

	// 2. 用户消费
	transaction, err := walletService.Consume(ctx, claims.UserID, 50.0, "购买书籍")
	assert.NoError(t, err, "消费应该成功")
	assert.NotNil(t, transaction, "交易记录不应为空")
	assert.Equal(t, "consume", transaction.Type)
	assert.Equal(t, -50.0, transaction.Amount)
	assert.Equal(t, 50.0, transaction.Balance, "消费后余额应为50")

	// 验证Mock调用
	authService.AssertExpectations(t)
	walletService.AssertExpectations(t)
}

// ============ 测试场景2: Admin + Wallet 集成 ============

// TestAdminWallet_ApproveWithdrawRequest 测试管理员审核提现请求（通过）
func TestAdminWallet_ApproveWithdrawRequest(t *testing.T) {
	ctx := context.Background()

	// 创建Mock服务
	adminService := new(MockAdminService)
	walletService := new(MockWalletService)

	withdrawID := "withdraw123"
	adminID := "admin001"

	// 模拟获取提现请求
	withdrawRequest := &wallet.WithdrawRequest{
		ID:        withdrawID,
		UserID:    "user123",
		Amount:    100.0,
		Account:   "bank_account_123",
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	walletService.On("GetWithdrawRequest", ctx, withdrawID).Return(withdrawRequest, nil)

	// 模拟管理员审核通过
	adminService.On("ReviewWithdraw", ctx, withdrawID, adminID, true, "").Return(nil)

	// 模拟钱包服务批准提现
	walletService.On("ApproveWithdraw", ctx, withdrawID, adminID).Return(nil)

	// 执行业务流程：获取请求 -> 管理员审核 -> 批准提现
	// 1. 获取提现请求
	request, err := walletService.GetWithdrawRequest(ctx, withdrawID)
	assert.NoError(t, err, "获取提现请求应该成功")
	assert.Equal(t, "pending", request.Status)
	assert.Equal(t, 100.0, request.Amount)

	// 2. 管理员审核
	err = adminService.ReviewWithdraw(ctx, withdrawID, adminID, true, "")
	assert.NoError(t, err, "管理员审核应该成功")

	// 3. 批准提现
	err = walletService.ApproveWithdraw(ctx, withdrawID, adminID)
	assert.NoError(t, err, "批准提现应该成功")

	// 验证Mock调用
	adminService.AssertExpectations(t)
	walletService.AssertExpectations(t)
}

// TestAdminWallet_RejectWithdrawRequest 测试管理员拒绝提现请求
func TestAdminWallet_RejectWithdrawRequest(t *testing.T) {
	ctx := context.Background()

	// 创建Mock服务
	adminService := new(MockAdminService)
	walletService := new(MockWalletService)

	withdrawID := "withdraw456"
	adminID := "admin001"
	rejectReason := "账户信息不符"

	// 模拟获取提现请求
	withdrawRequest := &wallet.WithdrawRequest{
		ID:        withdrawID,
		UserID:    "user456",
		Amount:    200.0,
		Account:   "invalid_account",
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	walletService.On("GetWithdrawRequest", ctx, withdrawID).Return(withdrawRequest, nil)

	// 模拟管理员审核拒绝
	adminService.On("ReviewWithdraw", ctx, withdrawID, adminID, false, rejectReason).Return(nil)

	// 模拟钱包服务拒绝提现（退还金额）
	walletService.On("RejectWithdraw", ctx, withdrawID, adminID, rejectReason).Return(nil)

	// 执行业务流程：获取请求 -> 管理员拒绝 -> 退还金额
	// 1. 获取提现请求
	request, err := walletService.GetWithdrawRequest(ctx, withdrawID)
	assert.NoError(t, err, "获取提现请求应该成功")
	assert.Equal(t, "pending", request.Status)

	// 2. 管理员拒绝
	err = adminService.ReviewWithdraw(ctx, withdrawID, adminID, false, rejectReason)
	assert.NoError(t, err, "管理员审核应该成功")

	// 3. 拒绝提现并退还金额
	err = walletService.RejectWithdraw(ctx, withdrawID, adminID, rejectReason)
	assert.NoError(t, err, "拒绝提现应该成功")

	// 验证Mock调用
	adminService.AssertExpectations(t)
	walletService.AssertExpectations(t)
}

// TestAdminWallet_ReviewWithdrawWithWalletCheck 测试审核提现时检查钱包状态
func TestAdminWallet_ReviewWithdrawWithWalletCheck(t *testing.T) {
	ctx := context.Background()

	// 创建Mock服务
	adminService := new(MockAdminService)
	walletService := new(MockWalletService)

	withdrawID := "withdraw789"
	adminID := "admin001"
	userID := "user789"

	// 模拟获取提现请求
	withdrawRequest := &wallet.WithdrawRequest{
		ID:        withdrawID,
		UserID:    userID,
		Amount:    500.0,
		Account:   "bank_account_789",
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	walletService.On("GetWithdrawRequest", ctx, withdrawID).Return(withdrawRequest, nil)

	// 模拟获取钱包信息（检查是否冻结）
	userWallet := &wallet.Wallet{
		ID:        "wallet789",
		UserID:    userID,
		Balance:   1000.0,
		Frozen:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	walletService.On("GetWallet", ctx, userID).Return(userWallet, nil)

	// 钱包正常，可以批准
	adminService.On("ReviewWithdraw", ctx, withdrawID, adminID, true, "").Return(nil)
	walletService.On("ApproveWithdraw", ctx, withdrawID, adminID).Return(nil)

	// 执行业务流程：获取请求 -> 检查钱包 -> 审核通过
	// 1. 获取提现请求
	request, err := walletService.GetWithdrawRequest(ctx, withdrawID)
	assert.NoError(t, err)
	assert.Equal(t, 500.0, request.Amount)

	// 2. 检查用户钱包状态
	userWalletInfo, err := walletService.GetWallet(ctx, request.UserID)
	assert.NoError(t, err)
	assert.False(t, userWalletInfo.Frozen, "钱包不应被冻结")
	assert.GreaterOrEqual(t, userWalletInfo.Balance, request.Amount, "余额应足够")

	// 3. 审核通过
	err = adminService.ReviewWithdraw(ctx, withdrawID, adminID, true, "")
	assert.NoError(t, err)

	// 4. 批准提现
	err = walletService.ApproveWithdraw(ctx, withdrawID, adminID)
	assert.NoError(t, err)

	// 验证Mock调用
	adminService.AssertExpectations(t)
	walletService.AssertExpectations(t)
}

// ============ 测试场景3: Storage + Auth 集成 ============

// TestStorageAuth_UploadFileWithAuth 测试用户上传文件（需要认证）
func TestStorageAuth_UploadFileWithAuth(t *testing.T) {
	ctx := context.Background()

	// 创建Mock服务
	authService := new(MockAuthService)
	storageService := new(MockStorageService)

	// 模拟Token验证
	token := "jwt_token_file_upload"
	tokenClaims := &auth.TokenClaims{
		UserID: "user999",
		Roles:  []string{"writer"},
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	authService.On("ValidateToken", ctx, token).Return(tokenClaims, nil)

	// 模拟文件上传
	uploadReq := &storage.UploadRequest{
		Filename:    "test_document.pdf",
		ContentType: "application/pdf",
		Size:        1024000,
		UserID:      "user999",
		IsPublic:    false,
		Category:    "document",
	}

	uploadedFile := &storage.FileInfo{
		ID:           "file123",
		Filename:     "test_document_uuid.pdf",
		OriginalName: "test_document.pdf",
		ContentType:  "application/pdf",
		Size:         1024000,
		Path:         "/uploads/documents/test_document_uuid.pdf",
		UserID:       "user999",
		IsPublic:     false,
		Category:     "document",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	storageService.On("Upload", ctx, uploadReq).Return(uploadedFile, nil)

	// 执行业务流程：验证Token -> 上传文件
	// 1. 验证用户Token
	claims, err := authService.ValidateToken(ctx, token)
	assert.NoError(t, err, "Token验证应该成功")
	assert.Equal(t, "user999", claims.UserID)

	// 2. 上传文件
	uploadReq.UserID = claims.UserID
	fileInfo, err := storageService.Upload(ctx, uploadReq)
	assert.NoError(t, err, "文件上传应该成功")
	assert.NotNil(t, fileInfo)
	assert.Equal(t, "user999", fileInfo.UserID)
	assert.False(t, fileInfo.IsPublic, "文件应为私有")

	// 验证Mock调用
	authService.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

// TestStorageAuth_CheckFileAccess 测试检查文件访问权限
func TestStorageAuth_CheckFileAccess(t *testing.T) {
	ctx := context.Background()

	// 创建Mock服务
	storageService := new(MockStorageService)

	fileID := "file456"
	fileOwnerID := "user888"
	accessUserID := "user999"

	// 模拟获取文件信息
	fileInfo := &storage.FileInfo{
		ID:           fileID,
		Filename:     "private_document.pdf",
		OriginalName: "private_document.pdf",
		UserID:       fileOwnerID,
		IsPublic:     false,
		Category:     "document",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	storageService.On("GetFileInfo", ctx, fileID).Return(fileInfo, nil)

	// 模拟检查访问权限（用户无权限）
	storageService.On("CheckAccess", ctx, fileID, accessUserID).Return(false, nil)

	// 执行业务流程：获取文件信息 -> 检查权限
	// 1. 获取文件信息
	file, err := storageService.GetFileInfo(ctx, fileID)
	assert.NoError(t, err)
	assert.False(t, file.IsPublic, "文件是私有的")

	// 2. 检查访问权限
	hasAccess, err := storageService.CheckAccess(ctx, fileID, accessUserID)
	assert.NoError(t, err)
	assert.False(t, hasAccess, "用户应该无权访问")

	// 验证Mock调用
	storageService.AssertExpectations(t)
}

// TestStorageAuth_GrantAccessAfterPermissionCheck 测试管理员授权文件访问
func TestStorageAuth_GrantAccessAfterPermissionCheck(t *testing.T) {
	ctx := context.Background()

	// 创建Mock服务
	authService := new(MockAuthService)
	storageService := new(MockStorageService)

	fileID := "file789"
	adminID := "admin001"
	targetUserID := "user777"

	// 模拟验证管理员权限
	authService.On("CheckPermission", ctx, adminID, "file:grant_access").Return(true, nil)

	// 模拟授权文件访问
	storageService.On("GrantAccess", ctx, fileID, targetUserID).Return(nil)

	// 验证授权后用户可以访问
	storageService.On("CheckAccess", ctx, fileID, targetUserID).Return(true, nil)

	// 执行业务流程：检查管理员权限 -> 授权访问 -> 验证访问
	// 1. 检查管理员权限
	hasPermission, err := authService.CheckPermission(ctx, adminID, "file:grant_access")
	assert.NoError(t, err)
	assert.True(t, hasPermission, "管理员应有授权权限")

	// 2. 授权文件访问
	err = storageService.GrantAccess(ctx, fileID, targetUserID)
	assert.NoError(t, err, "授权应该成功")

	// 3. 验证用户现在可以访问
	canAccess, err := storageService.CheckAccess(ctx, fileID, targetUserID)
	assert.NoError(t, err)
	assert.True(t, canAccess, "用户现在应该可以访问")

	// 验证Mock调用
	authService.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

// ============ 错误处理测试 ============

// TestAuthWallet_RegisterFailure 测试注册失败时不创建钱包
func TestAuthWallet_RegisterFailure(t *testing.T) {
	ctx := context.Background()

	// 创建Mock服务
	authService := new(MockAuthService)
	walletService := new(MockWalletService)

	// 模拟用户注册失败
	registerReq := &auth.RegisterRequest{
		Username: "duplicate_user",
		Email:    "duplicate@example.com",
		Password: "password123",
	}

	authService.On("Register", ctx, registerReq).Return(nil, errors.New("用户名已存在"))

	// 执行业务流程：注册失败 -> 不创建钱包
	user, err := authService.Register(ctx, registerReq)
	assert.Error(t, err, "注册应该失败")
	assert.Nil(t, user, "用户信息应为空")
	assert.Contains(t, err.Error(), "用户名已存在")

	// 不应该调用创建钱包（因为注册失败了）
	// walletService 不应该有任何调用

	// 验证Mock调用
	authService.AssertExpectations(t)
	walletService.AssertNotCalled(t, "CreateWallet", mock.Anything, mock.Anything)
}

// TestAuthWallet_InsufficientBalance 测试余额不足时消费失败
func TestAuthWallet_InsufficientBalance(t *testing.T) {
	ctx := context.Background()

	// 创建Mock服务
	authService := new(MockAuthService)
	walletService := new(MockWalletService)

	userID := "user_poor"
	token := "jwt_token_poor_user"

	// 验证Token
	tokenClaims := &auth.TokenClaims{
		UserID: userID,
		Roles:  []string{"reader"},
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	authService.On("ValidateToken", ctx, token).Return(tokenClaims, nil)

	// 模拟消费失败（余额不足）
	walletService.On("Consume", ctx, userID, 1000.0, "购买高级会员").
		Return(nil, errors.New("余额不足"))

	// 执行业务流程：验证Token -> 消费失败
	claims, err := authService.ValidateToken(ctx, token)
	assert.NoError(t, err)

	transaction, err := walletService.Consume(ctx, claims.UserID, 1000.0, "购买高级会员")
	assert.Error(t, err, "消费应该失败")
	assert.Nil(t, transaction, "交易记录应为空")
	assert.Contains(t, err.Error(), "余额不足")

	// 验证Mock调用
	authService.AssertExpectations(t)
	walletService.AssertExpectations(t)
}

// TestStorageAuth_UnauthorizedUpload 测试未授权用户无法上传文件
func TestStorageAuth_UnauthorizedUpload(t *testing.T) {
	ctx := context.Background()

	// 创建Mock服务
	authService := new(MockAuthService)
	storageService := new(MockStorageService)

	// 模拟Token验证失败
	invalidToken := "invalid_token"
	authService.On("ValidateToken", ctx, invalidToken).
		Return(nil, errors.New("无效的token"))

	// 执行业务流程：Token验证失败 -> 不上传文件
	claims, err := authService.ValidateToken(ctx, invalidToken)
	assert.Error(t, err, "Token验证应该失败")
	assert.Nil(t, claims, "Token声明应为空")

	// 不应该调用上传（因为认证失败）
	storageService.AssertNotCalled(t, "Upload", mock.Anything, mock.Anything)

	// 验证Mock调用
	authService.AssertExpectations(t)
}
