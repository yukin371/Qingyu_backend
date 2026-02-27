# API 示例代码

本目录包含青羽写作平台API的多种语言示例代码，帮助开发者快速集成和使用API。

## 目录结构

```
examples/api/
├── curl/           # Shell/cURL示例
├── python/         # Python示例
├── javascript/     # JavaScript/Node.js示例
└── README.md       # 本文件
```

## 快速开始

### 1. cURL 示例

适用于快速测试和Shell脚本集成。

**前置要求**:
- cURL
- jq (可选，用于格式化JSON输出)

```bash
cd curl
./auth.sh
```

详见: [curl/README.md](curl/README.md)

### 2. Python 示例

适用于Python后端集成和自动化脚本。

**前置要求**:
- Python 3.7+
- requests库

```bash
cd python
pip install -r requirements.txt
python auth_example.py
```

详见: [python/README.md](python/README.md)

### 3. JavaScript/Node.js 示例

适用于Node.js后端集成和前端开发参考。

**前置要求**:
- Node.js 14+
- npm/yarn

```bash
cd javascript
npm install
npm run auth
```

详见: [javascript/README.md](javascript/README.md)

## API分类

### 认证API

所有示例都包含完整的认证流程演示：
- 用户注册
- 用户登录
- Token管理
- 登出

**文件位置**:
- cURL: `curl/auth.sh`
- Python: `python/auth_example.py`
- JavaScript: `javascript/auth_example.js`

### 书城API

演示书城相关功能：
- 首页数据获取
- 书籍搜索和筛选
- 章节内容获取
- 评分功能

**文件位置**:
- cURL: `curl/bookstore.sh`
- Python: `python/bookstore_example.py`
- JavaScript: `javascript/bookstore_example.js`

### 管理员API

演示管理员功能：
- 用户管理
- 权限管理
- 统计分析
- 公告管理

**文件位置**:
- cURL: `curl/admin.sh`
- Python: `python/admin_example.py`
- JavaScript: `javascript/admin_example.js`

## 通用配置

### 服务器地址

默认使用 `http://localhost:9090/api/v1`，可通过以下方式修改：

**cURL**:
```bash
# 编辑脚本，修改BASE_URL变量
BASE_URL="http://your-server:port/api/v1"
```

**Python**:
```python
client = QingyuAPIClient(base_url="http://your-server:port/api/v1")
```

**JavaScript**:
```javascript
const client = new QingyuAPIClient('http://your-server:port/api/v1');
```

### Token设置

大多数API需要JWT认证。请先通过认证API获取Token：

```bash
# 登录获取Token
curl -X POST http://localhost:9090/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"your_username","password":"your_password"}'
```

然后设置Token：

**环境变量方式** (推荐):
```bash
export QINGYU_TOKEN=your_jwt_token_here
export QINGYU_ADMIN_TOKEN=your_admin_token_here
```

**直接修改脚本方式**:
编辑相应的示例文件，设置Token变量。

## 错误处理

所有示例都包含基本的错误处理：

### HTTP状态码

- 200: 成功
- 400: 请求参数错误
- 401: 未授权/Token无效
- 403: 权限不足
- 404: 资源不存在
- 500: 服务器错误

### 错误响应格式

```json
{
  "code": 400,
  "message": "错误描述",
  "details": "详细错误信息"
}
```

## 开发建议

### 生产环境

1. **使用HTTPS**: 生产环境必须使用HTTPS
2. **Token安全**: 不要在代码中硬编码Token，使用环境变量或配置文件
3. **错误重试**: 实现适当的重试机制
4. **日志记录**: 添加详细的日志记录
5. **限流处理**: 遵守API速率限制

### 测试环境

1. 使用测试账号进行测试
2. 不要在生产环境运行测试代码
3. 及时清理测试数据

## 示例用法

### Python - 获取书籍列表

```python
from python.bookstore_example import BookstoreAPIClient

# 创建客户端
client = BookstoreAPIClient(token="your_token")

# 获取书籍列表
books = client.get_books(page=1, limit=10)
print(books)
```

### JavaScript - 创建用户

```javascript
const { AdminAPIClient } = require('./javascript/api_client');

const client = new AdminAPIClient('http://localhost:9090/api/v1', 'admin_token');

async function createUser() {
  await client.createUser(
    'newuser',
    'user@example.com',
    'password123',
    '新用户'
  );
}

createUser();
```

### cURL - 搜索书籍

```bash
curl -X GET "http://localhost:9090/api/v1/bookstore/books/search/title?title=玄幻&page=1&limit=5" \
  -H "Authorization: Bearer your_token"
```

## 故障排除

### 无法连接到服务器

1. 检查服务器是否正在运行
2. 检查防火墙设置
3. 验证服务器地址和端口

### 401 未授权错误

1. 检查Token是否正确
2. 确认Token未过期
3. 尝试重新登录获取新Token

### 依赖安装失败

**Python**:
```bash
# 使用国内镜像
pip install -i https://pypi.tuna.tsinghua.edu.cn/simple -r requirements.txt
```

**Node.js**:
```bash
# 使用国内镜像
npm config set registry https://registry.npmmirror.com
npm install
```

## 更多资源

- [API参考文档](../../docs/api/reference.md)
- [错误代码说明](../../docs/api/reference.md#错误代码)
- [完整API文档](../../docs/swagger.yaml)

## 反馈和贡献

如有问题或建议，请：
1. 查看现有文档
2. 提交Issue
3. 提交Pull Request

## 许可证

MIT License
