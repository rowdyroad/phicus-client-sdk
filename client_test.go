package phicus

import (
	"bytes"
	"log"
	"testing"
)

func TestMain(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	c := NewHTTPClient("http://localhost:10080/api/measurings")
	measuringID, err := c.Send(NewMeasuring("test", "test", nil, nil, nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	log.Println("MeasuringID:", measuringID)

	fileID, err := c.Upload(bytes.NewBufferString("hello"))
	if err != nil {
		t.Fatal(err)
	}
	log.Println("FileID:", fileID)

	if err := c.Attach(measuringID, fileID); err != nil {
		t.Fatal(err)
	}

}
