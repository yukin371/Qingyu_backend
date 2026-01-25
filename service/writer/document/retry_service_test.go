package document

import (
	"errors"
	"testing"
	"time"
)

func TestRetryService_ShouldRetry_WithRetryableErrors(t *testing.T) {
	svc := NewRetryService()
	config := DefaultRetryConfig()

	err := errors.New("VERSION_CONFLICT")

	if !svc.ShouldRetry(err, config) {
		t.Error("VERSION_CONFLICT should be retryable")
	}
}

func TestRetryService_ShouldRetry_WithNonRetryableErrors(t *testing.T) {
	svc := NewRetryService()
	config := DefaultRetryConfig()

	err := errors.New("PERMISSION_DENIED")

	if svc.ShouldRetry(err, config) {
		t.Error("PERMISSION_DENIED should not be retryable")
	}
}

func TestRetryService_ShouldRetry_WithNilError(t *testing.T) {
	svc := NewRetryService()
	config := DefaultRetryConfig()

	if svc.ShouldRetry(nil, config) {
		t.Error("nil error should not be retryable")
	}
}

func TestRetryService_ShouldRetry_WithNilConfig(t *testing.T) {
	svc := NewRetryService()
	err := errors.New("VERSION_CONFLICT")

	if svc.ShouldRetry(err, nil) {
		t.Error("error with nil config should not be retryable")
	}
}

func TestRetryService_ShouldRetry_AllRetryableErrors(t *testing.T) {
	svc := NewRetryService()
	config := DefaultRetryConfig()

	retryableErrors := []string{
		"VERSION_CONFLICT",
		"NETWORK_ERROR",
		"TIMEOUT",
		"SERVICE_UNAVAILABLE",
	}

	for _, errMsg := range retryableErrors {
		err := errors.New(errMsg)
		if !svc.ShouldRetry(err, config) {
			t.Errorf("%s should be retryable", errMsg)
		}
	}
}

func TestRetryService_GetRetryDelay_ExponentialBackoff(t *testing.T) {
	svc := NewRetryService()
	config := &RetryConfig{
		RetryDelay: 1000, // 1秒
	}

	tests := []struct {
		name        string
		attempt     int
		expectedMin time.Duration
		expectedMax time.Duration
	}{
		{"attempt 0", 0, 1000 * time.Millisecond, 1000 * time.Millisecond},   // 2^0 * 1s = 1s
		{"attempt 1", 1, 2000 * time.Millisecond, 2000 * time.Millisecond},   // 2^1 * 1s = 2s
		{"attempt 2", 2, 4000 * time.Millisecond, 4000 * time.Millisecond},   // 2^2 * 1s = 4s
		{"attempt 3", 3, 8000 * time.Millisecond, 8000 * time.Millisecond},   // 2^3 * 1s = 8s
		{"attempt 10", 10, 60 * time.Second, 60 * time.Second},                // 超过最大值
		{"attempt 20", 20, 60 * time.Second, 60 * time.Second},                // 远超最大值
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delay := svc.GetRetryDelay(tt.attempt, config)
			if delay < tt.expectedMin || delay > tt.expectedMax {
				t.Errorf("attempt %d: expected delay between %v and %v, got %v",
					tt.attempt, tt.expectedMin, tt.expectedMax, delay)
			}
		})
	}
}

func TestRetryService_GetRetryDelay_WithNilConfig(t *testing.T) {
	svc := NewRetryService()

	delay := svc.GetRetryDelay(0, nil)
	if delay != 0 {
		t.Errorf("expected 0 delay with nil config, got %v", delay)
	}
}

func TestRetryService_GetRetryDelay_WithNegativeAttempt(t *testing.T) {
	svc := NewRetryService()
	config := &RetryConfig{
		RetryDelay: 1000,
	}

	delay := svc.GetRetryDelay(-1, config)
	if delay != 0 {
		t.Errorf("expected 0 delay with negative attempt, got %v", delay)
	}
}

func TestRetryService_CanRetry(t *testing.T) {
	svc := NewRetryService()
	config := &RetryConfig{
		MaxRetries: 3,
	}

	tests := []struct {
		name     string
		attempt  int
		expected bool
	}{
		{"attempt 0", 0, true},  // 还可以重试3次
		{"attempt 1", 1, true},  // 还可以重试2次
		{"attempt 2", 2, true},  // 还可以重试1次
		{"attempt 3", 3, false}, // 已达到最大次数
		{"attempt 4", 4, false}, // 超过最大次数
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.CanRetry(tt.attempt, config)
			if result != tt.expected {
				t.Errorf("attempt %d: expected %v, got %v", tt.attempt, tt.expected, result)
			}
		})
	}
}

func TestRetryService_CanRetry_WithNilConfig(t *testing.T) {
	svc := NewRetryService()

	result := svc.CanRetry(0, nil)
	if result {
		t.Error("expected false with nil config")
	}
}

func TestRetryConfig_DefaultValues(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxRetries != 3 {
		t.Errorf("expected MaxRetries=3, got %d", config.MaxRetries)
	}

	if config.RetryDelay != 1000 {
		t.Errorf("expected RetryDelay=1000, got %d", config.RetryDelay)
	}

	if len(config.RetryableErrors) == 0 {
		t.Error("expected non-empty RetryableErrors")
	}

	expectedErrors := []string{
		"VERSION_CONFLICT",
		"NETWORK_ERROR",
		"TIMEOUT",
		"SERVICE_UNAVAILABLE",
	}

	if len(config.RetryableErrors) != len(expectedErrors) {
		t.Errorf("expected %d retryable errors, got %d", len(expectedErrors), len(config.RetryableErrors))
	}

	for i, err := range expectedErrors {
		if config.RetryableErrors[i] != err {
			t.Errorf("expected error %d to be %s, got %s", i, err, config.RetryableErrors[i])
		}
	}
}

func TestRetryConfig_GetMaxDelay(t *testing.T) {
	config := &RetryConfig{}
	maxDelay := config.GetMaxDelay()

	if maxDelay != 60*time.Second {
		t.Errorf("expected max delay 60s, got %v", maxDelay)
	}
}

func TestRetryService_extractErrorCode(t *testing.T) {
	svc := NewRetryService()

	err := errors.New("TEST_ERROR")
	errorCode := svc.extractErrorCode(err)

	if errorCode != "TEST_ERROR" {
		t.Errorf("expected error code 'TEST_ERROR', got '%s'", errorCode)
	}
}

func TestRetryService_ExponentialBackoff_Formula(t *testing.T) {
	svc := NewRetryService()
	config := &RetryConfig{
		RetryDelay: 1000, // 1秒基准
	}

	// 验证公式：delay = baseDelay * 2^attempt
	// 2^0 = 1, 2^1 = 2, 2^2 = 4, 2^3 = 8
	expectedDelays := map[int]time.Duration{
		0: 1 * time.Second,
		1: 2 * time.Second,
		2: 4 * time.Second,
		3: 8 * time.Second,
		4: 16 * time.Second,
		5: 32 * time.Second,
	}

	for attempt, expected := range expectedDelays {
		t.Run("attempt_"+string(rune('0'+attempt)), func(t *testing.T) {
			delay := svc.GetRetryDelay(attempt, config)
			if delay != expected {
				t.Errorf("attempt %d: expected %v, got %v", attempt, expected, delay)
			}
		})
	}
}

func TestRetryService_MaxDelayConstraint(t *testing.T) {
	svc := NewRetryService()
	config := &RetryConfig{
		RetryDelay: 1000, // 1秒
	}

	// 测试即使尝试次数很大，延迟也不会超过60秒
	for attempt := 0; attempt < 20; attempt++ {
		delay := svc.GetRetryDelay(attempt, config)
		if delay > 60*time.Second {
			t.Errorf("attempt %d: delay %v exceeds max 60s", attempt, delay)
		}
	}
}

func TestRetryService_ShouldRetry_EmptyRetryableErrors(t *testing.T) {
	svc := NewRetryService()
	config := &RetryConfig{
		MaxRetries:      3,
		RetryDelay:      1000,
		RetryableErrors: []string{}, // 空列表
	}

	err := errors.New("VERSION_CONFLICT")

	if svc.ShouldRetry(err, config) {
		t.Error("should not retry with empty retryable errors list")
	}
}

func TestRetryService_CanRetry_ZeroMaxRetries(t *testing.T) {
	svc := NewRetryService()
	config := &RetryConfig{
		MaxRetries: 0,
	}

	if svc.CanRetry(0, config) {
		t.Error("should not allow any retry with MaxRetries=0")
	}
}
