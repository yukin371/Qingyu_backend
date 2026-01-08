package finance

import (
	"context"
	"testing"
	"time"

	financeModel "Qingyu_backend/models/finance"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ Mock 实现 ============

// MockMembershipRepository Mock会员仓储
type MockMembershipRepository struct {
	plans       map[string]*financeModel.MembershipPlan
	memberships map[string]*financeModel.UserMembership // userID -> membership
	cards       map[string]*financeModel.MembershipCard
	benefits    map[string][]*financeModel.MembershipBenefit // level -> benefits
	usages      map[string][]*financeModel.MembershipUsage    // userID -> usages
	nextID      int
}

func NewMockMembershipRepository() *MockMembershipRepository {
	repo := &MockMembershipRepository{
		plans:       make(map[string]*financeModel.MembershipPlan),
		memberships: make(map[string]*financeModel.UserMembership),
		cards:       make(map[string]*financeModel.MembershipCard),
		benefits:    make(map[string][]*financeModel.MembershipBenefit),
		usages:      make(map[string][]*financeModel.MembershipUsage),
		nextID:      1,
	}

	// 添加测试套餐
	monthlyPlan := &financeModel.MembershipPlan{
		ID:          primitive.NewObjectID(),
		Name:        "月度会员",
		Type:        financeModel.PlanTypeMonthly,
		Description: "30天会员",
		Price:       19.9,
		Duration:    30,
		IsEnabled:   true,
		Discount:    0,
		Benefits:    []string{"无限阅读", "免费章节"},
	}
	repo.plans[monthlyPlan.ID.Hex()] = monthlyPlan

	yearlyPlan := &financeModel.MembershipPlan{
		ID:          primitive.NewObjectID(),
		Name:        "年度会员",
		Type:        financeModel.PlanTypeYearly,
		Description: "365天会员",
		Price:       199.0,
		Duration:    365,
		IsEnabled:   true,
		Discount:    17,
		Benefits:    []string{"无限阅读", "免费章节", "专属标识"},
	}
	repo.plans[yearlyPlan.ID.Hex()] = yearlyPlan

	// 添加测试会员卡
	card := &financeModel.MembershipCard{
		ID:          primitive.NewObjectID(),
		Code:        "TEST_CARD_001",
		Type:        financeModel.PlanTypeMonthly,
		Duration:    30,
		Status:      financeModel.CardStatusActive,
		MaxUses:     100,
		UsedCount:   10,
		ExpiresAt:   time.Now().Add(30 * 24 * time.Hour),
		CreatedAt:   time.Now(),
	}
	repo.cards[card.ID.Hex()] = card

	// 添加测试权益
	benefits := []*financeModel.MembershipBenefit{
		{ID: primitive.NewObjectID(), Level: "normal", Name: "基础阅读", Description: "可以阅读所有免费书籍"},
		{ID: primitive.NewObjectID(), Level: "vip", Name: "无限阅读", Description: "可以阅读所有书籍包括付费内容"},
	}
	repo.benefits["vip"] = benefits

	return repo
}

func (m *MockMembershipRepository) ListPlans(ctx context.Context, enabledOnly bool) ([]*financeModel.MembershipPlan, error) {
	result := make([]*financeModel.MembershipPlan, 0)
	for _, plan := range m.plans {
		if !enabledOnly || plan.IsEnabled {
			result = append(result, plan)
		}
	}
	return result, nil
}

func (m *MockMembershipRepository) GetPlan(ctx context.Context, id primitive.ObjectID) (*financeModel.MembershipPlan, error) {
	return m.plans[id.Hex()], nil
}

func (m *MockMembershipRepository) CreateMembership(ctx context.Context, membership *financeModel.UserMembership) error {
	membership.ID = primitive.NewObjectID()
	m.memberships[membership.UserID] = membership
	return nil
}

func (m *MockMembershipRepository) GetMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error) {
	return m.memberships[userID], nil
}

func (m *MockMembershipRepository) UpdateMembership(ctx context.Context, membership *financeModel.UserMembership) error {
	m.memberships[membership.UserID] = membership
	return nil
}

func (m *MockMembershipRepository) DeleteMembership(ctx context.Context, userID string) error {
	delete(m.memberships, userID)
	return nil
}

func (m *MockMembershipRepository) GetCard(ctx context.Context, cardCode string) (*financeModel.MembershipCard, error) {
	for _, card := range m.cards {
		if card.Code == cardCode {
			return card, nil
		}
	}
	return nil, nil
}

func (m *MockMembershipRepository) UpdateCard(ctx context.Context, card *financeModel.MembershipCard) error {
	m.cards[card.ID.Hex()] = card
	return nil
}

func (m *MockMembershipRepository) ListCards(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.MembershipCard, int64, error) {
	result := make([]*financeModel.MembershipCard, 0)
	for _, card := range m.cards {
		result = append(result, card)
	}
	return result, int64(len(result)), nil
}

func (m *MockMembershipRepository) GetBenefits(ctx context.Context, level string) ([]*financeModel.MembershipBenefit, error) {
	return m.benefits[level], nil
}

func (m *MockMembershipRepository) GetUsage(ctx context.Context, userID string) ([]*financeModel.MembershipUsage, error) {
	return m.usages[userID], nil
}

func (m *MockMembershipRepository) CreateUsage(ctx context.Context, usage *financeModel.MembershipUsage) error {
	if m.usages[usage.UserID] == nil {
		m.usages[usage.UserID] = make([]*financeModel.MembershipUsage, 0)
	}
	m.usages[usage.UserID] = append(m.usages[usage.UserID], usage)
	return nil
}

// ============ 测试辅助函数 ============

func setupTestMembershipService() (*MembershipServiceImpl, *MockMembershipRepository) {
	repo := NewMockMembershipRepository()
	service := NewMembershipService(repo).(*MembershipServiceImpl)
	return service, repo
}

// ============ 测试用例 ============

// TestGetPlans 测试获取套餐列表
func TestGetPlans(t *testing.T) {
	service, repo := setupTestMembershipService()
	ctx := context.Background()

	plans, err := service.GetPlans(ctx)
	if err != nil {
		t.Fatalf("获取套餐列表失败: %v", err)
	}

	if len(plans) != len(repo.plans) {
		t.Errorf("套餐数量错误: 期望%d，实际%d", len(repo.plans), len(plans))
	}

	t.Logf("获取套餐列表成功: %d个套餐", len(plans))
}

// TestGetPlan 测试获取套餐详情
func TestGetPlan(t *testing.T) {
	service, repo := setupTestMembershipService()
	ctx := context.Background()

	// 获取第一个套餐ID
	var planID string
	for id := range repo.plans {
		planID = id
		break
	}

	plan, err := service.GetPlan(ctx, planID)
	if err != nil {
		t.Fatalf("获取套餐失败: %v", err)
	}

	if plan == nil {
		t.Fatal("套餐不应为空")
	}

	t.Logf("获取套餐成功: %s, 价格: %.1f", plan.Name, plan.Price)
}

// TestSubscribe 测试订阅会员
func TestSubscribe(t *testing.T) {
	service, repo := setupTestMembershipService()
	ctx := context.Background()

	// 获取第一个套餐ID
	var planID string
	for id := range repo.plans {
		planID = id
		break
	}

	userID := "user_123"
	paymentMethod := "wechat"

	membership, err := service.Subscribe(ctx, userID, planID, paymentMethod)
	if err != nil {
		t.Fatalf("订阅会员失败: %v", err)
	}

	if membership == nil {
		t.Fatal("会员记录不应为空")
	}
	if membership.UserID != userID {
		t.Errorf("用户ID错误: %s", membership.UserID)
	}
	if membership.Status != financeModel.MembershipStatusActive {
		t.Errorf("会员状态错误: %s", membership.Status)
	}

	t.Logf("订阅会员成功: 用户=%s, 套餐=%s, 到期时间=%s",
		userID, membership.PlanName, membership.EndTime.Format("2006-01-02"))
}

// TestSubscribe_Duplicate 测试重复订阅
func TestSubscribe_Duplicate(t *testing.T) {
	service, _ := setupTestMembershipService()
	ctx := context.Background()

	// 获取第一个套餐ID
	var planID string
	plans, _ := service.GetPlans(ctx)
	if len(plans) > 0 {
		planID = plans[0].ID.Hex()
	}

	userID := "user_456"

	// 第一次订阅
	_, err := service.Subscribe(ctx, userID, planID, "wechat")
	if err != nil {
		t.Fatalf("第一次订阅失败: %v", err)
	}

	// 第二次订阅（应该失败）
	_, err = service.Subscribe(ctx, userID, planID, "wechat")
	if err == nil {
		t.Fatal("应该拒绝重复订阅，但成功了")
	}

	t.Logf("正确拒绝了重复订阅: %v", err)
}

// TestGetMembership 测试获取会员状态
func TestGetMembership(t *testing.T) {
	service, _ := setupTestMembershipService()
	ctx := context.Background()

	// 先订阅
	var planID string
	plans, _ := service.GetPlans(ctx)
	if len(plans) > 0 {
		planID = plans[0].ID.Hex()
	}

	userID := "user_789"
	_, _ = service.Subscribe(ctx, userID, planID, "wechat")

	// 获取会员状态
	membership, err := service.GetMembership(ctx, userID)
	if err != nil {
		t.Fatalf("获取会员状态失败: %v", err)
	}

	if membership == nil {
		t.Fatal("会员记录不应为空")
	}

	t.Logf("获取会员状态成功: 用户=%s, 等级=%s, 状态=%s",
		userID, membership.Level, membership.Status)
}

// TestCancelMembership 测试取消会员
func TestCancelMembership(t *testing.T) {
	service, _ := setupTestMembershipService()
	ctx := context.Background()

	// 先订阅
	var planID string
	plans, _ := service.GetPlans(ctx)
	if len(plans) > 0 {
		planID = plans[0].ID.Hex()
	}

	userID := "user_cancel"
	_, _ = service.Subscribe(ctx, userID, planID, "wechat")

	// 取消会员
	err := service.CancelMembership(ctx, userID)
	if err != nil {
		t.Fatalf("取消会员失败: %v", err)
	}

	// 验证状态
	membership, _ := service.GetMembership(ctx, userID)
	if membership.Status != financeModel.MembershipStatusCancelled {
		t.Errorf("会员状态应该是已取消: %s", membership.Status)
	}

	t.Logf("取消会员成功")
}

// TestRenewMembership 测试续费会员
func TestRenewMembership(t *testing.T) {
	service, _ := setupTestMembershipService()
	ctx := context.Background()

	// 先订阅
	var planID string
	plans, _ := service.GetPlans(ctx)
	if len(plans) > 0 {
		planID = plans[0].ID.Hex()
	}

	userID := "user_renew"
	membership, _ := service.Subscribe(ctx, userID, planID, "wechat")

	// 记录原始到期时间
	originalEndTime := membership.EndTime

	// 续费
	renewed, err := service.RenewMembership(ctx, userID)
	if err != nil {
		t.Fatalf("续费会员失败: %v", err)
	}

	if renewed.EndTime.Before(originalEndTime) || renewed.EndTime.Equal(originalEndTime) {
		t.Errorf("续费后到期时间应该延后: 原始=%s, 新=%s",
			originalEndTime.Format("2006-01-02"), renewed.EndTime.Format("2006-01-02"))
	}

	t.Logf("续费会员成功: 新到期时间=%s", renewed.EndTime.Format("2006-01-02"))
}

// TestGetBenefits 测试获取会员权益
func TestGetBenefits(t *testing.T) {
	service, _ := setupTestMembershipService()
	ctx := context.Background()

	level := "vip"

	benefits, err := service.GetBenefits(ctx, level)
	if err != nil {
		t.Fatalf("获取会员权益失败: %v", err)
	}

	if len(benefits) == 0 {
		t.Error("权益列表不应为空")
	}

	t.Logf("获取会员权益成功: %d个权益", len(benefits))
}

// TestGetUsage 测试获取会员使用情况
func TestGetUsage(t *testing.T) {
	service, _ := setupTestMembershipService()
	ctx := context.Background()

	userID := "user_usage"

	usages, err := service.GetUsage(ctx, userID)
	if err != nil {
		t.Fatalf("获取使用情况失败: %v", err)
	}

	// 可能没有使用记录，这是正常的
	t.Logf("获取使用情况成功: %d条记录", len(usages))
}

// TestCheckMembership 测试检查会员资格
func TestCheckMembership(t *testing.T) {
	service, _ := setupTestMembershipService()
	ctx := context.Background()

	userID := "user_check"
	level := "vip"

	// 先订阅
	var planID string
	plans, _ := service.GetPlans(ctx)
	if len(plans) > 0 {
		planID = plans[0].ID.Hex()
	}
	_, _ = service.Subscribe(ctx, userID, planID, "wechat")

	// 检查会员资格
	hasMembership, err := service.CheckMembership(ctx, userID, level)
	if err != nil {
		t.Fatalf("检查会员资格失败: %v", err)
	}

	t.Logf("检查会员资格成功: 用户=%s, 等级=%s, 有会员=%v", userID, level, hasMembership)
}

// TestIsVIP 测试检查是否为VIP
func TestIsVIP(t *testing.T) {
	service, _ := setupTestMembershipService()
	ctx := context.Background()

	userID := "user_vip"

	// 先订阅
	var planID string
	plans, _ := service.GetPlans(ctx)
	if len(plans) > 0 {
		planID = plans[0].ID.Hex()
	}
	_, _ = service.Subscribe(ctx, userID, planID, "wechat")

	// 检查VIP
	isVIP, err := service.IsVIP(ctx, userID)
	if err != nil {
		t.Fatalf("检查VIP失败: %v", err)
	}

	if !isVIP {
		t.Error("用户应该是VIP")
	}

	t.Logf("检查VIP成功: 用户=%s, 是VIP=%v", userID, isVIP)
}

// TestActivateCard 测试激活会员卡
func TestActivateCard(t *testing.T) {
	service, repo := setupTestMembershipService()
	ctx := context.Background()

	// 获取测试卡号
	var cardCode string
	for _, card := range repo.cards {
		cardCode = card.Code
		break
	}

	userID := "user_card"

	membership, err := service.ActivateCard(ctx, userID, cardCode)
	if err != nil {
		t.Fatalf("激活会员卡失败: %v", err)
	}

	if membership == nil {
		t.Fatal("会员记录不应为空")
	}

	t.Logf("激活会员卡成功: 用户=%s, 卡号=%s", userID, cardCode)
}

// TestListCards 测试列出会员卡
func TestListCards(t *testing.T) {
	service, _ := setupTestMembershipService()
	ctx := context.Background()

	cards, total, err := service.ListCards(ctx, map[string]interface{}{}, 1, 10)
	if err != nil {
		t.Fatalf("列出会员卡失败: %v", err)
	}

	t.Logf("列出会员卡成功: %d张，总计%d", len(cards), total)
}

// BenchmarkGetPlans 性能测试：获取套餐列表
func BenchmarkGetPlans(b *testing.B) {
	service, _ := setupTestMembershipService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetPlans(ctx)
	}
}
