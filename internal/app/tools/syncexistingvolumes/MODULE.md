# 卷文档大纲同步工具

`internal/app/tools/syncexistingvolumes` 是 `cmd/tools/sync_existing_volumes` 的 owner，用于为缺少 `outline_node_id` 的卷文档创建 arc 大纲并回写关联。

默认只预览；写入必须同时使用 `--apply` 和 `QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA=true`。回写失败会尝试删除刚创建的大纲，避免留下孤立节点。
