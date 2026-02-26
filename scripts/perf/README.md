# k6 Performance Test (Thesis 6.4)

This script is for thesis section 6.4 and focuses on reproducible API performance evidence.

## 1) Prerequisites

- Backend service is running (default: `http://127.0.0.1:9090`)
- k6 installed

## 2) Script

- `scripts/perf/k6_ch6_4.js`

Default test APIs:

- `GET /health`
- `GET /api/v1/bookstore/homepage`
- `GET /api/v1/bookstore/books?page=1&pageSize=10`

Optional authenticated API (enabled only when login env is provided):

- `POST /api/v1/user/auth/login`
- `GET /api/v1/user/profile`

## 3) Run Commands

Baseline (10 VU, 1 min):

```bash
k6 run -e TEST_STAGE=baseline scripts/perf/k6_ch6_4.js
```

Normal load (50 VU, 3 min):

```bash
k6 run -e TEST_STAGE=normal scripts/perf/k6_ch6_4.js
```

Stress boundary (100 VU, 3 min):

```bash
k6 run -e TEST_STAGE=stress scripts/perf/k6_ch6_4.js
```

Export JSON summary (for screenshot evidence):

```bash
k6 run -e TEST_STAGE=normal --summary-export=reports/k6_ch6_4_normal.json scripts/perf/k6_ch6_4.js
```

## 4) Optional: test authenticated endpoint

```bash
k6 run -e TEST_STAGE=normal -e LOGIN_USERNAME=your_user -e LOGIN_PASSWORD=your_password scripts/perf/k6_ch6_4.js
```

## 5) Screenshot suggestions for thesis

1. k6 terminal summary (`avg`, `p95`, `http_req_failed`, `http_reqs`)
2. exported JSON file content (`reports/k6_ch6_4_normal.json`)
3. one run command with full parameters

