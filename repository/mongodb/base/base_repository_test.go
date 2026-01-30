package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/shared/types"
)

// TestBaseMongoRepository_ParseID 测试ParseID方法
func TestBaseMongoRepository_ParseID(t *testing.T) {
	// 创建测试用的repository（不需要真实连接）
	repo := &BaseMongoRepository{}

	tests := []struct {
		name    string
		input   string
		wantErr bool
		errType error
	}{
		{
			name:    "有效的ObjectID",
			input:   "507f1f77bcf86cd799439011",
			wantErr: false,
		},
		{
			name:    "空字符串",
			input:   "",
			wantErr: true,
			errType: types.ErrEmptyID,
		},
		{
			name:    "无效的hex格式",
			input:   "invalid-id-format",
			wantErr: true,
			errType: types.ErrInvalidIDFormat,
		},
		{
			name:    "长度不足",
			input:   "507f1f77bcf86cd79943901",
			wantErr: true,
			errType: types.ErrInvalidIDFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.ParseID(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				assert.True(t, result.IsZero(), "错误时应返回零值ObjectID")
			} else {
				assert.NoError(t, err)
				assert.False(t, result.IsZero(), "成功时应返回非零ObjectID")
				assert.Equal(t, tt.input, result.Hex())
			}
		})
	}
}

// TestBaseMongoRepository_ParseIDs 测试ParseIDs方法
func TestBaseMongoRepository_ParseIDs(t *testing.T) {
	repo := &BaseMongoRepository{}

	tests := []struct {
		name    string
		input   []string
		wantErr bool
	}{
		{
			name:    "全部有效的ID",
			input:   []string{"507f1f77bcf86cd799439011", "507f1f77bcf86cd799439012", "507f1f77bcf86cd799439013"},
			wantErr: false,
		},
		{
			name:    "包含无效ID",
			input:   []string{"507f1f77bcf86cd799439011", "invalid", "507f1f77bcf86cd799439013"},
			wantErr: true,
		},
		{
			name:    "空切片",
			input:   []string{},
			wantErr: false,
		},
		{
			name:    "包含空字符串",
			input:   []string{"507f1f77bcf86cd799439011", ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.ParseIDs(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if len(tt.input) == 0 {
					assert.Empty(t, result)
				} else {
					assert.Len(t, result, len(tt.input))
					// 验证转换的正确性
					for i, oid := range result {
						assert.Equal(t, tt.input[i], oid.Hex())
					}
				}
			}
		})
	}
}

// TestBaseMongoRepository_IDToHex 测试IDToHex方法
func TestBaseMongoRepository_IDToHex(t *testing.T) {
	repo := &BaseMongoRepository{}

	expectedHex := "507f1f77bcf86cd799439011"
	oid, _ := primitive.ObjectIDFromHex(expectedHex)

	result := repo.IDToHex(oid)
	assert.Equal(t, expectedHex, result)
}

// TestBaseMongoRepository_IDsToHex 测试IDsToHex方法
func TestBaseMongoRepository_IDsToHex(t *testing.T) {
	repo := &BaseMongoRepository{}

	oid1, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	oid2, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	oid3, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439013")

	oids := []primitive.ObjectID{oid1, oid2, oid3}

	result := repo.IDsToHex(oids)
	expected := []string{"507f1f77bcf86cd799439011", "507f1f77bcf86cd799439012", "507f1f77bcf86cd799439013"}

	assert.Equal(t, expected, result)
}

// TestBaseMongoRepository_IsValidID 测试IsValidID方法
func TestBaseMongoRepository_IsValidID(t *testing.T) {
	repo := &BaseMongoRepository{}

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "有效的ObjectID",
			input: "507f1f77bcf86cd799439011",
			want:  true,
		},
		{
			name:  "无效的hex格式",
			input: "invalid-id-format",
			want:  false,
		},
		{
			name:  "空字符串",
			input: "",
			want:  false,
		},
		{
			name:  "长度不足",
			input: "507f1f77bcf86cd79943901",
			want:  false,
		},
		{
			name:  "长度超出",
			input: "507f1f77bcf86cd7994390111111",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.IsValidID(tt.input)
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestBaseMongoRepository_GenerateID 测试GenerateID方法
func TestBaseMongoRepository_GenerateID(t *testing.T) {
	repo := &BaseMongoRepository{}

	// 生成两个ID，应该不同
	id1 := repo.GenerateID()
	id2 := repo.GenerateID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)

	// 验证生成的ID是有效的
	assert.True(t, repo.IsValidID(id1))
	assert.True(t, repo.IsValidID(id2))
}

// TestBaseMongoRepository_ParseID_RoundTrip 测试ID转换的往返一致性
func TestBaseMongoRepository_ParseID_RoundTrip(t *testing.T) {
	repo := &BaseMongoRepository{}

	originalHex := "507f1f77bcf86cd799439011"

	// 字符串 -> ObjectID
	oid, err := repo.ParseID(originalHex)
	require.NoError(t, err)

	// ObjectID -> 字符串
	resultHex := repo.IDToHex(oid)

	assert.Equal(t, originalHex, resultHex)
}

// TestBaseMongoRepository_ParseIDs_RoundTrip 测试批量ID转换的往返一致性
func TestBaseMongoRepository_ParseIDs_RoundTrip(t *testing.T) {
	repo := &BaseMongoRepository{}

	originalHexes := []string{"507f1f77bcf86cd799439011", "507f1f77bcf86cd799439012", "507f1f77bcf86cd799439013"}

	// 字符串 -> ObjectID
	oids, err := repo.ParseIDs(originalHexes)
	require.NoError(t, err)

	// ObjectID -> 字符串
	resultHexes := repo.IDsToHex(oids)

	assert.Equal(t, originalHexes, resultHexes)
}

// TestNewBaseMongoRepository 测试NewBaseMongoRepository构造函数
func TestNewBaseMongoRepository(t *testing.T) {
	// 这个测试需要真实的MongoDB连接，所以暂时跳过
	// 在集成测试中应该进行完整测试
	t.Skip("需要MongoDB连接，在集成测试中运行")
}

// TestBaseMongoRepository_GetCollection 测试GetCollection方法
func TestBaseMongoRepository_GetCollection(t *testing.T) {
	t.Skip("需要MongoDB连接，在集成测试中运行")
}

// 以下是需要MongoDB连接的集成测试
// 在CI/CD或本地开发环境中运行

// TestBaseMongoRepository_FindByID_Integration 集成测试FindByID
func TestBaseMongoRepository_FindByID_Integration(t *testing.T) {
	t.Skip("集成测试 - 需要MongoDB连接")
}

// TestBaseMongoRepository_UpdateByID_Integration 集成测试UpdateByID
func TestBaseMongoRepository_UpdateByID_Integration(t *testing.T) {
	t.Skip("集成测试 - 需要MongoDB连接")
}

// TestBaseMongoRepository_Find_Integration 集成测试Find
func TestBaseMongoRepository_Find_Integration(t *testing.T) {
	t.Skip("集成测试 - 需要MongoDB连接")
}

// TestBaseMongoRepository_Count_Integration 集成测试Count
func TestBaseMongoRepository_Count_Integration(t *testing.T) {
	t.Skip("集成测试 - 需要MongoDB连接")
}

// ===== 基准测试 =====

// BenchmarkParseID ParseID方法的基准测试
func BenchmarkParseID(b *testing.B) {
	repo := &BaseMongoRepository{}
	testID := "507f1f77bcf86cd799439011"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.ParseID(testID)
	}
}

// BenchmarkParseIDs ParseIDs方法的基准测试
func BenchmarkParseIDs(b *testing.B) {
	repo := &BaseMongoRepository{}
	testIDs := []string{
		"507f1f77bcf86cd799439011",
		"507f1f77bcf86cd799439012",
		"507f1f77bcf86cd799439013",
		"507f1f77bcf86cd799439014",
		"507f1f77bcf86cd799439015",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.ParseIDs(testIDs)
	}
}

// BenchmarkIDToHex IDToHex方法的基准测试
func BenchmarkIDToHex(b *testing.B) {
	repo := &BaseMongoRepository{}
	testOID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.IDToHex(testOID)
	}
}
