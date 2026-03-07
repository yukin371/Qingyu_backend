# CI/CD 配置指南

## 概述

本项目使用 GitHub Actions 实现完整的 CI/CD 流程，包括代码检查、测试、构建、发布等环节。

## 目录

- [工作流说明](#工作流说明)
- [环境配置](#环境配置)
- [Docker 镜像拉取问题解决方案](#docker-镜像拉取问题解决方案)
- [本地运行 CI 检查](#本地运行-ci-检查)
- [常见问题](#常见问题)

## 工作流说明

### 1. CI Pipeline (`.github/workflows/ci.yml`)

主要的持续集成流程，在推送代码和创建 PR 时触发。

**包含的任务：**

- **代码检查 (lint)**
  - 使用 `golangci-lint` 进行代码质量检查
  - 检查代码格式是否符合 `gofmt` 标准
  
- **安全扫描 (security)**
  - 使用本地安装的 `gosec` 而非 Docker 镜像（避免 Docker Hub 503 错误）
  - 生成安全报告并上传为 artifact

- **单元测试 (unit-tests)**
  - 运行所有单元测试
  - 生成代码覆盖率报告
  - 上传到 Codecov

- **集成测试 (integration-tests)**
  - 使用 GitHub Actions Services 启动 MongoDB
  - 运行集成测试
  - 包含健康检查和自动重试

- **API 测试 (api-tests)**
  - 测试 API 端点
  - 验证请求/响应格式

- **构建测试 (build)**
  - 跨平台构建验证（Linux、macOS、Windows）
  - 多架构支持（amd64、arm64）

- **依赖检查 (dependency-check)**
  - 使用 `govulncheck` 检查依赖漏洞
  - 验证 `go.mod` 和 `go.sum` 一致性

### 2. Docker 构建 (`.github/workflows/docker-build.yml`)

构建和推送 Docker 镜像到 GitHub Container Registry。

**特性：**
- 使用 Docker Buildx 构建多架构镜像
- 推送到 GitHub Container Registry（避免 Docker Hub 限制）
- Trivy 安全扫描
- 构建缓存优化

### 3. PR 检查 (`.github/workflows/pr-check.yml`)

Pull Request 专属的额外检查。

**包含：**
- PR 标题格式验证（遵循 Conventional Commits）
- 大文件检查（>5MB）
- 敏感数据扫描（使用 TruffleHog）
- 代码复杂度检查
- 测试覆盖率要求（>=60%）
- 自动标签添加

### 4. 发布流程 (`.github/workflows/release.yml`)

创建 Git tag 时自动发布。

**流程：**
1. 运行完整测试
2. 构建多平台二进制文件
3. 生成 SHA256 校验和
4. 创建 GitHub Release
5. 上传构建产物

### 5. CodeQL 分析 (`.github/workflows/codeql.yml`)

自动化的代码安全分析。

**触发条件：**
- 推送到 main/dev 分支
- PR 到 main 分支
- 每周一定时运行

## 环境配置

### 必需的环境变量

在 GitHub 仓库设置中配置以下 Secrets：

```bash
# MongoDB 配置（用于集成测试）
MONGODB_URI=mongodb://admin:password@localhost:27017
MONGODB_DATABASE=qingyu_test

# JWT 配置
JWT_SECRET=your-jwt-secret-key

# AI 服务配置（可选）
OPENAI_API_KEY=your-openai-api-key
DASHSCOPE_API_KEY=your-dashscope-api-key
```

### GitHub Token 权限

确保 `GITHUB_TOKEN` 具有以下权限：

- `contents: write` - 用于创建 release
- `packages: write` - 用于推送 Docker 镜像
- `security-events: write` - 用于上传安全扫描结果

在仓库设置中：Settings → Actions → General → Workflow permissions

## Docker 镜像拉取问题解决方案

### 问题描述

在 CI/CD 中遇到的主要问题：

```
Error response from daemon: Head "https://registry-1.docker.io/v2/.../manifests/...": 
received unexpected HTTP status: 503 Service Unavailable
```

### 解决方案

我们采用了以下策略来解决 Docker Hub 503 错误：

#### 1. 使用本地工具替代 Docker 镜像

**之前（使用 Docker 镜像）：**
```yaml
- name: Run gosec
  uses: securego/gosec@master
  with:
    image: docker://securego/gosec:2.22.10
```

**现在（使用本地安装）：**
```yaml
- name: Install gosec
  run: go install github.com/securego/gosec/v2/cmd/gosec@latest

- name: Run gosec
  run: gosec -fmt json -out gosec-report.json ./...
```

#### 2. 使用 GitHub Actions Services

**对于 MongoDB 等服务，使用 GitHub Actions Services：**

```yaml
services:
  mongodb:
    image: mongo:6.0
    ports:
      - 27017:27017
    env:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password
    options: >-
      --health-cmd "mongosh --eval 'db.adminCommand({ping: 1})' --quiet"
      --health-interval 10s
      --health-timeout 5s
      --health-retries 5
      --health-start-period 40s
```

**优势：**
- GitHub 自动处理镜像拉取和重试
- 内置健康检查
- 更稳定的网络连接

#### 3. 使用 GitHub Container Registry

推送 Docker 镜像到 GitHub Container Registry 而非 Docker Hub：

```yaml
env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

steps:
  - name: Log in to Container Registry
    uses: docker/login-action@v3
    with:
      registry: ${{ env.REGISTRY }}
      username: ${{ github.actor }}
      password: ${{ secrets.GITHUB_TOKEN }}
```

#### 4. 镜像加速器配置（国内环境）

如果在自托管的 Runner 上运行，可以配置镜像加速器：

**Linux/macOS:**

编辑 `/etc/docker/daemon.json`：

```json
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.ccs.tencentyun.com"
  ]
}
```

重启 Docker：
```bash
sudo systemctl daemon-reload
sudo systemctl restart docker
```

**Windows:**

在 Docker Desktop → Settings → Docker Engine 中添加：

```json
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn"
  ]
}
```

## 本地运行 CI 检查

### 1. 代码格式检查

```bash
# 检查格式
gofmt -l .

# 自动格式化
gofmt -w .

# 或使用 goimports
go install golang.org/x/tools/cmd/goimports@latest
goimports -w .
```

### 2. Linter 检查

```bash
# 安装 golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
  sh -s -- -b $(go env GOPATH)/bin v1.55.2

# 运行检查
golangci-lint run --timeout=10m

# 只检查新代码
golangci-lint run --new-from-rev=HEAD~1
```

### 3. 安全扫描

```bash
# 安装 gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest

# 运行扫描
gosec ./...

# 生成 JSON 报告
gosec -fmt json -out gosec-report.json ./...
```

### 4. 运行测试

```bash
# 单元测试
go test -v -race ./...

# 带覆盖率
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# 查看覆盖率
go tool cover -func=coverage.txt
go tool cover -html=coverage.txt -o coverage.html
```

### 5. 依赖检查

```bash
# 安装 govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest

# 检查漏洞
govulncheck ./...

# 验证依赖
go mod verify
go mod tidy
```

### 6. 使用 Makefile

项目提供了 Makefile 简化命令：

```bash
# 运行所有检查
make check

# 运行测试
make test

# 生成覆盖率报告
make coverage

# 构建
make build

# 清理
make clean
```

## 常见问题

### Q1: CI 中 MongoDB 连接失败

**A**: 确保使用了健康检查并添加了等待逻辑：

```yaml
- name: Wait for MongoDB
  run: |
    timeout 60 bash -c 'until mongosh --host localhost:27017 \
      -u admin -p password --eval "db.adminCommand({ping: 1})" --quiet; \
      do sleep 2; done'
```

### Q2: 测试超时

**A**: 增加超时时间并使用 `-timeout` 参数：

```bash
go test -v -race -timeout 10m ./...
```

### Q3: golangci-lint 运行缓慢

**A**: 
- 使用缓存：GitHub Actions 自动启用
- 调整超时：`golangci-lint run --timeout=10m`
- 只检查变更：`golangci-lint run --new-from-rev=HEAD~1`

### Q4: 如何跳过 CI 检查

**A**: 在 commit message 中添加 `[skip ci]` 或 `[ci skip]`：

```bash
git commit -m "docs: update README [skip ci]"
```

### Q5: 如何调试失败的 workflow

**A**: 
1. 查看详细日志：点击失败的 job 查看完整输出
2. 本地复现：使用相同的命令在本地运行
3. 启用 debug 日志：在仓库 Settings → Secrets 中添加：
   - Name: `ACTIONS_RUNNER_DEBUG`, Value: `true`
   - Name: `ACTIONS_STEP_DEBUG`, Value: `true`

### Q6: Docker 构建缓存问题

**A**: 
- GitHub Actions 自动使用 cache-from/cache-to
- 手动清理缓存：Settings → Actions → Caches
- 重新构建：在 workflow 中使用 `cache-from: type=gha`

## 性能优化建议

### 1. 并行化测试

```bash
# 使用并行测试
go test -v -race -parallel 4 ./...
```

### 2. 缓存依赖

GitHub Actions 配置已包含缓存：

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.21'
    cache: true  # 自动缓存 go modules
```

### 3. 按需运行测试

使用 path filter 只在相关文件变更时运行测试：

```yaml
- uses: dorny/paths-filter@v3
  with:
    filters: |
      go_files:
        - '**/*.go'
```

### 4. 复用构建产物

```yaml
- uses: actions/upload-artifact@v4
  with:
    name: binary
    path: ./build/

- uses: actions/download-artifact@v4
  with:
    name: binary
```

## 监控和通知

### GitHub Status Checks

所有 workflow 作为必需的状态检查：

Settings → Branches → Branch protection rules → Require status checks

### Slack/Discord 通知

添加通知步骤到 workflow：

```yaml
- name: Notify on failure
  if: failure()
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```

## 最佳实践

1. **保持 workflow 快速**：总运行时间应在 10 分钟内
2. **使用矩阵策略**：并行测试多个版本/平台
3. **合理使用缓存**：加速依赖下载
4. **明确的错误信息**：便于快速定位问题
5. **定期更新依赖**：使用 Dependabot 自动更新
6. **保护敏感信息**：使用 Secrets 管理密钥
7. **文档同步**：workflow 变更时更新文档

## 相关资源

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [golangci-lint 文档](https://golangci-lint.run/)
- [gosec 文档](https://github.com/securego/gosec)
- [Docker Buildx 文档](https://docs.docker.com/buildx/working-with-buildx/)
- [项目测试文档](../testing/测试指南.md)

## 更新历史

| 日期 | 版本 | 变更说明 |
|-----|------|---------|
| 2025-10-20 | 1.0 | 初始版本，解决 Docker 镜像拉取问题 |

---

**维护者**: 青羽后端团队  
**最后更新**: 2025-10-20

