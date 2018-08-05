package phicus

import (
	"bytes"
	"log"
	"testing"
)

func TestMain(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	c := NewHTTPClient("http://localhost:10088")
	measuring, _ := NewMeasuring("test", "test", nil, nil, nil, nil)
	measuringID, err := c.Send(*measuring)
	if err != nil {
		t.Error(err)
	}
	log.Println("MeasuringID:", measuringID)

	fileID, err := c.Upload(bytes.NewBufferString("hello"))
	if err != nil {
		t.Error(err)
	}
	log.Println("FileID:", fileID)

	if err := c.Attach(measuringID, fileID); err != nil {
		t.Error(err)
	}

}
