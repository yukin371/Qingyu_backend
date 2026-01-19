# 编辑器API文档

> **版本**: v1.0  
> **创建日期**: 2025-10-21

---

## 1. API概述

编辑器API提供文档编辑、自动保存、版本管理等功能。

**Base URL**: `/api/v1/editor`

---

## 2. 文档操作

### 2.1 获取文档内容

**接口**: `GET /editor/documents/:documentId`

**响应**：
```json
{
  "code": 200,
  "data": {
    "documentId": "doc_123",
    "title": "第一章",
    "content": "文档内容...",
    "wordCount": 3000,
    "version": 5,
    "updatedAt": "2025-10-21T10:00:00Z"
  }
}
```

### 2.2 保存文档

**接口**: `PUT /editor/documents/:documentId`

**请求**：
```json
{
  "content": "更新的文档内容",
  "isAutoSave": false,
  "comment": "修改了开头部分"
}
```

**响应**：
```json
{
  "code": 200,
  "message": "保存成功",
  "data": {
    "documentId": "doc_123",
    "version": 6,
    "wordCount": 3100
  }
}
```

---

## 3. 自动保存

### 3.1 自动保存

**接口**: `POST /editor/documents/:documentId/autosave`

**请求**：
```json
{
  "content": "当前编辑的内容",
  "cursorPosition": 1250
}
```

**说明**：
- 客户端每30秒自动调用
- 不增加版本号
- 仅更新内容

---

## 4. 版本管理

### 4.1 获取版本列表

**接口**: `GET /editor/documents/:documentId/versions`

**响应**：
```json
{
  "code": 200,
  "data": {
    "versions": [
      {
        "versionNum": 6,
        "comment": "修改了开头部分",
        "wordCount": 3100,
        "createdAt": "2025-10-21T10:00:00Z",
        "isAutoSave": false
      }
    ],
    "total": 6
  }
}
```

### 4.2 恢复历史版本

**接口**: `POST /editor/documents/:documentId/versions/:versionNum/restore`

**响应**：
```json
{
  "code": 200,
  "message": "恢复成功",
  "data": {
    "documentId": "doc_123",
    "version": 7
  }
}
```

---

## 5. AI辅助

### 5.1 AI续写

**接口**: `POST /editor/ai/continue`

**请求**：
```json
{
  "documentId": "doc_123",
  "context": "当前文本上下文...",
  "length": 500
}
```

**响应（SSE流式）**：
```
data: {"text": "AI生成的", "done": false}
data: {"text": "文本内容", "done": false}
data: {"text": "", "done": true}
```

---

## 6. 快捷键

### 6.1 获取用户快捷键配置

**接口**: `GET /editor/shortcuts`

**响应**：
```json
{
  "code": 200,
  "data": {
    "shortcuts": {
      "save": {"key": "Ctrl+S", "description": "保存"},
      "bold": {"key": "Ctrl+B", "description": "加粗"}
    }
  }
}
```

### 6.2 更新快捷键

**接口**: `PUT /editor/shortcuts/:action`

**请求**：
```json
{
  "key": "Ctrl+Shift+S"
}
```

---

**文档状态**: ✅ 已完成

