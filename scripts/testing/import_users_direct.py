#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
青羽后端 - 测试用户直接导入工具（纯Python版本）
直接连接MongoDB导入测试用户，避免Go脚本的数据库初始化问题
"""

import sys
import bcrypt
from datetime import datetime
from pymongo import MongoClient
from pymongo.errors import DuplicateKeyError

# 测试用户数据
TEST_USERS = [
    {
        "username": "admin",
        "email": "admin@qingyu.com",
        "password": "Admin@123456",
        "role": "admin",
        "description": "系统管理员"
    },
    {
        "username": "vip_user01",
        "email": "vip01@qingyu.com",
        "password": "Vip@123456",
        "role": "vip",
        "description": "VIP测试用户1"
    },
    {
        "username": "vip_user02",
        "email": "vip02@qingyu.com",
        "password": "Vip@123456",
        "role": "vip",
        "description": "VIP测试用户2"
    },
    {
        "username": "test_user01",
        "email": "test01@qingyu.com",
        "password": "Test@123456",
        "role": "user",
        "description": "普通测试用户1"
    },
    {
        "username": "test_user02",
        "email": "test02@qingyu.com",
        "password": "Test@123456",
        "role": "user",
        "description": "普通测试用户2"
    },
    {
        "username": "test_user03",
        "email": "test03@qingyu.com",
        "password": "Test@123456",
        "role": "user",
        "description": "普通测试用户3"
    },
    {
        "username": "test_user04",
        "email": "test04@qingyu.com",
        "password": "Test@123456",
        "role": "user",
        "description": "普通测试用户4"
    },
    {
        "username": "test_user05",
        "email": "test05@qingyu.com",
        "password": "Test@123456",
        "role": "user",
        "description": "普通测试用户5"
    }
]

def hash_password(password):
    """使用bcrypt加密密码"""
    # 使用bcrypt生成密码哈希（与Go的bcrypt兼容）
    hashed = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt())
    return hashed.decode('utf-8')

def import_users():
    """导入测试用户到MongoDB"""
    print("="*60)
    print("青羽后端 - 测试用户导入工具 (Python直连版)")
    print("="*60)
    print()

    # 连接MongoDB
    try:
        # 使用无认证连接（本地开发环境）
        mongo_uri = 'mongodb://localhost:27017'
        print(f"连接MongoDB: {mongo_uri}")
        client = MongoClient(mongo_uri)
        db = client['qingyu_test']
        users_collection = db['users']

        print("[OK] MongoDB连接成功")
        print()
    except Exception as e:
        print(f"[ERROR] MongoDB连接失败: {e}")
        print()
        print("请确保:")
        print("  1. MongoDB服务已启动")
        print("  2. 连接地址正确: mongodb://admin:password@localhost:27017")
        print("  3. 数据库名称正确: qingyu_test")
        print("  4. MongoDB认证信息正确")
        return 1

    # 导入用户
    success_count = 0
    exist_count = 0
    failed_count = 0

    print(f"开始导入 {len(TEST_USERS)} 个测试用户...")
    print()

    for i, user_data in enumerate(TEST_USERS, 1):
        print(f"[{i}/{len(TEST_USERS)}] {user_data['description']} ({user_data['username']})...", end=" ")

        try:
            # 检查用户是否已存在
            existing = users_collection.find_one({
                "$or": [
                    {"username": user_data["username"]},
                    {"email": user_data["email"]}
                ]
            })

            if existing:
                print("已存在")
                exist_count += 1
                continue

            # 创建用户文档
            now = datetime.utcnow()
            user_doc = {
                "username": user_data["username"],
                "email": user_data["email"],
                "password": hash_password(user_data["password"]),  # 注意：字段名是 password，不是 password_hash
                "role": user_data["role"],
                "status": "active",
                "email_verified": True,
                "phone_verified": False,
                "nickname": user_data["username"],
                "bio": user_data["description"],
                "avatar": "",
                "phone": "",
                "last_login_at": now,
                "last_login_ip": "127.0.0.1",
                "created_at": now,
                "updated_at": now
            }

            # 插入用户
            result = users_collection.insert_one(user_doc)
            print(f"[OK] 成功 (ID: {result.inserted_id})")
            success_count += 1

        except DuplicateKeyError:
            print("[EXIST] 已存在")
            exist_count += 1
        except Exception as e:
            print(f"[ERROR] 失败: {e}")
            failed_count += 1

    print()
    print("="*60)
    print("导入完成")
    print("="*60)
    print(f"成功: {success_count} 个")
    print(f"已存在: {exist_count} 个")
    print(f"失败: {failed_count} 个")
    print(f"总计: {len(TEST_USERS)} 个")
    print()

    if success_count > 0 or exist_count > 0:
        print("="*60)
        print("测试账号清单")
        print("="*60)
        print(f"{'角色':<12} {'用户名':<15} {'邮箱':<25} {'密码':<20}")
        print("-" * 75)
        for user in TEST_USERS:
            role_name = {"admin": "管理员", "vip": "VIP用户", "user": "普通用户"}[user["role"]]
            print(f"{role_name:<12} {user['username']:<15} {user['email']:<25} {user['password']:<20}")
        print()
        print("提示:")
        print("  - 登录时使用 username 字段，不是 email")
        print("  - 例如: username='test_user01', password='Test@123456'")
        print()

    client.close()
    return 0 if failed_count == 0 else 1

if __name__ == "__main__":
    try:
        sys.exit(import_users())
    except KeyboardInterrupt:
        print("\n\n操作已取消")
        sys.exit(1)
    except Exception as e:
        print(f"\n[ERROR] 发生错误: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

