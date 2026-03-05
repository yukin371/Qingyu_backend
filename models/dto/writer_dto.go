package dto

import "time"

// ===========================
// Writer DTO（符合分层架构规范）
// ===========================
//
// 本文件包含 Writer 模块的数据传输对象（DTO）
//
// 命名和标签规范：
// - DTO 结构体使用驼峰命名（PascalCase）
// - JSON 字段标签使用驼峰命名（camelCase）
// - 对应的 MongoDB 模型（位于 models/writer/）使用蛇形命名（snake_case）的 BSON 标签
//
// 用途：
// - 用于 Service 层和 API 层之间的数据传输
// - ID 和时间字段统一使用字符串类型
// - 避免直接暴露 MongoDB 模型到 API 层
