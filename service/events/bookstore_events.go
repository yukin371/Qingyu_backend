package events

import (
	"context"
	"log"
	"time"

	"Qingyu_backend/service/base"
)

// ============ 书城相关事件 ============

// 书城事件类型常量
const (
	// 购买事件
	EventBookPurchased    = "book.purchased"
	EventChapterPurchased = "chapter.purchased"
	EventRefundRequested  = "refund.requested"
	EventRefundApproved   = "refund.approved"
	EventRefundRejected   = "refund.rejected"

	// 订阅事件
	EventSubscriptionCreated   = "subscription.created"
	EventSubscriptionRenewed   = "subscription.renewed"
	EventSubscriptionExpired   = "subscription.expired"
	EventSubscriptionCancelled = "subscription.cancelled"

	// 打赏事件
	EventRewardCreated  = "reward.created"
	EventRewardReceived = "reward.received"

	// VIP事件
	EventVIPPurchased = "vip.purchased"
	EventVIPActivated = "vip.activated"
	EventVIPExpired   = "vip.expired"
)

// BookstoreEventData 书城事件数据
type BookstoreEventData struct {
	UserID    string                 `json:"user_id"`
	BookID    string                 `json:"book_id,omitempty"`
	ChapterID string                 `json:"chapter_id,omitempty"`
	Amount    float64                `json:"amount,omitempty"`
	Currency  string                 `json:"currency,omitempty"`
	Action    string                 `json:"action"`
	Time      time.Time              `json:"time"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ============ 购买事件 ============

// PurchaseEventData 购买事件数据
type PurchaseEventData struct {
	BookstoreEventData
	OrderID       string  `json:"order_id"`
	PaymentMethod string  `json:"payment_method"`
	Discount      float64 `json:"discount,omitempty"`
	FinalAmount   float64 `json:"final_amount"`
}

// NewBookPurchasedEvent 创建书籍购买事件
func NewBookPurchasedEvent(userID, bookID, orderID, paymentMethod string, amount, discount, finalAmount float64) base.Event {
	return &base.BaseEvent{
		EventType: EventBookPurchased,
		EventData: PurchaseEventData{
			BookstoreEventData: BookstoreEventData{
				UserID:   userID,
				BookID:   bookID,
				Amount:   amount,
				Currency: "CNY",
				Action:   "purchased",
				Time:     time.Now(),
			},
			OrderID:       orderID,
			PaymentMethod: paymentMethod,
			Discount:      discount,
			FinalAmount:   finalAmount,
		},
		Timestamp: time.Now(),
		Source:    "BookstoreService",
	}
}

// NewChapterPurchasedEvent 创建章节购买事件
func NewChapterPurchasedEvent(userID, bookID, chapterID, orderID, paymentMethod string, amount float64) base.Event {
	return &base.BaseEvent{
		EventType: EventChapterPurchased,
		EventData: PurchaseEventData{
			BookstoreEventData: BookstoreEventData{
				UserID:    userID,
				BookID:    bookID,
				ChapterID: chapterID,
				Amount:    amount,
				Currency:  "CNY",
				Action:    "chapter_purchased",
				Time:      time.Now(),
			},
			OrderID:       orderID,
			PaymentMethod: paymentMethod,
			FinalAmount:   amount,
		},
		Timestamp: time.Now(),
		Source:    "BookstoreService",
	}
}

// RefundEventData 退款事件数据
type RefundEventData struct {
	BookstoreEventData
	OrderID       string    `json:"order_id"`
	RefundID      string    `json:"refund_id"`
	Reason        string    `json:"reason"`
	RefundAmount  float64   `json:"refund_amount"`
	Status        string    `json:"status"`
	ProcessedTime time.Time `json:"processed_time,omitempty"`
}

// NewRefundRequestedEvent 创建退款请求事件
func NewRefundRequestedEvent(userID, orderID, reason string, amount float64) base.Event {
	return &base.BaseEvent{
		EventType: EventRefundRequested,
		EventData: RefundEventData{
			BookstoreEventData: BookstoreEventData{
				UserID:   userID,
				Amount:   amount,
				Currency: "CNY",
				Action:   "refund_requested",
				Time:     time.Now(),
			},
			OrderID:      orderID,
			Reason:       reason,
			RefundAmount: amount,
			Status:       "pending",
		},
		Timestamp: time.Now(),
		Source:    "BookstoreService",
	}
}

// NewRefundApprovedEvent 创建退款批准事件
func NewRefundApprovedEvent(userID, orderID, refundID string, amount float64) base.Event {
	return &base.BaseEvent{
		EventType: EventRefundApproved,
		EventData: RefundEventData{
			BookstoreEventData: BookstoreEventData{
				UserID:   userID,
				Amount:   amount,
				Currency: "CNY",
				Action:   "refund_approved",
				Time:     time.Now(),
			},
			OrderID:       orderID,
			RefundID:      refundID,
			RefundAmount:  amount,
			Status:        "approved",
			ProcessedTime: time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "BookstoreService",
	}
}

// ============ 订阅事件 ============

// SubscriptionEventData 订阅事件数据
type SubscriptionEventData struct {
	UserID         string                 `json:"user_id"`
	SubscriptionID string                 `json:"subscription_id"`
	PlanType       string                 `json:"plan_type"` // monthly/yearly
	Amount         float64                `json:"amount"`
	Currency       string                 `json:"currency"`
	Action         string                 `json:"action"`
	Time           time.Time              `json:"time"`
	StartTime      time.Time              `json:"start_time,omitempty"`
	EndTime        time.Time              `json:"end_time,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// NewSubscriptionCreatedEvent 创建订阅事件
func NewSubscriptionCreatedEvent(userID, subscriptionID, planType string, amount float64, startTime, endTime time.Time) base.Event {
	return &base.BaseEvent{
		EventType: EventSubscriptionCreated,
		EventData: SubscriptionEventData{
			UserID:         userID,
			SubscriptionID: subscriptionID,
			PlanType:       planType,
			Amount:         amount,
			Currency:       "CNY",
			Action:         "subscribed",
			Time:           time.Now(),
			StartTime:      startTime,
			EndTime:        endTime,
		},
		Timestamp: time.Now(),
		Source:    "BookstoreService",
	}
}

// NewSubscriptionRenewedEvent 创建订阅续费事件
func NewSubscriptionRenewedEvent(userID, subscriptionID, planType string, amount float64, endTime time.Time) base.Event {
	return &base.BaseEvent{
		EventType: EventSubscriptionRenewed,
		EventData: SubscriptionEventData{
			UserID:         userID,
			SubscriptionID: subscriptionID,
			PlanType:       planType,
			Amount:         amount,
			Currency:       "CNY",
			Action:         "renewed",
			Time:           time.Now(),
			EndTime:        endTime,
		},
		Timestamp: time.Now(),
		Source:    "BookstoreService",
	}
}

// NewSubscriptionExpiredEvent 创建订阅过期事件
func NewSubscriptionExpiredEvent(userID, subscriptionID string) base.Event {
	return &base.BaseEvent{
		EventType: EventSubscriptionExpired,
		EventData: SubscriptionEventData{
			UserID:         userID,
			SubscriptionID: subscriptionID,
			Action:         "expired",
			Time:           time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "BookstoreService",
	}
}

// ============ 打赏事件 ============

// RewardEventData 打赏事件数据
type RewardEventData struct {
	RewardID  string                 `json:"reward_id"`
	SponsorID string                 `json:"sponsor_id"` // 打赏者ID
	AuthorID  string                 `json:"author_id"`  // 作者ID
	BookID    string                 `json:"book_id"`
	ChapterID string                 `json:"chapter_id,omitempty"`
	Amount    float64                `json:"amount"`
	Currency  string                 `json:"currency"`
	Message   string                 `json:"message,omitempty"`
	Action    string                 `json:"action"`
	Time      time.Time              `json:"time"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// NewRewardCreatedEvent 创建打赏事件
func NewRewardCreatedEvent(rewardID, sponsorID, authorID, bookID, chapterID, message string, amount float64) base.Event {
	return &base.BaseEvent{
		EventType: EventRewardCreated,
		EventData: RewardEventData{
			RewardID:  rewardID,
			SponsorID: sponsorID,
			AuthorID:  authorID,
			BookID:    bookID,
			ChapterID: chapterID,
			Amount:    amount,
			Currency:  "CNY",
			Message:   message,
			Action:    "rewarded",
			Time:      time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "BookstoreService",
	}
}

// ============ VIP事件 ============

// VIPEventData VIP事件数据
type VIPEventData struct {
	UserID    string                 `json:"user_id"`
	VIPLevel  string                 `json:"vip_level"` // basic/silver/gold/platinum
	Duration  int                    `json:"duration"`  // 天数
	Amount    float64                `json:"amount,omitempty"`
	Currency  string                 `json:"currency,omitempty"`
	Action    string                 `json:"action"`
	Time      time.Time              `json:"time"`
	StartTime time.Time              `json:"start_time,omitempty"`
	EndTime   time.Time              `json:"end_time,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// NewVIPPurchasedEvent 创建VIP购买事件
func NewVIPPurchasedEvent(userID, vipLevel string, duration int, amount float64, startTime, endTime time.Time) base.Event {
	return &base.BaseEvent{
		EventType: EventVIPPurchased,
		EventData: VIPEventData{
			UserID:    userID,
			VIPLevel:  vipLevel,
			Duration:  duration,
			Amount:    amount,
			Currency:  "CNY",
			Action:    "vip_purchased",
			Time:      time.Now(),
			StartTime: startTime,
			EndTime:   endTime,
		},
		Timestamp: time.Now(),
		Source:    "BookstoreService",
	}
}

// NewVIPActivatedEvent 创建VIP激活事件
func NewVIPActivatedEvent(userID, vipLevel string, startTime, endTime time.Time) base.Event {
	return &base.BaseEvent{
		EventType: EventVIPActivated,
		EventData: VIPEventData{
			UserID:    userID,
			VIPLevel:  vipLevel,
			Action:    "vip_activated",
			Time:      time.Now(),
			StartTime: startTime,
			EndTime:   endTime,
		},
		Timestamp: time.Now(),
		Source:    "BookstoreService",
	}
}

// NewVIPExpiredEvent 创建VIP过期事件
func NewVIPExpiredEvent(userID, vipLevel string) base.Event {
	return &base.BaseEvent{
		EventType: EventVIPExpired,
		EventData: VIPEventData{
			UserID:   userID,
			VIPLevel: vipLevel,
			Action:   "vip_expired",
			Time:     time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "BookstoreService",
	}
}

// ============ 事件处理器 ============

// PurchaseStatisticsHandler 购买统计处理器
// 更新购买统计和畅销榜单
type PurchaseStatisticsHandler struct {
	name string
}

// NewPurchaseStatisticsHandler 创建购买统计处理器
func NewPurchaseStatisticsHandler() *PurchaseStatisticsHandler {
	return &PurchaseStatisticsHandler{
		name: "PurchaseStatisticsHandler",
	}
}

// Handle 处理事件
func (h *PurchaseStatisticsHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventBookPurchased:
		data, _ := event.GetEventData().(PurchaseEventData)
		log.Printf("[PurchaseStatistics] 书籍 %s 被购买，金额: %.2f", data.BookID, data.FinalAmount)
		// 更新书籍购买统计
		// 更新畅销榜单

	case EventChapterPurchased:
		data, _ := event.GetEventData().(PurchaseEventData)
		log.Printf("[PurchaseStatistics] 章节 %s 被购买，金额: %.2f", data.ChapterID, data.FinalAmount)
		// 更新章节购买统计
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *PurchaseStatisticsHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *PurchaseStatisticsHandler) GetSupportedEventTypes() []string {
	return []string{
		EventBookPurchased,
		EventChapterPurchased,
	}
}

// RoyaltyCalculationHandler 版税计算处理器
// 计算作者版税
type RoyaltyCalculationHandler struct {
	name string
}

// NewRoyaltyCalculationHandler 创建版税计算处理器
func NewRoyaltyCalculationHandler() *RoyaltyCalculationHandler {
	return &RoyaltyCalculationHandler{
		name: "RoyaltyCalculationHandler",
	}
}

// Handle 处理事件
func (h *RoyaltyCalculationHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventBookPurchased:
		data, _ := event.GetEventData().(PurchaseEventData)
		log.Printf("[RoyaltyCalculation] 计算书籍 %s 的版税，金额: %.2f", data.BookID, data.FinalAmount)
		// 计算作者版税（通常是销售额的50-70%）
		// 记录到作者收入账户

	case EventChapterPurchased:
		data, _ := event.GetEventData().(PurchaseEventData)
		log.Printf("[RoyaltyCalculation] 计算章节 %s 的版税，金额: %.2f", data.ChapterID, data.FinalAmount)
		// 计算章节版税

	case EventRewardCreated:
		data, _ := event.GetEventData().(RewardEventData)
		log.Printf("[RoyaltyCalculation] 记录打赏 %s 给作者 %s，金额: %.2f", data.RewardID, data.AuthorID, data.Amount)
		// 打赏金额通常全部给作者
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *RoyaltyCalculationHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *RoyaltyCalculationHandler) GetSupportedEventTypes() []string {
	return []string{
		EventBookPurchased,
		EventChapterPurchased,
		EventRewardCreated,
	}
}

// VIPBenefitHandler VIP权益处理器
// 激活和管理VIP权益
type VIPBenefitHandler struct {
	name string
}

// NewVIPBenefitHandler 创建VIP权益处理器
func NewVIPBenefitHandler() *VIPBenefitHandler {
	return &VIPBenefitHandler{
		name: "VIPBenefitHandler",
	}
}

// Handle 处理事件
func (h *VIPBenefitHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventVIPActivated:
		data, _ := event.GetEventData().(VIPEventData)
		log.Printf("[VIPBenefit] 激活用户 %s 的VIP权益，等级: %s", data.UserID, data.VIPLevel)
		// 激活VIP权益
		// 发送VIP激活通知

	case EventVIPExpired:
		data, _ := event.GetEventData().(VIPEventData)
		log.Printf("[VIPBenefit] 用户 %s 的VIP权益已过期", data.UserID)
		// 停用VIP权益
		// 发送VIP过期通知
		// 清理VIP相关缓存
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *VIPBenefitHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *VIPBenefitHandler) GetSupportedEventTypes() []string {
	return []string{
		EventVIPActivated,
		EventVIPExpired,
	}
}

// SubscriptionNotificationHandler 订阅通知处理器
// 发送订阅相关通知
type SubscriptionNotificationHandler struct {
	name string
}

// NewSubscriptionNotificationHandler 创建订阅通知处理器
func NewSubscriptionNotificationHandler() *SubscriptionNotificationHandler {
	return &SubscriptionNotificationHandler{
		name: "SubscriptionNotificationHandler",
	}
}

// Handle 处理事件
func (h *SubscriptionNotificationHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventSubscriptionCreated:
		data, _ := event.GetEventData().(SubscriptionEventData)
		log.Printf("[SubscriptionNotification] 用户 %s 订阅成功，%s订阅", data.UserID, data.PlanType)
		// 发送订阅成功通知

	case EventSubscriptionRenewed:
		data, _ := event.GetEventData().(SubscriptionEventData)
		log.Printf("[SubscriptionNotification] 用户 %s 续费成功，%s订阅", data.UserID, data.PlanType)
		// 发送续费成功通知

	case EventSubscriptionExpired:
		data, _ := event.GetEventData().(SubscriptionEventData)
		log.Printf("[SubscriptionNotification] 用户 %s 订阅已过期", data.UserID)
		// 发送订阅过期通知
		// 提示续费
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *SubscriptionNotificationHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *SubscriptionNotificationHandler) GetSupportedEventTypes() []string {
	return []string{
		EventSubscriptionCreated,
		EventSubscriptionRenewed,
		EventSubscriptionExpired,
	}
}
