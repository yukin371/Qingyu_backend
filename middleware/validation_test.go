package middleware

import (
	"testing"
)

// TestValidateEmail 测试邮箱验证
func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		want    bool
	}{
		{
			name:  "valid email",
			email: "test@example.com",
			want:  true,
		},
		{
			name:  "valid email with subdomain",
			email: "user@mail.example.com",
			want:  true,
		},
		{
			name:  "valid email with numbers",
			email: "user123@example123.com",
			want:  true,
		},
		{
			name:  "valid email with special chars",
			email: "user+tag@example.com",
			want:  true,
		},
		{
			name:  "valid email with dots",
			email: "first.last@example.com",
			want:  true,
		},
		{
			name:  "empty email",
			email: "",
			want:  false,
		},
		{
			name:  "no @ symbol",
			email: "invalidemail.com",
			want:  false,
		},
		{
			name:  "no domain",
			email: "user@",
			want:  false,
		},
		{
			name:  "no local part",
			email: "@example.com",
			want:  false,
		},
		{
			name:  "multiple @ symbols",
			email: "user@@example.com",
			want:  false,
		},
		{
			name:  "invalid domain",
			email: "user@.com",
			want:  false,
		},
		{
			name:  "space in email",
			email: "user @example.com",
			want:  false,
		},
		{
			name:  "no TLD",
			email: "user@example",
			want:  false,
		},
		{
			name:  "underscore in domain",
			email: "user@ex_ample.com",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateEmail(tt.email)
			if got != tt.want {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

// TestValidatePhone 测试手机号验证
func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  bool
	}{
		{
			name:  "valid phone - 13",
			phone: "13812345678",
			want:  true,
		},
		{
			name:  "valid phone - 15",
			phone: "15012345678",
			want:  true,
		},
		{
			name:  "valid phone - 18",
			phone: "18612345678",
			want:  true,
		},
		{
			name:  "empty phone",
			phone: "",
			want:  false,
		},
		{
			name:  "too short",
			phone: "1381234567",
			want:  false,
		},
		{
			name:  "too long",
			phone: "138123456789",
			want:  false,
		},
		{
			name:  "invalid prefix - 10",
			phone: "10123456789",
			want:  false,
		},
		{
			name:  "invalid prefix - 12",
			phone: "12123456789",
			want:  false,
		},
		{
			name:  "with letters",
			phone: "1381234567a",
			want:  false,
		},
		{
			name:  "with special chars",
			phone: "138-1234-5678",
			want:  false,
		},
		{
			name:  "with spaces",
			phone: "138 1234 5678",
			want:  false,
		},
		{
			name:  "starts with 0",
			phone: "013812345678",
			want:  false,
		},
		{
			name:  "not starting with 1",
			phone: "23812345678",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidatePhone(tt.phone)
			if got != tt.want {
				t.Errorf("ValidatePhone(%q) = %v, want %v", tt.phone, got, tt.want)
			}
		})
	}
}

// TestValidatePassword 测试密码验证
func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		password string
		wantErr bool
	}{
		{
			name:     "valid password - all requirements",
			password: "Abc123!@#",
			wantErr:  false,
		},
		{
			name:     "valid password - longer",
			password: "MySecureP@ssw0rd123",
			wantErr:  false,
		},
		{
			name:     "too short - 7 chars",
			password: "Abc123!",
			wantErr:  true,
		},
		{
			name:     "too long - 129 chars",
			password: string(make([]byte, 129)),
			wantErr:  true,
		},
		{
			name:     "no uppercase",
			password: "abc123!@#",
			wantErr:  true,
		},
		{
			name:     "no lowercase",
			password: "ABC123!@#",
			wantErr:  true,
		},
		{
			name:     "no numbers",
			password: "Abcdef!@#",
			wantErr:  true,
		},
		{
			name:     "no special chars",
			password: "Abc123456",
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "only lowercase",
			password: "abcdefgh",
			wantErr:  true,
		},
		{
			name:     "only numbers",
			password: "12345678",
			wantErr:  true,
		},
		{
			name:     "only special chars",
			password: "!@#$%^&*",
			wantErr:  true,
		},
		{
			name:     "minimum valid - 8 chars",
			password: "Aa1!Aa1!",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}

// TestValidateUsername 测试用户名验证
func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{
			name:     "valid username - alphanumeric",
			username: "user123",
			want:     true,
		},
		{
			name:     "valid username - with underscore",
			username: "user_name",
			want:     true,
		},
		{
			name:     "valid username - with dash",
			username: "user-name",
			want:     true,
		},
		{
			name:     "valid username - mixed",
			username: "User_Name-123",
			want:     true,
		},
		{
			name:     "empty username",
			username: "",
			want:     false,
		},
		{
			name:     "too short - 2 chars",
			username: "ab",
			want:     false,
		},
		{
			name:     "too long - 33 chars",
			username: string(make([]byte, 33)),
			want:     false,
		},
		{
			name:     "starts with number",
			username: "123user",
			want:     false,
		},
		{
			name:     "starts with special char",
			username: "_username",
			want:     false,
		},
		{
			name:     "contains space",
			username: "user name",
			want:     false,
		},
		{
			name:     "contains special char",
			username: "user@name",
			want:     false,
		},
		{
			name:     "Chinese characters",
			username: "用户名",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateUsername(tt.username)
			if got != tt.want {
				t.Errorf("ValidateUsername(%q) = %v, want %v", tt.username, got, tt.want)
			}
		})
	}
}

// TestValidateURL 测试URL验证
func TestValidateURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "valid URL - http",
			url:  "http://example.com",
			want: true,
		},
		{
			name: "valid URL - https",
			url:  "https://example.com",
			want: true,
		},
		{
			name: "valid URL - with path",
			url:  "https://example.com/path/to/resource",
			want: true,
		},
		{
			name: "valid URL - with query",
			url:  "https://example.com?query=param",
			want: true,
		},
		{
			name: "valid URL - with port",
			url:  "https://example.com:8080",
			want: true,
		},
		{
			name: "empty URL",
			url:  "",
			want: false,
		},
		{
			name: "no protocol",
			url:  "example.com",
			want: false,
		},
		{
			name: "invalid protocol",
			url:  "ftp://example.com",
			want: false,
		},
		{
			name: "space in URL",
			url:  "https://example .com",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateURL(tt.url)
			if got != tt.want {
				t.Errorf("ValidateURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

// TestSanitizeString 测试字符串清理
func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "normal string",
			input: "Hello World",
			want:  "Hello World",
		},
		{
			name:  "with leading/trailing spaces",
			input: "  Hello World  ",
			want:  "Hello World",
		},
		{
			name:  "with multiple spaces",
			input: "Hello    World",
			want:  "Hello World",
		},
		{
			name:  "with newlines",
			input: "Hello\nWorld",
			want:  "Hello World",
		},
		{
			name:  "with tabs",
			input: "Hello\tWorld",
			want:  "Hello World",
		},
		{
			name:  "with special HTML tags",
			input: "Hello<script>alert('xss')</script>World",
			want:  "Helloalert('xss')World",
		},
		{
			name:  "SQL injection attempt",
			input: "'; DROP TABLE users; --",
			want:  "'; TABLE users; --",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "only spaces",
			input: "     ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeString(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
