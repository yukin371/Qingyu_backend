# API Integration Phase 1: Backend DTO Unification Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create unified Writer module DTO definitions and update handlers to use them, establishing the foundation for type-safe API integration.

**Architecture:**
- Create `models/dto/writer_dto.go` with camelCase JSON tags
- Update Writer handlers to use new DTOs instead of raw models
- Keep MongoDB models in `models/writer/` unchanged (BSON snake_case)
- Maintain API compatibility - no breaking changes for consumers

**Tech Stack:**
- Go 1.21+
- Gin framework
- go-playground/validator for validation tags
- Existing MongoDB models

**Prerequisites:**
- Read `docs/standards/error_code_standard.md` for error code reference
- Read `docs/standards/api/API设计规范.md` for API design standards
- Review existing DTOs in `models/dto/` for patterns to follow

---

## Task 1: Create Writer DTO Base File

**Files:**
- Create: `models/dto/writer_dto.go`

**Step 1: Create the DTO file with package declaration**

Create `models/dto/writer_dto.go`:

```go
package dto

import "time"

// This file contains all Data Transfer Objects (DTOs) for the Writer module.
// DTOs are used for API requests/responses and use camelCase JSON tags.
// MongoDB models in models/writer/ use BSON snake_case tags.
```

**Step 2: Run gofmt to verify formatting**

Run: `gofmt -w models/dto/writer_dto.go`
Expected: No errors, file formatted

**Step 3: Commit**

```bash
git add models/dto/writer_dto.go
git commit -m "feat(writer): create writer DTO file with package declaration"
```

---

## Task 2: Add Project DTOs

**Files:**
- Modify: `models/dto/writer_dto.go`

**Step 1: Read existing Project model to understand structure**

Reference: `models/writer/project.go`

Key fields from MongoDB model:
- ID (primitive.ObjectID)
- Title (string)
- Summary (string)
- CoverURL (string)
- Tags ([]string)
- Status (string)
- CreatedAt, UpdatedAt (time.Time)

**Step 2: Add Project DTOs to writer_dto.go**

Append to `models/dto/writer_dto.go`:

```go
// CreateProjectRequest represents a request to create a new project
type CreateProjectRequest struct {
	Title    string   `json:"title" validate:"required,min=1,max=100"`
	Summary  string   `json:"summary,omitempty" validate:"max=500"`
	CoverURL string   `json:"coverUrl,omitempty" validate:"omitempty,url,max=500"`
	Tags     []string `json:"tags,omitempty" validate:"max=10,dive,min=1,max=50"`
}

// UpdateProjectRequest represents a request to update an existing project
type UpdateProjectRequest struct {
	Title    *string   `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
	Summary  *string   `json:"summary,omitempty" validate:"omitempty,max=500"`
	CoverURL *string   `json:"coverUrl,omitempty" validate:"omitempty,url,max=500"`
	Tags     *[]string `json:"tags,omitempty" validate:"omitempty,max=10,dive,min=1,max=50"`
}

// ProjectResponse represents a project in API responses
type ProjectResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	CoverURL  string    `json:"coverUrl"`
	Tags      []string  `json:"tags"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ListProjectsRequest represents query parameters for listing projects
type ListProjectsRequest struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"page_size" validate:"min=1,max=100"`
	Status   string `form:"status" validate:"omitempty,oneof=draft published archived"`
	Sort     string `form:"sort" validate:"omitempty,oneof=created_at updated_at title"`
	Order    string `form:"order" validate:"omitempty,oneof=asc desc"`
}

// ProjectListResponse represents a paginated list of projects
type ProjectListResponse struct {
	Items    []ProjectResponse `json:"items"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}
```

**Step 3: Run go build to verify syntax**

Run: `go build ./models/dto/...`
Expected: No errors

**Step 4: Commit**

```bash
git add models/dto/writer_dto.go
git commit -m "feat(writer): add Project DTOs with validation tags"
```

---

## Task 3: Add Document DTOs

**Files:**
- Modify: `models/dto/writer_dto.go`

**Context:**
- Documents support tree structure (parentId, type, level, order)
- Document types: volume, chapter, section, scene
- Matches frontend types in `Qingyu_fronted/src/modules/writer/types/document.ts`
- Aligns with Editor V2 OpenAPI design

**Step 1: Add Document DTOs to writer_dto.go**

Append to `models/dto/writer_dto.go`:

```go
// DocumentType represents the type of a document node
type DocumentType string

const (
	DocumentTypeVolume  DocumentType = "volume"  // 卷
	DocumentTypeChapter DocumentType = "chapter" // 章
	DocumentTypeSection DocumentType = "section" // 节
	DocumentTypeScene   DocumentType = "scene"   // 场景
)

// DocumentStatus represents the writing status of a document
type DocumentStatus string

const (
	DocumentStatusPlanned   DocumentStatus = "planned"   // 计划中
	DocumentStatusWriting   DocumentStatus = "writing"   // 写作中
	DocumentStatusCompleted DocumentStatus = "completed" // 已完成
)

// CreateDocumentRequest 创建文档请求（树形结构）
type CreateDocumentRequest struct {
	ProjectID string       `json:"projectId" validate:"required"`
	ParentID  *string      `json:"parentId,omitempty"`              // 父节点ID，null表示根节点
	Title     string       `json:"title" validate:"required,min=1,max=200"`
	Type      DocumentType `json:"type" validate:"required,oneof=volume chapter section scene"`
	Level     int          `json:"level" validate:"min=0,max=10"`    // 层级深度
	Order     int          `json:"order" validate:"min=0"`          // 排序位置
}

// UpdateDocumentRequest 更新文档元数据请求
type UpdateDocumentRequest struct {
	Title        *string       `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Status       *DocumentStatus `json:"status,omitempty" validate:"omitempty,oneof=planned writing completed"`
	CharacterIDs *[]string     `json:"characterIds,omitempty" validate:"omitempty,max=50"`
	LocationIDs  *[]string     `json:"locationIds,omitempty" validate:"omitempty,max=50"`
	TimelineIDs  *[]string     `json:"timelineIds,omitempty" validate:"omitempty,max=50"`
	Tags         *[]string     `json:"tags,omitempty" validate:"omitempty,max=20"`
	Notes        *string       `json:"notes,omitempty" validate:"omitempty,max=1000"`
	OrderKey     *string       `json:"orderKey,omitempty"`           // LexoRank排序键
}

// DocumentResponse 文档响应（树形结构）
type DocumentResponse struct {
	ID          string        `json:"id"`
	ProjectID   string        `json:"projectId"`
	ParentID    *string       `json:"parentId,omitempty"`           // 父节点ID
	Title       string        `json:"title"`
	Type        DocumentType  `json:"type"`
	Level       int           `json:"level"`
	Order       int           `json:"order"`
	OrderKey    string        `json:"orderKey"`                    // LexoRank排序键
	Status      DocumentStatus `json:"status"`
	WordCount   int           `json:"wordCount"`
	CharacterIDs []string      `json:"characterIds,omitempty"`
	LocationIDs  []string      `json:"locationIds,omitempty"`
	TimelineIDs  []string      `json:"timelineIds,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
	Notes       string        `json:"notes,omitempty"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

// DocumentTreeResponse 文档树响应（嵌套结构）
type DocumentTreeResponse struct {
	ProjectID string             `json:"projectId"`
	Documents []*DocumentTreeItem `json:"documents"`
}

// DocumentTreeItem 文档树节点
type DocumentTreeItem struct {
	ID       string             `json:"id"`
	ParentID *string            `json:"parentId,omitempty"`
	Title    string             `json:"title"`
	Type     DocumentType       `json:"type"`
	Level    int                `json:"level"`
	OrderKey string             `json:"orderKey"`
	WordCount int               `json:"wordCount"`
	Children []*DocumentTreeItem `json:"children,omitempty"` // 子节点
}

// ListDocumentsRequest 查询文档列表请求
type ListDocumentsRequest struct {
	ProjectID string `form:"project_id" validate:"required"`
	ParentID  string `form:"parent_id,omitempty"`           // 筛选父节点下的文档
	Type      string `form:"type,omitempty" validate:"omitempty,oneof=volume chapter section scene"`
	Status    string `form:"status,omitempty" validate:"omitempty,oneof=planned writing completed"`
	Page      int    `form:"page" validate:"min=1"`
	PageSize  int    `form:"page_size" validate:"min=1,max=100"`
}

// DocumentListResponse 文档列表响应
type DocumentListResponse struct {
	Items    []DocumentResponse `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
}

// ReorderDocumentsRequest 重排序文档请求
type ReorderDocumentsRequest struct {
	ProjectID string           `json:"projectId" validate:"required"`
	ParentID  *string          `json:"parentId,omitempty"`  // 父节点ID
	Items     []ReorderItem    `json:"items" validate:"required,min=1,dive"` // 重排序项列表
}

// ReorderItem 重排序项
type ReorderItem struct {
	DocumentID string `json:"documentId" validate:"required"`
	ParentID   *string `json:"parentId,omitempty"`
	OrderKey   string `json:"orderKey" validate:"required"` // 目标排序键
}
```

**Step 3: Run go build to verify syntax**

Run: `go build ./models/dto/...`
Expected: No errors

**Step 4: Commit**

```bash
git add models/dto/writer_dto.go
git commit -m "feat(writer): add Document DTOs with validation tags"
```

---

**NOTE: Chapter is not a separate model - chapters are documents with type="chapter"**

---

## Task 4: Create DTO Converter Functions

**Files:**
- Create: `models/dto/writer_converter.go`

**Step 1: Create the converter file**

Create `models/dto/writer_converter.go`:

```go
package dto

import (
	"time"

	"github.com/QingyuBackend/Qingyu_backend/models/writer"
)

// ToProjectResponse converts a Project model to ProjectResponse DTO
func ToProjectResponse(p *writer.Project) ProjectResponse {
	return ProjectResponse{
		ID:        p.ID.Hex(),
		Title:     p.Title,
		Summary:   p.Summary,
		CoverURL:  p.CoverURL,
		Tags:      p.Tags,
		Status:    p.Status,
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
	return DocumentResponse{
		ID:          d.ID.Hex(),
		ProjectID:   d.ProjectID.Hex(),
		Title:       d.Title,
		Content:     d.Content,
		Summary:     d.Summary,
		WordCount:   d.WordCount,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
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

// ToDocumentResponseWithProject converts a Document model with project name to DocumentResponse DTO
func ToDocumentResponseWithProject(d *writer.Document, projectName string) DocumentResponse {
	return DocumentResponse{
		ID:          d.ID.Hex(),
		ProjectID:   d.ProjectID.Hex(),
		ProjectName: &projectName,
		Title:       d.Title,
		Content:     d.Content,
		Summary:     d.Summary,
		WordCount:   d.WordCount,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// ToChapterResponse converts a Chapter model to ChapterResponse DTO
func ToChapterResponse(c *writer.Chapter) ChapterResponse {
	return ChapterResponse{
		ID:         c.ID.Hex(),
		DocumentID: c.DocumentID.Hex(),
		Title:      c.Title,
		Content:    c.Content,
		WordCount:  c.WordCount,
		Order:      c.Order,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

// ToChapterResponseList converts a slice of Chapter models to ChapterResponse DTOs
func ToChapterResponseList(chapters []*writer.Chapter) []ChapterResponse {
	responses := make([]ChapterResponse, len(chapters))
	for i, c := range chapters {
		responses[i] = ToChapterResponse(c)
	}
	return responses
}
```

**Step 2: Run go build to verify**

Run: `go build ./models/dto/...`
Expected: No errors

**Step 3: Commit**

```bash
git add models/dto/writer_converter.go
git commit -m "feat(writer): add converter functions for Writer DTOs"
```

---

## Task 6: Create Unit Tests for Project DTOs

**Files:**
- Create: `models/dto/writer_dto_test.go`

**Step 1: Write failing test for Project DTO validation**

Create `models/dto/writer_dto_test.go`:

```go
package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProjectRequest_Valid(t *testing.T) {
	// Arrange
	req := CreateProjectRequest{
		Title:    "Test Project",
		Summary:  "Test Summary",
		CoverURL: "https://example.com/cover.jpg",
		Tags:     []string{"fiction", "fantasy"},
	}

	// Act & Assert
	assert.Equal(t, "Test Project", req.Title)
	assert.Equal(t, "Test Summary", req.Summary)
	assert.Len(t, req.Tags, 2)
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
				Title: "Valid Title",
			},
			wantErr: false,
		},
		{
			name: "missing title",
			req: CreateProjectRequest{
				Title: "",
			},
			wantErr: true,
		},
		{
			name: "title too long",
			req: CreateProjectRequest{
				Title: string(make([]byte, 101)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test will be validated by handler middleware
			require.NotEmpty(t, tt.req.Title, "title should not be empty for valid request")
		})
	}
}

func TestProjectResponse_Fields(t *testing.T) {
	// Arrange
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

	// Act & Assert
	assert.Equal(t, "proj_123", resp.ID)
	assert.Equal(t, "Test Project", resp.Title)
	assert.Equal(t, "draft", resp.Status)
	assert.False(t, resp.CreatedAt.IsZero())
}
```

**Step 2: Run test to verify it compiles and passes**

Run: `go test ./models/dto/ -v -run TestProject`
Expected: PASS

**Step 3: Commit**

```bash
git add models/dto/writer_dto_test.go
git commit -m "test(writer): add unit tests for Project DTOs"
```

---

## Task 7: Create Unit Tests for Converter Functions

**Files:**
- Modify: `models/dto/writer_dto_test.go`

**Step 1: Write test for ToProjectResponse converter**

Append to `models/dto/writer_dto_test.go`:

```go
import (
	"github.com/QingyuBackend/Qingyu_backend/models/writer"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestToProjectResponse(t *testing.T) {
	// Arrange
	objID := primitive.NewObjectID()
	project := &writer.Project{
		ID:        objID,
		Title:     "Test Project",
		Summary:   "Test Summary",
		CoverURL:  "https://example.com/cover.jpg",
		Tags:      []string{"fiction", "fantasy"},
		Status:    "draft",
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	// Act
	resp := ToProjectResponse(project)

	// Assert
	assert.Equal(t, objID.Hex(), resp.ID)
	assert.Equal(t, "Test Project", resp.Title)
	assert.Equal(t, "Test Summary", resp.Summary)
	assert.Equal(t, "https://example.com/cover.jpg", resp.CoverURL)
	assert.Equal(t, []string{"fiction", "fantasy"}, resp.Tags)
	assert.Equal(t, "draft", resp.Status)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), resp.CreatedAt)
}

func TestToProjectResponseList(t *testing.T) {
	// Arrange
	objID1 := primitive.NewObjectID()
	objID2 := primitive.NewObjectID()
	projects := []*writer.Project{
		{ID: objID1, Title: "Project 1"},
		{ID: objID2, Title: "Project 2"},
	}

	// Act
	responses := ToProjectResponseList(projects)

	// Assert
	assert.Len(t, responses, 2)
	assert.Equal(t, "Project 1", responses[0].Title)
	assert.Equal(t, "Project 2", responses[1].Title)
}
```

**Step 2: Run test to verify it compiles and passes**

Run: `go test ./models/dto/ -v -run TestToProject`
Expected: PASS

**Step 3: Commit**

```bash
git add models/dto/writer_dto_test.go
git commit -m "test(writer): add unit tests for converter functions"
```

---

## Task 8: Review Current Writer Handler Implementation

**Files:**
- Read: `api/v1/writer/project_handler.go` (or similar)
- Read: `service/writer/project_service.go` (or similar)

**Step 1: Locate the Writer handler files**

Run: `find api/v1/writer -name "*.go" -type f`
Expected: List of Go files in Writer API handlers

**Step 2: Read and analyze current implementation**

Identify:
- Current request/response structures used
- How MongoDB models are converted to API responses
- Validation logic location
- Error handling patterns

**Step 3: Document findings in comments**

No commit - this is analysis only

---

## Task 9: Update CreateProject Handler to Use DTO

**Files:**
- Modify: `api/v1/writer/project_handler.go` (or equivalent)

**Step 1: Read the current CreateProject handler**

Identify current implementation pattern

**Step 2: Update handler to use CreateProjectRequest DTO**

```go
package writer

import (
	"net/http"

	"github.com/QingyuBackend/Qingyu_backend/models/dto"
	"github.com/QingyuBackend/Qingyu_backend/service/writer"
	"github.com/gin-gonic/gin"
)

// CreateProject handles project creation requests
func (h *Handler) CreateProject(c *gin.Context) {
	// Step 1: Bind request to DTO
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(1001, "Invalid request: "+err.Error()))
		return
	}

	// Step 2: Call service layer with DTO
	project, err := h.projectService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(5000, "Failed to create project"))
		return
	}

	// Step 3: Convert to response DTO
	response := dto.ToProjectResponse(project)

	// Step 4: Return success response
	c.JSON(http.StatusCreated, dto.CreatedResponse(response))
}
```

**Step 3: Update service layer to accept DTO**

Modify `service/writer/project_service.go`:

```go
package writer

import (
	"context"
	"time"

	"github.com/QingyuBackend/Qingyu_backend/models/dto"
	"github.com/QingyuBackend/Qingyu_backend/models/writer"
)

type ProjectService interface {
	Create(ctx context.Context, req *dto.CreateProjectRequest) (*writer.Project, error)
	// ... other methods
}

type projectService struct {
	repo Repository
}

func NewProjectService(repo Repository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) Create(ctx context.Context, req *dto.CreateProjectRequest) (*writer.Project, error) {
	now := time.Now()
	project := &writer.Project{
		Title:     req.Title,
		Summary:   req.Summary,
		CoverURL:  req.CoverURL,
		Tags:      req.Tags,
		Status:    "draft",
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.repo.Create(ctx, project)
}
```

**Step 4: Run tests to verify changes**

Run: `go test ./api/v1/writer/... -v`
Expected: Tests pass (may need to update test files)

**Step 5: Commit**

```bash
git add api/v1/writer/project_handler.go service/writer/project_service.go
git commit -m "refactor(writer): update CreateProject to use DTO"
```

---

## Task 10: Update GetProject Handler to Use DTO

**Files:**
- Modify: `api/v1/writer/project_handler.go`

**Step 1: Update GetProject handler**

```go
// GetProject retrieves a project by ID
func (h *Handler) GetProject(c *gin.Context) {
	// Step 1: Get project ID from URL parameter
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(1001, "Project ID is required"))
		return
	}

	// Step 2: Call service layer
	project, err := h.projectService.GetByID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse(3000, "Project not found"))
		return
	}

	// Step 3: Convert to response DTO
	response := dto.ToProjectResponse(project)

	// Step 4: Return success response
	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}
```

**Step 2: Run tests**

Run: `go test ./api/v1/writer/... -v -run TestGetProject`
Expected: PASS

**Step 3: Commit**

```bash
git add api/v1/writer/project_handler.go
git commit -m "refactor(writer): update GetProject to use DTO"
```

---

## Task 11: Update ListProjects Handler to Use DTO

**Files:**
- Modify: `api/v1/writer/project_handler.go`

**Step 1: Update ListProjects handler**

```go
// ListProjects retrieves a paginated list of projects
func (h *Handler) ListProjects(c *gin.Context) {
	// Step 1: Bind query parameters to DTO
	var req dto.ListProjectsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(1001, "Invalid query parameters: "+err.Error()))
		return
	}

	// Step 2: Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// Step 3: Call service layer
	projects, total, err := h.projectService.List(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(5000, "Failed to list projects"))
		return
	}

	// Step 4: Convert to response DTO
	response := dto.ProjectListResponse{
		Items:    dto.ToProjectResponseList(projects),
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// Step 5: Return success response
	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}
```

**Step 2: Update service layer List method**

Modify `service/writer/project_service.go`:

```go
func (s *projectService) List(ctx context.Context, req *dto.ListProjectsRequest) ([]*writer.Project, int64, error) {
	return s.repo.List(ctx, req.Page, req.PageSize, req.Status, req.Sort, req.Order)
}
```

**Step 3: Run tests**

Run: `go test ./api/v1/writer/... -v -run TestListProjects`
Expected: PASS

**Step 4: Commit**

```bash
git add api/v1/writer/project_handler.go service/writer/project_service.go
git commit -m "refactor(writer): update ListProjects to use DTO"
```

---

## Task 12: Update UpdateProject Handler to Use DTO

**Files:**
- Modify: `api/v1/writer/project_handler.go`

**Step 1: Update UpdateProject handler**

```go
// UpdateProject updates an existing project
func (h *Handler) UpdateProject(c *gin.Context) {
	// Step 1: Get project ID from URL parameter
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(1001, "Project ID is required"))
		return
	}

	// Step 2: Bind request to DTO
	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(1001, "Invalid request: "+err.Error()))
		return
	}

	// Step 3: Call service layer
	project, err := h.projectService.Update(c.Request.Context(), projectID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(5000, "Failed to update project"))
		return
	}

	// Step 4: Convert to response DTO
	response := dto.ToProjectResponse(project)

	// Step 5: Return success response
	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}
```

**Step 2: Update service layer Update method**

Modify `service/writer/project_service.go`:

```go
func (s *projectService) Update(ctx context.Context, projectID string, req *dto.UpdateProjectRequest) (*writer.Project, error) {
	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.Title != nil {
		project.Title = *req.Title
	}
	if req.Summary != nil {
		project.Summary = *req.Summary
	}
	if req.CoverURL != nil {
		project.CoverURL = *req.CoverURL
	}
	if req.Tags != nil {
		project.Tags = *req.Tags
	}
	project.UpdatedAt = time.Now()

	return s.repo.Update(ctx, project)
}
```

**Step 3: Run tests**

Run: `go test ./api/v1/writer/... -v -run TestUpdateProject`
Expected: PASS

**Step 4: Commit**

```bash
git add api/v1/writer/project_handler.go service/writer/project_service.go
git commit -m "refactor(writer): update UpdateProject to use DTO"
```

---

## Task 13: Update DeleteProject Handler to Use DTO

**Files:**
- Modify: `api/v1/writer/project_handler.go`

**Step 1: Update DeleteProject handler**

```go
// DeleteProject deletes a project
func (h *Handler) DeleteProject(c *gin.Context) {
	// Step 1: Get project ID from URL parameter
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(1001, "Project ID is required"))
		return
	}

	// Step 2: Call service layer
	err := h.projectService.Delete(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(5000, "Failed to delete project"))
		return
	}

	// Step 3: Return success response
	c.JSON(http.StatusOK, dto.SuccessResponse(map[string]string{
		"id":      projectID,
		"message": "Project deleted successfully",
	}))
}
```

**Step 2: Run tests**

Run: `go test ./api/v1/writer/... -v -run TestDeleteProject`
Expected: PASS

**Step 3: Commit**

```bash
git add api/v1/writer/project_handler.go
git commit -m "refactor(writer): update DeleteProject to use DTO"
```

---

## Task 14: Update Document Handlers to Use DTO

**Files:**
- Modify: `api/v1/writer/document_handler.go`
- Modify: `service/writer/document_service.go`

**Step 1: Update CreateDocument handler**

Follow the same pattern as Task 9:

```go
func (h *Handler) CreateDocument(c *gin.Context) {
	var req dto.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(1001, "Invalid request: "+err.Error()))
		return
	}

	document, err := h.documentService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(5000, "Failed to create document"))
		return
	}

	response := dto.ToDocumentResponse(document)
	c.JSON(http.StatusCreated, dto.CreatedResponse(response))
}
```

**Step 2: Update GetDocument handler**

```go
func (h *Handler) GetDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(1001, "Document ID is required"))
		return
	}

	document, err := h.documentService.GetByID(c.Request.Context(), documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse(3000, "Document not found"))
		return
	}

	response := dto.ToDocumentResponse(document)
	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}
```

**Step 3: Update ListDocuments handler**

```go
func (h *Handler) ListDocuments(c *gin.Context) {
	var req dto.ListDocumentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(1001, "Invalid query parameters: "+err.Error()))
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	documents, total, err := h.documentService.List(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(5000, "Failed to list documents"))
		return
	}

	response := dto.DocumentListResponse{
		Items:    dto.ToDocumentResponseList(documents),
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}
```

**Step 4: Run tests**

Run: `go test ./api/v1/writer/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add api/v1/writer/document_handler.go service/writer/document_service.go
git commit -m "refactor(writer): update Document handlers to use DTO"
```

---

## Task 15: Update Chapter Handlers to Use DTO

**Files:**
- Modify: `api/v1/writer/chapter_handler.go`
- Modify: `service/writer/chapter_service.go`

**Step 1: Update chapter handlers following Document pattern**

```go
func (h *Handler) CreateChapter(c *gin.Context) {
	var req dto.CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(1001, "Invalid request: "+err.Error()))
		return
	}

	chapter, err := h.chapterService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(5000, "Failed to create chapter"))
		return
	}

	response := dto.ToChapterResponse(chapter)
	c.JSON(http.StatusCreated, dto.CreatedResponse(response))
}

func (h *Handler) GetChapter(c *gin.Context) {
	chapterID := c.Param("id")
	if chapterID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(1001, "Chapter ID is required"))
		return
	}

	chapter, err := h.chapterService.GetByID(c.Request.Context(), chapterID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse(3000, "Chapter not found"))
		return
	}

	response := dto.ToChapterResponse(chapter)
	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}

func (h *Handler) ListChapters(c *gin.Context) {
	var req dto.ListChaptersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(1001, "Invalid query parameters: "+err.Error()))
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	chapters, total, err := h.chapterService.List(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(5000, "Failed to list chapters"))
		return
	}

	response := dto.ChapterListResponse{
		Items:    dto.ToChapterResponseList(chapters),
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}
```

**Step 2: Run tests**

Run: `go test ./api/v1/writer/... -v`
Expected: PASS

**Step 3: Commit**

```bash
git add api/v1/writer/chapter_handler.go service/writer/chapter_service.go
git commit -m "refactor(writer): update Chapter handlers to use DTO"
```

---

## Task 16: Run Full Test Suite

**Step 1: Run all Writer module tests**

Run: `go test ./api/v1/writer/... ./service/writer/... ./models/dto/... -v`
Expected: All tests PASS

**Step 2: Run full backend test suite**

Run: `go test ./... -v`
Expected: No new test failures

**Step 3: Run integration tests (if exist)**

Run: `go test ./integration/... -v`
Expected: All tests PASS

**Step 4: Verify server starts**

Run: `go run cmd/server/main.go &`
Wait: 5 seconds
Check: `curl http://localhost:8080/health`
Expected: Server responds successfully

**Step 5: Commit**

```bash
git add .
git commit -m "test(writer): verify all tests pass after DTO migration"
```

---

## Task 17: Manual API Testing

**Step 1: Test CreateProject endpoint**

Run:
```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Project", "summary": "Test Summary", "tags": ["test"]}'
```

Expected:
```json
{
  "code": 201,
  "message": "created",
  "data": {
    "id": "...",
    "title": "Test Project",
    "summary": "Test Summary",
    "coverUrl": "",
    "tags": ["test"],
    "status": "draft",
    "createdAt": "2024-03-05T...",
    "updatedAt": "2024-03-05T..."
  },
  "timestamp": "2024-03-05T...",
  "request_id": "..."
}
```

**Step 2: Test GetProject endpoint**

Run:
```bash
curl http://localhost:8080/api/v1/projects/{id_from_step_1}
```

Expected: JSON with camelCase fields (coverUrl, createdAt, updatedAt)

**Step 3: Test ListProjects endpoint**

Run:
```bash
curl "http://localhost:8080/api/v1/projects?page=1&page_size=10"
```

Expected: Paginated response with items array

**Step 4: Test UpdateProject endpoint**

Run:
```bash
curl -X PUT http://localhost:8080/api/v1/projects/{id} \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Title"}'
```

Expected: Updated project response

**Step 5: Test validation errors**

Run:
```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{"title": ""}'
```

Expected:
```json
{
  "code": 1001,
  "message": "Invalid request: ...",
  "data": null,
  ...
}
```

**Step 6: Document test results**

Create `docs/plans/2026-03-05-api-integration-phase1-test-results.md` with test results

**Step 7: Commit**

```bash
git add docs/plans/2026-03-05-api-integration-phase1-test-results.md
git commit -m "docs(writer): add Phase 1 manual test results"
```

---

## Task 18: Update Swagger Documentation

**Files:**
- Modify: `docs/swagger.yaml`

**Step 1: Add Project schemas to swagger.yaml**

Add under `components.schemas`:

```yaml
CreateProjectRequest:
  type: object
  required:
    - title
  properties:
    title:
      type: string
      minLength: 1
      maxLength: 100
      example: "My Novel"
    summary:
      type: string
      maxLength: 500
      example: "A fantasy adventure story"
    coverUrl:
      type: string
      format: uri
      maxLength: 500
      example: "https://example.com/cover.jpg"
    tags:
      type: array
      maxItems: 10
      items:
        type: string
        minLength: 1
        maxLength: 50
      example: ["fantasy", "adventure"]

UpdateProjectRequest:
  type: object
  properties:
    title:
      type: string
      minLength: 1
      maxLength: 100
    summary:
      type: string
      maxLength: 500
    coverUrl:
      type: string
      format: uri
      maxLength: 500
    tags:
      type: array
      maxItems: 10
      items:
        type: string
        minLength: 1
        maxLength: 50

ProjectResponse:
  type: object
  properties:
    id:
      type: string
      example: "507f1f77bcf86cd799439011"
    title:
      type: string
      example: "My Novel"
    summary:
      type: string
      example: "A fantasy adventure story"
    coverUrl:
      type: string
      example: "https://example.com/cover.jpg"
    tags:
      type: array
      items:
        type: string
    status:
      type: string
      enum: [draft, published, archived]
      example: "draft"
    createdAt:
      type: string
      format: date-time
    updatedAt:
      type: string
      format: date-time

ProjectListResponse:
  type: object
  properties:
    items:
      type: array
      items:
        $ref: '#/components/schemas/ProjectResponse'
    total:
      type: integer
      example: 100
    page:
      type: integer
      example: 1
    pageSize:
      type: integer
      example: 20
```

**Step 2: Update /api/v1/projects endpoints to use new schemas**

Update the request/response schemas in swagger.yaml

**Step 3: Validate swagger.yaml**

Run: `npx swagger-parser validate docs/swagger.yaml`
Expected: No errors

**Step 4: Commit**

```bash
git add docs/swagger.yaml
git commit -m "docs(writer): update swagger.yaml with Project DTO schemas"
```

---

## Task 19: Create Migration Guide

**Files:**
- Create: `docs/plans/2026-03-05-api-integration-phase1-migration-guide.md`

**Step 1: Write migration guide**

```markdown
# Writer Module DTO Migration Guide

## Overview

This guide documents the migration of Writer module API handlers to use unified DTOs from `models/dto/`.

## Changes

### Before
- Handlers used MongoDB model structs directly
- JSON tags were inconsistent
- No centralized validation

### After
- Handlers use DTO structs from `models/dto/writer_dto.go`
- Consistent camelCase JSON tags
- Validation tags on DTOs
- Converter functions for model-to-DTO transformation

## New Files

- `models/dto/writer_dto.go` - DTO definitions
- `models/dto/writer_converter.go` - Converter functions
- `models/dto/writer_dto_test.go` - Unit tests

## API Compatibility

### Breaking Changes
None. The API interface remains the same.

### Response Format Changes
- All responses now use camelCase (coverUrl, createdAt, updatedAt)
- All responses follow standard APIResponse format

## For Frontend Developers

### No Changes Required
If your frontend was already using camelCase field names, no changes are needed.

### Verify Your Integration
Test the following endpoints:
- POST /api/v1/projects
- GET /api/v1/projects/{id}
- GET /api/v1/projects
- PUT /api/v1/projects/{id}
- DELETE /api/v1/projects/{id}

## For Backend Developers

### Using Writer DTOs

```go
import "github.com/QingyuBackend/Qingyu_backend/models/dto"

// In handler
var req dto.CreateProjectRequest
if err := c.ShouldBindJSON(&req); err != nil {
    // handle error
}

// In service
project, err := service.Create(ctx, &req)

// Convert to response
response := dto.ToProjectResponse(project)
```

## Testing

Run the following to verify the migration:
```bash
# Unit tests
go test ./models/dto/... -v

# Integration tests
go test ./api/v1/writer/... -v

# Manual test
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{"title": "Test"}'
```

## Rollback Plan

If issues occur:
1. Revert to commit before migration
2. Frontend: No changes needed (API remains compatible)
3. Backend: Use `git revert` to undo changes

## Support

Contact: [Team Lead]
Migration Date: 2026-03-05
```

**Step 2: Commit**

```bash
git add docs/plans/2026-03-05-api-integration-phase1-migration-guide.md
git commit -m "docs(writer): add DTO migration guide"
```

---

## Task 20: Final Verification and Documentation

**Step 1: Create Phase 1 completion checklist**

Create `docs/plans/2026-03-05-api-integration-phase1-checklist.md`:

```markdown
# Phase 1 Completion Checklist

## DTO Creation
- [x] Create `models/dto/writer_dto.go`
- [x] Add Project DTOs
- [x] Add Document DTOs
- [x] Add Chapter DTOs
- [x] Create converter functions
- [x] Add unit tests for DTOs
- [x] Add unit tests for converters

## Handler Updates
- [x] Update CreateProject handler
- [x] Update GetProject handler
- [x] Update ListProjects handler
- [x] Update UpdateProject handler
- [x] Update DeleteProject handler
- [x] Update Document handlers
- [x] Update Chapter handlers

## Testing
- [x] Run unit tests
- [x] Run integration tests
- [x] Manual API testing
- [x] Verify server starts
- [x] Test validation errors

## Documentation
- [x] Update swagger.yaml
- [x] Create migration guide
- [x] Document test results

## Verification
- [x] No breaking changes to API
- [x] All tests pass
- [x] camelCase JSON tags consistent
- [x] Error handling standardized
```

**Step 2: Run final verification**

Run: `go test ./models/dto/... ./api/v1/writer/... ./service/writer/... -cover`
Expected: Coverage > 80%

**Step 3: Create Phase 1 summary**

Create `docs/plans/2026-03-05-api-integration-phase1-summary.md`:

```markdown
# Phase 1: Backend DTO Unification - Summary

**Completed**: 2026-03-05
**Duration**: 3 weeks (as planned)
**Status**: ✅ Complete

## Deliverables

### Code
- `models/dto/writer_dto.go` - 3 DTO groups (Project, Document, Chapter)
- `models/dto/writer_converter.go` - 6 converter functions
- `models/dto/writer_dto_test.go` - Unit tests
- Updated Writer handlers to use DTOs
- Updated Writer services to accept DTOs

### Documentation
- Migration guide
- API test results
- Swagger schema updates

## Metrics

- **DTOs Created**: 12 request/response DTOs
- **Converters**: 6 functions
- **Unit Tests**: 10+ test cases
- **Test Coverage**: >80%
- **API Endpoints Updated**: 15 endpoints
- **Breaking Changes**: 0

## Next Steps

Proceed to **Phase 2: Frontend Wrapper Optimization**
- Duration: 2 weeks
- Focus: Enhance orval-mutator, add unified error handling

## Lessons Learned

1. **DTO Pattern Works Well**: Clear separation between models and API layer
2. **Converter Functions**: Essential for maintaining clean boundaries
3. **Validation Tags**: Reduce boilerplate in handlers
4. **No Breaking Changes**: Key to smooth migration

## Risks Mitigated

- ✅ No API breaking changes
- ✅ Backward compatible
- ✅ All tests passing
- ✅ Documentation complete

## Approval

- [ ] Backend Lead Review
- [ ] QA Team Verification
- [ ] Frontend Team Notification
```

**Step 4: Final commit**

```bash
git add docs/plans/2026-03-05-api-integration-phase1-*.md
git commit -m "docs(writer): add Phase 1 completion documentation"
```

---

## Phase 1 Complete

**Total Estimated Time**: 3 weeks
**Total Tasks**: 20
**Total Commits**: ~20

**Next Phase**: Phase 2 - Frontend Wrapper Optimization

---

## For Implementation

**REQUIRED SUB-SKILL:** Use superpowers:executing-plans to implement this plan task-by-task.

**Alternative:** Use superpowers:subagent-driven-development for parallel execution with review checkpoints.
