# GitHub Actions 工作流说明

## 概述

本目录包含项目的所有 GitHub Actions 工作流配置。这些工作流提供了完整的 CI/CD 支持，包括代码检查、测试、构建和部署。

## 工作流列表

### 🔄 CI Pipeline (`ci.yml`)

**触发条件：**
- 推送到 `main`、`dev`、`develop` 分支
- 创建或更新 Pull Request

**包含任务：**
1. **代码检查（lint）** - golangci-lint 代码质量检查
2. **安全扫描（security）** - gosec 安全漏洞扫描
3. **单元测试（unit-tests）** - 快速单元测试，不依赖外部服务
4. **集成测试（integration-tests）** - 使用 MongoDB 的完整集成测试
5. **API 测试（api-tests）** - API 端点测试
6. **构建测试（build）** - 多平台交叉编译
7. **依赖检查（dependency-check）** - govulncheck 漏洞检查

**运行时间：** 约 8-12 分钟

### 🐳 Docker Build (`docker-build.yml`)

**触发条件：**
- 推送到 `main`、`dev` 分支
- 创建 tag（`v*`）
- PR 到 `main` 分支（仅构建，不推送）

**功能：**
- 构建多架构 Docker 镜像（amd64, arm64）
- 推送到 GitHub Container Registry (ghcr.io)
- Trivy 安全扫描
- 使用 GitHub Actions 缓存优化构建速度

**运行时间：** 约 5-8 分钟

### ✅ PR Check (`pr-check.yml`)

**触发条件：**
- Pull Request 打开、同步或重新打开时

**检查项：**
1. **PR 验证**
   - PR 标题格式（Conventional Commits）
   - 大文件检查（>5MB）
   - 敏感数据扫描（TruffleHog）

2. **代码质量**
   - 代码复杂度检查
   - 测试覆盖率验证（>=60%）

3. **变更检测**
   - 智能检测变更文件
   - 只测试受影响的模块

4. **自动标签**
   - 根据文件变更自动添加标签

**运行时间：** 约 5-7 分钟

### 🚀 Release (`release.yml`)

**触发条件：**
- 推送 Git tag（格式：`v*.*.*`）

**流程：**
1. 运行完整测试
2. 构建多平台二进制文件：
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)
3. 生成 SHA256 校验和
4. 创建 GitHub Release
5. 自动生成 Release Notes

**运行时间：** 约 10-15 分钟

### 🔒 CodeQL Analysis (`codeql.yml`)

**触发条件：**
- 推送到 `main`、`dev` 分支
- PR 到 `main` 分支
- 每周一定时运行

**功能：**
- 自动化代码安全分析
- 检测潜在的安全漏洞
- 结果上传到 GitHub Security

**运行时间：** 约 5-10 分钟

## 环境变量和 Secrets

### 环境变量

所有工作流中使用的环境变量：

```yaml
GO_VERSION: '1.21'           # Go 版本
MONGODB_VERSION: '6.0'       # MongoDB 版本
REGISTRY: ghcr.io            # 容器注册表
```

### 需要配置的 Secrets

在 GitHub 仓库设置中配置（Settings → Secrets and variables → Actions）：

| Secret 名称 | 说明 | 必需 |
|------------|------|------|
| `GITHUB_TOKEN` | 自动提供，无需配置 | ✅ |
| `CODECOV_TOKEN` | Codecov 上传 token | ⭕ 可选 |

### 环境特定配置

工作流使用以下环境变量（在 CI 中自动设置）：

```bash
MONGODB_URI=mongodb://admin:password@localhost:27017
MONGODB_DATABASE=qingyu_test
ENVIRONMENT=test
```

## 状态徽章

在 README 中添加状态徽章：

```markdown
[![CI](https://github.com/yourusername/Qingyu_backend/workflows/CI%20Pipeline/badge.svg)](https://github.com/yourusername/Qingyu_backend/actions/workflows/ci.yml)
[![Docker](https://github.com/yourusername/Qingyu_backend/workflows/Docker%20Build%20and%20Push/badge.svg)](https://github.com/yourusername/Qingyu_backend/actions/workflows/docker-build.yml)
[![CodeQL](https://github.com/yourusername/Qingyu_backend/workflows/CodeQL%20Analysis/badge.svg)](https://github.com/yourusername/Qingyu_backend/actions/workflows/codeql.yml)
```

## 本地调试

### 模拟 CI 检查

```bash
# 完整的 CI 流程
make ci-local

# PR 检查
make pr-check

# 单独检查
make lint           # 代码检查
make security       # 安全扫描
make vuln-check     # 漏洞检查
make test           # 运行测试
```

### 使用 act 本地运行 GitHub Actions

安装 [act](https://github.com/nektos/act)：

```bash
# macOS
brew install act

# Linux
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Windows
choco install act-cli
```

运行工作流：

```bash
# 运行所有工作流
act

# 运行特定工作流
act -W .github/workflows/ci.yml

# 运行特定 job
act -j unit-tests

# 使用特定事件触发
act pull_request
```

## 工作流优化

### 缓存策略

所有工作流都使用以下缓存：

1. **Go modules 缓存**
   ```yaml
   - uses: actions/setup-go@v5
     with:
       cache: true
   ```

2. **Docker layer 缓存**
   ```yaml
   cache-from: type=gha
   cache-to: type=gha,mode=max
   ```

### 并行化

- 多个 job 并行运行
- 使用矩阵策略构建多平台
- 智能的依赖关系管理

### 条件执行

```yaml
# 只在特定文件变更时运行
- uses: dorny/paths-filter@v3
  with:
    filters: |
      go_files:
        - '**/*.go'
```

## 故障排查

### 常见问题

#### 1. MongoDB 连接失败

**症状：** 集成测试失败，显示 "connection refused"

**解决方案：**
- 检查 Services 健康检查配置
- 确保等待 MongoDB 启动的脚本正确
- 验证环境变量设置

#### 2. 测试超时

**症状：** 测试运行超过 10 分钟

**解决方案：**
```yaml
# 增加超时时间
- run: go test -timeout 15m ./...
```

#### 3. 缓存问题

**症状：** 构建时间异常长

**解决方案：**
- 检查缓存键是否正确
- 手动清理缓存：Settings → Actions → Caches
- 重新运行工作流

#### 4. 权限错误

**症状：** "permission denied" 或 "403 Forbidden"

**解决方案：**
- 检查 GITHUB_TOKEN 权限：Settings → Actions → General → Workflow permissions
- 确保设置为 "Read and write permissions"

### 启用调试日志

在仓库 Settings → Secrets 中添加：

```
Name: ACTIONS_RUNNER_DEBUG
Value: true

Name: ACTIONS_STEP_DEBUG
Value: true
```

然后重新运行工作流查看详细日志。

## 分支保护规则

建议配置以下分支保护规则（Settings → Branches）：

### `main` 分支

- [x] Require a pull request before merging
- [x] Require approvals (至少 1 个)
- [x] Require status checks to pass before merging
  - lint
  - security
  - unit-tests
  - integration-tests
  - api-tests
  - build
  - dependency-check
- [x] Require branches to be up to date before merging
- [x] Require conversation resolution before merging
- [x] Do not allow bypassing the above settings

### `dev` 分支

- [x] Require status checks to pass before merging
  - lint
  - unit-tests
  - build

## 性能监控

### 查看工作流运行时间

```bash
# 使用 GitHub CLI
gh run list --workflow=ci.yml --limit 10

# 查看特定运行的详情
gh run view <run-id>
```

### 优化建议

1. **减少测试时间**
   - 使用 `-short` 标志跳过慢速测试
   - 增加测试并行度：`-parallel 4`
   - 只运行受影响的测试

2. **优化构建**
   - 使用多阶段构建
   - 最大化缓存利用
   - 减少镜像大小

3. **合理使用矩阵**
   - 只在必要时使用多版本测试
   - 考虑成本和时间平衡

## 维护检查清单

### 每周
- [ ] 检查并合并 Dependabot PR
- [ ] 查看失败的工作流并修复
- [ ] 清理旧的 workflow runs

### 每月
- [ ] 审查 Security Alerts
- [ ] 更新工具版本
- [ ] 优化缓存策略

### 每季度
- [ ] 审查并更新工作流配置
- [ ] 检查 GitHub Actions 最佳实践
- [ ] 评估新的 Actions 和工具

## 相关资源

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [CI/CD 配置指南](../../doc/ops/CI_CD配置指南.md)
- [CI/CD 问题解决方案](../../doc/ops/CI_CD问题解决方案.md)
- [项目 README](../../README.md)

## 贡献

如果需要修改工作流：

1. 在新分支上进行修改
2. 本地测试（使用 `act` 或 `make ci-local`）
3. 创建 PR 并描述变更原因
4. 等待 PR 检查通过
5. 请求代码审查

---

**最后更新：** 2025-10-20  
**维护者：** 青羽后端团队

