# CI/CD Docker 镜像拉取问题解决方案

## 问题描述

在 GitHub Actions CI/CD 流程中遇到 Docker 镜像拉取失败的问题：

```
Error response from daemon: Head "https://registry-1.docker.io/v2/securego/gosec/manifests/2.22.10": 
received unexpected HTTP status: 503 Service Unavailable
```

### 受影响的镜像

1. `securego/gosec:2.22.10` - Go 安全扫描工具
2. `mongo:6.0` - MongoDB 数据库

### 问题原因

- Docker Hub 服务暂时不可用（503 错误）
- Docker Hub 对免费用户有速率限制
- 网络连接问题
- 可能的区域性访问限制

## 解决方案

### 1. 使用本地安装的工具替代 Docker 镜像

对于 `gosec` 等开发工具，我们改用 Go 原生安装方式：

**之前的方式（依赖 Docker）：**
```yaml
- name: Run gosec
  uses: docker://securego/gosec:2.22.10
```

**改进后的方式（本地安装）：**
```yaml
- name: Install gosec
  run: go install github.com/securego/gosec/v2/cmd/gosec@latest

- name: Run gosec
  run: gosec -fmt json -out gosec-report.json ./...
```

**优势：**
- ✅ 不依赖 Docker Hub
- ✅ 安装速度更快
- ✅ 可以利用 Go modules 缓存
- ✅ 更稳定可靠

### 2. 使用 GitHub Actions Services

对于 MongoDB 等必须使用 Docker 的服务，采用 GitHub Actions Services：

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
- ✅ GitHub 自动处理镜像拉取和重试
- ✅ 内置健康检查机制
- ✅ 与 GitHub 基础设施直接集成
- ✅ 更好的网络连接

### 3. 使用 GitHub Container Registry

对于自定义镜像，推送到 GitHub Container Registry 而非 Docker Hub：

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

**优势：**
- ✅ 无速率限制
- ✅ 与 GitHub 深度集成
- ✅ 自动权限管理
- ✅ 更快的下载速度

### 4. 添加健康检查和等待逻辑

确保服务完全启动后再运行测试：

```yaml
- name: Wait for MongoDB
  run: |
    timeout 60 bash -c 'until mongosh --host localhost:27017 \
      -u admin -p password --eval "db.adminCommand({ping: 1})" --quiet; \
      do sleep 2; done'
```

## 实施的改进

### 1. 创建了完整的 CI/CD 工作流

**文件位置：** `.github/workflows/`

- **ci.yml** - 主要的持续集成流程
  - 代码检查（golangci-lint）
  - 安全扫描（本地 gosec）
  - 单元测试
  - 集成测试（使用 Services）
  - API 测试
  - 跨平台构建
  - 依赖检查

- **docker-build.yml** - Docker 镜像构建
  - 多架构支持（amd64, arm64）
  - 推送到 GitHub Container Registry
  - Trivy 安全扫描

- **pr-check.yml** - Pull Request 检查
  - PR 标题验证
  - 大文件检查
  - 敏感数据扫描
  - 代码质量检查
  - 覆盖率要求

- **release.yml** - 自动发布
  - 多平台二进制构建
  - 自动创建 GitHub Release

- **codeql.yml** - 代码安全分析
  - 自动化安全扫描
  - 定期运行（每周）

### 2. 配置文件

**golangci-lint 配置（`.golangci.yml`）：**
- 启用 20+ 个 linter
- 自定义规则和排除项
- 优化性能设置

**Dependabot 配置（`.github/dependabot.yml`）：**
- 自动更新 Go modules
- 自动更新 GitHub Actions
- 自动更新 Docker 镜像

**Issue 和 PR 模板：**
- Bug 报告模板
- 功能请求模板
- PR 检查清单
- 自动标签配置

### 3. 增强的 Makefile

新增命令：

```bash
make security        # 运行安全扫描
make vuln-check      # 检查依赖漏洞
make complexity      # 检查代码复杂度
make imports         # 整理导入
make ci-local        # 本地模拟完整 CI
make pr-check        # PR 提交前检查
make install-lint    # 安装 golangci-lint
```

## 使用指南

### 本地开发流程

1. **初始化环境**
   ```bash
   make init
   make install-lint
   ```

2. **开发前检查**
   ```bash
   make check  # 快速检查
   ```

3. **提交前检查**
   ```bash
   make pr-check  # 完整的 PR 检查
   ```

4. **本地模拟 CI**
   ```bash
   make ci-local  # 模拟 GitHub Actions
   ```

### CI/CD 流程

1. **推送代码到分支**
   - 自动运行 `ci.yml` 工作流
   - 所有检查必须通过

2. **创建 Pull Request**
   - 运行 `pr-check.yml` 额外检查
   - PR 标题格式验证
   - 自动添加标签

3. **合并到主分支**
   - 运行完整测试
   - 构建 Docker 镜像

4. **发布版本**
   - 创建 Git tag（如 `v1.0.0`）
   - 自动构建多平台二进制
   - 创建 GitHub Release

## 测试覆盖率要求

- **最低要求：** 60%
- **推荐目标：** 80%
- **检查方式：**
  ```bash
  make test-coverage-check
  ```

## 安全扫描

### 本地运行安全扫描

```bash
# 代码安全扫描
make security

# 依赖漏洞检查
make vuln-check
```

### CI 自动扫描

- 每次 PR 和推送都会运行
- 每周定时运行 CodeQL 分析
- 发现问题会自动创建 Security Alert

## 性能优化

### 1. 缓存策略

GitHub Actions 自动缓存：
- Go modules
- Go build cache
- Docker layers

### 2. 并行化

- 多个 job 并行运行
- 矩阵策略构建多平台
- 条件触发（只在相关文件变更时运行）

### 3. 测试优化

```bash
# 快速测试（跳过慢速测试）
make test-quick

# 只运行单元测试
make test-unit

# 使用并行
go test -parallel 4 ./...
```

## 故障排查

### Q: 为什么 CI 还是失败？

**A:** 检查以下几点：

1. **确保本地测试通过：**
   ```bash
   make ci-local
   ```

2. **检查 MongoDB 连接：**
   - 查看 Services 日志
   - 确认环境变量正确

3. **查看详细日志：**
   - 在 GitHub Actions 中点击失败的 job
   - 展开详细输出

### Q: 如何跳过某些检查？

**A:** 
- 在 commit message 中添加 `[skip ci]`
- 在 `.golangci.yml` 中配置排除规则
- 在测试文件中使用 build tags

### Q: 如何更新工具版本？

**A:**

1. **golangci-lint：** 编辑 `.github/workflows/ci.yml`
2. **Go 版本：** 修改 `GO_VERSION` 环境变量
3. **MongoDB 版本：** 修改 Services 中的镜像版本

## 监控和维护

### 定期检查

- [ ] 每周查看 Dependabot PR
- [ ] 每月审查 Security Alerts
- [ ] 每季度更新工具版本

### 性能监控

- 查看 workflow 运行时间
- 优化慢速测试
- 调整缓存策略

## 相关文档

- [CI/CD 配置指南](./CI_CD配置指南.md)
- [测试指南](../testing/测试指南.md)
- [部署指南](./部署指南.md)

## 总结

通过以下改进，我们完全解决了 Docker 镜像拉取问题：

✅ **使用本地工具代替 Docker 镜像**
✅ **使用 GitHub Actions Services**
✅ **推送到 GitHub Container Registry**
✅ **添加健康检查和重试机制**
✅ **完善的本地开发工具**
✅ **自动化的依赖管理**

现在 CI/CD 流程更加：
- **稳定可靠** - 不依赖外部 Docker Hub
- **快速高效** - 利用缓存和并行化
- **安全** - 自动安全扫描
- **易于维护** - 清晰的文档和工具

---

**创建日期：** 2025-10-20  
**维护者：** 青羽后端团队  
**状态：** ✅ 已解决

