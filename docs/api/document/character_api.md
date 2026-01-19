# Character API 设计文档

## 1. 概述

本文档定义了青羽项目中角色管理相关的API接口设计，包括角色卡片管理、角色关系管理等功能。

## 2. 数据模型

### 2.1 Character（角色）
```json
{
  "id": "string",
  "projectId": "string",
  "name": "string",
  "alias": ["string"],
  "summary": "string",
  "traits": ["string"],
  "background": "string",
  "avatarUrl": "string",
  "personalityPrompt": "string",
  "speechPattern": "string",
  "currentState": "string",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### 2.2 CharacterRelation（角色关系）
```json
{
  "id": "string",
  "projectId": "string",
  "fromId": "string",
  "toId": "string",
  "type": "朋友|家庭|敌人|恋人|盟友|其他",
  "strength": 85,
  "notes": "string",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

## 3. API 接口设计

### 3.1 角色管理

#### 3.1.1 创建角色
```
POST /api/v1/projects/{projectId}/characters
```

**请求体：**
```json
{
  "name": "李明",
  "alias": ["小李", "阿明"],
  "summary": "主角，程序员出身的冒险者",
  "traits": ["聪明", "勇敢", "有点内向"],
  "background": "出生在普通家庭，大学学习计算机科学...",
  "avatarUrl": "https://example.com/avatar.jpg",
  "personalityPrompt": "性格内向但关键时刻很勇敢，说话简洁有力",
  "speechPattern": "倾向于使用技术术语，语调平稳",
  "currentState": "刚刚进入异世界，正在适应新环境"
}
```

**响应：**
```json
{
  "code": 201,
  "message": "success",
  "data": {
    "id": "char_123456",
    "projectId": "proj_123456",
    "name": "李明",
    "alias": ["小李", "阿明"],
    "summary": "主角，程序员出身的冒险者",
    "traits": ["聪明", "勇敢", "有点内向"],
    "background": "出生在普通家庭，大学学习计算机科学...",
    "avatarUrl": "https://example.com/avatar.jpg",
    "personalityPrompt": "性格内向但关键时刻很勇敢，说话简洁有力",
    "speechPattern": "倾向于使用技术术语，语调平稳",
    "currentState": "刚刚进入异世界，正在适应新环境",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "req_123456"
}
```

#### 3.1.2 获取角色列表
```
GET /api/v1/projects/{projectId}/characters?page=1&page_size=20&search=李明&traits=勇敢
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "char_123456",
        "projectId": "proj_123456",
        "name": "李明",
        "alias": ["小李", "阿明"],
        "summary": "主角，程序员出身的冒险者",
        "traits": ["聪明", "勇敢", "有点内向"],
        "background": "出生在普通家庭，大学学习计算机科学...",
        "avatarUrl": "https://example.com/avatar.jpg",
        "personalityPrompt": "性格内向但关键时刻很勇敢，说话简洁有力",
        "speechPattern": "倾向于使用技术术语，语调平稳",
        "currentState": "刚刚进入异世界，正在适应新环境",
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

#### 3.1.3 获取角色详情
```
GET /api/v1/projects/{projectId}/characters/{characterId}
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "char_123456",
    "projectId": "proj_123456",
    "name": "李明",
    "alias": ["小李", "阿明"],
    "summary": "主角，程序员出身的冒险者",
    "traits": ["聪明", "勇敢", "有点内向"],
    "background": "出生在普通家庭，大学学习计算机科学...",
    "avatarUrl": "https://example.com/avatar.jpg",
    "personalityPrompt": "性格内向但关键时刻很勇敢，说话简洁有力",
    "speechPattern": "倾向于使用技术术语，语调平稳",
    "currentState": "刚刚进入异世界，正在适应新环境",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "req_123456"
}
```

#### 3.1.4 更新角色
```
PUT /api/v1/projects/{projectId}/characters/{characterId}
```

**请求体：**
```json
{
  "name": "李明",
  "alias": ["小李", "阿明", "明哥"],
  "summary": "主角，程序员出身的冒险者，现已成为团队领袖",
  "traits": ["聪明", "勇敢", "有点内向", "领导力"],
  "background": "出生在普通家庭，大学学习计算机科学，后来穿越到异世界...",
  "currentState": "已经适应异世界生活，正在寻找回家的方法"
}
```

#### 3.1.5 删除角色
```
DELETE /api/v1/projects/{projectId}/characters/{characterId}
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "deleted": true
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "req_123456"
}
```

### 3.2 角色关系管理

#### 3.2.1 创建角色关系
```
POST /api/v1/projects/{projectId}/character-relations
```

**请求体：**
```json
{
  "fromId": "char_123456",
  "toId": "char_789012",
  "type": "朋友",
  "strength": 85,
  "notes": "在冒险中结识的好友，互相信任"
}
```

**响应：**
```json
{
  "code": 201,
  "message": "success",
  "data": {
    "id": "rel_123456",
    "projectId": "proj_123456",
    "fromId": "char_123456",
    "toId": "char_789012",
    "type": "朋友",
    "strength": 85,
    "notes": "在冒险中结识的好友，互相信任",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "req_123456"
}
```

#### 3.2.2 获取角色关系列表
```
GET /api/v1/projects/{projectId}/character-relations?characterId=char_123456&type=朋友&page=1&page_size=20
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "rel_123456",
        "projectId": "proj_123456",
        "fromId": "char_123456",
        "toId": "char_789012",
        "type": "朋友",
        "strength": 85,
        "notes": "在冒险中结识的好友，互相信任",
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z",
        "fromCharacter": {
          "id": "char_123456",
          "name": "李明",
          "avatarUrl": "https://example.com/avatar1.jpg"
        },
        "toCharacter": {
          "id": "char_789012",
          "name": "王小红",
          "avatarUrl": "https://example.com/avatar2.jpg"
        }
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

#### 3.2.3 获取角色关系图
```
GET /api/v1/projects/{projectId}/character-relations/graph?characterId=char_123456&depth=2
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "nodes": [
      {
        "id": "char_123456",
        "name": "李明",
        "avatarUrl": "https://example.com/avatar1.jpg",
        "traits": ["聪明", "勇敢", "有点内向"]
      },
      {
        "id": "char_789012",
        "name": "王小红",
        "avatarUrl": "https://example.com/avatar2.jpg",
        "traits": ["活泼", "善良", "魔法师"]
      }
    ],
    "edges": [
      {
        "id": "rel_123456",
        "from": "char_123456",
        "to": "char_789012",
        "type": "朋友",
        "strength": 85,
        "notes": "在冒险中结识的好友，互相信任"
      }
    ]
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "req_123456"
}
```

#### 3.2.4 更新角色关系
```
PUT /api/v1/projects/{projectId}/character-relations/{relationId}
```

**请求体：**
```json
{
  "type": "恋人",
  "strength": 95,
  "notes": "从朋友发展为恋人关系"
}
```

#### 3.2.5 删除角色关系
```
DELETE /api/v1/projects/{projectId}/character-relations/{relationId}
```

### 3.3 角色分析

#### 3.3.1 获取角色统计
```
GET /api/v1/projects/{projectId}/characters/{characterId}/stats
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "characterId": "char_123456",
    "relationCount": 5,
    "relationTypes": {
      "朋友": 3,
      "敌人": 1,
      "恋人": 1
    },
    "appearanceCount": 15,
    "mentionCount": 32,
    "lastAppearance": "2024-01-01T00:00:00Z",
    "documents": [
      {
        "documentId": "doc_123",
        "title": "第一章：相遇",
        "appearanceCount": 5
      }
    ]
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "req_123456"
}
```

#### 3.3.2 角色关系强度分析
```
GET /api/v1/projects/{projectId}/character-relations/analysis?type=strength
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "averageStrength": 75.5,
    "strongestRelations": [
      {
        "id": "rel_123456",
        "fromCharacter": "李明",
        "toCharacter": "王小红",
        "type": "恋人",
        "strength": 95
      }
    ],
    "weakestRelations": [
      {
        "id": "rel_789012",
        "fromCharacter": "李明",
        "toCharacter": "张三",
        "type": "敌人",
        "strength": 20
      }
    ],
    "relationTypeDistribution": {
      "朋友": 40,
      "敌人": 20,
      "恋人": 15,
      "家庭": 15,
      "盟友": 10
    }
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "req_123456"
}
```

## 4. 查询参数说明

### 4.1 角色列表查询参数
- `page`: 页码，默认1
- `page_size`: 每页数量，默认20，最大100
- `search`: 搜索关键词（匹配姓名、别名、简介）
- `traits`: 性格标签过滤（支持多个，逗号分隔）
- `sort`: 排序字段（name, createdAt, updatedAt）
- `order`: 排序方向（asc, desc）

### 4.2 角色关系查询参数
- `characterId`: 指定角色ID（查询该角色的所有关系）
- `type`: 关系类型过滤
- `minStrength`: 最小关系强度
- `maxStrength`: 最大关系强度

## 5. 错误处理

### 5.1 特定错误码
| 错误码 | HTTP状态码 | 描述 |
|--------|------------|------|
| 42201 | 422 | 角色名称已存在 |
| 42202 | 422 | 无效的关系类型 |
| 42203 | 422 | 关系强度超出范围(0-100) |
| 42204 | 422 | 不能创建自己与自己的关系 |
| 42205 | 422 | 关系已存在 |

### 5.2 错误响应示例
```json
{
  "code": 422,
  "message": "角色名称已存在",
  "error": {
    "type": "character_name_exists",
    "details": [
      {
        "field": "name",
        "message": "项目中已存在名为'李明'的角色"
      }
    ]
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "req_123456"
}
```

## 6. 业务规则

### 6.1 角色管理规则
- 同一项目内角色名称不能重复
- 角色别名可以重复，但建议避免
- 删除角色时会同时删除相关的所有关系
- 角色的AI相关字段用于内容生成时的参考

### 6.2 关系管理规则
- 关系强度范围：0-100
- 支持单向关系，也可以创建双向关系
- 删除角色时自动删除相关关系
- 关系类型固定为预定义的几种类型

## 7. 扩展功能

### 7.1 角色模板
- 支持创建角色模板，快速创建相似角色
- 模板包含常用的性格标签和背景设定

### 7.2 角色导入导出
- 支持批量导入角色数据
- 支持导出角色关系图为图片或数据文件

### 7.3 AI集成
- 基于角色设定生成对话内容
- 根据角色关系推荐剧情发展
- 自动分析角色在文档中的出现频率