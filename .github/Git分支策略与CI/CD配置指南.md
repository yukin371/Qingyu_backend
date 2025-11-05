# Git 分支策略与 CI/CD 配置指南

**制定日期**: 2025-10-22  
**项目阶段**: 早期开发  
**版本**: 1.0

## 🌳 分支模型

### 推荐的分支结构

```
main (生产环境)
  ├─ tag: v1.0.0, v1.1.0 (版本标签)
  ↑
dev (开发环境，默认分支)
  ↑
feature/* (功能分支)
  ├─ feature/user-auth
  ├─ feature/document-api
  └─ feature/ai-integration

hotfix/* (紧急修复)
  ├─ hotfix/security-patch
  └─ hotfix/critical-bug

release/* (发布分支，可选)
  └─ release/v1.0.0
```

## 📋 分支详细说明

### 1. main 分支（生产/稳定分支）⭐

**用途**:
- 生产环境部署的代码
- 对外发布的稳定版本
- 每次合并都应该是一个可发布的版本

**保护策略**:
```yaml
✅ 需要 PR 才能合并
✅ 需要通过所有 CI 测试
✅ 需要代码审查（项目成熟后）
✅ 禁止直接 push
✅ 禁止 force push
```

**合并来源**:
- ← dev 分支（经过充分测试）
- ← hotfix/* 分支（紧急修复）
- ← release/* 分支（发布版本）

**部署**:
- 自动部署到生产环境
- 或手动批准后部署

**标签管理**:
```bash
# 每次发布打标签
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### 2. dev 分支（开发分支）⭐ 默认分支

**用途**:
- 日常开发的主分支
- 集成所有完成的功能
- 开发环境部署的代码

**保护策略**:
```yaml
✅ 需要 PR 才能合并（推荐）
✅ 需要通过 CI 测试
⚠️ 代码审查可选（早期阶段）
❌ 允许直接 push（早期阶段，可调整）
```

**合并来源**:
- ← feature/* 分支
- ← bugfix/* 分支
- ← 临时修复提交

**部署**:
- 自动部署到开发/测试环境

**设为默认分支**:
```
GitHub → Settings → Branches → Default branch
选择: dev
```

### 3. feature/* 分支（功能分支）

**命名规范**:
```
feature/user-authentication
feature/document-upload
feature/ai-text-generation
feature/payment-integration
```

**生命周期**:
```
1. 从 dev 创建
2. 开发功能
3. 本地测试
4. 创建 PR 到 dev
5. 通过 CI 和审查
6. 合并到 dev
7. 删除 feature 分支
```

**创建和使用**:
```bash
# 1. 从 dev 创建新分支
git checkout dev
git pull origin dev
git checkout -b feature/user-login

# 2. 开发功能
# ... 编写代码 ...

# 3. 提交
git add .
git commit -m "feat: add user login functionality"

# 4. 推送到远程
git push origin feature/user-login

# 5. 创建 PR
gh pr create --base dev --title "feat: add user login"

# 6. 合并后删除分支
git branch -d feature/user-login
git push origin --delete feature/user-login
```

### 4. hotfix/* 分支（热修复分支）

**用途**:
- 紧急修复生产环境的严重 bug
- 直接从 main 分支创建

**命名规范**:
```
hotfix/security-vulnerability
hotfix/payment-failure
hotfix/data-loss-bug
```

**生命周期**:
```
1. 从 main 创建
2. 修复 bug
3. 测试验证
4. 合并到 main（打标签）
5. 同时合并到 dev（保持同步）
6. 删除 hotfix 分支
```

**使用流程**:
```bash
# 1. 从 main 创建
git checkout main
git pull origin main
git checkout -b hotfix/critical-bug

# 2. 修复 bug
# ... 修复代码 ...

# 3. 提交
git commit -am "fix: resolve critical payment bug"

# 4. 合并到 main
git checkout main
git merge --no-ff hotfix/critical-bug
git tag -a v1.0.1 -m "Hotfix: critical bug"
git push origin main --tags

# 5. 同时合并到 dev
git checkout dev
git merge --no-ff hotfix/critical-bug
git push origin dev

# 6. 删除分支
git branch -d hotfix/critical-bug
git push origin --delete hotfix/critical-bug
```

### 5. release/* 分支（发布分支，可选）

**用途**:
- 准备新版本发布
- 版本号调整
- 最后的 bug 修复
- 发布文档更新

**命名规范**:
```
release/v1.0.0
release/v2.0.0
```

**生命周期**:
```
1. 从 dev 创建
2. 版本准备工作（文档、版本号等）
3. bug 修复（不添加新功能）
4. 测试验证
5. 合并到 main（打标签）
6. 同时合并到 dev
7. 删除 release 分支
```

## 🚀 CI/CD 覆盖策略

### 推荐的 CI 触发配置

#### 方案1：完整覆盖（推荐）⭐

```yaml
# .github/workflows/ci.yml
on:
  push:
    branches: 
      - main          # 生产分支
      - dev           # 开发分支
  pull_request:
    branches: 
      - main          # PR 到 main
      - dev           # PR 到 dev
```

**优势**:
- ✅ 确保 main 和 dev 代码质量
- ✅ PR 在合并前必须通过测试
- ✅ 及早发现问题

**CI 运行情况**:
| 操作 | 触发 CI | 说明 |
|------|--------|------|
| Push 到 main | ✅ | 运行完整测试 |
| Push 到 dev | ✅ | 运行完整测试 |
| Push 到 feature/* | ❌ | 不触发（节省资源） |
| PR 到 main | ✅ | 运行完整测试 |
| PR 到 dev | ✅ | 运行完整测试 |

#### 方案2：分级测试（高级）

```yaml
# .github/workflows/ci.yml
on:
  push:
    branches: 
      - main
  pull_request:
    branches: 
      - main
      - dev
  # 手动触发
  workflow_dispatch:
```

**特点**:
- main 分支：Push 时运行（确保生产质量）
- dev 分支：只在 PR 时运行（减少 CI 消耗）
- feature 分支：不自动运行（本地测试）

#### 方案3：最小化（早期项目）

```yaml
# .github/workflows/ci.yml
on:
  push:
    branches: 
      - main          # 只在 main 运行
  pull_request:
    branches: 
      - main
      - dev
```

**特点**:
- 只保护 main 分支
- PR 必须通过测试
- 减少 CI 消耗

## 📊 完整工作流程

### 日常功能开发

```
1. 从 dev 创建 feature 分支
   git checkout dev
   git pull origin dev
   git checkout -b feature/new-feature

2. 开发功能
   # 编写代码
   # 本地测试
   git commit -am "feat: add new feature"

3. 推送到远程
   git push origin feature/new-feature

4. 创建 PR 到 dev
   gh pr create --base dev --title "feat: add new feature"
   
5. CI 自动运行
   ✅ Linting
   ✅ Unit Tests
   ✅ Integration Tests
   ✅ API Tests
   
6. PR 审查（可选）
   
7. 合并到 dev
   # 自动或手动合并
   
8. 部署到测试环境
   # 自动部署
   
9. 删除 feature 分支
   git branch -d feature/new-feature
```

### 发布到生产

```
方式1: 直接从 dev 发布
dev → main → 生产环境

方式2: 使用 release 分支
dev → release/v1.0.0 → main → 生产环境

步骤：
1. 确保 dev 分支稳定
2. 创建 PR: dev → main
3. 通过所有 CI 测试
4. 代码审查
5. 合并到 main
6. 打标签: v1.0.0
7. 自动部署到生产
```

### 紧急修复

```
1. 从 main 创建 hotfix 分支
   git checkout main
   git checkout -b hotfix/critical-bug

2. 修复 bug
   git commit -am "fix: critical bug"

3. 合并到 main
   git checkout main
   git merge --no-ff hotfix/critical-bug
   git tag -a v1.0.1 -m "Hotfix"
   git push origin main --tags

4. 同步到 dev
   git checkout dev
   git merge --no-ff hotfix/critical-bug
   git push origin dev

5. 删除 hotfix 分支
```

## 🎯 分支策略决策表

### 新开分支应该基于哪个分支？

| 场景 | 基于分支 | 目标分支 | 示例 |
|------|---------|---------|------|
| **新功能开发** | dev | dev | feature/user-login |
| **Bug 修复** | dev | dev | bugfix/form-validation |
| **生产紧急修复** | main | main + dev | hotfix/security-patch |
| **版本发布准备** | dev | main | release/v1.0.0 |
| **实验性功能** | dev | - | experiment/new-ui |
| **文档更新** | dev | dev | docs/api-update |

### CI 测试应该覆盖哪些分支？

| 分支 | Push 时 | PR 时 | 原因 |
|------|---------|-------|------|
| **main** | ✅ 必须 | ✅ 必须 | 生产代码，必须保证质量 |
| **dev** | ✅ 推荐 | ✅ 必须 | 集成分支，需要验证 |
| **feature/*** | ❌ 不需要 | ✅ 在 PR 时 | 节省资源，PR 时验证 |
| **hotfix/*** | ❌ 不需要 | ✅ 在 PR 时 | PR 时验证即可 |
| **release/*** | ⚠️ 可选 | ✅ 必须 | 发布前最后验证 |

## 🔧 推荐配置实施

### 步骤1：设置默认分支

```bash
# GitHub 网页操作
Settings → Branches → Default branch → dev

# 或使用 GitHub CLI
gh repo edit --default-branch dev
```

### 步骤2：设置分支保护规则

**main 分支**:
```
GitHub → Settings → Branches → Add rule

Branch name pattern: main

保护规则：
☑️ Require a pull request before merging
  ☑️ Require approvals (1)
☑️ Require status checks to pass before merging
  ☑️ Require branches to be up to date
  选择: CI checks (lint, unit-tests, integration-tests)
☑️ Require conversation resolution before merging
☐ Require signed commits (可选)
☑️ Include administrators (推荐)
☑️ Restrict deletions
☐ Allow force pushes (禁用)
```

**dev 分支**（早期阶段，宽松配置）:
```
Branch name pattern: dev

保护规则：
☑️ Require status checks to pass before merging
  选择: CI checks
☐ Require pull request before merging (可选)
☐ Include administrators (允许管理员直接 push)
```

### 步骤3：更新 CI 配置

修改 `.github/workflows/ci.yml`:
```yaml
name: CI Pipeline

on:
  push:
    branches: 
      - main
      - dev
    # 可选：添加路径过滤
    # paths-ignore:
    #   - '**.md'
    #   - 'doc/**'
      
  pull_request:
    branches: 
      - main
      - dev
      
  # 允许手动触发
  workflow_dispatch:

env:
  GO_VERSION: '1.21'
```

### 步骤4：清理不需要的分支

```bash
# 查看所有分支
git branch -a

# 删除本地不需要的分支
git branch -d test  # 如果不再使用
git branch -d develop  # 如果与 dev 重复

# 删除远程分支
git push origin --delete test
git push origin --delete develop
```

## 📋 Git 工作流最佳实践

### 提交信息规范

使用 Conventional Commits:
```
feat: 新功能
fix: 修复 bug
docs: 文档更新
style: 代码格式
refactor: 重构
perf: 性能优化
test: 测试相关
build: 构建系统
ci: CI 配置
chore: 其他杂项

示例：
feat: add user authentication
fix: resolve login timeout issue
docs: update API documentation
```

### 分支命名规范

```
功能：feature/功能名称
修复：bugfix/问题描述  或  fix/问题描述
热修复：hotfix/问题描述
发布：release/版本号
实验：experiment/实验名称
文档：docs/文档主题
```

### 合并策略

**feature → dev**:
```bash
# 推荐：Squash and merge（合并为一个提交）
# 或：Merge commit（保留所有提交历史）
```

**dev → main**:
```bash
# 推荐：Merge commit（保留完整历史）
# 打标签标记版本
```

**hotfix → main**:
```bash
# 推荐：Merge commit（--no-ff）
# 立即打标签
```

---

**维护者**: 青羽后端团队  
**最后更新**: 2025-10-22  
**建议复审**: 项目进入稳定期时

