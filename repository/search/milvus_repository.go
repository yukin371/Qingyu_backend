package search

import (
	"context"
	"fmt"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.uber.org/zap"

	"Qingyu_backend/pkg/logger"
)

// MilvusRepository Milvus 数据访问层
type MilvusRepository struct {
	client         client.Client
	collectionName string
	dimension      int
	logger         *logger.Logger
}

// NewMilvusRepository 创建 Milvus 仓储（不连接，需单独调用 Connect）
func NewMilvusRepository(collectionName string, dimension int) (*MilvusRepository, error) {
	if collectionName == "" {
		collectionName = "qingyu_books"
	}
	if dimension == 0 {
		dimension = 1536
	}

	return &MilvusRepository{
		collectionName: collectionName,
		dimension:      dimension,
		logger:         logger.Get().WithModule("milvus-repository"),
	}, nil
}

// Connect 连接到 Milvus
func (r *MilvusRepository) Connect(ctx context.Context, host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)

	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c, err := client.NewClient(connectCtx, client.Config{
		Address: addr,
	})
	if err != nil {
		r.logger.Error("Failed to connect to Milvus",
			zap.String("addr", addr),
			zap.Error(err),
		)
		return fmt.Errorf("failed to connect to Milvus at %s: %w", addr, err)
	}

	r.client = c
	r.logger.Info("Connected to Milvus", zap.String("addr", addr))
	return nil
}

// Close 关闭 Milvus 连接
func (r *MilvusRepository) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// CreateCollection 创建 qingyu_books 集合
func (r *MilvusRepository) CreateCollection(ctx context.Context) error {
	if r.client == nil {
		return fmt.Errorf("milvus client not initialized")
	}

	has, err := r.client.HasCollection(ctx, r.collectionName)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}
	if has {
		r.logger.Info("Collection already exists", zap.String("collection", r.collectionName))
		return nil
	}

	schema := &entity.Schema{
		CollectionName: r.collectionName,
		Description:    "Qingyu books vector collection",
		AutoID:         false,
		Fields: []*entity.Field{
			{
				Name:       "book_id",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: true,
				TypeParams: map[string]string{
					"max_length": "24",
				},
			},
			{
				Name:       "title_vector",
				DataType:   entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": fmt.Sprintf("%d", r.dimension),
				},
			},
			{
				Name:       "status",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "20",
				},
			},
			{
				Name:     "is_private",
				DataType: entity.FieldTypeBool,
			},
		},
	}

	if err := r.client.CreateCollection(ctx, schema, entity.DefaultShardNumber); err != nil {
		return fmt.Errorf("failed to create collection %s: %w", r.collectionName, err)
	}

	r.logger.Info("Created Milvus collection",
		zap.String("collection", r.collectionName),
		zap.Int("dimension", r.dimension),
	)

	return nil
}

// Insert 插入向量数据
func (r *MilvusRepository) Insert(ctx context.Context, vectors []map[string]interface{}) error {
	if r.client == nil {
		return fmt.Errorf("milvus client not initialized")
	}
	if len(vectors) == 0 {
		return nil
	}

	n := len(vectors)
	bookIDs := make([]string, 0, n)
	statuses := make([]string, 0, n)
	privates := make([]bool, 0, n)
	vectorData := make([][]float32, 0, n)
	dim := 0

	for _, v := range vectors {
		bookID, _ := v["book_id"].(string)
		bookIDs = append(bookIDs, bookID)

		vec := extractFloat32Slice(v["title_vector"])
		if vec == nil {
			return fmt.Errorf("invalid or missing title_vector in vector data for book_id %s", bookID)
		}
		vectorData = append(vectorData, vec)
		if dim == 0 {
			dim = len(vec)
		}

		status, _ := v["status"].(string)
		if status == "" {
			status = "ongoing"
		}
		statuses = append(statuses, status)

		isPrivate, _ := v["is_private"].(bool)
		privates = append(privates, isPrivate)
	}

	bookIDCol := entity.NewColumnVarChar("book_id", bookIDs)
	vectorCol := entity.NewColumnFloatVector("title_vector", dim, vectorData)
	statusCol := entity.NewColumnVarChar("status", statuses)
	privateCol := entity.NewColumnBool("is_private", privates)

	_, err := r.client.Insert(
		ctx,
		r.collectionName,
		"", // partition name (default)
		bookIDCol,
		vectorCol,
		statusCol,
		privateCol,
	)
	if err != nil {
		return fmt.Errorf("failed to insert vectors: %w", err)
	}

	r.logger.Info("Inserted vectors into Milvus",
		zap.String("collection", r.collectionName),
		zap.Int("count", n),
	)

	return nil
}

// Search 执行向量搜索，返回 [(id, score), ...]
func (r *MilvusRepository) Search(ctx context.Context, collection string, vectors []float32, topK int) ([]map[string]interface{}, error) {
	if r.client == nil {
		return nil, fmt.Errorf("milvus client not initialized")
	}

	if topK <= 0 {
		topK = 10
	}

	collectionName := collection
	if collectionName == "" {
		collectionName = r.collectionName
	}

	// 构建搜索向量
	searchVec := []entity.Vector{entity.FloatVector(vectors)}

	// 搜索参数: IVF_FLAT, nprobe=16
	sp, err := entity.NewIndexIvfFlatSearchParam(16)
	if err != nil {
		return nil, fmt.Errorf("failed to create search params: %w", err)
	}

	// 执行搜索
	searchResults, err := r.client.Search(
		ctx,
		collectionName,
		[]string{},          // partition names
		"",                  // expression filter
		[]string{"book_id"}, // output fields
		searchVec,
		"title_vector",     // vector field name
		entity.COSINE,      // metric type
		topK,
		sp,
	)
	if err != nil {
		return nil, fmt.Errorf("milvus search failed: %w", err)
	}

	// 解析搜索结果
	results := make([]map[string]interface{}, 0)
	for _, sr := range searchResults {
		for i := 0; i < sr.ResultCount; i++ {
			item := map[string]interface{}{
				"score": sr.Scores[i],
			}
			// 提取 book_id
			if sr.Fields != nil {
				if col := sr.Fields.GetColumn("book_id"); col != nil {
					if bookIDCol, ok := col.(*entity.ColumnVarChar); ok {
						ids := bookIDCol.Data()
						if i < len(ids) {
							item["book_id"] = ids[i]
						}
					}
				}
			}
			results = append(results, item)
		}
	}

	r.logger.Info("Milvus search completed",
		zap.String("collection", collectionName),
		zap.Int("topK", topK),
		zap.Int("results", len(results)),
	)

	return results, nil
}

// DeleteVectors 删除向量
func (r *MilvusRepository) DeleteVectors(ctx context.Context, collection string, ids []string) error {
	if r.client == nil {
		return fmt.Errorf("milvus client not initialized")
	}
	if len(ids) == 0 {
		return nil
	}

	collectionName := collection
	if collectionName == "" {
		collectionName = r.collectionName
	}

	// Delete 使用表达式字符串: book_id in ["id1", "id2"]
	expr := fmt.Sprintf("book_id in %s", formatStringList(ids))

	if err := r.client.Delete(ctx, collectionName, "", expr); err != nil {
		return fmt.Errorf("failed to delete vectors: %w", err)
	}

	r.logger.Info("Deleted vectors from Milvus",
		zap.String("collection", collectionName),
		zap.Int("count", len(ids)),
	)

	return nil
}

// DropCollection 删除集合
func (r *MilvusRepository) DropCollection(ctx context.Context, collection string) error {
	if r.client == nil {
		return fmt.Errorf("milvus client not initialized")
	}

	collectionName := collection
	if collectionName == "" {
		collectionName = r.collectionName
	}

	if err := r.client.DropCollection(ctx, collectionName); err != nil {
		return fmt.Errorf("failed to drop collection %s: %w", collectionName, err)
	}

	r.logger.Info("Dropped Milvus collection", zap.String("collection", collectionName))
	return nil
}

// HealthCheck 健康检查
func (r *MilvusRepository) HealthCheck(ctx context.Context) error {
	if r.client == nil {
		return fmt.Errorf("milvus client not initialized")
	}

	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.client.ListCollections(healthCtx)
	if err != nil {
		return fmt.Errorf("milvus health check failed: %w", err)
	}

	return nil
}

// CreateIndex 在指定字段上创建向量索引
func (r *MilvusRepository) CreateIndex(ctx context.Context, fieldName string, idx entity.Index) error {
	if r.client == nil {
		return fmt.Errorf("milvus client not initialized")
	}

	if err := r.client.CreateIndex(ctx, r.collectionName, fieldName, idx, false); err != nil {
		return fmt.Errorf("failed to create index on field %s: %w", fieldName, err)
	}

	r.logger.Info("Created Milvus index",
		zap.String("collection", r.collectionName),
		zap.String("field", fieldName),
	)

	return nil
}

// extractFloat32Slice 从 interface{} 中提取 []float32
func extractFloat32Slice(v interface{}) []float32 {
	if v == nil {
		return nil
	}

	switch vec := v.(type) {
	case []float32:
		return vec
	case []float64:
		result := make([]float32, len(vec))
		for i, f := range vec {
			result[i] = float32(f)
		}
		return result
	case []interface{}:
		result := make([]float32, len(vec))
		for i, val := range vec {
			switch f := val.(type) {
			case float64:
				result[i] = float32(f)
			case float32:
				result[i] = f
			case int:
				result[i] = float32(f)
			default:
				return nil
			}
		}
		return result
	default:
		return nil
	}
}

// formatStringList 将字符串列表格式化为 Milvus 表达式中的 in 列表
func formatStringList(ids []string) string {
	result := "["
	for i, id := range ids {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("\"%s\"", id)
	}
	result += "]"
	return result
}
