# 青羽写作平台 - 导出和发布管理API文档

## 概述

本文档说明了青羽写作平台新增的导出功能和发布管理API的实现。

## 功能列表

### 1. 导出功能 API

#### 1.1 导出文档
- **端点**: `POST /api/v1/writer/documents/:id/export`
- **描述**: 将单个文档导出为指定格式（TXT/MD/DOCX）
- **请求参数**:
  - `id` (路径参数): 文档ID
  - `projectId` (查询参数): 项目ID
  - `format`: 导出格式（txt/md/docx）
  - `includeMeta`: 是否包含元数据
  - `options`: 导出选项（TOC、页码、字体等）

#### 1.2 导出项目
- **端点**: `POST /api/v1/writer/projects/:id/export`
- **描述**: 将整个项目导出为ZIP包
- **请求参数**:
  - `id` (路径参数): 项目ID
  - `includeDocuments`: 是否包含文档
  - `includeCharacters`: 是否包含角色
  - `includeLocations`: 是否包含地点
  - `documentFormats`: 文档导出格式

#### 1.3 查询导出任务状态
- **端点**: `GET /api/v1/writer/exports/:id`
- **描述**: 获取导出任务的详细状态
- **返回信息**:
  - 任务状态（pending/processing/completed/failed）
  - 进度百分比
  - 下载链接（完成后）

#### 1.4 下载导出文件
- **端点**: `GET /api/v1/writer/exports/:id/download`
- **描述**: 下载已完成的导出文件
- **返回**: 文件的签名URL（1小时有效）

#### 1.5 列出导出任务
- **端点**: `GET /api/v1/writer/projects/:projectId/exports`
- **描述**: 获取项目的所有导出任务列表
- **支持分页**: page, pageSize

#### 1.6 删除导出任务
- **端点**: `DELETE /api/v1/writer/exports/:id`
- **描述**: 删除导出任务及其文件

#### 1.7 取消导出任务
- **端点**: `POST /api/v1/writer/exports/:id/cancel`
- **描述**: 取消正在处理或等待中的导出任务

### 2. 发布管理 API

#### 2.1 发布项目
- **端点**: `POST /api/v1/writer/projects/:id/publish`
- **描述**: 将项目发布到指定书城平台
- **请求参数**:
  - `bookstoreId`: 书城ID
  - `categoryId`: 分类ID
  - `tags`: 标签列表
  - `description`: 作品简介
  - `coverImage`: 封面图URL
  - `publishType`: 发布方式（serial/complete）
  - `publishTime`: 定时发布时间（可选）
  - `price`: 定价（可选）
  - `freeChapters`: 免费章节数量
  - `authorNote`: 作者的话
  - `enableComment`: 是否开启评论
  - `enableShare`: 是否允许分享

#### 2.2 取消发布项目
- **端点**: `POST /api/v1/writer/projects/:id/unpublish`
- **描述**: 将项目从书城平台下架

#### 2.3 获取项目发布状态
- **端点**: `GET /api/v1/writer/projects/:id/publication-status`
- **描述**: 获取项目的发布状态和统计信息
- **返回信息**:
  - 是否已发布
  - 书城信息
  - 发布时间
  - 章节统计
  - 阅读统计（阅读量、点赞数、评论数等）

#### 2.4 发布章节
- **端点**: `POST /api/v1/writer/documents/:id/publish`
- **描述**: 发布单个文档（章节）到书城
- **请求参数**:
  - `chapterTitle`: 章节标题
  - `chapterNumber`: 章节序号
  - `isFree`: 是否免费
  - `publishTime`: 定时发布时间（可选）
  - `authorNote`: 章节作者的话

#### 2.5 更新章节发布状态
- **端点**: `PUT /api/v1/writer/documents/:id/publish-status`
- **描述**: 更新章节的发布状态
- **请求参数**:
  - `isPublished`: 是否已发布
  - `publishTime`: 发布时间
  - `unpublishReason`: 取消发布原因
  - `isFree`: 是否免费
  - `chapterNumber`: 章节序号

#### 2.6 批量发布章节
- **端点**: `POST /api/v1/writer/projects/:projectId/documents/batch-publish`
- **描述**: 批量发布多个章节
- **请求参数**:
  - `documentIds`: 文档ID列表（最多50个）
  - `autoNumbering`: 是否自动编号
  - `startNumber`: 起始章节号
  - `isFree`: 是否全部免费
  - `publishTime`: 发布时间

#### 2.7 获取发布记录列表
- **端点**: `GET /api/v1/writer/projects/:projectId/publications`
- **描述**: 获取项目的所有发布记录
- **支持分页**: page, pageSize

#### 2.8 获取发布记录详情
- **端点**: `GET /api/v1/writer/publications/:id`
- **描述**: 根据ID获取发布记录的详细信息

## 文件结构

### 创建的文件列表

#### Service接口层
1. `service/interfaces/export_service.go` - 导出服务接口定义
2. `service/interfaces/publish_service.go` - 发布服务接口定义

#### Service实现层
3. `service/writer/export_service.go` - 导出服务实现
4. `service/writer/publish_service.go` - 发布服务实现

#### API层
5. `api/v1/writer/export_api.go` - 导出API处理器
6. `api/v1/writer/publish_api.go` - 发布API处理器

#### 路由层
7. `router/writer/export_router.go` - 导出功能路由
8. `router/writer/publish_router.go` - 发布管理路由

### 修改的文件列表
1. `router/writer/writer.go` - 添加导出和发布路由的初始化调用
2. `router/writer/writer_router_init.go` - 添加导出和发布服务的创建和注入

## 需要手动实现的接口

为了使这些API完全工作，需要实现以下依赖接口：

### 导出功能所需接口

#### 1. ExportTaskRepository
在 `repository/interfaces/writing/` 或类似位置创建：
```go
type ExportTaskRepository interface {
    Create(ctx context.Context, task *interfaces.ExportTask) error
    FindByID(ctx context.Context, id string) (*interfaces.ExportTask, error)
    FindByProjectID(ctx context.Context, projectID string, page, pageSize int) ([]*interfaces.ExportTask, int64, error)
    Update(ctx context.Context, task *interfaces.ExportTask) error
    Delete(ctx context.Context, id string) error
}
```

#### 2. FileStorage
文件存储接口，用于上传和下载导出文件：
```go
type FileStorage interface {
    Upload(ctx context.Context, filename string, content io.Reader, mimeType string) (string, error)
    Download(ctx context.Context, url string) (io.ReadCloser, error)
    Delete(ctx context.Context, url string) error
    GetSignedURL(ctx context.Context, url string, expiration time.Duration) (string, error)
}
```

建议实现：
- 本地文件系统存储（开发环境）
- 阿里云OSS / AWS S3（生产环境）

### 发布管理所需接口

#### 1. PublicationRepository
发布记录仓储接口：
```go
type PublicationRepository interface {
    Create(ctx context.Context, record *interfaces.PublicationRecord) error
    FindByID(ctx context.Context, id string) (*interfaces.PublicationRecord, error)
    FindByProjectID(ctx context.Context, projectID string, page, pageSize int) ([]*interfaces.PublicationRecord, int64, error)
    FindByResourceID(ctx context.Context, resourceID string) (*interfaces.PublicationRecord, error)
    Update(ctx context.Context, record *interfaces.PublicationRecord) error
    Delete(ctx context.Context, id string) error
    FindPublishedByProjectID(ctx context.Context, projectID string) (*interfaces.PublicationRecord, error)
}
```

#### 2. BookstoreClient
书城客户端接口，用于与外部书城平台交互：
```go
type BookstoreClient interface {
    PublishProject(ctx context.Context, req *BookstorePublishProjectRequest) (*BookstorePublishResponse, error)
    UnpublishProject(ctx context.Context, projectID, bookstoreID string) error
    PublishChapter(ctx context.Context, req *BookstorePublishChapterRequest) (*BookstorePublishResponse, error)
    UnpublishChapter(ctx context.Context, chapterID, bookstoreID string) error
    UpdateChapter(ctx context.Context, req *BookstoreUpdateChapterRequest) error
    GetStatistics(ctx context.Context, projectID, bookstoreID string) (*interfaces.PublicationStatistics, error)
}
```

建议实现：
- 模拟书城客户端（开发和测试环境）
- 真实书城平台API客户端（生产环境）

## 配置步骤

### 步骤1: 实现依赖接口
按照上述说明实现 ExportTaskRepository、FileStorage、PublicationRepository 和 BookstoreClient 接口。

### 步骤2: 注册Repository工厂方法
在 `repository/factory/repository_factory.go` 或类似文件中添加：
```go
func (f *RepositoryFactory) CreateExportTaskRepository() writer.ExportTaskRepository {
    return writer.NewMongoExportTaskRepository(f.mongoDB)
}

func (f *RepositoryFactory) CreatePublicationRepository() writer.PublicationRepository {
    return writer.NewMongoPublicationRepository(f.mongoDB)
}
```

### 步骤3: 注册服务容器中的服务
在服务容器中添加：
```go
func (c *ServiceContainer) GetFileStorage() storage.FileStorage {
    return c.fileStorage
}

func (c *ServiceContainer) GetBookstoreClient() bookstore.BookstoreClient {
    return c.bookstoreClient
}
```

### 步骤4: 取消注释路由初始化代码
在 `router/writer/writer_router_init.go` 中，取消以下注释：
```go
// 导出服务
exportTaskRepo := repositoryFactory.CreateExportTaskRepository()
fileStorage := serviceContainer.GetFileStorage()
exportSvc := writerService.NewExportService(documentRepo, documentContentRepo, projectRepo, exportTaskRepo, fileStorage)

// 发布服务
publicationRepo := repositoryFactory.CreatePublicationRepository()
bookstoreClient := serviceContainer.GetBookstoreClient()
publishSvc := writerService.NewPublishService(projectRepo, documentRepo, publicationRepo, bookstoreClient, eventBus)
```

并将 `var exportSvc` 和 `var publishSvc` 替换为实际创建的服务实例。

## 测试建议

### 1. 导出功能测试
- 创建导出任务并检查状态
- 下载已完成的导出文件
- 测试不同格式的导出（TXT/MD/DOCX）
- 测试取消导出任务

### 2. 发布功能测试
- 发布项目到模拟书城
- 取消项目发布
- 发布单个章节
- 批量发布章节
- 更新章节发布状态
- 获取发布状态和统计信息

### 3. 错误处理测试
- 测试权限验证
- 测试资源不存在的情况
- 测试参数验证
- 测试并发场景

## 注意事项

1. **异步处理**: 导出和发布操作都是异步执行的，需要注意任务状态的轮询和通知机制。

2. **文件存储**: 导出文件有24小时的过期时间，需要定期清理过期文件。

3. **权限控制**: 所有API都需要JWT认证，并验证用户对资源的访问权限。

4. **书城集成**: BookstoreClient 的实现需要根据实际书城平台的API文档进行调整。

5. **错误处理**: 所有API都使用统一的响应格式 `shared.APIResponse`，错误信息会被正确封装。

6. **定时任务**: 定时发布功能需要额外的定时任务调度器支持。

## 后续优化建议

1. **WebSocket支持**: 添加WebSocket支持，实时推送导出和发布任务进度。

2. **任务队列**: 使用消息队列（如RabbitMQ、Redis）处理异步任务。

3. **缓存优化**: 缓存发布状态和统计信息，减少对书城API的调用。

4. **监控告警**: 添加任务失败监控和告警机制。

5. **批量操作优化**: 优化批量发布的并发处理，提高性能。

## 技术栈

- **Web框架**: Gin
- **数据库**: MongoDB
- **文件存储**: 本地/OSS
- **外部集成**: 书城API
- **认证**: JWT

## 相关文档

- [Gin框架文档](https://gin-gonic.com/docs/)
- [MongoDB Go驱动](https://www.mongodb.com/docs/drivers/go/)
- [项目主README](../README.md)
