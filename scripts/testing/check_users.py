#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""检查MongoDB中的测试用户"""

from pymongo import MongoClient

# 尝试两种连接方式
connections = [
    ("无认证", "mongodb://localhost:27017/", "qingyu_test"),
    ("带认证", "mongodb://admin:password@localhost:27017", "qingyu_test"),
]

for name, uri, db_name in connections:
    print(f"\n{'='*60}")
    print(f"尝试连接: {name}")
    print(f"URI: {uri}")
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
                print(f"{user.get('username', 'N/A'):<20} {user.get('email', 'N/A'):<30} {user.get('role', 'N/A'):<10}")
        else:
            print("\n没有找到用户数据")

        client.close()

    except Exception as e:
        print(f"[ERROR] 连接失败: {e}")


