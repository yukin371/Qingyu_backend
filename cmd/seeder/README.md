# 测试数据填充工具 (Seeder)

为青羽写作平台快速生成大规模测试数据。

## 功能特性

- 🚀 高性能批量数据生成
- 📊 支持多种数据规模（small/medium/large）
- 🔒 安全的 MongoDB 批量操作
- ✅ 内置数据完整性验证
- 🎙️ 友好的命令行界面

## 安装

```bash
cd Qingyu_backend/cmd/seeder
go mod download
go build -o seeder .
```

## 使用方法

### 查看帮助

```bash
./seeder --help
./seeder all --help
```

### 填充所有数据

```bash
# 使用默认规模 (medium)
./seeder all

# 指定规模
./seeder all --scale small   # 50用户, 100本书
./seeder all --scale medium  # 500用户, 500本书
./seeder all --scale large   # 2000用户, 1200本书

# 填充前清空现有数据
./seeder all --clean
```

### 填充特定模块

```bash
# 只填充用户
./seeder users -s medium

# 只填充书籍
./seeder bookstore -s large
```

### 清空数据

```bash
# 清空所有测试数据（需要确认）
./seeder clean
```

### 验证数据

```bash
./seeder verify
```

## 数据规模

| 规模   | 用户数 | 书籍数 | 作者数 |
| ------ | ------ | ------ | ------ |
| small  | 50     | 100    | 20     |
| medium | 500    | 500    | 100    |
| large  | 2000   | 1200   | 400    |

## 生成的数据

### 用户数据
- 真实测试账号（admin, author1, reader1, vip_user）
- 普通用户、作者、VIP 用户
- 随机用户名、邮箱、头像

### 书籍数据
- **分类比例**：
  - 仙侠: 30%
  - 都市: 25%
  - 科幻: 20%
  - 历史: 15%
  - 其他: 10%
- **热度等级**：
  - 高热度：评分 8.5-9.5，200-500 订阅
  - 中热度：评分 6.0-8.5，20-200 订阅
  - 低热度：评分 4.0-6.0，0-20 订阅

### 订阅关系
- 根据书籍评分智能分配订阅数
- 订阅时间随机分布在书籍发布时间之后
- 使用 Fisher-Yates 洗牌算法公平选择订阅用户

## 配置

默认配置：
- MongoDB URI: `mongodb://localhost:27017`
- Database: `qingyu`
- Batch Size: 100

可通过环境变量或配置文件修改（待实现）。

## 验证功能

验证工具检查：
- ✅ 用户名唯一性
- ✅ 书籍评分范围 (0-10)
- ✅ 订阅关系有效性（无孤儿记录）

## 真实数据文件

修改 `data/` 目录下的 JSON 文件来定制测试数据：

- `users.json` - 测试账号（admin, author1 等）
- `books.json` - 推荐书籍模板
- `banners.json` - 轮播图配置

## 项目结构

```
cmd/seeder/
├── main.go              # 主程序和 CLI
├── config/              # 配置管理
├── generators/          # 数据生成器
│   ├── base.go         # 基础生成器
│   ├── user.go         # 用户生成器
│   └── book.go         # 书籍生成器
├── relationships/       # 关联关系构建
│   └── builder.go      # 订阅关系构建器
├── utils/              # 工具函数
│   ├── mongodb.go      # MongoDB 操作
│   └── verify.go       # 数据验证
├── models/             # 数据模型
│   ├── user.go
│   └── book.go
├── data/               # 真实数据模板
└── seeder_*.go         # 数据填充器
```

## 技术栈

- Go 1.22+
- MongoDB 4.4+
- gofakeit v7 - 数据生成
- Cobra v1.8 - CLI 框架

## 开发

### 运行测试

```bash
cd Qingyu_backend/cmd/seeder
go test ./... -v
```

### 添加新的数据类型

1. 在 `models/` 创建模型定义
2. 在 `generators/` 创建生成器
3. 在 `seeder_*.go` 创建填充器
4. 在 `main.go` 添加 CLI 命令

## 常见问题

### MongoDB 连接失败

确保 MongoDB 服务正在运行：
```bash
# Windows
net start MongoDB

# Linux/Mac
brew services start mongodb  # Mac
sudo systemctl start mongod  # Linux
```

### 验证失败

运行验证命令查看详细错误：
```bash
./seeder verify
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
