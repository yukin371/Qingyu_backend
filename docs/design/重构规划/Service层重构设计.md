# Service层重构设计文档

## 1. 概述

### 1.1 重构目标

基于新的Repository层架构，重构Service层以实现更清晰的业务逻辑分离、更好的依赖注入和更强的可测试性。

### 1.2 现状分析

当前Service层存在以下问题：
- 直接依赖具体的Repository实现，耦合度高
- 业务逻辑与数据访问逻辑混合
- 缺乏统一的服务接口规范
- 事务管理分散在各个服务中
- 错误处理不统一

### 1.3 重构原则

- **依赖注入**：通过接口依赖Repository层，实现松耦合
- **单一职责**：每个Service专注于特定的业务领域
- **事务边界**：在Service层管理事务边界
- **错误转换**：将Repository错误转换为业务错误
- **API兼容**：保持现有Service API的向后兼容性

## 2. 架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                    API Layer                                │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   AI Router     │  │ Document Router │  │ User Router  │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Service Layer (Refactored)                   │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  Base Service   │  │Service Factory  │  │ Validator    │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  User Service   │  │   AI Service    │  │Doc Service   │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Repository Layer                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  User Repository│  │  AI Repository  │  │Doc Repository│ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 核心组件

#### 2.2.1 BaseService接口

```go
type BaseService[T any, ID comparable] interface {
    // 基础业务操作
    Create(ctx context.Context, entity *T) (*T, error)
    GetByID(ctx context.Context, id ID) (*T, error)
    Update(ctx context.Context, id ID, updates map[string]interface{}) (*T, error)
    Delete(ctx context.Context, id ID) error
    
    // 查询操作
    List(ctx context.Context, filter ServiceFilter) (*PagedResult[T], error)
    Count(ctx context.Context, filter ServiceFilter) (int64, error)
    Exists(ctx context.Context, id ID) (bool, error)
    
    // 批量操作
    BatchCreate(ctx context.Context, entities []*T) ([]*T, error)
    BatchUpdate(ctx context.Context, ids []ID, updates map[string]interface{}) error
    BatchDelete(ctx context.Context, ids []ID) error
    
    // 验证
    Validate(ctx context.Context, entity *T) error
}
```

#### 2.2.2 ServiceFilter和分页结果

```go
type ServiceFilter interface {
    ToRepositoryFilter() repository.Filter
    Validate() error
}

type PagedResult[T any] struct {
    Data       []*T  `json:"data"`
    Total      int64 `json:"total"`
    Page       int   `json:"page"`
    PageSize   int   `json:"pageSize"`
    TotalPages int   `json:"totalPages"`
}

type BaseServiceFilter struct {
    Page     int               `json:"page" validate:"min=1"`
    PageSize int               `json:"pageSize" validate:"min=1,max=100"`
    Sort     map[string]string `json:"sort,omitempty"`
    Search   string            `json:"search,omitempty"`
}
```

## 3. 专用Service设计

### 3.1 UserService重构

```go
type UserService interface {
    BaseService[system.User, string]
    
    // 用户认证相关
    Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error)
    Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
    Logout(ctx context.Context, userID string) error
    RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
    
    // 用户管理
    GetProfile(ctx context.Context, userID string) (*UserProfile, error)
    UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*UserProfile, error)
    ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error
    ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
    
    // 用户状态管理
    ActivateUser(ctx context.Context, userID string) error
    DeactivateUser(ctx context.Context, userID string) error
    GetActiveUsers(ctx context.Context, limit int) ([]*UserResponse, error)
}

type UserServiceImpl struct {
    userRepo        repository.UserRepository
    roleRepo        repository.RoleRepository
    transactionMgr  repository.TransactionManager
    validator       Validator
    passwordHasher  PasswordHasher
    tokenManager    TokenManager
    eventPublisher  EventPublisher
}
```

### 3.2 AIService重构

```go
type AIService interface {
    // 内容生成
    GenerateContent(ctx context.Context, req *GenerateContentRequest) (*GenerateContentResponse, error)
    GenerateContentStream(ctx context.Context, req *GenerateContentRequest) (<-chan *StreamResponse, error)
    
    // 内容分析
    AnalyzeContent(ctx context.Context, req *AnalyzeContentRequest) (*AnalyzeContentResponse, error)
    
    // 文本优化
    OptimizeText(ctx context.Context, req *OptimizeTextRequest) (*OptimizeTextResponse, error)
    ContinueWriting(ctx context.Context, req *ContinueWritingRequest) (*GenerateContentResponse, error)
    
    // 大纲生成
    GenerateOutline(ctx context.Context, req *GenerateOutlineRequest) (*OutlineResponse, error)
    
    // 上下文管理
    GetContextInfo(ctx context.Context, projectID, chapterID string) (*ContextInfo, error)
    UpdateContext(ctx context.Context, req *UpdateContextRequest) error
}

type AIServiceImpl struct {
    aiRepo          repository.AIRepository
    contextService  ContextService
    adapterManager  adapter.AdapterManager
    validator       Validator
    eventPublisher  EventPublisher
}
```

### 3.3 ChatService重构

```go
type ChatService interface {
    // 聊天会话管理
    CreateSession(ctx context.Context, req *CreateSessionRequest) (*ChatSessionResponse, error)
    GetSession(ctx context.Context, sessionID string) (*ChatSessionResponse, error)
    UpdateSession(ctx context.Context, sessionID string, req *UpdateSessionRequest) (*ChatSessionResponse, error)
    DeleteSession(ctx context.Context, sessionID string) error
    ListSessions(ctx context.Context, projectID string, filter *SessionFilter) (*PagedResult[ChatSessionResponse], error)
    
    // 聊天对话
    SendMessage(ctx context.Context, req *SendMessageRequest) (*ChatResponse, error)
    SendMessageStream(ctx context.Context, req *SendMessageRequest) (<-chan *StreamChatResponse, error)
    GetChatHistory(ctx context.Context, sessionID string, filter *MessageFilter) (*PagedResult[ChatMessageResponse], error)
    
    // 统计信息
    GetChatStatistics(ctx context.Context, projectID string) (*ChatStatistics, error)
}

type ChatServiceImpl struct {
    aiRepo         repository.AIRepository
    aiService      AIService
    contextService ContextService
    validator      Validator
    eventPublisher EventPublisher
}
```

### 3.4 DocumentService重构

```go
type DocumentService interface {
    BaseService[document.Document, string]
    
    // 文档管理
    CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*DocumentResponse, error)
    GetDocument(ctx context.Context, documentID string) (*DocumentResponse, error)
    UpdateDocument(ctx context.Context, documentID string, req *UpdateDocumentRequest) (*DocumentResponse, error)
    DeleteDocument(ctx context.Context, documentID string) error
    
    // 内容管理
    UpdateContent(ctx context.Context, documentID string, content string) (*DocumentResponse, error)
    GetContent(ctx context.Context, documentID string) (string, error)
    
    // 版本管理
    CreateVersion(ctx context.Context, documentID string, req *CreateVersionRequest) (*VersionResponse, error)
    GetVersions(ctx context.Context, documentID string) ([]*VersionResponse, error)
    RestoreVersion(ctx context.Context, documentID, versionID string) (*DocumentResponse, error)
    
    // 项目关联
    GetDocumentsByProject(ctx context.Context, projectID string, filter *DocumentFilter) (*PagedResult[DocumentResponse], error)
    GetDocumentsByUser(ctx context.Context, userID string, filter *DocumentFilter) (*PagedResult[DocumentResponse], error)
}

type DocumentServiceImpl struct {
    docRepo        repository.DocumentRepository
    versionRepo    repository.VersionRepository
    transactionMgr repository.TransactionManager
    validator      Validator
    eventPublisher EventPublisher
}
```

## 4. 依赖注入设计

### 4.1 ServiceContainer

```go
type ServiceContainer interface {
    // 注册服务
    RegisterService(name string, factory ServiceFactory) error
    RegisterSingleton(name string, instance interface{}) error
    
    // 获取服务
    GetService(name string) (interface{}, error)
    GetUserService() UserService
    GetAIService() AIService
    GetChatService() ChatService
    GetDocumentService() DocumentService
    
    // 生命周期管理
    Initialize() error
    Shutdown() error
}

type ServiceFactory func(container ServiceContainer) (interface{}, error)

type DefaultServiceContainer struct {
    services   map[string]ServiceFactory
    singletons map[string]interface{}
    instances  map[string]interface{}
    mu         sync.RWMutex
}
```

### 4.2 服务配置

```go
type ServiceConfig struct {
    // Repository配置
    RepositoryConfig repository.Config `yaml:"repository"`
    
    // 验证配置
    ValidationConfig ValidationConfig `yaml:"validation"`
    
    // 事件配置
    EventConfig EventConfig `yaml:"event"`
    
    // 缓存配置
    CacheConfig CacheConfig `yaml:"cache"`
}

func InitializeServices(config *ServiceConfig) (ServiceContainer, error) {
    container := NewServiceContainer()
    
    // 注册Repository层
    repoFactory, err := repository.NewRepositoryFactory(config.RepositoryConfig)
    if err != nil {
        return nil, err
    }
    container.RegisterSingleton("repositoryFactory", repoFactory)
    
    // 注册服务
    container.RegisterService("userService", func(c ServiceContainer) (interface{}, error) {
        repoFactory := c.GetService("repositoryFactory").(repository.RepositoryFactory)
        return NewUserService(repoFactory.CreateUserRepository()), nil
    })
    
    return container, nil
}
```

## 5. 事务管理

### 5.1 Service层事务边界

```go
type TransactionalService interface {
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// 在Service实现中使用事务
func (s *UserServiceImpl) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
    var result *UserResponse
    
    err := s.transactionMgr.ExecuteTransaction(ctx, func(txCtx repository.TransactionContext) error {
        // 1. 验证用户数据
        if err := s.validator.Validate(req); err != nil {
            return NewValidationError("用户数据验证失败", err)
        }
        
        // 2. 检查用户是否已存在
        exists, err := txCtx.UserRepository().ExistsByEmail(txCtx, req.Email)
        if err != nil {
            return err
        }
        if exists {
            return NewBusinessError("邮箱已被注册")
        }
        
        // 3. 创建用户
        user := &system.User{
            Username: req.Username,
            Email:    req.Email,
            Password: s.passwordHasher.Hash(req.Password),
        }
        
        if err := txCtx.UserRepository().Create(txCtx, user); err != nil {
            return err
        }
        
        // 4. 分配默认角色
        defaultRole, err := txCtx.RoleRepository().GetDefaultRole(txCtx)
        if err != nil {
            return err
        }
        
        if err := txCtx.RoleRepository().AssignRole(txCtx, user.ID, defaultRole.ID); err != nil {
            return err
        }
        
        // 5. 发布事件
        s.eventPublisher.Publish(&UserRegisteredEvent{
            UserID: user.ID,
            Email:  user.Email,
        })
        
        result = s.toUserResponse(user)
        return nil
    })
    
    return result, err
}
```

## 6. 错误处理

### 6.1 Service层错误类型

```go
type ServiceError struct {
    Type      ServiceErrorType `json:"type"`
    Message   string           `json:"message"`
    Code      string           `json:"code"`
    Details   map[string]interface{} `json:"details,omitempty"`
    Cause     error            `json:"-"`
    Timestamp time.Time        `json:"timestamp"`
}

type ServiceErrorType string

const (
    ServiceErrorTypeValidation   ServiceErrorType = "VALIDATION"
    ServiceErrorTypeBusiness     ServiceErrorType = "BUSINESS"
    ServiceErrorTypeNotFound     ServiceErrorType = "NOT_FOUND"
    ServiceErrorTypePermission   ServiceErrorType = "PERMISSION"
    ServiceErrorTypeConflict     ServiceErrorType = "CONFLICT"
    ServiceErrorTypeInternal     ServiceErrorType = "INTERNAL"
)

func NewValidationError(message string, cause error) *ServiceError
func NewBusinessError(message string) *ServiceError
func NewNotFoundError(resource string, id string) *ServiceError
func NewPermissionError(message string) *ServiceError
```

### 6.2 错误转换

```go
type ErrorTranslator interface {
    TranslateRepositoryError(err error) *ServiceError
    TranslateValidationError(err error) *ServiceError
}

func (s *BaseServiceImpl) handleError(err error) error {
    if err == nil {
        return nil
    }
    
    // 如果已经是Service错误，直接返回
    if serviceErr, ok := err.(*ServiceError); ok {
        return serviceErr
    }
    
    // 转换Repository错误
    if repoErr, ok := err.(*repository.RepositoryError); ok {
        return s.errorTranslator.TranslateRepositoryError(repoErr)
    }
    
    // 其他错误转换为内部错误
    return NewInternalError("内部服务错误", err)
}
```

## 7. 验证框架

### 7.1 Validator接口

```go
type Validator interface {
    Validate(obj interface{}) error
    ValidateStruct(obj interface{}) error
    ValidateField(field interface{}, tag string) error
}

type ValidationError struct {
    Field   string `json:"field"`
    Tag     string `json:"tag"`
    Value   string `json:"value"`
    Message string `json:"message"`
}

type ValidationErrors []*ValidationError

func (ve ValidationErrors) Error() string
```

### 7.2 业务规则验证

```go
type BusinessRuleValidator interface {
    ValidateUserRegistration(ctx context.Context, req *RegisterRequest) error
    ValidateDocumentCreation(ctx context.Context, req *CreateDocumentRequest) error
    ValidateAIRequest(ctx context.Context, req *GenerateContentRequest) error
}

// 在Service中使用
func (s *UserServiceImpl) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
    // 结构验证
    if err := s.validator.Validate(req); err != nil {
        return nil, NewValidationError("请求参数验证失败", err)
    }
    
    // 业务规则验证
    if err := s.businessValidator.ValidateUserRegistration(ctx, req); err != nil {
        return nil, err
    }
    
    // ... 继续处理
}
```

## 8. 事件系统

### 8.1 事件发布

```go
type Event interface {
    GetType() string
    GetTimestamp() time.Time
    GetData() interface{}
}

type EventPublisher interface {
    Publish(event Event) error
    PublishAsync(event Event) error
}

type EventSubscriber interface {
    Subscribe(eventType string, handler EventHandler) error
    Unsubscribe(eventType string, handler EventHandler) error
}

type EventHandler func(event Event) error

// 事件定义
type UserRegisteredEvent struct {
    UserID    string    `json:"userId"`
    Email     string    `json:"email"`
    Timestamp time.Time `json:"timestamp"`
}

func (e *UserRegisteredEvent) GetType() string { return "user.registered" }
func (e *UserRegisteredEvent) GetTimestamp() time.Time { return e.Timestamp }
func (e *UserRegisteredEvent) GetData() interface{} { return e }
```

## 9. API兼容性保证

### 9.1 适配器模式

```go
// 保持现有API接口
type LegacyAIService struct {
    newService AIService
}

func (s *LegacyAIService) GenerateContent(ctx context.Context, req *GenerateContentRequest) (*GenerateContentResponse, error) {
    // 转换请求格式
    newReq := s.convertRequest(req)
    
    // 调用新服务
    resp, err := s.newService.GenerateContent(ctx, newReq)
    if err != nil {
        return nil, s.convertError(err)
    }
    
    // 转换响应格式
    return s.convertResponse(resp), nil
}
```

### 9.2 渐进式迁移

```go
type ServiceMigrationManager struct {
    legacyServices map[string]interface{}
    newServices    map[string]interface{}
    migrationFlags map[string]bool
}

func (m *ServiceMigrationManager) GetService(name string) interface{} {
    if m.migrationFlags[name] {
        return m.newServices[name]
    }
    return m.legacyServices[name]
}
```

## 10. 测试策略

### 10.1 单元测试

```go
type MockUserRepository struct {
    users map[string]*system.User
}

func (m *MockUserRepository) Create(ctx context.Context, user *system.User) error {
    m.users[user.ID] = user
    return nil
}

func TestUserService_Register(t *testing.T) {
    // 准备Mock
    mockRepo := &MockUserRepository{users: make(map[string]*system.User)}
    mockValidator := &MockValidator{}
    
    // 创建服务
    service := NewUserService(mockRepo, mockValidator)
    
    // 测试用例
    req := &RegisterRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    resp, err := service.Register(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "testuser", resp.Username)
}
```

### 10.2 集成测试

```go
func TestUserServiceIntegration(t *testing.T) {
    // 使用真实的Repository和数据库
    container := setupTestContainer(t)
    userService := container.GetUserService()
    
    // 测试完整的注册流程
    req := &RegisterRequest{
        Username: "integrationtest",
        Email:    "integration@example.com",
        Password: "password123",
    }
    
    resp, err := userService.Register(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    
    // 验证数据库中的数据
    user, err := userService.GetByID(context.Background(), resp.ID)
    assert.NoError(t, err)
    assert.Equal(t, req.Username, user.Username)
}
```

## 11. 监控和指标

### 11.1 服务指标

```go
type ServiceMetrics interface {
    RecordRequest(serviceName, methodName string, duration time.Duration, success bool)
    RecordError(serviceName, methodName string, errorType string)
    GetMetrics() map[string]interface{}
}

// 在Service中使用
func (s *UserServiceImpl) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        s.metrics.RecordRequest("UserService", "Register", duration, true)
    }()
    
    // ... 业务逻辑
}
```

### 11.2 健康检查

```go
type HealthChecker interface {
    CheckHealth(ctx context.Context) *HealthStatus
}

type HealthStatus struct {
    Status     string            `json:"status"`
    Services   map[string]string `json:"services"`
    Timestamp  time.Time         `json:"timestamp"`
    Details    map[string]interface{} `json:"details,omitempty"`
}

func (s *UserServiceImpl) CheckHealth(ctx context.Context) *HealthStatus {
    status := &HealthStatus{
        Status:    "healthy",
        Services:  make(map[string]string),
        Timestamp: time.Now(),
    }
    
    // 检查Repository连接
    if err := s.userRepo.Health(ctx); err != nil {
        status.Status = "unhealthy"
        status.Services["userRepository"] = "unhealthy"
    } else {
        status.Services["userRepository"] = "healthy"
    }
    
    return status
}
```

这个Service层重构设计确保了与新Repository层的完美集成，同时保持了API的向后兼容性和系统的可扩展性。
