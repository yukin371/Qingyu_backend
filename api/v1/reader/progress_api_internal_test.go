package reader

import (
	"context"
	"net/http/httptest"
	"testing"

	"Qingyu_backend/models"
	"Qingyu_backend/service/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type recordingDevicePort struct {
	devices []*models.Device
}

func (p *recordingDevicePort) UpsertDevice(_ context.Context, device *models.Device) error {
	p.devices = append(p.devices, device)
	return nil
}

func (p *recordingDevicePort) GetByUserID(context.Context, string) ([]*models.Device, error) {
	return nil, nil
}

var _ interfaces.ReaderDevicePort = (*recordingDevicePort)(nil)

func TestProgressAPITrackDeviceSkipsInvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/api/v1/reader/progress", nil)
	ctx.Request.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120.0")

	port := &recordingDevicePort{}
	api := NewProgressAPI(nil, port)
	api.trackDevice(ctx, "not-an-object-id")

	require.Empty(t, port.devices)
}
