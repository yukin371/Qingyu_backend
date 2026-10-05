package reader

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAnnotationIDMatchSupportsCanonicalAndLegacyValues(t *testing.T) {
	id := primitive.NewObjectID()
	match, err := annotationIDMatch("用户ID", id.Hex())

	require.NoError(t, err)
	require.Equal(t, bson.M{"$in": []interface{}{id, id.Hex()}}, match)
}

func TestAnnotationIDMatchRejectsMalformedValue(t *testing.T) {
	_, err := annotationIDMatch("书籍ID", "malformed")

	require.EqualError(t, err, "无效的书籍ID")
}

func TestAnnotationOwnerFilterRejectsMalformedValueWithoutResourceLeak(t *testing.T) {
	_, err := annotationOwnerFilter("malformed")

	require.EqualError(t, err, "标注不存在")
}
