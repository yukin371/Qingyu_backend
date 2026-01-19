# 共享底层服务设计文档

> 基于模块化单体架构的共享服务设计
> 
> **架构模式**: 模块化单体（Modular Monolith）  
> **创建时间**: 2025-09-30  
> **最后更新**: 2025-10-17

---

## 📋 概述

共享底层服务是青羽轻量级阅读写作平台的基础设施层，为阅读端和写作端提供统一的账号管理、钱包系统、推荐服务等核心功能。

### 架构选择说明

**当前阶段采用模块化单体架构**，原因如下：

- ✅ 团队规模较小（<10人），单体架构开发效率更高
- ✅ 业务复杂度适中，暂不需要微服务的复杂性
- ✅ 初期流量不大，单体应用完全可以支撑
- ✅ 统一部署和事务管理，降低运维成本
- ✅ 通过模块化设计，为未来可能的微服务化预留空间

> 📖 详见：[微服务架构划分建议](../微服务架构划分建议.md)

---

## 🏗️ 服务架构

### 模块结构

```
service/shared/              # 共享服务模块
├── auth/                    # 账号与权限模块
│   ├── service.go          # 服务实现
│   ├── repository.go       # 数据访问
│   └── interfaces.go       # 模块接口（对外暴露）
│
├── wallet/                  # 钱包系统模块
│   ├── service.go
│   ├── repository.go
│   └── interfaces.go
│
├── recommendation/          # 推荐服务模块
│   ├── service.go
│   ├── repository.go
│   └── interfaces.go
│
├── messaging/               # 消息队列模块
│   ├── service.go
│   ├── repository.go
│   └── interfaces.go
│
├── storage/                 # 文件存储模块
│   ├── service.go
│   ├── repository.go
│   └── interfaces.go
│
└── admin/                   # 管理后台模块
    ├── service.go
    ├── repository.go
    └── interfaces.go
```

### 模块边界原则

1. **清晰的接口定义** - 每个模块通过 `interfaces.go` 暴露公开接口
2. **禁止跨模块直接调用** - 模块间只能通过接口交互
3. **独立的数据访问** - 每个模块有独立的 Repository 层
4. **共享代码最小化** - 共享代码统一放在 `pkg/` 或 `shared/utils/`

---

## 🎯 核心服务模块

### 1. 账号与权限模块 (auth)

#### 功能职责

- ✅ **用户注册**: 手机号 + 邮箱双通道注册
- ✅ **身份认证**: JWT Token + 刷新机制
- ✅ **角色管理**: Reader/Author/Admin 三角色体系
- ✅ **权限控制**: 基于 RBAC 的细粒度权限管理
- ✅ **会话管理**: Redis 存储会话状态

#### 技术实现

```go
// service/shared/auth/interfaces.go
package auth

type AuthService interface {
    // 用户注册
    Register(ctx context.Context, req *RegisterRequest) (*User, error)
    
    // 用户登录
    Login(ctx context.Context, username, password string) (token string, err error)
    
    // 验证Token
    ValidateToken(ctx context.Context, token string) (*UserClaims, error)
    
    // 权限检查
    CheckPermission(ctx context.Context, userID, permission string) (bool, error)
}
```

#### 数据模型

- `users` 集合：用户基本信息
- `roles` 集合：角色定义
- `permissions` 集合：权限定义
- Redis：Token黑名单、会话信息

---

### 2. 钱包系统模块 (wallet)

#### 功能职责

- ✅ **代币管理**: 充值、消费、余额查询
- ✅ **订单系统**: 雪花算法生成唯一订单号
- ✅ **支付集成**: 支持多种支付方式（支付宝、微信）
- ✅ **提现服务**: 支付宝提现功能
- ✅ **交易记录**: 完整的资金流水记录
- ✅ **风控系统**: 异常交易检测和风险控制

#### 技术实现

```go
// service/shared/wallet/interfaces.go
package wallet

type WalletService interface {
    // 创建钱包
    CreateWallet(ctx context.Context, userID string) (*Wallet, error)
    
    // 充值
    Recharge(ctx context.Context, userID string, amount float64) (*Transaction, error)
    
    // 消费
    Consume(ctx context.Context, userID string, amount float64, reason string) error
    
    // 提现
    Withdraw(ctx context.Context, userID string, amount float64) (*WithdrawRequest, error)
    
    // 查询余额
    GetBalance(ctx context.Context, userID string) (float64, error)
}
```

#### 数据模型

- `wallets` 集合：钱包信息
- `transactions` 集合：交易记录
- `withdraw_requests` 集合：提现申请

---

### 3. 推荐服务模块 (recommendation)

#### 功能职责

- ✅ **协同过滤**: 基于用户行为的推荐算法
- ✅ **内容标签**: 基于内容特征的推荐
- ✅ **热度推荐**: 基于统计的热门内容
- ✅ **个性化推荐**: 用户画像驱动的推荐
- ✅ **相似推荐**: 内容相似度计算

#### 技术实现

```go
// service/shared/recommendation/interfaces.go
package recommendation

type RecommendationService interface {
    // 获取个性化推荐
    GetPersonalizedRecommendations(ctx context.Context, userID string, limit int) ([]*RecommendedItem, error)
    
    // 获取相似内容
    GetSimilarItems(ctx context.Context, itemID string, limit int) ([]*RecommendedItem, error)
    
    // 记录用户行为（用于推荐模型训练）
    RecordUserBehavior(ctx context.Context, userID, itemID string, action string) error
}
```

#### 实现策略

- **简单算法优先**: 初期使用热度排序、标签匹配等简单算法
- **离线计算**: 定时任务预计算推荐结果，存储到 Redis
- **在线服务**: 从 Redis 快速读取预计算结果
- **逐步演进**: 后期根据数据积累引入机器学习算法

---

### 4. 消息队列模块 (messaging)

#### 功能职责

- ✅ **异步任务**: 后台任务处理（邮件、短信、数据统计）
- ✅ **事件发布订阅**: 模块间解耦通信
- ✅ **数据埋点**: 用户行为数据收集
- ✅ **延迟任务**: 定时任务调度

#### 技术实现（简化版）

**初期方案：使用 Go Channel + Redis**

```go
// service/shared/messaging/interfaces.go
package messaging

type MessagingService interface {
    // 发布消息
    Publish(ctx context.Context, topic string, message interface{}) error
    
    // 订阅消息
    Subscribe(ctx context.Context, topic string, handler MessageHandler) error
    
    // 发送延迟消息
    PublishDelayed(ctx context.Context, topic string, message interface{}, delay time.Duration) error
}
```

**技术选型**：
- **开发环境**: Go Channel（内存队列）
- **生产环境**: Redis Streams（轻量级消息队列）
- **未来升级**: 流量增长后考虑 Kafka/RabbitMQ

---

### 5. 文件存储模块 (storage)

#### 功能职责

- ✅ **文件上传**: 支持多种文件类型
- ✅ **文件下载**: 带权限控制的下载
- ✅ **文件管理**: 删除、重命名、版本管理
- ✅ **安全控制**: 访问权限、防盗链
- ✅ **存储优化**: 文件去重、压缩

#### 技术实现

```go
// service/shared/storage/interfaces.go
package storage

type StorageService interface {
    // 上传文件
    Upload(ctx context.Context, file *FileUpload) (*FileInfo, error)
    
    // 下载文件（返回URL或文件流）
    Download(ctx context.Context, fileID string) (string, error)
    
    // 删除文件
    Delete(ctx context.Context, fileID string) error
    
    // 获取文件信息
    GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error)
}
```

**存储策略**：
- **本地存储**: 开发环境使用本地文件系统
- **对象存储**: 生产环境使用云存储（OSS/S3）
- **CDN加速**: 静态资源通过 CDN 分发

---

### 6. 管理后台模块 (admin)

#### 功能职责

- ✅ **内容审核**: 书籍上线/驳回/下架管理
- ✅ **榜单管理**: 人工干预和操作日志
- ✅ **用户管理**: 用户信息和权限管理
- ✅ **财务管理**: 提现审核和风控
- ✅ **系统监控**: 服务状态和性能监控
- ✅ **操作日志**: 管理员操作审计

#### 技术实现

```go
// service/shared/admin/interfaces.go
package admin

type AdminService interface {
    // 内容审核
    ReviewContent(ctx context.Context, contentID string, action ReviewAction) error
    
    // 用户管理
    ManageUser(ctx context.Context, userID string, action UserAction) error
    
    // 提现审核
    ReviewWithdrawal(ctx context.Context, withdrawID string, approved bool) error
    
    // 记录操作日志
    LogOperation(ctx context.Context, adminID, operation string, details interface{}) error
}
```

---

## 🛠️ 技术架构（模块化单体）

### 应用层架构

```
┌─────────────────────────────────────────┐
│         Gin HTTP Server (单一入口)        │
├─────────────────────────────────────────┤
│              Middleware                  │
│  (Auth, Logging, CORS, RateLimit, etc.) │
├─────────────────────────────────────────┤
│               API Layer                  │
│        (RESTful API Handlers)            │
├─────────────────────────────────────────┤
│             Service Layer                │
│  ┌──────┬──────┬──────┬──────┬──────┐  │
│  │ Auth │Wallet│Recom.│Msg.  │Store │  │
│  └──────┴──────┴──────┴──────┴──────┘  │
├─────────────────────────────────────────┤
│           Repository Layer               │
│         (Data Access Objects)            │
├─────────────────────────────────────────┤
│             Data Layer                   │
│    MongoDB  │  Redis  │  FileSystem     │
└─────────────────────────────────────────┘
```

### 数据存储

#### MongoDB (主数据库)

```
数据库: qingyu_db
集合:
├── users                    # 用户信息
├── roles                    # 角色定义
├── permissions              # 权限定义
├── wallets                  # 钱包信息
├── transactions             # 交易记录
├── withdraw_requests        # 提现申请
├── user_behaviors           # 用户行为数据
├── files                    # 文件元数据
└── admin_logs              # 管理员操作日志
```

**优化策略**：
- 合理的索引设计
- 读写分离（主从复制）
- 数据归档（历史数据定期归档）

#### Redis (缓存 + 会话)

```
用途:
├── 会话管理        key: session:{token}
├── Token黑名单     key: blacklist:{token}
├── 缓存用户信息    key: user:{id}
├── 推荐结果缓存    key: recommend:{user_id}
├── 限流计数器      key: ratelimit:{ip}:{endpoint}
└── 消息队列        Redis Streams
```

#### 本地文件存储 (开发环境)

```
uploads/
├── avatars/        # 用户头像
├── books/          # 书籍文件
├── covers/         # 封面图片
└── attachments/    # 附件
```

**生产环境**: 迁移到云存储（阿里云OSS/AWS S3）

---

### 核心技术栈

| 组件 | 技术选型 | 说明 |
|------|---------|------|
| **Web框架** | Gin | 高性能 HTTP 框架 |
| **数据库** | MongoDB | 文档型数据库，灵活的 Schema |
| **缓存** | Redis | 内存缓存 + 会话存储 |
| **认证** | JWT | 无状态认证 |
| **密码加密** | bcrypt | 密码哈希 |
| **配置管理** | Viper | 配置文件管理（.env / YAML） |
| **日志** | Zap | 高性能结构化日志 |
| **参数验证** | validator | 请求参数验证 |
| **定时任务** | cron | 定时任务调度 |
| **消息队列** | Redis Streams | 轻量级消息队列 |

---

## 🔌 服务接口

### 内部模块通信

**原则**: 模块间只能通过接口调用，禁止直接访问

```go
// 示例：推荐服务调用用户服务
package recommendation

type RecommendationServiceImpl struct {
    userService  auth.AuthService      // 依赖注入
    bookService  bookstore.BookService // 依赖注入
}

func (s *RecommendationServiceImpl) GetPersonalizedRecommendations(
    ctx context.Context, 
    userID string, 
    limit int,
) ([]*RecommendedItem, error) {
    // 1. 通过接口获取用户信息
    user, err := s.userService.GetUserByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 2. 基于用户画像推荐
    // ...
    
    return recommendations, nil
}
```

### 外部 API 接口

**RESTful API**: 标准的 HTTP JSON API

```
GET    /api/v1/auth/profile           # 获取用户信息
POST   /api/v1/auth/login             # 用户登录
POST   /api/v1/wallet/recharge        # 钱包充值
GET    /api/v1/recommendations        # 获取推荐
POST   /api/v1/storage/upload         # 文件上传
```

---

## 📦 部署架构

### 开发环境

**单机部署 + Docker Compose**

```yaml
# docker-compose.dev.yml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - mongodb
      - redis
    
  mongodb:
    image: mongo:5.0
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    
  redis:
    image: redis:6.2
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  mongodb_data:
  redis_data:
```

**启动命令**:
```bash
docker-compose -f docker-compose.dev.yml up -d
go run main.go
```

---

### 生产环境（初期）

**单体应用 + 基础设施**

```
┌────────────────────────────────────┐
│         Nginx (反向代理)            │
│    - SSL终止                        │
│    - 静态资源                       │
│    - 负载均衡（可选）               │
└────────────────────────────────────┘
              ↓
┌────────────────────────────────────┐
│      Qingyu Backend (单体应用)     │
│         多实例部署（可选）          │
└────────────────────────────────────┘
              ↓
┌──────────┬──────────┬──────────────┐
│ MongoDB  │  Redis   │  File Storage│
│ (主从)   │ (哨兵)   │   (OSS)      │
└──────────┴──────────┴──────────────┘
```

**扩展策略**:
1. **垂直扩展**: 增加服务器配置（CPU、内存）
2. **水平扩展**: 多实例部署 + Nginx 负载均衡
3. **数据库扩展**: MongoDB 主从复制、分片（流量增长后）

---

### 未来微服务化准备

当满足以下条件时，考虑拆分微服务：

- ✅ 日活用户 > 10万
- ✅ 团队规模 > 15人
- ✅ 某个模块成为性能瓶颈

**拆分优先级**：

1. **推荐服务** - 计算密集型，易独立
2. **文件存储** - 流量独立，易拆分
3. **钱包系统** - 业务边界清晰
4. **消息队列** - 基础设施服务

> 📖 详见：[微服务架构划分建议](../微服务架构划分建议.md)

---

## 🔒 安全设计

### 数据安全

- ✅ **密码加密**: bcrypt 哈希存储
- ✅ **敏感数据**: 数据库字段加密（身份证、银行卡）
- ✅ **传输加密**: HTTPS (TLS 1.2+)
- ✅ **API签名**: 关键接口使用签名验证

### 访问控制

- ✅ **身份认证**: JWT Token
- ✅ **权限控制**: RBAC 基于角色的访问控制
- ✅ **接口限流**: Redis 计数器限流
- ✅ **IP白名单**: 管理后台 IP 访问控制

### 审计日志

```go
// 管理员操作日志
type AdminLog struct {
    ID          string    `bson:"_id,omitempty"`
    AdminID     string    `bson:"admin_id"`      // 管理员ID
    Operation   string    `bson:"operation"`     // 操作类型
    Target      string    `bson:"target"`        // 操作对象
    Details     string    `bson:"details"`       // 操作详情
    IP          string    `bson:"ip"`            // IP地址
    CreatedAt   time.Time `bson:"created_at"`    // 操作时间
}
```

---

## ⚡ 性能优化

### 缓存策略

**多级缓存**：

```
请求 → Redis缓存 → MongoDB → 返回
         ↓ 命中
       直接返回
```

**缓存内容**：
- 用户信息（TTL: 1小时）
- 推荐结果（TTL: 10分钟）
- 热门内容（TTL: 5分钟）

**缓存更新策略**：
- 写操作更新缓存
- 定时刷新热数据
- 缓存穿透防护（空对象缓存）

### 数据库优化

```javascript
// MongoDB 索引设计示例
db.users.createIndex({ "email": 1 }, { unique: true })
db.users.createIndex({ "username": 1 }, { unique: true })
db.transactions.createIndex({ "user_id": 1, "created_at": -1 })
db.wallets.createIndex({ "user_id": 1 }, { unique: true })
```

**查询优化**：
- 合理使用索引
- 避免全表扫描
- 分页查询限制返回条数
- 聚合查询优化

### 应用优化

- ✅ **连接池**: MongoDB、Redis 连接池管理
- ✅ **异步处理**: 耗时操作使用消息队列异步处理
- ✅ **批量操作**: 批量插入、批量更新减少数据库交互
- ✅ **并发控制**: 使用 Go 协程提升并发处理能力

---

## 📊 监控运维

### 日志管理

**结构化日志（Zap）**:

```go
logger.Info("用户登录",
    zap.String("user_id", userID),
    zap.String("ip", ip),
    zap.Duration("duration", duration),
)
```

**日志级别**：
- `DEBUG`: 开发调试
- `INFO`: 关键操作记录
- `WARN`: 警告信息
- `ERROR`: 错误信息

### 监控指标

**应用监控**：
- 请求量 (QPS)
- 响应时间 (P95, P99)
- 错误率
- 并发连接数

**资源监控**：
- CPU 使用率
- 内存使用率
- 磁盘 I/O
- 网络流量

**业务监控**：
- 用户注册数
- 充值金额
- 活跃用户数

### 健康检查

```go
// 健康检查接口
GET /health

// 响应示例
{
    "status": "healthy",
    "components": {
        "mongodb": "up",
        "redis": "up",
        "storage": "up"
    },
    "timestamp": "2025-09-30T12:00:00Z"
}
```

---

## 📚 相关文档

### 详细设计文档

#### MVP核心模块 ✅

- [账号权限系统设计](./账号权限系统设计.md) - JWT认证、RBAC权限
- [钱包系统设计](./钱包系统设计.md) - 充值提现、交易管理
- [推荐服务设计](./推荐服务设计.md) - 个性化推荐算法
- [消息队列设计](./消息队列设计.md) - 异步消息处理
- [文件存储设计](./文件存储设计.md) - 文件上传下载
- [管理后台设计](./管理后台设计.md) - 运营管理后台

#### 新增设计文档 🆕

- [通知服务设计](./通知服务设计.md) - 邮件、短信、站内消息
- [缓存策略设计](./缓存策略设计.md) - Redis缓存策略
- [会话管理设计](./会话管理设计.md) - 会话生命周期管理

#### 后续迭代 ⏸️

- [安全审计设计](./安全审计设计.md) - 操作日志、异常检测

### 架构文档

- [微服务架构划分建议](../微服务架构划分建议.md)
- [项目开发规则](../../architecture/项目开发规则.md)
- [架构设计规范](../../architecture/架构设计规范.md)

### 实施文档

- [测试框架实施进度](../../implementation/测试框架实施进度.md)
- [搁置任务清单](../../implementation/搁置任务清单.md)

---

## 🎯 总结

### 核心原则

1. **模块化设计** - 清晰的模块边界，为未来微服务化预留空间
2. **简单实用** - 技术选型以解决问题为导向，避免过度设计
3. **渐进演进** - 随着业务增长逐步优化和拆分
4. **质量优先** - 代码质量、测试覆盖率、文档完整性

### 演进路线

**当前阶段**（模块化单体）：
- 单一代码库，统一部署
- 模块边界清晰，接口定义明确
- 简单技术栈，快速迭代

**未来阶段**（按需拆分微服务）：
- 根据实际需求拆分独立服务
- 引入服务治理基础设施
- 云原生部署和运维

---

*本文档将随项目发展持续更新* 📝