# P0问题错误预防与最佳实践

**创建日期**: 2026-03-06
**创建者**: Kore
**版本**: 1.0
**来源**: P0问题审查总结

---

## 概述

本文档总结了P0问题审查中发现的系统性问题，提供预防和应对措施，避免类似问题再次发生。

### P0问题根源分析

| 问题类型 | 根本原因 | 影响范围 |
|----------|----------|----------|
| 类型定义不一致 | 缺少统一标准 | 数据一致性 |
| 架构分离无关联 | 设计不完整 | 功能缺失 |
| 层级职责混乱 | 缺少架构规范 | 可维护性 |
| 版本管理缺失 | 需求未考虑 | 数据追溯 |

---

## 一、类型定义一致性

### 1.1 问题描述

P0审查中发现**多处重复定义且不一致**的类型：

| 类型 | 定义位置 | 冲突 |
|------|----------|------|
| BookStatus | 3处定义 | `published` vs `ongoing` |
| CategoryIDs | 2种类型 | `[]ObjectID` vs `[]string` |
| ID字段 | 多种类型 | `string` vs `ObjectID` |

### 1.2 预防措施

#### 规则1: 单一数据源原则 (Single Source of Truth)

**所有共享类型必须在`internal/domain/`层定义唯一版本**

```
┌─────────────────────────────────────────────────────────┐
│                    类型定义层次                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │  internal/domain/ (唯一定义源)                   │   │
│  │  ├── enums.go       (枚举类型)                   │   │
│  │  ├── entities.go    (领域实体)                   │   │
│  │  └── value_objects.go (值对象)                   │   │
│  └─────────────────────────────────────────────────┘   │
│                      │ 引用                             │
│                      ▼                                  │
│  ┌─────────────────────────────────────────────────┐   │
│  │  models/ (仅数据模型，不定义业务类型)            │   │
│  │  ├── bookstore/                                 │   │
│  │  └── writer/                                    │   │
│  └─────────────────────────────────────────────────┘   │
│                      │ 使用                             │
│                      ▼                                  │
│  ┌─────────────────────────────────────────────────┐   │
│  │  api/ (DTO层，类型转换)                          │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

#### 规则2: 类型定义审查清单

**添加新类型时必须检查**：

- [ ] 检查`internal/domain/`是否已有类似定义
- [ ] 如有，使用现有定义；如无，在domain层创建
- [ ] 更新`docs/standards/type-definitions.md`（需创建）
- [ ] 运行类型一致性检查脚本

#### 规则3: 禁止重复定义

```go
// ❌ 错误：在models层定义枚举
package models

type BookStatus string
const (
    BookStatusOngoing BookStatus = "ongoing"
)

// ✅ 正确：引用domain层定义
package models

import "qingyu/internal/domain"

type Book struct {
    Status domain.BookStatus `bson:"status" json:"status"`
}
```

### 1.3 检测工具

#### 类型一致性检查脚本

```bash
# scripts/check-type-consistency.sh
#!/bin/bash

# 检查重复的枚举定义
echo "检查重复枚举定义..."
grep -r "type.*Status.*string" --include="*.go" internal/ models/ | cut -d: -f1 | sort | uniq -d

# 检查重复的常量定义
echo "检查重复常量定义..."
grep -r "const.*Status" --include="*.go" internal/ models/

# 检查ID类型一致性
echo "检查ID类型定义..."
grep -r "ID.*string" --include="*.go" models/
```

#### Git Pre-commit Hook

```bash
# .git/hooks/pre-commit
#!/bin/bash

# 运行类型一致性检查
bash scripts/check-type-consistency.sh
if [ $? -ne 0 ]; then
    echo "检测到类型定义不一致！请检查后再提交。"
    exit 1
fi
```

### 1.4 应对流程

**发现类型不一致时**：

1. **立即**：停止使用冲突的类型
2. **分析**：确定正确的定义（通常domain层）
3. **统一**：将所有引用改为使用统一定义
4. **测试**：确保没有破坏现有功能
5. **文档**：更新类型定义文档

---

## 二、架构设计完整性

### 2.1 问题描述

Project-Book分离架构**缺少关联关系**，导致：
- 发布时无法知道Book对应哪个Project
- Project更新时无法同步到Book
- 无法实现"从Project创建Book"功能

### 2.2 预防措施

#### 规则4: 架构设计必须考虑关联关系

**设计分离架构时的检查清单**：

- [ ] 两个实体之间是否需要关联？
- [ ] 如何建立关联？（外键、中间表、事件）
- [ ] 数据如何同步？
- [ ] 同步失败如何处理？
- [ ] 是否需要版本管理？

#### 规则5: 关联字段设计规范

```go
// 关联字段命名规范
type Entity struct {
    // 直接关联（1:1或N:1）
    RelatedEntityID *string `bson:"related_entity_id,omitempty"`

    // 关联元数据
    SourceType    SourceType `bson:"source_type"`
    SyncMode      SyncMode   `bson:"sync_mode"`
    LastSyncedAt  *time.Time `bson:"last_synced_at"`

    // 内容哈希（用于变更检测）
    ContentHash   string     `bson:"content_hash"`
}
```

#### 规则6: 新架构设计评审

**提交架构设计时必须包含**：

1. **ER图**：实体关系图
2. **数据流图**：数据如何流动
3. **状态图**：状态转换逻辑
4. **同步机制**：如何保持一致性
5. **异常处理**：失败如何恢复

### 2.3 架构设计模板

```markdown
# [功能名称] 架构设计

## 1. 实体关系

\`\`\`
┌──────────────┐         ┌──────────────┐
│    EntityA   │────────>│    EntityB   │
└──────────────┘         └──────────────┘
     关联字段: EntityBID
     关联类型: N:1
     同步方式: 实时/定时/手动
\`\`\`

## 2. 数据流

\`\`\`
A创建 ──> B创建 ──> 同步 ──> 更新
\`\`\`

## 3. 一致性保证

- 冲突检测：[方法]
- 冲突解决：[策略]
- 失败重试：[机制]

## 4. 边界场景

- A删除时B如何处理？
- A和B同时更新如何处理？
- 网络分区时如何处理？
```

### 2.4 应对流程

**发现架构缺陷时**：

1. **评估**：缺陷严重程度（P0/P1/P2）
2. **决策**：修复 / 遗留技术债 / 重构
3. **设计**：如果修复，设计完整方案
4. **实施**：按优先级实施
5. **验证**：确保不会引入新问题

---

## 三、层级职责清晰性

### 3.1 问题描述

Repository层承担Service层职责：
- 包含业务逻辑
- 数据转换和聚合
- 跨表join操作

### 3.2 预防措施

#### 规则7: 分层职责明确定义

| 层级 | 职责 | 不应该做 |
|------|------|----------|
| **API层** | 参数验证、DTO转换、响应格式 | 业务逻辑 |
| **Service层** | 业务逻辑、事务协调、跨Repository操作 | 直接访问数据库 |
| **Repository层** | 数据CRUD、查询构建 | 业务逻辑 |
| **Model层** | 数据结构定义 | 复杂方法 |

#### 规则8: Repository代码规范

```go
// ❌ 错误：Repository包含业务逻辑
func (r *BookRepository) GetBookWithAuthor(ctx context.Context, id string) (*BookDetailDTO, error) {
    // 跨表查询
    // 数据转换
    // 业务逻辑
}

// ✅ 正确：Repository只负责数据查询
func (r *BookRepository) GetByID(ctx context.Context, id string) (*Book, error) {
    var book Book
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&book)
    return &book, err
}

// Service层负责组装
func (s *BookService) GetBookDetail(ctx context.Context, id string) (*BookDetailDTO, error) {
    book, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    author, err := s.authorRepo.GetByID(ctx, book.AuthorID)
    if err != nil {
        return nil, err
    }
    return s.toDTO(book, author), nil
}
```

#### 规则9: 代码审查检查点

**审查Repository代码时检查**：

- [ ] 是否有业务逻辑？
- [ ] 是否跨多个Collection操作？
- [ ] 是否有复杂的数据转换？
- [ ] 方法名是否以`GetBy`/`Find`/`Create`/`Update`/`Delete`开头？

**如果任何一项为"是"，应该移到Service层**

### 3.3 检测工具

#### Repository职责检查脚本

```bash
# scripts/check-repository-responsibility.sh
#!/bin/bash

# 检查Repository中的可疑方法名
echo "检查Repository中不应存在的方法..."
grep -r "func.*Get.*With" --include="*_repository.go" repository/
grep -r "func.*Calculate" --include="*_repository.go" repository/
grep -r "func.*Process" --include="*_repository.go" repository/

# 检查Repository中的跨Collection操作
echo "检查跨Collection操作..."
grep -r "collection.*Find.*collection" --include="*_repository.go" repository/
```

---

## 四、版本管理需求

### 4.1 问题描述

缺少版本管理导致：
- 无法追溯历史变更
- 读者无法查看新旧差异
- 作者无法回滚错误修改

### 4.2 预防措施

#### 规则10: 内容类实体必须考虑版本管理

**以下实体类型需要版本管理**：

| 实体类型 | 是否需要版本管理 | 示例 |
|----------|------------------|------|
| 用户生成内容 | ✅ 是 | 文章、章节、评论 |
| 系统配置 | ✅ 是 | 配置文件、规则 |
| 业务单据 | ✅ 是 | 订单、发布单 |
| 日志数据 | ❌ 否 | 访问日志 |
| 缓存数据 | ❌ 否 | Redis缓存 |

#### 规则11: 版本管理设计要素

```go
// 版本管理实体必备字段
type VersionedEntity struct {
    ID           string     `bson:"_id"`
    Version      int        `bson:"version"`           // 版本号
    PreviousID   *string    `bson:"previous_id"`       // 前一版本
    ChangeType   ChangeType `bson:"change_type"`       // 变更类型
    ChangeReason string     `bson:"change_reason"`     // 变更原因
    ChangedBy    string     `bson:"changed_by"`        // 变更人
    ChangedAt    time.Time  `bson:"changed_at"`        // 变更时间
    SnapshotRef  *string    `bson:"snapshot_ref"`      // 快照引用
    DiffData     *DiffData  `bson:"diff_data"`         // 差异数据
}
```

#### 规则12: 版本管理实现检查清单

- [ ] 版本号自动递增
- [ ] 版本链完整（可追溯）
- [ ] 快照存储或diff存储
- [ ] 版本对比功能
- [ ] 版本回滚功能
- [ ] 版本查询API

---

## 五、数据迁移安全

### 5.1 预防措施

#### 规则13: 数据结构变更必须评估影响

**变更前检查清单**：

- [ ] 现有数据量有多大？
- [ ] 数据类型是否兼容？
- [ ] 是否需要数据迁移？
- [ ] 迁移失败如何回滚？
- [ ] 是否影响正在运行的服务？

#### 规则14: 安全迁移流程

```go
// 安全迁移模式
// 1. 添加新字段（保留旧字段）
type Entity struct {
    OldField string `bson:"old_field"` // 保留
    NewField string `bson:"new_field"` // 新增
}

// 2. 数据迁移脚本
func MigrateEntity(ctx context.Context, db *mongo.Database) error {
    // 分批迁移，避免锁表
    batchSize := 1000
    skip := 0

    for {
        cursor, err := db.Collection("entities").
            Find(ctx, bson.M{"new_field": bson.M{"$exists": false}},
                options.Find().SetSkip(int64(skip)).SetLimit(int64(batchSize)))
        if err != nil {
            return err
        }

        var results []Entity
        if err = cursor.All(ctx, &results); err != nil {
            return err
        }

        if len(results) == 0 {
            break
        }

        for _, r := range results {
            r.NewField = convertOldField(r.OldField)
            db.Collection("entities").UpdateByID(ctx, r.ID, bson.M{"$set": r})
        }

        skip += batchSize
    }

    return nil
}

// 3. 验证迁移结果
func VerifyMigration(ctx context.Context, db *mongo.Database) error {
    count, _ := db.Collection("entities").CountDocuments(ctx, bson.M{"new_field": bson.M{"$exists": false}})
    if count > 0 {
        return fmt.Errorf("迁移不完整，还有 %d 条记录未迁移", count)
    }
    return nil
}

// 4. 确认后删除旧字段
```

---

## 六、文档与代码同步

### 6.1 问题描述

设计文档与实际代码不一致：
- 文档定义的状态与代码不同
- API文档与实际行为不符
- 类型定义过时

### 6.2 预防措施

#### 规则15: 代码与文档同步更新

**提交代码时检查**：

- [ ] 如果修改了API，是否更新了OpenAPI文档？
- [ ] 如果修改了类型，是否更新了类型文档？
- [ ] 如果修改了业务流程，是否更新了设计文档？
- [ ] 如果修改了配置，是否更新了配置说明？

#### 规则16: 文档更新作为Code Review的一部分

**Code Review Checklist**：

- [ ] 代码变更与设计文档一致
- [ ] 相关文档已更新
- [ ] API文档已同步
- [ ] 变更日志已记录

---

## 七、测试覆盖要求

### 7.1 预防措施

#### 规则17: 核心业务逻辑必须有测试

**测试覆盖率要求**：

| 代码类型 | 最低覆盖率 | 推荐覆盖率 |
|----------|------------|------------|
| Domain层 | 90% | 95% |
| Service层 | 85% | 90% |
| Repository层 | 80% | 85% |
| API层 | 70% | 80% |

#### 规则18: 关键流程必须有集成测试

**必须测试的流程**：

- 用户注册/登录
- 内容创建/发布
- 支付流程
- 数据同步
- 事务回滚

---

## 八、最佳实践总结

### 8.1 设计阶段

| 实践 | 说明 |
|------|------|
| ✅ 先写设计文档 | 不设计就编码是问题之源 |
| ✅ 设计评审 | 让其他人review设计 |
| ✅ POC验证 | 复杂功能先做概念验证 |
| ✅ 考虑边界情况 | 正常路径只是冰山一角 |

### 8.2 编码阶段

| 实践 | 说明 |
|------|------|
| ✅ TDD | 先写测试，再写代码 |
| ✅ 代码审查 | 所有代码必须review |
| ✅ 遵循规范 | 使用统一的代码风格 |
| ✅ 小步提交 | 频繁提交，每步可运行 |

### 8.3 测试阶段

| 实践 | 说明 |
|------|------|
| ✅ 单元测试 | 覆盖率达标 |
| ✅ 集成测试 | 验证模块协作 |
| ✅ 压力测试 | 验证性能指标 |
| ✅ 安全测试 | 检查常见漏洞 |

### 8.4 发布阶段

| 实践 | 说明 |
|------|------|
| ✅ 灰度发布 | 逐步放量 |
| ✅ 监控告警 | 实时关注指标 |
| ✅ 回滚预案 | 出问题能快速恢复 |
| ✅ 发布总结 | 记录经验教训 |

---

## 九、错误应对决策树

```
                         发现问题
                            │
                            ▼
                    ┌───────────────┐
                    │  严重程度评估  │
                    └───────────────┘
                            │
            ┌───────────────┼───────────────┐
            ▼               ▼               ▼
         P0紧急           P1重要          P2一般
            │               │               │
            ▼               ▼               ▼
    ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
    │ 立即停止开发  │  │ 计划修复     │  │ 记录技术债   │
    │ 设计完整方案  │  │ 分配优先级   │  │ 有时间再处理 │
    │ 立即实施     │  │ 按计划实施   │  │              │
    └──────────────┘  └──────────────┘  └──────────────┘
```

---

## 十、工具和脚本

### 10.1 代码质量检查脚本

```bash
#!/bin/bash
# scripts/quality-check.sh

echo "==== Qingyu Backend 代码质量检查 ===="

# 1. 类型一致性检查
echo "1. 检查类型一致性..."
bash scripts/check-type-consistency.sh

# 2. Repository职责检查
echo "2. 检查Repository职责..."
bash scripts/check-repository-responsibility.sh

# 3. 测试覆盖率检查
echo "3. 检查测试覆盖率..."
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total

# 4. 代码格式检查
echo "4. 检查代码格式..."
gofmt -l .

# 5. 静态分析
echo "5. 静态分析..."
go vet ./...

echo "==== 检查完成 ===="
```

### 10.2 CI/CD集成

```yaml
# .github/workflows/quality-check.yml
name: Code Quality Check

on:
  pull_request:
    branches: [main, dev]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Run quality checks
        run: |
          bash scripts/quality-check.sh

      - name: Check coverage
        run: |
          go test -coverprofile=coverage.out ./...
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$COVERAGE < 80" | bc -l) )); then
            echo "覆盖率 $COVERAGE% 低于要求80%"
            exit 1
          fi
```

---

## 十一、持续改进

### 11.1 定期审查

| 活动 | 频率 | 目的 |
|------|------|------|
| 代码审查 | 每次提交 | 保证代码质量 |
| 架构审查 | 每月 | 发现架构问题 |
| 技术债清理 | 每季度 | 减少技术债积累 |
| 最佳实践更新 | 每半年 | 保持文档最新 |

### 11.2 经验教训记录

**模板**：

```markdown
# 问题记录

## 问题描述
[描述发现的问题]

## 根本原因
[分析根本原因]

## 影响范围
[评估影响]

## 解决方案
[实施的解决方案]

## 预防措施
[如何防止再次发生]

## 相关Issue
[链接到相关Issue]
```

---

## 十二、相关文档

- [实施路线图](../plans/2026-03-06-p0-implementation-roadmap.md)
- [代码规范](./coding/)
- [API设计规范](./api/)
- [测试规范](./testing/)

---

**文档版本**: 1.0
**最后更新**: 2026-03-06
**下次审查**: 2026-06-06（每季度审查）
