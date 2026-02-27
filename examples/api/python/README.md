# Python 示例

本目录包含使用Python调用青羽写作平台API的示例代码。

## 前置要求

- Python 3.7+
- requests库

## 安装依赖

```bash
pip install requests
```

或使用requirements.txt:

```bash
pip install -r requirements.txt
```

## 使用方法

### 认证示例

```bash
# 运行认证示例
python auth_example.py
```

**功能说明**:
- 用户注册
- 用户登录
- 获取用户信息
- 刷新Token
- 登出

### 书城API示例

```bash
# 设置Token环境变量
export QINGYU_TOKEN=your_jwt_token_here

# 运行书城API示例
python bookstore_example.py
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
# 设置管理员Token环境变量
export QINGYU_ADMIN_TOKEN=your_admin_jwt_token_here

# 运行管理员API示例
python admin_example.py
```

**功能说明**:
- 用户管理（列表、创建、更新、删除）
- 权限管理（列表、创建）
- 权限模板管理（列表、创建、应用）
- 审计日志查询
- 统计分析（用户增长、内容统计、系统概览）
- 公告管理

## 代码结构

### QingyuAPIClient / BookstoreAPIClient / AdminAPIClient

这些类封装了API调用逻辑，提供了：

- 统一的请求处理
- 自动添加认证头
- 错误处理
- 类型提示

### 方法说明

每个方法对应一个API端点，参数和返回值都使用了Python类型提示。

## 错误处理

所有方法都包含错误处理，当请求失败时会：

1. 打印错误信息
2. 打印响应内容（如果有）
3. 抛出异常

## 示例代码

### 基本使用

```python
from auth_example import QingyuAPIClient

# 创建客户端
client = QingyuAPIClient()

# 登录
result = client.login(username="testuser", password="password123")

# 访问受保护的API
if client.token:
    profile = client.get_profile()
```

### 自定义配置

```python
# 自定义API地址
client = QingyuAPIClient(base_url="https://api.example.com/api/v1")

# 使用现有Token创建客户端
client = BookstoreAPIClient(token="your_jwt_token")
```

## 注意事项

1. **Token管理**: Token有时效性，过期后需要重新登录或使用refresh_token
2. **错误处理**: 建议在实际应用中添加更完善的错误处理逻辑
3. **并发请求**: 如需并发请求，建议使用连接池或异步客户端
4. **环境变量**: 敏感信息（如Token）应使用环境变量，不要硬编码

## 更多信息

- [API参考文档](../../../docs/api/reference.md)
- [Python requests文档](https://docs.python-requests.org/)
