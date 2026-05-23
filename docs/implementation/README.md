# 后端实施文档入口

> 最后整理: 2026-05-22

## Page Role

- current-hub
- current-owner: `docs/implementation/`
- current-bounded: 当前后端实施文档主入口，负责落地过程、专项修复与阶段结果导航

## Recommended Read Path

1. 先按业务模块进入对应子目录。
2. 需要工程支撑时，再看 `infrastructure/`。
3. 需要验证结果时，再看 `testing/`。

## Quick Section Map

- 当前入口说明
- 当前真实目录
- Topic READMEs
- 根层重点文档
- 推荐阅读路径
- 边界说明

## Quick Takeaways

- 这是当前实施文档主入口，不是长期计划 owner。
- 长期治理计划在父仓 `docs/plans/submodules/backend/`。

## Skip Guide

- 只看长期计划：跳去父仓 `docs/plans/submodules/backend/README.md`。
- 只看当前架构事实：跳去 `../architecture/README.md`。

## 当前入口说明

- `Qingyu_backend/docs/implementation/` 保存已经落地的实施说明、阶段总结和专项修复记录。
- 这里不是跨仓库长期计划的 owner；新的长期治理方案请写到父仓 [../../../docs/plans/submodules/backend/README.md](../../../docs/plans/submodules/backend/README.md)。
- 历史上 README 中出现过的 `07-admin/`、`09-ai/`、`00进度指导/` 等目录当前并不在真实文件树中，不再作为可点击入口保留。

## 当前真实目录

| 分类 | 入口 | 说明 |
|------|------|------|
| 认证 | [01-auth/](./01-auth/) | 认证、权限、RBAC 相关实施 |
| 写作 | [02-writing/](./02-writing/) | 写作、导出、发布相关实施 |
| 阅读 | [03-reading/](./03-reading/) | 阅读统计等阅读侧实施 |
| 社交 | [04-social/](./04-social/) | 书单、社区等社交能力实施 |
| 通信 | [05-communication/](./05-communication/) | 通知与消息相关实施 |
| 书城 | [06-bookstore/](./06-bookstore/) | 书城接口、购买、搜索相关实施 |
| 财务 | [08-finance/](./08-finance/) | 支付与财务系统实施 |
| 基础设施 | [infrastructure/](./infrastructure/) | 中间件、路由、MCP 等工程支撑 |
| 测试专项 | [testing/](./testing/) | 实施过程中的验证报告与修复记录 |
| 历史归档 | [_archive/README.md](./_archive/README.md) | 旧资料归档说明 |

## Topic READMEs

- [01-auth/README.md](./01-auth/README.md)
- [02-writing/README.md](./02-writing/README.md)
- [03-reading/README.md](./03-reading/README.md)
- [04-social/README.md](./04-social/README.md)
- [05-communication/README.md](./05-communication/README.md)
- [06-bookstore/README.md](./06-bookstore/README.md)
- [08-finance/README.md](./08-finance/README.md)
- [infrastructure/README.md](./infrastructure/README.md)
- [testing/README.md](./testing/README.md)

## 根层重点文档

- [后端API实现状态核查.md](./后端API实现状态核查.md)
- [分阶段规划_文档索引.md](./分阶段规划_文档索引.md)
- [分阶段实施计划_v2.0.md](./分阶段实施计划_v2.0.md)
- [phase3_行动指南.md](./phase3_行动指南.md)
- [START_PYTHON_GRPC.md](./START_PYTHON_GRPC.md)
- [WINDOWS_QUICKSTART.md](./WINDOWS_QUICKSTART.md)
- [P0_REFACTOR_SUMMARY.md](./P0_REFACTOR_SUMMARY.md)
- [P1_ID_TYPE_ANALYSIS.md](./P1_ID_TYPE_ANALYSIS.md)
- [P2_CONTENT_SEPARATION_ANALYSIS.md](./P2_CONTENT_SEPARATION_ANALYSIS.md)

## 推荐阅读路径

1. 先看与当前需求最接近的模块目录，例如 [05-communication/](./05-communication/) 或 [06-bookstore/](./06-bookstore/)。
2. 如果涉及中间件、路由、运行环境，再看 [infrastructure/](./infrastructure/)。
3. 如果要确认历史修复是否已验证，再看 [testing/](./testing/) 和根层阶段总结文档。

## 边界说明

- Admin、AI 等能力目前没有在本目录中单独沉淀成 `07-admin/`、`09-ai/` 两个稳定专题；如需补齐，先确认是否已有更合适的 owner，再新增目录。
- 历史中文阶段目录当前记为 `TBD`，不在本 README 继续制造失效入口。
  - 确认路径：通过 `git log -- docs/implementation` 回溯旧结构，或检查父仓 `docs/plans/submodules/backend/` 是否已经承接对应阶段材料。
