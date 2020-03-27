package pmapi

import (
	"context"
	"fmt"

	"github.com/pinmonl/pinmonl/handler/api/pinl"
)

// Pinl defines the request and response structure.
type Pinl pinl.Body

// ListPinls retrieves Pinls from server.
func (c *Client) ListPinls(_ context.Context, pageOpts PageOpts) ([]Pinl, error) {
	res, err := c.doRequest("GET", c.url("/pinl", pageOpts.Query()), nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var out []Pinl
	if err := c.jsonUnmarshal(res.Body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetPinlPageInfo retrieves PageInfo of Pinl.
func (c *Client) GetPinlPageInfo(_ context.Context) (*PageInfo, error) {
	res, err := c.doRequest("GET", c.url("/pinl/page-info"), nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var out PageInfo
	if err := c.jsonUnmarshal(res.Body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FindPinlByID retrieves Pinl by ID.
func (c *Client) FindPinlByID(_ context.Context, in *Pinl) error {
	if in == nil {
		return fmt.Errorf("model cannot be empty")
	}
	res, err := c.doRequest("GET", c.url("/pinl/"+in.ID), nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var out Pinl
	if err := c.jsonUnmarshal(res.Body, &out); err != nil {
		return err
	}
	*in = out
	return nil
}

// CreatePinl creates Pinl by input.
func (c *Client) CreatePinl(_ context.Context, in *Pinl) error {
	res, err := c.doRequest("POST", c.url("/pinl"), in)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var out Pinl
	if err := c.jsonUnmarshal(res.Body, &out); err != nil {
		return err
	}
	*in = out
	return nil
}

// UpdatePinl updates Pinl by input.
func (c *Client) UpdatePinl(_ context.Context, in *Pinl) error {
	if in == nil {
		return fmt.Errorf("model cannot be empty")
	}
	res, err := c.doRequest("PUT", c.url("/pinl/"+in.ID), in)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var out Pinl
	if err := c.jsonUnmarshal(res.Body, &out); err != nil {
		return err
	}
	*in = out
	return nil
}

// DeletePinl deletes Pinl by ID.
func (c *Client) DeletePinl(_ context.Context, in *Pinl) error {
	if in == nil {
		return fmt.Errorf("model cannot be empty")
	}
	res, err := c.doRequest("DELETE", c.url("/pinl/"+in.ID), nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	switch {
	case res.StatusCode >= 200 && res.StatusCode < 300:
		return nil
	default:
		return fmt.Errorf("some error: [%d] %v", res.StatusCode, res.Body)
	}
}
