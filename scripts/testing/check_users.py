#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""检查MongoDB中的测试用户"""

import os

from pymongo import MongoClient

default_auth_uri = "mongodb://localhost:27017"
configured_auth_uri = os.getenv("QINGYU_TEST_MONGO_AUTH_URI", default_auth_uri)

# 尝试两种连接方式
connections = [
    ("无认证", "mongodb://localhost:27017/", "qingyu_test"),
    ("带认证", configured_auth_uri, "qingyu_test"),
]

def mask_uri(uri: str) -> str:
    if "@" not in uri:
        return uri
    prefix, suffix = uri.split("@", 1)
    if "://" in prefix:
        scheme, _ = prefix.split("://", 1)
        return f"{scheme}://***@{suffix}"
    return f"***@{suffix}"

def mask_email(email: str) -> str:
    if "@" not in email:
        return email
    name, domain = email.split("@", 1)
    if len(name) <= 1:
        return f"***@{domain}"
    return f"{name[:1]}***@{domain}"


if configured_auth_uri == default_auth_uri:
    print("[WARN] QINGYU_TEST_MONGO_AUTH_URI 未设置，‘带认证’连接当前仍使用默认本地 MongoDB 地址")

for name, uri, db_name in connections:
    print(f"\n{'='*60}")
    print(f"尝试连接: {name}")
    print(f"连接地址(脱敏): {mask_uri(uri)}")
    print('='*60)

    try:
        client = MongoClient(uri, serverSelectionTimeoutMS=2000)
        db = client[db_name]

        # 尝试读取用户
        users = list(db['users'].find({}, {"username": 1, "email": 1, "role": 1}))

        print(f"[OK] 连接成功!")
        print(f"数据库: {db_name}")
        print(f"用户数量: {len(users)}")

        if users:
            print(f"\n{'用户名':<20} {'邮箱':<30} {'角色':<10}")
            print("-"*65)
            for user in users:
                print(f"{user.get('username', 'N/A'):<20} {mask_email(user.get('email', 'N/A')):<30} {user.get('role', 'N/A'):<10}")
        else:
            print("\n没有找到用户数据")

        client.close()

    except Exception as e:
        print(f"[ERROR] 连接失败: {e}")


