# AI辅助功能测试套件

本目录包含青羽写作平台AI辅助功能的完整单元测试套件。

## 📁 测试文件结构

```
service/ai/
├── mocks/
│   └── ai_adapter_mock.go           # AI适配器Mock实现
├── summarize_service_test.go         # 内容总结服务测试
├── proofread_service_test.go         # 文本校对服务测试
├── sensitive_words_service_test.go   # 敏感词检测服务测试

api/v1/ai/
└── writing_assistant_api_test.go     # API层测试

根目录/
├── run_ai_writing_tests.bat          # Windows测试运行脚本
├── AI_WRITING_ASSISTANT_TESTS.md     # 测试文档
└── AI_TESTS_REPORT.md                # 测试报告
```

## 🚀 快速开始

### 1. 运行所有测试
```bash
# Windows
run_ai_writing_tests.bat

# 或手动运行
cd D:\Github\青羽\Qingyu_backend
go test -v ./service/ai/... ./api/v1/ai/...
```

### 2. 运行特定服务测试
```bash
# 内容总结服务
go test -v ./service/ai -run "TestSummarizeService.*"

# 文本校对服务
go test -v ./service/ai -run "TestProofreadService.*"

# 敏感词检测服务
go test -v ./service/ai -run "TestSensitiveWordsService.*"

# API层测试
go test -v ./api/v1/ai -run "TestWritingAssistantApi.*"
```

### 3. 运行性能测试
```bash
go test -bench=. -benchmem ./service/ai
```

### 4. 生成测试覆盖率报告
```bash
go test -coverprofile=coverage.out ./service/ai/...
go tool cover -html=coverage.out -o coverage.html
```

## 📋 测试覆盖范围

### ✅ 内容总结服务 (SummarizeService)
- ✅ 成功总结内容
- ✅ 空内容验证
- ✅ 不同总结类型（brief, detailed, keypoints）
- ✅ 关键点提取
- ✅ 压缩率计算
- ✅ 章节总结

### ✅ 文本校对服务 (ProofreadService)
- ✅ 成功校对内容
- ✅ 检查类型配置（spelling, grammar, punctuation, style）
- ✅ JSON/文本结果解析
- ✅ 统计信息生成
- ✅ 评分计算（0-100分）
- ✅ 长文本处理

### ✅ 敏感词检测服务 (SensitiveWordsService)
- ✅ 成功检测敏感词
- ✅ 分类检测（political, violence, adult）
- ✅ 自定义词库管理
- ✅ 词位置查找
- ✅ 上下文提取
- ✅ 风险级别评估
- ✅ 检测摘要生成

### ✅ API层测试 (WritingAssistantApi)
- ✅ 所有端点的成功场景
- ✅ 参数验证
- ✅ 错误处理
- ✅ HTTP方法验证
- ✅ 响应头验证

## 🎯 测试场景

### AI响应场景
- ✅ AI正常响应
- ⏳ AI服务超时（需要AdapterManager集成）
- ⏳ AI服务错误（需要AdapterManager集成）

### 数据验证场景
- ✅ 空内容处理
- ✅ 空白字符处理
- ✅ 参数类型验证
- ✅ 必需字段验证

### 业务逻辑场景
- ✅ 敏感词命中（政治、暴力、成人）
- ✅ 敏感词未命中
- ✅ 自定义敏感词
- ✅ 评分和统计

### 性能测试
- ✅ 关键点提取性能
- ✅ 评分计算性能
- ✅ 敏感词查找性能
- ✅ API请求性能

## 🔧 Mock适配器使用

### 创建Mock适配器
```go
import "Qingyu_backend/service/ai/mocks"

// 创建Mock实例
mockAdapter := mocks.NewMockAIAdapter("test-adapter")

// 配置成功响应
mockAdapter.SetTextResponse("这是AI生成的文本", 100)

// 配置失败
mockAdapter.ShouldFail = true
mockAdapter.FailureError = &adapter.AdapterError{
    Type:    "service_unavailable",
    Message: "服务不可用",
}

// 配置超时
mockAdapter.ResponseDelay = 5 * time.Second

// 重置状态
mockAdapter.Reset()
```

### 在测试中使用Mock
```go
func TestExample(t *testing.T) {
    // 1. 创建Mock
    mock := mocks.NewMockAIAdapter("test")

    // 2. 配置响应
    mock.SetTextResponse("预期响应", 50)

    // 3. 执行测试
    result, err := service.SomeMethod(ctx, req)

    // 4. 验证结果
    assert.NoError(t, err)
    assert.Equal(t, "预期响应", result.Text)

    // 5. 验证调用
    assert.Equal(t, 1, mock.CallCount)
}
```

## 📊 测试结果示例

### 成功输出
```
=== RUN   TestSummarizeService_SummarizeContent_EmptyContent
--- PASS: TestSummarizeService_SummarizeContent_EmptyContent (0.00s)
    assert.go:123: Error message not empty
PASS
coverage: 75.0% of statements
```

### 性能测试输出
```
BenchmarkSummarizeService_ExtractKeyPoints-8    500000    3.2 ns/op    128 B/op    2 allocs/op
BenchmarkProofreadService_CalculateScore-8      300000    4.5 ns/op     64 B/op    1 allocs/op
```

## 🐛 调试测试

### 启用详细输出
```bash
go test -v ./service/ai
```

### 只运行失败的测试
```bash
go test -v ./service/ai -run TestFailed
```

### 停在第一个失败
```bash
go test -v ./service/ai -failfast
```

## ⚠️ 已知限制

1. **AI适配器集成**: 当前Mock适配器独立，需要完善AdapterManager的依赖注入
2. **数据库操作**: 使用模拟数据，未测试实际数据库交互
3. **缓存系统**: 缓存功能待实现后需要添加测试
4. **配额系统**: 配额管理待实现后需要添加测试

## 🔄 持续集成

### GitHub Actions示例
```yaml
name: AI Writing Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Run tests
        run: |
          go test -v -race -coverprofile=coverage.out ./service/ai/...
          go tool cover -func=coverage.out
```

## 📚 相关文档

- [测试详细文档](../AI_WRITING_ASSISTANT_TESTS.md)
- [测试报告](../AI_TESTS_REPORT.md)
- [Go测试指南](https://golang.org/pkg/testing/)
- [Testify文档](https://github.com/stretchr/testify)

## 🤝 贡献指南

### 添加新测试
1. 在对应的`_test.go`文件中添加测试函数
2. 使用清晰的命名: `Test{Service}_{Method}_{Scenario}`
3. 添加必要的注释
4. 更新本README

### 测试命名规范
```go
// ✅ 好的命名
func TestSummarizeService_SummarizeContent_EmptyContent(t *testing.T)
func TestProofreadService_CalculateScore_PerfectContent(t *testing.T)

// ❌ 不好的命名
func TestSummarize1(t *testing.T)
func TestProof(t *testing.T)
```

### 断言使用
```go
// ✅ 使用testify断言
assert.NoError(t, err)
assert.Equal(t, expected, actual)
assert.True(t, condition)
assert.Contains(t, str, substring)

// ❌ 避免使用原生断言
if err != nil {
    t.Fatal(err)
}
```

## 📞 获取帮助

如有问题，请：
1. 查看[测试文档](../AI_WRITING_ASSISTANT_TESTS.md)
2. 查看[测试报告](../AI_TESTS_REPORT.md)
3. 检查测试代码中的注释
4. 提交Issue

## ✨ 致谢

感谢使用青羽写作平台的测试套件！

---

**最后更新**: 2026-01-03
**版本**: 1.0.0
**维护者**: 青羽开发团队
