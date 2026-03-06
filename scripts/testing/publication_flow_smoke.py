#!/usr/bin/env python3
import argparse
import os
import subprocess
import sys
import time
import urllib.error
import urllib.request
from pathlib import Path


def wait_for_http(url: str, timeout_seconds: int) -> None:
    deadline = time.time() + timeout_seconds
    last_error = None
    while time.time() < deadline:
        try:
            with urllib.request.urlopen(url) as resp:
                if 200 <= resp.status < 500:
                    return
        except Exception as exc:  # noqa: BLE001
            last_error = exc
        time.sleep(2)
    raise RuntimeError(f"timed out waiting for {url}: {last_error}")


def run_checked(command: list[str], cwd: Path, env: dict[str, str]) -> None:
    completed = subprocess.run(command, cwd=cwd, env=env, check=False)
    if completed.returncode != 0:
        raise RuntimeError(f"command failed ({completed.returncode}): {' '.join(command)}")


def terminate_process(process: subprocess.Popen) -> None:
    if process.poll() is not None:
        return
    process.terminate()
    try:
        process.wait(timeout=10)
    except subprocess.TimeoutExpired:
        process.kill()
        process.wait(timeout=10)


def main() -> int:
    parser = argparse.ArgumentParser(description="Run publication flow smoke test for CI or local automation")
    parser.add_argument("--repo-root", default=str(Path(__file__).resolve().parents[2]))
    parser.add_argument("--base-url", default=os.environ.get("PUBLICATION_SMOKE_BASE_URL", "http://localhost:9090"))
    parser.add_argument("--server-port", default=os.environ.get("PUBLICATION_SMOKE_PORT", "9090"))
    parser.add_argument("--server-health-path", default="/health")
    parser.add_argument("--health-timeout", type=int, default=90)
    parser.add_argument("--server-log", default=None)
    parser.add_argument("--approve-document", action="store_true")
    parser.add_argument("--reject-project", action="store_true")
    parser.add_argument("--reject-document", action="store_true")
    parser.add_argument("--retry-after-reject", action="store_true")
    parser.add_argument("--skip-seed", action="store_true")
    parser.add_argument("--skip-build", action="store_true")
    parser.add_argument("--keep-server", action="store_true")
    args = parser.parse_args()

    if args.reject_project and args.reject_document:
        raise RuntimeError("--reject-project and --reject-document cannot be used together")
    if args.approve_document and args.reject_project:
        raise RuntimeError("--approve-document cannot be combined with --reject-project")
    if args.retry_after_reject and (args.reject_project or args.reject_document):
        raise RuntimeError("--retry-after-reject cannot be combined with reject modes")

    repo_root = Path(args.repo_root).resolve()
    script_path = repo_root / "scripts" / "e2e_publication_flow.py"
    if not script_path.exists():
        raise RuntimeError(f"missing e2e script: {script_path}")

    env = os.environ.copy()
    if env.get("MONGODB_URI") and not env.get("QINGYU_DATABASE_PRIMARY_MONGODB_URI"):
        env["QINGYU_DATABASE_PRIMARY_MONGODB_URI"] = env["MONGODB_URI"]
    if env.get("MONGODB_DATABASE") and not env.get("QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE"):
        env["QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE"] = env["MONGODB_DATABASE"]
    env.setdefault("REDIS_ADDR", env.get("REDIS_ADDR", "localhost:6379"))
    env.setdefault("QINGYU_SERVER_PORT", args.server_port)

    if not args.skip_seed:
        print("[1/4] Seed test data")
        run_checked(["go", "run", "-tags", "auto", "./cmd/seed_data"], repo_root, env)

    binary_name = "publication_smoke_server.exe" if os.name == "nt" else "publication_smoke_server"
    binary_path = repo_root / binary_name
    if not args.skip_build:
        print("[2/4] Build server binary")
        run_checked(["go", "build", "-o", str(binary_path), "./cmd/server/main.go"], repo_root, env)

    log_path = Path(args.server_log) if args.server_log else repo_root / "tmp_publication_smoke_server.log"
    print(f"[3/4] Start server: {binary_path}")
    with log_path.open("w", encoding="utf-8") as log_file:
        process = subprocess.Popen(
            [str(binary_path)],
            cwd=repo_root,
            env=env,
            stdout=log_file,
            stderr=subprocess.STDOUT,
        )
    try:
        wait_for_http(args.base_url.rstrip("/") + args.server_health_path, args.health_timeout)
        print("[4/4] Run publication flow e2e")
        command = [sys.executable, str(script_path), "--base-url", args.base_url]
        if args.approve_document:
            command.append("--approve-document")
        if args.reject_project:
            command.append("--reject-project")
        if args.reject_document:
            command.append("--reject-document")
        if args.retry_after_reject:
            command.append("--retry-after-reject")
        run_checked(command, repo_root, env)
        print("publication flow smoke passed")
    except Exception:
        if log_path.exists():
            print("=== publication smoke server log ===", file=sys.stderr)
            print(log_path.read_text(encoding="utf-8", errors="replace"), file=sys.stderr)
        raise
    finally:
        if not args.keep_server:
            terminate_process(process)
        if binary_path.exists() and not args.keep_server:
            try:
                binary_path.unlink()
            except OSError:
                pass

    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:  # noqa: BLE001
        print(f"ERROR: {exc}", file=sys.stderr)
        raise SystemExit(1)
