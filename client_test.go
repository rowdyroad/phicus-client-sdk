package phicus

import (
	"bytes"
	"log"
	"testing"
)

func TestMain(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	c := NewHTTPClient("http://localhost/api/measurings")
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

func TestAutoAttach(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	c := NewHTTPClient("http://localhost/api/measurings")
	fileID, err := c.Upload(bytes.NewBufferString("hello"))
	if err != nil {
		t.Fatal(err)
	}
	log.Println("FileID:", fileID)

	measuringID, err := c.Send(NewMeasuring("test", "test", nil, nil, nil, []string{fileID}))
	if err != nil {
		t.Fatal(err)
	}
	log.Println("MeasuringID:", measuringID)

}

func TestDisplay(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	c := NewHTTPClient("http://localhost/api/measurings")
	fileID, err := c.UploadFile("./photo.jpg")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("FileID:", fileID)

	display := "<img src=\"{{UPLOAD}}/" + fileID + "\"/>"

	measuringID, err := c.Send(NewMeasuring("test", "test", nil, nil, &display, []string{fileID}))
	if err != nil {
		t.Fatal(err)
	}
	log.Println("MeasuringID:", measuringID)

}
