#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""直接测试登录API"""

import os
import requests

DEFAULT_TEST_PASSWORD = "password"


def warn_if_using_default_password(env_name: str, value: str) -> None:
    if value == DEFAULT_TEST_PASSWORD:
        print(f"[WARN] {env_name} 未设置，当前使用默认测试密码占位值")

# 测试登录
def test_login():
    url = "http://localhost:8080/api/v1/login"
    default_user_password = os.getenv("QINGYU_TEST_USER_PASSWORD", DEFAULT_TEST_PASSWORD)
    default_admin_password = os.getenv("QINGYU_TEST_ADMIN_PASSWORD", DEFAULT_TEST_PASSWORD)

    test_cases = [
        {"username": "test_user01", "password": default_user_password, "desc": "普通用户"},
        {"username": "admin", "password": default_admin_password, "desc": "管理员"},
    ]

    print("="*60)
    print("直接测试登录API")
    print("="*60)
    print("密码来源: QINGYU_TEST_USER_PASSWORD / QINGYU_TEST_ADMIN_PASSWORD")
    warn_if_using_default_password("QINGYU_TEST_USER_PASSWORD", default_user_password)
    warn_if_using_default_password("QINGYU_TEST_ADMIN_PASSWORD", default_admin_password)

    for case in test_cases:
        print(f"\n测试 {case['desc']} ({case['username']})...")

        try:
            response = requests.post(
                url,
                json={
                    "username": case["username"],
                    "password": case["password"]
                },
                headers={"Content-Type": "application/json"},
                timeout=5
            )

            print(f"  状态码: {response.status_code}")

            if response.status_code == 200:
                data = response.json()
                print(f"  [OK] 登录成功!")
                print("  响应体: 已隐藏敏感字段")
                print(f"  用户ID: {data.get('data', {}).get('user_id', 'N/A')}")
                print("  Token: 已获取（已隐藏）")
            else:
                data = response.json()
                print("  响应体: 已隐藏敏感字段")
                print(f"  [FAIL] 登录失败: {data.get('message', 'N/A')}")

        except requests.exceptions.ConnectionError:
            print(f"  [ERROR] 无法连接到服务器 (http://localhost:8080)")
            print(f"  请确保服务器正在运行")
            return
        except Exception as e:
            print(f"  [ERROR] 请求失败: {e}")

if __name__ == "__main__":
    test_login()


