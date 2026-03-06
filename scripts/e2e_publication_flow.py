#!/usr/bin/env python3
import argparse
import json
import sys
import urllib.error
import urllib.parse
import urllib.request

DEFAULT_AUTHOR_USERNAME = "author_new"
DEFAULT_AUTHOR_PASSWORD = "Author@123456"
DEFAULT_ADMIN_USERNAME = "admin"
DEFAULT_ADMIN_PASSWORD = "Admin@123456"
DEFAULT_PROJECT_TITLE = "联调发布示例项目"
DEFAULT_DOCUMENT_TITLE = "第1章 风起青川"


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
        for key in ("items", "records", "list", "chapters", "documents", "projects", "results"):
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


def login(base_url: str, username: str, password: str) -> tuple[str, dict]:
    response = request_json(
        "POST",
        build_url(base_url, "/api/v1/shared/auth/login"),
        body={"username": username, "password": password},
    )
    data = api_data(response)
    if not isinstance(data, dict) or not data.get("token"):
        raise RuntimeError(f"login did not return a token for {username}")
    return data["token"], data.get("user", {})


def list_writer_projects(base_url: str, author_token: str) -> list[dict]:
    response = request_json(
        "GET",
        build_url(base_url, "/api/v1/writer/projects", {"page": 1, "pageSize": 50}),
        token=author_token,
    )
    return normalize_items(api_data(response))


def list_writer_documents(base_url: str, author_token: str, project_id: str) -> list[dict]:
    response = request_json(
        "GET",
        build_url(base_url, f"/api/v1/writer/project/{project_id}/documents", {"page": 1, "pageSize": 100}),
        token=author_token,
    )
    return normalize_items(api_data(response))


def find_by_title(items: list[dict], title: str) -> dict:
    for item in items:
        item_title = item.get("title")
        if item_title is None and isinstance(item.get("data"), dict):
            item_title = item["data"].get("title")
        if item_title == title:
            return item
    raise RuntimeError(f"could not find item with title: {title}")


def resolve_project(base_url: str, author_token: str, project_id: str | None, project_title: str) -> dict:
    if project_id:
        response = request_json(
            "GET",
            build_url(base_url, f"/api/v1/writer/projects/{project_id}"),
            token=author_token,
        )
        project = api_data(response)
        if not isinstance(project, dict) or not project.get("id"):
            raise RuntimeError("project lookup returned no data")
        return project

    projects = list_writer_projects(base_url, author_token)
    return find_by_title(projects, project_title)


def resolve_document(base_url: str, author_token: str, project_id: str, document_id: str | None, document_title: str) -> dict:
    if document_id:
        response = request_json(
            "GET",
            build_url(base_url, f"/api/v1/writer/documents/{document_id}"),
            token=author_token,
        )
        document = api_data(response)
        if not isinstance(document, dict) or not document.get("documentId"):
            raise RuntimeError("document lookup returned no data")
        return {
            "id": document["documentId"],
            "title": document.get("title"),
            "type": document.get("type"),
        }

    search_response = request_json(
        "GET",
        build_url(
            base_url,
            "/api/v1/writer/search/documents",
            {"q": document_title, "project_id": project_id, "page": 1, "page_size": 20},
        ),
        token=author_token,
    )
    search_items = normalize_items(api_data(search_response))
    if search_items:
        match = find_by_title(search_items, document_title)
        if isinstance(match.get("data"), dict):
            return {
                "id": match.get("id"),
                "title": match["data"].get("title"),
                "type": match["data"].get("type"),
            }

    documents = list_writer_documents(base_url, author_token, project_id)
    return find_by_title(documents, document_title)


def resolve_book_id(reviewed_record: dict) -> str:
    external_id = reviewed_record.get("externalId")
    if is_object_id(external_id):
        return external_id

    bookstore_id = reviewed_record.get("bookstoreId")
    if is_object_id(bookstore_id):
        return bookstore_id

    raise RuntimeError("review response did not contain a usable book id")


def resolve_chapter(chapters: list, document_id: str, chapter_number: int, external_id: str | None = None):
    if external_id:
        for chapter in chapters:
            if chapter.get("id") == external_id:
                return chapter
    for chapter in chapters:
        if chapter.get("projectChapterId") == document_id:
            return chapter
    for chapter in chapters:
        chapter_num = chapter.get("chapterNum")
        if chapter_num is None:
            chapter_num = chapter.get("chapter_num")
        if chapter_num == chapter_number:
            return chapter
    raise RuntimeError("could not resolve chapter id from chapter list")


def approve_publication(base_url: str, admin_token: str, record_id: str, note: str) -> dict:
    response = request_json(
        "POST",
        build_url(base_url, f"/api/v1/admin/publications/{record_id}/review"),
        token=admin_token,
        body={"action": "approve", "note": note},
    )
    data = api_data(response)
    if not isinstance(data, dict):
        raise RuntimeError(f"review response did not contain a record for {record_id}")
    return data


def get_project_publications(base_url: str, author_token: str, project_id: str) -> list[dict]:
    response = request_json(
        "GET",
        build_url(base_url, f"/api/v1/writer/projects/{project_id}/publications", {"page": 1, "pageSize": 20}),
        token=author_token,
    )
    return normalize_items(api_data(response))


def maybe_reset_publication(base_url: str, author_token: str, project_id: str, enabled: bool) -> bool:
    if not enabled:
        return False

    records = get_project_publications(base_url, author_token, project_id)
    has_published_project = any(
        record.get("type") == "project" and record.get("status") == "published" for record in records
    )
    if not has_published_project:
        return False

    request_json(
        "POST",
        build_url(base_url, f"/api/v1/writer/projects/{project_id}/unpublish"),
        token=author_token,
        body={},
    )
    return True


def main() -> int:
    parser = argparse.ArgumentParser(description="E2E validate writer -> publish -> review -> reader flow")
    parser.add_argument("--base-url", default="http://localhost:8080")
    parser.add_argument("--author-token")
    parser.add_argument("--admin-token")
    parser.add_argument("--author-username", default=DEFAULT_AUTHOR_USERNAME)
    parser.add_argument("--author-password", default=DEFAULT_AUTHOR_PASSWORD)
    parser.add_argument("--admin-username", default=DEFAULT_ADMIN_USERNAME)
    parser.add_argument("--admin-password", default=DEFAULT_ADMIN_PASSWORD)
    parser.add_argument("--project-id")
    parser.add_argument("--document-id")
    parser.add_argument("--project-title", default=DEFAULT_PROJECT_TITLE)
    parser.add_argument("--document-title", default=DEFAULT_DOCUMENT_TITLE)
    parser.add_argument("--chapter-number", type=int, default=1)
    parser.add_argument("--approve-document", action="store_true")
    parser.add_argument("--skip-reset", action="store_true")
    args = parser.parse_args()

    print("0. Resolve auth tokens")
    author_user = {}
    admin_user = {}
    author_token = args.author_token
    admin_token = args.admin_token
    if not author_token:
        author_token, author_user = login(args.base_url, args.author_username, args.author_password)
    if not admin_token:
        admin_token, admin_user = login(args.base_url, args.admin_username, args.admin_password)
    print(f"   author user: {author_user.get('username', args.author_username)}")
    print(f"   admin user: {admin_user.get('username', args.admin_username)}")

    print("0.1 Resolve seed project and document")
    project = resolve_project(args.base_url, author_token, args.project_id, args.project_title)
    project_id = project.get("id") or project.get("projectId")
    if not project_id:
        raise RuntimeError("could not resolve project id")
    document = resolve_document(args.base_url, author_token, project_id, args.document_id, args.document_title)
    document_id = document.get("id") or document.get("documentId")
    if not document_id:
        raise RuntimeError("could not resolve document id")
    print(f"   project: {project.get('title')} ({project_id})")
    print(f"   document: {document.get('title')} ({document_id})")

    print("0.2 Reset published state when needed")
    reset_done = maybe_reset_publication(args.base_url, author_token, project_id, not args.skip_reset)
    print(f"   reset performed: {reset_done}")

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
        build_url(args.base_url, f"/api/v1/writer/projects/{project_id}/publish"),
        token=author_token,
        body=project_publish_body,
    )
    project_record = api_data(project_publish_resp)
    if not isinstance(project_record, dict) or not project_record.get("id"):
        raise RuntimeError("project publication record was not returned")
    print(f"   project record id: {project_record['id']}, status: {project_record.get('status')}")

    print("2. Author submits single document publication for review")
    document_publish_body = {
        "chapterTitle": document.get("title") or f"MVP Chapter {args.chapter_number}",
        "chapterNumber": args.chapter_number,
        "isFree": True,
        "authorNote": "Document publish from e2e python script",
    }
    document_publish_resp = request_json(
        "POST",
        build_url(
            args.base_url,
            f"/api/v1/writer/documents/{document_id}/publish",
            {"projectId": project_id},
        ),
        token=author_token,
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
        token=admin_token,
    )
    pending_items = normalize_items(api_data(pending_resp))
    print(f"   pending count (response page): {len(pending_items)}")

    print("4. Admin approves the project publication")
    reviewed_record = approve_publication(
        args.base_url,
        admin_token,
        project_record["id"],
        "Approved by e2e publication flow script",
    )
    print(f"   project review status: {reviewed_record.get('status')}")

    book_id = resolve_book_id(reviewed_record)

    reviewed_document_record = None
    if args.approve_document:
        print("4.1 Admin approves the document publication")
        reviewed_document_record = approve_publication(
            args.base_url,
            admin_token,
            document_record["id"],
            "Approved document by e2e publication flow script",
        )
        print(f"   document review status: {reviewed_document_record.get('status')}")

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
    chapter_external_id = None
    if isinstance(reviewed_document_record, dict):
        chapter_external_id = reviewed_document_record.get("externalId")
    chapter = resolve_chapter(chapter_items, document_id, args.chapter_number, chapter_external_id)

    print("6. Reader-side chapter access verification")
    reader_chapter_resp = request_json(
        "GET",
        build_url(args.base_url, f"/api/v1/reader/books/{bookstore_book['id']}/chapters/{chapter['id']}"),
        token=author_token,
    )
    reader_chapter = api_data(reader_chapter_resp)

    summary = {
        "author": author_user.get("username", args.author_username),
        "admin": admin_user.get("username", args.admin_username),
        "projectId": project_id,
        "documentId": document_id,
        "projectPublicationRecordId": project_record["id"],
        "documentPublicationRecordId": document_record["id"],
        "documentApproved": args.approve_document,
        "documentReviewRecord": reviewed_document_record,
        "bookId": bookstore_book["id"],
        "chapterId": chapter["id"],
        "resetPerformed": reset_done,
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
