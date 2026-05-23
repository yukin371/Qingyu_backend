# AI API 设计文档

## 1. 概述

### 1.1 项目背景
青羽智能写作系统是一个基于AI技术的智能写作平台，旨在为用户提供全方位的写作辅助服务。AI API模块是系统的核心组件，负责处理所有与人工智能相关的功能请求。

### 1.2 设计目标
- **智能化**: 提供高质量的AI写作辅助功能
- **上下文感知**: 基于项目上下文提供精准的AI服务
- **可扩展性**: 支持多种AI服务提供商和模型
- **高性能**: 优化响应时间和资源利用率
- **易用性**: 提供简洁明了的API接口

### 1.3 技术架构
- **框架**: Go + Gin Web Framework
- **AI服务**: OpenAI API / Azure OpenAI / 其他兼容服务
- **数据存储**: MongoDB (上下文数据) + PostgreSQL (业务数据)
- **配置管理**: Viper + 环境变量
- **认证授权**: JWT Token

## 2. 系统架构设计

### 2.1 整体架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   前端应用      │    │   移动端应用    │    │   第三方集成    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   API Gateway   │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │   AI Router     │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │   AI Controller │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │   AI Service    │
                    └─────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Context Service │    │External API Svc │    │  Cache Service  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   MongoDB       │    │   OpenAI API    │    │     Redis       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 2.2 模块划分

#### 2.2.1 路由层 (Router Layer)
- **文件**: `router/ai/ai_router.go`
- **职责**: 定义API路由规则，将HTTP请求路由到对应的控制器方法
- **主要功能**:
  - 路由注册和管理
  - 中间件集成
  - 参数预处理

#### 2.2.2 控制器层 (Controller Layer)
- **文件**: `api/v1/ai/ai_api.go`
- **职责**: 处理HTTP请求，参数验证，调用业务服务，返回响应
- **主要功能**:
  - 请求参数绑定和验证
  - 业务逻辑调用
  - 响应格式化
  - 错误处理

#### 2.2.3 服务层 (Service Layer)
- **文件**: `service/ai/ai_service.go`
- **职责**: 核心业务逻辑处理，协调各个子服务
- **主要功能**:
  - 业务流程编排
  - 数据转换和处理
  - 服务间调用协调

#### 2.2.4 外部API服务 (External API Service)
- **文件**: `service/ai/external_api_service.go`
- **职责**: 与外部AI服务提供商的API交互
- **主要功能**:
  - API请求构建和发送
  - 响应解析和处理
  - 错误重试和恢复

#### 2.2.5 上下文服务 (Context Service)
- **文件**: `service/ai/context_service.go`
- **职责**: 管理AI上下文信息，提供智能上下文构建
- **主要功能**:
  - 项目上下文收集
  - 上下文优化和压缩
  - 历史信息管理

## 3. API 接口设计

### 3.1 接口概览

| 接口路径 | HTTP方法 | 功能描述 | 认证要求 |
|---------|----------|----------|----------|
| `/api/v1/ai/generate` | POST | 生成内容 | 是 |
| `/api/v1/ai/continue` | POST | 续写内容 | 是 |
| `/api/v1/ai/optimize` | POST | 优化文本 | 是 |
| `/api/v1/ai/outline` | POST | 生成大纲 | 是 |
| `/api/v1/ai/analyze` | POST | 分析内容 | 否 |
| `/api/v1/ai/context/:projectId` | GET | 获取项目上下文 | 是 |
| `/api/v1/ai/context/:projectId/:chapterId` | GET | 获取章节上下文 | 是 |
| `/api/v1/ai/context/feedback` | POST | 更新上下文反馈 | 是 |

### 3.2 详细接口设计

#### 3.2.1 生成内容接口

**接口信息**
- **路径**: `POST /api/v1/ai/generate`
- **功能**: 基于项目上下文和用户提示生成新内容
- **认证**: 需要JWT Token

**请求参数**
```go
type GenerateContentRequest struct {
    ProjectID string                `json:"projectId"`           // 必填，项目ID
    ChapterID string                `json:"chapterId,omitempty"` // 可选，章节ID
    Prompt    string                `json:"prompt"`              // 必填，生成提示词
    Options   *ai.GenerateOptions   `json:"options,omitempty"`   // 可选，生成选项
}

type GenerateOptions struct {
    Temperature float64 `json:"temperature,omitempty"` // 创造性程度 (0-1)
    MaxTokens   int     `json:"maxTokens,omitempty"`   // 最大生成长度
    Style       string  `json:"style,omitempty"`       // 写作风格
    Genre       string  `json:"genre,omitempty"`       // 文体类型
}
```

**响应格式**
```go
type GenerateContentResponse struct {
    Content     string    `json:"content"`     // 生成的内容
    TokensUsed  int       `json:"tokensUsed"`  // 使用的Token数量
    Model       string    `json:"model"`       // 使用的AI模型
    GeneratedAt time.Time `json:"generatedAt"` // 生成时间
}
```

**业务流程**
1. 验证请求参数（项目ID和提示词必填）
2. 构建AI上下文信息
3. 设置默认生成选项
4. 调用外部AI API生成内容
5. 返回生成结果

#### 3.2.2 续写内容接口

**接口信息**
- **路径**: `POST /api/v1/ai/continue`
- **功能**: 基于当前文本内容进行智能续写
- **认证**: 需要JWT Token

**请求参数**
```go
type ContinueWritingRequest struct {
    ProjectID      string              `json:"projectId"`                // 必填，项目ID
    ChapterID      string              `json:"chapterId"`                // 必填，章节ID
    CurrentText    string              `json:"currentText"`              // 必填，当前文本内容
    ContinueLength int                 `json:"continueLength,omitempty"` // 可选，续写长度（字数）
    Options        *ai.GenerateOptions `json:"options,omitempty"`        // 可选，生成选项
}
```

**响应格式**
```go
// 使用与生成内容相同的响应格式
type GenerateContentResponse struct {
    Content     string    `json:"content"`
    TokensUsed  int       `json:"tokensUsed"`
    Model       string    `json:"model"`
    GeneratedAt time.Time `json:"generatedAt"`
}
```

#### 3.2.3 文本优化接口

**接口信息**
- **路径**: `POST /api/v1/ai/optimize`
- **功能**: 优化文本的语法、风格、流畅度等
- **认证**: 需要JWT Token

**请求参数**
```go
type OptimizeTextRequest struct {
    ProjectID      string              `json:"projectId"`                // 必填，项目ID
    ChapterID      string              `json:"chapterId,omitempty"`      // 可选，章节ID
    OriginalText   string              `json:"originalText"`             // 必填，原始文本
    OptimizeType   string              `json:"optimizeType"`             // 可选，优化类型
    Instructions   string              `json:"instructions,omitempty"`   // 可选，具体优化指示
    Options        *ai.GenerateOptions `json:"options,omitempty"`        // 可选，生成选项
}
```

**优化类型枚举**
- `grammar`: 语法优化
- `style`: 风格优化
- `flow`: 流畅度优化
- `dialogue`: 对话优化

#### 3.2.4 大纲生成接口

**接口信息**
- **路径**: `POST /api/v1/ai/outline`
- **功能**: 基于主题和要求生成故事大纲
- **认证**: 需要JWT Token

**请求参数**
```go
type GenerateOutlineRequest struct {
    ProjectID   string              `json:"projectId"`           // 必填，项目ID
    Theme       string              `json:"theme"`               // 必填，主题
    Genre       string              `json:"genre"`               // 可选，类型
    Length      string              `json:"length"`              // 可选，长度
    KeyElements []string            `json:"keyElements"`         // 可选，关键元素数组
    Options     *ai.GenerateOptions `json:"options,omitempty"`   // 可选，生成选项
}
```

#### 3.2.5 内容分析接口

**接口信息**
- **路径**: `POST /api/v1/ai/analyze`
- **功能**: 分析文本内容，提供专业分析和建议
- **认证**: 不需要（公开接口）

**请求参数**
```go
type AnalyzeContentRequest struct {
    Content      string `json:"content"`      // 必填，要分析的内容
    AnalysisType string `json:"analysisType"` // 可选，分析类型
}
```

**分析类型枚举**
- `plot`: 情节分析
- `character`: 角色分析
- `style`: 风格分析
- `general`: 综合分析（默认）

**响应格式**
```go
type AnalyzeContentResponse struct {
    Type        string    `json:"type"`        // 分析类型
    Analysis    string    `json:"analysis"`    // 分析结果
    TokensUsed  int       `json:"tokensUsed"`  // 使用的Token数量
    Model       string    `json:"model"`       // 使用的AI模型
    AnalyzedAt  time.Time `json:"analyzedAt"`  // 分析时间
}
```

#### 3.2.6 上下文管理接口

**获取上下文信息**
- **路径**: `GET /api/v1/ai/context/:projectId[/:chapterId]`
- **功能**: 获取项目或章节的AI上下文信息
- **认证**: 需要JWT Token

**响应格式**
```go
type AIContext struct {
    ProjectID        string             `json:"projectId"`
    CurrentChapter   *ChapterInfo       `json:"currentChapter"`
    ActiveCharacters []*CharacterInfo   `json:"activeCharacters"`
    CurrentLocations []*LocationInfo    `json:"currentLocations"`
    RelevantEvents   []*TimelineEvent   `json:"relevantEvents"`
    PreviousChapters []*ChapterSummary  `json:"previousChapters"`
    NextChapters     []*ChapterOutline  `json:"nextChapters"`
    WorldSettings    *WorldSettings     `json:"worldSettings"`
    PlotThreads      []*PlotThread      `json:"plotThreads"`
    TokenCount       int                `json:"tokenCount"`
}
```

**更新上下文反馈**
- **路径**: `POST /api/v1/ai/context/feedback`
- **功能**: 更新AI上下文的用户反馈信息
- **认证**: 需要JWT Token

**请求参数**
```go
type UpdateContextFeedbackRequest struct {
    ProjectID string `json:"projectId"` // 必填，项目ID
    ChapterID string `json:"chapterId"` // 可选，章节ID
    Feedback  string `json:"feedback"`  // 必填，反馈内容
}
```

## 4. 数据模型设计

### 4.1 AI模型信息

```go
type AIModel struct {
    ID          string    `bson:"_id,omitempty" json:"id"`
    Provider    string    `bson:"provider" json:"provider"`       // 服务提供商名称
    Name        string    `bson:"name" json:"name"`               // 模型名称
    Type        ModelType `bson:"type" json:"type"`               // 模型类型
    MaxTokens   int       `bson:"max_tokens" json:"maxTokens"`    // 最大令牌数
    Enabled     bool      `bson:"enabled" json:"enabled"`         // 是否启用
    Description string    `bson:"description" json:"description"` // 模型描述
    CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
    UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}

type ModelType string
const (
    ModelTypeChat  ModelType = "chat"
    ModelTypeImage ModelType = "image"
)
```

### 4.2 上下文数据结构

#### 4.2.1 章节信息
```go
type ChapterInfo struct {
    ID           string   `json:"id"`
    Title        string   `json:"title"`
    Summary      string   `json:"summary"`
    Content      string   `json:"content"`
    CharacterIDs []string `json:"characterIds"`
    LocationIDs  []string `json:"locationIds"`
    TimelineIDs  []string `json:"timelineIds"`
    PlotThreads  []string `json:"plotThreads"`
    KeyPoints    []string `json:"keyPoints"`
    WritingHints string   `json:"writingHints"`
}
```

#### 4.2.2 角色信息
```go
type CharacterInfo struct {
    ID                string   `json:"id"`
    Name              string   `json:"name"`
    Alias             []string `json:"alias,omitempty"`
    Summary           string   `json:"summary"`
    Traits            []string `json:"traits"`
    Background        string   `json:"background"`
    PersonalityPrompt string   `json:"personalityPrompt,omitempty"`
    SpeechPattern     string   `json:"speechPattern,omitempty"`
    CurrentState      string   `json:"currentState,omitempty"`
    Relationships     []string `json:"relationships,omitempty"`
}
```

#### 4.2.3 地点信息
```go
type LocationInfo struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Climate     string `json:"climate,omitempty"`
    Culture     string `json:"culture,omitempty"`
    Geography   string `json:"geography,omitempty"`
    Atmosphere  string `json:"atmosphere,omitempty"`
}
```

#### 4.2.4 世界观设定
```go
type WorldSettings struct {
    ID          string                 `json:"id"`
    ProjectID   string                 `json:"projectId"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Rules       []string               `json:"rules,omitempty"`
    History     string                 `json:"history,omitempty"`
    Geography   string                 `json:"geography,omitempty"`
    Culture     string                 `json:"culture,omitempty"`
    Magic       string                 `json:"magic,omitempty"`
    Technology  string                 `json:"technology,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
    CreatedAt   time.Time              `json:"createdAt"`
    UpdatedAt   time.Time              `json:"updatedAt"`
}
```

#### 4.2.5 情节线索
```go
type PlotThread struct {
    ID          string   `json:"id"`
    ProjectID   string   `json:"projectId"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Status      string   `json:"status"` // active, resolved, pending, suspended
    Priority    int      `json:"priority"` // 1-10, 10为最高优先级
    ChapterIDs  []string `json:"chapterIds"`
    Characters  []string `json:"characters,omitempty"`
    StartChapter string  `json:"startChapter,omitempty"`
    EndChapter   string  `json:"endChapter,omitempty"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

## 5. 配置管理设计

### 5.1 配置结构

```go
type AIConfig struct {
    Provider     string            `mapstructure:"provider"`
    ExternalAPI  ExternalAPIConfig `mapstructure:"external_api"`
    Context      ContextConfig     `mapstructure:"context"`
    Cache        CacheConfig       `mapstructure:"cache"`
    RateLimit    RateLimitConfig   `mapstructure:"rate_limit"`
}

type ExternalAPIConfig struct {
    APIKey       string        `mapstructure:"api_key"`
    BaseURL      string        `mapstructure:"base_url"`
    DefaultModel string        `mapstructure:"default_model"`
    Timeout      time.Duration `mapstructure:"timeout"`
    MaxRetries   int           `mapstructure:"max_retries"`
}

type ContextConfig struct {
    MaxTokens      int `mapstructure:"max_tokens"`
    OverlapTokens  int `mapstructure:"overlap_tokens"`
}

type CacheConfig struct {
    Enabled bool          `mapstructure:"enabled"`
    TTL     time.Duration `mapstructure:"ttl"`
}

type RateLimitConfig struct {
    Enabled            bool `mapstructure:"enabled"`
    RequestsPerMinute  int  `mapstructure:"requests_per_minute"`
    Burst              int  `mapstructure:"burst"`
}
```

### 5.2 环境变量配置

```bash
# AI服务提供商
AI_PROVIDER=openai

# API配置
AI_API_KEY=your_api_key_here
AI_BASE_URL=https://api.openai.com/v1
AI_DEFAULT_MODEL=gpt-3.5-turbo
AI_TIMEOUT=30s
AI_MAX_RETRIES=3

# 上下文配置
AI_CONTEXT_MAX_TOKENS=4000
AI_CONTEXT_OVERLAP_TOKENS=200

# 缓存配置
AI_CACHE_ENABLED=true
AI_CACHE_TTL=1h

# 限流配置
AI_RATE_LIMIT_ENABLED=true
AI_RATE_LIMIT_REQUESTS_PER_MINUTE=60
AI_RATE_LIMIT_BURST=10
```

## 6. 错误处理设计

### 6.1 错误码定义

| 错误码 | HTTP状态码 | 错误描述 | 处理建议 |
|--------|------------|----------|----------|
| 400 | 400 | 请求参数错误 | 检查请求参数格式和必填字段 |
| 401 | 401 | 未授权访问 | 提供有效的JWT Token |
| 403 | 403 | 权限不足 | 检查用户权限设置 |
| 404 | 404 | 资源不存在 | 确认项目ID或章节ID是否正确 |
| 429 | 429 | 请求频率超限 | 降低请求频率或等待后重试 |
| 500 | 500 | 服务器内部错误 | 联系技术支持 |
| 502 | 502 | 外部AI服务不可用 | 检查AI服务商状态 |
| 503 | 503 | 服务暂时不可用 | 稍后重试 |

### 6.2 错误响应格式

```go
type ErrorResponse struct {
    Code      int    `json:"code"`
    Message   string `json:"message"`
    Timestamp int64  `json:"timestamp"`
    Details   string `json:"details,omitempty"`
}
```

### 6.3 错误处理策略

1. **参数验证错误**: 立即返回400错误，提供详细的参数错误信息
2. **认证授权错误**: 返回401/403错误，引导用户重新认证
3. **资源不存在错误**: 返回404错误，提示资源标识符
4. **限流错误**: 返回429错误，提供重试建议
5. **外部服务错误**: 实现重试机制，超过重试次数后返回502错误
6. **系统内部错误**: 记录详细日志，返回500错误

## 7. 性能优化设计

### 7.1 缓存策略

#### 7.1.1 上下文缓存
- **缓存键**: `ai:context:{projectId}:{chapterId}`
- **缓存时间**: 1小时
- **更新策略**: 项目或章节更新时清除相关缓存

#### 7.1.2 生成结果缓存
- **缓存键**: `ai:generate:{hash(prompt+context)}`
- **缓存时间**: 30分钟
- **更新策略**: LRU淘汰策略

### 7.2 限流策略

#### 7.2.1 用户级限流
- 每分钟最多60次请求
- 突发请求最多10次
- 使用Token Bucket算法

#### 7.2.2 全局限流
- 每秒最多1000次请求
- 使用滑动窗口算法

### 7.3 异步处理

#### 7.3.1 长时间任务
- 大纲生成等耗时任务使用异步处理
- 返回任务ID，客户端轮询获取结果
- 使用消息队列处理任务

#### 7.3.2 批量处理
- 支持批量内容生成
- 使用协程池并发处理
- 限制并发数量避免资源耗尽

## 8. 安全设计

### 8.1 认证授权

#### 8.1.1 JWT Token认证
- 使用RS256算法签名
- Token有效期24小时
- 支持Token刷新机制

#### 8.1.2 权限控制
- 基于项目的权限控制
- 用户只能访问自己的项目数据
- 管理员可以访问所有数据

### 8.2 数据安全

#### 8.2.1 敏感信息保护
- API密钥加密存储
- 用户内容传输加密
- 日志脱敏处理

#### 8.2.2 输入验证
- 严格的参数类型检查
- 内容长度限制
- 特殊字符过滤

### 8.3 API安全

#### 8.3.1 HTTPS强制
- 所有API请求必须使用HTTPS
- 配置HSTS头部
- 证书定期更新

#### 8.3.2 CORS配置
- 限制允许的域名
- 配置预检请求
- 限制允许的HTTP方法

## 9. 监控和日志

### 9.1 监控指标

#### 9.1.1 业务指标
- API调用次数和成功率
- 平均响应时间
- AI Token使用量
- 用户活跃度

#### 9.1.2 技术指标
- 系统CPU和内存使用率
- 数据库连接池状态
- 缓存命中率
- 错误率统计

### 9.2 日志设计

#### 9.2.1 日志级别
- **DEBUG**: 详细的调试信息
- **INFO**: 一般信息记录
- **WARN**: 警告信息
- **ERROR**: 错误信息
- **FATAL**: 致命错误

#### 9.2.2 日志格式
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "INFO",
  "service": "ai-api",
  "traceId": "abc123",
  "userId": "user123",
  "projectId": "project456",
  "action": "generate_content",
  "duration": 1500,
  "tokensUsed": 256,
  "message": "Content generated successfully"
}
```

## 10. 部署和运维

### 10.1 部署架构

#### 10.1.1 容器化部署
- 使用Docker容器化应用
- 配置健康检查
- 资源限制设置

#### 10.1.2 负载均衡
- 使用Nginx作为反向代理
- 配置负载均衡策略
- 支持蓝绿部署

### 10.2 配置管理

#### 10.2.1 环境配置
- 开发、测试、生产环境分离
- 使用配置中心管理配置
- 敏感配置加密存储

#### 10.2.2 版本管理
- 使用语义化版本号
- API版本向后兼容
- 渐进式升级策略

### 10.3 备份和恢复

#### 10.3.1 数据备份
- 定期备份数据库
- 备份配置文件
- 异地备份策略

#### 10.3.2 灾难恢复
- 制定恢复计划
- 定期演练恢复流程
- 监控恢复时间目标

## 11. 测试策略

### 11.1 单元测试

#### 11.1.1 测试覆盖率
- 代码覆盖率目标80%以上
- 关键业务逻辑100%覆盖
- 使用Go内置测试框架

#### 11.1.2 测试用例设计
- 正常流程测试
- 异常情况测试
- 边界条件测试
- 性能测试

### 11.2 集成测试

#### 11.2.1 API测试
- 使用Postman或类似工具
- 自动化API测试脚本
- 测试环境数据准备

#### 11.2.2 外部服务测试
- Mock外部AI服务
- 测试网络异常情况
- 验证重试机制

### 11.3 性能测试

#### 11.3.1 压力测试
- 模拟高并发请求
- 测试系统极限性能
- 识别性能瓶颈

#### 11.3.2 稳定性测试
- 长时间运行测试
- 内存泄漏检测
- 资源使用监控

## 12. 未来扩展

### 12.1 功能扩展

#### 12.1.1 多模态支持
- 图像生成功能
- 语音合成功能
- 视频内容生成

#### 12.1.2 高级AI功能
- 智能推荐系统
- 个性化写作助手
- 协作写作功能

### 12.2 技术扩展

#### 12.2.1 微服务架构
- 服务拆分策略
- 服务间通信
- 分布式事务处理

#### 12.2.2 云原生支持
- Kubernetes部署
- 服务网格集成
- 云服务集成

## 13. 总结

本设计文档基于青羽智能写作系统的实际代码实现，详细描述了AI API模块的架构设计、接口定义、数据模型、配置管理等各个方面。该设计具有以下特点：

1. **模块化设计**: 清晰的分层架构，职责分离明确
2. **可扩展性**: 支持多种AI服务提供商，易于扩展新功能
3. **高性能**: 通过缓存、限流、异步处理等手段优化性能
4. **安全可靠**: 完善的认证授权、错误处理、监控日志机制
5. **易于维护**: 标准化的代码结构，完善的文档和测试

该设计为青羽智能写作系统提供了坚实的技术基础，能够满足当前业务需求，并为未来的功能扩展预留了充足的空间。