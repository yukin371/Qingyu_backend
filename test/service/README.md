# Service层单元测试

本目录包含项目所有Service层的单元测试，按功能模块组织成子包。

## 目录结构

```
test/service/
├── comment/           # 评论系统Service测试
│   ├── comment_service_test.go
│   └── README.md
├── like/              # 点赞系统Service测试
│   ├── like_service_test.go
│   └── README.md
├── collection/        # 收藏系统Service测试
│   ├── collection_service_test.go
│   └── README.md
└── README.md         # 本文件
```

## 测试组织原则

### 子包组织
- 每个业务系统对应一个独立的子包（comment, like, collection等）
- 每个子包包含该系统Service的单元测试和Mock实现
- 独立的包命名空间避免Mock类型冲突

### Mock接口管理
- 每个子包内包含所需的Mock Repository和Mock EventBus实现
- Mock实现遵循testify/mock的标准模式
- Mock方法与实际接口保持完全一致

### 测试规范
- 测试函数命名：`TestServiceName_MethodName`
- 使用子测试组织：`t.Run("Scenario_Name", func(t *testing.T) {...})`
- 每个测试后验证Mock期望：`mockRepo.AssertExpectations(t)`
- 使用日志输出测试通过信息：`t.Logf("✓ 测试通过")`

## 运行测试

### 运行所有Service测试
```bash
go test ./test/service/... -v
```

### 运行特定系统测试
```bash
# 评论系统
go test ./test/service/comment/... -v

# 点赞系统
go test ./test/service/like/... -v

# 收藏系统
go test ./test/service/collection/... -v
```

### 查看测试覆盖率
```bash
# 所有Service测试覆盖率
go test ./test/service/... -cover

# 生成覆盖率报告
go test ./test/service/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 运行特定测试
```bash
# 运行特定测试函数
go test ./test/service/comment/... -run TestCommentService_PublishComment -v

# 运行特定子测试
go test ./test/service/like/... -run TestLikeService_LikeBook/LikeBook_Success -v
```

## 测试覆盖范围

### 评论系统（comment）
- ✅ 发布评论（参数验证、敏感词过滤、自动审核）
- ✅ 发布回复（嵌套回复、回复计数）
- ✅ 审核评论（通过、拒绝）
- ✅ 获取评论列表（分页、排序）
- ✅ 删除评论（权限检查）
- ✅ 统计功能（评分统计）

### 点赞系统（like）
- ✅ 点赞书籍/评论（幂等性、防刷检测）
- ✅ 取消点赞（重复操作检测）
- ✅ 查询点赞状态和数量
- ✅ 获取用户点赞列表
- ✅ 点赞统计（用户统计）
- ✅ 并发点赞处理

### 收藏系统（collection）
- ✅ 添加收藏（重复检测、收藏夹关联）
- ✅ 更新/删除收藏（权限检查）
- ✅ 收藏夹管理（创建、更新、删除）
- ✅ 分享收藏（公开/私有切换）
- ✅ 收藏统计（用户统计）

## 编写新测试

### 1. 创建新的子包
```bash
mkdir test/service/yourmodule
```

### 2. 创建测试文件
```go
package yourmodule

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	
	"Qingyu_backend/service/yourmodule"
	"Qingyu_backend/service/base"
)

// MockYourRepository Mock Repository
type MockYourRepository struct {
	mock.Mock
}

// 实现Repository接口方法...

// MockEventBus Mock事件总线
type MockEventBus struct {
	events []base.Event
}

// 实现EventBus接口方法...

// TestYourService_YourMethod 测试用例
func TestYourService_YourMethod(t *testing.T) {
	mockRepo := new(MockYourRepository)
	mockEventBus := NewMockEventBus()
	
	service := yourmodule.NewYourService(mockRepo, mockEventBus)
	ctx := context.Background()
	
	t.Run("Success", func(t *testing.T) {
		// Setup mocks
		mockRepo.On("YourMethod", ctx, mock.Anything).Return(nil).Once()
		
		// Execute
		err := service.YourMethod(ctx, "param")
		
		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		
		t.Logf("✓ 测试通过")
	})
}
```

### 3. 创建README.md
每个子包都应该有自己的README.md，说明：
- 测试文件列表
- 测试覆盖范围
- Mock实现说明
- 运行方法
- 维护说明

## 最佳实践

### Mock设计
1. **完整性**：Mock实现必须包含Repository/EventBus的所有方法
2. **一致性**：Mock方法签名与实际接口完全一致
3. **独立性**：每个子包维护自己的Mock实现

### 测试设计
1. **单一职责**：每个测试只验证一个功能点
2. **数据隔离**：测试间不共享状态
3. **明确断言**：使用清晰的assert语句
4. **错误路径**：同时测试成功和失败场景

### 测试覆盖
1. **基础功能**：CRUD操作的正常流程
2. **边界条件**：参数验证、空值处理
3. **权限控制**：用户权限和所有权检查
4. **并发场景**：幂等性和竞态条件
5. **错误处理**：各种错误情况的正确处理

## 维护指南

### 添加新测试
1. 在对应子包中添加新的测试函数
2. 确保Mock覆盖所有依赖
3. 运行测试验证通过
4. 更新README文档

### 更新接口
1. 同步更新Mock实现
2. 修改相关测试用例
3. 验证所有测试通过
4. 检查测试覆盖率

### 重构测试
1. 保持测试独立性
2. 避免过度Mock
3. 提取公共测试辅助函数
4. 保持测试代码简洁

## 相关文档

- [项目开发规则](../../doc/architecture/项目开发规则.md)
- [软件工程规范](../../doc/engineering/软件工程规范_v2.0.md)
- [Repository层测试](../repository/README.md)
- [集成测试](../integration/README.md)
