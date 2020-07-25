package pinmonl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	addr   string
	client *http.Client
}

func NewClient(addr string, client *http.Client) *Client {
	return &Client{addr: addr, client: client}
}

func (c *Client) SetClient(client *http.Client) {
	c.client = client
}

func (c *Client) Info() (*ServerInfo, error) {
	dest := fmt.Sprintf("%s/api/info", c.addr)
	var out *ServerInfo
	_, err := c.get(dest, &out)
	return out, err
}

func (c *Client) Signup(user *User) (*Token, error) {
	dest := fmt.Sprintf("%s/api/signup", c.addr)
	var token *Token
	_, err := c.post(dest, user, &token)
	return token, err
}

func (c *Client) Login(user *User) (*Token, error) {
	dest := fmt.Sprintf("%s/api/login", c.addr)
	var token *Token
	_, err := c.post(dest, user, &token)
	return token, err
}

func (c *Client) MachineSignup() (*Token, error) {
	dest := fmt.Sprintf("%s/api/machine", c.addr)
	var token *Token
	_, err := c.post(dest, nil, &token)
	return token, err
}

func (c *Client) Alive() (*Token, error) {
	dest := fmt.Sprintf("%s/api/alive", c.addr)
	var token *Token
	_, err := c.post(dest, nil, &token)
	return token, err
}

func (c *Client) SharePrepare(slug string, in *Share) (*Share, error) {
	dest := fmt.Sprintf("%s/api/share/%s", c.addr, slug)
	var out *Share
	_, err := c.post(dest, in, &out)
	return out, err
}

func (c *Client) ShareDelete(slug string) error {
	dest := fmt.Sprintf("%s/api/share/%s", c.addr, slug)
	_, err := c.delete(dest, nil)
	return err
}

func (c *Client) SharePublish(slug string) (*Share, error) {
	dest := fmt.Sprintf("%s/api/share/%s/publish", c.addr, slug)
	var out *Share
	_, err := c.post(dest, nil, &out)
	return out, err
}

func (c *Client) SharetagCreate(slug string, in *Sharetag) (*Sharetag, error) {
	dest := fmt.Sprintf("%s/api/share/%s/tag", c.addr, slug)
	var out *Sharetag
	_, err := c.post(dest, in, &out)
	return out, err
}

func (c *Client) SharepinCreate(slug string, in *Sharepin) (*Sharepin, error) {
	dest := fmt.Sprintf("%s/api/share/%s/pinl", c.addr, slug)
	var out *Sharepin
	_, err := c.post(dest, in, &out)
	return out, err
}

func (c *Client) Sharing(user, slug string) (*Share, error) {
	dest := fmt.Sprintf("%s/api/sharing/%s/%s", c.addr, user, slug)
	var out *Share
	_, err := c.get(dest, &out)
	return out, err
}

func (c *Client) SharingPinlList(user, slug string, opts *PinlListOpts) ([]*Pinl, error) {
	dest := fmt.Sprintf("%s/api/sharing/%s/%s/pinl", c.addr, user, slug)
	if opts != nil {
		dest += "?" + opts.Encode()
	}

	var out []*Pinl
	_, err := c.get(dest, &out)
	return out, err
}

func (c *Client) SharingTagList(user, slug string, opts *TagListOpts) ([]*Tag, error) {
	dest := fmt.Sprintf("%s/api/sharing/%s/%s/tag", c.addr, user, slug)
	if opts != nil {
		dest += "?" + opts.Encode()
	}

	var out []*Tag
	_, err := c.get(dest, &out)
	return out, err
}

func (c *Client) Pkg(rawurl string) (*Pkg, error) {
	dest := fmt.Sprintf("%s/api/pkg/%s", c.addr, rawurl)
	var out *Pkg
	_, err := c.get(dest, &out)
	return out, err
}

func (c *Client) PkgList(rawurl string, opts *PkgListOpts) ([]*Pkg, error) {
	dest := fmt.Sprintf("%s/api/pkgs/%s", c.addr, rawurl)
	if opts != nil {
		dest += "?" + opts.Encode()
	}

	var out []*Pkg
	_, err := c.get(dest, &out)
	return out, err
}

func (c *Client) StatList(rawurl string, opts *StatListOpts) ([]*Stat, error) {
	dest := fmt.Sprintf("%s/api/stat/%s", c.addr, rawurl)
	if opts != nil {
		dest += "?" + opts.Encode()
	}

	var out []*Stat
	_, err := c.get(dest, &out)
	return out, err
}

func (c *Client) StatLatestList(rawurl string, opts *StatLatestListOpts) ([]*Stat, error) {
	dest := fmt.Sprintf("%s/api/stat/latest/%s", c.addr, rawurl)
	if opts != nil {
		dest += "?" + opts.Encode()
	}

	var out []*Stat
	_, err := c.get(dest, &out)
	return out, err
}

func (c *Client) PinlList(opts *PinlListOpts) ([]*Pinl, error) {
	dest := fmt.Sprintf("%s/api/pinl", c.addr)
	if opts != nil {
		dest += "?" + opts.Encode()
	}

	var out []*Pinl
	_, err := c.get(dest, &out)
	return out, err
}

func (c *Client) PinlClear() error {
	dest := fmt.Sprintf("%s/api/pinl", c.addr)
	_, err := c.delete(dest, nil)
	return err
}

func (c *Client) PinlCreate(in *Pinl) (*Pinl, error) {
	dest := fmt.Sprintf("%s/api/pinl", c.addr)
	var out *Pinl
	_, err := c.post(dest, in, &out)
	return out, err
}

func (c *Client) PinlDelete(pinlID string) error {
	dest := fmt.Sprintf("%s/api/pinl/%s", c.addr, pinlID)
	_, err := c.delete(dest, nil)
	return err
}

func (c *Client) get(rawurl string, out interface{}) (*http.Response, error) {
	return c.doRequest("GET", rawurl, nil, out)
}

func (c *Client) post(rawurl string, in, out interface{}) (*http.Response, error) {
	return c.doRequest("POST", rawurl, in, out)
}

func (c *Client) put(rawurl string, in, out interface{}) (*http.Response, error) {
	return c.doRequest("PUT", rawurl, in, out)
}

func (c *Client) delete(rawurl string, out interface{}) (*http.Response, error) {
	return c.doRequest("DELETE", rawurl, nil, out)
}

func (c *Client) doRequest(method, rawurl string, in, out interface{}) (*http.Response, error) {
	var reqbody io.Reader
	if in != nil {
		dec, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		reqbody = bytes.NewBuffer(dec)
	}
	req, err := http.NewRequest(method, rawurl, reqbody)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("client: err status %d, %s", resp.StatusCode, string(body))
	}
	if out != nil {
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(out)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}
