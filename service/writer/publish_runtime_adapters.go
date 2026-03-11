package writer

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	bookstoreModel "Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/dto"
	writerModel "Qingyu_backend/models/writer"
	"Qingyu_backend/repository"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	baseInterfaces "Qingyu_backend/service/interfaces/base"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PublishProjectRepositoryAdapter 适配现有项目仓储到发布服务接口。
type PublishProjectRepositoryAdapter struct {
	repo writerRepo.ProjectRepository
}

func NewPublishProjectRepositoryAdapter(repo writerRepo.ProjectRepository) *PublishProjectRepositoryAdapter {
	return &PublishProjectRepositoryAdapter{repo: repo}
}

func (a *PublishProjectRepositoryAdapter) FindByID(ctx context.Context, id string) (*writerModel.Project, error) {
	return a.repo.GetByID(ctx, id)
}

// PublishDocumentRepositoryAdapter 适配现有文档仓储到发布服务接口。
type PublishDocumentRepositoryAdapter struct {
	repo writerRepo.DocumentRepository
	db   *mongo.Database
}

func NewPublishDocumentRepositoryAdapter(repo writerRepo.DocumentRepository, db *mongo.Database) *PublishDocumentRepositoryAdapter {
	return &PublishDocumentRepositoryAdapter{repo: repo, db: db}
}

func (a *PublishDocumentRepositoryAdapter) FindByID(ctx context.Context, id string) (*writerModel.Document, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *PublishDocumentRepositoryAdapter) FindByProjectID(ctx context.Context, projectID string) ([]*writerModel.Document, error) {
	objectID, err := repository.ParseID(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}

	filter := bson.M{"project_id": objectID}
	opts := options.Find().
		SetSort(bson.D{{Key: "order", Value: 1}, {Key: "created_at", Value: 1}})

	cursor, err := a.db.Collection("documents").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []*writerModel.Document
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

// LocalBookstoreClient 直接写入本地书城集合，用于发布 MVP。
type LocalBookstoreClient struct {
	db *mongo.Database
}

type publishedParagraph struct {
	Order   int
	Content string
	Format  string
}

type PublishEventBusAdapter struct {
	bus baseInterfaces.EventBus
}

func NewLocalBookstoreClient(db *mongo.Database) *LocalBookstoreClient {
	return &LocalBookstoreClient{db: db}
}

func NewPublishEventBusAdapter(bus baseInterfaces.EventBus) *PublishEventBusAdapter {
	return &PublishEventBusAdapter{bus: bus}
}

func (a *PublishEventBusAdapter) PublishAsync(ctx context.Context, event interface{}) error {
	if a == nil || a.bus == nil || event == nil {
		return nil
	}
	if typed, ok := event.(baseInterfaces.Event); ok {
		return a.bus.PublishAsync(ctx, typed)
	}
	eventType, source := resolvePublishEventMetadata(event)
	return a.bus.PublishAsync(ctx, genericPublishEvent{
		eventType: eventType,
		data:      event,
		timestamp: time.Now(),
		source:    source,
	})
}

func resolvePublishEventMetadata(event interface{}) (string, string) {
	const (
		defaultEventType = "writer.publish"
		defaultSource    = "writer.publish_service"
	)

	payload, ok := event.(map[string]interface{})
	if !ok {
		return defaultEventType, defaultSource
	}

	eventType, _ := payload["eventType"].(string)
	if eventType == "" {
		eventType = defaultEventType
	}

	source, _ := payload["source"].(string)
	if source == "" {
		source = defaultSource
	}

	return eventType, source
}

func (c *LocalBookstoreClient) PublishProject(ctx context.Context, req *BookstorePublishProjectRequest) (*BookstorePublishResponse, error) {
	project, err := c.getProject(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}

	book, err := c.upsertBook(ctx, project, req)
	if err != nil {
		return nil, err
	}

	if err := c.rebuildProjectChapters(ctx, book, req.ProjectID, req.FreeChapters); err != nil {
		return nil, err
	}

	return &BookstorePublishResponse{
		Success:       true,
		BookstoreID:   "local",
		BookstoreName: "Local Bookstore",
		ExternalID:    book.ID.Hex(),
		Message:       "published",
	}, nil
}

func (c *LocalBookstoreClient) UnpublishProject(ctx context.Context, projectID, bookstoreID string) error {
	book, err := c.getBookByProjectID(ctx, projectID)
	if err != nil || book == nil {
		return err
	}

	chapterIDs, err := c.listChapterObjectIDs(ctx, book.ID.Hex())
	if err != nil {
		return err
	}
	if len(chapterIDs) > 0 {
		_, err = c.db.Collection("chapter_contents").DeleteMany(ctx, bson.M{"chapter_id": bson.M{"$in": chapterIDs}})
		if err != nil {
			return err
		}
	}
	if _, err := c.db.Collection("chapters").DeleteMany(ctx, bson.M{"book_id": book.ID.Hex()}); err != nil {
		return err
	}
	_, err = c.db.Collection("books").DeleteOne(ctx, bson.M{"_id": book.ID})
	return err
}

func (c *LocalBookstoreClient) PublishChapter(ctx context.Context, req *BookstorePublishChapterRequest) (*BookstorePublishResponse, error) {
	book, err := c.ensureProjectBook(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}
	chapterID, err := c.upsertSingleChapter(ctx, book, req)
	if err != nil {
		return nil, err
	}
	return &BookstorePublishResponse{
		Success:       true,
		BookstoreID:   book.ID.Hex(),
		BookstoreName: "Local Bookstore",
		ExternalID:    chapterID,
		Message:       "published",
	}, nil
}

func (c *LocalBookstoreClient) UnpublishChapter(ctx context.Context, chapterID, bookstoreID string) error {
	chapter, err := c.getChapterByProjectChapterID(ctx, chapterID)
	if err != nil || chapter == nil {
		return err
	}

	objectID := chapter.ID
	if _, err := c.db.Collection("chapter_contents").DeleteMany(ctx, bson.M{"chapter_id": objectID}); err != nil {
		return err
	}
	_, err = c.db.Collection("chapters").DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (c *LocalBookstoreClient) UpdateChapter(ctx context.Context, req *BookstoreUpdateChapterRequest) error {
	chapterObjectID, err := repository.ParseID(req.ChapterID)
	if err != nil {
		return fmt.Errorf("invalid chapter id: %w", err)
	}

	updates := bson.M{
		"title":       req.ChapterTitle,
		"chapter_num": req.ChapterNumber,
		"is_free":     req.IsFree,
		"updated_at":  time.Now(),
	}
	if _, err := c.db.Collection("chapters").UpdateOne(ctx, bson.M{"_id": chapterObjectID}, bson.M{"$set": updates}); err != nil {
		return err
	}
	_, err = c.db.Collection("chapter_contents").UpdateOne(
		ctx,
		bson.M{"chapter_id": chapterObjectID},
		bson.M{"$set": bson.M{"content": req.Content, "word_count": len([]rune(req.Content)), "updated_at": time.Now()}, "$inc": bson.M{"version": 1}},
	)
	return err
}

func (c *LocalBookstoreClient) GetStatistics(ctx context.Context, projectID, bookstoreID string) (*dto.PublicationStatistics, error) {
	book, err := c.getBookByProjectID(ctx, projectID)
	if err != nil || book == nil {
		return &dto.PublicationStatistics{}, err
	}

	return &dto.PublicationStatistics{
		TotalViews:     book.ViewCount,
		TotalLikes:     0,
		TotalComments:  0,
		TotalShares:    0,
		TotalFavorites: 0,
		LastSyncTime:   safeTime(book.LastSyncedAt),
	}, nil
}

func (c *LocalBookstoreClient) ensureProjectBook(ctx context.Context, projectID string) (*bookstoreModel.Book, error) {
	book, err := c.getBookByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if book != nil {
		return book, nil
	}
	project, err := c.getProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return c.upsertBook(ctx, project, &BookstorePublishProjectRequest{ProjectID: projectID, ProjectTitle: project.Title, AuthorID: project.AuthorID, Description: project.Summary, CoverImage: project.CoverURL, PublishType: "serial"})
}

func (c *LocalBookstoreClient) upsertBook(ctx context.Context, project *writerModel.Project, req *BookstorePublishProjectRequest) (*bookstoreModel.Book, error) {
	now := time.Now()
	status := bookstoreModel.BookStatusOngoing
	if req.PublishType == "complete" {
		status = bookstoreModel.BookStatusCompleted
	}
	projectID := project.ID.Hex()

	updates := bson.M{
		"title":          project.Title,
		"author":         project.AuthorID,
		"author_id":      project.AuthorID,
		"introduction":   pickString(req.Description, project.Summary),
		"cover":          pickString(req.CoverImage, project.CoverURL),
		"tags":           req.Tags,
		"project_id":     projectID,
		"source_type":    "writer_project",
		"status":         status,
		"price":          req.Price,
		"is_free":        req.Price <= 0,
		"published_at":   now,
		"last_update_at": now,
		"last_synced_at": now,
		"updated_at":     now,
	}

	existing, err := c.getBookByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if _, err := c.db.Collection("books").UpdateOne(ctx, bson.M{"_id": existing.ID}, bson.M{"$set": updates}); err != nil {
			return nil, err
		}
		for k, v := range updates {
			switch k {
			case "title":
				existing.Title = v.(string)
			case "author":
				existing.Author = v.(string)
			case "author_id":
				existing.AuthorID = v.(string)
			case "introduction":
				existing.Introduction = v.(string)
			case "cover":
				existing.Cover = v.(string)
			case "tags":
				existing.Tags = v.([]string)
			case "status":
				existing.Status = v.(bookstoreModel.BookStatus)
			case "price":
				existing.Price = v.(float64)
			case "is_free":
				existing.IsFree = v.(bool)
			}
		}
		existing.LastSyncedAt = &now
		existing.LastUpdateAt = &now
		if existing.PublishedAt == nil {
			existing.PublishedAt = &now
		}
		return existing, nil
	}

	book := &bookstoreModel.Book{
		Title:        project.Title,
		Author:       project.AuthorID,
		AuthorID:     project.AuthorID,
		Introduction: pickString(req.Description, project.Summary),
		Cover:        pickString(req.CoverImage, project.CoverURL),
		Tags:         req.Tags,
		ProjectID:    &projectID,
		SourceType:   "writer_project",
		Status:       status,
		Price:        req.Price,
		IsFree:       req.Price <= 0,
		PublishedAt:  &now,
		LastUpdateAt: &now,
		LastSyncedAt: &now,
	}
	book.CreatedAt = now
	book.UpdatedAt = now

	result, err := c.db.Collection("books").InsertOne(ctx, book)
	if err != nil {
		return nil, err
	}
	book.ID = result.InsertedID.(primitive.ObjectID)
	return book, nil
}

func (c *LocalBookstoreClient) rebuildProjectChapters(ctx context.Context, book *bookstoreModel.Book, projectID string, freeChapters int) error {
	chapterIDs, err := c.listChapterObjectIDs(ctx, book.ID.Hex())
	if err != nil {
		return err
	}
	if len(chapterIDs) > 0 {
		if _, err := c.db.Collection("chapter_contents").DeleteMany(ctx, bson.M{"chapter_id": bson.M{"$in": chapterIDs}}); err != nil {
			return err
		}
	}
	if _, err := c.db.Collection("chapters").DeleteMany(ctx, bson.M{"book_id": book.ID.Hex()}); err != nil {
		return err
	}

	documents, err := c.listProjectDocuments(ctx, projectID)
	if err != nil {
		return err
	}

	totalWords := int64(0)
	for idx, doc := range documents {
		paragraphs, err := c.getDocumentParagraphs(ctx, doc.ID.Hex())
		if err != nil {
			return err
		}
		content := joinPublishedParagraphs(paragraphs)
		chapterNum := idx + 1
		isFree := chapterNum <= freeChapters
		price := 0.0
		if !isFree {
			price = book.Price
		}

		projectChapterID := doc.ID.Hex()
		projectIDCopy := projectID
		chapter := &bookstoreModel.Chapter{
			BookID:           book.ID.Hex(),
			ProjectID:        &projectIDCopy,
			ProjectChapterID: &projectChapterID,
			Title:            doc.Title,
			ChapterNum:       chapterNum,
			WordCount:        len([]rune(content)),
			IsFree:           isFree,
			Price:            price,
			PublishTime:      time.Now(),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		insertResult, err := c.db.Collection("chapters").InsertOne(ctx, chapter)
		if err != nil {
			return err
		}
		chapterObjectID := insertResult.InsertedID.(primitive.ObjectID)
		totalWords += int64(chapter.WordCount)

		if err := c.replaceChapterContents(ctx, chapterObjectID, paragraphs, 1); err != nil {
			return err
		}
	}

	now := time.Now()
	_, err = c.db.Collection("books").UpdateOne(ctx, bson.M{"_id": book.ID}, bson.M{"$set": bson.M{"chapter_count": len(documents), "word_count": totalWords, "last_synced_at": now, "last_update_at": now, "updated_at": now}})
	return err
}

func (c *LocalBookstoreClient) upsertSingleChapter(ctx context.Context, book *bookstoreModel.Book, req *BookstorePublishChapterRequest) (string, error) {
	documentID, err := repository.ParseID(req.DocumentID)
	if err != nil {
		return "", fmt.Errorf("invalid document id: %w", err)
	}

	paragraphs, err := c.getDocumentParagraphs(ctx, req.DocumentID)
	if err != nil {
		return "", err
	}
	content := joinPublishedParagraphs(paragraphs)

	filter := bson.M{"book_id": book.ID.Hex(), "project_chapter_id": req.DocumentID}
	chapterDoc := bson.M{}
	err = c.db.Collection("chapters").FindOne(ctx, filter).Decode(&chapterDoc)
	if err != nil && err != mongo.ErrNoDocuments {
		return "", err
	}

	now := time.Now()
	updates := bson.M{
		"book_id":            book.ID.Hex(),
		"project_id":         req.ProjectID,
		"project_chapter_id": req.DocumentID,
		"title":              req.ChapterTitle,
		"chapter_num":        req.ChapterNumber,
		"word_count":         len([]rune(content)),
		"is_free":            req.IsFree,
		"price":              map[bool]float64{true: 0, false: book.Price}[req.IsFree],
		"publish_time":       now,
		"updated_at":         now,
	}

	var chapterObjectID primitive.ObjectID
	if err == mongo.ErrNoDocuments {
		updates["created_at"] = now
		result, err := c.db.Collection("chapters").InsertOne(ctx, updates)
		if err != nil {
			return "", err
		}
		chapterObjectID = result.InsertedID.(primitive.ObjectID)
	} else {
		chapterObjectID = chapterDoc["_id"].(primitive.ObjectID)
		if _, err := c.db.Collection("chapters").UpdateOne(ctx, bson.M{"_id": chapterObjectID}, bson.M{"$set": updates}); err != nil {
			return "", err
		}
	}

	version := 1
	if latestVersion, latestErr := c.getChapterContentVersion(ctx, chapterObjectID); latestErr == nil && latestVersion > 0 {
		version = latestVersion + 1
	}
	if err := c.replaceChapterContents(ctx, chapterObjectID, paragraphs, version); err != nil {
		return "", err
	}

	_, _ = c.db.Collection("books").UpdateOne(ctx, bson.M{"_id": book.ID}, bson.M{"$set": bson.M{"last_synced_at": now, "last_update_at": now, "updated_at": now}, "$inc": bson.M{"chapter_count": 0}})
	_ = documentID
	return chapterObjectID.Hex(), nil
}

func (c *LocalBookstoreClient) getProject(ctx context.Context, projectID string) (*writerModel.Project, error) {
	objectID, err := repository.ParseID(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}
	var project writerModel.Project
	if err := c.db.Collection("projects").FindOne(ctx, bson.M{"_id": objectID}).Decode(&project); err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *LocalBookstoreClient) getBookByProjectID(ctx context.Context, projectID string) (*bookstoreModel.Book, error) {
	var book bookstoreModel.Book
	if err := c.db.Collection("books").FindOne(ctx, bson.M{"project_id": projectID}).Decode(&book); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

func (c *LocalBookstoreClient) listProjectDocuments(ctx context.Context, projectID string) ([]*writerModel.Document, error) {
	objectID, err := repository.ParseID(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}

	cursor, err := c.db.Collection("documents").Find(ctx, bson.M{"project_id": objectID}, options.Find().SetSort(bson.D{{Key: "order", Value: 1}, {Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []*writerModel.Document
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return docs, nil
	}

	chapterDocs := make([]*writerModel.Document, 0)
	for _, doc := range docs {
		if doc.Type == writerModel.TypeChapter {
			chapterDocs = append(chapterDocs, doc)
		}
	}
	if len(chapterDocs) > 0 {
		sort.SliceStable(chapterDocs, func(i, j int) bool {
			if chapterDocs[i].Order == chapterDocs[j].Order {
				return chapterDocs[i].CreatedAt.Before(chapterDocs[j].CreatedAt)
			}
			return chapterDocs[i].Order < chapterDocs[j].Order
		})
		return chapterDocs, nil
	}
	return docs, nil
}

func (c *LocalBookstoreClient) getDocumentContent(ctx context.Context, documentID string) (string, error) {
	paragraphs, err := c.getDocumentParagraphs(ctx, documentID)
	if err != nil {
		return "", err
	}
	return joinPublishedParagraphs(paragraphs), nil
}

func (c *LocalBookstoreClient) getDocumentParagraphs(ctx context.Context, documentID string) ([]publishedParagraph, error) {
	objectID, err := repository.ParseID(documentID)
	if err != nil {
		return nil, fmt.Errorf("invalid document id: %w", err)
	}

	filter := bson.M{
		"document_id":  objectID,
		"content_type": bookstoreModel.ContentFormatTiptap,
	}
	opts := options.Find().SetSort(bson.D{{Key: "paragraph_order", Value: 1}, {Key: "created_at", Value: 1}})
	cursor, err := c.db.Collection("document_contents").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rows []struct {
		Content        string `bson:"content"`
		ContentType    string `bson:"content_type"`
		ParagraphOrder int    `bson:"paragraph_order"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, err
	}
	if len(rows) > 0 {
		paragraphs := make([]publishedParagraph, 0, len(rows))
		for idx, row := range rows {
			order := row.ParagraphOrder
			if order <= 0 {
				order = idx + 1
			}
			paragraphs = append(paragraphs, publishedParagraph{
				Order:   order,
				Content: strings.TrimSpace(row.Content),
				Format:  row.ContentType,
			})
		}
		return paragraphs, nil
	}

	var legacy struct {
		Content string `bson:"content"`
	}
	if err := c.db.Collection("document_contents").FindOne(ctx, bson.M{"document_id": objectID}).Decode(&legacy); err != nil {
		if err == mongo.ErrNoDocuments {
			return []publishedParagraph{}, nil
		}
		return nil, err
	}
	return splitPublishedParagraphs(legacy.Content), nil
}

func (c *LocalBookstoreClient) listChapterObjectIDs(ctx context.Context, bookID string) ([]primitive.ObjectID, error) {
	cursor, err := c.db.Collection("chapters").Find(ctx, bson.M{"book_id": bookID}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	ids := make([]primitive.ObjectID, 0)
	for cursor.Next(ctx) {
		var row struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}
		ids = append(ids, row.ID)
	}
	return ids, cursor.Err()
}

func (c *LocalBookstoreClient) getChapterByProjectChapterID(ctx context.Context, projectChapterID string) (*bookstoreModel.Chapter, error) {
	var chapter bookstoreModel.Chapter
	if err := c.db.Collection("chapters").FindOne(ctx, bson.M{"project_chapter_id": projectChapterID}).Decode(&chapter); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &chapter, nil
}

func pickString(primary, fallback string) string {
	if primary != "" {
		return primary
	}
	return fallback
}

func safeTime(ts *time.Time) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return *ts
}

func calculateContentHash(content string) string {
	cc := bookstoreModel.ChapterContent{Content: content}
	return cc.CalculateHash()
}

func (c *LocalBookstoreClient) getChapterContentVersion(ctx context.Context, chapterID primitive.ObjectID) (int, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "version", Value: -1}})
	var row struct {
		Version int `bson:"version"`
	}
	if err := c.db.Collection("chapter_contents").FindOne(ctx, bson.M{"chapter_id": chapterID}, opts).Decode(&row); err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}
	return row.Version, nil
}

func splitPublishedParagraphs(content string) []publishedParagraph {
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	normalized = strings.TrimSpace(normalized)
	if normalized == "" {
		return []publishedParagraph{}
	}

	blocks := strings.Split(normalized, "\n\n")
	paragraphs := make([]publishedParagraph, 0, len(blocks))
	order := 1
	for _, block := range blocks {
		paragraph := strings.TrimSpace(block)
		if paragraph == "" {
			continue
		}
		paragraphs = append(paragraphs, publishedParagraph{
			Order:   order,
			Content: paragraph,
			Format:  bookstoreModel.ContentFormatMarkdown,
		})
		order++
	}
	if len(paragraphs) == 0 {
		paragraphs = append(paragraphs, publishedParagraph{
			Order:   1,
			Content: normalized,
			Format:  bookstoreModel.ContentFormatMarkdown,
		})
	}
	return paragraphs
}

func joinPublishedParagraphs(paragraphs []publishedParagraph) string {
	if len(paragraphs) == 0 {
		return ""
	}
	parts := make([]string, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		parts = append(parts, paragraph.Content)
	}
	return strings.Join(parts, "\n\n")
}

func (c *LocalBookstoreClient) replaceChapterContents(ctx context.Context, chapterID primitive.ObjectID, paragraphs []publishedParagraph, version int) error {
	if _, err := c.db.Collection("chapter_contents").DeleteMany(ctx, bson.M{"chapter_id": chapterID}); err != nil {
		return err
	}
	if len(paragraphs) == 0 {
		return nil
	}

	now := time.Now()
	docs := make([]interface{}, 0, len(paragraphs))
	for idx, paragraph := range paragraphs {
		order := paragraph.Order
		if order <= 0 {
			order = idx + 1
		}
		content := bookstoreModel.ChapterContent{
			ID:             primitive.NewObjectID(),
			ChapterID:      chapterID,
			Content:        paragraph.Content,
			Format:         paragraph.Format,
			Version:        version,
			ParagraphOrder: order,
			WordCount:      len([]rune(paragraph.Content)),
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		content.Hash = calculateContentHash(paragraph.Content)
		docs = append(docs, content)
	}

	_, err := c.db.Collection("chapter_contents").InsertMany(ctx, docs)
	return err
}

type genericPublishEvent struct {
	eventType string
	data      interface{}
	timestamp time.Time
	source    string
}

func (e genericPublishEvent) GetEventType() string {
	return e.eventType
}

func (e genericPublishEvent) GetEventData() interface{} {
	return e.data
}

func (e genericPublishEvent) GetTimestamp() time.Time {
	return e.timestamp
}

func (e genericPublishEvent) GetSource() string {
	return e.source
}
