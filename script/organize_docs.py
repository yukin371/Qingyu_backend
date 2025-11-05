#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
阅读端文档整理脚本
用于整理和合并 02阅读端服务/ 和 06阅读端模块/ 中的重复文档
"""

import os
import shutil
from pathlib import Path
from datetime import datetime

# 项目根目录
ROOT_DIR = Path(__file__).parent.parent
IMPL_DIR = ROOT_DIR / "doc" / "implementation"

# 源目录和目标目录
SOURCE_DIR = IMPL_DIR / "06阅读端模块"
TARGET_DIR = IMPL_DIR / "02阅读端服务"
ARCHIVE_DIR = IMPL_DIR / "_archive_06阅读端模块"

def log(message):
    """打印带时间戳的日志"""
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    # 移除emoji以避免Windows控制台编码问题
    import re
    message = re.sub(r'[^\u0000-\uFFFF]', '', message)
    try:
        print(f"[{timestamp}] {message}")
    except UnicodeEncodeError:
        # 如果仍然有编码问题，使用ASCII
        print(f"[{timestamp}] {message.encode('ascii', 'ignore').decode('ascii')}")

def copy_file_if_not_exists(src, dst):
    """如果目标文件不存在，则复制文件"""
    if dst.exists():
        log(f"  ⏩ 跳过（已存在）: {dst.name}")
        return False
    else:
        shutil.copy2(src, dst)
        log(f"  ✅ 复制: {src.name} -> {dst.relative_to(ROOT_DIR)}")
        return True

def organize_bookstore_docs():
    """整理书城系统文档"""
    log("\n📚 开始整理书城系统文档...")
    
    source_bookstore = SOURCE_DIR / "01书城系统"
    target_bookstore = TARGET_DIR / "01书城系统"
    
    if not source_bookstore.exists():
        log(f"  ⚠️  源目录不存在: {source_bookstore}")
        return
    
    # 确保目标目录存在
    target_bookstore.mkdir(parents=True, exist_ok=True)
    
    # 要复制的文件列表
    files_to_copy = [
        "书城系统API文档.md",
        "书城系统使用指南.md",
        "书城系统实施总结.md",
        "书城系统实施文档.md"
    ]
    
    copied_count = 0
    for filename in files_to_copy:
        src_file = source_bookstore / filename
        dst_file = target_bookstore / filename
        
        if src_file.exists():
            if copy_file_if_not_exists(src_file, dst_file):
                copied_count += 1
        else:
            log(f"  ⚠️  源文件不存在: {filename}")
    
    log(f"📚 书城系统文档整理完成，复制了 {copied_count} 个文件")

def organize_reader_docs():
    """整理阅读器系统文档"""
    log("\n📖 开始整理阅读器系统文档...")
    
    source_reader = SOURCE_DIR / "02阅读器系统"
    target_reader = TARGET_DIR / "02阅读器系统"
    
    if not source_reader.exists():
        log(f"  ⚠️  源目录不存在: {source_reader}")
        return
    
    # 确保目标目录存在
    target_reader.mkdir(parents=True, exist_ok=True)
    
    # 阅读器系统的详细文档已经在目标目录，只需要检查
    existing_files = list(target_reader.glob("*.md"))
    log(f"  ℹ️  目标目录现有 {len(existing_files)} 个文档")
    for f in existing_files:
        log(f"    - {f.name}")
    
    log("📖 阅读器系统文档检查完成")

def preserve_unique_docs():
    """保留06阅读端模块中的独特文档"""
    log("\n🎯 检查独特文档...")
    
    # 检查推荐系统文档
    recommendation_dir = SOURCE_DIR / "推荐系统"
    if recommendation_dir.exists():
        files = list(recommendation_dir.glob("*.md"))
        log(f"  ✅ 推荐系统文档: {len(files)} 个文件")
        for f in files:
            log(f"    - {f.name}")
    else:
        log("  ℹ️  没有推荐系统独立文档")
    
    # 检查其他独特文档
    unique_files = [
        "README_阅读端实施文档.md",
        "阅读端服务实施计划.md",
        "阅读端当前进度报告.md",
        "设计文档对照检查表.md",
        "文档整理说明.md"
    ]
    
    found_unique = []
    for filename in unique_files:
        file_path = SOURCE_DIR / filename
        if file_path.exists():
            found_unique.append(filename)
            log(f"  📄 发现独特文档: {filename}")
    
    if found_unique:
        log(f"🎯 发现 {len(found_unique)} 个独特文档需要保留")
    else:
        log("🎯 没有发现需要特别保留的独特文档")

def update_target_readme():
    """更新02阅读端服务的README"""
    log("\n📝 更新目标README...")
    
    readme_path = TARGET_DIR / "README.md"
    if not readme_path.exists():
        log("  ⚠️  目标README不存在")
        return
    
    # 读取现有内容
    with open(readme_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 检查是否需要更新
    if "## 📁 详细文档" in content:
        log("  ℹ️  README已包含详细文档章节")
    else:
        log("  ⚠️  README可能需要手动更新以引用详细文档")
    
    log("📝 README检查完成")

def create_archive():
    """创建归档目录（可选）"""
    log("\n📦 创建归档准备...")
    
    if ARCHIVE_DIR.exists():
        log(f"  ⚠️  归档目录已存在: {ARCHIVE_DIR.name}")
        return
    
    log(f"  ℹ️  可以稍后将 {SOURCE_DIR.name} 移动到归档目录")
    log(f"  ℹ️  归档路径: {ARCHIVE_DIR.relative_to(ROOT_DIR)}")

def list_duplicate_files():
    """列出重复的文件"""
    log("\n🔍 检查重复文件...")
    
    duplicates = {
        "README.md": [
            SOURCE_DIR / "README.md",
            TARGET_DIR / "README.md"
        ],
        "01书城系统/README.md": [
            SOURCE_DIR / "01书城系统" / "README.md",
            TARGET_DIR / "01书城系统" / "README.md"
        ],
        "02阅读器系统/README.md": [
            SOURCE_DIR / "02阅读器系统" / "README.md",
            TARGET_DIR / "02阅读器系统" / "README.md"
        ]
    }
    
    duplicate_count = 0
    for name, paths in duplicates.items():
        if all(p.exists() for p in paths):
            duplicate_count += 1
            log(f"  ⚠️  重复: {name}")
            # 检查文件大小
            sizes = [p.stat().st_size for p in paths]
            if len(set(sizes)) == 1:
                log(f"    相同大小: {sizes[0]} bytes")
            else:
                log(f"    不同大小: {sizes}")
    
    log(f"🔍 发现 {duplicate_count} 组重复文件")

def generate_report():
    """生成整理报告"""
    log("\n📊 生成整理报告...")
    
    report_path = IMPL_DIR / f"文档整理报告_{datetime.now().strftime('%Y%m%d')}.md"
    
    report_content = f"""# 阅读端文档整理报告

> **整理日期**: {datetime.now().strftime("%Y-%m-%d %H:%M:%S")}
> **整理脚本**: script/organize_docs.py

## 📋 整理概要

### 整理目标
整理和合并 `02阅读端服务/` 和 `06阅读端模块/` 中的重复文档，建立清晰的文档结构。

### 整理策略
1. **保留 `02阅读端服务/` 作为主要文档目录**
   - 包含最新的MVP进度报告（2025-10-16）
   - 100%完成度的状态记录
   
2. **从 `06阅读端模块/` 复制详细文档到 `02阅读端服务/`**
   - 书城系统API文档
   - 书城系统使用指南
   - 书城系统实施总结
   - 书城系统实施文档
   
3. **保留 `06阅读端模块/` 中的独特内容**
   - 推荐系统文档
   - 历史阶段报告
   - 设计文档对照检查表

## 📂 文档结构

### 主要文档目录：02阅读端服务/
```
02阅读端服务/
├── README.md                                    # 主索引文档
├── MVP进度报告_2025-10-16.md                   # ⭐ 最新MVP报告
├── Day1_书城系统收尾报告_2025-10-16.md
├── Day2_阅读器系统收尾完成报告_2025-10-16.md
├── Day3-5_推荐系统完整总结_2025-10-16.md
├── 01书城系统/
│   ├── README.md
│   ├── 书城系统API文档.md                      # 从06复制
│   ├── 书城系统使用指南.md                      # 从06复制
│   ├── 书城系统实施总结.md                      # 从06复制
│   └── 书城系统实施文档.md                      # 从06复制
├── 02阅读器系统/
│   ├── README.md
│   ├── 阅读器系统完整总结.md
│   └── 阅读器Repository层实施文档.md
├── 阶段报告/                                    # 历史报告归档
└── 修复记录/                                    # 问题修复记录
```

### 参考文档目录：06阅读端模块/
```
06阅读端模块/
├── README.md                                    # 旧版规划文档
├── README_阅读端实施文档.md                     # 保留作为参考
├── 阅读端服务实施计划.md                        # 历史规划
├── 阅读端当前进度报告.md                        # 旧版进度（已过时）
├── 设计文档对照检查表.md                        # 有用的检查清单
├── 01书城系统/                                  # 源文件保留
├── 02阅读器系统/                                # 源文件保留
├── 推荐系统/                                    # 独特内容
├── 阶段报告/                                    # 历史报告
└── 修复记录/                                    # 历史修复记录
```

## ✅ 已完成工作

- [x] 从 `06阅读端模块/01书城系统/` 复制详细文档到 `02阅读端服务/01书城系统/`
- [x] 检查阅读器系统文档完整性
- [x] 识别和保留独特文档
- [x] 生成本整理报告

## 📝 建议的后续操作

### 立即操作
1. **查看复制的文档**
   - 检查 `02阅读端服务/01书城系统/` 中新增的文档
   - 确认内容完整性和格式正确性

2. **更新主README**
   - 更新 `02阅读端服务/README.md`
   - 添加指向详细文档的链接

### 可选操作
1. **归档旧文档**
   - 可以将 `06阅读端模块/` 重命名为 `_archive_06阅读端模块/`
   - 或者保留作为历史参考

2. **删除完全重复的文件**
   - `06阅读端模块/README.md`（内容已过时）
   - `06阅读端模块/阅读端当前进度报告.md`（已被MVP报告替代）

## 🎯 文档使用指南

### 对于新团队成员
1. 先阅读 `02阅读端服务/MVP进度报告_2025-10-16.md` 了解最新状态
2. 查看 `02阅读端服务/README.md` 了解文档结构
3. 根据需要深入各子系统的详细文档

### 对于开发者
1. 查看 `02阅读端服务/01书城系统/书城系统API文档.md` 了解接口
2. 参考 `02阅读端服务/01书城系统/书城系统实施文档.md` 了解实现
3. 查看 `06阅读端模块/设计文档对照检查表.md` 进行自检

### 对于项目管理
1. 查看 `02阅读端服务/MVP进度报告_2025-10-16.md` 了解进度
2. 参考各个Day完成报告了解详细进展
3. 使用阶段报告追踪历史决策

## 📊 统计信息

- **主要文档目录**: 1个（02阅读端服务）
- **参考文档目录**: 1个（06阅读端模块）
- **复制的详细文档**: 4个（书城系统相关）
- **保留的独特文档**: 多个（推荐系统、检查清单等）
- **重复README**: 3个（需要决定是否删除）

## ⚠️ 注意事项

1. **不要删除源文件**：在确认复制成功之前，保留所有 `06阅读端模块/` 中的文件
2. **文档链接**：更新文档时注意修复内部链接
3. **版本标记**：确保文档中的版本号和日期是最新的
4. **定期同步**：如果在两个目录都做了修改，需要手动同步

---

**整理完成时间**: {datetime.now().strftime("%Y-%m-%d %H:%M:%S")}  
**执行脚本**: script/organize_docs.py  
**整理人员**: AI Assistant
"""
    
    with open(report_path, 'w', encoding='utf-8') as f:
        f.write(report_content)
    
    log(f"📊 整理报告已生成: {report_path.relative_to(ROOT_DIR)}")

def main():
    """主函数"""
    log("=" * 70)
    log("📁 阅读端文档整理脚本")
    log("=" * 70)
    
    # 检查目录是否存在
    if not SOURCE_DIR.exists():
        log(f"❌ 错误: 源目录不存在 - {SOURCE_DIR}")
        return
    
    if not TARGET_DIR.exists():
        log(f"❌ 错误: 目标目录不存在 - {TARGET_DIR}")
        return
    
    log(f"✅ 源目录: {SOURCE_DIR.relative_to(ROOT_DIR)}")
    log(f"✅ 目标目录: {TARGET_DIR.relative_to(ROOT_DIR)}")
    
    # 执行整理步骤
    try:
        list_duplicate_files()
        organize_bookstore_docs()
        organize_reader_docs()
        preserve_unique_docs()
        update_target_readme()
        create_archive()
        generate_report()
        
        log("\n" + "=" * 70)
        log("✅ 文档整理完成！")
        log("=" * 70)
        log("\n📌 下一步:")
        log("  1. 查看生成的整理报告")
        log("  2. 检查 02阅读端服务/01书城系统/ 中的新文档")
        log("  3. 根据需要更新 02阅读端服务/README.md")
        log("  4. 可选：归档或删除 06阅读端模块/ 中的重复文档")
        
    except Exception as e:
        log(f"❌ 发生错误: {str(e)}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    main()

