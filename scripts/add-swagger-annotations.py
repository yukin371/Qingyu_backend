#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Swagger 注释自动生成脚本

功能：
1. 扫描 router/*/_router.go 文件
2. 解析 Gin 路由定义
3. 在对应 Handler 函数前插入 Swagger 注释模板

使用：
    python scripts/add-swagger-annotations.py --module reader
    python scripts/add-swagger-annotations.py --all
"""

import os
import re
import sys
import argparse
from pathlib import Path
from typing import List, Dict, Tuple

# 设置标准输出编码为UTF-8（Windows兼容）
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

# 项目根目录
ROOT_DIR = Path(__file__).parent.parent
ROUTER_DIR = ROOT_DIR / "router"
API_DIR = ROOT_DIR / "api/v1"


def main():
    parser = argparse.ArgumentParser(description="生成 Swagger 注释模板")
    parser.add_argument("--module", type=str, help="指定模块名称（如 reader）")
    parser.add_argument("--all", action="store_true", help="处理所有模块")
    parser.add_argument("--dry-run", action="store_true", help="预览模式，不修改文件")
    parser.add_argument("--verbose", action="store_true", help="详细输出")

    args = parser.parse_args()

    if args.all:
        modules = ["reader", "admin", "bookstore", "writer", "social", "ai", "user", "shared"]
    elif args.module:
        modules = [args.module]
    else:
        parser.print_help()
        sys.exit(1)

    for module in modules:
        process_module(module, args.dry_run, args.verbose)


def process_module(module_name: str, dry_run: bool, verbose: bool):
    """处理单个模块的路由"""
    print(f"\n{'='*60}")
    print(f"处理模块: {module_name}")
    print(f"{'='*60}")

    router_file = ROUTER_DIR / module_name / f"{module_name}_router.go"

    if not router_file.exists():
        print(f"❌ 路由文件不存在: {router_file}")
        return

    # TODO: 实现路由解析和注释生成
    print(f"✓ 找到路由文件: {router_file}")


if __name__ == "__main__":
    main()
