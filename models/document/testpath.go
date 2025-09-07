package document

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestBuildRelativePath(t *testing.T) {
	tests := []struct {
		name       string
		parentPath string
		nodeName   string
		want       string
		expectErr  bool
	}{
		{
			name:       "根目录创建",
			parentPath: "",
			nodeName:   "第一卷",
			want:       "第一卷",
			expectErr:  false,
		},
		{
			name:       "子目录创建",
			parentPath: "第一卷",
			nodeName:   "第一章.md",
			want:       "第一卷/第一章.md",
			expectErr:  false,
		},
		{
			name:       "深层目录创建",
			parentPath: "第一卷/章节",
			nodeName:   "第一节.md",
			want:       "第一卷/章节/第一节.md",
			expectErr:  false,
		},
		{
			name:       "空节点名",
			parentPath: "第一卷",
			nodeName:   "",
			want:       "",
			expectErr:  true,
		},
		{
			name:       "包含特殊字符的节点名",
			parentPath: "第一卷",
			nodeName:   "第一章/特殊.md",
			want:       "第一卷/第一章/特殊.md",
			expectErr:  false,
		},
		{
			name:       "节点名包含斜杠",
			parentPath: "第一卷",
			nodeName:   "第一章/第二章.md",
			want:       "第一卷/第一章/第二章.md",
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{Name: tt.nodeName}
			got, err := node.BuildRelativePath(tt.parentPath)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("BuildRelativePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoveNode(t *testing.T) {
	// 初始化测试服务
	service := &NodeService{}

	// 创建测试项目
	projectID := "test_project"

	// 测试用例
	tests := []struct {
		name         string
		nodeID       string
		newParentID  string
		expectError  bool
		expectedPath string
		setupFunc    func(t *testing.T) // 设置函数，用于创建测试数据
	}{
		{
			name:         "移动到根目录",
			nodeID:       "chapter1",
			newParentID:  "",
			expectError:  false,
			expectedPath: "第一章.md",
			setupFunc: func(t *testing.T) {
				// 创建测试节点
				_, err := nodeCol().InsertOne(context.Background(), bson.M{
					"_id":        "chapter1",
					"project_id": projectID,
					"name":       "第一章.md",
					"parent_id":  "volume1",
					"path":       "第一卷/第一章.md",
				})
				if err != nil {
					t.Fatalf("Failed to setup test node: %v", err)
				}
			},
		},
		{
			name:         "移动到子目录",
			nodeID:       "chapter1",
			newParentID:  "volume1",
			expectError:  false,
			expectedPath: "第一卷/第一章.md",
			setupFunc: func(t *testing.T) {
				// 创建测试节点
				_, err := nodeCol().InsertOne(context.Background(), bson.M{
					"_id":        "chapter1",
					"project_id": projectID,
					"name":       "第一章.md",
					"parent_id":  "",
					"path":       "第一章.md",
				})
				if err != nil {
					t.Fatalf("Failed to setup test node: %v", err)
				}

				// 创建父节点
				_, err = nodeCol().InsertOne(context.Background(), bson.M{
					"_id":        "volume1",
					"project_id": projectID,
					"name":       "第一卷",
					"parent_id":  "",
					"path":       "第一卷",
				})
				if err != nil {
					t.Fatalf("Failed to setup parent node: %v", err)
				}
			},
		},
		{
			name:         "无效节点ID",
			nodeID:       "invalid",
			newParentID:  "volume1",
			expectError:  true,
			expectedPath: "",
		},
		{
			name:         "移动到不存在的父节点",
			nodeID:       "chapter1",
			newParentID:  "nonexistent",
			expectError:  true,
			expectedPath: "",
			setupFunc: func(t *testing.T) {
				// 创建测试节点
				_, err := nodeCol().InsertOne(context.Background(), bson.M{
					"_id":        "chapter1",
					"project_id": projectID,
					"name":       "第一章.md",
					"parent_id":  "",
					"path":       "第一章.md",
				})
				if err != nil {
					t.Fatalf("Failed to setup test node: %v", err)
				}
			},
		},
		{
			name:         "移动到自己作为父节点",
			nodeID:       "chapter1",
			newParentID:  "chapter1",
			expectError:  true,
			expectedPath: "",
			setupFunc: func(t *testing.T) {
				// 创建测试节点
				_, err := nodeCol().InsertOne(context.Background(), bson.M{
					"_id":        "chapter1",
					"project_id": projectID,
					"name":       "第一章.md",
					"parent_id":  "",
					"path":       "第一章.md",
				})
				if err != nil {
					t.Fatalf("Failed to setup test node: %v", err)
				}
			},
		},
		{
			name:         "移动到自己的后代节点",
			nodeID:       "parent",
			newParentID:  "child",
			expectError:  true,
			expectedPath: "",
			setupFunc: func(t *testing.T) {
				// 创建父节点
				_, err := nodeCol().InsertOne(context.Background(), bson.M{
					"_id":        "parent",
					"project_id": projectID,
					"name":       "父目录",
					"parent_id":  "",
					"path":       "父目录",
				})
				if err != nil {
					t.Fatalf("Failed to setup parent node: %v", err)
				}

				// 创建子节点
				_, err = nodeCol().InsertOne(context.Background(), bson.M{
					"_id":        "child",
					"project_id": projectID,
					"name":       "子目录",
					"parent_id":  "parent",
					"path":       "父目录/子目录",
				})
				if err != nil {
					t.Fatalf("Failed to setup child node: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理之前的测试数据
			nodeCol().DeleteMany(context.Background(), bson.M{"project_id": projectID})

			// 如果有设置函数，执行它
			if tt.setupFunc != nil {
				tt.setupFunc(t)
			}

			err := service.MoveNode(projectID, tt.nodeID, tt.newParentID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// 验证路径是否正确更新
			var node Node
			err = nodeCol().FindOne(context.Background(), bson.M{"_id": tt.nodeID}).Decode(&node)
			if err != nil {
				t.Errorf("Failed to fetch moved node: %v", err)
				return
			}

			if node.RelativePath != tt.expectedPath {
				t.Errorf("Path mismatch: got %v, want %v", node.RelativePath, tt.expectedPath)
			}
		})
	}
}

func TestCascadePathUpdate(t *testing.T) {
	service := &NodeService{}
	projectID := "test_project"

	// 清理之前的测试数据
	nodeCol().DeleteMany(context.Background(), bson.M{"project_id": projectID})

	// 创建测试目录结构

	// 创建多层目录结构
	_, err := nodeCol().InsertMany(context.Background(), []interface{}{
		bson.M{
			"_id":        "volume1",
			"project_id": projectID,
			"name":       "第一卷",
			"parent_id":  "",
			"path":       "第一卷",
		},
		bson.M{
			"_id":        "chapter1",
			"project_id": projectID,
			"name":       "第一章",
			"parent_id":  "volume1",
			"path":       "第一卷/第一章",
		},
		bson.M{
			"_id":        "section1",
			"project_id": projectID,
			"name":       "第一节",
			"parent_id":  "chapter1",
			"path":       "第一卷/第一章/第一节",
		},
		bson.M{
			"_id":        "section2",
			"project_id": projectID,
			"name":       "第二节",
			"parent_id":  "chapter1",
			"path":       "第一卷/第一章/第二节",
		},
		bson.M{
			"_id":        "chapter2",
			"project_id": projectID,
			"name":       "第二章",
			"parent_id":  "volume1",
			"path":       "第一卷/第二章",
		},
	})
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	// 测试用例
	tests := []struct {
		name          string
		nodeID        string
		newParentID   string
		expectedPaths map[string]string // 节点ID到期望路径的映射
		expectError   bool
	}{
		{
			name:        "移动根目录到新位置",
			nodeID:      "volume1",
			newParentID: "new_parent",
			expectedPaths: map[string]string{
				"volume1":  "new_parent/第一卷",
				"chapter1": "new_parent/第一卷/第一章",
				"section1": "new_parent/第一卷/第一章/第一节",
				"section2": "new_parent/第一卷/第一章/第二节",
				"chapter2": "new_parent/第一卷/第二章",
			},
			expectError: false,
		},
		{
			name:        "移动中间层目录",
			nodeID:      "chapter1",
			newParentID: "chapter2",
			expectedPaths: map[string]string{
				"chapter1": "第一卷/第二章/第一章",
				"section1": "第一卷/第二章/第一章/第一节",
				"section2": "第一卷/第二章/第一章/第二节",
			},
			expectError: false,
		},
		{
			name:          "移动到不存在的父节点",
			nodeID:        "chapter1",
			newParentID:   "nonexistent",
			expectedPaths: map[string]string{},
			expectError:   true,
		},
		{
			name:          "移动到自己的后代节点",
			nodeID:        "volume1",
			newParentID:   "chapter1",
			expectedPaths: map[string]string{},
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置测试数据
			nodeCol().DeleteMany(context.Background(), bson.M{"project_id": projectID})
			_, err := nodeCol().InsertMany(context.Background(), []interface{}{
				bson.M{
					"_id":        "volume1",
					"project_id": projectID,
					"name":       "第一卷",
					"parent_id":  "",
					"path":       "第一卷",
				},
				bson.M{
					"_id":        "chapter1",
					"project_id": projectID,
					"name":       "第一章",
					"parent_id":  "volume1",
					"path":       "第一卷/第一章",
				},
				bson.M{
					"_id":        "section1",
					"project_id": projectID,
					"name":       "第一节",
					"parent_id":  "chapter1",
					"path":       "第一卷/第一章/第一节",
				},
				bson.M{
					"_id":        "section2",
					"project_id": projectID,
					"name":       "第二节",
					"parent_id":  "chapter1",
					"path":       "第一卷/第一章/第二节",
				},
				bson.M{
					"_id":        "chapter2",
					"project_id": projectID,
					"name":       "第二章",
					"parent_id":  "volume1",
					"path":       "第一卷/第二章",
				},
			})
			if err != nil {
				t.Fatalf("Failed to reset test data: %v", err)
			}

			err = service.MoveNode(projectID, tt.nodeID, tt.newParentID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// 验证所有相关节点的路径是否正确更新
			for nodeID, expectedPath := range tt.expectedPaths {
				var node Node
				err = nodeCol().FindOne(context.Background(), bson.M{"_id": nodeID}).Decode(&node)
				if err != nil {
					t.Errorf("Failed to fetch node %s: %v", nodeID, err)
					continue
				}

				if node.RelativePath != expectedPath {
					t.Errorf("Node %s path not updated correctly: got %v, want %v",
						nodeID, node.RelativePath, expectedPath)
				}
			}
		})
	}
}
