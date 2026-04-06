# 后端标准入口

> 更新日期: 2026-04-07  
> 适用范围: `Qingyu_backend/docs/standards/**`

本页是后端标准唯一入口，负责明确三件事：

1. 当前必须遵循的现行标准
2. 仅用于背景参考的历史补充
3. 新人和 AI 的最短阅读路径

## 新人 / AI 先看这 5 份

1. [2026-04-07-backend-architecture-and-module-standards.md](./2026-04-07-backend-architecture-and-module-standards.md)
2. [layer-service.md](./layer-service.md)
3. [layer-repository.md](./layer-repository.md)
4. [layer-api.md](./layer-api.md)
5. [layer-router.md](./layer-router.md)

这 5 份文档足以建立“分层边界 + 模块边界 + 高风险区域”的基础认知。

## 现行标准

### 项目级总则

- [2026-04-07-backend-architecture-and-module-standards.md](./2026-04-07-backend-architecture-and-module-standards.md)

### 分层标准（继续作为主体系）

- [layer-api.md](./layer-api.md)
- [layer-config.md](./layer-config.md)
- [layer-dto.md](./layer-dto.md)
- [layer-middleware.md](./layer-middleware.md)
- [layer-models.md](./layer-models.md)
- [layer-pkg.md](./layer-pkg.md)
- [layer-repository.md](./layer-repository.md)
- [layer-router.md](./layer-router.md)
- [layer-service.md](./layer-service.md)

### 横切标准

- [error_code_standard.md](./error_code_standard.md)
- [error_handling_guide.md](./error_handling_guide.md)
- [validation_standard.md](./validation_standard.md)
- [p0-error-prevention-and-best-practices.md](./p0-error-prevention-and-best-practices.md)

## 历史补充

以下文档保留用于历史追溯，不作为当前唯一依据：

- [backend-architecture-documentation-standard.md](./backend-architecture-documentation-standard.md)  
  说明: 与 2026-04-07 总则存在重叠，后续以新总则为准。
- [设计规范整理完成报告.md](./设计规范整理完成报告.md)  
  说明: 阶段性报告，不是规范正文。
- [archive/README.md](./archive/README.md) 及 `archive/*`  
  说明: 归档资料，只做背景参考。

## 与架构 / 审查文档的关系

`docs/standards/*` 负责“规则”，不负责“现状事实”。  
现状结构和风险请看：

- `Qingyu_backend/docs/architecture/*`
- `Qingyu_backend/docs/review/*`
- `Qingyu_backend/docs/issues/*`

当规则与代码现状冲突时，先在 `review` 或 `issues` 记录差异，再回写标准。

## 维护规则

1. 项目级约束统一写入 `2026-04-07-backend-architecture-and-module-standards.md`。
2. 某一层的实现规范只改对应 `layer-*.md`，避免在总则重复写细节。
3. 不再使用旧路径写法 `doc/standards/...`，统一使用 `docs/standards/...`。
4. 新增标准前先标注角色：现行标准、历史补充或阶段报告。
5. 继续保留的过时文档，必须在本 README 中标注“历史补充”。
