# MCP 浏览器服务器使用说明

## Page Role

- legacy-runbook
- current-owner: `docs/implementation/infrastructure/`
- current-bounded: MCP 浏览器服务器历史使用说明，只记录当时的配置与使用方式

## Recommended Read Path

1. 先读 `README.md`。
2. 需要回看浏览器 MCP 配置方式时，再读本文件。

## Boundary

- 本页是历史使用说明，不是当前工具 owner。
- 服务器名、包名和使用方式都带明显时间性。

## Quick Section Map

- 已配置的浏览器 MCP 服务器
- 使用方法
- 使用示例
- 配置文件位置
- 故障排除
- 更多资源
- 注意事项

## Quick Takeaways

- 这是浏览器 MCP 的历史使用说明。

## Skip Guide

- 只看当前工具 owner：跳过本文件。

## 📋 已配置的浏览器 MCP 服务器

### 1. Puppeteer MCP Server
- **服务器名**: `puppeteer`
- **包**: `@executeautomation/puppeteer-mcp-server`
- **模式**: Headless（无头模式）

### 2. Playwright MCP Server
- **服务器名**: `playwright`
- **包**: `@executeautomation/playwright-mcp-server`
- **模式**: Headless（无头模式）

## 🚀 使用方法

### 启动服务器
MCP 服务器会在你重启 Claude Code 后自动启动。重启后，你就可以使用浏览器功能了。

### 可用功能

#### Puppeteer/Playwright 功能：
1. **导航到网页**
2. **截图**
3. **点击元素**
4. **填写表单**
5. **提取页面内容**
6. **执行 JavaScript**
7. **页面交互**

## 💡 使用示例

### 示例 1：访问网页并截图
```
请帮我访问 https://www.google.com 并截图
```

### 示例 2：搜索内容
```
请用 puppeteer 打开 Google，搜索 "Claude AI"，然后截图
```

### 示例 3：提取页面信息
```
请访问 https://example.com 并提取页面标题
```

### 示例 4：自动化操作
```
请打开登录页面，填写用户名和密码，然后点击登录按钮
```

## ⚙️ 配置文件位置

1. **项目级别**: `D:/Github/青羽/.mcp.json`
2. **全局级别**: `C:/Users/yukin/.claude/mcp_servers.json`

## 🔧 故障排除

### 如果浏览器 MCP 无法启动：

1. **检查 Node.js 和 npm**
   ```bash
   node --version
   npm --version
   ```

2. **手动测试安装**
   ```bash
   npx -y @executeautomation/puppeteer-mcp-server
   ```

3. **查看日志**
   - 检查 Claude Code 的 MCP 服务器日志

4. **重启 Claude Code**
   - 完全关闭 Claude Code
   - 重新打开

## 📚 更多资源

- [Puppeteer MCP Server](https://github.com/executeautomation/puppeteer-mcp-server)
- [Playwright MCP Server](https://github.com/executeautomation/playwright-mcp-server)
- [Model Context Protocol](https://modelcontextprotocol.io/)

## ⚠️ 注意事项

1. **首次使用**：首次使用时会自动下载浏览器二进制文件，可能需要一些时间
2. **网络要求**：需要能够访问 npm registry
3. **权限**：某些操作可能需要用户权限批准
4. **资源占用**：浏览器会占用一定的系统资源
