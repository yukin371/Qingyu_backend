#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
青羽后端 - 本地测试数据初始化脚本（Python 版本）

本脚本用于快速初始化本地开发环境的测试数据，包括：
1. 导入小说数据（100本）
2. 创建内测用户（20个）

使用方法:
    python scripts/init/setup_local_test_data.py
    python scripts/init/setup_local_test_data.py --novels-file data/novels_500.json
"""

import argparse
import sys
import os
import subprocess
import json
from pathlib import Path
from typing import Optional

# ANSI 颜色代码
class Colors:
    RED = '\033[0;31m'
    GREEN = '\033[0;32m'
    YELLOW = '\033[1;33m'
    BLUE = '\033[0;34m'
    CYAN = '\033[0;36m'
    NC = '\033[0m'  # No Color

def print_color(message: str, color: str = Colors.NC):
    """打印彩色文本"""
    print(f"{color}{message}{Colors.NC}")

def print_header(title: str):
    """打印标题"""
    separator = "=" * 50
    print()
    print(separator)
    print(title)
    print(separator)
    print()

def print_step(step: int, total: int, description: str):
    """打印步骤信息"""
    print(f"[{step}/{total}] {description}...")
    print()

def check_file_exists(filepath: str) -> bool:
    """检查文件是否存在"""
    return Path(filepath).exists()

def run_command(cmd: list, description: str, cwd: Optional[str] = None) -> bool:
    """
    运行命令并处理结果

    Args:
        cmd: 命令列表
        description: 命令描述
        cwd: 工作目录

    Returns:
        是否成功
    """
    try:
        result = subprocess.run(
            cmd,
            cwd=cwd,
            capture_output=False,  # 直接显示输出
            text=True,
            check=True
        )
        return True
    except subprocess.CalledProcessError as e:
        print_color(f"✗ {description} 失败", Colors.RED)
        return False
    except FileNotFoundError:
        print_color(f"✗ 命令未找到: {cmd[0]}", Colors.RED)
        return False

def check_mongodb_connection() -> bool:
    """检查 MongoDB 连接"""
    print_step(1, 3, "检查 MongoDB 连接")

    cmd = ["go", "run", "cmd/migrate/main.go", "-command=status", "-config=."]
    success = run_command(cmd, "MongoDB 连接检查")

    if not success:
        print()
        print_color("MongoDB 连接失败！", Colors.RED)
        print("请确保：")
        print("1. MongoDB 服务已启动")
        print("2. config/config.local.yaml 中的数据库配置正确")
        print()
        print("启动 MongoDB 服务：")
        print("  Windows: net start MongoDB")
        print("  Linux:   sudo systemctl start mongod")
        print("  Mac:     brew services start mongodb-community")
        return False

    return True

def import_novels(novels_file: str) -> bool:
    """导入小说数据"""
    print_step(2, 3, f"导入小说数据（{novels_file}）")

    # 检查文件是否存在
    if not check_file_exists(novels_file):
        print_color(f"✗ 数据文件不存在: {novels_file}", Colors.RED)
        print()
        print("请先运行 Python 脚本生成数据：")
        print(f"  python scripts/data/import_novels.py --max-novels 100 --output {novels_file}")
        return False

    # 显示文件信息
    try:
        with open(novels_file, 'r', encoding='utf-8') as f:
            data = json.load(f)
            metadata = data.get('metadata', {})
            print(f"  数据来源: {metadata.get('source', 'Unknown')}")
            print(f"  小说数量: {metadata.get('total_novels', 0)}")
            print(f"  章节数量: {metadata.get('total_chapters', 0)}")
            print()
    except Exception as e:
        print_color(f"  警告: 无法读取文件元数据 - {e}", Colors.YELLOW)

    # 导入数据
    cmd = [
        "go", "run", "cmd/migrate/main.go",
        f"-command=import-novels",
        f"-file={novels_file}",
        "-config=."
    ]

    return run_command(cmd, "小说数据导入")

def create_beta_users() -> bool:
    """创建内测用户"""
    print_step(3, 3, "创建内测用户")

    cmd = ["go", "run", "cmd/create_beta_users/main.go"]
    return run_command(cmd, "内测用户创建")

def print_summary():
    """打印完成摘要"""
    print()
    print_header("测试数据初始化完成！")

    print_color("你现在可以：", Colors.GREEN)
    print()
    print("1. 启动服务器：")
    print_color("   go run cmd/server/main.go", Colors.CYAN)
    print()
    print("2. 使用以下测试账号登录：")
    print()
    print_color("   管理员账号：", Colors.YELLOW)
    print("     用户名: admin")
    print("     密码: Admin@123456")
    print()
    print_color("   VIP作家账号：", Colors.YELLOW)
    print("     用户名: vip_writer01")
    print("     密码: Vip@123456")
    print()
    print_color("   普通作家：", Colors.YELLOW)
    print("     用户名: writer_xuanhuan")
    print("     密码: Writer@123456")
    print()
    print_color("   普通读者：", Colors.YELLOW)
    print("     用户名: reader01")
    print("     密码: Reader@123456")
    print()
    print("详细账号列表请查看上方输出")
    print("=" * 50)
    print()

def main():
    """主函数"""
    parser = argparse.ArgumentParser(
        description='青羽后端 - 本地测试数据初始化',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 使用默认数据文件
  python scripts/init/setup_local_test_data.py

  # 使用自定义数据文件
  python scripts/init/setup_local_test_data.py --novels-file data/novels_500.json

  # 仅导入小说数据
  python scripts/init/setup_local_test_data.py --skip-users

  # 仅创建用户
  python scripts/init/setup_local_test_data.py --skip-novels
        """
    )

    parser.add_argument(
        '--novels-file',
        type=str,
        default='data/novels_100.json',
        help='小说数据文件路径（默认: data/novels_100.json）'
    )

    parser.add_argument(
        '--skip-novels',
        action='store_true',
        help='跳过小说数据导入'
    )

    parser.add_argument(
        '--skip-users',
        action='store_true',
        help='跳过用户创建'
    )

    args = parser.parse_args()

    # 打印欢迎信息
    print_header("青羽写作系统 - 本地测试数据初始化")

    # 检查是否在项目根目录
    if not check_file_exists("go.mod"):
        print_color("错误: 请在项目根目录运行此脚本", Colors.RED)
        print(f"当前目录: {os.getcwd()}")
        return 1

    # 步骤 1: 检查 MongoDB 连接
    if not check_mongodb_connection():
        return 1

    # 步骤 2: 导入小说数据
    if not args.skip_novels:
        if not import_novels(args.novels_file):
            return 1
    else:
        print_color("跳过小说数据导入", Colors.YELLOW)

    # 步骤 3: 创建内测用户
    if not args.skip_users:
        if not create_beta_users():
            return 1
    else:
        print_color("跳过用户创建", Colors.YELLOW)

    # 打印完成信息
    print_summary()

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

