package api

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
	"github.com/stretchr/testify/require"

	sharedAPI "Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/shared/wallet"
)

// === Mock WalletService ===

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
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockWalletService) UnfreezeWallet(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
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
	args := m.Called(ctx, withdrawID, adminID)
	return args.Error(0)
}

func (m *MockWalletService) RejectWithdraw(ctx context.Context, withdrawID, adminID, reason string) error {
	args := m.Called(ctx, withdrawID, adminID, reason)
	return args.Error(0)
}

func (m *MockWalletService) ProcessWithdraw(ctx context.Context, withdrawID string) error {
	args := m.Called(ctx, withdrawID)
	return args.Error(0)
}

func (m *MockWalletService) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// === 测试辅助函数 ===

func setupWalletTestRouter(walletService wallet.WalletService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置user_id
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := sharedAPI.NewWalletAPI(walletService)

	v1 := r.Group("/api/v1/shared/wallet")
	{
		v1.GET("/balance", api.GetBalance)
		v1.GET("", api.GetWallet)
		v1.POST("/recharge", api.Recharge)
		v1.POST("/consume", api.Consume)
		v1.POST("/transfer", api.Transfer)
		v1.GET("/transactions", api.GetTransactions)
		v1.POST("/withdraw", api.RequestWithdraw)
		v1.GET("/withdrawals", api.GetWithdrawRequests)
	}

	return r
}

func createTestWallet(userID string, balance float64) *wallet.Wallet {
	return &wallet.Wallet{
		UserID:  userID,
		Balance: balance,
		Frozen:  false,
	}
}

func createTestTransaction(userID string, amount float64, txType string) *wallet.Transaction {
	return &wallet.Transaction{
		ID:     "tx123",
		UserID: userID,
		Amount: amount,
		Type:   txType,
		Status: "completed",
	}
}

// === 测试用例 ===

func TestWalletAPI_GetBalance(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockWalletService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "成功获取余额",
			userID: "user123",
			setupMock: func(service *MockWalletService) {
				service.On("GetBalance", mock.Anything, "user123").Return(100.50, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "查询余额成功", resp["message"])
				assert.Equal(t, 100.50, resp["data"])
			},
		},
		{
			name:   "未认证用户",
			userID: "",
			setupMock: func(service *MockWalletService) {
				// 不设置Mock，因为不会调用Service
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(401), resp["code"])
				assert.Equal(t, "未认证", resp["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWalletService)
			tt.setupMock(mockService)

			router := setupWalletTestRouter(mockService, tt.userID)

			req := httptest.NewRequest("GET", "/api/v1/shared/wallet/balance", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}

func TestWalletAPI_GetWallet(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockWalletService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "成功获取钱包信息",
			userID: "user123",
			setupMock: func(service *MockWalletService) {
				testWallet := createTestWallet("user123", 100.50)
				service.On("GetWallet", mock.Anything, "user123").Return(testWallet, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "获取钱包信息成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "user123", data["user_id"])
				assert.Equal(t, 100.50, data["balance"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWalletService)
			tt.setupMock(mockService)

			router := setupWalletTestRouter(mockService, tt.userID)

			req := httptest.NewRequest("GET", "/api/v1/shared/wallet", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}

func TestWalletAPI_Recharge(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    sharedAPI.RechargeRequest
		setupMock      func(*MockWalletService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "成功充值",
			userID: "user123",
			requestBody: sharedAPI.RechargeRequest{
				Amount: 100.0,
				Method: "alipay",
			},
			setupMock: func(service *MockWalletService) {
				tx := createTestTransaction("user123", 100.0, "recharge")
				service.On("Recharge", mock.Anything, "user123", 100.0, "alipay").Return(tx, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "充值成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "tx123", data["id"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWalletService)
			tt.setupMock(mockService)

			router := setupWalletTestRouter(mockService, tt.userID)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/shared/wallet/recharge", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}

func TestWalletAPI_Consume(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    sharedAPI.ConsumeRequest
		setupMock      func(*MockWalletService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "成功消费",
			userID: "user123",
			requestBody: sharedAPI.ConsumeRequest{
				Amount: 50.0,
				Reason: "购买VIP章节",
			},
			setupMock: func(service *MockWalletService) {
				tx := createTestTransaction("user123", 50.0, "consume")
				service.On("Consume", mock.Anything, "user123", 50.0, "购买VIP章节").Return(tx, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "消费成功", resp["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWalletService)
			tt.setupMock(mockService)

			router := setupWalletTestRouter(mockService, tt.userID)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/shared/wallet/consume", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}

func TestWalletAPI_Transfer(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    sharedAPI.TransferRequest
		setupMock      func(*MockWalletService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "成功转账",
			userID: "user123",
			requestBody: sharedAPI.TransferRequest{
				ToUserID: "user456",
				Amount:   30.0,
				Reason:   "打赏",
			},
			setupMock: func(service *MockWalletService) {
				tx := createTestTransaction("user123", 30.0, "transfer")
				service.On("Transfer", mock.Anything, "user123", "user456", 30.0, "打赏").Return(tx, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "转账成功", resp["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWalletService)
			tt.setupMock(mockService)

			router := setupWalletTestRouter(mockService, tt.userID)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/shared/wallet/transfer", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}

func TestWalletAPI_GetTransactions(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		queryParams    string
		setupMock      func(*MockWalletService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:        "成功获取交易记录",
			userID:      "user123",
			queryParams: "?page=1&page_size=20",
			setupMock: func(service *MockWalletService) {
				transactions := []*wallet.Transaction{
					createTestTransaction("user123", 100.0, "recharge"),
					createTestTransaction("user123", 50.0, "consume"),
				}
				service.On("ListTransactions", mock.Anything, "user123", mock.AnythingOfType("*wallet.ListTransactionsRequest")).Return(transactions, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				data := resp["data"].([]interface{})
				assert.Len(t, data, 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWalletService)
			tt.setupMock(mockService)

			router := setupWalletTestRouter(mockService, tt.userID)

			req := httptest.NewRequest("GET", "/api/v1/shared/wallet/transactions"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}

func TestWalletAPI_RequestWithdraw(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    sharedAPI.WithdrawRequest
		setupMock      func(*MockWalletService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "成功申请提现",
			userID: "user123",
			requestBody: sharedAPI.WithdrawRequest{
				Amount:  200.0,
				Account: "alipay:user@example.com",
			},
			setupMock: func(service *MockWalletService) {
				withdrawReq := &wallet.WithdrawRequest{
					ID:      "wd123",
					UserID:  "user123",
					Amount:  200.0,
					Account: "alipay:user@example.com",
					Status:  "pending",
				}
				service.On("RequestWithdraw", mock.Anything, "user123", 200.0, "alipay:user@example.com").Return(withdrawReq, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "申请提现成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "wd123", data["id"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWalletService)
			tt.setupMock(mockService)

			router := setupWalletTestRouter(mockService, tt.userID)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/shared/wallet/withdraw", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}

func TestWalletAPI_GetWithdrawRequests(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		queryParams    string
		setupMock      func(*MockWalletService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:        "成功获取提现申请列表",
			userID:      "user123",
			queryParams: "?page=1&page_size=20",
			setupMock: func(service *MockWalletService) {
				requests := []*wallet.WithdrawRequest{
					{
						ID:     "wd123",
						UserID: "user123",
						Amount: 200.0,
						Status: "pending",
					},
				}
				service.On("ListWithdrawRequests", mock.Anything, mock.AnythingOfType("*wallet.ListWithdrawRequestsRequest")).Return(requests, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				data := resp["data"].([]interface{})
				assert.Len(t, data, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWalletService)
			tt.setupMock(mockService)

			router := setupWalletTestRouter(mockService, tt.userID)

			req := httptest.NewRequest("GET", "/api/v1/shared/wallet/withdrawals"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}
