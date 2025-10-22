package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// VIPLevel VIP等级
type VIPLevel int

const (
	VIPLevelNone  VIPLevel = 0 // 非VIP
	VIPLevelBasic VIPLevel = 1 // 基础VIP
	VIPLevelPlus  VIPLevel = 2 // VIP Plus
	VIPLevelPro   VIPLevel = 3 // VIP Pro
	VIPLevelUltra VIPLevel = 4 // VIP Ultra
)

// VIPConfig VIP配置
type VIPConfig struct {
	RequiredLevel VIPLevel // 最低VIP等级要求
	AllowTrial    bool     // 是否允许试用期用户
	Message       string   // 拒绝访问时的提示消息
}

// VIPUser VIP用户信息
type VIPUser struct {
	UserID      string    `json:"userId"`
	VIPLevel    VIPLevel  `json:"vipLevel"`
	VIPExpire   time.Time `json:"vipExpire"`
	IsTrialUser bool      `json:"isTrialUser"`
	TrialExpire time.Time `json:"trialExpire"`
	Permissions []string  `json:"permissions"`
}

// VIPService VIP服务接口
type VIPService interface {
	GetUserVIPInfo(ctx context.Context, userID string) (*VIPUser, error)
	CheckVIPPermission(ctx context.Context, userID string, requiredLevel VIPLevel) (bool, error)
}

// RequireVIP 需要VIP权限的中间件
func RequireVIP(requiredLevel VIPLevel) gin.HandlerFunc {
	return RequireVIPWithConfig(VIPConfig{
		RequiredLevel: requiredLevel,
		AllowTrial:    false,
		Message:       "该功能需要VIP权限",
	})
}

// RequireVIPWithConfig 带配置的VIP权限中间件
func RequireVIPWithConfig(config VIPConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID
		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":      40101,
				"message":   "用户未认证",
				"timestamp": time.Now().Unix(),
				"data":      nil,
			})
			c.Abort()
			return
		}

		// 2. 获取VIP信息（从context或数据库）
		vipUser, err := getUserVIPInfo(c, userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":      50001,
				"message":   "获取VIP信息失败",
				"timestamp": time.Now().Unix(),
				"data":      nil,
			})
			c.Abort()
			return
		}

		// 3. 检查VIP等级
		if !checkVIPLevel(vipUser, config) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    40302,
				"message": config.Message,
				"data": gin.H{
					"requiredLevel": config.RequiredLevel,
					"currentLevel":  vipUser.VIPLevel,
					"upgradeUrl":    "/vip/upgrade",
				},
				"timestamp": time.Now().Unix(),
			})
			c.Abort()
			return
		}

		// 4. 将VIP信息存入context
		c.Set("vipUser", vipUser)
		c.Set("vipLevel", vipUser.VIPLevel)

		c.Next()
	}
}

// RequireVIPPlus 需要VIP Plus及以上等级
func RequireVIPPlus() gin.HandlerFunc {
	return RequireVIP(VIPLevelPlus)
}

// RequireVIPPro 需要VIP Pro及以上等级
func RequireVIPPro() gin.HandlerFunc {
	return RequireVIP(VIPLevelPro)
}

// AllowVIPTrial 允许试用用户的VIP中间件
func AllowVIPTrial(requiredLevel VIPLevel) gin.HandlerFunc {
	return RequireVIPWithConfig(VIPConfig{
		RequiredLevel: requiredLevel,
		AllowTrial:    true,
		Message:       "该功能需要VIP权限或试用期",
	})
}

// CheckChapterVIPAccess 检查章节VIP访问权限
func CheckChapterVIPAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID
		userID, exists := c.Get("userId")
		if !exists {
			// 未登录用户只能访问免费章节
			c.Set("canAccessVIPChapter", false)
			c.Next()
			return
		}

		// 2. 获取VIP信息
		vipUser, err := getUserVIPInfo(c, userID.(string))
		if err != nil {
			c.Set("canAccessVIPChapter", false)
			c.Next()
			return
		}

		// 3. 检查VIP状态
		canAccess := false
		now := time.Now()

		// VIP用户且未过期
		if vipUser.VIPLevel > VIPLevelNone && vipUser.VIPExpire.After(now) {
			canAccess = true
		}

		// 试用用户且未过期
		if vipUser.IsTrialUser && vipUser.TrialExpire.After(now) {
			canAccess = true
		}

		// 4. 设置访问权限标志
		c.Set("canAccessVIPChapter", canAccess)
		c.Set("vipUser", vipUser)

		c.Next()
	}
}

// getUserVIPInfo 获取用户VIP信息
func getUserVIPInfo(c *gin.Context, userID string) (*VIPUser, error) {
	// 先尝试从context获取缓存的VIP信息
	if vipUser, exists := c.Get("vipUser"); exists {
		if user, ok := vipUser.(*VIPUser); ok {
			return user, nil
		}
	}

	// TODO: 从数据库或缓存中获取VIP信息
	// 这里暂时返回模拟数据，实际应该调用VIP服务
	return &VIPUser{
		UserID:      userID,
		VIPLevel:    VIPLevelNone,
		VIPExpire:   time.Now().Add(-24 * time.Hour), // 已过期
		IsTrialUser: false,
		TrialExpire: time.Time{},
		Permissions: []string{},
	}, nil
}

// checkVIPLevel 检查VIP等级是否满足要求
func checkVIPLevel(user *VIPUser, config VIPConfig) bool {
	now := time.Now()

	// 检查VIP等级
	if user.VIPLevel >= config.RequiredLevel && user.VIPExpire.After(now) {
		return true
	}

	// 如果允许试用，检查试用状态
	if config.AllowTrial && user.IsTrialUser && user.TrialExpire.After(now) {
		return true
	}

	return false
}

// GetVIPLevelName 获取VIP等级名称
func GetVIPLevelName(level VIPLevel) string {
	switch level {
	case VIPLevelNone:
		return "普通用户"
	case VIPLevelBasic:
		return "基础VIP"
	case VIPLevelPlus:
		return "VIP Plus"
	case VIPLevelPro:
		return "VIP Pro"
	case VIPLevelUltra:
		return "VIP Ultra"
	default:
		return "未知等级"
	}
}

// VIPContentFilter VIP内容过滤中间件
// 根据用户VIP等级过滤返回的内容
func VIPContentFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取VIP等级
		vipLevel := VIPLevelNone
		if level, exists := c.Get("vipLevel"); exists {
			if l, ok := level.(VIPLevel); ok {
				vipLevel = l
			}
		}

		// 设置内容过滤级别
		c.Set("contentFilterLevel", vipLevel)

		c.Next()
	}
}

// CheckBookVIPRequirement 检查书籍的VIP要求
func CheckBookVIPRequirement(bookID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 从数据库获取书籍的VIP要求
		// 这里是示例实现

		// 假设从数据库查询到书籍需要VIP Plus
		requiredLevel := VIPLevelPlus

		// 获取用户VIP等级
		vipLevel := VIPLevelNone
		if level, exists := c.Get("vipLevel"); exists {
			if l, ok := level.(VIPLevel); ok {
				vipLevel = l
			}
		}

		// 检查权限
		if vipLevel < requiredLevel {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    40302,
				"message": "该书籍需要" + GetVIPLevelName(requiredLevel),
				"data": gin.H{
					"bookId":        bookID,
					"requiredLevel": requiredLevel,
					"currentLevel":  vipLevel,
					"upgradeUrl":    "/vip/upgrade",
				},
				"timestamp": time.Now().Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// VIPRateLimiter VIP用户限流中间件
// 根据VIP等级设置不同的限流策略
func VIPRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		vipLevel := VIPLevelNone
		if level, exists := c.Get("vipLevel"); exists {
			if l, ok := level.(VIPLevel); ok {
				vipLevel = l
			}
		}

		// 根据VIP等级设置不同的限流策略
		var rateLimit int
		switch vipLevel {
		case VIPLevelNone:
			rateLimit = 100 // 普通用户：100次/小时
		case VIPLevelBasic:
			rateLimit = 500 // 基础VIP：500次/小时
		case VIPLevelPlus:
			rateLimit = 1000 // VIP Plus：1000次/小时
		case VIPLevelPro:
			rateLimit = 5000 // VIP Pro：5000次/小时
		case VIPLevelUltra:
			rateLimit = -1 // VIP Ultra：无限制
		default:
			rateLimit = 100
		}

		c.Set("rateLimit", rateLimit)

		c.Next()
	}
}

// VIPFeatureGate VIP功能开关
// 根据VIP等级启用不同的功能
type VIPFeature string

const (
	FeatureOfflineReading  VIPFeature = "offline_reading"  // 离线阅读
	FeatureHighlightExport VIPFeature = "highlight_export" // 导出笔记
	FeatureMultiDevice     VIPFeature = "multi_device"     // 多设备同步
	FeatureAdvancedSearch  VIPFeature = "advanced_search"  // 高级搜索
	FeatureCustomTheme     VIPFeature = "custom_theme"     // 自定义主题
	FeatureAdFree          VIPFeature = "ad_free"          // 无广告
	FeatureEarlyAccess     VIPFeature = "early_access"     // 抢先阅读
	FeaturePremiumSupport  VIPFeature = "premium_support"  // 优先客服
)

// RequireFeature 需要特定功能权限
func RequireFeature(feature VIPFeature) gin.HandlerFunc {
	return func(c *gin.Context) {
		vipLevel := VIPLevelNone
		if level, exists := c.Get("vipLevel"); exists {
			if l, ok := level.(VIPLevel); ok {
				vipLevel = l
			}
		}

		requiredLevel := getFeatureRequiredLevel(feature)

		if vipLevel < requiredLevel {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    40303,
				"message": "该功能需要" + GetVIPLevelName(requiredLevel),
				"data": gin.H{
					"feature":       feature,
					"requiredLevel": requiredLevel,
					"currentLevel":  vipLevel,
					"upgradeUrl":    "/vip/upgrade",
				},
				"timestamp": time.Now().Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getFeatureRequiredLevel 获取功能所需的VIP等级
func getFeatureRequiredLevel(feature VIPFeature) VIPLevel {
	featureMap := map[VIPFeature]VIPLevel{
		FeatureOfflineReading:  VIPLevelBasic,
		FeatureHighlightExport: VIPLevelBasic,
		FeatureMultiDevice:     VIPLevelPlus,
		FeatureAdvancedSearch:  VIPLevelPlus,
		FeatureCustomTheme:     VIPLevelPlus,
		FeatureAdFree:          VIPLevelPro,
		FeatureEarlyAccess:     VIPLevelPro,
		FeaturePremiumSupport:  VIPLevelUltra,
	}

	if level, exists := featureMap[feature]; exists {
		return level
	}

	return VIPLevelNone
}

// CheckVIPExpiration VIP过期检查中间件
func CheckVIPExpiration() gin.HandlerFunc {
	return func(c *gin.Context) {
		if vipUser, exists := c.Get("vipUser"); exists {
			if user, ok := vipUser.(*VIPUser); ok {
				now := time.Now()

				// 检查VIP是否即将过期（7天内）
				if user.VIPLevel > VIPLevelNone {
					expiresIn := user.VIPExpire.Sub(now)
					if expiresIn > 0 && expiresIn < 7*24*time.Hour {
						// VIP即将过期，添加提示信息到响应头
						c.Header("X-VIP-Expires-In", expiresIn.String())
						c.Header("X-VIP-Renewal-Url", "/vip/renewal")
					}
				}
			}
		}

		c.Next()
	}
}
