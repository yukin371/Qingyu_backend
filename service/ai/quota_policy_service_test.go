package ai

import (
	"context"
	"errors"
	"testing"

	"Qingyu_backend/config"
	aiModels "Qingyu_backend/models/ai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type quotaPolicyRepoStub struct {
	listPolicies           []*aiModels.QuotaPolicy
	createPolicies         []*aiModels.QuotaPolicy
	updatePolicies         []*aiModels.QuotaPolicy
	deletedIDs             []string
	policiesByID           map[string]*aiModels.QuotaPolicy
	policiesByRoleAndLevel map[string]*aiModels.QuotaPolicy
	listErr                error
	createErr              error
	getByIDErr             error
	getByRoleAndLevelErr   error
	updateErr              error
	deleteErr              error
}

func newQuotaPolicyRepoStub() *quotaPolicyRepoStub {
	return &quotaPolicyRepoStub{
		policiesByID:           make(map[string]*aiModels.QuotaPolicy),
		policiesByRoleAndLevel: make(map[string]*aiModels.QuotaPolicy),
	}
}

func policyRepoKey(role aiModels.UserRole, level aiModels.MembershipLevel) string {
	return string(role) + "|" + string(level)
}

func (s *quotaPolicyRepoStub) Create(ctx context.Context, policy *aiModels.QuotaPolicy) error {
	if s.createErr != nil {
		return s.createErr
	}
	if policy.ID.IsZero() {
		policy.ID = primitive.NewObjectID()
	}
	s.createPolicies = append(s.createPolicies, policy)
	s.policiesByID[policy.ID.Hex()] = policy
	s.policiesByRoleAndLevel[policyRepoKey(policy.UserRole, policy.MembershipLevel)] = policy
	return nil
}

func (s *quotaPolicyRepoStub) GetByID(ctx context.Context, id string) (*aiModels.QuotaPolicy, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	if policy, ok := s.policiesByID[id]; ok {
		return policy, nil
	}
	return nil, errors.New("not found")
}

func (s *quotaPolicyRepoStub) GetByRoleAndLevel(ctx context.Context, role aiModels.UserRole, level aiModels.MembershipLevel) (*aiModels.QuotaPolicy, error) {
	if s.getByRoleAndLevelErr != nil {
		return nil, s.getByRoleAndLevelErr
	}
	if policy, ok := s.policiesByRoleAndLevel[policyRepoKey(role, level)]; ok {
		return policy, nil
	}
	return nil, errors.New("not found")
}

func (s *quotaPolicyRepoStub) List(ctx context.Context, role string, status string, page, limit int) ([]*aiModels.QuotaPolicy, int64, error) {
	if s.listErr != nil {
		return nil, 0, s.listErr
	}
	return s.listPolicies, int64(len(s.listPolicies)), nil
}

func (s *quotaPolicyRepoStub) Update(ctx context.Context, policy *aiModels.QuotaPolicy) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.updatePolicies = append(s.updatePolicies, policy)
	s.policiesByID[policy.ID.Hex()] = policy
	s.policiesByRoleAndLevel[policyRepoKey(policy.UserRole, policy.MembershipLevel)] = policy
	return nil
}

func (s *quotaPolicyRepoStub) Delete(ctx context.Context, id string) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}
	s.deletedIDs = append(s.deletedIDs, id)
	delete(s.policiesByID, id)
	return nil
}

func (s *quotaPolicyRepoStub) Health(ctx context.Context) error {
	return nil
}

func TestQuotaPolicyServiceCreatePolicyAppliesDefaultsAndDetectsDuplicates(t *testing.T) {
	t.Run("rejects duplicate policy", func(t *testing.T) {
		repo := newQuotaPolicyRepoStub()
		existing := &aiModels.QuotaPolicy{
			ID:              primitive.NewObjectID(),
			UserRole:        aiModels.UserRoleReader,
			MembershipLevel: aiModels.MembershipLevelNormal,
			DailyQuota:      80,
		}
		repo.policiesByID[existing.ID.Hex()] = existing
		repo.policiesByRoleAndLevel[policyRepoKey(existing.UserRole, existing.MembershipLevel)] = existing

		service := NewQuotaPolicyService(repo)
		err := service.CreatePolicy(context.Background(), &aiModels.QuotaPolicy{
			UserRole:        aiModels.UserRoleReader,
			MembershipLevel: aiModels.MembershipLevelNormal,
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "已存在角色")
		assert.Len(t, repo.createPolicies, 0)
	})

	t.Run("fills default quotas for new policy", func(t *testing.T) {
		repo := newQuotaPolicyRepoStub()
		service := NewQuotaPolicyService(repo)
		policy := &aiModels.QuotaPolicy{
			UserRole:        aiModels.UserRoleWriter,
			MembershipLevel: aiModels.MembershipLevelVipMonthly,
			DailyQuota:      0,
			MonthlyQuota:    0,
			TotalQuota:      -5,
		}

		err := service.CreatePolicy(context.Background(), policy)
		require.NoError(t, err)
		require.Len(t, repo.createPolicies, 1)

		created := repo.createPolicies[0]
		assert.Equal(t, 1000, created.DailyQuota)
		assert.Equal(t, 30000, created.MonthlyQuota)
		assert.Equal(t, -1, created.TotalQuota)
	})
}

func TestQuotaPolicyServiceUpdatePolicyValidatesExistenceConflictsAndQuotaDefaults(t *testing.T) {
	t.Run("rejects nil policy", func(t *testing.T) {
		service := NewQuotaPolicyService(newQuotaPolicyRepoStub())

		err := service.UpdatePolicy(context.Background(), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "策略不能为空")
	})

	t.Run("rejects missing policy", func(t *testing.T) {
		service := NewQuotaPolicyService(newQuotaPolicyRepoStub())

		err := service.UpdatePolicy(context.Background(), &aiModels.QuotaPolicy{ID: primitive.NewObjectID()})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "策略不存在")
	})

	t.Run("rejects duplicate role and membership combination", func(t *testing.T) {
		repo := newQuotaPolicyRepoStub()
		existing := &aiModels.QuotaPolicy{
			ID:              primitive.NewObjectID(),
			UserRole:        aiModels.UserRoleReader,
			MembershipLevel: aiModels.MembershipLevelNormal,
		}
		conflict := &aiModels.QuotaPolicy{
			ID:              primitive.NewObjectID(),
			UserRole:        aiModels.UserRoleWriter,
			MembershipLevel: aiModels.MembershipLevelVipMonthly,
		}
		repo.policiesByID[existing.ID.Hex()] = existing
		repo.policiesByID[conflict.ID.Hex()] = conflict
		repo.policiesByRoleAndLevel[policyRepoKey(existing.UserRole, existing.MembershipLevel)] = existing
		repo.policiesByRoleAndLevel[policyRepoKey(conflict.UserRole, conflict.MembershipLevel)] = conflict

		service := NewQuotaPolicyService(repo)
		err := service.UpdatePolicy(context.Background(), &aiModels.QuotaPolicy{
			ID:              existing.ID,
			UserRole:        conflict.UserRole,
			MembershipLevel: conflict.MembershipLevel,
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "已存在角色")
		assert.Empty(t, repo.updatePolicies)
	})

	t.Run("fills default daily and monthly quota on update", func(t *testing.T) {
		repo := newQuotaPolicyRepoStub()
		existing := &aiModels.QuotaPolicy{
			ID:              primitive.NewObjectID(),
			UserRole:        aiModels.UserRoleWriter,
			MembershipLevel: aiModels.MembershipLevelVipMonthly,
			DailyQuota:      88,
			MonthlyQuota:    1234,
			TotalQuota:      5000,
		}
		repo.policiesByID[existing.ID.Hex()] = existing
		repo.policiesByRoleAndLevel[policyRepoKey(existing.UserRole, existing.MembershipLevel)] = existing

		service := NewQuotaPolicyService(repo)
		err := service.UpdatePolicy(context.Background(), &aiModels.QuotaPolicy{
			ID:              existing.ID,
			UserRole:        existing.UserRole,
			MembershipLevel: existing.MembershipLevel,
			DailyQuota:      0,
			MonthlyQuota:    -1,
			TotalQuota:      0,
		})

		require.NoError(t, err)
		require.Len(t, repo.updatePolicies, 1)
		updated := repo.updatePolicies[0]
		assert.Equal(t, 1000, updated.DailyQuota)
		assert.Equal(t, 30000, updated.MonthlyQuota)
		assert.Equal(t, 0, updated.TotalQuota)
	})

	t.Run("rejects unsupported negative total quota", func(t *testing.T) {
		repo := newQuotaPolicyRepoStub()
		existing := &aiModels.QuotaPolicy{
			ID:              primitive.NewObjectID(),
			UserRole:        aiModels.UserRoleAdmin,
			MembershipLevel: aiModels.MembershipLevelNormal,
		}
		repo.policiesByID[existing.ID.Hex()] = existing
		repo.policiesByRoleAndLevel[policyRepoKey(existing.UserRole, existing.MembershipLevel)] = existing

		service := NewQuotaPolicyService(repo)
		err := service.UpdatePolicy(context.Background(), &aiModels.QuotaPolicy{
			ID:              existing.ID,
			UserRole:        existing.UserRole,
			MembershipLevel: existing.MembershipLevel,
			DailyQuota:      10,
			MonthlyQuota:    100,
			TotalQuota:      -2,
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "总配额不能为负数")
		assert.Empty(t, repo.updatePolicies)
	})
}

func TestQuotaPolicyServiceInitializeDefaultPoliciesCreatesAllRoleLevelCombinations(t *testing.T) {
	previousConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{
		AIQuota: &config.AIQuotaConfig{
			DefaultQuotas: &config.DefaultQuotasConfig{
				Reader: map[string]int{
					"normal": 120,
					"vip":    240,
				},
				Writer: map[string]int{
					"novice": 360,
				},
				Admin: map[string]int{
					"normal": 480,
				},
			},
		},
	}
	t.Cleanup(func() {
		config.GlobalConfig = previousConfig
	})

	repo := newQuotaPolicyRepoStub()
	service := NewQuotaPolicyService(repo)

	err := service.InitializeDefaultPolicies(context.Background())
	require.NoError(t, err)
	require.Len(t, repo.createPolicies, 12)

	policies := make(map[string]*aiModels.QuotaPolicy, len(repo.createPolicies))
	for _, policy := range repo.createPolicies {
		require.NotNil(t, policy)
		policies[policyRepoKey(policy.UserRole, policy.MembershipLevel)] = policy
		assert.True(t, policy.IsDefault)
		assert.Equal(t, aiModels.QuotaPolicyStatusActive, policy.Status)
		assert.Equal(t, -1, policy.TotalQuota)
		assert.Equal(t, policy.DailyQuota*30, policy.MonthlyQuota)
	}

	assert.Equal(t, 120, policies[policyRepoKey(aiModels.UserRoleReader, aiModels.MembershipLevelNormal)].DailyQuota)
	assert.Equal(t, 240, policies[policyRepoKey(aiModels.UserRoleReader, aiModels.MembershipLevelVipYearly)].DailyQuota)
	assert.Equal(t, 360, policies[policyRepoKey(aiModels.UserRoleWriter, aiModels.MembershipLevelSuperVip)].DailyQuota)
	assert.Equal(t, 480, policies[policyRepoKey(aiModels.UserRoleAdmin, aiModels.MembershipLevelNormal)].DailyQuota)
}

func TestQuotaPolicyServiceInitializeDefaultPoliciesSkipsExistingPolicies(t *testing.T) {
	previousConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{
		AIQuota: &config.AIQuotaConfig{
			DefaultQuotas: &config.DefaultQuotasConfig{
				Reader: map[string]int{"normal": 120},
				Writer: map[string]int{"novice": 360},
				Admin:  map[string]int{"normal": 480},
			},
		},
	}
	t.Cleanup(func() {
		config.GlobalConfig = previousConfig
	})

	repo := newQuotaPolicyRepoStub()
	repo.listPolicies = []*aiModels.QuotaPolicy{
		{
			ID:              primitive.NewObjectID(),
			UserRole:        aiModels.UserRoleReader,
			MembershipLevel: aiModels.MembershipLevelNormal,
		},
	}

	service := NewQuotaPolicyService(repo)
	err := service.InitializeDefaultPolicies(context.Background())
	require.NoError(t, err)
	assert.Empty(t, repo.createPolicies)
}

func TestQuotaPolicyServiceGetEffectiveDailyQuotaUsesActivePolicyConfigAndModelFallbacks(t *testing.T) {
	t.Run("returns active policy quota", func(t *testing.T) {
		repo := newQuotaPolicyRepoStub()
		repo.policiesByRoleAndLevel[policyRepoKey(aiModels.UserRoleReader, aiModels.MembershipLevelNormal)] = &aiModels.QuotaPolicy{
			UserRole:        aiModels.UserRoleReader,
			MembershipLevel: aiModels.MembershipLevelNormal,
			DailyQuota:      777,
			Status:          aiModels.QuotaPolicyStatusActive,
		}
		service := NewQuotaPolicyService(repo)

		got, err := service.GetEffectiveDailyQuota(context.Background(), "reader", "normal")
		require.NoError(t, err)
		assert.Equal(t, 777, got)
	})

	t.Run("falls back to config when policy is absent or disabled", func(t *testing.T) {
		previousConfig := config.GlobalConfig
		config.GlobalConfig = &config.Config{
			AIQuota: &config.AIQuotaConfig{
				DefaultQuotas: &config.DefaultQuotasConfig{
					Reader: map[string]int{
						"normal": 120,
					},
					Writer: map[string]int{
						"novice": 360,
					},
					Admin: map[string]int{
						"normal": 480,
					},
				},
			},
		}
		t.Cleanup(func() {
			config.GlobalConfig = previousConfig
		})

		repo := newQuotaPolicyRepoStub()
		repo.policiesByRoleAndLevel[policyRepoKey(aiModels.UserRoleWriter, aiModels.MembershipLevelVipMonthly)] = &aiModels.QuotaPolicy{
			UserRole:        aiModels.UserRoleWriter,
			MembershipLevel: aiModels.MembershipLevelVipMonthly,
			DailyQuota:      999,
			Status:          aiModels.QuotaPolicyStatusDisabled,
		}
		service := NewQuotaPolicyService(repo)

		got, err := service.GetEffectiveDailyQuota(context.Background(), "writer", "vip_yearly")
		require.NoError(t, err)
		assert.Equal(t, 360, got)
	})

	t.Run("falls back to model default when config is missing", func(t *testing.T) {
		previousConfig := config.GlobalConfig
		config.GlobalConfig = nil
		t.Cleanup(func() {
			config.GlobalConfig = previousConfig
		})

		service := NewQuotaPolicyService(newQuotaPolicyRepoStub())

		got, err := service.GetEffectiveDailyQuota(context.Background(), "admin", "normal")
		require.NoError(t, err)
		assert.Equal(t, 5, got)
	})
}

func TestQuotaPolicyServiceRejectsInvalidInputsAndProtectedDefaults(t *testing.T) {
	t.Run("rejects invalid user role", func(t *testing.T) {
		service := NewQuotaPolicyService(newQuotaPolicyRepoStub())

		got, err := service.GetEffectiveDailyQuota(context.Background(), "ghost", "normal")
		require.Error(t, err)
		assert.Equal(t, 0, got)
		assert.Contains(t, err.Error(), "无效的用户角色")
	})

	t.Run("rejects invalid membership level", func(t *testing.T) {
		service := NewQuotaPolicyService(newQuotaPolicyRepoStub())

		got, err := service.GetEffectiveDailyQuota(context.Background(), "reader", "quarter")
		require.Error(t, err)
		assert.Equal(t, 0, got)
		assert.Contains(t, err.Error(), "无效的会员等级")
	})

	t.Run("rejects initialization when config is missing", func(t *testing.T) {
		previousConfig := config.GlobalConfig
		config.GlobalConfig = nil
		t.Cleanup(func() {
			config.GlobalConfig = previousConfig
		})

		service := NewQuotaPolicyService(newQuotaPolicyRepoStub())

		err := service.InitializeDefaultPolicies(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "配置未加载")
	})

	t.Run("rejects deleting default policy", func(t *testing.T) {
		repo := newQuotaPolicyRepoStub()
		defaultPolicy := &aiModels.QuotaPolicy{
			ID:        primitive.NewObjectID(),
			IsDefault: true,
		}
		repo.policiesByID[defaultPolicy.ID.Hex()] = defaultPolicy

		service := NewQuotaPolicyService(repo)

		err := service.DeletePolicy(context.Background(), defaultPolicy.ID.Hex())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "不能删除默认策略")
	})

	t.Run("deletes non-default policy", func(t *testing.T) {
		repo := newQuotaPolicyRepoStub()
		policy := &aiModels.QuotaPolicy{
			ID:        primitive.NewObjectID(),
			IsDefault: false,
		}
		repo.policiesByID[policy.ID.Hex()] = policy

		service := NewQuotaPolicyService(repo)

		err := service.DeletePolicy(context.Background(), policy.ID.Hex())
		require.NoError(t, err)
		assert.Equal(t, []string{policy.ID.Hex()}, repo.deletedIDs)
	})
}
