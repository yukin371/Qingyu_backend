#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
青羽写作平台 API - 认证示例
演示如何使用Python进行API认证和调用
"""

import requests
import json
from typing import Optional, Dict, Any

class QingyuAPIClient:
    """青羽写作平台API客户端"""

    def __init__(self, base_url: str = "http://localhost:9090/api/v1"):
        self.base_url = base_url.rstrip('/')
        self.token: Optional[str] = None
        self.session = requests.Session()

    def _request(
        self,
        method: str,
        endpoint: str,
        data: Optional[Dict[str, Any]] = None,
        params: Optional[Dict[str, Any]] = None,
        require_auth: bool = False
    ) -> Dict[str, Any]:
        """发送HTTP请求"""
        url = f"{self.base_url}/{endpoint.lstrip('/')}"
        headers = {"Content-Type": "application/json"}

        if require_auth and self.token:
            headers["Authorization"] = f"Bearer {self.token}"

        try:
            response = self.session.request(
                method=method,
                url=url,
                json=data,
                params=params,
                headers=headers
            )
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            print(f"请求失败: {e}")
            if hasattr(e.response, 'text'):
                print(f"响应内容: {e.response.text}")
            raise

    def register(self, username: str, email: str, password: str, nickname: str) -> Dict[str, Any]:
        """用户注册"""
        print("\n=== 用户注册 ===")
        data = {
            "username": username,
            "email": email,
            "password": password,
            "nickname": nickname
        }
        result = self._request("POST", "/auth/register", data=data)
        print(f"注册成功: {result}")
        return result

    def login(self, username: str, password: str) -> Dict[str, Any]:
        """用户登录"""
        print("\n=== 用户登录 ===")
        data = {
            "username": username,
            "password": password
        }
        result = self._request("POST", "/auth/login", data=data)

        # 保存token
        if "data" in result:
            self.token = result["data"].get("token") or result["data"].get("access_token")
            if self.token:
                print(f"登录成功，Token: {self.token[:20]}...")
            else:
                print("登录成功，但未获取到Token")

        return result

    def get_profile(self) -> Dict[str, Any]:
        """获取用户信息"""
        print("\n=== 获取用户信息 ===")
        result = self._request("GET", "/user/profile", require_auth=True)
        print(f"用户信息: {json.dumps(result, indent=2, ensure_ascii=False)}")
        return result

    def refresh_token(self) -> Dict[str, Any]:
        """刷新Token"""
        print("\n=== 刷新Token ===")
        result = self._request("POST", "/auth/refresh", require_auth=True)
        if "data" in result and "token" in result["data"]:
            self.token = result["data"]["token"]
            print(f"Token刷新成功: {self.token[:20]}...")
        return result

    def logout(self) -> Dict[str, Any]:
        """登出"""
        print("\n=== 登出 ===")
        result = self._request("POST", "/auth/logout", require_auth=True)
        self.token = None
        print("登出成功")
        return result


def main():
    """主函数 - 演示认证流程"""

    # 创建API客户端
    client = QingyuAPIClient()

    try:
        # 1. 用户注册
        print("\n【步骤1】用户注册")
        try:
            client.register(
                username="testuser_py",
                email="testpy@example.com",
                password="SecurePass123!",
                nickname="Python测试用户"
            )
        except Exception as e:
            print(f"注册可能失败（用户已存在）: {e}")

        # 2. 用户登录
        print("\n【步骤2】用户登录")
        client.login(username="testuser_py", password="SecurePass123!")

        # 3. 访问受保护的API
        print("\n【步骤3】访问受保护的API")
        if client.token:
            client.get_profile()
        else:
            print("未获取到Token，跳过")

        # 4. 刷新Token
        print("\n【步骤4】刷新Token")
        if client.token:
            client.refresh_token()

        # 5. 登出
        print("\n【步骤5】登出")
        if client.token:
            client.logout()

        print("\n=== 认证示例完成 ===")

    except Exception as e:
        print(f"\n错误: {e}")
        return 1

    return 0


if __name__ == "__main__":
    exit(main())
