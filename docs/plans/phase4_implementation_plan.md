# 第4阶段：实施计划 - 子代理任务清单

> **创建日期**: 2026-02-27
> **执行方式**: TDD子代理驱动
> **工作区**: 新worktree

---

## 工作区设置

```bash
# 创建worktree
git worktree add ../Qingyu_backend-phase4 feature/api-refactor-phase4-social

# 切换到worktree
cd ../Qingyu_backend-phase4
```

---

## Part 1: 代码质量优化任务

### Task 1.1: 优化 like_api.go

**目标**: 应用辅助函数，减少重复代码

**具体要求**:
1. 将 `c.Param("bookId")` + 空值检查 替换为 `shared.GetRequiredParam(c, "bookId", "书籍ID")`
2. 将 `c.Get("user_id")` + exists检查 替换为 `shared.GetUserID(c)`
3. 将手动分页结构体替换为 `shared.GetPaginationParamsStandard(c)`
4. 保持所有测试通过

**检查点**:
- [ ] like_api_test.go 所有测试通过
- [ ] 代码行数减少 > 30行

**预计时间**: 30分钟

---

### Task 1.2: 优化 follow_api.go

**目标**: 应用辅助函数，减少重复代码

**具体要求**:
1. 将路径参数获取改为 `shared.GetRequiredParam`
2. 将用户ID获取改为 `shared.GetUserID`
3. 统一分页参数处理
4. 保持所有测试通过

**检查点**:
- [ ] follow_api_test.go 所有测试通过
- [ ] 代码行数减少 > 40行

**预计时间**: 30分钟

---

### Task 1.3: 检查优化 comment_api.go

**目标**: 检查并应用辅助函数

**具体要求**:
1. 读取文件分析当前代码
2. 应用必要的辅助函数
3. 保持所有测试通过

**检查点**:
- [ ] comment_api_test.go 所有测试通过
- [ ] 代码质量提升

**预计时间**: 20分钟

---

### Task 1.4: 检查优化 rating_api.go

**目标**: 检查并应用辅助函数

**具体要求**:
1. 读取文件分析当前代码
2. 应用必要的辅助函数
3. 保持所有测试通过

**检查点**:
- [ ] rating_api_test.go 所有测试通过

**预计时间**: 20分钟

---

### Task 1.5: 检查优化 review_api.go

**目标**: 检查并应用辅助函数

**具体要求**:
1. 读取文件分析当前代码
2. 应用必要的辅助函数
3. 保持所有测试通过

**检查点**:
- [ ] review_api_test.go 所有测试通过

**预计时间**: 20分钟

---

### Task 1.6: 检查优化 booklist_api.go

**目标**: 检查并应用辅助函数

**具体要求**:
1. 读取文件分析当前代码
2. 应用必要的辅助函数
3. 保持所有测试通过

**检查点**:
- [ ] booklist_api_test.go 所有测试通过

**预计时间**: 20分钟

---

### Task 1.7: 全量测试验证

**目标**: 确保所有优化后测试通过

**具体要求**:
1. 运行 `go test ./api/v1/social/... -v`
2. 运行 `go test ./service/social/... -v`
3. 修复发现的所有问题

**检查点**:
- [ ] 所有测试通过
- [ ] 无新增警告

**预计时间**: 20分钟

---

## Part 2: 简单功能添加任务

### Task 2.1: 添加批量点赞API

**目标**: 支持批量点赞书籍

**TDD步骤**:
1. RED: 编写测试 `TestLikeAPI_BatchLikeBooks`
2. GREEN: 实现 `BatchLikeBooks` 方法
3. REFACTOR: 优化代码

**签名**:
```go
// POST /api/v1/reader/books/batch-like
func (api *LikeAPI) BatchLikeBooks(c *gin.Context)
```

**请求体**:
```go
type BatchLikeBooksRequest struct {
    BookIDs []string `json:"book_ids" binding:"required,min=1,max=50"`
}
```

**检查点**:
- [ ] 测试覆盖成功/失败场景
- [ ] 批量限制验证

**预计时间**: 1小时

---

### Task 2.2: 添加批量收藏API

**目标**: 支持批量添加收藏

**TDD步骤**:
1. RED: 编写测试
2. GREEN: 实现功能
3. REFACTOR: 优化代码

**签名**:
```go
// POST /api/v1/reader/collections/batch
func (api *CollectionAPI) BatchAddCollections(c *gin.Context)
```

**检查点**:
- [ ] 测试通过
- [ ] 与service层协调

**预计时间**: 1小时

---

### Task 2.3: 优化收藏分享功能

**目标**: 生成唯一分享链接

**具体要求**:
1. 为公开收藏生成唯一链接
2. 记录访问统计

**检查点**:
- [ ] 功能测试通过
- [ ] 链接唯一性验证

**预计时间**: 1小时

---

## Part 3: 代码审查任务

### Task 3.1: 代码审查

**审查内容**:
1. 辅助函数应用是否正确
2. 错误处理是否统一
3. 测试覆盖是否充分
4. 代码风格是否一致

**审查标准**:
- [ ] 符合项目规范
- [ ] 无明显bug
- [ ] 性能无明显下降

---

## TODO标记 (复杂功能)

以下功能标记为TODO，后续阶段处理：

### TODO-1: 用户标签系统
**复杂度**: 高
**依赖**: 数据库表设计
**预计**: 8小时
**标记位置**: `docs/plans/phase4_todo.md`

### TODO-2: 用户分组功能
**复杂度**: 高
**依赖**: 关系重构
**预计**: 12小时
**标记位置**: `docs/plans/phase4_todo.md`

### TODO-3: 热门内容缓存
**复杂度**: 中高
**依赖**: Redis架构
**预计**: 6小时
**标记位置**: `docs/plans/phase4_todo.md`

---

## 执行顺序

```
Part 1: 代码质量优化
  ├─ Task 1.1: like_api.go (30min)
  ├─ Task 1.2: follow_api.go (30min)
  ├─ Task 1.3: comment_api.go (20min)
  ├─ Task 1.4: rating_api.go (20min)
  ├─ Task 1.5: review_api.go (20min)
  ├─ Task 1.6: booklist_api.go (20min)
  └─ Task 1.7: 全量测试 (20min)
  ↓
Part 2: 简单功能添加
  ├─ Task 2.1: 批量点赞 (1h)
  ├─ Task 2.2: 批量收藏 (1h)
  └─ Task 2.3: 收藏分享优化 (1h)
  ↓
Part 3: 审查修复
  └─ Task 3.1: 代码审查
  ↓
Part 4: 提交
  └─ 创建PR
```

---

## 总预计时间

- Part 1: 2.5小时
- Part 2: 3小时
- Part 3: 1小时
- Part 4: 0.5小时

**总计**: 约7小时

---

**文档状态**: ✅ 计划完成
**下一步**: 创建worktree并派遣子代理
