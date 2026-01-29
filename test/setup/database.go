package setup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupTestDB 创建测试数据库连接
func SetupTestDB(t *testing.T) *mongo.Database {
	ctx := context.Background()

	// 连接到MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err, "测试数据库连接失败")

	// 测试数据库名称
	db := client.Database("qingyu_test")

	return db
}

// TeardownTestDB 清理测试数据库
func TeardownTestDB(t *testing.T, db *mongo.Database) {
	ctx := context.Background()

	// 删除测试数据库中的所有集合
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	assert.NoError(t, err, "列出集合失败")

	for _, collectionName := range collections {
		err := db.Collection(collectionName).Drop(ctx)
		assert.NoError(t, err, "删除集合失败: "+collectionName)
	}

	// 断开数据库连接
	client := db.Client()
	err = client.Disconnect(ctx)
	assert.NoError(t, err, "断开数据库连接失败")
}

// SeedTestData 插入测试数据
func SeedTestData(t *testing.T, db *mongo.Database, collection string, data []interface{}) {
	ctx := context.Background()
	coll := db.Collection(collection)

	_, err := coll.InsertMany(ctx, data)
	assert.NoError(t, err, "插入测试数据失败")
}
