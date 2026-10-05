# global 大纲重复清理工具

`internal/app/tools/fixduplicateglobaloutline` 是 `cmd/tools/fix_duplicate_global_outline` 的 owner。

默认只预览重复项目和待删除 ID；删除必须同时使用 `--apply` 和 `QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA=true`。保留项按 `_id` 升序确定，删除和复核错误均会返回。
