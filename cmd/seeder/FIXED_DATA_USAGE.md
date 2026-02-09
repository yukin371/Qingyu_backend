# 固定数据填充工具使用指南

## 概述

固定数据填充工具允许从YAML配置文件加载预定义的测试数据，用于E2E测试和演示。

## 数据结构

固定数据包括:
- **10个用户**: 5个读者、2个VIP读者、4个作者、1个管理员
- **10本书**: 包含完整的章节内容
- **默认密码**: 所有账号密码为 `password`

## 命令使用

### 1. 填充固定数据

```bash
# 填充固定数据（保留现有数据）
./seeder fixed

# 填充前清空现有固定数据
./seeder fixed --clean
```

### 2. 显示测试账号

```bash
./seeder show-accounts
```

输出示例:
```
========== 测试账号列表 ==========
默认密码: password

用户名: reader01        邮箱: reader01@qingyu.test    角色: reader
用户名: author01        邮箱: author01@qingyu.test    角色: reader author
用户名: admin01         邮箱: admin01@qingyu.test     角色: admin
===================================
```

## 测试账号说明

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

## 固定书籍列表

1. **修仙世界** (author01) - 仙侠类，3章
2. **都市之王** (author02) - 都市类，2章
3. **星际争霸** (author03) - 科幻类，1章
4. **大唐风华** (author04) - 历史类，1章
5. **修仙归来** (author01) - 仙侠/都市类，1章
6. **科技帝国** (author03) - 科幻/都市类，1章
7. **剑道独尊** (author01) - 仙侠类，1章
8. **商海沉浮** (author02) - 都市类，1章
9. **大明王朝** (author04) - 历史类，1章
10. **万古修仙** (author01) - 仙侠类，1章

## E2E测试使用

在E2E测试中，可以使用以下固定账号:

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

## 故障排除

### 问题: 配置文件未找到

```
错误: 读取配置文件失败: open data/fixed_data.yaml: no such file or directory
```

**解决**: 确保在seeder目录下运行命令，或提供完整路径。

### 问题: 用户已存在

```
用户 reader01 已存在，跳过
```

**说明**: 这是正常行为，工具会跳过已存在的用户。如需重新创建，使用 `--clean` 选项。
