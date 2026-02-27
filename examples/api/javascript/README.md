# JavaScript/Node.js 示例

本目录包含使用JavaScript/Node.js调用青羽写作平台API的示例代码。

## 前置要求

- Node.js 14+
- npm 或 yarn

## 安装依赖

```bash
npm install
```

或使用yarn:

```bash
yarn install
```

## 依赖包

- axios: HTTP客户端
- 其他依赖见package.json

## 使用方法

### 认证示例

```bash
# 运行认证示例
node auth_example.js
```

**功能说明**:
- 用户注册
- 用户登录
- 获取用户信息
- 刷新Token
- 登出

### 书城API示例

```bash
# 设置Token环境变量 (Linux/Mac)
export QINGYU_TOKEN=your_jwt_token_here

# 设置Token环境变量 (Windows)
set QINGYU_TOKEN=your_jwt_token_here

# 运行书城API示例
node bookstore_example.js
```

**功能说明**:
- 获取首页数据
- 获取书籍列表
- 按标题搜索
- 按作者搜索
- 按标签筛选
- 获取书籍详情
- 获取章节列表
- 获取章节内容
- 获取相似书籍
- 提交评分

### 管理员API示例

```bash
# 设置管理员Token环境变量 (Linux/Mac)
export QINGYU_ADMIN_TOKEN=your_admin_jwt_token_here

# 设置Token环境变量 (Windows)
set QINGYU_ADMIN_TOKEN=your_admin_jwt_token_here

# 运行管理员API示例
node admin_example.js
```

**功能说明**:
- 用户管理（列表、创建）
- 权限管理（列表、创建）
- 统计分析（用户增长、内容统计、系统概览）
- 公告管理

## 代码结构

### QingyuAPIClient (基础客户端)

提供基础的API客户端功能：
- 自动添加认证头
- 统一的错误处理
- 请求/响应拦截器

### BookstoreAPIClient (书城客户端)

继承自基础客户端，提供书城相关API：
- 书籍搜索和筛选
- 章节获取
- 评分功能

### AdminAPIClient (管理员客户端)

继承自基础客户端，提供管理员相关API：
- 用户管理
- 权限管理
- 统计分析
- 公告管理

## 示例代码

### 基本使用

```javascript
const { QingyuAPIClient } = require('./auth_example');

// 创建客户端
const client = new QingyuAPIClient();

// 登录
async function example() {
  await client.login('testuser', 'password123');

  // 访问受保护的API
  if (client.token) {
    const profile = await client.getProfile();
    console.log(profile);
  }
}

example();
```

### 自定义配置

```javascript
// 自定义API地址
const client = new QingyuAPIClient('https://api.example.com/api/v1');

// 使用现有Token创建客户端
const { BookstoreAPIClient } = require('./api_client');
const bookstoreClient = new BookstoreAPIClient('http://localhost:9090/api/v1', 'your_jwt_token');
```

### 错误处理

所有API调用都包含错误处理：

```javascript
try {
  await client.getBooks();
} catch (error) {
  console.error('获取书籍失败:', error.message);
  // 处理错误
}
```

## 在浏览器中使用

虽然这些示例是为Node.js设计的，但代码也可以在浏览器中使用：

```html
<!DOCTYPE html>
<html>
<head>
  <title>API示例</title>
  <script src="https://cdn.jsdelivr.net/npm/axios/dist/axios.min.js"></script>
</head>
<body>
  <script type="module">
    // 由于CORS限制，浏览器中直接调用API可能需要配置代理
    // 建议使用后端SDK或通过后端代理调用API
  </script>
</body>
</html>
```

## TypeScript支持

虽然示例是JavaScript，但可以很容易地转换为TypeScript。类型定义可以参考：

```typescript
interface APIResponse<T = any> {
  code: number;
  message: string;
  data?: T;
  timestamp?: number;
  request_id?: string;
}

interface Book {
  id: string;
  title: string;
  author: string;
  // ... 其他字段
}
```

## 注意事项

1. **Token管理**: Token有时效性，过期后需要重新登录
2. **CORS**: 浏览器中直接调用需要注意CORS配置
3. **环境变量**: 敏感信息应使用环境变量
4. **错误处理**: 建议在实际应用中添加更完善的错误处理

## 更多信息

- [API参考文档](../../../docs/api/reference.md)
- [axios文档](https://axios-http.com/docs/intro)
