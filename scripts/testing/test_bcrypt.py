#!/usr/bin/env python3
"""测试bcrypt密码加密"""

import bcrypt
from pymongo import MongoClient

password = "Test@123456"

# Python生成
python_hash = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt())
print(f"Python bcrypt hash: {python_hash}")
print(f"类型: {type(python_hash)}")
print(f"解码后: {python_hash.decode('utf-8')}")

# 验证
result = bcrypt.checkpw(password.encode('utf-8'), python_hash)
print(f"Python验证结果: {result}")

# 查看数据库中的密码哈希
print("\n" + "="*60)
print("数据库中的密码哈希:")
print("="*60)

client = MongoClient('mongodb://localhost:27017/')
db = client['qingyu_test']
user = db['users'].find_one({"username": "test_user01"})

if user:
    print(f"用户名: {user['username']}")
    print(f"密码哈希: {user.get('password_hash', 'N/A')}")
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


