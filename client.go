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
	url     string
	timeout *time.Duration
}

// NewHTTPClient Creates Client
func NewHTTPClient(url string) *Client {
	return &Client{url: url}
}

// SetTimeout http timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = &timeout
}

// Send measuring to Phicus Measuring API
func (c *Client) Send(measuring Measuring) (string, error) {
	var response measuringResponse
	if err := c.send(fmt.Sprintf("%s?timestamp=%d", c.url, time.Now().UnixNano()), measuring, &response); err != nil {
		return "", err
	}
	return response.MeasuringID, nil
}

// Upload file to Phicus Measuring API
func (c *Client) Upload(file io.Reader) (string, error) {
	var response uploadResponse
	if err := c.send(c.url+"/files", file, &response); err != nil {
		return "", err
	}
	return response.FileID, nil
}

// UploadFile file to Phicus Measuring API
func (c *Client) UploadFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return c.Upload(file)
}

// Attach uploaded file to exists measuring
func (c *Client) Attach(measuringID, fileID string) error {
	request := attachmentRequest{measuringID, fileID}
	return c.send(c.url+"/attachments", request, nil)
}

type attachmentRequest struct {
	MeasuringID string `json:"measuring_id"`
	FileID      string `json:"file_id"`
}

type errorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields"`
}

type uploadResponse struct {
	FileID string `json:"file_id"`
}

type measuringResponse struct {
	MeasuringID string `json:"measuring_id"`
}

func (c *Client) send(url string, data interface{}, response interface{}) error {
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
		req, err = http.NewRequest("POST", url, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		content, jsonErr := json.Marshal(data)
		if jsonErr != nil {
			return jsonErr
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(content))
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
