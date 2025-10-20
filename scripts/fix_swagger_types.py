#!/usr/bin/env python3
"""
修复Swagger注释中的类型引用问题
将所有 response.Response 替换为 shared.APIResponse
"""

import os
import sys
from pathlib import Path

def fix_swagger_types(file_path):
    """修复单个文件中的Swagger类型引用"""
    try:
        # 读取文件（使用UTF-8编码）
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()

        # 替换response.Response为shared.APIResponse
        original_content = content
        content = content.replace('response.Response', 'shared.APIResponse')

        # 如果内容有变化，写回文件
        if content != original_content:
            with open(file_path, 'w', encoding='utf-8', newline='\n') as f:
                f.write(content)
            print(f"[FIXED] {file_path}")
            return True
        else:
            print(f"[SKIP] {file_path}")
            return False
    except Exception as e:
        print(f"[ERROR] {file_path}: {e}")
        return False

def main():
    """主函数：遍历api/v1目录下的所有Go文件"""
    api_dir = Path('api/v1')

    if not api_dir.exists():
        print(f"Error: Directory {api_dir} not found!")
        sys.exit(1)

    # 查找所有.go文件
    go_files = list(api_dir.rglob('*.go'))

    if not go_files:
        print(f"No Go files found in {api_dir}")
        sys.exit(1)

    print(f"Found {len(go_files)} Go files in {api_dir}")
    print("=" * 60)

    fixed_count = 0
    for go_file in go_files:
        if fix_swagger_types(go_file):
            fixed_count += 1

    print("=" * 60)
    print(f"\nSummary: Fixed {fixed_count} out of {len(go_files)} files")

if __name__ == '__main__':
    main()

