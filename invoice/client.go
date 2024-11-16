package invoice

import (
	"io"
	"net/http"
)

type httpClient struct {
	c     http.Client
	token string
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Api-Token", c.token)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	return c.c.Do(req)
}

func (c *httpClient) Get(url string) (res *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *httpClient) Post(url string, contentType string, body io.Reader) (res *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)
}
