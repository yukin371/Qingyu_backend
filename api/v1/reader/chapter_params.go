package reader

// GetChapterContentParams 获取章节内容参数
type GetChapterContentParams struct {
	BookID    string `uri:"bookId" binding:"required" json:"bookId"`
	ChapterID string `uri:"chapterId" binding:"required" json:"chapterId"`
}

// GetChapterByNumberParams 根据章节号获取内容参数
type GetChapterByNumberParams struct {
	BookID     string `uri:"bookId" binding:"required" json:"bookId"`
	ChapterNum int    `uri:"chapterNum" binding:"required,min=1" json:"chapterNum"`
}

// GetNextChapterParams 获取下一章参数
type GetNextChapterParams struct {
	BookID    string `uri:"bookId" binding:"required" json:"bookId"`
	ChapterID string `uri:"chapterId" binding:"required" json:"chapterId"`
}

// GetPreviousChapterParams 获取上一章参数
type GetPreviousChapterParams struct {
	BookID    string `uri:"bookId" binding:"required" json:"bookId"`
	ChapterID string `uri:"chapterId" binding:"required" json:"chapterId"`
}

// GetChapterListParams 获取章节目录参数
type GetChapterListParams struct {
	BookID string `uri:"bookId" binding:"required" json:"bookId"`
	Page   int    `form:"page" binding:"omitempty,min=1" json:"page"`
	Size   int    `form:"size" binding:"omitempty,min=1,max=100" json:"size"`
}

// GetChapterInfoParams 获取章节信息参数
type GetChapterInfoParams struct {
	ChapterID string `uri:"chapterId" binding:"required" json:"chapterId"`
}
