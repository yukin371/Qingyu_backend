# 青羽写作平台 - 测试数据更新工具

## 功能介绍

这是一个用于快速创建和更新测试数据的工具，包含以下数据：

### 1. 用户数据 (users)
- **管理员**: 2个
  - admin / Admin@123456 (超级管理员)
  - admin02 / Admin@123456 (内容管理员)

- **作者**: 5个
  - author_famous / Author@123456 (知名作家)
  - author_new / Author@123456 (新人作家)
  - author_veteran / Author@123456 (资深作家)
  - author_pro / Author@123456 (全职作家)
  - author_parttime / Author@123456 (业余作家)

- **VIP读者**: 3个
  - reader_vip01 / Vip@123456
  - reader_vip02 / Vip@123456
  - reader_vip03 / Vip@123456

- **普通读者**: 5个
  - reader_normal01-05 / Reader@123456

- **特殊状态用户**: 2个
  - user_banned (被封禁)
  - user_inactive (未激活)

### 2. 书籍数据 (books)
- 包含5本示例书籍：
  - 修真世界 (仙侠, 已完结)
  - 诡秘之主 (奇幻, 已完结)
  - 全职高手 (游戏, 已完结)
  - 斗破苍穹 (玄幻, 免费)
  - 大奉打更人 (玄幻, 连载中)

### 3. 章节数据 (chapters)
- 每本书自动生成章节
- 完结书籍: 100-600章
- 连载书籍: 50-250章
- 前10章免费，后续收费
- 前20章包含完整示例内容

### 4. 钱包数据 (wallets)
- 所有活跃用户自动创建钱包
- 根据角色设置初始余额
  - 管理员: 10,000
  - 作者: 5,000-10,000
  - 读者: 100-1,000
- 自动生成交易记录

### 5. 社交数据
- **评论**: 每本书5-20条评论 + 随机回复
- **点赞**: 每本书20-80个点赞
- **收藏**: 每个用户5-30个收藏
- **关注**: 每个用户关注5-20个其他用户

## 使用方法

### 方式一：使用批处理脚本（推荐）

#### Windows
```bash
cd D:\Github\青羽\Qingyu_backend
scripts\update_test_data.bat
```

#### Linux/Mac
```bash
cd Qingyu_backend
chmod +x scripts/update_test_data.sh
./scripts/update_test_data.sh
```

### 方式二：直接运行

```bash
cd Qingyu_backend
go run cmd/seed_data/main.go
```

### 方式三：编译后运行

```bash
cd Qingyu_backend
go build -o seed_data.exe ./cmd/seed_data
./seed_data.exe  # Linux/Mac: ./seed_data
```

## 操作菜单

运行工具后，会显示以下菜单：

```
请选择操作：
1. 全部更新（清理旧数据 + 创建新数据）
2. 仅创建新数据（跳过已存在的）
3. 清理所有测试数据
4. 查看数据统计
5. 退出
```

### 选项说明

1. **全部更新**: 删除所有测试数据，重新创建全新的数据
2. **仅创建新数据**: 保留已存在的数据，只添加缺失的数据
3. **清理所有测试数据**: 删除所有测试数据，不创建新数据
4. **查看数据统计**: 显示当前数据库中各集合的记录数量
5. **退出**: 退出程序

## 注意事项

1. **数据库连接**
   - 确保 MongoDB 已启动
   - 检查 `config/config.yaml` 中的数据库配置

2. **数据安全**
   - 此工具会修改数据库数据
   - 建议在测试环境使用
   - 生产环境慎用

3. **密码管理**
   - 测试账号密码较为简单
   - 生产环境请使用强密码
   - 建议定期更换测试密码

4. **数据量**
   - 完整数据量约数千条记录
   - 根据书籍数量，章节数量会很大

## 文件结构

```
Qingyu_backend/
├── cmd/
│   └── seed_data/
│       └── main.go              # 主程序
├── migration/
│   └── seeds/
│       ├── enhanced_users.go    # 用户种子数据
│       ├── books.go             # 书籍种子数据
│       ├── chapters.go          # 章节种子数据
│       ├── wallets.go           # 钱包种子数据
│       ├── social.go            # 社交种子数据
│       └── README.md            # 本文件
└── scripts/
    ├── update_test_data.bat     # Windows批处理脚本
    └── update_test_data.sh      # Linux/Mac Shell脚本
```

## 故障排除

### 1. 数据库连接失败
```
❌ 初始化数据库失败: context deadline exceeded
```
**解决方案**:
- 检查 MongoDB 是否运行
- 验证配置文件中的连接字符串

### 2. 权限错误
```
❌ 创建用户失败: not authorized
```
**解决方案**:
- 确保数据库用户有创建权限
- 检查 MongoDB 认证配置

### 3. 端口占用
```
❌ 连接失败: connect ECONNREFUSED
```
**解决方案**:
- 检查 MongoDB 端口（默认27017）
- 验证防火墙设置

## 扩展开发

如需添加更多种子数据：

1. 在 `migration/seeds/` 目录创建新文件
2. 实现 `Seed*` 函数
3. 在 `cmd/seed_data/main.go` 中调用

示例：
```go
// migration/seeds/custom.go
package seeds

func SeedCustomData(ctx context.Context, db *mongo.Database) error {
    // 你的代码
    return nil
}
```

```go
// cmd/seed_data/main.go
// 在 createAllData 函数中添加
if err := seeds.SeedCustomData(ctx, global.DB); err != nil {
    fmt.Printf("❌ 创建自定义数据失败: %v\n", err)
}
```

## 版本历史

- **v1.0** (2026-01-01)
  - 初始版本
  - 支持用户、书籍、章节、钱包、社交数据
  - 交互式菜单

## 反馈与支持

如有问题或建议，请联系开发团队。
