package pmapi

import (
	"context"
	"fmt"

	"github.com/pinmonl/pinmonl/handler/api/apibody"
)

// Tag defines the request and response structure.
type Tag apibody.Tag

// ListTags retrieves Tags from server.
func (c *Client) ListTags(_ context.Context, pageOpts PageOpts) ([]Tag, error) {
	res, err := c.doRequest("GET", c.url("/tag", pageOpts.Query()), nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var out []Tag
	if err := c.jsonUnmarshal(res.Body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTagPageInfo retrieves PageInfo of Tag.
func (c *Client) GetTagPageInfo(_ context.Context) (*PageInfo, error) {
	res, err := c.doRequest("GET", c.url("/tag/page-info"), nil)
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

// FindTagByName retrieves Tag by name.
func (c *Client) FindTagByName(_ context.Context, in *Tag) error {
	res, err := c.doRequest("GET", c.url("/tag/"+in.Name), nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var out Tag
	if err := c.jsonUnmarshal(res.Body, &out); err != nil {
		return err
	}
	*in = out
	return nil
}

// CreateTag creates Tag by input.
func (c *Client) CreateTag(_ context.Context, in *Tag) error {
	res, err := c.doRequest("POST", c.url("/tag"), in)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var out Tag
	if err := c.jsonUnmarshal(res.Body, &out); err != nil {
		return err
	}
	*in = out
	return nil
}

// UpdateTag updates Tag by input.
func (c *Client) UpdateTag(_ context.Context, in *Tag) error {
	res, err := c.doRequest("PUT", c.url("/tag/"+in.ID), in)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var out Tag
	if err := c.jsonUnmarshal(res.Body, &out); err != nil {
		return err
	}
	*in = out
	return nil
}

// DeleteTag deletes Tag by ID.
func (c *Client) DeleteTag(_ context.Context, in *Tag) error {
	res, err := c.doRequest("DELETE", c.url("/tag/"+in.ID), nil)
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
