#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
青羽后端 - 小说导入测试脚本（Python 版本）

测试小说数据导入功能

使用方法:
    python scripts/data/test_novel_import.py
    python scripts/data/test_novel_import.py --file data/novels_500.json
"""

import argparse
import sys
import os
import subprocess
import json
from pathlib import Path
from typing import Dict, Any

# ANSI 颜色代码
class Colors:
    RED = '\033[0;31m'
    GREEN = '\033[0;32m'
    YELLOW = '\033[1;33m'
    BLUE = '\033[0;34m'
    CYAN = '\033[0;36m'
    NC = '\033[0m'

def print_color(message: str, color: str = Colors.NC):
    """打印彩色文本"""
    print(f"{color}{message}{Colors.NC}")

def print_header(title: str):
    """打印标题"""
    separator = "=" * 60
    print()
    print(separator)
    print(title)
    print(separator)
    print()

def validate_json_file(filepath: str) -> tuple[bool, Dict[str, Any]]:
    """
    验证 JSON 文件

    Returns:
        (是否有效, 数据或错误信息)
    """
    if not Path(filepath).exists():
        return False, {"error": f"文件不存在: {filepath}"}

    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            data = json.load(f)

        # 验证数据结构
        if 'metadata' not in data:
            return False, {"error": "缺少 metadata 字段"}

        if 'novels' not in data:
            return False, {"error": "缺少 novels 字段"}

        if not isinstance(data['novels'], list):
            return False, {"error": "novels 字段应该是数组"}

        return True, data

    except json.JSONDecodeError as e:
        return False, {"error": f"JSON 格式错误: {e}"}
    except Exception as e:
        return False, {"error": str(e)}

def print_data_info(data: Dict[str, Any]):
    """打印数据信息"""
    metadata = data.get('metadata', {})
    novels = data.get('novels', [])

    print_header("数据文件信息")

    print(f"来源: {metadata.get('source', 'Unknown')}")
    print(f"小说数量: {metadata.get('total_novels', len(novels))}")
    print(f"章节数量: {metadata.get('total_chapters', 0)}")
    print(f"生成时间: {metadata.get('generated_at', 'Unknown')}")
    print(f"章节大小: {metadata.get('chapter_size', 'Unknown')} 字")
    print()

    # 显示前几本小说信息
    if novels:
        print("前 3 本小说:")
        for i, novel in enumerate(novels[:3], 1):
            print(f"  {i}. {novel.get('title', 'Unknown')}")
            print(f"     作者: {novel.get('author', 'Unknown')}")
            print(f"     分类: {novel.get('category', 'Unknown')}")
            print(f"     字数: {novel.get('word_count', 0)}")
            print(f"     章节: {novel.get('chapter_count', 0)}")
            print()

def dry_run_import(filepath: str) -> bool:
    """
    试运行导入（不实际写入数据库）

    Returns:
        是否成功
    """
    print_header("试运行导入")

    cmd = [
        "go", "run", "cmd/migrate/main.go",
        "-command=import-novels",
        f"-file={filepath}",
        "-dry-run=true",
        "-config=."
    ]

    try:
        result = subprocess.run(cmd, check=False)

        if result.returncode == 0:
            print_color("✓ 试运行成功", Colors.GREEN)
            return True
        else:
            print_color("✗ 试运行失败", Colors.RED)
            return False
    except Exception as e:
        print_color(f"✗ 运行失败: {e}", Colors.RED)
        return False

def import_novels(filepath: str) -> bool:
    """
    正式导入小说数据

    Returns:
        是否成功
    """
    print_header("正式导入")

    print_color("警告: 这将向数据库插入数据", Colors.YELLOW)
    print("确认导入？(yes/no): ", end="", flush=True)

    try:
        confirmation = input().strip().lower()
    except EOFError:
        # 非交互式环境
        confirmation = "no"

    if confirmation != "yes":
        print_color("已取消导入", Colors.YELLOW)
        return False

    cmd = [
        "go", "run", "cmd/migrate/main.go",
        "-command=import-novels",
        f"-file={filepath}",
        "-config=."
    ]

    try:
        result = subprocess.run(cmd, check=False)

        if result.returncode == 0:
            print()
            print_color("✓ 导入成功", Colors.GREEN)
            return True
        else:
            print()
            print_color("✗ 导入失败", Colors.RED)
            return False
    except Exception as e:
        print()
        print_color(f"✗ 运行失败: {e}", Colors.RED)
        return False

def main():
    """主函数"""
    parser = argparse.ArgumentParser(
        description='青羽后端 - 小说导入测试',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 测试默认数据文件
  python scripts/data/test_novel_import.py

  # 测试自定义数据文件
  python scripts/data/test_novel_import.py --file data/novels_500.json

  # 仅验证文件，不导入
  python scripts/data/test_novel_import.py --validate-only

  # 直接导入，不试运行
  python scripts/data/test_novel_import.py --skip-dry-run
        """
    )

    parser.add_argument(
        '--file',
        type=str,
        default='data/novels_100.json',
        help='小说数据文件路径（默认: data/novels_100.json）'
    )

    parser.add_argument(
        '--validate-only',
        action='store_true',
        help='仅验证文件，不导入'
    )

    parser.add_argument(
        '--skip-dry-run',
        action='store_true',
        help='跳过试运行，直接导入'
    )

    parser.add_argument(
        '--auto-confirm',
        action='store_true',
        help='自动确认导入（用于自动化测试）'
    )

    args = parser.parse_args()

    # 打印欢迎信息
    print_header("青羽后端 - 小说导入测试")

    # 检查是否在项目根目录
    if not Path("go.mod").exists():
        print_color("错误: 请在项目根目录运行此脚本", Colors.RED)
        print(f"当前目录: {os.getcwd()}")
        return 1

    # 步骤 1: 验证 JSON 文件
    print(f"验证文件: {args.file}")
    valid, data_or_error = validate_json_file(args.file)

    if not valid:
        print_color(f"✗ 文件验证失败: {data_or_error.get('error')}", Colors.RED)
        return 1

    print_color("✓ 文件验证通过", Colors.GREEN)
    print()

    # 步骤 2: 显示数据信息
    print_data_info(data_or_error)

    # 如果只验证，到这里结束
    if args.validate_only:
        print_color("验证完成（仅验证模式）", Colors.GREEN)
        return 0

    # 步骤 3: 试运行导入
    if not args.skip_dry_run:
        if not dry_run_import(args.file):
            return 1

    # 步骤 4: 正式导入
    if not import_novels(args.file):
        return 1

    # 打印完成信息
    print()
    print_header("导入完成")
    print("你可以：")
    print("  1. 查看数据库中的数据")
    print("  2. 启动服务器测试: go run cmd/server/main.go")
    print("  3. 测试 API: python scripts/testing/test_reading_features.py")

    return 0

if __name__ == '__main__':
    try:
        sys.exit(main())
    except KeyboardInterrupt:
        print()
        print_color("用户中断操作", Colors.YELLOW)
        sys.exit(130)
    except Exception as e:
        print()
        print_color(f"错误: {e}", Colors.RED)
        import traceback
        traceback.print_exc()
        sys.exit(1)

