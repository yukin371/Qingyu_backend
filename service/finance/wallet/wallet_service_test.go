package wallet

import (
	"context"
	"testing"

	financeModel "Qingyu_backend/models/finance"
	"Qingyu_backend/models/shared/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalletService_CreateWallet(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		setup   func(*MockWalletRepositoryV2)
		wantErr bool
		errMsg  string
	}{
		{
			name:   "正常创建钱包",
			userID: "user123",
			setup: func(m *MockWalletRepositoryV2) {
				// 钱包不存在，可以创建
			},
			wantErr: false,
		},
		{
			name:   "用户已有钱包",
			userID: "user123",
			setup: func(m *MockWalletRepositoryV2) {
				m.wallets["user123"] = &financeModel.Wallet{
					UserID:  "user123",
					Balance: types.Money(1000),
				}
			},
			wantErr: true,
			errMsg:  "用户已有钱包",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockWalletRepositoryV2()
			tt.setup(repo)

			svc := NewWalletService(repo)
			result, err := svc.CreateWallet(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.userID, result.UserID)
				assert.Equal(t, int64(0), result.Balance)
				assert.False(t, result.Frozen)
			}
		})
	}
}

func TestWalletService_GetWallet(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		setup   func(*MockWalletRepositoryV2)
		wantErr bool
		errMsg  string
	}{
		{
			name:   "正常获取钱包",
			userID: "user123",
			setup: func(m *MockWalletRepositoryV2) {
				m.wallets["user123"] = &financeModel.Wallet{
					UserID:  "user123",
					Balance: types.Money(1000),
				}
			},
			wantErr: false,
		},
		{
			name:   "钱包不存在",
			userID: "user_not_exist",
			setup: func(m *MockWalletRepositoryV2) {
				// 不创建钱包
			},
			wantErr: true,
			errMsg:  "钱包不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockWalletRepositoryV2()
			tt.setup(repo)

			svc := NewWalletService(repo)
			result, err := svc.GetWallet(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.userID, result.UserID)
			}
		})
	}
}

func TestWalletService_GetBalance(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		setup      func(*MockWalletRepositoryV2)
		wantErr    bool
		errMsg     string
		wantAmount int64
	}{
		{
			name:   "正常获取余额",
			userID: "user123",
			setup: func(m *MockWalletRepositoryV2) {
				m.wallets["user123"] = &financeModel.Wallet{
					UserID:  "user123",
					Balance: types.Money(1000),
				}
			},
			wantErr:    false,
			wantAmount: 1000,
		},
		{
			name:   "钱包不存在",
			userID: "user_not_exist",
			setup: func(m *MockWalletRepositoryV2) {
				// 不创建钱包
			},
			wantErr: true,
			errMsg:  "钱包不存在",
		},
		{
			name:   "零余额",
			userID: "user_zero",
			setup: func(m *MockWalletRepositoryV2) {
				m.wallets["user_zero"] = &financeModel.Wallet{
					UserID:  "user_zero",
					Balance: types.Money(0),
				}
			},
			wantErr:    false,
			wantAmount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockWalletRepositoryV2()
			tt.setup(repo)

			svc := NewWalletService(repo)
			balance, err := svc.GetBalance(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantAmount, balance)
			}
		})
	}
}

func TestWalletService_FreezeWallet(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		reason  string
		setup   func(*MockWalletRepositoryV2)
		wantErr bool
		errMsg  string
	}{
		{
			name:   "正常冻结钱包",
			userID: "user123",
			reason: "违规操作",
			setup: func(m *MockWalletRepositoryV2) {
				m.wallets["user123"] = &financeModel.Wallet{
					UserID:  "user123",
					Balance: types.Money(1000),
					Frozen:  false,
				}
			},
			wantErr: false,
		},
		{
			name:   "钱包不存在",
			userID: "user_not_exist",
			reason: "测试",
			setup: func(m *MockWalletRepositoryV2) {
				// 不创建钱包
			},
			wantErr: true,
			errMsg:  "钱包不存在",
		},
		{
			name:   "已冻结的钱包再次冻结",
			userID: "user_frozen",
			reason: "再次冻结",
			setup: func(m *MockWalletRepositoryV2) {
				m.wallets["user_frozen"] = &financeModel.Wallet{
					UserID:  "user_frozen",
					Balance: types.Money(1000),
					Frozen:  true,
				}
			},
			wantErr: false, // 冻结操作是幂等的
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockWalletRepositoryV2()
			tt.setup(repo)

			svc := NewWalletService(repo)
			err := svc.FreezeWallet(context.Background(), tt.userID, tt.reason)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				// 验证钱包状态已更新
				wallet := repo.wallets[tt.userID]
				assert.True(t, wallet.Frozen)
			}
		})
	}
}

func TestWalletService_UnfreezeWallet(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		setup   func(*MockWalletRepositoryV2)
		wantErr bool
		errMsg  string
	}{
		{
			name:   "正常解冻钱包",
			userID: "user123",
			setup: func(m *MockWalletRepositoryV2) {
				m.wallets["user123"] = &financeModel.Wallet{
					UserID:  "user123",
					Balance: types.Money(1000),
					Frozen:  true,
				}
			},
			wantErr: false,
		},
		{
			name:   "钱包不存在",
			userID: "user_not_exist",
			setup: func(m *MockWalletRepositoryV2) {
				// 不创建钱包
			},
			wantErr: true,
			errMsg:  "钱包不存在",
		},
		{
			name:   "未冻结的钱包解冻",
			userID: "user_unfrozen",
			setup: func(m *MockWalletRepositoryV2) {
				m.wallets["user_unfrozen"] = &financeModel.Wallet{
					UserID:  "user_unfrozen",
					Balance: types.Money(1000),
					Frozen:  false,
				}
			},
			wantErr: false, // 解冻操作是幂等的
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockWalletRepositoryV2()
			tt.setup(repo)

			svc := NewWalletService(repo)
			err := svc.UnfreezeWallet(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				// 验证钱包状态已更新
				wallet := repo.wallets[tt.userID]
				assert.False(t, wallet.Frozen)
			}
		})
	}
}

func TestWalletService_GetWalletByID(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	svc := NewWalletService(repo)

	// GetWalletByID 当前返回未实现错误
	_, err := svc.GetWalletByID(context.Background(), "wallet123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "未实现")
}

// TestNewWalletService 测试构造函数
func TestNewWalletService(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	svc := NewWalletService(repo)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.walletRepo)
}

// TestConvertToWalletResponse 测试转换函数
func TestConvertToWalletResponse(t *testing.T) {
	wallet := &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(1000),
		Frozen:  true,
	}

	result := convertToWalletResponse(wallet)

	assert.NotNil(t, result)
	assert.Equal(t, "user123", result.UserID)
	assert.Equal(t, int64(1000), result.Balance)
	assert.True(t, result.Frozen)
	assert.NotEmpty(t, result.ID)
}
