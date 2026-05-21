package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Qingyu_backend/config"
	aiModels "Qingyu_backend/models/ai"
	aiRepo "Qingyu_backend/repository/interfaces/ai"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type quotaPolicyAPITestRepo struct {
	listItems []*aiModels.QuotaPolicy
	listTotal int64

	lastListRole   string
	lastListStatus string
	lastListPage   int
	lastListLimit  int

	createdPolicies []*aiModels.QuotaPolicy
	updatedPolicies []*aiModels.QuotaPolicy
	deletedIDs      []string

	policiesByID           map[string]*aiModels.QuotaPolicy
	policiesByRoleAndLevel map[string]*aiModels.QuotaPolicy

	listErr               error
	createErr             error
	getByIDErr            error
	getByRoleAndLevelErr  error
	updateErr             error
	deleteErr             error
}

func newQuotaPolicyAPITestRepo() *quotaPolicyAPITestRepo {
	return &quotaPolicyAPITestRepo{
		policiesByID:           make(map[string]*aiModels.QuotaPolicy),
		policiesByRoleAndLevel: make(map[string]*aiModels.QuotaPolicy),
	}
}

func quotaPolicyAPIRepoKey(role aiModels.UserRole, level aiModels.MembershipLevel) string {
	return string(role) + "|" + string(level)
}

func (s *quotaPolicyAPITestRepo) Create(ctx context.Context, policy *aiModels.QuotaPolicy) error {
	if s.createErr != nil {
		return s.createErr
	}
	if policy.ID.IsZero() {
		policy.ID = primitive.NewObjectID()
	}
	s.createdPolicies = append(s.createdPolicies, policy)
	s.policiesByID[policy.ID.Hex()] = policy
	s.policiesByRoleAndLevel[quotaPolicyAPIRepoKey(policy.UserRole, policy.MembershipLevel)] = policy
	return nil
}

func (s *quotaPolicyAPITestRepo) GetByID(ctx context.Context, id string) (*aiModels.QuotaPolicy, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	if policy, ok := s.policiesByID[id]; ok {
		return policy, nil
	}
	return nil, errors.New("not found")
}

func (s *quotaPolicyAPITestRepo) GetByRoleAndLevel(ctx context.Context, role aiModels.UserRole, level aiModels.MembershipLevel) (*aiModels.QuotaPolicy, error) {
	if s.getByRoleAndLevelErr != nil {
		return nil, s.getByRoleAndLevelErr
	}
	if policy, ok := s.policiesByRoleAndLevel[quotaPolicyAPIRepoKey(role, level)]; ok {
		return policy, nil
	}
	return nil, errors.New("not found")
}

func (s *quotaPolicyAPITestRepo) List(ctx context.Context, role string, status string, page, limit int) ([]*aiModels.QuotaPolicy, int64, error) {
	s.lastListRole = role
	s.lastListStatus = status
	s.lastListPage = page
	s.lastListLimit = limit
	if s.listErr != nil {
		return nil, 0, s.listErr
	}
	return s.listItems, s.listTotal, nil
}

func (s *quotaPolicyAPITestRepo) Update(ctx context.Context, policy *aiModels.QuotaPolicy) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.updatedPolicies = append(s.updatedPolicies, policy)
	s.policiesByID[policy.ID.Hex()] = policy
	s.policiesByRoleAndLevel[quotaPolicyAPIRepoKey(policy.UserRole, policy.MembershipLevel)] = policy
	return nil
}

func (s *quotaPolicyAPITestRepo) Delete(ctx context.Context, id string) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}
	s.deletedIDs = append(s.deletedIDs, id)
	delete(s.policiesByID, id)
	return nil
}

func (s *quotaPolicyAPITestRepo) Health(ctx context.Context) error {
	return nil
}

var _ aiRepo.QuotaPolicyRepository = (*quotaPolicyAPITestRepo)(nil)

func setupQuotaPolicyAPITestRouter(repo *quotaPolicyAPITestRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	api := NewQuotaPolicyAPI(aiService.NewQuotaPolicyService(repo))

	router := gin.New()
	router.GET("/policies", api.ListPolicies)
	router.POST("/policies", api.CreatePolicy)
	router.GET("/policies/:id", api.GetPolicy)
	router.PUT("/policies/:id", api.UpdatePolicy)
	router.DELETE("/policies/:id", api.DeletePolicy)
	router.POST("/policies/initialize", api.InitializeDefaultPolicies)
	return router
}

func TestQuotaPolicyAPIListPoliciesNormalizesPagingAndForwardsFilters(t *testing.T) {
	repo := newQuotaPolicyAPITestRepo()
	repo.listItems = []*aiModels.QuotaPolicy{
		{
			ID:         primitive.NewObjectID(),
			Name:       "reader-normal",
			UserRole:   aiModels.UserRoleReader,
			Status:     aiModels.QuotaPolicyStatusActive,
			DailyQuota: 120,
		},
	}
	repo.listTotal = 1
	router := setupQuotaPolicyAPITestRouter(repo)

	req, _ := http.NewRequest(http.MethodGet, "/policies?role=writer&status=active&page=0&limit=101", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "writer", repo.lastListRole)
	assert.Equal(t, "active", repo.lastListStatus)
	assert.Equal(t, 1, repo.lastListPage)
	assert.Equal(t, 20, repo.lastListLimit)

	var resp struct {
		Page  int              `json:"page"`
		Size  int              `json:"size"`
		Total int64            `json:"total"`
		Data  []map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 20, resp.Size)
	assert.Equal(t, int64(1), resp.Total)
	require.Len(t, resp.Data, 1)
}

func TestQuotaPolicyAPICreatePolicyMapsRequestAndSurfacesErrors(t *testing.T) {
	t.Run("maps request body to policy entity", func(t *testing.T) {
		repo := newQuotaPolicyAPITestRepo()
		router := setupQuotaPolicyAPITestRouter(repo)

		reqBody := `{"name":"VIP writer","userRole":"writer","membershipLevel":"vip_monthly","dailyQuota":321,"monthlyQuota":6543,"totalQuota":-1,"description":"custom"}`
		req, _ := http.NewRequest(http.MethodPost, "/policies", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Len(t, repo.createdPolicies, 1)
		created := repo.createdPolicies[0]
		assert.Equal(t, "VIP writer", created.Name)
		assert.Equal(t, aiModels.UserRoleWriter, created.UserRole)
		assert.Equal(t, aiModels.MembershipLevelVipMonthly, created.MembershipLevel)
		assert.Equal(t, 321, created.DailyQuota)
		assert.Equal(t, 6543, created.MonthlyQuota)
		assert.Equal(t, -1, created.TotalQuota)
		assert.Equal(t, "custom", created.Description)
	})

	t.Run("rejects malformed payload", func(t *testing.T) {
		router := setupQuotaPolicyAPITestRouter(newQuotaPolicyAPITestRepo())

		req, _ := http.NewRequest(http.MethodPost, "/policies", strings.NewReader("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "参数错误")
	})

	t.Run("surfaces create failure", func(t *testing.T) {
		repo := newQuotaPolicyAPITestRepo()
		repo.createErr = errors.New("create failed")
		router := setupQuotaPolicyAPITestRouter(repo)

		req, _ := http.NewRequest(http.MethodPost, "/policies", strings.NewReader(`{"name":"reader","userRole":"reader","membershipLevel":"normal"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "创建策略失败")
	})
}

func TestQuotaPolicyAPIGetUpdateDeleteAndInitializeErrorBranches(t *testing.T) {
	t.Run("rejects missing id for get update and delete", func(t *testing.T) {
		api := NewQuotaPolicyAPI(aiService.NewQuotaPolicyService(newQuotaPolicyAPITestRepo()))

		cases := []struct {
			name   string
			call   func(*gin.Context)
		}{
			{name: "get", call: api.GetPolicy},
			{name: "update", call: api.UpdatePolicy},
			{name: "delete", call: api.DeletePolicy},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				switch tt.name {
				case "update":
					c.Request = httptest.NewRequest(http.MethodPut, "/policies", strings.NewReader(`{"name":"x","userRole":"reader","membershipLevel":"normal"}`))
					c.Request.Header.Set("Content-Type", "application/json")
				case "delete":
					c.Request = httptest.NewRequest(http.MethodDelete, "/policies", nil)
				default:
					c.Request = httptest.NewRequest(http.MethodGet, "/policies", nil)
				}
				tt.call(c)

				require.Equal(t, http.StatusBadRequest, w.Code)
				assert.Contains(t, w.Body.String(), "策略ID不能为空")
			})
		}
	})

	t.Run("surfaces lookup failures", func(t *testing.T) {
		repo := newQuotaPolicyAPITestRepo()
		repo.getByIDErr = errors.New("lookup failed")
		router := setupQuotaPolicyAPITestRouter(repo)

		req, _ := http.NewRequest(http.MethodGet, "/policies/507f1f77bcf86cd799439011", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "获取策略失败")
	})

	t.Run("surfaces initialize failure when policy listing fails", func(t *testing.T) {
		previousConfig := config.GlobalConfig
		config.GlobalConfig = &config.Config{
			AIQuota: &config.AIQuotaConfig{
				DefaultQuotas: &config.DefaultQuotasConfig{
					Reader: map[string]int{"normal": 120},
					Writer: map[string]int{"normal": 240},
					Admin:  map[string]int{"normal": 360},
				},
			},
		}
		t.Cleanup(func() {
			config.GlobalConfig = previousConfig
		})

		repo := newQuotaPolicyAPITestRepo()
		repo.listErr = errors.New("list failed")
		router := setupQuotaPolicyAPITestRouter(repo)

		req, _ := http.NewRequest(http.MethodPost, "/policies/initialize", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "初始化默认策略失败")
	})
}

func TestQuotaPolicyAPIListPoliciesSurfacesRepositoryFailure(t *testing.T) {
	repo := newQuotaPolicyAPITestRepo()
	repo.listErr = errors.New("list failed")
	router := setupQuotaPolicyAPITestRouter(repo)

	req, _ := http.NewRequest(http.MethodGet, "/policies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "获取策略列表失败")
}

func TestQuotaPolicyAPIUpdatePolicyRejectsMalformedJSON(t *testing.T) {
	repo := newQuotaPolicyAPITestRepo()
	policyID := primitive.NewObjectID()
	repo.policiesByID[policyID.Hex()] = &aiModels.QuotaPolicy{
		ID:              policyID,
		Name:            "reader-normal",
		UserRole:        aiModels.UserRoleReader,
		MembershipLevel: aiModels.MembershipLevelNormal,
	}
	router := setupQuotaPolicyAPITestRouter(repo)

	req, _ := http.NewRequest(http.MethodPut, "/policies/"+policyID.Hex(), strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "参数错误")
}

func TestQuotaPolicyAPIUpdatePolicySurfacesGetFailure(t *testing.T) {
	repo := newQuotaPolicyAPITestRepo()
	policyID := primitive.NewObjectID()
	repo.getByIDErr = errors.New("get failed")
	router := setupQuotaPolicyAPITestRouter(repo)

	reqBody := `{"name":"reader-updated","userRole":"reader","membershipLevel":"normal","dailyQuota":222,"monthlyQuota":6666,"totalQuota":-1,"description":"updated"}`
	req, err := http.NewRequest(http.MethodPut, "/policies/"+policyID.Hex(), strings.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "获取策略失败")
	assert.Contains(t, w.Body.String(), "get failed")
}

func TestQuotaPolicyAPIUpdatePolicySurfacesUpdateFailure(t *testing.T) {
	repo := newQuotaPolicyAPITestRepo()
	policyID := primitive.NewObjectID()
	repo.policiesByID[policyID.Hex()] = &aiModels.QuotaPolicy{
		ID:              policyID,
		Name:            "reader-normal",
		UserRole:        aiModels.UserRoleReader,
		MembershipLevel: aiModels.MembershipLevelNormal,
	}
	repo.updateErr = errors.New("update failed")
	router := setupQuotaPolicyAPITestRouter(repo)

	reqBody := `{"name":"reader-updated","userRole":"reader","membershipLevel":"normal","dailyQuota":222,"monthlyQuota":6666,"totalQuota":-1,"description":"updated"}`
	req, err := http.NewRequest(http.MethodPut, "/policies/"+policyID.Hex(), strings.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "更新策略失败")
	assert.Contains(t, w.Body.String(), "update failed")
}

func TestQuotaPolicyAPIDeletePolicySurfacesDeleteFailure(t *testing.T) {
	repo := newQuotaPolicyAPITestRepo()
	policyID := primitive.NewObjectID()
	repo.policiesByID[policyID.Hex()] = &aiModels.QuotaPolicy{
		ID:        policyID,
		IsDefault: false,
	}
	repo.deleteErr = errors.New("delete failed")
	router := setupQuotaPolicyAPITestRouter(repo)

	req, _ := http.NewRequest(http.MethodDelete, "/policies/"+policyID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "删除策略失败")
	assert.Contains(t, w.Body.String(), "delete failed")
}

func TestQuotaPolicyAPIUpdateAndDeleteSuccess(t *testing.T) {
	t.Run("updates existing policy", func(t *testing.T) {
		repo := newQuotaPolicyAPITestRepo()
		policy := &aiModels.QuotaPolicy{
			ID:              primitive.NewObjectID(),
			Name:            "reader-normal",
			UserRole:        aiModels.UserRoleReader,
			MembershipLevel: aiModels.MembershipLevelNormal,
			DailyQuota:      100,
			MonthlyQuota:    3000,
		}
		repo.policiesByID[policy.ID.Hex()] = policy
		repo.policiesByRoleAndLevel[quotaPolicyAPIRepoKey(policy.UserRole, policy.MembershipLevel)] = policy
		router := setupQuotaPolicyAPITestRouter(repo)

		reqBody := `{"name":"reader-updated","userRole":"reader","membershipLevel":"normal","dailyQuota":222,"monthlyQuota":6666,"totalQuota":-1,"description":"updated"}`
		req, _ := http.NewRequest(http.MethodPut, "/policies/"+policy.ID.Hex(), strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Len(t, repo.updatedPolicies, 1)
		updated := repo.updatedPolicies[0]
		assert.Equal(t, "reader-updated", updated.Name)
		assert.Equal(t, 222, updated.DailyQuota)
		assert.Equal(t, 6666, updated.MonthlyQuota)
		assert.Equal(t, -1, updated.TotalQuota)
		assert.Equal(t, "updated", updated.Description)
	})

	t.Run("deletes non default policy", func(t *testing.T) {
		repo := newQuotaPolicyAPITestRepo()
		policy := &aiModels.QuotaPolicy{
			ID:        primitive.NewObjectID(),
			IsDefault: false,
		}
		repo.policiesByID[policy.ID.Hex()] = policy
		router := setupQuotaPolicyAPITestRouter(repo)

		req, _ := http.NewRequest(http.MethodDelete, "/policies/"+policy.ID.Hex(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, []string{policy.ID.Hex()}, repo.deletedIDs)
	})
}
