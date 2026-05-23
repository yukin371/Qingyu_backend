package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoEngineBuildSearchQuery_WildcardQueryOnlyUsesFilters(t *testing.T) {
	engine := &MongoEngine{}

	filter := engine.buildSearchQuery("*", &SearchOptions{
		Filter: bson.M{
			"author": "测试作者",
		},
	})

	assert.NotContains(t, filter, "$or")
	assert.Equal(t, "测试作者", filter["author"])
}

func TestMongoEngineBuildSearchQuery_EscapesRegexMetaCharacters(t *testing.T) {
	engine := &MongoEngine{}

	filter := engine.buildSearchQuery("a+b", nil)
	searchConditions, ok := filter["$or"].([]bson.M)
	if !ok {
		t.Fatalf("expected $or search conditions, got: %#v", filter)
	}

	assert.Equal(t, `a\+b`, searchConditions[0]["title"].(bson.M)["$regex"])
}
