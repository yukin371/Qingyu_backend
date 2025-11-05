#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""直接测试登录API"""

import requests
import json

# 测试登录
def test_login():
    url = "http://localhost:8080/api/v1/login"

    test_cases = [
        {"username": "test_user01", "password": "Test@123456", "desc": "普通用户"},
        {"username": "admin", "password": "Admin@123456", "desc": "管理员"},
    ]

    print("="*60)
    print("直接测试登录API")
    print("="*60)

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
            print(f"  响应体: {response.text[:200]}")

            if response.status_code == 200:
                data = response.json()
                print(f"  [OK] 登录成功!")
                print(f"  用户ID: {data.get('data', {}).get('user_id', 'N/A')}")
                print(f"  Token: {data.get('data', {}).get('token', 'N/A')[:50]}...")
            else:
                data = response.json()
                print(f"  [FAIL] 登录失败: {data.get('message', 'N/A')}")

        except requests.exceptions.ConnectionError:
            print(f"  [ERROR] 无法连接到服务器 (http://localhost:8080)")
            print(f"  请确保服务器正在运行")
            return
        except Exception as e:
            print(f"  [ERROR] 请求失败: {e}")

if __name__ == "__main__":
    test_login()


