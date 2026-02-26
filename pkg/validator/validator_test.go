package validator

import (
	"testing"
)

// TestValidatorInitialization 测试验证器初始化
func TestValidatorInitialization(t *testing.T) {
	v := GetValidator()

	if v == nil {
		t.Fatal("GetValidator() returned nil")
	}

	// 检查注册状态
	status := GetRegistrationStatus()
	if status.Total == 0 {
		t.Error("No validators were registered")
	}

	t.Logf("Validator registration status: %d/%d successful", status.Success, status.Total)
}

// TestRegistrationStatus 测试注册状态
func TestRegistrationStatus(t *testing.T) {
	status := GetRegistrationStatus()

	if status.Total != status.Success+status.Failed {
		t.Errorf("Total count mismatch: total=%d, success=%d, failed=%d",
			status.Total, status.Success, status.Failed)
	}

	if status.IsComplete() && status.Failed > 0 {
		t.Error("IsComplete() returned true when there were failures")
	}

	if !status.IsComplete() && status.Failed == 0 {
		t.Error("IsComplete() returned false when there were no failures")
	}

	t.Logf("Registration status: Total=%d, Success=%d, Failed=%d",
		status.Total, status.Success, status.Failed)

	if status.Failed > 0 {
		t.Logf("Failed validators: %v", status.FailedTags)
	}
}

// TestAmountValidation 使用实际验证器测试金额验证
func TestAmountValidation(t *testing.T) {
	type TestCase struct {
		Name    string
		Amount  float64 `validate:"amount"`
		Valid   bool
	}

	tests := []TestCase{
		{"valid positive amount", 100.50, true},
		{"valid zero", 0, true},
		{"negative amount", -50.00, false},
		{"valid with 2 decimals", 100.99, true},
		{"valid with 1 decimal", 100.5, true},
		{"valid without decimals", 100, true},
	}

	v := GetValidator()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.Struct(tt)
			if tt.Valid && err != nil {
				t.Errorf("Expected valid but got error: %v", err)
			}
			if !tt.Valid && err == nil {
				t.Error("Expected invalid but got no error")
			}
		})
	}
}

// TestPositiveAmountValidation 使用实际验证器测试正数金额验证
func TestPositiveAmountValidation(t *testing.T) {
	type TestCase struct {
		Name   string
		Amount float64 `validate:"positive_amount"`
		Valid  bool
	}

	tests := []TestCase{
		{"valid positive amount", 100.50, true},
		{"positive small amount", 0.01, true},
		{"zero", 0, false},
		{"negative amount", -50.00, false},
	}

	v := GetValidator()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.Struct(tt)
			if tt.Valid && err != nil {
				t.Errorf("Expected valid but got error: %v", err)
			}
			if !tt.Valid && err == nil {
				t.Error("Expected invalid but got no error")
			}
		})
	}
}

// TestAmountRangeValidation 使用实际验证器测试金额范围验证
func TestAmountRangeValidation(t *testing.T) {
	type TestCase struct {
		Name   string
		Amount float64 `validate:"amount_range"`
		Valid  bool
	}

	tests := []TestCase{
		{"valid amount", 100.50, true},
		{"minimum valid", 0.01, true},
		{"maximum valid", 1000000.00, true},
		{"below minimum", 0.001, false},
		{"above maximum", 1000000.01, false},
		{"zero", 0, false},
		{"negative", -100, false},
	}

	v := GetValidator()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.Struct(tt)
			if tt.Valid && err != nil {
				t.Errorf("Expected valid but got error: %v", err)
			}
			if !tt.Valid && err == nil {
				t.Error("Expected invalid but got no error")
			}
		})
	}
}

// TestFileTypeValidation 使用实际验证器测试文件类型验证
func TestFileTypeValidation(t *testing.T) {
	type TestCase struct {
		Name     string
		FileType string `validate:"file_type"`
		Valid    bool
	}

	tests := []TestCase{
		{"valid jpeg", "image/jpeg", true},
		{"valid png", "image/png", true},
		{"valid gif", "image/gif", true},
		{"valid webp", "image/webp", true},
		{"valid pdf", "application/pdf", true},
		{"valid doc", "application/msword", true},
		{"valid docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", true},
		{"valid text", "text/plain", true},
		{"valid zip", "application/zip", true},
		{"invalid type", "application/exe", false},
		{"empty type", "", false},
		{"invalid image", "image/bmp", false},
	}

	v := GetValidator()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.Struct(tt)
			if tt.Valid && err != nil {
				t.Errorf("Expected valid but got error: %v", err)
			}
			if !tt.Valid && err == nil {
				t.Error("Expected invalid but got no error")
			}
		})
	}
}

// TestFileSizeValidation 使用实际验证器测试文件大小验证
func TestFileSizeValidation(t *testing.T) {
	maxSize := int64(50 * 1024 * 1024) // 50MB

	type TestCase struct {
		Name  string
		Size  int64 `validate:"file_size"`
		Valid bool
	}

	tests := []TestCase{
		{"valid small file", 1024, true},
		{"valid 1MB file", 1024 * 1024, true},
		{"valid 10MB file", 10 * 1024 * 1024, true},
		{"valid 50MB file", maxSize, true},
		{"invalid zero size", 0, false},
		{"invalid negative size", -1024, false},
		{"invalid over 50MB", maxSize + 1, false},
		{"invalid huge file", 100 * 1024 * 1024, false},
	}

	v := GetValidator()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.Struct(tt)
			if tt.Valid && err != nil {
				t.Errorf("Expected valid but got error: %v", err)
			}
			if !tt.Valid && err == nil {
				t.Error("Expected invalid but got no error")
			}
		})
	}
}

// TestUsernameValidation 使用实际验证器测试用户名验证
func TestUsernameValidation(t *testing.T) {
	type TestCase struct {
		Name     string
		Username string `validate:"username"`
		Valid    bool
	}

	tests := []TestCase{
		{"valid simple", "user123", true},
		{"valid with underscore", "user_123", true},
		{"valid minimum length", "abc", true},
		{"valid maximum length", "12345678901234567890", true},
		{"invalid too short", "ab", false},
		{"invalid too long", "123456789012345678901", false},
		{"invalid with dash", "user-123", false},
		{"invalid with dot", "user.123", false},
		{"invalid with space", "user 123", false},
		{"invalid special chars", "user@123", false},
		{"invalid chinese", "用户名", false},
		{"empty", "", false},
	}

	v := GetValidator()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.Struct(tt)
			if tt.Valid && err != nil {
				t.Errorf("Expected valid but got error: %v", err)
			}
			if !tt.Valid && err == nil {
				t.Error("Expected invalid but got no error")
			}
		})
	}
}

// TestPhoneValidation 使用实际验证器测试手机号验证
func TestPhoneValidation(t *testing.T) {
	type TestCase struct {
		Name  string
		Phone string `validate:"phone"`
		Valid bool
	}

	tests := []TestCase{
		{"valid china mobile", "13812345678", true},
		{"valid china unicom", "13012345678", true},
		{"valid china telecom", "18112345678", true},
		{"valid with 19", "19112345678", true},
		{"invalid starts with 12", "12123456789", false},
		{"invalid starts with 10", "10123456789", false},
		{"invalid too short", "138123456", false},
		{"invalid too long", "138123456789", false},
		{"invalid with letters", "1381234567a", false},
		{"invalid with special chars", "1381234567#", false},
		{"empty", "", false},
	}

	v := GetValidator()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.Struct(tt)
			if tt.Valid && err != nil {
				t.Errorf("Expected valid but got error: %v", err)
			}
			if !tt.Valid && err == nil {
				t.Error("Expected invalid but got no error")
			}
		})
	}
}

// TestStrongPasswordValidation 使用实际验证器测试强密码验证
func TestStrongPasswordValidation(t *testing.T) {
	type TestCase struct {
		Name     string
		Password string `validate:"strong_password"`
		Valid    bool
	}

	tests := []TestCase{
		{"valid strong password", "Pass1234", true},
		{"valid with special chars", "Pass@123", true},
		{"valid longer", "StrongPass123", true},
		{"invalid too short", "Pass12", false},
		{"invalid no number", "Password", false},
		{"invalid no uppercase", "password123", false},
		{"invalid no lowercase", "PASSWORD123", false},
		{"invalid only numbers", "12345678", false},
		{"invalid only letters", "Password", false},
		{"empty", "", false},
		{"invalid only special chars", "!@#$%^&*", false},
	}

	v := GetValidator()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.Struct(tt)
			if tt.Valid && err != nil {
				t.Errorf("Expected valid but got error: %v", err)
			}
			if !tt.Valid && err == nil {
				t.Error("Expected invalid but got no error")
			}
		})
	}
}

// TestTransactionTypeValidation 使用实际验证器测试交易类型验证
func TestTransactionTypeValidation(t *testing.T) {
	type TestCase struct {
		Name string
		Type string `validate:"transaction_type"`
		Valid bool
	}

	tests := []TestCase{}

	validTypes := []string{"recharge", "consume", "transfer", "refund", "withdraw"}

	for _, validType := range validTypes {
		tests = append(tests, TestCase{"valid_" + validType, validType, true})
	}

	invalidTypes := []string{"deposit", "payment", "invalid", "", "RECHARGE"}

	for _, invalidType := range invalidTypes {
		tests = append(tests, TestCase{"invalid_" + invalidType, invalidType, false})
	}

	v := GetValidator()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.Struct(tt)
			if tt.Valid && err != nil {
				t.Errorf("Expected valid but got error: %v", err)
			}
			if !tt.Valid && err == nil {
				t.Error("Expected invalid but got no error")
			}
		})
	}
}

// TestWithdrawAccountValidation 使用实际验证器测试提现账号验证
func TestWithdrawAccountValidation(t *testing.T) {
	type TestCase struct {
		Name    string
		Account string `validate:"withdraw_account"`
		Valid   bool
	}

	tests := []TestCase{
		{"valid alipay", "alipay:user@example.com", true},
		{"valid wechat", "wechat:wxid_abc123", true},
		{"valid bank", "bank:6222021234567890", true},
		{"invalid no colon", "alipayuser@example.com", false},
		{"invalid empty account", "alipay:", false},
		{"invalid empty method", ":user@example.com", false},
		{"invalid method", "invalid:user@example.com", false},
		{"empty", "", false},
		{"invalid multiple colons", "alipay:user@example.com:extra", false},
	}

	v := GetValidator()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.Struct(tt)
			if tt.Valid && err != nil {
				t.Errorf("Expected valid but got error: %v", err)
			}
			if !tt.Valid && err == nil {
				t.Error("Expected invalid but got no error")
			}
		})
	}
}

// TestContentTypeValidation 使用实际验证器测试内容类型验证
func TestContentTypeValidation(t *testing.T) {
	type TestCase struct {
		Name       string
		ContentType string `validate:"content_type"`
		Valid      bool
	}

	tests := []TestCase{}

	validTypes := []string{"book", "chapter", "comment", "review"}

	for _, validType := range validTypes {
		tests = append(tests, TestCase{"valid_" + validType, validType, true})
	}

	invalidTypes := []string{"article", "post", "invalid", "", "BOOK"}

	for _, invalidType := range invalidTypes {
		tests = append(tests, TestCase{"invalid_" + invalidType, invalidType, false})
	}

	v := GetValidator()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.Struct(tt)
			if tt.Valid && err != nil {
				t.Errorf("Expected valid but got error: %v", err)
			}
			if !tt.Valid && err == nil {
				t.Error("Expected invalid but got no error")
			}
		})
	}
}

// TestValidateStructWithErrors 测试结构体验证
func TestValidateStructWithErrors(t *testing.T) {
	type TestRequest struct {
		Username string  `validate:"username"`
		Amount   float64 `validate:"positive_amount,amount_range"`
	}

	tests := []struct {
		name        string
		request     TestRequest
		expectError bool
	}{
		{
			name: "valid request",
			request: TestRequest{
				Username: "user123",
				Amount:   100.50,
			},
			expectError: false,
		},
		{
			name: "invalid username",
			request: TestRequest{
				Username: "ab",
				Amount:   100.50,
			},
			expectError: true,
		},
		{
			name: "invalid amount",
			request: TestRequest{
				Username: "user123",
				Amount:   -50.00,
			},
			expectError: true,
		},
		{
			name: "multiple invalid fields",
			request: TestRequest{
				Username: "ab",
				Amount:   -50.00,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStructWithErrors(tt.request)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateStructWithErrors() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

// BenchmarkValidateAmount 性能测试 - 金额验证
func BenchmarkValidateAmount(b *testing.B) {
	type TestCase struct {
		Amount float64 `validate:"amount"`
	}
	tc := TestCase{Amount: 100.50}
	v := GetValidator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Struct(tc)
	}
}

// BenchmarkValidatePhone 性能测试 - 手机号验证
func BenchmarkValidatePhone(b *testing.B) {
	type TestCase struct {
		Phone string `validate:"phone"`
	}
	tc := TestCase{Phone: "13812345678"}
	v := GetValidator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Struct(tc)
	}
}

// BenchmarkValidateStrongPassword 性能测试 - 强密码验证
func BenchmarkValidateStrongPassword(b *testing.B) {
	type TestCase struct {
		Password string `validate:"strong_password"`
	}
	tc := TestCase{Password: "Pass1234"}
	v := GetValidator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Struct(tc)
	}
}

// BenchmarkValidateComplex 性能测试 - 复杂结构体验证
func BenchmarkValidateComplex(b *testing.B) {
	type ComplexRequest struct {
		Username string  `validate:"username"`
		Phone    string  `validate:"phone"`
		Password string  `validate:"strong_password"`
		Amount   float64 `validate:"positive_amount,amount_range"`
	}
	req := ComplexRequest{
		Username: "user123",
		Phone:    "13812345678",
		Password: "Pass1234",
		Amount:   100.50,
	}
	v := GetValidator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Struct(req)
	}
}
