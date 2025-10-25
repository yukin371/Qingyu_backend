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
import requests

def print_header(text):
    print("\n" + "=" * 60)
    print(text)
    print("=" * 60 + "\n")

def check_server():
    """检查服务器是否运行"""
    try:
        response = requests.get("http://localhost:8080/api/v1/system/health", timeout=2)
        return response.status_code == 200
    except:
        return False

def start_server():
    """启动服务器"""
    print("正在后台启动服务器...")
    
    if sys.platform == 'win32':
        # Windows
        subprocess.Popen(
            "start /B go run cmd/server/main.go > server_test.log 2>&1",
            shell=True
        )
    else:
        # Linux/Mac
        subprocess.Popen(
            "go run cmd/server/main.go > server_test.log 2>&1 &",
            shell=True
        )
    
    # 等待服务器启动
    print("等待服务器启动...")
    for i in range(10):
        time.sleep(1)
        if check_server():
            print("✓ 服务器启动成功")
            return True
        print(f"  等待中... {i+1}/10")
    
    print("✗ 服务器启动失败")
    return False

def run_test(test_file):
    """运行测试"""
    print(f"\n运行测试: {test_file}\n")
    
    # 如果是运行所有测试，使用包路径而不是通配符
    if test_file == 'all':
        cmd = "go test -v ./test/integration/ -run Scenario"
    else:
        cmd = f"go test -v ./test/integration/{test_file}"
    
    result = subprocess.run(
        cmd,
        shell=True,
        encoding='utf-8',
        text=True
    )
    
    return result.returncode == 0

def main():
    print_header("青羽后端 - 集成测试运行工具")
    
    # 检查服务器
    if not check_server():
        print("⚠ 服务器未运行")
        choice = input("是否启动服务器？(yes/no): ")
        if choice.lower() == 'yes':
            if not start_server():
                print("无法启动服务器，请手动启动")
                print("命令: go run cmd/server/main.go")
                return 1
        else:
            print("请手动启动服务器后再运行测试")
            return 1
    else:
        print("✓ 服务器正在运行")
    
    # 选择测试
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
