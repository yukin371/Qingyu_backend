package finance

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	getEarningsFn          func(context.Context, string, int, int) ([]*financeModel.AuthorEarning, int64, error)
	getBookEarningsFn      func(context.Context, string, string) ([]*financeModel.AuthorEarning, int64, error)
	getRevenueDetailsFn    func(context.Context, string, int, int) ([]*financeModel.RevenueDetail, int64, error)
	getRevenueStatisticsFn func(context.Context, string, string) ([]*financeModel.RevenueStatistics, error)
	createWithdrawalReqFn  func(context.Context, string, float64, string, financeModel.WithdrawAccount) (*financeModel.WithdrawalRequest, error)
	getTaxInfoFn           func(context.Context, string) (*financeModel.TaxInfo, error)
	updateTaxInfoFn        func(context.Context, string, *financeModel.TaxInfo) error
}

var _ financeService.AuthorRevenueService = (*authorRevenueServiceStub)(nil)

func (s *authorRevenueServiceStub) GetEarnings(ctx context.Context, userID string, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error) {
	if s.getEarningsFn == nil {
		return nil, 0, nil
	}
	return s.getEarningsFn(ctx, userID, page, pageSize)
}

func (s *authorRevenueServiceStub) GetBookEarnings(ctx context.Context, userID string, bookID string) ([]*financeModel.AuthorEarning, int64, error) {
	if s.getBookEarningsFn == nil {
		return nil, 0, nil
	}
	return s.getBookEarningsFn(ctx, userID, bookID)
}

func (s *authorRevenueServiceStub) GetRevenueDetails(ctx context.Context, userID string, page, pageSize int) ([]*financeModel.RevenueDetail, int64, error) {
	if s.getRevenueDetailsFn == nil {
		return nil, 0, nil
	}
	return s.getRevenueDetailsFn(ctx, userID, page, pageSize)
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

func (s *authorRevenueServiceStub) GetTaxInfo(ctx context.Context, userID string) (*financeModel.TaxInfo, error) {
	if s.getTaxInfoFn == nil {
		return nil, nil
	}
	return s.getTaxInfoFn(ctx, userID)
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
	v1.GET("/earnings", api.GetEarnings)
	v1.GET("/earnings/:bookId", api.GetBookEarnings)
	v1.GET("/revenue-details", api.GetRevenueDetails)
	v1.GET("/revenue-statistics", api.GetRevenueStatistics)
	v1.POST("/withdraw", api.Withdraw)
	v1.GET("/tax-info", api.GetTaxInfo)
	v1.PUT("/tax-info", api.UpdateTaxInfo)
	return router
}

func TestAuthorRevenueAPIGetEarningsPaginationAndFailureBranches(t *testing.T) {
	t.Run("rejects unauthenticated requests", func(t *testing.T) {
		router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{}, "")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/earnings", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "未认证")
	})

	t.Run("defaults pagination and returns empty data", func(t *testing.T) {
		router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
			getEarningsFn: func(_ context.Context, userID string, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error) {
				assert.Equal(t, "author-1", userID)
				assert.Equal(t, 1, page)
				assert.Equal(t, 20, pageSize)
				return []*financeModel.AuthorEarning{}, 0, nil
			},
		}, "author-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/earnings", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "获取收入列表成功")

		var resp struct {
			Total int              `json:"total"`
			Data  []map[string]any `json:"data"`
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, 0, resp.Total)
		require.Len(t, resp.Data, 0)
	})

	t.Run("surfaces service errors", func(t *testing.T) {
		router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
			getEarningsFn: func(context.Context, string, int, int) ([]*financeModel.AuthorEarning, int64, error) {
				return nil, 0, errors.New("revenue chain failed")
			},
		}, "author-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/earnings?page=3&page_size=5", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "revenue chain failed")
	})
}

func TestAuthorRevenueAPIGetBookEarningsValidationAndMapping(t *testing.T) {
	t.Run("rejects unauthenticated requests", func(t *testing.T) {
		called := false
		router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
			getBookEarningsFn: func(context.Context, string, string) ([]*financeModel.AuthorEarning, int64, error) {
				called = true
				return nil, 0, nil
			},
		}, "")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/earnings/"+primitive.NewObjectID().Hex(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "未认证")
		assert.False(t, called)
	})

	t.Run("surfaces service failures", func(t *testing.T) {
		bookID := primitive.NewObjectID()
		router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
			getBookEarningsFn: func(_ context.Context, userID string, requestedBookID string) ([]*financeModel.AuthorEarning, int64, error) {
				assert.Equal(t, "author-1", userID)
				assert.Equal(t, bookID.Hex(), requestedBookID)
				return nil, 0, errors.New("book earnings query failed")
			},
		}, "author-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/earnings/"+bookID.Hex(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "book earnings query failed")
	})

	t.Run("maps service response for a specific book", func(t *testing.T) {
		now := time.Now()
		bookID := primitive.NewObjectID()
		chapterID := primitive.NewObjectID()
		router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
			getBookEarningsFn: func(_ context.Context, userID string, requestedBookID string) ([]*financeModel.AuthorEarning, int64, error) {
				assert.Equal(t, "author-1", userID)
				assert.Equal(t, bookID.Hex(), requestedBookID)
				return []*financeModel.AuthorEarning{
					{
						ID:             primitive.NewObjectID(),
						AuthorID:       userID,
						BookID:         bookID,
						BookTitle:      "星海回响",
						ChapterID:      chapterID,
						ChapterTitle:   "第一章",
						Type:           financeModel.EarningTypeChapterPurchase,
						Amount:         types.Money(4500),
						PlatformFee:    types.Money(900),
						AuthorIncome:   types.Money(3600),
						ReaderID:       "reader-9",
						ReaderNickname: "夜读者",
						IsSettled:      true,
						CreatedAt:      now,
						UpdatedAt:      now,
					},
				}, 1, nil
			},
		}, "author-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/earnings/"+bookID.Hex()+"?page=2&page_size=3", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "获取书籍收入成功")
		assert.Contains(t, w.Body.String(), `"book_title":"星海回响"`)
		assert.Contains(t, w.Body.String(), `"amount":36`)
		assert.Contains(t, w.Body.String(), `"gross_amount_cents":4500`)
		assert.Contains(t, w.Body.String(), `"platform_fee_cents":900`)
	})
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

func TestAuthorRevenueAPIGetRevenueDetailsMapsResponseAndServiceFailure(t *testing.T) {
	t.Run("maps details with explicit pagination", func(t *testing.T) {
		now := time.Now()
		firstAt := now.Add(-48 * time.Hour)
		lastAt := now.Add(-24 * time.Hour)
		router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
			getRevenueDetailsFn: func(_ context.Context, userID string, page, pageSize int) ([]*financeModel.RevenueDetail, int64, error) {
				assert.Equal(t, "author-1", userID)
				assert.Equal(t, 2, page)
				assert.Equal(t, 5, pageSize)
				return []*financeModel.RevenueDetail{
					{
						ID:               primitive.NewObjectID(),
						AuthorID:         userID,
						BookID:           primitive.NewObjectID(),
						BookTitle:        "测试作品",
						Type:             financeModel.EarningTypeReward,
						TotalAmount:      types.Money(3560),
						TotalIncome:      types.Money(2800),
						TransactionCount: 7,
						FirstEarningAt:   &firstAt,
						LastEarningAt:    &lastAt,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
				}, 1, nil
			},
		}, "author-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/revenue-details?page=2&page_size=5", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "获取收入明细成功")
		assert.Contains(t, w.Body.String(), `"book_title":"测试作品"`)
		assert.Contains(t, w.Body.String(), `"total_amount":35.6`)
		assert.Contains(t, w.Body.String(), `"total_income_cents":2800`)
	})

	t.Run("surfaces service errors", func(t *testing.T) {
		router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
			getRevenueDetailsFn: func(context.Context, string, int, int) ([]*financeModel.RevenueDetail, int64, error) {
				return nil, 0, errors.New("detail aggregation failed")
			},
		}, "author-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/revenue-details", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "detail aggregation failed")
	})
}

func TestAuthorRevenueAPIGetRevenueStatisticsSurfacesServiceFailure(t *testing.T) {
	router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
		getRevenueStatisticsFn: func(_ context.Context, userID string, period string) ([]*financeModel.RevenueStatistics, error) {
			assert.Equal(t, "author-1", userID)
			assert.Equal(t, "yearly", period)
			return nil, errors.New("statistics backend unavailable")
		},
	}, "author-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/revenue-statistics?period=yearly", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "statistics backend unavailable")
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

func TestAuthorRevenueAPIWithdrawSurfacesServiceFailure(t *testing.T) {
	router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
		createWithdrawalReqFn: func(_ context.Context, userID string, amount float64, method string, account financeModel.WithdrawAccount) (*financeModel.WithdrawalRequest, error) {
			assert.Equal(t, "author-1", userID)
			assert.Equal(t, 88.8, amount)
			assert.Equal(t, "alipay", method)
			assert.Equal(t, "收款人", account.AccountName)
			return nil, errors.New("insufficient settled balance")
		},
	}, "author-1")

	body := bytes.NewBufferString(`{"amount":88.8,"method":"alipay","account_type":"alipay","account_name":"收款人","account_no":"author@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/finance/author/withdraw", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "insufficient settled balance")
}

func TestAuthorRevenueAPIWithdrawRejectsInvalidPayload(t *testing.T) {
	called := false
	router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
		createWithdrawalReqFn: func(_ context.Context, userID string, amount float64, method string, account financeModel.WithdrawAccount) (*financeModel.WithdrawalRequest, error) {
			called = true
			return nil, nil
		},
	}, "author-1")

	body := bytes.NewBufferString(`{"amount":0,"method":"cash","account_type":"bank","account_name":"张三","account_no":"622200"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/finance/author/withdraw", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "请求参数错误")
	assert.Contains(t, w.Body.String(), "Amount")
	assert.False(t, called)
}

func TestAuthorRevenueAPIGetTaxInfoSuccessAndFailure(t *testing.T) {
	t.Run("rejects unauthenticated requests", func(t *testing.T) {
		called := false
		router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
			getTaxInfoFn: func(_ context.Context, userID string) (*financeModel.TaxInfo, error) {
				called = true
				return &financeModel.TaxInfo{UserID: userID}, nil
			},
		}, "")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/tax-info", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "未认证")
		assert.False(t, called)
	})

	t.Run("returns success with null data when profile is empty", func(t *testing.T) {
		router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
			getTaxInfoFn: func(_ context.Context, userID string) (*financeModel.TaxInfo, error) {
				assert.Equal(t, "author-1", userID)
				return nil, nil
			},
		}, "author-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/tax-info", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "获取税务信息成功")
		assert.Contains(t, w.Body.String(), `"data":null`)
	})

	t.Run("returns tax info payload", func(t *testing.T) {
		router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
			getTaxInfoFn: func(_ context.Context, userID string) (*financeModel.TaxInfo, error) {
				assert.Equal(t, "author-1", userID)
				return &financeModel.TaxInfo{
					UserID:     userID,
					IDType:     "id_card",
					IDNumber:   "310101199001011234",
					Name:       "张三",
					TaxType:    "individual",
					TaxRate:    0.12,
					IsVerified: true,
				}, nil
			},
		}, "author-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/tax-info", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "获取税务信息成功")
		assert.Contains(t, w.Body.String(), `"id_type":"id_card"`)
		assert.Contains(t, w.Body.String(), `"tax_rate":0.12`)
		assert.Contains(t, w.Body.String(), `"is_verified":true`)
	})

	t.Run("surfaces service failure", func(t *testing.T) {
		router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
			getTaxInfoFn: func(_ context.Context, userID string) (*financeModel.TaxInfo, error) {
				assert.Equal(t, "author-1", userID)
				return nil, errors.New("tax profile missing")
			},
		}, "author-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/author/tax-info", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "tax profile missing")
	})
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

func TestAuthorRevenueAPIUpdateTaxInfoSurfacesServiceFailure(t *testing.T) {
	router := setupAuthorRevenueAPITestRouter(&authorRevenueServiceStub{
		updateTaxInfoFn: func(_ context.Context, userID string, taxInfo *financeModel.TaxInfo) error {
			assert.Equal(t, "author-1", userID)
			require.NotNil(t, taxInfo)
			assert.Equal(t, "passport", taxInfo.IDType)
			assert.Equal(t, "company", taxInfo.TaxType)
			return errors.New("tax profile sync failed")
		},
	}, "author-1")

	body := bytes.NewBufferString(`{"id_type":"passport","id_number":"P1234567","name":"Shanghai Studio","tax_type":"company"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/finance/author/tax-info", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "tax profile sync failed")
}
