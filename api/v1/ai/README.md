# AI API 模块结构说明

## 📁 文件结构

```
api/v1/ai/
├── writing_api.go      # AI写作API（续写、改写）
├── chat_api.go         # AI聊天API（聊天、会话管理）
├── system_api.go       # AI系统API（健康检查、提供商、模型）
├── quota_api.go        # 配额管理API
└── README.md           # 本文件
```

## 🎯 模块职责划分

### 1. WritingApi (`writing_api.go`)

**职责**: AI智能写作功能

**核心功能**:
- ✅ 智能续写（标准/流式）
- ✅ 文本改写（扩写/缩写/润色，标准/流式）

**API端点**:
```
POST /api/v1/ai/writing/continue          # 智能续写
POST /api/v1/ai/writing/continue/stream   # 智能续写（流式）
POST /api/v1/ai/writing/rewrite           # 文本改写
POST /api/v1/ai/writing/rewrite/stream    # 文本改写（流式）
```

**依赖服务**:
- `aiService.Service` - AI核心服务
- `aiService.QuotaService` - 配额管理服务

---

### 2. ChatApi (`chat_api.go`)

**职责**: AI聊天助手功能

**核心功能**:
- ✅ 对话聊天（标准/流式）
- ✅ 会话管理（列表、历史、删除）

**API端点**:
```
POST   /api/v1/ai/chat                    # 聊天
POST   /api/v1/ai/chat/stream             # 聊天（流式）
GET    /api/v1/ai/chat/sessions           # 获取会话列表
GET    /api/v1/ai/chat/sessions/:id       # 获取会话历史
DELETE /api/v1/ai/chat/sessions/:id       # 删除会话
```

**依赖服务**:
- `aiService.ChatService` - 聊天服务
- `aiService.QuotaService` - 配额管理服务

---

### 3. SystemApi (`system_api.go`)

**职责**: AI系统功能

**核心功能**:
- ✅ 健康检查
- ✅ AI提供商管理
- ✅ AI模型查询

**API端点**:
```
GET /api/v1/ai/health      # 健康检查
GET /api/v1/ai/providers   # 获取提供商列表
GET /api/v1/ai/models      # 获取模型列表
```

**依赖服务**:
- `aiService.Service` - AI核心服务

---

### 4. QuotaApi (`quota_api.go`)

**职责**: 配额管理

**核心功能**:
- ✅ 配额查询（个人/所有类型）
- ✅ 配额统计
- ✅ 事务历史
- ✅ 管理员配额管理（更新/暂停/激活）

**API端点**:
```
# 用户API
GET /api/v1/ai/quota               # 获取配额信息
GET /api/v1/ai/quota/all           # 获取所有配额
GET /api/v1/ai/quota/statistics    # 获取配额统计
GET /api/v1/ai/quota/transactions  # 获取事务历史

# 管理员API
PUT  /api/v1/admin/quota/:userId           # 更新用户配额
POST /api/v1/admin/quota/:userId/suspend   # 暂停用户配额
POST /api/v1/admin/quota/:userId/activate  # 激活用户配额
```

**依赖服务**:
- `aiService.QuotaService` - 配额管理服务

---

## 🔄 API调用流程

### 标准流程（非流式）
```
客户端请求 
  → Router 
  → JWTAuth中间件 
  → QuotaCheck中间件 
  → API Handler 
  → Service层 
  → Repository层 
  → 数据库
```

### 流式响应流程（SSE）
```
客户端请求 
  → Router 
  → JWTAuth中间件 
  → QuotaCheck中间件 
  → API Handler 
  → Service层（生成Stream Channel）
  → 逐块推送给客户端
  → 完成后异步消费配额
```

---

## 🛡️ 中间件配置

### 1. 认证中间件
所有AI接口都需要JWT认证：
```go
aiGroup.Use(middleware.JWTAuth())
```

### 2. 配额检查中间件

**标准配额检查**（预估1000 tokens）:
```go
writingGroup.Use(middleware.QuotaCheckMiddleware(quotaService))
```

**轻量级配额检查**（预估300 tokens）:
```go
chatGroup.Use(middleware.LightQuotaCheckMiddleware(quotaService))
```

**重量级配额检查**（预估3000 tokens）:
```go
// 用于长文本生成
heavyGroup.Use(middleware.HeavyQuotaCheckMiddleware(quotaService))
```

---

## 📊 请求/响应示例

### 智能续写请求
```json
POST /api/v1/ai/writing/continue/stream
Content-Type: application/json
Authorization: Bearer <token>

{
  "projectId": "project-123",
  "chapterId": "chapter-456",
  "currentText": "在一个宁静的午后，李明独自坐在咖啡馆的角落...",
  "continueLength": 500,
  "options": {
    "temperature": 0.8,
    "maxTokens": 1000
  }
}
```

### SSE流式响应
```
event: message
data: {"requestId":"req-uuid","delta":"他","content":"他","tokens":1}

event: message
data: {"requestId":"req-uuid","delta":"端起","content":"他端起","tokens":3}

event: done
data: {"requestId":"req-uuid","content":"完整内容...","tokensUsed":450,"model":"gpt-4"}
```

---

## 🔧 设计原则

### 1. 单一职责原则
每个API文件只负责一个特定的功能领域，职责清晰、边界明确。

### 2. 依赖注入
通过构造函数注入依赖服务，便于单元测试和依赖管理。

### 3. RESTful风格
- 使用标准HTTP方法（GET/POST/PUT/DELETE）
- 资源路径清晰（/ai/writing/continue、/ai/chat/sessions）
- 状态码语义明确

### 4. 流式优先
所有AI生成接口都提供流式版本（/stream），降低用户感知延迟。

### 5. 统一响应格式
使用 `shared.Success` 和 `shared.Error` 统一响应格式。

---

## 📝 开发规范

### 1. 命名规范
- API结构体：`<功能>Api`（如 `WritingApi`、`ChatApi`）
- 构造函数：`New<功能>Api`（如 `NewWritingApi`）
- 方法名：动词+名词（如 `ContinueWriting`、`GetChatSessions`）

### 2. 错误处理
```go
if err != nil {
    shared.Error(c, http.StatusInternalServerError, "操作失败", err.Error())
    return
}
```

### 3. 参数验证
使用 `binding` 标签进行参数验证：
```go
type Request struct {
    Field string `json:"field" binding:"required"`
}
```

### 4. Context传递
在Gin Context中设置信息供中间件使用：
```go
c.Set("requestID", requestID)
c.Set("tokensUsed", tokensUsed)
c.Set("aiModel", model)
```

---

## 🚀 扩展建议

### 未来可添加的API模块

1. **OutlineApi** (`outline_api.go`)
   - 生成大纲
   - 扩展大纲
   - 优化大纲结构

2. **CharacterApi** (`character_api.go`)
   - 生成角色卡
   - 角色关系分析
   - 角色发展建议

3. **WorldbuildingApi** (`worldbuilding_api.go`)
   - 世界观设定生成
   - 地点描述
   - 背景故事

4. **AnalysisApi** (`analysis_api.go`)
   - 文本质量分析
   - 情节连贯性检查
   - 风格一致性分析

---

## 📚 相关文档

- [AI服务架构设计](../../../doc/design/ai/README.md)
- [AI写作API文档](../../../doc/api/ai/01.AI写作API.md)
- [配额管理设计](../../../doc/design/ai/quota/配额管理设计.md)
- [流式接口规范](../../../doc/design/ai/streaming/12.AI流式接口规范.md)

---

**版本**: v2.0  
**更新日期**: 2025-10-22  
**维护者**: AI模块开发组

