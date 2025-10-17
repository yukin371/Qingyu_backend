# 写作端模块设计文档总览

> **版本**: v2.0  
> **创建日期**: 2024-01-01  
> **最后更新**: 2025-10-16  
> **维护者**: 青羽写作团队

## 概述

写作端模块是青羽轻量级阅读写作平台的创作核心，为作者提供完整的写作工具链和创作辅助功能。本模块包含项目管理、文档管理、编辑器、设定百科、AI辅助、版本控制等完整功能体系。

## 模块架构

```plaintext
writing/
├── project/            # 项目管理模块
├── document/           # 文档管理模块
├── editor/             # 编辑器模块
├── version/            # 版本控制模块
├── encyclopedia/       # 设定百科模块
│   ├── character/      # 角色管理
│   ├── worldsetting/   # 世界观设定
│   ├── timeline/       # 时间线
│   └── location/       # 地理设定
├── ai_assistant/       # AI辅助功能
├── content_review/     # 内容审核
├── statistics/         # 数据统计
└── monetization/       # 稿费结算
```

## 核心功能模块

### 1. 项目管理模块 (project)

**功能概述**：
- 项目CRUD操作
- 项目分类和标签管理
- 项目封面和简介管理
- 项目统计（字数、章节数等）
- 项目权限管理
- 项目状态管理（草稿、连载、完结）

**详细设计文档**：[项目管理系统设计.md](./项目管理系统设计.md)

### 2. 文档管理模块 (document)

**功能概述**：
- 文档树形结构管理
- 文档CRUD操作
- 文档拖拽排序
- 文档移动和复制
- 文档层级管理（最多3层）
- 文档类型管理（卷、章、节、场景）

**详细设计文档**：[文档管理系统设计.md](./文档管理系统设计.md)

### 3. 编辑器模块 (editor)

**功能概述**：
- **双模式编辑**: Markdown + 富文本双Tab切换
- **自动保存**: 30秒间隔自动保存
- **文档内容存储**: MongoDB + GridFS大文件支持
- **字数统计**: 实时字数和字符统计
- **图片上传**: 支持图片插入
- **格式转换**: Markdown ↔ 富文本互转

**详细设计文档**：[编辑器系统设计.md](./编辑器系统设计.md)

### 4. 版本控制模块 (version)

**功能概述**：
- **版本历史**: 文档历史版本记录
- **版本对比**: Diff对比功能
- **版本回滚**: 一键恢复到历史版本
- **版本清理**: 自动清理30天前版本
- **版本标签**: 重要版本标记

**详细设计文档**：[编辑器系统设计.md](./编辑器系统设计.md#版本控制)

### 5. 设定百科模块 (encyclopedia)

#### 5.1 角色管理 (character)

**功能概述**：
- 角色卡管理
- 角色关系图谱
- 角色属性和性格设定
- 角色背景故事
- 角色关联查询

**详细设计文档**：[角色卡_关系图设计.md](./角色卡_关系图设计.md)

#### 5.2 世界观设定 (worldsetting)

**功能概述**：
- 多类型设定管理（修炼体系、势力、地理、物品等）
- 设定分类和标签
- 富文本设定内容
- 设定关联关系
- 设定搜索和筛选

**详细设计文档**：[世界观设定管理设计.md](./世界观设定管理设计.md)

#### 5.3 时间线和地图 (timeline & location)

**功能概述**：
- 时间线事件管理
- 地理位置管理
- 空间地图可视化
- 大纲关联

**详细设计文档**：[大纲_时间_空间地图设计.md](./大纲_时间_空间地图设计.md)

### 6. AI辅助功能 (ai_assistant)

**功能概述**：
- **智能续写**: 选中段落一键续写
- **模型切换**: 支持多种AI模型选择
- **创意建议**: 情节发展、角色塑造建议
- **文风分析**: 写作风格分析和改进建议
- **内容优化**: 语法检查、同义词替换

**详细设计文档**：[AI智能辅助系统.md](./AI智能辅助系统.md)

### 7. 内容审核 (content_review)

**功能概述**：
- **敏感词检测**: 实时检测并高亮提示
- **AI审核**: 调用ai-service进行内容审核
- **人工复审**: 疑似内容标记人工审核
- **合规建议**: 提供内容修改建议

### 8. 数据统计 (statistics)

**功能概述**：
- **阅读数据**: 每章完读率、跳章率分析
- **读者画像**: 读者年龄、性别、地域分布
- **收入统计**: 订阅、打赏、广告收入统计
- **趋势分析**: 数据趋势图表展示

### 9. 稿费结算 (monetization)

**功能概述**：
- **多元收入**: 订阅+打赏+广告三分账模式
- **实时结算**: 收入实时统计和分成计算
- **提现管理**: 支持支付宝提现
- **税务处理**: 自动计算税费和开票

## 技术架构

### 数据层
- **MongoDB**: 存储项目、文档、设定、版本等数据
- **GridFS**: 大文档内容存储（>1MB）
- **Redis**: 缓存编辑状态、文档树、统计数据
- **OSS**: 存储图片、附件等静态资源
- **Elasticsearch**: 全文搜索和内容检索（可选）

### 服务层架构

```
┌─────────────────────────────────────────────────────┐
│              Router Layer (路由层)                   │
│           /api/v1/projects/* /documents/*           │
├─────────────────────────────────────────────────────┤
│              API Layer (接口层)                      │
│   ProjectApi / DocumentApi / SettingApi ...         │
├─────────────────────────────────────────────────────┤
│            Service Layer (业务逻辑层)                │
│  ProjectService / DocumentService / SettingService  │
├─────────────────────────────────────────────────────┤
│          Repository Layer (数据访问层)               │
│  ProjectRepository / DocumentRepository ...         │
├─────────────────────────────────────────────────────┤
│            Model Layer (数据模型层)                  │
│     Project / Document / WorldSetting ...           │
├─────────────────────────────────────────────────────┤
│              MongoDB (存储层)                        │
└─────────────────────────────────────────────────────┘
```

### 核心服务模块

- **ProjectService**: 项目管理服务
- **DocumentService**: 文档管理服务
- **DocumentContentService**: 文档内容编辑服务
- **VersionService**: 版本控制服务
- **WorldSettingService**: 世界观设定服务
- **CharacterService**: 角色管理服务
- **AIService**: AI辅助服务
- **ReviewService**: 内容审核服务

## 相关设计文档

### 核心模块设计
- [项目管理系统设计](./项目管理系统设计.md) - 项目CRUD、权限管理、统计
- [文档管理系统设计](./文档管理系统设计.md) - 文档树、拖拽排序、移动
- [编辑器系统设计](./编辑器系统设计.md) - Markdown/富文本编辑、自动保存、版本控制
- [世界观设定管理设计](./世界观设定管理设计.md) - 多类型设定、关联管理

### 辅助功能设计
- [角色卡_关系图设计](./角色卡_关系图设计.md) - 角色管理、关系图谱
- [大纲_时间_空间地图设计](./大纲_时间_空间地图设计.md) - 大纲、时间线、地图
- [AI智能辅助系统](./AI智能辅助系统.md) - AI续写、内容分析

### 实施文档
- [写作端实施文档](../../implementation/04写作端模块/README_写作端实施文档.md)

## 开发优先级

### 第一阶段：项目管理模块（1周）
1. ✅ Project模型与Repository
2. ✅ ProjectService业务逻辑
3. ✅ 文档树结构设计
4. ✅ 文档管理Service
5. ✅ 项目管理API与前端

**预期交付**：项目CRUD、文档树功能

### 第二阶段：文档编辑器（1周）
1. ✅ 文档内容存储设计
2. ✅ 文档编辑API实现
3. ✅ Markdown编辑器集成
4. ✅ 富文本编辑器集成
5. ✅ 版本控制界面

**预期交付**：双编辑器、自动保存、版本控制

### 第三阶段：设定百科系统（3周）
1. 角色卡管理（1周）
2. 世界观设定（1周）
3. 设定模板库（1周）

**预期交付**：角色卡、世界观、模板系统

### MVP完成目标
- ✅ 项目和文档管理
- ✅ 双模式编辑器
- ✅ 自动保存和版本控制
- ✅ 基础设定管理
- 🔄 AI辅助续写
- 🔄 内容审核

### 后续迭代功能
1. 协作编辑功能
2. 完整的关系图谱系统
3. 深度数据分析
4. 完整的稿费结算系统
5. 高级AI辅助功能

## 数据模型概览

### 核心模型

```go
// Project - 项目模型
type Project struct {
    ID          string
    AuthorID    string
    Title       string
    Status      ProjectStatus  // draft | serializing | completed
    Statistics  ProjectStats
    Collaborators []Collaborator
}

// Document - 文档模型
type Document struct {
    ID          string
    ProjectID   string
    ParentID    string         // 树形结构
    Title       string
    Type        DocumentType   // volume | chapter | section
    Level       int            // 0-2层
    Order       int            // 排序
    ContentID   string         // 关联内容
}

// DocumentContent - 文档内容
type DocumentContent struct {
    ID          string
    DocumentID  string
    Content     string
    ContentType string         // markdown | richtext
    WordCount   int
    Version     int
    GridFSID    string         // 大文件存储
}

// Version - 版本记录
type Version struct {
    ID          string
    DocumentID  string
    VersionNum  int
    Content     string
    Comment     string
    CreatedAt   time.Time
}

// WorldSetting - 世界观设定
type WorldSetting struct {
    ID          string
    ProjectID   string
    Type        SettingType    // cultivation | faction | geography...
    Name        string
    Content     string
    Properties  map[string]interface{}
    RelatedSettings []string
}

// Character - 角色
type Character struct {
    ID          string
    ProjectID   string
    Name        string
    Traits      []string
    Background  string
    Relationships []Relationship
}
```

## 集成接口

### 与阅读端集成
- 作品发布到书城
- 章节内容同步
- 读者反馈收集

### 与底层服务集成
- 用户认证和权限管理（JWT）
- 钱包系统对接
- 推荐算法数据提供
- 统计数据上报

### 外部服务集成
- 第三方AI模型API（OpenAI、Claude、文心等）
- 支付宝提现接口
- 敏感词检测服务
- 图片存储服务（OSS）

## 技术规范

### 代码规范
- 遵循[项目开发规则](../../architecture/项目开发规则.md)
- 使用Repository模式进行数据访问
- Service层处理业务逻辑
- API层处理HTTP请求响应
- 统一错误处理和日志记录

### 测试规范
- 单元测试覆盖率≥80%
- Service层测试≥85%
- Repository层测试≥90%
- 集成测试覆盖核心流程

### 性能要求
- 文档加载：<500ms
- 自动保存：<200ms
- 文档树加载：<300ms
- 版本对比：<1s

---

**文档版本**: v2.0  
**创建日期**: 2024-01-01  
**最后更新**: 2025-10-16  
**维护者**: 青羽写作团队
