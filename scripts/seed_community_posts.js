// Community Posts Seeder - 社区动态数据填充脚本
// 使用方式: node seed_community_posts.js

import fs from 'node:fs/promises';

const BASE_URL = process.env.BASE_URL || 'http://localhost:9090';

// 测试用户数据
const TEST_USERS = [
  { id: 'user-test-001', username: '星河漫步', nickname: '星河漫步', avatar: 'https://picsum.photos/seed/user1/100/100', level: 12 },
  { id: 'user-test-002', username: '书虫小窝', nickname: '书虫小窝', avatar: 'https://picsum.photos/seed/user2/100/100', level: 8 },
  { id: 'user-test-003', username: '墨染年华', nickname: '墨染年华', avatar: 'https://picsum.photos/seed/user3/100/100', level: 15 },
  { id: 'user-test-004', username: '云端读者', nickname: '云端读者', avatar: 'https://picsum.photos/seed/user4/100/100', level: 6 },
  { id: 'user-test-005', username: '夜空守望', nickname: '夜空守望', avatar: 'https://picsum.photos/seed/user5/100/100', level: 20 }
];

// 测试动态数据
const TEST_POSTS = [
  { type: 'text', content: '今天看了一本非常精彩的仙侠小说，剧情紧凑，人物塑造也很到位！强烈推荐给大家~', topics: ['仙侠', '推荐'] },
  { type: 'text', content: '有人喜欢看都市异能类的书吗？最近书荒了，求推荐几本好看的！', topics: ['求推荐', '都市'] },
  { type: 'book_recommendation', content: '《星辰大海》这本书真的太赞了！讲述了人类探索宇宙的壮阔历程，看得我热血沸腾！', topics: ['科幻', '推荐'], book: { bookId: 'book-101', title: '星辰大海', cover: 'https://picsum.photos/seed/book101/200/300', author: '银河漫步' } },
  { type: 'reading_progress', content: '终于追完了一本追了三个月的小说！从筑基到飞升，经历了太多太多。感谢作者一路陪伴~', topics: ['读书感悟'], readingProgress: { bookId: 'book-202', chapterId: 'ch-202', chapterTitle: '第520章 大结局', progress: 100 } },
  { type: 'text', content: '周末宅家看书，一口气看了五章停不下来！这种感觉太美妙了，有没有人懂？', topics: ['日常', '阅读'] },
  { type: 'image', content: '分享一下最近入手的新书，封面设计太美了！已经迫不及待想要开始阅读了~', topics: ['晒书'], images: ['https://picsum.photos/seed/book1/400/300', 'https://picsum.photos/seed/book2/400/300'] },
  { type: 'text', content: '修仙小说的套路是不是都差不多啊？退婚、升级、打脸，看多了有点审美疲劳了...', topics: ['吐槽', '修仙'] },
  { type: 'book_recommendation', content: '给大家安利一本冷门好书《雾隐都市》，悬疑氛围营造得特别好，晚上看有点害怕但又停不下来！', topics: ['悬疑', '推荐'], book: { bookId: 'book-303', title: '雾隐都市', cover: 'https://picsum.photos/seed/book303/200/300', author: '暗夜行者' } },
  { type: 'text', content: '今天在书城发现了一本神作！作者文笔太厉害了，寥寥几笔就把人物写活了。', topics: ['惊喜', '推荐'] },
  { type: 'reading_progress', content: '追更《逆天改命》已经一年了，见证了主角从零开始一步步成长为强者，太励志了！', topics: ['追更', '热血'], readingProgress: { bookId: 'book-404', chapterId: 'ch-404', chapterTitle: '第1000章 巅峰之战', progress: 85 } },
  { type: 'text', content: '有没有人和我一样，喜欢在雨天窝在沙发上看书？这种感觉太惬意了~', topics: ['闲聊', '读书日常'] },
  { type: 'text', content: '最近在学习写小说，有没有大神可以指导一下？新手求带！', topics: ['新手求带', '写作'] },
  { type: 'book_recommendation', content: '重温经典《哈利波特》，每次看都有新的感悟。罗琳阿姨的魔法世界真的太伟大了！', topics: ['经典', '奇幻', '推荐'], book: { bookId: 'book-505', title: '哈利波特全集', cover: 'https://picsum.photos/seed/book505/200/300', author: 'J.K.罗琳' } },
  { type: 'text', content: '今天更新了三千字，虽然不多但也是进步！坚持就是胜利~', topics: ['写作', '日常'] },
  { type: 'image', content: '周末去了趟图书馆，发现这个角落超有氛围！顺便借了几本书回家~', topics: ['图书馆', '日常'], images: ['https://picsum.photos/seed/library1/400/300', 'https://picsum.photos/seed/library2/400/300', 'https://picsum.photos/seed/library3/400/300'] }
];

function randomPastTime() {
  const now = Date.now();
  const hoursAgo = Math.floor(Math.random() * 168);
  return new Date(now - hoursAgo * 60 * 60 * 1000);
}

async function seedPosts() {
  console.log('开始填充社区动态数据...');
  console.log(`API地址: ${BASE_URL}`);

  const createdPosts = [];

  for (let i = 0; i < TEST_POSTS.length; i++) {
    const postTemplate = TEST_POSTS[i];
    const user = TEST_USERS[i % TEST_USERS.length];
    const createdAt = randomPastTime();

    const post = {
      userId: user.id,
      userName: user.nickname,
      userAvatar: user.avatar,
      userLevel: user.level,
      type: postTemplate.type,
      content: postTemplate.content,
      images: postTemplate.images || [],
      bookId: postTemplate.book?.bookId || '',
      bookTitle: postTemplate.book?.title || '',
      bookCover: postTemplate.book?.cover || '',
      bookAuthor: postTemplate.book?.author || '',
      chapterId: postTemplate.readingProgress?.chapterId || '',
      chapterTitle: postTemplate.readingProgress?.chapterTitle || '',
      progress: postTemplate.readingProgress?.progress || 0,
      topics: postTemplate.topics || [],
      likeCount: Math.floor(Math.random() * 200) + 10,
      commentCount: Math.floor(Math.random() * 50) + 1,
      shareCount: Math.floor(Math.random() * 20),
      createdAt: createdAt.toISOString(),
      updatedAt: createdAt.toISOString()
    };

    createdPosts.push(post);
    console.log(`  [${i + 1}/${TEST_POSTS.length}] 创建动态: ${post.content.substring(0, 30)}...`);
  }

  console.log('\n数据生成完成！');
  console.log(`共生成 ${createdPosts.length} 条动态`);

  // 保存数据到文件
  const outputPath = './tmp/community_posts_seed.json';
  await fs.mkdir('./tmp', { recursive: true });
  await fs.writeFile(outputPath, JSON.stringify(createdPosts, null, 2));
  console.log(`\n数据已保存到: ${outputPath}`);

  return createdPosts;
}

// 运行
seedPosts()
  .then(() => {
    console.log('\n✅ 社区动态数据生成成功！');
    process.exit(0);
  })
  .catch((err) => {
    console.error('❌ 错误:', err);
    process.exit(1);
  });
