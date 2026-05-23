# 阶段报告索引

> 最后整理: 2026-05-22  
> 当前状态: `current-bounded`

这一组文档回答的是“这个阶段做了什么、进展到哪一步、验证结果如何”。它们比当前指南更偏结果和状态追踪。

## 推荐阅读顺序

1. [README.md](./README.md) - 先确认分类边界。
2. [测试实施进度报告_2025-10-27.md](./测试实施进度报告_2025-10-27.md) - 了解阶段性推进情况。
3. [测试覆盖率追踪报告.md](./测试覆盖率追踪报告.md) - 看覆盖率和追踪口径。
4. [API层测试实施计划.md](./API层测试实施计划.md) - 了解 API 测试阶段计划。
5. [TEST_FILES_SUMMARY.md](./TEST_FILES_SUMMARY.md) - 看文件级别的测试成果汇总。

## Report Buckets

| 分组 | 代表文档 | 内容类型 |
|------|------|------|
| 基础建设与追踪 | [TESTING_SETUP_COMPLETE.md](./TESTING_SETUP_COMPLETE.md) / [TEST_FILES_SUMMARY.md](./TEST_FILES_SUMMARY.md) / [测试覆盖率追踪报告.md](./测试覆盖率追踪报告.md) / [测试实施进度报告_2025-10-27.md](./测试实施进度报告_2025-10-27.md) | 测试基础设施、文件覆盖、阶段进度 |
| API / E2E / 社交专项 | [API层测试实施计划.md](./API层测试实施计划.md) / [API测试完成总结.md](./API测试完成总结.md) / [E2E_TEST_COMPLETE_REPORT.md](./E2E_TEST_COMPLETE_REPORT.md) / [SOCIAL_API_IMPLEMENTATION_SUMMARY.md](./SOCIAL_API_IMPLEMENTATION_SUMMARY.md) | HTTP API、E2E、社交接口专项结果 |
| gRPC / AI / 数据转换 | [gRPC集成测试执行报告_2025-10-31.md](./gRPC集成测试执行报告_2025-10-31.md) / [gRPC数据转换层优化报告_2025-10-31.md](./gRPC数据转换层优化报告_2025-10-31.md) / [数据转换层修复验证报告_2025-10-31.md](./数据转换层修复验证报告_2025-10-31.md) / [AI_TESTS_REPORT.md](./AI_TESTS_REPORT.md) / [AI_WRITING_TESTS_SUMMARY.md](./AI_WRITING_TESTS_SUMMARY.md) | AI 联调、gRPC 转换、修复验证 |
| 数据集 / TODO / 安全专题 | [CNNovel125K测试报告.md](./CNNovel125K测试报告.md) / [CNNovel125K功能总结.md](./CNNovel125K功能总结.md) / [TODO功能测试补充报告_2025-10-31.md](./TODO功能测试补充报告_2025-10-31.md) / [TODO测试补充状态_2025-10-31.md](./TODO测试补充状态_2025-10-31.md) / [SECURITY_TESTS_SUMMARY.md](./SECURITY_TESTS_SUMMARY.md) | 专题验证结果与阶段总结 |
| 覆盖率与分支专题 | [Service层测试覆盖率总结报告_2025-1027.md](./Service层测试覆盖率总结报告_2025-1027.md) / [2026-02-10-user-service-coverage-report.md](./2026-02-10-user-service-coverage-report.md) / [2026-02-10-user-service-tdd-migration-report.md](./2026-02-10-user-service-tdd-migration-report.md) / [2026-02-12-branch-merge-test-report.md](./2026-02-12-branch-merge-test-report.md) | 覆盖率、迁移验证、分支合并结果 |
| 汇总收口 | [COMPLETE_TEST_SUMMARY.md](./COMPLETE_TEST_SUMMARY.md) / [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md) / [集成测试完成总结.md](./集成测试完成总结.md) / [测试CI-CD完善总结_2025-10-18.md](./测试CI-CD完善总结_2025-10-18.md) / [测试TODO文档体系建设完成报告.md](./测试TODO文档体系建设完成报告.md) | 阶段汇总、交付总结、专项收口 |

## Routing Rules

1. 这一类文档适合保留“当前阶段仍有参考价值”的结果和状态信息。
2. 如果文档已经成为完整结题材料，建议转入 [README_历史总结.md](./README_历史总结.md) 的语义层，不再回灌到当前指南。
3. 同一专题若同时存在“快速开始/执行步骤”和“测试报告/功能总结”，前者留在 [README_当前指南.md](./README_当前指南.md)，后者优先挂到这里。
