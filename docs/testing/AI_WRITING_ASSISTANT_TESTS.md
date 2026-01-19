# 青羽写作平台 AI辅助功能 - 单元测试文档

## 概述

为青羽写作平台的AI辅助功能创建了完整的单元测试套件，包括内容总结、文本校对和敏感词检测服务。

## 测试文件清单

### 1. Mock 适配器
**文件**: `D:\Github\青羽\Qingyu_backend\service\ai\mocks\ai_adapter_mock.go`

**功能**:
- 模拟 AI 适配器行为
- 支持成功/失败/超时场景
- 可配置响应延迟和错误
- 记录调用次数和最后请求

**主要方法**:
```go
- NewMockAIAdapter(name string) *MockAIAdapter
- SetTextResponse(text string, tokens int)
- SetChatResponse(content string, tokens int)
- Reset()
```

### 2. 内容总结服务测试
**文件**: `D:\Github\青羽\Qingyu_backend\service\ai\summarize_service_test.go`

**测试场景**:

#### ✅ 基础功能测试
- `TestSummarizeService_SummarizeContent_Success` - 成功总结内容
- `TestSummarizeService_SummarizeContent_EmptyContent` - 空内容验证
- `TestSummarizeService_SummarizeContent_WhitespaceContent` - 仅空白字符验证

#### ✅ 总结类型测试
- `TestSummarizeService_SummarizeContent_DifferentTypes` - 测试不同总结类型
  - brief (简短摘要)
  - detailed (详细摘要)
  - keypoints (关键点提取)
  - 默认类型

#### ✅ 功能特性测试
- `TestSummarizeService_SummarizeContent_WithQuotes` - 包含引用的总结
- `TestSummarizeService_ExtractKeyPoints` - 关键点提取逻辑
- `TestSummarizeService_SummarizeContent_CompressionRate` - 压缩率计算

#### ✅ 错误场景测试
- `TestSummarizeService_SummarizeContent_AIError` - AI服务错误
- `TestSummarizeService_SummarizeContent_AITimeout` - AI服务超时

#### ✅ 章节总结测试
- `TestSummarizeService_SummarizeChapter_Success` - 章节总结成功
- `TestSummarizeService_SummarizeChapter_ChapterIDRequired` - 章节ID验证

#### 📊 性能测试
- `BenchmarkSummarizeService_ExtractKeyPoints` - 关键点提取性能

### 3. 文本校对服务测试
**文件**: `D:\Github\青羽\Qingyu_backend\service\ai\proofread_service_test.go`

**测试场景**:

#### ✅ 基础功能测试
- `TestProofreadService_ProofreadContent_Success` - 成功校对内容
- `TestProofreadService_ProofreadContent_EmptyContent` - 空内容验证
- `TestProofreadService_ProofreadContent_WhitespaceContent` - 空白内容验证

#### ✅ 检查类型测试
- `TestProofreadService_ProofreadContent_DefaultCheckTypes` - 默认检查类型
- `TestProofreadService_ProofreadContent_CustomCheckTypes` - 自定义检查类型
  - spelling (拼写)
  - grammar (语法)
  - punctuation (标点)
  - style (风格)

#### ✅ 结果解析测试
- `TestProofreadService_ParseProofreadResult_JSONFormat` - JSON格式解析
- `TestProofreadService_ParseProofreadResult_TextFormat` - 文本格式解析（后备方案）
- `TestProofreadService_ExtractIssuesFromText` - 从文本提取问题

#### ✅ 统计和评分测试
- `TestProofreadService_GenerateStatistics` - 统计信息生成
- `TestProofreadService_CalculateScore` - 评分计算
  - 完美内容: 100分
  - 仅错误: 每个-5分
  - 仅警告: 每个-2分
  - 仅建议: 每个-0.5分

#### ✅ 辅助功能测试
- `TestProofreadService_FindPositionInText` - 文本位置查找
- `TestProofreadService_GetProofreadSuggestion` - 获取校对建议
- `TestProofreadService_LongText` - 长文本处理

#### 📊 性能测试
- `BenchmarkProofreadService_CalculateScore` - 评分计算性能

### 4. 敏感词检测服务测试
**文件**: `D:\Github\青羽\Qingyu_backend\service\ai\sensitive_words_service_test.go`

**测试场景**:

#### ✅ 基础功能测试
- `TestSensitiveWordsService_CheckSensitiveWords_Success` - 成功检测敏感词
- `TestSensitiveWordsService_CheckSensitiveWords_EmptyContent` - 空内容验证
- `TestSensitiveWordsService_CheckSensitiveWords_NoMatch` - 未检测到敏感词

#### ✅ 分类检测测试
- `TestSensitiveWordsService_CheckSensitiveWords_PoliticalCategory` - 政治敏感词
- `TestSensitiveWordsService_CheckSensitiveWords_ViolenceCategory` - 暴力敏感词
- `TestSensitiveWordsService_CheckSensitiveWords_AdultCategory` - 成人内容敏感词
- `TestSensitiveWordsService_CheckSensitiveWords_AllCategories` - 所有分类

#### ✅ 自定义词库测试
- `TestSensitiveWordsService_CheckSensitiveWords_CustomWords` - 自定义敏感词
- `TestSensitiveWordsService_AddCustomWords` - 添加自定义词
- `TestSensitiveWordsService_RemoveCustomWords` - 移除自定义词

#### ✅ 检测逻辑测试
- `TestSensitiveWordsService_DetectSensitiveWords` - 敏感词检测逻辑
- `TestSensitiveWordsService_FindWordPositions` - 查找词位置
  - 单次出现
  - 多次出现
  - 未出现
  - 中文词组

#### ✅ 位置和上下文测试
- `TestSensitiveWordsService_CalculateLineColumn` - 行列计算
- `TestSensitiveWordsService_ExtractContext` - 上下文提取

#### ✅ 风险级别测试
- `TestSensitiveWordsService_DetermineWordLevel` - 风险级别确定
  - political: high
  - violence: medium
  - adult: high
  - custom: medium
  - unknown: low

#### ✅ 统计和分析测试
- `TestSensitiveWordsService_GenerateSuggestion` - 生成修改建议
- `TestSensitiveWordsService_GenerateCheckSummary` - 生成检测摘要
- `TestSensitiveWordsService_HasHighRiskWords` - 高风险词检测
- `TestSensitiveWordsService_GetSensitiveWordsDetail` - 获取检测详情

#### 📊 性能测试
- `BenchmarkSensitiveWordsService_FindWordPositions` - 查找位置性能
- `BenchmarkSensitiveWordsService_DetectSensitiveWords` - 检测性能

### 5. API层测试
**文件**: `D:\Github\青羽\Qingyu_backend\api\v1\ai\writing_assistant_api_test.go`

**测试场景**:

#### ✅ 内容总结API测试
- `TestWritingAssistantApi_SummarizeContent_Success` - 成功总结
- `TestWritingAssistantApi_SummarizeContent_InvalidJSON` - 无效JSON
- `TestWritingAssistantApi_SummarizeContent_EmptyContent` - 空内容

#### ✅ 章节总结API测试
- `TestWritingAssistantApi_SummarizeChapter_Success` - 成功总结章节
- `TestWritingAssistantApi_SummarizeChapter_MissingChapterID` - 缺少章节ID

#### ✅ 文本校对API测试
- `TestWritingAssistantApi_ProofreadContent_Success` - 成功校对
- `TestWritingAssistantApi_ProofreadContent_EmptyContent` - 空内容
- `TestWritingAssistantApi_GetProofreadSuggestion_Success` - 获取建议
- `TestWritingAssistantApi_GetProofreadSuggestion_EmptyID` - 空建议ID

#### ✅ 敏感词检测API测试
- `TestWritingAssistantApi_CheckSensitiveWords_Success` - 成功检测
- `TestWritingAssistantApi_CheckSensitiveWords_EmptyContent` - 空内容
- `TestWritingAssistantApi_GetSensitiveWordsDetail_Success` - 获取详情
- `TestWritingAssistantApi_GetSensitiveWordsDetail_EmptyID` - 空检测ID

#### ✅ HTTP协议测试
- `TestWritingAssistantApi_ResponseHeaders` - 响应头验证
- `TestWritingAssistantApi_HTTPMethods` - HTTP方法验证
- `TestWritingAssistantApi_MissingContentType` - 缺少Content-Type

#### ✅ 集成测试
- `TestWritingAssistantApi_Integration` - 完整API请求流程

#### 📊 性能测试
- `BenchmarkWritingAssistantApi_SummarizeContent` - API性能测试

## 测试命令

### 运行所有测试
```bash
cd D:\Github\青羽\Qingyu_backend
go test -v ./service/ai/... ./api/v1/ai/... -run "Test.*Service.*|TestWritingAssistantApi.*"
```

### 运行特定服务测试
```bash
# 内容总结服务
go test -v ./service/ai -run "TestSummarizeService.*"

# 文本校对服务
go test -v ./service/ai -run "TestProofreadService.*"

# 敏感词检测服务
go test -v ./service/ai -run "TestSensitiveWordsService.*"
```

### 运行API测试
```bash
go test -v ./api/v1/ai -run "TestWritingAssistantApi.*"
```

### 运行性能测试
```bash
go test -v ./service/ai -bench="Benchmark.*" -benchmem
```

### 使用测试脚本
```bash
# Windows
run_ai_writing_tests.bat
```

## 测试覆盖的场景

### ✅ AI正常响应
- Mock适配器返回成功响应
- 验证响应数据结构
- 验证Token使用统计

### ✅ AI服务超时
- Mock适配器配置超时延迟
- 验证超时错误处理
- 验证上下文取消

### ✅ AI服务错误
- Mock适配器配置失败状态
- 验证错误类型识别
- 验证错误信息传递

### ✅ 配额不足
- 配额检查（待实现）
- 配额扣减（待实现）
- 配额不足错误处理（待实现）

### ✅ 敏感词命中
- 政治敏感词检测
- 暴力敏感词检测
- 成人内容检测
- 高风险词识别

### ✅ 敏感词未命中
- 正常内容通过检测
- 验证IsSafe标志
- 验证TotalMatches为0

### ✅ 自定义敏感词
- 添加自定义词
- 移除自定义词
- 自定义词检测
- 用户隔离（不同用户的自定义词库独立）

### ✅ 缓存命中/未命中
- 结果缓存（待实现）
- 缓存命中验证（待实现）
- 缓存过期处理（待实现）

## 测试数据示例

### 内容总结请求示例
```json
{
  "content": "这是需要总结的完整文章内容...",
  "summaryType": "detailed",
  "maxLength": 1000,
  "includeQuotes": true
}
```

### 文本校对请求示例
```json
{
  "content": "这是需要校对的文本内容...",
  "checkTypes": ["grammar", "spelling", "punctuation"],
  "language": "zh-CN",
  "suggestions": true
}
```

### 敏感词检测请求示例
```json
{
  "content": "这是需要检测敏感词的内容...",
  "customWords": ["自定义词1", "自定义词2"],
  "category": "all"
}
```

## 测试覆盖率

根据Go测试覆盖率工具：
```bash
go test -cover ./service/ai/...
go test -coverprofile=coverage.out ./service/ai/...
go tool cover -html=coverage.out
```

## 待完善的功能

### TODO标记的测试
1. **适配器管理器集成** - 需要完善AdapterManager的Mock支持
2. **配额系统** - 需要实现配额检查和扣减逻辑
3. **结果缓存** - 需要实现缓存层并添加缓存测试
4. **数据库集成** - 需要Mock数据库层进行完整测试

### 建议的改进
1. 添加表驱动测试（Table-Driven Tests）
2. 添加模糊测试（Fuzzing Tests）
3. 添加竞争检测（Race Detection）
4. 添加集成测试（需要真实AI适配器）

## 测试最佳实践

1. **使用Mock隔离依赖** - 所有AI适配器调用都被Mock
2. **测试所有边界条件** - 空值、零值、最大值等
3. **验证错误处理** - 确保错误被正确传播
4. **性能基准测试** - 关键算法都有性能测试
5. **清晰的测试命名** - 测试名称清楚描述测试场景

## 贡献指南

添加新测试时：
1. 使用清晰的测试命名: `Test{Service}_{Method}_{Scenario}`
2. 添加必要的注释说明测试目的
3. 遵循现有测试结构
4. 更新本文档

## 联系方式

如有问题或建议，请通过项目Issue反馈。
