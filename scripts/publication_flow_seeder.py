#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Complete Publication Flow Seeder
Creates test data through the full publication workflow:
- Creates test authors
- Creates writer projects
- Creates chapter documents
- Submits for publication
- Admin approves
- Converts to bookstore books with chapters
"""

import argparse
import json
import sys
import urllib.request
import urllib.error
import urllib.parse
import time
import hashlib
import os

# Set UTF-8 encoding for Windows console
if sys.platform == 'win32':
    import codecs
    sys.stdout = codecs.getwriter('utf-8')(sys.stdout.buffer, errors='replace')
    sys.stderr = codecs.getwriter('utf-8')(sys.stderr.buffer, errors='replace')

# Default configuration
DEFAULT_BASE_URL = "http://localhost:9090"
DEFAULT_ADMIN_USERNAME = "testadmin001"
DEFAULT_ADMIN_PASSWORD = "password"

# Book configurations
BOOKS_PER_AUTHOR = 5  # Books per author
CHAPTERS_PER_BOOK = 3  # Chapters per book

# Test author configurations
TEST_AUTHORS = [
    {"username": "hot_author_01", "password": "Author@123456", "nickname": "HotAuthor01"},
    {"username": "hot_author_02", "password": "Author@123456", "nickname": "HotAuthor02"},
]


def build_url(base_url: str, path: str, query: dict | None = None) -> str:
    """Build URL with query parameters"""
    # Ensure base_url doesn't end with slash and path doesn't start with slash
    base_url = base_url.rstrip('/')
    path = path.lstrip('/')
    url = f"{base_url}/{path}"
    if query:
        url = f"{url}?{urllib.parse.urlencode(query)}"
    return url


def request_json(method: str, url: str, token: str | None = None, body: dict | None = None) -> dict | None:
    """Make HTTP request and return JSON response"""
    headers = {"Accept": "application/json"}
    data = None

    if token:
        headers["Authorization"] = f"Bearer {token}"
    if body:
        data = json.dumps(body).encode("utf-8")
        headers["Content-Type"] = "application/json"

    req = urllib.request.Request(url=url, data=data, headers=headers, method=method)

    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            payload = resp.read().decode("utf-8")
    except urllib.error.HTTPError as exc:
        payload = exc.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"{method} {url} failed: {exc.code} {payload}") from exc
    except urllib.error.URLError as exc:
        raise RuntimeError(f"{method} {url} failed: {exc}") from exc

    if not payload:
        return None

    try:
        return json.loads(payload)
    except json.JSONDecodeError:
        return None


def login(base_url: str, username: str, password: str) -> tuple[str, dict]:
    """Login and return user info"""
    response = request_json(
        "POST",
        build_url(base_url, "/api/v1/login"),
        body={"username": username, "password": password},
    )
    if not response:
        raise RuntimeError(f"login failed for {username}: no response")

    if response.get("code") != 0:
        raise RuntimeError(f"login failed for {username}: {response.get('message')}")

    data = response.get("data", {})
    if not isinstance(data, dict) or not data.get("token"):
        raise RuntimeError(f"login failed for {username}: no token in response")

    return data["token"], data.get("user", {})


def get_admin_token(base_url: str) -> str:
    """Get admin token using default credentials"""
    token, _ = login(base_url, DEFAULT_ADMIN_USERNAME, DEFAULT_ADMIN_PASSWORD)
    return token


def create_or_get_writer_project(
    base_url: str,
    author_token: str,
    author_id: str,
    title: str,
    introduction: str,
    category: str,
) -> dict:
    """Create or get a writer project"""
    # Try to find existing project
    projects_url = build_url(base_url, "/api/v1/writer/projects", {"page": 1, "pageSize": 50})
    projects_resp = request_json("GET", projects_url, token=author_token)

    if projects_resp and projects_resp.get("code") == 0:
        projects = projects_resp.get("data", {}).get("items", [])
        # Find project by title
        for project in projects:
            if project.get("title") == title:
                return project

    # Create new project
    create_url = build_url(base_url, "/api/v1/writer/projects")
    create_resp = request_json(
        "POST",
        create_url,
        token=author_token,
        body={
            "title": title,
            "summary": introduction,
            "writingType": "novel",
            "category": category,
        },
    )

    if create_resp and create_resp.get("code") == 0:
        return create_resp.get("data", {})
    return {}


def create_chapter_document(
    base_url: str,
    author_token: str,
    project_id: str,
    chapter_num: int,
    title: str,
    content: str,
    word_count: int,
) -> dict:
    """Create a chapter document"""
    doc_url = build_url(base_url, "/api/v1/writer/documents")
    create_resp = request_json(
        "POST",
        doc_url,
        token=author_token,
        body={
            "projectId": project_id,
            "title": title,
            "type": "chapter",
            "chapterNumber": chapter_num,
            "content": content,
            "wordCount": word_count,
            "status": "published",
        },
    )

    if create_resp and create_resp.get("code") == 0:
        return create_resp.get("data", {})
    return {}


def submit_project_publication(
    base_url: str,
    author_token: str,
    project_id: str,
) -> dict:
    """Submit project for publication"""
    pub_url = build_url(base_url, f"/api/v1/writer/projects/{project_id}/publish")
    pub_resp = request_json(
        "POST",
        pub_url,
        token=author_token,
        body={
            "publishType": "project",
            "categoryIds": [],
            "isFree": True,
            "price": 0,
        },
    )

    if pub_resp and pub_resp.get("code") == 0:
        return pub_resp.get("data", {})
    return {}


def submit_document_publication(
    base_url: str,
    author_token: str,
    document_id: str,
    chapter_number: int,
) -> dict:
    """Submit document for publication"""
    pub_url = build_url(base_url, f"/api/v1/writer/documents/{document_id}/publish")
    pub_resp = request_json(
        "POST",
        pub_url,
        token=author_token,
        body={
            "chapterNumber": chapter_number,
            "isFree": True,
            "price": 0,
        },
    )

    if pub_resp and pub_resp.get("code") == 0:
        return pub_resp.get("data", {})
    return {}


def approve_publication(
    base_url: str,
    admin_token: str,
    publication_id: str,
) -> dict:
    """Approve a publication"""
    review_url = build_url(base_url, f"/api/v1/admin/publications/{publication_id}/review")
    review_resp = request_json(
        "POST",
        review_url,
        token=admin_token,
        body={
            "action": "approve",
            "note": "Approved by publication flow seeder",
        },
    )

    if review_resp and review_resp.get("code") == 0:
        return review_resp.get("data", {})
    return {}


def get_bookstore_book(base_url: str, book_id: str) -> dict:
    """Get bookstore book by ID"""
    book_url = build_url(base_url, f"/api/v1/bookstore/books/{book_id}/detail")
    book_resp = request_json("GET", book_url)

    if book_resp and book_resp.get("code") == 0:
        return book_resp.get("data", {})
    return {}


def get_bookstore_chapters(base_url: str, book_id: str) -> list:
    """Get bookstore chapters"""
    chapters_url = build_url(base_url, f"/api/v1/bookstore/books/{book_id}/chapters", {"page": 1, "size": 20})
    chapters_resp = request_json("GET", chapters_url)

    if chapters_resp and chapters_resp.get("code") == 0:
        return chapters_resp.get("data", {}).get("items", [])
    return []


def generate_book_title(author_index: int, book_index: int) -> str:
    """Generate book title based on author and book index"""
    titles = [
        ["TianDao", "ZhiZun", "ShenJi", "XianDi", "WuShen"],
        ["CangQiong", "XingHe", "QianKun", "HunDun", "XuKong"],
        ["ChuanShuo", "JiYuan", "ShiDai", "WangChao", "DiGuo"],
    ]
    title_prefixes = titles[author_index % len(titles)]
    suffixes = ["Lu", "Zhuan", "Ji", "Zhi", "Shi"]
    return f"{title_prefixes[book_index % len(title_prefixes)]}{suffixes[book_index % len(suffixes)]}"


def generate_chapter_content(chapter_num: int, book_title: str) -> str:
    """Generate chapter content"""
    return f"""# Chapter {chapter_num}: {'The Beginning' if chapter_num == 1 else 'Rising Power' if chapter_num == 2 else 'The Climax'}

This is the world of {book_title}, a place full of miracles.

The protagonist Lin Feng stood on the mountain peak, looking at the vast land below.
His eyes shone with determination, for he knew that his era was about to begin.

"From today on, I will embark on the path of cultivation!"

Lin Feng took a deep breath, feeling the abundant spiritual energy between heaven and earth.
This energy surged into his body like a tide, flowing along his meridians,
finally gathering in his dantian.

This is the first step of immortal cultivation, and also the most important step.
Only by laying a solid foundation can one go further on the path of cultivation.

---
(This chapter is about 800 words, for testing purposes)
"""


def generate_test_comments(
    base_url: str,
    reader_token: str,
    book_id: str,
    num_comments: int,
) -> list:
    """Generate test comments for a book"""
    comments = []
    comment_templates = [
        "This book is amazing! Could not stop reading!",
        "Author please update faster! Waiting so hard!",
        "The plot is full of twists, characters are well developed, highly recommended!",
        "This is the best book I've read this year, bar none!",
        "Five stars! The story is captivating and addictive!",
    ]

    for i in range(num_comments):
        content = comment_templates[i % len(comment_templates)]
        rating = 4 + (i % 2)  # 4 or 5 stars

        comment_url = build_url(base_url, "/api/v1/reader/comments")
        try:
            comment_resp = request_json(
                "POST",
                comment_url,
                token=reader_token,
                body={
                    "bookId": book_id,
                    "content": content,
                    "rating": rating,
                },
            )

            if comment_resp and comment_resp.get("code") == 0:
                comments.append(comment_resp.get("data", {}))
        except Exception as e:
            print(f"      Warning: Comment creation failed: {e}")

    return comments


def seed_publication_flow(
    base_url: str = None,
    num_books_per_author: int = 5,
    chapters_per_book: int = 3,
):
    """Main function to seed publication flow data"""
    if base_url is None:
        base_url = DEFAULT_BASE_URL

    print(f"Starting publication flow seeder...")
    print(f"  Base URL: {base_url}")
    print(f"  Books per author: {num_books_per_author}")
    print(f"  Chapters per book: {chapters_per_book}")

    # 1. Get admin token
    print("\nStep 1: Getting admin token...")
    try:
        admin_token = get_admin_token(base_url)
        print(f"  [OK] Admin token obtained")
    except Exception as e:
        print(f"  [ERROR] Failed to get admin token: {e}")
        return

    # 2. Create test authors
    print("\nStep 2: Creating test authors...")
    test_authors = []
    reader_tokens = []

    for author_config in TEST_AUTHORS:
        print(f"  Creating author: {author_config['username']}")
        author_token = None
        author_id = None

        # Try to login first (user might already exist)
        try:
            login_url = build_url(base_url, "/api/v1/login")
            login_resp = request_json(
                "POST",
                login_url,
                body={
                    "username": author_config["username"],
                    "password": author_config["password"],
                },
            )

            if login_resp and login_resp.get("code") == 0:
                author_token = login_resp.get("data", {}).get("token")
                author_id = login_resp.get("data", {}).get("user", {}).get("user_id")
                print(f"    [OK] Login successful")
        except Exception:
            pass

        # If login failed, try to register (without email verification for test)
        if not author_token:
            try:
                # Use test registration endpoint if available, or direct DB creation
                print(f"    [INFO] User does not exist, will use existing test authors")
            except Exception as e:
                print(f"    [WARN] Could not create author: {e}")
                continue

        if author_token and author_id:
            test_authors.append({
                "username": author_config["username"],
                "nickname": author_config["nickname"],
                "token": author_token,
                "id": author_id,
            })
            reader_tokens.append(author_token)

    print(f"  [OK] {len(test_authors)} test authors ready")

    if len(test_authors) == 0:
        print("\n[WARN] No test authors available. Using existing author users from seeder...")
        # Try to use existing author users from seeder
        try:
            # Login as testauthor001 (created by user seeder)
            token, user = login(base_url, "testauthor001", "password")
            test_authors.append({
                "username": "testauthor001",
                "nickname": "TestAuthor01",
                "token": token,
                "id": user.get("user_id"),
            })
            reader_tokens.append(token)
            print(f"  [OK] Using testauthor001 as author")
        except Exception as e:
            print(f"  [ERROR] No authors available: {e}")
            return

    # 3. Create reader user for comments
    print("\nStep 3: Getting reader user...")
    reader_token = None

    try:
        token, _ = login(base_url, "testuser001", "password")
        reader_token = token
        print(f"  [OK] Reader user ready")
    except Exception as e:
        print(f"  [WARN] Could not get reader user: {e}")

    if reader_token:
        reader_tokens.append(reader_token)

    # 4. For each author, create books through publication flow
    print("\nStep 4: Creating books through publication flow...")
    categories = ["XuanHuan", "DuShi", "XianXia", "KeHuan", "LiShi", "WuXia", "YouXi", "QiHuan"]
    created_books = []

    for author_index, author in enumerate(test_authors):
        print(f"\n  === Author: {author['username']} ===")

        for book_i in range(num_books_per_author):
            # Determine category
            category = categories[(author_index * len(categories) + book_i) % len(categories)]

            # Generate book title
            book_title = generate_book_title(author_index, book_i)
            print(f"\n  Book {book_i + 1}: {book_title}")

            try:
                # 4.1 Create writer project
                project = create_or_get_writer_project(
                    base_url,
                    author["token"],
                    author["id"],
                    book_title,
                    f"A wonderful story about {category}.",
                    category,
                )

                if not project or not project.get("id"):
                    print(f"    [WARN] Failed to create project")
                    continue

                project_id = project["id"]
                print(f"    [OK] Project created: {project_id}")

                # 4.2 Create chapter documents
                chapter_ids = []
                for chapter_i in range(1, chapters_per_book + 1):
                    chapter_title = f"Chapter {chapter_i}: {'The Beginning' if chapter_i == 1 else 'Rising Power' if chapter_i == 2 else 'The Climax'}"
                    content = generate_chapter_content(chapter_i, book_title)

                    doc = create_chapter_document(
                        base_url,
                        author["token"],
                        project_id,
                        chapter_i,
                        chapter_title,
                        content,
                        800 + chapter_i * 100,
                    )

                    if doc and doc.get("id"):
                        chapter_ids.append(doc["id"])
                        print(f"      [OK] Chapter {chapter_i} created: {doc['id']}")

                if len(chapter_ids) == 0:
                    print(f"    [WARN] No chapters created")
                    continue

                # 4.3 Submit project publication
                print(f"    Submitting project publication...")
                pub_record = submit_project_publication(
                    base_url,
                    author["token"],
                    project_id,
                )

                if not pub_record or not pub_record.get("id"):
                    print(f"    [WARN] Failed to submit project publication")
                    continue

                pub_record_id = pub_record["id"]
                print(f"    [OK] Publication submitted: {pub_record_id}")

                # 4.4 Submit document publications
                doc_pub_ids = []
                for chapter_i, doc_id in enumerate(chapter_ids):
                    chapter_num = chapter_i + 1
                    print(f"      Submitting chapter {chapter_num} publication...")
                    doc_pub_record = submit_document_publication(
                        base_url,
                        author["token"],
                        doc_id,
                        chapter_num,
                    )

                    if doc_pub_record and doc_pub_record.get("id"):
                        doc_pub_ids.append(doc_pub_record["id"])
                        print(f"        [OK] Chapter {chapter_num} publication submitted")

                # 4.5 Admin approves project publication
                print(f"    Admin reviewing project publication...")
                approved_record = approve_publication(
                    base_url,
                    admin_token,
                    pub_record_id,
                )

                if not approved_record or approved_record.get("status") != "approved":
                    print(f"    [WARN] Project review failed: {approved_record}")
                    continue

                book_id = approved_record.get("externalId") or approved_record.get("bookstoreId")
                print(f"    [OK] Project approved, book ID: {book_id}")

                # 4.6 Admin approves document publications
                for doc_pub_id in doc_pub_ids:
                    print(f"      Reviewing chapter publication {doc_pub_id}...")
                    doc_approved = approve_publication(
                        base_url,
                        admin_token,
                        doc_pub_id,
                    )
                    if doc_approved and doc_approved.get("status") == "approved":
                        print(f"        [OK] Chapter publication approved")

                # 4.7 Verify bookstore book exists
                if book_id:
                    book = get_bookstore_book(base_url, book_id)
                    if book:
                        print(f"    [OK] Bookstore book created: {book.get('title')}")
                        created_books.append(book_id)

                        # Verify chapters
                        chapters = get_bookstore_chapters(base_url, book_id)
                        print(f"    [OK] Bookstore chapters: {len(chapters)}")
                    else:
                        print(f"    [WARN] Bookstore book not found")
                else:
                    print(f"    [WARN] No book ID returned")

            except Exception as e:
                print(f"    [ERROR] Error processing book {book_title}: {e}")
                continue

    # 5. Generate test comments
    if len(reader_tokens) > 0 and len(created_books) > 0:
        print("\nStep 5: Generating test comments...")
        for book_id in created_books:
            try:
                # Generate 2-3 comments per book
                num_comments = 2 + (hash(book_id) % 2)
                comments = generate_test_comments(
                    base_url,
                    reader_tokens[0],
                    book_id,
                    num_comments,
                )
                print(f"  [OK] Generated {len(comments)} comments for book {book_id}")
            except Exception as e:
                print(f"  [WARN] Comment generation failed: {e}")

    print("\n" + "=" * 50)
    print("[SUCCESS] Publication flow seeder completed!")
    print("=" * 50)

    # Summary
    print("\nSummary:")
    print(f"  Test authors: {len(test_authors)}")
    print(f"  Books per author: {num_books_per_author}")
    print(f"  Chapters per book: {chapters_per_book}")
    print(f"  Created books: {len(created_books)}")
    print(f"  Total chapters: {len(test_authors) * num_books_per_author * chapters_per_book}")

    return {
        "authors": len(test_authors),
        "books_per_author": num_books_per_author,
        "chapters_per_book": chapters_per_book,
        "created_books": len(created_books),
    }


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Create test data through complete publication flow")
    parser.add_argument("--base-url", default=DEFAULT_BASE_URL, help="Backend base URL")
    parser.add_argument("--books-per-author", type=int, default=BOOKS_PER_AUTHOR, help="Books per author")
    parser.add_argument("--chapters-per-book", type=int, default=CHAPTERS_PER_BOOK, help="Chapters per book")
    args = parser.parse_args()

    try:
        seed_publication_flow(
            base_url=args.base_url,
            num_books_per_author=args.books_per_author,
            chapters_per_book=args.chapters_per_book,
        )
    except Exception as e:
        print(f"\n[ERROR] {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
