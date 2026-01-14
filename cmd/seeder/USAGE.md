# Seeder CLI 使用说明

## 简介

seeder 是青羽写作平台的测试数据填充工具，使用 Cobra 框架构建的命令行界面。

## 编译

```bash
cd Qingyu_backend/cmd/seeder
go build -o seeder.exe .
```

## 命令列表

### 1. all - 执行所有数据填充

填充所有类型的测试数据（用户、书籍、订阅关系等）。

```bash
# 使用默认配置（medium 规模）
./seeder.exe all

# 指定数据规模
./seeder.exe all --scale small

# 填充前清空现有数据
./seeder.exe all --scale large --clean
```

### 2. users - 只填充用户数据

只填充用户相关的测试数据。

```bash
# 填充用户数据
./seeder.exe users

# 使用 small 规模填充
./seeder.exe users -s small

# 填充前清空用户数据
./seeder.exe users -s medium -c
```

### 3. bookstore - 只填充书籍数据

只填充书籍相关的测试数据。

```bash
# 填充书籍数据
./seeder.exe bookstore

# 使用 large 规模填充
./seeder.exe bookstore -s large

# 填充前清空书籍数据
./seeder.exe bookstore -c
```

### 4. clean - 清空所有测试数据

清空数据库中的所有测试数据。

```bash
# 执行清空操作（需要输入 YES 确认）
./seeder.exe clean

# 提示示例
# 警告: 此操作将清空所有测试数据!
# 请输入 'YES' 确认: YES
# 数据清空完成!

# 取消操作
# 请输入 'YES' 确认: NO
# 操作已取消
```

### 5. verify - 验证数据完整性

验证数据库中的测试数据是否完整和正确。

```bash
# 执行验证
./seeder.exe verify

# 输出示例
# 验证数据完整性...
#
# ✅ 用户数据: 通过
#    - 所有用户名唯一
#
# ✅ 书籍数据: 通过
#    - 所有书籍评分在有效范围内 (0-10)
#
# ✅ 订阅关系: 通过
#    - 所有订阅关系有效
#
# 总计: 3/3 验证通过
```

## 全局标志

所有命令都支持以下全局标志：

| 标志 | 简写 | 说明 | 默认值 |
|------|------|------|--------|
| --config | -c | 配置文件路径 | $HOME/.seeder.yaml |
| --scale | -s | 数据规模 (small/medium/large) | medium |
| --clean | -c | 填充前清空现有数据 | false |
| --help | -h | 显示帮助信息 | - |

## 数据规模说明

| 规模 | 用户数 | 书籍数 | 作者数 |
|------|--------|--------|--------|
| small | 50 | 100 | 20 |
| medium | 500 | 500 | 100 |
| large | 2000 | 1200 | 400 |

## 验证功能详情

verify 命令会执行以下验证：

1. **用户数据验证**
   - 检查用户名唯一性
   - 查找并报告重复的用户名

2. **书籍数据验证**
   - 检查评分是否在有效范围内 (0-10)
   - 统计并列出超出范围的书籍

3. **订阅关系验证**
   - 检查孤儿订阅（引用不存在的用户）
   - 使用 $lookup 进行关联查询验证

## 常用场景

### 场景 1: 初始化测试环境

```bash
# 清空现有数据并填充完整数据集
./seeder.exe clean --yes
./seeder.exe all --scale medium
```

### 场景 2: 只测试用户功能

```bash
# 只填充用户数据
./seeder.exe users --scale small
```

### 场景 3: 数据验证

```bash
# 填充后验证数据
./seeder.exe all
./seeder.exe verify
```

### 场景 4: 快速开发测试

```bash
# 使用小规模数据快速测试
./seeder.exe all --scale small --clean
```

## 注意事项

1. **MongoDB 连接**: 确保 MongoDB 服务正在运行，默认连接地址为 `mongodb://localhost:27017`
2. **数据库名称**: 默认使用 `qingyu` 数据库
3. **清空操作**: clean 命令需要手动输入 "YES" 确认，防止误操作
4. **批量插入**: 默认批次大小为 100，可在配置中修改
5. **退出码**: verify 命令在验证失败时会返回非零退出码

## 配置文件

支持通过配置文件自定义设置（可选）：

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

## 帮助信息

查看命令帮助：

```bash
# 查看主帮助
./seeder.exe --help

# 查看子命令帮助
./seeder.exe all --help
./seeder.exe users --help
./seeder.exe verify --help
```
