"""
专业Agent模块

包含专门用于创作任务的Agent：
- OutlineAgent: 大纲生成Agent
- CharacterAgent: 角色设计Agent
- PlotAgent: 情节安排Agent
"""

from agents.specialized.outline_agent import OutlineAgent
from agents.specialized.character_agent import CharacterAgent
from agents.specialized.plot_agent import PlotAgent

__all__ = [
    "OutlineAgent",
    "CharacterAgent",
    "PlotAgent",
]

