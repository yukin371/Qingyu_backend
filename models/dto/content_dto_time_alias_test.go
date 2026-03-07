package dto

import (
	"encoding/json"
	"testing"
	"time"
)

func TestReadingProgressResponse_MarshalIncludesUpdatedAtAlias(t *testing.T) {
	resp := ReadingProgressResponse{
		UserID:      "user-1",
		BookID:      "book-1",
		ChapterID:   "chapter-1",
		Progress:    0.5,
		ReadingTime: 120,
		UpdateTime:  1710000000,
		UpdatedAt:   1710000000,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	payload := map[string]any{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if payload["updateTime"] != float64(1710000000) {
		t.Fatalf("expected updateTime alias, got %#v", payload["updateTime"])
	}
	if payload["updatedAt"] != float64(1710000000) {
		t.Fatalf("expected updatedAt alias, got %#v", payload["updatedAt"])
	}
}

func TestChapterPublishStatus_MarshalIncludesPublishedAtAliases(t *testing.T) {
	now := time.Unix(1710000000, 0).UTC()
	resp := ChapterPublishStatus{
		ChapterID:     "chapter-1",
		IsPublished:   true,
		ChapterNumber: 1,
		PublishTime:   &now,
		PublishedAt:   &now,
		UpdateTime:    now,
		UpdatedAt:     now,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	payload := map[string]any{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if payload["publishTime"] != now.Format(time.RFC3339) {
		t.Fatalf("expected publishTime alias, got %#v", payload["publishTime"])
	}
	if payload["publishedAt"] != now.Format(time.RFC3339) {
		t.Fatalf("expected publishedAt alias, got %#v", payload["publishedAt"])
	}
	if payload["updateTime"] != now.Format(time.RFC3339) {
		t.Fatalf("expected updateTime alias, got %#v", payload["updateTime"])
	}
	if payload["updatedAt"] != now.Format(time.RFC3339) {
		t.Fatalf("expected updatedAt alias, got %#v", payload["updatedAt"])
	}
}
