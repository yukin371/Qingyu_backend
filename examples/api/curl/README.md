# cURL 示例

本目录包含使用cURL调用青羽写作平台API的示例脚本。

## 前置要求

- 安装 [cURL](https://curl.se/)
- 安装 [jq](https://stedolan.github.io/jq/) (用于JSON格式化输出)
- 运行后端服务器 (默认地址: http://localhost:9090)

## 使用方法

### 认证示例

```bash
# 赋予执行权限
chmod +x auth.sh

# 运行脚本
./auth.sh
```

### 书城API示例

```bash
# 首先在auth.sh中获取Token，然后更新bookstore.sh中的TOKEN变量
# 编辑bookstore.sh，设置TOKEN变量为你的JWT token

# 赋予执行权限
chmod +x bookstore.sh

# 运行脚本
./bookstore.sh
```

### 管理员API示例

```bash
# 编辑admin.sh，设置ADMIN_TOKEN变量为你的管理员JWT token

# 赋予执行权限
chmod +x admin.sh

# 运行脚本
./admin.sh
```

## 示例说明

### auth.sh - 认证流程示例

演示完整的认证流程：
1. 用户注册
2. 用户登录
3. 使用Token访问受保护的API
4. 刷新Token
5. 登出

### bookstore.sh - 书城API示例

演示书城相关功能：
1. 获取首页数据
2. 获取书籍列表
3. 搜索书籍（按标题）
4. 按作者搜索
5. 按标签筛选
6. 获取书籍详情
7. 获取章节列表
8. 获取章节内容
9. 获取相似书籍推荐
10. 提交书籍评分

### admin.sh - 管理员API示例

演示管理员功能：
1. 获取用户列表
2. 创建用户
3. 获取权限列表
4. 创建权限
5. 获取权限模板列表
6. 创建权限模板
7. 获取审计日志
8. 获取用户增长趋势
9. 获取内容统计
10. 创建公告
11. 导出书籍数据

## 注意事项

1. **Token设置**: 运行脚本前，请确保已设置正确的Token
2. **服务器地址**: 默认使用 `http://localhost:9090`，如需修改请编辑脚本中的 `BASE_URL` 变量
3. **错误处理**: 脚本会显示HTTP状态码，请检查是否为200（成功）
4. **JSON格式化**: 需要安装jq工具来格式化JSON输出

## 常见问题

### 无法连接到服务器

检查后端服务器是否正在运行：

```bash
# 测试连接
curl http://localhost:9090/health
```

### 401 未授权错误

Token可能过期或无效，请重新登录获取新Token。

### jq命令未找到

安装jq工具：

```bash
# macOS
brew install jq

# Ubuntu/Debian
sudo apt-get install jq

# Windows (使用Chocolatey)
choco install jq
```

## 更多信息

- [API参考文档](../../../docs/api/reference.md)
- [错误代码](../../../docs/api/reference.md#错误代码)
