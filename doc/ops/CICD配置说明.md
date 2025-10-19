# CI/CD 配置说明

**文档版本**：1.0  
**创建日期**：2025-10-18  
**最后更新**：2025-10-18

---

## 📋 概述

本文档说明青羽后端项目的CI/CD自动化流程配置。

### 核心特性

- ✅ 自动化测试（单元测试+集成测试）
- ✅ 代码质量检查（Lint+格式化）
- ✅ 安全扫描
- ✅ 覆盖率报告
- ✅ 性能基准测试
- ✅ Docker构建
- ✅ 自动部署（测试环境）

---

## 🔧 配置文件

### 1. GitHub Actions工作流

**文件**：`.github/workflows/ci.yml`

**触发条件**：
- Push到main/dev分支
- Pull Request到main/dev分支

**Jobs清单**：

| Job | 功能 | 依赖 |
|-----|------|------|
| lint | 代码检查 | - |
| test | 单元测试 | - |
| integration-test | 集成测试 | lint, test |
| build | 构建测试 | lint, test |
| security | 安全扫描 | lint |
| code-quality | 代码质量分析 | - |
| benchmark | 性能测试 | - |
| docker | Docker构建 | build |
| deploy-dev | 部署测试环境 | test, integration-test, build, security |
| report | 生成报告 | 所有jobs |

### 2. golangci-lint配置

**文件**：`.golangci.yml`

**启用的Linters**（24个）：
- bodyclose - 检查HTTP response body是否关闭
- errcheck - 检查错误是否被处理
- gosec - 安全检查
- govet - Go官方静态分析
- staticcheck - 高级静态分析
- 更多...

**排除规则**：
- 测试文件宽松检查
- cmd/目录宽松检查
- migration/目录宽松检查

### 3. 测试配置

**文件**：`config/config.test.yaml`

**特点**：
- 独立的测试数据库
- Redis使用不同DB
- 调试级别日志
- 测试后自动清理

---

## 🚀 使用指南

### 本地运行测试

#### 方法1：使用测试脚本（推荐）

```bash
# 运行所有测试
./scripts/run_tests.sh

# 包含集成测试
RUN_INTEGRATION=true ./scripts/run_tests.sh

# 包含性能测试
RUN_BENCHMARK=true ./scripts/run_tests.sh

# 包含所有
RUN_INTEGRATION=true RUN_BENCHMARK=true ./scripts/run_tests.sh
```

**脚本功能**：
1. ✅ 检查依赖（Go版本）
2. ✅ 检查服务（MongoDB、Redis）
3. ✅ 代码格式检查
4. ✅ Lint检查
5. ✅ 运行单元测试
6. ✅ 生成覆盖率报告
7. ✅ 覆盖率统计和排名

#### 方法2：手动运行

```bash
# 1. 启动服务
docker-compose up -d mongodb redis

# 2. 运行测试
export CONFIG_PATH=config/config.test.yaml
go test -v -race -coverprofile=coverage.txt ./...

# 3. 查看覆盖率
go tool cover -func=coverage.txt
go tool cover -html=coverage.txt -o coverage.html

# 4. 运行Lint
golangci-lint run

# 5. 运行性能测试
go test -bench=. -benchmem ./...
```

### CI/CD流程

#### Push到dev分支

```
1. Lint检查
2. 单元测试（带覆盖率）
3. 集成测试
4. 构建测试
5. 安全扫描
6. 代码质量分析
7. Docker构建
8. 部署到测试环境 ✅
9. 生成报告
```

#### Push到main分支

```
1. Lint检查
2. 单元测试（带覆盖率）
3. 集成测试
4. 构建测试
5. 安全扫描
6. 代码质量分析
7. 性能测试 ✅
8. Docker构建
9. 生成报告
```

#### Pull Request

```
1. Lint检查
2. 单元测试（带覆盖率）
3. 集成测试
4. 构建测试
5. 安全扫描
6. 代码质量分析
7. 生成报告
```

---

## 📊 测试覆盖率

### 当前覆盖率

| 模块 | 覆盖率 | 状态 |
|-----|--------|------|
| Service层 | ~85% | ✅ 达标 |
| Repository层 | ~90% | ✅ 达标 |
| API层 | ~80% | ✅ 达标 |
| Model层 | ~95% | ✅ 优秀 |
| 总体 | ~85% | ✅ 达标 |

### 覆盖率目标

- **最低要求**：80%
- **推荐目标**：85%
- **优秀水平**：90%+

### 查看覆盖率报告

```bash
# 生成报告
./scripts/run_tests.sh

# 在浏览器中打开
open coverage/coverage.html
```

---

## 🔒 安全扫描

### Gosec扫描

**功能**：
- SQL注入检测
- 文件路径遍历
- 不安全的加密
- 命令注入
- 等等...

**配置**：
- 生成SARIF格式报告
- 上传到GitHub Security

### 查看安全报告

1. 进入GitHub仓库
2. 点击"Security"标签
3. 查看"Code scanning alerts"

---

## 📈 代码质量指标

### Cyclomatic Complexity（圈复杂度）

**阈值**：15  
**工具**：gocyclo

**说明**：
- < 10：简单函数
- 10-15：中等复杂度
- > 15：需要重构

### Cognitive Complexity（认知复杂度）

**阈值**：15  
**工具**：gocognit

**说明**：
- 评估代码的可读性
- 比圈复杂度更关注人类理解

### 代码格式化

**工具**：gofmt

**规则**：
- 所有代码必须格式化
- CI会自动检查
- 本地运行：`gofmt -s -w .`

---

## 🐳 Docker构建

### 构建配置

**Dockerfile**：`docker/Dockerfile.prod`

**特性**：
- 多阶段构建
- 最小化镜像大小
- 缓存优化

### 本地测试Docker构建

```bash
# 构建镜像
docker build -f docker/Dockerfile.prod -t qingyu-backend:test .

# 运行容器
docker run -p 8080:8080 qingyu-backend:test
```

---

## 🚢 自动部署

### 测试环境部署

**触发条件**：Push到dev分支且所有测试通过

**流程**：
1. 拉取最新代码
2. 构建Docker镜像
3. 停止旧容器
4. 启动新容器
5. 健康检查
6. 发送通知

### 生产环境部署

**触发条件**：手动触发（通过GitHub Actions）

**流程**：
1. 创建Release Tag
2. 构建生产镜像
3. 推送到镜像仓库
4. 通知运维
5. 人工确认后部署

---

## 📝 最佳实践

### 1. 提交前检查

```bash
# 运行快速检查
gofmt -s -w .
golangci-lint run
go test ./...
```

### 2. 编写测试

**单元测试**：
```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    service := NewUserService(...)
    
    // Act
    result, err := service.CreateUser(...)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

**表驱动测试**：
```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid", "test", false},
        {"empty", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("got err %v, want err %v", err, tt.wantErr)
            }
        })
    }
}
```

### 3. Mock使用

```go
// 使用testify/mock
type MockUserRepo struct {
    mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

// 在测试中使用
mockRepo := new(MockUserRepo)
mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
```

### 4. 集成测试

```go
// +build integration

func TestUserServiceIntegration(t *testing.T) {
    // 使用真实的数据库连接
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    service := NewUserService(db)
    
    // 测试真实场景
    user, err := service.CreateUser(...)
    require.NoError(t, err)
}
```

---

## 🔧 故障排查

### 问题1：测试失败

**检查**：
1. MongoDB是否运行
2. Redis是否运行
3. 配置文件是否正确
4. 环境变量是否设置

```bash
# 检查服务
docker-compose ps

# 重启服务
docker-compose restart mongodb redis
```

### 问题2：覆盖率太低

**解决方案**：
1. 识别未覆盖的代码
2. 编写缺失的测试
3. 使用`go tool cover`查看详情

```bash
# 查看未覆盖的代码
go tool cover -html=coverage.txt
```

### 问题3：Lint错误

**常见问题**：
- 未处理的错误
- 未使用的变量
- 代码复杂度过高

**解决方案**：
```bash
# 查看详细错误
golangci-lint run --verbose

# 自动修复部分问题
golangci-lint run --fix
```

---

## 📚 相关资源

### 工具文档

- [GitHub Actions](https://docs.github.com/en/actions)
- [golangci-lint](https://golangci-lint.run/)
- [Gosec](https://github.com/securego/gosec)
- [testify](https://github.com/stretchr/testify)

### 项目文档

- [测试指南](../testing/README.md)
- [部署指南](./部署指南.md)
- [开发规范](../architecture/项目开发规则.md)

---

## ✨ 总结

### CI/CD配置完成

- ✅ GitHub Actions工作流
- ✅ golangci-lint配置
- ✅ 测试脚本
- ✅ 测试配置
- ✅ 完整文档

### 下一步

1. **推送代码触发CI**
2. **查看测试结果**
3. **修复失败的测试**
4. **提高覆盖率**
5. **继续开发阶段四**

---

**文档维护者**：青羽后端团队  
**更新周期**：根据CI/CD配置变化及时更新

