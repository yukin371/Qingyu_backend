package finance

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	walletsvc "Qingyu_backend/service/finance/wallet"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type walletServiceStub struct {
	getWalletFn        func(context.Context, string) (*walletsvc.Wallet, error)
	getBalanceFn       func(context.Context, string) (int64, error)
	rechargeFn         func(context.Context, string, int64, string) (*walletsvc.Transaction, error)
	listTransactionsFn func(context.Context, string, *walletsvc.ListTransactionsRequest) ([]*walletsvc.Transaction, error)
	requestWithdrawFn  func(context.Context, string, int64, string) (*walletsvc.WithdrawRequest, error)
	listWithdrawsFn    func(context.Context, *walletsvc.ListWithdrawRequestsRequest) ([]*walletsvc.WithdrawRequest, error)
}

func (m *walletServiceStub) CreateWallet(context.Context, string) (*walletsvc.Wallet, error) {
	return nil, nil
}
func (m *walletServiceStub) GetWallet(ctx context.Context, userID string) (*walletsvc.Wallet, error) {
	if m.getWalletFn == nil {
		return nil, nil
	}
	return m.getWalletFn(ctx, userID)
}
func (m *walletServiceStub) GetBalance(ctx context.Context, userID string) (int64, error) {
	if m.getBalanceFn == nil {
		return 0, nil
	}
	return m.getBalanceFn(ctx, userID)
}
func (m *walletServiceStub) FreezeWallet(context.Context, string) error   { return nil }
func (m *walletServiceStub) UnfreezeWallet(context.Context, string) error { return nil }
func (m *walletServiceStub) Recharge(ctx context.Context, userID string, amount int64, method string) (*walletsvc.Transaction, error) {
	if m.rechargeFn == nil {
		return nil, nil
	}
	return m.rechargeFn(ctx, userID, amount, method)
}
func (m *walletServiceStub) Consume(context.Context, string, int64, string) (*walletsvc.Transaction, error) {
	return nil, nil
}
func (m *walletServiceStub) Transfer(context.Context, string, string, int64, string) (*walletsvc.Transaction, error) {
	return nil, nil
}
func (m *walletServiceStub) GetTransaction(context.Context, string) (*walletsvc.Transaction, error) {
	return nil, nil
}
func (m *walletServiceStub) ListTransactions(ctx context.Context, userID string, req *walletsvc.ListTransactionsRequest) ([]*walletsvc.Transaction, error) {
	if m.listTransactionsFn == nil {
		return nil, nil
	}
	return m.listTransactionsFn(ctx, userID, req)
}
func (m *walletServiceStub) RequestWithdraw(ctx context.Context, userID string, amount int64, account string) (*walletsvc.WithdrawRequest, error) {
	if m.requestWithdrawFn == nil {
		return nil, nil
	}
	return m.requestWithdrawFn(ctx, userID, amount, account)
}
func (m *walletServiceStub) GetWithdrawRequest(context.Context, string) (*walletsvc.WithdrawRequest, error) {
	return nil, nil
}
func (m *walletServiceStub) ListWithdrawRequests(ctx context.Context, req *walletsvc.ListWithdrawRequestsRequest) ([]*walletsvc.WithdrawRequest, error) {
	if m.listWithdrawsFn == nil {
		return nil, nil
	}
	return m.listWithdrawsFn(ctx, req)
}
func (m *walletServiceStub) ApproveWithdraw(context.Context, string, string) error        { return nil }
func (m *walletServiceStub) RejectWithdraw(context.Context, string, string, string) error { return nil }
func (m *walletServiceStub) ProcessWithdraw(context.Context, string) error                { return nil }
func (m *walletServiceStub) Health(context.Context) error                                 { return nil }

func setupWalletAPITestRouter(service walletsvc.WalletService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := NewWalletAPI(service)
	v1 := router.Group("/api/v1/finance/wallet")
	v1.GET("", api.GetWallet)
	v1.GET("/balance", api.GetBalance)
	v1.POST("/recharge", api.Recharge)
	v1.GET("/transactions", api.GetTransactions)
	v1.POST("/withdraw", api.RequestWithdraw)
	v1.GET("/withdraws", api.GetWithdrawRequests)
	v1.GET("/withdrawals", api.GetWithdrawRequests)
	return router
}

func TestWalletAPIGetBalanceUnauthorized(t *testing.T) {
	router := setupWalletAPITestRouter(&walletServiceStub{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/wallet/balance", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWalletAPIGetBalanceSuccess(t *testing.T) {
	router := setupWalletAPITestRouter(&walletServiceStub{
		getBalanceFn: func(_ context.Context, userID string) (int64, error) {
			assert.Equal(t, "user-1", userID)
			return 1234, nil
		},
	}, "user-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/wallet/balance", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"data":1234`)
}

func TestWalletAPIGetWalletSuccess(t *testing.T) {
	router := setupWalletAPITestRouter(&walletServiceStub{
		getWalletFn: func(_ context.Context, userID string) (*walletsvc.Wallet, error) {
			assert.Equal(t, "user-1", userID)
			return &walletsvc.Wallet{
				ID:      "wallet-1",
				UserID:  userID,
				Balance: 2234,
			}, nil
		},
	}, "user-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/wallet", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "获取钱包信息成功")
	assert.Contains(t, w.Body.String(), `"balance":2234`)
}

func TestWalletAPIRechargeConvertsYuanToCents(t *testing.T) {
	router := setupWalletAPITestRouter(&walletServiceStub{
		rechargeFn: func(_ context.Context, userID string, amount int64, method string) (*walletsvc.Transaction, error) {
			assert.Equal(t, "user-1", userID)
			assert.Equal(t, int64(1234), amount)
			assert.Equal(t, "alipay", method)
			return &walletsvc.Transaction{
				ID:              "tx-1",
				UserID:          userID,
				Type:            "recharge",
				Amount:          amount,
				Method:          method,
				Status:          "success",
				TransactionTime: time.Now(),
				CreatedAt:       time.Now(),
			}, nil
		},
	}, "user-1")

	body := bytes.NewBufferString(`{"amount":12.34,"method":"alipay"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/finance/wallet/recharge", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "充值成功")
}

func TestWalletAPIGetTransactionsUsesDefaultPagination(t *testing.T) {
	router := setupWalletAPITestRouter(&walletServiceStub{
		listTransactionsFn: func(_ context.Context, userID string, req *walletsvc.ListTransactionsRequest) ([]*walletsvc.Transaction, error) {
			assert.Equal(t, "user-1", userID)
			assert.Equal(t, 1, req.Page)
			assert.Equal(t, 20, req.PageSize)
			assert.Equal(t, "", req.TransactionType)
			return []*walletsvc.Transaction{
				{
					ID:              "tx-1",
					UserID:          userID,
					Type:            "recharge",
					Amount:          500,
					Status:          "success",
					TransactionTime: time.Now(),
					CreatedAt:       time.Now(),
				},
			}, nil
		},
	}, "user-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/wallet/transactions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	assert.Equal(t, float64(0), payload["code"])
}

func TestWalletAPIGetTransactionsForwardsFilters(t *testing.T) {
	router := setupWalletAPITestRouter(&walletServiceStub{
		listTransactionsFn: func(_ context.Context, userID string, req *walletsvc.ListTransactionsRequest) ([]*walletsvc.Transaction, error) {
			assert.Equal(t, "user-1", userID)
			assert.Equal(t, 3, req.Page)
			assert.Equal(t, 50, req.PageSize)
			assert.Equal(t, "withdraw", req.TransactionType)
			return []*walletsvc.Transaction{}, nil
		},
	}, "user-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/wallet/transactions?page=3&page_size=50&type=withdraw", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestWalletAPIRequestWithdrawConvertsYuanToCents(t *testing.T) {
	router := setupWalletAPITestRouter(&walletServiceStub{
		requestWithdrawFn: func(_ context.Context, userID string, amount int64, account string) (*walletsvc.WithdrawRequest, error) {
			assert.Equal(t, "user-1", userID)
			assert.Equal(t, int64(1050), amount)
			assert.Equal(t, "acc-1", account)
			return &walletsvc.WithdrawRequest{
				ID:        "withdraw-1",
				UserID:    userID,
				Amount:    amount,
				Account:   account,
				Status:    "pending",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}, "user-1")

	body := bytes.NewBufferString(`{"amount":10.5,"method":"alipay","account":"acc-1","password":"secret"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/finance/wallet/withdraw", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "提现申请已提交")
}

func TestWalletAPIGetWithdrawRequestsSupportsCompatRoutes(t *testing.T) {
	for _, path := range []string{
		"/api/v1/finance/wallet/withdraws?page=2&page_size=5&status=pending",
		"/api/v1/finance/wallet/withdrawals?page=2&page_size=5&status=pending",
	} {
		t.Run(path, func(t *testing.T) {
			router := setupWalletAPITestRouter(&walletServiceStub{
				listWithdrawsFn: func(_ context.Context, req *walletsvc.ListWithdrawRequestsRequest) ([]*walletsvc.WithdrawRequest, error) {
					assert.Equal(t, "user-1", req.UserID)
					assert.Equal(t, "pending", req.Status)
					assert.Equal(t, 2, req.Page)
					assert.Equal(t, 5, req.PageSize)
					return []*walletsvc.WithdrawRequest{
						{
							ID:        "withdraw-1",
							UserID:    req.UserID,
							Amount:    1200,
							Account:   "acc-1",
							Status:    req.Status,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
					}, nil
				},
			}, "user-1")

			req := httptest.NewRequest(http.MethodGet, path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), "查询提现申请成功")
		})
	}
}
