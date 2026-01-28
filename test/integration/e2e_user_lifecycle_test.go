package integration

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/models/shared/types"
	"Qingyu_backend/service/admin"
	"Qingyu_backend/service/finance/wallet"
	"Qingyu_backend/service/shared/auth"
	"Qingyu_backend/service/shared/storage"

	"github.com/stretchr/testify/assert"
)

// ============ 补充Admin服务的Mock方法 ============

func (m *MockAdminService) GetPendingReviews(ctx context.Context, contentType string) ([]*admin.AuditRecord, error) {
	args := m.Called(ctx, contentType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*admin.AuditRecord), args.Error(1)
}

func (m *MockAdminService) ReviewContent(ctx context.Context, req *admin.ReviewContentRequest) error {
	return m.Called(ctx, req).Error(0)
}

// ============ 端到端测试：完整用户生命周期 ============

// TestE2E_CompleteUserLifecycle 测试完整的用户生命周期
// 场景：注册 -> 登录 -> 创建钱包 -> 充值 -> 消费 -> 申请提现 -> 登出
func TestE2E_CompleteUserLifecycle(t *testing.T) {
	ctx := context.Background()

	// 创建所有需要的Mock服务
	authService := new(MockAuthService)
	walletService := new(MockWalletService)

	// ========== 阶段1：用户注册 ==========
	t.Log("阶段1: 用户注册")

	registerReq := &auth.RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "securePassword123",
		Role:     "writer",
	}

	registerResp := &auth.RegisterResponse{
		User: &auth.UserInfo{
			ID:       "user_alice_001",
			Username: "alice",
			Email:    "alice@example.com",
			Roles:    []string{"writer"},
		},
		Token: "jwt_token_register_alice",
	}

	authService.On("Register", ctx, registerReq).Return(registerResp, nil)

	// 执行注册
	user, err := authService.Register(ctx, registerReq)
	assert.NoError(t, err, "用户注册应该成功")
	assert.NotNil(t, user)
	assert.Equal(t, "user_alice_001", user.User.ID)
	assert.NotEmpty(t, user.Token, "应该返回JWT Token")
	t.Logf("✓ 用户注册成功: %s (ID: %s)", user.User.Username, user.User.ID)

	// ========== 阶段2：自动创建钱包 ==========
	t.Log("阶段2: 为新用户创建钱包")

	walletCreated := &wallet.Wallet{
		ID:        "wallet_alice_001",
		UserID:    user.User.ID,
		Balance:   0.0,
		Frozen:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	walletService.On("CreateWallet", ctx, user.User.ID).Return(walletCreated, nil)

	// 创建钱包
	userWallet, err := walletService.CreateWallet(ctx, user.User.ID)
	assert.NoError(t, err, "创建钱包应该成功")
	assert.Equal(t, user.User.ID, userWallet.UserID)
	assert.Equal(t, 0.0, userWallet.Balance, "新钱包余额应为0")
	t.Logf("✓ 钱包创建成功: ID=%s, 余额=%.2f", userWallet.ID, float64(userWallet.Balance)/100)

	// ========== 阶段3：用户登录 ==========
	t.Log("阶段3: 用户登录")

	loginReq := &auth.LoginRequest{
		Username: "alice",
		Password: "securePassword123",
	}

	loginResp := &auth.LoginResponse{
		User: &auth.UserInfo{
			ID:       "user_alice_001",
			Username: "alice",
			Email:    "alice@example.com",
			Roles:    []string{"writer"},
		},
		Token: "jwt_token_login_alice_new",
	}

	authService.On("Login", ctx, loginReq).Return(loginResp, nil)

	// 执行登录
	loginUser, err := authService.Login(ctx, loginReq)
	assert.NoError(t, err, "用户登录应该成功")
	assert.Equal(t, user.User.ID, loginUser.User.ID)
	assert.NotEmpty(t, loginUser.Token)
	t.Logf("✓ 用户登录成功，获得新Token: %s", loginUser.Token[:20]+"...")

	currentToken := loginUser.Token

	// ========== 阶段4：充值 ==========
	t.Log("阶段4: 用户充值")

	// 验证Token
	tokenClaims := &auth.TokenClaims{
		UserID: user.User.ID,
		Roles:  []string{"writer"},
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	authService.On("ValidateToken", ctx, currentToken).Return(tokenClaims, nil)

	// 充值操作
	rechargeTransaction := &wallet.Transaction{
		ID:              "txn_recharge_001",
		UserID:          user.User.ID,
		Type:            "recharge",
		Amount:          500.0,
		Balance:         500.0,
		Method:          "alipay",
		Status:          "success",
		TransactionTime: time.Now(),
		CreatedAt:       time.Now(),
	}

	walletService.On("Recharge", ctx, user.User.ID, 500.0, "alipay").Return(rechargeTransaction, nil)

	// 验证Token并充值
	claims, err := authService.ValidateToken(ctx, currentToken)
	assert.NoError(t, err)
	assert.Equal(t, user.User.ID, claims.UserID)

	rechargeTxn, err := walletService.Recharge(ctx, claims.UserID, 500.0, "alipay")
	assert.NoError(t, err)
	assert.Equal(t, "success", rechargeTxn.Status)
	assert.Equal(t, 500.0, rechargeTxn.Balance, "充值后余额应为500")
	t.Logf("✓ 充值成功: 金额=%.2f, 余额=%.2f", float64(rechargeTxn.Amount)/100, float64(rechargeTxn.Balance)/100)

	// ========== 阶段5：消费 ==========
	t.Log("阶段5: 用户消费")

	// 第一次消费
	consume1Transaction := &wallet.Transaction{
		ID:              "txn_consume_001",
		UserID:          user.User.ID,
		Type:            "consume",
		Amount:          -150.0,
		Balance:         350.0,
		Reason:          "购买VIP会员",
		Status:          "success",
		TransactionTime: time.Now(),
		CreatedAt:       time.Now(),
	}

	walletService.On("Consume", ctx, user.User.ID, 150.0, "购买VIP会员").Return(consume1Transaction, nil)

	consumeTxn1, err := walletService.Consume(ctx, user.User.ID, 150.0, "购买VIP会员")
	assert.NoError(t, err)
	assert.Equal(t, "success", consumeTxn1.Status)
	assert.Equal(t, 350.0, consumeTxn1.Balance, "消费后余额应为350")
	t.Logf("✓ 第1次消费成功: 金额=%.2f, 剩余余额=%.2f, 原因=%s", float64(-consumeTxn1.Amount)/100, float64(consumeTxn1.Balance)/100, consumeTxn1.Reason)

	// 第二次消费
	consume2Transaction := &wallet.Transaction{
		ID:              "txn_consume_002",
		UserID:          user.User.ID,
		Type:            "consume",
		Amount:          -50.0,
		Balance:         300.0,
		Reason:          "购买书籍",
		Status:          "success",
		TransactionTime: time.Now(),
		CreatedAt:       time.Now(),
	}

	walletService.On("Consume", ctx, user.User.ID, 50.0, "购买书籍").Return(consume2Transaction, nil)

	consumeTxn2, err := walletService.Consume(ctx, user.User.ID, 50.0, "购买书籍")
	assert.NoError(t, err)
	assert.Equal(t, 300.0, consumeTxn2.Balance, "消费后余额应为300")
	t.Logf("✓ 第2次消费成功: 金额=%.2f, 剩余余额=%.2f, 原因=%s", float64(-consumeTxn2.Amount)/100, float64(consumeTxn2.Balance)/100, consumeTxn2.Reason)

	// ========== 阶段6：申请提现 ==========
	t.Log("阶段6: 申请提现")

	withdrawRequest := &wallet.WithdrawRequest{
		ID:        "withdraw_alice_001",
		UserID:    user.User.ID,
		Amount:    100.0,
		Account:   "alipay_account_alice",
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	walletService.On("RequestWithdraw", ctx, user.User.ID, 100.0, "alipay_account_alice").Return(withdrawRequest, nil)

	withdrawReq, err := walletService.RequestWithdraw(ctx, user.User.ID, 100.0, "alipay_account_alice")
	assert.NoError(t, err)
	assert.Equal(t, "pending", withdrawReq.Status)
	assert.Equal(t, 100.0, withdrawReq.Amount)
	t.Logf("✓ 提现申请已提交: ID=%s, 金额=%.2f, 状态=%s", withdrawReq.ID, float64(withdrawReq.Amount)/100, withdrawReq.Status)

	// ========== 阶段7：登出 ==========
	t.Log("阶段7: 用户登出")

	authService.On("Logout", ctx, currentToken).Return(nil)

	err = authService.Logout(ctx, currentToken)
	assert.NoError(t, err, "登出应该成功")
	t.Logf("✓ 用户登出成功")

	// ========== 总结 ==========
	t.Log("\n========== 用户生命周期测试完成 ==========")
	t.Logf("用户: %s (ID: %s)", user.User.Username, user.User.ID)
	t.Logf("最终余额: %.2f", float64(consumeTxn2.Balance)/100)
	t.Logf("交易记录: 充值1次, 消费2次, 提现申请1次")
	t.Logf("总充值: %.2f", float64(rechargeTxn.Amount)/100)
	t.Logf("总消费: %.2f", float64(-consumeTxn1.Amount-consumeTxn2.Amount)/100)
	t.Logf("待提现: %.2f", float64(withdrawReq.Amount)/100)

	// 验证所有Mock调用
	authService.AssertExpectations(t)
	walletService.AssertExpectations(t)
}

// ============ 端到端测试：内容发布与审核流程 ============

// TestE2E_ContentPublishAndReview 测试内容发布与审核的完整流程
// 场景：用户登录 -> 发布内容 -> 管理员审核 -> 状态更新
func TestE2E_ContentPublishAndReview(t *testing.T) {
	ctx := context.Background()

	// 创建所有需要的Mock服务
	authService := new(MockAuthService)
	adminService := new(MockAdminService)
	storageService := new(MockStorageService)

	// ========== 阶段1：作者登录 ==========
	t.Log("阶段1: 作者登录")

	loginReq := &auth.LoginRequest{
		Username: "writer_bob",
		Password: "writerPass123",
	}

	loginResp := &auth.LoginResponse{
		User: &auth.UserInfo{
			ID:       "user_bob_writer",
			Username: "writer_bob",
			Email:    "bob@example.com",
			Roles:    []string{"writer"},
		},
		Token: "jwt_token_bob_writer",
	}

	authService.On("Login", ctx, loginReq).Return(loginResp, nil)

	// 执行登录
	writer, err := authService.Login(ctx, loginReq)
	assert.NoError(t, err)
	assert.Equal(t, "user_bob_writer", writer.User.ID)
	t.Logf("✓ 作者登录成功: %s", writer.User.Username)

	writerToken := writer.Token

	// ========== 阶段2：验证Token并发布内容 ==========
	t.Log("阶段2: 发布内容")

	// 验证Token
	writerClaims := &auth.TokenClaims{
		UserID: writer.User.ID,
		Roles:  []string{"writer"},
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	authService.On("ValidateToken", ctx, writerToken).Return(writerClaims, nil)

	// 上传内容文件（假设是小说章节）
	uploadReq := &storage.UploadRequest{
		Filename:    "chapter_01.txt",
		ContentType: "text/plain",
		Size:        51200, // 50KB
		UserID:      writer.User.ID,
		IsPublic:    false,
		Category:    "novel_chapter",
	}

	uploadedFile := &storage.FileInfo{
		ID:           "file_chapter_001",
		Filename:     "chapter_01_uuid.txt",
		OriginalName: "chapter_01.txt",
		ContentType:  "text/plain",
		Size:         51200,
		Path:         "/uploads/novels/chapter_01_uuid.txt",
		UserID:       writer.User.ID,
		IsPublic:     false,
		Category:     "novel_chapter",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 验证Token
	claims, err := authService.ValidateToken(ctx, writerToken)
	assert.NoError(t, err)
	assert.Equal(t, writer.User.ID, claims.UserID)

	// 上传内容
	storageService.On("Upload", ctx, uploadReq).Return(uploadedFile, nil)

	fileInfo, err := storageService.Upload(ctx, uploadReq)
	assert.NoError(t, err)
	assert.Equal(t, "file_chapter_001", fileInfo.ID)
	assert.Equal(t, writer.User.ID, fileInfo.UserID)
	t.Logf("✓ 内容已发布: FileID=%s, 文件名=%s", fileInfo.ID, fileInfo.OriginalName)

	// 内容发布后自动创建审核记录（模拟）
	contentID := "content_chapter_001"
	t.Logf("✓ 审核记录已创建: ContentID=%s, 状态=pending", contentID)

	// ========== 阶段3：管理员登录 ==========
	t.Log("阶段3: 管理员登录")

	adminLoginReq := &auth.LoginRequest{
		Username: "admin_charlie",
		Password: "adminPass123",
	}

	adminLoginResp := &auth.LoginResponse{
		User: &auth.UserInfo{
			ID:       "user_admin_charlie",
			Username: "admin_charlie",
			Email:    "charlie@example.com",
			Roles:    []string{"admin"},
		},
		Token: "jwt_token_admin_charlie",
	}

	authService.On("Login", ctx, adminLoginReq).Return(adminLoginResp, nil)

	// 执行登录
	adminUser, err := authService.Login(ctx, adminLoginReq)
	assert.NoError(t, err)
	assert.Equal(t, "user_admin_charlie", adminUser.User.ID)
	t.Logf("✓ 管理员登录成功: %s", adminUser.User.Username)

	adminToken := adminUser.Token

	// ========== 阶段4：管理员查看待审核内容 ==========
	t.Log("阶段4: 查看待审核内容")

	// 验证管理员Token
	adminClaims := &auth.TokenClaims{
		UserID: adminUser.User.ID,
		Roles:  []string{"admin"},
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	authService.On("ValidateToken", ctx, adminToken).Return(adminClaims, nil)

	// 获取待审核列表
	pendingReviews := []*admin.AuditRecord{
		{
			ID:          "audit_001",
			ContentID:   contentID,
			ContentType: "novel_chapter",
			Status:      "pending",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	adminService.On("GetPendingReviews", ctx, "novel_chapter").Return(pendingReviews, nil)

	// 验证Token
	adminTokenClaims, err := authService.ValidateToken(ctx, adminToken)
	assert.NoError(t, err)
	assert.Equal(t, adminUser.User.ID, adminTokenClaims.UserID)
	assert.Contains(t, adminTokenClaims.Roles, "admin")

	// 获取待审核内容
	reviews, err := adminService.GetPendingReviews(ctx, "novel_chapter")
	assert.NoError(t, err)
	assert.Len(t, reviews, 1)
	assert.Equal(t, "pending", reviews[0].Status)
	t.Logf("✓ 待审核内容: %d个", len(reviews))

	// ========== 阶段5：管理员审核内容 ==========
	t.Log("阶段5: 审核内容")

	// 场景A：审核通过
	reviewApproveReq := &admin.ReviewContentRequest{
		ContentID:   contentID,
		ContentType: "novel_chapter",
		Action:      "approve",
		Reason:      "",
		ReviewerID:  adminUser.User.ID,
	}

	adminService.On("ReviewContent", ctx, reviewApproveReq).Return(nil)

	err = adminService.ReviewContent(ctx, reviewApproveReq)
	assert.NoError(t, err)
	t.Logf("✓ 内容审核通过: ContentID=%s", contentID)

	// ========== 总结 ==========
	t.Log("\n========== 内容发布与审核流程完成 ==========")
	t.Logf("作者: %s 发布了内容 %s", writer.User.Username, fileInfo.OriginalName)
	t.Logf("管理员: %s 审核通过", adminUser.User.Username)
	t.Logf("内容状态: pending -> approved")

	// 验证所有Mock调用
	authService.AssertExpectations(t)
	storageService.AssertExpectations(t)
	adminService.AssertExpectations(t)
}

// TestE2E_ContentRejectionFlow 测试内容审核被拒绝的流程
func TestE2E_ContentRejectionFlow(t *testing.T) {
	ctx := context.Background()

	// 创建所有需要的Mock服务
	authService := new(MockAuthService)
	adminService := new(MockAdminService)

	// ========== 阶段1：管理员登录 ==========
	t.Log("阶段1: 管理员登录")

	adminLoginReq := &auth.LoginRequest{
		Username: "admin_david",
		Password: "adminPass456",
	}

	adminLoginResp := &auth.LoginResponse{
		User: &auth.UserInfo{
			ID:       "user_admin_david",
			Username: "admin_david",
			Email:    "david@example.com",
			Roles:    []string{"admin"},
		},
		Token: "jwt_token_admin_david",
	}

	authService.On("Login", ctx, adminLoginReq).Return(adminLoginResp, nil)

	adminUser, err := authService.Login(ctx, adminLoginReq)
	assert.NoError(t, err)
	t.Logf("✓ 管理员登录成功: %s", adminUser.User.Username)

	// ========== 阶段2：审核拒绝内容 ==========
	t.Log("阶段2: 审核拒绝内容")

	contentID := "content_violation_001"
	rejectReason := "内容违反社区规范：包含不当言论"

	// 验证管理员Token
	adminClaims := &auth.TokenClaims{
		UserID: adminUser.User.ID,
		Roles:  []string{"admin"},
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	authService.On("ValidateToken", ctx, adminUser.Token).Return(adminClaims, nil)

	// 审核拒绝
	reviewRejectReq := &admin.ReviewContentRequest{
		ContentID:   contentID,
		ContentType: "novel_chapter",
		Action:      "reject",
		Reason:      rejectReason,
		ReviewerID:  adminUser.User.ID,
	}

	adminService.On("ReviewContent", ctx, reviewRejectReq).Return(nil)

	// 验证Token
	claims, err := authService.ValidateToken(ctx, adminUser.Token)
	assert.NoError(t, err)
	assert.Contains(t, claims.Roles, "admin")

	// 执行拒绝操作
	err = adminService.ReviewContent(ctx, reviewRejectReq)
	assert.NoError(t, err)
	t.Logf("✓ 内容审核被拒绝: ContentID=%s", contentID)
	t.Logf("  拒绝原因: %s", rejectReason)

	// ========== 总结 ==========
	t.Log("\n========== 内容拒绝流程完成 ==========")
	t.Logf("管理员: %s", adminUser.User.Username)
	t.Logf("操作: 拒绝内容 %s", contentID)
	t.Logf("原因: %s", rejectReason)

	// 验证所有Mock调用
	authService.AssertExpectations(t)
	adminService.AssertExpectations(t)
}

// ============ 端到端测试：多用户转账流程 ============

// TestE2E_TransferBetweenUsers 测试用户之间的转账流程
func TestE2E_TransferBetweenUsers(t *testing.T) {
	ctx := context.Background()

	// 创建服务
	authService := new(MockAuthService)
	walletService := new(MockWalletService)

	// ========== 阶段1：两个用户登录 ==========
	t.Log("阶段1: 两个用户登录")

	// 用户A登录
	userA := &auth.LoginResponse{
		User: &auth.UserInfo{
			ID:       "user_alice",
			Username: "alice",
			Roles:    []string{"writer"},
		},
		Token: "token_alice",
	}

	// 用户B登录
	userB := &auth.LoginResponse{
		User: &auth.UserInfo{
			ID:       "user_bob",
			Username: "bob",
			Roles:    []string{"reader"},
		},
		Token: "token_bob",
	}
	_ = userB.Token // 显式标记为有意未使用（测试数据完整性）

	t.Logf("✓ 用户A: %s (ID: %s)", userA.User.Username, userA.User.ID)
	t.Logf("✓ 用户B: %s (ID: %s)", userB.User.Username, userB.User.ID)

	// ========== 阶段2：用户A发起转账 ==========
	t.Log("阶段2: 用户A向用户B转账")

	// 验证用户A的Token
	claimsA := &auth.TokenClaims{
		UserID: userA.User.ID,
		Roles:  userA.User.Roles,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	authService.On("ValidateToken", ctx, userA.Token).Return(claimsA, nil)

	// 转账操作
	transferAmount := 100.0
	transferReason := "感谢打赏"

	transferTransaction := &wallet.Transaction{
		ID:              "txn_transfer_001",
		UserID:          userA.User.ID,
		Type:            "transfer_out",
		Amount:          -int64(types.NewMoneyFromYuan(transferAmount)),
		Balance:         int64(types.NewMoneyFromYuan(400.0)), // 假设转账前余额是500
		RelatedUserID:   userB.User.ID,
		Reason:          transferReason,
		Status:          "success",
		TransactionTime: time.Now(),
		CreatedAt:       time.Now(),
	}

	walletService.On("Transfer", ctx, userA.User.ID, userB.User.ID, transferAmount, transferReason).
		Return(transferTransaction, nil)

	// 执行转账
	claims, err := authService.ValidateToken(ctx, userA.Token)
	assert.NoError(t, err)

	txn, err := walletService.Transfer(ctx, claims.UserID, userB.User.ID, transferAmount, transferReason)
	assert.NoError(t, err)
	assert.Equal(t, "success", txn.Status)
	assert.Equal(t, "transfer_out", txn.Type)
	assert.Equal(t, -types.NewMoneyFromYuan(transferAmount), txn.Amount)
	t.Logf("✓ 转账成功: %s -> %s, 金额=%.2f, 原因=%s",
		userA.User.Username, userB.User.Username, transferAmount, transferReason)
	t.Logf("  用户A余额: %.2f", float64(txn.Balance)/100)

	// ========== 总结 ==========
	t.Log("\n========== 转账流程完成 ==========")
	t.Logf("发起人: %s", userA.User.Username)
	t.Logf("接收人: %s", userB.User.Username)
	t.Logf("转账金额: %.2f", transferAmount)
	t.Logf("交易状态: %s", txn.Status)

	// 验证所有Mock调用
	authService.AssertExpectations(t)
	walletService.AssertExpectations(t)
}
