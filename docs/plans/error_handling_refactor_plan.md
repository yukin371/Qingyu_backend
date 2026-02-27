# API层错误处理重构实施计划

## 📋 项目概述

**目标**: 统一API层错误处理方式，减少代码冗余40-60%

**分支**: `feature/error-handling-refactor`

**预计时间**: 2-3天（渐进式重构）

---

## 🎯 重构原则

### 1. 不改变的功能
- API接口路径
- 请求/响应格式
- 业务逻辑

### 2. 只改变的内容
- 错误处理方式（从手动改为中间件统一处理）
- 参数绑定方式（使用结构体+验证标签）
- 响应构造方式（使用辅助函数）

### 3. 保持兼容性
- 现有`shared.Success`等函数保持不变（添加新的辅助函数）
- 新旧方式共存，逐步迁移

---

## 📊 重构范围

### 优先级分类

| 优先级 | 模块 | 文件数 | 说明 |
|--------|------|--------|------|
| P0 | reader | 5 | 核心功能，使用频繁 |
| P1 | bookstore | 4 | 书城相关 |
| P1 | social | 3 | 社交功能 |
| P2 | writer | 3 | 作者端 |
| P2 | user | 3 | 用户相关 |
| P3 | admin | 5 | 管理员 |
| P3 | ai | 4 | AI功能 |

**总计**: 约27个API文件

---

## 🔄 重构模式

### 当前代码模式
```go
func (api *API) GetResource(c *gin.Context) {
    // 1. 手动参数绑定
    id, ok := shared.GetRequiredParam(c, "id", "ID")
    if !ok { return }

    // 2. 手动错误处理
    result, err := api.service.Get(id)
    if err != nil {
        if err == ErrNotFound {
            response.NotFound(c, "不存在")
            return
        }
        response.InternalError(c, err)
        return
    }

    // 3. 手动响应
    response.Success(c, result)
}
```

### 重构后代码模式
```go
func (api *API) GetResource(c *gin.Context) {
    // 1. 结构体参数绑定（自动验证）
    var params struct {
        ID string `uri:"id" binding:"required"`
    }
    if !helpers.BindParams(c, &params) {
        return
    }

    // 2. 简化的错误处理（交给中间件）
    result, err := api.service.Get(params.ID)
    if err != nil {
        c.Error(err)
        return
    }

    // 3. 简化的响应
    helpers.Success(c, result)
}
```

---

## 📅 实施步骤

### Phase 1: 准备工作（已完成）

- [x] 创建worktree
- [x] 创建辅助函数文件
- [x] 创建演示文档

### Phase 2: Reader模块重构（P0）

**目标**: 完成reader模块的重构作为模板

#### 2.1 创建Reader模块参数结构体

```go
// api/v1/reader/reader_params.go
package reader

type GetChapterParams struct {
    BookID    string `uri:"bookId" binding:"required"`
    ChapterID string `uri:"chapterId" binding:"required"`
}

type GetChapterByNumberParams struct {
    BookID     string `uri:"bookId" binding:"required"`
    ChapterNum int    `uri:"chapterNum" binding:"required,min=1"`
}

type SaveProgressParams struct {
    BookID    string  `form:"bookId" binding:"required"`
    ChapterID string  `form:"chapterId" binding:"required"`
    Progress  float64 `form:"progress" binding:"required,min=0,max=1"`
}
```

#### 2.2 重构ChapterAPI

**检查点1**: 重构后代码能编译通过
**检查点2**: 现有测试通过
**检查点3**: 手动测试章节阅读功能

#### 2.3 重构ProgressAPI

**检查点4**: 阅读进度保存正常工作

### Phase 3: Bookstore模块重构（P1）

#### 3.1 创建Bookstore参数结构体

```go
// api/v1/bookstore/bookstore_params.go
package bookstore

type GetBookDetailParams struct {
    ID string `uri:"id" binding:"required"`
}

type SearchBooksParams struct {
    Keyword   string `form:"keyword"`
    CategoryID string `form:"categoryId"`
    Author    string `form:"author"`
    Status    string `form:"status"`
    Page      int    `form:"page" binding:"min=1"`
    Size      int    `form:"size" binding:"min=1,max=100"`
}
```

#### 3.2 重构BookstoreAPI

**检查点5**: 搜索功能正常
**检查点6**: 详情页正常

### Phase 4: Social模块重构（P1）

**检查点7**: 评论发表正常
**检查点8**: 收藏功能正常

### Phase 5: 其他模块重构（P2-P3）

按优先级逐个模块进行

---

## 🧪 测试策略

### 1. 单元测试

确保每个重构的API函数有对应的测试：

```go
func TestChapterAPI_GetChapter(t *testing.T) {
    // Arrange
    api := setupTestAPI()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/reader/books/book1/chapters/chap1", nil)

    // Act
    api.router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, 200, w.Code)
    // ...
}
```

### 2. 集成测试

运行现有的集成测试确保功能不变：

```bash
# Reader模块测试
go test ./api/v1/reader/... -v

# Bookstore模块测试
go test ./api/v1/bookstore/... -v
```

### 3. 手动测试清单

- [ ] 章节阅读（正常/未登录/VIP章节）
- [ ] 阅读进度保存
- [ ] 书籍搜索
- [ ] 书籍详情
- [ ] 评论发表
- [ ] 收藏/取消收藏

---

## ⚠️ 风险控制

### 风险1: 破坏现有功能

**缓解措施**:
- 每个模块重构后立即运行测试
- 使用feature flag控制新旧方式切换
- 保留旧函数作为fallback

### 风险2: Service层错误类型不匹配

**缓解措施**:
- 检查Service层是否返回`UnifiedError`
- 如不匹配，创建适配器函数

### 风险3: 响应格式变化

**缓解措施**:
- 对比重构前后响应格式
- 使用API测试验证响应一致性

---

## 📝 每日检查清单

### 开始工作前
- [ ] 切换到worktree: `cd .claude/worktrees/error-handling-refactor/backend`
- [ ] 拉取最新代码: `git pull origin feature/api-refactor-phase5-ai-grpc`
- [ ] 检查当前分支: `git branch`

### 提交代码前
- [ ] 运行模块测试: `go test ./api/v1/[module]/... -v`
- [ ] 检查代码格式: `gofmt -l .`
- [ ] 确认编译通过: `go build ./...`

### 提交时
- [ ] 使用清晰的commit message:
  ```
  refactor(reader): 简化ChapterAPI错误处理

  - 使用helpers.BindParams替代手动参数绑定
  - 使用c.Error()替代手动错误处理
  - 减少约50%代码量
  ```

---

## 🎬 合并计划

### 前置条件
1. 所有模块重构完成
2. 测试覆盖率不低于重构前
3. 代码审查通过

### 合并步骤
1. 在worktree中完成最终测试
2. 推送到远程分支: `git push origin feature/error-handling-refactor`
3. 创建Pull Request到`feature/api-refactor-phase5-ai-grpc`
4. 代码审查
5. 合并后删除worktree

---

## 📈 进度跟踪

| 模块 | 状态 | 完成日期 | 备注 |
|------|------|----------|------|
| 准备工作 | ✅ | 2026-02-27 | worktree已创建 |
| reader | 🔄 | 待开始 | 优先级最高 |
| bookstore | ⏳ | 待开始 | - |
| social | ⏳ | 待开始 | - |
| writer | ⏳ | 待开始 | - |
| user | ⏳ | 待开始 | - |
| admin | ⏳ | 待开始 | - |
| ai | ⏳ | 待开始 | - |

**图例**: ✅ 已完成 | 🔄 进行中 | ⏳ 待开始 | ❌ 有问题

---

## 🔗 相关文档

- [API简化演示](./api_simplification_demo.md)
- [错误系统文档](../pkg/errors/README.md)
- [现有response函数](../api/v1/shared/response.go)

---

## 💡 最佳实践

### DO (推荐)
✅ 每次只重构一个模块
✅ 重构后立即测试
✅ 保留原有函数作为过渡
✅ 编写清晰的commit message

### DON'T (不推荐)
❌ 一次性重构所有模块
❌ 跳过测试直接提交
❌ 删除旧函数
❌ 混合重构和功能开发

---

*文档创建日期: 2026-02-27*
*最后更新: 2026-02-27*
