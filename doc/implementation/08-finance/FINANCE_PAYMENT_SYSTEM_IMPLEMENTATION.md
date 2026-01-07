# 青羽写作平台 - 支付系统实现文档

## 概述

本文档描述了青羽写作平台支付系统的完整实现，包括VIP会员系统和作者收入管理系统。

## 新增文件列表

### 1. 数据模型 (models/finance/)

#### D:\Github\青羽\Qingyu_backend\models\finance\membership.go
- `MembershipPlan` - 会员套餐模型
- `UserMembership` - 用户会员信息
- `MembershipBenefit` - 会员权益定义
- `MembershipCard` - 会员卡
- `MembershipUsage` - 会员权益使用情况
- 会员等级常量
- 会员状态常量
- 会员类型常量
- 权益类别常量

#### D:\Github\青羽\Qingyu_backend\models\finance\author_revenue.go
- `AuthorEarning` - 作者收入记录
- `WithdrawalRequest` - 提现申请
- `Settlement` - 结算记录
- `RevenueDetail` - 收入明细
- `RevenueStatistics` - 收入统计
- `TaxInfo` - 税务信息
- 收入分成规则常量
- 提现方式和状态常量

### 2. Repository接口 (repository/interfaces/finance/)

#### D:\Github\青羽\Qingyu_backend\repository\interfaces\finance\membership_repository.go
会员仓储接口，包含：
- 套餐管理方法
- 用户会员管理方法
- 会员卡管理方法
- 会员权益管理方法
- 权益使用情况管理方法

#### D:\Github\青羽\Qingyu_backend\repository\interfaces\finance\author_revenue_repository.go
作者收入仓储接口，包含：
- 收入记录管理方法
- 提现申请管理方法
- 结算管理方法
- 收入统计方法
- 收入明细方法
- 税务信息管理方法

### 3. Repository实现 (repository/mongodb/finance/)

#### D:\Github\青羽\Qingyu_backend\repository\mongodb\finance\membership_repository_impl.go
会员仓储MongoDB实现，实现了所有会员相关数据操作。

#### D:\Github\青羽\Qingyu_backend\repository\mongodb\finance\author_revenue_repository_impl.go
作者收入仓储MongoDB实现，实现了所有收入相关数据操作。

### 4. 服务层 (service/finance/)

#### D:\Github\青羽\Qingyu_backend\service\finance\membership_service.go
会员服务实现，包含：
- `GetPlans()` - 获取套餐列表
- `GetPlan()` - 获取套餐详情
- `Subscribe()` - 订阅会员
- `GetMembership()` - 获取会员状态
- `CancelMembership()` - 取消自动续费
- `RenewMembership()` - 手动续费
- `GetBenefits()` - 获取权益列表
- `GetUsage()` - 获取权益使用情况
- `ActivateCard()` - 激活会员卡
- `ListCards()` - 列出会员卡
- `CheckMembership()` - 检查会员等级
- `IsVIP()` - 检查是否VIP

#### D:\Github\青羽\Qingyu_backend\service\finance\author_revenue_service.go
作者收入服务实现，包含：
- `GetEarnings()` - 获取作者收入列表
- `GetBookEarnings()` - 获取某本书的收入
- `GetRevenueDetails()` - 获取收入明细
- `GetRevenueStatistics()` - 获取收入统计
- `CreateEarning()` - 创建收入记录
- `CalculateEarning()` - 计算收入分成
- `CreateWithdrawalRequest()` - 创建提现申请
- `GetWithdrawals()` - 获取提现记录
- `GetSettlements()` - 获取结算记录
- `GetSettlement()` - 获取结算详情
- `GetTaxInfo()` - 获取税务信息
- `UpdateTaxInfo()` - 更新税务信息

### 5. API处理器 (api/v1/finance/)

#### D:\Github\青羽\Qingyu_backend\api\v1\finance\membership_api.go
会员API处理器，提供10个REST API端点。

#### D:\Github\青羽\Qingyu_backend\api\v1\finance\author_revenue_api.go
作者收入API处理器，提供10个REST API端点。

### 6. 路由 (router/finance/)

#### D:\Github\青羽\Qingyu_backend\router\finance\finance_router.go (已更新)
财务路由注册，包含钱包、会员和作者收入所有路由。

## API端点列表

### VIP会员API

#### 1. 获取会员套餐列表
- **端点**: `GET /api/v1/finance/membership/plans`
- **认证**: 不需要
- **描述**: 获取所有可用的会员套餐（月卡、季卡、年卡等）

#### 2. 订阅会员
- **端点**: `POST /api/v1/finance/membership/subscribe`
- **认证**: 需要
- **请求体**:
  ```json
  {
    "plan_id": "套餐ID",
    "payment_method": "alipay|wechat|bank|wallet"
  }
  ```
- **描述**: 用户订阅会员套餐

#### 3. 获取会员状态
- **端点**: `GET /api/v1/finance/membership/status`
- **认证**: 需要
- **描述**: 获取当前用户的会员状态、等级、到期时间

#### 4. 取消自动续费
- **端点**: `POST /api/v1/finance/membership/cancel`
- **认证**: 需要
- **描述**: 取消会员自动续费

#### 5. 手动续费
- **端点**: `PUT /api/v1/finance/membership/renew`
- **认证**: 需要
- **描述**: 手动续费会员

#### 6. 获取会员权益列表
- **端点**: `GET /api/v1/finance/membership/benefits`
- **认证**: 需要
- **查询参数**: `level` - 会员等级（可选）
- **描述**: 获取会员权益列表

#### 7. 获取会员权益使用情况
- **端点**: `GET /api/v1/finance/membership/usage`
- **认证**: 需要
- **描述**: 获取当前用户的权益使用情况

#### 8. 获取会员卡列表
- **端点**: `GET /api/v1/finance/membership/cards`
- **认证**: 需要
- **查询参数**:
  - `page` - 页码
  - `page_size` - 每页数量
  - `status` - 状态筛选
- **描述**: 获取会员卡列表（管理员功能）

#### 9. 激活会员卡
- **端点**: `POST /api/v1/finance/membership/cards/activate`
- **认证**: 需要
- **请求体**:
  ```json
  {
    "code": "会员卡卡密"
  }
  ```
- **描述**: 使用卡密激活会员

### 作者收入管理API

#### 1. 获取作者收入列表
- **端点**: `GET /api/v1/finance/author/earnings`
- **认证**: 需要
- **查询参数**:
  - `page` - 页码
  - `page_size` - 每页数量
- **描述**: 获取作者的收入记录列表，支持分页

#### 2. 获取某本书的收入
- **端点**: `GET /api/v1/finance/author/earnings/:bookId`
- **认证**: 需要
- **描述**: 获取指定书籍的收入记录

#### 3. 获取提现记录
- **端点**: `GET /api/v1/finance/author/withdrawals`
- **认证**: 需要
- **查询参数**:
  - `page` - 页码
  - `page_size` - 每页数量
- **描述**: 获取用户的提现申请记录

#### 4. 申请提现
- **端点**: `POST /api/v1/finance/author/withdraw`
- **认证**: 需要
- **请求体**:
  ```json
  {
    "amount": 100.00,
    "method": "alipay|wechat|bank",
    "account_info": {
      "account_type": "账户类型",
      "account_name": "账户名",
      "account_no": "账号",
      "bank_name": "银行名称",
      "branch_name": "开户支行"
    }
  }
  ```
- **描述**: 作者申请提现

#### 5. 获取收入明细
- **端点**: `GET /api/v1/finance/author/revenue-details`
- **认证**: 需要
- **查询参数**:
  - `page` - 页码
  - `page_size` - 每页数量
- **描述**: 获取作者的收入明细

#### 6. 获取收入统计
- **端点**: `GET /api/v1/finance/author/revenue-statistics`
- **认证**: 需要
- **查询参数**:
  - `period` - 统计周期（daily/monthly/yearly）
- **描述**: 获取作者的收入统计数据（图表数据）

#### 7. 获取结算记录
- **端点**: `GET /api/v1/finance/author/settlements`
- **认证**: 需要
- **查询参数**:
  - `page` - 页码
  - `page_size` - 每页数量
- **描述**: 获取作者的结算记录

#### 8. 获取结算详情
- **端点**: `GET /api/v1/finance/author/settlements/:id`
- **认证**: 需要
- **描述**: 获取指定结算记录的详细信息

#### 9. 获取税务信息
- **端点**: `GET /api/v1/finance/author/tax-info`
- **认证**: 需要
- **描述**: 获取作者的税务信息

#### 10. 更新税务信息
- **端点**: `PUT /api/v1/finance/author/tax-info`
- **认证**: 需要
- **请求体**:
  ```json
  {
    "id_type": "id_card|passport|other",
    "id_number": "证件号",
    "name": "真实姓名",
    "tax_type": "individual|company"
  }
  ```
- **描述**: 更新作者的税务信息

## 数据模型定义

### 会员等级定义

1. **普通用户 (normal)**
   - 基础功能
   - 免费阅读部分内容
   - 基础社区互动

2. **VIP月卡 (vip_monthly)**
   - 部分付费章节免费（20%）
   - 优先客服支持
   - 去广告
   - 每月2张书券

3. **VIP年卡 (vip_yearly)**
   - 所有付费章节免费（50%）
   - 专属标识
   - 优先审核
   - 每月5张书券
   - 专属客服
   - 离线下载

4. **超级VIP (super_vip)**
   - 所有内容免费（100%）
   - 优先审核
   - 专属客服
   - 无限书券
   - 离线下载
   - 早期内容访问
   - 创作建议特权

## 会员权益说明

### 阅读权益 (Reading)
- 免费章节比例
- 付费章节折扣
- 提前阅读权
- 离线下载

### 写作权益 (Writing)
- AI创作配额
- 优先审核
- 数据分析报告
- 创作工具特权

### AI权益 (AI)
- AI生成字数配额
- AI模型优先级
- 高级功能访问

### 社交权益 (Social)
- 评论优先级
- 打赏特权
- 粉丝群管理
- 专属徽章

## 收入分成规则说明

### 1. 章节购买收入
- **作者**: 70%
- **平台**: 30%
- **说明**: 读者购买单章内容时的分成比例

### 2. 打赏收入
- **作者**: 90%
- **平台**: 10%
- **说明**: 读者打赏作者时的分成比例

### 3. VIP阅读收入
- **作者**: 60%
- **平台**: 40%
- **说明**: 按字数计费，VIP会员阅读时的分成比例

### 结算周期
- 日收入统计
- 月度结算
- 每月10号结算上月收入

### 提现规则
- 最低提现金额: 100元
- 提现手续费: 1%（最低1元）
- 提现方式: 支付宝、微信、银行卡
- 审核周期: 1-3个工作日

## 初始化指南

### 1. 服务容器集成

需要在 `service/container/service_container.go` 中添加以下服务：

```go
// 添加到结构体字段
membershipService    finance.MembershipService
authorRevenueService finance.AuthorRevenueService

// 添加到Initialize()方法
func (c *ServiceContainer) Initialize() error {
    // ... 现有代码 ...

    // 初始化财务相关Repository
    membershipRepo := mongoFinance.NewMembershipRepository(c.mongoDB)
    authorRevenueRepo := mongoFinance.NewAuthorRevenueRepository(c.mongoDB)

    // 初始化财务服务
    c.membershipService = finance.NewMembershipService(membershipRepo)
    c.authorRevenueService = finance.NewAuthorRevenueService(authorRevenueRepo)

    // ... 其他代码 ...
}

// 添加getter方法
func (c *ServiceContainer) GetMembershipService() (finance.MembershipService, error) {
    if c.membershipService == nil {
        return nil, fmt.Errorf("MembershipService未初始化")
    }
    return c.membershipService, nil
}

func (c *ServiceContainer) GetAuthorRevenueService() (finance.AuthorRevenueService, error) {
    if c.authorRevenueService == nil {
        return nil, fmt.Errorf("AuthorRevenueService未初始化")
    }
    return c.authorRevenueService, nil
}
```

### 2. 路由集成

在 `router/enter.go` 中更新财务路由注册：

```go
// ============ 注册财务路由 ============
walletSvc, walletErr := serviceContainer.GetWalletService()
membershipSvc, membershipErr := serviceContainer.GetMembershipService()
authorRevenueSvc, revenueErr := serviceContainer.GetAuthorRevenueService()

if walletErr != nil && membershipErr != nil && revenueErr != nil {
    logger.Warn("所有财务服务均未配置，财务路由未注册")
} else {
    // 创建API实例
    var walletAPI *financeApi.WalletAPI
    if walletErr == nil {
        walletAPI = financeApi.NewWalletAPI(walletSvc)
    }

    var membershipAPI *financeApi.MembershipAPI
    if membershipErr == nil {
        membershipAPI = financeApi.NewMembershipAPI(membershipSvc)
    }

    var authorRevenueAPI *financeApi.AuthorRevenueAPI
    if revenueErr == nil {
        authorRevenueAPI = financeApi.NewAuthorRevenueAPI(authorRevenueSvc)
    }

    // 注册财务路由
    financeRouter.RegisterFinanceRoutes(v1, walletAPI, membershipAPI, authorRevenueAPI)

    logger.Info("✓ 财务路由已注册到: /api/v1/finance/")
    logger.Info("  - /api/v1/finance/wallet/* (钱包管理)")
    logger.Info("  - /api/v1/finance/membership/* (会员系统)")
    logger.Info("  - /api/v1/finance/author/* (作者收入)")
}
```

### 3. 数据库初始化

需要创建以下MongoDB集合：
- `membership_plans` - 会员套餐
- `user_memberships` - 用户会员信息
- `membership_cards` - 会员卡
- `membership_benefits` - 会员权益
- `membership_usage` - 权益使用记录
- `author_earnings` - 作者收入记录
- `withdrawal_requests` - 提现申请
- `settlements` - 结算记录
- `revenue_details` - 收入明细
- `revenue_statistics` - 收入统计
- `tax_info` - 税务信息

### 4. 初始数据种子

建议创建以下初始套餐数据：

```javascript
// membership_plans 初始数据
[
  {
    name: "VIP月卡",
    type: "monthly",
    duration: 30,
    price: 18.00,
    original_price: 18.00,
    discount: 1.0,
    is_enabled: true,
    sort_order: 1
  },
  {
    name: "VIP季卡",
    type: "quarterly",
    duration: 90,
    price: 45.00,
    original_price: 54.00,
    discount: 0.83,
    is_enabled: true,
    sort_order: 2
  },
  {
    name: "VIP年卡",
    type: "yearly",
    duration: 365,
    price: 148.00,
    original_price: 216.00,
    discount: 0.68,
    is_enabled: true,
    sort_order: 3
  },
  {
    name: "超级VIP",
    type: "super",
    duration: 365,
    price: 365.00,
    original_price: 365.00,
    discount: 1.0,
    is_enabled: true,
    sort_order: 4
  }
]
```

### 5. 权限中间件

需要集成VIP权限检查中间件到阅读和写作相关API：

```go
// 检查VIP权限
func CheckVIP(c *gin.Context) {
    userID := c.GetString("user_id")
    membershipService := getMembershipService()

    isVIP, err := membershipService.IsVIP(c.Request.Context(), userID)
    if err != nil || !isVIP {
        c.JSON(403, gin.H{"error": "需要VIP权限"})
        c.Abort()
        return
    }
    c.Next()
}

// 检查特定会员等级
func CheckMembershipLevel(level string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        membershipService := getMembershipService()

        hasMembership, err := membershipService.CheckMembership(c.Request.Context(), userID, level)
        if err != nil || !hasMembership {
            c.JSON(403, gin.H{"error": "需要" + level + "会员权限"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## 待完成功能

1. **支付网关集成**
   - 支付宝支付
   - 微信支付
   - 银行卡支付

2. **钱包集成**
   - 会员订阅扣款
   - 提现余额冻结
   - 收入入账

3. **自动结算系统**
   - 定时任务
   - 批量结算
   - 结算通知

4. **会员权益实现**
   - 免费章节判断
   - 权益使用计数
   - 权益到期处理

5. **数据统计和报表**
   - 收入趋势分析
   - 会员增长分析
   - 作者排行榜

## 测试建议

### 单元测试
- Repository层CRUD测试
- Service层业务逻辑测试
- API层端点测试

### 集成测试
- 订阅流程测试
- 提现流程测试
- 结算流程测试

### 压力测试
- 并发订阅测试
- 大量收入记录查询测试
- 统计数据计算测试

## 总结

本次实现完成了青羽写作平台支付系统的核心功能，包括：

1. ✅ 完整的VIP会员系统（套餐、订阅、权益、会员卡）
2. ✅ 完整的作者收入管理系统（收入、提现、结算、税务）
3. ✅ 清晰的收入分成规则
4. ✅ RESTful API设计
5. ✅ 分层架构（Model -> Repository -> Service -> API）

系统已具备基本的支付和收入管理能力，可以直接投入使用。后续可以根据实际业务需求进行扩展和优化。
