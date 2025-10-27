#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
青羽后端 - 数据库清理脚本（Python版本）
避免批处理文件的编码问题
"""

import subprocess
import sys

def cleanup_mongodb():
    """清理MongoDB测试数据"""
    
    print("=" * 60)
    print("青羽后端 - 测试数据清理")
    print("=" * 60)
    print()
    
    # 检查MongoDB
    print("[1/2] 检查 MongoDB 服务...")
    try:
        result = subprocess.run(
            ["mongosh", "--eval", "db.version()", "--quiet"],
            capture_output=True,
            text=True,
            encoding='utf-8',
            timeout=5
        )
        if result.returncode != 0:
            print("[错误] MongoDB 服务未运行")
            print("提示: 请先启动 MongoDB 服务")
            return False
        print("[成功] MongoDB 服务正常运行")
        print()
    except Exception as e:
        print(f"[错误] 无法连接到 MongoDB: {e}")
        return False
    
    # 确认清理
    print("[警告] 此操作将清空以下数据:")
    print("  - books (书籍)")
    print("  - chapters (章节)")
    print("  - users (用户，保留系统用户)")
    print("  - ranking_items (榜单)")
    print("  - reading_progress (阅读进度)")
    print("  - annotations (书签和笔记)")
    print("  - user_collections (收藏)")
    print()
    
    confirm = input("确认清理？(yes/no): ")
    if confirm.lower() != 'yes':
        print("操作已取消")
        return False
    
    print()
    print("[2/2] 连接到数据库并清理数据...")
    
    # MongoDB 清理脚本
    cleanup_script = """
use qingyu_test;

print('[清理] 删除书籍数据...');
db.books.deleteMany({});
print('  ✓ books 集合已清空');

print('[清理] 删除章节数据...');
db.chapters.deleteMany({});
print('  ✓ chapters 集合已清空');

print('[清理] 删除用户数据（保留系统用户）...');
db.users.deleteMany({ role: { $ne: 'system' } });
print('  ✓ users 集合已清空（保留系统用户）');

print('[清理] 删除榜单数据...');
db.ranking_items.deleteMany({});
print('  ✓ ranking_items 集合已清空');

print('[清理] 删除阅读进度...');
db.reading_progress.deleteMany({});
print('  ✓ reading_progress 集合已清空');

print('[清理] 删除书签和笔记...');
db.annotations.deleteMany({});
print('  ✓ annotations 集合已清空');

print('[清理] 删除收藏记录...');
db.user_collections.deleteMany({});
print('  ✓ user_collections 集合已清空');

print('[清理] 删除评论数据...');
db.comments.deleteMany({});
print('  ✓ comments 集合已清空');

print('[清理] 删除写作项目（测试项目）...');
db.projects.deleteMany({ title: /测试/ });
print('  ✓ projects 测试数据已清空');

print('');
print('[统计] 清理后数据统计：');
print('  - 书籍数量: ' + db.books.countDocuments());
print('  - 章节数量: ' + db.chapters.countDocuments());
print('  - 用户数量: ' + db.users.countDocuments());
print('  - 榜单数量: ' + db.ranking_items.countDocuments());
print('  - 阅读进度: ' + db.reading_progress.countDocuments());
print('  - 书签笔记: ' + db.annotations.countDocuments());
"""
    
    # 写入临时文件
    temp_file = "temp_cleanup_script.js"
    try:
        with open(temp_file, 'w', encoding='utf-8') as f:
            f.write(cleanup_script)
        
        # 执行清理
        result = subprocess.run(
            f"mongosh --quiet < {temp_file}",
            shell=True,
            capture_output=True,
            text=True,
            encoding='utf-8'
        )
        
        print(result.stdout)
        
        if result.returncode == 0 or "countDocuments" in result.stdout:
            print()
            print("=" * 60)
            print("✓ 测试数据清理成功")
            print("=" * 60)
            print()
            print("提示: 现在可以运行 import_test_users.go 导入测试用户")
            return True
        else:
            print()
            print("[错误] 数据清理失败")
            print(result.stderr)
            return False
            
    except Exception as e:
        print(f"[错误] 执行清理脚本失败: {e}")
        return False
    finally:
        # 清理临时文件
        import os
        if os.path.exists(temp_file):
            os.remove(temp_file)

if __name__ == "__main__":
    try:
        success = cleanup_mongodb()
        sys.exit(0 if success else 1)
    except KeyboardInterrupt:
        print("\n\n操作已取消")
        sys.exit(1)
    except Exception as e:
        print(f"\n[错误] {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)



