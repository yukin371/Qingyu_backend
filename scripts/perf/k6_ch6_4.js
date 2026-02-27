import http from "k6/http";
import { check, group, sleep } from "k6";

const BASE_URL = __ENV.BASE_URL || "http://127.0.0.1:9090";
const LOGIN_USERNAME = __ENV.LOGIN_USERNAME || "";
const LOGIN_PASSWORD = __ENV.LOGIN_PASSWORD || "";

const STAGE = __ENV.TEST_STAGE || "normal";

const STAGE_MAP = {
  baseline: { vus: 10, duration: "1m" },
  normal: { vus: 50, duration: "3m" },
  stress: { vus: 100, duration: "3m" },
};

const selected = STAGE_MAP[STAGE] || STAGE_MAP.normal;

export const options = {
  vus: selected.vus,
  duration: selected.duration,
  thresholds: {
    http_req_failed: ["rate<0.05"],
    http_req_duration: ["p(95)<1200"],
  },
  summaryTrendStats: ["avg", "min", "med", "p(90)", "p(95)", "max"],
};

function jsonHeaders(token = "") {
  const headers = {
    "Content-Type": "application/json",
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  return headers;
}

function loginIfConfigured() {
  if (!LOGIN_USERNAME || !LOGIN_PASSWORD) {
    return "";
  }

  const payload = JSON.stringify({
    username: LOGIN_USERNAME,
    password: LOGIN_PASSWORD,
  });

  const res = http.post(`${BASE_URL}/api/v1/user/auth/login`, payload, {
    headers: jsonHeaders(),
    tags: { api: "user_login" },
  });

  const ok = check(res, {
    "login status is 200": (r) => r.status === 200,
    "login response has token": (r) => {
      try {
        const body = r.json();
        return !!(body && body.data && body.data.token);
      } catch (e) {
        return false;
      }
    },
  });

  if (!ok) {
    return "";
  }

  try {
    const body = res.json();
    return body.data.token || "";
  } catch (e) {
    return "";
  }
}

export default function () {
  group("public_endpoints", () => {
    const health = http.get(`${BASE_URL}/health`, {
      tags: { api: "health" },
    });
    check(health, {
      "GET /health status is 200": (r) => r.status === 200,
    });

    const homepage = http.get(`${BASE_URL}/api/v1/bookstore/homepage`, {
      tags: { api: "bookstore_homepage" },
    });
    check(homepage, {
      "GET /bookstore/homepage status is 200": (r) => r.status === 200,
    });

    const books = http.get(`${BASE_URL}/api/v1/bookstore/books?page=1&pageSize=10`, {
      tags: { api: "bookstore_books" },
    });
    check(books, {
      "GET /bookstore/books status is 200": (r) => r.status === 200,
    });
  });

  group("auth_optional_endpoints", () => {
    const token = loginIfConfigured();
    if (!token) {
      return;
    }

    const profile = http.get(`${BASE_URL}/api/v1/user/profile`, {
      headers: jsonHeaders(token),
      tags: { api: "user_profile" },
    });
    check(profile, {
      "GET /user/profile status is 200": (r) => r.status === 200,
    });
  });

  sleep(1);
}

