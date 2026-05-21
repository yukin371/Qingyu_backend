package ai

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuotaPolicyStatus 配额策略状态
type QuotaPolicyStatus string

const (
	QuotaPolicyStatusActive   QuotaPolicyStatus = "active"   // 激活
	QuotaPolicyStatusDisabled QuotaPolicyStatus = "disabled" // 禁用
)

// UserRole 用户角色
type UserRole string

const (
	UserRoleReader UserRole = "reader" // 读者
	UserRoleWriter UserRole = "writer" // 作者
	UserRoleAdmin  UserRole = "admin"  // 管理员
)

// MembershipLevel 会员等级
type MembershipLevel string

const (
	MembershipLevelNormal      MembershipLevel = "normal"       // 普通用户
	MembershipLevelVipMonthly  MembershipLevel = "vip_monthly"  // 月度VIP
	MembershipLevelVipYearly   MembershipLevel = "vip_yearly"   // 年度VIP
	MembershipLevelSuperVip    MembershipLevel = "super_vip"    // 超级VIP
)

// QuotaPolicy 配额策略模型
type QuotaPolicy struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`                                 // 策略名称
	UserRole        UserRole           `json:"userRole" bson:"user_role"`                         // 用户角色
	MembershipLevel MembershipLevel    `json:"membershipLevel" bson:"membership_level"`            // 会员等级
	DailyQuota      int                `json:"dailyQuota" bson:"daily_quota"`                     // 日配额
	MonthlyQuota    int                `json:"monthlyQuota" bson:"monthly_quota"`                 // 月配额
	TotalQuota      int                `json:"totalQuota" bson:"total_quota"`                     // 总配额（-1=无限）
	IsDefault       bool               `json:"isDefault" bson:"is_default"`                       // 是否为默认策略
	Status          QuotaPolicyStatus  `json:"status" bson:"status"`                              // 策略状态
	Description     string             `json:"description,omitempty" bson:"description,omitempty"` // 描述
	CreatedAt       time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updated_at"`
}

// CollectionName 指定集合名
func (QuotaPolicy) CollectionName() string {
	return "ai_quota_policies"
}

// BeforeCreate MongoDB钩子 - 创建前
func (p *QuotaPolicy) BeforeCreate() {
	if p.ID.IsZero() {
		p.ID = primitive.NewObjectID()
	}
	if p.Status == "" {
		p.Status = QuotaPolicyStatusActive
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// BeforeUpdate MongoDB钩子 - 更新前
func (p *QuotaPolicy) BeforeUpdate() {
	p.UpdatedAt = time.Now()
}

// IsActive 检查策略是否激活
func (p *QuotaPolicy) IsActive() bool {
	return p.Status == QuotaPolicyStatusActive
}

// IsUnlimited 检查配额是否无限
func (p *QuotaPolicy) IsUnlimited() bool {
	return p.TotalQuota == -1
}
