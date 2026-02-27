/**
 * 青羽写作平台 API - 书城示例
 * 演示如何使用JavaScript调用书城相关API
 */

const { BookstoreAPIClient } = require('./api_client');

async function main() {
  // 从环境变量获取token
  const token = process.env.QINGYU_TOKEN || 'your-jwt-token-here';

  if (token === 'your-jwt-token-here') {
    console.log('警告: 未设置Token，部分API将无法访问');
    console.log('请设置环境变量: export QINGYU_TOKEN=your_token');
  }

  // 创建API客户端
  const client = new BookstoreAPIClient('http://localhost:9090/api/v1', token);

  try {
    // 1. 获取首页数据
    await client.getHomeData();

    // 2. 获取书籍列表
    await client.getBooks(1, 5);

    // 3. 搜索书籍
    await client.searchByTitle('玄幻');

    // 4. 按作者搜索
    await client.searchByAuthor('唐家三少');

    // 5. 按标签筛选
    await client.filterByTags(['玄幻', '修真']);

    // 6. 获取书籍详情
    const bookId = 'book123';
    await client.getBookDetail(bookId);

    // 7. 获取章节列表
    await client.getBookChapters(bookId);

    // 8. 获取章节内容（需要认证）
    if (token !== 'your-jwt-token-here') {
      const chapterId = 'chapter123';
      await client.getChapterContent(chapterId);
    }

    // 9. 获取相似书籍
    await client.getSimilarBooks(bookId);

    // 10. 提交评分（需要认证）
    if (token !== 'your-jwt-token-here') {
      await client.rateBook(bookId, 5, '非常好看！');
    }

    console.log('\n=== 书城API示例完成 ===');
    return 0;
  } catch (error) {
    console.error('\n错误:', error.message);
    return 1;
  }
}

// 如果直接运行此文件
if (require.main === module) {
  main().then(exitCode => {
    process.exit(exitCode);
  });
}

module.exports = { main };
