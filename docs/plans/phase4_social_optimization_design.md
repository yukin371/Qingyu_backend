# 第4阶段：社交功能优化 - 功能设计文档

> **创建日期**: 2026-02-27
> **状态**: 设计阶段
> **目标**: 优化social模块API代码质量，添加简单功能

---

## 1. 现状分析

### 1.1 模块清单

| API文件 | 代码行数 | 测试覆盖 | 优化状态 |
|---------|---------|---------|---------|
| collection_api.go | ~608 | ✅ 有测试 | ✅ 已优化 |
| like_api.go | ~295 | ✅ 有测试 | ⚠️ 需优化 |
| follow_api.go | ~390 | ✅ 有测试 | ⚠️ 需优化 |
| comment_api.go | ~416 | ✅ 有测试 | ⚠️ 需检查 |
| rating_api.go | ? | ✅ 有测试 | ⚠️ 需检查 |
| review_api.go | ? | ✅ 有测试 | ⚠️ 需检查 |
| booklist_api.go | ? | ✅ 有测试 | ⚠️ 需检查 |
| message_api.go | ? | ✅ 有测试 | ⚠️ 需检查 |
| relation_api.go | ? | ✅ 有测试 | ⚠️ 需检查 |

### 1.2 代码质量问题

#### 问题1：未使用辅助函数
**影响文件**: like_api.go, follow_api.go, comment_api.go

**当前代码** (like_api.go):
```go
bookID := c.Param("bookId")
if bookID == "" {
    response.BadRequest(c, "参数错误", "书籍ID不能为空")
    return
}

userID, exists := c.Get("user_id")
if !exists {
    response.Unauthorized(c, "未授权")
    return
}
```

**优化后**:
```go
bookID, ok := shared.GetRequiredParam(c, "bookId", "书籍ID")
if !ok { return }

userID, ok := shared.GetUserID(c)
if !ok { return }
```

**收益**: 减少约15行重复代码/文件

#### 问题2：分页参数不统一
**当前代码** (like_api.go, follow_api.go):
```go
var params struct {
    Page int `form:"page" binding:"min=1"`
    Size int `form:"size" binding:"min=1,max=100"`
}
params.Page = 1
params.Size = 20

if err := c.ShouldBindQuery(&params); err != nil {
    response.BadRequest(c, "参数错误", err.Error())
    return
}
```

**优化后**:
```go
params := shared.GetPaginationParamsStandard(c)
```

**收益**: 减少10行重复代码/文件

#### 问题3：错误处理不统一
- like_api.go: 直接 response.InternalError
- collection_api.go: 根据错误类型返回不同响应
- follow_api.go: 字符串匹配错误类型

**优化方案**: 统一使用 response.HandleServiceError 或 pkg/errors 标准错误

---

## 2. 优化目标

### 2.1 代码质量优化 (P0 - 必须完成)

| 任务 | 目标 | 预计收益 |
|------|------|---------|
| 应用GetRequiredParam | 8个文件 | ~120行代码 |
| 应用GetUserID | 5个文件 | ~40行代码 |
| 统一分页参数 | 6个文件 | ~60行代码 |
| 统一响应格式 | 8个文件 | ~50行代码 |
| **总计** | - | **~270行代码** |

### 2.2 简单功能添加 (P1 - 尽量完成)

| 功能 | 复杂度 | 预计时间 | 优先级 |
|------|--------|----------|--------|
| 批量操作API | 简单 | 2h | 高 |
| 收藏分享优化 | 简单 | 1h | 中 |
| 统计信息缓存 | 中等 | 3h | 中 |

### 2.3 复杂功能标记TODO (P2 - 后续处理)

| 功能 | 复杂度 | 预计时间 | 依赖 |
|------|--------|----------|------|
| 用户标签系统 | 复杂 | 8h | 数据库设计 |
| 用户分组功能 | 复杂 | 12h | 关系重构 |
| 热门内容缓存 | 复杂 | 6h | Redis架构 |
| 性能监控 | 复杂 | 8h | 监控系统 |

---

## 3. 实施计划

### 3.1 Part 1: 代码质量优化 (2-3h)

#### 任务1.1: 优化 like_api.go
- [ ] 应用 GetRequiredParam
- [ ] 应用 GetUserID
- [ ] 统一分页参数
- [ ] 运行测试验证

#### 任务1.2: 优化 follow_api.go
- [ ] 应用 GetRequiredParam
- [ ] 应用 GetUserID
- [ ] 统一分页参数
- [ ] 运行测试验证

#### 任务1.3: 优化 comment_api.go
- [ ] 检查当前状态
- [ ] 应用辅助函数
- [ ] 运行测试验证

#### 任务1.4: 优化其他API文件
- [ ] rating_api.go
- [ ] review_api.go
- [ ] booklist_api.go
- [ ] message_api.go
- [ ] relation_api.go

#### 任务1.5: 全量测试
- [ ] 运行所有social模块测试
- [ ] 修复发现的问题

### 3.2 Part 2: 简单功能添加 (3-4h)

#### 任务2.1: 批量操作API
- [ ] 批量点赞
- [ ] 批量收藏
- [ ] 批量关注
- [ ] 编写测试

#### 任务2.2: 收藏分享优化
- [ ] 生成分享链接
- [ ] 访问统计
- [ ] 编写测试

#### 任务2.3: 统计信息缓存
- [ ] 缓存点赞数
- [ ] 缓存收藏数
- [ ] 缓存关注数

### 3.3 Part 3: 审查与修复 (1h)
- [ ] 代码审查
- [ ] 修复发现的问题
- [ ] 再次运行测试

### 3.4 Part 4: 提交 (0.5h)
- [ ] 创建PR
- [ ] 更新文档

---

## 4. 检查点设置

| 检查点 | 触发条件 | 检查内容 |
|--------|----------|----------|
| CP1 | 每个API优化后 | 测试通过 |
| CP2 | Part 1完成后 | 全量测试通过 |
| CP3 | Part 2完成后 | 新功能测试通过 |
| CP4 | 提交前 | 代码审查通过 |

---

## 5. TDD要求

### 5.1 测试优先原则
- 先写测试，观察失败
- 编写最小代码使测试通过
- 重构优化
- 确保测试持续通过

### 5.2 测试覆盖要求
- 现有功能: 保持测试通过
- 新功能: 100%测试覆盖
- 边界情况: 必须覆盖

---

## 6. 验收标准

### 6.1 功能验收
- [ ] 所有API使用统一辅助函数
- [ ] 所有测试通过
- [ ] 代码重复减少 > 200行

### 6.2 质量验收
- [ ] 代码审查通过
- [ ] 测试覆盖率不降低
- [ ] 无新增警告

### 6.3 文档验收
- [ ] 更新phase4_todo.md
- [ ] PR描述清晰

---

## 7. 风险与缓解

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| 测试失败 | 中 | 中 | 保持TDD，逐个文件优化 |
| 辅助函数不兼容 | 低 | 中 | 先验证单个文件 |
| 回归问题 | 低 | 高 | 全量测试覆盖 |

---

**文档状态**: ✅ 设计完成
**下一步**: 创建worktree并开始实施
