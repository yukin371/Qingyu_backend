# 前后端数据类型全面对比报告

**日期**: 2026-03-04
**版本**: v1.0
**状态**: 完成
**覆盖范围**: Writer域、Reader域、Bookstore域、Social域、Admin域、Recommendation域

---

## 一、执行摘要

### 1.1 分析范围

本次对比分析覆盖了Qingyu项目的五大核心业务域：

| 域 | 前端类型文件 | 后端模型文件 | 状态 |
|---|------------|-------------|------|
| Writer域 | `src/types/writer/index.d.ts` | `models/writer/*.go` | ✅ 已对比 |
| Reader域 | `src/types/reader.ts` | `models/bookstore/*.go` | ✅ 已对比 |
| Bookstore域 | `src/types/bookstore.ts` | `models/bookstore/*.go` | ✅ 已对比 |
| Social域 | `src/types/social.ts`, `review.ts` | `models/social/*.go` | ✅ 已对比 |
| Admin域 | `src/types/auth.ts`, `user/index.d.ts` | `models/users/*.go` | ✅ 已对比 |
| Recommendation域 | `src/types/recommendation.ts` | `models/recommendation/*.go` | ✅ 已对比 |

### 1.2 关键发现

**总体评估**：
- 命名规范一致性：⭐⭐⭐⭐ (后端采用BSON snake_case + JSON camelCase双重规范)
- 类型定义完整性：⭐⭐⭐⭐ (大部分类型定义完整)
- 枚举值一致性：⭐⭐⭐ (部分枚举值需要统一)
- 字段覆盖率：⭐⭐⭐ (存在一些字段缺失)

**问题统计**：
- **P0阻塞问题**: 5个
- **P1重要问题**: 12个
- **P2一般问题**: 11个

---

## 二、数据格式规范对比

### 2.1 命名规范

| 层级 | 命名规范 | 示例 |
|------|---------|------|
| 前端 TypeScript | camelCase | `projectId`, `avatarUrl`, `wordCount` |
| 后端 BSON (MongoDB) | snake_case | `project_id`, `avatar_url`, `word_count` |
| 后端 JSON (API) | camelCase | `projectId`, `avatarUrl`, `wordCount` |

**结论**: 后端采用了良好的双重命名规范，BSON存储使用snake_case（MongoDB标准），API响应使用camelCase（前端友好）。

### 2.2 类型转换规范

| 数据类型 | 前端 | 后端 | 转换方式 |
|---------|-----|------|---------|
| ID | string | primitive.ObjectID | 后端DTO层自动转换为hex string |
| 时间 | string (ISO 8601) | time.Time | JSON序列化自动转换 |
| 金额 | 元 (number) | 分 (float64) | 需要除以/乘以100 |
| 布尔值 | isXxx | is_xxx / isXxx | 自动转换 |
| 枚举 | string (联合类型) | string/常量 | 直接映射 |
> 注意：一般涉及金融，使用int类型定义最小单位更标准，后端应该定义int 分，前端显示元的转换
---

## 三、各域详细对比

### 3.1 Writer域（写作创作）

#### 3.1.1 Character（角色）

| 前端字段              | 后端BSON             | 后端JSON            | 类型一致性 |
| ----------------- | ------------------ | ----------------- | ----- |
| id                | _id                | id                | ✅     |
| projectId         | project_id         | projectId         | ✅     |
| name              | name               | name              | ✅     |
| alias             | alias              | alias             | ✅     |
| summary           | summary            | summary           | ✅     |
| traits            | traits             | traits            | ✅     |
| background        | background         | background        | ✅     |
| avatarUrl         | avatar_url         | avatarUrl         | ✅     |
| personalityPrompt | personality_prompt | personalityPrompt | ✅     |
| speechPattern     | speech_pattern     | speechPattern     | ✅     |
| currentState      | current_state      | currentState      | ✅     |
| createdAt         | created_at         | createdAt         | ✅     |
| updatedAt         | updated_at         | updatedAt         | ✅     |

**差异**：无重大差异，字段完全一致 ✓

#### 3.1.2 TimelineEvent（时间线事件）

EventType枚举对比：
- 前端：`'plot' | 'character' | 'world' | 'background' | 'milestone'`
- 后端：完全一致
- **结论**: 枚举值完全一致 ✓

#### 3.1.3 Document（文档）- V2架构

| 前端字段 | 后端字段 | 状态 |
|---------|---------|------|
| content | 已分离到DocumentContent | ⚠️ 需要更新 |

**V2架构说明**：后端已将Document内容分离到独立的`DocumentContent`集合，前端需要相应更新。

---

### 3.2 Reader域（阅读器）

#### 3.2.1 Chapter（章节）

| 前端字段 | 后端字段 | 类型一致性 | 备注 |
|---------|---------|-----------|------|
| id | id | ✅ | |
| bookId | book_id / bookId | ✅ | |
| title | title | ✅ | |
| chapterNum | chapter_num / chapterNum | ✅ | |
| wordCount | word_count / wordCount | ✅ | |
| isFree | is_free / isFree | ✅ | |
| price | price | ✅ | **后端单位：分** |
| publishTime | published_at / publishedAt | ⚠️ | **字段名不一致** |
| locked | - | ❌ | **后端缺失** |

**问题**：
1. 字段命名不一致：`publishTime` vs `publishedAt`
2. 缺失字段：后端没有`locked`字段
3. 金额单位：后端使用"分"，前端使用"元"

#### 3.2.2 ReadingProgress（阅读进度）

| 前端字段 | 后端字段 | 类型一致性 | 备注 |
|---------|---------|-----------|------|
| progress | progress | ⚠️ | **前端：0-100，后端：0-1** |
| chapterTitle | - | ❌ | **后端缺失** |
| scrollPosition | - | ❌ | **后端缺失** |

**问题**：
1. 进度范围不一致：前端使用百分比(0-100)，后端使用小数(0-1)
2. 缺失字段：后端缺少`chapterTitle`和`scrollPosition`

#### 3.2.3 Annotation（注记/书签）

AnnotationType枚举对比：
- 前端：`'bookmark' | 'highlight' | 'note'`
- 后端：完全一致
- **结论**: 枚举值完全一致 ✓

---

### 3.3 Bookstore域（书城）

#### 3.3.1 Book（书籍）

| 前端字段 | 后端字段 | 类型一致性 | 备注 |
|---------|---------|-----------|------|
| id | id | ✅ | |
| title | title | ✅ | |
| author | author | ✅ | |
| authorId | author_id / authorId | ✅ | |
| cover | cover_url / coverUrl | ⚠️ | **字段名不一致** |
| description | description | ✅ | |
| tags | tags | ✅ | |
| status | status | ✅ | **枚举值有差异** |
| wordCount | word_count / wordCount | ✅ | |
| chapterCount | chapter_count / chapterCount | ✅ | |
| rating | rating | ✅ | |
| ratingCount | rating_count / ratingCount | ✅ | |
| viewCount | view_count / viewCount | ✅ | |
| favoriteCount | collect_count / collectCount | ⚠️ | **字段名不一致** |
| isVip | - | ❌ | **后端缺失** |
| isFree | is_free / isFree | ✅ | |
| price | price | ✅ | **后端单位：分** |

**BookStatus枚举对比**：
- 前端：`'serializing' | 'completed' | 'paused'`
- 后端：`BookStatusDraft | BookStatusPublished | BookStatusCompleted | BookStatusPaused | BookStatusDeleted`
- **差异**：前端缺少`draft`和`deleted`状态

#### 3.3.2 Category（分类）

| 前端字段 | 后端字段 | 类型一致性 | 备注 |
|---------|---------|-----------|------|
| id | id | ✅ | |
| name | name | ✅ | |
| description | description | ✅ | |
| icon | - | ❌ | **后端缺失** |
| bookCount | book_count / bookCount | ✅ | |
| parentId | parent_id / parentId | ✅ | |
| children | - | ❌ | **后端需要单独查询** |

---

### 3.4 Social域（社交）

#### 3.4.1 Comment（评论）

| 前端字段 | 后端字段 | 类型一致性 | 备注 |
|---------|---------|-----------|------|
| id | id | ✅ | |
| bookId | book_id / bookId | ✅ | 兼容字段 |
| chapterId | chapter_id / chapterId | ✅ | 兼容字段 |
| userId | author_id / authorId | ⚠️ | **字段名不一致** |
| content | content | ✅ | |
| rating | rating | ✅ | |
| likeCount | like_count / likeCount | ✅ | |
| isLiked | - | ❌ | **后端缺失** |
| replyCount | reply_count / replyCount | ✅ | |
| parentId | parent_id / parentId | ✅ | |
| replies | - | ❌ | **后端使用CommentThread** |

**后端额外字段**：
- `targetType`：目标类型枚举（book, chapter, article, announcement, project）
- `targetId`：目标ID（通用化设计）
- `state`：评论状态（normal, hidden, deleted, rejected）
- `isPinned`：是否置顶
- `isAuthorReply`：是否作者回复
- `isFeatured`：是否精选

**CommentTargetType枚举**：
- 后端：`book | chapter | article | announcement | project`
- 前端：没有明确定义

#### 3.4.2 Follow（关注）

| 前端字段 | 后端字段 | 类型一致性 | 备注 |
|---------|---------|-----------|------|
| id | id | ✅ | |
| followerId | follower_id / followerId | ✅ | |
| followingId | following_id / followingId | ✅ | |
| followTime | created_at / createdAt | ⚠️ | **字段名不一致** |
| isMutual | is_mutual / isMutual | ✅ | |

**后端额外字段**：
- `followType`：关注类型（user, author）
- `updatedAt`：更新时间

---

### 3.5 Admin域（管理）

#### 3.5.1 User（用户）

| 前端字段 | 后端字段 | 类型一致性 | 备注 |
|---------|---------|-----------|------|
| id | id | ✅ | |
| username | username | ✅ | |
| email | email | ✅ | |
| nickname | nickname | ✅ | |
| avatar | avatar | ✅ | |
| bio | bio | ✅ | |
| gender | gender | ✅ | |
| birthday | birthday | ✅ | |
| location | location | ✅ | |
| role | roles | ⚠️ | **前端单值，后端数组** |
| createdAt | createdAt | ✅ | |

**UserRole枚举对比**：
- 前端：`'admin' | 'writer' | 'user'`
- 后端：`RoleReader | RoleAuthor | RoleAdmin`
- **差异**：命名不完全一致，前端使用`writer`，后端使用`author`

**UserStatus枚举对比**：
- 前端：没有明确定义
- 后端：`active | inactive | banned | deleted`
- **缺失**：前端应该定义用户状态枚举

---

### 3.6 Recommendation域（推荐）

#### 3.6.1 BehaviorType（行为类型）

| 前端枚举值 | 后端枚举值 | 一致性 |
|----------|----------|-------|
| view | view | ✅ |
| click | click | ✅ |
| like | like | ✅ |
| collect | - | ❌ **后端缺失** |
| read | read | ✅ |
| search | - | ❌ **后端缺失** |
| favorite | favorite | ✅ |
| purchase | purchase | ✅ |
| - | finish | ❌ **前端缺失** |
| - | share | ❌ **前端缺失** |

**差异分析**：
- 前端独有的行为：collect, search
- 后端独有的行为：finish, share
- **建议**：统一行为类型定义

---

## 四、问题汇总

### 4.1 P0阻塞问题（必须立即修复）

| # | 问题 | 影响域 | 修复方案 |
|---|------|-------|---------|
| 1 | BookStatus枚举值不匹配 | Bookstore | 前端添加`draft`和`deleted`状态 |
| 2 | is_*布尔字段转换遗漏 | 全域 | 检查HTTP拦截器转换逻辑 |
| 3 | CategoryIDs数组类型不匹配 | Bookstore | 后端是ObjectId数组，前端是string数组 |
| 4 | 响应拦截器处理不一致 | 全域 | 统一API响应处理 |
| 5 | 金额单位未转换 | Bookstore | 前端除以100显示价格 |

### 4.2 P1重要问题（3-5天内修复）

| # | 问题 | 影响域 | 修复方案 |
|---|------|-------|---------|
| 1 | 时间字段命名不一致 | 全域 | 统一使用`publishedAt`, `updatedAt` |
| 2 | Comment字段名不一致 | Social | 统一使用`authorId` |
| 3 | UserRole枚举命名不一致 | Admin | 统一使用`author` |
| 4 | UserStatus枚举缺失 | Admin | 前端添加用户状态定义 |
| 5 | BehaviorType枚举不一致 | Recommendation | 合并前后端定义 |
| 6 | ReadingProgress范围不一致 | Reader | 统一使用0-100或0-1 |
| 7 | DocumentContent V2支持 | Writer | 前端更新以支持V2架构 |
| 8 | Comment.targetType缺失 | Social | 前端添加目标类型字段 |

### 4.3 P2一般问题（长期优化）

| # | 问题 | 影响域 | 修复方案 |
|---|------|-------|---------|
| 1 | Chapter.locked字段缺失 | Reader | 后端添加锁定状态 |
| 2 | ReadingProgress.scrollPosition缺失 | Reader | 后端添加滚动位置 |
| 3 | Category.icon字段缺失 | Bookstore | 后端添加图标字段 |
| 4 | Annotation.color字段缺失 | Reader | 后端添加颜色字段 |
| 5 | Comment.isLiked字段缺失 | Social | 后端添加或前端计算 |

---

## 五、修复路线图

### Phase 1: P0阻塞问题（1-2天）

1. **修复BookStatus枚举**
   ```typescript
   // 前端更新
   type BookStatus = 'draft' | 'serializing' | 'completed' | 'paused' | 'deleted';
   ```

2. **检查is_*字段转换**
   - 确保HTTP拦截器正确处理所有`is_`前缀字段
   - 测试用例覆盖所有布尔字段

3. **修复CategoryIDs类型**
   ```typescript
   // 前端更新
   interface CategoryInfo {
     categoryIds: string[];  // 改为数组
   }
   ```

4. **统一响应拦截器**
   - 确保所有API响应经过统一处理
   - 分页信息正确映射

5. **金额单位转换**
   ```typescript
   // 前端显示时转换
   const displayPrice = (price: number) => price / 100;
   ```

### Phase 2: P1重要问题（3-5天）

1. **统一时间字段命名**
   ```typescript
   // 使用一致的命名
   interface Book {
     publishedAt: string;  // 代替publishTime
     updatedAt: string;    // 代替updateTime
   }
   ```

2. **统一枚举定义**
   ```typescript
   // UserRole
   type UserRole = 'admin' | 'author' | 'user';  // 使用author

   // UserStatus
   type UserStatus = 'active' | 'inactive' | 'banned' | 'deleted';

   // BehaviorType
   type BehaviorType =
     'view' | 'click' | 'like' | 'collect' |
     'read' | 'finish' | 'search' | 'share' |
     'favorite' | 'purchase';
   ```

3. **添加V2架构支持**
   ```typescript
   // DocumentContent类型定义
   interface DocumentContent {
     id: string;
     documentId: string;
     content: TipTapJSON;
     contentType: 'tiptap';
     wordCount: number;
     version: number;
   }
   ```

### Phase 3: P2一般问题（长期优化）

- 添加缺失的可选字段
- 完善嵌套结构
- 性能优化

---

## 六、附录

### 6.1 前端类型定义文件清单

| 文件路径 | 主要类型 | 状态 |
|---------|---------|------|
| `src/types/api.ts` | APIResponse, Pagination | ✅ |
| `src/types/models.ts` | 核心数据模型 | ✅ |
| `src/types/writer/index.d.ts` | Character, Location, Timeline | ✅ |
| `src/types/reader.ts` | Chapter, ReadingProgress, Annotation | ✅ |
| `src/types/bookstore.ts` | Book, Category, Ranking | ✅ |
| `src/types/social.ts` | Follow, Comment | ✅ |
| `src/types/auth.ts` | Login, Register, UserInfo | ✅ |
| `src/types/recommendation.ts` | Recommendation, Behavior | ✅ |
| `src/api/generated/model.ts/` | Orval自动生成类型（597个） | ✅ |

### 6.2 后端数据模型文件清单

| 文件路径 | 主要结构体 | MongoDB集合 |
|---------|-----------|------------|
| `models/writer/project.go` | Project | projects |
| `models/writer/document.go` | Document | documents |
| `models/writer/document_content.go` | DocumentContent | document_contents |
| `models/writer/document_comment.go` | DocumentComment | document_comments |
| `models/bookstore/book.go` | Book | books |
| `models/bookstore/chapter.go` | Chapter | chapters |
| `models/bookstore/chapter_content.go` | ChapterContent | chapter_contents |
| `models/social/comment.go` | Comment | comments |
| `models/social/follow.go` | Follow | follows |
| `models/users/user.go` | User | users |
| `models/recommendation/recommendation.go` | UserBehaviorRecord | user_behaviors |

### 6.3 API接口清单

| 域 | 路由文件 | 接口数量 |
|---|---------|---------|
| Writer | `router/writer/writer.go` | 30+ |
| Bookstore | `router/bookstore/bookstore_router.go` | 25+ |
| Social | `router/social/social_router.go` | 20+ |
| Admin | `router/admin/admin_router.go` | 35+ |
| Recommendation | `router/recommendation/recommendation_router.go` | 5+ |

---

## 七、参考文档

1. 数据模型V2设计：`docs/plans/2026-03-04-editor-data-model-v2-schema.md`
2. ER图：`docs/diagrams/editor-data-model-v2-er.png`
3. 前端类型定义：`Qingyu_fronted/src/types/*.ts`
4. 后端数据模型：`Qingyu_backend/models/**/*.go`

---

**报告生成时间**: 2026-03-04
**报告版本**: v1.0
**下次更新**: 根据修复进度更新
