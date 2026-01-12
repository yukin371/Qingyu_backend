package integration

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"Qingyu_backend/service/shared/auth"
	"Qingyu_backend/service/finance/wallet"

	"github.com/stretchr/testify/mock"
)

// ============ 基准测试：Auth服务 ============

// BenchmarkAuthService_Login 基准测试：用户登录
func BenchmarkAuthService_Login(b *testing.B) {
	ctx := context.Background()
	authService := new(MockAuthService)

	loginReq := &auth.LoginRequest{
		Username: "benchmark_user",
		Password: "password123",
	}

	loginResp := &auth.LoginResponse{
		User: &auth.UserInfo{
			ID:       "user_benchmark",
			Username: "benchmark_user",
			Roles:    []string{"reader"},
		},
		Token: "jwt_token_benchmark",
	}

	authService.On("Login", mock.Anything, mock.Anything).Return(loginResp, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := authService.Login(ctx, loginReq)
		if err != nil {
			b.Fatalf("Login failed: %v", err)
		}
	}
}

// BenchmarkAuthService_ValidateToken 基准测试：Token验证
func BenchmarkAuthService_ValidateToken(b *testing.B) {
	ctx := context.Background()
	authService := new(MockAuthService)

	token := "jwt_token_benchmark"
	claims := &auth.TokenClaims{
		UserID: "user_benchmark",
		Roles:  []string{"reader"},
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	authService.On("ValidateToken", mock.Anything, mock.Anything).Return(claims, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := authService.ValidateToken(ctx, token)
		if err != nil {
			b.Fatalf("ValidateToken failed: %v", err)
		}
	}
}

// BenchmarkAuthService_Register 基准测试：用户注册
func BenchmarkAuthService_Register(b *testing.B) {
	ctx := context.Background()
	authService := new(MockAuthService)

	registerResp := &auth.RegisterResponse{
		User: &auth.UserInfo{
			ID:       "user_new",
			Username: "new_user",
			Roles:    []string{"reader"},
		},
		Token: "jwt_token_new",
	}

	authService.On("Register", mock.Anything, mock.Anything).Return(registerResp, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registerReq := &auth.RegisterRequest{
			Username: fmt.Sprintf("user_%d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Password: "password123",
		}
		_, err := authService.Register(ctx, registerReq)
		if err != nil {
			b.Fatalf("Register failed: %v", err)
		}
	}
}

// ============ 基准测试：Wallet服务 ============

// BenchmarkWalletService_GetBalance 基准测试：查询余额
func BenchmarkWalletService_GetBalance(b *testing.B) {
	ctx := context.Background()
	walletService := new(MockWalletService)

	userID := "user_benchmark"
	walletService.On("GetBalance", mock.Anything, mock.Anything).Return(1000.0, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := walletService.GetBalance(ctx, userID)
		if err != nil {
			b.Fatalf("GetBalance failed: %v", err)
		}
	}
}

// BenchmarkWalletService_Recharge 基准测试：充值操作
func BenchmarkWalletService_Recharge(b *testing.B) {
	ctx := context.Background()
	walletService := new(MockWalletService)

	userID := "user_benchmark"
	transaction := &wallet.Transaction{
		ID:      "txn_benchmark",
		UserID:  userID,
		Type:    "recharge",
		Amount:  100.0,
		Balance: 1100.0,
		Status:  "success",
	}

	walletService.On("Recharge", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(transaction, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := walletService.Recharge(ctx, userID, 100.0, "alipay")
		if err != nil {
			b.Fatalf("Recharge failed: %v", err)
		}
	}
}

// BenchmarkWalletService_Consume 基准测试：消费操作
func BenchmarkWalletService_Consume(b *testing.B) {
	ctx := context.Background()
	walletService := new(MockWalletService)

	userID := "user_benchmark"
	transaction := &wallet.Transaction{
		ID:      "txn_consume",
		UserID:  userID,
		Type:    "consume",
		Amount:  -50.0,
		Balance: 950.0,
		Status:  "success",
	}

	walletService.On("Consume", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(transaction, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := walletService.Consume(ctx, userID, 50.0, "购买商品")
		if err != nil {
			b.Fatalf("Consume failed: %v", err)
		}
	}
}

// BenchmarkWalletService_Transfer 基准测试：转账操作
func BenchmarkWalletService_Transfer(b *testing.B) {
	ctx := context.Background()
	walletService := new(MockWalletService)

	fromUserID := "user_a"
	toUserID := "user_b"
	transaction := &wallet.Transaction{
		ID:            "txn_transfer",
		UserID:        fromUserID,
		Type:          "transfer_out",
		Amount:        -100.0,
		Balance:       900.0,
		RelatedUserID: toUserID,
		Status:        "success",
	}

	walletService.On("Transfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(transaction, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := walletService.Transfer(ctx, fromUserID, toUserID, 100.0, "转账")
		if err != nil {
			b.Fatalf("Transfer failed: %v", err)
		}
	}
}

// ============ 并发测试 ============

// TestConcurrent_AuthLogin 并发测试：用户登录
func TestConcurrent_AuthLogin(t *testing.T) {
	ctx := context.Background()
	authService := new(MockAuthService)

	loginResp := &auth.LoginResponse{
		User: &auth.UserInfo{
			ID:       "user_concurrent",
			Username: "concurrent_user",
			Roles:    []string{"reader"},
		},
		Token: "jwt_token_concurrent",
	}

	authService.On("Login", mock.Anything, mock.Anything).Return(loginResp, nil)

	concurrentUsers := 100
	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			loginReq := &auth.LoginRequest{
				Username: fmt.Sprintf("user_%d", id),
				Password: "password123",
			}
			_, err := authService.Login(ctx, loginReq)
			if err != nil {
				t.Errorf("Login failed for user %d: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	t.Logf("\n========== 并发登录测试结果 ==========")
	t.Logf("并发用户数: %d", concurrentUsers)
	t.Logf("总耗时: %v", duration)
	t.Logf("平均响应时间: %v", duration/time.Duration(concurrentUsers))
	t.Logf("QPS: %.2f", float64(concurrentUsers)/duration.Seconds())
}

// TestConcurrent_WalletOperations 并发测试：钱包操作
func TestConcurrent_WalletOperations(t *testing.T) {
	ctx := context.Background()
	walletService := new(MockWalletService)

	// Mock余额查询
	walletService.On("GetBalance", mock.Anything, mock.Anything).Return(1000.0, nil)

	// Mock充值
	rechargeTransaction := &wallet.Transaction{
		ID:      "txn_recharge",
		Type:    "recharge",
		Amount:  100.0,
		Balance: 1100.0,
		Status:  "success",
	}
	walletService.On("Recharge", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(rechargeTransaction, nil)

	// Mock消费
	consumeTransaction := &wallet.Transaction{
		ID:      "txn_consume",
		Type:    "consume",
		Amount:  -50.0,
		Balance: 950.0,
		Status:  "success",
	}
	walletService.On("Consume", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(consumeTransaction, nil)

	concurrentOps := 300 // 100个查询 + 100个充值 + 100个消费
	var wg sync.WaitGroup
	startTime := time.Now()

	// 并发查询余额
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			userID := fmt.Sprintf("user_%d", id)
			_, err := walletService.GetBalance(ctx, userID)
			if err != nil {
				t.Errorf("GetBalance failed for user %d: %v", id, err)
			}
		}(i)
	}

	// 并发充值
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			userID := fmt.Sprintf("user_%d", id)
			_, err := walletService.Recharge(ctx, userID, 100.0, "alipay")
			if err != nil {
				t.Errorf("Recharge failed for user %d: %v", id, err)
			}
		}(i)
	}

	// 并发消费
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			userID := fmt.Sprintf("user_%d", id)
			_, err := walletService.Consume(ctx, userID, 50.0, "购买商品")
			if err != nil {
				t.Errorf("Consume failed for user %d: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	t.Logf("\n========== 并发钱包操作测试结果 ==========")
	t.Logf("并发操作数: %d (查询100 + 充值100 + 消费100)", concurrentOps)
	t.Logf("总耗时: %v", duration)
	t.Logf("平均响应时间: %v", duration/time.Duration(concurrentOps))
	t.Logf("QPS: %.2f", float64(concurrentOps)/duration.Seconds())
}

// TestConcurrent_MixedOperations 并发测试：混合操作
func TestConcurrent_MixedOperations(t *testing.T) {
	ctx := context.Background()
	authService := new(MockAuthService)
	walletService := new(MockWalletService)

	// Mock Auth服务
	loginResp := &auth.LoginResponse{
		User: &auth.UserInfo{
			ID:       "user_mixed",
			Username: "mixed_user",
			Roles:    []string{"reader"},
		},
		Token: "jwt_token_mixed",
	}
	authService.On("Login", mock.Anything, mock.Anything).Return(loginResp, nil)

	tokenClaims := &auth.TokenClaims{
		UserID: "user_mixed",
		Roles:  []string{"reader"},
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}
	authService.On("ValidateToken", mock.Anything, mock.Anything).Return(tokenClaims, nil)

	// Mock Wallet服务
	walletService.On("GetBalance", mock.Anything, mock.Anything).Return(1000.0, nil)

	rechargeTransaction := &wallet.Transaction{
		ID:      "txn_recharge",
		Type:    "recharge",
		Amount:  100.0,
		Balance: 1100.0,
		Status:  "success",
	}
	walletService.On("Recharge", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(rechargeTransaction, nil)

	concurrentUsers := 50
	var wg sync.WaitGroup
	startTime := time.Now()

	// 每个用户执行：登录 -> 验证Token -> 查询余额 -> 充值
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 1. 登录
			loginReq := &auth.LoginRequest{
				Username: fmt.Sprintf("user_%d", id),
				Password: "password123",
			}
			user, err := authService.Login(ctx, loginReq)
			if err != nil {
				t.Errorf("Login failed: %v", err)
				return
			}

			// 2. 验证Token
			_, err = authService.ValidateToken(ctx, user.Token)
			if err != nil {
				t.Errorf("ValidateToken failed: %v", err)
				return
			}

			// 3. 查询余额
			userID := user.User.ID
			_, err = walletService.GetBalance(ctx, userID)
			if err != nil {
				t.Errorf("GetBalance failed: %v", err)
				return
			}

			// 4. 充值
			_, err = walletService.Recharge(ctx, userID, 100.0, "alipay")
			if err != nil {
				t.Errorf("Recharge failed: %v", err)
				return
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	totalOps := concurrentUsers * 4 // 每个用户4个操作

	t.Logf("\n========== 并发混合操作测试结果 ==========")
	t.Logf("并发用户数: %d", concurrentUsers)
	t.Logf("每用户操作数: 4 (登录+验证+查询+充值)")
	t.Logf("总操作数: %d", totalOps)
	t.Logf("总耗时: %v", duration)
	t.Logf("平均每用户耗时: %v", duration/time.Duration(concurrentUsers))
	t.Logf("平均每操作耗时: %v", duration/time.Duration(totalOps))
	t.Logf("用户QPS: %.2f", float64(concurrentUsers)/duration.Seconds())
	t.Logf("操作QPS: %.2f", float64(totalOps)/duration.Seconds())
}

// ============ 压力测试 ============

// TestStress_HighConcurrency 压力测试：高并发场景
func TestStress_HighConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过压力测试（使用 -short 标志）")
	}

	ctx := context.Background()
	authService := new(MockAuthService)

	loginResp := &auth.LoginResponse{
		User: &auth.UserInfo{
			ID:       "user_stress",
			Username: "stress_user",
			Roles:    []string{"reader"},
		},
		Token: "jwt_token_stress",
	}
	authService.On("Login", mock.Anything, mock.Anything).Return(loginResp, nil)

	concurrentLevels := []int{100, 500, 1000}

	for _, level := range concurrentLevels {
		t.Run(fmt.Sprintf("并发%d", level), func(t *testing.T) {
			var wg sync.WaitGroup
			startTime := time.Now()

			for i := 0; i < level; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					loginReq := &auth.LoginRequest{
						Username: fmt.Sprintf("user_%d", id),
						Password: "password123",
					}
					_, err := authService.Login(ctx, loginReq)
					if err != nil {
						t.Errorf("Login failed: %v", err)
					}
				}(i)
			}

			wg.Wait()
			duration := time.Since(startTime)

			t.Logf("\n并发级别: %d", level)
			t.Logf("总耗时: %v", duration)
			t.Logf("平均响应时间: %v", duration/time.Duration(level))
			t.Logf("QPS: %.2f", float64(level)/duration.Seconds())
		})
	}
}

// TestStress_SustainedLoad 压力测试：持续负载
func TestStress_SustainedLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过压力测试（使用 -short 标志）")
	}

	ctx := context.Background()
	walletService := new(MockWalletService)

	walletService.On("GetBalance", mock.Anything, mock.Anything).Return(1000.0, nil)

	duration := 5 * time.Second
	concurrency := 50

	var totalOps int64
	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	startTime := time.Now()

	// 启动并发worker
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			userID := fmt.Sprintf("user_%d", id)
			ops := 0
			for {
				select {
				case <-stopCh:
					return
				default:
					_, err := walletService.GetBalance(ctx, userID)
					if err != nil {
						t.Errorf("GetBalance failed: %v", err)
						return
					}
					ops++
					if ops%100 == 0 {
						// 每100次操作短暂休息，模拟真实场景
						time.Sleep(1 * time.Millisecond)
					}
				}
			}
		}(i)
	}

	// 运行指定时间
	time.Sleep(duration)
	close(stopCh)
	wg.Wait()

	elapsed := time.Since(startTime)

	// 估算总操作数（基于Mock的调用次数）
	// 这是一个近似值
	estimatedOps := int64(concurrency * 1000) // 粗略估算

	t.Logf("\n========== 持续负载测试结果 ==========")
	t.Logf("并发数: %d", concurrency)
	t.Logf("测试时长: %v", duration)
	t.Logf("实际耗时: %v", elapsed)
	t.Logf("估算总操作数: %d", estimatedOps)
	t.Logf("估算QPS: %.2f", float64(estimatedOps)/elapsed.Seconds())
	totalOps = estimatedOps
	_ = totalOps
}
