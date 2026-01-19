# Phase3 创作API文档

**版本**: v1.0  
**日期**: 2025-10-30

---

## 📋 概述

Phase3创作API提供了基于AI的智能创作工作流，包括大纲生成、角色设计、情节安排等功能。

## 🔗 基础信息

- **Base URL**: `/api/v1/ai/creative`
- **认证**: 需要JWT Token（除健康检查外）
- **Content-Type**: `application/json`

---

## 📡 API列表

### 1. 健康检查

**GET** `/health`

检查Phase3 AI服务的健康状态。

**请求**:
无需参数（公开接口）

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "healthy",
    "checks": {
      "outline_agent": "healthy",
      "character_agent": "healthy",
      "plot_agent": "healthy"
    }
  }
}
```

---

### 2. 生成大纲

**POST** `/outline`

根据任务描述生成完整的故事大纲。

**请求头**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "task": "创作一个修仙小说大纲，主角是天才少年，包含5章内容",
  "project_id": "project_123",
  "workspace_context": {
    "genre": "修仙",
    "style": "热血"
  },
  "correction_prompt": ""
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|-----|------|
| task | string | 是 | 创作任务描述 |
| project_id | string | 否 | 项目ID |
| workspace_context | object | 否 | 工作区上下文 |
| correction_prompt | string | 否 | 修正提示 |

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "outline": {
      "title": "天命之子",
      "genre": "修仙",
      "core_theme": "成长与复仇",
      "target_audience": "青年读者",
      "estimated_total_words": 500000,
      "chapters": [
        {
          "chapter_id": 1,
          "title": "第一章：天降奇才",
          "summary": "主角林轩在一个平凡的小镇中出生，却展现出惊人的修炼天赋...",
          "key_events": ["发现天赋", "拜师学艺", "初遇危机"],
          "characters_involved": ["林轩", "玄天老祖"],
          "conflict_type": "人物冲突",
          "emotional_tone": "激昂",
          "estimated_word_count": 10000,
          "chapter_goal": "引出主角，展现天赋",
          "cliffhanger": "神秘黑衣人出现"
        }
      ],
      "story_arc": {
        "setup": [1],
        "rising_action": [2, 3],
        "climax": [4],
        "falling_action": [],
        "resolution": [5]
      }
    },
    "execution_time": 9.8
  }
}
```

**耗时**: 8-12秒

---

### 3. 生成角色

**POST** `/characters`

根据大纲生成角色设定，包含角色关系网络。

**请求头**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "task": "根据大纲创建主要角色",
  "project_id": "project_123",
  "outline": {
    "title": "天命之子",
    "chapters": [...]
  },
  "workspace_context": {},
  "correction_prompt": ""
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|-----|------|
| task | string | 是 | 创作任务描述 |
| project_id | string | 否 | 项目ID |
| outline | object | 建议 | 大纲数据（提供后质量更好） |
| workspace_context | object | 否 | 工作区上下文 |
| correction_prompt | string | 否 | 修正提示 |

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "characters": {
      "characters": [
        {
          "character_id": "char_001",
          "name": "林轩",
          "role_type": "protagonist",
          "importance": "major",
          "age": 16,
          "gender": "男",
          "appearance": "英俊少年，眉宇间带有坚毅之色",
          "personality": {
            "traits": ["勇敢", "聪慧", "执着"],
            "strengths": ["天赋异禀", "心志坚定"],
            "weaknesses": ["过于自信", "缺乏经验"],
            "core_values": "保护亲人，追求强大",
            "fears": "失去所爱之人"
          },
          "background": {
            "summary": "出生于小镇，父母早逝...",
            "family": "孤儿",
            "education": "自学为主",
            "key_experiences": ["父母被杀", "发现天赋"]
          },
          "motivation": "为父母报仇，保护所爱之人",
          "relationships": [
            {
              "character": "玄天老祖",
              "relation_type": "mentor",
              "description": "师徒关系",
              "dynamics": "亦师亦父"
            }
          ],
          "development_arc": {
            "starting_point": "稚嫩少年",
            "turning_points": ["拜师", "突破瓶颈", "大仇得报"],
            "ending_point": "一代宗师",
            "growth_theme": "从复仇到守护"
          },
          "role_in_story": "故事主线推动者",
          "first_appearance": 1,
          "chapters_involved": [1, 2, 3, 4, 5]
        }
      ],
      "relationship_network": {
        "alliances": [["林轩", "玄天老祖"]],
        "conflicts": [["林轩", "血煞魔君"]],
        "mentorships": [
          {
            "mentor": "玄天老祖",
            "student": "林轩"
          }
        ]
      }
    },
    "execution_time": 12.3
  }
}
```

**耗时**: 10-15秒

---

### 4. 生成情节

**POST** `/plot`

根据大纲和角色生成情节，包含时间线事件、情节线索等。

**请求头**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "task": "根据大纲和角色设计情节",
  "project_id": "project_123",
  "outline": {...},
  "characters": {...},
  "workspace_context": {},
  "correction_prompt": ""
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|-----|------|
| task | string | 是 | 创作任务描述 |
| project_id | string | 否 | 项目ID |
| outline | object | 建议 | 大纲数据 |
| characters | object | 建议 | 角色数据 |
| workspace_context | object | 否 | 工作区上下文 |
| correction_prompt | string | 否 | 修正提示 |

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "plot": {
      "timeline_events": [
        {
          "event_id": "event_001",
          "timestamp": "第1天",
          "location": "小镇",
          "title": "天赋觉醒",
          "description": "林轩在危机中觉醒修炼天赋",
          "participants": ["林轩"],
          "event_type": "转折",
          "impact": {
            "on_plot": "开启主线剧情",
            "on_characters": {
              "林轩": "开启修炼之路"
            },
            "emotional_impact": "震撼、激动"
          },
          "causes": ["生死危机"],
          "consequences": ["踏上修仙路"],
          "chapter_id": 1
        }
      ],
      "plot_threads": [
        {
          "thread_id": "thread_main",
          "title": "复仇之路",
          "description": "主角为父母复仇的主线",
          "type": "main",
          "events": ["event_001", "event_005", "event_010"],
          "starting_event": "event_001",
          "climax_event": "event_010",
          "resolution_event": "event_015",
          "characters_involved": ["林轩", "血煞魔君"]
        }
      ],
      "conflicts": [
        {
          "conflict_id": "conflict_001",
          "type": "人物冲突",
          "parties": ["林轩", "血煞魔君"],
          "description": "主角与仇敌的对抗",
          "escalation_events": ["event_005", "event_008"],
          "resolution_event": "event_010"
        }
      ],
      "key_plot_points": {
        "inciting_incident": "event_001",
        "plot_point_1": "event_005",
        "midpoint": "event_008",
        "plot_point_2": "event_009",
        "climax": "event_010",
        "resolution": "event_015"
      }
    },
    "execution_time": 14.7
  }
}
```

**耗时**: 12-18秒

---

### 5. 执行完整工作流

**POST** `/workflow`

一次性执行 Outline → Characters → Plot 完整创作流程。

**请求头**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "task": "创作一个都市爱情小说的完整设定，包含3章内容",
  "project_id": "project_123",
  "max_reflections": 3,
  "enable_human_review": false,
  "workspace_context": {}
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|-----|------|-----|--------|------|
| task | string | 是 | - | 创作任务描述 |
| project_id | string | 否 | - | 项目ID |
| max_reflections | int | 否 | 3 | 最大反思次数 |
| enable_human_review | bool | 否 | false | 是否启用人工审核 |
| workspace_context | object | 否 | {} | 工作区上下文 |

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "execution_id": "exec_uuid_xxx",
    "review_passed": true,
    "reflection_count": 0,
    "outline": {...},
    "characters": {...},
    "plot": {...},
    "diagnostic_report": null,
    "reasoning": [
      "OutlineAgent: 生成了3章大纲",
      "CharacterAgent: 生成了2个角色",
      "PlotAgent: 生成了15个时间线事件"
    ],
    "execution_times": {
      "outline": 9.8,
      "character": 12.3,
      "plot": 14.7
    },
    "tokens_used": 9000
  }
}
```

**耗时**: 30-45秒

---

## 🔒 认证

除 `/health` 接口外，所有接口都需要JWT认证：

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## ⚠️ 错误码

| 错误码 | 说明 |
|-------|------|
| 200 | 成功 |
| 400 | 参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 500 | 服务器错误 |
| 503 | AI服务不可用 |

**错误响应示例**:
```json
{
  "code": 500,
  "message": "大纲生成失败",
  "data": {
    "error": "连接AI服务失败: connection refused"
  }
}
```

---

## 📊 性能指标

| 接口 | 平均耗时 | Token消耗 |
|-----|---------|----------|
| GenerateOutline | 8-12秒 | ~2000 |
| GenerateCharacters | 10-15秒 | ~3000 |
| GeneratePlot | 12-18秒 | ~4000 |
| ExecuteCreativeWorkflow | 30-45秒 | ~9000 |

---

## 💡 使用建议

### 1. 分步使用

推荐先单独调用各个接口，逐步完善：

```
1. POST /outline → 生成大纲
2. POST /characters (传入outline) → 生成角色
3. POST /plot (传入outline + characters) → 生成情节
```

### 2. 一键生成

如需快速原型，使用完整工作流：

```
POST /workflow → 一次性生成所有内容
```

### 3. 工作区上下文

利用`workspace_context`传递额外信息：

```json
{
  "workspace_context": {
    "genre": "修仙",
    "style": "热血",
    "target_words": "500000",
    "existing_characters": "林轩,玄天老祖"
  }
}
```

---

## 🔗 前端集成示例

### JavaScript/Axios

```javascript
import axios from 'axios';

// 生成大纲
const generateOutline = async (task, token) => {
  try {
    const response = await axios.post(
      '/api/v1/ai/creative/outline',
      {
        task: task,
        project_id: 'my_project'
      },
      {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      }
    );
    
    console.log('大纲:', response.data.data.outline);
    return response.data.data;
  } catch (error) {
    console.error('生成失败:', error.response.data);
  }
};

// 执行完整工作流
const executeWorkflow = async (task, token) => {
  try {
    const response = await axios.post(
      '/api/v1/ai/creative/workflow',
      {
        task: task,
        max_reflections: 3
      },
      {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      }
    );
    
    const data = response.data.data;
    console.log('执行ID:', data.execution_id);
    console.log('大纲:', data.outline);
    console.log('角色:', data.characters);
    console.log('情节:', data.plot);
    
    return data;
  } catch (error) {
    console.error('执行失败:', error.response.data);
  }
};
```

---

## 📚 相关文档

- [Go客户端使用](../../cmd/test_phase3_grpc/README.md)
- [Python gRPC服务](../../python_ai_service/GRPC_INTEGRATION_GUIDE.md)
- [完整集成报告](../implementation/00进度指导/Phase3_Go集成完成总结_2025-10-30.md)

---

**最后更新**: 2025-10-30  
**维护者**: 青羽后端架构团队

