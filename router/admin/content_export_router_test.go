package admin

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"Qingyu_backend/models/bookstore"
	serviceinterfaces "Qingyu_backend/service/interfaces"
)

type contentBookPortStub struct{}

func (contentBookPortStub) List(context.Context, serviceinterfaces.AdminContentFilter) ([]*bookstore.Book, error) {
	return nil, nil
}

type contentChapterPortStub struct{}

func (contentChapterPortStub) List(context.Context, serviceinterfaces.AdminContentFilter) ([]*bookstore.Chapter, error) {
	return nil, nil
}

func TestRegisterContentExportRoutesRegistersProtectedRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v1 := engine.Group("/api/v1")

	RegisterContentExportRoutes(v1, contentBookPortStub{}, contentChapterPortStub{})

	paths := collectRoutePaths(engine)
	assert.True(t, paths["/api/v1/admin/content/books/export"])
	assert.True(t, paths["/api/v1/admin/content/books/export/template"])
	assert.True(t, paths["/api/v1/admin/content/chapters/export"])
	assert.True(t, paths["/api/v1/admin/content/chapters/export/template"])
}

func TestRegisterContentExportRoutesSkipsWhenPortMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v1 := engine.Group("/api/v1")

	RegisterContentExportRoutes(v1, contentBookPortStub{}, nil)

	paths := collectRoutePaths(engine)
	assert.False(t, paths["/api/v1/admin/content/books/export"])
	assert.False(t, paths["/api/v1/admin/content/chapters/export"])
}
