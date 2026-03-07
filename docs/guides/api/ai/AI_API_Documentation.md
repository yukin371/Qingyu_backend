# AI API 文档

## 概述

青羽智能写作系统提供了一套完整的AI辅助写作API，支持内容生成、文本分析、续写优化、大纲生成等功能。

## 基础信息

- **基础URL**: `/api/v1/ai`
- **请求格式**: JSON
- **响应格式**: JSON
- **认证方式**: JWT Token（如果需要）

## 通用响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1640995200
}
```

- `code`: 状态码，0表示成功，非0表示错误
- `message`: 响应消息
- `data`: 响应数据
- `timestamp`: 时间戳

## API 接口

### 1. 生成内容

**接口**: `POST /api/v1/ai/generate`

**描述**: 基于项目上下文和用户提示生成新内容

**请求参数**:
```json
{
  "projectId": "string",           // 必填，项目ID
  "chapterId": "string",           // 可选，章节ID
  "prompt": "string",              // 必填，生成提示词
  "options": {                     // 可选，生成选项
    "temperature": 0.7,            // 创造性程度 (0-1)
    "maxTokens": 1000,             // 最大生成长度
    "style": "string",             // 写作风格
    "genre": "string"              // 文体类型
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "content": "生成的内容文本...",
    "tokensUsed": 256,
    "model": "gpt-3.5-turbo",
    "generatedAt": "2024-01-01T12:00:00Z"
  },
  "timestamp": 1640995200
}
```

### 2. 续写内容

**接口**: `POST /api/v1/ai/continue`

**描述**: 基于当前文本内容进行智能续写

**请求参数**:
```json
{
  "projectId": "string",           // 必填，项目ID
  "chapterId": "string",           // 必填，章节ID
  "currentText": "string",         // 必填，当前文本内容
  "continueLength": 500,           // 可选，续写长度（字数）
  "options": {                     // 可选，生成选项
    "temperature": 0.7,
    "maxTokens": 1000,
    "style": "string"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "content": "续写的内容...",
    "tokensUsed": 180,
    "model": "gpt-3.5-turbo",
    "generatedAt": "2024-01-01T12:00:00Z"
  },
  "timestamp": 1640995200
}
```

### 3. 分析内容

**接口**: `POST /api/v1/ai/analyze`

**描述**: 分析文本内容，提供情节、角色、风格等方面的分析

**请求参数**:
```json
{
  "content": "string",             // 必填，要分析的内容
  "analysisType": "string"         // 可选，分析类型：plot/character/style/general
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "type": "plot",
    "analysis": "分析结果...",
    "tokensUsed": 120,
    "model": "gpt-3.5-turbo",
    "analyzedAt": "2024-01-01T12:00:00Z"
  },
  "timestamp": 1640995200
}
```

### 4. 优化文本

**接口**: `POST /api/v1/ai/optimize`

**描述**: 优化文本的语法、风格、流畅度等

**请求参数**:
```json
{
  "projectId": "string",           // 必填，项目ID
  "chapterId": "string",           // 可选，章节ID
  "originalText": "string",        // 必填，原始文本
  "optimizeType": "string",        // 可选，优化类型：grammar/style/flow/dialogue
  "instructions": "string",        // 可选，具体优化指示
  "options": {                     // 可选，生成选项
    "temperature": 0.3,
    "maxTokens": 1500
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "content": "优化后的文本...",
    "tokensUsed": 200,
    "model": "gpt-3.5-turbo",
    "generatedAt": "2024-01-01T12:00:00Z"
  },
  "timestamp": 1640995200
}
```

### 5. 生成大纲

**接口**: `POST /api/v1/ai/outline`

**描述**: 基于主题和要求生成故事大纲

**请求参数**:
```json
{
  "projectId": "string",           // 必填，项目ID
  "theme": "string",               // 必填，主题
  "genre": "string",               // 可选，类型
  "length": "string",              // 可选，长度：short/medium/long
  "keyElements": ["string"],       // 可选，关键元素数组
  "options": {                     // 可选，生成选项
    "temperature": 0.8,
    "maxTokens": 2000
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "content": "生成的大纲内容...",
    "tokensUsed": 300,
    "model": "gpt-3.5-turbo",
    "generatedAt": "2024-01-01T12:00:00Z"
  },
  "timestamp": 1640995200
}
```

### 6. 获取上下文信息

**接口**: `GET /api/v1/ai/context/:projectId/:chapterId`

**描述**: 获取项目或章节的AI上下文信息

**路径参数**:
- `projectId`: 项目ID（必填）
- `chapterId`: 章节ID（可选）

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "projectId": "project123",
    "currentChapter": {
      "id": "chapter1",
      "title": "第一章",
      "content": "章节内容...",
      "plotThreads": ["主线情节"],
      "keyPoints": ["关键点1", "关键点2"]
    },
    "activeCharacters": [
      {
        "id": "char1",
        "name": "主角",
        "summary": "角色描述",
        "traits": ["勇敢", "善良"],
        "personalityPrompt": "性格提示词",
        "speechPattern": "说话方式",
        "currentState": "当前状态"
      }
    ],
    "currentLocations": [
      {
        "id": "loc1",
        "name": "城市",
        "description": "地点描述",
        "atmosphere": "氛围"
      }
    ],
    "worldSettings": {
      "name": "世界观名称",
      "description": "世界观描述",
      "rules": ["规则1", "规则2"]
    },
    "tokenCount": 1500
  },
  "timestamp": 1640995200
}
```

### 7. 更新上下文反馈

**接口**: `POST /api/v1/ai/context/feedback`

**描述**: 更新AI上下文的用户反馈信息

**请求参数**:
```json
{
  "projectId": "string",           // 必填，项目ID
  "chapterId": "string",           // 可选，章节ID
  "feedback": "string"             // 必填，反馈内容
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "timestamp": 1640995200
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权访问 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 429 | 请求频率超限 |
| 500 | 服务器内部错误 |
| 502 | 外部AI服务不可用 |
| 503 | 服务暂时不可用 |

## 使用示例

### JavaScript/Node.js

```javascript
// 生成内容示例
const generateContent = async (projectId, prompt) => {
  const response = await fetch('/api/v1/ai/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer your-jwt-token'
    },
    body: JSON.stringify({
      projectId: projectId,
      prompt: prompt,
      options: {
        temperature: 0.7,
        maxTokens: 1000
      }
    })
  });
  
  const result = await response.json();
  return result.data;
};

// 续写内容示例
const continueWriting = async (projectId, chapterId, currentText) => {
  const response = await fetch('/api/v1/ai/continue', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer your-jwt-token'
    },
    body: JSON.stringify({
      projectId: projectId,
      chapterId: chapterId,
      currentText: currentText,
      continueLength: 500
    })
  });
  
  const result = await response.json();
  return result.data;
};
```

### Python

```python
import requests

def generate_content(project_id, prompt, token):
    url = '/api/v1/ai/generate'
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {token}'
    }
    data = {
        'projectId': project_id,
        'prompt': prompt,
        'options': {
            'temperature': 0.7,
            'maxTokens': 1000
        }
    }
    
    response = requests.post(url, json=data, headers=headers)
    return response.json()['data']

def analyze_content(content, analysis_type='general', token=None):
    url = '/api/v1/ai/analyze'
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {token}' if token else None
    }
    data = {
        'content': content,
        'analysisType': analysis_type
    }
    
    response = requests.post(url, json=data, headers=headers)
    return response.json()['data']
```

## 配置说明

### 环境变量配置

在使用AI服务前，需要配置以下环境变量：

```bash
# AI服务提供商
AI_PROVIDER=openai

# API密钥（必填）
AI_API_KEY=your_api_key_here

# API基础URL
AI_BASE_URL=https://api.openai.com/v1

# 默认模型
AI_DEFAULT_MODEL=gpt-3.5-turbo

# 请求超时时间（秒）
AI_TIMEOUT=30

# 最大重试次数
AI_MAX_RETRIES=3
```

### 限流说明

为了保护服务稳定性，API实施了以下限流策略：

- 每分钟最多60次请求
- 每小时最多1000次请求
- 每天最多10000次请求
- 突发请求最多10次

## 最佳实践

1. **合理设置参数**: 根据需求调整temperature和maxTokens参数
2. **错误处理**: 实现完善的错误处理和重试机制
3. **缓存策略**: 对相同请求进行缓存以提高性能
4. **上下文管理**: 合理利用项目上下文信息提高生成质量
5. **分批处理**: 对于大量文本，建议分批处理避免超时

## 常见问题

### Q: 如何提高生成内容的质量？
A: 
- 提供详细的项目上下文信息
- 使用具体明确的提示词
- 合理设置temperature参数
- 利用角色和场景信息

### Q: 如何处理API请求超时？
A: 
- 检查网络连接
- 适当增加超时时间配置
- 实现重试机制
- 考虑分批处理大量内容

### Q: 如何优化API调用成本？
A: 
- 使用缓存避免重复请求
- 合理设置maxTokens参数
- 选择合适的模型
- 批量处理相关请求

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 支持基础的内容生成、分析、续写功能
- 实现上下文管理和配置系统