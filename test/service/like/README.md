# 点赞系统Service层单元测试

本目录包含点赞系统Service层的单元测试。

## 测试文件

- `like_service_test.go`: LikeService的单元测试

## 测试覆盖范围

### 基础功能测试
- ✅ 点赞书籍（成功、已点赞、参数验证）
- ✅ 取消点赞书籍（成功、未点赞）
- ✅ 点赞评论（成功、已点赞）
- ✅ 取消点赞评论（成功、未点赞）

### 查询功能测试
- ✅ 获取书籍点赞数
- ✅ 获取用户点赞书籍列表
- ✅ 获取用户点赞统计

### 高级功能测试
- ✅ 并发点赞（幂等性验证）
- ✅ 防刷机制（快速点赞检测）
- ✅ 批量操作（预留扩展）

## Mock实现

### MockLikeRepository
模拟点赞Repository，提供所有点赞数据访问方法的Mock实现。

### MockCommentRepository
模拟评论Repository，用于测试点赞评论时更新评论点赞数的逻辑。

### MockEventBus
模拟事件总线，用于验证点赞事件的发布。

## 运行测试

```bash
# 运行点赞系统测试
go test ./test/service/like/... -v

# 查看测试覆盖率
go test ./test/service/like/... -cover

# 生成覆盖率报告
go test ./test/service/like/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 测试规范

1. **命名规范**：测试函数使用 `TestServiceName_MethodName` 格式
2. **子测试**：使用 `t.Run()` 组织相关测试场景
3. **Mock验证**：每个测试后调用 `AssertExpectations()` 验证Mock调用
4. **日志输出**：使用 `t.Logf()` 输出测试通过信息
5. **错误检查**：使用 `assert` 包进行断言验证

## 维护说明

- 新增Service方法时，同步添加对应的单元测试
- Mock接口变更时，及时更新Mock实现
- 保持测试的独立性，避免测试间相互依赖

