# AI写作辅助API测试示例

本文档提供了所有AI写作辅助API的详细测试示例，包括使用curl和Postman的测试方法。

## 前置准备

### 1. 获取认证Token
```bash
# 首先需要登录获取JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test@example.com",
    "password": "password123"
  }'

# 保存返回的token
export TOKEN="your_jwt_token_here"
```

### 2. 设置环境变量
```bash
export API_BASE="http://localhost:8080"
export API_HEADER="Authorization: Bearer $TOKEN"
```

---

## 一、内容总结API测试

### 1.1 总结文档内容 - 简短摘要

**请求示例**:
```bash
curl -X POST $API_BASE/api/v1/ai/writing/summarize \
  -H "$API_HEADER" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "人工智能（Artificial Intelligence，简称AI）是计算机科学的一个分支，它企图了解智能的实质，并生产出一种新的能以人类智能相似的方式做出反应的智能机器。该领域的研究包括机器人、语言识别、图像识别、自然语言处理和专家系统等。人工智能从诞生以来，理论和技术日益成熟，应用领域也不断扩大，可以设想，未来人工智能带来的科技产品，将会是人类智慧的"容器"。人工智能是对人的意识、思维的信息过程的模拟。人工智能不是人的智能，但能像人那样思考、也可能超过人的智能。",
    "summaryType": "brief",
    "maxLength": 100,
    "includeQuotes": false
  }'
```

**Postman配置**:
- Method: `POST`
- URL: `{{baseUrl}}/api/v1/ai/writing/summarize`
- Headers:
  - `Authorization`: `Bearer {{token}}`
  - `Content-Type`: `application/json`
- Body (raw JSON):
```json
{
  "content": "人工智能（Artificial Intelligence，简称AI）是计算机科学的一个分支...",
  "summaryType": "brief",
  "maxLength": 100,
  "includeQuotes": false
}
```

### 1.2 总结文档内容 - 详细摘要

**请求示例**:
```bash
curl -X POST $API_BASE/api/v1/ai/writing/summarize \
  -H "$API_HEADER" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "长篇文章内容...",
    "summaryType": "detailed",
    "maxLength": 500,
    "includeQuotes": true
  }'
```

### 1.3 总结文档内容 - 关键点提取

**请求示例**:
```bash
curl -X POST $API_BASE/api/v1/ai/writing/summarize \
  -H "$API_HEADER" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "需要提取关键点的文章内容...",
    "summaryType": "keypoints"
  }'
```

### 1.4 总结章节内容

**请求示例**:
```bash
curl -X POST $API_BASE/api/v1/ai/writing/summarize-chapter \
  -H "$API_HEADER" \
  -H "Content-Type: application/json" \
  -d '{
    "projectId": "proj_12345",
    "chapterId": "ch_67890",
    "outlineLevel": 3
  }'
```

**预期响应**:
```json
{
  "code": 200,
  "message": "章节总结成功",
  "data": {
    "chapterId": "ch_67890",
    "chapterTitle": "第三章：命运的转折",
    "summary": "本章描述了主人公在关键时刻做出的抉择...",
    "keyPoints": [
      "主人公发现了重要线索",
      "与反派角色初次交锋"
    ],
    "plotOutline": [],
    "characters": [],
    "tokensUsed": 1200,
    "processedAt": "2026-01-03T10:30:00Z"
  }
}
```

---

## 二、文本校对API测试

### 2.1 完整校对（检查所有类型）

**请求示例**:
```bash
curl -X POST $API_BASE/api/v1/ai/writing/proofread \
  -H "$API_HEADER" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "他们是一群热爱编程的年轻人。每天都在努力学习新的技术。他们相信，通过不懈的努力，一定能够实现自己的梦想。这个项目对他们来说意义重大。他们希望通过这个项目，能够帮助到更多的人。",
    "checkTypes": ["spelling", "grammar", "punctuation", "style"],
    "suggestions": true
  }'
```

**Postman配置**:
- Method: `POST`
- URL: `{{baseUrl}}/api/v1/ai/writing/proofread`
- Headers:
  - `Authorization`: `Bearer {{token}}`
  - `Content-Type`: `application/json`
- Body (raw JSON):
```json
{
  "content": "他们是一群热爱编程的年轻人...",
  "checkTypes": ["spelling", "grammar", "punctuation"],
  "suggestions": true
}
```

**预期响应**:
```json
{
  "code": 200,
  "message": "校对完成",
  "data": {
    "originalContent": "原文...",
    "issues": [
      {
        "id": "issue_001",
        "type": "grammar",
        "severity": "error",
        "message": "语法错误",
        "position": {
          "line": 1,
          "column": 1,
          "start": 0,
          "end": 10,
          "length": 10
        },
        "originalText": "原文片段",
        "suggestions": ["修改建议1", "修改建议2"]
      }
    ],
    "score": 85.5,
    "statistics": {
      "totalIssues": 5,
      "errorCount": 2,
      "warningCount": 2,
      "suggestionCount": 1,
      "issuesByType": {
        "grammar": 2,
        "punctuation": 2,
        "spelling": 1
      },
      "wordCount": 50,
      "characterCount": 200
    },
    "tokensUsed": 600,
    "model": "gpt-3.5-turbo",
    "processedAt": "2026-01-03T10:30:00Z"
  }
}
```

### 2.2 只检查语法和拼写

**请求示例**:
```bash
curl -X POST $API_BASE/api/v1/ai/writing/proofread \
  -H "$API_HEADER" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "需要检查的文本内容...",
    "checkTypes": ["grammar", "spelling"]
  }'
```

### 2.3 获取校对建议详情

**请求示例**:
```bash
curl -X GET $API_BASE/api/v1/ai/writing/suggestions/issue_001 \
  -H "$API_HEADER"
```

**预期响应**:
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
    "examples": ["正确示例1", "正确示例2"]
  }
}
```

---

## 三、敏感词检测API测试

### 3.1 检测所有类型的敏感词

**请求示例**:
```bash
curl -X POST $API_BASE/api/v1/ai/audit/sensitive-words \
  -H "$API_HEADER" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "这是一篇正常的文章内容，用于测试敏感词检测功能。",
    "category": "all"
  }'
```

**Postman配置**:
- Method: `POST`
- URL: `{{baseUrl}}/api/v1/ai/audit/sensitive-words`
- Headers:
  - `Authorization`: `Bearer {{token}}`
  - `Content-Type`: `application/json`
- Body (raw JSON):
```json
{
  "content": "待检测的内容...",
  "category": "all",
  "customWords": []
}
```

**预期响应**:
```json
{
  "code": 200,
  "message": "检测完成",
  "data": {
    "checkId": "check_abc123",
    "isSafe": true,
    "totalMatches": 0,
    "sensitiveWords": [],
    "summary": {
      "byCategory": {},
      "byLevel": {},
      "highRiskCount": 0,
      "mediumRiskCount": 0,
      "lowRiskCount": 0
    },
    "tokensUsed": 0,
    "processedAt": "2026-01-03T10:30:00Z"
  }
}
```

### 3.2 只检测政治敏感词

**请求示例**:
```bash
curl -X POST $API_BASE/api/v1/ai/audit/sensitive-words \
  -H "$API_HEADER" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "待检测的内容...",
    "category": "political"
  }'
```

### 3.3 使用自定义敏感词检测

**请求示例**:
```bash
curl -X POST $API_BASE/api/v1/ai/audit/sensitive-words \
  -H "$API_HEADER" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "这篇文章包含一些特定词汇...",
    "customWords": ["测试词1", "测试词2", "测试词3"],
    "category": "all"
  }'
```

### 3.4 获取敏感词检测结果

**请求示例**:
```bash
curl -X GET $API_BASE/api/v1/ai/audit/sensitive-words/check_abc123 \
  -H "$API_HEADER"
```

**预期响应**:
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
        "category": "custom",
        "level": "medium",
        "position": {
          "line": 1,
          "column": 10,
          "start": 10,
          "end": 15,
          "length": 5
        },
        "context": "...前文...敏感词...后文...",
        "suggestion": "建议修改或删除敏感词"
      }
    ],
    "customWords": ["测试词1", "测试词2"],
    "summary": {
      "byCategory": {
        "custom": 1
      },
      "byLevel": {
        "medium": 1
      },
      "highRiskCount": 0,
      "mediumRiskCount": 1,
      "lowRiskCount": 0
    },
    "createdAt": "2026-01-03T10:30:00Z",
    "expiresAt": "2026-02-02T10:30:00Z"
  }
}
```

---

## 四、错误场景测试

### 4.1 未授权访问

**请求示例**:
```bash
# 不提供token
curl -X POST $API_BASE/api/v1/ai/writing/summarize \
  -H "Content-Type: application/json" \
  -d '{"content": "测试内容"}'
```

**预期响应**:
```json
{
  "code": 401,
  "message": "未授权",
  "error": "请先登录或提供有效的访问凭证"
}
```

### 4.2 参数验证失败

**请求示例**:
```bash
curl -X POST $API_BASE/api/v1/ai/writing/summarize \
  -H "$API_HEADER" \
  -H "Content-Type: application/json" \
  -d '{
    "summaryType": "brief"
  }'
```

**预期响应**:
```json
{
  "code": 400,
  "message": "参数错误",
  "error": "Content is required"
}
```

### 4.3 配额不足

**请求示例**:
```bash
# 使用配额不足的账号
curl -X POST $API_BASE/api/v1/ai/writing/summarize \
  -H "Authorization: Bearer low_quota_token" \
  -H "Content-Type: application/json" \
  -d '{"content": "测试内容"}'
```

**预期响应**:
```json
{
  "code": 403,
  "message": "配额不足",
  "error": "您的AI配额已用完，请充值后再试"
}
```

---

## 五、性能测试

### 5.1 长文本总结测试

**请求示例**:
```bash
# 测试5000字长文本
curl -X POST $API_BASE/api/v1/ai/writing/summarize \
  -H "$API_HEADER" \
  -H "Content-Type: application/json" \
  -d @long_text.json
```

### 5.2 并发测试

**使用Apache Bench**:
```bash
# 100个并发请求，总共1000个请求
ab -n 1000 -c 100 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -p summarize_request.json \
  $API_BASE/api/v1/ai/writing/summarize
```

**使用wrk**:
```bash
# 10线程，100个连接，持续30秒
wrk -t10 -c100 -d30s \
  -H "Authorization: Bearer $TOKEN" \
  -s summarize.lua \
  $API_BASE/api/v1/ai/writing/summarize
```

---

## 六、Postman测试集合

### 导入JSON

创建文件 `qingyu_ai_tests.postman_collection.json`:

```json
{
  "info": {
    "name": "青羽AI写作辅助API测试",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080"
    },
    {
      "key": "token",
      "value": "your_jwt_token_here"
    }
  ],
  "item": [
    {
      "name": "内容总结",
      "item": [
        {
          "name": "简短摘要",
          "request": {
            "method": "POST",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{token}}"
              },
              {
                "key": "Content-Type",
                "value": "application/json"
              }
            ],
            "url": "{{baseUrl}}/api/v1/ai/writing/summarize",
            "body": {
              "mode": "raw",
              "raw": "{\n  \"content\": \"测试内容...\",\n  \"summaryType\": \"brief\"\n}"
            }
          }
        }
      ]
    },
    {
      "name": "文本校对",
      "item": [
        {
          "name": "完整校对",
          "request": {
            "method": "POST",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{token}}"
              },
              {
                "key": "Content-Type",
                "value": "application/json"
              }
            ],
            "url": "{{baseUrl}}/api/v1/ai/writing/proofread",
            "body": {
              "mode": "raw",
              "raw": "{\n  \"content\": \"待校对文本...\",\n  \"checkTypes\": [\"grammar\", \"spelling\"]\n}"
            }
          }
        }
      ]
    },
    {
      "name": "敏感词检测",
      "item": [
        {
          "name": "全面检测",
          "request": {
            "method": "POST",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{token}}"
              },
              {
                "key": "Content-Type",
                "value": "application/json"
              }
            ],
            "url": "{{baseUrl}}/api/v1/ai/audit/sensitive-words",
            "body": {
              "mode": "raw",
              "raw": "{\n  \"content\": \"待检测内容...\",\n  \"category\": \"all\"\n}"
            }
          }
        }
      ]
    }
  ]
}
```

---

## 七、测试脚本

### Bash测试脚本

创建文件 `test_ai_apis.sh`:

```bash
#!/bin/bash

# 配置
API_BASE="http://localhost:8080"
TOKEN="your_jwt_token_here"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local name=$1
    local url=$2
    local data=$3

    echo -e "\n${GREEN}测试: $name${NC}"
    response=$(curl -s -X POST "$API_BASE$url" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "$data")

    if echo "$response" | grep -q '"code":200'; then
        echo -e "${GREEN}✓ 通过${NC}"
    else
        echo -e "${RED}✗ 失败${NC}"
        echo "$response"
    fi
}

# 执行测试
test_api "简短摘要" "/api/v1/ai/writing/summarize" \
    '{"content": "测试内容", "summaryType": "brief"}'

test_api "文本校对" "/api/v1/ai/writing/proofread" \
    '{"content": "测试内容", "checkTypes": ["grammar"]}'

test_api "敏感词检测" "/api/v1/ai/audit/sensitive-words" \
    '{"content": "测试内容", "category": "all"}'

echo -e "\n${GREEN}测试完成${NC}"
```

**使用方法**:
```bash
chmod +x test_ai_apis.sh
./test_ai_apis.sh
```

---

## 八、响应时间基准

### 预期响应时间

| API端点 | 内容长度 | 预期时间 |
|--------|---------|---------|
| 简短摘要 | 500字 | < 2秒 |
| 详细摘要 | 5000字 | < 5秒 |
| 章节总结 | 3000字 | < 4秒 |
| 文本校对 | 1000字 | < 3秒 |
| 敏感词检测 | 5000字 | < 1秒 |

### 性能优化建议

1. **缓存结果**: 对相同内容的重复请求使用缓存
2. **批量处理**: 对多个小文本合并处理
3. **流式响应**: 对长文本使用流式返回
4. **异步处理**: 对耗时操作使用异步任务队列

---

## 九、监控和日志

### 查看API日志

```bash
# 查看实时日志
tail -f /var/log/qingyu/backend.log | grep "ai_service"

# 查看错误日志
grep "ERROR" /var/log/qingyu/backend.log | grep "ai"

# 统计API调用次数
grep "summarize" /var/log/qingyu/backend.log | wc -l
```

### 监控指标

- **响应时间**: p50, p95, p99
- **成功率**: 200状态码占比
- **错误率**: 4xx/5xx状态码占比
- **Token消耗**: 每日Token使用量
- **并发数**: 同时处理的请求数

---

## 总结

本文档提供了完整的API测试指南，包括：

- ✅ 6个API端点的详细测试示例
- ✅ curl和Postman两种测试方式
- ✅ 正常场景和错误场景测试
- ✅ 性能测试和基准
- ✅ 自动化测试脚本
- ✅ 监控和日志查看方法

建议按照本文档的顺序进行测试，确保所有功能正常运行。
