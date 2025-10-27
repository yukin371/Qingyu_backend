# Viper配置管理机制详解

## 📋 概述

本项目使用 **Viper** 作为配置管理工具，支持多种配置来源和优先级。**不使用 `.env` 文件**，而是采用 **YAML配置文件 + 环境变量** 的方式。

**核心特性**：
- ✅ 支持YAML配置文件
- ✅ 支持环境变量覆盖
- ✅ 支持配置热重载
- ✅ 支持多环境配置
- ❌ 不依赖 `.env` 文件

---

## 🔧 Viper在项目中的使用

### 1. 配置加载流程

```
启动应用
    ↓
检查 CONFIG_FILE 环境变量
    ↓
    ├─有：直接加载指定文件
    │     (如: export CONFIG_FILE=./config/config.test.yaml)
    │
    └─无：按顺序搜索配置文件
          1. ./config/config.yaml
          2. ../../config/config.yaml (从cmd/server运行时)
          3. ./config.yaml
          4. /app/config/config.yaml (Docker环境)
    ↓
读取YAML文件内容
    ↓
应用默认值 (setDefaults)
    ↓
环境变量覆盖 (AutomaticEnv)
    ↓
验证配置 (ValidateConfig)
    ↓
设置全局配置 (GlobalConfig)
    ↓
配置加载完成
```

### 2. 核心代码实现

#### config/config.go - 主配置加载

```go
func LoadConfig(configPath string) (*Config, error) {
    v = viper.New()

    // 1. 设置默认值
    setDefaults()

    // 2. 配置文件路径
    if strings.HasSuffix(configPath, ".yaml") {
        v.SetConfigFile(configPath) // 直接指定文件
    } else {
        v.SetConfigName("config")
        v.SetConfigType("yaml")
        v.AddConfigPath(configPath)
        v.AddConfigPath("./config")
        v.AddConfigPath("../../config")
        v.AddConfigPath(".")
    }

    // 3. 环境变量支持
    v.AutomaticEnv()                                   // 自动读取环境变量
    v.SetEnvPrefix("QINGYU")                           // 前缀: QINGYU_
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // . → _

    // 4. 读取配置文件
    if err := v.ReadInConfig(); err != nil {
        // 处理错误
    }

    // 5. 解析到结构体
    config := &Config{}
    v.Unmarshal(config)

    // 6. 验证配置
    ValidateConfig(config)

    // 7. 设置全局配置
    GlobalConfig = config

    return config, nil
}
```

### 3. 默认值设置

```go
func setDefaults() {
    // 数据库默认配置
    v.SetDefault("database.primary.type", "mongodb")
    v.SetDefault("database.primary.mongodb.uri", "mongodb://localhost:27017")
    v.SetDefault("database.primary.mongodb.database", "qingyu")
    v.SetDefault("database.primary.mongodb.max_pool_size", 100)

    // 服务器默认配置
    v.SetDefault("server.port", "8080")
    v.SetDefault("server.mode", "debug")

    // JWT默认配置
    v.SetDefault("jwt.secret", "qingyu_secret_key")
    v.SetDefault("jwt.expiration_hours", 24)

    // AI默认配置
    v.SetDefault("ai.base_url", "https://api.openai.com/v1")
    v.SetDefault("ai.max_tokens", 2000)
}
```

### 4. 配置热重载

```go
func WatchConfig(onChange func()) {
    v.WatchConfig()
    v.OnConfigChange(func(e fsnotify.Event) {
        fmt.Printf("Config file changed: %s\n", e.Name)

        // 重新加载配置
        config := &Config{}
        v.Unmarshal(config)

        // 验证配置
        ValidateConfig(config)

        // 更新全局配置
        GlobalConfig = config

        // 调用回调
        onChange()
    })
}
```

---

## 🎯 配置优先级

Viper的配置优先级（从高到低）：

```
1. 环境变量 (最高优先级)
   ↓
2. 配置文件中的值
   ↓
3. 默认值 (最低优先级)
```

### 优先级示例

假设有以下配置：

**config.yaml**:
```yaml
server:
  port: "8080"
  mode: "debug"
```

**环境变量**:
```bash
export QINGYU_SERVER_PORT="9090"
```

**默认值**:
```go
v.SetDefault("server.port", "8080")
v.SetDefault("server.mode", "debug")
```

**最终结果**:
```yaml
server:
  port: "9090"  # 环境变量覆盖
  mode: "debug" # 配置文件值
```

---

## 🌍 环境变量覆盖机制

### 1. 命名规则

**YAML配置键** → **环境变量名**

| YAML路径 | 环境变量 |
|---------|---------|
| `server.port` | `QINGYU_SERVER_PORT` |
| `database.primary.mongodb.uri` | `QINGYU_DATABASE_PRIMARY_MONGODB_URI` |
| `jwt.secret` | `QINGYU_JWT_SECRET` |
| `ai.api_key` | `QINGYU_AI_API_KEY` |

**转换规则**：
1. 添加前缀 `QINGYU_`
2. 所有字母大写
3. `.` 替换为 `_`

### 2. 使用示例

#### 开发环境
```bash
# 使用默认配置
./qingyu_backend
```

#### 测试环境
```bash
# 指定测试配置文件
export CONFIG_FILE="./config/config.test.yaml"
./qingyu_backend
```

#### 生产环境
```bash
# 使用生产配置 + 环境变量覆盖敏感信息
export QINGYU_JWT_SECRET="super-secure-production-secret"
export QINGYU_DATABASE_PRIMARY_MONGODB_URI="mongodb://prod-user:prod-pass@prod-host:27017"
export QINGYU_SERVER_MODE="release"
./qingyu_backend
```

#### Docker环境
```yaml
# docker-compose.yml
services:
  backend:
    image: qingyu-backend:latest
    environment:
      - CONFIG_FILE=/app/config/config.docker.yaml
      - QINGYU_JWT_SECRET=production-secret
      - QINGYU_SERVER_MODE=release
    volumes:
      - ./config:/app/config
```

---

## ❌ 关于 .env 文件

### 项目中不使用 .env 文件

本项目 **不依赖 `.env` 文件**，原因如下：

#### 1. Viper原生支持环境变量

Viper的 `AutomaticEnv()` 方法可以直接读取系统环境变量，无需额外的 `.env` 文件。

```go
v.AutomaticEnv()  // 自动读取所有环境变量
v.SetEnvPrefix("QINGYU")  // 只读取 QINGYU_ 前缀的
```

#### 2. .gitignore 中已忽略

```gitignore
# 环境变量文件
.env
.env.local
.env.development.local
.env.test.local
.env.production.local
```

#### 3. 使用YAML配置更清晰

**对比**：

**.env 方式**（不推荐）：
```env
# .env
DATABASE_URI=mongodb://localhost:27017
DATABASE_NAME=qingyu
JWT_SECRET=secret
SERVER_PORT=8080
```

**YAML方式**（推荐）：
```yaml
# config.yaml
database:
  primary:
    mongodb:
      uri: mongodb://localhost:27017
      database: qingyu

jwt:
  secret: qingyu_secret_key

server:
  port: "8080"
```

**优势**：
- ✅ 层次结构清晰
- ✅ 支持注释说明
- ✅ 支持复杂数据结构
- ✅ IDE语法高亮支持
- ✅ 类型提示更好

#### 4. 环境变量优先级更高

即使不使用 `.env` 文件，环境变量仍然可以覆盖配置文件：

```bash
# 直接设置环境变量（无需.env文件）
export QINGYU_JWT_SECRET="production-secret"
export QINGYU_SERVER_PORT="9090"
./qingyu_backend
```

---

## 🔄 多环境配置方案

### 方案1: 多配置文件（推荐）

为不同环境创建独立的配置文件：

```
config/
├── config.yaml          # 默认配置（开发）
├── config.test.yaml     # 测试配置
├── config.docker.yaml   # Docker配置
└── config.prod.yaml     # 生产配置（不提交到Git）
```

**使用方式**：

```bash
# 开发环境
./qingyu_backend

# 测试环境
export CONFIG_FILE="./config/config.test.yaml"
./qingyu_backend

# 生产环境
export CONFIG_FILE="./config/config.prod.yaml"
./qingyu_backend
```

### 方案2: 配置文件 + 环境变量

基于一个基础配置文件，通过环境变量覆盖特定值：

```bash
# 使用开发配置，但覆盖数据库和密钥
export CONFIG_FILE="./config/config.yaml"
export QINGYU_DATABASE_PRIMARY_MONGODB_URI="mongodb://prod:27017"
export QINGYU_JWT_SECRET="production-secret"
./qingyu_backend
```

### 方案3: 完全环境变量（适合容器）

在容器环境中，完全通过环境变量配置：

```yaml
# docker-compose.yml
services:
  backend:
    environment:
      - QINGYU_SERVER_PORT=8080
      - QINGYU_SERVER_MODE=release
      - QINGYU_DATABASE_PRIMARY_MONGODB_URI=mongodb://mongodb:27017
      - QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE=qingyu_prod
      - QINGYU_JWT_SECRET=${JWT_SECRET}
      - QINGYU_AI_API_KEY=${AI_API_KEY}
```

---

## 📊 配置来源对比

| 配置方式 | 优先级 | 适用场景 | 推荐度 |
|---------|-------|---------|--------|
| **YAML配置文件** | 中 | 基础配置、默认值 | ⭐⭐⭐⭐⭐ |
| **环境变量** | 高 | 敏感信息、环境特定值 | ⭐⭐⭐⭐⭐ |
| **.env文件** | 低 | 本地开发（可选） | ⭐⭐ |
| **默认值** | 最低 | 兜底配置 | ⭐⭐⭐⭐ |

---

## 🛠️ 实际使用指南

### 场景1: 本地开发

```bash
# 1. 使用默认配置
cd Qingyu_backend
go run cmd/server/main.go

# 配置加载顺序：
# ./config/config.yaml → 默认值
```

### 场景2: 运行测试

```bash
# 方式1: 指定测试配置
export CONFIG_FILE="./config/config.test.yaml"
go test ./test/...

# 方式2: 在测试代码中指定
_, err := config.LoadConfig("../../config/config.test.yaml")
```

### 场景3: Docker部署

```yaml
# docker-compose.yml
services:
  backend:
    build: .
    environment:
      - CONFIG_FILE=/app/config/config.docker.yaml
      - QINGYU_JWT_SECRET=${JWT_SECRET}  # 从宿主机环境变量读取
    volumes:
      - ./config:/app/config
```

```bash
# 启动前设置敏感信息
export JWT_SECRET="production-jwt-secret"
export AI_API_KEY="your-ai-api-key"
docker-compose up
```

### 场景4: 生产环境

```bash
# 1. 准备生产配置文件（不提交到Git）
cp config/config.yaml config/config.prod.yaml
# 编辑 config.prod.yaml，设置生产环境的基础配置

# 2. 通过环境变量覆盖敏感信息
export QINGYU_JWT_SECRET="$(openssl rand -base64 32)"
export QINGYU_DATABASE_PRIMARY_MONGODB_URI="mongodb://prod-user:$(cat /secrets/db-password)@prod-host:27017"
export QINGYU_AI_API_KEY="$(cat /secrets/ai-key)"

# 3. 启动服务
export CONFIG_FILE="./config/config.prod.yaml"
./qingyu_backend
```

---

## 🔍 配置验证

### 启动时验证

项目会在启动时自动验证配置：

```go
// config/validation.go
func ValidateConfig(cfg *Config) error {
    // 验证必填项
    if cfg.Server.Port == "" {
        return errors.New("server.port is required")
    }

    // 验证JWT密钥强度
    if len(cfg.JWT.Secret) < 16 {
        return errors.New("jwt.secret must be at least 16 characters")
    }

    // 验证数据库配置
    if cfg.Database == nil {
        return errors.New("database configuration is required")
    }

    return nil
}
```

### 查看当前配置

使用配置管理API查看生效的配置：

```bash
# 获取管理员Token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"your_password"}' | jq -r '.data.token')

# 查看所有配置
curl -X GET http://localhost:8080/api/v1/admin/config \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📝 最佳实践

### 1. 配置文件管理

```bash
# ✅ 提交到Git
config/config.yaml          # 基础配置，包含默认值
config/config.docker.yaml   # Docker环境配置
config/config.test.yaml     # 测试环境配置

# ❌ 不要提交到Git
config/config.local.yaml    # 本地开发个人配置
config/config.prod.yaml     # 生产环境配置
.env                        # 环境变量文件（如果使用）
.env.local                  # 本地环境变量
```

### 2. 敏感信息管理

```bash
# ❌ 不要在配置文件中硬编码
jwt:
  secret: "my-secret-key"  # 不安全

# ✅ 使用环境变量
jwt:
  secret: "default-dev-key"  # 仅用于开发

# 生产环境通过环境变量覆盖
export QINGYU_JWT_SECRET="$(openssl rand -base64 32)"
```

### 3. 配置分层

```yaml
# config.yaml - 基础配置
server:
  port: "8080"
  mode: "debug"

# 通过环境变量覆盖特定值
# export QINGYU_SERVER_MODE="release"
```

### 4. 文档化配置

在配置文件中添加注释：

```yaml
# config.yaml
server:
  port: "8080"  # 服务器监听端口
  mode: "debug" # 运行模式: debug/release

jwt:
  secret: "qingyu_secret_key"  # JWT密钥，生产环境必须修改
  expiration_hours: 24          # Token有效期（小时）
```

---

## 🆚 对比：.env vs Viper环境变量

| 特性 | .env文件 | Viper环境变量 |
|-----|---------|--------------|
| **加载方式** | 需要 `godotenv` 包 | 原生支持 |
| **优先级** | 需要手动实现 | 自动处理 |
| **类型支持** | 仅字符串 | 自动类型转换 |
| **嵌套配置** | 不支持 | 支持 |
| **热重载** | 不支持 | 支持 |
| **默认值** | 需要手动处理 | 自动处理 |
| **验证** | 需要手动实现 | 集成验证 |
| **Docker集成** | 需要额外处理 | 原生支持 |

---

## 📚 相关文档

- `doc/配置文件加载顺序和优先级说明.md` - 配置加载详细说明
- `doc/api/admin/配置管理API文档.md` - 配置管理API
- `doc/usage/配置管理工具使用指南.md` - 配置管理工具
- `config/config.go` - Viper配置代码
- `config/viper_integration.go` - Viper集成实现

---

## 💡 总结

1. **不需要 `.env` 文件**：项目使用 Viper + YAML + 环境变量
2. **环境变量覆盖**：通过 `QINGYU_` 前缀的环境变量覆盖配置
3. **多环境支持**：通过 `CONFIG_FILE` 环境变量指定配置文件
4. **配置优先级**：环境变量 > 配置文件 > 默认值
5. **敏感信息**：通过环境变量传递，不要硬编码在配置文件中

**推荐做法**：
- ✅ 基础配置放在 YAML 文件中
- ✅ 敏感信息通过环境变量传递
- ✅ 不同环境使用不同配置文件
- ✅ 使用配置管理API实时调整配置

---

**维护者**：青羽后端团队
**最后更新**：2025-10-25
**版本**：v1.0

