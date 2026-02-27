#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
青羽写作平台 API - 书城示例
演示如何使用Python调用书城相关API
"""

import requests
import json
from typing import Optional, Dict, Any, List


class BookstoreAPIClient:
    """书城API客户端"""

    def __init__(self, base_url: str = "http://localhost:9090/api/v1", token: Optional[str] = None):
        self.base_url = base_url.rstrip('/')
        self.token = token
        self.session = requests.Session()

    def _request(
        self,
        method: str,
        endpoint: str,
        data: Optional[Dict[str, Any]] = None,
        params: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """发送HTTP请求"""
        url = f"{self.base_url}/{endpoint.lstrip('/')}"
        headers = {"Content-Type": "application/json"}

        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"

        try:
            response = self.session.request(
                method=method,
                url=url,
                json=data,
                params=params,
                headers=headers,
                timeout=30
            )
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            print(f"请求失败: {e}")
            if hasattr(e.response, 'text'):
                print(f"响应内容: {e.response.text}")
            raise

    def get_home_data(self) -> Dict[str, Any]:
        """获取首页数据"""
        print("\n=== 获取首页数据 ===")
        result = self._request("GET", "/bookstore/home")
        print(f"首页数据: {json.dumps(result, indent=2, ensure_ascii=False)}")
        return result

    def get_books(self, page: int = 1, limit: int = 10) -> Dict[str, Any]:
        """获取书籍列表"""
        print(f"\n=== 获取书籍列表 (第{page}页) ===")
        result = self._request("GET", "/bookstore/books", params={"page": page, "limit": limit})
        if "data" in result and "books" in result["data"]:
            print(f"获取到 {len(result['data']['books'])} 本书籍")
        return result

    def search_by_title(self, title: str, page: int = 1, limit: int = 5) -> Dict[str, Any]:
        """按标题搜索书籍"""
        print(f"\n=== 搜索书籍（标题: {title}） ===")
        result = self._request(
            "GET",
            "/bookstore/books/search/title",
            params={"title": title, "page": page, "limit": limit}
        )
        if "data" in result and "books" in result["data"]:
            print(f"找到 {len(result['data']['books'])} 本相关书籍")
        return result

    def search_by_author(self, author: str, page: int = 1, limit: int = 5) -> Dict[str, Any]:
        """按作者搜索"""
        print(f"\n=== 搜索作者（作者: {author}） ===")
        result = self._request(
            "GET",
            "/bookstore/books/search/author",
            params={"author": author, "page": page, "limit": limit}
        )
        if "data" in result and "books" in result["data"]:
            print(f"找到 {len(result['data']['books'])} 本相关书籍")
        return result

    def filter_by_tags(self, tags: List[str], page: int = 1, limit: int = 5) -> Dict[str, Any]:
        """按标签筛选"""
        print(f"\n=== 按标签筛选（标签: {', '.join(tags)}） ===")
        result = self._request(
            "GET",
            "/bookstore/books/tags",
            params={"tags": ",".join(tags), "page": page, "limit": limit}
        )
        if "data" in result and "books" in result["data"]:
            print(f"找到 {len(result['data']['books'])} 本相关书籍")
        return result

    def get_book_detail(self, book_id: str) -> Dict[str, Any]:
        """获取书籍详情"""
        print(f"\n=== 获取书籍详情 (ID: {book_id}) ===")
        result = self._request("GET", f"/bookstore/books/{book_id}")
        print(f"书籍标题: {result.get('data', {}).get('title', 'N/A')}")
        return result

    def get_book_chapters(self, book_id: str, page: int = 1, limit: int = 20) -> Dict[str, Any]:
        """获取书籍章节列表"""
        print(f"\n=== 获取章节列表 (书籍ID: {book_id}) ===")
        result = self._request(
            "GET",
            f"/bookstore/books/{book_id}/chapters",
            params={"page": page, "limit": limit}
        )
        if "data" in result and "chapters" in result["data"]:
            print(f"获取到 {len(result['data']['chapters'])} 个章节")
        return result

    def get_chapter_content(self, chapter_id: str) -> Dict[str, Any]:
        """获取章节内容（需要认证）"""
        print(f"\n=== 获取章节内容 (ID: {chapter_id}) ===")
        result = self._request("GET", f"/bookstore/chapters/{chapter_id}")
        if "data" in result:
            content = result["data"].get("content", "")
            print(f"章节内容预览: {content[:100]}...")
        return result

    def get_similar_books(self, book_id: str, limit: int = 10) -> Dict[str, Any]:
        """获取相似书籍推荐"""
        print(f"\n=== 获取相似书籍 (书籍ID: {book_id}) ===")
        result = self._request(
            "GET",
            f"/bookstore/books/{book_id}/similar",
            params={"limit": limit}
        )
        if "data" in result and "books" in result["data"]:
            print(f"找到 {len(result['data']['books'])} 本相似书籍")
        return result

    def rate_book(self, book_id: str, rating: int, comment: str = "") -> Dict[str, Any]:
        """提交书籍评分（需要认证）"""
        print(f"\n=== 提交评分 (书籍ID: {book_id}, 评分: {rating}) ===")
        data = {"rating": rating}
        if comment:
            data["comment"] = comment
        result = self._request("POST", f"/bookstore/books/{book_id}/rating", data=data)
        print("评分提交成功")
        return result


def main():
    """主函数 - 演示书城API"""

    # 从环境变量或命令行获取token
    import os
    token = os.getenv("QINGYU_TOKEN", "your-jwt-token-here")

    if token == "your-jwt-token-here":
        print("警告: 未设置Token，部分API将无法访问")
        print("请设置环境变量: export QINGYU_TOKEN=your_token")
        token = None

    # 创建API客户端
    client = BookstoreAPIClient(token=token)

    try:
        # 1. 获取首页数据
        client.get_home_data()

        # 2. 获取书籍列表
        client.get_books(page=1, limit=5)

        # 3. 搜索书籍
        client.search_by_title("玄幻")

        # 4. 按作者搜索
        client.search_by_author("唐家三少")

        # 5. 按标签筛选
        client.filter_by_tags(["玄幻", "修真"])

        # 6. 获取书籍详情
        book_id = "book123"
        client.get_book_detail(book_id)

        # 7. 获取章节列表
        client.get_book_chapters(book_id)

        # 8. 获取章节内容（需要认证）
        if token:
            chapter_id = "chapter123"
            client.get_chapter_content(chapter_id)

        # 9. 获取相似书籍
        client.get_similar_books(book_id)

        # 10. 提交评分（需要认证）
        if token:
            client.rate_book(book_id, rating=5, comment="非常好看！")

        print("\n=== 书城API示例完成 ===")

    except Exception as e:
        print(f"\n错误: {e}")
        return 1

    return 0


if __name__ == "__main__":
    exit(main())
