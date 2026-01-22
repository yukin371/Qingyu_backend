package finance

import (
	"Qingyu_backend/models/shared/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MembershipPlan 会员套餐
type MembershipPlan struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`                     // 套餐名称：月卡、季卡、年卡、超级VIP
	Type          string             `bson:"type" json:"type"`                     // 套餐类型: monthly, quarterly, yearly, super
	Duration      int                `bson:"duration" json:"duration"`             // 有效期（天）
	Price         types.Money        `bson:"price_cents" json:"-"`                // 价格（分）
	OriginalPrice types.Money        `bson:"original_price_cents" json:"-"`       // 原价（分）
	Discount      int                `bson:"discount_percent" json:"discount"`     // 折扣百分比（0-100）
	Benefits      []string           `bson:"benefits" json:"benefits"`             // 权益列表
	Description   string             `bson:"description" json:"description"`       // 套餐描述
	IsEnabled     bool               `bson:"is_enabled" json:"$1$2"`         // 是否启用
	SortOrder     int                `bson:"sort_order" json:"$1$2"`         // 排序顺序
	CreatedAt     time.Time          `bson:"created_at" json:"$1$2"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"$1$2"`
}

// UserMembership 用户会员信息
type UserMembership struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"$1$2"`
	PlanID      primitive.ObjectID `bson:"plan_id" json:"$1$2"`
	PlanName    string             `bson:"plan_name" json:"$1$2"`
	PlanType    string             `bson:"plan_type" json:"$1$2"`                       // 会员类型
	Level       string             `bson:"level" json:"level"`                               // 等级：normal, vip_monthly, vip_yearly, super_vip
	StartTime   time.Time          `bson:"start_time" json:"$1$2"`                     // 开始时间
	EndTime     time.Time          `bson:"end_time" json:"$1$2"`                         // 结束时间
	AutoRenew   bool               `bson:"auto_renew" json:"$1$2"`                     // 是否自动续费
	Status      string             `bson:"status" json:"status"`                             // 状态：active, expired, cancelled
	PaymentID   primitive.ObjectID `bson:"payment_id,omitempty" json:"payment_id,omitempty"` // 支付记录ID
	ActivatedAt time.Time          `bson:"activated_at" json:"$1$2"`                 // 激活时间
	CreatedAt   time.Time          `bson:"created_at" json:"$1$2"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"$1$2"`
}

// MembershipBenefit 会员权益定义
type MembershipBenefit struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Level       string             `bson:"level" json:"level"`             // 权益等级：normal, vip_monthly, vip_yearly, super_vip
	Code        string             `bson:"code" json:"code"`               // 权益代码
	Name        string             `bson:"name" json:"name"`               // 权益名称
	Description string             `bson:"description" json:"description"` // 权益描述
	Value       string             `bson:"value" json:"value"`             // 权益值
	Category    string             `bson:"category" json:"category"`       // 权益类别：reading, writing, ai, social
	IsEnabled   bool               `bson:"is_enabled" json:"$1$2"`   // 是否启用
	CreatedAt   time.Time          `bson:"created_at" json:"$1$2"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"$1$2"`
}

// MembershipCard 会员卡
type MembershipCard struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code        string             `bson:"code" json:"code"`                                     // 卡密
	PlanID      primitive.ObjectID `bson:"plan_id" json:"$1$2"`                               // 关联套餐ID
	PlanType    string             `bson:"plan_type" json:"$1$2"`                           // 套餐类型
	Duration    int                `bson:"duration" json:"duration"`                             // 有效期（天）
	BatchID     string             `bson:"batch_id" json:"$1$2"`                             // 批次ID
	Status      string             `bson:"status" json:"status"`                                 // 状态：unused, used, expired, disabled
	ActivatedBy string             `bson:"activated_by,omitempty" json:"activated_by,omitempty"` // 激活用户ID
	ActivatedAt *time.Time         `bson:"activated_at,omitempty" json:"activated_at,omitempty"` // 激活时间
	ExpireAt    *time.Time         `bson:"expire_at,omitempty" json:"expire_at,omitempty"`       // 卡过期时间
	CreatedBy   string             `bson:"created_by" json:"$1$2"`                         // 创建人
	Note        string             `bson:"note,omitempty" json:"note,omitempty"`                 // 备注
	CreatedAt   time.Time          `bson:"created_at" json:"$1$2"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"$1$2"`
}

// MembershipUsage 会员权益使用情况
type MembershipUsage struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"$1$2"`
	BenefitCode string             `bson:"benefit_code" json:"$1$2"`                     // 权益代码
	UsageCount  int                `bson:"usage_count" json:"$1$2"`                       // 使用次数
	LastUsedAt  *time.Time         `bson:"last_used_at,omitempty" json:"last_used_at,omitempty"` // 最后使用时间
	Period      string             `bson:"period" json:"period"`                                 // 周期：daily, monthly, yearly
	PeriodStart time.Time          `bson:"period_start" json:"$1$2"`                     // 周期开始
	PeriodEnd   time.Time          `bson:"period_end" json:"$1$2"`                         // 周期结束
	CreatedAt   time.Time          `bson:"created_at" json:"$1$2"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"$1$2"`
}

// MembershipLevel 会员等级枚举
const (
	MembershipLevelNormal     = "normal"      // 普通用户
	MembershipLevelVIPMonthly = "vip_monthly" // VIP月卡
	MembershipLevelVIPYearly  = "vip_yearly"  // VIP年卡
	MembershipLevelSuperVIP   = "super_vip"   // 超级VIP
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
