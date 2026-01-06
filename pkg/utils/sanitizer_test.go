package utils

import (
	"testing"
)

// TestMaskString 测试字符串脱敏
func TestMaskString(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		visibleChars int
		want         string
	}{
		{
			name:         "显示前2个字符",
			input:        "hello",
			visibleChars: 2,
			want:         "he***",
		},
		{
			name:         "显示前4个字符",
			input:        "world",
			visibleChars: 4,
			want:         "worl*",
		},
		{
			name:         "空字符串",
			input:        "",
			visibleChars: 2,
			want:         "",
		},
		{
			name:         "字符串太短",
			input:        "hi",
			visibleChars: 4,
			want:         "h*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskString(tt.input, tt.visibleChars); got != tt.want {
				t.Errorf("MaskString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMaskEmail 测试邮箱脱敏
func TestMaskEmail(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "test@example.com",
			want:  "te***@example.com",
		},
		{
			input: "a@example.com",
			want:  "***@example.com",
		},
		{
			input: "",
			want:  "",
		},
		{
			input: "invalid-email",
			want:  "invalid-email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := MaskEmail(tt.input); got != tt.want {
				t.Errorf("MaskEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMaskPhone 测试手机号脱敏
func TestMaskPhone(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "13812345678",
			want:  "138****5678",
		},
		{
			input: "+86 138-1234-5678",
			want:  "138****5678",
		},
		{
			input: "12345678901", // 11位但不是手机号格式
			want:  "123****8901",
		},
		{
			input: "123",
			want:  "123", // 太短不脱敏
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := MaskPhone(tt.input); got != tt.want {
				t.Errorf("MaskPhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMaskIDCard 测试身份证脱敏
func TestMaskIDCard(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "110101199001011234",
			want:  "110***********1234",
		},
		{
			input: "12345678",
			want:  "123****5678",
		},
		{
			input: "123",
			want:  "***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := MaskIDCard(tt.input); got != tt.want {
				t.Errorf("MaskIDCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMaskBankCard 测试银行卡脱敏
func TestMaskBankCard(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "6222021234567890123",
			want:  "6222************0123",
		},
		{
			input: "6222-2021-2345-6789-0123",
			want:  "6222************0123",
		},
		{
			input: "12345678",
			want:  "1234****5678",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := MaskBankCard(tt.input); got != tt.want {
				t.Errorf("MaskBankCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMaskToken 测试Token脱敏
func TestMaskToken(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			want:  "eyJhbGci************VCJ9",
		},
		{
			input: "short",
			want:  "********",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := MaskToken(tt.input); got != tt.want {
				t.Errorf("MaskToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestUserSanitizer 测试用户信息脱敏
func TestUserSanitizer(t *testing.T) {
	user := &UserSanitizer{
		Email:     "test@example.com",
		Phone:     "13812345678",
		IDCard:    "110101199001011234",
		BankCard:  "6222021234567890123",
		Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		Password:  "password123",
	}

	sanitized := user.Sanitize()

	if sanitized.Email != "te***@example.com" {
		t.Errorf("Email = %v, want te***@example.com", sanitized.Email)
	}

	if sanitized.Phone != "138****5678" {
		t.Errorf("Phone = %v, want 138****5678", sanitized.Phone)
	}

	if sanitized.IDCard != "110***********1234" {
		t.Errorf("IDCard = %v, want 110***********1234", sanitized.IDCard)
	}

	if sanitized.BankCard != "6222************0123" {
		t.Errorf("BankCard = %v, want 6222************0123", sanitized.BankCard)
	}

	if sanitized.Token != "eyJhbGci************VCJ9" {
		t.Errorf("Token = %v, want eyJhbGci************VCJ9", sanitized.Token)
	}

	if sanitized.Password != "***********" {
		t.Errorf("Password = %v, want ***********", sanitized.Password)
	}
}

// BenchmarkMaskString 性能测试
func BenchmarkMaskString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MaskString("hello world", 4)
	}
}

// BenchmarkMaskEmail 性能测试
func BenchmarkMaskEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MaskEmail("test@example.com")
	}
}
