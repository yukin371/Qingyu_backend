package community

import socialmodels "Qingyu_backend/models/social"

// 向后兼容的类型别名 - 指向 social 模块

// Like 点赞模型
type Like = socialmodels.Like

// 常量导出
const (
	LikeTargetTypeBook    = socialmodels.LikeTargetTypeBook
	LikeTargetTypeComment = socialmodels.LikeTargetTypeComment
	LikeTargetTypeChapter = socialmodels.LikeTargetTypeChapter
)
