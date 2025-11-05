# 收藏系统Service层单元测试

本目录包含收藏系统Service层的单元测试。

## 测试文件

- `collection_service_test.go`: CollectionService的单元测试

## 测试覆盖范围

### 收藏功能测试
- ✅ 添加收藏（成功、已收藏、参数验证）
- ✅ 添加到收藏夹（成功、权限检查）
- ✅ 删除收藏（成功、权限检查）
- ✅ 更新收藏（成功、权限检查）
- ✅ 获取用户收藏列表

### 收藏夹功能测试
- ✅ 创建收藏夹（成功、名称验证）
- ✅ 更新收藏夹（成功、权限检查）
- ✅ 删除收藏夹（成功、非空检测）

### 高级功能测试
- ✅ 分享收藏（公开设置）
- ✅ 获取收藏统计（用户统计）

## Mock实现

### MockCollectionRepository
模拟收藏Repository，提供所有收藏和收藏夹数据访问方法的Mock实现。

### MockEventBus
模拟事件总线，用于验证收藏事件的发布。

## 运行测试

```bash
# 运行收藏系统测试
go test ./test/service/collection/... -v

# 查看测试覆盖率
go test ./test/service/collection/... -cover

# 生成覆盖率报告
go test ./test/service/collection/... -coverprofile=coverage.out
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

