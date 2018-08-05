package phicus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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
func (c *Client) Send(measuring Measuring) error {
	content, err := json.Marshal(measuring)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(content))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

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
	if resp.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("Error, repsonse %s", resp.Status)
}

// Upload file to Phicus Measuring API
func (c *Client) Upload(file *io.Reader) (string, error) {
	panic(errors.New("Not implemented"))
}

// Attach uploaded file to exists measuring
func (c *Client) Attach(uploadID, measuringID string) (string, error) {
	panic(errors.New("Not implemented"))
}
