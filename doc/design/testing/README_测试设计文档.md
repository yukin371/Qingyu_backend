# 测试设计文档

本目录包含青羽写作平台的测试相关设计文档，涵盖自动化测试、性能测试等内容。

## 📁 文档目录

### 自动化测试
- [自动化测试总体方案](./自动化测试总体方案.md) - 单元测试、集成测试、端到端测试的整体方案

### 性能测试
- [性能测试模型与k6脚本](./性能测试模型与k6脚本.md) - 性能测试场景、指标、k6脚本设计

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

## 🔧 测试工具链

### 单元测试
```go
// 使用 testify 进行断言
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    service := NewUserService(mockRepo)
    
    // Act
    err := service.CreateUser(user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user.ID)
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

- 2025-01-01: 创建测试设计文档目录，整理现有设计文档
