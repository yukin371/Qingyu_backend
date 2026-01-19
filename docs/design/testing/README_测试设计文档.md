# 测试设计文档

本目录包含青羽写作平台的测试相关设计文档，涵盖自动化测试、性能测试等内容。

## 📁 文档目录

### 自动化测试
- [自动化测试总体方案](./自动化测试总体方案.md) ✅ - 单元测试、集成测试、端到端测试的整体方案
- [自动化测试总体方案（简化版）](./自动化测试总体方案-简化版.md) - 快速上手的简化方案
- [测试最佳实践指南](./测试最佳实践指南.md) 🆕 - Table-Driven测试、Mock使用、测试模式

### 性能测试
- [性能测试模型与k6脚本](./性能测试模型与k6脚本.md) - 性能测试场景、指标、k6脚本设计
- [性能测试模型与k6脚本（简化版）](./性能测试模型与k6脚本-简化版.md) - 简化的性能测试方案

### 模块测试设计
- [书城系统测试设计](./书城系统测试设计.md) - 书城模块的详细测试设计

### 测试示例
- [Service层测试示例](../../test/examples/service_test_example.go) 🆕 - 完整的Service测试示例
- [Repository层测试示例](../../test/examples/repository_test_example.go) 🆕 - 完整的Repository测试示例

## 🎯 测试目标

### 质量保证
- 代码覆盖率 ≥ 80%
- 核心功能覆盖率 100%
- 缺陷逃逸率 < 5%
- 回归测试自动化率 ≥ 90%

### 性能保证
- API响应时间 < 200ms (P95)
- 并发支持 ≥ 1000 QPS
- 系统可用性 ≥ 99.9%
- 数据库查询时间 < 100ms

## 🧪 测试类型

### 单元测试 (Unit Test)
- **工具**: Go testing, testify
- **范围**: 函数、方法级别
- **目标**: 验证单个单元的正确性
- **覆盖**: Service层、Repository层、Utils

### 集成测试 (Integration Test)
- **工具**: Go testing, Docker Compose
- **范围**: 模块间交互
- **目标**: 验证模块集成的正确性
- **覆盖**: API + Service + Repository

### 端到端测试 (E2E Test)
- **工具**: Playwright, Cypress
- **范围**: 完整业务流程
- **目标**: 验证用户场景的正确性
- **覆盖**: 关键业务流程

### 性能测试 (Performance Test)
- **工具**: k6, JMeter
- **类型**: 负载测试、压力测试、稳定性测试
- **目标**: 验证系统性能指标
- **场景**: 高并发、大数据量

### 安全测试 (Security Test)
- **工具**: OWASP ZAP, SonarQube
- **类型**: 漏洞扫描、渗透测试
- **目标**: 发现安全隐患
- **范围**: 全系统

## 📊 测试策略

### 测试金字塔
```
      ┌─────────┐
      │  E2E    │ 10%
      ├─────────┤
      │ 集成测试 │ 30%
      ├─────────┤
      │ 单元测试 │ 60%
      └─────────┘
```

### 测试左移
- 需求阶段：可测性分析
- 开发阶段：TDD开发、单元测试
- 集成阶段：集成测试、API测试
- 发布阶段：E2E测试、性能测试

## 🚀 快速开始

### 第一次使用

```bash
# 1. 安装测试工具
bash scripts/setup-test-env.sh

# 2. 启动测试服务（Docker）
docker run -d -p 27017:27017 --name test-mongo mongo:5.0
docker run -d -p 6379:6379 --name test-redis redis:6.2-alpine

# 3. 运行测试
make test-unit           # 运行单元测试
make test-coverage       # 生成覆盖率报告
```

### 日常开发

```bash
# 运行所有测试
make test

# 只运行单元测试（快速）
make test-unit

# 生成测试模板
make test-gen file=service/user/user_service.go

# 查看覆盖率
make test-coverage

# 清理测试缓存
make test-clean
```

### Makefile命令总览

| 命令 | 说明 |
|------|------|
| `make test` | 运行所有测试（带竞态检测） |
| `make test-unit` | 运行单元测试（Service和Repository层） |
| `make test-integration` | 运行集成测试 |
| `make test-api` | 运行API测试 |
| `make test-coverage` | 生成覆盖率报告（HTML） |
| `make test-coverage-check` | 检查覆盖率是否达到80% |
| `make test-gen file=xxx.go` | 为指定文件生成测试模板 |
| `make test-clean` | 清理测试缓存和覆盖率文件 |
| `make test-watch` | 监视文件变化并自动运行测试 |

## 🔧 测试工具链

### 核心框架
- **Go testing** - 标准库测试框架
- **testify** - 断言和Mock框架 (`github.com/stretchr/testify`)
- **gotests** - 测试代码生成工具

### 测试辅助库
- **testutil** (`test/testutil/helpers.go`) - 测试助手函数
- **fixtures** (`test/fixtures/factory.go`) - 测试数据工厂

### CI/CD工具
- **GitHub Actions** - 自动化测试流程
- **Docker** - 测试环境服务（MongoDB, Redis）

### 单元测试示例
```go
// 使用 testify 进行断言
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    service := NewUserService(mockRepo)
    user := testutil.CreateTestUser()
    
    // Act
    err := service.CreateUser(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)
    mockRepo.AssertExpectations(t)
}
```

### Table-Driven测试示例
```go
func TestUserService_CreateUser_TableDriven(t *testing.T) {
    tests := []struct {
        name    string
        input   *CreateUserRequest
        mock    func(*MockUserRepository)
        wantErr bool
    }{
        {
            name: "成功创建用户",
            input: &CreateUserRequest{Username: "test"},
            mock: func(m *MockUserRepository) {
                m.On("Create", mock.Anything, mock.Anything).Return(nil)
            },
            wantErr: false,
        },
        // ... 更多测试用例
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ... 测试逻辑
        })
    }
}
```

### 性能测试
```javascript
// k6 性能测试脚本
import http from 'k6/http';
import { check } from 'k6';

export let options = {
    stages: [
        { duration: '1m', target: 100 },
        { duration: '3m', target: 500 },
        { duration: '1m', target: 0 },
    ],
};

export default function() {
    let res = http.get('https://api.example.com/books');
    check(res, {
        'status is 200': (r) => r.status === 200,
        'response time < 200ms': (r) => r.timings.duration < 200,
    });
}
```

## 📋 测试流程

### 开发阶段
1. 编写单元测试
2. 运行本地测试
3. 代码覆盖率检查
4. 提交代码

### CI阶段
1. 自动运行单元测试
2. 自动运行集成测试
3. 代码质量扫描
4. 构建Docker镜像

### 发布阶段
1. 冒烟测试
2. 回归测试
3. 性能测试
4. 安全测试

## 📈 测试指标

### 质量指标
- 代码覆盖率
- 测试通过率
- 缺陷密度
- 缺陷修复时间

### 性能指标
- 响应时间 (P50, P95, P99)
- 吞吐量 (QPS, TPS)
- 并发用户数
- 错误率

## 🔗 相关文档

- [书城系统测试设计](../测试/书城系统测试设计.md) - 书城模块的测试设计
- [自动化测试总体方案](./自动化测试总体方案.md) - 详细的自动化测试方案

## 📝 更新日志

- 2025-10-17: 完善测试工具链，添加快速开始指南、Makefile命令、测试示例
- 2025-01-01: 创建测试设计文档目录，整理现有设计文档
