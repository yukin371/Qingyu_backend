# 自动分配审查者配置说明

**创建日期**: 2025-10-22  
**配置版本**: 1.0

## 📋 配置文件概览

| 文件 | 作用 | 优先级 |
|------|------|--------|
| `CODEOWNERS` | 基于文件路径自动分配审查者 | 高 ⭐ |
| `auto_assign.yml` | 配置自动分配规则 | 中 |
| `workflows/auto-assign.yml` | 执行自动分配的 Action | 中 |

## 🎯 CODEOWNERS 文件

### 作用
- ✅ GitHub 原生支持，无需额外配置
- ✅ 基于文件路径自动分配审查者
- ✅ 适用于所有 PR（包括 Dependabot）
- ✅ 强制审查（可配合分支保护使用）

### 工作原理

```
1. PR 创建
   ↓
2. GitHub 检测修改的文件
   ↓
3. 匹配 CODEOWNERS 规则
   ↓
4. 自动添加对应的审查者
   ↓
5. 审查者收到通知
```

### 规则示例

```
# 通配符规则
* @yukin371                    # 所有文件

# 目录规则
api/ @yukin371                 # api 目录下所有文件
service/ @yukin371             # service 目录下所有文件

# 文件类型规则
*.go @yukin371                 # 所有 .go 文件
*_test.go @yukin371           # 所有测试文件

# 特定文件规则
go.mod @yukin371              # go.mod 文件
docker-compose*.yml @yukin371 # 所有 docker-compose 文件

# 优先级：更具体的规则优先
# 例如：如果修改了 api/user.go
# 会匹配 api/ 规则，而不是 * 规则
```

### 验证配置

```bash
# 本地验证语法（如果安装了 GitHub CLI）
gh codeowners check

# 查看特定文件的所有者
gh codeowners show path/to/file
```

## 🤖 Auto Assign Action

### 作用
- ✅ 更灵活的自动分配规则
- ✅ 支持随机分配多个审查者
- ✅ 支持跳过关键字（如 WIP）
- ✅ 可以分配 Assignees（负责人）

### 配置选项

```yaml
# auto_assign.yml

# 基础配置
addReviewers: true              # 添加审查者
addAssignees: true              # 添加负责人
numberOfReviewers: 1            # 分配数量

# 用户列表
reviewers:
  - yukin371

# 跳过规则
skipKeywords:
  - wip                         # PR 标题包含 wip 时跳过
  - "[WIP]"
  - "[skip assign]"
```

### 工作流程

```
1. PR 打开或重新打开
   ↓
2. 检查是否是草稿 PR
   ↓ (不是草稿)
3. 检查 PR 标题/描述是否包含跳过关键字
   ↓ (不包含)
4. 从配置中选择审查者
   ↓
5. 自动添加到 PR
   ↓
6. 发送通知
```

## 🔀 两种方式的对比

| 特性 | CODEOWNERS | Auto Assign Action |
|------|-----------|-------------------|
| GitHub 原生支持 | ✅ | ❌ 需要 Action |
| 基于文件路径 | ✅ | ⚠️ 可配置 |
| 随机分配 | ❌ | ✅ |
| 跳过草稿 PR | ❌ | ✅ |
| 强制审查 | ✅ 配合分支保护 | ❌ |
| 分配负责人 | ❌ | ✅ |
| 配置复杂度 | 简单 | 中等 |

## 📊 推荐配置策略

### 策略1：仅使用 CODEOWNERS（推荐）⭐

**优势**：
- ✅ 简单可靠
- ✅ GitHub 原生支持
- ✅ 无需额外 Action

**适用场景**：
- 小团队（1-3人）
- 规则简单
- 主要关注代码审查

**配置**：
```
只需 CODEOWNERS 文件即可
可以删除 auto_assign.yml 和 workflows/auto-assign.yml
```

### 策略2：组合使用（功能最全）

**优势**：
- ✅ CODEOWNERS 处理代码审查
- ✅ Auto Assign 处理负责人分配
- ✅ 支持跳过草稿 PR

**适用场景**：
- 中大团队（4+人）
- 需要区分审查者和负责人
- 复杂的分配规则

**配置**：
```
CODEOWNERS → 基于文件路径的审查者
Auto Assign → 自动分配负责人和跳过规则
```

### 策略3：仅使用 Auto Assign

**优势**：
- ✅ 灵活性高
- ✅ 支持所有高级功能

**劣势**：
- ❌ 需要维护 GitHub Action
- ❌ 不是 GitHub 原生功能

**适用场景**：
- 需要复杂的分配逻辑
- 需要随机分配功能

## 🚀 使用方法

### 基础使用（CODEOWNERS）

1. **无需额外操作**
   - CODEOWNERS 文件已创建
   - 下次 PR 自动生效

2. **验证**
   ```bash
   # 创建测试 PR
   git checkout -b test-codeowners
   echo "test" >> README.md
   git commit -am "test: verify CODEOWNERS"
   git push origin test-codeowners
   
   # 在 GitHub 创建 PR
   # 检查是否自动添加了审查者
   ```

3. **查看效果**
   - PR 页面会显示 "Reviewers" 部分
   - 自动添加的审查者会显示
   - 审查者会收到邮件通知

### 高级使用（Auto Assign Action）

1. **启用 Action**
   - 文件已创建
   - 下次 PR 时自动运行

2. **查看运行结果**
   ```
   PR → Checks → Auto Assign
   查看是否成功分配
   ```

3. **调试**
   ```bash
   # 查看 Action 日志
   # GitHub → Actions → Auto Assign Reviewers and Assignees
   # 点击具体的运行记录查看日志
   ```

## ⚙️ 自定义配置

### 添加更多审查者

**CODEOWNERS**:
```
# 不同目录不同审查者
api/ @yukin371 @backend-team
service/ @yukin371
test/ @yukin371 @qa-team
```

**auto_assign.yml**:
```yaml
reviewers:
  - yukin371
  - user2
  - user3

numberOfReviewers: 2  # 随机选2个
```

### 基于文件类型分配

**CODEOWNERS**:
```
# 后端代码
*.go @backend-team

# 前端代码
*.ts @frontend-team
*.tsx @frontend-team

# 文档
*.md @doc-team

# 配置文件
*.yml @devops-team
*.yaml @devops-team
docker/* @devops-team
```

### 特殊规则

**跳过某些 PR**:
```yaml
# auto_assign.yml
skipKeywords:
  - "[skip assign]"
  - "[bot]"
  - "chore(deps)"  # 跳过依赖更新
```

**PR 标题**:
```
chore(deps): update dependencies [skip assign]
```

## 🔍 故障排查

### CODEOWNERS 不工作

**检查**:
1. 文件位置是否正确：`.github/CODEOWNERS`
2. 语法是否正确（每行一个规则）
3. 用户名是否正确（@username）
4. 用户是否有仓库访问权限

**解决**:
```bash
# 验证语法
cat .github/CODEOWNERS

# 检查用户权限
# Settings → Collaborators
```

### Auto Assign Action 失败

**常见原因**:
1. `GITHUB_TOKEN` 权限不足
2. 配置文件路径错误
3. 用户不存在

**解决**:
```yaml
# 检查权限配置
permissions:
  contents: read
  pull-requests: write  # 必需

# 检查配置路径
configuration-path: ".github/auto_assign.yml"  # 确保正确
```

## 📈 最佳实践

1. **从简单开始**
   - 先只使用 CODEOWNERS
   - 需要时再添加 Auto Assign

2. **保持规则清晰**
   - 添加注释说明规则用途
   - 定期审查和更新规则

3. **测试配置**
   - 创建测试 PR 验证
   - 确保通知正常工作

4. **团队协作**
   - 让团队知道自动分配规则
   - 培训如何正确使用

5. **监控效果**
   - 定期检查是否正确分配
   - 收集团队反馈优化规则

## 🎯 针对 Dependabot 的优化

**CODEOWNERS**:
```
# Dependabot 会修改这些文件
go.mod @yukin371
go.sum @yukin371
.github/workflows/* @yukin371
docker/Dockerfile* @yukin371

# Dependabot 创建的 PR 会自动分配审查者
```

**效果**:
```
Dependabot 创建 PR
   ↓
自动分配 @yukin371 为审查者
   ↓
@yukin371 收到邮件通知
   ↓
审查并合并
```

## 📚 相关文档

- [GitHub CODEOWNERS 文档](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)
- [Auto Assign Action](https://github.com/kentaro-m/auto-assign-action)
- [分支保护规则](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)

## 🔄 维护清单

- [ ] 定期审查 CODEOWNERS 规则
- [ ] 团队成员变动时更新配置
- [ ] 测试新规则是否生效
- [ ] 收集团队反馈
- [ ] 优化自动分配策略

---

**维护者**: @yukin371  
**最后更新**: 2025-10-22  
**配置版本**: 1.0

## 💡 快速参考

**只想要基础功能？**
→ 保留 CODEOWNERS，删除其他文件

**需要更多控制？**
→ 同时使用 CODEOWNERS 和 Auto Assign

**配置不生效？**
→ 查看本文档的"故障排查"部分

