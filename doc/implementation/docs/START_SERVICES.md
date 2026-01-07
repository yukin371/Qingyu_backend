# 服务启动指南

## 快速启动服务

### 方法一：使用 PowerShell 脚本（推荐）

#### 1. 启动后端服务

打开第一个 PowerShell 终端：

```powershell
cd "D:\Github\青羽\Qingyu_backend"
go run main.go
```

等待看到类似以下的输出信息：
```
[GIN-debug] Listening and serving HTTP on :8080
数据库连接成功
```

#### 2. 启动前端服务

打开第二个 PowerShell 终端（保持第一个终端运行）：

```powershell
cd "D:\Github\青羽\Qingyu"
npm run dev
```

等待看到类似以下的输出：
```
  VITE v5.x.x  ready in xxx ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

### 方法二：使用批处理脚本

创建 `start-backend.bat`：
```batch
@echo off
cd /d "D:\Github\青羽\Qingyu_backend"
echo Starting backend service...
go run main.go
pause
```

创建 `start-frontend.bat`：
```batch
@echo off
cd /d "D:\Github\青羽\Qingyu"
echo Starting frontend service...
npm run dev
pause
```

双击这两个批处理文件即可启动服务。

## 访问应用

### 前端应用
- **主页**: http://localhost:5173/
- **Shared API 测试页面**: http://localhost:5173/shared-api-test

### 后端API
- **基础URL**: http://localhost:8080/api/v1
- **Shared Auth**: http://localhost:8080/api/v1/shared/auth/
- **Shared Wallet**: http://localhost:8080/api/v1/shared/wallet/
- **Shared Storage**: http://localhost:8080/api/v1/shared/storage/
- **Shared Admin**: http://localhost:8080/api/v1/shared/admin/

## 验证服务状态

### 检查后端服务

在浏览器中访问或使用 PowerShell：
```powershell
Invoke-WebRequest -Uri "http://localhost:8080" -Method GET
```

### 检查前端服务

在浏览器中访问：
```
http://localhost:5173
```

## 测试 Shared API

1. 访问测试页面: http://localhost:5173/shared-api-test

2. 按照以下顺序测试：

   **a. 用户注册**
   - 填写用户名、邮箱、密码
   - 点击"测试注册"

   **b. 用户登录**
   - 填写用户名、密码
   - 点击"测试登录"
   - Token 会自动保存

   **c. 钱包操作**
   - 点击"查询余额"
   - 测试充值、消费等功能

   **d. 文件存储**
   - 上传测试文件
   - 查看文件列表

3. 查看测试结果
   - 所有操作结果会显示在页面底部
   - 绿色表示成功，红色表示失败

## 常见问题

### 1. 后端启动失败

**检查 MongoDB**
```powershell
netstat -ano | findstr "27017"
```

如果没有输出，需要启动 MongoDB：
```powershell
# 方法1：使用服务
net start MongoDB

# 方法2：直接启动
mongod --dbpath "C:\data\db"
```

**检查端口占用**
```powershell
netstat -ano | findstr "8080"
```

如果端口被占用，可以在 `config.yaml` 中修改端口。

### 2. 前端启动失败

**检查 Node.js 版本**
```powershell
node --version
npm --version
```

推荐使用 Node.js v18 或更高版本。

**重新安装依赖**
```powershell
cd "D:\Github\青羽\Qingyu"
rm -r node_modules
npm install
```

### 3. 前后端连接失败

**检查 API 基础URL**

查看 `Qingyu/.env` 文件或 `vite.config.js`：
```javascript
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

**检查 CORS 配置**

后端应该已经配置了 CORS 允许前端访问。

### 4. 401 未授权错误

先进行登录操作，获取 Token 后再测试其他需要认证的 API。

## 停止服务

### 停止后端
在后端终端中按 `Ctrl + C`

### 停止前端
在前端终端中按 `Ctrl + C`

## 开发提示

### 实时查看日志

**后端日志**
- 后端终端会实时显示请求日志
- 包括 API 调用、数据库操作等

**前端日志**
- 打开浏览器开发者工具 (F12)
- 查看 Console 标签页
- 查看 Network 标签页监控 API 请求

### 热重载

- **后端**: 需要手动重启 (或使用 `air` 工具实现热重载)
- **前端**: Vite 自动热重载，修改代码后自动刷新

## 生产部署

生产环境部署请参考：
- 后端: `Qingyu_backend/README.md`
- 前端: `Qingyu/README.md`

## 技术支持

如遇问题，请检查：
1. 后端终端的错误信息
2. 前端终端的错误信息
3. 浏览器控制台的错误信息
4. MongoDB 是否正常运行

## 相关文档

- [Shared API 快速启动](Qingyu/SHARED_API_QUICKSTART.md)
- [Shared API 测试指南](Qingyu/doc/shared-api-test-guide.md)
- [API 使用文档](Qingyu/src/api/shared/README.md)

