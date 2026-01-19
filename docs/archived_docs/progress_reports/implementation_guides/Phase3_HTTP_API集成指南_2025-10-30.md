# Phase3 HTTP API集成指南

**日期**: 2025-10-30  
**状态**: ✅ API层已完成，待初始化

---

## 📋 已完成的工作

✅ **API层完整实现**：
- `api/v1/ai/creative_models.go` - 请求/响应模型
- `api/v1/ai/creative_api.go` - API处理器
- `api/v1/ai/creative_converters.go` - 数据转换
- `router/ai/creative.go` - 路由注册
- `doc/api/Phase3创作API文档.md` - API文档

✅ **路由已集成**：
- `router/ai/ai_router.go` - 已更新支持Phase3
- `router/enter.go` - 已集成Phase3Client获取

✅ **服务容器支持**：
- `service/container/service_container.go` - 已添加Phase3Client字段和获取方法

---

## 🚀 初始化步骤

### 步骤1: 配置gRPC服务地址

在 `config/config.yaml` 中添加配置：

```yaml
# AI服务配置
ai:
  # Phase3 gRPC服务配置
  phase3:
    enabled: true
    grpc_address: "localhost:50051"
    timeout: 120 # 秒
```

### 步骤2: 在服务容器中初始化Phase3Client

修改 `service/container/service_container.go` 的初始化方法：

在 `SetupDefaultServices()` 方法中添加：

```go
// 初始化Phase3 gRPC客户端（如果配置启用）
phase3Enabled := viper.GetBool("ai.phase3.enabled")
if phase3Enabled {
    grpcAddr := viper.GetString("ai.phase3.grpc_address")
    if grpcAddr == "" {
        grpcAddr = "localhost:50051"
    }
    
    logger.Info("初始化Phase3 gRPC客户端", zap.String("address", grpcAddr))
    
    phase3Client, err := aiService.NewPhase3Client(grpcAddr)
    if err != nil {
        logger.Warn("Phase3 gRPC客户端初始化失败", zap.Error(err))
        c.phase3Client = nil
    } else {
        c.phase3Client = phase3Client
        logger.Info("✅ Phase3 gRPC客户端初始化成功")
    }
} else {
    logger.Info("Phase3服务未启用")
}
```

### 步骤3: 启动Python gRPC服务

在启动Go后端之前，先启动Python AI服务：

```powershell
# 终端1
cd python_ai_service
$env:GOOGLE_API_KEY="your_api_key_here"
python run_grpc_server.py
```

### 步骤4: 启动Go后端

```bash
go run cmd/server/main.go
```

查看日志确认Phase3路由已注册：

```
✓ AI服务路由已注册到: /api/v1/ai/
  - /api/v1/ai/writing/* (续写、改写)
  - /api/v1/ai/chat/* (聊天)
  - /api/v1/ai/quota/* (配额管理)
  - /api/v1/ai/creative/* (Phase3创作工作流)  ← 应该看到这行
```

---

## 🧪 测试API

### 方式1: 使用curl

```bash
# 1. 登录获取token
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'

# 2. 健康检查（无需token）
curl http://localhost:8080/api/v1/ai/creative/health

# 3. 生成大纲
curl -X POST http://localhost:8080/api/v1/ai/creative/outline \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "task": "创作一个修仙小说大纲，主角是天才少年"
  }'

# 4. 执行完整工作流
curl -X POST http://localhost:8080/api/v1/ai/creative/workflow \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "task": "创作一个都市爱情小说设定",
    "max_reflections": 3
  }'
```

### 方式2: 使用Postman/Apifox

1. 导入API文档: `doc/api/Phase3创作API文档.md`
2. 设置Authorization为Bearer Token
3. 调用各个接口测试

### 方式3: 前端调用

```javascript
// 生成大纲
async function generateOutline() {
  const response = await fetch('/api/v1/ai/creative/outline', {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer ' + token,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      task: '创作一个科幻小说大纲'
    })
  });
  
  const data = await response.json();
  console.log('大纲:', data.data.outline);
}
```

---

## 📁 文件清单

### API层（3个文件）

1. **api/v1/ai/creative_models.go** (~300行)
   - 所有请求/响应模型定义
   - OutlineData, CharactersData, PlotData等

2. **api/v1/ai/creative_api.go** (~200行)
   - 5个API处理器
   - 数据验证和转换
   - gRPC客户端调用

3. **api/v1/ai/creative_converters.go** (~300行)
   - Proto ↔ Model 双向转换
   - 类型安全的数据转换

### 路由层（2个文件）

4. **router/ai/creative.go** (~30行)
   - Phase3创作路由注册
   - 认证中间件配置

5. **router/ai/ai_router.go** (已修改)
   - 集成Phase3路由

6. **router/enter.go** (已修改)
   - 从服务容器获取Phase3Client

### 服务容器（1个文件）

7. **service/container/service_container.go** (已修改)
   - 添加phase3Client字段
   - 实现GetPhase3Client()方法

### 文档（2个文件）

8. **doc/api/Phase3创作API文档.md**
   - 完整的API文档
   - 请求/响应示例
   - 前端集成示例

9. **doc/implementation/.../Phase3_HTTP_API集成指南_2025-10-30.md** (本文件)
   - 集成指南

**总计**: ~900行代码 + 完整文档

---

## 🎯 API路由结构

```
/api/v1/ai/creative/
├── GET  /health                    # 健康检查（公开）
├── POST /outline                   # 生成大纲（需认证）
├── POST /characters                # 生成角色（需认证）
├── POST /plot                      # 生成情节（需认证）
└── POST /workflow                  # 完整工作流（需认证）
```

---

## 📊 数据流

```
前端请求
    ↓
HTTP API (/api/v1/ai/creative/*)
    ↓
API Handler (creative_api.go)
    ↓
Data Converter (Model → Proto)
    ↓
Phase3Client (gRPC客户端)
    ↓
Python gRPC服务 (localhost:50051)
    ↓
Phase3 Agents (Outline/Character/Plot)
    ↓
Gemini 2.0 Flash API
    ↓
← 返回结果
    ↓
Data Converter (Proto → Model)
    ↓
← HTTP响应
```

---

## 🔧 配置参考

完整的配置文件示例：

```yaml
# config/config.yaml

server:
  port: 8080
  mode: debug

ai:
  # 现有AI配置
  providers:
    - name: gemini
      api_key: ${GEMINI_API_KEY}
      model: gemini-pro
  
  # Phase3配置（新增）
  phase3:
    enabled: true                    # 是否启用Phase3
    grpc_address: "localhost:50051"  # gRPC服务地址
    timeout: 120                     # 超时时间（秒）
    max_retries: 3                   # 最大重试次数
    
database:
  uri: ${MONGODB_URI}
  
jwt:
  secret: ${JWT_SECRET}
  expire: 7200
```

环境变量：

```bash
export GOOGLE_API_KEY="your_gemini_api_key"
export MONGODB_URI="mongodb://localhost:27017/qingyu"
export JWT_SECRET="your_jwt_secret"
```

---

## 🐛 故障排查

### 问题1: Phase3路由未显示

**症状**:
```
✓ AI服务路由已注册到: /api/v1/ai/
  - /api/v1/ai/writing/* (续写、改写)
  # 没有 /api/v1/ai/creative/*
```

**解决**:
1. 检查`GetPhase3Client()`是否返回错误
2. 确认Python gRPC服务是否启动
3. 检查gRPC地址配置是否正确

### 问题2: API调用失败

**症状**:
```json
{
  "code": 500,
  "message": "大纲生成失败",
  "data": {
    "error": "connection refused"
  }
}
```

**解决**:
1. 确认Python gRPC服务正在运行
2. 检查端口50051是否被占用
3. 查看Go后端日志

### 问题3: 超时错误

**症状**:
```
context deadline exceeded
```

**解决**:
1. 增加超时时间（修改`phase3_client.go`）
2. 检查AI服务响应速度
3. 查看Python服务日志

---

## 📚 相关文档

- [Phase3 Go客户端](Phase3_Go集成完成总结_2025-10-30.md)
- [Phase3 gRPC集成](Phase3_gRPC集成完成报告_2025-10-30.md)
- [API文档](../../api/Phase3创作API文档.md)
- [快速开始](../../../PHASE3_QUICK_START_GO.md)

---

## ✅ 检查清单

### 代码层面
- [x] API模型定义完成
- [x] API处理器实现完成
- [x] 数据转换器完成
- [x] 路由注册完成
- [x] 服务容器支持完成

### 配置层面
- [ ] 添加config.yaml配置
- [ ] 服务容器初始化Phase3Client
- [ ] 环境变量配置

### 测试层面
- [ ] Python gRPC服务启动
- [ ] Go后端启动
- [ ] API健康检查通过
- [ ] 各个接口测试通过

### 文档层面
- [x] API文档完成
- [x] 集成指南完成
- [x] 使用示例完成

---

**状态**: 🎯 **API层完成，待配置和测试**

**下一步**:
1. 修改配置文件
2. 初始化Phase3Client
3. 启动服务测试
4. 前端对接

---

**完成时间**: 2025-10-30  
**维护者**: 青羽后端架构团队

