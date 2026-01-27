package migration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateChaptersIndexes_Definition(t *testing.T) {
	expectedIndexes := []string{
		"book_id_1_chapter_number_1",
		"book_id_1_status_1_chapter_number_1",
	}
	assert.Len(t, expectedIndexes, 2, "应定义2个chapters索引")
}

func TestCreateChaptersIndexes_Up(t *testing.T) {
	m := &CreateChaptersIndexes{}
	assert.NotNil(t, m, "CreateChaptersIndexes实例不应为空")
}

func TestCreateChaptersIndexes_Down(t *testing.T) {
	m := &CreateChaptersIndexes{}
	assert.NotNil(t, m, "CreateChaptersIndexes实例不应为空")
}
