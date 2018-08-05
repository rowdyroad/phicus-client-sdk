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

type Client struct {
	url     string
	timeout *time.Duration
}

func NewHTTPClient(url string) *Client {
	return &Client{url: url}
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = &timeout
}

func (c *Client) Send(measuring *Measuring) error {
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
	return fmt.Errorf("Error, repsonse %d", resp.Status)
}

func (c *Client) Attach(file *io.Reader) (string, error) {
	panic(errors.New("Not implemented"))
	return "", nil
}
