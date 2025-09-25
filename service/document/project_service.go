package document

import (
	"context"
	"errors"
	"time"

	"Qingyu_backend/global"
	model "Qingyu_backend/models/document"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProjectService struct{}

func projectCol() *mongo.Collection { return global.DB.Collection("projects") }

// EnsureIndexes 创建项目相关的 MongoDB 索引（幂等）
func (s *ProjectService) EnsureIndexes(ctx context.Context) error {
	// 获取集合
	collection := projectCol()

	// 创建索引模型
	indexModels := []mongo.IndexModel{
		{
			// 项目ID唯一索引
			Keys:    bson.M{"_id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			// 所有者ID索引，用于按所有者查询项目
			Keys: bson.M{"owner_id": 1},
		},
		{
			// 状态索引，用于按状态过滤项目
			Keys: bson.M{"status": 1},
		},
		{
			// 复合索引：所有者ID和状态，用于同时按所有者和状态过滤
			Keys: bson.M{"owner_id": 1, "status": 1},
		},
		{
			// 创建时间索引，用于排序
			Keys: bson.M{"created_at": -1},
		},
		{
			// 删除时间索引，用于软删除查询
			Keys:    bson.M{"deleted_at": 1},
			Options: options.Index().SetSparse(true), // 稀疏索引，因为只有被软删除的文档才有这个字段
		},
	}

	// 创建索引
	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return err
	}

	return nil
}

// CreateProject 创建项目
func (s *ProjectService) CreateProject(ctx context.Context, p *model.Project) (*model.Project, error) {
	if p.Name == "" || p.ID == "" {
		return nil, errors.New("invalid arguments")
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt
	p.ID = primitive.NewObjectID().Hex()
	if _, err := projectCol().InsertOne(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// GetProjectByID 根据项目id获取项目详情
func (s *ProjectService) GetProjectByID(ctx context.Context, projectID string) (*model.Project, error) {
	if projectID == "" {
		return nil, errors.New("未提供项目id")
	}
	var p model.Project
	if err := projectCol().FindOne(ctx, bson.M{"_id": projectID}).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

// ProjectList 获取项目列表(支持owner过滤/状态过滤)分页
func (s *ProjectService) GetProjectList(ctx context.Context, ownerID string, status string, limit, offset int64) ([]*model.Project, error) {
	filter := bson.M{}
	if ownerID != "" {
		filter["owner_id"] = ownerID
	}
	if status != "" {
		filter["status"] = status
	}
	opt := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(limit).SetSkip(offset)
	cur, err := projectCol().Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	list := make([]*model.Project, 0)
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateProject 更新项目（只允许改 name/description/status）
func (s *ProjectService) UpdateProjectByID(ctx context.Context, projectID, ownerID string, upd *model.Project) error {
	if upd == nil {
		return errors.New("invalid arguments")
	}
	set := bson.M{"updated_at": time.Now()}
	if upd.Name != "" {
		set["name"] = upd.Name
	}
	if upd.Description != "" {
		set["description"] = upd.Description
	}
	if upd.Status != "" {
		set["status"] = upd.Status
	}
	res, err := projectCol().UpdateOne(ctx,
		bson.M{"_id": projectID, "owner_id": ownerID}, // 强制只能改自己的
		bson.M{"$set": set})
	if res.MatchedCount == 0 {
		return errors.New("project not found")
	}
	if mongo.IsDuplicateKeyError(err) {
		return errors.New("project name duplicate")
	}
	return err
}

// 软删除（软删标记，防真删）
func (s *ProjectService) DeleteProjectByID(ctx context.Context, projectID, ownerID string) error {
	if projectID == "" || ownerID == "" {
		return errors.New("invalid arguments")
	}
	res, err := projectCol().UpdateOne(ctx,
		bson.M{"_id": projectID, "owner_id": ownerID},
		bson.M{"$set": bson.M{"deleted_at": time.Now(), "status": "deleted"}})
	if res.MatchedCount == 0 {
		return errors.New("project not found")
	}
	return err
}

// RestoreProjectByID 恢复项目（软删恢复）
func (s *ProjectService) RestoreProjectByID(ctx context.Context, projectID, ownerID string) error {
	if projectID == "" || ownerID == "" {
		return errors.New("invalid arguments")
	}
	return s.UpdateProjectByID(ctx, projectID, ownerID, &model.Project{
		Status: "active",
	})
}

// DeleteHard 硬删除（管理后台用）
func (s *ProjectService) DeleteHard(ctx context.Context, projectID string) error {
	_, err := projectCol().DeleteOne(ctx, bson.M{"_id": projectID})
	if err != nil {
		return err
	}
	return nil
}

// IsOwner 判断用户是否 owner（权限切面用）
func (s *ProjectService) IsOwner(ctx context.Context, projectID, userID string) bool {
	return projectCol().FindOne(ctx, bson.M{"_id": projectID, "owner_id": userID}).Err() == nil
}

// CreateWithRootNode 事务示例：创建工程同时初始化根节点（跨表）
func (s *ProjectService) CreateWithRootNode(ctx context.Context, p *model.Project, rootNode *model.Node) error {
	return global.MongoClient.UseSession(ctx, func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			return err
		}
		defer sc.EndSession(ctx)

		if _, err := s.CreateProject(ctx, p); err != nil {
			sc.AbortTransaction(sc)
			return err
		}
		rootNode.ProjectID = p.ID
		rootNode.TouchForCreate()
		if _, err := global.DB.Collection("nodes").InsertOne(sc, rootNode); err != nil {
			sc.AbortTransaction(sc)
			return err
		}
		return sc.CommitTransaction(sc)
	})
}
