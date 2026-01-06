# AI写作辅助API文档

## 概述

本文档介绍了青羽写作平台新增的AI辅助功能API，包括内容总结、文本校对和敏感词检测功能。

## 创建的文件列表

### 1. 服务接口层
- **D:\Github\青羽\Qingyu_backend\service\interfaces\ai\writing_assistant_service.go**
  - 定义了WritingAssistantService接口
  - 包含总结、校对、敏感词检测的所有方法签名

### 2. 数据传输对象(DTO)
- **D:\Github\青羽\Qingyu_backend\service\ai\dto\writing_assistant_dto.go**
  - SummarizeRequest/Response - 内容总结
  - ChapterSummaryRequest/Response - 章节总结
  - ProofreadRequest/Response - 文本校对
  - SensitiveWordsCheckRequest/Response - 敏感词检测
  - 相关辅助结构体

### 3. 服务实现层
- **D:\Github\青羽\Qingyu_backend\service\ai\summarize_service.go**
  - SummarizeService 实现
  - 支持文档总结和章节总结
  - 自动提取关键点

- **D:\Github\青羽\Qingyu_backend\service\ai\proofread_service.go**
  - ProofreadService 实现
  - 检查拼写、语法、标点错误
  - 生成整体评分和统计信息

- **D:\Github\青羽\Qingyu_backend\service\ai\sensitive_words_service.go**
  - SensitiveWordsService 实现
  - 内置敏感词库管理
  - 支持自定义敏感词
  - AI语义分析（可选）

### 4. API处理层
- **D:\Github\青羽\Qingyu_backend\api\v1\ai\writing_assistant_api.go**
  - WritingAssistantApi 实现
  - 所有API端点的处理函数
  - 包含完整的Swagger注释

### 5. 路由配置
- **D:\Github\青羽\Qingyu_backend\router\ai\ai_router.go** (已更新)
  - 新增写作辅助路由
  - 新增内容审核路由

### 6. 服务核心
- **D:\Github\青羽\Qingyu_backend\service\ai\ai_service.go** (已更新)
  - 新增GetAdapterManager()方法

---

## API接口详细说明

### 一、内容总结API

#### 1.1 总结文档内容

**接口路径**: `POST /api/v1/ai/writing/summarize`

**功能**: 使用AI总结文档内容，支持简短摘要、详细摘要、关键点提取等多种模式

**请求参数**:
```json
{
  "content": "要总结的文档内容...",
  "projectId": "项目ID（可选）",
  "chapterId": "章节ID（可选）",
  "maxLength": 1000,           // 摘要最大长度（可选）
  "summaryType": "detailed",   // 总结类型: brief, detailed, keypoints
  "includeQuotes": true        // 是否包含关键引用（可选）
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "总结成功",
  "data": {
    "summary": "这是AI生成的摘要内容...",
    "keyPoints": [
      "关键点1",
      "关键点2",
      "关键点3"
    ],
    "originalLength": 5000,
    "summaryLength": 200,
    "compressionRate": 0.04,
    "tokensUsed": 850,
    "model": "gpt-3.5-turbo",
    "processedAt": "2026-01-03T10:30:00Z"
  },
  "timestamp": 1704274200
}
```

#### 1.2 总结章节内容

**接口路径**: `POST /api/v1/ai/writing/summarize-chapter`

**功能**: 自动提取章节要点、情节大纲、涉及角色等详细信息

**请求参数**:
```json
{
  "chapterId": "章节ID",
  "projectId": "项目ID",
  "outlineLevel": 3  // 大纲级别（可选）
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "章节总结成功",
  "data": {
    "chapterId": "ch_12345",
    "chapterTitle": "第三章：命运的转折",
    "summary": "本章描述了主人公在关键时刻做出的抉择...",
    "keyPoints": [
      "主人公发现了重要线索",
      "与反派角色初次交锋",
      "埋下了后续伏笔"
    ],
    "plotOutline": [
      {
        "level": 1,
        "title": "开场",
        "description": "平静的日常被打破"
      }
    ],
    "characters": [
      {
        "name": "张三",
        "role": "主角",
        "appearance": "全章主要角色"
      }
    ],
    "tokensUsed": 1200,
    "model": "gpt-3.5-turbo",
    "processedAt": "2026-01-03T10:30:00Z"
  }
}
```

---

### 二、文本校对API

#### 2.1 文本校对

**接口路径**: `POST /api/v1/ai/writing/proofread`

**功能**: 检查拼写、语法、标点错误，返回修改建议列表和整体评分

**请求参数**:
```json
{
  "content": "要校对的文本内容...",
  "projectId": "项目ID（可选）",
  "chapterId": "章节ID（可选）",
  "checkTypes": ["spelling", "grammar", "punctuation"],  // 检查类型
  "language": "zh-CN",    // 语言（自动检测，可选）
  "suggestions": true     // 是否提供修改建议
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "校对完成",
  "data": {
    "originalContent": "原文内容...",
    "issues": [
      {
        "id": "issue_001",
        "type": "grammar",
        "severity": "error",
        "message": "语法错误：主谓不一致",
        "position": {
          "line": 5,
          "column": 10,
          "start": 120,
          "end": 125,
          "length": 5
        },
        "originalText": "他们",
        "suggestions": ["它们"],
        "category": "语法",
        "rule": "subject_verb_agreement"
      },
      {
        "id": "issue_002",
        "type": "punctuation",
        "severity": "warning",
        "message": "标点符号建议：句号后应有空格",
        "position": {
          "line": 8,
          "column": 15,
          "start": 200,
          "end": 201,
          "length": 1
        },
        "originalText": "。",
        "suggestions": ["。 "],
        "category": "标点"
      }
    ],
    "score": 85.5,
    "statistics": {
      "totalIssues": 2,
      "errorCount": 1,
      "warningCount": 1,
      "suggestionCount": 0,
      "issuesByType": {
        "grammar": 1,
        "punctuation": 1
      },
      "wordCount": 150,
      "characterCount": 500
    },
    "tokensUsed": 600,
    "model": "gpt-3.5-turbo",
    "processedAt": "2026-01-03T10:30:00Z"
  }
}
```

#### 2.2 获取校对建议详情

**接口路径**: `GET /api/v1/ai/writing/suggestions/:id`

**功能**: 根据建议ID获取详细的修改建议和说明

**请求参数**:
- URL参数: `id` - 建议ID

**响应示例**:
```json
{
  "code": 200,
  "message": "获取建议成功",
  "data": {
    "issueId": "issue_001",
    "type": "grammar",
    "message": "建议修改语法错误",
    "position": {
      "line": 1,
      "column": 10,
      "start": 10,
      "end": 20,
      "length": 10
    },
    "originalText": "原文示例",
    "suggestions": [
      {
        "text": "建议文本",
        "confidence": 0.95,
        "reason": "语法更通顺"
      }
    ],
    "explanation": "这是一个语法错误的示例说明",
    "examples": [
      "正确示例1",
      "正确示例2"
    ]
  }
}
```

---

### 三、敏感词检测API

#### 3.1 检测敏感词

**接口路径**: `POST /api/v1/ai/audit/sensitive-words`

**功能**: 检测文本中的敏感词，返回敏感词列表、位置和修改建议

**请求参数**:
```json
{
  "content": "要检测的内容...",
  "projectId": "项目ID（可选）",
  "chapterId": "章节ID（可选）",
  "customWords": ["自定义敏感词1", "自定义敏感词2"],  // 自定义敏感词（可选）
  "category": "all"  // 检测分类: all, political, violence, adult, other
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "检测完成",
  "data": {
    "checkId": "check_abc123",
    "isSafe": false,
    "totalMatches": 2,
    "sensitiveWords": [
      {
        "id": "sw_001",
        "word": "敏感词示例",
        "category": "political",
        "level": "high",
        "position": {
          "line": 10,
          "column": 5,
          "start": 250,
          "end": 256,
          "length": 6
        },
        "context": "...前文50字符...敏感词示例...后文50字符...",
        "suggestion": "建议修改或删除敏感词「敏感词示例」",
        "reason": "该词汇属于政治敏感类别"
      }
    ],
    "summary": {
      "byCategory": {
        "political": 1,
        "violence": 1
      },
      "byLevel": {
        "high": 1,
        "medium": 1
      },
      "highRiskCount": 1,
      "mediumRiskCount": 1,
      "lowRiskCount": 0
    },
    "tokensUsed": 0,
    "processedAt": "2026-01-03T10:30:00Z"
  }
}
```

#### 3.2 获取敏感词检测结果

**接口路径**: `GET /api/v1/ai/audit/sensitive-words/:id`

**功能**: 根据检测ID获取详细的敏感词检测结果

**请求参数**:
- URL参数: `id` - 检测ID

**响应示例**:
```json
{
  "code": 200,
  "message": "获取检测结果成功",
  "data": {
    "checkId": "check_abc123",
    "content": "检测的内容...",
    "isSafe": false,
    "matches": [
      {
        "id": "sw_001",
        "word": "敏感词",
        "category": "political",
        "level": "high",
        "position": {...},
        "context": "...",
        "suggestion": "..."
      }
    ],
    "customWords": ["自定义词1"],
    "summary": {...},
    "createdAt": "2026-01-03T10:30:00Z",
    "expiresAt": "2026-02-02T10:30:00Z"
  }
}
```

---

## 技术实现要点

### 1. 统一的错误处理

所有API都使用统一的响应格式：

```go
// 成功响应
shared.Success(c, http.StatusOK, "操作成功", data)

// 错误响应
shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
```

### 2. 配额管理

所有AI功能都集成了配额检查中间件：

```go
writingGroup.Use(middleware.QuotaCheckMiddleware(quotaService))
auditGroup.Use(middleware.QuotaCheckMiddleware(quotaService))
```

### 3. 请求追踪

每个请求都生成唯一的requestID用于追踪：

```go
requestID := uuid.New().String()
c.Set("requestID", requestID)
c.Set("aiService", "service_name")
c.Set("tokensUsed", tokensUsed)
c.Set("aiModel", model)
```

### 4. AI适配器复用

所有服务都复用现有的AI适配器管理器：

```go
summarizeService := ai.NewSummarizeService(aiService.GetAdapterManager())
proofreadService := ai.NewProofreadService(aiService.GetAdapterManager())
sensitiveWordsService := ai.NewSensitiveWordsService(aiService.GetAdapterManager())
```

### 5. 敏感词库设计

- **内置词库**: 政治、暴力、成人内容分类
- **自定义词库**: 支持用户添加自定义敏感词
- **风险级别**: high/medium/low三级分类
- **位置精确定位**: 行号、列号、字符位置
- **上下文提取**: 自动提取敏感词前后50字符

### 6. 文本校对特性

- **多类型检查**: 拼写、语法、标点、风格
- **严重程度分级**: error/warning/suggestion
- **整体评分**: 0-100分，根据错误数量和严重程度计算
- **统计信息**: 总问题数、分类统计、词数字数统计

---

## 测试建议

### 1. 单元测试

建议为每个服务编写单元测试：

```go
// 测试内容总结
func TestSummarizeService_SummarizeContent(t *testing.T) {
    // 测试代码
}

// 测试文本校对
func TestProofreadService_ProofreadContent(t *testing.T) {
    // 测试代码
}

// 测试敏感词检测
func TestSensitiveWordsService_CheckSensitiveWords(t *testing.T) {
    // 测试代码
}
```

### 2. 集成测试

使用Postman或类似工具进行API测试：

```bash
# 测试总结API
curl -X POST http://localhost:8080/api/v1/ai/writing/summarize \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"content": "测试内容..."}'

# 测试校对API
curl -X POST http://localhost:8080/api/v1/ai/writing/proofread \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"content": "测试内容...", "checkTypes": ["grammar", "spelling"]}'

# 测试敏感词检测API
curl -X POST http://localhost:8080/api/v1/ai/audit/sensitive-words \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"content": "测试内容..."}'
```

---

## 后续优化建议

### 1. 缓存机制

对于重复的内容总结和校对，可以添加缓存机制：

```go
// 使用Redis缓存结果
cacheKey := fmt.Sprintf("summarize:%s", hash(content))
if cached, found := cache.Get(cacheKey); found {
    return cached.(*dto.SummarizeResponse), nil
}
```

### 2. 批量处理

支持批量文档处理：

```go
// 批量总结
BatchSummarize(ctx context.Context, reqs []*dto.SummarizeRequest) ([]*dto.SummarizeResponse, error)

// 批量校对
BatchProofread(ctx context.Context, reqs []*dto.ProofreadRequest) ([]*dto.ProofreadResponse, error)
```

### 3. 结果持久化

将检测和校对结果保存到数据库，便于历史查询：

```go
// 定义Repository接口
type ProofreadResultRepository interface {
    Save(ctx context.Context, result *dto.ProofreadResponse) error
    GetByID(ctx context.Context, id string) (*dto.ProofreadResponse, error)
    ListByUser(ctx context.Context, userID string, page, pageSize int) ([]*dto.ProofreadResponse, error)
}
```

### 4. 实时流式响应

为长文本校对和敏感词检测添加流式响应：

```go
// 流式校对
func (api *WritingAssistantApi) ProofreadContentStream(c *gin.Context) {
    // SSE流式实现
}
```

### 5. 多语言支持

扩展服务以支持多语言文本处理：

```go
// 自动检测语言
language := detectLanguage(req.Content)

// 根据语言选择对应的AI模型和规则
model := selectModelByLanguage(language)
rules := loadRulesByLanguage(language)
```

---

## 注意事项

1. **API认证**: 所有API都需要JWT认证
2. **配额限制**: AI功能消耗配额，需要确保用户有足够配额
3. **内容长度**: 建议限制单次处理的文本长度（如5000字）
4. **并发控制**: 高并发时需要考虑限流和队列机制
5. **敏感词库**: 敏感词库应该定期更新和维护
6. **错误处理**: AI服务可能失败，需要优雅降级处理

---

## Swagger文档集成

所有API都已添加Swagger注释，可以通过以下路径访问API文档：

```
http://localhost:8080/swagger/index.html
```

或者使用Swagger CLI生成文档：

```bash
swag init -g cmd/server/main.go
```

---

## 总结

本次实现为青羽写作平台新增了三大AI辅助功能：

1. **内容总结** - 智能总结文档和章节内容
2. **文本校对** - 全面的语法、拼写、标点检查
3. **敏感词检测** - 多分类敏感词识别和风险评级

所有功能都：
- 复用了现有的AI服务架构
- 使用统一的API响应格式
- 完整的错误处理和日志记录
- 支持配额管理和请求追踪
- 提供完整的Swagger文档

这些API将为作者提供强大的AI辅助工具，提升创作效率和内容质量。
