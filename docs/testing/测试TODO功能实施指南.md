# 测试TODO功能实施指南

**文档版本**: v1.0  
**创建日期**: 2025-10-27  
**最后更新**: 2025-10-27  
**关联计划**: `doc/implementation/00进度指导/计划/2025-10-25测试TODO功能实施计划.md`

---

## 📋 文档说明

本文档是基于集成测试发现的功能缺失和路由问题制定的测试实施指南，为开发人员提供清晰的测试任务清单和测试覆盖率目标。

### 文档目的

- 🎯 明确各功能模块的测试需求
- 📊 跟踪测试覆盖率进度
- ✅ 提供测试实施的最佳实践
- 🔄 确保测试与功能开发同步

### 适用范围

- 后端开发人员（实施功能和单元测试）
- QA测试人员（执行集成测试和回归测试）
- 项目经理（跟踪测试进度）

---

## 📊 测试覆盖率概览

### 当前状态（2025-10-27）

| 测试类别 | 总数 | 通过 | 失败 | 跳过 | 覆盖率 | 目标覆盖率 |
|---------|------|------|------|------|--------|-----------|
| 集成测试 - 互动功能 | 10 | 6 | 0 | 4 | 60% | 100% |
| 集成测试 - 阅读流程 | 8 | 6 | 2 | 0 | 75% | 100% |
| 单元测试 - 评论系统 | 0 | 0 | 0 | 0 | 0% | 85%+ |
| 单元测试 - 点赞系统 | 0 | 0 | 0 | 0 | 0% | 85%+ |
| 单元测试 - 收藏系统 | 0 | 0 | 0 | 0 | 0% | 85%+ |
| 单元测试 - 阅读历史 | 0 | 0 | 0 | 0 | 0% | 85%+ |
| **总体覆盖率** | - | - | - | - | **45%** | **90%+** |

### 测试分布

```
总测试数: 18
├── 通过: 12 (67%)
├── 失败: 2 (11%)
├── 跳过: 4 (22%)
└── 待新增: ~50 (评论、点赞、收藏、历史单元测试)
```

### 优先级分布

- 🔥 **P0 高优先级**: 评论系统测试、点赞系统测试、API路由修复测试
- ⚠️ **P1 中优先级**: 收藏系统测试、阅读历史测试
- 📝 **P2 低优先级**: 性能测试、边界测试

---

## 🎯 测试TODO清单

### 阶段一：核心互动功能测试（P0 🔥）

#### 1.1 评论系统测试

**优先级**: P0 🔥  
**预计工作量**: 3-4天  
**目标覆盖率**: 85%+  
**关联文件**: 
- `test/repository/comment_repository_test.go`
- `test/service/comment_service_test.go`
- `test/api/comment_api_test.go`
- `test/integration/comment_integration_test.go`

##### Repository层测试清单

**测试文件**: `test/repository/comment_repository_test.go`

- [ ] **基础CRUD测试**
  - [ ] `TestCreateComment` - 测试创建评论
    - 验证评论内容保存
    - 验证时间戳设置
    - 验证默认状态为pending
  - [ ] `TestGetCommentByID` - 测试获取单条评论
    - 验证正确获取
    - 验证不存在的ID返回错误
  - [ ] `TestGetCommentsByBookID` - 测试获取书籍评论列表
    - 验证按bookId筛选
    - 验证分页功能
    - 验证按创建时间排序
  - [ ] `TestGetCommentsByUserID` - 测试获取用户评论历史
    - 验证按userId筛选
    - 验证分页功能
  - [ ] `TestGetRepliesByCommentID` - 测试获取评论回复
    - 验证按parentId筛选
    - 验证回复嵌套关系
  - [ ] `TestUpdateComment` - 测试更新评论
    - 验证内容更新
    - 验证更新时间自动设置
  - [ ] `TestDeleteComment` - 测试删除评论
    - 验证软删除（标记status为deleted）
    - 验证物理删除（可选）

- [ ] **审核功能测试**
  - [ ] `TestUpdateCommentStatus` - 测试更新审核状态
    - 验证状态从pending到approved
    - 验证状态从pending到rejected
    - 验证拒绝原因字段
  - [ ] `TestGetPendingComments` - 测试获取待审核评论
    - 验证只返回status=pending的评论
    - 验证分页和排序

- [ ] **统计功能测试**
  - [ ] `TestIncrementLikeCount` - 测试增加点赞数
    - 验证点赞数正确递增
    - 验证并发安全（使用MongoDB原子操作）
  - [ ] `TestIncrementReplyCount` - 测试增加回复数
    - 验证回复数正确递增
  - [ ] `TestGetBookRatingStats` - 测试获取书籍评分统计
    - 验证平均评分计算
    - 验证评分分布统计

- [ ] **边界和异常测试**
  - [ ] `TestCreateCommentWithInvalidData` - 测试无效数据
    - 空内容
    - 超长内容
    - 无效bookId
  - [ ] `TestConcurrentLikeIncrement` - 测试并发点赞
    - 验证原子操作正确性
    - 验证数据一致性

##### Service层测试清单

**测试文件**: `test/service/comment_service_test.go`

- [ ] **业务逻辑测试**
  - [ ] `TestPublishComment` - 测试发表评论
    - 验证内容长度限制（10-500字）
    - 验证评分范围（1-5星）
    - 验证敏感词过滤
    - 验证自动审核逻辑
    - 验证事件发布（CommentCreatedEvent）
  - [ ] `TestReplyComment` - 测试回复评论
    - 验证回复关系建立（parentId, rootId）
    - 验证回复数统计更新
    - 验证回复嵌套层级限制
  - [ ] `TestGetCommentList` - 测试获取评论列表
    - 验证排序（最新优先）
    - 验证排序（最热优先 - 按点赞数）
    - 验证过滤（只显示已审核评论）
    - 验证用户信息附加
  - [ ] `TestUpdateComment` - 测试编辑评论
    - 验证只能编辑自己的评论
    - 验证编辑时间窗口（15分钟内）
    - 验证编辑后重新审核
  - [ ] `TestDeleteComment` - 测试删除评论
    - 验证只能删除自己的评论
    - 验证管理员可删除任何评论
    - 验证删除后级联处理回复

- [ ] **审核功能测试**
  - [ ] `TestReviewComment` - 测试审核评论
    - 验证只有管理员可审核
    - 验证审核通过
    - 验证审核拒绝
    - 验证事件发布（CommentReviewedEvent）

- [ ] **点赞功能测试**
  - [ ] `TestLikeComment` - 测试点赞评论
    - 验证点赞成功
    - 验证防重复点赞
    - 验证点赞数更新
  - [ ] `TestUnlikeComment` - 测试取消点赞
    - 验证取消成功
    - 验证点赞数减少

- [ ] **统计功能测试**
  - [ ] `TestGetBookCommentStats` - 测试书籍评论统计
    - 验证总评论数
    - 验证平均评分
    - 验证评分分布
  - [ ] `TestGetUserCommentStats` - 测试用户评论统计
    - 验证用户总评论数
    - 验证用户获赞数

- [ ] **Mock测试**
  - 使用Mock Repository
  - 使用Mock EventBus
  - 验证依赖注入正确

##### API层测试清单

**测试文件**: `test/api/comment_api_test.go`

- [ ] **端点测试**
  - [ ] `TestPostComment` - POST /api/v1/reader/comments
    - 验证201响应
    - 验证请求参数绑定
    - 验证参数验证错误返回400
    - 验证认证检查
  - [ ] `TestGetComments` - GET /api/v1/reader/comments
    - 验证200响应
    - 验证分页参数
    - 验证排序参数
    - 验证返回数据格式
  - [ ] `TestGetCommentDetail` - GET /api/v1/reader/comments/:id
    - 验证200响应
    - 验证不存在返回404
  - [ ] `TestUpdateComment` - PUT /api/v1/reader/comments/:id
    - 验证200响应
    - 验证权限检查（只能编辑自己的）
    - 验证403错误
  - [ ] `TestDeleteComment` - DELETE /api/v1/reader/comments/:id
    - 验证204响应
    - 验证权限检查
  - [ ] `TestReplyComment` - POST /api/v1/reader/comments/:id/reply
    - 验证201响应
    - 验证回复关系建立
  - [ ] `TestLikeComment` - POST /api/v1/reader/comments/:id/like
    - 验证200响应
    - 验证防重复点赞
  - [ ] `TestUnlikeComment` - DELETE /api/v1/reader/comments/:id/like
    - 验证204响应

- [ ] **管理员API测试**
  - [ ] `TestGetPendingComments` - GET /api/v1/admin/comments/pending
    - 验证管理员权限
    - 验证返回待审核列表
  - [ ] `TestReviewComment` - POST /api/v1/admin/comments/:id/review
    - 验证管理员权限
    - 验证审核操作

##### 集成测试清单

**测试文件**: `test/integration/comment_integration_test.go`

- [ ] **端到端测试**
  - [ ] `TestCommentE2EScenario` - 完整评论流程
    - 用户发表评论
    - 管理员审核通过
    - 其他用户查看评论
    - 用户点赞评论
    - 用户回复评论
    - 用户删除自己的评论
  - [ ] **修复现有集成测试**
    - [ ] 修复 `TestInteractionScenario/4.评论_发表书籍评论`
    - [ ] 修复 `TestInteractionScenario/5.评论_获取书籍评论列表`

**预期成果**:
- ✅ 所有测试通过
- ✅ 测试覆盖率 ≥ 85%
- ✅ 集成测试不再跳过

---

#### 1.2 点赞系统测试

**优先级**: P0 🔥  
**预计工作量**: 2天  
**目标覆盖率**: 85%+  
**关联文件**: 
- `test/repository/like_repository_test.go`
- `test/service/like_service_test.go`
- `test/api/like_api_test.go`
- `test/integration/like_integration_test.go`

##### Repository层测试清单

**测试文件**: `test/repository/like_repository_test.go`

- [ ] **基础操作测试**
  - [ ] `TestAddLike` - 测试添加点赞
    - 验证点赞记录创建
    - 验证唯一索引（防重复点赞）
    - 验证时间戳设置
  - [ ] `TestRemoveLike` - 测试取消点赞
    - 验证点赞记录删除
    - 验证不存在时返回错误
  - [ ] `TestIsLiked` - 测试检查点赞状态
    - 验证已点赞返回true
    - 验证未点赞返回false
  - [ ] `TestGetUserLikes` - 测试获取用户点赞列表
    - 验证按targetType筛选
    - 验证分页功能
    - 验证按时间排序
  - [ ] `TestGetLikeCount` - 测试获取点赞数
    - 验证正确计数
    - 验证不同targetType独立计数

- [ ] **批量操作测试**
  - [ ] `TestGetLikesCountBatch` - 测试批量获取点赞数
    - 验证批量查询性能
    - 验证返回数据正确性
  - [ ] `TestGetUserLikeStatusBatch` - 测试批量检查点赞状态
    - 验证批量查询
    - 验证返回格式

- [ ] **并发和边界测试**
  - [ ] `TestConcurrentAddLike` - 测试并发点赞
    - 验证唯一索引防重
    - 验证数据一致性
  - [ ] `TestAddLikeWithInvalidData` - 测试无效数据
    - 空userId
    - 空targetId
    - 无效targetType

##### Service层测试清单

**测试文件**: `test/service/like_service_test.go`

- [ ] **业务逻辑测试**
  - [ ] `TestLikeBook` - 测试点赞书籍
    - 验证点赞成功
    - 验证防重复点赞
    - 验证书籍点赞数更新
    - 验证事件发布（BookLikedEvent）
  - [ ] `TestUnlikeBook` - 测试取消点赞书籍
    - 验证取消成功
    - 验证书籍点赞数减少
    - 验证事件发布（BookUnlikedEvent）
  - [ ] `TestLikeComment` - 测试点赞评论
    - 验证点赞成功
    - 验证评论点赞数更新
  - [ ] `TestUnlikeComment` - 测试取消点赞评论
    - 验证取消成功
    - 验证评论点赞数减少
  - [ ] `TestGetBookLikeCount` - 测试获取书籍点赞数
    - 验证正确返回
  - [ ] `TestGetUserLikeStatus` - 测试检查用户点赞状态
    - 验证已点赞状态
    - 验证未点赞状态
  - [ ] `TestGetUserLikedBooks` - 测试获取用户点赞的书籍
    - 验证列表返回
    - 验证分页功能

- [ ] **防刷机制测试**
  - [ ] `TestRateLimitLike` - 测试点赞频率限制
    - 验证1秒内不能重复点赞同一对象
    - 验证频率限制错误返回

- [ ] **Mock测试**
  - 使用Mock Repository
  - 使用Mock EventBus
  - 验证依赖注入

##### API层测试清单

**测试文件**: `test/api/like_api_test.go`

- [ ] **端点测试**
  - [ ] `TestLikeBook` - POST /api/v1/reader/books/:id/like
    - 验证200响应
    - 验证认证检查
    - 验证防重复点赞
  - [ ] `TestUnlikeBook` - DELETE /api/v1/reader/books/:id/like
    - 验证204响应
    - 验证未点赞时返回错误
  - [ ] `TestGetLikeStatus` - GET /api/v1/reader/books/:id/like/status
    - 验证返回点赞状态
    - 验证返回点赞数
  - [ ] `TestLikeComment` - POST /api/v1/reader/comments/:id/like
    - 验证200响应
  - [ ] `TestUnlikeComment` - DELETE /api/v1/reader/comments/:id/like
    - 验证204响应

- [ ] **参数验证测试**
  - [ ] 无效ID格式
  - [ ] 不存在的目标对象

##### 集成测试清单

**测试文件**: `test/integration/like_integration_test.go`

- [ ] **端到端测试**
  - [ ] `TestLikeE2EScenario` - 完整点赞流程
    - 用户点赞书籍
    - 检查点赞状态
    - 取消点赞
    - 再次检查状态
  - [ ] **修复现有集成测试**
    - [ ] 修复 `TestInteractionScenario/6.点赞_点赞书籍`
    - [ ] 修复 `TestInteractionScenario/7.点赞_取消点赞`

**预期成果**:
- ✅ 所有测试通过
- ✅ 测试覆盖率 ≥ 85%
- ✅ 集成测试不再跳过
- ✅ 并发测试通过

---

### 阶段二：API路由修复测试（P0 🔥）

#### 2.1 书籍详情API测试

**优先级**: P0 🔥  
**预计工作量**: 0.5天  
**关联文件**: 
- `test/integration/scenario_reading_test.go`
- `test/api/bookstore_api_test.go`

##### 测试清单

- [ ] **路由测试**
  - [ ] `TestBookDetailRoute` - 测试路由注册
    - 验证GET /api/v1/bookstore/books/:id路由存在
    - 验证路由参数正确解析
  - [ ] `TestGetBookByID` - 测试获取书籍详情
    - 验证200响应
    - 验证返回完整书籍信息
    - 验证返回JSON格式
  - [ ] `TestGetBookByInvalidID` - 测试无效ID
    - 验证400响应（无效ObjectID格式）
  - [ ] `TestGetBookByNonExistentID` - 测试不存在的ID
    - 验证404响应

- [ ] **修复集成测试**
  - [ ] 修复 `TestReadingScenario/1.书籍详情_获取书籍信息`

**预期成果**:
- ✅ API返回200
- ✅ ObjectID正确解析
- ✅ 集成测试通过

---

#### 2.2 章节列表API测试

**优先级**: P1  
**预计工作量**: 0.5天  
**关联文件**: 
- `test/integration/scenario_reading_test.go`
- `test/api/chapter_api_test.go`

##### 测试清单

- [ ] **路由测试**
  - [ ] `TestChapterListRoute` - 测试路由注册
    - 验证GET /api/v1/reader/chapters路由存在
    - 验证查询参数正确解析
  - [ ] `TestGetChapterList` - 测试获取章节列表
    - 验证200响应
    - 验证返回JSON格式（而非HTML）
    - 验证返回章节列表
  - [ ] `TestGetChapterListWithPagination` - 测试分页
    - 验证page和pageSize参数
    - 验证返回分页信息
  - [ ] `TestGetChapterListInvalidBookID` - 测试无效bookId
    - 验证400响应

- [ ] **修复集成测试**
  - [ ] 修复 `TestReadingScenario/2.书籍详情_获取章节列表`

**预期成果**:
- ✅ API返回JSON而非HTML
- ✅ 分页功能正常
- ✅ 集成测试通过

---

### 阶段三：功能完善测试（P1）

#### 3.1 独立收藏系统测试

**优先级**: P1  
**预计工作量**: 2天  
**目标覆盖率**: 85%+  
**关联文件**: 
- `test/repository/collection_repository_test.go`
- `test/service/collection_service_test.go`
- `test/api/collection_api_test.go`
- `test/integration/collection_integration_test.go`

##### Repository层测试清单

- [ ] **基础CRUD测试**
  - [ ] `TestAddCollection` - 测试添加收藏
  - [ ] `TestRemoveCollection` - 测试取消收藏
  - [ ] `TestGetCollectionsByUserID` - 测试获取用户收藏列表
  - [ ] `TestGetCollectionByID` - 测试获取单条收藏
  - [ ] `TestUpdateCollection` - 测试更新收藏（笔记、标签）

- [ ] **收藏夹管理测试**
  - [ ] `TestCreateFolder` - 测试创建收藏夹
  - [ ] `TestGetFoldersByUserID` - 测试获取收藏夹列表
  - [ ] `TestUpdateFolder` - 测试更新收藏夹
  - [ ] `TestDeleteFolder` - 测试删除收藏夹
  - [ ] `TestMoveCollectionToFolder` - 测试移动收藏到文件夹

- [ ] **查询测试**
  - [ ] `TestGetCollectionsByFolder` - 测试按文件夹筛选
  - [ ] `TestGetCollectionsByTag` - 测试按标签筛选
  - [ ] `TestSearchCollections` - 测试搜索收藏

##### Service层测试清单

- [ ] **业务逻辑测试**
  - [ ] `TestAddToCollection` - 测试添加收藏
    - 验证防重复收藏
    - 验证收藏数更新
  - [ ] `TestRemoveFromCollection` - 测试取消收藏
  - [ ] `TestUpdateCollectionNote` - 测试更新笔记
  - [ ] `TestAddCollectionTags` - 测试添加标签
  - [ ] `TestShareCollection` - 测试分享收藏
    - 验证公开/私有设置

- [ ] **收藏夹管理测试**
  - [ ] `TestCreateCollectionFolder` - 测试创建收藏夹
  - [ ] `TestRenameFolder` - 测试重命名收藏夹
  - [ ] `TestDeleteFolderWithCollections` - 测试删除包含收藏的文件夹
    - 验证级联处理

##### API层测试清单

- [ ] **端点测试**
  - [ ] POST /api/v1/reader/collections - 添加收藏
  - [ ] GET /api/v1/reader/collections - 获取收藏列表
  - [ ] GET /api/v1/reader/collections/:id - 获取收藏详情
  - [ ] PUT /api/v1/reader/collections/:id - 更新收藏
  - [ ] DELETE /api/v1/reader/collections/:id - 取消收藏
  - [ ] POST /api/v1/reader/collections/folders - 创建收藏夹
  - [ ] GET /api/v1/reader/collections/folders - 获取收藏夹列表

##### 集成测试清单

- [ ] **端到端测试**
  - [ ] `TestCollectionE2EScenario` - 完整收藏流程
    - 创建收藏夹
    - 添加收藏到文件夹
    - 添加标签和笔记
    - 移动收藏到其他文件夹
    - 分享收藏
    - 取消收藏

**预期成果**:
- ✅ 独立收藏系统测试完成
- ✅ 测试覆盖率 ≥ 85%
- ✅ 与书架系统区分清晰

---

#### 3.2 独立阅读历史系统测试

**优先级**: P1  
**预计工作量**: 1.5天  
**目标覆盖率**: 85%+  
**关联文件**: 
- `test/repository/history_repository_test.go`
- `test/service/history_service_test.go`
- `test/api/history_api_test.go`
- `test/integration/history_integration_test.go`

##### Repository层测试清单

- [ ] **基础CRUD测试**
  - [ ] `TestRecordReadingHistory` - 测试记录阅读历史
  - [ ] `TestGetReadingHistoryByUserID` - 测试获取用户历史
  - [ ] `TestGetReadingHistoryByID` - 测试获取单条历史
  - [ ] `TestDeleteReadingHistory` - 测试删除历史
  - [ ] `TestClearReadingHistory` - 测试清空历史

- [ ] **查询测试**
  - [ ] `TestGetReadingHistoryByBookID` - 测试按书籍筛选
  - [ ] `TestGetReadingHistoryByDateRange` - 测试按时间范围筛选
  - [ ] `TestGetReadingHistoryWithPagination` - 测试分页

- [ ] **统计测试**
  - [ ] `TestGetReadingStats` - 测试阅读统计
    - 总阅读时长
    - 阅读书籍数
    - 阅读章节数

##### Service层测试清单

- [ ] **业务逻辑测试**
  - [ ] `TestRecordReading` - 测试记录阅读
    - 验证自动创建历史记录
    - 验证阅读时长计算
    - 验证进度更新
  - [ ] `TestGetUserReadingHistory` - 测试获取历史
    - 验证按时间排序
    - 验证分页
  - [ ] `TestGetReadingStats` - 测试统计
    - 验证总时长统计
    - 验证日/周/月统计
  - [ ] `TestCleanupOldHistory` - 测试历史清理
    - 验证90天前记录自动清理

##### API层测试清单

- [ ] **端点测试**
  - [ ] GET /api/v1/reader/history - 获取阅读历史
  - [ ] GET /api/v1/reader/history/stats - 获取阅读统计
  - [ ] DELETE /api/v1/reader/history - 清空历史
  - [ ] DELETE /api/v1/reader/history/:id - 删除单条历史

##### 集成测试清单

- [ ] **端到端测试**
  - [ ] `TestReadingHistoryE2EScenario` - 完整历史流程
    - 用户阅读章节
    - 自动记录历史
    - 查看历史列表
    - 查看统计数据
    - 删除部分历史
    - 清空历史
  - [ ] **修复现有集成测试**
    - [ ] 修复 `TestInteractionScenario/8.阅读历史_查看阅读历史`

**预期成果**:
- ✅ 独立阅读历史系统测试完成
- ✅ 测试覆盖率 ≥ 85%
- ✅ 与阅读进度系统区分清晰
- ✅ 集成测试不再跳过

---

## 📋 测试实施最佳实践

### 1. 测试组织原则

#### 文件组织

```
test/
├── repository/              # Repository层单元测试
│   ├── comment_repository_test.go
│   ├── like_repository_test.go
│   ├── collection_repository_test.go
│   └── history_repository_test.go
├── service/                 # Service层单元测试
│   ├── comment_service_test.go
│   ├── like_service_test.go
│   ├── collection_service_test.go
│   └── history_service_test.go
├── api/                     # API层单元测试
│   ├── comment_api_test.go
│   ├── like_api_test.go
│   ├── collection_api_test.go
│   └── history_api_test.go
└── integration/             # 集成测试
    ├── comment_integration_test.go
    ├── like_integration_test.go
    ├── collection_integration_test.go
    ├── history_integration_test.go
    ├── scenario_interaction_test.go  # 现有互动场景测试
    └── scenario_reading_test.go      # 现有阅读场景测试
```

#### 测试命名规范

```go
// Repository层测试
func TestCreateComment(t *testing.T) { ... }
func TestGetCommentByID(t *testing.T) { ... }

// Service层测试
func TestPublishComment(t *testing.T) { ... }
func TestReviewComment(t *testing.T) { ... }

// API层测试
func TestPostComment(t *testing.T) { ... }
func TestGetComments(t *testing.T) { ... }

// 集成测试
func TestCommentE2EScenario(t *testing.T) { ... }
```

### 2. 测试数据准备

#### 使用测试工具

```go
// 使用现有的testutil包
import "Qingyu_backend/test/testutil"

func TestCreateComment(t *testing.T) {
    // 初始化测试环境
    testEnv := testutil.SetupTestEnvironment(t)
    defer testEnv.Cleanup()
    
    // 创建测试数据
    user := testEnv.CreateTestUser("testuser")
    book := testEnv.CreateTestBook("testbook")
    
    // 执行测试
    comment := &models.Comment{
        BookID:  book.ID.Hex(),
        UserID:  user.ID.Hex(),
        Content: "这是一条测试评论",
        Rating:  5,
    }
    
    err := repo.CreateComment(context.Background(), comment)
    assert.NoError(t, err)
    assert.NotEmpty(t, comment.ID)
}
```

#### 测试数据清理

```go
func TestWithCleanup(t *testing.T) {
    // 使用t.Cleanup确保测试后清理
    t.Cleanup(func() {
        // 清理测试数据
        testutil.CleanupTestData(t)
    })
    
    // 测试逻辑
}
```

### 3. Mock使用规范

#### Service层Mock Repository

```go
// 使用mockery生成Mock
type MockCommentRepository struct {
    mock.Mock
}

func (m *MockCommentRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
    args := m.Called(ctx, comment)
    return args.Error(0)
}

// 在测试中使用
func TestCommentService(t *testing.T) {
    mockRepo := new(MockCommentRepository)
    mockEventBus := new(MockEventBus)
    
    service := NewCommentService(mockRepo, mockEventBus)
    
    // 设置Mock期望
    mockRepo.On("CreateComment", mock.Anything, mock.Anything).Return(nil)
    mockEventBus.On("Publish", mock.Anything).Return(nil)
    
    // 执行测试
    err := service.PublishComment(context.Background(), &CommentRequest{...})
    
    // 验证
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
    mockEventBus.AssertExpectations(t)
}
```

### 4. 集成测试规范

#### 使用场景测试

```go
func TestCommentE2EScenario(t *testing.T) {
    // 子测试组织
    t.Run("发表评论", func(t *testing.T) {
        // 测试发表评论
    })
    
    t.Run("审核评论", func(t *testing.T) {
        // 测试审核流程
    })
    
    t.Run("点赞评论", func(t *testing.T) {
        // 测试点赞功能
    })
    
    t.Run("回复评论", func(t *testing.T) {
        // 测试回复功能
    })
    
    t.Run("删除评论", func(t *testing.T) {
        // 测试删除功能
    })
}
```

#### HTTP请求测试

```go
func TestCommentAPI(t *testing.T) {
    // 使用httptest
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

### 5. 并发测试

```go
func TestConcurrentLike(t *testing.T) {
    // 准备测试环境
    repo := setupTestRepository(t)
    bookID := "test_book_id"
    
    // 并发点赞
    const numGoroutines = 100
    var wg sync.WaitGroup
    wg.Add(numGoroutines)
    
    for i := 0; i < numGoroutines; i++ {
        go func(userID string) {
            defer wg.Done()
            err := repo.AddLike(context.Background(), userID, "book", bookID)
            assert.NoError(t, err)
        }(fmt.Sprintf("user_%d", i))
    }
    
    wg.Wait()
    
    // 验证点赞数
    count, err := repo.GetLikeCount(context.Background(), "book", bookID)
    assert.NoError(t, err)
    assert.Equal(t, numGoroutines, count)
}
```

### 6. 性能测试

```go
func BenchmarkCreateComment(b *testing.B) {
    repo := setupTestRepository(b)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        comment := &models.Comment{
            BookID:  "book_123",
            UserID:  "user_123",
            Content: "Benchmark test comment",
            Rating:  5,
        }
        repo.CreateComment(context.Background(), comment)
    }
}
```

---

## 📊 测试覆盖率追踪

### 使用go test工具

```bash
# 运行测试并生成覆盖率报告
go test ./test/repository/... -coverprofile=coverage_repo.out
go test ./test/service/... -coverprofile=coverage_service.out
go test ./test/api/... -coverprofile=coverage_api.out
go test ./test/integration/... -coverprofile=coverage_integration.out

# 查看覆盖率
go tool cover -func=coverage_repo.out
go tool cover -html=coverage_repo.out -o coverage_repo.html

# 合并覆盖率报告
gocovmerge coverage_*.out > coverage_total.out
go tool cover -func=coverage_total.out
```

### 覆盖率目标

| 测试层级 | 目标覆盖率 | 当前覆盖率 | 状态 |
|---------|-----------|-----------|------|
| Repository层 | 85% | 0% | ❌ 待实现 |
| Service层 | 85% | 0% | ❌ 待实现 |
| API层 | 80% | 0% | ❌ 待实现 |
| 集成测试 | 100% | 67% | 🟡 进行中 |
| **总体覆盖率** | **90%** | **45%** | 🟡 进行中 |

### 覆盖率提升计划

#### 第1周目标
- Repository层覆盖率达到 70%+
- Service层覆盖率达到 70%+
- 集成测试通过率达到 80%+

#### 第2周目标
- Repository层覆盖率达到 85%+
- Service层覆盖率达到 85%+
- API层覆盖率达到 70%+
- 集成测试通过率达到 90%+

#### 第3周目标
- API层覆盖率达到 80%+
- 集成测试通过率达到 100%
- **总体覆盖率达到 90%+**

---

## 🚀 快速开始

### 1. 环境准备

```bash
# 1. 确保MongoDB和Redis运行
docker-compose -f docker/docker-compose.test.yml up -d

# 2. 配置测试环境变量
cp config/config.test.yaml.example config/config.test.yaml

# 3. 初始化测试数据
go run cmd/prepare_test_data/main.go
```

### 2. 运行测试

```bash
# 运行所有测试
make test

# 运行单元测试
make test-unit

# 运行集成测试
make test-integration

# 运行特定测试
go test ./test/repository/comment_repository_test.go -v
go test ./test/integration/comment_integration_test.go -v

# 运行测试并生成覆盖率报告
make test-coverage
```

### 3. 查看测试报告

```bash
# 查看覆盖率报告
open test_results/coverage.html

# 查看测试结果
cat test_results/test_results.txt
```

---

## 📝 测试报告模板

### 单元测试报告

```markdown
## 评论系统单元测试报告

**测试日期**: 2025-11-01  
**测试人员**: [姓名]  
**测试范围**: Repository层、Service层、API层

### Repository层测试结果
- 测试用例数: 15
- 通过: 15
- 失败: 0
- 覆盖率: 88%

### Service层测试结果
- 测试用例数: 20
- 通过: 20
- 失败: 0
- 覆盖率: 92%

### API层测试结果
- 测试用例数: 12
- 通过: 12
- 失败: 0
- 覆盖率: 85%

### 问题记录
无

### 总结
评论系统单元测试全部通过，覆盖率达标。
```

### 集成测试报告

```markdown
## 评论系统集成测试报告

**测试日期**: 2025-11-01  
**测试人员**: [姓名]  
**测试场景**: 完整评论流程

### 测试场景
1. 用户发表评论 ✅
2. 管理员审核通过 ✅
3. 其他用户查看评论 ✅
4. 用户点赞评论 ✅
5. 用户回复评论 ✅
6. 用户删除评论 ✅

### 性能指标
- 发表评论响应时间: 150ms
- 获取评论列表响应时间: 180ms
- 点赞评论响应时间: 80ms

### 问题记录
无

### 总结
评论系统集成测试全部通过，性能指标符合要求。
```

---

## 🔍 常见问题

### Q1: 如何处理测试数据隔离？

**A**: 使用独立的测试数据库和测试数据，每个测试用例使用唯一的测试数据，测试后自动清理。

```go
func TestWithIsolation(t *testing.T) {
    // 创建唯一的测试数据
    testID := uuid.New().String()
    user := createTestUser(t, "user_"+testID)
    
    // 使用t.Cleanup确保清理
    t.Cleanup(func() {
        deleteTestUser(t, user.ID)
    })
    
    // 测试逻辑
}
```

### Q2: 如何加速测试执行？

**A**: 
1. 使用并行测试 `t.Parallel()`
2. Mock外部依赖（数据库、Redis、外部API）
3. 使用内存数据库（如SQLite）进行Repository测试
4. 合理使用测试缓存

```go
func TestParallel(t *testing.T) {
    t.Parallel() // 标记为可并行测试
    
    // 测试逻辑
}
```

### Q3: 如何测试异步操作（如EventBus）？

**A**: 使用Mock EventBus或等待机制验证事件发布。

```go
func TestEventPublish(t *testing.T) {
    mockEventBus := new(MockEventBus)
    service := NewCommentService(repo, mockEventBus)
    
    // 设置期望
    mockEventBus.On("Publish", mock.MatchedBy(func(event base.Event) bool {
        return event.GetEventType() == "CommentCreatedEvent"
    })).Return(nil)
    
    // 执行操作
    service.PublishComment(ctx, req)
    
    // 验证事件已发布
    mockEventBus.AssertExpectations(t)
}
```

### Q4: 集成测试失败如何调试？

**A**:
1. 添加详细日志
2. 使用`t.Logf`输出调试信息
3. 使用Postman或curl手动测试API
4. 检查测试数据是否正确

```go
func TestWithDebug(t *testing.T) {
    t.Logf("开始测试: %s", t.Name())
    
    resp, err := makeRequest(t, "POST", "/api/v1/comments", body)
    t.Logf("响应状态: %d", resp.StatusCode)
    t.Logf("响应内容: %s", resp.Body)
    
    // 断言
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

---

## 📚 参考资源

### 项目内文档
- `doc/testing/测试最佳实践.md` - 测试最佳实践指南
- `doc/testing/测试架构设计规范.md` - 测试架构设计
- `doc/testing/集成测试使用指南.md` - 集成测试指南
- `test/README.md` - 测试运行指南
- `doc/implementation/00进度指导/计划/2025-10-25测试TODO功能实施计划.md` - 详细实施计划

### Go测试官方文档
- [Go Testing](https://golang.org/pkg/testing/)
- [Go Test Coverage](https://blog.golang.org/cover)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)

### 测试工具
- [testify](https://github.com/stretchr/testify) - 断言和Mock库
- [mockery](https://github.com/vektra/mockery) - Mock生成工具
- [gocovmerge](https://github.com/wadey/gocovmerge) - 覆盖率合并工具

---

## 📅 测试进度跟踪

### 总体进度

```
进度: ████░░░░░░░░░░░░░░░░ 20% (0/50 完成)

阶段一: 评论系统和点赞系统测试 ░░░░░░░░░░ 0%
阶段二: API路由修复测试 ░░░░░░░░░░ 0%
阶段三: 收藏和历史系统测试 ░░░░░░░░░░ 0%
```

### 周进度更新

**本周完成** (2025-10-27 ~ 11-01):
- [ ] 评论系统Repository层测试
- [ ] 评论系统Service层测试
- [ ] 评论系统API层测试
- [ ] 点赞系统完整测试

**下周计划** (2025-11-04 ~ 11-08):
- [ ] API路由修复测试
- [ ] 收藏系统测试开始

---

## ✅ 验收标准

### 功能完整性
- [ ] 所有TODO功能已实现测试
- [ ] 所有集成测试通过（0失败，0跳过）
- [ ] 所有单元测试通过

### 测试覆盖率
- [ ] Repository层覆盖率 ≥ 85%
- [ ] Service层覆盖率 ≥ 85%
- [ ] API层覆盖率 ≥ 80%
- [ ] 总体覆盖率 ≥ 90%

### 代码质量
- [ ] 遵循测试最佳实践
- [ ] 测试代码清晰可维护
- [ ] 无linter错误
- [ ] 测试文档完整

### 性能指标
- [ ] 单元测试执行时间 < 30秒
- [ ] 集成测试执行时间 < 5分钟
- [ ] API响应时间符合要求

---

**文档维护者**: Qingyu后端测试团队  
**最后更新**: 2025-10-27  
**文档状态**: ✅ 完整

**下一步行动**:
1. ✅ 评审本测试指南
2. 🚀 开始阶段一测试实施（评论和点赞系统）
3. 📊 每周更新测试进度
4. 🎯 确保按时完成测试目标

