#!/usr/bin/env python3
"""测试bcrypt密码加密"""

import os

import bcrypt
from pymongo import MongoClient

DEFAULT_TEST_PASSWORD = "password"
password = os.getenv("QINGYU_TEST_USER_PASSWORD", DEFAULT_TEST_PASSWORD)


def mask_hash(value: str) -> str:
    if len(value) <= 8:
        return "***"
    return f"{value[:4]}***{value[-4:]}"


if password == DEFAULT_TEST_PASSWORD:
    print("[WARN] QINGYU_TEST_USER_PASSWORD 未设置，当前使用默认测试密码占位值")

# Python生成
python_hash = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt())
print(f"Python bcrypt hash: {mask_hash(python_hash.decode('utf-8'))}")
print(f"类型: {type(python_hash)}")

# 验证
result = bcrypt.checkpw(password.encode('utf-8'), python_hash)
print(f"Python验证结果: {result}")

# 查看数据库中的密码哈希
print("\n" + "="*60)
print("数据库中的密码哈希摘要:")
print("="*60)

client = MongoClient('mongodb://localhost:27017/')
db = client['qingyu_test']
user = db['users'].find_one({"username": "test_user01"})

if user:
    print(f"用户名: {user['username']}")
    print(f"密码哈希摘要: {mask_hash(user.get('password_hash', 'N/A'))}")
    print(f"密码哈希长度: {len(user.get('password_hash', ''))}")

    # 尝试验证
    try:
        db_hash = user['password_hash'].encode('utf-8')
        result = bcrypt.checkpw(password.encode('utf-8'), db_hash)
        print(f"\n数据库密码验证结果: {result}")
    except Exception as e:
        print(f"\n验证失败: {e}")
else:
    print("未找到test_user01用户")

client.close()


