#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
青羽后端 - 测试运行脚本（Python 版本）

运行项目测试套件

使用方法:
    python scripts/testing/run_tests.py
    python scripts/testing/run_tests.py --coverage
    python scripts/testing/run_tests.py --package=service
"""

import argparse
import sys
import os
import subprocess
import time
from pathlib import Path
from datetime import datetime

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

def run_go_test(package: str = "./...", coverage: bool = False, verbose: bool = False, timeout: str = "10m") -> int:
    """
    运行 Go 测试

    Args:
        package: 测试包路径
        coverage: 是否生成覆盖率报告
        verbose: 是否显示详细输出
        timeout: 超时时间

    Returns:
        退出代码
    """
    cmd = ["go", "test"]

    # 添加包路径
    cmd.append(package)

    # 添加选项
    if verbose:
        cmd.append("-v")

    cmd.extend(["-timeout", timeout])

    # 添加覆盖率选项
    if coverage:
        coverage_file = "coverage"
        cmd.extend(["-coverprofile", coverage_file])

    # 运行测试
    try:
        result = subprocess.run(cmd, check=False)
        return result.returncode
    except KeyboardInterrupt:
        print_color("\n测试被用户中断", Colors.YELLOW)
        return 130
    except Exception as e:
        print_color(f"运行测试失败: {e}", Colors.RED)
        return 1

def generate_coverage_report(output_format: str = "html"):
    """生成覆盖率报告"""
    coverage_file = "coverage"

    if not Path(coverage_file).exists():
        print_color("覆盖率文件不存在", Colors.YELLOW)
        return

    print()
    print_header("生成覆盖率报告")

    if output_format == "html":
        # 生成 HTML 报告
        output_file = "coverage.html"
        cmd = ["go", "tool", "cover", f"-html={coverage_file}", f"-o", output_file]

        try:
            subprocess.run(cmd, check=True)
            print_color(f"✓ HTML 覆盖率报告已生成: {output_file}", Colors.GREEN)

            # 尝试在浏览器中打开
            try:
                import webbrowser
                abs_path = Path(output_file).absolute()
                webbrowser.open(f"file://{abs_path}")
                print_color("✓ 已在浏览器中打开报告", Colors.GREEN)
            except:
                pass
        except Exception as e:
            print_color(f"生成 HTML 报告失败: {e}", Colors.RED)

    # 显示覆盖率摘要
    cmd = ["go", "tool", "cover", f"-func={coverage_file}"]
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, check=True)
        print()
        print("覆盖率摘要:")
        print(result.stdout)
    except Exception as e:
        print_color(f"显示覆盖率摘要失败: {e}", Colors.RED)

def main():
    """主函数"""
    parser = argparse.ArgumentParser(
        description='青羽后端 - 测试运行脚本',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 运行所有测试
  python scripts/testing/run_tests.py

  # 运行特定包的测试
  python scripts/testing/run_tests.py --package=./service/...
  python scripts/testing/run_tests.py --package=./api/v1/...

  # 生成覆盖率报告
  python scripts/testing/run_tests.py --coverage

  # 详细输出
  python scripts/testing/run_tests.py -v

  # 快速测试（跳过需要数据库的测试）
  python scripts/testing/run_tests.py --short
        """
    )

    parser.add_argument(
        '--package',
        type=str,
        default='./...',
        help='要测试的包路径（默认: ./...，表示所有包）'
    )

    parser.add_argument(
        '--coverage',
        action='store_true',
        help='生成覆盖率报告'
    )

    parser.add_argument(
        '-v', '--verbose',
        action='store_true',
        help='显示详细输出'
    )

    parser.add_argument(
        '--short',
        action='store_true',
        help='运行快速测试（跳过长时间运行的测试）'
    )

    parser.add_argument(
        '--timeout',
        type=str,
        default='10m',
        help='测试超时时间（默认: 10m）'
    )

    args = parser.parse_args()

    # 打印欢迎信息
    print_header("青羽后端 - 测试运行器")

    # 检查是否在项目根目录
    if not Path("go.mod").exists():
        print_color("错误: 请在项目根目录运行此脚本", Colors.RED)
        print(f"当前目录: {os.getcwd()}")
        return 1

    # 显示配置
    print("测试配置:")
    print(f"  包路径: {args.package}")
    print(f"  覆盖率: {'是' if args.coverage else '否'}")
    print(f"  详细输出: {'是' if args.verbose else '否'}")
    print(f"  快速模式: {'是' if args.short else '否'}")
    print(f"  超时时间: {args.timeout}")
    print()

    # 记录开始时间
    start_time = time.time()

    # 运行测试
    print_header("运行测试")

    exit_code = run_go_test(
        package=args.package,
        coverage=args.coverage,
        verbose=args.verbose,
        timeout=args.timeout
    )

    # 记录结束时间
    elapsed_time = time.time() - start_time

    # 生成覆盖率报告
    if args.coverage and exit_code == 0:
        generate_coverage_report()

    # 打印结果
    print()
    print_header("测试结果")

    print(f"运行时间: {elapsed_time:.2f} 秒")

    if exit_code == 0:
        print_color("✓ 所有测试通过", Colors.GREEN)
        return 0
    else:
        print_color("✗ 测试失败", Colors.RED)
        return exit_code

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

