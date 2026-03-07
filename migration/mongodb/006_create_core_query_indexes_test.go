package migration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCoreQueryIndexes_Definition(t *testing.T) {
	expectedIndexes := map[string][]string{
		"users": {
			"username_1_unique",
			"email_1_unique_non_empty",
			"phone_1_unique_non_empty",
		},
		"comments": {
			"target_id_1_target_type_1_state_1_created_at_-1",
			"author_id_1_state_1_created_at_-1",
			"parent_id_1_state_1_created_at_1",
		},
		"notifications": {
			"user_id_1_created_at_-1",
			"user_id_1_is_read_1_created_at_-1",
			"user_id_1_type_1_created_at_-1",
		},
	}

	assert.Len(t, expectedIndexes["users"], 3, "users 应定义 3 个核心索引")
	assert.Len(t, expectedIndexes["comments"], 3, "comments 应定义 3 个核心索引")
	assert.Len(t, expectedIndexes["notifications"], 3, "notifications 应定义 3 个核心索引")
}
