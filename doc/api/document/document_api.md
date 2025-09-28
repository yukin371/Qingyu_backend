# Document API 设计文档

## 1. 概述

本文档定义了青羽项目中文档管理相关的API接口设计，包括项目管理、文档结构、版本控制等功能。

## 2. 数据模型

### 2.1 核心实体

#### Project（项目）
```json
{
  "id": "string",
  "ownerId": "string",
  "name": "string",
  "status": "public|private",
  "description": "string",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Node（节点）
```json
{
  "id": "string",
  "projectId": "string",
  "parentId": "string",
  "type": "folder|file",
  "name": "string",
  "slug": "string",
  "relativePath": "string",
  "order": 0,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z",
  "children": []
}
```

#### Document（文档）
```json
{
  "id": "string",
  "projectId": "string",
  "nodeId": "string",
  "title": "string",
  "content": "string",
  "format": "markdown|txt|json",
  "words": 0,
  "version": 1,
  "aiContext": "string",
  "plotThreads": ["string"],
  "keyPoints": ["string"],
  "writingHints": "string",
  "characterIds": ["string"],
  "locationIds": ["string"],
  "timelineIds": ["string"],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

## 3. API 接口设计

### 3.1 项目管理

#### 3.1.1 创建项目
```
POST /api/v1/projects
```

**请求体：**
```json
{
  "name": "我的小说项目",
  "status": "private",
  "description": "这是一个科幻小说项目"
}
```

**响应：**
```json
{
  "code": 201,
  "message": "success",
  "data": {
    "id": "proj_123456",
    "ownerId": "user_123",
    "name": "我的小说项目",
    "status": "private",
    "description": "这是一个科幻小说项目",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "req_123456"
}
```

#### 3.1.2 获取项目列表
```
GET /api/v1/projects?page=1&page_size=20&status=private&search=小说
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "proj_123456",
        "ownerId": "user_123",
        "name": "我的小说项目",
        "status": "private",
        "description": "这是一个科幻小说项目",
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 1,
      "total_pages": 1
    }
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "req_123456"
}
```

#### 3.1.3 获取项目详情
```
GET /api/v1/projects/{projectId}
```

#### 3.1.4 更新项目
```
PUT /api/v1/projects/{projectId}
```

#### 3.1.5 删除项目
```
DELETE /api/v1/projects/{projectId}
```

### 3.2 节点管理（目录/文件结构）

#### 3.2.1 获取项目结构树
```
GET /api/v1/projects/{projectId}/nodes?type=tree
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "node_123",
      "projectId": "proj_123456",
      "parentId": null,
      "type": "folder",
      "name": "第一卷",
      "relativePath": "第一卷",
      "order": 1,
      "children": [
        {
          "id": "node_124",
          "projectId": "proj_123456",
          "parentId": "node_123",
          "type": "file",
          "name": "第一章.md",
          "relativePath": "第一卷/第一章.md",
          "order": 1,
          "children": []
        }
      ]
    }
  ]
}
```

#### 3.2.2 创建节点
```
POST /api/v1/projects/{projectId}/nodes
```

**请求体：**
```json
{
  "parentId": "node_123",
  "type": "file",
  "name": "第二章.md",
  "order": 2
}
```

#### 3.2.3 更新节点
```
PUT /api/v1/projects/{projectId}/nodes/{nodeId}
```

#### 3.2.4 移动节点
```
PUT /api/v1/projects/{projectId}/nodes/{nodeId}/move
```

**请求体：**
```json
{
  "newParentId": "node_456",
  "newOrder": 3
}
```

#### 3.2.5 删除节点
```
DELETE /api/v1/projects/{projectId}/nodes/{nodeId}
```

### 3.3 文档管理

#### 3.3.1 创建文档
```
POST /api/v1/projects/{projectId}/documents
```

**请求体：**
```json
{
  "nodeId": "node_124",
  "title": "第一章：开端",
  "content": "这是第一章的内容...",
  "format": "markdown",
  "aiContext": "科幻背景，主角是程序员",
  "plotThreads": ["主线剧情"],
  "keyPoints": ["角色介绍", "世界观设定"],
  "writingHints": "注意节奏控制",
  "characterIds": ["char_123"],
  "locationIds": ["loc_123"]
}
```

#### 3.3.2 获取文档详情
```
GET /api/v1/projects/{projectId}/documents/{documentId}
```

#### 3.3.3 更新文档
```
PUT /api/v1/projects/{projectId}/documents/{documentId}
```

**请求体：**
```json
{
  "title": "第一章：新的开端",
  "content": "更新后的内容...",
  "expectedVersion": 1,
  "aiContext": "更新的AI上下文",
  "plotThreads": ["主线剧情", "支线剧情"],
  "keyPoints": ["角色介绍", "世界观设定", "冲突设置"]
}
```

#### 3.3.4 删除文档
```
DELETE /api/v1/projects/{projectId}/documents/{documentId}
```

#### 3.3.5 获取文档列表
```
GET /api/v1/projects/{projectId}/documents?page=1&page_size=20&search=第一章
```

### 3.4 版本控制

#### 3.4.1 创建版本
```
POST /api/v1/projects/{projectId}/nodes/{nodeId}/versions
```

**请求体：**
```json
{
  "authorId": "user_123",
  "message": "更新第一章内容",
  "content": "新的文档内容...",
  "expectedVersion": 1
}
```

#### 3.4.2 获取版本历史
```
GET /api/v1/projects/{projectId}/nodes/{nodeId}/versions?page=1&page_size=20
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "rev_123",
        "projectId": "proj_123456",
        "nodeId": "node_124",
        "commitId": "commit_123",
        "version": 2,
        "authorId": "user_123",
        "message": "更新第一章内容",
        "parentVersion": 1,
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 2,
      "total_pages": 1
    }
  }
}
```

#### 3.4.3 版本回滚
```
POST /api/v1/projects/{projectId}/nodes/{nodeId}/rollback
```

**请求体：**
```json
{
  "authorId": "user_123",
  "version": 1,
  "message": "回滚到初始版本"
}
```

#### 3.4.4 获取当前版本
```
GET /api/v1/projects/{projectId}/nodes/{nodeId}/versions/current
```

### 3.5 批量提交

#### 3.5.1 创建批量提交
```
POST /api/v1/projects/{projectId}/commits
```

**请求体：**
```json
{
  "authorId": "user_123",
  "message": "批量更新多个章节",
  "files": [
    {
      "nodeId": "node_124",
      "content": "第一章更新内容...",
      "expectedVersion": 1
    },
    {
      "nodeId": "node_125",
      "content": "第二章更新内容...",
      "expectedVersion": 1
    }
  ]
}
```

#### 3.5.2 获取提交历史
```
GET /api/v1/projects/{projectId}/commits?page=1&page_size=20&authorId=user_123
```

#### 3.5.3 获取提交详情
```
GET /api/v1/projects/{projectId}/commits/{commitId}
```

### 3.6 冲突检测与解决

#### 3.6.1 检测版本冲突
```
POST /api/v1/projects/{projectId}/conflicts/detect
```

**请求体：**
```json
{
  "files": [
    {
      "nodeId": "node_124",
      "expectedVersion": 1
    },
    {
      "nodeId": "node_125",
      "expectedVersion": 2
    }
  ]
}
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "projectId": "proj_123456",
    "hasConflicts": true,
    "conflicts": {
      "node_124": {
        "hasConflict": true,
        "currentVersion": 3,
        "expectedVersion": 1,
        "conflictingRevisions": [
          {
            "id": "rev_124",
            "version": 2,
            "authorId": "user_456",
            "message": "其他用户的修改",
            "createdAt": "2024-01-01T01:00:00Z"
          }
        ],
        "lastModified": "2024-01-01T01:00:00Z"
      }
    }
  }
}
```

#### 3.6.2 解决冲突
```
POST /api/v1/projects/{projectId}/conflicts/resolve
```

**请求体：**
```json
{
  "projectId": "proj_123456",
  "authorId": "user_123",
  "message": "解决版本冲突",
  "resolutions": {
    "node_124": {
      "strategy": "manual",
      "resolvedBy": "user_123",
      "resolution": "手动合并冲突内容",
      "mergedContent": "合并后的最终内容..."
    }
  }
}
```

### 3.7 补丁管理

#### 3.7.1 创建补丁
```
POST /api/v1/projects/{projectId}/nodes/{nodeId}/patches
```

**请求体：**
```json
{
  "createdBy": "ai_assistant",
  "diffFormat": "unified",
  "diffPayload": "@@ -1,3 +1,4 @@\n 第一行\n+新增行\n 第二行\n 第三行",
  "baseVersion": 1,
  "message": "AI建议的内容优化"
}
```

#### 3.7.2 应用补丁
```
POST /api/v1/projects/{projectId}/nodes/{nodeId}/patches/{patchId}/apply
```

**查询参数：**
- `applier`: 应用者ID

#### 3.7.3 获取补丁列表
```
GET /api/v1/projects/{projectId}/nodes/{nodeId}/patches?status=pending&page=1&page_size=20
```

## 4. 错误处理

### 4.1 常见错误码

| 错误码 | HTTP状态码 | 描述 |
|--------|------------|------|
| 40001 | 400 | 请求参数错误 |
| 40101 | 401 | 未认证 |
| 40301 | 403 | 无权限访问 |
| 40401 | 404 | 资源不存在 |
| 40901 | 409 | 版本冲突 |
| 42201 | 422 | 业务逻辑错误 |
| 50001 | 500 | 服务器内部错误 |

### 4.2 错误响应格式

```json
{
  "code": 409,
  "message": "版本冲突",
  "error": {
    "type": "version_conflict",
    "details": [
      {
        "field": "version",
        "message": "当前版本为3，期望版本为1"
      }
    ]
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "req_123456"
}
```

## 5. 认证与授权

### 5.1 认证方式
- 使用JWT Bearer Token
- Token在Authorization头中传递：`Authorization: Bearer {token}`

### 5.2 权限控制
- 项目所有者拥有完全权限
- 协作者根据角色拥有不同权限
- 公开项目允许只读访问

## 6. 性能优化

### 6.1 分页
- 默认分页大小：20
- 最大分页大小：100
- 使用cursor分页优化大数据集查询

### 6.2 缓存策略
- 项目结构树缓存30分钟
- 文档内容缓存10分钟
- 版本历史缓存1小时

### 6.3 限流
- 每用户每分钟最多100次API调用
- 批量操作限制每次最多50个文件

## 7. 监控指标

- API响应时间
- 版本冲突率
- 文档更新频率
- 用户活跃度