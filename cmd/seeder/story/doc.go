// Package story 提供演示数据"星际觉醒"的所有故事内容
//
// 这是一个完整的科幻小说演示项目，包含：
// - 1个项目：星际觉醒
// - 4卷：火星觉醒、独立战争、第一次接触、联盟纪元
// - 24章：每卷6章，完整故事线
// - 12个角色：主角、配角、反派
// - 30+条关系边：角色之间的复杂关系网络
// - 8个道具、6个地点、3条时间线
//
// 使用方式：
//
//	projectID := createProject(adminID)
//	docIDMap := createDocuments(projectID)
//	createOutlines(projectID, docIDMap)
//	characterIDs := createCharacters(projectID)
//	createRelations(projectID, characterIDs)
//	createAssets(projectID, characterIDs)
package story
