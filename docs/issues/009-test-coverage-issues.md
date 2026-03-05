# Issue #009: 测试覆盖率不足

**优先级**: 高 (P0)
**类型**: 测试问题
**状态**: 待处理
**创建日期**: 2026-03-05
**来源报告**: [后端综合审计报告](../reports/archived/backend-comprehensive-audit-summary-2026-01-26.md)、[后端测试分析](../reports/archived/backend-testing-analysis-2026-01-26.md)

---

## 问题描述

多个核心模块测试覆盖率为 0%，存在严重的质量风险。

### 具体问题

#### 1. AI 服务测试覆盖率 0% 🔴 P0

**问题**: AI 模块没有任何测试。

**影响**:
- AI 功能质量无保障
- 重构风险高
- Bug 难以定位

**模块**:
- `service/ai/`
- `service/ai/chat/`
- `service/ai/generation/`

#### 2. 财务系统测试覆盖率 0% 🔴 P0

**问题**: 财务相关功能没有任何测试。

**影响**:
- **资金安全风险**
- 无法验证计算正确性
- 审计追踪困难

**模块**:
- `service/finance/`
- `service/payment/`
- `service/wallet/`

#### 3. 4 个测试编译失败 🔴 P0

**问题**: 部分测试代码无法编译。

**影响**:
- CI/CD 流程阻塞
- 测试质量数据不准确

#### 4. Service 层覆盖率 60-70% 🟡 P1

**目标**: ≥ 80%

**差距**: 10-20%

#### 5. API 层覆盖率 40-50% 🟡 P1

**目标**: ≥ 60%

**差距**: 10-20%

---

## 测试现状

### 当前覆盖率统计

| 模块 | 覆盖率 | 目标 | 状态 |
|------|--------|------|------|
| AI 服务 | 0% | 80% | 🔴 严重不足 |
| 财务系统 | 0% | 80% | 🔴 严重不足 |
| 推荐系统 | 0% | 70% | 🟡 待补充 |
| Service 层 | 60-70% | 80% | 🟡 需提升 |
| API 层 | 40-50% | 60% | 🟡 需提升 |
| Repository 层 | 75-85% | 80% | ✅ 基本达标 |

### 测试文件统计

| 指标 | 数值 |
|------|------|
| 测试文件 | 150 个 |
| 测试代码 | 68,012 行 |
| 编译失败 | 4 个 |

---

## 解决方案

### 1. 修复编译失败的测试

```bash
# 1. 找出编译失败的测试
go test ./... -v 2>&1 | grep "cannot load"

# 2. 修复导入问题
# 3. 修复类型不匹配
# 4. 修复依赖问题
```

### 2. AI 服务测试框架

```go
// service/ai/chat/chat_service_test.go
package chat_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock AI 客户端
type MockAIClient struct {
    mock.Mock
}

func (m *MockAIClient) GenerateChat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
    args := m.Called(ctx, req)
    return args.Get(0).(*ChatResponse), args.Error(1)
}

// 测试聊天服务
func TestChatService_GenerateChat(t *testing.T) {
    // 安排
    mockClient := new(MockAIClient)
    service := NewChatService(mockClient)

    mockClient.On("GenerateChat", mock.Anything, mock.Anything).
        Return(&ChatResponse{Content: "Test response"}, nil)

    // 执行
    resp, err := service.GenerateChat(context.Background(), &ChatRequest{
        Message: "Hello",
    })

    // 断言
    assert.NoError(t, err)
    assert.Equal(t, "Test response", resp.Content)
    mockClient.AssertExpectations(t)
}
```

### 3. 财务系统测试

```go
// service/finance/wallet_service_test.go
package finance_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestWalletService_Recharge(t *testing.T) {
    // 安排
    repo := NewMockWalletRepository()
    service := NewWalletService(repo)
    userID := "test-user"

    // 执行 - 充值 100 元
    balance, err := service.Recharge(context.Background(), userID, 10000)

    // 断言
    require.NoError(t, err)
    assert.Equal(t, int64(10000), balance)

    // 验证余额已更新
    wallet, _ := repo.GetWallet(context.Background(), userID)
    assert.Equal(t, int64(10000), wallet.Balance)
}

func TestWalletService_Withdraw_InsufficientBalance(t *testing.T) {
    // 安排
    repo := NewMockWalletRepository()
    service := NewWalletService(repo)
    userID := "test-user"

    // 执行 - 尝试提取超过余额的金额
    _, err := service.Withdraw(context.Background(), userID, 10000)

    // 断言
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "insufficient balance")
}
```

### 4. 表驱动测试模式

```go
// service/payment/payment_service_test.go
func TestPaymentService_ProcessPayment(t *testing.T) {
    tests := []struct {
        name        string
        amount      int64
        paymentType string
        wantErr     bool
        errMsg      string
    }{
        {
            name:        "valid payment",
            amount:      10000,
            paymentType: "alipay",
            wantErr:     false,
        },
        {
            name:        "invalid amount",
            amount:      -100,
            paymentType: "alipay",
            wantErr:     true,
            errMsg:      "invalid amount",
        },
        {
            name:        "unsupported payment type",
            amount:      10000,
            paymentType: "unknown",
            wantErr:     true,
            errMsg:      "unsupported payment type",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := setupTestService(t)

            err := service.ProcessPayment(
                context.Background(),
                tt.amount,
                tt.paymentType,
            )

            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

---

## 实施计划

### Phase 1: 修复编译问题（1 天）

1. 找出所有编译失败的测试
2. 修复导入和依赖问题
3. 验证所有测试可编译

### Phase 2: AI 服务测试（2 周）

**优先级**: 高

1. **Mock AI 服务客户端**
   - [ ] 定义 Mock 接口
   - [ ] 实现 Mock 客户端

2. **核心功能测试**
   - [ ] 聊天服务测试
   - [ ] 内容生成测试
   - [ ] 总结服务测试

3. **集成测试**
   - [ ] 端到端测试
   - [ ] 性能测试

### Phase 3: 财务系统测试（1 周）

**优先级**: 高（涉及资金安全）

1. **钱包服务测试**
   - [ ] 充值测试
   - [ ] 提现测试
   - [ ] 转账测试
   - [ ] 余额不足测试

2. **支付服务测试**
   - [ ] 支付流程测试
   - [ ] 回调处理测试
   - [ ] 异常场景测试

3. **财务计算测试**
   - [ ] 金额计算精度测试
   - [ ] 手续费计算测试
   - [ ] 税费计算测试

### Phase 4: Service 和 API 层提升（2 周）

1. **Service 层测试**
   - 识别覆盖率低的模块
   - 补充单元测试
   - 添加集成测试

2. **API 层测试**
   - 添加 HTTP 测试
   - 验证请求/响应格式
   - 测试错误场景

---

## 测试工具推荐

### 1. 测试框架

```go
// 安装测试依赖
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/stretchr/testify/suite
```

### 2. 测试辅助工具

```go
// test/mocks/repository.go
package mocks

type MockWalletRepository struct {
    mock.Mock
}

func (m *MockWalletRepository) GetBalance(ctx context.Context, userID string) (int64, error) {
    args := m.Called(ctx, userID)
    return args.Get(0).(int64), args.Error(1)
}

// 生成 mock 工具
// go install github.com/golang/mock/mockgen@latest
// mockgen -source=repository/wallet_repository.go -destination=test/mocks/wallet_repository_mock.go
```

### 3. 测试覆盖率工具

```bash
# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# 查看覆盖率
go tool cover -func=coverage.out | grep total
```

---

## 测试最佳实践

### 1. 测试命名

```go
// ✅ 好的命名
func TestWalletService_Recharge_Success(t *testing.T)
func TestWalletService_Recharge_InvalidAmount(t *testing.T)
func TestWalletService_Recharge_NegativeAmount(t *testing.T)

// ❌ 不好的命名
func TestWallet1(t *testing.T)
func TestWalletRechargeTest(t *testing.T)
```

### 2. 测试结构

```go
func TestXxx(t *testing.T) {
    // Arrange - 准备测试数据
    service := setupTestService(t)
    userID := "test-user"

    // Act - 执行被测试的函数
    result, err := service.DoSomething(userID)

    // Assert - 验证结果
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### 3. 测试隔离

```go
// 使用 t.Cleanup 确保清理
func TestXxx(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()  // 简单场景

    // 复杂场景使用 t.Cleanup
    t.Cleanup(func() {
        // 清理逻辑
        db.DropAllCollections()
    })
}
```

---

## 检查清单

### 编译问题
- [ ] 所有测试可编译
- [ ] 无导入错误
- [ ] 无类型错误

### AI 服务测试
- [ ] Mock 客户端实现
- [ ] 聊天服务测试
- [ ] 生成服务测试
- [ ] 集成测试

### 财务系统测试
- [ ] 钱包服务测试
- [ ] 支付服务测试
- [ ] 计算精度测试
- [ ] 异常场景测试

### 覆盖率提升
- [ ] Service 层达到 80%
- [ ] API 层达到 60%
- [ ] 覆盖率报告生成

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [后端测试分析](../reports/archived/backend-testing-analysis-2026-01-26.md) | 测试问题详细分析 |
| [Go 测试最佳实践](https://go.dev/doc/tutorial/add-a-test) | Go 官方测试教程 |
| [Testify 文档](https://github.com/stretchr/testify) | Testify 使用文档 |

---

## 相关 Issue

- [#003: 测试基础设施改进](./003-test-infrastructure-improvements.md)
- [#007: Service 层事务管理缺失](./007-transaction-management.md)
