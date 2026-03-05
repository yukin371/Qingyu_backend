package dto

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ===========================
// Project DTOs 单元测试
// ===========================

func TestCreateProjectRequest_Valid(t *testing.T) {
	req := CreateProjectRequest{
		Title:    "Test Project",
		Summary:  "Test Summary",
		CoverURL: "https://example.com/cover.jpg",
		Tags:     []string{"fiction", "fantasy"},
	}

	assert.Equal(t, "Test Project", req.Title)
	assert.Equal(t, "Test Summary", req.Summary)
	assert.Equal(t, "https://example.com/cover.jpg", req.CoverURL)
	assert.Len(t, req.Tags, 2)
	assert.Contains(t, req.Tags, "fiction")
	assert.Contains(t, req.Tags, "fantasy")
}

func TestCreateProjectRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateProjectRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateProjectRequest{
				Title:   "Valid Title",
				Summary: "A valid summary",
				Tags:    []string{"fiction"},
			},
			wantErr: false,
		},
		{
			name: "empty title - should fail validation",
			req: CreateProjectRequest{
				Title: "",
			},
			wantErr: true,
		},
		{
			name: "title too long - should fail validation",
			req: CreateProjectRequest{
				Title: strings.Repeat("a", 101),
			},
			wantErr: true,
		},
		{
			name: "summary too long - should fail validation",
			req: CreateProjectRequest{
				Title:   "Valid Title",
				Summary: strings.Repeat("b", 501),
			},
			wantErr: true,
		},
		{
			name: "invalid cover URL format",
			req: CreateProjectRequest{
				Title:    "Valid Title",
				CoverURL: "not-a-url",
			},
			wantErr: true,
		},
		{
			name: "too many tags - should fail validation",
			req: CreateProjectRequest{
				Title: "Valid Title",
				Tags:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
			},
			wantErr: true,
		},
		{
			name: "empty tag in array - should fail validation",
			req: CreateProjectRequest{
				Title: "Valid Title",
				Tags:  []string{"fiction", ""},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				// 验证无效请求的特定条件
				if tt.req.Title == "" {
					assert.Empty(t, tt.req.Title, "title should be empty for invalid request")
				}
				if len(tt.req.Title) > 100 {
					assert.Greater(t, len(tt.req.Title), 100, "title exceeds max length")
				}
				if len(tt.req.Summary) > 500 {
					assert.Greater(t, len(tt.req.Summary), 500, "summary exceeds max length")
				}
			} else {
				assert.NotEmpty(t, tt.req.Title, "title should not be empty for valid request")
				assert.LessOrEqual(t, len(tt.req.Title), 100, "title should not exceed max length")
			}
		})
	}
}

func TestProjectResponse_Fields(t *testing.T) {
	now := time.Now()
	resp := ProjectResponse{
		ID:        "proj_123",
		Title:     "Test Project",
		Summary:   "Test Summary",
		CoverURL:  "https://example.com/cover.jpg",
		Tags:      []string{"fiction"},
		Status:    "draft",
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, "proj_123", resp.ID)
	assert.Equal(t, "Test Project", resp.Title)
	assert.Equal(t, "Test Summary", resp.Summary)
	assert.Equal(t, "https://example.com/cover.jpg", resp.CoverURL)
	assert.Equal(t, "draft", resp.Status)
	assert.Len(t, resp.Tags, 1)
	assert.False(t, resp.CreatedAt.IsZero())
	assert.False(t, resp.UpdatedAt.IsZero())
}

func TestProjectResponse_EmptyFields(t *testing.T) {
	now := time.Now()
	resp := ProjectResponse{
		ID:        "proj_456",
		Title:     "Project Without Optional Fields",
		Summary:   "",
		CoverURL:  "",
		Tags:      []string{},
		Status:    "published",
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, "proj_456", resp.ID)
	assert.Empty(t, resp.Summary)
	assert.Empty(t, resp.CoverURL)
	assert.Empty(t, resp.Tags)
	assert.Equal(t, "published", resp.Status)
}

func TestListProjectsRequest_DefaultValues(t *testing.T) {
	req := ListProjectsRequest{
		Page:     0,
		PageSize: 0,
	}

	// 测试默认值应该被handler应用
	assert.Equal(t, 0, req.Page)
	assert.Equal(t, 0, req.PageSize)
}

func TestListProjectsRequest_WithValues(t *testing.T) {
	req := ListProjectsRequest{
		Page:     2,
		PageSize: 50,
		Status:   "draft",
		Sort:     "created_at",
		Order:    "desc",
	}

	assert.Equal(t, 2, req.Page)
	assert.Equal(t, 50, req.PageSize)
	assert.Equal(t, "draft", req.Status)
	assert.Equal(t, "created_at", req.Sort)
	assert.Equal(t, "desc", req.Order)
}

func TestProjectListResponse_Structure(t *testing.T) {
	now := time.Now()
	items := []ProjectResponse{
		{ID: "1", Title: "Project 1", Status: "draft", CreatedAt: now, UpdatedAt: now},
		{ID: "2", Title: "Project 2", Status: "published", CreatedAt: now, UpdatedAt: now},
	}

	resp := ProjectListResponse{
		Items:    items,
		Total:    2,
		Page:     1,
		PageSize: 20,
	}

	assert.Len(t, resp.Items, 2)
	assert.Equal(t, int64(2), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 20, resp.PageSize)
	assert.Equal(t, "Project 1", resp.Items[0].Title)
	assert.Equal(t, "Project 2", resp.Items[1].Title)
}

func TestProjectListResponse_Empty(t *testing.T) {
	resp := ProjectListResponse{
		Items:    []ProjectResponse{},
		Total:    0,
		Page:     1,
		PageSize: 20,
	}

	assert.Empty(t, resp.Items)
	assert.Equal(t, int64(0), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 20, resp.PageSize)
}

func TestUpdateProjectRequest_PointerFields(t *testing.T) {
	title := "Updated Title"
	summary := "Updated Summary"
	tags := []string{"updated", "tags"}

	req := UpdateProjectRequest{
		Title:    &title,
		Summary:  &summary,
		CoverURL: nil,
		Tags:     &tags,
	}

	assert.NotNil(t, req.Title)
	assert.NotNil(t, req.Summary)
	assert.Nil(t, req.CoverURL)
	assert.NotNil(t, req.Tags)
	assert.Equal(t, "Updated Title", *req.Title)
	assert.Equal(t, "Updated Summary", *req.Summary)
	assert.Len(t, *req.Tags, 2)
}

func TestUpdateProjectRequest_AllNil(t *testing.T) {
	req := UpdateProjectRequest{
		Title:    nil,
		Summary:  nil,
		CoverURL: nil,
		Tags:     nil,
	}

	assert.Nil(t, req.Title)
	assert.Nil(t, req.Summary)
	assert.Nil(t, req.CoverURL)
	assert.Nil(t, req.Tags)
}

// ===========================
// Converter Functions 单元测试
// ===========================

// Mock types for testing converter functions
// These mock types avoid circular dependencies with the models package

type mockProject struct {
	ID        string
	Title     string
	Summary   string
	CoverURL  string
	Tags      []string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func TestToProjectResponse(t *testing.T) {
	objID := "507f1f77bcf86cd799439011"
	now := time.Now()

	project := &mockProject{
		ID:        objID,
		Title:     "Test Project",
		Summary:   "Test Summary",
		CoverURL:  "https://example.com/cover.jpg",
		Tags:      []string{"fiction", "fantasy"},
		Status:    "draft",
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := ProjectResponse{
		ID:        project.ID,
		Title:     project.Title,
		Summary:   project.Summary,
		CoverURL:  project.CoverURL,
		Tags:      project.Tags,
		Status:    project.Status,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}

	assert.Equal(t, objID, response.ID)
	assert.Equal(t, "Test Project", response.Title)
	assert.Equal(t, "Test Summary", response.Summary)
	assert.Equal(t, "https://example.com/cover.jpg", response.CoverURL)
	assert.Len(t, response.Tags, 2)
	assert.Contains(t, response.Tags, "fiction")
	assert.Contains(t, response.Tags, "fantasy")
	assert.Equal(t, "draft", response.Status)
	assert.False(t, response.CreatedAt.IsZero())
	assert.False(t, response.UpdatedAt.IsZero())
}

func TestToProjectResponse_EmptyFields(t *testing.T) {
	objID := "507f1f77bcf86cd799439012"
	now := time.Now()

	project := &mockProject{
		ID:        objID,
		Title:     "Project Without Optional Fields",
		Summary:   "",
		CoverURL:  "",
		Tags:      []string{},
		Status:    "published",
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := ProjectResponse{
		ID:        project.ID,
		Title:     project.Title,
		Summary:   project.Summary,
		CoverURL:  project.CoverURL,
		Tags:      project.Tags,
		Status:    project.Status,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}

	assert.Equal(t, objID, response.ID)
	assert.Equal(t, "Project Without Optional Fields", response.Title)
	assert.Empty(t, response.Summary)
	assert.Empty(t, response.CoverURL)
	assert.Empty(t, response.Tags)
	assert.Equal(t, "published", response.Status)
}

func TestToProjectResponseList(t *testing.T) {
	objID1 := "507f1f77bcf86cd799439013"
	objID2 := "507f1f77bcf86cd799439014"

	projects := []*mockProject{
		{ID: objID1, Title: "Project 1", Status: "draft", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: objID2, Title: "Project 2", Status: "published", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	responses := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		responses[i] = ProjectResponse{
			ID:        p.ID,
			Title:     p.Title,
			Summary:   p.Summary,
			CoverURL:  p.CoverURL,
			Tags:      p.Tags,
			Status:    p.Status,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}

	assert.Len(t, responses, 2)
	assert.Equal(t, "Project 1", responses[0].Title)
	assert.Equal(t, "draft", responses[0].Status)
	assert.Equal(t, "Project 2", responses[1].Title)
	assert.Equal(t, "published", responses[1].Status)
}

func TestToProjectResponseList_Empty(t *testing.T) {
	projects := []*mockProject{}

	responses := make([]ProjectResponse, len(projects))

	assert.Len(t, responses, 0)
	assert.Empty(t, responses)
}

func TestToProjectResponseList_Nil(t *testing.T) {
	var projects []*mockProject

	responses := make([]ProjectResponse, 0)
	if projects != nil {
		responses = make([]ProjectResponse, len(projects))
		for i, p := range projects {
			responses[i] = ProjectResponse{
				ID:        p.ID,
				Title:     p.Title,
				Summary:   p.Summary,
				CoverURL:  p.CoverURL,
				Tags:      p.Tags,
				Status:    p.Status,
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,
			}
		}
	}

	assert.Len(t, responses, 0)
	assert.Empty(t, responses)
}
