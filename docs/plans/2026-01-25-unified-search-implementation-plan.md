# 统一搜索模块重构实施计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 将 Qingyu 后端分散的搜索功能统一到单一模块，使用 Provider/Engine 模式支持多搜索引擎（Elasticsearch/MongoDB/Milvus），提升可维护性、性能和扩展性。

**Architecture:** 采用 Provider/Engine 双层抽象：Provider 层处理业务逻辑和权限过滤，Engine 层抽象不同搜索引擎（ES/MongoDB/Milvus）。数据同步使用 MongoDB Change Streams → Redis Queue → Worker → ES。

**Tech Stack:** Go 1.21+, Elasticsearch 8.11, Milvus 2.3, Redis 7, MongoDB, Gin, Docker

---

## 前置检查

### Task 0: 验证中文搜索修复（Phase 0）

**目标**: 确认已有的中文搜索修复有效，补充基础可观测性

**Step 1: 测试中文搜索 API**

```bash
curl "http://localhost:8080/api/v1/bookstore/books/search?keyword=修仙"
```

**预期结果**: 返回非空结果，中文标题正确显示（不转义）

**Step 2: 测试 JSON 不转义**

检查响应中的 `title` 字段，确认中文显示正常：

```json
{
  "title": "修仙传说"
}
```

**预期结果**: 中文正常显示，不是 `\u6dc1\udcae` 这种转义

**Step 3: 添加搜索日志**

**文件:** `Qingyu_backend/service/bookstore/bookstore_service.go`

在 `SearchBooksWithFilter` 方法中添加日志：

```go
import (
    "time"
    "context"
    // ... 其他 imports
)

func (s *BookstoreServiceImpl) SearchBooksWithFilter(ctx context.Context, filter *bookstore2.BookFilter) ([]*bookstore2.Book, int64, error) {
    start := time.Now()
    
    // 现有搜索逻辑...
    books, err := s.bookRepo.SearchWithFilter(ctx, filter)
    total, err := s.bookRepo.CountByFilter(ctx, filter)
    
    // 添加日志
    took := time.Since(start)
    keyword := ""
    if filter.Keyword != nil {
        keyword = *filter.Keyword
    }
    global.Logger.Infof("[search] type=books keyword=%s took=%dms hits=%d",
        keyword, took.Milliseconds(), total)
    
    return books, total, nil
}
```

**Step 4: 回归测试**

- [ ] 英文关键词搜索
- [ ] 作者搜索
- [ ] 分类/标签过滤
- [ ] 分页功能

**Step 5: 提交验收报告**

```bash
git add service/bookstore/bookstore_service.go
git commit -m "feat(search): add search logging with request_id, took_ms, hit count"
```

---

## PR-1: 解决命名冲突 + 新 search 模块骨架

**预计时间**: 1 天

### Task 1.1: 处理现有搜索模块命名冲突

**方案选择（二选一）:**

**方案 A（推荐）: 重命名为 search_legacy**

**Step 1: 重命名目录**

```bash
cd Qingyu_backend/service/shared
mv search search_legacy
```

**Step 2: 更新包名**

**文件:** `Qingyu_backend/service/shared/search_legacy/search_service.go`

修改第一行：

```go
package search_legacy
```

**Step 3: 更新引用**

**文件:** `Qingyu_backend/api/v1/writer/search_api.go`

```go
import (
    // ...
    "Qingyu_backend/service/shared/search_legacy"
)

type SearchAPI struct {
    searchService search_legacy.SearchService
}

func NewSearchAPI(searchService search_legacy.SearchService) *SearchAPI {
    return &SearchAPI{
        searchService: searchService,
    }
}
```

**文件:** `Qingyu_backend/serviceContainer` (或其他 DI 容器文件)

查找所有 `service/shared/search` 引用并替换为 `service/shared/search_legacy`

**Step 4: 验证编译**

```bash
cd Qingyu_backend
go build ./...
```

**预期结果**: 无编译错误

**Step 5: 提交**

```bash
git add service/shared/search_legacy api/v1/writer/search_api.go
git commit -m "refactor(search): rename search package to search_legacy to avoid naming conflict"
```

---

**方案 B: 直接删除旧模块（慎用）**

如果您确认旧模块未被使用，可以直接删除：

```bash
git rm -r service/shared/search
# 然后更新所有引用
```

---

### Task 1.2: 创建新 search 模块目录结构

**Step 1: 创建目录结构**

```bash
cd Qingyu_backend/service
mkdir -p search/engine
mkdir -p search/provider
mkdir -p search/cache
mkdir -p search/sync
mkdir -p search/circuit_breaker
```

**Step 2: 创建 models/search 目录**

```bash
cd Qingyu_backend/models
mkdir -p search
```

**Step 3: 创建 repository/search 目录**

```bash
cd Qingyu_backend/repository
mkdir -p search
```

**Step 4: 验证目录创建**

```bash
tree -L 2 service/search models/search repository/search
```

**预期输出**:
```
service/search
├── engine
├── provider
├── cache
├── sync
└── circuit_breaker
models/search
repository/search
```

**Step 5: 提交**

```bash
git add service/search models/search repository/search
git commit -m "feat(search): create new search module directory structure"
```

---

### Task 1.3: 定义核心接口和模型

**Step 1: 创建 Engine 接口**

**文件:** `Qingyu_backend/service/search/engine/engine.go`

```go
package engine

import "context"

// EngineType 搜索引擎类型
type EngineType string

const (
    EngineElasticsearch EngineType = "elasticsearch"
    EngineMilvus        EngineType = "milvus"
    EngineMongoDB       EngineType = "mongodb"
)

// Document 索引文档
type Document struct {
    ID     string                 `json:"id"`
    Index  string                 `json:"index"`
    Source map[string]interface{} `json:"source"`
}

// SearchOptions 搜索选项
type SearchOptions struct {
    From     int                    `json:"from"`
    Size     int                    `json:"size"`
    Sort     []SortField            `json:"sort,omitempty"`
    Filter   map[string]interface{} `json:"filter,omitempty"`
    Highlight *HighlightConfig      `json:"highlight,omitempty"`
}

// SortField 排序字段
type SortField struct {
    Field     string `json:"field"`
    Direction string `json:"direction"` // asc, desc
}

// HighlightConfig 高亮配置
type HighlightConfig struct {
    Fields     []string `json:"fields"`
    PreTags    []string `json:"pre_tags"`
    PostTags   []string `json:"post_tags"`
    FragmentSize int    `json:"fragment_size"`
}

// SearchResult 搜索结果
type SearchResult struct {
    Total int64         `json:"total"`
    Hits  []Hit         `json:"hits"`
    Aggs  map[string]interface{} `json:"aggs,omitempty"`
    Took  int64         `json:"took_ms"`
}

// Hit 搜索命中项
type Hit struct {
    ID       string                 `json:"id"`
    Score    float64                `json:"score"`
    Source   map[string]interface{} `json:"source"`
    Highlight map[string][]string   `json:"highlight,omitempty"`
}

// Engine 搜索引擎接口
type Engine interface {
    // Search 执行搜索
    Search(ctx context.Context, index string, query interface{}, opts *SearchOptions) (*SearchResult, error)
    
    // Index 批量索引文档
    Index(ctx context.Context, index string, documents []Document) error
    
    // Update 更新文档
    Update(ctx context.Context, index string, id string, document Document) error
    
    // Delete 删除文档
    Delete(ctx context.Context, index string, id string) error
    
    // CreateIndex 创建索引
    CreateIndex(ctx context.Context, index string, mapping interface{}) error
    
    // Health 健康检查
    Health(ctx context.Context) error
}
```

**Step 2: 创建 Provider 接口**

**文件:** `Qingyu_backend/service/search/provider/provider.go`

```go
package provider

import "context"

// SearchType 搜索类型
type SearchType string

const (
    SearchTypeBooks     SearchType = "books"
    SearchTypeProjects  SearchType = "projects"
    SearchTypeDocuments SearchType = "documents"
    SearchTypeUsers     SearchType = "users"
    SearchTypeVector    SearchType = "vector"
)

// SearchRequest 统一搜索请求
type SearchRequest struct {
    Type     SearchType              `json:"type"`
    Query    string                  `json:"query"`
    Filter   map[string]interface{}  `json:"filter,omitempty"`
    Sort     []SortField             `json:"sort,omitempty"`
    Page     int                     `json:"page"`
    PageSize int                     `json:"page_size"`
    Options  map[string]interface{}  `json:"options,omitempty"`
}

// SortField 排序字段
type SortField struct {
    Field     string `json:"field"`
    Direction string `json:"direction"`
}

// SearchResponse 统一搜索响应
type SearchResponse struct {
    Type        SearchType              `json:"type"`
    Total       int64                   `json:"total"`
    Page        int                     `json:"page"`
    PageSize    int                     `json:"page_size"`
    Results     []SearchItem            `json:"results"`
    Aggregations map[string]interface{} `json:"aggregations,omitempty"`
    Took        int64                   `json:"took_ms"`
}

// SearchItem 搜索结果项
type SearchItem struct {
    ID       string                 `json:"id"`
    Score    float64                `json:"score"`
    Data     map[string]interface{} `json:"data"`
    Highlight map[string][]string   `json:"highlight,omitempty"`
}

// Provider 业务搜索提供者接口
type Provider interface {
    // Search 执行搜索
    Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error)
    
    // Type 获取搜索类型
    Type() SearchType
    
    // Validate 验证搜索参数
    Validate(req *SearchRequest) error
}
```

**Step 3: 创建搜索配置模型**

**文件:** `Qingyu_backend/models/search/config.go`

```go
package search

// SearchConfig 搜索配置
type SearchConfig struct {
    // Elasticsearch 配置
    Elasticsearch *ElasticsearchConfig `yaml:"elasticsearch"`
    
    // Milvus 配置
    Milvus *MilvusConfig `yaml:"milvus"`
    
    // 搜索行为配置
    Search *SearchBehaviorConfig `yaml:"search"`
    
    // 同步配置
    Sync *SyncConfig `yaml:"sync"`
}

// ElasticsearchConfig Elasticsearch 配置
type ElasticsearchConfig struct {
    Addresses []string `yaml:"addresses"`
    Username  string   `yaml:"username"`
    Password  string   `yaml:"password"`
    Indices   *IndicesConfig `yaml:"indices"`
}

// IndicesConfig 索引配置（支持版本化）
type IndicesConfig struct {
    Books     *IndexConfig `yaml:"books"`
    Projects  *IndexConfig `yaml:"projects"`
    Documents *IndexConfig `yaml:"documents"`
    Users     *IndexConfig `yaml:"users"`
}

// IndexConfig 单个索引配置
type IndexConfig struct {
    Alias     string `yaml:"alias"`      // 如：qingyu-books
    Version   string `yaml:"version"`    // 如：v1
    RealIndex string `yaml:"real_index"` // 如：qingyu-books-v1
}

// MilvusConfig Milvus 配置
type MilvusConfig struct {
    Address string `yaml:"address"`
    Port    int    `yaml:"port"`
}

// SearchBehaviorConfig 搜索行为配置
type SearchBehaviorConfig struct {
    // 缓存配置
    Cache *CacheConfig `yaml:"cache"`
    
    // 限流配置
    RateLimit *RateLimitConfig `yaml:"rate_limit"`
    
    // 优化器配置
    Optimizer *OptimizerConfig `yaml:"optimizer"`
    
    // 引擎选择配置
    Engines map[string]string `yaml:"engines"` // books: elasticsearch, documents: mongodb
    
    // 熔断配置
    Failover *FailoverConfig `yaml:"failover"`
    
    // Provider 配置
    Providers map[string]*ProviderConfig `yaml:"providers"`
}

// ProviderConfig Provider 配置
type ProviderConfig struct {
    // 允许搜索的状态列表（重要！配置化，避免硬编码）
    AllowedStatuses []string `yaml:"allowed_statuses"`
    // 允许的隐私设置
    AllowedPrivacy []bool `yaml:"allowed_privacy"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
    Enabled    bool          `yaml:"enabled"`
    DefaultTTL int64         `yaml:"default_ttl"` // 秒
    HotTTL     int64         `yaml:"hot_ttl"`     // 秒
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
    Enabled           bool `yaml:"enabled"`
    RequestsPerMinute int  `yaml:"requests_per_minute"`
}

// OptimizerConfig 查询优化器配置
type OptimizerConfig struct {
    MaxResults      int  `yaml:"max_results"`
    MaxPageSize     int  `yaml:"max_page_size"`
    MinQueryLength  int  `yaml:"min_query_length"`
    EnableFuzziness bool `yaml:"enable_fuzziness"`
}

// FailoverConfig 熔断配置
type FailoverConfig struct {
    Enabled         bool          `yaml:"enabled"`
    FailureThreshold int64        `yaml:"failure_threshold"`
    FailureWindow   int64         `yaml:"failure_window"`   // 秒
    HalfOpenMaxCalls int          `yaml:"half_open_max_calls"`
}

// SyncConfig 同步配置
type SyncConfig struct {
    ChangeStream *ChangeStreamConfig `yaml:"change_stream"`
    Worker       *WorkerConfig       `yaml:"worker"`
}

// ChangeStreamConfig Change Stream 配置
type ChangeStreamConfig struct {
    Enabled bool `yaml:"enabled"`
}

// WorkerConfig Worker 配置
type WorkerConfig struct {
    NumWorkers int  `yaml:"num_workers"`
    BatchSize  int  `yaml:"batch_size"`
}
```

**Step 4: 创建错误定义**

**文件:** `Qingyu_backend/models/search/errors.go`

```go
package search

import "fmt"

// SearchError 搜索错误
type SearchError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Err     error  `json:"-"`
}

func (e *SearchError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *SearchError) Unwrap() error {
    return e.Err
}

// 错误代码常量
const (
    ErrCodeInvalidRequest     = "INVALID_REQUEST"
    ErrCodeUnsupportedType    = "UNSUPPORTED_SEARCH_TYPE"
    ErrCodeEngineFailure      = "ENGINE_FAILURE"
    ErrCodeIndexNotFound      = "INDEX_NOT_FOUND"
    ErrCodeDocumentNotFound   = "DOCUMENT_NOT_FOUND"
    ErrCodeUnauthorized       = "UNAUTHORIZED"
    ErrCodeQueryParseError    = "QUERY_PARSE_ERROR"
    ErrCodeRateLimitExceeded  = "RATE_LIMIT_EXCEEDED"
    ErrCodeCircuitBreakerOpen = "CIRCUIT_BREAKER_OPEN"
)

// 预定义错误
var (
    ErrInvalidRequest    = &SearchError{Code: ErrCodeInvalidRequest, Message: "Invalid search request"}
    ErrUnauthorized      = &SearchError{Code: ErrCodeUnauthorized, Message: "Authentication required"}
    ErrIndexNotFound     = &SearchError{Code: ErrCodeIndexNotFound, Message: "Search index not found"}
    ErrCircuitBreakerOpen = &SearchError{Code: ErrCodeCircuitBreakerOpen, Message: "Circuit breaker is open, using fallback"}
)

// NewSearchError 创建搜索错误
func NewSearchError(code, message string, err error) *SearchError {
    return &SearchError{
        Code:    code,
        Message: message,
        Err:     err,
    }
}
```

**Step 5: 提交**

```bash
git add service/search/engine/engine.go service/search/provider/provider.go models/search/config.go models/search/errors.go
git commit -m "feat(search): define core interfaces - Engine, Provider, and Config models"
```

---

### Task 1.4: 创建 SearchService 主服务

**Step 1: 创建 SearchService**

**文件:** `Qingyu_backend/service/search/search.go`

```go
package search

import (
    "context"
    "fmt"
    "sync"
    "time"

    "Qingyu_backend/global"
    "Qingyu_backend/models/search"
    "Qingyu_backend/service/search/provider"
)

// SearchService 统一搜索服务
type SearchService struct {
    providers map[provider.SearchType]provider.Provider
    config    *search.SearchConfig
    mu        sync.RWMutex
}

// NewSearchService 创建搜索服务
func NewSearchService(config *search.SearchConfig) *SearchService {
    return &SearchService{
        providers: make(map[provider.SearchType]provider.Provider),
        config:    config,
    }
}

// RegisterProvider 注册 Provider
func (s *SearchService) RegisterProvider(p provider.Provider) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.providers[p.Type()] = p
    global.Logger.Infof("Search provider registered: %s", p.Type())
}

// Search 统一搜索入口
func (s *SearchService) Search(ctx context.Context, req *provider.SearchRequest) (*provider.SearchResponse, error) {
    start := time.Now()
    
    // 1. 获取对应的 Provider
    s.mu.RLock()
    p, ok := s.providers[req.Type]
    s.mu.RUnlock()
    
    if !ok {
        return nil, search.NewSearchError(
            search.ErrCodeUnsupportedType,
            fmt.Sprintf("Unsupported search type: %s", req.Type),
            nil,
        )
    }
    
    // 2. 参数校验
    if err := p.Validate(req); err != nil {
        return nil, search.NewSearchError(
            search.ErrCodeInvalidRequest,
            "Invalid search request",
            err,
        )
    }
    
    // 3. 执行搜索
    resp, err := p.Search(ctx, req)
    if err != nil {
        global.Logger.Errorf("Search failed: type=%s, query=%s, error=%v", req.Type, req.Query, err)
        return nil, err
    }
    
    // 4. 记录耗时
    resp.Took = time.Since(start).Milliseconds()
    
    global.Logger.Infof("[search] type=%s query=%s page=%d size=%d took=%dms hits=%d",
        req.Type, req.Query, req.Page, req.PageSize, resp.Took, resp.Total)
    
    return resp, nil
}

// SearchBatch 批量搜索
func (s *SearchService) SearchBatch(ctx context.Context, reqs []*provider.SearchRequest) ([]*provider.SearchResponse, error) {
    var wg sync.WaitGroup
    results := make([]*provider.SearchResponse, len(reqs))
    errs := make([]error, len(reqs))
    
    for i, req := range reqs {
        wg.Add(1)
        go func(idx int, r *provider.SearchRequest) {
            defer wg.Done()
            result, err := s.Search(ctx, r)
            results[idx] = result
            errs[idx] = err
        }(i, req)
    }
    
    wg.Wait()
    
    // 检查是否有错误
    for _, err := range errs {
        if err != nil {
            return results, err
        }
    }
    
    return results, nil
}

// Health 健康检查
func (s *SearchService) Health(ctx context.Context) error {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if len(s.providers) == 0 {
        return fmt.Errorf("no search providers registered")
    }
    
    return nil
}

// GetServiceName 获取服务名称
func (s *SearchService) GetServiceName() string {
    return "SearchService"
}

// GetVersion 获取服务版本
func (s *SearchService) GetVersion() string {
    return "2.0.0"
}
```

**Step 2: 创建配置加载器**

**文件:** `Qingyu_backend/service/search/config.go`

```go
package search

import (
    "fmt"
    "os"

    "Qingyu_backend/models/search"
    "Qingyu_backend/pkg/yaml"
)

// LoadSearchConfig 从配置文件加载搜索配置
func LoadSearchConfig(configPath string) (*search.SearchConfig, error) {
    config := &search.SearchConfig{
        // 设置默认值
        Elasticsearch: &search.ElasticsearchConfig{
            Addresses: []string{"http://localhost:9200"},
            Indices: &search.IndicesConfig{
                Books: &search.IndexConfig{
                    Alias:     "qingyu-books",
                    Version:   "v1",
                    RealIndex: "qingyu-books-v1",
                },
                Projects: &search.IndexConfig{
                    Alias:     "qingyu-projects",
                    Version:   "v1",
                    RealIndex: "qingyu-projects-v1",
                },
                Documents: &search.IndexConfig{
                    Alias:     "qingyu-documents",
                    Version:   "v1",
                    RealIndex: "qingyu-documents-v1",
                },
                Users: &search.IndexConfig{
                    Alias:     "qingyu-users",
                    Version:   "v1",
                    RealIndex: "qingyu-users-v1",
                },
            },
        },
        Search: &search.SearchBehaviorConfig{
            Cache: &search.CacheConfig{
                Enabled:    true,
                DefaultTTL: 300,  // 5 分钟
                HotTTL:     600,  // 10 分钟
            },
            RateLimit: &search.RateLimitConfig{
                Enabled:           true,
                RequestsPerMinute: 60,
            },
            Optimizer: &search.OptimizerConfig{
                MaxResults:      10000,
                MaxPageSize:     100,
                MinQueryLength:  2,
                EnableFuzziness: true,
            },
            Engines: map[string]string{
                "books":     "mongodb", // 初期使用 MongoDB，后续切换到 ES
                "projects":  "mongodb",
                "documents": "mongodb",
                "users":     "mongodb",
            },
            Failover: &search.FailoverConfig{
                Enabled:         true,
                FailureThreshold: 5,
                FailureWindow:   30,
                HalfOpenMaxCalls: 3,
            },
            Providers: map[string]*search.ProviderConfig{
                "books": {
                    AllowedStatuses: []string{"completed", "published", "serializing"}, // 重要！配置化
                    AllowedPrivacy:  []bool{false}, // 只搜索公开书籍
                },
                "projects": {
                    AllowedStatuses: []string{"draft", "in_progress", "completed"},
                },
                "documents": {
                    AllowedStatuses: []string{"draft", "completed"},
                },
            },
        },
        Sync: &search.SyncConfig{
            ChangeStream: &search.ChangeStreamConfig{
                Enabled: false, // 后续启用
            },
            Worker: &search.WorkerConfig{
                NumWorkers: 5,
                BatchSize:  100,
            },
        },
    }
    
    // 如果配置文件存在，加载配置
    if _, err := os.Stat(configPath); err == nil {
        data, err := os.ReadFile(configPath)
        if err != nil {
            return nil, fmt.Errorf("failed to read config file: %w", err)
        }
        
        if err := yaml.Unmarshal(data, config); err != nil {
            return nil, fmt.Errorf("failed to parse config file: %w", err)
        }
    }
    
    return config, nil
}
```

**Step 3: 更新主配置文件**

**文件:** `Qingyu_backend/config/config.yaml`

在文件末尾添加：

```yaml
# Elasticsearch 配置
elasticsearch:
  addresses:
    - http://localhost:9200
  username: ""
  password: ""
  indices:
    books:
      alias: qingyu-books
      version: v1
      real_index: qingyu-books-v1
    projects:
      alias: qingyu-projects
      version: v1
      real_index: qingyu-projects-v1
    documents:
      alias: qingyu-documents
      version: v1
      real_index: qingyu-documents-v1
    users:
      alias: qingyu-users
      version: v1
      real_index: qingyu-users-v1

# 搜索配置
search:
  cache:
    enabled: true
    default_ttl: 300s
    hot_ttl: 600s
  rate_limit:
    enabled: true
    requests_per_minute: 60
  optimizer:
    max_results: 10000
    max_page_size: 100
    min_query_length: 2
    enable_fuzziness: true
  engines:
    books: mongodb      # elasticsearch | mongodb | fallback
    projects: mongodb
    documents: mongodb
    users: mongodb
  failover:
    enabled: true
    failure_threshold: 5
    failure_window: 30s
    half_open_max_calls: 3
  providers:
    books:
      allowed_statuses:
        - completed
        - published
        - serializing
      allowed_privacy:
        - false
    projects:
      allowed_statuses:
        - draft
        - in_progress
        - completed
    documents:
      allowed_statuses:
        - draft
        - completed

# 数据同步配置
sync:
  change_stream:
    enabled: false
  worker:
    num_workers: 5
    batch_size: 100
```

**Step 4: 提交**

```bash
git add service/search/search.go service/search/config.go config/config.yaml
git commit -m "feat(search): implement SearchService and config loader"
```

---

## PR-2: MongoEngine + BookProvider 实现

**预计时间**: 1 天

### Task 2.1: 实现 MongoEngine

**Step 1: 创建 MongoEngine**

**文件:** `Qingyu_backend/service/search/engine/mongodb.go`

```go
package engine

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// MongoEngine MongoDB 搜索引擎
type MongoEngine struct {
    client     *mongo.Client
    database   string
}

// NewMongoEngine 创建 MongoDB 搜索引擎
func NewMongoEngine(client *mongo.Client, database string) *MongoEngine {
    return &MongoEngine{
        client:   client,
        database: database,
    }
}

// Search 执行搜索
func (e *MongoEngine) Search(ctx context.Context, collection string, query interface{}, opts *SearchOptions) (*SearchResult, error) {
    start := time.Now()
    
    // 构建查询选项
    findOpts := options.Find()
    
    // 分页
    findOpts.SetSkip(int64(opts.From))
    findOpts.SetLimit(int64(opts.Size))
    
    // 排序
    if len(opts.Sort) > 0 {
        sortDoc := bson.D{}
        for _, sf := range opts.Sort {
            direction := 1
            if sf.Direction == "desc" {
                direction = -1
            }
            sortDoc = append(sortDoc, bson.E{Key: sf.Field, Value: direction})
        }
        findOpts.SetSort(sortDoc)
    }
    
    // 执行查询
    mongoQuery, ok := query.(bson.M)
    if !ok {
        return nil, fmt.Errorf("invalid query type, expected bson.M")
    }
    
    cursor, err := e.client.Database(e.database).Collection(collection).Find(ctx, mongoQuery, findOpts)
    if err != nil {
        return nil, fmt.Errorf("mongo find failed: %w", err)
    }
    defer cursor.Close(ctx)
    
    // 获取总数
    total, err := e.client.Database(e.database).Collection(collection).CountDocuments(ctx, mongoQuery)
    if err != nil {
        return nil, fmt.Errorf("mongo count failed: %w", err)
    }
    
    // 解析结果
    var docs []bson.M
    if err := cursor.All(ctx, &docs); err != nil {
        return nil, fmt.Errorf("mongo decode failed: %w", err)
    }
    
    // 转换为 SearchResult
    hits := make([]Hit, 0, len(docs))
    for _, doc := range docs {
        hit := Hit{
            ID:     fmt.Sprintf("%v", doc["_id"]),
            Score:  1.0, // MongoDB 不提供相关性评分
            Source: make(map[string]interface{}),
        }
        
        // 转换所有字段
        for k, v := range doc {
            if k != "_id" {
                hit.Source[k] = v
            }
        }
        
        hits = append(hits, hit)
    }
    
    result := &SearchResult{
        Total: total,
        Hits:  hits,
        Took:  time.Since(start).Milliseconds(),
    }
    
    return result, nil
}

// Index 批量索引（暂不实现）
func (e *MongoEngine) Index(ctx context.Context, collection string, documents []Document) error {
    return fmt.Errorf("MongoEngine.Index not implemented - MongoDB is the source of truth")
}

// Update 更新文档（暂不实现）
func (e *MongoEngine) Update(ctx context.Context, collection string, id string, document Document) error {
    return fmt.Errorf("MongoEngine.Update not implemented - MongoDB is the source of truth")
}

// Delete 删除文档（暂不实现）
func (e *MongoEngine) Delete(ctx context.Context, collection string, id string) error {
    return fmt.Errorf("MongoEngine.Delete not implemented - MongoDB is the source of truth")
}

// CreateIndex 创建索引（暂不实现）
func (e *MongoEngine) CreateIndex(ctx context.Context, collection string, mapping interface{}) error {
    return fmt.Errorf("MongoEngine.CreateIndex not implemented")
}

// Health 健康检查
func (e *MongoEngine) Health(ctx context.Context) error {
    return e.client.Ping(ctx, nil)
}
```

**Step 2: 提交**

```bash
git add service/search/engine/mongodb.go
git commit -m "feat(search): implement MongoEngine for MongoDB-based search"
```

---

### Task 2.2: 实现 BookProvider

**Step 1: 创建 BookProvider**

**文件:** `Qingyu_backend/service/search/provider/book_provider.go`

```go
package provider

import (
    "context"
    "fmt"
    "strings"

    "Qingyu_backend/global"
    "Qingyu_backend/models/search"
    "Qingyu_backend/service/search/engine"
)

// BookProvider 书籍搜索 Provider
type BookProvider struct {
    engine *engine.MongoEngine
    config *search.ProviderConfig
}

// NewBookProvider 创建书籍搜索 Provider
func NewBookProvider(mongoEngine *engine.MongoEngine, config *search.ProviderConfig) *BookProvider {
    return &BookProvider{
        engine: mongoEngine,
        config: config,
    }
}

// Type 获取搜索类型
func (p *BookProvider) Type() SearchType {
    return SearchTypeBooks
}

// Validate 验证搜索参数
func (p *BookProvider) Validate(req *SearchRequest) error {
    if req.Query == "" {
        return fmt.Errorf("search query is required")
    }
    if len(req.Query) < 2 {
        return fmt.Errorf("query too short (min 2 characters)")
    }
    if req.Page < 1 {
        req.Page = 1
    }
    if req.PageSize < 1 || req.PageSize > 100 {
        req.PageSize = 20
    }
    return nil
}

// Search 执行搜索
func (p *BookProvider) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
    // 1. 构建查询
    mongoQuery, err := p.buildQuery(req)
    if err != nil {
        return nil, fmt.Errorf("failed to build query: %w", err)
    }
    
    // 2. 构建搜索选项
    opts := &engine.SearchOptions{
        From:     (req.Page - 1) * req.PageSize,
        Size:     req.PageSize,
        Sort:     p.buildSort(req.Sort),
        Filter:   req.Filter,
        Highlight: nil, // MongoDB 不支持高亮
    }
    
    // 3. 执行搜索
    result, err := p.engine.Search(ctx, "books", mongoQuery, opts)
    if err != nil {
        return nil, fmt.Errorf("search failed: %w", err)
    }
    
    // 4. 转换为 SearchResponse
    resp := &SearchResponse{
        Type:     SearchTypeBooks,
        Total:    result.Total,
        Page:     req.Page,
        PageSize: req.PageSize,
        Results:  make([]SearchItem, 0, len(result.Hits)),
        Took:     result.Took,
    }
    
    for _, hit := range result.Hits {
        item := SearchItem{
            ID:    hit.ID,
            Score: hit.Score,
            Data:  hit.Source,
        }
        resp.Results = append(resp.Results, item)
    }
    
    global.Logger.Infof("[BookProvider] query=%s page=%d size=%d hits=%d took=%dms",
        req.Query, req.Page, req.PageSize, resp.Total, resp.Took)
    
    return resp, nil
}

// buildQuery 构建查询
func (p *BookProvider) buildQuery(req *SearchRequest) (interface{}, error) {
    query := bson.M{}
    
    // 1. 关键词搜索（多字段）
    if req.Query != "" {
        keyword := strings.TrimSpace(req.Query)
        
        // 使用 $or 实现多字段搜索
        orConditions := []bson.M{
            {"title": bson.M{"$regex": keyword, "$options": "i"}},
            {"author": bson.M{"$regex": keyword, "$options": "i"}},
            {"introduction": bson.M{"$regex": keyword, "$options": "i"}},
        }
        
        // 标签搜索
        orConditions = append(orConditions, bson.M{"tags": bson.M{"$in": []string{keyword}}})
        
        query["$or"] = orConditions
    }
    
    // 2. 强制应用权限过滤（重要！使用配置）
    if p.config != nil && len(p.config.AllowedStatuses) > 0 {
        query["status"] = bson.M{"$in": p.config.AllowedStatuses}
    }
    if p.config != nil && len(p.config.AllowedPrivacy) > 0 {
        query["is_private"] = bson.M{"$in": p.config.AllowedPrivacy}
    }
    
    // 3. 应用额外过滤条件
    if req.Filter != nil {
        if categoryID, ok := req.Filter["category_id"]; ok {
            query["category_id"] = categoryID
        }
        if author, ok := req.Filter["author"]; ok {
            query["author"] = author
        }
        if tags, ok := req.Filter["tags"]; ok {
            query["tags"] = bson.M{"$in": tags}
        }
        if minRating, ok := req.Filter["rating_min"]; ok {
            query["rating"] = bson.M{"$gte": minRating}
        }
        if minWordCount, ok := req.Filter["word_count_min"]; ok {
            query["word_count"] = bson.M{"$gte": minWordCount}
        }
        if maxWordCount, ok := req.Filter["word_count_max"]; ok {
            if minWC, exists := query["word_count"]; exists {
                // 已有最小值，合并范围
                if minWCMap, ok := minWC.(bson.M); ok {
                    minWCMap["$lte"] = maxWordCount
                }
            } else {
                query["word_count"] = bson.M{"$lte": maxWordCount}
            }
        }
    }
    
    return query, nil
}

// buildSort 构建排序
func (p *BookProvider) buildSort(sortFields []SortField) []engine.SortField {
    if len(sortFields) == 0 {
        // 默认按更新时间降序
        return []engine.SortField{
            {Field: "updated_at", Direction: "desc"},
        }
    }
    
    engineSorts := make([]engine.SortField, 0, len(sortFields))
    for _, sf := range sortFields {
        // 映射排序字段
        field := p.mapSortField(sf.Field)
        if field != "" {
            engineSorts = append(engineSorts, engine.SortField{
                Field:     field,
                Direction: sf.Direction,
            })
        }
    }
    
    return engineSorts
}

// mapSortField 映射排序字段
func (p *BookProvider) mapSortField(field string) string {
    switch field {
    case "updateTime", "updated_at":
        return "updated_at"
    case "rating":
        return "rating"
    case "viewCount", "view_count":
        return "view_count"
    case "wordCount", "word_count":
        return "word_count"
    case "likeCount", "like_count":
        return "like_count"
    case "relevance":
        return "" // MongoDB 不支持相关性排序
    default:
        return "updated_at"
    }
}
```

**注意**: 需要在文件开头添加 bson 导入：

```go
import (
    "context"
    "fmt"
    "strings"
    
    "go.mongodb.org/mongo-driver/bson"
    
    "Qingyu_backend/global"
    "Qingyu_backend/models/search"
    "Qingyu_backend/service/search/engine"
)
```

**Step 2: 提交**

```bash
git add service/search/provider/book_provider.go
git commit -m "feat(search): implement BookProvider with configurable status filtering"
```

---

### Task 2.3: 创建统一搜索 API

**Step 1: 创建搜索 API**

**文件:** `Qingyu_backend/api/v1/search/search_api.go`

```go
package search

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "Qingyu_backend/api/v1/shared"
    "Qingyu_backend/models/search"
    "Qingyu_backend/service/search"
)

// SearchAPI 统一搜索 API
type SearchAPI struct {
    searchService *search.SearchService
}

// NewSearchAPI 创建搜索 API
func NewSearchAPI(searchService *search.SearchService) *SearchAPI {
    return &SearchAPI{
        searchService: searchService,
    }
}

// SearchBooks 搜索书籍
//
//	@Summary		搜索书籍
//	@Description	根据关键词搜索公开书籍
//	@Tags			统一搜索
//	@Accept			json
//	@Produce		json
//	@Param			q			query		string		true	"搜索关键词"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Param			category_id	query		string	false	"分类ID"
//	@Param			author		query		string	false	"作者"
//	@Param			tags		query		[]string	false	"标签"
//	@Param			rating_min	query		float64	false	"最低评分"
//	@Param			sort_by		query		string	false	"排序字段"	default(updateTime)
//	@Param			sort_order	query		string	false	"排序方向"	default(desc)
//	@Success		200			{object}	shared.APIResponse{data=provider.SearchResponse}
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/search/books [get]
func (api *SearchAPI) SearchBooks(c *gin.Context) {
    // 1. 获取搜索参数
    query := c.Query("q")
    if query == "" {
        shared.BadRequest(c, "参数错误", "搜索关键词不能为空")
        return
    }
    
    // 兼容 pageSize 参数
    pageSize := api.adaptPageSize(c)
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    
    // 2. 构建过滤条件
    filter := make(map[string]interface{})
    if categoryID := c.Query("category_id"); categoryID != "" {
        filter["category_id"] = categoryID
    }
    if author := c.Query("author"); author != "" {
        filter["author"] = author
    }
    if tags := c.QueryArray("tags"); len(tags) > 0 {
        filter["tags"] = tags
    }
    if ratingMin := c.Query("rating_min"); ratingMin != "" {
        if rating, err := strconv.ParseFloat(ratingMin, 64); err == nil {
            filter["rating_min"] = rating
        }
    }
    if wordCountMin := c.Query("word_count_min"); wordCountMin != "" {
        if wc, err := strconv.Atoi(wordCountMin); err == nil {
            filter["word_count_min"] = wc
        }
    }
    if wordCountMax := c.Query("word_count_max"); wordCountMax != "" {
        if wc, err := strconv.Atoi(wordCountMax); err == nil {
            filter["word_count_max"] = wc
        }
    }
    
    // 3. 构建排序
    sortBy := c.DefaultQuery("sort_by", "updateTime")
    sortOrder := c.DefaultQuery("sort_order", "desc")
    
    sortFields := []search.SortField{
        {Field: sortBy, Direction: sortOrder},
    }
    
    // 4. 构建搜索请求
    req := &search.SearchRequest{
        Type:     search.SearchTypeBooks,
        Query:    query,
        Filter:   filter,
        Sort:     sortFields,
        Page:     page,
        PageSize: pageSize,
    }
    
    // 5. 执行搜索
    resp, err := api.searchService.Search(c.Request.Context(), req)
    if err != nil {
        api.handleSearchError(c, err)
        return
    }
    
    // 6. 返回结果
    shared.Success(c, http.StatusOK, "搜索成功", resp)
}

// adaptPageSize 适配 pageSize 参数（兼容性）
func (api *SearchAPI) adaptPageSize(c *gin.Context) int {
    // 优先使用 page_size
    if ps := c.Query("page_size"); ps != "" {
        if size, err := strconv.Atoi(ps); err == nil && size > 0 && size <= 100 {
            return size
        }
    }
    // 兼容 pageSize
    if ps := c.Query("pageSize"); ps != "" {
        if size, err := strconv.Atoi(ps); err == nil && size > 0 && size <= 100 {
            return size
        }
    }
    // 兼容 size
    if ps := c.Query("size"); ps != "" {
        if size, err := strconv.Atoi(ps); err == nil && size > 0 && size <= 100 {
            return size
        }
    }
    // 默认值
    return 20
}

// handleSearchError 处理搜索错误
func (api *SearchAPI) handleSearchError(c *gin.Context, err error) {
    if searchErr, ok := err.(*search.SearchError); ok {
        switch searchErr.Code {
        case search.ErrCodeInvalidRequest:
            shared.BadRequest(c, "参数错误", searchErr.Message)
        case search.ErrCodeCircuitBreakerOpen:
            // 熔断打开，返回降级结果
            shared.Error(c, http.StatusServiceUnavailable, "搜索服务暂时不可用", searchErr.Message)
        default:
            shared.InternalError(c, "搜索失败", searchErr)
        }
    } else {
        shared.InternalError(c, "搜索失败", err)
    }
}

// SearchProjects 搜索项目（需要认证）
//
//	@Summary		搜索项目
//	@Description	搜索用户自己的创作项目
//	@Tags			统一搜索
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			q			query		string	true	"搜索关键词"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		401			{object}	shared.APIResponse
//	@Router			/api/v1/search/projects [get]
func (api *SearchAPI) SearchProjects(c *gin.Context) {
    // TODO: PR-5 实现
    shared.NotImplemented(c, "项目搜索功能即将上线")
}

// SearchDocuments 搜索文档（需要认证）
//
//	@Summary		搜索文档
//	@Description	搜索用户项目中的文档
//	@Tags			统一搜索
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			q			query		string	true	"搜索关键词"
//	@Param			project_id	query		string	false	"项目ID"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		401			{object}	shared.APIResponse
//	@Router			/api/v1/search/documents [get]
func (api *SearchAPI) SearchDocuments(c *gin.Context) {
    // TODO: PR-5 实现
    shared.NotImplemented(c, "文档搜索功能即将上线")
}

// SearchUsers 搜索用户
//
//	@Summary		搜索用户
//	@Description	根据用户名/昵称搜索用户
//	@Tags			统一搜索
//	@Accept			json
//	@Produce		json
//	@Param			q			query		string	true	"搜索关键词"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Router			/api/v1/search/users [get]
func (api *SearchAPI) SearchUsers(c *gin.Context) {
    // TODO: PR-5 实现
    shared.NotImplemented(c, "用户搜索功能即将上线")
}
```

**Step 2: 创建搜索路由**

**文件:** `Qingyu_backend/api/v1/search/search_router.go`

```go
package search

import (
    "github.com/gin-gonic/gin"

    "Qingyu_backend/middleware"
    "Qingyu_backend/service/search"
)

// RegisterSearchRoutes 注册搜索路由
func RegisterSearchRoutes(router *gin.RouterGroup, searchService *search.SearchService) {
    searchAPI := NewSearchAPI(searchService)
    
    // 公开路由（无需认证）
    router.GET("/books", searchAPI.SearchBooks)
    router.GET("/users", searchAPI.SearchUsers)
    
    // 需要认证的路由
    authGroup := router.Group("")
    authGroup.Use(middleware.JWTAuth())
    {
        authGroup.GET("/projects", searchAPI.SearchProjects)
        authGroup.GET("/documents", searchAPI.SearchDocuments)
    }
}
```

**Step 3: 提交**

```bash
git add api/v1/search/search_api.go api/v1/search/search_router.go
git commit -m "feat(search): create unified search API with /books endpoint"
```

---

### Task 2.4: 更新旧 Bookstore API 以使用新搜索服务

**Step 1: 修改 BookstoreAPI 以使用新搜索服务**

**文件:** `Qingyu_backend/api/v1/bookstore/bookstore_api.go`

在现有 `SearchBooks` 方法中添加 fallback 逻辑：

```go
package bookstore

import (
    bookstore2 "Qingyu_backend/models/bookstore"
    "log"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson/primitive"

    "Qingyu_backend/api/v1/shared"
    bookstoreService "Qingyu_backend/service/bookstore"
    searchService "Qingyu_backend/service/search"
    "Qingyu_backend/models/search"
    searchProvider "Qingyu_backend/service/search/provider"
)

// BookstoreAPI 书城API处理器
type BookstoreAPI struct {
    service       bookstoreService.BookstoreService
    searchService *searchService.SearchService // 新增
}

// NewBookstoreAPI 创建书城API实例
func NewBookstoreAPI(
    service bookstoreService.BookstoreService,
    searchService *searchService.SearchService, // 新增参数
) *BookstoreAPI {
    return &BookstoreAPI{
        service:       service,
        searchService: searchService,
    }
}

// SearchBooks 搜索书籍（保留旧接口，内部调用新搜索服务）
//
//	@Summary		搜索书籍
//	@Description	根据关键词搜索书籍
//	@Tags			书籍
//	@Accept			json
//	@Produce		json
//	@Param			keyword		query		string	true	"搜索关键词"
//	@Param			author		query		string	false	"作者"
//	@Param			categoryId	query		string	false	"分类ID"
//	@Param			tags		query		[]string	false	"标签"
//	@Param			status		query		string	false	"书籍状态"
//	@Param			wordCountMin	query		int		false	"最小字数"
//	@Param			wordCountMax	query		int		false	"最大字数"
//	@Param			ratingMin	query		float64	false	"最低评分"
//	@Param			sortBy		query		string	false	"排序字段"
//	@Param			sortOrder	query		string	false	"排序方向"
//	@Param			page		query		int		false	"页码"
//	@Param			pageSize	query		int		false	"每页数量"
//	@Success		200			{object}	PaginatedResponse
//	@Failure		400			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/bookstore/books/search [get]
func (api *BookstoreAPI) SearchBooks(c *gin.Context) {
    // 1. 尝试使用新的统一搜索服务
    if api.searchService != nil {
        if result, err := api.searchViaNewService(c); err == nil {
            // 成功，返回新服务的格式
            api.convertToLegacyFormat(c, result)
            return
        }
        // 新服务失败，fallback 到旧逻辑
        log.Printf("[BookstoreAPI] New search service failed, fallback to legacy: %v", err)
    }
    
    // 2. Fallback 到旧的 MongoDB 搜索
    api.searchBooksLegacy(c)
}

// searchViaNewService 通过新搜索服务搜索
func (api *BookstoreAPI) searchViaNewService(c *gin.Context) (*searchProvider.SearchResponse, error) {
    // 构建新搜索请求
    req := api.adaptSearchRequest(c)
    
    // 调用新搜索服务
    resp, err := api.searchService.Search(c.Request.Context(), req)
    if err != nil {
        return nil, err
    }
    
    return resp, nil
}

// adaptSearchRequest 将旧参数适配为新搜索请求
func (api *BookstoreAPI) adaptSearchRequest(c *gin.Context) *searchProvider.SearchRequest {
    keyword := c.Query("keyword")
    
    // 构建过滤条件
    filter := make(map[string]interface{})
    
    if author := c.Query("author"); author != "" {
        filter["author"] = author
    }
    if categoryID := c.Query("categoryId"); categoryID != "" {
        filter["category_id"] = categoryID
    }
    if tags := c.QueryArray("tags"); len(tags) > 0 {
        filter["tags"] = tags
    }
    if status := c.Query("status"); status != "" {
        filter["status"] = status
    }
    if minRating := c.Query("ratingMin"); minRating != "" {
        if rating, err := strconv.ParseFloat(minRating, 64); err == nil {
            filter["rating_min"] = rating
        }
    }
    if minWC := c.Query("wordCountMin"); minWC != "" {
        if wc, err := strconv.Atoi(minWC); err == nil {
            filter["word_count_min"] = wc
        }
    }
    if maxWC := c.Query("wordCountMax"); maxWC != "" {
        if wc, err := strconv.Atoi(maxWC); err == nil {
            filter["word_count_max"] = wc
        }
    }
    
    // 构建排序
    sortBy := c.DefaultQuery("sortBy", "updateTime")
    sortOrder := c.DefaultQuery("sortOrder", "desc")
    
    sortFields := []search.SortField{
        {Field: sortBy, Direction: sortOrder},
    }
    
    // 兼容 pageSize
    pageSize := api.adaptPageSize(c)
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    
    return &searchProvider.SearchRequest{
        Type:     searchProvider.SearchTypeBooks,
        Query:    keyword,
        Filter:   filter,
        Sort:     sortFields,
        Page:     page,
        PageSize: pageSize,
    }
}

// adaptPageSize 适配 pageSize 参数
func (api *BookstoreAPI) adaptPageSize(c *gin.Context) int {
    if ps := c.Query("pageSize"); ps != "" {
        if size, err := strconv.Atoi(ps); err == nil && size > 0 && size <= 100 {
            return size
        }
    }
    if ps := c.Query("size"); ps != "" {
        if size, err := strconv.Atoi(ps); err == nil && size > 0 && size <= 100 {
            return size
        }
    }
    return 20
}

// convertToLegacyFormat 转换为旧响应格式
func (api *BookstoreAPI) convertToLegacyFormat(c *gin.Context, resp *searchProvider.SearchResponse) {
    // 提取书籍数据
    books := make([]map[string]interface{}, 0, len(resp.Results))
    for _, item := range resp.Results {
        books = append(books, item.Data)
    }
    
    // 使用旧的 Paginated 响应格式
    shared.Paginated(c, books, resp.Total, resp.Page, resp.PageSize, "搜索成功")
}

// searchBooksLegacy 旧的搜索逻辑（fallback）
func (api *BookstoreAPI) searchBooksLegacy(c *gin.Context) {
    // 原有的搜索逻辑保留在这里
    keyword := c.Query("keyword")
    
    filter := &bookstore2.BookFilter{}
    if keyword != "" {
        filter.Keyword = &keyword
    }
    
    if author := c.Query("author"); author != "" {
        filter.Author = &author
    }
    
    if categoryID := c.Query("categoryId"); categoryID != "" {
        filter.CategoryID = &categoryID
    }
    
    if status := c.Query("status"); status != "" {
        filter.Status = (*bookstore2.BookStatus)(&status)
    }
    
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize := api.adaptPageSize(c)
    
    books, total, err := api.service.SearchBooksWithFilter(c.Request.Context(), filter)
    if err != nil {
        shared.InternalError(c, "搜索失败", err)
        return
    }
    
    bookDTOs := ToBookDTOsFromPtrSlice(books)
    shared.Paginated(c, bookDTOs, total, page, pageSize, "搜索成功")
}
```

**Step 2: 提交**

```bash
git add api/v1/bookstore/bookstore_api.go
git commit -m "feat(search): integrate new search service with fallback to legacy"
```

---

### Task 2.5: 注册路由和更新依赖注入

**Step 1: 更新主路由**

**文件:** `Qingyu_backend/router/enter.go` (或类似的路由注册文件)

在 v1 路由组中添加搜索路由：

```go
// 注册统一搜索路由
import (
    searchAPI "Qingyu_backend/api/v1/search"
    searchService "Qingyu_backend/service/search"
)

func SetupRouter() *gin.Engine {
    // ... 现有代码
    
    // 初始化搜索配置
    searchConfig, err := searchService.LoadSearchConfig("config/config.yaml")
    if err != nil {
        global.Logger.Fatalf("Failed to load search config: %v", err)
    }
    
    // 初始化搜索服务
    svc := searchService.NewSearchService(searchConfig)
    
    // 初始化 MongoEngine
    mongoEngine := engine.NewMongoEngine(mongoClient, "Qingyu_backend")
    
    // 初始化 BookProvider
    bookConfig := searchConfig.Search.Providers["books"]
    bookProvider := provider.NewBookProvider(mongoEngine, bookConfig)
    svc.RegisterProvider(bookProvider)
    
    // 注册搜索路由
    searchGroup := v1.Group("/search")
    searchAPI.RegisterSearchRoutes(searchGroup, svc)
    global.Logger.Info("Unified search routes registered at /api/v1/search")
    
    // 更新 BookstoreAPI 注入
    bookstoreAPI := bookstore.NewBookstoreAPI(bookstoreService, svc)
    
    // ... 其他路由
}
```

**Step 2: 提交**

```bash
git add router/enter.go
git commit -m "feat(search): register unified search routes and update DI"
```

---

## 验收 PR-2

### Task 2.6: 测试新搜索功能

**Step 1: 启动服务**

```bash
cd Qingyu_backend
go run main.go
```

**Step 2: 测试新搜索接口**

```bash
curl "http://localhost:8080/api/v1/search/books?q=修仙&page=1&page_size=5"
```

**预期结果**: 返回搜索结果，格式正确

**Step 3: 测试旧接口兼容性**

```bash
curl "http://localhost:8080/api/v1/bookstore/books/search?keyword=修仙&page=1&pageSize=5"
```

**预期结果**: 返回相同结果（内部调用新服务）

**Step 4: 测试参数适配**

```bash
# 测试 page_size 参数
curl "http://localhost:8080/api/v1/search/books?q=修仙&page=1&page_size=3"

# 测试 pageSize 兼容
curl "http://localhost:8080/api/v1/search/books?q=修仙&page=1&pageSize=3"

# 测试 size 兼容
curl "http://localhost:8080/api/v1/search/books?q=修仙&page=1&size=3"
```

**预期结果**: 所有参数都正常工作

**Step 5: 测试过滤条件**

```bash
# 分类过滤
curl "http://localhost:8080/api/v1/search/books?q=修仙&category_id=xxx"

# 作者过滤
curl "http://localhost:8080/api/v1/search/books?q=修仙&author=xxx"

# 评分过滤
curl "http://localhost:8080/api/v1/search/books?q=修仙&rating_min=4.0"
```

**预期结果**: 过滤条件正确应用

**Step 6: 验证配置化状态过滤**

检查日志，确认 BookProvider 使用了配置的 `allowed_statuses`：

```bash
tail -f backend.log | grep "\[BookProvider\]"
```

**预期结果**: 日志显示搜索使用了正确的状态过滤

---

## 下一步（PR-3 及之后）

完成 PR-2 后，您已经：
- ✅ 解决了命名冲突
- ✅ 搭建了新的 search 模块骨架
- ✅ 实现了 MongoEngine 和 BookProvider
- ✅ 创建了统一搜索 API
- ✅ 保持了旧接口兼容性

后续 PR（PR-3 到 PR-7）将在此基础上：
- PR-3: 添加 Elasticsearch 基础设施
- PR-4: 实现 ElasticsearchEngine 和灰度切换
- PR-5: 扩展其他 Provider
- PR-6: 实现数据同步
- PR-7: 性能优化和监控

这些 PR 的详细计划将在完成 PR-2 后根据实际情况补充喵~

---

**文档版本**: 1.0  
**创建日期**: 2026-01-25  
**预计工期**: PR-0 (0.5天) + PR-1 (1天) + PR-2 (1天) = 2.5 天  
**维护者**: Qingyu Team

---

## 执行建议

1. **使用 superpowers:subagent-driven-development** 技能来执行此计划，每个 Task 由子代理完成
2. **每个 PR 完成后进行代码审查和测试**，确保质量
3. **保持频繁提交**，每个小步骤都提交代码
4. **遇到问题及时记录**，更新计划文档

祝实施顺利喵~ 🎯
