package reader_test

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 测试辅助函数：生成固定的ObjectID用于测试
func getTestObjectIDs() (userID, bookID, chapterID primitive.ObjectID) {
	userID, _ = primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	bookID, _ = primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	chapterID, _ = primitive.ObjectIDFromHex("507f1f77bcf86cd799439013")
	return
}
