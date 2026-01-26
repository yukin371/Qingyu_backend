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


def parse_gin_routes(router_file: Path) -> List[Dict]:
    """
    解析 Gin 路由文件，提取路由信息

    返回格式：
    [
        {
            "method": "GET",
            "path": "/books/:bookId",
            "handler": "GetBooks",
            "line": 42
        },
        ...
    ]
    """
    routes = []
    content = router_file.read_text(encoding="utf-8")

    # 匹配路由定义，例如：
    # readerGroup.GET("/books/:bookId", readerApiHandler.GetBooks)
    pattern = re.compile(
        r'(\w+)\.(GET|POST|PUT|DELETE|PATCH)\(["\']([^"\']+)["\'],\s*\w+\.(\w+)\)',
        re.MULTILINE
    )

    for match in pattern.finditer(content):
        group_name, method, path, handler_name = match.groups()
        routes.append({
            "group": group_name,
            "method": method.upper(),
            "path": path,
            "handler": handler_name,
            "line": content[:match.start()].count('\n') + 1
        })

    return routes


def extract_path_params(path: str) -> List[str]:
    """提取路径参数，例如 /books/:bookId -> ["bookId"]"""
    return re.findall(r':(\w+)', path)


def gin_to_openapi_path(path: str) -> str:
    """
    转换 Gin 风格路径为 OpenAPI 风格
    例如：/books/:bookId -> /books/{bookId}
    """
    return re.sub(r':(\w+)', r'{\1}', path)


def process_module(module_name: str, dry_run: bool, verbose: bool):
    """处理单个模块的路由"""
    print(f"\n{'='*60}")
    print(f"处理模块: {module_name}")
    print(f"{'='*60}")

    router_file = ROUTER_DIR / module_name / f"{module_name}_router.go"

    if not router_file.exists():
        print(f"❌ 路由文件不存在: {router_file}")
        return

    routes = parse_gin_routes(router_file)

    if not routes:
        print(f"⚠️  未找到任何路由定义")
        return

    print(f"✓ 找到 {len(routes)} 个路由定义")

    if verbose:
        for route in routes:
            path_params = extract_path_params(route["path"])
            openapi_path = gin_to_openapi_path(route["path"])
            print(f"  {route['method']:6} {route['path']:40} -> {route['handler']}")
            if path_params:
                print(f"         路径参数: {', '.join(path_params)}")
                print(f"         OpenAPI:  {openapi_path}")


if __name__ == "__main__":
    main()
