#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
CNNovel125K 数据集导入脚本
从 Hugging Face 加载小说数据并转换为 JSON 格式供 Go 导入
"""

import json
import re
import os
from typing import List, Dict, Any
from datetime import datetime

try:
    from datasets import load_dataset
    print("[OK] datasets 库已加载")
except ImportError:
    print("[ERROR] 请先安装 datasets 库")
    print("运行命令: pip install datasets")
    exit(1)


class NovelImporter:
    """小说导入器"""

    def __init__(self, max_novels: int = 1000, chapter_size: int = 3000):
        """
        初始化导入器

        Args:
            max_novels: 最大导入小说数量
            chapter_size: 每章字数（默认3000字）
        """
        self.max_novels = max_novels
        self.chapter_size = chapter_size
        self.sensitive_words = self._load_sensitive_words()

    def _load_sensitive_words(self) -> set:
        """加载敏感词列表（简化版本，实际使用需要更完整的词库）"""
        # 这里是简化版本，实际应该使用完整的敏感词库
        return {
            "敏感词1", "敏感词2",  # 示例
        }

    def _clean_text(self, text: str) -> str:
        """
        清洗文本

        Args:
            text: 原始文本

        Returns:
            清洗后的文本
        """
        if not text:
            return ""

        # 移除多余空白
        text = re.sub(r'\s+', ' ', text)
        text = text.strip()

        # 简单的敏感词过滤（实际使用需要更完善的方案）
        for word in self.sensitive_words:
            if word in text:
                text = text.replace(word, "*" * len(word))

        return text

    def _split_into_chapters(self, content: str, book_title: str) -> List[Dict[str, Any]]:
        """
        将长文本拆分为章节

        Args:
            content: 书籍全文
            book_title: 书名

        Returns:
            章节列表
        """
        chapters = []

        # 尝试按章节标题分割
        chapter_pattern = r'第[零一二三四五六七八九十百千万\d]+章|第\d+节|Chapter\s*\d+'
        chapter_splits = re.split(f'({chapter_pattern})', content)

        if len(chapter_splits) > 3:
            # 找到了章节标记
            current_title = "序章"
            for i, part in enumerate(chapter_splits):
                if re.match(chapter_pattern, part):
                    current_title = part.strip()
                elif part.strip():
                    chapter_content = self._clean_text(part)
                    if chapter_content:
                        chapters.append({
                            "title": current_title,
                            "content": chapter_content,
                            "word_count": len(chapter_content),
                        })
        else:
            # 没有明显的章节标记，按字数分割
            cleaned_content = self._clean_text(content)
            total_length = len(cleaned_content)
            chapter_num = 1

            for start in range(0, total_length, self.chapter_size):
                end = min(start + self.chapter_size, total_length)
                chapter_content = cleaned_content[start:end]

                if chapter_content:
                    chapters.append({
                        "title": f"第{chapter_num}章",
                        "content": chapter_content,
                        "word_count": len(chapter_content),
                    })
                    chapter_num += 1

        return chapters

    def _map_category(self, original_category: str) -> str:
        """
        映射原始分类到系统分类

        Args:
            original_category: 原始分类

        Returns:
            系统分类
        """
        # 分类映射表
        category_mapping = {
            "玄幻": "玄幻",
            "奇幻": "奇幻",
            "武侠": "武侠",
            "仙侠": "仙侠",
            "都市": "都市",
            "言情": "言情",
            "历史": "历史",
            "军事": "军事",
            "科幻": "科幻",
            "灵异": "灵异",
            "游戏": "游戏",
            "竞技": "竞技",
            "同人": "同人",
            "轻小说": "轻小说",
        }

        for key in category_mapping:
            if key in original_category:
                return category_mapping[key]

        return "其他"

    def _validate_novel(self, novel: Dict[str, Any]) -> bool:
        """
        验证小说数据是否有效

        Args:
            novel: 小说数据

        Returns:
            是否有效
        """
        # 必须有标题和内容
        if not novel.get('title') or not novel.get('content'):
            return False

        # 标题长度验证
        if len(novel['title']) < 1 or len(novel['title']) > 100:
            return False

        # 内容长度验证（至少1000字）
        if len(novel.get('content', '')) < 1000:
            return False

        return True

    def _validate_novel_with_reason(self, novel: Dict[str, Any]) -> str:
        """
        验证小说数据并返回失败原因

        Args:
            novel: 小说数据

        Returns:
            验证失败原因，如果验证通过则返回空字符串
        """
        # 必须有标题和内容
        if not novel.get('title'):
            return "标题为空"
        if not novel.get('content'):
            return "内容为空"

        # 标题长度验证
        title_len = len(novel['title'])
        if title_len < 1:
            return "标题过短"
        if title_len > 100:
            return f"标题过长({title_len}字符)"

        # 内容长度验证（至少1000字）
        content_len = len(novel.get('content', ''))
        if content_len < 1000:
            return f"内容过短({content_len}字符，需要至少1000字符)"

        return ""

    def load_and_process(self) -> List[Dict[str, Any]]:
        """
        加载并处理数据集

        Returns:
            处理后的小说列表
        """
        print(f"开始加载 CNNovel125K 数据集（流式模式）...")
        print(f"最大导入数量: {self.max_novels} 本")

        try:
            # 使用流式加载，避免下载整个数据集
            ds = load_dataset("RyokoAI/CNNovel125K", split="train", streaming=True)
            print(f"[OK] 数据集连接成功，开始流式处理...")
        except Exception as e:
            print(f"[ERROR] 加载数据集失败: {e}")
            return []

        processed_novels = []
        processed_count = 0
        skipped_count = 0

        # 流式数据集，取出需要的数量（考虑可能有无效数据，多取一些）
        max_fetch = self.max_novels * 2  # 多取一倍以应对无效数据

        for i, item in enumerate(ds.take(max_fetch)):
            if processed_count >= self.max_novels:
                break

            # 第一条数据时打印字段信息
            if i == 0:
                print(f"\n数据字段: {list(item.keys())}")
                meta = item.get('meta', {})
                print(f"meta字段: {list(meta.keys()) if isinstance(meta, dict) else 'not a dict'}")

            # 提取字段 - 数据在 meta 和 text 中
            meta = item.get('meta', {}) if isinstance(item.get('meta'), dict) else {}
            novel_data = {
                'title': meta.get('title', meta.get('书名', '')).strip(),
                'author': meta.get('author', meta.get('作者', '佚名')).strip(),
                'content': item.get('text', ''),
                'category': meta.get('category', meta.get('分类', '其他')),
                'word_count': meta.get('word_count', meta.get('字数', 0)),
                'rating': float(meta.get('rating', meta.get('评分', 0.0) or 0.0)),
            }

            # 验证数据
            validation_error = self._validate_novel_with_reason(novel_data)
            if validation_error:
                if i < 5:  # 只打印前5个错误
                    print(f"  跳过第{i+1}本: {validation_error}")
                skipped_count += 1
                continue

            # 清洗标题和作者
            novel_data['title'] = self._clean_text(novel_data['title'])
            novel_data['author'] = self._clean_text(novel_data['author'])

            # 映射分类
            novel_data['category'] = self._map_category(novel_data['category'])

            # 拆分章节
            chapters = self._split_into_chapters(novel_data['content'], novel_data['title'])

            if not chapters:
                skipped_count += 1
                continue

            # 计算实际字数
            actual_word_count = sum(ch['word_count'] for ch in chapters)

            # 生成简介（取前500字）
            introduction = novel_data['content'][:500] + "..."
            introduction = self._clean_text(introduction)

            # 构建最终数据
            processed_novel = {
                'title': novel_data['title'],
                'author': novel_data['author'],
                'introduction': introduction,
                'category': novel_data['category'],
                'word_count': actual_word_count,
                'chapter_count': len(chapters),
                'rating': float(novel_data.get('rating', 0.0)),
                'status': 'completed',  # 全本小说标记为已完结
                'is_free': True,  # 默认免费
                'chapters': chapters,
            }

            processed_novels.append(processed_novel)
            processed_count += 1

            # 每处理100本输出一次进度
            if processed_count % 100 == 0:
                print(f"已处理 {processed_count} 本小说...")

        print(f"\n处理完成:")
        print(f"  [OK] 成功处理: {processed_count} 本")
        print(f"  [SKIP] 跳过无效: {skipped_count} 本")
        print(f"  总章节数: {sum(novel['chapter_count'] for novel in processed_novels)}")

        return processed_novels

    def save_to_json(self, novels: List[Dict[str, Any]], output_file: str):
        """
        保存为 JSON 文件

        Args:
            novels: 小说列表
            output_file: 输出文件路径
        """
        # 确保输出目录存在
        os.makedirs(os.path.dirname(output_file) if os.path.dirname(output_file) else '.', exist_ok=True)

        # 添加元数据
        output_data = {
            'metadata': {
                'source': 'CNNovel125K',
                'total_novels': len(novels),
                'total_chapters': sum(novel['chapter_count'] for novel in novels),
                'generated_at': datetime.now().isoformat() + 'Z',  # 添加UTC时区标识
                'chapter_size': self.chapter_size,
            },
            'novels': novels,
        }

        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump(output_data, f, ensure_ascii=False, indent=2)

        print(f"\n[OK] 数据已保存到: {output_file}")
        print(f"  文件大小: {os.path.getsize(output_file) / 1024 / 1024:.2f} MB")


def main():
    """主函数"""
    import argparse

    parser = argparse.ArgumentParser(description='CNNovel125K 数据集导入脚本')
    parser.add_argument('--max-novels', type=int, default=1000, help='最大导入小说数量（默认1000）')
    parser.add_argument('--chapter-size', type=int, default=3000, help='每章字数（默认3000）')
    parser.add_argument('--output', type=str, default='data/novels.json', help='输出文件路径')

    args = parser.parse_args()

    # 创建导入器
    importer = NovelImporter(
        max_novels=args.max_novels,
        chapter_size=args.chapter_size
    )

    # 加载和处理数据
    novels = importer.load_and_process()

    if not novels:
        print("[ERROR] 没有处理任何小说数据")
        return

    # 保存为 JSON
    importer.save_to_json(novels, args.output)

    print("\n数据处理完成！")
    print(f"下一步: 使用 Go 程序导入数据到 MongoDB")
    print(f"运行命令: go run cmd/migrate/main.go --action=import-novels --file={args.output}")


if __name__ == '__main__':
    main()

