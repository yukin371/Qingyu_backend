# 测试脚本目录入口

> 最后整理: 2026-05-22

本目录保存测试准备、集成测试、gRPC 联调与辅助校验脚本。当前目录入口以本页为准，不再由单个 `README_*.md` 文件隐性充当总入口。

## Recommended Read Path

1. [README_gRPC测试.md](./README_gRPC测试.md) - gRPC 集成测试脚本说明
2. [README_Python脚本说明.md](./README_Python脚本说明.md) - Python 辅助脚本说明

## Key Script Groups

- `test_grpc_integration.bat`, `run_grpc_tests.bat`: gRPC 集成测试
- `setup_integration_tests.py`, `run_tests.py`, `cleanup_database.py`: Python 测试辅助脚本
- `ensure_test_data.*`, `prepare_test_data.*`, `import_novels.bat`: 测试数据准备
- `quick_verify.*`, `publication_flow_smoke.py`: 快速校验与专项 smoke

## Boundary

- `Qingyu_backend/docs/testing/README.md`: 当前测试文档总入口
- `Qingyu_backend/migration/seeds/README.md`: 测试数据和导入种子入口

## Rule

1. 新增测试脚本时，优先在本页补目录级说明。
2. 具体专题说明可继续放到子 README 或专题文档，但不应替代本页成为目录默认入口。
