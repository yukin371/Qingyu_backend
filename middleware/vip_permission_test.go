package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetVIPLevelName(t *testing.T) {
	tests := []struct {
		level    VIPLevel
		expected string
	}{
		{VIPLevelNone, "普通用户"},
		{VIPLevelBasic, "基础VIP"},
		{VIPLevelPlus, "VIP Plus"},
		{VIPLevelPro, "VIP Pro"},
		{VIPLevelUltra, "VIP Ultra"},
		{VIPLevel(99), "未知等级"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GetVIPLevelName(tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCheckVIPLevel(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		user     *VIPUser
		config   VIPConfig
		expected bool
	}{
		{
			name: "VIP等级满足且未过期",
			user: &VIPUser{
				VIPLevel:  VIPLevelPlus,
				VIPExpire: now.Add(30 * 24 * time.Hour),
			},
			config: VIPConfig{
				RequiredLevel: VIPLevelPlus,
			},
			expected: true,
		},
		{
			name: "VIP等级不足",
			user: &VIPUser{
				VIPLevel:  VIPLevelBasic,
				VIPExpire: now.Add(30 * 24 * time.Hour),
			},
			config: VIPConfig{
				RequiredLevel: VIPLevelPlus,
			},
			expected: false,
		},
		{
			name: "VIP已过期",
			user: &VIPUser{
				VIPLevel:  VIPLevelPlus,
				VIPExpire: now.Add(-24 * time.Hour),
			},
			config: VIPConfig{
				RequiredLevel: VIPLevelPlus,
			},
			expected: false,
		},
		{
			name: "试用用户-允许试用",
			user: &VIPUser{
				VIPLevel:    VIPLevelNone,
				IsTrialUser: true,
				TrialExpire: now.Add(7 * 24 * time.Hour),
			},
			config: VIPConfig{
				RequiredLevel: VIPLevelBasic,
				AllowTrial:    true,
			},
			expected: true,
		},
		{
			name: "试用用户-不允许试用",
			user: &VIPUser{
				VIPLevel:    VIPLevelNone,
				IsTrialUser: true,
				TrialExpire: now.Add(7 * 24 * time.Hour),
			},
			config: VIPConfig{
				RequiredLevel: VIPLevelBasic,
				AllowTrial:    false,
			},
			expected: false,
		},
		{
			name: "试用已过期",
			user: &VIPUser{
				VIPLevel:    VIPLevelNone,
				IsTrialUser: true,
				TrialExpire: now.Add(-24 * time.Hour),
			},
			config: VIPConfig{
				RequiredLevel: VIPLevelBasic,
				AllowTrial:    true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkVIPLevel(tt.user, tt.config)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRequireVIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		requiredLevel  VIPLevel
		expectedStatus int
		expectedAbort  bool
	}{
		{
			name: "未认证用户",
			setupContext: func(c *gin.Context) {
				// 不设置userId
			},
			requiredLevel:  VIPLevelBasic,
			expectedStatus: http.StatusUnauthorized,
			expectedAbort:  true,
		},
		{
			name: "VIP等级满足",
			setupContext: func(c *gin.Context) {
				c.Set("userId", "user123")
				c.Set("vipUser", &VIPUser{
					UserID:    "user123",
					VIPLevel:  VIPLevelPlus,
					VIPExpire: time.Now().Add(30 * 24 * time.Hour),
				})
			},
			requiredLevel:  VIPLevelPlus,
			expectedStatus: http.StatusOK,
			expectedAbort:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			tt.setupContext(c)

			middleware := RequireVIP(tt.requiredLevel)
			middleware(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedAbort, c.IsAborted())
		})
	}
}

func TestGetFeatureRequiredLevel(t *testing.T) {
	tests := []struct {
		feature  VIPFeature
		expected VIPLevel
	}{
		{FeatureOfflineReading, VIPLevelBasic},
		{FeatureHighlightExport, VIPLevelBasic},
		{FeatureMultiDevice, VIPLevelPlus},
		{FeatureAdvancedSearch, VIPLevelPlus},
		{FeatureCustomTheme, VIPLevelPlus},
		{FeatureAdFree, VIPLevelPro},
		{FeatureEarlyAccess, VIPLevelPro},
		{FeaturePremiumSupport, VIPLevelUltra},
		{VIPFeature("unknown"), VIPLevelNone},
	}

	for _, tt := range tests {
		t.Run(string(tt.feature), func(t *testing.T) {
			result := getFeatureRequiredLevel(tt.feature)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCheckChapterVIPAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	now := time.Now()

	tests := []struct {
		name              string
		setupContext      func(*gin.Context)
		expectedCanAccess bool
	}{
		{
			name: "未登录用户",
			setupContext: func(c *gin.Context) {
				// 不设置userId
			},
			expectedCanAccess: false,
		},
		{
			name: "VIP用户未过期",
			setupContext: func(c *gin.Context) {
				c.Set("userId", "user123")
				c.Set("vipUser", &VIPUser{
					UserID:    "user123",
					VIPLevel:  VIPLevelPlus,
					VIPExpire: now.Add(30 * 24 * time.Hour),
				})
			},
			expectedCanAccess: true,
		},
		{
			name: "VIP用户已过期",
			setupContext: func(c *gin.Context) {
				c.Set("userId", "user123")
				c.Set("vipUser", &VIPUser{
					UserID:    "user123",
					VIPLevel:  VIPLevelPlus,
					VIPExpire: now.Add(-24 * time.Hour),
				})
			},
			expectedCanAccess: false,
		},
		{
			name: "试用用户未过期",
			setupContext: func(c *gin.Context) {
				c.Set("userId", "user123")
				c.Set("vipUser", &VIPUser{
					UserID:      "user123",
					VIPLevel:    VIPLevelNone,
					IsTrialUser: true,
					TrialExpire: now.Add(7 * 24 * time.Hour),
				})
			},
			expectedCanAccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			tt.setupContext(c)

			middleware := CheckChapterVIPAccess()
			middleware(c)

			canAccess, exists := c.Get("canAccessVIPChapter")
			assert.True(t, exists)
			assert.Equal(t, tt.expectedCanAccess, canAccess)
		})
	}
}

func TestVIPRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name              string
		vipLevel          VIPLevel
		expectedRateLimit int
	}{
		{"普通用户", VIPLevelNone, 100},
		{"基础VIP", VIPLevelBasic, 500},
		{"VIP Plus", VIPLevelPlus, 1000},
		{"VIP Pro", VIPLevelPro, 5000},
		{"VIP Ultra", VIPLevelUltra, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)
			c.Set("vipLevel", tt.vipLevel)

			middleware := VIPRateLimiter()
			middleware(c)

			rateLimit, exists := c.Get("rateLimit")
			assert.True(t, exists)
			assert.Equal(t, tt.expectedRateLimit, rateLimit)
		})
	}
}
