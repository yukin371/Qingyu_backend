package validator

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
)

// TestValidationErrors 测试ValidationErrors
func TestValidationErrors(t *testing.T) {
	// 创建一个会导致验证错误的请求
	type TestRequest struct {
		Username string `validate:"username"`
		Amount   float64 `validate:"positive_amount"`
	}

	req := TestRequest{
		Username: "ab", // 太短
		Amount:   -50,  // 负数
	}

	err := ValidateStruct(req)
	if err == nil {
		t.Fatal("Expected validation error but got nil")
	}

	// 测试TranslateError
	validationErrors := TranslateError(err)
	if validationErrors == nil {
		t.Fatal("TranslateError returned nil")
	}

	// 测试GetFieldErrors
	fieldErrors := validationErrors.GetFieldErrors()
	if len(fieldErrors) == 0 {
		t.Error("GetFieldErrors returned empty map")
	}

	t.Logf("Field errors: %v", fieldErrors)
}

// TestFieldError 测试FieldError
func TestFieldError(t *testing.T) {
	v := GetValidator()

	type TestRequest struct {
		Username string `validate:"username"`
	}

	req := TestRequest{
		Username: "ab",
	}

	err := v.Struct(req)
	if err == nil {
		t.Fatal("Expected validation error but got nil")
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatal("Error is not validator.ValidationErrors type")
	}

	for _, e := range validationErrors {
		field := e.Field()
		tag := e.Tag()
		value := e.Value()
		param := e.Param()

		t.Logf("Field: %s, Tag: %s, Value: %v, Param: %s", field, tag, value, param)

		if field != "Username" {
			t.Errorf("Expected field 'Username', got '%s'", field)
		}
		if tag != "username" {
			t.Errorf("Expected tag 'username', got '%s'", tag)
		}
	}
}

// TestGetInitError 测试GetInitError
func TestGetInitError(t *testing.T) {
	// 第一次调用会初始化
	v := GetValidator()
	if v == nil {
		t.Fatal("GetValidator returned nil")
	}

	// 获取初始化错误（应该为nil）
	err := GetInitError()
	// 目前没有实现实际的错误返回，只是测试函数存在
	if err != nil {
		t.Logf("Init error: %v", err)
	}
}

// TestGetValidatorConcurrency 测试GetValidator并发安全
func TestGetValidatorConcurrency(t *testing.T) {
	done := make(chan bool)

	// 并发调用100次
	for i := 0; i < 100; i++ {
		go func() {
			v := GetValidator()
			if v == nil {
				t.Error("GetValidator returned nil")
			}
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 100; i++ {
		<-done
	}

	t.Log("Concurrent GetValidator calls completed successfully")
}

// TestValidateStructWithErrorsEmpty 测试ValidateStructWithErrors空错误
func TestValidateStructWithErrorsEmpty(t *testing.T) {
	type TestRequest struct {
		Username string `validate:"username"`
	}

	req := TestRequest{
		Username: "user123",
	}

	err := ValidateStructWithErrors(req)
	if err != nil {
		t.Errorf("Expected nil error for valid input, got %v", err)
	}
}

// TestRegistrationStatusMethods 测试RegistrationStatus方法
func TestRegistrationStatusMethods(t *testing.T) {
	status := GetRegistrationStatus()

	// 测试IsComplete
	isComplete := status.IsComplete()
	t.Logf("IsComplete: %v", isComplete)

	// 测试GetFailedCount
	failedCount := status.GetFailedCount()
	t.Logf("Failed count: %d", failedCount)

	// 测试GetFailedTags
	failedTags := status.GetFailedTags()
	t.Logf("Failed tags: %v", failedTags)

	// 验证一致性
	if failedCount != len(failedTags) {
		t.Errorf("Failed count %d != len(failedTags) %d", failedCount, len(failedTags))
	}
}

// TestMultipleValidators 测试多个验证器组合
func TestMultipleValidators(t *testing.T) {
	type TestRequest struct {
		Username string  `validate:"username"`
		Phone    string  `validate:"phone"`
		Amount   float64 `validate:"positive_amount,amount_range"`
	}

	tests := []struct {
		name        string
		request     TestRequest
		expectError bool
	}{
		{
			name: "all valid",
			request: TestRequest{
				Username: "user123",
				Phone:    "13812345678",
				Amount:   100.50,
			},
			expectError: false,
		},
		{
			name: "all invalid",
			request: TestRequest{
				Username: "ab",
				Phone:    "12345",
				Amount:   -50,
			},
			expectError: true,
		},
		{
			name: "partial invalid",
			request: TestRequest{
				Username: "user123",
				Phone:    "12345",
				Amount:   100.50,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.request)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateStruct() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	type TestRequest struct {
		StringField string `validate:"username"`
		FloatField  float64 `validate:"amount"`
		IntField    int64   `validate:"file_size"`
	}

	tests := []struct {
		name        string
		request     TestRequest
		expectError bool
	}{
		{
			name: "zero float",
			request: TestRequest{
				StringField: "user123",
				FloatField:  0,
				IntField:    1000,
			},
			expectError: false, // amount allows 0
		},
		{
			name: "very large int",
			request: TestRequest{
				StringField: "user123",
				FloatField:  100,
				IntField:    9223372036854775807, // max int64
			},
			expectError: true, // file_size max is 50MB
		},
		{
			name: "negative int",
			request: TestRequest{
				StringField: "user123",
				FloatField:  100,
				IntField:    -1,
			},
			expectError: true, // file_size requires > 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.request)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateStruct() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

// TestValidatorSingleton 测试验证器单例
func TestValidatorSingleton(t *testing.T) {
	v1 := GetValidator()
	v2 := GetValidator()

	// 验证返回的是同一个实例
	if v1 != v2 {
		t.Error("GetValidator() did not return the same instance")
	}

	// 验证指针相等
	if reflect.ValueOf(v1).Pointer() != reflect.ValueOf(v2).Pointer() {
		t.Error("GetValidator() pointers are not equal")
	}
}
