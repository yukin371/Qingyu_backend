# 共享服务 API 参考

> **版本**: v1.3 ⭐️已更新  
> **最后更新**: 2025-10-25  
> **基础路径**: `/api/v1/shared`

---

## 1. 概述

共享服务提供平台通用功能，包括钱包管理、文件存储、认证服务、管理员功能等，为阅读端和写作端提供基础支持。

### 1.1 基础信息

- **认证要求**: 所有接口均需登录认证
- **响应格式**: 统一 JSON 格式
- **Token**: 在请求头中添加 `Authorization: Bearer <token>`
- **文件上传**: 支持 multipart/form-data

### 1.2 功能特性

- ✅ 钱包管理（余额查询、充值、消费、转账）
- ✅ 交易记录查询
- ✅ 提现申请和管理
- ✅ 文件上传/下载/删除
- ✅ 文件列表和信息查询
- ✅ 认证服务（登录、登出、Token 刷新）
- ✅ 权限和角色管理
- ✅ 管理员审核功能

---

## 1.3 统一响应格式 ⭐️v1.3更新

### 成功响应
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    // 业务数据
  },
  "timestamp": 1729875123,
  "request_id": "req-12345-abcde"
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "参数错误",
  "error": "详细错误信息",
  "timestamp": 1729875123,
  "request_id": "req-12345-abcde"
}
```

**新增字段说明**:
- `timestamp`: Unix时间戳，服务器响应时间
- `request_id`: 请求追踪ID，便于调试和日志追踪（可选字段）

---

## 1.4 TypeScript 类型定义 ⭐️v1.3新增

```typescript
// src/types/shared.ts
import type { APIResponse, PaginatedResponse } from './api';

// 钱包信息
export interface WalletInfo {
  userId: string;
  balance: number;
  frozenBalance: number;
  totalIncome: number;
  totalExpense: number;
  currency: string;
}

// 交易记录
export interface Transaction {
  id: string;
  userId: string;
  type: 'recharge' | 'consume' | 'transfer' | 'withdraw';
  amount: number;
  balance: number;
  reason: string;
  status: 'pending' | 'success' | 'failed';
  createdAt: string;
}

// 文件信息
export interface FileInfo {
  id: string;
  filename: string;
  originalName: string;
  size: number;
  mimeType: string;
  url: string;
  category: string;
  uploadedAt: string;
}

// 提现申请
export interface WithdrawRequest {
  id: string;
  userId: string;
  amount: number;
  account: string;
  status: 'pending' | 'approved' | 'rejected';
  reason?: string;
  createdAt: string;
  processedAt?: string;
}

// API函数
export const getWalletBalance = () => {
  return request.get<APIResponse<{ balance: number }>>('/shared/wallet/balance');
};

export const recharge = (amount: number, method: string) => {
  return request.post<APIResponse<Transaction>>('/shared/wallet/recharge', { amount, method });
};

export const uploadFile = (file: File, path?: string) => {
  const formData = new FormData();
  formData.append('file', file);
  if (path) formData.append('path', path);
  
  return request.post<APIResponse<FileInfo>>('/shared/storage/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  });
};

export const getTransactionHistory = (params: { page?: number; page_size?: number; type?: string }) => {
  return request.get<PaginatedResponse<Transaction>>('/shared/wallet/transactions', { params });
};
```

---

## 2. 钱包服务 API

### 2.1 查询余额

**接口说明**: 查询用户钱包余额

**请求**
```
GET /api/v1/shared/wallet/balance
```

**请求头**
```
Authorization: Bearer <token>
```

**响应示例**
```json
{
  "code": 200,
  "message": "查询余额成功",
  "data": {
    "balance": 1250.50,
    "frozenAmount": 50.00,
    "availableAmount": 1200.50
  }
}
```

**响应字段说明**

| 字段 | 类型 | 说明 |
|------|------|------|
| balance | float | 总余额 |
| frozenAmount | float | 冻结金额 |
| availableAmount | float | 可用金额 |

**错误响应**
- `401` - 未认证
- `500` - 服务器错误

**JavaScript/Axios 示例**
```javascript
const getBalance = async () => {
  try {
    const response = await axios.get('/api/v1/shared/wallet/balance');
    return response.data.data.balance;
  } catch (error) {
    console.error('查询余额失败:', error);
    throw error;
  }
};
```

---

### 2.2 获取钱包信息

**接口说明**: 获取用户完整钱包信息

**请求**
```
GET /api/v1/shared/wallet
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取钱包信息成功",
  "data": {
    "walletId": "wallet_123456",
    "userId": "user_123",
    "balance": 1250.50,
    "frozenAmount": 50.00,
    "availableAmount": 1200.50,
    "totalIncome": 5000.00,
    "totalExpense": 3749.50,
    "createdAt": "2025-01-01T00:00:00Z",
    "updatedAt": "2025-10-20T10:00:00Z"
  }
}
```

---

### 2.3 充值

**接口说明**: 用户钱包充值

**请求**
```
POST /api/v1/shared/wallet/recharge
```

**请求体**
```json
{
  "amount": 100.00,
  "method": "alipay"
}
```

**请求参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| amount | float | 是 | 充值金额（0.01-10000） |
| method | string | 是 | 支付方式（alipay/wechat/bank） |

**响应示例**
```json
{
  "code": 200,
  "message": "充值成功",
  "data": {
    "transactionId": "txn_123456",
    "amount": 100.00,
    "method": "alipay",
    "status": "completed",
    "createdAt": "2025-10-20T10:00:00Z"
  }
}
```

**错误响应**
- `400` - 参数错误（金额不合法）
- `401` - 未认证
- `500` - 充值失败

**JavaScript/Axios 示例**
```javascript
const recharge = async (amount, method = 'alipay') => {
  try {
    const response = await axios.post('/api/v1/shared/wallet/recharge', {
      amount,
      method
    });
    
    // 返回交易ID
    return response.data.data.transactionId;
  } catch (error) {
    if (error.response?.status === 400) {
      throw new Error('充值金额不合法');
    }
    throw error;
  }
};
```

---

### 2.4 消费

**接口说明**: 用户钱包消费

**请求**
```
POST /api/v1/shared/wallet/consume
```

**请求体**
```json
{
  "amount": 50.00,
  "reason": "购买章节"
}
```

**请求参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| amount | float | 是 | 消费金额 |
| reason | string | 是 | 消费原因（1-200字） |

**响应示例**
```json
{
  "code": 200,
  "message": "消费成功",
  "data": {
    "transactionId": "txn_123457",
    "amount": 50.00,
    "reason": "购买章节",
    "remainingBalance": 1150.50,
    "createdAt": "2025-10-20T10:05:00Z"
  }
}
```

**错误响应**
- `400` - 参数错误
- `401` - 未认证
- `402` - 余额不足
- `500` - 消费失败

**JavaScript/Axios 示例**
```javascript
const consume = async (amount, reason) => {
  try {
    const response = await axios.post('/api/v1/shared/wallet/consume', {
      amount,
      reason
    });
    
    return response.data.data;
  } catch (error) {
    if (error.response?.status === 402) {
      throw new Error('余额不足');
    }
    throw error;
  }
};
```

---

### 2.5 转账

**接口说明**: 向其他用户转账

**请求**
```
POST /api/v1/shared/wallet/transfer
```

**请求体**
```json
{
  "to_user_id": "user_456",
  "amount": 100.00,
  "reason": "赠送"
}
```

**请求参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| to_user_id | string | 是 | 接收方用户ID |
| amount | float | 是 | 转账金额 |
| reason | string | 否 | 转账原因（最多200字） |

**响应示例**
```json
{
  "code": 200,
  "message": "转账成功",
  "data": {
    "transactionId": "txn_123458",
    "fromUserId": "user_123",
    "toUserId": "user_456",
    "amount": 100.00,
    "reason": "赠送",
    "createdAt": "2025-10-20T10:10:00Z"
  }
}
```

**错误响应**
- `400` - 参数错误（用户ID无效）
- `402` - 余额不足
- `404` - 接收方用户不存在

---

### 2.6 查询交易记录

**接口说明**: 查询用户交易记录列表

**请求**
```
GET /api/v1/shared/wallet/transactions
```

**Query 参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |
| type | string | 否 | - | 交易类型（recharge/consume/transfer） |

**响应示例**
```json
{
  "code": 200,
  "message": "查询交易记录成功",
  "data": [
    {
      "transactionId": "txn_123456",
      "type": "recharge",
      "amount": 100.00,
      "balance": 1250.50,
      "description": "支付宝充值",
      "createdAt": "2025-10-20T10:00:00Z"
    },
    {
      "transactionId": "txn_123457",
      "type": "consume",
      "amount": -50.00,
      "balance": 1200.50,
      "description": "购买章节",
      "createdAt": "2025-10-20T10:05:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "page_size": 20
}
```

**JavaScript/Axios 示例**
```javascript
const getTransactions = async (page = 1, pageSize = 20, type = '') => {
  try {
    const response = await axios.get('/api/v1/shared/wallet/transactions', {
      params: {
        page,
        page_size: pageSize,
        type
      }
    });
    
    return response.data.data;
  } catch (error) {
    console.error('查询交易记录失败:', error);
    throw error;
  }
};
```

---

### 2.7 申请提现

**接口说明**: 用户申请提现

**请求**
```
POST /api/v1/shared/wallet/withdraw
```

**请求体**
```json
{
  "amount": 500.00,
  "account": "支付宝账号:example@example.com"
}
```

**请求参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| amount | float | 是 | 提现金额 |
| account | string | 是 | 提现账号 |

**响应示例**
```json
{
  "code": 200,
  "message": "申请提现成功",
  "data": {
    "withdrawId": "withdraw_123456",
    "amount": 500.00,
    "account": "支付宝账号:example@example.com",
    "status": "pending",
    "createdAt": "2025-10-20T10:00:00Z"
  }
}
```

**错误响应**
- `400` - 参数错误
- `402` - 余额不足
- `429` - 提现次数超限

**JavaScript/Axios 示例**
```javascript
const requestWithdraw = async (amount, account) => {
  try {
    const response = await axios.post('/api/v1/shared/wallet/withdraw', {
      amount,
      account
    });
    
    return response.data.data;
  } catch (error) {
    if (error.response?.status === 402) {
      throw new Error('余额不足');
    } else if (error.response?.status === 429) {
      throw new Error('提现次数超限，请稍后再试');
    }
    throw error;
  }
};
```

---

### 2.8 查询提现申请

**接口说明**: 查询用户提现申请列表

**请求**
```
GET /api/v1/shared/wallet/withdrawals
```

**Query 参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |
| status | string | 否 | - | 状态（pending/approved/rejected） |

**响应示例**
```json
{
  "code": 200,
  "message": "查询提现申请成功",
  "data": [
    {
      "withdrawId": "withdraw_123456",
      "amount": 500.00,
      "account": "支付宝账号:example@example.com",
      "status": "pending",
      "createdAt": "2025-10-20T10:00:00Z",
      "reviewedAt": null,
      "reason": ""
    }
  ],
  "total": 5,
  "page": 1,
  "page_size": 20
}
```

---

## 3. 文件存储 API

### 3.1 上传文件

**接口说明**: 上传文件到存储服务

**请求**
```
POST /api/v1/shared/storage/upload
```

**Content-Type**: `multipart/form-data`

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 上传的文件 |
| path | string | 否 | 存储路径（默认: uploads/{userId}） |

**响应示例**
```json
{
  "code": 200,
  "message": "上传文件成功",
  "data": {
    "fileId": "file_123456",
    "filename": "cover.jpg",
    "size": 102400,
    "contentType": "image/jpeg",
    "url": "https://storage.qingyu.com/uploads/user_123/cover.jpg",
    "createdAt": "2025-10-20T10:00:00Z"
  }
}
```

**错误响应**
- `400` - 文件格式不支持或文件过大
- `401` - 未认证
- `500` - 上传失败

**JavaScript/Axios 示例**
```javascript
const uploadFile = async (file, path = '') => {
  try {
    const formData = new FormData();
    formData.append('file', file);
    if (path) {
      formData.append('path', path);
    }
    
    const response = await axios.post('/api/v1/shared/storage/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total);
        console.log(`上传进度: ${percent}%`);
      }
    });
    
    return response.data.data;
  } catch (error) {
    console.error('上传失败:', error);
    throw error;
  }
};
```

**Vue 3 文件上传组件示例**
```vue
<template>
  <div class="upload-container">
    <input 
      type="file" 
      @change="handleFileChange" 
      ref="fileInput"
      accept="image/*"
    />
    <button @click="uploadFile" :disabled="!selectedFile || uploading">
      {{ uploading ? '上传中...' : '上传文件' }}
    </button>
    <div v-if="uploadProgress > 0" class="progress-bar">
      <div :style="{ width: uploadProgress + '%' }"></div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';

const fileInput = ref(null);
const selectedFile = ref(null);
const uploading = ref(false);
const uploadProgress = ref(0);

const handleFileChange = (e) => {
  selectedFile.value = e.target.files[0];
};

const uploadFile = async () => {
  if (!selectedFile.value) return;
  
  uploading.value = true;
  uploadProgress.value = 0;
  
  try {
    const formData = new FormData();
    formData.append('file', selectedFile.value);
    
    const response = await axios.post('/api/v1/shared/storage/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (e) => {
        uploadProgress.value = Math.round((e.loaded * 100) / e.total);
      }
    });
    
    console.log('上传成功:', response.data.data.url);
    selectedFile.value = null;
  } catch (error) {
    console.error('上传失败:', error);
  } finally {
    uploading.value = false;
  }
};
</script>
```

---

### 3.2 下载文件

**接口说明**: 下载指定的文件

**请求**
```
GET /api/v1/shared/storage/download/{file_id}
```

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file_id | string | 是 | 文件ID |

**响应**: 文件二进制流

**错误响应**
- `400` - 文件ID无效
- `401` - 未认证
- `404` - 文件不存在
- `500` - 下载失败

**JavaScript/Axios 示例**
```javascript
const downloadFile = async (fileId, filename) => {
  try {
    const response = await axios.get(`/api/v1/shared/storage/download/${fileId}`, {
      responseType: 'blob'
    });
    
    // 创建下载链接
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', filename);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  } catch (error) {
    console.error('下载失败:', error);
    throw error;
  }
};
```

---

### 3.3 删除文件

**接口说明**: 删除指定的文件

**请求**
```
DELETE /api/v1/shared/storage/files/{file_id}
```

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file_id | string | 是 | 文件ID |

**响应示例**
```json
{
  "code": 200,
  "message": "删除文件成功"
}
```

**错误响应**
- `400` - 文件ID无效
- `401` - 未认证
- `403` - 无权限删除该文件
- `500` - 删除失败

---

### 3.4 获取文件信息

**接口说明**: 获取指定文件的详细信息

**请求**
```
GET /api/v1/shared/storage/files/{file_id}
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取文件信息成功",
  "data": {
    "fileId": "file_123456",
    "filename": "cover.jpg",
    "size": 102400,
    "contentType": "image/jpeg",
    "url": "https://storage.qingyu.com/uploads/user_123/cover.jpg",
    "userId": "user_123",
    "category": "uploads/user_123",
    "createdAt": "2025-10-20T10:00:00Z"
  }
}
```

---

### 3.5 列出文件

**接口说明**: 列出用户的文件列表

**请求**
```
GET /api/v1/shared/storage/files
```

**Query 参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |
| category | string | 否 | - | 文件分类/路径前缀 |

**响应示例**
```json
{
  "code": 200,
  "message": "列出文件成功",
  "data": [
    {
      "fileId": "file_123456",
      "filename": "cover.jpg",
      "size": 102400,
      "contentType": "image/jpeg",
      "url": "https://storage.qingyu.com/uploads/user_123/cover.jpg",
      "createdAt": "2025-10-20T10:00:00Z"
    }
  ],
  "total": 15,
  "page": 1,
  "page_size": 20
}
```

---

### 3.6 获取文件访问URL

**接口说明**: 获取文件的临时访问URL

**请求**
```
GET /api/v1/shared/storage/files/{file_id}/url
```

**Query 参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| expire | int | 否 | 3600 | 过期时间（秒） |

**响应示例**
```json
{
  "code": 200,
  "message": "获取文件URL成功",
  "data": {
    "url": "https://storage.qingyu.com/uploads/user_123/cover.jpg?token=xxx&expire=1698000000"
  }
}
```

**JavaScript/Axios 示例**
```javascript
const getFileURL = async (fileId, expire = 3600) => {
  try {
    const response = await axios.get(`/api/v1/shared/storage/files/${fileId}/url`, {
      params: { expire }
    });
    
    return response.data.data.url;
  } catch (error) {
    console.error('获取文件URL失败:', error);
    throw error;
  }
};
```

---

## 4. 认证服务 API

### 4.1 用户注册

**接口说明**: 注册新用户账号

**请求**
```
POST /api/v1/shared/auth/register
```

**请求体**
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

**请求参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名（3-50字符） |
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码（6-50字符） |

**响应示例**
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "userId": "user_123456",
    "username": "testuser",
    "email": "test@example.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**错误响应**
- `400` - 参数错误（用户名或邮箱格式不正确）
- `409` - 用户名或邮箱已存在
- `500` - 注册失败

---

### 4.2 用户登录

**接口说明**: 用户登录获取Token

**请求**
```
POST /api/v1/shared/auth/login
```

**请求体**
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**响应示例**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "userId": "user_123456",
    "username": "testuser",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresAt": "2025-10-21T10:00:00Z"
  }
}
```

**错误响应**
- `400` - 参数错误
- `401` - 用户名或密码错误
- `500` - 登录失败

**JavaScript/Axios 示例**
```javascript
const login = async (username, password) => {
  try {
    const response = await axios.post('/api/v1/shared/auth/login', {
      username,
      password
    });
    
    const { token, userId } = response.data.data;
    
    // 保存 Token
    localStorage.setItem('token', token);
    localStorage.setItem('userId', userId);
    
    // 设置 axios 默认请求头
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    
    return response.data.data;
  } catch (error) {
    if (error.response?.status === 401) {
      throw new Error('用户名或密码错误');
    }
    throw error;
  }
};
```

---

### 4.3 用户登出

**接口说明**: 用户登出，使Token失效

**请求**
```
POST /api/v1/shared/auth/logout
```

**请求头**
```
Authorization: Bearer <token>
```

**响应示例**
```json
{
  "code": 200,
  "message": "登出成功"
}
```

**JavaScript/Axios 示例**
```javascript
const logout = async () => {
  try {
    await axios.post('/api/v1/shared/auth/logout');
    
    // 清除本地存储
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    
    // 清除 axios 默认请求头
    delete axios.defaults.headers.common['Authorization'];
    
    // 跳转到登录页
    window.location.href = '/login';
  } catch (error) {
    console.error('登出失败:', error);
  }
};
```

---

### 4.4 刷新Token

**接口说明**: 使用当前Token获取新Token

**请求**
```
POST /api/v1/shared/auth/refresh
```

**请求头**
```
Authorization: Bearer <token>
```

**响应示例**
```json
{
  "code": 200,
  "message": "Token刷新成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**错误响应**
- `401` - Token无效或已过期
- `500` - 刷新失败

**JavaScript/Axios 示例**
```javascript
// Token 自动刷新拦截器
axios.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config;
    
    // 如果是401错误且未重试过
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        // 刷新 Token
        const response = await axios.post('/api/v1/shared/auth/refresh');
        const newToken = response.data.data.token;
        
        // 保存新 Token
        localStorage.setItem('token', newToken);
        axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`;
        
        // 重试原请求
        originalRequest.headers['Authorization'] = `Bearer ${newToken}`;
        return axios(originalRequest);
      } catch (refreshError) {
        // Token刷新失败，跳转登录
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }
    
    return Promise.reject(error);
  }
);
```

---

### 4.5 获取用户权限

**接口说明**: 获取当前用户的权限列表

**请求**
```
GET /api/v1/shared/auth/permissions
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取权限成功",
  "data": [
    "read:books",
    "write:documents",
    "manage:projects",
    "admin:users"
  ]
}
```

---

### 4.6 获取用户角色

**接口说明**: 获取当前用户的角色列表

**请求**
```
GET /api/v1/shared/auth/roles
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取角色成功",
  "data": [
    {
      "roleId": "role_123",
      "roleName": "author",
      "displayName": "作者",
      "permissions": ["write:documents", "manage:projects"]
    }
  ]
}
```

---

## 5. 管理员 API

### 5.1 获取待审核内容

**接口说明**: 获取待审核内容列表（管理员接口）

**请求**
```
GET /api/v1/shared/admin/reviews/pending
```

**Query 参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content_type | string | 否 | 内容类型（document/comment） |

**响应示例**
```json
{
  "code": 200,
  "message": "获取待审核内容成功",
  "data": [
    {
      "reviewId": "review_123",
      "contentId": "doc_123",
      "contentType": "document",
      "title": "待审核文档",
      "submittedBy": "user_123",
      "submittedAt": "2025-10-20T10:00:00Z"
    }
  ]
}
```

**错误响应**
- `401` - 未认证
- `403` - 无权限访问（非管理员）

---

### 5.2 审核内容

**接口说明**: 审核用户提交的内容（管理员接口）

**请求**
```
POST /api/v1/shared/admin/reviews
```

**请求体**
```json
{
  "content_id": "doc_123",
  "content_type": "document",
  "action": "approve",
  "reason": ""
}
```

**请求参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content_id | string | 是 | 内容ID |
| content_type | string | 是 | 内容类型 |
| action | string | 是 | 审核操作（approve/reject） |
| reason | string | 否 | 拒绝原因（action为reject时必填） |

**响应示例**
```json
{
  "code": 200,
  "message": "审核成功"
}
```

---

### 5.3 审核提现

**接口说明**: 审核用户提现申请（管理员接口）

**请求**
```
POST /api/v1/shared/admin/withdraw/review
```

**请求体**
```json
{
  "withdraw_id": "withdraw_123",
  "approved": true,
  "reason": ""
}
```

**响应示例**
```json
{
  "code": 200,
  "message": "批准提现成功"
}
```

---

### 5.4 获取用户统计

**接口说明**: 获取指定用户的统计信息（管理员接口）

**请求**
```
GET /api/v1/shared/admin/users/{user_id}/statistics
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取用户统计成功",
  "data": {
    "userId": "user_123",
    "totalProjects": 5,
    "totalDocuments": 50,
    "totalWords": 500000,
    "totalRevenue": 10000.00,
    "registeredAt": "2025-01-01T00:00:00Z"
  }
}
```

---

### 5.5 获取操作日志

**接口说明**: 获取管理员操作日志（管理员接口）

**请求**
```
GET /api/v1/shared/admin/operation-logs
```

**Query 参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |
| admin_id | string | 否 | - | 管理员ID |
| operation | string | 否 | - | 操作类型 |

**响应示例**
```json
{
  "code": 200,
  "message": "获取操作日志成功",
  "data": [
    {
      "logId": "log_123",
      "adminId": "admin_1",
      "adminName": "管理员",
      "operation": "approve_content",
      "targetType": "document",
      "targetId": "doc_123",
      "details": "批准文档发布",
      "createdAt": "2025-10-20T10:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

---

## 6. 完整示例代码

### 6.1 Vue 3 钱包管理组件

```vue
<template>
  <div class="wallet-container">
    <!-- 余额卡片 -->
    <div class="balance-card">
      <h3>我的余额</h3>
      <p class="balance">¥{{ balance.toFixed(2) }}</p>
      <div class="actions">
        <button @click="showRechargeDialog">充值</button>
        <button @click="showWithdrawDialog">提现</button>
      </div>
    </div>
    
    <!-- 交易记录 -->
    <div class="transaction-list">
      <h3>交易记录</h3>
      <div v-for="txn in transactions" :key="txn.transactionId" class="transaction-item">
        <div class="txn-info">
          <span class="txn-type">{{ getTxnTypeName(txn.type) }}</span>
          <span class="txn-desc">{{ txn.description }}</span>
        </div>
        <div class="txn-amount" :class="{ income: txn.amount > 0, expense: txn.amount < 0 }">
          {{ txn.amount > 0 ? '+' : '' }}{{ txn.amount.toFixed(2) }}
        </div>
        <div class="txn-date">{{ formatDate(txn.createdAt) }}</div>
      </div>
    </div>
    
    <!-- 充值对话框 -->
    <el-dialog v-model="rechargeDialogVisible" title="充值">
      <el-form :model="rechargeForm" label-width="100px">
        <el-form-item label="充值金额">
          <el-input v-model.number="rechargeForm.amount" type="number" placeholder="请输入充值金额" />
        </el-form-item>
        <el-form-item label="支付方式">
          <el-radio-group v-model="rechargeForm.method">
            <el-radio label="alipay">支付宝</el-radio>
            <el-radio label="wechat">微信</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rechargeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRecharge">确认充值</el-button>
      </template>
    </el-dialog>
    
    <!-- 提现对话框 -->
    <el-dialog v-model="withdrawDialogVisible" title="提现">
      <el-form :model="withdrawForm" label-width="100px">
        <el-form-item label="提现金额">
          <el-input v-model.number="withdrawForm.amount" type="number" placeholder="请输入提现金额" />
        </el-form-item>
        <el-form-item label="提现账号">
          <el-input v-model="withdrawForm.account" placeholder="支付宝账号或银行卡号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="withdrawDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleWithdraw">确认提现</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { ElMessage } from 'element-plus';

const balance = ref(0);
const transactions = ref([]);
const rechargeDialogVisible = ref(false);
const withdrawDialogVisible = ref(false);

const rechargeForm = ref({
  amount: 0,
  method: 'alipay'
});

const withdrawForm = ref({
  amount: 0,
  account: ''
});

// 获取余额
const fetchBalance = async () => {
  try {
    const response = await axios.get('/api/v1/shared/wallet/balance');
    balance.value = response.data.data.balance;
  } catch (error) {
    ElMessage.error('获取余额失败');
  }
};

// 获取交易记录
const fetchTransactions = async () => {
  try {
    const response = await axios.get('/api/v1/shared/wallet/transactions', {
      params: { page: 1, page_size: 20 }
    });
    transactions.value = response.data.data;
  } catch (error) {
    ElMessage.error('获取交易记录失败');
  }
};

// 充值
const handleRecharge = async () => {
  try {
    await axios.post('/api/v1/shared/wallet/recharge', {
      amount: rechargeForm.value.amount,
      method: rechargeForm.value.method
    });
    
    ElMessage.success('充值成功');
    rechargeDialogVisible.value = false;
    
    // 刷新余额和交易记录
    fetchBalance();
    fetchTransactions();
  } catch (error) {
    ElMessage.error('充值失败');
  }
};

// 提现
const handleWithdraw = async () => {
  try {
    await axios.post('/api/v1/shared/wallet/withdraw', {
      amount: withdrawForm.value.amount,
      account: withdrawForm.value.account
    });
    
    ElMessage.success('提现申请已提交');
    withdrawDialogVisible.value = false;
    
    // 刷新余额
    fetchBalance();
  } catch (error) {
    if (error.response?.status === 402) {
      ElMessage.error('余额不足');
    } else {
      ElMessage.error('提现失败');
    }
  }
};

const showRechargeDialog = () => {
  rechargeDialogVisible.value = true;
};

const showWithdrawDialog = () => {
  withdrawDialogVisible.value = true;
};

const getTxnTypeName = (type) => {
  const types = {
    recharge: '充值',
    consume: '消费',
    transfer: '转账',
    withdraw: '提现'
  };
  return types[type] || type;
};

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString('zh-CN');
};

onMounted(() => {
  fetchBalance();
  fetchTransactions();
});
</script>

<style scoped>
.wallet-container {
  padding: 2rem;
}

.balance-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 2rem;
  border-radius: 12px;
  margin-bottom: 2rem;
}

.balance {
  font-size: 3rem;
  font-weight: bold;
  margin: 1rem 0;
}

.actions {
  display: flex;
  gap: 1rem;
}

.actions button {
  padding: 0.5rem 1.5rem;
  background: white;
  color: #667eea;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.transaction-list {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.transaction-item {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  padding: 1rem 0;
  border-bottom: 1px solid #eee;
}

.txn-amount.income {
  color: #67C23A;
}

.txn-amount.expense {
  color: #F56C6C;
}
</style>
```

---

## 7. 最佳实践

### 7.1 文件上传最佳实践

```javascript
// 文件大小限制
const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB

// 支持的文件类型
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif'];

// 文件验证
const validateFile = (file) => {
  if (file.size > MAX_FILE_SIZE) {
    throw new Error('文件大小不能超过10MB');
  }
  
  if (!ALLOWED_TYPES.includes(file.type)) {
    throw new Error('不支持的文件类型');
  }
  
  return true;
};

// 带压缩的图片上传
const uploadImageWithCompression = async (file) => {
  // 验证文件
  validateFile(file);
  
  // 压缩图片
  const compressedFile = await compressImage(file, {
    maxWidth: 1920,
    maxHeight: 1080,
    quality: 0.8
  });
  
  // 上传
  return await uploadFile(compressedFile);
};
```

---

### 7.2 Token 管理最佳实践

```javascript
// Token 存储工具类
class TokenManager {
  static TOKEN_KEY = 'auth_token';
  static USER_ID_KEY = 'user_id';
  
  static saveToken(token, userId) {
    localStorage.setItem(this.TOKEN_KEY, token);
    localStorage.setItem(this.USER_ID_KEY, userId);
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
  }
  
  static getToken() {
    return localStorage.getItem(this.TOKEN_KEY);
  }
  
  static getUserId() {
    return localStorage.getItem(this.USER_ID_KEY);
  }
  
  static clearToken() {
    localStorage.removeItem(this.TOKEN_KEY);
    localStorage.removeItem(this.USER_ID_KEY);
    delete axios.defaults.headers.common['Authorization'];
  }
  
  static isAuthenticated() {
    return !!this.getToken();
  }
}

// 使用示例
const login = async (username, password) => {
  const response = await axios.post('/api/v1/shared/auth/login', {
    username,
    password
  });
  
  const { token, userId } = response.data.data;
  TokenManager.saveToken(token, userId);
  
  return response.data.data;
};
```

---

### 7.3 钱包操作安全实践

```javascript
// 操作确认
const confirmOperation = async (message) => {
  return new Promise((resolve) => {
    if (window.confirm(message)) {
      resolve(true);
    } else {
      resolve(false);
    }
  });
};

// 安全的充值流程
const safeRecharge = async (amount, method) => {
  // 1. 金额验证
  if (amount < 1 || amount > 10000) {
    throw new Error('充值金额必须在1-10000之间');
  }
  
  // 2. 确认操作
  const confirmed = await confirmOperation(`确认充值 ¥${amount} 吗？`);
  if (!confirmed) return;
  
  // 3. 执行充值
  try {
    const result = await axios.post('/api/v1/shared/wallet/recharge', {
      amount,
      method
    });
    
    // 4. 记录日志
    console.log('充值成功:', result.data.data.transactionId);
    
    return result.data.data;
  } catch (error) {
    // 5. 错误处理
    if (error.response?.data?.message) {
      throw new Error(error.response.data.message);
    }
    throw error;
  }
};
```

---

## 8. 常见问题

### Q1: Token过期怎么办？

**A**: 
1. 使用 Token 自动刷新机制（见 4.4 刷新Token）
2. 在拦截器中捕获401错误并自动刷新
3. 如果刷新失败，跳转到登录页

### Q2: 文件上传失败怎么办？

**A**:
1. 检查文件大小是否超限
2. 检查文件类型是否支持
3. 检查网络连接
4. 查看具体错误信息

### Q3: 余额不足怎么办？

**A**:
1. 先充值再消费
2. 检查冻结金额
3. 联系客服了解详情

### Q4: 提现多久到账？

**A**:
1. 提现申请提交后需要管理员审核
2. 审核通过后1-3个工作日到账
3. 可以在提现记录中查看状态

---

## 9. 相关文档

- [API 快速参考](./API快速参考.md) - 快速查找接口
- [前端集成指南](./前端集成指南.md) - 详细集成步骤
- [写作系统 API 参考](./写作系统API参考.md) - 项目、文档、编辑器等
- [用户系统 API 参考](./用户系统API参考.md) - 用户管理

---

**最后更新**: 2025-10-20  
**维护者**: 青羽后端团队  
**文档版本**: v1.0

