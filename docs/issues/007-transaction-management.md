# Issue #007: Service 层事务管理缺失

**优先级**: 高 (P0)
**类型**: 架构问题
**状态**: ⚠️ 部分修复
**创建日期**: 2026-03-05
**来源报告**: [后端综合审计报告](../reports/archived/backend-comprehensive-audit-summary-2026-01-26.md)、[后端 Service 分析](../reports/archived/backend-service-analysis-2026-01-26.md)
**审查日期**: 2026-03-05
**审查报告**: [P0问题审查报告](../reports/2026-03-05-p0-issue-audit-report.md)

---

## 审查结果

**状态**: ⚠️ 问题仍存在，但高风险财务链路已基本收敛

### 审查发现

1. ⚠️ **仓库级事务管理仍未全域统一，但已新增通用 Mongo transaction runner**
2. ⚠️ **Service层仍未形成全域统一 `RunInTransaction` 规范**
3. ✅ **wallet 交易/提现服务已接入事务执行**
4. ✅ **作者收入提现申请已接入钱包冻结与双写回滚**
5. ✅ **会员订阅/续费已接入钱包扣款与回滚**
6. ✅ **发布事件总线已切换为持久化存储，不再只依赖内存分发**
7. ✅ **`transaction_service.go` 中的余额回滚 TODO 已消除**

### 证据代码

```go
// service/finance/wallet/transaction_service.go:182-189
if err := s.walletRepo.UpdateBalance(ctx, fromWalletID, -amount); err != nil {
    return fmt.Errorf("更新源钱包余额失败: %w", err)
}
if err := s.walletRepo.UpdateBalance(ctx, toWalletID, amount); err != nil {
    // TODO: 需要回滚  ← 确认问题存在
    return fmt.Errorf("更新目标钱包余额失败: %w", err)
}
```

**影响**: 原先如果第二步失败，第一步已扣款无法回滚；该财务完整性风险现已修复

### 本轮已落地

1. ✅ `WalletRepository` 新增 `RunInTransaction`
2. ✅ Mongo 钱包仓储已实现 `StartSession + WithTransaction`
3. ✅ `Recharge / Consume / Transfer` 已改为事务执行
4. ✅ `CreateWithdrawRequest / ApproveWithdraw / RejectWithdraw` 已改为事务执行
5. ✅ wallet 域已引入显式 `TransactionRunner`，不再由业务方法直接依赖仓储事务细节
6. ✅ `ServiceContainer` 已通过 provider 注入 wallet 事务入口
7. ✅ 已补单测，验证余额更新、状态更新、交易记录失败时会整体回滚
8. ✅ `pkg/transaction` 已提供领域无关的 `Runner`，为后续跨域统一事务入口打底
9. ✅ `AuthorRevenueService.CreateWithdrawalRequest` 已改为统一事务内同时创建作者提现单、钱包提现单并冻结余额
10. ✅ `ServiceContainer` 已为 `AuthorRevenueService` 注入通用 Mongo transaction runner
11. ✅ 已补作者提现单测，验证钱包写失败和扣款失败时会整体回滚
12. ✅ `MembershipService.Subscribe / RenewMembership` 已改为统一事务内完成会员生效与钱包扣款
13. ✅ `ServiceContainer` 已为 `MembershipService` 注入通用 Mongo transaction runner
14. ✅ 已补会员订阅/续费回滚测试，验证扣款失败和流水失败时不会留下已生效会员
15. ✅ `FollowService` 的关注/取关已改为事务内统一更新关注关系和互关状态
16. ✅ 已补关注关系回滚测试，验证互关状态更新失败时不会留下半成功关系
17. ✅ `CommentService` 的回复/删除回复已改为事务内统一维护子评论与父评论回复计数
18. ✅ 已补评论回滚测试，验证回复计数更新失败时不会留下脏回复或错误计数
19. ✅ `PublishService` 对发布/下架事件失败已改为可审计，不再静默吞错
20. ✅ 已补发布事件失败测试，验证失败信息会回写到 `PublicationRecord`
21. ✅ `ServiceContainer` 已切换为 `PersistedEventBus + MongoEventStore`，admin 事件回放接口恢复真实数据源
22. ✅ `PublishEventBusAdapter` 已保留真实业务事件类型，避免 `project.published` / `document.published` 被统一折叠为通用事件名

---

## 设计方案

**设计文档**: [事务管理器实现设计方案](../plans/2026-03-05-transaction-manager-design.md)

### 设计要点

1. **统一事务入口** - 收敛现有散落的事务实现
2. **Service层集成** - 将高风险跨仓储操作统一迁移到 `RunInTransaction`
3. **Repository层支持** - 补足需要事务的仓储接口
4. **依赖注入配置** - 在 ServiceContainer 中注册统一事务管理器

**预计实施时间**: 15小时（2个工作日）

---

## 问题描述

Service 层缺少事务管理机制，导致复杂业务场景下的数据一致性风险。

### 具体问题

#### 1. 缺少事务支持 🔴 P0

**问题**: 当前 Service 层没有事务管理，跨多个 Repository 的操作无法保证原子性。

**示例场景**:
```go
// 创建书单并添加书籍（需要事务）
func (s *BookListService) CreateBookListWithBooks(ctx context.Context, req *CreateRequest) error {
    // 1. 创建书单
    bookList := &BookList{...}
    if err := s.repo.CreateBookList(ctx, bookList); err != nil {
        return err
    }

    // 2. 添加书籍（如果失败，书单已创建，无法回滚）
    for _, bookID := range req.BookIDs {
        if err := s.repo.AddBookToList(ctx, bookList.ID, bookID); err != nil {
            // 💥 这里失败后，上面的书单创建无法回滚
            return err
        }
    }

    return nil
}
```

**影响**:
- 数据不一致
- 无法实现复杂的业务逻辑
- 错误恢复困难

#### 2. 事件驱动架构实施不足 🔴 P0

**问题**: 事件发布和业务操作不在同一事务中，可能导致事件丢失或数据不一致。

```go
// 当前问题模式
func (s *BookService) PublishBook(ctx context.Context, bookID string) error {
    // 1. 更新书籍状态
    if err := s.repo.UpdateBookStatus(ctx, bookID, "published"); err != nil {
        return err
    }

    // 2. 发布事件（如果失败，状态已更新，无法回滚）
    if err := s.eventBus.Publish("book.published", bookID); err != nil {
        // 💥 事件发布失败，但书籍状态已更新
        log.Error("Failed to publish event", err)
        // 这里应该回滚，但没有事务支持
    }

    return nil
}
```

---

## 解决方案

### 1. 事务管理器实现

```go
// pkg/transaction/manager.go
package transaction

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
)

type Manager interface {
    // Begin 开始一个新事务
    Begin(ctx context.Context) (context.Context, error)

    // Commit 提交事务
    Commit(ctx context.Context) error

    // Rollback 回滚事务
    Rollback(ctx context.Context) error

    // RunInTransaction 在事务中执行函数
    RunInTransaction(ctx context.Context, fn func(context.Context) error) error
}

type mongoTransactionManager struct {
    client *mongo.Client
    db     *mongo.Database
}

func (m *mongoTransactionManager) RunInTransaction(
    ctx context.Context,
    fn func(context.Context) error,
) error {
    session, err := m.client.StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(ctx)

    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        return nil, fn(sessCtx)
    })

    return err
}
```

### 2. Service 层事务支持

```go
// service/book/book_service.go
func (s *BookService) CreateBookWithChapters(
    ctx context.Context,
    req *CreateBookRequest,
) error {
    return s.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {
        // 1. 创建书籍
        book := &Book{
            Title:       req.Title,
            Description: req.Description,
            AuthorID:    req.AuthorID,
            Status:      "draft",
        }
        if err := s.repo.CreateBook(txCtx, book); err != nil {
            return err
        }

        // 2. 创建章节
        for i, chapterReq := range req.Chapters {
            chapter := &Chapter{
                BookID:        book.ID,
                ChapterNumber: i + 1,
                Title:         chapterReq.Title,
                Content:       chapterReq.Content,
            }
            if err := s.repo.CreateChapter(txCtx, chapter); err != nil {
                // ✅ 在事务中，会自动回滚
                return err
            }
        }

        // 3. 发布事件
        if err := s.eventBus.Publish(txCtx, "book.created", book); err != nil {
            return err
        }

        return nil
    })
}
```

### 3. Repository 层事务支持

```go
// repository/mongodb/base/transactional_repository.go
package base

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
)

type TransactionalRepository struct {
    client *mongo.Client
    db     *mongo.Database
}

// GetCollection 获取事务感知的集合
func (r *TransactionalRepository) GetCollection(ctx context.Context) *mongo.Collection {
    // 如果在事务中，使用会话
    if session, ok := ctx.Value("mongoSession").(mongo.Session); ok {
        return r.db.Collection("books").WithSession(session)
    }
    return r.db.Collection("books")
}
```

---

## 实施计划

### Phase 1: 基础设施（1 周）

1. **统一事务管理器**
   - [ ] 收敛已有事务接口/实现
   - [ ] 统一 MongoDB 事务管理器入口
   - [ ] 编写单元测试

2. **集成到 ServiceContainer**
   - [ ] 注册事务管理器
   - [ ] 更新依赖注入

### Phase 2: 核心业务迁移（2-3 周）

**优先级排序**:
1. 书籍管理（创建书籍+章节）
2. 订单管理（创建订单+支付）
3. 用户充值（创建流水+更新余额）
4. 社交互动（创建评论+更新计数）

### Phase 3: 事件集成（1-2 周）

1. **事务性事件发布**
   - [ ] 实现事务事件队列
   - [ ] 事务提交后发布事件
   - [ ] 处理事件发布失败

### Phase 4: 测试和验证（1 周）

1. **集成测试**
   - [ ] 事务回滚测试
   - [ ] 并发事务测试
   - [ ] 性能测试

---

## 使用示例

### 创建带事务的 Service

```go
// service/order/order_service.go
type OrderService struct {
    txManager  transaction.Manager
    orderRepo  repository.OrderRepository
    paymentSvc *PaymentService
}

func (s *OrderService) CreateOrder(
    ctx context.Context,
    req *CreateOrderRequest,
) (*Order, error) {
    var order *Order

    err := s.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {
        // 1. 创建订单
        order = &Order{
            UserID:  req.UserID,
            BookID:  req.BookID,
            Amount:  req.Amount,
            Status:  "pending",
        }
        if err := s.orderRepo.CreateOrder(txCtx, order); err != nil {
            return err
        }

        // 2. 创建支付记录
        payment := &Payment{
            OrderID: order.ID,
            Amount:  order.Amount,
            Status:  "pending",
        }
        if err := s.paymentSvc.CreatePayment(txCtx, payment); err != nil {
            return err  // ✅ 自动回滚订单创建
        }

        // 3. 更新书籍销量
        if err := s.bookService.IncrementSales(txCtx, order.BookID, 1); err != nil {
            return err  // ✅ 自动回滚所有操作
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    return order, nil
}
```

---

## 注意事项

### MongoDB 事务限制

1. **事务大小限制**: 16MB
2. **嵌套限制**: 不支持嵌套事务
3. **集合限制**: 事务中的所有集合必须在同一分片（如果使用分片集群）
4. **性能影响**: 事务操作有性能开销

### 最佳实践

1. **保持事务简短**: 尽量减少事务中的操作数量
2. **避免长时间等待**: 不要在事务中进行外部 API 调用
3. **正确处理错误**: 确保错误时回滚
4. **使用重试机制**: 处理事务冲突

---

## 检查清单

### 基础设施
- [x] 统一事务管理器基础实现
- [x] wallet 财务链路单元测试覆盖
- [x] 财务高风险链路已集成到 DI 容器

### 业务迁移
- [ ] 书籍管理事务支持
- [ ] 订单管理事务支持
- [x] 钱包充值/消费/转账/提现事务支持
- [x] 作者收入提现申请事务支持
- [x] 会员订阅/续费事务支持
- [x] 关注/取关事务支持
- [x] 评论回复/删除回复事务支持
- [x] 发布事件失败可见化
- [ ] 点赞等其他社交互动事务支持

### 测试验证
- [x] 高风险财务回滚测试
- [ ] 并发测试
- [ ] 性能测试

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [后端 Service 分析](../reports/archived/backend-service-analysis-2026-01-26.md) | Service 层详细分析 |
| [设计审查 - Service 层](../reports/archived/design-review-block6-service-layer-20260127.md) | Service 层设计审查 |

---

## 相关Issue

### 相关Issue（联合处理）
- [#010: Repository 层业务逻辑渗透](./010-repository-layer-business-logic-permeation.md) - ⚠️ Repository层中的跨表事务需要依赖事务管理器
- [#005: API 标准化问题](./005-api-standardization-issues.md) - 事务失败后的错误响应需要标准化

### 依赖关系
- 本Issue是 #010 中的跨表事务问题的前置依赖
- 实现事务管理器后，才能将Repository层的事务逻辑移到Service层
