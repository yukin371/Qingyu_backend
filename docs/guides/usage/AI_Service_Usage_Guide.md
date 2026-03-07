# AI服务使用指南

## 概述

青语智能写作系统的AI服务为用户提供了强大的写作辅助功能，包括内容生成、智能续写、文本分析、内容优化等。本指南将详细介绍如何使用这些功能。

## 快速开始

### 1. 环境配置

在开始使用AI服务前，请确保已正确配置环境变量：

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑环境变量文件
# 设置你的AI API密钥和其他配置
```

必需的配置项：
- `AI_API_KEY`: 你的AI服务API密钥
- `AI_PROVIDER`: AI服务提供商（如：openai）
- `AI_BASE_URL`: API基础URL
- `AI_DEFAULT_MODEL`: 默认使用的AI模型

### 2. 启动服务

```bash
# 启动后端服务
go run main.go

# 服务将在配置的端口启动（默认8080）
```

## 功能详解

### 1. 智能内容生成

#### 功能描述
基于用户提供的提示词和项目上下文，生成符合要求的文本内容。

#### 使用场景
- 创作新的章节内容
- 生成角色对话
- 创建场景描述
- 补充情节细节

#### 使用方法

**前端调用示例**：
```javascript
// 生成内容
const generateContent = async (projectId, prompt, options = {}) => {
  try {
    const response = await fetch('/api/v1/ai/generate', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken}`
      },
      body: JSON.stringify({
        projectId: projectId,
        prompt: prompt,
        options: {
          temperature: options.temperature || 0.7,
          maxTokens: options.maxTokens || 1000,
          style: options.style || 'narrative',
          genre: options.genre || 'general'
        }
      })
    });
    
    const result = await response.json();
    if (result.code === 0) {
      return result.data.content;
    } else {
      throw new Error(result.message);
    }
  } catch (error) {
    console.error('生成内容失败:', error);
    throw error;
  }
};

// 使用示例
const content = await generateContent(
  'project123',
  '描述主角第一次见到反派时的心理活动',
  {
    temperature: 0.8,
    maxTokens: 500,
    style: 'psychological',
    genre: 'fantasy'
  }
);
```

**参数说明**：
- `temperature`: 创造性程度（0-1），值越高越有创意
- `maxTokens`: 生成内容的最大长度
- `style`: 写作风格（narrative/dialogue/descriptive/psychological等）
- `genre`: 文体类型（fantasy/romance/mystery/sci-fi等）

### 2. 智能续写

#### 功能描述
基于当前文本内容，智能续写后续情节，保持风格和逻辑的连贯性。

#### 使用场景
- 章节续写
- 情节推进
- 对话延续
- 场景扩展

#### 使用方法

```javascript
const continueWriting = async (projectId, chapterId, currentText, length = 500) => {
  try {
    const response = await fetch('/api/v1/ai/continue', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken}`
      },
      body: JSON.stringify({
        projectId: projectId,
        chapterId: chapterId,
        currentText: currentText,
        continueLength: length,
        options: {
          temperature: 0.7,
          maxTokens: length * 2 // 预留更多token
        }
      })
    });
    
    const result = await response.json();
    return result.data.content;
  } catch (error) {
    console.error('续写失败:', error);
    throw error;
  }
};

// 使用示例
const currentChapterText = "主角走进了神秘的森林，四周静得可怕...";
const continuation = await continueWriting(
  'project123',
  'chapter1',
  currentChapterText,
  800
);
```

### 3. 文本分析

#### 功能描述
分析文本内容，提供情节、角色、风格等方面的专业分析和建议。

#### 分析类型
- `plot`: 情节分析
- `character`: 角色分析
- `style`: 风格分析
- `general`: 综合分析

#### 使用方法

```javascript
const analyzeContent = async (content, analysisType = 'general') => {
  try {
    const response = await fetch('/api/v1/ai/analyze', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        content: content,
        analysisType: analysisType
      })
    });
    
    const result = await response.json();
    return result.data.analysis;
  } catch (error) {
    console.error('分析失败:', error);
    throw error;
  }
};

// 使用示例
const chapterContent = "这里是你的章节内容...";

// 情节分析
const plotAnalysis = await analyzeContent(chapterContent, 'plot');
console.log('情节分析:', plotAnalysis);

// 角色分析
const characterAnalysis = await analyzeContent(chapterContent, 'character');
console.log('角色分析:', characterAnalysis);
```

### 4. 内容优化

#### 功能描述
优化文本的语法、风格、流畅度，提升文本质量。

#### 优化类型
- `grammar`: 语法优化
- `style`: 风格优化
- `flow`: 流畅度优化
- `dialogue`: 对话优化

#### 使用方法

```javascript
const optimizeText = async (projectId, originalText, optimizeType = 'general', instructions = '') => {
  try {
    const response = await fetch('/api/v1/ai/optimize', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken}`
      },
      body: JSON.stringify({
        projectId: projectId,
        originalText: originalText,
        optimizeType: optimizeType,
        instructions: instructions,
        options: {
          temperature: 0.3, // 优化时使用较低的创造性
          maxTokens: originalText.length * 2
        }
      })
    });
    
    const result = await response.json();
    return result.data.content;
  } catch (error) {
    console.error('优化失败:', error);
    throw error;
  }
};

// 使用示例
const originalText = "他走的很快，因为他很着急。";
const optimizedText = await optimizeText(
  'project123',
  originalText,
  'style',
  '请让这段文字更生动，增加细节描述'
);
```

### 5. 大纲生成

#### 功能描述
基于主题和要求生成详细的故事大纲。

#### 使用方法

```javascript
const generateOutline = async (projectId, theme, options = {}) => {
  try {
    const response = await fetch('/api/v1/ai/outline', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken}`
      },
      body: JSON.stringify({
        projectId: projectId,
        theme: theme,
        genre: options.genre || 'general',
        length: options.length || 'medium',
        keyElements: options.keyElements || [],
        options: {
          temperature: 0.8,
          maxTokens: 2000
        }
      })
    });
    
    const result = await response.json();
    return result.data.content;
  } catch (error) {
    console.error('大纲生成失败:', error);
    throw error;
  }
};

// 使用示例
const outline = await generateOutline(
  'project123',
  '一个关于时间旅行的科幻故事',
  {
    genre: 'sci-fi',
    length: 'long',
    keyElements: ['时间机器', '平行宇宙', '蝴蝶效应']
  }
);
```

## 上下文管理

### 获取项目上下文

AI服务会自动获取项目的上下文信息，包括：
- 当前章节内容
- 活跃角色信息
- 场景设定
- 世界观设定
- 情节线索

```javascript
const getContext = async (projectId, chapterId = null) => {
  try {
    const url = chapterId 
      ? `/api/v1/ai/context/${projectId}/${chapterId}`
      : `/api/v1/ai/context/${projectId}`;
      
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${userToken}`
      }
    });
    
    const result = await response.json();
    return result.data;
  } catch (error) {
    console.error('获取上下文失败:', error);
    throw error;
  }
};
```

### 更新上下文反馈

```javascript
const updateContextFeedback = async (projectId, chapterId, feedback) => {
  try {
    const response = await fetch('/api/v1/ai/context/feedback', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken}`
      },
      body: JSON.stringify({
        projectId: projectId,
        chapterId: chapterId,
        feedback: feedback
      })
    });
    
    const result = await response.json();
    return result.code === 0;
  } catch (error) {
    console.error('更新反馈失败:', error);
    throw error;
  }
};
```

## 最佳实践

### 1. 提示词优化

**好的提示词示例**：
```
描述主角李明在雨夜中追赶逃犯的场景。要求：
- 突出紧张刺激的氛围
- 描写雨水和夜色的环境
- 展现主角的决心和技能
- 长度约300字
```

**避免的提示词**：
```
写一段追逐戏
```

### 2. 参数调优

- **创意写作**：temperature = 0.7-0.9
- **技术文档**：temperature = 0.1-0.3
- **对话生成**：temperature = 0.6-0.8
- **内容优化**：temperature = 0.2-0.4

### 3. 错误处理

```javascript
const safeAICall = async (apiCall, retries = 3) => {
  for (let i = 0; i < retries; i++) {
    try {
      return await apiCall();
    } catch (error) {
      if (i === retries - 1) throw error;
      
      // 根据错误类型决定是否重试
      if (error.status === 429) {
        // 限流错误，等待后重试
        await new Promise(resolve => setTimeout(resolve, 1000 * (i + 1)));
      } else if (error.status >= 500) {
        // 服务器错误，短暂等待后重试
        await new Promise(resolve => setTimeout(resolve, 500));
      } else {
        // 客户端错误，不重试
        throw error;
      }
    }
  }
};
```

### 4. 性能优化

```javascript
// 使用缓存避免重复请求
const cache = new Map();

const cachedAICall = async (cacheKey, apiCall) => {
  if (cache.has(cacheKey)) {
    return cache.get(cacheKey);
  }
  
  const result = await apiCall();
  cache.set(cacheKey, result);
  
  // 设置缓存过期时间
  setTimeout(() => cache.delete(cacheKey), 5 * 60 * 1000); // 5分钟
  
  return result;
};
```

## 常见问题解决

### Q1: AI生成的内容质量不高怎么办？

**解决方案**：
1. 优化提示词，提供更详细的要求
2. 确保项目上下文信息完整
3. 调整temperature参数
4. 使用更高级的AI模型

### Q2: API调用频繁失败

**解决方案**：
1. 检查API密钥是否正确
2. 确认网络连接正常
3. 检查是否触发限流
4. 实现重试机制

### Q3: 生成内容与预期不符

**解决方案**：
1. 提供更具体的提示词
2. 使用示例来指导生成
3. 分步骤生成复杂内容
4. 利用内容优化功能

### Q4: 如何控制生成成本？

**解决方案**：
1. 合理设置maxTokens参数
2. 使用缓存避免重复请求
3. 选择合适的AI模型
4. 批量处理相关请求

## 进阶用法

### 1. 批量处理

```javascript
const batchProcess = async (items, processor, batchSize = 5) => {
  const results = [];
  
  for (let i = 0; i < items.length; i += batchSize) {
    const batch = items.slice(i, i + batchSize);
    const batchPromises = batch.map(processor);
    const batchResults = await Promise.all(batchPromises);
    results.push(...batchResults);
    
    // 避免过快请求
    if (i + batchSize < items.length) {
      await new Promise(resolve => setTimeout(resolve, 1000));
    }
  }
  
  return results;
};
```

### 2. 流式生成

```javascript
const streamGenerate = async (projectId, prompt, onChunk) => {
  const response = await fetch('/api/v1/ai/generate/stream', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${userToken}`
    },
    body: JSON.stringify({
      projectId: projectId,
      prompt: prompt,
      stream: true
    })
  });
  
  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  
  while (true) {
    const { done, value } = await reader.read();
    if (done) break;
    
    const chunk = decoder.decode(value);
    onChunk(chunk);
  }
};
```

## 更新说明

### v1.0.0 (2024-01-01)
- 初始版本发布
- 支持基础AI功能
- 提供完整的使用文档

---

如有其他问题，请参考API文档或联系技术支持。