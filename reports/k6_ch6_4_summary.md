# k6 6.4 Performance Summary

GeneratedAt: 2026-02-23 18:25:31

| Stage | VU | Duration | Avg(ms) | P95(ms) | RPS | Error Rate | Requests |
|---|---:|---|---:|---:|---:|---:|---:|
| baseline | 10 | 1min | 8.56 | 23.64 | 29.12 | 0% | 1770 |
| normal | 50 | 3min | 10.26 | 29.73 | 145.01 | 0% | 26247 |
| stress | 100 | 3min | 78.67 | 380.04 | 241.62 | 0% | 43737 |

## Threshold Review
- **baseline**: P95<800ms `PASS`; ErrorRate<1% `PASS`
- **normal**: P95<800ms `PASS`; ErrorRate<1% `PASS`
- **stress**: P95<800ms `PASS`; ErrorRate<1% `PASS`
