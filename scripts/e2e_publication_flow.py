#!/usr/bin/env python3
import argparse
import json
import sys
import urllib.error
import urllib.parse
import urllib.request


def build_url(base_url: str, path: str, query: dict | None = None) -> str:
    base = base_url.rstrip("/")
    url = f"{base}{path}"
    if query:
        url = f"{url}?{urllib.parse.urlencode(query)}"
    return url


def request_json(method: str, url: str, token: str | None = None, body: dict | None = None):
    headers = {"Accept": "application/json"}
    data = None
    if token:
        headers["Authorization"] = f"Bearer {token}"
    if body is not None:
        data = json.dumps(body).encode("utf-8")
        headers["Content-Type"] = "application/json"

    req = urllib.request.Request(url=url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req) as resp:
            payload = resp.read().decode("utf-8")
    except urllib.error.HTTPError as exc:
        payload = exc.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"{method} {url} failed: {exc.code} {payload}") from exc
    except urllib.error.URLError as exc:
        raise RuntimeError(f"{method} {url} failed: {exc}") from exc

    if not payload:
        return None
    return json.loads(payload)


def api_data(response):
    if response is None:
        return None
    if isinstance(response, dict) and "data" in response:
        return response["data"]
    return response


def normalize_items(value):
    if value is None:
        return []
    if isinstance(value, list):
        return value
    if isinstance(value, dict):
        for key in ("items", "records", "list", "chapters"):
            nested = value.get(key)
            if isinstance(nested, list):
                return nested
        return [value]
    return []


def is_object_id(value):
    if not isinstance(value, str) or len(value) != 24:
        return False
    try:
        int(value, 16)
        return True
    except ValueError:
        return False


def resolve_book_id(reviewed_record: dict) -> str:
    external_id = reviewed_record.get("externalId")
    if is_object_id(external_id):
        return external_id

    bookstore_id = reviewed_record.get("bookstoreId")
    if is_object_id(bookstore_id):
        return bookstore_id

    raise RuntimeError("review response did not contain a usable book id")


def resolve_chapter(chapters: list, document_id: str, chapter_number: int):
    for chapter in chapters:
        if chapter.get("projectChapterId") == document_id:
            return chapter
    for chapter in chapters:
        if chapter.get("chapterNum") == chapter_number:
            return chapter
    raise RuntimeError("could not resolve chapter id from chapter list")


def main() -> int:
    parser = argparse.ArgumentParser(description="E2E validate writer -> publish -> review -> reader flow")
    parser.add_argument("--base-url", default="http://localhost:8080")
    parser.add_argument("--author-token", required=True)
    parser.add_argument("--admin-token", required=True)
    parser.add_argument("--project-id", required=True)
    parser.add_argument("--document-id", required=True)
    parser.add_argument("--chapter-number", type=int, default=1)
    args = parser.parse_args()

    project_publish_body = {
        "bookstoreId": "local",
        "categoryId": "mvp-category",
        "tags": ["mvp", "publication"],
        "description": "Publication MVP end-to-end flow",
        "coverImage": "",
        "publishType": "serial",
        "freeChapters": 1,
        "authorNote": "Submitted by e2e python script",
        "enableComment": True,
        "enableShare": True,
    }

    print("1. Author submits project publication for review")
    project_publish_resp = request_json(
        "POST",
        build_url(args.base_url, f"/api/v1/writer/projects/{args.project_id}/publish"),
        token=args.author_token,
        body=project_publish_body,
    )
    project_record = api_data(project_publish_resp)
    if not isinstance(project_record, dict) or not project_record.get("id"):
        raise RuntimeError("project publication record was not returned")
    print(f"   project record id: {project_record['id']}, status: {project_record.get('status')}")

    print("2. Author submits single document publication for review")
    document_publish_body = {
        "chapterTitle": f"MVP Chapter {args.chapter_number}",
        "chapterNumber": args.chapter_number,
        "isFree": True,
        "authorNote": "Document publish from e2e python script",
    }
    document_publish_resp = request_json(
        "POST",
        build_url(
            args.base_url,
            f"/api/v1/writer/documents/{args.document_id}/publish",
            {"projectId": args.project_id},
        ),
        token=args.author_token,
        body=document_publish_body,
    )
    document_record = api_data(document_publish_resp)
    if not isinstance(document_record, dict) or not document_record.get("id"):
        raise RuntimeError("document publication record was not returned")
    print(f"   document record id: {document_record['id']}, status: {document_record.get('status')}")

    print("3. Admin fetches pending publication queue")
    pending_resp = request_json(
        "GET",
        build_url(args.base_url, "/api/v1/admin/publications/pending", {"page": 1, "pageSize": 20}),
        token=args.admin_token,
    )
    pending_items = normalize_items(api_data(pending_resp))
    print(f"   pending count (response page): {len(pending_items)}")

    print("4. Admin approves the project publication")
    review_resp = request_json(
        "POST",
        build_url(args.base_url, f"/api/v1/admin/publications/{project_record['id']}/review"),
        token=args.admin_token,
        body={"action": "approve", "note": "Approved by e2e publication flow script"},
    )
    reviewed_record = api_data(review_resp)
    if not isinstance(reviewed_record, dict):
        raise RuntimeError("project review response did not contain a record")
    print(f"   project review status: {reviewed_record.get('status')}")

    book_id = resolve_book_id(reviewed_record)

    print("5. Read-side verification through bookstore routes")
    bookstore_book_resp = request_json(
        "GET",
        build_url(args.base_url, f"/api/v1/bookstore/books/{book_id}"),
    )
    bookstore_book = api_data(bookstore_book_resp)
    if not isinstance(bookstore_book, dict) or not bookstore_book.get("id"):
        raise RuntimeError("bookstore book lookup returned no data")
    print(f"   bookstore book id: {bookstore_book['id']}, title: {bookstore_book.get('title')}")

    chapter_list_resp = request_json(
        "GET",
        build_url(args.base_url, f"/api/v1/bookstore/books/{bookstore_book['id']}/chapters"),
    )
    chapter_items = normalize_items(api_data(chapter_list_resp))
    chapter = resolve_chapter(chapter_items, args.document_id, args.chapter_number)

    print("6. Reader-side chapter access verification")
    reader_chapter_resp = request_json(
        "GET",
        build_url(args.base_url, f"/api/v1/reader/books/{bookstore_book['id']}/chapters/{chapter['id']}"),
        token=args.author_token,
    )
    reader_chapter = api_data(reader_chapter_resp)

    summary = {
        "projectPublicationRecordId": project_record["id"],
        "documentPublicationRecordId": document_record["id"],
        "bookId": bookstore_book["id"],
        "chapterId": chapter["id"],
        "chapterList": chapter_items,
        "readerChapter": reader_chapter,
    }
    print(json.dumps(summary, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        print(f"ERROR: {exc}", file=sys.stderr)
        raise SystemExit(1)
