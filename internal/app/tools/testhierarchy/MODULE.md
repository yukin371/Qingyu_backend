# 大纲层级 API 集成工具

`internal/app/tools/testhierarchy` 是 `cmd/tools/test_hierarchy` 的 owner，用于验证注册、项目、卷、章节和大纲树 API 链路。

默认只预览；真实写入流程必须使用 `--apply`。服务地址可通过 `QINGYU_TEST_BASE_URL` 覆盖，默认 `http://localhost:9090/api/v1`。

## TODO

- `verifyStructure` 当前只输出“未实现”并跳过 MongoDB 结构复核。迁移阶段保留原工具的 API 验证范围；后续应接入只读数据库校验或独立 `test_outline_sync`，确认路径为补充可用测试数据库配置后增加集成测试。
