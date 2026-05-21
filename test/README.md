# 青羽平台测试套件

## 📋 概述

本目录包含青羽平台的完整测试套件，涵盖单元测试、集成测试、性能测试等多个层面，确保系统的质量、性能和可靠性。

## 🏗️ 测试结构

```
test/
├── api/                          # API接口测试
├── integration/                  # 集成测试
│   ├── ai_integration_test.go
│   ├── bookstore_integration_test.go
│   ├── stream_benchmark_test.go
│   └── stream_test.go
├── repository/                   # 仓储层测试
├── service/                      # 服务层测试
│   └── shared/                   # 共享服务测试
│       ├── auth_service_test.go
│       ├── wallet_service_test.go
│       ├── recommendation_service_test.go
│       ├── storage_service_test.go
│       ├── messaging_service_test.go
│       └── admin_service_test.go
├── bookstore_api_test.go         # 书城API测试
├── bookstore_cache_test.go       # 书城缓存测试
├── bookstore_ranking_test.go     # 书城排行测试
├── bookstore_service_test.go     # 书城服务测试
├── compatibility_test.go         # 兼容性测试
├── new_architecture_test.go      # 新架构测试
└── README.md                     # 本文件
```

## 🧪 测试类型

### 1. 单元测试 (Unit Tests)
- **位置**: `service/`, `repository/`
- **目的**: 测试单个函数或方法的功能
- **特点**: 快速执行，隔离依赖，使用Mock对象

### 2. 集成测试 (Integration Tests)
- **位置**: `integration/`
- **目的**: 测试多个组件间的集成
- **特点**: 使用真实数据库和外部服务

### 3. API测试 (API Tests)
- **位置**: `api/`, `*_api_test.go`
- **目的**: 测试HTTP API接口
- **特点**: 端到端测试，验证完整请求响应流程

### 4. 性能测试 (Performance Tests)
- **位置**: `*_benchmark_test.go`
- **目的**: 测试系统性能和资源使用
- **特点**: 基准测试，压力测试，并发测试

## 🚀 快速开始

### 环境准备

1. **安装依赖**
```bash
go mod download
```

2. **启动测试环境**
```bash
# 使用Docker Compose启动测试依赖服务
docker-compose -f docker-compose.test.yml up -d

# 等待服务启动
sleep 30
```

3. **配置测试环境变量**
```bash
export TEST_ENV=true
export MONGODB_URI="mongodb://test:test123@localhost:27017/qingyu_test"
export REDIS_ADDR="localhost:6379"
export KAFKA_BROKERS="localhost:9092"
```

### 运行测试

#### 运行所有测试
```bash
go test ./...
```

#### 运行特定包的测试
```bash
# 运行共享服务测试
go test ./test/service/shared/

# 运行集成测试
go test ./test/integration/

# 运行API测试
go test ./test/api/
```

#### 运行认证服务测试
```bash
go test ./service/auth/...
```

#### 运行特定测试用例
```bash
go test -run TestAuthService_Register_Success ./service/auth/
```

#### 生成测试覆盖率报告
```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看覆盖率统计
go tool cover -func=coverage.out

# 生成HTML覆盖率报告
go tool cover -html=coverage.out -o coverage.html
```

#### 运行性能测试
```bash
# 运行基准测试
go test -bench=. ./...

# 运行特定基准测试
go test -bench=BenchmarkAuthService_Login ./test/integration/

# 生成性能分析文件
go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof ./...
```

## 📊 共享服务测试详解

### 1. 账号权限系统测试 (`auth_service_test.go`)

#### 测试覆盖范围
- ✅ 用户注册功能
- ✅ 用户登录验证
- ✅ 权限获取和验证
- ✅ JWT Token管理
- ✅ 异常情况处理

#### 关键测试用例
```go
// 用户注册成功
TestAuthService_Register_Success

// 用户名已存在
TestAuthService_Register_UsernameExists

// 登录成功
TestAuthService_Login_Success

// 登录失败
TestAuthService_Login_InvalidCredentials

// 获取用户权限
TestAuthService_GetUserPermissions_Success
```

### 2. 钱包系统测试 (`wallet_service_test.go`)

#### 测试覆盖范围
- ✅ 钱包创建和管理
- ✅ 充值功能
- ✅ 消费功能
- ✅ 提现功能
- ✅ 交易记录管理
- ✅ 余额一致性验证

#### 关键测试用例
```go
// 创建钱包
TestWalletService_CreateWallet_Success

// 获取余额
TestWalletService_GetBalance_Success

// 充值成功
TestWalletService_Recharge_Success

// 消费成功
TestWalletService_Consume_Success

// 余额不足
TestWalletService_Consume_InsufficientBalance

// 提现申请
TestWalletService_Withdraw_Success
```

### 3. 推荐服务测试 (`recommendation_service_test.go`)

#### 测试覆盖范围
- ✅ 个性化推荐算法
- ✅ 协同过滤推荐
- ✅ 相似内容推荐
- ✅ 用户行为记录
- ✅ 用户画像更新

#### 关键测试用例
```go
// 个性化推荐（有用户画像）
TestRecommendationService_GetPersonalizedRecommendations_WithProfile

// 个性化推荐（无用户画像）
TestRecommendationService_GetPersonalizedRecommendations_NoProfile

// 相似推荐
TestRecommendationService_GetSimilarRecommendations_Success

// 协同过滤推荐
TestRecommendationService_GetCollaborativeRecommendations_Success

// 记录用户行为
TestRecommendationService_RecordUserBehavior_Success
```

### 4. 文件存储服务测试 (`storage_service_test.go`)

#### 测试覆盖范围
- ✅ 文件上传功能
- ✅ 文件下载功能
- ✅ 文件权限管理
- ✅ 文件版本控制
- ✅ 重复文件处理

#### 关键测试用例
```go
// 文件上传成功
TestStorageService_UploadFile_Success

// 重复文件处理
TestStorageService_UploadFile_DuplicateFile

// 获取文件信息
TestStorageService_GetFile_Success

// 权限拒绝
TestStorageService_GetFile_PermissionDenied

// 删除文件
TestStorageService_DeleteFile_Success
```

### 5. 消息队列服务测试 (`messaging_service_test.go`)

#### 测试覆盖范围
- ✅ 主题创建和管理
- ✅ 消息发布功能
- ✅ 消息订阅功能
- ✅ 消费者组管理
- ✅ 消息处理机制

#### 关键测试用例
```go
// 创建主题
TestMessagingService_CreateTopic_Success

// 发布消息
TestMessagingService_PublishMessage_Success

// 订阅主题
TestMessagingService_SubscribeToTopic_Success

// 处理待处理消息
TestMessagingService_ProcessPendingMessages_Success

// 删除主题
TestMessagingService_DeleteTopic_Success
```

### 6. 管理后台服务测试 (`admin_service_test.go`)

#### 测试覆盖范围
- ✅ 管理员用户管理
- ✅ 操作日志记录
- ✅ 系统配置管理
- ✅ 数据统计功能
- ✅ 仪表盘数据

#### 关键测试用例
```go
// 创建管理员用户
TestAdminService_CreateAdminUser_Success

// 更新管理员用户
TestAdminService_UpdateAdminUser_Success

// 记录操作日志
TestAdminService_LogOperation_Success

// 更新系统配置
TestAdminService_UpdateSystemConfig_Success

// 获取仪表盘统计
TestAdminService_GetDashboardStats_Success
```

## 🛠️ 测试工具和框架

### 主要依赖
```go
import (
    "testing"                           // Go标准测试框架
    "github.com/stretchr/testify/assert" // 断言库
    "github.com/stretchr/testify/mock"   // Mock框架
    "github.com/stretchr/testify/suite"  // 测试套件
    "go.mongodb.org/mongo-driver/mongo"  // MongoDB驱动
)
```

### Mock框架使用
```go
// 创建Mock对象
type MockUserRepository struct {
    mock.Mock
}

// 实现接口方法
func (m *MockUserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*User), args.Error(1)
}

// 在测试中设置期望
userRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

// 验证Mock调用
userRepo.AssertExpectations(t)
```

## 📈 测试指标和质量门禁

### 覆盖率目标
- **单元测试覆盖率**: ≥ 80%
- **集成测试覆盖率**: ≥ 60%
- **API测试覆盖率**: = 100%

### 性能指标
- **单元测试执行时间**: < 100ms/用例
- **集成测试执行时间**: < 5s/用例
- **API测试响应时间**: < 200ms

### 质量门禁
```bash
# 检查测试覆盖率
if [ $(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//') -lt 80 ]; then
    echo "Test coverage is below 80%"
    exit 1
fi

# 检查测试通过率
if ! go test ./...; then
    echo "Tests failed"
    exit 1
fi
```

## 🔧 测试配置

### 测试环境配置文件
```yaml
# config.test.yaml
database:
  mongodb:
    uri: "mongodb://test:test123@localhost:27017/qingyu_test"
    database: "qingyu_test"
  redis:
    addr: "localhost:6379"
    password: ""
    db: 1

kafka:
  brokers: ["localhost:9092"]
  
storage:
  provider: "local"
  local:
    path: "./test_uploads"

logging:
  level: "debug"
  output: "stdout"
```

### Docker测试环境
```yaml
# docker-compose.test.yml
version: '3.8'
services:
  mongodb:
    image: mongo:5.0
    environment:
      MONGO_INITDB_ROOT_USERNAME: test
      MONGO_INITDB_ROOT_PASSWORD: test123
    ports:
      - "27017:27017"
    volumes:
      - mongodb_test_data:/data/db

  redis:
    image: redis:6.2
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis_test_data:/data

  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

  kafka:
    image: confluentinc/cp-kafka:latest
    depends_on:
      - zookeeper
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    ports:
      - "9092:9092"

volumes:
  mongodb_test_data:
  redis_test_data:
```

## 🎯 最佳实践

### 1. 测试编写规范
- **命名规范**: `Test[Function]_[Scenario]_[ExpectedResult]`
- **结构清晰**: 使用Arrange-Act-Assert模式
- **独立性**: 测试用例间相互独立
- **可读性**: 代码清晰，注释完整

### 2. Mock使用原则
- **隔离外部依赖**: 数据库、网络请求、文件系统
- **模拟异常情况**: 网络错误、数据库错误、超时等
- **验证交互**: 确保正确调用了依赖的方法

### 3. 测试数据管理
- **使用工厂模式**: 创建测试数据
- **数据隔离**: 每个测试使用独立数据
- **清理机制**: 测试完成后清理数据

### 4. 性能测试建议
- **基准测试**: 使用`go test -bench`
- **内存分析**: 使用`-memprofile`
- **CPU分析**: 使用`-cpuprofile`
- **并发测试**: 使用`-race`检测竞态条件

## 🚨 常见问题

### 1. 测试环境问题
**Q: 测试时数据库连接失败**
```bash
# 确保测试数据库服务正在运行
docker-compose -f docker-compose.test.yml ps

# 检查连接配置
export MONGODB_URI="mongodb://test:test123@localhost:27017/qingyu_test"
```

**Q: Redis连接超时**
```bash
# 检查Redis服务状态
docker-compose -f docker-compose.test.yml logs redis

# 重启Redis服务
docker-compose -f docker-compose.test.yml restart redis
```

### 2. 测试执行问题
**Q: 测试覆盖率不足**
```bash
# 查看详细覆盖率报告
go tool cover -html=coverage.out

# 识别未覆盖的代码行
go tool cover -func=coverage.out | grep -v "100.0%"
```

**Q: 测试执行缓慢**
```bash
# 并行执行测试
go test -parallel 4 ./...

# 只运行快速测试
go test -short ./...
```

### 3. Mock相关问题
**Q: Mock期望设置错误**
```go
// 确保参数匹配
userRepo.On("GetByID", mock.Anything, userID).Return(expectedUser, nil)

// 使用mock.MatchedBy进行复杂匹配
userRepo.On("Create", mock.Anything, mock.MatchedBy(func(user *User) bool {
    return user.Username == "testuser"
})).Return(nil)
```

## 📚 相关文档

- [测试组织规范](../doc/testing/测试组织规范.md) - 测试分类和组织原则 ⭐ 必读
- [测试运行指南](./README_测试运行指南.md) - 快速参考手册 ⭐ 必读
- [测试最佳实践](../doc/testing/测试最佳实践.md) - 编写高质量测试的指南
- [API测试指南](../doc/testing/API测试指南.md) - API测试完整指南
- [性能测试规范](../doc/testing/性能测试规范.md) - 基准测试和压力测试
- [共享服务测试文档](../doc/testing/共享服务测试文档.md) - 共享服务测试示例
- [Postman测试指南](../doc/testing/Postman测试指南.md) - Postman使用指南

## 🤝 贡献指南

### 添加新测试
1. 在相应目录创建测试文件
2. 遵循命名规范和代码风格
3. 确保测试覆盖率达标
4. 更新相关文档

### 测试审查清单
- [ ] 测试用例覆盖正常和异常场景
- [ ] Mock对象正确设置和验证
- [ ] 测试数据独立且可重复
- [ ] 性能测试包含基准测试
- [ ] 文档更新完整
