# 08 - 财务模块

> **模块编号**: 08
> **模块名称**: Finance & Payment
> **负责功能**: 充值、提现、收入结算、会员系统
> **完成度**: 🟡 60%

## 📋 目录结构

```
财务模块/
├── api/v1/
│   └── finance/                  # 财务API
│       ├── wallet_api.go        # 钱包管理
│       ├── payment_api.go       # 支付管理
│       ├── membership_api.go    # 会员管理
│       ├── revenue_api.go       # 收入管理
│       └── order_api.go         # 订单管理
├── service/finance/              # 财务服务层
│   ├── wallet_service.go       # 钱包服务
│   ├── payment_service.go      # 支付服务
│   ├── membership_service.go   # 会员服务
│   └── revenue_service.go      # 收入服务
├── repository/interfaces/finance/ # 仓储接口
├── repository/mongodb/finance/    # MongoDB仓储实现
└── models/finance/                # 数据模型
    ├── wallet.go                 # 钱包
    ├── transaction.go           # 交易记录
    ├── membership.go            # 会员
    ├── order.go                 # 订单
    └── settlement.go            # 结算
```

## 🎯 核心功能

### 1. 钱包管理

- **余额查询**: 查询钱包余额
- **充值**: 充值到钱包
- **提现**: 从钱包提现
- **交易记录**: 查询交易明细
- **账单管理**: 月度账单

### 2. 支付管理

- **支付方式**: 微信、支付宝、银行卡
- **支付回调**: 支付结果通知
- **退款处理**: 订单退款
- **对账系统**: 财务对账

### 3. 会员系统

- **会员套餐**: 月卡、季卡、年卡
- **会员订阅**: 购买会员
- **会员权益**: 免费阅读、专属标识
- **自动续费**: 自动订阅续费
- **会员卡**: 激活会员卡

### 4. 作者收入

- **收入统计**: 查看收入明细
- **分成结算**: 按规则分成
- **提现申请**: 申请提现
- **收入报表**: 收入报表

### 5. 订单管理

- **创建订单**: 创建购买订单
- **订单查询**: 查询订单状态
- **订单取消**: 取消未支付订单
- **订单退款**: 订单退款处理

## 📊 数据模型

### Wallet (钱包)

```go
type Wallet struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    Balance         int                  `bson:"balance" json:"balance"`           // 余额（分）
    FrozenBalance   int                  `bson:"frozen_balance" json:"frozenBalance"` // 冻结余额
    TotalIncome     int64                `bson:"total_income" json:"totalIncome"`
    TotalExpense    int64                `bson:"total_expense" json:"totalExpense"`
    UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
}
```

### Transaction (交易记录)

```go
type Transaction struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    Type            TransactionType      `bson:"type" json:"type"`
    Amount          int                  `bson:"amount" json:"amount"`             // 金额（分）
    BalanceBefore   int                  `bson:"balance_before" json:"balanceBefore"`
    BalanceAfter    int                  `bson:"balance_after" json:"balanceAfter"`
    Description     string               `bson:"description" json:"description"`
    OrderID         *string              `bson:"order_id,omitempty" json:"orderId,omitempty"`
    Status          TransactionStatus    `bson:"status" json:"status"`
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
}

type TransactionType string
const (
    TransactionTypeRecharge    TransactionType = "recharge"
    TransactionTypeConsume     TransactionType = "consume"
    TransactionTypeWithdraw    TransactionType = "withdraw"
    TransactionTypeRefund      TransactionType = "refund"
    TransactionTypeIncome      TransactionType = "income"
)

type TransactionStatus string
const (
    TransactionStatusPending   TransactionStatus = "pending"
    TransactionStatusSuccess   TransactionStatus = "success"
    TransactionStatusFailed    TransactionStatus = "failed"
)
```

### Membership (会员)

```go
type Membership struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    Level           MembershipLevel      `bson:"level" json:"level"`
    StartTime       time.Time            `bson:"start_time" json:"startTime"`
    EndTime         time.Time            `bson:"end_time" json:"endTime"`
    AutoRenew       bool                 `bson:"auto_renew" json:"autoRenew"`
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
}

type MembershipLevel string
const (
    MembershipLevelFree    MembershipLevel = "free"
    MembershipLevelMonth   MembershipLevel = "month"
    MembershipLevelQuarter MembershipLevel = "quarter"
    MembershipLevelYear    MembershipLevel = "year"
)
```

### Order (订单)

```go
type Order struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    OrderNo         string               `bson:"order_no" json:"orderNo"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    Type            OrderType            `bson:"type" json:"type"`
    Amount          int                  `bson:"amount" json:"amount"`
    PaymentMethod   string               `bson:"payment_method" json:"paymentMethod"`
    Status          OrderStatus          `bson:"status" json:"status"`
    PaidAt          *time.Time           `bson:"paid_at,omitempty" json:"paidAt,omitempty"`
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
    UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
}

type OrderType string
const (
    OrderTypeRecharge    OrderType = "recharge"
    OrderTypeMembership  OrderType = "membership"
    OrderTypeChapter     OrderType = "chapter"
    OrderTypeBook        OrderType = "book"
)

type OrderStatus string
const (
    OrderStatusPending   OrderStatus = "pending"
    OrderStatusPaid      OrderStatus = "paid"
    OrderStatusCancelled OrderStatus = "cancelled"
    OrderStatusRefunded  OrderStatus = "refunded"
)
```

### Settlement (结算)

```go
type Settlement struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    AuthorID        primitive.ObjectID   `bson:"author_id" json:"authorId"`
    Period          string               `bson:"period" json:"period"`             // 结算周期 YYYY-MM
    GrossIncome     int                  `bson:"gross_income" json:"grossIncome"`   // 总收入
    PlatformFee     int                  `bson:"platform_fee" json:"platformFee"`   // 平台费用
    NetIncome       int                  `bson:"net_income" json:"netIncome"`       // 净收入
    Status          SettlementStatus      `bson:"status" json:"status"`
    SettledAt       *time.Time           `bson:"settled_at,omitempty" json:"settledAt,omitempty"`
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
}

type SettlementStatus string
const (
    SettlementStatusPending   SettlementStatus = "pending"
    SettlementStatusSettled   SettlementStatus = "settled"
)
```

## 🌐 API端点

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/finance/wallet | 获取钱包信息 | 是 |
| GET | /api/v1/finance/transactions | 获取交易记录 | 是 |
| POST | /api/v1/finance/recharge | 发起充值 | 是 |
| POST | /api/v1/finance/withdraw | 申请提现 | 是 |
| GET | /api/v1/finance/membership | 获取会员信息 | 是 |
| POST | /api/v1/finance/membership/subscribe | 购买会员 | 是 |
| PUT | /api/v1/finance/membership/auto-renew | 设置自动续费 | 是 |
| GET | /api/v1/finance/revenue | 获取收入信息 | 是 |
| GET | /api/v1/finance/revenue/settlements | 获取结算记录 | 是 |
| GET | /api/v1/finance/orders | 获取订单列表 | 是 |
| GET | /api/v1/finance/orders/:id | 获取订单详情 | 是 |
| POST | /api/v1/finance/orders/:id/cancel | 取消订单 | 是 |
| POST | /api/v1/finance/payment/callback | 支付回调 | 否 |

## 🔐 安全考虑

### 支付安全

- 签名验证
- 金额校验
- 重复支付检测
- 支付密码

### 提现安全

- 实名认证
- 提现限额
- 审核流程
- 防刷机制

## 🔧 依赖关系

### 依赖的模块
- **01 - 认证授权**: 用户身份验证
- **02 - 写作创作**: 作者作品收入计算
- **06 - 书城**: 章节购买

### 外部服务
- **支付网关**: 微信支付、支付宝
- **银行接口**: 银行卡提现

## 📈 扩展点

1. **优惠券系统**
   - 优惠券发放
   - 优惠券使用
   - 满减活动

2. **积分系统**
   - 积分获取
   - 积分兑换
   - 积分商城

3. **财务报表**
   - 收支报表
   - 流水报表
   - 对账报表

---

**文档维护**: 青羽后端架构团队
**最后更新**: 2025-01-06
**对应实现**: `../../Qingyu_backend/api/v1/finance/`
