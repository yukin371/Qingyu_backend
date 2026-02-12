# API 层服务初始化示例

本目录包含如何使用新的工厂类和 Port 接口来初始化服务的示例代码。

## 工厂模式概述

新架构使用工厂模式来组装服务，主要优势：
1. **松耦合**：Handler 只依赖接口，不依赖具体实现
2. **易于测试**：可以注入 Mock 实现
3. **清晰的依赖关系**：通过工厂明确组装流程
4. **可维护性**：Port 接口定义清晰，每个 Port 职责单一

## 文件说明

### service_initialization.go
用户服务（UserService）的初始化示例。

### writer_service_initialization.go
写作服务（WriterService）的初始化示例。

### reader_service_initialization.go
阅读服务（ReaderService）的初始化示例。

### api_handler_usage.go
展示如何在 API Handler 中使用工厂创建的服务，包括依赖注入模式的示例。

### migration_guide.go
从旧架构迁移到新工厂模式的指南，包含逐步迁移策略。

## 使用方式

### 1. 用户服务初始化

```go
import (
    useriface "Qingyu_backend/service/interfaces/user"
    "Qingyu_backend/service/user"
    "Qingyu_backend/service/user/impl"
    userRepo "Qingyu_backend/repository/interfaces/user"
)

// 使用工厂创建用户服务
func CreateUserService(userRepo userRepo.UserRepository, roleRepo userRepo.RoleRepository) useriface.UserService {
    factory := user.NewUserServiceFactory()

    // 创建各个 Port 实现（注意：某些 Port 有依赖关系）
    passwordPort := impl.NewPasswordManagementImpl(userRepo)
    managementPort := impl.NewUserManagementImpl(userRepo)
    authPort := impl.NewUserAuthImpl(userRepo, roleRepo, passwordPort) // Auth 依赖 PasswordPort
    emailPort := impl.NewEmailManagementImpl(userRepo)
    permissionPort := impl.NewUserPermissionImpl(userRepo, roleRepo)
    statusPort := impl.NewUserStatusImpl(userRepo)

    // 使用工厂组装
    return factory.CreateWithPorts(
        managementPort,
        authPort,
        passwordPort,
        emailPort,
        permissionPort,
        statusPort,
    )
}
```

### 2. Handler 中使用服务

```go
import (
    "Qingyu_backend/api/v1/user/handler"
    useriface "Qingyu_backend/service/interfaces/user"
)

// 在应用初始化时创建服务
userService := CreateUserService(userRepo, roleRepo)

// 注入到 Handler
authHandler := handler.NewAuthHandler(userService)
profileHandler := handler.NewProfileHandler(userService)
```

## 注意事项

1. **Port 之间的依赖关系**：
   - `UserAuthPort` 依赖 `PasswordManagementPort`
   - 创建时需要先创建被依赖的 Port

2. **参数类型**：
   - 某些 impl 需要具体的服务实例（如 Reader 需要具体的 ChapterService）
   - 某些需要接口类型（如 Writer 的 CommentService）

3. **现有服务适配**：
   - 当前的 impl 实际上是包装现有服务的适配器
   - 迁移时需要确保现有服务已正确初始化

## 下一步

1. 参考 `api_handler_usage.go` 了解完整的依赖注入模式
2. 参考 `migration_guide.go` 了解如何逐步迁移现有代码
3. 运行测试验证：`go test ./api/v1/examples/...`
