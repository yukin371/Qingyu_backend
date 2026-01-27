# 性能基线报告

**日期**: 2026-01-27
**测试工具**: wrk（待安装）
**测试配置**: 4线程, 100连接, 30秒

## 基线状态

⚠️ wrk工具未安装，性能基线待建立

当前基线文件：`test_results/baselines/baseline_20260127_152420.json`

```json
{
  "timestamp": "20260127_152420",
  "note": "wrk未安装，性能基线未建立",
  "tests": {}
}
```

## 下一步

请安装wrk后重新运行测试：

```bash
# Windows (WSL)
sudo apt-get install wrk

# 或从源码编译
git clone https://github.com/wg/wrk.git wrk
cd wrk
make

# 运行测试
cd Qingyu_backend
./scripts/performance_baseline.sh
```

## 环境信息

- **Go版本**: go1.25.1 windows/amd64
- **操作系统**: Windows
- **后端服务状态**: 未运行（测试时）
- **测试脚本**: `scripts/performance_baseline.sh`

## 注意事项

- 此基线用于后续重构阶段的性能对比
- 重构后P95延迟增加应<10%
- QPS下降应<5%
- 当前为占位符基线，待wrk安装后建立实际性能数据

## 相关文档

- 架构优化设计：`docs/plans/2026-01-26-block1-architecture-opt-design.md`
- 实施计划：`docs/plans/2026-01-27-p1-bookstore-implementation.md`

## 测试覆盖范围

性能基线测试将覆盖以下API端点：

1. **健康检查**: `GET /health`
2. **公共接口**: `GET /api/v1/books`
3. **认证接口**: `GET /api/v1/books/1`（需要JWT认证）
4. **管理接口**: `GET /api/v1/admin/books`（需要管理员权限）

## 基线指标

建立基线后，将记录以下关键指标：

- **QPS** (Queries Per Second): 每秒请求数
- **P50延迟**: 50%请求的响应时间
- **P95延迟**: 95%请求的响应时间
- **P99延迟**: 99%请求的响应时间
- **错误率**: 请求失败百分比

---

**创建时间**: 2026-01-27 15:24:20
**负责人**: 猫娘Kore
**任务状态**: ✅ 基线文件已创建（待wrk安装后重新测试）
