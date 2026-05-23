package finance

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

type membershipServiceStub struct {
	getPlansFn         func(context.Context) ([]*financeModel.MembershipPlan, error)
	subscribeFn        func(context.Context, string, string, string) (*financeModel.UserMembership, error)
	getMembershipFn    func(context.Context, string) (*financeModel.UserMembership, error)
	cancelMembershipFn func(context.Context, string) error
	renewMembershipFn  func(context.Context, string) (*financeModel.UserMembership, error)
	listCardsFn        func(context.Context, map[string]interface{}, int, int) ([]*financeModel.MembershipCard, int64, error)
}

var _ financeService.MembershipService = (*membershipServiceStub)(nil)

func (s *membershipServiceStub) GetPlans(ctx context.Context) ([]*financeModel.MembershipPlan, error) {
	if s.getPlansFn != nil {
		return s.getPlansFn(ctx)
	}
	return nil, nil
}

func (s *membershipServiceStub) GetPlan(context.Context, string) (*financeModel.MembershipPlan, error) {
	return nil, nil
}

func (s *membershipServiceStub) Subscribe(ctx context.Context, userID, planID, paymentMethod string) (*financeModel.UserMembership, error) {
	if s.subscribeFn != nil {
		return s.subscribeFn(ctx, userID, planID, paymentMethod)
	}
	return nil, nil
}

func (s *membershipServiceStub) GetMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error) {
	if s.getMembershipFn != nil {
		return s.getMembershipFn(ctx, userID)
	}
	return nil, nil
}

func (s *membershipServiceStub) CancelMembership(ctx context.Context, userID string) error {
	if s.cancelMembershipFn != nil {
		return s.cancelMembershipFn(ctx, userID)
	}
	return nil
}

func (s *membershipServiceStub) RenewMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error) {
	if s.renewMembershipFn != nil {
		return s.renewMembershipFn(ctx, userID)
	}
	return nil, nil
}

func (s *membershipServiceStub) GetBenefits(context.Context, string) ([]*financeModel.MembershipBenefit, error) {
	return nil, nil
}

func (s *membershipServiceStub) GetUsage(context.Context, string) ([]*financeModel.MembershipUsage, error) {
	return nil, nil
}

func (s *membershipServiceStub) ActivateCard(context.Context, string, string) (*financeModel.UserMembership, error) {
	return nil, nil
}

func (s *membershipServiceStub) ListCards(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.MembershipCard, int64, error) {
	if s.listCardsFn != nil {
		return s.listCardsFn(ctx, filter, page, pageSize)
	}
	return nil, 0, nil
}

func (s *membershipServiceStub) CheckMembership(context.Context, string, string) (bool, error) {
	return false, nil
}

func (s *membershipServiceStub) IsVIP(context.Context, string) (bool, error) {
	return false, nil
}

func setupMembershipAPITestRouter(service financeService.MembershipService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := NewMembershipAPI(service)
	v1 := router.Group("/api/v1/finance/membership")
	v1.GET("/plans", api.GetPlans)
	v1.POST("/subscribe", api.Subscribe)
	v1.GET("/status", api.GetStatus)
	v1.POST("/cancel", api.Cancel)
	v1.PUT("/renew", api.Renew)
	v1.GET("/cards", api.ListCards)
	return router
}

func TestMembershipAPIGetPlansListCardsAndAuthGuards(t *testing.T) {
	t.Run("get plans returns data", func(t *testing.T) {
		router := setupMembershipAPITestRouter(&membershipServiceStub{
			getPlansFn: func(context.Context) ([]*financeModel.MembershipPlan, error) {
				return []*financeModel.MembershipPlan{
					{
						Name:      "月度VIP",
						Type:      financeModel.MembershipTypeMonthly,
						Price:     types.NewMoneyFromYuan(19.9),
						IsEnabled: true,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil
			},
		}, "user-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/membership/plans", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "获取套餐列表成功")
		assert.Contains(t, w.Body.String(), "月度VIP")
	})

	t.Run("list cards uses defaults and forwards filters", func(t *testing.T) {
		router := setupMembershipAPITestRouter(&membershipServiceStub{
			listCardsFn: func(_ context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.MembershipCard, int64, error) {
				assert.Equal(t, 1, page)
				assert.Equal(t, 20, pageSize)
				assert.Equal(t, "enabled", filter["status"])
				return []*financeModel.MembershipCard{
					{Code: "CARD001", Status: financeModel.CardStatusUnused, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				}, 1, nil
			},
		}, "admin-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/membership/cards?status=enabled", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "获取会员卡列表成功")
	})

	t.Run("subscribe rejects unauthenticated request", func(t *testing.T) {
		router := setupMembershipAPITestRouter(&membershipServiceStub{}, "")
		req := httptest.NewRequest(http.MethodPost, "/api/v1/finance/membership/subscribe", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "未认证")
	})

	t.Run("get status returns membership for authenticated user", func(t *testing.T) {
		now := time.Now()
		router := setupMembershipAPITestRouter(&membershipServiceStub{
			getMembershipFn: func(_ context.Context, userID string) (*financeModel.UserMembership, error) {
				assert.Equal(t, "reader-1", userID)
				return &financeModel.UserMembership{
					UserID:    userID,
					PlanName:  "年度VIP",
					PlanType:  financeModel.MembershipTypeYearly,
					Level:     "vip_yearly",
					AutoRenew: true,
					Status:    "active",
					StartTime: now,
					EndTime:   now.Add(365 * 24 * time.Hour),
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
		}, "reader-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/membership/status", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "获取会员状态成功")
		assert.Contains(t, w.Body.String(), "\"plan_name\":\"年度VIP\"")
		assert.Contains(t, w.Body.String(), "\"status\":\"active\"")
	})

	t.Run("cancel surfaces service failure for authenticated user", func(t *testing.T) {
		router := setupMembershipAPITestRouter(&membershipServiceStub{
			cancelMembershipFn: func(_ context.Context, userID string) error {
				assert.Equal(t, "reader-1", userID)
				return errors.New("cancel membership failed")
			},
		}, "reader-1")

		req := httptest.NewRequest(http.MethodPost, "/api/v1/finance/membership/cancel", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "cancel membership failed")
	})

	t.Run("renew returns updated membership", func(t *testing.T) {
		now := time.Now()
		router := setupMembershipAPITestRouter(&membershipServiceStub{
			renewMembershipFn: func(_ context.Context, userID string) (*financeModel.UserMembership, error) {
				assert.Equal(t, "reader-1", userID)
				return &financeModel.UserMembership{
					UserID:    userID,
					PlanName:  "月度VIP",
					PlanType:  financeModel.MembershipTypeMonthly,
					Level:     "vip_monthly",
					AutoRenew: false,
					Status:    "active",
					StartTime: now,
					EndTime:   now.Add(30 * 24 * time.Hour),
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
		}, "reader-1")

		req := httptest.NewRequest(http.MethodPut, "/api/v1/finance/membership/renew", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "续费成功")
		assert.Contains(t, w.Body.String(), "\"plan_type\":\"monthly\"")
		assert.Contains(t, w.Body.String(), "\"level\":\"vip_monthly\"")
	})
}

func TestMembershipAPIGetPlansFailureAndAuthGuards(t *testing.T) {
	t.Run("get plans surfaces service failure", func(t *testing.T) {
		router := setupMembershipAPITestRouter(&membershipServiceStub{
			getPlansFn: func(context.Context) ([]*financeModel.MembershipPlan, error) {
				return nil, errors.New("list plans failed")
			},
		}, "user-1")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/membership/plans", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "list plans failed")
	})

	t.Run("get status rejects unauthenticated request", func(t *testing.T) {
		router := setupMembershipAPITestRouter(&membershipServiceStub{}, "")
		req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/membership/status", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "未认证")
	})

	t.Run("renew rejects unauthenticated request", func(t *testing.T) {
		router := setupMembershipAPITestRouter(&membershipServiceStub{}, "")
		req := httptest.NewRequest(http.MethodPut, "/api/v1/finance/membership/renew", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "未认证")
	})
}

func TestMembershipAPISubscribeSuccess(t *testing.T) {
	now := time.Now()
	planID := primitive.NewObjectID()
	router := setupMembershipAPITestRouter(&membershipServiceStub{
		subscribeFn: func(_ context.Context, userID, requestedPlanID, paymentMethod string) (*financeModel.UserMembership, error) {
			assert.Equal(t, "reader-1", userID)
			assert.Equal(t, "plan-yearly", requestedPlanID)
			assert.Equal(t, "wallet", paymentMethod)
			return &financeModel.UserMembership{
				UserID:    userID,
				PlanID:    planID,
				PlanName:  "年度VIP",
				PlanType:  financeModel.MembershipTypeYearly,
				Level:     "vip_yearly",
				AutoRenew: true,
				Status:    "active",
				StartTime: now,
				EndTime:   now.Add(365 * 24 * time.Hour),
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}, "reader-1")

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/finance/membership/subscribe",
		strings.NewReader(`{"plan_id":"plan-yearly","payment_method":"wallet"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "订阅成功")
	assert.Contains(t, w.Body.String(), "\"plan_id\":\""+planID.Hex()+"\"")
	assert.Contains(t, w.Body.String(), "\"plan_name\":\"年度VIP\"")
}

func TestMembershipAPISubscribeValidationFailure(t *testing.T) {
	serviceCalled := false
	router := setupMembershipAPITestRouter(&membershipServiceStub{
		subscribeFn: func(_ context.Context, _, _, _ string) (*financeModel.UserMembership, error) {
			serviceCalled = true
			return nil, nil
		},
	}, "reader-1")

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/finance/membership/subscribe",
		strings.NewReader(`{"plan_id":"plan-yearly","payment_method":"cash"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, serviceCalled)
	assert.Contains(t, w.Body.String(), "请求参数错误")
	assert.Contains(t, w.Body.String(), "oneof")
}

func TestMembershipAPISubscribeMissingPlanID(t *testing.T) {
	serviceCalled := false
	router := setupMembershipAPITestRouter(&membershipServiceStub{
		subscribeFn: func(_ context.Context, _, _, _ string) (*financeModel.UserMembership, error) {
			serviceCalled = true
			return nil, nil
		},
	}, "reader-1")

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/finance/membership/subscribe",
		strings.NewReader(`{"payment_method":"wallet"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, serviceCalled)
	assert.Contains(t, w.Body.String(), "请求参数错误")
	assert.Contains(t, w.Body.String(), "PlanID")
}

func TestMembershipAPISubscribeInvalidJSON(t *testing.T) {
	serviceCalled := false
	router := setupMembershipAPITestRouter(&membershipServiceStub{
		subscribeFn: func(_ context.Context, _, _, _ string) (*financeModel.UserMembership, error) {
			serviceCalled = true
			return nil, nil
		},
	}, "reader-1")

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/finance/membership/subscribe",
		strings.NewReader(`{"plan_id":"plan-yearly","payment_method":"wallet"`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, serviceCalled)
	assert.Contains(t, w.Body.String(), "请求参数错误")
	assert.Contains(t, w.Body.String(), "unexpected EOF")
}

func TestMembershipAPIGetStatusServiceFailure(t *testing.T) {
	router := setupMembershipAPITestRouter(&membershipServiceStub{
		getMembershipFn: func(_ context.Context, userID string) (*financeModel.UserMembership, error) {
			assert.Equal(t, "reader-1", userID)
			return nil, errors.New("membership lookup failed")
		},
	}, "reader-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/membership/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "membership lookup failed")
}

func TestMembershipAPIGetStatusNilMembership(t *testing.T) {
	router := setupMembershipAPITestRouter(&membershipServiceStub{
		getMembershipFn: func(_ context.Context, userID string) (*financeModel.UserMembership, error) {
			assert.Equal(t, "reader-1", userID)
			return nil, nil
		},
	}, "reader-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/finance/membership/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "获取会员状态成功")
	assert.Contains(t, w.Body.String(), "\"data\":null")
}

func TestMembershipAPICancelSuccess(t *testing.T) {
	router := setupMembershipAPITestRouter(&membershipServiceStub{
		cancelMembershipFn: func(_ context.Context, userID string) error {
			assert.Equal(t, "reader-1", userID)
			return nil
		},
	}, "reader-1")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/finance/membership/cancel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "取消自动续费成功")
}

func TestMembershipAPICancelUnauthorized(t *testing.T) {
	router := setupMembershipAPITestRouter(&membershipServiceStub{}, "")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/finance/membership/cancel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "未认证")
}

func TestMembershipAPIRenewServiceFailure(t *testing.T) {
	router := setupMembershipAPITestRouter(&membershipServiceStub{
		renewMembershipFn: func(_ context.Context, userID string) (*financeModel.UserMembership, error) {
			assert.Equal(t, "reader-1", userID)
			return nil, errors.New("renew membership failed")
		},
	}, "reader-1")

	req := httptest.NewRequest(http.MethodPut, "/api/v1/finance/membership/renew", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "renew membership failed")
}
