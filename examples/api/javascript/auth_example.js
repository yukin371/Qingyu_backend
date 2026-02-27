/**
 * 青羽写作平台 API - 认证示例
 * 演示如何使用JavaScript/Node.js进行API认证和调用
 */

const axios = require('axios');

class QingyuAPIClient {
  constructor(baseURL = 'http://localhost:9090/api/v1') {
    this.baseURL = baseURL.replace(/\/$/, '');
    this.token = null;
    this.client = axios.create({
      baseURL: this.baseURL,
      headers: {
        'Content-Type': 'application/json'
      }
    });

    // 请求拦截器 - 自动添加Token
    this.client.interceptors.request.use(config => {
      if (this.token) {
        config.headers.Authorization = `Bearer ${this.token}`;
      }
      return config;
    });

    // 响应拦截器 - 统一错误处理
    this.client.interceptors.response.use(
      response => response.data,
      error => {
        console.error('请求失败:', error.message);
        if (error.response) {
          console.error('响应内容:', error.response.data);
        }
        throw error;
      }
    );
  }

  async register(username, email, password, nickname) {
    console.log('\n=== 用户注册 ===');
    try {
      const response = await this.client.post('/auth/register', {
        username,
        email,
        password,
        nickname
      });
      console.log('注册成功:', response);
      return response;
    } catch (error) {
      throw error;
    }
  }

  async login(username, password) {
    console.log('\n=== 用户登录 ===');
    try {
      const response = await this.client.post('/auth/login', {
        username,
        password
      });

      if (response.data && response.data.token) {
        this.token = response.data.token;
        console.log(`登录成功，Token: ${this.token.substring(0, 20)}...`);
      } else if (response.data && response.data.access_token) {
        this.token = response.data.access_token;
        console.log(`登录成功，Token: ${this.token.substring(0, 20)}...`);
      } else {
        console.log('登录成功，但未获取到Token');
      }

      return response;
    } catch (error) {
      throw error;
    }
  }

  async getProfile() {
    console.log('\n=== 获取用户信息 ===');
    try {
      const response = await this.client.get('/user/profile');
      console.log('用户信息:', JSON.stringify(response, null, 2));
      return response;
    } catch (error) {
      throw error;
    }
  }

  async refreshToken() {
    console.log('\n=== 刷新Token ===');
    try {
      const response = await this.client.post('/auth/refresh');
      if (response.data && response.data.token) {
        this.token = response.data.token;
        console.log(`Token刷新成功: ${this.token.substring(0, 20)}...`);
      }
      return response;
    } catch (error) {
      throw error;
    }
  }

  async logout() {
    console.log('\n=== 登出 ===');
    try {
      const response = await this.client.post('/auth/logout');
      this.token = null;
      console.log('登出成功');
      return response;
    } catch (error) {
      throw error;
    }
  }
}

// 主函数 - 演示认证流程
async function main() {
  const client = new QingyuAPIClient();

  try {
    // 1. 用户注册
    console.log('\n【步骤1】用户注册');
    try {
      await client.register(
        'testuser_js',
        'testjs@example.com',
        'SecurePass123!',
        'JavaScript测试用户'
      );
    } catch (error) {
      console.log('注册可能失败（用户已存在）:', error.message);
    }

    // 2. 用户登录
    console.log('\n【步骤2】用户登录');
    await client.login('testuser_js', 'SecurePass123!');

    // 3. 访问受保护的API
    console.log('\n【步骤3】访问受保护的API');
    if (client.token) {
      await client.getProfile();
    } else {
      console.log('未获取到Token，跳过');
    }

    // 4. 刷新Token
    console.log('\n【步骤4】刷新Token');
    if (client.token) {
      await client.refreshToken();
    }

    // 5. 登出
    console.log('\n【步骤5】登出');
    if (client.token) {
      await client.logout();
    }

    console.log('\n=== 认证示例完成 ===');
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

module.exports = { QingyuAPIClient };
