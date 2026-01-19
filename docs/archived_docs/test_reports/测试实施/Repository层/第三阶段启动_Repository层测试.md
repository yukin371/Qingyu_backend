# 第三阶段启动 - Repository层测试

**日期**: 2025-10-19  
**前置条件**: 第二阶段Service层测试已完成  
**状态**: 待启动

---

## 背景

第二阶段已成功完成，新增105个Service层测试用例，测试通过率100%。现在进入第三阶段：Repository层测试。

---

## 目标

### 覆盖率目标
- Bookstore Repository: 70%+
- Writing Repository: 70%+
- Shared Repository: 70%+
- 整体Repository层: 70%+

### 测试策略
- **Mock测试**: 使用Mock数据库连接测试Repository逻辑
- **集成测试**: 使用真实MongoDB测试完整数据流
- **双重覆盖**: Mock保证逻辑正确，集成保证数据库操作正确

---

## 实施计划

### 3.1 Bookstore Repository测试 (优先级P0)

**预计工作量**: 4-5小时  
**预计测试用例**: 40-50个

#### 待创建测试文件

1. **BookRepository测试**
   - 文件: `test/repository/bookstore/book_repository_test.go`
   - 覆盖: CRUD操作、状态过滤、搜索、分页
   - 预计用例: 8-10个

2. **CategoryRepository测试**
   - 文件: `test/repository/bookstore/category_repository_test.go`
   - 覆盖: 分类树操作、父子关系、类型过滤
   - 预计用例: 6-8个

3. **BannerRepository测试**
   - 文件: `test/repository/bookstore/banner_repository_test.go`
   - 覆盖: Banner CRUD、状态管理、点击统计
   - 预计用例: 5-7个

4. **BookDetailRepository测试**
   - 文件: `test/repository/bookstore/book_detail_repository_test.go`
   - 覆盖: 详情CRUD、按作者/分类/标签查询、搜索
   - 预计用例: 10-12个

5. **BookStatisticsRepository测试**
   - 文件: `test/repository/bookstore/book_statistics_repository_test.go`
   - 覆盖: 统计CRUD、Top排行、热度计算
   - 预计用例: 8-10个

6. **BookRatingRepository测试**
   - 文件: `test/repository/bookstore/book_rating_repository_test.go`
   - 覆盖: 评分CRUD、统计分布、用户评分历史
   - 预计用例: 6-8个

7. **ChapterRepository测试**
   - 文件: `test/repository/bookstore/chapter_repository_test.go`
   - 覆盖: 章节CRUD、按书籍查询、VIP章节
   - 预计用例: 8-10个

---

### 3.2 Writing Repository测试 (优先级P1)

**预计工作量**: 3-4小时  
**预计测试用例**: 25-30个

#### 待创建测试文件

1. **ProjectRepository测试**
   - 文件: `test/repository/writing/project_repository_test.go`
   - 覆盖: 项目CRUD、权限验证、软删除、统计更新
   - 预计用例: 10-12个

2. **DocumentRepository测试**
   - 文件: `test/repository/writing/document_repository_test.go`
   - 覆盖: 文档CRUD、层级关系、树形结构
   - 预计用例: 8-10个

3. **DocumentContentRepository测试**
   - 文件: `test/repository/writing/document_content_repository_test.go`
   - 覆盖: 内容存储、版本管理、Delta操作
   - 预计用例: 8-10个

---

### 3.3 Shared Repository测试 (优先级P2)

**预计工作量**: 2-3小时  
**预计测试用例**: 15-20个

#### 待创建测试文件

1. **AuthRepository测试**
   - 文件: `test/repository/shared/auth_repository_test.go`
   - 覆盖: Token管理、Session存储
   - 预计用例: 5-7个

2. **WalletRepository测试**
   - 文件: `test/repository/shared/wallet_repository_test.go`
   - 覆盖: 钱包CRUD、交易记录、余额更新
   - 预计用例: 6-8个

3. **RecommendationRepository测试**
   - 文件: `test/repository/shared/recommendation_repository_test.go`
   - 覆盖: 推荐算法、用户偏好
   - 预计用例: 5-7个

---

## 技术准备

### Mock策略

#### MongoDB Mock方案
由于Repository层已经充分测试过（部分已有测试），可以采用以下策略：

1. **优先级1**: 为新功能和复杂查询编写测试
2. **优先级2**: 为核心CRUD操作补充测试
3. **优先级3**: 集成测试验证完整流程

#### 测试工具
- 使用 `testify/mock` 进行Mock
- 使用真实MongoDB进行集成测试
- 使用testcontainers创建测试数据库（可选）

---

## 预期成果

### 数量指标
- **新增测试文件**: 13个
- **新增测试用例**: 80-100个
- **测试通过率**: 100%
- **测试覆盖率**: Repository层70%+

### 质量指标
- ✅ 覆盖所有核心Repository操作
- ✅ 验证数据库查询正确性
- ✅ 测试边界条件和错误场景
- ✅ 符合Go测试最佳实践

---

## 里程碑

- [ ] M1: Bookstore Repository测试完成 (预计+3天)
- [ ] M2: Writing Repository测试完成 (预计+5天)
- [ ] M3: Shared Repository测试完成 (预计+7天)
- [ ] M4: 第三阶段验收完成 (预计+8天)

---

## 风险与挑战

### 潜在挑战
1. MongoDB连接管理（需要测试环境配置）
2. 复杂查询的Mock实现
3. 事务测试的编写
4. 大数据量场景测试

### 应对策略
1. 使用Docker提供稳定的MongoDB测试环境
2. 简化Mock，专注核心逻辑
3. 使用真实数据库进行集成测试
4. 使用表驱动测试处理多场景

---

**文档创建时间**: 2025-10-19  
**待启动**: 等待用户确认后开始第三阶段
