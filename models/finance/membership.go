package finance

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MembershipPlan 会员套餐
type MembershipPlan struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`                     // 套餐名称：月卡、季卡、年卡、超级VIP
	Type        string             `bson:"type" json:"type"`                     // 套餐类型: monthly, quarterly, yearly, super
	Duration    int                `bson:"duration" json:"duration"`             // 有效期（天）
	Price       float64            `bson:"price" json:"price"`                   // 价格（元）
	OriginalPrice float64          `bson:"original_price" json:"original_price"` // 原价（元）
	Discount    float64            `bson:"discount" json:"discount"`             // 折扣
	Benefits    []string           `bson:"benefits" json:"benefits"`             // 权益列表
	Description string             `bson:"description" json:"description"`       // 套餐描述
	IsEnabled   bool               `bson:"is_enabled" json:"is_enabled"`         // 是否启用
	SortOrder   int                `bson:"sort_order" json:"sort_order"`         // 排序顺序
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// UserMembership 用户会员信息
type UserMembership struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           string             `bson:"user_id" json:"user_id"`
	PlanID           primitive.ObjectID `bson:"plan_id" json:"plan_id"`
	PlanName         string             `bson:"plan_name" json:"plan_name"`
	PlanType         string             `bson:"plan_type" json:"plan_type"`         // 会员类型
	Level            string             `bson:"level" json:"level"`                 // 等级：normal, vip_monthly, vip_yearly, super_vip
	StartTime        time.Time          `bson:"start_time" json:"start_time"`       // 开始时间
	EndTime          time.Time          `bson:"end_time" json:"end_time"`           // 结束时间
	AutoRenew        bool               `bson:"auto_renew" json:"auto_renew"`       // 是否自动续费
	Status           string             `bson:"status" json:"status"`               // 状态：active, expired, cancelled
	PaymentID        primitive.ObjectID `bson:"payment_id,omitempty" json:"payment_id,omitempty"` // 支付记录ID
	ActivatedAt      time.Time          `bson:"activated_at" json:"activated_at"`   // 激活时间
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

// MembershipBenefit 会员权益定义
type MembershipBenefit struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Level       string             `bson:"level" json:"level"`           // 权益等级：normal, vip_monthly, vip_yearly, super_vip
	Code        string             `bson:"code" json:"code"`             // 权益代码
	Name        string             `bson:"name" json:"name"`             // 权益名称
	Description string             `bson:"description" json:"description"` // 权益描述
	Value       string             `bson:"value" json:"value"`           // 权益值
	Category    string             `bson:"category" json:"category"`     // 权益类别：reading, writing, ai, social
	IsEnabled   bool               `bson:"is_enabled" json:"is_enabled"` // 是否启用
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// MembershipCard 会员卡
type MembershipCard struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code           string             `bson:"code" json:"code"`                     // 卡密
	PlanID         primitive.ObjectID `bson:"plan_id" json:"plan_id"`               // 关联套餐ID
	PlanType       string             `bson:"plan_type" json:"plan_type"`           // 套餐类型
	Duration       int                `bson:"duration" json:"duration"`             // 有效期（天）
	BatchID        string             `bson:"batch_id" json:"batch_id"`             // 批次ID
	Status         string             `bson:"status" json:"status"`                 // 状态：unused, used, expired, disabled
	ActivatedBy    string             `bson:"activated_by,omitempty" json:"activated_by,omitempty"` // 激活用户ID
	ActivatedAt    *time.Time         `bson:"activated_at,omitempty" json:"activated_at,omitempty"` // 激活时间
	ExpireAt       *time.Time         `bson:"expire_at,omitempty" json:"expire_at,omitempty"`       // 卡过期时间
	CreatedBy      string             `bson:"created_by" json:"created_by"`         // 创建人
	Note           string             `bson:"note,omitempty" json:"note,omitempty"` // 备注
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

// MembershipUsage 会员权益使用情况
type MembershipUsage struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         string             `bson:"user_id" json:"user_id"`
	BenefitCode    string             `bson:"benefit_code" json:"benefit_code"`   // 权益代码
	UsageCount     int                `bson:"usage_count" json:"usage_count"`     // 使用次数
	LastUsedAt     *time.Time         `bson:"last_used_at,omitempty" json:"last_used_at,omitempty"` // 最后使用时间
	Period         string             `bson:"period" json:"period"`               // 周期：daily, monthly, yearly
	PeriodStart    time.Time          `bson:"period_start" json:"period_start"`   // 周期开始
	PeriodEnd      time.Time          `bson:"period_end" json:"period_end"`       // 周期结束
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

// MembershipLevel 会员等级枚举
const (
	MembershipLevelNormal      = "normal"       // 普通用户
	MembershipLevelVIPMonthly  = "vip_monthly"  // VIP月卡
	MembershipLevelVIPYearly   = "vip_yearly"   // VIP年卡
	MembershipLevelSuperVIP    = "super_vip"    // 超级VIP
)

// MembershipStatus 会员状态枚举
const (
	MembershipStatusActive    = "active"    // 激活
	MembershipStatusExpired   = "expired"   // 过期
	MembershipStatusCancelled = "cancelled" // 取消
)

// MembershipCardStatus 会员卡状态枚举
const (
	CardStatusUnused   = "unused"   // 未使用
	CardStatusUsed     = "used"     // 已使用
	CardStatusExpired  = "expired"  // 已过期
	CardStatusDisabled = "disabled" // 已禁用
)

// MembershipType 会员类型枚举
const (
	MembershipTypeMonthly   = "monthly"   // 月卡
	MembershipTypeQuarterly = "quarterly" // 季卡
	MembershipTypeYearly    = "yearly"    // 年卡
	MembershipTypeSuper     = "super"     // 超级VIP
)

// BenefitCategory 权益类别枚举
const (
	BenefitCategoryReading = "reading" // 阅读权益
	BenefitCategoryWriting = "writing" // 写作权益
	BenefitCategoryAI      = "ai"      // AI权益
	BenefitCategorySocial  = "social"  // 社交权益
)
