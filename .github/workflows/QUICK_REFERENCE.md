# CI/CD 工作流快速参考

## 📋 当前活跃的工作流

| 工作流 | 触发条件 | 用途 | 运行时间 |
|--------|---------|------|---------|
| **ci.yml** ⭐ | Push/PR | 主CI测试 | ~15分钟 |
| **pr-check.yml** | 仅PR | PR特有检查 | ~3分钟 |
| **codeql.yml** | 定时/Push/PR | 安全扫描 | ~10分钟 |
| **docker-build.yml** | Tags/手动 | Docker构建 | ~8分钟 |
| **ci-simple.yml** 🚫 | 手动触发 | 紧急备用 | ~15分钟 |

## 🎯 工作流选择指南

### 我应该关注哪个工作流？

```
你的场景是什么？
│
├─ 日常开发 → 关注 ci.yml
│   └─ 包含所有必要的测试
│
├─ 创建PR → 关注 ci.yml + pr-check.yml
│   ├─ ci.yml: 完整测试
│   └─ pr-check.yml: PR格式检查
│
├─ 发布版本 → 关注 docker-build.yml
│   └─ 构建生产Docker镜像
│
└─ 安全审计 → 关注 codeql.yml
    └─ 代码安全扫描
```

## 🔍 ci.yml - 主CI工作流

### 包含的检查

| 检查项 | 描述 | 失败影响 | 耗时 |
|--------|------|---------|------|
| **Linting** | 代码规范检查 | 🔴 阻塞 | 2分钟 |
| **Security** | 安全扫描 | 🟡 警告 | 2分钟 |
| **Unit Tests** | 单元测试 | 🔴 阻塞 | 3分钟 |
| **Integration Tests** | 集成测试 | 🔴 阻塞 | 5分钟 |
| **API Tests** | API测试 | 🔴 阻塞 | 5分钟 |
| **Dependency Check** | 依赖检查 | 🟡 警告 | 2分钟 |

### 何时运行

```yaml
✅ Push到 main, dev, develop
✅ Pull Request 到 main, dev, develop
❌ 其他分支（需要手动触发）
```

### 依赖检查特殊规则

```
feature分支 → 跳过依赖检查
dev分支     → 运行依赖检查
main分支    → 运行依赖检查
```

### 本地测试命令

```bash
# 快速验证（推荐在提交前运行）
golangci-lint run --timeout=10m
go test -v -short ./...

# 完整测试（等同CI）
./scripts/run_tests_with_docker.sh
```

## 🎨 pr-check.yml - PR检查

### 包含的检查

| 检查项 | 描述 | 示例 |
|--------|------|------|
| **PR标题** | Semantic格式 | `feat: add login` |
| **大文件** | 检测>5MB文件 | 避免提交大文件 |
| **敏感数据** | 检测密钥泄露 | API keys, 密码等 |
| **代码复杂度** | 函数复杂度<15 | 重构复杂函数 |
| **快速检查** | go vet + build | 语法检查 |
| **Docker构建** | 仅Docker变更时 | 验证构建 |

### PR标题格式

**格式**: `<type>: <description>`

**允许的类型**:
```
feat      - 新功能
fix       - 修复bug
docs      - 文档更新
style     - 代码格式
refactor  - 重构
perf      - 性能优化
test      - 测试相关
build     - 构建相关
ci        - CI配置
chore     - 其他杂项
revert    - 回退
```

**示例**:
```
✅ feat: add user authentication
✅ fix: resolve login timeout issue
✅ docs: update API documentation
❌ Add new feature (缺少type)
❌ Feat: Add login (大写F)
```

### 大文件检查

```bash
# 本地检查
find . -type f -size +5M -not -path "./.git/*"

# 使用Git LFS
git lfs track "*.pdf"
git lfs track "*.zip"
```

## 🚦 CI状态理解

### 状态图标

| 图标 | 状态 | 含义 | 行动 |
|------|------|------|------|
| 🟢 | Success | 全部通过 | 可以合并 |
| 🔴 | Failure | 有失败 | 必须修复 |
| 🟡 | Warning | 有警告 | 建议修复 |
| ⚪ | Pending | 运行中 | 等待 |
| ⭕ | Skipped | 已跳过 | 正常 |

### 检查结果分级

**核心检查** (必须通过):
- Linting
- Unit Tests
- Integration Tests
- API Tests

**辅助检查** (可以警告):
- Security Scan
- Dependency Check

```
核心检查失败 → 🔴 PR被阻塞
辅助检查失败 → 🟡 可合并但有警告
```

## 🔧 常用操作

### 重新运行失败的检查

```
1. 进入PR页面
2. 点击"Details"查看失败原因
3. 修复问题后重新提交
4. 或点击"Re-run jobs"重新运行
```

### 手动触发CI

```
1. 进入Actions页面
2. 选择工作流
3. 点击"Run workflow"
4. 选择分支
5. 点击"Run workflow"按钮
```

### 查看详细日志

```
1. 进入Actions页面
2. 点击运行记录
3. 点击失败的Job
4. 展开失败的步骤
5. 查看错误信息
```

## 📊 性能优化建议

### 减少CI时间

**代码提交前**:
```bash
# 1. 本地运行linter
golangci-lint run

# 2. 格式化代码
gofmt -w .

# 3. 本地测试
go test -short ./...

# 4. 提交代码
git commit
```

**PR创建时**:
```
1. 确保PR标题格式正确
2. 不包含大文件
3. 代码已经过本地测试
4. 提交前运行go vet
```

### 避免重复CI

```
❌ 避免：频繁提交小改动
✅ 推荐：积累改动后一次提交

❌ 避免：提交后立即再提交
✅ 推荐：等待CI结果后再提交

❌ 避免：在PR中频繁force push
✅ 推荐：本地测试后再push
```

## 🐛 常见问题

### Q1: 为什么我的PR运行了两个CI？

**A**: 正常现象
- `ci.yml`: 主CI，运行完整测试
- `pr-check.yml`: PR检查，运行PR特有验证

### Q2: 依赖检查为什么被跳过？

**A**: 优化策略
- 只在 `main` 和 `dev` 分支运行
- Feature分支会跳过以节省时间

### Q3: 如何让CI运行得更快？

**A**: 
1. 本地测试后再提交
2. 避免频繁小提交
3. 使用 `go test -short` 跳过慢速测试
4. 优化测试代码

### Q4: CI失败但我认为是误报？

**A**:
1. 查看详细日志确认
2. 如果确实是误报，联系维护者
3. 可以临时在代码中添加忽略指令
4. 提issue讨论是否需要调整规则

### Q5: 如何在不触发CI的情况下提交？

**A**: 
```bash
# 在commit message中添加 [skip ci]
git commit -m "docs: update README [skip ci]"
```

**注意**: 仅用于文档更新等不影响代码的改动

## 📚 进一步阅读

- [CI优化详细说明](README_OPTIMIZATION.md)
- [Docker测试环境](../../docker/README_TEST.md)
- [快速测试指南](../../docker/QUICK_TEST_GUIDE.md)

## 🆘 获取帮助

**CI相关问题**:
1. 查看本文档
2. 查看 Actions 日志
3. 查看 README_OPTIMIZATION.md
4. 联系团队维护者

**紧急情况**:
- 如果 `ci.yml` 有问题，可以手动触发 `ci-simple.yml`
- 查看 [故障排查指南](README_OPTIMIZATION.md#故障排查)

---

**最后更新**: 2025-10-22  
**维护者**: 青羽后端团队

