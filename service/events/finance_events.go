package events

import (
	"context"
	"log"
	"time"

	"Qingyu_backend/service/base"
)

// ============ 财务相关事件 ============

// 财务事件类型常量
const (
	// 充值事件
	EventDepositCreated   = "deposit.created"
	EventDepositCompleted = "deposit.completed"
	EventDepositFailed    = "deposit.failed"

	// 提现事件
	EventWithdrawalCreated   = "withdrawal.created"
	EventWithdrawalApproved  = "withdrawal.approved"
	EventWithdrawalRejected  = "withdrawal.rejected"
	EventWithdrawalCompleted = "withdrawal.completed"

	// 结算事件
	EventSettlementGenerated = "settlement.generated"
	EventSettlementPaid      = "settlement.paid"

	// 收入事件
	EventRevenueEarned  = "revenue.earned"
	EventRevenueSettled = "revenue.settled"
)

// FinanceEventData 财务事件数据
type FinanceEventData struct {
	UserID    string                 `json:"user_id"`
	Amount    float64                `json:"amount"`
	Currency  string                 `json:"currency"`
	Status    string                 `json:"status"`
	Action    string                 `json:"action"`
	Time      time.Time              `json:"time"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ============ 充值事件 ============

// DepositEventData 充值事件数据
type DepositEventData struct {
	DepositID     string                 `json:"deposit_id"`
	UserID        string                 `json:"user_id"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	PaymentMethod string                 `json:"payment_method"`
	Status        string                 `json:"status"`        // pending/completed/failed
	Action        string                 `json:"action"`
	Time          time.Time              `json:"time"`
	CompletedTime time.Time              `json:"completed_time,omitempty"`
	FailureReason string                 `json:"failure_reason,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// NewDepositCreatedEvent 创建充值请求事件
func NewDepositCreatedEvent(depositID, userID string, amount float64, paymentMethod string) base.Event {
	return &base.BaseEvent{
		EventType: EventDepositCreated,
		EventData: DepositEventData{
			DepositID:     depositID,
			UserID:        userID,
			Amount:        amount,
			Currency:      "CNY",
			PaymentMethod: paymentMethod,
			Status:        "pending",
			Action:        "deposit_requested",
			Time:          time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "FinanceService",
	}
}

// NewDepositCompletedEvent 创建充值完成事件
func NewDepositCompletedEvent(depositID, userID string, amount float64) base.Event {
	return &base.BaseEvent{
		EventType: EventDepositCompleted,
		EventData: DepositEventData{
			DepositID:     depositID,
			UserID:        userID,
			Amount:        amount,
			Currency:      "CNY",
			Status:        "completed",
			Action:        "deposit_completed",
			Time:          time.Now(),
			CompletedTime: time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "FinanceService",
	}
}

// NewDepositFailedEvent 创建充值失败事件
func NewDepositFailedEvent(depositID, userID string, amount float64, reason string) base.Event {
	return &base.BaseEvent{
		EventType: EventDepositFailed,
		EventData: DepositEventData{
			DepositID:     depositID,
			UserID:        userID,
			Amount:        amount,
			Currency:      "CNY",
			Status:        "failed",
			Action:        "deposit_failed",
			Time:          time.Now(),
			FailureReason: reason,
		},
		Timestamp: time.Now(),
		Source:    "FinanceService",
	}
}

// ============ 提现事件 ============

// WithdrawalEventData 提现事件数据
type WithdrawalEventData struct {
	WithdrawalID  string                 `json:"withdrawal_id"`
	UserID        string                 `json:"user_id"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	BankAccount   string                 `json:"bank_account"`
	BankName      string                 `json:"bank_name"`
	AccountName   string                 `json:"account_name"`
	Status        string                 `json:"status"`        // pending/approved/rejected/completed
	Action        string                 `json:"action"`
	Time          time.Time              `json:"time"`
	ProcessedTime time.Time              `json:"processed_time,omitempty"`
	ProcessedBy   string                 `json:"processed_by,omitempty"`
	RejectionReason string               `json:"rejection_reason,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// NewWithdrawalCreatedEvent 创建提现请求事件
func NewWithdrawalCreatedEvent(withdrawalID, userID string, amount float64, bankAccount, bankName, accountName string) base.Event {
	return &base.BaseEvent{
		EventType: EventWithdrawalCreated,
		EventData: WithdrawalEventData{
			WithdrawalID: withdrawalID,
			UserID:       userID,
			Amount:       amount,
			Currency:     "CNY",
			BankAccount:  bankAccount,
			BankName:     bankName,
			AccountName:  accountName,
			Status:       "pending",
			Action:       "withdrawal_requested",
			Time:         time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "FinanceService",
	}
}

// NewWithdrawalApprovedEvent 创建提现批准事件
func NewWithdrawalApprovedEvent(withdrawalID, userID string, amount float64, processedBy string) base.Event {
	return &base.BaseEvent{
		EventType: EventWithdrawalApproved,
		EventData: WithdrawalEventData{
			WithdrawalID:  withdrawalID,
			UserID:        userID,
			Amount:        amount,
			Currency:      "CNY",
			Status:        "approved",
			Action:        "withdrawal_approved",
			Time:          time.Now(),
			ProcessedTime: time.Now(),
			ProcessedBy:   processedBy,
		},
		Timestamp: time.Now(),
		Source:    "FinanceService",
	}
}

// NewWithdrawalRejectedEvent 创建提现拒绝事件
func NewWithdrawalRejectedEvent(withdrawalID, userID string, amount float64, reason, processedBy string) base.Event {
	return &base.BaseEvent{
		EventType: EventWithdrawalRejected,
		EventData: WithdrawalEventData{
			WithdrawalID:    withdrawalID,
			UserID:          userID,
			Amount:          amount,
			Currency:        "CNY",
			Status:          "rejected",
			Action:          "withdrawal_rejected",
			Time:            time.Now(),
			ProcessedTime:   time.Now(),
			ProcessedBy:     processedBy,
			RejectionReason: reason,
		},
		Timestamp: time.Now(),
		Source:    "FinanceService",
	}
}

// NewWithdrawalCompletedEvent 创建提现完成事件
func NewWithdrawalCompletedEvent(withdrawalID, userID string, amount float64) base.Event {
	return &base.BaseEvent{
		EventType: EventWithdrawalCompleted,
		EventData: WithdrawalEventData{
			WithdrawalID:  withdrawalID,
			UserID:        userID,
			Amount:        amount,
			Currency:      "CNY",
			Status:        "completed",
			Action:        "withdrawal_completed",
			Time:          time.Now(),
			ProcessedTime: time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "FinanceService",
	}
}

// ============ 结算事件 ============

// SettlementEventData 结算事件数据
type SettlementEventData struct {
	SettlementID   string                 `json:"settlement_id"`
	UserID         string                 `json:"user_id"`         // 作者ID
	Period         string                 `json:"period"`          // 结算周期，如 2024-01
	TotalRevenue   float64                `json:"total_revenue"`   // 总收入
	Commission     float64                `json:"commission"`      // 平台佣金
	NetRevenue     float64                `json:"net_revenue"`     // 净收入
	Tax            float64                `json:"tax"`             // 税费
	FinalAmount    float64                `json:"final_amount"`    // 最终结算金额
	Currency       string                 `json:"currency"`
	Status         string                 `json:"status"`          // generated/paid
	Action         string                 `json:"action"`
	Time           time.Time              `json:"time"`
	PaidTime       time.Time              `json:"paid_time,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// NewSettlementGeneratedEvent 创建结算单生成事件
func NewSettlementGeneratedEvent(settlementID, userID, period string, totalRevenue, commission, netRevenue, tax, finalAmount float64) base.Event {
	return &base.BaseEvent{
		EventType: EventSettlementGenerated,
		EventData: SettlementEventData{
			SettlementID: settlementID,
			UserID:       userID,
			Period:       period,
			TotalRevenue: totalRevenue,
			Commission:   commission,
			NetRevenue:   netRevenue,
			Tax:          tax,
			FinalAmount:  finalAmount,
			Currency:     "CNY",
			Status:       "generated",
			Action:       "settlement_generated",
			Time:         time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "FinanceService",
	}
}

// NewSettlementPaidEvent 创建结算支付完成事件
func NewSettlementPaidEvent(settlementID, userID string, amount float64) base.Event {
	return &base.BaseEvent{
		EventType: EventSettlementPaid,
		EventData: SettlementEventData{
			SettlementID: settlementID,
			UserID:       userID,
			FinalAmount:  amount,
			Currency:     "CNY",
			Status:       "paid",
			Action:       "settlement_paid",
			Time:         time.Now(),
			PaidTime:     time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "FinanceService",
	}
}

// ============ 收入事件 ============

// RevenueEventData 收入事件数据
type RevenueEventData struct {
	RevenueID      string                 `json:"revenue_id"`
	UserID         string                 `json:"user_id"`         // 作者ID
	Source         string                 `json:"source"`          // 收入来源: book_purchase/chapter_purchase/reward
	SourceID       string                 `json:"source_id"`       // 来源ID: bookID/chapterID/rewardID
	Amount         float64                `json:"amount"`
	Currency       string                 `json:"currency"`
	RoyaltyRate    float64                `json:"royalty_rate"`    // 版税率
	NetAmount      float64                `json:"net_amount"`      // 净收入
	Status         string                 `json:"status"`          // earned/settled
	Action         string                 `json:"action"`
	Time           time.Time              `json:"time"`
	SettledTime    time.Time              `json:"settled_time,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// NewRevenueEarnedEvent 创建收入事件
func NewRevenueEarnedEvent(revenueID, userID, source, sourceID string, amount, royaltyRate, netAmount float64) base.Event {
	return &base.BaseEvent{
		EventType: EventRevenueEarned,
		EventData: RevenueEventData{
			RevenueID:   revenueID,
			UserID:      userID,
			Source:      source,
			SourceID:    sourceID,
			Amount:      amount,
			Currency:    "CNY",
			RoyaltyRate: royaltyRate,
			NetAmount:   netAmount,
			Status:      "earned",
			Action:      "revenue_earned",
			Time:        time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "FinanceService",
	}
}

// NewRevenueSettledEvent 创建收入结算事件
func NewRevenueSettledEvent(revenueID, userID string, amount float64) base.Event {
	return &base.BaseEvent{
		EventType: EventRevenueSettled,
		EventData: RevenueEventData{
			RevenueID:   revenueID,
			UserID:      userID,
			NetAmount:   amount,
			Currency:    "CNY",
			Status:      "settled",
			Action:      "revenue_settled",
			Time:        time.Now(),
			SettledTime: time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "FinanceService",
	}
}

// ============ 事件处理器 ============

// TransactionHandler 交易处理器
// 处理充值和提现交易
type TransactionHandler struct {
	name string
}

// NewTransactionHandler 创建交易处理器
func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{
		name: "TransactionHandler",
	}
}

// Handle 处理事件
func (h *TransactionHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventDepositCreated:
		data, _ := event.GetEventData().(DepositEventData)
		log.Printf("[Transaction] 用户 %s 发起充值请求 %.2f 元，支付方式: %s", data.UserID, data.Amount, data.PaymentMethod)
		// 调用支付网关处理充值

	case EventDepositCompleted:
		data, _ := event.GetEventData().(DepositEventData)
		log.Printf("[Transaction] 用户 %s 充值成功，金额: %.2f 元", data.UserID, data.Amount)
		// 更新用户账户余额
		// 发送充值成功通知

	case EventDepositFailed:
		data, _ := event.GetEventData().(DepositEventData)
		log.Printf("[Transaction] 用户 %s 充值失败，原因: %s", data.UserID, data.FailureReason)
		// 发送充值失败通知

	case EventWithdrawalCreated:
		data, _ := event.GetEventData().(WithdrawalEventData)
		log.Printf("[Transaction] 用户 %s 发起提现请求 %.2f 元，银行: %s", data.UserID, data.Amount, data.BankName)
		// 通知管理员审核

	case EventWithdrawalApproved:
		data, _ := event.GetEventData().(WithdrawalEventData)
		log.Printf("[Transaction] 提现请求 %s 已批准，处理人: %s", data.WithdrawalID, data.ProcessedBy)
		// 执行银行转账
		// 冻结用户账户余额

	case EventWithdrawalRejected:
		data, _ := event.GetEventData().(WithdrawalEventData)
		log.Printf("[Transaction] 提现请求 %s 已拒绝，原因: %s", data.WithdrawalID, data.RejectionReason)
		// 发送拒绝通知

	case EventWithdrawalCompleted:
		data, _ := event.GetEventData().(WithdrawalEventData)
		log.Printf("[Transaction] 提现请求 %s 完成，金额: %.2f 元", data.WithdrawalID, data.Amount)
		// 扣除用户账户余额
		// 发送提现成功通知
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *TransactionHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *TransactionHandler) GetSupportedEventTypes() []string {
	return []string{
		EventDepositCreated,
		EventDepositCompleted,
		EventDepositFailed,
		EventWithdrawalCreated,
		EventWithdrawalApproved,
		EventWithdrawalRejected,
		EventWithdrawalCompleted,
	}
}

// SettlementHandler 结算处理器
// 处理作者结算
type SettlementHandler struct {
	name string
}

// NewSettlementHandler 创建结算处理器
func NewSettlementHandler() *SettlementHandler {
	return &SettlementHandler{
		name: "SettlementHandler",
	}
}

// Handle 处理事件
func (h *SettlementHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventSettlementGenerated:
		data, _ := event.GetEventData().(SettlementEventData)
		log.Printf("[Settlement] 生成用户 %s 的 %s 结算单，总金额: %.2f 元", data.UserID, data.Period, data.FinalAmount)
		// 生成结算单PDF
		// 发送结算单通知给作者
		// 等待作者确认

	case EventSettlementPaid:
		data, _ := event.GetEventData().(SettlementEventData)
		log.Printf("[Settlement] 结算单 %s 已支付，金额: %.2f 元", data.SettlementID, data.FinalAmount)
		// 更新结算状态
		// 发送支付成功通知
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *SettlementHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *SettlementHandler) GetSupportedEventTypes() []string {
	return []string{
		EventSettlementGenerated,
		EventSettlementPaid,
	}
}

// RevenueStatisticsHandler 收入统计处理器
// 统计和分析收入数据
type RevenueStatisticsHandler struct {
	name string
}

// NewRevenueStatisticsHandler 创建收入统计处理器
func NewRevenueStatisticsHandler() *RevenueStatisticsHandler {
	return &RevenueStatisticsHandler{
		name: "RevenueStatisticsHandler",
	}
}

// Handle 处理事件
func (h *RevenueStatisticsHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventRevenueEarned:
		data, _ := event.GetEventData().(RevenueEventData)
		log.Printf("[RevenueStatistics] 用户 %s 获得收入 %.2f 元，来源: %s:%s", data.UserID, data.NetAmount, data.Source, data.SourceID)
		// 更新作者收入统计
		// 更新作品收入排名
		// 生成收入报表

	case EventRevenueSettled:
		data, _ := event.GetEventData().(RevenueEventData)
		log.Printf("[RevenueStatistics] 用户 %s 收入 %.2f 元已结算", data.UserID, data.NetAmount)
		// 更新已结算收入统计
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *RevenueStatisticsHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *RevenueStatisticsHandler) GetSupportedEventTypes() []string {
	return []string{
		EventRevenueEarned,
		EventRevenueSettled,
	}
}

// FinanceNotificationHandler 财务通知处理器
// 发送财务相关通知
type FinanceNotificationHandler struct {
	name string
}

// NewFinanceNotificationHandler 创建财务通知处理器
func NewFinanceNotificationHandler() *FinanceNotificationHandler {
	return &FinanceNotificationHandler{
		name: "FinanceNotificationHandler",
	}
}

// Handle 处理事件
func (h *FinanceNotificationHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventDepositCompleted:
		data, _ := event.GetEventData().(DepositEventData)
		log.Printf("[FinanceNotification] 充值成功通知发送给用户 %s", data.UserID)

	case EventWithdrawalCompleted:
		data, _ := event.GetEventData().(WithdrawalEventData)
		log.Printf("[FinanceNotification] 提现完成通知发送给用户 %s", data.UserID)

	case EventSettlementGenerated:
		data, _ := event.GetEventData().(SettlementEventData)
		log.Printf("[FinanceNotification] 结算单生成通知发送给用户 %s", data.UserID)
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *FinanceNotificationHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *FinanceNotificationHandler) GetSupportedEventTypes() []string {
	return []string{
		EventDepositCompleted,
		EventWithdrawalCompleted,
		EventSettlementGenerated,
		EventSettlementPaid,
	}
}
