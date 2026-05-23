# 审核API文档

> **版本**: v1.0  
> **创建日期**: 2025-10-21

## Page Role

- current-reference
- current-owner: `docs/api/`
- current-bounded: 当前审核 API 参考文档，负责审核相关接口说明

## Recommended Read Path

1. 先读 `README.md`。
2. 再读本页。

## Boundary

- 本页是审核接口参考，不是审核能力总 owner。
- 当前审核相关架构与实现仍需结合其他 current owner 文档确认。

## Quick Section Map

- API概述
- 内容审核
- 敏感词管理
- 违规记录

## Quick Takeaways

- 这是当前审核 API 参考文档。

## Skip Guide

- 只看总入口：回 `README.md`。

---

## 1. API概述

审核API提供内容审核、敏感词管理、违规记录等功能。

**Base URL**: `/api/v1/audit`

---

## 2. 内容审核

### 2.1 审核文本

**接口**: `POST /audit/text`

**请求**：
```json
{
  "targetType": "document",
  "targetId": "doc_123",
  "content": "待审核的文本内容",
  "authorId": "user_123"
}
```

**响应**：
```json
{
  "code": 200,
  "message": "审核完成",
  "data": {
    "auditId": "audit_456",
    "status": "approved",  // approved/rejected/pending
    "riskLevel": 1,        // 1-5
    "riskScore": 12.5,
    "violations": []
  }
}
```

---

## 3. 敏感词管理

### 3.1 添加敏感词

**接口**: `POST /audit/sensitive-words`

**权限**: 管理员

**请求**：
```json
{
  "word": "敏感词",
  "category": "政治",
  "level": 5,
  "replacement": "***"
}
```

### 3.2 获取敏感词列表

**接口**: `GET /audit/sensitive-words`

**参数**：
- `category` - 分类
- `page` - 页码
- `pageSize` - 每页数量

---

## 4. 违规记录

### 4.1 获取用户违规记录

**接口**: `GET /audit/violations/user/:userId`

**响应**：
```json
{
  "code": 200,
  "data": {
    "violations": [
      {
        "id": "violation_123",
        "targetType": "document",
        "violationType": "sensitive_word",
        "violationLevel": 3,
        "isPenalized": true,
        "createdAt": "2025-10-21T10:00:00Z"
      }
    ],
    "total": 5
  }
}
```

---

**文档状态**: ✅ 已完成

