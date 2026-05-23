# 测试TODO快速参考卡

**快速参考** | **版本**: v1.0 | **更新**: 2025-10-27

> 📖 **详细文档**: 参见 [测试TODO功能实施指南](./测试TODO功能实施指南.md) 和 [测试覆盖率追踪报告](./测试覆盖率追踪报告.md)

---

## 🎯 当前状态一览

```
总体覆盖率: 45% / 目标: 90%
进度: ████████░░░░░░░░░░░░ 50%

集成测试: 12通过 / 2失败 / 4跳过 (共18个)
```

---

## 📋 核心TODO清单

### 🔥 P0 高优先级（本周必须完成）

#### 1. 评论系统 [0%]
- [ ] Repository层测试 (22个用例)
  - 文件: `test/repository/comment_repository_test.go`
  - 关键: CRUD、审核、统计、并发测试
- [ ] Service层测试 (27个用例)
  - 文件: `test/service/comment_service_test.go`
  - 关键: 业务逻辑、权限、事件发布
- [ ] API层测试 (18个用例)
  - 文件: `test/api/comment_api_test.go`
  - 关键: 端点测试、参数验证、权限控制
- [ ] 集成测试修复
  - `TestInteractionScenario/4.评论_发表书籍评论` ⏭️ 跳过
  - `TestInteractionScenario/5.评论_获取书籍评论列表` ⏭️ 跳过

**预计**: 4天 | **目标覆盖率**: 85%+

#### 2. 点赞系统 [0%]
- [ ] Repository层测试 (14个用例)
  - 文件: `test/repository/like_repository_test.go`
  - 关键: 添加/取消点赞、防重、并发测试
- [ ] Service层测试 (13个用例)
  - 文件: `test/service/like_service_test.go`
  - 关键: 点赞书籍/评论、防刷机制
- [ ] API层测试 (9个用例)
  - 文件: `test/api/like_api_test.go`
  - 关键: 点赞端点、取消点赞、查询状态
- [ ] 集成测试修复
  - `TestInteractionScenario/6.点赞_点赞书籍` ⏭️ 跳过
  - `TestInteractionScenario/7.点赞_取消点赞` ⏭️ 跳过

**预计**: 2天 | **目标覆盖率**: 85%+

#### 3. 书籍详情API修复 [0%]
- [ ] 修复路由问题
  - 问题: GET /api/v1/bookstore/books/:id 返回404
  - 检查: `router/bookstore/book_routes.go`
  - 检查: `api/v1/bookstore/book_api.go`
- [ ] 添加测试
  - 有效ID、无效ID、不存在ID测试
- [ ] 修复集成测试
  - `TestReadingScenario/1.书籍详情_获取书籍信息` ❌ 失败

**预计**: 0.5天

---

### ⚠️ P1 中优先级（下周完成）

#### 4. 章节列表API修复 [0%]
- [ ] 修复返回格式问题
  - 问题: GET /api/v1/reader/chapters 返回HTML而非JSON
  - 检查: `router/reader/chapter_routes.go`
  - 检查: `api/v1/reader/chapter_api.go`
- [ ] 修复集成测试
  - `TestReadingScenario/2.书籍详情_获取章节列表` ❌ 失败

**预计**: 0.5天

#### 5. 独立收藏系统 [0%]
- [ ] Repository层测试 (15个用例)
- [ ] Service层测试 (15个用例)
- [ ] API层测试 (8个用例)
- [ ] 集成测试创建

**预计**: 3天 | **目标覆盖率**: 85%+

#### 6. 独立阅读历史系统 [40%]
- [ ] Repository层测试 (11个用例)
- [ ] Service层测试 (11个用例)
- [ ] API层测试 (6个用例)
- [ ] 集成测试修复
  - `TestInteractionScenario/8.阅读历史_查看阅读历史` ⏭️ 跳过

**预计**: 2天 | **目标覆盖率**: 85%+

---

## 📅 3周时间表

### 第1周 (10/28-11/01) - 评论和点赞

| 日期 | 任务 | 产出 |
|------|------|------|
| 周一 | 评论Repository层 | 15个用例，60%覆盖率 |
| 周二 | 评论Repository完成 | 22个用例，85%+覆盖率 |
| 周三 | 评论Service层 | 15个用例，60%覆盖率 |
| 周四 | 评论完成+API层 | 45个用例，85%+覆盖率 |
| 周五 | 点赞系统完整测试 | 36个用例，85%+覆盖率 |

**里程碑**: ✅ 评论和点赞上线，4个集成测试通过

### 第2周 (11/04-11/08) - API修复和收藏

| 日期 | 任务 | 产出 |
|------|------|------|
| 周一 | 书籍详情API修复 | API测试通过 |
| 周二 | 章节列表API修复 | API测试通过 |
| 周三-四 | 收藏系统实现 | 38个用例，85%+覆盖率 |
| 周五 | 集成测试和回归 | 所有集成测试通过 |

**里程碑**: ✅ API修复，收藏系统上线，2个失败测试修复

### 第3周 (11/11-11/15) - 历史和验收

| 日期 | 任务 | 产出 |
|------|------|------|
| 周一-二 | 阅读历史测试 | 28个用例，85%+覆盖率 |
| 周三 | 集成测试完善 | 所有集成测试通过 |
| 周四 | 补充遗漏测试 | 总体覆盖率90%+ |
| 周五 | 最终验收 | 项目验收通过 ✅ |

**里程碑**: ✅ 全部功能上线，总体覆盖率90%+

---

## 🚀 快速开始

### 环境准备

```bash
# 1. 启动测试环境
docker-compose -f docker/docker-compose.test.yml up -d

# 2. 配置测试环境
cp config/config.test.yaml.example config/config.test.yaml

# 3. 初始化测试数据
go run cmd/prepare_test_data/main.go
```

### 运行测试

```bash
# 运行所有测试
make test

# 运行单元测试
make test-unit

# 运行集成测试
make test-integration

# 生成覆盖率报告
make test-coverage

# 运行特定测试
go test ./test/repository/comment_repository_test.go -v
```

### 查看报告

```bash
# 查看覆盖率报告
open test_results/coverage.html

# 查看测试结果
cat test_results/test_results.txt
```

---

## ✅ 测试模板

### Repository层测试模板

```go
package repository

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "Qingyu_backend/test/testutil"
)

func TestCreateComment(t *testing.T) {
    // 准备测试环境
    testEnv := testutil.SetupTestEnvironment(t)
    defer testEnv.Cleanup()
    
    repo := testEnv.CommentRepository
    
    // 准备测试数据
    comment := &models.Comment{
        BookID:  "book_123",
        UserID:  "user_123",
        Content: "测试评论",
        Rating:  5,
    }
    
    // 执行测试
    err := repo.CreateComment(context.Background(), comment)
    
    // 验证结果
    assert.NoError(t, err)
    assert.NotEmpty(t, comment.ID)
}
```

### Service层测试模板

```go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestPublishComment(t *testing.T) {
    // 创建Mock
    mockRepo := new(MockCommentRepository)
    mockEventBus := new(MockEventBus)
    
    service := NewCommentService(mockRepo, mockEventBus)
    
    // 设置Mock期望
    mockRepo.On("CreateComment", mock.Anything, mock.Anything).Return(nil)
    mockEventBus.On("Publish", mock.Anything).Return(nil)
    
    // 执行测试
    req := &CommentRequest{
        BookID:  "book_123",
        Content: "测试评论",
        Rating:  5,
    }
    err := service.PublishComment(context.Background(), req)
    
    // 验证结果
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
    mockEventBus.AssertExpectations(t)
}
```

### API层测试模板

```go
package api

import (
    "testing"
    "net/http/httptest"
    "github.com/stretchr/testify/assert"
)

func TestPostComment(t *testing.T) {
    // 准备测试环境
    router := setupTestRouter()
    
    // 准备请求
    reqBody := `{"book_id":"123","content":"测试评论","rating":5}`
    req := httptest.NewRequest("POST", "/api/v1/reader/comments", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+testToken)
    
    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // 验证响应
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "success", response["status"])
}
```

---

## 📊 覆盖率指标

### 目标覆盖率

| 层级 | 目标 | 当前 | 状态 |
|-----|------|------|------|
| Repository | 85% | 75% | 🟡 |
| Service | 85% | 42% | 🔴 |
| API | 80% | 38% | 🔴 |
| 集成测试 | 100% | 67% | 🟡 |
| **总体** | **90%** | **45%** | 🔴 |

### 每周目标

- **第1周末**: 总体覆盖率 ≥ 60%
- **第2周末**: 总体覆盖率 ≥ 75%
- **第3周末**: 总体覆盖率 ≥ 90% ✅

---

## 🚨 关键注意事项

### ⚠️ 高风险项

1. **评论系统复杂度高** - 67个测试用例，需要4天时间
2. **API路由问题** - 原因不明，可能需要额外调试时间
3. **并发测试** - 点赞和评论需要正确处理并发场景

### ✅ 最佳实践

1. **测试隔离** - 每个测试使用独立数据，使用`t.Cleanup`清理
2. **并行测试** - 使用`t.Parallel()`加速测试执行
3. **Mock使用** - Service层使用Mock Repository和EventBus
4. **详细日志** - 使用`t.Logf`输出调试信息
5. **覆盖率追踪** - 每天更新覆盖率报告

---

## 🔍 常见问题

### Q: 如何运行单个测试？

```bash
go test ./test/repository/comment_repository_test.go -v -run TestCreateComment
```

### Q: 如何查看测试覆盖率？

```bash
go test ./test/repository/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Q: 如何调试失败的测试？

1. 添加`t.Logf`输出调试信息
2. 使用`-v`参数运行测试查看详细输出
3. 检查测试数据是否正确
4. 使用Postman手动测试API

### Q: 如何处理测试数据隔离？

使用`testutil.SetupTestEnvironment(t)`和`t.Cleanup()`确保测试隔离和清理。

---

## 📚 相关文档

### 必读文档
- [测试TODO功能实施指南](./测试TODO功能实施指南.md) - 详细任务清单
- [测试覆盖率追踪报告](./测试覆盖率追踪报告.md) - 详细覆盖率统计
- [测试最佳实践](./测试最佳实践.md) - 测试规范和技巧

### 参考文档
- [集成测试使用指南](./集成测试使用指南.md)
- [测试架构设计规范](../standards/testing/测试架构设计规范.md)
- [详细实施指南](./测试TODO功能实施指南.md)

---

## 📞 获取帮助

**遇到问题？**
1. 查看 [测试最佳实践](./测试最佳实践.md)
2. 查看 [测试TODO功能实施指南](./测试TODO功能实施指南.md)
3. 在项目仓库创建Issue，标签: `testing`, `help-wanted`

---

## 📌 本周重点（2025-10-28 ~ 11-01）

### 必须完成
- [x] 评论系统Repository层测试 (22个用例)
- [x] 评论系统Service层测试 (27个用例)
- [x] 评论系统API层测试 (18个用例)
- [x] 点赞系统完整测试 (36个用例)
- [x] 修复4个跳过的集成测试

### 目标指标
- [x] 评论系统覆盖率 ≥ 85%
- [x] 点赞系统覆盖率 ≥ 85%
- [x] 总体覆盖率 ≥ 60%
- [x] 集成测试通过率 ≥ 80%

---

**版本**: v1.0  
**最后更新**: 2025-10-27  
**下次更新**: 每周五

**当前状态**: 🔴 待开始  
**下一步**: 立即开始评论系统Repository层测试 🚀

