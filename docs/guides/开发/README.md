# 开发指南专题入口

> 最后整理: 2026-05-22

本目录保留开发协作相关的操作型指南，属于 `current-bounded` 专题区，不承担 API、测试或数据库的总入口职责。

## Current Documents

| 文档 | 状态 | 用途 |
|------|------|------|
| [API文档使用指南.md](./API文档使用指南.md) | `current-bounded` | API 文档生成、查看和使用方式 |
| [单元测试指南.md](./单元测试指南.md) | `current-bounded` | 开发阶段单元测试使用说明 |
| [缓存策略集成指南.md](./缓存策略集成指南.md) | `current-bounded` | 缓存策略接入说明 |

## Boundary

- `../README.md`：guides 总入口 owner。
- `../../api/`：API 参考和 Swagger 文档 owner。
- `../../testing/`：测试操作与报告入口 owner。
- `../../database/`：数据库专题 owner。

## Practical Rules

1. 如果文档在讲“开发时如何操作”，可以留在这里。
2. 如果文档已经演变成 API 参考、测试总指南或数据库总规范，应分别回写到对应 canonical owner。
3. 不要在本目录平铺阶段总结或完成报告。
