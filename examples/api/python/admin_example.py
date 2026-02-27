#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
青羽写作平台 API - 管理员示例
演示如何使用Python调用管理员API
"""

import requests
import json
from typing import Optional, Dict, Any, List


class AdminAPIClient:
    """管理员API客户端"""

    def __init__(self, base_url: str = "http://localhost:9090/api/v1", admin_token: Optional[str] = None):
        self.base_url = base_url.rstrip('/')
        self.admin_token = admin_token
        self.session = requests.Session()

    def _request(
        self,
        method: str,
        endpoint: str,
        data: Optional[Dict[str, Any]] = None,
        params: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """发送HTTP请求"""
        url = f"{self.base_url}/{endpoint.lstrip('/')}"
        headers = {"Content-Type": "application/json"}

        if self.admin_token:
            headers["Authorization"] = f"Bearer {self.admin_token}"

        try:
            response = self.session.request(
                method=method,
                url=url,
                json=data,
                params=params,
                headers=headers,
                timeout=30
            )
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            print(f"请求失败: {e}")
            if hasattr(e.response, 'text'):
                print(f"响应内容: {e.response.text}")
            raise

    # ==================== 用户管理 ====================

    def get_users(self, page: int = 1, limit: int = 20, status: Optional[str] = None) -> Dict[str, Any]:
        """获取用户列表"""
        print(f"\n=== 获取用户列表 (第{page}页) ===")
        params = {"page": page, "limit": limit}
        if status:
            params["status"] = status
        result = self._request("GET", "/admin/users", params=params)
        if "data" in result and "users" in result["data"]:
            print(f"获取到 {len(result['data']['users'])} 个用户")
        return result

    def create_user(
        self,
        username: str,
        email: str,
        password: str,
        nickname: str,
        role: str = "user"
    ) -> Dict[str, Any]:
        """创建用户"""
        print(f"\n=== 创建用户 ({username}) ===")
        data = {
            "username": username,
            "email": email,
            "password": password,
            "nickname": nickname,
            "role": role
        }
        result = self._request("POST", "/admin/users", data=data)
        print(f"用户创建成功: {result.get('message', 'OK')}")
        return result

    def update_user(self, user_id: str, **kwargs) -> Dict[str, Any]:
        """更新用户信息"""
        print(f"\n=== 更新用户信息 (ID: {user_id}) ===")
        result = self._request("PUT", f"/admin/users/{user_id}", data=kwargs)
        print("用户信息更新成功")
        return result

    def delete_user(self, user_id: str) -> Dict[str, Any]:
        """删除用户"""
        print(f"\n=== 删除用户 (ID: {user_id}) ===")
        result = self._request("DELETE", f"/admin/users/{user_id}")
        print("用户删除成功")
        return result

    # ==================== 权限管理 ====================

    def get_permissions(self) -> Dict[str, Any]:
        """获取所有权限"""
        print("\n=== 获取权限列表 ===")
        result = self._request("GET", "/admin/permissions")
        if "data" in result:
            print(f"获取到 {len(result['data'])} 个权限")
        return result

    def create_permission(
        self,
        code: str,
        name: str,
        description: str,
        category: str
    ) -> Dict[str, Any]:
        """创建权限"""
        print(f"\n=== 创建权限 ({code}) ===")
        data = {
            "code": code,
            "name": name,
            "description": description,
            "category": category
        }
        result = self._request("POST", "/admin/permissions", data=data)
        print("权限创建成功")
        return result

    # ==================== 权限模板管理 ====================

    def get_permission_templates(self, category: Optional[str] = None) -> Dict[str, Any]:
        """获取权限模板列表"""
        print("\n=== 获取权限模板列表 ===")
        params = {}
        if category:
            params["category"] = category
        result = self._request("GET", "/admin/permission-templates", params=params)
        if "data" in result:
            print(f"获取到 {len(result['data'])} 个模板")
        return result

    def create_permission_template(
        self,
        name: str,
        code: str,
        description: str,
        permissions: List[str],
        category: str
    ) -> Dict[str, Any]:
        """创建权限模板"""
        print(f"\n=== 创建权限模板 ({name}) ===")
        data = {
            "name": name,
            "code": code,
            "description": description,
            "permissions": permissions,
            "category": category
        }
        result = self._request("POST", "/admin/permission-templates", data=data)
        print("权限模板创建成功")
        return result

    def apply_template_to_role(self, template_id: str, role_id: str) -> Dict[str, Any]:
        """应用权限模板到角色"""
        print(f"\n=== 应用模板到角色 (模板: {template_id}, 角色: {role_id}) ===")
        data = {"roleId": role_id}
        result = self._request("POST", f"/admin/permission-templates/{template_id}/apply", data=data)
        print("模板应用成功")
        return result

    # ==================== 审计日志 ====================

    def get_audit_logs(
        self,
        page: int = 1,
        size: int = 20,
        operation: Optional[str] = None,
        resource_type: Optional[str] = None
    ) -> Dict[str, Any]:
        """获取审计日志"""
        print(f"\n=== 获取审计日志 (第{page}页) ===")
        params = {"page": page, "size": size}
        if operation:
            params["operation"] = operation
        if resource_type:
            params["resource_type"] = resource_type
        result = self._request("GET", "/admin/audit/trail", params=params)
        if "data" in result and "logs" in result["data"]:
            print(f"获取到 {len(result['data']['logs'])} 条日志")
        return result

    # ==================== 统计分析 ====================

    def get_user_growth_trend(
        self,
        start_date: str,
        end_date: str,
        interval: str = "daily"
    ) -> Dict[str, Any]:
        """获取用户增长趋势"""
        print(f"\n=== 获取用户增长趋势 ({start_date} 至 {end_date}) ===")
        params = {
            "start_date": start_date,
            "end_date": end_date,
            "interval": interval
        }
        result = self._request("GET", "/admin/analytics/user-growth", params=params)
        if "data" in result and "trend" in result["data"]:
            print(f"获取到 {len(result['data']['trend'])} 条趋势数据")
        return result

    def get_content_statistics(
        self,
        start_date: Optional[str] = None,
        end_date: Optional[str] = None
    ) -> Dict[str, Any]:
        """获取内容统计"""
        print("\n=== 获取内容统计 ===")
        params = {}
        if start_date:
            params["start_date"] = start_date
        if end_date:
            params["end_date"] = end_date
        result = self._request("GET", "/admin/analytics/content-statistics", params=params)
        print("内容统计获取成功")
        return result

    def get_system_overview(self) -> Dict[str, Any]:
        """获取系统概览"""
        print("\n=== 获取系统概览 ===")
        result = self._request("GET", "/admin/analytics/system-overview")
        print("系统概览获取成功")
        return result

    # ==================== 公告管理 ====================

    def create_announcement(
        self,
        title: str,
        content: str,
        announcement_type: str,
        priority: str = "normal",
        is_pinned: bool = False
    ) -> Dict[str, Any]:
        """创建公告"""
        print(f"\n=== 创建公告 ({title}) ===")
        data = {
            "title": title,
            "content": content,
            "type": announcement_type,
            "priority": priority,
            "is_pinned": is_pinned
        }
        result = self._request("POST", "/admin/announcements", data=data)
        print("公告创建成功")
        return result

    def get_announcements(self) -> Dict[str, Any]:
        """获取公告列表"""
        print("\n=== 获取公告列表 ===")
        result = self._request("GET", "/admin/announcements")
        if "data" in result and "announcements" in result["data"]:
            print(f"获取到 {len(result['data']['announcements'])} 个公告")
        return result


def main():
    """主函数 - 演示管理员API"""

    import os

    # 从环境变量获取管理员token
    admin_token = os.getenv("QINGYU_ADMIN_TOKEN")

    if not admin_token:
        print("错误: 未设置管理员Token")
        print("请设置环境变量: export QINGYU_ADMIN_TOKEN=your_admin_token")
        return 1

    # 创建API客户端
    client = AdminAPIClient(admin_token=admin_token)

    try:
        # 1. 获取用户列表
        client.get_users(page=1, limit=10)

        # 2. 创建用户
        try:
            client.create_user(
                username="newuser_py",
                email="newuser_py@example.com",
                password="SecurePass123!",
                nickname="Python创建的用户",
                role="user"
            )
        except Exception as e:
            print(f"创建用户失败（可能已存在）: {e}")

        # 3. 获取权限列表
        client.get_permissions()

        # 4. 创建权限
        try:
            client.create_permission(
                code="test.permission.py",
                name="Python测试权限",
                description="通过Python脚本创建的测试权限",
                category="test"
            )
        except Exception as e:
            print(f"创建权限失败（可能已存在）: {e}")

        # 5. 获取权限模板
        client.get_permission_templates()

        # 6. 创建权限模板
        try:
            client.create_permission_template(
                name="Python测试模板",
                code="python_test_template",
                description="通过Python脚本创建的测试模板",
                permissions=["user.read", "content.read"],
                category="test"
            )
        except Exception as e:
            print(f"创建模板失败（可能已存在）: {e}")

        # 7. 获取审计日志
        client.get_audit_logs(page=1, size=10)

        # 8. 获取用户增长趋势
        client.get_user_growth_trend(
            start_date="2026-01-01",
            end_date="2026-01-31",
            interval="daily"
        )

        # 9. 获取内容统计
        client.get_content_statistics()

        # 10. 获取系统概览
        client.get_system_overview()

        # 11. 创建公告
        client.create_announcement(
            title="Python测试公告",
            content="这是一个通过Python脚本创建的测试公告",
            announcement_type="info",
            priority="normal",
            is_pinned=False
        )

        print("\n=== 管理员API示例完成 ===")

    except Exception as e:
        print(f"\n错误: {e}")
        return 1

    return 0


if __name__ == "__main__":
    exit(main())
