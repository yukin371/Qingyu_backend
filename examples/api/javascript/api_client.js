/**
 * 青羽写作平台 API - 通用客户端
 * 提供书城和管理员API的封装
 */

const { QingyuAPIClient: BaseClient } = require('./auth_example');

class BookstoreAPIClient extends BaseClient {
  constructor(baseURL, token) {
    super(baseURL);
    this.token = token;
  }

  async getHomeData() {
    console.log('\n=== 获取首页数据 ===');
    const response = await this.client.get('/bookstore/home');
    console.log('首页数据:', JSON.stringify(response, null, 2));
    return response;
  }

  async getBooks(page = 1, limit = 10) {
    console.log(`\n=== 获取书籍列表 (第${page}页) ===`);
    const response = await this.client.get('/bookstore/books', {
      params: { page, limit }
    });
    if (response.data && response.data.books) {
      console.log(`获取到 ${response.data.books.length} 本书籍`);
    }
    return response;
  }

  async searchByTitle(title, page = 1, limit = 5) {
    console.log(`\n=== 搜索书籍（标题: ${title}） ===`);
    const response = await this.client.get('/bookstore/books/search/title', {
      params: { title, page, limit }
    });
    if (response.data && response.data.books) {
      console.log(`找到 ${response.data.books.length} 本相关书籍`);
    }
    return response;
  }

  async searchByAuthor(author, page = 1, limit = 5) {
    console.log(`\n=== 搜索作者（作者: ${author}） ===`);
    const response = await this.client.get('/bookstore/books/search/author', {
      params: { author, page, limit }
    });
    if (response.data && response.data.books) {
      console.log(`找到 ${response.data.books.length} 本相关书籍`);
    }
    return response;
  }

  async filterByTags(tags, page = 1, limit = 5) {
    console.log(`\n=== 按标签筛选（标签: ${tags.join(', ')}） ===`);
    const response = await this.client.get('/bookstore/books/tags', {
      params: { tags: tags.join(','), page, limit }
    });
    if (response.data && response.data.books) {
      console.log(`找到 ${response.data.books.length} 本相关书籍`);
    }
    return response;
  }

  async getBookDetail(bookId) {
    console.log(`\n=== 获取书籍详情 (ID: ${bookId}) ===`);
    const response = await this.client.get(`/bookstore/books/${bookId}`);
    console.log('书籍标题:', response.data?.title || 'N/A');
    return response;
  }

  async getBookChapters(bookId, page = 1, limit = 20) {
    console.log(`\n=== 获取章节列表 (书籍ID: ${bookId}) ===`);
    const response = await this.client.get(`/bookstore/books/${bookId}/chapters`, {
      params: { page, limit }
    });
    if (response.data && response.data.chapters) {
      console.log(`获取到 ${response.data.chapters.length} 个章节`);
    }
    return response;
  }

  async getChapterContent(chapterId) {
    console.log(`\n=== 获取章节内容 (ID: ${chapterId}) ===`);
    const response = await this.client.get(`/bookstore/chapters/${chapterId}`);
    if (response.data && response.data.content) {
      console.log(`章节内容预览: ${response.data.content.substring(0, 100)}...`);
    }
    return response;
  }

  async getSimilarBooks(bookId, limit = 10) {
    console.log(`\n=== 获取相似书籍 (书籍ID: ${bookId}) ===`);
    const response = await this.client.get(`/bookstore/books/${bookId}/similar`, {
      params: { limit }
    });
    if (response.data && response.data.books) {
      console.log(`找到 ${response.data.books.length} 本相似书籍`);
    }
    return response;
  }

  async rateBook(bookId, rating, comment = '') {
    console.log(`\n=== 提交评分 (书籍ID: ${bookId}, 评分: ${rating}) ===`);
    const data = { rating };
    if (comment) {
      data.comment = comment;
    }
    const response = await this.client.post(`/bookstore/books/${bookId}/rating`, data);
    console.log('评分提交成功');
    return response;
  }
}

class AdminAPIClient extends BaseClient {
  constructor(baseURL, adminToken) {
    super(baseURL);
    this.token = adminToken;
  }

  // ==================== 用户管理 ====================

  async getUsers(page = 1, limit = 20, status = null) {
    console.log(`\n=== 获取用户列表 (第${page}页) ===`);
    const params = { page, limit };
    if (status) {
      params.status = status;
    }
    const response = await this.client.get('/admin/users', { params });
    if (response.data && response.data.users) {
      console.log(`获取到 ${response.data.users.length} 个用户`);
    }
    return response;
  }

  async createUser(username, email, password, nickname, role = 'user') {
    console.log(`\n=== 创建用户 (${username}) ===`);
    const response = await this.client.post('/admin/users', {
      username,
      email,
      password,
      nickname,
      role
    });
    console.log('用户创建成功');
    return response;
  }

  // ==================== 权限管理 ====================

  async getPermissions() {
    console.log('\n=== 获取权限列表 ===');
    const response = await this.client.get('/admin/permissions');
    if (response.data) {
      console.log(`获取到 ${response.data.length} 个权限`);
    }
    return response;
  }

  async createPermission(code, name, description, category) {
    console.log(`\n=== 创建权限 (${code}) ===`);
    const response = await this.client.post('/admin/permissions', {
      code,
      name,
      description,
      category
    });
    console.log('权限创建成功');
    return response;
  }

  // ==================== 统计分析 ====================

  async getUserGrowthTrend(startDate, endDate, interval = 'daily') {
    console.log(`\n=== 获取用户增长趋势 (${startDate} 至 ${endDate}) ===`);
    const response = await this.client.get('/admin/analytics/user-growth', {
      params: { start_date: startDate, end_date: endDate, interval }
    });
    if (response.data && response.data.trend) {
      console.log(`获取到 ${response.data.trend.length} 条趋势数据`);
    }
    return response;
  }

  async getContentStatistics(startDate = null, endDate = null) {
    console.log('\n=== 获取内容统计 ===');
    const params = {};
    if (startDate) params.start_date = startDate;
    if (endDate) params.end_date = endDate;
    const response = await this.client.get('/admin/analytics/content-statistics', { params });
    console.log('内容统计获取成功');
    return response;
  }

  async getSystemOverview() {
    console.log('\n=== 获取系统概览 ===');
    const response = await this.client.get('/admin/analytics/system-overview');
    console.log('系统概览获取成功');
    return response;
  }

  // ==================== 公告管理 ====================

  async createAnnouncement(title, content, type, priority = 'normal', isPinned = false) {
    console.log(`\n=== 创建公告 (${title}) ===`);
    const response = await this.client.post('/admin/announcements', {
      title,
      content,
      type,
      priority,
      is_pinned: isPinned
    });
    console.log('公告创建成功');
    return response;
  }

  async getAnnouncements() {
    console.log('\n=== 获取公告列表 ===');
    const response = await this.client.get('/admin/announcements');
    if (response.data && response.data.announcements) {
      console.log(`获取到 ${response.data.announcements.length} 个公告`);
    }
    return response;
  }
}

module.exports = { BookstoreAPIClient, AdminAPIClient };
