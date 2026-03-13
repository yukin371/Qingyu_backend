# 测试数据填充工具 (Seeder)

为青羽写作平台快速生成大规模测试数据的统一工具。

## 功能特性

- 🚀 高性能批量数据生成
- 📊 支持多种数据规模（small/medium/large）
- 🔒 安全的 MongoDB 批量操作
- ✅ 内置数据完整性验证
- 🎙️ 友好的命令行界面
- 📚 **一站式数据填充** - 整合了所有分散的数据填充工具

## 安装

```bash
cd Qingyu_backend/cmd/seeder
go mod download
go build -o seeder .
```

## 命令列表

### 分层命令

| 层级 | 适用场景 | 推荐入口 |
|------|---------|---------|
| 精选演示层 | 首页、榜单、详情页演示数据 | `showcase` |
| 基线数据层 | 本地开发、联调基线、CI 初始化 | `baseline` |
| 扩展数据层 | 钱包/通知/消息/财务等完整业务测试 | `full` |
| 链路增强层 | 发布后作者统计、评分、评论、阅读行为补数 | `node scripts/bootstrap_test_data.mjs --mode author-stats` |

### 核心命令

| 命令 | 说明 | 依赖 |
|------|------|------|
| `all` | 兼容旧入口，等价于 `baseline` | - |
| `showcase` | 只填充少量精选演示书籍，并让榜单优先可见 | users, categories |
| `baseline` | 构建联调基线数据（用户、书城、章节、订阅、社交、阅读、统计） | - |
| `full` | 构建完整测试数据（基线 + 钱包/通知/消息/财务/AI配额） | - |
| `users` | 只填充用户数据 | - |
| `categories` | 填充标准分类数据（8个分类） | - |
| `bookstore` | 只填充书籍数据和Banner | - |
| `subscriptions` | 只填充书籍订阅关系 | users, books |
| `chapters` | 填充章节数据和内容 | books |
| `social` | 填充社交数据（评论、点赞、收藏、关注） | users, books |
| `wallets` | 填充钱包和交易数据 | users |
| `rankings` | 填充榜单数据 | books |
| `reader` | 填充阅读数据（阅读历史、书架、订阅） | users, books |
| `notifications` | 填充通知数据 | users |
| `messaging` | 填充消息数据（私信、对话、公告） | users |
| `stats` | 填充统计数据（书籍、章节统计） | books, chapters |
| `finance` | 填充财务数据（作者收益、会员） | users, books |
| `ai-quota` | 激活用户AI配额 | users |
| `import` | 从JSON文件导入小说数据 | - |
| `clean` | 清空所有测试数据 | - |
| `verify` | 验证数据完整性 | - |
| `test` | 填充E2E测试所需的特定数据 | - |

### 全局标志

所有命令都支持以下全局标志：

| 标志 | 简写 | 说明 | 默认值 |
|------|------|------|--------|
| --scale | -s | 数据规模 (small/medium/large) | medium |
| --clean | -c | 填充前清空现有数据 | false |

## 使用方法

### 查看帮助

```bash
./seeder --help
./seeder all --help
```

### 一键初始化完整测试环境

```bash
# 推荐：构建联调基线
./seeder baseline --scale medium --clean

# 推荐：先构建精选演示数据，用于首页和榜单演示
./seeder showcase --clean

# 推荐：构建完整业务测试数据
./seeder full --scale medium --clean
```

### 填充特定模块

```bash
# 只填充用户
./seeder users -s medium

# 只填充分类
./seeder categories --clean

# 只填充书籍
./seeder bookstore -s large

# 只填充精选演示书籍
./seeder showcase --clean

# 只刷新订阅关系
./seeder subscriptions

# 填充章节数据（需要先有书籍）
./seeder chapters

# 填充社交数据（需要先有用户和书籍）
./seeder social

# 填充钱包数据（需要先有用户）
./seeder wallets

# 填充榜单数据（需要先有书籍）
./seeder rankings

# 填充阅读数据（需要先有用户和书籍）
./seeder reader

# 填充通知数据（需要先有用户）
./seeder notifications

# 填充消息数据（需要先有用户）
./seeder messaging

# 填充统计数据（需要先有书籍和章节）
./seeder stats

# 填充财务数据（需要先有用户和书籍）
./seeder finance

# 激活AI配额（需要先有用户）
./seeder ai-quota
```

### 导入小说数据

```bash
# 从JSON文件导入100本小说
./seeder import

# 自定义文件路径
./seeder import --file data/novels_100.json
```

### 清空数据

```bash
# 清空所有测试数据（需要输入 YES 确认）
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

### 用户数据 (users)
- 真实测试账号（admin, author1, reader1, vip_user）
- 普通用户、作者、VIP 用户
- 随机用户名、邮箱、头像

### 书籍数据 (bookstore)
- **数据策略**：
  - 少量精选演示书籍：手工编写元数据，适合首页、榜单、详情页演示
  - 大量随机填充书籍：用于列表、搜索、联调和压力测试
- **分类比例**：
  - 仙侠: 25%
  - 都市: 20%
  - 科幻: 15%
  - 历史: 10%
  - 玄幻: 10%
  - 武侠: 8%
  - 游戏: 7%
  - 奇幻: 5%
- **热度等级**：
  - 高热度：评分 8.5-9.5，200-500 订阅
  - 中热度：评分 6.0-8.5，20-200 订阅
  - 低热度：评分 4.0-6.0，0-20 订阅
- **Banner 数据**: 2个轮播图

### 分类数据 (categories)
- 标准 8 个顶级分类
- 书籍写入真实 `category_ids`
- 同时保留 `categories` 名称快照便于展示

### 精选演示数据 (showcase)
- 默认内置 5 本精选作品
- 带 `seed:showcase` 标签，榜单会轻量优先这些作品
- 适合作为：
  - 首页推荐
  - 榜单头部
  - 详情页、阅读页联调演示
- 当前只有 `云海问剑录` 完整补了 banner + 专属章节标题 + 前 3 章手工正文
- `TODO(showcase-next)` 已写在 [`showcase_content.go`](/E:/Github/Qingyu/Qingyu_backend/cmd/seeder/showcase_content.go)，后续补其他书直接按模板追加

### 章节数据 (chapters)
- 为每本书生成 30-500 章（根据书籍状态）
- 每章 1000-2000 字
- 前10章免费
- 章节内容自动生成

### 社交数据 (social)
- **评论**: 每本书 5-20 条评论
- **点赞**: 每本书 10-100 个点赞
- **收藏**: 每个用户 5-30 个收藏
- **关注**: 每个用户 5-50 个关注

### 钱包数据 (wallets)
- 管理员：10000 初始余额
- VIP用户：5000-10000 初始余额
- 作者：3000-6000 初始余额
- 普通用户：100-1000 初始余额
- 自动生成充值和消费记录

### 榜单数据 (rankings)
- 实时榜、日榜、周榜、月榜
- 新人榜、完结榜
- 根据书籍评分和浏览量排序

### AI配额 (ai-quota)
- 管理员：999999（无限配额）
- VIP用户：100000
- 普通用户：10000

### 阅读数据 (reader)
- **阅读历史**: 每个用户 5-30 条阅读记录
- **书架**: 每个用户 5-20 本收藏书籍
- **订阅**: 每个用户 3-15 个订阅作者
- **阅读进度**: 自动生成章节阅读进度

### 通知数据 (notifications)
- **评论通知**: 40% - 评论作品相关通知
- **点赞通知**: 30% - 点赞作品相关通知
- **关注通知**: 20% - 新粉丝通知
- **系统通知**: 10% - 平台公告和系统消息
- 每个用户 20-50 条通知
- 70% 已读，30% 未读

### 消息数据 (messaging)
- **私信**: 用户间私信消息
  - 40% 用户有私信记录
  - 每个用户 3-10 个对话
  - 每个对话 2-20 条消息
  - 60% 已读，40% 未读
- **公告**: 平台公告
  - 3-5 条系统公告
  - 2-3 条活动公告
  - 1-2 条更新公告
- 消息类型：文本(70%)、图片(15%)、系统(10%)、其他(5%)

### 统计数据 (stats)
- **书籍统计**: 每本书的浏览、阅读、收藏、分享数据
  - 使用运行时真实模型填充 `book_stats` / `book_stats_daily`
- **章节统计**: 每章的详细统计，填充 `chapter_stats`
- **读者行为**: 自动补齐 `reader_behaviors`
- **留存统计**: 自动补齐 `reader_retentions`

## JS 链路增强脚本

适用于“基线数据已有，但需要让某个作者项目的统计页马上有数据”的场景。

```bash
# 先准备基线，再对指定作者项目补评分/评论/书签/阅读历史/阅读行为，并回刷统计
node scripts/bootstrap_test_data.mjs --mode bootstrap --scale small --projectId <writer_project_id>

# 只做作者统计链路增强，不重跑基线
node scripts/bootstrap_test_data.mjs --mode author-stats --projectId <writer_project_id>

# 按作者自动扫描多个项目并批量增强
node scripts/bootstrap_test_data.mjs --mode author-stats-all --authorUsername testauthor001 --limit 5

# 自动补齐“作者发布 + 管理员审核 + 读者互动”
# 未发布但已有文档的项目会被自动发布并审核通过
node scripts/bootstrap_test_data.mjs --mode author-stats-all --authorUsername testauthor001 --adminUsername testadmin001 --limit 5

# 只跑基线
node scripts/bootstrap_test_data.mjs --mode baseline --scale small --clean
```

脚本职责：
- 登录作者和测试读者
- 解析 `project_id -> published book_id`
- 必要时自动补齐“作者发布 + 管理员审核”
- 通过真实 API 写入评分、评论、书签、阅读历史、阅读行为
- 最后执行 `seeder stats` 回刷作者统计页所需聚合

### 财务数据 (finance)
- **作者收益**: 每本书的收益记录
  - 订阅收益：60%
  - 打赏收益：25%
  - 广告收益：10%
  - 其他收益：5%
  - 月度收益记录
- **会员**: 用户会员信息
  - 20% 用户为VIP会员
  - 月度会员：60%
  - 年度会员：30%
  - 终身会员：10%

## 配置

默认配置：
- MongoDB URI: `mongodb://localhost:27017`
- Database: `qingyu`
- Batch Size: 100

可通过环境变量或配置文件修改（待实现）。

## 项目结构

```
cmd/seeder/
├── main.go                  # 主程序和 CLI
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
├── seeder_*.go              # 数据填充器
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
│   └── seeder_import.go
└── README.md
```

## 技术栈

- Go 1.22+
- MongoDB 4.4+
- gofakeit v7 - 数据生成
- Cobra v1.8 - CLI 框架

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

### 命令依赖关系

某些命令需要先运行其他命令：
- `chapters` 需要先运行 `bookstore`
- `bookstore` 会自动确保 `categories` 已填充
- `social` 需要先运行 `users` 和 `bookstore`
- `wallets` 需要先运行 `users`
- `rankings` 需要先运行 `bookstore`
- `reader` 需要先运行 `users` 和 `bookstore`
- `notifications` 需要先运行 `users`
- `messaging` 需要先运行 `users`
- `stats` 需要先运行 `bookstore` 和 `chapters`
- `finance` 需要先运行 `users` 和 `bookstore`
- `ai-quota` 需要先运行 `users`

## 完整初始化流程

```bash
# 1. 编译工具
cd Qingyu_backend/cmd/seeder
go build -o seeder .

# 2. 填充所有基础数据
./seeder all --scale medium --clean

# 3. 填充扩展数据
./seeder chapters
./seeder social
./seeder wallets
./seeder rankings
./seeder reader
./seeder notifications
./seeder messaging
./seeder stats
./seeder finance
./seeder ai-quota

# 4. 验证数据
./seeder verify
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
