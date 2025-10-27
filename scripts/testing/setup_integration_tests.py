#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
青羽后端 - 集成测试准备脚本
用途：准备集成测试所需的所有数据，避免编码问题
"""

import os
import sys
import json
import subprocess
import time
import requests
from pathlib import Path

# 颜色输出
class Colors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

def print_header(message):
    print(f"\n{Colors.HEADER}{'='*60}{Colors.ENDC}")
    print(f"{Colors.HEADER}{message}{Colors.ENDC}")
    print(f"{Colors.HEADER}{'='*60}{Colors.ENDC}\n")

def print_success(message):
    print(f"{Colors.OKGREEN}✓ {message}{Colors.ENDC}")

def print_error(message):
    print(f"{Colors.FAIL}✗ {message}{Colors.ENDC}")

def print_info(message):
    print(f"{Colors.OKCYAN}➜ {message}{Colors.ENDC}")

def print_warning(message):
    print(f"{Colors.WARNING}⚠ {message}{Colors.ENDC}")

def run_command(cmd, shell=True, check=True):
    """运行命令并返回结果"""
    try:
        result = subprocess.run(
            cmd,
            shell=shell,
            check=check,
            capture_output=True,
            text=True,
            encoding='utf-8'
        )
        return result.returncode == 0, result.stdout, result.stderr
    except subprocess.CalledProcessError as e:
        return False, e.stdout, e.stderr
    except Exception as e:
        return False, "", str(e)

def check_mongodb():
    """检查MongoDB是否运行"""
    print_info("检查 MongoDB 服务...")
    success, stdout, stderr = run_command("mongosh --eval \"db.version()\" --quiet", check=False)
    if success:
        print_success("MongoDB 服务正常运行")
        return True
    else:
        print_error("MongoDB 服务未运行")
        print_warning("请先启动 MongoDB 服务")
        return False

def cleanup_test_data():
    """清理测试数据"""
    print_header("第1步: 清理测试数据")
    
    if not check_mongodb():
        return False
    
    confirm = input("确认清理测试数据？这将删除所有测试数据 (yes/no): ")
    if confirm.lower() != 'yes':
        print_warning("操作已取消")
        return False
    
    # MongoDB 清理命令
    cleanup_script = """
    use qingyu_test;
    
    print('清理书籍数据...');
    db.books.deleteMany({});
    
    print('清理章节数据...');
    db.chapters.deleteMany({});
    
    print('清理用户数据（保留系统用户）...');
    db.users.deleteMany({ role: { $ne: 'system' } });
    
    print('清理榜单数据...');
    db.ranking_items.deleteMany({});
    
    print('清理阅读进度...');
    db.reading_progress.deleteMany({});
    
    print('清理书签和笔记...');
    db.annotations.deleteMany({});
    
    print('清理收藏记录...');
    db.user_collections.deleteMany({});
    
    print('清理评论数据...');
    db.comments.deleteMany({});
    
    print('清理写作项目（测试项目）...');
    db.projects.deleteMany({ title: /测试/ });
    
    print('');
    print('清理后数据统计：');
    print('  - 书籍数量: ' + db.books.countDocuments());
    print('  - 章节数量: ' + db.chapters.countDocuments());
    print('  - 用户数量: ' + db.users.countDocuments());
    """
    
    # 写入临时文件
    temp_file = "temp_cleanup.js"
    with open(temp_file, 'w', encoding='utf-8') as f:
        f.write(cleanup_script)
    
    # 执行清理
    print_info("执行数据清理...")
    success, stdout, stderr = run_command(f'mongosh --quiet < {temp_file}', check=False)
    
    # 删除临时文件
    if os.path.exists(temp_file):
        os.remove(temp_file)
    
    if success or "countDocuments" in stdout:
        print_success("数据清理完成")
        print(stdout)
        return True
    else:
        print_error("数据清理失败")
        print(stderr)
        return False

def import_test_users():
    """导入测试用户"""
    print_header("第2步: 导入测试用户")
    
    print_info("运行用户导入脚本...")
    success, stdout, stderr = run_command(
        "go run scripts/testing/import_test_users.go",
        check=False
    )
    
    if success or "成功创建" in stdout:
        print_success("测试用户导入成功")
        # 提取并显示关键信息
        for line in stdout.split('\n'):
            if '成功创建' in line or '已存在' in line or '用户名:' in line:
                print(f"  {line.strip()}")
        return True
    else:
        print_error("测试用户导入失败")
        print(stderr)
        return False

def import_novels():
    """导入小说数据"""
    print_header("第3步: 导入小说数据")
    
    # 检查数据文件
    data_file = "data/novels_100.json"
    if not os.path.exists(data_file):
        print_error(f"数据文件不存在: {data_file}")
        print_warning("请先运行 Python 脚本生成数据文件")
        return False
    
    print_success(f"找到数据文件: {data_file}")
    
    # 验证数据（试运行）
    print_info("验证数据格式（试运行）...")
    success, stdout, stderr = run_command(
        f"go run cmd/migrate/main.go -command=import-novels -file={data_file} -dry-run=true",
        check=False
    )
    
    if not success and "验证通过" not in stdout:
        print_error("数据验证失败")
        print(stderr)
        return False
    
    print_success("数据验证通过")
    
    # 正式导入
    print_info("正式导入数据...")
    success, stdout, stderr = run_command(
        f"go run cmd/migrate/main.go -command=import-novels -file={data_file}",
        check=False
    )
    
    if success or "导入完成" in stdout or "成功" in stdout:
        print_success("小说数据导入成功")
        # 显示统计信息
        for line in stdout.split('\n'):
            if '书籍' in line or '章节' in line or '成功' in line:
                print(f"  {line.strip()}")
        return True
    else:
        print_error("小说数据导入失败")
        print(stderr)
        return False

def check_server():
    """检查服务器是否运行"""
    print_header("检查服务器状态")
    
    try:
        response = requests.get("http://localhost:8080/api/v1/system/health", timeout=2)
        if response.status_code == 200:
            print_success("服务器正在运行")
            return True
    except:
        pass
    
    print_warning("服务器未运行")
    print_info("请在另一个终端运行: go run cmd/server/main.go")
    
    start_server = input("是否在后台启动服务器？(yes/no): ")
    if start_server.lower() == 'yes':
        if sys.platform == 'win32':
            subprocess.Popen(
                "start /B go run cmd/server/main.go",
                shell=True,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL
            )
        else:
            subprocess.Popen(
                "go run cmd/server/main.go > server.log 2>&1 &",
                shell=True,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL
            )
        
        print_info("等待服务器启动...")
        for i in range(10):
            time.sleep(1)
            try:
                response = requests.get("http://localhost:8080/api/v1/system/health", timeout=1)
                if response.status_code == 200:
                    print_success("服务器启动成功")
                    return True
            except:
                print(f"  等待中... {i+1}/10")
        
        print_error("服务器启动失败")
        return False
    
    return False

def run_tests():
    """运行集成测试"""
    print_header("执行集成测试")
    
    print("可用的测试场景:")
    print("  1. 书城流程测试")
    print("  2. 搜索功能测试")
    print("  3. 阅读流程测试")
    print("  4. AI生成测试")
    print("  5. 认证流程测试")
    print("  6. 写作流程测试")
    print("  7. 互动功能测试")
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
        print_error("无效的选择")
        return False
    
    test_file = test_files[choice]
    print_info(f"运行测试: {test_file}")
    
    # 根据选择构建命令
    if test_file == 'all':
        cmd = "go test -v ./test/integration/ -run Scenario"
    else:
        cmd = f"go test -v ./test/integration/{test_file}"
    
    success, stdout, stderr = run_command(cmd, check=False)
    
    print(stdout)
    if stderr:
        print(stderr)
    
    if success or "PASS" in stdout:
        print_success("测试执行完成")
        return True
    else:
        print_warning("部分测试可能失败，请查看详细输出")
        return False

def show_test_accounts():
    """显示测试账号信息"""
    print_header("测试账号清单")
    
    accounts = [
        {"角色": "管理员", "邮箱": "admin@qingyu.com", "密码": "Admin@123456"},
        {"角色": "VIP用户", "邮箱": "vip01@qingyu.com", "密码": "Vip@123456"},
        {"角色": "VIP用户", "邮箱": "vip02@qingyu.com", "密码": "Vip@123456"},
        {"角色": "普通用户", "邮箱": "test01@qingyu.com", "密码": "Test@123456"},
        {"角色": "普通用户", "邮箱": "test02@qingyu.com", "密码": "Test@123456"},
        {"角色": "普通用户", "邮箱": "test03@qingyu.com", "密码": "Test@123456"},
        {"角色": "普通用户", "邮箱": "test04@qingyu.com", "密码": "Test@123456"},
        {"角色": "普通用户", "邮箱": "test05@qingyu.com", "密码": "Test@123456"},
    ]
    
    print(f"{'角色':<12} {'邮箱':<25} {'密码':<20}")
    print("-" * 60)
    for account in accounts:
        print(f"{account['角色']:<12} {account['邮箱']:<25} {account['密码']:<20}")
    print()

def main():
    """主函数"""
    print_header("青羽后端 - 集成测试准备工具")
    
    # 切换到项目根目录
    script_dir = Path(__file__).parent
    project_root = script_dir.parent.parent
    os.chdir(project_root)
    
    print(f"项目根目录: {os.getcwd()}\n")
    
    # 询问是否准备数据
    prepare_data = input("是否需要准备测试数据？(首次运行选yes) [yes/no]: ")
    
    if prepare_data.lower() == 'yes':
        # 1. 清理数据
        if not cleanup_test_data():
            print_error("数据清理失败，终止执行")
            return 1
        
        # 2. 导入用户
        if not import_test_users():
            print_error("用户导入失败，终止执行")
            return 1
        
        # 3. 导入小说
        if not import_novels():
            print_error("小说导入失败，终止执行")
            return 1
        
        # 显示测试账号
        show_test_accounts()
    
    # 检查服务器
    if not check_server():
        print_warning("服务器未运行，部分测试可能失败")
    
    # 运行测试
    run_test = input("\n是否运行集成测试？[yes/no]: ")
    if run_test.lower() == 'yes':
        run_tests()
    
    print_header("完成")
    print_success("集成测试准备完成！")
    print()
    print("提示:")
    print("  - 查看使用指南: doc/testing/集成测试使用指南.md")
    print("  - 手动运行测试: go test -v ./test/integration/scenario_*.go")
    print("  - 停止服务器: Ctrl+C 或手动关闭")
    print()
    
    return 0

if __name__ == "__main__":
    try:
        sys.exit(main())
    except KeyboardInterrupt:
        print("\n\n操作已取消")
        sys.exit(1)
    except Exception as e:
        print_error(f"发生错误: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

