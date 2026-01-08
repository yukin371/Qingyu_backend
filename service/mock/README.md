# Mock接口使用说明

本目录包含使用testify/mock生成的Mock接口实现，用于测试。

## 📁 目录结构

```
service/mock/
├── ai/                           # AI服务相关Mock
│   └── quota_repository_mock.go  # 配额仓储Mock
├── bookstore/                    # 书店服务相关Mock
│   └── book_detail_repository_mock.go  # 书籍详情仓储Mock
├── writer/                       # 写作服务相关Mock
│   └── document_content_repository_mock.go  # 文档内容仓储Mock
└── shared/                       # 共享服务相关Mock
    └── (待添加)
```

## 🚀 使用方法

### 1. 在测试中导入Mock包

```go
import (
    "Qingyu_backend/service/mock"
    "Qingyu_backend/service/ai"

    "github.com/stretchr/testify/mock"
)
```

### 2. 创建Mock实例

```go
func TestQuotaService(t *testing.T) {
    mockRepo := new(mock.MockQuotaRepository)
    service := ai.NewQuotaService(mockRepo)

    // 设置期望的调用
    mockRepo.On("GetQuotaByUserID", mock.Anything, "user123", mock.Anything).Return(
        &ai.UserQuota{UserID: "user123", TotalQuota: 1000},
        nil,
    )

    // 调用被测试的方法
    quota, err := service.GetQuotaInfo(context.Background(), "user123")

    // 验证结果
    assert.NoError(t, err)
    assert.Equal(t, "user123", quota.UserID)

    // 验证Mock调用
    mockRepo.AssertExpectations(t)
}
```

### 3. 设置返回值

```go
// 返回单个值
mockRepo.On("GetQuotaByUserID", ctx, userID, quotaType).Return(quota, nil)

// 返回错误
mockRepo.On("CreateQuota", ctx, mock.Anything).Return(errors.New("创建失败"))

// 返回nil
mockRepo.On("GetByID", ctx, "nonexistent").Return(nil, nil)
```

### 4. 验证调用

```go
// 验证方法被调用
mockRepo.AssertCalled(t, "UpdateQuota", ctx, mock.Anything)

// 验证方法调用次数
mockRepo.AssertNumberOfCalls(t, "UpdateQuota", 1)

// 验证所有期望
mockRepo.AssertExpectations(t)
```

## 📝 已实现的Mock接口

| Mock文件 | 接口 | 用途 |
|---------|------|------|
| `quota_repository_mock.go` | `ai.QuotaRepository` | AI配额管理 |
| `book_detail_repository_mock.go` | `bookstore.BookDetailRepository` | 书籍详情管理 |
| `document_content_repository_mock.go` | `writer.DocumentContentRepository` | 文档内容管理 |

## 🔧 添加新的Mock接口

### 方法1: 手动创建（适用于简单接口）

```go
// service/mock/your_module/your_repository_mock.go
package mock

import (
    "context"
    "github.com/stretchr/testify/mock"
)

type MockYourRepository struct {
    mock.Mock
}

func (m *MockYourRepository) YourMethod(ctx context.Context, param string) error {
    args := m.Called(ctx, param)
    return args.Error(0)
}
```

### 方法2: 使用mockgen工具（推荐）

```bash
# 1. 安装mockgen
go install github.com/golang/mock/mockgen@latest

# 2. 生成Mock文件
mockgen -source=repository/interfaces/your_module/your_repository_interface.go \
        -destination=service/mock/your_module/your_repository_mock.go \
        -package=mock

# 3. 转换为testify/mock格式（如果需要）
# 手动替换gomock为testify/mock
```

## ⚠️ 注意事项

1. **接口匹配**: Mock方法签名必须与实际接口完全匹配
2. **参数匹配**: 使用`mock.Anything`匹配任意参数，或使用具体值匹配
3. **返回值**: 返回值类型和数量必须与接口定义一致
4. **nil检查**: 对于返回指针的方法，需要检查`args.Get(0)`是否为nil

## 🐛 常见问题

### 问题1: 接口不匹配
```
error: MockQuotaRepository does not implement QuotaRepository (missing method GetTotalConsumption)
```
**解决**: 确保Mock实现了所有接口方法

### 问题2: 参数类型错误
```
cannot use mockRepo (type *MockQuotaRepository) as type QuotaRepository
```
**解决**: 检查方法签名是否完全匹配，包括参数类型和返回值类型

### 问题3: Mock返回值问题
```
panic: runtime error: invalid memory address or nil pointer dereference
```
**解决**: 在Return中提供非nil的返回值，或在测试中检查nil

## 📚 参考资料

- [Testify Mock文档](https://github.com/stretchr/testify#mock-package)
- [Go Mock最佳实践](https://github.com/golang/mock)
- [项目测试文档](../../doc/测试/)
