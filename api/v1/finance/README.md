# Finance API 模块结构说明

> finance 模块是后端唯一的财务域 owner，负责钱包管理、会员系统、作者收入等全部财务相关功能。作者侧（writer）等模块仅消费展示，不持有财务逻辑。

## 文件结构

```
api/v1/finance/
├── wallet_api.go           # 钱包管理API（余额、充值、消费、转账、提现）
├── membership_api.go       # 会员系统API（套餐、订阅、续费、权益、卡密）
├── author_revenue_api.go   # 作者收入API（收入查询、提现、结算、税务）
├── types.go                # 公共类型定义
└── README.md               # 本文件
```

## 模块职责划分

### 1. WalletAPI (`wallet_api.go`)

**职责**: 用户钱包管理，包括余额查询、充值、消费、转账和提现等核心资金操作

**核心功能**:
- ✅ 查询钱包余额
- ✅ 获取钱包完整信息（兼容 `/wallet` 和 `/wallet/detail` 两条路径）
- ✅ 钱包充值（支持支付宝、微信、银行卡）
- ✅ 钱包消费（余额扣减）
- ✅ 用户间转账
- ✅ 交易记录查询（分页、按类型筛选）
- ✅ 申请提现
- ✅ 提现申请列表查询（兼容 `/withdraws` 和 `/withdrawals`）

**API端点**:
```
GET  /api/v1/finance/wallet                # 获取钱包信息
GET  /api/v1/finance/wallet/balance        # 查询钱包余额
GET  /api/v1/finance/wallet/detail         # 获取钱包详情（与 /wallet 等价）
POST /api/v1/finance/wallet/recharge       # 钱包充值
POST /api/v1/finance/wallet/consume        # 钱包消费
POST /api/v1/finance/wallet/transfer       # 用户转账
GET  /api/v1/finance/wallet/transactions   # 获取交易记录
POST /api/v1/finance/wallet/withdraw       # 申请提现
GET  /api/v1/finance/wallet/withdraws      # 提现申请列表
GET  /api/v1/finance/wallet/withdrawals    # 提现申请列表（兼容路径）
```

**依赖服务**:
- `wallet.WalletService` - 钱包核心服务（余额管理、交易、提现）

---

### 2. MembershipAPI (`membership_api.go`)

**职责**: 会员系统管理，包括套餐展示、订阅/续费/取消、权益查询和卡密激活

**核心功能**:
- ✅ 获取会员套餐列表（公开接口，无需认证）
- ✅ 订阅会员（支持支付宝、微信、银行卡、钱包支付）
- ✅ 获取会员状态
- ✅ 取消自动续费
- ✅ 手动续费
- ✅ 获取会员权益列表
- ✅ 获取权益使用情况
- ✅ 获取会员卡列表
- ✅ 激活会员卡（卡密兑换）

**API端点**:
```
GET  /api/v1/finance/membership/plans             # 获取会员套餐列表
GET  /api/v1/finance/membership/status            # 获取会员状态
POST /api/v1/finance/membership/subscribe         # 订阅会员
POST /api/v1/finance/membership/cancel            # 取消自动续费
PUT  /api/v1/finance/membership/renew             # 手动续费
GET  /api/v1/finance/membership/benefits          # 获取会员权益列表
GET  /api/v1/finance/membership/usage             # 获取权益使用情况
GET  /api/v1/finance/membership/cards             # 获取会员卡列表
POST /api/v1/finance/membership/cards/activate    # 激活会员卡
```

**依赖服务**:
- `finance.MembershipService` - 会员核心服务（套餐管理、订阅、权益）

---

### 3. AuthorRevenueAPI (`author_revenue_api.go`)

**职责**: 作者收入管理，包括收入记录查询、收入统计、提现申请、结算管理和税务信息维护

**核心功能**:
- ✅ 获取作者收入列表（分页）
- ✅ 获取指定书籍的收入记录
- ✅ 获取收入明细
- ✅ 获取收入统计（支持按日/月/年统计）
- ✅ 获取提现记录（分页）
- ✅ 申请提现（支持支付宝、微信、银行卡）
- ✅ 获取结算记录列表
- ✅ 获取结算详情
- ✅ 获取税务信息
- ✅ 更新税务信息（身份证/护照、个人/企业）

**API端点**:
```
# 收入查询
GET  /api/v1/finance/author/earnings               # 获取作者收入列表
GET  /api/v1/finance/author/earnings/:bookId       # 获取指定书籍收入
GET  /api/v1/finance/author/revenue-details        # 获取收入明细
GET  /api/v1/finance/author/revenue-statistics     # 获取收入统计

# 提现管理
GET  /api/v1/finance/author/withdrawals            # 获取提现记录
POST /api/v1/finance/author/withdraw               # 申请提现

# 结算管理
GET  /api/v1/finance/author/settlements            # 获取结算记录列表
GET  /api/v1/finance/author/settlements/:id        # 获取结算详情

# 税务信息
GET  /api/v1/finance/author/tax-info               # 获取税务信息
PUT  /api/v1/finance/author/tax-info               # 更新税务信息
```

**依赖服务**:
- `finance.AuthorRevenueService` - 作者收入核心服务（收入查询、提现、结算、税务）

---

## 中间件配置

所有 finance 路由统一挂载在 `/api/v1/finance` 下，并应用以下中间件：

| 中间件 | 说明 |
|--------|------|
| `auth.JWTAuth()` | JWT 认证，所有 finance 接口均需登录（plans 除外） |
| `ratelimit.RateLimitMiddlewareSimple(50, 60)` | 限流，每分钟 50 次请求 |

---

## 数据模型映射

`author_revenue_api.go` 内部使用 `mapXxx()` 系列函数将 Model 层对象转换为 API 响应格式：

| 映射函数 | 源模型 | 输出 |
|----------|--------|------|
| `mapAuthorEarning` | `financeModel.AuthorEarning` | 收入记录（含金额元/分双格式） |
| `mapAuthorWithdrawal` | `financeModel.WithdrawalRequest` | 提现记录（含账户信息嵌套） |
| `mapRevenueDetail` | `financeModel.RevenueDetail` | 收入明细汇总 |
| `mapRevenueStatistics` | `financeModel.RevenueStatistics` | 收入统计（按周期） |
| `mapSettlement` | `financeModel.Settlement` | 结算记录（含平台费、税费） |

---

## 金额处理规范

- API 层接收金额单位为 **元**（float64），内部自动转换为 **分**（int64）传递给服务层
- 响应中同时返回元和分两种格式（如 `amount` 和 `amount_cents`），前端按需使用
- 金额校验使用自定义验证器 `positive_amount`、`amount_range`

---

## 路由注册入口

路由在 `router/finance/finance_router.go` 的 `RegisterFinanceRoutes` 函数中统一注册，三个 API 处理器通过参数注入，支持 nil 安全（未启用的模块跳过注册）。

---

**版本**: v1.0
**更新日期**: 2026-05-21
**维护者**: yukin371
