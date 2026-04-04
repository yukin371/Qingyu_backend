package audit

import (
	"context"
	"testing"
)

func TestGetPendingReviews_WithNilAuditRecordRepo_ReturnsEmptyList(t *testing.T) {
	service := NewContentAuditService(nil, nil, nil, nil)

	records, err := service.GetPendingReviews(context.Background(), 10)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected empty records, got %d", len(records))
	}
}

func TestGetPendingReviews_WithNilReceiver_ReturnsEmptyList(t *testing.T) {
	var service *ContentAuditService

	records, err := service.GetPendingReviews(context.Background(), 10)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected empty records, got %d", len(records))
	}
}

func TestGetHighRiskAudits_WithNilAuditRecordRepo_ReturnsEmptyList(t *testing.T) {
	service := NewContentAuditService(nil, nil, nil, nil)

	records, err := service.GetHighRiskAudits(context.Background(), 3, 10)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected empty records, got %d", len(records))
	}
}

func TestGetHighRiskAudits_WithNilReceiver_ReturnsEmptyList(t *testing.T) {
	var service *ContentAuditService

	records, err := service.GetHighRiskAudits(context.Background(), 3, 10)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected empty records, got %d", len(records))
	}
}
