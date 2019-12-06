package phicus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// Client for Phicus Measuring API
type Client struct {
	url      string
	timeout  *time.Duration
	timediff time.Duration
}

// NewClient Creates Client
func NewClient(url string, timeout *time.Duration) *Client {
	var timediff time.Duration
	if resp, err := http.Head(url); err == nil {
		if t, err := time.Parse(time.RFC3339Nano, resp.Header.Get("x-time")); err == nil {
			timediff = time.Now().Sub(t)
		}
	}
	return &Client{url: url, timediff: timediff, timeout: timeout}
}

func (c *Client) Send(key, value string) (string, error) {
	return c.sendWithParams(key, value, nil, nil, nil)
}

func (c *Client) SendWithLatLng(key, value string, lat, lng float64) (string, error) {
	return c.sendWithParams(key, value, &lat, &lng, nil)
}

func (c *Client) SendWithDisplay(key, value string, display string) (string, error) {
	return c.sendWithParams(key, value, nil, nil, &display)
}

func (c *Client) SendWithParams(key, value string, lat, lng float64, display string) (string, error) {
	return c.sendWithParams(key, value, &lat, &lng, &display)
}

func (c *Client) sendWithParams(key, value string, lat, lng *float64, display *string) (string, error) {
	m := measuring{
		Value:   value,
		Lat:     lat,
		Lng:     lng,
		Display: display,
		Time:    time.Now(),
	}
	var response struct {
		MeasuringID string `json:"measuring_id"`
	}
	if err := c.send("POST", fmt.Sprintf("%s/%s", c.url, key), m, &response); err != nil {
		return "", err
	}
	return response.MeasuringID, nil
}

// Upload file to Phicus Measuring API
func (c *Client) Upload(key string, file io.Reader) (string, error) {
	var response struct {
		FileID string `json:"file_id"`
	}
	if err := c.send("PUT", fmt.Sprintf("%s/%s", c.url, key), file, &response); err != nil {
		return "", err
	}
	return response.FileID, nil
}

// UploadFile file to Phicus Measuring API
func (c *Client) UploadFile(key, filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return c.Upload(key, file)
}

// Attach uploaded file to exists measuring
func (c *Client) Attach(key, measuringID, fileID string) error {
	request := attachmentRequest{FileID: fileID}
	return c.send("PATCH", fmt.Sprintf("%s/%s/%s", c.url, key, measuringID), request, nil)
}

type attachmentRequest struct {
	FileID string `json:"file_id"`
}

type errorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields"`
}

func (c *Client) send(method string, url string, data interface{}, response interface{}) error {
	var req *http.Request
	var err error
	if reader, ok := data.(io.Reader); ok {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", uuid.New().String())
		if err != nil {
			return err
		}
		if _, err := io.Copy(part, reader); err != nil {
			return err
		}
		if err := writer.Close(); err != nil {
			return err
		}
		req, err = http.NewRequest(method, url, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		content, jsonErr := json.Marshal(data)
		if jsonErr != nil {
			return jsonErr
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(content))
		req.Header.Set("Content-Type", "application/json")
	}

	if err != nil {
		return err
	}

	client := &http.Client{}
	if c.timeout != nil {
		client.Timeout = *c.timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		decoder := json.NewDecoder(resp.Body)
		var result errorResponse
		if err := decoder.Decode(&result); err != nil {
			return err
		}
		return errors.New(result.Error)
	}
	if response != nil {
		decoder := json.NewDecoder(resp.Body)
		return decoder.Decode(&response)
	}
	return nil
}

type measuring struct {
	Value   string    `json:"value"`
	Lat     *float64  `json:"lat"`
	Lng     *float64  `json:"lng"`
	Display *string   `json:"display"`
	Time    time.Time `json:"time"`
}
