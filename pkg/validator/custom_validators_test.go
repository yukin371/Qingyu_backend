package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateAmount(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{"正数金额", 100.00, false},
		{"小数金额", 50.50, false},
		{"最小金额", 0.01, false},
		{"零金额", 0.00, false},
		{"负数金额", -10.00, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Amount float64 `validate:"amount"`
			}
			s := TestStruct{Amount: tt.amount}
			err := v.Struct(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePositiveAmount(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{"正数", 100.00, false},
		{"最小正数", 0.01, false},
		{"零", 0.00, true},
		{"负数", -10.00, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Amount float64 `validate:"positive_amount"`
			}
			s := TestStruct{Amount: tt.amount}
			err := v.Struct(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAmountRange(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{"最小值", 0.01, false},
		{"中间值", 500.00, false},
		{"最大值", 1000000.00, false},
		{"小于最小值", 0.001, true},
		{"大于最大值", 1000001.00, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Amount float64 `validate:"amount_range"`
			}
			s := TestStruct{Amount: tt.amount}
			err := v.Struct(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"有效用户名", "user123", false},
		{"下划线", "user_name", false},
		{"最短", "abc", false},
		{"最长", "12345678901234567890", false},
		{"太短", "ab", true},
		{"太长", "123456789012345678901", true},
		{"包含特殊字符", "user@123", true},
		{"包含空格", "user 123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Username string `validate:"username"`
			}
			s := TestStruct{Username: tt.username}
			err := v.Struct(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePhone(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"有效手机号", "13812345678", false},
		{"移动", "15912345678", false},
		{"电信", "18912345678", false},
		{"联通", "13012345678", false},
		{"太短", "1381234567", true},
		{"太长", "138123456789", true},
		{"错误前缀", "12812345678", true},
		{"包含字母", "138abcd5678", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Phone string `validate:"phone"`
			}
			s := TestStruct{Phone: tt.phone}
			err := v.Struct(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateStrongPassword(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"强密码", "Password123", false},
		{"最短强密码", "Pass123", true}, // 少于8位
		{"只有小写", "password123", true},
		{"只有大写", "PASSWORD123", true},
		{"只有字母", "PasswordABC", true},
		{"包含特殊字符", "Pass@123", false},
		{"超长强密码", "Password123456789", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Password string `validate:"strong_password"`
			}
			s := TestStruct{Password: tt.password}
			err := v.Struct(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTransactionType(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	tests := []struct {
		name    string
		txType  string
		wantErr bool
	}{
		{"充值", "recharge", false},
		{"消费", "consume", false},
		{"转账", "transfer", false},
		{"退款", "refund", false},
		{"提现", "withdraw", false},
		{"无效类型", "invalid", true},
		{"空字符串", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Type string `validate:"transaction_type"`
			}
			s := TestStruct{Type: tt.txType}
			err := v.Struct(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateWithdrawAccount(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	tests := []struct {
		name    string
		account string
		wantErr bool
	}{
		{"支付宝", "alipay:user@example.com", false},
		{"微信", "wechat:wxid_123456", false},
		{"银行卡", "bank:6222021234567890", false},
		{"无效格式", "alipayuser@example.com", true},
		{"无效方式", "paypal:user@example.com", true},
		{"空账号", "alipay:", true},
		{"空字符串", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Account string `validate:"withdraw_account"`
			}
			s := TestStruct{Account: tt.account}
			err := v.Struct(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateContentType(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	tests := []struct {
		name        string
		contentType string
		wantErr     bool
	}{
		{"书籍", "book", false},
		{"章节", "chapter", false},
		{"评论", "comment", false},
		{"评价", "review", false},
		{"无效类型", "article", true},
		{"空字符串", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Type string `validate:"content_type"`
			}
			s := TestStruct{Type: tt.contentType}
			err := v.Struct(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// 测试综合场景
func TestComplexValidation(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	type WalletRequest struct {
		Amount  float64 `validate:"positive_amount,amount_range"`
		Account string  `validate:"withdraw_account"`
	}

	tests := []struct {
		name    string
		req     WalletRequest
		wantErr bool
	}{
		{
			name:    "有效请求",
			req:     WalletRequest{Amount: 100.00, Account: "alipay:user@example.com"},
			wantErr: false,
		},
		{
			name:    "金额无效",
			req:     WalletRequest{Amount: -10.00, Account: "alipay:user@example.com"},
			wantErr: true,
		},
		{
			name:    "账号无效",
			req:     WalletRequest{Amount: 100.00, Account: "invalid"},
			wantErr: true,
		},
		{
			name:    "全部无效",
			req:     WalletRequest{Amount: 0.00, Account: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
