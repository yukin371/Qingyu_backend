#!/usr/bin/env python3
"""
JSON Tag修复工具
用法: python scripts/fix_json_tags.py models/
"""
import re
import sys
import os
from pathlib import Path


def camel_to_snake(name):
    """将驼峰命名转换为蛇形命名"""
    result = []
    for i, char in enumerate(name):
        if char.isupper() and i > 0:
            prev = name[i-1]
            # 前一个字符是小写，或前一个是大写但下一个是小写
            if prev.islower():
                result.append('_')
            elif i < len(name) - 1 and name[i+1].islower():
                # XMLParser -> XML_Parser
                result.append('_')
        result.append(char.lower())
    return ''.join(result)


def fix_file(filepath):
    """修复单个文件中的JSON tag"""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    lines = content.split('\n')
    new_lines = []
    fixed_count = 0

    # 匹配结构体字段行中的 json:"$1$2"
    # 例如: 	FieldName string `bson:"field_name" json:"$1$2"`
    pattern = re.compile(r'\t([A-Z][a-zA-Z0-9]*)\s+[^`]+`[^`]*json:"\$1\$2"[^`]*`')

    for line in lines:
        new_line = line
        # 首先检查是否包含json:"$1$2"
        if 'json:"$1$2"' in line:
            # 尝试提取字段名
            # 格式1: FieldName type `json:"$1$2" ...`
            # 格式2: FieldName type `... json:"$1$2" ...`
            field_match = re.search(r'\t([A-Z][a-zA-Z0-9]*)\s+', line)
            if field_match:
                field_name = field_match.group(1)
                snake_name = camel_to_snake(field_name)
                new_line = line.replace('json:"$1$2"', f'json:"{snake_name}"')
                if new_line != line:
                    fixed_count += 1
        new_lines.append(new_line)

    if fixed_count > 0:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write('\n'.join(new_lines))
        return fixed_count
    return 0


def find_go_files(path):
    """递归查找所有.go文件"""
    path = Path(path)
    if path.is_file() and path.suffix == '.go':
        return [path]
    if path.is_dir():
        files = []
        for item in path.rglob('*.go'):
            files.append(item)
        return files
    return []


def main():
    if len(sys.argv) < 2:
        print("用法: python scripts/fix_json_tags.py <文件或目录>...")
        sys.exit(1)

    total_fixed = 0

    for target in sys.argv[1:]:
        files = find_go_files(target)
        for filepath in files:
            try:
                fixed = fix_file(filepath)
                if fixed > 0:
                    print(f"[OK] {filepath} ({fixed} fixes)")
                    total_fixed += fixed
            except Exception as e:
                print(f"[ERROR] {filepath}: {e}", file=sys.stderr)

    print(f"\n完成! 共修复 {total_fixed} 处")


if __name__ == '__main__':
    main()
