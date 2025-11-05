#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
青羽阅读端功能自动化测试脚本
"""

import os
import sys
import json
import time
from datetime import datetime
from typing import Dict, Any, Tuple, Optional

try:
    import requests
except ImportError:
    print("错误: 请先安装 requests 库")
    print("运行命令: pip install requests")
    sys.exit(1)

# 配置
BASE_URL = os.getenv("BASE_URL", "http://localhost:8080")
API_PREFIX = "/api/v1"
OUTPUT_DIR = "test_results"

# 颜色输出（Windows 支持）
try:
    from colorama import init, Fore, Style
    init(autoreset=True)
    HAS_COLOR = True
except ImportError:
    HAS_COLOR = False
    class Fore:
        GREEN = RED = YELLOW = CYAN = ""
    class Style:
        RESET_ALL = ""

# 测试统计
test_stats = {
    "total": 0,
    "passed": 0,
    "failed": 0,
    "start_time": None,
    "end_time": None
}


def print_header(text: str):
    """打印标题"""
    print("\n" + "=" * 50)
    print(text)
    print("=" * 50)


def print_result(name: str, success: bool, response_time: float = 0):
    """打印测试结果"""
    test_stats["total"] += 1
    if success:
        test_stats["passed"] += 1
        status = f"{Fore.GREEN}[OK]{Style.RESET_ALL}" if HAS_COLOR else "[OK]"
        print(f"{status} {name} ({response_time:.0f}ms)")
    else:
        test_stats["failed"] += 1
        status = f"{Fore.RED}[FAIL]{Style.RESET_ALL}" if HAS_COLOR else "[FAIL]"
        print(f"{status} {name}")


def test_api(name: str, url: str, method: str = "GET",
             data: Optional[Dict] = None, save_response: bool = True) -> Tuple[bool, Optional[Dict], float]:
    """
    测试 API 接口

    返回: (是否成功, 响应数据, 响应时间ms)
    """
    full_url = f"{BASE_URL}{API_PREFIX}{url}"

    try:
        start_time = time.time()

        if method == "GET":
            response = requests.get(full_url, timeout=10)
        elif method == "POST":
            response = requests.post(full_url, json=data, timeout=10)
        else:
            response = requests.request(method, full_url, json=data, timeout=10)

        response_time = (time.time() - start_time) * 1000  # 转换为毫秒

        # 检查状态码
        if response.status_code >= 200 and response.status_code < 300:
            try:
                response_data = response.json()

                # 保存响应
                if save_response:
                    save_file = os.path.join(OUTPUT_DIR, f"{name.replace(' ', '_')}.json")
                    with open(save_file, 'w', encoding='utf-8') as f:
                        json.dump(response_data, f, ensure_ascii=False, indent=2)

                return True, response_data, response_time
            except json.JSONDecodeError:
                return False, None, response_time
        else:
            return False, None, response_time

    except requests.RequestException as e:
        print(f"  错误: {e}")
        return False, None, 0


def extract_book_id(response: Optional[Dict]) -> Optional[str]:
    """从响应中提取第一本书的 ID"""
    if not response or not response.get("data"):
        return None

    data = response["data"]
    if isinstance(data, list) and len(data) > 0:
        return data[0].get("id")
    return None


def extract_chapter_id(response: Optional[Dict]) -> Optional[str]:
    """从响应中提取第一个章节的 ID"""
    if not response or not response.get("data"):
        return None

    data = response["data"]
    if isinstance(data, list) and len(data) > 0:
        return data[0].get("id")
    return None


def run_tests():
    """运行所有测试"""
    # 创建输出目录
    os.makedirs(OUTPUT_DIR, exist_ok=True)

    print("=" * 50)
    print("青羽阅读端功能自动化测试")
    print("=" * 50)
    print(f"测试服务器: {BASE_URL}")
    print(f"测试时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")

    test_stats["start_time"] = time.time()

    # ========================================
    # 一、服务器健康检查
    # ========================================
    print_header("0. 服务器健康检查")

    try:
        response = requests.get(f"{BASE_URL}/ping", timeout=5)
        if response.status_code == 200:
            print(f"{Fore.GREEN}[OK]{Style.RESET_ALL} 服务器运行正常" if HAS_COLOR else "[OK] 服务器运行正常")
        else:
            print(f"{Fore.RED}[FAIL]{Style.RESET_ALL} 服务器响应异常" if HAS_COLOR else "[FAIL] 服务器响应异常")
            return
    except requests.RequestException:
        print(f"{Fore.RED}[FAIL]{Style.RESET_ALL} 无法连接到服务器，请确保服务器正在运行" if HAS_COLOR else "[FAIL] 无法连接到服务器")
        return

    # ========================================
    # 二、书城浏览功能测试
    # ========================================
    print_header("1. 书城浏览功能测试")

    # 推荐书籍
    success, response, rt = test_api("获取推荐书籍", "/bookstore/books/recommended")
    print_result("获取推荐书籍", success, rt)
    book_id = extract_book_id(response)

    # 精选书籍
    success, response, rt = test_api("获取精选书籍", "/bookstore/books/featured")
    print_result("获取精选书籍", success, rt)
    if not book_id:
        book_id = extract_book_id(response)

    # 搜索书籍
    success, response, rt = test_api("搜索书籍-关键词书", "/bookstore/books/search?keyword=书&page=1&limit=5")
    print_result("搜索书籍", success, rt)
    if not book_id:
        book_id = extract_book_id(response)

    success, response, rt = test_api("搜索书籍-关键词仙", "/bookstore/books/search?keyword=仙&page=1&limit=5")
    print_result("搜索书籍(关键词:仙)", success, rt)

    # 首页数据
    success, response, rt = test_api("获取书城首页", "/bookstore/homepage")
    print_result("获取书城首页", success, rt)

    # ========================================
    # 三、榜单功能测试
    # ========================================
    print_header("2. 榜单功能测试")

    # 各种榜单
    success, response, rt = test_api("获取实时榜", "/bookstore/rankings/realtime")
    print_result("实时榜", success, rt)

    success, response, rt = test_api("获取周榜", "/bookstore/rankings/weekly")
    print_result("周榜", success, rt)

    success, response, rt = test_api("获取月榜", "/bookstore/rankings/monthly")
    print_result("月榜", success, rt)

    success, response, rt = test_api("获取新人榜", "/bookstore/rankings/newbie")
    print_result("新人榜", success, rt)

    # ========================================
    # 四、分类功能测试
    # ========================================
    print_header("3. 分类功能测试")

    success, response, rt = test_api("获取分类树", "/bookstore/categories/tree")
    print_result("获取分类树", success, rt)

    # 提取第一个分类ID
    category_id = None
    if response and response.get("data"):
        categories = response["data"]
        if isinstance(categories, list) and len(categories) > 0:
            category_id = categories[0].get("id")

    if category_id:
        success, response, rt = test_api("获取分类详情", f"/bookstore/categories/{category_id}")
        print_result("获取分类详情", success, rt)

        success, response, rt = test_api("获取分类书籍", f"/bookstore/categories/{category_id}/books?page=1&limit=10")
        print_result("获取分类书籍", success, rt)

    # ========================================
    # 五、书籍详情测试
    # ========================================
    print_header("4. 书籍详情测试")

    if book_id:
        print(f"使用书籍 ID: {book_id}")

        success, response, rt = test_api("获取书籍详情", f"/bookstore/books/{book_id}")
        print_result("获取书籍详情", success, rt)
    else:
        print(f"{Fore.YELLOW}[跳过]{Style.RESET_ALL} 未找到可用的书籍 ID" if HAS_COLOR else "[跳过] 未找到可用的书籍 ID")

    # ========================================
    # 六、Banner测试
    # ========================================
    print_header("5. Banner功能测试")

    success, response, rt = test_api("获取活动Banner", "/bookstore/banners")
    print_result("获取活动Banner", success, rt)

    # ========================================
    # 测试总结
    # ========================================
    test_stats["end_time"] = time.time()
    total_time = test_stats["end_time"] - test_stats["start_time"]

    print_header("测试总结")

    print(f"\n总测试项: {test_stats['total']}")
    print(f"{Fore.GREEN}通过: {test_stats['passed']}{Style.RESET_ALL}" if HAS_COLOR else f"通过: {test_stats['passed']}")
    print(f"{Fore.RED}失败: {test_stats['failed']}{Style.RESET_ALL}" if HAS_COLOR else f"失败: {test_stats['failed']}")
    print(f"总耗时: {total_time:.2f} 秒")

    if test_stats["failed"] == 0:
        print(f"\n{Fore.GREEN}[SUCCESS] 所有测试通过！{Style.RESET_ALL}" if HAS_COLOR else "\n[SUCCESS] 所有测试通过！")
        pass_rate = 100
    else:
        pass_rate = (test_stats["passed"] * 100) // test_stats["total"]
        print(f"\n{Fore.YELLOW}通过率: {pass_rate}%{Style.RESET_ALL}" if HAS_COLOR else f"\n通过率: {pass_rate}%")

    # 生成测试报告
    report_file = os.path.join(OUTPUT_DIR, f"test_report_{datetime.now().strftime('%Y%m%d_%H%M%S')}.txt")
    with open(report_file, 'w', encoding='utf-8') as f:
        f.write("=" * 50 + "\n")
        f.write("青羽阅读端功能测试报告\n")
        f.write("=" * 50 + "\n")
        f.write(f"测试时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
        f.write(f"测试服务器: {BASE_URL}\n")
        f.write("\n")
        f.write("测试结果统计:\n")
        f.write(f"- 总测试项: {test_stats['total']}\n")
        f.write(f"- 通过项: {test_stats['passed']}\n")
        f.write(f"- 失败项: {test_stats['failed']}\n")
        f.write(f"- 通过率: {pass_rate}%\n")
        f.write(f"- 总耗时: {total_time:.2f} 秒\n")
        f.write("\n")
        f.write(f"详细结果请查看: {OUTPUT_DIR} 目录\n")

    print(f"\n测试报告已保存: {report_file}")
    print(f"详细响应数据保存在: {OUTPUT_DIR} 目录\n")

    return test_stats["failed"] == 0


def show_sample_data():
    """显示一些示例数据"""
    print_header("示例数据展示")

    # 读取推荐书籍数据
    recommend_file = os.path.join(OUTPUT_DIR, "获取推荐书籍.json")
    if os.path.exists(recommend_file):
        with open(recommend_file, 'r', encoding='utf-8') as f:
            data = json.load(f)
            if data.get("data"):
                books = data["data"][:3]  # 只显示前3本
                print("\n推荐书籍示例:")
                for i, book in enumerate(books, 1):
                    print(f"\n{i}. {book.get('title', 'N/A')}")
                    print(f"   作者: {book.get('author', 'N/A')}")
                    print(f"   分类: {', '.join(book.get('categories', []))}")
                    print(f"   字数: {book.get('word_count', 0):,}")
                    print(f"   章节: {book.get('chapter_count', 0)}")


def main():
    """主函数"""
    print("\n" + "=" * 50)
    print("青羽阅读端功能自动化测试")
    print("=" * 50 + "\n")

    # 检查是否安装了 colorama（用于彩色输出）
    if not HAS_COLOR:
        print("提示: 安装 colorama 可以获得彩色输出效果")
        print("运行命令: pip install colorama\n")

    # 运行测试
    success = run_tests()

    # 显示示例数据
    try:
        show_sample_data()
    except Exception as e:
        pass

    # 退出码
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()

