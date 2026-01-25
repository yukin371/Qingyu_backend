package finance

import (
	"context"
	"testing"
	"time"

	financeModel "Qingyu_backend/models/finance"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =========================
// Mock Repository实现
// =========================

// MockMembershipRepository Mock会员仓储
type MockMembershipRepository struct {
	mock.Mock
}

// 套餐相关方法
func (m *MockMembershipRepository) ListPlans(ctx context.Context, isEnabled bool) ([]*financeModel.MembershipPlan, error) {
	args := m.Called(ctx, isEnabled)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*financeModel.MembershipPlan), args.Error(1)
}

func (m *MockMembershipRepository) GetPlan(ctx context.Context, id primitive.ObjectID) (*financeModel.MembershipPlan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*financeModel.MembershipPlan), args.Error(1)
}

func (m *MockMembershipRepository) GetPlanByType(ctx context.Context, planType string) (*financeModel.MembershipPlan, error) {
	args := m.Called(ctx, planType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*financeModel.MembershipPlan), args.Error(1)
}

func (m *MockMembershipRepository) CreatePlan(ctx context.Context, plan *financeModel.MembershipPlan) error {
	args := m.Called(ctx, plan)
	return args.Error(0)
}

func (m *MockMembershipRepository) UpdatePlan(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockMembershipRepository) DeletePlan(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// 会员相关方法
func (m *MockMembershipRepository) GetMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*financeModel.UserMembership), args.Error(1)
}

func (m *MockMembershipRepository) GetMembershipByID(ctx context.Context, id primitive.ObjectID) (*financeModel.UserMembership, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*financeModel.UserMembership), args.Error(1)
}

func (m *MockMembershipRepository) CreateMembership(ctx context.Context, membership *financeModel.UserMembership) error {
	args := m.Called(ctx, membership)
	return args.Error(0)
}

func (m *MockMembershipRepository) UpdateMembership(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockMembershipRepository) DeleteMembership(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMembershipRepository) ListMemberships(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.UserMembership, error) {
	args := m.Called(ctx, filter, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*financeModel.UserMembership), args.Error(1)
}

func (m *MockMembershipRepository) CountMemberships(ctx context.Context, filter map[string]interface{}) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// 权益相关方法
func (m *MockMembershipRepository) ListBenefits(ctx context.Context, level string, isEnabled bool) ([]*financeModel.MembershipBenefit, error) {
	args := m.Called(ctx, level, isEnabled)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*financeModel.MembershipBenefit), args.Error(1)
}

func (m *MockMembershipRepository) GetBenefit(ctx context.Context, id primitive.ObjectID) (*financeModel.MembershipBenefit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*financeModel.MembershipBenefit), args.Error(1)
}

func (m *MockMembershipRepository) GetBenefitByCode(ctx context.Context, code string) (*financeModel.MembershipBenefit, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*financeModel.MembershipBenefit), args.Error(1)
}

func (m *MockMembershipRepository) CreateBenefit(ctx context.Context, benefit *financeModel.MembershipBenefit) error {
	args := m.Called(ctx, benefit)
	return args.Error(0)
}

func (m *MockMembershipRepository) UpdateBenefit(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockMembershipRepository) DeleteBenefit(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// 权益使用相关方法
func (m *MockMembershipRepository) ListUsages(ctx context.Context, userID string) ([]*financeModel.MembershipUsage, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*financeModel.MembershipUsage), args.Error(1)
}

func (m *MockMembershipRepository) CreateUsage(ctx context.Context, usage *financeModel.MembershipUsage) error {
	args := m.Called(ctx, usage)
	return args.Error(0)
}

func (m *MockMembershipRepository) UpdateUsage(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockMembershipRepository) GetUsage(ctx context.Context, userID string, benefitID string) (*financeModel.MembershipUsage, error) {
	args := m.Called(ctx, userID, benefitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*financeModel.MembershipUsage), args.Error(1)
}

// 会员卡相关方法
func (m *MockMembershipRepository) ListMembershipCards(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.MembershipCard, error) {
	args := m.Called(ctx, filter, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*financeModel.MembershipCard), args.Error(1)
}

func (m *MockMembershipRepository) CountMembershipCards(ctx context.Context, filter map[string]interface{}) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMembershipRepository) GetMembershipCard(ctx context.Context, id primitive.ObjectID) (*financeModel.MembershipCard, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*financeModel.MembershipCard), args.Error(1)
}

func (m *MockMembershipRepository) GetMembershipCardByCode(ctx context.Context, code string) (*financeModel.MembershipCard, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*financeModel.MembershipCard), args.Error(1)
}

func (m *MockMembershipRepository) CreateMembershipCard(ctx context.Context, card *financeModel.MembershipCard) error {
	args := m.Called(ctx, card)
	return args.Error(0)
}

func (m *MockMembershipRepository) UpdateMembershipCard(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockMembershipRepository) DeleteMembershipCard(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMembershipRepository) BatchCreateMembershipCards(ctx context.Context, cards []*financeModel.MembershipCard) error {
	args := m.Called(ctx, cards)
	return args.Error(0)
}

func (m *MockMembershipRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// =========================
// 测试辅助函数
// =========================

// setupMembershipService 创建测试用的MembershipService实例
func setupMembershipService() (*MembershipServiceImpl, *MockMembershipRepository) {
	mockRepo := new(MockMembershipRepository)

	service := NewMembershipService(mockRepo)

	return service.(*MembershipServiceImpl), mockRepo
}

// createTestPlan 创建测试用套餐
func createTestPlan(name, planType string, duration int, price float64, enabled bool) *financeModel.MembershipPlan {
	return &financeModel.MembershipPlan{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Type:      planType,
		Duration:  duration,
		Price:     types.NewMoneyFromYuan(price),
		IsEnabled: enabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// createTestMembership 创建测试用会员
func createTestMembership(userID string, planID primitive.ObjectID, level string, startTime, endTime time.Time) *financeModel.UserMembership {
	return &financeModel.UserMembership{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		PlanID:      planID,
		PlanName:    "测试套餐",
		PlanType:    financeModel.MembershipTypeMonthly,
		Level:       level,
		StartTime:   startTime,
		EndTime:     endTime,
		AutoRenew:   false,
		Status:      financeModel.MembershipStatusActive,
		ActivatedAt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// createTestMembershipCard 创建测试用会员卡
func createTestMembershipCard(code string, planType string, duration int, status string) *financeModel.MembershipCard {
	expiredAt := time.Now().AddDate(0, 0, 30)
	return &financeModel.MembershipCard{
		ID:        primitive.NewObjectID(),
		Code:      code,
		PlanType:  planType,
		PlanID:    primitive.NewObjectID(),
		Duration:  duration,
		Status:    status,
		ExpireAt:  &expiredAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// =========================
// 套餐管理相关测试
// =========================

// TestMembershipService_GetPlans_Success 测试获取套餐列表成功
func TestMembershipService_GetPlans_Success(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	expectedPlans := []*financeModel.MembershipPlan{
		createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true),
		createTestPlan("年度VIP", financeModel.MembershipTypeYearly, 365, 199.9, true),
	}

	mockRepo.On("ListPlans", ctx, true).Return(expectedPlans, nil)

	// Act
	plans, err := service.GetPlans(ctx)

	// Assert
	require.NoError(t, err)
	assert.Len(t, plans, 2)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_GetPlans_RepositoryError 测试获取套餐列表-仓储错误
func TestMembershipService_GetPlans_RepositoryError(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	mockRepo.On("ListPlans", ctx, true).Return(nil, assert.AnError)

	// Act
	plans, err := service.GetPlans(ctx)

	// Assert
	require.Error(t, err)
	assert.Nil(t, plans)
	assert.Contains(t, err.Error(), "获取套餐列表失败")

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_GetPlan_Success 测试获取套餐详情成功
func TestMembershipService_GetPlan_Success(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	planID := plan.ID.Hex()

	mockRepo.On("GetPlan", ctx, mock.AnythingOfType("primitive.ObjectID")).Return(plan, nil)

	// Act
	result, err := service.GetPlan(ctx, planID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, plan, result)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_GetPlan_InvalidID 测试获取套餐详情-无效ID
func TestMembershipService_GetPlan_InvalidID(t *testing.T) {
	// Arrange
	service, _ := setupMembershipService()
	ctx := context.Background()

	// Act
	result, err := service.GetPlan(ctx, "invalid-id")

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "无效的套餐ID")
}

// TestMembershipService_GetPlan_NotFound 测试获取套餐详情-不存在
func TestMembershipService_GetPlan_NotFound(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	planID := primitive.NewObjectID().Hex()

	mockRepo.On("GetPlan", ctx, mock.AnythingOfType("primitive.ObjectID")).Return(nil, assert.AnError)

	// Act
	result, err := service.GetPlan(ctx, planID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "获取套餐失败")

	mockRepo.AssertExpectations(t)
}

// =========================
// 订阅管理相关测试
// =========================

// TestMembershipService_Subscribe_Success 测试订阅会员成功
func TestMembershipService_Subscribe_Success(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	planID := plan.ID.Hex()

	mockRepo.On("GetPlan", ctx, mock.AnythingOfType("primitive.ObjectID")).Return(plan, nil)
	mockRepo.On("GetMembership", ctx, userID).Return(nil, assert.AnError)
	mockRepo.On("CreateMembership", ctx, mock.MatchedBy(func(m *financeModel.UserMembership) bool {
		return m.UserID == userID && m.Status == financeModel.MembershipStatusActive
	})).Return(nil)

	// Act
	membership, err := service.Subscribe(ctx, userID, planID, "alipay")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, membership)
	assert.Equal(t, userID, membership.UserID)
	assert.Equal(t, financeModel.MembershipStatusActive, membership.Status)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_Subscribe_PlanDisabled 测试订阅会员-套餐已下架
func TestMembershipService_Subscribe_PlanDisabled(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, false) // disabled
	planID := plan.ID.Hex()

	mockRepo.On("GetPlan", ctx, mock.AnythingOfType("primitive.ObjectID")).Return(plan, nil)

	// Act
	membership, err := service.Subscribe(ctx, userID, planID, "alipay")

	// Assert
	require.Error(t, err)
	assert.Nil(t, membership)
	assert.Contains(t, err.Error(), "套餐已下架")

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_Subscribe_AlreadyMember 测试订阅会员-已是会员
func TestMembershipService_Subscribe_AlreadyMember(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	planID := plan.ID.Hex()

	existingMembership := createTestMembership(userID, plan.ID, financeModel.MembershipLevelVIPMonthly, time.Now(), time.Now().Add(30*24*time.Hour))
	existingMembership.Status = financeModel.MembershipStatusActive

	mockRepo.On("GetPlan", ctx, mock.AnythingOfType("primitive.ObjectID")).Return(plan, nil)
	mockRepo.On("GetMembership", ctx, userID).Return(existingMembership, nil)

	// Act
	membership, err := service.Subscribe(ctx, userID, planID, "alipay")

	// Assert
	require.Error(t, err)
	assert.Nil(t, membership)
	assert.Contains(t, err.Error(), "您已是会员")

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_Subscribe_Renewal 测试订阅会员-续费
func TestMembershipService_Subscribe_Renewal(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	planID := plan.ID.Hex()

	oldEndTime := time.Now().Add(10 * 24 * time.Hour) // 10天后过期
	existingMembership := createTestMembership(userID, plan.ID, financeModel.MembershipLevelVIPMonthly, time.Now(), oldEndTime)
	existingMembership.Status = financeModel.MembershipStatusActive

	mockRepo.On("GetPlan", ctx, mock.AnythingOfType("primitive.ObjectID")).Return(plan, nil)
	mockRepo.On("GetMembership", ctx, userID).Return(existingMembership, nil)

	// Act
	membership, err := service.Subscribe(ctx, userID, planID, "alipay")

	// Assert - Service returns error for already active members
	require.Error(t, err)
	assert.Nil(t, membership)
	assert.Contains(t, err.Error(), "您已是会员")

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_GetMembership_Success 测试获取会员状态成功
func TestMembershipService_GetMembership_Success(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	expectedMembership := createTestMembership(userID, plan.ID, financeModel.MembershipLevelVIPMonthly, time.Now(), time.Now().Add(30*24*time.Hour))

	mockRepo.On("GetMembership", ctx, userID).Return(expectedMembership, nil)

	// Act
	membership, err := service.GetMembership(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedMembership, membership)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_GetMembership_Expired 测试获取会员状态-已过期
func TestMembershipService_GetMembership_Expired(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	expiredMembership := createTestMembership(userID, plan.ID, financeModel.MembershipLevelVIPMonthly, time.Now().Add(-60*24*time.Hour), time.Now().Add(-30*24*time.Hour))
	expiredMembership.Status = financeModel.MembershipStatusActive

	mockRepo.On("GetMembership", ctx, userID).Return(expiredMembership, nil)
	mockRepo.On("UpdateMembership", ctx, expiredMembership.ID, mock.MatchedBy(func(u map[string]interface{}) bool {
		return u["status"] == financeModel.MembershipStatusExpired
	})).Return(nil)

	// Act
	membership, err := service.GetMembership(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, financeModel.MembershipStatusExpired, membership.Status)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_CancelMembership_Success 测试取消自动续费成功
func TestMembershipService_CancelMembership_Success(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	membership := createTestMembership(userID, plan.ID, financeModel.MembershipLevelVIPMonthly, time.Now(), time.Now().Add(30*24*time.Hour))
	membership.AutoRenew = true
	membership.Status = financeModel.MembershipStatusActive

	mockRepo.On("GetMembership", ctx, userID).Return(membership, nil)
	mockRepo.On("UpdateMembership", ctx, membership.ID, mock.MatchedBy(func(u map[string]interface{}) bool {
		return u["auto_renew"] == false
	})).Return(nil)

	// Act
	err := service.CancelMembership(ctx, userID)

	// Assert
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_CancelMembership_NotActive 测试取消自动续费-会员未激活
func TestMembershipService_CancelMembership_NotActive(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	membership := createTestMembership(userID, plan.ID, financeModel.MembershipLevelVIPMonthly, time.Now(), time.Now().Add(-10*24*time.Hour))
	membership.Status = financeModel.MembershipStatusExpired

	mockRepo.On("GetMembership", ctx, userID).Return(membership, nil)

	// Act
	err := service.CancelMembership(ctx, userID)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "会员未激活或已过期")

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_RenewMembership_Success 测试手动续费成功
func TestMembershipService_RenewMembership_Success(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	oldMembership := createTestMembership(userID, plan.ID, financeModel.MembershipLevelVIPMonthly, time.Now(), time.Now().Add(5*24*time.Hour))

	newEndTime := time.Now().Add(35 * 24 * time.Hour)
	updatedMembership := createTestMembership(userID, plan.ID, financeModel.MembershipLevelVIPMonthly, time.Now(), newEndTime)
	updatedMembership.Status = financeModel.MembershipStatusActive

	mockRepo.On("GetMembership", ctx, userID).Return(oldMembership, nil)
	mockRepo.On("GetPlan", ctx, oldMembership.PlanID).Return(plan, nil)
	mockRepo.On("UpdateMembership", ctx, oldMembership.ID, mock.MatchedBy(func(u map[string]interface{}) bool {
		return u["end_time"] != nil && u["status"] == financeModel.MembershipStatusActive
	})).Return(nil)
	mockRepo.On("GetMembership", ctx, userID).Return(updatedMembership, nil)

	// Act
	membership, err := service.RenewMembership(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, membership)
	assert.Equal(t, financeModel.MembershipStatusActive, membership.Status)

	mockRepo.AssertExpectations(t)
}

// =========================
// 会员权益相关测试
// =========================

// TestMembershipService_GetBenefits_Success 测试获取权益列表成功
func TestMembershipService_GetBenefits_Success(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	level := financeModel.MembershipLevelVIPMonthly
	expectedBenefits := []*financeModel.MembershipBenefit{
		{ID: primitive.NewObjectID(), Level: level, Name: "免费阅读", Description: "所有书籍免费阅读"},
		{ID: primitive.NewObjectID(), Level: level, Name: "专属徽章", Description: "VIP专属徽章"},
	}

	mockRepo.On("ListBenefits", ctx, level, true).Return(expectedBenefits, nil)

	// Act
	benefits, err := service.GetBenefits(ctx, level)

	// Assert
	require.NoError(t, err)
	assert.Len(t, benefits, 2)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_GetUsage_Success 测试获取权益使用情况成功
func TestMembershipService_GetUsage_Success(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	expectedUsages := []*financeModel.MembershipUsage{
		{ID: primitive.NewObjectID(), UserID: userID, BenefitCode: "free_reading", UsageCount: 10},
	}

	mockRepo.On("ListUsages", ctx, userID).Return(expectedUsages, nil)

	// Act
	usages, err := service.GetUsage(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, usages, 1)

	mockRepo.AssertExpectations(t)
}

// =========================
// 会员卡管理相关测试
// =========================

// TestMembershipService_ActivateCard_Success 测试激活会员卡成功
func TestMembershipService_ActivateCard_Success(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	cardCode := "CARD123456"
	card := createTestMembershipCard(cardCode, financeModel.MembershipTypeMonthly, 30, financeModel.CardStatusUnused)

	mockRepo.On("GetMembershipCardByCode", ctx, cardCode).Return(card, nil)
	mockRepo.On("GetMembership", ctx, userID).Return(nil, assert.AnError)
	mockRepo.On("CreateMembership", ctx, mock.MatchedBy(func(m *financeModel.UserMembership) bool {
		return m.UserID == userID && m.Status == financeModel.MembershipStatusActive
	})).Return(nil)
	mockRepo.On("UpdateMembershipCard", ctx, card.ID, mock.MatchedBy(func(u map[string]interface{}) bool {
		status, ok := u["status"].(string)
		return ok && status == financeModel.CardStatusUsed && u["activated_by"] == userID
	})).Return(nil)

	// Act
	membership, err := service.ActivateCard(ctx, userID, cardCode)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, membership)
	assert.Equal(t, userID, membership.UserID)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_ActivateCard_CardNotFound 测试激活会员卡-卡不存在
func TestMembershipService_ActivateCard_CardNotFound(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	cardCode := "INVALID_CODE"

	mockRepo.On("GetMembershipCardByCode", ctx, cardCode).Return(nil, assert.AnError)

	// Act
	membership, err := service.ActivateCard(ctx, userID, cardCode)

	// Assert
	require.Error(t, err)
	assert.Nil(t, membership)
	assert.Contains(t, err.Error(), "会员卡不存在或已失效")

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_ActivateCard_CardUsed 测试激活会员卡-卡已使用
func TestMembershipService_ActivateCard_CardUsed(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	cardCode := "USED_CARD"
	card := createTestMembershipCard(cardCode, financeModel.MembershipTypeMonthly, 30, financeModel.CardStatusUsed)

	mockRepo.On("GetMembershipCardByCode", ctx, cardCode).Return(card, nil)

	// Act
	membership, err := service.ActivateCard(ctx, userID, cardCode)

	// Assert
	require.Error(t, err)
	assert.Nil(t, membership)
	assert.Contains(t, err.Error(), "会员卡已被使用")

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_ActivateCard_CardExpired 测试激活会员卡-卡已过期
func TestMembershipService_ActivateCard_CardExpired(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	cardCode := "EXPIRED_CARD"
	expiredTime := time.Now().Add(-24 * time.Hour)
	card := createTestMembershipCard(cardCode, financeModel.MembershipTypeMonthly, 30, financeModel.CardStatusUnused)
	card.ExpireAt = &expiredTime

	mockRepo.On("GetMembershipCardByCode", ctx, cardCode).Return(card, nil)

	// Act
	membership, err := service.ActivateCard(ctx, userID, cardCode)

	// Assert
	require.Error(t, err)
	assert.Nil(t, membership)
	assert.Contains(t, err.Error(), "会员卡已过期")

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_ActivateCard_ExtendExisting 测试激活会员卡-延长现有会员
func TestMembershipService_ActivateCard_ExtendExisting(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	cardCode := "CARD123456"
	card := createTestMembershipCard(cardCode, financeModel.MembershipTypeMonthly, 30, financeModel.CardStatusUnused)
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)

	oldEndTime := time.Now().Add(10 * 24 * time.Hour)
	existingMembership := createTestMembership(userID, plan.ID, financeModel.MembershipLevelVIPMonthly, time.Now(), oldEndTime)
	existingMembership.Status = financeModel.MembershipStatusActive

	mockRepo.On("GetMembershipCardByCode", ctx, cardCode).Return(card, nil)
	mockRepo.On("GetMembership", ctx, userID).Return(existingMembership, nil)
	mockRepo.On("UpdateMembership", ctx, existingMembership.ID, mock.MatchedBy(func(u map[string]interface{}) bool {
		// 新的结束时间应该在旧结束时间之后30天
		newEndTime := oldEndTime.Add(30 * 24 * time.Hour)
		endTime, ok := u["end_time"].(time.Time)
		return ok && endTime.Equal(newEndTime) || endTime.After(newEndTime.Add(-time.Hour))
	})).Return(nil)
	mockRepo.On("GetMembershipByID", ctx, existingMembership.ID).Return(existingMembership, nil)
	mockRepo.On("UpdateMembershipCard", ctx, card.ID, mock.AnythingOfType("map[string]interface {}")).Return(nil)

	// Act
	membership, err := service.ActivateCard(ctx, userID, cardCode)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, membership)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_ListCards_Success 测试列出会员卡成功
func TestMembershipService_ListCards_Success(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	filter := map[string]interface{}{"status": financeModel.CardStatusUnused}
	page := 1
	pageSize := 20

	expectedCards := []*financeModel.MembershipCard{
		createTestMembershipCard("CARD001", financeModel.MembershipTypeMonthly, 30, financeModel.CardStatusUnused),
		createTestMembershipCard("CARD002", financeModel.MembershipTypeYearly, 365, financeModel.CardStatusUnused),
	}

	mockRepo.On("ListMembershipCards", ctx, filter, page, pageSize).Return(expectedCards, nil)
	mockRepo.On("CountMembershipCards", ctx, filter).Return(int64(2), nil)

	// Act
	cards, total, err := service.ListCards(ctx, filter, page, pageSize)

	// Assert
	require.NoError(t, err)
	assert.Len(t, cards, 2)
	assert.Equal(t, int64(2), total)

	mockRepo.AssertExpectations(t)
}

// =========================
// 会员检查相关测试
// =========================

// TestMembershipService_CheckMembership_Success 测试检查会员等级成功
func TestMembershipService_CheckMembership_Success(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	level := financeModel.MembershipLevelVIPMonthly
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	membership := createTestMembership(userID, plan.ID, level, time.Now(), time.Now().Add(30*24*time.Hour))
	membership.Status = financeModel.MembershipStatusActive

	mockRepo.On("GetMembership", ctx, userID).Return(membership, nil)

	// Act
	isMember, err := service.CheckMembership(ctx, userID, level)

	// Assert
	require.NoError(t, err)
	assert.True(t, isMember)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_CheckMembership_NotMember 测试检查会员等级-不是会员
func TestMembershipService_CheckMembership_NotMember(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	level := financeModel.MembershipLevelVIPMonthly

	mockRepo.On("GetMembership", ctx, userID).Return(nil, assert.AnError)

	// Act
	isMember, err := service.CheckMembership(ctx, userID, level)

	// Assert
	require.NoError(t, err)
	assert.False(t, isMember)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_CheckMembership_WrongLevel 测试检查会员等级-等级不符
func TestMembershipService_CheckMembership_WrongLevel(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	membership := createTestMembership(userID, plan.ID, financeModel.MembershipLevelVIPMonthly, time.Now(), time.Now().Add(30*24*time.Hour))
	membership.Status = financeModel.MembershipStatusActive

	mockRepo.On("GetMembership", ctx, userID).Return(membership, nil)

	// Act - 检查更高的等级
	isMember, err := service.CheckMembership(ctx, userID, financeModel.MembershipLevelSuperVIP)

	// Assert
	require.NoError(t, err)
	assert.False(t, isMember)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_CheckMembership_Expired 测试检查会员等级-已过期
func TestMembershipService_CheckMembership_Expired(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	level := financeModel.MembershipLevelVIPMonthly
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	membership := createTestMembership(userID, plan.ID, level, time.Now().Add(-60*24*time.Hour), time.Now().Add(-30*24*time.Hour))
	membership.Status = financeModel.MembershipStatusActive

	mockRepo.On("GetMembership", ctx, userID).Return(membership, nil)

	// Act
	isMember, err := service.CheckMembership(ctx, userID, level)

	// Assert
	require.NoError(t, err)
	assert.False(t, isMember)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_IsVIP_True 测试检查是否是VIP-是VIP
func TestMembershipService_IsVIP_True(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	plan := createTestPlan("月度VIP", financeModel.MembershipTypeMonthly, 30, 19.9, true)
	membership := createTestMembership(userID, plan.ID, financeModel.MembershipLevelVIPMonthly, time.Now(), time.Now().Add(30*24*time.Hour))
	membership.Status = financeModel.MembershipStatusActive

	mockRepo.On("GetMembership", ctx, userID).Return(membership, nil)

	// Act
	isVIP, err := service.IsVIP(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.True(t, isVIP)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_IsVIP_False 测试检查是否是VIP-不是VIP
func TestMembershipService_IsVIP_False(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"
	plan := createTestPlan("普通套餐", financeModel.MembershipTypeMonthly, 30, 9.9, true)
	membership := createTestMembership(userID, plan.ID, financeModel.MembershipLevelNormal, time.Now(), time.Now().Add(30*24*time.Hour))
	membership.Status = financeModel.MembershipStatusActive

	mockRepo.On("GetMembership", ctx, userID).Return(membership, nil)

	// Act
	isVIP, err := service.IsVIP(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.False(t, isVIP)

	mockRepo.AssertExpectations(t)
}

// TestMembershipService_IsVIP_NoMembership 测试检查是否是VIP-无会员
func TestMembershipService_IsVIP_NoMembership(t *testing.T) {
	// Arrange
	service, mockRepo := setupMembershipService()
	ctx := context.Background()

	userID := "user123"

	mockRepo.On("GetMembership", ctx, userID).Return(nil, assert.AnError)

	// Act
	isVIP, err := service.IsVIP(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.False(t, isVIP)

	mockRepo.AssertExpectations(t)
}
