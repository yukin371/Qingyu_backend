package seeds

import (
	"context"
	"fmt"
	"strings"
	"time"

	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/users"
	writerModel "Qingyu_backend/models/writer"
	writerBase "Qingyu_backend/models/writer/base"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	publicationFlowAuthorUsername  = "author_new"
	publicationFlowAdminUsername   = "admin"
	publicationFlowProjectTitle    = "联调发布示例项目"
	publicationFlowProjectSummary  = "用于平台写作、发布、审核、阅读主链路联调的最小项目数据。"
	publicationFlowPrimaryDocTitle = "第1章 风起青川"
	publicationFlowSecondDocTitle  = "第2章 初入书城"
)

type PublicationFlowSeedInfo struct {
	AuthorUsername string
	AdminUsername  string
	ProjectID      string
	DocumentIDs    []string
}

type publicationFlowDocumentSeed struct {
	Title   string
	Order   int
	Content string
}

var publicationFlowDocuments = []publicationFlowDocumentSeed{
	{
		Title: publicationFlowPrimaryDocTitle,
		Order: 1,
		Content: strings.TrimSpace(`青川城的钟声在黎明前敲了三下，沈砚放下笔，重新检查面前的稿纸。

这是他第一次准备把作品公开发布。

窗外天色将明，街巷里还带着夜雨后的凉意。他知道，只要按下提交，这个只存在于案头的故事，就会真正走向读者。`),
	},
	{
		Title: publicationFlowSecondDocTitle,
		Order: 2,
		Content: strings.TrimSpace(`书城的审核列表刷新时，新的待审稿件静静躺在最上方。

另一端，读者首页仍是一片安静，直到管理员点击通过，作品才会沿着既定流程进入书籍与章节集合。

这正是联调需要验证的路径。`),
	},
}

func SeedPublicationFlowData(ctx context.Context, db *mongo.Database) (*PublicationFlowSeedInfo, error) {
	usersCollection := db.Collection("users")
	projectsCollection := db.Collection("projects")
	documentsCollection := db.Collection("documents")
	contentsCollection := db.Collection("document_contents")

	author, err := findSeedUserByUsername(ctx, usersCollection, publicationFlowAuthorUsername)
	if err != nil {
		return nil, fmt.Errorf("find publication author: %w", err)
	}

	admin, err := findSeedUserByUsername(ctx, usersCollection, publicationFlowAdminUsername)
	if err != nil {
		return nil, fmt.Errorf("find publication admin: %w", err)
	}

	project, created, err := ensurePublicationFlowProject(ctx, projectsCollection, author)
	if err != nil {
		return nil, err
	}
	if created {
		fmt.Printf("✓ 创建联调项目: %s (%s)\n", project.Title, project.ID.Hex())
	} else {
		fmt.Printf("✓ 复用联调项目: %s (%s)\n", project.Title, project.ID.Hex())
	}

	documentIDs := make([]string, 0, len(publicationFlowDocuments))
	totalWords := 0
	for _, seedDoc := range publicationFlowDocuments {
		doc, docCreated, err := ensurePublicationFlowDocument(ctx, documentsCollection, project.ID, seedDoc)
		if err != nil {
			return nil, err
		}
		if docCreated {
			fmt.Printf("  ├─ 创建文档: %s (%s)\n", doc.Title, doc.ID.Hex())
		} else {
			fmt.Printf("  ├─ 复用文档: %s (%s)\n", doc.Title, doc.ID.Hex())
		}

		if err := ensurePublicationFlowContent(ctx, contentsCollection, doc.ID, author.ID.Hex(), seedDoc.Content); err != nil {
			return nil, err
		}

		documentIDs = append(documentIDs, doc.ID.Hex())
		totalWords += countSeedContentUnits(seedDoc.Content)
	}

	if err := updatePublicationFlowProjectStats(ctx, projectsCollection, project.ID, len(publicationFlowDocuments), totalWords); err != nil {
		return nil, err
	}

	return &PublicationFlowSeedInfo{
		AuthorUsername: author.Username,
		AdminUsername:  admin.Username,
		ProjectID:      project.ID.Hex(),
		DocumentIDs:    documentIDs,
	}, nil
}

func findSeedUserByUsername(ctx context.Context, collection *mongo.Collection, username string) (*users.User, error) {
	var user users.User
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user %s not found", username)
		}
		return nil, err
	}
	return &user, nil
}

func ensurePublicationFlowProject(ctx context.Context, collection *mongo.Collection, author *users.User) (*writerModel.Project, bool, error) {
	var existing writerModel.Project
	err := collection.FindOne(ctx, bson.M{
		"author_id": author.ID.Hex(),
		"title":     publicationFlowProjectTitle,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}).Decode(&existing)
	if err == nil {
		return &existing, false, nil
	}
	if err != mongo.ErrNoDocuments {
		return nil, false, fmt.Errorf("find publication flow project: %w", err)
	}

	now := time.Now()
	project := &writerModel.Project{
		IdentifiedEntity: writerBase.IdentifiedEntity{ID: primitive.NewObjectID()},
		OwnedEntity:      writerBase.OwnedEntity{AuthorID: author.ID.Hex()},
		TitledEntity:     shared.TitledEntity{Title: publicationFlowProjectTitle},
		Timestamps:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
		WritingType:      "novel",
		Summary:          publicationFlowProjectSummary,
		Status:           writerModel.StatusSerializing,
		Category:         "玄幻",
		Tags:             []string{"联调", "发布", "审核"},
		Visibility:       writerModel.VisibilityPublic,
		Statistics: writerModel.ProjectStats{
			LastUpdateAt: now,
		},
		Settings: writerModel.ProjectSettings{
			AutoBackup:     true,
			BackupInterval: 24,
			WordCountGoal:  3000,
		},
	}

	if err := project.Validate(); err != nil {
		return nil, false, fmt.Errorf("validate publication flow project: %w", err)
	}

	if _, err := collection.InsertOne(ctx, project); err != nil {
		return nil, false, fmt.Errorf("insert publication flow project: %w", err)
	}
	return project, true, nil
}

func ensurePublicationFlowDocument(ctx context.Context, collection *mongo.Collection, projectID primitive.ObjectID, seedDoc publicationFlowDocumentSeed) (*writerModel.Document, bool, error) {
	var existing writerModel.Document
	err := collection.FindOne(ctx, bson.M{
		"project_id": projectID,
		"title":      seedDoc.Title,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}).Decode(&existing)
	if err == nil {
		return &existing, false, nil
	}
	if err != mongo.ErrNoDocuments {
		return nil, false, fmt.Errorf("find publication flow document %s: %w", seedDoc.Title, err)
	}

	doc := &writerModel.Document{
		ProjectID: projectID,
		Title:     seedDoc.Title,
		Type:      writerModel.TypeChapter,
		Level:     0,
		Order:     seedDoc.Order,
		Status:    writerModel.DocumentStatusCompleted,
		WordCount: countSeedContentUnits(seedDoc.Content),
	}
	doc.TouchForCreate()

	if err := doc.Validate("novel"); err != nil {
		return nil, false, fmt.Errorf("validate publication flow document %s: %w", seedDoc.Title, err)
	}

	if _, err := collection.InsertOne(ctx, doc); err != nil {
		return nil, false, fmt.Errorf("insert publication flow document %s: %w", seedDoc.Title, err)
	}

	return doc, true, nil
}

func ensurePublicationFlowContent(ctx context.Context, collection *mongo.Collection, documentID primitive.ObjectID, editorID string, content string) error {
	wordCount := countSeedContentUnits(content)
	now := time.Now()

	var existing writerModel.DocumentContent
	err := collection.FindOne(ctx, bson.M{"document_id": documentID}).Decode(&existing)
	if err == nil {
		_, updateErr := collection.UpdateOne(ctx, bson.M{"_id": existing.ID}, bson.M{
			"$set": bson.M{
				"content":        content,
				"content_type":   "markdown",
				"word_count":     wordCount,
				"char_count":     wordCount,
				"last_edited_by": editorID,
				"last_saved_at":  now,
				"updated_at":     now,
			},
		})
		if updateErr != nil {
			return fmt.Errorf("update publication flow content: %w", updateErr)
		}
		return nil
	}
	if err != mongo.ErrNoDocuments {
		return fmt.Errorf("find publication flow content: %w", err)
	}

	docContent := &writerModel.DocumentContent{
		DocumentID:   documentID,
		Content:      content,
		ContentType:  "markdown",
		WordCount:    wordCount,
		CharCount:    wordCount,
		LastEditedBy: editorID,
	}
	docContent.TouchForCreate()

	if err := docContent.Validate(); err != nil {
		return fmt.Errorf("validate publication flow content: %w", err)
	}

	if _, err := collection.InsertOne(ctx, docContent); err != nil {
		return fmt.Errorf("insert publication flow content: %w", err)
	}

	return nil
}

func updatePublicationFlowProjectStats(ctx context.Context, collection *mongo.Collection, projectID primitive.ObjectID, documentCount int, totalWords int) error {
	now := time.Now()
	_, err := collection.UpdateOne(ctx, bson.M{"_id": projectID}, bson.M{
		"$set": bson.M{
			"statistics.document_count": documentCount,
			"statistics.chapter_count":  documentCount,
			"statistics.total_words":    totalWords,
			"statistics.last_update_at": now,
			"updated_at":                now,
		},
	})
	if err != nil {
		return fmt.Errorf("update publication flow project stats: %w", err)
	}
	return nil
}

func countSeedContentUnits(content string) int {
	return len([]rune(strings.TrimSpace(content)))
}
