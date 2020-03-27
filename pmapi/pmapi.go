package pmapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client communicates with server API.
type Client struct {
	endpoint string
	client   *http.Client
}

// NewClient creates Client and communicates with server at endpoint.
func NewClient(endpoint string, client *http.Client) *Client {
	if client == nil {
		client = &http.Client{}
	}
	return &Client{
		endpoint: endpoint,
		client:   client,
	}
}

func (c *Client) url(subpath string, queries ...url.Values) string {
	dest := c.endpoint + "/api" + subpath
	if len(queries) > 0 {
		mergeQ := url.Values{}
		for _, query := range queries {
			for k, v := range query {
				mergeQ[k] = v
			}
		}
		dest += "?" + mergeQ.Encode()
	}
	return dest
}

func (c *Client) doRequest(method, url string, body interface{}) (*http.Response, error) {
	var (
		bodyR io.Reader
		err   error
	)
	if body != nil {
		br, err := c.jsonMarshal(body)
		if err != nil {
			return nil, err
		}
		bodyR = br
	}
	req, err := http.NewRequest(method, url, bodyR)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 400 {
		return res, nil
	}
	defer res.Body.Close()

	var eb ErrorBody
	if err := c.jsonUnmarshal(res.Body, &eb); err != nil {
		return nil, err
	}
	return nil, eb
}

func (c *Client) jsonMarshal(v interface{}) (io.Reader, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *Client) jsonUnmarshal(r io.Reader, v interface{}) error {
	dec := json.NewDecoder(r)
	return dec.Decode(v)
}

// Ping reports error when status code is not ok.
func (c *Client) Ping() error {
	req, err := http.NewRequest("GET", c.endpoint+"/ping", nil)
	if err != nil {
		return err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("pmapi: ping returns status code %d", res.StatusCode)
	}
	return nil
}
