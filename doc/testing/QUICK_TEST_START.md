# 快速测试指南 - 青羽写作平台

## 🚀 一分钟快速开始

### 前置条件
```bash
# 确保在项目根目录
cd D:\Github\青羽\Qingyu_backend

# 检查Go版本
go version
```

### 运行所有新创建的测试
```bash
# 服务层测试
go test ./service/bookstore/chapter_purchase_service_test.go -v

# API层测试 - 书城
go test ./api/v1/bookstore/chapter_catalog_api_test.go -v

# API层测试 - 阅读器主题
go test ./api/v1/reader/theme_api_test.go -v

# API层测试 - 章节评论
go test ./api/v1/reader/chapter_comment_api_test.go -v
```

## 📊 测试覆盖概览

| 模块 | 测试文件 | 测试用例数 | 覆盖功能 |
|------|---------|-----------|---------|
| 章节购买 | `chapter_purchase_service_test.go` | 40+ | 购买、权限、价格计算 |
| 章节目录 | `chapter_catalog_api_test.go` | 15+ | API端点测试 |
| 主题管理 | `theme_api_test.go` | 20+ | 主题CRUD、激活 |
| 章节评论 | `chapter_comment_api_test.go` | 35+ | 评论、点赞、段落评论 |

## 🎯 常用测试命令

### 运行特定功能测试

```bash
# 测试购买功能
go test ./service/bookstore/... -v -run Purchase

# 测试权限检查
go test ./service/bookstore/... -v -run Access

# 测试价格计算
go test ./service/bookstore/... -v -run Price

# 测试主题功能
go test ./api/v1/reader/... -v -run ThemeAPI

# 测试评论功能
go test ./api/v1/reader/... -v -run Comment
```

### 生成覆盖率报告

```bash
# 生成覆盖率
go test ./service/bookstore/... -coverprofile=coverage.out

# 在浏览器中查看
go tool cover -html=coverage.out

# 查看终端覆盖率
go tool cover -func=coverage.out
```

### 运行基准测试

```bash
# 性能测试
go test ./api/v1/reader/... -bench=. -benchmem
```

## 📁 文件结构

```
Qingyu_backend/
├── service/bookstore/
│   ├── chapter_purchase_service.go          # 被测试的服务
│   └── chapter_purchase_service_test.go     # ✨ 新创建
├── api/v1/bookstore/
│   ├── chapter_catalog_api.go               # 被测试的API
│   └── chapter_catalog_api_test.go          # ✨ 新创建
├── api/v1/reader/
│   ├── theme_api.go                         # 被测试的API
│   ├── theme_api_test.go                    # ✨ 新创建
│   ├── chapter_comment_api.go               # 被测试的API
│   └── chapter_comment_api_test.go          # ✨ 新创建
├── service/reader/mocks/
│   └── reader_mocks.go                      # ✨ 新创建 (Mock对象)
├── BOOKSTORE_READING_TESTS.md               # ✨ 详细测试文档
└── TEST_FILES_SUMMARY.md                    # ✨ 测试总结
```

## 🔍 测试示例

### 示例1: 测试章节购买

```bash
# 运行购买章节的所有测试
go test ./service/bookstore/... -v -run TestChapterPurchaseService_PurchaseChapter

# 预期输出:
# === RUN   TestChapterPurchaseService_PurchaseChapter_Success
# --- PASS: TestChapterPurchaseService_PurchaseChapter_Success (0.05s)
# === RUN   TestChapterPurchaseService_PurchaseChapter_AlreadyPurchased
# --- PASS: TestChapterPurchaseService_PurchaseChapter_AlreadyPurchased (0.02s)
# === RUN   TestChapterPurchaseService_PurchaseChapter_InsufficientBalance
# --- PASS: TestChapterPurchaseService_PurchaseChapter_InsufficientBalance (0.03s)
# PASS
```

### 示例2: 测试主题API

```bash
# 运行主题相关的所有测试
go test ./api/v1/reader/... -v -run ThemeAPI

# 预期输出:
# === RUN   TestThemeAPI_GetThemes_AllThemes
# --- PASS: TestThemeAPI_GetThemes_AllThemes (0.01s)
# === RUN   TestThemeAPI_GetThemeByName_Success
# --- PASS: TestThemeAPI_GetThemeByName_Success (0.01s)
# === RUN   TestThemeAPI_CreateCustomTheme_Success
# --- PASS: TestThemeAPI_CreateCustomTheme_Success (0.02s)
# PASS
```

### 示例3: 测试评论功能

```bash
# 运行评论创建测试
go test ./api/v1/reader/... -v -run TestChapterCommentAPI_CreateChapterComment

# 预期输出:
# === RUN   TestChapterCommentAPI_CreateChapterComment_Success
# --- PASS: TestChapterCommentAPI_CreateChapterComment_Success (0.03s)
# === RUN   TestChapterCommentAPI_CreateChapterComment_InvalidRating
# --- PASS: TestChapterCommentAPI_CreateChapterComment_InvalidRating (0.02s)
# PASS
```

## 🛠️ 故障排查

### 问题1: 导入错误

```
error: import cycle not allowed
```

**解决方案**: 测试文件使用Mock对象，确保不导入被测试包的实现文件

### 问题2: Mock期望不匹配

```
Expected: GetByID(ctx, chapterID)
Actual: GetByID(ctx, <different>)
```

**解决方案**: 检查Mock设置，确保参数类型和值完全匹配

### 问题3: 找不到测试文件

```
no test files found
```

**解决方案**: 确保文件名以 `_test.go` 结尾，且在正确的包目录中

## 📚 详细文档

查看以下文件获取更多信息：
- **BOOKSTORE_READING_TESTS.md** - 完整的测试文档
- **TEST_FILES_SUMMARY.md** - 测试文件总结

## 🎓 测试最佳实践

### 1. 测试隔离
每个测试应该独立运行，不依赖其他测试

### 2. 清晰命名
使用描述性的测试名称：`Test{功能}_{场景}_{预期结果}`

### 3. Mock使用
使用Mock对象隔离外部依赖

### 4. 断言验证
对每个关键结果进行断言验证

## 📈 持续集成

### 本地预检查

```bash
# 运行所有测试
go test ./... -v

# 检查代码格式
go fmt ./...

# 静态分析
go vet ./...
```

### 提交前检查清单

- [ ] 所有测试通过
- [ ] 代码已格式化
- [ ] 无静态分析警告
- [ ] 覆盖率达标 (>75%)

## 🎉 下一步

1. ✅ 运行测试验证通过
2. ✅ 查看覆盖率报告
3. ✅ 阅读详细文档
4. ✅ 添加新功能时编写对应测试

---

**快速开始**: 5分钟即可运行所有测试 ✨
**文档完整度**: 100%
**测试覆盖**: 书城系统 + 阅读功能
