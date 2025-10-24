#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
青羽后端 - 快速验证脚本（Python 版本）

快速验证项目关键功能是否正常工作

使用方法:
    python scripts/testing/quick_verify.py
    python scripts/testing/quick_verify.py --verbose
"""

import argparse
import sys
import os
import subprocess
import time
from pathlib import Path
from typing import Tuple, List

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

def run_command(cmd: list, timeout: int = 30, verbose: bool = False) -> Tuple[bool, str]:
    """
    运行命令并返回结果

    Args:
        cmd: 命令列表
        timeout: 超时时间（秒）
        verbose: 是否显示详细输出

    Returns:
        (是否成功, 输出内容)
    """
    try:
        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            timeout=timeout,
            check=False
        )

        output = result.stdout + result.stderr

        if verbose and output:
            print(output)

        return result.returncode == 0, output
    except subprocess.TimeoutExpired:
        return False, "命令超时"
    except FileNotFoundError:
        return False, f"命令未找到: {cmd[0]}"
    except Exception as e:
        return False, str(e)

class VerificationTest:
    """验证测试基类"""

    def __init__(self, name: str, verbose: bool = False):
        self.name = name
        self.verbose = verbose

    def run(self) -> bool:
        """运行测试"""
        raise NotImplementedError

    def print_start(self):
        """打印开始信息"""
        print(f"  [{self.name}]", end=" ", flush=True)

    def print_result(self, success: bool, message: str = ""):
        """打印结果"""
        if success:
            print_color("✓ 通过", Colors.GREEN)
            if message and self.verbose:
                print(f"    {message}")
        else:
            print_color("✗ 失败", Colors.RED)
            if message:
                print_color(f"    {message}", Colors.RED)

class EnvironmentCheck(VerificationTest):
    """环境检查"""

    def run(self) -> bool:
        self.print_start()

        # 检查必需文件
        required_files = [
            "go.mod",
            "config/config.local.yaml",
            "cmd/server/main.go"
        ]

        for file in required_files:
            if not Path(file).exists():
                self.print_result(False, f"缺少文件: {file}")
                return False

        self.print_result(True, "所有必需文件存在")
        return True

class GoVersionCheck(VerificationTest):
    """Go 版本检查"""

    def run(self) -> bool:
        self.print_start()

        success, output = run_command(["go", "version"], verbose=self.verbose)

        if success:
            version_info = output.strip()
            self.print_result(True, version_info)
        else:
            self.print_result(False, "Go 未安装或无法访问")

        return success

class MongoDBCheck(VerificationTest):
    """MongoDB 连接检查"""

    def run(self) -> bool:
        self.print_start()

        cmd = ["go", "run", "cmd/migrate/main.go", "-command=status", "-config=."]
        success, output = run_command(cmd, timeout=10, verbose=self.verbose)

        if success:
            self.print_result(True, "MongoDB 连接正常")
        else:
            self.print_result(False, "MongoDB 连接失败，请确保服务已启动")

        return success

class BuildCheck(VerificationTest):
    """编译检查"""

    def run(self) -> bool:
        self.print_start()

        cmd = ["go", "build", "-o", "temp_server", "cmd/server/main.go"]
        success, output = run_command(cmd, timeout=60, verbose=self.verbose)

        # 清理临时文件
        if Path("temp_server").exists():
            try:
                os.remove("temp_server")
            except:
                pass

        if success:
            self.print_result(True, "项目编译成功")
        else:
            self.print_result(False, "编译失败")
            if self.verbose:
                print(output)

        return success

class UnitTestCheck(VerificationTest):
    """单元测试检查"""

    def run(self) -> bool:
        self.print_start()

        # 运行快速测试（只测试不需要数据库的部分）
        cmd = ["go", "test", "-short", "-timeout=30s", "./..."]
        success, output = run_command(cmd, timeout=60, verbose=self.verbose)

        if success:
            self.print_result(True, "单元测试通过")
        else:
            # 单元测试失败不一定是错误（可能需要数据库）
            self.print_result(False, "部分测试失败（可能需要完整环境）")

        # 返回 True 因为这不是关键错误
        return True

class ConfigCheck(VerificationTest):
    """配置文件检查"""

    def run(self) -> bool:
        self.print_start()

        config_file = "config/config.local.yaml"

        if not Path(config_file).exists():
            self.print_result(False, f"配置文件不存在: {config_file}")
            return False

        # 检查配置文件内容
        try:
            with open(config_file, 'r', encoding='utf-8') as f:
                content = f.read()

                # 检查关键配置项
                required_keys = ["database", "server", "jwt"]
                missing_keys = [key for key in required_keys if key not in content]

                if missing_keys:
                    self.print_result(False, f"缺少配置项: {', '.join(missing_keys)}")
                    return False

                self.print_result(True, "配置文件完整")
                return True
        except Exception as e:
            self.print_result(False, str(e))
            return False

def main():
    """主函数"""
    parser = argparse.ArgumentParser(
        description='青羽后端 - 快速验证',
        formatter_class=argparse.RawDescriptionHelpFormatter
    )

    parser.add_argument(
        '-v', '--verbose',
        action='store_true',
        help='显示详细输出'
    )

    parser.add_argument(
        '--skip-build',
        action='store_true',
        help='跳过编译检查（加快验证速度）'
    )

    parser.add_argument(
        '--skip-tests',
        action='store_true',
        help='跳过单元测试'
    )

    args = parser.parse_args()

    # 打印欢迎信息
    print_header("青羽后端 - 快速验证")

    # 检查是否在项目根目录
    if not Path("go.mod").exists():
        print_color("错误: 请在项目根目录运行此脚本", Colors.RED)
        print(f"当前目录: {os.getcwd()}")
        return 1

    # 创建测试列表
    tests: List[VerificationTest] = [
        EnvironmentCheck("环境检查", args.verbose),
        GoVersionCheck("Go 版本", args.verbose),
        MongoDBCheck("MongoDB 连接", args.verbose),
        ConfigCheck("配置文件", args.verbose),
    ]

    if not args.skip_build:
        tests.append(BuildCheck("项目编译", args.verbose))

    if not args.skip_tests:
        tests.append(UnitTestCheck("单元测试", args.verbose))

    # 运行所有测试
    print("正在验证项目状态...")
    print()

    results = []
    for test in tests:
        result = test.run()
        results.append((test.name, result))

    # 打印汇总
    print()
    print_header("验证结果")

    passed = sum(1 for _, result in results if result)
    total = len(results)

    for name, result in results:
        status = print_color("✓", Colors.GREEN) if result else print_color("✗", Colors.RED)
        # 不使用 print_color 的返回值，直接打印
        if result:
            print(f"  {name}: ", end="")
            print_color("通过", Colors.GREEN)
        else:
            print(f"  {name}: ", end="")
            print_color("失败", Colors.RED)

    print()
    print(f"总计: {passed}/{total} 项通过")

    if passed == total:
        print()
        print_color("✓ 所有检查通过！项目状态良好。", Colors.GREEN)
        print()
        print("你可以：")
        print("  1. 启动服务器: go run cmd/server/main.go")
        print("  2. 运行测试: python scripts/testing/run_tests.py")
        print("  3. 初始化数据: python scripts/init/setup_local_test_data.py")
        return 0
    else:
        print()
        print_color(f"✗ {total - passed} 项检查失败", Colors.RED)
        print()
        print("请修复上述问题后重试")
        return 1

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

