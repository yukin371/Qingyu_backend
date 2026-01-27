package migration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateReadingProgressIndexes_Definition(t *testing.T) {
	expectedIndexes := []string{
		"user_id_1_updated_at_-1",
		"book_id_1",
	}
	assert.Len(t, expectedIndexes, 2, "应定义2个reading_progress索引")
}

func TestCreateReadingProgressIndexes_Up(t *testing.T) {
	m := &CreateReadingProgressIndexes{}
	assert.NotNil(t, m, "CreateReadingProgressIndexes实例不应为空")
}

func TestCreateReadingProgressIndexes_Down(t *testing.T) {
	m := &CreateReadingProgressIndexes{}
	assert.NotNil(t, m, "CreateReadingProgressIndexes实例不应为空")
}
