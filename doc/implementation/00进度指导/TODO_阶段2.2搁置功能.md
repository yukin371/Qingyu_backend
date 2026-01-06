# 阶段 2.2 搁置功能 TODO

**创建日期**: 2025-10-28  
**优先级**: P1-P2（性能优化阶段）  
**预计时间**: 1.5小时

---

## 📋 待完成功能清单

### 1. Reranker 重排序器 ⏳

**优先级**: P1  
**预计时间**: 45分钟  
**依赖**: `sentence-transformers` 或 `FlagEmbedding`

**实现任务**:
- [ ] 创建 `python_ai_service/src/rag/reranker.py`
- [ ] 实现 `CrossEncoderReranker` 类
  - [ ] 模型加载（BAAI/bge-reranker-large）
  - [ ] `rerank()` 方法 - 重排序文档
  - [ ] `get_scores()` 方法 - 获取相关性分数
  - [ ] 批量处理支持
  - [ ] GPU加速支持
- [ ] 集成到 RAGPipeline
  - [ ] 更新 `retrieve_and_rerank()` 方法
  - [ ] 添加配置开关
- [ ] 编写单元测试
  - [ ] 测试重排序效果
  - [ ] 测试分数范围
  - [ ] 性能测试

**预期效果**:
- 检索准确率提升 15-20%
- Top-3 准确率提升 25%+

**代码模板**:
```python
from FlagEmbedding import FlagReranker

class CrossEncoderReranker:
    def __init__(self, model_name: str = "BAAI/bge-reranker-large"):
        self.model = FlagReranker(model_name, use_fp16=True)
    
    async def rerank(
        self, 
        query: str, 
        documents: List[str], 
        top_k: int = 5
    ) -> List[Tuple[int, float]]:
        # 构建query-doc对
        pairs = [[query, doc] for doc in documents]
        # 获取分数
        scores = self.model.compute_score(pairs)
        # 排序
        sorted_indices = sorted(
            enumerate(scores), 
            key=lambda x: x[1], 
            reverse=True
        )
        return sorted_indices[:top_k]
```

---

### 2. 混合检索系统 ⏳

**优先级**: P1  
**预计时间**: 40分钟  
**依赖**: `rank-bm25`

**实现任务**:

#### 2.1 BM25检索器
- [ ] 创建 `python_ai_service/src/rag/bm25_retriever.py`
- [ ] 实现 `BM25Retriever` 类
  - [ ] `index_documents()` - 索引文档
  - [ ] `search()` - BM25检索
  - [ ] 中文分词支持（jieba）
  - [ ] 持久化索引（pickle）

#### 2.2 RRF融合算法
- [ ] 创建 `python_ai_service/src/rag/fusion.py`
- [ ] 实现 `ReciprocalRankFusion` 类
  - [ ] `fuse()` - 融合多个排序结果
  - [ ] 支持自定义k值
  - [ ] 归一化分数

#### 2.3 混合检索器
- [ ] 创建 `python_ai_service/src/rag/hybrid_retriever.py`
- [ ] 实现 `HybridRetriever` 类
  - [ ] 向量检索 + BM25检索
  - [ ] RRF融合
  - [ ] 权重融合（可选）
  - [ ] 配置化权重

#### 2.4 集成测试
- [ ] 对比测试：向量 vs BM25 vs 混合
- [ ] 性能测试
- [ ] 中文查询测试

**预期效果**:
- 召回率提升 10-15%
- 关键词查询准确率提升 20%+

**代码模板**:
```python
from rank_bm25 import BM25Okapi
import jieba

class BM25Retriever:
    def __init__(self):
        self.corpus = []
        self.bm25 = None
        
    def index_documents(self, documents: List[str]):
        # 中文分词
        tokenized_corpus = [list(jieba.cut(doc)) for doc in documents]
        self.bm25 = BM25Okapi(tokenized_corpus)
    
    def search(self, query: str, top_k: int = 10):
        tokenized_query = list(jieba.cut(query))
        scores = self.bm25.get_scores(tokenized_query)
        top_indices = np.argsort(scores)[::-1][:top_k]
        return [(idx, scores[idx]) for idx in top_indices]
```

---

### 3. Citation Manager 引用管理 ⏳

**优先级**: P2  
**预计时间**: 25分钟

**实现任务**:
- [ ] 创建 `python_ai_service/src/rag/citation.py`
- [ ] 实现 `CitationManager` 类
  - [ ] `extract_citations()` - 提取引用标记 [1], [2]
  - [ ] `format_citations()` - 格式化引用列表
  - [ ] `validate_citations()` - 验证引用有效性
  - [ ] `add_citation_metadata()` - 添加引用元数据
- [ ] 集成到 ContextBuilder
- [ ] 编写测试用例

**代码模板**:
```python
import re

class CitationManager:
    def extract_citations(self, text: str) -> List[int]:
        """提取 [1], [2] 等引用标记"""
        pattern = r'\[(\d+)\]'
        citations = re.findall(pattern, text)
        return [int(c) for c in citations]
    
    def format_citations(self, results: List[RetrievalResult]) -> str:
        """格式化引用列表"""
        citations = []
        for i, result in enumerate(results, 1):
            citation = (
                f"[{i}] {result.source}\n"
                f"    {result.get_citation_text()}\n"
            )
            citations.append(citation)
        return "\n".join(citations)
```

---

### 4. 完整测试覆盖 ⏳

**优先级**: P1  
**预计时间**: 30分钟

**测试文件**:
- [ ] `tests/test_rag_pipeline.py`
  - [ ] 测试基础检索
  - [ ] 测试元数据过滤
  - [ ] 测试上下文构建
  - [ ] 测试Token限制
  - [ ] 测试空结果处理
  
- [ ] `tests/test_context_builder.py`
  - [ ] 测试Token计数
  - [ ] 测试智能截断
  - [ ] 测试自定义模板
  - [ ] 测试引用标注

- [ ] `tests/test_reranker.py`（Reranker实现后）
  - [ ] 测试重排序效果
  - [ ] 测试分数准确性
  
- [ ] `tests/test_hybrid_retrieval.py`（混合检索实现后）
  - [ ] 测试BM25检索
  - [ ] 测试RRF融合
  - [ ] 对比性能测试

---

## 📊 实施优先级

### 高优先级（P0-P1）- 性能提升明显
1. **Reranker** - 准确率提升20%+
2. **混合检索** - 召回率提升15%+
3. **完整测试** - 保证质量

### 中优先级（P2）- 用户体验
4. **Citation Manager** - 引用管理

---

## 🎯 实施建议

### 方案A：性能优化阶段统一实施
- 在Phase3完成后，专门安排性能优化阶段
- 集中实施所有优化功能
- 完整的A/B测试对比

### 方案B：按需实施
- 根据实际使用中的问题决定优先级
- 逐步实施，持续优化
- 更灵活但可能不够系统

**推荐**: 方案A，在Phase3完成后统一优化

---

## 📁 依赖包

需要额外安装：
```bash
# Reranker
pip install FlagEmbedding

# BM25
pip install rank-bm25

# 中文分词
pip install jieba

# Token计数（已有但可选）
pip install tiktoken
```

或添加到 `pyproject.toml`:
```toml
[tool.poetry.dependencies]
FlagEmbedding = {version = "^1.2.0", optional = true}
rank-bm25 = {version = "^0.2.2", optional = true}
jieba = {version = "^0.42.1", optional = true}

[tool.poetry.extras]
rag-advanced = ["FlagEmbedding", "rank-bm25", "jieba"]
```

---

## 📈 预期收益

| 功能 | 准确率提升 | 召回率提升 | 用户满意度 |
|------|-----------|-----------|-----------|
| Reranker | +20% | +5% | ⭐⭐⭐⭐⭐ |
| 混合检索 | +10% | +15% | ⭐⭐⭐⭐ |
| Citation | - | - | ⭐⭐⭐ |

**综合提升**: 准确率+30%，召回率+20%

---

## 🔗 相关文档

- [阶段2.2实施计划](./计划/阶段2.2_结构化RAG实现_实施计划.md)
- [阶段2.2进度报告](./阶段2.2进度报告_2025-10-28.md)
- [RAGPipeline源码](../../python_ai_service/src/rag/rag_pipeline.py)

---

**创建时间**: 2025-10-28  
**维护者**: 青羽后端架构团队  
**状态**: 待实施（性能优化阶段）

