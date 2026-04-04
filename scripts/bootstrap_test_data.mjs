#!/usr/bin/env node

import { execFileSync } from 'node:child_process';
import { parseArgs } from 'node:util';

const DEFAULT_BASE_URL = 'http://localhost:9090';
const DEFAULT_AUTHOR_USERNAME = 'testauthor001';
const DEFAULT_AUTHOR_PASSWORD = 'password';
const DEFAULT_ADMIN_USERNAME = 'testadmin001';
const DEFAULT_ADMIN_PASSWORD = 'password';
const DEFAULT_READER_USERS = ['testuser001', 'testuser002', 'testuser003', 'testuser004'];

const { values } = parseArgs({
  options: {
    mode: { type: 'string', default: 'bootstrap' },
    scale: { type: 'string', default: 'small' },
    baseUrl: { type: 'string', default: DEFAULT_BASE_URL },
    authorUsername: { type: 'string', default: DEFAULT_AUTHOR_USERNAME },
    authorPassword: { type: 'string', default: DEFAULT_AUTHOR_PASSWORD },
    adminUsername: { type: 'string', default: DEFAULT_ADMIN_USERNAME },
    adminPassword: { type: 'string', default: DEFAULT_ADMIN_PASSWORD },
    projectId: { type: 'string' },
    projectTitle: { type: 'string' },
    bookId: { type: 'string' },
    bookTitle: { type: 'string' },
    limit: { type: 'string', default: '10' },
    readers: { type: 'string', default: DEFAULT_READER_USERS.join(',') },
    clean: { type: 'boolean', default: false },
  },
});

const rootDir = new URL('..', import.meta.url);
const backendDir = filePath(rootDir);

async function main() {
  const mode = values.mode;
  if (mode === 'baseline') {
    runSeeder('baseline', values.clean);
    return;
  }
  if (mode === 'full') {
    runSeeder('full', values.clean);
    return;
  }
  if (mode === 'bootstrap') {
    runSeeder('baseline', values.clean);
  }

  if (mode === 'bootstrap' || mode === 'author-stats') {
    const readers = values.readers
      .split(',')
      .map((item) => item.trim())
      .filter(Boolean);
    await enrichAuthorStats({
      baseUrl: values.baseUrl,
      authorUsername: values.authorUsername,
      authorPassword: values.authorPassword,
      adminUsername: values.adminUsername,
      adminPassword: values.adminPassword,
      projectId: values.projectId,
      projectTitle: values.projectTitle,
      readers,
    });
    runSeeder('stats', false);
    return;
  }

  if (mode === 'author-stats-all') {
    const readers = values.readers
      .split(',')
      .map((item) => item.trim())
      .filter(Boolean);
    await enrichAllAuthorProjects({
      baseUrl: values.baseUrl,
      authorUsername: values.authorUsername,
      authorPassword: values.authorPassword,
      adminUsername: values.adminUsername,
      adminPassword: values.adminPassword,
      readers,
      limit: Number.parseInt(values.limit, 10) || 10,
    });
    runSeeder('stats', false);
    return;
  }

  if (mode === 'showcase-book') {
    const readers = values.readers
      .split(',')
      .map((item) => item.trim())
      .filter(Boolean);
    await enrichShowcaseBook({
      baseUrl: values.baseUrl,
      bookId: values.bookId || values.projectId,
      bookTitle: values.bookTitle || values.projectTitle,
      readers,
    });
    runSeeder('stats', false);
    return;
  }

  throw new Error(`unsupported mode: ${mode}`);
}

function runSeeder(command, clean) {
  const args = ['run', './cmd/seeder', command, '--scale', values.scale];
  if (clean) {
    args.push('--clean');
  }
  console.log(`> go ${args.join(' ')}`);
  execFileSync('go', args, { cwd: backendDir, stdio: 'inherit' });
}

async function enrichAuthorStats({
  baseUrl,
  authorUsername,
  authorPassword,
  adminUsername,
  adminPassword,
  projectId,
  projectTitle,
  readers,
}) {
  const authorToken = await login(baseUrl, authorUsername, authorPassword);
  const adminToken = await login(baseUrl, adminUsername, adminPassword);
  const project = await resolveProject(baseUrl, authorToken, projectId, projectTitle);
  await enrichProjectStats(baseUrl, authorToken, adminToken, project, readers);
}

async function enrichAllAuthorProjects({
  baseUrl,
  authorUsername,
  authorPassword,
  adminUsername,
  adminPassword,
  readers,
  limit,
}) {
  const authorToken = await login(baseUrl, authorUsername, authorPassword);
  const adminToken = await login(baseUrl, adminUsername, adminPassword);
  const response = await requestJson('GET', `${baseUrl}/api/v1/writer/projects?page=1&pageSize=${Math.max(limit, 1)}`, {
    token: authorToken,
  });
  const projects = normalizeItems(response.data).slice(0, Math.max(limit, 1));
  if (!projects.length) {
    throw new Error(`no writer projects found for ${authorUsername}`);
  }

  for (const project of projects) {
    try {
      await enrichProjectStats(baseUrl, authorToken, adminToken, project, readers);
    } catch (error) {
      console.warn(`skip project ${project?.id ?? '<unknown>'}: ${error.message}`);
    }
  }
}

async function enrichShowcaseBook({
  baseUrl,
  bookId,
  bookTitle,
  readers,
}) {
  await ensureReaderAccounts(baseUrl, readers);
  const targetBook = await resolveShowcaseBook(baseUrl, bookId, bookTitle);
  const chapters = await requestJson(
    'GET',
    `${baseUrl}/api/v1/bookstore/books/${targetBook.id}/chapters?page=1&size=20`,
  );
  const chapterItems = normalizeItems(chapters.data);
  const firstChapter = chapterItems[0];
  if (!firstChapter?.id) {
    throw new Error(`no chapter found for showcase book ${targetBook.id}`);
  }

  console.log(`showcase book: ${targetBook.title} (${targetBook.id})`);

  for (const [index, readerUsername] of readers.entries()) {
    const readerToken = await login(baseUrl, readerUsername, 'password');
    await createOrUpdateRating(
      baseUrl,
      readerToken,
      targetBook.id,
      [5, 5, 4, 4][index % 4],
      `${readerUsername}：剑墟线索铺陈扎实，适合追更。`,
    );
    await ensureBookComment(
      baseUrl,
      readerToken,
      readerUsername,
      targetBook.id,
      firstChapter.id,
      `${readerUsername}：顾照临这一卷开篇很稳，海港和旧山门的气质立住了，愿意继续追更。`,
      [5, 5, 4, 4][index % 4],
    );
    await ensureBookCollection(baseUrl, readerToken, targetBook.id);
    await createBookmark(baseUrl, readerToken, targetBook.id, firstChapter.id, targetBook.title, index);
    await recordReadingHistory(baseUrl, readerToken, targetBook.id, firstChapter.id, index);
    await recordReaderBehaviors(baseUrl, readerToken, targetBook.id, firstChapter.id, index);
  }

  console.log(`showcase book enrichment completed for ${targetBook.title}`);
}

async function enrichProjectStats(baseUrl, authorToken, adminToken, project, readers) {
  const publication = await ensurePublishedBook(baseUrl, authorToken, adminToken, project);
  const bookId = publication.bookId;
  const chapters = await requestJson('GET', `${baseUrl}/api/v1/bookstore/books/${bookId}/chapters`);
  const chapterItems = normalizeItems(chapters.data);
  const firstChapter = chapterItems[0];
  if (!firstChapter?.id) {
    throw new Error(`no published chapter found for book ${bookId}`);
  }

  console.log(`author project: ${project.title} (${project.id})`);
  console.log(`published book: ${bookId}`);

  for (const [index, readerUsername] of readers.entries()) {
    const readerToken = await login(baseUrl, readerUsername, 'password');
    await createOrUpdateRating(baseUrl, readerToken, bookId, 4 + (index % 2), `${readerUsername} rating for ${project.title}`);
    await createBookComment(
      baseUrl,
      readerToken,
      bookId,
      firstChapter.id,
      `${readerUsername} automated feedback for ${project.title}. 这是一条用于联调统计链路的测试评论。`,
      4 + (index % 2),
    );
    await createBookmark(baseUrl, readerToken, bookId, firstChapter.id, project.title, index);
    await recordReadingHistory(baseUrl, readerToken, bookId, firstChapter.id, index);
    await recordReaderBehaviors(baseUrl, readerToken, bookId, firstChapter.id, index);
  }

  console.log(`author stats enrichment completed for ${project.id}`);
}

async function login(baseUrl, username, password) {
  const response = await requestJson('POST', `${baseUrl}/api/v1/user/auth/login`, {
    body: { username, password },
  });
  const token = response?.data?.token;
  if (!token) {
    throw new Error(`login failed for ${username}`);
  }
  return token;
}

async function ensureReaderAccounts(baseUrl, readers) {
  for (const username of readers) {
    try {
      await login(baseUrl, username, 'password');
      continue;
    } catch (_error) {
      await registerReader(baseUrl, username);
      await login(baseUrl, username, 'password');
    }
  }
}

async function registerReader(baseUrl, username) {
  try {
    await requestJson('POST', `${baseUrl}/api/v1/user/auth/register`, {
      body: {
        username,
        email: `${username}@qingyu.test`,
        password: 'password',
      },
    });
  } catch (error) {
    const text = String(error.message || '');
    if (text.includes('409') || text.includes('用户名已被注册') || text.includes('邮箱已被注册')) {
      return;
    }
    throw error;
  }
}

async function resolveProject(baseUrl, token, projectId, projectTitle) {
  if (projectId) {
    const response = await requestJson('GET', `${baseUrl}/api/v1/writer/projects/${projectId}`, { token });
    if (response?.data?.id) {
      return response.data;
    }
    throw new Error(`project not found: ${projectId}`);
  }

  const response = await requestJson('GET', `${baseUrl}/api/v1/writer/projects?page=1&pageSize=50`, { token });
  const items = normalizeItems(response.data);
  if (!items.length) {
    throw new Error('writer has no projects');
  }

  if (!projectTitle) {
    return items[0];
  }

  const project = items.find((item) => item.title === projectTitle);
  if (!project) {
    throw new Error(`project not found by title: ${projectTitle}`);
  }
  return project;
}

async function resolvePublishedBook(baseUrl, token, projectId) {
  const response = await requestJson(
    'GET',
    `${baseUrl}/api/v1/writer/projects/${projectId}/publications?page=1&pageSize=50`,
    { token },
  );
  const items = normalizeItems(response.data);
  const projectRecord = items.find((item) => item.type === 'project' && item.status === 'published');
  if (!projectRecord) {
    throw new Error(`published project record not found for ${projectId}`);
  }

  const bookId = projectRecord.externalId || projectRecord.bookstoreId;
  if (!bookId) {
    throw new Error(`published book id not found for ${projectId}`);
  }
  return { record: projectRecord, bookId };
}

async function ensurePublishedBook(baseUrl, authorToken, adminToken, project) {
  try {
    return await resolvePublishedBook(baseUrl, authorToken, project.id);
  } catch (_error) {
    const documents = await requestJson(
      'GET',
      `${baseUrl}/api/v1/writer/project/${project.id}/documents?page=1&pageSize=100`,
      { token: authorToken },
    );
    const items = normalizeItems(documents.data);
    const document = items[0];
    if (!document?.id) {
      throw new Error(`project ${project.id} has no documents to publish`);
    }

    const projectRecord = await requestJson('POST', `${baseUrl}/api/v1/writer/projects/${project.id}/publish`, {
      token: authorToken,
      body: {
        bookstoreId: 'local',
        categoryId: project.categoryId ?? 'mvp-category',
        tags: ['seeded', 'cross-role'],
        description: `${project.title ?? 'project'} automated publication flow`,
        coverImage: '',
        publishType: 'serial',
        freeChapters: 1,
        authorNote: 'Automated cross-role bootstrap publish',
        enableComment: true,
        enableShare: true,
      },
    });
    const projectRecordId = projectRecord?.data?.id ?? projectRecord?.id;
    if (!projectRecordId) {
      throw new Error(`project publication record missing for ${project.id}`);
    }

    const chapterNumber = document.chapterNum ?? document.chapter_num ?? 1;
    const documentTitle = document.title ?? `Chapter ${chapterNumber}`;
    const documentRecord = await requestJson(
      'POST',
      `${baseUrl}/api/v1/writer/documents/${document.id}/publish?projectId=${project.id}`,
      {
        token: authorToken,
        body: {
          chapterTitle: documentTitle,
          chapterNumber,
          isFree: true,
          authorNote: 'Automated cross-role bootstrap chapter publish',
        },
      },
    );
    const documentRecordId = documentRecord?.data?.id ?? documentRecord?.id;
    if (!documentRecordId) {
      throw new Error(`document publication record missing for ${document.id}`);
    }

    const reviewedProject = await reviewPublication(baseUrl, adminToken, projectRecordId, 'approve');
    await reviewPublication(baseUrl, adminToken, documentRecordId, 'approve');

    const bookId = reviewedProject.externalId || reviewedProject.bookstoreId;
    if (!bookId) {
      throw new Error(`reviewed project record missing bookstore id for ${project.id}`);
    }
    return { record: reviewedProject, bookId };
  }
}

async function reviewPublication(baseUrl, adminToken, recordId, action) {
  const response = await requestJson('POST', `${baseUrl}/api/v1/admin/publications/${recordId}/review`, {
    token: adminToken,
    body: {
      action,
      note: `Automated ${action} by bootstrap_test_data.mjs`,
    },
  });
  return response?.data ?? response;
}

async function createOrUpdateRating(baseUrl, token, bookId, score, review) {
  await requestJson('POST', `${baseUrl}/api/v1/bookstore/ratings`, {
    token,
    body: { bookId, score, review, tags: ['seeded', 'author-stats'] },
  });
}

async function createBookComment(baseUrl, token, bookId, chapterId, content, rating) {
  await requestJson('POST', `${baseUrl}/api/v1/reader/comments`, {
    token,
    body: { book_id: bookId, chapter_id: chapterId, content, rating },
  });
}

async function ensureBookComment(baseUrl, token, readerUsername, bookId, chapterId, content, rating) {
  const response = await requestJson(
    'GET',
    `${baseUrl}/api/v1/reader/comments?book_id=${bookId}&page=1&size=100&sortBy=latest`,
    { token },
  );
  const items = normalizeItems(response.data);
  const existing = items.find((item) => String(item?.content || '').includes(`${readerUsername}：`));
  if (existing) {
    return;
  }
  await createBookComment(baseUrl, token, bookId, chapterId, content, rating);
}

async function ensureBookCollection(baseUrl, token, bookId) {
  const collected = await requestJson(
    'GET',
    `${baseUrl}/api/v1/social/collections/check?book_id=${bookId}`,
    { token },
  );
  if (collected?.data?.is_collected || collected?.is_collected) {
    return;
  }
  await requestJson('POST', `${baseUrl}/api/v1/social/collections`, {
    token,
    body: { book_id: bookId },
  });
}

async function createBookmark(baseUrl, token, bookId, chapterId, projectTitle, offset) {
  const position = 120 + offset * 40;
  try {
    await requestJson('POST', `${baseUrl}/api/v1/reader/books/${bookId}/bookmarks`, {
      token,
      body: {
        bookId,
        chapterId,
        position,
        note: `${projectTitle} 自动书签 ${offset + 1}`,
        color: offset % 2 === 0 ? 'amber' : 'blue',
        quote: `${projectTitle} bootstrap bookmark`,
        isPublic: offset % 2 === 0,
        tags: ['seeded', 'reader-engagement'],
      },
    });
  } catch (error) {
    if (String(error.message).includes('409')) {
      return;
    }
    throw error;
  }
}

async function recordReadingHistory(baseUrl, token, bookId, chapterId, offset) {
  const now = new Date();
  const startTime = new Date(now.getTime() - (offset + 1) * 15 * 60 * 1000);
  const endTime = new Date(startTime.getTime() + (8 + offset) * 60 * 1000);
  await requestJson('POST', `${baseUrl}/api/v1/reader/reading-history`, {
    token,
    body: {
      book_id: bookId,
      chapter_id: chapterId,
      start_time: startTime.toISOString(),
      end_time: endTime.toISOString(),
      progress: 68 + offset * 6,
      device_type: offset % 2 === 0 ? 'mobile' : 'desktop',
      device_id: `bootstrap-reader-${offset + 1}`,
    },
  });
}

async function recordReaderBehaviors(baseUrl, token, bookId, chapterId, offset) {
  const now = new Date();
  const baseReadAt = new Date(now.getTime() - offset * 60 * 60 * 1000);
  const viewPayload = {
    book_id: bookId,
    chapter_id: chapterId,
    behavior_type: 'view',
    start_position: 0,
    end_position: 600 + offset * 30,
    progress: 0.72,
    read_duration: 480 + offset * 25,
    read_at: baseReadAt.toISOString(),
    device_type: 'mobile',
    client_ip: `10.0.0.${20 + offset}`,
    source: 'bookshelf',
    referrer: '/bookshelf',
  };
  const completePayload = {
    ...viewPayload,
    behavior_type: 'complete',
    end_position: 1200 + offset * 50,
    progress: 1,
    read_duration: 820 + offset * 30,
    read_at: new Date(baseReadAt.getTime() + 120000).toISOString(),
  };

  await requestJson('POST', `${baseUrl}/api/v1/writer/reader/behavior`, { token, body: viewPayload });
  await requestJson('POST', `${baseUrl}/api/v1/writer/reader/behavior`, { token, body: completePayload });
}

async function requestJson(method, url, { token, body } = {}) {
  const headers = { Accept: 'application/json' };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  if (body !== undefined) {
    headers['Content-Type'] = 'application/json';
  }

  const response = await fetch(url, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  const text = await response.text();
  const payload = text ? JSON.parse(text) : null;
  if (!response.ok) {
    throw new Error(`${method} ${url} failed: ${response.status} ${text}`);
  }
  return payload;
}

async function resolveShowcaseBook(baseUrl, bookId, bookTitle) {
  if (bookId) {
    const response = await requestJson('GET', `${baseUrl}/api/v1/bookstore/books/${bookId}/detail`);
    const detail = response?.data ?? response;
    if (detail?.id) {
      return detail;
    }
    throw new Error(`showcase book not found: ${bookId}`);
  }

  if (!bookTitle) {
    throw new Error('showcase-book mode requires --projectId <bookId> or --projectTitle <title>');
  }

  const response = await requestJson(
    'GET',
    `${baseUrl}/api/v1/bookstore/books/search/title?title=${encodeURIComponent(bookTitle)}`,
  );
  const items = normalizeItems(response.data);
  const target = items.find((item) => item.title === bookTitle);
  if (!target?.id) {
    throw new Error(`showcase book not found by title: ${bookTitle}`);
  }
  return target;
}

function normalizeItems(value) {
  if (Array.isArray(value)) {
    return value;
  }
  if (!value || typeof value !== 'object') {
    return [];
  }
  for (const key of ['items', 'records', 'list', 'chapters', 'projects', 'results']) {
    if (Array.isArray(value[key])) {
      return value[key];
    }
  }
  return [];
}

function filePath(url) {
  return decodeURIComponent(url.pathname.replace(/^\//, '').replace(/\//g, '\\'));
}

main().catch((error) => {
  console.error(error.message || error);
  process.exit(1);
});
