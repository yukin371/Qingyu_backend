# Saga事务模式设计

> 最后整理: 2026-05-22  
> 当前状态: `legacy-live`
> 历史版本: v1.0  
> 创建时间: 2025-10-21  
> 旧状态标记: 设计阶段

## Page Role

- 这里负责：Saga 事务模式的历史设计方案，包括数据模型、核心业务场景、补偿机制与实施建议。
- 不负责：当前事务协调实现事实、现行服务边界、当前分布式事务方案结论。

## Recommended Read Path

1. [README.md](./README.md)
2. [../../../architecture/README.md](../../../architecture/README.md)
3. [../../../implementation/README.md](../../../implementation/README.md)
4. [../../../standards/README.md](../../../standards/README.md)

## Boundary

- 如果你要找“当前架构与服务边界”，优先看 [../../../architecture/README.md](../../../architecture/README.md)。
- 如果你要找“已落地的实施记录”，优先看 [../../../implementation/README.md](../../../implementation/README.md)。
- 如果你要找“当前规则和标准”，优先看 [../../../standards/README.md](../../../standards/README.md)。
- 本页更适合解释 Saga 当时如何规划，不适合作为当前实现事实的唯一依据。

## Quick Section Map

| 如果你想看 | 直接跳到 |
|------|------|
| Saga 模式背景与目标 | 见“文档概述”“设计目标”“一、Saga模式概述” |
| 数据模型 | 见“二、数据模型设计” |
| 关键业务场景 | 见“三、核心业务场景设计” |
| MQ 补偿与最终一致性 | 见“四、基于消息队列的最终一致性” |
| Repository 接口与实施建议 | 见“五、Repository接口”“六、实施建议” |

## Quick Takeaways

- 这篇最有价值的地方，是它把 Saga 数据模型、业务场景和补偿流程放成了一份完整的最终一致性方案。
- 如果你只关心 Saga 主线，优先看“一、Saga模式概述”“三、核心业务场景设计”“四、基于消息队列的最终一致性”。
- 文中的跨服务场景和实施建议属于历史设计阶段语境，不应直接视为当前实现事实。

## Skip Guide

- 只想知道“Saga 适不适合”：看“一、Saga模式概述”。
- 只想知道“购买章节/提现怎么编排”：看“三、核心业务场景设计”。
- 如果你当前只关心现行架构或落地状态，请优先回到 `docs/architecture/`、`docs/implementation/`。

## 📋 文档概述

本文档设计青羽平台的Saga事务模式，用于保证跨服务的分布式事务最终一致性，特别是在Monetization和Wallet Service之间的关键业务场景。

## 🎯 设计目标

1. **最终一致性**：保证跨服务事务的最终一致性
2. **补偿机制**：每个步骤都有对应的补偿操作
3. **可恢复性**：支持从失败中恢复
4. **可追踪性**：完整的Saga执行日志

---

## 一、Saga模式概述

### 1.1 Saga vs 2PC/3PC

| 特性 | Saga | 2PC/3PC |
|------|------|---------|
| **一致性** | 最终一致性 | 强一致性 |
| **性能** | 高（无锁） | 低（有锁） |
| **可用性** | 高 | 低（协调者单点故障） |
| **实现复杂度** | 中 | 高 |
| **适用场景** | 跨服务长事务 | 单数据库短事务 |

### 1.2 Saga工作流

```
正向流程（成功）：
Step1 → Step2 → Step3 → Complete

补偿流程（失败）：
Step1 → Step2 → Step3(失败) → Compensate3 → Compensate2 → Compensate1 → Abort
```

---

## 二、数据模型设计

### 2.1 SagaLog数据模型

```go
package models

type SagaLog struct {
    ID          string        `bson:"_id" json:"id"`
    SagaType    string        `bson:"saga_type" json:"sagaType"` // purchase_chapter, withdrawal
    Status      SagaStatus    `bson:"status" json:"status"`
    Steps       []SagaStep    `bson:"steps" json:"steps"`
    Context     SagaContext   `bson:"context" json:"context"`
    CreatedAt   time.Time     `bson:"created_at" json:"createdAt"`
    UpdatedAt   time.Time     `bson:"updated_at" json:"updatedAt"`
    CompletedAt *time.Time    `bson:"completed_at,omitempty" json:"completedAt,omitempty"`
}

type SagaStatus string

const (
    SagaStatusPending    SagaStatus = "pending"     // 待执行
    SagaStatusRunning    SagaStatus = "running"     // 执行中
    SagaStatusCompleted  SagaStatus = "completed"   // 完成
    SagaStatusFailed     SagaStatus = "failed"      // 失败
    SagaStatusCompensating SagaStatus = "compensating" // 补偿中
    SagaStatusAborted    SagaStatus = "aborted"     // 已中止
)

type SagaStep struct {
    StepID         string          `bson:"step_id" json:"stepId"`
    StepName       string          `bson:"step_name" json:"stepName"`
    Service        string          `bson:"service" json:"service"`
    Action         string          `bson:"action" json:"action"`
    Status         StepStatus      `bson:"status" json:"status"`
    Input          json.RawMessage `bson:"input" json:"input"`
    Output         json.RawMessage `bson:"output,omitempty" json:"output,omitempty"`
    Error          string          `bson:"error,omitempty" json:"error,omitempty"`
    Compensatable  bool            `bson:"compensatable" json:"compensatable"`
    Compensated    bool            `bson:"compensated" json:"compensated"`
    StartedAt      *time.Time      `bson:"started_at,omitempty" json:"startedAt,omitempty"`
    CompletedAt    *time.Time      `bson:"completed_at,omitempty" json:"completedAt,omitempty"`
    RetryCount     int             `bson:"retry_count" json:"retryCount"`
}

type StepStatus string

const (
    StepStatusPending    StepStatus = "pending"
    StepStatusRunning    StepStatus = "running"
    StepStatusCompleted  StepStatus = "completed"
    StepStatusFailed     StepStatus = "failed"
    StepStatusCompensated StepStatus = "compensated"
)

type SagaContext struct {
    UserID    string                 `bson:"user_id" json:"userId"`
    Data      map[string]interface{} `bson:"data" json:"data"`
    Metadata  map[string]string      `bson:"metadata" json:"metadata"`
}
```

---

## 三、核心业务场景设计

### 3.1 场景1：用户购买章节

#### 3.1.1 Saga流程设计

```go
package saga

// PurchaseChapterSaga 购买章节Saga
type PurchaseChapterSaga struct {
    sagaRepo            SagaRepository
    walletService       wallet.WalletService
    bookstoreService    bookstore.BookstoreService
    monetizationService monetization.MonetizationService
    eventBus            base.EventBus
}

// PurchaseChapterRequest 购买请求
type PurchaseChapterRequest struct {
    UserID    string  `json:"userId"`
    ChapterID string  `json:"chapterId"`
    BookID    string  `json:"bookId"`
    AuthorID  string  `json:"authorId"`
    Amount    float64 `json:"amount"`
}

// Execute 执行Saga
func (saga *PurchaseChapterSaga) Execute(ctx context.Context, req *PurchaseChapterRequest) error {
    // 1. 创建Saga Log
    sagaLog := &SagaLog{
        ID:       uuid.New().String(),
        SagaType: "purchase_chapter",
        Status:   SagaStatusRunning,
        Steps:    []SagaStep{},
        Context: SagaContext{
            UserID: req.UserID,
            Data: map[string]interface{}{
                "chapter_id": req.ChapterID,
                "book_id":    req.BookID,
                "author_id":  req.AuthorID,
                "amount":     req.Amount,
            },
        },
        CreatedAt: time.Now(),
    }
    
    if err := saga.sagaRepo.Create(ctx, sagaLog); err != nil {
        return err
    }
    
    // 2. Step 1: 扣除读者代币
    step1 := saga.createStep("deduct_reader_tokens", "WalletService", "deduct", req)
    if err := saga.executeStep(ctx, sagaLog, step1, func() error {
        _, err := saga.walletService.Consume(ctx, req.UserID, req.Amount, fmt.Sprintf("purchase_chapter_%s", req.ChapterID))
        return err
    }); err != nil {
        return saga.compensate(ctx, sagaLog)
    }
    
    // 3. Step 2: 解锁章节权限
    step2 := saga.createStep("unlock_chapter", "BookstoreService", "unlock", req)
    if err := saga.executeStep(ctx, sagaLog, step2, func() error {
        return saga.bookstoreService.UnlockChapter(ctx, req.UserID, req.ChapterID)
    }); err != nil {
        return saga.compensate(ctx, sagaLog)
    }
    
    // 4. Step 3: 增加作者收益（异步，允许失败）
    step3 := saga.createStep("author_revenue", "MonetizationService", "add_revenue", req)
    step3.Compensatable = false // 非关键步骤，失败不回滚
    if err := saga.executeStep(ctx, sagaLog, step3, func() error {
        shareAmount := req.Amount * 0.7 // 70%分成
        return saga.monetizationService.RecordRevenue(ctx, req.AuthorID, req.BookID, req.ChapterID, shareAmount)
    }); err != nil {
        // 发布补偿事件，后台处理
        saga.publishCompensationEvent(ctx, sagaLog.ID, step3)
        log.Warn("作者收益记录失败，已加入补偿队列", "error", err)
    }
    
    // 5. 完成Saga
    sagaLog.Status = SagaStatusCompleted
    now := time.Now()
    sagaLog.CompletedAt = &now
    saga.sagaRepo.Update(ctx, sagaLog)
    
    // 6. 发布事件
    saga.eventBus.PublishAsync(ctx, &ChapterPurchasedEvent{
        UserID:    req.UserID,
        ChapterID: req.ChapterID,
        Amount:    req.Amount,
    })
    
    return nil
}

// executeStep 执行步骤
func (saga *PurchaseChapterSaga) executeStep(
    ctx context.Context,
    sagaLog *SagaLog,
    step *SagaStep,
    action func() error,
) error {
    // 1. 更新步骤状态
    step.Status = StepStatusRunning
    now := time.Now()
    step.StartedAt = &now
    sagaLog.Steps = append(sagaLog.Steps, *step)
    saga.sagaRepo.Update(ctx, sagaLog)
    
    // 2. 执行操作
    err := action()
    
    // 3. 记录结果
    if err != nil {
        step.Status = StepStatusFailed
        step.Error = err.Error()
        step.RetryCount++
        saga.sagaRepo.UpdateStep(ctx, sagaLog.ID, step)
        return err
    }
    
    step.Status = StepStatusCompleted
    step.CompletedAt = &now
    saga.sagaRepo.UpdateStep(ctx, sagaLog.ID, step)
    
    return nil
}

// compensate 补偿操作
func (saga *PurchaseChapterSaga) compensate(ctx context.Context, sagaLog *SagaLog) error {
    sagaLog.Status = SagaStatusCompensating
    saga.sagaRepo.Update(ctx, sagaLog)
    
    // 逆序执行补偿
    for i := len(sagaLog.Steps) - 1; i >= 0; i-- {
        step := &sagaLog.Steps[i]
        
        // 只补偿已完成且可补偿的步骤
        if step.Status != StepStatusCompleted || !step.Compensatable {
            continue
        }
        
        // 执行补偿
        if err := saga.compensateStep(ctx, sagaLog, step); err != nil {
            log.Error("补偿步骤失败", "step", step.StepName, "error", err)
            // 继续尝试补偿其他步骤
        }
    }
    
    sagaLog.Status = SagaStatusAborted
    saga.sagaRepo.Update(ctx, sagaLog)
    
    return fmt.Errorf("saga执行失败，已补偿")
}

// compensateStep 补偿单个步骤
func (saga *PurchaseChapterSaga) compensateStep(ctx context.Context, sagaLog *SagaLog, step *SagaStep) error {
    switch step.StepName {
    case "deduct_reader_tokens":
        // 补偿：退还代币
        amount := sagaLog.Context.Data["amount"].(float64)
        userID := sagaLog.Context.UserID
        _, err := saga.walletService.Recharge(ctx, userID, amount, fmt.Sprintf("refund_saga_%s", sagaLog.ID))
        if err != nil {
            return err
        }
        
    case "unlock_chapter":
        // 补偿：撤销解锁
        userID := sagaLog.Context.UserID
        chapterID := sagaLog.Context.Data["chapter_id"].(string)
        err := saga.bookstoreService.RevokeChapterAccess(ctx, userID, chapterID)
        if err != nil {
            return err
        }
    }
    
    step.Compensated = true
    step.Status = StepStatusCompensated
    saga.sagaRepo.UpdateStep(ctx, sagaLog.ID, step)
    
    return nil
}

// createStep 创建步骤
func (saga *PurchaseChapterSaga) createStep(name, service, action string, req *PurchaseChapterRequest) *SagaStep {
    input, _ := json.Marshal(req)
    return &SagaStep{
        StepID:        uuid.New().String(),
        StepName:      name,
        Service:       service,
        Action:        action,
        Status:        StepStatusPending,
        Input:         input,
        Compensatable: true,
        RetryCount:    0,
    }
}
```

### 3.2 场景2：作者提现

```go
// WithdrawalSaga 提现Saga
type WithdrawalSaga struct {
    sagaRepo            SagaRepository
    walletService       wallet.WalletService
    monetizationService monetization.MonetizationService
    paymentGateway      payment.PaymentGateway
}

func (saga *WithdrawalSaga) Execute(ctx context.Context, req *WithdrawalRequest) error {
    sagaLog := &SagaLog{
        ID:       uuid.New().String(),
        SagaType: "withdrawal",
        Status:   SagaStatusRunning,
        Context: SagaContext{
            UserID: req.AuthorID,
            Data: map[string]interface{}{
                "amount":  req.Amount,
                "account": req.Account,
                "method":  req.Method,
            },
        },
    }
    
    saga.sagaRepo.Create(ctx, sagaLog)
    
    // Step 1: 冻结资金
    step1 := saga.createStep("freeze_funds", "WalletService", "freeze")
    if err := saga.executeStep(ctx, sagaLog, step1, func() error {
        return saga.walletService.FreezeForWithdrawal(ctx, req.AuthorID, req.Amount, sagaLog.ID)
    }); err != nil {
        return saga.compensate(ctx, sagaLog)
    }
    
    // Step 2: 创建提现订单
    var orderID string
    step2 := saga.createStep("create_order", "MonetizationService", "create_withdrawal_order")
    if err := saga.executeStep(ctx, sagaLog, step2, func() error {
        order, err := saga.monetizationService.CreateWithdrawalOrder(ctx, req)
        if err != nil {
            return err
        }
        orderID = order.ID
        sagaLog.Context.Data["order_id"] = orderID
        return nil
    }); err != nil {
        return saga.compensate(ctx, sagaLog)
    }
    
    // Step 3: 调用支付接口
    step3 := saga.createStep("payment_gateway", "PaymentGateway", "transfer")
    if err := saga.executeStep(ctx, sagaLog, step3, func() error {
        return saga.paymentGateway.Transfer(ctx, &payment.TransferRequest{
            Amount:  req.Amount,
            Account: req.Account,
            Method:  req.Method,
        })
    }); err != nil {
        return saga.compensate(ctx, sagaLog)
    }
    
    // Step 4: 处理提现（扣除余额）
    step4 := saga.createStep("process_withdrawal", "WalletService", "process")
    if err := saga.executeStep(ctx, sagaLog, step4, func() error {
        return saga.walletService.ProcessWithdrawal(ctx, orderID)
    }); err != nil {
        return saga.compensate(ctx, sagaLog)
    }
    
    // Step 5: 更新订单状态
    step5 := saga.createStep("update_order", "MonetizationService", "complete_order")
    step5.Compensatable = false
    saga.executeStep(ctx, sagaLog, step5, func() error {
        return saga.monetizationService.CompleteWithdrawalOrder(ctx, orderID)
    })
    
    sagaLog.Status = SagaStatusCompleted
    saga.sagaRepo.Update(ctx, sagaLog)
    
    return nil
}

func (saga *WithdrawalSaga) compensate(ctx context.Context, sagaLog *SagaLog) error {
    // 逆序补偿
    for i := len(sagaLog.Steps) - 1; i >= 0; i-- {
        step := &sagaLog.Steps[i]
        
        if step.Status != StepStatusCompleted || !step.Compensatable {
            continue
        }
        
        saga.compensateStep(ctx, sagaLog, step)
    }
    
    sagaLog.Status = SagaStatusAborted
    saga.sagaRepo.Update(ctx, sagaLog)
    
    return fmt.Errorf("提现失败，已补偿")
}

func (saga *WithdrawalSaga) compensateStep(ctx context.Context, sagaLog *SagaLog, step *SagaStep) error {
    switch step.StepName {
    case "freeze_funds":
        // 解冻资金
        authorID := sagaLog.Context.UserID
        return saga.walletService.CancelWithdrawal(ctx, sagaLog.ID)
        
    case "create_order":
        // 取消订单
        orderID := sagaLog.Context.Data["order_id"].(string)
        return saga.monetizationService.CancelWithdrawalOrder(ctx, orderID)
        
    case "payment_gateway":
        // 请求退款
        // 注意：这可能是异步的
        amount := sagaLog.Context.Data["amount"].(float64)
        account := sagaLog.Context.Data["account"].(string)
        return saga.paymentGateway.Refund(ctx, &payment.RefundRequest{
            Amount:  amount,
            Account: account,
        })
    }
    
    step.Compensated = true
    saga.sagaRepo.UpdateStep(ctx, sagaLog.ID, step)
    
    return nil
}
```

---

## 四、基于消息队列的最终一致性

### 4.1 补偿事件发布

```go
// 发布补偿事件
func (saga *PurchaseChapterSaga) publishCompensationEvent(ctx context.Context, sagaID string, step *SagaStep) {
    event := &SagaCompensationEvent{
        SagaID:    sagaID,
        StepID:    step.StepID,
        StepName:  step.StepName,
        Action:    "compensate",
        Timestamp: time.Now(),
    }
    
    saga.eventBus.Publish(ctx, "saga.compensation", event)
}
```

### 4.2 补偿事件处理器

```go
// SagaCompensationHandler 补偿事件处理器
type SagaCompensationHandler struct {
    sagaRepo            SagaRepository
    monetizationService monetization.MonetizationService
}

func (h *SagaCompensationHandler) Handle(ctx context.Context, event base.Event) error {
    compensationEvent := event.GetEventData().(*SagaCompensationEvent)
    
    // 1. 获取Saga Log
    sagaLog, err := h.sagaRepo.GetByID(ctx, compensationEvent.SagaID)
    if err != nil {
        return err
    }
    
    // 2. 找到对应的步骤
    var step *SagaStep
    for i := range sagaLog.Steps {
        if sagaLog.Steps[i].StepID == compensationEvent.StepID {
            step = &sagaLog.Steps[i]
            break
        }
    }
    
    if step == nil {
        return fmt.Errorf("步骤不存在")
    }
    
    // 3. 执行补偿
    if step.StepName == "author_revenue" {
        // 重试记录作者收益
        authorID := sagaLog.Context.Data["author_id"].(string)
        bookID := sagaLog.Context.Data["book_id"].(string)
        chapterID := sagaLog.Context.Data["chapter_id"].(string)
        amount := sagaLog.Context.Data["amount"].(float64) * 0.7
        
        err := h.monetizationService.RecordRevenue(ctx, authorID, bookID, chapterID, amount)
        if err != nil {
            // 如果仍然失败，记录到死信队列
            log.Error("补偿失败，加入死信队列", "saga_id", sagaLog.ID, "step", step.StepName)
            // TODO: 加入死信队列
            return err
        }
        
        // 更新步骤状态
        step.Status = StepStatusCompleted
        h.sagaRepo.UpdateStep(ctx, sagaLog.ID, step)
    }
    
    return nil
}
```

---

## 五、Repository接口

```go
package saga

type SagaRepository interface {
    Create(ctx context.Context, saga *SagaLog) error
    GetByID(ctx context.Context, id string) (*SagaLog, error)
    Update(ctx context.Context, saga *SagaLog) error
    UpdateStep(ctx context.Context, sagaID string, step *SagaStep) error
    
    // 查询
    FindByStatus(ctx context.Context, status SagaStatus, limit int) ([]*SagaLog, error)
    FindPendingCompensation(ctx context.Context, limit int) ([]*SagaLog, error)
    
    // 健康检查
    Health(ctx context.Context) error
}
```

---

## 六、实施建议

### 6.1 开发阶段

| 阶段 | 时间 | 任务 |
|------|------|------|
| **阶段1** | 2天 | 数据模型、Repository实现 |
| **阶段2** | 3天 | 购买章节Saga实现 |
| **阶段3** | 3天 | 提现Saga实现 |
| **阶段4** | 2天 | 补偿机制、消息队列 |
| **阶段5** | 2天 | 测试与优化 |

### 6.2 技术风险

| 风险 | 等级 | 应对措施 |
|------|------|---------|
| 补偿失败 | 高 | 死信队列、人工介入 |
| 性能问题 | 中 | 异步补偿、批量处理 |
| 数据不一致 | 高 | 完善的日志、监控告警 |

---

**文档版本**: v1.0  
**创建时间**: 2025-10-21  
**负责人**: 架构组  
**审核状态**: 待评审

