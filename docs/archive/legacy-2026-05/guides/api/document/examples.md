# Document API 使用示例

本文档提供了青羽文档管理系统API的详细使用示例，包括请求和响应格式。

## 目录

- [认证](#认证)
- [项目管理](#项目管理)
- [节点管理](#节点管理)
- [文档管理](#文档管理)
- [版本控制](#版本控制)
- [角色管理](#角色管理)
- [错误处理](#错误处理)

## 认证

所有API请求都需要在请求头中包含JWT令牌：

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 项目管理

### 创建项目

**请求示例：**

```http
POST /api/v1/projects
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "我的小说项目",
  "status": "private",
  "description": "这是一个科幻小说项目"
}
```

**响应示例：**

```json
{
  "code": 201,
  "message": "项目创建成功",
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "req_abc123",
  "data": {
    "id": "proj_123456",
    "ownerId": "user_789",
    "name": "我的小说项目",
    "status": "private",
    "description": "这是一个科幻小说项目",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  }
}
```

### 获取项目列表

**请求示例：**

```http
GET /api/v1/projects?page=1&page_size=10&status=private&search=小说
Authorization: Bearer <token>
```

**响应示例：**

```json
{
  "code": 200,
  "message": "获取项目列表成功",
  "timestamp": "2024-01-15T10:35:00Z",
  "request_id": "req_def456",
  "data": {
    "items": [
      {
        "id": "proj_123456",
        "ownerId": "user_789",
        "name": "我的小说项目",
        "status": "private",
        "description": "这是一个科幻小说项目",
        "createdAt": "2024-01-15T10:30:00Z",
        "updatedAt": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

### 更新项目

**请求示例：**

```http
PUT /api/v1/projects/proj_123456
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "我的科幻小说项目",
  "description": "这是一个关于未来世界的科幻小说项目"
}
```

**响应示例：**

```json
{
  "code": 200,
  "message": "项目更新成功",
  "timestamp": "2024-01-15T10:40:00Z",
  "request_id": "req_ghi789",
  "data": {
    "id": "proj_123456",
    "ownerId": "user_789",
    "name": "我的科幻小说项目",
    "status": "private",
    "description": "这是一个关于未来世界的科幻小说项目",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:40:00Z"
  }
}
```

## 节点管理

### 创建文件夹节点

**请求示例：**

```http
POST /api/v1/projects/proj_123456/nodes
Content-Type: application/json
Authorization: Bearer <token>

{
  "type": "folder",
  "name": "第一章",
  "order": 1
}
```

**响应示例：**

```json
{
  "code": 201,
  "message": "节点创建成功",
  "timestamp": "2024-01-15T11:00:00Z",
  "request_id": "req_jkl012",
  "data": {
    "id": "node_folder_001",
    "projectId": "proj_123456",
    "parentId": null,
    "type": "folder",
    "name": "第一章",
    "slug": "chapter-1",
    "relativePath": "/第一章",
    "order": 1,
    "createdAt": "2024-01-15T11:00:00Z",
    "updatedAt": "2024-01-15T11:00:00Z"
  }
}
```

### 创建文件节点

**请求示例：**

```http
POST /api/v1/projects/proj_123456/nodes
Content-Type: application/json
Authorization: Bearer <token>

{
  "parentId": "node_folder_001",
  "type": "file",
  "name": "开场",
  "order": 1
}
```

**响应示例：**

```json
{
  "code": 201,
  "message": "节点创建成功",
  "timestamp": "2024-01-15T11:05:00Z",
  "request_id": "req_mno345",
  "data": {
    "id": "node_file_001",
    "projectId": "proj_123456",
    "parentId": "node_folder_001",
    "type": "file",
    "name": "开场",
    "slug": "opening",
    "relativePath": "/第一章/开场",
    "order": 1,
    "createdAt": "2024-01-15T11:05:00Z",
    "updatedAt": "2024-01-15T11:05:00Z"
  }
}
```

### 获取项目结构树

**请求示例：**

```http
GET /api/v1/projects/proj_123456/nodes?type=tree
Authorization: Bearer <token>
```

**响应示例：**

```json
{
  "code": 200,
  "message": "获取项目结构成功",
  "timestamp": "2024-01-15T11:10:00Z",
  "request_id": "req_pqr678",
  "data": [
    {
      "id": "node_folder_001",
      "projectId": "proj_123456",
      "parentId": null,
      "type": "folder",
      "name": "第一章",
      "slug": "chapter-1",
      "relativePath": "/第一章",
      "order": 1,
      "createdAt": "2024-01-15T11:00:00Z",
      "updatedAt": "2024-01-15T11:00:00Z",
      "children": [
        {
          "id": "node_file_001",
          "projectId": "proj_123456",
          "parentId": "node_folder_001",
          "type": "file",
          "name": "开场",
          "slug": "opening",
          "relativePath": "/第一章/开场",
          "order": 1,
          "createdAt": "2024-01-15T11:05:00Z",
          "updatedAt": "2024-01-15T11:05:00Z",
          "children": []
        }
      ]
    }
  ]
}
```

## 文档管理

### 创建文档

**请求示例：**

```http
POST /api/v1/projects/proj_123456/documents
Content-Type: application/json
Authorization: Bearer <token>

{
  "nodeId": "node_file_001",
  "title": "开场",
  "content": "# 开场\n\n2087年，地球已经不再是人类唯一的家园...",
  "format": "markdown",
  "aiContext": "科幻小说，未来世界背景",
  "plotThreads": ["主角登场", "世界观介绍"],
  "keyPoints": ["时间设定", "主角身份"],
  "writingHints": "注意营造未来感",
  "characterIds": ["char_001"],
  "locationIds": ["loc_001"],
  "timelineIds": ["timeline_001"]
}
```

**响应示例：**

```json
{
  "code": 201,
  "message": "文档创建成功",
  "timestamp": "2024-01-15T11:15:00Z",
  "request_id": "req_stu901",
  "data": {
    "id": "doc_001",
    "projectId": "proj_123456",
    "nodeId": "node_file_001",
    "title": "开场",
    "content": "# 开场\n\n2087年，地球已经不再是人类唯一的家园...",
    "format": "markdown",
    "words": 156,
    "version": 1,
    "aiContext": "科幻小说，未来世界背景",
    "plotThreads": ["主角登场", "世界观介绍"],
    "keyPoints": ["时间设定", "主角身份"],
    "writingHints": "注意营造未来感",
    "characterIds": ["char_001"],
    "locationIds": ["loc_001"],
    "timelineIds": ["timeline_001"],
    "createdAt": "2024-01-15T11:15:00Z",
    "updatedAt": "2024-01-15T11:15:00Z"
  }
}
```

### 更新文档（乐观锁）

**请求示例：**

```http
PUT /api/v1/projects/proj_123456/documents/doc_001
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "开场 - 修订版",
  "content": "# 开场 - 修订版\n\n2087年，地球已经不再是人类唯一的家园。随着科技的发展...",
  "expectedVersion": 1,
  "plotThreads": ["主角登场", "世界观介绍", "科技背景"],
  "keyPoints": ["时间设定", "主角身份", "科技水平"]
}
```

**响应示例：**

```json
{
  "code": 200,
  "message": "文档更新成功",
  "timestamp": "2024-01-15T11:20:00Z",
  "request_id": "req_vwx234",
  "data": {
    "id": "doc_001",
    "projectId": "proj_123456",
    "nodeId": "node_file_001",
    "title": "开场 - 修订版",
    "content": "# 开场 - 修订版\n\n2087年，地球已经不再是人类唯一的家园。随着科技的发展...",
    "format": "markdown",
    "words": 189,
    "version": 2,
    "aiContext": "科幻小说，未来世界背景",
    "plotThreads": ["主角登场", "世界观介绍", "科技背景"],
    "keyPoints": ["时间设定", "主角身份", "科技水平"],
    "writingHints": "注意营造未来感",
    "characterIds": ["char_001"],
    "locationIds": ["loc_001"],
    "timelineIds": ["timeline_001"],
    "createdAt": "2024-01-15T11:15:00Z",
    "updatedAt": "2024-01-15T11:20:00Z"
  }
}
```

### 版本冲突处理

**请求示例（版本冲突）：**

```http
PUT /api/v1/projects/proj_123456/documents/doc_001
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "开场 - 另一个版本",
  "content": "不同的内容...",
  "expectedVersion": 1
}
```

**响应示例（冲突）：**

```json
{
  "code": 409,
  "message": "版本冲突",
  "timestamp": "2024-01-15T11:25:00Z",
  "request_id": "req_yzab567",
  "error": {
    "type": "version_conflict",
    "details": [
      {
        "field": "version",
        "message": "当前版本为2，期望版本为1"
      }
    ]
  }
}
```

## 版本控制

### 创建版本

**请求示例：**

```http
POST /api/v1/projects/proj_123456/nodes/node_file_001/versions
Content-Type: application/json
Authorization: Bearer <token>

{
  "authorId": "user_789",
  "message": "添加更多世界观描述",
  "content": "# 开场 - 最新版\n\n2087年，地球已经不再是人类唯一的家园...",
  "expectedVersion": 2
}
```

**响应示例：**

```json
{
  "code": 201,
  "message": "版本创建成功",
  "timestamp": "2024-01-15T11:30:00Z",
  "request_id": "req_cdef890",
  "data": {
    "id": "rev_001",
    "projectId": "proj_123456",
    "nodeId": "node_file_001",
    "commitId": "commit_001",
    "version": 3,
    "authorId": "user_789",
    "message": "添加更多世界观描述",
    "parentVersion": 2,
    "createdAt": "2024-01-15T11:30:00Z"
  }
}
```

### 获取版本历史

**请求示例：**

```http
GET /api/v1/projects/proj_123456/nodes/node_file_001/versions?page=1&page_size=10
Authorization: Bearer <token>
```

**响应示例：**

```json
{
  "code": 200,
  "message": "获取版本历史成功",
  "timestamp": "2024-01-15T11:35:00Z",
  "request_id": "req_ghij123",
  "data": {
    "items": [
      {
        "id": "rev_001",
        "projectId": "proj_123456",
        "nodeId": "node_file_001",
        "commitId": "commit_001",
        "version": 3,
        "authorId": "user_789",
        "message": "添加更多世界观描述",
        "parentVersion": 2,
        "createdAt": "2024-01-15T11:30:00Z"
      },
      {
        "id": "rev_002",
        "projectId": "proj_123456",
        "nodeId": "node_file_001",
        "commitId": "commit_002",
        "version": 2,
        "authorId": "user_789",
        "message": "修订开场内容",
        "parentVersion": 1,
        "createdAt": "2024-01-15T11:20:00Z"
      },
      {
        "id": "rev_003",
        "projectId": "proj_123456",
        "nodeId": "node_file_001",
        "commitId": "commit_003",
        "version": 1,
        "authorId": "user_789",
        "message": "初始版本",
        "parentVersion": 0,
        "createdAt": "2024-01-15T11:15:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 3,
      "total_pages": 1
    }
  }
}
```

### 版本回滚

**请求示例：**

```http
POST /api/v1/projects/proj_123456/nodes/node_file_001/rollback
Content-Type: application/json
Authorization: Bearer <token>

{
  "authorId": "user_789",
  "version": 2,
  "message": "回滚到修订版"
}
```

**响应示例：**

```json
{
  "code": 200,
  "message": "回滚成功",
  "timestamp": "2024-01-15T11:40:00Z",
  "request_id": "req_klmn456",
  "data": {
    "id": "rev_004",
    "projectId": "proj_123456",
    "nodeId": "node_file_001",
    "commitId": "commit_004",
    "version": 4,
    "authorId": "user_789",
    "message": "回滚到修订版",
    "parentVersion": 3,
    "createdAt": "2024-01-15T11:40:00Z"
  }
}
```

## 角色管理

### 创建角色

**请求示例：**

```http
POST /api/v1/projects/proj_123456/characters
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "李明",
  "alias": ["小李", "明哥"],
  "summary": "主角，一名年轻的太空工程师",
  "traits": ["聪明", "勇敢", "好奇"],
  "background": "出生于地球，从小对太空充满向往...",
  "avatarUrl": "https://example.com/avatar/liming.jpg",
  "personalityPrompt": "性格开朗，喜欢探索未知",
  "speechPattern": "说话直接，偶尔会用技术术语",
  "currentState": "刚刚加入太空站工作"
}
```

**响应示例：**

```json
{
  "code": 201,
  "message": "角色创建成功",
  "timestamp": "2024-01-15T12:00:00Z",
  "request_id": "req_opqr789",
  "data": {
    "id": "char_001",
    "projectId": "proj_123456",
    "name": "李明",
    "alias": ["小李", "明哥"],
    "summary": "主角，一名年轻的太空工程师",
    "traits": ["聪明", "勇敢", "好奇"],
    "background": "出生于地球，从小对太空充满向往...",
    "avatarUrl": "https://example.com/avatar/liming.jpg",
    "personalityPrompt": "性格开朗，喜欢探索未知",
    "speechPattern": "说话直接，偶尔会用技术术语",
    "currentState": "刚刚加入太空站工作",
    "createdAt": "2024-01-15T12:00:00Z",
    "updatedAt": "2024-01-15T12:00:00Z"
  }
}
```

### 获取角色列表

**请求示例：**

```http
GET /api/v1/projects/proj_123456/characters?page=1&page_size=10&search=李&traits=聪明
Authorization: Bearer <token>
```

**响应示例：**

```json
{
  "code": 200,
  "message": "获取角色列表成功",
  "timestamp": "2024-01-15T12:05:00Z",
  "request_id": "req_stuv012",
  "data": {
    "items": [
      {
        "id": "char_001",
        "projectId": "proj_123456",
        "name": "李明",
        "alias": ["小李", "明哥"],
        "summary": "主角，一名年轻的太空工程师",
        "traits": ["聪明", "勇敢", "好奇"],
        "background": "出生于地球，从小对太空充满向往...",
        "avatarUrl": "https://example.com/avatar/liming.jpg",
        "personalityPrompt": "性格开朗，喜欢探索未知",
        "speechPattern": "说话直接，偶尔会用技术术语",
        "currentState": "刚刚加入太空站工作",
        "createdAt": "2024-01-15T12:00:00Z",
        "updatedAt": "2024-01-15T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

### 更新角色

**请求示例：**

```http
PUT /api/v1/projects/proj_123456/characters/char_001
Content-Type: application/json
Authorization: Bearer <token>

{
  "summary": "主角，一名经验丰富的太空工程师",
  "currentState": "已经在太空站工作了一年",
  "traits": ["聪明", "勇敢", "好奇", "经验丰富"]
}
```

**响应示例：**

```json
{
  "code": 200,
  "message": "角色更新成功",
  "timestamp": "2024-01-15T12:10:00Z",
  "request_id": "req_wxyz345",
  "data": {
    "id": "char_001",
    "projectId": "proj_123456",
    "name": "李明",
    "alias": ["小李", "明哥"],
    "summary": "主角，一名经验丰富的太空工程师",
    "traits": ["聪明", "勇敢", "好奇", "经验丰富"],
    "background": "出生于地球，从小对太空充满向往...",
    "avatarUrl": "https://example.com/avatar/liming.jpg",
    "personalityPrompt": "性格开朗，喜欢探索未知",
    "speechPattern": "说话直接，偶尔会用技术术语",
    "currentState": "已经在太空站工作了一年",
    "createdAt": "2024-01-15T12:00:00Z",
    "updatedAt": "2024-01-15T12:10:00Z"
  }
}
```

## 错误处理

### 参数验证错误

**请求示例（缺少必填字段）：**

```http
POST /api/v1/projects
Content-Type: application/json
Authorization: Bearer <token>

{
  "status": "private"
}
```

**响应示例：**

```json
{
  "code": 400,
  "message": "请求参数错误",
  "timestamp": "2024-01-15T12:15:00Z",
  "request_id": "req_abcd678",
  "error": {
    "type": "validation_error",
    "details": [
      {
        "field": "name",
        "message": "名称不能为空"
      }
    ]
  }
}
```

### 业务逻辑错误

**请求示例（项目名称重复）：**

```http
POST /api/v1/projects
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "我的科幻小说项目",
  "status": "private"
}
```

**响应示例：**

```json
{
  "code": 422,
  "message": "业务逻辑错误",
  "timestamp": "2024-01-15T12:20:00Z",
  "request_id": "req_efgh901",
  "error": {
    "type": "business_logic_error",
    "details": [
      {
        "field": "name",
        "message": "项目名称已存在"
      }
    ]
  }
}
```

### 资源不存在

**请求示例：**

```http
GET /api/v1/projects/proj_nonexistent
Authorization: Bearer <token>
```

**响应示例：**

```json
{
  "code": 404,
  "message": "资源不存在",
  "timestamp": "2024-01-15T12:25:00Z",
  "request_id": "req_ijkl234",
  "error": {
    "type": "not_found_error",
    "details": []
  }
}
```

### 认证失败

**请求示例（无效token）：**

```http
GET /api/v1/projects
Authorization: Bearer invalid_token
```

**响应示例：**

```json
{
  "code": 401,
  "message": "未认证",
  "timestamp": "2024-01-15T12:30:00Z",
  "request_id": "req_mnop567",
  "error": {
    "type": "authentication_error",
    "details": []
  }
}
```

### 服务器内部错误

**响应示例：**

```json
{
  "code": 500,
  "message": "服务器内部错误",
  "timestamp": "2024-01-15T12:35:00Z",
  "request_id": "req_qrst890",
  "error": {
    "type": "internal_server_error",
    "details": []
  }
}
```

## 注意事项

1. **认证**：所有API请求都需要有效的JWT令牌
2. **版本控制**：更新文档时建议使用乐观锁（expectedVersion）避免冲突
3. **分页**：列表接口支持分页，默认每页20条记录
4. **搜索**：支持关键词搜索，会在名称、标题等字段中匹配
5. **错误处理**：根据HTTP状态码和错误类型进行相应处理
6. **请求ID**：每个响应都包含request_id，便于问题追踪
7. **时间格式**：所有时间字段使用ISO 8601格式（UTC时区）

## SDK示例

### JavaScript/TypeScript

```javascript
// 创建项目
const createProject = async (projectData) => {
  const response = await fetch('/api/v1/projects', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify(projectData)
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message);
  }
  
  return await response.json();
};

// 更新文档（带版本控制）
const updateDocument = async (projectId, documentId, updateData) => {
  try {
    const response = await fetch(`/api/v1/projects/${projectId}/documents/${documentId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(updateData)
    });
    
    if (response.status === 409) {
      // 处理版本冲突
      const conflict = await response.json();
      console.log('版本冲突:', conflict.error.details);
      // 可以选择重新获取最新版本或提示用户
      return null;
    }
    
    return await response.json();
  } catch (error) {
    console.error('更新文档失败:', error);
    throw error;
  }
};
```

### Python

```python
import requests
import json

class DocumentAPI:
    def __init__(self, base_url, token):
        self.base_url = base_url
        self.headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {token}'
        }
    
    def create_project(self, project_data):
        response = requests.post(
            f'{self.base_url}/projects',
            headers=self.headers,
            json=project_data
        )
        response.raise_for_status()
        return response.json()
    
    def update_document(self, project_id, document_id, update_data):
        try:
            response = requests.put(
                f'{self.base_url}/projects/{project_id}/documents/{document_id}',
                headers=self.headers,
                json=update_data
            )
            
            if response.status_code == 409:
                # 处理版本冲突
                conflict = response.json()
                print(f"版本冲突: {conflict['error']['details']}")
                return None
            
            response.raise_for_status()
            return response.json()
            
        except requests.exceptions.RequestException as e:
            print(f"更新文档失败: {e}")
            raise

# 使用示例
api = DocumentAPI('https://api.qingyu.com/api/v1', 'your_token_here')

# 创建项目
project = api.create_project({
    'name': '我的小说',
    'status': 'private',
    'description': '一个精彩的故事'
})

print(f"项目创建成功: {project['data']['id']}")
```

这些示例展示了如何正确使用青羽文档管理系统的API，包括错误处理、版本控制和认证等关键功能。