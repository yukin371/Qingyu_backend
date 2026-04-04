# 后端层级设计文档规划

> **版本**: v1.0
> **创建日期**: 2026-03-19
> **状态**: 进行中

---

## 1. 背景

为后端每个层级创建设计说明文档，方便后续开发时参考已有设计，避免重复设计，保持架构一致性。

---

## 2. 目标

- 为每个层级提供完整的设计指南（5-10页）
- 包含职责边界、命名规范、设计模式、接口规范、测试策略、代码示例
- 两处存放：详细版在 `docs/standards/`，快速参考版在各层级目录内

---

## 3. 文档结构

### 3.1 详细版文档（`docs/standards/`）

```
docs/standards/
├── layer-models.md          # Models层设计说明
├── layer-repository.md      # Repository层设计说明
├── layer-service.md         # Service层设计说明
├── layer-api.md             # API层设计说明
├── layer-dto.md             # DTO层设计说明
├── layer-middleware.md      # Middleware层设计说明
├── layer-router.md          # Router层设计说明
├── layer-config.md          # Config层设计说明
└── layer-pkg.md             # PKG层设计说明
```

### 3.2 快速参考版（各层级目录内）

```
models/README.md             # Models层快速参考
repository/README.md         # Repository层快速参考
service/README.md            # Service层快速参考
api/v1/README.md             # API层快速参考（已存在，需更新）
models/dto/README.md         # DTO层快速参考
internal/middleware/README.md # Middleware层快速参考
router/README.md             # Router层快速参考
config/README.md             # Config层快速参考
pkg/README.md                # PKG层快速参考
```

---

## 4. 文档模板

### 4.1 详细版模板

```markdown
# X层设计说明

> 版本: v1.0 | 创建日期: YYYY-MM-DD

## 1. 职责边界与依赖关系
- 该层级的职责定义
- 与上下层级的交互边界
- 依赖关系图

## 2. 命名与代码规范
- 文件命名规范
- 函数/方法命名规范
- 目录组织规范
- 导出规则

## 3. 设计模式与最佳实践
- 该层级常用的设计模式
- 推荐实现方式
- 反模式警示

## 4. 接口与契约规范
- 输入参数规范
- 返回值规范
- 错误码与错误处理

## 5. 测试策略
- 单元测试编写指南
- Mock/Stub使用
- 测试覆盖率要求

## 6. 完整代码示例
- 典型场景示例
- 完整实现参考
```

### 4.2 快速参考版模板

```markdown
# X层快速参考

## 职责
- 一句话职责说明

## 目录结构
```
layer/
├── module1/
└── module2/
```

## 命名规范
| 类型 | 规范 | 示例 |
|------|------|------|
| 文件 | xxx | xxx.go |

## 快速示例
```go
// 最小可行示例
```

## 详见
完整设计文档: [docs/standards/layer-xxx.md](../docs/standards/layer-xxx.md)
```

---

## 5. 分批实施计划

### 5.1 批次划分

| 批次 | 层级 | 理由 |
|------|------|------|
| **P0 核心** | Models → Repository → Service | 数据模型是基础，Repository 和 Service 是业务核心 |
| **P1 接口** | API → DTO | 承接 Service，对外接口层 |
| **P2 基础设施** | Middleware → Router → Config | 横切关注点和配置 |
| **P3 工具** | PKG → Utils | 通用工具包 |

### 5.2 P0 批次详细规划

#### Models 层

| 章节 | 内容要点 |
|------|----------|
| 职责边界 | 数据结构定义、MongoDB模型映射、字段验证标签 |
| 命名规范 | 文件名小写、结构体 PascalCase、字段 camelCase + bson tag |
| 设计模式 | 基础模型嵌入（BaseModel）、枚举常量定义、软删除模式 |
| 接口规范 | 模型不直接暴露，通过 DTO 转换后输出 |
| 测试策略 | 模型验证测试、JSON序列化测试 |
| 代码示例 | Book、User、Chapter 等核心模型完整示例 |

#### Repository 层

| 章节 | 内容要点 |
|------|----------|
| 职责边界 | 数据库 CRUD 操作、查询构建、缓存策略 |
| 命名规范 | 接口 `IXxxRepository`、实现 `xxxRepository`、文件按领域划分 |
| 设计模式 | Repository 接口模式、Query Builder 模式、ID 转换器 |
| 接口规范 | 返回 `(*Model, error)`、使用自定义错误类型、分页统一封装 |
| 测试策略 | 集成测试为主、使用 test container 或 mock 接口 |
| 代码示例 | BaseRepository 实现、典型 CRUD 操作、复杂查询示例 |

#### Service 层

| 章节 | 内容要点 |
|------|----------|
| 职责边界 | 业务逻辑编排、跨 Repository 协调、DTO 转换 |
| 命名规范 | 接口 `IXxxService`、实现 `xxxService`、按业务模块划分 |
| 设计模式 | 依赖注入、Facade 模式、事件驱动（可选） |
| 接口规范 | 输入 DTO、返回 DTO、统一错误包装 |
| 测试策略 | 单元测试 + Mock Repository、业务场景测试 |
| 代码示例 | 典型 Service 实现、跨模块调用、错误处理示例 |

### 5.3 后续批次概览

| 批次 | 层级 | 重点工作 |
|------|------|----------|
| **P1 接口** | API → DTO | 请求处理流程、响应格式规范、DTO 转换规则、Swagger 文档集成 |
| **P2 基础设施** | Middleware → Router → Config | 中间件链设计、路由组织规范、配置热加载 |
| **P3 工具** | PKG → Utils | 通用工具函数规范、错误处理工具、日志工具 |

---

## 6. 交付物清单

| 批次 | 详细版文档 | 快速参考版 |
|------|-----------|-----------|
| P0 | 3 份 | 3 份 |
| P1 | 2 份 | 2 份 |
| P2 | 3 份 | 3 份 |
| P3 | 2 份 | 2 份 |
| **总计** | **10 份** | **10 份** |

---

## 7. 进度跟踪

| 批次 | 层级 | 详细版 | 快速参考版 | 状态 |
|------|------|--------|-----------|------|
| P0 | Models | ✅ docs/standards/layer-models.md | ✅ models/README.md | **已完成** |
| P0 | Repository | ✅ docs/standards/layer-repository.md | ✅ repository/README.md | **已完成** |
| P0 | Service | ✅ docs/standards/layer-service.md | ✅ service/README.md | **已完成** |
| P1 | API | ✅ docs/standards/layer-api.md | ✅ api/v1/README.md | **已完成** |
| P1 | DTO | ✅ docs/standards/layer-dto.md | ✅ models/dto/README.md | **已完成** |
| P2 | Middleware | ✅ docs/standards/layer-middleware.md | ✅ internal/middleware/README.md | **已完成** |
| P2 | Router | ✅ docs/standards/layer-router.md | ✅ router/README.md | **已完成** |
| P2 | Config | ✅ docs/standards/layer-config.md | ✅ config/README.md | **已完成** |
| P3 | PKG (含Utils) | ✅ docs/standards/layer-pkg.md | ✅ pkg/README.md | **已完成** |
| P3 | Utils | - | - | (合并到PKG层) |

---

*最后更新：2026-03-19*
