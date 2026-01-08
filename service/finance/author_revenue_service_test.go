package finance

import (
	"context"
	"testing"
	"time"

	financeModel "Qingyu_backend/models/finance"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ Mock 实现 ============

// MockAuthorRevenueRepository Mock作者收入仓储
type MockAuthorRevenueRepository struct {
	earnings        map[string][]*financeModel.AuthorEarning  // authorID -> earnings
	revenueDetails  map[string][]*financeModel.RevenueDetail   // authorID -> details
	revenueStats    map[string][]*financeModel.RevenueStatistics // authorID -> stats
	withdrawals     map[string]*financeModel.WithdrawalRequest  // userID -> request
	settlements     map[string][]*financeModel.Settlement       // authorID -> settlements
	taxInfos        map[string]*financeModel.TaxInfo            // userID -> taxInfo
	nextID          int
}

func NewMockAuthorRevenueRepository() *MockAuthorRevenueRepository {
	return &MockAuthorRevenueRepository{
		earnings:       make(map[string][]*financeModel.AuthorEarning),
		revenueDetails: make(map[string][]*financeModel.RevenueDetail),
		revenueStats:   make(map[string][]*financeModel.RevenueStatistics),
		withdrawals:    make(map[string]*financeModel.WithdrawalRequest),
		settlements:    make(map[string][]*financeModel.Settlement),
		taxInfos:       make(map[string]*financeModel.TaxInfo),
		nextID:         1,
	}
}

func (m *MockAuthorRevenueRepository) GetEarningsByAuthor(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error) {
	earnings := m.earnings[authorID]
	if earnings == nil {
		return []*financeModel.AuthorEarning{}, 0, nil
	}
	return earnings, int64(len(earnings)), nil
}

func (m *MockAuthorRevenueRepository) GetEarningsByBook(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error) {
	result := make([]*financeModel.AuthorEarning, 0)
	for _, earnings := range m.earnings {
		for _, earning := range earnings {
			if earning.BookID == bookID {
				result = append(result, earning)
			}
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockAuthorRevenueRepository) CreateEarning(ctx context.Context, earning *financeModel.AuthorEarning) error {
	earning.ID = primitive.NewObjectID()
	earning.CreatedAt = time.Now()

	authorID := earning.AuthorID
	if m.earnings[authorID] == nil {
		m.earnings[authorID] = make([]*financeModel.AuthorEarning, 0)
	}
	m.earnings[authorID] = append(m.earnings[authorID], earning)

	return nil
}

func (m *MockAuthorRevenueRepository) GetRevenueDetails(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.RevenueDetail, int64, error) {
	details := m.revenueDetails[authorID]
	if details == nil {
		return []*financeModel.RevenueDetail{}, 0, nil
	}
	return details, int64(len(details)), nil
}

func (m *MockAuthorRevenueRepository) CreateRevenueDetail(ctx context.Context, detail *financeModel.RevenueDetail) error {
	authorID := detail.AuthorID
	if m.revenueDetails[authorID] == nil {
		m.revenueDetails[authorID] = make([]*financeModel.RevenueDetail, 0)
	}
	detail.ID = primitive.NewObjectID()
	detail.CreatedAt = time.Now()
	m.revenueDetails[authorID] = append(m.revenueDetails[authorID], detail)
	return nil
}

func (m *MockAuthorRevenueRepository) GetRevenueStatistics(ctx context.Context, authorID string, period string, limit int) ([]*financeModel.RevenueStatistics, error) {
	stats := m.revenueStats[authorID]
	if stats == nil {
		return []*financeModel.RevenueStatistics{}, nil
	}
	return stats, nil
}

func (m *MockAuthorRevenueRepository) CreateWithdrawalRequest(ctx context.Context, request *financeModel.WithdrawalRequest) error {
	request.ID = primitive.NewObjectID()
	request.CreatedAt = time.Now()
	request.Status = financeModel.WithdrawalStatusPending
	m.withdrawals[request.UserID] = request
	return nil
}

func (m *MockAuthorRevenueRepository) GetWithdrawalRequests(ctx context.Context, userID string, page, pageSize int) ([]*financeModel.WithdrawalRequest, int64, error) {
	request := m.withdrawals[userID]
	if request == nil {
		return []*financeModel.WithdrawalRequest{}, 0, nil
	}
	return []*financeModel.WithdrawalRequest{request}, 1, nil
}

func (m *MockAuthorRevenueRepository) GetSettlements(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.Settlement, int64, error) {
	settlements := m.settlements[authorID]
	if settlements == nil {
		return []*financeModel.Settlement{}, 0, nil
	}
	return settlements, int64(len(settlements)), nil
}

func (m *MockAuthorRevenueRepository) GetSettlement(ctx context.Context, settlementID string) (*financeModel.Settlement, error) {
	for _, settlements := range m.settlements {
		for _, settlement := range settlements {
			if settlement.ID.Hex() == settlementID {
				return settlement, nil
			}
		}
	}
	return nil, nil
}

func (m *MockAuthorRevenueRepository) GetTaxInfo(ctx context.Context, userID string) (*financeModel.TaxInfo, error) {
	return m.taxInfos[userID], nil
}

func (m *MockAuthorRevenueRepository) UpdateTaxInfo(ctx context.Context, taxInfo *financeModel.TaxInfo) error {
	m.taxInfos[taxInfo.UserID] = taxInfo
	return nil
}

// ============ 测试辅助函数 ============

func setupTestAuthorRevenueService() (*AuthorRevenueServiceImpl, *MockAuthorRevenueRepository) {
	repo := NewMockAuthorRevenueRepository()
	service := NewAuthorRevenueService(repo).(*AuthorRevenueServiceImpl)
	return service, repo
}

func createTestEarning(authorID, bookID string, amount float64, earningType string) *financeModel.AuthorEarning {
	bookOID, _ := primitive.ObjectIDFromHex(bookID)
	return &financeModel.AuthorEarning{
		AuthorID:      authorID,
		BookID:        bookOID,
		BookTitle:     "测试书籍",
		Type:          earningType,
		Amount:        amount,
		AuthorIncome:  amount * 0.7, // 假设作者分成70%
		PlatformIncome: amount * 0.3,
		Status:        financeModel.EarningStatusSettled,
	}
}

// ============ 测试用例 ============

// TestGetEarnings 测试获取作者收入列表
func TestGetEarnings(t *testing.T) {
	service, repo := setupTestAuthorRevenueService()
	ctx := context.Background()

	authorID := "author_123"

	// 创建测试收入
	earning1 := createTestEarning(authorID, "507f1f77bcf86cd799439011", 100.0, financeModel.EarningTypeChapterPurchase)
	earning2 := createTestEarning(authorID, "507f1f77bcf86cd799439012", 50.0, financeModel.EarningTypeReward)
	_ = repo.CreateEarning(ctx, earning1)
	_ = repo.CreateEarning(ctx, earning2)

	// 获取收入列表
	earnings, total, err := service.GetEarnings(ctx, authorID, 1, 10)
	if err != nil {
		t.Fatalf("获取收入列表失败: %v", err)
	}

	if len(earnings) != 2 {
		t.Errorf("收入数量错误: 期望2条，实际%d条", len(earnings))
	}
	if total != 2 {
		t.Errorf("总数错误: 期望2，实际%d", total)
	}

	t.Logf("获取收入列表成功: %d条，总计%.2f元", len(earnings), total)
}

// TestGetBookEarnings 测试获取书籍收入
func TestGetBookEarnings(t *testing.T) {
	service, repo := setupTestAuthorRevenueService()
	ctx := context.Background()

	authorID := "author_456"
	bookID := "507f1f77bcf86cd799439013"

	// 创建测试收入
	earning1 := createTestEarning(authorID, bookID, 100.0, financeModel.EarningTypeChapterPurchase)
	earning2 := createTestEarning(authorID, bookID, 50.0, financeModel.EarningTypeReward)
	_ = repo.CreateEarning(ctx, earning1)
	_ = repo.CreateEarning(ctx, earning2)

	// 获取书籍收入
	earnings, total, err := service.GetBookEarnings(ctx, authorID, bookID)
	if err != nil {
		t.Fatalf("获取书籍收入失败: %v", err)
	}

	if len(earnings) != 2 {
		t.Errorf("收入数量错误: 期望2条，实际%d条", len(earnings))
	}
	if total != 2 {
		t.Errorf("总数错误: 期望2，实际%d", total)
	}

	t.Logf("获取书籍收入成功: %d条，总计%.2f元", len(earnings), total)
}

// TestGetRevenueDetails 测试获取收入明细
func TestGetRevenueDetails(t *testing.T) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	authorID := "author_789"

	details, total, err := service.GetRevenueDetails(ctx, authorID, 1, 10)
	if err != nil {
		t.Fatalf("获取收入明细失败: %v", err)
	}

	t.Logf("获取收入明细成功: %d条，总计%d", len(details), total)
}

// TestGetRevenueStatistics 测试获取收入统计
func TestGetRevenueStatistics(t *testing.T) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	authorID := "author_stats"
	period := "2024-01"

	stats, err := service.GetRevenueStatistics(ctx, authorID, period)
	if err != nil {
		t.Fatalf("获取收入统计失败: %v", err)
	}

	t.Logf("获取收入统计成功: %d条记录", len(stats))
}

// TestCreateEarning 测试创建收入记录
func TestCreateEarning(t *testing.T) {
	service, repo := setupTestAuthorRevenueService()
	ctx := context.Background()

	authorID := "author_create"
	bookID := "507f1f77bcf86cd799439014"

	earning := createTestEarning(authorID, bookID, 100.0, financeModel.EarningTypeChapterPurchase)

	err := service.CreateEarning(ctx, earning)
	if err != nil {
		t.Fatalf("创建收入记录失败: %v", err)
	}

	// 验证收入记录已创建
	earnings := repo.earnings[authorID]
	if len(earnings) != 1 {
		t.Errorf("收入记录数量错误: 期望1条，实际%d条", len(earnings))
	}

	t.Logf("创建收入记录成功: 作者=%s, 金额=%.2f元", authorID, earning.Amount)
}

// TestCalculateEarning_ChapterPurchase 测试计算章节购买收入分成
func TestCalculateEarning_ChapterPurchase(t *testing.T) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	authorID := "author_calculate"
	bookID := "507f1f77bcf86cd799439015"
	bookOID, _ := primitive.ObjectIDFromHex(bookID)

	amount := 100.0
	authorIncome, platformIncome, err := service.CalculateEarning(ctx, financeModel.EarningTypeChapterPurchase, amount, authorID, bookOID)
	if err != nil {
		t.Fatalf("计算收入分成失败: %v", err)
	}

	expectedAuthorIncome := amount * financeModel.ChapterPurchaseAuthorRate
	expectedPlatformIncome := amount * financeModel.ChapterPurchasePlatformRate

	if authorIncome != expectedAuthorIncome {
		t.Errorf("作者收入错误: 期望%.2f，实际%.2f", expectedAuthorIncome, authorIncome)
	}
	if platformIncome != expectedPlatformIncome {
		t.Errorf("平台收入错误: 期望%.2f，实际%.2f", expectedPlatformIncome, platformIncome)
	}

	t.Logf("计算章节购买收入分成成功: 总金额=%.2f, 作者收入=%.2f, 平台收入=%.2f",
		amount, authorIncome, platformIncome)
}

// TestCalculateEarning_Reward 测试计算打赏收入分成
func TestCalculateEarning_Reward(t *testing.T) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	authorID := "author_reward"
	bookID := "507f1f77bcf86cd799439016"
	bookOID, _ := primitive.ObjectIDFromHex(bookID)

	amount := 50.0
	authorIncome, platformIncome, err := service.CalculateEarning(ctx, financeModel.EarningTypeReward, amount, authorID, bookOID)
	if err != nil {
		t.Fatalf("计算打赏分成失败: %v", err)
	}

	expectedAuthorIncome := amount * financeModel.RewardAuthorRate
	expectedPlatformIncome := amount * financeModel.RewardPlatformRate

	if authorIncome != expectedAuthorIncome {
		t.Errorf("作者收入错误: 期望%.2f，实际%.2f", expectedAuthorIncome, authorIncome)
	}

	t.Logf("计算打赏分成成功: 总金额=%.2f, 作者收入=%.2f, 平台收入=%.2f",
		amount, authorIncome, platformIncome)
}

// TestCalculateEarning_VIPReading 测试计算VIP阅读收入分成
func TestCalculateEarning_VIPReading(t *testing.T) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	authorID := "author_vip"
	bookID := "507f1f77bcf86cd799439017"
	bookOID, _ := primitive.ObjectIDFromHex(bookID)

	amount := 10.0
	authorIncome, platformIncome, err := service.CalculateEarning(ctx, financeModel.EarningTypeVIPReading, amount, authorID, bookOID)
	if err != nil {
		t.Fatalf("计算VIP阅读分成失败: %v", err)
	}

	expectedAuthorIncome := amount * financeModel.VIPReadingAuthorRate
	expectedPlatformIncome := amount * financeModel.VIPReadingPlatformRate

	if authorIncome != expectedAuthorIncome {
		t.Errorf("作者收入错误: 期望%.2f，实际%.2f", expectedAuthorIncome, authorIncome)
	}

	t.Logf("计算VIP阅读分成成功: 总金额=%.2f, 作者收入=%.2f, 平台收入=%.2f",
		amount, authorIncome, platformIncome)
}

// TestCalculateEarning_UnknownType 测试计算未知类型收入分成
func TestCalculateEarning_UnknownType(t *testing.T) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	authorID := "author_unknown"
	bookID := "507f1f77bcf86cd799439018"
	bookOID, _ := primitive.ObjectIDFromHex(bookID)

	amount := 100.0
	_, _, err := service.CalculateEarning(ctx, "unknown_type", amount, authorID, bookOID)
	if err == nil {
		t.Fatal("应该拒绝未知收入类型，但成功了")
	}

	t.Logf("正确拒绝了未知收入类型: %v", err)
}

// TestCreateWithdrawalRequest 测试创建提现申请
func TestCreateWithdrawalRequest(t *testing.T) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	userID := "user_withdraw"
	amount := 100.0
	method := "alipay"
	account := financeModel.WithdrawAccount{
		Type:     "支付宝",
		Account:  "test@example.com",
		RealName: "测试用户",
	}

	request, err := service.CreateWithdrawalRequest(ctx, userID, amount, method, account)
	if err != nil {
		t.Fatalf("创建提现申请失败: %v", err)
	}

	if request == nil {
		t.Fatal("提现申请不应为空")
	}
	if request.UserID != userID {
		t.Errorf("用户ID错误: %s", request.UserID)
	}
	if request.Amount != amount {
		t.Errorf("提现金额错误: 期望%.2f，实际%.2f", amount, request.Amount)
	}
	if request.Status != financeModel.WithdrawalStatusPending {
		t.Errorf("状态应该是待处理: %s", request.Status)
	}

	t.Logf("创建提现申请成功: 用户=%s, 金额=%.2f元, 方式=%s", userID, amount, method)
}

// TestCreateWithdrawalRequest_ZeroAmount 测试创建零金额提现申请
func TestCreateWithdrawalRequest_ZeroAmount(t *testing.T) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	userID := "user_zero"
	amount := 0.0
	method := "wechat"
	account := financeModel.WithdrawAccount{}

	_, err := service.CreateWithdrawalRequest(ctx, userID, amount, method, account)
	if err == nil {
		t.Fatal("应该拒绝零金额提现，但成功了")
	}

	t.Logf("正确拒绝了零金额提现: %v", err)
}

// TestGetWithdrawals 测试获取提现申请列表
func TestGetWithdrawals(t *testing.T) {
	service, repo := setupTestAuthorRevenueService()
	ctx := context.Background()

	userID := "user_list_withdraw"

	// 创建提现申请
	request := &financeModel.WithdrawalRequest{
		UserID:   userID,
		Amount:   100.0,
		Method:   "alipay",
		Status:   financeModel.WithdrawalStatusPending,
	}
	_ = repo.CreateWithdrawalRequest(ctx, request)

	// 获取提现申请列表
	requests, total, err := service.GetWithdrawals(ctx, userID, 1, 10)
	if err != nil {
		t.Fatalf("获取提现申请列表失败: %v", err)
	}

	if len(requests) != 1 {
		t.Errorf("提现申请数量错误: 期望1条，实际%d条", len(requests))
	}
	if total != 1 {
		t.Errorf("总数错误: 期望1，实际%d", total)
	}

	t.Logf("获取提现申请列表成功: %d条，总计%d", len(requests), total)
}

// TestGetSettlements 测试获取结算列表
func TestGetSettlements(t *testing.T) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	authorID := "author_settlement"

	settlements, total, err := service.GetSettlements(ctx, authorID, 1, 10)
	if err != nil {
		t.Fatalf("获取结算列表失败: %v", err)
	}

	t.Logf("获取结算列表成功: %d条，总计%d", len(settlements), total)
}

// TestGetSettlement 测试获取结算详情
func TestGetSettlement(t *testing.T) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	settlementID := "507f1f77bcf86cd799439019"

	settlement, err := service.GetSettlement(ctx, settlementID)
	if err != nil {
		t.Fatalf("获取结算详情失败: %v", err)
	}

	if settlement != nil {
		t.Logf("获取结算详情成功: ID=%s", settlement.ID.Hex())
	} else {
		t.Logf("结算详情不存在: ID=%s", settlementID)
	}
}

// TestGetTaxInfo 测试获取税务信息
func TestGetTaxInfo(t *testing.T) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	userID := "user_tax"

	taxInfo, err := service.GetTaxInfo(ctx, userID)
	if err != nil {
		t.Fatalf("获取税务信息失败: %v", err)
	}

	if taxInfo != nil {
		t.Logf("获取税务信息成功: 用户=%s", taxInfo.UserID)
	} else {
		t.Logf("税务信息不存在: 用户=%s", userID)
	}
}

// TestUpdateTaxInfo 测试更新税务信息
func TestUpdateTaxInfo(t *testing.T) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	userID := "user_update_tax"
	taxInfo := &financeModel.TaxInfo{
		UserID:        userID,
		IDType:        "身份证",
		IDNumber:      "123456789012345678",
		RealName:      "测试用户",
		TaxType:       financeModel.TaxTypeIndividual,
		IsVerified:    false,
	}

	err := service.UpdateTaxInfo(ctx, userID, taxInfo)
	if err != nil {
		t.Fatalf("更新税务信息失败: %v", err)
	}

	t.Logf("更新税务信息成功: 用户=%s, 姓名=%s", userID, taxInfo.RealName)
}

// BenchmarkCalculateEarning 性能测试：计算收入分成
func BenchmarkCalculateEarning(b *testing.B) {
	service, _ := setupTestAuthorRevenueService()
	ctx := context.Background()

	authorID := "author_bench"
	bookID := "507f1f77bcf86cd799439020"
	bookOID, _ := primitive.ObjectIDFromHex(bookID)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = service.CalculateEarning(ctx, financeModel.EarningTypeChapterPurchase, 100.0, authorID, bookOID)
	}
}
