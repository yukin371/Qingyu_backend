# 写作系统 API 参考

> **版本**: v1.0  
> **最后更新**: 2025-10-20  
> **基础路径**: `/api/v1`

---

## 1. 概述

写作系统为作者提供完整的创作工具集，包括项目管理、文档编辑、版本控制、内容审核、数据统计等功能。

### 1.1 基础信息

- **认证要求**: 所有接口均需登录认证
- **响应格式**: 统一 JSON 格式
- **分页支持**: 是
- **Token**: 在请求头中添加 `Authorization: Bearer <token>`

### 1.2 功能特性

- ✅ 项目管理（创建、组织、统计）
- ✅ 文档树形结构管理
- ✅ 编辑器自动保存（版本冲突检测）
- ✅ 字数统计（支持 Markdown 过滤）
- ✅ 快捷键自定义
- ✅ 内容审核（实时检测、全文审核、申诉）
- ✅ 数据统计（阅读、收入、读者行为）
- ✅ 版本历史（对比、恢复）

---

## 2. 项目管理 API

### 2.1 创建项目

**接口说明**: 创建一个新的写作项目

**请求**
```
POST /api/v1/projects
```

**请求头**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
```json
{
  "title": "我的小说",
  "description": "这是一部玄幻小说",
  "category": "玄幻",
  "type": "novel",
  "settings": {
    "targetWordCount": 1000000,
    "dailyGoal": 3000
  }
}
```

**请求参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 项目标题（1-200字） |
| description | string | 否 | 项目简介 |
| category | string | 否 | 项目分类 |
| type | string | 否 | 项目类型（novel/essay/others） |
| settings | object | 否 | 项目设置 |

**响应示例**
```json
{
  "code": 201,
  "message": "创建成功",
  "data": {
    "projectId": "proj_123456",
    "title": "我的小说",
    "createdAt": "2025-10-20T10:00:00Z"
  }
}
```

**错误响应**
- `400` - 参数错误（标题为空或过长）
- `401` - 未认证
- `500` - 服务器错误

**JavaScript/Axios 示例**
```javascript
const createProject = async (projectData) => {
  try {
    const response = await axios.post('/api/v1/projects', {
      title: projectData.title,
      description: projectData.description,
      category: projectData.category,
      type: 'novel',
      settings: {
        targetWordCount: 1000000,
        dailyGoal: 3000
      }
    });
    return response.data.data;
  } catch (error) {
    console.error('创建项目失败:', error);
    throw error;
  }
};
```

---

### 2.2 获取项目列表

**接口说明**: 获取当前用户的项目列表，支持分页和筛选

**请求**
```
GET /api/v1/projects
```

**Query 参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| pageSize | int | 否 | 10 | 每页数量 |
| status | string | 否 | - | 项目状态（draft/writing/completed） |
| category | string | 否 | - | 项目分类 |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "projects": [
      {
        "projectId": "proj_123456",
        "title": "我的小说",
        "category": "玄幻",
        "wordCount": 50000,
        "chapterCount": 10,
        "status": "writing",
        "updatedAt": "2025-10-20T09:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "pageSize": 10
  }
}
```

**JavaScript/Axios 示例**
```javascript
const getProjects = async (page = 1, pageSize = 10, filters = {}) => {
  try {
    const response = await axios.get('/api/v1/projects', {
      params: {
        page,
        pageSize,
        ...filters
      }
    });
    return response.data.data;
  } catch (error) {
    console.error('获取项目列表失败:', error);
    throw error;
  }
};
```

---

### 2.3 获取项目详情

**接口说明**: 根据项目ID获取项目的详细信息

**请求**
```
GET /api/v1/projects/{id}
```

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 项目ID |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "projectId": "proj_123456",
    "title": "我的小说",
    "description": "这是一部玄幻小说",
    "category": "玄幻",
    "type": "novel",
    "wordCount": 50000,
    "chapterCount": 10,
    "status": "writing",
    "settings": {
      "targetWordCount": 1000000,
      "dailyGoal": 3000
    },
    "createdAt": "2025-10-01T00:00:00Z",
    "updatedAt": "2025-10-20T09:00:00Z"
  }
}
```

---

### 2.4 更新项目

**接口说明**: 更新项目信息

**请求**
```
PUT /api/v1/projects/{id}
```

**请求体**
```json
{
  "title": "修改后的标题",
  "description": "修改后的简介",
  "status": "completed"
}
```

**响应示例**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

---

### 2.5 删除项目

**接口说明**: 删除项目（软删除）

**请求**
```
DELETE /api/v1/projects/{id}
```

**响应示例**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

---

### 2.6 更新项目统计信息

**接口说明**: 重新计算并更新项目的统计信息（字数、章节数等）

**请求**
```
PUT /api/v1/projects/{id}/statistics
```

**响应示例**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

---

## 3. 文档管理 API

### 3.1 创建文档

**接口说明**: 在项目中创建新文档（章节）

**请求**
```
POST /api/v1/projects/{projectId}/documents
```

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| projectId | string | 是 | 项目ID |

**请求体**
```json
{
  "title": "第一章 开篇",
  "content": "章节内容...",
  "parentId": "",
  "type": "chapter"
}
```

**请求参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 文档标题 |
| content | string | 否 | 文档内容 |
| parentId | string | 否 | 父文档ID（用于创建子章节） |
| type | string | 否 | 文档类型（chapter/folder） |

**响应示例**
```json
{
  "code": 201,
  "message": "创建成功",
  "data": {
    "documentId": "doc_123456",
    "title": "第一章 开篇",
    "createdAt": "2025-10-20T10:00:00Z"
  }
}
```

**JavaScript/Axios 示例**
```javascript
const createDocument = async (projectId, docData) => {
  try {
    const response = await axios.post(`/api/v1/projects/${projectId}/documents`, {
      title: docData.title,
      content: docData.content || '',
      parentId: docData.parentId || '',
      type: 'chapter'
    });
    return response.data.data;
  } catch (error) {
    console.error('创建文档失败:', error);
    throw error;
  }
};
```

---

### 3.2 获取文档列表

**接口说明**: 获取项目下的文档列表（支持分页）

**请求**
```
GET /api/v1/projects/{projectId}/documents
```

**Query 参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| pageSize | int | 否 | 20 | 每页数量 |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "documents": [
      {
        "documentId": "doc_123456",
        "title": "第一章 开篇",
        "wordCount": 3000,
        "updatedAt": "2025-10-20T09:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "pageSize": 20
  }
}
```

---

### 3.3 获取文档树

**接口说明**: 获取项目的文档树形结构

**请求**
```
GET /api/v1/projects/{projectId}/documents/tree
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "tree": [
      {
        "documentId": "doc_123",
        "title": "卷一",
        "type": "folder",
        "children": [
          {
            "documentId": "doc_456",
            "title": "第一章",
            "type": "chapter",
            "children": []
          }
        ]
      }
    ]
  }
}
```

**JavaScript/Axios 示例**
```javascript
const getDocumentTree = async (projectId) => {
  try {
    const response = await axios.get(`/api/v1/projects/${projectId}/documents/tree`);
    return response.data.data.tree;
  } catch (error) {
    console.error('获取文档树失败:', error);
    throw error;
  }
};
```

---

### 3.4 获取文档详情

**接口说明**: 根据文档ID获取文档详细信息

**请求**
```
GET /api/v1/documents/{id}
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "documentId": "doc_123456",
    "projectId": "proj_123",
    "title": "第一章 开篇",
    "content": "章节内容...",
    "wordCount": 3000,
    "status": "published",
    "createdAt": "2025-10-01T00:00:00Z",
    "updatedAt": "2025-10-20T09:00:00Z"
  }
}
```

---

### 3.5 更新文档

**接口说明**: 更新文档信息

**请求**
```
PUT /api/v1/documents/{id}
```

**请求体**
```json
{
  "title": "修改后的标题",
  "content": "修改后的内容",
  "status": "published"
}
```

**响应示例**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

---

### 3.6 删除文档

**接口说明**: 删除文档（软删除）

**请求**
```
DELETE /api/v1/documents/{id}
```

**响应示例**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

---

### 3.7 移动文档

**接口说明**: 移动文档到新的父节点

**请求**
```
PUT /api/v1/documents/{id}/move
```

**请求体**
```json
{
  "newParentId": "doc_parent_456",
  "newOrder": 1
}
```

**响应示例**
```json
{
  "code": 200,
  "message": "移动成功"
}
```

---

### 3.8 重新排序文档

**接口说明**: 批量更新同级文档的顺序

**请求**
```
PUT /api/v1/projects/{projectId}/documents/reorder
```

**请求体**
```json
{
  "orders": [
    { "documentId": "doc_1", "order": 1 },
    { "documentId": "doc_2", "order": 2 },
    { "documentId": "doc_3", "order": 3 }
  ]
}
```

**响应示例**
```json
{
  "code": 200,
  "message": "排序成功"
}
```

---

## 4. 编辑器 API

### 4.1 自动保存文档

**接口说明**: 自动保存文档内容，支持版本冲突检测

**请求**
```
POST /api/v1/documents/{id}/autosave
```

**请求体**
```json
{
  "content": "最新的文档内容...",
  "version": 5,
  "wordCount": 3500
}
```

**请求参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content | string | 是 | 文档内容 |
| version | int | 是 | 当前版本号（用于冲突检测） |
| wordCount | int | 否 | 字数统计 |

**响应示例**
```json
{
  "code": 200,
  "message": "保存成功",
  "data": {
    "version": 6,
    "savedAt": "2025-10-20T10:00:00Z",
    "wordCount": 3500
  }
}
```

**错误响应**
- `409` - 版本冲突（文档已被其他用户修改）

```json
{
  "code": 409,
  "message": "版本冲突",
  "error": "文档已被其他用户修改，请刷新后重试"
}
```

**JavaScript/Axios 示例**
```javascript
// 自动保存钩子（每30秒自动保存）
const useAutoSave = (documentId, content, version) => {
  useEffect(() => {
    const timer = setInterval(async () => {
      try {
        await axios.post(`/api/v1/documents/${documentId}/autosave`, {
          content,
          version,
          wordCount: content.length
        });
        console.log('自动保存成功');
      } catch (error) {
        if (error.response?.status === 409) {
          alert('文档已被修改，请刷新页面');
        }
      }
    }, 30000); // 30秒

    return () => clearInterval(timer);
  }, [documentId, content, version]);
};
```

---

### 4.2 获取保存状态

**接口说明**: 获取文档的保存状态和最后保存时间

**请求**
```
GET /api/v1/documents/{id}/save-status
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "isSaved": true,
    "lastSavedAt": "2025-10-20T10:00:00Z",
    "version": 6,
    "hasUnsavedChanges": false
  }
}
```

---

### 4.3 获取文档内容

**接口说明**: 获取文档的完整内容（用于编辑器加载）

**请求**
```
GET /api/v1/documents/{id}/content
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "documentId": "doc_123456",
    "content": "文档内容...",
    "version": 6,
    "wordCount": 3500,
    "updatedAt": "2025-10-20T10:00:00Z"
  }
}
```

---

### 4.4 更新文档内容

**接口说明**: 手动更新文档内容（非自动保存）

**请求**
```
PUT /api/v1/documents/{id}/content
```

**请求体**
```json
{
  "content": "更新后的内容...",
  "wordCount": 3600
}
```

**响应示例**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

---

### 4.5 计算字数

**接口说明**: 计算文档内容的字数统计（支持 Markdown 过滤）

**请求**
```
POST /api/v1/documents/{id}/word-count
```

**请求体**
```json
{
  "content": "文档内容...",
  "filterMarkdown": true
}
```

**请求参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content | string | 是 | 文档内容 |
| filterMarkdown | bool | 否 | 是否过滤 Markdown 语法（默认 false） |

**响应示例**
```json
{
  "code": 200,
  "message": "计算成功",
  "data": {
    "totalWords": 3500,
    "totalCharacters": 10500,
    "totalLines": 120,
    "totalParagraphs": 50
  }
}
```

**JavaScript/Axios 示例**
```javascript
const calculateWordCount = async (documentId, content, filterMarkdown = true) => {
  try {
    const response = await axios.post(`/api/v1/documents/${documentId}/word-count`, {
      content,
      filterMarkdown
    });
    return response.data.data;
  } catch (error) {
    console.error('计算字数失败:', error);
    throw error;
  }
};
```

---

### 4.6 获取用户快捷键配置

**接口说明**: 获取当前用户的快捷键配置（包括自定义和默认）

**请求**
```
GET /api/v1/user/shortcuts
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "shortcuts": {
      "save": "Ctrl+S",
      "bold": "Ctrl+B",
      "italic": "Ctrl+I",
      "undo": "Ctrl+Z",
      "redo": "Ctrl+Y"
    }
  }
}
```

---

### 4.7 更新用户快捷键配置

**接口说明**: 更新用户的自定义快捷键配置

**请求**
```
PUT /api/v1/user/shortcuts
```

**请求体**
```json
{
  "shortcuts": {
    "save": "Ctrl+Alt+S",
    "bold": "Ctrl+B"
  }
}
```

**响应示例**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

---

### 4.8 重置用户快捷键配置

**接口说明**: 重置用户快捷键为默认配置

**请求**
```
POST /api/v1/user/shortcuts/reset
```

**响应示例**
```json
{
  "code": 200,
  "message": "重置成功"
}
```

---

### 4.9 获取快捷键帮助

**接口说明**: 获取快捷键帮助文档（按分类）

**请求**
```
GET /api/v1/user/shortcuts/help
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "category": "编辑",
      "shortcuts": [
        { "name": "保存", "key": "Ctrl+S" },
        { "name": "撤销", "key": "Ctrl+Z" }
      ]
    },
    {
      "category": "格式",
      "shortcuts": [
        { "name": "加粗", "key": "Ctrl+B" },
        { "name": "斜体", "key": "Ctrl+I" }
      ]
    }
  ]
}
```

---

## 5. 内容审核 API

### 5.1 实时检测内容

**接口说明**: 快速检测内容是否包含违规信息（不创建审核记录）

**请求**
```
POST /api/v1/audit/check
```

**请求体**
```json
{
  "content": "要检测的文本内容..."
}
```

**响应示例**
```json
{
  "code": 200,
  "message": "检测完成",
  "data": {
    "passed": false,
    "riskLevel": 3,
    "violations": [
      {
        "type": "sensitive_word",
        "content": "违规词",
        "position": 10,
        "suggestion": "建议修改为..."
      }
    ]
  }
}
```

**JavaScript/Axios 示例**
```javascript
// 实时内容检测（输入时防抖）
const useContentCheck = (content) => {
  const [checkResult, setCheckResult] = useState(null);
  
  useEffect(() => {
    const timer = setTimeout(async () => {
      if (content.length > 0) {
        try {
          const response = await axios.post('/api/v1/audit/check', { content });
          setCheckResult(response.data.data);
        } catch (error) {
          console.error('内容检测失败:', error);
        }
      }
    }, 500); // 500ms 防抖
    
    return () => clearTimeout(timer);
  }, [content]);
  
  return checkResult;
};
```

---

### 5.2 全文审核文档

**接口说明**: 对文档进行全文审核并创建审核记录

**请求**
```
POST /api/v1/documents/{id}/audit
```

**请求体**
```json
{
  "content": "文档全文内容..."
}
```

**响应示例**
```json
{
  "code": 200,
  "message": "审核完成",
  "data": {
    "auditId": "audit_123456",
    "status": "rejected",
    "riskLevel": 4,
    "violations": [
      {
        "type": "政治敏感",
        "count": 2,
        "positions": [10, 25]
      }
    ],
    "createdAt": "2025-10-20T10:00:00Z"
  }
}
```

---

### 5.3 获取审核结果

**接口说明**: 根据文档ID获取审核结果

**请求**
```
GET /api/v1/documents/{id}/audit-result
```

**Query 参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| targetType | string | 否 | document | 目标类型 |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "auditId": "audit_123456",
    "status": "pending",
    "riskLevel": 2,
    "reviewedAt": "2025-10-20T10:00:00Z"
  }
}
```

---

### 5.4 提交申诉

**接口说明**: 对审核结果提交申诉

**请求**
```
POST /api/v1/audit/{id}/appeal
```

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 审核记录ID |

**请求体**
```json
{
  "reason": "申诉理由..."
}
```

**响应示例**
```json
{
  "code": 200,
  "message": "申诉已提交，等待复核"
}
```

---

### 5.5 获取用户违规记录

**接口说明**: 获取指定用户的所有违规记录

**请求**
```
GET /api/v1/users/{userId}/violations
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "violationId": "vio_123",
      "type": "敏感词",
      "content": "违规内容",
      "status": "confirmed",
      "createdAt": "2025-10-15T00:00:00Z"
    }
  ]
}
```

---

### 5.6 获取用户违规统计

**接口说明**: 获取指定用户的违规统计信息

**请求**
```
GET /api/v1/users/{userId}/violation-summary
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "totalViolations": 5,
    "violationsByType": {
      "敏感词": 3,
      "政治敏感": 2
    },
    "lastViolationAt": "2025-10-15T00:00:00Z"
  }
}
```

---

## 6. 数据统计 API

### 6.1 获取作品统计数据

**接口说明**: 获取作品的完整统计信息，包括阅读、收入、互动等数据

**请求**
```
GET /api/v1/writer/books/{book_id}/stats
```

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| book_id | string | 是 | 作品ID |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "bookId": "book_123",
    "totalViews": 100000,
    "totalLikes": 5000,
    "totalComments": 1200,
    "totalRevenue": 15000.50,
    "averageReadingTime": 1800,
    "completionRate": 0.75,
    "updatedAt": "2025-10-20T10:00:00Z"
  }
}
```

**JavaScript/Axios 示例**
```javascript
const getBookStats = async (bookId) => {
  try {
    const response = await axios.get(`/api/v1/writer/books/${bookId}/stats`);
    return response.data.data;
  } catch (error) {
    console.error('获取作品统计失败:', error);
    throw error;
  }
};
```

---

### 6.2 获取章节统计数据

**接口说明**: 获取单个章节的统计信息

**请求**
```
GET /api/v1/writer/chapters/{chapter_id}/stats
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "chapterId": "chapter_123",
    "viewCount": 5000,
    "likeCount": 200,
    "commentCount": 50,
    "completionRate": 0.80,
    "averageReadingTime": 600,
    "dropOffRate": 0.15
  }
}
```

---

### 6.3 获取阅读热力图

**接口说明**: 获取作品各章节的阅读热度分布

**请求**
```
GET /api/v1/writer/books/{book_id}/heatmap
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "chapterId": "chapter_1",
      "chapterNumber": 1,
      "heat": 0.95
    },
    {
      "chapterId": "chapter_2",
      "chapterNumber": 2,
      "heat": 0.85
    }
  ]
}
```

**JavaScript/Axios 示例**
```javascript
// 使用热力图数据渲染可视化图表
const renderHeatmap = async (bookId) => {
  try {
    const response = await axios.get(`/api/v1/writer/books/${bookId}/heatmap`);
    const heatmapData = response.data.data;
    
    // 使用 ECharts 或其他图表库渲染
    const option = {
      series: [{
        type: 'heatmap',
        data: heatmapData.map(item => [item.chapterNumber, item.heat])
      }]
    };
    
    return heatmapData;
  } catch (error) {
    console.error('获取热力图失败:', error);
    throw error;
  }
};
```

---

### 6.4 获取收入统计

**接口说明**: 获取作品的收入细分数据

**请求**
```
GET /api/v1/writer/books/{book_id}/revenue
```

**Query 参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| start_date | string | 否 | 最近30天 | 开始日期（YYYY-MM-DD） |
| end_date | string | 否 | 今天 | 结束日期（YYYY-MM-DD） |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "totalRevenue": 15000.50,
    "subscriptionRevenue": 10000.00,
    "rewardRevenue": 3000.50,
    "adRevenue": 2000.00,
    "breakdown": [
      {
        "date": "2025-10-15",
        "revenue": 500.00
      }
    ]
  }
}
```

---

### 6.5 获取热门章节

**接口说明**: 获取作品的热门章节统计（阅读量最高、收入最高、完读率最低、跳出率最高）

**请求**
```
GET /api/v1/writer/books/{book_id}/top-chapters
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "mostViewed": [
      {
        "chapterId": "chapter_1",
        "title": "第一章",
        "viewCount": 10000
      }
    ],
    "highestRevenue": [...],
    "lowestCompletion": [...],
    "highestDropOff": [...]
  }
}
```

---

### 6.6 获取每日统计

**接口说明**: 获取作品最近N天的每日统计数据

**请求**
```
GET /api/v1/writer/books/{book_id}/daily-stats
```

**Query 参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| days | int | 否 | 7 | 天数（1-365） |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "date": "2025-10-20",
      "views": 1000,
      "likes": 50,
      "revenue": 200.00
    },
    {
      "date": "2025-10-19",
      "views": 950,
      "likes": 45,
      "revenue": 180.00
    }
  ]
}
```

**JavaScript/Axios 示例**
```javascript
// 获取最近7天的数据并渲染折线图
const getDailyStatsChart = async (bookId, days = 7) => {
  try {
    const response = await axios.get(`/api/v1/writer/books/${bookId}/daily-stats`, {
      params: { days }
    });
    
    const data = response.data.data;
    
    // 准备图表数据
    const chartData = {
      dates: data.map(item => item.date),
      views: data.map(item => item.views),
      likes: data.map(item => item.likes),
      revenue: data.map(item => item.revenue)
    };
    
    return chartData;
  } catch (error) {
    console.error('获取每日统计失败:', error);
    throw error;
  }
};
```

---

### 6.7 获取跳出点分析

**接口说明**: 获取跳出率最高的章节列表

**请求**
```
GET /api/v1/writer/books/{book_id}/drop-off-points
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "chapterId": "chapter_5",
      "title": "第五章",
      "dropOffRate": 0.35,
      "viewCount": 5000,
      "completedCount": 3250
    }
  ]
}
```

---

### 6.8 获取留存率

**接口说明**: 获取作品的N日留存率

**请求**
```
GET /api/v1/writer/books/{book_id}/retention
```

**Query 参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| days | int | 否 | 7 | 天数（1-90） |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "book_id": "book_123",
    "days": 7,
    "retention_rate": 0.45
  }
}
```

---

## 7. 版本管理 API

### 7.1 获取版本历史

**接口说明**: 获取文档的版本历史列表

**请求**
```
GET /api/v1/documents/{documentId}/versions
```

**Query 参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| pageSize | int | 否 | 20 | 每页数量 |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "versions": [
      {
        "versionId": "v_123",
        "versionNumber": 10,
        "createdAt": "2025-10-20T10:00:00Z",
        "comment": "修改第三段内容",
        "wordCount": 3500
      }
    ],
    "total": 10,
    "page": 1,
    "pageSize": 20
  }
}
```

**JavaScript/Axios 示例**
```javascript
const getVersionHistory = async (documentId, page = 1, pageSize = 20) => {
  try {
    const response = await axios.get(`/api/v1/documents/${documentId}/versions`, {
      params: { page, pageSize }
    });
    return response.data.data;
  } catch (error) {
    console.error('获取版本历史失败:', error);
    throw error;
  }
};
```

---

### 7.2 获取特定版本

**接口说明**: 获取文档的特定版本内容

**请求**
```
GET /api/v1/documents/{documentId}/versions/{versionId}
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "versionId": "v_123",
    "versionNumber": 10,
    "content": "版本10的文档内容...",
    "wordCount": 3500,
    "createdAt": "2025-10-20T10:00:00Z"
  }
}
```

---

### 7.3 比较版本

**接口说明**: 比较两个版本的差异

**请求**
```
GET /api/v1/documents/{documentId}/versions/compare
```

**Query 参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| fromVersion | string | 是 | 源版本ID |
| toVersion | string | 是 | 目标版本ID |

**响应示例**
```json
{
  "code": 200,
  "message": "比较成功",
  "data": {
    "fromVersion": "v_8",
    "toVersion": "v_10",
    "diff": [
      {
        "type": "add",
        "position": 100,
        "content": "新增的内容"
      },
      {
        "type": "delete",
        "position": 200,
        "content": "删除的内容"
      },
      {
        "type": "modify",
        "position": 300,
        "oldContent": "旧内容",
        "newContent": "新内容"
      }
    ]
  }
}
```

**JavaScript/Axios 示例**
```javascript
// 版本对比并高亮显示差异
const compareVersions = async (documentId, fromVersion, toVersion) => {
  try {
    const response = await axios.get(
      `/api/v1/documents/${documentId}/versions/compare`,
      {
        params: { fromVersion, toVersion }
      }
    );
    
    const diff = response.data.data.diff;
    
    // 渲染差异（使用 diff 库或自定义高亮）
    return diff;
  } catch (error) {
    console.error('版本对比失败:', error);
    throw error;
  }
};
```

---

### 7.4 恢复版本

**接口说明**: 将文档恢复到特定版本

**请求**
```
POST /api/v1/documents/{documentId}/versions/{versionId}/restore
```

**响应示例**
```json
{
  "code": 200,
  "message": "恢复成功"
}
```

**JavaScript/Axios 示例**
```javascript
const restoreVersion = async (documentId, versionId) => {
  try {
    await axios.post(
      `/api/v1/documents/${documentId}/versions/${versionId}/restore`
    );
    
    alert('版本恢复成功');
    // 刷新文档内容
    window.location.reload();
  } catch (error) {
    console.error('恢复版本失败:', error);
    throw error;
  }
};
```

---

## 8. 完整示例代码

### 8.1 Vue 3 写作项目管理组件

```vue
<template>
  <div class="writer-dashboard">
    <!-- 项目列表 -->
    <div class="project-list">
      <h2>我的项目</h2>
      <button @click="showCreateProject">创建新项目</button>
      
      <div v-for="project in projects" :key="project.projectId" class="project-item">
        <h3>{{ project.title }}</h3>
        <p>字数: {{ project.wordCount }} | 章节数: {{ project.chapterCount }}</p>
        <button @click="openProject(project.projectId)">打开</button>
        <button @click="deleteProject(project.projectId)">删除</button>
      </div>
    </div>
    
    <!-- 创建项目对话框 -->
    <el-dialog v-model="createDialogVisible" title="创建新项目">
      <el-form :model="newProject" label-width="100px">
        <el-form-item label="项目标题">
          <el-input v-model="newProject.title" placeholder="请输入项目标题" />
        </el-form-item>
        <el-form-item label="项目简介">
          <el-input type="textarea" v-model="newProject.description" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="newProject.category">
            <el-option label="玄幻" value="玄幻" />
            <el-option label="都市" value="都市" />
            <el-option label="历史" value="历史" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="createProject">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { ElMessage } from 'element-plus';

const projects = ref([]);
const createDialogVisible = ref(false);
const newProject = ref({
  title: '',
  description: '',
  category: '玄幻'
});

// 获取项目列表
const fetchProjects = async () => {
  try {
    const response = await axios.get('/api/v1/projects', {
      params: { page: 1, pageSize: 20 }
    });
    projects.value = response.data.data.projects;
  } catch (error) {
    ElMessage.error('获取项目列表失败');
  }
};

// 创建项目
const createProject = async () => {
  try {
    await axios.post('/api/v1/projects', {
      title: newProject.value.title,
      description: newProject.value.description,
      category: newProject.value.category,
      type: 'novel'
    });
    
    ElMessage.success('创建成功');
    createDialogVisible.value = false;
    fetchProjects();
  } catch (error) {
    ElMessage.error('创建失败');
  }
};

// 删除项目
const deleteProject = async (projectId) => {
  try {
    await ElMessageBox.confirm('确定删除该项目吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await axios.delete(`/api/v1/projects/${projectId}`);
    ElMessage.success('删除成功');
    fetchProjects();
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败');
    }
  }
};

const showCreateProject = () => {
  createDialogVisible.value = true;
};

const openProject = (projectId) => {
  // 跳转到编辑器页面
  window.location.href = `/editor/${projectId}`;
};

onMounted(() => {
  fetchProjects();
});
</script>
```

---

### 8.2 Vue 3 编辑器组件（带自动保存）

```vue
<template>
  <div class="editor-container">
    <div class="editor-header">
      <h2>{{ document.title }}</h2>
      <span class="save-status">
        {{ saveStatus }}
      </span>
    </div>
    
    <div class="editor-body">
      <textarea
        v-model="content"
        @input="handleContentChange"
        placeholder="开始写作..."
      />
    </div>
    
    <div class="editor-footer">
      <span>字数: {{ wordCount }}</span>
      <button @click="manualSave">手动保存</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import axios from 'axios';
import { ElMessage } from 'element-plus';
import { useRoute } from 'vue-router';

const route = useRoute();
const documentId = route.params.id;

const document = ref({});
const content = ref('');
const version = ref(0);
const saveStatus = ref('已保存');
const autoSaveTimer = ref(null);

const wordCount = computed(() => content.value.length);

// 加载文档内容
const loadDocument = async () => {
  try {
    const response = await axios.get(`/api/v1/documents/${documentId}/content`);
    const data = response.data.data;
    
    document.value = data;
    content.value = data.content;
    version.value = data.version;
  } catch (error) {
    ElMessage.error('加载文档失败');
  }
};

// 自动保存
const autoSave = async () => {
  try {
    saveStatus.value = '保存中...';
    
    const response = await axios.post(`/api/v1/documents/${documentId}/autosave`, {
      content: content.value,
      version: version.value,
      wordCount: wordCount.value
    });
    
    version.value = response.data.data.version;
    saveStatus.value = '已保存';
  } catch (error) {
    if (error.response?.status === 409) {
      ElMessage.error('版本冲突，请刷新页面');
      saveStatus.value = '保存失败';
    } else {
      ElMessage.error('自动保存失败');
      saveStatus.value = '保存失败';
    }
  }
};

// 手动保存
const manualSave = async () => {
  await autoSave();
  ElMessage.success('保存成功');
};

// 内容变化时标记为未保存
const handleContentChange = () => {
  saveStatus.value = '未保存';
  
  // 重置自动保存计时器
  if (autoSaveTimer.value) {
    clearTimeout(autoSaveTimer.value);
  }
  
  autoSaveTimer.value = setTimeout(() => {
    autoSave();
  }, 30000); // 30秒后自动保存
};

// 页面卸载前保存
onBeforeUnmount(() => {
  if (saveStatus.value === '未保存') {
    autoSave();
  }
  
  if (autoSaveTimer.value) {
    clearTimeout(autoSaveTimer.value);
  }
});

onMounted(() => {
  loadDocument();
});
</script>

<style scoped>
.editor-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.editor-header {
  padding: 1rem;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.save-status {
  color: #666;
  font-size: 0.9rem;
}

.editor-body {
  flex: 1;
  padding: 1rem;
}

.editor-body textarea {
  width: 100%;
  height: 100%;
  border: none;
  outline: none;
  font-size: 1rem;
  line-height: 1.8;
  resize: none;
}

.editor-footer {
  padding: 1rem;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
```

---

### 8.3 数据统计仪表板组件

```vue
<template>
  <div class="stats-dashboard">
    <h2>作品数据</h2>
    
    <!-- 概览卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <h3>总阅读量</h3>
        <p class="stat-value">{{ stats.totalViews?.toLocaleString() }}</p>
      </div>
      <div class="stat-card">
        <h3>总收入</h3>
        <p class="stat-value">¥{{ stats.totalRevenue?.toFixed(2) }}</p>
      </div>
      <div class="stat-card">
        <h3>总点赞</h3>
        <p class="stat-value">{{ stats.totalLikes?.toLocaleString() }}</p>
      </div>
      <div class="stat-card">
        <h3>完读率</h3>
        <p class="stat-value">{{ (stats.completionRate * 100)?.toFixed(1) }}%</p>
      </div>
    </div>
    
    <!-- 每日趋势图 -->
    <div class="chart-container">
      <h3>最近7天趋势</h3>
      <div ref="dailyChart" style="height: 300px;"></div>
    </div>
    
    <!-- 热力图 -->
    <div class="chart-container">
      <h3>章节热力图</h3>
      <div ref="heatmapChart" style="height: 400px;"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import * as echarts from 'echarts';
import { useRoute } from 'vue-router';

const route = useRoute();
const bookId = route.params.bookId;

const stats = ref({});
const dailyChart = ref(null);
const heatmapChart = ref(null);

// 获取作品统计
const fetchBookStats = async () => {
  try {
    const response = await axios.get(`/api/v1/writer/books/${bookId}/stats`);
    stats.value = response.data.data;
  } catch (error) {
    console.error('获取统计数据失败:', error);
  }
};

// 获取每日统计并渲染图表
const fetchDailyStats = async () => {
  try {
    const response = await axios.get(`/api/v1/writer/books/${bookId}/daily-stats`, {
      params: { days: 7 }
    });
    
    const data = response.data.data;
    
    const chart = echarts.init(dailyChart.value);
    const option = {
      tooltip: { trigger: 'axis' },
      legend: { data: ['阅读量', '点赞数', '收入'] },
      xAxis: {
        type: 'category',
        data: data.map(item => item.date)
      },
      yAxis: [
        { type: 'value', name: '阅读/点赞' },
        { type: 'value', name: '收入' }
      ],
      series: [
        {
          name: '阅读量',
          type: 'line',
          data: data.map(item => item.views)
        },
        {
          name: '点赞数',
          type: 'line',
          data: data.map(item => item.likes)
        },
        {
          name: '收入',
          type: 'bar',
          yAxisIndex: 1,
          data: data.map(item => item.revenue)
        }
      ]
    };
    
    chart.setOption(option);
  } catch (error) {
    console.error('获取每日统计失败:', error);
  }
};

// 获取热力图并渲染
const fetchHeatmap = async () => {
  try {
    const response = await axios.get(`/api/v1/writer/books/${bookId}/heatmap`);
    const data = response.data.data;
    
    const chart = echarts.init(heatmapChart.value);
    const option = {
      tooltip: {},
      xAxis: {
        type: 'category',
        data: data.map(item => `第${item.chapterNumber}章`)
      },
      yAxis: { type: 'value' },
      series: [{
        type: 'bar',
        data: data.map(item => item.heat),
        itemStyle: {
          color: (params) => {
            const heat = params.value;
            if (heat > 0.8) return '#67C23A';
            if (heat > 0.5) return '#E6A23C';
            return '#F56C6C';
          }
        }
      }]
    };
    
    chart.setOption(option);
  } catch (error) {
    console.error('获取热力图失败:', error);
  }
};

onMounted(() => {
  fetchBookStats();
  fetchDailyStats();
  fetchHeatmap();
});
</script>

<style scoped>
.stats-dashboard {
  padding: 2rem;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.stat-card h3 {
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 0.5rem;
}

.stat-value {
  font-size: 2rem;
  font-weight: bold;
  color: #333;
}

.chart-container {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  margin-bottom: 1.5rem;
}
</style>
```

---

## 9. 最佳实践

### 9.1 自动保存策略

```javascript
// 推荐的自动保存策略
const AUTO_SAVE_INTERVAL = 30000; // 30秒
const DEBOUNCE_DELAY = 500; // 500ms 防抖

// 使用防抖 + 定时器组合
let autoSaveTimer = null;
let debounceTimer = null;

const scheduleAutoSave = () => {
  // 清除旧的定时器
  if (autoSaveTimer) clearTimeout(autoSaveTimer);
  
  // 设置新的定时器
  autoSaveTimer = setTimeout(() => {
    performAutoSave();
  }, AUTO_SAVE_INTERVAL);
};

const handleContentChange = () => {
  // 防抖：用户停止输入500ms后才标记为未保存
  if (debounceTimer) clearTimeout(debounceTimer);
  
  debounceTimer = setTimeout(() => {
    saveStatus.value = '未保存';
    scheduleAutoSave();
  }, DEBOUNCE_DELAY);
};
```

---

### 9.2 版本冲突处理

```javascript
// 处理版本冲突的最佳实践
const handleVersionConflict = async (localContent, localVersion) => {
  try {
    // 1. 获取服务器最新版本
    const response = await axios.get(`/api/v1/documents/${documentId}/content`);
    const serverContent = response.data.data.content;
    const serverVersion = response.data.data.version;
    
    // 2. 提示用户选择
    const choice = await showConflictDialog({
      local: localContent,
      server: serverContent,
      localVersion,
      serverVersion
    });
    
    if (choice === 'useLocal') {
      // 强制使用本地版本
      await axios.put(`/api/v1/documents/${documentId}/content`, {
        content: localContent,
        force: true
      });
    } else if (choice === 'useServer') {
      // 使用服务器版本
      content.value = serverContent;
      version.value = serverVersion;
    } else if (choice === 'merge') {
      // 手动合并（显示对比界面）
      showMergeDialog(localContent, serverContent);
    }
  } catch (error) {
    console.error('处理版本冲突失败:', error);
  }
};
```

---

### 9.3 内容审核集成

```javascript
// 实时内容检测（输入时）
const useContentAudit = (content) => {
  const [violations, setViolations] = useState([]);
  
  useEffect(() => {
    const timer = setTimeout(async () => {
      if (content.length > 100) { // 至少100字才检测
        try {
          const response = await axios.post('/api/v1/audit/check', { content });
          setViolations(response.data.data.violations || []);
        } catch (error) {
          console.error('内容检测失败:', error);
        }
      }
    }, 1000); // 1秒防抖
    
    return () => clearTimeout(timer);
  }, [content]);
  
  return violations;
};

// 发布前全文审核
const auditBeforePublish = async (documentId, content) => {
  try {
    const response = await axios.post(`/api/v1/documents/${documentId}/audit`, {
      content
    });
    
    const result = response.data.data;
    
    if (result.status === 'rejected') {
      // 显示违规详情
      showViolationDialog(result.violations);
      return false;
    }
    
    return true;
  } catch (error) {
    console.error('审核失败:', error);
    return false;
  }
};
```

---

## 10. 常见问题

### Q1: 自动保存失败怎么办？

**A**: 
1. 检查网络连接
2. 查看浏览器控制台错误信息
3. 如果是版本冲突（409错误），刷新页面获取最新版本
4. 如果持续失败，手动复制内容到本地

### Q2: 如何避免版本冲突？

**A**:
1. 不要在多个设备同时编辑同一文档
2. 使用实时协作功能（如果支持）
3. 编辑前先刷新页面获取最新版本
4. 定期手动保存

### Q3: 内容被审核驳回怎么办？

**A**:
1. 查看具体违规原因
2. 修改违规内容后重新提交
3. 如果认为审核有误，可以提交申诉
4. 等待人工复核结果

### Q4: 如何提高写作效率？

**A**:
1. 使用快捷键（可自定义）
2. 设置每日字数目标
3. 查看数据统计，了解读者喜好
4. 使用版本管理，大胆修改不怕丢失

### Q5: 数据统计多久更新一次？

**A**:
- 实时统计: 即时更新
- 每日统计: 每天凌晨1点更新
- 热力图: 每小时更新
- 收入统计: 每天更新

---

## 11. 相关文档

- [API 快速参考](./API快速参考.md) - 快速查找接口
- [前端集成指南](./前端集成指南.md) - 详细集成步骤
- [共享服务 API 参考](./共享服务API参考.md) - 钱包、存储等
- [用户系统 API 参考](./用户系统API参考.md) - 用户管理

---

**最后更新**: 2025-10-20  
**维护者**: 青羽后端团队  
**文档版本**: v1.0

