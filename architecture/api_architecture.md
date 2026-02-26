# API 架构图文档

> **版本**: v1.0
> **创建日期**: 2026-02-26
> **最后更新**: 2026-02-26

---

## 目录

- [API层整体架构图](#api层整体架构图)
- [目录结构图](#目录结构图)
- [请求处理流程图](#请求处理流程图)
- [错误处理流程图](#错误处理流程图)
- [重构前后对比图](#重构前后对比图)
- [模块依赖关系图](#模块依赖关系图)
- [数据流图](#数据流图)

---

## API层整体架构图

```mermaid
graph TB
    subgraph "客户端层"
        Web[Web前端]
        Mobile[移动端]
        ThirdParty[第三方应用]
    end

    subgraph "API网关 / 路由层"
        Router[gin.Engine Router]
        CORS[CORS中间件]
        Logger[日志中间件]
        Recovery[恢复中间件]
    end

    subgraph "认证授权层"
        AuthMiddleware[JWT认证中间件]
        PermissionMiddleware[权限检查中间件]
        RateLimitMiddleware[限流中间件]
    end

    subgraph "API v1 层 (业务模块)"
        AdminAPI[admin<br/>管理员功能]
        AIAPI[ai<br/>AI服务]
        AuthAPI[auth<br/>认证授权]
        BookstoreAPI[bookstore<br/>书店]
        FinanceAPI[finance<br/>财务]
        MessagesAPI[messages<br/>消息]
        NotificationsAPI[notifications<br/>通知]
        RecommendationAPI[recommendation<br/>推荐]
        SearchAPI[search<br/>搜索]
        SharedAPI[shared<br/>共享工具]
        SocialAPI[social<br/>社交]
        StatsAPI[stats<br/>统计]
        SystemAPI[system<br/>系统]
    end

    subgraph "共享层"
        ResponseHelper[响应助手]
        Validator[请求验证器]
        CommonDTOs[公共DTO]
        Converters[数据转换器]
    end

    subgraph "服务层"
        UserService[用户服务]
        ContentService[内容服务]
        SocialService[社交服务]
        AIService[AI服务]
        FinanceService[财务服务]
    end

    subgraph "数据访问层"
        UserRepo[用户仓储]
        ContentRepo[内容仓储]
        SocialRepo[社交仓储]
    end

    subgraph "基础设施层"
        MongoDB[(MongoDB)]
        Redis[(Redis)]
        Milvus[(Milvus)]
        MinIO[(MinIO)]
    end

    Web --> Router
    Mobile --> Router
    ThirdParty --> Router

    Router --> CORS
    CORS --> Logger
    Logger --> Recovery

    Recovery --> AuthMiddleware

    AuthMiddleware --> AdminAPI
    AuthMiddleware --> AIAPI
    AuthMiddleware --> AuthAPI
    AuthMiddleware --> BookstoreAPI
    AuthMiddleware --> FinanceAPI
    AuthMiddleware --> MessagesAPI
    AuthMiddleware --> NotificationsAPI
    AuthMiddleware --> RecommendationAPI
    AuthMiddleware --> SearchAPI
    AuthMiddleware --> SocialAPI
    AuthMiddleware --> StatsAPI
    AuthMiddleware --> SystemAPI

    AuthMiddleware --> PermissionMiddleware
    PermissionMiddleware --> RateLimitMiddleware

    AdminAPI -.使用.-> SharedAPI
    AIAPI -.使用.-> SharedAPI
    AuthAPI -.使用.-> SharedAPI
    BookstoreAPI -.使用.-> SharedAPI
    FinanceAPI -.使用.-> SharedAPI
    MessagesAPI -.使用.-> SharedAPI
    NotificationsAPI -.使用.-> SharedAPI
    RecommendationAPI -.使用.-> SharedAPI
    SearchAPI -.使用.-> SharedAPI
    SocialAPI -.使用.-> SharedAPI
    StatsAPI -.使用.-> SharedAPI
    SystemAPI -.使用.-> SharedAPI

    SharedAPI --> ResponseHelper
    SharedAPI --> Validator
    SharedAPI --> CommonDTOs
    SharedAPI --> Converters

    AdminAPI --> UserService
    AIAPI --> AIService
    AuthAPI --> UserService
    BookstoreAPI --> ContentService
    FinanceAPI --> FinanceService
    MessagesAPI --> SocialService
    NotificationsAPI --> SocialService
    SocialAPI --> SocialService

    UserService --> UserRepo
    ContentService --> ContentRepo
    SocialService --> SocialRepo

    UserRepo --> MongoDB
    ContentRepo --> MongoDB
    SocialRepo --> MongoDB

    UserService --> Redis
    ContentService --> Milvus
    FinanceService --> MinIO

    classDef client fill:#e1f5fe
    classDef gateway fill:#fff3e0
    classDef auth fill:#f3e5f5
    classDef api fill:#e8f5e9
    classDef shared fill:#fce4ec
    classDef service fill:#fff9c4
    classDef repo fill:#e0f2f1
    classDef db fill:#efebe9

    class Web,Mobile,ThirdParty client
    class Router,CORS,Logger,Recovery gateway
    class AuthMiddleware,PermissionMiddleware,RateLimitMiddleware auth
    class AdminAPI,AIAPI,AuthAPI,BookstoreAPI,FinanceAPI,MessagesAPI,NotificationsAPI,RecommendationAPI,SearchAPI,SharedAPI,SocialAPI,StatsAPI,SystemAPI api
    class ResponseHelper,Validator,CommonDTOs,Converters shared
    class UserService,ContentService,SocialService,AIService,FinanceService service
    class UserRepo,ContentRepo,SocialRepo repo
    class MongoDB,Redis,Milvus,MinIO db
```

---

## 目录结构图

```mermaid
graph LR
    subgraph "api/v1/"
        A[admin/]
        B[ai/]
        C[announcements/]
        D[audit/]
        E[auth/]
        F[bookstore/]
        G[examples/]
        H[finance/]
        I[messages/]
        J[notifications/]
        K[reader/]
        L[recommendation/]
        M[search/]
        N[shared/]
        O[social/]
        P[stats/]
        Q[system/]
        R[user/]
        S[writer/]
        T[version_api.go]
    end

    subgraph "标准模块结构"
        M1[handler/]
        M2[dto/]
        M3[routes.go]
        M4[README.md]
    end

    subgraph "shared/ 特殊结构"
        S1[response.go]
        S2[request_validator.go]
        S3[types.go]
        S4[converter.go]
        S5[README.md]
    end

    N --> S1
    N --> S2
    N --> S3
    N --> S4
    N --> S5

    A -.推荐.-> M1
    B -.推荐.-> M1
    C -.推荐.-> M1
    D -.推荐.-> M1
    E -.推荐.-> M1
    F -.推荐.-> M1
    H -.推荐.-> M1
    I -.推荐.-> M1
    J -.推荐.-> M1
    L -.推荐.-> M1
    M -.推荐.-> M1
    O -.推荐.-> M1
    P -.推荐.-> M1
    Q -.推荐.-> M1

    K -.待重构.-> M1
    R -.待重构.-> M1
    S -.待重构.-> M1

    classDef active fill:#c8e6c9
    classDef pending fill:#ffccbc
    classDef sharedClass fill:#b2dfdb
    classDef standard fill:#e1bee7

    class A,B,C,D,E,F,H,I,J,L,M,O,P,Q active
    class K,R,S pending
    class N sharedClass
    class M1,M2,M3,M4,S1,S2,S3,S4,S5 standard
```

---

## 请求处理流程图

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant Router as 路由器
    participant MW1 as 全局中间件
    participant MW2 as 认证中间件
    participant Handler as API Handler
    participant Service as 服务层
    participant Repo as 数据访问层
    participant DB as 数据库

    Client->>Router: HTTP请求
    Router->>MW1: CORS检查
    MW1->>MW1: 日志记录
    MW1->>MW1: 恢复处理

    MW1->>MW2: 认证检查
    alt Token有效
        MW2->>MW2: 设置用户上下文
        MW2->>Handler: 调用Handler

        Handler->>Handler: 参数验证
        alt 验证失败
            Handler->>Client: 400 Bad Request
        else 验证成功
            Handler->>Service: 调用业务逻辑
            Service->>Repo: 查询数据
            Repo->>DB: 数据库操作
            DB-->>Repo: 返回数据
            Repo-->>Service: 返回实体
            Service-->>Handler: 返回DTO
            Handler->>Handler: 格式化响应
            Handler->>Client: 200 OK + 数据
        end
    else Token无效
        MW2->>Client: 401 Unauthorized
    end

    Note over Client,DB: 成功流程

    Client->>Router: HTTP请求
    Router->>MW1: 全局中间件处理
    MW1->>MW2: 认证中间件
    MW2->>Handler: 执行Handler
    Handler->>Service: 业务逻辑
    Service->>Repo: 数据操作
    Repo->>DB: 数据库访问
    DB-->>Repo: 返回结果
    Repo-->>Service: 返回实体
    Service-->>Handler: 返回DTO
    Handler->>Client: 统一响应格式
```

---

## 错误处理流程图

```mermaid
graph TB
    Start[收到请求] --> Validate{参数验证}

    Validate -->|失败| BadRequest[400 Bad Request<br/>参数错误]
    Validate -->|成功| Auth{认证检查}

    Auth -->|失败| Unauthorized[401 Unauthorized<br/>未认证]
    Auth -->|成功| Permission{权限检查}

    Permission -->|失败| Forbidden[403 Forbidden<br/>无权限]
    Permission -->|成功| Business{业务逻辑}

    Business -->|资源不存在| NotFound[404 Not Found<br/>资源不存在]
    Business -->|资源冲突| Conflict[409 Conflict<br/>资源已存在]
    Business -->|验证失败| Unprocessable[422 Unprocessable Entity<br/>业务规则验证失败]
    Business -->|请求过多| RateLimit[429 Too Many Requests<br/>超过限流]
    Business -->|成功| Success[200/201 OK/Created]

    Business -->|系统错误| InternalError[500 Internal Server Error<br/>服务器错误]
    InternalError --> Log{记录日志}
    Log --> Notify{通知告警}
    Notify --> ErrorResponse[返回错误响应]

    BadRequest --> ErrorResponse
    Unauthorized --> ErrorResponse
    Forbidden --> ErrorResponse
    NotFound --> ErrorResponse
    Conflict --> ErrorResponse
    Unprocessable --> ErrorResponse
    RateLimit --> ErrorResponse
    Success --> SuccessResponse[返回成功响应]

    ErrorResponse --> End[结束]
    SuccessResponse --> End

    classDef error fill:#ffcdd2
    classDef success fill:#c8e6c9
    classDef process fill:#fff9c4

    class BadRequest,Unauthorized,Forbidden,NotFound,Conflict,Unprocessable,RateLimit,InternalError error
    class Success success
    class Validate,Auth,Permission,Business,Log,Notify process
```

---

## 重构前后对比图

### 架构对比

```mermaid
graph TB
    subgraph "重构前 - 按角色划分"
        OldAdmin[admin<br/>管理员API]
        OldUser[user<br/>用户API]
        OldReader[reader<br/>读者API]
        OldWriter[writer<br/>作者API]
    end

    subgraph "问题"
        P1[功能分散]
        P2[代码重复]
        P3[边界不清]
    end

    subgraph "重构后 - 按功能划分"
        NewUserManagement[user-management<br/>用户管理]
        NewContent[content-management<br/>内容管理]
        NewSocial[social<br/>社交功能]
        NewFinance[finance<br/>财务管理]
        NewAdmin[admin<br/>系统管理]
    end

    OldAdmin -.-> P1
    OldUser -.-> P1
    OldReader -.-> P2
    OldWriter -.-> P2

    P1 -->|重构| NewUserManagement
    P2 -->|重构| NewContent
    P3 -->|重构| NewSocial

    classDef old fill:#ffccbc
    classDef problem fill:#ffe0b2
    classDef new fill:#c8e6c9

    class OldAdmin,OldUser,OldReader,OldWriter old
    class P1,P2,P3 problem
    class NewUserManagement,NewContent,NewSocial,NewFinance,NewAdmin new
```

### 代码对比

```mermaid
graph LR
    subgraph "重构前代码统计"
        OldFiles[API文件数]
        OldDup[重复代码行]
        OldDTO[DTO定义位置]
        OldModules[按角色模块]
    end

    subgraph "改进"
        Arrow1[↓ 36%]
        Arrow2[↓ 80%]
        Arrow3[↓ 67%]
        Arrow4[→ 功能划分]
    end

    subgraph "重构后代码统计"
        NewFiles[API文件数]
        NewDup[重复代码行]
        NewDTO[DTO定义位置]
        NewModules[按功能模块]
    end

    OldFiles -->|156| Arrow1
    Arrow1 -->|~100| NewFiles

    OldDup -->|~2900| Arrow2
    Arrow2 -->|<500| NewDup

    OldDTO -->|3+处| Arrow3
    Arrow3 -->|1处| NewDTO

    OldModules -->|8个角色模块| Arrow4
    Arrow4 -->|12个功能模块| NewModules

    classDef old fill:#ffcdd2
    classDef arrow fill:#fff9c4
    classDef new fill:#c8e6c9

    class OldFiles,OldDup,OldDTO,OldModules old
    class Arrow1,Arrow2,Arrow3,Arrow4 arrow
    class NewFiles,NewDup,NewDTO,NewModules new
```

### 目录结构对比

```mermaid
graph TB
    subgraph "重构前目录结构"
        OldDir1[api/v1/]
        OldDir2[├── admin/<br/>├── user/<br/>├── reader/<br/>├── writer/<br/>├── shared/]
    end

    subgraph "重构后目录结构"
        NewDir1[api/v1/]
        NewDir2[├── admin/<br/>├── user-management/<br/>├── content-management/<br/>├── social/<br/>├── finance/<br/>├── shared/]
    end

    OldDir1 -.演变.-> NewDir1
    OldDir2 -.重组.-> NewDir2

    classDef old fill:#ffccbc
    classDef new fill:#c8e6c9

    class OldDir1,OldDir2 old
    class NewDir1,NewDir2 new
```

---

## 模块依赖关系图

```mermaid
graph TB
    subgraph "业务模块层"
        Admin[admin]
        AI[ai]
        Auth[auth]
        Bookstore[bookstore]
        Finance[finance]
        Messages[messages]
        Notifications[notifications]
        Recommendation[recommendation]
        Search[search]
        Social[social]
        Stats[stats]
        System[system]
    end

    subgraph "共享层"
        Shared[shared]
    end

    subgraph "中间件层"
        AuthMW[auth middleware]
        PermissionMW[permission middleware]
        RateLimitMW[rate limit middleware]
        LoggerMW[logger middleware]
    end

    subgraph "服务层"
        UserService[user service]
        ContentService[content service]
        SocialService[social service]
        AIService[AI service]
        FinanceService[finance service]
    end

    subgraph "工具层"
        Errors[errors package]
        Utils[utils package]
        Models[models package]
    end

    Admin --> Shared
    AI --> Shared
    Auth --> Shared
    Bookstore --> Shared
    Finance --> Shared
    Messages --> Shared
    Notifications --> Shared
    Recommendation --> Shared
    Search --> Shared
    Social --> Shared
    Stats --> Shared
    System --> Shared

    Admin --> UserService
    Auth --> UserService
    Bookstore --> ContentService
    Messages --> SocialService
    Notifications --> SocialService
    AI --> AIService
    Finance --> FinanceService
    Social --> SocialService

    AuthMW --> Shared
    PermissionMW --> Shared
    RateLimitMW --> Shared
    LoggerMW --> Shared

    UserService --> Errors
    ContentService --> Errors
    SocialService --> Errors
    AIService --> Errors
    FinanceService --> Errors

    Shared --> Utils
    Shared --> Models

    classDef business fill:#e8f5e9
    classDef shared fill:#fff9c4
    classDef middleware fill:#e1f5fe
    classDef service fill:#f3e5f5
    classDef utils fill:#fce4ec

    class Admin,AI,Auth,Bookstore,Finance,Messages,Notifications,Recommendation,Search,Social,Stats,System business
    class Shared shared
    class AuthMW,PermissionMW,RateLimitMW,LoggerMW middleware
    class UserService,ContentService,SocialService,AIService,FinanceService service
    class Errors,Utils,Models utils
```

---

## 数据流图

```mermaid
sequenceDiagram
    participant Client
    participant Router
    participant Middleware
    participant Handler
    participant Validator
    participant Service
    participant Repository
    participant Database

    Client->>Router: POST /api/v1/auth/login
    Router->>Middleware: 全局中间件链

    Middleware->>Middleware: 1. CORS检查
    Middleware->>Middleware: 2. 日志记录
    Middleware->>Middleware: 3. 恢复处理

    Middleware->>Handler: 调用Handler

    Handler->>Validator: 验证请求体
    Note over Validator: binding:"required"

    alt 验证失败
        Validator->>Client: 400 Bad Request
    else 验证成功
        Handler->>Service: 调用业务逻辑
        Note over Service: Login(username, password)

        Service->>Repository: 查询用户
        Repository->>Database: MongoDB查询
        Database-->>Repository: 返回用户实体
        Repository-->>Service: 返回用户领域对象

        Service->>Service: 验证密码
        Service->>Service: 生成JWT Token

        Service-->>Handler: 返回TokenDTO
        Handler->>Handler: 转换为响应格式

        Handler->>Client: 200 OK
        Note over Client: {<br/>  "code": 200,<br/>  "message": "登录成功",<br/>  "data": {<br/>    "token": "...",<br/>    "expires_in": 3600<br/>  }<br/>}
    end
```

---

## 分页请求流程图

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Service
    participant Repository
    participant Database

    Client->>Handler: GET /api/v1/users?page=1&page_size=20
    Note over Client: 查询参数

    Handler->>Handler: 验证查询参数
    Note over Handler: ValidateQueryParams()

    Handler->>Service: 调用列表服务
    Note over Service: ListUsers(page, pageSize, keyword, role)

    Service->>Repository: 查询数据
    Repository->>Database: 执行分页查询
    Note over Database: find().skip((page-1)*pageSize).limit(pageSize)

    Database-->>Repository: 返回数据列表
    Repository->>Database: 查询总数
    Database-->>Repository: 返回总数

    Repository-->>Service: 返回列表和总数
    Service->>Service: 业务逻辑处理
    Service-->>Handler: 返回结果

    Handler->>Handler: 格式化分页响应
    Note over Handler: Paginated(data, total, page, pageSize)

    Handler->>Client: 返回分页响应
    Note over Client: {<br/>  "code": 200,<br/>  "data": [...],<br/>  "pagination": {<br/>    "page": 1,<br/>    "page_size": 20,<br/>    "total": 100,<br/>    "total_pages": 5<br/>  }<br/>}
```

---

## 错误处理详细流程图

```mermaid
graph TB
    Start[Handler执行] --> Try{try-catch}

    Try -->|正常执行| BusinessLogic[业务逻辑]
    Try -->|Panic| CatchPanic[捕获panic]

    BusinessLogic --> ServiceCall[调用服务]
    ServiceCall --> CheckError{检查错误}

    CheckError -->|无错误| Success[返回成功]
    CheckError -->|有错误| ClassifyError{分类错误}

    ClassifyError -->|NotFound| NotFoundError[404错误]
    ClassifyError -->|ValidationError| ValidationErr[400错误]
    ClassifyError -->|PermissionError| PermissionErr[403错误]
    ClassifyError -->|ConflictError| ConflictErr[409错误]
    ClassifyError -->|其他| InternalErr[500错误]

    CatchPanic --> LogPanic[记录panic日志]
    LogPanic --> InternalErr

    NotFoundError --> Response1[shared.NotFound]
    ValidationErr --> Response2[shared.BadRequest]
    PermissionErr --> Response3[shared.Forbidden]
    ConflictErr --> Response4[shared.Error 409]
    InternalErr --> Response5[shared.InternalError]

    Success --> End[返回响应]
    Response1 --> End
    Response2 --> End
    Response3 --> End
    Response4 --> End
    Response5 --> LogError[记录错误日志]
    LogError --> End

    classDef normal fill:#c8e6c9
    classDef error fill:#ffcdd2
    classDef process fill:#fff9c4
    classDef response fill:#e1f5fe

    class BusinessLogic,ServiceCall,Success normal
    class NotFoundError,ValidationErr,PermissionErr,ConflictErr,InternalErr,CatchPanic error
    class Try,CheckError,ClassifyError,LogPanic,LogError process
    class Response1,Response2,Response3,Response4,Response5 response
```

---

## 文档说明

### 图例说明

- **绿色**: 表示活跃的/正确的/新的模块或流程
- **红色**: 表示错误的/废弃的/旧的模块或流程
- **黄色**: 表示处理过程或中间状态
- **蓝色**: 表示响应或返回值

### 使用说明

1. **API层整体架构图**: 展示完整的系统架构和各层之间的关系
2. **目录结构图**: 展示API层的目录组织和模块划分
3. **请求处理流程图**: 展示从请求到响应的完整流程
4. **错误处理流程图**: 展示各种错误情况的处理逻辑
5. **重构前后对比图**: 对比重构前后的架构差异
6. **模块依赖关系图**: 展示各模块之间的依赖关系
7. **数据流图**: 展示数据在系统中的流动

### 更新记录

| 日期 | 版本 | 变更内容 |
|------|------|----------|
| 2026-02-26 | v1.0 | 初始版本，创建所有架构图 |

---

**维护者**: Backend Architecture Team
**最后更新**: 2026-02-26
