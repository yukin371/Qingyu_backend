# 写作端API完整文档

> **版本**: v1.0  
> **最后更新**: 2025-10-19  
> **基础路径**: `/api/v1/writer`

---

## 📋 目录

1. [项目管理API](#项目管理api)
2. [文档管理API](#文档管理api)
3. [版本控制API](#版本控制api)
4. [编辑器API](#编辑器api)
5. [内容审核API](#内容审核api)
6. [数据统计API](#数据统计api)
7. [通用规范](#通用规范)

---

## 🔐 认证说明

所有API都需要JWT认证，请在请求头中携带：

```http
Authorization: Bearer {your_jwt_token}
```

---

## 1. 项目管理API

### 1.1 创建项目

**接口**: `POST /api/v1/writer/projects`

**描述**: 创建新的写作项目

**请求体**:
```json
{
  "title": "我的第一本小说",
  "description": "这是一本精彩的玄幻小说",
  "genre": "玄幻",
  "tags": ["修仙", "热血", "爽文"],
  "cover_url": "https://example.com/cover.jpg"
}
```

**响应示例**:
```json
{
  "code": 201,
  "message": "创建成功",
  "data": {
    "project_id": "proj_123456",
    "title": "我的第一本小说",
    "status": "draft",
    "created_at": "2025-10-19T10:00:00Z"
  }
}
```

**状态码**:
- `201` - 创建成功
- `400` - 参数错误
- `401` - 未授权
- `500` - 服务器错误

---

### 1.2 获取项目详情

**接口**: `GET /api/v1/writer/projects/{project_id}`

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "project_id": "proj_123456",
    "title": "我的第一本小说",
    "description": "这是一本精彩的玄幻小说",
    "genre": "玄幻",
    "tags": ["修仙", "热血", "爽文"],
    "status": "draft",
    "word_count": 50000,
    "chapter_count": 25,
    "created_at": "2025-10-19T10:00:00Z",
    "updated_at": "2025-10-19T15:30:00Z"
  }
}
```

---

### 1.3 更新项目

**接口**: `PUT /api/v1/writer/projects/{project_id}`

**请求体**:
```json
{
  "title": "修改后的标题",
  "description": "修改后的简介",
  "tags": ["修仙", "热血", "爽文", "新标签"]
}
```

---

### 1.4 删除项目

**接口**: `DELETE /api/v1/writer/projects/{project_id}`

**描述**: 软删除项目

**响应**:
```json
{
  "code": 200,
  "message": "删除成功"
}
```

---

### 1.5 获取项目列表

**接口**: `GET /api/v1/writer/projects`

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20，最大100）
- `status`: 项目状态（draft/published/completed）
- `sort_by`: 排序字段（created_at/updated_at/word_count）
- `order`: 排序方向（asc/desc）

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "projects": [
      {
        "project_id": "proj_123456",
        "title": "我的第一本小说",
        "genre": "玄幻",
        "status": "draft",
        "word_count": 50000,
        "chapter_count": 25,
        "created_at": "2025-10-19T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 45,
      "total_pages": 3
    }
  }
}
```

---

### 1.6 更新项目统计

**接口**: `POST /api/v1/writer/projects/{project_id}/statistics`

**描述**: 手动触发项目统计数据重算

**响应**:
```json
{
  "code": 200,
  "message": "统计更新成功",
  "data": {
    "word_count": 52000,
    "chapter_count": 26,
    "updated_at": "2025-10-19T16:00:00Z"
  }
}
```

---

## 2. 文档管理API

### 2.1 创建文档

**接口**: `POST /api/v1/writer/documents`

**请求体**:
```json
{
  "project_id": "proj_123456",
  "title": "第一章：穿越异世",
  "type": "chapter",
  "parent_id": "",
  "order": 1,
  "content": "章节内容..."
}
```

**文档类型**:
- `chapter` - 章节
- `volume` - 卷
- `section` - 节
- `note` - 笔记

---

### 2.2 获取文档详情

**接口**: `GET /api/v1/writer/documents/{document_id}`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "document_id": "doc_789012",
    "project_id": "proj_123456",
    "title": "第一章：穿越异世",
    "type": "chapter",
    "content": "章节内容...",
    "word_count": 2000,
    "version": 5,
    "status": "draft",
    "created_at": "2025-10-19T10:30:00Z",
    "updated_at": "2025-10-19T15:45:00Z"
  }
}
```

---

### 2.3 获取文档树

**接口**: `GET /api/v1/writer/projects/{project_id}/documents/tree`

**描述**: 获取项目的文档树结构

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "tree": [
      {
        "document_id": "vol_001",
        "title": "第一卷：新手村",
        "type": "volume",
        "order": 1,
        "children": [
          {
            "document_id": "chap_001",
            "title": "第一章：穿越异世",
            "type": "chapter",
            "order": 1,
            "word_count": 2000
          }
        ]
      }
    ]
  }
}
```

---

### 2.4 更新文档

**接口**: `PUT /api/v1/writer/documents/{document_id}`

**请求体**:
```json
{
  "title": "修改后的标题",
  "content": "修改后的内容...",
  "version": 5
}
```

**注意**: 必须提供正确的版本号，防止版本冲突

---

### 2.5 移动文档

**接口**: `POST /api/v1/writer/documents/{document_id}/move`

**请求体**:
```json
{
  "new_parent_id": "vol_002",
  "new_order": 3
}
```

---

### 2.6 文档排序

**接口**: `POST /api/v1/writer/projects/{project_id}/documents/reorder`

**请求体**:
```json
{
  "parent_id": "vol_001",
  "document_orders": [
    {"document_id": "chap_002", "order": 1},
    {"document_id": "chap_001", "order": 2},
    {"document_id": "chap_003", "order": 3}
  ]
}
```

---

### 2.7 删除文档

**接口**: `DELETE /api/v1/writer/documents/{document_id}`

**查询参数**:
- `cascade`: 是否级联删除子文档（true/false，默认false）

---

### 2.8 获取文档列表

**接口**: `GET /api/v1/writer/projects/{project_id}/documents`

**查询参数**:
- `type`: 文档类型筛选
- `status`: 状态筛选
- `page`: 页码
- `page_size`: 每页数量

---

## 3. 版本控制API

### 3.1 获取版本历史

**接口**: `GET /api/v1/writer/documents/{document_id}/versions`

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "versions": [
      {
        "version": 5,
        "content_hash": "abc123def456",
        "word_count": 2100,
        "change_summary": "修改了第三段内容",
        "created_at": "2025-10-19T15:45:00Z",
        "created_by": "author_001"
      },
      {
        "version": 4,
        "content_hash": "xyz789uvw012",
        "word_count": 2000,
        "created_at": "2025-10-19T14:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 5
    }
  }
}
```

---

### 3.2 获取特定版本

**接口**: `GET /api/v1/writer/documents/{document_id}/versions/{version}`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "version": 4,
    "content": "这是版本4的完整内容...",
    "word_count": 2000,
    "created_at": "2025-10-19T14:30:00Z"
  }
}
```

---

### 3.3 版本对比

**接口**: `GET /api/v1/writer/documents/{document_id}/versions/compare`

**查询参数**:
- `from_version`: 起始版本号
- `to_version`: 目标版本号

**响应**:
```json
{
  "code": 200,
  "message": "对比成功",
  "data": {
    "from_version": 3,
    "to_version": 5,
    "diff": {
      "additions": 150,
      "deletions": 50,
      "changes": [
        {
          "type": "added",
          "position": 120,
          "content": "新增的一段文字..."
        },
        {
          "type": "deleted",
          "position": 80,
          "content": "删除的文字..."
        }
      ]
    }
  }
}
```

---

### 3.4 恢复版本

**接口**: `POST /api/v1/writer/documents/{document_id}/versions/{version}/restore`

**描述**: 将文档恢复到指定版本

**响应**:
```json
{
  "code": 200,
  "message": "恢复成功",
  "data": {
    "new_version": 6,
    "restored_from": 3
  }
}
```

---

## 4. 编辑器API

### 4.1 自动保存

**接口**: `POST /api/v1/writer/documents/{document_id}/autosave`

**请求体**:
```json
{
  "content": "当前编辑的内容...",
  "version": 5,
  "cursor_position": 1250
}
```

**响应**:
```json
{
  "code": 200,
  "message": "保存成功",
  "data": {
    "success": true,
    "new_version": 6,
    "save_time": "2025-10-19T16:30:00Z"
  }
}
```

**错误响应（版本冲突）**:
```json
{
  "code": 409,
  "message": "版本冲突",
  "error": "文档已被其他客户端修改，请刷新后重试",
  "data": {
    "current_version": 7,
    "your_version": 5
  }
}
```

---

### 4.2 获取保存状态

**接口**: `GET /api/v1/writer/documents/{document_id}/save-status`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "status": "saved",
    "last_save_time": "2025-10-19T16:30:00Z",
    "has_unsaved_changes": false,
    "current_version": 6
  }
}
```

**状态说明**:
- `saved` - 已保存
- `saving` - 保存中
- `unsaved` - 有未保存更改
- `error` - 保存失败

---

### 4.3 获取文档内容

**接口**: `GET /api/v1/writer/documents/{document_id}/content`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "content": "文档的完整内容...",
    "version": 6,
    "word_count": 2150,
    "updated_at": "2025-10-19T16:30:00Z"
  }
}
```

---

### 4.4 更新文档内容

**接口**: `PUT /api/v1/writer/documents/{document_id}/content`

**请求体**:
```json
{
  "content": "更新后的内容...",
  "version": 6
}
```

---

### 4.5 计算字数

**接口**: `POST /api/v1/writer/word-count`

**请求体**:
```json
{
  "content": "要计算字数的内容...",
  "filter_markdown": true
}
```

**响应**:
```json
{
  "code": 200,
  "message": "计算成功",
  "data": {
    "total_words": 2150,
    "chinese_chars": 2000,
    "english_words": 50,
    "digits": 100,
    "paragraphs": 15,
    "estimated_read_time": 4.3
  }
}
```

---

### 4.6 获取用户快捷键配置

**接口**: `GET /api/v1/writer/user/shortcuts`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "shortcuts": [
      {
        "action": "save",
        "key": "Ctrl+S",
        "description": "保存文档"
      },
      {
        "action": "undo",
        "key": "Ctrl+Z",
        "description": "撤销"
      },
      {
        "action": "redo",
        "key": "Ctrl+Y",
        "description": "重做"
      }
    ]
  }
}
```

---

### 4.7 更新用户快捷键

**接口**: `PUT /api/v1/writer/user/shortcuts`

**请求体**:
```json
{
  "shortcuts": [
    {"action": "save", "key": "Ctrl+S"},
    {"action": "bold", "key": "Ctrl+B"}
  ]
}
```

---

### 4.8 重置快捷键

**接口**: `POST /api/v1/writer/user/shortcuts/reset`

**描述**: 重置为默认快捷键配置

---

### 4.9 获取快捷键帮助

**接口**: `GET /api/v1/writer/shortcuts/help`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "categories": [
      {
        "name": "文件操作",
        "shortcuts": [
          {"action": "save", "key": "Ctrl+S", "description": "保存文档"},
          {"action": "save_as", "key": "Ctrl+Shift+S", "description": "另存为"}
        ]
      },
      {
        "name": "编辑操作",
        "shortcuts": [
          {"action": "undo", "key": "Ctrl+Z"},
          {"action": "redo", "key": "Ctrl+Y"},
          {"action": "cut", "key": "Ctrl+X"},
          {"action": "copy", "key": "Ctrl+C"},
          {"action": "paste", "key": "Ctrl+V"}
        ]
      }
    ]
  }
}
```

---

## 5. 内容审核API

### 5.1 实时内容检测

**接口**: `POST /api/v1/writer/audit/check`

**请求体**:
```json
{
  "content": "要检测的内容文本...",
  "content_id": "doc_789012",
  "author_id": "author_001"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "检测完成",
  "data": {
    "safe": false,
    "risk_level": "high",
    "matched_words": [
      {"word": "敏感词1", "category": "politics", "position": [10, 15]},
      {"word": "敏感词2", "category": "violence", "position": [120, 125]}
    ],
    "suggestions": [
      "建议删除或修改包含敏感词的段落"
    ]
  }
}
```

**风险等级**:
- `low` - 低风险（通过）
- `medium` - 中风险（警告）
- `high` - 高风险（拒绝）

---

### 5.2 全文审核

**接口**: `POST /api/v1/writer/documents/{document_id}/audit`

**请求体**:
```json
{
  "audit_type": "full",
  "options": {
    "check_sensitive_words": true,
    "check_rules": true,
    "auto_fix": false
  }
}
```

**响应**:
```json
{
  "code": 200,
  "message": "审核完成",
  "data": {
    "audit_id": "audit_123456",
    "status": "rejected",
    "risk_level": "high",
    "issues_count": 3,
    "audit_time": "2025-10-19T17:00:00Z"
  }
}
```

---

### 5.3 获取审核结果

**接口**: `GET /api/v1/writer/documents/{document_id}/audit-result`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "audit_id": "audit_123456",
    "document_id": "doc_789012",
    "status": "rejected",
    "risk_level": "high",
    "issues": [
      {
        "type": "sensitive_word",
        "severity": "high",
        "position": 120,
        "content": "敏感内容片段",
        "suggestion": "建议删除或修改"
      }
    ],
    "audited_at": "2025-10-19T17:00:00Z",
    "auditor": "system"
  }
}
```

---

### 5.4 提交申诉

**接口**: `POST /api/v1/writer/audit/{audit_id}/appeal`

**请求体**:
```json
{
  "reason": "申诉理由...",
  "evidence": "相关证据和说明...",
  "contact": "联系方式"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "申诉已提交",
  "data": {
    "appeal_id": "appeal_123456",
    "status": "pending",
    "submitted_at": "2025-10-19T17:30:00Z",
    "estimated_review_time": "24小时内"
  }
}
```

---

### 5.5 获取用户违规记录

**接口**: `GET /api/v1/writer/user/violations`

**查询参数**:
- `page`: 页码
- `page_size`: 每页数量

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "violations": [
      {
        "violation_id": "vio_123456",
        "content_id": "doc_789012",
        "type": "sensitive_content",
        "risk_level": "high",
        "penalty": "warning",
        "created_at": "2025-10-19T16:00:00Z",
        "status": "active"
      }
    ],
    "summary": {
      "total_violations": 2,
      "active_penalties": 1,
      "warning_count": 1,
      "ban_count": 0
    }
  }
}
```

---

### 5.6 获取待审核列表（管理员）

**接口**: `GET /api/v1/writer/admin/audit/pending`

**权限**: 管理员

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "pending_audits": [
      {
        "audit_id": "audit_123456",
        "document_id": "doc_789012",
        "author_id": "author_001",
        "risk_level": "high",
        "submitted_at": "2025-10-19T17:00:00Z"
      }
    ],
    "total": 15
  }
}
```

---

### 5.7 人工审核（管理员）

**接口**: `POST /api/v1/writer/admin/audit/{audit_id}/review`

**权限**: 管理员

**请求体**:
```json
{
  "decision": "approved",
  "comment": "审核意见...",
  "reviewer_id": "admin_001"
}
```

**决策选项**:
- `approved` - 通过
- `rejected` - 拒绝
- `needs_revision` - 需要修改

---

## 6. 数据统计API

### 6.1 获取作品统计

**接口**: `GET /api/v1/writer/books/{book_id}/stats`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "book_id": "proj_123456",
    "total_views": 150000,
    "unique_readers": 50000,
    "total_chapters": 50,
    "avg_completion_rate": 0.72,
    "avg_drop_off_rate": 0.28,
    "avg_reading_duration": 180,
    "total_revenue": 12500.00,
    "view_trend": "up",
    "revenue_trend": "up",
    "stat_date": "2025-10-19"
  }
}
```

---

### 6.2 获取章节统计

**接口**: `GET /api/v1/writer/books/{book_id}/chapters/{chapter_id}/stats`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "chapter_id": "chap_001",
    "chapter_title": "第一章：穿越异世",
    "view_count": 8000,
    "unique_viewers": 5000,
    "completion_rate": 0.75,
    "drop_off_rate": 0.25,
    "avg_read_duration": 180,
    "revenue": 250.00,
    "comments_count": 120,
    "likes_count": 300
  }
}
```

---

### 6.3 生成章节阅读热力图

**接口**: `GET /api/v1/writer/books/{book_id}/chapters/{chapter_id}/heatmap`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "heatmap": [
      {"position": 0, "intensity": 100},
      {"position": 500, "intensity": 95},
      {"position": 1000, "intensity": 80},
      {"position": 1500, "intensity": 60},
      {"position": 2000, "intensity": 45}
    ],
    "analysis": {
      "hottest_section": {"start": 0, "end": 600},
      "coldest_section": {"start": 1800, "end": 2200},
      "avg_intensity": 76
    }
  }
}
```

---

### 6.4 获取收入明细

**接口**: `GET /api/v1/writer/books/{book_id}/revenue`

**查询参数**:
- `start_date`: 开始日期（YYYY-MM-DD）
- `end_date`: 结束日期（YYYY-MM-DD）

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "book_id": "proj_123456",
    "start_date": "2025-10-01",
    "end_date": "2025-10-19",
    "chapter_revenue": 8000.00,
    "subscribe_revenue": 3500.00,
    "reward_revenue": 1000.00,
    "total_revenue": 12500.00,
    "daily_breakdown": [
      {"date": "2025-10-19", "revenue": 650.00},
      {"date": "2025-10-18", "revenue": 580.00}
    ]
  }
}
```

---

### 6.5 获取热门章节

**接口**: `GET /api/v1/writer/books/{book_id}/top-chapters`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "most_viewed": [
      {"chapter_id": "chap_001", "title": "第一章", "view_count": 8000}
    ],
    "highest_revenue": [
      {"chapter_id": "chap_025", "title": "第二十五章", "revenue": 500.00}
    ],
    "lowest_completion": [
      {"chapter_id": "chap_015", "title": "第十五章", "completion_rate": 0.45}
    ],
    "highest_drop_off": [
      {"chapter_id": "chap_012", "title": "第十二章", "drop_off_rate": 0.55}
    ]
  }
}
```

---

### 6.6 获取每日统计

**接口**: `GET /api/v1/writer/books/{book_id}/daily-stats`

**查询参数**:
- `days`: 天数（默认7，最大365）

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "daily_stats": [
      {
        "date": "2025-10-19",
        "daily_views": 5000,
        "new_readers": 500,
        "daily_revenue": 650.00,
        "avg_read_time": 12.5
      }
    ]
  }
}
```

---

### 6.7 获取章节跳出点

**接口**: `GET /api/v1/writer/books/{book_id}/chapters/{chapter_id}/drop-off-points`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "drop_off_points": [
      {
        "position": 1200,
        "drop_off_count": 150,
        "drop_off_rate": 0.15,
        "context": "...周围的文字内容..."
      }
    ]
  }
}
```

---

### 6.8 记录读者行为

**接口**: `POST /api/v1/writer/reader-behavior`

**请求体**:
```json
{
  "user_id": "reader_001",
  "book_id": "proj_123456",
  "chapter_id": "chap_001",
  "behavior_type": "read",
  "read_duration": 180,
  "read_position": 1500,
  "read_at": "2025-10-19T18:00:00Z"
}
```

---

### 6.9 计算留存率

**接口**: `GET /api/v1/writer/books/{book_id}/retention-rate`

**查询参数**:
- `days`: 天数（默认7，最大90）

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "book_id": "proj_123456",
    "days": 7,
    "retention_rate": 0.65,
    "initial_readers": 1000,
    "retained_readers": 650,
    "calculated_at": "2025-10-19T18:00:00Z"
  }
}
```

---

## 7. 通用规范

### 7.1 统一响应格式

所有API响应都遵循统一格式：

**成功响应**:
```json
{
  "code": 200,
  "message": "操作描述",
  "data": {
    // 响应数据
  }
}
```

**错误响应**:
```json
{
  "code": 400,
  "message": "错误描述",
  "error": "详细错误信息"
}
```

---

### 7.2 HTTP状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 204 | 删除成功（无返回内容） |
| 400 | 参数错误 |
| 401 | 未授权 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突（如版本冲突） |
| 422 | 验证失败 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |
| 503 | 服务不可用 |

---

### 7.3 分页参数

所有列表接口支持统一的分页参数：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码，从1开始 |
| page_size | int | 20 | 每页数量，最大100 |
| sort_by | string | created_at | 排序字段 |
| order | string | desc | 排序方向（asc/desc） |

---

### 7.4 时间格式

所有时间字段使用 ISO 8601 格式：

```
2025-10-19T18:30:00Z
```

---

### 7.5 错误码

| 错误码 | 说明 |
|--------|------|
| INVALID_PARAMETER | 参数错误 |
| UNAUTHORIZED | 未授权 |
| FORBIDDEN | 权限不足 |
| NOT_FOUND | 资源不存在 |
| VERSION_CONFLICT | 版本冲突 |
| RATE_LIMIT_EXCEEDED | 请求频率超限 |
| INTERNAL_ERROR | 服务器内部错误 |
| AUDIT_REJECTED | 审核未通过 |

---

### 7.6 速率限制

| API类型 | 限制 |
|---------|------|
| 读取操作 | 1000次/小时 |
| 写入操作 | 100次/小时 |
| 自动保存 | 600次/小时（每分钟10次） |
| 审核检测 | 200次/小时 |

---

## 📚 使用示例

### 完整工作流示例

```javascript
// 1. 创建项目
const project = await fetch('/api/v1/writer/projects', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_TOKEN',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    title: '我的小说',
    genre: '玄幻'
  })
});

const { project_id } = await project.json();

// 2. 创建章节
const chapter = await fetch('/api/v1/writer/documents', {
  method: 'POST',
  body: JSON.stringify({
    project_id: project_id,
    title: '第一章',
    type: 'chapter'
  })
});

const { document_id } = await chapter.json();

// 3. 自动保存
const save = await fetch(`/api/v1/writer/documents/${document_id}/autosave`, {
  method: 'POST',
  body: JSON.stringify({
    content: '章节内容...',
    version: 1
  })
});

// 4. 内容审核
const audit = await fetch('/api/v1/writer/audit/check', {
  method: 'POST',
  body: JSON.stringify({
    content: '章节内容...',
    content_id: document_id
  })
});

// 5. 查看统计
const stats = await fetch(`/api/v1/writer/books/${project_id}/stats`);
```

---

## 📞 技术支持

如有问题，请联系：

- **GitHub**: https://github.com/qingyu/backend

---

**最后更新**: 2025-10-19  
**文档版本**: v1.0
