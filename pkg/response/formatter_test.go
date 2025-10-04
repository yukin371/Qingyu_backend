package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterSensitiveFields_Struct(t *testing.T) {
	type User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"` // 敏感字段
		APIKey   string `json:"api_key"`  // 敏感字段
	}

	user := User{
		ID:       "user_001",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "secret123",
		APIKey:   "sk_test_123456",
	}

	filtered := FilterSensitiveFields(user)
	result, ok := filtered.(map[string]interface{})
	assert.True(t, ok)

	// 验证非敏感字段保留
	assert.Equal(t, "user_001", result["id"])
	assert.Equal(t, "testuser", result["username"])
	assert.Equal(t, "test@example.com", result["email"])

	// 验证敏感字段被过滤
	assert.NotContains(t, result, "password")
	assert.NotContains(t, result, "api_key")
}

func TestFilterSensitiveFields_Slice(t *testing.T) {
	type User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	users := []User{
		{Username: "user1", Password: "pass1"},
		{Username: "user2", Password: "pass2"},
	}

	filtered := FilterSensitiveFields(users)
	result, ok := filtered.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(result))

	// 验证每个元素的密码都被过滤
	for _, item := range result {
		userMap, ok := item.(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, userMap, "username")
		assert.NotContains(t, userMap, "password")
	}
}

func TestFilterSensitiveFields_Map(t *testing.T) {
	data := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "secret123",
		"token":    "abc123",
	}

	filtered := FilterSensitiveFields(data)
	result, ok := filtered.(map[string]interface{})
	assert.True(t, ok)

	// 验证非敏感字段保留
	assert.Equal(t, "testuser", result["username"])
	assert.Equal(t, "test@example.com", result["email"])

	// 验证敏感字段被过滤
	assert.NotContains(t, result, "password")
	assert.NotContains(t, result, "token")
}

func TestFilterSensitiveFields_NestedStruct(t *testing.T) {
	type Credentials struct {
		Password string `json:"password"`
		APIKey   string `json:"api_key"`
	}

	type User struct {
		Username string      `json:"username"`
		Creds    Credentials `json:"credentials"`
	}

	user := User{
		Username: "testuser",
		Creds: Credentials{
			Password: "secret123",
			APIKey:   "sk_123",
		},
	}

	filtered := FilterSensitiveFields(user)
	result, ok := filtered.(map[string]interface{})
	assert.True(t, ok)

	// 验证用户名保留
	assert.Equal(t, "testuser", result["username"])

	// 验证嵌套的敏感字段被过滤
	creds, ok := result["credentials"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotContains(t, creds, "password")
	assert.NotContains(t, creds, "api_key")
}

func TestIsSensitiveField(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		expected  bool
	}{
		{"密码字段", "password", true},
		{"大写密码字段", "Password", true},
		{"包含密码", "user_password", true},
		{"密钥", "secret", true},
		{"Token", "access_token", true},
		{"API密钥", "api_key", true},
		{"普通字段", "username", false},
		{"邮箱", "email", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSensitiveField(tt.fieldName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAddSensitiveField(t *testing.T) {
	// 保存原始列表
	original := make([]string, len(SensitiveFields))
	copy(original, SensitiveFields)
	defer func() {
		SensitiveFields = original
	}()

	// 添加新字段
	AddSensitiveField("credit_card")

	// 验证字段被添加
	found := false
	for _, field := range SensitiveFields {
		if field == "credit_card" {
			found = true
			break
		}
	}
	assert.True(t, found, "新字段应该被添加到敏感字段列表")

	// 测试过滤
	type Payment struct {
		Amount     float64 `json:"amount"`
		CreditCard string  `json:"credit_card"`
	}

	payment := Payment{Amount: 100.00, CreditCard: "4111111111111111"}
	filtered := FilterSensitiveFields(payment)
	result, ok := filtered.(map[string]interface{})
	assert.True(t, ok)

	assert.Contains(t, result, "amount")
	assert.NotContains(t, result, "credit_card")
}

func TestRemoveSensitiveField(t *testing.T) {
	// 保存原始列表
	original := make([]string, len(SensitiveFields))
	copy(original, SensitiveFields)
	defer func() {
		SensitiveFields = original
	}()

	// 移除字段
	RemoveSensitiveField("password")

	// 验证字段被移除
	found := false
	for _, field := range SensitiveFields {
		if field == "password" {
			found = true
			break
		}
	}
	assert.False(t, found, "字段应该被从敏感字段列表移除")
}

func TestFilterSensitiveFields_NilValue(t *testing.T) {
	result := FilterSensitiveFields(nil)
	assert.Nil(t, result)
}

func TestFilterSensitiveFields_PrimitiveTypes(t *testing.T) {
	// 字符串
	assert.Equal(t, "test", FilterSensitiveFields("test"))

	// 数字
	assert.Equal(t, 123, FilterSensitiveFields(123))
	assert.Equal(t, 45.67, FilterSensitiveFields(45.67))

	// 布尔
	assert.Equal(t, true, FilterSensitiveFields(true))
}

func BenchmarkFilterSensitiveFields(b *testing.B) {
	type User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Token    string `json:"token"`
	}

	user := User{
		ID:       "user_001",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "secret123",
		Token:    "token_abc123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FilterSensitiveFields(user)
	}
}
