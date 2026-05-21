package finance

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	financeModel "Qingyu_backend/models/finance"
	"Qingyu_backend/models/shared/types"
	financeService "Qingyu_backend/service/finance"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type authorRevenueServiceStub struct {
	getRevenueStatisticsFn func(context.Context, string, string) ([]*financeModel.RevenueStatistics, error)
	createWithdrawalReqFn  func(context.Context, string, float64, string, financeModel.WithdrawAccount) (*financeModel.WithdrawalRequest, error)
	updateTaxInfoFn        func(context.Context, string, *financeModel.TaxInfo) error
}

var _ financeService.AuthorRevenueService = (*authorRevenueServiceStub)(nil)

func (s *authorRevenueServiceStub) GetEarnings(context.Context, string, int, int) ([]*financeModel.AuthorEarning, int64, error) {
	return nil, 0, nil
}

func (s *authorRevenueServiceStub) GetBookEarnings(context.Context, string, string) ([]*financeModel.AuthorEarning, int64, error) {
	return nil, 0, nil
}

func (s *authorRevenueServiceStub) GetRevenueDetails(context.Context, string, int, int) ([]*financeModel.RevenueDetail, int64, error) {
	return nil, 0, nil
}

func (s *authorRevenueServiceStub) GetRevenueStatistics(ctx context.Context, userID string, period string) ([]*financeModel.RevenueStatistics, error) {
	if s.getRevenueStatisticsFn == nil {
		return nil, nil
	}
	return s.getRevenueStatisticsFn(ctx, userID, period)
}

func (s *authorRevenueServiceStub) CreateEarning(context.Context, *financeModel.AuthorEarning) error {
	return nil
}

func (s *authorRevenueServiceStub) CalculateEarning(context.Context, string, float64, string, primitive.ObjectID) (float64, float64, error) {
	return 0, 0, nil
}

func (s *authorRevenueServiceStub) CreateWithdrawalRequest(ctx context.Context, userID string, amount float64, method string, account financeModel.WithdrawAccount) (*financeModel.WithdrawalRequest, error) {
	if s.createWithdrawalReqFn == nil {
		return nil, nil
	}
	return s.createWithdrawalReqFn(ctx, userID, amount, method, account)
}

func (s *authorRevenueServiceStub) GetWithdrawals(context.Context, string, int, int) ([]*financeModel.WithdrawalRequest, int64, error) {
	return nil, 0, nil
}

func (s *authorRevenueServiceStub) GetSettlements(context.Context, string, int, int) ([]*financeModel.Settlement, int64, error) {
	return nil, 0, nil
}

func (s *authorRevenueServiceStub) GetSettlement(context.Context, string) (*financeModel.Settlement, error) {
	return nil, nil
}

func (s *authorRevenueServiceStub) GetTaxInfo(context.Context, string) (*financeModel.TaxInfo, error) {
	return nil, nil
}

func (s *authorRevenueServiceStub) UpdateTaxInfo(ctx context.Context, userID string, taxInfo *financeModel.TaxInfo) error {
	if s.updateTaxInfoFn == nil {
		return nil
	}
	return s.updateTaxInfoFn(ctx, userID, taxInfo)
}

func setupAuthorRevenueAPITestRouter(service financeService.AuthorRevenueService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := NewAuthorRevenueAPI(service)
	v1 := router.Group("/api/v1/finance/author")
	v1.GET("/revenue-statistics", api.GetRevenueStatistics)
	v1.POST("/withdraw", api.Withdraw)
	v1.PUT("/tax-info", api.UpdateTaxInfo)
	return router
}

func TestAuthorRevenueAPIGetRevenueStatisticsDefaultsPeriod(t *testing.T) {
	router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
		getRevenueStatisticsFn: func(_ context.Context, userID string, period string) ([]*financeModel.RevenueStatistics, error) {
			assert.Equal(t, "author-1", userID)
			assert.Equal(t, "monthly", period)
			return []*financeModel.RevenueStatistics{
				{
					AuthorID:         userID,
					Period:           period,
					TotalRevenue:     types.Money(12800),
					ChapterIncome:    types.Money(10000),
					RewardIncome:     types.Money(1800),
					VIPReadingIncome: types.Money(1000),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
			}, nil
		},
	}, "author-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/revenue-statistics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "获取收入统计成功")
	assert.Contains(t, w.Body.String(), `"period":"monthly"`)
}

func TestAuthorRevenueAPIWithdrawPassesAccountInfo(t *testing.T) {
	router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
		createWithdrawalReqFn: func(_ context.Context, userID string, amount float64, method string, account financeModel.WithdrawAccount) (*financeModel.WithdrawalRequest, error) {
			assert.Equal(t, "author-1", userID)
			assert.Equal(t, 12.5, amount)
			assert.Equal(t, "bank", method)
			assert.Equal(t, "bank", account.AccountType)
			assert.Equal(t, "张三", account.AccountName)
			assert.Equal(t, "622200", account.AccountNo)
			assert.Equal(t, "招商银行", account.BankName)
			assert.Equal(t, "上海分行", account.BranchName)
			return &financeModel.WithdrawalRequest{
				ID:           primitive.NewObjectID(),
				UserID:       userID,
				Amount:       types.Money(1250),
				ActualAmount: types.Money(1250),
				Method:       method,
				AccountInfo:  account,
				Status:       financeModel.WithdrawStatusPending,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}, nil
		},
	}, "author-1")

	body := bytes.NewBufferString(`{"amount":12.5,"method":"bank","account_type":"bank","account_name":"张三","account_no":"622200","bank_name":"招商银行","branch_name":"上海分行"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/finance/author/withdraw", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "申请提现成功")
	assert.Contains(t, w.Body.String(), `"account_name":"张三"`)
}

func TestAuthorRevenueAPIUpdateTaxInfoMapsRequest(t *testing.T) {
	router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
		updateTaxInfoFn: func(_ context.Context, userID string, taxInfo *financeModel.TaxInfo) error {
			assert.Equal(t, "author-1", userID)
			require.NotNil(t, taxInfo)
			assert.Equal(t, "id_card", taxInfo.IDType)
			assert.Equal(t, "310101199001011234", taxInfo.IDNumber)
			assert.Equal(t, "张三", taxInfo.Name)
			assert.Equal(t, "individual", taxInfo.TaxType)
			assert.Equal(t, 0.0, taxInfo.TaxRate)
			return nil
		},
	}, "author-1")

	body := bytes.NewBufferString(`{"id_type":"id_card","id_number":"310101199001011234","name":"张三","tax_type":"individual"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/finance/author/tax-info", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "更新税务信息成功")
}
