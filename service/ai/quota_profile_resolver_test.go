package ai

import (
	"context"
	"errors"
	"testing"
	"time"

	aiModels "Qingyu_backend/models/ai"
	financeModel "Qingyu_backend/models/finance"
	usersModel "Qingyu_backend/models/users"
	aiRepo "Qingyu_backend/repository/interfaces/ai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type quotaRepoStub struct {
	quotas map[string]*aiModels.UserQuota
}

func newQuotaRepoStub() *quotaRepoStub {
	return &quotaRepoStub{quotas: make(map[string]*aiModels.UserQuota)}
}

func (s *quotaRepoStub) key(userID string, quotaType aiModels.QuotaType) string {
	return userID + ":" + string(quotaType)
}

func (s *quotaRepoStub) CreateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	s.quotas[s.key(quota.UserID, quota.QuotaType)] = quota
	return nil
}

func (s *quotaRepoStub) GetQuotaByUserID(ctx context.Context, userID string, quotaType aiModels.QuotaType) (*aiModels.UserQuota, error) {
	quota, ok := s.quotas[s.key(userID, quotaType)]
	if !ok {
		return nil, aiModels.ErrQuotaNotFound
	}
	return quota, nil
}

func (s *quotaRepoStub) UpdateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	s.quotas[s.key(quota.UserID, quota.QuotaType)] = quota
	return nil
}

func (s *quotaRepoStub) DeleteQuota(ctx context.Context, userID string, quotaType aiModels.QuotaType) error {
	delete(s.quotas, s.key(userID, quotaType))
	return nil
}

func (s *quotaRepoStub) GetAllQuotasByUserID(ctx context.Context, userID string) ([]*aiModels.UserQuota, error) {
	var quotas []*aiModels.UserQuota
	for key, quota := range s.quotas {
		if len(key) >= len(userID)+1 && key[:len(userID)+1] == userID+":" {
			quotas = append(quotas, quota)
		}
	}
	return quotas, nil
}

func (s *quotaRepoStub) BatchResetQuotas(ctx context.Context, quotaType aiModels.QuotaType) error {
	return nil
}

func (s *quotaRepoStub) CreateTransaction(ctx context.Context, transaction *aiModels.QuotaTransaction) error {
	return nil
}

func (s *quotaRepoStub) GetTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*aiModels.QuotaTransaction, error) {
	return nil, nil
}

func (s *quotaRepoStub) GetTransactionsByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*aiModels.QuotaTransaction, error) {
	return nil, nil
}

func (s *quotaRepoStub) GetQuotaStatistics(ctx context.Context, userID string) (*aiRepo.QuotaStatistics, error) {
	return nil, nil
}

func (s *quotaRepoStub) GetTotalConsumption(ctx context.Context, userID string, quotaType aiModels.QuotaType, startTime, endTime time.Time) (int, error) {
	return 0, nil
}

func (s *quotaRepoStub) GetDashboardSummary(ctx context.Context) (*aiModels.DashboardSummary, error) {
	return nil, nil
}

func (s *quotaRepoStub) GetQuotaDistribution(ctx context.Context) (*aiModels.QuotaDistribution, error) {
	return nil, nil
}

func (s *quotaRepoStub) GetTopConsumers(ctx context.Context, limit int) ([]aiModels.UserQuotaRanking, error) {
	return nil, nil
}

func (s *quotaRepoStub) GetConsumptionTrend(ctx context.Context, days int) ([]aiModels.TrendPoint, error) {
	return nil, nil
}

func (s *quotaRepoStub) GetConsumptionSummary(ctx context.Context, startTime, endTime time.Time, workflowType, groupBy string, page, pageSize int) (*aiModels.QuotaConsumptionSummary, error) {
	return nil, nil
}

func (s *quotaRepoStub) ListUserQuotas(ctx context.Context, role, status, search string, page, limit int) ([]*aiModels.UserQuotaListItem, int64, error) {
	return nil, 0, nil
}

func (s *quotaRepoStub) Health(ctx context.Context) error {
	return nil
}

type quotaUserReaderStub struct {
	user *usersModel.User
	err  error
}

func (s *quotaUserReaderStub) GetByID(ctx context.Context, id string) (*usersModel.User, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.user, nil
}

type quotaMembershipReaderStub struct {
	membership *financeModel.UserMembership
	err        error
}

func (s *quotaMembershipReaderStub) GetMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.membership, nil
}

func TestQuotaProfileResolverResolveAuthorVIPProfile(t *testing.T) {
	resolver := NewQuotaProfileResolver(
		&quotaUserReaderStub{user: &usersModel.User{Roles: []string{"author"}}},
		&quotaMembershipReaderStub{membership: &financeModel.UserMembership{
			UserID:    "u-1",
			Level:     financeModel.MembershipLevelVIPYearly,
			Status:    financeModel.MembershipStatusActive,
			EndTime:   time.Now().Add(24 * time.Hour),
			StartTime: time.Now().Add(-24 * time.Hour),
		}},
	)

	profile, err := resolver.Resolve(context.Background(), "u-1")
	require.NoError(t, err)
	assert.Equal(t, "writer", profile.UserRole)
	assert.Equal(t, financeModel.MembershipLevelVIPYearly, profile.MembershipLevel)
}

func TestQuotaServiceCheckQuotaInitializesResolvedVIPReaderQuota(t *testing.T) {
	repo := newQuotaRepoStub()
	service := NewQuotaService(repo)
	service.SetProfileResolver(NewQuotaProfileResolver(
		&quotaUserReaderStub{user: &usersModel.User{Roles: []string{"reader"}}},
		&quotaMembershipReaderStub{membership: &financeModel.UserMembership{
			UserID:    "reader-1",
			Level:     financeModel.MembershipLevelVIPMonthly,
			Status:    financeModel.MembershipStatusActive,
			EndTime:   time.Now().Add(24 * time.Hour),
			StartTime: time.Now().Add(-24 * time.Hour),
		}},
	))

	err := service.CheckQuota(context.Background(), "reader-1", 1)
	require.NoError(t, err)

	quota, err := repo.GetQuotaByUserID(context.Background(), "reader-1", aiModels.QuotaTypeDaily)
	require.NoError(t, err)
	assert.Equal(t, 50, quota.TotalQuota)
	require.NotNil(t, quota.Metadata)
	assert.Equal(t, "reader", quota.Metadata.UserRole)
	assert.Equal(t, financeModel.MembershipLevelVIPMonthly, quota.Metadata.MembershipLevel)
}

func TestQuotaServiceRefreshUserQuotaProfilePreservesManualOverride(t *testing.T) {
	repo := newQuotaRepoStub()
	quota := &aiModels.UserQuota{
		UserID:         "reader-2",
		QuotaType:      aiModels.QuotaTypeDaily,
		TotalQuota:     999,
		UsedQuota:      10,
		RemainingQuota: 989,
		Status:         aiModels.QuotaStatusActive,
		ResetAt:        time.Now().Add(24 * time.Hour),
		Metadata: &aiModels.QuotaMetadata{
			UserRole:        "reader",
			MembershipLevel: "normal",
		},
	}
	setQuotaManualOverride(quota, true)
	require.NoError(t, repo.CreateQuota(context.Background(), quota))

	service := NewQuotaService(repo)
	service.SetProfileResolver(NewQuotaProfileResolver(
		&quotaUserReaderStub{user: &usersModel.User{Roles: []string{"reader"}}},
		&quotaMembershipReaderStub{membership: &financeModel.UserMembership{
			UserID:    "reader-2",
			Level:     financeModel.MembershipLevelVIPYearly,
			Status:    financeModel.MembershipStatusActive,
			EndTime:   time.Now().Add(24 * time.Hour),
			StartTime: time.Now().Add(-24 * time.Hour),
		}},
	))

	err := service.RefreshUserQuotaProfile(context.Background(), "reader-2")
	require.NoError(t, err)

	refreshed, err := repo.GetQuotaByUserID(context.Background(), "reader-2", aiModels.QuotaTypeDaily)
	require.NoError(t, err)
	assert.Equal(t, 999, refreshed.TotalQuota)
	assert.Equal(t, 989, refreshed.RemainingQuota)
	assert.Equal(t, financeModel.MembershipLevelVIPYearly, refreshed.Metadata.MembershipLevel)
}

func TestQuotaProfileResolverFallsBackToDefaults(t *testing.T) {
	resolver := NewQuotaProfileResolver(
		&quotaUserReaderStub{err: errors.New("user unavailable")},
		&quotaMembershipReaderStub{err: errors.New("membership unavailable")},
	)

	profile, err := resolver.Resolve(context.Background(), "u-2")
	require.NoError(t, err)
	assert.Equal(t, "reader", profile.UserRole)
	assert.Equal(t, financeModel.MembershipLevelNormal, profile.MembershipLevel)
}
