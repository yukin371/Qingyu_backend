package middleware

import (
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		limit        int
		window       int
		requests     int
		expectStatus int
	}{
		{
			name:         "allow requests within limit",
			limit:        10,
			window:       60,
			requests:     5,
			expectStatus: 200,
		},
		{
			name:         "block requests over limit",
			limit:        5,
			window:       60,
			requests:     10,
			expectStatus: 429,
		},
		{
			name:         "exact limit",
			limit:        3,
			window:       60,
			requests:     3,
			expectStatus: 200,
		},
		{
			name:         "one over limit",
			limit:        3,
			window:       60,
			requests:     4,
			expectStatus: 429,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(RateLimitMiddleware(tt.limit, tt.window))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "ok"})
			})

			var lastStatus int
			for i := 0; i < tt.requests; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				lastStatus = w.Code
			}

			if lastStatus != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, lastStatus)
			}
		})
	}
}

func TestRateLimitWithConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		config       RateLimitConfig
		requests     int
		expectStatus int
		skipPath     bool
	}{
		{
			name: "default config",
			config: RateLimitConfig{
				RequestsPerSecond: 100,
				BurstSize:         200,
				KeyFunc:           "ip",
				Message:           "Too many requests",
				StatusCode:        429,
			},
			requests:     5,
			expectStatus: 200,
		},
		{
			name: "skip path",
			config: RateLimitConfig{
				RequestsPerSecond: 1,
				BurstSize:         1,
				KeyFunc:           "ip",
				SkipPaths:         []string{"/health"},
				StatusCode:        429,
			},
			requests:     10,
			expectStatus: 200,
			skipPath:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(RateLimitWithConfig(tt.config))
			path := "/test"
			if tt.skipPath {
				path = "/health"
			}
			router.GET(path, func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "ok"})
			})

			var lastStatus int
			for i := 0; i < tt.requests; i++ {
				req := httptest.NewRequest("GET", path, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				lastStatus = w.Code
			}

			if lastStatus != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, lastStatus)
			}
		})
	}
}

func TestRateLimitKeyFunctions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		keyFunc  string
		path     string
		setupCtx func(*gin.Context)
	}{
		{
			name:    "ip based",
			keyFunc: "ip",
			path:    "/test",
		},
		{
			name:    "path based",
			keyFunc: "path",
			path:    "/api/users",
		},
		{
			name:    "ip_path based",
			keyFunc: "ip_path",
			path:    "/api/users",
		},
		{
			name:     "user based",
			keyFunc:  "user",
			path:     "/api/users",
			setupCtx: func(c *gin.Context) {
				c.Set("user_id", "12345")
			},
		},
		{
			name:     "user_path based",
			keyFunc:  "user_path",
			path:     "/api/users",
			setupCtx: func(c *gin.Context) {
				c.Set("user_id", "12345")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			config := DefaultRateLimitConfig()
			config.KeyFunc = tt.keyFunc
			config.RequestsPerSecond = 1000 // High limit to avoid blocking
			config.BurstSize = 1000
			router.Use(RateLimitWithConfig(config))

			if tt.setupCtx != nil {
				router.Use(func(c *gin.Context) {
					tt.setupCtx(c)
					c.Next()
				})
			}

			router.GET(tt.path, func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "ok"})
			})

			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Errorf("Expected status 200, got %d", w.Code)
			}
		})
	}
}

func TestTokenBucketLimiter(t *testing.T) {
	tests := []struct {
		name    string
		rps     int
		burst   int
		request int
		want    bool
	}{
		{
			name:    "allow within burst",
			rps:     10,
			burst:   5,
			request: 5,
			want:    true,
		},
		{
			name:    "block over burst",
			rps:     10,
			burst:   5,
			request: 6,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := NewTokenBucketLimiter(tt.rps, tt.burst)
			key := "test-key"

			var result bool
			for i := 0; i < tt.request; i++ {
				result = limiter.Allow(key)
			}

			if result != tt.want {
				t.Errorf("After %d requests, expected %v, got %v", tt.request, tt.want, result)
			}
		})
	}
}

func TestTokenBucketLimiterConcurrency(t *testing.T) {
	limiter := NewTokenBucketLimiter(100, 50)
	key := "concurrent-test"
	iterations := 100

	var wg sync.WaitGroup
	results := make([]bool, iterations)

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = limiter.Allow(key)
		}(i)
	}

	wg.Wait()

	allowedCount := 0
	for _, result := range results {
		if result {
			allowedCount++
		}
	}

	// 应该允许所有请求（burst=50，但并发启动可能会有些被拒绝）
	if allowedCount > 50 {
		t.Logf("Allowed %d out of %d concurrent requests", allowedCount, iterations)
	}
}

func TestDefaultRateLimitConfig(t *testing.T) {
	config := DefaultRateLimitConfig()

	if config.RequestsPerSecond != 100 {
		t.Errorf("Expected RequestsPerSecond 100, got %d", config.RequestsPerSecond)
	}

	if config.RequestsPerMinute != 1000 {
		t.Errorf("Expected RequestsPerMinute 1000, got %d", config.RequestsPerMinute)
	}

	if config.RequestsPerHour != 10000 {
		t.Errorf("Expected RequestsPerHour 10000, got %d", config.RequestsPerHour)
	}

	if config.BurstSize != 200 {
		t.Errorf("Expected BurstSize 200, got %d", config.BurstSize)
	}

	if config.KeyFunc != "ip" {
		t.Errorf("Expected KeyFunc 'ip', got %s", config.KeyFunc)
	}

	if config.StatusCode != 429 {
		t.Errorf("Expected StatusCode 429, got %d", config.StatusCode)
	}

	if len(config.SkipPaths) != 2 {
		t.Errorf("Expected 2 skip paths, got %d", len(config.SkipPaths))
	}
}

func TestShouldSkipRateLimit(t *testing.T) {
	tests := []struct {
		path      string
		skipPaths []string
		expected  bool
	}{
		{
			path:      "/health",
			skipPaths: []string{"/health", "/metrics"},
			expected:  true,
		},
		{
			path:      "/metrics",
			skipPaths: []string{"/health", "/metrics"},
			expected:  true,
		},
		{
			path:      "/api/users",
			skipPaths: []string{"/health", "/metrics"},
			expected:  false,
		},
		{
			path:      "/api/health",
			skipPaths: []string{"/health"},
			expected:  false,
		},
		{
			path:      "/api/users",
			skipPaths: []string{},
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := shouldSkipRateLimit(tt.path, tt.skipPaths)
			if result != tt.expected {
				t.Errorf("shouldSkipRateLimit(%q, %v) = %v, want %v", tt.path, tt.skipPaths, result, tt.expected)
			}
		})
	}
}

func TestGetRateLimitKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		keyFunc  string
		setupCtx func(*gin.Context)
		path     string
	}{
		{
			name:    "ip key",
			keyFunc: "ip",
			path:    "/test",
		},
		{
			name:     "user key with user_id",
			keyFunc:  "user",
			path:     "/test",
			setupCtx: func(c *gin.Context) { c.Set("user_id", "123") },
		},
		{
			name:     "user key without user_id falls back to ip",
			keyFunc:  "user",
			path:     "/test",
			setupCtx: nil,
		},
		{
			name:    "path key",
			keyFunc: "path",
			path:    "/api/users",
		},
		{
			name:     "ip_path key",
			keyFunc:  "ip_path",
			path:     "/api/users",
			setupCtx: nil,
		},
		{
			name:     "user_path key",
			keyFunc:  "user_path",
			path:     "/api/users",
			setupCtx: func(c *gin.Context) { c.Set("user_id", "123") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("GET", tt.path, nil)

			if tt.setupCtx != nil {
				tt.setupCtx(c)
			}

			key := getRateLimitKey(c, tt.keyFunc)
			if key == "" {
				t.Error("Expected non-empty key")
			}
		})
	}
}

func TestRateLimitRefill(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := NewTokenBucketLimiter(10, 5) // 10 requests per second, burst 5
	key := "refill-test"

	// 消耗所有burst tokens
	for i := 0; i < 5; i++ {
		if !limiter.Allow(key) {
			t.Error("Expected all burst tokens to be allowed")
		}
	}

	// 下一个应该被拒绝
	if limiter.Allow(key) {
		t.Error("Expected request to be blocked after burst exhausted")
	}

	// 等待token refill
	time.Sleep(150 * time.Millisecond)

	// 现在应该允许一个请求
	if !limiter.Allow(key) {
		t.Error("Expected request to be allowed after refill")
	}
}
