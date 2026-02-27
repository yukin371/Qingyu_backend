/**
 * 青羽写作平台 API - 管理员示例
 * 演示如何使用JavaScript调用管理员API
 */

const { AdminAPIClient } = require('./api_client');

async function main() {
  // 从环境变量获取管理员token
  const adminToken = process.env.QINGYU_ADMIN_TOKEN;

  if (!adminToken) {
    console.error('错误: 未设置管理员Token');
    console.log('请设置环境变量: export QINGYU_ADMIN_TOKEN=your_admin_token');
    return 1;
  }

  // 创建API客户端
  const client = new AdminAPIClient('http://localhost:9090/api/v1', adminToken);

  try {
    // 1. 获取用户列表
    await client.getUsers(1, 10);

    // 2. 创建用户
    try {
      await client.createUser(
        'newuser_js',
        'newuser_js@example.com',
        'SecurePass123!',
        'JavaScript创建的用户',
        'user'
      );
    } catch (error) {
      console.log('创建用户失败（可能已存在）:', error.message);
    }

    // 3. 获取权限列表
    await client.getPermissions();

    // 4. 创建权限
    try {
      await client.createPermission(
        'test.permission.js',
        'JavaScript测试权限',
        '通过JavaScript脚本创建的测试权限',
        'test'
      );
    } catch (error) {
      console.log('创建权限失败（可能已存在）:', error.message);
    }

    // 5. 获取用户增长趋势
    await client.getUserGrowthTrend('2026-01-01', '2026-01-31', 'daily');

    // 6. 获取内容统计
    await client.getContentStatistics();

    // 7. 获取系统概览
    await client.getSystemOverview();

    // 8. 创建公告
    await client.createAnnouncement(
      'JavaScript测试公告',
      '这是一个通过JavaScript脚本创建的测试公告',
      'info',
      'normal',
      false
    );

    // 9. 获取公告列表
    await client.getAnnouncements();

    console.log('\n=== 管理员API示例完成 ===');
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
