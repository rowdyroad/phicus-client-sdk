package phicus

import (
	"testing"

	"github.com/influxdb/pkg/testing/assert"
)

func TestSimple(t *testing.T) {
	c := NewClient("http://localhost:9050", nil)

	measuringID, err := c.Send("send", "send")
	assert.NoError(t, err)
	assert.Equal(t, 64, len(measuringID))

	measuringID, err = c.SendWithLatLng("send", "send", 55.1, 37.1)
	assert.NoError(t, err)
	assert.Equal(t, 64, len(measuringID))

	measuringID, err = c.SendWithDisplay("send", "send", "hello")
	assert.NoError(t, err)
	assert.Equal(t, 64, len(measuringID))

	measuringID, err = c.SendWithParams("send", "send", 55.1, 37.1, "hello")
	assert.NoError(t, err)
	assert.Equal(t, 64, len(measuringID))
}

func TestUpload(t *testing.T) {
	c := NewClient("http://localhost:9050", nil)
	measuringID, err := c.Send("send", "send")
	assert.NoError(t, err)
	assert.Equal(t, 64, len(measuringID))

	fileID, err := c.UploadFile("send", "photo.jpg")
	assert.NoError(t, err)
	assert.Equal(t, 64, len(fileID))

	err = c.Attach("send", measuringID, fileID)
	assert.NoError(t, err)
}
