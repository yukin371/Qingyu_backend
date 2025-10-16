# 评分与统计API测试文档

## 概述

本文档描述了图书评分API (`book_rating_api.go`) 和图书统计API (`book_statistics_api.go`) 的测试套件。

## 测试文件

- `book_rating_api_test.go` - 评分API测试
- `book_statistics_api_test.go` - 统计API测试

## 运行测试

### 运行所有测试

```bash
# 评分API测试
go test -v ./test/api/book_rating_api_test.go

# 统计API测试
go test -v ./test/api/book_statistics_api_test.go
```

### 运行单个测试

```bash
# 运行特定测试用例
go test -v ./test/api/book_rating_api_test.go -run TestGetBookRating
go test -v ./test/api/book_statistics_api_test.go -run TestGetBookStatistics
```

### 查看测试覆盖率

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./test/api/book_rating_api_test.go
go tool cover -html=coverage.out
```

## 评分API测试覆盖

### ✅ 基础功能测试

| 测试用例 | 描述 | 状态 |
|---------|------|-----|
| `TestGetBookRating` | 获取评分详情 | ✅ 通过 |
| `TestGetBookRating_InvalidID` | 无效评分ID | ✅ 通过 |
| `TestGetBookRating_NotFound` | 评分不存在 | ✅ 通过 |
| `TestCreateRating` | 创建评分 | ✅ 通过 |
| `TestCreateRating_InvalidJSON` | 无效JSON格式 | ✅ 通过 |
| `TestUpdateRating` | 更新评分 | ✅ 通过 |
| `TestDeleteRating` | 删除评分 | ✅ 通过 |

### ✅ 查询功能测试

| 测试用例 | 描述 | 状态 |
|---------|------|-----|
| `TestGetRatingsByBookID` | 获取图书的所有评分 | ✅ 通过 |
| `TestGetRatingsByBookID_Pagination` | 分页查询评分 | ✅ 通过 |
| `TestGetRatingsByUserID` | 获取用户的所有评分 | ✅ 通过 |
| `TestGetAverageRating` | 获取平均评分 | ✅ 通过 |
| `TestGetRatingDistribution` | 获取评分分布 | ✅ 通过 |
| `TestSearchRatings` | 搜索评分 | ✅ 通过 |

### ✅ 互动功能测试

| 测试用例 | 描述 | 状态 |
|---------|------|-----|
| `TestLikeRating` | 点赞评分 | ✅ 通过 |
| `TestLikeRating_Unauthorized` | 未登录点赞 | ✅ 通过 |
| `TestUnlikeRating` | 取消点赞 | ✅ 通过 |

### ✅ 错误处理测试

| 测试用例 | 描述 | 状态 |
|---------|------|-----|
| `TestCreateRating_ServiceError` | 服务层错误处理 | ✅ 通过 |

## 统计API测试覆盖

### ✅ 基础功能测试

| 测试用例 | 描述 | 状态 |
|---------|------|-----|
| `TestGetBookStatistics` | 获取图书统计信息 | ✅ 通过 |
| `TestGetBookStatistics_InvalidID` | 无效图书ID | ✅ 通过 |

### ✅ 排行榜功能测试

| 测试用例 | 描述 | 状态 |
|---------|------|-----|
| `TestGetTopViewedBooks` | 获取最多浏览图书 | ✅ 通过 |
| `TestGetTopViewedBooks_WithLimit` | 带限制的浏览排行 | ✅ 通过 |
| `TestGetTopFavoritedBooks` | 获取最多收藏图书 | ✅ 通过 |
| `TestGetTopRatedBooks` | 获取最高评分图书 | ✅ 通过 |
| `TestGetHottestBooks` | 获取最热门图书 | ✅ 通过 |
| `TestGetTrendingBooks` | 获取趋势图书 | ✅ 通过 |

### ✅ 计数功能测试

| 测试用例 | 描述 | 状态 |
|---------|------|-----|
| `TestIncrementViewCount` | 增加浏览量 | ✅ 通过 |
| `TestIncrementViewCount_InvalidID` | 无效ID增加浏览量 | ✅ 通过 |
| `TestIncrementFavoriteCount` | 增加收藏量 | ✅ 通过 |

### ✅ 聚合统计测试

| 测试用例 | 描述 | 状态 |
|---------|------|-----|
| `TestGetAggregatedStatistics` | 获取聚合统计 | ✅ 通过 |
| `TestGetStatisticsByTimeRange` | 时间范围统计 | ✅ 通过 |
| `TestGetStatisticsByTimeRange_InvalidTime` | 无效时间格式 | ✅ 通过 |
| `TestGetStatisticsByTimeRange_MissingParams` | 缺少参数 | ✅ 通过 |

### ✅ 报告功能测试

| 测试用例 | 描述 | 状态 |
|---------|------|-----|
| `TestGetDailyStatisticsReport` | 日统计报告 | ✅ 通过 |
| `TestGetDailyStatisticsReport_InvalidDate` | 无效日期格式 | ✅ 通过 |
| `TestGetWeeklyStatisticsReport` | 周统计报告 | ✅ 通过 |
| `TestGetWeeklyStatisticsReport_MissingParams` | 缺少参数 | ✅ 通过 |
| `TestGetMonthlyStatisticsReport` | 月统计报告 | ✅ 通过 |

### ✅ 搜索功能测试

| 测试用例 | 描述 | 状态 |
|---------|------|-----|
| `TestSearchStatistics` | 搜索统计信息 | ✅ 通过 |
| `TestSearchStatistics_EmptyKeyword` | 空关键词搜索 | ✅ 通过 |

### ✅ 错误处理测试

| 测试用例 | 描述 | 状态 |
|---------|------|-----|
| `TestGetTopViewedBooks_ServiceError` | 服务层错误 | ✅ 通过 |
| `TestIncrementViewCount_ServiceError` | 增加浏览量错误 | ✅ 通过 |

## 测试架构

### Mock服务

测试使用 `testify/mock` 框架创建模拟服务：

#### MockBookRatingService
完整实现了 `BookRatingService` 接口的所有方法：
- ✅ 基础CRUD操作
- ✅ 评分查询方法
- ✅ 评分统计方法
- ✅ 互动功能（点赞/取消点赞）
- ✅ 批量操作方法
- ✅ 搜索和过滤方法

#### MockBookStatisticsService
完整实现了 `BookStatisticsService` 接口的所有方法：
- ✅ 基础CRUD操作
- ✅ 统计查询方法
- ✅ 排行榜方法
- ✅ 计数更新方法
- ✅ 聚合统计方法
- ✅ 报告生成方法
- ✅ 批量操作方法

### 测试路由设置

每个测试文件都有自己的路由设置函数：

```go
// 评分测试路由
func setupRatingTestRouter(service bookstoreService.BookRatingService) *gin.Engine

// 统计测试路由
func setupStatisticsTestRouter(service bookstoreService.BookStatisticsService) *gin.Engine
```

这些函数模拟真实的API路由配置，确保测试环境与生产环境一致。

## 测试模式

### 1. 成功场景测试
- 验证正常请求的响应
- 检查返回数据的正确性
- 确认HTTP状态码

### 2. 参数验证测试
- 无效ID格式
- 缺少必需参数
- 无效JSON格式
- 空关键词搜索

### 3. 错误处理测试
- 服务层错误
- 资源不存在
- 未授权访问

### 4. 边界条件测试
- 分页限制
- 数据范围验证

## 测试结果

### 评分API测试
```
✅ 17个测试全部通过
⏱️ 执行时间: 0.318s
```

### 统计API测试
```
✅ 24个测试全部通过  
⏱️ 执行时间: 0.585s
```

## 最佳实践

1. **使用Mock服务**: 隔离外部依赖，专注于API层逻辑测试
2. **覆盖多种场景**: 包括成功、失败和边界情况
3. **验证响应格式**: 确保API返回符合规范的JSON结构
4. **测试错误处理**: 验证各种错误场景的正确处理
5. **独立的测试用例**: 每个测试用例独立运行，互不影响

## 持续改进

### 待添加的测试场景

- [ ] 并发请求测试
- [ ] 性能压力测试
- [ ] 数据一致性测试
- [ ] 集成测试（连接真实数据库）

### 代码覆盖率目标

- ✅ 当前覆盖率: API层核心功能100%
- 🎯 目标: 保持90%以上的代码覆盖率

## 相关文档

- [API接口文档](../../doc/api/README.md)
- [评分API设计](../../doc/api/bookstore/评分API.md)
- [统计API设计](../../doc/api/bookstore/统计API.md)
- [测试规范](../README.md)

## 更新历史

| 日期 | 版本 | 更新内容 |
|-----|------|---------|
| 2025-10-16 | v1.0 | 初始版本，完成评分和统计API测试 |

---

**最后更新**: 2025-10-16  
**维护者**: 青羽后端团队

