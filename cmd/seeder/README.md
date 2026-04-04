# 测试数据填充工具 (Seeder)

青羽写作平台统一的测试数据生成工具，整合了所有分散的数据填充功能。

## 目录

- [功能特性](#功能特性)
- [安装与编译](#安装与编译)
- [快速开始](#快速开始)
- [命令参考](#命令参考)
- [数据规模](#数据规模)
- [生成的数据详情](#生成的数据详情)
- [测试账号](#测试账号)
- [JS链路增强脚本](#js链路增强脚本)
- [配置](#配置)
- [项目结构](#项目结构)
- [常见问题](#常见问题)
- [迁移指南](#迁移指南)

---

## 功能特性

- 🚀 高性能批量数据生成
- 📊 支持多种数据规模（small/medium/large）
- 🔒 安全的 MongoDB 批量操作
- ✅ 内置数据完整性验证
- 🎙️ 友好的命令行界面（Cobra框架）
- 📚 一站式数据填充 - 整合所有分散的数据填充工具
- 🎯 分层命令 - 精选演示、基线数据、扩展数据、链路增强

---

## 安装与编译

### 环境要求

- Go 1.22+
- MongoDB 4.4+

### 编译

```bash
cd Qingyu_backend/cmd/seeder
go mod download
go build -o seeder.exe .    # Windows
go build -o seeder .         # Linux/Mac
```

### 默认配置

- MongoDB URI: `mongodb://localhost:27017`
- Database: `qingyu`
- Batch Size: 100

---

## 快速开始

### 一键初始化

```bash
# 推荐：构建联调基线（本地开发、联调、CI初始化）
./seeder.exe baseline --scale medium --clean

# 推荐：构建精选演示数据（首页、榜单、详情页演示）
./seeder.exe showcase --clean

# 推荐：构建完整业务测试数据（基线 + 钱包/通知/消息/财务/AI配额）
./seeder.exe full --scale medium --clean

# 兼容旧入口，等价于 baseline
./seeder.exe all --scale medium --clean
```

### 查看帮助

```bash
./seeder.exe --help
./seeder.exe baseline --help
```

---

## 命令参考

### 全局标志

| 标志 | 简写 | 说明 | 默认值 |
|------|------|------|--------|
| --scale | -s | 数据规模 (small/medium/large) | medium |
| --clean | -c | 填充前清空现有数据 | false |
| --help | -h | 显示帮助信息 | - |

### 分层命令概览

| 层级 | 命令 | 适用场景 | 依赖 |
|------|------|---------|------|
| 精选演示层 | `showcase` | 首页、榜单、详情页演示数据 | users, categories |
| 基线数据层 | `baseline` | 本地开发、联调基线、CI 初始化 | - |
| 扩展数据层 | `full` | 钱包/通知/消息/财务等完整业务测试 | - |
| 链路增强层 | JS脚本 | 发布后作者统计、评分、评论、阅读行为补数 | - |

### 核心命令

#### baseline - 构建联调基线数据

填充用户、设置、书城、章节、订阅、社交、阅读和统计数据。

```bash
./seeder.exe baseline --scale small --clean
```

#### full - 构建完整测试数据

在 `baseline` 基础上继续填充钱包、通知、消息、财务和 AI 配额数据。

```bash
./seeder.exe full --scale medium --clean
```

#### showcase - 构建精选演示数据

适合首页、榜单、详情页演示。插入少量手工编写的精选书籍，榜单优先展示这些作品。

```bash
./seeder.exe showcase --clean
```

#### all - 兼容旧入口

等价于 `baseline`。

```bash
./seeder.exe all --scale medium --clean
```

### 模块命令

#### users - 填充用户数据

```bash
./seeder.exe users -s medium
./seeder.exe users -s small -c    # 填充前清空
```

#### categories - 填充分类数据

填充标准8个分类。

```bash
./seeder.exe categories --clean
```

#### bookstore - 填充书籍数据

只填充书籍和Banner数据。

```bash
./seeder.exe bookstore -s large
./seeder.exe bookstore -c
```

#### chapters - 填充章节数据

为现有书籍生成章节数据和内容。**依赖：bookstore**

```bash
./seeder.exe chapters
./seeder.exe chapters -c
```

#### subscriptions - 刷新订阅关系

单独刷新书籍订阅关系。

```bash
./seeder.exe subscriptions
./seeder.exe subscriptions --clean
```

#### social - 填充社交数据

填充评论、点赞、收藏、关注等社交数据。**依赖：users, bookstore**

```bash
./seeder.exe social
```

#### wallets - 填充钱包数据

填充用户钱包和交易记录数据。**依赖：users**

```bash
./seeder.exe wallets
```

#### rankings - 填充榜单数据

填充各种榜单数据（实时榜、日榜、周榜、月榜等）。**依赖：bookstore**

```bash
./seeder.exe rankings
```

#### reader - 填充阅读数据

填充阅读历史、书架、订阅、阅读进度等数据。**依赖：users, bookstore**

```bash
./seeder.exe reader
./seeder.exe reader -c
```

#### notifications - 填充通知数据

填充用户通知消息。**依赖：users**

```bash
./seeder.exe notifications
./seeder.exe notifications -c
```

#### messaging - 填充消息数据

填充私信和公告数据。**依赖：users**

```bash
./seeder.exe messaging
./seeder.exe messaging -c
```

#### stats - 填充统计数据

填充书籍和章节统计数据。**依赖：bookstore, chapters**

```bash
./seeder.exe stats
./seeder.exe stats -c
```

#### finance - 填充财务数据

填充作者收益和会员数据。**依赖：users, bookstore**

```bash
./seeder.exe finance
./seeder.exe finance -c
```

#### ai-quota - 激活AI配额

为所有用户激活AI写作配额。**依赖：users**

```bash
./seeder.exe ai-quota
```

#### import - 导入小说数据

从JSON文件导入大量小说数据。

```bash
./seeder.exe import                           # 默认路径
./seeder.exe import --file data/novels_100.json
```

### 管理命令

#### clean - 清空所有测试数据

```bash
./seeder.exe clean
# 警告: 此操作将清空所有测试数据!
# 请输入 'YES' 确认: YES
# 数据清空完成!
```

#### verify - 验证数据完整性

```bash
./seeder.exe verify
# ✅ 用户数据: 通过
#    - 所有用户名唯一
# ✅ 书籍数据: 通过
#    - 所有书籍评分在有效范围内 (0-10)
# ✅ 订阅关系: 通过
# 总计: 3/3 验证通过
```

验证功能详情：
1. **用户数据验证** - 检查用户名唯一性
2. **书籍数据验证** - 检查评分范围 (0-10)
3. **订阅关系验证** - 检查孤儿订阅

#### test - 填充E2E测试数据

填充E2E测试所需的特定数据。

```bash
./seeder.exe test
```

### 命令依赖关系

| 命令 | 前置依赖 |
|------|---------|
| chapters | bookstore |
| social | users, bookstore |
| wallets | users |
| rankings | bookstore |
| reader | users, bookstore |
| notifications | users |
| messaging | users |
| stats | bookstore, chapters |
| finance | users, bookstore |
| ai-quota | users |

---

## 数据规模

| 规模 | 用户数 | 书籍数 | 作者数 |
|------|--------|--------|--------|
| small | 50 | 100 | 20 |
| medium | 500 | 500 | 100 |
| large | 2000 | 1200 | 400 |

---

## 生成的数据详情

### 用户数据 (users)

- 真实测试账号（admin, author1, reader1, vip_user）
- 普通用户、作者、VIP 用户
- 随机用户名、邮箱、头像

### 分类数据 (categories)

- 标准 8 个顶级分类
- 书籍写入真实 `category_ids`
- 同时保留 `categories` 名称快照便于展示

### 书籍数据 (bookstore)

**数据策略**：
- 少量精选演示书籍：手工编写元数据，适合首页、榜单、详情页演示
- 大量随机填充书籍：用于列表、搜索、联调和压力测试

**分类比例**：
| 分类 | 比例 |
|------|------|
| 仙侠 | 25% |
| 都市 | 20% |
| 科幻 | 15% |
| 历史 | 10% |
| 玄幻 | 10% |
| 武侠 | 8% |
| 游戏 | 7% |
| 奇幻 | 5% |

**热度等级**：
- 高热度：评分 8.5-9.5，200-500 订阅
- 中热度：评分 6.0-8.5，20-200 订阅
- 低热度：评分 4.0-6.0，0-20 订阅

**Banner 数据**: 2个轮播图

### 精选演示数据 (showcase)

- 默认内置 5 本精选作品
- 带 `seed:showcase` 标签，榜单会轻量优先这些作品
- 适合作为：首页推荐、榜单头部、详情页/阅读页联调演示
- 当前只有 `云海问剑录` 完整补了 banner + 专属章节标题 + 前 3 章手工正文

### 章节数据 (chapters)

- 为每本书生成 30-500 章（根据书籍状态）
- 每章 1000-2000 字
- 前10章免费
- 章节内容自动生成

### 社交数据 (social)

| 类型 | 数量 |
|------|------|
| 评论 | 每本书 5-20 条 |
| 点赞 | 每本书 10-100 个 |
| 收藏 | 每个用户 5-30 个 |
| 关注 | 每个用户 5-50 个 |

### 钱包数据 (wallets)

| 用户类型 | 初始余额 |
|---------|---------|
| 管理员 | 10000 |
| VIP用户 | 5000-10000 |
| 作者 | 3000-6000 |
| 普通用户 | 100-1000 |

自动生成充值和消费记录。

### 榜单数据 (rankings)

- 实时榜、日榜、周榜、月榜
- 新人榜、完结榜
- 根据书籍评分和浏览量排序

### 阅读数据 (reader)

| 数据类型 | 数量 |
|---------|------|
| 阅读历史 | 每个用户 5-30 条 |
| 书架 | 每个用户 5-20 本 |
| 订阅 | 每个用户 3-15 个作者 |
| 书签 | 30% 用户有，每人 1-10 个 |
| 批注 | 20% 用户有，每人 1-20 条 |

### 通知数据 (notifications)

| 类型 | 占比 |
|------|------|
| 评论通知 | 40% |
| 点赞通知 | 30% |
| 关注通知 | 20% |
| 系统通知 | 10% |

- 每个用户 20-50 条通知
- 70% 已读，30% 未读

### 消息数据 (messaging)

**私信**：
- 40% 用户有私信记录
- 每个用户 3-10 个对话
- 每个对话 2-20 条消息
- 60% 已读，40% 未读

**公告**：
- 3-5 条系统公告
- 2-3 条活动公告
- 1-2 条更新公告

**消息类型**：文本(70%)、图片(15%)、系统(10%)、其他(5%)

### 统计数据 (stats)

**书籍统计**：
- 高热度书籍：5000-20000 浏览
- 中热度书籍：1000-5000 浏览
- 低热度书籍：100-1000 浏览

**章节统计**：
- 平均停留时间：30-180 秒
- 跳章率：10%-40%
- 完读率：40%-90%

### 财务数据 (finance)

**作者收益**：
| 类型 | 占比 |
|------|------|
| 订阅收益 | 60% |
| 打赏收益 | 25% |
| 广告收益 | 10% |
| 其他收益 | 5% |

**会员**：
- 20% 用户为VIP会员
- 月度会员：60%
- 年度会员：30%
- 终身会员：10%

### AI配额 (ai-quota)

| 用户类型 | 配额 |
|---------|------|
| 管理员 | 999999（无限） |
| VIP用户 | 100000 |
| 普通用户 | 10000 |

---

## 测试账号

### 固定测试账号（E2E测试用）

通过 `./seeder.exe test` 或固定数据填充生成。

**默认密码**: `password`

### 读者账号 (5个)

| 用户名 | 邮箱 | 角色 | 说明 |
|--------|------|------|------|
| reader01 | reader01@qingyu.test | reader | 普通读者 |
| reader02 | reader02@qingyu.test | reader | 普通读者 |
| reader03 | reader03@qingyu.test | reader | 普通读者 |
| vipreader01 | vipreader01@qingyu.test | reader, vip | VIP读者 |
| vipreader02 | vipreader02@qingyu.test | reader, vip | VIP读者 |

### 作者账号 (4个)

| 用户名 | 邮箱 | 角色 | 说明 |
|--------|------|------|------|
| author01 | author01@qingyu.test | reader, author | 修仙作者 |
| author02 | author02@qingyu.test | reader, author | 都市作者 |
| author03 | author03@qingyu.test | reader, author | 科幻作者 |
| author04 | author04@qingyu.test | reader, author | 历史作者 |

### 管理员账号 (1个)

| 用户名 | 邮箱 | 角色 | 说明 |
|--------|------|------|------|
| admin01 | admin01@qingyu.test | admin | 系统管理员 |

### E2E测试示例

```typescript
// 登录测试
await page.fill('#username', 'reader01')
await page.fill('#password', 'password')
await page.click('button[type="submit"]')

// 作者测试
await page.fill('#username', 'author01')
await page.fill('#password', 'password')
await page.click('button[type="submit"]')

// 管理员测试
await page.fill('#username', 'admin01')
await page.fill('#password', 'password')
await page.click('button[type="submit"]')
```

---

## JS链路增强脚本

适用于"基线数据已有，但需要让某个作者项目的统计页马上有数据"的场景。

### 脚本位置

`Qingyu_backend/scripts/bootstrap_test_data.mjs`

### 使用方法

```bash
# 先准备基线，再对指定作者项目补评分/评论/书签/阅读历史/阅读行为，并回刷统计
node scripts/bootstrap_test_data.mjs --mode bootstrap --scale small --projectId <writer_project_id>

# 只做作者统计链路增强，不重跑基线
node scripts/bootstrap_test_data.mjs --mode author-stats --projectId <writer_project_id>

# 按作者自动扫描多个项目并批量增强
node scripts/bootstrap_test_data.mjs --mode author-stats-all --authorUsername testauthor001 --limit 5

# 自动补齐"作者发布 + 管理员审核 + 读者互动"
# 未发布但已有文档的项目会被自动发布并审核通过
node scripts/bootstrap_test_data.mjs --mode author-stats-all --authorUsername testauthor001 --adminUsername testadmin001 --limit 5

# 只跑基线
node scripts/bootstrap_test_data.mjs --mode baseline --scale small --clean
```

### 脚本职责

1. 登录作者和测试读者
2. 解析 `project_id -> published book_id`
3. 必要时自动补齐"作者发布 + 管理员审核"
4. 通过真实 API 写入评分、评论、书签、阅读历史、阅读行为
5. 最后执行 `seeder stats` 回刷作者统计页所需聚合

---

## 配置

### 环境变量（待实现）

可通过环境变量修改默认配置。

### 配置文件（可选）

```yaml
mongodb:
  uri: "mongodb://localhost:27017"
  database: "qingyu"

scale: "medium"
batch_size: 100
```

使用配置文件：

```bash
./seeder.exe all --config /path/to/config.yaml
```

---

## 项目结构

```
cmd/seeder/
├── main.go                  # 主程序和 CLI
├── README.md                # 本文档
├── config/                  # 配置管理
├── generators/              # 数据生成器
│   ├── base.go             # 基础生成器
│   ├── user.go             # 用户生成器
│   ├── reader.go           # 阅读数据生成器
│   └── book.go             # 书籍生成器
├── utils/                   # 工具函数
│   ├── mongodb.go          # MongoDB 操作
│   └── verify.go           # 数据验证
├── models/                  # 数据模型
│   ├── user.go
│   ├── book.go
│   ├── reader.go           # 阅读相关模型
│   ├── messaging.go        # 消息相关模型
│   ├── notification.go     # 通知相关模型
│   ├── writer.go           # 作者相关模型
│   ├── finance.go          # 财务相关模型
│   └── stats.go            # 统计相关模型
├── relationships/           # 关系构建
│   └── builder.go
├── validator/              # 数据验证
│   └── validator.go
├── seeder_*.go             # 数据填充器
│   ├── seeder_users.go
│   ├── seeder_bookstore.go
│   ├── seeder_chapters.go
│   ├── seeder_social.go
│   ├── seeder_wallets.go
│   ├── seeder_rankings.go
│   ├── seeder_reader.go
│   ├── seeder_notification.go
│   ├── seeder_messaging.go
│   ├── seeder_stats.go
│   ├── seeder_finance.go
│   ├── seeder_ai_quota.go
│   ├── seeder_import.go
│   ├── seeder_test_data.go
│   ├── seeder_categories.go
│   ├── seeder_subscriptions.go
│   ├── seeder_publication_flow.go
│   └── seeder_settings.go
├── showcase_content.go     # 精选演示内容
├── showcase_books.go       # 精选演示书籍
├── audit_reader.go         # 读者数据审计
└── audit_author.go         # 作者数据审计
```

---

## 常见问题

### MongoDB 连接失败

确保 MongoDB 服务正在运行：

```bash
# Windows
net start MongoDB

# Linux
sudo systemctl start mongod

# Mac
brew services start mongodb
```

### 配置文件未找到

```
错误: 读取配置文件失败: open data/fixed_data.yaml: no such file or directory
```

**解决**: 确保在seeder目录下运行命令，或提供完整路径。

### 用户已存在

```
用户 reader01 已存在，跳过
```

**说明**: 这是正常行为，工具会跳过已存在的用户。如需重新创建，使用 `--clean` 选项。

### 命令依赖错误

某些命令需要先运行其他命令，请参考[命令依赖关系](#命令依赖关系)。

---

## 常用场景

### 场景 1: 初始化完整测试环境

```bash
./seeder.exe clean
echo "YES" | ./seeder.exe clean
./seeder.exe full --scale medium
./seeder.exe verify
```

### 场景 2: 只测试用户功能

```bash
./seeder.exe users --scale small
```

### 场景 3: 测试阅读流程

```bash
./seeder.exe bookstore
./seeder.exe chapters
```

### 场景 4: 测试社交功能

```bash
./seeder.exe baseline
./seeder.exe social
```

### 场景 5: 快速开发测试

```bash
./seeder.exe baseline --scale small --clean
```

### 场景 6: 首页/榜单演示

```bash
./seeder.exe showcase --clean
```

---

## 迁移指南

如果你之前使用分散的命令，可以迁移到新的 seeder 工具：

| 旧命令 | 新命令 |
|--------|--------|
| `go run cmd/seed_data/main.go` | `./seeder baseline && ./seeder chapters && ./seeder social && ./seeder wallets` |
| `go run cmd/seed_bookstore/main.go` | `./seeder bookstore` |
| `go run cmd/import_novels_auto/main.go` | `./seeder import` |
| `go run cmd/init_test_data/main.go` | `./seeder baseline && ./seeder ai-quota` |

---

## 技术栈

- Go 1.22+
- MongoDB 4.4+
- gofakeit v7 - 数据生成
- Cobra v1.8 - CLI 框架

---

## 许可证

MIT License
