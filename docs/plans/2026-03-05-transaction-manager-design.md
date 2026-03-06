# 事务管理器实现设计方案

**设计日期**: 2026-03-05
**设计者**: Kore
**优先级**: 🔴 P0
**问题来源**: Issue #007: Service 层事务管理缺失

---

## 问题描述

### 当前问题

1. **缺少事务管理器**：`pkg/transaction/` 目录不存在
2. **无统一事务接口**：Service层无法使用 `RunInTransaction` 模式
3. **事务回滚缺失**：`transaction_service.go:187` 存在 `// TODO: 需要回滚` 注释

### 问题证据

```go
// service/finance/wallet/transaction_service.go:182-189
// 8. 更新余额
if err := s.walletRepo.UpdateBalance(ctx, fromWalletID, -amount); err != nil {
    return fmt.Errorf("更新源钱包余额失败: %w", err)
}

if err := s.walletRepo.UpdateBalance(ctx, toWalletID, amount); err != nil {
    // TODO: 需要回滚  ← 确认问题存在
    return fmt.Errorf("更新目标钱包余额失败: %w", err)
}
```

**影响**：
- 如果第二步失败，第一步已扣款无法回滚
- 用户余额会不一致
- 财务数据完整性风险

---

## 设计方案

### 架构设计

```
┌─────────────────────────────────────────────┐
│            Service 层                          │
│  ┌───────────────────────────────────────┐  │
│  │  TransactionService                    │  │
│  │                                          │  │
│  │  func (s *TransactionService)           │  │
│  │      Transfer(ctx, from, to, amount) {  │  │
│  │                                          │  │
│  │      s.txManager.RunInTransaction(      │  │
│  │          ctx,                            │  │
│  │          func(txCtx context.Context) {    │  │
│  │                                              │
│  │              // 1. 验证                      │  │
│  │              // 2. 扣款                    │  │
│  │              // 3. 入账                    │  │
│  │              // 4. 记录                    │  │
│  │          }                               │  │
│  │      )                                   │  │
│  │  }                                       │  │
│  └───────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────┐
│         Transaction Manager                    │
│  ┌───────────────────────────────────────┐  │
│  │  TxManager (interface)                 │  │
│  │                                          │  │
│  │  - Begin(ctx) (context, error)          │  │
│  │  - Commit(ctx) error                    │  │
│  │  - Rollback(ctx) error                  │  │
│  │  - RunInTransaction(ctx, fn) error       │  │
│  └───────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────┐
│         MongoDB Session                      │
│  ┌───────────────────────────────────────┐  │
│  │  session.WithTransaction(...)           │  │
│  │                                          │  │
│  │  - 自动开始事务                          │  │
│  │  - fn返回nil → Commit                   │  │
│  │  - fn返回error → Rollback               │  │
│  └───────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
```

---

## 接口设计

### 1. 事务管理器接口

```go
// pkg/transaction/manager.go
package transaction

import (
    "context"
)

// Manager 事务管理器接口
type Manager interface {
    // Begin 开始一个新事务
    // 返回: 事务上下文, 错误
    Begin(ctx context.Context) (context.Context, error)

    // Commit 提交事务
    // 返回: 错误
    Commit(ctx context.Context) error

    // Rollback 回滚事务
    // 返回: 错误
    Rollback(ctx context.Context) error

    // RunInTransaction 在事务中执行函数
    // fn 返回nil → 提交事务
    // fn 返回error → 回滚事务
    // 返回: fn的返回错误
    RunInTransaction(ctx context.Context, fn func(context.Context) error) error
}
```

### 2. MongoDB事务管理器实现

```go
// pkg/transaction/mongo_manager.go
package transaction

import (
    "context"
    "fmt"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type mongoTransactionManager struct {
    client *mongo.Client
    db     *mongo.Database
}

func NewMongoManager(client *mongo.Client, database string) Manager {
    return &mongoTransactionManager{
        client: client,
        db:     client.Database(database),
    }
}

// Begin 开始事务
func (m *mongoTransactionManager) Begin(ctx context.Context) (context.Context, error) {
    session, err := m.client.StartSession()
    if err != nil {
        return nil, fmt.Errorf("failed to start session: %w", err)
    }

    // 将session存储到context中
    txCtx := context.WithValue(ctx, transactionSessionKey{}, session)

    return txCtx, nil
}

// Commit 提交事务
func (m *mongoTransactionManager) Commit(ctx context.Context) error {
    session, ok := getTransactionSession(ctx)
    if !ok {
        return fmt.Errorf("no transaction session in context")
    }

    // 使用 WithTransaction自动提交
    _, err := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        return nil, session.CommitTransaction(sessCtx)
    })

    // 清理session
    session.EndSession(ctx)

    return err
}

// Rollback 回滚事务
func (m *mongoTransactionManager) Rollback(ctx context.Context) error {
    session, ok := getTransactionSession(ctx)
    if !ok {
        return fmt.Errorf("no transaction session in context")
    }

    // 使用WithTransaction自动回滚
    _, err := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        return nil, session.AbortTransaction(sessCtx)
    })

    // 清理session
    session.EndSession(ctx)

    return err
}

// RunInTransaction 在事务中执行函数
func (m *mongoTransactionManager) RunInTransaction(
    ctx context.Context,
    fn func(context.Context) error,
) error {
    session, err := m.client.StartSession()
    if err != nil {
        return fmt.Errorf("failed to start session: %w", err)
    }
    defer session.EndSession(ctx)

    // 使用 WithTransaction简化事务管理
    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        // 在事务中执行用户函数
        if err := fn(sessCtx); err != nil {
            return nil, err  // 返回错误会导致自动回滚
        }
        return nil, nil  // 返回nil会导致自动提交
    })

    return err
}

// 事务session的context key
type transactionSessionKey struct{}

// getTransactionSession 从context获取事务session
func getTransactionSession(ctx context.Context) (mongo.Session, bool) {
    session, ok := ctx.Value(transactionSessionKey{}).(mongo.Session)
    return session, ok
}
```

---

## Service层集成

### 1. 更新TransactionService

```go
// service/finance/wallet/transaction_service.go
type TransactionServiceImpl struct {
    walletRepo   interfaces.WalletRepository
    txManager    transaction.Manager  // ← 新增
    transactionRepo interfaces.TransactionRepository
}

func NewTransactionService(
    walletRepo interfaces.WalletRepository,
    txManager transaction.Manager,  // ← 新增参数
    transactionRepo interfaces.TransactionRepository,
) *TransactionServiceImpl {
    return &TransactionServiceImpl{
        walletRepo:   walletRepo,
        txManager:    txManager,
        transactionRepo: transactionRepo,
    }
}

// Transfer 转账（使用事务）
func (s *TransactionServiceImpl) Transfer(
    ctx context.Context,
    req *TransferRequest,
) error {
    var txError error

    // 使用事务执行转账
    err := s.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {
        // 1. 验证钱包存在
        fromWallet, err := s.walletRepo.GetWallet(txCtx, req.FromUserID)
        if err != nil {
            return fmt.Errorf("源钱包不存在: %w", err)
        }

        toWallet, err := s.walletRepo.GetWallet(txCtx, req.ToUserID)
        if err != nil {
            return fmt.Errorf("目标钱包不存在: %w", err)
        }

        // 2. 检查钱包状态
        if fromWallet.Frozen || toWallet.Frozen {
            return fmt.Errorf("钱包已冻结，无法转账")
        }

        // 3. 检查余额
        if fromWallet.Balance < types.Money(req.Amount) {
            return fmt.Errorf("余额不足")
        }

        // 4. 创建转出交易记录
        outTransaction := &financeModel.Transaction{
            UserID:        req.FromUserID,
            Type:          "transfer_out",
            Amount:        types.Money(-req.Amount),
            Status:        "success",
            Reason:        req.Reason + " → " + req.ToUserID,
            RelatedUserID: req.ToUserID,
        }

        if err := s.transactionRepo.Create(txCtx, outTransaction); err != nil {
            return fmt.Errorf("创建转出记录失败: %w", err)
        }

        // 5. 创建转入交易记录
        inTransaction := &financeModel.Transaction{
            UserID:        req.ToUserID,
            Type:          "transfer_in",
            Amount:        types.Money(req.Amount),
            Status:        "success",
            Reason:        req.FromUserID + " → " + req.Reason,
            RelatedUserID: req.FromUserID,
        }

        if err := s.transactionRepo.Create(txCtx, inTransaction); err != nil {
            return fmt.Errorf("创建转入记录失败: %w", err)
        }

        // 6. 更新余额（在事务中）
        if err := s.walletRepo.UpdateBalance(txCtx, req.FromUserID, -req.Amount); err != nil {
            return fmt.Errorf("更新源钱包余额失败: %w", err)
        }

        if err := s.walletRepo.UpdateBalance(txCtx, req.ToUserID, req.Amount); err != nil {
            // ✅ 在事务中，会自动回滚上面的操作
            return fmt.Errorf("更新目标钱包余额失败: %w", err)
        }

        return nil
    })

    if err != nil {
        txError = err
    }

    return txError
}
```

---

## Repository层事务支持

### 事务感知的Repository

Repository需要支持传入事务上下文：

```go
// repository/mongodb/base/transactional_repository.go
package base

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
)

// TransactionalRepository 事务感知的Repository基类
type TransactionalRepository struct {
    client *mongo.Client
    db     *mongo.Database
}

// GetCollection 获取事务感知的集合
func (r *TransactionalRepository) GetCollection(ctx context.Context) *mongo.Collection {
    collection := r.db.Collection("books")

    // 检查是否在事务中
    if session, ok := ctx.Value("transactionSession").(mongo.Session); ok {
        // 使用事务session
        return collection.WithSession(session)
    }

    return collection
}
```

---

## 依赖注入配置

### ServiceContainer配置

```go
// internal/container/service_container.go
func NewServiceContainer(db *mongo.Database, client *mongo.Client) *ServiceContainer {
    // 创建事务管理器
    txManager := transaction.NewMongoManager(client, db.Database().Name())

    return &ServiceContainer{
        // ...
        txManager: txManager,
        // ...
    }
}

// 注入TransactionService
func (c *ServiceContainer) TransactionService() interfaces.TransactionService {
    return transaction.NewTransactionService(
        c.WalletRepository(),
        c.txManager,  // ← 注入事务管理器
        c.TransactionRepository(),
    )
}
```

---

## 迁移步骤

### Phase 1: 创建事务管理器（1天）

1. **创建目录结构**:
   ```
   pkg/transaction/
   ├── manager.go           # 接口定义
   ├── mongo_manager.go     # MongoDB实现
   └── README.md
   ```

2. **实现接口**:
   - `Manager` 接口
   - `MongoManager` 实现
   - 单元测试

3. **集成到DI容器**:
   - 在 `ServiceContainer` 中注册事务管理器
   - 更新 `NewTransactionService` 构造函数

### Phase 2: 修复现有服务（2-3天）

1. **修复 Transfer 方法**:
   - 使用 `RunInTransaction` 包装转账逻辑
   - 移除 `// TODO: 需要回滚` 注释

2. **检查其他需要事务的服务**:
   - `Consume` 方法
   - `Withdraw` 方法
   - 创建带事务的Service方法

### Phase 3: Repository层事务支持（1天）

1. **创建事务感知Repository基类**:
   - `TransactionalRepository` 基类
   - `GetCollection` 方法支持事务上下文

2. **更新现有Repository**:
   - 继承 `TransactionalRepository`
   - 使用 `GetCollection(ctx)` 获取集合

### Phase 4: 测试验证（1天）

1. **单元测试**:
   - 事务管理器测试
   - Service层事务测试

2. **集成测试**:
   - 事务回滚测试
   - 并发事务测试

3. **API测试**:
   - 转账API测试
   - 失败回滚验证

---

## 验证清单

### 事务管理器
- [ ] `pkg/transaction/` 目录已创建
- [ ] `Manager` 接口已定义
- [ ] `MongoManager` 实现完成
- [ ] 单元测试通过

### Service层
- [ ] `TransactionService` 使用事务管理器
- [ ] TODO注释已移除
- [ ] Transfer方法在事务中执行

### Repository层
- [ ] Repository支持事务上下文
- [ ] 集成测试验证事务隔离

### 测试验证
- [ ] 事务回滚测试通过
- [ ] 并发事务测试通过
- [ ] API测试验证

---

## 使用示例

### Service层使用事务

```go
// 1. 简单事务
func (s *SomeService) CreateWithTransaction(ctx context.Context) error {
    return s.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {
        // 所有数据库操作在txCtx中执行
        if err := s.repoA.Create(txCtx, dataA); err != nil {
            return err  // 自动回滚
        }
        if err := s.repoB.Create(txCtx, dataB); err != nil {
            return err  // 自动回滚
        }
        return nil  // 自动提交
    })
}

// 2. 带验证的事务
func (s *TransactionService) Transfer(ctx context.Context, req *TransferRequest) error {
    return s.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {
        // 1. 业务验证
        if err := s.validateTransfer(txCtx, req); err != nil {
            return err
        }

        // 2. 执行业务逻辑
        if err := s.executeTransfer(txCtx, req); err != nil {
            return err
        }

        return nil
    })
}
```

---

## 注意事项

### MongoDB事务限制

1. **事务大小限制**: 16MB
2. **嵌套限制**: 不支持嵌套事务
3. **集合限制**: 事务中的所有集合必须在同一分片（如果使用分片集群）
4. **性能影响**: 事务操作有性能开销

### 最佳实践

1. **保持事务简短**: 尽量减少事务中的操作数量
2. **避免长时间等待**: 不要在事务中进行外部API调用
3. **正确处理错误**: 确保错误时回滚
4. **使用重试机制**: 处理事务冲突

---

## 实施计划

| 步骤 | 任务 | 预计时间 |
|------|------|----------|
| 1. 创建事务管理器 | 实现 Manager 接口和 MongoManager | 4小时 |
| 2. 集成到DI容器 | 更新 ServiceContainer | 2小时 |
| 3. 修复 Transfer 方法 | 使用事务管理器 | 2小时 |
| 4. 创建Repository基类 | TransactionalRepository | 2小时 |
| 5. 测试验证 | 单元测试、集成测试 | 4小时 |
| 6. 文档更新 | README和使用示例 | 1小时 |

**总预计时间**: 15小时（2个工作日）

---

## 相关文档

- [Issue #007: Service 层事务管理缺失](../issues/007-transaction-management.md)
- [Issue #010: Repository 层业务逻辑渗透](../issues/010-repository-layer-business-logic-permeation.md)
- [P0问题审查报告](../reports/2026-03-05-p0-issue-audit-report.md)

---

**设计完成时间**: 2026-03-05
**预计实施时间**: 2个工作日
**建议执行者**: 后端团队
