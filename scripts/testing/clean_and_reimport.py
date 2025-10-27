#!/usr/bin/env python3
"""清理旧用户数据并重新导入"""

from pymongo import MongoClient
import subprocess
import sys

print("="*60)
print("步骤1: 清理旧用户数据")
print("="*60)

try:
    client = MongoClient('mongodb://localhost:27017')
    db = client['qingyu_test']
    result = db['users'].delete_many({})
    print(f"[OK] 删除了 {result.deleted_count} 个旧用户")
    client.close()
except Exception as e:
    print(f"[ERROR] 清理失败: {e}")
    sys.exit(1)

print("\n" + "="*60)
print("步骤2: 重新导入用户")
print("="*60)
print()

# 运行导入脚本
result = subprocess.run(
    [sys.executable, "scripts/testing/import_users_direct.py"],
    capture_output=False
)

sys.exit(result.returncode)


