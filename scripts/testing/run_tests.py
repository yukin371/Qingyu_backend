#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
青羽后端 - 集成测试运行脚本（Python版本）
支持跨平台，避免编码问题
"""

import subprocess
import sys
import os
import time
import io

# 设置标准输出编码为UTF-8（解决Windows PowerShell编码问题）
if sys.platform == 'win32':
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

try:
    import requests
except ImportError:
    print("错误: 缺少 requests 库")
    print("请运行: pip install requests")
    sys.exit(1)

def print_header(text):
    print("\n" + "=" * 60)
    print(text)
    print("=" * 60 + "\n")

def check_server():
    """检查服务器是否运行"""
    try:
        response = requests.get("http://localhost:8080/api/v1/system/health", timeout=3)
        return response.status_code == 200
    except requests.exceptions.ConnectionError:
        return False
    except requests.exceptions.Timeout:
        return False
    except Exception as e:
        print(f"检查服务器时出错: {e}")
        return False

def start_server():
    """启动服务器"""
    print("正在后台启动服务器...")

    try:
        # 打开日志文件
        log_file = open("server_test.log", 'w', encoding='utf-8')

        if sys.platform == 'win32':
            # Windows - 使用 CREATE_NEW_PROCESS_GROUP 避免继承控制台
            subprocess.Popen(
                ["go", "run", "cmd/server/main.go"],
                stdout=log_file,
                stderr=subprocess.STDOUT,
                creationflags=subprocess.CREATE_NEW_PROCESS_GROUP
            )
        else:
            # Linux/Mac
            subprocess.Popen(
                ["go", "run", "cmd/server/main.go"],
                stdout=log_file,
                stderr=subprocess.STDOUT
            )

        # 等待服务器启动（增加到15秒）
        print("等待服务器启动...")
        for i in range(15):
            time.sleep(1)
            if check_server():
                print("✓ 服务器启动成功")
                return True
            print(f"  等待中... {i+1}/15")

        print("✗ 服务器启动失败")
        print("请查看 server_test.log 了解详情")
        return False
    except Exception as e:
        print(f"✗ 启动服务器时出错: {e}")
        return False

def run_test(test_file):
    """运行测试"""
    print(f"\n运行测试: {test_file}\n")

    try:
        # 如果是运行所有测试，使用包路径而不是通配符
        if test_file == 'all':
            cmd = ["go", "test", "-v", "./test/integration/", "-run", "Scenario"]
        else:
            cmd = ["go", "test", "-v", f"./test/integration/{test_file}"]

        result = subprocess.run(
            cmd,
            encoding='utf-8',
            text=True,
            capture_output=False
        )

        return result.returncode == 0
    except Exception as e:
        print(f"✗ 运行测试时出错: {e}")
        return False

def main():
    print_header("青羽后端 - 集成测试运行工具")

    # 检查命令行参数
    import argparse
    parser = argparse.ArgumentParser(description='青羽后端集成测试运行工具')
    parser.add_argument('--test', type=str, choices=['1', '2', '3', '4', '5', '6', '7', '8', 'all'],
                        help='直接运行指定的测试（1-8 或 all）')
    parser.add_argument('--no-start', action='store_true',
                        help='不自动启动服务器')
    args = parser.parse_args()

    # 检查 Go 环境
    try:
        result = subprocess.run(["go", "version"], capture_output=True, text=True)
        if result.returncode != 0:
            print("✗ 错误: 未找到 Go 环境")
            print("请确保已安装 Go 并添加到 PATH")
            return 1
        print(f"✓ Go 环境: {result.stdout.strip()}")
    except Exception as e:
        print(f"✗ 错误: 无法检查 Go 环境: {e}")
        return 1

    # 检查服务器
    print("\n检查服务器状态...")
    if not check_server():
        print("⚠ 服务器未运行")
        if not args.no_start:
            if args.test:
                # 非交互模式，自动启动服务器
                print("自动启动服务器...")
                if not start_server():
                    print("\n✗ 无法启动服务器，请手动启动")
                    print("命令: go run cmd/server/main.go")
                    return 1
            else:
                # 交互模式，询问用户
                choice = input("是否启动服务器？(yes/no): ")
                if choice.lower() in ['yes', 'y']:
                    if not start_server():
                        print("\n✗ 无法启动服务器，请手动启动")
                        print("命令: go run cmd/server/main.go")
                        return 1
                else:
                    print("\n请手动启动服务器后再运行测试")
                    print("命令: go run cmd/server/main.go")
                    return 1
        else:
            print("\n✗ 服务器未运行且指定了 --no-start")
            print("请手动启动服务器: go run cmd/server/main.go")
            return 1
    else:
        print("✓ 服务器正在运行")

    # 选择测试
    if args.test:
        # 命令行指定了测试
        if args.test == 'all':
            choice = '8'
        else:
            choice = args.test
    else:
        # 交互式选择
        print("\n可用的测试场景:")
        print("  1. 书城流程测试 (scenario_bookstore_test.go)")
        print("  2. 搜索功能测试 (scenario_search_test.go)")
        print("  3. 阅读流程测试 (scenario_reading_test.go)")
        print("  4. AI生成测试 (scenario_ai_generation_test.go)")
        print("  5. 认证流程测试 (scenario_auth_test.go)")
        print("  6. 写作流程测试 (scenario_writing_test.go)")
        print("  7. 互动功能测试 (scenario_interaction_test.go)")
        print("  8. 全部测试")
        print()

        choice = input("请选择要执行的测试 (1-8): ")

    test_files = {
        '1': 'scenario_bookstore_test.go',
        '2': 'scenario_search_test.go',
        '3': 'scenario_reading_test.go',
        '4': 'scenario_ai_generation_test.go',
        '5': 'scenario_auth_test.go',
        '6': 'scenario_writing_test.go',
        '7': 'scenario_interaction_test.go',
        '8': 'all'  # 特殊标记，表示运行所有测试
    }

    if choice not in test_files:
        print("无效的选择")
        return 1

    test_file = test_files[choice]

    print_header("执行集成测试")

    success = run_test(test_file)

    print_header("测试完成")

    if success:
        print("✓ 测试执行成功")
    else:
        print("⚠ 部分测试可能失败，请查看详细输出")

    # 询问是否查看服务器日志
    if os.path.exists("server_test.log"):
        view_log = input("\n是否查看服务器日志？(yes/no): ")
        if view_log.lower() == 'yes':
            with open("server_test.log", 'r', encoding='utf-8') as f:
                print("\n" + "=" * 60)
                print("服务器日志")
                print("=" * 60)
                print(f.read())

    print("\n提示:")
    print("- 测试结果已显示在上方")
    print("- 服务器日志: server_test.log")
    print("- 如需停止服务器，请手动关闭或使用 Ctrl+C")
    print()

    return 0 if success else 1

if __name__ == "__main__":
    try:
        sys.exit(main())
    except KeyboardInterrupt:
        print("\n\n操作已取消")
        sys.exit(1)
