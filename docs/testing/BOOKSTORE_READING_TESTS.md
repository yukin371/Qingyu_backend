# 青羽写作平台 - 书城系统与阅读功能测试文档

## 概述

本目录包含青羽写作平台书城系统和阅读功能的完整单元测试套件。测试覆盖了章节购买、目录管理、主题系统和评论功能。

## 测试文件结构

```
Qingyu_backend/
├── service/bookstore/
│   └── chapter_purchase_service_test.go          # 章节购买服务单元测试
├── api/v1/bookstore/
│   └── chapter_catalog_api_test.go               # 章节目录API测试
├── api/v1/reader/
│   ├── theme_api_test.go                         # 主题管理API测试
│   └── chapter_comment_api_test.go               # 章节评论API测试
└── service/reader/mocks/
    └── reader_mocks.go                           # Mock接口实现
```

## 测试覆盖范围

### 1. 章节购买服务测试 (`chapter_purchase_service_test.go`)

#### 测试的功能：
- **章节目录管理**
  - ✅ 获取章节目录（成功、书籍不存在、空书籍ID）
  - ✅ 获取试读章节（成功、默认试读数量）
  - ✅ 获取VIP章节

- **章节购买**
  - ✅ 购买单个章节（成功、已购买、余额不足、免费章节）
  - ✅ 批量购买章节（成功）
  - ✅ 购买全书（成功、已购买）
  - ✅ 权限检查（免费章节、已购买、无权限）

- **购买记录查询**
  - ✅ 获取章节购买记录（分页）
  - ✅ 获取书籍购买记录
  - ✅ 获取所有购买记录

- **价格计算**
  - ✅ 获取章节价格
  - ✅ 计算全书价格（原价、折扣价）

#### Mock的服务：
- `ChapterRepository` - 章节数据访问
- `ChapterPurchaseRepository` - 购买记录数据访问
- `BookStoreRepository` - 书籍数据访问
- `WalletService` - 钱包服务（扣费、退款）
- `CacheService` - 缓存服务

### 2. 章节目录API测试 (`chapter_catalog_api_test.go`)

#### 测试的API端点：

**GET** `/api/v1/bookstore/books/:id/chapters`
- ✅ 获取章节目录成功
- ✅ 无效的书籍ID格式
- ✅ 空书籍ID

**GET** `/api/v1/bookstore/books/:id/chapters/:chapterId`
- ✅ 获取章节信息成功
- ✅ 章节不存在

**GET** `/api/v1/bookstore/books/:id/trial-chapters`
- ✅ 获取试读章节成功
- ✅ 自定义试读数量

**GET** `/api/v1/bookstore/books/:id/vip-chapters`
- ✅ 获取VIP章节成功

**GET** `/api/v1/bookstore/chapters/:chapterId/price`
- ✅ 获取章节价格成功

**POST** `/api/v1/reader/chapters/:chapterId/purchase`
- ✅ 购买章节成功
- ✅ 余额不足处理

**POST** `/api/v1/reader/books/:id/buy-all`
- ✅ 购买全书成功

**GET** `/api/v1/reader/purchases`
- ✅ 获取购买记录（分页）

**GET** `/api/v1/reader/purchases/:bookId`
- ✅ 获取特定书籍的购买记录

**GET** `/api/v1/bookstore/chapters/:chapterId/access`
- ✅ 检查章节访问权限

### 3. 主题管理API测试 (`theme_api_test.go`)

#### 测试的API端点：

**GET** `/api/v1/reader/themes`
- ✅ 获取所有主题
- ✅ 仅获取内置主题
- ✅ 仅获取公开主题

**GET** `/api/v1/reader/themes/:name`
- ✅ 根据名称获取主题（light、dark、sepia、eye-care）
- ✅ 主题不存在
- ✅ 空主题名称

**POST** `/api/v1/reader/themes`
- ✅ 创建自定义主题成功
- ✅ 缺少必填字段
- ✅ 未授权访问

**PUT** `/api/v1/reader/themes/:id`
- ✅ 更新主题成功
- ✅ 空主题ID
- ✅ 未授权访问

**DELETE** `/api/v1/reader/themes/:id`
- ✅ 删除主题成功
- ✅ 空主题ID
- ✅ 未授权访问

**POST** `/api/v1/reader/themes/:name/activate`
- ✅ 激活主题成功
- ✅ 无效主题
- ✅ 空主题名称
- ✅ 未授权访问

#### 主题颜色验证：
- ✅ 主题颜色配置完整性
- ✅ 内置主题必需颜色字段检查

### 4. 章节评论API测试 (`chapter_comment_api_test.go`)

#### 测试的API端点：

**GET** `/api/v1/reader/chapters/:chapterId/comments`
- ✅ 获取章节评论列表（成功、分页、排序）
- ✅ 无效章节ID
- ✅ 空章节ID
- ✅ 无效排序字段（应默认）
- ✅ 过滤顶级评论
- ✅ 页码大小限制

**POST** `/api/v1/reader/chapters/:chapterId/comments`
- ✅ 发表评论成功
- ✅ 回复评论
- ✅ 无效章节ID
- ✅ 评分范围验证（0-5）
- ✅ 空评论内容
- ✅ 未授权访问
- ✅ 无效父评论ID

**GET** `/api/v1/reader/comments/:commentId`
- ✅ 获取单条评论详情
- ✅ 无效评论ID
- ✅ 空评论ID
- ✅ 评论不存在

**PUT** `/api/v1/reader/comments/:commentId`
- ✅ 更新评论成功
- ✅ 更新评分
- ✅ 无效评分范围
- ✅ 未授权访问

**DELETE** `/api/v1/reader/comments/:commentId`
- ✅ 删除评论成功
- ✅ 未授权访问

**POST** `/api/v1/reader/comments/:commentId/like`
- ✅ 点赞评论成功
- ✅ 未授权访问

**DELETE** `/api/v1/reader/comments/:commentId/like`
- ✅ 取消点赞成功
- ✅ 未授权访问

**GET** `/api/v1/reader/chapters/:chapterId/paragraphs/:paragraphIndex/comments`
- ✅ 获取段落评论成功
- ✅ 无效章节ID
- ✅ 无效段落索引
- ✅ 负数段落索引

**POST** `/api/v1/reader/chapters/:chapterId/paragraph-comments`
- ✅ 发表段落评论成功
- ✅ 缺少段落索引
- ✅ 未授权访问

**GET** `/api/v1/reader/chapters/:chapterId/paragraph-comments`
- ✅ 获取章节所有段落评论概览

#### 模型方法测试：
- ✅ `IsParagraphComment()` - 判断是否为段落级评论
- ✅ `IsTopLevel()` - 判断是否为顶级评论
- ✅ `CanEdit()` - 判断是否可编辑（30分钟内）
  - ✅ 在可编辑时间内
  - ✅ 超过可编辑时间
  - ✅ 已删除评论

## 运行测试

### 前置要求

1. **安装依赖**
```bash
cd D:\Github\青羽\Qingyu_backend
go mod download
```

2. **安装测试框架**
```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
```

### 运行所有测试

```bash
# 在项目根目录
go test ./service/bookstore/... -v
go test ./api/v1/bookstore/... -v
go test ./api/v1/reader/... -v
```

### 运行特定测试文件

```bash
# 章节购买服务测试
go test ./service/bookstore/chapter_purchase_service_test.go -v

# 章节目录API测试
go test ./api/v1/bookstore/chapter_catalog_api_test.go -v

# 主题API测试
go test ./api/v1/reader/theme_api_test.go -v

# 章节评论API测试
go test ./api/v1/reader/chapter_comment_api_test.go -v
```

### 运行特定测试用例

```bash
# 运行单个测试
go test ./service/bookstore/... -v -run TestChapterPurchaseService_PurchaseChapter_Success

# 运行包含特定关键词的所有测试
go test ./service/bookstore/... -v -run Purchase

# 运行基准测试
go test ./api/v1/reader/... -bench=. -benchmem
```

### 测试覆盖率

```bash
# 生成覆盖率报告
go test ./service/bookstore/... -coverprofile=coverage.out
go test ./api/v1/bookstore/... -coverprofile=coverage.out
go test ./api/v1/reader/... -coverprofile=coverage.out

# 查看覆盖率
go tool cover -html=coverage.out

# 查看覆盖率百分比
go tool cover -func=coverage.out
```

## 测试数据准备

### Mock数据

测试使用Mock对象模拟依赖服务：

```go
// 示例：创建Mock章节
chapter := &bookstore.Chapter{
    ID:         primitive.NewObjectID(),
    BookID:     primitive.NewObjectID(),
    Title:      "Test Chapter",
    ChapterNum: 1,
    WordCount:  2000,
    IsFree:     false,
    Price:      1.99,
    PublishTime: time.Now(),
}

// 示例：创建Mock购买记录
purchase := &bookstore.ChapterPurchase{
    ID:        primitive.NewObjectID(),
    UserID:    primitive.NewObjectID(),
    ChapterID: primitive.NewObjectID(),
    Price:     1.99,
}
```

## 关键业务逻辑测试

### 1. 购买流程测试

#### 单章购买流程：
1. ✅ 检查章节是否已购买
2. ✅ 验证章节存在性
3. ✅ 检查是否为免费章节
4. ✅ 验证书籍存在性
5. ✅ 检查用户余额
6. ✅ 执行事务扣费
7. ✅ 创建购买记录
8. ✅ 清除缓存

#### 全书购买流程：
1. ✅ 检查全书是否已购买
2. ✅ 获取所有付费章节
3. ✅ 计算折扣价格
4. ✅ 检查用户余额
5. ✅ 执行事务扣费
6. ✅ 创建全书购买记录
7. ✅ 为所有章节创建购买记录
8. ✅ 清除缓存

### 2. 权限控制测试

#### 章节访问权限：
- ✅ 免费章节：直接访问
- ✅ 已购买章节：可访问
- ✅ 全书购买：所有章节可访问
- ✅ VIP章节：VIP用户可访问
- ✅ 未购买付费章节：不可访问

### 3. 分页和排序测试

#### 分页验证：
- ✅ 默认页码：1
- ✅ 默认每页数量：20
- ✅ 最大每页数量：100
- ✅ 总数统计
- ✅ 总页数计算

#### 排序验证：
- ✅ 按创建时间排序
- ✅ 按点赞数排序
- ✅ 按评分排序
- ✅ 升序/降序

### 4. 边界条件测试

#### 输入验证：
- ✅ 空ID处理
- ✅ 无效ID格式
- ✅ 负数页码
- ✅ 超大页码大小
- ✅ 评分范围（0-5）
- ✅ 空评论内容

#### 异常处理：
- ✅ 资源不存在
- ✅ 余额不足
- ✅ 重复购买
- ✅ 未授权访问
- ✅ 软删除处理

## Mock接口说明

### MockWalletService
模拟钱包服务，用于测试购买流程中的扣费和退款操作。

```go
type MockWalletService struct {
    mock.Mock
}

// 主要方法：
- GetBalance(ctx, userID) -> (float64, error)
- Consume(ctx, userID, amount, description) -> (string, error)
- Refund(ctx, userID, amount, description) -> (string, error)
```

### MockChapterRepository
模拟章节数据访问，用于测试章节查询功能。

```go
type MockChapterRepository struct {
    mock.Mock
}

// 主要方法：
- GetByID(ctx, id) -> (*Chapter, error)
- GetByBookID(ctx, bookID, limit, offset) -> ([]*Chapter, error)
- GetFreeChapters(ctx, bookID, limit, offset) -> ([]*Chapter, error)
- GetPaidChapters(ctx, bookID, limit, offset) -> ([]*Chapter, error)
```

### MockChapterPurchaseRepository
模拟购买记录数据访问，用于测试购买和查询功能。

```go
type MockChapterPurchaseRepository struct {
    mock.Mock
}

// 主要方法：
- GetByUserAndChapter(ctx, userID, chapterID) -> (*ChapterPurchase, error)
- Create(ctx, purchase) -> error
- GetByUser(ctx, userID, page, pageSize) -> ([]*ChapterPurchase, int64, error)
- CheckUserPurchasedChapter(ctx, userID, chapterID) -> (bool, error)
- Transaction(ctx, fn) -> error
```

## CI/CD集成

### GitHub Actions配置示例

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Install dependencies
        run: go mod download
      - name: Run tests
        run: |
          go test ./service/bookstore/... -v -cover
          go test ./api/v1/bookstore/... -v -cover
          go test ./api/v1/reader/... -v -cover
      - name: Upload coverage
        uses: codecov/codecov-action@v2
```

## 测试最佳实践

### 1. 测试命名

使用描述性的测试名称：
```
✅ TestChapterPurchaseService_PurchaseChapter_Success
✅ TestChapterPurchaseService_PurchaseChapter_InsufficientBalance
❌ TestPurchaseChapter1
```

### 2. 测试结构

遵循AAA模式（Arrange-Act-Assert）：

```go
func TestChapterPurchaseService_PurchaseChapter_Success(t *testing.T) {
    // Arrange - 准备测试数据和环境
    chapterRepo := new(MockChapterRepository)
    purchaseRepo := new(MockChapterPurchaseRepository)
    // ...

    // Act - 执行被测试的操作
    purchase, err := service.PurchaseChapter(ctx, userID, chapterID)

    // Assert - 验证结果
    assert.NoError(t, err)
    assert.NotNil(t, purchase)
    assert.Equal(t, userID, purchase.UserID)
}
```

### 3. Mock期望设置

确保设置所有必要的Mock期望：

```go
chapterRepo.On("GetByID", ctx, chapterID).Return(chapter, nil)
chapterRepo.On("GetByBookID", ctx, bookID, 10000, 0).Return(chapters, nil)
purchaseRepo.On("GetByUserAndChapter", ctx, userID, chapterID).Return(nil, errors.New("not found"))

// 执行测试后验证期望
chapterRepo.AssertExpectations(t)
purchaseRepo.AssertExpectations(t)
```

### 4. 测试隔离

每个测试应该独立，不依赖其他测试：

```go
func TestChapterPurchaseService_PurchaseChapter_Success(t *testing.T) {
    // 每个测试创建新的Mock实例
    chapterRepo := new(MockChapterRepository)
    purchaseRepo := new(MockChapterPurchaseRepository)
    // ...
}
```

## 性能基准测试

```bash
# 运行基准测试
go test ./api/v1/reader/... -bench=BenchmarkChapterCommentAPI_GetChapterComments -benchmem

# 输出示例：
# BenchmarkChapterCommentAPI_GetChapterComments-8   	   10000	    123456 ns/op	   12345 B/op	    123 allocs/op
```

## 故障排查

### 常见问题

1. **导入错误**
```
import cycle not allowed
```
**解决方案**：检查是否有循环依赖，确保mock文件与实现文件分离

2. **Mock不匹配**
```
 Expected: GetByID(ctx, chapterID)
 Actual: GetByID(ctx, <different>)
```
**解决方案**：确保Mock期望的参数与实际调用参数完全一致

3. **数据库连接**
```
connection refused
```
**解决方案**：单元测试不应连接真实数据库，使用Mock对象

## 贡献指南

### 添加新测试

1. 在相应的测试文件中添加测试函数
2. 遵循命名规范：`Test{Service/API}_{Method}_{Scenario}`
3. 包含成功场景和失败场景
4. 使用Mock对象隔离依赖
5. 添加必要的注释说明测试目的

### 代码审查清单

- [ ] 测试覆盖了主要功能路径
- [ ] 包含边界条件测试
- [ ] Mock对象正确设置
- [ ] 断言清晰明确
- [ ] 测试名称描述性强
- [ ] 无硬编码的测试数据
- [ ] 测试之间相互独立

## 联系方式

如有问题或建议，请联系开发团队。

---

**最后更新**: 2026-01-03
**版本**: 1.0.0
**维护者**: 青羽开发团队
