window.BENCHMARK_DATA = {
  "lastUpdate": 1772984557149,
  "repoUrl": "https://github.com/yukin371/Qingyu_backend",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "yukin3430@gmail.com",
            "name": "Alias",
            "username": "yukin3430"
          },
          "committer": {
            "email": "yukin3430@gmail.com",
            "name": "Alias",
            "username": "yukin3430"
          },
          "distinct": true,
          "id": "38323e4288dda267d6f1dbe9f2243028b7f0bcd1",
          "message": "fix(ci): 修复CI工作流中的lint和integration tests问题\n\n1. test.yml: 为benchmark job添加contents: write权限，修复gh-pages推送失败\n2. integration tests: 禁用strict logging检查，因为集成测试可能产生预期的错误日志\n3. comment_like_integration_test.go: 优化MongoDB连接超时，在不可用时优雅跳过\n4. batch_operation_atomic_false_test.go: 修复errcheck警告\n\nCo-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>",
          "timestamp": "2026-03-08T23:29:26+08:00",
          "tree_id": "5af00298c49ae3ba7fd04685a40b704bb40c1d17",
          "url": "https://github.com/yukin371/Qingyu_backend/commit/38323e4288dda267d6f1dbe9f2243028b7f0bcd1"
        },
        "date": 1772984556477,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared)",
            "value": 7119,
            "unit": "ns/op\t    2608 B/op\t      22 allocs/op",
            "extra": "231126 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - ns/op",
            "value": 7119,
            "unit": "ns/op",
            "extra": "231126 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - B/op",
            "value": 2608,
            "unit": "B/op",
            "extra": "231126 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - allocs/op",
            "value": 22,
            "unit": "allocs/op",
            "extra": "231126 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared)",
            "value": 6288,
            "unit": "ns/op\t    2608 B/op\t      22 allocs/op",
            "extra": "162704 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - ns/op",
            "value": 6288,
            "unit": "ns/op",
            "extra": "162704 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - B/op",
            "value": 2608,
            "unit": "B/op",
            "extra": "162704 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - allocs/op",
            "value": 22,
            "unit": "allocs/op",
            "extra": "162704 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared)",
            "value": 6541,
            "unit": "ns/op\t    2608 B/op\t      22 allocs/op",
            "extra": "153109 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - ns/op",
            "value": 6541,
            "unit": "ns/op",
            "extra": "153109 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - B/op",
            "value": 2608,
            "unit": "B/op",
            "extra": "153109 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - allocs/op",
            "value": 22,
            "unit": "allocs/op",
            "extra": "153109 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared)",
            "value": 7111,
            "unit": "ns/op\t    2608 B/op\t      22 allocs/op",
            "extra": "199482 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - ns/op",
            "value": 7111,
            "unit": "ns/op",
            "extra": "199482 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - B/op",
            "value": 2608,
            "unit": "B/op",
            "extra": "199482 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - allocs/op",
            "value": 22,
            "unit": "allocs/op",
            "extra": "199482 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared)",
            "value": 8337,
            "unit": "ns/op\t    2608 B/op\t      22 allocs/op",
            "extra": "193002 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - ns/op",
            "value": 8337,
            "unit": "ns/op",
            "extra": "193002 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - B/op",
            "value": 2608,
            "unit": "B/op",
            "extra": "193002 times\n4 procs"
          },
          {
            "name": "BenchmarkSuccess (Qingyu_backend/api/v1/shared) - allocs/op",
            "value": 22,
            "unit": "allocs/op",
            "extra": "193002 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth)",
            "value": 1515,
            "unit": "ns/op\t    1776 B/op\t      17 allocs/op",
            "extra": "789360 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - ns/op",
            "value": 1515,
            "unit": "ns/op",
            "extra": "789360 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - B/op",
            "value": 1776,
            "unit": "B/op",
            "extra": "789360 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "789360 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth)",
            "value": 1539,
            "unit": "ns/op\t    1776 B/op\t      17 allocs/op",
            "extra": "745036 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - ns/op",
            "value": 1539,
            "unit": "ns/op",
            "extra": "745036 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - B/op",
            "value": 1776,
            "unit": "B/op",
            "extra": "745036 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "745036 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth)",
            "value": 1535,
            "unit": "ns/op\t    1776 B/op\t      17 allocs/op",
            "extra": "735752 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - ns/op",
            "value": 1535,
            "unit": "ns/op",
            "extra": "735752 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - B/op",
            "value": 1776,
            "unit": "B/op",
            "extra": "735752 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "735752 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth)",
            "value": 1533,
            "unit": "ns/op\t    1776 B/op\t      17 allocs/op",
            "extra": "732566 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - ns/op",
            "value": 1533,
            "unit": "ns/op",
            "extra": "732566 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - B/op",
            "value": 1776,
            "unit": "B/op",
            "extra": "732566 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "732566 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth)",
            "value": 1523,
            "unit": "ns/op\t    1776 B/op\t      17 allocs/op",
            "extra": "711243 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - ns/op",
            "value": 1523,
            "unit": "ns/op",
            "extra": "711243 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - B/op",
            "value": 1776,
            "unit": "B/op",
            "extra": "711243 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionMiddleware (Qingyu_backend/internal/middleware/auth) - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "711243 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1467,
            "unit": "ns/op\t    1600 B/op\t      18 allocs/op",
            "extra": "790504 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1467,
            "unit": "ns/op",
            "extra": "790504 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1600,
            "unit": "B/op",
            "extra": "790504 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 18,
            "unit": "allocs/op",
            "extra": "790504 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1557,
            "unit": "ns/op\t    1600 B/op\t      18 allocs/op",
            "extra": "766639 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1557,
            "unit": "ns/op",
            "extra": "766639 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1600,
            "unit": "B/op",
            "extra": "766639 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 18,
            "unit": "allocs/op",
            "extra": "766639 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1600,
            "unit": "ns/op\t    1600 B/op\t      18 allocs/op",
            "extra": "777270 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1600,
            "unit": "ns/op",
            "extra": "777270 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1600,
            "unit": "B/op",
            "extra": "777270 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 18,
            "unit": "allocs/op",
            "extra": "777270 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1461,
            "unit": "ns/op\t    1600 B/op\t      18 allocs/op",
            "extra": "769900 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1461,
            "unit": "ns/op",
            "extra": "769900 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1600,
            "unit": "B/op",
            "extra": "769900 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 18,
            "unit": "allocs/op",
            "extra": "769900 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1469,
            "unit": "ns/op\t    1600 B/op\t      18 allocs/op",
            "extra": "773335 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1469,
            "unit": "ns/op",
            "extra": "773335 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1600,
            "unit": "B/op",
            "extra": "773335 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressionMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 18,
            "unit": "allocs/op",
            "extra": "773335 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 2520,
            "unit": "ns/op\t    1841 B/op\t      26 allocs/op",
            "extra": "451525 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 2520,
            "unit": "ns/op",
            "extra": "451525 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1841,
            "unit": "B/op",
            "extra": "451525 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 26,
            "unit": "allocs/op",
            "extra": "451525 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 2525,
            "unit": "ns/op\t    1841 B/op\t      26 allocs/op",
            "extra": "461907 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 2525,
            "unit": "ns/op",
            "extra": "461907 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1841,
            "unit": "B/op",
            "extra": "461907 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 26,
            "unit": "allocs/op",
            "extra": "461907 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 2541,
            "unit": "ns/op\t    1841 B/op\t      26 allocs/op",
            "extra": "450913 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 2541,
            "unit": "ns/op",
            "extra": "450913 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1841,
            "unit": "B/op",
            "extra": "450913 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 26,
            "unit": "allocs/op",
            "extra": "450913 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 2548,
            "unit": "ns/op\t    1841 B/op\t      26 allocs/op",
            "extra": "453163 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 2548,
            "unit": "ns/op",
            "extra": "453163 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1841,
            "unit": "B/op",
            "extra": "453163 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 26,
            "unit": "allocs/op",
            "extra": "453163 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 2612,
            "unit": "ns/op\t    1841 B/op\t      26 allocs/op",
            "extra": "450222 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 2612,
            "unit": "ns/op",
            "extra": "450222 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1841,
            "unit": "B/op",
            "extra": "450222 times\n4 procs"
          },
          {
            "name": "BenchmarkCORSMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 26,
            "unit": "allocs/op",
            "extra": "450222 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1250,
            "unit": "ns/op\t    1440 B/op\t      15 allocs/op",
            "extra": "874779 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1250,
            "unit": "ns/op",
            "extra": "874779 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1440,
            "unit": "B/op",
            "extra": "874779 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "874779 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1246,
            "unit": "ns/op\t    1440 B/op\t      15 allocs/op",
            "extra": "872559 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1246,
            "unit": "ns/op",
            "extra": "872559 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1440,
            "unit": "B/op",
            "extra": "872559 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "872559 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1294,
            "unit": "ns/op\t    1440 B/op\t      15 allocs/op",
            "extra": "826953 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1294,
            "unit": "ns/op",
            "extra": "826953 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1440,
            "unit": "B/op",
            "extra": "826953 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "826953 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1244,
            "unit": "ns/op\t    1440 B/op\t      15 allocs/op",
            "extra": "829974 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1244,
            "unit": "ns/op",
            "extra": "829974 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1440,
            "unit": "B/op",
            "extra": "829974 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "829974 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1240,
            "unit": "ns/op\t    1440 B/op\t      15 allocs/op",
            "extra": "889926 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1240,
            "unit": "ns/op",
            "extra": "889926 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1440,
            "unit": "B/op",
            "extra": "889926 times\n4 procs"
          },
          {
            "name": "BenchmarkErrorHandlerMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "889926 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1594,
            "unit": "ns/op\t    1440 B/op\t      15 allocs/op",
            "extra": "775345 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1594,
            "unit": "ns/op",
            "extra": "775345 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1440,
            "unit": "B/op",
            "extra": "775345 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "775345 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1587,
            "unit": "ns/op\t    1440 B/op\t      15 allocs/op",
            "extra": "739629 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1587,
            "unit": "ns/op",
            "extra": "739629 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1440,
            "unit": "B/op",
            "extra": "739629 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "739629 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1506,
            "unit": "ns/op\t    1440 B/op\t      15 allocs/op",
            "extra": "703461 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1506,
            "unit": "ns/op",
            "extra": "703461 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1440,
            "unit": "B/op",
            "extra": "703461 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "703461 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1552,
            "unit": "ns/op\t    1440 B/op\t      15 allocs/op",
            "extra": "667488 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1552,
            "unit": "ns/op",
            "extra": "667488 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1440,
            "unit": "B/op",
            "extra": "667488 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "667488 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin)",
            "value": 1626,
            "unit": "ns/op\t    1440 B/op\t      15 allocs/op",
            "extra": "627564 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - ns/op",
            "value": 1626,
            "unit": "ns/op",
            "extra": "627564 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - B/op",
            "value": 1440,
            "unit": "B/op",
            "extra": "627564 times\n4 procs"
          },
          {
            "name": "BenchmarkRecoveryMiddleware (Qingyu_backend/internal/middleware/builtin) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "627564 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics)",
            "value": 327.5,
            "unit": "ns/op\t       3 B/op\t       1 allocs/op",
            "extra": "3621064 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 327.5,
            "unit": "ns/op",
            "extra": "3621064 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - B/op",
            "value": 3,
            "unit": "B/op",
            "extra": "3621064 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3621064 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics)",
            "value": 394.4,
            "unit": "ns/op\t       3 B/op\t       1 allocs/op",
            "extra": "3612801 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 394.4,
            "unit": "ns/op",
            "extra": "3612801 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - B/op",
            "value": 3,
            "unit": "B/op",
            "extra": "3612801 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3612801 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics)",
            "value": 329.5,
            "unit": "ns/op\t       3 B/op\t       1 allocs/op",
            "extra": "3737936 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 329.5,
            "unit": "ns/op",
            "extra": "3737936 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - B/op",
            "value": 3,
            "unit": "B/op",
            "extra": "3737936 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3737936 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics)",
            "value": 330.4,
            "unit": "ns/op\t       3 B/op\t       1 allocs/op",
            "extra": "3861136 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 330.4,
            "unit": "ns/op",
            "extra": "3861136 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - B/op",
            "value": 3,
            "unit": "B/op",
            "extra": "3861136 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3861136 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics)",
            "value": 325.1,
            "unit": "ns/op\t       3 B/op\t       1 allocs/op",
            "extra": "3768454 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 325.1,
            "unit": "ns/op",
            "extra": "3768454 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - B/op",
            "value": 3,
            "unit": "B/op",
            "extra": "3768454 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordHttpRequest (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3768454 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics)",
            "value": 162.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8180648 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 162.3,
            "unit": "ns/op",
            "extra": "8180648 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8180648 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8180648 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics)",
            "value": 158.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7261357 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 158.7,
            "unit": "ns/op",
            "extra": "7261357 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7261357 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7261357 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics)",
            "value": 162.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8174433 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 162.6,
            "unit": "ns/op",
            "extra": "8174433 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8174433 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8174433 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics)",
            "value": 153.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8073234 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 153.9,
            "unit": "ns/op",
            "extra": "8073234 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8073234 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8073234 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics)",
            "value": 150.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8146296 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 150.5,
            "unit": "ns/op",
            "extra": "8146296 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8146296 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordDbQuery (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8146296 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics)",
            "value": 194.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5820948 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 194.8,
            "unit": "ns/op",
            "extra": "5820948 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5820948 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5820948 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics)",
            "value": 198.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5903972 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 198.8,
            "unit": "ns/op",
            "extra": "5903972 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5903972 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5903972 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics)",
            "value": 197.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6094528 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 197.3,
            "unit": "ns/op",
            "extra": "6094528 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6094528 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6094528 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics)",
            "value": 200.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6159526 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 200.2,
            "unit": "ns/op",
            "extra": "6159526 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6159526 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6159526 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics)",
            "value": 198.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6200918 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - ns/op",
            "value": 198.1,
            "unit": "ns/op",
            "extra": "6200918 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6200918 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSearch (Qingyu_backend/pkg/metrics) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6200918 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor)",
            "value": 3e-7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - ns/op",
            "value": 3e-7,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor)",
            "value": 2e-7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - ns/op",
            "value": 2e-7,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor)",
            "value": 2e-7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - ns/op",
            "value": 2e-7,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor)",
            "value": 2e-7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - ns/op",
            "value": 2e-7,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor)",
            "value": 5e-7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - ns/op",
            "value": 5e-7,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkOrphanedRecords (Qingyu_backend/pkg/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator)",
            "value": 239.1,
            "unit": "ns/op\t      40 B/op\t       3 allocs/op",
            "extra": "5073939 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - ns/op",
            "value": 239.1,
            "unit": "ns/op",
            "extra": "5073939 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - B/op",
            "value": 40,
            "unit": "B/op",
            "extra": "5073939 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5073939 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator)",
            "value": 219.2,
            "unit": "ns/op\t      40 B/op\t       3 allocs/op",
            "extra": "5375151 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - ns/op",
            "value": 219.2,
            "unit": "ns/op",
            "extra": "5375151 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - B/op",
            "value": 40,
            "unit": "B/op",
            "extra": "5375151 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5375151 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator)",
            "value": 203.7,
            "unit": "ns/op\t      40 B/op\t       3 allocs/op",
            "extra": "5848666 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - ns/op",
            "value": 203.7,
            "unit": "ns/op",
            "extra": "5848666 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - B/op",
            "value": 40,
            "unit": "B/op",
            "extra": "5848666 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5848666 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator)",
            "value": 203.3,
            "unit": "ns/op\t      40 B/op\t       3 allocs/op",
            "extra": "5896173 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - ns/op",
            "value": 203.3,
            "unit": "ns/op",
            "extra": "5896173 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - B/op",
            "value": 40,
            "unit": "B/op",
            "extra": "5896173 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5896173 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator)",
            "value": 203.3,
            "unit": "ns/op\t      40 B/op\t       3 allocs/op",
            "extra": "5879317 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - ns/op",
            "value": 203.3,
            "unit": "ns/op",
            "extra": "5879317 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - B/op",
            "value": 40,
            "unit": "B/op",
            "extra": "5879317 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateAmount (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5879317 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator)",
            "value": 4495,
            "unit": "ns/op\t    5178 B/op\t      67 allocs/op",
            "extra": "251956 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - ns/op",
            "value": 4495,
            "unit": "ns/op",
            "extra": "251956 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - B/op",
            "value": 5178,
            "unit": "B/op",
            "extra": "251956 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "251956 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator)",
            "value": 4556,
            "unit": "ns/op\t    5178 B/op\t      67 allocs/op",
            "extra": "260097 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - ns/op",
            "value": 4556,
            "unit": "ns/op",
            "extra": "260097 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - B/op",
            "value": 5178,
            "unit": "B/op",
            "extra": "260097 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "260097 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator)",
            "value": 5003,
            "unit": "ns/op\t    5178 B/op\t      67 allocs/op",
            "extra": "206746 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - ns/op",
            "value": 5003,
            "unit": "ns/op",
            "extra": "206746 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - B/op",
            "value": 5178,
            "unit": "B/op",
            "extra": "206746 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "206746 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator)",
            "value": 4483,
            "unit": "ns/op\t    5178 B/op\t      67 allocs/op",
            "extra": "261528 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - ns/op",
            "value": 4483,
            "unit": "ns/op",
            "extra": "261528 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - B/op",
            "value": 5178,
            "unit": "B/op",
            "extra": "261528 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "261528 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator)",
            "value": 4465,
            "unit": "ns/op\t    5178 B/op\t      67 allocs/op",
            "extra": "266910 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - ns/op",
            "value": 4465,
            "unit": "ns/op",
            "extra": "266910 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - B/op",
            "value": 5178,
            "unit": "B/op",
            "extra": "266910 times\n4 procs"
          },
          {
            "name": "BenchmarkValidatePhone (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "266910 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator)",
            "value": 2280,
            "unit": "ns/op\t    2318 B/op\t      28 allocs/op",
            "extra": "502102 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - ns/op",
            "value": 2280,
            "unit": "ns/op",
            "extra": "502102 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - B/op",
            "value": 2318,
            "unit": "B/op",
            "extra": "502102 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "502102 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator)",
            "value": 2272,
            "unit": "ns/op\t    2319 B/op\t      28 allocs/op",
            "extra": "489982 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - ns/op",
            "value": 2272,
            "unit": "ns/op",
            "extra": "489982 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - B/op",
            "value": 2319,
            "unit": "B/op",
            "extra": "489982 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "489982 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator)",
            "value": 2271,
            "unit": "ns/op\t    2318 B/op\t      28 allocs/op",
            "extra": "501069 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - ns/op",
            "value": 2271,
            "unit": "ns/op",
            "extra": "501069 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - B/op",
            "value": 2318,
            "unit": "B/op",
            "extra": "501069 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "501069 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator)",
            "value": 2275,
            "unit": "ns/op\t    2320 B/op\t      28 allocs/op",
            "extra": "513394 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - ns/op",
            "value": 2275,
            "unit": "ns/op",
            "extra": "513394 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - B/op",
            "value": 2320,
            "unit": "B/op",
            "extra": "513394 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "513394 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator)",
            "value": 2258,
            "unit": "ns/op\t    2317 B/op\t      28 allocs/op",
            "extra": "508466 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - ns/op",
            "value": 2258,
            "unit": "ns/op",
            "extra": "508466 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - B/op",
            "value": 2317,
            "unit": "B/op",
            "extra": "508466 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateStrongPassword (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "508466 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator)",
            "value": 9792,
            "unit": "ns/op\t   10005 B/op\t     127 allocs/op",
            "extra": "120963 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - ns/op",
            "value": 9792,
            "unit": "ns/op",
            "extra": "120963 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - B/op",
            "value": 10005,
            "unit": "B/op",
            "extra": "120963 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 127,
            "unit": "allocs/op",
            "extra": "120963 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator)",
            "value": 10898,
            "unit": "ns/op\t   10008 B/op\t     127 allocs/op",
            "extra": "120001 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - ns/op",
            "value": 10898,
            "unit": "ns/op",
            "extra": "120001 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - B/op",
            "value": 10008,
            "unit": "B/op",
            "extra": "120001 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 127,
            "unit": "allocs/op",
            "extra": "120001 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator)",
            "value": 9798,
            "unit": "ns/op\t   10006 B/op\t     127 allocs/op",
            "extra": "118626 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - ns/op",
            "value": 9798,
            "unit": "ns/op",
            "extra": "118626 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - B/op",
            "value": 10006,
            "unit": "B/op",
            "extra": "118626 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 127,
            "unit": "allocs/op",
            "extra": "118626 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator)",
            "value": 9779,
            "unit": "ns/op\t   10004 B/op\t     127 allocs/op",
            "extra": "120528 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - ns/op",
            "value": 9779,
            "unit": "ns/op",
            "extra": "120528 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - B/op",
            "value": 10004,
            "unit": "B/op",
            "extra": "120528 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 127,
            "unit": "allocs/op",
            "extra": "120528 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator)",
            "value": 9901,
            "unit": "ns/op\t   10008 B/op\t     127 allocs/op",
            "extra": "120813 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - ns/op",
            "value": 9901,
            "unit": "ns/op",
            "extra": "120813 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - B/op",
            "value": 10008,
            "unit": "B/op",
            "extra": "120813 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateComplex (Qingyu_backend/pkg/validator) - allocs/op",
            "value": 127,
            "unit": "allocs/op",
            "extra": "120813 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base)",
            "value": 19.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "59722388 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 19.26,
            "unit": "ns/op",
            "extra": "59722388 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "59722388 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "59722388 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base)",
            "value": 19.24,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "60685005 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 19.24,
            "unit": "ns/op",
            "extra": "60685005 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "60685005 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "60685005 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base)",
            "value": 19.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "62082076 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 19.27,
            "unit": "ns/op",
            "extra": "62082076 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "62082076 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "62082076 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base)",
            "value": 19.38,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "61792062 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 19.38,
            "unit": "ns/op",
            "extra": "61792062 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "61792062 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "61792062 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base)",
            "value": 20.76,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "62150690 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 20.76,
            "unit": "ns/op",
            "extra": "62150690 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "62150690 times\n4 procs"
          },
          {
            "name": "BenchmarkParseID (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "62150690 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base)",
            "value": 184.8,
            "unit": "ns/op\t     112 B/op\t       2 allocs/op",
            "extra": "6727774 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 184.8,
            "unit": "ns/op",
            "extra": "6727774 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "6727774 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6727774 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base)",
            "value": 184.2,
            "unit": "ns/op\t     112 B/op\t       2 allocs/op",
            "extra": "6516530 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 184.2,
            "unit": "ns/op",
            "extra": "6516530 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "6516530 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6516530 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base)",
            "value": 183,
            "unit": "ns/op\t     112 B/op\t       2 allocs/op",
            "extra": "6555451 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 183,
            "unit": "ns/op",
            "extra": "6555451 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "6555451 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6555451 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base)",
            "value": 184,
            "unit": "ns/op\t     112 B/op\t       2 allocs/op",
            "extra": "6533637 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 184,
            "unit": "ns/op",
            "extra": "6533637 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "6533637 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6533637 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base)",
            "value": 183.1,
            "unit": "ns/op\t     112 B/op\t       2 allocs/op",
            "extra": "6445618 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 183.1,
            "unit": "ns/op",
            "extra": "6445618 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "6445618 times\n4 procs"
          },
          {
            "name": "BenchmarkParseIDs (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6445618 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base)",
            "value": 18.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "63119365 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 18.97,
            "unit": "ns/op",
            "extra": "63119365 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "63119365 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "63119365 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base)",
            "value": 18.96,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "63182496 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 18.96,
            "unit": "ns/op",
            "extra": "63182496 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "63182496 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "63182496 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base)",
            "value": 19.92,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "63196593 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 19.92,
            "unit": "ns/op",
            "extra": "63196593 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "63196593 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "63196593 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base)",
            "value": 19.36,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "63615217 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 19.36,
            "unit": "ns/op",
            "extra": "63615217 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "63615217 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "63615217 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base)",
            "value": 19.11,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "60680577 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - ns/op",
            "value": 19.11,
            "unit": "ns/op",
            "extra": "60680577 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "60680577 times\n4 procs"
          },
          {
            "name": "BenchmarkIDToHex (Qingyu_backend/repository/mongodb/base) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "60680577 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor)",
            "value": 54.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21771465 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - ns/op",
            "value": 54.29,
            "unit": "ns/op",
            "extra": "21771465 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21771465 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21771465 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor)",
            "value": 55.44,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22222713 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - ns/op",
            "value": 55.44,
            "unit": "ns/op",
            "extra": "22222713 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "22222713 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "22222713 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor)",
            "value": 55.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21090608 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - ns/op",
            "value": 55.83,
            "unit": "ns/op",
            "extra": "21090608 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21090608 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21090608 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor)",
            "value": 55.25,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22130436 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - ns/op",
            "value": 55.25,
            "unit": "ns/op",
            "extra": "22130436 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "22130436 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "22130436 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor)",
            "value": 55.89,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18672388 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - ns/op",
            "value": 55.89,
            "unit": "ns/op",
            "extra": "18672388 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18672388 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordSlowQuery (Qingyu_backend/repository/mongodb/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18672388 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor)",
            "value": 85.81,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15096621 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - ns/op",
            "value": 85.81,
            "unit": "ns/op",
            "extra": "15096621 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15096621 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15096621 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor)",
            "value": 81.05,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14626191 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - ns/op",
            "value": 81.05,
            "unit": "ns/op",
            "extra": "14626191 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14626191 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14626191 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor)",
            "value": 79.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14998352 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - ns/op",
            "value": 79.31,
            "unit": "ns/op",
            "extra": "14998352 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14998352 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14998352 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor)",
            "value": 79.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15030538 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - ns/op",
            "value": 79.66,
            "unit": "ns/op",
            "extra": "15030538 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15030538 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15030538 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor)",
            "value": 79.23,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15131914 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - ns/op",
            "value": 79.23,
            "unit": "ns/op",
            "extra": "15131914 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15131914 times\n4 procs"
          },
          {
            "name": "BenchmarkRecordQueryDuration (Qingyu_backend/repository/mongodb/monitor) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15131914 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies)",
            "value": 253.2,
            "unit": "ns/op\t     168 B/op\t       3 allocs/op",
            "extra": "4747900 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - ns/op",
            "value": 253.2,
            "unit": "ns/op",
            "extra": "4747900 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - B/op",
            "value": 168,
            "unit": "B/op",
            "extra": "4747900 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4747900 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies)",
            "value": 267.6,
            "unit": "ns/op\t     168 B/op\t       3 allocs/op",
            "extra": "4742649 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - ns/op",
            "value": 267.6,
            "unit": "ns/op",
            "extra": "4742649 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - B/op",
            "value": 168,
            "unit": "B/op",
            "extra": "4742649 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4742649 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies)",
            "value": 257.2,
            "unit": "ns/op\t     168 B/op\t       3 allocs/op",
            "extra": "4708867 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - ns/op",
            "value": 257.2,
            "unit": "ns/op",
            "extra": "4708867 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - B/op",
            "value": 168,
            "unit": "B/op",
            "extra": "4708867 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4708867 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies)",
            "value": 279.2,
            "unit": "ns/op\t     168 B/op\t       3 allocs/op",
            "extra": "3874890 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - ns/op",
            "value": 279.2,
            "unit": "ns/op",
            "extra": "3874890 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - B/op",
            "value": 168,
            "unit": "B/op",
            "extra": "3874890 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "3874890 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies)",
            "value": 258.4,
            "unit": "ns/op\t     168 B/op\t       3 allocs/op",
            "extra": "4710082 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - ns/op",
            "value": 258.4,
            "unit": "ns/op",
            "extra": "4710082 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - B/op",
            "value": 168,
            "unit": "B/op",
            "extra": "4710082 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckFile (Qingyu_backend/scripts/check-dependencies) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4710082 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai)",
            "value": 7.009,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "171444157 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - ns/op",
            "value": 7.009,
            "unit": "ns/op",
            "extra": "171444157 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "171444157 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "171444157 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai)",
            "value": 6.979,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "171443449 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - ns/op",
            "value": 6.979,
            "unit": "ns/op",
            "extra": "171443449 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "171443449 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "171443449 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai)",
            "value": 7.014,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "171940759 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - ns/op",
            "value": 7.014,
            "unit": "ns/op",
            "extra": "171940759 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "171940759 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "171940759 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai)",
            "value": 7.016,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "171714194 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - ns/op",
            "value": 7.016,
            "unit": "ns/op",
            "extra": "171714194 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "171714194 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "171714194 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai)",
            "value": 7.128,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "162167911 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - ns/op",
            "value": 7.128,
            "unit": "ns/op",
            "extra": "162167911 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "162167911 times\n4 procs"
          },
          {
            "name": "BenchmarkCircuitBreaker_AllowRequest (Qingyu_backend/service/ai) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "162167911 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync)",
            "value": 4492,
            "unit": "ns/op\t    2000 B/op\t      54 allocs/op",
            "extra": "259506 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - ns/op",
            "value": 4492,
            "unit": "ns/op",
            "extra": "259506 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - B/op",
            "value": 2000,
            "unit": "B/op",
            "extra": "259506 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - allocs/op",
            "value": 54,
            "unit": "allocs/op",
            "extra": "259506 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync)",
            "value": 4490,
            "unit": "ns/op\t    2016 B/op\t      54 allocs/op",
            "extra": "259855 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - ns/op",
            "value": 4490,
            "unit": "ns/op",
            "extra": "259855 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - B/op",
            "value": 2016,
            "unit": "B/op",
            "extra": "259855 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - allocs/op",
            "value": 54,
            "unit": "allocs/op",
            "extra": "259855 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync)",
            "value": 4558,
            "unit": "ns/op\t    2032 B/op\t      54 allocs/op",
            "extra": "261205 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - ns/op",
            "value": 4558,
            "unit": "ns/op",
            "extra": "261205 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "261205 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - allocs/op",
            "value": 54,
            "unit": "allocs/op",
            "extra": "261205 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync)",
            "value": 4519,
            "unit": "ns/op\t    2016 B/op\t      54 allocs/op",
            "extra": "250287 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - ns/op",
            "value": 4519,
            "unit": "ns/op",
            "extra": "250287 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - B/op",
            "value": 2016,
            "unit": "B/op",
            "extra": "250287 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - allocs/op",
            "value": 54,
            "unit": "allocs/op",
            "extra": "250287 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync)",
            "value": 4501,
            "unit": "ns/op\t    2016 B/op\t      54 allocs/op",
            "extra": "258462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - ns/op",
            "value": 4501,
            "unit": "ns/op",
            "extra": "258462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - B/op",
            "value": 2016,
            "unit": "B/op",
            "extra": "258462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertEvent (Qingyu_backend/service/search/sync) - allocs/op",
            "value": 54,
            "unit": "allocs/op",
            "extra": "258462 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync)",
            "value": 869.9,
            "unit": "ns/op\t     712 B/op\t      11 allocs/op",
            "extra": "1240375 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - ns/op",
            "value": 869.9,
            "unit": "ns/op",
            "extra": "1240375 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - B/op",
            "value": 712,
            "unit": "B/op",
            "extra": "1240375 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "1240375 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync)",
            "value": 1011,
            "unit": "ns/op\t     712 B/op\t      11 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - ns/op",
            "value": 1011,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - B/op",
            "value": 712,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync)",
            "value": 864,
            "unit": "ns/op\t     712 B/op\t      11 allocs/op",
            "extra": "1371033 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - ns/op",
            "value": 864,
            "unit": "ns/op",
            "extra": "1371033 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - B/op",
            "value": 712,
            "unit": "B/op",
            "extra": "1371033 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "1371033 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync)",
            "value": 866.6,
            "unit": "ns/op\t     712 B/op\t      11 allocs/op",
            "extra": "1408758 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - ns/op",
            "value": 866.6,
            "unit": "ns/op",
            "extra": "1408758 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - B/op",
            "value": 712,
            "unit": "B/op",
            "extra": "1408758 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "1408758 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync)",
            "value": 862.6,
            "unit": "ns/op\t     712 B/op\t      11 allocs/op",
            "extra": "1370383 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - ns/op",
            "value": 862.6,
            "unit": "ns/op",
            "extra": "1370383 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - B/op",
            "value": 712,
            "unit": "B/op",
            "extra": "1370383 times\n4 procs"
          },
          {
            "name": "BenchmarkCheckConsistency (Qingyu_backend/service/search/sync) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "1370383 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache)",
            "value": 110437,
            "unit": "ns/op\t   21301 B/op\t     822 allocs/op",
            "extra": "9670 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 110437,
            "unit": "ns/op",
            "extra": "9670 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - B/op",
            "value": 21301,
            "unit": "B/op",
            "extra": "9670 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 822,
            "unit": "allocs/op",
            "extra": "9670 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache)",
            "value": 124649,
            "unit": "ns/op\t   21302 B/op\t     822 allocs/op",
            "extra": "8598 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 124649,
            "unit": "ns/op",
            "extra": "8598 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - B/op",
            "value": 21302,
            "unit": "B/op",
            "extra": "8598 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 822,
            "unit": "allocs/op",
            "extra": "8598 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache)",
            "value": 112125,
            "unit": "ns/op\t   21302 B/op\t     822 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 112125,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - B/op",
            "value": 21302,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 822,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache)",
            "value": 122111,
            "unit": "ns/op\t   21302 B/op\t     822 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 122111,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - B/op",
            "value": 21302,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 822,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache)",
            "value": 112372,
            "unit": "ns/op\t   21302 B/op\t     822 allocs/op",
            "extra": "10923 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 112372,
            "unit": "ns/op",
            "extra": "10923 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - B/op",
            "value": 21302,
            "unit": "B/op",
            "extra": "10923 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MGet (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 822,
            "unit": "allocs/op",
            "extra": "10923 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache)",
            "value": 233291,
            "unit": "ns/op\t    7999 B/op\t     334 allocs/op",
            "extra": "4995 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 233291,
            "unit": "ns/op",
            "extra": "4995 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - B/op",
            "value": 7999,
            "unit": "B/op",
            "extra": "4995 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 334,
            "unit": "allocs/op",
            "extra": "4995 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache)",
            "value": 231716,
            "unit": "ns/op\t    7997 B/op\t     334 allocs/op",
            "extra": "4998 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 231716,
            "unit": "ns/op",
            "extra": "4998 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - B/op",
            "value": 7997,
            "unit": "B/op",
            "extra": "4998 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 334,
            "unit": "allocs/op",
            "extra": "4998 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache)",
            "value": 230710,
            "unit": "ns/op\t    7998 B/op\t     334 allocs/op",
            "extra": "5018 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 230710,
            "unit": "ns/op",
            "extra": "5018 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - B/op",
            "value": 7998,
            "unit": "B/op",
            "extra": "5018 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 334,
            "unit": "allocs/op",
            "extra": "5018 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache)",
            "value": 232276,
            "unit": "ns/op\t    7997 B/op\t     334 allocs/op",
            "extra": "5073 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 232276,
            "unit": "ns/op",
            "extra": "5073 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - B/op",
            "value": 7997,
            "unit": "B/op",
            "extra": "5073 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 334,
            "unit": "allocs/op",
            "extra": "5073 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache)",
            "value": 231829,
            "unit": "ns/op\t    7997 B/op\t     334 allocs/op",
            "extra": "4952 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 231829,
            "unit": "ns/op",
            "extra": "4952 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - B/op",
            "value": 7997,
            "unit": "B/op",
            "extra": "4952 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_MSet (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 334,
            "unit": "allocs/op",
            "extra": "4952 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache)",
            "value": 61764,
            "unit": "ns/op\t     400 B/op\t      20 allocs/op",
            "extra": "20136 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 61764,
            "unit": "ns/op",
            "extra": "20136 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - B/op",
            "value": 400,
            "unit": "B/op",
            "extra": "20136 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 20,
            "unit": "allocs/op",
            "extra": "20136 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache)",
            "value": 60292,
            "unit": "ns/op\t     400 B/op\t      20 allocs/op",
            "extra": "19800 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 60292,
            "unit": "ns/op",
            "extra": "19800 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - B/op",
            "value": 400,
            "unit": "B/op",
            "extra": "19800 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 20,
            "unit": "allocs/op",
            "extra": "19800 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache)",
            "value": 60267,
            "unit": "ns/op\t     400 B/op\t      20 allocs/op",
            "extra": "20145 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 60267,
            "unit": "ns/op",
            "extra": "20145 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - B/op",
            "value": 400,
            "unit": "B/op",
            "extra": "20145 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 20,
            "unit": "allocs/op",
            "extra": "20145 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache)",
            "value": 58888,
            "unit": "ns/op\t     400 B/op\t      20 allocs/op",
            "extra": "20166 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 58888,
            "unit": "ns/op",
            "extra": "20166 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - B/op",
            "value": 400,
            "unit": "B/op",
            "extra": "20166 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 20,
            "unit": "allocs/op",
            "extra": "20166 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache)",
            "value": 59687,
            "unit": "ns/op\t     400 B/op\t      20 allocs/op",
            "extra": "20418 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 59687,
            "unit": "ns/op",
            "extra": "20418 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - B/op",
            "value": 400,
            "unit": "B/op",
            "extra": "20418 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Get (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 20,
            "unit": "allocs/op",
            "extra": "20418 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache)",
            "value": 60246,
            "unit": "ns/op\t     796 B/op\t      35 allocs/op",
            "extra": "19108 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 60246,
            "unit": "ns/op",
            "extra": "19108 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - B/op",
            "value": 796,
            "unit": "B/op",
            "extra": "19108 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "19108 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache)",
            "value": 65308,
            "unit": "ns/op\t     796 B/op\t      35 allocs/op",
            "extra": "18934 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 65308,
            "unit": "ns/op",
            "extra": "18934 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - B/op",
            "value": 796,
            "unit": "B/op",
            "extra": "18934 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "18934 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache)",
            "value": 63838,
            "unit": "ns/op\t     796 B/op\t      35 allocs/op",
            "extra": "19398 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 63838,
            "unit": "ns/op",
            "extra": "19398 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - B/op",
            "value": 796,
            "unit": "B/op",
            "extra": "19398 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "19398 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache)",
            "value": 62739,
            "unit": "ns/op\t     796 B/op\t      35 allocs/op",
            "extra": "19812 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 62739,
            "unit": "ns/op",
            "extra": "19812 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - B/op",
            "value": 796,
            "unit": "B/op",
            "extra": "19812 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "19812 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache)",
            "value": 60319,
            "unit": "ns/op\t     797 B/op\t      35 allocs/op",
            "extra": "18105 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 60319,
            "unit": "ns/op",
            "extra": "18105 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - B/op",
            "value": 797,
            "unit": "B/op",
            "extra": "18105 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_Set (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "18105 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache)",
            "value": 63483,
            "unit": "ns/op\t     876 B/op\t      30 allocs/op",
            "extra": "18860 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 63483,
            "unit": "ns/op",
            "extra": "18860 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - B/op",
            "value": 876,
            "unit": "B/op",
            "extra": "18860 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "18860 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache)",
            "value": 64019,
            "unit": "ns/op\t     876 B/op\t      30 allocs/op",
            "extra": "19165 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 64019,
            "unit": "ns/op",
            "extra": "19165 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - B/op",
            "value": 876,
            "unit": "B/op",
            "extra": "19165 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "19165 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache)",
            "value": 63571,
            "unit": "ns/op\t     877 B/op\t      30 allocs/op",
            "extra": "17688 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 63571,
            "unit": "ns/op",
            "extra": "17688 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - B/op",
            "value": 877,
            "unit": "B/op",
            "extra": "17688 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "17688 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache)",
            "value": 64107,
            "unit": "ns/op\t     876 B/op\t      30 allocs/op",
            "extra": "19059 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 64107,
            "unit": "ns/op",
            "extra": "19059 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - B/op",
            "value": 876,
            "unit": "B/op",
            "extra": "19059 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "19059 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache)",
            "value": 61928,
            "unit": "ns/op\t     876 B/op\t      30 allocs/op",
            "extra": "19402 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 61928,
            "unit": "ns/op",
            "extra": "19402 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - B/op",
            "value": 876,
            "unit": "B/op",
            "extra": "19402 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZAdd (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "19402 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache)",
            "value": 62623,
            "unit": "ns/op\t    1736 B/op\t      59 allocs/op",
            "extra": "18813 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 62623,
            "unit": "ns/op",
            "extra": "18813 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - B/op",
            "value": 1736,
            "unit": "B/op",
            "extra": "18813 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 59,
            "unit": "allocs/op",
            "extra": "18813 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache)",
            "value": 63294,
            "unit": "ns/op\t    1736 B/op\t      59 allocs/op",
            "extra": "18990 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 63294,
            "unit": "ns/op",
            "extra": "18990 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - B/op",
            "value": 1736,
            "unit": "B/op",
            "extra": "18990 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 59,
            "unit": "allocs/op",
            "extra": "18990 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache)",
            "value": 66170,
            "unit": "ns/op\t    1736 B/op\t      59 allocs/op",
            "extra": "20390 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 66170,
            "unit": "ns/op",
            "extra": "20390 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - B/op",
            "value": 1736,
            "unit": "B/op",
            "extra": "20390 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 59,
            "unit": "allocs/op",
            "extra": "20390 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache)",
            "value": 60353,
            "unit": "ns/op\t    1736 B/op\t      59 allocs/op",
            "extra": "19119 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 60353,
            "unit": "ns/op",
            "extra": "19119 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - B/op",
            "value": 1736,
            "unit": "B/op",
            "extra": "19119 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 59,
            "unit": "allocs/op",
            "extra": "19119 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache)",
            "value": 62256,
            "unit": "ns/op\t    1736 B/op\t      59 allocs/op",
            "extra": "18729 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - ns/op",
            "value": 62256,
            "unit": "ns/op",
            "extra": "18729 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - B/op",
            "value": 1736,
            "unit": "B/op",
            "extra": "18729 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisCacheService_ZRange (Qingyu_backend/service/shared/cache) - allocs/op",
            "value": 59,
            "unit": "allocs/op",
            "extra": "18729 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social)",
            "value": 32505,
            "unit": "ns/op\t   10107 B/op\t     113 allocs/op",
            "extra": "38551 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - ns/op",
            "value": 32505,
            "unit": "ns/op",
            "extra": "38551 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - B/op",
            "value": 10107,
            "unit": "B/op",
            "extra": "38551 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - allocs/op",
            "value": 113,
            "unit": "allocs/op",
            "extra": "38551 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social)",
            "value": 31863,
            "unit": "ns/op\t   10099 B/op\t     113 allocs/op",
            "extra": "38720 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - ns/op",
            "value": 31863,
            "unit": "ns/op",
            "extra": "38720 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - B/op",
            "value": 10099,
            "unit": "B/op",
            "extra": "38720 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - allocs/op",
            "value": 113,
            "unit": "allocs/op",
            "extra": "38720 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social)",
            "value": 31068,
            "unit": "ns/op\t   10125 B/op\t     113 allocs/op",
            "extra": "38194 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - ns/op",
            "value": 31068,
            "unit": "ns/op",
            "extra": "38194 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - B/op",
            "value": 10125,
            "unit": "B/op",
            "extra": "38194 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - allocs/op",
            "value": 113,
            "unit": "allocs/op",
            "extra": "38194 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social)",
            "value": 30968,
            "unit": "ns/op\t   10087 B/op\t     113 allocs/op",
            "extra": "39072 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - ns/op",
            "value": 30968,
            "unit": "ns/op",
            "extra": "39072 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - B/op",
            "value": 10087,
            "unit": "B/op",
            "extra": "39072 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - allocs/op",
            "value": 113,
            "unit": "allocs/op",
            "extra": "39072 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social)",
            "value": 33511,
            "unit": "ns/op\t   10197 B/op\t     113 allocs/op",
            "extra": "36385 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - ns/op",
            "value": 33511,
            "unit": "ns/op",
            "extra": "36385 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - B/op",
            "value": 10197,
            "unit": "B/op",
            "extra": "36385 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectionService_AddToCollection (Qingyu_backend/service/social) - allocs/op",
            "value": 113,
            "unit": "allocs/op",
            "extra": "36385 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social)",
            "value": 26256,
            "unit": "ns/op\t    8672 B/op\t      95 allocs/op",
            "extra": "46129 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - ns/op",
            "value": 26256,
            "unit": "ns/op",
            "extra": "46129 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - B/op",
            "value": 8672,
            "unit": "B/op",
            "extra": "46129 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "46129 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social)",
            "value": 25573,
            "unit": "ns/op\t    8646 B/op\t      95 allocs/op",
            "extra": "47341 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - ns/op",
            "value": 25573,
            "unit": "ns/op",
            "extra": "47341 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - B/op",
            "value": 8646,
            "unit": "B/op",
            "extra": "47341 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "47341 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social)",
            "value": 26243,
            "unit": "ns/op\t    8664 B/op\t      95 allocs/op",
            "extra": "46350 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - ns/op",
            "value": 26243,
            "unit": "ns/op",
            "extra": "46350 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - B/op",
            "value": 8664,
            "unit": "B/op",
            "extra": "46350 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "46350 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social)",
            "value": 26208,
            "unit": "ns/op\t    8655 B/op\t      95 allocs/op",
            "extra": "46689 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - ns/op",
            "value": 26208,
            "unit": "ns/op",
            "extra": "46689 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - B/op",
            "value": 8655,
            "unit": "B/op",
            "extra": "46689 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "46689 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social)",
            "value": 26007,
            "unit": "ns/op\t    8686 B/op\t      95 allocs/op",
            "extra": "45907 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - ns/op",
            "value": 26007,
            "unit": "ns/op",
            "extra": "45907 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - B/op",
            "value": 8686,
            "unit": "B/op",
            "extra": "45907 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_LikeBook (Qingyu_backend/service/social) - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "45907 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social)",
            "value": 10959,
            "unit": "ns/op\t    4013 B/op\t      42 allocs/op",
            "extra": "103082 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - ns/op",
            "value": 10959,
            "unit": "ns/op",
            "extra": "103082 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - B/op",
            "value": 4013,
            "unit": "B/op",
            "extra": "103082 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "103082 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social)",
            "value": 11979,
            "unit": "ns/op\t    4020 B/op\t      42 allocs/op",
            "extra": "102476 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - ns/op",
            "value": 11979,
            "unit": "ns/op",
            "extra": "102476 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - B/op",
            "value": 4020,
            "unit": "B/op",
            "extra": "102476 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "102476 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social)",
            "value": 10821,
            "unit": "ns/op\t    4045 B/op\t      42 allocs/op",
            "extra": "100183 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - ns/op",
            "value": 10821,
            "unit": "ns/op",
            "extra": "100183 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - B/op",
            "value": 4045,
            "unit": "B/op",
            "extra": "100183 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "100183 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social)",
            "value": 10898,
            "unit": "ns/op\t    3992 B/op\t      42 allocs/op",
            "extra": "105350 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - ns/op",
            "value": 10898,
            "unit": "ns/op",
            "extra": "105350 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - B/op",
            "value": 3992,
            "unit": "B/op",
            "extra": "105350 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "105350 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social)",
            "value": 11008,
            "unit": "ns/op\t    4009 B/op\t      42 allocs/op",
            "extra": "103592 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - ns/op",
            "value": 11008,
            "unit": "ns/op",
            "extra": "103592 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - B/op",
            "value": 4009,
            "unit": "B/op",
            "extra": "103592 times\n4 procs"
          },
          {
            "name": "BenchmarkLikeService_GetBookLikeCount (Qingyu_backend/service/social) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "103592 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social)",
            "value": 14027,
            "unit": "ns/op\t    4445 B/op\t      50 allocs/op",
            "extra": "83564 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - ns/op",
            "value": 14027,
            "unit": "ns/op",
            "extra": "83564 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - B/op",
            "value": 4445,
            "unit": "B/op",
            "extra": "83564 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - allocs/op",
            "value": 50,
            "unit": "allocs/op",
            "extra": "83564 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social)",
            "value": 13923,
            "unit": "ns/op\t    4435 B/op\t      50 allocs/op",
            "extra": "84300 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - ns/op",
            "value": 13923,
            "unit": "ns/op",
            "extra": "84300 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - B/op",
            "value": 4435,
            "unit": "B/op",
            "extra": "84300 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - allocs/op",
            "value": 50,
            "unit": "allocs/op",
            "extra": "84300 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social)",
            "value": 13656,
            "unit": "ns/op\t    4438 B/op\t      50 allocs/op",
            "extra": "83989 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - ns/op",
            "value": 13656,
            "unit": "ns/op",
            "extra": "83989 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - B/op",
            "value": 4438,
            "unit": "B/op",
            "extra": "83989 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - allocs/op",
            "value": 50,
            "unit": "allocs/op",
            "extra": "83989 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social)",
            "value": 13785,
            "unit": "ns/op\t    4438 B/op\t      50 allocs/op",
            "extra": "84254 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - ns/op",
            "value": 13785,
            "unit": "ns/op",
            "extra": "84254 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - B/op",
            "value": 4438,
            "unit": "B/op",
            "extra": "84254 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - allocs/op",
            "value": 50,
            "unit": "allocs/op",
            "extra": "84254 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social)",
            "value": 15044,
            "unit": "ns/op\t    4464 B/op\t      50 allocs/op",
            "extra": "82126 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - ns/op",
            "value": 15044,
            "unit": "ns/op",
            "extra": "82126 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "82126 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheHit (Qingyu_backend/service/social) - allocs/op",
            "value": 50,
            "unit": "allocs/op",
            "extra": "82126 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social)",
            "value": 40320,
            "unit": "ns/op\t   13224 B/op\t     138 allocs/op",
            "extra": "29946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - ns/op",
            "value": 40320,
            "unit": "ns/op",
            "extra": "29946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - B/op",
            "value": 13224,
            "unit": "B/op",
            "extra": "29946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - allocs/op",
            "value": 138,
            "unit": "allocs/op",
            "extra": "29946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social)",
            "value": 40113,
            "unit": "ns/op\t   13296 B/op\t     138 allocs/op",
            "extra": "29277 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - ns/op",
            "value": 40113,
            "unit": "ns/op",
            "extra": "29277 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - B/op",
            "value": 13296,
            "unit": "B/op",
            "extra": "29277 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - allocs/op",
            "value": 138,
            "unit": "allocs/op",
            "extra": "29277 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social)",
            "value": 41069,
            "unit": "ns/op\t   13216 B/op\t     138 allocs/op",
            "extra": "30048 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - ns/op",
            "value": 41069,
            "unit": "ns/op",
            "extra": "30048 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - B/op",
            "value": 13216,
            "unit": "B/op",
            "extra": "30048 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - allocs/op",
            "value": 138,
            "unit": "allocs/op",
            "extra": "30048 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social)",
            "value": 39982,
            "unit": "ns/op\t   13218 B/op\t     138 allocs/op",
            "extra": "29976 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - ns/op",
            "value": 39982,
            "unit": "ns/op",
            "extra": "29976 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - B/op",
            "value": 13218,
            "unit": "B/op",
            "extra": "29976 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - allocs/op",
            "value": 138,
            "unit": "allocs/op",
            "extra": "29976 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social)",
            "value": 40572,
            "unit": "ns/op\t   13318 B/op\t     138 allocs/op",
            "extra": "28975 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - ns/op",
            "value": 40572,
            "unit": "ns/op",
            "extra": "28975 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - B/op",
            "value": 13318,
            "unit": "B/op",
            "extra": "28975 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_CacheMiss (Qingyu_backend/service/social) - allocs/op",
            "value": 138,
            "unit": "allocs/op",
            "extra": "28975 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social)",
            "value": 43403,
            "unit": "ns/op\t   14096 B/op\t     161 allocs/op",
            "extra": "28029 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - ns/op",
            "value": 43403,
            "unit": "ns/op",
            "extra": "28029 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - B/op",
            "value": 14096,
            "unit": "B/op",
            "extra": "28029 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - allocs/op",
            "value": 161,
            "unit": "allocs/op",
            "extra": "28029 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social)",
            "value": 43190,
            "unit": "ns/op\t   14161 B/op\t     161 allocs/op",
            "extra": "27573 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - ns/op",
            "value": 43190,
            "unit": "ns/op",
            "extra": "27573 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - B/op",
            "value": 14161,
            "unit": "B/op",
            "extra": "27573 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - allocs/op",
            "value": 161,
            "unit": "allocs/op",
            "extra": "27573 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social)",
            "value": 43447,
            "unit": "ns/op\t   14121 B/op\t     161 allocs/op",
            "extra": "27920 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - ns/op",
            "value": 43447,
            "unit": "ns/op",
            "extra": "27920 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - B/op",
            "value": 14121,
            "unit": "B/op",
            "extra": "27920 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - allocs/op",
            "value": 161,
            "unit": "allocs/op",
            "extra": "27920 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social)",
            "value": 43331,
            "unit": "ns/op\t   14166 B/op\t     161 allocs/op",
            "extra": "27433 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - ns/op",
            "value": 43331,
            "unit": "ns/op",
            "extra": "27433 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - B/op",
            "value": 14166,
            "unit": "B/op",
            "extra": "27433 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - allocs/op",
            "value": 161,
            "unit": "allocs/op",
            "extra": "27433 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social)",
            "value": 43350,
            "unit": "ns/op\t   14119 B/op\t     161 allocs/op",
            "extra": "27964 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - ns/op",
            "value": 43350,
            "unit": "ns/op",
            "extra": "27964 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - B/op",
            "value": 14119,
            "unit": "B/op",
            "extra": "27964 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Book_CacheMiss (Qingyu_backend/service/social) - allocs/op",
            "value": 161,
            "unit": "allocs/op",
            "extra": "27964 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social)",
            "value": 18751,
            "unit": "ns/op\t    4864 B/op\t      59 allocs/op",
            "extra": "62271 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - ns/op",
            "value": 18751,
            "unit": "ns/op",
            "extra": "62271 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - B/op",
            "value": 4864,
            "unit": "B/op",
            "extra": "62271 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - allocs/op",
            "value": 59,
            "unit": "allocs/op",
            "extra": "62271 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social)",
            "value": 18870,
            "unit": "ns/op\t    4862 B/op\t      59 allocs/op",
            "extra": "62409 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - ns/op",
            "value": 18870,
            "unit": "ns/op",
            "extra": "62409 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - B/op",
            "value": 4862,
            "unit": "B/op",
            "extra": "62409 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - allocs/op",
            "value": 59,
            "unit": "allocs/op",
            "extra": "62409 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social)",
            "value": 20592,
            "unit": "ns/op\t    5006 B/op\t      59 allocs/op",
            "extra": "53323 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - ns/op",
            "value": 20592,
            "unit": "ns/op",
            "extra": "53323 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - B/op",
            "value": 5006,
            "unit": "B/op",
            "extra": "53323 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - allocs/op",
            "value": 59,
            "unit": "allocs/op",
            "extra": "53323 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social)",
            "value": 18905,
            "unit": "ns/op\t    4867 B/op\t      59 allocs/op",
            "extra": "62058 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - ns/op",
            "value": 18905,
            "unit": "ns/op",
            "extra": "62058 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - B/op",
            "value": 4867,
            "unit": "B/op",
            "extra": "62058 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - allocs/op",
            "value": 59,
            "unit": "allocs/op",
            "extra": "62058 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social)",
            "value": 19082,
            "unit": "ns/op\t    4880 B/op\t      59 allocs/op",
            "extra": "61165 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - ns/op",
            "value": 19082,
            "unit": "ns/op",
            "extra": "61165 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - B/op",
            "value": 4880,
            "unit": "B/op",
            "extra": "61165 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRatingStats_Concurrent_CacheHit (Qingyu_backend/service/social) - allocs/op",
            "value": 59,
            "unit": "allocs/op",
            "extra": "61165 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social)",
            "value": 1737,
            "unit": "ns/op\t     720 B/op\t      14 allocs/op",
            "extra": "645255 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - ns/op",
            "value": 1737,
            "unit": "ns/op",
            "extra": "645255 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - B/op",
            "value": 720,
            "unit": "B/op",
            "extra": "645255 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "645255 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social)",
            "value": 1765,
            "unit": "ns/op\t     720 B/op\t      14 allocs/op",
            "extra": "645678 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - ns/op",
            "value": 1765,
            "unit": "ns/op",
            "extra": "645678 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - B/op",
            "value": 720,
            "unit": "B/op",
            "extra": "645678 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "645678 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social)",
            "value": 1734,
            "unit": "ns/op\t     720 B/op\t      14 allocs/op",
            "extra": "632529 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - ns/op",
            "value": 1734,
            "unit": "ns/op",
            "extra": "632529 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - B/op",
            "value": 720,
            "unit": "B/op",
            "extra": "632529 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "632529 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social)",
            "value": 1732,
            "unit": "ns/op\t     720 B/op\t      14 allocs/op",
            "extra": "641550 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - ns/op",
            "value": 1732,
            "unit": "ns/op",
            "extra": "641550 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - B/op",
            "value": 720,
            "unit": "B/op",
            "extra": "641550 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "641550 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social)",
            "value": 1736,
            "unit": "ns/op\t     720 B/op\t      14 allocs/op",
            "extra": "638721 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - ns/op",
            "value": 1736,
            "unit": "ns/op",
            "extra": "638721 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - B/op",
            "value": 720,
            "unit": "B/op",
            "extra": "638721 times\n4 procs"
          },
          {
            "name": "BenchmarkSerializeStats (Qingyu_backend/service/social) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "638721 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social)",
            "value": 4245,
            "unit": "ns/op\t     768 B/op\t      17 allocs/op",
            "extra": "305378 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - ns/op",
            "value": 4245,
            "unit": "ns/op",
            "extra": "305378 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - B/op",
            "value": 768,
            "unit": "B/op",
            "extra": "305378 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "305378 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social)",
            "value": 3793,
            "unit": "ns/op\t     768 B/op\t      17 allocs/op",
            "extra": "304836 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - ns/op",
            "value": 3793,
            "unit": "ns/op",
            "extra": "304836 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - B/op",
            "value": 768,
            "unit": "B/op",
            "extra": "304836 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "304836 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social)",
            "value": 3790,
            "unit": "ns/op\t     768 B/op\t      17 allocs/op",
            "extra": "304398 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - ns/op",
            "value": 3790,
            "unit": "ns/op",
            "extra": "304398 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - B/op",
            "value": 768,
            "unit": "B/op",
            "extra": "304398 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "304398 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social)",
            "value": 3869,
            "unit": "ns/op\t     768 B/op\t      17 allocs/op",
            "extra": "305229 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - ns/op",
            "value": 3869,
            "unit": "ns/op",
            "extra": "305229 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - B/op",
            "value": 768,
            "unit": "B/op",
            "extra": "305229 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "305229 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social)",
            "value": 3811,
            "unit": "ns/op\t     768 B/op\t      17 allocs/op",
            "extra": "301538 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - ns/op",
            "value": 3811,
            "unit": "ns/op",
            "extra": "301538 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - B/op",
            "value": 768,
            "unit": "B/op",
            "extra": "301538 times\n4 procs"
          },
          {
            "name": "BenchmarkDeserializeStats (Qingyu_backend/service/social) - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "301538 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social)",
            "value": 10701,
            "unit": "ns/op\t    3822 B/op\t      42 allocs/op",
            "extra": "108097 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - ns/op",
            "value": 10701,
            "unit": "ns/op",
            "extra": "108097 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - B/op",
            "value": 3822,
            "unit": "B/op",
            "extra": "108097 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "108097 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social)",
            "value": 11064,
            "unit": "ns/op\t    3863 B/op\t      42 allocs/op",
            "extra": "103548 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - ns/op",
            "value": 11064,
            "unit": "ns/op",
            "extra": "103548 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - B/op",
            "value": 3863,
            "unit": "B/op",
            "extra": "103548 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "103548 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social)",
            "value": 10714,
            "unit": "ns/op\t    3830 B/op\t      42 allocs/op",
            "extra": "107067 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - ns/op",
            "value": 10714,
            "unit": "ns/op",
            "extra": "107067 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - B/op",
            "value": 3830,
            "unit": "B/op",
            "extra": "107067 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "107067 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social)",
            "value": 10935,
            "unit": "ns/op\t    3840 B/op\t      42 allocs/op",
            "extra": "106117 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - ns/op",
            "value": 10935,
            "unit": "ns/op",
            "extra": "106117 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - B/op",
            "value": 3840,
            "unit": "B/op",
            "extra": "106117 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "106117 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social)",
            "value": 11865,
            "unit": "ns/op\t    3796 B/op\t      42 allocs/op",
            "extra": "87948 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - ns/op",
            "value": 11865,
            "unit": "ns/op",
            "extra": "87948 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - B/op",
            "value": 3796,
            "unit": "B/op",
            "extra": "87948 times\n4 procs"
          },
          {
            "name": "BenchmarkInvalidateCache (Qingyu_backend/service/social) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "87948 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user)",
            "value": 145.4,
            "unit": "ns/op\t     112 B/op\t       3 allocs/op",
            "extra": "8230773 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - ns/op",
            "value": 145.4,
            "unit": "ns/op",
            "extra": "8230773 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "8230773 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "8230773 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user)",
            "value": 154.5,
            "unit": "ns/op\t     112 B/op\t       3 allocs/op",
            "extra": "8102551 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - ns/op",
            "value": 154.5,
            "unit": "ns/op",
            "extra": "8102551 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "8102551 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "8102551 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user)",
            "value": 145.5,
            "unit": "ns/op\t     112 B/op\t       3 allocs/op",
            "extra": "7317631 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - ns/op",
            "value": 145.5,
            "unit": "ns/op",
            "extra": "7317631 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "7317631 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "7317631 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user)",
            "value": 146.1,
            "unit": "ns/op\t     112 B/op\t       3 allocs/op",
            "extra": "8166284 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - ns/op",
            "value": 146.1,
            "unit": "ns/op",
            "extra": "8166284 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "8166284 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "8166284 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user)",
            "value": 146.5,
            "unit": "ns/op\t     112 B/op\t       3 allocs/op",
            "extra": "8182762 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - ns/op",
            "value": 146.5,
            "unit": "ns/op",
            "extra": "8182762 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "8182762 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorCreation (Qingyu_backend/service/user) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "8182762 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user)",
            "value": 0.3137,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - ns/op",
            "value": 0.3137,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user)",
            "value": 0.3123,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - ns/op",
            "value": 0.3123,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user)",
            "value": 0.3115,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - ns/op",
            "value": 0.3115,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user)",
            "value": 0.3119,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - ns/op",
            "value": 0.3119,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user)",
            "value": 0.3265,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - ns/op",
            "value": 0.3265,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkUserErrorWithCause (Qingyu_backend/service/user) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user)",
            "value": 2638,
            "unit": "ns/op\t    2337 B/op\t      28 allocs/op",
            "extra": "407002 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - ns/op",
            "value": 2638,
            "unit": "ns/op",
            "extra": "407002 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - B/op",
            "value": 2337,
            "unit": "B/op",
            "extra": "407002 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "407002 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user)",
            "value": 2640,
            "unit": "ns/op\t    2337 B/op\t      28 allocs/op",
            "extra": "438373 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - ns/op",
            "value": 2640,
            "unit": "ns/op",
            "extra": "438373 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - B/op",
            "value": 2337,
            "unit": "B/op",
            "extra": "438373 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "438373 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user)",
            "value": 2622,
            "unit": "ns/op\t    2337 B/op\t      28 allocs/op",
            "extra": "438464 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - ns/op",
            "value": 2622,
            "unit": "ns/op",
            "extra": "438464 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - B/op",
            "value": 2337,
            "unit": "B/op",
            "extra": "438464 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "438464 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user)",
            "value": 2618,
            "unit": "ns/op\t    2336 B/op\t      28 allocs/op",
            "extra": "444103 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - ns/op",
            "value": 2618,
            "unit": "ns/op",
            "extra": "444103 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - B/op",
            "value": 2336,
            "unit": "B/op",
            "extra": "444103 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "444103 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user)",
            "value": 2621,
            "unit": "ns/op\t    2338 B/op\t      28 allocs/op",
            "extra": "447364 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - ns/op",
            "value": 2621,
            "unit": "ns/op",
            "extra": "447364 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - B/op",
            "value": 2338,
            "unit": "B/op",
            "extra": "447364 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_ValidateStrength (Qingyu_backend/service/user) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "447364 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user)",
            "value": 4495,
            "unit": "ns/op\t    3354 B/op\t      41 allocs/op",
            "extra": "261315 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - ns/op",
            "value": 4495,
            "unit": "ns/op",
            "extra": "261315 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - B/op",
            "value": 3354,
            "unit": "B/op",
            "extra": "261315 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - allocs/op",
            "value": 41,
            "unit": "allocs/op",
            "extra": "261315 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user)",
            "value": 4508,
            "unit": "ns/op\t    3351 B/op\t      41 allocs/op",
            "extra": "262989 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - ns/op",
            "value": 4508,
            "unit": "ns/op",
            "extra": "262989 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - B/op",
            "value": 3351,
            "unit": "B/op",
            "extra": "262989 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - allocs/op",
            "value": 41,
            "unit": "allocs/op",
            "extra": "262989 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user)",
            "value": 4755,
            "unit": "ns/op\t    3352 B/op\t      41 allocs/op",
            "extra": "263842 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - ns/op",
            "value": 4755,
            "unit": "ns/op",
            "extra": "263842 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - B/op",
            "value": 3352,
            "unit": "B/op",
            "extra": "263842 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - allocs/op",
            "value": 41,
            "unit": "allocs/op",
            "extra": "263842 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user)",
            "value": 4568,
            "unit": "ns/op\t    3353 B/op\t      41 allocs/op",
            "extra": "257625 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - ns/op",
            "value": 4568,
            "unit": "ns/op",
            "extra": "257625 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - B/op",
            "value": 3353,
            "unit": "B/op",
            "extra": "257625 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - allocs/op",
            "value": 41,
            "unit": "allocs/op",
            "extra": "257625 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user)",
            "value": 4517,
            "unit": "ns/op\t    3351 B/op\t      41 allocs/op",
            "extra": "262197 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - ns/op",
            "value": 4517,
            "unit": "ns/op",
            "extra": "262197 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - B/op",
            "value": 3351,
            "unit": "B/op",
            "extra": "262197 times\n4 procs"
          },
          {
            "name": "BenchmarkPasswordValidator_GetStrengthScore (Qingyu_backend/service/user) - allocs/op",
            "value": 41,
            "unit": "allocs/op",
            "extra": "262197 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document)",
            "value": 1026,
            "unit": "ns/op\t    1795 B/op\t      12 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - ns/op",
            "value": 1026,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - B/op",
            "value": 1795,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document)",
            "value": 1039,
            "unit": "ns/op\t    1795 B/op\t      12 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - ns/op",
            "value": 1039,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - B/op",
            "value": 1795,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document)",
            "value": 1037,
            "unit": "ns/op\t    1795 B/op\t      12 allocs/op",
            "extra": "1164044 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - ns/op",
            "value": 1037,
            "unit": "ns/op",
            "extra": "1164044 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - B/op",
            "value": 1795,
            "unit": "B/op",
            "extra": "1164044 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "1164044 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document)",
            "value": 1037,
            "unit": "ns/op\t    1795 B/op\t      12 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - ns/op",
            "value": 1037,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - B/op",
            "value": 1795,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document)",
            "value": 1041,
            "unit": "ns/op\t    1795 B/op\t      12 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - ns/op",
            "value": 1041,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - B/op",
            "value": 1795,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCreateDuplicateDocument (Qingyu_backend/service/writer/document) - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline)",
            "value": 2.296,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "548124333 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - ns/op",
            "value": 2.296,
            "unit": "ns/op",
            "extra": "548124333 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "548124333 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "548124333 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline)",
            "value": 2.195,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "547471911 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - ns/op",
            "value": 2.195,
            "unit": "ns/op",
            "extra": "547471911 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "547471911 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "547471911 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline)",
            "value": 2.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "550683588 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - ns/op",
            "value": 2.2,
            "unit": "ns/op",
            "extra": "550683588 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "550683588 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "550683588 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline)",
            "value": 2.194,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "546698116 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - ns/op",
            "value": 2.194,
            "unit": "ns/op",
            "extra": "546698116 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "546698116 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "546698116 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline)",
            "value": 2.184,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "548250936 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - ns/op",
            "value": 2.184,
            "unit": "ns/op",
            "extra": "548250936 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "548250936 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/ValidateToken (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "548250936 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline)",
            "value": 3.435,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "348530104 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - ns/op",
            "value": 3.435,
            "unit": "ns/op",
            "extra": "348530104 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "348530104 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "348530104 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline)",
            "value": 3.437,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "349530252 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - ns/op",
            "value": 3.437,
            "unit": "ns/op",
            "extra": "349530252 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "349530252 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "349530252 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline)",
            "value": 3.659,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "348786868 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - ns/op",
            "value": 3.659,
            "unit": "ns/op",
            "extra": "348786868 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "348786868 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "348786868 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline)",
            "value": 3.432,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "342856590 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - ns/op",
            "value": 3.432,
            "unit": "ns/op",
            "extra": "342856590 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "342856590 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "342856590 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline)",
            "value": 3.448,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "344845482 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - ns/op",
            "value": 3.448,
            "unit": "ns/op",
            "extra": "344845482 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "344845482 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenOperations/CheckPermission (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "344845482 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline)",
            "value": 4.115,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "295252326 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 4.115,
            "unit": "ns/op",
            "extra": "295252326 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "295252326 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "295252326 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline)",
            "value": 4.066,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "296040940 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 4.066,
            "unit": "ns/op",
            "extra": "296040940 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "296040940 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "296040940 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline)",
            "value": 4.074,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "294363769 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 4.074,
            "unit": "ns/op",
            "extra": "294363769 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "294363769 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "294363769 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline)",
            "value": 4.151,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "295710974 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 4.151,
            "unit": "ns/op",
            "extra": "295710974 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "295710974 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "295710974 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline)",
            "value": 4.081,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "266095605 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 4.081,
            "unit": "ns/op",
            "extra": "266095605 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "266095605 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadSmallFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "266095605 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline)",
            "value": 3.754,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "320124598 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 3.754,
            "unit": "ns/op",
            "extra": "320124598 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "320124598 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "320124598 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline)",
            "value": 3.759,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "320176807 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 3.759,
            "unit": "ns/op",
            "extra": "320176807 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "320176807 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "320176807 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline)",
            "value": 3.741,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "319164530 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 3.741,
            "unit": "ns/op",
            "extra": "319164530 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "319164530 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "319164530 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline)",
            "value": 3.744,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "320049316 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 3.744,
            "unit": "ns/op",
            "extra": "320049316 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "320049316 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "320049316 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline)",
            "value": 3.748,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "319235109 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 3.748,
            "unit": "ns/op",
            "extra": "319235109 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "319235109 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/UploadLargeFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "319235109 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline)",
            "value": 15.2,
            "unit": "ns/op\t       4 B/op\t       1 allocs/op",
            "extra": "84000036 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 15.2,
            "unit": "ns/op",
            "extra": "84000036 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - B/op",
            "value": 4,
            "unit": "B/op",
            "extra": "84000036 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "84000036 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline)",
            "value": 14.03,
            "unit": "ns/op\t       4 B/op\t       1 allocs/op",
            "extra": "85094349 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 14.03,
            "unit": "ns/op",
            "extra": "85094349 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - B/op",
            "value": 4,
            "unit": "B/op",
            "extra": "85094349 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "85094349 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline)",
            "value": 13.99,
            "unit": "ns/op\t       4 B/op\t       1 allocs/op",
            "extra": "74250633 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 13.99,
            "unit": "ns/op",
            "extra": "74250633 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - B/op",
            "value": 4,
            "unit": "B/op",
            "extra": "74250633 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "74250633 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline)",
            "value": 13.96,
            "unit": "ns/op\t       4 B/op\t       1 allocs/op",
            "extra": "75948549 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 13.96,
            "unit": "ns/op",
            "extra": "75948549 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - B/op",
            "value": 4,
            "unit": "B/op",
            "extra": "75948549 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "75948549 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline)",
            "value": 13.99,
            "unit": "ns/op\t       4 B/op\t       1 allocs/op",
            "extra": "73387555 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - ns/op",
            "value": 13.99,
            "unit": "ns/op",
            "extra": "73387555 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - B/op",
            "value": 4,
            "unit": "B/op",
            "extra": "73387555 times\n4 procs"
          },
          {
            "name": "BenchmarkStorageOperations/DownloadFile (Qingyu_backend/test/baseline) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "73387555 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration)",
            "value": 8348,
            "unit": "ns/op\t    3227 B/op\t      32 allocs/op",
            "extra": "151801 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - ns/op",
            "value": 8348,
            "unit": "ns/op",
            "extra": "151801 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - B/op",
            "value": 3227,
            "unit": "B/op",
            "extra": "151801 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - allocs/op",
            "value": 32,
            "unit": "allocs/op",
            "extra": "151801 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration)",
            "value": 8108,
            "unit": "ns/op\t    3219 B/op\t      32 allocs/op",
            "extra": "152908 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - ns/op",
            "value": 8108,
            "unit": "ns/op",
            "extra": "152908 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - B/op",
            "value": 3219,
            "unit": "B/op",
            "extra": "152908 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - allocs/op",
            "value": 32,
            "unit": "allocs/op",
            "extra": "152908 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration)",
            "value": 8131,
            "unit": "ns/op\t    3237 B/op\t      32 allocs/op",
            "extra": "150048 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - ns/op",
            "value": 8131,
            "unit": "ns/op",
            "extra": "150048 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - B/op",
            "value": 3237,
            "unit": "B/op",
            "extra": "150048 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - allocs/op",
            "value": 32,
            "unit": "allocs/op",
            "extra": "150048 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration)",
            "value": 9285,
            "unit": "ns/op\t    3226 B/op\t      32 allocs/op",
            "extra": "151707 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - ns/op",
            "value": 9285,
            "unit": "ns/op",
            "extra": "151707 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - B/op",
            "value": 3226,
            "unit": "B/op",
            "extra": "151707 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - allocs/op",
            "value": 32,
            "unit": "allocs/op",
            "extra": "151707 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration)",
            "value": 8310,
            "unit": "ns/op\t    3220 B/op\t      32 allocs/op",
            "extra": "152809 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - ns/op",
            "value": 8310,
            "unit": "ns/op",
            "extra": "152809 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - B/op",
            "value": 3220,
            "unit": "B/op",
            "extra": "152809 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Login (Qingyu_backend/test/integration) - allocs/op",
            "value": 32,
            "unit": "allocs/op",
            "extra": "152809 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration)",
            "value": 8186,
            "unit": "ns/op\t    3130 B/op\t      30 allocs/op",
            "extra": "151706 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - ns/op",
            "value": 8186,
            "unit": "ns/op",
            "extra": "151706 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - B/op",
            "value": 3130,
            "unit": "B/op",
            "extra": "151706 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "151706 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration)",
            "value": 8054,
            "unit": "ns/op\t    3129 B/op\t      30 allocs/op",
            "extra": "151820 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - ns/op",
            "value": 8054,
            "unit": "ns/op",
            "extra": "151820 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - B/op",
            "value": 3129,
            "unit": "B/op",
            "extra": "151820 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "151820 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration)",
            "value": 8100,
            "unit": "ns/op\t    3132 B/op\t      30 allocs/op",
            "extra": "151561 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - ns/op",
            "value": 8100,
            "unit": "ns/op",
            "extra": "151561 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - B/op",
            "value": 3132,
            "unit": "B/op",
            "extra": "151561 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "151561 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration)",
            "value": 8070,
            "unit": "ns/op\t    3199 B/op\t      30 allocs/op",
            "extra": "140713 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - ns/op",
            "value": 8070,
            "unit": "ns/op",
            "extra": "140713 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - B/op",
            "value": 3199,
            "unit": "B/op",
            "extra": "140713 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "140713 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration)",
            "value": 8051,
            "unit": "ns/op\t    3117 B/op\t      30 allocs/op",
            "extra": "153950 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - ns/op",
            "value": 8051,
            "unit": "ns/op",
            "extra": "153950 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - B/op",
            "value": 3117,
            "unit": "B/op",
            "extra": "153950 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_ValidateToken (Qingyu_backend/test/integration) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "153950 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration)",
            "value": 8989,
            "unit": "ns/op\t    3565 B/op\t      39 allocs/op",
            "extra": "140583 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - ns/op",
            "value": 8989,
            "unit": "ns/op",
            "extra": "140583 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - B/op",
            "value": 3565,
            "unit": "B/op",
            "extra": "140583 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - allocs/op",
            "value": 39,
            "unit": "allocs/op",
            "extra": "140583 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration)",
            "value": 9859,
            "unit": "ns/op\t    3617 B/op\t      39 allocs/op",
            "extra": "133147 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - ns/op",
            "value": 9859,
            "unit": "ns/op",
            "extra": "133147 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - B/op",
            "value": 3617,
            "unit": "B/op",
            "extra": "133147 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - allocs/op",
            "value": 39,
            "unit": "allocs/op",
            "extra": "133147 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration)",
            "value": 9006,
            "unit": "ns/op\t    3573 B/op\t      39 allocs/op",
            "extra": "139226 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - ns/op",
            "value": 9006,
            "unit": "ns/op",
            "extra": "139226 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - B/op",
            "value": 3573,
            "unit": "B/op",
            "extra": "139226 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - allocs/op",
            "value": 39,
            "unit": "allocs/op",
            "extra": "139226 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration)",
            "value": 8993,
            "unit": "ns/op\t    3553 B/op\t      39 allocs/op",
            "extra": "142240 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - ns/op",
            "value": 8993,
            "unit": "ns/op",
            "extra": "142240 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - B/op",
            "value": 3553,
            "unit": "B/op",
            "extra": "142240 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - allocs/op",
            "value": 39,
            "unit": "allocs/op",
            "extra": "142240 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration)",
            "value": 8882,
            "unit": "ns/op\t    3565 B/op\t      39 allocs/op",
            "extra": "140420 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - ns/op",
            "value": 8882,
            "unit": "ns/op",
            "extra": "140420 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - B/op",
            "value": 3565,
            "unit": "B/op",
            "extra": "140420 times\n4 procs"
          },
          {
            "name": "BenchmarkAuthService_Register (Qingyu_backend/test/integration) - allocs/op",
            "value": 39,
            "unit": "allocs/op",
            "extra": "140420 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration)",
            "value": 7954,
            "unit": "ns/op\t    3108 B/op\t      30 allocs/op",
            "extra": "151419 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - ns/op",
            "value": 7954,
            "unit": "ns/op",
            "extra": "151419 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - B/op",
            "value": 3108,
            "unit": "B/op",
            "extra": "151419 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "151419 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration)",
            "value": 8081,
            "unit": "ns/op\t    3297 B/op\t      30 allocs/op",
            "extra": "156642 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - ns/op",
            "value": 8081,
            "unit": "ns/op",
            "extra": "156642 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - B/op",
            "value": 3297,
            "unit": "B/op",
            "extra": "156642 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "156642 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration)",
            "value": 8018,
            "unit": "ns/op\t    3087 B/op\t      30 allocs/op",
            "extra": "155162 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - ns/op",
            "value": 8018,
            "unit": "ns/op",
            "extra": "155162 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - B/op",
            "value": 3087,
            "unit": "B/op",
            "extra": "155162 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "155162 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration)",
            "value": 8005,
            "unit": "ns/op\t    3093 B/op\t      30 allocs/op",
            "extra": "154203 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - ns/op",
            "value": 8005,
            "unit": "ns/op",
            "extra": "154203 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - B/op",
            "value": 3093,
            "unit": "B/op",
            "extra": "154203 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "154203 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration)",
            "value": 8654,
            "unit": "ns/op\t    3141 B/op\t      30 allocs/op",
            "extra": "145964 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - ns/op",
            "value": 8654,
            "unit": "ns/op",
            "extra": "145964 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - B/op",
            "value": 3141,
            "unit": "B/op",
            "extra": "145964 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_GetBalance (Qingyu_backend/test/integration) - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "145964 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration)",
            "value": 9865,
            "unit": "ns/op\t    3976 B/op\t      44 allocs/op",
            "extra": "126741 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - ns/op",
            "value": 9865,
            "unit": "ns/op",
            "extra": "126741 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - B/op",
            "value": 3976,
            "unit": "B/op",
            "extra": "126741 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "126741 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration)",
            "value": 9827,
            "unit": "ns/op\t    3969 B/op\t      44 allocs/op",
            "extra": "127860 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - ns/op",
            "value": 9827,
            "unit": "ns/op",
            "extra": "127860 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - B/op",
            "value": 3969,
            "unit": "B/op",
            "extra": "127860 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "127860 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration)",
            "value": 9770,
            "unit": "ns/op\t    3963 B/op\t      44 allocs/op",
            "extra": "128264 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - ns/op",
            "value": 9770,
            "unit": "ns/op",
            "extra": "128264 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - B/op",
            "value": 3963,
            "unit": "B/op",
            "extra": "128264 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "128264 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration)",
            "value": 9752,
            "unit": "ns/op\t    3963 B/op\t      44 allocs/op",
            "extra": "128356 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - ns/op",
            "value": 9752,
            "unit": "ns/op",
            "extra": "128356 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - B/op",
            "value": 3963,
            "unit": "B/op",
            "extra": "128356 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "128356 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration)",
            "value": 9716,
            "unit": "ns/op\t    3960 B/op\t      44 allocs/op",
            "extra": "128577 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - ns/op",
            "value": 9716,
            "unit": "ns/op",
            "extra": "128577 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - B/op",
            "value": 3960,
            "unit": "B/op",
            "extra": "128577 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Recharge (Qingyu_backend/test/integration) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "128577 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration)",
            "value": 9812,
            "unit": "ns/op\t    3953 B/op\t      44 allocs/op",
            "extra": "128448 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - ns/op",
            "value": 9812,
            "unit": "ns/op",
            "extra": "128448 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - B/op",
            "value": 3953,
            "unit": "B/op",
            "extra": "128448 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "128448 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration)",
            "value": 10284,
            "unit": "ns/op\t    3971 B/op\t      44 allocs/op",
            "extra": "126513 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - ns/op",
            "value": 10284,
            "unit": "ns/op",
            "extra": "126513 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - B/op",
            "value": 3971,
            "unit": "B/op",
            "extra": "126513 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "126513 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration)",
            "value": 9687,
            "unit": "ns/op\t    3942 B/op\t      44 allocs/op",
            "extra": "103418 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - ns/op",
            "value": 9687,
            "unit": "ns/op",
            "extra": "103418 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - B/op",
            "value": 3942,
            "unit": "B/op",
            "extra": "103418 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "103418 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration)",
            "value": 9890,
            "unit": "ns/op\t    3977 B/op\t      44 allocs/op",
            "extra": "125666 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - ns/op",
            "value": 9890,
            "unit": "ns/op",
            "extra": "125666 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - B/op",
            "value": 3977,
            "unit": "B/op",
            "extra": "125666 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "125666 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration)",
            "value": 9726,
            "unit": "ns/op\t    3963 B/op\t      44 allocs/op",
            "extra": "127570 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - ns/op",
            "value": 9726,
            "unit": "ns/op",
            "extra": "127570 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - B/op",
            "value": 3963,
            "unit": "B/op",
            "extra": "127570 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Consume (Qingyu_backend/test/integration) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "127570 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration)",
            "value": 10467,
            "unit": "ns/op\t    4218 B/op\t      51 allocs/op",
            "extra": "118765 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - ns/op",
            "value": 10467,
            "unit": "ns/op",
            "extra": "118765 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - B/op",
            "value": 4218,
            "unit": "B/op",
            "extra": "118765 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - allocs/op",
            "value": 51,
            "unit": "allocs/op",
            "extra": "118765 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration)",
            "value": 10451,
            "unit": "ns/op\t    4208 B/op\t      51 allocs/op",
            "extra": "119865 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - ns/op",
            "value": 10451,
            "unit": "ns/op",
            "extra": "119865 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - B/op",
            "value": 4208,
            "unit": "B/op",
            "extra": "119865 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - allocs/op",
            "value": 51,
            "unit": "allocs/op",
            "extra": "119865 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration)",
            "value": 10419,
            "unit": "ns/op\t    4213 B/op\t      51 allocs/op",
            "extra": "119085 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - ns/op",
            "value": 10419,
            "unit": "ns/op",
            "extra": "119085 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - B/op",
            "value": 4213,
            "unit": "B/op",
            "extra": "119085 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - allocs/op",
            "value": 51,
            "unit": "allocs/op",
            "extra": "119085 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration)",
            "value": 10508,
            "unit": "ns/op\t    4207 B/op\t      51 allocs/op",
            "extra": "120051 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - ns/op",
            "value": 10508,
            "unit": "ns/op",
            "extra": "120051 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - B/op",
            "value": 4207,
            "unit": "B/op",
            "extra": "120051 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - allocs/op",
            "value": 51,
            "unit": "allocs/op",
            "extra": "120051 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration)",
            "value": 11015,
            "unit": "ns/op\t    4208 B/op\t      51 allocs/op",
            "extra": "119935 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - ns/op",
            "value": 11015,
            "unit": "ns/op",
            "extra": "119935 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - B/op",
            "value": 4208,
            "unit": "B/op",
            "extra": "119935 times\n4 procs"
          },
          {
            "name": "BenchmarkWalletService_Transfer (Qingyu_backend/test/integration) - allocs/op",
            "value": 51,
            "unit": "allocs/op",
            "extra": "119935 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration)",
            "value": 148766,
            "unit": "ns/op\t     250 B/op\t       7 allocs/op",
            "extra": "7840 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - ns/op",
            "value": 148766,
            "unit": "ns/op",
            "extra": "7840 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - B/op",
            "value": 250,
            "unit": "B/op",
            "extra": "7840 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "7840 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration)",
            "value": 147656,
            "unit": "ns/op\t     250 B/op\t       7 allocs/op",
            "extra": "8125 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - ns/op",
            "value": 147656,
            "unit": "ns/op",
            "extra": "8125 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - B/op",
            "value": 250,
            "unit": "B/op",
            "extra": "8125 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "8125 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration)",
            "value": 148674,
            "unit": "ns/op\t     251 B/op\t       7 allocs/op",
            "extra": "7934 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - ns/op",
            "value": 148674,
            "unit": "ns/op",
            "extra": "7934 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - B/op",
            "value": 251,
            "unit": "B/op",
            "extra": "7934 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "7934 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration)",
            "value": 148933,
            "unit": "ns/op\t     250 B/op\t       7 allocs/op",
            "extra": "8318 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - ns/op",
            "value": 148933,
            "unit": "ns/op",
            "extra": "8318 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - B/op",
            "value": 250,
            "unit": "B/op",
            "extra": "8318 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "8318 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration)",
            "value": 147153,
            "unit": "ns/op\t     250 B/op\t       7 allocs/op",
            "extra": "8361 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - ns/op",
            "value": 147153,
            "unit": "ns/op",
            "extra": "8361 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - B/op",
            "value": 250,
            "unit": "B/op",
            "extra": "8361 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Set (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "8361 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration)",
            "value": 146258,
            "unit": "ns/op\t     208 B/op\t       7 allocs/op",
            "extra": "8382 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - ns/op",
            "value": 146258,
            "unit": "ns/op",
            "extra": "8382 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "8382 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "8382 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration)",
            "value": 149408,
            "unit": "ns/op\t     208 B/op\t       7 allocs/op",
            "extra": "8163 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - ns/op",
            "value": 149408,
            "unit": "ns/op",
            "extra": "8163 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "8163 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "8163 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration)",
            "value": 145198,
            "unit": "ns/op\t     208 B/op\t       7 allocs/op",
            "extra": "8505 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - ns/op",
            "value": 145198,
            "unit": "ns/op",
            "extra": "8505 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "8505 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "8505 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration)",
            "value": 152134,
            "unit": "ns/op\t     208 B/op\t       7 allocs/op",
            "extra": "8026 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - ns/op",
            "value": 152134,
            "unit": "ns/op",
            "extra": "8026 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "8026 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "8026 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration)",
            "value": 144886,
            "unit": "ns/op\t     208 B/op\t       7 allocs/op",
            "extra": "8226 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - ns/op",
            "value": 144886,
            "unit": "ns/op",
            "extra": "8226 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "8226 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Get (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "8226 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration)",
            "value": 144881,
            "unit": "ns/op\t     185 B/op\t       6 allocs/op",
            "extra": "8256 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - ns/op",
            "value": 144881,
            "unit": "ns/op",
            "extra": "8256 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - B/op",
            "value": 185,
            "unit": "B/op",
            "extra": "8256 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "8256 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration)",
            "value": 144902,
            "unit": "ns/op\t     184 B/op\t       6 allocs/op",
            "extra": "7999 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - ns/op",
            "value": 144902,
            "unit": "ns/op",
            "extra": "7999 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - B/op",
            "value": 184,
            "unit": "B/op",
            "extra": "7999 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "7999 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration)",
            "value": 143779,
            "unit": "ns/op\t     184 B/op\t       6 allocs/op",
            "extra": "8274 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - ns/op",
            "value": 143779,
            "unit": "ns/op",
            "extra": "8274 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - B/op",
            "value": 184,
            "unit": "B/op",
            "extra": "8274 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "8274 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration)",
            "value": 144576,
            "unit": "ns/op\t     184 B/op\t       6 allocs/op",
            "extra": "8094 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - ns/op",
            "value": 144576,
            "unit": "ns/op",
            "extra": "8094 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - B/op",
            "value": 184,
            "unit": "B/op",
            "extra": "8094 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "8094 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration)",
            "value": 144134,
            "unit": "ns/op\t     184 B/op\t       6 allocs/op",
            "extra": "8274 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - ns/op",
            "value": 144134,
            "unit": "ns/op",
            "extra": "8274 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - B/op",
            "value": 184,
            "unit": "B/op",
            "extra": "8274 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/Incr (Qingyu_backend/test/integration) - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "8274 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration)",
            "value": 148306,
            "unit": "ns/op\t     248 B/op\t       7 allocs/op",
            "extra": "7786 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - ns/op",
            "value": 148306,
            "unit": "ns/op",
            "extra": "7786 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "7786 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "7786 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration)",
            "value": 150191,
            "unit": "ns/op\t     248 B/op\t       7 allocs/op",
            "extra": "7822 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - ns/op",
            "value": 150191,
            "unit": "ns/op",
            "extra": "7822 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "7822 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "7822 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration)",
            "value": 150808,
            "unit": "ns/op\t     248 B/op\t       7 allocs/op",
            "extra": "7230 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - ns/op",
            "value": 150808,
            "unit": "ns/op",
            "extra": "7230 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "7230 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "7230 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration)",
            "value": 147128,
            "unit": "ns/op\t     249 B/op\t       7 allocs/op",
            "extra": "8056 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - ns/op",
            "value": 147128,
            "unit": "ns/op",
            "extra": "8056 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - B/op",
            "value": 249,
            "unit": "B/op",
            "extra": "8056 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "8056 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration)",
            "value": 147035,
            "unit": "ns/op\t     248 B/op\t       7 allocs/op",
            "extra": "7978 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - ns/op",
            "value": 147035,
            "unit": "ns/op",
            "extra": "7978 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "7978 times\n4 procs"
          },
          {
            "name": "BenchmarkRedisOperations/HSet (Qingyu_backend/test/integration) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "7978 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance)",
            "value": 25.23,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "48797498 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - ns/op",
            "value": 25.23,
            "unit": "ns/op",
            "extra": "48797498 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "48797498 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "48797498 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance)",
            "value": 24.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "49626939 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - ns/op",
            "value": 24.2,
            "unit": "ns/op",
            "extra": "49626939 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "49626939 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "49626939 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance)",
            "value": 24.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "49714563 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - ns/op",
            "value": 24.5,
            "unit": "ns/op",
            "extra": "49714563 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "49714563 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "49714563 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance)",
            "value": 24.33,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "49648804 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - ns/op",
            "value": 24.33,
            "unit": "ns/op",
            "extra": "49648804 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "49648804 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "49648804 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance)",
            "value": 24.24,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "49571914 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - ns/op",
            "value": 24.24,
            "unit": "ns/op",
            "extra": "49571914 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "49571914 times\n4 procs"
          },
          {
            "name": "BenchmarkGetHomepageData (Qingyu_backend/test/performance) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "49571914 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance)",
            "value": 68.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16852789 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - ns/op",
            "value": 68.97,
            "unit": "ns/op",
            "extra": "16852789 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16852789 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16852789 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance)",
            "value": 68.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17451375 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - ns/op",
            "value": 68.97,
            "unit": "ns/op",
            "extra": "17451375 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17451375 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17451375 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance)",
            "value": 68.94,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17420493 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - ns/op",
            "value": 68.94,
            "unit": "ns/op",
            "extra": "17420493 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17420493 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17420493 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance)",
            "value": 69.22,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17444240 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - ns/op",
            "value": 69.22,
            "unit": "ns/op",
            "extra": "17444240 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17444240 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17444240 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance)",
            "value": 69.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14774446 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - ns/op",
            "value": 69.29,
            "unit": "ns/op",
            "extra": "14774446 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14774446 times\n4 procs"
          },
          {
            "name": "BenchmarkGetBookByID (Qingyu_backend/test/performance) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14774446 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance)",
            "value": 2427,
            "unit": "ns/op\t    3040 B/op\t      21 allocs/op",
            "extra": "488815 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - ns/op",
            "value": 2427,
            "unit": "ns/op",
            "extra": "488815 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - B/op",
            "value": 3040,
            "unit": "B/op",
            "extra": "488815 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "488815 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance)",
            "value": 2432,
            "unit": "ns/op\t    3040 B/op\t      21 allocs/op",
            "extra": "489392 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - ns/op",
            "value": 2432,
            "unit": "ns/op",
            "extra": "489392 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - B/op",
            "value": 3040,
            "unit": "B/op",
            "extra": "489392 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "489392 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance)",
            "value": 2426,
            "unit": "ns/op\t    3040 B/op\t      21 allocs/op",
            "extra": "470430 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - ns/op",
            "value": 2426,
            "unit": "ns/op",
            "extra": "470430 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - B/op",
            "value": 3040,
            "unit": "B/op",
            "extra": "470430 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "470430 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance)",
            "value": 2460,
            "unit": "ns/op\t    3040 B/op\t      21 allocs/op",
            "extra": "477811 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - ns/op",
            "value": 2460,
            "unit": "ns/op",
            "extra": "477811 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - B/op",
            "value": 3040,
            "unit": "B/op",
            "extra": "477811 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "477811 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance)",
            "value": 2466,
            "unit": "ns/op\t    3040 B/op\t      21 allocs/op",
            "extra": "470443 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - ns/op",
            "value": 2466,
            "unit": "ns/op",
            "extra": "470443 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - B/op",
            "value": 3040,
            "unit": "B/op",
            "extra": "470443 times\n4 procs"
          },
          {
            "name": "BenchmarkGetRankings (Qingyu_backend/test/performance) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "470443 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance)",
            "value": 782734,
            "unit": "ns/op\t  284553 B/op\t    6534 allocs/op",
            "extra": "1507 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - ns/op",
            "value": 782734,
            "unit": "ns/op",
            "extra": "1507 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - B/op",
            "value": 284553,
            "unit": "B/op",
            "extra": "1507 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - allocs/op",
            "value": 6534,
            "unit": "allocs/op",
            "extra": "1507 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance)",
            "value": 783176,
            "unit": "ns/op\t  284570 B/op\t    6534 allocs/op",
            "extra": "1587 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - ns/op",
            "value": 783176,
            "unit": "ns/op",
            "extra": "1587 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - B/op",
            "value": 284570,
            "unit": "B/op",
            "extra": "1587 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - allocs/op",
            "value": 6534,
            "unit": "allocs/op",
            "extra": "1587 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance)",
            "value": 870379,
            "unit": "ns/op\t  284527 B/op\t    6533 allocs/op",
            "extra": "1543 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - ns/op",
            "value": 870379,
            "unit": "ns/op",
            "extra": "1543 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - B/op",
            "value": 284527,
            "unit": "B/op",
            "extra": "1543 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - allocs/op",
            "value": 6533,
            "unit": "allocs/op",
            "extra": "1543 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance)",
            "value": 780823,
            "unit": "ns/op\t  284547 B/op\t    6534 allocs/op",
            "extra": "1503 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - ns/op",
            "value": 780823,
            "unit": "ns/op",
            "extra": "1503 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - B/op",
            "value": 284547,
            "unit": "B/op",
            "extra": "1503 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - allocs/op",
            "value": 6534,
            "unit": "allocs/op",
            "extra": "1503 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance)",
            "value": 779511,
            "unit": "ns/op\t  284538 B/op\t    6534 allocs/op",
            "extra": "1513 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - ns/op",
            "value": 779511,
            "unit": "ns/op",
            "extra": "1513 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - B/op",
            "value": 284538,
            "unit": "B/op",
            "extra": "1513 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_NoCursor (Qingyu_backend/test/performance) - allocs/op",
            "value": 6534,
            "unit": "allocs/op",
            "extra": "1513 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance)",
            "value": 259065,
            "unit": "ns/op\t   55891 B/op\t    1079 allocs/op",
            "extra": "4454 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - ns/op",
            "value": 259065,
            "unit": "ns/op",
            "extra": "4454 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - B/op",
            "value": 55891,
            "unit": "B/op",
            "extra": "4454 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - allocs/op",
            "value": 1079,
            "unit": "allocs/op",
            "extra": "4454 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance)",
            "value": 259822,
            "unit": "ns/op\t   55884 B/op\t    1079 allocs/op",
            "extra": "4401 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - ns/op",
            "value": 259822,
            "unit": "ns/op",
            "extra": "4401 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - B/op",
            "value": 55884,
            "unit": "B/op",
            "extra": "4401 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - allocs/op",
            "value": 1079,
            "unit": "allocs/op",
            "extra": "4401 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance)",
            "value": 256747,
            "unit": "ns/op\t   55894 B/op\t    1079 allocs/op",
            "extra": "4464 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - ns/op",
            "value": 256747,
            "unit": "ns/op",
            "extra": "4464 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - B/op",
            "value": 55894,
            "unit": "B/op",
            "extra": "4464 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - allocs/op",
            "value": 1079,
            "unit": "allocs/op",
            "extra": "4464 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance)",
            "value": 258209,
            "unit": "ns/op\t   55880 B/op\t    1079 allocs/op",
            "extra": "4485 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - ns/op",
            "value": 258209,
            "unit": "ns/op",
            "extra": "4485 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - B/op",
            "value": 55880,
            "unit": "B/op",
            "extra": "4485 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - allocs/op",
            "value": 1079,
            "unit": "allocs/op",
            "extra": "4485 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance)",
            "value": 257477,
            "unit": "ns/op\t   55891 B/op\t    1079 allocs/op",
            "extra": "4424 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - ns/op",
            "value": 257477,
            "unit": "ns/op",
            "extra": "4424 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - B/op",
            "value": 55891,
            "unit": "B/op",
            "extra": "4424 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearch_WithCursor (Qingyu_backend/test/performance) - allocs/op",
            "value": 1079,
            "unit": "allocs/op",
            "extra": "4424 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance)",
            "value": 1703,
            "unit": "ns/op\t     984 B/op\t      21 allocs/op",
            "extra": "670614 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - ns/op",
            "value": 1703,
            "unit": "ns/op",
            "extra": "670614 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - B/op",
            "value": 984,
            "unit": "B/op",
            "extra": "670614 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "670614 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance)",
            "value": 1887,
            "unit": "ns/op\t     984 B/op\t      21 allocs/op",
            "extra": "674380 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - ns/op",
            "value": 1887,
            "unit": "ns/op",
            "extra": "674380 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - B/op",
            "value": 984,
            "unit": "B/op",
            "extra": "674380 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "674380 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance)",
            "value": 1695,
            "unit": "ns/op\t     984 B/op\t      21 allocs/op",
            "extra": "671071 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - ns/op",
            "value": 1695,
            "unit": "ns/op",
            "extra": "671071 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - B/op",
            "value": 984,
            "unit": "B/op",
            "extra": "671071 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "671071 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance)",
            "value": 1706,
            "unit": "ns/op\t     984 B/op\t      21 allocs/op",
            "extra": "686632 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - ns/op",
            "value": 1706,
            "unit": "ns/op",
            "extra": "686632 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - B/op",
            "value": 984,
            "unit": "B/op",
            "extra": "686632 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "686632 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance)",
            "value": 1710,
            "unit": "ns/op\t     984 B/op\t      21 allocs/op",
            "extra": "678102 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - ns/op",
            "value": 1710,
            "unit": "ns/op",
            "extra": "678102 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - B/op",
            "value": 984,
            "unit": "B/op",
            "extra": "678102 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorEncoding (Qingyu_backend/test/performance) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "678102 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance)",
            "value": 1370,
            "unit": "ns/op\t     712 B/op\t      16 allocs/op",
            "extra": "819109 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - ns/op",
            "value": 1370,
            "unit": "ns/op",
            "extra": "819109 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - B/op",
            "value": 712,
            "unit": "B/op",
            "extra": "819109 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "819109 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance)",
            "value": 1394,
            "unit": "ns/op\t     712 B/op\t      16 allocs/op",
            "extra": "809959 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - ns/op",
            "value": 1394,
            "unit": "ns/op",
            "extra": "809959 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - B/op",
            "value": 712,
            "unit": "B/op",
            "extra": "809959 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "809959 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance)",
            "value": 1371,
            "unit": "ns/op\t     712 B/op\t      16 allocs/op",
            "extra": "830946 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - ns/op",
            "value": 1371,
            "unit": "ns/op",
            "extra": "830946 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - B/op",
            "value": 712,
            "unit": "B/op",
            "extra": "830946 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "830946 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance)",
            "value": 1372,
            "unit": "ns/op\t     712 B/op\t      16 allocs/op",
            "extra": "838878 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - ns/op",
            "value": 1372,
            "unit": "ns/op",
            "extra": "838878 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - B/op",
            "value": 712,
            "unit": "B/op",
            "extra": "838878 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "838878 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance)",
            "value": 1380,
            "unit": "ns/op\t     712 B/op\t      16 allocs/op",
            "extra": "819871 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - ns/op",
            "value": 1380,
            "unit": "ns/op",
            "extra": "819871 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - B/op",
            "value": 712,
            "unit": "B/op",
            "extra": "819871 times\n4 procs"
          },
          {
            "name": "BenchmarkCursorDecoding (Qingyu_backend/test/performance) - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "819871 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance)",
            "value": 224663,
            "unit": "ns/op\t   92832 B/op\t    2310 allocs/op",
            "extra": "4532 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - ns/op",
            "value": 224663,
            "unit": "ns/op",
            "extra": "4532 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - B/op",
            "value": 92832,
            "unit": "B/op",
            "extra": "4532 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - allocs/op",
            "value": 2310,
            "unit": "allocs/op",
            "extra": "4532 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance)",
            "value": 202260,
            "unit": "ns/op\t   92832 B/op\t    2310 allocs/op",
            "extra": "6079 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - ns/op",
            "value": 202260,
            "unit": "ns/op",
            "extra": "6079 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - B/op",
            "value": 92832,
            "unit": "B/op",
            "extra": "6079 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - allocs/op",
            "value": 2310,
            "unit": "allocs/op",
            "extra": "6079 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance)",
            "value": 203756,
            "unit": "ns/op\t   92833 B/op\t    2310 allocs/op",
            "extra": "5830 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - ns/op",
            "value": 203756,
            "unit": "ns/op",
            "extra": "5830 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - B/op",
            "value": 92833,
            "unit": "B/op",
            "extra": "5830 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - allocs/op",
            "value": 2310,
            "unit": "allocs/op",
            "extra": "5830 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance)",
            "value": 204960,
            "unit": "ns/op\t   92832 B/op\t    2310 allocs/op",
            "extra": "5815 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - ns/op",
            "value": 204960,
            "unit": "ns/op",
            "extra": "5815 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - B/op",
            "value": 92832,
            "unit": "B/op",
            "extra": "5815 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - allocs/op",
            "value": 2310,
            "unit": "allocs/op",
            "extra": "5815 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance)",
            "value": 203297,
            "unit": "ns/op\t   92832 B/op\t    2310 allocs/op",
            "extra": "5761 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - ns/op",
            "value": 203297,
            "unit": "ns/op",
            "extra": "5761 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - B/op",
            "value": 92832,
            "unit": "B/op",
            "extra": "5761 times\n4 procs"
          },
          {
            "name": "BenchmarkNDJSONParsing (Qingyu_backend/test/performance) - allocs/op",
            "value": 2310,
            "unit": "allocs/op",
            "extra": "5761 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance)",
            "value": 37491181,
            "unit": "ns/op\t22936521 B/op\t  503158 allocs/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - ns/op",
            "value": 37491181,
            "unit": "ns/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - B/op",
            "value": 22936521,
            "unit": "B/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - allocs/op",
            "value": 503158,
            "unit": "allocs/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance)",
            "value": 36549440,
            "unit": "ns/op\t22935928 B/op\t  503154 allocs/op",
            "extra": "32 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - ns/op",
            "value": 36549440,
            "unit": "ns/op",
            "extra": "32 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - B/op",
            "value": 22935928,
            "unit": "B/op",
            "extra": "32 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - allocs/op",
            "value": 503154,
            "unit": "allocs/op",
            "extra": "32 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance)",
            "value": 36798035,
            "unit": "ns/op\t22937798 B/op\t  503171 allocs/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - ns/op",
            "value": 36798035,
            "unit": "ns/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - B/op",
            "value": 22937798,
            "unit": "B/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - allocs/op",
            "value": 503171,
            "unit": "allocs/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance)",
            "value": 37825630,
            "unit": "ns/op\t22936159 B/op\t  503151 allocs/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - ns/op",
            "value": 37825630,
            "unit": "ns/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - B/op",
            "value": 22936159,
            "unit": "B/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - allocs/op",
            "value": 503151,
            "unit": "allocs/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance)",
            "value": 36127921,
            "unit": "ns/op\t22936032 B/op\t  503152 allocs/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - ns/op",
            "value": 36127921,
            "unit": "ns/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - B/op",
            "value": 22936032,
            "unit": "B/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkStreamSearchLargeDataset (Qingyu_backend/test/performance) - allocs/op",
            "value": 503152,
            "unit": "allocs/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance)",
            "value": 101963,
            "unit": "ns/op\t   58034 B/op\t    1088 allocs/op",
            "extra": "12126 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - ns/op",
            "value": 101963,
            "unit": "ns/op",
            "extra": "12126 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - B/op",
            "value": 58034,
            "unit": "B/op",
            "extra": "12126 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - allocs/op",
            "value": 1088,
            "unit": "allocs/op",
            "extra": "12126 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance)",
            "value": 102920,
            "unit": "ns/op\t   58074 B/op\t    1088 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - ns/op",
            "value": 102920,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - B/op",
            "value": 58074,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - allocs/op",
            "value": 1088,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance)",
            "value": 103111,
            "unit": "ns/op\t   58067 B/op\t    1088 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - ns/op",
            "value": 103111,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - B/op",
            "value": 58067,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - allocs/op",
            "value": 1088,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance)",
            "value": 101929,
            "unit": "ns/op\t   58041 B/op\t    1088 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - ns/op",
            "value": 101929,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - B/op",
            "value": 58041,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - allocs/op",
            "value": 1088,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance)",
            "value": 100511,
            "unit": "ns/op\t   57986 B/op\t    1088 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - ns/op",
            "value": 100511,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - B/op",
            "value": 57986,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentStreamRequests (Qingyu_backend/test/performance) - allocs/op",
            "value": 1088,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user)",
            "value": 341.7,
            "unit": "ns/op\t     192 B/op\t       3 allocs/op",
            "extra": "3546368 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - ns/op",
            "value": 341.7,
            "unit": "ns/op",
            "extra": "3546368 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "3546368 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "3546368 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user)",
            "value": 368.1,
            "unit": "ns/op\t     192 B/op\t       3 allocs/op",
            "extra": "3548650 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - ns/op",
            "value": 368.1,
            "unit": "ns/op",
            "extra": "3548650 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "3548650 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "3548650 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user)",
            "value": 343.6,
            "unit": "ns/op\t     192 B/op\t       3 allocs/op",
            "extra": "3540235 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - ns/op",
            "value": 343.6,
            "unit": "ns/op",
            "extra": "3540235 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "3540235 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "3540235 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user)",
            "value": 338.1,
            "unit": "ns/op\t     192 B/op\t       3 allocs/op",
            "extra": "3464808 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - ns/op",
            "value": 338.1,
            "unit": "ns/op",
            "extra": "3464808 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "3464808 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "3464808 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user)",
            "value": 337.3,
            "unit": "ns/op\t     192 B/op\t       3 allocs/op",
            "extra": "3525456 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - ns/op",
            "value": 337.3,
            "unit": "ns/op",
            "extra": "3525456 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "3525456 times\n4 procs"
          },
          {
            "name": "BenchmarkGenerateToken (Qingyu_backend/test/service/user) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "3525456 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user)",
            "value": 70.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15425497 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - ns/op",
            "value": 70.6,
            "unit": "ns/op",
            "extra": "15425497 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15425497 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15425497 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user)",
            "value": 70.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17008630 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - ns/op",
            "value": 70.66,
            "unit": "ns/op",
            "extra": "17008630 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17008630 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17008630 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user)",
            "value": 71.17,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16996071 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - ns/op",
            "value": 71.17,
            "unit": "ns/op",
            "extra": "16996071 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16996071 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16996071 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user)",
            "value": 71.82,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16982858 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - ns/op",
            "value": 71.82,
            "unit": "ns/op",
            "extra": "16982858 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16982858 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16982858 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user)",
            "value": 70.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15686750 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - ns/op",
            "value": 70.67,
            "unit": "ns/op",
            "extra": "15686750 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15686750 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateToken (Qingyu_backend/test/service/user) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15686750 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark)",
            "value": 267764,
            "unit": "ns/op\t    4597 B/op\t      58 allocs/op",
            "extra": "4558 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267764,
            "unit": "ns/op",
            "extra": "4558 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4597,
            "unit": "B/op",
            "extra": "4558 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4558 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark)",
            "value": 266422,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4758 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 266422,
            "unit": "ns/op",
            "extra": "4758 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4758 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4758 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark)",
            "value": 266998,
            "unit": "ns/op\t    4595 B/op\t      58 allocs/op",
            "extra": "4772 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 266998,
            "unit": "ns/op",
            "extra": "4772 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4595,
            "unit": "B/op",
            "extra": "4772 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4772 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark)",
            "value": 262169,
            "unit": "ns/op\t    4597 B/op\t      58 allocs/op",
            "extra": "4760 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 262169,
            "unit": "ns/op",
            "extra": "4760 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4597,
            "unit": "B/op",
            "extra": "4760 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4760 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark)",
            "value": 267535,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4765 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267535,
            "unit": "ns/op",
            "extra": "4765 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4765 times\n4 procs"
          },
          {
            "name": "BenchmarkBookSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4765 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark)",
            "value": 266817,
            "unit": "ns/op\t    4595 B/op\t      58 allocs/op",
            "extra": "4780 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 266817,
            "unit": "ns/op",
            "extra": "4780 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4595,
            "unit": "B/op",
            "extra": "4780 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4780 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark)",
            "value": 268069,
            "unit": "ns/op\t    4595 B/op\t      58 allocs/op",
            "extra": "4753 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 268069,
            "unit": "ns/op",
            "extra": "4753 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4595,
            "unit": "B/op",
            "extra": "4753 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4753 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark)",
            "value": 267362,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4540 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267362,
            "unit": "ns/op",
            "extra": "4540 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4540 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4540 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark)",
            "value": 267222,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4588 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267222,
            "unit": "ns/op",
            "extra": "4588 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4588 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4588 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark)",
            "value": 268006,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4567 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 268006,
            "unit": "ns/op",
            "extra": "4567 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4567 times\n4 procs"
          },
          {
            "name": "BenchmarkProjectSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4567 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark)",
            "value": 268046,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4576 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 268046,
            "unit": "ns/op",
            "extra": "4576 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4576 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4576 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark)",
            "value": 267212,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4590 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267212,
            "unit": "ns/op",
            "extra": "4590 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4590 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4590 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark)",
            "value": 263969,
            "unit": "ns/op\t    4595 B/op\t      58 allocs/op",
            "extra": "4586 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 263969,
            "unit": "ns/op",
            "extra": "4586 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4595,
            "unit": "B/op",
            "extra": "4586 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4586 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark)",
            "value": 267557,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4771 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267557,
            "unit": "ns/op",
            "extra": "4771 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4771 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4771 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark)",
            "value": 267178,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4780 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267178,
            "unit": "ns/op",
            "extra": "4780 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4780 times\n4 procs"
          },
          {
            "name": "BenchmarkDocumentSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4780 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark)",
            "value": 266704,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4568 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 266704,
            "unit": "ns/op",
            "extra": "4568 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4568 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4568 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark)",
            "value": 266817,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4584 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 266817,
            "unit": "ns/op",
            "extra": "4584 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4584 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4584 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark)",
            "value": 267797,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4573 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267797,
            "unit": "ns/op",
            "extra": "4573 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4573 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4573 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark)",
            "value": 267375,
            "unit": "ns/op\t    4595 B/op\t      58 allocs/op",
            "extra": "4684 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267375,
            "unit": "ns/op",
            "extra": "4684 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4595,
            "unit": "B/op",
            "extra": "4684 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4684 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark)",
            "value": 267189,
            "unit": "ns/op\t    4595 B/op\t      58 allocs/op",
            "extra": "4568 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267189,
            "unit": "ns/op",
            "extra": "4568 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4595,
            "unit": "B/op",
            "extra": "4568 times\n4 procs"
          },
          {
            "name": "BenchmarkUserSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4568 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark)",
            "value": 266721,
            "unit": "ns/op\t    4597 B/op\t      58 allocs/op",
            "extra": "4560 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 266721,
            "unit": "ns/op",
            "extra": "4560 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4597,
            "unit": "B/op",
            "extra": "4560 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4560 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark)",
            "value": 267017,
            "unit": "ns/op\t    4595 B/op\t      58 allocs/op",
            "extra": "4555 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267017,
            "unit": "ns/op",
            "extra": "4555 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4595,
            "unit": "B/op",
            "extra": "4555 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4555 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark)",
            "value": 267267,
            "unit": "ns/op\t    4595 B/op\t      58 allocs/op",
            "extra": "4958 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267267,
            "unit": "ns/op",
            "extra": "4958 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4595,
            "unit": "B/op",
            "extra": "4958 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4958 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark)",
            "value": 269097,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4777 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 269097,
            "unit": "ns/op",
            "extra": "4777 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4777 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4777 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark)",
            "value": 267380,
            "unit": "ns/op\t    4596 B/op\t      58 allocs/op",
            "extra": "4777 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267380,
            "unit": "ns/op",
            "extra": "4777 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4596,
            "unit": "B/op",
            "extra": "4777 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentSearch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "4777 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark)",
            "value": 5343,
            "unit": "ns/op\t    4582 B/op\t      57 allocs/op",
            "extra": "268509 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 5343,
            "unit": "ns/op",
            "extra": "268509 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4582,
            "unit": "B/op",
            "extra": "268509 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 57,
            "unit": "allocs/op",
            "extra": "268509 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark)",
            "value": 4620,
            "unit": "ns/op\t    4581 B/op\t      57 allocs/op",
            "extra": "259947 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 4620,
            "unit": "ns/op",
            "extra": "259947 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4581,
            "unit": "B/op",
            "extra": "259947 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 57,
            "unit": "allocs/op",
            "extra": "259947 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark)",
            "value": 4707,
            "unit": "ns/op\t    4580 B/op\t      57 allocs/op",
            "extra": "258172 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 4707,
            "unit": "ns/op",
            "extra": "258172 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4580,
            "unit": "B/op",
            "extra": "258172 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 57,
            "unit": "allocs/op",
            "extra": "258172 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark)",
            "value": 4770,
            "unit": "ns/op\t    4580 B/op\t      57 allocs/op",
            "extra": "262969 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 4770,
            "unit": "ns/op",
            "extra": "262969 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4580,
            "unit": "B/op",
            "extra": "262969 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 57,
            "unit": "allocs/op",
            "extra": "262969 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark)",
            "value": 4682,
            "unit": "ns/op\t    4580 B/op\t      57 allocs/op",
            "extra": "237794 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 4682,
            "unit": "ns/op",
            "extra": "237794 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4580,
            "unit": "B/op",
            "extra": "237794 times\n4 procs"
          },
          {
            "name": "BenchmarkHighQPSLoad (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 57,
            "unit": "allocs/op",
            "extra": "237794 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark)",
            "value": 266051,
            "unit": "ns/op\t    4515 B/op\t      56 allocs/op",
            "extra": "4767 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 266051,
            "unit": "ns/op",
            "extra": "4767 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4515,
            "unit": "B/op",
            "extra": "4767 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 56,
            "unit": "allocs/op",
            "extra": "4767 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark)",
            "value": 267889,
            "unit": "ns/op\t    4514 B/op\t      56 allocs/op",
            "extra": "4584 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267889,
            "unit": "ns/op",
            "extra": "4584 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4514,
            "unit": "B/op",
            "extra": "4584 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 56,
            "unit": "allocs/op",
            "extra": "4584 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark)",
            "value": 266991,
            "unit": "ns/op\t    4514 B/op\t      56 allocs/op",
            "extra": "4773 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 266991,
            "unit": "ns/op",
            "extra": "4773 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4514,
            "unit": "B/op",
            "extra": "4773 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 56,
            "unit": "allocs/op",
            "extra": "4773 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark)",
            "value": 267266,
            "unit": "ns/op\t    4514 B/op\t      56 allocs/op",
            "extra": "4573 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267266,
            "unit": "ns/op",
            "extra": "4573 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4514,
            "unit": "B/op",
            "extra": "4573 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 56,
            "unit": "allocs/op",
            "extra": "4573 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark)",
            "value": 267347,
            "unit": "ns/op\t    4514 B/op\t      56 allocs/op",
            "extra": "4562 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267347,
            "unit": "ns/op",
            "extra": "4562 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4514,
            "unit": "B/op",
            "extra": "4562 times\n4 procs"
          },
          {
            "name": "BenchmarkGinRouting (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 56,
            "unit": "allocs/op",
            "extra": "4562 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark)",
            "value": 3925,
            "unit": "ns/op\t    2670 B/op\t      44 allocs/op",
            "extra": "258690 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 3925,
            "unit": "ns/op",
            "extra": "258690 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - B/op",
            "value": 2670,
            "unit": "B/op",
            "extra": "258690 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "258690 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark)",
            "value": 4142,
            "unit": "ns/op\t    3146 B/op\t      44 allocs/op",
            "extra": "270002 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 4142,
            "unit": "ns/op",
            "extra": "270002 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - B/op",
            "value": 3146,
            "unit": "B/op",
            "extra": "270002 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "270002 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark)",
            "value": 4079,
            "unit": "ns/op\t    3133 B/op\t      44 allocs/op",
            "extra": "273534 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 4079,
            "unit": "ns/op",
            "extra": "273534 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - B/op",
            "value": 3133,
            "unit": "B/op",
            "extra": "273534 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "273534 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark)",
            "value": 4163,
            "unit": "ns/op\t    3151 B/op\t      44 allocs/op",
            "extra": "268546 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 4163,
            "unit": "ns/op",
            "extra": "268546 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - B/op",
            "value": 3151,
            "unit": "B/op",
            "extra": "268546 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "268546 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark)",
            "value": 4127,
            "unit": "ns/op\t    3133 B/op\t      44 allocs/op",
            "extra": "273501 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 4127,
            "unit": "ns/op",
            "extra": "273501 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - B/op",
            "value": 3133,
            "unit": "B/op",
            "extra": "273501 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONSerialization (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 44,
            "unit": "allocs/op",
            "extra": "273501 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark)",
            "value": 412.5,
            "unit": "ns/op\t     272 B/op\t       4 allocs/op",
            "extra": "2913859 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 412.5,
            "unit": "ns/op",
            "extra": "2913859 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 272,
            "unit": "B/op",
            "extra": "2913859 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2913859 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark)",
            "value": 413.4,
            "unit": "ns/op\t     272 B/op\t       4 allocs/op",
            "extra": "2909522 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 413.4,
            "unit": "ns/op",
            "extra": "2909522 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 272,
            "unit": "B/op",
            "extra": "2909522 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2909522 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark)",
            "value": 409.7,
            "unit": "ns/op\t     272 B/op\t       4 allocs/op",
            "extra": "2931226 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 409.7,
            "unit": "ns/op",
            "extra": "2931226 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 272,
            "unit": "B/op",
            "extra": "2931226 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2931226 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark)",
            "value": 410.1,
            "unit": "ns/op\t     272 B/op\t       4 allocs/op",
            "extra": "2928928 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 410.1,
            "unit": "ns/op",
            "extra": "2928928 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 272,
            "unit": "B/op",
            "extra": "2928928 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2928928 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark)",
            "value": 412.2,
            "unit": "ns/op\t     272 B/op\t       4 allocs/op",
            "extra": "2944580 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 412.2,
            "unit": "ns/op",
            "extra": "2944580 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - B/op",
            "value": 272,
            "unit": "B/op",
            "extra": "2944580 times\n4 procs"
          },
          {
            "name": "BenchmarkContextSwitch (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2944580 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark)",
            "value": 266080,
            "unit": "ns/op\t    4515 B/op\t      56 allocs/op",
            "extra": "4752 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 266080,
            "unit": "ns/op",
            "extra": "4752 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4515,
            "unit": "B/op",
            "extra": "4752 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 56,
            "unit": "allocs/op",
            "extra": "4752 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark)",
            "value": 265963,
            "unit": "ns/op\t    4515 B/op\t      56 allocs/op",
            "extra": "4708 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 265963,
            "unit": "ns/op",
            "extra": "4708 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4515,
            "unit": "B/op",
            "extra": "4708 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 56,
            "unit": "allocs/op",
            "extra": "4708 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark)",
            "value": 266682,
            "unit": "ns/op\t    4514 B/op\t      56 allocs/op",
            "extra": "4578 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 266682,
            "unit": "ns/op",
            "extra": "4578 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4514,
            "unit": "B/op",
            "extra": "4578 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 56,
            "unit": "allocs/op",
            "extra": "4578 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark)",
            "value": 266549,
            "unit": "ns/op\t    4514 B/op\t      56 allocs/op",
            "extra": "4759 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 266549,
            "unit": "ns/op",
            "extra": "4759 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4514,
            "unit": "B/op",
            "extra": "4759 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 56,
            "unit": "allocs/op",
            "extra": "4759 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark)",
            "value": 267093,
            "unit": "ns/op\t    4515 B/op\t      56 allocs/op",
            "extra": "4779 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - ns/op",
            "value": 267093,
            "unit": "ns/op",
            "extra": "4779 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - B/op",
            "value": 4515,
            "unit": "B/op",
            "extra": "4779 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequest (Qingyu_backend/tests/benchmark) - allocs/op",
            "value": 56,
            "unit": "allocs/op",
            "extra": "4779 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e)",
            "value": 2784,
            "unit": "ns/op\t    3105 B/op\t      28 allocs/op",
            "extra": "396104 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - ns/op",
            "value": 2784,
            "unit": "ns/op",
            "extra": "396104 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - B/op",
            "value": 3105,
            "unit": "B/op",
            "extra": "396104 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "396104 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e)",
            "value": 2802,
            "unit": "ns/op\t    3105 B/op\t      28 allocs/op",
            "extra": "410230 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - ns/op",
            "value": 2802,
            "unit": "ns/op",
            "extra": "410230 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - B/op",
            "value": 3105,
            "unit": "B/op",
            "extra": "410230 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "410230 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e)",
            "value": 2802,
            "unit": "ns/op\t    3105 B/op\t      28 allocs/op",
            "extra": "396153 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - ns/op",
            "value": 2802,
            "unit": "ns/op",
            "extra": "396153 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - B/op",
            "value": 3105,
            "unit": "B/op",
            "extra": "396153 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "396153 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e)",
            "value": 2812,
            "unit": "ns/op\t    3105 B/op\t      28 allocs/op",
            "extra": "416785 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - ns/op",
            "value": 2812,
            "unit": "ns/op",
            "extra": "416785 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - B/op",
            "value": 3105,
            "unit": "B/op",
            "extra": "416785 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "416785 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e)",
            "value": 3172,
            "unit": "ns/op\t    3105 B/op\t      28 allocs/op",
            "extra": "411396 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - ns/op",
            "value": 3172,
            "unit": "ns/op",
            "extra": "411396 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - B/op",
            "value": 3105,
            "unit": "B/op",
            "extra": "411396 times\n4 procs"
          },
          {
            "name": "BenchmarkPermissionCheck (Qingyu_backend/tests/e2e) - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "411396 times\n4 procs"
          }
        ]
      }
    ]
  }
}