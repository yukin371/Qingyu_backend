package dto

import (
	"Qingyu_backend/models/writer"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToProjectResponse converts a Project model to ProjectResponse DTO
func ToProjectResponse(p *writer.Project) ProjectResponse {
	return ProjectResponse{
		ID:        p.ID.Hex(),
		Title:     p.Title,
		Summary:   p.Summary,
		CoverURL:  p.CoverURL,
		Tags:      p.Tags,
		Status:    string(p.Status),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// ToProjectResponseList converts a slice of Project models to ProjectResponse DTOs
func ToProjectResponseList(projects []*writer.Project) []ProjectResponse {
	responses := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		responses[i] = ToProjectResponse(p)
	}
	return responses
}

// ToDocumentResponse converts a Document model to DocumentResponse DTO
func ToDocumentResponse(d *writer.Document) DocumentResponse {
	var parentID *string
	if !d.ParentID.IsZero() {
		s := d.ParentID.Hex()
		parentID = &s
	}

	return DocumentResponse{
		ID:           d.ID.Hex(),
		ProjectID:    d.ProjectID.Hex(),
		ParentID:     parentID,
		Title:        d.Title,
		Type:         DocumentType(d.Type),
		Level:        d.Level,
		Order:        d.Order,
		OrderKey:     d.OrderKey,
		Status:       DocumentStatus(d.Status),
		WordCount:    d.WordCount,
		CharacterIDs: convertObjectIDSlice(d.CharacterIDs),
		LocationIDs:  convertObjectIDSlice(d.LocationIDs),
		TimelineIDs:  convertObjectIDSlice(d.TimelineIDs),
		Tags:         d.Tags,
		Notes:        d.Notes,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

// ToDocumentResponseList converts a slice of Document models to DocumentResponse DTOs
func ToDocumentResponseList(documents []*writer.Document) []DocumentResponse {
	responses := make([]DocumentResponse, len(documents))
	for i, d := range documents {
		responses[i] = ToDocumentResponse(d)
	}
	return responses
}

// ToDocumentTreeItem converts a Document model to DocumentTreeItem DTO
func ToDocumentTreeItem(d *writer.Document) *DocumentTreeItem {
	if d == nil {
		return nil
	}

	var parentID *string
	if !d.ParentID.IsZero() {
		s := d.ParentID.Hex()
		parentID = &s
	}

	return &DocumentTreeItem{
		ID:        d.ID.Hex(),
		ParentID:  parentID,
		Title:     d.Title,
		Type:      DocumentType(d.Type),
		Level:     d.Level,
		OrderKey:  d.OrderKey,
		WordCount: d.WordCount,
		Children:  nil, // 需要调用方填充
	}
}

// ToDocumentTreeItemList converts a slice of Document models to DocumentTreeItem DTOs
func ToDocumentTreeItemList(documents []*writer.Document) []*DocumentTreeItem {
	items := make([]*DocumentTreeItem, len(documents))
	for i, d := range documents {
		items[i] = ToDocumentTreeItem(d)
	}
	return items
}

// Helper function to convert []primitive.ObjectID to []string
func convertObjectIDSlice(ids []primitive.ObjectID) []string {
	if ids == nil {
		return nil
	}
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = id.Hex()
	}
	return result
}
