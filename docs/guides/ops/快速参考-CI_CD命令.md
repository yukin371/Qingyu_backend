# CI/CD 快速参考

## 📋 命令速查表

### 开发环境设置

```bash
# 初始化项目（首次使用）
make init

# 安装开发工具
make install-tools

# 安装 golangci-lint
make install-lint
```

### 日常开发

```bash
# 启动开发服务器
make run

# 启动热重载模式
make dev

# 构建应用
make build

# 清理构建文件
make clean
```

### 代码质量

```bash
# 代码格式化
make fmt

# 整理导入
make imports

# 运行 linter
make lint

# Go vet 检查
make vet

# 快速检查（格式+vet+lint+单元测试）
make check
```

### 安全和质量

```bash
# 安全扫描（gosec）
make security

# 依赖漏洞检查（govulncheck）
make vuln-check

# 代码复杂度检查
make complexity
```

### 测试

```bash
# 运行所有测试
make test

# 单元测试
make test-unit

# 集成测试
make test-integration

# API 测试
make test-api

# 快速测试（跳过慢速）
make test-quick

# 测试覆盖率
make test-coverage

# 检查覆盖率是否达标（>=60%）
make test-coverage-check

# 生成详细测试报告
make test-report

# 基准测试
make test-bench
```

### CI/CD

```bash
# 完整 CI 流程（推荐）
make ci

# 本地模拟 GitHub Actions
make ci-local

# PR 提交前检查
make pr-check
```

### 依赖管理

```bash
# 下载依赖
make deps

# 更新依赖
make deps-update

# 创建 vendor 目录
make deps-vendor
```

### Docker

```bash
# 构建 Docker 镜像
make docker-build

# 运行 Docker 容器
make docker-run

# 清理 Docker 镜像
make docker-clean
```

## 🔄 工作流程

### 开始新功能

```bash
# 1. 创建新分支
git checkout -b feat/your-feature-name

# 2. 开发...

# 3. 提交前检查
make pr-check

# 4. 提交代码
git add .
git commit -m "feat: your feature description"

# 5. 推送到远程
git push origin feat/your-feature-name

# 6. 创建 Pull Request
```

### 修复 Bug

```bash
# 1. 创建修复分支
git checkout -b fix/bug-description

# 2. 修复...

# 3. 运行测试
make test

# 4. 提交
git commit -m "fix: bug description"

# 5. 推送并创建 PR
git push origin fix/bug-description
```

### 发布新版本

```bash
# 1. 确保在 main 分支
git checkout main
git pull

# 2. 运行完整检查
make ci-local

# 3. 创建 tag
git tag -a v1.0.0 -m "Release v1.0.0"

# 4. 推送 tag
git push origin v1.0.0

# GitHub Actions 会自动：
# - 运行测试
# - 构建多平台二进制
# - 创建 GitHub Release
```

## 🎯 常用场景

### 场景 1: 快速验证代码

```bash
make check
```

包含：格式化 → vet → lint → 单元测试

### 场景 2: 提交 PR 前的完整检查

```bash
make pr-check
```

包含：格式化 → 导入整理 → lint → 测试 → 覆盖率 → 依赖验证

### 场景 3: 本地模拟 CI 流程

```bash
make ci-local
```

包含：格式 → vet → lint → 安全扫描 → 漏洞检查 → 测试 → 覆盖率

### 场景 4: 只运行受影响的测试

```bash
# 运行单元测试（不需要 MongoDB）
make test-unit

# 运行 API 测试
make test-api
```

## 📊 测试覆盖率

### 查看覆盖率

```bash
# 生成 HTML 报告
make test-coverage

# 在浏览器中打开 coverage.html
```

### 检查是否达标

```bash
# 要求 >=60%
make test-coverage-check

# 要求 >=80%（修改 Makefile 中的阈值）
make test-coverage-check
```

## 🔍 调试技巧

### 运行特定测试

```bash
# 运行特定包的测试
go test -v ./service/user/...

# 运行特定测试函数
go test -v -run TestCreateUser ./service/user/...

# 运行匹配模式的测试
go test -v -run "Test.*User" ./...
```

### 查看详细输出

```bash
# 详细测试输出
make test-verbose

# 查看 linter 的详细信息
golangci-lint run --verbose
```

### 本地运行 GitHub Actions

```bash
# 安装 act
brew install act  # macOS
# 或
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# 运行所有工作流
act

# 运行特定工作流
act -W .github/workflows/ci.yml

# 运行特定 job
act -j unit-tests
```

## 🚨 故障排查

### Linter 失败

```bash
# 自动修复可修复的问题
golangci-lint run --fix

# 查看具体错误
golangci-lint run --verbose
```

### 测试失败

```bash
# 清理测试缓存
make test-clean

# 重新运行测试
make test

# 只运行失败的测试
make test-fix
```

### 依赖问题

```bash
# 清理并重新下载
go clean -modcache
go mod download

# 验证依赖
go mod verify

# 整理依赖
go mod tidy
```

## 📝 提交信息规范

### 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响代码运行）
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试相关
- `build`: 构建系统或外部依赖
- `ci`: CI 配置文件和脚本
- `chore`: 其他改动

### 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### 示例

```bash
# 简单提交
git commit -m "feat: 添加用户注册功能"

# 带作用域
git commit -m "fix(auth): 修复 JWT token 过期问题"

# 带详细说明
git commit -m "refactor(service): 重构用户服务

- 使用依赖注入
- 改进错误处理
- 添加单元测试

Closes #123"
```

## 🔗 相关链接

### 文档

- [CI/CD 配置指南](CI_CD配置指南.md)
- [CI/CD 问题解决方案](CI_CD问题解决方案.md)
- [GitHub Actions 工作流](.github/workflows/README.md)
- [测试指南](../testing/测试指南.md)

### 工具

- [golangci-lint](https://golangci-lint.run/)
- [gosec](https://github.com/securego/gosec)
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
- [act - 本地运行 GitHub Actions](https://github.com/nektos/act)

### 规范

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## 💡 小贴士

1. **提交前总是运行** `make pr-check`
2. **保持测试覆盖率** >=60%（推荐 80%）
3. **遵循提交规范** 便于自动生成 changelog
4. **频繁提交** 保持提交小而专注
5. **本地测试** 在推送前本地运行所有检查
6. **查看 CI 日志** 如果 CI 失败，仔细查看日志
7. **更新文档** 代码变更时同步更新文档
8. **安全第一** 定期运行安全扫描和漏洞检查

---

**最后更新**: 2025-10-20  
**维护者**: 青羽后端团队

