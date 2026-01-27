#!/usr/bin/env python3
"""解析测试结果JSON"""

import sys
import json
from pathlib import Path

def load_result(filename):
    """加载测试结果"""
    with open(filename, 'r') as f:
        return json.load(f)

def print_summary(result):
    """打印结果摘要"""
    print(f"场景: {result.get('scenario', 'N/A')}")
    print(f"总请求数: {result.get('total_requests', 'N/A')}")
    print(f"成功: {result.get('success_count', 'N/A')}")
    print(f"失败: {result.get('error_count', 'N/A')}")
    print(f"平均延迟: {result.get('avg_latency', 'N/A')}")
    print(f"P95延迟: {result.get('p95_latency', 'N/A')}")
    print(f"P99延迟: {result.get('p99_latency', 'N/A')}")
    print(f"吞吐量: {result.get('throughput', 'N/A')} req/s")

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("Usage: python parse_ab_result.py <result_json_file>")
        sys.exit(1)

    result = load_result(sys.argv[1])
    print_summary(result)
