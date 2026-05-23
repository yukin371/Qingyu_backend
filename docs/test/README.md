# Test Legacy Notice

`Qingyu_backend/docs/test` 是历史文档路径，不是当前测试文档入口，也不是仓库根层 `Qingyu_backend/test/` 测试源码目录。

## Archive Location

- [../archive/legacy-2026-05/test/](../archive/legacy-2026-05/test/)

## Current Source Of Truth

- 当前测试入口以 [../testing/README.md](../testing/README.md) 为准。
- `docs/testing/` 负责持续维护中的测试指南、测试报告和测试专题。
- `docs/standards/testing/` 负责测试规范、分层规则与底层约束。
- 仓库根层 `../../test/` 是 Go 测试源码目录，不属于文档 owner。

## Current Rule

- 不要继续把新的测试指南、测试总结或测试报告写入 `docs/test/`。
- 不要把 `docs/test/` 与仓库根层 `test/` 混为一谈。
- 若文档在回答“怎么跑测试”，写到 [../testing/README.md](../testing/README.md) 所属体系。
- 若文档在回答“测试必须遵循什么规则”，写到 [../standards/testing/README.md](../standards/testing/README.md) 所属体系。
