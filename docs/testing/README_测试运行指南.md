# 测试运行指南

快速参考：如何运行青羽后端的各种测试

---

## 🚀 快速开始

### 运行所有测试
```bash
go test ./...
```

### 运行所有测试（详细输出）
```bash
go test -v ./...
```

### 运行测试并显示覆盖率
```bash
go test -cover ./...
```

---

## 📂 按目录运行

### 单元测试（源代码目录）
```bash
# AI服务单元测试
go test ./service/ai/...

# 项目服务单元测试
go test ./service/project/...

# 所有服务单元测试
go test ./service/...

# 中间件测试
go test ./middleware/...

# 工具包测试
go test ./pkg/...
```

### 集成测试
```bash
# 所有集成测试
go test ./test/integration/...

# 版本控制集成测试
go test -v -run TestUpdateContentWithVersion ./test/integration/

# 用户生命周期端到端测试
go test -v -run TestE2E_UserLifecycle ./test/integration/
```

### API测试
```bash
# 所有API测试
go test ./test/api/...

# 书店API测试
go test -v ./test/api/bookstore_api_test.go

# 阅读器API测试
go test -v ./test/api/reader_api_test.go
```

### 性能测试
```bash
# 运行所有性能基准测试
go test -bench=. ./test/performance/...

# 运行特定基准测试
go test -bench=BenchmarkBookstore ./test/performance/...

# 性能测试 + CPU分析
go test -bench=. -cpuprofile=cpu.prof ./test/performance/...

# 性能测试 + 内存分析
go test -bench=. -memprofile=mem.prof ./test/performance/...
```

---

## 🎯 按测试类型运行

### 运行特定测试函数
```bash
# 运行名称匹配的测试
go test -run TestChatService_StartChat ./service/ai/

# 运行多个匹配的测试（正则表达式）
go test -run "TestChat.*" ./service/ai/

# 运行集成测试中的特定测试
go test -run TestVersionService ./test/integration/
```

### 跳过慢速测试
```bash
# 使用 -short 标志跳过长时间运行的测试
go test -short ./...

# 在测试中使用
if testing.Short() {
    t.Skip("跳过集成测试")
}
```

### 并行运行测试
```bash
# 指定并发数
go test -parallel 4 ./...
```

---

## 📊 测试覆盖率

### 生成覆盖率报告
```bash
# 简单覆盖率
go test -cover ./...

# 详细覆盖率（按包）
go test -coverprofile=coverage.out ./...

# 查看覆盖率报告
go tool cover -func=coverage.out

# 生成HTML覆盖率报告
go tool cover -html=coverage.out -o coverage.html
```

### 查看特定包的覆盖率
```bash
# AI服务覆盖率
go test -cover ./service/ai/...

# 集成测试覆盖率
go test -cover ./test/integration/...
```

---

## 🔍 调试测试

### 详细输出
```bash
# 显示所有测试输出
go test -v ./...

# 显示测试日志（即使测试通过）
go test -v -args -test.v
```

### 运行失败的测试
```bash
# 第一个失败后停止
go test -failfast ./...

# 显示完整的错误堆栈
go test -v ./... 2>&1 | more
```

### 超时控制
```bash
# 设置测试超时（默认10分钟）
go test -timeout 30s ./...

# 单个测试的超时
go test -timeout 5m ./test/integration/...
```

---

## 🐳 Docker环境测试

### 使用Docker运行MongoDB集成测试
```bash
# 启动测试用MongoDB
docker-compose -f docker-compose.test.yml up -d

# 运行集成测试
go test ./test/integration/...

# 清理
docker-compose -f docker-compose.test.yml down
```

### Repository层测试脚本
```bash
# Windows
cd test/repository/user
./run_docker_test.ps1

# Linux/Mac
cd test/repository/user
./run_docker_test.sh
```

---

## 🔧 Mock测试

### 运行使用Mock的测试
```bash
# AI聊天服务Mock测试
go test -v ./service/ai/chat_service_test.go

# 确保Mock被正确调用
go test -v -run TestChatService_StartChat ./service/ai/
```

---

## 📈 持续集成（CI）

### GitHub Actions测试命令
```yaml
# .github/workflows/test.yml
- name: Run Unit Tests
  run: go test -v -cover ./service/... ./pkg/... ./middleware/...

- name: Run Integration Tests
  run: go test -v -cover ./test/integration/...

- name: Run API Tests
  run: go test -v -cover ./test/api/...
```

---

## 💡 常用测试命令组合

### 开发中快速测试
```bash
# 测试当前修改的包（快速反馈）
go test -v ./service/ai/

# 测试并显示覆盖率
go test -v -cover ./service/ai/
```

### 提交前完整测试
```bash
# 运行所有测试 + 覆盖率
go test -v -cover ./...

# 检查代码格式
go fmt ./...

# 运行Linter
golangci-lint run
```

### 性能分析
```bash
# 生成性能分析文件
go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof ./test/performance/...

# 分析CPU性能
go tool pprof cpu.prof

# 分析内存使用
go tool pprof mem.prof
```

---

## 🐛 测试失败排查

### 查看详细错误
```bash
# 显示完整的测试输出
go test -v -run TestFailingTest ./...

# 显示panic堆栈
go test -v ./... 2>&1 | grep -A 10 "panic"
```

### 隔离问题
```bash
# 只运行失败的测试
go test -run TestSpecificFailure ./service/ai/

# 多次运行以检测不稳定测试
for i in {1..10}; do go test -run TestFlaky ./...; done
```

---

## 📝 测试命名约定

### 测试文件命名
```
单元测试：      xxx_test.go
集成测试：      xxx_integration_test.go
端到端测试：    xxx_e2e_test.go
性能测试：      xxx_benchmark_test.go
```

### 测试函数命名
```
单元测试：      Test<ServiceName>_<MethodName>_<Scenario>
集成测试：      Test<Feature>_Integration
端到端测试：    TestE2E_<Feature>
性能测试：      Benchmark<Feature>
```

---

## 🎓 推荐测试流程

### 1. 本地开发
```bash
# 快速测试当前工作的包
go test -v ./service/ai/

# 确认无破坏性变更
go test ./...
```

### 2. 提交前
```bash
# 完整测试套件
go test -v -cover ./...

# 检查覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
```

### 3. Pull Request前
```bash
# 运行所有测试（包括集成测试）
go test -v ./...

# 运行Linter
golangci-lint run

# 检查代码格式
go fmt ./...
go vet ./...
```

---

## 📚 更多资源

- [测试组织规范](../doc/testing/测试组织规范.md) - 详细的测试分类和最佳实践
- [单元测试示例](./examples/service_test_example.go) - 单元测试示例代码
- [集成测试示例](./integration/README.md) - 集成测试指南
- [性能测试指南](./performance/README.md) - 性能基准测试

---

**最后更新**: 2025-10-17  
**维护团队**: 青羽后端团队

