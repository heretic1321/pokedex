package pokeapi

import (
	"fmt"
	"github.com/heretic1321/pokedex/internal/errorhandler"
	"io"
	"net/http"
)

type Client struct {
	Base string
	http *http.Client
}

const defaultBaseUrl = "https://pokeapi.co/api/v2/"

func New(httpClient *http.Client, base string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if base == "" {
		base = defaultBaseUrl
	}
	return &Client{Base: base, http: httpClient}
}

func (c *Client) DoJSON(url string, dst any) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errorhandler.Handle(err)
	}
	res, err := c.http.Do(req)

	if err != nil {
		return nil, errorhandler.Handle(err)
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}

	if bytes, err := io.ReadAll(res.Body); err != nil {
		return nil, errorhandler.Handle(err)
	} else {
		return bytes, nil
	}
}
