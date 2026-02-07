# Qingyu Backend 组件分析

## 1. 组件依赖图

### 1.1 完整依赖关系图

```mermaid
graph LR
    subgraph "API Layer"
        BA[BookstoreAPI]
        RA[ReaderAPI]
        WA[WriterAPI]
        SA[SocialAPI]
        AA[AdminAPI]
    end

    subgraph "Service Layer"
        BSvc[BookstoreService]
        RSvc[ReaderService]
        WSvc[WriterService]
        SSvc[SocialService]
        ASvc[AIService]
        FSvc[FinanceService]
        CSvc[ComplianceService]
        Container[ServiceContainer]
        Events[EventBus]
    end

    subgraph "Repository Layer"
        BRepo[BookstoreRepository]
        RRepo[ReaderRepository]
        WRepo[WriterRepository]
        SRepo[SocialRepository]
        FRepo[FinanceRepository]
    end

    subgraph "Infrastructure"
        Cache[Cache]
        Auth[Auth]
        Queue[MessageQueue]
    end

    subgraph "Data Layer"
        Mongo[(MongoDB)]
        Redis[(Redis)]
        Milvus[(Milvus)]
        Minio[(MinIO)]
        AIgrpc[(AI gRPC)]
    end

    BA --> BSvc
    RA --> RSvc
    WA --> WSvc
    SA --> SSvc
    AA --> ASvc

    BSvc --> Container
    RSvc --> Container
    WSvc --> Container
    SSvc --> Container
    ASvc --> Container
    FSvc --> Container

    BSvc --> BRepo
    RSvc --> RRepo
    WSvc --> WRepo
    SSvc --> SRepo
    FSvc --> FRepo

    RSvc --> BSvc
    SSvc --> RSvc
    WSvc --> BSvc
    WSvc --> ASvc
    ASvc --> Events
    CSvc --> WSvc

    BSvc --> Cache
    RSvc --> Cache
    SSvc --> Cache

    BRepo --> Mongo
    RRepo --> Mongo
    WRepo --> Mongo
    SRepo --> Mongo
    FRepo --> Mongo

    Cache --> Redis
    ASvc --> AIgrpc
    WSvc --> Milvus
    BSvc --> Minio
```

### 1.2 模块间依赖矩阵

|  | Bookstore | Reader | Writer | Social | AI | Finance |
|---|---|---:|---|---|---|---|
| **Bookstore** | ✓ | - | - | - | - | - |
| **Reader** | → | ✓ | - | - | - | - |
| **Writer** | → | - | ✓ | - | → | - |
| **Social** | - | → | - | ✓ | - | - |
| **AI** | - | - | - | - | ✓ | - |
| **Finance** | - | - | → | - | - | ✓ |

→ 表示依赖关系

## 2. 关键接口分析

### 2.1 核心API端点

| API路径 | 方法 | 功能 | 调用链 |
|---------|------|------|--------|
| `/api/v1/bookstore/books` | GET | 获取书籍列表 | API → BookstoreSvc → BookstoreRepo → MongoDB |
| `/api/v1/reader/history` | GET | 获取阅读历史 | API → ReaderSvc → ReaderRepo → MongoDB |
| `/api/v1/writer/chapter` | POST | 创建章节 | API → WriterSvc → WriterRepo → MongoDB |
| `/api/v1/social/comment` | POST | 发表评论 | API → SocialSvc → SocialRepo → MongoDB → EventBus |
| `/api/v1/ai/chat` | POST | AI对话 | API → AISvc → AI gRPC |

### 2.2 Service层接口

```go
// 核心服务接口示例
type BookstoreService interface {
    GetBook(ctx context.Context, id string) (*models.Book, error)
    ListBooks(ctx context.Context, filter BookFilter) ([]*models.Book, error)
    CreateBook(ctx context.Context, book *models.Book) error
    UpdateBook(ctx context.Context, id string, book *models.Book) error
    DeleteBook(ctx context.Context, id string) error
}

type ReaderService interface {
    GetReadingHistory(ctx context.Context, userID string) ([]*models.ReadingHistory, error)
    AddToBookshelf(ctx context.Context, userID, bookID string) error
    UpdateProgress(ctx context.Context, userID, bookID string, progress int) error
}
```

### 2.3 Repository层接口

```go
// Repository接口示例
type BookRepository interface {
    FindByID(ctx context.Context, id string) (*models.Book, error)
    Find(ctx context.Context, filter BookFilter) ([]*models.Book, error)
    Create(ctx context.Context, book *models.Book) error
    Update(ctx context.Context, book *models.Book) error
    Delete(ctx context.Context, id string) error
    Exists(ctx context.Context, id string) (bool, error)
}
```

## 3. 数据流分析

### 3.1 典型请求流程：读者阅读小说

```mermaid
sequenceDiagram
    participant User as 用户
    participant API as ReaderAPI
    participant Auth as AuthMiddleware
    participant Svc as ReaderService
    participant BSvc as BookstoreService
    participant Repo as ReaderRepository
    participant Cache as Redis Cache
    participant DB as MongoDB

    User->>API: GET /api/v1/reader/books/{id}/read
    API->>Auth: 验证Token
    Auth-->>API: 用户信息
    API->>Svc: GetBookForReading(bookID)
    Svc->>Cache: 检查缓存
    alt 缓存命中
        Cache-->>Svc: 书籍数据
    else 缓存未命中
        Svc->>BSvc: GetBook(bookID)
        BSvc->>Repo: FindByID(bookID)
        Repo->>DB: 查询
        DB-->>Repo: 书籍数据
        Repo-->>BSvc: 书籍
        BSvc-->>Svc: 书籍
        Svc->>Cache: 写入缓存
    end
    Svc->>Repo: GetReadingProgress(userID, bookID)
    Repo->>DB: 查询阅读进度
    DB-->>Repo: 阅读进度
    Repo-->>Svc: 阅读进度
    Svc-->>API: 书籍+进度
    API-->>User: 返回章节内容
```

### 3.2 典型请求流程：作者发布章节

```mermaid
sequenceDiagram
    participant Writer as 作者
    participant API as WriterAPI
    participant Svc as WriterService
    participant AI as AIService
    participant EventBus as EventBus
    participant NotifSvc as NotificationService
    participant Repo as WriterRepository
    participant DB as MongoDB

    Writer->>API: POST /api/v1/writer/chapter
    API->>Svc: CreateChapter(chapter)
    Svc->>AI: 请求AI辅助内容
    AI-->>Svc: AI建议
    Svc->>Repo: SaveChapter(chapter)
    Repo->>DB: 插入章节
    DB-->>Repo: 章节ID
    Repo-->>Svc: 章节
    Svc->>EventBus: Publish(ChapterCreatedEvent)
    EventBus->>NotifSvc: 订阅者处理
    NotifSvc->>Repo: 创建推送通知
    Svc-->>API: 章节创建成功
    API-->>Writer: 返回章节ID
```

### 3.3 事件流分析

```mermaid
graph TB
    subgraph "事件发布者"
        WriterSvc[WriterService]
        ReaderSvc[ReaderService]
        SocialSvc[SocialService]
    end

    subgraph "事件总线"
        EventBus[EventBus]
    end

    subgraph "事件订阅者"
        NotifSvc[NotificationService]
        StatSvc[StatisticsService]
        AuditSvc[AuditService]
        AIRecSvc[AIRecommendService]
    end

    WriterSvc -->|ChapterCreated| EventBus
    WriterSvc -->|BookPublished| EventBus
    ReaderSvc -->|BookAdded| EventBus
    SocialSvc -->|CommentCreated| EventBus
    SocialSvc -->|LikeAction| EventBus

    EventBus -->|ChapterCreated| NotifSvc
    EventBus -->|BookPublished| NotifSvc
    EventBus -->|CommentCreated| NotifSvc
    EventBus -->|LikeAction| StatSvc
    EventBus -->|BookPublished| AIRecSvc
    EventBus --> AllEvents --> AuditSvc
```

## 4. 组件耦合度分析

### 4.1 高耦合组件

| 组件 | 依赖数量 | 说明 | 建议 |
|------|----------|------|------|
| **WriterService** | 5+ | 依赖Bookstore, AI, Events, Notification等 | 考虑拆分或使用领域事件 |
| **ServiceContainer** | 全部 | 集中管理所有服务 | 保持现状，但加强文档 |
| **QuotaMiddleware** | Service层 | 中间件直接依赖业务服务 | 引入抽象层 |

### 4.2 低耦合组件

| 组件 | 说明 |
|------|------|
| **BookstoreService** | 核心模块，依赖最少 |
| **Repository层** | 接口隔离良好 |

---

**文档版本**: v1.0
**最后更新**: 2026-02-07
**维护者**: yukin371
