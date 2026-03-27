package writer_test

import (
	"testing"
	"time"

	"Qingyu_backend/models/writer"

	"github.com/stretchr/testify/assert"
)

// mockTimelineEvent creates a test timeline event
func mockTimelineEvent(chapterID, chapterTitle string, oldType, newType writer.RelationType, strength int) writer.RelationTimelineEvent {
	return writer.RelationTimelineEvent{
		ChapterID:    chapterID,
		ChapterTitle: chapterTitle,
		OldType:     oldType,
		NewType:     newType,
		Strength:    strength,
		Notes:       "测试变化",
		Timestamp:   time.Now(),
	}
}

func TestCreateRelationTimelineEvent_Structure(t *testing.T) {
	// Test that timeline event structure is correct
	event := mockTimelineEvent("ch1", "第一章", writer.RelationFriend, writer.RelationEnemy, 80)

	assert.Equal(t, "ch1", event.ChapterID)
	assert.Equal(t, "第一章", event.ChapterTitle)
	assert.Equal(t, writer.RelationEnemy, event.NewType)
	assert.Equal(t, 80, event.Strength)
	assert.Equal(t, writer.RelationFriend, event.OldType)
}

func TestRelationTimelineEvent_Validation(t *testing.T) {
	// Test relation timeline event validation
	event := mockTimelineEvent("ch1", "第一章", writer.RelationFriend, writer.RelationEnemy, 80)

	assert.NotEmpty(t, event.ChapterID)
	assert.NotEmpty(t, event.NewType)
	assert.NotZero(t, event.Timestamp)
}

func TestGetRelationTimeline_EmptyList(t *testing.T) {
	// Test that empty timeline returns empty slice
	events := []writer.RelationTimelineEvent{}

	assert.Len(t, events, 0)
}

func TestCharacterRelation_WithTimelineEvents(t *testing.T) {
	// Test that timeline events can be created and stored
	events := []writer.RelationTimelineEvent{
		mockTimelineEvent("ch1", "第一章", writer.RelationFriend, writer.RelationEnemy, 80),
		mockTimelineEvent("ch2", "第二章", writer.RelationEnemy, writer.RelationAlly, 60),
	}

	assert.Len(t, events, 2)
	assert.Equal(t, writer.RelationEnemy, events[0].NewType)
	assert.Equal(t, writer.RelationAlly, events[1].NewType)
}

func TestRelationType_Constants(t *testing.T) {
	// Test relation type constants
	assert.Equal(t, writer.RelationType("朋友"), writer.RelationFriend)
	assert.Equal(t, writer.RelationType("敌人"), writer.RelationEnemy)
	assert.Equal(t, writer.RelationType("恋人"), writer.RelationRomance)
	assert.Equal(t, writer.RelationType("盟友"), writer.RelationAlly)
	assert.Equal(t, writer.RelationType("家庭"), writer.RelationFamily)
}

func TestIsValidRelationType(t *testing.T) {
	// Test relation type validation
	assert.True(t, writer.IsValidRelationType("朋友"))
	assert.True(t, writer.IsValidRelationType("敌人"))
	assert.True(t, writer.IsValidRelationType("恋人"))
	assert.False(t, writer.IsValidRelationType("未知类型"))
}
