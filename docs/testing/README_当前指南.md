# 当前指南索引

> 最后整理: 2026-05-22  
> 当前状态: `current-owner`

这一组文档回答的是“现在要怎么做测试”。如果你要启动测试、搭环境、查命令、找操作步骤，先从这里进。

## 推荐阅读顺序

1. [README.md](./README.md) - 总入口和分类规则。
2. [README_测试运行指南.md](./README_测试运行指南.md) - 常用 `go test` 运行方式与排障。
3. [QUICK_TEST_START.md](./QUICK_TEST_START.md) - 一分钟快速开始。
4. [测试最佳实践.md](./测试最佳实践.md) - 测试方法和协作约定。
5. [API测试指南.md](./API测试指南.md) - API 相关测试入口。
6. [集成测试使用指南.md](./集成测试使用指南.md) - 集成测试入口。
7. [../../scripts/testing/README.md](../../scripts/testing/README.md) - 测试脚本目录入口。
8. [../../migration/seeds/README.md](../../migration/seeds/README.md) - 测试数据与导入种子入口。

## Current Workflows

| 工作流 | 推荐入口 | 适用场景 |
|------|------|------|
| 基础运行 | [README_测试运行指南.md](./README_测试运行指南.md) / [QUICK_TEST_START.md](./QUICK_TEST_START.md) | 快速启动、常用 `go test`、覆盖率、调试 |
| API / Postman | [API测试指南.md](./API测试指南.md) / [API测试最佳实践指南.md](./API测试最佳实践指南.md) / [Postman测试指南.md](./Postman测试指南.md) | HTTP API 联调、接口验证、Postman 场景 |
| 集成与权限 | [集成测试使用指南.md](./集成测试使用指南.md) / [permission-test-setup.md](./permission-test-setup.md) | 多模块联调、权限测试、集成测试组织方式 |
| gRPC / AI 联调 | [gRPC测试专题索引.md](./gRPC测试专题索引.md) / [AI写作测试专题索引.md](./AI写作测试专题索引.md) | 后端与 AI 服务联调、AI 写作链路验证 |
| 数据准备 | [本地测试数据初始化指南.md](./本地测试数据初始化指南.md) / [内测账号快速参考.md](./内测账号快速参考.md) / [../../migration/seeds/README.md](../../migration/seeds/README.md) | 初始化本地数据、查测试账号、导入种子 |
| 阅读与数据集专题 | [CNNovel125K测试专题索引.md](./CNNovel125K测试专题索引.md) / [BOOKSTORE_READING_TESTS.md](./BOOKSTORE_READING_TESTS.md) | 阅读链路、数据集验证、专题资源查找 |
| TODO / 安全 / 发布专项 | [TODO测试专题索引.md](./TODO测试专题索引.md) / [SECURITY_TESTS_QUICK_GUIDE.md](./SECURITY_TESTS_QUICK_GUIDE.md) / [WRITER_PUBLISH_EXPORT_TESTS.md](./WRITER_PUBLISH_EXPORT_TESTS.md) | 特定功能面专项测试 |

## Topic Shortcuts

| 专题 | 快速入口 | 进一步阅读 |
|------|------|------|
| gRPC 联调 | [gRPC测试专题索引.md](./gRPC测试专题索引.md) | [gRPC集成测试指南.md](./gRPC集成测试指南.md) / [gRPC集成测试快速启动.md](./gRPC集成测试快速启动.md) / [gRPC测试结果示例.md](./gRPC测试结果示例.md) |
| AI 写作 | [AI写作测试专题索引.md](./AI写作测试专题索引.md) | [AI_WRITING_ASSISTANT_TESTS.md](./AI_WRITING_ASSISTANT_TESTS.md) / [AI微服务测试执行说明_2025-10-31.md](./AI微服务测试执行说明_2025-10-31.md) |
| TODO | [TODO测试专题索引.md](./TODO测试专题索引.md) | [测试TODO功能实施指南.md](./测试TODO功能实施指南.md) / [测试TODO快速参考卡.md](./测试TODO快速参考卡.md) |
| CNNovel125K / 阅读 | [CNNovel125K测试专题索引.md](./CNNovel125K测试专题索引.md) | [CNNovel125K快速开始.md](./CNNovel125K快速开始.md) / [执行步骤_CNNovel125K测试.md](./执行步骤_CNNovel125K测试.md) / [阅读端测试快速指南.md](./阅读端测试快速指南.md) |
| 共享与工具 | [../../scripts/testing/README.md](../../scripts/testing/README.md) | [共享服务测试文档.md](./共享服务测试文档.md) / [测试工具设计文档.md](./测试工具设计文档.md) / [测试资源总览.md](./测试资源总览.md) |

## Routing Rules

1. 这一类文档重心是“现在怎么跑、怎么验”，不是阶段结论。
2. 如果内容开始变成进度统计、验收结果或完成总结，应转到 [README_阶段报告.md](./README_阶段报告.md) 或 [README_历史总结.md](./README_历史总结.md)。
3. `AI_TESTS_REPORT.md`、`SECURITY_TESTS_SUMMARY.md` 这类“包含结果”的文档不再作为当前指南主入口。
4. 需要脚本或种子数据时，优先从 `../../scripts/testing/README.md` 与 `../../migration/seeds/README.md` 进入，不要直接记旧专题文件名。
