#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Swagger 注释自动生成脚本

功能：
1. 扫描 router/*/_router.go 文件
2. 解析 Gin 路由定义
3. 在对应 Handler 函数前插入 Swagger 注释模板

使用：
    python scripts/add-swagger-annotations.py --module reader
    python scripts/add-swagger-annotations.py --all
"""

import os
import re
import sys
import argparse
from pathlib import Path
from typing import List, Dict, Tuple

# 设置标准输出编码为UTF-8（Windows兼容）
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

# 项目根目录
ROOT_DIR = Path(__file__).parent.parent
ROUTER_DIR = ROOT_DIR / "router"
API_DIR = ROOT_DIR / "api/v1"


def main():
    parser = argparse.ArgumentParser(description="生成 Swagger 注释模板")
    parser.add_argument("--module", type=str, help="指定模块名称（如 reader）")
    parser.add_argument("--all", action="store_true", help="处理所有模块")
    parser.add_argument("--dry-run", action="store_true", help="预览模式，不修改文件")
    parser.add_argument("--verbose", action="store_true", help="详细输出")

    args = parser.parse_args()

    if args.all:
        modules = ["reader", "admin", "bookstore", "writer", "social", "ai", "user", "shared"]
    elif args.module:
        modules = [args.module]
    else:
        parser.print_help()
        sys.exit(1)

    for module in modules:
        process_module(module, args.dry_run, args.verbose)


def parse_gin_routes(router_file: Path) -> List[Dict]:
    """
    解析 Gin 路由文件，提取路由信息

    返回格式：
    [
        {
            "method": "GET",
            "path": "/books/:bookId",
            "handler": "GetBooks",
            "line": 42
        },
        ...
    ]
    """
    routes = []
    content = router_file.read_text(encoding="utf-8")

    # 匹配路由定义，例如：
    # readerGroup.GET("/books/:bookId", readerApiHandler.GetBooks)
    pattern = re.compile(
        r'(\w+)\.(GET|POST|PUT|DELETE|PATCH)\(["\']([^"\']+)["\'],\s*\w+\.(\w+)\)',
        re.MULTILINE
    )

    for match in pattern.finditer(content):
        group_name, method, path, handler_name = match.groups()
        routes.append({
            "group": group_name,
            "method": method.upper(),
            "path": path,
            "handler": handler_name,
            "line": content[:match.start()].count('\n') + 1
        })

    return routes


def extract_path_params(path: str) -> List[str]:
    """提取路径参数，例如 /books/:bookId -> ["bookId"]"""
    return re.findall(r':(\w+)', path)


def gin_to_openapi_path(path: str) -> str:
    """
    转换 Gin 风格路径为 OpenAPI 风格
    例如：/books/:bookId -> /books/{bookId}
    """
    return re.sub(r':(\w+)', r'{\1}', path)


def generate_swagger_comment(route: Dict, module_name: str, requires_auth: bool = True) -> str:
    """
    生成 Swagger 注释模板

    Args:
        route: 路由信息字典
        module_name: 模块名称
        requires_auth: 是否需要鉴权（默认 True）

    Returns:
        Swagger 注释字符串
    """
    method = route["method"]
    path = route["path"]
    handler_name = route["handler"]
    path_params = extract_path_params(path)
    openapi_path = gin_to_openapi_path(path)

    # 从 handler 名推断 summary
    summary = infer_summary_from_handler(handler_name)

    comment_lines = [
        f"// @Summary {summary}",
        "// @Description TODO: 补充详细描述",
        f"// @Tags {module_name}",
        "// @Accept json",
        "// @Produce json",
    ]

    # 添加鉴权标注（关键：与项目 JWTAuth 中间件对齐）
    if requires_auth:
        comment_lines.append("// @Security Bearer")

    # 添加路径参数
    for param in path_params:
        param_name = param[0].upper() + param[1:]  # bookId -> BookId
        comment_lines.append(f"// @Param {param} path string true \"{param_name}\"")

    # 添加响应
    comment_lines.append(f"// @Success 200 {{object}} response.APIResponse{{data=TODO}}")
    comment_lines.append(f"// @Failure 400 {{object}} response.APIResponse")

    # 添加路由（关键：使用 OpenAPI 标准 {param}，禁止 Gin 风格的 :param）
    comment_lines.append(f"// @Router /{module_name}{openapi_path} [{method.lower()}]")

    return "\n".join(comment_lines) + "\n"


def infer_summary_from_handler(handler_name: str) -> str:
    """从 Handler 函数名推断 Summary"""
    # 移除常见前缀
    name = handler_name
    for prefix in ["Get", "Create", "Update", "Delete", "List", "Fetch"]:
        if name.startswith(prefix):
            name = name[len(prefix):]
            break

    # 转换为中文（简单映射，实际使用时可扩展）
    summaries = {
        "ChapterContent": "获取章节内容",
        "Books": "获取书籍列表",
        "BookDetail": "获取书籍详情",
        "PublicCollections": "获取公开书架",
        "RecentReading": "获取最近阅读",
        "UnfinishedBooks": "获取未完成书籍",
        "FinishedBooks": "获取已完成书籍",
        "AddToBookshelf": "添加到书架",
        "RemoveFromBookshelf": "从书架移除",
        "UpdateBookStatus": "更新书籍状态",
        "BatchUpdateBookStatus": "批量更新书籍状态",
        "LikeBook": "点赞书籍",
        "UnlikeBook": "取消点赞",
        "BookLikeInfo": "获取书籍点赞信息",
        "ChapterByNumber": "按章节号获取章节",
        "NextChapter": "获取下一章",
        "PreviousChapter": "获取上一章",
        # TODO: 添加更多映射
    }

    return summaries.get(name, f"{handler_name} 操作")


def find_handler_file(module_name: str, handler_name: str) -> Path:
    """
    查找 Handler 函数所在的文件（健壮版本）

    查找策略：
    1. 扫描 api/v1/{module_name}/**/*.go
    2. 匹配函数签名：func (xxx *YYY) HandlerName(
    3. 支持多文件拆分、多 API struct、不同 handler 命名文件

    Returns:
        Path to handler file, or None if not found
    """
    module_dir = API_DIR / module_name

    if not module_dir.exists():
        return None

    # 递归扫描模块目录下所有 .go 文件
    for go_file in module_dir.rglob("*.go"):
        # 跳过测试文件和生成的文件
        if "_test.go" in go_file.name or go_file.name.startswith("mock_"):
            continue

        content = go_file.read_text(encoding="utf-8")

        # 匹配 Handler 函数签名（支持接收者）
        # 例如：func (h *Handler) GetBooks( 或 func (r *ReaderAPI) GetChapterContent(
        pattern = rf'func\s+\([^)]*\*\s*\w+\)\s+{re.escape(handler_name)}\('
        if re.search(pattern, content):
            return go_file

    return None


def process_module(module_name: str, dry_run: bool, verbose: bool):
    """处理单个模块的路由"""
    print(f"\n{'='*60}")
    print(f"处理模块: {module_name}")
    print(f"{'='*60}")

    router_file = ROUTER_DIR / module_name / f"{module_name}_router.go"

    if not router_file.exists():
        print(f"❌ 路由文件不存在: {router_file}")
        return

    routes = parse_gin_routes(router_file)

    if not routes:
        print(f"⚠️  未找到任何路由定义")
        return

    print(f"✓ 找到 {len(routes)} 个路由定义")

    # 生成 Swagger 注释
    generated_count = 0
    skipped_count = 0

    for route in routes:
        handler_file = find_handler_file(module_name, route["handler"])

        if not handler_file:
            if verbose:
                print(f"  ⚠️  未找到 Handler {route['handler']} 的文件")
            skipped_count += 1
            continue

        # 检查是否已有 Swagger 注释
        content = handler_file.read_text(encoding="utf-8")
        handler_pattern = rf'func \(.*?\) {route["handler"]}\('
        handler_match = re.search(handler_pattern, content)

        if not handler_match:
            if verbose:
                print(f"  ⚠️  未找到 Handler 函数 {route['handler']}")
            skipped_count += 1
            continue

        # 检查函数前是否已有 @Summary
        handler_start = handler_match.start()
        preceding_text = content[max(0, handler_start - 500):handler_start]

        if "@Summary" in preceding_text:
            if verbose:
                print(f"  ⊙ 跳过 {route['handler']} (已有注释)")
            skipped_count += 1
            continue

        # 生成注释
        comment = generate_swagger_comment(route, module_name)

        if dry_run:
            print(f"  [DRY RUN] 将在 {handler_file.name} 中为 {route['handler']} 添加注释")
            if verbose:
                print(f"    {route['method']} {route['path']}")
        else:
            # 插入注释
            new_content = (
                content[:handler_start] +
                comment +
                "\n" +
                content[handler_start:]
            )
            handler_file.write_text(new_content, encoding="utf-8")
            generated_count += 1

            if verbose:
                print(f"  ✓ 为 {route['handler']} 添加注释")

    print(f"\n✓ 处理完成:")
    print(f"  - 新增: {generated_count}")
    print(f"  - 跳过: {skipped_count}")

    if generated_count > 0:
        print(f"\n⚠️  请手动补充以下内容:")
        print(f"  - 详细的 @Description")
        print(f"  - 请求参数的具体字段 (query/body)")
        print(f"  - @Success 和 @Failure 的具体类型")


if __name__ == "__main__":
    main()
