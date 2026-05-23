# 后端测试文档总入口

> 最后整理: 2026-05-22

`docs/testing/` 现在只做一件事: 作为测试文档的当前入口。这里保留“怎么跑、怎么验、阶段产物在哪里看”的索引，不再让不同类型的文档混在一起。

## 先看哪里

1. [README_当前指南.md](./README_当前指南.md) - 运行指南、快速开始、最佳实践、测试工具与专题使用说明。
2. [README_阶段报告.md](./README_阶段报告.md) - 进度、实施计划、覆盖率、完成总结等阶段性产物。
3. [README_历史总结.md](./README_历史总结.md) - 已结题、已收口的历史总结与阶段终稿。

## 分类规则

| 分类 | 适合放什么 | 判断关键词 |
|------|------------|------------|
| 当前指南 | 运行方式、联调步骤、快速开始、测试设计、最佳实践、账号/数据准备 | `指南`、`快速开始`、`使用说明`、`最佳实践`、`实施计划` |
| 阶段报告 | 进度跟踪、覆盖率统计、实施结果、阶段完成总结、验证报告 | `报告`、`追踪`、`完成总结`、`实施`、`验证` |
| 历史总结 | 已结束阶段的总评、日报/晚报、阶段终稿、复盘类文档 | `今日工作`、`最终报告`、`Phase`、`总结`、`完结` |

## 当前入口约定

- 如果是“我现在要怎么跑测试”，先看 [README_当前指南.md](./README_当前指南.md)。
- 如果是“这个阶段做到了什么”，先看 [README_阶段报告.md](./README_阶段报告.md)。
- 如果是“这个事情已经结束了，想回看结果”，先看 [README_历史总结.md](./README_历史总结.md)。

## 专题入口

| 专题 | 入口 |
|------|------|
| gRPC 联调 | [gRPC测试专题索引.md](./gRPC测试专题索引.md) |
| TODO 测试 | [TODO测试专题索引.md](./TODO测试专题索引.md) |
| AI 写作 | [AI写作测试专题索引.md](./AI写作测试专题索引.md) |
| CNNovel125K / 阅读 | [CNNovel125K测试专题索引.md](./CNNovel125K测试专题索引.md) |

## 边界说明

| 路径 | 角色 | 当前规则 |
|------|------|----------|
| `docs/testing/` | 当前测试文档入口 | 只放测试指南、阶段报告和历史总结的索引与正文 |
| `docs/standards/testing/` | 测试标准入口 | 负责分层规则、底层约束、标准口径 |
| `docs/test/` | 历史文档路径 | 只保留迁移说明，不再作为正文入口 |
| 根层 `test/` | 测试源码目录 | 放 Go 测试代码，不是文档目录 |

## 常用入口

| 用途 | 推荐入口 |
|------|----------|
| 测试规范总览 | [../standards/testing/README.md](../standards/testing/README.md) |
| 日常最佳实践 | [测试最佳实践.md](./测试最佳实践.md) |
| 集成测试 | [集成测试使用指南.md](./集成测试使用指南.md) / [集成测试完成总结.md](./集成测试完成总结.md) |
| API / Postman | [API测试指南.md](./API测试指南.md) / [Postman测试指南.md](./Postman测试指南.md) |
| gRPC 联调 | [gRPC集成测试指南.md](./gRPC集成测试指南.md) / [gRPC集成测试快速启动.md](./gRPC集成测试快速启动.md) |
| 数据与账号 | [本地测试数据初始化指南.md](./本地测试数据初始化指南.md) / [内测账号快速参考.md](./内测账号快速参考.md) |
| 测试专题 | [BOOKSTORE_READING_TESTS.md](./BOOKSTORE_READING_TESTS.md) / [AI_WRITING_ASSISTANT_TESTS.md](./AI_WRITING_ASSISTANT_TESTS.md) |

## 已知整理结果

- 旧版 README 中引用的 `测试架构设计规范V1.0.md`、`测试组织规范.md`、`性能测试规范.md` 当前并不在本目录。
- 对应规范入口已统一收口到 [../standards/testing/README.md](../standards/testing/README.md)。
- 文中出现的 `./test/...`、`../../test/...` 命令路径指向仓库根层测试源码目录，不等于 `docs/test/`。
- 如果后续要新增同类文档，优先先放进对应分类索引，再补正文。
- 当前已补四条高频专题入口，后续同主题文档优先挂到专题索引，不要继续只在根层平铺。

## 相关入口

- [../api/README.md](../api/README.md)
- [../implementation/README.md](../implementation/README.md)
- [../ops/README.md](../ops/README.md)
- [../deployment/README_TEST.md](../deployment/README_TEST.md)
