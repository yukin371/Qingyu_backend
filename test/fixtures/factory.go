package fixtures

import (
	"fmt"
	"time"

	"Qingyu_backend/models/document"
	"Qingyu_backend/models/users"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ UserFactory ============

// UserFactory 用户数据工厂
type UserFactory struct {
	counter int
}

// NewUserFactory 创建用户工厂
func NewUserFactory() *UserFactory {
	return &UserFactory{counter: 0}
}

// Create 创建用户（支持自定义选项）
func (f *UserFactory) Create(opts ...func(*users.User)) *users.User {
	f.counter++
	user := &users.User{
		ID:        primitive.NewObjectID().Hex(),
		Username:  fmt.Sprintf("user%d", f.counter),
		Email:     fmt.Sprintf("user%d@test.com", f.counter),
		Password:  "hashed_password_" + fmt.Sprint(f.counter),
		Role:      "user",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 应用自定义选项
	for _, opt := range opts {
		opt(user)
	}

	return user
}

// CreateBatch 批量创建用户
func (f *UserFactory) CreateBatch(count int) []*users.User {
	result := make([]*users.User, count)
	for i := 0; i < count; i++ {
		result[i] = f.Create()
	}
	return result
}

// CreateAdmin 创建管理员用户
func (f *UserFactory) CreateAdmin() *users.User {
	return f.Create(func(u *users.User) {
		u.Role = "admin"
		u.Username = fmt.Sprintf("admin%d", f.counter)
	})
}

// CreateAuthor 创建作者用户
func (f *UserFactory) CreateAuthor() *users.User {
	return f.Create(func(u *users.User) {
		u.Role = "author"
		u.Username = fmt.Sprintf("author%d", f.counter)
	})
}

// ============ ProjectFactory ============

// ProjectFactory 项目数据工厂
type ProjectFactory struct {
	counter int
}

// NewProjectFactory 创建项目工厂
func NewProjectFactory() *ProjectFactory {
	return &ProjectFactory{counter: 0}
}

// Create 创建项目
func (f *ProjectFactory) Create(authorID string, opts ...func(*document.Project)) *document.Project {
	f.counter++
	project := &document.Project{
		ID:        primitive.NewObjectID().Hex(),
		Title:     fmt.Sprintf("测试项目 %d", f.counter),
		Summary:   fmt.Sprintf("这是第%d个测试项目", f.counter),
		AuthorID:  authorID,
		Status:    document.StatusDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 应用自定义选项
	for _, opt := range opts {
		opt(project)
	}

	return project
}

// CreateBatch 批量创建项目
func (f *ProjectFactory) CreateBatch(authorID string, count int) []*document.Project {
	result := make([]*document.Project, count)
	for i := 0; i < count; i++ {
		result[i] = f.Create(authorID)
	}
	return result
}

// CreateNovel 创建小说项目
func (f *ProjectFactory) CreateNovel(authorID string) *document.Project {
	return f.Create(authorID, func(p *document.Project) {
		p.Title = fmt.Sprintf("小说项目 %d", f.counter)
		p.Summary = "这是一个小说项目"
	})
}

// ============ DocumentFactory ============

// DocumentFactory 文档数据工厂
type DocumentFactory struct {
	counter int
}

// NewDocumentFactory 创建文档工厂
func NewDocumentFactory() *DocumentFactory {
	return &DocumentFactory{counter: 0}
}

// Create 创建文档元数据
// 注意：此方法只创建Document元数据，不包含内容
// 如需创建文档内容，请使用CreateDocumentContent方法
func (f *DocumentFactory) Create(projectID string, opts ...func(*document.Document)) *document.Document {
	f.counter++
	doc := &document.Document{
		ID:        primitive.NewObjectID().Hex(),
		ProjectID: projectID,
		Title:     fmt.Sprintf("第%d章", f.counter),
		Type:      document.TypeChapter,
		Status:    "draft",
		WordCount: 1000 + f.counter*100,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 应用自定义选项
	for _, opt := range opts {
		opt(doc)
	}

	return doc
}

// CreateDocumentContent 创建文档内容
func (f *DocumentFactory) CreateDocumentContent(documentID string) *document.DocumentContent {
	return &document.DocumentContent{
		ID:          primitive.NewObjectID().Hex(),
		DocumentID:  documentID,
		Content:     fmt.Sprintf("这是文档%s的内容...", documentID),
		ContentType: "markdown",
		WordCount:   1000,
		CharCount:   1000,
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		LastSavedAt: time.Now(),
	}
}

// CreateBatch 批量创建文档
func (f *DocumentFactory) CreateBatch(projectID string, count int) []*document.Document {
	result := make([]*document.Document, count)
	for i := 0; i < count; i++ {
		result[i] = f.Create(projectID)
	}
	return result
}

// CreatePublished 创建已发布文档
func (f *DocumentFactory) CreatePublished(projectID string) *document.Document {
	return f.Create(projectID, func(d *document.Document) {
		d.Status = "published"
	})
}

// ============ 使用示例 ============

/*
使用示例：

func TestExample(t *testing.T) {
    // 创建工厂
    userFactory := fixtures.NewUserFactory()
    projectFactory := fixtures.NewProjectFactory()
    bookFactory := fixtures.NewBookFactory()

    // 创建用户
    user1 := userFactory.Create()
    admin := userFactory.CreateAdmin()
    users := userFactory.CreateBatch(5)

    // 创建自定义用户
    customUser := userFactory.Create(func(u *users.User) {
        u.Username = "custom"
        u.Email = "custom@test.com"
    })

    // 创建项目
    project1 := projectFactory.Create(user1.ID)
    projects := projectFactory.CreateBatch(user1.ID, 3)

    // 创建书籍
    book1 := bookFactory.Create()
    popularBook := bookFactory.CreatePopular()
    books := bookFactory.CreateBatch(10)
}
*/
